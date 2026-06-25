package evmfuzz

import (
	"bytes"
	"testing"
)

const stateAccessSeed = "0xdeadbeefcafe"

// TestStateAccessDeterminism asserts that the same (seed,txID) produces
// byte-identical bytecode across independent generator instances.
func TestStateAccessDeterminism(t *testing.T) {
	for txID := uint64(0); txID < 16; txID++ {
		g1 := NewOpcodeGenerator(txID, stateAccessSeed, 512, 1000000)
		g1.SetFuzzMode("state-access")
		a := g1.Generate()

		g2 := NewOpcodeGenerator(txID, stateAccessSeed, 512, 1000000)
		g2.SetFuzzMode("state-access")
		b := g2.Generate()

		if !bytes.Equal(a, b) {
			t.Fatalf("txID %d: non-deterministic output (%d vs %d bytes)", txID, len(a), len(b))
		}
	}
}

// TestStateAccessBounds verifies output respects size budget and that the
// virtual stack stays accurate (net-zero state-access sequences must not drive
// stackSize negative or beyond the EVM limit during generation).
func TestStateAccessBounds(t *testing.T) {
	maxSize := 512
	for txID := uint64(0); txID < 32; txID++ {
		g := NewOpcodeGenerator(txID, stateAccessSeed, maxSize, 1000000)
		g.SetFuzzMode("state-access")
		bc := g.Generate()

		if len(bc) > maxSize {
			t.Fatalf("txID %d: bytecode %d exceeds maxSize %d", txID, len(bc), maxSize)
		}
		if g.stackSize < 0 || g.stackSize > 1024 {
			t.Fatalf("txID %d: final stackSize %d out of bounds", txID, g.stackSize)
		}
	}
}

// TestStateAccessEmitsStorageOps verifies the mode actually emits storage/access
// opcodes (SSTORE/SLOAD/EXTCODE*) beyond the seed/txID fingerprint prologue.
func TestStateAccessEmitsStorageOps(t *testing.T) {
	found := false
	for txID := uint64(0); txID < 64 && !found; txID++ {
		g := NewOpcodeGenerator(txID, stateAccessSeed, 512, 1000000)
		g.SetFuzzMode("state-access")
		bc := g.Generate()
		for _, op := range bc {
			switch op {
			case 0x55, 0x54, 0x3b, 0x31, 0x3f, 0x3c: // SSTORE/SLOAD/EXTCODESIZE/BALANCE/EXTCODEHASH/EXTCODECOPY
				found = true
			}
		}
	}
	if !found {
		t.Fatal("state-access mode never emitted any storage/account-access opcode")
	}
}

// TestStateAccessSelectorsNotInOtherModes verifies the virtual selectors
// (0x120/0x121) are gated out of the other fuzz modes' valid opcode lists.
func TestStateAccessSelectorsNotInOtherModes(t *testing.T) {
	for _, mode := range []string{"all", "opcodes", "precompiles"} {
		g := NewOpcodeGenerator(1, stateAccessSeed, 512, 1000000)
		g.SetFuzzMode(mode)
		for _, op := range g.validOpcodes {
			if g.isStateAccessOpcode(op.Opcode) {
				t.Fatalf("mode %q wrongly includes state-access selector 0x%x", mode, op.Opcode)
			}
		}
	}
}
