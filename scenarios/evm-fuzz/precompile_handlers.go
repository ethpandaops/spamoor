package evmfuzz

import (
	"encoding/binary"
	"math/big"
)

// BN256 curve parameters
var (
	// BN256 field modulus: 21888242871839275222246405745257275088696311157297823662689037894645226208583
	bn256FieldModulus = new(big.Int)
	bn256One          = big.NewInt(1)
)

func init() {
	// Initialize BN256 field modulus
	bn256FieldModulus.SetString("21888242871839275222246405745257275088696311157297823662689037894645226208583", 10)
}

// Precompile call templates
func ecrecoverTemplate(g *OpcodeGenerator) []byte { return g.generatePrecompileCall(1) }
func sha256Template(g *OpcodeGenerator) []byte    { return g.generatePrecompileCall(2) }
func ripemd160Template(g *OpcodeGenerator) []byte { return g.generatePrecompileCall(3) }
func identityTemplate(g *OpcodeGenerator) []byte  { return g.generatePrecompileCall(4) }
func modexpTemplate(g *OpcodeGenerator) []byte    { return g.generatePrecompileCall(5) }
func ecAddTemplate(g *OpcodeGenerator) []byte     { return g.generatePrecompileCall(6) }
func ecMulTemplate(g *OpcodeGenerator) []byte     { return g.generatePrecompileCall(7) }
func ecPairingTemplate(g *OpcodeGenerator) []byte { return g.generatePrecompileCall(8) }
func blake2fTemplate(g *OpcodeGenerator) []byte   { return g.generatePrecompileCall(9) }
func pointEvalTemplate(g *OpcodeGenerator) []byte { return g.generatePrecompileCall(10) }

// BLS12-381 precompile templates
func bls12G1AddTemplate(g *OpcodeGenerator) []byte     { return g.generatePrecompileCall(11) } // 0x0b
func bls12G1MSMTemplate(g *OpcodeGenerator) []byte     { return g.generatePrecompileCall(12) } // 0x0c
func bls12G2AddTemplate(g *OpcodeGenerator) []byte     { return g.generatePrecompileCall(13) } // 0x0d
func bls12G2MSMTemplate(g *OpcodeGenerator) []byte     { return g.generatePrecompileCall(14) } // 0x0e
func bls12PairingTemplate(g *OpcodeGenerator) []byte   { return g.generatePrecompileCall(15) } // 0x0f
func bls12MapFpToG1Template(g *OpcodeGenerator) []byte { return g.generatePrecompileCall(16) } // 0x10
func bls12MapFp2G2Template(g *OpcodeGenerator) []byte  { return g.generatePrecompileCall(17) } // 0x11

// generateValidBN256Point generates a BN256 curve point - mostly valid, sometimes invalid for testing
func (g *OpcodeGenerator) generateValidBN256Point() []byte {
	point := make([]byte, 64) // 32 bytes x + 32 bytes y

	// 4% chance to generate invalid BN256 points for testing edge cases
	if g.rng.Float64() < 0.04 {
		invalidType := g.rng.Intn(4)
		switch invalidType {
		case 0:
			// Coordinates >= field modulus
			copy(point[:32], g.rng.Bytes(32))
			copy(point[32:], g.rng.Bytes(32))
			// Ensure coordinates are >= field modulus
			point[0] |= 0x80  // Make x large
			point[32] |= 0x80 // Make y large
		case 1:
			// All 0xFF bytes (definitely > modulus)
			for i := range point {
				point[i] = 0xFF
			}
		case 2:
			// Valid x, invalid y (point not on curve)
			x := big.NewInt(int64(g.rng.Intn(1000) + 1))
			xBytes := make([]byte, 32)
			x.FillBytes(xBytes)
			copy(point[:32], xBytes)
			// Generate deliberately wrong y coordinate
			copy(point[32:], g.rng.Bytes(32))
			point[32] |= 0x80 // Ensure invalid
		default:
			// Coordinates that would cause overflow in calculations
			// Use values near the field modulus
			modBytes := bn256FieldModulus.Bytes()
			copy(point[:32], modBytes)
			copy(point[32:], modBytes)
			// Add small random increment to go over modulus
			point[31] += byte(g.rng.Intn(10) + 1)
			point[63] += byte(g.rng.Intn(10) + 1)
		}
		return point
	}

	// 10% chance of zero point (point at infinity)
	if g.rng.Float64() < 0.1 {
		return point // All zeros = point at infinity
	}

	// Generate mostly valid curve point by trying random x values
	// For efficiency, we'll generate a few known good patterns
	patterns := []func() []byte{
		// Pattern 1: Use small x values (often valid)
		func() []byte {
			x := big.NewInt(int64(g.rng.Intn(1000) + 1))
			y := g.calculateBN256Y(x)
			if y == nil {
				// Fallback to generator point
				return g.getBN256Generator()
			}
			xBytes := make([]byte, 32)
			yBytes := make([]byte, 32)
			x.FillBytes(xBytes)
			y.FillBytes(yBytes)
			return append(xBytes, yBytes...)
		},
		// Pattern 2: Use generator point multiples
		func() []byte {
			return g.getBN256Generator()
		},
		// Pattern 3: Use modular reduction of random values
		func() []byte {
			randomX := new(big.Int).SetBytes(g.rng.Bytes(32))
			x := new(big.Int).Mod(randomX, bn256FieldModulus)
			y := g.calculateBN256Y(x)
			if y == nil {
				return g.getBN256Generator()
			}
			xBytes := make([]byte, 32)
			yBytes := make([]byte, 32)
			x.FillBytes(xBytes)
			y.FillBytes(yBytes)
			return append(xBytes, yBytes...)
		},
	}

	return patterns[g.rng.Intn(len(patterns))]()
}

// calculateBN256Y calculates y coordinate for given x on BN256 curve: y² = x³ + 3
func (g *OpcodeGenerator) calculateBN256Y(x *big.Int) *big.Int {
	// y² = x³ + 3 (mod p)
	x3 := new(big.Int).Exp(x, big.NewInt(3), bn256FieldModulus)
	x3Plus3 := new(big.Int).Add(x3, big.NewInt(3))
	x3Plus3.Mod(x3Plus3, bn256FieldModulus)

	// Check if x³ + 3 is a quadratic residue
	y := g.modularSqrt(x3Plus3, bn256FieldModulus)
	return y
}

// modularSqrt calculates square root modulo prime using Tonelli-Shanks algorithm (simplified)
func (g *OpcodeGenerator) modularSqrt(a, p *big.Int) *big.Int {
	// For BN256 field, p ≡ 3 (mod 4), so we can use simple formula: y = a^((p+1)/4) mod p
	exp := new(big.Int).Add(p, bn256One)
	exp.Div(exp, big.NewInt(4))
	y := new(big.Int).Exp(a, exp, p)

	// Verify: y² ≡ a (mod p)
	ySquared := new(big.Int).Exp(y, big.NewInt(2), p)
	if ySquared.Cmp(a) != 0 {
		return nil // Not a quadratic residue
	}
	return y
}

// getBN256Generator returns the BN256 generator point
func (g *OpcodeGenerator) getBN256Generator() []byte {
	// BN256 generator point (1, 2)
	point := make([]byte, 64)
	point[31] = 1 // x = 1
	point[63] = 2 // y = 2
	return point
}

// generateValidBN256Scalar generates a valid scalar for BN256 operations
func (g *OpcodeGenerator) generateValidBN256Scalar() []byte {
	// Generate random scalar modulo curve order
	// For simplicity, just use a 32-byte value (will be reduced by precompile)
	scalar := g.rng.Bytes(32)

	// 5% chance of small scalars (often interesting edge cases)
	if g.rng.Float64() < 0.05 {
		small := big.NewInt(int64(g.rng.Intn(100)))
		smallBytes := make([]byte, 32)
		small.FillBytes(smallBytes)
		return smallBytes
	}

	return scalar
}

// generatePrecompileCall creates a precompile call with proper calldata setup
func (g *OpcodeGenerator) generatePrecompileCall(precompileAddr uint8) []byte {
	// Calculate calldata size and number of 32-byte words needed
	var calldataSize uint16
	var wordCount int

	switch precompileAddr {
	case 1: // ecrecover - hash, v, r, s
		calldataSize = 128 // 4 * 32 bytes
		wordCount = 4
	case 2: // sha256 - variable length data
		calldataSize = 64 // Use 64 bytes
		wordCount = 2
	case 3: // ripemd160 - variable length data
		calldataSize = 64 // Use 64 bytes
		wordCount = 2
	case 4: // identity - variable length data
		calldataSize = 32 // Use 32 bytes
		wordCount = 1
	case 5: // modexp - special handling with dynamic lengths
		calldataSize = 384 // 96 bytes for lengths + up to 288 bytes for data
		wordCount = 0      // Special case - don't use stack values
	case 6: // ecAdd - x1, y1, x2, y2 (BN256)
		calldataSize = 128 // 4 * 32 bytes
		wordCount = 0      // Special BN256 handling
	case 7: // ecMul - x, y, scalar (BN256)
		calldataSize = 96 // 3 * 32 bytes
		wordCount = 0     // Special BN256 handling
	case 8: // ecPairing - pairs of BN256 points
		calldataSize = 192 // 6 * 32 bytes minimum (1 pair)
		wordCount = 0      // Special BN256 handling
	case 9: // blake2f - rounds, h, m, t, f
		calldataSize = 213 // Fixed size
		wordCount = 7      // Approximate for memory setup
	case 10: // pointEval - versioned_hash, x, y, commitment, proof
		calldataSize = 192 // 6 * 32 bytes
		wordCount = 0      // Special KZG handling
	case 11: // bls12_g1add - two G1 points (x1, y1, x2, y2)
		calldataSize = 128 // 4 * 32 bytes
		wordCount = 0      // Special BLS12 handling
	case 12: // bls12_g1mul - G1 point and scalar (x, y, scalar)
		calldataSize = 96 // 3 * 32 bytes
		wordCount = 0     // Special BLS12 handling
	case 13: // bls12_g1msm - multiple G1 points and scalars
		calldataSize = 384 // Variable, use 12 * 32 bytes (4 pairs)
		wordCount = 0      // Special handling
	case 14: // bls12_g2add - two G2 points
		calldataSize = 256 // 8 * 32 bytes
		wordCount = 0      // Special BLS12 handling
	case 15: // bls12_g2mul - G2 point and scalar
		calldataSize = 160 // 5 * 32 bytes
		wordCount = 0      // Special BLS12 handling
	case 16: // bls12_g2msm - multiple G2 points and scalars
		calldataSize = 640 // Variable, use 20 * 32 bytes (4 pairs)
		wordCount = 0      // Special handling
	case 17: // bls12_pairing - pairs of G1 and G2 points
		calldataSize = 384 // 12 * 32 bytes (2 pairs minimum)
		wordCount = 0      // Special BLS12 handling
	case 18: // bls12_map_fp2_to_g2 - field element
		calldataSize = 128 // 4 * 32 bytes
		wordCount = 0      // Special BLS12 handling
	default:
		calldataSize = 32
		wordCount = 1
	}

	var bytecode []byte

	// Special handling for different precompile types
	switch precompileAddr {
	case 5: // MODEXP
		return g.generateModexpCall()
	case 6: // BN256 ecAdd
		return g.generateBN256EcAddCall()
	case 7: // BN256 ecMul
		return g.generateBN256EcMulCall()
	case 8: // BN256 ecPairing
		return g.generateBN256PairingCall()
	case 10: // KZG Point Evaluation
		return g.generateKZGPointEvalCall()
	case 12: // BLS12 G1 MSM
		return g.generateBLS12G1MSMCall()
	case 14: // BLS12 G2 MSM
		return g.generateBLS12G2MSMCall()
	case 11, 13, 15, 16, 17: // Other BLS12 operations
		return g.generateBLS12Call(precompileAddr)
	}

	// Standard precompile handling for simple cases
	// Write stack inputs to memory at offset 0
	for i := 0; i < wordCount && i < len(g.stack); i++ {
		memOffset := i * 32

		// If we have known stack values, use them; otherwise generate random data
		if i < len(g.stack) && g.stack[len(g.stack)-1-i].Known {
			// Use the known value from stack (reverse order since stack is LIFO)
			value := g.stack[len(g.stack)-1-i].Value

			// PUSH32 the value
			bytecode = append(bytecode, 0x7f) // PUSH32
			bytecode = append(bytecode, value...)
		} else {
			// Generate random 32-byte value
			randomData := g.rng.Bytes(32)
			bytecode = append(bytecode, 0x7f) // PUSH32
			bytecode = append(bytecode, randomData...)
		}

		// PUSH memory offset
		if memOffset < 256 {
			bytecode = append(bytecode, 0x60, byte(memOffset)) // PUSH1 offset
		} else {
			bytecode = append(bytecode, 0x61, byte(memOffset>>8), byte(memOffset)) // PUSH2 offset
		}

		// MSTORE to write the 32-byte value to memory
		bytecode = append(bytecode, 0x52) // MSTORE
	}

	// If we need more data than available stack items, fill with random data
	for i := len(g.stack); i < wordCount; i++ {
		memOffset := i * 32

		// Generate random 32-byte value
		randomData := g.rng.Bytes(32)
		bytecode = append(bytecode, 0x7f) // PUSH32
		bytecode = append(bytecode, randomData...)

		// PUSH memory offset
		if memOffset < 256 {
			bytecode = append(bytecode, 0x60, byte(memOffset)) // PUSH1 offset
		} else {
			bytecode = append(bytecode, 0x61, byte(memOffset>>8), byte(memOffset)) // PUSH2 offset
		}

		// MSTORE
		bytecode = append(bytecode, 0x52) // MSTORE
	}

	// Now setup the CALL to the precompile
	// CALL(gas, address, value, argsOffset, argsSize, retOffset, retSize)
	// Push in reverse order so gas is at top of stack (position 0)

	// Return data size (32 bytes for most precompiles)
	retSize := 32
	if precompileAddr == 8 { // ecPairing returns 32 bytes (0 or 1)
		retSize = 32
	} else if precompileAddr == 9 { // blake2f returns 64 bytes
		retSize = 64
	}
	bytecode = append(bytecode, 0x60, byte(retSize)) // PUSH1 retSize

	// Return data offset (store return after input data)
	returnOffset := (wordCount * 32) + 32 // After our input data
	if returnOffset < 256 {
		bytecode = append(bytecode, 0x60, byte(returnOffset)) // PUSH1 retOffset
	} else {
		bytecode = append(bytecode, 0x61, byte(returnOffset>>8), byte(returnOffset)) // PUSH2 retOffset
	}

	// Arguments size
	if calldataSize < 256 {
		bytecode = append(bytecode, 0x60, byte(calldataSize)) // PUSH1 calldataSize
	} else {
		bytecode = append(bytecode, 0x61, byte(calldataSize>>8), byte(calldataSize)) // PUSH2 calldataSize
	}

	// Arguments offset (0 - we stored data starting at memory 0)
	bytecode = append(bytecode, 0x60, 0x00) // PUSH1 0

	// Value (always 0 for precompiles)
	bytecode = append(bytecode, 0x60, 0x00) // PUSH1 0

	// Precompile address
	bytecode = append(bytecode, 0x60, precompileAddr) // PUSH1 precompileAddr

	// Gas for the call - 90% use GAS opcode, 10% use fixed value
	if g.rng.Float64() < 0.9 {
		bytecode = append(bytecode, 0x5a) // GAS opcode - use remaining gas
	} else {
		bytecode = append(bytecode, 0x61, 0x13, 0x88) // PUSH2 5000
	}

	// CALL opcode
	bytecode = append(bytecode, 0xf1) // CALL

	// Add LOG0 to log the result when in precompiles-only mode
	bytecode = append(bytecode, g.addLogIfPrecompilesOnly(returnOffset, retSize)...)

	return bytecode
}

// generateBN256EcAddCall creates a BN256 elliptic curve point addition call
func (g *OpcodeGenerator) generateBN256EcAddCall() []byte {
	var bytecode []byte

	// Generate two valid BN256 points
	point1 := g.generateValidBN256Point()
	point2 := g.generateValidBN256Point()

	// Write first point (x1, y1) to memory at offset 0
	bytecode = append(bytecode, 0x7f)           // PUSH32
	bytecode = append(bytecode, point1[:32]...) // x1
	bytecode = append(bytecode, 0x60, 0x00)     // PUSH1 0
	bytecode = append(bytecode, 0x52)           // MSTORE

	bytecode = append(bytecode, 0x7f)           // PUSH32
	bytecode = append(bytecode, point1[32:]...) // y1
	bytecode = append(bytecode, 0x60, 0x20)     // PUSH1 32
	bytecode = append(bytecode, 0x52)           // MSTORE

	// Write second point (x2, y2) to memory at offset 64
	bytecode = append(bytecode, 0x7f)           // PUSH32
	bytecode = append(bytecode, point2[:32]...) // x2
	bytecode = append(bytecode, 0x60, 0x40)     // PUSH1 64
	bytecode = append(bytecode, 0x52)           // MSTORE

	bytecode = append(bytecode, 0x7f)           // PUSH32
	bytecode = append(bytecode, point2[32:]...) // y2
	bytecode = append(bytecode, 0x60, 0x60)     // PUSH1 96
	bytecode = append(bytecode, 0x52)           // MSTORE

	// Setup CALL to ecAdd precompile
	bytecode = append(bytecode, 0x60, 0x40) // PUSH1 64 (return size)
	bytecode = append(bytecode, 0x60, 0x80) // PUSH1 128 (return offset)
	bytecode = append(bytecode, 0x60, 0x80) // PUSH1 128 (args size)
	bytecode = append(bytecode, 0x60, 0x00) // PUSH1 0 (args offset)
	bytecode = append(bytecode, 0x60, 0x00) // PUSH1 0 (value)
	bytecode = append(bytecode, 0x60, 0x06) // PUSH1 6 (precompile address)
	bytecode = append(bytecode, 0x5a)       // GAS
	bytecode = append(bytecode, 0xf1)       // CALL

	// Add LOG0 to log the result when in precompiles-only mode (return at offset 128, size 64)
	bytecode = append(bytecode, g.addLogIfPrecompilesOnly(128, 64)...)

	return bytecode
}

// generateBN256EcMulCall creates a BN256 elliptic curve scalar multiplication call
func (g *OpcodeGenerator) generateBN256EcMulCall() []byte {
	var bytecode []byte

	// Generate valid BN256 point and scalar
	point := g.generateValidBN256Point()
	scalar := g.generateValidBN256Scalar()

	// Write point (x, y) to memory at offset 0
	bytecode = append(bytecode, 0x7f)          // PUSH32
	bytecode = append(bytecode, point[:32]...) // x
	bytecode = append(bytecode, 0x60, 0x00)    // PUSH1 0
	bytecode = append(bytecode, 0x52)          // MSTORE

	bytecode = append(bytecode, 0x7f)          // PUSH32
	bytecode = append(bytecode, point[32:]...) // y
	bytecode = append(bytecode, 0x60, 0x20)    // PUSH1 32
	bytecode = append(bytecode, 0x52)          // MSTORE

	// Write scalar to memory at offset 64
	bytecode = append(bytecode, 0x7f)       // PUSH32
	bytecode = append(bytecode, scalar...)  // scalar
	bytecode = append(bytecode, 0x60, 0x40) // PUSH1 64
	bytecode = append(bytecode, 0x52)       // MSTORE

	// Setup CALL to ecMul precompile
	bytecode = append(bytecode, 0x60, 0x40) // PUSH1 64 (return size)
	bytecode = append(bytecode, 0x60, 0x60) // PUSH1 96 (return offset)
	bytecode = append(bytecode, 0x60, 0x60) // PUSH1 96 (args size)
	bytecode = append(bytecode, 0x60, 0x00) // PUSH1 0 (args offset)
	bytecode = append(bytecode, 0x60, 0x00) // PUSH1 0 (value)
	bytecode = append(bytecode, 0x60, 0x07) // PUSH1 7 (precompile address)
	bytecode = append(bytecode, 0x5a)       // GAS
	bytecode = append(bytecode, 0xf1)       // CALL

	// Add LOG0 to log the result when in precompiles-only mode (return at offset 96, size 64)
	bytecode = append(bytecode, g.addLogIfPrecompilesOnly(96, 64)...)

	return bytecode
}

// generateBN256PairingCall creates a BN256 pairing check call
func (g *OpcodeGenerator) generateBN256PairingCall() []byte {
	var bytecode []byte

	// Generate 1-3 pairs of G1 and G2 points
	pairCount := 1 + g.rng.Intn(3) // 1, 2, or 3 pairs
	totalLen := pairCount * 192    // Each pair: 64 bytes (G1) + 128 bytes (G2)

	// Write pairs to memory
	for i := 0; i < pairCount; i++ {
		memOffset := i * 192

		// Generate G1 point (64 bytes)
		g1Point := g.generateValidBN256Point()
		bytecode = append(bytecode, 0x7f)            // PUSH32
		bytecode = append(bytecode, g1Point[:32]...) // G1 x
		if memOffset < 256 {
			bytecode = append(bytecode, 0x60, byte(memOffset)) // PUSH1 offset
		} else {
			bytecode = append(bytecode, 0x61, byte(memOffset>>8), byte(memOffset)) // PUSH2 offset
		}
		bytecode = append(bytecode, 0x52) // MSTORE

		bytecode = append(bytecode, 0x7f)            // PUSH32
		bytecode = append(bytecode, g1Point[32:]...) // G1 y
		g1yOffset := memOffset + 32
		if g1yOffset < 256 {
			bytecode = append(bytecode, 0x60, byte(g1yOffset)) // PUSH1 offset
		} else {
			bytecode = append(bytecode, 0x61, byte(g1yOffset>>8), byte(g1yOffset)) // PUSH2 offset
		}
		bytecode = append(bytecode, 0x52) // MSTORE

		// Generate G2 point (128 bytes: x1, x2, y1, y2)
		g2Point := g.generateValidBN256G2Point()
		for j := 0; j < 4; j++ {
			bytecode = append(bytecode, 0x7f) // PUSH32
			bytecode = append(bytecode, g2Point[j*32:(j+1)*32]...)
			g2Offset := memOffset + 64 + j*32
			if g2Offset < 256 {
				bytecode = append(bytecode, 0x60, byte(g2Offset)) // PUSH1 offset
			} else {
				bytecode = append(bytecode, 0x61, byte(g2Offset>>8), byte(g2Offset)) // PUSH2 offset
			}
			bytecode = append(bytecode, 0x52) // MSTORE
		}
	}

	// Setup CALL to ecPairing precompile
	bytecode = append(bytecode, 0x60, 0x20) // PUSH1 32 (return size)
	returnOffset := totalLen + 32
	if returnOffset < 256 {
		bytecode = append(bytecode, 0x60, byte(returnOffset)) // PUSH1 retOffset
	} else {
		bytecode = append(bytecode, 0x61, byte(returnOffset>>8), byte(returnOffset)) // PUSH2 retOffset
	}
	if totalLen < 256 {
		bytecode = append(bytecode, 0x60, byte(totalLen)) // PUSH1 totalLen
	} else {
		bytecode = append(bytecode, 0x61, byte(totalLen>>8), byte(totalLen)) // PUSH2 totalLen
	}
	bytecode = append(bytecode, 0x60, 0x00) // PUSH1 0 (args offset)
	bytecode = append(bytecode, 0x60, 0x00) // PUSH1 0 (value)
	bytecode = append(bytecode, 0x60, 0x08) // PUSH1 8 (precompile address)
	bytecode = append(bytecode, 0x5a)       // GAS
	bytecode = append(bytecode, 0xf1)       // CALL

	// Add LOG0 to log the result when in precompiles-only mode (return at returnOffset, size 32)
	bytecode = append(bytecode, g.addLogIfPrecompilesOnly(returnOffset, 32)...)

	return bytecode
}

// generateValidBN256G2Point generates a valid BN256 G2 point
func (g *OpcodeGenerator) generateValidBN256G2Point() []byte {
	// For simplicity, use the G2 generator point or zero point
	if g.rng.Float64() < 0.2 {
		// Return zero point (point at infinity)
		return make([]byte, 128)
	}

	// Return BN256 G2 generator point coordinates
	// This is a known valid point on the G2 curve
	point := make([]byte, 128)
	// Use a simplified valid G2 point (in practice, these would be specific field elements)
	// For fuzzing, we'll use patterns that are likely to be valid
	binary.BigEndian.PutUint64(point[24:32], 1)   // x1 = 1
	binary.BigEndian.PutUint64(point[56:64], 2)   // x2 = 2
	binary.BigEndian.PutUint64(point[88:96], 1)   // y1 = 1
	binary.BigEndian.PutUint64(point[120:128], 2) // y2 = 2
	return point
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
	bytecode = append(bytecode, g.addLogIfPrecompilesOnly(0xe0, 64)...)

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

// generateValidBLS12BaseField generates a BLS12-381 base field element (64 bytes) - mostly valid, sometimes invalid for testing
func (g *OpcodeGenerator) generateValidBLS12BaseField() []byte {
	// BLS12-381 base field modulus (p)
	// 0x1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaaab
	fieldElement := make([]byte, 64)

	// 5% chance to generate invalid field elements for testing edge cases
	if g.rng.Float64() < 0.05 {
		invalidType := g.rng.Intn(3)
		switch invalidType {
		case 0:
			// Field element >= modulus (set high bits)
			copy(fieldElement, g.rng.Bytes(64))
			fieldElement[0] |= 0xE0 // Set top 3 bits to ensure >= p
		case 1:
			// All 0xFF bytes (definitely > modulus)
			for i := range fieldElement {
				fieldElement[i] = 0xFF
			}
		default:
			// Modulus + small random value
			copy(fieldElement, g.rng.Bytes(64))
			fieldElement[0] = 0x1A // Start with modulus high byte
			fieldElement[1] = 0x01
			fieldElement[2] = 0x11
			fieldElement[3] = 0xEA
		}
		return fieldElement
	}

	// Generate mostly valid patterns
	patternType := g.rng.Intn(4)
	switch patternType {
	case 0:
		// Small values (often useful for testing)
		smallVal := g.rng.Intn(1000)
		binary.BigEndian.PutUint64(fieldElement[56:], uint64(smallVal))
	case 1:
		// Powers of 2 (common edge cases)
		power := g.rng.Intn(63)
		if power < 56 {
			binary.BigEndian.PutUint64(fieldElement[56:], 1<<power)
		} else {
			binary.BigEndian.PutUint64(fieldElement[48:56], 1<<(power-56))
		}
	case 2:
		// Random but constrained to be less than modulus
		// Use simple constraint: clear high bits to ensure < p
		copy(fieldElement, g.rng.Bytes(64))
		fieldElement[0] &= 0x1F // Clear top 3 bits to ensure < p
	default:
		// Medium random values
		copy(fieldElement[32:], g.rng.Bytes(32)) // Only use lower 32 bytes
	}

	return fieldElement
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
	bytecode = append(bytecode, g.addLogIfPrecompilesOnly(returnOffset, modLen)...)

	return bytecode
}

// generateBLS12Call creates calls to BLS12-381 precompiles with proper field element generation
func (g *OpcodeGenerator) generateBLS12Call(precompileAddr uint8) []byte {
	var bytecode []byte

	switch precompileAddr {
	case 11: // BLS12_G1ADD
		return g.generateBLS12G1AddCall()
	case 13: // BLS12_G2ADD
		return g.generateBLS12G2AddCall()
	case 15: // BLS12_PAIRING_CHECK
		return g.generateBLS12PairingCall()
	case 16: // BLS12_MAP_FP_TO_G1
		return g.generateBLS12MapFpToG1Call()
	case 17: // BLS12_MAP_FP2_TO_G2
		return g.generateBLS12MapFp2G2Call()
	default:
		// Fallback to simple random data
		bytecode = append(bytecode, 0x7f) // PUSH32
		bytecode = append(bytecode, g.rng.Bytes(32)...)
		bytecode = append(bytecode, 0x60, 0x00)           // PUSH1 0
		bytecode = append(bytecode, 0x52)                 // MSTORE
		bytecode = append(bytecode, 0x60, 0x20)           // PUSH1 32 (return size)
		bytecode = append(bytecode, 0x60, 0x20)           // PUSH1 32 (return offset)
		bytecode = append(bytecode, 0x60, 0x20)           // PUSH1 32 (args size)
		bytecode = append(bytecode, 0x60, 0x00)           // PUSH1 0 (args offset)
		bytecode = append(bytecode, 0x60, 0x00)           // PUSH1 0 (value)
		bytecode = append(bytecode, 0x60, precompileAddr) // PUSH1 precompileAddr
		bytecode = append(bytecode, 0x5a)                 // GAS
		bytecode = append(bytecode, 0xf1)                 // CALL

		// Add LOG0 to log the result when in precompiles-only mode
		bytecode = append(bytecode, g.addLogIfPrecompilesOnly(32, 32)...)
	}

	return bytecode
}

// generateValidBLS12G1PointUncompressed generates a BLS12-381 G1 point in uncompressed format (128 bytes) - mostly valid, sometimes invalid for testing
func (g *OpcodeGenerator) generateValidBLS12G1PointUncompressed() []byte {
	point := make([]byte, 128) // 128 bytes uncompressed format: 64 bytes x + 64 bytes y

	// 3% chance to generate invalid G1 points for testing edge cases
	if g.rng.Float64() < 0.03 {
		invalidType := g.rng.Intn(4)
		switch invalidType {
		case 0:
			// Invalid field elements (> modulus)
			copy(point[:64], g.rng.Bytes(64))
			copy(point[64:], g.rng.Bytes(64))
			point[0] |= 0xE0  // Ensure x >= modulus
			point[64] |= 0xE0 // Ensure y >= modulus
		case 1:
			// All 0xFF bytes
			for i := range point {
				point[i] = 0xFF
			}
		case 2:
			// Mixed valid x, invalid y
			xCoord := g.generateValidBLS12BaseField()
			copy(point[:64], xCoord)
			copy(point[64:], g.rng.Bytes(64))
			point[64] |= 0xE0 // Invalid y
		default:
			// Invalid point format (random high bits set in wrong places)
			xCoord := g.generateValidBLS12BaseField()
			yCoord := g.generateValidBLS12BaseField()
			copy(point[:64], xCoord)
			copy(point[64:], yCoord)
			// Set compression/infinity flags incorrectly
			point[0] |= 0x80 // Set compression flag (should be 0 for uncompressed)
		}
		return point
	}

	// 10% chance of zero point (point at infinity)
	if g.rng.Float64() < 0.1 {
		// Zero point: all zeros represents point at infinity in uncompressed format
		return point
	}

	// Generate mostly valid uncompressed G1 point using proper field elements
	patternType := g.rng.Intn(3)
	switch patternType {
	case 0:
		// Use generator point pattern with valid field elements
		xCoord := g.generateValidBLS12BaseField()
		yCoord := g.generateValidBLS12BaseField()
		copy(point[:64], xCoord)
		copy(point[64:], yCoord)
	case 1:
		// Small coordinate values within field
		xCoord := make([]byte, 64)
		yCoord := make([]byte, 64)
		binary.BigEndian.PutUint64(xCoord[56:64], uint64(g.rng.Intn(1000)+1)) // x coordinate
		binary.BigEndian.PutUint64(yCoord[56:64], uint64(g.rng.Intn(1000)+1)) // y coordinate
		copy(point[:64], xCoord)
		copy(point[64:], yCoord)
	default:
		// Valid random field elements
		xCoord := g.generateValidBLS12BaseField()
		yCoord := g.generateValidBLS12BaseField()
		copy(point[:64], xCoord)
		copy(point[64:], yCoord)
	}

	return point
}

// generateValidBLS12G2PointUncompressed generates a BLS12-381 G2 point in uncompressed format (256 bytes) - mostly valid, sometimes invalid for testing
func (g *OpcodeGenerator) generateValidBLS12G2PointUncompressed() []byte {
	point := make([]byte, 256) // 256 bytes uncompressed format: 64+64+64+64 bytes for x_c0,x_c1,y_c0,y_c1

	// 3% chance to generate invalid G2 points for testing edge cases
	if g.rng.Float64() < 0.03 {
		invalidType := g.rng.Intn(5)
		switch invalidType {
		case 0:
			// Invalid field elements (> modulus)
			copy(point, g.rng.Bytes(256))
			for i := 0; i < 4; i++ {
				point[i*64] |= 0xE0 // Ensure all field elements >= modulus
			}
		case 1:
			// All 0xFF bytes
			for i := range point {
				point[i] = 0xFF
			}
		case 2:
			// Mixed valid and invalid field elements
			xC0 := g.generateValidBLS12BaseField()
			xC1 := g.generateValidBLS12BaseField()
			copy(point[0:64], xC0)
			copy(point[64:128], xC1)
			// Invalid y coordinates
			copy(point[128:192], g.rng.Bytes(64))
			copy(point[192:256], g.rng.Bytes(64))
			point[128] |= 0xE0 // Invalid y_c0
			point[192] |= 0xE0 // Invalid y_c1
		case 3:
			// Invalid point format flags
			xC0 := g.generateValidBLS12BaseField()
			xC1 := g.generateValidBLS12BaseField()
			yC0 := g.generateValidBLS12BaseField()
			yC1 := g.generateValidBLS12BaseField()
			copy(point[0:64], xC0)
			copy(point[64:128], xC1)
			copy(point[128:192], yC0)
			copy(point[192:256], yC1)
			// Set compression flag incorrectly
			point[0] |= 0x80 // Set compression flag (should be 0 for uncompressed)
		default:
			// Wrong field element ordering
			xC0 := g.generateValidBLS12BaseField()
			xC1 := g.generateValidBLS12BaseField()
			yC0 := g.generateValidBLS12BaseField()
			yC1 := g.generateValidBLS12BaseField()
			// Deliberately swap field elements to test validation
			copy(point[0:64], yC1)    // Wrong position
			copy(point[64:128], yC0)  // Wrong position
			copy(point[128:192], xC1) // Wrong position
			copy(point[192:256], xC0) // Wrong position
		}
		return point
	}

	// 10% chance of zero point (point at infinity)
	if g.rng.Float64() < 0.1 {
		// Zero point: all zeros represents point at infinity in uncompressed format
		return point
	}

	// Generate mostly valid uncompressed G2 point using proper field elements
	patternType := g.rng.Intn(3)
	switch patternType {
	case 0:
		// Use generator point pattern with valid field elements
		xC0 := g.generateValidBLS12BaseField()
		xC1 := g.generateValidBLS12BaseField()
		yC0 := g.generateValidBLS12BaseField()
		yC1 := g.generateValidBLS12BaseField()
		copy(point[0:64], xC0)    // x_c0
		copy(point[64:128], xC1)  // x_c1
		copy(point[128:192], yC0) // y_c0
		copy(point[192:256], yC1) // y_c1
	case 1:
		// Small coordinate values within field
		xC0 := make([]byte, 64)
		xC1 := make([]byte, 64)
		yC0 := make([]byte, 64)
		yC1 := make([]byte, 64)
		binary.BigEndian.PutUint64(xC0[56:64], uint64(g.rng.Intn(1000)+1)) // x_c0
		binary.BigEndian.PutUint64(xC1[56:64], uint64(g.rng.Intn(1000)+1)) // x_c1
		binary.BigEndian.PutUint64(yC0[56:64], uint64(g.rng.Intn(1000)+1)) // y_c0
		binary.BigEndian.PutUint64(yC1[56:64], uint64(g.rng.Intn(1000)+1)) // y_c1
		copy(point[0:64], xC0)                                             // x_c0
		copy(point[64:128], xC1)                                           // x_c1
		copy(point[128:192], yC0)                                          // y_c0
		copy(point[192:256], yC1)                                          // y_c1
	default:
		// Valid random field elements
		xC0 := g.generateValidBLS12BaseField()
		xC1 := g.generateValidBLS12BaseField()
		yC0 := g.generateValidBLS12BaseField()
		yC1 := g.generateValidBLS12BaseField()
		copy(point[0:64], xC0)    // x_c0
		copy(point[64:128], xC1)  // x_c1
		copy(point[128:192], yC0) // y_c0
		copy(point[192:256], yC1) // y_c1
	}

	return point
}

// generateBLS12MapFpToG1Call creates a BLS12-381 map field point to G1 call
func (g *OpcodeGenerator) generateBLS12MapFpToG1Call() []byte {
	var bytecode []byte

	// Generate valid BLS12-381 base field element (64 bytes)
	fieldElement := g.generateValidBLS12BaseField()

	// Write field element using 2 MSTORE operations (64 bytes = 2 × 32 bytes)
	bytecode = append(bytecode, 0x7f) // PUSH32
	bytecode = append(bytecode, fieldElement[:32]...)
	bytecode = append(bytecode, 0x60, 0x00) // PUSH1 0
	bytecode = append(bytecode, 0x52)       // MSTORE

	bytecode = append(bytecode, 0x7f) // PUSH32
	bytecode = append(bytecode, fieldElement[32:]...)
	bytecode = append(bytecode, 0x60, 0x20) // PUSH1 32
	bytecode = append(bytecode, 0x52)       // MSTORE

	// Setup CALL to BLS12_MAP_FP_TO_G1 precompile
	bytecode = append(bytecode, 0x60, 0x80) // PUSH1 128 (return size - uncompressed G1 point)
	bytecode = append(bytecode, 0x60, 0x40) // PUSH1 64 (return offset)
	bytecode = append(bytecode, 0x60, 0x40) // PUSH1 64 (args size)
	bytecode = append(bytecode, 0x60, 0x00) // PUSH1 0 (args offset)
	bytecode = append(bytecode, 0x60, 0x00) // PUSH1 0 (value)
	bytecode = append(bytecode, 0x60, 0x10) // PUSH1 16 (precompile address)
	bytecode = append(bytecode, 0x5a)       // GAS
	bytecode = append(bytecode, 0xf1)       // CALL

	// Add LOG0 to log the result when in precompiles-only mode (return at 0x40, size 128)
	bytecode = append(bytecode, g.addLogIfPrecompilesOnly(0x40, 128)...)

	return bytecode
}

// generateBLS12G1AddCall creates a BLS12-381 G1 point addition call
func (g *OpcodeGenerator) generateBLS12G1AddCall() []byte {
	var bytecode []byte

	// BLS12_G1ADD input: 256 bytes (64+64+64+64 bytes for x1,y1,x2,y2)
	// Generate two valid BLS12 G1 points (128 bytes each, uncompressed format)
	point1 := g.generateValidBLS12G1PointUncompressed()
	point2 := g.generateValidBLS12G1PointUncompressed()

	// Write points to memory using 8 MSTORE operations (256 bytes = 8 × 32 bytes)
	for i := 0; i < 8; i++ {
		var data []byte
		if i < 4 {
			// First point (128 bytes = 4 × 32 bytes)
			data = point1[i*32 : (i+1)*32]
		} else {
			// Second point (128 bytes = 4 × 32 bytes)
			data = point2[(i-4)*32 : (i-3)*32]
		}

		bytecode = append(bytecode, 0x7f) // PUSH32
		bytecode = append(bytecode, data...)
		offset := i * 32
		if offset < 256 {
			bytecode = append(bytecode, 0x60, byte(offset)) // PUSH1 offset
		} else {
			bytecode = append(bytecode, 0x61, byte(offset>>8), byte(offset)) // PUSH2 offset
		}
		bytecode = append(bytecode, 0x52) // MSTORE
	}

	// Setup CALL to BLS12_G1ADD precompile
	bytecode = append(bytecode, 0x60, 0x80)       // PUSH1 128 (return size - uncompressed G1 point)
	bytecode = append(bytecode, 0x61, 0x01, 0x00) // PUSH2 256 (return offset)
	bytecode = append(bytecode, 0x61, 0x01, 0x00) // PUSH2 256 (args size - two 128-byte points)
	bytecode = append(bytecode, 0x60, 0x00)       // PUSH1 0 (args offset)
	bytecode = append(bytecode, 0x60, 0x00)       // PUSH1 0 (value)
	bytecode = append(bytecode, 0x60, 0x0b)       // PUSH1 11 (precompile address)
	bytecode = append(bytecode, 0x5a)             // GAS
	bytecode = append(bytecode, 0xf1)             // CALL

	// Add LOG0 to log the result when in precompiles-only mode (return at 0x100, size 128)
	bytecode = append(bytecode, g.addLogIfPrecompilesOnly(0x100, 128)...)

	return bytecode
}

// generateBLS12G2AddCall creates a BLS12-381 G2 point addition call
func (g *OpcodeGenerator) generateBLS12G2AddCall() []byte {
	var bytecode []byte

	// BLS12_G2ADD input: 512 bytes (64+64+64+64+64+64+64+64 bytes for x1_c0,x1_c1,y1_c0,y1_c1,x2_c0,x2_c1,y2_c0,y2_c1)
	// Generate two valid BLS12 G2 points (256 bytes each, uncompressed format)
	point1 := g.generateValidBLS12G2PointUncompressed()
	point2 := g.generateValidBLS12G2PointUncompressed()

	// Write points to memory using 16 MSTORE operations (512 bytes = 16 × 32 bytes)
	for i := 0; i < 16; i++ {
		var data []byte
		if i < 8 {
			// First point (256 bytes = 8 × 32 bytes)
			data = point1[i*32 : (i+1)*32]
		} else {
			// Second point (256 bytes = 8 × 32 bytes)
			data = point2[(i-8)*32 : (i-7)*32]
		}

		bytecode = append(bytecode, 0x7f) // PUSH32
		bytecode = append(bytecode, data...)
		offset := i * 32
		if offset < 256 {
			bytecode = append(bytecode, 0x60, byte(offset)) // PUSH1 offset
		} else {
			bytecode = append(bytecode, 0x61, byte(offset>>8), byte(offset)) // PUSH2 offset
		}
		bytecode = append(bytecode, 0x52) // MSTORE
	}

	// Setup CALL to BLS12_G2ADD precompile
	bytecode = append(bytecode, 0x61, 0x01, 0x00) // PUSH2 256 (return size - uncompressed G2 point)
	bytecode = append(bytecode, 0x61, 0x02, 0x00) // PUSH2 512 (return offset)
	bytecode = append(bytecode, 0x61, 0x02, 0x00) // PUSH2 512 (args size - two 256-byte points)
	bytecode = append(bytecode, 0x60, 0x00)       // PUSH1 0 (args offset)
	bytecode = append(bytecode, 0x60, 0x00)       // PUSH1 0 (value)
	bytecode = append(bytecode, 0x60, 0x0d)       // PUSH1 13 (precompile address)
	bytecode = append(bytecode, 0x5a)             // GAS
	bytecode = append(bytecode, 0xf1)             // CALL

	// Add LOG0 to log the result when in precompiles-only mode (return at 0x200, size 256)
	bytecode = append(bytecode, g.addLogIfPrecompilesOnly(0x200, 256)...)

	return bytecode
}

// generateBLS12PairingCall creates a BLS12-381 pairing check call
func (g *OpcodeGenerator) generateBLS12PairingCall() []byte {
	var bytecode []byte

	// Generate 1-2 pairs of G1 and G2 points
	pairCount := 1 + g.rng.Intn(2) // 1 or 2 pairs
	totalLen := pairCount * 384    // Each pair: 128 bytes (G1 uncompressed) + 256 bytes (G2 uncompressed)

	// Write pairs to memory using optimal MSTORE operations
	for i := 0; i < pairCount; i++ {
		memOffset := i * 384

		// Generate G1 point (128 bytes uncompressed)
		g1Point := g.generateValidBLS12G1PointUncompressed()

		// Write G1 point using 4 MSTORE operations (128 bytes = 4 × 32 bytes)
		for j := 0; j < 4; j++ {
			bytecode = append(bytecode, 0x7f) // PUSH32
			bytecode = append(bytecode, g1Point[j*32:(j+1)*32]...)
			g1Offset := memOffset + j*32
			if g1Offset < 256 {
				bytecode = append(bytecode, 0x60, byte(g1Offset)) // PUSH1 offset
			} else {
				bytecode = append(bytecode, 0x61, byte(g1Offset>>8), byte(g1Offset)) // PUSH2 offset
			}
			bytecode = append(bytecode, 0x52) // MSTORE
		}

		// Generate G2 point (256 bytes uncompressed)
		g2Point := g.generateValidBLS12G2PointUncompressed()

		// Write G2 point using 8 MSTORE operations (256 bytes = 8 × 32 bytes)
		for j := 0; j < 8; j++ {
			bytecode = append(bytecode, 0x7f) // PUSH32
			bytecode = append(bytecode, g2Point[j*32:(j+1)*32]...)
			g2Offset := memOffset + 128 + j*32
			if g2Offset < 256 {
				bytecode = append(bytecode, 0x60, byte(g2Offset)) // PUSH1 offset
			} else {
				bytecode = append(bytecode, 0x61, byte(g2Offset>>8), byte(g2Offset)) // PUSH2 offset
			}
			bytecode = append(bytecode, 0x52) // MSTORE
		}
	}

	// Setup CALL to BLS12_PAIRING precompile
	bytecode = append(bytecode, 0x60, 0x20) // PUSH1 32 (return size)
	returnOffset := totalLen + 32
	if returnOffset < 256 {
		bytecode = append(bytecode, 0x60, byte(returnOffset)) // PUSH1 retOffset
	} else {
		bytecode = append(bytecode, 0x61, byte(returnOffset>>8), byte(returnOffset)) // PUSH2 retOffset
	}
	if totalLen < 256 {
		bytecode = append(bytecode, 0x60, byte(totalLen)) // PUSH1 totalLen
	} else {
		bytecode = append(bytecode, 0x61, byte(totalLen>>8), byte(totalLen)) // PUSH2 totalLen
	}
	bytecode = append(bytecode, 0x60, 0x00) // PUSH1 0 (args offset)
	bytecode = append(bytecode, 0x60, 0x00) // PUSH1 0 (value)
	bytecode = append(bytecode, 0x60, 0x0f) // PUSH1 15 (precompile address)
	bytecode = append(bytecode, 0x5a)       // GAS
	bytecode = append(bytecode, 0xf1)       // CALL

	// Add LOG0 to log the result when in precompiles-only mode (return at returnOffset, size 32)
	bytecode = append(bytecode, g.addLogIfPrecompilesOnly(returnOffset, 32)...)

	return bytecode
}

// generateBLS12MapFp2G2Call creates a BLS12-381 map field point to G2 call
func (g *OpcodeGenerator) generateBLS12MapFp2G2Call() []byte {
	var bytecode []byte

	// Generate field element (fp2: 2 * 64 bytes = 128 bytes total)
	// Each component must be a valid base field element
	fp2Element := make([]byte, 128)
	c0 := g.generateValidBLS12BaseField()
	c1 := g.generateValidBLS12BaseField()
	copy(fp2Element[:64], c0)
	copy(fp2Element[64:], c1)

	// Write field element to memory using 4 MSTORE operations (128 bytes = 4 × 32 bytes)
	for i := 0; i < 4; i++ {
		bytecode = append(bytecode, 0x7f) // PUSH32
		bytecode = append(bytecode, fp2Element[i*32:(i+1)*32]...)
		offset := i * 32
		bytecode = append(bytecode, 0x60, byte(offset)) // PUSH1 offset
		bytecode = append(bytecode, 0x52)               // MSTORE
	}

	// Setup CALL to BLS12_MAP_FP2_TO_G2 precompile
	bytecode = append(bytecode, 0x61, 0x01, 0x00) // PUSH2 256 (return size - uncompressed G2 point)
	bytecode = append(bytecode, 0x60, 0x80)       // PUSH1 128 (return offset)
	bytecode = append(bytecode, 0x60, 0x80)       // PUSH1 128 (args size)
	bytecode = append(bytecode, 0x60, 0x00)       // PUSH1 0 (args offset)
	bytecode = append(bytecode, 0x60, 0x00)       // PUSH1 0 (value)
	bytecode = append(bytecode, 0x60, 0x11)       // PUSH1 17 (precompile address)
	bytecode = append(bytecode, 0x5a)             // GAS
	bytecode = append(bytecode, 0xf1)             // CALL

	// Add LOG0 to log the result when in precompiles-only mode (return at 0x80, size 256)
	bytecode = append(bytecode, g.addLogIfPrecompilesOnly(0x80, 256)...)

	return bytecode
}

// generateBLS12G1MSMCall creates a BLS12-381 G1 multi-scalar multiplication call
func (g *OpcodeGenerator) generateBLS12G1MSMCall() []byte {
	var bytecode []byte

	// G1 MSM input format: pairs of (G1_point, scalar)
	// Each G1 point: 128 bytes (64 bytes x + 64 bytes y, uncompressed format)
	// Each scalar: 32 bytes
	// Total per pair: 160 bytes

	// Generate 2-3 pairs (fewer due to larger size)
	pairCount := 2 + g.rng.Intn(2) // 2 or 3 pairs
	totalLen := pairCount * 160

	// Write pairs to memory
	for i := 0; i < pairCount; i++ {
		memOffset := i * 160

		// Generate valid G1 point (128 bytes uncompressed)
		pointData := g.generateValidBLS12G1PointUncompressed()

		// Write G1 point using 4 MSTORE operations (128 bytes = 4 × 32 bytes)
		for j := 0; j < 4; j++ {
			bytecode = append(bytecode, 0x7f) // PUSH32
			bytecode = append(bytecode, pointData[j*32:(j+1)*32]...)
			coordOffset := memOffset + j*32
			if coordOffset < 256 {
				bytecode = append(bytecode, 0x60, byte(coordOffset)) // PUSH1 offset
			} else {
				bytecode = append(bytecode, 0x61, byte(coordOffset>>8), byte(coordOffset)) // PUSH2 offset
			}
			bytecode = append(bytecode, 0x52) // MSTORE
		}

		// Generate scalar (32 bytes)
		scalarData := g.rng.Bytes(32)
		bytecode = append(bytecode, 0x7f) // PUSH32
		bytecode = append(bytecode, scalarData...)
		scalarOffset := memOffset + 128 // Point is 128 bytes, scalar starts after
		if scalarOffset < 256 {
			bytecode = append(bytecode, 0x60, byte(scalarOffset)) // PUSH1 offset
		} else {
			bytecode = append(bytecode, 0x61, byte(scalarOffset>>8), byte(scalarOffset)) // PUSH2 offset
		}
		bytecode = append(bytecode, 0x52) // MSTORE
	}

	// Setup CALL to BLS12_G1MSM precompile
	bytecode = append(bytecode, 0x60, 0x80) // PUSH1 128 (return size - uncompressed G1 point)
	returnOffset := totalLen + 32
	if returnOffset < 256 {
		bytecode = append(bytecode, 0x60, byte(returnOffset)) // PUSH1 retOffset
	} else {
		bytecode = append(bytecode, 0x61, byte(returnOffset>>8), byte(returnOffset)) // PUSH2 retOffset
	}
	if totalLen < 256 {
		bytecode = append(bytecode, 0x60, byte(totalLen)) // PUSH1 totalLen
	} else {
		bytecode = append(bytecode, 0x61, byte(totalLen>>8), byte(totalLen)) // PUSH2 totalLen
	}
	bytecode = append(bytecode, 0x60, 0x00) // PUSH1 0 (args offset)
	bytecode = append(bytecode, 0x60, 0x00) // PUSH1 0 (value)
	bytecode = append(bytecode, 0x60, 0x0c) // PUSH1 12 (precompile address)
	bytecode = append(bytecode, 0x5a)       // GAS
	bytecode = append(bytecode, 0xf1)       // CALL

	// Add LOG0 to log the result when in precompiles-only mode (return at returnOffset, size 128)
	bytecode = append(bytecode, g.addLogIfPrecompilesOnly(returnOffset, 128)...)

	return bytecode
}

// generateBLS12G2MSMCall creates a BLS12-381 G2 multi-scalar multiplication call
func (g *OpcodeGenerator) generateBLS12G2MSMCall() []byte {
	var bytecode []byte

	// G2 MSM input format: pairs of (G2_point, scalar)
	// Each G2 point: 256 bytes (64+64+64+64 bytes for x_c0,x_c1,y_c0,y_c1, uncompressed format)
	// Each scalar: 32 bytes
	// Total per pair: 288 bytes

	// Generate 2 pairs (fewer due to very large size)
	pairCount := 2 // 2 pairs
	totalLen := pairCount * 288

	// Write pairs to memory using optimal MSTORE operations
	for i := 0; i < pairCount; i++ {
		memOffset := i * 288

		// Generate valid G2 point (256 bytes uncompressed)
		pointData := g.generateValidBLS12G2PointUncompressed()

		// Write G2 point using 8 MSTORE operations (256 bytes = 8 × 32 bytes)
		for j := 0; j < 8; j++ {
			bytecode = append(bytecode, 0x7f) // PUSH32
			bytecode = append(bytecode, pointData[j*32:(j+1)*32]...)
			coordOffset := memOffset + j*32
			if coordOffset < 256 {
				bytecode = append(bytecode, 0x60, byte(coordOffset)) // PUSH1 offset
			} else {
				bytecode = append(bytecode, 0x61, byte(coordOffset>>8), byte(coordOffset)) // PUSH2 offset
			}
			bytecode = append(bytecode, 0x52) // MSTORE
		}

		// Generate scalar (32 bytes)
		scalarData := g.rng.Bytes(32)
		bytecode = append(bytecode, 0x7f) // PUSH32
		bytecode = append(bytecode, scalarData...)
		scalarOffset := memOffset + 256
		if scalarOffset < 256 {
			bytecode = append(bytecode, 0x60, byte(scalarOffset)) // PUSH1 offset
		} else {
			bytecode = append(bytecode, 0x61, byte(scalarOffset>>8), byte(scalarOffset)) // PUSH2 offset
		}
		bytecode = append(bytecode, 0x52) // MSTORE
	}

	// Setup CALL to BLS12_G2MSM precompile
	bytecode = append(bytecode, 0x61, 0x01, 0x00) // PUSH2 256 (return size - uncompressed G2 point)
	returnOffset := totalLen + 32
	if returnOffset < 256 {
		bytecode = append(bytecode, 0x60, byte(returnOffset)) // PUSH1 retOffset
	} else {
		bytecode = append(bytecode, 0x61, byte(returnOffset>>8), byte(returnOffset)) // PUSH2 retOffset
	}
	if totalLen < 256 {
		bytecode = append(bytecode, 0x60, byte(totalLen)) // PUSH1 totalLen
	} else {
		bytecode = append(bytecode, 0x61, byte(totalLen>>8), byte(totalLen)) // PUSH2 totalLen
	}
	bytecode = append(bytecode, 0x60, 0x00) // PUSH1 0 (args offset)
	bytecode = append(bytecode, 0x60, 0x00) // PUSH1 0 (value)
	bytecode = append(bytecode, 0x60, 0x0e) // PUSH1 14 (precompile address)
	bytecode = append(bytecode, 0x5a)       // GAS
	bytecode = append(bytecode, 0xf1)       // CALL

	// Add LOG0 to log the result when in precompiles-only mode (return at returnOffset, size 256)
	bytecode = append(bytecode, g.addLogIfPrecompilesOnly(returnOffset, 256)...)

	return bytecode
}

// addLogIfPrecompilesOnly adds LOG0 to log precompile results when in precompiles-only mode
func (g *OpcodeGenerator) addLogIfPrecompilesOnly(memOffset, size int) []byte {
	// Only add LOG0 when fuzzing precompiles only
	if g.fuzzMode != "precompiles" {
		return []byte{}
	}

	var bytecode []byte

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

	return bytecode
}
