package evmfuzz

import "fmt"

// OpcodeInfo defines properties of an EVM opcode
type OpcodeInfo struct {
	Name        string
	Opcode      uint16 // uint16 to support precompiles (≥0x100) and regular opcodes (≤0xff)
	StackInput  int    // Number of items consumed from stack
	StackOutput int    // Number of items pushed to stack
	GasCost     uint64
	Template    func() []byte // Function to generate valid sequence
	Probability float64       // Relative probability weight for selection (1.0 = normal)
}

// simpleOpcode creates a template function for a single-byte opcode
func simpleOpcode(opcode byte) func() []byte {
	return func() []byte { return []byte{opcode} }
}

// sanitizeInput is a generic sanitization wrapper that applies masks to stack positions
// before executing base template. masks[i] = mask for stack position i (0 = no sanitization)
func sanitizeInput(baseTemplate func() []byte, masks []uint64) func() []byte {
	return func() []byte {
		var bytecode []byte

		// Apply sanitization to each specified stack position
		for i, mask := range masks {
			if mask == 0 {
				continue // Skip positions that don't need sanitization
			}

			// Generate appropriate PUSH instruction for the mask
			var pushBytes []byte
			if mask <= 0xFF {
				pushBytes = []byte{0x60, byte(mask)} // PUSH1
			} else if mask <= 0xFFFF {
				pushBytes = []byte{0x61, byte(mask >> 8), byte(mask)} // PUSH2
			} else if mask <= 0x2ffff {
				pushBytes = []byte{0x62, byte(mask >> 16), byte(mask >> 8), byte(mask)} // PUSH3
			} else {
				pushBytes = []byte{0x63, byte(mask >> 24), byte(mask >> 16), byte(mask >> 8), byte(mask)} // PUSH4
			}

			if i == 0 {
				// Sanitize stack[0] (top of stack)
				bytecode = append(bytecode, pushBytes...)
				bytecode = append(bytecode, 0x16) // AND
			} else {
				// Sanitize stack[i] - need to use DUP to copy to top, sanitize, then SWAP back
				dupOpcode := byte(0x80 + i)  // DUP(i+1)
				swapOpcode := byte(0x90 + i) // SWAP(i+1)

				bytecode = append(bytecode, dupOpcode)    // DUP(i+1) - copy stack[i] to top
				bytecode = append(bytecode, pushBytes...) // PUSH mask
				bytecode = append(bytecode, 0x16)         // AND (sanitize)
				bytecode = append(bytecode, swapOpcode)   // SWAP(i+1) - put sanitized value back
				bytecode = append(bytecode, 0x50)         // POP (remove original unsanitized value)
			}
		}

		// Execute the base template
		baseBytes := baseTemplate()
		bytecode = append(bytecode, baseBytes...)

		return bytecode
	}
}

// getBaseOpcodeDefinitions returns the static opcode definitions that don't require generator reference
func getBaseOpcodeDefinitions() []*OpcodeInfo {
	return []*OpcodeInfo{
		{"STOP", 0x00, 0, 0, 0, simpleOpcode(0x00), 0},

		// Arithmetic operations
		{"ADD", 0x01, 2, 1, 3, simpleOpcode(0x01), 1.0},
		{"MUL", 0x02, 2, 1, 5, simpleOpcode(0x02), 1.0},
		{"SUB", 0x03, 2, 1, 3, simpleOpcode(0x03), 1.0},
		{"DIV", 0x04, 2, 1, 5, simpleOpcode(0x04), 1.0},
		{"SDIV", 0x05, 2, 1, 5, simpleOpcode(0x05), 1.0},
		{"MOD", 0x06, 2, 1, 5, simpleOpcode(0x06), 1.0},
		{"SMOD", 0x07, 2, 1, 5, simpleOpcode(0x07), 1.0},
		{"ADDMOD", 0x08, 3, 1, 8, simpleOpcode(0x08), 1.0},
		{"MULMOD", 0x09, 3, 1, 8, simpleOpcode(0x09), 1.0},
		{"EXP", 0x0a, 2, 1, 10, simpleOpcode(0x0a), 1.0},
		{"SIGNEXTEND", 0x0b, 2, 1, 5, simpleOpcode(0x0b), 1.0},

		// Comparison operations
		{"LT", 0x10, 2, 1, 3, simpleOpcode(0x10), 1.0},
		{"GT", 0x11, 2, 1, 3, simpleOpcode(0x11), 1.0},
		{"SLT", 0x12, 2, 1, 3, simpleOpcode(0x12), 1.0},
		{"SGT", 0x13, 2, 1, 3, simpleOpcode(0x13), 1.0},
		{"EQ", 0x14, 2, 1, 3, simpleOpcode(0x14), 1.0},
		{"ISZERO", 0x15, 1, 1, 3, simpleOpcode(0x15), 1.0},
		{"AND", 0x16, 2, 1, 3, simpleOpcode(0x16), 1.0},
		{"OR", 0x17, 2, 1, 3, simpleOpcode(0x17), 1.0},
		{"XOR", 0x18, 2, 1, 3, simpleOpcode(0x18), 1.0},
		{"NOT", 0x19, 1, 1, 3, simpleOpcode(0x19), 1.0},
		{"BYTE", 0x1a, 2, 1, 3, sanitizeInput(simpleOpcode(0x1a), []uint64{0x3f}), 1.0},
		{"SHL", 0x1b, 2, 1, 3, simpleOpcode(0x1b), 1.0},
		{"SHR", 0x1c, 2, 1, 3, simpleOpcode(0x1c), 1.0},
		{"SAR", 0x1d, 2, 1, 3, simpleOpcode(0x1d), 1.0},
		{"CLZ", 0x1e, 1, 1, 3, simpleOpcode(0x1e), 1.5},

		// Crypto operations
		{"KECCAK256", 0x20, 2, 1, 30, sanitizeInput(simpleOpcode(0x20), []uint64{0x2ffff, 0xffff}), 1.5},

		// Environmental information
		{"ADDRESS", 0x30, 0, 1, 2, simpleOpcode(0x30), 2.0},
		{"BALANCE", 0x31, 1, 1, 100, simpleOpcode(0x31), 2.2},
		{"ORIGIN", 0x32, 0, 1, 2, simpleOpcode(0x32), 2.0},
		{"CALLER", 0x33, 0, 1, 2, simpleOpcode(0x33), 2.0},
		{"CALLVALUE", 0x34, 0, 1, 2, simpleOpcode(0x34), 2.0},
		{"CALLDATALOAD", 0x35, 1, 1, 3, simpleOpcode(0x35), 2.0},
		{"CALLDATASIZE", 0x36, 0, 1, 2, simpleOpcode(0x36), 2.0},
		{"CALLDATACOPY", 0x37, 3, 0, 3, sanitizeInput(simpleOpcode(0x37), []uint64{0x2ffff, 0x2ffff, 0xffff}), 2.0},
		{"CODESIZE", 0x38, 0, 1, 2, simpleOpcode(0x38), 2.0},
		{"CODECOPY", 0x39, 3, 0, 3, sanitizeInput(simpleOpcode(0x39), []uint64{0x2ffff, 0x2ffff, 0xffff}), 2.0},
		{"GASPRICE", 0x3a, 0, 1, 2, simpleOpcode(0x3a), 2.0},
		{"EXTCODESIZE", 0x3b, 1, 1, 100, simpleOpcode(0x3b), 2.0},
		{"EXTCODECOPY", 0x3c, 4, 0, 100, sanitizeInput(simpleOpcode(0x3c), []uint64{0, 0x2ffff, 0x2ffff, 0xffff}), 2.0},
		{"RETURNDATASIZE", 0x3d, 0, 1, 2, simpleOpcode(0x3d), 2.0},
		{"RETURNDATACOPY", 0x3e, 3, 0, 3, sanitizeInput(simpleOpcode(0x3e), []uint64{0x2ffff, 0x2ffff, 0xffff}), 2.0},
		{"EXTCODEHASH", 0x3f, 1, 1, 100, simpleOpcode(0x3f), 2.0},

		// Block information
		{"BLOCKHASH", 0x40, 1, 1, 20, simpleOpcode(0x40), 2.0},
		{"COINBASE", 0x41, 0, 1, 2, simpleOpcode(0x41), 2.0},
		{"TIMESTAMP", 0x42, 0, 1, 2, simpleOpcode(0x42), 2.0},
		{"NUMBER", 0x43, 0, 1, 2, simpleOpcode(0x43), 2.0},
		{"DIFFICULTY", 0x44, 0, 1, 2, simpleOpcode(0x44), 2.0},
		{"GASLIMIT", 0x45, 0, 1, 2, simpleOpcode(0x45), 2.0},
		{"CHAINID", 0x46, 0, 1, 2, simpleOpcode(0x46), 2.5},
		{"SELFBALANCE", 0x47, 0, 1, 5, simpleOpcode(0x47), 2.5},
		{"BASEFEE", 0x48, 0, 1, 2, simpleOpcode(0x48), 2.5},
		{"BLOBHASH", 0x49, 1, 1, 3, simpleOpcode(0x49), 2.0},
		{"BLOBBASEFEE", 0x4a, 0, 1, 2, simpleOpcode(0x4a), 2.0},
		{"SLOTNUM", 0x4b, 0, 1, 2, simpleOpcode(0x4b), 2.0},

		// Stack operations with memory/storage sanitization
		{"POP", 0x50, 1, 0, 2, simpleOpcode(0x50), 1.0},
		{"MLOAD", 0x51, 1, 1, 3, sanitizeInput(simpleOpcode(0x51), []uint64{0x2ffff}), 1.0},
		{"MSTORE", 0x52, 2, 0, 3, sanitizeInput(simpleOpcode(0x52), []uint64{0x2ffff}), 1.0},
		{"MSTORE8", 0x53, 2, 0, 3, sanitizeInput(simpleOpcode(0x53), []uint64{0x2ffff}), 1.0},
		{"SLOAD", 0x54, 1, 1, 100, sanitizeInput(simpleOpcode(0x54), []uint64{0x2ffff}), 1.2},
		{"SSTORE", 0x55, 2, 0, 100, sanitizeInput(simpleOpcode(0x55), []uint64{0x2ffff}), 1.2},
		// JUMP and JUMPI are added by generator (need generator reference)
		{"PC", 0x58, 0, 1, 2, simpleOpcode(0x58), 1.0},
		{"MSIZE", 0x59, 0, 1, 2, simpleOpcode(0x59), 1.0},
		{"GAS", 0x5a, 0, 1, 2, simpleOpcode(0x5a), 1.0},
		{"JUMPDEST", 0x5b, 0, 0, 1, simpleOpcode(0x5b), 1.0},
		{"TLOAD", 0x5c, 1, 1, 100, sanitizeInput(simpleOpcode(0x5c), []uint64{0x2ffff}), 2.0},
		{"TSTORE", 0x5d, 2, 0, 100, sanitizeInput(simpleOpcode(0x5d), []uint64{0x2ffff}), 2.0},
		{"MCOPY", 0x5e, 3, 0, 3, sanitizeInput(simpleOpcode(0x5e), []uint64{0x2ffff, 0x2ffff, 0xffff}), 2.0},

		// Push operations
		{"PUSH0", 0x5f, 0, 1, 2, simpleOpcode(0x5f), 1.0},
		// PUSH1-PUSH32 are added dynamically by generator

		// Log operations
		{"LOG0", 0xa0, 2, 0, 375, sanitizeInput(simpleOpcode(0xa0), []uint64{0x2ffff, 0xffff}), 1.0},
		{"LOG1", 0xa1, 3, 0, 750, sanitizeInput(simpleOpcode(0xa1), []uint64{0x2ffff, 0xffff}), 1.0},
		{"LOG2", 0xa2, 4, 0, 1125, sanitizeInput(simpleOpcode(0xa2), []uint64{0x2ffff, 0xffff}), 1.0},
		{"LOG3", 0xa3, 5, 0, 1500, sanitizeInput(simpleOpcode(0xa3), []uint64{0x2ffff, 0xffff}), 1.0},
		{"LOG4", 0xa4, 6, 0, 1875, sanitizeInput(simpleOpcode(0xa4), []uint64{0x2ffff, 0xffff}), 1.0},

		// Contract operations (static ones)
		{"CALL", 0xf1, 6, 1, 100, sanitizeInput(simpleOpcode(0xf1), []uint64{0, 0xffff, 0, 0x2ffff, 0xffff, 0x2ffff, 0xffff}), 1.3},
		{"CALLCODE", 0xf2, 6, 1, 100, sanitizeInput(simpleOpcode(0xf2), []uint64{0, 0xffff, 0, 0x2ffff, 0xffff, 0x2ffff, 0xffff}), 1.0},
		{"RETURN", 0xf3, 2, 0, 0, sanitizeInput(simpleOpcode(0xf3), []uint64{0x2ffff, 0xffff}), 0.1},
		{"DELEGATECALL", 0xf4, 6, 1, 100, sanitizeInput(simpleOpcode(0xf4), []uint64{0, 0, 0x2ffff, 0xffff, 0x2ffff, 0xffff}), 1.3},
		// CREATE and CREATE2 are added by generator (need generator reference)
		{"STATICCALL", 0xfa, 6, 1, 100, sanitizeInput(simpleOpcode(0xfa), []uint64{0, 0, 0x2ffff, 0xffff, 0x2ffff, 0xffff}), 1.3},
		{"REVERT", 0xfd, 2, 0, 0, sanitizeInput(simpleOpcode(0xfd), []uint64{0x2ffff, 0xffff}), 0.1},
		{"SELFDESTRUCT", 0xff, 1, 0, 5000, simpleOpcode(0xff), 1.1},
	}
}

// getDupOpcodeDefinitions returns DUP1-DUP16 opcode definitions
func getDupOpcodeDefinitions() []*OpcodeInfo {
	opcodes := make([]*OpcodeInfo, 0, 16)
	for i := 1; i <= 16; i++ {
		dupOpcode := uint16(0x7f + i)
		dupDepth := i
		opcodes = append(opcodes, &OpcodeInfo{
			Name:        fmt.Sprintf("DUP%d", i),
			Opcode:      dupOpcode,
			StackInput:  dupDepth,
			StackOutput: dupDepth + 1,
			GasCost:     3,
			Template:    simpleOpcode(byte(dupOpcode)),
			Probability: 1.0,
		})
	}
	return opcodes
}

// getSwapOpcodeDefinitions returns SWAP1-SWAP16 opcode definitions
func getSwapOpcodeDefinitions() []*OpcodeInfo {
	opcodes := make([]*OpcodeInfo, 0, 16)
	for i := 1; i <= 16; i++ {
		swapOpcode := uint16(0x8f + i)
		swapDepth := i + 1
		opcodes = append(opcodes, &OpcodeInfo{
			Name:        fmt.Sprintf("SWAP%d", i),
			Opcode:      swapOpcode,
			StackInput:  swapDepth,
			StackOutput: swapDepth,
			GasCost:     3,
			Template:    simpleOpcode(byte(swapOpcode)),
			Probability: 1.0,
		})
	}
	return opcodes
}
