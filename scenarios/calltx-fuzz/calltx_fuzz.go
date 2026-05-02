package calltxfuzz

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/utils"
)

// ScenarioOptions configures the calltx-fuzz scenario.
type ScenarioOptions struct {
	// Standard options
	TotalCount  uint64  `yaml:"total_count"`
	Throughput  uint64  `yaml:"throughput"`
	MaxPending  uint64  `yaml:"max_pending"`
	MaxWallets  uint64  `yaml:"max_wallets"`
	Rebroadcast uint64  `yaml:"rebroadcast"`
	BaseFee     float64 `yaml:"base_fee"`
	TipFee      float64 `yaml:"tip_fee"`
	BaseFeeWei  string  `yaml:"base_fee_wei"`
	TipFeeWei   string  `yaml:"tip_fee_wei"`
	GasLimit    uint64  `yaml:"gas_limit"`
	Timeout     string  `yaml:"timeout"`
	ClientGroup string  `yaml:"client_group"`
	LogTxs      bool    `yaml:"log_txs"`

	// Bytecode generation
	MaxCodeSize uint64 `yaml:"max_code_size"`
	MinCodeSize uint64 `yaml:"min_code_size"`
	PayloadSeed string `yaml:"payload_seed"`
	TxIdOffset  uint64 `yaml:"tx_id_offset"`

	// Mode selection
	FuzzType         string  `yaml:"fuzz_type"`
	ContractPoolSize uint64  `yaml:"contract_pool_size"`
	DeployRatio      float64 `yaml:"deploy_ratio"`
	CalldataMaxSize  uint64  `yaml:"calldata_max_size"`

	// SetCode-specific
	SetCodeFuzzMode   string  `yaml:"setcode_fuzz_mode"`
	MinAuthorizations uint64  `yaml:"min_authorizations"`
	MaxAuthorizations uint64  `yaml:"max_authorizations"`
	MaxDelegators     uint64  `yaml:"max_delegators"`
	InvalidAuthRatio  float64 `yaml:"invalid_auth_ratio"`

	// Frame-specific (stub)
	FrameFuzzMode string `yaml:"frame_fuzz_mode"`
	MaxFrames     uint64 `yaml:"max_frames"`
}

// Scenario implements the calltx-fuzz transaction scenario.
type Scenario struct {
	options      ScenarioOptions
	logger       *logrus.Entry
	walletPool   *spamoor.WalletPool
	seed         string
	contractPool *ContractPool

	delegatorSeed []byte
	delegatorMu   sync.Mutex
	delegators    []*spamoor.Wallet
}

// ScenarioName is the registered name for this scenario.
var ScenarioName = "calltx-fuzz"

// ScenarioDefaultOptions contains sensible defaults for the scenario.
var ScenarioDefaultOptions = ScenarioOptions{
	TotalCount:        0,
	Throughput:        50,
	MaxPending:        100,
	MaxWallets:        0,
	Rebroadcast:       30,
	BaseFee:           20,
	TipFee:            2,
	GasLimit:          5000000,
	MaxCodeSize:       1024,
	MinCodeSize:       200,
	PayloadSeed:       "",
	TxIdOffset:        0,
	FuzzType:          "setcode",
	ContractPoolSize:  50,
	DeployRatio:       0.1,
	CalldataMaxSize:   256,
	SetCodeFuzzMode:   "mixed",
	MinAuthorizations: 1,
	MaxAuthorizations: 5,
	MaxDelegators:     100,
	InvalidAuthRatio:  0.1,
	FrameFuzzMode:     "mixed",
	MaxFrames:         10,
}

// ScenarioDescriptor is the registration entry for this scenario.
var ScenarioDescriptor = scenario.Descriptor{
	Name: ScenarioName,
	Description: "Deploy fuzzed contracts and call them via " +
		"Type 2 (call), Type 4 (setcode delegation), " +
		"or Type 6 (frame tx stub) transactions",
	DefaultOptions: ScenarioDefaultOptions,
	NewScenario:    newScenario,
}

func newScenario(logger logrus.FieldLogger) scenario.Scenario {
	return &Scenario{
		options: ScenarioDefaultOptions,
		logger:  logger.WithField("scenario", ScenarioName),
	}
}

// Flags registers CLI flags for the scenario.
func (s *Scenario) Flags(flags *pflag.FlagSet) error {
	flags.Uint64VarP(&s.options.TotalCount, "count", "c", ScenarioDefaultOptions.TotalCount, "Total number of transactions to send (0 = unlimited)")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t", ScenarioDefaultOptions.Throughput, "Number of transactions per slot")
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", ScenarioDefaultOptions.Rebroadcast, "Enable reliable rebroadcast with exponential backoff")
	flags.Float64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas in gwei")
	flags.Float64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas in gwei")
	flags.StringVar(&s.options.BaseFeeWei, "basefee-wei", "", "Max fee per gas in wei (overrides --basefee)")
	flags.StringVar(&s.options.TipFeeWei, "tipfee-wei", "", "Max tip per gas in wei (overrides --tipfee)")
	flags.Uint64Var(&s.options.GasLimit, "gaslimit", ScenarioDefaultOptions.GasLimit, "Gas limit per transaction")
	flags.StringVar(&s.options.Timeout, "timeout", ScenarioDefaultOptions.Timeout, "Scenario timeout (e.g. '1h', '30m')")
	flags.StringVar(&s.options.ClientGroup, "client-group", ScenarioDefaultOptions.ClientGroup, "Client group for sending transactions")
	flags.BoolVar(&s.options.LogTxs, "log-txs", ScenarioDefaultOptions.LogTxs, "Log all submitted transactions")

	// Bytecode generation
	flags.Uint64Var(&s.options.MaxCodeSize, "max-code-size", ScenarioDefaultOptions.MaxCodeSize, "Maximum runtime bytecode size")
	flags.Uint64Var(&s.options.MinCodeSize, "min-code-size", ScenarioDefaultOptions.MinCodeSize, "Minimum runtime bytecode size")
	flags.StringVar(&s.options.PayloadSeed, "payload-seed", ScenarioDefaultOptions.PayloadSeed, "Hex seed for reproducible fuzzing (e.g. 0x1234)")
	flags.Uint64Var(&s.options.TxIdOffset, "tx-id-offset", ScenarioDefaultOptions.TxIdOffset, "Start from a specific transaction ID")

	// Mode selection
	flags.StringVar(&s.options.FuzzType, "fuzz-type", ScenarioDefaultOptions.FuzzType, "Transaction type: 'call', 'setcode', or 'frame'")
	flags.Uint64Var(&s.options.ContractPoolSize, "contract-pool-size", ScenarioDefaultOptions.ContractPoolSize, "Number of contracts in the pool")
	flags.Float64Var(&s.options.DeployRatio, "deploy-ratio", ScenarioDefaultOptions.DeployRatio, "Fraction of txs that deploy new contracts (0.0-1.0)")
	flags.Uint64Var(&s.options.CalldataMaxSize, "calldata-max-size", ScenarioDefaultOptions.CalldataMaxSize, "Maximum calldata size in bytes")

	// SetCode-specific
	flags.StringVar(&s.options.SetCodeFuzzMode, "setcode-fuzz-mode", ScenarioDefaultOptions.SetCodeFuzzMode, "SetCode sub-mode: 'delegation', 'execution', 'storage', 'mixed'")
	flags.Uint64Var(&s.options.MinAuthorizations, "min-authorizations", ScenarioDefaultOptions.MinAuthorizations, "Minimum authorizations per setcode tx")
	flags.Uint64Var(&s.options.MaxAuthorizations, "max-authorizations", ScenarioDefaultOptions.MaxAuthorizations, "Maximum authorizations per setcode tx")
	flags.Uint64Var(&s.options.MaxDelegators, "max-delegators", ScenarioDefaultOptions.MaxDelegators, "Maximum delegator wallets to reuse (0 = unlimited)")
	flags.Float64Var(&s.options.InvalidAuthRatio, "invalid-auth-ratio", ScenarioDefaultOptions.InvalidAuthRatio, "Fraction of deliberately invalid authorizations")

	// Frame-specific
	flags.StringVar(&s.options.FrameFuzzMode, "frame-fuzz-mode", ScenarioDefaultOptions.FrameFuzzMode, "Frame sub-mode: 'approval', 'isolation', 'ordering', 'mixed'")
	flags.Uint64Var(&s.options.MaxFrames, "max-frames", ScenarioDefaultOptions.MaxFrames, "Maximum frames per frame tx")

	return nil
}

// Init validates configuration and sets up the scenario.
func (s *Scenario) Init(options *scenario.Options) error {
	s.walletPool = options.WalletPool

	if options.Config != "" {
		err := scenario.ParseAndValidateConfig(
			&ScenarioDescriptor, options.Config, &s.options, s.logger,
		)
		if err != nil {
			return err
		}
	}

	// Configure wallet count
	if s.options.MaxWallets > 0 {
		s.walletPool.SetWalletCount(s.options.MaxWallets)
	} else if s.options.TotalCount > 0 {
		maxWallets := s.options.TotalCount / 50
		if maxWallets < 10 {
			maxWallets = 10
		} else if maxWallets > 1000 {
			maxWallets = 1000
		}
		s.walletPool.SetWalletCount(maxWallets)
	} else {
		walletCount := s.options.Throughput * 10
		if walletCount > 1000 {
			walletCount = 1000
		}
		s.walletPool.SetWalletCount(walletCount)
	}

	// Configure deployer wallet
	s.walletPool.SetRefillAmount(utils.EtherToWei(uint256.NewInt(5)))
	s.walletPool.SetRefillBalance(utils.EtherToWei(uint256.NewInt(1)))
	s.walletPool.AddWellKnownWallet(&spamoor.WellKnownWalletConfig{
		Name:          "deployer",
		RefillAmount:  utils.EtherToWei(uint256.NewInt(50)),
		RefillBalance: utils.EtherToWei(uint256.NewInt(10)),
		VeryWellKnown: false,
	})

	// Validate options
	if s.options.TotalCount == 0 && s.options.Throughput == 0 {
		return fmt.Errorf("neither total count nor throughput set, must define at least one")
	}

	if s.options.MinCodeSize > s.options.MaxCodeSize {
		return fmt.Errorf("min code size cannot be larger than max code size")
	}

	if s.options.GasLimit > utils.MaxGasLimitPerTx {
		s.logger.Warnf("Gas limit %d exceeds %d and will likely be dropped",
			s.options.GasLimit, utils.MaxGasLimitPerTx)
	}

	// Validate seed
	if s.options.PayloadSeed != "" {
		if err := validateSeed(s.options.PayloadSeed); err != nil {
			return fmt.Errorf("invalid payload seed: %w", err)
		}
	}

	// Validate fuzz type
	validTypes := map[string]bool{"call": true, "setcode": true, "frame": true}
	if !validTypes[s.options.FuzzType] {
		return fmt.Errorf("invalid fuzz-type %q, must be 'call', 'setcode', or 'frame'", s.options.FuzzType)
	}

	// Validate deploy ratio
	if s.options.DeployRatio < 0 || s.options.DeployRatio > 1 {
		return fmt.Errorf("deploy-ratio must be between 0.0 and 1.0")
	}

	// Initialize delegator seed for setcode mode
	s.delegatorSeed = make([]byte, 32)
	rand.Read(s.delegatorSeed)

	if s.options.MaxDelegators > 0 {
		s.delegators = make([]*spamoor.Wallet, 0, s.options.MaxDelegators)
	}

	return nil
}

// Run executes the calltx-fuzz scenario.
func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Infof("starting scenario: %s (fuzz-type=%s)", ScenarioName, s.options.FuzzType)
	defer s.logger.Infof("scenario %s finished.", ScenarioName)

	// Generate seed
	s.seed = s.options.PayloadSeed
	if s.seed == "" {
		randomBytes := make([]byte, 32)
		rand.Read(randomBytes)
		s.seed = hex.EncodeToString(randomBytes)
		s.logger.Infof("Generated random seed: 0x%s", s.seed)
	} else {
		s.logger.Infof("Using provided seed: %s", s.seed)
	}

	// Initialize contract pool
	s.contractPool = NewContractPool(
		s.options.ContractPoolSize,
		s.seed,
		s.options.MaxCodeSize,
		s.options.MinCodeSize,
		s.options.GasLimit,
		s.options.DeployRatio,
		s.logger,
	)

	// Deploy initial pool
	err := s.contractPool.InitPool(
		ctx,
		s.walletPool,
		s.options.TxIdOffset,
		s.options.BaseFee, s.options.TipFee,
		s.options.BaseFeeWei, s.options.TipFeeWei,
		s.options.GasLimit,
		s.options.ClientGroup,
	)
	if err != nil {
		return fmt.Errorf("contract pool initialization failed: %w", err)
	}

	if s.contractPool.Size() == 0 {
		return fmt.Errorf("no contracts deployed in pool, cannot proceed")
	}

	// Frame mode warning
	if s.options.FuzzType == "frame" {
		s.logger.Warn("Frame transaction mode generates tx data only " +
			"— signing and submission require client support for EIP-8141")
	}

	// Configure max pending
	maxPending := s.options.MaxPending
	if maxPending == 0 {
		maxPending = s.options.Throughput * 10
		if maxPending == 0 {
			maxPending = 4000
		}
		if maxPending > s.walletPool.GetConfiguredWalletCount()*10 {
			maxPending = s.walletPool.GetConfiguredWalletCount() * 10
		}
	}

	// Parse timeout
	var timeout time.Duration
	if s.options.Timeout != "" {
		timeout, err = time.ParseDuration(s.options.Timeout)
		if err != nil {
			return fmt.Errorf("invalid timeout: %w", err)
		}
		s.logger.Infof("Timeout set to %v", timeout)
	}

	s.logger.WithFields(logrus.Fields{
		"total":          s.options.TotalCount,
		"throughput":     s.options.Throughput,
		"maxPending":     maxPending,
		"poolSize":       s.contractPool.Size(),
		"deployInterval": s.contractPool.GetDeployInterval(),
		"codeSize":       fmt.Sprintf("%d-%d", s.options.MinCodeSize, s.options.MaxCodeSize),
		"fuzzType":       s.options.FuzzType,
	}).Info("starting transaction fuzzing")

	err = scenario.RunTransactionScenario(ctx, scenario.TransactionScenarioOptions{
		TotalCount: s.options.TotalCount,
		Throughput: s.options.Throughput,
		MaxPending: maxPending,
		Timeout:    timeout,
		WalletPool: s.walletPool,
		Logger:     s.logger,
		ProcessNextTxFn: func(ctx context.Context, params *scenario.ProcessNextTxParams) error {
			return s.processNextTx(ctx, params)
		},
	})

	return err
}

// processNextTx dispatches to the appropriate mode handler or deploys a new contract.
func (s *Scenario) processNextTx(ctx context.Context, params *scenario.ProcessNextTxParams) error {
	logger := s.logger

	// Frame mode: generate data only, no actual transaction
	if s.options.FuzzType == "frame" {
		ftx := s.generateFrameTx(params.TxIdx)
		params.NotifySubmitted()
		params.OrderedLogCb(func() {
			logger.Debugf("frame tx #%6d generated: %d frames",
				params.TxIdx+1, len(ftx.Frames))
		})
		return nil
	}

	// Deterministic deployment check based on effective tx ID
	effectiveTxID := params.TxIdx + s.options.TxIdOffset
	shouldDeploy := s.contractPool.ShouldDeploy(effectiveTxID)

	var (
		receiptChan scenario.ReceiptChan
		tx          *types.Transaction
		client      *spamoor.Client
		wallet      *spamoor.Wallet
		err         error
	)

	if shouldDeploy {
		receiptChan, tx, client, wallet, err = s.contractPool.DeployNew(
			ctx, s.walletPool, effectiveTxID,
			s.options.BaseFee, s.options.TipFee,
			s.options.BaseFeeWei, s.options.TipFeeWei,
			s.options.GasLimit,
			s.options.ClientGroup,
			params.TxIdx,
		)
	} else {
		switch s.options.FuzzType {
		case "call":
			receiptChan, tx, client, wallet, err = s.sendCallTx(ctx, params.TxIdx)
		case "setcode":
			receiptChan, tx, client, wallet, err = s.sendSetCodeTx(ctx, params.TxIdx)
		default:
			return fmt.Errorf("unsupported fuzz type: %s", s.options.FuzzType)
		}
	}

	if client != nil {
		logger = logger.WithField("rpc", client.GetName())
	}
	if tx != nil {
		logger = logger.WithField("nonce", tx.Nonce())
	}
	if wallet != nil {
		logger = logger.WithField("wallet", s.walletPool.GetWalletName(wallet.GetAddress()))
	}

	params.NotifySubmitted()
	params.OrderedLogCb(func() {
		action := s.options.FuzzType
		if shouldDeploy {
			action = "deploy"
		}
		if err != nil {
			logger.Warnf("%s tx #%6d failed: %v", action, params.TxIdx+1, err)
		} else if s.options.LogTxs {
			logger.Infof("%s tx #%6d sent: %v", action, params.TxIdx+1, tx.Hash().String())
		} else {
			logger.Debugf("%s tx #%6d sent: %v", action, params.TxIdx+1, tx.Hash().String())
		}
	})

	// Wait for receipt
	if receiptChan != nil {
		if _, waitErr := receiptChan.Wait(ctx); waitErr != nil {
			return waitErr
		}
	}

	return err
}

// validateSeed validates a hex seed string.
func validateSeed(seed string) error {
	clean := strings.TrimPrefix(seed, "0x")
	if _, err := hex.DecodeString(clean); err != nil {
		return fmt.Errorf("seed must be valid hex: %w", err)
	}

	return nil
}
