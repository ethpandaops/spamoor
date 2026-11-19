package spamoor

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"

	"github.com/ethpandaops/spamoor/txtypes"
	"github.com/ethpandaops/spamoor/utils"
)

// CostPerStateByte is the EIP-8037 cost_per_state_byte value spamoor uses
// when the Amsterdam fee model is active. State-gas charges (account
// creation, code deposit, authorization delegators) multiply this by the
// per-charge size constants below.
const CostPerStateByte uint64 = 1530

// feeWindowSize is the number of recent blocks tracked for dynamic fee headroom calculation.
const feeWindowSize = 10

// BlockInfo represents information about a processed block including
// hash, parent hash, gas limit, and timestamp for chain reorganization detection.
type BlockInfo struct {
	Number     uint64
	Hash       common.Hash
	ParentHash common.Hash
	Timestamp  uint64
	GasLimit   uint64
	SlotNumber uint64
}

// BlockWithHash represents a block with its hash.
type BlockWithHash struct {
	Hash  common.Hash
	Block *txtypes.Block
}

// TxInfo represents information about a confirmed transaction including
// the transaction details, associated wallets, and send options.
type TxInfo struct {
	TxHash     common.Hash
	From       common.Address
	To         *common.Address
	Tx         *txtypes.Transaction
	TxFees     *utils.TxFees
	FromWallet *Wallet
	ToWallet   *Wallet
	Options    *SendTransactionOptions

	receipt    *txtypes.Receipt
	walletPool *WalletPool
	tracked    *trackedTx
}

// BlockStatsCallback is called when a block is processed with statistics for a specific wallet pool
type BlockStatsCallback func(blockNumber uint64, walletPoolStats *WalletPoolBlockStats)

// BulkBlockStatsCallback is called when a block is processed with statistics for ALL wallet pools
type BulkBlockStatsCallback func(blockNumber uint64, globalBlockStats *GlobalBlockStats)

// WalletPoolBlockStats contains transaction statistics for a wallet pool in a specific block
type WalletPoolBlockStats struct {
	ConfirmedTxCount uint64
	TotalTxFees      *big.Int
	AffectedWallets  int
	Block            *txtypes.Block
	ConfirmedTxs     []*TxInfo
	Receipts         []*txtypes.Receipt
}

type GlobalBlockStats struct {
	Block           *txtypes.Block
	Receipts        []*txtypes.Receipt
	WalletPoolStats map[*WalletPool]*WalletPoolBlockStats
}

// BlockSubscription represents a subscription to block updates for a specific wallet pool
type BlockSubscription struct {
	ID         uint64
	WalletPool *WalletPool
	Callback   BlockStatsCallback
}

// BulkBlockSubscription represents a subscription to block updates for ALL wallet pools
type BulkBlockSubscription struct {
	ID       uint64
	Callback BulkBlockStatsCallback
}

// TxPool manages transaction submission, confirmation tracking, and chain reorganization handling.
// It monitors blockchain blocks, tracks transaction confirmations, handles reorgs by re-submitting
// affected transactions, and provides transaction awaiting functionality with automatic rebroadcasting.
type TxPool struct {
	options          *TxPoolOptions
	processStaleChan chan uint64
	lastBlockNumber  uint64
	slotDuration     atomic.Int64

	// wallet and wallet pool tracking
	walletsMutex sync.RWMutex
	wallets      map[common.Address]*txPoolWalletRegistration
	walletPools  map[*WalletPool]struct{}

	// Block tracking for reorg detection
	blocksMutex sync.RWMutex
	blocks      map[uint64]*BlockInfo
	reorgDepth  int // Number of blocks to keep in memory for reorg tracking

	// Transaction tracking for reorg recovery
	txsMutex     sync.RWMutex
	confirmedTxs map[uint64][]*TxInfo

	// Current block gas limit tracking
	currentGasLimit uint64
	currentBaseFee  *big.Int
	blockStatsMutex sync.RWMutex

	// Frame transaction capability, probed from the predeploys the EIPs install.
	frameSupport frameSupportState

	// preAmsterdamFeeModel selects the legacy (pre-EIP-8037) gas/fee model
	// when true. Defaults to false: the Amsterdam fee model is active.
	preAmsterdamFeeModel bool

	// Recent block tracking for dynamic fee calculation.
	// Sliding window of gas usage and baseFee to adjust feeCap:
	// - Gas fill ratio controls headroom (how much above baseFee)
	// - Average baseFee creates lag that self-corrects rising baseFee
	recentGasUsed  [feeWindowSize]uint64
	recentGasLimit [feeWindowSize]uint64
	recentBaseFees [feeWindowSize]*big.Int
	recentGasIdx   int
	recentGasLen   int

	// Block update subscriptions
	subscriptionsMutex sync.RWMutex
	subscriptions      map[uint64]*BlockSubscription
	bulkSubscriptions  map[uint64]*BulkBlockSubscription
	nextSubscriptionID atomic.Uint64

	// Hash watches for transactions whose sender the pool does not track, keyed by
	// transaction hash. hashWatcherCount lets block processing skip the lock entirely
	// in the common case of no watches.
	hashWatchersMutex sync.Mutex
	hashWatchers      map[common.Hash]func(*txtypes.Transaction, *txtypes.Receipt)
	hashWatcherCount  atomic.Int64

	// Pending transaction tracking, see txtracker.go.
	tracker      txTracker
	recoverySem  chan struct{}
	rebroadcasts workQueue

	// Block that failed to process and the attempts made at it, see retryBlock.
	retryBlockNumber uint64
	retryBlockCount  int
}

type txPoolWalletRegistration struct {
	ctxs   map[context.Context]struct{}
	wallet *Wallet
}

// TxPoolOptions contains configuration options for the transaction pool.
type TxPoolOptions struct {
	Context             context.Context
	Logger              *logrus.Entry
	ClientPool          *ClientPool
	ReorgDepth          int // Number of blocks to keep in memory for reorg tracking
	ChainId             *big.Int
	ExternalBlockSource *ExternalBlockSource
	// SlotDuration is the slot duration used while the chain reports no EIP-7843
	// slot numbers (pre-Amsterdam chains, L2s). Defaults to 12s.
	SlotDuration time.Duration
	// PreAmsterdamFeeModel forces the legacy (pre-EIP-8037) gas/fee model.
	// When false (default) spamoor uses the Amsterdam fee model, which is
	// safe to keep on non-Amsterdam chains because its gas budgets are
	// strictly higher than the legacy ones.
	PreAmsterdamFeeModel bool
}

type ExternalBlockSource struct {
	SubscribeBlocks func(ctx context.Context, capacity int) chan *ExternalBlockEvent
}

type ExternalBlockEvent struct {
	Number   uint64
	Clients  []*Client
	Block    *BlockWithHash
	Receipts []*txtypes.Receipt
}

// NewTxPool creates a new transaction pool with the specified options.
// It starts background goroutines for block processing and stale transaction handling.
// The pool automatically begins monitoring the blockchain for new blocks and managing
// transaction confirmations and reorgs.
func NewTxPool(options *TxPoolOptions) *TxPool {
	pool := &TxPool{
		options:              options,
		wallets:              make(map[common.Address]*txPoolWalletRegistration),
		walletPools:          make(map[*WalletPool]struct{}),
		processStaleChan:     make(chan uint64, 1),
		recoverySem:          make(chan struct{}, recoveryConcurrency),
		blocks:               map[uint64]*BlockInfo{},
		confirmedTxs:         map[uint64][]*TxInfo{},
		reorgDepth:           10, // Default value
		subscriptions:        map[uint64]*BlockSubscription{},
		bulkSubscriptions:    map[uint64]*BulkBlockSubscription{},
		preAmsterdamFeeModel: options.PreAmsterdamFeeModel,
	}

	if options.Context == nil {
		options.Context = context.Background()
	}

	if options.ReorgDepth > 0 {
		pool.reorgDepth = options.ReorgDepth
	}

	if options.SlotDuration <= 0 {
		options.SlotDuration = 12 * time.Second
	}

	// Block stats (gas limit, base fee) are loaded by the caller via
	// InitializeBlockStats before any scenario starts. Construction itself
	// does not perform RPC calls.

	go pool.runTxPoolLoop()
	go pool.processStaleTransactionsLoop()

	return pool
}

// RegisterWallet registers a wallet with the transaction pool, keyed by address.
// Multiple scenarios (e.g. a spammer that was paused and restarted with the same seed)
// may register the same address; each registering context is tracked, and the
// registration is auto-removed only once ALL of its contexts have been cancelled. This
// ensures a restarted spammer's re-registration keeps the wallet alive instead of it
// being evicted when the old, cancelled context's cleanup fires.
// WatchTransaction registers a one-shot callback fired when a transaction with the given
// hash is included in a processed block.
//
// It is for transactions whose sender the pool does not track through a wallet, such as a
// keyless EIP-8141 contract sender: block processing matches confirmations to wallets by
// their from-address, so a transaction from an unregistered address would otherwise be
// invisible. The watch is removed when it fires or when ctx is cancelled, so callers pass
// a bounded context to keep an unincluded transaction from lingering.
func (pool *TxPool) WatchTransaction(ctx context.Context, txHash common.Hash, callback func(*txtypes.Transaction, *txtypes.Receipt)) {
	pool.hashWatchersMutex.Lock()
	if pool.hashWatchers == nil {
		pool.hashWatchers = make(map[common.Hash]func(*txtypes.Transaction, *txtypes.Receipt))
	}

	if _, exists := pool.hashWatchers[txHash]; !exists {
		pool.hashWatcherCount.Add(1)
	}

	pool.hashWatchers[txHash] = callback
	pool.hashWatchersMutex.Unlock()

	context.AfterFunc(ctx, func() {
		pool.hashWatchersMutex.Lock()
		if _, ok := pool.hashWatchers[txHash]; ok {
			delete(pool.hashWatchers, txHash)
			pool.hashWatcherCount.Add(-1)
		}
		pool.hashWatchersMutex.Unlock()
	})
}

// fireHashWatcher invokes and removes the watch for a transaction hash, if one is set.
func (pool *TxPool) fireHashWatcher(txHash common.Hash, tx *txtypes.Transaction, receipt *txtypes.Receipt) {
	if pool.hashWatcherCount.Load() == 0 {
		return
	}

	pool.hashWatchersMutex.Lock()
	callback := pool.hashWatchers[txHash]
	if callback != nil {
		delete(pool.hashWatchers, txHash)
		pool.hashWatcherCount.Add(-1)
	}
	pool.hashWatchersMutex.Unlock()

	if callback != nil {
		callback(tx, receipt)
	}
}

func (pool *TxPool) RegisterWallet(wallet *Wallet, ctx context.Context) *Wallet {
	if ctx.Err() != nil {
		return nil
	}

	pool.walletsMutex.Lock()
	registration, ok := pool.wallets[wallet.address]
	if ok {
		if registration.wallet.privkey == nil && wallet.privkey != nil {
			registration.wallet.privkey = wallet.privkey
		}
		if _, exists := registration.ctxs[ctx]; exists {
			regWallet := registration.wallet
			pool.walletsMutex.Unlock()
			return regWallet
		}
		registration.ctxs[ctx] = struct{}{}
	} else {
		registration = &txPoolWalletRegistration{
			ctxs:   map[context.Context]struct{}{ctx: {}},
			wallet: wallet,
		}
		pool.wallets[wallet.address] = registration
	}
	regWallet := registration.wallet
	pool.walletsMutex.Unlock()

	// Auto-cleanup when this context is cancelled. The identity check (reg == registration)
	// guards against a later re-registration having replaced the entry.
	context.AfterFunc(ctx, func() {
		pool.walletsMutex.Lock()
		if reg, ok := pool.wallets[wallet.address]; ok && reg == registration {
			delete(reg.ctxs, ctx)
			if len(reg.ctxs) == 0 {
				delete(pool.wallets, wallet.address)
			}
		}
		pool.walletsMutex.Unlock()
	})

	return regWallet
}

// RegisterWalletPool registers a wallet pool with the transaction pool.
// It is used to track the wallet pools that are active and need to be processed.
// The pool is auto-removed when its scenario context is cancelled (e.g. on pause or
// restart), so a restarted spammer's stale pool does not linger and cause its
// transactions to be double-attributed across two pools sharing the same spammer id.
func (pool *TxPool) RegisterWalletPool(walletPool *WalletPool) {
	if walletPool.ctx.Err() != nil {
		return
	}

	pool.walletsMutex.Lock()
	pool.walletPools[walletPool] = struct{}{}
	pool.walletsMutex.Unlock()

	context.AfterFunc(walletPool.ctx, func() {
		pool.walletsMutex.Lock()
		delete(pool.walletPools, walletPool)
		pool.walletsMutex.Unlock()
	})
}

// GetRegisteredWallet returns a wallet by address.
func (pool *TxPool) GetRegisteredWallet(address common.Address) *Wallet {
	pool.walletsMutex.Lock()
	defer pool.walletsMutex.Unlock()

	registration, found := pool.wallets[address]
	if found {
		// Lazily prune cancelled contexts (belt-and-suspenders alongside the AfterFunc
		// cleanup registered in RegisterWallet).
		for c := range registration.ctxs {
			if c.Err() != nil {
				delete(registration.ctxs, c)
			}
		}
		if len(registration.ctxs) == 0 {
			delete(pool.wallets, address)
			return nil
		}
		return registration.wallet
	}

	return nil
}

// runTxPoolLoop continuously monitors for new blocks and processes them.
// It tracks the highest block number across all clients and processes new blocks
// sequentially. Also triggers stale transaction processing when new blocks arrive.
// Runs until the context is cancelled and recovers from panics with logging.
func (pool *TxPool) runTxPoolLoop() {
	// recover only stops an unwinding panic when the recovering function is
	// deferred directly, so utils.RecoverPanic must not be wrapped in a closure.
	defer utils.RecoverPanic(pool.options.Logger, "TxPool.runTxPoolLoop", pool.runTxPoolLoop)

	if pool.options.ExternalBlockSource != nil {
		pool.runExternalBlockSourceLoop()
		return
	}

	highestBlockNumber := uint64(0)
	for {
		newHighestBlockNumber, clients := pool.getHighestBlockNumber()
		if newHighestBlockNumber > highestBlockNumber {
			// Skip processing historical blocks on startup unless blockchain is young (< 10 blocks)
			// This prevents processing millions of blocks when connecting to a long running chain
			if highestBlockNumber == 0 && newHighestBlockNumber > 10 {
				highestBlockNumber = newHighestBlockNumber - 1
			}

			for blockNumber := highestBlockNumber + 1; blockNumber <= newHighestBlockNumber; blockNumber++ {
				processedBlock := false
				for _, client := range clients {
					err := pool.processBlock(pool.options.Context, client, blockNumber, nil, nil)
					if err != nil {
						logrus.WithField("client", client.GetName()).WithError(err).Errorf("error processing block %v", blockNumber)
						continue
					}

					highestBlockNumber = blockNumber
					processedBlock = true
					break
				}

				if !processedBlock {
					if !pool.retryBlock(blockNumber) {
						break
					}

					highestBlockNumber = blockNumber
				}
			}

			select {
			case <-pool.options.Context.Done():
				return
			case pool.processStaleChan <- highestBlockNumber:
			default:
			}
		}

		select {
		case <-pool.options.Context.Done():
			return
		case <-time.After(3 * time.Second):
		}
	}
}

func (pool *TxPool) runExternalBlockSourceLoop() {
	// recover only stops an unwinding panic when the recovering function is
	// deferred directly, so utils.RecoverPanic must not be wrapped in a closure.
	defer utils.RecoverPanic(pool.options.Logger, "TxPool.runExternalBlockSourceLoop", pool.runExternalBlockSourceLoop)

	blockChan := pool.options.ExternalBlockSource.SubscribeBlocks(pool.options.Context, 10)
	highestBlockNumber := uint64(0)
	for {
		select {
		case <-pool.options.Context.Done():
			return
		case blockEvent := <-blockChan:
			for _, client := range blockEvent.Clients {
				client.recordBlockHeight(blockEvent.Number)
			}

			if blockEvent.Number > highestBlockNumber {
				// Skip processing historical blocks on startup unless blockchain is young (< 10 blocks)
				// This prevents processing millions of blocks when connecting to a long running chain
				if highestBlockNumber == 0 && blockEvent.Number > 10 {
					highestBlockNumber = blockEvent.Number - 1
				}

				for blockNumber := highestBlockNumber + 1; blockNumber <= blockEvent.Number; blockNumber++ {
					processedBlock := false
					for _, client := range blockEvent.Clients {
						var blockWithHash *BlockWithHash
						if blockNumber == blockEvent.Number {
							blockWithHash = blockEvent.Block
						}
						err := pool.processBlock(pool.options.Context, client, blockNumber, blockWithHash, nil)
						if err != nil {
							logrus.WithField("client", client.GetName()).WithError(err).Errorf("error processing block %v", blockNumber)
							continue
						}

						highestBlockNumber = blockNumber
						processedBlock = true
						break
					}

					if !processedBlock {
						if !pool.retryBlock(blockNumber) {
							break
						}

						highestBlockNumber = blockNumber
					}
				}

				select {
				case <-pool.options.Context.Done():
					return
				case pool.processStaleChan <- highestBlockNumber:
				default:
				}
			}
		}
	}
}

// processStaleTransactionsLoop handles stale transaction confirmation checking.
// It listens for block number updates and processes stale confirmations for all
// active wallets. Runs until the context is cancelled and recovers from panics.
func (pool *TxPool) processStaleTransactionsLoop() {
	// recover only stops an unwinding panic when the recovering function is
	// deferred directly, so utils.RecoverPanic must not be wrapped in a closure.
	defer utils.RecoverPanic(pool.options.Logger, "TxPool.processStaleTransactionsLoop", pool.processStaleTransactionsLoop)

	for {
		select {
		case <-pool.options.Context.Done():
			return
		case blockNumber := <-pool.processStaleChan:
			pool.processStaleSenders(blockNumber)
		}
	}
}

// maxBlockRetries is how often a block no client could serve is tried again before it
// is skipped.
const maxBlockRetries = 5

// retryBlock counts a failed attempt to process a block and reports whether to give up
// on it. Until then the block loop stops at it and tries again on its next round, so no
// confirmation is missed. A skipped block may have held some, so every sender with
// pending transactions gets its on-chain state checked.
func (pool *TxPool) retryBlock(blockNumber uint64) bool {
	if pool.retryBlockNumber != blockNumber {
		pool.retryBlockNumber = blockNumber
		pool.retryBlockCount = 0
	}

	pool.retryBlockCount++

	if pool.retryBlockCount <= maxBlockRetries {
		logrus.Warnf("failed to process block %v, retrying (%v/%v)", blockNumber, pool.retryBlockCount, maxBlockRetries)

		return false
	}

	logrus.Errorf("failed to process block %v, skipping it", blockNumber)
	pool.markSendersSuspect()

	return true
}

// lookupWallet returns the registered wallet of an address, or nil.
func (pool *TxPool) lookupWallet(address common.Address) *Wallet {
	pool.walletsMutex.RLock()
	defer pool.walletsMutex.RUnlock()

	if registration := pool.wallets[address]; registration != nil {
		return registration.wallet
	}

	return nil
}

// processBlock processes a single block for transaction confirmations and reorg detection.
// It loads the block body, checks for chain reorganizations by comparing parent hashes,
// stores block information for reorg tracking, and processes all transactions in the block.
// Also handles cleanup of old block data based on the reorg depth setting.
func (pool *TxPool) processBlock(ctx context.Context, client *Client, blockNumber uint64, blockWithHash *BlockWithHash, receipts []*txtypes.Receipt) error {
	pool.lastBlockNumber = blockNumber

	pool.walletsMutex.RLock()
	walletCount := len(pool.wallets)
	pool.walletsMutex.RUnlock()

	if walletCount == 0 {
		return nil
	}

	if blockWithHash == nil || blockWithHash.Block == nil {
		blockWithHash = pool.getBlockBody(ctx, client, blockNumber)
		if blockWithHash == nil {
			return fmt.Errorf("could not load block body")
		}
	}

	// Check for reorg by comparing parent hash
	pool.blocksMutex.RLock()
	lastBlock, hasLastBlock := pool.blocks[blockNumber-1]
	pool.blocksMutex.RUnlock()

	if hasLastBlock && lastBlock.Hash != blockWithHash.Block.ParentHash {
		logrus.Warnf("Detected chain reorganization at block %d. Parent hash mismatch: expected %s, got %s",
			blockNumber, lastBlock.Hash.Hex(), blockWithHash.Block.ParentHash.Hex())

		// Handle reorg
		pool.handleReorg(ctx, client, blockNumber, blockWithHash)
	}

	// Store block info
	pool.blocksMutex.Lock()
	pool.blocks[blockNumber] = &BlockInfo{
		Number:     blockNumber,
		Hash:       blockWithHash.Hash,
		ParentHash: blockWithHash.Block.ParentHash,
		Timestamp:  blockWithHash.Block.Timestamp,
		GasLimit:   blockWithHash.Block.GasLimit,
		SlotNumber: blockWithHash.Block.SlotNumber,
	}
	if hasLastBlock && lastBlock.Hash == blockWithHash.Block.ParentHash {
		pool.updateSlotDuration(lastBlock, pool.blocks[blockNumber])
	}

	// Update current gas limit and recent block gas tracking
	pool.blockStatsMutex.Lock()
	pool.currentGasLimit = blockWithHash.Block.GasLimit
	pool.currentBaseFee = new(big.Int).Set(blockWithHash.Block.BaseFee)
	pool.recentGasUsed[pool.recentGasIdx] = blockWithHash.Block.GasUsed
	pool.recentGasLimit[pool.recentGasIdx] = blockWithHash.Block.GasLimit
	pool.recentBaseFees[pool.recentGasIdx] = new(big.Int).Set(blockWithHash.Block.BaseFee)
	pool.recentGasIdx = (pool.recentGasIdx + 1) % feeWindowSize
	if pool.recentGasLen < feeWindowSize {
		pool.recentGasLen++
	}
	pool.blockStatsMutex.Unlock()

	// Clean up old blocks
	if blockNumber > uint64(pool.reorgDepth) {
		delete(pool.blocks, blockNumber-uint64(pool.reorgDepth))
	}
	pool.blocksMutex.Unlock()

	return pool.processBlockTxs(ctx, client, blockNumber, blockWithHash, receipts)
}

// processBlockTxs processes all transactions in a block for confirmation tracking.
// It loads block receipts, decodes transaction senders, updates wallet states for
// confirmed transactions, and tracks transaction information for reorg recovery.
// Also handles cleanup of old confirmed transaction data.
func (pool *TxPool) processBlockTxs(ctx context.Context, client *Client, blockNumber uint64, blockWithHash *BlockWithHash, receipts []*txtypes.Receipt) error {
	t1 := time.Now()
	txCount := len(blockWithHash.Block.Transactions)
	if receipts == nil {
		var err error
		receipts, err = pool.getBlockReceipts(ctx, client, blockWithHash.Hash, txCount)
		if err != nil {
			return fmt.Errorf("could not load block receipts: %w", err)
		}
	}

	loadingTime := time.Since(t1)
	t1 = time.Now()

	confirmCount := 0
	affectedWalletMap := map[common.Address]bool{}
	pool.txsMutex.Lock()
	pool.confirmedTxs[blockNumber] = []*TxInfo{}
	pool.txsMutex.Unlock()

	var completions []txCompletion

	for idx, tx := range blockWithHash.Block.Transactions {
		receipt := receipts[idx]
		if receipt == nil {
			logrus.Warnf("missing receipt for tx %v in block %v", idx, blockNumber)
			continue
		}

		txFrom, err := tx.From(pool.options.ChainId)
		if err != nil {
			logrus.Warnf("error decoding tx sender (block %v, tx %v): %v", blockNumber, idx, err)
			continue
		}

		txHash := tx.Hash()
		pool.fireHashWatcher(txHash, tx, receipt)

		txFees := utils.GetTransactionFees(tx, receipt)
		fromWallet := pool.lookupWallet(txFrom)
		toAddr := tx.To()
		toWallet := (*Wallet)(nil)
		if toAddr != nil {
			toWallet = pool.lookupWallet(*toAddr)
		}

		var tracked *trackedTx

		tracked, completions = pool.processTransactionInclusion(blockNumber, txFrom, tx, receipt, completions)

		if fromWallet != nil || toWallet != nil || tracked != nil {
			txInfo := &TxInfo{
				TxHash:     txHash,
				From:       txFrom,
				To:         tx.To(),
				Tx:         tx,
				TxFees:     txFees,
				FromWallet: fromWallet,
				ToWallet:   toWallet,
				receipt:    receipt,
				tracked:    tracked,
			}

			if tracked != nil {
				txInfo.Options = tracked.Options
				txInfo.walletPool = tracked.walletPool
			} else if fromWallet != nil {
				txInfo.walletPool = fromWallet.getOwnerPool()
			}

			pool.txsMutex.Lock()
			pool.confirmedTxs[blockNumber] = append(pool.confirmedTxs[blockNumber], txInfo)
			pool.txsMutex.Unlock()
		}

		if fromWallet != nil {
			confirmCount++
			affectedWalletMap[txFrom] = true
			fromWallet.confirmTx(tx, txFees)
		}

		if toWallet != nil {
			toWallet.AddBalance(tx.Value())
		}

		completions = pool.processAuthorizations(tx, completions)
	}

	// Clean up old confirmed transactions
	if blockNumber > uint64(pool.reorgDepth) {
		oldBlock := blockNumber - uint64(pool.reorgDepth)
		pool.txsMutex.Lock()
		for blockNumber := range pool.confirmedTxs {
			if blockNumber < oldBlock {
				delete(pool.confirmedTxs, blockNumber)
			}
		}
		pool.txsMutex.Unlock()
	}

	logrus.Infof("processed block %v:  %v total tx, %v tx confirmed from %v wallets (%v, %v)", blockNumber, txCount, confirmCount, len(affectedWalletMap), loadingTime, time.Since(t1))

	// Notify block subscribers with wallet-specific stats
	pool.notifyBlockSubscribers(blockNumber, blockWithHash.Block, receipts)

	// The receipts were only kept for the statistics.
	pool.txsMutex.RLock()
	for _, txInfo := range pool.confirmedTxs[blockNumber] {
		txInfo.receipt = nil
	}
	pool.txsMutex.RUnlock()

	// Transaction callbacks run last, so they see the block fully processed.
	for _, completion := range completions {
		go pool.notifyTx(completion.tx, completion.receipt, nil)
	}

	return nil
}

// getHighestBlockNumber queries all good clients to find the highest block number.
// It runs concurrent queries to all available clients and returns the highest
// block number found along with the clients that reported that height.
// Builder clients are excluded as they don't support eth_blockNumber.
func (pool *TxPool) getHighestBlockNumber() (uint64, []*Client) {
	clientCount := len(pool.options.ClientPool.GetAllGoodClients())
	wg := &sync.WaitGroup{}

	highestBlockNumber := uint64(0)
	highestBlockNumberMutex := sync.Mutex{}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	highestBlockNumberClients := []*Client{}

	for i := 0; i < clientCount; i++ {
		client := pool.options.ClientPool.GetClient(WithClientSelectionMode(SelectClientByIndex, i))
		if client == nil {
			continue
		}

		// Skip builder clients as they don't support eth_blockNumber
		if client.IsBuilder() {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			blockNumber, err := client.GetBlockNumber(ctx)
			if err != nil {
				return
			}

			highestBlockNumberMutex.Lock()
			if blockNumber > highestBlockNumber {
				highestBlockNumber = blockNumber
				highestBlockNumberClients = []*Client{client}
			} else if blockNumber == highestBlockNumber {
				highestBlockNumberClients = append(highestBlockNumberClients, client)
			}
			highestBlockNumberMutex.Unlock()
		}()
	}

	wg.Wait()
	return highestBlockNumber, highestBlockNumberClients
}

// getBlockBody retrieves a block body from the specified client.
// Builder clients are not supported as they don't provide eth_getBlockByNumber.
func (pool *TxPool) getBlockBody(ctx context.Context, client *Client, blockNumber uint64) *BlockWithHash {
	if client.IsBuilder() {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	block, err := client.GetBlockByNumber(ctx, blockNumber)
	if err != nil || block == nil {
		return nil
	}

	if block.Hash == (common.Hash{}) {
		return nil
	}

	return &BlockWithHash{
		Hash:  block.Hash,
		Block: block,
	}
}

// getBlockReceipts retrieves all transaction receipts for a block and validates that
// their count matches the block's transaction count.
// Builder clients are not supported as they don't provide eth_getBlockReceipts.
func (pool *TxPool) getBlockReceipts(ctx context.Context, client *Client, blockHash common.Hash, txCount int) ([]*txtypes.Receipt, error) {
	if client.IsBuilder() {
		return nil, fmt.Errorf("builder clients do not support eth_getBlockReceipts")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	blockReceipts, err := client.GetBlockReceipts(ctx, blockHash)
	if err != nil {
		return nil, err
	}

	if len(blockReceipts) != txCount {
		return nil, fmt.Errorf("block %v has %v receipts, expected %v", blockHash.Hex(), len(blockReceipts), txCount)
	}

	return blockReceipts, nil
}

// SubscribeToBlockUpdates subscribes to block update notifications for a specific wallet pool.
// Returns a unique subscription ID that can be used to unsubscribe later.
func (pool *TxPool) SubscribeToBlockUpdates(walletPool *WalletPool, callback BlockStatsCallback) uint64 {
	pool.subscriptionsMutex.Lock()
	defer pool.subscriptionsMutex.Unlock()

	// Generate unique subscription ID
	id := pool.nextSubscriptionID.Add(1)

	pool.subscriptions[id] = &BlockSubscription{
		ID:         id,
		WalletPool: walletPool,
		Callback:   callback,
	}

	return id
}

// UnsubscribeFromBlockUpdates removes a block update subscription.
func (pool *TxPool) UnsubscribeFromBlockUpdates(id uint64) {
	pool.subscriptionsMutex.Lock()
	defer pool.subscriptionsMutex.Unlock()

	delete(pool.subscriptions, id)
}

// GetActiveWalletPools returns all active wallet pools
func (pool *TxPool) GetActiveWalletPools() []*WalletPool {
	pool.walletsMutex.RLock()
	defer pool.walletsMutex.RUnlock()

	walletPools := make([]*WalletPool, 0, len(pool.walletPools))
	for walletPool := range pool.walletPools {
		walletPools = append(walletPools, walletPool)
	}
	return walletPools
}

// SubscribeToBulkBlockUpdates subscribes to block update notifications for ALL wallet pools.
// Returns a unique subscription ID that can be used to unsubscribe later.
func (pool *TxPool) SubscribeToBulkBlockUpdates(callback BulkBlockStatsCallback) uint64 {
	pool.subscriptionsMutex.Lock()
	defer pool.subscriptionsMutex.Unlock()

	// Generate unique subscription ID
	id := pool.nextSubscriptionID.Add(1)

	pool.bulkSubscriptions[id] = &BulkBlockSubscription{
		ID:       id,
		Callback: callback,
	}

	return id
}

// UnsubscribeFromBulkBlockUpdates removes a bulk block update subscription.
func (pool *TxPool) UnsubscribeFromBulkBlockUpdates(id uint64) {
	pool.subscriptionsMutex.Lock()
	defer pool.subscriptionsMutex.Unlock()

	delete(pool.bulkSubscriptions, id)
}

// notifyBlockSubscribers notifies all subscribers about a processed block with wallet-specific stats.
func (pool *TxPool) notifyBlockSubscribers(blockNumber uint64, block *txtypes.Block, receipts []*txtypes.Receipt) {
	pool.subscriptionsMutex.RLock()

	// Copy subscriptions for concurrent access
	subscriptions := make(map[uint64]*BlockSubscription, len(pool.subscriptions))
	for id, sub := range pool.subscriptions {
		subscriptions[id] = sub
	}

	bulkSubscriptions := make(map[uint64]*BulkBlockSubscription, len(pool.bulkSubscriptions))
	for id, sub := range pool.bulkSubscriptions {
		bulkSubscriptions[id] = sub
	}

	pool.subscriptionsMutex.RUnlock()

	pool.txsMutex.RLock()
	blockConfirmedTxs := pool.confirmedTxs[blockNumber]
	pool.txsMutex.RUnlock()

	// If we have bulk subscribers, calculate stats for all wallet pools once
	var globalBlockStats *GlobalBlockStats
	if len(bulkSubscriptions) > 0 {
		globalBlockStats = pool.calculateAllWalletPoolStats(blockConfirmedTxs, block, receipts)

		// Notify bulk subscribers
		for _, bulkSubscription := range bulkSubscriptions {
			bulkSubscription.Callback(blockNumber, globalBlockStats)
		}
	}

	// Notify individual subscribers (reuse stats if already calculated)
	for _, subscription := range subscriptions {
		var stats *WalletPoolBlockStats
		if globalBlockStats != nil {
			// Reuse already calculated stats
			stats = globalBlockStats.WalletPoolStats[subscription.WalletPool]
		}

		if stats == nil {
			// Calculate stats individually
			stats = pool.calculateWalletPoolStats(subscription.WalletPool, blockConfirmedTxs, block)
		}
		subscription.Callback(blockNumber, stats)
	}
}

// calculateWalletPoolStats calculates transaction statistics for a specific wallet pool.
func (pool *TxPool) calculateWalletPoolStats(walletPool *WalletPool, confirmedTxs []*TxInfo, block *txtypes.Block) *WalletPoolBlockStats {
	stats := newWalletPoolBlockStats(block, len(confirmedTxs))
	affectedWallets := make(map[common.Address]struct{}, len(confirmedTxs))

	for _, txInfo := range confirmedTxs {
		if txInfo.walletPool == walletPool {
			stats.addTx(txInfo)
			affectedWallets[txInfo.From] = struct{}{}
		}
	}

	stats.AffectedWallets = len(affectedWallets)

	return stats
}

// calculateAllWalletPoolStats calculates statistics for ALL active wallet pools in a single pass.
func (pool *TxPool) calculateAllWalletPoolStats(confirmedTxs []*TxInfo, block *txtypes.Block, receipts []*txtypes.Receipt) *GlobalBlockStats {
	activeWalletPools := pool.GetActiveWalletPools()
	allStats := make(map[*WalletPool]*WalletPoolBlockStats, len(activeWalletPools))
	affectedWallets := make(map[*WalletPool]map[common.Address]struct{}, len(activeWalletPools))

	for _, walletPool := range activeWalletPools {
		allStats[walletPool] = newWalletPoolBlockStats(block, 0)
	}

	for _, txInfo := range confirmedTxs {
		stats := allStats[txInfo.walletPool]
		if stats == nil {
			continue
		}

		stats.addTx(txInfo)

		if affectedWallets[txInfo.walletPool] == nil {
			affectedWallets[txInfo.walletPool] = make(map[common.Address]struct{}, 16)
		}

		affectedWallets[txInfo.walletPool][txInfo.From] = struct{}{}
	}

	for walletPool, stats := range allStats {
		stats.AffectedWallets = len(affectedWallets[walletPool])
	}

	return &GlobalBlockStats{
		Block:           block,
		Receipts:        receipts,
		WalletPoolStats: allStats,
	}
}

func newWalletPoolBlockStats(block *txtypes.Block, capacity int) *WalletPoolBlockStats {
	return &WalletPoolBlockStats{
		TotalTxFees:  big.NewInt(0),
		Block:        block,
		Receipts:     make([]*txtypes.Receipt, 0, capacity),
		ConfirmedTxs: make([]*TxInfo, 0, capacity),
	}
}

// addTx counts a transaction the pool sent.
func (stats *WalletPoolBlockStats) addTx(txInfo *TxInfo) {
	stats.ConfirmedTxCount++

	if txInfo.TxFees != nil {
		totalFee := new(big.Int).Add(&txInfo.TxFees.FeeAmount, &txInfo.TxFees.BlobFeeAmount)
		stats.TotalTxFees.Add(stats.TotalTxFees, totalFee)
	}

	stats.Receipts = append(stats.Receipts, txInfo.receipt)
	stats.ConfirmedTxs = append(stats.ConfirmedTxs, txInfo)
}

// recoveryConcurrency bounds the stale sender checks running at once, and separately
// the rebroadcasts (including resubmissions after a reorg). It keeps a large backlog
// from flooding the clients.
const recoveryConcurrency = 8

// maxRebroadcastNonces is the number of lowest pending nonces of a sender considered
// for rebroadcasting.
const maxRebroadcastNonces = 2

// txCompletion is a transaction that completed while a block was processed.
type txCompletion struct {
	tx      *trackedTx
	receipt *txtypes.Receipt
}

// submitTransaction handles the core transaction submission logic.
// It tracks the transaction until it completes and submits it to clients. The submitNow
// parameter controls whether to immediately submit or just set up confirmation tracking.
func (pool *TxPool) submitTransaction(ctx context.Context, wallet *Wallet, tx *txtypes.Transaction, options *SendTransactionOptions, submitNow bool) error {
	_, err := pool.trackTransaction(ctx, wallet, tx, options, submitNow, nil)

	return err
}

// trackTransaction is submitTransaction with the wallet pool the transaction is
// attributed to. A nil walletPool derives it from the context and the wallet.
func (pool *TxPool) trackTransaction(ctx context.Context, wallet *Wallet, tx *txtypes.Transaction, options *SendTransactionOptions, submitNow bool, walletPool *WalletPool) (*trackedTx, error) {
	if options == nil {
		options = &SendTransactionOptions{}
	}

	if ctx.Err() != nil {
		// OnComplete is documented as always being called once processing ends.
		// Callers such as WalletPool.ReclaimFunds release a WaitGroup from it and
		// would hang forever if a pre-cancelled context skipped the callback.
		options.invokeComplete(tx, nil, ctx.Err())

		return nil, ctx.Err()
	}

	if walletPool == nil {
		walletPool = attributeTx(ctx, wallet)
	}

	tracked := newTrackedTx(wallet, tx, options)
	tracked.walletPool = walletPool

	if walletPool != nil {
		walletPool.pendingTxCount.Add(1)
	}

	// A transaction whose account nonce is already used cannot be confirmed through
	// its lane anymore; its receipt is looked up instead.
	laned := true
	if tracked.keys == nil {
		wallet.txNonceMutex.Lock()
		laned = wallet.confirmedTxCount <= tracked.seq
		wallet.txNonceMutex.Unlock()
	}

	pool.tracker.add(tracked, laned, pool.lastBlockNumber)

	stopCtx := context.AfterFunc(ctx, func() {
		if pool.finishTx(tracked) {
			pool.notifyTx(tracked, nil, nil)
		}
	})
	tracked.stopCtx.Store(&stopCtx)

	if tracked.done.Load() {
		tracked.stopCtx.Store(nil)
		stopCtx()
	}

	if submitNow {
		submitted, err := pool.sendToClients(ctx, tx, options)
		if err != nil {
			pool.completeTx(tracked, nil, err)

			return nil, err
		}

		if submitted {
			wallet.IncrementSubmittedTxCount()

			if walletPool != nil {
				walletPool.submittedTxCount.Add(1)
			}
		}
	}

	if !laned {
		go func() {
			pool.completeTx(tracked, pool.loadTransactionReceipt(ctx, tx), nil)
		}()
	}

	return tracked, nil
}

// sendToClients submits a transaction to the clients the options select, until
// SubmitCount of them accepted it. It reports whether any did, and fails if all refused.
func (pool *TxPool) sendToClients(ctx context.Context, tx *txtypes.Transaction, options *SendTransactionOptions) (bool, error) {
	var submitErr error

	submitTx := func(client *Client) error {
		if options.hasEncodeFn() {
			txBytes, err := options.invokeEncode(tx)
			if err != nil {
				return fmt.Errorf("failed to encode transaction: %w", err)
			}

			if len(txBytes) > 0 {
				return client.SendRawTransaction(ctx, txBytes)
			}
		}

		return client.SendTransaction(ctx, tx)
	}

	submitCount := options.SubmitCount
	if submitCount == 0 {
		submitCount = 3
	}

	// Derive start offset for redundant client selection.
	// When ClientsStartOffset is 0 (the default, never explicitly set by any scenario),
	// use the tx hash to distribute redundant submissions uniformly across all clients
	// instead of always targeting clients at index 1 and 2.
	startOffset := options.ClientsStartOffset
	if startOffset == 0 && tx != nil {
		h := tx.Hash()
		startOffset = int(h[0]) | (int(h[1]) << 8)
	}

	success := false
	clientCount := len(pool.options.ClientPool.GetAllGoodClients())
	for i := 0; i < clientCount; i++ {
		client := options.Client
		if client == nil || i > 0 {
			if len(options.ClientList) > 0 {
				client = options.ClientList[(i+startOffset)%len(options.ClientList)]
			} else {
				var clientOpts []ClientSelectionOption
				if options.ClientGroup != "" {
					clientOpts = append(clientOpts, WithClientGroup(options.ClientGroup))
				}
				clientOpts = append(clientOpts, WithClientSelectionMode(SelectClientByIndex, i+startOffset))
				client = pool.options.ClientPool.GetClient(clientOpts...)
			}
		}
		if client == nil {
			continue
		}

		err := submitTx(client)

		if options.LogFn != nil {
			options.LogFn(client, i, 0, err)
		}
		if err == nil || (strings.Contains(err.Error(), "already known") || strings.Contains(err.Error(), "Known transaction")) {
			success = true
			submitCount--
			if submitCount == 0 {
				break
			}
		} else if submitErr == nil {
			submitErr = err
		}
	}

	if success {
		return true, nil
	}

	return false, submitErr
}

// finishTx ends the tracking of a transaction. It reports false if the transaction
// was already finished.
func (pool *TxPool) finishTx(tracked *trackedTx) bool {
	if !tracked.done.CompareAndSwap(false, true) {
		return false
	}

	pool.tracker.remove(tracked)

	// The finished transaction stays referenced for a reorg, its context hook is not
	// needed anymore.
	if stopCtx := tracked.stopCtx.Swap(nil); stopCtx != nil {
		(*stopCtx)()
	}

	if tracked.walletPool != nil {
		tracked.walletPool.pendingTxCount.Add(-1)
	}

	return true
}

// notifyTx runs the callbacks of a finished transaction. A nil receipt with a nil
// error means the transaction was cancelled or replaced.
func (pool *TxPool) notifyTx(tracked *trackedTx, receipt *txtypes.Receipt, err error) {
	if receipt != nil {
		tracked.Options.invokeConfirm(tracked.Tx, receipt)
	}

	tracked.Options.invokeComplete(tracked.Tx, receipt, err)

	if walletPool := tracked.walletPool; walletPool != nil && walletPool.ctx.Err() == nil {
		if tracker := walletPool.GetTransactionTracker(); tracker != nil {
			tracker(err)
		}
	}
}

// completeTx finishes a transaction and runs its callbacks on their own goroutine.
func (pool *TxPool) completeTx(tracked *trackedTx, receipt *txtypes.Receipt, err error) {
	if pool.finishTx(tracked) {
		go pool.notifyTx(tracked, receipt, err)
	}
}

// retireTxs finishes the transactions an included transaction retired from their lanes:
// the included one is confirmed, the others were replaced by it.
func (pool *TxPool) retireTxs(retired []*trackedTx, txHash common.Hash, receipt *txtypes.Receipt, completions []txCompletion) []txCompletion {
	for _, tracked := range retired {
		if !pool.finishTx(tracked) {
			continue
		}

		completion := txCompletion{tx: tracked}
		if tracked.hash == txHash {
			completion.receipt = receipt
		}

		completions = append(completions, completion)
	}

	return completions
}

// processTransactionInclusion handles a transaction included in a block: it completes
// the tracked transactions with its hash, and every tracked transaction of the same
// sender that its nonce makes impossible. Returns the tracked transaction it confirms,
// if any, and the completions to notify once the block is processed.
func (pool *TxPool) processTransactionInclusion(blockNumber uint64, txFrom common.Address, tx *txtypes.Transaction, receipt *txtypes.Receipt, completions []txCompletion) (*trackedTx, []txCompletion) {
	txHash := tx.Hash()
	confirmed := pool.tracker.byHash(txHash)

	if sender := pool.tracker.getSender(txFrom); sender != nil {
		keys, seq := txLanePosition(tx)

		sender.mu.Lock()
		retired := sender.retire(keys, seq)
		sender.lastProgress = blockNumber
		sender.staleChecks = 0
		sender.mu.Unlock()

		completions = pool.retireTxs(retired, txHash, receipt, completions)
	}

	// Transactions outside the lanes are only known by their hash.
	completions = pool.retireTxs(confirmed, txHash, receipt, completions)

	if len(confirmed) == 0 {
		return nil, completions
	}

	return confirmed[0], completions
}

// processAuthorizations applies the EIP-7702 authorizations of an included transaction.
// A valid authorization uses up the authority's account nonce without a transaction from
// it, which retires whatever the authority had pending on that nonce.
func (pool *TxPool) processAuthorizations(tx *txtypes.Transaction, completions []txCompletion) []txCompletion {
	authList := tx.AuthList()
	if len(authList) == 0 {
		return completions
	}

	chainID := tx.ChainId()

	for i := range authList {
		auth := &authList[i]
		if !auth.ChainID.IsZero() && (chainID == nil || auth.ChainID.CmpBig(chainID) != 0) {
			continue
		}

		authority, err := auth.Authority()
		if err != nil {
			continue
		}

		wallet := pool.lookupWallet(authority)
		if wallet == nil || !wallet.applyAuthorization(auth.Nonce) {
			continue
		}

		if sender := pool.tracker.getSender(authority); sender != nil {
			sender.mu.Lock()
			retired := sender.retire(nil, auth.Nonce)
			sender.mu.Unlock()

			completions = pool.retireTxs(retired, common.Hash{}, nil, completions)
		}
	}

	return completions
}

// runRecovery runs a stale sender check within the recovery concurrency limit.
func (pool *TxPool) runRecovery(fn func()) {
	if pool.recoverySem == nil {
		fn()

		return
	}

	select {
	case pool.recoverySem <- struct{}{}:
	case <-pool.options.Context.Done():
		return
	}

	defer func() { <-pool.recoverySem }()

	fn()
}

// workQueue runs queued functions in order on a bounded number of workers, which
// exist only while there is work.
type workQueue struct {
	mu      sync.Mutex
	items   []func()
	workers int
}

func (queue *workQueue) push(limit int, fn func()) {
	queue.mu.Lock()
	queue.items = append(queue.items, fn)

	spawn := queue.workers < limit
	if spawn {
		queue.workers++
	}
	queue.mu.Unlock()

	if spawn {
		go queue.work()
	}
}

func (queue *workQueue) work() {
	for {
		queue.mu.Lock()
		if len(queue.items) == 0 {
			queue.items = nil
			queue.workers--
			queue.mu.Unlock()

			return
		}

		fn := queue.items[0]
		queue.items[0] = nil
		queue.items = queue.items[1:]
		queue.mu.Unlock()

		fn()
	}
}

// queueRebroadcast rebroadcasts a tracked transaction through the rebroadcast queue.
// fanout sends it to every eligible client instead of the first accepting one.
func (pool *TxPool) queueRebroadcast(tracked *trackedTx, fanout bool) {
	if !tracked.rebroadcasting.CompareAndSwap(false, true) {
		return
	}

	tracked.LastRebroadcast = time.Now()
	tracked.RebroadcastCount++
	retryCount := tracked.RebroadcastCount

	pool.rebroadcasts.push(recoveryConcurrency, func() {
		defer tracked.rebroadcasting.Store(false)

		if tracked.done.Load() || pool.options.Context.Err() != nil {
			return
		}

		pool.rebroadcastTransaction(pool.options.Context, tracked.Tx, tracked.Options, retryCount, fanout)
	})
}

// staleCheckInterval is the number of blocks a sender may go without a confirmation
// before its pending transactions are checked, growing with every fruitless check.
func staleCheckInterval(staleChecks uint64) uint64 {
	return 10 << min(staleChecks, 2)
}

// processStaleSenders recovers senders whose pending transactions made no progress:
// it re-reads their on-chain nonces to catch confirmations and replacements block
// processing did not see, and rebroadcasts what is still pending.
func (pool *TxPool) processStaleSenders(blockNumber uint64) {
	var wg sync.WaitGroup

	for _, sender := range pool.tracker.allSenders() {
		sender.mu.Lock()
		due := sender.suspect || blockNumber > sender.lastProgress+staleCheckInterval(sender.staleChecks)
		if due {
			sender.suspect = false
			sender.lastProgress = blockNumber
			sender.staleChecks++
		}
		sender.mu.Unlock()

		if !due {
			continue
		}

		wg.Add(1)

		go func() {
			defer wg.Done()

			pool.runRecovery(func() { pool.recoverSender(blockNumber, sender) })
		}()
	}

	wg.Wait()
}

// markSendersSuspect makes the next stale pass check every sender, after a block could
// not be processed and its confirmations may have been missed.
func (pool *TxPool) markSendersSuspect() {
	for _, sender := range pool.tracker.allSenders() {
		sender.mu.Lock()
		sender.suspect = true
		sender.mu.Unlock()
	}
}

func (pool *TxPool) recoverSender(blockNumber uint64, sender *txSender) {
	sender.mu.Lock()
	hasAccountTxs := !sender.account.empty()

	var keyedHeads []*trackedTx

	for _, lane := range sender.keyed {
		for _, tracked := range lane.slots[0].txs {
			if !slices.Contains(keyedHeads, tracked) {
				keyedHeads = append(keyedHeads, tracked)
			}
		}
	}
	sender.mu.Unlock()

	if hasAccountTxs {
		pool.recoverAccountLane(blockNumber, sender)
	}

	for _, tracked := range keyedHeads {
		pool.recoverKeyedTx(sender, tracked)
	}
}

// recoverAccountLane reconciles a sender's account nonce lane with the on-chain nonce.
func (pool *TxPool) recoverAccountLane(blockNumber uint64, sender *txSender) {
	wallet := sender.wallet

	var onChainNonce uint64
	var err error
	for retry := 0; retry < 3; retry++ {
		client := pool.options.ClientPool.GetClient(WithClientSelectionMode(SelectClientRandom))
		if client == nil {
			continue
		}

		onChainNonce, err = client.GetNonceAt(pool.options.Context, sender.address, big.NewInt(int64(blockNumber)))
		if err == nil {
			break
		}
	}

	if err != nil {
		logrus.WithError(err).Warnf("failed to get on-chain nonce for %v", sender.address.String())
		return
	}

	var (
		confirmedTxs     []*trackedTx
		txsToRebroadcast []*trackedTx
		expiredTxs       []*trackedTx
		nonceGaps        []uint64
	)

	chainTime := pool.blockTime(blockNumber)

	sender.mu.Lock()

	if onChainNonce > 0 {
		confirmedTxs = sender.account.retire(onChainNonce-1, nil)
	}

	for _, slot := range sender.account.slots {
		if slot.seq > onChainNonce+(maxRebroadcastNonces-1) {
			break
		}

		if len(slot.txs) == 0 {
			continue
		}

		mostRecentTx := slot.txs[len(slot.txs)-1]

		// A transaction past its own deadline is dead rather than slow: it is
		// replaced rather than rebroadcast, and regardless of whether the
		// caller asked for rebroadcasts, because until its nonce is retired
		// nothing else this wallet sends can confirm.
		if isExpiredTx(mostRecentTx.Tx, chainTime) {
			expiredTxs = append(expiredTxs, mostRecentTx)

			continue
		}

		if mostRecentTx.Options.Rebroadcast && time.Since(mostRecentTx.LastRebroadcast) >= pool.calculateBackoffDelay(mostRecentTx.RebroadcastCount) {
			txsToRebroadcast = append(txsToRebroadcast, mostRecentTx)
		}
	}

	if !sender.account.empty() {
		for nonce := onChainNonce; nonce < sender.account.slots[0].seq; nonce++ {
			nonceGaps = append(nonceGaps, nonce)
		}
	}

	sender.mu.Unlock()

	wallet.raiseConfirmedNonce(onChainNonce)

	// These confirmed in a block that was not processed. Whether as the tracked
	// transaction or as a replacement of it, only the receipt can tell.
	for _, tracked := range confirmedTxs {
		logrus.Debugf("recovering stale confirmed transaction for %v (nonce %v)", sender.address.String(), tracked.seq)
		pool.completeTx(tracked, pool.loadTransactionReceipt(pool.options.Context, tracked.Tx), nil)
	}

	// Retire expired transactions, which hold a nonce nothing can ever confirm. The
	// attempt is counted like a rebroadcast, so a replacement a client refuses as
	// underpriced is retried at a higher price rather than at the same one forever.
	for _, expired := range expiredTxs {
		expired.RebroadcastCount++
		expired.LastRebroadcast = time.Now()

		go pool.replaceExpiredTransaction(pool.options.Context, wallet, &expired.PendingTx, expired.RebroadcastCount)
	}

	if len(nonceGaps) > 0 {
		logrus.WithFields(logrus.Fields{
			"wallet": sender.address.Hex(),
			"gaps":   nonceGaps,
		}).Warnf("detected nonce gaps, filling with dummy transactions")
		go pool.fillNonceGaps(pool.options.Context, wallet, nonceGaps, nil)
	}

	for _, tracked := range txsToRebroadcast {
		logrus.WithFields(logrus.Fields{
			"wallet": sender.address.Hex(),
			"nonce":  tracked.seq,
			"txhash": tracked.hash.Hex(),
		}).Debugf("rebroadcasting stale transaction")

		pool.queueRebroadcast(tracked, true)
	}
}

// recoverKeyedTx recovers the lowest pending transaction of an EIP-8250 keyed lane. Its
// receipt tells whether it was included, the key's on-chain sequence whether another
// transaction used it up.
func (pool *TxPool) recoverKeyedTx(sender *txSender, tracked *trackedTx) {
	if tracked.done.Load() {
		return
	}

	receipt := pool.loadTransactionReceipt(pool.options.Context, tracked.Tx)

	if receipt != nil || pool.isKeyedLanePassed(sender.address, tracked) {
		logrus.Debugf("recovering stale transaction for %v (hash %v)", sender.address.String(), tracked.hash.Hex())

		sender.mu.Lock()
		retired := sender.retire(tracked.keys, tracked.seq)
		sender.mu.Unlock()

		for _, completion := range pool.retireTxs(retired, tracked.hash, receipt, nil) {
			go pool.notifyTx(completion.tx, completion.receipt, nil)
		}

		return
	}

	if !tracked.Options.Rebroadcast || time.Since(tracked.LastRebroadcast) < pool.calculateBackoffDelay(tracked.RebroadcastCount) {
		return
	}

	logrus.WithFields(logrus.Fields{
		"wallet": sender.address.Hex(),
		"txhash": tracked.hash.Hex(),
	}).Debugf("rebroadcasting stale transaction")

	pool.queueRebroadcast(tracked, true)
}

// isKeyedLanePassed reports whether the first nonce key of a transaction has moved past
// the sequence number the transaction needs it at.
func (pool *TxPool) isKeyedLanePassed(address common.Address, tracked *trackedTx) bool {
	if len(tracked.keys) == 0 {
		return false
	}

	client := pool.options.ClientPool.GetClient(WithClientSelectionMode(SelectClientRandom))
	if client == nil {
		return false
	}

	value, err := client.GetStorageAt(pool.options.Context, txtypes.NonceManager, txtypes.NonceManagerSlot(address, &tracked.keys[0]))
	if err != nil {
		return false
	}

	current := new(uint256.Int).SetBytes(value[:])
	if tracked.seq == txtypes.MaxNonceSeq {
		// single-use: the key has to be unused
		return !current.IsZero()
	}

	return !current.IsUint64() || current.Uint64() > tracked.seq
}

// loadTransactionReceipt attempts to load a transaction receipt from multiple clients.
// It retries up to 5 times with different clients and includes exponential backoff.
// Returns nil if the receipt cannot be loaded after all retries.
func (pool *TxPool) loadTransactionReceipt(ctx context.Context, tx *txtypes.Transaction) *txtypes.Receipt {
	retryCount := uint64(0)

	for {
		client := pool.options.ClientPool.GetClient(WithClientSelectionMode(SelectClientRandom))
		if client == nil {
			return nil
		}

		reqCtx, reqCtxCancel := context.WithTimeout(ctx, 5*time.Second)

		//nolint:gocritic // ignore
		defer reqCtxCancel()

		receipt, err := client.GetTransactionReceipt(reqCtx, tx.Hash())
		if err == nil {
			return receipt
		}

		// A transaction the client does not know has no receipt; that is an answer,
		// not a failure to retry.
		if ctx.Err() != nil || errors.Is(err, ethereum.NotFound) {
			return nil
		}

		if retryCount > 2 {
			logrus.WithFields(logrus.Fields{
				"client": client.GetName(),
				"txhash": tx.Hash(),
			}).Warnf("could not load tx receipt: %v", err)
		}

		if retryCount < 5 {
			time.Sleep(1 * time.Second)

			retryCount++
		} else {
			return nil
		}
	}
}

// handleReorg handles a detected chain reorganization by re-submitting affected transactions.
// It finds the common ancestor, identifies reorged-out transactions, resets wallet nonces,
// and re-submits the affected transactions as pending. Also processes the new canonical blocks.
func (pool *TxPool) handleReorg(ctx context.Context, client *Client, blockNumber uint64, newBlockWithHash *BlockWithHash) error {
	type newBlockInfo struct {
		blockWithHash *BlockWithHash
	}

	newBlockParents := []newBlockInfo{}

	// let's assume a reorg of 2 blocks:
	// old chain: 1 -> 2 -> 3a -> 4a
	// new chain: 1 -> 2 -> 3b -> 4b

	// we'll find out about this reorg when receiving block 4b
	// we need to:
	// 1. find the common ancestor of the two blocks (block 2), by loading the new parents (3b)
	// 2. find all the transactions that were reorged out in block 3a - 4a
	// 3. add reorged out transactions as pending txs, reset affected wallets nonce state
	// 4. remove blocks 3a - 4a from the pool
	// 5. re-process block 3b (4b will be processed after the reorg processing completes)

	// find the common ancestor
	block := newBlockWithHash
	for {
		if block.Block.NumberU64() == 0 {
			break
		}

		blockNumber := block.Block.NumberU64() - 1
		if pool.blocks[blockNumber] == nil {
			break
		}

		if pool.blocks[blockNumber].Hash == block.Block.ParentHash {
			break
		}

		parentBlockBody := pool.getBlockBody(ctx, client, blockNumber)
		if parentBlockBody == nil {
			return fmt.Errorf("could not load block body for new parent block %v", blockNumber)
		}

		newBlockParents = append(newBlockParents, newBlockInfo{
			blockWithHash: parentBlockBody,
		})
		block = parentBlockBody
	}

	reorgBaseBlock := block

	// find all the transactions that were reorged out
	pool.txsMutex.Lock()
	reorgedOutTxs := []*TxInfo{}
	for blockNum := reorgBaseBlock.Block.NumberU64(); blockNum <= blockNumber; blockNum++ {
		blockTxs := pool.confirmedTxs[blockNum]
		if blockTxs == nil {
			continue
		}

		reorgedOutTxs = append(reorgedOutTxs, blockTxs...)
	}

	// remove reorged out blocks & txs from cache
	pool.blocksMutex.Lock()
	for blockNum := reorgBaseBlock.Block.NumberU64() + 1; blockNum <= blockNumber; blockNum++ {
		delete(pool.blocks, blockNum)
		delete(pool.confirmedTxs, blockNum)
	}
	pool.blocksMutex.Unlock()
	pool.txsMutex.Unlock()

	// add reorged out txs as pending txs, reset affected wallets nonce state
	resetWallets := map[common.Address]bool{}
	resubmitted := map[common.Address]bool{}

	for _, tx := range reorgedOutTxs {
		wallet := tx.FromWallet
		if wallet == nil && tx.tracked != nil {
			wallet = tx.tracked.wallet
		}

		if tx.FromWallet != nil {
			// A transaction outside the account nonce sequence carries no information
			// about where that sequence should be rewound to.
			if !resetWallets[tx.From] && tx.Tx.UsesAccountNonce() {
				resetWallets[tx.From] = true

				tx.FromWallet.txNonceMutex.Lock()
				tx.FromWallet.confirmedTxCount = tx.Tx.Nonce()
				tx.FromWallet.txNonceMutex.Unlock()
			}

			totalAmount := new(big.Int).Add(tx.Tx.Value(), &tx.TxFees.FeeAmount)
			totalAmount = new(big.Int).Add(totalAmount, &tx.TxFees.BlobFeeAmount)
			tx.FromWallet.AddBalance(totalAmount)
		}

		if wallet != nil {
			// add tx as pending tx
			txOptions := &SendTransactionOptions{
				Client:      client,
				Rebroadcast: true,
				OnComplete: func(tx *txtypes.Transaction, receipt *txtypes.Receipt, err error) {
					if err == nil && receipt != nil {
						logrus.Infof("reorged out tx %v confirmed", tx.Hash())
					}
				},
			}

			if tx.Options != nil {
				txOptions.LogFn = tx.Options.LogFn
				txOptions.ClientGroup = tx.Options.ClientGroup
				txOptions.ClientList = tx.Options.ClientList
				txOptions.OnEncode = tx.Options.OnEncode
				txOptions.OnGethEncode = tx.Options.OnGethEncode
			}

			// The transaction as it was sent still has what the block's copy lacks,
			// such as blob sidecars.
			reorgedTx := tx.Tx
			if tx.tracked != nil {
				reorgedTx = tx.tracked.Tx
			}

			tracked, err := pool.trackTransaction(ctx, wallet, reorgedTx, txOptions, false, tx.walletPool)
			if err != nil {
				logrus.WithError(err).Errorf("error adding pending transaction for reorged out tx %v", tx.Tx.Hash())
			} else {
				// Clients may have dropped it with the orphaned block, so it is sent
				// again right away: the first one of each sender, which the others
				// depend on, to every client.
				pool.queueRebroadcast(tracked, !resubmitted[tx.From])
				resubmitted[tx.From] = true
			}
		}

		// An authorization's nonce bump is undone with the block as well.
		for i := range tx.Tx.AuthList() {
			if authority, err := tx.Tx.AuthList()[i].Authority(); err == nil {
				if authorityWallet := pool.lookupWallet(authority); authorityWallet != nil {
					authorityWallet.MarkNeedResync()
				}
			}
		}

		if tx.ToWallet != nil {
			// reverse processTransactionReceival
			tx.ToWallet.SubBalance(tx.Tx.Value())
		}
	}

	// re-process the new parent blocks
	slices.Reverse(newBlockParents)
	for _, parentBlockInfo := range newBlockParents {
		pool.processBlockTxs(ctx, client, parentBlockInfo.blockWithHash.Block.NumberU64(), parentBlockInfo.blockWithHash, nil)
	}

	return nil
}

// updateSlotDuration derives the slot duration from the EIP-7843 slot numbers and
// timestamps of two consecutive blocks.
func (pool *TxPool) updateSlotDuration(parent *BlockInfo, block *BlockInfo) {
	if parent.SlotNumber == 0 || block.SlotNumber <= parent.SlotNumber || block.Timestamp <= parent.Timestamp {
		return
	}

	duration := time.Duration(block.Timestamp-parent.Timestamp) * time.Second / time.Duration(block.SlotNumber-parent.SlotNumber)
	if time.Duration(pool.slotDuration.Swap(int64(duration))) != duration {
		pool.options.Logger.Infof("slot duration changed to %v (slot %v)", duration, block.SlotNumber)
	}
}

// GetSlotDuration returns the slot duration derived from EIP-7843 block slot numbers,
// or the configured fallback while the chain does not report slot numbers.
func (pool *TxPool) GetSlotDuration() time.Duration {
	if duration := time.Duration(pool.slotDuration.Load()); duration > 0 {
		return duration
	}

	return pool.options.SlotDuration
}

// GetCurrentGasLimit returns the current gas limit of the chain.
func (pool *TxPool) GetCurrentGasLimit() uint64 {
	pool.blockStatsMutex.RLock()
	defer pool.blockStatsMutex.RUnlock()
	return pool.currentGasLimit
}

// IsAmsterdam reports whether spamoor uses the Amsterdam (EIP-8037) gas/fee
// model. Defaults to true; set TxPoolOptions.PreAmsterdamFeeModel to true to
// force the legacy model. The Amsterdam model is safe to keep on non-Amsterdam
// chains because its gas budgets are strictly higher than the legacy ones.
func (pool *TxPool) IsAmsterdam() bool {
	return !pool.preAmsterdamFeeModel
}

// HasObservedHead reports whether the pool has loaded chain state via
// InitializeBlockStats or processBlock. With the current startup sequence this
// always returns true once scenarios are running; the method is retained for
// plugin ABI compatibility.
func (pool *TxPool) HasObservedHead() bool {
	pool.blockStatsMutex.RLock()
	defer pool.blockStatsMutex.RUnlock()
	return pool.currentGasLimit > 0
}

// GetCostPerStateByte returns the EIP-8037 cost_per_state_byte spamoor charges
// for state-creation gas budgeting. Returns the flat CostPerStateByte constant
// on Amsterdam (default), or 0 when the pre-Amsterdam fee model is forced.
func (pool *TxPool) GetCostPerStateByte() uint64 {
	if pool.preAmsterdamFeeModel {
		return 0
	}
	return CostPerStateByte
}

// MaxTxGas returns the effective per-tx gas cap for the connected chain.
// On Amsterdam, EIP-8037 caps only the state-creation portion at EIP-7825's
// TX_MAX_GAS_LIMIT (16,777,216) while regular (compute) gas is bounded only by
// the block gas limit, so MaxTxGas returns the block gas limit. Pre-Amsterdam
// the cap is utils.MaxGasLimitPerTx, clamped to the block limit on chains
// whose block limit is below 16.7M.
func (pool *TxPool) MaxTxGas() uint64 {
	blockLimit := pool.GetCurrentGasLimit()
	if pool.IsAmsterdam() {
		return blockLimit
	}
	cap := uint64(utils.MaxGasLimitPerTx)
	if blockLimit > 0 && blockLimit < cap {
		cap = blockLimit
	}
	return cap
}

// MinIntrinsicGas returns the base intrinsic gas cost for a simple ETH
// transfer to an existing recipient (no calldata, no account creation).
// This is the minimum tx.Gas value any transaction must set, and the value
// our filler / funding paths use as the per-tx baseline.
//
// Currently always returns 21,000. Geth removed the prior +10% txpool
// intrinsic buffer, so the raw base cost is what the validator enforces.
//
// EIP-2780 (Reduce intrinsic transaction gas) will lower this base cost
// once activated (4,500 base + accessory charges that vary by tx shape;
// ~7,756 for a value transfer to an existing account). Centralized here
// so consumers update through a single method when EL clients ship 2780.
func (pool *TxPool) MinIntrinsicGas() uint64 {
	return 21_000
}

// GetCurrentBaseFee returns a copy of the current base fee of the chain.
func (pool *TxPool) GetCurrentBaseFee() *big.Int {
	pool.blockStatsMutex.RLock()
	defer pool.blockStatsMutex.RUnlock()

	if pool.currentBaseFee == nil {
		return nil
	}

	return new(big.Int).Set(pool.currentBaseFee)
}

// GetRecentBlockGasTotals returns the sum of gas used and gas limits
// over the recent block window. Must not hold blockStatsMutex.
func (pool *TxPool) GetRecentBlockGasTotals() (totalGasUsed, totalGasLimit uint64) {
	pool.blockStatsMutex.RLock()
	defer pool.blockStatsMutex.RUnlock()
	for i := range pool.recentGasLen {
		totalGasUsed += pool.recentGasUsed[i]
		totalGasLimit += pool.recentGasLimit[i]
	}
	return
}

// getWindowAvgBaseFee returns the average baseFee over the recent block window.
// When baseFee is rising, the average lags behind the current baseFee, creating
// natural backpressure that lets baseFee recover. Must not hold blockStatsMutex.
func (pool *TxPool) getWindowAvgBaseFee() *big.Int {
	pool.blockStatsMutex.RLock()
	defer pool.blockStatsMutex.RUnlock()
	count := 0
	sum := new(big.Int)
	for i := range pool.recentGasLen {
		if pool.recentBaseFees[i] != nil && pool.recentBaseFees[i].Sign() > 0 {
			sum.Add(sum, pool.recentBaseFees[i])
			count++
		}
	}
	if count == 0 {
		return nil
	}
	return sum.Div(sum, big.NewInt(int64(count)))
}

// GetCurrentGasLimitWithInit returns the current gas limit. With the current
// startup sequence (InitializeBlockStats runs before scenarios start) the
// state is always populated, but the lazy-init fallback is retained for
// plugin ABI compatibility.
func (pool *TxPool) GetCurrentGasLimitWithInit() (uint64, error) {
	gasLimit := pool.GetCurrentGasLimit()
	if gasLimit == 0 {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := pool.tryInitBlockStats(ctx); err != nil {
			return 0, err
		}
		gasLimit = pool.GetCurrentGasLimit()
	}
	return gasLimit, nil
}

// GetCurrentBaseFeeWithInit returns the current base fee. With the current
// startup sequence (InitializeBlockStats runs before scenarios start) the
// state is always populated, but the lazy-init fallback is retained for
// plugin ABI compatibility.
func (pool *TxPool) GetCurrentBaseFeeWithInit() (*big.Int, error) {
	baseFee := pool.GetCurrentBaseFee()
	if baseFee == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := pool.tryInitBlockStats(ctx); err != nil {
			return nil, err
		}
		baseFee = pool.GetCurrentBaseFee()
	}
	return baseFee, nil
}

// InitializeBlockStats fetches the current block stats from the network and
// must be called once during startup, after the client pool is prepared and
// before any scenario starts using the txpool. It blocks (with backoff) until
// either a head block is loaded or ctx is cancelled, so that downstream code
// can rely on GetCurrentBaseFee / GetCurrentGasLimit / IsAmsterdam returning
// real chain state. The call is idempotent: a second invocation with already
// populated state returns immediately.
func (pool *TxPool) InitializeBlockStats(ctx context.Context) error {
	backoff := 100 * time.Millisecond
	const maxBackoff = 2 * time.Second
	for {
		err := pool.tryInitBlockStats(ctx)
		if err == nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("failed to initialize block stats: %w (last error: %v)", ctx.Err(), err)
		case <-time.After(backoff):
		}

		if backoff < maxBackoff {
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}
	}
}

// tryInitBlockStats performs a single attempt to fetch the latest block and
// populate the txpool's chain-state fields. Returns immediately with nil if
// the state is already populated.
func (pool *TxPool) tryInitBlockStats(ctx context.Context) error {
	pool.blockStatsMutex.Lock()
	defer pool.blockStatsMutex.Unlock()

	if pool.currentGasLimit > 0 && pool.currentBaseFee != nil {
		return nil
	}

	client := pool.options.ClientPool.GetClient(WithoutBuilder())
	if client == nil {
		return fmt.Errorf("no non-builder client available to fetch gas limit")
	}

	fetchCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Header-only: we only need the gas limit and base fee here.
	latestHeader, err := client.GetHeaderByNumber(fetchCtx, nil)
	if err != nil {
		return fmt.Errorf("failed to fetch latest header for gas limit: %w", err)
	}

	pool.currentGasLimit = latestHeader.GasLimit
	pool.currentBaseFee = latestHeader.BaseFee
	logrus.Infof("initialized block stats from latest block: gasLimit=%v baseFee=%v amsterdamFeeModel=%v",
		pool.currentGasLimit, pool.currentBaseFee, !pool.preAmsterdamFeeModel)

	return nil
}

// GetSuggestedFees returns the suggested fees for a transaction.
// If baseFeeWei and tipFeeWei are provided (non-nil), they are used directly.
// If not provided (nil), the fees are fetched from the client.
//
// With FeeStrategy "adaptive", the dynamic fee mechanism computes feeCap from
// the current baseFee with adaptive headroom and normal distribution spread,
// using baseFeeWei as a maximum cap instead.
func (pool *TxPool) GetSuggestedFees(client *Client, baseFeeWei *big.Int, tipFeeWei *big.Int, feeStrategy ...string) (feeCap *big.Int, tipCap *big.Int, err error) {
	strategy := ""
	if len(feeStrategy) > 0 {
		strategy = feeStrategy[0]
	}
	if strategy == "adaptive" {
		return pool.getSuggestedFeesAdaptive(client, baseFeeWei, tipFeeWei)
	}

	// Legacy strategy: use provided fees or fetch from network.
	if baseFeeWei != nil && baseFeeWei.Sign() > 0 {
		feeCap = new(big.Int).Set(baseFeeWei)
	}
	if tipFeeWei != nil && tipFeeWei.Sign() > 0 {
		tipCap = new(big.Int).Set(tipFeeWei)
	}

	if feeCap == nil || tipCap == nil {
		networkFeeCap, networkTipCap, fetchErr := client.GetSuggestedFee(pool.options.Context)
		if fetchErr != nil {
			return nil, nil, fetchErr
		}
		if feeCap == nil {
			feeCap = networkFeeCap
		}
		if tipCap == nil {
			tipCap = networkTipCap
		}
	}

	return feeCap, tipCap, nil
}

// getSuggestedFeesAdaptive computes fees using dynamic headroom with normal
// distribution spread. baseFeeWei acts as a maximum feeCap cap (0 or nil = no cap).
// tipFeeWei overrides the network-suggested tip (0 or nil = use network tip).
func (pool *TxPool) getSuggestedFeesAdaptive(client *Client, baseFeeWei *big.Int, tipFeeWei *big.Int) (feeCap *big.Int, tipCap *big.Int, err error) {
	if tipFeeWei != nil && tipFeeWei.Sign() > 0 {
		tipCap = new(big.Int).Set(tipFeeWei)
	}

	if tipCap == nil {
		if client == nil {
			return nil, nil, fmt.Errorf("client required for adaptive fee strategy when tipFee not provided")
		}
		_, networkTipCap, fetchErr := client.GetSuggestedFee(pool.options.Context)
		if fetchErr != nil {
			return nil, nil, fetchErr
		}
		tipCap = networkTipCap
	}

	// Dynamic fee calculation: compute feeCap from current chain state.
	// eth_gasPrice returns baseFee + tipCap, which provides zero headroom for
	// baseFee increases. After a full block, baseFee rises 12.5% and ALL pending
	// tx become underpriced, causing alternating full/empty blocks.
	//
	// Dynamic headroom based on recent block fill over a sliding window:
	//   headroom = baseFee * (totalGasLimit - totalGasUsed) / (4 * totalGasLimit)
	//
	// This scales linearly with how underfull recent blocks are:
	//   50% avg fill (target) → baseFee/8  headroom (1 block of max EIP-1559 change)
	//   100% avg fill         → clamped to baseFee/16 minimum (prevents oscillation)
	//   0% avg fill           → baseFee/4  headroom (2 blocks of max change)
	currentBaseFee := pool.GetCurrentBaseFee()
	if currentBaseFee != nil && currentBaseFee.Sign() > 0 {
		totalGasUsed, totalGasLimit := pool.GetRecentBlockGasTotals()

		var headroom *big.Int
		if totalGasLimit > 0 && totalGasUsed <= totalGasLimit {
			headroom = new(big.Int).Set(currentBaseFee)
			headroom.Mul(headroom, new(big.Int).SetUint64(totalGasLimit-totalGasUsed))
			headroom.Div(headroom, new(big.Int).SetUint64(4*totalGasLimit))
		} else {
			headroom = new(big.Int).Div(new(big.Int).Set(currentBaseFee), big.NewInt(8))
		}

		// Clamp to [baseFee/16, baseFee/4]
		minH := new(big.Int).Div(new(big.Int).Set(currentBaseFee), big.NewInt(16))
		maxH := new(big.Int).Div(new(big.Int).Set(currentBaseFee), big.NewInt(4))
		if headroom.Cmp(minH) < 0 {
			headroom = minH
		}
		if headroom.Cmp(maxH) > 0 {
			headroom = maxH
		}

		// Normal distribution centered on the AVERAGE baseFee from the window.
		// When baseFee is rising, the average lags behind current, shifting more
		// tx below currentBaseFee → rejected by EL → nonce gaps → wallet stalls →
		// reduced throughput → baseFee recovers. Self-correcting feedback loop.
		// sigma = headroom / 1.28 → ~10% canary tx at steady state.
		avgBaseFee := pool.getWindowAvgBaseFee()
		if avgBaseFee == nil {
			avgBaseFee = currentBaseFee
		}

		sigma := new(big.Int).Mul(headroom, big.NewInt(25))
		sigma.Div(sigma, big.NewInt(32)) // ≈ headroom / 1.28

		sample := rand.NormFloat64()
		scaledSample := int64(math.Round(sample * 10000))
		offset := new(big.Int).Mul(sigma, big.NewInt(scaledSample))
		offset.Div(offset, big.NewInt(10000))

		feeCap = new(big.Int).Add(avgBaseFee, headroom)
		feeCap.Add(feeCap, offset)

		// Ensure feeCap >= tipCap (EIP-1559 requirement: maxFeePerGas >= maxPriorityFeePerGas)
		if tipCap != nil && feeCap.Cmp(tipCap) < 0 {
			feeCap.Set(tipCap)
		}
		if feeCap.Sign() <= 0 {
			feeCap.SetInt64(1)
		}
	} else if feeCap == nil {
		// Fallback: no tracked baseFee yet, use network suggestion
		if client == nil {
			return nil, nil, fmt.Errorf("client required for adaptive fee strategy when baseFee not yet tracked")
		}
		feeCap, _, _ = client.GetSuggestedFee(pool.options.Context)
	}

	// Apply baseFeeWei as a maximum cap if provided.
	// Scenario defaults (e.g., 20 gwei) act as a safety ceiling — the dynamic
	// mechanism runs freely below it, and the cap only triggers if baseFee
	// rises unexpectedly high.
	if feeCap != nil && baseFeeWei != nil && baseFeeWei.Sign() > 0 && feeCap.Cmp(baseFeeWei) > 0 {
		feeCap.Set(baseFeeWei)
	}

	return feeCap, tipCap, nil
}

// ResolveFees resolves fee values with precedence: Wei > Gwei > nil (network fetch).
// Wei values are parsed from strings to support precise sub-Gwei fees on L2s.
func ResolveFees(baseFeeGwei, tipFeeGwei float64, baseFeeWei, tipFeeWei string) (baseFee, tipFee *big.Int) {
	// Parse Wei strings first (highest precedence)
	if baseFeeWei != "" && baseFeeWei != "0" {
		baseFee, _ = new(big.Int).SetString(baseFeeWei, 10)
	}
	if tipFeeWei != "" && tipFeeWei != "0" {
		tipFee, _ = new(big.Int).SetString(tipFeeWei, 10)
	}

	// Fall back to Gwei conversion if Wei not provided
	if baseFee == nil && baseFeeGwei > 0 {
		baseFee = new(big.Int).SetUint64(uint64(baseFeeGwei * 1e9))
	}
	if tipFee == nil && tipFeeGwei > 0 {
		tipFee = new(big.Int).SetUint64(uint64(tipFeeGwei * 1e9))
	}

	return baseFee, tipFee
}

// calculateBackoffDelay calculates the exponential backoff delay for rebroadcast attempts.
// Uses 30s base delay, 1.5x multiplier, with 10min maximum delay.
func (pool *TxPool) calculateBackoffDelay(retryCount uint64) time.Duration {
	const (
		baseDelay  = 20 * time.Second
		multiplier = 1.5
		maxDelay   = 5 * time.Minute
	)

	delay := time.Duration(float64(baseDelay) * math.Pow(multiplier, float64(retryCount)))
	if delay > maxDelay {
		delay = maxDelay
	}
	return delay
}

// rebroadcastStartOffset picks the client a rebroadcast attempt starts at.
//
// The transaction hash spreads the starting point across the client pool, the way
// submitTransaction does, and the retry count advances it, so every attempt begins at the
// next client rather than at the one the previous attempt already reached. That matters
// because the loop stops at the first client that takes the transaction and a client
// already holding it takes it again: without the retry offset a stuck transaction is
// re-sent to the same node forever. Observed on frames-devnet-0, where reth keeps frame
// transactions in its queued set and never includes or forwards them -- every retry went
// back into the same node and the wallet stayed blocked.
func rebroadcastStartOffset(configured int, hash common.Hash, retryCount uint64) int {
	offset := configured
	if offset == 0 {
		offset = int(hash[0]) | (int(hash[1]) << 8)
	}

	return offset + int(retryCount)
}

// rebroadcastTransaction performs the actual rebroadcast of a transaction.
// With fanout it goes to every eligible client: peers that already saw the hash do not
// fetch it again from gossip, so a node that dropped it only gets it back directly.
// Without, it stops at the first client that takes it.
func (pool *TxPool) rebroadcastTransaction(ctx context.Context, tx *txtypes.Transaction, options *SendTransactionOptions, retryCount uint64, fanout bool) {
	submitTx := func(client *Client) error {
		if options.hasEncodeFn() {
			txBytes, err := options.invokeEncode(tx)
			if err != nil {
				return fmt.Errorf("failed to encode transaction: %w", err)
			}

			if len(txBytes) > 0 {
				return client.SendRawTransaction(ctx, txBytes)
			}
		}

		return client.SendTransaction(ctx, tx)
	}

	clientCount := len(pool.options.ClientPool.GetAllGoodClients())
	if len(options.ClientList) > 0 {
		clientCount = len(options.ClientList)
	}
	if !fanout && clientCount > 5 {
		clientCount = 5
	}

	startOffset := rebroadcastStartOffset(options.ClientsStartOffset, tx.Hash(), retryCount)

	for j := 0; j < clientCount; j++ {
		if ctx.Err() != nil {
			break
		}

		var client *Client
		if len(options.ClientList) > 0 {
			client = options.ClientList[(j+startOffset)%len(options.ClientList)]
		} else {
			var clientOpts []ClientSelectionOption
			if options.ClientGroup != "" {
				clientOpts = append(clientOpts, WithClientGroup(options.ClientGroup))
			}
			clientOpts = append(clientOpts, WithClientSelectionMode(SelectClientByIndex, j+startOffset+1))
			client = pool.options.ClientPool.GetClient(clientOpts...)
		}

		if client == nil {
			break
		}

		err := submitTx(client)

		if options.LogFn != nil {
			options.LogFn(client, j, int(retryCount), err)
		}

		if !fanout && (err == nil || strings.Contains(err.Error(), "already known")) {
			break
		}
	}
}

// blockTime returns the timestamp of a processed block.
//
// The newest block still tracked answers for one that has already been pruned, since a
// later timestamp only makes a deadline more certainly past. Zero means no block has been
// processed yet, and nothing is judged expired on that.
func (pool *TxPool) blockTime(blockNumber uint64) uint64 {
	pool.blocksMutex.RLock()
	defer pool.blocksMutex.RUnlock()

	if info := pool.blocks[blockNumber]; info != nil {
		return info.Timestamp
	}

	newest := uint64(0)
	for _, info := range pool.blocks {
		if info != nil && info.Timestamp > newest {
			newest = info.Timestamp
		}
	}

	return newest
}

// isExpiredTx reports whether a transaction carries a deadline the chain has already
// passed.
//
// The comparison is against the head block's timestamp rather than the local clock: a
// later block can only carry a later timestamp, so a deadline that the head has reached
// can never be met again. Without a known chain time nothing is treated as expired.
func isExpiredTx(tx *txtypes.Transaction, chainTime uint64) bool {
	if chainTime == 0 {
		return false
	}

	deadline, ok := tx.ExpiryDeadline()

	return ok && deadline <= chainTime
}

// replacementFeeBumpPercent is how far a replacement transaction is priced above the one
// it replaces. Clients require a bump before accepting a replacement at the same nonce --
// 10% in geth -- and the extra margin covers the ones that ask for more.
const replacementFeeBumpPercent = 25

// replacementFees prices a transaction meant to replace another one at the same nonce.
// The bump grows with the attempt count, so a replacement a client refuses as underpriced
// is not retried at the same price forever.
func (pool *TxPool) replacementFees(replaced *txtypes.Transaction, attempt uint64) (gasFeeCap, gasTipCap *big.Int) {
	if attempt < 1 {
		attempt = 1
	}

	percent := big.NewInt(100 + replacementFeeBumpPercent*int64(attempt))

	bump := func(value *big.Int) *big.Int {
		if value == nil {
			return big.NewInt(0)
		}

		bumped := new(big.Int).Mul(value, percent)

		return bumped.Div(bumped, big.NewInt(100))
	}

	gasFeeCap = bump(replaced.GasFeeCap())
	gasTipCap = bump(replaced.GasTipCap())

	// The bump alone can still leave the replacement below what the chain currently
	// costs, which would strand it just as thoroughly as the transaction it replaces.
	// The base fee is unknown only before the first block was processed, and stale
	// processing is driven by block processing, so the guard is for the helper's sake.
	if baseFee := pool.GetCurrentBaseFee(); baseFee != nil {
		minFeeCap := new(big.Int).Mul(baseFee, big.NewInt(2))
		if gasFeeCap.Cmp(minFeeCap) < 0 {
			gasFeeCap = minFeeCap
		}
	}

	minTipCap := big.NewInt(1e9)
	if gasTipCap.Cmp(minTipCap) < 0 {
		gasTipCap = minTipCap
	}

	if gasTipCap.Cmp(gasFeeCap) > 0 {
		gasTipCap = new(big.Int).Set(gasFeeCap)
	}

	return gasFeeCap, gasTipCap
}

// replaceExpiredTransaction takes over the nonce of a transaction whose deadline passed.
//
// Such a transaction is not slow, it is finished: an EIP-8141 expiry verifier frame
// reverts once its deadline is behind the chain, so no amount of rebroadcasting will ever
// get it included -- while it keeps occupying its nonce, and every later transaction the
// wallet sends queues behind it for the rest of the run. A dummy transaction at the same
// nonce, priced above the original so the pools take the replacement, retires it the same
// way fillNonceGaps retires a gap.
func (pool *TxPool) replaceExpiredTransaction(ctx context.Context, wallet *Wallet, expired *PendingTx, attempt uint64) {
	nonce := expired.Tx.Nonce()
	deadline, _ := expired.Tx.ExpiryDeadline()

	gasFeeCap, gasTipCap := pool.replacementFees(expired.Tx, attempt)

	fillerTx, err := wallet.BuildFillerTx(nonce, gasTipCap, gasFeeCap, pool.MinIntrinsicGas())
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"wallet": wallet.GetAddress().Hex(),
			"nonce":  nonce,
		}).WithError(err).Warnf("failed to build replacement for expired transaction")

		return
	}

	options := &SendTransactionOptions{
		Rebroadcast: true,
		OnComplete: func(tx *txtypes.Transaction, receipt *txtypes.Receipt, err error) {
			fields := logrus.Fields{
				"wallet": wallet.GetAddress().Hex(),
				"nonce":  tx.Nonce(),
				"txhash": tx.Hash().Hex(),
			}

			switch {
			case err != nil:
				logrus.WithFields(fields).WithError(err).Warnf("replacement for expired transaction failed")
			case receipt != nil:
				logrus.WithFields(fields).Infof("replacement for expired transaction confirmed")
			}
		},
	}
	if expired.Options != nil {
		options.ClientGroup = expired.Options.ClientGroup
	}

	logrus.WithFields(logrus.Fields{
		"wallet":   wallet.GetAddress().Hex(),
		"nonce":    nonce,
		"txhash":   expired.Tx.Hash().Hex(),
		"deadline": deadline,
		"attempt":  attempt,
	}).Warnf("transaction expired in the mempool, replacing it to unblock the nonce")

	if err := pool.submitTransaction(ctx, wallet, fillerTx, options, true); err != nil {
		logrus.WithFields(logrus.Fields{
			"wallet": wallet.GetAddress().Hex(),
			"nonce":  nonce,
		}).WithError(err).Warnf("failed to submit replacement for expired transaction")
	}
}

// fillNonceGaps creates and submits dummy transactions to fill nonce gaps.
func (pool *TxPool) fillNonceGaps(ctx context.Context, wallet *Wallet, gaps []uint64, baseOptions *SendTransactionOptions) {
	baseFee := pool.GetCurrentBaseFee()

	// Use 2x base fee for fee cap and reasonable tip
	gasFeeCap := new(big.Int).Mul(baseFee, big.NewInt(2))
	gasTipCap := big.NewInt(1e9) // 1 gwei tip
	if gasFeeCap.Cmp(gasTipCap) < 0 {
		gasTipCap = gasFeeCap
	}

	for _, nonce := range gaps {
		select {
		case <-ctx.Done():
			return
		default:
		}

		// Build and sign the filler transaction
		fillerTx, err := wallet.BuildFillerTx(nonce, gasTipCap, gasFeeCap, pool.MinIntrinsicGas())
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"wallet": wallet.GetAddress().Hex(),
				"nonce":  nonce,
			}).WithError(err).Warnf("failed to build filler transaction")
			continue
		}

		// Create options for the filler transaction
		fillerOptions := &SendTransactionOptions{
			Rebroadcast: true,
			OnComplete: func(tx *txtypes.Transaction, receipt *txtypes.Receipt, err error) {
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"wallet": wallet.GetAddress().Hex(),
						"nonce":  tx.Nonce(),
					}).WithError(err).Warnf("filler transaction failed")
				} else {
					logrus.WithFields(logrus.Fields{
						"wallet": wallet.GetAddress().Hex(),
						"nonce":  tx.Nonce(),
						"txhash": tx.Hash().Hex(),
					}).Infof("filler transaction confirmed")
				}
			},
		}

		if baseOptions != nil {
			fillerOptions.ClientGroup = baseOptions.ClientGroup
		}

		// Submit the filler transaction
		err = pool.submitTransaction(ctx, wallet, fillerTx, fillerOptions, true)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"wallet": wallet.GetAddress().Hex(),
				"nonce":  nonce,
			}).WithError(err).Warnf("failed to submit filler transaction")
		}
	}
}
