package factorydeploytx

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/scenarios/factorydeploytx/contract"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/utils"
)

type ScenarioOptions struct {
	TotalCount       uint64 `yaml:"total_count"`
	Throughput       uint64 `yaml:"throughput"`
	MaxPending       uint64 `yaml:"max_pending"`
	MaxWallets       uint64 `yaml:"max_wallets"`
	Rebroadcast      uint64 `yaml:"rebroadcast"`
	BaseFee          uint64 `yaml:"base_fee"`
	TipFee           uint64 `yaml:"tip_fee"`
	GasLimit         uint64 `yaml:"gas_limit"`
	FactoryAddress   string `yaml:"factory_address"`
	InitCode         string `yaml:"init_code"`
	StartSalt        uint64 `yaml:"start_salt"`
	WellKnownFactory bool   `yaml:"well_known_factory"`
	Timeout          string `yaml:"timeout"`
	ClientGroup      string `yaml:"client_group"`
	LogTxs           bool   `yaml:"log_txs"`
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	factoryAddr   common.Address
	initCodeBytes []byte

	pendingWGroup sync.WaitGroup
}

var ScenarioName = "factorydeploytx"
var ScenarioDefaultOptions = ScenarioOptions{
	TotalCount:       0,
	Throughput:       50,
	MaxPending:       0,
	MaxWallets:       0,
	Rebroadcast:      1,
	BaseFee:          20,
	TipFee:           2,
	GasLimit:         2000000,
	FactoryAddress:   "",
	InitCode:         "",
	StartSalt:        0,
	WellKnownFactory: true,
	Timeout:          "",
	ClientGroup:      "",
	LogTxs:           false,
}
var ScenarioDescriptor = scenario.Descriptor{
	Name:           ScenarioName,
	Description:    "Deploy contracts using CREATE2 factory",
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
	flags.Uint64VarP(&s.options.TotalCount, "count", "c", ScenarioDefaultOptions.TotalCount, "Total number of contracts to deploy")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t", ScenarioDefaultOptions.Throughput, "Number of deployment transactions to send per slot")
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", ScenarioDefaultOptions.Rebroadcast, "Enable reliable rebroadcast system")
	flags.Uint64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use in transactions (in gwei)")
	flags.Uint64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use in transactions (in gwei)")
	flags.Uint64Var(&s.options.GasLimit, "gaslimit", ScenarioDefaultOptions.GasLimit, "Gas limit to use in transactions")
	flags.StringVar(&s.options.FactoryAddress, "factory-address", ScenarioDefaultOptions.FactoryAddress, "Address of existing CREATE2 factory (optional)")
	flags.StringVar(&s.options.InitCode, "init-code", ScenarioDefaultOptions.InitCode, "Hex-encoded init code of contract to deploy")
	flags.Uint64Var(&s.options.StartSalt, "start-salt", ScenarioDefaultOptions.StartSalt, "Starting salt value for deployments")
	flags.BoolVar(&s.options.WellKnownFactory, "well-known-factory", ScenarioDefaultOptions.WellKnownFactory, "Use well-known factory deployer wallet")
	flags.StringVar(&s.options.Timeout, "timeout", ScenarioDefaultOptions.Timeout, "Timeout for the scenario (e.g. '1h', '30m', '5s') - empty means no timeout")
	flags.StringVar(&s.options.ClientGroup, "client-group", ScenarioDefaultOptions.ClientGroup, "Client group to use for sending transactions")
	flags.BoolVar(&s.options.LogTxs, "log-txs", ScenarioDefaultOptions.LogTxs, "Log all submitted transactions")
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

	// Set up well-known factory wallet if enabled
	if s.options.WellKnownFactory {
		s.walletPool.AddWellKnownWallet(&spamoor.WellKnownWalletConfig{
			Name:          "create2-factory-deployer",
			RefillAmount:  uint256.NewInt(10000000000000000000), // 10 ETH
			RefillBalance: uint256.NewInt(1000000000000000000),  // 1 ETH
			VeryWellKnown: true,
		})
	}

	if s.options.MaxWallets > 0 {
		s.walletPool.SetWalletCount(s.options.MaxWallets)
	} else if s.options.TotalCount > 0 {
		if s.options.TotalCount < 1000 {
			s.walletPool.SetWalletCount(s.options.TotalCount)
		} else {
			s.walletPool.SetWalletCount(1000)
		}
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

	if s.options.InitCode == "" {
		return errors.New("init-code parameter is required")
	}

	s.initCodeBytes = common.FromHex(s.options.InitCode)
	if len(s.initCodeBytes) == 0 {
		return errors.New("invalid init-code provided")
	}

	// Calculate and log the init code hash for subsequent scenarios
	initCodeHash := crypto.Keccak256Hash(s.initCodeBytes)
	s.logger.Infof("Init code hash: %s", initCodeHash.Hex())

	return nil
}

func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Infof("starting scenario: %s", ScenarioName)
	defer s.logger.Infof("scenario %s finished.", ScenarioName)

	// Deploy or connect to factory
	factoryAddr, err := s.deployFactory(ctx)
	if err != nil {
		s.logger.Errorf("could not deploy/connect to factory: %v", err)
		return err
	}
	s.factoryAddr = factoryAddr
	s.logger.Infof("using CREATE2 factory at: %v", s.factoryAddr.String())

	// Start deploying contracts
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
		var err error
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
		ProcessNextTxFn: func(ctx context.Context, txIdx uint64, onComplete func()) (func(), error) {
			logger := s.logger
			tx, client, wallet, err := s.sendTx(ctx, txIdx, onComplete)
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
					logger.Infof("sent deployment tx #%6d: %v", txIdx+1, tx.Hash().String())
				} else {
					logger.Debugf("sent deployment tx #%6d: %v", txIdx+1, tx.Hash().String())
				}
			}, err
		},
	})

	s.pendingWGroup.Wait()
	s.logger.Infof("finished sending transactions, awaiting block inclusion...")

	return err
}

func (s *Scenario) deployFactory(ctx context.Context) (common.Address, error) {
	if s.options.FactoryAddress != "" {
		return common.HexToAddress(s.options.FactoryAddress), nil
	}

	factoryWallet := s.walletPool.GetWellKnownWallet("create2-factory-deployer")
	if factoryWallet == nil {
		return common.Address{}, fmt.Errorf("factory deployer wallet not available")
	}

	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, s.options.ClientGroup)
	if client == nil {
		return common.Address{}, fmt.Errorf("no client available")
	}

	// Check if factory already exists by checking deployer nonce
	deployerNonce, err := client.GetEthClient().NonceAt(ctx, factoryWallet.GetAddress(), nil)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to get deployer nonce: %w", err)
	}

	if deployerNonce > 0 {
		// Factory was already deployed with nonce 0, calculate its address
		factoryAddr := crypto.CreateAddress(factoryWallet.GetAddress(), 0)
		code, err := client.GetEthClient().CodeAt(ctx, factoryAddr, nil)
		if err == nil && len(code) > 0 {
			s.logger.Infof("Factory already exists at %v", factoryAddr.String())
			return factoryAddr, nil
		}
	}

	feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(client, s.options.BaseFee, s.options.TipFee)
	if err != nil {
		return common.Address{}, err
	}

	tx, err := factoryWallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       2000000,
		Value:     uint256.NewInt(0),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		_, deployTx, _, err := contract.DeployCREATE2Factory(transactOpts, client.GetEthClient())
		return deployTx, err
	})

	if err != nil {
		return common.Address{}, err
	}

	receipt, err := s.walletPool.GetSubmitter().SendAndAwait(ctx, factoryWallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		Rebroadcast: true,
	})
	if err != nil {
		return common.Address{}, err
	}

	if receipt == nil || receipt.ContractAddress == (common.Address{}) {
		return common.Address{}, fmt.Errorf("factory deployment failed")
	}

	s.logger.Infof("deployed CREATE2 factory at: %v (confirmed in block #%v)", receipt.ContractAddress.String(), receipt.BlockNumber.String())
	return receipt.ContractAddress, nil
}

func (s *Scenario) sendTx(ctx context.Context, txIdx uint64, onComplete func()) (*types.Transaction, *spamoor.Client, *spamoor.Wallet, error) {
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, int(txIdx), s.options.ClientGroup)
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(txIdx))
	transactionSubmitted := false

	defer func() {
		if !transactionSubmitted {
			onComplete()
		}
	}()

	if client == nil {
		return nil, client, wallet, fmt.Errorf("no client available")
	}

	feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(client, s.options.BaseFee, s.options.TipFee)
	if err != nil {
		return nil, client, wallet, err
	}

	factory, err := contract.NewCREATE2Factory(s.factoryAddr, client.GetEthClient())
	if err != nil {
		return nil, nil, wallet, err
	}

	salt := [32]byte{}
	saltValue := s.options.StartSalt + txIdx
	// Convert uint64 to bytes32 (big endian)
	for i := 0; i < 8; i++ {
		salt[31-i] = byte(saltValue >> (8 * i))
	}

	tx, err := wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       s.options.GasLimit,
		Value:     uint256.NewInt(0),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		return factory.Deploy(transactOpts, salt, s.initCodeBytes)
	})
	if err != nil {
		return nil, nil, wallet, err
	}

	s.pendingWGroup.Add(1)
	transactionSubmitted = true
	err = s.walletPool.GetSubmitter().Send(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		Rebroadcast: s.options.Rebroadcast > 0,
		OnComplete: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			onComplete()
			s.pendingWGroup.Done()
		},
		OnConfirm: func(tx *types.Transaction, receipt *types.Receipt) {
			if receipt == nil {
				return
			}

			// Calculate deployed contract address
			deployedAddr := common.Address{}
			if len(receipt.Logs) > 0 {
				// Try to extract from ContractDeployed event
				for _, log := range receipt.Logs {
					if len(log.Topics) > 1 {
						deployedAddr = common.HexToAddress(log.Topics[1].Hex())
						break
					}
				}
			}

			txFees := utils.GetTransactionFees(tx, receipt)
			s.logger.WithField("rpc", client.GetName()).Debugf("deployment tx %d confirmed in block #%v. deployed at: %v, total fee: %v gwei (base: %v)",
				txIdx+1, receipt.BlockNumber.String(), deployedAddr.String(), txFees.TotalFeeGwei(), txFees.TxBaseFeeGwei())
		},
		LogFn: spamoor.GetDefaultLogFn(s.logger, "", fmt.Sprintf("%6d", txIdx+1), tx),
	})
	if err != nil {
		// reset nonce if tx was not sent
		wallet.ResetPendingNonce(ctx, client)

		return nil, client, wallet, err
	}

	return tx, client, wallet, nil
}
