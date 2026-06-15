package txfuzzinvalid

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/spamoor"

	"github.com/holiman/uint256"
)

const (
	unstuckMaxTries = 20
	unstuckBatch    = 16
	unstuckWait     = 60 * time.Second
	summaryEvery    = 500
)

type ScenarioOptions struct {
	TotalCount  uint64 `yaml:"total_count"`
	Throughput  uint64 `yaml:"throughput"`
	MaxPending  uint64 `yaml:"max_pending"`
	MaxWallets  uint64 `yaml:"max_wallets"`
	Timeout     string `yaml:"timeout"`
	ClientGroup string `yaml:"client_group"`
	LogTxs      bool   `yaml:"log_txs"`

	Categories  string `yaml:"categories"`   // comma list of invalid categories (or "all")
	PayloadSeed string `yaml:"payload_seed"` // optional hex seed for reproducible fuzzing
}

type Scenario struct {
	options    ScenarioOptions
	logger     logrus.FieldLogger
	walletPool *spamoor.WalletPool
	seed       []byte
	enabled    []invalidCategory

	statesMu sync.Mutex
	states   map[common.Address]*walletState

	// stats
	submitted atomic.Uint64
	accepted  atomic.Uint64
	statsMu   sync.Mutex
	catStats  map[string]*catStat
}

type catStat struct {
	sent     uint64
	accepted uint64
}

// walletState tracks the out-of-pool nonce/balance for a burner wallet. These
// wallets are funded by the wallet pool but never submit through it, so we
// manage their nonce ourselves and unstuck them when an invalid tx unexpectedly
// lands.
type walletState struct {
	mu           sync.Mutex
	baseNonce    uint64
	balance      *big.Int
	loaded       bool
	needsUnstuck bool
}

var ScenarioName = "tx-fuzz-invalid"
var ScenarioDefaultOptions = ScenarioOptions{
	TotalCount:  0,
	Throughput:  20,
	MaxPending:  50,
	MaxWallets:  0,
	Timeout:     "",
	ClientGroup: "",
	LogTxs:      false,
	Categories:  "all",
	PayloadSeed: "",
}

var ScenarioDescriptor = scenario.Descriptor{
	Name:           ScenarioName,
	Description:    "Fuzzes transaction validation by firing deliberately-invalid transactions (bad chainid/nonce/gas/sig/fees, malformed RLP, etc.) from disposable burner wallets, out-of-pool and fire-and-forget",
	DefaultOptions: ScenarioDefaultOptions,
	NewScenario:    newScenario,
}

func newScenario(logger logrus.FieldLogger) scenario.Scenario {
	return &Scenario{
		options:  ScenarioDefaultOptions,
		logger:   logger.WithField("scenario", ScenarioName),
		states:   make(map[common.Address]*walletState),
		catStats: make(map[string]*catStat),
	}
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
	flags.Uint64VarP(&s.options.TotalCount, "count", "c", ScenarioDefaultOptions.TotalCount, "Total number of invalid transactions to send")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t", ScenarioDefaultOptions.Throughput, "Number of invalid transactions to send per slot")
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum number of concurrent in-flight submissions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of burner wallets to use")
	flags.StringVar(&s.options.Timeout, "timeout", ScenarioDefaultOptions.Timeout, "Timeout for the scenario (e.g. '1h', '30m', '5s') - empty means no timeout")
	flags.StringVar(&s.options.ClientGroup, "client-group", ScenarioDefaultOptions.ClientGroup, "Client group to use for sending transactions")
	flags.BoolVar(&s.options.LogTxs, "log-txs", ScenarioDefaultOptions.LogTxs, "Log every submission and its rejection reason")
	flags.StringVar(&s.options.Categories, "categories", ScenarioDefaultOptions.Categories, "Comma-separated invalid categories to fuzz (or 'all'): "+strings.Join(categoryNames(), ", "))
	flags.StringVar(&s.options.PayloadSeed, "payload-seed", ScenarioDefaultOptions.PayloadSeed, "Custom hex seed for reproducible fuzzing (e.g. 0x1234abcd, empty means random)")
	return nil
}

func (s *Scenario) Init(options *scenario.Options) error {
	s.walletPool = options.WalletPool

	if options.Config != "" {
		if err := scenario.ParseAndValidateConfig(&ScenarioDescriptor, options.Config, &s.options, s.logger); err != nil {
			return err
		}
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
		count := s.options.Throughput * 5
		if count < 10 {
			count = 10
		} else if count > 1000 {
			count = 1000
		}
		s.walletPool.SetWalletCount(count)
	}

	// Burner wallets are a small, REUSED pool - never one funded wallet per tx -
	// which keeps tx load on the shared root wallet low (per-tx funding would
	// flood the root and break in combination with other scenarios). They also
	// barely spend: most invalid txs are rejected pre-execution and never touch
	// balance; only unstuck self-transfers and the rare accepted tx cost gas. So
	// fund them minimally and infrequently.
	s.walletPool.SetRefillAmount(uint256.NewInt(50_000_000_000_000_000))  // 0.05 ETH
	s.walletPool.SetRefillBalance(uint256.NewInt(10_000_000_000_000_000)) // 0.01 ETH

	if s.options.TotalCount == 0 && s.options.Throughput == 0 {
		return fmt.Errorf("neither total count nor throughput limit set, must define at least one of them (see --help for list of all flags)")
	}

	enabled, err := parseCategories(s.options.Categories)
	if err != nil {
		return err
	}
	s.enabled = enabled

	if s.options.PayloadSeed != "" {
		if err := validateSeed(s.options.PayloadSeed); err != nil {
			return fmt.Errorf("invalid payload seed: %v", err)
		}
	}

	return nil
}

func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Infof("starting scenario: %s", ScenarioName)
	defer func() {
		s.logSummary()
		s.logger.Infof("scenario %s finished.", ScenarioName)
	}()

	seedStr := s.options.PayloadSeed
	if seedStr == "" {
		randomBytes := make([]byte, 32)
		rand.Read(randomBytes)
		seedStr = hex.EncodeToString(randomBytes)
		s.logger.Infof("Generated random seed for this run: 0x%s", seedStr)
	} else {
		s.logger.Infof("Using provided seed: %s", seedStr)
	}
	s.seed = []byte(seedStr)

	maxPending := s.options.MaxPending
	if maxPending == 0 {
		maxPending = s.options.Throughput * 5
		if maxPending == 0 {
			maxPending = 1000
		}
	}

	var timeout time.Duration
	var err error
	if s.options.Timeout != "" {
		timeout, err = time.ParseDuration(s.options.Timeout)
		if err != nil {
			return fmt.Errorf("invalid timeout value: %v", err)
		}
		s.logger.Infof("Timeout set to %v", timeout)
	}

	s.logger.WithFields(logrus.Fields{
		"total":      s.options.TotalCount,
		"throughput": s.options.Throughput,
		"maxPending": maxPending,
		"categories": s.options.Categories,
	}).Info("Starting invalid-transaction fuzzer (out-of-pool, fire-and-forget)")

	return scenario.RunTransactionScenario(ctx, scenario.TransactionScenarioOptions{
		TotalCount: s.options.TotalCount,
		Throughput: s.options.Throughput,
		MaxPending: maxPending,
		Timeout:    timeout,
		WalletPool: s.walletPool,
		Logger:     s.logger.(*logrus.Entry),
		ProcessNextTxFn: func(ctx context.Context, params *scenario.ProcessNextTxParams) error {
			category, client, wallet, sendErr := s.sendInvalid(ctx, params.TxIdx)

			logger := s.logger
			if client != nil {
				logger = logger.WithField("rpc", client.GetName())
			}
			if wallet != nil {
				logger = logger.WithField("wallet", s.walletPool.GetWalletName(wallet.GetAddress()))
			}

			params.NotifySubmitted()
			params.OrderedLogCb(func() {
				if s.options.LogTxs {
					logger.Infof("invalid tx %6d (%s): %v", params.TxIdx+1, category, rejectionText(sendErr))
				} else {
					logger.Debugf("invalid tx %6d (%s): %v", params.TxIdx+1, category, rejectionText(sendErr))
				}
			})

			if n := s.submitted.Load(); n > 0 && n%summaryEvery == 0 {
				s.logSummary()
			}

			// fire-and-forget: an RPC rejection is the expected outcome, not a
			// scenario failure, so we never surface it as an error.
			return nil
		},
	})
}

// sendInvalid builds and raw-submits one invalid transaction from a burner
// wallet, recording whether the node (correctly) rejected it or (notably)
// accepted it.
func (s *Scenario) sendInvalid(ctx context.Context, txIdx uint64) (string, *spamoor.Client, *spamoor.Wallet, error) {
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(txIdx))
	if wallet == nil {
		return "", nil, nil, fmt.Errorf("no wallet available")
	}
	client := s.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, int(txIdx)),
		spamoor.WithClientGroup(s.options.ClientGroup),
	)
	if client == nil {
		return "", nil, wallet, fmt.Errorf("no client available")
	}

	st := s.stateFor(wallet.GetAddress())
	st.mu.Lock()
	defer st.mu.Unlock()

	// Recycle: if a prior tx unexpectedly occupied a nonce, clear it before reuse.
	if st.needsUnstuck {
		s.unstuck(ctx, client, wallet, st)
	}
	if !st.loaded {
		s.loadState(ctx, client, wallet, st)
	}

	rng := rngFor(s.seed, txIdx)
	res, err := generateInvalid(rng, s.enabled, genInput{
		key:           wallet.GetPrivateKey(),
		from:          wallet.GetAddress(),
		chainID:       wallet.GetChainId().Uint64(),
		baseNonce:     st.baseNonce,
		balance:       st.balance,
		blockGasLimit: s.walletPool.GetTxPool().GetCurrentGasLimit(),
	})
	if err != nil {
		return "", client, wallet, err
	}

	sendErr := client.SendRawTransaction(ctx, res.raw)
	s.record(res.category.name, sendErr)

	if sendErr == nil {
		// The node accepted a transaction we built to be invalid. That is the
		// interesting signal this scenario hunts for - surface it loudly.
		s.accepted.Add(1)
		s.logger.Warnf("node %s ACCEPTED invalid tx (category=%s) from %s - possible validation gap",
			client.GetName(), res.category.name, wallet.GetAddress().Hex())
		// Accepted txs may occupy a nonce; force a resync (and unstuck if it can
		// stick at the current nonce) before this wallet is reused.
		st.loaded = false
		if res.category.canStick {
			st.needsUnstuck = true
		}
	}

	return res.category.name, client, wallet, nil
}

func (s *Scenario) stateFor(addr common.Address) *walletState {
	s.statesMu.Lock()
	defer s.statesMu.Unlock()
	st := s.states[addr]
	if st == nil {
		st = &walletState{balance: big.NewInt(0)}
		s.states[addr] = st
	}
	return st
}

// loadState fetches the on-chain nonce and balance for a burner wallet. Must be
// called with st.mu held.
func (s *Scenario) loadState(ctx context.Context, client *spamoor.Client, wallet *spamoor.Wallet, st *walletState) {
	if nonce, err := client.GetNonceAt(ctx, wallet.GetAddress(), nil); err == nil {
		st.baseNonce = nonce
	}
	if bal, err := client.GetBalanceAt(ctx, wallet.GetAddress()); err == nil {
		st.balance = bal
	}
	st.loaded = true
}

// unstuck mirrors tx-fuzz's tryUnstuck: while the wallet's pending nonce is
// ahead of its latest (confirmed) nonce, it replaces the stuck nonces with
// fee-bumped self-transfers until the wallet is clean again. Must be called with
// st.mu held.
func (s *Scenario) unstuck(ctx context.Context, client *spamoor.Client, wallet *spamoor.Wallet, st *walletState) {
	addr := wallet.GetAddress()
	for try := 0; try < unstuckMaxTries; try++ {
		if ctx.Err() != nil {
			return
		}
		latest, err := client.GetNonceAt(ctx, addr, nil)
		if err != nil {
			return
		}
		pending, err := client.GetPendingNonceAt(ctx, addr)
		if err != nil {
			return
		}
		if pending <= latest {
			st.baseNonce = latest
			st.needsUnstuck = false
			st.loaded = true
			return
		}

		feeCap, tip := s.unstuckFees(ctx, client)
		end := pending
		if end > latest+unstuckBatch {
			end = latest + unstuckBatch
		}

		var last *types.Transaction
		for n := latest; n < end; n++ {
			raw, signed, err := buildSelfTransfer(wallet, n, feeCap, tip)
			if err != nil {
				continue
			}
			// ignore submit errors ("already known", "replacement underpriced") -
			// the next loop iteration re-reads state and bumps again if needed.
			_ = client.SendRawTransaction(ctx, raw)
			last = signed
		}

		if last != nil {
			wctx, cancel := context.WithTimeout(ctx, unstuckWait)
			_, _ = bind.WaitMined(wctx, client.GetEthClient(), last)
			cancel()
		}
	}

	s.logger.Warnf("could not unstuck burner wallet %s after %d tries", addr.Hex(), unstuckMaxTries)
	if latest, err := client.GetNonceAt(ctx, addr, nil); err == nil {
		st.baseNonce = latest
		st.loaded = true
	}
	st.needsUnstuck = false
}

// unstuckFees returns aggressively-bumped fees so the self-transfers replace any
// underpriced txs sitting at those nonces.
func (s *Scenario) unstuckFees(ctx context.Context, client *spamoor.Client) (*big.Int, *big.Int) {
	const gwei = 1_000_000_000
	feeCap := big.NewInt(100 * gwei)
	tip := big.NewInt(10 * gwei)
	if gasCap, tipCap, err := client.GetSuggestedFee(ctx); err == nil && gasCap != nil && tipCap != nil {
		feeCap = new(big.Int).Mul(gasCap, big.NewInt(2))
		tip = new(big.Int).Mul(tipCap, big.NewInt(2))
		if tip.Sign() == 0 {
			tip = big.NewInt(gwei)
		}
	}
	return feeCap, tip
}

// buildSelfTransfer builds and signs a minimal 1-wei self-transfer used to fill
// or replace a stuck nonce.
func buildSelfTransfer(wallet *spamoor.Wallet, nonce uint64, feeCap, tip *big.Int) ([]byte, *types.Transaction, error) {
	to := wallet.GetAddress()
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   wallet.GetChainId(),
		Nonce:     nonce,
		GasTipCap: tip,
		GasFeeCap: feeCap,
		Gas:       21000,
		To:        &to,
		Value:     big.NewInt(1),
	})
	signed, err := types.SignTx(tx, types.LatestSignerForChainID(wallet.GetChainId()), wallet.GetPrivateKey())
	if err != nil {
		return nil, nil, err
	}
	raw, err := signed.MarshalBinary()
	if err != nil {
		return nil, nil, err
	}
	return raw, signed, nil
}

func (s *Scenario) record(category string, sendErr error) {
	s.submitted.Add(1)
	s.statsMu.Lock()
	defer s.statsMu.Unlock()
	cs := s.catStats[category]
	if cs == nil {
		cs = &catStat{}
		s.catStats[category] = cs
	}
	cs.sent++
	if sendErr == nil {
		cs.accepted++
	}
}

func (s *Scenario) logSummary() {
	s.statsMu.Lock()
	defer s.statsMu.Unlock()

	total := s.submitted.Load()
	accepted := s.accepted.Load()
	parts := make([]string, 0, len(categories))
	for _, c := range categories {
		if cs := s.catStats[c.name]; cs != nil {
			parts = append(parts, fmt.Sprintf("%s=%d/%d", c.name, cs.accepted, cs.sent))
		}
	}
	s.logger.WithFields(logrus.Fields{
		"submitted":       total,
		"accepted":        accepted,
		"accepted_by_cat": strings.Join(parts, " "),
	}).Info("invalid fuzzer summary (accepted/sent per category; accepted>0 is a potential finding)")
}

func categoryNames() []string {
	names := make([]string, 0, len(categories))
	for _, c := range categories {
		names = append(names, c.name)
	}
	return names
}

func parseCategories(s string) ([]invalidCategory, error) {
	if s == "" || s == "all" {
		return categories, nil
	}
	var out []invalidCategory
	seen := map[string]bool{}
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(strings.ToLower(part))
		if part == "" {
			continue
		}
		c, ok := categoriesByName[part]
		if !ok {
			return nil, fmt.Errorf("unknown invalid category %q (valid: %s, all)", part, strings.Join(categoryNames(), ", "))
		}
		if !seen[part] {
			out = append(out, c)
			seen[part] = true
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no valid categories selected")
	}
	return out, nil
}

func validateSeed(seed string) error {
	clean := strings.TrimPrefix(seed, "0x")
	if _, err := hex.DecodeString(clean); err != nil {
		return fmt.Errorf("seed must be a valid hex string (with or without 0x prefix): %v", err)
	}
	return nil
}

// rejectionText renders a submit result for logging.
func rejectionText(err error) string {
	if err == nil {
		return "ACCEPTED (unexpected!)"
	}
	return "rejected: " + err.Error()
}
