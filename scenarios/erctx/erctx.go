package erctx

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"

	"github.com/ethpandaops/spamoor/scenarios/erctx/contract"
	"github.com/ethpandaops/spamoor/scenariotypes"
	"github.com/ethpandaops/spamoor/spamoortypes"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/utils"
)

type ScenarioOptions struct {
	TotalCount   uint64 `yaml:"total_count"`
	Throughput   uint64 `yaml:"throughput"`
	MaxPending   uint64 `yaml:"max_pending"`
	MaxWallets   uint64 `yaml:"max_wallets"`
	Rebroadcast  uint64 `yaml:"rebroadcast"`
	BaseFee      uint64 `yaml:"base_fee"`
	TipFee       uint64 `yaml:"tip_fee"`
	Amount       uint64 `yaml:"amount"`
	RandomAmount bool   `yaml:"random_amount"`
	RandomTarget bool   `yaml:"random_target"`
	Timeout      string `yaml:"timeout"`
	ClientGroup  string `yaml:"client_group"`
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool spamoortypes.WalletPool

	contractAddr common.Address

	pendingWGroup sync.WaitGroup
}

var ScenarioName = "erctx"
var ScenarioDefaultOptions = ScenarioOptions{
	TotalCount:   0,
	Throughput:   200,
	MaxPending:   0,
	MaxWallets:   0,
	Rebroadcast:  1,
	BaseFee:      20,
	TipFee:       2,
	Amount:       20,
	RandomAmount: false,
	RandomTarget: false,
	Timeout:      "",
	ClientGroup:  "",
}
var ScenarioDescriptor = scenariotypes.ScenarioDescriptor{
	Name:           ScenarioName,
	Description:    "Send ERC20 transactions with different configurations",
	DefaultOptions: ScenarioDefaultOptions,
	NewScenario:    newScenario,
}

func newScenario(logger logrus.FieldLogger) scenariotypes.Scenario {
	return &Scenario{
		options: ScenarioDefaultOptions,
		logger:  logger.WithField("scenario", ScenarioName),
	}
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
	flags.Uint64VarP(&s.options.TotalCount, "count", "c", ScenarioDefaultOptions.TotalCount, "Total number of transfer transactions to send")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t", ScenarioDefaultOptions.Throughput, "Number of transfer transactions to send per slot")
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", ScenarioDefaultOptions.Rebroadcast, "Enable reliable rebroadcast system")
	flags.Uint64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use in transfer transactions (in gwei)")
	flags.Uint64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use in transfer transactions (in gwei)")
	flags.Uint64Var(&s.options.Amount, "amount", ScenarioDefaultOptions.Amount, "Transfer amount per transaction (in gwei)")
	flags.BoolVar(&s.options.RandomAmount, "random-amount", ScenarioDefaultOptions.RandomAmount, "Use random amounts for transactions (with --amount as limit)")
	flags.BoolVar(&s.options.RandomTarget, "random-target", ScenarioDefaultOptions.RandomTarget, "Use random to addresses for transactions")
	flags.StringVar(&s.options.Timeout, "timeout", ScenarioDefaultOptions.Timeout, "Timeout for the scenario (e.g. '1h', '30m', '5s') - empty means no timeout")
	flags.StringVar(&s.options.ClientGroup, "client-group", ScenarioDefaultOptions.ClientGroup, "Client group to use for sending transactions")
	return nil
}

func (s *Scenario) Init(options *scenariotypes.ScenarioOptions) error {
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

	return nil
}

func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Infof("starting scenario: %s", ScenarioName)
	defer s.logger.Infof("scenario %s finished.", ScenarioName)

	// deploy erc20 contract
	contractReceipt, _, err := s.sendDeploymentTx(ctx)
	if err != nil {
		s.logger.Errorf("could not deploy token contract: %v", err)
		return err
	}
	if contractReceipt == nil {
		return fmt.Errorf("could not deploy token contract: %w", err)
	}
	s.contractAddr = contractReceipt.ContractAddress
	s.logger.Infof("deployed token contract: %v (confirmed in block #%v)", s.contractAddr.String(), contractReceipt.BlockNumber.String())

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
		var err error
		timeout, err = time.ParseDuration(s.options.Timeout)
		if err != nil {
			return fmt.Errorf("invalid timeout value: %v", err)
		}
		s.logger.Infof("Timeout set to %v", timeout)
	}

	err = utils.RunTransactionScenario(ctx, utils.TransactionScenarioOptions{
		TotalCount:                  s.options.TotalCount,
		Throughput:                  s.options.Throughput,
		MaxPending:                  maxPending,
		ThroughputIncrementInterval: 0,
		Timeout:                     timeout,

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
				} else {
					logger.Infof("sent tx #%6d: %v", txIdx+1, tx.Hash().String())
				}
			}, err
		},
	})

	s.logger.Infof("finished sending transactions, awaiting block inclusion...")
	s.pendingWGroup.Wait()

	return err
}

func (s *Scenario) sendDeploymentTx(ctx context.Context) (*types.Receipt, spamoortypes.Client, error) {
	client := s.walletPool.GetClient(spamoortypes.SelectClientByIndex, 0, s.options.ClientGroup)
	wallet := s.walletPool.GetWallet(spamoortypes.SelectWalletByIndex, 0)

	var feeCap *big.Int
	var tipCap *big.Int

	if client == nil {
		return nil, client, fmt.Errorf("no client available")
	}

	if s.options.BaseFee > 0 {
		feeCap = new(big.Int).Mul(big.NewInt(int64(s.options.BaseFee)), big.NewInt(1000000000))
	}
	if s.options.TipFee > 0 {
		tipCap = new(big.Int).Mul(big.NewInt(int64(s.options.TipFee)), big.NewInt(1000000000))
	}

	if feeCap == nil || tipCap == nil {
		var err error
		feeCap, tipCap, err = client.GetSuggestedFee(s.walletPool.GetContext())
		if err != nil {
			return nil, client, err
		}
	}

	if feeCap.Cmp(big.NewInt(1000000000)) < 0 {
		feeCap = big.NewInt(1000000000)
	}
	if tipCap.Cmp(big.NewInt(1000000000)) < 0 {
		tipCap = big.NewInt(1000000000)
	}

	tx, err := wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       2000000,
		Value:     uint256.NewInt(0),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		_, deployTx, _, err := contract.DeployContract(transactOpts, client.GetEthClient())
		return deployTx, err
	})

	if err != nil {
		return nil, nil, err
	}

	var txReceipt *types.Receipt
	var txErr error
	txWg := sync.WaitGroup{}
	txWg.Add(1)

	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &spamoortypes.SendTransactionOptions{
		Client:      client,
		Rebroadcast: true,
		OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			defer func() {
				txWg.Done()
			}()

			txErr = err
			txReceipt = receipt
		},
	})
	if err != nil {
		return nil, client, err
	}

	txWg.Wait()
	if txErr != nil {
		return nil, client, err
	}
	return txReceipt, client, nil
}

func (s *Scenario) sendTx(ctx context.Context, txIdx uint64, onComplete func()) (*types.Transaction, spamoortypes.Client, spamoortypes.Wallet, error) {
	client := s.walletPool.GetClient(spamoortypes.SelectClientByIndex, int(txIdx), s.options.ClientGroup)
	wallet := s.walletPool.GetWallet(spamoortypes.SelectWalletByIndex, int(txIdx))
	transactionSubmitted := false

	defer func() {
		if !transactionSubmitted {
			onComplete()
		}
	}()

	if client == nil {
		return nil, client, wallet, fmt.Errorf("no client available")
	}

	var feeCap *big.Int
	var tipCap *big.Int

	if s.options.BaseFee > 0 {
		feeCap = new(big.Int).Mul(big.NewInt(int64(s.options.BaseFee)), big.NewInt(1000000000))
	}
	if s.options.TipFee > 0 {
		tipCap = new(big.Int).Mul(big.NewInt(int64(s.options.TipFee)), big.NewInt(1000000000))
	}

	if feeCap == nil || tipCap == nil {
		var err error
		feeCap, tipCap, err = client.GetSuggestedFee(s.walletPool.GetContext())
		if err != nil {
			return nil, client, wallet, err
		}
	}

	if feeCap.Cmp(big.NewInt(1000000000)) < 0 {
		feeCap = big.NewInt(1000000000)
	}
	if tipCap.Cmp(big.NewInt(1000000000)) < 0 {
		tipCap = big.NewInt(1000000000)
	}

	amount := uint256.NewInt(s.options.Amount)
	amount = amount.Mul(amount, uint256.NewInt(1000000000))
	if s.options.RandomAmount {
		n, err := rand.Int(rand.Reader, amount.ToBig())
		if err == nil {
			amount = uint256.MustFromBig(n)
		}
	}

	toAddr := s.walletPool.GetWallet(spamoortypes.SelectWalletByIndex, int(txIdx)+1).GetAddress()
	if s.options.RandomTarget {
		addrBytes := make([]byte, 20)
		rand.Read(addrBytes)
		toAddr = common.Address(addrBytes)
	}

	testToken, err := contract.NewContract(s.contractAddr, client.GetEthClient())
	if err != nil {
		return nil, nil, wallet, err
	}

	tx, err := wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       100000,
		Value:     uint256.NewInt(0),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		return testToken.TransferMint(transactOpts, toAddr, amount.ToBig())
	})
	if err != nil {
		return nil, nil, wallet, err
	}

	s.pendingWGroup.Add(1)
	transactionSubmitted = true
	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &spamoortypes.SendTransactionOptions{
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

			effectiveGasPrice := receipt.EffectiveGasPrice
			if effectiveGasPrice == nil {
				effectiveGasPrice = big.NewInt(0)
			}
			feeAmount := new(big.Int).Mul(effectiveGasPrice, big.NewInt(int64(receipt.GasUsed)))
			totalAmount := new(big.Int).Add(tx.Value(), feeAmount)
			wallet.SubBalance(totalAmount)

			gweiTotalFee := new(big.Int).Div(feeAmount, big.NewInt(1000000000))
			gweiBaseFee := new(big.Int).Div(effectiveGasPrice, big.NewInt(1000000000))

			s.logger.WithField("rpc", client.GetName()).Debugf(" transaction %d confirmed in block #%v. total fee: %v gwei (base: %v) logs: %v", txIdx+1, receipt.BlockNumber.String(), gweiTotalFee, gweiBaseFee, len(receipt.Logs))
		},
		LogFn: func(client spamoortypes.Client, retry int, rebroadcast int, err error) {
			logger := s.logger.WithField("rpc", client.GetName())
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
