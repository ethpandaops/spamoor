package spamoor

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"sync"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/utils"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// WalletSelectionMode defines how wallets are selected from the pool.
type WalletSelectionMode uint8

var (
	// SelectWalletByIndex selects a wallet by index (modulo pool size).
	SelectWalletByIndex WalletSelectionMode = 0
	// SelectWalletRandom selects a random wallet from the pool.
	SelectWalletRandom WalletSelectionMode = 1
	// SelectWalletRoundRobin selects wallets in round-robin fashion.
	SelectWalletRoundRobin WalletSelectionMode = 2
	// SelectWalletByPendingTxCount selects a wallet by pending tx count (lowest pending tx count first).
	SelectWalletByPendingTxCount WalletSelectionMode = 3
)

// WalletPoolConfig contains configuration settings for the wallet pool including
// wallet count, funding amounts, and automatic refill behavior.
type WalletPoolConfig struct {
	WalletCount     uint64       `yaml:"wallet_count,omitempty"`
	RefillAmount    *uint256.Int `yaml:"refill_amount"`
	RefillBalance   *uint256.Int `yaml:"refill_balance"`
	RefillInterval  uint64       `yaml:"refill_interval"`
	WalletSeed      string       `yaml:"seed"`
	FundingGasLimit uint64       `yaml:"funding_gas_limit"`
	FeeStrategy     string       `yaml:"fee_strategy,omitempty"` // Fee calculation strategy: "" (legacy) or "adaptive"
}

// WellKnownWalletConfig defines configuration for a named wallet with custom funding settings.
// Well-known wallets have specific names and can have different refill amounts than regular wallets.
type WellKnownWalletConfig struct {
	Name          string
	RefillAmount  *uint256.Int
	RefillBalance *uint256.Int
	VeryWellKnown bool
	PrivateKey    string
}

// WalletPool manages a pool of child wallets derived from a root wallet with automatic funding
// and balance monitoring. It provides wallet selection strategies, automatic refills when balances
// drop below thresholds, and batch funding operations for efficiency.
type WalletPool struct {
	ctx         context.Context
	config      WalletPoolConfig
	logger      logrus.FieldLogger
	rootWallet  *RootWallet
	clientPool  *ClientPool
	txpool      *TxPool
	runFundings bool

	childWallets         []*Wallet
	wellKnownNames       []*WellKnownWalletConfig
	wellKnownWallets     map[string]*Wallet
	customWalletPrivkeys []string
	selectionMutex       sync.Mutex
	rrWalletIdx          int
	reclaimedFunds       bool

	// Optional callback to track transaction results for metrics
	transactionTracker func(err error)

	// Spammer ID for metrics tracking
	spammerID uint64

	// Low balance notification system
	lowBalanceNotifyChan chan struct{}
	lastFundingTime      time.Time

	// Stop funding loop signal
	stopFundingChan chan struct{}
	fundingLoopDone chan struct{}
}

// EIP-8037 spec constants reused by FundingGasFor / batcherGasFor here and by
// scenarios that need to size per-authorization / per-account state gas
// (e.g. scenarios/setcodetx).
const (
	// AccountCreationSize is the EIP-8037 charge unit count for creating a
	// new account. State gas = AccountCreationSize * cpsb.
	AccountCreationSize = 112
	// AuthorizationCreationSize is the EIP-8037 charge unit count applied to
	// each EIP-7702 SetCode authorization that creates a delegator account.
	AuthorizationCreationSize = 23
)

// FundingRequest represents a request to fund a wallet with a specific amount.
// Used internally for batch funding operations.
type FundingRequest struct {
	Wallet *Wallet
	Amount *uint256.Int
	// IsEmpty is true when the target wallet was observed as nonce==0 and
	// balance==0 at request-build time. Under EIP-8037/EIP-2780 a transfer
	// to an empty account additionally costs AccountCreationSize*cpsb state
	// gas, so the funding path must reserve more gas for it.
	IsEmpty bool
}

// GetDefaultWalletConfig returns default wallet pool configuration for a given scenario.
// It generates a random seed and sets reasonable defaults for refill amounts and intervals.
func GetDefaultWalletConfig(scenarioName string) *WalletPoolConfig {
	return &WalletPoolConfig{
		WalletSeed:     fmt.Sprintf("%v-%v", scenarioName, rand.Intn(1000000)),
		WalletCount:    0,
		RefillAmount:   uint256.NewInt(5000000000000000000),
		RefillBalance:  uint256.NewInt(1000000000000000000),
		RefillInterval: 600,
	}
}

// NewWalletPool creates a new wallet pool with the specified dependencies.
// The pool must be configured and prepared with PrepareWallets() before use.
func NewWalletPool(ctx context.Context, logger logrus.FieldLogger, rootWallet *RootWallet, clientPool *ClientPool, txpool *TxPool) *WalletPool {
	pool := &WalletPool{
		ctx:                  ctx,
		logger:               logger,
		rootWallet:           rootWallet,
		clientPool:           clientPool,
		txpool:               txpool,
		childWallets:         make([]*Wallet, 0),
		wellKnownWallets:     make(map[string]*Wallet),
		runFundings:          true,
		lowBalanceNotifyChan: make(chan struct{}, 100), // buffered channel
		stopFundingChan:      make(chan struct{}),
		fundingLoopDone:      make(chan struct{}),
	}

	txpool.RegisterWalletPool(pool)

	return pool
}

// GetContext returns the context associated with this wallet pool.
func (pool *WalletPool) GetContext() context.Context {
	return pool.ctx
}

// GetTxPool returns the transaction pool used by this wallet pool.
func (pool *WalletPool) GetTxPool() *TxPool {
	return pool.txpool
}

// GetSuggestedFees returns suggested fees using this pool's configured fee strategy.
func (pool *WalletPool) GetSuggestedFees(client *Client, baseFeeWei *big.Int, tipFeeWei *big.Int) (feeCap *big.Int, tipCap *big.Int, err error) {
	return pool.txpool.GetSuggestedFees(client, baseFeeWei, tipFeeWei, pool.config.FeeStrategy)
}

// SetFeeStrategy sets the fee calculation strategy for this wallet pool.
func (pool *WalletPool) SetFeeStrategy(strategy string) {
	pool.config.FeeStrategy = strategy
}

// GetClientPool returns the client pool used for blockchain interactions.
func (pool *WalletPool) GetClientPool() *ClientPool {
	return pool.clientPool
}

// GetRootWallet returns the root wallet that funds all child wallets.
func (pool *WalletPool) GetRootWallet() *RootWallet {
	return pool.rootWallet
}

// GetChainId returns the chain ID from the root wallet.
func (pool *WalletPool) GetChainId() *big.Int {
	return pool.rootWallet.wallet.GetChainId()
}

// LoadConfig loads wallet pool configuration from YAML string.
func (pool *WalletPool) LoadConfig(configYaml string) error {
	err := yaml.Unmarshal([]byte(configYaml), &pool.config)
	if err != nil {
		return err
	}

	return nil
}

// MarshalConfig returns the current configuration as a YAML string.
func (pool *WalletPool) MarshalConfig() (string, error) {
	yamlBytes, err := yaml.Marshal(&pool.config)
	if err != nil {
		return "", err
	}

	return string(yamlBytes), nil
}

// SetWalletCount sets the number of child wallets to create.
func (pool *WalletPool) SetWalletCount(count uint64) {
	pool.config.WalletCount = count
}

// SetWalletPrivkeys sets custom wallet private keys used for child wallets.
// The privkey parameter accepts hex strings with or without "0x" prefix.
// You should supply as many privkeys as you want to create wallets for (SetWalletCount),
// otherwise the remaining wallets will be generated deterministically based on the root wallet & seed.
func (pool *WalletPool) SetWalletPrivkeys(privkeys []string) {
	pool.customWalletPrivkeys = privkeys
}

// SetRunFundings enables or disables automatic wallet funding.
// When disabled, wallets will not be automatically refilled when their balance drops.
// This should be called before PrepareWallets(). To stop an already running funding loop,
// use StopFunding() instead.
func (pool *WalletPool) SetRunFundings(runFundings bool) {
	pool.runFundings = runFundings
}

// StopFunding stops the automatic wallet funding loop.
// This can be called after PrepareWallets() to stop the background funding goroutine.
// It blocks until the funding loop has fully stopped.
func (pool *WalletPool) StopFunding() {
	if !pool.runFundings {
		return
	}

	pool.runFundings = false

	select {
	case <-pool.stopFundingChan:
		// Already closed
	default:
		close(pool.stopFundingChan)
	}

	// Wait for the funding loop to finish
	select {
	case <-pool.fundingLoopDone:
	case <-pool.ctx.Done():
	}
}

// AddWellKnownWallet adds a named wallet with custom funding configuration.
// Well-known wallets are created alongside regular numbered wallets.
func (pool *WalletPool) AddWellKnownWallet(config *WellKnownWalletConfig) {
	pool.wellKnownNames = append(pool.wellKnownNames, config)
}

// SetRefillAmount sets the amount sent to wallets when they need funding.
func (pool *WalletPool) SetRefillAmount(amount *uint256.Int) {
	pool.config.RefillAmount = amount
}

// SetRefillBalance sets the balance threshold below which wallets are automatically refilled.
func (pool *WalletPool) SetRefillBalance(balance *uint256.Int) {
	pool.config.RefillBalance = balance
}

// SetWalletSeed sets the seed used for deterministic wallet generation.
// The same seed will always generate the same set of wallets.
func (pool *WalletPool) SetWalletSeed(seed string) {
	pool.config.WalletSeed = seed
}

// SetRefillInterval sets the interval in seconds between automatic balance checks.
func (pool *WalletPool) SetRefillInterval(interval uint64) {
	pool.config.RefillInterval = interval
}

// SetFundingGasLimit sets an explicit per-recipient gas limit for wallet
// funding transactions, overriding the automatic Amsterdam/EIP-8037 detection.
// Use 100000+ for L2 chains like Arbitrum/Optimism. When set to 0 (default),
// spamoor derives the limit from the current chain state: TxPool.MinIntrinsicGas
// (21,000 today; reduced once EIP-2780 ships), plus AccountCreationSize·cpsb
// for empty-target funding under EIP-8037.
func (pool *WalletPool) SetFundingGasLimit(gasLimit uint64) {
	pool.config.FundingGasLimit = gasLimit
}

// FundingGasFor returns the gas limit for a single direct EOA funding tx to a
// target with the given emptiness. Honors an explicit SetFundingGasLimit
// override when configured; otherwise computes from live chain state, which
// is guaranteed populated by InitializeBlockStats during startup.
//
// Base: TxPool.MinIntrinsicGas (21,000 today; will follow EIP-2780 reduction).
// Empty recipient on Amsterdam: + AccountCreationSize·cpsb (EIP-8037 state
// gas for new account creation). EIP-2780 will eventually replace the
// cpsb-based surcharge with a flat GAS_NEW_ACCOUNT (25,000) charge — handled
// at the same call site once the EL ships 2780.
func (pool *WalletPool) FundingGasFor(isEmpty bool) uint64 {
	if pool.config.FundingGasLimit != 0 {
		return pool.config.FundingGasLimit
	}
	base := pool.txpool.MinIntrinsicGas()
	if isEmpty {
		if cpsb := pool.txpool.GetCostPerStateByte(); cpsb > 0 {
			base += AccountCreationSize * cpsb
		}
	}
	return base
}

// GetFundingGasLimit returns the gas limit for a worst-case (empty target)
// direct funding transaction. Preserved for backwards compatibility with
// callers that need a single conservative number.
func (pool *WalletPool) GetFundingGasLimit() uint64 {
	return pool.FundingGasFor(true)
}

// batcherGasFor returns the per-recipient gas contribution for a single
// funding request inside the batcher contract. Honors an explicit
// SetFundingGasLimit override when it exceeds BatcherDefaultGasPerTx;
// otherwise computes from live chain state.
func (pool *WalletPool) batcherGasFor(req *FundingRequest) uint64 {
	if pool.config.FundingGasLimit > BatcherDefaultGasPerTx {
		return pool.config.FundingGasLimit
	}
	cpsb := pool.txpool.GetCostPerStateByte()
	if cpsb == 0 {
		return BatcherDefaultGasPerTx
	}
	base := uint64(callRegularGas)
	if req.IsEmpty {
		base += AccountCreationSize * cpsb
	}
	return base
}

// batcherBaseGas returns the per-batch overhead (tx intrinsic + batcher
// dispatch). Currently a flat constant; kept as a method so EIP-2780's
// reduced base intrinsic can be folded in here once activated.
func (pool *WalletPool) batcherBaseGas() uint64 {
	return BatcherBaseGas
}

// packFundingBatches greedily groups funding requests into batcher calls so
// that each batch stays under BatcherRPCGasCap. It also returns the total gas
// across all batches, used by the caller to reserve fee on the root wallet.
// With mixed empty/non-empty recipients this replaces the former fixed-size
// chunking, since per-recipient gas varies by a factor of ~10 on Amsterdam.
func (pool *WalletPool) packFundingBatches(reqs []*FundingRequest) ([][]*FundingRequest, uint64) {
	if len(reqs) == 0 {
		return nil, 0
	}
	baseGas := pool.batcherBaseGas()
	batches := make([][]*FundingRequest, 0, 1)
	var totalGas uint64
	curStart := 0
	curGas := baseGas
	for i, r := range reqs {
		reqGas := pool.batcherGasFor(r)
		// Flush the current batch if adding this request would exceed the cap,
		// but only if it already contains at least one request. Otherwise the
		// request goes in alone and we accept whatever over-cap gas it brings
		// (preserves the old "at least one recipient per tx" guarantee).
		if i > curStart && curGas+reqGas > BatcherRPCGasCap {
			batches = append(batches, reqs[curStart:i])
			totalGas += curGas
			curStart = i
			curGas = baseGas
		}
		curGas += reqGas
	}
	batches = append(batches, reqs[curStart:])
	totalGas += curGas
	return batches, totalGas
}

// EstimateDeployGas determines the gas limit for a contract-creation transaction
// with the given initcode (data). Primary path: eth_estimateGas via the provided
// client. Fallback path (RPC error or nil client): a formula-based upper bound
// sized for post-Amsterdam state-creation costs.
//
// Deployments are rare (typically one-time init steps), so the RPC round-trip
// is cheap. Static guesses under EIP-8037 frequently under-allocate because
// state-creation cost scales with both account creation (112·cpsb) and code
// deposit (per-byte cpsb), making the fallback intentionally conservative.
func (pool *WalletPool) EstimateDeployGas(ctx context.Context, client *Client, from common.Address, value *uint256.Int, data []byte) uint64 {
	var val *big.Int
	if value != nil {
		val = value.ToBig()
	}
	return estimateTxGas(ctx, client, from, nil, val, data, pool.txpool.IsAmsterdam(), pool.txpool.GetCostPerStateByte(), pool.txpool.MaxTxGas(), pool.logger)
}

// fallbackDeployGas returns a safe upper-bound gas limit for a contract
// deployment when eth_estimateGas is unavailable. Sized to cover worst-case
// code deposit (deployed == initcode) plus a generous execution heuristic and,
// on Amsterdam, the EIP-8037 state-gas components.
func fallbackDeployGas(initcodeLen int, isAmsterdam bool, cpsb uint64) uint64 {
	n := uint64(initcodeLen)
	out := uint64(100_000)
	// Rough per-byte budget covering initcode execution and code-deposit regular
	// gas. 200/byte covers the pre-Amsterdam deposit cost; post-Amsterdam the
	// deposit regular cost drops to ~0.2/byte, but we keep the margin for
	// initcode execution which static analysis cannot predict.
	out += n * 200
	// ABI-encoded calldata intrinsic (conservative avg of zero/nonzero bytes).
	out += n * 10
	if isAmsterdam {
		// Intrinsic: TxGas + CreateGasAmsterdam.
		out += 30_000
		// Intrinsic state gas for new account creation.
		out += AccountCreationSize * cpsb
		// Per-byte state gas for code deposit (upper bound: deployed == initcode).
		out += n * cpsb
	} else {
		// Pre-Amsterdam intrinsic for contract creation (53k includes CreateGas 32k).
		out += 50_000
	}
	// Cap at MaxTxGas so we never produce a tx rejected at EIP-7825's 16,777,216
	// regular-gas limit. Real deploys that would exceed this should specify an
	// explicit Gas value and rely on multi-tx factory patterns.
	if out > math.MaxUint32 {
		out = math.MaxUint32
	}
	return out
}

// estimateTxGas runs eth_estimateGas via the client and applies a defensive
// buffer on top of the result. Falls back to a size-based upper bound when
// estimation is unavailable. Handles both contract creation (to == nil) and
// call txs; the formula fallback only applies a deploy-style over-estimate
// when to == nil since per-call costs are highly contract-specific.
//
// maxTxGas is the effective per-tx gas cap for the chain (computed by the
// caller via TxPool.MaxTxGas, or set to utils.MaxGasLimitPerTx when no
// txpool context is available). Under EIP-8037 with the reservoir validator
// active, this can exceed 16.7M and approach the block gas limit so large
// contract deployments (Uniswap factory etc.) are not rejected.
//
// Package-level so both WalletPool (with live Amsterdam/cpsb state) and
// Wallet (with conservative defaults) can share the estimation logic.
// logger may be nil.
func estimateTxGas(ctx context.Context, client *Client, from common.Address, to *common.Address, value *big.Int, data []byte, isAmsterdam bool, cpsb uint64, maxTxGas uint64, logger logrus.FieldLogger) uint64 {
	if maxTxGas == 0 {
		maxTxGas = utils.MaxGasLimitPerTx
	}

	clamp := func(gas uint64) uint64 {
		if gas > maxTxGas {
			return maxTxGas
		}
		return gas
	}

	if client != nil {
		gas, err := client.EstimateGas(ctx, ethereum.CallMsg{
			From:  from,
			To:    to,
			Value: value,
			Data:  data,
		})
		if err == nil && gas > 0 {
			// If the raw estimate already exceeds the effective cap, no buffer
			// will make the tx includable — warn loudly so the operator knows
			// the contract/call is fundamentally too big rather than silently
			// shipping a tx that will be rejected.
			if gas > maxTxGas {
				if logger != nil {
					logger.WithFields(logrus.Fields{"estimate": gas, "cap": maxTxGas, "data_bytes": len(data), "to": to}).
						Warnf("tx gas estimate exceeds per-tx cap; contract is too large to deploy in a single tx under current chain params")
				}
				return maxTxGas
			}
			var buffered uint64
			if isAmsterdam {
				buffered = gas * 12 / 10
			} else {
				buffered = gas * 11 / 10
			}
			return clamp(buffered)
		}
		if err != nil && logger != nil {
			logger.WithError(err).Debugf("gas estimation via RPC failed; falling back to formula (data=%d bytes, to=%v)", len(data), to)
		}
	}

	// Fallback: for deploys we have a solid upper-bound formula; for calls
	// we can only apply a very rough one because per-call costs are contract-
	// specific. Return a generous cap so the tx at least goes through.
	if to == nil {
		return clamp(fallbackDeployGas(len(data), isAmsterdam, cpsb))
	}
	return clamp(1_000_000 + uint64(len(data))*200)
}

// SetTransactionTracker sets the optional callback to track transaction results for metrics.
func (pool *WalletPool) SetTransactionTracker(tracker func(err error)) {
	pool.transactionTracker = tracker
}

// GetTransactionTracker returns the transaction tracking callback if set.
func (pool *WalletPool) GetTransactionTracker() func(err error) {
	return pool.transactionTracker
}

// GetClient returns a client from the client pool using the specified options.
func (pool *WalletPool) GetClient(options ...ClientSelectionOption) *Client {
	return pool.clientPool.GetClient(options...)
}

// GetWallet returns a wallet from the pool using the specified selection strategy.
// Returns nil if no wallets are available.
func (pool *WalletPool) GetWallet(mode WalletSelectionMode, input int) *Wallet {
	pool.selectionMutex.Lock()
	defer pool.selectionMutex.Unlock()

	if len(pool.childWallets) == 0 {
		return nil
	}

	switch mode {
	case SelectWalletByIndex:
		input = input % len(pool.childWallets)
	case SelectWalletRandom:
		input = rand.Intn(len(pool.childWallets))
	case SelectWalletRoundRobin:
		input = pool.rrWalletIdx
		pool.rrWalletIdx++
		if pool.rrWalletIdx >= len(pool.childWallets) {
			pool.rrWalletIdx = 0
		}
	case SelectWalletByPendingTxCount:
		minPendingCount := uint64(math.MaxUint64)
		minPendingIndexes := []int{}
		for i, wallet := range pool.childWallets {
			pendingCount := wallet.GetNonce() - wallet.GetConfirmedNonce()
			if pendingCount < minPendingCount {
				minPendingCount = pendingCount
				minPendingIndexes = []int{i}
			} else if pendingCount == minPendingCount {
				minPendingIndexes = append(minPendingIndexes, i)
			}
		}
		input = input % len(minPendingIndexes)
		input = minPendingIndexes[input]
	}
	return pool.childWallets[input]
}

// GetWellKnownWallet returns a well-known wallet by name.
// Returns nil if the wallet doesn't exist.
func (pool *WalletPool) GetWellKnownWallet(name string) *Wallet {
	return pool.wellKnownWallets[name]
}

// GetVeryWellKnownWalletAddress derives the address of a "very well known" wallet
// without registering it. Very well known wallets are derived only from the root
// wallet's private key and the wallet name, without any scenario seed.
// This makes them consistent across different scenario runs.
func (pool *WalletPool) GetVeryWellKnownWalletAddress(name string) common.Address {
	idxBytes := make([]byte, len(name))
	copy(idxBytes, name)
	// VeryWellKnown wallets don't use the seed, so we skip adding it

	parentKey := crypto.FromECDSA(pool.rootWallet.wallet.privkey)
	childKey := sha256.Sum256(append(parentKey, idxBytes...))

	// Derive private key and then address
	privateKey, err := crypto.HexToECDSA(fmt.Sprintf("%x", childKey))
	if err != nil {
		return common.Address{}
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}
	}

	return crypto.PubkeyToAddress(*publicKeyECDSA)
}

// GetWalletName returns a human-readable name for the given wallet address.
// Returns "root" for the root wallet, numbered names for child wallets,
// custom names for well-known wallets, or "unknown" if not found.
func (pool *WalletPool) GetWalletName(address common.Address) string {
	if pool.rootWallet.wallet.GetAddress() == address {
		return "root"
	}

	for i, wallet := range pool.childWallets {
		if wallet.GetAddress() == address {
			return fmt.Sprintf("%d", i+1)
		}
	}

	for _, config := range pool.wellKnownNames {
		wallet := pool.wellKnownWallets[config.Name]
		if wallet != nil && wallet.GetAddress() == address {
			return config.Name
		}
	}

	return "unknown"
}

// GetAllWallets returns a slice containing all wallets (well-known and child wallets).
// The root wallet is not included in this list.
func (pool *WalletPool) GetAllWallets() []*Wallet {
	if pool.childWallets == nil {
		return []*Wallet{}
	}

	wallets := make([]*Wallet, 0, len(pool.childWallets)+len(pool.wellKnownNames))
	wallets = append(wallets, pool.childWallets...)
	for _, config := range pool.wellKnownNames {
		if pool.wellKnownWallets[config.Name] != nil {
			wallets = append(wallets, pool.wellKnownWallets[config.Name])
		}
	}

	return wallets
}

// GetConfiguredWalletCount returns the configured number of child wallets.
func (pool *WalletPool) GetConfiguredWalletCount() uint64 {
	return pool.config.WalletCount
}

// GetWalletCount returns the actual number of child wallets created.
func (pool *WalletPool) GetWalletCount() uint64 {
	return uint64(len(pool.childWallets))
}

// PrepareWallets creates all configured wallets and funds them if needed.
// It generates deterministic wallets based on the root wallet and seed,
// then funds any wallets below the refill threshold. Also starts the
// automatic balance monitoring if funding is enabled.
func (pool *WalletPool) PrepareWallets() error {
	if len(pool.childWallets) > 0 {
		return nil
	}

	seed := pool.config.WalletSeed

	if pool.config.WalletCount == 0 && len(pool.wellKnownWallets) == 0 {
		pool.childWallets = make([]*Wallet, 0)
	} else {
		var client *Client
		var fundingReqs []*FundingRequest
		var childWallets []*Wallet

		for i := 0; i < 3; i++ {
			client = pool.clientPool.GetClient(WithClientSelectionMode(SelectClientRandom), WithoutBuilder()) // send all preparation transactions via this client to avoid rejections due to nonces
			if client == nil {
				return fmt.Errorf("no client available")
			}

			childWallets = make([]*Wallet, pool.config.WalletCount)
			fundingReqs = make([]*FundingRequest, 0, pool.config.WalletCount)

			var walletErr error
			wg := &sync.WaitGroup{}
			wl := make(chan bool, 50)
			walletsMutex := &sync.Mutex{}

			for _, config := range pool.wellKnownNames {
				wg.Add(1)
				wl <- true
				go func(config *WellKnownWalletConfig) {
					defer func() {
						<-wl
						wg.Done()
					}()
					if walletErr != nil {
						return
					}

					childWallet, fundingReq, err := pool.prepareWellKnownWallet(config, client, seed)
					if err != nil {
						pool.logger.Errorf("could not prepare well known wallet %v: %v", config.Name, err)
						walletErr = err
						return
					}

					walletsMutex.Lock()
					if fundingReq != nil {
						fundingReqs = append(fundingReqs, fundingReq)
					}
					pool.wellKnownWallets[config.Name] = childWallet
					walletsMutex.Unlock()
				}(config)
			}

			if walletErr != nil {
				return fmt.Errorf("well known wallet preparation failed: %w", walletErr)
			}

			for childIdx := uint64(0); childIdx < pool.config.WalletCount; childIdx++ {
				wg.Add(1)
				wl <- true
				go func(childIdx uint64) {
					defer func() {
						<-wl
						wg.Done()
					}()
					if walletErr != nil {
						return
					}

					childWallet, fundingReq, err := pool.prepareChildWallet(childIdx, client, seed)
					if err != nil {
						pool.logger.Errorf("could not prepare child wallet %v: %v", childIdx, err)
						walletErr = err
						return
					}

					walletsMutex.Lock()
					childWallets[childIdx] = childWallet
					if fundingReq != nil {
						fundingReqs = append(fundingReqs, fundingReq)
					}
					walletsMutex.Unlock()
				}(childIdx)
			}
			wg.Wait()

			if walletErr != nil {
				return fmt.Errorf("child wallet preparation failed: %w", walletErr)
			}

			if len(childWallets) > 0 {
				break
			}
		}

		pool.childWallets = childWallets

		if pool.runFundings && len(fundingReqs) > 0 {
			err := pool.processFundingRequests(fundingReqs)
			if err != nil {
				return err
			}
		}

		for childIdx, childWallet := range pool.childWallets {
			pool.logger.Debugf(
				"initialized child wallet %4d (addr: %v, balance: %v ETH, nonce: %v)",
				childIdx,
				childWallet.GetAddress().String(),
				utils.WeiToEther(uint256.MustFromBig(childWallet.GetBalance())).Uint64(),
				childWallet.GetNonce(),
			)
		}

		pool.logger.Infof("initialized %v child wallets", pool.config.WalletCount)
	}

	// watch wallet balances
	if pool.runFundings {
		go pool.watchWalletBalancesLoop()
	}

	return nil
}

// prepareChildWallet creates a child wallet derived from the root wallet using deterministic key generation.
// It generates a private key by hashing the root wallet's private key with the child index and seed.
// Returns the wallet, funding request (if needed), and any error.
func (pool *WalletPool) prepareChildWallet(childIdx uint64, client *Client, seed string) (*Wallet, *FundingRequest, error) {
	walletPrivkey := ""
	if childIdx < uint64(len(pool.customWalletPrivkeys)) {
		walletPrivkey = pool.customWalletPrivkeys[childIdx]
	} else {
		idxBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(idxBytes, childIdx)
		if seed != "" {
			seedBytes := []byte(seed)
			idxBytes = append(idxBytes, seedBytes...)
		}
		parentKey := crypto.FromECDSA(pool.rootWallet.wallet.privkey)
		childKey := sha256.Sum256(append(parentKey, idxBytes...))
		walletPrivkey = fmt.Sprintf("%x", childKey)
	}

	return pool.prepareWallet(walletPrivkey, client, pool.config.RefillAmount, pool.config.RefillBalance)
}

// prepareWellKnownWallet creates a named wallet derived from the root wallet using deterministic key generation.
// It generates a private key by hashing the root wallet's private key with the wallet name and seed.
// If a private key is provided in the config, it uses that instead of generating a new one.
// Uses custom refill amounts from the config if specified, otherwise falls back to pool defaults.
func (pool *WalletPool) prepareWellKnownWallet(config *WellKnownWalletConfig, client *Client, seed string) (*Wallet, *FundingRequest, error) {
	walletPrivkey := config.PrivateKey
	if walletPrivkey == "" {
		idxBytes := make([]byte, len(config.Name))
		copy(idxBytes, config.Name)
		if seed != "" && !config.VeryWellKnown {
			seedBytes := []byte(seed)
			idxBytes = append(idxBytes, seedBytes...)
		}
		parentKey := crypto.FromECDSA(pool.rootWallet.wallet.privkey)
		childKey := sha256.Sum256(append(parentKey, idxBytes...))
		walletPrivkey = fmt.Sprintf("%x", childKey)
	}

	refillAmount := pool.config.RefillAmount
	refillBalance := pool.config.RefillBalance

	if config.RefillAmount != nil {
		refillAmount = config.RefillAmount
	}
	if config.RefillBalance != nil {
		refillBalance = config.RefillBalance
	}

	return pool.prepareWallet(walletPrivkey, client, refillAmount, refillBalance)
}

// calculateFundingAmount calculates the amount to fund a wallet.
// It returns the maximum of refillAmount and (refillBalance - currentBalance),
// ensuring the wallet reaches at least the refillBalance threshold.
func calculateFundingAmount(currentBalance *uint256.Int, refillAmount *uint256.Int, refillBalance *uint256.Int) *uint256.Int {
	fundingAmount := new(uint256.Int).Set(refillAmount)

	// If refillBalance > currentBalance, calculate the difference
	if refillBalance.Cmp(currentBalance) > 0 {
		neededAmount := new(uint256.Int).Sub(refillBalance, currentBalance)
		if neededAmount.Cmp(fundingAmount) > 0 {
			fundingAmount = neededAmount
		}
	}

	return fundingAmount
}

// prepareWallet creates a wallet from a private key and checks if it needs funding.
// Updates the wallet's state from the blockchain and creates a funding request if
// the wallet's balance is below the specified refill threshold.
func (pool *WalletPool) prepareWallet(privkey string, client *Client, refillAmount *uint256.Int, refillBalance *uint256.Int) (*Wallet, *FundingRequest, error) {
	var err error

	privateKey, address, err := LoadPrivateKey(privkey)
	if err != nil {
		return nil, nil, err
	}

	childWallet := pool.txpool.GetRegisteredWallet(address)
	if childWallet == nil {
		childWallet = NewWallet(privateKey, address)

		childWallet = pool.txpool.RegisterWallet(childWallet, pool.ctx)
		if childWallet == nil {
			// RegisterWallet returns nil when the context is already cancelled.
			// Surface that as an error instead of dereferencing the nil wallet
			// in the UpdateWallet call below.
			if ctxErr := pool.ctx.Err(); ctxErr != nil {
				return nil, nil, fmt.Errorf("wallet registration cancelled: %w", ctxErr)
			}

			return nil, nil, fmt.Errorf("wallet registration failed for %s", address.Hex())
		}
	}

	err = client.UpdateWallet(pool.ctx, childWallet)
	if err != nil {
		return nil, nil, err
	}

	// Set up low balance notification
	if pool.runFundings && refillBalance != nil {
		childWallet.setLowBalanceNotification(pool.lowBalanceNotifyChan, refillBalance.ToBig())
	}

	var fundingReq *FundingRequest
	if pool.runFundings && refillBalance != nil && childWallet.GetBalance().Cmp(refillBalance.ToBig()) < 0 {
		currentBalance := uint256.MustFromBig(childWallet.GetBalance())
		fundingReq = &FundingRequest{
			Wallet:  childWallet,
			Amount:  calculateFundingAmount(currentBalance, refillAmount, refillBalance),
			IsEmpty: isWalletEmpty(childWallet),
		}
	}
	return childWallet, fundingReq, nil
}

// isWalletEmpty reports whether the wallet's current on-chain state qualifies
// as "empty" per EIP-2780 / EIP-8037: nonce == 0 and balance == 0. We skip the
// code check because spamoor wallets are always EOAs without deployed code;
// treating a coded account as empty would only over-allocate gas, which has
// no on-chain cost on success.
func isWalletEmpty(w *Wallet) bool {
	return w.GetNonce() == 0 && w.GetBalance().Sign() == 0
}

// watchWalletBalancesLoop runs continuously to monitor and refill wallet balances.
// It periodically checks all wallets and funds those below the refill threshold.
// Also listens for low balance notifications to trigger immediate checks.
// Exits when the context is cancelled, funds have been reclaimed, or StopFunding is called.
func (pool *WalletPool) watchWalletBalancesLoop() {
	defer close(pool.fundingLoopDone)

	sleepTime := time.Duration(pool.config.RefillInterval) * time.Second
	timer := time.NewTimer(sleepTime)
	defer timer.Stop()

	// Aggregation timer for notifications
	var aggregationTimer *time.Timer
	aggregationChan := make(chan struct{})

	resupplyWallets := func() {
		if pool.reclaimedFunds {
			return
		}

		if time.Since(pool.lastFundingTime) < 30*time.Second {
			return
		}

		pool.lastFundingTime = time.Now()
		if aggregationTimer != nil {
			aggregationTimer.Stop()
			aggregationTimer = nil
		}

		err := pool.resupplyChildWallets()
		if err != nil {
			pool.logger.Warnf("could not check & resupply child wallets: %v", err)
			sleepTime = 1 * time.Minute
		} else {
			sleepTime = time.Duration(pool.config.RefillInterval) * time.Second
		}
		timer.Reset(sleepTime)
	}

	for {
		select {
		case <-pool.ctx.Done():
			return

		case <-pool.stopFundingChan:
			pool.logger.Debugf("funding loop stopped")
			return

		case <-timer.C:
			// Regular interval check
			resupplyWallets()

		case <-pool.lowBalanceNotifyChan:
			// Low balance notification received
			if aggregationTimer == nil {
				// First notification, start aggregation timer
				aggregationTimer = time.AfterFunc(30*time.Second, func() {
					select {
					case aggregationChan <- struct{}{}:
					default:
					}
				})
			}
			// Drain any additional notifications to avoid blocking
		drainloop:
			for {
				select {
				case <-pool.lowBalanceNotifyChan:
				default:
					break drainloop
				}
			}

		case <-aggregationChan:
			// Aggregation period ended, trigger funding
			resupplyWallets()
		}
	}
}

// resupplyChildWallets checks all wallets and creates funding requests for those below the refill threshold.
// It updates wallet states from the blockchain and processes any needed funding requests in batch.
func (pool *WalletPool) resupplyChildWallets() error {
	client := pool.clientPool.GetClient(WithClientSelectionMode(SelectClientRandom))
	if client == nil {
		return fmt.Errorf("no client available")
	}

	err := client.UpdateWallet(pool.ctx, pool.rootWallet.wallet)
	if err != nil {
		return err
	}

	var walletErr error
	wg := &sync.WaitGroup{}
	wl := make(chan bool, 50)

	fundingReqs := make([]*FundingRequest, 0, len(pool.childWallets)+len(pool.wellKnownWallets))
	reqsMutex := &sync.Mutex{}

	for idx, config := range pool.wellKnownNames {
		wellKnownWallet := pool.wellKnownWallets[config.Name]
		if wellKnownWallet == nil {
			continue
		}

		wg.Add(1)
		wl <- true
		go func(idx int, childWallet *Wallet, config *WellKnownWalletConfig) {
			defer func() {
				<-wl
				wg.Done()
			}()
			if walletErr != nil {
				return
			}

			refillAmount := pool.config.RefillAmount
			refillBalance := pool.config.RefillBalance

			if config.RefillAmount != nil {
				refillAmount = config.RefillAmount
			}
			if config.RefillBalance != nil {
				refillBalance = config.RefillBalance
			}

			err := client.UpdateWallet(pool.ctx, childWallet)
			if err != nil {
				walletErr = err
				return
			}

			if refillBalance != nil && childWallet.GetBalance().Cmp(refillBalance.ToBig()) < 0 {
				currentBalance := uint256.MustFromBig(childWallet.GetBalance())
				reqsMutex.Lock()
				fundingReqs = append(fundingReqs, &FundingRequest{
					Wallet:  childWallet,
					Amount:  calculateFundingAmount(currentBalance, refillAmount, refillBalance),
					IsEmpty: isWalletEmpty(childWallet),
				})
				reqsMutex.Unlock()
			}
		}(idx, wellKnownWallet, config)
	}

	for childIdx := 0; childIdx < len(pool.childWallets); childIdx++ {
		wg.Add(1)
		wl <- true
		go func(childIdx int) {
			defer func() {
				<-wl
				wg.Done()
			}()
			if walletErr != nil {
				return
			}

			childWallet := pool.childWallets[childIdx]
			err := client.UpdateWallet(pool.ctx, childWallet)
			if err != nil {
				walletErr = err
				return
			}
			if pool.config.RefillBalance != nil && childWallet.GetBalance().Cmp(pool.config.RefillBalance.ToBig()) < 0 {
				currentBalance := uint256.MustFromBig(childWallet.GetBalance())
				reqsMutex.Lock()
				fundingReqs = append(fundingReqs, &FundingRequest{
					Wallet:  childWallet,
					Amount:  calculateFundingAmount(currentBalance, pool.config.RefillAmount, pool.config.RefillBalance),
					IsEmpty: isWalletEmpty(childWallet),
				})
				reqsMutex.Unlock()
			}
		}(childIdx)
	}
	wg.Wait()
	if walletErr != nil {
		return walletErr
	}

	if len(fundingReqs) > 0 {
		err := pool.processFundingRequests(fundingReqs)
		if err != nil {
			return err
		}
	} else {
		pool.logger.Infof("checked child wallets (no funding needed)")
	}

	return nil
}

// processFundingRequests handles a batch of funding requests by creating and sending transactions.
// It can use either individual transactions or batch transactions via the batcher contract for efficiency.
// Processes transactions in chunks to avoid overwhelming the network.
// It now checks if the root wallet has sufficient balance before attempting to send funding transactions.
func (pool *WalletPool) processFundingRequests(fundingReqs []*FundingRequest) error {
	client := pool.clientPool.GetClient(WithClientSelectionMode(SelectClientRandom))
	if client == nil {
		return fmt.Errorf("no client available")
	}

	// Calculate total funding amount needed
	totalFundingAmount := uint256.NewInt(0)
	for _, req := range fundingReqs {
		totalFundingAmount = totalFundingAmount.Add(totalFundingAmount, req.Amount)
	}

	feeCap, _, err := client.GetSuggestedFee(pool.ctx)
	if err != nil {
		return err
	}

	batcher := pool.rootWallet.GetTxBatcher()
	feeAmount := big.NewInt(0).Set(feeCap)

	// Partition requests into batches that each fit under the RPC gas cap.
	// batchTxCount is the number of outer transactions we'll submit, and
	// totalGas is the sum of all their gas limits (used to reserve fees on
	// the root wallet).
	var batches [][]*FundingRequest
	var totalGas uint64
	if batcher != nil {
		batches, totalGas = pool.packFundingBatches(fundingReqs)
	} else {
		// Individual funding: one tx per request, gas sized per recipient.
		batches = make([][]*FundingRequest, 0, len(fundingReqs))
		for _, req := range fundingReqs {
			batches = append(batches, []*FundingRequest{req})
			totalGas += pool.FundingGasFor(req.IsEmpty)
		}
	}
	batchTxCount := len(batches)
	feeAmount = big.NewInt(0).Mul(feeAmount, big.NewInt(int64(totalGas)))

	totalFundingAmount = totalFundingAmount.Add(totalFundingAmount, uint256.MustFromBig(feeAmount))

	return pool.rootWallet.WithWalletLock(pool.ctx, batchTxCount, totalFundingAmount, pool.clientPool, func(reason string) {
		pool.logger.Infof("root wallet is locked, %s", reason)
	}, func() error {
		if batcher != nil {
			err := batcher.Deploy(pool.ctx, pool.rootWallet.wallet, client)
			if err != nil {
				return fmt.Errorf("failed to deploy batcher: %v", err)
			}
		}

		txList := make([]*types.Transaction, 0, batchTxCount)
		batchTxMap := map[common.Hash][]*FundingRequest{}
		if batcher != nil {
			for _, batch := range batches {
				tx, err := pool.buildWalletFundingBatchTx(batch, client, batcher)
				if err != nil {
					return err
				}
				txList = append(txList, tx)
				batchTxMap[tx.Hash()] = batch
			}
		} else {
			for _, req := range fundingReqs {
				tx, err := pool.buildWalletFundingTx(req, client)
				if err != nil {
					return err
				}
				txList = append(txList, tx)
			}
		}

		receipts, err := pool.txpool.SendTransactionBatch(pool.ctx, pool.rootWallet.wallet, txList, &BatchOptions{
			SendTransactionOptions: SendTransactionOptions{
				Client: client,
			},
			MaxRetries:   3,
			PendingLimit: 200,
			LogFn: func(confirmedCount int, totalCount int) {
				pool.logger.Infof("funding child wallets... (%v/%v)", confirmedCount, totalCount)
			},
			LogInterval: 200,
		})
		if err != nil {
			return fmt.Errorf("could not send funding txs: %w", err)
		}

		for _, receipt := range receipts {
			if receipt != nil && receipt.Status == types.ReceiptStatusSuccessful {
				batch, ok := batchTxMap[receipt.TxHash]
				if ok {
					for _, req := range batch {
						req.Wallet.AddBalance(req.Amount.ToBig())
					}
				}
			}
		}

		return nil
	})
}

// buildWalletFundingTx creates a direct-EOA transaction to fund a single
// wallet. The per-tx gas limit is sized via FundingGasFor(req.IsEmpty), so
// non-empty refills use the cheap 21k path while initial fundings to empty
// targets reserve enough for the EIP-2780/EIP-8037 state gas charge.
func (pool *WalletPool) buildWalletFundingTx(req *FundingRequest, client *Client) (*types.Transaction, error) {
	if client == nil {
		client = pool.clientPool.GetClient(WithClientSelectionMode(SelectClientByIndex, 0))
		if client == nil {
			return nil, fmt.Errorf("no client available")
		}
	}
	feeCap, tipCap, err := client.GetSuggestedFee(pool.ctx)
	if err != nil {
		return nil, err
	}
	if feeCap.Cmp(big.NewInt(100000000000)) > 0 {
		feeCap = big.NewInt(100000000000) // 100 gwei
	} else {
		feeCap = feeCap.Add(feeCap, big.NewInt(2000000000)) // +2 gwei
	}
	if tipCap.Cmp(big.NewInt(100000000000)) > 0 {
		tipCap = big.NewInt(100000000000) // 100 gwei
	} else {
		tipCap = tipCap.Add(tipCap, big.NewInt(100000000)) // +0.1 gwei
	}
	// Ensure tipCap never exceeds feeCap (EIP-1559 requirement)
	if tipCap.Cmp(feeCap) > 0 {
		tipCap = new(big.Int).Set(feeCap)
	}

	toAddr := req.Wallet.GetAddress()
	refillTx, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       pool.FundingGasFor(req.IsEmpty),
		To:        &toAddr,
		Value:     req.Amount,
	})
	if err != nil {
		return nil, err
	}
	tx, err := pool.rootWallet.wallet.BuildDynamicFeeTx(refillTx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// buildWalletFundingBatchTx creates a transaction to fund multiple wallets using the batcher contract.
// It calculates the total amount needed, encodes the funding requests as calldata,
// and builds a transaction to the batcher contract with appropriate gas limits.
func (pool *WalletPool) buildWalletFundingBatchTx(requests []*FundingRequest, client *Client, batcher *TxBatcher) (*types.Transaction, error) {
	if client == nil {
		client = pool.clientPool.GetClient(WithClientSelectionMode(SelectClientByIndex, 0))
		if client == nil {
			return nil, fmt.Errorf("no client available")
		}
	}
	feeCap, tipCap, err := client.GetSuggestedFee(pool.ctx)
	if err != nil {
		return nil, err
	}
	if feeCap.Cmp(big.NewInt(100000000000)) > 0 {
		feeCap = big.NewInt(100000000000) // 100 gwei
	} else {
		feeCap = feeCap.Add(feeCap, big.NewInt(2000000000)) // +2 gwei
	}
	if tipCap.Cmp(big.NewInt(100000000000)) > 0 {
		tipCap = big.NewInt(100000000000) // 100 gwei
	} else {
		tipCap = tipCap.Add(tipCap, big.NewInt(100000000)) // +0.1 gwei
	}
	// Ensure tipCap never exceeds feeCap (EIP-1559 requirement)
	if tipCap.Cmp(feeCap) > 0 {
		tipCap = new(big.Int).Set(feeCap)
	}

	totalAmount := uint256.NewInt(0)
	gasTotal := pool.batcherBaseGas()
	for _, req := range requests {
		totalAmount = totalAmount.Add(totalAmount, req.Amount)
		gasTotal += pool.batcherGasFor(req)
	}

	batchData, err := batcher.GetRequestCalldata(requests)
	if err != nil {
		return nil, err
	}

	toAddr := batcher.GetAddress()
	refillTx, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       gasTotal,
		To:        &toAddr,
		Value:     totalAmount,
		Data:      batchData,
	})
	if err != nil {
		return nil, err
	}
	tx, err := pool.rootWallet.wallet.BuildDynamicFeeTx(refillTx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// reclaimTx holds a reclaim transaction and its associated wallet for fund recovery operations.
type reclaimTx struct {
	tx     *types.Transaction
	wallet *Wallet
}

// buildWalletReclaimTx creates a transaction to reclaim funds from a child wallet back to the target address.
// It calculates the maximum amount that can be reclaimed after accounting for transaction fees.
// Returns nil if the wallet doesn't have enough balance to cover fees.
func (pool *WalletPool) buildWalletReclaimTx(ctx context.Context, childWallet *Wallet, client *Client, target common.Address) (*types.Transaction, error) {
	if client == nil {
		client = pool.clientPool.GetClient(WithClientSelectionMode(SelectClientByIndex, 0))
		if client == nil {
			return nil, fmt.Errorf("no client available")
		}
	}
	feeCap, tipCap, err := client.GetSuggestedFee(ctx)
	if err != nil {
		return nil, err
	}
	if feeCap.Cmp(big.NewInt(100000000000)) > 0 {
		feeCap = big.NewInt(100000000000) // 100 gwei
	} else {
		feeCap = feeCap.Add(feeCap, big.NewInt(1000000000)) // +1 gwei
	}
	if tipCap.Cmp(feeCap) < 0 {
		// set tip to the same as fee, this ensures the wallet balance is cleared completely
		tipCap = feeCap
	}

	gasLimit := pool.GetFundingGasLimit()
	feeAmount := big.NewInt(0).Mul(tipCap, big.NewInt(int64(gasLimit)))
	reclaimAmount := big.NewInt(0).Sub(childWallet.GetBalance(), feeAmount)

	if reclaimAmount.Cmp(big.NewInt(0)) <= 0 {
		return nil, nil
	}

	reclaimTx, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       gasLimit,
		To:        &target,
		Value:     uint256.MustFromBig(reclaimAmount),
	})
	if err != nil {
		return nil, err
	}

	tx, err := childWallet.BuildDynamicFeeTx(reclaimTx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// CheckChildWalletBalance checks and refills a specific wallet if needed.
// This can be used to manually trigger funding for a single wallet.
func (pool *WalletPool) CheckChildWalletBalance(childWallet *Wallet) error {
	client := pool.clientPool.GetClient(WithClientSelectionMode(SelectClientRandom))
	if client == nil {
		return fmt.Errorf("no client available")
	}

	balance, err := client.GetBalanceAt(pool.ctx, childWallet.GetAddress())
	if err != nil {
		return err
	}
	childWallet.SetBalance(balance)

	refillAmount := pool.config.RefillAmount
	refillBalance := pool.config.RefillBalance

	for _, config := range pool.wellKnownNames {
		wellKnownWallet := pool.wellKnownWallets[config.Name]
		if wellKnownWallet != nil && wellKnownWallet.GetAddress().String() == childWallet.GetAddress().String() {
			if config.RefillAmount != nil {
				refillAmount = config.RefillAmount
			}
			if config.RefillBalance != nil {
				refillBalance = config.RefillBalance
			}
		}
	}

	if refillBalance == nil || childWallet.GetBalance().Cmp(refillBalance.ToBig()) >= 0 {
		return nil
	}

	currentBalance := uint256.MustFromBig(childWallet.GetBalance())

	return pool.processFundingRequests([]*FundingRequest{
		{
			Wallet: childWallet,
			Amount: calculateFundingAmount(currentBalance, refillAmount, refillBalance),
		},
	})
}

// ReclaimFunds reclaims all funds from child wallets back to the root wallet.
// This is typically called when shutting down to consolidate remaining funds.
// After calling this, automatic funding is disabled.
func (pool *WalletPool) ReclaimFunds(ctx context.Context, client *Client) error {
	pool.reclaimedFunds = true

	if client == nil {
		client = pool.clientPool.GetClient(WithClientSelectionMode(SelectClientRandom))
	}
	if client == nil {
		return fmt.Errorf("no client available")
	}

	reclaimMtx := sync.Mutex{}
	reclaimTxs := []*reclaimTx{}
	reclaimWg := sync.WaitGroup{}
	reclaimChan := make(chan struct{}, 100)

	reclaimWallet := func(wallet *Wallet) {
		reclaimWg.Add(1)
		reclaimChan <- struct{}{}

		go func() {
			defer func() {
				<-reclaimChan
				reclaimWg.Done()
			}()

			err := client.UpdateWallet(ctx, wallet)
			if err != nil {
				return
			}

			balance := wallet.GetBalance()
			if balance.Cmp(big.NewInt(0)) == 0 {
				return
			}

			tx, err := pool.buildWalletReclaimTx(ctx, wallet, client, pool.rootWallet.wallet.GetAddress())
			if err != nil || tx == nil {
				return
			}

			reclaimMtx.Lock()
			reclaimTxs = append(reclaimTxs, &reclaimTx{
				tx:     tx,
				wallet: wallet,
			})
			reclaimMtx.Unlock()
		}()
	}

	for _, wallet := range pool.childWallets {
		reclaimWallet(wallet)
	}
	for _, wallet := range pool.wellKnownWallets {
		reclaimWallet(wallet)
	}
	reclaimWg.Wait()

	if len(reclaimTxs) > 0 {
		wg := sync.WaitGroup{}
		wg.Add(len(reclaimTxs))
		for _, tx := range reclaimTxs {
			go func(tx *reclaimTx) {
				pool.logger.Infof("sending reclaim tx %v (%v)", tx.tx.Hash().String(), utils.ReadableAmount(uint256.MustFromBig(tx.tx.Value())))
				pool.txpool.SendTransaction(ctx, tx.wallet, tx.tx, &SendTransactionOptions{
					Client: client,
					OnComplete: func(_ *types.Transaction, receipt *types.Receipt, err error) {
						wg.Done()
						if err != nil {
							pool.logger.Warnf("reclaim tx %v failed: %v", tx.tx.Hash().String(), err)
						}
					},
				})
			}(tx)
		}
		wg.Wait()
	}

	return nil
}

// SetSpammerID sets the spammer ID for metrics tracking
func (pool *WalletPool) SetSpammerID(spammerID uint64) {
	pool.spammerID = spammerID
}

// GetSpammerID returns the spammer ID for metrics tracking
func (pool *WalletPool) GetSpammerID() uint64 {
	return pool.spammerID
}
