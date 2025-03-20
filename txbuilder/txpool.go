package txbuilder

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"runtime/debug"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/sirupsen/logrus"
)

type TxPool struct {
	options          *TxPoolOptions
	wallets          map[common.Address]*Wallet
	walletsMutex     sync.RWMutex
	processStaleChan chan uint64
	lastBlockNumber  uint64
}

type TxPoolOptions struct {
	Context          context.Context
	GetClientFn      func(index int, random bool) *Client
	GetClientCountFn func() int
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
}

func NewTxPool(options *TxPoolOptions) *TxPool {
	pool := &TxPool{
		options:          options,
		wallets:          map[common.Address]*Wallet{},
		processStaleChan: make(chan uint64, 1),
	}

	if options.Context == nil {
		options.Context = context.Background()
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
		newHighestBlockNumber := pool.getHighestBlockNumber()
		if newHighestBlockNumber > highestBlockNumber {
			if highestBlockNumber > 0 || newHighestBlockNumber < 10 {
				for blockNumber := highestBlockNumber + 1; blockNumber <= newHighestBlockNumber; blockNumber++ {
					err := pool.processBlock(pool.options.Context, blockNumber)
					if err != nil {
						logrus.WithError(err).Errorf("error processing block %v", blockNumber)
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

func (pool *TxPool) processBlock(ctx context.Context, blockNumber uint64) error {
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

	t1 := time.Now()
	blockBody := pool.getBlockBody(ctx, blockNumber)
	if blockBody == nil {
		return fmt.Errorf("could not load block body")
	}

	txCount := len(blockBody.Transactions())
	receipts, err := pool.getBlockReceipts(ctx, blockNumber, txCount)
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

		fromWallet := pool.getWallet(txFrom)
		if fromWallet != nil {
			confirmCount++
			walletMap[txFrom] = true
			pool.processTransactionInclusion(blockNumber, fromWallet, tx, receipt)
		}

		toAddr := tx.To()
		if toAddr != nil {
			toWallet := pool.getWallet(*toAddr)
			if toWallet != nil {
				pool.processTransactionReceival(toWallet, tx)
			}
		}
	}

	logrus.Infof("processed block %v:  %v total tx, %v tx confirmed from %v wallets (%v, %v)", blockNumber, txCount, confirmCount, len(walletMap), loadingTime, time.Since(t1))

	return nil
}

func (pool *TxPool) getHighestBlockNumber() uint64 {
	clientCount := pool.options.GetClientCountFn()
	wg := &sync.WaitGroup{}

	highestBlockNumber := uint64(0)
	highestBlockNumberMutex := sync.Mutex{}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

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
			}
			highestBlockNumberMutex.Unlock()
		}()
	}

	wg.Wait()
	return highestBlockNumber
}

func (pool *TxPool) getBlockBody(ctx context.Context, blockNumber uint64) *types.Block {
	clientCount := pool.options.GetClientCountFn()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	for i := 0; i < clientCount; i++ {
		client := pool.options.GetClientFn(i, false)
		if client == nil {
			continue
		}

		blockBody, err := client.client.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
		if err == nil {
			return blockBody
		}
	}

	return nil
}

func (pool *TxPool) getBlockReceipts(ctx context.Context, blockNumber uint64, txCount int) ([]*types.Receipt, error) {
	clientCount := pool.options.GetClientCountFn()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var receiptErr error
	blockNum := rpc.BlockNumber(blockNumber)

	for i := 0; i < clientCount; i++ {
		client := pool.options.GetClientFn(i, false)
		if client == nil {
			continue
		}

		blockReceipts, err := client.client.BlockReceipts(ctx, rpc.BlockNumberOrHash{
			BlockNumber: &blockNum,
		})
		if err != nil {
			receiptErr = err
		} else {
			if len(blockReceipts) != txCount {
				logrus.Warnf("block %v has %v receipts, expected %v", blockNumber, len(blockReceipts), txCount)
				continue
			}

			return blockReceipts, nil
		}
	}

	return nil, receiptErr
}

func (pool *TxPool) getWallet(address common.Address) *Wallet {
	pool.walletsMutex.RLock()
	defer pool.walletsMutex.RUnlock()
	return pool.wallets[address]
}

func (pool *TxPool) SendTransaction(ctx context.Context, wallet *Wallet, tx *types.Transaction, options *SendTransactionOptions) error {
	var confirmCtx context.Context

	var confirmCancel context.CancelFunc

	if options.OnConfirm != nil || options.MaxRebroadcasts > 0 {
		confirmCtx, confirmCancel = context.WithCancel(ctx)
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
		}()

		wg.Wait()
	}

	var err error

	clientCount := pool.options.GetClientCountFn()
	for i := 0; i < clientCount; i++ {
		client := options.Client
		if client == nil || i > 0 {
			client = pool.options.GetClientFn(i+options.ClientsStartOffset, false)
		}
		if client == nil {
			continue
		}

		err = client.SendTransactionCtx(ctx, tx)

		if options.LogFn != nil {
			options.LogFn(client, i, 0, err)
		}

		if err == nil {
			break
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

					err = client.SendTransactionCtx(ctx, tx)

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
