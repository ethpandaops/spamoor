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

const SecondsPerSlot uint64 = 12

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
	// - The onComplete callback must be called when transaction processing is complete
	ProcessNextTxFn func(ctx context.Context, txIdx uint64, onComplete func()) (func(), error)
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

	var lastChan chan bool
	var maxPending atomic.Uint64
	var pendingMutex sync.Mutex
	var pendingCond *sync.Cond

	pendingWg := sync.WaitGroup{}

	if options.MaxPending > 0 {
		maxPending.Store(options.MaxPending)
		pendingCond = sync.NewCond(&pendingMutex)
	}

	initialRate := rate.Limit(float64(options.Throughput) / float64(SecondsPerSlot))
	if initialRate == 0 {
		initialRate = rate.Inf
	}
	limiter := rate.NewLimiter(initialRate, 1)

	// Subscribe to block updates for stats reporting
	var lastSubmittedCount uint64
	if options.WalletPool != nil && options.WalletPool.GetTxPool() != nil {
		txPool := options.WalletPool.GetTxPool()
		subscriptionID := txPool.SubscribeToBlockUpdates(options.WalletPool, func(blockNumber uint64, walletPoolStats *spamoor.WalletPoolBlockStats) {
			currentSubmitted := txCount.Load()
			submittedThisBlock := currentSubmitted - lastSubmittedCount
			lastSubmittedCount = currentSubmitted
			pending := pendingCount.Load()

			// Record confirmed transactions for this block
			throughputTracker.recordCompletion(blockNumber, walletPoolStats.ConfirmedTxCount)

			// Calculate average transactions per block over different ranges
			throughput5B := throughputTracker.getAverageThroughput(5, blockNumber)
			throughput20B := throughputTracker.getAverageThroughput(20, blockNumber)
			throughput60B := throughputTracker.getAverageThroughput(60, blockNumber)

			options.Logger.WithField("wallets", walletPoolStats.AffectedWallets).Infof(
				"block %d: submitted=%d, pending=%d, confirmed=%d, throughput: 5B=%.2f tx/B, 20B=%.2f tx/B, 60B=%.2f tx/B",
				blockNumber, submittedThisBlock, pending, walletPoolStats.ConfirmedTxCount, throughput5B, throughput20B, throughput60B,
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

					limiter.SetLimit(rate.Limit(float64(newThroughput) / float64(SecondsPerSlot)))
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

		currentChan := make(chan bool, 1)

		go func(txIdx uint64, lastChan, currentChan chan bool) {
			defer func() {
				utils.RecoverPanic(options.Logger, "scenario.processNextTxFn", nil)
				currentChan <- true
			}()

			completed := false

			logcb, err := options.ProcessNextTxFn(ctx, txIdx, func() {
				if completed {
					return
				}

				completed = true

				pendingWg.Done()
				pendingCount.Add(-1)
				txCount.Add(1)
				if pendingCond != nil {
					pendingCond.Signal()
				}
			})

			if lastChan != nil {
				<-lastChan
				close(lastChan)
			}

			if logcb != nil {
				logcb()
			} else if err != nil {
				options.Logger.Warnf("process next tx failed: %v", err)
			}
		}(txIdx, lastChan, currentChan)

		lastChan = currentChan

		count := txCount.Load() + uint64(pendingCount.Load())
		if options.TotalCount > 0 && count >= options.TotalCount {
			break
		}
	}

	if lastChan != nil {
		<-lastChan
		close(lastChan)
	}

	if !options.NoAwaitTransactions {
		pendingWg.Wait()
	}

	return nil
}
