package frametxfuzz

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

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

	PayloadSeed  string  `yaml:"payload_seed"`
	TxIdOffset   uint64  `yaml:"tx_id_offset"`
	Recipe       string  `yaml:"recipe"`
	Axes         string  `yaml:"axes"`
	Envelope     string  `yaml:"envelope"`
	PostTx       string  `yaml:"post_tx"`
	MaxFrames    uint64  `yaml:"max_frames"`
	InvalidRatio float64 `yaml:"invalid_ratio"`
	LogFrames    bool    `yaml:"log_frames"`
	Data         string  `yaml:"data"`

	UserOpGas uint64 `yaml:"user_op_gas"`
	VerifyGas uint64 `yaml:"verify_gas"`
	StateGas  uint64 `yaml:"state_gas"`
	Amount    uint64 `yaml:"amount"`
	Expiry    uint64 `yaml:"expiry_offset"`
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	// seed is the effective payload seed, generated once when none was configured.
	seed string

	axes             axisWeights
	pinnedExtensions *txtypes.FrameExtensions
	replay           *Recipe
	callData         []byte

	env      *environment
	coverage *coverage

	submitted   atomic.Uint64
	rootCounter atomic.Uint64
}

var ScenarioName = "frametx-fuzz"
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

	PayloadSeed:  "",
	TxIdOffset:   0,
	Axes:         "all",
	Envelope:     "auto",
	PostTx:       "auto",
	MaxFrames:    6,
	InvalidRatio: 0.05,
	LogFrames:    false,

	UserOpGas: 60000,
	VerifyGas: 5000,
	StateGas:  0,
	Amount:    20,
	Expiry:    600,
}

var ScenarioDescriptor = scenario.Descriptor{
	Name:           ScenarioName,
	Description:    "Fuzz EIP-8141 frame transactions across their envelope, prefix, batch, signature and introspection dimensions, checking every landed transaction against what the spec says it must produce",
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
	flags.Uint64VarP(&s.options.TotalCount, "count", "c", ScenarioDefaultOptions.TotalCount, "Total number of frame transactions to generate")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t", ScenarioDefaultOptions.Throughput, "Number of frame transactions to generate per slot")
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", ScenarioDefaultOptions.Rebroadcast, "Enable reliable rebroadcast system")
	flags.Float64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use (in gwei)")
	flags.Float64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use (in gwei)")
	flags.StringVar(&s.options.BaseFeeWei, "basefee-wei", ScenarioDefaultOptions.BaseFeeWei, "Max fee per gas in wei (overrides --basefee)")
	flags.StringVar(&s.options.TipFeeWei, "tipfee-wei", ScenarioDefaultOptions.TipFeeWei, "Max tip per gas in wei (overrides --tipfee)")
	flags.StringVar(&s.options.Timeout, "timeout", ScenarioDefaultOptions.Timeout, "Timeout for the scenario (e.g. '1h', '30m')")
	flags.StringVar(&s.options.ClientGroup, "client-group", ScenarioDefaultOptions.ClientGroup, "Client group to use for sending transactions")
	flags.BoolVar(&s.options.LogTxs, "log-txs", ScenarioDefaultOptions.LogTxs, "Log all submitted transactions")

	flags.StringVar(&s.options.PayloadSeed, "payload-seed", ScenarioDefaultOptions.PayloadSeed, "Hex seed for reproducible generation (empty means random, and the generated seed is logged)")
	flags.Uint64Var(&s.options.TxIdOffset, "tx-id-offset", ScenarioDefaultOptions.TxIdOffset, "Start generating from a specific transaction id")
	flags.StringVar(&s.options.Recipe, "recipe", ScenarioDefaultOptions.Recipe, "Replay a single recipe, as reported with a finding")
	flags.StringVar(&s.options.Axes, "axes", ScenarioDefaultOptions.Axes,
		fmt.Sprintf("Weighted list of dimensions to fuzz, e.g. 'nonces:5,roots:2'. Known axes: %s", strings.Join(axisNamesText(), ", ")))
	flags.StringVar(&s.options.Envelope, "envelope", ScenarioDefaultOptions.Envelope, "Envelope shape to encode: auto, full, keyed, roots, base")
	flags.StringVar(&s.options.PostTx, "post-tx", ScenarioDefaultOptions.PostTx, "EIP-7906 POST_TX frames: auto (probe the chain), on, off")
	flags.Uint64Var(&s.options.MaxFrames, "max-frames", ScenarioDefaultOptions.MaxFrames, "Maximum number of body frames per transaction")
	flags.Float64Var(&s.options.InvalidRatio, "invalid-ratio", ScenarioDefaultOptions.InvalidRatio,
		fmt.Sprintf("Share of the stream carrying a deliberate violation. Known violations: %s", strings.Join(violationNames(), ", ")))
	flags.BoolVar(&s.options.LogFrames, "log-frames", ScenarioDefaultOptions.LogFrames, "Log the per-frame result of every landed transaction")
	flags.StringVar(&s.options.Data, "data", ScenarioDefaultOptions.Data, "Call data for frames that do not carry a generated script")

	flags.Uint64Var(&s.options.UserOpGas, "user-op-gas", ScenarioDefaultOptions.UserOpGas, "Execution gas limit per body frame")
	flags.Uint64Var(&s.options.VerifyGas, "verify-gas", ScenarioDefaultOptions.VerifyGas, "Execution gas limit for validation frames")
	flags.Uint64Var(&s.options.StateGas, "state-gas", ScenarioDefaultOptions.StateGas, "Additional state gas budget per body frame")
	flags.Uint64Var(&s.options.Amount, "amount", ScenarioDefaultOptions.Amount, "Transfer amount per value-bearing frame (in gwei)")
	flags.Uint64Var(&s.options.Expiry, "expiry-offset", ScenarioDefaultOptions.Expiry, "Seconds ahead of now for an expiry frame's deadline")

	return nil
}

func (s *Scenario) Init(options *scenario.Options) error {
	s.walletPool = options.WalletPool

	if options.Config != "" {
		if err := scenario.ParseAndValidateConfig(&ScenarioDescriptor, options.Config, &s.options, s.logger); err != nil {
			return err
		}
	}

	axes, err := parseAxes(s.options.Axes)
	if err != nil {
		return err
	}

	s.axes = axes

	if err := s.parseEnvelope(); err != nil {
		return err
	}

	switch s.options.PostTx {
	case "auto", "on", "off":
	default:
		return fmt.Errorf("unknown --post-tx value %q, expected auto, on or off", s.options.PostTx)
	}

	if s.options.Recipe != "" {
		s.replay, err = ParseRecipe(s.options.Recipe)
		if err != nil {
			return err
		}
	}

	if s.options.PayloadSeed != "" {
		if _, err := utils.ParseHexSeed(s.options.PayloadSeed); err != nil {
			return fmt.Errorf("invalid payload seed: %w", err)
		}
	}

	if s.options.Data != "" {
		callData, err := txbuilder.ParseBlobRefsBytes(strings.Split(s.options.Data, ","), nil)
		if err != nil {
			return fmt.Errorf("failed to parse data: %w", err)
		}

		s.callData = callData
	}

	s.setWalletCount()

	if s.options.TotalCount == 0 && s.options.Throughput == 0 {
		return fmt.Errorf("neither total count nor throughput limit set, must define at least one of them (see --help for list of all flags)")
	}

	if s.options.InvalidRatio > 0 {
		// A deliberately invalid transaction never lands, so it consumes no nonce and
		// one wallet can fire them all. Keeping them off the generating wallets keeps
		// the managed pool, which assumes every transaction eventually confirms, out
		// of it entirely.
		s.walletPool.AddWellKnownWallet(&spamoor.WellKnownWalletConfig{
			Name:          BurnerWalletName,
			RefillAmount:  uint256.NewInt(200000000000000000),
			RefillBalance: uint256.NewInt(50000000000000000),
		})
	}

	if s.axes.enabled(axisRoots) {
		// The root source is a wallet of its own: a source is identified by the address
		// that wrote it, so keeping it fixed lets a rerun reference roots an earlier run
		// committed.
		s.walletPool.AddWellKnownWallet(&spamoor.WellKnownWalletConfig{
			Name:          RootSourceWalletName,
			RefillAmount:  uint256.NewInt(500000000000000000),
			RefillBalance: uint256.NewInt(200000000000000000),
		})
	}

	if s.axes.enabled(axisProbe) {
		// The probe contract's deployer also plays the paymaster, so it needs enough
		// balance to cover a sponsored transaction's maximum cost.
		s.walletPool.AddWellKnownWallet(&spamoor.WellKnownWalletConfig{
			Name:          ProbeDeployerWalletName,
			RefillAmount:  uint256.NewInt(2000000000000000000),
			RefillBalance: uint256.NewInt(1000000000000000000),
		})

		// Delegations are carried by a separate wallet: an EIP-7702 authorization
		// signed by the transaction's own sender needs the nonce after the
		// transaction's, and keeping both sequences in one wallet leaves a nonce gap.
		s.walletPool.AddWellKnownWallet(&spamoor.WellKnownWalletConfig{
			Name:          ProbeDelegatorWalletName,
			RefillAmount:  uint256.NewInt(500000000000000000),
			RefillBalance: uint256.NewInt(200000000000000000),
		})
	}

	return nil
}

// setWalletCount sizes the wallet pool.
//
// The public mempool keeps one pending frame transaction per sender, so throughput is
// wallets divided by confirmation latency rather than a function of pending depth. A
// pool sized like an ordinary scenario's starves the generator at a few transactions per
// block regardless of the configured throughput.
func (s *Scenario) setWalletCount() {
	if s.options.MaxWallets > 0 {
		s.walletPool.SetWalletCount(s.options.MaxWallets)

		return
	}

	count := s.options.Throughput * 12
	if s.options.TotalCount > 0 && s.options.TotalCount/4 < count {
		count = s.options.TotalCount / 4
	}

	if count < 20 {
		count = 20
	} else if count > 1000 {
		count = 1000
	}

	s.walletPool.SetWalletCount(count)
}

// parseEnvelope resolves the --envelope option, which pins the payload shape instead of
// taking the one the chain's predeploys imply.
func (s *Scenario) parseEnvelope() error {
	choices := map[string]txtypes.FrameExtensions{
		"base":  0,
		"keyed": txtypes.FrameExtKeyedNonces,
		"roots": txtypes.FrameExtRecentRoots,
		"full":  txtypes.FrameExtAll,
	}

	spec := strings.ToLower(strings.TrimSpace(s.options.Envelope))
	if spec == "" || spec == "auto" {
		return nil
	}

	extensions, ok := choices[spec]
	if !ok {
		return fmt.Errorf("unknown envelope %q, known: auto, base, keyed, roots, full", spec)
	}

	s.pinnedExtensions = &extensions

	return nil
}

func (s *Scenario) Run(ctx context.Context) error {
	s.seed = s.options.PayloadSeed
	if s.seed == "" {
		random := make([]byte, 32)
		if _, err := rand.Read(random); err != nil {
			return err
		}

		s.seed = hex.EncodeToString(random)
		s.logger.Infof("generated payload seed for this run: 0x%s", s.seed)
	}

	env, err := s.setupEnvironment(ctx)
	if err != nil {
		return err
	}

	s.env = env
	s.coverage = newCoverage(s.logger, strings.TrimPrefix(s.seed, "0x"))

	s.probeCapabilities(ctx, env)

	if env.roots != nil {
		go s.maintainRoots(ctx)
	}

	maxPending := s.options.MaxPending
	if maxPending == 0 {
		maxPending = s.walletPool.GetWalletCount()
	}

	err = scenario.RunTransactionScenario(ctx, scenario.TransactionScenarioOptions{
		TotalCount: s.options.TotalCount,
		Throughput: s.options.Throughput,
		MaxPending: maxPending,
		Timeout:    s.parseTimeout(),
		WalletPool: s.walletPool,
		Logger:     s.logger,
		ProcessNextTxFn: func(ctx context.Context, params *scenario.ProcessNextTxParams) error {
			return s.processTx(ctx, params)
		},
	})

	s.logger.Infof("frametx-fuzz coverage: %s", s.coverage.summary())

	for _, line := range s.coverage.refusals() {
		s.logger.Infof("refused %s", line)
	}

	return err
}

// processTx generates, submits and awaits one transaction.
func (s *Scenario) processTx(ctx context.Context, params *scenario.ProcessNextTxParams) error {
	receiptChan, result, client, err := s.sendTx(ctx, params.TxIdx)

	logger := s.logger
	if client != nil {
		logger = logger.WithField("rpc", client.GetName())
	}

	params.NotifySubmitted()
	params.OrderedLogCb(func() {
		switch {
		case err != nil:
			logger.Warnf("could not send frame transaction %d: %v", params.TxIdx+1, err)
		case s.options.LogTxs:
			logger.Infof("sent frame tx #%6d (%s)", params.TxIdx+1, result.recipe)
		default:
			logger.Debugf("sent frame tx #%6d (%s)", params.TxIdx+1, result.recipe)
		}
	})

	if receiptChan != nil {
		if _, waitErr := receiptChan.Wait(ctx); waitErr != nil {
			return waitErr
		}
	}

	return err
}

// recipeFor draws the recipe for a transaction index, or returns the pinned replay.
func (s *Scenario) recipeFor(txIdx uint64) *Recipe {
	if s.replay != nil {
		return s.replay
	}

	index := txIdx + s.options.TxIdOffset
	rng := utils.NewDeterministicRNGWithSeed(index, s.seed)

	return Draw(rng, index, DrawOptions{
		Axes:          s.axes,
		MaxBodyFrames: int(s.options.MaxFrames),
		AllowPostTx:   s.env.allowPostTx,
		AllowContract: s.env.contractCount > 0,
		AllowRoots:    s.env.roots != nil,
		AllowKeyed:    s.env.nonces != nil,
		AllowProbe:    s.env.probe != nil,

		InvalidChance: s.options.InvalidRatio,
		Violations:    violationNames(),
	})
}

// sendTx builds and submits one generated transaction.
func (s *Scenario) sendTx(ctx context.Context, txIdx uint64) (scenario.ReceiptChan, *build, *spamoor.Client, error) {
	client := s.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, int(txIdx)),
		spamoor.WithClientGroup(s.options.ClientGroup),
	)
	if client == nil {
		return nil, nil, nil, scenario.ErrNoClients
	}

	feeCap, tipCap, err := s.fees(client)
	if err != nil {
		return nil, nil, client, err
	}

	recipe := s.recipeFor(txIdx)

	result, err := s.buildRecipe(ctx, client, s.env, recipe, feeCap, tipCap)
	if err != nil {
		return nil, nil, client, err
	}

	if err := result.sender.ResetNoncesIfNeeded(ctx, client); err != nil {
		return nil, result, client, err
	}

	s.coverage.record(result.coverage)

	if recipe.Invalid != "" {
		if err := result.sender.PrepareFrameTx(result.tx); err != nil {
			return nil, result, client, err
		}

		return nil, result, client, s.sendInvalid(ctx, client, result)
	}

	if err := result.sender.PrepareFrameTx(result.tx); err != nil {
		return nil, result, client, err
	}

	if result.p256 != nil {
		// The P256 entry's signer is part of the signature hash, so it has to be
		// signed before the secp256k1 entries.
		if err := result.tx.SignEntryP256(p256Index(result.tx), result.p256); err != nil {
			result.sender.MarkSkippedNonce(result.tx.NonceSeq)

			return nil, result, client, err
		}
	}

	tx, err := result.sender.SignFrameTx(result.tx)
	if err != nil {
		return nil, result, client, err
	}

	signed, ok := tx.Inner().(*txtypes.FrameTx)
	if !ok {
		result.sender.MarkSkippedNonce(tx.Nonce())

		return nil, result, client, fmt.Errorf("built transaction is not a frame transaction")
	}

	receiptChan := make(scenario.ReceiptChan, 1)

	err = s.walletPool.GetTxPool().SendTransaction(ctx, result.sender, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		ClientGroup: s.options.ClientGroup,
		Rebroadcast: s.options.Rebroadcast > 0,
		OnComplete: func(tx *txtypes.Transaction, receipt *txtypes.Receipt, err error) {
			receiptChan <- receipt
		},
		OnConfirm: func(tx *txtypes.Transaction, receipt *txtypes.Receipt) {
			s.onConfirm(result, signed, receipt)
		},
		LogFn: spamoor.GetDefaultLogFn(s.logger, ScenarioName, fmt.Sprintf("%6d", txIdx+1), tx),
	})
	if err != nil {
		result.sender.MarkSkippedNonce(tx.Nonce())
		s.coverage.refusedOne(result.recipe, err.Error())

		return nil, result, client, err
	}

	s.submitted.Add(1)

	return receiptChan, result, client, nil
}

// onConfirm records what a landed transaction produced.
//
// The per-frame result is logged rather than judged: a status vector this scenario
// disagreed with would say more about the reading baked in here than about the client.
func (s *Scenario) onConfirm(result *build, tx *txtypes.FrameTx, receipt *txtypes.Receipt) {
	if result.nonces != nil {
		s.env.nonces.consumed(result.sender.GetAddress(), result.nonces)
	}

	s.coverage.confirmedOne()

	if !s.options.LogFrames {
		return
	}

	extra := receipt.FrameExtra()
	if extra == nil {
		s.logger.Debugf("frame transaction %v confirmed without a frame receipt section", receipt.TxHash)

		return
	}

	s.logger.Debugf("frame transaction %v: payer %v, frames %s, durable %s, exec gas %d, state gas %d",
		receipt.TxHash, extra.Payer, renderStatuses(extra), renderDurable(extra, tx),
		extra.TotalExecutionGas(), extra.TotalStateGas())
}

// p256Index returns the index of the transaction's P256 entry.
func p256Index(tx *txtypes.FrameTx) int {
	for i, sig := range tx.Signatures {
		if sig.Scheme == txtypes.SigSchemeP256 {
			return i
		}
	}

	return -1
}

// deadline returns the expiry frame's deadline.
func (s *Scenario) deadline() uint64 {
	return uint64(time.Now().Unix()) + s.options.Expiry
}

// currentSlot returns the slot the next block is expected to be in, used to decide
// which committed roots may still be referenced.
func (s *Scenario) currentSlot() uint64 {
	if s.env == nil || s.env.roots == nil {
		return 0
	}

	return s.env.roots.slotFor(uint64(time.Now().Unix()))
}

// amount returns the value a transfer frame carries.
func (s *Scenario) amount(recipe *Recipe) *uint256.Int {
	amount := uint256.NewInt(s.options.Amount)
	amount = amount.Mul(amount, uint256.NewInt(1000000000))

	if amount.IsZero() {
		return amount
	}

	// Derived from the recipe rather than drawn at random, so a replayed recipe moves
	// the same value.
	spread := uint256.NewInt(recipe.Index%97 + 1)

	return new(uint256.Int).Div(new(uint256.Int).Mul(amount, spread), uint256.NewInt(97))
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

// generateP256Key returns a fresh NIST P-256 key for a P256 signature entry.
func generateP256Key() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

// axisNamesText renders the axis names for help text.
func axisNamesText() []string {
	names := make([]string, len(axisNames))
	for i, name := range axisNames {
		names[i] = string(name)
	}

	return names
}

// parseAxes parses a "name:weight,name:weight" selection. An empty value, or "all",
// enables every axis with equal weight.
func parseAxes(spec string) (axisWeights, error) {
	spec = strings.TrimSpace(spec)
	weights := axisWeights{}

	if spec == "" || spec == "all" {
		for _, name := range axisNames {
			weights[name] = 1
		}

		return weights, nil
	}

	for _, entry := range strings.Split(spec, ",") {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		name := entry
		weight := uint64(1)

		if idx := strings.LastIndex(entry, ":"); idx >= 0 {
			name = strings.TrimSpace(entry[:idx])

			parsed, err := strconv.ParseUint(strings.TrimSpace(entry[idx+1:]), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid weight in %q: %w", entry, err)
			}

			weight = parsed
		}

		if !knownAxis(axis(name)) {
			return nil, fmt.Errorf("unknown axis %q, known axes: %s", name, strings.Join(axisNamesText(), ", "))
		}

		weights[axis(name)] = weight
	}

	if len(weights) == 0 {
		return nil, fmt.Errorf("no axes enabled")
	}

	return weights, nil
}

// knownAxis reports whether name is an axis the generator understands.
func knownAxis(name axis) bool {
	for _, known := range axisNames {
		if known == name {
			return true
		}
	}

	return false
}
