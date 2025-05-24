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
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenarios/statebloat/contract_deploy/contract"
	"github.com/ethpandaops/spamoor/scenariotypes"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
)

const (
	// GAS_PER_CONTRACT is the estimated gas cost for deploying one StateBloatToken contract
	GAS_PER_CONTRACT = 4970000 // ~4.97M gas per contract
	// BLOCK_GAS_LIMIT is the typical block gas limit for Ethereum mainnet
	BLOCK_GAS_LIMIT = 30000000 // 30M gas
)

type ScenarioOptions struct {
	MaxPending      uint64 `yaml:"max_pending"`
	MaxWallets      uint64 `yaml:"max_wallets"`
	Rebroadcast     uint64 `yaml:"rebroadcast"`
	BaseFee         uint64 `yaml:"base_fee"`
	TipFee          uint64 `yaml:"tip_fee"`
	ClientGroup     string `yaml:"client_group"`
	MaxTransactions uint64 `yaml:"max_transactions"` // Maximum number of transactions to send (0 for unlimited)
	BlockGasLimit   uint64 `yaml:"block_gas_limit"`  // Block gas limit for batching (0 = use default)
}

// ContractDeployment tracks a deployed contract with its deployer info
type ContractDeployment struct {
	ContractAddress string `json:"contract_address"`
	PrivateKey      string `json:"private_key"`
	TxHash          string `json:"tx_hash"`
	GasUsed         uint64 `json:"gas_used"`
	BlockNumber     uint64 `json:"block_number"`
	BytecodeSize    int    `json:"bytecode_size"`
	GasPerByte      string `json:"gas_per_byte"`
}

// PendingTransaction tracks a transaction until it's mined
type PendingTransaction struct {
	TxHash     common.Hash
	Wallet     *txbuilder.Wallet
	TxIndex    uint64
	SentAt     time.Time
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
	maxConcurrent   int

	// Results tracking
	deployedContracts []ContractDeployment
	contractsMutex    sync.Mutex

	// Nonce management per wallet
	walletNonces map[common.Address]uint64
	nonceMutex   sync.RWMutex
}

var ScenarioName = "contract-deploy"
var ScenarioDefaultOptions = ScenarioOptions{
	MaxPending:      10, // Limit concurrent transactions
	MaxWallets:      0,
	Rebroadcast:     1,
	BaseFee:         20,
	TipFee:          2,
	ClientGroup:     "default",
	MaxTransactions: 0, // Default to unlimited
	BlockGasLimit:   0, // Use default BLOCK_GAS_LIMIT
}
var ScenarioDescriptor = scenariotypes.ScenarioDescriptor{
	Name:           ScenarioName,
	Description:    "Deploy contracts to create state bloat",
	DefaultOptions: ScenarioDefaultOptions,
	NewScenario:    newScenario,
}

func newScenario(logger logrus.FieldLogger) scenariotypes.Scenario {
	return &Scenario{
		logger:        logger.WithField("scenario", ScenarioName),
		pendingTxs:    make(map[common.Hash]*PendingTransaction),
		walletNonces:  make(map[common.Address]uint64),
		maxConcurrent: 5, // Conservative default
	}
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", ScenarioDefaultOptions.Rebroadcast, "Number of seconds to wait before re-broadcasting a transaction")
	flags.Uint64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use in transactions (in gwei)")
	flags.Uint64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use in transactions (in gwei)")
	flags.StringVar(&s.options.ClientGroup, "client-group", ScenarioDefaultOptions.ClientGroup, "Client group to use for sending transactions")
	flags.Uint64Var(&s.options.MaxTransactions, "max-transactions", ScenarioDefaultOptions.MaxTransactions, "Maximum number of transactions to send (0 for unlimited)")
	flags.Uint64Var(&s.options.BlockGasLimit, "block-gas-limit", ScenarioDefaultOptions.BlockGasLimit, "Block gas limit for batching (0 = use default)")
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

	// Set concurrent transaction limit
	if s.options.MaxPending > 0 {
		s.maxConcurrent = int(s.options.MaxPending)
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
		if s.chainIDError == nil {
			s.logger.Infof("Cached chain ID: %s", s.chainID.String())
		}
	})
	return s.chainID, s.chainIDError
}

// getWalletNonce gets the current nonce for a wallet (without incrementing)
func (s *Scenario) getWalletNonce(ctx context.Context, wallet *txbuilder.Wallet, client *txbuilder.Client) (uint64, error) {
	s.nonceMutex.Lock()
	defer s.nonceMutex.Unlock()

	addr := crypto.PubkeyToAddress(wallet.GetPrivateKey().PublicKey)

	// Always refresh nonce from network to avoid gaps
	nonce, err := client.GetEthClient().PendingNonceAt(ctx, addr)
	if err != nil {
		return 0, fmt.Errorf("failed to get nonce for %s: %w", addr.Hex(), err)
	}

	// Store the current network nonce
	s.walletNonces[addr] = nonce
	s.logger.Debugf("Current nonce for wallet %s: %d", addr.Hex(), nonce)

	return nonce, nil
}

// incrementWalletNonce increments the local nonce counter after successful tx send
func (s *Scenario) incrementWalletNonce(wallet *txbuilder.Wallet) {
	s.nonceMutex.Lock()
	defer s.nonceMutex.Unlock()

	addr := crypto.PubkeyToAddress(wallet.GetPrivateKey().PublicKey)
	s.walletNonces[addr]++
	s.logger.Debugf("Incremented nonce for wallet %s to: %d", addr.Hex(), s.walletNonces[addr])
}

// waitForPendingTxSlot waits until we have capacity for another transaction
func (s *Scenario) waitForPendingTxSlot(ctx context.Context) {
	for {
		s.pendingTxsMutex.RLock()
		count := len(s.pendingTxs)
		s.pendingTxsMutex.RUnlock()

		if count < s.maxConcurrent {
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
	defer s.pendingTxsMutex.Unlock()

	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, s.options.ClientGroup)
	if client == nil {
		return
	}

	ethClient := client.GetEthClient()
	var confirmedTxs []common.Hash

	for txHash, pendingTx := range s.pendingTxs {
		receipt, err := ethClient.TransactionReceipt(ctx, txHash)
		if err != nil {
			// Transaction still pending or error retrieving receipt
			continue
		}

		confirmedTxs = append(confirmedTxs, txHash)

		// Process successful deployment
		if receipt.Status == 1 && receipt.ContractAddress != (common.Address{}) {
			s.recordDeployedContract(receipt, pendingTx)
		} else if receipt.Status != 1 {
			s.logger.Warnf("Transaction failed: %s (gas used: %d)", txHash.Hex(), receipt.GasUsed)
		}
	}

	// Remove confirmed transactions from pending map
	for _, txHash := range confirmedTxs {
		delete(s.pendingTxs, txHash)
	}

	if len(confirmedTxs) > 0 {
		s.logger.Debugf("Processed %d confirmed transactions, %d still pending",
			len(confirmedTxs), len(s.pendingTxs))
	}
}

// recordDeployedContract records a successfully deployed contract
func (s *Scenario) recordDeployedContract(receipt *types.Receipt, pendingTx *PendingTransaction) {
	s.contractsMutex.Lock()
	defer s.contractsMutex.Unlock()

	// Get the actual transaction to calculate real bytecode size
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, s.options.ClientGroup)
	var txBytes int
	var gasPerByte float64

	if client != nil {
		tx, _, err := client.GetEthClient().TransactionByHash(context.Background(), receipt.TxHash)
		if err == nil {
			txBytes = len(tx.Data())
			gasPerByte = float64(receipt.GasUsed) / float64(max(txBytes, 1))
		} else {
			// Fallback to estimated size
			txBytes = 24564 // Approximate StateBloatToken bytecode size
			gasPerByte = float64(receipt.GasUsed) / float64(txBytes)
		}
	} else {
		// Fallback values
		txBytes = 24564
		gasPerByte = float64(receipt.GasUsed) / float64(txBytes)
	}

	deployment := ContractDeployment{
		ContractAddress: receipt.ContractAddress.Hex(),
		PrivateKey:      fmt.Sprintf("0x%x", crypto.FromECDSA(pendingTx.PrivateKey)),
		TxHash:          receipt.TxHash.Hex(),
		GasUsed:         receipt.GasUsed,
		BlockNumber:     receipt.BlockNumber.Uint64(),
		BytecodeSize:    txBytes,
		GasPerByte:      fmt.Sprintf("%.2f", gasPerByte),
	}

	s.deployedContracts = append(s.deployedContracts, deployment)

	s.logger.WithFields(logrus.Fields{
		"contract_address": deployment.ContractAddress,
		"tx_hash":          deployment.TxHash,
		"gas_used":         deployment.GasUsed,
		"block_number":     deployment.BlockNumber,
		"bytecode_size":    deployment.BytecodeSize,
		"gas_per_byte":     deployment.GasPerByte,
		"total_contracts":  len(s.deployedContracts),
	}).Info("Contract successfully deployed and recorded")

	// Save the deployments.json file each time a contract is confirmed
	if err := s.saveDeploymentsMapping(); err != nil {
		s.logger.Warnf("Failed to save deployments.json: %v", err)
	}
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

	s.logger.WithFields(logrus.Fields{
		"total_deployers": len(deploymentMap),
		"total_contracts": len(s.deployedContracts),
		"file":            "deployments.json",
	}).Info("Updated deployments mapping file")

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

	// Determine block gas limit to use
	blockGasLimit := s.options.BlockGasLimit
	if blockGasLimit == 0 {
		blockGasLimit = BLOCK_GAS_LIMIT
	}

	// Calculate how many contracts can fit in one block
	contractsPerBlock := int(blockGasLimit / GAS_PER_CONTRACT)
	if contractsPerBlock == 0 {
		contractsPerBlock = 1
	}

	s.logger.Infof("Chain ID: %s, Block gas limit: %d, Gas per contract: %d, Contracts per block: %d, Max concurrent: %d",
		chainID.String(), blockGasLimit, GAS_PER_CONTRACT, contractsPerBlock, s.maxConcurrent)

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

	// Save results to JSON
	s.contractsMutex.Lock()
	totalContracts := len(s.deployedContracts)
	s.contractsMutex.Unlock()

	s.logger.WithFields(logrus.Fields{
		"test":            "contract-deploy",
		"total_txs":       totalTxCount.Load(),
		"total_contracts": totalContracts,
	}).Info("All transactions completed")

	// Log final summary
	return s.logFinalSummary()
}

// sendTransaction sends a single contract deployment transaction with retry logic
func (s *Scenario) sendTransaction(ctx context.Context, txIdx uint64) error {
	maxRetries := 3
	baseGasMultiplier := 1.0

	for attempt := 0; attempt < maxRetries; attempt++ {
		err := s.attemptTransaction(ctx, txIdx, baseGasMultiplier)
		if err == nil {
			return nil
		}

		// Check if it's an underpriced error
		if strings.Contains(err.Error(), "replacement transaction underpriced") ||
			strings.Contains(err.Error(), "transaction underpriced") {
			baseGasMultiplier += 0.2 // Increase gas by 20% each retry
			s.logger.Warnf("Transaction %d underpriced, retrying with %.1fx gas (attempt %d/%d)",
				txIdx, baseGasMultiplier, attempt+1, maxRetries)
			time.Sleep(time.Duration(attempt+1) * 500 * time.Millisecond) // Exponential backoff
			continue
		}

		// For other errors, return immediately
		return err
	}

	return fmt.Errorf("failed to send transaction after %d attempts", maxRetries)
}

// attemptTransaction makes a single attempt to send a transaction
func (s *Scenario) attemptTransaction(ctx context.Context, txIdx uint64, gasMultiplier float64) error {
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
	nonce, err := s.getWalletNonce(ctx, wallet, client)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %w", err)
	}

	// Create transaction auth
	auth, err := bind.NewKeyedTransactorWithChainID(wallet.GetPrivateKey(), chainID)
	if err != nil {
		return fmt.Errorf("failed to create auth: %w", err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = GAS_PER_CONTRACT + GAS_PER_CONTRACT*5/100

	// Set EIP-1559 fee parameters with dynamic adjustment
	baseFee := s.options.BaseFee
	tipFee := s.options.TipFee

	if gasMultiplier > 1.0 {
		baseFee = uint64(float64(baseFee) * gasMultiplier)
		tipFee = uint64(float64(tipFee) * gasMultiplier)
	}

	if baseFee > 0 {
		auth.GasFeeCap = new(big.Int).Mul(big.NewInt(int64(baseFee)), big.NewInt(1000000000))
	}
	if tipFee > 0 {
		auth.GasTipCap = new(big.Int).Mul(big.NewInt(int64(tipFee)), big.NewInt(1000000000))
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

	address, tx, _, err := contract.DeployContract(auth, ethClient, saltInt)
	if err != nil {
		return fmt.Errorf("failed to deploy contract: %w", err)
	}

	// Track pending transaction
	pendingTx := &PendingTransaction{
		TxHash:     tx.Hash(),
		Wallet:     wallet,
		TxIndex:    txIdx,
		SentAt:     time.Now(),
		PrivateKey: wallet.GetPrivateKey(),
	}

	s.pendingTxsMutex.Lock()
	s.pendingTxs[tx.Hash()] = pendingTx
	s.pendingTxsMutex.Unlock()

	s.logger.WithFields(logrus.Fields{
		"tx_index":         txIdx,
		"tx_hash":          tx.Hash().Hex(),
		"expected_address": address.Hex(),
		"nonce":            nonce,
		"wallet":           crypto.PubkeyToAddress(wallet.GetPrivateKey().PublicKey).Hex(),
		"gas_limit":        auth.GasLimit,
		"base_fee_gwei":    baseFee,
		"tip_fee_gwei":     tipFee,
		"gas_multiplier":   fmt.Sprintf("%.1fx", gasMultiplier),
	}).Info("Transaction sent successfully")

	// Only increment nonce after successful send
	s.incrementWalletNonce(wallet)

	return nil
}

// logFinalSummary logs the final summary of the scenario
func (s *Scenario) logFinalSummary() error {
	s.contractsMutex.Lock()
	defer s.contractsMutex.Unlock()

	if len(s.deployedContracts) == 0 {
		s.logger.Info("No contracts were deployed")
		return nil
	}

	// Log final summary
	var totalGasUsed uint64
	for _, contract := range s.deployedContracts {
		totalGasUsed += contract.GasUsed
	}

	s.logger.WithFields(logrus.Fields{
		"total_contracts":      len(s.deployedContracts),
		"total_gas_used":       totalGasUsed,
		"avg_gas_per_contract": totalGasUsed / uint64(len(s.deployedContracts)),
		"deployments_file":     "deployments.json",
	}).Info("Final deployment summary - only deployments.json file created")

	return nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
