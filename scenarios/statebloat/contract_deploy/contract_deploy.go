package sbcontractdeploy

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
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/scenarios/statebloat/contract_deploy/contract"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
)

type ScenarioOptions struct {
	MaxPending      uint64 `yaml:"max_pending"`
	MaxWallets      uint64 `yaml:"max_wallets"`
	BaseFee         uint64 `yaml:"base_fee"`
	TipFee          uint64 `yaml:"tip_fee"`
	ClientGroup     string `yaml:"client_group"`
	MaxTransactions uint64 `yaml:"max_transactions"`
	DeploymentsFile string `yaml:"deployments_file"`
}

// ContractDeployment tracks a deployed contract with its deployer info
type ContractDeployment struct {
	ContractAddress string `json:"contract_address"`
	PrivateKey      string `json:"private_key"`
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

	// Results tracking
	deployedContracts []ContractDeployment
	contractsMutex    sync.Mutex

	// Block-level statistics tracking
	blockStats      map[uint64]*BlockDeploymentStats
	blockStatsMutex sync.Mutex
	lastLoggedBlock uint64
}

var ScenarioName = "statebloat-contract-deploy"
var ScenarioDefaultOptions = ScenarioOptions{
	MaxWallets:      0, // Use root wallet only by default
	BaseFee:         5, // Moderate base fee (5 gwei)
	TipFee:          1, // Priority fee (1 gwei)
	MaxTransactions: 0,
	DeploymentsFile: "deployments.json",
}
var ScenarioDescriptor = scenario.Descriptor{
	Name:           ScenarioName,
	Description:    "Deploy contracts to create state bloat",
	DefaultOptions: ScenarioDefaultOptions,
	NewScenario:    newScenario,
}

func newScenario(logger logrus.FieldLogger) scenario.Scenario {
	return &Scenario{
		logger: logger.WithField("scenario", ScenarioName),
	}
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Base fee per gas in gwei")
	flags.Uint64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Tip fee per gas in gwei")
	flags.StringVar(&s.options.ClientGroup, "client-group", ScenarioDefaultOptions.ClientGroup, "Client group to use for sending transactions")
	flags.Uint64Var(&s.options.MaxTransactions, "max-transactions", ScenarioDefaultOptions.MaxTransactions, "Maximum number of transactions to send (0 = use rate limiting based on block gas limit)")
	flags.StringVar(&s.options.DeploymentsFile, "deployments-file", ScenarioDefaultOptions.DeploymentsFile, "File to save deployments to")
	return nil
}

func (s *Scenario) Init(options *scenario.Options) error {
	s.walletPool = options.WalletPool

	if options.Config != "" {
		err := yaml.Unmarshal([]byte(options.Config), &s.options)
		if err != nil {
			return fmt.Errorf("failed to unmarshal config: %w", err)
		}
	}

	if s.options.MaxWallets > 0 {
		s.walletPool.SetWalletCount(s.options.MaxWallets)
	} else {
		s.walletPool.SetWalletCount(10)
	}

	return nil
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
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, s.options.ClientGroup)
	// TODO: This should be a constant documented on how this number is obtained.
	var bytecodeSize int = 23914

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
	if s.options.DeploymentsFile == "" {
		return nil
	}

	// Create a map from private key to array of contract addresses
	deploymentMap := make(map[string][]string)

	for _, contract := range s.deployedContracts {
		privateKey := contract.PrivateKey
		contractAddr := contract.ContractAddress
		deploymentMap[privateKey] = append(deploymentMap[privateKey], contractAddr)
	}

	// Create or overwrite the deployments.json file
	deploymentsFile, err := os.Create(s.options.DeploymentsFile)
	if err != nil {
		return fmt.Errorf("failed to create %v file: %w", s.options.DeploymentsFile, err)
	}
	defer deploymentsFile.Close()

	// Write the mapping as JSON with pretty formatting
	encoder := json.NewEncoder(deploymentsFile)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(deploymentMap)
	if err != nil {
		return fmt.Errorf("failed to write %v: %w", s.options.DeploymentsFile, err)
	}

	return nil
}

// startBlockMonitor starts a background goroutine that monitors for new blocks
// and logs block deployment summaries immediately when blocks are mined
func (s *Scenario) startBlockMonitor(ctx context.Context) {
	go func() {
		client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, s.options.ClientGroup)
		if client == nil {
			s.logger.Warn("No client available for block monitoring")
			return
		}

		ethClient := client.GetEthClient()
		ticker := time.NewTicker(2 * time.Second) // Poll every 2 seconds
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// Get current block number
				latestBlock, err := ethClient.BlockByNumber(ctx, nil)
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

func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Infof("starting scenario: %s", ScenarioName)
	defer s.logger.Infof("scenario %s finished.", ScenarioName)

	// Start block monitoring for real-time logging
	s.startBlockMonitor(ctx)

	// Calculate rate limiting based on block gas limit if max-transactions is 0
	var maxTxsPerBlock uint64
	var maxPending uint64 = 100
	var totalTxCount uint64 = 0

	if s.options.MaxTransactions == 0 {
		// Get block gas limit from the network
		client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, s.options.ClientGroup)
		if client == nil {
			return fmt.Errorf("no client available for gas limit query")
		}

		latestBlock, err := client.GetEthClient().BlockByNumber(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to get latest block: %w", err)
		}

		blockGasLimit := latestBlock.GasLimit()
		// TODO: This should be a constant.
		estimatedGasPerContract := uint64(4949468) // Updated estimate based on contract size reduction
		maxTxsPerBlock = blockGasLimit / estimatedGasPerContract

		s.logger.Infof("Rate limiting enabled: block gas limit %d, gas per contract %d, max txs per block %d",
			blockGasLimit, estimatedGasPerContract, maxTxsPerBlock)
	}

	if s.options.MaxPending > 0 {
		maxPending = s.options.MaxPending
	}

	err := scenario.RunTransactionScenario(ctx, scenario.TransactionScenarioOptions{
		TotalCount: s.options.MaxTransactions,
		Throughput: maxTxsPerBlock,
		MaxPending: maxPending,
		WalletPool: s.walletPool,

		Logger: s.logger,
		ProcessNextTxFn: func(ctx context.Context, txIdx uint64, onComplete func()) (func(), error) {
			logger := s.logger
			tx, err := s.sendTransaction(ctx, txIdx, onComplete)

			atomic.AddUint64(&totalTxCount, 1)

			return func() {
				if err != nil {
					logger.Warnf("could not send transaction: %v", err)
				} else {
					logger.Debugf("sent deployment tx #%6d: %v", txIdx+1, tx.Hash().String())
				}
			}, err
		},
	})

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
		"total_txs":       totalTxCount,
		"total_contracts": totalContracts,
	}).Info("All transactions completed")

	return err
}

// sendTransaction sends a single contract deployment transaction
func (s *Scenario) sendTransaction(ctx context.Context, txIdx uint64, onComplete func()) (*types.Transaction, error) {
	maxRetries := 3

	for attempt := 0; attempt < maxRetries; attempt++ {
		tx, err := s.attemptTransaction(ctx, txIdx, attempt, onComplete)
		if err == nil {
			return tx, nil
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
		return nil, err
	}

	return nil, fmt.Errorf("failed to send transaction after %d attempts", maxRetries)
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
func (s *Scenario) attemptTransaction(ctx context.Context, txIdx uint64, attempt int, onComplete func()) (*types.Transaction, error) {
	// Get client and wallet
	client := s.walletPool.GetClient(spamoor.SelectClientRoundRobin, 0, "")
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(txIdx))

	if client == nil {
		return nil, fmt.Errorf("no client available")
	}
	if wallet == nil {
		return nil, fmt.Errorf("no wallet available")
	}

	// Set EIP-1559 fee parameters
	feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(client, s.options.BaseFee, s.options.TipFee)
	if err != nil {
		return nil, fmt.Errorf("failed to get suggested fees: %w", err)
	}

	// Generate random salt for unique contract
	salt := make([]byte, 32)
	_, err = rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	saltInt := new(big.Int).SetBytes(salt)

	// Deploy the contract
	ethClient := client.GetEthClient()
	if ethClient == nil {
		return nil, fmt.Errorf("failed to get eth client")
	}

	tx, err := wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       5200000,
		Value:     uint256.NewInt(0),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		_, deployTx, _, err := contract.DeployStateBloatToken(transactOpts, client.GetEthClient(), saltInt)
		return deployTx, err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create deployment transaction: %w", err)
	}

	mu := sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()

	var callOnComplete bool

	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		Rebroadcast: true,
		OnComplete: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			defer func() {
				mu.Lock()
				defer mu.Unlock()

				if callOnComplete {
					onComplete()
				}
			}()

			if receipt != nil {
				s.recordDeployedContract(receipt.ContractAddress, wallet.GetPrivateKey(), receipt, tx.Hash())
			}
		},
	})

	callOnComplete = err == nil
	if err != nil {
		return nil, err
	}

	return tx, nil
}
