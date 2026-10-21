package scenario

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

// GlobalSlotDuration is the global setting for slot/block duration used in rate limiting.
// This can be set via CLI flag (--slot-duration) and applies to all scenarios.
// For L1 chains, this is typically 12s. For L2s like Arbitrum, it can be 250ms or less.
// It is the fallback for chains that do not report EIP-7843 slot numbers.
var GlobalSlotDuration time.Duration = 12 * time.Second

// SlotDuration returns the chain's slot duration derived from EIP-7843 block slot
// numbers, falling back to GlobalSlotDuration.
func SlotDuration(walletPool *spamoor.WalletPool) time.Duration {
	if walletPool != nil {
		if txPool := walletPool.GetTxPool(); txPool != nil {
			if duration := txPool.GetSlotDuration(); duration > 0 {
				return duration
			}
		}
	}
	return GlobalSlotDuration
}

// TransactionScenarioOptions configures how the transaction scenario is executed.
type TransactionScenarioOptions struct {
	TotalCount                  uint64
	Throughput                  uint64
	MaxPending                  uint64
	ThroughputIncrementInterval uint64
	Timeout                     time.Duration // Maximum duration for scenario execution (0 = no timeout)
	WalletPool                  *spamoor.WalletPool
	NoAwaitTransactions         bool // If true, the scenario will not wait for transactions to be included in a block

	// Logger for scenario execution information
	Logger *logrus.Entry

	// ProcessNextTxFn handles transaction execution with the given index
	// It should return:
	// - A callback function to log transaction results (can be nil)
	// - An error if transaction creation failed
	ProcessNextTxFn func(ctx context.Context, params *ProcessNextTxParams) error
}

type ProcessNextTxParams struct {
	TxIdx           uint64
	OrderedLogCb    func(logFunc func()) // Optional callback to log transaction results in order
	NotifySubmitted func()               // Optional callback to notify that the transaction has been submitted
}

// RunTransactionScenario executes a controlled transaction scenario with rate limiting
// and concurrency management. It processes transactions according to the specified options
// until either context cancellation or reaching TotalCount (if > 0).
//
// Features:
// - Rate-limited execution based on Throughput
// - Optional dynamic throughput increases
// - Concurrency control via MaxPending
// - Sequential logging via transaction chaining
//
// Returns an error only if the scenario cannot be started. Transaction failures
// should be handled within ProcessNextTxFn.
func RunTransactionScenario(ctx context.Context, options TransactionScenarioOptions) error {
	slotDuration := SlotDuration(options.WalletPool)

	txIdxCounter := uint64(0)
	pendingCount := atomic.Int64{}
	txCount := atomic.Uint64{}
	throughputTracker := newThroughputTracker()

	// Apply timeout if specified
	if options.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, options.Timeout)
		defer cancel()

		options.Logger.Infof("scenario will timeout after %v", options.Timeout)
	}

	// Logging synchronization structures
	type txState struct {
		logCb     []func()
		submitted bool
		done      bool
	}

	txStates := make(map[uint64]*txState)
	txStatesMutex := sync.Mutex{}
	nextLogIdx := uint64(0)

	// Helper function to process pending logs in order.
	// The caller must hold txStatesMutex.
	processPendingLogs := func() {
		for {
			state, exists := txStates[nextLogIdx]
			if !exists || !state.submitted {
				break
			}

			// Execute the log callback
			if len(state.logCb) > 0 {
				for _, logcb := range state.logCb {
					logcb()
				}
				state.logCb = nil
			}

			// The transaction has been logged in order, so its state is no longer
			// needed here. The worker goroutine keeps its own reference for the
			// rest of its teardown, so removing it from the map does not affect
			// in-flight work and keeps the map from growing with every transaction.
			delete(txStates, nextLogIdx)

			nextLogIdx++
		}
	}

	var maxPending atomic.Uint64
	var pendingMutex sync.Mutex
	var pendingCond *sync.Cond

	pendingWg := sync.WaitGroup{}

	if options.MaxPending > 0 {
		maxPending.Store(options.MaxPending)
		pendingCond = sync.NewCond(&pendingMutex)
	}

	throughput := float64(options.Throughput)
	var throughputMutex sync.Mutex

	slotRate := func() rate.Limit {
		if throughput == 0 {
			return rate.Inf
		}
		return rate.Limit(throughput / slotDuration.Seconds())
	}
	limiter := rate.NewLimiter(slotRate(), 1)

	isErrorMode := atomic.Bool{}
	errorLimiter := rate.NewLimiter(rate.Limit(0.5), 1) // 2 sec interval when in error mode

	// Subscribe to block updates for stats reporting
	var lastSubmittedCount uint64
	if options.WalletPool != nil && options.WalletPool.GetTxPool() != nil {
		txPool := options.WalletPool.GetTxPool()
		subscriptionID := txPool.SubscribeToBlockUpdates(options.WalletPool, func(blockNumber uint64, walletPoolStats *spamoor.WalletPoolBlockStats) {
			throughputMutex.Lock()
			if newSlotDuration := SlotDuration(options.WalletPool); newSlotDuration != slotDuration {
				slotDuration = newSlotDuration
				limiter.SetLimit(slotRate())
			}
			throughputMutex.Unlock()

			pendingCount := uint64(0)
			submittedCount := uint64(0)
			for _, wallet := range options.WalletPool.GetAllWallets() {
				// Get pending count
				pendingNonce := wallet.GetNonce()
				confirmedNonce := wallet.GetConfirmedNonce()
				if pendingNonce > confirmedNonce {
					pendingCount += pendingNonce - confirmedNonce
				}

				// Get submitted count
				submittedCount += wallet.GetSubmittedTxCount()
			}

			submittedThisBlock := submittedCount - lastSubmittedCount
			lastSubmittedCount = submittedCount

			// Record confirmed transactions for this block
			throughputTracker.recordCompletion(blockNumber, walletPoolStats.ConfirmedTxCount)

			// Calculate average transactions per block over different ranges
			throughput5B := throughputTracker.getAverageThroughput(5, blockNumber)
			throughput20B := throughputTracker.getAverageThroughput(20, blockNumber)
			throughput60B := throughputTracker.getAverageThroughput(60, blockNumber)

			options.Logger.WithField("wallets", walletPoolStats.AffectedWallets).Infof(
				"block %d: submitted=%d, pending=%d, confirmed=%d, throughput: 5B=%.2f tx/B, 20B=%.2f tx/B, 60B=%.2f tx/B",
				blockNumber, submittedThisBlock, pendingCount, walletPoolStats.ConfirmedTxCount, throughput5B, throughput20B, throughput60B,
			)
		})

		defer txPool.UnsubscribeFromBlockUpdates(subscriptionID)
	}

	if options.ThroughputIncrementInterval != 0 {
		// Calculate the ratio between MaxPending and Throughput
		pendingRatio := float64(options.MaxPending) / float64(options.Throughput)

		go func() {
			ticker := time.NewTicker(time.Duration(options.ThroughputIncrementInterval) * time.Second)
			for {
				select {
				case <-ticker.C:
					throughputMutex.Lock()
					oldThroughput := throughput
					throughput++
					newMaxPending := uint64(throughput * pendingRatio)

					options.Logger.Infof("Increasing throughput from %.3f to %.3f and max pending from %d to %d",
						oldThroughput, throughput, maxPending.Load(), newMaxPending)

					limiter.SetLimit(slotRate())
					throughputMutex.Unlock()
					maxPending.Store(newMaxPending)

					// Signal one waiting goroutine that capacity has increased
					if pendingCond != nil {
						pendingCond.Signal()
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	for {
		if err := limiter.Wait(ctx); err != nil {
			if ctx.Err() != nil {
				break
			}

			options.Logger.Debugf("rate limited: %s", err.Error())
			time.Sleep(100 * time.Millisecond)
			continue
		}

		if isErrorMode.Load() {
			if err := errorLimiter.Wait(ctx); err != nil {
				if ctx.Err() != nil {
					break
				}

				options.Logger.Debugf("rate limited: %s", err.Error())
				time.Sleep(100 * time.Millisecond)
				continue // retry
			}
		}

		txIdx := txIdxCounter
		txIdxCounter++

		if options.MaxPending > 0 {
			pendingMutex.Lock()
			for pendingCount.Load() >= int64(maxPending.Load()) {
				pendingCond.Wait()
				if ctx.Err() != nil {
					pendingMutex.Unlock()
					return nil
				}
			}
			pendingMutex.Unlock()
		}
		pendingCount.Add(1)
		pendingWg.Add(1)

		// Initialize state for this transaction
		state := &txState{
			submitted: false,
			done:      false,
		}

		// Register the transaction state
		txStatesMutex.Lock()
		txStates[txIdx] = state
		txStatesMutex.Unlock()

		go func(txIdx uint64, state *txState) {
			defer func() {
				// Mark transaction as done
				txStatesMutex.Lock()
				state.done = true

				// If not submitted yet, mark as submitted to unblock logging
				if !state.submitted {
					state.submitted = true
					processPendingLogs()
				}
				txStatesMutex.Unlock()

				pendingWg.Done()
				pendingCount.Add(-1)
				if pendingCond != nil {
					pendingCond.Signal()
				}
			}()

			// recover only stops an unwinding panic when the recovering function
			// is deferred directly, so this must not be moved into the teardown
			// closure above. Deferred last, it runs first and recovers before the
			// teardown executes.
			defer utils.RecoverPanic(options.Logger, "scenario.processNextTxFn", nil)

			params := &ProcessNextTxParams{
				TxIdx: txIdx,
				NotifySubmitted: func() {
					txStatesMutex.Lock()
					defer txStatesMutex.Unlock()

					if state.submitted {
						return
					}

					state.submitted = true
					processPendingLogs()
				},
				OrderedLogCb: func(logcb func()) {
					txStatesMutex.Lock()
					if state.submitted {
						// Already submitted, so ordered logging has passed this
						// transaction: run the callback immediately.
						txStatesMutex.Unlock()
						logcb()

						return
					}

					state.logCb = append(state.logCb, logcb)
					txStatesMutex.Unlock()
				},
			}

			err := options.ProcessNextTxFn(ctx, params)

			// Every returned call counts toward TotalCount, regardless of whether it
			// submitted a transaction, failed, or completed without ever calling
			// NotifySubmitted. Counting only submissions and errors let a
			// ProcessNextTxFn that finishes successfully without calling
			// NotifySubmitted run forever, since the count never advanced for it.
			txCount.Add(1)

			if err == ErrNoClients {
				isErrorMode.Store(true)
			} else {
				isErrorMode.CompareAndSwap(true, false)
			}
		}(txIdx, state)

		count := txCount.Load() + uint64(pendingCount.Load())
		if options.TotalCount > 0 && count >= options.TotalCount {
			break
		}
	}

	if !options.NoAwaitTransactions {
		pendingWg.Wait()
	}

	return nil
}
