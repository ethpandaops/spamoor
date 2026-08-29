package frametxfuzz

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"

	"github.com/ethpandaops/spamoor/txtypes"
)

// ProbeScript selectors, matching the dispatch table in contract.geas.
const (
	opStop           = 0x01
	opRevert         = 0x02
	opSStore         = 0x03
	opLog            = 0x04
	opBurn           = 0x05
	opApprove        = 0x06
	opReadTxParam    = 0x07
	opReadFrameParam = 0x08
	opReadSigParam   = 0x09
	opReadFrameData  = 0x0a
	opReadSigData    = 0x0b
	opReadRootRef    = 0x0c
	opCall           = 0x0d
	opWriteRoot      = 0x0e
	opStoreTxParam   = 0x0f
	opReadSLoad      = 0x10
)

// ProbeRecordSize is the width of one script record: the selector and three arguments, one
// 32-byte word each.
const ProbeRecordSize = 128

// probeEntryGas covers a frame's EIP-2929 charge for accessing the probe contract plus
// the interpreter's fixed overhead, which every script pays before its first operation.
const probeEntryGas = 5_000

// interpreterGasPerRecord covers the loop's bounds check, the calldata loads and the
// dispatch chain. It is an over-estimate for early selectors and about right for late
// ones, which is the direction a frame's gas budget should err in.
const interpreterGasPerRecord = 250

// ErrNotPrefixSafe is returned when a script would run an operation that EIP-8141's
// validation trace rules forbid inside the validation prefix.
var ErrNotPrefixSafe = fmt.Errorf("script operation is not allowed in the validation prefix")

// record is one encoded script operation together with what it costs and where it may
// run.
type record struct {
	op         byte
	a, b, c    common.Hash
	name       string
	prefixSafe bool
	execGas    uint64
	stateGas   uint64
}

// ProbeScript is a program for the probe contract, built op by op.
//
// The builder is chainable and never fails: an operation that cannot be expressed is
// rejected by Validate rather than by a constructor, so a caller assembling a script
// from generated parameters has one place to check.
type ProbeScript struct {
	records []record
}

// NewProbeScript returns an empty script.
func NewProbeScript() *ProbeScript { return &ProbeScript{} }

// add appends a record.
func (s *ProbeScript) add(r record) *ProbeScript {
	s.records = append(s.records, r)

	return s
}

// Bytes encodes the script as the contract's calldata.
func (s *ProbeScript) Bytes() []byte {
	data := make([]byte, 0, len(s.records)*ProbeRecordSize)

	for _, r := range s.records {
		var selector common.Hash
		selector[31] = r.op

		data = append(data, selector.Bytes()...)
		data = append(data, r.a.Bytes()...)
		data = append(data, r.b.Bytes()...)
		data = append(data, r.c.Bytes()...)
	}

	return data
}

// Len returns the number of operations.
func (s *ProbeScript) Len() int { return len(s.records) }

// PrefixSafe reports whether every operation may run inside a validation prefix frame.
func (s *ProbeScript) PrefixSafe() bool {
	for _, r := range s.records {
		if !r.prefixSafe {
			return false
		}
	}

	return true
}

// Validate checks the script against where it is going to run.
//
// EIP-8141 bans a list of opcodes during validation-prefix execution and restricts
// storage access to the sender's own, so a script that would be refused by a public
// mempool node is caught here rather than as an opaque rejection. Frames outside the
// prefix have no such restriction.
func (s *ProbeScript) Validate(inPrefix bool) error {
	if !inPrefix {
		return nil
	}

	for i, r := range s.records {
		if !r.prefixSafe {
			return fmt.Errorf("%w: operation %d (%s)", ErrNotPrefixSafe, i, r.name)
		}
	}

	return nil
}

// ExecutionGas returns a conservative execution gas estimate for the script, which a
// caller uses to size the frame's limits.execution.
func (s *ProbeScript) ExecutionGas() uint64 {
	total := uint64(0)

	for _, r := range s.records {
		total += interpreterGasPerRecord + r.execGas
	}

	return total
}

// StateGas returns a conservative EIP-8037 state gas estimate for the script.
func (s *ProbeScript) StateGas() uint64 {
	total := uint64(0)

	for _, r := range s.records {
		total += r.stateGas
	}

	return total
}

// Stop ends the script successfully. It is implicit at the end of the calldata; an
// explicit stop is only needed to cut a script short.
func (s *ProbeScript) Stop() *ProbeScript {
	return s.add(record{op: opStop, name: "stop", prefixSafe: true})
}

// Revert fails the frame deterministically, which is how a shape produces a failed
// frame without depending on a gas boundary.
func (s *ProbeScript) Revert() *ProbeScript {
	return s.add(record{op: opRevert, name: "revert", prefixSafe: true})
}

// SStore writes a storage slot.
//
// Not prefix-safe: EIP-8141 bans SSTORE during the validation prefix outside the first
// deploy frame. A new slot costs STATE_BYTES_PER_STORAGE_SET state gas, which is the
// dominant cost of the operation under EIP-8037.
func (s *ProbeScript) SStore(slot, value common.Hash) *ProbeScript {
	return s.add(record{
		op:       opSStore,
		a:        slot,
		b:        value,
		name:     "sstore",
		execGas:  2_200,
		stateGas: txtypes.StateBytesPerStorageSet * txtypes.CostPerStateByte,
	})
}

// Log emits a single-topic event with dataLength zero bytes of data, which is what a
// per-frame log attribution check needs.
func (s *ProbeScript) Log(topic common.Hash, dataLength uint64) *ProbeScript {
	return s.add(record{
		op:      opLog,
		a:       topic,
		b:       hashFromUint64(dataLength),
		name:    "log",
		execGas: 1_150 + 8*dataLength + memoryGas(0x100+dataLength),
	})
}

// Burn spins until the frame's remaining gas reaches gasTarget.
//
// Not prefix-safe: it reads GAS, which is banned during the validation prefix except
// immediately before a call.
func (s *ProbeScript) Burn(gasTarget uint64) *ProbeScript {
	return s.add(record{
		op:      opBurn,
		a:       hashFromUint64(gasTarget),
		name:    "burn",
		execGas: 0,
	})
}

// Approve calls APPROVE with the given scope and exits the frame, so anything after it
// in the script is unreachable.
//
// This is what makes the contract usable as a sender account or as a paymaster: the
// protocol's default code approves only for an account with no code of its own.
func (s *ProbeScript) Approve(scope uint8) *ProbeScript {
	return s.add(record{
		op:         opApprove,
		a:          hashFromUint64(uint64(scope)),
		name:       "approve",
		prefixSafe: true,
		execGas:    100,
	})
}

// ReadTxParam executes TXPARAM and discards the result.
//
// The introspection operations exist to make the instruction run inside a frame, not to
// check what it returns. Comparing against an expected value would bake one reading of
// the spec into every transaction, and a client that disagreed with that reading would
// show up as a failing frame rather than as a disagreement.
func (s *ProbeScript) ReadTxParam(param uint8) *ProbeScript {
	return s.add(record{
		op:         opReadTxParam,
		a:          hashFromUint64(uint64(param)),
		name:       "read_txparam",
		prefixSafe: true,
		execGas:    50,
	})
}

// ReadFrameParam executes FRAMEPARAM for one frame.
//
// Status and gas-used parameters exceptionally halt when read for the current or a later
// frame, so a caller reading those must name a frame that has already completed.
func (s *ProbeScript) ReadFrameParam(param uint8, frameIndex int) *ProbeScript {
	return s.add(record{
		op:         opReadFrameParam,
		a:          hashFromUint64(uint64(param)),
		b:          hashFromUint64(uint64(frameIndex)),
		name:       "read_frameparam",
		prefixSafe: true,
		execGas:    50,
	})
}

// ReadSigParam executes SIGPARAM for one signature entry.
//
// The scheme decides what is legal to ask: the resolved signer halts on an ARBITRARY
// entry, and the raw length halts on every protocol-validated scheme, whose bytes are
// deliberately not introspectable.
func (s *ProbeScript) ReadSigParam(param uint8, sigIndex int) *ProbeScript {
	return s.add(record{
		op:         opReadSigParam,
		a:          hashFromUint64(uint64(param)),
		b:          hashFromUint64(uint64(sigIndex)),
		name:       "read_sigparam",
		prefixSafe: true,
		execGas:    50,
	})
}

// ReadFrameData executes FRAMEDATALOAD against another frame's data.
func (s *ProbeScript) ReadFrameData(frameIndex int, offset uint64) *ProbeScript {
	return s.add(record{
		op:         opReadFrameData,
		a:          hashFromUint64(uint64(frameIndex)),
		b:          hashFromUint64(offset),
		name:       "read_framedata",
		prefixSafe: true,
		execGas:    50,
	})
}

// ReadSigData executes SIGDATACOPY against an ARBITRARY signature entry. Protocol-
// validated schemes are not introspectable and halt.
func (s *ProbeScript) ReadSigData(sigIndex int, offset uint64) *ProbeScript {
	return s.add(record{
		op:         opReadSigData,
		a:          hashFromUint64(uint64(sigIndex)),
		b:          hashFromUint64(offset),
		name:       "read_sigdata",
		prefixSafe: true,
		execGas:    60,
	})
}

// ReadRootRef executes RECENTROOTREFLOAD against a declared reference.
//
// Its opcode byte is the one EIP-8272 shares with EIP-8141's SIGDATACOPY, so a script
// using both operations puts a chain running both EIPs in front of that collision.
func (s *ProbeScript) ReadRootRef(index int, field uint8) *ProbeScript {
	return s.add(record{
		op:         opReadRootRef,
		a:          hashFromUint64(uint64(index)),
		b:          hashFromUint64(uint64(field)),
		name:       "read_rootref",
		prefixSafe: true,
		execGas:    50,
	})
}

// Call makes a plain call with no value and no data, reverting if it fails.
//
// Not prefix-safe: the validation trace rules restrict what may be called and forbid
// calling an address that is neither an existing contract nor a precompile, which the
// script cannot know on the caller's behalf.
func (s *ProbeScript) Call(target common.Address, gas uint64) *ProbeScript {
	return s.add(record{
		op:      opCall,
		a:       common.BytesToHash(target.Bytes()),
		b:       hashFromUint64(gas),
		name:    "call",
		execGas: 2_700 + gas,
	})
}

// WriteRoot makes the 64-byte call that commits a recent root, turning the calling
// contract into an EIP-8272 root source.
//
// Not prefix-safe: it is a state write and reads GAS.
func (s *ProbeScript) WriteRoot(salt, root common.Hash) *ProbeScript {
	return s.add(record{
		op:       opWriteRoot,
		a:        salt,
		b:        root,
		name:     "write_root",
		execGas:  10_000,
		stateGas: txtypes.StateBytesPerStorageSet * txtypes.CostPerStateByte,
	})
}

// StoreTxParam records a TXPARAM value in storage instead of asserting it, for values
// the caller wants to read back rather than predict.
func (s *ProbeScript) StoreTxParam(param uint8, slot common.Hash) *ProbeScript {
	return s.add(record{
		op:       opStoreTxParam,
		a:        hashFromUint64(uint64(param)),
		b:        slot,
		name:     "store_txparam",
		execGas:  2_200,
		stateGas: txtypes.StateBytesPerStorageSet * txtypes.CostPerStateByte,
	})
}

// ReadSLoad reads a storage slot and discards the value.
//
// Prefix-safe only against the sender's own storage, which the caller is responsible
// for: the validation trace rules allow SLOAD nowhere else.
func (s *ProbeScript) ReadSLoad(slot common.Hash) *ProbeScript {
	return s.add(record{
		op:         opReadSLoad,
		a:          slot,
		name:       "read_sload",
		prefixSafe: true,
		execGas:    2_150,
	})
}

// hashFromUint64 renders a number as a right-aligned word.
func hashFromUint64(v uint64) common.Hash {
	return common.BigToHash(new(uint256.Int).SetUint64(v).ToBig())
}

// hashFromU256 renders a 256-bit value as a word.
func hashFromU256(v *uint256.Int) common.Hash {
	if v == nil {
		return common.Hash{}
	}

	return common.BigToHash(v.ToBig())
}

// memoryGas returns the EVM memory expansion cost for a region ending at size bytes.
func memoryGas(size uint64) uint64 {
	words := (size + 31) / 32

	return 3*words + words*words/512
}
