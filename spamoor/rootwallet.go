package spamoor

import (
	"context"
	"sync"

	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"

	"github.com/ethpandaops/spamoor/utils"
)

type RootWallet struct {
	wallet      *Wallet
	txbatcher   *TxBatcher
	txSemMutex  sync.Mutex
	txSemaphore chan struct{}
}

func InitRootWallet(ctx context.Context, privkey string, client *Client, logger logrus.FieldLogger) (*RootWallet, error) {
	rootWallet, err := NewWallet(privkey)
	if err != nil {
		return nil, err
	}

	err = client.UpdateWallet(ctx, rootWallet)
	if err != nil {
		return nil, err
	}

	if logger != nil {
		logger.Infof(
			"initialized root wallet (addr: %v balance: %v ETH, nonce: %v)",
			rootWallet.GetAddress().String(),
			utils.WeiToEther(uint256.MustFromBig(rootWallet.GetBalance())).Uint64(),
			rootWallet.GetNonce(),
		)
	}

	return &RootWallet{
		wallet:      rootWallet,
		txSemaphore: make(chan struct{}, 200),
	}, nil
}

func (wallet *RootWallet) GetWallet() *Wallet {
	return wallet.wallet
}

func (wallet *RootWallet) WithWalletLock(ctx context.Context, txCount int, lockedLogFn func(), lockedFn func() error) error {
	acquiredCount := 0
	acquireLock := func() error {
		wallet.txSemMutex.Lock()
		defer wallet.txSemMutex.Unlock()

		for i := 0; i < txCount; i++ {
			if lockedLogFn != nil {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case wallet.txSemaphore <- struct{}{}:
					acquiredCount++
					continue
				default:
					lockedLogFn()
					lockedLogFn = nil
				}
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case wallet.txSemaphore <- struct{}{}:
				acquiredCount++
				continue
			}
		}

		return nil
	}

	defer func() {
		for i := 0; i < acquiredCount; i++ {
			<-wallet.txSemaphore
		}
	}()

	err := acquireLock()
	if err != nil {
		return err
	}

	return lockedFn()
}

func (wallet *RootWallet) GetTxBatcher() *TxBatcher {
	return wallet.txbatcher
}

func (wallet *RootWallet) InitTxBatcher(ctx context.Context, txpool *TxPool) {
	wallet.txbatcher = NewTxBatcher(txpool)
}
