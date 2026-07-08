package aavelending

import (
	"context"
	"fmt"
	"math/big"
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
	MinAmount         string  `yaml:"min_amount"`
	MaxAmount         string  `yaml:"max_amount"`
	SeedAmount        string  `yaml:"seed_amount"`
	WalletFunding     string  `yaml:"wallet_funding"`
	RiskyRatio        uint64  `yaml:"risky_ratio"`
	Liquidations      bool    `yaml:"liquidations"`
	PriceTickInterval uint64  `yaml:"price_tick_interval"`
	PriceVolatility   uint64  `yaml:"price_volatility"`
	GasLimit          uint64  `yaml:"gas_limit"`
	Timeout           string  `yaml:"timeout"`
	ClientGroup       string  `yaml:"client_group"`
	DeployClientGroup string  `yaml:"deploy_client_group"`
	LogTxs            bool    `yaml:"log_txs"`
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	deployment *DeploymentInfo

	minAmount     *big.Int
	maxAmount     *big.Int
	seedAmount    *big.Int
	walletFunding *big.Int
	walletCount   uint64

	// prices mirrors the oracle answer the scenario maintains per reserve (8
	// decimals). The price-tick action random-walks these and the action engine
	// reads them to size borrows and find liquidatable positions.
	prices   [2]*big.Int
	pricesMu sync.RWMutex

	// liqCursor round-robins liquidation scans across the risky wallets.
	liqCursor uint64
}

// aaveActionGasLimit is the static gas limit for the supply/borrow/repay/
// withdraw/liquidate transactions on the hot path. The actions skip per-tx gas
// estimation to avoid the extra RPC round trip; this limit leaves headroom for
// the heaviest action (a borrow or liquidation that writes fresh debt/collateral
// slots and updates indexes) under the Amsterdam state-creation fee schedule.
// Override with --gas-limit.
const aaveActionGasLimit = 1200000

var ScenarioName = "aave-lending"
var ScenarioDefaultOptions = ScenarioOptions{
	TotalCount:        0,
	Throughput:        10,
	MaxPending:        0,
	MaxWallets:        0,
	Rebroadcast:       1,
	BaseFee:           20,
	TipFee:            2,
	BaseFeeWei:        "",
	TipFeeWei:         "",
	MinAmount:         "1000000000000000000",       // 1 token
	MaxAmount:         "2000000000000000000000",    // 2000 tokens
	SeedAmount:        "1000000000000000000000000", // 1,000,000 tokens of initial borrowable liquidity per reserve
	WalletFunding:     "100000000000000000000000",  // 100,000 tokens minted to each wallet
	RiskyRatio:        6,                           // 1 in 6 wallets runs a near-max-LTV position (liquidation targets); 0 disables
	Liquidations:      true,
	PriceTickInterval: 40,   // random-walk an oracle price every N transactions; 0 disables
	PriceVolatility:   2000, // +/- 20% price band around $1 (deep enough for occasional liquidations)
	GasLimit:          0,
	Timeout:           "",
	ClientGroup:       "",
	DeployClientGroup: "",
	LogTxs:            false,
}
var ScenarioDescriptor = scenario.Descriptor{
	Name:           ScenarioName,
	Description:    "Deploy a full Aave V3 market and have wallets organically supply/borrow/repay/withdraw two tokens, with oracle price moves and occasional liquidations",
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
	flags.Uint64VarP(&s.options.TotalCount, "count", "c", ScenarioDefaultOptions.TotalCount, "Total number of lending transactions to send")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t", ScenarioDefaultOptions.Throughput, "Number of lending transactions to send per slot")
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", ScenarioDefaultOptions.Rebroadcast, "Enable reliable rebroadcast system")
	flags.Float64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use in lending transactions (in gwei)")
	flags.Float64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use in lending transactions (in gwei)")
	flags.StringVar(&s.options.BaseFeeWei, "basefee-wei", "", "Max fee per gas in wei (overrides --basefee for L2 sub-gwei fees)")
	flags.StringVar(&s.options.TipFeeWei, "tipfee-wei", "", "Max tip per gas in wei (overrides --tipfee for L2 sub-gwei fees)")
	flags.StringVar(&s.options.MinAmount, "min-amount", ScenarioDefaultOptions.MinAmount, "Minimum amount per lending action (in wei)")
	flags.StringVar(&s.options.MaxAmount, "max-amount", ScenarioDefaultOptions.MaxAmount, "Maximum amount per lending action (in wei)")
	flags.StringVar(&s.options.SeedAmount, "seed-amount", ScenarioDefaultOptions.SeedAmount, "Initial borrowable liquidity supplied to each reserve by the deployer (in wei)")
	flags.StringVar(&s.options.WalletFunding, "wallet-funding", ScenarioDefaultOptions.WalletFunding, "Amount of each token minted to every child wallet (in wei)")
	flags.Uint64Var(&s.options.RiskyRatio, "risky-ratio", ScenarioDefaultOptions.RiskyRatio, "1 in N wallets runs a near-max-LTV position (liquidation targets); 0 disables risky positions")
	flags.BoolVar(&s.options.Liquidations, "liquidations", ScenarioDefaultOptions.Liquidations, "Enable opportunistic liquidations of underwater positions")
	flags.Uint64Var(&s.options.PriceTickInterval, "price-tick-interval", ScenarioDefaultOptions.PriceTickInterval, "Random-walk an oracle price every N transactions; 0 disables price moves")
	flags.Uint64Var(&s.options.PriceVolatility, "price-volatility", ScenarioDefaultOptions.PriceVolatility, "Oracle price band in basis points around the $1 base (e.g. 1500 = +/-15%)")
	flags.Uint64Var(&s.options.GasLimit, "gas-limit", ScenarioDefaultOptions.GasLimit, "Gas limit per lending transaction (0 = use built-in default)")
	flags.StringVar(&s.options.Timeout, "timeout", ScenarioDefaultOptions.Timeout, "Timeout for the scenario (e.g. '1h', '30m', '5s') - empty means no timeout")
	flags.StringVar(&s.options.ClientGroup, "client-group", ScenarioDefaultOptions.ClientGroup, "Client group to use for sending transactions")
	flags.StringVar(&s.options.DeployClientGroup, "deploy-client-group", ScenarioDefaultOptions.DeployClientGroup, "Client group to use for deployments")
	flags.BoolVar(&s.options.LogTxs, "log-txs", ScenarioDefaultOptions.LogTxs, "Log all submitted transactions")
	return nil
}

func (s *Scenario) Init(options *scenario.Options) error {
	s.walletPool = options.WalletPool

	if options.Config != "" {
		err := scenario.ParseAndValidateConfig(&ScenarioDescriptor, options.Config, &s.options, s.logger)
		if err != nil {
			return err
		}
	}

	if s.options.TotalCount == 0 && s.options.Throughput == 0 {
		return fmt.Errorf("neither total count nor throughput limit set, must define at least one of them (see --help for list of all flags)")
	}

	var ok bool
	if s.minAmount, ok = new(big.Int).SetString(s.options.MinAmount, 10); !ok || s.minAmount.Sign() <= 0 {
		return fmt.Errorf("invalid min-amount: %s", s.options.MinAmount)
	}
	if s.maxAmount, ok = new(big.Int).SetString(s.options.MaxAmount, 10); !ok || s.maxAmount.Cmp(s.minAmount) < 0 {
		return fmt.Errorf("invalid max-amount: %s (must be >= min-amount)", s.options.MaxAmount)
	}
	if s.seedAmount, ok = new(big.Int).SetString(s.options.SeedAmount, 10); !ok || s.seedAmount.Sign() <= 0 {
		return fmt.Errorf("invalid seed-amount: %s", s.options.SeedAmount)
	}
	if s.walletFunding, ok = new(big.Int).SetString(s.options.WalletFunding, 10); !ok || s.walletFunding.Sign() <= 0 {
		return fmt.Errorf("invalid wallet-funding: %s", s.options.WalletFunding)
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

	// The deployer funds the heavy Aave deployment, seeds reserve liquidity and
	// acts as the reserve treasury, so it needs a generous balance.
	s.walletPool.AddWellKnownWallet(&spamoor.WellKnownWalletConfig{
		Name:          "deployer",
		RefillAmount:  uint256.NewInt(0).Mul(uint256.NewInt(20), uint256.NewInt(1000000000000000000)), // 20 ETH
		RefillBalance: uint256.NewInt(0).Mul(uint256.NewInt(10), uint256.NewInt(1000000000000000000)), // 10 ETH
	})

	return nil
}

func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Infof("starting scenario: %s", ScenarioName)
	defer s.logger.Infof("scenario %s finished.", ScenarioName)

	deployment, err := s.DeployAaveMarket(ctx)
	if err != nil {
		s.logger.Errorf("could not deploy aave market: %v", err)
		return err
	}
	s.deployment = deployment

	// seed the price mirror with the on-chain base price for both reserves
	s.prices[0] = big.NewInt(oraclePriceAnswer)
	s.prices[1] = big.NewInt(oraclePriceAnswer)

	if err := s.FundAndApproveWallets(ctx, deployment); err != nil {
		s.logger.Errorf("could not fund wallets: %v", err)
		return err
	}

	s.walletCount = s.walletPool.GetConfiguredWalletCount()
	if s.walletCount == 0 {
		return fmt.Errorf("no child wallets configured")
	}

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

	var timeout time.Duration
	if s.options.Timeout != "" {
		timeout, err = time.ParseDuration(s.options.Timeout)
		if err != nil {
			return fmt.Errorf("invalid timeout value: %v", err)
		}
		s.logger.Infof("timeout set to %v", timeout)
	}

	return scenario.RunTransactionScenario(ctx, scenario.TransactionScenarioOptions{
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

			if receiptChan != nil {
				if _, err := receiptChan.Wait(ctx); err != nil {
					return err
				}
			}

			return err
		},
	})
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

	tx, action, err := s.buildActionTx(ctx, wallet, txIdx, feeCap, tipCap)
	if err != nil {
		return nil, nil, client, wallet, err
	}

	receiptChan := make(scenario.ReceiptChan, 1)
	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		ClientGroup: s.options.ClientGroup,
		Rebroadcast: s.options.Rebroadcast > 0,
		OnComplete: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			receiptChan <- receipt
		},
		OnConfirm: func(tx *types.Transaction, receipt *types.Receipt) {
			txFees := utils.GetTransactionFees(tx, receipt)
			outcome := "ok"
			if receipt.Status != types.ReceiptStatusSuccessful {
				outcome = "reverted"
			}
			s.logger.WithField("rpc", client.GetName()).Debugf(
				" transaction %d (%s) %s in block #%v. total fee: %v gwei",
				txIdx+1, action, outcome, receipt.BlockNumber.String(), txFees.TotalFeeGweiString(),
			)
		},
		LogFn: spamoor.GetDefaultLogFn(s.logger, "", fmt.Sprintf("%6d", txIdx+1), tx),
	})
	if err != nil {
		wallet.MarkSkippedNonce(tx.Nonce())
		return nil, nil, client, wallet, err
	}

	return receiptChan, tx, client, wallet, nil
}

// actionGasLimit returns the configured per-tx gas limit, falling back to the
// built-in default when unset.
func (s *Scenario) actionGasLimit() uint64 {
	if s.options.GasLimit > 0 {
		return s.options.GasLimit
	}
	return aaveActionGasLimit
}
