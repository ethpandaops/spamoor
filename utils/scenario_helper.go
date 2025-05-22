package utils

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

// TransactionScenarioOptions configures how the transaction scenario is executed.
type TransactionScenarioOptions struct {
	TotalCount                  uint64
	Throughput                  uint64
	MaxPending                  uint64
	ThroughputIncrementInterval uint64

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

	var lastChan chan bool
	var pendingChan chan bool

	if options.MaxPending > 0 {
		pendingChan = make(chan bool, options.MaxPending)
	}

	initialRate := rate.Limit(float64(options.Throughput) / float64(SecondsPerSlot))
	if initialRate == 0 {
		initialRate = rate.Inf
	}
	limiter := rate.NewLimiter(initialRate, 1)

	if options.ThroughputIncrementInterval != 0 {
		go func() {
			ticker := time.NewTicker(time.Duration(options.ThroughputIncrementInterval) * time.Second)
			for {
				select {
				case <-ticker.C:
					throughput := limiter.Limit() * 12
					newThroughput := throughput + 1
					options.Logger.Infof("Increasing throughput from %.3f to %.3f", throughput, newThroughput)
					limiter.SetLimit(rate.Limit(float64(newThroughput) / float64(SecondsPerSlot)))
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

		if pendingChan != nil {
			// await pending transactions
			pendingChan <- true
		}
		pendingCount.Add(1)

		currentChan := make(chan bool, 1)

		go func(txIdx uint64, lastChan, currentChan chan bool) {
			defer func() {
				RecoverPanic(options.Logger, "scenario.processNextTxFn")
				currentChan <- true
			}()

			logcb, err := options.ProcessNextTxFn(ctx, txIdx, func() {
				pendingCount.Add(-1)
				if pendingChan != nil {
					time.Sleep(10 * time.Millisecond)
					<-pendingChan
				}
			})

			if lastChan != nil {
				<-lastChan
				close(lastChan)
			}

			if err == nil {
				txCount.Add(1)
			}

			if logcb != nil {
				logcb()
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

	return nil
}
