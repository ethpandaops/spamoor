package erc20bloater

import (
	"testing"

	"github.com/ethpandaops/spamoor/spamoor"
)

func TestAddressesPerBloatTx(t *testing.T) {
	tests := []struct {
		name             string
		costPerStateByte uint64
		want             uint64
	}{
		{
			name:             "legacy fee model keeps the pre-Amsterdam batch size",
			costPerStateByte: 0,
			want:             MaxBloatedAddressesPerTx,
		},
		{
			// (16,700,000 - 300,000) / (44,400 + 2 * 64 * 1530) = 68
			name:             "amsterdam flat cost per state byte",
			costPerStateByte: spamoor.CostPerStateByte,
			want:             68,
		},
		{
			name:             "tiny cost per state byte stays close to the pre-Amsterdam batch size",
			costPerStateByte: 1,
			want:             368,
		},
		{
			name:             "huge cost per state byte still makes progress",
			costPerStateByte: 1 << 40,
			want:             1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addressesPerBloatTx(tt.costPerStateByte); got != tt.want {
				t.Fatalf("addressesPerBloatTx(%d) = %d, want %d", tt.costPerStateByte, got, tt.want)
			}
		})
	}
}

// The batch must always fit the fixed gas limit once the state-creation gas of every
// fresh slot is added to the compute estimate.
func TestAddressesPerBloatTxFitsGasLimit(t *testing.T) {
	for _, cpsb := range []uint64{100, 1000, spamoor.CostPerStateByte, 5000, 20000} {
		addresses := addressesPerBloatTx(cpsb)
		stateGas := addresses * SlotsPerBloatCycle * StateCreationBytesPerSlot * cpsb
		total := addresses*PreAmsterdamGasPerAddress + stateGas + BloatTxOverheadGas

		if total > FixedGasLimitPerTx {
			t.Fatalf("cpsb %d: %d addresses need %d gas, exceeding the %d gas limit",
				cpsb, addresses, total, FixedGasLimitPerTx)
		}
	}
}

func TestBatchSizeShrinkAndGrow(t *testing.T) {
	s := &Scenario{maxAddressesPerTx: 100, addressesPerTx: 100}
	s.logger = newTestLogger()

	s.shrinkBatchSize()
	if s.addressesPerTx != 75 {
		t.Fatalf("expected batch size 75 after shrink, got %d", s.addressesPerTx)
	}

	for range batchGrowAfterRounds - 1 {
		s.growBatchSize()
	}
	if s.addressesPerTx != 75 {
		t.Fatalf("batch size grew too early: %d", s.addressesPerTx)
	}

	s.growBatchSize()
	if s.addressesPerTx != 82 {
		t.Fatalf("expected batch size 82 after growth step, got %d", s.addressesPerTx)
	}

	for range 10 * batchGrowAfterRounds {
		s.growBatchSize()
	}
	if s.addressesPerTx != 100 {
		t.Fatalf("expected batch size to settle at the maximum 100, got %d", s.addressesPerTx)
	}

	s.addressesPerTx = 1
	s.shrinkBatchSize()
	if s.addressesPerTx != 1 {
		t.Fatalf("batch size must never drop below 1, got %d", s.addressesPerTx)
	}
}
