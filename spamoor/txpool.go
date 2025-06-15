package spamoor

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"math/big"
	"runtime/debug"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"

	"github.com/ethpandaops/spamoor/spamoortypes"
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
	FromWallet spamoortypes.Wallet
	ToWallet   spamoortypes.Wallet
	Options    *spamoortypes.SendTransactionOptions
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
	gasLimitMutex   sync.RWMutex
}

// TxPoolOptions contains configuration options for the transaction pool.
type TxPoolOptions struct {
	Context              context.Context
	ClientPool           spamoortypes.ClientPool
	ReorgDepth           int // Number of blocks to keep in memory for reorg tracking
	GetActiveWalletPools func() []spamoortypes.WalletPool
}

// NewTxPool creates a new transaction pool with the specified options.
// It starts background goroutines for block processing and stale transaction handling.
// The pool automatically begins monitoring the blockchain for new blocks and managing
// transaction confirmations and reorgs.
func NewTxPool(options *TxPoolOptions) *TxPool {
	pool := &TxPool{
		options:          options,
		processStaleChan: make(chan uint64, 1),
		blocks:           map[uint64]*BlockInfo{},
		confirmedTxs:     map[uint64][]*TxInfo{},
		reorgDepth:       10, // Default value
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

// runTxPoolLoop continuously monitors for new blocks and processes them.
// It tracks the highest block number across all clients and processes new blocks
// sequentially. Also triggers stale transaction processing when new blocks arrive.
// Runs until the context is cancelled and recovers from panics with logging.
func (pool *TxPool) runTxPoolLoop() {
	defer func() {
		if err := recover(); err != nil {
			logrus.WithError(err.(error)).Errorf("uncaught panic in TxPool.runTxPoolLoop subroutine: %v, stack: %v", err, string(debug.Stack()))
		}
	}()

	highestBlockNumber := uint64(0)
	for {
		newHighestBlockNumber, clients := pool.getHighestBlockNumber()
		if newHighestBlockNumber > highestBlockNumber {
			if highestBlockNumber > 0 || newHighestBlockNumber < 10 {
				blockNumber := highestBlockNumber + 1
				for _, client := range clients {
					for ; blockNumber <= newHighestBlockNumber; blockNumber++ {
						err := pool.processBlock(pool.options.Context, client, blockNumber)
						if err != nil {
							logrus.WithError(err).Errorf("error processing block %v", blockNumber)
						}
					}
					if blockNumber == newHighestBlockNumber {
						break
					}
				}
			}
			highestBlockNumber = newHighestBlockNumber

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
		if err := recover(); err != nil {
			logrus.WithError(err.(error)).Errorf("uncaught panic in TxPool.processStaleTransactionsLoop subroutine: %v, stack: %v", err, string(debug.Stack()))
		}
	}()

	for {
		select {
		case <-pool.options.Context.Done():
			return
		case blockNumber := <-pool.processStaleChan:
			for _, wallet := range pool.getWalletMap() {
				if wallet.GetPendingTxCount() > 0 && blockNumber > wallet.GetLastConfirmation()+10 {
					wallet.SetLastConfirmation(blockNumber)

					var lastNonce uint64
					var err error
					for retry := 0; retry < 3; retry++ {
						client := pool.options.ClientPool.GetClient(spamoortypes.SelectClientRandom, retry, "")
						if client == nil {
							continue
						}

						lastNonce, err = client.GetNonceAt(pool.options.Context, wallet.GetAddress(), big.NewInt(int64(blockNumber)))
						if err == nil {
							break
						}
					}

					wallet.ProcessStaleTransactions(blockNumber, lastNonce)
				}
			}
		}
	}
}

// getWalletMap collects all wallets from active wallet pools into a single map.
// It iterates through all active wallet pools and calls their collectPoolWallets
// method to build a comprehensive address-to-wallet mapping.
func (pool *TxPool) getWalletMap() map[common.Address]spamoortypes.Wallet {
	walletMap := map[common.Address]spamoortypes.Wallet{}
	walletPools := pool.options.GetActiveWalletPools()
	for _, walletPool := range walletPools {
		walletPool.CollectPoolWallets(walletMap)
	}
	return walletMap
}

// processBlock processes a single block for transaction confirmations and reorg detection.
// It loads the block body, checks for chain reorganizations by comparing parent hashes,
// stores block information for reorg tracking, and processes all transactions in the block.
// Also handles cleanup of old block data based on the reorg depth setting.
func (pool *TxPool) processBlock(ctx context.Context, client spamoortypes.Client, blockNumber uint64) error {
	pool.lastBlockNumber = blockNumber

	walletPools := pool.options.GetActiveWalletPools()
	walletMap := map[common.Address]spamoortypes.Wallet{}

	var chainId *big.Int
	for _, walletPool := range walletPools {
		if walletPool.GetChainId() == nil {
			continue
		}
		chainId = walletPool.GetChainId()
		walletPool.CollectPoolWallets(walletMap)
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
	pool.gasLimitMutex.Lock()
	pool.currentGasLimit = blockBody.GasLimit()
	pool.gasLimitMutex.Unlock()

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
func (pool *TxPool) processBlockTxs(ctx context.Context, client spamoortypes.Client, blockNumber uint64, blockBody *types.Block, chainId *big.Int, walletMap map[common.Address]spamoortypes.Wallet) error {
	t1 := time.Now()
	txCount := len(blockBody.Transactions())
	receipts, err := pool.getBlockReceipts(ctx, client, blockNumber, txCount)
	if receipts == nil {
		return fmt.Errorf("could not load block receipts: %w", err)
	}

	loadingTime := time.Since(t1)
	t1 = time.Now()

	signer := types.LatestSignerForChainID(chainId)
	confirmCount := 0
	affectedWalletMap := map[common.Address]bool{}

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
		fromWallet := walletMap[txFrom]
		toAddr := tx.To()
		toWallet := (spamoortypes.Wallet)(nil)
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
				FromWallet: fromWallet,
				ToWallet:   toWallet,
			})
			pool.txsMutex.Unlock()
		}

		if fromWallet != nil {
			confirmCount++
			affectedWalletMap[txFrom] = true
			fromWallet.ProcessTransactionInclusion(blockNumber, tx, receipt)
		}

		if toWallet != nil {
			toWallet.ProcessTransactionReceival(tx)
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

	return nil
}

// getHighestBlockNumber queries all good clients to find the highest block number.
// It runs concurrent queries to all available clients and returns the highest
// block number found along with the clients that reported that height.
func (pool *TxPool) getHighestBlockNumber() (uint64, []spamoortypes.Client) {
	clientCount := len(pool.options.ClientPool.GetAllGoodClients())
	wg := &sync.WaitGroup{}

	highestBlockNumber := uint64(0)
	highestBlockNumberMutex := sync.Mutex{}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	highestBlockNumberClients := []spamoortypes.Client{}

	for i := 0; i < clientCount; i++ {
		client := pool.options.ClientPool.GetClient(spamoortypes.SelectClientByIndex, i, "")
		if client == nil {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			blockNumber, err := client.GetBlockHeight(ctx)
			if err != nil {
				return
			}

			highestBlockNumberMutex.Lock()
			if blockNumber > highestBlockNumber {
				highestBlockNumber = blockNumber
				highestBlockNumberClients = []spamoortypes.Client{client}
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
func (pool *TxPool) getBlockBody(ctx context.Context, client spamoortypes.Client, blockNumber uint64) *types.Block {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	blockBody, err := client.GetBlock(ctx, blockNumber)
	if err == nil {
		return blockBody
	}

	return nil
}

// getBlockReceipts retrieves all transaction receipts for a block.
// It validates that the number of receipts matches the expected transaction count
// and uses a 5-second timeout for the request.
func (pool *TxPool) getBlockReceipts(ctx context.Context, client spamoortypes.Client, blockNumber uint64, txCount int) ([]*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	blockReceipts, receiptErr := client.GetBlockReceipts(ctx, blockNumber)
	if receiptErr == nil {
		if len(blockReceipts) != txCount {
			return nil, fmt.Errorf("block %v has %v receipts, expected %v", blockNumber, len(blockReceipts), txCount)
		}

		return blockReceipts, nil
	}

	return nil, receiptErr
}

// SendTransaction submits a transaction to the network with the specified options.
// It handles client selection, rebroadcasting, confirmation tracking, and error handling.
// The transaction is automatically rebroadcast according to the options until confirmed.
func (pool *TxPool) SendTransaction(ctx context.Context, wallet spamoortypes.Wallet, tx *types.Transaction, options *spamoortypes.SendTransactionOptions) error {
	return pool.addPendingTransaction(ctx, wallet, tx, options, true)
}

// addPendingTransaction handles the core transaction submission logic.
// It starts a confirmation tracking goroutine, submits the transaction to clients,
// and optionally sets up automatic rebroadcasting. The submitNow parameter controls
// whether to immediately submit or just set up confirmation tracking.
func (pool *TxPool) addPendingTransaction(ctx context.Context, wallet spamoortypes.Wallet, tx *types.Transaction, options *spamoortypes.SendTransactionOptions, submitNow bool) error {
	confirmCtx, confirmCancel := context.WithCancel(ctx)
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		var receipt *types.Receipt
		var err error

		defer confirmCancel()

		if options.OnConfirm != nil {
			defer func() {
				options.OnConfirm(tx, receipt, err)
			}()
		}

		// Track transaction result for metrics
		defer func() {
			walletPools := pool.options.GetActiveWalletPools()
			for _, walletPool := range walletPools {
				if tracker := walletPool.GetTransactionTracker(); tracker != nil {
					// Check if this wallet belongs to this pool
					allWallets := walletPool.GetAllWallets()
					for _, poolWallet := range allWallets {
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

	submitTx := func(client spamoortypes.Client) error {
		if options.TransactionBytes != nil {
			return client.SendRawTransaction(ctx, options.TransactionBytes)
		}

		return client.SendTransaction(ctx, tx)
	}

	if submitNow {
		clientCount := len(pool.options.ClientPool.GetAllGoodClients())
		for i := 0; i < clientCount; i++ {
			client := options.Client
			if client == nil || i > 0 {
				client = pool.options.ClientPool.GetClient(spamoortypes.SelectClientByIndex, i+options.ClientsStartOffset, options.ClientGroup)
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

	if err != nil {
		if confirmCancel != nil {
			confirmCancel()
		}

		// Call OnConfirm callback immediately for submission failures
		if options.OnConfirm != nil {
			options.OnConfirm(tx, nil, err)
		}

		// Track initial transaction submission failure for metrics
		walletPools := pool.options.GetActiveWalletPools()
		for _, walletPool := range walletPools {
			if tracker := walletPool.GetTransactionTracker(); tracker != nil {
				// Check if this wallet belongs to this pool
				allWallets := walletPool.GetAllWallets()
				for _, poolWallet := range allWallets {
					if poolWallet.GetAddress() == wallet.GetAddress() {
						tracker(err)
						break
					}
				}
			}
		}

		return err
	}

	// Start reliable rebroadcast if enabled
	if options.Rebroadcast {
		pool.startReliableRebroadcast(ctx, confirmCtx, wallet, tx, options)
	}

	return nil
}

// AwaitTransaction waits for a transaction to be confirmed and returns its receipt.
// It monitors the blockchain for the transaction and handles reorgs by continuing
// to wait if the transaction gets reorged out of the chain.
func (pool *TxPool) AwaitTransaction(ctx context.Context, wallet spamoortypes.Wallet, tx *types.Transaction) (*types.Receipt, error) {
	return pool.awaitTransaction(ctx, wallet, tx, nil)
}

// awaitTransaction waits for a specific transaction to be confirmed.
// It uses the wallet's nonce channel system to wait for confirmation and
// handles cases where the transaction might be replaced or reorged.
// The wg parameter is signaled when confirmation tracking is set up.
func (pool *TxPool) awaitTransaction(ctx context.Context, wallet spamoortypes.Wallet, tx *types.Transaction, wg *sync.WaitGroup) (*types.Receipt, error) {
	txHash := tx.Hash()
	nonceChan, isFirstPendingTx := wallet.GetTxNonceChan(tx.Nonce())

	if isFirstPendingTx && pool.lastBlockNumber > wallet.GetLastConfirmation()+1 {
		wallet.SetLastConfirmation(pool.lastBlockNumber - 1)
	}

	if wg != nil {
		wg.Done()
	}

	if nonceChan != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-nonceChan.Channel:
		}

		receipt := nonceChan.Receipt
		if receipt != nil {
			if bytes.Equal(receipt.TxHash[:], txHash[:]) {
				return receipt, nil
			}

			return nil, nil
		}
	}

	return pool.loadTransactionReceipt(ctx, tx), nil
}

// SendAndAwaitTxRange sends multiple transactions and waits for all of them to be confirmed.
// It automatically handles fee calculation, balance updates, and provides confirmation
// callbacks for each transaction. All transactions are processed concurrently.
func (pool *TxPool) SendAndAwaitTxRange(ctx context.Context, wallet spamoortypes.Wallet, txs []*types.Transaction, options *spamoortypes.SendTransactionOptions) error {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	for idx := range txs {
		err := func(idx int) error {
			tx := txs[idx]

			return pool.SendTransaction(ctx, wallet, tx, &spamoortypes.SendTransactionOptions{
				Client:      options.Client,
				Rebroadcast: true,
				OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
					defer wg.Done()

					if err != nil {
						if options.OnConfirm != nil {
							options.OnConfirm(tx, receipt, err)
						}
						return
					}

					feeAmount := big.NewInt(0)
					if receipt != nil {
						effectiveGasPrice := receipt.EffectiveGasPrice
						if effectiveGasPrice == nil {
							effectiveGasPrice = big.NewInt(0)
						}
						feeAmount = feeAmount.Mul(effectiveGasPrice, big.NewInt(int64(receipt.GasUsed)))
					}

					totalAmount := big.NewInt(0).Add(tx.Value(), feeAmount)
					wallet.SubBalance(totalAmount)

					if options.OnConfirm != nil {
						options.OnConfirm(tx, receipt, nil)
					}
				},
			})
		}(idx)

		if err != nil {
			return err
		}
		wg.Add(1)
	}

	wg.Done()
	wg.Wait()
	return nil
}

// loadTransactionReceipt attempts to load a transaction receipt from multiple clients.
// It retries up to 5 times with different clients and includes exponential backoff.
// Returns nil if the receipt cannot be loaded after all retries.
func (pool *TxPool) loadTransactionReceipt(ctx context.Context, tx *types.Transaction) *types.Receipt {
	retryCount := uint64(0)

	for {
		client := pool.options.ClientPool.GetClient(spamoortypes.SelectClientRandom, int(retryCount), "")
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
func (pool *TxPool) handleReorg(ctx context.Context, client spamoortypes.Client, blockNumber uint64, newBlock *types.Block, chainId *big.Int, walletMap map[common.Address]spamoortypes.Wallet) error {
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

				tx.FromWallet.SetConfirmedNonce(tx.Tx.Nonce())
			}

			// add tx as pending tx
			txOptions := &spamoortypes.SendTransactionOptions{
				Client:      client,
				Rebroadcast: true,
				OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
					if err != nil {
						logrus.WithError(err).Errorf("error confirming reorged out tx %v", tx.Hash())
					} else {
						logrus.Infof("reorged out tx %v confirmed", tx.Hash())
					}
				},
			}

			if tx.Options != nil {
				txOptions.LogFn = tx.Options.LogFn
			}

			err := pool.addPendingTransaction(ctx, tx.FromWallet, tx.Tx, txOptions, false)
			if err != nil {
				logrus.WithError(err).Errorf("error adding pending transaction for reorged out tx %v", tx.Tx.Hash())
			}
		}

		if tx.ToWallet != nil {
			tx.ToWallet.RevertTransactionReceival(tx.Tx)
		}
	}

	// re-process the new parent blocks
	slices.Reverse(newBlockParents)
	for _, parentBlock := range newBlockParents {
		pool.processBlockTxs(ctx, client, parentBlock.NumberU64(), parentBlock, chainId, walletMap)
	}

	return nil
}

// GetCurrentGasLimit returns the current gas limit of the transaction pool.
func (pool *TxPool) GetCurrentGasLimit() uint64 {
	pool.gasLimitMutex.RLock()
	defer pool.gasLimitMutex.RUnlock()
	return pool.currentGasLimit
}

// GetCurrentGasLimitWithInit returns the current gas limit, initializing it from RPC if needed.
// This is a convenience method that combines GetCurrentGasLimit and InitializeGasLimit.
func (pool *TxPool) GetCurrentGasLimitWithInit() (uint64, error) {
	gasLimit := pool.GetCurrentGasLimit()
	if gasLimit == 0 {
		if err := pool.InitializeGasLimit(); err != nil {
			return 0, err
		}
		gasLimit = pool.GetCurrentGasLimit()
	}
	return gasLimit, nil
}

// InitializeGasLimit fetches the current block gas limit from the network if not already set.
// This is useful during startup when the pool hasn't processed any blocks yet.
func (pool *TxPool) InitializeGasLimit() error {
	pool.gasLimitMutex.Lock()
	defer pool.gasLimitMutex.Unlock()

	// If we already have a gas limit, don't fetch it again
	if pool.currentGasLimit > 0 {
		return nil
	}

	// Try to get a client to fetch the latest block
	client := pool.options.ClientPool.GetClient(spamoortypes.SelectClientRandom, 0, "")
	if client == nil {
		return fmt.Errorf("no client available to fetch gas limit")
	}

	// Fetch the latest block to get the gas limit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	gasLimit, err := client.GetLatestGasLimit(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch latest block for gas limit: %w", err)
	}

	pool.currentGasLimit = gasLimit
	logrus.Infof("initialized gas limit from latest block: %v", pool.currentGasLimit)

	return nil
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

// startReliableRebroadcast starts a reliable rebroadcast goroutine for a transaction.
// It uses exponential backoff and unlimited retries, only rebroadcasting when the
// transaction is blocking wallet progress.
func (pool *TxPool) startReliableRebroadcast(ctx context.Context, confirmCtx context.Context, fromWallet spamoortypes.Wallet, tx *types.Transaction, options *spamoortypes.SendTransactionOptions) {
	retryCount := uint64(0)

	go func() {
		for {
			// Calculate backoff delay
			backoffDelay := pool.calculateBackoffDelay(retryCount)

			select {
			case <-confirmCtx.Done():
				return // Transaction confirmed
			case <-time.After(backoffDelay):
				// Check if this transaction is blocking wallet progress
				if tx.Nonce() != fromWallet.GetNextNonce() {
					// Transaction is not the next required nonce, stop rebroadcasting
					continue
				}

				// Perform rebroadcast using existing logic
				pool.rebroadcastTransaction(ctx, tx, options, retryCount)
				retryCount++
			}
		}
	}()
}

// rebroadcastTransaction performs the actual rebroadcast of a transaction.
// This method encapsulates the existing rebroadcast logic for reuse.
func (pool *TxPool) rebroadcastTransaction(ctx context.Context, tx *types.Transaction, options *spamoortypes.SendTransactionOptions, retryCount uint64) {
	submitTx := func(client spamoortypes.Client) error {
		var err error
		if options.TransactionBytes != nil {
			err = client.SendRawTransaction(ctx, options.TransactionBytes)
		} else {
			err = client.SendTransaction(ctx, tx)
		}
		return err
	}

	clientCount := len(pool.options.ClientPool.GetAllGoodClients())
	for j := 0; j < clientCount; j++ {
		client := pool.options.ClientPool.GetClient(spamoortypes.SelectClientByIndex, j+options.ClientsStartOffset+1, options.ClientGroup)
		if client == nil {
			continue
		}

		if options.OnRebroadcast != nil {
			options.OnRebroadcast(tx, options, client)
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
