package spamoor

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

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

// ValidateFundingOptions contains configuration for root wallet funding validation.
type ValidateFundingOptions struct {
	MinBalance      *big.Int      // Minimum required balance in Wei
	RetryInterval   time.Duration // Time to wait between balance checks
	MaxRetries      int           // Maximum number of retry attempts (0 = unlimited)
	TimeoutDuration time.Duration // Maximum time to wait for funding (0 = unlimited)
}

// GetDefaultValidateFundingOptions returns default options for root wallet funding validation.
func GetDefaultValidateFundingOptions() *ValidateFundingOptions {
	minBalance := new(big.Int)
	minBalance.SetString("10000000000000000000", 10)

	return &ValidateFundingOptions{
		MinBalance:      minBalance,
		RetryInterval:   30 * time.Second,
		MaxRetries:      0,
		TimeoutDuration: 0,
	}
}

// ValidateFunding validates that the root wallet has sufficient funds before starting scenarios.
// It checks the wallet balance against the minimum required balance and waits for funding if needed.
// The function respects context cancellation and returns an error if funding validation fails.
func (wallet *RootWallet) ValidateFunding(ctx context.Context, client *Client, options *ValidateFundingOptions, logger logrus.FieldLogger) error {
	if options == nil {
		options = GetDefaultValidateFundingOptions()
	}

	if options.TimeoutDuration > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, options.TimeoutDuration)
		defer cancel()

		if logger != nil {
			logger.Infof("root wallet funding validation will timeout after %v", options.TimeoutDuration)
		}
	}

	retryCount := 0
	for {
		if ctx.Err() != nil {
			return fmt.Errorf("context cancelled while waiting for root wallet funding: %w", ctx.Err())
		}

		err := client.UpdateWallet(ctx, wallet.wallet)
		if err != nil {
			if logger != nil {
				logger.WithError(err).Warnf("failed to update root wallet balance, retrying in %v", options.RetryInterval)
			}
			select {
			case <-ctx.Done():
				return fmt.Errorf("context cancelled while updating root wallet: %w", ctx.Err())
			case <-time.After(options.RetryInterval):
				retryCount++
				if options.MaxRetries > 0 && retryCount >= options.MaxRetries {
					return fmt.Errorf("failed to update root wallet balance after %d retries: %w", retryCount, err)
				}
				continue
			}
		}

		balance := wallet.wallet.GetBalance()
		balanceETH := utils.WeiToEther(uint256.MustFromBig(balance))

		if balance.Cmp(options.MinBalance) >= 0 {
			if logger != nil {
				logger.Infof("root wallet funding validation successful (balance: %v ETH, required: %v ETH)",
					balanceETH,
					utils.WeiToEther(uint256.MustFromBig(options.MinBalance)))
			}
			return nil
		}

		requiredETH := utils.WeiToEther(uint256.MustFromBig(options.MinBalance))
		if logger != nil {
			if retryCount == 0 {
				logger.Warnf("root wallet has insufficient funds (balance: %v ETH, required: %v ETH)", balanceETH, requiredETH)
				logger.Infof("waiting for root wallet funding (address: %s)...", wallet.wallet.GetAddress().String())
			} else {
				logger.Debugf("root wallet still underfunded (balance: %v ETH, required: %v ETH), continuing to wait...", balanceETH, requiredETH)
			}
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled while waiting for root wallet funding: %w", ctx.Err())
		case <-time.After(options.RetryInterval):
			retryCount++
			if options.MaxRetries > 0 && retryCount >= options.MaxRetries {
				return fmt.Errorf("root wallet funding validation failed after %d retries (balance: %v ETH, required: %v ETH)",
					retryCount, balanceETH, requiredETH)
			}
		}
	}
}
