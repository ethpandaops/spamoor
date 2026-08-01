package erc20bloater

import (
	"testing"

	"github.com/ethpandaops/spamoor/utils"
)

// Pre-Amsterdam behavior must stay exactly what it is today: the fixed
// constants, regardless of whatever gas figures the chain reports.
func TestCalculateBatchParams_PreAmsterdamUnchanged(t *testing.T) {
	gasLimit, addresses := calculateBatchParams(false, 45_000_000, 1530, nil)
	if gasLimit != FixedGasLimitPerTx {
		t.Fatalf("expected gas limit %d, got %d", FixedGasLimitPerTx, gasLimit)
	}
	if addresses != MaxBloatedAddressesPerTx {
		t.Fatalf("expected %d addresses, got %d", MaxBloatedAddressesPerTx, addresses)
	}
}

// On Amsterdam, the batch must be sized from the real state-gas cost per
// address (2 new slots x 64 bytes x cost_per_state_byte) and the requested
// gas limit must never exceed the chain's reported ceiling.
func TestCalculateBatchParams_AmsterdamSizesFromLiveGasLimit(t *testing.T) {
	const costPerStateByte = 1530
	const maxTxGas = 45_000_000

	gasLimit, addresses := calculateBatchParams(true, maxTxGas, costPerStateByte, nil)

	statePerAddr := SlotsPerBloatCycle * uint64(64) * uint64(costPerStateByte)
	wantAddresses := (uint64(maxTxGas) - utils.MaxGasLimitPerTx) / statePerAddr
	wantGasLimit := utils.MaxGasLimitPerTx + wantAddresses*statePerAddr

	if addresses != wantAddresses {
		t.Fatalf("expected %d addresses, got %d", wantAddresses, addresses)
	}
	if gasLimit != wantGasLimit {
		t.Fatalf("expected gas limit %d, got %d", wantGasLimit, gasLimit)
	}
	if gasLimit > maxTxGas {
		t.Fatalf("computed gas limit %d exceeds the chain's reported ceiling %d", gasLimit, maxTxGas)
	}
	if addresses == 0 {
		t.Fatal("expected at least one address per tx")
	}
}

// If the chain's block gas limit doesn't even clear the legacy EIP-7825
// ceiling, there is nowhere for a state-gas reservoir to come from at all;
// this must fall back to the pre-Amsterdam sizing rather than divide by zero
// or hand back a zero-address batch that can never make progress.
func TestCalculateBatchParams_FallsBackWhenNoRoomForStateGas(t *testing.T) {
	gasLimit, addresses := calculateBatchParams(true, 10_000_000, 1530, nil)
	if gasLimit != FixedGasLimitPerTx || addresses != MaxBloatedAddressesPerTx {
		t.Fatalf("expected fallback to (%d, %d), got (%d, %d)",
			FixedGasLimitPerTx, MaxBloatedAddressesPerTx, gasLimit, addresses)
	}
}

// A tiny excess over the EIP-7825 ceiling (less than one address's worth of
// state gas) must still produce at least one address per tx, not zero -
// zero would mean the scenario silently stalls forever instead of making
// slow progress.
func TestCalculateBatchParams_NeverZeroAddresses(t *testing.T) {
	_, addresses := calculateBatchParams(true, utils.MaxGasLimitPerTx+100, 1530, nil)
	if addresses != 1 {
		t.Fatalf("expected exactly 1 address for a near-zero excess budget, got %d", addresses)
	}
}

// At an unrealistically large gas ceiling, the regular-gas need for the
// batch (which is tiny per address but scales with address count) must still
// be guarded so it never exceeds the mandatory EIP-7825 regular allotment.
func TestCalculateBatchParams_RegularGasNeverExceedsCeiling(t *testing.T) {
	const sstoreResetGas = uint64(5000)
	const loopOverhead = uint64(400)
	const callOverhead = uint64(300000)
	regularPerAddr := SlotsPerBloatCycle*sstoreResetGas + loopOverhead

	_, addresses := calculateBatchParams(true, 1_000_000_000, 1530, nil)

	if regularPerAddr*addresses+callOverhead > utils.MaxGasLimitPerTx {
		t.Fatalf("regular gas need %d for %d addresses exceeds the EIP-7825 regular allotment %d",
			regularPerAddr*addresses+callOverhead, addresses, utils.MaxGasLimitPerTx)
	}
}
