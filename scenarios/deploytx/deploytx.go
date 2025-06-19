package deploytx

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
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

type ScenarioOptions struct {
	TotalCount    uint64 `yaml:"total_count"`
	Throughput    uint64 `yaml:"throughput"`
	MaxPending    uint64 `yaml:"max_pending"`
	MaxWallets    uint64 `yaml:"max_wallets"`
	Rebroadcast   uint64 `yaml:"rebroadcast"`
	GasLimit      uint64 `yaml:"gas_limit"`
	BaseFee       uint64 `yaml:"base_fee"`
	TipFee        uint64 `yaml:"tip_fee"`
	Bytecodes     string `yaml:"bytecodes"`
	BytecodesFile string `yaml:"bytecodes_file"`
	Timeout       string `yaml:"timeout"`
	ClientGroup   string `yaml:"client_group"`
	LogTxs        bool   `yaml:"log_txs"`
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	bytecodes [][]byte

	pendingWGroup sync.WaitGroup
}

var ScenarioName = "deploytx"
var ScenarioDefaultOptions = ScenarioOptions{
	TotalCount:    0,
	Throughput:    10,
	MaxPending:    0,
	MaxWallets:    0,
	Rebroadcast:   1,
	GasLimit:      1000000,
	BaseFee:       20,
	TipFee:        2,
	Bytecodes:     "",
	BytecodesFile: "",
	Timeout:       "",
	ClientGroup:   "",
	LogTxs:        false,
}
var ScenarioDescriptor = scenario.Descriptor{
	Name:           ScenarioName,
	Description:    "Send deploy transactions",
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
	flags.Uint64VarP(&s.options.TotalCount, "count", "c", ScenarioDefaultOptions.TotalCount, "Total number of deployment transactions to send")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t", ScenarioDefaultOptions.Throughput, "Number of deployment transactions to send per slot")
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", ScenarioDefaultOptions.Rebroadcast, "Enable reliable rebroadcast system")
	flags.Uint64Var(&s.options.GasLimit, "gaslimit", ScenarioDefaultOptions.GasLimit, "Gas limit to use in deployment transactions (in gwei)")
	flags.Uint64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use in deployment transactions (in gwei)")
	flags.Uint64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use in deployment transactions (in gwei)")
	flags.StringVar(&s.options.Bytecodes, "bytecodes", ScenarioDefaultOptions.Bytecodes, "Bytecodes to deploy (, separated list of hex bytecodes)")
	flags.StringVar(&s.options.BytecodesFile, "bytecodes-file", ScenarioDefaultOptions.BytecodesFile, "File with bytecodes to deploy (list with hex bytecodes)")
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

	s.bytecodes = [][]byte{}
	if s.options.Bytecodes != "" {
		for _, hexStr := range strings.Split(s.options.Bytecodes, ",") {
			s.bytecodes = append(s.bytecodes, common.FromHex(hexStr))
		}
	}
	if s.options.BytecodesFile != "" {
		fp, err := os.Open(s.options.BytecodesFile)
		if err != nil {
			return fmt.Errorf("cannot open bytecodes list: %w", err)
		}
		defer fp.Close()
		scanner := bufio.NewScanner(fp)
		for scanner.Scan() {
			hexStr := strings.Trim(scanner.Text(), " \t")
			if strings.HasPrefix(hexStr, "#") || hexStr == "" {
				continue
			}
			s.bytecodes = append(s.bytecodes, common.FromHex(hexStr))
		}
	}

	if len(s.bytecodes) == 0 {
		return errors.New(`no bytecodes provided, please specify at least one bytecode to deploy(eg: --bytecodes "0x1234,0x5678")`)
	}

	return nil
}

func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Infof("starting scenario: %s", ScenarioName)
	defer s.logger.Infof("scenario %s finished.", ScenarioName)

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

	err := scenario.RunTransactionScenario(ctx, scenario.TransactionScenarioOptions{
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
					logger.Infof("sent tx #%6d: %v", txIdx+1, tx.Hash().String())
				} else {
					logger.Debugf("sent tx #%6d: %v", txIdx+1, tx.Hash().String())
				}
			}, err
		},
	})

	s.pendingWGroup.Wait()
	s.logger.Infof("finished sending transactions, awaiting block inclusion...")

	return err
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

	// Determine gas limit: use block gas limit if GasLimit is 0
	gasLimit := s.options.GasLimit
	if gasLimit == 0 {
		var err error
		gasLimit, err = s.walletPool.GetTxPool().GetCurrentGasLimitWithInit()
		if err != nil {
			s.logger.Warnf("tx %6d: failed to fetch current gas limit: %v, using fallback", txIdx+1, err)
			gasLimit = 30000000
		} else if gasLimit == 0 {
			// Final fallback to a reasonable default if no block gas limit is available
			gasLimit = 30000000
			s.logger.Warnf("tx %6d: no gas limit available, using fallback %v", txIdx+1, gasLimit)
		} else {
			s.logger.Debugf("tx %6d: using block gas limit %v", txIdx+1, gasLimit)
		}
	}

	deployData := s.bytecodes[int(txIdx)%len(s.bytecodes)]
	txData, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       s.options.GasLimit,
		To:        nil,
		Value:     uint256.NewInt(0),
		Data:      deployData,
	})
	if err != nil {
		return nil, nil, wallet, err
	}

	tx, err := wallet.BuildDynamicFeeTx(txData)
	if err != nil {
		return nil, nil, wallet, err
	}

	s.pendingWGroup.Add(1)
	transactionSubmitted = true
	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		Rebroadcast: s.options.Rebroadcast > 0,
		OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			defer func() {
				onComplete()
				s.pendingWGroup.Done()
			}()

			if err != nil {
				s.logger.WithField("rpc", client.GetName()).Warnf("tx %6d: await receipt failed: %v", txIdx+1, err)
				return
			}
			if receipt == nil {
				return
			}

			txFees := utils.GetTransactionFees(tx, receipt)
			s.logger.WithField("rpc", client.GetName()).Debugf(" transaction %d confirmed in block #%v. total fee: %v gwei (base: %v) logs: %v", txIdx+1, receipt.BlockNumber.String(), txFees.TotalFeeGwei(), txFees.TxBaseFeeGwei(), len(receipt.Logs))
		},
		LogFn: func(client *spamoor.Client, retry int, rebroadcast int, err error) {
			logger := s.logger.WithField("rpc", client.GetName()).WithField("nonce", tx.Nonce())
			if retry == 0 && rebroadcast > 0 {
				logger.Infof("rebroadcasting tx %6d", txIdx+1)
			}
			if retry > 0 {
				logger = logger.WithField("retry", retry)
			}
			if rebroadcast > 0 {
				logger = logger.WithField("rebroadcast", rebroadcast)
			}
			if err != nil {
				logger.Debugf("failed sending tx %6d: %v", txIdx+1, err)
			} else if retry > 0 || rebroadcast > 0 {
				logger.Debugf("successfully sent tx %6d", txIdx+1)
			}
		},
	})
	if err != nil {
		// reset nonce if tx was not sent
		wallet.ResetPendingNonce(ctx, client)

		return nil, client, wallet, err
	}

	return tx, client, wallet, nil
}
