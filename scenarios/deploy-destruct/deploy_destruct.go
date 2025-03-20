package deploydestruct

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenarios/deploy-destruct/contract"
	"github.com/ethpandaops/spamoor/scenariotypes"
	"github.com/ethpandaops/spamoor/tester"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/utils"
)

type ScenarioOptions struct {
	TotalCount   uint64
	Throughput   uint64
	MaxPending   uint64
	MaxWallets   uint64
	Rebroadcast  uint64
	BaseFee      uint64
	TipFee       uint64
	Amount       uint64
	GasLimit     uint64
	RandomAmount bool
}

type Scenario struct {
	options ScenarioOptions
	logger  *logrus.Entry
	tester  *tester.Tester

	contractAddr common.Address

	pendingChan   chan bool
	pendingWGroup sync.WaitGroup
}

func NewScenario() scenariotypes.Scenario {
	return &Scenario{
		logger: logrus.WithField("scenario", "erctx"),
	}
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
	flags.Uint64VarP(&s.options.TotalCount, "count", "c", 0, "Total number of transfer transactions to send")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t", 0, "Number of transfer transactions to send per slot")
	flags.Uint64Var(&s.options.MaxPending, "max-pending", 0, "Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", 0, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", 120, "Number of seconds to wait before re-broadcasting a transaction")
	flags.Uint64Var(&s.options.BaseFee, "basefee", 20, "Max fee per gas to use in transfer transactions (in gwei)")
	flags.Uint64Var(&s.options.TipFee, "tipfee", 2, "Max tip per gas to use in transfer transactions (in gwei)")
	flags.Uint64Var(&s.options.Amount, "amount", 20, "Transfer amount per transaction (in gwei)")
	flags.Uint64Var(&s.options.GasLimit, "gaslimit", 10000000, "The gas limit for each deployment test tx")
	flags.BoolVar(&s.options.RandomAmount, "random-amount", false, "Use random amounts for transactions (with --amount as limit)")
	return nil
}

func (s *Scenario) Init(testerCfg *tester.TesterConfig) error {
	if s.options.TotalCount == 0 && s.options.Throughput == 0 {
		return fmt.Errorf("neither total count nor throughput limit set, must define at least one of them (see --help for list of all flags)")
	}

	if s.options.MaxWallets > 0 {
		testerCfg.WalletCount = s.options.MaxWallets
	} else if s.options.TotalCount > 0 {
		if s.options.TotalCount < 1000 {
			testerCfg.WalletCount = s.options.TotalCount
		} else {
			testerCfg.WalletCount = 1000
		}
	} else {
		if s.options.Throughput*10 < 1000 {
			testerCfg.WalletCount = s.options.Throughput * 10
		} else {
			testerCfg.WalletCount = 1000
		}
	}

	if s.options.MaxPending > 0 {
		s.pendingChan = make(chan bool, s.options.MaxPending)
	}

	return nil
}

func (s *Scenario) Run(tester *tester.Tester) error {
	s.tester = tester
	txIdxCounter := uint64(0)
	pendingCount := atomic.Int64{}
	txCount := atomic.Uint64{}
	var lastChan chan bool

	s.logger.Infof("starting scenario: deploy-destruct")
	contractReceipt, _, err := s.sendDeploymentTx()
	if err != nil {
		s.logger.Errorf("could not deploy test contract: %v", err)
		return err
	}
	s.contractAddr = contractReceipt.ContractAddress
	s.logger.Infof("deployed test contract: %v (confirmed in block #%v)", s.contractAddr.String(), contractReceipt.BlockNumber.String())

	initialRate := rate.Limit(float64(s.options.Throughput) / float64(utils.SecondsPerSlot))
	if initialRate == 0 {
		initialRate = rate.Inf
	}
	limiter := rate.NewLimiter(initialRate, 1)

	for {
		if err := limiter.Wait(context.Background()); err != nil {
			s.logger.Debugf("rate limited: %s", err.Error())
			time.Sleep(100 * time.Millisecond)
			continue
		}

		txIdx := txIdxCounter
		txIdxCounter++

		if s.pendingChan != nil {
			// await pending transactions
			s.pendingChan <- true
		}
		pendingCount.Add(1)
		currentChan := make(chan bool, 1)

		go func(txIdx uint64, lastChan, currentChan chan bool) {
			defer func() {
				pendingCount.Add(-1)
				currentChan <- true
			}()

			logger := s.logger
			tx, client, wallet, err := s.sendTx(txIdx)
			if client != nil {
				logger = logger.WithField("rpc", client.GetName())
			}
			if tx != nil {
				logger = logger.WithField("nonce", tx.Nonce())
			}
			if wallet != nil {
				logger = logger.WithField("wallet", s.tester.GetWalletIndex(wallet.GetAddress()))
			}
			if lastChan != nil {
				<-lastChan
				close(lastChan)
			}
			if err != nil {
				logger.Warnf("could not send transaction: %v", err)
				<-s.pendingChan
				return
			}

			txCount.Add(1)
			if lastChan != nil {
				<-lastChan
				close(lastChan)
			}
			logger.Infof("sent tx #%6d: %v", txIdx+1, tx.Hash().String())
		}(txIdx, lastChan, currentChan)

		lastChan = currentChan

		count := txCount.Load() + uint64(pendingCount.Load())
		if s.options.TotalCount > 0 && count >= s.options.TotalCount {
			break
		}
	}
	s.pendingWGroup.Wait()
	s.logger.Infof("finished sending transactions, awaiting block inclusion...")

	return nil
}

func (s *Scenario) sendDeploymentTx() (*types.Receipt, *txbuilder.Client, error) {
	client := s.tester.GetClient(tester.SelectByIndex, 0)
	wallet := s.tester.GetWallet(tester.SelectByIndex, 0)

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
		feeCap, tipCap, err = client.GetSuggestedFee(s.tester.GetContext())
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

	tx, err := wallet.BuildBoundTx(&txbuilder.TxMetadata{
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

	err = s.tester.GetTxPool().SendTransaction(context.Background(), wallet, tx, &txbuilder.SendTransactionOptions{
		Client:              client,
		MaxRebroadcasts:     10,
		RebroadcastInterval: 30 * time.Second,
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

func (s *Scenario) sendTx(txIdx uint64) (*types.Transaction, *txbuilder.Client, *txbuilder.Wallet, error) {
	client := s.tester.GetClient(tester.SelectByIndex, int(txIdx))
	wallet := s.tester.GetWallet(tester.SelectByIndex, int(txIdx))

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
		feeCap, tipCap, err = client.GetSuggestedFee(s.tester.GetContext())
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

	testToken, err := contract.NewContract(s.contractAddr, client.GetEthClient())
	if err != nil {
		return nil, nil, wallet, err
	}

	tx, err := wallet.BuildBoundTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       s.options.GasLimit,
		Value:     uint256.MustFromBig(amount.ToBig()),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		return testToken.Test(transactOpts, big.NewInt(0))
	})
	if err != nil {
		return nil, nil, wallet, err
	}

	rebroadcast := 0
	if s.options.Rebroadcast > 0 {
		rebroadcast = 10
	}

	err = s.tester.GetTxPool().SendTransaction(context.Background(), wallet, tx, &txbuilder.SendTransactionOptions{
		Client:              client,
		MaxRebroadcasts:     rebroadcast,
		RebroadcastInterval: time.Duration(s.options.Rebroadcast) * time.Second,
		OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			defer func() {
				if s.pendingChan != nil {
					time.Sleep(100 * time.Millisecond)
					<-s.pendingChan
				}
				s.pendingWGroup.Done()
			}()

			if err != nil {
				s.logger.WithField("client", client.GetName()).Warnf("tx %6d: await receipt failed: %v", txIdx+1, err)
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

			s.logger.WithField("client", client.GetName()).Debugf(" transaction %d confirmed in block #%v. total fee: %v gwei (base: %v) logs: %v", txIdx+1, receipt.BlockNumber.String(), gweiTotalFee, gweiBaseFee, len(receipt.Logs))
		},
		LogFn: func(client *txbuilder.Client, retry int, rebroadcast int, err error) {
			logger := s.logger.WithField("client", client.GetName())
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
		wallet.ResetPendingNonce(s.tester.GetContext(), client)

		return nil, client, wallet, err
	}

	s.pendingWGroup.Add(1)

	return tx, client, wallet, nil
}
