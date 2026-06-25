package evmfuzz

import (
	"bytes"
	"testing"
)

// TestSystemModeDeterministic asserts byte-identical output for a given (seed, txID).
func TestSystemModeDeterministic(t *testing.T) {
	const seed = "0xdeadbeef"
	for txID := uint64(0); txID < 20; txID++ {
		g1 := NewOpcodeGenerator(txID, seed, 512, 1_000_000)
		g1.SetFuzzMode("system")
		out1 := g1.Generate()

		g2 := NewOpcodeGenerator(txID, seed, 512, 1_000_000)
		g2.SetFuzzMode("system")
		out2 := g2.Generate()

		if !bytes.Equal(out1, out2) {
			t.Fatalf("txID %d: non-deterministic output:\n%x\n%x", txID, out1, out2)
		}
	}
}

// TestSystemModeBudgets asserts emitted bytecode stays within the size budget.
func TestSystemModeBudgets(t *testing.T) {
	const maxSize = 512
	for txID := uint64(0); txID < 50; txID++ {
		g := NewOpcodeGenerator(txID, "0x01", maxSize, 1_000_000)
		g.SetFuzzMode("system")
		out := g.Generate()
		if len(out) > maxSize {
			t.Fatalf("txID %d: bytecode %d exceeds max size %d", txID, len(out), maxSize)
		}
		if len(out) == 0 {
			t.Fatalf("txID %d: empty bytecode", txID)
		}
	}
}

// TestSystemModeOnlySystemOpcodes asserts buildValidOpcodeList gates to system ops only.
func TestSystemModeOnlySystemOpcodes(t *testing.T) {
	g := NewOpcodeGenerator(7, "0x02", 512, 1_000_000)
	g.SetFuzzMode("system")
	if len(g.validOpcodes) != 3 {
		t.Fatalf("expected exactly 3 system opcodes, got %d", len(g.validOpcodes))
	}
	for _, op := range g.validOpcodes {
		if !g.isSystemOpcode(op.Opcode) {
			t.Fatalf("non-system opcode %s (0x%x) present in system mode", op.Name, op.Opcode)
		}
	}
}

// TestSystemModeEmitsContractCalls asserts system mode actually emits the predeploy
// addresses (PUSH20 of each system contract) and CALL opcodes.
func TestSystemModeEmitsContractCalls(t *testing.T) {
	seenDeposit, seenExit, seenFactory := false, false, false
	for txID := uint64(0); txID < 100; txID++ {
		g := NewOpcodeGenerator(txID, "0x03", 512, 2_000_000)
		g.SetFuzzMode("system")
		out := g.Generate()
		if bytes.Contains(out, builderDepositAddr[:]) {
			seenDeposit = true
		}
		if bytes.Contains(out, builderExitAddr[:]) {
			seenExit = true
		}
		if bytes.Contains(out, create2FactoryAddr[:]) {
			seenFactory = true
		}
	}
	if !seenDeposit || !seenExit || !seenFactory {
		t.Fatalf("missing system-contract calls: deposit=%v exit=%v factory=%v", seenDeposit, seenExit, seenFactory)
	}
}

// TestSystemHandlersStackEffect asserts each handler's net stack effect matches its
// declared StackOutput (1) in "all" mode, keeping the generator's stack model honest.
func TestSystemHandlersStackEffect(t *testing.T) {
	cases := []struct {
		name string
		gen  func(*OpcodeGenerator) []byte
	}{
		{"deposit", (*OpcodeGenerator).generateBuilderDepositCall},
		{"exit", (*OpcodeGenerator).generateBuilderExitCall},
		{"factory", (*OpcodeGenerator).generateCreate2FactoryCall},
	}
	for _, c := range cases {
		for txID := uint64(0); txID < 40; txID++ {
			g := NewOpcodeGenerator(txID, "0x04", 512, 2_000_000)
			g.SetFuzzMode("all") // MLOAD branch: result pushed once
			bc := c.gen(g)
			if got := netStackEffect(bc); got != 1 {
				t.Fatalf("%s txID %d: net stack effect %d, want 1\n%x", c.name, txID, got, bc)
			}
		}
	}
}

// netStackEffect computes the net stack delta of a self-contained bytecode sequence
// covering exactly the opcodes the system handlers emit.
func netStackEffect(bc []byte) int {
	stack := 0
	pc := 0
	for pc < len(bc) {
		op := bc[pc]
		switch {
		case op == 0x5f: // PUSH0
			stack++
			pc++
		case op >= 0x60 && op <= 0x7f: // PUSH1-PUSH32
			stack++
			pc += int(op-0x5f) + 1
		case op == 0x52: // MSTORE
			stack -= 2
			pc++
		case op == 0x51: // MLOAD
			pc++ // consumes 1 produces 1: net 0
		case op == 0x50: // POP
			stack--
			pc++
		case op == 0x34: // CALLVALUE
			stack++
			pc++
		case op == 0x5a: // GAS
			stack++
			pc++
		case op == 0xf1: // CALL
			stack -= 6 // 7 in, 1 out
			pc++
		case op == 0xa0: // LOG0
			stack -= 2
			pc++
		default:
			pc++
		}
	}
	return stack
}
