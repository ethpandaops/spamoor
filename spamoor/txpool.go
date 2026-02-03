package spamoor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethpandaops/spamoor/utils"
	"github.com/sirupsen/logrus"
)

// BlockInfo represents information about a processed block including
// hash, parent hash, gas limit, and timestamp for chain reorganization detection.
type BlockInfo struct {
	Number     uint64
	Hash       common.Hash
	ParentHash common.Hash
	Timestamp  uint64
	GasLimit   uint64
}

// BlockWithHash represents a block with its hash.
type BlockWithHash struct {
	Hash  common.Hash
	Block *types.Block
}

// TxInfo represents information about a confirmed transaction including
// the transaction details, associated wallets, and send options.
type TxInfo struct {
	TxHash     common.Hash
	From       common.Address
	To         *common.Address
	Tx         *types.Transaction
	TxFees     *utils.TxFees
	FromWallet *Wallet
	ToWallet   *Wallet
	Options    *SendTransactionOptions
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
	Block            *types.Block
	ConfirmedTxs     []*TxInfo
	Receipts         []*types.Receipt
}

type GlobalBlockStats struct {
	Block           *types.Block
	Receipts        []*types.Receipt
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

	// Block update subscriptions
	subscriptionsMutex sync.RWMutex
	subscriptions      map[uint64]*BlockSubscription
	bulkSubscriptions  map[uint64]*BulkBlockSubscription
	nextSubscriptionID atomic.Uint64
}

type txPoolWalletRegistration struct {
	ctxs   []context.Context
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
}

type ExternalBlockSource struct {
	SubscribeBlocks func(ctx context.Context, capacity int) chan *ExternalBlockEvent
}

type ExternalBlockEvent struct {
	Number   uint64
	Clients  []*Client
	Block    *BlockWithHash
	Receipts []*types.Receipt
}

// NewTxPool creates a new transaction pool with the specified options.
// It starts background goroutines for block processing and stale transaction handling.
// The pool automatically begins monitoring the blockchain for new blocks and managing
// transaction confirmations and reorgs.
func NewTxPool(options *TxPoolOptions) *TxPool {
	pool := &TxPool{
		options:           options,
		wallets:           make(map[common.Address]*txPoolWalletRegistration),
		walletPools:       make(map[*WalletPool]struct{}),
		processStaleChan:  make(chan uint64, 1),
		blocks:            map[uint64]*BlockInfo{},
		confirmedTxs:      map[uint64][]*TxInfo{},
		reorgDepth:        10, // Default value
		subscriptions:     map[uint64]*BlockSubscription{},
		bulkSubscriptions: map[uint64]*BulkBlockSubscription{},
	}

	if options.Context == nil {
		options.Context = context.Background()
	}

	if options.ReorgDepth > 0 {
		pool.reorgDepth = options.ReorgDepth
	}

	go pool.runTxPoolLoop()
	go pool.processStaleTransactionsLoop()

	return pool
}

// RegisterWallet registers a wallet with the transaction pool.
// It is used to track the wallets that are active and need to be processed.
// The registration is removed when the context is cancelled.
func (pool *TxPool) RegisterWallet(wallet *Wallet, ctx context.Context) *Wallet {
	if ctx.Err() != nil {
		return nil
	}

	pool.walletsMutex.Lock()
	defer pool.walletsMutex.Unlock()

	if registration, ok := pool.wallets[wallet.address]; ok {
		if registration.wallet.privkey == nil && wallet.privkey != nil {
			registration.wallet.privkey = wallet.privkey
		}

		for _, regctx := range registration.ctxs {
			if regctx == ctx {
				return registration.wallet
			}
		}

		registration.ctxs = append(registration.ctxs, ctx)
		return registration.wallet
	}

	pool.wallets[wallet.address] = &txPoolWalletRegistration{
		ctxs:   []context.Context{ctx},
		wallet: wallet,
	}

	return wallet
}

// RegisterWalletPool registers a wallet pool with the transaction pool.
// It is used to track the wallet pools that are active and need to be processed.
func (pool *TxPool) RegisterWalletPool(walletPool *WalletPool) {
	if walletPool.ctx.Err() != nil {
		return
	}

	pool.walletsMutex.Lock()
	defer pool.walletsMutex.Unlock()

	pool.walletPools[walletPool] = struct{}{}
}

// GetRegisteredWallet returns a wallet by address.
func (pool *TxPool) GetRegisteredWallet(address common.Address) *Wallet {
	pool.walletsMutex.RLock()
	defer pool.walletsMutex.RUnlock()

	registration, found := pool.wallets[address]
	if found {
		for {
			if registration.ctxs[0].Err() != nil {
				registration.ctxs = registration.ctxs[1:]
				if len(registration.ctxs) == 0 {
					delete(pool.wallets, registration.wallet.address)
					break
				}
			} else {
				break
			}
		}

		if len(registration.ctxs) > 0 {
			return registration.wallet
		}
	}

	return nil
}

// runTxPoolLoop continuously monitors for new blocks and processes them.
// It tracks the highest block number across all clients and processes new blocks
// sequentially. Also triggers stale transaction processing when new blocks arrive.
// Runs until the context is cancelled and recovers from panics with logging.
func (pool *TxPool) runTxPoolLoop() {
	defer func() {
		utils.RecoverPanic(pool.options.Logger, "TxPool.runTxPoolLoop", pool.runTxPoolLoop)
	}()

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
					logrus.Errorf("failed to process block %v", blockNumber)
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
	defer func() {
		utils.RecoverPanic(pool.options.Logger, "TxPool.runExternalBlockSourceLoop", pool.runExternalBlockSourceLoop)
	}()

	blockChan := pool.options.ExternalBlockSource.SubscribeBlocks(pool.options.Context, 10)
	highestBlockNumber := uint64(0)
	for {
		select {
		case <-pool.options.Context.Done():
			return
		case blockEvent := <-blockChan:

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
						logrus.Errorf("failed to process block %v", blockNumber)
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
	defer func() {
		utils.RecoverPanic(pool.options.Logger, "TxPool.processStaleTransactionsLoop", pool.processStaleTransactionsLoop)
	}()

	for {
		select {
		case <-pool.options.Context.Done():
			return
		case blockNumber := <-pool.processStaleChan:
			for _, wallet := range pool.getWalletMap() {
				pool.processStaleConfirmations(blockNumber, wallet)
			}
		}
	}
}

// getWalletMap collects all registered wallets into a single map.
func (pool *TxPool) getWalletMap() map[common.Address]*Wallet {
	pool.walletsMutex.RLock()
	defer pool.walletsMutex.RUnlock()

	walletMap := make(map[common.Address]*Wallet)
	for _, registration := range pool.wallets {
		walletMap[registration.wallet.address] = registration.wallet
	}
	return walletMap
}

// processBlock processes a single block for transaction confirmations and reorg detection.
// It loads the block body, checks for chain reorganizations by comparing parent hashes,
// stores block information for reorg tracking, and processes all transactions in the block.
// Also handles cleanup of old block data based on the reorg depth setting.
func (pool *TxPool) processBlock(ctx context.Context, client *Client, blockNumber uint64, blockWithHash *BlockWithHash, receipts []*types.Receipt) error {
	pool.lastBlockNumber = blockNumber

	walletMap := pool.getWalletMap()
	if len(walletMap) == 0 {
		return nil
	}

	txSkipMap := make(map[uint32]bool)
	if blockWithHash == nil || blockWithHash.Block == nil {
		blockWithHash, txSkipMap = pool.getBlockBody(ctx, client, blockNumber)
		if blockWithHash == nil {
			return fmt.Errorf("could not load block body")
		}
	}

	// Check for reorg by comparing parent hash
	pool.blocksMutex.RLock()
	lastBlock, hasLastBlock := pool.blocks[blockNumber-1]
	pool.blocksMutex.RUnlock()

	if hasLastBlock && lastBlock.Hash != blockWithHash.Block.ParentHash() {
		logrus.Warnf("Detected chain reorganization at block %d. Parent hash mismatch: expected %s, got %s",
			blockNumber, lastBlock.Hash.Hex(), blockWithHash.Block.ParentHash().Hex())

		// Handle reorg
		pool.handleReorg(ctx, client, blockNumber, blockWithHash, walletMap)
	}

	// Store block info
	pool.blocksMutex.Lock()
	pool.blocks[blockNumber] = &BlockInfo{
		Number:     blockNumber,
		Hash:       blockWithHash.Hash,
		ParentHash: blockWithHash.Block.ParentHash(),
		Timestamp:  blockWithHash.Block.Time(),
		GasLimit:   blockWithHash.Block.GasLimit(),
	}

	// Update current gas limit
	pool.blockStatsMutex.Lock()
	pool.currentGasLimit = blockWithHash.Block.GasLimit()
	pool.currentBaseFee = blockWithHash.Block.BaseFee()
	pool.blockStatsMutex.Unlock()

	// Clean up old blocks
	if blockNumber > uint64(pool.reorgDepth) {
		delete(pool.blocks, blockNumber-uint64(pool.reorgDepth))
	}
	pool.blocksMutex.Unlock()

	return pool.processBlockTxs(ctx, client, blockNumber, blockWithHash, walletMap, txSkipMap, receipts)
}

// processBlockTxs processes all transactions in a block for confirmation tracking.
// It loads block receipts, decodes transaction senders, updates wallet states for
// confirmed transactions, and tracks transaction information for reorg recovery.
// Also handles cleanup of old confirmed transaction data.
func (pool *TxPool) processBlockTxs(ctx context.Context, client *Client, blockNumber uint64, blockWithHash *BlockWithHash, walletMap map[common.Address]*Wallet, txSkipMap map[uint32]bool, receipts []*types.Receipt) error {
	t1 := time.Now()
	txCount := len(blockWithHash.Block.Transactions())
	if receipts == nil {
		var err error
		receipts, err = pool.getBlockReceipts(ctx, client, blockWithHash.Hash, txCount, txSkipMap)
		if err != nil {
			return fmt.Errorf("could not load block receipts: %w", err)
		}
	}

	loadingTime := time.Since(t1)
	t1 = time.Now()

	signer := types.LatestSignerForChainID(pool.options.ChainId)
	confirmCount := 0
	affectedWalletMap := map[common.Address]bool{}
	pool.txsMutex.Lock()
	pool.confirmedTxs[blockNumber] = []*TxInfo{}
	pool.txsMutex.Unlock()

	for idx, tx := range blockWithHash.Block.Transactions() {
		receipt := receipts[idx]
		if receipt == nil {
			logrus.Warnf("missing receipt for tx %v in block %v", idx, blockNumber)
			continue
		}

		txFrom, err := types.Sender(signer, tx)
		if err != nil {
			logrus.Warnf("error decoding tx sender (block %v, tx %v): %v", blockNumber, idx, err)
			continue
		}

		txHash := tx.Hash()
		txFees := utils.GetTransactionFees(tx, receipt)
		fromWallet := walletMap[txFrom]
		toAddr := tx.To()
		toWallet := (*Wallet)(nil)
		if toAddr != nil {
			toWallet = walletMap[*toAddr]
		}

		if fromWallet != nil || toWallet != nil {
			pool.txsMutex.Lock()
			pool.confirmedTxs[blockNumber] = append(pool.confirmedTxs[blockNumber], &TxInfo{
				TxHash:     txHash,
				From:       txFrom,
				To:         tx.To(),
				Tx:         tx,
				TxFees:     txFees,
				FromWallet: fromWallet,
				ToWallet:   toWallet,
			})
			pool.txsMutex.Unlock()
		}

		if fromWallet != nil {
			confirmCount++
			affectedWalletMap[txFrom] = true
			pool.processTransactionInclusion(blockNumber, fromWallet, tx, receipt, txFees)
		}

		if toWallet != nil {
			toWallet.AddBalance(tx.Value())
		}
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
			blockNumber, err := client.client.BlockNumber(ctx)
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
// It uses a 5-second timeout and returns the block if successful, nil otherwise.
// Builder clients are not supported as they don't provide eth_getBlockByNumber.
func (pool *TxPool) getBlockBody(ctx context.Context, client *Client, blockNumber uint64) (*BlockWithHash, map[uint32]bool) {
	// Builder clients don't support eth_getBlockByNumber
	if client.IsBuilder() {
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	/*
		blockBody, err := client.client.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
		if err == nil {
			return blockBody
		}
	*/

	var raw json.RawMessage
	err := client.client.Client().CallContext(ctx, &raw, "eth_getBlockByNumber", rpc.BlockNumber(blockNumber), true)
	if err != nil {
		return nil, nil
	}

	// Decode header and transactions.
	var head *types.Header
	if err := json.Unmarshal(raw, &head); err != nil {
		return nil, nil
	}
	// When the block is not found, the API returns JSON null.
	if head == nil {
		return nil, nil
	}

	var body struct {
		Hash         common.Hash       `json:"hash"`
		Transactions []json.RawMessage `json:"transactions"`
	}
	if err := json.Unmarshal(raw, &body); err != nil {
		return nil, nil
	}

	transactions := make([]*types.Transaction, 0, len(body.Transactions))
	txSkipMap := make(map[uint32]bool)
	for idx, rawTx := range body.Transactions {
		var txHeader struct {
			Type hexutil.Uint64 `json:"type"`
		}

		isValid := false

		if err := json.Unmarshal(rawTx, &txHeader); err == nil {
			switch txHeader.Type {
			case types.LegacyTxType, types.AccessListTxType, types.DynamicFeeTxType, types.BlobTxType, types.SetCodeTxType:
				isValid = true
			}
		}

		if isValid {
			var tx types.Transaction
			if err := json.Unmarshal(rawTx, &tx); err != nil {
				isValid = false
			}
			transactions = append(transactions, &tx)
		}

		if !isValid {
			txSkipMap[uint32(idx)] = true
			continue
		}
	}

	block := types.NewBlockWithHeader(head).WithBody(
		types.Body{
			Transactions: transactions,
		},
	)

	if body.Hash.Cmp(common.Hash{}) == 0 {
		body.Hash = block.Hash()
	}

	return &BlockWithHash{
		Hash:  body.Hash,
		Block: block,
	}, txSkipMap
}

// getBlockReceipts retrieves all transaction receipts for a block.
// It validates that the number of receipts matches the expected transaction count
// and uses a 5-second timeout for the request.
// Builder clients are not supported as they don't provide eth_getBlockReceipts.
func (pool *TxPool) getBlockReceipts(ctx context.Context, client *Client, blockHash common.Hash, txCount int, txSkipMap map[uint32]bool) ([]*types.Receipt, error) {
	// Builder clients don't support eth_getBlockReceipts
	if client.IsBuilder() {
		return nil, fmt.Errorf("builder clients do not support eth_getBlockReceipts")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var receiptErr error

	blockReceipts, err := client.client.BlockReceipts(ctx, rpc.BlockNumberOrHash{
		BlockHash: &blockHash,
	})
	if err != nil {
		receiptErr = err
	} else {
		if len(txSkipMap) > 0 {
			filteredBlockReceipts := make([]*types.Receipt, 0, txCount)
			for idx, receipt := range blockReceipts {
				if !txSkipMap[uint32(idx)] {
					filteredBlockReceipts = append(filteredBlockReceipts, receipt)
				}
			}
			blockReceipts = filteredBlockReceipts
		}

		if len(blockReceipts) != txCount {
			return nil, fmt.Errorf("block %v has %v receipts, expected %v", blockHash.Hex(), len(blockReceipts), txCount)
		}

		return blockReceipts, nil
	}

	return nil, receiptErr
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
func (pool *TxPool) notifyBlockSubscribers(blockNumber uint64, block *types.Block, receipts []*types.Receipt) {
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

	confirmedTxMap := make(map[common.Hash]*TxInfo)
	for _, txInfo := range blockConfirmedTxs {
		confirmedTxMap[txInfo.TxHash] = txInfo
	}

	// If we have bulk subscribers, calculate stats for all wallet pools once
	var globalBlockStats *GlobalBlockStats
	if len(bulkSubscriptions) > 0 {
		globalBlockStats = pool.calculateAllWalletPoolStats(confirmedTxMap, block, receipts)

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
		} else {
			// Calculate stats individually
			stats = pool.calculateWalletPoolStats(subscription.WalletPool, confirmedTxMap, block, receipts)
		}
		subscription.Callback(blockNumber, stats)
	}
}

// calculateWalletPoolStats calculates transaction statistics for a specific wallet pool.
func (pool *TxPool) calculateWalletPoolStats(walletPool *WalletPool, confirmedTxMap map[common.Hash]*TxInfo, block *types.Block, receipts []*types.Receipt) *WalletPoolBlockStats {
	stats := &WalletPoolBlockStats{
		TotalTxFees:  big.NewInt(0),
		Block:        block,
		Receipts:     make([]*types.Receipt, 0, len(confirmedTxMap)),
		ConfirmedTxs: make([]*TxInfo, 0, len(confirmedTxMap)),
	}

	affectedWallets := make(map[common.Address]bool)
	allWallets := walletPool.GetAllWallets()
	walletSet := make(map[common.Address]bool)

	for _, wallet := range allWallets {
		if wallet == nil {
			continue
		}

		walletSet[wallet.GetAddress()] = true
	}

	for idx, tx := range block.Transactions() {
		txInfo, ok := confirmedTxMap[tx.Hash()]
		if !ok {
			continue
		}

		if txInfo.FromWallet != nil && walletSet[txInfo.FromWallet.GetAddress()] {
			stats.ConfirmedTxCount++
			if txInfo.TxFees != nil {
				totalFee := new(big.Int).Add(&txInfo.TxFees.FeeAmount, &txInfo.TxFees.BlobFeeAmount)
				stats.TotalTxFees.Add(stats.TotalTxFees, totalFee)
			}
			affectedWallets[txInfo.FromWallet.GetAddress()] = true

			stats.Receipts = append(stats.Receipts, receipts[idx])
			stats.ConfirmedTxs = append(stats.ConfirmedTxs, txInfo)
		}
	}

	stats.AffectedWallets = len(affectedWallets)

	return stats
}

// calculateAllWalletPoolStats efficiently calculates statistics for ALL active wallet pools in a single pass.
// This avoids recalculating the same data multiple times for bulk subscribers.
func (pool *TxPool) calculateAllWalletPoolStats(confirmedTxMap map[common.Hash]*TxInfo, block *types.Block, receipts []*types.Receipt) *GlobalBlockStats {
	// Get all active wallet pools
	activeWalletPools := pool.GetActiveWalletPools()

	// Initialize result map
	allStats := make(map[*WalletPool]*WalletPoolBlockStats, len(activeWalletPools))

	// Create wallet address to wallet pool mapping for efficient lookup
	walletToPool := make(map[common.Address]*WalletPool)
	for _, walletPool := range activeWalletPools {
		allWallets := walletPool.GetAllWallets()
		for _, wallet := range allWallets {
			if wallet == nil {
				continue
			}

			walletToPool[wallet.GetAddress()] = walletPool
		}

		// Initialize stats for each wallet pool
		allStats[walletPool] = &WalletPoolBlockStats{
			TotalTxFees:  big.NewInt(0),
			Block:        block,
			Receipts:     make([]*types.Receipt, 0, len(confirmedTxMap)),
			ConfirmedTxs: make([]*TxInfo, 0, len(confirmedTxMap)),
		}
	}

	// Track affected wallets per pool
	affectedWallets := make(map[*WalletPool]map[common.Address]bool)
	for _, walletPool := range activeWalletPools {
		affectedWallets[walletPool] = make(map[common.Address]bool)
	}

	// Single pass through confirmed transactions
	for idx, tx := range block.Transactions() {
		txInfo, ok := confirmedTxMap[tx.Hash()]
		if !ok {
			continue
		}

		if txInfo.FromWallet != nil {
			walletAddr := txInfo.FromWallet.GetAddress()
			if walletPool, exists := walletToPool[walletAddr]; exists {
				stats := allStats[walletPool]

				// Update transaction count
				stats.ConfirmedTxCount++

				// Update total fees
				if txInfo.TxFees != nil {
					totalFee := new(big.Int).Add(&txInfo.TxFees.FeeAmount, &txInfo.TxFees.BlobFeeAmount)
					stats.TotalTxFees.Add(stats.TotalTxFees, totalFee)
				}

				affectedWallets[walletPool][walletAddr] = true
				stats.Receipts = append(stats.Receipts, receipts[idx])
				stats.ConfirmedTxs = append(stats.ConfirmedTxs, txInfo)
			}
		}
	}

	// Set affected wallet counts
	for walletPool, stats := range allStats {
		stats.AffectedWallets = len(affectedWallets[walletPool])
	}

	return &GlobalBlockStats{
		Block:           block,
		Receipts:        receipts,
		WalletPoolStats: allStats,
	}
}

// submitTransaction handles the core transaction submission logic.
// It starts a confirmation tracking goroutine, submits the transaction to clients,
// and optionally sets up automatic rebroadcasting. The submitNow parameter controls
// whether to immediately submit or just set up confirmation tracking.
func (pool *TxPool) submitTransaction(ctx context.Context, wallet *Wallet, tx *types.Transaction, options *SendTransactionOptions, submitNow bool) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	confirmCtx, confirmCancel := context.WithCancel(ctx)
	wg := &sync.WaitGroup{}
	wg.Add(1)

	submissionComplete := make(chan error, 1)

	go func() {
		var receipt *types.Receipt

		var err error

		defer confirmCancel()

		defer func() {
			submissionError := <-submissionComplete
			if submissionError != nil {
				err = submissionError
				// special case: since we never sucessfully submitted the tx, we need to drop it from the pending txs as it is not known to the network
				wallet.dropPendingTx(tx)
			} else if options.OnConfirm != nil && receipt != nil {
				options.OnConfirm(tx, receipt)
			}

			if options.OnComplete != nil {
				options.OnComplete(tx, receipt, err)
			}

			// Track transaction result for metrics
			walletPools := pool.GetActiveWalletPools()
			for _, walletPool := range walletPools {
				if tracker := walletPool.GetTransactionTracker(); tracker != nil {
					// Check if this wallet belongs to this pool
					allWallets := walletPool.GetAllWallets()
					for _, poolWallet := range allWallets {
						if poolWallet == nil {
							continue
						}

						if poolWallet.GetAddress() == wallet.GetAddress() {
							tracker(err)
							break
						}
					}
				}
			}
		}()

		receipt, err = pool.awaitTransaction(confirmCtx, wallet, tx, options, wg)
		if confirmCtx.Err() != nil {
			err = nil
		}

		if receipt != nil {
			pool.txsMutex.Lock()
			for _, tx := range pool.confirmedTxs[receipt.BlockNumber.Uint64()] {
				if tx.TxHash == receipt.TxHash {
					tx.Options = options
					break
				}
			}
			pool.txsMutex.Unlock()
		}
	}()

	wg.Wait()

	var submitErr error

	submitTx := func(client *Client) error {
		if options.OnEncode != nil {
			txBytes, err := options.OnEncode(tx)
			if err != nil {
				return fmt.Errorf("failed to encode transaction: %w", err)
			}

			if len(txBytes) > 0 {
				return client.SendRawTransaction(ctx, txBytes)
			}
		}

		return client.SendTransaction(ctx, tx)
	}

	if submitNow {
		submitCount := options.SubmitCount
		if submitCount == 0 {
			submitCount = 3
		}

		success := false
		clientCount := len(pool.options.ClientPool.GetAllGoodClients())
		for i := 0; i < clientCount; i++ {
			client := options.Client
			if client == nil || i > 0 {
				var clientOpts []ClientSelectionOption
				if options.ClientGroup != "" {
					clientOpts = append(clientOpts, WithClientGroup(options.ClientGroup))
				}
				clientOpts = append(clientOpts, WithClientSelectionMode(SelectClientByIndex, i+options.ClientsStartOffset))
				client = pool.options.ClientPool.GetClient(clientOpts...)
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
			wallet.IncrementSubmittedTxCount()

			submitErr = nil
		}
	}

	submissionComplete <- submitErr

	if submitErr != nil {
		confirmCancel()
		return submitErr
	}

	return nil
}

// awaitTransaction waits for a specific transaction to be confirmed.
// It uses the wallet's nonce channel system to wait for confirmation and
// handles cases where the transaction might be replaced or reorged.
// The wg parameter is signaled when confirmation tracking is set up.
func (pool *TxPool) awaitTransaction(ctx context.Context, wallet *Wallet, tx *types.Transaction, options *SendTransactionOptions, wg *sync.WaitGroup) (*types.Receipt, error) {
	txHash := tx.Hash()
	nonceChan, isFirstPendingTx := wallet.getTxNonceChan(tx, options)

	if isFirstPendingTx && pool.lastBlockNumber > wallet.lastConfirmation+1 {
		wallet.lastConfirmation = pool.lastBlockNumber - 1
	}

	if wg != nil {
		wg.Done()
	}

	if nonceChan != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-nonceChan.channel:
		}

		receipt := nonceChan.receipt
		if receipt != nil {
			if bytes.Equal(receipt.TxHash[:], txHash[:]) {
				return receipt, nil
			}

			return nil, nil
		}
	}

	return pool.loadTransactionReceipt(ctx, tx), nil
}

// processTransactionInclusion handles the confirmation of a transaction from a tracked wallet.
// It updates the wallet's nonce state, signals any waiting confirmation channels,
// and cleans up completed nonce channels. Updates the wallet's confirmation tracking.
// Also updates the wallet's balance by subtracting the transaction fee and blob fee.
func (pool *TxPool) processTransactionInclusion(blockNumber uint64, wallet *Wallet, tx *types.Transaction, receipt *types.Receipt, txFees *utils.TxFees) {
	totalAmount := new(big.Int).Add(tx.Value(), &txFees.FeeAmount)
	totalAmount = new(big.Int).Add(totalAmount, &txFees.BlobFeeAmount)
	wallet.SubBalance(totalAmount)

	wallet.txNonceMutex.Lock()
	defer wallet.txNonceMutex.Unlock()

	nonce := tx.Nonce()
	if nonceChan := wallet.txNonceChans[nonce]; nonceChan != nil {
		nonceChan.receipt = receipt
	}

	wallet.confirmedTxCount = nonce + 1
	if nonce+1 > wallet.pendingTxCount.Load() {
		wallet.pendingTxCount.Store(nonce + 1)
	}
	if blockNumber > wallet.lastConfirmation {
		wallet.lastConfirmation = blockNumber
	}

	for n := range wallet.txNonceChans {
		if n <= nonce {
			close(wallet.txNonceChans[n].channel)
			delete(wallet.txNonceChans, n)
		}
	}

}

// processStaleConfirmations recovers stale transactions that may have been missed
// and handles rebroadcasting of pending transactions.
// It checks if a wallet has pending transactions that are older than 10 blocks,
// queries the current nonce from the blockchain, recovers any confirmed transactions
// that weren't properly tracked, and rebroadcasts pending transactions if enabled.
func (pool *TxPool) processStaleConfirmations(blockNumber uint64, wallet *Wallet) {
	if len(wallet.txNonceChans) == 0 || blockNumber <= wallet.lastConfirmation+10 {
		return
	}

	wallet.lastConfirmation = blockNumber

	// Query on-chain nonce
	var onChainNonce uint64
	var err error
	for retry := 0; retry < 3; retry++ {
		client := pool.options.ClientPool.GetClient(WithClientSelectionMode(SelectClientRandom))
		if client == nil {
			continue
		}

		onChainNonce, err = client.GetNonceAt(pool.options.Context, wallet.address, big.NewInt(int64(blockNumber)))
		if err == nil {
			break
		}
	}

	if err != nil {
		logrus.WithError(err).Warnf("failed to get on-chain nonce for %v", wallet.address.String())
		return
	}

	wallet.txNonceMutex.Lock()

	// Collect pending nonces sorted
	pendingNonces := make([]uint64, 0, len(wallet.txNonceChans))
	for n, nc := range wallet.txNonceChans {
		if len(nc.txs) > 0 {
			pendingNonces = append(pendingNonces, n)
		}
	}
	slices.Sort(pendingNonces)

	// Find lowest pending nonce
	lowestPendingNonce := uint64(0)
	if len(pendingNonces) > 0 {
		lowestPendingNonce = pendingNonces[0]
	}

	name := wallet.GetAddress().String()
	logrus.Debugf("processing stale confirmations for %v (on-chain nonce: %v, confirmed: %v, pending count: %v, lowest pending: %v)",
		name, onChainNonce, wallet.confirmedTxCount, len(wallet.txNonceChans), lowestPendingNonce)

	// Close channels for confirmed transactions and collect txs to rebroadcast
	var txsToRebroadcast []*PendingTx
	var nonceGaps []uint64

	// Only consider rebroadcasting the 2 lowest pending nonces
	const maxRebroadcastNonces = 2
	pendingNoncesChecked := 0

	for _, nonce := range pendingNonces {
		nonceChan := wallet.txNonceChans[nonce]
		if nonce < onChainNonce {
			// Transaction confirmed - close channel and clean up
			logrus.Debugf("recovering stale confirmed transaction for %v (nonce %v)", wallet.address.String(), nonce)
			close(nonceChan.channel)
			delete(wallet.txNonceChans, nonce)
		} else if nonce <= onChainNonce+(maxRebroadcastNonces-1) {
			pendingNoncesChecked++
			// Get the most recent pending tx for this nonce (last in the list)
			if len(nonceChan.txs) > 0 {
				mostRecentTx := nonceChan.txs[len(nonceChan.txs)-1]
				if mostRecentTx.Options != nil && mostRecentTx.Options.Rebroadcast {
					// Check if enough time has passed since last rebroadcast
					backoffDelay := pool.calculateBackoffDelay(mostRecentTx.RebroadcastCount)
					if time.Since(mostRecentTx.LastRebroadcast) >= backoffDelay {
						txsToRebroadcast = append(txsToRebroadcast, mostRecentTx)
					}
				}
			}
		}
	}

	// Update confirmed nonce
	if onChainNonce > wallet.confirmedTxCount {
		wallet.confirmedTxCount = onChainNonce
	}

	// Check for nonce gaps between on-chain nonce and lowest pending nonce
	if lowestPendingNonce > onChainNonce {
		for nonce := onChainNonce; nonce < lowestPendingNonce; nonce++ {
			nonceGaps = append(nonceGaps, nonce)
		}
	}

	wallet.txNonceMutex.Unlock()

	// Fill nonce gaps if any
	if len(nonceGaps) > 0 {
		logrus.WithFields(logrus.Fields{
			"wallet": wallet.GetAddress().Hex(),
			"gaps":   nonceGaps,
		}).Warnf("detected nonce gaps, filling with dummy transactions")
		go pool.fillNonceGaps(pool.options.Context, wallet, nonceGaps, nil)
	}

	// Rebroadcast pending transactions
	for _, pendingTx := range txsToRebroadcast {
		logrus.WithFields(logrus.Fields{
			"wallet": wallet.GetAddress().Hex(),
			"nonce":  pendingTx.Tx.Nonce(),
			"txhash": pendingTx.Tx.Hash().Hex(),
		}).Debugf("rebroadcasting stale transaction")

		pendingTx.LastRebroadcast = time.Now()
		pendingTx.RebroadcastCount++
		go pool.rebroadcastTransaction(pool.options.Context, pendingTx.Tx, pendingTx.Options, pendingTx.RebroadcastCount)
	}
}

// loadTransactionReceipt attempts to load a transaction receipt from multiple clients.
// It retries up to 5 times with different clients and includes exponential backoff.
// Returns nil if the receipt cannot be loaded after all retries.
func (pool *TxPool) loadTransactionReceipt(ctx context.Context, tx *types.Transaction) *types.Receipt {
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

		if ctx.Err() != nil {
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
func (pool *TxPool) handleReorg(ctx context.Context, client *Client, blockNumber uint64, newBlockWithHash *BlockWithHash, walletMap map[common.Address]*Wallet) error {
	type newBlockInfo struct {
		blockWithHash *BlockWithHash
		txSkipMap     map[uint32]bool
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

		if pool.blocks[blockNumber].Hash == block.Block.ParentHash() {
			break
		}

		parentBlockBody, txSkipMap := pool.getBlockBody(ctx, client, blockNumber)
		if parentBlockBody == nil {
			return fmt.Errorf("could not load block body for new parent block %v", blockNumber)
		}

		newBlockParents = append(newBlockParents, newBlockInfo{
			blockWithHash: parentBlockBody,
			txSkipMap:     txSkipMap,
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
	for _, tx := range reorgedOutTxs {
		if tx.FromWallet != nil {
			if !resetWallets[tx.From] {
				resetWallets[tx.From] = true

				tx.FromWallet.txNonceMutex.Lock()
				tx.FromWallet.confirmedTxCount = tx.Tx.Nonce()
				tx.FromWallet.txNonceMutex.Unlock()
			}

			totalAmount := new(big.Int).Add(tx.Tx.Value(), &tx.TxFees.FeeAmount)
			totalAmount = new(big.Int).Add(totalAmount, &tx.TxFees.BlobFeeAmount)
			tx.FromWallet.AddBalance(totalAmount)

			// add tx as pending tx
			txOptions := &SendTransactionOptions{
				Client:      client,
				Rebroadcast: true,
				OnComplete: func(tx *types.Transaction, receipt *types.Receipt, err error) {
					if err == nil {
						logrus.Infof("reorged out tx %v confirmed", tx.Hash())
					}
				},
			}

			if tx.Options != nil {
				txOptions.LogFn = tx.Options.LogFn
				txOptions.ClientGroup = tx.Options.ClientGroup
			}

			err := pool.submitTransaction(ctx, tx.FromWallet, tx.Tx, txOptions, false)
			if err != nil {
				logrus.WithError(err).Errorf("error adding pending transaction for reorged out tx %v", tx.Tx.Hash())
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
		pool.processBlockTxs(ctx, client, parentBlockInfo.blockWithHash.Block.NumberU64(), parentBlockInfo.blockWithHash, walletMap, parentBlockInfo.txSkipMap, nil)
	}

	return nil
}

// GetCurrentGasLimit returns the current gas limit of the chain.
func (pool *TxPool) GetCurrentGasLimit() uint64 {
	pool.blockStatsMutex.RLock()
	defer pool.blockStatsMutex.RUnlock()
	return pool.currentGasLimit
}

// GetCurrentBaseFee returns the current base fee of the chain.
func (pool *TxPool) GetCurrentBaseFee() *big.Int {
	pool.blockStatsMutex.RLock()
	defer pool.blockStatsMutex.RUnlock()
	return pool.currentBaseFee
}

// GetCurrentGasLimitWithInit returns the current gas limit, initializing it from RPC if needed.
// This is a convenience method that combines GetCurrentGasLimit and InitializeGasLimit.
func (pool *TxPool) GetCurrentGasLimitWithInit() (uint64, error) {
	gasLimit := pool.GetCurrentGasLimit()
	if gasLimit == 0 {
		if err := pool.initBlockStats(); err != nil {
			return 0, err
		}
		gasLimit = pool.GetCurrentGasLimit()
	}
	return gasLimit, nil
}

// GetCurrentBaseFeeWithInit returns the current base fee, initializing it from RPC if needed.
func (pool *TxPool) GetCurrentBaseFeeWithInit() (*big.Int, error) {
	baseFee := pool.GetCurrentBaseFee()
	if baseFee == nil {
		if err := pool.initBlockStats(); err != nil {
			return nil, err
		}
		baseFee = pool.GetCurrentBaseFee()
	}
	return baseFee, nil
}

// initBlockStats fetches the current block stats from the network if not already set.
// This is useful during startup when the pool hasn't processed any blocks yet.
func (pool *TxPool) initBlockStats() error {
	pool.blockStatsMutex.Lock()
	defer pool.blockStatsMutex.Unlock()

	// If we already have a gas limit, don't fetch it again
	if pool.currentGasLimit > 0 && pool.currentBaseFee != nil {
		return nil
	}

	// Try to get a non-builder client to fetch the latest block
	client := pool.options.ClientPool.GetClient(WithoutBuilder())
	if client == nil {
		return fmt.Errorf("no non-builder client available to fetch gas limit")
	}

	// Fetch the latest block to get the gas limit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	latestBlock, err := client.client.BlockByNumber(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to fetch latest block for gas limit: %w", err)
	}

	pool.currentGasLimit = latestBlock.GasLimit()
	pool.currentBaseFee = latestBlock.BaseFee()
	logrus.Infof("initialized block stats from latest block: %v, %v", pool.currentGasLimit, pool.currentBaseFee)

	return nil
}

// GetSuggestedFees returns the suggested fees for a transaction.
// If baseFeeWei and tipFeeWei are provided (non-nil), they are used directly.
// If not provided (nil), the fees are fetched from the client.
func (pool *TxPool) GetSuggestedFees(client *Client, baseFeeWei *big.Int, tipFeeWei *big.Int) (feeCap *big.Int, tipCap *big.Int, err error) {
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

// rebroadcastTransaction performs the actual rebroadcast of a transaction.
// This method encapsulates the existing rebroadcast logic for reuse.
func (pool *TxPool) rebroadcastTransaction(ctx context.Context, tx *types.Transaction, options *SendTransactionOptions, retryCount uint64) {
	submitTx := func(client *Client) error {
		if options.OnEncode != nil {
			txBytes, err := options.OnEncode(tx)
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
	if clientCount > 5 {
		clientCount = 5
	}

	for j := 0; j < clientCount; j++ {
		if ctx.Err() != nil {
			break
		}

		var clientOpts []ClientSelectionOption
		if options.ClientGroup != "" {
			clientOpts = append(clientOpts, WithClientGroup(options.ClientGroup))
		}
		clientOpts = append(clientOpts, WithClientSelectionMode(SelectClientByIndex, j+options.ClientsStartOffset+1))
		client := pool.options.ClientPool.GetClient(clientOpts...)
		if client == nil {
			break
		}

		err := submitTx(client)

		if options.LogFn != nil {
			options.LogFn(client, j, int(retryCount), err)
		}

		if err == nil || strings.Contains(err.Error(), "already known") {
			break
		}
	}
}

// fillNonceGaps creates and submits dummy transactions to fill nonce gaps.
func (pool *TxPool) fillNonceGaps(ctx context.Context, wallet *Wallet, gaps []uint64, baseOptions *SendTransactionOptions) {
	// Get current gas prices from the pool
	baseFee := pool.GetCurrentBaseFee()
	if baseFee == nil {
		baseFee = big.NewInt(1e9) // fallback to 1 gwei
	}

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
		fillerTx, err := wallet.BuildFillerTx(nonce, gasTipCap, gasFeeCap)
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
			OnComplete: func(tx *types.Transaction, receipt *types.Receipt, err error) {
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
