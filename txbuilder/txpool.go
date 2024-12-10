package txbuilder

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
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
}

type TxPoolOptions struct {
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

	go pool.runTxPoolLoop()
	go pool.processStaleTransactionsLoop()

	return pool
}

func (pool *TxPool) runTxPoolLoop() {
	highestBlockNumber := uint64(0)
	for {
		newHighestBlockNumber := pool.getHighestBlockNumber()
		if newHighestBlockNumber > highestBlockNumber {
			if highestBlockNumber > 0 || newHighestBlockNumber < 10 {
				for blockNumber := highestBlockNumber + 1; blockNumber <= newHighestBlockNumber; blockNumber++ {
					err := pool.processBlock(blockNumber)
					if err != nil {
						logrus.WithError(err).Errorf("error processing block %v", blockNumber)
					}
				}

			}
			highestBlockNumber = newHighestBlockNumber

			select {
			case pool.processStaleChan <- highestBlockNumber:
			default:
			}
		}
		time.Sleep(3 * time.Second)
	}
}

func (pool *TxPool) processStaleTransactionsLoop() {
	for blockNumber := range pool.processStaleChan {
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

func (pool *TxPool) processBlock(blockNumber uint64) error {
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

	blockBody := pool.getBlockBody(blockNumber)
	if blockBody == nil {
		return fmt.Errorf("could not load block body")
	}

	receipts := pool.getBlockReceipts(blockNumber)
	if receipts == nil {
		return fmt.Errorf("could not load block receipts")
	}

	signer := types.LatestSignerForChainID(chainId)

	for idx, tx := range blockBody.Transactions() {
		receipt := receipts[idx]
		if receipt == nil {
			continue
		}

		txFrom, err := types.Sender(signer, tx)
		if err != nil {
			logrus.Warnf("error decoding tx sender (block %v, tx %v): %v", blockNumber, idx, err)
			continue
		}

		fromWallet := pool.getWallet(txFrom)
		if fromWallet != nil {
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

func (pool *TxPool) getBlockBody(blockNumber uint64) *types.Block {
	clientCount := pool.options.GetClientCountFn()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

func (pool *TxPool) getBlockReceipts(blockNumber uint64) []*types.Receipt {
	clientCount := pool.options.GetClientCountFn()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i := 0; i < clientCount; i++ {
		client := pool.options.GetClientFn(i, false)
		if client == nil {
			continue
		}

		blockNum := rpc.BlockNumber(int64(blockNumber))
		blockReceipts, err := client.client.BlockReceipts(ctx, rpc.BlockNumberOrHash{
			BlockNumber: &blockNum,
		})
		if err == nil {
			return blockReceipts
		}
	}

	return nil
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
	nonceChan := wallet.getTxNonceChan(tx.Nonce())

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

	if wallet.confirmedNonce > nonce {
		return
	}

	for n := range wallet.txNonceChans {
		if n == nonce {
			wallet.txNonceChans[n].receipt = receipt
		}

		if n <= nonce {
			close(wallet.txNonceChans[n].channel)
			delete(wallet.txNonceChans, n)
		}
	}

	wallet.confirmedNonce = nonce + 1
	if wallet.confirmedNonce > wallet.pendingNonce {
		wallet.pendingNonce = wallet.confirmedNonce
	}
	if blockNumber > wallet.lastConfirmation {
		wallet.lastConfirmation = blockNumber
	}
}

func (pool *TxPool) processStaleConfirmations(blockNumber uint64, wallet *Wallet) {
	if len(wallet.txNonceChans) > 0 && blockNumber > wallet.lastConfirmation+10 {
		wallet.lastConfirmation = blockNumber

		var lastNonce uint64
		var err error
		for retry := 0; retry < 3; retry++ {
			lastNonce, err = pool.options.GetClientFn(retry, true).GetNonceAt(wallet.address, big.NewInt(int64(blockNumber)))
			if err == nil {
				break
			}
		}

		wallet.txNonceMutex.Lock()
		defer wallet.txNonceMutex.Unlock()

		if wallet.confirmedNonce >= lastNonce {
			return
		}

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
