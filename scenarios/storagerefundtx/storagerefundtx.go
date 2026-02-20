package storagerefundtx

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/scenarios/storagerefundtx/contract"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/utils"
)

// ScenarioOptions defines the configurable options for the storage refund scenario.
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
	SlotsPerCall      uint64  `yaml:"slots_per_call"`
	GasLimit          uint64  `yaml:"gas_limit"`
	Timeout           string  `yaml:"timeout"`
	ClientGroup       string  `yaml:"client_group"`
	DeployClientGroup string  `yaml:"deploy_client_group"`
	LogTxs            bool    `yaml:"log_txs"`
}

// Scenario implements the storage refund transaction scenario.
type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	contractAddr common.Address
}

// ScenarioName is the unique identifier for this scenario.
var ScenarioName = "storagerefundtx"

// ScenarioDefaultOptions defines the default configuration values.
var ScenarioDefaultOptions = ScenarioOptions{
	TotalCount:        0,
	Throughput:        10,
	MaxPending:        0,
	MaxWallets:        0,
	Rebroadcast:       1,
	BaseFee:           20,
	TipFee:            2,
	SlotsPerCall:      500,
	GasLimit:          0,
	Timeout:           "",
	ClientGroup:       "",
	DeployClientGroup: "",
	LogTxs:            false,
}

// ScenarioDescriptor is the registration descriptor for the scenario system.
var ScenarioDescriptor = scenario.Descriptor{
	Name: ScenarioName,
	Description: "Send transactions that write and clear storage slots " +
		"to trigger gas refunds (for EIP-7778 testing)",
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
	flags.Uint64VarP(&s.options.TotalCount, "count", "c",
		ScenarioDefaultOptions.TotalCount,
		"Total number of transactions to send")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t",
		ScenarioDefaultOptions.Throughput,
		"Number of transactions to send per slot")
	flags.Uint64Var(&s.options.MaxPending, "max-pending",
		ScenarioDefaultOptions.MaxPending,
		"Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets",
		ScenarioDefaultOptions.MaxWallets,
		"Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast",
		ScenarioDefaultOptions.Rebroadcast,
		"Enable reliable rebroadcast system")
	flags.Float64Var(&s.options.BaseFee, "basefee",
		ScenarioDefaultOptions.BaseFee,
		"Max fee per gas to use in transactions (in gwei)")
	flags.Float64Var(&s.options.TipFee, "tipfee",
		ScenarioDefaultOptions.TipFee,
		"Max tip per gas to use in transactions (in gwei)")
	flags.StringVar(&s.options.BaseFeeWei, "basefee-wei", "",
		"Max fee per gas in wei (overrides --basefee for L2 sub-gwei fees)")
	flags.StringVar(&s.options.TipFeeWei, "tipfee-wei", "",
		"Max tip per gas in wei (overrides --tipfee for L2 sub-gwei fees)")
	flags.Uint64Var(&s.options.SlotsPerCall, "slots-per-call",
		ScenarioDefaultOptions.SlotsPerCall,
		"Number of storage slots to write and clear per transaction")
	flags.Uint64Var(&s.options.GasLimit, "gaslimit",
		ScenarioDefaultOptions.GasLimit,
		"Gas limit per transaction (0 = auto-estimate based on slots-per-call)")
	flags.StringVar(&s.options.Timeout, "timeout",
		ScenarioDefaultOptions.Timeout,
		"Timeout for the scenario (e.g. '1h', '30m', '5s') - empty means no timeout")
	flags.StringVar(&s.options.ClientGroup, "client-group",
		ScenarioDefaultOptions.ClientGroup,
		"Client group to use for sending transactions")
	flags.StringVar(&s.options.DeployClientGroup, "deploy-client-group",
		ScenarioDefaultOptions.DeployClientGroup,
		"Client group to use for deployments")
	flags.BoolVar(&s.options.LogTxs, "log-txs",
		ScenarioDefaultOptions.LogTxs,
		"Log all submitted transactions")
	return nil
}

// Init initializes the scenario with the given options.
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

	s.walletPool.AddWellKnownWallet(&spamoor.WellKnownWalletConfig{
		Name:          "deployer",
		RefillAmount:  uint256.NewInt(1000000000000000000), // 1 ETH
		RefillBalance: uint256.NewInt(500000000000000000),  // 0.5 ETH
	})

	if s.options.TotalCount == 0 && s.options.Throughput == 0 {
		return fmt.Errorf(
			"neither total count nor throughput limit set, " +
				"must define at least one of them (see --help for list of all flags)",
		)
	}

	if s.options.SlotsPerCall == 0 {
		return fmt.Errorf("slots-per-call must be greater than 0")
	}

	return nil
}

// Run deploys the contract and starts sending transactions.
func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Infof("starting scenario: %s", ScenarioName)
	defer s.logger.Infof("scenario %s finished.", ScenarioName)

	// deploy storage refund contract
	contractReceipt, _, err := s.sendDeploymentTx(ctx)
	if err != nil {
		s.logger.Errorf("could not deploy storage refund contract: %v", err)
		return err
	}

	if contractReceipt == nil {
		return fmt.Errorf("could not deploy storage refund contract: receipt is nil")
	}

	s.contractAddr = contractReceipt.ContractAddress
	s.logger.Infof(
		"deployed storage refund contract: %v (confirmed in block #%v)",
		s.contractAddr.String(), contractReceipt.BlockNumber.String(),
	)

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
			return fmt.Errorf("invalid timeout value: %w", err)
		}

		s.logger.Infof("timeout set to %v", timeout)
	}

	err = scenario.RunTransactionScenario(
		ctx, scenario.TransactionScenarioOptions{
			TotalCount:                  s.options.TotalCount,
			Throughput:                  s.options.Throughput,
			MaxPending:                  maxPending,
			ThroughputIncrementInterval: 0,
			Timeout:                     timeout,
			WalletPool:                  s.walletPool,

			Logger: s.logger,
			ProcessNextTxFn: func(
				ctx context.Context,
				params *scenario.ProcessNextTxParams,
			) error {
				logger := s.logger
				receiptChan, tx, client, wallet, txErr := s.sendTx(
					ctx, params.TxIdx,
				)
				if client != nil {
					logger = logger.WithField("rpc", client.GetName())
				}
				if tx != nil {
					logger = logger.WithField("nonce", tx.Nonce())
				}
				if wallet != nil {
					logger = logger.WithField(
						"wallet",
						s.walletPool.GetWalletName(wallet.GetAddress()),
					)
				}

				params.NotifySubmitted()
				params.OrderedLogCb(func() {
					if txErr != nil {
						logger.Warnf("could not send transaction: %v", txErr)
					} else if s.options.LogTxs {
						logger.Infof(
							"sent tx #%6d: %v",
							params.TxIdx+1, tx.Hash().String(),
						)
					} else {
						logger.Debugf(
							"sent tx #%6d: %v",
							params.TxIdx+1, tx.Hash().String(),
						)
					}
				})

				// wait for receipt
				if _, waitErr := receiptChan.Wait(ctx); waitErr != nil {
					return waitErr
				}

				return txErr
			},
		},
	)

	return err
}

func (s *Scenario) sendDeploymentTx(
	ctx context.Context,
) (*types.Receipt, *spamoor.Client, error) {
	deployClientGroup := s.options.DeployClientGroup
	if deployClientGroup == "" {
		deployClientGroup = s.options.ClientGroup
	}

	client := s.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0),
		spamoor.WithClientGroup(deployClientGroup),
	)
	wallet := s.walletPool.GetWellKnownWallet("deployer")

	if client == nil {
		return nil, client, scenario.ErrNoClients
	}

	if wallet == nil {
		return nil, client, scenario.ErrNoWallet
	}

	baseFeeWei, tipFeeWei := spamoor.ResolveFees(
		s.options.BaseFee, s.options.TipFee,
		s.options.BaseFeeWei, s.options.TipFeeWei,
	)

	feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(
		client, baseFeeWei, tipFeeWei,
	)
	if err != nil {
		return nil, client, err
	}

	tx, err := wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       2000000,
		Value:     uint256.NewInt(0),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		_, deployTx, _, deployErr := contract.DeployStorageRefund(
			transactOpts, client.GetEthClient(),
		)
		return deployTx, deployErr
	})
	if err != nil {
		return nil, nil, err
	}

	receipt, err := s.walletPool.GetTxPool().SendAndAwaitTransaction(
		ctx, wallet, tx, &spamoor.SendTransactionOptions{
			Client:      client,
			ClientGroup: deployClientGroup,
			Rebroadcast: true,
		},
	)
	if err != nil {
		return nil, client, err
	}

	return receipt, client, nil
}

// gasLimitForSlots estimates the gas limit needed for a given number of
// storage slots to write and clear. Accounts for EIP-2929 cold storage
// access surcharges on top of base SSTORE costs:
//   - Write (zero to non-zero, cold): 20,000 + 2,100 = 22,100 gas
//   - Clear (non-zero to zero, cold): 2,900 + 2,100 = 5,000 gas
//   - Mapping keccak256 + loop overhead: ~200 gas per slot
//
// Plus fixed overhead for function call, state variable reads/writes.
func gasLimitForSlots(slotsPerCall uint64) uint64 {
	// Steady-state per-slot cost: write + clear + overhead
	writeGas := uint64(22300)   // SSTORE_SET(20k) + COLD_SLOAD(2.1k) + margin
	clearGas := uint64(5200)    // SSTORE_RESET(2.9k) + COLD_SLOAD(2.1k) + margin
	loopOverhead := uint64(200) // keccak256, stack ops, jump per iteration
	overhead := uint64(100000)  // function dispatch, state var access, margin

	return slotsPerCall*(writeGas+clearGas+loopOverhead) + overhead
}

func (s *Scenario) sendTx(
	ctx context.Context, txIdx uint64,
) (scenario.ReceiptChan, *types.Transaction, *spamoor.Client, *spamoor.Wallet, error) {
	client := s.walletPool.GetClient(
		spamoor.WithClientSelectionMode(
			spamoor.SelectClientByIndex, int(txIdx),
		),
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

	baseFeeWei, tipFeeWei := spamoor.ResolveFees(
		s.options.BaseFee, s.options.TipFee,
		s.options.BaseFeeWei, s.options.TipFeeWei,
	)

	feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(
		client, baseFeeWei, tipFeeWei,
	)
	if err != nil {
		return nil, nil, client, wallet, err
	}

	gasLimit := s.options.GasLimit
	if gasLimit == 0 {
		gasLimit = gasLimitForSlots(s.options.SlotsPerCall)
	}

	storageRefund, err := contract.NewStorageRefund(
		s.contractAddr, client.GetEthClient(),
	)
	if err != nil {
		return nil, nil, client, wallet, err
	}

	tx, err := wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       gasLimit,
		Value:     uint256.NewInt(0),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		return storageRefund.Execute(
			transactOpts,
			new(big.Int).SetUint64(s.options.SlotsPerCall),
		)
	})
	if err != nil {
		return nil, nil, client, wallet, err
	}

	receiptChan := make(scenario.ReceiptChan, 1)
	err = s.walletPool.GetTxPool().SendTransaction(
		ctx, wallet, tx, &spamoor.SendTransactionOptions{
			Client:      client,
			ClientGroup: s.options.ClientGroup,
			Rebroadcast: s.options.Rebroadcast > 0,
			OnComplete: func(
				tx *types.Transaction,
				receipt *types.Receipt,
				err error,
			) {
				receiptChan <- receipt
			},
			OnConfirm: func(tx *types.Transaction, receipt *types.Receipt) {
				txFees := utils.GetTransactionFees(tx, receipt)
				s.logger.WithField("rpc", client.GetName()).Debugf(
					" transaction %d confirmed in block #%v. "+
						"total fee: %v gwei (base: %v) "+
						"gas used: %v",
					txIdx+1,
					receipt.BlockNumber.String(),
					txFees.TotalFeeGweiString(),
					txFees.TxBaseFeeGweiString(),
					receipt.GasUsed,
				)
			},
			LogFn: spamoor.GetDefaultLogFn(
				s.logger, "", fmt.Sprintf("%6d", txIdx+1), tx,
			),
		},
	)
	if err != nil {
		wallet.MarkSkippedNonce(tx.Nonce())
		return nil, nil, client, wallet, err
	}

	return receiptChan, tx, client, wallet, nil
}
