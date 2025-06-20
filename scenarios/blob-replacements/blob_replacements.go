package blobreplacements

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenario"
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
	Timeout                     string                   `yaml:"timeout"`
	ClientGroup                 string                   `yaml:"client_group"`
	LogTxs                      bool                     `yaml:"log_txs"`
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool
}

var ScenarioName = "blob-replacements"
var ScenarioDefaultOptions = ScenarioOptions{
	TotalCount:                  0,
	Throughput:                  6,
	Sidecars:                    1,
	MaxPending:                  0,
	MaxWallets:                  0,
	Replace:                     10,
	MaxReplacements:             5,
	Rebroadcast:                 1,
	BaseFee:                     20,
	TipFee:                      2,
	BlobFee:                     20,
	BlobV1Percent:               100,
	FuluActivation:              math.MaxInt64,
	ThroughputIncrementInterval: 0,
	Timeout:                     "",
	ClientGroup:                 "",
	LogTxs:                      false,
}
var ScenarioDescriptor = scenario.Descriptor{
	Name:           ScenarioName,
	Description:    "Send blob transactions with replacements",
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
	flags.Uint64VarP(&s.options.TotalCount, "count", "c", ScenarioDefaultOptions.TotalCount, "Total number of blob transactions to send")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t", ScenarioDefaultOptions.Throughput, "Number of blob transactions to send per slot")
	flags.Uint64VarP(&s.options.Sidecars, "sidecars", "b", ScenarioDefaultOptions.Sidecars, "Number of blob sidecars per blob transactions")
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.Replace, "replace", ScenarioDefaultOptions.Replace, "Number of seconds to wait before replace a transaction")
	flags.Uint64Var(&s.options.MaxReplacements, "max-replace", ScenarioDefaultOptions.MaxReplacements, "Maximum number of replacement transactions")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", ScenarioDefaultOptions.Rebroadcast, "Enable reliable rebroadcast system")
	flags.Uint64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use in blob transactions (in gwei)")
	flags.Uint64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use in blob transactions (in gwei)")
	flags.Uint64Var(&s.options.BlobFee, "blobfee", ScenarioDefaultOptions.BlobFee, "Max blob fee to use in blob transactions (in gwei)")
	flags.Uint64Var(&s.options.BlobV1Percent, "blob-v1-percent", ScenarioDefaultOptions.BlobV1Percent, "Percentage of blob transactions to be submitted with the v1 wrapper format")
	flags.Uint64Var((*uint64)(&s.options.FuluActivation), "fulu-activation", uint64(ScenarioDefaultOptions.FuluActivation), "Unix timestamp of the Fulu activation")
	flags.Uint64Var(&s.options.ThroughputIncrementInterval, "throughput-increment-interval", ScenarioDefaultOptions.ThroughputIncrementInterval, "Increment the throughput every interval (in sec).")
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

	if options.GlobalCfg != nil {
		if v, ok := options.GlobalCfg["fulu_activation"]; ok && s.options.FuluActivation == ScenarioDefaultOptions.FuluActivation {
			s.options.FuluActivation = utils.FlexibleJsonUInt64(v.(uint64))
		}
	}

	if s.options.MaxWallets > 0 {
		s.walletPool.SetWalletCount(s.options.MaxWallets)
	} else if s.options.TotalCount > 0 {
		maxWallets := s.options.TotalCount / 3
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

	return nil
}

func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Infof("starting scenario: %s", ScenarioName)
	defer s.logger.Infof("scenario %s finished.", ScenarioName)

	maxPending := s.options.MaxPending
	if maxPending == 0 {
		maxPending = s.options.Throughput * 3
		if maxPending == 0 {
			maxPending = 1000
		}

		if maxPending > s.walletPool.GetConfiguredWalletCount()*2 {
			maxPending = s.walletPool.GetConfiguredWalletCount() * 2
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
		ThroughputIncrementInterval: s.options.ThroughputIncrementInterval,
		Timeout:                     timeout,
		WalletPool:                  s.walletPool,

		Logger: s.logger,
		ProcessNextTxFn: func(ctx context.Context, txIdx uint64, onComplete func()) (func(), error) {
			logger := s.logger
			tx, client, wallet, txVersion, err := s.sendBlobTx(ctx, txIdx, 0, 0, onComplete)
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
					logger.Warnf("blob tx %6d.0 failed: %v", txIdx+1, err)
				} else if s.options.LogTxs {
					logger.Infof("blob tx %6d.0 sent:  %v (%v sidecars, v%v)", txIdx+1, tx.Hash().String(), len(tx.BlobTxSidecar().Blobs), txVersion)
				} else {
					logger.Debugf("blob tx %6d.0 sent:  %v (%v sidecars, v%v)", txIdx+1, tx.Hash().String(), len(tx.BlobTxSidecar().Blobs), txVersion)
				}
			}, err
		},
	})

	return err
}

func (s *Scenario) sendBlobTx(ctx context.Context, txIdx uint64, replacementIdx uint64, txNonce uint64, onComplete func()) (*types.Transaction, *spamoor.Client, *spamoor.Wallet, uint8, error) {
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, int(txIdx), s.options.ClientGroup)
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByPendingTxCount, int(txIdx))
	transactionSubmitted := false

	defer func() {
		if !transactionSubmitted {
			onComplete()
		}
	}()

	if rand.Intn(100) < 50 {
		// 50% chance to send transaction via another client
		// will cause some replacement txs being sent via different clients than the original tx
		client = s.walletPool.GetClient(spamoor.SelectClientRandom, 0, s.options.ClientGroup)
	}

	if client == nil {
		return nil, client, wallet, 0, fmt.Errorf("no client available")
	}

	feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(client, s.options.BaseFee, s.options.TipFee)
	if err != nil {
		return nil, client, wallet, 0, err
	}

	var blobFee *big.Int
	if s.options.BlobFee > 0 {
		blobFee = new(big.Int).Mul(big.NewInt(int64(s.options.BlobFee)), big.NewInt(1000000000))
	} else {
		blobFee = new(big.Int).Mul(feeCap, big.NewInt(1000000000))
	}

	for i := 0; i < int(replacementIdx); i++ {
		// x3 fee for each replacement tx
		feeCap = feeCap.Mul(feeCap, big.NewInt(3))
		tipCap = tipCap.Mul(tipCap, big.NewInt(3))
		blobFee = blobFee.Mul(blobFee, big.NewInt(3))
	}

	blobCount := s.options.Sidecars
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

	var blobCellProofs []kzg4844.Proof

	if s.options.BlobV1Percent > 0 {
		// generate cell proofs here to avoid heavy recomputation on each submission
		blobCellProofs, err = txbuilder.GenerateCellProofs(tx.BlobTxSidecar())
		if err != nil {
			s.logger.Warnf("failed to generate cell proofs: %v", err)
		}
	}

	getTxBytes := func() ([]byte, uint8) {
		var txBytes []byte
		txVersion := uint8(0)
		sendAsV1 := time.Now().Unix() > int64(s.options.FuluActivation) && rand.Intn(100) < int(s.options.BlobV1Percent)
		if sendAsV1 {
			txBytes, err = txbuilder.MarshalBlobV1Tx(tx, blobCellProofs)
			if err != nil {
				s.logger.Warnf("failed to marshal blob tx as v1: %v", err)
			} else {
				txVersion = 1
			}
		}
		return txBytes, txVersion
	}

	_, txVersion := getTxBytes()

	var awaitConfirmation bool = true
	transactionSubmitted = true
	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		Rebroadcast: s.options.Rebroadcast > 0,
		OnComplete: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			awaitConfirmation = false
			onComplete()
		},
		OnConfirm: func(tx *types.Transaction, receipt *types.Receipt) {
			txFees := utils.GetTransactionFees(tx, receipt)
			s.logger.WithField("rpc", client.GetName()).Debugf(
				"blob tx %6d.%v confirmed in block #%v!  total fee: %v gwei (tx: %v/%v, blob: %v/%v)",
				txIdx+1,
				replacementIdx,
				receipt.BlockNumber.String(),
				txFees.TotalFeeGweiString(),
				txFees.TxFeeGweiString(),
				txFees.TxBaseFeeGweiString(),
				txFees.BlobFeeGweiString(),
				txFees.BlobBaseFeeGweiString(),
			)
		},
		LogFn: spamoor.GetDefaultLogFn(s.logger, "blob", fmt.Sprintf("%6d.%v", txIdx+1, replacementIdx), tx),
		OnEncode: func(tx *types.Transaction) ([]byte, error) {
			txBytes, _ := getTxBytes()
			return txBytes, nil
		},
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

	logger := s.logger

	replaceTx, client, wallet, txVersion, err := s.sendBlobTx(ctx, txIdx, replacementIdx+1, tx.Nonce(), func() {})
	if tx != nil {
		logger = logger.WithField("nonce", tx.Nonce())
	}
	if wallet != nil {
		logger = logger.WithField("wallet", s.walletPool.GetWalletName(wallet.GetAddress()))
	}
	if err != nil {
		logger.WithField("rpc", client.GetName()).Warnf("blob tx %6d.%v replacement failed: %v", txIdx+1, replacementIdx+1, err)
		return
	}

	if s.options.LogTxs {
		logger.WithField("rpc", client.GetName()).Infof("blob tx %6d.%v sent:  %v (%v sidecars, v%v)", txIdx+1, replacementIdx+1, replaceTx.Hash().String(), len(tx.BlobTxSidecar().Blobs), txVersion)
	} else {
		logger.WithField("rpc", client.GetName()).Debugf("blob tx %6d.%v sent:  %v (%v sidecars, v%v)", txIdx+1, replacementIdx+1, replaceTx.Hash().String(), len(tx.BlobTxSidecar().Blobs), txVersion)
	}
}
