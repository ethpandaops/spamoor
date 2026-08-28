package frametx

import (
	"context"
	"crypto/rand"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
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
	Envelope      string `yaml:"envelope"`
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
	extensions txtypes.FrameExtensions
	autoDetect bool

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
	Envelope:     "auto",
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
	flags.StringVar(&s.options.Envelope, "envelope", ScenarioDefaultOptions.Envelope,
		"Envelope shape to encode: auto, full (8141+8250+8272), keyed (8141+8250), roots (8141+8272), base (8141)")
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

	extensions, autoDetect, err := parseEnvelope(s.options.Envelope)
	if err != nil {
		return err
	}

	s.extensions = extensions
	s.autoDetect = autoDetect

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

	if s.shapes.requiresPostTx() {
		s.logger.Warnf("EIP-7906 POST_TX shapes were selected; a chain without EIP-7906 rejects frame mode 3")
	}

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

// checkFrameTxSupport establishes what the chain supports before the scenario sends
// anything.
//
// Everything here is read from chain state rather than inferred from error text. Each
// of these EIPs installs a predeploy at activation and requires the address to be
// empty beforehand, so the presence of that account is an exact, client-independent
// signal:
//
//	EIP-8141  EXPIRY_VERIFIER      0x…8141
//	EIP-8250  NONCE_MANAGER        0x…8250
//	EIP-8272  RECENT_ROOT_ADDRESS  0x…8272
//
// The envelope shape follows from which of the two extensions are active, which is
// what makes it safe to encode: a wrong guess would fail to decode on every
// transaction rather than once at startup.
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

	frames, err := predeployActive(ctx, client, txtypes.ExpiryVerifier)
	if err != nil {
		return fmt.Errorf("failed reading the expiry verifier predeploy: %w", err)
	}

	if !frames {
		return fmt.Errorf("%s has no code at the EIP-8141 expiry verifier predeploy %s: this chain does not implement frame transactions",
			client.GetName(), txtypes.ExpiryVerifier)
	}

	if !s.autoDetect {
		s.logger.Infof("using pinned frame transaction envelope: %s", s.extensions)

		return nil
	}

	extensions, err := detectEnvelope(ctx, client)
	if err != nil {
		return err
	}

	s.extensions = extensions
	s.logger.Infof("detected frame transaction envelope: %s", extensions)

	return nil
}

// detectEnvelope reads the active envelope extensions from the predeploys each one
// installs at activation.
func detectEnvelope(ctx context.Context, client *spamoor.Client) (txtypes.FrameExtensions, error) {
	var extensions txtypes.FrameExtensions

	keyed, err := predeployActive(ctx, client, txtypes.NonceManager)
	if err != nil {
		return 0, fmt.Errorf("failed reading the EIP-8250 nonce manager predeploy: %w", err)
	}

	if keyed {
		extensions |= txtypes.FrameExtKeyedNonces
	}

	roots, err := predeployActive(ctx, client, txtypes.RecentRootAddress)
	if err != nil {
		return 0, fmt.Errorf("failed reading the EIP-8272 recent root predeploy: %w", err)
	}

	if roots {
		extensions |= txtypes.FrameExtRecentRoots
	}

	return extensions, nil
}

// predeployActive reports whether a predeploy account has been installed.
//
// Activation sets both code and nonce 1, and the fork configuration must pick an
// address that is empty beforehand. Nonce is checked as well as code because one of
// these codes is still TBD in its EIP, and an account with nonce 1 is unambiguous
// either way.
func predeployActive(ctx context.Context, client *spamoor.Client, address common.Address) (bool, error) {
	code, err := client.GetCodeAt(ctx, address)
	if err != nil {
		return false, err
	}

	if len(code) > 0 {
		return true, nil
	}

	nonce, err := client.GetNonceAt(ctx, address, nil)
	if err != nil {
		return false, err
	}

	return nonce > 0, nil
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

	frameTx := txtypes.NewFrameTxWithExtensions(s.extensions, nil, wallet.GetAddress(), 0,
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

	frameTx, _ := tx.Inner().(*txtypes.FrameTx)

	bodyReverted := false
	if frameTx != nil {
		bodyReverted = extra.PostTxReverted(frameTx)
	}

	logger.Debugf(" frame tx %d (%s) confirmed in block #%v. payer: %v, frames: %v, durable: %v, exec gas: %v, state gas: %v, total fee: %v gwei",
		txIdx+1, shape.name, receipt.BlockNumber.String(), extra.Payer.Hex(),
		formatStatuses(extra), formatDurable(extra, frameTx),
		extra.TotalExecutionGas(), extra.TotalStateGas(),
		txFees.TotalFeeGweiString())

	if !s.options.VerifyFrames {
		return
	}

	s.framesChecked.Add(1)

	if mismatch := compareStatuses(shape.expectedStatus, extra); mismatch != "" {
		s.framesMismatch.Add(1)
		logger.Warnf("frame receipt mismatch for %s tx %v: %s", shape.name, tx.Hash().Hex(), mismatch)

		return
	}

	// A POST_TX failure must revert the whole execution body, which the per-frame
	// statuses alone do not show.
	if bodyReverted != shape.expectBodyReverted {
		s.framesMismatch.Add(1)
		logger.Warnf("frame body revert mismatch for %s tx %v: body reverted = %v, expected %v",
			shape.name, tx.Hash().Hex(), bodyReverted, shape.expectBodyReverted)
	}
}
