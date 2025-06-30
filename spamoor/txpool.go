package spamoor

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"math/big"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
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
type BulkBlockStatsCallback func(blockNumber uint64, allWalletPoolStats map[*WalletPool]*WalletPoolBlockStats)

// WalletPoolBlockStats contains transaction statistics for a wallet pool in a specific block
type WalletPoolBlockStats struct {
	ConfirmedTxCount uint64
	TotalTxFees      *big.Int
	AffectedWallets  int
	Block            *types.Block
	ConfirmedTxs     []*TxInfo
	Receipts         []*types.Receipt
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

// RebroadcastRequest represents a transaction that needs to be rebroadcast
type RebroadcastRequest struct {
	confirmCtx  context.Context
	fromWallet  *Wallet
	tx          *types.Transaction
	options     *SendTransactionOptions
	retryCount  uint64
	nextAttempt time.Time
}

// TxPool manages transaction submission, confirmation tracking, and chain reorganization handling.
// It monitors blockchain blocks, tracks transaction confirmations, handles reorgs by re-submitting
// affected transactions, and provides transaction awaiting functionality with automatic rebroadcasting.
type TxPool struct {
	options          *TxPoolOptions
	processStaleChan chan uint64
	lastBlockNumber  uint64

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

	// Rebroadcast management
	rebroadcastRequests     []*RebroadcastRequest
	rebroadcastRequestsChan chan *RebroadcastRequest
	rebroadcastMutex        sync.RWMutex
}

// TxPoolOptions contains configuration options for the transaction pool.
type TxPoolOptions struct {
	Context              context.Context
	Logger               *logrus.Entry
	ClientPool           *ClientPool
	ReorgDepth           int // Number of blocks to keep in memory for reorg tracking
	GetActiveWalletPools func() []*WalletPool
}

// NewTxPool creates a new transaction pool with the specified options.
// It starts background goroutines for block processing and stale transaction handling.
// The pool automatically begins monitoring the blockchain for new blocks and managing
// transaction confirmations and reorgs.
func NewTxPool(options *TxPoolOptions) *TxPool {
	pool := &TxPool{
		options:           options,
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

	// Initialize rebroadcast system
	pool.rebroadcastRequests = []*RebroadcastRequest{}
	pool.rebroadcastRequestsChan = make(chan *RebroadcastRequest, 100)

	go pool.runTxPoolLoop()
	go pool.processStaleTransactionsLoop()
	go pool.runRebroadcastLoop()

	return pool
}

// runTxPoolLoop continuously monitors for new blocks and processes them.
// It tracks the highest block number across all clients and processes new blocks
// sequentially. Also triggers stale transaction processing when new blocks arrive.
// Runs until the context is cancelled and recovers from panics with logging.
func (pool *TxPool) runTxPoolLoop() {
	defer func() {
		utils.RecoverPanic(pool.options.Logger, "TxPool.runTxPoolLoop", pool.runTxPoolLoop)
	}()

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
					err := pool.processBlock(pool.options.Context, client, blockNumber)
					if err != nil {
						logrus.WithError(err).Errorf("error processing block %v", blockNumber)
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

// getWalletMap collects all wallets from active wallet pools into a single map.
// It iterates through all active wallet pools and calls their collectPoolWallets
// method to build a comprehensive address-to-wallet mapping.
func (pool *TxPool) getWalletMap() map[common.Address]*Wallet {
	walletMap := map[common.Address]*Wallet{}
	walletPools := pool.options.GetActiveWalletPools()
	for _, walletPool := range walletPools {
		walletPool.collectPoolWallets(walletMap)
	}
	return walletMap
}

// processBlock processes a single block for transaction confirmations and reorg detection.
// It loads the block body, checks for chain reorganizations by comparing parent hashes,
// stores block information for reorg tracking, and processes all transactions in the block.
// Also handles cleanup of old block data based on the reorg depth setting.
func (pool *TxPool) processBlock(ctx context.Context, client *Client, blockNumber uint64) error {
	pool.lastBlockNumber = blockNumber

	walletPools := pool.options.GetActiveWalletPools()
	walletMap := map[common.Address]*Wallet{}

	var chainId *big.Int
	for _, walletPool := range walletPools {
		if walletPool.GetChainId() == nil {
			continue
		}
		chainId = walletPool.GetChainId()
		walletPool.collectPoolWallets(walletMap)
	}

	if len(walletMap) == 0 {
		return nil
	}

	blockBody := pool.getBlockBody(ctx, client, blockNumber)
	if blockBody == nil {
		return fmt.Errorf("could not load block body")
	}

	// Check for reorg by comparing parent hash
	pool.blocksMutex.RLock()
	lastBlock, hasLastBlock := pool.blocks[blockNumber-1]
	pool.blocksMutex.RUnlock()

	if hasLastBlock && lastBlock.Hash != blockBody.ParentHash() {
		logrus.Warnf("Detected chain reorganization at block %d. Parent hash mismatch: expected %s, got %s",
			blockNumber, lastBlock.Hash.Hex(), blockBody.ParentHash().Hex())

		// Handle reorg
		pool.handleReorg(ctx, client, blockNumber, blockBody, chainId, walletMap)
	}

	// Store block info
	pool.blocksMutex.Lock()
	pool.blocks[blockNumber] = &BlockInfo{
		Number:     blockNumber,
		Hash:       blockBody.Hash(),
		ParentHash: blockBody.ParentHash(),
		Timestamp:  blockBody.Time(),
		GasLimit:   blockBody.GasLimit(),
	}

	// Update current gas limit
	pool.blockStatsMutex.Lock()
	pool.currentGasLimit = blockBody.GasLimit()
	pool.currentBaseFee = blockBody.BaseFee()
	pool.blockStatsMutex.Unlock()

	// Clean up old blocks
	if blockNumber > uint64(pool.reorgDepth) {
		delete(pool.blocks, blockNumber-uint64(pool.reorgDepth))
	}
	pool.blocksMutex.Unlock()

	return pool.processBlockTxs(ctx, client, blockNumber, blockBody, chainId, walletMap)
}

// processBlockTxs processes all transactions in a block for confirmation tracking.
// It loads block receipts, decodes transaction senders, updates wallet states for
// confirmed transactions, and tracks transaction information for reorg recovery.
// Also handles cleanup of old confirmed transaction data.
func (pool *TxPool) processBlockTxs(ctx context.Context, client *Client, blockNumber uint64, blockBody *types.Block, chainId *big.Int, walletMap map[common.Address]*Wallet) error {
	t1 := time.Now()
	txCount := len(blockBody.Transactions())
	receipts, err := pool.getBlockReceipts(ctx, client, blockBody.Hash(), txCount)
	if err != nil {
		return fmt.Errorf("could not load block receipts: %w", err)
	}

	loadingTime := time.Since(t1)
	t1 = time.Now()

	signer := types.LatestSignerForChainID(chainId)
	confirmCount := 0
	affectedWalletMap := map[common.Address]bool{}
	pool.txsMutex.Lock()
	pool.confirmedTxs[blockNumber] = []*TxInfo{}
	pool.txsMutex.Unlock()

	for idx, tx := range blockBody.Transactions() {
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
	pool.txsMutex.RLock()
	blockConfirmedTxs := pool.confirmedTxs[blockNumber]
	pool.txsMutex.RUnlock()

	pool.notifyBlockSubscribers(blockNumber, blockConfirmedTxs, blockBody, receipts)

	return nil
}

// getHighestBlockNumber queries all good clients to find the highest block number.
// It runs concurrent queries to all available clients and returns the highest
// block number found along with the clients that reported that height.
func (pool *TxPool) getHighestBlockNumber() (uint64, []*Client) {
	clientCount := len(pool.options.ClientPool.GetAllGoodClients())
	wg := &sync.WaitGroup{}

	highestBlockNumber := uint64(0)
	highestBlockNumberMutex := sync.Mutex{}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	highestBlockNumberClients := []*Client{}

	for i := 0; i < clientCount; i++ {
		client := pool.options.ClientPool.GetClient(SelectClientByIndex, i, "")
		if client == nil {
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
func (pool *TxPool) getBlockBody(ctx context.Context, client *Client, blockNumber uint64) *types.Block {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	blockBody, err := client.client.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
	if err == nil {
		return blockBody
	}

	return nil
}

// getBlockReceipts retrieves all transaction receipts for a block.
// It validates that the number of receipts matches the expected transaction count
// and uses a 5-second timeout for the request.
func (pool *TxPool) getBlockReceipts(ctx context.Context, client *Client, blockHash common.Hash, txCount int) ([]*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var receiptErr error

	blockReceipts, err := client.client.BlockReceipts(ctx, rpc.BlockNumberOrHash{
		BlockHash: &blockHash,
	})
	if err != nil {
		receiptErr = err
	} else {
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
	if pool.options.GetActiveWalletPools != nil {
		return pool.options.GetActiveWalletPools()
	}
	return []*WalletPool{}
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
func (pool *TxPool) notifyBlockSubscribers(blockNumber uint64, confirmedTxs []*TxInfo, block *types.Block, receipts []*types.Receipt) {
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

	// If we have bulk subscribers, calculate stats for all wallet pools once
	var allWalletPoolStats map[*WalletPool]*WalletPoolBlockStats
	if len(bulkSubscriptions) > 0 {
		allWalletPoolStats = pool.calculateAllWalletPoolStats(confirmedTxs, block, receipts)

		// Notify bulk subscribers
		for _, bulkSubscription := range bulkSubscriptions {
			bulkSubscription.Callback(blockNumber, allWalletPoolStats)
		}
	}

	// Notify individual subscribers (reuse stats if already calculated)
	for _, subscription := range subscriptions {
		var stats *WalletPoolBlockStats
		if allWalletPoolStats != nil {
			// Reuse already calculated stats
			stats = allWalletPoolStats[subscription.WalletPool]
		} else {
			// Calculate stats individually
			stats = pool.calculateWalletPoolStats(subscription.WalletPool, confirmedTxs, block, receipts)
		}
		subscription.Callback(blockNumber, stats)
	}
}

// calculateWalletPoolStats calculates transaction statistics for a specific wallet pool.
func (pool *TxPool) calculateWalletPoolStats(walletPool *WalletPool, confirmedTxs []*TxInfo, block *types.Block, receipts []*types.Receipt) *WalletPoolBlockStats {
	stats := &WalletPoolBlockStats{
		TotalTxFees:  big.NewInt(0),
		Block:        block,
		Receipts:     make([]*types.Receipt, 0, len(receipts)),
		ConfirmedTxs: make([]*TxInfo, 0, len(confirmedTxs)),
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

	for idx, txInfo := range confirmedTxs {
		if txInfo.FromWallet != nil && walletSet[txInfo.FromWallet.GetAddress()] {
			stats.ConfirmedTxCount++
			if txInfo.TxFees != nil {
				totalFee := new(big.Int).Add(&txInfo.TxFees.FeeAmount, &txInfo.TxFees.BlobFeeAmount)
				stats.TotalTxFees.Add(stats.TotalTxFees, totalFee)
			}
			affectedWallets[txInfo.FromWallet.GetAddress()] = true

			stats.Receipts = append(stats.Receipts, receipts[idx])
			stats.ConfirmedTxs = append(stats.ConfirmedTxs, confirmedTxs[idx])
		}
	}

	stats.AffectedWallets = len(affectedWallets)

	return stats
}

// calculateAllWalletPoolStats efficiently calculates statistics for ALL active wallet pools in a single pass.
// This avoids recalculating the same data multiple times for bulk subscribers.
func (pool *TxPool) calculateAllWalletPoolStats(confirmedTxs []*TxInfo, block *types.Block, receipts []*types.Receipt) map[*WalletPool]*WalletPoolBlockStats {
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
			Receipts:     make([]*types.Receipt, 0, len(receipts)),
			ConfirmedTxs: make([]*TxInfo, 0, len(confirmedTxs)),
		}
	}

	// Track affected wallets per pool
	affectedWallets := make(map[*WalletPool]map[common.Address]bool)
	for _, walletPool := range activeWalletPools {
		affectedWallets[walletPool] = make(map[common.Address]bool)
	}

	// Single pass through confirmed transactions
	for idx, txInfo := range confirmedTxs {
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
				stats.ConfirmedTxs = append(stats.ConfirmedTxs, confirmedTxs[idx])
			}
		}
	}

	// Set affected wallet counts
	for walletPool, stats := range allStats {
		stats.AffectedWallets = len(affectedWallets[walletPool])
	}

	return allStats
}

// submitTransaction handles the core transaction submission logic.
// It starts a confirmation tracking goroutine, submits the transaction to clients,
// and optionally sets up automatic rebroadcasting. The submitNow parameter controls
// whether to immediately submit or just set up confirmation tracking.
func (pool *TxPool) submitTransaction(ctx context.Context, wallet *Wallet, tx *types.Transaction, options *SendTransactionOptions, submitNow bool) error {
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
			} else if options.OnConfirm != nil && receipt != nil {
				options.OnConfirm(tx, receipt)
			}

			if options.OnComplete != nil {
				options.OnComplete(tx, receipt, err)
			}
		}()

		// Track transaction result for metrics
		defer func() {
			walletPools := pool.options.GetActiveWalletPools()
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

		receipt, err = pool.awaitTransaction(confirmCtx, wallet, tx, wg)
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

	var err error

	submitTx := func(client *Client) error {
		if options.OnEncode != nil {
			txBytes, err := options.OnEncode(tx)
			if err != nil {
				return fmt.Errorf("failed to encode transaction: %w", err)
			}

			if txBytes != nil {
				return client.SendRawTransaction(ctx, txBytes)
			}
		}

		return client.SendTransaction(ctx, tx)
	}

	if submitNow {
		clientCount := len(pool.options.ClientPool.GetAllGoodClients())
		for i := 0; i < clientCount; i++ {
			client := options.Client
			if client == nil || i > 0 {
				client = pool.options.ClientPool.GetClient(SelectClientByIndex, i+options.ClientsStartOffset, options.ClientGroup)
			}
			if client == nil {
				continue
			}

			err = submitTx(client)

			if options.LogFn != nil {
				options.LogFn(client, i, 0, err)
			}

			if err == nil {
				break
			}
		}
	}

	submissionComplete <- err

	if err != nil {
		confirmCancel()

		// Track initial transaction submission failure for metrics
		walletPools := pool.options.GetActiveWalletPools()
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

		return err
	}

	// Increment wallet's submitted transaction counter
	wallet.IncrementSubmittedTxCount()

	// Start reliable rebroadcast if enabled
	if options.Rebroadcast {
		pool.scheduleRebroadcast(confirmCtx, wallet, tx, options)
	}

	return nil
}

// awaitTransaction waits for a specific transaction to be confirmed.
// It uses the wallet's nonce channel system to wait for confirmation and
// handles cases where the transaction might be replaced or reorged.
// The wg parameter is signaled when confirmation tracking is set up.
func (pool *TxPool) awaitTransaction(ctx context.Context, wallet *Wallet, tx *types.Transaction, wg *sync.WaitGroup) (*types.Receipt, error) {
	txHash := tx.Hash()
	nonceChan, isFirstPendingTx := wallet.getTxNonceChan(tx.Nonce())

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

// processStaleConfirmations recovers stale transactions that may have been missed.
// It checks if a wallet has pending transactions that are older than 10 blocks,
// queries the current nonce from the blockchain, and recovers any confirmed
// transactions that weren't properly tracked.
func (pool *TxPool) processStaleConfirmations(blockNumber uint64, wallet *Wallet) {
	if len(wallet.txNonceChans) > 0 && blockNumber > wallet.lastConfirmation+10 {
		wallet.lastConfirmation = blockNumber

		var lastNonce uint64
		var err error
		for retry := 0; retry < 3; retry++ {
			client := pool.options.ClientPool.GetClient(SelectClientRandom, retry, "")
			if client == nil {
				continue
			}

			lastNonce, err = client.GetNonceAt(pool.options.Context, wallet.address, big.NewInt(int64(blockNumber)))
			if err == nil {
				break
			}
		}

		if lastNonce == wallet.confirmedTxCount {
			return
		}

		pendingNonce := 0
		for n := range wallet.txNonceChans {
			pendingNonce = int(n)
			break
		}

		logrus.Debugf("recovering stale transactions for %v (tx count: %v, current nonce %v, cache nonce %v, first pending nonce: %v)", wallet.address.String(), len(wallet.txNonceChans), lastNonce, wallet.confirmedTxCount, pendingNonce)

		wallet.txNonceMutex.Lock()
		defer wallet.txNonceMutex.Unlock()

		for n := range wallet.txNonceChans {
			if n < lastNonce {
				logrus.WithError(err).Warnf("recovering stale confirmed transactions for %v (nonce %v)", wallet.address.String(), n)
				close(wallet.txNonceChans[n].channel)
				delete(wallet.txNonceChans, n)
			}
		}
	}
}

// loadTransactionReceipt attempts to load a transaction receipt from multiple clients.
// It retries up to 5 times with different clients and includes exponential backoff.
// Returns nil if the receipt cannot be loaded after all retries.
func (pool *TxPool) loadTransactionReceipt(ctx context.Context, tx *types.Transaction) *types.Receipt {
	retryCount := uint64(0)

	for {
		client := pool.options.ClientPool.GetClient(SelectClientRandom, int(retryCount), "")
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
func (pool *TxPool) handleReorg(ctx context.Context, client *Client, blockNumber uint64, newBlock *types.Block, chainId *big.Int, walletMap map[common.Address]*Wallet) error {
	newBlockParents := []*types.Block{}

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
	block := newBlock
	for {
		if block.NumberU64() == 0 {
			break
		}

		blockNumber := block.NumberU64() - 1
		if pool.blocks[blockNumber] == nil {
			break
		}

		if pool.blocks[blockNumber].Hash == block.ParentHash() {
			break
		}

		parentBlockBody := pool.getBlockBody(ctx, client, blockNumber)
		if parentBlockBody == nil {
			return fmt.Errorf("could not load block body for new parent block %v", blockNumber)
		}

		newBlockParents = append(newBlockParents, parentBlockBody)
		block = parentBlockBody
	}

	reorgBaseBlock := block

	// find all the transactions that were reorged out
	pool.txsMutex.Lock()
	reorgedOutTxs := []*TxInfo{}
	for blockNum := reorgBaseBlock.NumberU64(); blockNum <= blockNumber; blockNum++ {
		blockTxs := pool.confirmedTxs[blockNum]
		if blockTxs == nil {
			continue
		}

		reorgedOutTxs = append(reorgedOutTxs, blockTxs...)
	}

	// remove reorged out blocks & txs from cache
	pool.blocksMutex.Lock()
	for blockNum := reorgBaseBlock.NumberU64() + 1; blockNum <= blockNumber; blockNum++ {
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
	for _, parentBlock := range newBlockParents {
		pool.processBlockTxs(ctx, client, parentBlock.NumberU64(), parentBlock, chainId, walletMap)
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

	// Try to get a client to fetch the latest block
	client := pool.options.ClientPool.GetClient(SelectClientRandom, 0, "")
	if client == nil {
		return fmt.Errorf("no client available to fetch gas limit")
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
// If baseFeeGwei and tipFeeGwei are provided, they are used as the base fee and tip fee.
// If not provided, the fees are fetched from the client. The fees are returned in wei.
func (pool *TxPool) GetSuggestedFees(client *Client, baseFeeGwei float64, tipFeeGwei float64) (feeCap *big.Int, tipCap *big.Int, err error) {
	if baseFeeGwei > 0 {
		feeCap = new(big.Int).SetUint64(uint64(baseFeeGwei * 1e9))
	}
	if tipFeeGwei > 0 {
		tipCap = new(big.Int).SetUint64(uint64(tipFeeGwei * 1e9))
	}

	if feeCap == nil || tipCap == nil {
		feeCap, tipCap, err = client.GetSuggestedFee(pool.options.Context)
		if err != nil {
			return nil, nil, err
		}
	}

	return feeCap, tipCap, nil
}

// calculateBackoffDelay calculates the exponential backoff delay for rebroadcast attempts.
// Uses 30s base delay, 1.5x multiplier, with 10min maximum delay.
func (pool *TxPool) calculateBackoffDelay(retryCount uint64) time.Duration {
	const (
		baseDelay  = 30 * time.Second
		multiplier = 1.5
		maxDelay   = 10 * time.Minute
	)

	delay := time.Duration(float64(baseDelay) * math.Pow(multiplier, float64(retryCount)))
	if delay > maxDelay {
		delay = maxDelay
	}
	return delay
}

// isTransactionBlocking checks if a transaction is blocking wallet progress.
// Returns true if the transaction nonce is the next required nonce for the wallet.
func (pool *TxPool) isTransactionBlocking(wallet *Wallet, txNonce uint64) bool {
	wallet.txNonceMutex.Lock()
	defer wallet.txNonceMutex.Unlock()

	// This transaction is the next one that needs to be included
	return txNonce <= wallet.confirmedTxCount
}

// rebroadcastTransaction performs the actual rebroadcast of a transaction.
// This method encapsulates the existing rebroadcast logic for reuse.
func (pool *TxPool) rebroadcastTransaction(ctx context.Context, tx *types.Transaction, options *SendTransactionOptions, retryCount uint64) {
	submitTx := func(client *Client) error {
		var err error
		if options.OnEncode != nil {
			txBytes, encodeErr := options.OnEncode(tx)
			if encodeErr != nil {
				return fmt.Errorf("failed to encode transaction: %w", encodeErr)
			}
			err = client.SendRawTransaction(ctx, txBytes)
		} else {
			err = client.SendTransaction(ctx, tx)
		}
		return err
	}

	clientCount := len(pool.options.ClientPool.GetAllGoodClients())
	for j := 0; j < clientCount; j++ {
		client := pool.options.ClientPool.GetClient(SelectClientByIndex, j+options.ClientsStartOffset+1, options.ClientGroup)
		if client == nil {
			continue
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

// runRebroadcastLoop runs the main rebroadcast processing loop
func (pool *TxPool) runRebroadcastLoop() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-pool.options.Context.Done():
			return
		case newRequest := <-pool.rebroadcastRequestsChan:
			pool.addRebroadcastRequest(newRequest)
		case <-ticker.C:
			nextAttempt := pool.processRebroadcastRequests()
			if !nextAttempt.IsZero() {
				duration := time.Until(nextAttempt)
				if duration > 0 {
					ticker.Reset(duration)
				} else {
					ticker.Reset(50 * time.Millisecond)
				}
			} else {
				ticker.Reset(10 * time.Second)
			}
		}
	}
}

// addRebroadcastRequest adds a new rebroadcast request to the queue
func (pool *TxPool) addRebroadcastRequest(req *RebroadcastRequest) {
	pool.rebroadcastMutex.Lock()
	defer pool.rebroadcastMutex.Unlock()

	// Set initial next attempt time
	backoffDelay := pool.calculateBackoffDelay(req.retryCount)
	req.nextAttempt = time.Now().Add(backoffDelay)

	pool.rebroadcastRequests = append(pool.rebroadcastRequests, req)
}

// processRebroadcastRequests processes all pending rebroadcast requests
func (pool *TxPool) processRebroadcastRequests() time.Time {
	pool.rebroadcastMutex.Lock()
	defer pool.rebroadcastMutex.Unlock()

	now := time.Now()
	activeRequests := make([]*RebroadcastRequest, 0, len(pool.rebroadcastRequests))
	lowestNextAttempt := time.Time{}

	for _, req := range pool.rebroadcastRequests {
		select {
		case <-req.confirmCtx.Done():
			// Transaction confirmed, remove from queue
			continue
		default:
		}

		// Check if it's time to rebroadcast
		if now.After(req.nextAttempt) {
			// Check if this transaction is blocking wallet progress
			if pool.isTransactionBlocking(req.fromWallet, req.tx.Nonce()) {
				// Perform rebroadcast
				go pool.rebroadcastTransaction(req.confirmCtx, req.tx, req.options, req.retryCount)
				req.retryCount++

				// Calculate next attempt time
				backoffDelay := pool.calculateBackoffDelay(req.retryCount)
				req.nextAttempt = now.Add(backoffDelay)
			}

			if lowestNextAttempt.IsZero() || req.nextAttempt.Before(lowestNextAttempt) {
				lowestNextAttempt = req.nextAttempt
			}
		}

		// Keep request in queue for next iteration
		activeRequests = append(activeRequests, req)
	}

	pool.rebroadcastRequests = activeRequests
	return lowestNextAttempt
}

// scheduleRebroadcast schedules a transaction for rebroadcasting
func (pool *TxPool) scheduleRebroadcast(confirmCtx context.Context, fromWallet *Wallet, tx *types.Transaction, options *SendTransactionOptions) {
	req := &RebroadcastRequest{
		confirmCtx: confirmCtx,
		fromWallet: fromWallet,
		tx:         tx,
		options:    options,
		retryCount: 0,
	}

	select {
	case pool.rebroadcastRequestsChan <- req:
	case <-pool.options.Context.Done():
	}
}
