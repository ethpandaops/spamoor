package frametx

import (
	"context"
	"crypto/rand"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/txtypes"
	"github.com/ethpandaops/spamoor/utils"
)

type ScenarioOptions struct {
	TotalCount  uint64  `yaml:"total_count"`
	Throughput  uint64  `yaml:"throughput"`
	MaxPending  uint64  `yaml:"max_pending"`
	MaxWallets  uint64  `yaml:"max_wallets"`
	Rebroadcast uint64  `yaml:"rebroadcast"`
	BaseFee     float64 `yaml:"base_fee"`
	TipFee      float64 `yaml:"tip_fee"`
	BaseFeeWei  string  `yaml:"base_fee_wei"`
	TipFeeWei   string  `yaml:"tip_fee_wei"`
	Timeout     string  `yaml:"timeout"`
	ClientGroup string  `yaml:"client_group"`
	LogTxs      bool    `yaml:"log_txs"`

	Shapes        string `yaml:"shapes"`
	FramesPerTx   uint64 `yaml:"frames_per_tx"`
	UserOpGas     uint64 `yaml:"user_op_gas"`
	VerifyGas     uint64 `yaml:"verify_gas"`
	StateGas      uint64 `yaml:"state_gas"`
	Amount        uint64 `yaml:"amount"`
	RandomAmount  bool   `yaml:"random_amount"`
	RandomTarget  bool   `yaml:"random_target"`
	Data          string `yaml:"data"`
	ExpiryOffset  uint64 `yaml:"expiry_offset"`
	VerifyFrames  bool   `yaml:"verify_frames"`
	SkipMempoolOK bool   `yaml:"skip_mempool_check"`
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool
	shapes     *weightedShapes
	callData   []byte

	// Conformance counters: how many confirmed transactions had frame receipts that
	// matched the shape's expectation, and how many did not.
	framesChecked  atomic.Uint64
	framesMismatch atomic.Uint64
	missingReceipt atomic.Uint64
}

var ScenarioName = "frametx"
var ScenarioDefaultOptions = ScenarioOptions{
	TotalCount:  0,
	Throughput:  10,
	MaxPending:  0,
	MaxWallets:  0,
	Rebroadcast: 1,
	BaseFee:     20,
	TipFee:      2,
	ClientGroup: "",
	LogTxs:      false,

	Shapes:       "all",
	FramesPerTx:  4,
	UserOpGas:    30000,
	VerifyGas:    5000,
	StateGas:     0,
	Amount:       20,
	RandomAmount: true,
	RandomTarget: false,
	ExpiryOffset: 600,
	VerifyFrames: true,
}

var ScenarioDescriptor = scenario.Descriptor{
	Name:           ScenarioName,
	Description:    "Send EIP-8141 frame transactions in a mix of frame shapes and check the per-frame receipts against what each shape should produce",
	DefaultOptions: ScenarioDefaultOptions,
	NewScenario:    newScenario,
}

func newScenario(logger logrus.FieldLogger) scenario.Scenario {
	return &Scenario{
		logger: logger.WithField("scenario", ScenarioName),
	}
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
	flags.Uint64VarP(&s.options.TotalCount, "count", "c", ScenarioDefaultOptions.TotalCount, "Total number of frame transactions to send")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t", ScenarioDefaultOptions.Throughput, "Number of frame transactions to send per slot")
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", ScenarioDefaultOptions.Rebroadcast, "Enable reliable rebroadcast system")
	flags.Float64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use in frame transactions (in gwei)")
	flags.Float64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use in frame transactions (in gwei)")
	flags.StringVar(&s.options.BaseFeeWei, "basefee-wei", ScenarioDefaultOptions.BaseFeeWei, "Max fee per gas in wei (overrides --basefee)")
	flags.StringVar(&s.options.TipFeeWei, "tipfee-wei", ScenarioDefaultOptions.TipFeeWei, "Max tip per gas in wei (overrides --tipfee)")
	flags.StringVar(&s.options.Timeout, "timeout", ScenarioDefaultOptions.Timeout, "Timeout for the scenario (e.g. '1h', '30m')")
	flags.StringVar(&s.options.ClientGroup, "client-group", ScenarioDefaultOptions.ClientGroup, "Client group to use for sending transactions")
	flags.BoolVar(&s.options.LogTxs, "log-txs", ScenarioDefaultOptions.LogTxs, "Log all submitted transactions")

	flags.StringVar(&s.options.Shapes, "shapes", ScenarioDefaultOptions.Shapes,
		fmt.Sprintf("Frame shapes to send as a weighted list, e.g. 'self-verify:10,atomic:2'. Known shapes: %s", strings.Join(shapeNames(), ", ")))
	flags.Uint64Var(&s.options.FramesPerTx, "frames-per-tx", ScenarioDefaultOptions.FramesPerTx, "Number of user operation frames for the 'batch' shape")
	flags.Uint64Var(&s.options.UserOpGas, "user-op-gas", ScenarioDefaultOptions.UserOpGas, "Execution gas limit per user operation frame")
	flags.Uint64Var(&s.options.VerifyGas, "verify-gas", ScenarioDefaultOptions.VerifyGas, "Execution gas limit for validation frames")
	flags.Uint64Var(&s.options.StateGas, "state-gas", ScenarioDefaultOptions.StateGas, "Additional state gas budget per user operation frame")
	flags.Uint64Var(&s.options.Amount, "amount", ScenarioDefaultOptions.Amount, "Transfer amount per value-bearing frame (in gwei)")
	flags.BoolVar(&s.options.RandomAmount, "random-amount", ScenarioDefaultOptions.RandomAmount, "Use random amounts up to --amount")
	flags.BoolVar(&s.options.RandomTarget, "random-target", ScenarioDefaultOptions.RandomTarget, "Send to random (non-existent) target addresses")
	flags.StringVar(&s.options.Data, "data", ScenarioDefaultOptions.Data, "Call data for user operation frames")
	flags.Uint64Var(&s.options.ExpiryOffset, "expiry-offset", ScenarioDefaultOptions.ExpiryOffset, "Seconds ahead of now to set the deadline of the 'expiry' shape")
	flags.BoolVar(&s.options.VerifyFrames, "verify-frames", ScenarioDefaultOptions.VerifyFrames, "Check per-frame receipt statuses against what each shape should produce")
	flags.BoolVar(&s.options.SkipMempoolOK, "skip-mempool-check", ScenarioDefaultOptions.SkipMempoolOK, "Skip the local public-mempool policy check before sending")

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

	shapes, err := parseShapes(s.options.Shapes)
	if err != nil {
		return err
	}

	s.shapes = shapes

	if s.options.Data != "" {
		callData, err := txbuilder.ParseBlobRefsBytes(strings.Split(s.options.Data, ","), nil)
		if err != nil {
			return fmt.Errorf("failed to parse data: %w", err)
		}

		s.callData = callData
	}

	if s.options.MaxWallets > 0 {
		s.walletPool.SetWalletCount(s.options.MaxWallets)
	} else if s.options.TotalCount > 0 {
		maxWallets := s.options.TotalCount / 50
		if maxWallets < 10 {
			maxWallets = 10
		} else if maxWallets > 1000 {
			maxWallets = 1000
		}

		s.walletPool.SetWalletCount(maxWallets)
	} else {
		walletCount := s.options.Throughput * 2
		if walletCount < 10 {
			walletCount = 10
		} else if walletCount > 1000 {
			walletCount = 1000
		}

		s.walletPool.SetWalletCount(walletCount)
	}

	if s.options.TotalCount == 0 && s.options.Throughput == 0 {
		return fmt.Errorf("neither total count nor throughput limit set, must define at least one of them (see --help for list of all flags)")
	}

	return nil
}

func (s *Scenario) Run(ctx context.Context) error {
	if err := s.checkFrameTxSupport(ctx); err != nil {
		return err
	}

	s.logger.Infof("starting frame transaction scenario with shapes: %s", strings.Join(shapeNames(), ", "))

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

	err := scenario.RunTransactionScenario(ctx, scenario.TransactionScenarioOptions{
		TotalCount:                  s.options.TotalCount,
		Throughput:                  s.options.Throughput,
		MaxPending:                  maxPending,
		ThroughputIncrementInterval: 0,
		Timeout:                     s.parseTimeout(),
		WalletPool:                  s.walletPool,

		Logger: s.logger,
		ProcessNextTxFn: func(ctx context.Context, params *scenario.ProcessNextTxParams) error {
			logger := s.logger
			receiptChan, tx, client, wallet, err := s.sendTx(ctx, params.TxIdx)

			if client != nil {
				logger = logger.WithField("rpc", client.GetName())
			}

			if tx != nil {
				logger = logger.WithField("nonce", tx.Nonce())
			}

			if wallet != nil {
				logger = logger.WithField("wallet", s.walletPool.GetWalletName(wallet.GetAddress()))
			}

			params.NotifySubmitted()
			params.OrderedLogCb(func() {
				if err != nil {
					logger.Warnf("could not send frame transaction: %v", err)
				} else if s.options.LogTxs {
					logger.Infof("sent frame tx #%6d: %v", params.TxIdx+1, tx.Hash().String())
				} else {
					logger.Debugf("sent frame tx #%6d: %v", params.TxIdx+1, tx.Hash().String())
				}
			})

			if _, waitErr := receiptChan.Wait(ctx); waitErr != nil {
				return waitErr
			}

			return err
		},
	})

	checked := s.framesChecked.Load()
	mismatch := s.framesMismatch.Load()
	missing := s.missingReceipt.Load()

	if s.options.VerifyFrames && checked > 0 {
		s.logger.Infof("frame receipt check: %d transactions verified, %d mismatched, %d without frame receipts",
			checked, mismatch, missing)
	}

	return err
}

// checkFrameTxSupport probes whether the chain implements EIP-8141 before the scenario
// starts sending.
//
// It submits a frame transaction from a throwaway key that holds no funds. A chain
// that implements the type rejects it for a reason of its own (unknown account, fee or
// balance); one that does not rejects the type itself. Failing here beats emitting a
// stream of rejections that look like a scenario bug.
func (s *Scenario) checkFrameTxSupport(ctx context.Context) error {
	txpool := s.walletPool.GetTxPool()
	if txpool == nil {
		return nil
	}

	if !txpool.IsAmsterdam() {
		return fmt.Errorf("frame transactions need the Amsterdam (EIP-8037) gas model, but --pre-amsterdam-fee-model is set")
	}

	client := s.walletPool.GetClient(spamoor.WithClientGroup(s.options.ClientGroup))
	if client == nil {
		return scenario.ErrNoClients
	}

	chainID, err := client.GetChainId(ctx)
	if err != nil {
		return fmt.Errorf("failed reading the chain id: %w", err)
	}

	key, err := crypto.GenerateKey()
	if err != nil {
		return err
	}

	probe := txtypes.NewFrameTx(
		uint256.MustFromBig(chainID),
		crypto.PubkeyToAddress(key.PublicKey), 0,
		txtypes.FrameFees{GasTipCap: uint256.NewInt(1), GasFeeCap: uint256.NewInt(1e9)},
		[]*txtypes.Frame{
			txtypes.SelfVerifyFrame(txtypes.FrameLimits{Execution: s.options.VerifyGas}),
			txtypes.UserOpFrame(&txtypes.EntryPoint, nil, nil, txtypes.FrameLimits{Execution: 21000}),
		},
		[]*txtypes.FrameSignature{txtypes.SenderSignature()},
	)

	signed, err := txtypes.SignTx(txtypes.NewTx(probe), chainID, key)
	if err != nil {
		return fmt.Errorf("failed building the frame transaction probe: %w", err)
	}

	err = client.SendTransaction(ctx, signed)
	if err != nil && isUnsupportedTypeError(err) {
		return fmt.Errorf("%s rejected an EIP-8141 frame transaction (%v): this chain does not implement frame transactions", client.GetName(), err)
	}

	return nil
}

// isUnsupportedTypeError reports whether an RPC error means the client cannot process
// the transaction type at all, as opposed to rejecting this particular transaction.
//
// The markers are deliberately narrow. An error raised once the client is simulating
// the validation prefix means the type is understood, and the probe's throwaway key
// draws exactly that: an unfunded sender makes APPROVE revert for want of balance.
func isUnsupportedTypeError(err error) bool {
	message := strings.ToLower(err.Error())

	for _, marker := range []string{
		"transaction type not supported",
		"unsupported transaction type",
		"unknown transaction type",
		"typed transaction too short",
		"not supported before", // fork gate: the type exists but is not yet active
		"invalid rlp",
		"error decoding",
	} {
		if strings.Contains(message, marker) {
			return true
		}
	}

	return false
}

func (s *Scenario) parseTimeout() time.Duration {
	if s.options.Timeout == "" {
		return 0
	}

	timeout, err := time.ParseDuration(s.options.Timeout)
	if err != nil {
		s.logger.Warnf("invalid timeout value %q, ignoring", s.options.Timeout)

		return 0
	}

	return timeout
}

func (s *Scenario) sendTx(ctx context.Context, txIdx uint64) (scenario.ReceiptChan, *txtypes.Transaction, *spamoor.Client, *spamoor.Wallet, error) {
	client := s.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, int(txIdx)),
		spamoor.WithClientGroup(s.options.ClientGroup),
	)
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(txIdx))

	if client == nil {
		return nil, nil, client, wallet, scenario.ErrNoClients
	}

	if wallet == nil {
		return nil, nil, client, wallet, scenario.ErrNoWallet
	}

	if err := wallet.ResetNoncesIfNeeded(ctx, client); err != nil {
		return nil, nil, client, wallet, err
	}

	baseFeeWei, tipFeeWei := spamoor.ResolveFees(s.options.BaseFee, s.options.TipFee, s.options.BaseFeeWei, s.options.TipFeeWei)

	feeCap, tipCap, err := s.walletPool.GetSuggestedFees(client, baseFeeWei, tipFeeWei)
	if err != nil {
		return nil, nil, client, wallet, err
	}

	shape, err := s.buildShapeFor(txIdx, wallet)
	if err != nil {
		return nil, nil, client, wallet, err
	}

	frameTx := txtypes.NewFrameTx(nil, wallet.GetAddress(), 0,
		txtypes.FrameFees{
			GasTipCap: uint256.MustFromBig(tipCap),
			GasFeeCap: uint256.MustFromBig(feeCap),
		},
		shape.frames,
		[]*txtypes.FrameSignature{txtypes.SenderSignature()},
	)

	tx, err := wallet.BuildFrameTx(frameTx)
	if err != nil {
		return nil, nil, client, wallet, err
	}

	// SignTx signs a copy, so validate what was actually built: the checks cover the
	// signature encoding and only mean anything once the entries are filled in.
	signedFrame, ok := tx.Inner().(*txtypes.FrameTx)
	if !ok {
		wallet.MarkSkippedNonce(tx.Nonce())

		return nil, nil, client, wallet, fmt.Errorf("built transaction is not a frame transaction")
	}

	if err := s.validate(signedFrame); err != nil {
		wallet.MarkSkippedNonce(tx.Nonce())

		return nil, nil, client, wallet, err
	}

	receiptChan := make(scenario.ReceiptChan, 1)

	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		ClientGroup: s.options.ClientGroup,
		Rebroadcast: s.options.Rebroadcast > 0,
		OnComplete: func(tx *txtypes.Transaction, receipt *txtypes.Receipt, err error) {
			receiptChan <- receipt
		},
		OnConfirm: func(tx *txtypes.Transaction, receipt *txtypes.Receipt) {
			s.onConfirm(client, shape, txIdx, tx, receipt)
		},
		LogFn: spamoor.GetDefaultLogFn(s.logger, string(shape.name), fmt.Sprintf("%6d", txIdx+1), tx),
	})
	if err != nil {
		// mark nonce as skipped if tx was not sent
		wallet.MarkSkippedNonce(tx.Nonce())

		return nil, nil, client, wallet, err
	}

	return receiptChan, tx, client, wallet, nil
}

// buildShapeFor assembles the frames for the shape selected for this transaction.
func (s *Scenario) buildShapeFor(txIdx uint64, wallet *spamoor.Wallet) (*builtShape, error) {
	target := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(txIdx)+1).GetAddress()
	targetEmpty := false

	if s.options.RandomTarget {
		addrBytes := make([]byte, 20)
		if _, err := rand.Read(addrBytes); err != nil {
			return nil, err
		}

		target = common.Address(addrBytes)
		targetEmpty = true
	}

	amount := uint256.NewInt(s.options.Amount)
	amount = amount.Mul(amount, uint256.NewInt(1000000000))

	if s.options.RandomAmount && !amount.IsZero() {
		n, err := rand.Int(rand.Reader, amount.ToBig())
		if err == nil {
			amount = uint256.MustFromBig(n)
		}
	}

	return buildShape(s.shapes.pick(txIdx), shapeParams{
		sender:      wallet.GetAddress(),
		target:      target,
		amount:      amount,
		data:        s.callData,
		userOpGas:   s.options.UserOpGas,
		verifyGas:   s.options.VerifyGas,
		stateGas:    s.options.StateGas,
		frameCount:  s.options.FramesPerTx,
		expiryAt:    uint64(time.Now().Unix()) + s.options.ExpiryOffset,
		targetEmpty: targetEmpty,
	})
}

// validate runs the local validity and mempool policy checks so a malformed shape is
// caught here rather than as an opaque RPC rejection.
func (s *Scenario) validate(tx *txtypes.FrameTx) error {
	if err := tx.ValidatePayload(); err != nil {
		return err
	}

	if err := tx.VerifySignatures(); err != nil {
		return err
	}

	if s.options.SkipMempoolOK {
		return nil
	}

	return tx.ValidateMempoolPrefix()
}

// onConfirm logs the confirmation and compares the frame receipts against what the
// shape should have produced.
func (s *Scenario) onConfirm(client *spamoor.Client, shape *builtShape, txIdx uint64, tx *txtypes.Transaction, receipt *txtypes.Receipt) {
	txFees := utils.GetTransactionFees(tx, receipt)
	logger := s.logger.WithField("rpc", client.GetName())

	extra := receipt.FrameExtra()
	if extra == nil {
		s.missingReceipt.Add(1)
		logger.Debugf(" frame tx %d (%s) confirmed in block #%v without frame receipts. total fee: %v gwei",
			txIdx+1, shape.name, receipt.BlockNumber.String(), txFees.TotalFeeGweiString())

		return
	}

	logger.Debugf(" frame tx %d (%s) confirmed in block #%v. payer: %v, frames: %v, exec gas: %v, state gas: %v, total fee: %v gwei",
		txIdx+1, shape.name, receipt.BlockNumber.String(), extra.Payer.Hex(),
		formatStatuses(extra), extra.TotalExecutionGas(), extra.TotalStateGas(),
		txFees.TotalFeeGweiString())

	if !s.options.VerifyFrames {
		return
	}

	s.framesChecked.Add(1)

	if mismatch := compareStatuses(shape.expectedStatus, extra); mismatch != "" {
		s.framesMismatch.Add(1)
		logger.Warnf("frame receipt mismatch for %s tx %v: %s", shape.name, tx.Hash().Hex(), mismatch)
	}
}
