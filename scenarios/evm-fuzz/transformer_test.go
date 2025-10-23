package evmfuzz

import (
	"testing"
)

func TestInputTransformers(t *testing.T) {
	rng := NewDeterministicRNGWithSeed(42, "test")
	transformer := NewInputTransformer(rng)

	// Test basic transformation
	input := []byte{0x01, 0x02, 0x03, 0x04}
	transformed := transformer.TransformInput(input, 4)

	// Should either be the same (85% chance) or transformed
	if len(transformed) == 0 && len(input) > 0 {
		t.Logf("Input was transformed to empty (size variation)")
	} else if len(transformed) != len(input) {
		t.Logf("Input size changed from %d to %d", len(input), len(transformed))
	} else {
		t.Logf("Input transformation: %x -> %x", input, transformed)
	}

	// Test precompile-specific transformations
	ecrecoverInput := make([]byte, 128) // ECRECOVER format
	ecrecoverTransformed := transformer.TransformPrecompileInput(ecrecoverInput, 0x01)
	if len(ecrecoverTransformed) != 128 {
		t.Logf("ECRECOVER input size changed from 128 to %d", len(ecrecoverTransformed))
	}

	// Test BLS12 transformations
	bls12Input := make([]byte, 128) // BLS12 G1 point
	bls12Transformed := transformer.TransformPrecompileInput(bls12Input, 0x0b)
	t.Logf("BLS12 transformation applied, size: %d -> %d", len(bls12Input), len(bls12Transformed))

	// Test size transformations specifically
	sizeTransformed := transformer.applySizeTransformations(input, 4)
	t.Logf("Size transformation: %d -> %d bytes", len(input), len(sizeTransformed))

	// Test data corruption
	corrupted := transformer.applyDataCorruption(input)
	t.Logf("Data corruption: %x -> %x", input, corrupted)
}
