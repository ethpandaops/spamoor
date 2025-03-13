package blobcombined

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenariotypes"
	"github.com/ethpandaops/spamoor/tester"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/utils"
)

type ScenarioOptions struct {
	TotalCount      uint64
	Throughput      uint64
	Sidecars        uint64
	MaxPending      uint64
	MaxWallets      uint64
	Replace         uint64
	MaxReplacements uint64
	Rebroadcast     uint64
	BaseFee         uint64
	TipFee          uint64
	BlobFee         uint64
}

type Scenario struct {
	options ScenarioOptions
	logger  *logrus.Entry
	tester  *tester.Tester

	pendingChan   chan bool
	pendingWGroup sync.WaitGroup
}

func NewScenario() scenariotypes.Scenario {
	return &Scenario{
		logger: logrus.WithField("scenario", "blob-combined"),
	}
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
	flags.Uint64VarP(&s.options.TotalCount, "count", "c", 0, "Total number of blob transactions to send")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t", 0, "Number of blob transactions to send per slot")
	flags.Uint64VarP(&s.options.Sidecars, "sidecars", "b", 3, "Maximum number of blob sidecars per blob transactions")
	flags.Uint64Var(&s.options.MaxPending, "max-pending", 0, "Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", 0, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.Replace, "replace", 30, "Number of seconds to wait before replace a transaction")
	flags.Uint64Var(&s.options.MaxReplacements, "max-replace", 4, "Maximum number of replacement transactions")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", 30, "Number of seconds to wait before re-broadcasting a transaction")
	flags.Uint64Var(&s.options.BaseFee, "basefee", 20, "Max fee per gas to use in blob transactions (in gwei)")
	flags.Uint64Var(&s.options.TipFee, "tipfee", 2, "Max tip per gas to use in blob transactions (in gwei)")
	flags.Uint64Var(&s.options.BlobFee, "blobfee", 20, "Max blob fee to use in blob transactions (in gwei)")
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

	s.logger.Infof("starting scenario: blob-combined")

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
			tx, client, wallet, err := s.sendBlobTx(txIdx, 0, 0)
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
				logger.Warnf("blob tx %6d.0 failed: %v", txIdx+1, err)
				if s.pendingChan != nil {
					<-s.pendingChan
				}
				return
			}

			txCount.Add(1)
			logger.Infof("blob tx %6d.0 sent:  %v (%v sidecars)", txIdx+1, tx.Hash().String(), len(tx.BlobTxSidecar().Blobs))
		}(txIdx, lastChan, currentChan)

		lastChan = currentChan

		count := txCount.Load() + uint64(pendingCount.Load())
		if s.options.TotalCount > 0 && count >= s.options.TotalCount {
			break
		}
	}
	<-lastChan
	close(lastChan)

	s.logger.Infof("finished sending transactions, awaiting block inclusion...")
	s.pendingWGroup.Wait()
	s.logger.Infof("all transactions included!")

	return nil
}

func (s *Scenario) sendBlobTx(txIdx uint64, replacementIdx uint64, txNonce uint64) (*types.Transaction, *txbuilder.Client, *txbuilder.Wallet, error) {
	client := s.tester.GetClient(tester.SelectByIndex, int(txIdx))
	wallet := s.tester.GetWallet(tester.SelectByIndex, int(txIdx))

	if rand.Intn(100) < 20 {
		// 20% chance to send transaction via another client
		// will cause some replacement txs being sent via different clients than the original tx
		client = s.tester.GetClient(tester.SelectRandom, 0)
	}

	var feeCap *big.Int
	var tipCap *big.Int
	var blobFee *big.Int

	if s.options.BaseFee > 0 {
		feeCap = new(big.Int).Mul(big.NewInt(int64(s.options.BaseFee)), big.NewInt(1000000000))
	}
	if s.options.TipFee > 0 {
		tipCap = new(big.Int).Mul(big.NewInt(int64(s.options.TipFee)), big.NewInt(1000000000))
	}
	if s.options.BlobFee > 0 {
		blobFee = new(big.Int).Mul(big.NewInt(int64(s.options.BlobFee)), big.NewInt(1000000000))
	}

	if feeCap == nil || tipCap == nil {
		// get suggested fee from client
		var err error
		feeCap, tipCap, err = client.GetSuggestedFee()
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
	if blobFee == nil {
		blobFee = big.NewInt(1000000000)
	}

	for i := 0; i < int(replacementIdx); i++ {
		// x3 fee for each replacement tx
		feeCap = feeCap.Mul(feeCap, big.NewInt(3))
		tipCap = tipCap.Mul(tipCap, big.NewInt(3))
		blobFee = blobFee.Mul(blobFee, big.NewInt(3))
	}

	blobCount := uint64(rand.Int63n(int64(s.options.Sidecars)) + 1)
	blobRefs := make([][]string, blobCount)
	for i := 0; i < int(blobCount); i++ {
		blobLabel := fmt.Sprintf("0x1611AA0000%08dFF%02dFF%04dFEED", txIdx, i, replacementIdx)

		specialBlob := rand.Intn(50)
		switch specialBlob {
		case 0: // special blob commitment - all 0x0
			blobRefs[i] = []string{"0x0"}
		case 1, 2: // reuse well known blob
			blobRefs[i] = []string{"repeat:0x42:1337"}
		case 3, 4: // duplicate commitment
			if i == 0 {
				blobRefs[i] = []string{blobLabel, "random"}
			} else {
				blobRefs[i] = []string{"copy:0"}
			}

		default: // random blob data
			blobRefs[i] = []string{blobLabel, "random"}
		}
	}

	toAddr := s.tester.GetWallet(tester.SelectByIndex, int(txIdx)+1).GetAddress()
	blobTx, err := txbuilder.BuildBlobTx(&txbuilder.TxMetadata{
		GasFeeCap:  uint256.MustFromBig(feeCap),
		GasTipCap:  uint256.MustFromBig(tipCap),
		BlobFeeCap: uint256.MustFromBig(blobFee),
		Gas:        21000,
		To:         &toAddr,
		Value:      uint256.NewInt(0),
	}, blobRefs)
	if err != nil {
		return nil, client, wallet, err
	}

	var tx *types.Transaction
	if replacementIdx == 0 {
		tx, err = wallet.BuildBlobTx(blobTx)
	} else {
		tx, err = wallet.ReplaceBlobTx(blobTx, txNonce)
	}
	if err != nil {
		return nil, client, wallet, err
	}

	rebroadcast := 0
	if s.options.Rebroadcast > 0 {
		rebroadcast = 10
	}

	var awaitConfirmation bool = true
	s.pendingWGroup.Add(1)
	err = s.tester.GetTxPool().SendTransaction(context.Background(), wallet, tx, &txbuilder.SendTransactionOptions{
		Client:              client,
		MaxRebroadcasts:     rebroadcast,
		RebroadcastInterval: time.Duration(s.options.Rebroadcast) * time.Second,
		OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			defer func() {
				awaitConfirmation = false
				if replacementIdx == 0 {
					if s.pendingChan != nil {
						time.Sleep(100 * time.Millisecond)
						<-s.pendingChan
					}
				}
				s.pendingWGroup.Done()
			}()

			if err != nil {
				s.logger.WithField("client", client.GetName()).Warnf("blob tx %6d.%v: await receipt failed: %v", txIdx+1, replacementIdx, err)
				return
			}
			if receipt == nil {
				return
			}

			effectiveGasPrice := receipt.EffectiveGasPrice
			if effectiveGasPrice == nil {
				effectiveGasPrice = big.NewInt(0)
			}
			blobGasPrice := receipt.BlobGasPrice
			if blobGasPrice == nil {
				blobGasPrice = big.NewInt(0)
			}
			feeAmount := new(big.Int).Mul(effectiveGasPrice, big.NewInt(int64(receipt.GasUsed)))
			blobFeeAmount := new(big.Int).Mul(blobGasPrice, big.NewInt(int64(receipt.BlobGasUsed)))
			totalAmount := new(big.Int).Add(tx.Value(), feeAmount)
			totalAmount = new(big.Int).Add(totalAmount, blobFeeAmount)
			wallet.SubBalance(totalAmount)

			gweiTotalFee := new(big.Int).Div(feeAmount, big.NewInt(1000000000))
			gweiBaseFee := new(big.Int).Div(effectiveGasPrice, big.NewInt(1000000000))
			gweiBlobFee := new(big.Int).Div(blobGasPrice, big.NewInt(1000000000))

			s.logger.WithField("client", client.GetName()).Debugf("blob tx %6d.%v confirmed in block #%v!  total fee: %v gwei (base: %v, blob: %v)", txIdx+1, replacementIdx, receipt.BlockNumber.String(), gweiTotalFee, gweiBaseFee, gweiBlobFee)
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
				logger.Debugf("failed sending blob tx %6d.%v: %v", txIdx+1, replacementIdx, err)
			} else if retry > 0 || rebroadcast > 0 {
				logger.Debugf("successfully sent blob tx %6d.%v", txIdx+1, replacementIdx)
			}
		},
	})
	if err != nil {
		if replacementIdx == 0 {
			// reset nonce if tx was not sent
			wallet.ResetPendingNonce(client)
		}

		return nil, client, wallet, err
	}

	if s.options.Replace > 0 && replacementIdx < s.options.MaxReplacements && rand.Intn(100) < 70 {
		go s.delayedReplace(txIdx, tx, &awaitConfirmation, replacementIdx)
	}

	return tx, client, wallet, nil
}

func (s *Scenario) delayedReplace(txIdx uint64, tx *types.Transaction, awaitConfirmation *bool, replacementIdx uint64) {
	time.Sleep(time.Duration(rand.Intn(int(s.options.Replace))+2) * time.Second)

	if !*awaitConfirmation {
		return
	}

	replaceTx, client, wallet, err := s.sendBlobTx(txIdx, replacementIdx+1, tx.Nonce())
	if err != nil {
		s.logger.WithField("client", client.GetName()).Warnf("blob tx %6d.%v replacement failed: %v", txIdx+1, replacementIdx+1, err)
		return
	}
	s.logger.WithField("client", client.GetName()).Infof("blob tx %6d.%v sent:  %v (%v sidecars), wallet: %v, nonce: %v", txIdx+1, replacementIdx+1, replaceTx.Hash().String(), len(tx.BlobTxSidecar().Blobs), s.tester.GetWalletIndex(wallet.GetAddress()), tx.Nonce())
}
