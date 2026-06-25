package evmfuzz

// ETH value-transfer and SELFDESTRUCT generators for the "transfers" fuzz mode.
// These exercise every value-bearing path so EIP-7708 (transfer logs) and
// EIP-8246 (SELFDESTRUCT burn/credit) become observable for cross-client diffing.
// Each method emits self-contained bytecode with strict stack accounting:
// StackInput=0 and a fixed StackOutput, with the template managing all real
// stack args internally (mirrors the precompile handler convention).

// isTransferOpcode checks if an opcode is a synthetic ETH-transfer generator.
func (g *OpcodeGenerator) isTransferOpcode(opcode uint16) bool {
	return opcode >= 0x120 && opcode <= 0x125
}

// isValueBaseOpcode lets the transfers mode also pick stack-building and
// value-relevant base opcodes so generated programs have operands and variety.
func (g *OpcodeGenerator) isValueBaseOpcode(opcode uint16) bool {
	switch opcode {
	case 0x30, 0x32, 0x33, 0x34, 0x47, 0x31, // ADDRESS,ORIGIN,CALLER,CALLVALUE,SELFBALANCE,BALANCE
		0x00, 0x50, 0x5f, 0x52, 0x5b, // STOP,POP,PUSH0,MSTORE,JUMPDEST
		0xa0, 0xa1: // LOG0,LOG1 (so non-transfer logs interleave)
		return true
	}
	// PUSH1..PUSH32 give operands
	return opcode >= 0x60 && opcode <= 0x7f
}

// pushBeneficiary emits a PUSH that places a target address on the stack.
// Deterministic 6-way set covering self/caller/origin, random EOAs, a
// known-nonexistent address, and a precompile/system address. Each branch
// nets exactly one stack item.
func (g *OpcodeGenerator) pushBeneficiary() []byte {
	switch g.rng.Intn(6) {
	case 0:
		return []byte{0x30} // ADDRESS (self-transfer)
	case 1:
		return []byte{0x33} // CALLER
	case 2:
		return []byte{0x32} // ORIGIN
	case 3:
		return append([]byte{0x73}, g.rng.Bytes(20)...) // PUSH20 random EOA
	case 4:
		// Known-nonexistent: high byte forced 0xde, almost surely empty
		addr := g.rng.Bytes(20)
		addr[0] = 0xde
		return append([]byte{0x73}, addr...)
	default:
		return []byte{0x60, byte(1 + g.rng.Intn(10))} // PUSH1 precompile/system addr
	}
}

// pushTransferValue pushes an ETH amount: weighted toward 0-value and tiny
// amounts so transfers succeed and the 0-value EIP-7708 log path is hit.
// Each branch pushes exactly one stack item.
func (g *OpcodeGenerator) pushTransferValue() []byte {
	switch g.rng.Intn(4) {
	case 0:
		return []byte{0x5f} // PUSH0 -> 0-value transfer (still emits a 7708 log)
	case 1:
		return []byte{0x60, byte(1 + g.rng.Intn(0xff))} // PUSH1 small wei
	case 2:
		v := g.rng.Intn(0x6000) + 0x100
		return []byte{0x61, byte(v >> 8), byte(v)} // PUSH2
	default:
		return []byte{0x34} // CALLVALUE -> forward this tx's value
	}
}

// generateValueCall: CALL with value (0x120, StackOutput=1).
// Capped gas (10000) avoids draining the frame; transfer success/flag pushed.
func (g *OpcodeGenerator) generateValueCall() []byte {
	return g.emitValueCall(0xf1)
}

// generateValueCallcode: CALLCODE with value (0x121, StackOutput=1).
func (g *OpcodeGenerator) generateValueCallcode() []byte {
	return g.emitValueCall(0xf2)
}

// emitValueCall builds a value-bearing CALL/CALLCODE. Net stack effect: +1.
// CALL stack (top->bottom): gas, addr, value, argsOffset, argsSize, retOffset, retSize.
func (g *OpcodeGenerator) emitValueCall(callOp byte) []byte {
	var bc []byte
	bc = append(bc, 0x60, 0x00)               // PUSH1 0  retSize
	bc = append(bc, 0x60, 0x00)               // PUSH1 0  retOffset
	bc = append(bc, 0x60, 0x00)               // PUSH1 0  argsSize
	bc = append(bc, 0x60, 0x00)               // PUSH1 0  argsOffset
	bc = append(bc, g.pushTransferValue()...) // value
	bc = append(bc, g.pushBeneficiary()...)   // addr
	bc = append(bc, 0x61, 0x27, 0x10)         // PUSH2 10000  capped gas
	bc = append(bc, callOp)                   // CALL/CALLCODE -> success flag
	return bc
}

// generateValueCreate: CREATE with value (0x122, StackOutput=1).
func (g *OpcodeGenerator) generateValueCreate() []byte {
	return g.emitValueCreate(false)
}

// generateValueCreate2: CREATE2 with value (0x123, StackOutput=1).
func (g *OpcodeGenerator) generateValueCreate2() []byte {
	return g.emitValueCreate(true)
}

// emitValueCreate stores empty-runtime initcode at mem[0] then CREATE/CREATE2
// with an ETH endowment. Net stack effect: +1 (address or 0).
func (g *OpcodeGenerator) emitValueCreate(create2 bool) []byte {
	var bc []byte
	// initcode: 60 00 60 00 f3 (PUSH1 0 PUSH1 0 RETURN -> deploys empty runtime)
	word := make([]byte, 32)
	copy(word, []byte{0x60, 0x00, 0x60, 0x00, 0xf3})
	bc = append(bc, 0x7f)    // PUSH32
	bc = append(bc, word...) // initcode word
	bc = append(bc, 0x5f)    // PUSH0  mem offset 0
	bc = append(bc, 0x52)    // MSTORE
	if create2 {
		bc = append(bc, 0x63)              // PUSH4
		bc = append(bc, g.rng.Bytes(4)...) // salt
	}
	bc = append(bc, 0x60, 0x05)               // PUSH1 5  size (initcode length)
	bc = append(bc, 0x5f)                     // PUSH0    offset 0
	bc = append(bc, g.pushTransferValue()...) // value (endowment)
	if create2 {
		bc = append(bc, 0xf5) // CREATE2 -> address
	} else {
		bc = append(bc, 0xf0) // CREATE -> address
	}
	return bc
}

// generateRevertingTransfer deploys a runtime that always REVERTs, then CALLs
// it with value (0x124, StackOutput=1). The EIP-7708 transfer-log must roll
// back when the inner frame reverts; clients are diffed on this. Net: +1.
func (g *OpcodeGenerator) generateRevertingTransfer() []byte {
	var bc []byte
	// initcode returns a 5-byte REVERT runtime (60 00 60 00 fd):
	//   PUSH1 05 PUSH1 0c PUSH1 00 CODECOPY  PUSH1 05 PUSH1 00 RETURN  || runtime
	initcode := []byte{
		0x60, 0x05, 0x60, 0x0c, 0x60, 0x00, 0x39, // copy 5 bytes from offset 12
		0x60, 0x05, 0x60, 0x00, 0xf3, // return mem[0:5]
		0x60, 0x00, 0x60, 0x00, 0xfd, // runtime: PUSH1 0 PUSH1 0 REVERT
	} // 17 bytes
	word := make([]byte, 32)
	copy(word, initcode)
	bc = append(bc, 0x7f)       // PUSH32
	bc = append(bc, word...)    // initcode
	bc = append(bc, 0x5f)       // PUSH0 offset
	bc = append(bc, 0x52)       // MSTORE
	bc = append(bc, 0x60, 0x11) // PUSH1 17 size
	bc = append(bc, 0x5f)       // PUSH0 offset
	bc = append(bc, 0x5f)       // PUSH0 value (0 endowment for the deploy)
	bc = append(bc, 0xf0)       // CREATE -> address  (stack: [addr])
	// CALL it with value; it reverts, rolling back the transfer log.
	bc = append(bc, 0x60, 0x00)               // retSize
	bc = append(bc, 0x60, 0x00)               // retOffset
	bc = append(bc, 0x60, 0x00)               // argsSize
	bc = append(bc, 0x60, 0x00)               // argsOffset
	bc = append(bc, g.pushTransferValue()...) // value
	bc = append(bc, 0x85)                     // DUP6 -> copy the CREATE address
	bc = append(bc, 0x61, 0x27, 0x10)         // PUSH2 10000 gas
	bc = append(bc, 0xf1)                     // CALL -> success flag (0)
	bc = append(bc, 0x50)                     // POP flag, leaving addr as +1 result
	return bc
}

// generateSelfdestructSweep exercises EIP-8246 burn/credit behavior across
// beneficiary classes, in both normal and same-tx-created contexts (0x125,
// StackOutput=0).
func (g *OpcodeGenerator) generateSelfdestructSweep() []byte {
	if g.rng.Float64() < 0.5 {
		return g.emitInlineSelfdestruct()
	}
	return g.emitCreatedSelfdestruct()
}

// emitInlineSelfdestruct: SELFDESTRUCT in the current (normal) frame.
// Pushes a chosen beneficiary then SELFDESTRUCT consumes it. Net: 0.
func (g *OpcodeGenerator) emitInlineSelfdestruct() []byte {
	var bc []byte
	bc = append(bc, g.pushBeneficiary()...) // +1 beneficiary
	bc = append(bc, 0xff)                   // SELFDESTRUCT (consumes 1, halts)
	return bc
}

// emitCreatedSelfdestruct: CREATE a sub-contract whose INITCODE selfdestructs
// to a chosen beneficiary (same-tx-created context for EIP-8246). Net: 0.
func (g *OpcodeGenerator) emitCreatedSelfdestruct() []byte {
	var bc []byte
	ben := g.rng.Bytes(20)
	if g.rng.Float64() < 0.4 {
		ben = make([]byte, 20) // 40%: address(0) - known-nonexistent class
	}
	// initcode (runs at deploy): 73 <20-byte ben> ff  => 22 bytes
	initcode := append([]byte{0x73}, ben...)
	initcode = append(initcode, 0xff)
	word := make([]byte, 32)
	copy(word, initcode)
	bc = append(bc, 0x7f)                     // PUSH32
	bc = append(bc, word...)                  // initcode word
	bc = append(bc, 0x5f)                     // PUSH0 offset
	bc = append(bc, 0x52)                     // MSTORE
	bc = append(bc, 0x60, 0x16)               // PUSH1 22 size
	bc = append(bc, 0x5f)                     // PUSH0 offset
	bc = append(bc, g.pushTransferValue()...) // value endowment to the doomed contract
	bc = append(bc, 0xf0)                     // CREATE -> address (or 0)
	bc = append(bc, 0x50)                     // POP the result
	return bc
}
