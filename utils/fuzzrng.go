package utils

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"strings"
)

// DeterministicRNG is a seeded pseudo-random source for fuzzing scenarios.
//
// Fuzzers need their output to be reproducible from a short command line: a run is
// identified by a seed, and every transaction index draws from its own stream derived
// from that seed, so a single failing transaction can be regenerated on its own without
// replaying the ones before it.
type DeterministicRNG struct {
	state   uint64
	counter uint64
}

// ParseHexSeed decodes a hex seed, tolerating a 0x prefix and an odd digit count.
func ParseHexSeed(seed string) ([]byte, error) {
	seed = strings.TrimPrefix(seed, "0x")

	if len(seed)%2 == 1 {
		seed = "0" + seed
	}

	return hex.DecodeString(seed)
}

// NewDeterministicRNGWithSeed creates a generator for one transaction index.
//
// The seed is hashed together with the index so that neighbouring indices produce
// unrelated streams; a seed that is not valid hex is used as raw bytes rather than
// rejected, so an operator can pass a memorable string.
func NewDeterministicRNGWithSeed(txID uint64, baseSeed string) *DeterministicRNG {
	h := sha256.New()

	if baseSeed != "" {
		seedBytes, err := ParseHexSeed(baseSeed)
		if err != nil {
			h.Write([]byte(baseSeed))
		} else {
			h.Write(seedBytes)
		}
	} else {
		binary.Write(h, binary.LittleEndian, uint64(42))
	}

	binary.Write(h, binary.LittleEndian, txID)
	binary.Write(h, binary.LittleEndian, uint64(0x1337DEADBEEF))

	seed := binary.LittleEndian.Uint64(h.Sum(nil)[:8])
	if seed == 0 {
		seed = 1
	}

	return &DeterministicRNG{
		state:   seed,
		counter: 0,
	}
}

// Uint64 returns the next value from the stream, using xorshift64*.
func (r *DeterministicRNG) Uint64() uint64 {
	r.counter++
	r.state ^= r.state >> 12
	r.state ^= r.state << 25
	r.state ^= r.state >> 27

	return r.state * 0x2545F4914F6CDD1D
}

// Intn returns a value in [0, n), or 0 when n is not positive.
func (r *DeterministicRNG) Intn(n int) int {
	if n <= 0 {
		return 0
	}

	return int(r.Uint64() % uint64(n))
}

// Uint64n returns a value in [0, n), or 0 when n is zero.
func (r *DeterministicRNG) Uint64n(n uint64) uint64 {
	if n == 0 {
		return 0
	}

	return r.Uint64() % n
}

// Float64 returns a value in [0, 1].
func (r *DeterministicRNG) Float64() float64 {
	return float64(r.Uint64()) / float64(^uint64(0))
}

// Bytes returns n pseudo-random bytes.
func (r *DeterministicRNG) Bytes(n int) []byte {
	result := make([]byte, n)

	for i := 0; i < n; i += 8 {
		val := r.Uint64()
		for j := 0; j < 8 && i+j < n; j++ {
			result[i+j] = byte(val >> (j * 8))
		}
	}

	return result
}

// Draws returns how many values have been taken from the stream, which lets a caller
// assert that a generation pass consumed a stable number of draws.
func (r *DeterministicRNG) Draws() uint64 { return r.counter }
