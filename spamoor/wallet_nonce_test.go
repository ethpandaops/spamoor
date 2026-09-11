package spamoor

import "testing"

func TestAdvancePendingTxCount(t *testing.T) {
	tests := []struct {
		name    string
		current uint64
		target  uint64
		want    uint64
	}{
		{name: "advances when target is ahead", current: 10, target: 13, want: 13},
		{name: "never moves backwards", current: 15, target: 13, want: 15},
		{name: "no-op when equal", current: 13, target: 13, want: 13},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Wallet{}
			w.pendingTxCount.Store(tt.current)

			w.advancePendingTxCount(tt.target)

			if got := w.pendingTxCount.Load(); got != tt.want {
				t.Fatalf("pendingTxCount = %d, want %d", got, tt.want)
			}
		})
	}
}
