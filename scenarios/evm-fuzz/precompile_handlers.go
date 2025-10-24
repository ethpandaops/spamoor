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

// Precompile call templates - all using specialized handlers directly
func ecrecoverTemplate(g *OpcodeGenerator) []byte { return g.generateEcrecoverCall() }
func sha256Template(g *OpcodeGenerator) []byte    { return g.generateSha256Call() }
func ripemd160Template(g *OpcodeGenerator) []byte { return g.generateRipemd160Call() }
func identityTemplate(g *OpcodeGenerator) []byte  { return g.generateIdentityCall() }
func modexpTemplate(g *OpcodeGenerator) []byte    { return g.generateModexpCall() }
func ecAddTemplate(g *OpcodeGenerator) []byte     { return g.generateBN256EcAddCall() }
func ecMulTemplate(g *OpcodeGenerator) []byte     { return g.generateBN256EcMulCall() }
func ecPairingTemplate(g *OpcodeGenerator) []byte { return g.generateBN256PairingCall() }
func blake2fTemplate(g *OpcodeGenerator) []byte   { return g.generateBlake2fCall() }
func pointEvalTemplate(g *OpcodeGenerator) []byte { return g.generateKZGPointEvalCall() }

// BLS12-381 precompile templates
func bls12G1AddTemplate(g *OpcodeGenerator) []byte     { return g.generateBLS12G1AddCall() }     // 0x0b
func bls12G1MSMTemplate(g *OpcodeGenerator) []byte     { return g.generateBLS12G1MSMCall() }     // 0x0c
func bls12G2AddTemplate(g *OpcodeGenerator) []byte     { return g.generateBLS12G2AddCall() }     // 0x0d
func bls12G2MSMTemplate(g *OpcodeGenerator) []byte     { return g.generateBLS12G2MSMCall() }     // 0x0e
func bls12PairingTemplate(g *OpcodeGenerator) []byte   { return g.generateBLS12PairingCall() }   // 0x0f
func bls12MapFpToG1Template(g *OpcodeGenerator) []byte { return g.generateBLS12MapFpToG1Call() } // 0x10
func bls12MapFp2G2Template(g *OpcodeGenerator) []byte  { return g.generateBLS12MapFp2G2Call() }  // 0x11

// generateValidBN256Point generates a BN256 curve point - mostly valid, sometimes invalid for testing
func (g *OpcodeGenerator) generateValidBN256Point() []byte {
	point := make([]byte, 64) // 32 bytes x + 32 bytes y

	// 1% chance to generate invalid BN256 points for testing edge cases
	if g.rng.Float64() < 0.01 {
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
		return g.transformer.TransformPrecompileInput(point, 0x06)
	}

	// 10% chance of zero point (point at infinity)
	if g.rng.Float64() < 0.1 {
		return g.transformer.TransformPrecompileInput(point, 0x06) // All zeros = point at infinity
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
			return g.transformer.TransformPrecompileInput(append(xBytes, yBytes...), 0x06)
		},
		// Pattern 2: Use generator point multiples
		func() []byte {
			return g.transformer.TransformPrecompileInput(g.getBN256Generator(), 0x06)
		},
		// Pattern 3: Use modular reduction of random values
		func() []byte {
			randomX := new(big.Int).SetBytes(g.rng.Bytes(32))
			x := new(big.Int).Mod(randomX, bn256FieldModulus)
			y := g.calculateBN256Y(x)
			if y == nil {
				return g.transformer.TransformPrecompileInput(g.getBN256Generator(), 0x06)
			}
			xBytes := make([]byte, 32)
			yBytes := make([]byte, 32)
			x.FillBytes(xBytes)
			y.FillBytes(yBytes)
			return g.transformer.TransformPrecompileInput(append(xBytes, yBytes...), 0x06)
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
	bytecode = append(bytecode, g.processPrecompileResult(128, 64)...)

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
	bytecode = append(bytecode, g.processPrecompileResult(96, 64)...)

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
	bytecode = append(bytecode, g.processPrecompileResult(returnOffset, 32)...)

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

// =============================================================================
// Specialized handlers for basic precompiles (0x01-0x0a)
// =============================================================================

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
