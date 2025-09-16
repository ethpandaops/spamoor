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

// GlobalSecondsPerSlot is the global setting for seconds per slot used in rate limiting.
// This can be set via CLI flag (--seconds-per-slot) and applies to all scenarios.
var GlobalSecondsPerSlot uint64 = 12

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
	OrderedLogCb    func(logFunc func())
	NotifySubmitted func()
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
	// Use global SecondsPerSlot
	secondsPerSlot := GlobalSecondsPerSlot

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

	// Helper function to process pending logs in order
	processPendingLogs := func() {
		txStatesMutex.Lock()
		defer txStatesMutex.Unlock()

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

			// Clean up if the transaction is done
			if state.done {
				delete(txStates, nextLogIdx)
			}

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

	initialRate := rate.Limit(float64(options.Throughput) / float64(secondsPerSlot))
	if initialRate == 0 {
		initialRate = rate.Inf
	}
	limiter := rate.NewLimiter(initialRate, 1)

	isErrorMode := false
	errorLimiter := rate.NewLimiter(rate.Limit(0.5), 1) // 2 sec interval when in error mode

	// Subscribe to block updates for stats reporting
	var lastSubmittedCount uint64
	if options.WalletPool != nil && options.WalletPool.GetTxPool() != nil {
		txPool := options.WalletPool.GetTxPool()
		subscriptionID := txPool.SubscribeToBlockUpdates(options.WalletPool, func(blockNumber uint64, walletPoolStats *spamoor.WalletPoolBlockStats) {
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
					throughput := limiter.Limit() * 12
					newThroughput := throughput + 1
					newMaxPending := uint64(float64(newThroughput) * pendingRatio)

					options.Logger.Infof("Increasing throughput from %.3f to %.3f and max pending from %d to %d",
						throughput, newThroughput, maxPending.Load(), newMaxPending)

					limiter.SetLimit(rate.Limit(float64(newThroughput) / float64(secondsPerSlot)))
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

		if isErrorMode {
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
				utils.RecoverPanic(options.Logger, "scenario.processNextTxFn", nil)

				// Mark transaction as done
				state.done = true

				// If not submitted yet, mark as submitted to unblock logging
				if !state.submitted {
					state.submitted = true
					processPendingLogs()
				}

				pendingWg.Done()
				pendingCount.Add(-1)
				if pendingCond != nil {
					pendingCond.Signal()
				}
			}()

			params := &ProcessNextTxParams{
				TxIdx: txIdx,
				NotifySubmitted: func() {
					if state.submitted {
						return
					}

					state.submitted = true
					processPendingLogs()
				},
				OrderedLogCb: func(logcb func()) {
					// If already submitted, process logs immediately
					if state.submitted {
						logcb()
					} else {
						state.logCb = append(state.logCb, logcb)
					}
				},
			}

			err := options.ProcessNextTxFn(ctx, params)

			if err != nil || state.submitted {
				txCount.Add(1)
			}

			if err == ErrNoClients {
				isErrorMode = true
			} else if isErrorMode {
				isErrorMode = false
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
