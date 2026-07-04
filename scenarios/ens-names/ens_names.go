package ensnames

import (
	"context"
	cryptorand "crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/utils"
)

// wallet naming service modes
const (
	walletNamingAll  = "all"
	walletNamingPool = "pool"
	walletNamingOff  = "off"
)

// minRegistrationDuration is the ETHRegistrarController's hardcoded
// MIN_REGISTRATION_DURATION (28 days).
const minRegistrationDuration = 28 * 24 * 3600

// ScenarioOptions configures the ens-names scenario: the standard
// throughput/fee knobs shared by all scenarios plus the ENS-specific
// registration targets, per-operation weights and the wallet naming service.
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

	DeploymentSeed       string `yaml:"deployment_seed"`
	NamesPerWallet       uint64 `yaml:"names_per_wallet"`
	RegistrationDuration uint64 `yaml:"registration_duration"`
	RenewalDuration      uint64 `yaml:"renewal_duration"`
	MinCommitmentAge     uint64 `yaml:"min_commitment_age"`
	MaxCommitmentAge     uint64 `yaml:"max_commitment_age"`

	RotationWeight     uint64 `yaml:"rotation_weight"`
	RenewWeight        uint64 `yaml:"renew_weight"`
	RecordUpdateWeight uint64 `yaml:"record_update_weight"`
	TransferWeight     uint64 `yaml:"transfer_weight"`
	AbandonWeight      uint64 `yaml:"abandon_weight"`
	ReverseWeight      uint64 `yaml:"reverse_weight"`
	WrapWeight         uint64 `yaml:"wrap_weight"`
	ChurnWeight        uint64 `yaml:"churn_weight"`

	WalletNaming  string `yaml:"wallet_naming"`
	NamingPerSlot uint64 `yaml:"naming_per_slot"`

	Timeout           string `yaml:"timeout"`
	ClientGroup       string `yaml:"client_group"`
	DeployClientGroup string `yaml:"deploy_client_group"`
	LogTxs            bool   `yaml:"log_txs"`
}

// Scenario deploys the full ENS stack (registry, .eth registrar, commit-reveal
// controller, public resolver, reverse registrars, name wrapper) at
// deterministic addresses and drives organic ENS usage from child wallets:
// commit-reveal registrations, renewals, record updates, transfers, wrapping
// and short-lived churn names. An opt-out background service additionally
// registers forward + reverse names for every spamoor wallet on the host.
type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	// runID is a random per-run id embedded in all generated labels, so names
	// from different runs (and concurrent scenario instances) never collide.
	runID      string
	deployment *DeploymentInfo

	walletStates    map[common.Address]*walletState
	walletStatesMtx sync.Mutex
}

var ScenarioName = "ens-names"
var ScenarioDefaultOptions = ScenarioOptions{
	TotalCount:  0,
	Throughput:  10,
	MaxPending:  0,
	MaxWallets:  0,
	Rebroadcast: 1,
	BaseFee:     20,
	TipFee:      2,

	DeploymentSeed:       "",
	NamesPerWallet:       3,
	RegistrationDuration: 90 * 24 * 3600,
	RenewalDuration:      30 * 24 * 3600,
	MinCommitmentAge:     10,
	MaxCommitmentAge:     86400,

	RotationWeight:     10,
	RenewWeight:        20,
	RecordUpdateWeight: 25,
	TransferWeight:     10,
	AbandonWeight:      5,
	ReverseWeight:      10,
	WrapWeight:         15,
	ChurnWeight:        15,

	WalletNaming:  walletNamingAll,
	NamingPerSlot: 2,

	Timeout:           "",
	ClientGroup:       "",
	DeployClientGroup: "",
	LogTxs:            false,
}

var ScenarioDescriptor = scenario.Descriptor{
	Name:           ScenarioName,
	Description:    "Deploy the ENS stack and send name registrations, renewals, record updates & transfers",
	DefaultOptions: ScenarioDefaultOptions,
	NewScenario:    newScenario,
}

func newScenario(logger logrus.FieldLogger) scenario.Scenario {
	return &Scenario{
		options:      ScenarioDefaultOptions,
		logger:       logger.WithField("scenario", ScenarioName),
		walletStates: make(map[common.Address]*walletState, 64),
	}
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
	flags.Uint64VarP(&s.options.TotalCount, "count", "c", ScenarioDefaultOptions.TotalCount, "Total number of ENS transactions to send")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t", ScenarioDefaultOptions.Throughput, "Number of ENS transactions to send per slot")
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", ScenarioDefaultOptions.Rebroadcast, "Enable reliable rebroadcast system")
	flags.Float64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use in transactions (in gwei)")
	flags.Float64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use in transactions (in gwei)")
	flags.StringVar(&s.options.BaseFeeWei, "basefee-wei", "", "Max fee per gas in wei (overrides --basefee for L2 sub-gwei fees)")
	flags.StringVar(&s.options.TipFeeWei, "tipfee-wei", "", "Max tip per gas in wei (overrides --tipfee for L2 sub-gwei fees)")
	flags.StringVar(&s.options.DeploymentSeed, "deployment-seed", ScenarioDefaultOptions.DeploymentSeed, "Seed for the deterministic ENS stack addresses (empty = one shared stack per root key)")
	flags.Uint64Var(&s.options.NamesPerWallet, "names-per-wallet", ScenarioDefaultOptions.NamesPerWallet, "Number of .eth names each child wallet registers and maintains")
	flags.Uint64Var(&s.options.RegistrationDuration, "registration-duration", ScenarioDefaultOptions.RegistrationDuration, "Registration duration in seconds (min 28 days, enforced by the ENS controller)")
	flags.Uint64Var(&s.options.RenewalDuration, "renewal-duration", ScenarioDefaultOptions.RenewalDuration, "Renewal duration in seconds")
	flags.Uint64Var(&s.options.MinCommitmentAge, "min-commitment-age", ScenarioDefaultOptions.MinCommitmentAge, "Controller min commitment age in seconds (constructor arg, changes the stack addresses)")
	flags.Uint64Var(&s.options.MaxCommitmentAge, "max-commitment-age", ScenarioDefaultOptions.MaxCommitmentAge, "Controller max commitment age in seconds (constructor arg, changes the stack addresses)")
	flags.Uint64Var(&s.options.RotationWeight, "rotation-weight", ScenarioDefaultOptions.RotationWeight, "Weight of rotation registrations: keep registering fresh names, retiring the oldest from active management (0 = disabled)")
	flags.Uint64Var(&s.options.RenewWeight, "renew-weight", ScenarioDefaultOptions.RenewWeight, "Weight of renew operations (0 = disabled)")
	flags.Uint64Var(&s.options.RecordUpdateWeight, "record-update-weight", ScenarioDefaultOptions.RecordUpdateWeight, "Weight of resolver record update operations (0 = disabled)")
	flags.Uint64Var(&s.options.TransferWeight, "transfer-weight", ScenarioDefaultOptions.TransferWeight, "Weight of name transfer operations (0 = disabled)")
	flags.Uint64Var(&s.options.AbandonWeight, "abandon-weight", ScenarioDefaultOptions.AbandonWeight, "Weight of name abandon operations (0 = disabled)")
	flags.Uint64Var(&s.options.ReverseWeight, "reverse-weight", ScenarioDefaultOptions.ReverseWeight, "Weight of default reverse record update operations (0 = disabled)")
	flags.Uint64Var(&s.options.WrapWeight, "wrap-weight", ScenarioDefaultOptions.WrapWeight, "Weight of NameWrapper wrap/unwrap operations (0 = disabled)")
	flags.Uint64Var(&s.options.ChurnWeight, "churn-weight", ScenarioDefaultOptions.ChurnWeight, "Weight of short-lived churn registrations (0 = disabled unless no other op is feasible)")
	flags.StringVar(&s.options.WalletNaming, "wallet-naming", ScenarioDefaultOptions.WalletNaming, "Wallet naming service scope: all (every spamoor wallet on the host), pool (own wallets only) or off")
	flags.Uint64Var(&s.options.NamingPerSlot, "naming-per-slot", ScenarioDefaultOptions.NamingPerSlot, "Maximum wallets the naming service names per slot")
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

	if s.options.MinCommitmentAge >= s.options.MaxCommitmentAge {
		return fmt.Errorf("min-commitment-age must be lower than max-commitment-age")
	}
	if s.options.RegistrationDuration < minRegistrationDuration {
		return fmt.Errorf("registration-duration must be at least 28 days (%d seconds, enforced by the ENS controller)", minRegistrationDuration)
	}

	switch s.options.WalletNaming {
	case walletNamingAll, walletNamingPool, walletNamingOff:
	default:
		return fmt.Errorf("invalid wallet-naming mode %q (must be all, pool or off)", s.options.WalletNaming)
	}

	// The deployer funds the executor + ENS stack deployment (the NameWrapper
	// deploy alone is a multi-million-gas tx); the nameservice wallet sends
	// the wallet naming txs.
	s.walletPool.AddWellKnownWallet(&spamoor.WellKnownWalletConfig{
		Name:          "deployer",
		RefillAmount:  uint256.NewInt(5000000000000000000), // 5 ETH
		RefillBalance: uint256.NewInt(2000000000000000000), // 2 ETH
	})
	s.walletPool.AddWellKnownWallet(&spamoor.WellKnownWalletConfig{
		Name:          "nameservice",
		RefillAmount:  uint256.NewInt(2000000000000000000), // 2 ETH
		RefillBalance: uint256.NewInt(1000000000000000000), // 1 ETH
	})

	runID := make([]byte, 4)
	if _, err := cryptorand.Read(runID); err != nil {
		return fmt.Errorf("could not generate run id: %w", err)
	}
	s.runID = hex.EncodeToString(runID)

	return nil
}

func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Infof("starting scenario: %s", ScenarioName)
	defer s.logger.Infof("scenario %s finished.", ScenarioName)

	deployment, err := s.DeployContracts(ctx)
	if err != nil {
		s.logger.Errorf("could not deploy ens stack: %v", err)
		return err
	}
	s.deployment = deployment

	// Background wallet naming service: registers forward + reverse names for
	// all spamoor wallets (opt-out via --wallet-naming=off).
	if s.options.WalletNaming != walletNamingOff {
		go s.runNamingService(ctx)
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
		s.logger.Infof("Timeout set to %v", timeout)
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
			receiptChan, tx, client, wallet, desc, err := s.sendTx(ctx, params.TxIdx)
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
					logger.Infof("sent tx #%6d: %v (%s)", params.TxIdx+1, tx.Hash().String(), desc)
				} else {
					logger.Debugf("sent tx #%6d: %v (%s)", params.TxIdx+1, tx.Hash().String(), desc)
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

// sendTx builds and submits one ENS action for the given tx index, wiring the
// action's name-state bookkeeping to the transaction result.
func (s *Scenario) sendTx(ctx context.Context, txIdx uint64) (scenario.ReceiptChan, *types.Transaction, *spamoor.Client, *spamoor.Wallet, string, error) {
	client := s.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, int(txIdx)),
		spamoor.WithClientGroup(s.options.ClientGroup),
	)
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(txIdx))

	if client == nil {
		return nil, nil, client, wallet, "", scenario.ErrNoClients
	}
	if wallet == nil {
		return nil, nil, client, wallet, "", scenario.ErrNoWallet
	}

	if err := wallet.ResetNoncesIfNeeded(ctx, client); err != nil {
		return nil, nil, client, wallet, "", err
	}

	baseFeeWei, tipFeeWei := spamoor.ResolveFees(s.options.BaseFee, s.options.TipFee, s.options.BaseFeeWei, s.options.TipFeeWei)
	feeCap, tipCap, err := s.walletPool.GetSuggestedFees(client, baseFeeWei, tipFeeWei)
	if err != nil {
		return nil, nil, client, wallet, "", err
	}

	walletIdx := txIdx % s.walletPool.GetWalletCount()
	tx, onResult, desc, err := s.buildActionTx(ctx, wallet, walletIdx, txIdx, feeCap, tipCap)
	if err != nil {
		return nil, nil, client, wallet, "", err
	}

	receiptChan := make(scenario.ReceiptChan, 1)
	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		ClientGroup: s.options.ClientGroup,
		Rebroadcast: s.options.Rebroadcast > 0,
		OnComplete: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			onResult(err == nil && receipt != nil && receipt.Status == types.ReceiptStatusSuccessful)
			receiptChan <- receipt
		},
		OnConfirm: func(tx *types.Transaction, receipt *types.Receipt) {
			txFees := utils.GetTransactionFees(tx, receipt)
			logger := s.logger.WithField("rpc", client.GetName())
			if receipt.Status != types.ReceiptStatusSuccessful {
				logger.Warnf("transaction %d reverted: %v (%s)", txIdx+1, tx.Hash().Hex(), desc)
				return
			}
			logger.Debugf(" transaction %d confirmed in block #%v. total fee: %v gwei (base: %v) logs: %v (%s)",
				txIdx+1, receipt.BlockNumber.String(), txFees.TotalFeeGweiString(), txFees.TxBaseFeeGweiString(), len(receipt.Logs), desc)
		},
		LogFn: spamoor.GetDefaultLogFn(s.logger, "", fmt.Sprintf("%6d", txIdx+1), tx),
	})
	if err != nil {
		// The tx never went out: roll back the name state and free the nonce.
		onResult(false)
		wallet.MarkSkippedNonce(tx.Nonce())
		return nil, nil, client, wallet, desc, err
	}

	return receiptChan, tx, client, wallet, desc, nil
}

// deployClientGroup returns the client group used for deployments, falling
// back to the general client group when unset.
func (s *Scenario) deployClientGroup() string {
	if s.options.DeployClientGroup != "" {
		return s.options.DeployClientGroup
	}
	return s.options.ClientGroup
}
