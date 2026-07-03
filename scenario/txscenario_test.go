package scenario

import (
	"context"
	"io"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func discardLogger() *logrus.Entry {
	lg := logrus.New()
	lg.SetOutput(io.Discard)
	return lg.WithField("test", "txscenario")
}

// TestTransactionStateReleased runs a large number of transactions through the
// scenario and checks that per-transaction bookkeeping does not accumulate. A
// leak would show up as live heap objects growing roughly in step with the
// number of processed transactions.
func TestTransactionStateReleased(t *testing.T) {
	const total = 20000

	var baseline, peak uint64

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := RunTransactionScenario(ctx, TransactionScenarioOptions{
		TotalCount: total,
		Throughput: 60000, // paced so only a few transactions are in flight at a time
		MaxPending: 0,
		Logger:     discardLogger(),
		ProcessNextTxFn: func(_ context.Context, params *ProcessNextTxParams) error {
			params.NotifySubmitted()
			switch params.TxIdx {
			case 2000:
				runtime.GC()
				var m runtime.MemStats
				runtime.ReadMemStats(&m)
				baseline = m.HeapObjects
			case total - 1:
				runtime.GC()
				var m runtime.MemStats
				runtime.ReadMemStats(&m)
				peak = m.HeapObjects
			}
			return nil
		},
	})
	if err != nil {
		t.Fatalf("scenario error: %v", err)
	}
	if baseline == 0 || peak == 0 {
		t.Fatalf("measurements not captured (baseline=%d peak=%d)", baseline, peak)
	}

	growth := int64(peak) - int64(baseline)
	if growth > total/4 {
		t.Fatalf("live heap objects grew by %d over %d transactions; per-transaction state is not being released", growth, total-2000)
	}
	t.Logf("live heap objects grew by %d over %d transactions (bounded)", growth, total-2000)
}

// TestOrderedLoggingInOrder verifies that ordered log callbacks still fire
// exactly once, in submission order, so the cleanup change does not affect
// ordered logging.
func TestOrderedLoggingInOrder(t *testing.T) {
	const total = 1000

	var mu sync.Mutex
	order := make([]uint64, 0, total)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := RunTransactionScenario(ctx, TransactionScenarioOptions{
		TotalCount: total,
		Throughput: 60000,
		MaxPending: 0,
		Logger:     discardLogger(),
		ProcessNextTxFn: func(_ context.Context, params *ProcessNextTxParams) error {
			idx := params.TxIdx
			params.OrderedLogCb(func() {
				mu.Lock()
				order = append(order, idx)
				mu.Unlock()
			})
			params.NotifySubmitted()
			return nil
		},
	})
	if err != nil {
		t.Fatalf("scenario error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(order) != total {
		t.Fatalf("expected %d ordered log callbacks, got %d", total, len(order))
	}
	for i, v := range order {
		if v != uint64(i) {
			t.Fatalf("ordered log out of order at position %d: got %d", i, v)
		}
	}
}

// TestScenarioOrderedLogAfterSubmit verifies that a log callback registered
// after the transaction has been submitted runs immediately instead of being
// buffered.
func TestScenarioOrderedLogAfterSubmit(t *testing.T) {
	var logged atomic.Uint64

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := RunTransactionScenario(ctx, TransactionScenarioOptions{
		TotalCount: 100,
		Throughput: 60000,
		Logger:     discardLogger(),
		ProcessNextTxFn: func(_ context.Context, params *ProcessNextTxParams) error {
			params.NotifySubmitted()
			params.OrderedLogCb(func() {
				logged.Add(1)
			})
			return nil
		},
	})
	if err != nil {
		t.Fatalf("scenario error: %v", err)
	}
	if logged.Load() != 100 {
		t.Fatalf("expected 100 log callbacks, got %d", logged.Load())
	}
}

// TestScenarioTimeout runs an unbounded scenario with a timeout and verifies
// that it stops on its own once the timeout elapses.
func TestScenarioTimeout(t *testing.T) {
	var processed atomic.Uint64
	start := time.Now()

	err := RunTransactionScenario(context.Background(), TransactionScenarioOptions{
		TotalCount: 0,
		Throughput: 60000,
		Timeout:    300 * time.Millisecond,
		Logger:     discardLogger(),
		ProcessNextTxFn: func(_ context.Context, params *ProcessNextTxParams) error {
			params.NotifySubmitted()
			processed.Add(1)
			return nil
		},
	})
	if err != nil {
		t.Fatalf("scenario error: %v", err)
	}
	if processed.Load() == 0 {
		t.Fatal("expected transactions to be processed before the timeout")
	}
	if elapsed := time.Since(start); elapsed > 10*time.Second {
		t.Fatalf("scenario did not stop after the timeout (ran %v)", elapsed)
	}
}

// TestScenarioMaxPendingLimit verifies that no more than MaxPending
// transactions are processed concurrently.
func TestScenarioMaxPendingLimit(t *testing.T) {
	const total = 12
	const maxPending = 3

	var inFlight, maxObserved atomic.Int64

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := RunTransactionScenario(ctx, TransactionScenarioOptions{
		TotalCount: total,
		Throughput: 60000,
		MaxPending: maxPending,
		Logger:     discardLogger(),
		ProcessNextTxFn: func(_ context.Context, params *ProcessNextTxParams) error {
			current := inFlight.Add(1)
			for {
				peak := maxObserved.Load()
				if current <= peak || maxObserved.CompareAndSwap(peak, current) {
					break
				}
			}

			time.Sleep(10 * time.Millisecond)
			inFlight.Add(-1)
			params.NotifySubmitted()
			return nil
		},
	})
	if err != nil {
		t.Fatalf("scenario error: %v", err)
	}
	if peak := maxObserved.Load(); peak > maxPending {
		t.Fatalf("observed %d concurrent transactions, limit is %d", peak, maxPending)
	}
}

// TestScenarioMaxPendingCancel verifies that a scenario blocked on the
// pending limit stops when its context is cancelled.
func TestScenarioMaxPendingCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	started := make(chan struct{}, 1)
	done := make(chan error, 1)

	go func() {
		done <- RunTransactionScenario(ctx, TransactionScenarioOptions{
			TotalCount: 0,
			Throughput: 60000,
			MaxPending: 1,
			Logger:     discardLogger(),
			ProcessNextTxFn: func(txCtx context.Context, _ *ProcessNextTxParams) error {
				select {
				case started <- struct{}{}:
				default:
				}
				<-txCtx.Done()
				return txCtx.Err()
			},
		})
	}()

	<-started
	// Give the main loop time to block on the pending limit before cancelling
	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("scenario error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("scenario did not stop after context cancellation")
	}
}

// TestScenarioClientErrorRecovery verifies that a transaction failing with
// ErrNoClients puts the scenario into error mode and that it recovers once
// transactions succeed again.
func TestScenarioClientErrorRecovery(t *testing.T) {
	var succeeded atomic.Uint64

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := RunTransactionScenario(ctx, TransactionScenarioOptions{
		TotalCount: 4,
		Throughput: 240, // paced so each result lands before the next transaction starts
		Logger:     discardLogger(),
		ProcessNextTxFn: func(_ context.Context, params *ProcessNextTxParams) error {
			if params.TxIdx == 0 {
				return ErrNoClients
			}

			params.NotifySubmitted()
			succeeded.Add(1)
			return nil
		},
	})
	if err != nil {
		t.Fatalf("scenario error: %v", err)
	}
	if count := succeeded.Load(); count < 2 || count > 3 {
		t.Fatalf("expected 2-3 successful transactions after error recovery, got %d", count)
	}
}

// TestScenarioPanicRecovery verifies that a panicking transaction handler
// does not crash the scenario and that processing continues afterwards.
func TestScenarioPanicRecovery(t *testing.T) {
	var succeeded atomic.Uint64

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := RunTransactionScenario(ctx, TransactionScenarioOptions{
		TotalCount: 3,
		Throughput: 240, // paced so each result lands before the next transaction starts
		Logger:     discardLogger(),
		ProcessNextTxFn: func(_ context.Context, params *ProcessNextTxParams) error {
			if params.TxIdx == 0 {
				panic("transaction handler panic")
			}

			params.NotifySubmitted()
			succeeded.Add(1)
			return nil
		},
	})
	if err != nil {
		t.Fatalf("scenario error: %v", err)
	}
	if succeeded.Load() < 2 {
		t.Fatalf("expected the scenario to continue after a panic, got %d successful transactions", succeeded.Load())
	}
}

// TestScenarioNoAwaitTransactions verifies that the scenario returns without
// waiting for in-flight transactions when NoAwaitTransactions is set.
func TestScenarioNoAwaitTransactions(t *testing.T) {
	release := make(chan struct{})
	finished := sync.WaitGroup{}
	finished.Add(2)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := RunTransactionScenario(ctx, TransactionScenarioOptions{
		TotalCount:          2,
		Throughput:          60000,
		NoAwaitTransactions: true,
		Logger:              discardLogger(),
		ProcessNextTxFn: func(_ context.Context, params *ProcessNextTxParams) error {
			defer finished.Done()
			params.NotifySubmitted()
			<-release
			return nil
		},
	})
	if err != nil {
		t.Fatalf("scenario error: %v", err)
	}

	// The scenario returned while both transactions are still blocked;
	// release them and wait so no goroutines outlive the test.
	close(release)
	finished.Wait()
}

// TestScenarioThroughputIncrement runs a scenario with a periodic throughput
// increase and verifies that it keeps processing transactions.
func TestScenarioThroughputIncrement(t *testing.T) {
	var processed atomic.Uint64

	err := RunTransactionScenario(context.Background(), TransactionScenarioOptions{
		TotalCount:                  0,
		Throughput:                  60, // 5 tx/s at the default 12s slot time
		MaxPending:                  3,
		ThroughputIncrementInterval: 1,
		Timeout:                     2200 * time.Millisecond,
		Logger:                      discardLogger(),
		ProcessNextTxFn: func(_ context.Context, params *ProcessNextTxParams) error {
			params.NotifySubmitted()
			processed.Add(1)
			return nil
		},
	})
	if err != nil {
		t.Fatalf("scenario error: %v", err)
	}
	if processed.Load() < 5 {
		t.Fatalf("expected at least 5 transactions to be processed, got %d", processed.Load())
	}
}
