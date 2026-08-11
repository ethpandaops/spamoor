package erc20bloater

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/params"

	"github.com/ethpandaops/spamoor/spamoor"
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

// wouldBeRejectedPreAmsterdam mirrors, verbatim, the EIP-7825 check
// go-ethereum's mempool applies at core/txpool/validation.go: on a chain
// that has activated Osaka but not yet Amsterdam, any tx.Gas() above
// params.MaxTxGas is rejected before the transaction ever reaches a block.
func wouldBeRejectedPreAmsterdam(rules params.Rules, txGas uint64) bool {
	return rules.IsOsaka && !rules.IsAmsterdam && txGas > params.MaxTxGas
}

// TestCalculateBatchParams_DefaultFlagRejectedOnPreAmsterdamOsakaChain is a
// mechanical reproduction of the finding that spamoor's IsAmsterdam() is a
// static operator preference (defaults to true) rather than live on-chain
// fork detection. On a chain that is Osaka-active but has not yet activated
// Amsterdam - exactly the state of any devnet before it crosses the fork
// boundary described in the underlying issue - the batch size this scenario
// computes by default would be rejected outright by a real node, not merely
// reverted on-chain.
func TestCalculateBatchParams_DefaultFlagRejectedOnPreAmsterdamOsakaChain(t *testing.T) {
	for _, blockGasLimit := range []uint64{30_000_000, 45_000_000, 60_000_000} {
		// isAmsterdam=true is exactly what TxPool.IsAmsterdam() returns with
		// default spamoor flags; no misconfiguration is needed to reach this.
		gasLimitPerTx, addressesPerTx := calculateBatchParams(true, blockGasLimit, spamoor.CostPerStateByte, nil)

		if gasLimitPerTx <= utils.MaxGasLimitPerTx {
			t.Fatalf("block gas limit %d: expected gasLimitPerTx to exceed the EIP-7825 ceiling (%d), got %d (addresses=%d) - test assumption invalid",
				blockGasLimit, utils.MaxGasLimitPerTx, gasLimitPerTx, addressesPerTx)
		}

		// The exact rejection condition from go-ethereum's real mempool
		// validation, applied to the value this scenario would actually put
		// in the transaction's Gas field.
		preAmsterdamOsaka := params.Rules{IsOsaka: true, IsAmsterdam: false}
		if !wouldBeRejectedPreAmsterdam(preAmsterdamOsaka, gasLimitPerTx) {
			t.Fatalf("block gas limit %d: expected gasLimitPerTx=%d to be rejected on a pre-Amsterdam Osaka chain, but the real validation condition says it would be accepted",
				blockGasLimit, gasLimitPerTx)
		}

		// Control: the identical value is accepted once Amsterdam has
		// actually activated, proving the rejection is specifically about
		// the fork mismatch and not just "large numbers always fail."
		amsterdam := params.Rules{IsOsaka: true, IsAmsterdam: true}
		if wouldBeRejectedPreAmsterdam(amsterdam, gasLimitPerTx) {
			t.Fatalf("block gas limit %d: gasLimitPerTx=%d unexpectedly rejected even once the chain has activated Amsterdam",
				blockGasLimit, gasLimitPerTx)
		}

		t.Logf("CONFIRMED: block gas limit %d -> computed gasLimitPerTx=%d (%d addresses) exceeds the EIP-7825 ceiling %d; rejected pre-Amsterdam, accepted post-Amsterdam",
			blockGasLimit, gasLimitPerTx, addressesPerTx, utils.MaxGasLimitPerTx)
	}
}

// TestIsGasLimitTooHighError verifies the detection helper recognizes the
// real go-ethereum rejection text through the exact wrapping chain
// SendMultiTransactionBatch and attemptBloatRound apply to it, and does not
// false-positive on unrelated errors - including other errors that also
// mention "gas" or "limit" individually.
func TestIsGasLimitTooHighError(t *testing.T) {
	// Mirrors the real wrapping chain: core.ErrGasLimitTooHigh, wrapped by
	// the per-tx submit retry loop, wrapped by SendMultiTransactionBatch's
	// per-wallet aggregation, wrapped by attemptBloatRound.
	wrapped := fmt.Errorf("failed to send transaction batch: %w",
		fmt.Errorf("wallet 0xabc transaction 0 failed: %w",
			fmt.Errorf("failed to submit after 1 attempts: %w",
				fmt.Errorf("transaction gas limit too high (cap: %d, tx: %d)", utils.MaxGasLimitPerTx, 45_000_000))))

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil error", nil, false},
		{"exact geth error text", errors.New("transaction gas limit too high"), true},
		{"geth error with cap/tx detail", fmt.Errorf("transaction gas limit too high (cap: %d, tx: %d)", utils.MaxGasLimitPerTx, 45_000_000), true},
		{"wrapped through the real call chain", wrapped, true},
		{"unrelated context error", errors.New("context deadline exceeded"), false},
		{"unrelated error that also mentions gas", errors.New("insufficient funds for gas * price + value"), false},
		{"unrelated error that also mentions limit", errors.New("rate limit exceeded"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isGasLimitTooHighError(tt.err); got != tt.want {
				t.Fatalf("isGasLimitTooHighError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

// TestCalculateBatchParams_FallbackSizingNeverExceedsLegacyCeiling is the
// other half of the fix's proof: TestCalculateBatchParams_DefaultFlagRejectedOnPreAmsterdamOsakaChain
// shows the optimistic (Amsterdam-assumed) sizing gets rejected on a
// pre-Amsterdam Osaka chain. This shows the exact fallback values Run() uses
// once that rejection is detected - the same constants calculateBatchParams
// itself returns when isAmsterdam is false - are never subject to that
// rejection at all, at any block gas limit.
func TestCalculateBatchParams_FallbackSizingNeverExceedsLegacyCeiling(t *testing.T) {
	if FixedGasLimitPerTx > utils.MaxGasLimitPerTx {
		t.Fatalf("FixedGasLimitPerTx (%d) exceeds the EIP-7825 ceiling (%d)", FixedGasLimitPerTx, utils.MaxGasLimitPerTx)
	}

	preAmsterdamOsaka := params.Rules{IsOsaka: true, IsAmsterdam: false}
	if wouldBeRejectedPreAmsterdam(preAmsterdamOsaka, FixedGasLimitPerTx) {
		t.Fatalf("FixedGasLimitPerTx (%d) would be rejected on a pre-Amsterdam Osaka chain", FixedGasLimitPerTx)
	}

	// Same fallback values regardless of how large the (possibly stale or
	// misreported) block gas limit that triggered the fallback was.
	for _, blockGasLimit := range []uint64{10_000_000, 30_000_000, 45_000_000, 60_000_000} {
		gasLimit, addresses := calculateBatchParams(false, blockGasLimit, spamoor.CostPerStateByte, nil)
		if gasLimit != FixedGasLimitPerTx || addresses != MaxBloatedAddressesPerTx {
			t.Fatalf("block gas limit %d: expected fallback (%d, %d), got (%d, %d)",
				blockGasLimit, FixedGasLimitPerTx, MaxBloatedAddressesPerTx, gasLimit, addresses)
		}
	}
}
