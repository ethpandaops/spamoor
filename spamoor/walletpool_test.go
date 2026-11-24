package spamoor

import (
	"context"
	"math/big"
	"testing"
	"time"
)

// Funding fees are derived from the client's cached fee suggestion without changing it,
// and the whole transaction stays below the fee cap clients enforce on their RPC.
func TestGetFundingFees(t *testing.T) {
	gwei := func(n int64) *big.Int { return new(big.Int).Mul(big.NewInt(n), big.NewInt(1000000000)) }

	tests := []struct {
		name       string
		suggested  *big.Int
		tip        *big.Int
		gas        uint64
		wantFeeCap *big.Int
		wantTipCap *big.Int
	}{
		{"headroom added", gwei(20), gwei(1), 21000, gwei(22), big.NewInt(1100000000)},
		{"capped at 100 gwei", gwei(500), gwei(1), 21000, gwei(100), big.NewInt(1100000000)},
		{"capped by the transaction fee", gwei(60), gwei(1), 16000000, big.NewInt(61875000000), big.NewInt(1100000000)},
		{"tip never above the fee cap", gwei(1), gwei(50), 21000, gwei(3), gwei(3)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				lastGasSuggestion: time.Now(),
				lastGasCap:        new(big.Int).Set(tt.suggested),
				lastTipCap:        new(big.Int).Set(tt.tip),
			}
			pool := &WalletPool{ctx: context.Background()}

			for range 3 {
				feeCap, tipCap, err := pool.getFundingFees(client, tt.gas)
				if err != nil {
					t.Fatal(err)
				}

				if feeCap.Cmp(tt.wantFeeCap) != 0 || tipCap.Cmp(tt.wantTipCap) != 0 {
					t.Fatalf("got fee cap %v / tip cap %v, want %v / %v", feeCap, tipCap, tt.wantFeeCap, tt.wantTipCap)
				}

				if fee := new(big.Int).Mul(feeCap, new(big.Int).SetUint64(tt.gas)); fee.Cmp(maxFundingTxFee) > 0 {
					t.Fatalf("transaction fee %v exceeds %v", fee, maxFundingTxFee)
				}
			}

			if client.lastGasCap.Cmp(tt.suggested) != 0 || client.lastTipCap.Cmp(tt.tip) != 0 {
				t.Fatalf("cached suggestion changed to %v / %v", client.lastGasCap, client.lastTipCap)
			}
		})
	}
}

// A funding transaction too large to pay the chain's fee within the fee limit is not
// sent underpriced, where it would hold a root wallet nonce without ever being included.
func TestGetFundingFeesRejectsUnpayableTx(t *testing.T) {
	client := &Client{
		lastGasSuggestion: time.Now(),
		lastGasCap:        big.NewInt(80000000000),
		lastTipCap:        big.NewInt(1000000000),
	}
	pool := &WalletPool{ctx: context.Background()}

	if _, _, err := pool.getFundingFees(client, 16000000); err == nil {
		t.Fatal("expected an error for a transaction that cannot pay 80 gwei within the fee limit")
	}

	if _, _, err := pool.getFundingFees(client, 10000000); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
