package scenario

import (
	"context"
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
)

func TestReceiptChanWaitNilChannel(t *testing.T) {
	var rc ReceiptChan

	receipt, err := rc.Wait(context.Background())
	if err != nil {
		t.Fatalf("expected no error for nil channel, got %v", err)
	}
	if receipt != nil {
		t.Fatalf("expected nil receipt for nil channel, got %v", receipt)
	}
}

func TestReceiptChanWaitReceivesReceipt(t *testing.T) {
	rc := make(ReceiptChan, 1)
	want := &types.Receipt{Status: types.ReceiptStatusSuccessful}
	rc <- want

	got, err := rc.Wait(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got != want {
		t.Fatalf("expected the sent receipt to be returned, got %v", got)
	}
}

func TestReceiptChanWaitContextCancelled(t *testing.T) {
	rc := make(ReceiptChan)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	receipt, err := rc.Wait(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if receipt != nil {
		t.Fatalf("expected nil receipt on cancellation, got %v", receipt)
	}
}
