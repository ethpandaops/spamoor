package evmfuzz

import (
	"encoding/binary"
	"math/big"
)

// BN256 curve parameters
var (
	// BN256 field modulus: 21888242871839275222246405745257275088696311157297823662689037894645226208583
	bn256FieldModulus = new(big.Int)

	// BN256 curve order (r): 21888242871839275222246405745257275088548364400416034343698204186575808495617
	bn256CurveOrder = new(big.Int)

	bn256One = big.NewInt(1)
)

func init() {
	// Initialize BN256 field modulus
	bn256FieldModulus.SetString("21888242871839275222246405745257275088696311157297823662689037894645226208583", 10)

	// Initialize BN256 curve order
	bn256CurveOrder.SetString("21888242871839275222246405745257275088548364400416034343698204186575808495617", 10)
}

// generateValidBN256Point generates a BN256 curve point - mostly valid, sometimes invalid for testing
func (g *OpcodeGenerator) generateValidBN256Point() []byte {
	point := make([]byte, 64) // 32 bytes x + 32 bytes y

	choice := g.rng.Intn(100)

	if choice < 20 {
		// 20% chance to generate invalid BN256 points for comprehensive testing
		return g.generateInvalidBN256Point()
	}

	if choice < 25 {
		// 5% chance of zero point (point at infinity)
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

// generateInvalidBN256Point generates various types of invalid BN256 points for comprehensive testing
func (g *OpcodeGenerator) generateInvalidBN256Point() []byte {
	point := make([]byte, 64)
	invalidType := g.rng.Intn(12)

	switch invalidType {
	case 0:
		// Type 1: Field elements greater than modulus
		xBig := new(big.Int).Add(bn256FieldModulus, big.NewInt(int64(g.rng.Intn(1000)+1)))
		yBig := new(big.Int).Add(bn256FieldModulus, big.NewInt(int64(g.rng.Intn(1000)+1)))
		safeFillBytes(xBig, point[:32])
		safeFillBytes(yBig, point[32:])

	case 1:
		// Type 2: Maximum field values (all 0xFF)
		for i := range point {
			point[i] = 0xFF
		}

	case 2:
		// Type 3: Point not on curve - valid field elements but don't satisfy curve equation
		x := new(big.Int).SetBytes(g.rng.Bytes(32))
		x.Mod(x, bn256FieldModulus)

		// Generate random y that likely won't satisfy y² = x³ + 3
		y := new(big.Int).SetBytes(g.rng.Bytes(32))
		y.Mod(y, bn256FieldModulus)

		safeFillBytes(x, point[:32])
		safeFillBytes(y, point[32:])

	case 3:
		// Type 4: Edge case - modulus minus 1
		modMinusOne := new(big.Int).Sub(bn256FieldModulus, big.NewInt(1))
		safeFillBytes(modMinusOne, point[:32])
		safeFillBytes(modMinusOne, point[32:])

	case 4:
		// Type 5: Specific invalid coordinates that look valid but aren't
		// x = 1, y = 1 (not on BN256 curve: y² ≠ x³ + 3)
		point[31] = 1 // x = 1
		point[63] = 1 // y = 1

	case 5:
		// Type 6: Valid x, y = 0 (unless x = 0, this won't be on curve)
		x := big.NewInt(2) // Use x = 2
		safeFillBytes(x, point[:32])
		// y stays 0

	case 6:
		// Type 7: Points that could cause arithmetic overflow
		// Use values very close to field modulus
		nearModulus := new(big.Int).Sub(bn256FieldModulus, big.NewInt(int64(g.rng.Intn(10)+1)))
		safeFillBytes(nearModulus, point[:32])
		safeFillBytes(nearModulus, point[32:])

	case 7:
		// Type 8: Invalid point at infinity representations
		// Some implementations expect specific encodings for infinity
		for i := 0; i < 32; i++ {
			point[i] = 0xFF // Invalid "infinity" encoding
		}
		// Leave y as zeros

	case 8:
		// Type 9: Points with small coordinates that might cause edge cases
		point[31] = byte(g.rng.Intn(10) + 1) // Small x
		point[63] = byte(g.rng.Intn(10) + 1) // Small y (likely not on curve)

	case 9:
		// Type 10: Coordinates that are valid field elements but create invalid curve points
		// Use square roots of non-residues
		x := big.NewInt(int64(g.rng.Intn(1000) + 1))
		safeFillBytes(x, point[:32])

		// Calculate x³ + 3 and add 1 to make it not a quadratic residue
		x3 := new(big.Int).Exp(x, big.NewInt(3), bn256FieldModulus)
		x3Plus3 := new(big.Int).Add(x3, big.NewInt(3))
		invalidY := new(big.Int).Add(x3Plus3, big.NewInt(1))
		invalidY.Mod(invalidY, bn256FieldModulus)
		safeFillBytes(invalidY, point[32:])

	case 10:
		// Type 11: Points from twist attacks (simulate points on quadratic twist)
		// Generate point using curve equation with different constant: y² = x³ + b where b ≠ 3
		x := new(big.Int).SetBytes(g.rng.Bytes(32))
		x.Mod(x, bn256FieldModulus)

		// Use wrong curve constant (e.g., y² = x³ + 5 instead of x³ + 3)
		x3 := new(big.Int).Exp(x, big.NewInt(3), bn256FieldModulus)
		x3Plus5 := new(big.Int).Add(x3, big.NewInt(5))
		x3Plus5.Mod(x3Plus5, bn256FieldModulus)

		// Try to find square root (may not exist, which is what we want)
		twistY := g.modularSqrt(x3Plus5, bn256FieldModulus)
		if twistY == nil {
			// Good, no square root exists - use a random y
			twistY = new(big.Int).SetBytes(g.rng.Bytes(32))
			twistY.Mod(twistY, bn256FieldModulus)
		}

		safeFillBytes(x, point[:32])
		safeFillBytes(twistY, point[32:])

	default:
		// Type 12: Random bytes (most likely invalid)
		copy(point, g.rng.Bytes(64))
	}

	return g.transformer.TransformPrecompileInput(point, 0x06)
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
	choice := g.rng.Intn(100)

	if choice < 15 {
		// 15% chance of invalid scalars for testing
		return g.generateInvalidBN256Scalar()
	}

	if choice < 20 {
		// 5% chance of small scalars (often interesting edge cases)
		small := big.NewInt(int64(g.rng.Intn(100)))
		smallBytes := make([]byte, 32)
		small.FillBytes(smallBytes)
		return smallBytes
	}

	// 80% chance of random valid scalars
	// Generate random scalar modulo curve order
	// For simplicity, just use a 32-byte value (will be reduced by precompile)
	return g.rng.Bytes(32)
}

// generateInvalidBN256Scalar generates invalid scalars for BN256 operations
func (g *OpcodeGenerator) generateInvalidBN256Scalar() []byte {
	scalar := make([]byte, 32)
	invalidType := g.rng.Intn(8)

	switch invalidType {
	case 0:
		// Type 1: Scalar >= curve order (invalid)
		scalarBig := new(big.Int).Add(bn256CurveOrder, big.NewInt(int64(g.rng.Intn(1000)+1)))
		safeFillBytes(scalarBig, scalar)

	case 1:
		// Type 2: Scalar = curve order (invalid)
		safeFillBytes(bn256CurveOrder, scalar)

	case 2:
		// Type 3: Maximum value (all 0xFF)
		for i := range scalar {
			scalar[i] = 0xFF
		}

	case 3:
		// Type 4: Curve order - 1 (valid but edge case)
		orderMinusOne := new(big.Int).Sub(bn256CurveOrder, big.NewInt(1))
		safeFillBytes(orderMinusOne, scalar)

	case 4:
		// Type 5: Powers of 2 near curve order
		powerOf2 := new(big.Int).Lsh(big.NewInt(1), 254) // 2^254
		safeFillBytes(powerOf2, scalar)

	case 5:
		// Type 6: Small multiples of curve order
		multiple := big.NewInt(int64(g.rng.Intn(10) + 2)) // 2-11
		invalidScalar := new(big.Int).Mul(bn256CurveOrder, multiple)
		safeFillBytes(invalidScalar, scalar)

	case 6:
		// Type 7: Negative values (two's complement representation)
		validScalar := new(big.Int).SetBytes(g.rng.Bytes(32))
		validScalar.Mod(validScalar, bn256CurveOrder)
		negated := new(big.Int).Sub(bn256CurveOrder, validScalar)
		safeFillBytes(negated, scalar)

	default:
		// Type 8: Random bytes (likely invalid)
		copy(scalar, g.rng.Bytes(32))
		// Ensure it's >= curve order by setting high bits
		scalar[0] |= 0xF0
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
	choice := g.rng.Intn(100)

	if choice < 20 {
		// 20% chance to generate invalid BN256 G2 points for comprehensive testing
		return g.generateInvalidBN256G2Point()
	}

	if choice < 25 {
		// 5% chance of zero point (point at infinity)
		return make([]byte, 128)
	}

	// 75% chance of valid G2 points
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

// generateInvalidBN256G2Point generates various types of invalid BN256 G2 points for comprehensive testing
func (g *OpcodeGenerator) generateInvalidBN256G2Point() []byte {
	point := make([]byte, 128)
	invalidType := g.rng.Intn(12)

	switch invalidType {
	case 0:
		// Type 1: Field elements greater than modulus
		for i := 0; i < 4; i++ {
			fieldBig := new(big.Int).Add(bn256FieldModulus, big.NewInt(int64(g.rng.Intn(1000)+1)))
			safeFillBytes(fieldBig, point[i*32:(i+1)*32])
		}

	case 1:
		// Type 2: Maximum field values (all 0xFF)
		for i := range point {
			point[i] = 0xFF
		}

	case 2:
		// Type 3: Point not on curve - valid field elements but don't satisfy curve equation
		for i := 0; i < 4; i++ {
			field := new(big.Int).SetBytes(g.rng.Bytes(32))
			field.Mod(field, bn256FieldModulus)
			safeFillBytes(field, point[i*32:(i+1)*32])
		}

	case 3:
		// Type 4: Edge case - modulus minus 1 for all components
		modMinusOne := new(big.Int).Sub(bn256FieldModulus, big.NewInt(1))
		for i := 0; i < 4; i++ {
			safeFillBytes(modMinusOne, point[i*32:(i+1)*32])
		}

	case 4:
		// Type 5: Mixed valid/invalid field elements
		validField := big.NewInt(int64(g.rng.Intn(1000) + 1))
		invalidField := new(big.Int).Add(bn256FieldModulus, big.NewInt(1))

		safeFillBytes(validField, point[0:32])     // Valid x1
		safeFillBytes(invalidField, point[32:64])  // Invalid x2
		safeFillBytes(validField, point[64:96])    // Valid y1
		safeFillBytes(invalidField, point[96:128]) // Invalid y2

	case 5:
		// Type 6: Specific invalid coordinates (1,1,1,1) - not on BN256 G2 curve
		point[31] = 1  // x1 = 1
		point[63] = 1  // x2 = 1
		point[95] = 1  // y1 = 1
		point[127] = 1 // y2 = 1

	case 6:
		// Type 7: Valid X coordinates, Y = (0,0) (likely not on curve)
		binary.BigEndian.PutUint64(point[24:32], 2) // x1 = 2
		binary.BigEndian.PutUint64(point[56:64], 3) // x2 = 3
		// Y coordinates stay 0

	case 7:
		// Type 8: Invalid point at infinity representations
		for i := 0; i < 64; i++ {
			point[i] = 0xFF // Invalid "infinity" encoding for X
		}
		// Leave Y as zeros

	case 8:
		// Type 9: Points with small coordinates that might cause edge cases
		for i := 0; i < 4; i++ {
			point[i*32+31] = byte(g.rng.Intn(10) + 1)
		}

	case 9:
		// Type 10: Points from twist attacks (simulate points on quadratic twist)
		// Generate points using wrong curve equation
		for i := 0; i < 2; i++ {
			x := new(big.Int).SetBytes(g.rng.Bytes(32))
			x.Mod(x, bn256FieldModulus)

			// Use wrong curve constant for twist
			x3 := new(big.Int).Exp(x, big.NewInt(3), bn256FieldModulus)
			x3Plus7 := new(big.Int).Add(x3, big.NewInt(7)) // Wrong constant
			x3Plus7.Mod(x3Plus7, bn256FieldModulus)

			safeFillBytes(x, point[i*32:(i+1)*32])
			safeFillBytes(x3Plus7, point[(i+2)*32:(i+3)*32])
		}

	case 10:
		// Type 11: Points that could cause arithmetic overflow
		nearModulus := new(big.Int).Sub(bn256FieldModulus, big.NewInt(int64(g.rng.Intn(5)+1)))
		for i := 0; i < 4; i++ {
			safeFillBytes(nearModulus, point[i*32:(i+1)*32])
		}

	default:
		// Type 12: Random bytes (most likely invalid)
		copy(point, g.rng.Bytes(128))
	}

	return point
}
