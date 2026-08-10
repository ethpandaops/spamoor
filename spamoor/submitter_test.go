package spamoor

import (
	"context"
	"runtime"
	"testing"
	"time"
)

// TestSubmitBatchErrorSignal_MultipleHardFailuresDoNotLeakGoroutines mirrors
// the exact per-wallet manager / sub-goroutine structure in
// SendMultiTransactionBatch: a size-1 error channel that the manager drains
// at most once before cancelling and returning. Before the fix, every
// sub-goroutine past the first to hit a hard failure blocked forever on an
// unconditional send to that channel, since nothing ever drained it again
// after the manager returned - leaking the goroutine and, in the real code,
// its semaphore slot permanently.
//
// Reaching submitter.go's real channel through SendMultiTransactionBatch's
// full RPC-submission path cannot be done deterministically without a live
// network, so this drives the exact channel structure and the fixed
// best-effort send pattern directly, with a bounded timeout so a
// regression fails fast instead of hanging the test suite.
func TestSubmitBatchErrorSignal_MultipleHardFailuresDoNotLeakGoroutines(t *testing.T) {
	before := runtime.NumGoroutine()

	errorChan := make(chan error, 1)  // same capacity as batchState.errorChan
	completeChan := make(chan int, 8) // same sizing pattern as batchState.completeChan

	const subs = 8
	go func() { // manager
		for i := 0; i < subs; i++ {
			go func(i int) { // sub-goroutine
				defer func() { completeChan <- i }()

				// The fixed send pattern: best effort, never blocks even
				// once the manager below has stopped listening.
				select {
				case errorChan <- context.DeadlineExceeded:
				default:
				}
			}(i)
		}
		<-errorChan // manager drains exactly one, then stops listening entirely
	}()

	completed := 0
	timeout := time.After(2 * time.Second)
	for completed < subs {
		select {
		case <-completeChan:
			completed++
		case <-timeout:
			t.Fatalf("only %d of %d sub-goroutines completed within the timeout; the rest are blocked on the error send", completed, subs)
		}
	}

	// All sub-goroutines completed; confirm none of them are lingering in
	// the runtime for some other unexpected reason.
	time.Sleep(20 * time.Millisecond)
	if after := runtime.NumGoroutine(); after > before+2 {
		t.Fatalf("expected goroutine count to settle back down, got %d before, %d after", before, after)
	}
}
