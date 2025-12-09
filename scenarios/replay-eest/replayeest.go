package replayeest

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"

	"github.com/ethpandaops/spamoor/eestconv"
	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
)

// ScenarioOptions holds all configurable parameters
type ScenarioOptions struct {
	TotalCount      uint64  `yaml:"total_count"`
	Throughput      uint64  `yaml:"throughput"`
	MaxPending      uint64  `yaml:"max_pending"`
	MaxWallets      uint64  `yaml:"max_wallets"`
	Rebroadcast     uint64  `yaml:"rebroadcast"`
	BaseFee         float64 `yaml:"base_fee"`
	TipFee          float64 `yaml:"tip_fee"`
	Timeout         string  `yaml:"timeout"`
	PayloadTimeout  string  `yaml:"payload_timeout"` // Timeout for each payload (default: 30m)
	ClientGroup     string  `yaml:"client_group"`
	PayloadFile     string  `yaml:"payload_file"`
	FixturesRelease string  `yaml:"fixtures_release"` // URL to .tar.gz file containing EEST fixtures
	FixturesPattern string  `yaml:"fixtures_pattern"` // Regex pattern to include fixtures by path/name
	FixturesExclude string  `yaml:"fixtures_exclude"` // Regex pattern to exclude fixtures by path/name
	StartOffset     uint64  `yaml:"start_offset"`
	SkipPostChecks  bool    `yaml:"skip_postchecks"`
	LogTxs          bool    `yaml:"log_txs"`
}

// Scenario implements the replay-eest scenario
type Scenario struct {
	options        ScenarioOptions
	logger         *logrus.Entry
	walletPool     *spamoor.WalletPool
	payloadTimeout time.Duration

	// Loaded payloads
	payloads []Payload

	// Wallet locking
	walletLock     sync.Mutex
	lockedWallets  map[common.Address]bool
	walletUnlockCh chan common.Address

	// Stats tracking
	statsLock          sync.Mutex
	currentPayloadIdx  uint64
	successCount       uint64
	failInvalidCount   uint64 // Failed due to invalid tx (wallet acquisition, build errors)
	failRevertedCount  uint64 // Failed due to tx revert
	failPostcheckCount uint64 // Failed due to postcheck failures
	failTimeoutCount   uint64 // Failed due to payload timeout
}

var ScenarioName = "replay-eest"
var ScenarioDefaultOptions = ScenarioOptions{
	TotalCount:      0,
	Throughput:      1,
	MaxPending:      10,
	MaxWallets:      100,
	Rebroadcast:     1,
	BaseFee:         20,
	TipFee:          2,
	Timeout:         "",
	PayloadTimeout:  "30m",
	ClientGroup:     "",
	PayloadFile:     "",
	FixturesRelease: "",
	FixturesPattern: "",
	FixturesExclude: "",
	StartOffset:     0,
	SkipPostChecks:  false,
	LogTxs:          false,
}

var ScenarioDescriptor = scenario.Descriptor{
	Name:           ScenarioName,
	Description:    "Replay EEST test fixtures from intermediate representation",
	DefaultOptions: ScenarioDefaultOptions,
	NewScenario:    newScenario,
}

func newScenario(logger logrus.FieldLogger) scenario.Scenario {
	return &Scenario{
		options:        ScenarioDefaultOptions,
		logger:         logger.WithField("scenario", ScenarioName),
		lockedWallets:  make(map[common.Address]bool),
		walletUnlockCh: make(chan common.Address, 1000),
	}
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
	flags.Uint64VarP(&s.options.TotalCount, "count", "c", ScenarioDefaultOptions.TotalCount,
		"Total number of test cases to run (0 = all)")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t", ScenarioDefaultOptions.Throughput,
		"Number of test cases to run per slot")
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending,
		"Maximum number of pending test cases")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets,
		"Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", ScenarioDefaultOptions.Rebroadcast,
		"Enable rebroadcasting (0 = disabled)")
	flags.Float64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee,
		"Base fee in gwei (0 = use suggested)")
	flags.Float64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee,
		"Tip fee in gwei (0 = use suggested)")
	flags.StringVar(&s.options.Timeout, "timeout", ScenarioDefaultOptions.Timeout,
		"Timeout for the scenario (e.g., '1h', '30m')")
	flags.StringVar(&s.options.PayloadTimeout, "payload-timeout", ScenarioDefaultOptions.PayloadTimeout,
		"Timeout for each payload before cancellation (e.g., '30m', '1h')")
	flags.StringVar(&s.options.ClientGroup, "client-group", ScenarioDefaultOptions.ClientGroup,
		"Client group to use for transactions")
	flags.StringVarP(&s.options.PayloadFile, "payload", "f", ScenarioDefaultOptions.PayloadFile,
		"Path or URL to the YAML payload file")
	flags.StringVar(&s.options.FixturesRelease, "fixtures-release", ScenarioDefaultOptions.FixturesRelease,
		"URL or path to a .tar.gz file containing EEST fixtures (alternative to --payload)")
	flags.StringVar(&s.options.FixturesPattern, "fixtures-pattern", ScenarioDefaultOptions.FixturesPattern,
		"Regex pattern to include fixtures by path/name (used with --fixtures-release)")
	flags.StringVar(&s.options.FixturesExclude, "fixtures-exclude", ScenarioDefaultOptions.FixturesExclude,
		"Regex pattern to exclude fixtures by path/name (used with --fixtures-release)")
	flags.Uint64Var(&s.options.StartOffset, "start-offset", ScenarioDefaultOptions.StartOffset,
		"Start offset for the scenario")
	flags.BoolVar(&s.options.SkipPostChecks, "skip-postchecks", ScenarioDefaultOptions.SkipPostChecks,
		"Skip post-check validation")
	flags.BoolVar(&s.options.LogTxs, "log-txs", ScenarioDefaultOptions.LogTxs,
		"Log individual transactions")
	return nil
}

func (s *Scenario) Init(options *scenario.Options) error {
	s.walletPool = options.WalletPool

	// Parse YAML config if provided
	if options.Config != "" {
		err := scenario.ParseAndValidateConfig(&ScenarioDescriptor, options.Config, &s.options, s.logger)
		if err != nil {
			return err
		}
	}

	// Validate options: need either PayloadFile or FixturesRelease
	if s.options.PayloadFile == "" && s.options.FixturesRelease == "" {
		return fmt.Errorf("either payload file or fixtures release is required (use --payload or --fixtures-release)")
	}

	if s.options.PayloadFile != "" && s.options.FixturesRelease != "" {
		return fmt.Errorf("cannot specify both payload file and fixtures release")
	}

	if s.options.PayloadTimeout != "" {
		var err error
		s.payloadTimeout, err = time.ParseDuration(s.options.PayloadTimeout)
		if err != nil {
			return fmt.Errorf("invalid payload timeout: %w", err)
		}
	}

	// Configure wallet count: we need enough wallets for parallel test cases
	minWallets := 5 * s.options.MaxPending
	if s.options.MaxWallets > 0 && s.options.MaxWallets < minWallets {
		s.logger.Warnf("max-wallets (%d) is less than recommended (%d), may cause contention",
			s.options.MaxWallets, minWallets)
		s.walletPool.SetWalletCount(s.options.MaxWallets)
	} else if s.options.MaxWallets > 0 {
		s.walletPool.SetWalletCount(s.options.MaxWallets)
	} else {
		s.walletPool.SetWalletCount(minWallets)
	}

	return nil
}

func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Infof("starting scenario: %s", ScenarioName)
	defer s.logger.Infof("scenario %s finished", ScenarioName)

	// Load payloads from appropriate source
	var payloads []Payload
	var err error

	if s.options.FixturesRelease != "" {
		// Load from fixtures release (.tar.gz)
		payloads, err = s.loadFromFixturesRelease(ctx, s.options.FixturesRelease, s.options.FixturesPattern, s.options.FixturesExclude)
		if err != nil {
			return fmt.Errorf("failed to load from fixtures release: %w", err)
		}
	} else {
		// Load from payload file
		payloads, err = s.loadPayloads(s.options.PayloadFile)
		if err != nil {
			return fmt.Errorf("failed to load payloads: %w", err)
		}
	}

	s.payloads = payloads
	s.logger.Infof("loaded %d test payloads", len(s.payloads))

	// Check context after loading
	if err := ctx.Err(); err != nil {
		return err
	}

	// Calculate max senders needed across all payloads
	maxSenders := 0
	for _, p := range s.payloads {
		senderCount := p.GetSenderCount()
		if senderCount > maxSenders {
			maxSenders = senderCount
		}
	}

	// Start wallet unlock processor
	go s.processWalletUnlocks(ctx)

	// Start stats summary printer
	go s.printStatsSummary(ctx)

	totalCount := s.options.TotalCount
	if totalCount == 0 || totalCount > uint64(len(s.payloads)) {
		totalCount = uint64(len(s.payloads))
	}

	maxPending := s.options.MaxPending
	if maxPending == 0 {
		maxPending = 10
	}

	var timeout time.Duration
	if s.options.Timeout != "" {
		var err error
		timeout, err = time.ParseDuration(s.options.Timeout)
		if err != nil {
			return fmt.Errorf("invalid timeout: %w", err)
		}
	}

	err = scenario.RunTransactionScenario(ctx, scenario.TransactionScenarioOptions{
		TotalCount: totalCount,
		Throughput: s.options.Throughput,
		MaxPending: maxPending,
		Timeout:    timeout,
		WalletPool: s.walletPool,
		Logger:     s.logger,
		ProcessNextTxFn: func(ctx context.Context, params *scenario.ProcessNextTxParams) error {
			return s.processTestCase(ctx, params)
		},
	})

	// Print final summary
	s.statsLock.Lock()
	totalFailed := s.failInvalidCount + s.failRevertedCount + s.failPostcheckCount + s.failTimeoutCount
	s.logger.Infof("final summary: current=%d, success=%d, failed=%d (invalid=%d, reverted=%d, postcheck=%d, timeout=%d), total=%d",
		s.currentPayloadIdx, s.successCount, totalFailed,
		s.failInvalidCount, s.failRevertedCount, s.failPostcheckCount, s.failTimeoutCount, len(s.payloads))
	s.statsLock.Unlock()

	return err
}

func (s *Scenario) processTestCase(ctx context.Context, params *scenario.ProcessNextTxParams) error {
	payloadIdx := (int(params.TxIdx) + int(s.options.StartOffset)) % len(s.payloads)
	payload := s.payloads[payloadIdx]

	// Update current payload index
	s.statsLock.Lock()
	s.currentPayloadIdx = uint64(payloadIdx)
	s.statsLock.Unlock()

	logger := s.logger.WithField("payload", payloadIdx)

	// Determine how many wallets we need
	senderCount := payload.GetSenderCount()
	totalWalletsNeeded := senderCount + 1 // +1 for deployer

	// Acquire all wallets atomically
	wallets, err := s.acquireWallets(ctx, totalWalletsNeeded)
	if err != nil {
		params.NotifySubmitted()
		params.OrderedLogCb(func() {
			logger.Warnf("failed to acquire wallets: %v", err)
		})
		s.statsLock.Lock()
		s.failInvalidCount++
		s.statsLock.Unlock()
		return err
	}

	// Notify that we've started processing
	params.NotifySubmitted()

	// Execute the test case
	result := s.executeTestCase(ctx, logger, payload, wallets)

	// Release all wallets
	s.releaseWallets(wallets)

	// Update stats
	s.statsLock.Lock()
	if result.timedOut {
		s.failTimeoutCount++
	} else if result.err != nil {
		if result.reverted {
			s.failRevertedCount++
		} else {
			s.failInvalidCount++
		}
	} else if !s.options.SkipPostChecks && len(result.postCheckFailures) > 0 {
		s.failPostcheckCount++
	} else {
		s.successCount++
	}
	s.statsLock.Unlock()

	// Log result
	params.OrderedLogCb(func() {
		if result.timedOut {
			logger.WithField("test", payload.Name).Warnf("test case timed out after payload timeout exceeded")
		} else if result.err != nil {
			logger.WithField("test", payload.Name).Warnf("test case failed: %w", result.err)
		} else if len(result.postCheckFailures) > 0 {
			errMsg := make([]string, 0, len(result.postCheckFailures))
			errMsg = append(errMsg, result.postCheckFailures...)
			logger.WithField("test", payload.Name).Warnf("test case completed with %d post-check failures: %s", len(result.postCheckFailures), strings.Join(errMsg, ", "))
		} else {
			if s.options.LogTxs {
				logger.Infof("test case #%d completed successfully", params.TxIdx+1)
			} else {
				logger.Debugf("test case #%d completed successfully", params.TxIdx+1)
			}
		}
	})

	return result.err
}

type testCaseResult struct {
	err               error
	reverted          bool // True if the error was due to tx revert
	timedOut          bool // True if the payload timed out
	postCheckFailures []string
}

// senderGasCosts tracks the fixture and actual gas costs per sender for balance comparison
type senderGasCosts struct {
	fixtureGasCost *big.Int // Gas * FixtureBaseFee for each tx
	actualGasCost  *big.Int // Actual gas spent on chain
}

func (s *Scenario) executeTestCase(
	ctx context.Context,
	logger *logrus.Entry,
	payload Payload,
	wallets []*spamoor.Wallet,
) testCaseResult {
	result := testCaseResult{}

	// wallets[0] is the deployer, wallets[1:] are senders
	deployer := wallets[0]
	senderWallets := wallets[1:]

	// Get a client
	client := s.walletPool.GetClient(spamoor.WithClientGroup(s.options.ClientGroup))
	if client == nil {
		result.err = fmt.Errorf("no client available")
		return result
	}

	// Record initial balances for senders (for post-check comparison)
	initialBalances := make(map[int]*big.Int)
	for i, wallet := range senderWallets {
		initialBalances[i+1] = new(big.Int).Set(wallet.GetBalance())
	}

	// Track gas costs per sender for balance scaling
	gasCosts := make(map[int]*senderGasCosts)

	// Build address mapping for placeholder replacement
	// Contract addresses are computed from deployer nonce
	deployerNonce := deployer.GetNonce()

	addressMap := s.buildAddressMap(deployer, senderWallets, deployerNonce, payload)

	// Get suggested fees
	feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(client, s.options.BaseFee, s.options.TipFee)
	if err != nil {
		result.err = fmt.Errorf("failed to get suggested fees: %w", err)
		return result
	}

	// Group consecutive transactions by sender for batching
	// We need to preserve order: all deployer txs first, then sender txs in original order
	// but we can batch consecutive txs from the same sender

	var deployerTxs []Tx
	var senderTxGroups []struct {
		from string
		txs  []Tx
	}

	for _, tx := range payload.Txs {
		if tx.From == "deployer" {
			deployerTxs = append(deployerTxs, tx)
		} else {
			// Check if we can append to the last group (same sender)
			if len(senderTxGroups) > 0 && senderTxGroups[len(senderTxGroups)-1].from == tx.From {
				senderTxGroups[len(senderTxGroups)-1].txs = append(senderTxGroups[len(senderTxGroups)-1].txs, tx)
			} else {
				// Start a new group
				senderTxGroups = append(senderTxGroups, struct {
					from string
					txs  []Tx
				}{from: tx.From, txs: []Tx{tx}})
			}
		}
	}

	// Execute deployer transactions first
	if len(deployerTxs) > 0 {
		reverted, timedOut, err := s.executeDeployerTxs(ctx, logger, deployer, client, deployerTxs, addressMap, senderWallets, feeCap, tipCap)
		if err != nil {
			result.err = fmt.Errorf("deployer transactions failed: %w", err)
			result.reverted = reverted
			result.timedOut = timedOut
			return result
		}
	}

	// Execute sender transaction groups in order (preserving original tx order)
	for _, group := range senderTxGroups {
		senderIdx := parseSenderIndex(group.from)
		if senderIdx < 1 || senderIdx > len(senderWallets) {
			result.err = fmt.Errorf("invalid sender index: %s", group.from)
			return result
		}
		senderWallet := senderWallets[senderIdx-1]

		fixtureGas, actualGas, reverted, timedOut, err := s.executeSenderTxs(ctx, logger, senderWallet, client, group.txs, addressMap, senderWallets, feeCap, tipCap)
		if err != nil {
			result.err = fmt.Errorf("sender[%d] transactions failed: %w", senderIdx, err)
			result.reverted = reverted
			result.timedOut = timedOut
			return result
		}

		// Accumulate gas costs for this sender
		if gasCosts[senderIdx] == nil {
			gasCosts[senderIdx] = &senderGasCosts{
				fixtureGasCost: big.NewInt(0),
				actualGasCost:  big.NewInt(0),
			}
		}
		gasCosts[senderIdx].fixtureGasCost.Add(gasCosts[senderIdx].fixtureGasCost, fixtureGas)
		gasCosts[senderIdx].actualGasCost.Add(gasCosts[senderIdx].actualGasCost, actualGas)
	}

	// Perform post-checks
	if !s.options.SkipPostChecks {
		result.postCheckFailures = s.performPostChecks(ctx, logger, client, senderWallets, payload, addressMap, initialBalances, gasCosts)
	}

	return result
}

func (s *Scenario) executeDeployerTxs(
	ctx context.Context,
	logger *logrus.Entry,
	deployer *spamoor.Wallet,
	client *spamoor.Client,
	txs []Tx,
	addressMap map[string]common.Address,
	senderWallets []*spamoor.Wallet,
	feeCap, tipCap *big.Int,
) (bool, bool, error) {
	// Reset nonces if needed
	if err := deployer.ResetNoncesIfNeeded(ctx, client); err != nil {
		return false, false, fmt.Errorf("failed to reset deployer nonces: %w", err)
	}

	// Build, submit, and wait for each transaction one by one
	receipts := make([]*types.Receipt, len(txs))
	errors := make([]error, len(txs))
	pendingTxs := 0
	pendingChan := make(chan struct{}, 1)
	txObjects := make([]*types.Transaction, len(txs))

	var timer *time.Timer
	if s.payloadTimeout > 0 {
		timer = time.NewTimer(s.payloadTimeout)
	} else {
		timer = time.NewTimer(time.Hour * 2)
	}
	defer timer.Stop()

	for i, tx := range txs {
		// Build transaction
		signedTx, err := s.buildTransaction(ctx, deployer, client, tx, addressMap, senderWallets, feeCap, tipCap)
		if err != nil {
			return false, false, fmt.Errorf("failed to build deployer tx %d: %w", i, err)
		}

		txObjects[i] = signedTx

		// Submit transaction with OnComplete callback
		pendingTxs++
		txIdx := i
		err = s.walletPool.GetTxPool().SendTransaction(ctx, deployer, signedTx, &spamoor.SendTransactionOptions{
			Client:      client,
			ClientGroup: s.options.ClientGroup,
			Rebroadcast: s.options.Rebroadcast > 0,
			OnComplete: func(_ *types.Transaction, receipt *types.Receipt, err error) {
				receipts[txIdx] = receipt
				errors[txIdx] = err
				pendingChan <- struct{}{}
			},
		})
		if err != nil {
			return false, false, fmt.Errorf("failed to submit deployer tx %d: %w", i, err)
		}
	}

	// Wait for all transactions to complete or timeout
waitloop:
	select {
	case <-pendingChan:
		pendingTxs--
		if pendingTxs == 0 {
			break waitloop
		}
	case <-timer.C:
		err := s.replaceWithDummyTxs(ctx, txObjects, deployer, client, feeCap, tipCap)
		return false, true, fmt.Errorf("deployer transactions timed out (replaced: %w)", err)
	case <-ctx.Done():
		return false, false, ctx.Err()
	}

	// Check for any errors
	for i, err := range errors {
		if err != nil {
			return false, false, fmt.Errorf("deployer tx %d failed: %w", i, err)
		}
	}

	reverted := false
	for _, receipt := range receipts {
		if receipt == nil || receipt.Status == 0 {
			reverted = true
		}
	}

	if s.options.LogTxs {
		logger.Debugf("deployer transactions confirmed: %d txs", len(txs))
	}

	return reverted, false, nil
}

func (s *Scenario) executeSenderTxs(
	ctx context.Context,
	logger *logrus.Entry,
	sender *spamoor.Wallet,
	client *spamoor.Client,
	txs []Tx,
	addressMap map[string]common.Address,
	senderWallets []*spamoor.Wallet,
	feeCap, tipCap *big.Int,
) (fixtureGasCost, actualGasCost *big.Int, reverted bool, timedOut bool, err error) {
	fixtureGasCost = big.NewInt(0)
	actualGasCost = big.NewInt(0)

	// Reset nonces if needed
	if err := sender.ResetNoncesIfNeeded(ctx, client); err != nil {
		return nil, nil, false, false, fmt.Errorf("failed to reset sender nonces: %w", err)
	}

	// Build, submit, and wait for each transaction one by one
	receipts := make([]*types.Receipt, len(txs))
	errors := make([]error, len(txs))
	pendingTxs := 0
	pendingChan := make(chan struct{}, 1)
	txObjects := make([]*types.Transaction, len(txs))

	var timer *time.Timer
	if s.payloadTimeout > 0 {
		timer = time.NewTimer(s.payloadTimeout)
	} else {
		timer = time.NewTimer(time.Hour * 2)
	}
	defer timer.Stop()

	for i, tx := range txs {
		// Build transaction
		signedTx, err := s.buildTransaction(ctx, sender, client, tx, addressMap, senderWallets, feeCap, tipCap)
		if err != nil {
			return nil, nil, false, false, fmt.Errorf("failed to build sender tx %d: %w", i, err)
		}

		txObjects[i] = signedTx

		// Submit transaction with OnComplete callback
		pendingTxs++
		txIdx := i
		err = s.walletPool.GetTxPool().SendTransaction(ctx, sender, signedTx, &spamoor.SendTransactionOptions{
			Client:      client,
			ClientGroup: s.options.ClientGroup,
			Rebroadcast: s.options.Rebroadcast > 0,
			OnComplete: func(_ *types.Transaction, receipt *types.Receipt, err error) {
				receipts[txIdx] = receipt
				errors[txIdx] = err
				pendingChan <- struct{}{}
			},
		})
		if err != nil {
			return nil, nil, false, false, fmt.Errorf("failed to submit sender tx %d: %w", i, err)
		}
	}

	// Wait for all transactions to complete
waitloop:
	select {
	case <-pendingChan:
		pendingTxs--
		if pendingTxs == 0 {
			break waitloop
		}
	case <-timer.C:
		err := s.replaceWithDummyTxs(ctx, txObjects, sender, client, feeCap, tipCap)
		return nil, nil, false, true, fmt.Errorf("sender transactions timed out (replaced: %w)", err)
	case <-ctx.Done():
		return nil, nil, false, false, ctx.Err()
	}

	// Check for errors and calculate gas costs from receipts
	for i, receipt := range receipts {
		if errors[i] != nil {
			return nil, nil, false, false, fmt.Errorf("sender tx %d failed: %w", i, errors[i])
		}
		if receipt == nil || receipt.Status == 0 {
			reverted = true
		}
		if receipt != nil {
			// Fixture gas cost = gasUsed * fixtureBaseFee
			txFixtureCost := new(big.Int).Mul(big.NewInt(int64(receipt.GasUsed)), big.NewInt(int64(txs[i].GasPrice)))
			fixtureGasCost.Add(fixtureGasCost, txFixtureCost)

			// Actual gas cost = gasUsed * effectiveGasPrice
			txActualCost := new(big.Int).Mul(big.NewInt(int64(receipt.GasUsed)), feeCap)
			actualGasCost.Add(actualGasCost, txActualCost)
		}
	}

	if s.options.LogTxs {
		logger.Debugf("sender transactions confirmed: %d txs", len(txs))
	}

	return fixtureGasCost, actualGasCost, reverted, false, nil
}

func (s *Scenario) buildTransaction(
	ctx context.Context,
	wallet *spamoor.Wallet,
	client *spamoor.Client,
	tx Tx,
	addressMap map[string]common.Address,
	senderWallets []*spamoor.Wallet,
	feeCap, tipCap *big.Int,
) (*types.Transaction, error) {
	// Replace placeholders in data
	data, err := s.replacePlaceholders(tx.Data, addressMap)
	if err != nil {
		return nil, fmt.Errorf("failed to replace placeholders in data: %w", err)
	}

	// Parse data as hex
	var dataBytes []byte
	if data != "" && data != "0x" {
		dataBytes, err = hex.DecodeString(strings.TrimPrefix(data, "0x"))
		if err != nil {
			return nil, fmt.Errorf("failed to decode data: %w", err)
		}
	}

	// Determine target address
	var toAddr *common.Address
	if tx.To != "" {
		toStr, err := s.replacePlaceholders(tx.To, addressMap)
		if err != nil {
			return nil, fmt.Errorf("failed to replace placeholders in to: %w", err)
		}
		addr := common.HexToAddress(toStr)
		toAddr = &addr
	}

	// Parse value
	value := uint256.NewInt(0)
	if tx.Value != "" && tx.Value != "0x0" && tx.Value != "0x00" {
		valueBig, ok := new(big.Int).SetString(strings.TrimPrefix(tx.Value, "0x"), 16)
		if !ok {
			return nil, fmt.Errorf("failed to parse value: %s", tx.Value)
		}
		value = uint256.MustFromBig(valueBig)
	}

	// Determine gas limit
	gasLimit := tx.Gas
	if gasLimit == 0 {
		gasLimit = 1000000 // default
	}

	// Convert access list if present
	var accessList types.AccessList
	if len(tx.AccessList) > 0 {
		accessList = make(types.AccessList, len(tx.AccessList))
		for i, entry := range tx.AccessList {
			// Replace placeholders in access list address
			addrStr, err := s.replacePlaceholders(entry.Address, addressMap)
			if err != nil {
				return nil, fmt.Errorf("failed to replace placeholders in access list address: %w", err)
			}
			addr := common.HexToAddress(addrStr)

			storageKeys := make([]common.Hash, len(entry.StorageKeys))
			for j, key := range entry.StorageKeys {
				storageKeys[j] = common.HexToHash(key)
			}

			accessList[i] = types.AccessTuple{
				Address:     addr,
				StorageKeys: storageKeys,
			}
		}
	}

	// Build transaction based on type
	var signedTx *types.Transaction

	switch tx.Type {
	case 0: // Legacy
		txData, err := txbuilder.LegacyTx(&txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			Gas:       gasLimit,
			To:        toAddr,
			Value:     value,
			Data:      dataBytes,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to build legacy tx: %w", err)
		}
		signedTx, err = wallet.BuildLegacyTx(txData)
		if err != nil {
			return nil, fmt.Errorf("failed to build legacy tx: %w", err)
		}
	case 1: // Access List
		txData, err := txbuilder.AccessListTx(&txbuilder.TxMetadata{
			GasFeeCap:  uint256.MustFromBig(feeCap),
			Gas:        gasLimit,
			To:         toAddr,
			Value:      value,
			Data:       dataBytes,
			AccessList: accessList,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to build access list tx: %w", err)
		}
		signedTx, err = wallet.BuildAccessListTx(txData)
		if err != nil {
			return nil, fmt.Errorf("failed to build access list tx: %w", err)
		}
	case 2: // EIP-1559
		txData, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
			GasFeeCap:  uint256.MustFromBig(feeCap),
			GasTipCap:  uint256.MustFromBig(tipCap),
			Gas:        gasLimit,
			To:         toAddr,
			Value:      value,
			Data:       dataBytes,
			AccessList: accessList,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to build dynamic fee tx: %w", err)
		}
		signedTx, err = wallet.BuildDynamicFeeTx(txData)
		if err != nil {
			return nil, fmt.Errorf("failed to build dynamic fee tx: %w", err)
		}
	case 3: // Blob transaction (EIP-4844)
		if toAddr == nil {
			return nil, fmt.Errorf("blob transactions require a 'to' address")
		}

		// Use a high blob fee cap (same approach as other scenarios)
		blobFee := new(big.Int).Mul(feeCap, big.NewInt(1000000000))

		// Generate random blobs for the specified count
		blobRefs := make([][]string, tx.BlobCount)
		for i := 0; i < tx.BlobCount; i++ {
			blobRefs[i] = []string{"random"}
		}

		blobTx, err := txbuilder.BuildBlobTx(&txbuilder.TxMetadata{
			GasFeeCap:  uint256.MustFromBig(feeCap),
			GasTipCap:  uint256.MustFromBig(tipCap),
			BlobFeeCap: uint256.MustFromBig(blobFee),
			Gas:        gasLimit,
			To:         toAddr,
			Value:      value,
			Data:       dataBytes,
			AccessList: accessList,
		}, blobRefs)
		if err != nil {
			return nil, fmt.Errorf("failed to build blob tx: %w", err)
		}
		signedTx, err = wallet.BuildBlobTx(blobTx)
		if err != nil {
			return nil, fmt.Errorf("failed to sign blob tx: %w", err)
		}
	case 4: // SetCode transaction (EIP-7702)
		if toAddr == nil {
			return nil, fmt.Errorf("setcode transactions require a 'to' address")
		}

		// Build authorization list
		authList := make([]types.SetCodeAuthorization, len(tx.AuthorizationList))
		for i, authItem := range tx.AuthorizationList {
			// Resolve the target address (may have placeholders like $contract[1])
			addrStr, err := s.replacePlaceholders(authItem.Address, addressMap)
			if err != nil {
				return nil, fmt.Errorf("failed to replace placeholders in auth address: %w", err)
			}
			targetAddr := common.HexToAddress(addrStr)

			auth := types.SetCodeAuthorization{
				ChainID: *uint256.NewInt(authItem.ChainID),
				Address: targetAddr,
				Nonce:   authItem.Nonce,
			}

			// Check if signer is a sender placeholder (e.g., "sender[1]")
			if strings.HasPrefix(authItem.Signer, "sender[") {
				senderIdx := parseSenderIndex(authItem.Signer)
				if senderIdx < 1 || senderIdx > len(senderWallets) {
					return nil, fmt.Errorf("invalid sender index in authorization: %s", authItem.Signer)
				}
				signerWallet := senderWallets[senderIdx-1]

				// Sign the authorization with the sender's private key
				auth, err = types.SignSetCode(signerWallet.GetPrivateKey(), auth)
				if err != nil {
					return nil, fmt.Errorf("failed to sign authorization: %w", err)
				}
			} else {
				// For non-sender signers, use the original signature from the fixture
				auth.V = parseHexUint8(authItem.V)
				auth.R = *parseHexUint256(authItem.R)
				auth.S = *parseHexUint256(authItem.S)
			}

			authList[i] = auth
		}

		setCodeTx, err := txbuilder.SetCodeTx(&txbuilder.TxMetadata{
			GasFeeCap:  uint256.MustFromBig(feeCap),
			GasTipCap:  uint256.MustFromBig(tipCap),
			Gas:        gasLimit,
			To:         toAddr,
			Value:      value,
			Data:       dataBytes,
			AccessList: accessList,
			AuthList:   authList,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to build setcode tx: %w", err)
		}
		signedTx, err = wallet.BuildSetCodeTx(setCodeTx)
		if err != nil {
			return nil, fmt.Errorf("failed to sign setcode tx: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported transaction type: %d", tx.Type)
	}

	return signedTx, nil
}

// replaceWithDummyTxs
func (s *Scenario) replaceWithDummyTxs(ctx context.Context, txs []*types.Transaction, wallet *spamoor.Wallet, client *spamoor.Client, feeCap, tipCap *big.Int) error {
	doubleFeeCap := new(big.Int).Mul(feeCap, big.NewInt(2))
	doubleTipCap := new(big.Int).Mul(tipCap, big.NewInt(2))
	walletAddr := wallet.GetAddress()

	replaceTxs := make([]*types.Transaction, len(txs))
	for i, tx := range txs {
		txData, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(doubleFeeCap),
			GasTipCap: uint256.MustFromBig(doubleTipCap),
			Gas:       25000,
			To:        &walletAddr,
			Value:     uint256.NewInt(0),
			Data:      []byte{},
		})
		if err != nil {
			return fmt.Errorf("failed to build dynamic fee tx: %w", err)
		}
		signedTx, err := wallet.ReplaceDynamicFeeTx(txData, tx.Nonce())
		if err != nil {
			return fmt.Errorf("failed to build dynamic fee tx: %w", err)
		}

		replaceTxs[i] = signedTx
	}

	_, err := s.walletPool.GetTxPool().SendTransactionBatch(ctx, wallet, replaceTxs, &spamoor.BatchOptions{
		PendingLimit: 100,
		ClientPool:   s.walletPool.GetClientPool(),
		ClientGroup:  s.options.ClientGroup,
		MaxRetries:   3,
	})

	return err
}

func (s *Scenario) buildAddressMap(
	deployer *spamoor.Wallet,
	senders []*spamoor.Wallet,
	deployerNonce uint64,
	payload Payload,
) map[string]common.Address {
	addressMap := make(map[string]common.Address)

	// Map sender addresses
	for i, wallet := range senders {
		key := fmt.Sprintf("$sender[%d]", i+1)
		addressMap[key] = wallet.GetAddress()
	}

	// Count deployer transactions to determine contract addresses
	contractIdx := 0
	for _, tx := range payload.Txs {
		if tx.From == "deployer" && tx.To == "" {
			// This is a contract creation
			contractIdx++
			contractAddr := crypto.CreateAddress(deployer.GetAddress(), deployerNonce+uint64(contractIdx)-1)
			key := fmt.Sprintf("$contract[%d]", contractIdx)
			addressMap[key] = contractAddr
		}
	}

	return addressMap
}

func (s *Scenario) replacePlaceholders(input string, addressMap map[string]common.Address) (string, error) {
	if input == "" {
		return input, nil
	}

	result := input

	// Replace all placeholders like $contract[1], $sender[1]
	placeholderRegex := regexp.MustCompile(`\$(?:contract|sender)\[\d+\]`)
	result = placeholderRegex.ReplaceAllStringFunc(result, func(placeholder string) string {
		if addr, ok := addressMap[placeholder]; ok {
			// Return address without 0x prefix for bytecode embedding
			return strings.ToLower(addr.Hex()[2:])
		}
		return placeholder
	})

	return result, nil
}

func (s *Scenario) performPostChecks(
	ctx context.Context,
	logger *logrus.Entry,
	client *spamoor.Client,
	senderWallets []*spamoor.Wallet,
	payload Payload,
	addressMap map[string]common.Address,
	initialBalances map[int]*big.Int,
	gasCosts map[int]*senderGasCosts,
) []string {
	var failures []string

	for key, check := range payload.PostCheck {
		// Determine the address for this check
		var addr common.Address
		var isSender bool
		var senderIdx int

		if strings.HasPrefix(key, "contract[") {
			placeholder := "$" + key
			if a, ok := addressMap[placeholder]; ok {
				addr = a
			} else {
				failures = append(failures, fmt.Sprintf("%s: unknown contract", key))
				continue
			}
		} else if strings.HasPrefix(key, "sender[") {
			placeholder := "$" + key
			if a, ok := addressMap[placeholder]; ok {
				addr = a
				isSender = true
				senderIdx = parseSenderIndex(key)
			} else {
				failures = append(failures, fmt.Sprintf("%s: unknown sender", key))
				continue
			}
		} else {
			failures = append(failures, fmt.Sprintf("%s: unknown address type", key))
			continue
		}

		// Check storage slots
		for slot, expectedValue := range check.Storage {
			slotHash := common.HexToHash(slot)

			// Replace placeholders in expected value
			expectedValueReplaced, err := s.replacePlaceholders(expectedValue, addressMap)
			if err != nil {
				failures = append(failures, fmt.Sprintf("%s storage[%s]: failed to replace placeholders: %v", key, slot, err))
				continue
			}

			actualValue, err := client.GetEthClient().StorageAt(ctx, addr, slotHash, nil)
			if err != nil {
				failures = append(failures, fmt.Sprintf("%s storage[%s]: failed to get storage: %v", key, slot, err))
				continue
			}

			// Normalize values for comparison
			expectedBytes := common.FromHex(expectedValueReplaced)
			expectedHash := common.BytesToHash(expectedBytes)
			actualHash := common.BytesToHash(actualValue)

			if actualHash != expectedHash {
				failures = append(failures, fmt.Sprintf("%s storage[%s]: expected %s, got %s",
					key, slot, expectedHash.Hex(), actualHash.Hex()))
			}
		}

		// Check balance (relative comparison for senders)
		if check.Balance != "" {
			expectedBalance, ok := new(big.Int).SetString(strings.TrimPrefix(check.Balance, "0x"), 16)
			if !ok {
				failures = append(failures, fmt.Sprintf("%s balance: failed to parse expected balance: %s", key, check.Balance))
				continue
			}

			if isSender {
				// For senders, check relative balance change
				// If test expected balance to decrease by X, our wallet should also decrease by X
				if senderIdx < 1 || senderIdx > len(senderWallets) {
					failures = append(failures, fmt.Sprintf("%s balance: invalid sender index %d", key, senderIdx))
					continue
				}

				actualBalance := senderWallets[senderIdx-1].GetBalance()

				initialBalance, ok := initialBalances[senderIdx]
				if !ok {
					failures = append(failures, fmt.Sprintf("%s balance: no initial balance recorded", key))
					continue
				}

				// Get the prerequisite balance (what the test started with)
				prereqKey := fmt.Sprintf("sender[%d]", senderIdx)
				prereqBalanceStr, ok := payload.Prerequisites[prereqKey]
				if !ok {
					// No prerequisite, skip balance check
					continue
				}

				prereqBalance, ok := new(big.Int).SetString(strings.TrimPrefix(prereqBalanceStr, "0x"), 16)
				if !ok {
					failures = append(failures, fmt.Sprintf("%s balance: failed to parse prerequisite balance", key))
					continue
				}

				// Expected change in fixture = prereqBalance - expectedBalance
				// This includes: value transfers + gas at fixture gas cost
				expectedChange := new(big.Int).Sub(prereqBalance, expectedBalance)

				// Actual change = initialBalance - actualBalance
				// This includes: value transfers + gas cost at actual base fee
				actualChange := new(big.Int).Sub(initialBalance, actualBalance)

				// Adjust expected change to account for different gas costs:
				// expectedChange includes fixtureGasCost, but we paid actualGasCost
				// So: adjustedExpectedChange = expectedChange - fixtureGasCost + actualGasCost
				adjustedExpectedChange := new(big.Int).Set(expectedChange)
				if costs, ok := gasCosts[senderIdx]; ok && costs != nil {
					// Subtract fixture gas cost and add actual gas cost
					adjustedExpectedChange.Sub(adjustedExpectedChange, costs.fixtureGasCost)
					adjustedExpectedChange.Add(adjustedExpectedChange, costs.actualGasCost)
				}

				// Compare actual vs adjusted expected
				diff := new(big.Int).Sub(actualChange, adjustedExpectedChange)
				diff.Abs(diff)

				if diff.Cmp(big.NewInt(0)) != 0 {
					failures = append(failures, fmt.Sprintf(
						"%s balance: expected change %s (adjusted from %s), got %s (diff: %s)",
						key, adjustedExpectedChange.String(), expectedChange.String(),
						actualChange.String(), diff.String()))
				}
			} else {
				// For contracts, check absolute balance
				actualBalance, err := client.GetEthClient().BalanceAt(ctx, addr, nil)
				if err != nil {
					failures = append(failures, fmt.Sprintf("%s balance: failed to get balance: %v", key, err))
					continue
				}
				if actualBalance.Cmp(expectedBalance) != 0 {
					failures = append(failures, fmt.Sprintf("%s balance: expected %s, got %s",
						key, expectedBalance.String(), actualBalance.String()))
				}
			}
		}
	}

	return failures
}

// Wallet locking

func (s *Scenario) acquireWallets(ctx context.Context, count int) ([]*spamoor.Wallet, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		wallets := s.tryAcquireWallets(count)
		if wallets != nil {
			return wallets, nil
		}

		// Wait for a wallet to be released or timeout
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-s.walletUnlockCh:
			// A wallet was released, try again
		case <-time.After(100 * time.Millisecond):
			// Periodic retry
		}
	}
}

func (s *Scenario) tryAcquireWallets(count int) []*spamoor.Wallet {
	s.walletLock.Lock()
	defer s.walletLock.Unlock()

	// Find 'count' unlocked wallets
	var candidates []*spamoor.Wallet
	walletCount := s.walletPool.GetConfiguredWalletCount()

	for i := uint64(0); i < walletCount && len(candidates) < count; i++ {
		wallet := s.walletPool.GetWallet(spamoor.SelectWalletRoundRobin, int(i))
		if wallet == nil {
			continue
		}

		addr := wallet.GetAddress()
		if !s.lockedWallets[addr] {
			candidates = append(candidates, wallet)
		}
	}

	// Only lock if we found enough wallets (atomic acquisition)
	if len(candidates) < count {
		return nil
	}

	// Lock all candidates
	for _, wallet := range candidates {
		s.lockedWallets[wallet.GetAddress()] = true
	}

	return candidates
}

func (s *Scenario) releaseWallets(wallets []*spamoor.Wallet) {
	s.walletLock.Lock()
	defer s.walletLock.Unlock()

	for _, wallet := range wallets {
		addr := wallet.GetAddress()
		delete(s.lockedWallets, addr)

		// Notify that a wallet was released
		select {
		case s.walletUnlockCh <- addr:
		default:
			// Channel full, that's ok
		}
	}
}

func (s *Scenario) processWalletUnlocks(ctx context.Context) {
	// Just drain the channel to prevent blocking
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.walletUnlockCh:
			// Drained
		}
	}
}

func (s *Scenario) printStatsSummary(ctx context.Context) {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.statsLock.Lock()
			totalFailed := s.failInvalidCount + s.failRevertedCount + s.failPostcheckCount + s.failTimeoutCount
			s.logger.Infof("progress: current=%d/%d, success=%d, failed=%d (invalid=%d, reverted=%d, postcheck=%d, timeout=%d)",
				s.currentPayloadIdx, len(s.payloads), s.successCount, totalFailed,
				s.failInvalidCount, s.failRevertedCount, s.failPostcheckCount, s.failTimeoutCount)
			s.statsLock.Unlock()
		}
	}
}

// Payload loading

func (s *Scenario) loadPayloads(pathOrURL string) ([]Payload, error) {
	var data []byte
	var err error

	if strings.HasPrefix(pathOrURL, "http://") || strings.HasPrefix(pathOrURL, "https://") {
		// Load from URL
		resp, err := http.Get(pathOrURL)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch URL: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("HTTP error: %s", resp.Status)
		}

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		}
	} else {
		// Load from file
		data, err = os.ReadFile(pathOrURL)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}
	}

	var output PayloadFile
	if err := yaml.Unmarshal(data, &output); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return output.Payloads, nil
}

// loadFromFixturesRelease loads fixtures from a .tar.gz file (URL or local path),
// extracts it to a temp directory, converts the fixtures using eestconv,
// and returns the payloads.
func (s *Scenario) loadFromFixturesRelease(ctx context.Context, pathOrURL string, testPattern string, excludePattern string) ([]Payload, error) {
	s.logger.WithField("source", pathOrURL).Info("loading fixtures release")

	// Create temp directory for extraction
	tempDir, err := os.MkdirTemp("", "eest-fixtures-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer func() {
		s.logger.WithField("path", tempDir).Debug("cleaning up temp directory")
		os.RemoveAll(tempDir)
	}()

	// Check context before starting
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// Get tar.gz reader (from URL or file)
	var tarReader io.Reader
	var contentLength int64
	var cleanup func()

	if strings.HasPrefix(pathOrURL, "http://") || strings.HasPrefix(pathOrURL, "https://") {
		s.logger.Info("downloading fixtures from URL")

		// Create request with context
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, pathOrURL, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch URL: %w", err)
		}
		cleanup = func() { resp.Body.Close() }

		if resp.StatusCode != http.StatusOK {
			cleanup()
			return nil, fmt.Errorf("HTTP error: %s", resp.Status)
		}
		contentLength = resp.ContentLength
		tarReader = resp.Body

		// Log total size if known
		if contentLength > 0 {
			sizeMB := float64(contentLength) / (1024 * 1024)
			s.logger.Infof("downloading %.1f MB", sizeMB)
		}
	} else {
		s.logger.Info("loading fixtures from file")
		file, err := os.Open(pathOrURL)
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}
		cleanup = func() { file.Close() }

		// Get file size
		fileInfo, err := file.Stat()
		if err == nil {
			contentLength = fileInfo.Size()
			sizeMB := float64(contentLength) / (1024 * 1024)
			s.logger.Infof("reading %.1f MB", sizeMB)
		}
		tarReader = file
	}
	defer cleanup()

	// Wrap reader with progress tracking and context awareness
	progressReader := newProgressReader(ctx, tarReader, contentLength, s.logger, "downloading/reading")

	// Extract tar.gz with progress reporting
	s.logger.Info("extracting fixtures archive")
	fileCount, err := extractTarGzWithProgress(ctx, progressReader, tempDir, s.logger)
	if err != nil {
		return nil, fmt.Errorf("failed to extract tar.gz: %w", err)
	}
	s.logger.WithField("files", fileCount).Info("extraction complete")

	// Find the blockchain_tests directory
	blockchainTestsDir := filepath.Join(tempDir, "fixtures", "blockchain_tests")
	if _, err := os.Stat(blockchainTestsDir); os.IsNotExist(err) {
		// Try without "fixtures" subdirectory
		blockchainTestsDir = filepath.Join(tempDir, "blockchain_tests")
		if _, err := os.Stat(blockchainTestsDir); os.IsNotExist(err) {
			return nil, fmt.Errorf("could not find blockchain_tests directory in archive")
		}
	}

	s.logger.WithField("path", blockchainTestsDir).Info("found blockchain_tests directory")

	// Convert fixtures using eestconv
	converter, err := eestconv.NewConverter(s.logger, eestconv.ConvertOptions{
		TestPattern:    testPattern,
		ExcludePattern: excludePattern,
		Verbose:        false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create converter: %w", err)
	}

	output, err := converter.ConvertDirectory(blockchainTestsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to convert fixtures: %w", err)
	}

	// Convert eestconv payloads to scenario Payloads
	payloads := make([]Payload, len(output.Payloads))
	for i, p := range output.Payloads {
		payloads[i] = convertEestconvPayload(p)
	}

	return payloads, nil
}

// progressReader wraps an io.Reader and reports progress periodically
type progressReader struct {
	ctx          context.Context
	reader       io.Reader
	totalSize    int64
	bytesRead    int64
	logger       *logrus.Entry
	operation    string
	lastReport   time.Time
	reportPeriod time.Duration
	startTime    time.Time
}

func newProgressReader(ctx context.Context, reader io.Reader, totalSize int64, logger *logrus.Entry, operation string) *progressReader {
	return &progressReader{
		ctx:          ctx,
		reader:       reader,
		totalSize:    totalSize,
		logger:       logger,
		operation:    operation,
		lastReport:   time.Now(),
		reportPeriod: 30 * time.Second,
		startTime:    time.Now(),
	}
}

func (pr *progressReader) Read(p []byte) (int, error) {
	// Check context before reading
	if err := pr.ctx.Err(); err != nil {
		return 0, err
	}

	n, err := pr.reader.Read(p)
	pr.bytesRead += int64(n)

	// Report progress every 30 seconds
	if time.Since(pr.lastReport) >= pr.reportPeriod {
		pr.reportProgress()
		pr.lastReport = time.Now()
	}

	return n, err
}

func (pr *progressReader) reportProgress() {
	elapsed := time.Since(pr.startTime)
	bytesReadMB := float64(pr.bytesRead) / (1024 * 1024)

	if pr.totalSize > 0 {
		percent := float64(pr.bytesRead) / float64(pr.totalSize) * 100
		totalMB := float64(pr.totalSize) / (1024 * 1024)
		pr.logger.Infof("%s progress: %.1f MB / %.1f MB (%.1f%%) - elapsed: %s",
			pr.operation, bytesReadMB, totalMB, percent, elapsed.Round(time.Second))
	} else {
		pr.logger.Infof("%s progress: %.1f MB - elapsed: %s",
			pr.operation, bytesReadMB, elapsed.Round(time.Second))
	}
}

// extractTarGzWithProgress extracts a tar.gz archive to a destination directory,
// only extracting files under fixtures/blockchain_tests/ or blockchain_tests/,
// and reports progress periodically. It respects context cancellation.
func extractTarGzWithProgress(ctx context.Context, reader io.Reader, destDir string, logger *logrus.Entry) (int, error) {
	gzReader, err := gzip.NewReader(reader)
	if err != nil {
		return 0, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	fileCount := 0
	lastReport := time.Now()
	startTime := time.Now()
	reportPeriod := 30 * time.Second

	for {
		// Check context before processing each file
		if err := ctx.Err(); err != nil {
			return fileCount, err
		}

		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fileCount, fmt.Errorf("failed to read tar header: %w", err)
		}

		// Only extract files under blockchain_tests directories
		name := header.Name
		if !strings.Contains(name, "blockchain_tests/") &&
			!strings.HasSuffix(name, "blockchain_tests") {
			continue
		}

		// Security: prevent path traversal
		targetPath := filepath.Join(destDir, header.Name)
		if !strings.HasPrefix(targetPath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fileCount, fmt.Errorf("invalid file path in archive: %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return fileCount, fmt.Errorf("failed to create directory: %w", err)
			}
		case tar.TypeReg:
			// Create parent directory if needed
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return fileCount, fmt.Errorf("failed to create parent directory: %w", err)
			}

			outFile, err := os.Create(targetPath)
			if err != nil {
				return fileCount, fmt.Errorf("failed to create file: %w", err)
			}

			// Limit copy size for security (4GB max per file)
			if _, err := io.CopyN(outFile, tarReader, 4*1024*1024*1024); err != nil && err != io.EOF {
				outFile.Close()
				return fileCount, fmt.Errorf("failed to write file: %w", err)
			}
			outFile.Close()
			fileCount++

			// Report extraction progress every 30 seconds
			if time.Since(lastReport) >= reportPeriod {
				elapsed := time.Since(startTime)
				logger.Infof("extraction progress: %d files extracted - elapsed: %s",
					fileCount, elapsed.Round(time.Second))
				lastReport = time.Now()
			}
		}
	}

	return fileCount, nil
}

// convertEestconvPayload converts an eestconv.ConvertedPayload to the scenario's Payload type
func convertEestconvPayload(p eestconv.ConvertedPayload) Payload {
	txs := make([]Tx, len(p.Txs))
	for i, t := range p.Txs {
		accessList := make([]AccessListItem, len(t.AccessList))
		for j, al := range t.AccessList {
			accessList[j] = AccessListItem{
				Address:     al.Address,
				StorageKeys: al.StorageKeys,
			}
		}

		authList := make([]AuthorizationItem, len(t.AuthorizationList))
		for j, auth := range t.AuthorizationList {
			authList[j] = AuthorizationItem{
				ChainID: auth.ChainID,
				Address: auth.Address,
				Nonce:   auth.Nonce,
				Signer:  auth.Signer,
				V:       auth.V,
				R:       auth.R,
				S:       auth.S,
			}
		}

		txs[i] = Tx{
			From:                 t.From,
			Type:                 t.Type,
			To:                   t.To,
			Data:                 t.Data,
			Gas:                  t.Gas,
			GasPrice:             t.GasPrice,
			MaxFeePerGas:         t.MaxFeePerGas,
			MaxPriorityFeePerGas: t.MaxPriorityFeePerGas,
			Value:                t.Value,
			AccessList:           accessList,
			BlobCount:            t.BlobCount,
			AuthorizationList:    authList,
			FixtureBaseFee:       t.FixtureBaseFee,
		}
	}

	postCheck := make(map[string]PostCheckEntry, len(p.PostCheck))
	for k, v := range p.PostCheck {
		postCheck[k] = PostCheckEntry{
			Balance: v.Balance,
			Storage: v.Storage,
		}
	}

	return Payload{
		Name:          p.Name,
		Prerequisites: p.Prerequisites,
		Txs:           txs,
		PostCheck:     postCheck,
	}
}

// Helper functions

func parseSenderIndex(s string) int {
	// Parse "sender[1]" -> 1
	re := regexp.MustCompile(`sender\[(\d+)\]`)
	matches := re.FindStringSubmatch(s)
	if len(matches) < 2 {
		return 0
	}
	var idx int
	fmt.Sscanf(matches[1], "%d", &idx)
	return idx
}

func parseHexUint8(s string) uint8 {
	s = strings.TrimPrefix(s, "0x")
	if s == "" {
		return 0
	}
	val, err := hex.DecodeString(padHexString(s))
	if err != nil {
		return 0
	}
	if len(val) == 0 {
		return 0
	}
	return val[len(val)-1]
}

func parseHexUint256(s string) *uint256.Int {
	s = strings.TrimPrefix(s, "0x")
	if s == "" {
		return uint256.NewInt(0)
	}
	val := new(uint256.Int)
	val.SetFromHex("0x" + s)
	return val
}

func padHexString(s string) string {
	if len(s)%2 != 0 {
		return "0" + s
	}
	return s
}
