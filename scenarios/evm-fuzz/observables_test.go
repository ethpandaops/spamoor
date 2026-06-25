package evmfuzz

import (
	"bytes"
	"encoding/binary"
	"testing"
)

const testSeed = "0xdeadbeefcafe"

// genWithProbes builds bytecode for a given txID with gas probes enabled.
func genWithProbes(txID uint64) []byte {
	g := NewOpcodeGenerator(txID, testSeed, 4096, 30000000)
	g.SetFuzzMode("all")
	g.SetGasProbes(true)
	return g.Generate()
}

// TestGasProbesDeterministic asserts byte-identical output for the same (seed,txID).
func TestGasProbesDeterministic(t *testing.T) {
	for txID := uint64(0); txID < 20; txID++ {
		a := genWithProbes(txID)
		b := genWithProbes(txID)
		if !bytes.Equal(a, b) {
			t.Fatalf("txID %d: probe generation not deterministic", txID)
		}
	}
}

// TestGasProbesComposeWithModes asserts probes are deterministic under every fuzz mode.
func TestGasProbesComposeWithModes(t *testing.T) {
	for _, mode := range []string{"all", "opcodes", "precompiles"} {
		gen := func() []byte {
			g := NewOpcodeGenerator(7, testSeed, 4096, 30000000)
			g.SetFuzzMode(mode)
			g.SetGasProbes(true)
			return g.Generate()
		}
		if !bytes.Equal(gen(), gen()) {
			t.Fatalf("mode %s: probe generation not deterministic", mode)
		}
	}
}

// TestGasProbesRespectBudgets asserts emitted bytecode stays within size/opcode budgets.
func TestGasProbesRespectBudgets(t *testing.T) {
	maxSize := 4096
	g := NewOpcodeGenerator(3, testSeed, maxSize, 30000000)
	g.SetFuzzMode("all")
	g.SetGasProbes(true)
	code := g.Generate()

	if len(code) > maxSize {
		t.Fatalf("bytecode size %d exceeds maxSize %d", len(code), maxSize)
	}
	if got := g.countOpcodesInBytecode(code); got > g.maxOpcodeCount {
		t.Fatalf("opcode count %d exceeds maxOpcodeCount %d", got, g.maxOpcodeCount)
	}
}

// countProbeTopics scans for PUSH32 topic words matching txID + a known probe kind.
func countProbeTopics(code []byte, txID uint64) (checkpoints, deltas int) {
	prefix := make([]byte, 8)
	binary.BigEndian.PutUint64(prefix, txID)
	for i := 0; i+33 <= len(code); i++ {
		if code[i] != 0x7f { // PUSH32
			continue
		}
		w := code[i+1 : i+33]
		if !bytes.Equal(w[0:8], prefix) {
			continue
		}
		switch w[8] {
		case probeKindCheckpoint:
			checkpoints++
		case probeKindDelta:
			deltas++
		}
	}
	return
}

// TestGasProbesActuallyEmitted asserts probes appear when enabled and never when disabled.
func TestGasProbesActuallyEmitted(t *testing.T) {
	// Disabled: no probe topics should ever appear.
	for txID := uint64(0); txID < 20; txID++ {
		g := NewOpcodeGenerator(txID, testSeed, 4096, 30000000)
		g.SetFuzzMode("all")
		g.SetGasProbes(false)
		code := g.Generate()
		if c, d := countProbeTopics(code, txID); c+d != 0 {
			t.Fatalf("txID %d: probes emitted while disabled (cp=%d delta=%d)", txID, c, d)
		}
	}

	// Enabled: across a range of txIDs at least some probes of each kind must appear.
	totalCp, totalDelta := 0, 0
	for txID := uint64(0); txID < 50; txID++ {
		code := genWithProbes(txID)
		cp, d := countProbeTopics(code, txID)
		totalCp += cp
		totalDelta += d
	}
	if totalCp == 0 {
		t.Fatalf("no checkpoint probes emitted across 50 txIDs")
	}
	if totalDelta == 0 {
		t.Fatalf("no delta probes emitted across 50 txIDs")
	}
}
