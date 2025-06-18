package sbeoadelegation

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
	EoaFile  string `yaml:"eoa_file"`
	LogTxs   bool   `yaml:"log_txs"`
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
}

var ScenarioName = "statebloat-eoa-delegation"
var ScenarioDefaultOptions = ScenarioOptions{
	BaseFee:  20,
	TipFee:   2,
	CodeAddr: "",
	EoaFile:  "",
	LogTxs:   false,
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
	flags.StringVar(&s.options.EoaFile, "eoa-file", ScenarioDefaultOptions.EoaFile, "File to write EOAs to")
	flags.BoolVar(&s.options.LogTxs, "log-txs", ScenarioDefaultOptions.LogTxs, "Log transactions")
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

	// In max-bloating mode, use 100 wallets for funding delegators
	s.walletPool.SetWalletCount(100)
	s.walletPool.SetRefillAmount(uint256.NewInt(0).Mul(uint256.NewInt(20), uint256.NewInt(1000000000000000000)))  // 20 ETH
	s.walletPool.SetRefillBalance(uint256.NewInt(0).Mul(uint256.NewInt(10), uint256.NewInt(1000000000000000000))) // 10 ETH

	// register well known wallets
	s.walletPool.AddWellKnownWallet(&spamoor.WellKnownWalletConfig{
		Name:          "bloater",
		RefillAmount:  uint256.NewInt(0).Mul(uint256.NewInt(20), uint256.NewInt(1000000000000000000)), // 20 ETH
		RefillBalance: uint256.NewInt(0).Mul(uint256.NewInt(10), uint256.NewInt(1000000000000000000)), // 10 ETH
	})

	// Initialize FIFO queue and worker for EOA management
	s.eoaQueue = make([]EOAEntry, 0)

	return nil
}

func (s *Scenario) Config() string {
	yamlBytes, _ := yaml.Marshal(&s.options)
	return string(yamlBytes)
}

func (s *Scenario) prepareDelegator(delegatorIndex uint64) (*spamoor.Wallet, error) {
	delegatorSeed := make([]byte, 8)
	binary.BigEndian.PutUint64(delegatorSeed, delegatorIndex)
	delegatorSeed = append(delegatorSeed, s.walletPool.GetWalletSeed()...)
	delegatorSeed = append(delegatorSeed, s.walletPool.GetRootWallet().GetWallet().GetAddress().Bytes()...)
	childKey := sha256.Sum256(delegatorSeed)
	return spamoor.NewWallet(fmt.Sprintf("%x", childKey))
}

func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Infof("starting max bloating mode: self-adjusting to target block gas limit, continuous operation")

	go s.eoaWorker(ctx)

	// Get the actual network block gas limit
	networkGasLimit, err := s.walletPool.GetTxPool().GetCurrentGasLimitWithInit()
	if err != nil {
		s.logger.Errorf("failed to get current gas limit: %v", err)
		return err
	}

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
			_, err := s.fundMaxBloatingDelegators(ctx, currentAuthorizations, blockCounter, networkGasLimit)
			if err != nil {
				s.logger.Errorf("failed to fund delegators for initial iteration: %v", err)
				time.Sleep(RetryDelay) // Wait before retry
				blockCounter--         // Retry the same iteration
				continue
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
		_, err = s.fundMaxBloatingDelegators(ctx, currentAuthorizations, blockCounter+1, networkGasLimit)
		if err != nil {
			s.logger.Errorf("failed to fund delegators for next iteration: %v", err)
			// Don't fail the entire loop, just log the error and continue
		}
	}
}

func (s *Scenario) fundMaxBloatingDelegators(ctx context.Context, targetCount int, iteration int, gasLimit uint64) (int64, error) {
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, "")
	if client == nil {
		return 0, fmt.Errorf("no client available for funding delegators")
	}

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
	wg := sync.WaitGroup{}

	delegatorIndexBase := uint64(iteration * FundingIterationOffset) // Large offset per iteration to avoid conflicts

	fundNextDelegator := func(ctx context.Context, txIdx uint64, onComplete func()) (*types.Transaction, *spamoor.Client, *spamoor.Wallet, error) {
		wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(txIdx))
		if wallet == nil {
			return nil, nil, nil, fmt.Errorf("no wallet available for funding")
		}

		client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, "")
		if client == nil {
			return nil, nil, nil, fmt.Errorf("no client available for funding delegators")
		}

		delegatorIndex := delegatorIndexBase + txIdx
		transactionSubmitted := false

		defer func() {
			if !transactionSubmitted {
				onComplete()
			}
		}()

		delegator, err := s.prepareDelegator(delegatorIndex)
		if err != nil {
			s.logger.Errorf("could not prepare delegator %v for funding: %v", delegatorIndex, err)
			return nil, nil, nil, err
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
			return nil, nil, nil, err
		}

		tx, err := wallet.BuildDynamicFeeTx(txData)
		if err != nil {
			s.logger.Errorf("failed to build funding transaction for delegator %d: %v", delegatorIndex, err)
			return nil, nil, nil, err
		}

		// Send funding transaction with no retries to avoid duplicates
		transactionSubmitted = true
		wg.Add(1)
		err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
			Client:      client,
			Rebroadcast: false, // No retries to avoid duplicates
			OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
				defer wg.Done()
				defer onComplete()

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

		return tx, client, wallet, err
	}

	// Calculate approximate transactions per block based on gas limit
	// Standard transfer = BaseTransferCost gas.
	maxTxsPerBlock := gasLimit / uint64(BaseTransferCost)

	scenario.RunTransactionScenario(ctx, scenario.TransactionScenarioOptions{
		TotalCount: uint64(targetCount),
		Throughput: maxTxsPerBlock * 2,
		MaxPending: maxTxsPerBlock * 2,
		WalletPool: s.walletPool,
		ProcessNextTxFn: func(ctx context.Context, txIdx uint64, onComplete func()) (func(), error) {
			logger := s.logger
			tx, client, wallet, err := fundNextDelegator(ctx, txIdx, onComplete)
			if client != nil {
				logger = logger.WithField("rpc", client.GetName())
			}
			if tx != nil {
				logger = logger.WithField("nonce", tx.Nonce())
			}
			if wallet != nil {
				logger = logger.WithField("wallet", s.walletPool.GetWalletName(wallet.GetAddress()))
			}

			return func() {
				if err != nil {
					logger.Warnf("could not send transaction: %v", err)
				} else if s.options.LogTxs {
					logger.Infof("sent tx #%6d: %v", txIdx+1, tx.Hash().String())
				} else {
					logger.Debugf("sent tx #%6d: %v", txIdx+1, tx.Hash().String())
				}
			}, err
		},
	})

	// Return the confirmed count
	wg.Wait()
	confirmed := atomic.LoadInt64(&confirmedCount)
	return confirmed, nil
}

func (s *Scenario) sendMaxBloatingTransaction(ctx context.Context, targetAuthorizations int, targetGasLimit uint64, blockCounter int) (uint64, string, int, float64, float64, string, error) {
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, "")
	if client == nil {
		return 0, "", 0, 0, 0, "", fmt.Errorf("no client available for sending max bloating transaction")
	}

	// Use bloater wallet
	wallet := s.walletPool.GetWellKnownWallet("bloater")
	if wallet == nil {
		return 0, "", 0, 0, 0, "", fmt.Errorf("no bloater wallet available")
	}

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
	toAddr := wallet.GetAddress()

	// No call data for max bloating transactions
	txCallData := []byte{}

	// Build the authorizations for maximum state bloat
	authorizations := s.buildMaxBloatingAuthorizations(targetAuthorizations, blockCounter)

	// Check transaction size and split into batches if needed
	batches := s.splitAuthorizationsBatches(authorizations, len(txCallData))

	if len(batches) == 1 {
		// Single transaction - use existing logic
		return s.sendSingleMaxBloatingTransaction(ctx, batches[0], txCallData, feeCap, tipCap, amount, toAddr, targetGasLimit, wallet, client)
	} else {
		// Multiple transactions needed - send them as a batch
		return s.sendBatchedMaxBloatingTransactions(ctx, batches, txCallData, feeCap, tipCap, amount, toAddr, targetGasLimit, wallet, client)
	}
}

// buildMaxBloatingAuthorizations builds the authorizations for the max bloating transaction
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

	chainId := s.walletPool.GetChainId()

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
			ChainID: *uint256.MustFromBig(chainId),
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

	wg := sync.WaitGroup{}
	wg.Add(1)

	var txreceipt *types.Receipt
	var txerr error

	// Send the transaction
	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		Rebroadcast: true,
		OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			txreceipt = receipt
			txerr = err
			wg.Done()
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
	wg.Wait()

	if txerr != nil {
		return 0, "", 0, 0, 0, "", fmt.Errorf("failed to send max bloating transaction: %w", txerr)
	}

	if txreceipt == nil {
		return 0, "", 0, 0, 0, "", fmt.Errorf("no receipt received")
	}

	effectiveGasPrice := txreceipt.EffectiveGasPrice
	if effectiveGasPrice == nil {
		effectiveGasPrice = big.NewInt(0)
	}
	feeAmount := new(big.Int).Mul(effectiveGasPrice, big.NewInt(int64(txreceipt.GasUsed)))
	totalAmount := new(big.Int).Add(tx.Value(), feeAmount)
	wallet.SubBalance(totalAmount)

	gweiTotalFee := new(big.Int).Div(feeAmount, big.NewInt(1000000000))

	// Calculate efficiency metrics
	authCount := len(authorizations)
	gasPerAuth := float64(txreceipt.GasUsed) / float64(authCount)
	gasPerByte := gasPerAuth / EstimatedBytesPerAuth

	return txreceipt.GasUsed, txreceipt.BlockNumber.String(), authCount, gasPerAuth, gasPerByte, gweiTotalFee.String(), nil
}

// sendBatchedMaxBloatingTransactions sends multiple transactions when size limit is exceeded
func (s *Scenario) sendBatchedMaxBloatingTransactions(ctx context.Context, batches [][]types.SetCodeAuthorization, txCallData []byte, feeCap, tipCap *big.Int, amount *uint256.Int, toAddr common.Address, targetGasLimit uint64, wallet *spamoor.Wallet, client *spamoor.Client) (uint64, string, int, float64, float64, string, error) {
	// Aggregate results
	var totalGasUsed uint64
	var totalAuthCount int
	var totalFees *big.Int = big.NewInt(0)
	var lastBlockNumber string

	// Create result channels for all batches upfront
	wg := sync.WaitGroup{}
	wg.Add(len(batches))

	txreceipts := make([]*types.Receipt, len(batches))
	txerrs := make([]error, len(batches))

	// Send all batches quickly with minimal delay to increase chance of same block inclusion
	for batchIndex, batch := range batches {

		err := func(batchIndex int, batch []types.SetCodeAuthorization) error {
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
				return fmt.Errorf("failed to build batch %d transaction metadata: %w", batchIndex+1, err)
			}

			tx, err := wallet.BuildSetCodeTx(txData)
			if err != nil {
				return fmt.Errorf("failed to build batch %d transaction: %w", batchIndex+1, err)
			}

			// Send the transaction immediately without waiting for confirmation
			err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
				Client:      client,
				Rebroadcast: true,
				OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
					txreceipts[batchIndex] = receipt
					txerrs[batchIndex] = err
					wg.Done()
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
				return fmt.Errorf("failed to send batch %d transaction: %w", batchIndex+1, err)
			}

			return nil
		}(batchIndex, batch)

		if err != nil {
			return 0, "", 0, 0, 0, "", fmt.Errorf("failed to send batch %d transaction: %w", batchIndex+1, err)
		}
	}

	// Now wait for all batch confirmations
	wg.Wait()

	blockNumbers := make(map[string]int)         // Track which blocks contain our transactions
	batchDetails := make([]string, len(batches)) // Store details of each batch
	for batchIndex := range batches {
		if txerrs[batchIndex] != nil {
			return 0, "", 0, 0, 0, "", fmt.Errorf("batch %d failed: %w", batchIndex+1, txerrs[batchIndex])
		}

		txreceipt := txreceipts[batchIndex]
		if txreceipt == nil {
			return 0, "", 0, 0, 0, "", fmt.Errorf("batch %d: no receipt received", batchIndex+1)
		}

		effectiveGasPrice := txreceipt.EffectiveGasPrice
		if effectiveGasPrice == nil {
			effectiveGasPrice = big.NewInt(0)
		}
		feeAmount := new(big.Int).Mul(effectiveGasPrice, big.NewInt(int64(txreceipt.GasUsed)))

		gweiTotalFee := new(big.Int).Div(feeAmount, big.NewInt(GweiPerEth))

		// Aggregate successful results
		totalGasUsed += txreceipt.GasUsed
		totalAuthCount += len(batches[batchIndex])
		lastBlockNumber = txreceipt.BlockNumber.String()

		// Track block numbers
		blockNumbers[txreceipt.BlockNumber.String()]++

		// Parse and add fee
		totalFees.Add(totalFees, gweiTotalFee)

		// Store batch details
		batchGasInM := float64(txreceipt.GasUsed) / GasPerMillion
		gasPerAuthBatch := float64(txreceipt.GasUsed) / float64(len(batches[batchIndex]))
		gasPerByteBatch := gasPerAuthBatch / EstimatedBytesPerAuth

		// Calculate tx size based on authorizations
		txSize := s.calculateTransactionSize(len(batches[batchIndex]), len(txCallData))
		sizeKiB := float64(txSize) / BytesPerKiB

		batchDetails[batchIndex] = fmt.Sprintf("Batch %d/%d: %.2fM gas, %d auths, %.2f KiB, %.2f gas/auth, %.2f gas/byte, (block %s)",
			batchIndex+1, len(batches), batchGasInM, len(batches[batchIndex]), sizeKiB, gasPerAuthBatch, gasPerByteBatch, txreceipt.BlockNumber.String())
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
func (s *Scenario) eoaWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			// Shutdown signal received
			return
		case <-time.After(30 * time.Second): // flush every 30 seconds
			s.processEOAQueue()
		}
	}
}

// processEOAQueue drains the EOA queue and writes entries to EOAs.json
func (s *Scenario) processEOAQueue() {
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

	if s.options.EoaFile != "" {
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
