package evmfuzz

import (
	"encoding/binary"
)

// EIP-8282 builder-request / EIP-7997 CREATE2-factory system-contract predeploys.
var (
	builderDepositAddr = [20]byte{0x00, 0x00, 0x88, 0x4d, 0x2a, 0xa3, 0x2e, 0xaa, 0x15, 0x5f, 0x59, 0xa2, 0xf2, 0x4e, 0xfa, 0x73, 0xd9, 0x00, 0x82, 0x82}
	builderExitAddr    = [20]byte{0x00, 0x00, 0x14, 0x57, 0x4a, 0x74, 0xc8, 0x05, 0x59, 0x0a, 0xff, 0x94, 0x99, 0xfc, 0x7a, 0x69, 0x0f, 0x00, 0x82, 0x82}
	create2FactoryAddr = [20]byte{0x4e, 0x59, 0xb4, 0x48, 0x47, 0xb3, 0x79, 0x57, 0x85, 0x88, 0x92, 0x0c, 0xa7, 0x8f, 0xbf, 0x26, 0xc0, 0xb4, 0x95, 0x6c}
)

// isSystemOpcode reports whether an opcode is an EIP-8282/7997 system-contract call.
func (g *OpcodeGenerator) isSystemOpcode(opcode uint16) bool {
	return opcode >= 0x120 && opcode <= 0x122
}

// pushAddress emits PUSH20 <addr>; the address is the literal predeploy, never masked.
func (g *OpcodeGenerator) pushAddress(addr [20]byte) []byte {
	return append([]byte{0x73}, addr[:]...)
}

// mstoreWords emits PUSH32 word; PUSH off; MSTORE for each 32-byte word of data
// starting at baseOffset. len(data) must be a multiple of 32.
func (g *OpcodeGenerator) mstoreWords(data []byte, baseOffset int) []byte {
	var bc []byte
	for i := 0; i+32 <= len(data); i += 32 {
		bc = append(bc, 0x7f)            // PUSH32
		bc = append(bc, data[i:i+32]...) // word
		off := baseOffset + i
		if off < 256 {
			bc = append(bc, 0x60, byte(off)) // PUSH1 off
		} else {
			bc = append(bc, 0x61, byte(off>>8), byte(off)) // PUSH2 off
		}
		bc = append(bc, 0x52) // MSTORE
	}
	return bc
}

// generateBuilderDepositCall emits a CALL to the EIP-8282 builder deposit system
// contract (request type 0x03). Calldata = pubkey(48) ++ withdrawal_credentials(32)
// ++ amount(8 BE gwei) ++ signature(96) = 184 bytes. Most calls revert (excess
// inhibitor / underpayment / amount < 1 ETH), which is a valid fuzz path.
func (g *OpcodeGenerator) generateBuilderDepositCall() []byte {
	var bytecode []byte

	const payloadLen = 184
	payload := make([]byte, 192) // padded to 6 words for MSTORE; last 8 bytes unused

	// pubkey [0:48] — opaque BLS pubkey, 10% all-zero edge case.
	if g.rng.Float64() >= 0.1 {
		copy(payload[0:48], g.rng.Bytes(48))
	}

	// withdrawal_credentials [48:80] — 90% well-formed 0x01/0x02 prefix.
	wc := g.rng.Bytes(32)
	if g.rng.Float64() < 0.9 {
		wc[0] = byte(0x01 + g.rng.Intn(2))
		for i := 1; i < 12; i++ {
			wc[i] = 0x00
		}
	}
	copy(payload[48:80], wc)

	// amount [80:88] — 8 bytes big-endian gwei; mix valid and edge values.
	var amount uint64
	switch g.rng.Intn(4) {
	case 0:
		amount = 0 // below-minimum edge
	case 1:
		amount = 1_000_000_000 // exactly 1 ETH (minimum)
	case 2:
		amount = 32_000_000_000 // 32 ETH (typical)
	default:
		amount = g.rng.Uint64() // fully random
	}
	binary.BigEndian.PutUint64(payload[80:88], amount)

	// signature [88:184] — 96 opaque bytes.
	copy(payload[88:184], g.rng.Bytes(96))

	bytecode = append(bytecode, g.mstoreWords(payload, 0)...)

	const retOffset = 192
	bytecode = append(bytecode, 0x60, 0x20)       // PUSH1 32  (retSize)
	bytecode = append(bytecode, 0x60, retOffset)  // PUSH1 192 (retOffset)
	bytecode = append(bytecode, 0x60, payloadLen) // PUSH1 184 (argsSize)
	bytecode = append(bytecode, 0x60, 0x00)       // PUSH1 0   (argsOffset)
	if g.rng.Float64() < 0.5 {
		bytecode = append(bytecode, 0x34) // CALLVALUE (forward our value)
	} else {
		bytecode = append(bytecode, 0x5f) // PUSH0 (value 0)
	}
	bytecode = append(bytecode, g.pushAddress(builderDepositAddr)...) // PUSH20 addr
	bytecode = append(bytecode, 0x5a)                                 // GAS
	bytecode = append(bytecode, 0xf1)                                 // CALL
	bytecode = append(bytecode, 0x50)                                 // POP (drop success flag)

	bytecode = append(bytecode, g.processPrecompileResult(retOffset, 32)...)
	return bytecode
}

// generateBuilderExitCall emits a CALL to the EIP-8282 builder exit system contract
// (request type 0x04). Calldata = pubkey(48); 30% send EMPTY calldata to hit the
// fee-getter path (empty calldata returns the 32-byte fee).
func (g *OpcodeGenerator) generateBuilderExitCall() []byte {
	var bytecode []byte

	argsSize := 48
	if g.rng.Float64() < 0.3 {
		argsSize = 0 // fee-getter path
	} else {
		pubkey := make([]byte, 64) // 48 real + 16 pad => 2 words
		copy(pubkey, g.rng.Bytes(48))
		bytecode = append(bytecode, g.mstoreWords(pubkey, 0)...)
	}

	const retOffset = 64
	bytecode = append(bytecode, 0x60, 0x20)           // PUSH1 32 (retSize)
	bytecode = append(bytecode, 0x60, retOffset)      // PUSH1 64 (retOffset)
	bytecode = append(bytecode, 0x60, byte(argsSize)) // PUSH1 argsSize (0 or 48)
	bytecode = append(bytecode, 0x60, 0x00)           // PUSH1 0  (argsOffset)
	if g.rng.Float64() < 0.5 {
		bytecode = append(bytecode, 0x34) // CALLVALUE (cover fee)
	} else {
		bytecode = append(bytecode, 0x5f) // PUSH0 (value 0)
	}
	bytecode = append(bytecode, g.pushAddress(builderExitAddr)...) // PUSH20 addr
	bytecode = append(bytecode, 0x5a)                              // GAS
	bytecode = append(bytecode, 0xf1)                              // CALL
	bytecode = append(bytecode, 0x50)                              // POP (drop success flag)

	bytecode = append(bytecode, g.processPrecompileResult(retOffset, 32)...)
	return bytecode
}

// generateCreate2FactoryCall emits a CALL to the EIP-7997 deterministic CREATE2
// factory. Calldata = salt(32) ++ initcode(N); returns the 20-byte deployed address.
// 10% send undersized (<32 byte) calldata to exercise the factory revert path; 30%
// of well-formed calls repeat the identical (salt,initcode) to hit the CREATE2
// collision path. Init code is bounded to keep deploy gas/size small.
func (g *OpcodeGenerator) generateCreate2FactoryCall() []byte {
	var bytecode []byte

	const retOffset = 0x100 // 256: clear of the salt+initcode region (<=128)

	// 10%: malformed undersized calldata (< 32 bytes) => factory reverts.
	if g.rng.Float64() < 0.1 {
		short := make([]byte, 32)
		copy(short, g.rng.Bytes(g.rng.Intn(31)))
		argsSize := g.rng.Intn(31) // deliberately < 32
		bytecode = append(bytecode, g.mstoreWords(short, 0)...)
		bytecode = append(bytecode, g.factoryCallTail(byte(argsSize), retOffset)...)
		bytecode = append(bytecode, g.processPrecompileResult(retOffset, 32)...)
		return bytecode
	}

	salt := g.rng.Bytes(32)
	icLen := 1 + g.rng.Intn(64)
	initcode := g.rng.Bytes(icLen)

	totalLen := 32 + icLen
	buf := make([]byte, ((totalLen+31)/32)*32)
	copy(buf[0:32], salt)
	copy(buf[32:totalLen], initcode)
	bytecode = append(bytecode, g.mstoreWords(buf, 0)...)

	bytecode = append(bytecode, g.factoryCallTail(byte(min(totalLen, 255)), retOffset)...)

	// 30%: identical second call to test the (salt,initcode) collision path.
	if g.rng.Float64() < 0.3 {
		bytecode = append(bytecode, g.factoryCallTail(byte(min(totalLen, 255)), retOffset)...)
	}

	bytecode = append(bytecode, g.processPrecompileResult(retOffset, 32)...)
	return bytecode
}

// factoryCallTail emits the CALL to the CREATE2 factory for argsSize bytes at mem[0]
// and pops the success flag. Net stack effect: 0 (the result word is handled once by
// the caller via processPrecompileResult).
func (g *OpcodeGenerator) factoryCallTail(argsSize byte, retOffset int) []byte {
	var bc []byte
	bc = append(bc, 0x60, 0x20)                                // PUSH1 32 (retSize: low 20 bytes = addr)
	bc = append(bc, 0x61, byte(retOffset>>8), byte(retOffset)) // PUSH2 retOffset
	bc = append(bc, 0x60, argsSize)                            // PUSH1 argsSize
	bc = append(bc, 0x60, 0x00)                                // PUSH1 0 (argsOffset)
	if g.rng.Float64() < 0.5 {
		bc = append(bc, 0x34) // CALLVALUE (forward to deployment)
	} else {
		bc = append(bc, 0x5f) // PUSH0 (value 0)
	}
	bc = append(bc, g.pushAddress(create2FactoryAddr)...) // PUSH20 addr
	bc = append(bc, 0x5a)                                 // GAS
	bc = append(bc, 0xf1)                                 // CALL
	bc = append(bc, 0x50)                                 // POP (drop success flag)
	return bc
}
