package blobs

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
	TotalCount  uint64 `yaml:"total_count"`
	Throughput  uint64 `yaml:"throughput"`
	Sidecars    uint64 `yaml:"sidecars"`
	MaxPending  uint64 `yaml:"max_pending"`
	MaxWallets  uint64 `yaml:"max_wallets"`
	Rebroadcast uint64 `yaml:"rebroadcast"`
	BaseFee     uint64 `yaml:"base_fee"`
	TipFee      uint64 `yaml:"tip_fee"`
	BlobFee     uint64 `yaml:"blob_fee"`
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	pendingChan   chan bool
	pendingWGroup sync.WaitGroup
}

var ScenarioName = "blobs"
var ScenarioDefaultOptions = ScenarioOptions{
	TotalCount:  0,
	Throughput:  0,
	Sidecars:    1,
	MaxPending:  0,
	MaxWallets:  0,
	Rebroadcast: 120,
	BaseFee:     20,
	TipFee:      2,
	BlobFee:     20,
}
var ScenarioDescriptor = scenariotypes.ScenarioDescriptor{
	Name:           ScenarioName,
	Description:    "Send simple blob transactions",
	DefaultOptions: ScenarioDefaultOptions,
	NewScenario:    newScenario,
}

func newScenario(logger logrus.FieldLogger) scenariotypes.Scenario {
	return &Scenario{
		logger: logger.WithField("scenario", ScenarioName),
	}
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
	flags.Uint64VarP(&s.options.TotalCount, "count", "c", ScenarioDefaultOptions.TotalCount, "Total number of blob transactions to send")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t", ScenarioDefaultOptions.Throughput, "Number of blob transactions to send per slot")
	flags.Uint64VarP(&s.options.Sidecars, "sidecars", "b", ScenarioDefaultOptions.Sidecars, "Number of blob sidecars per blob transactions")
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", ScenarioDefaultOptions.Rebroadcast, "Number of seconds to wait before re-broadcasting a transaction")
	flags.Uint64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use in blob transactions (in gwei)")
	flags.Uint64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use in blob transactions (in gwei)")
	flags.Uint64Var(&s.options.BlobFee, "blobfee", ScenarioDefaultOptions.BlobFee, "Max blob fee to use in blob transactions (in gwei)")
	return nil
}

func (s *Scenario) Init(walletPool *spamoor.WalletPool, config string) error {
	s.walletPool = walletPool

	if config != "" {
		err := yaml.Unmarshal([]byte(config), &s.options)
		if err != nil {
			return fmt.Errorf("failed to unmarshal config: %w", err)
		}
	}

	if s.options.MaxWallets > 0 {
		walletPool.SetWalletCount(s.options.MaxWallets)
	} else if s.options.TotalCount > 0 {
		if s.options.TotalCount < 1000 {
			walletPool.SetWalletCount(s.options.TotalCount)
		} else {
			walletPool.SetWalletCount(1000)
		}
	} else {
		if s.options.Throughput*10 < 1000 {
			walletPool.SetWalletCount(s.options.Throughput * 10)
		} else {
			walletPool.SetWalletCount(1000)
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
			tx, client, wallet, err := s.sendBlobTx(ctx, txIdx)
			if client != nil {
				logger = logger.WithField("rpc", client.GetName())
			}
			if tx != nil {
				logger = logger.WithField("nonce", tx.Nonce())
			}
			if wallet != nil {
				logger = logger.WithField("wallet", s.walletPool.GetWalletIndex(wallet.GetAddress()))
			}
			if lastChan != nil {
				<-lastChan
				close(lastChan)
			}
			if err != nil {
				logger.Warnf("could not send blob transaction: %v", err)
				<-s.pendingChan
				return
			}

			txCount.Add(1)
			logger.Infof("sent blob tx #%6d: %v (%v sidecars)", txIdx+1, tx.Hash().String(), len(tx.BlobTxSidecar().Blobs))
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

func (s *Scenario) sendBlobTx(ctx context.Context, txIdx uint64) (*types.Transaction, *txbuilder.Client, *txbuilder.Wallet, error) {
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, int(txIdx))
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(txIdx))

	var feeCap *big.Int
	var tipCap *big.Int
	var blobFee *big.Int

	if s.options.BaseFee > 0 {
		feeCap = new(big.Int).Mul(big.NewInt(int64(s.options.BaseFee)), big.NewInt(1000000))
	}
	if s.options.TipFee > 0 {
		tipCap = new(big.Int).Mul(big.NewInt(int64(s.options.TipFee)), big.NewInt(1000000))
	}
	if s.options.BlobFee > 0 {
		blobFee = new(big.Int).Mul(big.NewInt(int64(s.options.BlobFee)), big.NewInt(1000000))
	}

	if feeCap == nil || tipCap == nil {
		var err error
		feeCap, tipCap, err = client.GetSuggestedFee(s.walletPool.GetContext())
		if err != nil {
			return nil, client, wallet, err
		}
	}

	if feeCap.Cmp(big.NewInt(1000000)) < 0 {
		feeCap = big.NewInt(1000000)
	}
	if tipCap.Cmp(big.NewInt(1000000)) < 0 {
		tipCap = big.NewInt(1000000)
	}
	if blobFee == nil {
		blobFee = big.NewInt(1000000)
	}

	blobCount := uint64(rand.Int63n(int64(s.options.Sidecars)) + 1)
	blobRefs := make([][]string, blobCount)
	for i := 0; i < int(blobCount); i++ {
		blobLabel := fmt.Sprintf("0x1611AA0000%08dFF%02dFF%04dFEED", txIdx, i, 0)

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
		return nil, client, wallet, err
	}

	tx, err := wallet.BuildBlobTx(blobTx)
	if err != nil {
		return nil, nil, wallet, err
	}

	rebroadcast := 0
	if s.options.Rebroadcast > 0 {
		rebroadcast = 10
	}

	s.pendingWGroup.Add(1)
	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &txbuilder.SendTransactionOptions{
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
				s.logger.WithField("client", client.GetName()).Warnf("blob tx %6d: await receipt failed: %v", txIdx+1, err)
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

			gweiTotalFee := new(big.Int).Div(feeAmount, big.NewInt(1000000000))
			gweiBaseFee := new(big.Int).Div(effectiveGasPrice, big.NewInt(1000000000))
			gweiBlobFee := new(big.Int).Div(blobGasPrice, big.NewInt(1000000000))

			s.logger.WithField("client", client.GetName()).Debugf(" transaction %d confirmed in block #%v. total fee: %v gwei (base: %v, blob: %v)", txIdx+1, receipt.BlockNumber.String(), gweiTotalFee, gweiBaseFee, gweiBlobFee)
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
				logger.Debugf("failed sending blob tx %6d: %v", txIdx+1, err)
			} else if retry > 0 || rebroadcast > 0 {
				logger.Debugf("successfully sent blob tx %6d", txIdx+1)
			}
		},
	})
	if err != nil {
		// reset nonce if tx was not sent
		wallet.ResetPendingNonce(s.walletPool.GetContext(), client)

		return nil, client, wallet, err
	}

	return tx, client, wallet, nil
}
