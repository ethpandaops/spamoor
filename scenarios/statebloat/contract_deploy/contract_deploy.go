package contractdeploy

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethpandaops/spamoor/scenarios/statebloat/contract_deploy/contract"
	"github.com/ethpandaops/spamoor/scenariotypes"
	"github.com/ethpandaops/spamoor/spamoor"
)

type ScenarioOptions struct {
	MaxPending      uint64 `yaml:"max_pending"`
	MaxWallets      uint64 `yaml:"max_wallets"`
	BaseFee         uint64 `yaml:"base_fee"`
	TipFee          uint64 `yaml:"tip_fee"`
	ClientGroup     string `yaml:"client_group"`
	MaxTransactions uint64 `yaml:"max_transactions"`
}

// ContractDeployment tracks a deployed contract with its deployer info
type ContractDeployment struct {
	ContractAddress string `json:"contract_address"`
	PrivateKey      string `json:"private_key"`
}

// PendingTransaction tracks a transaction until it's mined
type PendingTransaction struct {
	TxHash     common.Hash
	PrivateKey *ecdsa.PrivateKey
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	// Cached chain ID to avoid repeated RPC calls
	chainID      *big.Int
	chainIDOnce  sync.Once
	chainIDError error

	// Transaction tracking
	pendingTxs      map[common.Hash]*PendingTransaction
	pendingTxsMutex sync.RWMutex

	// Results tracking
	deployedContracts []ContractDeployment
	contractsMutex    sync.Mutex
}

var ScenarioName = "contract-deploy"
var ScenarioDefaultOptions = ScenarioOptions{
	MaxPending:      10,
	MaxWallets:      0,
	BaseFee:         20,
	TipFee:          2,
	ClientGroup:     "default",
	MaxTransactions: 0,
}
var ScenarioDescriptor = scenariotypes.ScenarioDescriptor{
	Name:           ScenarioName,
	Description:    "Deploy contracts to create state bloat",
	DefaultOptions: ScenarioDefaultOptions,
	NewScenario:    newScenario,
}

func newScenario(logger logrus.FieldLogger) scenariotypes.Scenario {
	return &Scenario{
		logger:     logger.WithField("scenario", ScenarioName),
		pendingTxs: make(map[common.Hash]*PendingTransaction),
	}
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use in transactions (in gwei)")
	flags.Uint64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use in transactions (in gwei)")
	flags.StringVar(&s.options.ClientGroup, "client-group", ScenarioDefaultOptions.ClientGroup, "Client group to use for sending transactions")
	flags.Uint64Var(&s.options.MaxTransactions, "max-transactions", ScenarioDefaultOptions.MaxTransactions, "Maximum number of transactions to send (0 for unlimited)")
	return nil
}

func (s *Scenario) Init(walletPool *spamoor.WalletPool, config string) error {
	s.walletPool = walletPool

	if config != "" {
		err := yaml.Unmarshal([]byte(config), &s.options)
		if err != nil {
			return fmt.Errorf("failed to unmarshal config: %w", err)
		}
	}

	if s.options.MaxWallets > 0 {
		walletPool.SetWalletCount(s.options.MaxWallets)
	} else {
		walletPool.SetWalletCount(1000)
	}

	return nil
}

func (s *Scenario) Config() string {
	yamlBytes, _ := yaml.Marshal(&s.options)
	return string(yamlBytes)
}

// getChainID caches the chain ID to avoid repeated RPC calls
func (s *Scenario) getChainID(ctx context.Context) (*big.Int, error) {
	s.chainIDOnce.Do(func() {
		client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, s.options.ClientGroup)
		if client == nil {
			s.chainIDError = fmt.Errorf("no client available for chain ID")
			return
		}
		s.chainID, s.chainIDError = client.GetChainId(ctx)
	})
	return s.chainID, s.chainIDError
}

// waitForPendingTxSlot waits until we have capacity for another transaction
func (s *Scenario) waitForPendingTxSlot(ctx context.Context) {
	for {
		s.pendingTxsMutex.RLock()
		count := len(s.pendingTxs)
		s.pendingTxsMutex.RUnlock()

		if count < int(s.options.MaxPending) {
			return
		}

		// Check and clean up confirmed transactions
		s.processPendingTransactions(ctx)
		time.Sleep(500 * time.Millisecond)
	}
}

// processPendingTransactions checks for transaction confirmations and updates state
func (s *Scenario) processPendingTransactions(ctx context.Context) {
	s.pendingTxsMutex.Lock()

	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, s.options.ClientGroup)
	if client == nil {
		s.pendingTxsMutex.Unlock()
		return
	}

	ethClient := client.GetEthClient()
	var confirmedTxs []common.Hash
	var successfulDeployments []struct {
		ContractAddress common.Address
		PrivateKey      *ecdsa.PrivateKey
		Receipt         *types.Receipt
		TxHash          common.Hash
	}

	for txHash, pendingTx := range s.pendingTxs {
		receipt, err := ethClient.TransactionReceipt(ctx, txHash)
		if err != nil {
			// Transaction still pending or error retrieving receipt
			continue
		}

		confirmedTxs = append(confirmedTxs, txHash)

		// Process successful deployment
		if receipt.Status == 1 && receipt.ContractAddress != (common.Address{}) {
			successfulDeployments = append(successfulDeployments, struct {
				ContractAddress common.Address
				PrivateKey      *ecdsa.PrivateKey
				Receipt         *types.Receipt
				TxHash          common.Hash
			}{
				ContractAddress: receipt.ContractAddress,
				PrivateKey:      pendingTx.PrivateKey,
				Receipt:         receipt,
				TxHash:          txHash,
			})
		}
	}

	// Remove confirmed transactions from pending map
	for _, txHash := range confirmedTxs {
		delete(s.pendingTxs, txHash)
	}

	s.pendingTxsMutex.Unlock()

	// Process successful deployments after releasing the lock
	for _, deployment := range successfulDeployments {
		s.recordDeployedContract(deployment.ContractAddress, deployment.PrivateKey, deployment.Receipt, deployment.TxHash)
	}
}

// recordDeployedContract records a successfully deployed contract
func (s *Scenario) recordDeployedContract(contractAddress common.Address, privateKey *ecdsa.PrivateKey, receipt *types.Receipt, txHash common.Hash) {
	s.contractsMutex.Lock()
	defer s.contractsMutex.Unlock()

	// Keep the JSON structure simple - only contract address and private key
	deployment := ContractDeployment{
		ContractAddress: contractAddress.Hex(),
		PrivateKey:      fmt.Sprintf("0x%x", crypto.FromECDSA(privateKey)),
	}

	s.deployedContracts = append(s.deployedContracts, deployment)

	// Calculate detailed information for logging
	walletAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	// Get the actual transaction to calculate real bytecode size
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, s.options.ClientGroup)
	var bytecodeSize int = 24564 // Default fallback
	var gasPerByte float64

	if client != nil {
		tx, _, err := client.GetEthClient().TransactionByHash(context.Background(), txHash)
		if err == nil {
			bytecodeSize = len(tx.Data())
		}
	}

	// Calculate gas per byte
	gasPerByte = float64(receipt.GasUsed) / float64(max(bytecodeSize, 1))

	// Log with detailed information
	s.logger.WithFields(logrus.Fields{
		"block_number":     receipt.BlockNumber.Uint64(),
		"bytecode_size":    bytecodeSize,
		"contract_address": deployment.ContractAddress,
		"gas_per_byte":     fmt.Sprintf("%.2f", gasPerByte),
		"gas_used":         receipt.GasUsed,
		"total_contracts":  len(s.deployedContracts),
		"tx_hash":          txHash.Hex(),
		"wallet_address":   walletAddress.Hex(),
	}).Info("Contract successfully deployed and recorded")

	// Save the deployments.json file each time a contract is confirmed
	if err := s.saveDeploymentsMapping(); err != nil {
		s.logger.Warnf("Failed to save deployments.json: %v", err)
	}
}

// Helper function for max calculation
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// saveDeploymentsMapping creates/updates deployments.json with private key to contract address mapping
func (s *Scenario) saveDeploymentsMapping() error {
	// Create a map from private key to array of contract addresses
	deploymentMap := make(map[string][]string)

	for _, contract := range s.deployedContracts {
		privateKey := contract.PrivateKey
		contractAddr := contract.ContractAddress
		deploymentMap[privateKey] = append(deploymentMap[privateKey], contractAddr)
	}

	// Create or overwrite the deployments.json file
	deploymentsFile, err := os.Create("deployments.json")
	if err != nil {
		return fmt.Errorf("failed to create deployments.json file: %w", err)
	}
	defer deploymentsFile.Close()

	// Write the mapping as JSON with pretty formatting
	encoder := json.NewEncoder(deploymentsFile)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(deploymentMap)
	if err != nil {
		return fmt.Errorf("failed to write deployments.json: %w", err)
	}

	return nil
}

func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Infof("starting scenario: %s", ScenarioName)
	defer s.logger.Infof("scenario %s finished.", ScenarioName)

	// Cache chain ID at startup
	chainID, err := s.getChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}

	s.logger.Infof("Chain ID: %s", chainID.String())

	txIdxCounter := uint64(0)
	totalTxCount := atomic.Uint64{}

	for {
		// Check if we've reached max transactions
		if s.options.MaxTransactions > 0 && txIdxCounter >= s.options.MaxTransactions {
			s.logger.Infof("reached maximum number of transactions (%d)", s.options.MaxTransactions)
			break
		}

		// Wait for available slot
		s.waitForPendingTxSlot(ctx)

		// Send a single transaction
		err := s.sendTransaction(ctx, txIdxCounter)
		if err != nil {
			s.logger.Warnf("failed to send transaction %d: %v", txIdxCounter, err)
			time.Sleep(1 * time.Second)
			continue
		}

		txIdxCounter++
		totalTxCount.Add(1)

		// Process pending transactions periodically
		if txIdxCounter%10 == 0 {
			s.processPendingTransactions(ctx)

			s.contractsMutex.Lock()
			contractCount := len(s.deployedContracts)
			s.contractsMutex.Unlock()

			s.logger.Infof("Progress: sent %d txs, deployed %d contracts", txIdxCounter, contractCount)
		}

		// Small delay to prevent overwhelming the RPC
		time.Sleep(100 * time.Millisecond)
	}

	// Wait for all pending transactions to complete
	s.logger.Info("Waiting for remaining transactions to complete...")
	for {
		s.processPendingTransactions(ctx)

		s.pendingTxsMutex.RLock()
		pendingCount := len(s.pendingTxs)
		s.pendingTxsMutex.RUnlock()

		if pendingCount == 0 {
			break
		}

		s.logger.Infof("Waiting for %d pending transactions...", pendingCount)
		time.Sleep(2 * time.Second)
	}

	// Log final summary
	s.contractsMutex.Lock()
	totalContracts := len(s.deployedContracts)
	s.contractsMutex.Unlock()

	s.logger.WithFields(logrus.Fields{
		"total_txs":       totalTxCount.Load(),
		"total_contracts": totalContracts,
	}).Info("All transactions completed")

	return nil
}

// sendTransaction sends a single contract deployment transaction
func (s *Scenario) sendTransaction(ctx context.Context, txIdx uint64) error {
	maxRetries := 3

	for attempt := 0; attempt < maxRetries; attempt++ {
		err := s.attemptTransaction(ctx, txIdx, attempt)
		if err == nil {
			return nil
		}

		// Check if it's a base fee error
		if strings.Contains(err.Error(), "max fee per gas less than block base fee") {
			s.logger.Warnf("Transaction %d base fee too low, adjusting fees and retrying (attempt %d/%d)",
				txIdx, attempt+1, maxRetries)

			// Update fees based on current network conditions
			if updateErr := s.updateDynamicFees(ctx); updateErr != nil {
				s.logger.Warnf("Failed to update dynamic fees: %v", updateErr)
			}

			time.Sleep(time.Duration(attempt+1) * 500 * time.Millisecond) // Exponential backoff
			continue
		}

		// For other errors, return immediately
		return err
	}

	return fmt.Errorf("failed to send transaction after %d attempts", maxRetries)
}

// updateDynamicFees queries the network and updates base fee and tip fee
func (s *Scenario) updateDynamicFees(ctx context.Context) error {
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, s.options.ClientGroup)
	if client == nil {
		return fmt.Errorf("no client available")
	}

	ethClient := client.GetEthClient()

	// Get the latest block to check current base fee
	latestBlock, err := ethClient.BlockByNumber(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to get latest block: %w", err)
	}

	if latestBlock.BaseFee() != nil {
		// Convert base fee from wei to gwei
		currentBaseFeeGwei := new(big.Int).Div(latestBlock.BaseFee(), big.NewInt(1000000000))

		// Set new base fee to current base fee + 20% buffer
		newBaseFeeGwei := new(big.Int).Mul(currentBaseFeeGwei, big.NewInt(120))
		newBaseFeeGwei = new(big.Int).Div(newBaseFeeGwei, big.NewInt(100))

		// Ensure minimum increase of 5 gwei
		minIncrease := big.NewInt(5)
		if newBaseFeeGwei.Cmp(new(big.Int).Add(big.NewInt(int64(s.options.BaseFee)), minIncrease)) < 0 {
			newBaseFeeGwei = new(big.Int).Add(big.NewInt(int64(s.options.BaseFee)), minIncrease)
		}

		s.options.BaseFee = newBaseFeeGwei.Uint64()

		// Also increase tip fee slightly to ensure competitive priority
		if s.options.TipFee+1 > 3 {
			s.options.TipFee = s.options.TipFee + 1
		} else {
			s.options.TipFee = 3 // Minimum 3 gwei tip
		}

		s.logger.Infof("Updated dynamic fees - Base fee: %d gwei, Tip fee: %d gwei (network base fee: %s gwei)",
			s.options.BaseFee, s.options.TipFee, currentBaseFeeGwei.String())
	}

	return nil
}

// attemptTransaction makes a single attempt to send a transaction
func (s *Scenario) attemptTransaction(ctx context.Context, txIdx uint64, attempt int) error {
	// Get client and wallet
	client := s.walletPool.GetClient(spamoor.SelectClientRoundRobin, 0, s.options.ClientGroup)
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletRoundRobin, 0)

	if client == nil {
		return fmt.Errorf("no client available")
	}
	if wallet == nil {
		return fmt.Errorf("no wallet available")
	}

	// Get cached chain ID
	chainID, err := s.getChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Get current nonce for this wallet
	addr := crypto.PubkeyToAddress(wallet.GetPrivateKey().PublicKey)
	nonce, err := client.GetEthClient().PendingNonceAt(ctx, addr)
	if err != nil {
		return fmt.Errorf("failed to get nonce for %s: %w", addr.Hex(), err)
	}

	// Create transaction auth
	auth, err := bind.NewKeyedTransactorWithChainID(wallet.GetPrivateKey(), chainID)
	if err != nil {
		return fmt.Errorf("failed to create auth: %w", err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = 5200000 // Fixed gas limit for contract deployment

	// Set EIP-1559 fee parameters
	if s.options.BaseFee > 0 {
		auth.GasFeeCap = new(big.Int).Mul(big.NewInt(int64(s.options.BaseFee)), big.NewInt(1000000000))
	}
	if s.options.TipFee > 0 {
		auth.GasTipCap = new(big.Int).Mul(big.NewInt(int64(s.options.TipFee)), big.NewInt(1000000000))
	}

	// Generate random salt for unique contract
	salt := make([]byte, 32)
	_, err = rand.Read(salt)
	if err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}
	saltInt := new(big.Int).SetBytes(salt)

	// Deploy the contract
	ethClient := client.GetEthClient()
	if ethClient == nil {
		return fmt.Errorf("failed to get eth client")
	}

	_, tx, _, err := contract.DeployContract(auth, ethClient, saltInt)
	if err != nil {
		return fmt.Errorf("failed to deploy contract: %w", err)
	}

	// Track pending transaction
	pendingTx := &PendingTransaction{
		TxHash:     tx.Hash(),
		PrivateKey: wallet.GetPrivateKey(),
	}

	s.pendingTxsMutex.Lock()
	s.pendingTxs[tx.Hash()] = pendingTx
	s.pendingTxsMutex.Unlock()

	s.logger.WithFields(logrus.Fields{
		"tx_hash":       tx.Hash().Hex(),
		"nonce":         nonce,
		"base_fee_gwei": s.options.BaseFee,
		"tip_fee_gwei":  s.options.TipFee,
		"attempt":       attempt + 1,
	}).Info("Transaction sent")

	return nil
}
