package scenario

import (
	"context"
	"io"
	"runtime"
	"sync"
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
