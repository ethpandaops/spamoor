package frametxfuzz

import (
	"bytes"
	"testing"

	"github.com/ethpandaops/spamoor/txtypes"
)

// The legacy fee flag selects bogota-devnet-1's EIP-7906 table everywhere the fuzzer
// encodes it: the probe contract, the generated contracts and the assertion reads.
func TestAssertTableFollowsTheFeeRevision(t *testing.T) {
	for _, legacy := range []bool{false, true} {
		previous := txtypes.LegacyFrameFees
		txtypes.LegacyFrameFees = legacy

		ops := assertOps()

		runtime, err := compileProbeRuntime(ops)
		if err != nil {
			t.Fatalf("legacy=%v: %v", legacy, err)
		}

		instructions := instructionsIn(runtime)
		if !bytes.Contains(instructions, []byte{ops.eventDataCopy}) {
			t.Errorf("legacy=%v: probe has no EVENTDATACOPY 0x%02x", legacy, ops.eventDataCopy)
		}

		if hasB8 := bytes.Contains(instructions, []byte{0xb8}); hasB8 != legacy {
			t.Errorf("legacy=%v: probe executes 0xb8 = %v", legacy, hasB8)
		}

		if want := map[bool]byte{false: 0xb7, true: 0xb8}[legacy]; ops.eventDataCopy != want {
			t.Errorf("legacy=%v: EVENTDATACOPY = 0x%02x, want 0x%02x", legacy, ops.eventDataCopy, want)
		}

		hasTxDiff := false
		for _, def := range frameOpcodeDefinitions(nil) {
			if def.Name == "TXDIFF" {
				hasTxDiff = def.Opcode == 0xb7
			}
		}

		if hasTxDiff != legacy {
			t.Errorf("legacy=%v: generator TXDIFF present = %v", legacy, hasTxDiff)
		}

		script := NewProbeScript()
		appendAssertReads(script, &Recipe{}, [20]byte{})

		readsTxDiff := false
		for _, record := range script.records {
			readsTxDiff = readsTxDiff || record.op == opReadTxDiff
		}

		if readsTxDiff != legacy {
			t.Errorf("legacy=%v: assertion sweep reads TXDIFF = %v", legacy, readsTxDiff)
		}

		txtypes.LegacyFrameFees = previous
	}
}
