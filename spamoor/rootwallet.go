package spamoor

import (
	"context"
	"sync"

	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"

	"github.com/ethpandaops/spamoor/utils"
)

// RootWallet represents a primary wallet with transaction rate limiting and batching capabilities.
// It wraps a standard Wallet with a semaphore-based transaction limiter and optional transaction batcher
// for managing high-volume transaction scenarios.
type RootWallet struct {
	wallet      *Wallet
	txbatcher   *TxBatcher
	txSemMutex  sync.Mutex
	txSemaphore chan struct{}
	txSemLimit  int
}

// InitRootWallet creates and initializes a new RootWallet from a private key.
// It creates the underlying wallet, updates its state from the blockchain,
// and sets up transaction rate limiting with a default limit of 200 concurrent transactions.
// Returns the initialized RootWallet and logs wallet information if logger is provided.
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
		txSemLimit:  200,
	}, nil
}

// GetWallet returns the underlying Wallet instance.
func (wallet *RootWallet) GetWallet() *Wallet {
	return wallet.wallet
}

// WithWalletLock executes a function while holding transaction semaphore locks.
// It acquires the specified number of transaction slots from the semaphore,
// calls the optional lockedLogFn when waiting for locks, then executes lockedFn.
// The locks are automatically released when the function returns.
//
// Parameters:
//   - ctx: context for cancellation
//   - txCount: number of transaction slots to acquire
//   - lockedLogFn: optional function called once when waiting for locks (can be nil)
//   - lockedFn: function to execute while holding the locks
func (wallet *RootWallet) WithWalletLock(ctx context.Context, txCount int, lockedLogFn func(), lockedFn func() error) error {
	acquiredCount := 0
	acquireLock := func() error {
		wallet.txSemMutex.Lock()
		defer wallet.txSemMutex.Unlock()

		for i := 0; i < txCount; i++ {

			if acquiredCount >= wallet.txSemLimit {
				return nil
			}

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

// GetTxBatcher returns the transaction batcher instance, or nil if not initialized.
func (wallet *RootWallet) GetTxBatcher() *TxBatcher {
	return wallet.txbatcher
}

// InitTxBatcher initializes the transaction batcher with the specified transaction pool.
// This enables batched transaction processing for improved efficiency.
func (wallet *RootWallet) InitTxBatcher(ctx context.Context, txpool *TxPool) {
	wallet.txbatcher = NewTxBatcher(txpool)
}
