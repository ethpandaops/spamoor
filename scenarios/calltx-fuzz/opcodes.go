package calltxfuzz

import (
	evmfuzz "github.com/ethpandaops/spamoor/scenarios/evm-fuzz"
)

// op is a shorthand constructor avoiding verbose keyed-field literals
// for every opcode definition.
func op(name string, opcode uint16, stackIn, stackOut int, gas uint64, tmpl func() []byte, prob float64) *evmfuzz.OpcodeInfo {
	return &evmfuzz.OpcodeInfo{
		Name:        name,
		Opcode:      opcode,
		StackInput:  stackIn,
		StackOutput: stackOut,
		GasCost:     gas,
		Template:    tmpl,
		Probability: prob,
	}
}

// getCalltxFuzzOpcodeDefinitions returns opcode definitions with weights tuned
// for delegation-aware fuzzing. Storage, transient storage, cross-contract calls,
// and identity opcodes are weighted higher than in evm-fuzz to exercise
// delegation-relevant EVM behavior.
func getCalltxFuzzOpcodeDefinitions() []*evmfuzz.OpcodeInfo {
	s := evmfuzz.SimpleOpcode
	san := evmfuzz.SanitizeInput

	return []*evmfuzz.OpcodeInfo{
		op("STOP", 0x00, 0, 0, 0, s(0x00), 0),

		// Arithmetic operations (same weights as evm-fuzz)
		op("ADD", 0x01, 2, 1, 3, s(0x01), 1.0),
		op("MUL", 0x02, 2, 1, 5, s(0x02), 1.0),
		op("SUB", 0x03, 2, 1, 3, s(0x03), 1.0),
		op("DIV", 0x04, 2, 1, 5, s(0x04), 1.0),
		op("SDIV", 0x05, 2, 1, 5, s(0x05), 1.0),
		op("MOD", 0x06, 2, 1, 5, s(0x06), 1.0),
		op("SMOD", 0x07, 2, 1, 5, s(0x07), 1.0),
		op("ADDMOD", 0x08, 3, 1, 8, s(0x08), 1.0),
		op("MULMOD", 0x09, 3, 1, 8, s(0x09), 1.0),
		op("EXP", 0x0a, 2, 1, 10, s(0x0a), 1.0),
		op("SIGNEXTEND", 0x0b, 2, 1, 5, s(0x0b), 1.0),

		// Comparison operations
		op("LT", 0x10, 2, 1, 3, s(0x10), 1.0),
		op("GT", 0x11, 2, 1, 3, s(0x11), 1.0),
		op("SLT", 0x12, 2, 1, 3, s(0x12), 1.0),
		op("SGT", 0x13, 2, 1, 3, s(0x13), 1.0),
		op("EQ", 0x14, 2, 1, 3, s(0x14), 1.0),
		op("ISZERO", 0x15, 1, 1, 3, s(0x15), 1.0),
		op("AND", 0x16, 2, 1, 3, s(0x16), 1.0),
		op("OR", 0x17, 2, 1, 3, s(0x17), 1.0),
		op("XOR", 0x18, 2, 1, 3, s(0x18), 1.0),
		op("NOT", 0x19, 1, 1, 3, s(0x19), 1.0),
		op("BYTE", 0x1a, 2, 1, 3, san(s(0x1a), []uint64{0x3f}), 1.0),
		op("SHL", 0x1b, 2, 1, 3, s(0x1b), 1.0),
		op("SHR", 0x1c, 2, 1, 3, s(0x1c), 1.0),
		op("SAR", 0x1d, 2, 1, 3, s(0x1d), 1.0),
		op("CLZ", 0x1e, 1, 1, 3, s(0x1e), 1.5),

		// Crypto
		op("KECCAK256", 0x20, 2, 1, 30, san(s(0x20), []uint64{0x2ffff, 0xffff}), 1.5),

		// Environmental information — higher weights for identity-related opcodes
		op("ADDRESS", 0x30, 0, 1, 2, s(0x30), 4.0),
		op("BALANCE", 0x31, 1, 1, 100, s(0x31), 4.0),
		op("ORIGIN", 0x32, 0, 1, 2, s(0x32), 4.0),
		op("CALLER", 0x33, 0, 1, 2, s(0x33), 4.0),
		op("CALLVALUE", 0x34, 0, 1, 2, s(0x34), 2.0),
		op("CALLDATALOAD", 0x35, 1, 1, 3, s(0x35), 2.0),
		op("CALLDATASIZE", 0x36, 0, 1, 2, s(0x36), 2.0),
		op("CALLDATACOPY", 0x37, 3, 0, 3, san(s(0x37), []uint64{0x2ffff, 0x2ffff, 0xffff}), 2.0),
		op("CODESIZE", 0x38, 0, 1, 2, s(0x38), 2.0),
		op("CODECOPY", 0x39, 3, 0, 3, san(s(0x39), []uint64{0x2ffff, 0x2ffff, 0xffff}), 2.0),
		op("GASPRICE", 0x3a, 0, 1, 2, s(0x3a), 2.0),
		op("EXTCODESIZE", 0x3b, 1, 1, 100, s(0x3b), 3.5),
		op("EXTCODECOPY", 0x3c, 4, 0, 100, san(s(0x3c), []uint64{0, 0x2ffff, 0x2ffff, 0xffff}), 3.5),
		op("RETURNDATASIZE", 0x3d, 0, 1, 2, s(0x3d), 2.0),
		op("RETURNDATACOPY", 0x3e, 3, 0, 3, san(s(0x3e), []uint64{0x2ffff, 0xff, 0xff}), 2.0),
		op("EXTCODEHASH", 0x3f, 1, 1, 100, s(0x3f), 3.5),

		// Block information
		op("BLOCKHASH", 0x40, 1, 1, 20, s(0x40), 2.0),
		op("COINBASE", 0x41, 0, 1, 2, s(0x41), 2.0),
		op("TIMESTAMP", 0x42, 0, 1, 2, s(0x42), 2.0),
		op("NUMBER", 0x43, 0, 1, 2, s(0x43), 2.0),
		op("DIFFICULTY", 0x44, 0, 1, 2, s(0x44), 2.0),
		op("GASLIMIT", 0x45, 0, 1, 2, s(0x45), 2.0),
		op("CHAINID", 0x46, 0, 1, 2, s(0x46), 3.0),
		op("SELFBALANCE", 0x47, 0, 1, 5, s(0x47), 4.0),
		op("BASEFEE", 0x48, 0, 1, 2, s(0x48), 2.0),
		op("BLOBHASH", 0x49, 1, 1, 3, s(0x49), 2.0),
		op("BLOBBASEFEE", 0x4a, 0, 1, 2, s(0x4a), 2.0),

		// Stack/memory/storage — higher weights for storage and transient storage
		op("POP", 0x50, 1, 0, 2, s(0x50), 1.0),
		op("MLOAD", 0x51, 1, 1, 3, san(s(0x51), []uint64{0x2ffff}), 1.0),
		op("MSTORE", 0x52, 2, 0, 3, san(s(0x52), []uint64{0x2ffff}), 1.0),
		op("MSTORE8", 0x53, 2, 0, 3, san(s(0x53), []uint64{0x2ffff}), 1.0),
		op("SLOAD", 0x54, 1, 1, 100, san(s(0x54), []uint64{0x2ffff}), 4.0),
		op("SSTORE", 0x55, 2, 0, 100, san(s(0x55), []uint64{0x2ffff}), 4.0),
		op("PC", 0x58, 0, 1, 2, s(0x58), 1.0),
		op("MSIZE", 0x59, 0, 1, 2, s(0x59), 1.0),
		op("GAS", 0x5a, 0, 1, 2, s(0x5a), 1.0),
		op("JUMPDEST", 0x5b, 0, 0, 1, s(0x5b), 1.0),
		op("TLOAD", 0x5c, 1, 1, 100, san(s(0x5c), []uint64{0x2ffff}), 5.0),
		op("TSTORE", 0x5d, 2, 0, 100, san(s(0x5d), []uint64{0x2ffff}), 5.0),
		op("MCOPY", 0x5e, 3, 0, 3, san(s(0x5e), []uint64{0x2ffff, 0x2ffff, 0xffff}), 2.0),

		// Push operations
		op("PUSH0", 0x5f, 0, 1, 2, s(0x5f), 1.0),

		// Log operations
		op("LOG0", 0xa0, 2, 0, 375, san(s(0xa0), []uint64{0x2ffff, 0xffff}), 1.0),
		op("LOG1", 0xa1, 3, 0, 750, san(s(0xa1), []uint64{0x2ffff, 0xffff}), 1.0),
		op("LOG2", 0xa2, 4, 0, 1125, san(s(0xa2), []uint64{0x2ffff, 0xffff}), 1.0),
		op("LOG3", 0xa3, 5, 0, 1500, san(s(0xa3), []uint64{0x2ffff, 0xffff}), 1.0),
		op("LOG4", 0xa4, 6, 0, 1875, san(s(0xa4), []uint64{0x2ffff, 0xffff}), 1.0),

		// Contract call operations — much higher weights for delegation testing
		// Gas masks (0x4fffff ~5M) prevent uint64 overflow on the gas argument.
		// CALL/CALLCODE take 7 stack inputs (gas,addr,value,argsOff,argsLen,retOff,retLen).
		// Value mask 0xffff allows small value transfers; most will succeed since
		// the transaction itself sends ETH. "Insufficient balance" still occurs
		// when contracts have zero balance (25% of txs send zero value).
		op("CALL", 0xf1, 7, 1, 100, san(s(0xf1), []uint64{0x4fffff, 0xffff, 0xffff, 0x2ffff, 0xffff, 0x2ffff, 0xffff}), 5.0),
		op("CALLCODE", 0xf2, 7, 1, 100, san(s(0xf2), []uint64{0x4fffff, 0xffff, 0xffff, 0x2ffff, 0xffff, 0x2ffff, 0xffff}), 2.0),
		op("RETURN", 0xf3, 2, 0, 0, san(s(0xf3), []uint64{0x2ffff, 0xffff}), 0.1),
		// DELEGATECALL/STATICCALL take 6 stack inputs (gas,addr,argsOff,argsLen,retOff,retLen).
		op("DELEGATECALL", 0xf4, 6, 1, 100, san(s(0xf4), []uint64{0x4fffff, 0, 0x2ffff, 0xffff, 0x2ffff, 0xffff}), 5.0),
		// Lower STATICCALL weight: callees often write (SSTORE/TSTORE/LOG),
		// causing write-protection errors. Keep it available but less frequent.
		op("STATICCALL", 0xfa, 6, 1, 100, san(s(0xfa), []uint64{0x4fffff, 0, 0x2ffff, 0xffff, 0x2ffff, 0xffff}), 2.0),
		op("REVERT", 0xfd, 2, 0, 0, san(s(0xfd), []uint64{0x2ffff, 0xffff}), 0.1),
		op("SELFDESTRUCT", 0xff, 1, 0, 5000, s(0xff), 1.1),

		// CREATE/CREATE2 — slightly higher for delegation context.
		// Small value mask (0xff) allows occasional value transfers.
		op("CREATE", 0xf0, 3, 1, 32000, san(s(0xf0), []uint64{0xff, 0x2ffff, 0xffff}), 2.5),
		op("CREATE2", 0xf5, 4, 1, 32000, san(s(0xf5), []uint64{0xff, 0x2ffff, 0xffff}), 2.5),
	}
}
