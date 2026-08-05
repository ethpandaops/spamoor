package spamoor

import (
	"context"
	"math/big"
	"testing"

	"github.com/holiman/uint256"
)

func rootWalletWithBalance(balanceWei *big.Int) *RootWallet {
	return &RootWallet{
		wallet:              &Wallet{balance: balanceWei},
		pendingFundingTotal: uint256.NewInt(0),
		txSemaphore:         make(chan struct{}, 16),
		txSemLimit:          16,
	}
}

func weiForEth(n int64) *uint256.Int {
	return uint256.MustFromBig(new(big.Int).Mul(big.NewInt(n), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)))
}

// A context cancelled while WithWalletLock is still waiting for sufficient
// balance must leave pendingFundingTotal untouched: the reservation was never
// added, so the deferred release must not subtract it either.
func TestWithWalletLock_CancelledWaitLeavesAccountingIntact(t *testing.T) {
	// Zero balance means waitForSufficientBalance can never succeed, so it
	// blocks until ctx is done.
	rw := rootWalletWithBalance(big.NewInt(0))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	lockedFnRan := false
	err := rw.WithWalletLock(ctx, 1, weiForEth(5), nil, nil, func() error {
		lockedFnRan = true
		return nil
	})

	if err == nil {
		t.Fatal("expected WithWalletLock to fail on a cancelled wait")
	}
	if lockedFnRan {
		t.Fatal("lockedFn must not run when acquiring the lock fails")
	}
	if !rw.pendingFundingTotal.IsZero() {
		t.Fatalf("expected pendingFundingTotal to stay at 0, got %s", rw.pendingFundingTotal.Dec())
	}
}

// A normal, successful call still nets the reservation back to zero once
// released - the fix must not change behavior for the case that already
// worked.
func TestWithWalletLock_SuccessfulCallReleasesReservation(t *testing.T) {
	rw := rootWalletWithBalance(weiForEth(10).ToBig())

	lockedFnRan := false
	err := rw.WithWalletLock(context.Background(), 1, weiForEth(5), nil, nil, func() error {
		lockedFnRan = true
		if rw.pendingFundingTotal.Cmp(weiForEth(5)) != 0 {
			t.Fatalf("expected pendingFundingTotal to be reserved at 5 ETH while locked, got %s", rw.pendingFundingTotal.Dec())
		}
		return nil
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !lockedFnRan {
		t.Fatal("expected lockedFn to run")
	}
	if !rw.pendingFundingTotal.IsZero() {
		t.Fatalf("expected pendingFundingTotal to be released back to 0, got %s", rw.pendingFundingTotal.Dec())
	}
}

// The balance guard must still correctly refuse an oversized request after a
// cancelled wait - this is the actual safety property the accounting exists
// to protect, and the fix must restore it rather than just fixing the number.
func TestWithWalletLock_BalanceGuardStillWorksAfterCancelledWait(t *testing.T) {
	rw := rootWalletWithBalance(weiForEth(3).ToBig()) // wallet genuinely holds 3 ETH

	// A cancelled wait against an unrelated, unsatisfiable request must not
	// corrupt accounting used by any later request on this wallet.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_ = rw.WithWalletLock(ctx, 1, weiForEth(100), nil, nil, func() error { return nil })

	if !rw.pendingFundingTotal.IsZero() {
		t.Fatalf("expected pendingFundingTotal to be 0 after the cancelled wait, got %s", rw.pendingFundingTotal.Dec())
	}

	ok, _ := rw.hasSufficientBalance(weiForEth(50), 1)
	if ok {
		t.Fatal("expected the guard to correctly refuse 50 ETH against a real balance of 3 ETH")
	}
}
