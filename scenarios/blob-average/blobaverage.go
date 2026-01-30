package blobaverage

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"time"

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
	TargetAverage   float64                  `yaml:"target_average"`
	TrackingSeconds uint64                   `yaml:"tracking_seconds"`
	MaxPending      uint64                   `yaml:"max_pending"`
	MaxWallets      uint64                   `yaml:"max_wallets"`
	Sidecars        uint64                   `yaml:"sidecars"`
	Rebroadcast     uint64                   `yaml:"rebroadcast"`
	BaseFee         float64                  `yaml:"base_fee"`
	TipFee          float64                  `yaml:"tip_fee"`
	BaseFeeWei      string                   `yaml:"base_fee_wei"`
	TipFeeWei       string                   `yaml:"tip_fee_wei"`
	BlobFee         float64                  `yaml:"blob_fee"`
	BlobV1Percent   uint64                   `yaml:"blob_v1_percent"`
	FuluActivation  utils.FlexibleJsonUInt64 `yaml:"fulu_activation"`
	BlobData        string                   `yaml:"blob_data"`
	ClientGroup     string                   `yaml:"client_group"`
	SubmitCount     uint64                   `yaml:"submit_count"`
	LogTxs          bool                     `yaml:"log_txs"`
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	// Blob tracking
	blobTracker      *BlobTracker
	lastBlockNumber  uint64
	pendingBlobCount atomic.Int64
	blockChan        chan uint64
}

var ScenarioName = "blob-average"
var ScenarioDefaultOptions = ScenarioOptions{
	TargetAverage:   3,
	TrackingSeconds: 3600, // 1 hour
	MaxPending:      10,
	MaxWallets:      20,
	Sidecars:        1,
	Rebroadcast:     1,
	BaseFee:         20,
	TipFee:          2,
	BlobFee:         20,
	BlobV1Percent:   100,
	FuluActivation:  math.MaxInt64,
	BlobData:        "",
	ClientGroup:     "",
	SubmitCount:     3,
	LogTxs:          false,
}
var ScenarioDescriptor = scenario.Descriptor{
	Name:           ScenarioName,
	Description:    "Send blob transactions to maintain a network-wide average blob count",
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
	flags.Float64VarP(&s.options.TargetAverage, "target-average", "a", ScenarioDefaultOptions.TargetAverage, "Target average blob count per block")
	flags.Uint64Var(&s.options.TrackingSeconds, "tracking-seconds", ScenarioDefaultOptions.TrackingSeconds, "Time window in seconds for tracking blob averages (default 1h)")
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum number of pending blob transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets to use")
	flags.Uint64VarP(&s.options.Sidecars, "sidecars", "b", ScenarioDefaultOptions.Sidecars, "Number of blob sidecars per blob transaction")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", ScenarioDefaultOptions.Rebroadcast, "Enable reliable rebroadcast system")
	flags.Float64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use in blob transactions (in gwei)")
	flags.Float64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use in blob transactions (in gwei)")
	flags.StringVar(&s.options.BaseFeeWei, "basefee-wei", "", "Max fee per gas in wei (overrides --basefee for L2 sub-gwei fees)")
	flags.StringVar(&s.options.TipFeeWei, "tipfee-wei", "", "Max tip per gas in wei (overrides --tipfee for L2 sub-gwei fees)")
	flags.Float64Var(&s.options.BlobFee, "blobfee", ScenarioDefaultOptions.BlobFee, "Max blob fee to use in blob transactions (in gwei)")
	flags.Uint64Var(&s.options.BlobV1Percent, "blob-v1-percent", ScenarioDefaultOptions.BlobV1Percent, "Percentage of blob transactions to be submitted with the v1 wrapper format")
	flags.Uint64Var((*uint64)(&s.options.FuluActivation), "fulu-activation", uint64(ScenarioDefaultOptions.FuluActivation), "Unix timestamp of the Fulu activation")
	flags.StringVar(&s.options.BlobData, "blob-data", ScenarioDefaultOptions.BlobData, "Blob data to use in blob transactions")
	flags.StringVar(&s.options.ClientGroup, "client-group", ScenarioDefaultOptions.ClientGroup, "Client group to use for sending transactions")
	flags.Uint64Var(&s.options.SubmitCount, "submit-count", ScenarioDefaultOptions.SubmitCount, "Number of times to submit each transaction (to increase chance of inclusion)")
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

	if options.GlobalCfg != nil {
		if v, ok := options.GlobalCfg["fulu_activation"]; ok && s.options.FuluActivation == ScenarioDefaultOptions.FuluActivation {
			s.options.FuluActivation = utils.FlexibleJsonUInt64(v.(uint64))
		}
	}

	if s.options.MaxWallets > 0 {
		s.walletPool.SetWalletCount(s.options.MaxWallets)
	} else {
		s.walletPool.SetWalletCount(20)
	}

	if s.options.Sidecars > 6 {
		s.logger.Warnf("Transactions with more than 6 blobs will most likely be dropped by the execution layer client. Got %d sidecars, limiting to 6.", s.options.Sidecars)
	}

	if s.options.TargetAverage <= 0 {
		return fmt.Errorf("target-average must be greater than 0")
	}

	if s.options.TrackingSeconds == 0 {
		return fmt.Errorf("tracking-seconds must be greater than 0")
	}

	return nil
}

func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Infof("starting scenario: %s", ScenarioName)
	s.logger.Infof("target average: %.2f blobs/block, tracking window: %d seconds", s.options.TargetAverage, s.options.TrackingSeconds)
	defer s.logger.Infof("scenario %s finished.", ScenarioName)

	// Initialize blob tracker
	s.blobTracker = NewBlobTracker(time.Duration(s.options.TrackingSeconds) * time.Second)
	s.blockChan = make(chan uint64, 1)

	// Subscribe to bulk block updates
	txPool := s.walletPool.GetTxPool()
	subscriptionID := txPool.SubscribeToBulkBlockUpdates(func(blockNumber uint64, globalBlockStats *spamoor.GlobalBlockStats) {
		s.processBlock(blockNumber, globalBlockStats)
	})
	defer txPool.UnsubscribeFromBulkBlockUpdates(subscriptionID)

	// Run the custom transaction loop
	return s.runBlobAverageLoop(ctx)
}

// processBlock processes a new block and updates the blob tracker
func (s *Scenario) processBlock(blockNumber uint64, globalBlockStats *spamoor.GlobalBlockStats) {
	if globalBlockStats == nil || globalBlockStats.Block == nil {
		return
	}

	block := globalBlockStats.Block
	blobCount := uint64(0)

	// Count blobs from all transactions in the block
	for _, tx := range block.Transactions() {
		if tx.Type() == types.BlobTxType {
			blobCount += uint64(len(tx.BlobHashes()))
		}
	}

	s.blobTracker.AddBlock(blockNumber, block.Time(), blobCount)
	s.lastBlockNumber = blockNumber

	avgBlobs := s.blobTracker.GetAverageBlobCount()
	deficit := s.blobTracker.GetBlobDeficit(s.options.TargetAverage)
	pendingBlobs := s.pendingBlobCount.Load()

	s.logger.Debugf("block %d: %d blobs, avg: %.2f, deficit: %.2f, pending: %d", blockNumber, blobCount, avgBlobs, deficit, pendingBlobs)

	select {
	case s.blockChan <- blockNumber:
	default:
	}
}

// runBlobAverageLoop runs the main loop that sends blob transactions to maintain the target average
func (s *Scenario) runBlobAverageLoop(ctx context.Context) error {
	slotDuration := scenario.GlobalSlotDuration

	pendingCount := atomic.Int64{}
	txIdxCounter := atomic.Uint64{}

	var pendingMutex sync.Mutex
	var pendingCond *sync.Cond
	pendingWg := sync.WaitGroup{}

	if s.options.MaxPending > 0 {
		pendingCond = sync.NewCond(&pendingMutex)
	}

outer:
	for {
		select {
		case <-ctx.Done():
			// Wait for pending transactions to complete before exiting
			break outer
		case <-s.blockChan:
		case <-time.After(slotDuration * 2): // Fallback timeout if block notification missed
		}

		// Check if we need to send more blobs
		deficit := s.blobTracker.GetBlobDeficit(s.options.TargetAverage)

		// Subtract pending blobs (blobs we've submitted but not yet confirmed)
		pendingBlobs := s.pendingBlobCount.Load()
		effectiveDeficit := deficit - float64(pendingBlobs)

		// Calculate how many blob transactions we need to send
		// Each transaction has s.options.Sidecars blobs
		blobsPerTx := s.options.Sidecars
		if blobsPerTx == 0 {
			blobsPerTx = 1
		}

		// Respect max pending limit
		currentPending := uint64(pendingCount.Load())
		if s.options.MaxPending > 0 && currentPending >= s.options.MaxPending {
			s.logger.Debugf("max pending reached (%d/%d), waiting", currentPending, s.options.MaxPending)
			continue
		}

		// Limit txs to send based on available pending slots
		if effectiveDeficit < 0 {
			effectiveDeficit = 0
		}
		txsNeeded := uint64(math.Ceil(effectiveDeficit / float64(blobsPerTx)))
		if s.options.MaxPending > 0 {
			availableSlots := s.options.MaxPending - currentPending
			if txsNeeded > availableSlots {
				txsNeeded = availableSlots
			}
		}

		s.logger.Infof("last block: %d, blob deficit: %.2f, pending: %d, sending: %d txs (%d blobs each)", s.lastBlockNumber, deficit, pendingBlobs, txsNeeded, blobsPerTx)

		// If we're above or at target (accounting for pending), don't send more blobs
		if effectiveDeficit <= 0 {
			avgBlobs := s.blobTracker.GetAverageBlobCount()
			s.logger.Debugf("average %.2f + %d pending >= target %.2f, not sending blobs", avgBlobs, pendingBlobs, s.options.TargetAverage)
			continue
		}

		if txsNeeded == 0 {
			continue
		}

		// Send the required number of blob transactions
		for i := uint64(0); i < txsNeeded; i++ {
			if ctx.Err() != nil {
				break
			}

			// Check pending limit
			if s.options.MaxPending > 0 {
				pendingMutex.Lock()
				for pendingCount.Load() >= int64(s.options.MaxPending) {
					pendingCond.Wait()
					if ctx.Err() != nil {
						pendingMutex.Unlock()
						break
					}
				}
				pendingMutex.Unlock()
			}

			if ctx.Err() != nil {
				break
			}

			pendingCount.Add(1)
			pendingWg.Add(1)

			txIdx := txIdxCounter.Add(1) - 1
			blobsInTx := int64(blobsPerTx)

			go func(txIdx uint64, blobsInTx int64) {
				defer func() {
					utils.RecoverPanic(s.logger, "blobaverage.sendBlobTx", nil)
					pendingWg.Done()
					pendingCount.Add(-1)
					if pendingCond != nil {
						pendingCond.Signal()
					}
				}()

				receiptChan, tx, client, wallet, txVersion, err := s.sendBlobTx(ctx, txIdx)
				logger := s.logger
				if client != nil {
					logger = logger.WithField("rpc", client.GetName())
				}
				if tx != nil {
					logger = logger.WithField("nonce", tx.Nonce())
				}
				if wallet != nil {
					logger = logger.WithField("wallet", s.walletPool.GetWalletName(wallet.GetAddress()))
				}

				if err != nil {
					logger.Warnf("blob tx %6d failed: %v", txIdx+1, err)
					return
				}

				// Track pending blobs
				s.pendingBlobCount.Add(blobsInTx)

				if s.options.LogTxs {
					logger.Infof("blob tx %6d sent: %v (%v sidecars, v%v)", txIdx+1, tx.Hash().String(), len(tx.BlobTxSidecar().Blobs), txVersion)
				} else {
					logger.Debugf("blob tx %6d sent: %v (%v sidecars, v%v)", txIdx+1, tx.Hash().String(), len(tx.BlobTxSidecar().Blobs), txVersion)
				}

				// Wait for receipt
				if _, err := receiptChan.Wait(ctx); err != nil {
					logger.Debugf("blob tx %6d receipt wait error: %v", txIdx+1, err)
				}

				// Decrement pending blob count when transaction completes
				s.pendingBlobCount.Add(-blobsInTx)
			}(txIdx, blobsInTx)
		}
	}

	// Wait for pending transactions to complete
	pendingWg.Wait()

	return ctx.Err()
}

func (s *Scenario) sendBlobTx(ctx context.Context, txIdx uint64) (scenario.ReceiptChan, *types.Transaction, *spamoor.Client, *spamoor.Wallet, uint8, error) {
	client := s.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, int(txIdx)),
		spamoor.WithClientGroup(s.options.ClientGroup),
	)
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByPendingTxCount, int(txIdx))

	if client == nil {
		return nil, nil, client, wallet, 0, scenario.ErrNoClients
	}

	if wallet == nil {
		return nil, nil, client, wallet, 0, scenario.ErrNoWallet
	}

	if err := wallet.ResetNoncesIfNeeded(ctx, client); err != nil {
		return nil, nil, client, wallet, 0, err
	}

	baseFeeWei, tipFeeWei := spamoor.ResolveFees(s.options.BaseFee, s.options.TipFee, s.options.BaseFeeWei, s.options.TipFeeWei)
	feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(client, baseFeeWei, tipFeeWei)
	if err != nil {
		return nil, nil, client, wallet, 0, err
	}

	var blobFee *big.Int
	if s.options.BlobFee > 0 {
		blobFee = big.NewInt(int64(s.options.BlobFee * 1e9))
	} else {
		blobFee = new(big.Int).Mul(feeCap, big.NewInt(1000000000))
	}

	blobCount := s.options.Sidecars
	blobRefs := make([][]string, blobCount)
	for i := 0; i < int(blobCount); i++ {
		blobLabel := fmt.Sprintf("0x1611AA0000%08dFF%02dFF%04dFEED", txIdx, i, 0)

		if s.options.BlobData != "" {
			blobRefs[i] = []string{}
			for _, blob := range strings.Split(s.options.BlobData, ",") {
				if blob == "label" {
					blob = blobLabel
				}
				blobRefs[i] = append(blobRefs[i], blob)
			}

		} else {
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
				blobRefs[i] = []string{blobLabel, "random:full"}
			}
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
		return nil, nil, client, wallet, 0, err
	}

	tx, err := wallet.BuildBlobTx(blobTx)
	if err != nil {
		return nil, nil, client, wallet, 0, err
	}

	isBlobV1 := false
	receiptChan := make(scenario.ReceiptChan, 1)

	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		ClientGroup: s.options.ClientGroup,
		Rebroadcast: s.options.Rebroadcast > 0,
		SubmitCount: int(s.options.SubmitCount),
		OnComplete: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			receiptChan <- receipt
		},
		OnConfirm: func(tx *types.Transaction, receipt *types.Receipt) {
			txFees := utils.GetTransactionFees(tx, receipt)
			s.logger.WithField("rpc", client.GetName()).Debugf(
				" transaction %d confirmed in block #%v. total fee: %v gwei (tx: %v/%v, blob: %v/%v)",
				txIdx+1,
				receipt.BlockNumber.String(),
				txFees.TotalFeeGweiString(),
				txFees.TxFeeGweiString(),
				txFees.TxBaseFeeGweiString(),
				txFees.BlobFeeGweiString(),
				txFees.BlobBaseFeeGweiString(),
			)
		},
		LogFn: func(client *spamoor.Client, retry int, rebroadcast int, err error) {
			logger := s.logger.WithField("rpc", client.GetName()).WithField("nonce", tx.Nonce())
			if retry == 0 && rebroadcast > 0 {
				logger.Infof("rebroadcasting blob tx %6d", txIdx+1)
			}
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
		OnEncode: func(tx *types.Transaction) ([]byte, error) {
			sendAsV1 := uint64(time.Now().Unix()) > uint64(s.options.FuluActivation) && rand.Intn(100) < int(s.options.BlobV1Percent)
			if sendAsV1 && !isBlobV1 {
				err := tx.BlobTxSidecar().ToV1()
				if err != nil {
					return nil, err
				}

				isBlobV1 = true
			}
			return nil, nil
		},
	})
	if err != nil {
		wallet.MarkSkippedNonce(tx.Nonce())
		return nil, nil, client, wallet, 0, err
	}

	return receiptChan, tx, client, wallet, tx.BlobTxSidecar().Version, nil
}
