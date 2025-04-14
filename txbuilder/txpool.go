package txbuilder

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"runtime/debug"
	"slices"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/sirupsen/logrus"
)

// BlockInfo represents information about a processed block
type BlockInfo struct {
	Number     uint64
	Hash       common.Hash
	ParentHash common.Hash
	Timestamp  uint64
}

// TxInfo represents information about a confirmed transaction
type TxInfo struct {
	TxHash     common.Hash
	From       common.Address
	To         *common.Address
	Tx         *types.Transaction
	FromWallet *Wallet
	ToWallet   *Wallet
	Options    *SendTransactionOptions
}

type TxPool struct {
	options          *TxPoolOptions
	wallets          map[common.Address]*Wallet
	walletsMutex     sync.RWMutex
	processStaleChan chan uint64
	lastBlockNumber  uint64

	// Block tracking for reorg detection
	blocksMutex sync.RWMutex
	blocks      map[uint64]*BlockInfo
	reorgDepth  int // Number of blocks to keep in memory for reorg tracking

	// Transaction tracking for reorg recovery
	txsMutex     sync.RWMutex
	confirmedTxs map[uint64][]*TxInfo
}

type TxPoolOptions struct {
	Context          context.Context
	GetClientFn      func(index int, random bool) *Client
	GetClientCountFn func() int
	ReorgDepth       int // Number of blocks to keep in memory for reorg tracking
}

type TxConfirmFn func(tx *types.Transaction, receipt *types.Receipt, err error)
type TxLogFn func(client *Client, retry int, rebroadcast int, err error)

type SendTransactionOptions struct {
	Client             *Client
	ClientsStartOffset int

	OnConfirm TxConfirmFn
	LogFn     TxLogFn

	MaxRebroadcasts     int
	RebroadcastInterval time.Duration
	TransactionBytes    []byte
}

func NewTxPool(options *TxPoolOptions) *TxPool {
	pool := &TxPool{
		options:          options,
		wallets:          map[common.Address]*Wallet{},
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
			pool.walletsMutex.RLock()
			wallets := make([]*Wallet, 0, len(pool.wallets))
			for _, wallet := range pool.wallets {
				wallets = append(wallets, wallet)
			}
			pool.walletsMutex.RUnlock()

			for _, wallet := range wallets {
				pool.processStaleConfirmations(blockNumber, wallet)
			}
		}
	}
}

func (pool *TxPool) processBlock(ctx context.Context, client *Client, blockNumber uint64) error {
	pool.lastBlockNumber = blockNumber

	pool.walletsMutex.RLock()
	walletsLen := len(pool.wallets)
	var chainId *big.Int
	if walletsLen > 0 {
		for _, wallet := range pool.wallets {
			if wallet.chainid == nil {
				continue
			}
			chainId = wallet.chainid
			break
		}
	}
	pool.walletsMutex.RUnlock()

	if walletsLen == 0 {
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
		pool.handleReorg(ctx, client, blockNumber, blockBody, chainId)
	}

	// Store block info
	pool.blocksMutex.Lock()
	pool.blocks[blockNumber] = &BlockInfo{
		Number:     blockNumber,
		Hash:       blockBody.Hash(),
		ParentHash: blockBody.ParentHash(),
		Timestamp:  blockBody.Time(),
	}

	// Clean up old blocks
	if blockNumber > uint64(pool.reorgDepth) {
		delete(pool.blocks, blockNumber-uint64(pool.reorgDepth))
	}
	pool.blocksMutex.Unlock()

	return pool.processBlockTxs(ctx, client, blockNumber, blockBody, chainId)
}

func (pool *TxPool) processBlockTxs(ctx context.Context, client *Client, blockNumber uint64, blockBody *types.Block, chainId *big.Int) error {
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
	walletMap := map[common.Address]bool{}

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
		fromWallet := pool.getWallet(txFrom)
		toAddr := tx.To()
		toWallet := (*Wallet)(nil)
		if toAddr != nil {
			toWallet = pool.getWallet(*toAddr)
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
			walletMap[txFrom] = true
			pool.processTransactionInclusion(blockNumber, fromWallet, tx, receipt)
		}

		if toWallet != nil {
			pool.processTransactionReceival(toWallet, tx)
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

	logrus.Infof("processed block %v:  %v total tx, %v tx confirmed from %v wallets (%v, %v)", blockNumber, txCount, confirmCount, len(walletMap), loadingTime, time.Since(t1))

	return nil
}

func (pool *TxPool) getHighestBlockNumber() (uint64, []*Client) {
	clientCount := pool.options.GetClientCountFn()
	wg := &sync.WaitGroup{}

	highestBlockNumber := uint64(0)
	highestBlockNumberMutex := sync.Mutex{}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	highestBlockNumberClients := []*Client{}

	for i := 0; i < clientCount; i++ {
		client := pool.options.GetClientFn(i, false)
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

func (pool *TxPool) getBlockBody(ctx context.Context, client *Client, blockNumber uint64) *types.Block {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	blockBody, err := client.client.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
	if err == nil {
		return blockBody
	}

	return nil
}

func (pool *TxPool) getBlockReceipts(ctx context.Context, client *Client, blockNumber uint64, txCount int) ([]*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var receiptErr error
	blockNum := rpc.BlockNumber(blockNumber)

	blockReceipts, err := client.client.BlockReceipts(ctx, rpc.BlockNumberOrHash{
		BlockNumber: &blockNum,
	})
	if err != nil {
		receiptErr = err
	} else {
		if len(blockReceipts) != txCount {
			return nil, fmt.Errorf("block %v has %v receipts, expected %v", blockNumber, len(blockReceipts), txCount)
		}

		return blockReceipts, nil
	}

	return nil, receiptErr
}

func (pool *TxPool) getWallet(address common.Address) *Wallet {
	pool.walletsMutex.RLock()
	defer pool.walletsMutex.RUnlock()
	return pool.wallets[address]
}

func (pool *TxPool) SendTransaction(ctx context.Context, wallet *Wallet, tx *types.Transaction, options *SendTransactionOptions) error {
	return pool.addPendingTransaction(ctx, wallet, tx, options, true)
}

func (pool *TxPool) addPendingTransaction(ctx context.Context, wallet *Wallet, tx *types.Transaction, options *SendTransactionOptions, submitNow bool) error {
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
		if options.TransactionBytes != nil {
			return client.SendRawTransactionCtx(ctx, options.TransactionBytes)
		}

		return client.SendTransactionCtx(ctx, tx)
	}

	if submitNow {
		clientCount := pool.options.GetClientCountFn()
		for i := 0; i < clientCount; i++ {
			client := options.Client
			if client == nil || i > 0 {
				client = pool.options.GetClientFn(i+options.ClientsStartOffset, false)
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

		return err
	}

	if options.MaxRebroadcasts > 0 {
		go func() {
			for i := 0; i < options.MaxRebroadcasts; i++ {
				select {
				case <-confirmCtx.Done():
					return
				case <-time.After(options.RebroadcastInterval):
				}

				clientCount := pool.options.GetClientCountFn()
				for j := 0; j < clientCount; j++ {
					client := pool.options.GetClientFn(i+j+options.ClientsStartOffset+1, false)
					if client == nil {
						continue
					}

					err = submitTx(client)

					if options.LogFn != nil {
						options.LogFn(client, j, i+1, err)
					}

					if err == nil {
						break
					}
				}
			}
		}()
	}

	return nil
}

func (pool *TxPool) AwaitTransaction(ctx context.Context, wallet *Wallet, tx *types.Transaction) (*types.Receipt, error) {
	return pool.awaitTransaction(ctx, wallet, tx, nil)
}

func (pool *TxPool) awaitTransaction(ctx context.Context, wallet *Wallet, tx *types.Transaction, wg *sync.WaitGroup) (*types.Receipt, error) {
	pool.walletsMutex.Lock()
	pool.wallets[wallet.address] = wallet
	pool.walletsMutex.Unlock()

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

func (pool *TxPool) processTransactionInclusion(blockNumber uint64, wallet *Wallet, tx *types.Transaction, receipt *types.Receipt) {
	nonce := tx.Nonce()

	wallet.txNonceMutex.Lock()
	defer wallet.txNonceMutex.Unlock()

	if nonceChan := wallet.txNonceChans[nonce]; nonceChan != nil {
		nonceChan.receipt = receipt
	}

	wallet.confirmedNonce = nonce + 1
	if nonce+1 > wallet.pendingNonce.Load() {
		wallet.pendingNonce.Store(nonce + 1)
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

func (pool *TxPool) processStaleConfirmations(blockNumber uint64, wallet *Wallet) {
	if len(wallet.txNonceChans) > 0 && blockNumber > wallet.lastConfirmation+10 {
		wallet.lastConfirmation = blockNumber

		var lastNonce uint64
		var err error
		for retry := 0; retry < 3; retry++ {
			lastNonce, err = pool.options.GetClientFn(retry, true).GetNonceAt(pool.options.Context, wallet.address, big.NewInt(int64(blockNumber)))
			if err == nil {
				break
			}
		}

		pendingNonce := 0
		for n := range wallet.txNonceChans {
			pendingNonce = int(n)
			break
		}

		logrus.Debugf("recovering %v stale transactions for %v (current nonce %v, cache nonce %v, first pending nonce: %v)", len(wallet.txNonceChans), wallet.address.String(), lastNonce, wallet.confirmedNonce, pendingNonce)

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

func (pool *TxPool) processTransactionReceival(wallet *Wallet, tx *types.Transaction) {
	wallet.balance = wallet.balance.Add(wallet.balance, tx.Value())
}

func (pool *TxPool) loadTransactionReceipt(ctx context.Context, tx *types.Transaction) *types.Receipt {
	retryCount := uint64(0)

	for {
		client := pool.options.GetClientFn(int(retryCount), true)

		reqCtx, reqCtxCancel := context.WithTimeout(ctx, 5*time.Second)

		//nolint:gocritic // ignore
		defer reqCtxCancel()

		receipt, err := client.GetTransactionReceiptCtx(reqCtx, tx.Hash())
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

// handleReorg handles a chain reorganization
func (pool *TxPool) handleReorg(ctx context.Context, client *Client, blockNumber uint64, newBlock *types.Block, chainId *big.Int) error {
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
				tx.FromWallet.confirmedNonce = tx.Tx.Nonce()
				tx.FromWallet.txNonceMutex.Unlock()
			}

			// add tx as pending tx
			txOptions := &SendTransactionOptions{
				Client: client,
				OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
					if err != nil {
						logrus.WithError(err).Errorf("error confirming reorged out tx %v", tx.Hash())
					} else {
						logrus.Infof("reorged out tx %v confirmed", tx.Hash())
					}
				},
				RebroadcastInterval: 30 * time.Second,
				MaxRebroadcasts:     10,
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
			// reverse processTransactionReceival
			tx.ToWallet.balance = tx.ToWallet.balance.Sub(tx.ToWallet.balance, tx.Tx.Value())
		}
	}

	// re-process the new parent blocks
	slices.Reverse(newBlockParents)
	for _, parentBlock := range newBlockParents {
		pool.processBlockTxs(ctx, client, parentBlock.NumberU64(), parentBlock, chainId)
	}

	return nil
}
