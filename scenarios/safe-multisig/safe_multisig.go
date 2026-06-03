package safemultisig

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/scenarios/safe-multisig/contract"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/utils"
)

type ScenarioOptions struct {
	TotalCount        uint64  `yaml:"total_count"`
	Throughput        uint64  `yaml:"throughput"`
	MaxPending        uint64  `yaml:"max_pending"`
	MaxWallets        uint64  `yaml:"max_wallets"`
	Rebroadcast       uint64  `yaml:"rebroadcast"`
	BaseFee           float64 `yaml:"base_fee"`
	TipFee            float64 `yaml:"tip_fee"`
	BaseFeeWei        string  `yaml:"base_fee_wei"`
	TipFeeWei         string  `yaml:"tip_fee_wei"`
	MinOwners         uint64  `yaml:"min_owners"`
	MaxOwners         uint64  `yaml:"max_owners"`
	Threshold         uint64  `yaml:"threshold"`
	SafesPerWallet    uint64  `yaml:"safes_per_wallet"`
	ContractRatio     float64 `yaml:"contract_ratio"`
	RecreateRate      float64 `yaml:"recreate_rate"`
	BurnRounds        uint64  `yaml:"burn_rounds"`
	EoaValue          uint64  `yaml:"eoa_value"`
	FundingInterval   uint64  `yaml:"funding_interval"`
	GasLimit          uint64  `yaml:"gas_limit"`
	Timeout           string  `yaml:"timeout"`
	ClientGroup       string  `yaml:"client_group"`
	DeployClientGroup string  `yaml:"deploy_client_group"`
	LogTxs            bool    `yaml:"log_txs"`
}

// safeEntry is one deployed multisig managed by a single executor wallet: its
// address, bound instance, owner signers (sorted ascending), threshold, current
// Safe nonce, and cached EIP-712 domain separator.
type safeEntry struct {
	addr      common.Address
	instance  *contract.Safe
	owners    []*signer
	threshold int
	nonce     uint64
	domainSep [32]byte
	// needResync is set when an execTransaction for this safe reverts. A revert
	// leaves the on-chain Safe nonce unchanged while the local counter advanced,
	// so the next use re-reads the nonce from chain to recover instead of
	// cascading further failures.
	needResync bool
	// pending is true for a freshly created safe whose creation transaction has
	// not confirmed yet; pending safes are skipped when selecting a safe.
	pending bool
	// removing is true while a teardown (fund-forwarding) transaction for this
	// safe is in flight; such safes are skipped when selecting a safe.
	removing bool
	// balance is the approximate ETH balance held by the safe, refreshed from
	// chain by the top-up loop and adjusted optimistically on value transfers.
	// Used to decide whether an EOA value transfer can be sourced from the safe.
	balance *big.Int
}

// walletSafePool holds the safes managed by a single executor (child) wallet.
// All of a wallet's safes are only ever exercised via execTransaction sent by
// that same wallet, so the mutex - held across safe-nonce assignment and the
// BuildBoundTx that assigns the executor EOA nonce - keeps each safe's internal
// nonce ordered consistently with the executor EOA nonce.
type walletSafePool struct {
	mu      sync.Mutex
	safes   []*safeEntry
	opCount uint64
	// needCreate forces the next iteration for this wallet to create a
	// replacement safe (skipping the recreate-rate draw) after a teardown sweep
	// dropped one.
	needCreate bool
	// drainSafe, when set, is an over-funded safe mid-teardown: it has already
	// forwarded a small balance to a sibling and the next iteration must drain
	// its remaining balance to the root wallet before dropping it.
	drainSafe *safeEntry
}

// usable reports whether the safe can currently be selected (created and not
// being torn down).
func (e *safeEntry) usable() bool {
	return !e.pending && !e.removing
}

// pickUsableLocked returns the next usable safe round-robin, advancing opCount,
// or nil if none is usable. Must be called with pool.mu held.
func (p *walletSafePool) pickUsableLocked() *safeEntry {
	n := len(p.safes)
	for i := 0; i < n; i++ {
		entry := p.safes[p.opCount%uint64(n)]
		p.opCount++
		if entry.usable() {
			return entry
		}
	}
	return nil
}

// usableCountLocked returns the number of currently usable safes. Must be called
// with pool.mu held.
func (p *walletSafePool) usableCountLocked() int {
	count := 0
	for _, e := range p.safes {
		if e.usable() {
			count++
		}
	}
	return count
}

// oldestUsableLocked returns the oldest usable safe (lowest slice index), or nil.
// Must be called with pool.mu held.
func (p *walletSafePool) oldestUsableLocked() *safeEntry {
	for _, e := range p.safes {
		if e.usable() {
			return e
		}
	}
	return nil
}

// lowestBalanceOtherLocked returns the usable safe (other than exclude) with the
// lowest tracked balance, or nil if there is no other usable safe. Used to pick
// the recipient for value transfers and teardown fund-forwarding so balances
// stay spread across the pool. Must be called with pool.mu held.
func (p *walletSafePool) lowestBalanceOtherLocked(exclude *safeEntry) *safeEntry {
	var best *safeEntry
	for _, e := range p.safes {
		if e == exclude || !e.usable() {
			continue
		}
		if best == nil || safeBalance(e).Cmp(safeBalance(best)) < 0 {
			best = e
		}
	}
	return best
}

// removeSafeLocked removes the given entry from the pool. Must be called with
// pool.mu held.
func (p *walletSafePool) removeSafeLocked(entry *safeEntry) {
	for i, e := range p.safes {
		if e == entry {
			p.safes = append(p.safes[:i], p.safes[i+1:]...)
			return
		}
	}
}

// safeBalance returns the entry's tracked balance, treating nil as zero.
func safeBalance(e *safeEntry) *big.Int {
	if e.balance == nil {
		return big.NewInt(0)
	}
	return e.balance
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	deploymentInfo *DeploymentInfo
	burnABI        *abi.ABI

	// safePools maps an executor wallet address -> *walletSafePool.
	safePools sync.Map
}

var ScenarioName = "safe-multisig"
var ScenarioDefaultOptions = ScenarioOptions{
	TotalCount:        0,
	Throughput:        10,
	MaxPending:        0,
	MaxWallets:        0,
	Rebroadcast:       1,
	BaseFee:           20,
	TipFee:            2,
	MinOwners:         1,
	MaxOwners:         5,
	Threshold:         0,
	SafesPerWallet:    3,
	ContractRatio:     0.5,
	RecreateRate:      0,
	BurnRounds:        1000,
	EoaValue:          0,
	FundingInterval:   32,
	GasLimit:          0,
	Timeout:           "",
	ClientGroup:       "",
	DeployClientGroup: "",
	LogTxs:            false,
}
var ScenarioDescriptor = scenario.Descriptor{
	Name:           ScenarioName,
	Aliases:        []string{"safe", "gnosis-safe"},
	Description:    "Create Safe multisigs with varying owner counts and drive EOA & contract-call execTransactions",
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
	flags.Uint64VarP(&s.options.TotalCount, "count", "c", ScenarioDefaultOptions.TotalCount, "Total number of execTransaction transactions to send")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t", ScenarioDefaultOptions.Throughput, "Number of execTransaction transactions to send per slot")
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", ScenarioDefaultOptions.Rebroadcast, "Enable reliable rebroadcast system")
	flags.Float64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use in transactions (in gwei)")
	flags.Float64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use in transactions (in gwei)")
	flags.StringVar(&s.options.BaseFeeWei, "basefee-wei", "", "Max fee per gas in wei (overrides --basefee for L2 sub-gwei fees)")
	flags.StringVar(&s.options.TipFeeWei, "tipfee-wei", "", "Max tip per gas in wei (overrides --tipfee for L2 sub-gwei fees)")
	flags.Uint64Var(&s.options.MinOwners, "min-owners", ScenarioDefaultOptions.MinOwners, "Minimum number of owners per safe (owner count is randomized per safe)")
	flags.Uint64Var(&s.options.MaxOwners, "max-owners", ScenarioDefaultOptions.MaxOwners, "Maximum number of owners per safe (clamped to the wallet pool size)")
	flags.Uint64Var(&s.options.Threshold, "threshold", ScenarioDefaultOptions.Threshold, "Signing threshold per safe (0 = n-of-n, otherwise clamped to the owner count)")
	flags.Uint64Var(&s.options.SafesPerWallet, "safes-per-wallet", ScenarioDefaultOptions.SafesPerWallet, "Number of safes to create per child wallet")
	flags.Float64Var(&s.options.ContractRatio, "contract-ratio", ScenarioDefaultOptions.ContractRatio, "Fraction of execTransactions that call the gas-burner contract (rest are EOA calls)")
	flags.Float64Var(&s.options.RecreateRate, "recreate-rate", ScenarioDefaultOptions.RecreateRate, "Probability (0..1) that an iteration re-creates a safe with a new random shape instead of executing a transaction (1 = creation-only, for state bloat)")
	flags.Uint64Var(&s.options.BurnRounds, "burn-rounds", ScenarioDefaultOptions.BurnRounds, "Upper bound on hashing rounds for the gas-burner contract (per-call gas is seed-derived up to this bound)")
	flags.Uint64Var(&s.options.EoaValue, "eoa-value", ScenarioDefaultOptions.EoaValue, "Value (in gwei) for EOA-call execTransactions (transfers target sibling safes); >0 keeps safes funded via low incremental top-ups")
	flags.Uint64Var(&s.options.FundingInterval, "funding-interval", ScenarioDefaultOptions.FundingInterval, "Interval (in slots) between safe balance top-up checks")
	flags.Uint64Var(&s.options.GasLimit, "gas-limit", ScenarioDefaultOptions.GasLimit, "Gas limit for execTransaction txs (0 = auto-compute per transaction)")
	flags.StringVar(&s.options.Timeout, "timeout", ScenarioDefaultOptions.Timeout, "Timeout for the scenario (e.g. '1h', '30m', '5s') - empty means no timeout")
	flags.StringVar(&s.options.ClientGroup, "client-group", ScenarioDefaultOptions.ClientGroup, "Client group to use for sending transactions")
	flags.StringVar(&s.options.DeployClientGroup, "deploy-client-group", ScenarioDefaultOptions.DeployClientGroup, "Client group to use for deployments")
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

	if s.options.SafesPerWallet == 0 {
		s.options.SafesPerWallet = 1
	}
	if s.options.MinOwners == 0 {
		s.options.MinOwners = 1
	}
	if s.options.MaxOwners < s.options.MinOwners {
		s.options.MaxOwners = s.options.MinOwners
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
		if s.options.Throughput*10 < 1000 {
			s.walletPool.SetWalletCount(s.options.Throughput * 10)
		} else {
			s.walletPool.SetWalletCount(1000)
		}
	}

	// deployer funds the shared infra deployment in one batch.
	s.walletPool.AddWellKnownWallet(&spamoor.WellKnownWalletConfig{
		Name:          "deployer",
		RefillAmount:  uint256.NewInt(2000000000000000000), // 2 ETH
		RefillBalance: uint256.NewInt(1000000000000000000), // 1 ETH
	})

	if s.options.TotalCount == 0 && s.options.Throughput == 0 {
		return fmt.Errorf("neither total count nor throughput limit set, must define at least one of them (see --help for list of all flags)")
	}

	return nil
}

func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Infof("starting scenario: %s", ScenarioName)
	defer s.logger.Infof("scenario %s finished.", ScenarioName)

	// deploy the shared Safe infrastructure (singleton, factory, gas burner).
	deploymentInfo, err := s.DeployContracts(ctx)
	if err != nil {
		s.logger.Errorf("could not deploy safe contracts: %v", err)
		return err
	}
	s.deploymentInfo = deploymentInfo
	s.logger.Infof("deployed singleton=%v factory=%v gasBurner=%v",
		deploymentInfo.SingletonAddr.Hex(), deploymentInfo.FactoryAddr.Hex(), deploymentInfo.GasBurnerAddr.Hex())

	burnABI, err := contract.GasBurnerMetaData.GetAbi()
	if err != nil {
		return fmt.Errorf("could not load gas burner abi: %w", err)
	}
	s.burnABI = burnABI

	// derive and create the per-wallet multisigs.
	if err := s.setupSafes(ctx); err != nil {
		s.logger.Errorf("could not set up safes: %v", err)
		return err
	}
	s.logger.Infof("set up %d safe multisigs", s.countSafes())

	// When EOA value transfers are enabled, keep the safes funded with small
	// incremental top-ups instead of a large lump sum. The top-up loop is driven
	// by its own context so it is torn down when Run returns (the parent ctx is
	// only cancelled by the caller after Run returns).
	if s.options.EoaValue > 0 {
		if err := s.topUpSafes(ctx); err != nil {
			return fmt.Errorf("could not fund safes: %w", err)
		}

		topUpCtx, topUpCancel := context.WithCancel(ctx)
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.runSafeTopUp(topUpCtx)
		}()
		defer func() {
			topUpCancel()
			wg.Wait()
		}()
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

	var timeout time.Duration
	if s.options.Timeout != "" {
		timeout, err = time.ParseDuration(s.options.Timeout)
		if err != nil {
			return fmt.Errorf("invalid timeout value: %v", err)
		}
		s.logger.Infof("Timeout set to %v", timeout)
	}

	err = scenario.RunTransactionScenario(ctx, scenario.TransactionScenarioOptions{
		TotalCount:                  s.options.TotalCount,
		Throughput:                  s.options.Throughput,
		MaxPending:                  maxPending,
		ThroughputIncrementInterval: 0,
		Timeout:                     timeout,
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
					logger.Warnf("could not send transaction: %v", err)
				} else if s.options.LogTxs {
					logger.Infof("sent tx #%6d: %v", params.TxIdx+1, tx.Hash().String())
				} else {
					logger.Debugf("sent tx #%6d: %v", params.TxIdx+1, tx.Hash().String())
				}
			})

			if _, err := receiptChan.Wait(ctx); err != nil {
				return err
			}

			return err
		},
	})

	return err
}

func (s *Scenario) sendTx(ctx context.Context, txIdx uint64) (scenario.ReceiptChan, *types.Transaction, *spamoor.Client, *spamoor.Wallet, error) {
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

	pool := s.poolForWallet(wallet.GetAddress())

	// In recreate mode the safes are grown and churned lazily; otherwise they are
	// created upfront and every iteration executes through an existing safe.
	if s.options.RecreateRate > 0 {
		return s.sendManagedTx(ctx, txIdx, client, wallet, pool, feeCap, tipCap)
	}
	return s.sendExecThrough(ctx, txIdx, client, wallet, pool, feeCap, tipCap)
}

// sendManagedTx implements the recreate-mode workload for one executor wallet:
// grow the pool up to --safes-per-wallet, then per iteration draw --recreate-rate
// to decide between tearing down the oldest safe and executing a transaction.
//
//   - below the safe-count limit: create a new safe (grow / replace after a drop);
//   - at the limit + churn draw, oldest safe funded: forward its balance to the
//     lowest-balance sibling (one teardown tx); the next iteration grows back;
//   - at the limit + churn draw, oldest safe empty: drop it and create the
//     replacement directly in this iteration (no teardown tx);
//   - otherwise: execute a transaction.
//
// At --recreate-rate 1 this is a pure safe-creation (state-bloat) workload.
func (s *Scenario) sendManagedTx(ctx context.Context, txIdx uint64, client *spamoor.Client, wallet *spamoor.Wallet, pool *walletSafePool, feeCap, tipCap *big.Int) (scenario.ReceiptChan, *types.Transaction, *spamoor.Client, *spamoor.Wallet, error) {
	target := int(s.options.SafesPerWallet)

	pool.mu.Lock()

	// An over-funded safe mid-teardown must drain its remainder to root next.
	if pool.drainSafe != nil {
		d := pool.drainSafe
		pool.drainSafe = nil
		pool.mu.Unlock()
		return s.drainSafeToRoot(ctx, txIdx, client, wallet, pool, d, feeCap, tipCap)
	}

	// A teardown dropped a safe; refill it now, skipping the recreate-rate draw.
	if pool.needCreate {
		pool.needCreate = false
		pool.mu.Unlock()
		return s.createSafe(ctx, txIdx, client, wallet, pool, feeCap, tipCap)
	}

	// Need at least one usable safe to execute through.
	if pool.usableCountLocked() == 0 {
		pool.mu.Unlock()
		return s.createSafe(ctx, txIdx, client, wallet, pool, feeCap, tipCap)
	}

	// Growth and teardown both happen only at the configured recreate rate;
	// otherwise execute a transaction through an existing safe.
	if !s.shouldRecreate() {
		pool.mu.Unlock()
		return s.sendExecThrough(ctx, txIdx, client, wallet, pool, feeCap, tipCap)
	}

	if len(pool.safes) < target {
		// Below the per-wallet limit: grow.
		pool.mu.Unlock()
		return s.createSafe(ctx, txIdx, client, wallet, pool, feeCap, tipCap)
	}

	// At the limit: tear down the oldest safe. If it holds funds, forward them to
	// the lowest-balance sibling via one execTransaction and let the next
	// iteration create the replacement; if it is empty, drop it and create the
	// replacement directly in this iteration.
	oldest := pool.oldestUsableLocked()
	if oldest != nil && s.options.EoaValue > 0 && safeBalance(oldest).Sign() > 0 {
		oldest.removing = true
		pool.mu.Unlock()
		return s.sweepAndTeardown(ctx, txIdx, client, wallet, pool, oldest, feeCap, tipCap)
	}
	if oldest != nil {
		pool.removeSafeLocked(oldest)
	}
	pool.mu.Unlock()
	return s.createSafe(ctx, txIdx, client, wallet, pool, feeCap, tipCap)
}

// sendExecThrough builds and submits a Safe execTransaction through one of the
// executor wallet's usable safes. The inner action is a gas-burner contract call
// or an EOA value transfer to a sibling safe; an EOA transfer is only used when
// it is funded (value transfers stay within the safe set so funds are not lost),
// otherwise it falls back to a contract call.
func (s *Scenario) sendExecThrough(ctx context.Context, txIdx uint64, client *spamoor.Client, wallet *spamoor.Wallet, pool *walletSafePool, feeCap, tipCap *big.Int) (scenario.ReceiptChan, *types.Transaction, *spamoor.Client, *spamoor.Wallet, error) {
	pool.mu.Lock()

	safe := pool.pickUsableLocked()
	if safe == nil {
		pool.mu.Unlock()
		return nil, nil, client, wallet, fmt.Errorf("no usable safe for wallet %s", wallet.GetAddress().Hex())
	}

	if err := s.resyncLocked(ctx, safe); err != nil {
		pool.mu.Unlock()
		return nil, nil, client, wallet, err
	}

	safeNonce := safe.nonce
	safe.nonce++

	// Decide the inner action. An EOA transfer normally targets a sibling safe
	// (so funds stay within the pool); when the source safe is over-funded it
	// drains to the root wallet instead. A value transfer is only used when the
	// source has enough tracked balance, otherwise it falls back to a gas-burner
	// contract call.
	isContract := s.isContractCall()
	value := big.NewInt(0)
	var recipient *safeEntry  // sibling safe to credit, nil when draining to root
	var toAddr common.Address // value-transfer destination
	if !isContract {
		if s.options.EoaValue > 0 {
			value = s.eoaCallValue()
		}
		switch {
		case s.options.EoaValue > 0 && safeBalance(safe).Cmp(value) < 0:
			isContract = true // unfunded
		case s.overFunded(safeBalance(safe)):
			toAddr = s.rootAddr() // drain excess to root
		default:
			recipient = pool.lowestBalanceOtherLocked(safe)
			if recipient == nil {
				isContract = true // no sibling to receive
			} else {
				toAddr = recipient.addr
			}
		}
	}

	var params *safeTxParams
	var burnRounds uint64
	if isContract {
		// Seed the burn deterministically from the tx index so the per-call gas
		// is reproducible across runs yet varies between calls.
		seed := new(big.Int).SetUint64(txIdx + 1)
		burnRounds = computeBurnRounds(seed, s.options.BurnRounds)
		data, err := s.burnABI.Pack("burn", seed, new(big.Int).SetUint64(s.options.BurnRounds))
		if err != nil {
			pool.mu.Unlock()
			return nil, nil, client, wallet, fmt.Errorf("could not pack burn call: %w", err)
		}
		params = &safeTxParams{
			To:        s.deploymentInfo.GasBurnerAddr,
			Value:     big.NewInt(0),
			Data:      data,
			Operation: 0,
			Nonce:     new(big.Int).SetUint64(safeNonce),
		}
	} else {
		params = &safeTxParams{
			To:        toAddr,
			Value:     value,
			Data:      nil,
			Operation: 0,
			Nonce:     new(big.Int).SetUint64(safeNonce),
		}
		// Optimistic balance bookkeeping; the top-up loop reconciles from chain.
		safe.balance = new(big.Int).Sub(safeBalance(safe), value)
		if recipient != nil {
			recipient.balance = new(big.Int).Add(safeBalance(recipient), value)
		}
	}

	tx, err := s.buildExecTx(ctx, wallet, safe, params, s.execGasLimit(isContract, safe.threshold, burnRounds, safeNonce == 0), feeCap, tipCap)
	if err != nil {
		pool.mu.Unlock()
		return nil, nil, client, wallet, err
	}
	pool.mu.Unlock()

	action := "eoa"
	if isContract {
		action = "burn"
	}
	receiptChan, err := s.submit(ctx, txIdx, client, wallet, tx, func(receipt *types.Receipt) {
		if receipt != nil && receipt.Status != types.ReceiptStatusSuccessful {
			// A reverted execTransaction did not advance the on-chain nonce and
			// did not move funds; flag for resync and roll back the optimistic
			// balance change.
			pool.mu.Lock()
			safe.needResync = true
			if !isContract && value.Sign() > 0 {
				safe.balance = new(big.Int).Add(safeBalance(safe), value)
				if recipient != nil {
					recipient.balance = new(big.Int).Sub(safeBalance(recipient), value)
				}
			}
			pool.mu.Unlock()
			s.logger.Warnf("execTransaction reverted (safe %s, action %v); will resync nonce", safe.addr.Hex(), action)
			return
		}
		if receipt != nil {
			txFees := utils.GetTransactionFees(tx, receipt)
			s.logger.WithField("rpc", client.GetName()).Debugf(
				" transaction %d confirmed in block #%v. action: %v owners: %d/%d fee: %v gwei logs: %v",
				txIdx+1, receipt.BlockNumber.String(), action, safe.threshold, len(safe.owners),
				txFees.TotalFeeGweiString(), len(receipt.Logs))
		}
	})
	if err != nil {
		return nil, nil, client, wallet, err
	}
	return receiptChan, tx, client, wallet, nil
}

// sweepAndTeardown forwards the full balance of the oldest safe to a sibling safe
// (the lowest-balance one) via a single execTransaction, then untracks the old
// safe on confirmation. The abandoned safe stays on-chain, contributing to state
// growth, while its funds are preserved within the pool.
func (s *Scenario) sweepAndTeardown(ctx context.Context, txIdx uint64, client *spamoor.Client, wallet *spamoor.Wallet, pool *walletSafePool, oldest *safeEntry, feeCap, tipCap *big.Int) (scenario.ReceiptChan, *types.Transaction, *spamoor.Client, *spamoor.Wallet, error) {
	// Sync the safe's actual current balance from chain before tearing down.
	onchain, err := client.GetEthClient().BalanceAt(ctx, oldest.addr, nil)
	if err != nil {
		s.clearRemoving(pool, oldest)
		return nil, nil, client, wallet, fmt.Errorf("could not read safe balance: %w", err)
	}

	pool.mu.Lock()
	if err := s.resyncLocked(ctx, oldest); err != nil {
		pool.mu.Unlock()
		s.clearRemoving(pool, oldest)
		return nil, nil, client, wallet, err
	}

	// Take the conservative minimum of the on-chain balance and the locally
	// tracked balance so an unconfirmed in/out transfer cannot make the sweep
	// exceed what the safe actually holds when it executes (the on-chain read
	// misses pending outgoing transfers; the tracked value misses pending
	// incoming ones).
	balance := onchain
	if tracked := safeBalance(oldest); tracked.Cmp(balance) < 0 {
		balance = tracked
	}
	oldest.balance = balance
	if balance.Sign() <= 0 {
		// Nothing to forward after all - drop and create the replacement.
		pool.removeSafeLocked(oldest)
		pool.mu.Unlock()
		return s.createSafe(ctx, txIdx, client, wallet, pool, feeCap, tipCap)
	}

	// Decide where the funds go. Normally they are forwarded to a sibling safe so
	// they stay within the pool. If the safe is over-funded (> 10x the funding
	// target), only a small balance goes to the sibling and the remainder is
	// drained back to the root wallet in a follow-up (a single execTransaction
	// targets one address). With no sibling at all, everything goes to root.
	sibling := pool.lowestBalanceOtherLocked(oldest)
	overFunded := s.overFunded(balance) && sibling != nil

	var recipient *safeEntry
	var toAddr common.Address
	value := balance
	if sibling != nil {
		recipient = sibling
		toAddr = sibling.addr
		if overFunded {
			// small balance to the sibling; the rest drains to root next.
			value = s.fundingTarget()
			if value.Cmp(balance) > 0 {
				value = balance
			}
		}
	} else {
		toAddr = s.rootAddr()
	}

	safeNonce := oldest.nonce
	oldest.nonce++
	params := &safeTxParams{
		To:        toAddr,
		Value:     value,
		Data:      nil,
		Operation: 0,
		Nonce:     new(big.Int).SetUint64(safeNonce),
	}
	tx, err := s.buildExecTx(ctx, wallet, oldest, params, s.execGasLimit(false, oldest.threshold, 0, safeNonce == 0), feeCap, tipCap)
	if err != nil {
		pool.mu.Unlock()
		s.clearRemoving(pool, oldest)
		return nil, nil, client, wallet, err
	}
	pool.mu.Unlock()

	receiptChan, err := s.submit(ctx, txIdx, client, wallet, tx, func(receipt *types.Receipt) {
		pool.mu.Lock()
		if receipt != nil && receipt.Status == types.ReceiptStatusSuccessful {
			if recipient != nil {
				recipient.balance = new(big.Int).Add(safeBalance(recipient), value)
			}
			if overFunded {
				// Remainder drains to root next iteration; keep the safe (removing).
				pool.drainSafe = oldest
			} else {
				pool.removeSafeLocked(oldest)
				pool.needCreate = true
			}
			pool.mu.Unlock()
			s.logger.WithField("rpc", client.GetName()).Debugf(
				" transaction %d confirmed in block #%v. action: teardown safe: %v -> %v (%v wei)",
				txIdx+1, receipt.BlockNumber.String(), oldest.addr.Hex(), toAddr.Hex(), value.String())
			return
		}
		// Sweep reverted - keep the safe and allow a retry.
		oldest.removing = false
		oldest.needResync = true
		pool.mu.Unlock()
		s.logger.Warnf("safe teardown sweep reverted (safe %s); will retry", oldest.addr.Hex())
	})
	if err != nil {
		s.clearRemoving(pool, oldest)
		return nil, nil, client, wallet, err
	}

	return receiptChan, tx, client, wallet, nil
}

// drainSafeToRoot completes the teardown of an over-funded safe by forwarding its
// remaining balance to the root wallet, then dropping it and scheduling a
// replacement. It runs only after the first teardown sweep has confirmed, so the
// balance it reads already excludes the small amount sent to the sibling.
func (s *Scenario) drainSafeToRoot(ctx context.Context, txIdx uint64, client *spamoor.Client, wallet *spamoor.Wallet, pool *walletSafePool, oldest *safeEntry, feeCap, tipCap *big.Int) (scenario.ReceiptChan, *types.Transaction, *spamoor.Client, *spamoor.Wallet, error) {
	balance, err := client.GetEthClient().BalanceAt(ctx, oldest.addr, nil)
	if err != nil {
		// Retry the drain on the next iteration.
		pool.mu.Lock()
		pool.drainSafe = oldest
		pool.mu.Unlock()
		return nil, nil, client, wallet, fmt.Errorf("could not read safe balance: %w", err)
	}
	if balance.Sign() == 0 {
		// Nothing left to drain - drop and create the replacement.
		pool.mu.Lock()
		pool.removeSafeLocked(oldest)
		pool.mu.Unlock()
		return s.createSafe(ctx, txIdx, client, wallet, pool, feeCap, tipCap)
	}

	rootAddr := s.rootAddr()
	pool.mu.Lock()
	if err := s.resyncLocked(ctx, oldest); err != nil {
		pool.drainSafe = oldest // retry the drain next iteration
		pool.mu.Unlock()
		return nil, nil, client, wallet, err
	}
	safeNonce := oldest.nonce
	oldest.nonce++
	params := &safeTxParams{
		To:        rootAddr,
		Value:     balance,
		Data:      nil,
		Operation: 0,
		Nonce:     new(big.Int).SetUint64(safeNonce),
	}
	tx, err := s.buildExecTx(ctx, wallet, oldest, params, s.execGasLimit(false, oldest.threshold, 0, safeNonce == 0), feeCap, tipCap)
	if err != nil {
		pool.drainSafe = oldest // retry the drain next iteration
		pool.mu.Unlock()
		return nil, nil, client, wallet, err
	}
	pool.mu.Unlock()

	receiptChan, err := s.submit(ctx, txIdx, client, wallet, tx, func(receipt *types.Receipt) {
		pool.mu.Lock()
		if receipt != nil && receipt.Status == types.ReceiptStatusSuccessful {
			pool.removeSafeLocked(oldest)
			pool.needCreate = true
			pool.mu.Unlock()
			s.logger.WithField("rpc", client.GetName()).Debugf(
				" transaction %d confirmed in block #%v. action: drain safe: %v -> root (%v wei)",
				txIdx+1, receipt.BlockNumber.String(), oldest.addr.Hex(), balance.String())
			return
		}
		// Drain reverted - retry on a later iteration.
		oldest.needResync = true
		pool.drainSafe = oldest
		pool.mu.Unlock()
		s.logger.Warnf("safe drain-to-root reverted (safe %s); will retry", oldest.addr.Hex())
	})
	if err != nil {
		pool.mu.Lock()
		pool.drainSafe = oldest
		pool.mu.Unlock()
		return nil, nil, client, wallet, err
	}

	return receiptChan, tx, client, wallet, nil
}

// createSafe creates a fresh safe with a new pseudo-random shape (owner set,
// threshold, salt) via SafeProxyFactory and appends it to the executor's pool.
// The safe is held pending until its creation transaction confirms, so it is not
// selected for execution before it exists on-chain.
func (s *Scenario) createSafe(ctx context.Context, txIdx uint64, client *spamoor.Client, wallet *spamoor.Wallet, pool *walletSafePool, feeCap, tipCap *big.Int) (scenario.ReceiptChan, *types.Transaction, *spamoor.Client, *spamoor.Wallet, error) {
	walletCount := s.walletPool.GetConfiguredWalletCount()

	// Pseudo-random shape + salt. Recreated safes are ephemeral churn, so they
	// need not be reproducible across restarts - a random salt also guarantees a
	// fresh, collision-free CREATE2 address.
	seed := make([]byte, 32)
	if _, err := rand.Read(seed); err != nil {
		return nil, nil, client, wallet, fmt.Errorf("could not generate safe seed: %w", err)
	}
	saltBytes := make([]byte, 32)
	if _, err := rand.Read(saltBytes); err != nil {
		return nil, nil, client, wallet, fmt.Errorf("could not generate safe salt: %w", err)
	}
	saltNonce := new(big.Int).SetBytes(saltBytes)

	spec, err := s.buildSafeSpecFromSeed(wallet, seed, saltNonce, walletCount)
	if err != nil {
		return nil, nil, client, wallet, err
	}

	initializer, err := s.encodeSetup(spec)
	if err != nil {
		return nil, nil, client, wallet, err
	}

	// createProxyWithNonce has no safe-nonce dependency, so gas estimation here is
	// reliable (a fresh-salt proxy creation always simulates cleanly).
	tx, err := wallet.BuildBoundTxWithEstimate(ctx, client, s.walletPool.GetTxPool(), &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Value:     uint256.NewInt(0),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		return s.deploymentInfo.Factory.CreateProxyWithNonce(transactOpts, s.deploymentInfo.SingletonAddr, initializer, spec.saltNonce)
	})
	if err != nil {
		return nil, nil, client, wallet, fmt.Errorf("could not build safe creation tx: %w", err)
	}

	instance, err := contract.NewSafe(spec.addr, client.GetEthClient())
	if err != nil {
		return nil, nil, client, wallet, fmt.Errorf("could not create safe instance: %w", err)
	}
	domainSep, err := computeDomainSeparator(s.walletPool.GetChainId(), spec.addr)
	if err != nil {
		return nil, nil, client, wallet, err
	}
	newEntry := &safeEntry{
		addr:      spec.addr,
		instance:  instance,
		owners:    spec.owners,
		threshold: spec.threshold,
		nonce:     0,
		domainSep: domainSep,
		pending:   true,
		balance:   big.NewInt(0),
	}

	pool.mu.Lock()
	pool.safes = append(pool.safes, newEntry)
	pool.opCount++
	pool.mu.Unlock()

	receiptChan, err := s.submit(ctx, txIdx, client, wallet, tx, func(receipt *types.Receipt) {
		pool.mu.Lock()
		if receipt != nil && receipt.Status == types.ReceiptStatusSuccessful {
			newEntry.pending = false
			pool.mu.Unlock()
			s.logger.WithField("rpc", client.GetName()).Debugf(
				" transaction %d confirmed in block #%v. action: create owners: %d/%d safe: %v",
				txIdx+1, receipt.BlockNumber.String(), spec.threshold, len(spec.owners), spec.addr.Hex())
			return
		}
		if receipt != nil {
			// Creation reverted - drop the dead entry so it is never selected.
			pool.removeSafeLocked(newEntry)
		}
		pool.mu.Unlock()
	})
	if err != nil {
		pool.mu.Lock()
		pool.removeSafeLocked(newEntry)
		pool.mu.Unlock()
		return nil, nil, client, wallet, err
	}
	return receiptChan, tx, client, wallet, nil
}

// resyncLocked re-reads the safe's nonce from chain when a prior transaction
// reverted (which left the on-chain nonce behind the local counter). Must be
// called with the owning pool's mutex held.
func (s *Scenario) resyncLocked(ctx context.Context, safe *safeEntry) error {
	if !safe.needResync {
		return nil
	}
	onchainNonce, err := safe.instance.Nonce(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("could not resync safe nonce: %w", err)
	}
	safe.nonce = onchainNonce.Uint64()
	safe.needResync = false
	return nil
}

// clearRemoving clears the removing flag on a safe (used when a teardown could
// not be submitted).
func (s *Scenario) clearRemoving(pool *walletSafePool, safe *safeEntry) {
	pool.mu.Lock()
	safe.removing = false
	pool.mu.Unlock()
}

// buildExecTx builds (without sending) a Safe execTransaction for params, signed
// by the safe's lowest-addressed threshold owners. Must be called with the
// owning pool's mutex held (it reads the safe's owners/domain separator).
func (s *Scenario) buildExecTx(ctx context.Context, wallet *spamoor.Wallet, safe *safeEntry, params *safeTxParams, gasLimit uint64, feeCap, tipCap *big.Int) (*types.Transaction, error) {
	safeTxHash, err := computeSafeTxHash(safe.domainSep, params)
	if err != nil {
		return nil, err
	}
	signatures, err := signSafeTx(safeTxHash, safe.owners[:safe.threshold])
	if err != nil {
		return nil, err
	}

	zeroAddr := common.Address{}
	return wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       gasLimit,
		Value:     uint256.NewInt(0),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		return safe.instance.ExecTransaction(transactOpts,
			params.To, params.Value, params.Data, params.Operation,
			big.NewInt(0), big.NewInt(0), big.NewInt(0),
			zeroAddr, zeroAddr, signatures)
	})
}

// submit sends tx asynchronously, wiring a receipt channel and an optional
// confirmation callback, and returns the channel the caller waits on.
func (s *Scenario) submit(ctx context.Context, txIdx uint64, client *spamoor.Client, wallet *spamoor.Wallet, tx *types.Transaction, onConfirm func(*types.Receipt)) (scenario.ReceiptChan, error) {
	receiptChan := make(scenario.ReceiptChan, 1)
	err := s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		ClientGroup: s.options.ClientGroup,
		Rebroadcast: s.options.Rebroadcast > 0,
		OnComplete: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			receiptChan <- receipt
		},
		OnConfirm: func(tx *types.Transaction, receipt *types.Receipt) {
			if onConfirm != nil {
				onConfirm(receipt)
			}
		},
		LogFn: spamoor.GetDefaultLogFn(s.logger, "", fmt.Sprintf("%6d", txIdx+1), tx),
	})
	if err != nil {
		wallet.MarkSkippedNonce(tx.Nonce())
		return nil, err
	}
	return receiptChan, nil
}

// safeRef is a snapshot reference to a tracked safe and its owning pool, used by
// the top-up loop so balance reads (RPC) happen outside the pool lock.
type safeRef struct {
	pool  *walletSafePool
	entry *safeEntry
	addr  common.Address
}

// snapshotSafes returns a reference to every currently tracked safe across all
// executor pools.
func (s *Scenario) snapshotSafes() []safeRef {
	var refs []safeRef
	s.safePools.Range(func(_, v any) bool {
		pool := v.(*walletSafePool)
		pool.mu.Lock()
		for _, e := range pool.safes {
			refs = append(refs, safeRef{pool: pool, entry: e, addr: e.addr})
		}
		pool.mu.Unlock()
		return true
	})
	return refs
}

// countSafes returns the total number of tracked safes across all pools.
func (s *Scenario) countSafes() int {
	count := 0
	s.safePools.Range(func(_, v any) bool {
		pool := v.(*walletSafePool)
		pool.mu.Lock()
		count += len(pool.safes)
		pool.mu.Unlock()
		return true
	})
	return count
}

// poolForWallet returns the safe pool for the given executor wallet, creating an
// empty one on first use.
func (s *Scenario) poolForWallet(executor common.Address) *walletSafePool {
	if v, ok := s.safePools.Load(executor); ok {
		return v.(*walletSafePool)
	}
	v, _ := s.safePools.LoadOrStore(executor, &walletSafePool{})
	return v.(*walletSafePool)
}

// deployClientGroup returns the client group to use for deployments, falling back
// to the transaction client group when unset.
func (s *Scenario) deployClientGroup() string {
	if s.options.DeployClientGroup != "" {
		return s.options.DeployClientGroup
	}
	return s.options.ClientGroup
}

// shouldRecreate draws whether this iteration should re-create a safe (with a
// new random shape) instead of executing a transaction, per --recreate-rate.
func (s *Scenario) shouldRecreate() bool {
	if s.options.RecreateRate <= 0 {
		return false
	}
	if s.options.RecreateRate >= 1 {
		return true
	}
	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		return false
	}
	return float64(n.Int64())/1_000_000 < s.options.RecreateRate
}

// isContractCall draws whether the next execTransaction targets the gas-burner
// contract (true) or performs an EOA call (false), per the configured ratio.
func (s *Scenario) isContractCall() bool {
	if s.options.ContractRatio >= 1 {
		return true
	}
	if s.options.ContractRatio <= 0 {
		return false
	}
	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		return false
	}
	return float64(n.Int64())/1_000_000 < s.options.ContractRatio
}

// eoaCallValue returns the value to send in an EOA-call execTransaction: zero
// when EoaValue is unset, otherwise a random amount in (0, EoaValue] gwei.
func (s *Scenario) eoaCallValue() *big.Int {
	if s.options.EoaValue == 0 {
		return big.NewInt(0)
	}
	max := gweiToWei(s.options.EoaValue)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return max
	}
	// shift into (0, max]
	return new(big.Int).Add(n, big.NewInt(1))
}

// Gas budgeting constants for execTransaction. They are intentionally generous -
// the limit only reserves block space (EIP-1559 charges gas actually used) - and
// are sized for the Amsterdam fee schedule's higher state-access costs, measured
// empirically against a Safe v1.4.1 on an Amsterdam devnet:
//   - each owner signature costs ~28k (cold owner SLOAD + ecrecover + linked-list
//     walk) under Amsterdam, so safeGasPerOwner is rounded up to 30k;
//   - the gas-burner loop costs ~229 gas/round, rounded up to 270 to cover the
//     64/63 gas forwarded across the Safe -> target call boundary.
const (
	safeExecBaseGas   = uint64(80000)  // intrinsic + nonce SSTORE + ExecutionSuccess event + base
	safeGasPerOwner   = uint64(30000)  // per verified signature
	firstExecStateGas = uint64(100000) // first execTransaction creates the nonce slot (Amsterdam state-creation surcharge)
	burnGasPerRound   = uint64(270)    // gas-burner hashing loop, per round
	burnCallOverhead  = uint64(15000)  // gas-burner call framing
	// eoaCallOverheadG covers the inner value-transfer call. The recipient is a
	// sibling Safe proxy whose receive() delegatecalls its singleton and emits an
	// event (~28k under Amsterdam), so this is sized well above a plain-EOA call.
	eoaCallOverheadG = uint64(50000)
)

// fundingTarget is the per-safe balance the top-up loop refills to (in wei):
// safeRefillTargetCalls worth of max-value transfers. Zero when value transfers
// are disabled.
func (s *Scenario) fundingTarget() *big.Int {
	return new(big.Int).Mul(gweiToWei(s.options.EoaValue), big.NewInt(safeRefillTargetCalls))
}

// overFunded reports whether a safe holds more than 10x the funding target, in
// which case value transfers / teardown sweeps drain to the root wallet instead
// of piling funds onto a sibling safe.
func (s *Scenario) overFunded(balance *big.Int) bool {
	if s.options.EoaValue == 0 {
		return false
	}
	return balance.Cmp(new(big.Int).Mul(s.fundingTarget(), big.NewInt(10))) > 0
}

// rootAddr returns the root wallet address (drain target for over-funded safes).
func (s *Scenario) rootAddr() common.Address {
	return s.walletPool.GetRootWallet().GetWallet().GetAddress()
}

// execGasLimit returns the gas limit for an execTransaction. When the gas-limit
// option is set it is used verbatim; otherwise it is computed from the signing
// threshold (each signature costs an ecrecover plus owner storage reads), the
// gas-burner round count (for contract calls), and whether this is the safe's
// first execTransaction - the first one writes the nonce slot 0->1, which carries
// the Amsterdam state-creation surcharge that later transactions do not pay. The
// result is clamped to the per-tx gas cap.
func (s *Scenario) execGasLimit(isContract bool, threshold int, burnRounds uint64, firstExec bool) uint64 {
	if s.options.GasLimit > 0 {
		return clampGas(s.options.GasLimit)
	}

	g := safeExecBaseGas + uint64(threshold)*safeGasPerOwner
	if firstExec {
		g += firstExecStateGas
	}
	if isContract {
		g += burnCallOverhead + burnRounds*burnGasPerRound
	} else {
		g += eoaCallOverheadG
	}
	return clampGas(g)
}

func clampGas(g uint64) uint64 {
	if g > utils.MaxGasLimitPerTx {
		return utils.MaxGasLimitPerTx
	}
	return g
}

// gweiToWei converts a gwei amount to wei.
func gweiToWei(gwei uint64) *big.Int {
	return new(big.Int).Mul(new(big.Int).SetUint64(gwei), big.NewInt(1_000_000_000))
}
