package evmfuzz

import (
	"math/big"

	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fp"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
)

// generateValidBLS12G1PointUncompressed generates a BLS12-381 G1 point using gnark-crypto library
func (g *OpcodeGenerator) generateValidBLS12G1PointUncompressed() []byte {
	point := make([]byte, 128) // 128 bytes uncompressed format: 64 bytes x + 64 bytes y

	choice := g.rng.Intn(100)

	if choice < 5 {
		// 5% chance to generate invalid G1 points for testing edge cases
		invalidType := g.rng.Intn(2)
		switch invalidType {
		case 0:
			// All 0xFF bytes (definitely invalid)
			for i := range point {
				point[i] = 0xFF
			}
		default:
			// Invalid field elements (> modulus)
			copy(point[:64], g.rng.Bytes(64))
			copy(point[64:], g.rng.Bytes(64))
			point[0] |= 0xE0  // Ensure x >= modulus
			point[64] |= 0xE0 // Ensure y >= modulus
		}
		return g.transformer.TransformPrecompileInput(point, 0x0b)
	}

	if choice < 15 {
		// 10% chance of zero point (point at infinity)
		// Zero point: all zeros represents point at infinity in uncompressed format
		return g.transformer.TransformPrecompileInput(point, 0x0b)
	}

	// 85% of the time: Generate random valid points using gnark-crypto library
	var g1Point bls12381.G1Affine

	// Generate a random scalar (this ensures we get a valid point in the subgroup)
	var scalar fr.Element
	scalarBytes := g.rng.Bytes(32)
	scalar.SetBytes(scalarBytes)

	// Get the generator point and multiply by random scalar
	_, _, g1Gen, _ := bls12381.Generators()

	// Scalar multiplication: scalar * generator (guaranteed to be in correct subgroup)
	g1Point.ScalarMultiplication(&g1Gen, scalar.BigInt(new(big.Int)))

	// Convert to uncompressed bytes format (128 bytes: 64 bytes X + 64 bytes Y)
	// Extract X and Y coordinates with consistent 48-byte field elements
	xBytes := g1Point.X.Bytes()
	yBytes := g1Point.Y.Bytes()

	// Always place field elements at the end of 64-byte slots with proper padding
	// Field elements should be exactly 48 bytes, placed at bytes [16:64] and [80:128]
	if len(xBytes) == 48 {
		copy(point[16:64], xBytes[:])
	} else {
		// Handle variable length by right-aligning
		copy(point[64-len(xBytes):64], xBytes[:])
	}

	if len(yBytes) == 48 {
		copy(point[80:128], yBytes[:])
	} else {
		// Handle variable length by right-aligning
		copy(point[128-len(yBytes):128], yBytes[:])
	}

	return g.transformer.TransformPrecompileInput(point, 0x0b)
}

// generateBLS12G1FieldElements generates raw BLS12-381 G1 field elements (48 bytes each)
func (g *OpcodeGenerator) generateBLS12G1FieldElements() ([]byte, []byte) {
	choice := g.rng.Intn(100)

	if choice < 5 {
		// 5% chance to generate invalid field elements for testing
		invalidX := make([]byte, 48)
		invalidY := make([]byte, 48)
		copy(invalidX, g.rng.Bytes(48))
		copy(invalidY, g.rng.Bytes(48))
		// Ensure they're > modulus
		invalidX[0] |= 0xE0
		invalidY[0] |= 0xE0
		transformedX := g.transformer.TransformPrecompileInput(invalidX, 0x0b)
		transformedY := g.transformer.TransformPrecompileInput(invalidY, 0x0b)
		return transformedX, transformedY
	}

	if choice < 15 {
		// 10% chance of zero point (point at infinity)
		zeroX := make([]byte, 48)
		zeroY := make([]byte, 48)
		transformedX := g.transformer.TransformPrecompileInput(zeroX, 0x0b)
		transformedY := g.transformer.TransformPrecompileInput(zeroY, 0x0b)
		return transformedX, transformedY
	}

	// 85% of the time: Generate random valid points using gnark-crypto library
	var g1Point bls12381.G1Affine

	// Generate a random scalar (this ensures we get a valid point in the subgroup)
	var scalar fr.Element
	scalarBytes := g.rng.Bytes(32)
	scalar.SetBytes(scalarBytes)

	// Get the generator point and multiply by random scalar
	_, _, g1Gen, _ := bls12381.Generators()

	// Scalar multiplication: scalar * generator (guaranteed to be in correct subgroup)
	g1Point.ScalarMultiplication(&g1Gen, scalar.BigInt(new(big.Int)))

	// Extract raw field elements (48 bytes each)
	xBytes := g1Point.X.Bytes()
	yBytes := g1Point.Y.Bytes()

	transformedX := g.transformer.TransformPrecompileInput(xBytes[:], 0x0b)
	transformedY := g.transformer.TransformPrecompileInput(yBytes[:], 0x0b)
	return transformedX, transformedY
}

// generateValidBLS12G2PointUncompressed generates a BLS12-381 G2 point using gnark-crypto library
func (g *OpcodeGenerator) generateValidBLS12G2PointUncompressed() []byte {
	point := make([]byte, 256) // 256 bytes uncompressed format

	choice := g.rng.Intn(100)

	if choice < 5 {
		// 5% chance to generate invalid G2 points for testing edge cases
		invalidType := g.rng.Intn(2)
		switch invalidType {
		case 0:
			// All 0xFF bytes (definitely invalid)
			for i := range point {
				point[i] = 0xFF
			}
		default:
			// Invalid field elements (> modulus)
			copy(point, g.rng.Bytes(256))
			for i := 0; i < 4; i++ {
				point[i*64] |= 0xE0 // Ensure all field elements >= modulus
			}
		}
		return g.transformer.TransformPrecompileInput(point, 0x0d)
	}

	if choice < 15 {
		// 10% chance of zero point (point at infinity)
		// Zero point: all zeros represents point at infinity in uncompressed format
		return g.transformer.TransformPrecompileInput(point, 0x0d)
	}

	// 85% of the time: Generate random valid points using gnark-crypto library
	var g2Point bls12381.G2Affine

	// Generate a random scalar
	var scalar fr.Element
	scalarBytes := g.rng.Bytes(32)
	scalar.SetBytes(scalarBytes)

	// Get the generator point and multiply by random scalar
	_, _, _, g2Gen := bls12381.Generators()

	// Scalar multiplication: scalar * generator
	g2Point.ScalarMultiplication(&g2Gen, scalar.BigInt(new(big.Int)))

	// Convert to uncompressed bytes format (256 bytes: X_c0, X_c1, Y_c0, Y_c1)
	// Extract coordinates with consistent 48-byte field elements
	xc0Bytes := g2Point.X.A0.Bytes()
	xc1Bytes := g2Point.X.A1.Bytes()
	yc0Bytes := g2Point.Y.A0.Bytes()
	yc1Bytes := g2Point.Y.A1.Bytes()

	// Always place field elements with proper padding
	// Place each 48-byte field element at the end of its 64-byte slot
	copy(point[16:64], xc0Bytes[:])   // X_c0: bytes 16-63
	copy(point[80:128], xc1Bytes[:])  // X_c1: bytes 80-127
	copy(point[144:192], yc0Bytes[:]) // Y_c0: bytes 144-191
	copy(point[208:256], yc1Bytes[:]) // Y_c1: bytes 208-255

	return g.transformer.TransformPrecompileInput(point, 0x0d)
}

// generateBLS12G1AddCall creates a BLS12-381 G1 addition call using gnark-crypto library
func (g *OpcodeGenerator) generateBLS12G1AddCall() []byte {
	var bytecode []byte

	// Generate two BLS12-381 G1 points (128 bytes each = 256 bytes total)
	point1 := g.generateValidBLS12G1PointUncompressed()
	point2 := g.generateValidBLS12G1PointUncompressed()

	// Store points in memory using multiple MSTORE operations
	// Each MSTORE handles 32 bytes, so we need 8 operations for 256 bytes
	for i := 0; i < 8; i++ {
		bytecode = append(bytecode, 0x7f) // PUSH32
		if i < 4 {
			// First point (128 bytes)
			bytecode = append(bytecode, point1[i*32:(i+1)*32]...)
		} else {
			// Second point (128 bytes)
			bytecode = append(bytecode, point2[(i-4)*32:(i-4+1)*32]...)
		}
		bytecode = append(bytecode, 0x60, byte(i*32)) // PUSH1 offset
		bytecode = append(bytecode, 0x52)             // MSTORE
	}

	// Setup CALL to BLS12_G1ADD precompile
	bytecode = append(bytecode, 0x61, 0x01, 0x00) // PUSH2 256 (return size - uncompressed G1 point)
	bytecode = append(bytecode, 0x61, 0x01, 0x00) // PUSH2 256 (return offset)
	bytecode = append(bytecode, 0x61, 0x01, 0x00) // PUSH2 256 (args size - two G1 points)
	bytecode = append(bytecode, 0x60, 0x00)       // PUSH1 0 (args offset)
	bytecode = append(bytecode, 0x60, 0x00)       // PUSH1 0 (value)
	bytecode = append(bytecode, 0x60, 0x0b)       // PUSH1 11 (precompile address)
	bytecode = append(bytecode, 0x5a)             // GAS
	bytecode = append(bytecode, 0xf1)             // CALL

	returnOffset := 0x100
	returnSize := 0x80

	// Process precompile result: LOG in precompiles-only mode, load to stack in normal mode
	bytecode = append(bytecode, g.processPrecompileResult(returnOffset, returnSize)...)

	return bytecode
}

// generateBLS12G1MSMCall creates a BLS12-381 G1 multi-scalar multiplication call
func (g *OpcodeGenerator) generateBLS12G1MSMCall() []byte {
	var bytecode []byte

	// G1MSM input format per EIP-2537: [(point1, scalar1), (point2, scalar2), ...]
	// Each entry: 128 bytes G1 point + 32 bytes scalar = 160 bytes per entry
	// Point format: [0-63] X coordinate, [64-127] Y coordinate, [128-159] scalar
	numPairs := g.rng.Intn(3) + 1 // 1-3 pairs
	totalSize := numPairs * 160

	// Generate point-scalar pairs
	for i := 0; i < numPairs; i++ {
		point := g.generateValidBLS12G1PointUncompressed()
		scalar := g.rng.Bytes(32)

		// Store point coordinates directly with proper alignment
		// Each coordinate is 64 bytes, stored as 2 x 32-byte chunks

		// Store X coordinate (bytes 0-63 of point)
		for j := 0; j < 2; j++ {
			bytecode = append(bytecode, 0x7f) // PUSH32
			bytecode = append(bytecode, point[j*32:(j+1)*32]...)
			offset := i*160 + j*32
			if offset < 256 {
				bytecode = append(bytecode, 0x60, byte(offset)) // PUSH1
			} else {
				bytecode = append(bytecode, 0x61, byte(offset>>8), byte(offset)) // PUSH2
			}
			bytecode = append(bytecode, 0x52) // MSTORE
		}

		// Store Y coordinate (bytes 64-127 of point)
		for j := 0; j < 2; j++ {
			bytecode = append(bytecode, 0x7f) // PUSH32
			bytecode = append(bytecode, point[64+j*32:64+(j+1)*32]...)
			offset := i*160 + 64 + j*32
			if offset < 256 {
				bytecode = append(bytecode, 0x60, byte(offset)) // PUSH1
			} else {
				bytecode = append(bytecode, 0x61, byte(offset>>8), byte(offset)) // PUSH2
			}
			bytecode = append(bytecode, 0x52) // MSTORE
		}

		// Store scalar after point (32 bytes)
		bytecode = append(bytecode, 0x7f) // PUSH32
		bytecode = append(bytecode, scalar...)
		offset := i*160 + 128
		if offset < 256 {
			bytecode = append(bytecode, 0x60, byte(offset)) // PUSH1
		} else {
			bytecode = append(bytecode, 0x61, byte(offset>>8), byte(offset)) // PUSH2
		}
		bytecode = append(bytecode, 0x52) // MSTORE
	}

	// Setup CALL to BLS12_G1MSM precompile
	bytecode = append(bytecode, 0x60, 0x80)                                // PUSH1 128 (return size - G1 point)
	bytecode = append(bytecode, 0x61, 0x01, 0x00)                          // PUSH2 256 (return offset)
	bytecode = append(bytecode, 0x61, byte(totalSize>>8), byte(totalSize)) // PUSH2 input size
	bytecode = append(bytecode, 0x60, 0x00)                                // PUSH1 0 (args offset)
	bytecode = append(bytecode, 0x60, 0x00)                                // PUSH1 0 (value)
	bytecode = append(bytecode, 0x60, 0x0c)                                // PUSH1 12 (precompile address)
	bytecode = append(bytecode, 0x5a)                                      // GAS
	bytecode = append(bytecode, 0xf1)                                      // CALL

	returnOffset := 0x100
	returnSize := 0x80
	// Process precompile result: LOG in precompiles-only mode, load to stack in normal mode
	bytecode = append(bytecode, g.processPrecompileResult(returnOffset, returnSize)...)
	return bytecode
}

// generateBLS12G2AddCall creates a BLS12-381 G2 addition call
func (g *OpcodeGenerator) generateBLS12G2AddCall() []byte {
	var bytecode []byte

	// Generate two BLS12-381 G2 points (256 bytes each = 512 bytes total)
	point1 := g.generateValidBLS12G2PointUncompressed()
	point2 := g.generateValidBLS12G2PointUncompressed()

	// Store points in memory using multiple MSTORE operations
	// Each MSTORE handles 32 bytes, so we need 16 operations for 512 bytes
	for i := 0; i < 16; i++ {
		bytecode = append(bytecode, 0x7f) // PUSH32
		if i < 8 {
			// First point (256 bytes)
			bytecode = append(bytecode, point1[i*32:(i+1)*32]...)
		} else {
			// Second point (256 bytes)
			bytecode = append(bytecode, point2[(i-8)*32:(i-8+1)*32]...)
		}
		bytecode = append(bytecode, 0x60, byte(i*32)) // PUSH1 offset
		bytecode = append(bytecode, 0x52)             // MSTORE
	}

	// Setup CALL to BLS12_G2ADD precompile
	bytecode = append(bytecode, 0x61, 0x01, 0x00) // PUSH2 256 (return size - uncompressed G2 point)
	bytecode = append(bytecode, 0x61, 0x02, 0x00) // PUSH2 512 (return offset)
	bytecode = append(bytecode, 0x61, 0x02, 0x00) // PUSH2 512 (args size - two G2 points)
	bytecode = append(bytecode, 0x60, 0x00)       // PUSH1 0 (args offset)
	bytecode = append(bytecode, 0x60, 0x00)       // PUSH1 0 (value)
	bytecode = append(bytecode, 0x60, 0x0d)       // PUSH1 13 (precompile address)
	bytecode = append(bytecode, 0x5a)             // GAS
	bytecode = append(bytecode, 0xf1)             // CALL

	returnOffset := 0x200
	returnSize := 0x100
	// Process precompile result: LOG in precompiles-only mode, load to stack in normal mode
	bytecode = append(bytecode, g.processPrecompileResult(returnOffset, returnSize)...)
	return bytecode
}

// generateBLS12G2MSMCall creates a BLS12-381 G2 multi-scalar multiplication call
func (g *OpcodeGenerator) generateBLS12G2MSMCall() []byte {
	var bytecode []byte

	// G2MSM input format per EIP-2537: [(point1, scalar1), (point2, scalar2), ...]
	// Each entry: 256 bytes G2 point + 32 bytes scalar = 288 bytes per entry
	// Point format: [0-255] G2 point, [256-287] scalar
	numPairs := g.rng.Intn(3) + 1 // 1-3 pairs
	totalSize := numPairs * 288

	// Generate point-scalar pairs
	for i := 0; i < numPairs; i++ {
		point := g.generateValidBLS12G2PointUncompressed()
		scalar := g.rng.Bytes(32)

		// Store point first (256 bytes = 8 x 32-byte chunks)
		for j := 0; j < 8; j++ {
			bytecode = append(bytecode, 0x7f) // PUSH32
			bytecode = append(bytecode, point[j*32:(j+1)*32]...)
			offset := i*288 + j*32
			if offset < 256 {
				bytecode = append(bytecode, 0x60, byte(offset)) // PUSH1
			} else {
				bytecode = append(bytecode, 0x61, byte(offset>>8), byte(offset)) // PUSH2
			}
			bytecode = append(bytecode, 0x52) // MSTORE
		}

		// Store scalar after point (32 bytes)
		bytecode = append(bytecode, 0x7f) // PUSH32
		bytecode = append(bytecode, scalar...)
		offset := i*288 + 256 // 256 bytes after start of this pair
		if offset < 256 {
			bytecode = append(bytecode, 0x60, byte(offset)) // PUSH1
		} else {
			bytecode = append(bytecode, 0x61, byte(offset>>8), byte(offset)) // PUSH2
		}
		bytecode = append(bytecode, 0x52) // MSTORE
	}

	// Setup CALL to BLS12_G2MSM precompile
	bytecode = append(bytecode, 0x61, 0x01, 0x00)                          // PUSH2 256 (return size - G2 point)
	bytecode = append(bytecode, 0x61, 0x03, 0x00)                          // PUSH2 768 (return offset)
	bytecode = append(bytecode, 0x61, byte(totalSize>>8), byte(totalSize)) // PUSH2 input size
	bytecode = append(bytecode, 0x60, 0x00)                                // PUSH1 0 (args offset)
	bytecode = append(bytecode, 0x60, 0x00)                                // PUSH1 0 (value)
	bytecode = append(bytecode, 0x60, 0x0e)                                // PUSH1 14 (precompile address)
	bytecode = append(bytecode, 0x5a)                                      // GAS
	bytecode = append(bytecode, 0xf1)                                      // CALL

	returnOffset := 0x300
	returnSize := 0x100
	// Process precompile result: LOG in precompiles-only mode, load to stack in normal mode
	bytecode = append(bytecode, g.processPrecompileResult(returnOffset, returnSize)...)
	return bytecode
}

// generateBLS12PairingCall creates a BLS12-381 pairing call
func (g *OpcodeGenerator) generateBLS12PairingCall() []byte {
	var bytecode []byte

	// Pairing input format: [(G1_point1, G2_point1), (G1_point2, G2_point2), ...]
	// Each pair: 128 bytes G1 + 256 bytes G2 = 384 bytes per pair
	numPairs := g.rng.Intn(3) + 1 // 1-3 pairs
	totalSize := numPairs * 384

	// Generate G1-G2 pairs
	for i := 0; i < numPairs; i++ {
		g1Point := g.generateValidBLS12G1PointUncompressed()
		g2Point := g.generateValidBLS12G2PointUncompressed()

		// Store G1 point (128 bytes = 4 x 32-byte chunks)
		for j := 0; j < 4; j++ {
			bytecode = append(bytecode, 0x7f) // PUSH32
			bytecode = append(bytecode, g1Point[j*32:(j+1)*32]...)
			offset := i*384 + j*32
			bytecode = append(bytecode, 0x61, byte(offset>>8), byte(offset)) // PUSH2 offset
			bytecode = append(bytecode, 0x52)                                // MSTORE
		}

		// Store G2 point (256 bytes = 8 x 32-byte chunks)
		for j := 0; j < 8; j++ {
			bytecode = append(bytecode, 0x7f) // PUSH32
			bytecode = append(bytecode, g2Point[j*32:(j+1)*32]...)
			offset := i*384 + 128 + j*32
			bytecode = append(bytecode, 0x61, byte(offset>>8), byte(offset)) // PUSH2 offset
			bytecode = append(bytecode, 0x52)                                // MSTORE
		}
	}

	// Setup CALL to BLS12_PAIRING precompile
	bytecode = append(bytecode, 0x60, 0x20)                                // PUSH1 32 (return size - boolean result)
	bytecode = append(bytecode, 0x61, 0x04, 0x00)                          // PUSH2 1024 (return offset)
	bytecode = append(bytecode, 0x61, byte(totalSize>>8), byte(totalSize)) // PUSH2 input size
	bytecode = append(bytecode, 0x60, 0x00)                                // PUSH1 0 (args offset)
	bytecode = append(bytecode, 0x60, 0x00)                                // PUSH1 0 (value)
	bytecode = append(bytecode, 0x60, 0x0f)                                // PUSH1 15 (precompile address)
	bytecode = append(bytecode, 0x5a)                                      // GAS
	bytecode = append(bytecode, 0xf1)                                      // CALL

	returnOffset := 0x400
	returnSize := 0x20
	// Process precompile result: LOG in precompiles-only mode, load to stack in normal mode
	bytecode = append(bytecode, g.processPrecompileResult(returnOffset, returnSize)...)
	return bytecode
}

// generateBLS12MapFpToG1Call creates a BLS12-381 map field element to G1 call
func (g *OpcodeGenerator) generateBLS12MapFpToG1Call() []byte {
	var bytecode []byte

	// Generate a field element (64 bytes)
	fieldElement := make([]byte, 64)

	choice := g.rng.Intn(100)
	if choice < 10 {
		// 10% chance of invalid field elements for testing
		copy(fieldElement, g.rng.Bytes(64))
		fieldElement[0] |= 0xF0 // Ensure > modulus
	} else {
		// 90% chance of valid field elements
		var fpElem fp.Element
		fpBytes := g.rng.Bytes(48) // 48 bytes for BLS12-381 field element
		fpElem.SetBytes(fpBytes)
		fpResult := fpElem.Bytes()
		// Pad to 64 bytes: place 48-byte field element at end
		copy(fieldElement[64-len(fpResult):64], fpResult[:])
	}

	// Store field element in memory (64 bytes = 2 x 32-byte chunks)
	for i := 0; i < 2; i++ {
		bytecode = append(bytecode, 0x7f) // PUSH32
		bytecode = append(bytecode, fieldElement[i*32:(i+1)*32]...)
		bytecode = append(bytecode, 0x60, byte(i*32)) // PUSH1 offset
		bytecode = append(bytecode, 0x52)             // MSTORE
	}

	// Setup CALL to BLS12_MAP_FP_TO_G1 precompile
	bytecode = append(bytecode, 0x60, 0x80) // PUSH1 128 (return size - G1 point)
	bytecode = append(bytecode, 0x60, 0x40) // PUSH1 64 (return offset)
	bytecode = append(bytecode, 0x60, 0x40) // PUSH1 64 (args size - field element)
	bytecode = append(bytecode, 0x60, 0x00) // PUSH1 0 (args offset)
	bytecode = append(bytecode, 0x60, 0x00) // PUSH1 0 (value)
	bytecode = append(bytecode, 0x60, 0x10) // PUSH1 16 (precompile address)
	bytecode = append(bytecode, 0x5a)       // GAS
	bytecode = append(bytecode, 0xf1)       // CALL

	returnOffset := 0x40
	returnSize := 0x80
	// Process precompile result: LOG in precompiles-only mode, load to stack in normal mode
	bytecode = append(bytecode, g.processPrecompileResult(returnOffset, returnSize)...)
	return bytecode
}

// generateBLS12MapFp2G2Call creates a BLS12-381 map field extension element to G2 call
func (g *OpcodeGenerator) generateBLS12MapFp2G2Call() []byte {
	var bytecode []byte

	// Generate a field extension element (128 bytes: 64 bytes c0 + 64 bytes c1)
	fieldElement := make([]byte, 128)

	choice := g.rng.Intn(100)
	if choice < 10 {
		// 10% chance of invalid field elements for testing
		copy(fieldElement, g.rng.Bytes(128))
		fieldElement[0] |= 0xF0  // Ensure c0 > modulus
		fieldElement[64] |= 0xF0 // Ensure c1 > modulus
	} else {
		// 90% chance of valid field elements
		var fp2 bls12381.E2
		// Generate c0 and c1 components
		c0Bytes := g.rng.Bytes(48)
		c1Bytes := g.rng.Bytes(48)
		fp2.A0.SetBytes(c0Bytes)
		fp2.A1.SetBytes(c1Bytes)

		c0Result := fp2.A0.Bytes()
		c1Result := fp2.A1.Bytes()
		// Pad to 64 bytes each: place 48-byte field elements at end
		copy(fieldElement[64-len(c0Result):64], c0Result[:])   // First 64 bytes: c0 component
		copy(fieldElement[128-len(c1Result):128], c1Result[:]) // Second 64 bytes: c1 component
	}

	// Store field extension element in memory (128 bytes = 4 x 32-byte chunks)
	for i := 0; i < 4; i++ {
		bytecode = append(bytecode, 0x7f) // PUSH32
		bytecode = append(bytecode, fieldElement[i*32:(i+1)*32]...)
		bytecode = append(bytecode, 0x60, byte(i*32)) // PUSH1 offset
		bytecode = append(bytecode, 0x52)             // MSTORE
	}

	// Setup CALL to BLS12_MAP_FP2_TO_G2 precompile
	bytecode = append(bytecode, 0x61, 0x01, 0x00) // PUSH2 256 (return size - G2 point)
	bytecode = append(bytecode, 0x60, 0x80)       // PUSH1 128 (return offset)
	bytecode = append(bytecode, 0x60, 0x80)       // PUSH1 128 (args size - field extension element)
	bytecode = append(bytecode, 0x60, 0x00)       // PUSH1 0 (args offset)
	bytecode = append(bytecode, 0x60, 0x00)       // PUSH1 0 (value)
	bytecode = append(bytecode, 0x60, 0x11)       // PUSH1 17 (precompile address)
	bytecode = append(bytecode, 0x5a)             // GAS
	bytecode = append(bytecode, 0xf1)             // CALL

	returnOffset := 0x80
	returnSize := 0x100
	// Process precompile result: LOG in precompiles-only mode, load to stack in normal mode
	bytecode = append(bytecode, g.processPrecompileResult(returnOffset, returnSize)...)
	return bytecode
}
