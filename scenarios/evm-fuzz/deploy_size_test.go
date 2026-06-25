package evmfuzz

import (
	"bytes"
	"testing"
)

// newDeploySizeGen builds a generator in deploy-size mode for the given txID.
func newDeploySizeGen(txID uint64) *OpcodeGenerator {
	g := NewOpcodeGenerator(txID, "0xdeadbeef", 512, 1_000_000)
	g.SetFuzzMode("deploy-size")
	return g
}

func TestDeploySizeStructure(t *testing.T) {
	g := newDeploySizeGen(7)
	out := g.Generate()
	n := g.deploySizeN

	if n < 1 {
		t.Fatalf("deploySizeN must be >= 1, got %d", n)
	}
	if len(out) != deploySizeInitOverhead+n {
		t.Fatalf("len(out)=%d, want %d (overhead %d + N %d)", len(out), deploySizeInitOverhead+n, deploySizeInitOverhead, n)
	}

	// Fingerprint: two PUSH32 then POP POP.
	if out[0] != 0x7f || out[33] != 0x7f {
		t.Fatalf("expected PUSH32 seed/txID prologue, got 0x%02x 0x%02x", out[0], out[33])
	}
	if out[66] != 0x50 || out[67] != 0x50 {
		t.Fatalf("expected POP POP after fingerprint, got 0x%02x 0x%02x", out[66], out[67])
	}

	// Prologue must end with PUSH0 RETURN right before the filler.
	if out[deploySizeInitOverhead-2] != 0x5f || out[deploySizeInitOverhead-1] != 0xf3 {
		t.Fatalf("prologue must end PUSH0 RETURN, got 0x%02x 0x%02x", out[deploySizeInitOverhead-2], out[deploySizeInitOverhead-1])
	}

	// Both PUSH3 N operands must decode to N.
	copyLen := int(out[69])<<16 | int(out[70])<<8 | int(out[71])
	retLen := int(out[79])<<16 | int(out[80])<<8 | int(out[81])
	if copyLen != n || retLen != n {
		t.Fatalf("PUSH3 N operands decode to copy=%d ret=%d, want %d", copyLen, retLen, n)
	}

	// Filler is the deployed runtime: exactly N bytes of 0xfe.
	for i := deploySizeInitOverhead; i < len(out); i++ {
		if out[i] != 0xfe {
			t.Fatalf("filler byte %d = 0x%02x, want 0xfe", i, out[i])
		}
	}
}

func TestDeploySizeDeterministic(t *testing.T) {
	for _, txID := range []uint64{0, 1, 42, 9999} {
		a := newDeploySizeGen(txID).Generate()
		b := newDeploySizeGen(txID).Generate()
		if !bytes.Equal(a, b) {
			t.Fatalf("txID %d: non-deterministic output (%d vs %d bytes)", txID, len(a), len(b))
		}
	}
}

func TestDeploySizeBoundaryCoverage(t *testing.T) {
	want := make(map[int]bool)
	for _, b := range deploySizeBoundaries {
		want[b] = true
	}

	for txID := uint64(0); txID < 10000; txID++ {
		g := newDeploySizeGen(txID)
		g.Generate()
		delete(want, g.deploySizeN)
		if len(want) == 0 {
			return
		}
	}
	t.Fatalf("boundary sizes never emitted across 10000 txIDs: %v", want)
}
