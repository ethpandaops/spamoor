package delegatedcounter

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/txtypes"
	"github.com/ethpandaops/spamoor/utils"
)

// counterInitCode deploys `push0 sload push1 1 add push0 sstore stop`, which increments slot 0 on every call.
var counterInitCode = common.FromHex("0x675f546001015f55005f5260086018f3")

type ScenarioOptions struct {
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
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	authority common.Address
}

var ScenarioName = "delegated-counter"
var ScenarioDefaultOptions = ScenarioOptions{
	Throughput:  10,
	Rebroadcast: 1,
	BaseFee:     20,
	TipFee:      2,
	GasLimit:    100000,
}

var ScenarioDescriptor = scenario.Descriptor{
	Name:           ScenarioName,
	Description:    "Delegates an EOA to a counter contract via EIP-7702, then sends txs to the EOA so each one increments its storage",
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
	flags.Uint64VarP(&s.options.TotalCount, "count", "c", ScenarioDefaultOptions.TotalCount, "Total number of counter transactions to send")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t", ScenarioDefaultOptions.Throughput, "Number of counter transactions to send per slot")
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", ScenarioDefaultOptions.Rebroadcast, "Enable reliable rebroadcast system")
	flags.Float64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use in transactions (in gwei)")
	flags.Float64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use in transactions (in gwei)")
	flags.StringVar(&s.options.BaseFeeWei, "basefee-wei", "", "Max fee per gas in wei (overrides --basefee)")
	flags.StringVar(&s.options.TipFeeWei, "tipfee-wei", "", "Max tip per gas in wei (overrides --tipfee)")
	flags.Uint64Var(&s.options.GasLimit, "gaslimit", ScenarioDefaultOptions.GasLimit, "Gas limit for counter transactions")
	flags.StringVar(&s.options.Timeout, "timeout", ScenarioDefaultOptions.Timeout, "Timeout for the scenario (e.g. '1h', '30m', '5s') - empty means no timeout")
	flags.StringVar(&s.options.ClientGroup, "client-group", ScenarioDefaultOptions.ClientGroup, "Client group to use for sending transactions")
	flags.BoolVar(&s.options.LogTxs, "log-txs", ScenarioDefaultOptions.LogTxs, "Log all submitted transactions")
	return nil
}

func (s *Scenario) Init(options *scenario.Options) error {
	s.walletPool = options.WalletPool

	if options.Config != "" {
		if err := scenario.ParseAndValidateConfig(&ScenarioDescriptor, options.Config, &s.options, s.logger); err != nil {
			return err
		}
	}

	if s.options.TotalCount == 0 && s.options.Throughput == 0 {
		return fmt.Errorf("neither total count nor throughput limit set, must define at least one of them (see --help for list of all flags)")
	}

	walletCount := s.options.MaxWallets
	if walletCount == 0 {
		walletCount = s.options.Throughput * 10
		if walletCount < 10 {
			walletCount = 10
		} else if walletCount > 1000 {
			walletCount = 1000
		}
	}
	s.walletPool.SetWalletCount(walletCount)

	s.walletPool.AddWellKnownWallet(&spamoor.WellKnownWalletConfig{
		Name:          "deployer",
		RefillAmount:  utils.EtherToWei(uint256.NewInt(1)),
		RefillBalance: utils.EtherToWei(uint256.NewInt(1)),
	})
	// The authority only signs the delegation, so it is never funded: any ETH sent to it would run the counter.
	s.walletPool.AddWellKnownWallet(&spamoor.WellKnownWalletConfig{
		Name:          "authority",
		RefillAmount:  uint256.NewInt(0),
		RefillBalance: uint256.NewInt(0),
	})

	return nil
}

func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Infof("starting scenario: %s", ScenarioName)
	defer s.logger.Infof("scenario %s finished.", ScenarioName)

	// The setup is a receipt-dependent chain (deploy -> delegate to the deployed address), which a
	// spammer config cannot express.
	counter, err := s.deployCounter(ctx)
	if err != nil {
		return err
	}
	if err := s.delegate(ctx, counter); err != nil {
		return err
	}
	s.logger.Infof("authority %s delegated to counter %s", s.authority.Hex(), counter.Hex())

	maxPending := s.options.MaxPending
	if maxPending == 0 {
		maxPending = s.walletPool.GetConfiguredWalletCount() * 10
	}

	var timeout time.Duration
	if s.options.Timeout != "" {
		if timeout, err = time.ParseDuration(s.options.Timeout); err != nil {
			return fmt.Errorf("invalid timeout value: %w", err)
		}
	}

	return scenario.RunTransactionScenario(ctx, scenario.TransactionScenarioOptions{
		TotalCount: s.options.TotalCount,
		Throughput: s.options.Throughput,
		MaxPending: maxPending,
		Timeout:    timeout,
		WalletPool: s.walletPool,
		Logger:     s.logger,
		ProcessNextTxFn: func(ctx context.Context, params *scenario.ProcessNextTxParams) error {
			receiptChan, tx, err := s.sendCounterTx(ctx, params.TxIdx)

			params.NotifySubmitted()
			params.OrderedLogCb(func() {
				if err != nil {
					s.logger.Warnf("could not send transaction: %v", err)
				} else if s.options.LogTxs {
					s.logger.Infof("sent tx #%6d: %v", params.TxIdx+1, tx.Hash().String())
				}
			})

			if _, waitErr := receiptChan.Wait(ctx); waitErr != nil {
				return waitErr
			}
			return err
		},
	})
}

func (s *Scenario) setupClient() (*spamoor.Client, *uint256.Int, *uint256.Int, error) {
	client := s.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0),
		spamoor.WithClientGroup(s.options.ClientGroup),
	)
	if client == nil {
		return nil, nil, nil, scenario.ErrNoClients
	}
	baseFeeWei, tipFeeWei := spamoor.ResolveFees(s.options.BaseFee, s.options.TipFee, s.options.BaseFeeWei, s.options.TipFeeWei)
	feeCap, tipCap, err := s.walletPool.GetSuggestedFees(client, baseFeeWei, tipFeeWei)
	if err != nil {
		return nil, nil, nil, err
	}
	return client, uint256.MustFromBig(feeCap), uint256.MustFromBig(tipCap), nil
}

func (s *Scenario) deployCounter(ctx context.Context) (common.Address, error) {
	deployer := s.walletPool.GetWellKnownWallet("deployer")
	client, feeCap, tipCap, err := s.setupClient()
	if err != nil {
		return common.Address{}, err
	}

	txData, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
		GasFeeCap: feeCap,
		GasTipCap: tipCap,
		Gas:       s.walletPool.EstimateDeployGas(ctx, client, deployer.GetAddress(), nil, counterInitCode),
		Value:     uint256.NewInt(0),
		Data:      counterInitCode,
	})
	if err != nil {
		return common.Address{}, err
	}
	tx, err := deployer.BuildDynamicFeeTx(txData)
	if err != nil {
		return common.Address{}, err
	}

	receipt, err := s.walletPool.GetTxPool().SendAndAwaitTransaction(ctx, deployer, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		Rebroadcast: true,
	})
	if err != nil {
		return common.Address{}, fmt.Errorf("counter deployment failed: %w", err)
	}
	if receipt.Status != txtypes.ReceiptStatusSuccessful {
		return common.Address{}, fmt.Errorf("counter deployment reverted")
	}
	return receipt.ContractAddress, nil
}

func (s *Scenario) delegate(ctx context.Context, counter common.Address) error {
	deployer := s.walletPool.GetWellKnownWallet("deployer")
	authority := s.walletPool.GetWellKnownWallet("authority")
	s.authority = authority.GetAddress()

	client, feeCap, tipCap, err := s.setupClient()
	if err != nil {
		return err
	}

	auth, err := txtypes.SignAuthorization(txtypes.SetCodeAuthorization{
		ChainID: *uint256.MustFromBig(s.walletPool.GetChainId()),
		Address: counter,
		Nonce:   authority.GetNextNonce(),
	}, authority.GetPrivateKey())
	if err != nil {
		return fmt.Errorf("could not sign authorization: %w", err)
	}

	// Authorizations apply before execution, so sending this tx to the authority would already bump the counter.
	deployerAddr := deployer.GetAddress()
	txData, err := txbuilder.SetCodeTx(&txbuilder.TxMetadata{
		GasFeeCap: feeCap,
		GasTipCap: tipCap,
		Gas:       200000,
		To:        &deployerAddr,
		Value:     uint256.NewInt(0),
		AuthList:  []txtypes.SetCodeAuthorization{auth},
	})
	if err != nil {
		return err
	}
	tx, err := deployer.BuildSetCodeTx(txData)
	if err != nil {
		return err
	}

	receipt, err := s.walletPool.GetTxPool().SendAndAwaitTransaction(ctx, deployer, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		Rebroadcast: true,
	})
	if err != nil {
		return fmt.Errorf("delegation failed: %w", err)
	}
	if receipt.Status != txtypes.ReceiptStatusSuccessful {
		return fmt.Errorf("delegation reverted")
	}
	return nil
}

func (s *Scenario) sendCounterTx(ctx context.Context, txIdx uint64) (scenario.ReceiptChan, *txtypes.Transaction, error) {
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(txIdx))
	client := s.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, int(txIdx)),
		spamoor.WithClientGroup(s.options.ClientGroup),
	)
	if wallet == nil {
		return nil, nil, scenario.ErrNoWallet
	}
	if client == nil {
		return nil, nil, scenario.ErrNoClients
	}

	baseFeeWei, tipFeeWei := spamoor.ResolveFees(s.options.BaseFee, s.options.TipFee, s.options.BaseFeeWei, s.options.TipFeeWei)
	feeCap, tipCap, err := s.walletPool.GetSuggestedFees(client, baseFeeWei, tipFeeWei)
	if err != nil {
		return nil, nil, err
	}

	txData, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       s.options.GasLimit,
		To:        &s.authority,
		Value:     uint256.NewInt(0),
	})
	if err != nil {
		return nil, nil, err
	}
	tx, err := wallet.BuildDynamicFeeTx(txData)
	if err != nil {
		return nil, nil, err
	}

	receiptChan := make(scenario.ReceiptChan, 1)
	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		ClientGroup: s.options.ClientGroup,
		Rebroadcast: s.options.Rebroadcast > 0,
		OnComplete: func(tx *txtypes.Transaction, receipt *txtypes.Receipt, err error) {
			receiptChan <- receipt
		},
		LogFn: spamoor.GetDefaultLogFn(s.logger, "", fmt.Sprintf("%6d", txIdx+1), tx),
	})
	if err != nil {
		wallet.MarkSkippedNonce(tx.Nonce())
		return nil, nil, err
	}
	return receiptChan, tx, nil
}
