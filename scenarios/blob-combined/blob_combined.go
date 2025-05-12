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
	"gopkg.in/yaml.v3"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenariotypes"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/utils"
)

type ScenarioOptions struct {
	TotalCount                  uint64                   `yaml:"total_count"`
	Throughput                  uint64                   `yaml:"throughput"`
	Sidecars                    uint64                   `yaml:"sidecars"`
	MaxPending                  uint64                   `yaml:"max_pending"`
	MaxWallets                  uint64                   `yaml:"max_wallets"`
	Replace                     uint64                   `yaml:"replace"`
	MaxReplacements             uint64                   `yaml:"max_replacements"`
	Rebroadcast                 uint64                   `yaml:"rebroadcast"`
	BaseFee                     uint64                   `yaml:"base_fee"`
	TipFee                      uint64                   `yaml:"tip_fee"`
	BlobFee                     uint64                   `yaml:"blob_fee"`
	BlobV1Percent               uint64                   `yaml:"blob_v1_percent"`
	FuluActivation              utils.FlexibleJsonUInt64 `yaml:"fulu_activation"`
	ThroughputIncrementInterval uint64                   `yaml:"throughput_increment_interval"`
	ClientGroup                 string                   `yaml:"client_group"`
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	pendingChan   chan bool
	pendingWGroup sync.WaitGroup
}

var ScenarioName = "blob-combined"
var ScenarioDefaultOptions = ScenarioOptions{
	TotalCount:                  0,
	Throughput:                  0,
	Sidecars:                    3,
	MaxPending:                  0,
	MaxWallets:                  0,
	Replace:                     30,
	MaxReplacements:             4,
	Rebroadcast:                 30,
	BaseFee:                     20,
	TipFee:                      2,
	BlobFee:                     20,
	BlobV1Percent:               100,
	FuluActivation:              0,
	ThroughputIncrementInterval: 0,
	ClientGroup:                 "",
}
var ScenarioDescriptor = scenariotypes.ScenarioDescriptor{
	Name:           ScenarioName,
	Description:    "Send blob transactions with different configurations",
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
	flags.Uint64VarP(&s.options.TotalCount, "count", "c", ScenarioDefaultOptions.TotalCount, "Total number of blob transactions to send")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t", ScenarioDefaultOptions.Throughput, "Number of blob transactions to send per slot")
	flags.Uint64VarP(&s.options.Sidecars, "sidecars", "b", ScenarioDefaultOptions.Sidecars, "Maximum number of blob sidecars per blob transactions")
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.Replace, "replace", ScenarioDefaultOptions.Replace, "Number of seconds to wait before replace a transaction")
	flags.Uint64Var(&s.options.MaxReplacements, "max-replace", ScenarioDefaultOptions.MaxReplacements, "Maximum number of replacement transactions")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", ScenarioDefaultOptions.Rebroadcast, "Number of seconds to wait before re-broadcasting a transaction")
	flags.Uint64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use in blob transactions (in gwei)")
	flags.Uint64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use in blob transactions (in gwei)")
	flags.Uint64Var(&s.options.BlobFee, "blobfee", ScenarioDefaultOptions.BlobFee, "Max blob fee to use in blob transactions (in gwei)")
	flags.Uint64Var(&s.options.BlobV1Percent, "blob-v1-percent", ScenarioDefaultOptions.BlobV1Percent, "Percentage of blob transactions to be submitted with the v1 wrapper format")
	flags.Uint64Var((*uint64)(&s.options.FuluActivation), "fulu-activation", uint64(ScenarioDefaultOptions.FuluActivation), "Unix timestamp of the Fulu activation")
	flags.Uint64Var(&s.options.ThroughputIncrementInterval, "throughput-increment-interval", ScenarioDefaultOptions.ThroughputIncrementInterval, "Increment the throughput every interval (in sec).")
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

	if options.GlobalCfg != nil {
		if v, ok := options.GlobalCfg["fulu_activation"]; ok && s.options.FuluActivation == 0 {
			s.options.FuluActivation = utils.FlexibleJsonUInt64(v.(uint64))
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

	if s.options.MaxPending > 0 {
		s.pendingChan = make(chan bool, s.options.MaxPending)
	}

	return nil
}

func (s *Scenario) Config() string {
	yamlBytes, _ := yaml.Marshal(&s.options)
	return string(yamlBytes)
}

func (s *Scenario) Run(ctx context.Context) error {
	txIdxCounter := uint64(0)
	pendingCount := atomic.Int64{}
	txCount := atomic.Uint64{}
	var lastChan chan bool

	s.logger.Infof("starting scenario: %s", ScenarioName)
	defer s.logger.Infof("scenario %s finished.", ScenarioName)

	initialRate := rate.Limit(float64(s.options.Throughput) / float64(utils.SecondsPerSlot))
	if initialRate == 0 {
		initialRate = rate.Inf
	}
	limiter := rate.NewLimiter(initialRate, 1)

	if s.options.ThroughputIncrementInterval != 0 {
		go func() {
			ticker := time.NewTicker(time.Duration(s.options.ThroughputIncrementInterval) * time.Second)
			for {
				select {
				case <-ticker.C:
					throughput := limiter.Limit() * 12
					newThroughput := throughput + 1
					s.logger.Infof("Increasing throughput from %.3f to %.3f", throughput, newThroughput)
					limiter.SetLimit(rate.Limit(float64(newThroughput) / float64(utils.SecondsPerSlot)))
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	for {
		if err := limiter.Wait(ctx); err != nil {
			if ctx.Err() != nil {
				break
			}

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
			tx, client, wallet, txVersion, err := s.sendBlobTx(ctx, txIdx, 0, 0, func() {
				if s.pendingChan != nil {
					time.Sleep(100 * time.Millisecond)
					<-s.pendingChan
				}
			})
			if client != nil {
				logger = logger.WithField("rpc", client.GetName())
			}
			if tx != nil {
				logger = logger.WithField("nonce", tx.Nonce())
			}
			if wallet != nil {
				logger = logger.WithField("wallet", s.walletPool.GetWalletName(wallet.GetAddress()))
			}
			if lastChan != nil {
				<-lastChan
				close(lastChan)
			}
			if err != nil {
				logger.Warnf("blob tx %6d.0 failed: %v", txIdx+1, err)
				return
			}

			txCount.Add(1)
			logger.Infof("blob tx %6d.0 sent:  %v (%v sidecars, v%v)", txIdx+1, tx.Hash().String(), len(tx.BlobTxSidecar().Blobs), txVersion)
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

	return nil
}

func (s *Scenario) sendBlobTx(ctx context.Context, txIdx uint64, replacementIdx uint64, txNonce uint64, onComplete func()) (*types.Transaction, *txbuilder.Client, *txbuilder.Wallet, uint8, error) {
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, int(txIdx), s.options.ClientGroup)
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(txIdx))
	transactionSubmitted := false

	defer func() {
		if !transactionSubmitted {
			onComplete()
		}
	}()

	if rand.Intn(100) < 20 {
		// 20% chance to send transaction via another client
		// will cause some replacement txs being sent via different clients than the original tx
		client = s.walletPool.GetClient(spamoor.SelectClientRandom, 0, s.options.ClientGroup)
	}

	if client == nil {
		return nil, client, wallet, 0, fmt.Errorf("no client available")
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
		feeCap, tipCap, err = client.GetSuggestedFee(s.walletPool.GetContext())
		if err != nil {
			return nil, client, wallet, 0, err
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

	toAddr := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(txIdx)+1).GetAddress()
	blobTx, err := txbuilder.BuildBlobTx(&txbuilder.TxMetadata{
		GasFeeCap:  uint256.MustFromBig(feeCap),
		GasTipCap:  uint256.MustFromBig(tipCap),
		BlobFeeCap: uint256.MustFromBig(blobFee),
		Gas:        21000,
		To:         &toAddr,
		Value:      uint256.NewInt(0),
	}, blobRefs)
	if err != nil {
		return nil, client, wallet, 0, err
	}

	var tx *types.Transaction
	if replacementIdx == 0 {
		tx, err = wallet.BuildBlobTx(blobTx)
	} else {
		tx, err = wallet.ReplaceBlobTx(blobTx, txNonce)
	}
	if err != nil {
		return nil, client, wallet, 0, err
	}

	var txBytes []byte
	txVersion := uint8(0)
	sendAsV1 := time.Now().Unix() > int64(s.options.FuluActivation) && rand.Intn(100) < int(s.options.BlobV1Percent)
	if sendAsV1 {
		txBytes, err = txbuilder.MarshalBlobV1Tx(tx)
		if err != nil {
			s.logger.Warnf("failed to marshal blob tx as v1: %v", err)
		} else {
			txVersion = 1
		}
	}

	rebroadcast := 0
	if s.options.Rebroadcast > 0 {
		rebroadcast = 10
	}

	var awaitConfirmation bool = true
	transactionSubmitted = true
	s.pendingWGroup.Add(1)
	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &txbuilder.SendTransactionOptions{
		Client:              client,
		MaxRebroadcasts:     rebroadcast,
		RebroadcastInterval: time.Duration(s.options.Rebroadcast) * time.Second,
		OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			defer func() {
				awaitConfirmation = false
				onComplete()
				s.pendingWGroup.Done()
			}()

			if err != nil {
				s.logger.WithField("rpc", client.GetName()).Warnf("blob tx %6d.%v: await receipt failed: %v", txIdx+1, replacementIdx, err)
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

			s.logger.WithField("rpc", client.GetName()).Debugf("blob tx %6d.%v confirmed in block #%v!  total fee: %v gwei (base: %v, blob: %v)", txIdx+1, replacementIdx, receipt.BlockNumber.String(), gweiTotalFee, gweiBaseFee, gweiBlobFee)
		},
		LogFn: func(client *txbuilder.Client, retry int, rebroadcast int, err error) {
			logger := s.logger.WithField("rpc", client.GetName())
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
		TransactionBytes: txBytes,
	})
	if err != nil {
		if replacementIdx == 0 {
			// reset nonce if tx was not sent
			wallet.ResetPendingNonce(s.walletPool.GetContext(), client)
		}

		return nil, client, wallet, 0, err
	}

	if s.options.Replace > 0 && replacementIdx < s.options.MaxReplacements && rand.Intn(100) < 70 {
		go s.delayedReplace(ctx, txIdx, tx, &awaitConfirmation, replacementIdx)
	}

	return tx, client, wallet, txVersion, nil
}

func (s *Scenario) delayedReplace(ctx context.Context, txIdx uint64, tx *types.Transaction, awaitConfirmation *bool, replacementIdx uint64) {
	time.Sleep(time.Duration(rand.Intn(int(s.options.Replace))+2) * time.Second)

	if !*awaitConfirmation {
		return
	}

	replaceTx, client, wallet, txVersion, err := s.sendBlobTx(ctx, txIdx, replacementIdx+1, tx.Nonce(), func() {})
	if err != nil {
		s.logger.WithField("rpc", client.GetName()).Warnf("blob tx %6d.%v replacement failed: %v", txIdx+1, replacementIdx+1, err)
		return
	}
	s.logger.WithFields(logrus.Fields{
		"rpc":    client.GetName(),
		"wallet": s.walletPool.GetWalletName(wallet.GetAddress()),
		"nonce":  tx.Nonce(),
	}).Infof("blob tx %6d.%v sent:  %v (%v sidecars, v%v)", txIdx+1, replacementIdx+1, replaceTx.Hash().String(), len(tx.BlobTxSidecar().Blobs), txVersion)
}
