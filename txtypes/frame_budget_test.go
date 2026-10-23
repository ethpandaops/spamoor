package txtypes

import (
	"bytes"
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/holiman/uint256"
)

// Tests for the EIP-8141 revision bogota-devnet-2 runs: the max_fee budget
// (EIPs#12414), signature lengths in the signature hash (EIPs#12502) and EIP-8250's
// single-use pricing exemption (EIPs#12534).

// setLegacyFrameFees selects the builders' fee revision for the test's duration.
func setLegacyFrameFees(t *testing.T, value bool) {
	t.Helper()

	previous := LegacyFrameFees
	LegacyFrameFees = value

	t.Cleanup(func() { LegacyFrameFees = previous })
}

// feeListLength returns the number of fields in an encoded payload's fee list.
func feeListLength(t *testing.T, tx *FrameTx) int {
	t.Helper()

	var buf bytes.Buffer
	if err := tx.EncodePayload(&buf); err != nil {
		t.Fatalf("failed encoding: %v", err)
	}

	var env frameEnvelopeKeyed
	if err := rlp.DecodeBytes(buf.Bytes(), &env); err != nil {
		t.Fatalf("failed decoding envelope: %v", err)
	}

	content, _, err := rlp.SplitList(env.Fees)
	if err != nil {
		t.Fatalf("fee field is not a list: %v", err)
	}

	count, err := rlp.CountValues(content)
	if err != nil {
		t.Fatalf("failed counting fee fields: %v", err)
	}

	return count
}

// TestFrameFeeRevisionsRoundTrip checks that both fee layouts encode as their EIP text
// says, and that decoding, by binary and by JSON, recovers the revision and the hash.
func TestFrameFeeRevisionsRoundTrip(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed generating key: %v", err)
	}

	sender := crypto.PubkeyToAddress(key.PublicKey)

	for _, legacy := range []bool{false, true} {
		setLegacyFrameFees(t, legacy)

		tx := buildTransferTx(t, sender)
		if err := tx.SignPayload(testChainID.ToBig(), key); err != nil {
			t.Fatalf("legacy=%v: failed signing: %v", legacy, err)
		}

		if want := map[bool]int{false: 2, true: 3}[legacy]; feeListLength(t, tx) != want {
			t.Fatalf("legacy=%v: fee list has %d fields, want %d", legacy, feeListLength(t, tx), want)
		}

		if !legacy && (tx.Fees.MaxFee == nil || tx.Fees.MaxFee.IsZero()) {
			t.Fatal("signing did not derive max_fee")
		}

		if err := tx.ValidatePayload(); err != nil {
			t.Fatalf("legacy=%v: rejected: %v", legacy, err)
		}

		sent := NewTx(tx)

		encoded, err := sent.MarshalBinary()
		if err != nil {
			t.Fatalf("legacy=%v: failed encoding: %v", legacy, err)
		}

		decoded, err := DecodeTx(encoded)
		if err != nil {
			t.Fatalf("legacy=%v: failed decoding: %v", legacy, err)
		}

		back, _ := decoded.Inner().(*FrameTx)
		if back.LegacyFees != legacy || decoded.Hash() != sent.Hash() || back.SigHash() != tx.SigHash() {
			t.Fatalf("legacy=%v: decoded LegacyFees=%v, hash or sig hash changed", legacy, back.LegacyFees)
		}

		asJSON, err := sent.MarshalJSON()
		if err != nil {
			t.Fatalf("legacy=%v: failed JSON encoding: %v", legacy, err)
		}

		if bytes.Contains(asJSON, []byte(`"maxFee"`)) == legacy {
			t.Fatalf("legacy=%v: JSON is %s", legacy, asJSON)
		}

		fromJSON, err := UnmarshalJSONTx(asJSON)
		if err != nil {
			t.Fatalf("legacy=%v: failed JSON decoding: %v", legacy, err)
		}

		if fromJSON.Hash() != sent.Hash() {
			t.Fatalf("legacy=%v: hash changed across JSON", legacy)
		}
	}
}

// TestFrameSigHashCommitsToLength checks that an elided signature is replaced by its
// length: rewriting its bytes keeps the hash, changing its length does not. Under
// LegacyFees neither changes it.
func TestFrameSigHashCommitsToLength(t *testing.T) {
	sender := crypto.PubkeyToAddress(mustKey(t).PublicKey)

	for _, legacy := range []bool{false, true} {
		setLegacyFrameFees(t, legacy)

		tx := buildTransferTx(t, sender)
		tx.Signatures = append(tx.Signatures, ArbitrarySignature(bytes.Repeat([]byte{0x01}, 32)))

		base := tx.SigHash()

		tx.Signatures[1].Signature = bytes.Repeat([]byte{0x02}, 32)
		if tx.SigHash() != base {
			t.Fatalf("legacy=%v: rewriting elided bytes changed the hash", legacy)
		}

		tx.Signatures[1].Signature = bytes.Repeat([]byte{0x02}, 33)
		if changed := tx.SigHash() != base; changed == legacy {
			t.Fatalf("legacy=%v: padding the elided bytes changed the hash: %v", legacy, changed)
		}
	}
}

// TestFrameSigHashMatchesSpec rebuilds compute_sig_hash from the EIP text: each
// empty-msg entry's signature replaced by its length as an RLP integer.
func TestFrameSigHashMatchesSpec(t *testing.T) {
	key := mustKey(t)
	tx := buildTransferTx(t, crypto.PubkeyToAddress(key.PublicKey))
	tx.Signatures = append(tx.Signatures, ArbitrarySignature([]byte{0x00, 0x01, 0x02}))

	if err := tx.SignPayload(testChainID.ToBig(), key); err != nil {
		t.Fatalf("failed signing: %v", err)
	}

	type specSignature struct {
		Scheme    FrameSigScheme
		Signer    []byte
		Msg       []byte
		Signature uint64
	}

	type specFees struct{ MaxFee, MaxPriorityFeePerGas *uint256.Int }

	preimage := []any{
		tx.ChainID, tx.NonceKeys, tx.NonceSeq, tx.Sender, tx.Frames,
		[]*specSignature{{SigSchemeSecp256k1, nil, nil, 65}, {SigSchemeArbitrary, nil, nil, 3}},
		&specFees{tx.Fees.MaxFee, tx.Fees.GasTipCap}, tx.BlobHashes,
	}

	encoded, err := rlp.EncodeToBytes(preimage)
	if err != nil {
		t.Fatalf("failed encoding preimage: %v", err)
	}

	if want := crypto.Keccak256Hash(append([]byte{FrameTxType}, encoded...)); tx.SigHash() != want {
		t.Fatalf("sig hash = %s, want %s", tx.SigHash(), want)
	}
}

// TestFrameElidedSignaturePricing checks that elided signature bytes are priced by
// length, as if all non-zero, while explicit-msg entries keep the weighted price.
func TestFrameElidedSignaturePricing(t *testing.T) {
	sender := crypto.PubkeyToAddress(mustKey(t).PublicKey)
	zeros := make([]byte, 10)

	for _, legacy := range []bool{false, true} {
		setLegacyFrameFees(t, legacy)

		tx := buildTransferTx(t, sender)
		without := tx.IntrinsicGas()

		tx.Signatures = append(tx.Signatures, ArbitrarySignature(zeros))

		want := uint64(sigGasArbitrary) + StandardTokenCost*uint64(len(zeros))
		if !legacy {
			want = sigGasArbitrary + StandardTokenCost*4*uint64(len(zeros))
		}

		if got := tx.IntrinsicGas() - without; got != want {
			t.Fatalf("legacy=%v: ten zero bytes cost %d, want %d", legacy, got, want)
		}

		// An explicit digest commits the bytes, so they keep the weighted price.
		tx.Signatures[1].Msg = make([]byte, 32)
		want = sigGasArbitrary + StandardTokenCost*uint64(32+len(zeros))

		if got := tx.IntrinsicGas() - without; got != want {
			t.Fatalf("legacy=%v: explicit-msg entry cost %d, want %d", legacy, got, want)
		}
	}
}

// TestFrameDerivedMaxFeeCoversSigning checks that the budget derived before signing
// still covers the signed transaction, and that signing does not change it.
func TestFrameDerivedMaxFeeCoversSigning(t *testing.T) {
	key := mustKey(t)
	tx := buildTransferTx(t, crypto.PubkeyToAddress(key.PublicKey))
	tx.BlobHashes = nil

	before := tx.SigHash()

	if err := tx.SignPayload(testChainID.ToBig(), key); err != nil {
		t.Fatalf("failed signing: %v", err)
	}

	if tx.SigHash() != before {
		t.Fatal("signing changed the signature hash")
	}

	required := new(big.Int).Mul(new(big.Int).SetUint64(tx.MaxGas()), tx.Fees.GasFeeCap.ToBig())
	if tx.GetMaxFee().Cmp(required) != 0 {
		t.Fatalf("max fee = %s, want MaxGas at the fee cap %s", tx.GetMaxFee(), required)
	}

	// GetGasFeeCap reports the budget per gas, so the engine's Cost covers max_fee.
	cost := NewTx(tx).Cost()
	cost.Sub(cost, NewTx(tx).Value())

	if cost.Cmp(tx.GetMaxFee()) < 0 {
		t.Fatalf("cost %s below max fee %s", cost, tx.GetMaxFee())
	}
}

// TestFrameMaxCost checks the fee budget's collection rule: the worst-case base fee
// charge plus the tip, capped by what remains of max_fee.
func TestFrameMaxCost(t *testing.T) {
	tx := buildTransferTx(t, testTarget)
	maxGas := new(big.Int).SetUint64(tx.MaxGas())
	gas := func(perGas int64) *big.Int { return new(big.Int).Mul(maxGas, big.NewInt(perGas)) }

	tx.Fees.MaxFee = uint256.MustFromBig(gas(10))
	tx.Fees.GasTipCap = uint256.NewInt(2)

	for _, tc := range []struct {
		baseFee int64
		want    *big.Int
	}{
		{5, gas(7)},  // full tip fits
		{9, gas(10)}, // tip clipped to the remaining budget
		{10, gas(10)},
		{11, nil}, // below the required max fee
	} {
		got := tx.MaxCost(big.NewInt(tc.baseFee), nil)
		if (got == nil) != (tc.want == nil) || (got != nil && got.Cmp(tc.want) != 0) {
			t.Fatalf("base fee %d: max cost %v, want %v", tc.baseFee, got, tc.want)
		}
	}
}

// TestFrameSingleUseNonceSeqUnpriced checks EIPs#12534: the single-use sequence is not
// priced as transaction data, while the keys still are.
func TestFrameSingleUseNonceSeqUnpriced(t *testing.T) {
	setLegacyNonceKeyZero(t, false)

	key := mustKey(t)
	sequential := signedKeyedTx(t, key, keys(7), 0)
	singleUse := signedKeyedTx(t, key, keys(7), MaxNonceSeq)

	if !bytes.Equal(singleUse.extensionCalldata(), []byte{0xc1, 0x07}) {
		t.Fatalf("single-use nonce calldata = %x, want rlp([7])", singleUse.extensionCalldata())
	}

	// rlp(0) = 0x80 is one non-zero byte.
	singleUse.Signatures[0].Signature = sequential.Signatures[0].Signature
	if got, want := singleUse.IntrinsicGas(), sequential.IntrinsicGas()-StandardTokenCost*4; got != want {
		t.Fatalf("single-use intrinsic gas = %d, want %d", got, want)
	}
}

func mustKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()

	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed generating key: %v", err)
	}

	return key
}
