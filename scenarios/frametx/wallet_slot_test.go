package frametx

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestWalletSlotWaitsForConfirmation(t *testing.T) {
	s := &Scenario{walletSlots: []chan struct{}{make(chan struct{}, 1), make(chan struct{}, 1)}}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	release, err := s.acquireWallet(ctx, 0)
	if err != nil {
		t.Fatal(err)
	}
	// Another sender can progress while the first waits for confirmation.
	releaseOther, err := s.acquireWallet(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	releaseOther()

	// Wrapping the round-robin index must not issue another nonce for the
	// first sender until its original transaction completes.
	waitCtx, stop := context.WithTimeout(ctx, 20*time.Millisecond)
	defer stop()
	if releaseAgain, err := s.acquireWallet(waitCtx, 2); !errors.Is(err, context.DeadlineExceeded) {
		if releaseAgain != nil {
			releaseAgain()
		}
		t.Fatalf("second transaction was not held until confirmation: %v", err)
	}
	release()
	releaseAgain, err := s.acquireWallet(ctx, 2)
	if err != nil {
		t.Fatal(err)
	}
	releaseAgain()
}
