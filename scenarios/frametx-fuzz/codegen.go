package frametxfuzz

import (
	evmfuzz "github.com/ethpandaops/spamoor/scenarios/evm-fuzz"
	"github.com/ethpandaops/spamoor/utils"
)

// Generated contracts: fuzzed EVM code that one frame deploys and another calls, built
// by evm-fuzz's stack-aware generator with the frame instructions added to its table.

// Frame instruction opcodes. RECENTROOTREFLOAD gets a synthetic key because EIP-8141 and
// EIP-8272 both assign it the byte 0xb5 and the generator's table is keyed by opcode; the
// template still emits the real byte.
const (
	opcodeApprove           = 0xaa
	opcodeTxParam           = 0xb0
	opcodeFrameDataLoad     = 0xb1
	opcodeFrameDataCopy     = 0xb2
	opcodeFrameParam        = 0xb3
	opcodeSigParam          = 0xb4
	opcodeSigDataCopy       = 0xb5
	opcodeRecentRootRefLoad = 0x1b5

	// recentRootRefLoadByte is what RECENTROOTREFLOAD actually assembles to. It is the
	// same byte as SIGDATACOPY on purpose.
	recentRootRefLoadByte = 0xb5
)

// push1 emits a single-byte push.
func push1(v byte) []byte { return []byte{0x60, v} }

// concat joins operand pushes and the instruction itself.
func concat(parts ...[]byte) []byte {
	out := []byte{}
	for _, part := range parts {
		out = append(out, part...)
	}

	return out
}

// frameOpcodeDefinitions describes the frame instructions to the generator.
//
// Every template pushes its own operands, so the generator's stack tracking stays exact.
// Operands are drawn small but not always in range: an index past the end of the frame or
// signature list halts, which is as much part of the space as a valid one.
func frameOpcodeDefinitions(rng *utils.DeterministicRNG) []*evmfuzz.OpcodeInfo {
	// Kept small so a copy does not spend the frame's budget on memory expansion.
	memOffset := func() []byte { return push1(byte(rng.Intn(0x80))) }
	length := func() []byte { return push1(byte(rng.Intn(0x40))) }
	frameIndex := func() []byte { return push1(byte(rng.Intn(8))) }
	sigIndex := func() []byte { return push1(byte(rng.Intn(4))) }

	return []*evmfuzz.OpcodeInfo{
		{
			Name: "TXPARAM", Opcode: opcodeTxParam, StackInput: 0, StackOutput: 1, GasCost: 2,
			Probability: 1.5,
			Template: func() []byte {
				// 0x00-0x10 spans EIP-8141's indices and those EIP-8250 and EIP-8272
				// add, plus one past the end.
				return concat(push1(byte(rng.Intn(0x12))), []byte{opcodeTxParam})
			},
		},
		{
			Name: "FRAMEPARAM", Opcode: opcodeFrameParam, StackInput: 0, StackOutput: 1, GasCost: 2,
			Probability: 1.5,
			Template: func() []byte {
				// param is second from top, frameIndex on top.
				return concat(push1(byte(rng.Intn(0x0d))), frameIndex(), []byte{opcodeFrameParam})
			},
		},
		{
			Name: "SIGPARAM", Opcode: opcodeSigParam, StackInput: 0, StackOutput: 1, GasCost: 2,
			Probability: 1.0,
			Template: func() []byte {
				return concat(push1(byte(rng.Intn(0x05))), sigIndex(), []byte{opcodeSigParam})
			},
		},
		{
			Name: "FRAMEDATALOAD", Opcode: opcodeFrameDataLoad, StackInput: 0, StackOutput: 1, GasCost: 3,
			Probability: 1.0,
			Template: func() []byte {
				// frameIndex is second from top, offset on top.
				return concat(frameIndex(), push1(byte(rng.Intn(0x40))), []byte{opcodeFrameDataLoad})
			},
		},
		{
			Name: "FRAMEDATACOPY", Opcode: opcodeFrameDataCopy, StackInput: 0, StackOutput: 0, GasCost: 6,
			Probability: 0.8,
			Template: func() []byte {
				// Deepest first: frameIndex, length, dataOffset, memOffset.
				return concat(frameIndex(), length(), push1(byte(rng.Intn(0x40))), memOffset(),
					[]byte{opcodeFrameDataCopy})
			},
		},
		{
			Name: "SIGDATACOPY", Opcode: opcodeSigDataCopy, StackInput: 0, StackOutput: 0, GasCost: 6,
			Probability: 0.8,
			Template: func() []byte {
				return concat(sigIndex(), length(), push1(byte(rng.Intn(0x40))), memOffset(),
					[]byte{opcodeSigDataCopy})
			},
		},
		{
			Name: "RECENTROOTREFLOAD", Opcode: opcodeRecentRootRefLoad, StackInput: 0, StackOutput: 1, GasCost: 3,
			Probability: 0.8,
			Template: func() []byte {
				// index is second from top, field on top; field > 2 halts.
				return concat(push1(byte(rng.Intn(4))), push1(byte(rng.Intn(4))),
					[]byte{recentRootRefLoadByte})
			},
		},
		{
			Name: "APPROVE", Opcode: opcodeApprove, StackInput: 0, StackOutput: 0, GasCost: 0,
			// Rare: APPROVE exits the call frame, so everything after it is unreachable.
			Probability: 0.2,
			Template: func() []byte {
				// Deepest first: scope, length, offset.
				return concat(push1(byte(rng.Intn(4))), push1(0), push1(0), []byte{opcodeApprove})
			},
		},
	}
}

// GenerateContractCode returns fuzzed runtime code for a contract a frame deploys.
//
// The stream is derived from the seed, transaction index and frame position, so a
// replayed recipe deploys byte-identical code without shifting the recipe draw.
func GenerateContractCode(seed string, txID uint64, frameIndex int, maxSize int, gasLimit uint64) []byte {
	streamID := txID*64 + uint64(frameIndex)

	generator := evmfuzz.NewOpcodeGenerator(streamID, seed, maxSize, gasLimit)
	generator.AddOpcodes(frameOpcodeDefinitions(utils.NewDeterministicRNGWithSeed(streamID, seed+"-frame")))
	generator.SetFuzzMode("all")

	return generator.Generate()
}

// accountCodeSize bounds the fuzzed part of an account contract.
//
// Small on purpose. The code runs inside the validation prefix, which is capped at
// MaxVerifyGas across every frame in it, and the longer the prologue the less often it
// falls through to the APPROVE that lets the transaction land at all.
const accountCodeSize = 48

// approveTail reads the approval scope from the first calldata word and approves it,
// which exits the frame. It runs after the fuzzed prologue on the role path.
var approveTail = []byte{
	0x60, 0x00, 0x35, // PUSH1 0x00, CALLDATALOAD -- scope
	0x60, 0x00, // PUSH1 0x00 -- length
	0x60, 0x00, // PUSH1 0x00 -- offset
	opcodeApprove,
	0x00, // STOP -- unreachable after APPROVE, guards against fall-through
}

// stopTail ends the role path for an account that does not approve.
var stopTail = []byte{0x00}

// GenerateAccountCode returns runtime code for an account that plays the sender or
// paymaster role and can be swept afterwards.
//
// Its calldata's first word selects the behaviour: a non-zero scope runs the fuzzed
// prologue and then, unless approves is false, APPROVEs that scope; a zero scope jumps
// straight to a tail that sends the account's whole balance to CALLER. The role path is
// what puts generated code inside the validation prefix, where EIP-8141's banned-opcode
// and storage rules apply; the wipe path is how the funding deployed with the account is
// reclaimed once it has been used -- a later SENDER-mode frame calls it, so CALLER is that
// transaction's sender and the funds return to a tracked wallet.
func GenerateAccountCode(seed string, txID uint64, frameIndex int, gasLimit uint64, approves bool) []byte {
	prologue := GenerateContractCode(seed+"-account", txID, frameIndex, accountCodeSize, gasLimit)

	roleTail := approveTail
	if !approves {
		roleTail = stopTail
	}

	// The wipe JUMPDEST sits after the dispatch, the prologue and the role tail.
	wipePC := len(accountDispatch) + len(prologue) + len(roleTail)

	code := make([]byte, 0, wipePC+len(wipeTail))
	code = append(code, accountDispatch...)
	code = append(code, prologue...)
	code = append(code, roleTail...)
	code = append(code, wipeTail...)

	// Patch the dispatch's PUSH2 with the wipe JUMPDEST offset.
	code[dispatchSweepPCOffset] = byte(wipePC >> 8)
	code[dispatchSweepPCOffset+1] = byte(wipePC)

	return code
}

// accountDispatch reads the scope word and jumps to the wipe tail when it is zero,
// leaving the scope on the stack for the role path otherwise.
var accountDispatch = []byte{
	0x60, 0x00, 0x35, // PUSH1 0x00, CALLDATALOAD -- scope
	0x80,       // DUP1
	0x15,       // ISZERO
	0x61, 0, 0, // PUSH2 <sweepPC> (patched)
	0x57, // JUMPI
}

// dispatchSweepPCOffset is where accountDispatch carries the PUSH2 wipe offset.
const dispatchSweepPCOffset = 6

// wipeTail sends the account's whole balance to CALLER with an empty call, then stops.
// CALLER is the wiping transaction's sender, so the funding returns to a tracked wallet.
var wipeTail = []byte{
	0x5b,       // JUMPDEST
	0x60, 0x00, // PUSH1 0 -- retLength
	0x60, 0x00, // PUSH1 0 -- retOffset
	0x60, 0x00, // PUSH1 0 -- argsLength
	0x60, 0x00, // PUSH1 0 -- argsOffset
	0x47, // SELFBALANCE -- value
	0x33, // CALLER -- destination
	0x5a, // GAS
	0xf1, // CALL
	0x50, // POP
	0x00, // STOP
}
