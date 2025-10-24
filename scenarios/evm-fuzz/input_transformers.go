package evmfuzz

import (
	"crypto/rand"
	"encoding/binary"
)

// InputTransformer applies various transformations to inputs to generate edge cases
type InputTransformer struct {
	rng *DeterministicRNG
}

// NewInputTransformer creates a new input transformer
func NewInputTransformer(rng *DeterministicRNG) *InputTransformer {
	return &InputTransformer{rng: rng}
}

// TransformInput applies random transformations to input data based on probability
func (t *InputTransformer) TransformInput(input []byte, expectedSize int) []byte {
	// 15% chance to apply transformations (keeping most inputs valid)
	if t.rng.Float64() > 0.15 {
		return input
	}

	transformType := t.rng.Intn(100)

	switch {
	case transformType < 20: // 20% of transformations: Size variations
		return t.applySizeTransformations(input, expectedSize)
	case transformType < 40: // 20% of transformations: Data corruption
		return t.applyDataCorruption(input)
	case transformType < 60: // 20% of transformations: Boundary conditions
		return t.applyBoundaryConditions(input, expectedSize)
	case transformType < 80: // 20% of transformations: Format violations
		return t.applyFormatViolations(input, expectedSize)
	default: // 20% of transformations: Extreme edge cases
		return t.applyExtremeEdgeCases(input, expectedSize)
	}
}

// applySizeTransformations generates inputs with wrong sizes
func (t *InputTransformer) applySizeTransformations(input []byte, expectedSize int) []byte {
	sizeType := t.rng.Intn(100)

	switch {
	case sizeType < 20: // Too short by 1-16 bytes
		reduction := t.rng.Intn(16) + 1
		if len(input) <= reduction {
			return []byte{} // Empty input
		}
		return input[:len(input)-reduction]

	case sizeType < 40: // Too long by 1-64 bytes
		extension := t.rng.Intn(64) + 1
		extended := make([]byte, len(input)+extension)
		copy(extended, input)
		// Fill extension with random or pattern data
		if t.rng.Intn(2) == 0 {
			copy(extended[len(input):], t.rng.Bytes(extension))
		} else {
			// Pattern fill (useful for detecting buffer overruns)
			for i := len(input); i < len(extended); i++ {
				extended[i] = 0xAA
			}
		}
		return extended

	case sizeType < 60: // Completely wrong size
		wrongSizes := []int{0, 1, 3, 7, 15, 31, 33, 63, 65, 127, 129, 255, 257, 511, 513, 1023, 1025}
		wrongSize := wrongSizes[t.rng.Intn(len(wrongSizes))]
		result := make([]byte, wrongSize)
		if wrongSize > 0 {
			copy(result, t.rng.Bytes(wrongSize))
		}
		return result

	case sizeType < 80: // Power-of-2 boundary violations
		boundaries := []int{2, 4, 8, 16, 32, 64, 128, 256, 512, 1024}
		boundary := boundaries[t.rng.Intn(len(boundaries))]
		// Off-by-one from power of 2
		if t.rng.Intn(2) == 0 {
			boundary-- // One less than power of 2
		} else {
			boundary++ // One more than power of 2
		}
		result := make([]byte, boundary)
		copy(result, t.rng.Bytes(boundary))
		return result

	default: // Extremely large inputs (test memory limits)
		largeSizes := []int{2048, 4096, 8192, 16384, 32768, 65536}
		largeSize := largeSizes[t.rng.Intn(len(largeSizes))]
		result := make([]byte, largeSize)
		// Fill with patterns to detect issues
		pattern := byte(t.rng.Intn(256))
		for i := range result {
			result[i] = pattern
		}
		return result
	}
}

// applyDataCorruption corrupts data in various ways
func (t *InputTransformer) applyDataCorruption(input []byte) []byte {
	if len(input) == 0 {
		return input
	}

	result := make([]byte, len(input))
	copy(result, input)

	corruptionType := t.rng.Intn(100)

	switch {
	case corruptionType < 15: // Bit flips
		numFlips := t.rng.Intn(8) + 1 // 1-8 bit flips
		for i := 0; i < numFlips; i++ {
			byteIdx := t.rng.Intn(len(result))
			bitIdx := t.rng.Intn(8)
			result[byteIdx] ^= (1 << bitIdx)
		}

	case corruptionType < 30: // Byte corruption
		numBytes := t.rng.Intn(len(result)/4 + 1) // Up to 25% of bytes
		for i := 0; i < numBytes; i++ {
			idx := t.rng.Intn(len(result))
			result[idx] = byte(t.rng.Intn(256))
		}

	case corruptionType < 45: // Zero out sections
		start := t.rng.Intn(len(result))
		length := t.rng.Intn(len(result)-start) + 1
		for i := start; i < start+length; i++ {
			result[i] = 0x00
		}

	case corruptionType < 60: // Fill sections with 0xFF
		start := t.rng.Intn(len(result))
		length := t.rng.Intn(len(result)-start) + 1
		for i := start; i < start+length; i++ {
			result[i] = 0xFF
		}

	case corruptionType < 75: // Swap bytes
		numSwaps := t.rng.Intn(len(result)/2) + 1
		for i := 0; i < numSwaps; i++ {
			idx1 := t.rng.Intn(len(result))
			idx2 := t.rng.Intn(len(result))
			result[idx1], result[idx2] = result[idx2], result[idx1]
		}

	case corruptionType < 90: // Insert/delete bytes (shift corruption)
		if t.rng.Intn(2) == 0 {
			// Insert random byte
			idx := t.rng.Intn(len(result))
			newResult := make([]byte, len(result)+1)
			copy(newResult[:idx], result[:idx])
			newResult[idx] = byte(t.rng.Intn(256))
			copy(newResult[idx+1:], result[idx:])
			result = newResult
		} else if len(result) > 1 {
			// Delete byte
			idx := t.rng.Intn(len(result))
			newResult := make([]byte, len(result)-1)
			copy(newResult[:idx], result[:idx])
			copy(newResult[idx:], result[idx+1:])
			result = newResult
		}

	default: // Pattern corruption (create recognizable patterns)
		patterns := []byte{0xAA, 0x55, 0xCC, 0x33, 0xF0, 0x0F}
		pattern := patterns[t.rng.Intn(len(patterns))]
		start := t.rng.Intn(len(result))
		length := t.rng.Intn(len(result)-start) + 1
		for i := start; i < start+length; i++ {
			result[i] = pattern
		}
	}

	return result
}

// applyBoundaryConditions generates inputs that test boundary conditions
func (t *InputTransformer) applyBoundaryConditions(input []byte, expectedSize int) []byte {
	boundaryType := t.rng.Intn(100)

	switch {
	case boundaryType < 25: // All zeros
		result := make([]byte, expectedSize)
		return result

	case boundaryType < 50: // All ones (0xFF)
		result := make([]byte, expectedSize)
		for i := range result {
			result[i] = 0xFF
		}
		return result

	case boundaryType < 70: // Alternating patterns
		result := make([]byte, expectedSize)
		patterns := [][]byte{
			{0x00, 0xFF}, {0xAA, 0x55}, {0xCC, 0x33}, {0xF0, 0x0F},
		}
		pattern := patterns[t.rng.Intn(len(patterns))]
		for i := range result {
			result[i] = pattern[i%len(pattern)]
		}
		return result

	case boundaryType < 85: // Single bit set
		result := make([]byte, expectedSize)
		if expectedSize > 0 {
			byteIdx := t.rng.Intn(expectedSize)
			bitIdx := t.rng.Intn(8)
			result[byteIdx] = 1 << bitIdx
		}
		return result

	default: // Incrementing/decrementing patterns
		result := make([]byte, expectedSize)
		if t.rng.Intn(2) == 0 {
			// Incrementing
			for i := range result {
				result[i] = byte(i % 256)
			}
		} else {
			// Decrementing
			for i := range result {
				result[i] = byte(255 - (i % 256))
			}
		}
		return result
	}
}

// applyFormatViolations creates inputs that violate expected formats
func (t *InputTransformer) applyFormatViolations(input []byte, expectedSize int) []byte {
	violationType := t.rng.Intn(100)

	switch {
	case violationType < 20: // Wrong endianness (swap byte order in chunks)
		result := make([]byte, len(input))
		copy(result, input)
		chunkSize := 4 // Swap 4-byte chunks
		if t.rng.Intn(2) == 0 {
			chunkSize = 8 // Or 8-byte chunks
		}

		for i := 0; i < len(result)-chunkSize+1; i += chunkSize {
			for j := 0; j < chunkSize/2; j++ {
				result[i+j], result[i+chunkSize-1-j] = result[i+chunkSize-1-j], result[i+j]
			}
		}
		return result

	case violationType < 40: // Misaligned data (shift by 1-7 bytes)
		if len(input) < 8 {
			return input
		}
		shift := t.rng.Intn(7) + 1
		result := make([]byte, len(input))
		copy(result[shift:], input[:len(input)-shift])
		// Fill shifted area with random data
		copy(result[:shift], t.rng.Bytes(shift))
		return result

	case violationType < 60: // Invalid field modulus (for curve points)
		result := make([]byte, expectedSize)
		copy(result, input)
		if expectedSize >= 32 {
			// Make first 32 bytes > field modulus by setting high bits
			result[0] |= 0xF0
			for i := 1; i < 4; i++ {
				result[i] = 0xFF
			}
		}
		return result

	case violationType < 80: // Invalid coordinate combinations
		result := make([]byte, expectedSize)
		copy(result, input)
		// For point inputs, create impossible coordinate combinations
		if expectedSize >= 64 { // At least x,y coordinates
			// Set x = 0, y = 1 (not a valid curve point for most curves)
			for i := 0; i < 32; i++ {
				result[i] = 0x00 // x = 0
			}
			for i := 32; i < 64; i++ {
				result[i] = 0x00 // y = 0 except last byte
			}
			result[63] = 0x01 // y = 1
		}
		return result

	default: // Invalid encoding flags/headers
		result := make([]byte, expectedSize)
		copy(result, input)
		if expectedSize > 0 {
			// Set invalid encoding flags in first byte
			invalidFlags := []byte{0x80, 0xC0, 0xE0, 0xF0, 0xF8, 0xFC, 0xFE, 0xFF}
			result[0] = invalidFlags[t.rng.Intn(len(invalidFlags))]
		}
		return result
	}
}

// applyExtremeEdgeCases generates extreme edge cases
func (t *InputTransformer) applyExtremeEdgeCases(input []byte, expectedSize int) []byte {
	edgeType := t.rng.Intn(100)

	switch {
	case edgeType < 20: // Completely random data
		result := make([]byte, expectedSize)
		if expectedSize > 0 {
			// Use crypto/rand for truly random data (tests randomness handling)
			rand.Read(result)
		}
		return result

	case edgeType < 35: // Structured but invalid (valid structure, invalid values)
		result := make([]byte, expectedSize)
		// Create structured data that looks valid but isn't
		for i := 0; i < expectedSize; i += 8 {
			if i+8 <= expectedSize {
				binary.BigEndian.PutUint64(result[i:i+8], uint64(i)*0xDEADBEEFCAFEBABE)
			} else {
				// Fill remaining bytes with pattern
				for j := i; j < expectedSize; j++ {
					result[j] = byte(0xDE)
				}
				break
			}
		}
		return result

	case edgeType < 50: // Near-valid values (off by one from valid)
		result := make([]byte, len(input))
		copy(result, input)
		if len(result) > 0 {
			// Increment/decrement last byte to make near-valid
			if t.rng.Intn(2) == 0 {
				result[len(result)-1]++
			} else {
				result[len(result)-1]--
			}
		}
		return result

	case edgeType < 65: // Repeated patterns that might cause issues
		result := make([]byte, expectedSize)
		patterns := [][]byte{
			{0xDE, 0xAD, 0xBE, 0xEF},
			{0xCA, 0xFE, 0xBA, 0xBE},
			{0x13, 0x37},
			{0x00, 0x00, 0x00, 0x01},
			{0xFF, 0xFF, 0xFF, 0xFF},
		}
		pattern := patterns[t.rng.Intn(len(patterns))]
		for i := 0; i < expectedSize; i++ {
			result[i] = pattern[i%len(pattern)]
		}
		return result

	case edgeType < 80: // Arithmetic overflow/underflow values
		result := make([]byte, expectedSize)
		// Fill with values that might cause arithmetic issues
		overflowPatterns := [][]byte{
			{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, // Max uint64
			{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // Min int64
			{0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, // Max int64
		}
		pattern := overflowPatterns[t.rng.Intn(len(overflowPatterns))]
		for i := 0; i < expectedSize; i += len(pattern) {
			remaining := expectedSize - i
			if remaining >= len(pattern) {
				copy(result[i:i+len(pattern)], pattern)
			} else {
				copy(result[i:i+remaining], pattern[:remaining])
			}
		}
		return result

	default: // Mixed corruption (combine multiple techniques)
		result := t.applyDataCorruption(input)
		result = t.applyBoundaryConditions(result, len(result))
		return result
	}
}

// TransformPrecompileInput applies specific transformations for precompile inputs
func (t *InputTransformer) TransformPrecompileInput(input []byte, precompileAddr byte) []byte {
	// For critical precompile operations, we need to be more careful about size preservation
	preserveSize := t.shouldPreserveSize(precompileAddr, len(input))

	var transformed []byte
	if preserveSize {
		// Apply only data content transformations, preserve size
		transformed = t.transformContentOnly(input)
	} else {
		// Apply general transformations including size changes
		transformed = t.TransformInput(input, len(input))
	}

	// Ensure minimum required size for critical operations
	transformed = t.ensureMinimumSize(transformed, input, precompileAddr)

	// Apply precompile-specific transformations with low probability
	if t.rng.Float64() > 0.05 { // 5% chance for precompile-specific transforms
		return transformed
	}

	switch precompileAddr {
	case 0x01: // ECRECOVER
		return t.transformEcrecoverInput(transformed)
	case 0x02, 0x03: // SHA256, RIPEMD160
		return t.transformHashInput(transformed)
	case 0x05: // MODEXP
		return t.transformModexpInput(transformed)
	case 0x06, 0x07, 0x08: // BN256 operations
		return t.transformBN256Input(transformed)
	case 0x09: // BLAKE2F
		return t.transformBlake2fInput(transformed)
	case 0x0A: // KZG Point Evaluation
		return t.transformKZGInput(transformed)
	case 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11: // BLS12-381 operations
		return t.transformBLS12Input(transformed)
	default:
		return transformed
	}
}

// transformEcrecoverInput creates invalid ECRECOVER inputs
func (t *InputTransformer) transformEcrecoverInput(input []byte) []byte {
	// ECRECOVER expects 128 bytes: hash(32) + v(32) + r(32) + s(32)
	result := make([]byte, 128)
	copy(result, input)

	edgeCase := t.rng.Intn(100)
	switch {
	case edgeCase < 25: // Invalid v value (should be 27 or 28)
		binary.BigEndian.PutUint64(result[56:64], uint64(t.rng.Intn(256)+1)) // Invalid v
	case edgeCase < 50: // Zero r or s (invalid signature)
		if t.rng.Intn(2) == 0 {
			for i := 64; i < 96; i++ {
				result[i] = 0
			} // Zero r
		} else {
			for i := 96; i < 128; i++ {
				result[i] = 0
			} // Zero s
		}
	case edgeCase < 75: // r or s >= curve order
		target := 64 + t.rng.Intn(2)*32 // r or s
		result[target] = 0xFF           // Make it large
		for i := target + 1; i < target+8; i++ {
			result[i] = 0xFF
		}
	default: // Malformed signature values
		for i := 64; i < 128; i++ {
			result[i] = 0xFF
		}
	}
	return result
}

// transformHashInput creates edge cases for hash functions
func (t *InputTransformer) transformHashInput(input []byte) []byte {
	edgeCase := t.rng.Intn(100)
	switch {
	case edgeCase < 20: // Empty input
		return []byte{}
	case edgeCase < 30: // Single byte
		return []byte{byte(t.rng.Intn(256))}
	case edgeCase < 50: // Very large input
		large := make([]byte, 100000+t.rng.Intn(100000))
		pattern := byte(t.rng.Intn(256))
		for i := range large {
			large[i] = pattern
		}
		return large
	case edgeCase < 70: // Pattern that might cause hash collisions
		return []byte{0x61, 0x62, 0x63} // "abc" - common test input
	default: // Keep original input but apply content transformations
		return t.applyDataCorruption(input)
	}
}

// transformModexpInput creates invalid MODEXP inputs
func (t *InputTransformer) transformModexpInput(input []byte) []byte {
	// MODEXP format: base_len(32) + exp_len(32) + mod_len(32) + base + exp + mod
	result := make([]byte, 96) // At least the length headers

	edgeCase := t.rng.Intn(100)
	switch {
	case edgeCase < 20: // Zero modulus
		binary.BigEndian.PutUint64(result[24:32], 32) // mod_len = 32
		result = append(result, make([]byte, 64)...)  // base + exp (zeros)
		result = append(result, make([]byte, 32)...)  // mod = 0
	case edgeCase < 40: // Modulus = 1
		binary.BigEndian.PutUint64(result[24:32], 32) // mod_len = 32
		result = append(result, make([]byte, 64)...)  // base + exp
		modulus := make([]byte, 32)
		modulus[31] = 1 // mod = 1
		result = append(result, modulus...)
	case edgeCase < 60: // Extremely large lengths
		binary.BigEndian.PutUint64(result[0:8], 0xFFFFFFFF)   // Huge base_len
		binary.BigEndian.PutUint64(result[8:16], 0xFFFFFFFF)  // Huge exp_len
		binary.BigEndian.PutUint64(result[16:24], 0xFFFFFFFF) // Huge mod_len
	default: // Mismatched lengths
		binary.BigEndian.PutUint64(result[0:8], 64)   // base_len = 64
		binary.BigEndian.PutUint64(result[8:16], 32)  // exp_len = 32
		binary.BigEndian.PutUint64(result[16:24], 16) // mod_len = 16
		result = append(result, make([]byte, 32)...)  // Only 32 bytes of data
	}
	return result
}

// transformBN256Input creates invalid BN256 curve inputs
func (t *InputTransformer) transformBN256Input(input []byte) []byte {
	// BN256 points are 64 bytes (32 x + 32 y), pairings use multiple points
	if len(input) < 64 {
		return input
	}

	result := make([]byte, len(input))
	copy(result, input)

	edgeCase := t.rng.Intn(100)
	switch {
	case edgeCase < 30: // Point not on curve
		// Set coordinates that don't satisfy curve equation
		for i := 0; i < 64; i++ {
			result[i] = byte(i) // Sequential pattern unlikely to be on curve
		}
	case edgeCase < 60: // Coordinates >= field modulus
		result[0] = 0xFF // Make x coordinate very large
		for i := 1; i < 8; i++ {
			result[i] = 0xFF
		}
	default: // Invalid point at infinity representation
		// Some clients expect specific encoding for point at infinity
		for i := 0; i < 32; i++ {
			result[i] = 0xFF // Invalid "infinity" encoding
		}
	}
	return result
}

// transformBlake2fInput creates invalid BLAKE2F inputs
func (t *InputTransformer) transformBlake2fInput(input []byte) []byte {
	// BLAKE2F expects exactly 213 bytes
	result := make([]byte, 213)
	copy(result, input)

	edgeCase := t.rng.Intn(100)
	switch {
	case edgeCase < 50: // Invalid rounds (too high)
		binary.BigEndian.PutUint32(result[0:4], 0xFFFFFFFF) // Max rounds
	default: // Invalid final flag (should be 0 or 1)
		result[212] = byte(t.rng.Intn(254) + 2) // 2-255 (invalid)
	}
	return result
}

// transformKZGInput creates invalid KZG Point Evaluation inputs
func (t *InputTransformer) transformKZGInput(input []byte) []byte {
	// KZG expects 192 bytes: versioned_hash(32) + z(32) + y(32) + commitment(48) + proof(48)
	result := make([]byte, 192)
	copy(result, input)

	edgeCase := t.rng.Intn(100)
	switch {
	case edgeCase < 25: // Invalid versioned hash prefix
		result[0] = byte(t.rng.Intn(255) + 1) // Wrong version byte
	case edgeCase < 50: // Invalid BLS field elements
		for i := 144; i < 192; i++ {
			result[i] = 0xFF // Make proof field elements invalid
		}
	case edgeCase < 75: // Mismatched commitment/proof
		// Commitment for one polynomial, proof for another
		for i := 96; i < 144; i++ {
			result[i] = byte(i % 256)
		}
	default: // Invalid point on BLS curve
		result[96] = 0xFF // Invalid commitment
		for i := 97; i < 104; i++ {
			result[i] = 0xFF
		}
	}
	return result
}

// transformBLS12Input creates invalid BLS12-381 inputs
func (t *InputTransformer) transformBLS12Input(input []byte) []byte {
	result := make([]byte, len(input))
	copy(result, input)

	edgeCase := t.rng.Intn(100)
	switch {
	case edgeCase < 20: // Invalid field elements (> modulus)
		// BLS12-381 field elements are 48 bytes, but stored in 64-byte slots
		for i := 0; i < len(result); i += 64 {
			if i+48 <= len(result) {
				result[i] = 0xFF // Ensure > field modulus
				for j := i + 1; j < i+8; j++ {
					if j < len(result) {
						result[j] = 0xFF
					}
				}
			}
		}
	case edgeCase < 40: // Points not in correct subgroup
		// Valid field elements but not in prime-order subgroup
		for i := 16; i < len(result) && i < 64; i += 48 {
			if i+32 < len(result) {
				for j := i; j < i+32; j++ {
					result[j] = byte(j % 256) // Structured pattern unlikely to be in subgroup
				}
			}
		}
	case edgeCase < 60: // Invalid point encoding
		// Wrong compression flags or invalid point representation
		if len(result) > 0 {
			result[0] |= 0xE0 // Set invalid compression flags
		}
	case edgeCase < 80: // Pairing input count mismatch
		// For pairing, input should be multiple of 384 bytes (128 G1 + 256 G2)
		// Create non-multiple size to test validation
		if len(result) >= 384 {
			return result[:len(result)-1] // Remove one byte
		}
	default: // MSM input format violation
		// MSM expects (point, scalar) pairs
		// Create input where point/scalar boundaries are misaligned
		if len(result) >= 160 { // 128 point + 32 scalar
			// Shift data to misalign point/scalar boundaries
			shifted := make([]byte, len(result))
			copy(shifted[1:], result[:len(result)-1])
			result = shifted
		}
	}
	return result
}

// shouldPreserveSize determines if we should preserve input size for specific precompiles
func (t *InputTransformer) shouldPreserveSize(precompileAddr byte, inputSize int) bool {
	switch precompileAddr {
	case 0x01: // ECRECOVER - expects exactly 128 bytes
		return inputSize == 128
	case 0x06, 0x07: // BN256 EC operations - expect 64 bytes for points
		return inputSize == 64 || inputSize == 96 || inputSize == 128
	case 0x08: // BN256 pairing - expects multiple of 192 bytes
		return inputSize > 0 && inputSize%192 == 0
	case 0x09: // BLAKE2F - expects exactly 213 bytes
		return inputSize == 213
	case 0x0A: // KZG Point Evaluation - expects exactly 192 bytes
		return inputSize == 192
	case 0x0B, 0x0D: // BLS12 G1/G2 Add - expect 128/256 bytes
		return inputSize == 128 || inputSize == 256
	case 0x0C, 0x0E: // BLS12 G1/G2 MSM - expect multiples of 160/288 bytes
		return (inputSize > 0 && inputSize%160 == 0) || (inputSize > 0 && inputSize%288 == 0)
	case 0x0F: // BLS12 Pairing - expects multiple of 384 bytes
		return inputSize > 0 && inputSize%384 == 0
	case 0x10: // BLS12 Map Fp to G1 - expects 64 bytes
		return inputSize == 64
	case 0x11: // BLS12 Map Fp2 to G2 - expects 128 bytes
		return inputSize == 128
	default:
		return false // For hash functions and MODEXP, size changes are okay
	}
}

// transformContentOnly applies transformations that preserve input size
func (t *InputTransformer) transformContentOnly(input []byte) []byte {
	// 10% chance to apply content-only transformations
	if t.rng.Float64() > 0.10 {
		return input
	}

	transformType := t.rng.Intn(100)

	switch {
	case transformType < 30: // Data corruption (preserves size)
		return t.applyDataCorruption(input)
	case transformType < 60: // Boundary conditions (preserves size)
		return t.applyBoundaryConditions(input, len(input))
	case transformType < 80: // Format violations (preserves size)
		return t.applyFormatViolations(input, len(input))
	default: // Pattern-based transformations
		result := make([]byte, len(input))
		copy(result, input)

		// Apply pattern transformations
		patternType := t.rng.Intn(4)
		switch patternType {
		case 0: // All zeros
			for i := range result {
				result[i] = 0x00
			}
		case 1: // All ones
			for i := range result {
				result[i] = 0xFF
			}
		case 2: // Alternating pattern
			for i := range result {
				if i%2 == 0 {
					result[i] = 0xAA
				} else {
					result[i] = 0x55
				}
			}
		default: // Incremental pattern
			for i := range result {
				result[i] = byte(i % 256)
			}
		}
		return result
	}
}

// ensureMinimumSize ensures transformed data meets minimum size requirements
func (t *InputTransformer) ensureMinimumSize(transformed []byte, original []byte, precompileAddr byte) []byte {
	minSize := t.getMinimumRequiredSize(precompileAddr, len(original))

	if len(transformed) >= minSize {
		return transformed
	}

	// If transformed is too small, pad it to minimum size
	result := make([]byte, minSize)
	copy(result, transformed)

	// Fill remaining bytes with pattern or zeros
	for i := len(transformed); i < minSize; i++ {
		if t.rng.Intn(2) == 0 {
			result[i] = 0x00 // Zero padding
		} else {
			result[i] = 0xAA // Pattern padding
		}
	}

	return result
}

// getMinimumRequiredSize returns the minimum size needed for a precompile operation
func (t *InputTransformer) getMinimumRequiredSize(precompileAddr byte, originalSize int) int {
	switch precompileAddr {
	case 0x01: // ECRECOVER
		return 128
	case 0x06, 0x07: // BN256 EC operations
		return 64
	case 0x08: // BN256 pairing
		return 192
	case 0x09: // BLAKE2F
		return 213
	case 0x0A: // KZG Point Evaluation
		return 192
	case 0x0B: // BLS12 G1 Add
		return 128
	case 0x0C: // BLS12 G1 MSM
		return 160
	case 0x0D: // BLS12 G2 Add
		return 256
	case 0x0E: // BLS12 G2 MSM
		return 288
	case 0x0F: // BLS12 Pairing
		return 384
	case 0x10: // BLS12 Map Fp to G1
		return 64
	case 0x11: // BLS12 Map Fp2 to G2
		return 128
	default:
		return 0 // No minimum size requirement
	}
}
