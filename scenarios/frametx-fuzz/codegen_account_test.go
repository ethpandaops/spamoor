package frametxfuzz

import (
	"testing"

	"github.com/ethereum/go-ethereum/core/vm/runtime"
)

// TestAccountCodeDispatch checks that an account contract's calldata selects between the
// role path and the sweep path.
//
// APPROVE (0xaa) is undefined in a bare EVM, so it is a marker: reaching it faults, and
// skipping it returns cleanly. A non-zero scope must run the role path and hit APPROVE; a
// zero scope must jump over it to the sweep tail and return.
func TestAccountCodeDispatch(t *testing.T) {
	code := GenerateAccountCode("0x2", 2, 0, 400_000, true)

	// The patched dispatch must point at a JUMPDEST.
	sweepPC := int(code[dispatchSweepPCOffset])<<8 | int(code[dispatchSweepPCOffset+1])
	if sweepPC >= len(code) || code[sweepPC] != 0x5b {
		t.Fatalf("sweep offset %d does not land on a JUMPDEST", sweepPC)
	}

	scope := make([]byte, 32)
	scope[31] = 3

	if _, _, err := runtime.Execute(code, scope, &runtime.Config{GasLimit: 5_000_000}); err == nil {
		t.Error("a non-zero scope returned cleanly; expected it to reach the undefined APPROVE")
	}

	if _, _, err := runtime.Execute(code, []byte{}, &runtime.Config{GasLimit: 5_000_000}); err != nil {
		t.Errorf("a zero scope faulted (%v); expected it to sweep and return", err)
	}
}
