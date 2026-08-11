package spamoor

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

// TestSendTransaction_OnCompleteFiresWhenContextAlreadyCancelled drives the
// real SendTransaction/submitTransaction path with a context that is already
// cancelled before the call is made. OnComplete is documented as "always
// called regardless of success/failure", but the early ctx.Err() guard used
// to return before ever setting up the goroutine that calls it. Any caller
// that relies on OnComplete firing to release a WaitGroup or semaphore (see
// WalletPool.ReclaimFunds) would hang forever in that case.
func TestSendTransaction_OnCompleteFiresWhenContextAlreadyCancelled(t *testing.T) {
	pool := &TxPool{}
	wallet := &Wallet{}
	tx := types.NewTx(&types.LegacyTx{Nonce: 0})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var called bool
	var gotErr error
	done := make(chan struct{})

	err := pool.SendTransaction(ctx, wallet, tx, &SendTransactionOptions{
		OnComplete: func(_ *types.Transaction, _ *types.Receipt, err error) {
			called = true
			gotErr = err
			close(done)
		},
	})

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("OnComplete was never called for an already-cancelled context")
	}

	if !called {
		t.Fatal("OnComplete was not called")
	}
	if gotErr == nil {
		t.Fatal("expected OnComplete to receive the context error, got nil")
	}
	if err == nil {
		t.Fatal("expected SendTransaction to return the context error")
	}
}

// TestReclaimFundsPattern_DoesNotHangWhenContextAlreadyCancelled mirrors the
// exact synchronization structure ReclaimFunds (walletpool.go) uses: one
// goroutine per transaction calls the real SendTransaction, and wg.Done() is
// only ever called from inside OnComplete - the dispatching goroutine itself
// does not call it. ReclaimFunds needs a live client and funded wallets to
// reach this point, which is infeasible to set up here, so this drives the
// same wg/OnComplete wiring directly against the real SendTransaction with a
// context that is already cancelled, exactly as it would be if a daemon
// shutdown or a spammer Pause cancelled the context while a batch of reclaim
// transactions was still being dispatched.
func TestReclaimFundsPattern_DoesNotHangWhenContextAlreadyCancelled(t *testing.T) {
	pool := &TxPool{}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	const txCount = 50
	var wg sync.WaitGroup
	wg.Add(txCount)

	for i := 0; i < txCount; i++ {
		wallet := &Wallet{}
		tx := types.NewTx(&types.LegacyTx{Nonce: uint64(i)})

		go func() {
			pool.SendTransaction(ctx, wallet, tx, &SendTransactionOptions{
				OnComplete: func(_ *types.Transaction, _ *types.Receipt, _ error) {
					wg.Done()
				},
			})
		}()
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("ReclaimFunds-style wg.Wait() hung after context was already cancelled before dispatch")
	}
}
