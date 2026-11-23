package txtypes

import (
	"bytes"
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/holiman/uint256"
)

// wrappedTx builds a signed self-relayed transaction declaring one leanSPHINCS
// dependency over its own signature hash under key, and returns it with that
// dependency in its resolved form.
func wrappedTx(t *testing.T, nonce uint64, key common.Hash) (*Transaction, FrameDependency) {
	t.Helper()

	signer, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed generating key: %v", err)
	}

	frameTx := NewFrameTx(testChainID, crypto.PubkeyToAddress(signer.PublicKey), nonce,
		FrameFees{GasFeeCap: uint256.NewInt(1e9)},
		[]*Frame{
			SelfVerifyFrame(FrameLimits{Execution: 5_000}),
			DependencyFrame(FrameDependency{Scheme: DepSchemeLeanSphincs, VerificationKey: key}),
			UserOpFrame(&testTarget, nil, nil, FrameLimits{Execution: 30_000}),
		},
		[]*FrameSignature{SenderSignature()})

	if err := frameTx.SignPayload(testChainID.ToBig(), signer); err != nil {
		t.Fatalf("failed signing: %v", err)
	}

	deps, err := frameTx.Dependencies()
	if err != nil || len(deps) != 1 {
		t.Fatalf("dependencies = %+v, %v", deps, err)
	}

	return NewTx(frameTx), deps[0]
}

// TestDirectProofWrapperEncoding checks the wire form of a direct wrapper against the
// layout the devnet clients decode: [[raw_tx...], 0, [[dependency...], [proof...]]],
// with transactions in ascending hash order and the resolved dependencies.
func TestDirectProofWrapperEncoding(t *testing.T) {
	txA, depA := wrappedTx(t, 0, common.HexToHash("0x0a"))
	txB, depB := wrappedTx(t, 1, common.HexToHash("0x0b"))

	proofs := map[FrameDependency][]byte{
		depA: bytes.Repeat([]byte{0xaa}, 6208),
		depB: bytes.Repeat([]byte{0xbb}, 6208),
	}

	wrapper, err := NewDirectProofWrapper([]*Transaction{txA, txB}, proofs)
	if err != nil {
		t.Fatalf("failed building wrapper: %v", err)
	}

	if err := wrapper.Validate(); err != nil {
		t.Fatalf("wrapper rejected: %v", err)
	}

	encoded, err := wrapper.MarshalBinary()
	if err != nil {
		t.Fatalf("failed encoding: %v", err)
	}

	var wire struct {
		Transactions [][]byte
		Mode         uint64
		Content      struct {
			Deps   [][]byte
			Proofs [][]byte
		}
	}

	if err := rlp.DecodeBytes(encoded, &wire); err != nil {
		t.Fatalf("wire form does not decode: %v", err)
	}

	if wire.Mode != ProofWrapperModeDirect || len(wire.Transactions) != 2 {
		t.Fatalf("mode %d with %d transactions", wire.Mode, len(wire.Transactions))
	}

	first, second := txA, txB
	if hashA, hashB := txA.Hash(), txB.Hash(); bytes.Compare(hashA[:], hashB[:]) > 0 {
		first, second = txB, txA
	}

	for i, tx := range []*Transaction{first, second} {
		raw, err := tx.MarshalNetwork()
		if err != nil {
			t.Fatalf("failed encoding transaction: %v", err)
		}

		if !bytes.Equal(wire.Transactions[i], raw) {
			t.Errorf("transaction %d is not %s in its network form", i, tx.Hash())
		}
	}

	wantDeps := CanonicalDependencies([]FrameDependency{depA, depB})

	for i, dep := range wantDeps {
		if !bytes.Equal(wire.Content.Deps[i], EncodeDependencies([]FrameDependency{dep})) {
			t.Errorf("dependency %d is not the resolved triple", i)
		}

		if !bytes.Equal(wire.Content.Proofs[i], proofs[dep]) {
			t.Errorf("proof %d does not belong to dependency %d", i, i)
		}
	}

	var decoded ProofWrapper
	if err := decoded.UnmarshalBinary(encoded); err != nil {
		t.Fatalf("failed decoding: %v", err)
	}

	if err := decoded.Validate(); err != nil {
		t.Fatalf("decoded wrapper rejected: %v", err)
	}

	reencoded, err := decoded.MarshalBinary()
	if err != nil || !bytes.Equal(reencoded, encoded) {
		t.Fatalf("round trip changed the encoding (%v)", err)
	}
}

// TestRecursiveProofWrapperEncoding checks the wire form of a recursive wrapper, whose
// content is the dependency list followed by the proof itself.
func TestRecursiveProofWrapperEncoding(t *testing.T) {
	tx, dep := wrappedTx(t, 0, common.HexToHash("0x0a"))
	stark := bytes.Repeat([]byte{0x17}, 1024)

	wrapper, err := NewRecursiveProofWrapper([]*Transaction{tx}, stark)
	if err != nil {
		t.Fatalf("failed building wrapper: %v", err)
	}

	encoded, err := wrapper.MarshalBinary()
	if err != nil {
		t.Fatalf("failed encoding: %v", err)
	}

	var wire struct {
		Transactions [][]byte
		Mode         uint64
		Content      struct {
			Deps  [][]byte
			Proof []byte
		}
	}

	if err := rlp.DecodeBytes(encoded, &wire); err != nil {
		t.Fatalf("wire form does not decode: %v", err)
	}

	if wire.Mode != ProofWrapperModeRecursive || !bytes.Equal(wire.Content.Proof, stark) {
		t.Fatalf("mode %d, %d proof bytes", wire.Mode, len(wire.Content.Proof))
	}

	var decoded ProofWrapper
	if err := decoded.UnmarshalBinary(encoded); err != nil {
		t.Fatalf("failed decoding: %v", err)
	}

	if err := decoded.Validate(); err != nil {
		t.Fatalf("decoded wrapper rejected: %v", err)
	}

	if len(decoded.Dependencies) != 1 || decoded.Dependencies[0] != dep {
		t.Fatalf("dependencies = %+v, want %+v", decoded.Dependencies, dep)
	}
}

// TestProofWrapperRejections checks what a wrapper may not be.
func TestProofWrapperRejections(t *testing.T) {
	txA, depA := wrappedTx(t, 0, common.HexToHash("0x0a"))
	txB, depB := wrappedTx(t, 1, common.HexToHash("0x0b"))
	proofs := map[FrameDependency][]byte{depA: {1}, depB: {2}}

	if _, err := NewDirectProofWrapper([]*Transaction{txA, txB}, map[FrameDependency][]byte{depA: {1}}); !errors.Is(err, ErrInvalidProofWrapper) {
		t.Errorf("missing proof: got %v", err)
	}

	plain := NewTx(&DynamicFeeTx{ChainID: testChainID.ToBig()})
	if _, err := NewDirectProofWrapper([]*Transaction{plain}, nil); !errors.Is(err, ErrInvalidProofWrapper) {
		t.Errorf("non-frame transaction: got %v", err)
	}

	build := func() *ProofWrapper {
		wrapper, err := NewDirectProofWrapper([]*Transaction{txA, txB}, proofs)
		if err != nil {
			t.Fatalf("failed building wrapper: %v", err)
		}

		return wrapper
	}

	for name, mutate := range map[string]func(*ProofWrapper){
		"no transactions":       func(w *ProofWrapper) { w.Transactions = nil },
		"descending hash order": func(w *ProofWrapper) { w.Transactions[0], w.Transactions[1] = w.Transactions[1], w.Transactions[0] },
		"repeated transaction":  func(w *ProofWrapper) { w.Transactions[1] = w.Transactions[0] },
		"missing dependency":    func(w *ProofWrapper) { w.Dependencies = w.Dependencies[:1] },
		"foreign dependency":    func(w *ProofWrapper) { w.Dependencies[0].VerificationKey = common.HexToHash("0xff") },
		"missing proof":         func(w *ProofWrapper) { w.Proofs = w.Proofs[:1] },
		"unknown mode":          func(w *ProofWrapper) { w.Mode = 2 },
		"recursive without proof": func(w *ProofWrapper) {
			w.Mode, w.Proofs = ProofWrapperModeRecursive, nil
		},
	} {
		t.Run(name, func(t *testing.T) {
			wrapper := build()
			mutate(wrapper)

			if err := wrapper.Validate(); !errors.Is(err, ErrInvalidProofWrapper) {
				t.Fatalf("got %v, want %v", err, ErrInvalidProofWrapper)
			}
		})
	}

	// A direct wrapper holds a bounded number of dependencies.
	txs := make([]*Transaction, 0, MaxLeanSigDepsPerWrapper+1)
	many := make(map[FrameDependency][]byte, MaxLeanSigDepsPerWrapper+1)

	for i := 0; i <= MaxLeanSigDepsPerWrapper; i++ {
		tx, dep := wrappedTx(t, uint64(i), common.BytesToHash([]byte{byte(i + 1)}))
		txs = append(txs, tx)
		many[dep] = []byte{byte(i)}
	}

	wrapper, err := NewDirectProofWrapper(txs, many)
	if err != nil {
		t.Fatalf("failed building wrapper: %v", err)
	}

	if err := wrapper.Validate(); !errors.Is(err, ErrInvalidProofWrapper) {
		t.Errorf("over the direct dependency limit: got %v", err)
	}

	// Malformed encodings.
	encoded, err := build().MarshalBinary()
	if err != nil {
		t.Fatalf("failed encoding: %v", err)
	}

	var decoded ProofWrapper
	if err := decoded.UnmarshalBinary(encoded[:len(encoded)-1]); !errors.Is(err, ErrInvalidProofWrapper) {
		t.Errorf("truncated wrapper: got %v", err)
	}

	if err := decoded.UnmarshalBinary(make([]byte, MaxWrapperBytes+1)); !errors.Is(err, ErrInvalidProofWrapper) {
		t.Errorf("oversized wrapper: got %v", err)
	}

	oversized := build()
	oversized.Proofs[0] = make([]byte, MaxWrapperBytes)

	if _, err := oversized.MarshalBinary(); !errors.Is(err, ErrInvalidProofWrapper) {
		t.Errorf("oversized encoding: got %v", err)
	}
}
