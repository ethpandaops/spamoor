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

// Constants for the contract deployment scenario
const (
	// EstimatedGasPerContract is the estimated gas used per contract deployment
	// This value is based on the actual gas usage of deploying the StateBloatToken contract
	// which has a bytecode size of ~23.9KB
	EstimatedGasPerContract = uint64(4949468)

	// DeploymentGasLimit is the gas limit for contract deployment transactions
	DeploymentGasLimit = uint64(4980000)

	// EstimatedDeployedBytecodeSize is the estimated size of the deployed contract bytecode
	// This is the actual bytecode size on-chain after deployment
	EstimatedDeployedBytecodeSize = 23914

	// MaxPendingMultiplier is the multiplier for calculating max pending transactions
	// We allow 2x the block capacity to handle network delays and variations
	MaxPendingMultiplier = 2

	// DefaultMaxPending is the default max pending transactions for manual mode
	// This is only used when max-transactions is set (not using adaptive max pending)
	DefaultMaxPending = 10
)

type ScenarioOptions struct {
	MaxWallets      uint64 `yaml:"max_wallets"`
	BaseFee         uint64 `yaml:"base_fee"`
	TipFee          uint64 `yaml:"tip_fee"`
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
	Timestamp  time.Time
}

// BlockDeploymentStats tracks deployment statistics per block
type BlockDeploymentStats struct {
	BlockNumber       uint64
	ContractCount     int
	TotalGasUsed      uint64
	TotalBytecodeSize int
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

	// Block-level statistics tracking
	blockStats      map[uint64]*BlockDeploymentStats
	blockStatsMutex sync.Mutex
	lastLoggedBlock uint64

	// Block monitoring for real-time logging
	blockMonitorCancel context.CancelFunc
	blockMonitorDone   chan struct{}

	// Auto-calculated max pending transactions
	maxPending uint64
}

var ScenarioName = "contract-deploy"
var ScenarioDefaultOptions = ScenarioOptions{
	MaxWallets:      0, // Use root wallet only by default
	BaseFee:         5, // Moderate base fee (5 gwei)
	TipFee:          1, // Priority fee (1 gwei)
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
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Base fee per gas in gwei")
	flags.Uint64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Tip fee per gas in gwei")
	flags.Uint64Var(&s.options.MaxTransactions, "max-transactions", ScenarioDefaultOptions.MaxTransactions, "Maximum number of transactions to send (0 = unlimited with adaptive max pending based on block gas limit)")
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
		// Use only root wallet by default for better efficiency
		// This avoids child wallet funding overhead
		walletPool.SetWalletCount(0)
	}

	// Initialize default max pending for manual mode
	s.maxPending = DefaultMaxPending

	return nil
}

func (s *Scenario) Config() string {
	yamlBytes, _ := yaml.Marshal(&s.options)
	return string(yamlBytes)
}

// getChainID caches the chain ID to avoid repeated RPC calls
func (s *Scenario) getChainID(ctx context.Context) (*big.Int, error) {
	s.chainIDOnce.Do(func() {
		client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, "")
		if client == nil {
			s.chainIDError = fmt.Errorf("no client available for chain ID")
			return
		}
		s.chainID, s.chainIDError = client.GetChainId(ctx)
	})
	return s.chainID, s.chainIDError
}

// waitForPendingTxSlot waits until we have capacity for another transaction
// This is the primary flow control mechanism to prevent overwhelming the network
func (s *Scenario) waitForPendingTxSlot(ctx context.Context) {
	for {
		s.pendingTxsMutex.RLock()
		count := len(s.pendingTxs)
		s.pendingTxsMutex.RUnlock()

		if count < int(s.maxPending) {
			return
		}

		// Check and clean up confirmed transactions
		s.processPendingTransactions(ctx)
		time.Sleep(1 * time.Second)
	}
}

// processPendingTransactions checks for transaction confirmations and updates state
// This function is critical for maintaining accurate pending count and flow control
func (s *Scenario) processPendingTransactions(ctx context.Context) {
	s.pendingTxsMutex.Lock()

	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, "")
	if client == nil {
		s.pendingTxsMutex.Unlock()
		return
	}

	ethClient := client.GetEthClient()
	var confirmedTxs []common.Hash
	var timedOutTxs []common.Hash
	var successfulDeployments []struct {
		ContractAddress common.Address
		PrivateKey      *ecdsa.PrivateKey
		Receipt         *types.Receipt
		TxHash          common.Hash
	}

	for txHash, pendingTx := range s.pendingTxs {
		// Check if transaction is too old (1 minute timeout)
		if time.Since(pendingTx.Timestamp) > 1*time.Minute {
			s.logger.Warnf("Transaction %s timed out after 1 minute, removing from pending", txHash.Hex())
			timedOutTxs = append(timedOutTxs, txHash)
			continue
		}

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

	// Remove timed out transactions from pending map
	for _, txHash := range timedOutTxs {
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

	// Get the actual deployed contract bytecode size
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, "")
	var bytecodeSize int = EstimatedDeployedBytecodeSize

	if client != nil {
		// Get the actual deployed bytecode size using eth_getCode
		contractCode, err := client.GetEthClient().CodeAt(context.Background(), contractAddress, nil)
		if err == nil {
			bytecodeSize = len(contractCode)
		}
	}

	blockNumber := receipt.BlockNumber.Uint64()

	// Debug logging for block tracking
	s.logger.WithFields(logrus.Fields{
		"tx_block":        blockNumber,
		"existing_blocks": len(s.blockStats),
	}).Debug("Recording contract deployment")

	// Update block-level statistics
	s.blockStatsMutex.Lock()
	defer s.blockStatsMutex.Unlock()

	if s.blockStats == nil {
		s.blockStats = make(map[uint64]*BlockDeploymentStats)
	}

	// Create or update current block stats (removed the old logging logic)
	if s.blockStats[blockNumber] == nil {
		s.blockStats[blockNumber] = &BlockDeploymentStats{
			BlockNumber: blockNumber,
		}
		s.logger.WithField("block_number", blockNumber).Debug("Created new block stats")
	}

	blockStat := s.blockStats[blockNumber]
	blockStat.ContractCount++
	blockStat.TotalGasUsed += receipt.GasUsed
	blockStat.TotalBytecodeSize += bytecodeSize

	s.logger.WithFields(logrus.Fields{
		"block_number":       blockNumber,
		"contracts_in_block": blockStat.ContractCount,
		"gas_used":           blockStat.TotalGasUsed,
		"bytecode_size":      blockStat.TotalBytecodeSize,
	}).Debug("Updated block stats")

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

// startBlockMonitor starts a background goroutine that monitors for new blocks
// and logs block deployment summaries immediately when blocks are mined
func (s *Scenario) startBlockMonitor(ctx context.Context) {
	monitorCtx, cancel := context.WithCancel(ctx)
	s.blockMonitorCancel = cancel
	s.blockMonitorDone = make(chan struct{})

	go func() {
		defer close(s.blockMonitorDone)

		client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, "")
		if client == nil {
			s.logger.Warn("No client available for block monitoring")
			return
		}

		ethClient := client.GetEthClient()
		ticker := time.NewTicker(2 * time.Second) // Poll every 2 seconds
		defer ticker.Stop()

		for {
			select {
			case <-monitorCtx.Done():
				return
			case <-ticker.C:
				// Get current block number
				latestBlock, err := ethClient.BlockByNumber(monitorCtx, nil)
				if err != nil {
					s.logger.WithError(err).Debug("Failed to get latest block for monitoring")
					continue
				}

				currentBlockNumber := latestBlock.Number().Uint64()

				// Log any completed blocks that haven't been logged yet
				s.blockStatsMutex.Lock()
				for bn := s.lastLoggedBlock + 1; bn < currentBlockNumber; bn++ {
					if stats, exists := s.blockStats[bn]; exists && stats.ContractCount > 0 {
						avgGasPerByte := float64(stats.TotalGasUsed) / float64(max(stats.TotalBytecodeSize, 1))

						s.contractsMutex.Lock()
						totalContracts := len(s.deployedContracts)
						s.contractsMutex.Unlock()

						s.logger.WithFields(logrus.Fields{
							"block_number":        bn,
							"contracts_deployed":  stats.ContractCount,
							"total_gas_used":      stats.TotalGasUsed,
							"total_bytecode_size": stats.TotalBytecodeSize,
							"avg_gas_per_byte":    fmt.Sprintf("%.2f", avgGasPerByte),
							"total_contracts":     totalContracts,
						}).Info("Block deployment summary")

						s.lastLoggedBlock = bn
					}
				}
				s.blockStatsMutex.Unlock()
			}
		}
	}()
}

// stopBlockMonitor stops the block monitoring goroutine
func (s *Scenario) stopBlockMonitor() {
	if s.blockMonitorCancel != nil {
		s.blockMonitorCancel()
	}
	if s.blockMonitorDone != nil {
		<-s.blockMonitorDone // Wait for goroutine to finish
	}
}

func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Infof("starting scenario: %s", ScenarioName)
	defer s.logger.Infof("scenario %s finished.", ScenarioName)

	// Start block monitoring for real-time logging
	s.startBlockMonitor(ctx)
	defer s.stopBlockMonitor()

	// Cache chain ID at startup
	chainID, err := s.getChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}

	s.logger.Infof("Chain ID: %s", chainID.String())

	// Calculate adaptive max pending based on block gas limit if max-transactions is 0
	if s.options.MaxTransactions == 0 {
		// Get block gas limit from the network
		client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, "")
		if client == nil {
			return fmt.Errorf("no client available for gas limit query")
		}

		latestBlock, err := client.GetEthClient().BlockByNumber(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to get latest block: %w", err)
		}

		blockGasLimit := latestBlock.GasLimit()
		maxTxsPerBlock := blockGasLimit / EstimatedGasPerContract

		// Auto-calculate max pending transactions based on block gas limit
		// Allow up to 2x the block capacity to handle network delays
		s.maxPending = maxTxsPerBlock * MaxPendingMultiplier

		s.logger.Infof("Adaptive max pending enabled: block gas limit %d, gas per contract %d, max txs per block %d, max pending %d",
			blockGasLimit, EstimatedGasPerContract, maxTxsPerBlock, s.maxPending)
	}

	txIdxCounter := uint64(0)
	totalTxCount := atomic.Uint64{}

	for {
		// Check if we've reached max transactions (if set)
		if s.options.MaxTransactions > 0 && txIdxCounter >= s.options.MaxTransactions {
			s.logger.Infof("reached maximum number of transactions (%d)", s.options.MaxTransactions)
			break
		}

		// Wait for available slot based on max pending
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

		// Process pending transactions periodically with 1 second intervals
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

	// Wait for all pending transactions to complete with 1 second intervals
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
		time.Sleep(1 * time.Second) // Changed from 2 seconds to 1 second
	}

	// Stop block monitoring before final cleanup
	s.stopBlockMonitor()

	// Log any remaining unlogged blocks (final blocks) - keep this as final safety net
	s.blockStatsMutex.Lock()
	for bn, stats := range s.blockStats {
		if bn > s.lastLoggedBlock && stats.ContractCount > 0 {
			avgGasPerByte := float64(stats.TotalGasUsed) / float64(max(stats.TotalBytecodeSize, 1))

			s.logger.WithFields(logrus.Fields{
				"block_number":        bn,
				"contracts_deployed":  stats.ContractCount,
				"total_gas_used":      stats.TotalGasUsed,
				"total_bytecode_size": stats.TotalBytecodeSize,
				"avg_gas_per_byte":    fmt.Sprintf("%.2f", avgGasPerByte),
				"total_contracts":     len(s.deployedContracts),
			}).Info("Block deployment summary")
		}
	}
	s.blockStatsMutex.Unlock()

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
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, "")
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

		newBaseFeeGwei := new(big.Int).Add(currentBaseFeeGwei, big.NewInt(100))

		s.options.BaseFee = newBaseFeeGwei.Uint64()

		// Also increase tip fee slightly to ensure competitive priority
		if s.options.TipFee+1 > 3 {
			s.options.TipFee = s.options.TipFee + 1
		} else {
			s.options.TipFee = 2 // Minimum 3 gwei tip
		}

		s.logger.Infof("Updated dynamic fees - Base fee: %d gwei, Tip fee: %d gwei (network base fee: %s gwei)",
			s.options.BaseFee, s.options.TipFee, currentBaseFeeGwei.String())
	}

	return nil
}

// attemptTransaction makes a single attempt to send a transaction
func (s *Scenario) attemptTransaction(ctx context.Context, txIdx uint64, attempt int) error {
	// Get client and wallet
	client := s.walletPool.GetClient(spamoor.SelectClientRoundRobin, 0, "")
	wallet := s.walletPool.GetRootWallet()

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
	auth.GasLimit = DeploymentGasLimit // Fixed gas limit for contract deployment

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
		Timestamp:  time.Now(),
	}

	s.pendingTxsMutex.Lock()
	s.pendingTxs[tx.Hash()] = pendingTx
	s.pendingTxsMutex.Unlock()

	return nil
}
