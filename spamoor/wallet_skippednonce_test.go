package spamoor

import "testing"

// A transaction that fails to submit keeps a nonce at or above the confirmed
// count. MarkSkippedNonce should record it so the next allocation reuses it
// instead of leaving a permanent gap.
func TestMarkSkippedNonceReusesFailedNonce(t *testing.T) {
	w := &Wallet{}
	w.confirmedTxCount = 100
	w.pendingTxCount.Store(100)

	n := w.GetNextNonce() // allocates 100, pending count moves to 101
	if n != 100 {
		t.Fatalf("expected first allocation 100, got %d", n)
	}

	w.MarkSkippedNonce(n) // submit failed, make the nonce reusable

	if reused := w.GetNextNonce(); reused != 100 {
		t.Fatalf("expected nonce 100 to be reused, got %d", reused)
	}
}

// Nonces at or above the confirmed count are still pending, so they should be
// recorded for reuse.
func TestMarkSkippedNonceRecordsPendingNonces(t *testing.T) {
	w := &Wallet{}
	w.confirmedTxCount = 100
	w.pendingTxCount.Store(103)

	w.MarkSkippedNonce(101)
	w.MarkSkippedNonce(102)

	if len(w.skippedNonces) != 2 {
		t.Fatalf("expected 101 and 102 to be recorded, got %v", w.skippedNonces)
	}
}

// Nonces below the confirmed count are already confirmed on chain and cannot be
// reused, so they should be ignored.
func TestMarkSkippedNonceIgnoresConfirmedNonces(t *testing.T) {
	w := &Wallet{}
	w.confirmedTxCount = 100
	w.pendingTxCount.Store(105)

	w.MarkSkippedNonce(50)

	if len(w.skippedNonces) != 0 {
		t.Fatalf("expected a confirmed nonce to be ignored, got %v", w.skippedNonces)
	}
}
