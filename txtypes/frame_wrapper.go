package txtypes

import (
	"bytes"
	"errors"
	"fmt"
	"sort"

	"github.com/ethereum/go-ethereum/rlp"
)

// EIP-8288 mempool wrappers.
//
// A transaction with dependency frames does not travel alone: it is submitted in a
// wrapper carrying the proofs of its dependencies, through eth_sendProofWrapper rather
// than eth_sendRawTransaction.
const (
	// ProofWrapperModeDirect carries one proof per dependency.
	ProofWrapperModeDirect = 0

	// ProofWrapperModeRecursive carries one recursive STARK over all dependencies.
	ProofWrapperModeRecursive = 1

	// MaxWrapperBytes bounds the complete wrapper encoding in either mode.
	MaxWrapperBytes = 2 * 1024 * 1024

	// Dependency limits of a direct wrapper.
	MaxLeanSigDepsPerWrapper   = 16
	MaxLeanStarkDepsPerWrapper = 1
)

// ErrInvalidProofWrapper is the base error for a malformed mempool wrapper.
var ErrInvalidProofWrapper = errors.New("invalid proof wrapper")

// ProofWrapper is the mempool wrapper [transactions, mode, [deps, proof_content]].
//
// Transactions are frame transactions in strictly ascending hash order. Dependencies is
// the canonical union of their dependencies. A direct wrapper holds one entry of Proofs
// per dependency, in the same order; a recursive one holds StarkProof, whose public
// input is DependenciesHash of the dependencies.
type ProofWrapper struct {
	Transactions []*Transaction
	Mode         uint8
	Dependencies []FrameDependency
	Proofs       [][]byte
	StarkProof   []byte
}

// NewDirectProofWrapper builds the direct wrapper of txs: it orders the transactions,
// takes the union of their dependencies and looks up the proof of each one in proofs,
// which is keyed by the resolved dependency, the one carrying a transaction's signature
// hash where its frame carries zero.
func NewDirectProofWrapper(txs []*Transaction, proofs map[FrameDependency][]byte) (*ProofWrapper, error) {
	wrapper := &ProofWrapper{Mode: ProofWrapperModeDirect}

	if err := wrapper.setTransactions(txs); err != nil {
		return nil, err
	}

	wrapper.Proofs = make([][]byte, 0, len(wrapper.Dependencies))

	for _, dep := range wrapper.Dependencies {
		proof, ok := proofs[dep]
		if !ok {
			return nil, fmt.Errorf("%w: no proof for the %s dependency over %s under %s",
				ErrInvalidProofWrapper, dep.Scheme, dep.DataHash, dep.VerificationKey)
		}

		wrapper.Proofs = append(wrapper.Proofs, proof)
	}

	return wrapper, nil
}

// NewRecursiveProofWrapper builds the recursive wrapper of txs around starkProof, a
// recursive STARK over the union of their dependencies.
func NewRecursiveProofWrapper(txs []*Transaction, starkProof []byte) (*ProofWrapper, error) {
	wrapper := &ProofWrapper{Mode: ProofWrapperModeRecursive, StarkProof: starkProof}

	if err := wrapper.setTransactions(txs); err != nil {
		return nil, err
	}

	return wrapper, nil
}

// setTransactions stores txs in wrapper order together with the union of their
// dependencies.
func (w *ProofWrapper) setTransactions(txs []*Transaction) error {
	sorted := make([]*Transaction, len(txs))
	copy(sorted, txs)

	sort.Slice(sorted, func(i, j int) bool {
		a, b := sorted[i].Hash(), sorted[j].Hash()

		return bytes.Compare(a[:], b[:]) < 0
	})

	deps, err := wrappedDependencies(sorted)
	if err != nil {
		return err
	}

	w.Transactions = sorted
	w.Dependencies = deps

	return nil
}

// wrappedDependencies returns the canonical union of the transactions' dependencies.
func wrappedDependencies(txs []*Transaction) ([]FrameDependency, error) {
	union := make([]FrameDependency, 0, len(txs))

	for i, tx := range txs {
		frameTx, ok := tx.Inner().(*FrameTx)
		if !ok {
			return nil, fmt.Errorf("%w: transaction %d is not a frame transaction", ErrInvalidProofWrapper, i)
		}

		deps, err := frameTx.Dependencies()
		if err != nil {
			return nil, err
		}

		union = append(union, deps...)
	}

	return CanonicalDependencies(union), nil
}

// Validate checks the wrapper against the rules a mempool node applies before it looks
// at the proofs themselves: the transaction order, the dependency list being the
// canonical union of the transactions' dependencies, and the mode's own shape and
// limits.
func (w *ProofWrapper) Validate() error {
	if len(w.Transactions) == 0 {
		return fmt.Errorf("%w: no transactions", ErrInvalidProofWrapper)
	}

	for i := 1; i < len(w.Transactions); i++ {
		previous, current := w.Transactions[i-1].Hash(), w.Transactions[i].Hash()
		if bytes.Compare(previous[:], current[:]) >= 0 {
			return fmt.Errorf("%w: transactions are not in strictly ascending hash order", ErrInvalidProofWrapper)
		}
	}

	deps, err := wrappedDependencies(w.Transactions)
	if err != nil {
		return err
	}

	if len(deps) != len(w.Dependencies) {
		return fmt.Errorf("%w: %d dependencies listed, the transactions declare %d",
			ErrInvalidProofWrapper, len(w.Dependencies), len(deps))
	}

	for i := range deps {
		if deps[i] != w.Dependencies[i] {
			return fmt.Errorf("%w: dependency %d is not the one the transactions declare",
				ErrInvalidProofWrapper, i)
		}
	}

	switch w.Mode {
	case ProofWrapperModeDirect:
		return w.validateDirect()
	case ProofWrapperModeRecursive:
		if len(w.StarkProof) == 0 {
			return fmt.Errorf("%w: recursive wrapper without a proof", ErrInvalidProofWrapper)
		}

		return nil
	default:
		return fmt.Errorf("%w: unknown mode %d", ErrInvalidProofWrapper, w.Mode)
	}
}

// validateDirect checks a direct wrapper's proof count and dependency limits.
func (w *ProofWrapper) validateDirect() error {
	if len(w.Proofs) != len(w.Dependencies) {
		return fmt.Errorf("%w: %d proofs for %d dependencies",
			ErrInvalidProofWrapper, len(w.Proofs), len(w.Dependencies))
	}

	sigs, starks := 0, 0

	for _, dep := range w.Dependencies {
		if dep.Scheme == DepSchemeLeanStark {
			starks++
		} else {
			sigs++
		}
	}

	if sigs > MaxLeanSigDepsPerWrapper {
		return fmt.Errorf("%w: %d leanSPHINCS dependencies, a direct wrapper holds %d",
			ErrInvalidProofWrapper, sigs, MaxLeanSigDepsPerWrapper)
	}

	if starks > MaxLeanStarkDepsPerWrapper {
		return fmt.Errorf("%w: %d leanSTARK dependencies, a direct wrapper holds %d",
			ErrInvalidProofWrapper, starks, MaxLeanStarkDepsPerWrapper)
	}

	return nil
}

// rlpProofWrapper is the wrapper's wire form. The content list differs by mode, so it
// is carried raw.
type rlpProofWrapper struct {
	Transactions [][]byte
	Mode         uint8
	Content      rlp.RawValue
}

// rlpDirectContent is the content of a direct wrapper.
type rlpDirectContent struct {
	Dependencies [][]byte
	Proofs       [][]byte
}

// rlpRecursiveContent is the content of a recursive wrapper.
type rlpRecursiveContent struct {
	Dependencies [][]byte
	StarkProof   []byte
}

// MarshalBinary returns the wrapper's RLP encoding, the argument of
// eth_sendProofWrapper. Transactions are encoded in their network form.
func (w *ProofWrapper) MarshalBinary() ([]byte, error) {
	wire := rlpProofWrapper{
		Transactions: make([][]byte, 0, len(w.Transactions)),
		Mode:         w.Mode,
	}

	for i, tx := range w.Transactions {
		encoded, err := tx.MarshalNetwork()
		if err != nil {
			return nil, fmt.Errorf("failed encoding transaction %d: %w", i, err)
		}

		wire.Transactions = append(wire.Transactions, encoded)
	}

	deps := make([][]byte, 0, len(w.Dependencies))
	for _, dep := range w.Dependencies {
		deps = append(deps, EncodeDependencies([]FrameDependency{dep}))
	}

	var content any

	switch w.Mode {
	case ProofWrapperModeDirect:
		proofs := w.Proofs
		if proofs == nil {
			proofs = [][]byte{}
		}

		content = &rlpDirectContent{Dependencies: deps, Proofs: proofs}
	case ProofWrapperModeRecursive:
		content = &rlpRecursiveContent{Dependencies: deps, StarkProof: w.StarkProof}
	default:
		return nil, fmt.Errorf("%w: unknown mode %d", ErrInvalidProofWrapper, w.Mode)
	}

	var err error

	if wire.Content, err = rlp.EncodeToBytes(content); err != nil {
		return nil, err
	}

	encoded, err := rlp.EncodeToBytes(&wire)
	if err != nil {
		return nil, err
	}

	if len(encoded) > MaxWrapperBytes {
		return nil, fmt.Errorf("%w: %d bytes, limit is %d", ErrInvalidProofWrapper, len(encoded), MaxWrapperBytes)
	}

	return encoded, nil
}

// UnmarshalBinary decodes a wrapper from its RLP encoding. It checks the encoding only;
// Validate checks the wrapper against its transactions.
func (w *ProofWrapper) UnmarshalBinary(data []byte) error {
	if len(data) > MaxWrapperBytes {
		return fmt.Errorf("%w: %d bytes, limit is %d", ErrInvalidProofWrapper, len(data), MaxWrapperBytes)
	}

	var wire rlpProofWrapper
	if err := rlp.DecodeBytes(data, &wire); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidProofWrapper, err)
	}

	decoded := ProofWrapper{
		Transactions: make([]*Transaction, 0, len(wire.Transactions)),
		Mode:         wire.Mode,
	}

	for i, raw := range wire.Transactions {
		tx, err := DecodeTx(raw)
		if err != nil {
			return fmt.Errorf("%w: transaction %d: %v", ErrInvalidProofWrapper, i, err)
		}

		decoded.Transactions = append(decoded.Transactions, tx)
	}

	var deps [][]byte

	switch wire.Mode {
	case ProofWrapperModeDirect:
		var content rlpDirectContent
		if err := rlp.DecodeBytes(wire.Content, &content); err != nil {
			return fmt.Errorf("%w: %v", ErrInvalidProofWrapper, err)
		}

		deps, decoded.Proofs = content.Dependencies, content.Proofs
	case ProofWrapperModeRecursive:
		var content rlpRecursiveContent
		if err := rlp.DecodeBytes(wire.Content, &content); err != nil {
			return fmt.Errorf("%w: %v", ErrInvalidProofWrapper, err)
		}

		deps, decoded.StarkProof = content.Dependencies, content.StarkProof
	default:
		return fmt.Errorf("%w: unknown mode %d", ErrInvalidProofWrapper, wire.Mode)
	}

	decoded.Dependencies = make([]FrameDependency, 0, len(deps))

	for i, triple := range deps {
		if len(triple) != DependencyTripleBytes || !bytes.Equal(triple[:31], make([]byte, 31)) {
			return fmt.Errorf("%w: dependency %d is not a %d-byte triple", ErrInvalidProofWrapper, i, DependencyTripleBytes)
		}

		dep := (&Frame{Mode: FrameModeDepVerify, Data: triple}).Dependencies()[0]
		if !dep.Scheme.Known() {
			return fmt.Errorf("%w: dependency %d has unknown scheme 0x%02x", ErrInvalidProofWrapper, i, uint8(dep.Scheme))
		}

		decoded.Dependencies = append(decoded.Dependencies, dep)
	}

	*w = decoded

	return nil
}
