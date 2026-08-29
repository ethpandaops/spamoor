package frametxfuzz

import (
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/core/vm/runtime"
	"github.com/ethpandaops/spamoor/txtypes"
)

// runScript executes a script against the compiled probe contract in a bare EVM.
//
// No devnet implements EIP-8141 yet, so the script interpreter is checked here instead:
// everything except the frame instructions themselves is ordinary EVM, and a bug in the
// dispatch loop or the argument offsets would otherwise only show up as an unexplained
// frame failure on a devnet.
func runScript(t *testing.T, script *ProbeScript) (*runtime.Config, error) {
	t.Helper()

	code, err := ProbeRuntimeCode()
	if err != nil {
		t.Fatalf("failed to compile probe contract: %v", err)
	}

	cfg := &runtime.Config{GasLimit: 10_000_000}

	_, state, err := runtime.Execute(code, script.Bytes(), cfg)
	cfg.State = state

	return cfg, err
}

func TestScriptEmptyReturnsSuccessfully(t *testing.T) {
	if _, err := runScript(t, NewProbeScript()); err != nil {
		t.Fatalf("empty script failed: %v", err)
	}
}

func TestScriptStop(t *testing.T) {
	// Everything after an explicit stop must be unreachable, so the revert that follows
	// it does not fire.
	if _, err := runScript(t, NewProbeScript().Stop().Revert()); err != nil {
		t.Fatalf("stop did not end the script: %v", err)
	}
}

func TestScriptRevert(t *testing.T) {
	_, err := runScript(t, NewProbeScript().Revert())
	if !errors.Is(err, vm.ErrExecutionReverted) {
		t.Fatalf("revert produced %v, want an execution revert", err)
	}
}

func TestScriptUnknownSelectorReverts(t *testing.T) {
	// An operation the contract does not implement must fail loudly rather than be
	// skipped, or a caller would read a passing frame as a passing assertion.
	script := NewProbeScript()
	script.add(record{op: 0x7f, name: "unknown"})

	if _, err := runScript(t, script); !errors.Is(err, vm.ErrExecutionReverted) {
		t.Fatalf("unknown selector produced %v, want an execution revert", err)
	}
}

func TestScriptSStore(t *testing.T) {
	slot := common.HexToHash("0x11")
	value := common.HexToHash("0x2222")

	cfg, err := runScript(t, NewProbeScript().SStore(slot, value))
	if err != nil {
		t.Fatalf("sstore failed: %v", err)
	}

	contract := common.BytesToAddress([]byte("contract"))
	if got := cfg.State.GetState(contract, slot); got != value {
		t.Errorf("slot holds %s, want %s", got, value)
	}
}

func TestScriptReadSLoad(t *testing.T) {
	slot := common.HexToHash("0x11")
	value := common.HexToHash("0x2222")

	if _, err := runScript(t, NewProbeScript().SStore(slot, value).ReadSLoad(slot)); err != nil {
		t.Fatalf("read after write failed: %v", err)
	}
}

func TestScriptLog(t *testing.T) {
	topic := common.HexToHash("0xabcd")

	cfg, err := runScript(t, NewProbeScript().Log(topic, 64))
	if err != nil {
		t.Fatalf("log failed: %v", err)
	}

	logs := cfg.State.Logs()
	if len(logs) != 1 {
		t.Fatalf("emitted %d logs, want 1", len(logs))
	}

	if len(logs[0].Topics) != 1 || logs[0].Topics[0] != topic {
		t.Errorf("log topics are %v, want [%s]", logs[0].Topics, topic)
	}

	if len(logs[0].Data) != 64 {
		t.Errorf("log carries %d data bytes, want 64", len(logs[0].Data))
	}
}

func TestScriptBurn(t *testing.T) {
	// Burning must consume gas without failing, and must stop at the target rather than
	// running out.
	if _, err := runScript(t, NewProbeScript().Burn(1_000_000)); err != nil {
		t.Fatalf("burn failed: %v", err)
	}
}

func TestScriptMultipleRecords(t *testing.T) {
	// Every record is the same width, so a bad cursor advance shows up as soon as more
	// than one operation is present.
	script := NewProbeScript().
		SStore(common.HexToHash("0x01"), common.HexToHash("0xaa")).
		SStore(common.HexToHash("0x02"), common.HexToHash("0xbb")).
		ReadSLoad(common.HexToHash("0x01")).
		ReadSLoad(common.HexToHash("0x02"))

	cfg, err := runScript(t, script)
	if err != nil {
		t.Fatalf("multi-record script failed: %v", err)
	}

	contract := common.BytesToAddress([]byte("contract"))
	if got := cfg.State.GetState(contract, common.HexToHash("0x02")); got != common.HexToHash("0xbb") {
		t.Errorf("second slot holds %s, want 0xbb", got)
	}
}

func TestScriptFrameInstructionsAreReached(t *testing.T) {
	// The frame instructions do not exist outside a frame transaction, so on a plain EVM
	// they must halt as undefined. That proves the dispatch reaches them: a script that
	// silently did nothing would return successfully.
	for _, tc := range []struct {
		name   string
		script *ProbeScript
	}{
		{"txparam", NewProbeScript().ReadTxParam(TxParamFrameCount)},
		{"frameparam", NewProbeScript().ReadFrameParam(FrameParamMode, 0)},
		{"sigparam", NewProbeScript().ReadSigParam(SigParamScheme, 0)},
		{"framedata", NewProbeScript().ReadFrameData(0, 0)},
		{"sigdata", NewProbeScript().ReadSigData(0, 0)},
		{"rootref", NewProbeScript().ReadRootRef(0, RecentRootFieldRoot)},
		{"approve", NewProbeScript().Approve(txtypes.ApproveExecutionAndPayment)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := runScript(t, tc.script); err == nil {
				t.Error("frame instruction executed on a plain EVM, want an undefined-opcode halt")
			}
		})
	}
}

func TestScriptEncoding(t *testing.T) {
	script := NewProbeScript().SStore(common.HexToHash("0x01"), common.HexToHash("0x02"))

	encoded := script.Bytes()
	if len(encoded) != ProbeRecordSize {
		t.Fatalf("one operation encoded to %d bytes, want %d", len(encoded), ProbeRecordSize)
	}

	if encoded[31] != opSStore {
		t.Errorf("selector byte is 0x%02x, want 0x%02x", encoded[31], opSStore)
	}

	if common.BytesToHash(encoded[32:64]) != common.HexToHash("0x01") {
		t.Error("first argument is not in the second word")
	}

	if common.BytesToHash(encoded[64:96]) != common.HexToHash("0x02") {
		t.Error("second argument is not in the third word")
	}
}

func TestScriptPrefixSafety(t *testing.T) {
	safe := NewProbeScript().Approve(txtypes.ApproveExecutionAndPayment)
	if !safe.PrefixSafe() {
		t.Error("approve must be allowed in the validation prefix")
	}

	if err := safe.Validate(true); err != nil {
		t.Errorf("prefix-safe script rejected: %v", err)
	}

	// SSTORE, GAS and calls are all banned during validation-prefix execution.
	for name, script := range map[string]*ProbeScript{
		"sstore":     NewProbeScript().SStore(common.Hash{}, common.Hash{}),
		"burn":       NewProbeScript().Burn(1000),
		"call":       NewProbeScript().Call(common.Address{}, 1000),
		"write_root": NewProbeScript().WriteRoot(common.Hash{}, common.Hash{}),
	} {
		if script.PrefixSafe() {
			t.Errorf("%s must not be allowed in the validation prefix", name)
		}

		if err := script.Validate(true); !errors.Is(err, ErrNotPrefixSafe) {
			t.Errorf("%s validated with %v, want ErrNotPrefixSafe", name, err)
		}

		if err := script.Validate(false); err != nil {
			t.Errorf("%s rejected outside the prefix: %v", name, err)
		}
	}
}
