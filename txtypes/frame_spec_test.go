package txtypes

import (
	"crypto/ecdsa"
	"math"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
)

// Tests for the EIP texts as of 2026-10-18: EIP-8250's empty legacy key set and
// single-use form, EIP-7906's ASSERT batch rules, EIP-8272's verification gas, and the
// static validation issues found while aligning with them.

// setLegacyNonceKeyZero selects the EIP-8250 revision for the test's duration. The
// key-zero revision only exists alongside the per-gas fees, so it selects those too.
func setLegacyNonceKeyZero(t *testing.T, value bool) {
	t.Helper()

	previous, previousFees := LegacyNonceKeyZero, LegacyFrameFees
	LegacyNonceKeyZero, LegacyFrameFees = value, value

	t.Cleanup(func() { LegacyNonceKeyZero, LegacyFrameFees = previous, previousFees })
}

// signedKeyedTx builds a self-relayed transaction with the given key set and sequence,
// signed by key.
func signedKeyedTx(t *testing.T, key *ecdsa.PrivateKey, keys []*uint256.Int, seq uint64) *FrameTx {
	t.Helper()

	sender := crypto.PubkeyToAddress(key.PublicKey)

	tx := NewFrameTx(testChainID, sender, 0, FrameFees{GasFeeCap: uint256.NewInt(1e9)},
		[]*Frame{
			SelfVerifyFrame(FrameLimits{Execution: 5_000}),
			UserOpFrame(&testTarget, nil, nil, FrameLimits{Execution: 30_000}),
		},
		[]*FrameSignature{SenderSignature()})

	tx.NonceKeys = keys
	tx.NonceSeq = seq

	if err := tx.SignPayload(testChainID.ToBig(), key); err != nil {
		t.Fatalf("failed signing: %v", err)
	}

	return tx
}

func keys(values ...uint64) []*uint256.Int {
	out := make([]*uint256.Int, 0, len(values))
	for _, v := range values {
		out = append(out, uint256.NewInt(v))
	}

	return out
}

// TestFrameNonceForms checks EIP-8250's key set rules in both revisions: which forms
// are valid, which select the account nonce, and what each one pays.
func TestFrameNonceForms(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed generating key: %v", err)
	}

	sixteen := make([]*uint256.Int, 0, MaxNonceKeys)
	for i := 0; i < MaxNonceKeys; i++ {
		sixteen = append(sixteen, uint256.NewInt(uint64(i)))
	}

	seventeen := append(append([]*uint256.Int{}, sixteen...), uint256.NewInt(MaxNonceKeys))

	cases := []struct {
		name       string
		keyZero    bool
		keys       []*uint256.Int
		seq        uint64
		valid      bool
		legacy     bool
		singleUse  bool
		accessCost uint64
	}{
		{"empty set is the account nonce", false, nil, 3, true, true, false, 0},
		{"key zero is an ordinary key", false, keys(0), 3, true, false, false, 12_100},
		{"key zero combines with others", false, keys(0, 5), 3, true, false, false, 2 * 12_100},
		{"sixteen keys", false, sixteen, 0, true, false, false, 16 * 12_100},
		{"seventeen keys", false, seventeen, 0, false, false, false, 17 * 12_100},
		{"unsorted keys", false, keys(9, 7), 0, false, false, false, 2 * 12_100},
		{"duplicate keys", false, keys(7, 7), 0, false, false, false, 2 * 12_100},
		{"single use", false, keys(7), MaxNonceSeq, true, false, true, 12_100},
		{"single use without keys", false, nil, MaxNonceSeq, false, true, false, 0},

		{"key-zero revision: [0] is the account nonce", true, keys(0), 3, true, true, false, 0},
		{"key-zero revision: empty set is invalid", true, nil, 3, false, false, false, 0},
		{"key-zero revision: zero must be alone", true, keys(0, 5), 3, false, false, false, 2 * 12_100},
		{"key-zero revision: no single use", true, keys(7), MaxNonceSeq, false, false, false, 12_100},
		{"key-zero revision: ordinary key", true, keys(7), 3, true, false, false, 12_100},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			setLegacyNonceKeyZero(t, tc.keyZero)

			tx := signedKeyedTx(t, key, tc.keys, tc.seq)

			err := tx.ValidatePayload()
			if tc.valid && err != nil {
				t.Fatalf("rejected: %v", err)
			}

			if !tc.valid && err == nil {
				t.Fatal("accepted")
			}

			if got := tx.UsesLegacyNonce(); got != tc.legacy {
				t.Errorf("UsesLegacyNonce = %v, want %v", got, tc.legacy)
			}

			if got := tx.UsesSingleUseNonce(); got != tc.singleUse {
				t.Errorf("UsesSingleUseNonce = %v, want %v", got, tc.singleUse)
			}

			if got := tx.keyedNonceAccessCost(); got != tc.accessCost {
				t.Errorf("access cost = %d, want %d", got, tc.accessCost)
			}
		})
	}

	setLegacyNonceKeyZero(t, false)

	single := signedKeyedTx(t, key, nil, 0).WithSingleUseNonceKeys(keys(7, 9))
	if !single.UsesSingleUseNonce() || single.NonceSeq != MaxNonceSeq {
		t.Fatal("WithSingleUseNonceKeys did not select the single-use form")
	}
}

// TestFrameLegacyNonceEncoding checks that the two spellings of the account nonce encode
// differently, decode to what was sent, and that the builders pick the right one.
func TestFrameLegacyNonceEncoding(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed generating key: %v", err)
	}

	sender := crypto.PubkeyToAddress(key.PublicKey)

	build := func(keyZero bool) *FrameTx {
		setLegacyNonceKeyZero(t, keyZero)

		tx := NewFrameTx(testChainID, sender, 7, FrameFees{GasFeeCap: uint256.NewInt(1e9)},
			[]*Frame{SelfVerifyFrame(FrameLimits{Execution: 5_000})},
			[]*FrameSignature{SenderSignature()})

		if err := tx.SignPayload(testChainID.ToBig(), key); err != nil {
			t.Fatalf("failed signing: %v", err)
		}

		return tx
	}

	current := build(false)
	if len(current.NonceKeys) != 0 || !current.UsesLegacyNonce() {
		t.Fatalf("current revision built keys %v", current.NonceKeys)
	}

	keyZero := build(true)
	if len(keyZero.NonceKeys) != 1 || !keyZero.NonceKeys[0].IsZero() || !keyZero.UsesLegacyNonce() {
		t.Fatalf("key-zero revision built keys %v", keyZero.NonceKeys)
	}

	// The same account nonce in the two spellings is two different transactions.
	if NewTx(current).Hash() == NewTx(keyZero).Hash() {
		t.Fatal("[] and [0] must not hash alike")
	}

	for name, tx := range map[string]*FrameTx{"empty": current, "key zero": keyZero} {
		encoded, err := NewTx(tx).MarshalBinary()
		if err != nil {
			t.Fatalf("%s: failed encoding: %v", name, err)
		}

		decoded, err := DecodeTx(encoded)
		if err != nil {
			t.Fatalf("%s: failed decoding: %v", name, err)
		}

		back, _ := decoded.Inner().(*FrameTx)
		if len(back.NonceKeys) != len(tx.NonceKeys) {
			t.Fatalf("%s: decoded %d keys, sent %d", name, len(back.NonceKeys), len(tx.NonceKeys))
		}

		if decoded.Hash() != NewTx(tx).Hash() {
			t.Fatalf("%s: hash changed across encoding", name)
		}

		// JSON keeps the key set as well.
		asJSON, err := NewTx(tx).MarshalJSON()
		if err != nil {
			t.Fatalf("%s: failed JSON encoding: %v", name, err)
		}

		fromJSON, err := UnmarshalJSONTx(asJSON)
		if err != nil {
			t.Fatalf("%s: failed JSON decoding: %v", name, err)
		}

		if fromJSON.Hash() != NewTx(tx).Hash() {
			t.Fatalf("%s: hash changed across JSON", name)
		}
	}

	// WithLegacyNonce spells the account nonce for the selected revision.
	setLegacyNonceKeyZero(t, true)

	relabelled := build(true).WithNonceKeys(keys(7), 1).WithLegacyNonce()
	if len(relabelled.NonceKeys) != 1 || !relabelled.UsesLegacyNonce() {
		t.Fatalf("WithLegacyNonce under the key-zero revision gave %v", relabelled.NonceKeys)
	}

	setLegacyNonceKeyZero(t, false)

	relabelled.WithLegacyNonce()
	if len(relabelled.NonceKeys) != 0 || !relabelled.UsesLegacyNonce() {
		t.Fatalf("WithLegacyNonce under the current revision gave %v", relabelled.NonceKeys)
	}
}

// TestFrameNonceFloorTokensAreFlat checks that the nonce encoding is priced with weighted
// tokens in the intrinsic cost but flat ones in the calldata floor, as EIP-8250 does.
func TestFrameNonceFloorTokensAreFlat(t *testing.T) {
	setLegacyNonceKeyZero(t, false)

	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed generating key: %v", err)
	}

	legacy := signedKeyedTx(t, key, nil, 0)

	// rlp([0x0100]) || rlp(0) = c3 82 01 00 80: one zero byte among five.
	keyed := signedKeyedTx(t, key, keys(0x0100), 0)

	// The two signatures differ in their zero bytes, which the weighted intrinsic cost
	// prices; only the nonce encoding is under test here.
	keyed.Signatures[0].Signature = legacy.Signatures[0].Signature

	legacyExt, keyedExt := legacy.extensionCalldata(), keyed.extensionCalldata()
	if len(legacyExt) != 2 || len(keyedExt) != 5 {
		t.Fatalf("extension calldata lengths %d and %d", len(legacyExt), len(keyedExt))
	}

	wantIntrinsic := legacy.IntrinsicGas() + StandardTokenCost*(17-8) + KeyedNonceAccessCost
	if got := keyed.IntrinsicGas(); got != wantIntrinsic {
		t.Fatalf("intrinsic gas = %d, want %d", got, wantIntrinsic)
	}

	wantFloor := legacy.CalldataFloorGas() + TotalCostFloorPerToken*StandardTokenCost*(5-2) + KeyedNonceAccessCost
	if got := keyed.CalldataFloorGas(); got != wantFloor {
		t.Fatalf("calldata floor gas = %d, want %d", got, wantFloor)
	}
}

// TestFrameGasOverflowRejected checks that frame limits overflowing 64 bits are
// rejected and that the gas helpers saturate rather than report a small wrapped total.
func TestFrameGasOverflowRejected(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed generating key: %v", err)
	}

	sender := crypto.PubkeyToAddress(key.PublicKey)

	build := func(frames ...*Frame) *FrameTx {
		tx := NewFrameTx(testChainID, sender, 0, FrameFees{GasFeeCap: uint256.NewInt(1)},
			frames, []*FrameSignature{SenderSignature()})

		if err := tx.SignPayload(testChainID.ToBig(), key); err != nil {
			t.Fatalf("failed signing: %v", err)
		}

		return tx
	}

	for name, tx := range map[string]*FrameTx{
		"one frame wraps": build(
			SelfVerifyFrame(FrameLimits{Execution: math.MaxUint64, State: 2}),
			UserOpFrame(&testTarget, nil, nil, FrameLimits{Execution: 21_000})),
		"the total wraps": build(
			SelfVerifyFrame(FrameLimits{Execution: 1 << 63}),
			UserOpFrame(&testTarget, nil, nil, FrameLimits{Execution: 1 << 63})),
		"intrinsic plus execution wraps": build(
			SelfVerifyFrame(FrameLimits{Execution: math.MaxUint64 - 100}),
			UserOpFrame(&testTarget, nil, nil, FrameLimits{Execution: 0})),
	} {
		t.Run(name, func(t *testing.T) {
			if err := tx.ValidatePayload(); err == nil {
				t.Fatal("overflowing frame gas was accepted")
			}

			if tx.ExecutionGas() <= TxMaxGasLimit {
				t.Fatalf("ExecutionGas wrapped to %d", tx.ExecutionGas())
			}

			if tx.MaxGas() <= TxMaxGasLimit {
				t.Fatalf("MaxGas wrapped to %d", tx.MaxGas())
			}
		})
	}
}

// TestFrameAssertBatchRules checks EIP-7906's rule that the atomic batch flag is valid
// neither on an ASSERT frame nor on the frame directly before one.
func TestFrameAssertBatchRules(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed generating key: %v", err)
	}

	sender := crypto.PubkeyToAddress(key.PublicKey)
	assertions := common.HexToAddress("0xeeee000000000000000000000000000000000006")

	build := func(frames ...*Frame) *FrameTx {
		tx := NewFrameTx(testChainID, sender, 0, FrameFees{GasFeeCap: uint256.NewInt(1e9)},
			frames, []*FrameSignature{SenderSignature()})

		if err := tx.SignPayload(testChainID.ToBig(), key); err != nil {
			t.Fatalf("failed signing: %v", err)
		}

		return tx
	}

	verify := func() *Frame { return SelfVerifyFrame(FrameLimits{Execution: 5_000}) }
	userOp := func() *Frame { return UserOpFrame(&testTarget, nil, nil, FrameLimits{Execution: 30_000}) }
	assert := func() *Frame { return AssertFrame(assertions, nil, FrameLimits{Execution: 20_000}) }

	for name, tc := range map[string]struct {
		tx    *FrameTx
		valid bool
	}{
		"batch closed before the suffix": {build(verify(), userOp().WithAtomicBatch(), userOp(), assert()), true},
		"batch into an assert frame":     {build(verify(), userOp().WithAtomicBatch(), assert(), assert()), false},
		"batched assert frame":           {build(verify(), userOp(), assert().WithAtomicBatch(), assert()), false},
	} {
		t.Run(name, func(t *testing.T) {
			err := tc.tx.ValidatePayload()
			if tc.valid && err != nil {
				t.Fatalf("rejected: %v", err)
			}

			if !tc.valid && err == nil {
				t.Fatal("accepted")
			}
		})
	}
}

// TestFrameP256Verification checks that VerifySignatures verifies P256 entries, not just
// their encoding.
func TestFrameP256Verification(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed generating key: %v", err)
	}

	p256Key, err := generateP256Key()
	if err != nil {
		t.Fatalf("failed generating P256 key: %v", err)
	}

	sender := crypto.PubkeyToAddress(key.PublicKey)

	tx := NewFrameTx(testChainID, sender, 0, FrameFees{GasFeeCap: uint256.NewInt(1e9)},
		[]*Frame{SelfVerifyFrame(FrameLimits{Execution: 5_000})},
		[]*FrameSignature{SenderSignature(), {Scheme: SigSchemeP256}})

	if err := tx.SignEntryP256(1, p256Key); err != nil {
		t.Fatalf("failed signing P256 entry: %v", err)
	}

	if err := tx.SignPayload(testChainID.ToBig(), key); err != nil {
		t.Fatalf("failed signing: %v", err)
	}

	if err := tx.VerifySignatures(); err != nil {
		t.Fatalf("valid signatures rejected: %v", err)
	}

	// A flipped signature byte keeps the key and signer consistent but no longer
	// verifies, which only a real P256VERIFY notices.
	tampered := tx.CopyTx().(*FrameTx)
	tampered.Signatures[1].Signature[5] ^= 0x01

	if err := tampered.VerifySignatures(); err == nil {
		t.Fatal("a tampered P256 signature verified")
	}

	// A point off the curve, with a signer derived from it so the format check passes.
	offCurve := tx.CopyTx().(*FrameTx)
	for i := 64; i < 128; i++ {
		offCurve.Signatures[1].Signature[i] = 0
	}

	offCurve.Signatures[1].Signature[95] = 1
	offCurve.Signatures[1].Signature[127] = 1
	offCurve.Signatures[1].Signer = P256Signer(uint256.NewInt(1).ToBig(), uint256.NewInt(1).ToBig()).Bytes()

	if err := offCurve.VerifySignatures(); err == nil {
		t.Fatal("a public key off the curve verified")
	}
}

// TestRecentRootVerifyGasFitsCap checks that a full recent root frame still leaves room
// for the account validation inside MaxVerifyGas now that it counts.
func TestRecentRootVerifyGasFitsCap(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed generating key: %v", err)
	}

	sender := crypto.PubkeyToAddress(key.PublicKey)

	references := make([]*RecentRootReference, 0, MaxRecentRootReferences)
	for i := 0; i < MaxRecentRootReferences; i++ {
		references = append(references, &RecentRootReference{Slot: uint64(i + 1)})
	}

	// Sixteen cold tuples cost 3119 + 16*2396 = 41455 by the bytecode.
	budget := RecentRootVerifyGas(MaxRecentRootReferences)
	if budget < 41_455 || budget > 45_000 {
		t.Fatalf("budget for sixteen tuples = %d", budget)
	}

	tx := NewFrameTx(testChainID, sender, 0, FrameFees{GasFeeCap: uint256.NewInt(1e9)},
		[]*Frame{
			RecentRootVerifyFrame(references, budget),
			SelfVerifyFrame(FrameLimits{Execution: 5_000}),
			UserOpFrame(&testTarget, nil, nil, FrameLimits{Execution: 30_000}),
		},
		[]*FrameSignature{SenderSignature()})

	if err := tx.SignPayload(testChainID.ToBig(), key); err != nil {
		t.Fatalf("failed signing: %v", err)
	}

	if err := tx.ValidateMempoolPrefix(); err != nil {
		t.Fatalf("a full recent root frame does not fit the cap: %v", err)
	}
}
