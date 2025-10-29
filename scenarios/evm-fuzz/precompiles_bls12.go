package evmfuzz

import (
	"math/big"

	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fp"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
)

// BLS12-381 curve parameters and constants for generating invalid inputs
var (
	// BLS12-381 field modulus (q): 4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894272559787
	bls12381FieldModulus = new(big.Int)

	// BLS12-381 curve order (r): 52435875175126190479447740508185965837690552500527637822603658699938581184513
	bls12381CurveOrder = new(big.Int)

	// BLS12-381 cofactor for G1: 0x396c8c005555e1568c00aaab0000aaab
	bls12381G1Cofactor = new(big.Int)

	// BLS12-381 cofactor for G2: 0x5d543a95414e7f1091d50792876a202cd91de4547085abaa68a205b2e5a7ddfa628f1cb4d9e82ef21537e293a6691ae1616ec6e786f0c70cf1c38e31c7238e5
	bls12381G2Cofactor = new(big.Int)
)

func init() {
	// Initialize BLS12-381 field modulus (q)
	bls12381FieldModulus.SetString("4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894272559787", 10)

	// Initialize BLS12-381 curve order (r) - order of prime subgroup
	bls12381CurveOrder.SetString("52435875175126190479447740508185965837690552500527637822603658699938581184513", 10)

	// Initialize G1 cofactor
	bls12381G1Cofactor.SetString("396c8c005555e1568c00aaab0000aaab", 16)

	// Initialize G2 cofactor
	bls12381G2Cofactor.SetString("5d543a95414e7f1091d50792876a202cd91de4547085abaa68a205b2e5a7ddfa628f1cb4d9e82ef21537e293a6691ae1616ec6e786f0c70cf1c38e31c7238e5", 16)
}

// generateValidBLS12G1PointUncompressed generates a BLS12-381 G1 point using gnark-crypto library
func (g *OpcodeGenerator) generateValidBLS12G1PointUncompressed() []byte {
	point := make([]byte, 128) // 128 bytes uncompressed format: 64 bytes x + 64 bytes y

	choice := g.rng.Intn(100)

	if choice < 20 {
		// 20% chance to generate invalid G1 points for comprehensive testing
		return g.generateInvalidBLS12G1Point()
	}

	if choice < 25 {
		// 5% chance of zero point (point at infinity)
		// Zero point: all zeros represents point at infinity in uncompressed format
		return g.transformer.TransformPrecompileInput(point, 0x0b)
	}

	// 75% of the time: Generate random valid points using gnark-crypto library
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
	copy(point[16:64], xBytes[:])
	copy(point[80:128], yBytes[:])

	return g.transformer.TransformPrecompileInput(point, 0x0b)
}

// generateInvalidBLS12G1Point generates various types of invalid G1 points for comprehensive testing
func (g *OpcodeGenerator) generateInvalidBLS12G1Point() []byte {
	point := make([]byte, 128)
	invalidType := g.rng.Intn(12) // Increased to 12 for more subgroup attack variants

	switch invalidType {
	case 0:
		// Type 1: Field elements greater than modulus
		// Generate x coordinate > field modulus
		xBig := new(big.Int).Add(bls12381FieldModulus, big.NewInt(int64(g.rng.Intn(1000)+1)))
		xBytes := xBig.FillBytes(make([]byte, 48))
		copy(point[16:64], xBytes)

		// Generate y coordinate > field modulus
		yBig := new(big.Int).Add(bls12381FieldModulus, big.NewInt(int64(g.rng.Intn(1000)+1)))
		yBytes := yBig.FillBytes(make([]byte, 48))
		copy(point[80:128], yBytes)

	case 1:
		// Type 2: Maximum field values (all 0xFF in field element positions)
		for i := 16; i < 64; i++ {
			point[i] = 0xFF
		}
		for i := 80; i < 128; i++ {
			point[i] = 0xFF
		}

	case 2:
		// Type 3: Point not on curve - valid field elements but don't satisfy curve equation
		// Use valid x coordinate, calculate wrong y
		var validX fp.Element
		validX.SetBytes(g.rng.Bytes(48))
		xBytes := validX.Bytes()
		copy(point[16:64], xBytes[:])

		// Set y to a random valid field element (likely not satisfying curve equation)
		var invalidY fp.Element
		invalidY.SetBytes(g.rng.Bytes(48))
		yBytes := invalidY.Bytes()
		copy(point[80:128], yBytes[:])

	case 3:
		// Type 4: Point on curve but not in prime subgroup (random multiple of cofactor)
		var g1Point bls12381.G1Affine
		_, _, g1Gen, _ := bls12381.Generators()

		// Generate random multiple of cofactor (1 to 100) to get various points not in prime subgroup
		multiplier := big.NewInt(int64(g.rng.Intn(100) + 1))
		cofactorMultiple := new(big.Int).Mul(bls12381G1Cofactor, multiplier)

		// Multiply by (cofactor * random) to get point on curve but not in prime subgroup
		g1Point.ScalarMultiplication(&g1Gen, cofactorMultiple)

		xBytes := g1Point.X.Bytes()
		yBytes := g1Point.Y.Bytes()
		copy(point[16:64], xBytes[:])
		copy(point[80:128], yBytes[:])

	case 4:
		// Type 5: Invalid encoding - non-zero bytes in padding areas
		// Generate valid point first
		var g1Point bls12381.G1Affine
		var scalar fr.Element
		scalar.SetBytes(g.rng.Bytes(32))
		_, _, g1Gen, _ := bls12381.Generators()
		g1Point.ScalarMultiplication(&g1Gen, scalar.BigInt(new(big.Int)))

		xBytes := g1Point.X.Bytes()
		yBytes := g1Point.Y.Bytes()
		copy(point[16:64], xBytes[:])
		copy(point[80:128], yBytes[:])

		// Corrupt padding areas (should be zero)
		for i := 0; i < 16; i++ {
			point[i] = byte(g.rng.Intn(256))
		}
		for i := 64; i < 80; i++ {
			point[i] = byte(g.rng.Intn(256))
		}

	case 5:
		// Type 6: Edge case - modulus minus 1
		modMinusOne := new(big.Int).Sub(bls12381FieldModulus, big.NewInt(1))
		xBytes := modMinusOne.FillBytes(make([]byte, 48))
		yBytes := modMinusOne.FillBytes(make([]byte, 48))
		copy(point[16:64], xBytes)
		copy(point[80:128], yBytes)

	case 6:
		// Type 7: Specific invalid coordinates that look valid but aren't
		// x = 1, y = 1 (not on BLS12-381 curve)
		point[63] = 1  // x = 1
		point[127] = 1 // y = 1

	case 7:
		// Type 8: Valid x, y = 0 (unless x = 0, this won't be on curve)
		var validX fp.Element
		// Use a non-zero x to ensure (x,0) is not on curve
		validX.SetUint64(2)
		xBytes := validX.Bytes()
		copy(point[16:64], xBytes[:])
		// y stays 0

	case 8:
		// Type 9: Compressed point flags in uncompressed format (invalid)
		var g1Point bls12381.G1Affine
		var scalar fr.Element
		scalar.SetBytes(g.rng.Bytes(32))
		_, _, g1Gen, _ := bls12381.Generators()
		g1Point.ScalarMultiplication(&g1Gen, scalar.BigInt(new(big.Int)))

		xBytes := g1Point.X.Bytes()
		yBytes := g1Point.Y.Bytes()
		copy(point[16:64], xBytes[:])
		copy(point[80:128], yBytes[:])

		// Set compression flags (invalid for uncompressed format)
		point[0] |= 0x80 // Set compression bit

	case 10:
		// Type 11: Point from random scalar multiplication (may not be in subgroup)
		// Generate point using random scalar that's not reduced modulo curve order
		var g1Point bls12381.G1Affine
		_, _, g1Gen, _ := bls12381.Generators()

		// Use a scalar that's specifically designed to potentially create subgroup issues
		// Multiply by (curve_order + cofactor * random)
		randomOffset := big.NewInt(int64(g.rng.Intn(1000) + 1))
		cofactorOffset := new(big.Int).Mul(bls12381G1Cofactor, randomOffset)
		invalidScalar := new(big.Int).Add(bls12381CurveOrder, cofactorOffset)

		g1Point.ScalarMultiplication(&g1Gen, invalidScalar)

		xBytes := g1Point.X.Bytes()
		yBytes := g1Point.Y.Bytes()
		copy(point[16:64], xBytes[:])
		copy(point[80:128], yBytes[:])

	case 11:
		// Type 12: Point from cofactor clearing failure simulation
		// Start with a valid point, then multiply by a value that simulates cofactor clearing failure
		var g1Point bls12381.G1Affine
		_, _, g1Gen, _ := bls12381.Generators()

		// Generate a scalar that when multiplied by cofactor gives interesting edge cases
		baseScalar := big.NewInt(int64(g.rng.Intn(1000) + 1))

		// Create a scalar that's cofactor * base + small_offset (simulates partial cofactor clearing)
		smallOffset := big.NewInt(int64(g.rng.Intn(10) + 1))
		cofactorBase := new(big.Int).Mul(bls12381G1Cofactor, baseScalar)
		partialScalar := new(big.Int).Add(cofactorBase, smallOffset)

		g1Point.ScalarMultiplication(&g1Gen, partialScalar)

		xBytes := g1Point.X.Bytes()
		yBytes := g1Point.Y.Bytes()
		copy(point[16:64], xBytes[:])
		copy(point[80:128], yBytes[:])

	default:
		// Type 10: Random bytes (most likely invalid)
		copy(point, g.rng.Bytes(128))
	}

	return g.transformer.TransformPrecompileInput(point, 0x0b)
}

// generateValidBLS12G2PointUncompressed generates a BLS12-381 G2 point using gnark-crypto library
func (g *OpcodeGenerator) generateValidBLS12G2PointUncompressed() []byte {
	point := make([]byte, 256) // 256 bytes uncompressed format

	choice := g.rng.Intn(100)

	if choice < 20 {
		// 20% chance to generate invalid G2 points for comprehensive testing
		return g.generateInvalidBLS12G2Point()
	}

	if choice < 25 {
		// 5% chance of zero point (point at infinity)
		// Zero point: all zeros represents point at infinity in uncompressed format
		return g.transformer.TransformPrecompileInput(point, 0x0d)
	}

	// 75% of the time: Generate random valid points using gnark-crypto library
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

// generateInvalidBLS12G2Point generates various types of invalid G2 points for comprehensive testing
func (g *OpcodeGenerator) generateInvalidBLS12G2Point() []byte {
	point := make([]byte, 256)
	invalidType := g.rng.Intn(15) // Increased to 15 for more subgroup attack variants

	switch invalidType {
	case 0:
		// Type 1: Field elements greater than modulus
		for i := 0; i < 4; i++ {
			fieldBig := new(big.Int).Add(bls12381FieldModulus, big.NewInt(int64(g.rng.Intn(1000)+1)))
			fieldBytes := fieldBig.FillBytes(make([]byte, 48))
			offset := i*64 + 16
			copy(point[offset:offset+48], fieldBytes)
		}

	case 1:
		// Type 2: Maximum field values (all 0xFF in field element positions)
		for i := 0; i < 4; i++ {
			offset := i*64 + 16
			for j := offset; j < offset+48; j++ {
				point[j] = 0xFF
			}
		}

	case 2:
		// Type 3: Point not on curve - valid field elements but don't satisfy curve equation
		var xc0, xc1, yc0, yc1 fp.Element
		xc0.SetBytes(g.rng.Bytes(48))
		xc1.SetBytes(g.rng.Bytes(48))
		yc0.SetBytes(g.rng.Bytes(48))
		yc1.SetBytes(g.rng.Bytes(48))

		xc0Bytes := xc0.Bytes()
		xc1Bytes := xc1.Bytes()
		yc0Bytes := yc0.Bytes()
		yc1Bytes := yc1.Bytes()
		copy(point[16:64], xc0Bytes[:])
		copy(point[80:128], xc1Bytes[:])
		copy(point[144:192], yc0Bytes[:])
		copy(point[208:256], yc1Bytes[:])

	case 3:
		// Type 4: Point on curve but not in prime subgroup (random multiple of cofactor)
		var g2Point bls12381.G2Affine
		_, _, _, g2Gen := bls12381.Generators()

		// Generate random multiple of cofactor (1 to 50) to get various points not in prime subgroup
		// Using smaller range for G2 since cofactor is much larger
		multiplier := big.NewInt(int64(g.rng.Intn(50) + 1))
		cofactorMultiple := new(big.Int).Mul(bls12381G2Cofactor, multiplier)

		// Multiply by (cofactor * random) to get point on curve but not in prime subgroup
		g2Point.ScalarMultiplication(&g2Gen, cofactorMultiple)

		xc0Bytes := g2Point.X.A0.Bytes()
		xc1Bytes := g2Point.X.A1.Bytes()
		yc0Bytes := g2Point.Y.A0.Bytes()
		yc1Bytes := g2Point.Y.A1.Bytes()
		copy(point[16:64], xc0Bytes[:])
		copy(point[80:128], xc1Bytes[:])
		copy(point[144:192], yc0Bytes[:])
		copy(point[208:256], yc1Bytes[:])

	case 4:
		// Type 5: Invalid encoding - non-zero bytes in padding areas
		var g2Point bls12381.G2Affine
		var scalar fr.Element
		scalar.SetBytes(g.rng.Bytes(32))
		_, _, _, g2Gen := bls12381.Generators()
		g2Point.ScalarMultiplication(&g2Gen, scalar.BigInt(new(big.Int)))

		xc0Bytes := g2Point.X.A0.Bytes()
		xc1Bytes := g2Point.X.A1.Bytes()
		yc0Bytes := g2Point.Y.A0.Bytes()
		yc1Bytes := g2Point.Y.A1.Bytes()
		copy(point[16:64], xc0Bytes[:])
		copy(point[80:128], xc1Bytes[:])
		copy(point[144:192], yc0Bytes[:])
		copy(point[208:256], yc1Bytes[:])

		// Corrupt padding areas (should be zero)
		for i := 0; i < 16; i++ {
			point[i] = byte(g.rng.Intn(256))
		}
		for i := 64; i < 80; i++ {
			point[i] = byte(g.rng.Intn(256))
		}
		for i := 128; i < 144; i++ {
			point[i] = byte(g.rng.Intn(256))
		}
		for i := 192; i < 208; i++ {
			point[i] = byte(g.rng.Intn(256))
		}

	case 5:
		// Type 6: Edge case - modulus minus 1 for all components
		modMinusOne := new(big.Int).Sub(bls12381FieldModulus, big.NewInt(1))
		fieldBytes := modMinusOne.FillBytes(make([]byte, 48))
		copy(point[16:64], fieldBytes)
		copy(point[80:128], fieldBytes)
		copy(point[144:192], fieldBytes)
		copy(point[208:256], fieldBytes)

	case 6:
		// Type 7: Mixed valid/invalid field elements
		var validField fp.Element
		validField.SetBytes(g.rng.Bytes(48))
		validBytes := validField.Bytes()

		invalidBig := new(big.Int).Add(bls12381FieldModulus, big.NewInt(1))
		invalidBytes := invalidBig.FillBytes(make([]byte, 48))

		copy(point[16:64], validBytes[:])   // Valid X_c0
		copy(point[80:128], invalidBytes)   // Invalid X_c1
		copy(point[144:192], validBytes[:]) // Valid Y_c0
		copy(point[208:256], invalidBytes)  // Invalid Y_c1

	case 7:
		// Type 8: Specific invalid coordinates (1,1,1,1) - not on BLS12-381 G2 curve
		point[63] = 1  // X_c0 = 1
		point[127] = 1 // X_c1 = 1
		point[191] = 1 // Y_c0 = 1
		point[255] = 1 // Y_c1 = 1

	case 8:
		// Type 9: Valid X coordinates, Y = (0,0) (unless special case, won't be on curve)
		var xc0, xc1 fp.Element
		xc0.SetUint64(2)
		xc1.SetUint64(3)
		xc0Bytes := xc0.Bytes()
		xc1Bytes := xc1.Bytes()
		copy(point[16:64], xc0Bytes[:])
		copy(point[80:128], xc1Bytes[:])
		// Y coordinates stay 0

	case 9:
		// Type 10: Compressed point flags in uncompressed format (invalid)
		var g2Point bls12381.G2Affine
		var scalar fr.Element
		scalar.SetBytes(g.rng.Bytes(32))
		_, _, _, g2Gen := bls12381.Generators()
		g2Point.ScalarMultiplication(&g2Gen, scalar.BigInt(new(big.Int)))

		xc0Bytes := g2Point.X.A0.Bytes()
		xc1Bytes := g2Point.X.A1.Bytes()
		yc0Bytes := g2Point.Y.A0.Bytes()
		yc1Bytes := g2Point.Y.A1.Bytes()
		copy(point[16:64], xc0Bytes[:])
		copy(point[80:128], xc1Bytes[:])
		copy(point[144:192], yc0Bytes[:])
		copy(point[208:256], yc1Bytes[:])

		// Set compression flags (invalid for uncompressed format)
		point[0] |= 0x80 // Set compression bit

	case 10:
		// Type 11: Infinity flag with non-zero coordinates (invalid)
		var g2Point bls12381.G2Affine
		var scalar fr.Element
		scalar.SetBytes(g.rng.Bytes(32))
		_, _, _, g2Gen := bls12381.Generators()
		g2Point.ScalarMultiplication(&g2Gen, scalar.BigInt(new(big.Int)))

		xc0Bytes := g2Point.X.A0.Bytes()
		xc1Bytes := g2Point.X.A1.Bytes()
		yc0Bytes := g2Point.Y.A0.Bytes()
		yc1Bytes := g2Point.Y.A1.Bytes()
		copy(point[16:64], xc0Bytes[:])
		copy(point[80:128], xc1Bytes[:])
		copy(point[144:192], yc0Bytes[:])
		copy(point[208:256], yc1Bytes[:])

		// Set infinity flag with non-zero coordinates (invalid)
		point[0] |= 0x40 // Set infinity bit

	case 12:
		// Type 13: Point from random scalar multiplication (may not be in subgroup)
		var g2Point bls12381.G2Affine
		_, _, _, g2Gen := bls12381.Generators()

		// Use a scalar that's specifically designed to potentially create subgroup issues
		// Multiply by (curve_order + cofactor * random)
		randomOffset := big.NewInt(int64(g.rng.Intn(100) + 1)) // Smaller for G2
		cofactorOffset := new(big.Int).Mul(bls12381G2Cofactor, randomOffset)
		invalidScalar := new(big.Int).Add(bls12381CurveOrder, cofactorOffset)

		g2Point.ScalarMultiplication(&g2Gen, invalidScalar)

		xc0Bytes := g2Point.X.A0.Bytes()
		xc1Bytes := g2Point.X.A1.Bytes()
		yc0Bytes := g2Point.Y.A0.Bytes()
		yc1Bytes := g2Point.Y.A1.Bytes()
		copy(point[16:64], xc0Bytes[:])
		copy(point[80:128], xc1Bytes[:])
		copy(point[144:192], yc0Bytes[:])
		copy(point[208:256], yc1Bytes[:])

	case 13:
		// Type 14: Point from cofactor clearing failure simulation
		var g2Point bls12381.G2Affine
		_, _, _, g2Gen := bls12381.Generators()

		// Generate a scalar that when multiplied by cofactor gives interesting edge cases
		baseScalar := big.NewInt(int64(g.rng.Intn(100) + 1))

		// Create a scalar that's cofactor * base + small_offset (simulates partial cofactor clearing)
		smallOffset := big.NewInt(int64(g.rng.Intn(5) + 1))
		cofactorBase := new(big.Int).Mul(bls12381G2Cofactor, baseScalar)
		partialScalar := new(big.Int).Add(cofactorBase, smallOffset)

		g2Point.ScalarMultiplication(&g2Gen, partialScalar)

		xc0Bytes := g2Point.X.A0.Bytes()
		xc1Bytes := g2Point.X.A1.Bytes()
		yc0Bytes := g2Point.Y.A0.Bytes()
		yc1Bytes := g2Point.Y.A1.Bytes()
		copy(point[16:64], xc0Bytes[:])
		copy(point[80:128], xc1Bytes[:])
		copy(point[144:192], yc0Bytes[:])
		copy(point[208:256], yc1Bytes[:])

	case 14:
		// Type 15: Point from twist attack simulation
		// Generate point that might be on the quadratic/sextic twist
		var g2Point bls12381.G2Affine
		_, _, _, g2Gen := bls12381.Generators()

		// Use a scalar that's a multiple of a small prime factor of the twist order
		// This simulates points that could come from twist attacks
		twistFactor := big.NewInt(int64(g.rng.Intn(997) + 2)) // Random small prime-like number
		cofactorTwist := new(big.Int).Mul(bls12381G2Cofactor, twistFactor)

		g2Point.ScalarMultiplication(&g2Gen, cofactorTwist)

		xc0Bytes := g2Point.X.A0.Bytes()
		xc1Bytes := g2Point.X.A1.Bytes()
		yc0Bytes := g2Point.Y.A0.Bytes()
		yc1Bytes := g2Point.Y.A1.Bytes()
		copy(point[16:64], xc0Bytes[:])
		copy(point[80:128], xc1Bytes[:])
		copy(point[144:192], yc0Bytes[:])
		copy(point[208:256], yc1Bytes[:])

	default:
		// Type 12: Random bytes (most likely invalid)
		copy(point, g.rng.Bytes(256))
	}

	return g.transformer.TransformPrecompileInput(point, 0x0d)
}

// generateInvalidScalar generates invalid scalars for MSM operations
func (g *OpcodeGenerator) generateInvalidScalar() []byte {
	scalar := make([]byte, 32)
	invalidType := g.rng.Intn(8)

	switch invalidType {
	case 0:
		// Type 1: Scalar >= curve order (invalid)
		scalarBig := new(big.Int).Add(bls12381CurveOrder, big.NewInt(int64(g.rng.Intn(1000)+1)))
		safeFillBytes(scalarBig, scalar)

	case 1:
		// Type 2: Scalar = curve order (invalid)
		safeFillBytes(bls12381CurveOrder, scalar)

	case 2:
		// Type 3: Maximum value (all 0xFF)
		for i := range scalar {
			scalar[i] = 0xFF
		}

	case 3:
		// Type 4: Curve order - 1 (valid but edge case)
		orderMinusOne := new(big.Int).Sub(bls12381CurveOrder, big.NewInt(1))
		safeFillBytes(orderMinusOne, scalar)

	case 4:
		// Type 5: Powers of 2 near curve order
		powerOf2 := new(big.Int).Lsh(big.NewInt(1), 255) // 2^255
		safeFillBytes(powerOf2, scalar)

	case 5:
		// Type 6: Negative values (two's complement representation)
		// Generate valid scalar first, then negate
		var validScalar fr.Element
		validScalar.SetBytes(g.rng.Bytes(32))
		validBig := validScalar.BigInt(new(big.Int))
		negated := new(big.Int).Sub(bls12381CurveOrder, validBig)
		safeFillBytes(negated, scalar)

	case 6:
		// Type 7: Small multiples of curve order
		multiple := big.NewInt(int64(g.rng.Intn(10) + 2)) // 2-11
		invalidScalar := new(big.Int).Mul(bls12381CurveOrder, multiple)
		safeFillBytes(invalidScalar, scalar)

	default:
		// Type 8: Random bytes (likely invalid)
		copy(scalar, g.rng.Bytes(32))
		// Ensure it's >= curve order by setting high bits
		scalar[0] |= 0xF0
	}

	return scalar
}

// generateInvalidFieldElement generates invalid field elements for mapping operations
func (g *OpcodeGenerator) generateInvalidFieldElement() []byte {
	fieldElement := make([]byte, 64)
	invalidType := g.rng.Intn(8)

	switch invalidType {
	case 0:
		// Type 1: Field element >= modulus
		fieldBig := new(big.Int).Add(bls12381FieldModulus, big.NewInt(int64(g.rng.Intn(1000)+1)))
		fieldBytes := fieldBig.FillBytes(make([]byte, 48))
		copy(fieldElement[16:64], fieldBytes) // Place at end of 64-byte slot

	case 1:
		// Type 2: Field element = modulus
		fieldBytes := bls12381FieldModulus.FillBytes(make([]byte, 48))
		copy(fieldElement[16:64], fieldBytes)

	case 2:
		// Type 3: Maximum value (all 0xFF)
		for i := 16; i < 64; i++ {
			fieldElement[i] = 0xFF
		}

	case 3:
		// Type 4: Invalid padding (non-zero bytes in padding area)
		var validField fp.Element
		validField.SetBytes(g.rng.Bytes(48))
		fieldBytes := validField.Bytes()
		copy(fieldElement[16:64], fieldBytes[:])

		// Corrupt padding area (should be zero)
		for i := 0; i < 16; i++ {
			fieldElement[i] = byte(g.rng.Intn(256))
		}

	case 4:
		// Type 5: Modulus - 1 (valid but edge case)
		modMinusOne := new(big.Int).Sub(bls12381FieldModulus, big.NewInt(1))
		fieldBytes := modMinusOne.FillBytes(make([]byte, 48))
		copy(fieldElement[16:64], fieldBytes)

	case 5:
		// Type 6: Powers of 2 near modulus
		powerOf2 := new(big.Int).Lsh(big.NewInt(1), 381) // 2^381 (close to field size)
		fieldBytes := powerOf2.FillBytes(make([]byte, 48))
		copy(fieldElement[16:64], fieldBytes)

	case 6:
		// Type 7: Small values with high bits set (definitely > modulus)
		copy(fieldElement[16:64], g.rng.Bytes(48))
		fieldElement[16] |= 0xF0 // Ensure high bits are set

	default:
		// Type 8: Random bytes (likely invalid)
		copy(fieldElement, g.rng.Bytes(64))
	}

	return fieldElement
}

// generateInvalidFieldExtensionElement generates invalid Fp2 elements for G2 mapping
func (g *OpcodeGenerator) generateInvalidFieldExtensionElement() []byte {
	fieldElement := make([]byte, 128)
	invalidType := g.rng.Intn(10)

	switch invalidType {
	case 0:
		// Type 1: Both components >= modulus
		for i := 0; i < 2; i++ {
			fieldBig := new(big.Int).Add(bls12381FieldModulus, big.NewInt(int64(g.rng.Intn(1000)+1)))
			fieldBytes := fieldBig.FillBytes(make([]byte, 48))
			offset := i*64 + 16
			copy(fieldElement[offset:offset+48], fieldBytes)
		}

	case 1:
		// Type 2: One component valid, one invalid
		var validField fp.Element
		validField.SetBytes(g.rng.Bytes(48))
		validBytes := validField.Bytes()
		copy(fieldElement[16:64], validBytes[:]) // Valid c0

		invalidBig := new(big.Int).Add(bls12381FieldModulus, big.NewInt(1))
		invalidBytes := invalidBig.FillBytes(make([]byte, 48))
		copy(fieldElement[80:128], invalidBytes) // Invalid c1

	case 2:
		// Type 3: Both components = modulus
		fieldBytes := bls12381FieldModulus.FillBytes(make([]byte, 48))
		copy(fieldElement[16:64], fieldBytes)
		copy(fieldElement[80:128], fieldBytes)

	case 3:
		// Type 4: Maximum values (all 0xFF)
		for i := 16; i < 64; i++ {
			fieldElement[i] = 0xFF
		}
		for i := 80; i < 128; i++ {
			fieldElement[i] = 0xFF
		}

	case 4:
		// Type 5: Invalid padding in both components
		var c0, c1 fp.Element
		c0.SetBytes(g.rng.Bytes(48))
		c1.SetBytes(g.rng.Bytes(48))
		c0Bytes := c0.Bytes()
		c1Bytes := c1.Bytes()
		copy(fieldElement[16:64], c0Bytes[:])
		copy(fieldElement[80:128], c1Bytes[:])

		// Corrupt padding areas
		for i := 0; i < 16; i++ {
			fieldElement[i] = byte(g.rng.Intn(256))
		}
		for i := 64; i < 80; i++ {
			fieldElement[i] = byte(g.rng.Intn(256))
		}

	case 5:
		// Type 6: Edge case - modulus - 1 for both
		modMinusOne := new(big.Int).Sub(bls12381FieldModulus, big.NewInt(1))
		fieldBytes := modMinusOne.FillBytes(make([]byte, 48))
		copy(fieldElement[16:64], fieldBytes)
		copy(fieldElement[80:128], fieldBytes)

	case 6:
		// Type 7: Zero c0, invalid c1
		// c0 stays zero
		invalidBig := new(big.Int).Add(bls12381FieldModulus, big.NewInt(1))
		invalidBytes := invalidBig.FillBytes(make([]byte, 48))
		copy(fieldElement[80:128], invalidBytes)

	case 7:
		// Type 8: Invalid c0, zero c1
		invalidBig := new(big.Int).Add(bls12381FieldModulus, big.NewInt(1))
		invalidBytes := invalidBig.FillBytes(make([]byte, 48))
		copy(fieldElement[16:64], invalidBytes)
		// c1 stays zero

	case 8:
		// Type 9: Powers of 2 for both components
		powerOf2 := new(big.Int).Lsh(big.NewInt(1), 381)
		fieldBytes := powerOf2.FillBytes(make([]byte, 48))
		copy(fieldElement[16:64], fieldBytes)
		copy(fieldElement[80:128], fieldBytes)

	default:
		// Type 10: Random bytes (likely invalid)
		copy(fieldElement, g.rng.Bytes(128))
	}

	return fieldElement
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

		// 15% chance of invalid scalar for testing
		var scalar []byte
		if g.rng.Intn(100) < 15 {
			scalar = g.generateInvalidScalar()
		} else {
			scalar = g.rng.Bytes(32)
		}

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

		// 15% chance of invalid scalar for testing
		var scalar []byte
		if g.rng.Intn(100) < 15 {
			scalar = g.generateInvalidScalar()
		} else {
			scalar = g.rng.Bytes(32)
		}

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
	var fieldElement []byte

	choice := g.rng.Intn(100)
	if choice < 20 {
		// 20% chance of invalid field elements for comprehensive testing
		fieldElement = g.generateInvalidFieldElement()
	} else {
		// 80% chance of valid field elements
		fieldElement = make([]byte, 64)
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
	var fieldElement []byte

	choice := g.rng.Intn(100)
	if choice < 20 {
		// 20% chance of invalid field extension elements for comprehensive testing
		fieldElement = g.generateInvalidFieldExtensionElement()
	} else {
		// 80% chance of valid field extension elements
		fieldElement = make([]byte, 128)
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
