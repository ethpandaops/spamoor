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

// approveEpilogue reads the approval scope from the frame's first calldata word and
// approves it, which exits the frame.
var approveEpilogue = []byte{
	0x60, 0x00, 0x35, // PUSH1 0x00, CALLDATALOAD -- scope
	0x60, 0x00, // PUSH1 0x00 -- length
	0x60, 0x00, // PUSH1 0x00 -- offset
	opcodeApprove,
}

// accountCodeSize bounds the fuzzed part of an account contract.
//
// It is small on purpose: the generated code runs before the epilogue, so the longer it
// is the less often it falls through to the APPROVE that lets the transaction land at
// all. Both outcomes are worth reaching, and a short prologue keeps the balance.
const accountCodeSize = 64

// GenerateAccountCode returns runtime code for an account that runs fuzzed code in its
// validation frame and then approves.
//
// It is what puts generated code inside the validation prefix, where EIP-8141's
// banned-opcode and storage rules apply and a public mempool node has to simulate it.
func GenerateAccountCode(seed string, variant int, gasLimit uint64) []byte {
	prologue := GenerateContractCode(seed+"-account", uint64(variant), 0, accountCodeSize, gasLimit)

	code := make([]byte, 0, len(prologue)+len(approveEpilogue))
	code = append(code, prologue...)
	code = append(code, approveEpilogue...)

	return code
}
