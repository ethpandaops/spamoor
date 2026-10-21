package spamoor

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/ethpandaops/spamoor/txtypes"
)

// OnComplete is documented as always being called once processing ends. The early
// return for an already-cancelled context used to skip it, which left callers that
// release a WaitGroup from the callback (WalletPool.ReclaimFunds) hanging forever.
func TestSendTransactionCallsOnCompleteWhenContextCancelled(t *testing.T) {
	pool := &TxPool{}
	wallet := &Wallet{}
	tx := txtypes.NewTx(&txtypes.LegacyTx{Nonce: 0})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	done := make(chan error, 1)
	err := pool.SendTransaction(ctx, wallet, tx, &SendTransactionOptions{
		OnComplete: func(_ *txtypes.Transaction, receipt *txtypes.Receipt, err error) {
			if receipt != nil {
				t.Errorf("expected no receipt for a cancelled submission")
			}
			done <- err
		},
	})
	if err == nil {
		t.Fatal("expected SendTransaction to return the context error")
	}

	select {
	case cbErr := <-done:
		if cbErr == nil {
			t.Fatal("expected OnComplete to receive the context error")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("OnComplete was not called for an already-cancelled context")
	}
}

// Mirrors the ReclaimFunds wiring: one send per wallet, wg.Done only from OnComplete.
func TestSendTransactionReleasesWaitGroupWhenContextCancelled(t *testing.T) {
	pool := &TxPool{}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	const txCount = 20
	var wg sync.WaitGroup
	wg.Add(txCount)

	for i := range txCount {
		wallet := &Wallet{}
		tx := txtypes.NewTx(&txtypes.LegacyTx{Nonce: uint64(i)})

		go func() {
			_ = pool.SendTransaction(ctx, wallet, tx, &SendTransactionOptions{
				OnComplete: func(_ *txtypes.Transaction, _ *txtypes.Receipt, _ error) {
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
		t.Fatal("wait group was never released after the context was cancelled")
	}
}

func TestUpdateSlotDuration(t *testing.T) {
	pool := &TxPool{options: &TxPoolOptions{Logger: logrus.NewEntry(logrus.New()), SlotDuration: 6 * time.Second}}

	pool.updateSlotDuration(&BlockInfo{Timestamp: 1000}, &BlockInfo{Timestamp: 1012})
	if pool.GetSlotDuration() != 6*time.Second {
		t.Fatalf("expected the configured fallback without slot numbers, got %v", pool.GetSlotDuration())
	}

	pool.updateSlotDuration(&BlockInfo{Timestamp: 1000, SlotNumber: 100}, &BlockInfo{Timestamp: 1012, SlotNumber: 101})
	if pool.GetSlotDuration() != 12*time.Second {
		t.Fatalf("expected 12s, got %v", pool.GetSlotDuration())
	}

	pool.updateSlotDuration(&BlockInfo{Timestamp: 2000, SlotNumber: 200}, &BlockInfo{Timestamp: 2027, SlotNumber: 203})
	if pool.GetSlotDuration() != 9*time.Second {
		t.Fatalf("expected 9s across missed slots, got %v", pool.GetSlotDuration())
	}
}
