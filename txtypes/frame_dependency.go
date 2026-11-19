package txtypes

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
)

// EIP-8288 dependency verification frames.
//
// A dependency frame runs no code. Its data declares (scheme, data_hash,
// verification_key_hash) triples the block has to prove with one recursive STARK; the
// signatures and proofs behind them travel in mempool wrappers and never reach the
// chain.
const (
	DependencyTripleBytes   = 96
	MaxDependenciesPerFrame = 256

	LeanSphincsVerificationGas = 3_000
	LeanStarkVerificationGas   = 30_000

	// Public mempool limits, counted over a transaction's deduplicated dependencies.
	MaxSigsPerTx   = 16
	MaxStarksPerTx = 1
)

// DependencyScheme identifies what a dependency's proof is: the low byte of the
// triple's first word.
type DependencyScheme uint8

// Dependency schemes.
const (
	DepSchemeLeanSphincs DependencyScheme = 0x10

	// DepSchemeLeanStark is reserved by the EIP until get_deps_hash is defined for it.
	// The devnet clients accept it, with the digest DependenciesHash describes.
	DepSchemeLeanStark DependencyScheme = 0x11
)

// leanStarkDepsDomain closes the digest of a dependency list holding leanSTARK entries.
const leanStarkDepsDomain = "eip8288/leanstark-deps"

// Known reports whether the scheme is one a dependency frame may declare.
func (s DependencyScheme) Known() bool {
	return s == DepSchemeLeanSphincs || s == DepSchemeLeanStark
}

// VerificationGas returns the execution gas a dependency frame budgets per triple of
// this scheme.
func (s DependencyScheme) VerificationGas() uint64 {
	if s == DepSchemeLeanStark {
		return LeanStarkVerificationGas
	}

	return LeanSphincsVerificationGas
}

// String names the scheme.
func (s DependencyScheme) String() string {
	switch s {
	case DepSchemeLeanSphincs:
		return "leanSPHINCS"
	case DepSchemeLeanStark:
		return "leanSTARK"
	default:
		return fmt.Sprintf("unknown(0x%02x)", uint8(s))
	}
}

// FrameDependency is one triple of a dependency frame.
//
// DataHash is the message hash of a leanSPHINCS signature or the public inputs hash of
// a leanSTARK. VerificationKey is the leanSPHINCS public key itself, PK.seed || PK.root,
// or the hash of a leanSTARK verification key.
type FrameDependency struct {
	Scheme          DependencyScheme
	DataHash        common.Hash
	VerificationKey common.Hash
}

// SignsTransaction reports whether the dependency is declared over the containing
// transaction's signature hash, which a frame spells as a zero data hash: the hash
// covers the frame's data, so the frame cannot carry it.
func (d FrameDependency) SignsTransaction() bool {
	return d.DataHash == (common.Hash{})
}

// EncodeDependencies returns the triples in the 96-byte form a dependency frame's data
// carries: the scheme as a 32-byte big-endian word, then the two hashes.
func EncodeDependencies(deps []FrameDependency) []byte {
	data := make([]byte, len(deps)*DependencyTripleBytes)

	for i, dep := range deps {
		triple := data[i*DependencyTripleBytes:]
		triple[31] = uint8(dep.Scheme)
		copy(triple[32:64], dep.DataHash[:])
		copy(triple[64:96], dep.VerificationKey[:])
	}

	return data
}

// DependencyFrame builds a dependency verification frame declaring deps, with the
// execution budget the triples fix. A zero DataHash declares the dependency over the
// transaction's own signature hash.
func DependencyFrame(deps ...FrameDependency) *Frame {
	gas := uint64(0)
	for _, dep := range deps {
		gas += dep.Scheme.VerificationGas()
	}

	return &Frame{
		Mode:   FrameModeDepVerify,
		Flags:  ApproveNone,
		Target: nil,
		Limits: FrameLimits{Execution: gas},
		Value:  new(uint256.Int),
		Data:   EncodeDependencies(deps),
	}
}

// IsDependency reports whether the frame is a dependency verification frame.
func (f *Frame) IsDependency() bool {
	return f.Mode == FrameModeDepVerify
}

// Dependencies returns the triples a dependency frame declares, in frame order and as
// written: duplicates are kept and a zero data hash is not resolved. It returns nil for
// any other frame, and ignores a trailing partial triple, which ValidatePayload rejects.
func (f *Frame) Dependencies() []FrameDependency {
	if !f.IsDependency() {
		return nil
	}

	count := len(f.Data) / DependencyTripleBytes
	deps := make([]FrameDependency, 0, count)

	for i := 0; i < count; i++ {
		triple := f.Data[i*DependencyTripleBytes : (i+1)*DependencyTripleBytes]

		deps = append(deps, FrameDependency{
			Scheme:          DependencyScheme(triple[31]),
			DataHash:        common.BytesToHash(triple[32:64]),
			VerificationKey: common.BytesToHash(triple[64:96]),
		})
	}

	return deps
}

// DependencyCount returns the number of triples the transaction's dependency frames
// declare, duplicates included.
func (tx *FrameTx) DependencyCount() int {
	count := 0

	for _, frame := range tx.Frames {
		if frame.IsDependency() {
			count += len(frame.Data) / DependencyTripleBytes
		}
	}

	return count
}

// RecursiveStarkGas returns the transaction's share of the block's recursive STARK,
// part of its intrinsic cost: LeanStarkVerificationGas per declared triple of either
// scheme, duplicates included. It comes on top of the budgets of the dependency frames.
func (tx *FrameTx) RecursiveStarkGas() uint64 {
	return saturatingMul(uint64(tx.DependencyCount()), LeanStarkVerificationGas)
}

// Dependencies returns the transaction's dependencies as the block has to prove them:
// every triple of every dependency frame, a zero data hash replaced by the
// transaction's signature hash, sorted and deduplicated.
func (tx *FrameTx) Dependencies() ([]FrameDependency, error) {
	count := tx.DependencyCount()
	if count == 0 {
		return nil, nil
	}

	deps := make([]FrameDependency, 0, count)
	sigHash := common.Hash{}
	resolved := false

	for _, frame := range tx.Frames {
		for _, dep := range frame.Dependencies() {
			if dep.SignsTransaction() {
				if !resolved {
					hash, err := tx.SigHash()
					if err != nil {
						return nil, err
					}

					sigHash, resolved = hash, true
				}

				dep.DataHash = sigHash
			}

			deps = append(deps, dep)
		}
	}

	return CanonicalDependencies(deps), nil
}

// CanonicalDependencies returns deps sorted by scheme, data hash and verification key
// with duplicates removed, the order the dependency digest is taken in. The union of
// several transactions' dependencies is canonicalized the same way.
func CanonicalDependencies(deps []FrameDependency) []FrameDependency {
	sorted := make([]FrameDependency, len(deps))
	copy(sorted, deps)

	sort.Slice(sorted, func(i, j int) bool { return compareDependencies(sorted[i], sorted[j]) < 0 })

	unique := sorted[:0]

	for _, dep := range sorted {
		if len(unique) == 0 || dep != unique[len(unique)-1] {
			unique = append(unique, dep)
		}
	}

	return unique
}

// compareDependencies orders two triples by scheme, data hash and verification key.
func compareDependencies(a, b FrameDependency) int {
	if a.Scheme != b.Scheme {
		if a.Scheme < b.Scheme {
			return -1
		}

		return 1
	}

	if order := bytes.Compare(a.DataHash[:], b.DataHash[:]); order != 0 {
		return order
	}

	return bytes.Compare(a.VerificationKey[:], b.VerificationKey[:])
}

// DependenciesHash returns get_deps_hash of a dependency list: the digest a recursive
// STARK commits to, and a block's block_deps_hash when taken over the union of its
// transactions' dependencies. The list is canonicalized first.
//
// The EIP leaves DEPS_HASH open and defines the digest for leanSPHINCS lists only. This
// is what the devnet clients compute: Keccak-256 over the data_hash || verification_key
// entries, and for a list holding any other scheme the digests of the leanSPHINCS
// entries and of the remaining ones, joined under a domain tag.
func DependenciesHash(deps []FrameDependency) common.Hash {
	deps = CanonicalDependencies(deps)

	sphincs := make([]byte, 0, len(deps)*2*common.HashLength)
	others := make([]byte, 0, 2*common.HashLength)

	for _, dep := range deps {
		if dep.Scheme == DepSchemeLeanSphincs {
			sphincs = append(sphincs, dep.DataHash[:]...)
			sphincs = append(sphincs, dep.VerificationKey[:]...)

			continue
		}

		others = append(others, dep.DataHash[:]...)
		others = append(others, dep.VerificationKey[:]...)
	}

	sphincsHash := crypto.Keccak256Hash(sphincs)
	if len(others) == 0 {
		return sphincsHash
	}

	domain := make([]byte, common.HashLength)
	copy(domain, leanStarkDepsDomain)

	return crypto.Keccak256Hash(sphincsHash[:], crypto.Keccak256(others), domain)
}

// validateDependencyFrame checks the validity rules of the dependency frame at index.
// The rule that it may not close an atomic batch is covered by the batch rules, which
// only let the flag point at a DEFAULT or SENDER frame.
func (tx *FrameTx) validateDependencyFrame(index int, frame *Frame) error {
	if frame.Target != nil {
		return fmt.Errorf("%w: dependency frame %d has a target", ErrInvalidFrameTx, index)
	}

	if frame.Flags != ApproveNone {
		return fmt.Errorf("%w: dependency frame %d sets flags 0x%02x", ErrInvalidFrameTx, index, frame.Flags)
	}

	if frame.Value != nil && !frame.Value.IsZero() {
		return fmt.Errorf("%w: dependency frame %d carries value", ErrInvalidFrameTx, index)
	}

	if len(frame.Data) == 0 || len(frame.Data)%DependencyTripleBytes != 0 {
		return fmt.Errorf("%w: dependency frame %d has %d data bytes, expected whole %d-byte triples",
			ErrInvalidFrameTx, index, len(frame.Data), DependencyTripleBytes)
	}

	count := len(frame.Data) / DependencyTripleBytes
	if count > MaxDependenciesPerFrame {
		return fmt.Errorf("%w: dependency frame %d declares %d dependencies, limit is %d",
			ErrInvalidFrameTx, index, count, MaxDependenciesPerFrame)
	}

	gas := uint64(0)

	for i := 0; i < count; i++ {
		triple := frame.Data[i*DependencyTripleBytes : (i+1)*DependencyTripleBytes]

		for _, b := range triple[:31] {
			if b != 0 {
				return fmt.Errorf("%w: dependency %d of frame %d has a non-zero scheme padding",
					ErrInvalidFrameTx, i, index)
			}
		}

		scheme := DependencyScheme(triple[31])
		if !scheme.Known() {
			return fmt.Errorf("%w: dependency %d of frame %d has unknown scheme 0x%02x",
				ErrInvalidFrameTx, i, index, uint8(scheme))
		}

		gas += scheme.VerificationGas()
	}

	if frame.Limits.State != 0 {
		return fmt.Errorf("%w: dependency frame %d budgets state gas", ErrInvalidFrameTx, index)
	}

	if frame.Limits.Execution != gas {
		return fmt.Errorf("%w: dependency frame %d budgets %d execution gas, its dependencies fix %d",
			ErrInvalidFrameTx, index, frame.Limits.Execution, gas)
	}

	return nil
}

// validateDependencyLimits checks the public mempool's per-transaction dependency
// limits, which count the deduplicated list: repeating a declaration is legal, only
// wasteful.
func (tx *FrameTx) validateDependencyLimits() error {
	deps, err := tx.Dependencies()
	if err != nil {
		return err
	}

	sigs, starks := 0, 0

	for _, dep := range deps {
		if dep.Scheme == DepSchemeLeanStark {
			starks++
		} else {
			sigs++
		}
	}

	if sigs > MaxSigsPerTx {
		return fmt.Errorf("%w: %d leanSPHINCS dependencies, limit is %d", ErrMempoolPolicy, sigs, MaxSigsPerTx)
	}

	if starks > MaxStarksPerTx {
		return fmt.Errorf("%w: %d leanSTARK dependencies, limit is %d", ErrMempoolPolicy, starks, MaxStarksPerTx)
	}

	return nil
}
