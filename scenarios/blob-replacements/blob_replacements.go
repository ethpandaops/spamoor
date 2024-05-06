package blobreplacements

import (
	"fmt"
	"math/big"
	"math/rand"
	"sync"
	"time"

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
		logger: logrus.WithField("scenario", "blob-replacements"),
	}
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
	flags.Uint64VarP(&s.options.TotalCount, "count", "c", 0, "Total number of blob transactions to send")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t", 0, "Number of blob transactions to send per slot")
	flags.Uint64VarP(&s.options.Sidecars, "sidecars", "b", 1, "Number of blob sidecars per blob transactions")
	flags.Uint64Var(&s.options.MaxPending, "max-pending", 0, "Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", 0, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.Replace, "replace", 10, "Number of seconds to wait before replace a transaction")
	flags.Uint64Var(&s.options.MaxReplacements, "max-replace", 5, "Maximum number of replacement transactions")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", 120, "Number of seconds to wait before re-broadcasting a transaction")
	flags.Uint64Var(&s.options.BaseFee, "basefee", 20, "Max fee per gas to use in blob transactions (in gwei)")
	flags.Uint64Var(&s.options.TipFee, "tipfee", 2, "Max tip per gas to use in blob transactions (in gwei)")
	flags.Uint64Var(&s.options.BlobFee, "blobfee", 20, "Max blob fee to use in blob transactions (in gwei)")
	return nil
}

func (s *Scenario) Init(testerCfg *tester.TesterConfig) error {
	if s.options.TotalCount == 0 && s.options.Throughput == 0 {
		return fmt.Errorf("neither total count nor throughput limit set, must define at least one of them")
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
	counterMutex := sync.Mutex{}
	waitGroup := sync.WaitGroup{}
	pendingCount := uint64(0)
	txCount := uint64(0)
	startTime := time.Now()

	s.logger.Infof("starting scenario: replacements")

	for {
		txIdx := txIdxCounter
		txIdxCounter++

		if s.pendingChan != nil {
			// await pending transactions
			s.pendingChan <- true
		}
		waitGroup.Add(1)
		counterMutex.Lock()
		pendingCount++
		counterMutex.Unlock()

		go func(txIdx uint64) {
			defer func() {
				counterMutex.Lock()
				pendingCount--
				counterMutex.Unlock()
				waitGroup.Done()
			}()

			logger := s.logger
			tx, client, err := s.sendBlobTx(txIdx, 0, 0)
			if client != nil {
				logger = logger.WithField("rpc", client.GetName())
			}
			if err != nil {
				logger.Warnf("blob tx %6d.0 failed: %v", txIdx+1, err)
				if s.pendingChan != nil {
					<-s.pendingChan
				}
				return
			}

			counterMutex.Lock()
			txCount++
			counterMutex.Unlock()
			logger.Infof("blob tx %6d.0 sent:  %v (%v sidecars)", txIdx+1, tx.Hash().String(), len(tx.BlobTxSidecar().Blobs))
		}(txIdx)

		count := txCount + pendingCount
		if s.options.TotalCount > 0 && count >= s.options.TotalCount {
			break
		}
		if s.options.Throughput > 0 {
			for count/((uint64(time.Since(startTime).Seconds())/utils.SecondsPerSlot)+1) >= s.options.Throughput {
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
	waitGroup.Wait()

	s.logger.Infof("finished sending transactions, awaiting block inclusion...")
	s.pendingWGroup.Wait()
	s.logger.Infof("all transactions included!")

	return nil
}

func (s *Scenario) sendBlobTx(txIdx uint64, replacementIdx uint64, txNonce uint64) (*types.Transaction, *txbuilder.Client, error) {
	client := s.tester.GetClient(tester.SelectByIndex, int(txIdx))
	wallet := s.tester.GetWallet(tester.SelectByIndex, int(txIdx))

	if rand.Intn(100) < 50 {
		// 50% chance to send transaction via another client
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
			return nil, client, err
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
		return nil, client, err
	}

	var tx *types.Transaction
	if replacementIdx == 0 {
		tx, err = wallet.BuildBlobTx(blobTx)
	} else {
		tx, err = wallet.ReplaceBlobTx(blobTx, txNonce)
	}
	if err != nil {
		return nil, client, err
	}

	err = client.SendTransaction(tx)
	if err != nil {
		return nil, client, err
	}

	s.pendingWGroup.Add(1)
	go s.awaitTx(txIdx, tx, client, wallet, replacementIdx)

	return tx, client, nil
}

func (s *Scenario) awaitTx(txIdx uint64, tx *types.Transaction, client *txbuilder.Client, wallet *txbuilder.Wallet, replacementIdx uint64) {
	var awaitConfirmation bool = true
	defer func() {
		awaitConfirmation = false
		if replacementIdx == 0 {
			if s.pendingChan != nil {
				<-s.pendingChan
			}
		}
		s.pendingWGroup.Done()
	}()
	if s.options.Replace > 0 && replacementIdx < s.options.MaxReplacements && rand.Intn(100) < 70 {
		go s.delayedReplace(txIdx, tx, &awaitConfirmation, replacementIdx)
	} else if s.options.Rebroadcast > 0 {
		go s.delayedResend(txIdx, tx, &awaitConfirmation, replacementIdx)
	}

	receipt, blockNum, err := client.AwaitTransaction(tx)
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
	totalAmount := new(big.Int).Add(tx.Value(), feeAmount)
	wallet.SubBalance(totalAmount)

	gweiTotalFee := new(big.Int).Div(totalAmount, big.NewInt(1000000000))
	gweiBaseFee := new(big.Int).Div(effectiveGasPrice, big.NewInt(1000000000))
	gweiBlobFee := new(big.Int).Div(blobGasPrice, big.NewInt(1000000000))

	s.logger.WithField("client", client.GetName()).Infof("blob tx %6d.%v confirmed in block #%v!  total fee: %v gwei (base: %v, blob: %v)", txIdx+1, replacementIdx, blockNum, gweiTotalFee, gweiBaseFee, gweiBlobFee)
}

func (s *Scenario) delayedResend(txIdx uint64, tx *types.Transaction, awaitConfirmation *bool, replacementIdx uint64) {
	for {
		time.Sleep(time.Duration(s.options.Rebroadcast) * time.Second)

		if !*awaitConfirmation {
			break
		}

		client := s.tester.GetClient(tester.SelectRandom, 0)
		client.SendTransaction(tx)
		s.logger.WithField("client", client.GetName()).Debugf("blob tx %6d.%v re-broadcasted.", txIdx+1, replacementIdx)
	}
}

func (s *Scenario) delayedReplace(txIdx uint64, tx *types.Transaction, awaitConfirmation *bool, replacementIdx uint64) {
	time.Sleep(time.Duration(rand.Intn(int(s.options.Replace))+2) * time.Second)

	if !*awaitConfirmation {
		return
	}

	replaceTx, client, err := s.sendBlobTx(txIdx, replacementIdx+1, tx.Nonce())
	if err != nil {
		s.logger.WithField("client", client.GetName()).Warnf("blob tx %6d.%v replacement failed: %v", txIdx+1, replacementIdx+1, err)
		s.delayedResend(txIdx, tx, awaitConfirmation, replacementIdx)
		return
	}
	s.logger.WithField("client", client.GetName()).Infof("blob tx %6d.%v sent:  %v (%v sidecars)", txIdx+1, replacementIdx+1, replaceTx.Hash().String(), len(tx.BlobTxSidecar().Blobs))
}
