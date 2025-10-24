package evmfuzz

import (
	"encoding/binary"
)

// processPrecompileResult handles precompile results based on fuzzing mode
// - In precompiles-only mode: adds LOG0 to log results
// - In normal mode: loads results to stack in 32-byte words for further fuzzing
func (g *OpcodeGenerator) processPrecompileResult(memOffset, size int) []byte {
	var bytecode []byte

	if g.fuzzMode == "precompiles" {
		// Precompiles-only mode: LOG the result
		// PUSH the size of the data to log
		if size < 256 {
			bytecode = append(bytecode, 0x60, byte(size)) // PUSH1 size
		} else {
			bytecode = append(bytecode, 0x61, byte(size>>8), byte(size)) // PUSH2 size
		}

		// PUSH the memory offset where the data is stored
		if memOffset < 256 {
			bytecode = append(bytecode, 0x60, byte(memOffset)) // PUSH1 offset
		} else {
			bytecode = append(bytecode, 0x61, byte(memOffset>>8), byte(memOffset)) // PUSH2 offset
		}

		// LOG0 opcode to log the precompile result
		bytecode = append(bytecode, 0xa0) // LOG0
	} else {
		// Normal fuzzing mode: load result data to stack for further processing
		// Load result in 32-byte words using MLOAD
		numWords := (size + 31) / 32 // Round up to get number of 32-byte words

		// Load each 32-byte word from memory to stack
		for i := 0; i < numWords; i++ {
			wordOffset := memOffset + (i * 32)

			// PUSH the memory offset for this word
			if wordOffset < 256 {
				bytecode = append(bytecode, 0x60, byte(wordOffset)) // PUSH1 offset
			} else {
				bytecode = append(bytecode, 0x61, byte(wordOffset>>8), byte(wordOffset)) // PUSH2 offset
			}

			// MLOAD to load 32 bytes from memory to stack
			bytecode = append(bytecode, 0x51) // MLOAD
		}
	}

	return bytecode
}

// generateKZGPointEvalCall creates a KZG point evaluation precompile call
func (g *OpcodeGenerator) generateKZGPointEvalCall() []byte {
	var bytecode []byte

	// KZG Point Evaluation input format:
	// [0:32]   - versioned_hash (32 bytes)
	// [32:64]  - z (field element, 32 bytes)
	// [64:96]  - y (field element, 32 bytes)
	// [96:144] - commitment (48 bytes)
	// [144:192] - proof (48 bytes)
	// Total: 192 bytes

	// Generate versioned hash (32 bytes)
	versionedHash := g.rng.Bytes(32)
	// Set version byte to 0x01 for KZG commitments
	versionedHash[0] = 0x01
	bytecode = append(bytecode, 0x7f) // PUSH32
	bytecode = append(bytecode, versionedHash...)
	bytecode = append(bytecode, 0x60, 0x00) // PUSH1 0
	bytecode = append(bytecode, 0x52)       // MSTORE

	// Generate z field element (32 bytes)
	// Use BLS12-381 scalar field modulus for validity
	zFieldElement := g.generateValidBLS12ScalarField()
	bytecode = append(bytecode, 0x7f) // PUSH32
	bytecode = append(bytecode, zFieldElement...)
	bytecode = append(bytecode, 0x60, 0x20) // PUSH1 32
	bytecode = append(bytecode, 0x52)       // MSTORE

	// Generate y field element (32 bytes)
	yFieldElement := g.generateValidBLS12ScalarField()
	bytecode = append(bytecode, 0x7f) // PUSH32
	bytecode = append(bytecode, yFieldElement...)
	bytecode = append(bytecode, 0x60, 0x40) // PUSH1 64
	bytecode = append(bytecode, 0x52)       // MSTORE

	// Generate commitment (48 bytes = G1 point in compressed form)
	commitment := g.generateValidKZGCommitment()
	// Write first 32 bytes of commitment
	bytecode = append(bytecode, 0x7f) // PUSH32
	bytecode = append(bytecode, commitment[:32]...)
	bytecode = append(bytecode, 0x60, 0x60) // PUSH1 96
	bytecode = append(bytecode, 0x52)       // MSTORE

	// Write remaining 16 bytes of commitment (padded to 32 bytes)
	commitmentPart2 := make([]byte, 32)
	copy(commitmentPart2, commitment[32:])
	bytecode = append(bytecode, 0x7f) // PUSH32
	bytecode = append(bytecode, commitmentPart2...)
	bytecode = append(bytecode, 0x60, 0x80) // PUSH1 128
	bytecode = append(bytecode, 0x52)       // MSTORE

	// Generate proof (48 bytes = G1 point in compressed form)
	proof := g.generateValidKZGCommitment() // Same format as commitment
	// Write first 32 bytes of proof
	bytecode = append(bytecode, 0x7f) // PUSH32
	bytecode = append(bytecode, proof[:32]...)
	bytecode = append(bytecode, 0x60, 0xa0) // PUSH1 160
	bytecode = append(bytecode, 0x52)       // MSTORE

	// Write remaining 16 bytes of proof (padded to 32 bytes)
	proofPart2 := make([]byte, 32)
	copy(proofPart2, proof[32:])
	bytecode = append(bytecode, 0x7f) // PUSH32
	bytecode = append(bytecode, proofPart2...)
	bytecode = append(bytecode, 0x60, 0xc0) // PUSH1 192
	bytecode = append(bytecode, 0x52)       // MSTORE

	// Setup CALL to KZG point evaluation precompile
	bytecode = append(bytecode, 0x60, 0x40) // PUSH1 64 (return size - field element + success flag)
	bytecode = append(bytecode, 0x60, 0xe0) // PUSH1 224 (return offset)
	bytecode = append(bytecode, 0x60, 0xc0) // PUSH1 192 (args size)
	bytecode = append(bytecode, 0x60, 0x00) // PUSH1 0 (args offset)
	bytecode = append(bytecode, 0x60, 0x00) // PUSH1 0 (value)
	bytecode = append(bytecode, 0x60, 0x0a) // PUSH1 10 (precompile address)
	bytecode = append(bytecode, 0x5a)       // GAS
	bytecode = append(bytecode, 0xf1)       // CALL

	// Add LOG0 to log the result when in precompiles-only mode (return at 0xe0, size 64)
	bytecode = append(bytecode, g.processPrecompileResult(0xe0, 64)...)

	return bytecode
}

// generateValidBLS12ScalarField generates a valid BLS12-381 scalar field element
func (g *OpcodeGenerator) generateValidBLS12ScalarField() []byte {
	// BLS12-381 scalar field modulus (order of the curve)
	// 0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001
	scalar := make([]byte, 32)

	// Generate patterns that are likely to be valid
	patternType := g.rng.Intn(4)
	switch patternType {
	case 0:
		// Small values (often useful for testing)
		smallVal := g.rng.Intn(1000)
		binary.BigEndian.PutUint64(scalar[24:], uint64(smallVal))
	case 1:
		// Powers of 2 (common edge cases)
		power := g.rng.Intn(31)
		binary.BigEndian.PutUint64(scalar[24:], 1<<power)
	case 2:
		// Random but constrained to be less than modulus
		// Use simple constraint: set high bit to 0
		copy(scalar, g.rng.Bytes(32))
		scalar[0] &= 0x7F // Clear high bit
	default:
		// Completely random (let precompile handle modular reduction)
		copy(scalar, g.rng.Bytes(32))
	}

	return scalar
}

// generateValidKZGCommitment generates a valid KZG commitment (compressed G1 point)
func (g *OpcodeGenerator) generateValidKZGCommitment() []byte {
	commitment := make([]byte, 48)

	// 10% chance of identity/zero commitment
	if g.rng.Float64() < 0.1 {
		// Zero commitment (point at infinity in compressed form)
		commitment[0] = 0xc0 // Compressed point at infinity marker
		return commitment
	}

	// Generate valid compressed G1 point
	// For BLS12-381, compressed points are 48 bytes with specific format
	patternType := g.rng.Intn(3)
	switch patternType {
	case 0:
		// Use generator point pattern
		commitment[0] = 0x97 // Compressed point marker + y-coordinate bit
		// Fill with known-good pattern
		for i := 1; i < 48; i++ {
			commitment[i] = byte(i % 256)
		}
	case 1:
		// Small coordinate values
		commitment[0] = 0x80 // Compressed point marker
		binary.BigEndian.PutUint64(commitment[40:48], uint64(g.rng.Intn(1000)+1))
	default:
		// Random but constrained
		copy(commitment[1:], g.rng.Bytes(47))
		// Set compression flag and ensure valid format
		if g.rng.Float64() < 0.5 {
			commitment[0] = 0x80 // Compressed, y-coord even
		} else {
			commitment[0] = 0xa0 // Compressed, y-coord odd
		}
	}

	return commitment
}

// generateModexpCall creates a special MODEXP precompile call with proper dynamic length format
func (g *OpcodeGenerator) generateModexpCall() []byte {
	var bytecode []byte

	// MODEXP input format:
	// [0:32]   - base_len (big-endian 32-byte integer)
	// [32:64]  - exp_len (big-endian 32-byte integer)
	// [64:96]  - mod_len (big-endian 32-byte integer)
	// [96:96+base_len] - base (big-endian integer)
	// [96+base_len:96+base_len+exp_len] - exponent (big-endian integer)
	// [96+base_len+exp_len:96+base_len+exp_len+mod_len] - modulus (big-endian integer)

	// Generate reasonable lengths (between 1 and 96 bytes each)
	baseLen := 1 + g.rng.Intn(95) // 1-96 bytes
	expLen := 1 + g.rng.Intn(95)  // 1-96 bytes
	modLen := 1 + g.rng.Intn(95)  // 1-96 bytes

	totalLen := 96 + baseLen + expLen + modLen

	// Write base_len at offset 0 (32 bytes, big-endian)
	baseLenBytes := make([]byte, 32)
	baseLenBytes[31] = byte(baseLen)  // Simple encoding for small values
	bytecode = append(bytecode, 0x7f) // PUSH32
	bytecode = append(bytecode, baseLenBytes...)
	bytecode = append(bytecode, 0x60, 0x00) // PUSH1 0 (offset)
	bytecode = append(bytecode, 0x52)       // MSTORE

	// Write exp_len at offset 32
	expLenBytes := make([]byte, 32)
	expLenBytes[31] = byte(expLen)
	bytecode = append(bytecode, 0x7f) // PUSH32
	bytecode = append(bytecode, expLenBytes...)
	bytecode = append(bytecode, 0x60, 0x20) // PUSH1 32 (offset)
	bytecode = append(bytecode, 0x52)       // MSTORE

	// Write mod_len at offset 64
	modLenBytes := make([]byte, 32)
	modLenBytes[31] = byte(modLen)
	bytecode = append(bytecode, 0x7f) // PUSH32
	bytecode = append(bytecode, modLenBytes...)
	bytecode = append(bytecode, 0x60, 0x40) // PUSH1 64 (offset)
	bytecode = append(bytecode, 0x52)       // MSTORE

	// Write base data starting at offset 96
	baseData := g.rng.Bytes(baseLen)
	for i := 0; i < baseLen; i += 32 {
		remaining := baseLen - i
		if remaining >= 32 {
			// Write full 32 bytes
			bytecode = append(bytecode, 0x7f) // PUSH32
			bytecode = append(bytecode, baseData[i:i+32]...)
		} else {
			// Write partial bytes using appropriate PUSH
			chunk := make([]byte, 32)
			copy(chunk, baseData[i:])
			bytecode = append(bytecode, 0x7f) // PUSH32
			bytecode = append(bytecode, chunk...)
		}

		offset := 96 + i
		if offset < 256 {
			bytecode = append(bytecode, 0x60, byte(offset)) // PUSH1 offset
		} else {
			bytecode = append(bytecode, 0x61, byte(offset>>8), byte(offset)) // PUSH2 offset
		}
		bytecode = append(bytecode, 0x52) // MSTORE
	}

	// Write exponent data
	expData := g.rng.Bytes(expLen)
	for i := 0; i < expLen; i += 32 {
		remaining := expLen - i
		if remaining >= 32 {
			bytecode = append(bytecode, 0x7f) // PUSH32
			bytecode = append(bytecode, expData[i:i+32]...)
		} else {
			chunk := make([]byte, 32)
			copy(chunk, expData[i:])
			bytecode = append(bytecode, 0x7f) // PUSH32
			bytecode = append(bytecode, chunk...)
		}

		offset := 96 + baseLen + i
		if offset < 256 {
			bytecode = append(bytecode, 0x60, byte(offset)) // PUSH1 offset
		} else {
			bytecode = append(bytecode, 0x61, byte(offset>>8), byte(offset)) // PUSH2 offset
		}
		bytecode = append(bytecode, 0x52) // MSTORE
	}

	// Write modulus data
	modData := g.rng.Bytes(modLen)
	for i := 0; i < modLen; i += 32 {
		remaining := modLen - i
		if remaining >= 32 {
			bytecode = append(bytecode, 0x7f) // PUSH32
			bytecode = append(bytecode, modData[i:i+32]...)
		} else {
			chunk := make([]byte, 32)
			copy(chunk, modData[i:])
			bytecode = append(bytecode, 0x7f) // PUSH32
			bytecode = append(bytecode, chunk...)
		}

		offset := 96 + baseLen + expLen + i
		if offset < 256 {
			bytecode = append(bytecode, 0x60, byte(offset)) // PUSH1 offset
		} else {
			bytecode = append(bytecode, 0x61, byte(offset>>8), byte(offset)) // PUSH2 offset
		}
		bytecode = append(bytecode, 0x52) // MSTORE
	}

	// Now setup the CALL to MODEXP precompile
	// CALL(gas, address, value, argsOffset, argsSize, retOffset, retSize)
	// Push in reverse order so gas is at top of stack (position 0)

	// Return data size (equal to modulus length)
	if modLen < 256 {
		bytecode = append(bytecode, 0x60, byte(modLen)) // PUSH1 modLen
	} else {
		bytecode = append(bytecode, 0x61, byte(modLen>>8), byte(modLen)) // PUSH2 modLen
	}

	// Return data offset (store return after input data)
	returnOffset := totalLen + 32
	if returnOffset < 256 {
		bytecode = append(bytecode, 0x60, byte(returnOffset)) // PUSH1 retOffset
	} else {
		bytecode = append(bytecode, 0x61, byte(returnOffset>>8), byte(returnOffset)) // PUSH2 retOffset
	}

	// Arguments size (total length)
	if totalLen < 256 {
		bytecode = append(bytecode, 0x60, byte(totalLen)) // PUSH1 totalLen
	} else {
		bytecode = append(bytecode, 0x61, byte(totalLen>>8), byte(totalLen)) // PUSH2 totalLen
	}

	// Arguments offset (0 - we stored data starting at memory 0)
	bytecode = append(bytecode, 0x60, 0x00) // PUSH1 0

	// Value (always 0 for precompiles)
	bytecode = append(bytecode, 0x60, 0x00) // PUSH1 0

	// Precompile address (5)
	bytecode = append(bytecode, 0x60, 0x05) // PUSH1 5

	// Gas for the call - 90% use GAS opcode, 10% use fixed value
	if g.rng.Float64() < 0.9 {
		bytecode = append(bytecode, 0x5a) // GAS opcode - use remaining gas
	} else {
		bytecode = append(bytecode, 0x61, 0x13, 0x88) // PUSH2 5000
	}

	// CALL opcode
	bytecode = append(bytecode, 0xf1) // CALL

	// Add LOG0 to log the result when in precompiles-only mode (return at returnOffset, size modLen)
	bytecode = append(bytecode, g.processPrecompileResult(returnOffset, modLen)...)

	return bytecode
}

// generateEcrecoverCall creates an ECRECOVER precompile call with specialized input generation
func (g *OpcodeGenerator) generateEcrecoverCall() []byte {
	var bytecode []byte

	// ECRECOVER takes 128 bytes: hash(32) + v(32) + r(32) + s(32)
	// Returns 32 bytes (recovered address or zero)

	// Generate message hash (32 bytes)
	hash := g.rng.Bytes(32)
	hash = g.transformer.TransformPrecompileInput(hash, 0x01)
	bytecode = append(bytecode, 0x7f) // PUSH32
	bytecode = append(bytecode, hash...)
	bytecode = append(bytecode, 0x60, 0x00) // PUSH1 0 (offset)
	bytecode = append(bytecode, 0x52)       // MSTORE

	// Generate v parameter (recovery ID: 27 or 28, or invalid values for testing)
	var v [32]byte
	if g.rng.Float64() < 0.01 { // 1% chance of invalid v
		invalidV := g.rng.Intn(255)
		if invalidV == 27 || invalidV == 28 {
			invalidV = 29 // Force invalid
		}
		v[31] = byte(invalidV)
	} else {
		v[31] = byte(27 + g.rng.Intn(2)) // Valid: 27 or 28
	}
	bytecode = append(bytecode, 0x7f) // PUSH32
	bytecode = append(bytecode, v[:]...)
	bytecode = append(bytecode, 0x60, 0x20) // PUSH1 32 (offset)
	bytecode = append(bytecode, 0x52)       // MSTORE

	// Generate r parameter (32 bytes)
	r := g.rng.Bytes(32)
	bytecode = append(bytecode, 0x7f) // PUSH32
	bytecode = append(bytecode, r...)
	bytecode = append(bytecode, 0x60, 0x40) // PUSH1 64 (offset)
	bytecode = append(bytecode, 0x52)       // MSTORE

	// Generate s parameter (32 bytes)
	s := g.rng.Bytes(32)
	bytecode = append(bytecode, 0x7f) // PUSH32
	bytecode = append(bytecode, s...)
	bytecode = append(bytecode, 0x60, 0x60) // PUSH1 96 (offset)
	bytecode = append(bytecode, 0x52)       // MSTORE

	// Setup CALL to ECRECOVER precompile
	bytecode = append(bytecode, 0x60, 0x20) // PUSH1 32 (return size)
	bytecode = append(bytecode, 0x60, 0x80) // PUSH1 128 (return offset)
	bytecode = append(bytecode, 0x60, 0x80) // PUSH1 128 (args size)
	bytecode = append(bytecode, 0x60, 0x00) // PUSH1 0 (args offset)
	bytecode = append(bytecode, 0x60, 0x00) // PUSH1 0 (value)
	bytecode = append(bytecode, 0x60, 0x01) // PUSH1 1 (precompile address)
	bytecode = append(bytecode, 0x5a)       // GAS
	bytecode = append(bytecode, 0xf1)       // CALL

	// Push result to stack - load return value from memory
	bytecode = append(bytecode, 0x60, 0x80) // PUSH1 128 (return offset)
	bytecode = append(bytecode, 0x51)       // MLOAD

	// Add LOG0 to log the result when in precompiles-only mode
	bytecode = append(bytecode, g.processPrecompileResult(0x80, 32)...)

	return bytecode
}

// generateSha256Call creates a SHA256 precompile call with variable-length input
func (g *OpcodeGenerator) generateSha256Call() []byte {
	var bytecode []byte

	// Generate variable-length input data (1-1024 bytes)
	inputLen := g.rng.Intn(1024) + 1
	inputData := g.rng.Bytes(inputLen)
	// Apply transformations for edge case testing
	inputData = g.transformer.TransformPrecompileInput(inputData, 0x02)

	// Update inputLen to reflect the actual transformed data size
	actualInputLen := len(inputData)

	// Write input data to memory in 32-byte chunks
	memOffset := 0
	for i := 0; i < actualInputLen; i += 32 {
		chunk := make([]byte, 32)
		chunkLen := actualInputLen - i
		if chunkLen > 32 {
			chunkLen = 32
		}
		if i+chunkLen <= len(inputData) {
			copy(chunk, inputData[i:i+chunkLen])
		}

		bytecode = append(bytecode, 0x7f) // PUSH32
		bytecode = append(bytecode, chunk...)
		bytecode = append(bytecode, 0x60, byte(memOffset)) // PUSH1 offset
		bytecode = append(bytecode, 0x52)                  // MSTORE
		memOffset += 32
	}

	// Setup CALL to SHA256 precompile
	returnOffset := memOffset + 32
	bytecode = append(bytecode, 0x60, 0x20)                                          // PUSH1 32 (return size)
	bytecode = append(bytecode, 0x60, byte(returnOffset))                            // PUSH1 returnOffset
	bytecode = append(bytecode, 0x61, byte(actualInputLen>>8), byte(actualInputLen)) // PUSH2 actualInputLen
	bytecode = append(bytecode, 0x60, 0x00)                                          // PUSH1 0 (args offset)
	bytecode = append(bytecode, 0x60, 0x00)                                          // PUSH1 0 (value)
	bytecode = append(bytecode, 0x60, 0x02)                                          // PUSH1 2 (precompile address)
	bytecode = append(bytecode, 0x5a)                                                // GAS
	bytecode = append(bytecode, 0xf1)                                                // CALL

	// Push result to stack
	bytecode = append(bytecode, 0x60, byte(returnOffset)) // PUSH1 returnOffset
	bytecode = append(bytecode, 0x51)                     // MLOAD

	// Add LOG0 to log the result when in precompiles-only mode
	bytecode = append(bytecode, g.processPrecompileResult(returnOffset, 32)...)

	return bytecode
}

// generateRipemd160Call creates a RIPEMD160 precompile call
func (g *OpcodeGenerator) generateRipemd160Call() []byte {
	var bytecode []byte

	// Generate variable-length input data (1-512 bytes)
	inputLen := g.rng.Intn(512) + 1
	inputData := g.rng.Bytes(inputLen)

	// Write input data to memory in 32-byte chunks
	memOffset := 0
	for i := 0; i < inputLen; i += 32 {
		chunk := make([]byte, 32)
		chunkLen := inputLen - i
		if chunkLen > 32 {
			chunkLen = 32
		}
		copy(chunk, inputData[i:i+chunkLen])

		bytecode = append(bytecode, 0x7f) // PUSH32
		bytecode = append(bytecode, chunk...)
		bytecode = append(bytecode, 0x60, byte(memOffset)) // PUSH1 offset
		bytecode = append(bytecode, 0x52)                  // MSTORE
		memOffset += 32
	}

	// Setup CALL to RIPEMD160 precompile
	returnOffset := memOffset + 32
	bytecode = append(bytecode, 0x60, 0x20)                              // PUSH1 32 (return size)
	bytecode = append(bytecode, 0x60, byte(returnOffset))                // PUSH1 returnOffset
	bytecode = append(bytecode, 0x61, byte(inputLen>>8), byte(inputLen)) // PUSH2 inputLen
	bytecode = append(bytecode, 0x60, 0x00)                              // PUSH1 0 (args offset)
	bytecode = append(bytecode, 0x60, 0x00)                              // PUSH1 0 (value)
	bytecode = append(bytecode, 0x60, 0x03)                              // PUSH1 3 (precompile address)
	bytecode = append(bytecode, 0x5a)                                    // GAS
	bytecode = append(bytecode, 0xf1)                                    // CALL

	// Push result to stack
	bytecode = append(bytecode, 0x60, byte(returnOffset)) // PUSH1 returnOffset
	bytecode = append(bytecode, 0x51)                     // MLOAD

	// Add LOG0 to log the result when in precompiles-only mode
	bytecode = append(bytecode, g.processPrecompileResult(returnOffset, 32)...)

	return bytecode
}

// generateIdentityCall creates an IDENTITY precompile call
func (g *OpcodeGenerator) generateIdentityCall() []byte {
	var bytecode []byte

	// Generate variable-length input data (0-2048 bytes)
	inputLen := g.rng.Intn(2048)
	inputData := g.rng.Bytes(inputLen)

	// Write input data to memory in 32-byte chunks
	memOffset := 0
	if inputLen > 0 {
		for i := 0; i < inputLen; i += 32 {
			chunk := make([]byte, 32)
			chunkLen := inputLen - i
			if chunkLen > 32 {
				chunkLen = 32
			}
			copy(chunk, inputData[i:i+chunkLen])

			bytecode = append(bytecode, 0x7f) // PUSH32
			bytecode = append(bytecode, chunk...)
			bytecode = append(bytecode, 0x60, byte(memOffset)) // PUSH1 offset
			bytecode = append(bytecode, 0x52)                  // MSTORE
			memOffset += 32
		}
	}

	// Setup CALL to IDENTITY precompile
	returnOffset := memOffset + 32
	returnSize := inputLen // Identity returns same size as input
	if returnSize == 0 {
		returnSize = 32 // At least push one word to stack
	}

	bytecode = append(bytecode, 0x61, byte(returnSize>>8), byte(returnSize)) // PUSH2 returnSize
	bytecode = append(bytecode, 0x60, byte(returnOffset))                    // PUSH1 returnOffset
	bytecode = append(bytecode, 0x61, byte(inputLen>>8), byte(inputLen))     // PUSH2 inputLen
	bytecode = append(bytecode, 0x60, 0x00)                                  // PUSH1 0 (args offset)
	bytecode = append(bytecode, 0x60, 0x00)                                  // PUSH1 0 (value)
	bytecode = append(bytecode, 0x60, 0x04)                                  // PUSH1 4 (precompile address)
	bytecode = append(bytecode, 0x5a)                                        // GAS
	bytecode = append(bytecode, 0xf1)                                        // CALL

	// Push result to stack
	bytecode = append(bytecode, 0x60, byte(returnOffset)) // PUSH1 returnOffset
	bytecode = append(bytecode, 0x51)                     // MLOAD

	// Add LOG0 to log the result when in precompiles-only mode
	bytecode = append(bytecode, g.processPrecompileResult(returnOffset, returnSize)...)

	return bytecode
}

// generateBlake2fCall creates a BLAKE2F precompile call
func (g *OpcodeGenerator) generateBlake2fCall() []byte {
	var bytecode []byte

	// BLAKE2F requires exactly 213 bytes input:
	// - rounds (4 bytes)
	// - h (64 bytes) - state vector
	// - m (128 bytes) - message block
	// - t (16 bytes) - offset counters
	// - f (1 byte) - final block indicator

	// Generate rounds (4 bytes) - usually a reasonable number
	rounds := make([]byte, 32)
	if g.rng.Float64() < 0.05 { // 5% chance of extreme values
		binary.BigEndian.PutUint32(rounds[28:], uint32(g.rng.Intn(0xFFFFFF)))
	} else {
		binary.BigEndian.PutUint32(rounds[28:], uint32(g.rng.Intn(100)+1))
	}
	bytecode = append(bytecode, 0x7f) // PUSH32
	bytecode = append(bytecode, rounds...)
	bytecode = append(bytecode, 0x60, 0x00) // PUSH1 0
	bytecode = append(bytecode, 0x52)       // MSTORE

	// Generate h (state vector) - 64 bytes in 2 chunks
	h1 := g.rng.Bytes(32)
	h2 := g.rng.Bytes(32)
	bytecode = append(bytecode, 0x7f) // PUSH32
	bytecode = append(bytecode, h1...)
	bytecode = append(bytecode, 0x60, 0x04) // PUSH1 4 (offset)
	bytecode = append(bytecode, 0x52)       // MSTORE
	bytecode = append(bytecode, 0x7f)       // PUSH32
	bytecode = append(bytecode, h2...)
	bytecode = append(bytecode, 0x60, 0x24) // PUSH1 36 (offset)
	bytecode = append(bytecode, 0x52)       // MSTORE

	// Generate m (message block) - 128 bytes in 4 chunks
	for i := 0; i < 4; i++ {
		m := g.rng.Bytes(32)
		bytecode = append(bytecode, 0x7f) // PUSH32
		bytecode = append(bytecode, m...)
		offset := 0x44 + i*0x20
		bytecode = append(bytecode, 0x60, byte(offset)) // PUSH1 offset
		bytecode = append(bytecode, 0x52)               // MSTORE
	}

	// Generate t (offset counters) - 16 bytes
	t := make([]byte, 32)
	copy(t[16:], g.rng.Bytes(16))
	bytecode = append(bytecode, 0x7f) // PUSH32
	bytecode = append(bytecode, t...)
	bytecode = append(bytecode, 0x60, 0xc4) // PUSH1 196 (offset)
	bytecode = append(bytecode, 0x52)       // MSTORE

	// Generate f (final block indicator) - 1 byte
	f := make([]byte, 32)
	f[31] = byte(g.rng.Intn(2))       // 0 or 1
	bytecode = append(bytecode, 0x7f) // PUSH32
	bytecode = append(bytecode, f...)
	bytecode = append(bytecode, 0x60, 0xd4) // PUSH1 212 (offset)
	bytecode = append(bytecode, 0x52)       // MSTORE

	// Setup CALL to BLAKE2F precompile
	bytecode = append(bytecode, 0x60, 0x40) // PUSH1 64 (return size)
	bytecode = append(bytecode, 0x60, 0xf5) // PUSH1 245 (return offset)
	bytecode = append(bytecode, 0x60, 0xd5) // PUSH1 213 (args size)
	bytecode = append(bytecode, 0x60, 0x00) // PUSH1 0 (args offset)
	bytecode = append(bytecode, 0x60, 0x00) // PUSH1 0 (value)
	bytecode = append(bytecode, 0x60, 0x09) // PUSH1 9 (precompile address)
	bytecode = append(bytecode, 0x5a)       // GAS
	bytecode = append(bytecode, 0xf1)       // CALL

	// Push first 32 bytes of result to stack
	bytecode = append(bytecode, 0x60, 0xf5) // PUSH1 245 (return offset)
	bytecode = append(bytecode, 0x51)       // MLOAD

	// Push second 32 bytes of result to stack
	bytecode = append(bytecode, 0x61, 0x01, 0x15) // PUSH2 277 (return offset + 32)
	bytecode = append(bytecode, 0x51)             // MLOAD

	// Add LOG0 to log the result when in precompiles-only mode
	bytecode = append(bytecode, g.processPrecompileResult(0xf5, 64)...)

	return bytecode
}
