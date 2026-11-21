package bloatedeoa

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm/runtime"
)

// deployAttack compiles the attack contract for the given options and installs
// it as the code of a fresh account in a new statedb.
func deployAttack(t *testing.T, options ScenarioOptions) (*state.StateDB, common.Address, []byte) {
	t.Helper()

	options.LoopGasThreshold = ScenarioDefaultOptions.LoopGasThreshold
	code, err := compileAttackContract(&options)
	if err != nil {
		t.Fatal(err)
	}

	statedb, err := state.New(types.EmptyRootHash, state.NewDatabaseForTesting())
	if err != nil {
		t.Fatal(err)
	}

	addr := common.BytesToAddress([]byte("attack"))
	statedb.CreateAccount(addr)
	statedb.SetCode(addr, code, tracing.CodeChangeUnspecified)

	return statedb, addr, code
}

func slot(statedb *state.StateDB, addr common.Address, key int64) int64 {
	return statedb.GetState(addr, common.BigToHash(big.NewInt(key))).Big().Int64()
}

// TestSloadAttack verifies the sload runtime (test_sload_bloated): the counter
// lives in slot 0, calls without calldata run a gas-bounded loop, and the
// counter advances and continues across calls.
func TestSloadAttack(t *testing.T) {
	const startSlot = 100

	statedb, addr, _ := deployAttack(t, ScenarioOptions{Mode: ModeSload})
	cfg := &runtime.Config{GasLimit: 2_000_000, State: statedb}

	// pre-state: s[i] = i (SLOAD does not change it)
	for i := int64(1); i < 100000; i++ {
		statedb.SetState(addr, common.BigToHash(big.NewInt(i)), common.BigToHash(big.NewInt(i)))
	}

	// delegation tx: set the counter via calldata
	if _, _, err := runtime.Call(addr, common.BigToHash(big.NewInt(startSlot)).Bytes(), cfg); err != nil {
		t.Fatalf("counter setup failed: %v", err)
	}
	if got := slot(statedb, addr, 0); got != startSlot {
		t.Fatalf("counter = %d, want %d", got, startSlot)
	}

	// first attack call advances the counter well past the start
	if _, _, err := runtime.Call(addr, nil, cfg); err != nil {
		t.Fatalf("attack call failed: %v", err)
	}
	end := slot(statedb, addr, 0)
	if end <= startSlot+10 {
		t.Fatalf("counter did not advance: %d", end)
	}

	// second call continues from where the first one stopped
	if _, _, err := runtime.Call(addr, nil, cfg); err != nil {
		t.Fatalf("second attack call failed: %v", err)
	}
	if next := slot(statedb, addr, 0); next <= end {
		t.Fatalf("counter did not continue: %d <= %d", next, end)
	}

	t.Logf("slots accessed per 2M gas call: %d", end-startSlot)
}

// TestSstoreAttack verifies the sstore runtime (test_sstore_bloated): the
// [start, end) range comes from calldata, each variant writes the EELS value,
// slot 0 is never touched, and slots outside the range are untouched.
func TestSstoreAttack(t *testing.T) {
	const (
		prefill = 1000 // distinct pre-state so every write is observable
		start   = 200
		count   = 150
		end     = start + count
	)

	tests := []struct {
		name    string
		options ScenarioOptions
		expect  func(key int64) int64
	}{
		{"existing", ScenarioOptions{Mode: ModeSstore, ExistingSlots: true}, func(key int64) int64 { return key }},
		{"existing-new", ScenarioOptions{Mode: ModeSstore, ExistingSlots: true, WriteNewValue: true}, func(key int64) int64 { return key + 1 }},
		{"empty-new", ScenarioOptions{Mode: ModeSstore, ExistingSlots: false, WriteNewValue: true}, func(key int64) int64 { return key + 1 }},
		{"empty", ScenarioOptions{Mode: ModeSstore, ExistingSlots: false}, func(key int64) int64 { return 0 }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			statedb, addr, _ := deployAttack(t, test.options)
			cfg := &runtime.Config{GasLimit: 5_000_000, State: statedb}

			// pre-state: s[i] = i + prefill for the whole span, so a correct
			// write to any value (including 0) is observable.
			for i := int64(start - 1); i <= end; i++ {
				statedb.SetState(addr, common.BigToHash(big.NewInt(i)), common.BigToHash(big.NewInt(i+prefill)))
			}

			calldata := append(common.BigToHash(big.NewInt(start)).Bytes(), common.BigToHash(big.NewInt(end)).Bytes()...)
			if _, _, err := runtime.Call(addr, calldata, cfg); err != nil {
				t.Fatalf("attack call failed: %v", err)
			}

			// slot 0 must never be touched (no on-chain counter)
			if got := slot(statedb, addr, 0); got != 0 {
				t.Fatalf("slot 0 was written: %d", got)
			}

			// the range [start, end) holds the expected value
			for key := int64(start); key < end; key++ {
				if got := slot(statedb, addr, key); got != test.expect(key) {
					t.Fatalf("s[%d] = %d, want %d", key, got, test.expect(key))
				}
			}

			// boundaries stay at their pre-state
			for _, key := range []int64{start - 1, end} {
				if got := slot(statedb, addr, key); got != key+prefill {
					t.Fatalf("out-of-range slot %d was touched: %d", key, got)
				}
			}
		})
	}
}

// TestSstoreDelegationNoop verifies the calldata guard: an sstore contract
// called with empty calldata (the delegation tx) touches no storage.
func TestSstoreDelegationNoop(t *testing.T) {
	statedb, addr, _ := deployAttack(t, ScenarioOptions{Mode: ModeSstore, ExistingSlots: true})
	cfg := &runtime.Config{GasLimit: 1_000_000, State: statedb}

	if _, _, err := runtime.Call(addr, nil, cfg); err != nil {
		t.Fatalf("no-op delegation call failed: %v", err)
	}

	for _, key := range []int64{0, 1, 2} {
		if got := slot(statedb, addr, key); got != 0 {
			t.Fatalf("empty-calldata call wrote s[%d] = %d", key, got)
		}
	}
}

func TestSstoreRangesDisjoint(t *testing.T) {
	s := &Scenario{startCounter: big.NewInt(1), slotsPerTx: 150}

	var prevEnd *big.Int
	for txIdx := uint64(0); txIdx < 5; txIdx++ {
		start, end := s.sstoreRange(txIdx)

		if new(big.Int).Sub(end, start).Uint64() != s.slotsPerTx {
			t.Fatalf("tx %d range width = %v, want %d", txIdx, new(big.Int).Sub(end, start), s.slotsPerTx)
		}
		if prevEnd != nil && start.Cmp(prevEnd) != 0 {
			t.Fatalf("tx %d start %v does not continue from previous end %v", txIdx, start, prevEnd)
		}
		prevEnd = end
	}
}

func TestComputeSlotsPerTxOverride(t *testing.T) {
	s := &Scenario{options: ScenarioOptions{SlotsPerTx: 42}}
	if got := s.computeSlotsPerTx(); got != 42 {
		t.Fatalf("slots per tx = %d, want 42", got)
	}
}

func TestNonExistingStartSlot(t *testing.T) {
	want := common.HexToHash("0xf3cf193bb4af1022af7d2089f37d8bae7157b85f").Big()
	if nonExistingStartSlot.Cmp(want) != 0 {
		t.Fatalf("start slot = 0x%x, want 0x%x", nonExistingStartSlot, want)
	}
}

func TestLoopGasThresholdOverride(t *testing.T) {
	code, err := compileAttackContract(&ScenarioOptions{Mode: ModeSload, LoopGasThreshold: 0x123456})
	if err != nil {
		t.Fatal(err)
	}

	// PUSH3 0x123456
	if !bytes.Contains(code, []byte{0x62, 0x12, 0x34, 0x56}) {
		t.Fatalf("threshold override not applied: 0x%x", code)
	}
}

func TestInitcodeDeploysRuntime(t *testing.T) {
	initcode, err := compileInitcode()
	if err != nil {
		t.Fatal(err)
	}

	runtimeCode, err := compileAttackContract(&ScenarioOptions{Mode: ModeSstore, LoopGasThreshold: 0xffff})
	if err != nil {
		t.Fatal(err)
	}

	deployed, _, _, err := runtime.Create(append(initcode, runtimeCode...), &runtime.Config{GasLimit: 1_000_000})
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(deployed, runtimeCode) {
		t.Fatalf("deployed code 0x%x, want 0x%x", deployed, runtimeCode)
	}
}
