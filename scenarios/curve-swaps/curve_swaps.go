package curveswaps

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/utils"
)

type ScenarioOptions struct {
	TotalCount        uint64  `yaml:"total_count"`
	Throughput        uint64  `yaml:"throughput"`
	MaxPending        uint64  `yaml:"max_pending"`
	MaxWallets        uint64  `yaml:"max_wallets"`
	Rebroadcast       uint64  `yaml:"rebroadcast"`
	BaseFee           float64 `yaml:"base_fee"`
	TipFee            float64 `yaml:"tip_fee"`
	BaseFeeWei        string  `yaml:"base_fee_wei"`
	TipFeeWei         string  `yaml:"tip_fee_wei"`
	PoolCount         uint64  `yaml:"pool_count"`
	Amplification     uint64  `yaml:"amplification"`
	Fee               uint64  `yaml:"fee"`
	SeedAmount        string  `yaml:"seed_amount"`
	WalletFunding     string  `yaml:"wallet_funding"`
	MinSwapAmount     string  `yaml:"min_swap_amount"`
	MaxSwapAmount     string  `yaml:"max_swap_amount"`
	Slippage          uint64  `yaml:"slippage"`
	Timeout           string  `yaml:"timeout"`
	ClientGroup       string  `yaml:"client_group"`
	DeployClientGroup string  `yaml:"deploy_client_group"`
	LogTxs            bool    `yaml:"log_txs"`
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	curve      *Curve
	deployment *CurveDeploymentInfo
}

// swapGasLimit is the static gas limit used for all swap (spam) transactions.
// A StableSwap exchange runs the full Newton's-method get_D/get_y iterations
// plus two ERC20 transfers; this keeps comfortable headroom for that work and
// for fresh state creation under the Amsterdam fee schedule, while avoiding a
// per-tx eth_estimateGas round trip on the hot path.
const swapGasLimit = 600000

var ScenarioName = "curve-swaps"
var ScenarioDefaultOptions = ScenarioOptions{
	TotalCount:        0,
	Throughput:        10,
	MaxPending:        0,
	MaxWallets:        0,
	Rebroadcast:       1,
	BaseFee:           20,
	TipFee:            2,
	PoolCount:         1,
	Amplification:     200,
	Fee:               4000000,                     // 0.04% (denominated in 1e10)
	SeedAmount:        "1000000000000000000000000", // 1,000,000 tokens per coin
	WalletFunding:     "10000000000000000000000",   // 10,000 tokens (coin 0) per wallet
	MinSwapAmount:     "1000000000000000000",       // 1 token
	MaxSwapAmount:     "1000000000000000000000",    // 1,000 tokens
	Slippage:          100,                         // 1%
	Timeout:           "",
	ClientGroup:       "",
	DeployClientGroup: "",
	LogTxs:            false,
}
var ScenarioDescriptor = scenario.Descriptor{
	Name:           ScenarioName,
	Description:    "Send Curve StableSwap exchanges across self-deployed 3-coin stable pools",
	DefaultOptions: ScenarioDefaultOptions,
	NewScenario:    newScenario,
}

func newScenario(logger logrus.FieldLogger) scenario.Scenario {
	return &Scenario{
		options: ScenarioDefaultOptions,
		logger:  logger.WithField("scenario", ScenarioName),
	}
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
	flags.Uint64VarP(&s.options.TotalCount, "count", "c", ScenarioDefaultOptions.TotalCount, "Total number of swap transactions to send")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t", ScenarioDefaultOptions.Throughput, "Number of swap transactions to send per slot")
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", ScenarioDefaultOptions.Rebroadcast, "Enable reliable rebroadcast system")
	flags.Float64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use in swap transactions (in gwei)")
	flags.Float64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use in swap transactions (in gwei)")
	flags.StringVar(&s.options.BaseFeeWei, "basefee-wei", "", "Max fee per gas in wei (overrides --basefee for L2 sub-gwei fees)")
	flags.StringVar(&s.options.TipFeeWei, "tipfee-wei", "", "Max tip per gas in wei (overrides --tipfee for L2 sub-gwei fees)")
	flags.Uint64Var(&s.options.PoolCount, "pool-count", ScenarioDefaultOptions.PoolCount, "Number of StableSwap pools to deploy")
	flags.Uint64Var(&s.options.Amplification, "amplification", ScenarioDefaultOptions.Amplification, "StableSwap amplification coefficient (A)")
	flags.Uint64Var(&s.options.Fee, "fee", ScenarioDefaultOptions.Fee, "StableSwap swap fee, denominated in 1e10 (e.g. 4000000 = 0.04%)")
	flags.StringVar(&s.options.SeedAmount, "seed-amount", ScenarioDefaultOptions.SeedAmount, "Liquidity seeded per coin into each pool (in wei)")
	flags.StringVar(&s.options.WalletFunding, "wallet-funding", ScenarioDefaultOptions.WalletFunding, "Initial coin balance minted to each wallet (in wei)")
	flags.StringVar(&s.options.MinSwapAmount, "min-swap", ScenarioDefaultOptions.MinSwapAmount, "Minimum swap amount in wei")
	flags.StringVar(&s.options.MaxSwapAmount, "max-swap", ScenarioDefaultOptions.MaxSwapAmount, "Maximum swap amount in wei")
	flags.Uint64Var(&s.options.Slippage, "slippage", ScenarioDefaultOptions.Slippage, "Slippage tolerance in basis points")
	flags.StringVar(&s.options.Timeout, "timeout", ScenarioDefaultOptions.Timeout, "Timeout for the scenario (e.g. '1h', '30m', '5s') - empty means no timeout")
	flags.StringVar(&s.options.ClientGroup, "client-group", ScenarioDefaultOptions.ClientGroup, "Client group to use for sending transactions")
	flags.StringVar(&s.options.DeployClientGroup, "deploy-client-group", ScenarioDefaultOptions.DeployClientGroup, "Client group to use for deployments")
	flags.BoolVar(&s.options.LogTxs, "log-txs", ScenarioDefaultOptions.LogTxs, "Log all submitted transactions")
	return nil
}

func (s *Scenario) Init(options *scenario.Options) error {
	s.walletPool = options.WalletPool

	if options.Config != "" {
		// Use the generalized config validation and parsing helper
		err := scenario.ParseAndValidateConfig(&ScenarioDescriptor, options.Config, &s.options, s.logger)
		if err != nil {
			return err
		}
	}

	if s.options.PoolCount == 0 {
		return fmt.Errorf("pool-count must be at least 1")
	}

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
		if s.options.Throughput*10 < 1000 {
			s.walletPool.SetWalletCount(s.options.Throughput * 10)
		} else {
			s.walletPool.SetWalletCount(1000)
		}
	}

	if s.options.TotalCount == 0 && s.options.Throughput == 0 {
		return fmt.Errorf("neither total count nor throughput limit set, must define at least one of them (see --help for list of all flags)")
	}

	// register well known wallets
	s.walletPool.AddWellKnownWallet(&spamoor.WellKnownWalletConfig{
		Name:          "deployer",
		RefillAmount:  uint256.NewInt(2000000000000000000), // 2 ETH
		RefillBalance: uint256.NewInt(1000000000000000000), // 1 ETH
	})

	return nil
}

func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Infof("starting scenario: %s", ScenarioName)
	defer s.logger.Infof("scenario %s finished.", ScenarioName)

	seedAmount, overflow := uint256.FromBig(mustParseWei(s.options.SeedAmount))
	if overflow || seedAmount.IsZero() {
		return fmt.Errorf("invalid seed-amount: %s", s.options.SeedAmount)
	}
	walletFunding, overflow := uint256.FromBig(mustParseWei(s.options.WalletFunding))
	if overflow || walletFunding.IsZero() {
		return fmt.Errorf("invalid wallet-funding: %s", s.options.WalletFunding)
	}

	deployClientGroup := s.options.DeployClientGroup
	if deployClientGroup == "" {
		deployClientGroup = s.options.ClientGroup
	}

	s.curve = NewCurve(ctx, s.walletPool, s.logger, CurveOptions{
		BaseFee:       s.options.BaseFee,
		TipFee:        s.options.TipFee,
		BaseFeeWei:    s.options.BaseFeeWei,
		TipFeeWei:     s.options.TipFeeWei,
		PoolCount:     s.options.PoolCount,
		Amplification: s.options.Amplification,
		Fee:           s.options.Fee,
		SeedAmount:    seedAmount,
		WalletFunding: walletFunding,
		ClientGroup:   deployClientGroup,
	})

	deploymentInfo, err := s.curve.DeployCurvePools()
	if err != nil {
		s.logger.Errorf("could not deploy curve pools: %v", err)
		return err
	}
	if deploymentInfo == nil {
		return fmt.Errorf("could not deploy curve pools")
	}
	s.deployment = deploymentInfo

	if err := s.curve.InitializeContracts(deploymentInfo); err != nil {
		s.logger.Errorf("could not initialize curve contracts: %v", err)
		return err
	}

	if err := s.curve.FundAndApproveWallets(); err != nil {
		s.logger.Errorf("could not fund and approve wallets: %v", err)
		return err
	}

	s.curve.InitializeTokenBalances()

	// send transactions
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
			return fmt.Errorf("invalid timeout value: %v", err)
		}
		s.logger.Infof("Timeout set to %v", timeout)
	}

	err = scenario.RunTransactionScenario(ctx, scenario.TransactionScenarioOptions{
		TotalCount:                  s.options.TotalCount,
		Throughput:                  s.options.Throughput,
		MaxPending:                  maxPending,
		ThroughputIncrementInterval: 0,
		Timeout:                     timeout,
		WalletPool:                  s.walletPool,

		Logger: s.logger,
		ProcessNextTxFn: func(ctx context.Context, params *scenario.ProcessNextTxParams) error {
			logger := s.logger
			receiptChan, tx, client, wallet, err := s.sendTx(ctx, params.TxIdx)
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
				if err != nil {
					logger.Warnf("could not send transaction: %v", err)
				} else if s.options.LogTxs {
					logger.Infof("sent tx #%6d: %v", params.TxIdx+1, tx.Hash().String())
				} else {
					logger.Debugf("sent tx #%6d: %v", params.TxIdx+1, tx.Hash().String())
				}
			})

			// wait for receipt
			if _, err := receiptChan.Wait(ctx); err != nil {
				return err
			}

			return err
		},
	})

	return err
}

func (s *Scenario) sendTx(ctx context.Context, txIdx uint64) (scenario.ReceiptChan, *types.Transaction, *spamoor.Client, *spamoor.Wallet, error) {
	client := s.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, int(txIdx)),
		spamoor.WithClientGroup(s.options.ClientGroup),
	)
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(txIdx))

	if client == nil {
		return nil, nil, client, wallet, scenario.ErrNoClients
	}
	if wallet == nil {
		return nil, nil, client, wallet, scenario.ErrNoWallet
	}

	if err := wallet.ResetNoncesIfNeeded(ctx, client); err != nil {
		return nil, nil, client, wallet, err
	}

	baseFeeWei, tipFeeWei := spamoor.ResolveFees(s.options.BaseFee, s.options.TipFee, s.options.BaseFeeWei, s.options.TipFeeWei)
	feeCap, tipCap, err := s.walletPool.GetSuggestedFees(client, baseFeeWei, tipFeeWei)
	if err != nil {
		return nil, nil, client, wallet, err
	}

	tx, err := s.buildSwapTx(ctx, wallet, feeCap, tipCap)
	if err != nil {
		return nil, nil, client, wallet, err
	}

	return s.submitSwapTx(ctx, txIdx, client, wallet, tx)
}

// submitSwapTx sends a built swap transaction and returns a receipt channel.
func (s *Scenario) submitSwapTx(ctx context.Context, txIdx uint64, client *spamoor.Client, wallet *spamoor.Wallet, tx *types.Transaction) (scenario.ReceiptChan, *types.Transaction, *spamoor.Client, *spamoor.Wallet, error) {
	receiptChan := make(scenario.ReceiptChan, 1)
	err := s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		ClientGroup: s.options.ClientGroup,
		Rebroadcast: s.options.Rebroadcast > 0,
		OnComplete: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			receiptChan <- receipt
		},
		OnConfirm: func(tx *types.Transaction, receipt *types.Receipt) {
			txFees := utils.GetTransactionFees(tx, receipt)
			s.logger.WithField("rpc", client.GetName()).Debugf(
				" transaction %d confirmed in block #%v. total fee: %v gwei (base: %v) logs: %v",
				txIdx+1,
				receipt.BlockNumber.String(),
				txFees.TotalFeeGweiString(),
				txFees.TxBaseFeeGweiString(),
				len(receipt.Logs),
			)
		},
		LogFn: spamoor.GetDefaultLogFn(s.logger, "", fmt.Sprintf("%6d", txIdx+1), tx),
	})
	if err != nil {
		// mark nonce as skipped if tx was not sent
		wallet.MarkSkippedNonce(tx.Nonce())

		return nil, nil, client, wallet, err
	}

	return receiptChan, tx, client, wallet, nil
}

// mustParseWei parses a base-10 wei string, returning 0 on failure so the caller
// can surface a validation error.
func mustParseWei(s string) *big.Int {
	v, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return big.NewInt(0)
	}
	return v
}
