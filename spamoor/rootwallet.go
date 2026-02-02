package spamoor

import (
	"context"
	"fmt"
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
	wallet              *Wallet
	txbatcher           *TxBatcher
	txSemMutex          sync.Mutex
	txSemaphore         chan struct{}
	txSemLimit          int
	balanceUpdateCtx    context.Context
	balanceUpdateCancel context.CancelFunc
	logger              logrus.FieldLogger

	// Track pending funding amounts to avoid overshooting available balance
	pendingFundingMutex sync.Mutex
	pendingFundingTotal *uint256.Int
}

// InitRootWallet creates and initializes a new RootWallet from a private key.
// It creates the underlying wallet, updates its state from the blockchain,
// and sets up transaction rate limiting with a default limit of 200 concurrent transactions.
// Returns the initialized RootWallet and logs wallet information if logger is provided.
func InitRootWallet(ctx context.Context, privkey string, clientPool *ClientPool, txpool *TxPool, logger logrus.FieldLogger) (*RootWallet, error) {
	privateKey, address, err := LoadPrivateKey(privkey)
	if err != nil {
		return nil, err
	}
	rootWallet := NewWallet(privateKey, address)

	txpool.RegisterWallet(rootWallet, ctx)

	client := clientPool.GetClient()
	if client == nil {
		return nil, fmt.Errorf("no client available")
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

	balanceUpdateCtx, balanceUpdateCancel := context.WithCancel(ctx)

	rootWalletInstance := &RootWallet{
		wallet:              rootWallet,
		txSemaphore:         make(chan struct{}, 200),
		txSemLimit:          200,
		balanceUpdateCtx:    balanceUpdateCtx,
		balanceUpdateCancel: balanceUpdateCancel,
		logger:              logger,
		pendingFundingTotal: uint256.NewInt(0),
	}

	go rootWalletInstance.balanceUpdateLoop(clientPool)

	return rootWalletInstance, nil
}

// GetWallet returns the underlying Wallet instance.
func (wallet *RootWallet) GetWallet() *Wallet {
	return wallet.wallet
}

// WithWalletLock executes a function while holding transaction semaphore locks and ensuring sufficient balance.
// It first waits for sufficient balance (including pending amounts), then acquires the specified number
// of transaction slots from the semaphore, reserves the funding amount, and executes the locked function.
// The locks and funding reservation are automatically released when the function returns.
//
// Parameters:
//   - ctx: context for cancellation
//   - txCount: number of transaction slots to acquire
//   - fundingAmount: total amount to reserve from wallet balance (nil to skip balance check)
//   - clientPool: client pool for balance updates during waiting
//   - lockedLogFn: optional function called once when waiting for locks (can be nil)
//   - lockedFn: function to execute while holding the locks
func (wallet *RootWallet) WithWalletLock(ctx context.Context, txCount int, fundingAmount *uint256.Int, clientPool *ClientPool, lockedLogFn func(reason string), lockedFn func() error) error {
	acquiredCount := 0

	acquireLock := func() error {
		wallet.txSemMutex.Lock()
		defer wallet.txSemMutex.Unlock()

		// Await & reserve funding amount before acquiring transaction locks
		if fundingAmount != nil {
			err := wallet.waitForSufficientBalance(ctx, fundingAmount, txCount, clientPool, lockedLogFn)
			if err != nil {
				return err
			}

			wallet.pendingFundingMutex.Lock()
			wallet.pendingFundingTotal = wallet.pendingFundingTotal.Add(wallet.pendingFundingTotal, fundingAmount)
			wallet.pendingFundingMutex.Unlock()
		}

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
					lockedLogFn("waiting for other funding txs to finish")
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
		// Release transaction semaphore locks
		for i := 0; i < acquiredCount; i++ {
			<-wallet.txSemaphore
		}

		// Release funding reservation
		if fundingAmount != nil {
			wallet.pendingFundingMutex.Lock()
			wallet.pendingFundingTotal = wallet.pendingFundingTotal.Sub(wallet.pendingFundingTotal, fundingAmount)
			wallet.pendingFundingMutex.Unlock()
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

// balanceUpdateLoop runs in the background and updates the wallet balance every 4 blocks (~48 seconds).
// Uses a separate client connection to avoid interfering with transaction operations.
func (wallet *RootWallet) balanceUpdateLoop(clientPool *ClientPool) {
	// Update every 4 blocks (assuming 12 second block times = ~48 seconds)
	updateInterval := 48 * time.Second
	ticker := time.NewTicker(updateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-wallet.balanceUpdateCtx.Done():
			return
		case <-ticker.C:
			client := clientPool.GetClient(WithClientSelectionMode(SelectClientRandom))
			if client != nil {
				err := client.UpdateWallet(wallet.balanceUpdateCtx, wallet.wallet)
				if err != nil && wallet.logger != nil {
					wallet.logger.Debugf("failed to update root wallet balance: %v", err)
				}
			}
		}
	}
}

// Shutdown gracefully shuts down the root wallet and stops the balance updater.
func (wallet *RootWallet) Shutdown() {
	if wallet.balanceUpdateCancel != nil {
		wallet.balanceUpdateCancel()
	}
}

// hasSufficientBalance checks if the root wallet has enough balance for the specified amount plus gas costs and reserve.
// Returns true if sufficient funds are available, false otherwise.
// Includes a 1 ETH reserve buffer for batcher contract deployment and other operations.
// Also accounts for pending funding amounts that are already allocated but not yet reflected in the balance.
func (wallet *RootWallet) hasSufficientBalance(requiredAmount *uint256.Int, txCount int) (bool, *uint256.Int) {
	// Add gas buffer: estimate 100 Gwei gas price * 100k gas per tx as safety margin
	gasBuffer := uint256.NewInt(100000000000)                    // 100 Gwei
	gasBuffer = gasBuffer.Mul(gasBuffer, uint256.NewInt(100000)) // 100k gas
	gasBuffer = gasBuffer.Mul(gasBuffer, uint256.NewInt(uint64(txCount)))

	// Add 1 ETH reserve buffer for batcher contract deployment and other operations
	reserveBuffer := utils.EtherToWei(uint256.NewInt(1)) // 1 ETH

	totalRequired := uint256.NewInt(0)
	totalRequired = totalRequired.Add(requiredAmount, gasBuffer)
	totalRequired = totalRequired.Add(totalRequired, reserveBuffer)

	// Get current balance and subtract pending funding amounts
	currentBalance := uint256.MustFromBig(wallet.wallet.GetBalance())

	wallet.pendingFundingMutex.Lock()
	availableBalance := uint256.NewInt(0)
	availableBalance = availableBalance.Sub(currentBalance, wallet.pendingFundingTotal)
	wallet.pendingFundingMutex.Unlock()

	return availableBalance.Cmp(totalRequired) >= 0, totalRequired
}

// waitForSufficientBalance blocks until the root wallet has sufficient balance for the specified requirements.
// Returns error only if context is cancelled, otherwise blocks until funds are available.
func (wallet *RootWallet) waitForSufficientBalance(ctx context.Context, requiredAmount *uint256.Int, txCount int, clientPool *ClientPool, lockedLogFn func(reason string)) error {
	sufficient, totalRequired := wallet.hasSufficientBalance(requiredAmount, txCount)
	if sufficient {
		return nil
	}

	if lockedLogFn != nil {
		lockedLogFn(fmt.Sprintf("insufficient root wallet balance. awaiting addtional funds. needed: %v ETH", utils.WeiToEther(totalRequired).Uint64()))
	}

	// Wait for balance to be sufficient
	ticker := time.NewTicker(12 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled while waiting for root wallet funding")
		case <-ticker.C:
			if clientPool != nil {
				// Update balance manually here to avoid waiting for the background updater
				client := clientPool.GetClient(WithClientSelectionMode(SelectClientRandom))
				if client == nil {
					if wallet.logger != nil {
						wallet.logger.Debugf("no client available while waiting for root wallet funding")
					}
					continue // no client available, skip
				}

				err := client.UpdateWallet(ctx, wallet.wallet)
				if err != nil && wallet.logger != nil {
					wallet.logger.Debugf("failed to update root wallet balance while waiting for funding: %v", err)
					continue
				}
			}

			currentBalance := uint256.MustFromBig(wallet.wallet.GetBalance())

			wallet.pendingFundingMutex.Lock()
			availableBalance := uint256.NewInt(0)
			availableBalance = availableBalance.Sub(currentBalance, wallet.pendingFundingTotal)
			wallet.pendingFundingMutex.Unlock()

			if availableBalance.Cmp(totalRequired) >= 0 {
				return nil
			}

			// Log current status
			if wallet.logger != nil {
				deficit := uint256.NewInt(0)
				deficit = deficit.Sub(totalRequired, availableBalance)
				wallet.logger.Infof(
					"waiting for root wallet funding. Current: %v ETH, still need: %v ETH",
					utils.WeiToEther(availableBalance).Uint64(),
					utils.WeiToEther(deficit).Uint64(),
				)
			}
		}
	}
}
