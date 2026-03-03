package calltxfuzz

import (
	evmfuzz "github.com/ethpandaops/spamoor/scenarios/evm-fuzz"
)

// CalldataGenerator generates fuzzed calldata for calling deployed contracts.
type CalldataGenerator struct {
	rng     *evmfuzz.DeterministicRNG
	maxSize uint64
}

// NewCalldataGenerator creates a new calldata generator.
func NewCalldataGenerator(rng *evmfuzz.DeterministicRNG, maxSize uint64) *CalldataGenerator {
	return &CalldataGenerator{
		rng:     rng,
		maxSize: maxSize,
	}
}

// Generate produces fuzzed calldata with various patterns.
func (g *CalldataGenerator) Generate() []byte {
	pattern := g.rng.Intn(100)

	switch {
	case pattern < 30:
		return g.generateRandomBytes()
	case pattern < 55:
		return g.generateABILike()
	case pattern < 70:
		return g.generateEdgeCase()
	case pattern < 85:
		return g.generateStorageKeyPattern()
	default:
		return g.generateRepeatedPattern()
	}
}

// generateRandomBytes produces uniform random calldata.
func (g *CalldataGenerator) generateRandomBytes() []byte {
	size := g.rng.Intn(int(g.maxSize) + 1)
	if size == 0 {
		return nil
	}

	return g.rng.Bytes(size)
}

// generateABILike produces calldata that looks like an ABI-encoded call:
// 4-byte function selector + 32-byte-aligned parameters.
func (g *CalldataGenerator) generateABILike() []byte {
	// Function selector (4 bytes)
	selector := g.rng.Bytes(4)

	// 1-8 parameters, each 32 bytes
	paramCount := g.rng.Intn(8) + 1
	size := 4 + paramCount*32
	if uint64(size) > g.maxSize {
		paramCount = int(g.maxSize-4) / 32
		if paramCount < 1 {
			return selector
		}
		size = 4 + paramCount*32
	}

	data := make([]byte, size)
	copy(data, selector)

	for i := 0; i < paramCount; i++ {
		offset := 4 + i*32
		// 50% chance of random data, 50% chance of structured values
		if g.rng.Float64() < 0.5 {
			copy(data[offset:offset+32], g.rng.Bytes(32))
		} else {
			// Structured value: small integer in last 8 bytes
			val := g.rng.Uint64()
			for j := 0; j < 8; j++ {
				data[offset+24+j] = byte(val >> (56 - j*8))
			}
		}
	}

	return data
}

// generateEdgeCase produces edge-case calldata sizes.
func (g *CalldataGenerator) generateEdgeCase() []byte {
	choice := g.rng.Intn(5)
	switch choice {
	case 0: // Empty
		return nil
	case 1: // Single byte
		return g.rng.Bytes(1)
	case 2: // Exactly 4 bytes (just a selector)
		return g.rng.Bytes(4)
	case 3: // 32 bytes (one word)
		return g.rng.Bytes(32)
	default: // Max size
		if g.maxSize == 0 {
			return nil
		}
		return g.rng.Bytes(int(g.maxSize))
	}
}

// generateStorageKeyPattern produces calldata encoding storage keys
// matching the generator's deterministic key patterns (keys 0-31).
func (g *CalldataGenerator) generateStorageKeyPattern() []byte {
	// ABI-like: selector + key as uint256
	data := make([]byte, 36)
	copy(data[:4], g.rng.Bytes(4))

	// Storage key in the low byte (matches generateStoragePattern)
	data[35] = byte(g.rng.Intn(32))

	return data
}

// generateRepeatedPattern produces calldata with repeated byte patterns.
func (g *CalldataGenerator) generateRepeatedPattern() []byte {
	size := g.rng.Intn(int(g.maxSize) + 1)
	if size == 0 {
		return nil
	}

	data := make([]byte, size)
	choice := g.rng.Intn(3)

	switch choice {
	case 0: // Same byte repeated
		b := byte(g.rng.Intn(256))
		for i := range data {
			data[i] = b
		}
	case 1: // Alternating pattern
		a := byte(g.rng.Intn(256))
		b := byte(g.rng.Intn(256))
		for i := range data {
			if i%2 == 0 {
				data[i] = a
			} else {
				data[i] = b
			}
		}
	default: // Incrementing
		for i := range data {
			data[i] = byte(i % 256)
		}
	}

	return data
}
