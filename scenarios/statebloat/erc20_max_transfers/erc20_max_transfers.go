package sberc20maxtransfers

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/ethereum/go-ethereum/accounts/abi"
	//"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenario"
	contract "github.com/ethpandaops/spamoor/scenarios/statebloat/contract_deploy/contract"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
)

// Constants for ERC20 transfer operations
const (
	// ERC20TransferGasCost - gas cost for a standard ERC20 transfer (updated to 70K)
	ERC20TransferGasCost = 70000
	// DefaultBaseFeeGwei - default base fee in gwei
	DefaultBaseFeeGwei = 10
	// DefaultTipFeeGwei - default tip fee in gwei
	DefaultTipFeeGwei = 5
	// TokenTransferAmount - amount of tokens to transfer (1 token in smallest unit)
	TokenTransferAmount = 1
	// DefaultTargetGasRatio - target percentage of block gas limit to use (99.5% for minimal safety margin)
	DefaultTargetGasRatio = 0.995
	// FallbackBlockGasLimit - fallback gas limit if network query fails
	FallbackBlockGasLimit = 30000000
	// GweiPerEth - conversion factor from Gwei to Wei
	GweiPerEth = 1000000000
	// BlockMiningTimeout - timeout for waiting for a new block to be mined
	BlockMiningTimeout = 30 * time.Second
	// BlockPollingInterval - interval for checking new blocks
	BlockPollingInterval = 1 * time.Second
	// TransactionBatchSize - number of transactions to send in a batch
	TransactionBatchSize = 100
	// TransactionBatchThreshold - threshold for continuing to fill a block
	TransactionBatchThreshold = 50
	// InitialTransactionDelay - delay between initial transactions
	InitialTransactionDelay = 10 * time.Millisecond
	// OptimizedTransactionDelay - reduced delay after initial batch
	OptimizedTransactionDelay = 5 * time.Millisecond
	// ConfirmationDelay - delay before checking confirmations
	ConfirmationDelay = 2 * time.Second
	// MaxRebroadcasts - maximum number of times to rebroadcast a transaction
	MaxRebroadcasts = 10
	// RetryDelay - delay before retrying failed operations
	RetryDelay = 5 * time.Second
	// GasPerMillion - divisor for converting gas to millions
	GasPerMillion = 1_000_000.0
	// BytesPerKiB - bytes in a kibibyte
	BytesPerKiB = 1024.0
	// EstimatedStateGrowthPerTransfer - estimated state growth in bytes per new recipient
	EstimatedStateGrowthPerTransfer = 100
	// BloatingSummaryFileName - name of the bloating summary file
	BloatingSummaryFileName = "erc20_bloating_summary.json"
)

// ScenarioOptions defines the configuration options for the scenario
type ScenarioOptions struct {
	BaseFee         uint64 `yaml:"base_fee"`
	TipFee          uint64 `yaml:"tip_fee"`
	Contract        string `yaml:"contract"`
	DeploymentsFile string `yaml:"deployments_file"`
}

// DeploymentEntry represents a contract deployment from deployments.json
type DeploymentEntry map[string][]string

// ContractBloatStats tracks unique recipients per contract
type ContractBloatStats struct {
	UniqueRecipients int `json:"unique_recipients"`
}

// BloatingSummary represents the JSON file structure
type BloatingSummary struct {
	Contracts       map[string]*ContractBloatStats `json:"contracts"`
	TotalRecipients int                            `json:"total_recipients"`
	LastBlockNumber string                         `json:"last_block_number"`
	LastBlockUpdate time.Time                      `json:"last_block_update"`
}

// Scenario implements the ERC20 max transfers scenario
type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	// Deployed contracts and private key
	deployerPrivateKey   string
	deployerAddress      common.Address
	deployerWallet       *spamoor.Wallet // Store the wallet instance
	deployedContracts    []common.Address
	currentRoundContract common.Address // Contract being used for current round
	contractsLock        sync.Mutex

	// Transfer function ABI
	transferABI abi.Method
	contractABI abi.ABI

	// Used addresses tracking
	usedAddresses     map[common.Address]bool
	usedAddressesLock sync.Mutex

	// Bloating statistics tracking
	contractStats     map[common.Address]*ContractBloatStats
	contractStatsLock sync.Mutex
}

var ScenarioName = "statebloat-erc20-max-transfers"
var ScenarioDefaultOptions = ScenarioOptions{
	BaseFee:         DefaultBaseFeeGwei,
	TipFee:          DefaultTipFeeGwei,
	Contract:        "",
	DeploymentsFile: "deployments.json",
}
var ScenarioDescriptor = scenario.Descriptor{
	Name:           ScenarioName,
	Description:    "Maximum ERC20 transfers per block to unique addresses",
	DefaultOptions: ScenarioDefaultOptions,
	NewScenario:    newScenario,
}

func newScenario(logger logrus.FieldLogger) scenario.Scenario {
	return &Scenario{
		logger:        logger.WithField("scenario", ScenarioName),
		usedAddresses: make(map[common.Address]bool),
		contractStats: make(map[common.Address]*ContractBloatStats),
	}
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
	flags.Uint64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use in transactions (in gwei)")
	flags.Uint64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use in transactions (in gwei)")
	flags.StringVar(&s.options.Contract, "contract", ScenarioDefaultOptions.Contract, "Specific contract address to use (default: rotate through all)")
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

	// Load deployed contracts from deployments.json
	err := s.loadDeployedContracts()
	if err != nil {
		return fmt.Errorf("failed to load deployed contracts: %w", err)
	}

	// Load transfer function ABI
	err = s.loadTransferABI()
	if err != nil {
		return fmt.Errorf("failed to load transfer ABI: %w", err)
	}

	// We'll use the deployer wallet which we'll prepare in runMaxTransfersMode
	s.walletPool.SetWalletCount(0)

	return nil
}

func (s *Scenario) Config() string {
	yamlBytes, _ := yaml.Marshal(&s.options)
	return string(yamlBytes)
}

// loadDeployedContracts loads contract addresses and private key from deployments.json
func (s *Scenario) loadDeployedContracts() error {
	if s.options.DeploymentsFile == "" {
		return fmt.Errorf("deployments file is not set")
	}

	data, err := os.ReadFile(s.options.DeploymentsFile)
	if err != nil {
		return fmt.Errorf("failed to read %v: %w", s.options.DeploymentsFile, err)
	}

	var deployments DeploymentEntry
	err = json.Unmarshal(data, &deployments)
	if err != nil {
		return fmt.Errorf("failed to parse deployments.json: %w", err)
	}

	// Get the first (and only) entry
	for privateKey, addresses := range deployments {
		// Trim 0x prefix if present
		privateKey = strings.TrimPrefix(privateKey, "0x")

		s.deployerPrivateKey = privateKey
		s.deployedContracts = make([]common.Address, len(addresses))
		for i, addr := range addresses {
			s.deployedContracts[i] = common.HexToAddress(addr)
		}
		break // Only process the first entry
	}

	if s.deployerPrivateKey == "" || len(s.deployedContracts) == 0 {
		return fmt.Errorf("no valid deployments found in deployments.json")
	}

	s.logger.Infof("Loaded %d deployed contracts from deployments.json", len(s.deployedContracts))

	// Initialize contract stats for all deployed contracts
	for _, contractAddr := range s.deployedContracts {
		s.contractStats[contractAddr] = &ContractBloatStats{
			UniqueRecipients: 0,
		}
	}

	// If specific contract requested, validate it exists
	if s.options.Contract != "" {
		contractAddr := common.HexToAddress(s.options.Contract)
		found := false
		for _, addr := range s.deployedContracts {
			if addr == contractAddr {
				found = true
				s.deployedContracts = []common.Address{contractAddr} // Use only this contract
				break
			}
		}
		if !found {
			return fmt.Errorf("specified contract %s not found in deployments", s.options.Contract)
		}
		s.logger.Infof("Using specific contract: %s", contractAddr.Hex())
	}

	return nil
}

// loadTransferABI loads the transfer function ABI from the contract
func (s *Scenario) loadTransferABI() error {
	// Parse the contract ABI to get the transfer method
	contractABI, err := abi.JSON(strings.NewReader(contract.StateBloatTokenMetaData.ABI))
	if err != nil {
		return fmt.Errorf("failed to parse contract ABI: %w", err)
	}

	transferMethod, exists := contractABI.Methods["transfer"]
	if !exists {
		return fmt.Errorf("transfer method not found in contract ABI")
	}

	s.transferABI = transferMethod
	s.contractABI = contractABI
	return nil
}

// getNetworkBlockGasLimit retrieves the current block gas limit from the network
// It waits for a new block to be mined (with timeout) to ensure fresh data
func (s *Scenario) getNetworkBlockGasLimit(ctx context.Context, client *spamoor.Client) uint64 {
	// Create a timeout context for the entire operation
	timeoutCtx, cancel := context.WithTimeout(ctx, BlockMiningTimeout)
	defer cancel()

	// Get the current block number first
	currentBlockNumber, err := client.GetEthClient().BlockNumber(timeoutCtx)
	if err != nil {
		s.logger.Warnf("failed to get current block number: %v, using fallback: %d", err, FallbackBlockGasLimit)
		return FallbackBlockGasLimit
	}

	s.logger.Debugf("waiting for new block to be mined (current: %d, timeout: %v)", currentBlockNumber, BlockMiningTimeout)

	// Wait for a new block to be mined
	ticker := time.NewTicker(BlockPollingInterval)
	defer ticker.Stop()

	var latestBlock *types.Block
	for {
		select {
		case <-timeoutCtx.Done():
			s.logger.Warnf("timeout waiting for new block to be mined, using fallback: %d", FallbackBlockGasLimit)
			return FallbackBlockGasLimit
		case <-ticker.C:
			// Check for a new block
			newBlockNumber, err := client.GetEthClient().BlockNumber(timeoutCtx)
			if err != nil {
				s.logger.Debugf("error checking block number: %v", err)
				continue
			}

			// If we have a new block, get its details
			if newBlockNumber > currentBlockNumber {
				latestBlock, err = client.GetEthClient().BlockByNumber(timeoutCtx, nil)
				if err != nil {
					s.logger.Debugf("error getting latest block details: %v", err)
					continue
				}
				s.logger.Debugf("new block mined: %d", newBlockNumber)
				goto blockFound
			}
		}
	}

blockFound:
	gasLimit := latestBlock.GasLimit()
	s.logger.Debugf("network block gas limit from fresh block #%d: %d", latestBlock.NumberU64(), gasLimit)
	return gasLimit
}

// generateRecipient generates a deterministic recipient address based on index
func (s *Scenario) generateRecipient(recipientIndex uint64) common.Address {
	idxBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(idxBytes, recipientIndex)
	// Use deployer address as seed for deterministic generation
	hash := sha256.Sum256(append(s.deployerAddress.Bytes(), idxBytes...))
	return common.BytesToAddress(hash[12:]) // Use last 20 bytes as address
}

// loadBloatingSummary loads the bloating summary from file or creates a new one
func (s *Scenario) loadBloatingSummary() (*BloatingSummary, error) {
	data, err := os.ReadFile(BloatingSummaryFileName)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, return new summary
			return &BloatingSummary{
				Contracts:       make(map[string]*ContractBloatStats),
				TotalRecipients: 0,
			}, nil
		}
		return nil, fmt.Errorf("failed to read bloating summary: %w", err)
	}

	var summary BloatingSummary
	if err := json.Unmarshal(data, &summary); err != nil {
		return nil, fmt.Errorf("failed to unmarshal bloating summary: %w", err)
	}

	// Ensure contracts map is initialized
	if summary.Contracts == nil {
		summary.Contracts = make(map[string]*ContractBloatStats)
	}

	return &summary, nil
}

// saveBloatingSummary saves the bloating summary to file
func (s *Scenario) saveBloatingSummary(summary *BloatingSummary) error {
	data, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal bloating summary: %w", err)
	}

	if err := os.WriteFile(BloatingSummaryFileName, data, 0644); err != nil {
		return fmt.Errorf("failed to write bloating summary: %w", err)
	}

	return nil
}

// updateContractStats updates the statistics for a contract when a transfer is confirmed
func (s *Scenario) updateContractStats(contractAddr common.Address) {
	s.contractStatsLock.Lock()
	defer s.contractStatsLock.Unlock()

	stats, exists := s.contractStats[contractAddr]
	if !exists {
		stats = &ContractBloatStats{
			UniqueRecipients: 0,
		}
		s.contractStats[contractAddr] = stats
	}
	stats.UniqueRecipients++
}

// updateAndSaveBloatingSummary updates the bloating summary with current stats and saves to file
func (s *Scenario) updateAndSaveBloatingSummary(blockNumber string) error {
	// Load existing summary
	summary, err := s.loadBloatingSummary()
	if err != nil {
		return err
	}

	// Update with current stats
	s.contractStatsLock.Lock()
	totalRecipients := 0
	for contractAddr, stats := range s.contractStats {
		contractHex := contractAddr.Hex()
		summary.Contracts[contractHex] = &ContractBloatStats{
			UniqueRecipients: stats.UniqueRecipients,
		}
		totalRecipients += stats.UniqueRecipients
	}
	s.contractStatsLock.Unlock()

	// Update summary metadata
	summary.TotalRecipients = totalRecipients
	summary.LastBlockNumber = blockNumber
	summary.LastBlockUpdate = time.Now()

	// Save to file
	return s.saveBloatingSummary(summary)
}

// getContractBloatingSummaryForBlock returns a formatted string with contract bloating info for latest block
func (s *Scenario) getContractBloatingSummaryForBlock() string {
	s.contractStatsLock.Lock()
	defer s.contractStatsLock.Unlock()

	// Get current round contract
	s.contractsLock.Lock()
	currentContract := s.currentRoundContract
	s.contractsLock.Unlock()

	if currentContract == (common.Address{}) {
		return "No contract selected for current round"
	}

	// Get stats for current contract
	stats, exists := s.contractStats[currentContract]
	if !exists {
		return fmt.Sprintf("CONTRACT BLOATING STATUS:\n  Round Contract: %s - No transfers yet", currentContract.Hex())
	}

	return fmt.Sprintf("CONTRACT BLOATING STATUS:\n  Round Contract: %s - %d unique recipients",
		currentContract.Hex(), stats.UniqueRecipients)
}

// selectRandomContract randomly selects a contract from the deployed contracts
func (s *Scenario) selectRandomContract() (common.Address, error) {
	s.contractsLock.Lock()
	defer s.contractsLock.Unlock()

	if len(s.deployedContracts) == 0 {
		return common.Address{}, fmt.Errorf("no deployed contracts available")
	}

	// If only one contract, return it
	if len(s.deployedContracts) == 1 {
		return s.deployedContracts[0], nil
	}

	// Generate random index
	max := big.NewInt(int64(len(s.deployedContracts)))
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to generate random number: %w", err)
	}

	return s.deployedContracts[n.Int64()], nil
}

func (s *Scenario) Run(ctx context.Context) error {
	return s.runMaxTransfersMode(ctx)
}

func (s *Scenario) runMaxTransfersMode(ctx context.Context) error {
	s.logger.Infof("starting max transfers mode: self-adjusting to target block gas limit, continuous operation")

	// Get a client for network operations
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, "")

	// Get the actual network block gas limit
	networkGasLimit := s.getNetworkBlockGasLimit(ctx, client)
	targetGas := uint64(float64(networkGasLimit) * DefaultTargetGasRatio)

	// Calculate initial transfer count based on network gas limit and known gas cost per transfer
	initialTransfers := int(targetGas / ERC20TransferGasCost)

	// Dynamic transfer count - starts based on network parameters and adjusts based on actual performance
	currentTransfers := initialTransfers

	// Prepare the deployer wallet if not already done
	if s.deployerWallet == nil {
		// Create wallet from deployer private key
		deployerWallet, err := spamoor.NewWallet(s.deployerPrivateKey)
		if err != nil {
			return fmt.Errorf("failed to create deployer wallet: %w", err)
		}

		// Update wallet with chain info using the client
		err = client.UpdateWallet(ctx, deployerWallet)
		if err != nil {
			return fmt.Errorf("failed to update deployer wallet: %w", err)
		}

		// Store the wallet instance
		s.deployerWallet = deployerWallet
		s.deployerAddress = deployerWallet.GetAddress()

		s.logger.Infof("Initialized deployer wallet - Address: %s, Nonce: %d, Balance: %s ETH",
			s.deployerAddress.Hex(), deployerWallet.GetNonce(), new(big.Int).Div(deployerWallet.GetBalance(), big.NewInt(1e18)).String())
	}

	var blockCounter int
	var totalTransfers uint64
	var totalUniqueRecipients uint64

	for {
		select {
		case <-ctx.Done():
			s.logger.Errorf("max transfers mode stopping due to context cancellation")
			return ctx.Err()
		default:
		}

		blockCounter++

		// Send the max transfer transactions and wait for confirmation
		s.logger.Infof("════════════════ TRANSFER PHASE #%d ════════════════", blockCounter)
		actualGasUsed, blockNumber, transferCount, gasPerTransfer, uniqueRecipients, err := s.sendMaxTransfers(ctx, s.deployerWallet, currentTransfers, targetGas, blockCounter, client)
		if err != nil {
			s.logger.Errorf("failed to send max transfers for iteration %d: %v", blockCounter, err)
			time.Sleep(RetryDelay) // Wait before retry
			continue
		}

		// Update totals
		totalTransfers += uint64(transferCount)
		totalUniqueRecipients += uint64(uniqueRecipients)

		s.logger.Infof("%%%%%%%%%%%%%%%%%%%% ANALYSIS PHASE #%d %%%%%%%%%%%%%%%%%%%%", blockCounter)

		// Calculate metrics
		blockGasLimit := float64(networkGasLimit)
		gasUtilization := (float64(actualGasUsed) / blockGasLimit) * 100
		estimatedStateGrowth := uniqueRecipients * EstimatedStateGrowthPerTransfer

		s.logger.WithField("scenario", ScenarioName).Infof("TRANSFER METRICS - Block #%s | Transfers: %d | Unique recipients: %d | Gas used: %.2fM | Block utilization: %.2f%% | Gas/transfer: %.1f | Est. state growth: %.2f KiB",
			blockNumber, transferCount, uniqueRecipients, float64(actualGasUsed)/GasPerMillion, gasUtilization, gasPerTransfer, float64(estimatedStateGrowth)/BytesPerKiB)

		// Log contract-specific bloating info
		s.logger.WithField("scenario", ScenarioName).Info(s.getContractBloatingSummaryForBlock())

		// Log cumulative metrics
		s.logger.WithField("scenario", ScenarioName).Infof("CUMULATIVE TOTALS - Total transfers: %d | Total unique recipients: %d | Avg transfers/block: %.1f",
			totalTransfers, totalUniqueRecipients, float64(totalTransfers)/float64(blockCounter))

		// Update and save bloating summary
		err = s.updateAndSaveBloatingSummary(blockNumber)
		if err != nil {
			s.logger.Warnf("Failed to update bloating summary: %v", err)
		}

		// Self-adjust transfer count based on actual performance
		if actualGasUsed > 0 && transferCount > 0 {
			avgGasPerTransfer := float64(actualGasUsed) / float64(transferCount)
			targetTransfers := int(float64(targetGas) / avgGasPerTransfer)

			// Calculate the adjustment needed
			transferDifference := targetTransfers - transferCount

			if actualGasUsed < targetGas {
				// We're under target, increase transfer count with a slight safety margin
				newTransfers := currentTransfers + transferDifference - 1

				if newTransfers > currentTransfers {
					s.logger.Infof("Adjusting transfers: %d → %d (need %d more for target)",
						currentTransfers, newTransfers, transferDifference)
					currentTransfers = newTransfers
				}
			} else if actualGasUsed > targetGas {
				// We're over target, reduce to reach max block utilization
				excess := actualGasUsed - targetGas
				excessTransfers := int(float64(excess) / avgGasPerTransfer)
				newTransfers := currentTransfers - excessTransfers

				s.logger.Infof("Reducing transfers: %d → %d (excess: %d gas, ~%d transfers)",
					currentTransfers, newTransfers, excess, excessTransfers)
				currentTransfers = newTransfers

			} else {
				s.logger.Infof("Target achieved! Gas Used: %d / Target: %d", actualGasUsed, targetGas)
			}
		}
	}
}

func (s *Scenario) sendMaxTransfers(ctx context.Context, deployerWallet *spamoor.Wallet, targetTransfers int, targetGasLimit uint64, blockCounter int, client *spamoor.Client) (uint64, string, int, float64, int, error) {
	// Select a random contract for this round
	contractForRound, err := s.selectRandomContract()
	if err != nil {
		return 0, "", 0, 0, 0, fmt.Errorf("failed to select contract for round: %w", err)
	}

	// Update current round contract
	s.contractsLock.Lock()
	s.currentRoundContract = contractForRound
	s.contractsLock.Unlock()

	s.logger.Infof("Selected contract for round #%d: %s", blockCounter, contractForRound.Hex())

	// Get suggested fees or use configured values
	feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(client, s.options.BaseFee, s.options.TipFee)
	if err != nil {
		return 0, "", 0, 0, 0, fmt.Errorf("failed to get suggested fees: %w", err)
	}

	// Send transfers in batches
	return s.sendTransferBatch(ctx, deployerWallet, targetTransfers, targetGasLimit, blockCounter, client, feeCap, tipCap)
}

func (s *Scenario) sendTransferBatch(ctx context.Context, wallet *spamoor.Wallet, targetTransfers int, targetGasLimit uint64, iteration int, client *spamoor.Client, feeCap, tipCap *big.Int) (uint64, string, int, float64, int, error) {
	var confirmedCount int64
	var uniqueRecipientsCount int64
	var totalGasUsed uint64
	var lastBlockNumber string

	sentCount := 0
	recipientIndex := uint64(iteration * 1000000) // Large offset per iteration to avoid conflicts

	// Calculate approximate transactions per block based on gas limit
	maxTxsPerBlock := int(targetGasLimit / ERC20TransferGasCost)

	// Track confirmations
	type confirmResult struct {
		gasUsed      uint64
		blockNumber  string
		recipient    common.Address
		contractUsed common.Address
	}
	// Make channel buffered with enough capacity for all transactions
	confirmChan := make(chan confirmResult, targetTransfers*2) // Double buffer to be safe

	// Send transactions
	for sentCount < targetTransfers {
		// Generate unique recipient address
		var recipient common.Address
		for {
			recipient = s.generateRecipient(recipientIndex)
			recipientIndex++

			// Check if address already used
			s.usedAddressesLock.Lock()
			if !s.usedAddresses[recipient] {
				s.usedAddresses[recipient] = true
				s.usedAddressesLock.Unlock()
				break
			}
			s.usedAddressesLock.Unlock()
		}

		// Use the contract selected for this round
		s.contractsLock.Lock()
		contractAddr := s.currentRoundContract
		s.contractsLock.Unlock()

		// Encode transfer call data
		transferAmount := big.NewInt(TokenTransferAmount)
		callData, err := s.contractABI.Pack("transfer", recipient, transferAmount)
		if err != nil {
			s.logger.Errorf("failed to pack transfer call data: %v", err)
			continue
		}

		// Build transaction
		txMetadata := &txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Gas:       ERC20TransferGasCost,
			To:        &contractAddr,
			Value:     uint256.NewInt(0), // No ETH value for ERC20 transfer
			Data:      callData,
		}

		txData, err := txbuilder.DynFeeTx(txMetadata)
		if err != nil {
			s.logger.Errorf("failed to create tx data: %v", err)
			continue
		}

		tx, err := wallet.BuildDynamicFeeTx(txData)
		if err != nil {
			s.logger.Errorf("failed to build transaction: %v", err)
			continue
		}

		// Capture values for closure
		capturedRecipient := recipient
		capturedContract := contractAddr

		// Send transaction
		err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
			Client:      client,
			Rebroadcast: false, // No retries to avoid duplicates
			OnComplete: func(tx *types.Transaction, receipt *types.Receipt, err error) {
				if err != nil {
					return // Don't log individual failures
				}
				if receipt != nil && receipt.Status == 1 {
					atomic.AddInt64(&confirmedCount, 1)
					atomic.AddInt64(&uniqueRecipientsCount, 1)

					// Update contract stats
					s.updateContractStats(capturedContract)

					// Send result to channel with captured values
					confirmChan <- confirmResult{
						gasUsed:      receipt.GasUsed,
						blockNumber:  receipt.BlockNumber.String(),
						recipient:    capturedRecipient,
						contractUsed: capturedContract,
					}
				}
			},
			LogFn: func(client *spamoor.Client, retry int, rebroadcast int, err error) {
				// Only log actual send failures
				if err != nil {
					s.logger.Debugf("transfer tx send failed: %v", err)
				}
			},
		})

		if err != nil {
			continue
		}

		sentCount++

		// Small delay between transactions to ensure proper nonce ordering
		if sentCount < TransactionBatchSize {
			time.Sleep(InitialTransactionDelay)
		} else if sentCount%maxTxsPerBlock < TransactionBatchThreshold {
			time.Sleep(OptimizedTransactionDelay)
		}

		// Add context cancellation check
		select {
		case <-ctx.Done():
			return 0, "", 0, 0, 0, ctx.Err()
		default:
		}
	}

	// Wait for confirmations
	s.logger.Infof("Sent %d transfer transactions, waiting for confirmations...", sentCount)
	time.Sleep(ConfirmationDelay)

	// Log initial confirmation status
	initialConfirmed := atomic.LoadInt64(&confirmedCount)
	if initialConfirmed > 0 {
		s.logger.Debugf("Already have %d confirmations before collection", initialConfirmed)
	}

	// Collect results - wait for all sent transactions or timeout
	confirmTimeout := time.After(30 * time.Second)
	resultCount := 0

collectResults:
	for resultCount < sentCount {
		select {
		case result := <-confirmChan:
			totalGasUsed += result.gasUsed
			lastBlockNumber = result.blockNumber
			resultCount++

		case <-confirmTimeout:
			// Final check for any remaining confirmations
			confirmed := atomic.LoadInt64(&confirmedCount)
			s.logger.Warnf("Timeout waiting for confirmations, received %d results, %d confirmed, %d sent", resultCount, confirmed, sentCount)
			break collectResults

		case <-ctx.Done():
			return 0, "", 0, 0, 0, ctx.Err()
		}
	}

	// Drain any remaining results from the channel (non-blocking)
	for {
		select {
		case result := <-confirmChan:
			totalGasUsed += result.gasUsed
			lastBlockNumber = result.blockNumber
			resultCount++
		default:
			// No more results available
			goto done
		}
	}
done:

	// Calculate metrics
	confirmed := atomic.LoadInt64(&confirmedCount)
	uniqueRecipients := atomic.LoadInt64(&uniqueRecipientsCount)

	// Log detailed confirmation statistics
	s.logger.Debugf("Confirmation stats: sent=%d, confirmed=%d, results=%d, gas=%d",
		sentCount, confirmed, resultCount, totalGasUsed)

	if confirmed == 0 {
		return 0, "", 0, 0, 0, fmt.Errorf("no transfers confirmed")
	}

	gasPerTransfer := float64(totalGasUsed) / float64(confirmed)

	return totalGasUsed, lastBlockNumber, int(confirmed), gasPerTransfer, int(uniqueRecipients), nil
}
