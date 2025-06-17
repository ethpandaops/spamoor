package eoadelegation

import (
	"context"
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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/utils"
)

// EIP-7702 gas cost constants
const (
	// PER_AUTH_BASE_COST (EIP-7702) - upper bound on gas cost per authorization when delegating to existing contract in this scenario
	GasPerAuthorization = 26000
	// EstimatedBytesPerAuth - estimated state change in bytes per EOA delegation
	EstimatedBytesPerAuth = 135.0
	// ActualBytesPerAuth - actual observed RLP-encoded bytes per authorization in transaction
	// RLP encoding breakdown for typical authorization:
	// - ChainID: ~3 bytes (RLP encodes integers efficiently, not fixed 8 bytes)
	// - Address: 21 bytes (0x94 prefix + 20 bytes)
	// - Nonce: 1 byte (0x80 for nonce 0, small values)
	// - YParity: 1 byte (0x00 or 0x01)
	// - R: 33 bytes (0xa0 prefix + 32 bytes)
	// - S: 33 bytes (0xa0 prefix + 32 bytes)
	// - List overhead: 2 bytes (0xf8 + length byte)
	// Total: ~94 bytes (confirmed by empirical data: 89KiB for 969 auths)
	ActualBytesPerAuth = 94
	// DefaultTargetGasRatio - target percentage of block gas limit to use (99.5% for minimal safety margin)
	DefaultTargetGasRatio = 0.995
	// FallbackBlockGasLimit - fallback gas limit if network query fails
	FallbackBlockGasLimit = 30000000
	// BaseTransferCost - gas cost for a standard ETH transfer
	BaseTransferCost = 21000
	// MaxTransactionSize - Ethereum transaction size limit in bytes (128KiB)
	MaxTransactionSize = 131072 // 128 * 1024
	// GweiPerEth - conversion factor from Gwei to Wei
	GweiPerEth = 1000000000
	// BlockMiningTimeout - timeout for waiting for a new block to be mined
	BlockMiningTimeout = 30 * time.Second
	// BlockPollingInterval - interval for checking new blocks
	BlockPollingInterval = 1 * time.Second
	// MaxRebroadcasts - maximum number of times to rebroadcast a transaction
	MaxRebroadcasts = 10
	// TransactionBatchSize - used for batching funding transactions
	TransactionBatchSize = 100
	// TransactionBatchThreshold - threshold for continuing to fill a block
	TransactionBatchThreshold = 50
	// InitialTransactionDelay - delay between initial funding transactions
	InitialTransactionDelay = 10 * time.Millisecond
	// OptimizedTransactionDelay - reduced delay after initial batch
	OptimizedTransactionDelay = 5 * time.Millisecond
	// FundingConfirmationDelay - delay before checking funding confirmations
	FundingConfirmationDelay = 3 * time.Second
	// RetryDelay - delay before retrying failed operations
	RetryDelay = 5 * time.Second
	// FundingIterationOffset - large offset to avoid delegator index conflicts between iterations
	FundingIterationOffset = 1000000
	// TransactionBaseOverhead - base RLP encoding overhead for a transaction
	TransactionBaseOverhead = 200
	// TransactionExtraOverhead - additional RLP encoding overhead
	TransactionExtraOverhead = 50
	// GasPerCallDataByte - gas cost per byte of calldata (16 gas per non-zero byte)
	GasPerCallDataByte = 16
	// BytesPerKiB - bytes in a kibibyte
	BytesPerKiB = 1024.0
	// GasPerMillion - divisor for converting gas to millions
	GasPerMillion = 1_000_000.0
	// TransactionSizeSafetyFactor - safety factor for transaction size (95%)
	TransactionSizeSafetyFactor = 95
)

type ScenarioOptions struct {
	BaseFee  uint64 `yaml:"base_fee"`
	TipFee   uint64 `yaml:"tip_fee"`
	CodeAddr string `yaml:"code_addr"`
}

// EOAEntry represents a funded EOA account
type EOAEntry struct {
	Address    string `json:"address"`
	PrivateKey string `json:"private_key"`
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	// FIFO queue for funded accounts
	eoaQueue      []EOAEntry
	eoaQueueMutex sync.Mutex

	// Semaphore for worker control
	workerSemaphore chan struct{}
	workerDone      chan struct{}
	workerWg        sync.WaitGroup
}

var ScenarioName = "eoa-delegation"
var ScenarioDefaultOptions = ScenarioOptions{
	BaseFee:  20,
	TipFee:   2,
	CodeAddr: "",
}
var ScenarioDescriptor = scenario.Descriptor{
	Name:           ScenarioName,
	Description:    "Maximum state bloating via EIP-7702 EOA delegations",
	DefaultOptions: ScenarioDefaultOptions,
	NewScenario:    newScenario,
}

func newScenario(logger logrus.FieldLogger) scenario.Scenario {
	return &Scenario{
		logger: logger.WithField("scenario", ScenarioName),
	}
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
	flags.Uint64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use in transactions (in gwei)")
	flags.Uint64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use in transactions (in gwei)")
	flags.StringVar(&s.options.CodeAddr, "code-addr", ScenarioDefaultOptions.CodeAddr, "Code delegation target address to use for transactions (default: ecrecover precompile)")
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

	if s.options.CodeAddr == "" {
		s.logger.Infof("no --code-addr specified, using ecrecover precompile as delegate: %s", common.HexToAddress("0x0000000000000000000000000000000000000001"))
	}

	// In max-bloating mode, use 1 wallet (the root wallet)
	s.walletPool.SetWalletCount(1)

	// Initialize FIFO queue and worker for EOA management
	s.eoaQueue = make([]EOAEntry, 0)
	s.workerSemaphore = make(chan struct{}, 1) // Buffered channel for semaphore
	s.workerDone = make(chan struct{})

	// Start the worker goroutine for writing EOAs to file
	s.workerWg.Add(1)
	go s.eoaWorker()

	return nil
}

func (s *Scenario) Config() string {
	yamlBytes, _ := yaml.Marshal(&s.options)
	return string(yamlBytes)
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

func (s *Scenario) Run(ctx context.Context) error {
	// This scenario only runs in max-bloating mode
	return s.runMaxBloatingMode(ctx)
}

func (s *Scenario) prepareDelegator(delegatorIndex uint64) (*spamoor.Wallet, error) {
	idxBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(idxBytes, delegatorIndex)
	childKey := sha256.Sum256(append(common.FromHex(s.walletPool.GetRootWallet().GetWallet().GetAddress().Hex()), idxBytes...))
	return spamoor.NewWallet(fmt.Sprintf("%x", childKey))
}

func (s *Scenario) buildMaxBloatingAuthorizations(targetCount int, iteration int) []types.SetCodeAuthorization {
	authorizations := make([]types.SetCodeAuthorization, 0, targetCount)

	// Use a fixed delegate contract address for maximum efficiency
	// In max bloating mode, we want all EOAs to delegate to the same existing contract
	// to benefit from reduced gas costs (PER_AUTH_BASE_COST vs PER_EMPTY_ACCOUNT_COST)
	// Precompiles are ideal as they're guaranteed to exist with code on all networks
	var codeAddr common.Address
	if s.options.CodeAddr != "" {
		codeAddr = common.HexToAddress(s.options.CodeAddr)
	} else {
		// Default to using the ecrecover precompile (0x1) as delegate target
		codeAddr = common.HexToAddress("0x0000000000000000000000000000000000000001")
	}

	chainId := s.walletPool.GetChainId().Uint64()

	for i := 0; i < targetCount; i++ {
		// Create a unique delegator for each authorization
		// Include iteration counter to ensure different addresses for each iteration
		delegatorIndex := uint64(iteration*targetCount + i)

		delegator, err := s.prepareDelegator(delegatorIndex)
		if err != nil {
			s.logger.Errorf("could not prepare delegator %v: %v", delegatorIndex, err)
			continue
		}

		// Each EOA uses auth_nonce = 0 (assuming first EIP-7702 operation)
		// This creates maximum new state as each EOA gets its first delegation
		authorization := types.SetCodeAuthorization{
			ChainID: *uint256.NewInt(chainId),
			Address: codeAddr,
			Nonce:   0, // First delegation for each EOA
		}

		// Sign the authorization with the delegator's private key
		signedAuth, err := types.SignSetCode(delegator.GetPrivateKey(), authorization)
		if err != nil {
			s.logger.Errorf("could not sign set code authorization for delegator %v: %v", delegatorIndex, err)
			continue
		}

		authorizations = append(authorizations, signedAuth)
	}

	return authorizations
}

func (s *Scenario) runMaxBloatingMode(ctx context.Context) error {
	s.logger.Infof("starting max bloating mode: self-adjusting to target block gas limit, continuous operation")

	// Get a client for network operations
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, "")

	// Get the actual network block gas limit
	networkGasLimit := s.getNetworkBlockGasLimit(ctx, client)
	targetGas := uint64(float64(networkGasLimit) * DefaultTargetGasRatio)

	// Calculate initial authorization count based on network gas limit and known gas cost per authorization
	initialAuthorizations := int(targetGas / GasPerAuthorization)

	// Dynamic authorization count - starts based on network parameters and adjusts based on actual performance
	currentAuthorizations := initialAuthorizations

	var blockCounter int

	for {
		select {
		case <-ctx.Done():
			s.logger.Errorf("max bloating mode stopping due to context cancellation")
			return ctx.Err()
		default:
		}

		blockCounter++

		// For the first iteration, we need to fund delegators before bloating
		// For subsequent iterations, funding happens after analysis
		if blockCounter == 1 {
			s.logger.Infof("════════════════ INITIAL FUNDING PHASE ════════════════")
			confirmedCount, err := s.fundMaxBloatingDelegators(ctx, currentAuthorizations, blockCounter, networkGasLimit)
			if err != nil {
				s.logger.Errorf("failed to fund delegators for initial iteration: %v", err)
				time.Sleep(RetryDelay) // Wait before retry
				blockCounter--         // Retry the same iteration
				continue
			}

			// Wait for funding transactions to be confirmed and included in blocks
			err = s.waitForFundingConfirmations(ctx, confirmedCount)
			if err != nil {
				s.logger.Errorf("error waiting for funding confirmations: %v", err)
			}
		}

		// Send the max bloating transaction and wait for confirmation
		s.logger.Infof("════════════════ BLOATING PHASE #%d ════════════════", blockCounter)
		actualGasUsed, _, authCount, gasPerAuth, gasPerByte, _, err := s.sendMaxBloatingTransaction(ctx, currentAuthorizations, targetGas, blockCounter)
		if err != nil {
			s.logger.Errorf("failed to send max bloating transaction for iteration %d: %v", blockCounter, err)
			time.Sleep(RetryDelay) // Wait before retry
			continue
		}

		// Open semaphore (green light) during analysis phase to allow worker to process EOA queue
		s.openWorkerSemaphore()

		s.logger.Infof("%%%%%%%%%%%%%%%%%%%% ANALYSIS PHASE #%d %%%%%%%%%%%%%%%%%%%%", blockCounter)

		// Calculate total bytes written to state
		totalBytesWritten := authCount * int(EstimatedBytesPerAuth)

		// Get block gas limit for utilization calculation
		blockGasLimit := float64(networkGasLimit)
		gasUtilization := (float64(actualGasUsed) / blockGasLimit) * 100

		s.logger.WithField("scenario", "eoa-delegation").Infof("STATE BLOATING METRICS - Total bytes written: %.2f KiB, Gas used: %.2fM, Block utilization: %.2f%%, Authorizations: %d, Gas/auth: %.1f, Gas/byte: %.1f",
			float64(totalBytesWritten)/BytesPerKiB, float64(actualGasUsed)/GasPerMillion, gasUtilization, authCount, gasPerAuth, gasPerByte)

		// Self-adjust authorization count based on actual performance
		if actualGasUsed > 0 && authCount > 0 {
			gasPerAuth := float64(actualGasUsed) / float64(authCount)
			targetAuths := int(float64(targetGas) / gasPerAuth)

			// Calculate the adjustment needed
			authDifference := targetAuths - authCount

			if actualGasUsed < targetGas {
				// We're under target, increase authorization count with a slight safety margin
				newAuthorizations := currentAuthorizations + authDifference - 1

				if newAuthorizations > currentAuthorizations {
					s.logger.Infof("Adjusting authorizations: %d → %d (need %d more for target)",
						currentAuthorizations, newAuthorizations, authDifference)
					currentAuthorizations = newAuthorizations
				}
			} else if actualGasUsed > targetGas {
				// We're over target, reduce to reach max block utilization
				excess := actualGasUsed - targetGas
				newAuthorizations := currentAuthorizations - int(excess) + 1

				s.logger.Infof("Reducing authorizations: %d → %d (excess: %d gas)",
					currentAuthorizations, newAuthorizations, excess)
				currentAuthorizations = newAuthorizations

			} else {
				s.logger.Infof("Target achieved! Gas Used: %d / Target: %d", actualGasUsed, targetGas)
			}
		}

		// Now fund delegators for the next iteration (except on the last iteration)
		// This ensures funding happens AFTER bloating transactions are confirmed
		s.logger.Infof("════════════════ FUNDING PHASE #%d (for next iteration) ════════════════", blockCounter)
		confirmedCount, err := s.fundMaxBloatingDelegators(ctx, currentAuthorizations, blockCounter+1, networkGasLimit)
		if err != nil {
			s.logger.Errorf("failed to fund delegators for next iteration: %v", err)
			// Don't fail the entire loop, just log the error and continue
		}

		// Wait for funding transactions to be confirmed before next bloating phase
		if confirmedCount > 0 {
			err = s.waitForFundingConfirmations(ctx, confirmedCount)
			if err != nil {
				s.logger.Errorf("error waiting for funding confirmations: %v", err)
			}
		}
	}
}

func (s *Scenario) fundMaxBloatingDelegators(ctx context.Context, targetCount int, iteration int, gasLimit uint64) (int64, error) {
	// Close semaphore (red light) during funding phase
	s.closeWorkerSemaphore()

	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, "")
	if client == nil {
		return 0, fmt.Errorf("no client available for funding delegators")
	}

	// Use root wallet since we set child wallet count to 0 in max-bloating mode
	wallet := s.walletPool.GetRootWallet()

	// Get suggested fees for funding transactions
	var feeCap *big.Int
	var tipCap *big.Int

	if s.options.BaseFee > 0 {
		feeCap = new(big.Int).Mul(big.NewInt(int64(s.options.BaseFee)), big.NewInt(GweiPerEth))
	}
	if s.options.TipFee > 0 {
		tipCap = new(big.Int).Mul(big.NewInt(int64(s.options.TipFee)), big.NewInt(GweiPerEth))
	}

	if feeCap == nil || tipCap == nil {
		var err error
		feeCap, tipCap, err = client.GetSuggestedFee(s.walletPool.GetContext())
		if err != nil {
			return 0, fmt.Errorf("failed to get suggested fees for funding: %w", err)
		}
	}

	// Minimum gas prices
	if feeCap.Cmp(big.NewInt(GweiPerEth)) < 0 {
		feeCap = big.NewInt(GweiPerEth)
	}
	if tipCap.Cmp(big.NewInt(GweiPerEth)) < 0 {
		tipCap = big.NewInt(GweiPerEth)
	}

	// Fund with 1 wei as requested by user
	fundingAmount := uint256.NewInt(1)

	var confirmedCount int64
	sentCount := uint64(0)
	delegatorIndex := uint64(iteration * FundingIterationOffset) // Large offset per iteration to avoid conflicts

	// Calculate approximate transactions per block based on gas limit
	// Standard transfer = BaseTransferCost gas.
	var maxTxsPerBlock = gasLimit / uint64(BaseTransferCost)

	for {
		// Check if we have enough confirmed transactions
		confirmed := atomic.LoadInt64(&confirmedCount)
		if confirmed >= int64(targetCount) {
			// We have minimum required, but let's check if we should fill the current block
			// If we've sent transactions recently, wait a bit to see if block gets filled
			if sentCount > 0 && (sentCount%TransactionBatchSize) > TransactionBatchThreshold {
				// Continue to fill the block
			} else {
				break
			}
		}

		// Generate unique delegator address
		delegator, err := s.prepareDelegator(delegatorIndex)
		if err != nil {
			s.logger.Errorf("could not prepare delegator %v for funding: %v", delegatorIndex, err)
			delegatorIndex++
			continue
		}

		// Build funding transaction
		delegatorAddr := delegator.GetAddress()
		txData, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Gas:       BaseTransferCost, // Standard ETH transfer gas
			To:        &delegatorAddr,
			Value:     fundingAmount,
			Data:      []byte{},
		})
		if err != nil {
			s.logger.Errorf("failed to build funding tx for delegator %d: %v", delegatorIndex, err)
			delegatorIndex++
			continue
		}

		tx, err := wallet.GetWallet().BuildDynamicFeeTx(txData)
		if err != nil {
			s.logger.Errorf("failed to build funding transaction for delegator %d: %v", delegatorIndex, err)
			delegatorIndex++
			continue
		}

		// Send funding transaction with no retries to avoid duplicates
		err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet.GetWallet(), tx, &spamoor.SendTransactionOptions{
			Client:      client,
			Rebroadcast: false, // No retries to avoid duplicates
			OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
				if err != nil {
					return // Don't log individual failures
				}
				if receipt != nil && receipt.Status == 1 {
					atomic.AddInt64(&confirmedCount, 1)

					// Add successfully funded delegator to EOA queue
					s.addEOAToQueue(delegator.GetAddress().Hex(), fmt.Sprintf("%x", delegator.GetPrivateKey().D))

					// No progress logging - only log when target is reached
				}
			},
			LogFn: func(client *spamoor.Client, retry int, rebroadcast int, err error) {
				// Only log actual send failures, not confirmation failures
				if err != nil {
					s.logger.Debugf("funding tx send failed: %v", err)
				}
			},
		})

		if err != nil {
			delegatorIndex++
			continue
		}

		sentCount++
		delegatorIndex++

		// Check if we should continue filling the block
		confirmed = atomic.LoadInt64(&confirmedCount)
		if confirmed >= int64(targetCount) && sentCount%maxTxsPerBlock < TransactionBatchSize {
			// We have enough confirmed and we're at the end of a block cycle, stop for now
			break
		}

		// Small delay between transactions to ensure proper nonce ordering
		// Reduce delay as we get more efficient
		if sentCount < TransactionBatchSize {
			time.Sleep(InitialTransactionDelay)
		} else {
			time.Sleep(OptimizedTransactionDelay)
		}

		// Add context cancellation check
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
		}
	}

	// Wait for any remaining transactions to be included
	time.Sleep(FundingConfirmationDelay)

	// Return the confirmed count
	confirmed := atomic.LoadInt64(&confirmedCount)
	return confirmed, nil
}

// waitForFundingConfirmations waits for funding transactions to be confirmed by monitoring for new blocks
func (s *Scenario) waitForFundingConfirmations(ctx context.Context, targetConfirmations int64) error {
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, "")
	if client == nil {
		return fmt.Errorf("no client available for monitoring blocks")
	}

	s.logger.Infof("Waiting for funding transactions to be confirmed (expecting ~%d confirmations)...", targetConfirmations)

	// Get the starting block number
	startBlock, err := client.GetEthClient().BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to get starting block number: %w", err)
	}

	// Monitor until we see at least 1 new block to ensure funding txs are included
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	blocksWaited := uint64(0)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			// Check current block number
			currentBlock, err := client.GetEthClient().BlockNumber(ctx)
			if err != nil {
				s.logger.Debugf("Error getting block number: %v", err)
				continue
			}

			// If we have new blocks
			if currentBlock > startBlock {
				blocksWaited = currentBlock - startBlock
				s.logger.Debugf("New block %d mined (%d blocks since funding started)", currentBlock, blocksWaited)

				// Wait for at least 1 block to ensure funding transactions are included
				if blocksWaited >= 1 {
					s.logger.Infof("Funding transactions should be confirmed (waited %d blocks)", blocksWaited)
					return nil
				}
			}
		}
	}
}

func (s *Scenario) sendMaxBloatingTransaction(ctx context.Context, targetAuthorizations int, targetGasLimit uint64, blockCounter int) (uint64, string, int, float64, float64, string, error) {
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, "")
	if client == nil {
		return 0, "", 0, 0, 0, "", fmt.Errorf("no client available for sending max bloating transaction")
	}

	// Use root wallet since we set child wallet count to 0 in max-bloating mode
	wallet := s.walletPool.GetRootWallet()

	// Get suggested fees or use configured values
	var feeCap *big.Int
	var tipCap *big.Int

	if s.options.BaseFee > 0 {
		feeCap = new(big.Int).Mul(big.NewInt(int64(s.options.BaseFee)), big.NewInt(GweiPerEth))
	}
	if s.options.TipFee > 0 {
		tipCap = new(big.Int).Mul(big.NewInt(int64(s.options.TipFee)), big.NewInt(GweiPerEth))
	}

	if feeCap == nil || tipCap == nil {
		var err error
		feeCap, tipCap, err = client.GetSuggestedFee(s.walletPool.GetContext())
		if err != nil {
			return 0, "", 0, 0, 0, "", fmt.Errorf("failed to get suggested fees: %w", err)
		}
	}

	// Ensure minimum gas prices for inclusion
	if feeCap.Cmp(big.NewInt(GweiPerEth)) < 0 {
		feeCap = big.NewInt(GweiPerEth)
	}
	if tipCap.Cmp(big.NewInt(GweiPerEth)) < 0 {
		tipCap = big.NewInt(GweiPerEth)
	}

	// Use minimal amount for max bloating (focus on authorizations, not value transfer)
	amount := uint256.NewInt(0) // No value transfer needed

	// Target address - use our own wallet for simplicity
	toAddr := wallet.GetWallet().GetAddress()

	// No call data for max bloating transactions
	txCallData := []byte{}

	// Build the authorizations for maximum state bloat
	authorizations := s.buildMaxBloatingAuthorizations(targetAuthorizations, blockCounter)

	// Check transaction size and split into batches if needed
	batches := s.splitAuthorizationsBatches(authorizations, len(txCallData))

	if len(batches) == 1 {
		// Single transaction - use existing logic
		return s.sendSingleMaxBloatingTransaction(ctx, batches[0], txCallData, feeCap, tipCap, amount, toAddr, targetGasLimit, wallet.GetWallet(), client)
	} else {
		// Multiple transactions needed - send them as a batch
		return s.sendBatchedMaxBloatingTransactions(ctx, batches, txCallData, feeCap, tipCap, amount, toAddr, targetGasLimit, wallet.GetWallet(), client)
	}
}

// sendSingleMaxBloatingTransaction sends a single transaction (original logic)
func (s *Scenario) sendSingleMaxBloatingTransaction(ctx context.Context, authorizations []types.SetCodeAuthorization, txCallData []byte, feeCap, tipCap *big.Int, amount *uint256.Int, toAddr common.Address, targetGasLimit uint64, wallet *spamoor.Wallet, client *spamoor.Client) (uint64, string, int, float64, float64, string, error) {
	txData, err := txbuilder.SetCodeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       targetGasLimit,
		To:        &toAddr,
		Value:     amount,
		Data:      txCallData,
		AuthList:  authorizations,
	})
	if err != nil {
		return 0, "", 0, 0, 0, "", fmt.Errorf("failed to build transaction metadata: %w", err)
	}

	tx, err := wallet.BuildSetCodeTx(txData)
	if err != nil {
		return 0, "", 0, 0, 0, "", fmt.Errorf("failed to build transaction: %w", err)
	}

	// Log actual transaction size
	txSize := len(tx.Data())
	if encoded, err := tx.MarshalBinary(); err == nil {
		txSize = len(encoded)
	}
	sizeKiB := float64(txSize) / BytesPerKiB
	exceedsLimit := txSize > MaxTransactionSize
	limitKiB := float64(MaxTransactionSize) / BytesPerKiB

	s.logger.WithField("scenario", "eoa-delegation").Infof("MAX BLOATING TX SIZE: %d bytes (%.2f KiB) | Limit: %d bytes (%.1f KiB) | %d authorizations | Exceeds limit: %v",
		txSize, sizeKiB, MaxTransactionSize, limitKiB, len(authorizations), exceedsLimit)

	// Use channels to capture transaction results
	resultChan := make(chan struct {
		gasUsed      uint64
		blockNumber  string
		authCount    int
		gasPerAuth   float64
		gasPerByte   float64
		gweiTotalFee string
		err          error
	}, 1)

	// Send the transaction
	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		Rebroadcast: true,
		OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			if err != nil {
				s.logger.WithField("rpc", client.GetName()).Errorf("max bloating tx failed: %v", err)
				resultChan <- struct {
					gasUsed      uint64
					blockNumber  string
					authCount    int
					gasPerAuth   float64
					gasPerByte   float64
					gweiTotalFee string
					err          error
				}{0, "", 0, 0, 0, "", err}
				return
			}
			if receipt == nil {
				resultChan <- struct {
					gasUsed      uint64
					blockNumber  string
					authCount    int
					gasPerAuth   float64
					gasPerByte   float64
					gweiTotalFee string
					err          error
				}{0, "", 0, 0, 0, "", fmt.Errorf("no receipt received")}
				return
			}

			effectiveGasPrice := receipt.EffectiveGasPrice
			if effectiveGasPrice == nil {
				effectiveGasPrice = big.NewInt(0)
			}
			feeAmount := new(big.Int).Mul(effectiveGasPrice, big.NewInt(int64(receipt.GasUsed)))
			totalAmount := new(big.Int).Add(tx.Value(), feeAmount)
			wallet.SubBalance(totalAmount)

			gweiTotalFee := new(big.Int).Div(feeAmount, big.NewInt(1000000000))

			// Calculate efficiency metrics
			authCount := len(authorizations)
			gasPerAuth := float64(receipt.GasUsed) / float64(authCount)
			gasPerByte := gasPerAuth / EstimatedBytesPerAuth

			resultChan <- struct {
				gasUsed      uint64
				blockNumber  string
				authCount    int
				gasPerAuth   float64
				gasPerByte   float64
				gweiTotalFee string
				err          error
			}{receipt.GasUsed, receipt.BlockNumber.String(),
				authCount,
				gasPerAuth,
				gasPerByte,
				gweiTotalFee.String(), nil}
		},
		LogFn: func(client *spamoor.Client, retry int, rebroadcast int, err error) {
			logger := s.logger.WithField("rpc", client.GetName())
			if retry > 0 {
				logger = logger.WithField("retry", retry)
			}
			if rebroadcast > 0 {
				logger = logger.WithField("rebroadcast", rebroadcast)
			}
			if err != nil {
				logger.Errorf("failed sending max bloating tx: %v", err)
			} else if retry > 0 || rebroadcast > 0 {
				logger.Infof("successfully sent max bloating tx")
			}
		},
	})

	if err != nil {
		wallet.ResetPendingNonce(ctx, client)
		return 0, "", 0, 0, 0, "", fmt.Errorf("failed to send max bloating transaction: %w", err)
	}

	// Wait for transaction confirmation
	result := <-resultChan
	return result.gasUsed, result.blockNumber, result.authCount, result.gasPerAuth, result.gasPerByte, result.gweiTotalFee, result.err
}

// sendBatchedMaxBloatingTransactions sends multiple transactions when size limit is exceeded
func (s *Scenario) sendBatchedMaxBloatingTransactions(ctx context.Context, batches [][]types.SetCodeAuthorization, txCallData []byte, feeCap, tipCap *big.Int, amount *uint256.Int, toAddr common.Address, targetGasLimit uint64, wallet *spamoor.Wallet, client *spamoor.Client) (uint64, string, int, float64, float64, string, error) {
	// Aggregate results
	var totalGasUsed uint64
	var totalAuthCount int
	var totalFees *big.Int = big.NewInt(0)
	var lastBlockNumber string

	// Create result channels for all batches upfront
	resultChans := make([]chan struct {
		gasUsed      uint64
		blockNumber  string
		authCount    int
		gweiTotalFee string
		err          error
	}, len(batches))

	// Send all batches quickly with minimal delay to increase chance of same block inclusion
	for batchIndex, batch := range batches {
		// Create result channel for this batch
		resultChans[batchIndex] = make(chan struct {
			gasUsed      uint64
			blockNumber  string
			authCount    int
			gweiTotalFee string
			err          error
		}, 1)

		// Calculate appropriate gas limit for this batch based on authorization count
		// Each authorization needs ~26000 gas, plus some overhead for the transaction itself
		batchGasLimit := uint64(len(batch))*GasPerAuthorization + BaseTransferCost + uint64(len(txCallData)*GasPerCallDataByte)

		// Ensure we don't exceed the target limit per transaction
		maxGasPerTx := targetGasLimit
		if batchGasLimit > maxGasPerTx {
			batchGasLimit = maxGasPerTx
		}

		// Build the transaction for this batch
		txData, err := txbuilder.SetCodeTx(&txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Gas:       batchGasLimit,
			To:        &toAddr,
			Value:     amount,
			Data:      txCallData,
			AuthList:  batch,
		})
		if err != nil {
			return 0, "", 0, 0, 0, "", fmt.Errorf("failed to build batch %d transaction metadata: %w", batchIndex+1, err)
		}

		tx, err := wallet.BuildSetCodeTx(txData)
		if err != nil {
			return 0, "", 0, 0, 0, "", fmt.Errorf("failed to build batch %d transaction: %w", batchIndex+1, err)
		}

		// Send the transaction immediately without waiting for confirmation
		resultChan := resultChans[batchIndex]
		err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
			Client:      client,
			Rebroadcast: true,
			OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
				if err != nil {
					s.logger.WithField("rpc", client.GetName()).Errorf("batch %d tx failed: %v", batchIndex+1, err)
					resultChan <- struct {
						gasUsed      uint64
						blockNumber  string
						authCount    int
						gweiTotalFee string
						err          error
					}{0, "", 0, "", err}
					return
				}
				if receipt == nil {
					resultChan <- struct {
						gasUsed      uint64
						blockNumber  string
						authCount    int
						gweiTotalFee string
						err          error
					}{0, "", 0, "", fmt.Errorf("batch %d: no receipt received", batchIndex+1)}
					return
				}

				txFee := utils.GetTransactionFees(tx, receipt)
				totalFee := txFee.TotalFeeGwei()

				resultChan <- struct {
					gasUsed      uint64
					blockNumber  string
					authCount    int
					gweiTotalFee string
					err          error
				}{receipt.GasUsed, receipt.BlockNumber.String(), len(batch), totalFee.String(), nil}
			},
			LogFn: func(client *spamoor.Client, retry int, rebroadcast int, err error) {
				logger := s.logger.WithField("rpc", client.GetName())
				if err != nil {
					logger.Errorf("failed sending batch %d tx: %v", batchIndex+1, err)
				} else if retry > 0 || rebroadcast > 0 {
					logger.Debugf("successfully sent batch %d tx (retry/rebroadcast)", batchIndex+1)
				}
			},
		})

		if err != nil {
			wallet.ResetPendingNonce(ctx, client)
			return 0, "", 0, 0, 0, "", fmt.Errorf("failed to send batch %d transaction: %w", batchIndex+1, err)
		}

		// No delay between batches - send as fast as possible to ensure same block inclusion
	}

	// Now wait for all batch confirmations
	blockNumbers := make(map[string]int)         // Track which blocks contain our transactions
	batchDetails := make([]string, len(batches)) // Store details of each batch
	for batchIndex := range batches {
		result := <-resultChans[batchIndex]
		if result.err != nil {
			return 0, "", 0, 0, 0, "", fmt.Errorf("batch %d failed: %w", batchIndex+1, result.err)
		}

		// Aggregate successful results
		totalGasUsed += result.gasUsed
		totalAuthCount += result.authCount
		lastBlockNumber = result.blockNumber

		// Track block numbers
		blockNumbers[result.blockNumber]++

		// Parse and add fee
		if feeGwei, ok := new(big.Int).SetString(result.gweiTotalFee, 10); ok {
			totalFees.Add(totalFees, feeGwei)
		}

		// Store batch details
		batchGasInM := float64(result.gasUsed) / GasPerMillion
		gasPerAuthBatch := float64(result.gasUsed) / float64(result.authCount)
		gasPerByteBatch := gasPerAuthBatch / EstimatedBytesPerAuth

		// Calculate tx size based on authorizations
		txSize := s.calculateTransactionSize(len(batches[batchIndex]), len(txCallData))
		sizeKiB := float64(txSize) / BytesPerKiB

		batchDetails[batchIndex] = fmt.Sprintf("Batch %d/%d: %.2fM gas, %d auths, %.2f KiB, %.2f gas/auth, %.2f gas/byte, (block %s)",
			batchIndex+1, len(batches), batchGasInM, result.authCount, sizeKiB, gasPerAuthBatch, gasPerByteBatch, result.blockNumber)
	}

	// Calculate aggregate metrics
	gasPerAuth := float64(totalGasUsed) / float64(totalAuthCount)
	gasPerByte := gasPerAuth / EstimatedBytesPerAuth
	totalGasInM := float64(totalGasUsed) / GasPerMillion

	// Build block distribution summary
	var blockDistribution strings.Builder
	for blockNum, txCount := range blockNumbers {
		if blockDistribution.Len() > 0 {
			blockDistribution.WriteString(", ")
		}
		blockDistribution.WriteString(fmt.Sprintf("Block #%s: %d tx", blockNum, txCount))
	}

	// Create comprehensive summary log with decorative border
	s.logger.WithField("scenario", "eoa-delegation").Infof(`════════════════ BATCHED MAX BLOATING SUMMARY ════════════════
Individual Batches:
%s

Block Distribution: %s

Aggregate Metrics:
- Total Gas Used: %.2fM
- Total Authorizations: %d
- Gas per Auth: %.2f
- Gas per Byte: %.2f`,
		strings.Join(batchDetails, "\n"),
		blockDistribution.String(),
		totalGasInM,
		totalAuthCount,
		gasPerAuth,
		gasPerByte)

	return totalGasUsed, lastBlockNumber, totalAuthCount, gasPerAuth, gasPerByte, totalFees.String(), nil
}

// eoaWorker runs in a separate goroutine and writes funded EOAs to EOAs.json
// when the semaphore is open (green). It sleeps when the semaphore is closed (red).
func (s *Scenario) eoaWorker() {
	defer s.workerWg.Done()

	for {
		select {
		case <-s.workerDone:
			// Shutdown signal received
			return
		case <-s.workerSemaphore:
			// Semaphore is green, process the queue
			s.processEOAQueue()
		}
	}
}

// processEOAQueue drains the EOA queue and writes entries to EOAs.json
func (s *Scenario) processEOAQueue() {
	for {
		// Check if there are items in the queue
		s.eoaQueueMutex.Lock()
		if len(s.eoaQueue) == 0 {
			s.eoaQueueMutex.Unlock()
			return // Queue is empty, exit processing
		}

		// Dequeue all items (FIFO)
		eoasToWrite := make([]EOAEntry, len(s.eoaQueue))
		copy(eoasToWrite, s.eoaQueue)
		s.eoaQueue = s.eoaQueue[:0] // Clear the queue
		s.eoaQueueMutex.Unlock()

		// Write to file
		err := s.writeEOAsToFile(eoasToWrite)
		if err != nil {
			s.logger.Errorf("failed to write EOAs to file: %v", err)
			// Re-queue the items if write failed
			s.eoaQueueMutex.Lock()
			s.eoaQueue = append(eoasToWrite, s.eoaQueue...)
			s.eoaQueueMutex.Unlock()
			return
		}

	}
}

// writeEOAsToFile appends EOA entries to EOAs.json file
func (s *Scenario) writeEOAsToFile(eoas []EOAEntry) error {
	if len(eoas) == 0 {
		return nil
	}

	fileName := "EOAs.json"

	// Read existing entries if file exists
	var existingEntries []EOAEntry
	if data, err := os.ReadFile(fileName); err == nil {
		json.Unmarshal(data, &existingEntries)
	}

	// Append new entries
	allEntries := append(existingEntries, eoas...)

	// Write back to file
	data, err := json.MarshalIndent(allEntries, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal EOA entries: %w", err)
	}

	err = os.WriteFile(fileName, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write EOAs.json: %w", err)
	}

	return nil
}

// addEOAToQueue adds a funded EOA to the queue
func (s *Scenario) addEOAToQueue(address, privateKey string) {
	s.eoaQueueMutex.Lock()
	defer s.eoaQueueMutex.Unlock()

	entry := EOAEntry{
		Address:    address,
		PrivateKey: privateKey,
	}

	s.eoaQueue = append(s.eoaQueue, entry)
}

// openWorkerSemaphore opens the semaphore (green light) allowing the worker to process
func (s *Scenario) openWorkerSemaphore() {
	select {
	case s.workerSemaphore <- struct{}{}:
		// Semaphore opened successfully
	default:
		// Semaphore already open, do nothing
	}
}

// closeWorkerSemaphore closes the semaphore (red light) putting the worker to sleep
func (s *Scenario) closeWorkerSemaphore() {
	select {
	case <-s.workerSemaphore:
		// Semaphore closed successfully
	default:
		// Semaphore already closed, do nothing
	}
}

// shutdownWorker signals the worker to stop and waits for it to finish
func (s *Scenario) shutdownWorker() {
	close(s.workerDone)
	s.workerWg.Wait()
}

// calculateTransactionSize estimates the serialized size of a transaction with given authorizations
func (s *Scenario) calculateTransactionSize(authCount int, callDataSize int) int {
	// Estimation based on empirical data and RLP encoding structure:
	// - Base transaction overhead: ~200 bytes
	// - Each SetCodeAuthorization: ~94 bytes (based on actual observed data)
	// - Call data: variable size
	// - Additional RLP encoding overhead: ~50 bytes

	baseSize := TransactionBaseOverhead + callDataSize + TransactionExtraOverhead
	authSize := authCount * ActualBytesPerAuth
	return baseSize + authSize
}

// splitAuthorizationsBatches splits authorizations into batches that fit within transaction size limit
func (s *Scenario) splitAuthorizationsBatches(authorizations []types.SetCodeAuthorization, callDataSize int) [][]types.SetCodeAuthorization {
	if len(authorizations) == 0 {
		return [][]types.SetCodeAuthorization{authorizations}
	}

	// To get closer to 128KiB limit, we need to adjust our estimate
	// Using a safety factor of 0.95 to stay just under the limit
	targetSize := MaxTransactionSize * TransactionSizeSafetyFactor / 100 // Safety margin

	maxAuthsPerTx := (targetSize - TransactionBaseOverhead - callDataSize) / ActualBytesPerAuth
	if maxAuthsPerTx <= 0 {
		s.logger.Warnf("Transaction call data too large, using minimal batch size of 1")
		maxAuthsPerTx = 1
	}

	// If all authorizations fit in one transaction, return as single batch
	if len(authorizations) <= maxAuthsPerTx {
		estimatedSize := s.calculateTransactionSize(len(authorizations), callDataSize)
		s.logger.Infof("All %d authorizations fit in single transaction (estimated size: %d bytes)", len(authorizations), estimatedSize)
		return [][]types.SetCodeAuthorization{authorizations}
	}

	// Split into multiple batches
	var batches [][]types.SetCodeAuthorization
	for i := 0; i < len(authorizations); i += maxAuthsPerTx {
		end := i + maxAuthsPerTx
		if end > len(authorizations) {
			end = len(authorizations)
		}
		batch := authorizations[i:end]
		batches = append(batches, batch)
	}

	s.logger.Infof("Split %d authorizations into %d batches (max %d auths per batch, target size: %.2f KiB)",
		len(authorizations), len(batches), maxAuthsPerTx, float64(targetSize)/BytesPerKiB)
	return batches
}
