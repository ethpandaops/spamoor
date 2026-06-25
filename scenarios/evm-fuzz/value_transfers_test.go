package evmfuzz

import (
	"bytes"
	"testing"
)

// TestTransfersModeDeterministic asserts that the transfers mode produces
// byte-identical bytecode for the same (seed, txID) across runs.
func TestTransfersModeDeterministic(t *testing.T) {
	const seed = "0xdeadbeef"
	for txID := uint64(0); txID < 25; txID++ {
		g1 := NewOpcodeGenerator(txID, seed, 512, 1000000)
		g1.SetFuzzMode("transfers")
		out1 := g1.Generate()

		g2 := NewOpcodeGenerator(txID, seed, 512, 1000000)
		g2.SetFuzzMode("transfers")
		out2 := g2.Generate()

		if !bytes.Equal(out1, out2) {
			t.Fatalf("txID %d: non-deterministic output: %x != %x", txID, out1, out2)
		}
		if len(out1) > 512 {
			t.Fatalf("txID %d: bytecode exceeds maxSize: %d > 512", txID, len(out1))
		}
		if len(out1) == 0 {
			t.Fatalf("txID %d: empty bytecode", txID)
		}
	}
}

// TestTransfersGeneratorsStackBalance verifies every transfer generator nets
// exactly its declared StackOutput when run in isolation. A single off-by-one
// here desyncs the whole generator.
func TestTransfersGeneratorsStackBalance(t *testing.T) {
	cases := []struct {
		name   string
		opcode uint16
		out    int
		fn     func(g *OpcodeGenerator) []byte
	}{
		{"XFER_CALL", 0x120, 1, (*OpcodeGenerator).generateValueCall},
		{"XFER_CALLCODE", 0x121, 1, (*OpcodeGenerator).generateValueCallcode},
		{"XFER_CREATE", 0x122, 1, (*OpcodeGenerator).generateValueCreate},
		{"XFER_CREATE2", 0x123, 1, (*OpcodeGenerator).generateValueCreate2},
		{"XFER_REVERTING", 0x124, 1, (*OpcodeGenerator).generateRevertingTransfer},
		{"SELFDESTRUCT_SWEEP", 0x125, 0, (*OpcodeGenerator).generateSelfdestructSweep},
	}

	for _, tc := range cases {
		// Run across several seeds to cover all internal random branches.
		for s := 0; s < 50; s++ {
			g := NewOpcodeGenerator(uint64(s), "0x01", 512, 1000000)
			bc := tc.fn(g)
			if len(bc) == 0 {
				t.Fatalf("%s seed %d: emitted no bytecode", tc.name, s)
			}
			got := simulateStackDelta(bc)
			if got != tc.out {
				t.Fatalf("%s seed %d: stack delta %d, want %d (bytecode %x)", tc.name, s, got, tc.out, bc)
			}
		}

		// Verify the op is registered with the declared StackOutput.
		info := NewOpcodeGenerator(0, "0x01", 512, 1000000).opcodeInfos[tc.opcode]
		if info == nil {
			t.Fatalf("%s: opcode 0x%x not registered", tc.name, tc.opcode)
		}
		if info.StackInput != 0 || info.StackOutput != tc.out {
			t.Fatalf("%s: registered StackInput=%d StackOutput=%d, want 0/%d", tc.name, info.StackInput, info.StackOutput, tc.out)
		}
	}
}

// TestTransfersModeGating ensures transfer ops appear only in transfers mode
// and never leak into the other modes.
func TestTransfersModeGating(t *testing.T) {
	hasTransfer := func(mode string) bool {
		g := NewOpcodeGenerator(0, "0x01", 512, 1000000)
		g.SetFuzzMode(mode)
		for _, op := range g.validOpcodes {
			if g.isTransferOpcode(op.Opcode) {
				return true
			}
		}
		return false
	}

	if !hasTransfer("transfers") {
		t.Fatal("transfers mode missing transfer generators")
	}
	for _, mode := range []string{"all", "opcodes", "precompiles"} {
		if hasTransfer(mode) {
			t.Fatalf("mode %q must not contain transfer generators", mode)
		}
	}
}

// simulateStackDelta computes the net stack change of a self-contained bytecode
// blob using accurate per-opcode push/pop semantics for the opcodes our
// generators emit. SELFDESTRUCT/STOP/RETURN/REVERT halt the frame; we stop
// counting there since nothing after executes (the virtual model nets there).
func simulateStackDelta(bc []byte) int {
	delta := 0
	pc := 0
	for pc < len(bc) {
		op := bc[pc]
		switch {
		case op == 0x5f: // PUSH0
			delta++
			pc++
		case op >= 0x60 && op <= 0x7f: // PUSH1..PUSH32
			delta++
			pc += 1 + int(op-0x5f)
		case op >= 0x80 && op <= 0x8f: // DUP1..DUP16
			delta++
			pc++
		case op >= 0x90 && op <= 0x9f: // SWAP1..SWAP16
			pc++
		case op == 0x30 || op == 0x32 || op == 0x33 || op == 0x34: // ADDRESS/ORIGIN/CALLER/CALLVALUE
			delta++
			pc++
		case op == 0x50: // POP
			delta--
			pc++
		case op == 0x52: // MSTORE
			delta -= 2
			pc++
		case op == 0xf0: // CREATE: -3 +1
			delta -= 2
			pc++
		case op == 0xf5: // CREATE2: -4 +1
			delta -= 3
			pc++
		case op == 0xf1 || op == 0xf2: // CALL/CALLCODE: -7 +1
			delta -= 6
			pc++
		case op == 0xff: // SELFDESTRUCT: -1, halts
			delta--
			return delta
		default:
			pc++
		}
	}
	return delta
}
