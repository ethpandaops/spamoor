package txtypes

import (
	"crypto/ecdsa"
	"math"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
)

// Tests for the EIP text as of 2026-10-18: EIP-8250's empty legacy key set and
// single-use form, EIP-7906's ASSERT batch rules, EIP-8272's verification gas, and the
// static validation and signing issues found while aligning with them.

// signedKeyedTx builds a self-relayed transaction in the given envelope revision with the
// given key set and sequence, signed by key.
func signedKeyedTx(t *testing.T, key *ecdsa.PrivateKey, extensions FrameExtensions, keys []*uint256.Int, seq uint64) *FrameTx {
	t.Helper()

	sender := crypto.PubkeyToAddress(key.PublicKey)

	tx := NewFrameTxWithExtensions(extensions, testChainID, sender, 0,
		FrameFees{GasFeeCap: uint256.NewInt(1e9)},
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

	const keyZero = FrameExtKeyedNonces | FrameExtLegacyNonceKeyZero

	cases := []struct {
		name       string
		extensions FrameExtensions
		keys       []*uint256.Int
		seq        uint64
		valid      bool
		legacy     bool
		singleUse  bool
		accessCost uint64
	}{
		{"empty set is the account nonce", FrameExtKeyedNonces, nil, 3, true, true, false, 0},
		{"key zero is an ordinary key", FrameExtKeyedNonces, keys(0), 3, true, false, false, 12_100},
		{"key zero combines with others", FrameExtKeyedNonces, keys(0, 5), 3, true, false, false, 2 * 12_100},
		{"sixteen keys", FrameExtKeyedNonces, sixteen, 0, true, false, false, 16 * 12_100},
		{"seventeen keys", FrameExtKeyedNonces, seventeen, 0, false, false, false, 17 * 12_100},
		{"unsorted keys", FrameExtKeyedNonces, keys(9, 7), 0, false, false, false, 2 * 12_100},
		{"duplicate keys", FrameExtKeyedNonces, keys(7, 7), 0, false, false, false, 2 * 12_100},
		{"single use", FrameExtKeyedNonces, keys(7), MaxNonceSeq, true, false, true, 12_100},
		{"single use without keys", FrameExtKeyedNonces, nil, MaxNonceSeq, false, true, false, 0},

		{"key-zero revision: [0] is the account nonce", keyZero, keys(0), 3, true, true, false, 0},
		{"key-zero revision: empty set is invalid", keyZero, nil, 3, false, false, false, 0},
		{"key-zero revision: zero must be alone", keyZero, keys(0, 5), 3, false, false, false, 2 * 12_100},
		{"key-zero revision: no single use", keyZero, keys(7), MaxNonceSeq, false, false, false, 12_100},
		{"key-zero revision: ordinary key", keyZero, keys(7), 3, true, false, false, 12_100},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tx := signedKeyedTx(t, key, tc.extensions, tc.keys, tc.seq)

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

	// WithSingleUseNonceKeys is the single-use form.
	single := signedKeyedTx(t, key, FrameExtKeyedNonces, nil, 0).WithSingleUseNonceKeys(keys(7, 9))
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

	build := func(extensions FrameExtensions) *FrameTx {
		return NewFrameTxWithExtensions(extensions, testChainID, sender, 7,
			FrameFees{GasFeeCap: uint256.NewInt(1e9)},
			[]*Frame{SelfVerifyFrame(FrameLimits{Execution: 5_000})},
			[]*FrameSignature{SenderSignature()})
	}

	current := build(FrameExtKeyedNonces)
	if len(current.NonceKeys) != 0 || !current.UsesLegacyNonce() {
		t.Fatalf("current revision built keys %v", current.NonceKeys)
	}

	keyZero := build(FrameExtKeyedNonces | FrameExtLegacyNonceKeyZero)
	if len(keyZero.NonceKeys) != 1 || !keyZero.NonceKeys[0].IsZero() || !keyZero.UsesLegacyNonce() {
		t.Fatalf("key-zero revision built keys %v", keyZero.NonceKeys)
	}

	if current.Extensions.String() != "8141+8250" || keyZero.Extensions.String() != "8141+8250/key0" {
		t.Fatalf("extension names %q and %q", current.Extensions, keyZero.Extensions)
	}

	// The same account nonce in the two spellings is two different transactions.
	for _, tx := range []*FrameTx{current, keyZero} {
		if err := tx.SignPayload(testChainID.ToBig(), key); err != nil {
			t.Fatalf("failed signing: %v", err)
		}
	}

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

		// The revision marker is not on the wire: a decoded [0] is read as key zero.
		if back.Extensions != FrameExtKeyedNonces {
			t.Fatalf("%s: decoded extensions %s", name, back.Extensions)
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

	// WithLegacyNonce spells the account nonce for the transaction's revision.
	relabelled := build(FrameExtKeyedNonces|FrameExtLegacyNonceKeyZero).WithNonceKeys(keys(7), 1)
	relabelled.WithLegacyNonce()

	if len(relabelled.NonceKeys) != 1 || !relabelled.UsesLegacyNonce() {
		t.Fatalf("WithLegacyNonce under the key-zero revision gave %v", relabelled.NonceKeys)
	}

	relabelled.Extensions = FrameExtKeyedNonces
	relabelled.WithLegacyNonce()

	if len(relabelled.NonceKeys) != 0 || !relabelled.UsesLegacyNonce() {
		t.Fatalf("WithLegacyNonce under the current revision gave %v", relabelled.NonceKeys)
	}
}

// TestFrameNonceFloorTokensAreFlat checks that the nonce encoding is priced with weighted
// tokens in the intrinsic cost but flat ones in the calldata floor, as EIP-8250 does.
func TestFrameNonceFloorTokensAreFlat(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed generating key: %v", err)
	}

	legacy := signedKeyedTx(t, key, FrameExtKeyedNonces, nil, 0)

	// rlp([0x0100]) || rlp(0) = c3 82 01 00 80: one zero byte among five.
	keyed := signedKeyedTx(t, key, FrameExtKeyedNonces, keys(0x0100), 0)

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

// TestFrameGasOverflowRejected reproduces the wrapped-sum validity hole: frame limits
// that overflow 64 bits must be rejected, and the gas helpers must saturate rather than
// report a small wrapped total.
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
// their encoding, and that P256 entries have to be signed before the canonical-hash
// entries whose digest they change.
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

	build := func() *FrameTx {
		return NewFrameTx(testChainID, sender, 0, FrameFees{GasFeeCap: uint256.NewInt(1e9)},
			[]*Frame{SelfVerifyFrame(FrameLimits{Execution: 5_000})},
			[]*FrameSignature{SenderSignature(), {Scheme: SigSchemeP256}})
	}

	tx := build()
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

	// Signing the secp256k1 entry first leaves it over a hash the P256 signer changes.
	late := build()
	if err := late.SignPayload(testChainID.ToBig(), key); err != nil {
		t.Fatalf("failed signing: %v", err)
	}

	if err := late.SignEntryP256(1, p256Key); err == nil {
		t.Fatal("signing a P256 entry after a canonical-hash entry was allowed")
	}
}

// TestFrameSignPayloadOrder checks that one key owning both an explicit-msg entry and a
// canonical-hash entry ends up with both verifying: the explicit one has to be signed
// first because its bytes are part of the canonical hash.
func TestFrameSignPayloadOrder(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed generating key: %v", err)
	}

	sender := crypto.PubkeyToAddress(key.PublicKey)
	digest := common.HexToHash("0x1122334455667788112233445566778811223344556677881122334455667788")

	tx := NewFrameTx(testChainID, sender, 0, FrameFees{GasFeeCap: uint256.NewInt(1e9)},
		[]*Frame{SelfVerifyFrame(FrameLimits{Execution: 5_000})},
		[]*FrameSignature{
			SenderSignature(),
			{Scheme: SigSchemeSecp256k1, Signer: sender.Bytes(), Msg: digest.Bytes()},
		})

	if err := tx.SignPayload(testChainID.ToBig(), key); err != nil {
		t.Fatalf("failed signing: %v", err)
	}

	if err := tx.VerifySignatures(); err != nil {
		t.Fatalf("signatures do not verify: %v", err)
	}
}

// TestSignTxKeepsExplicitSender checks that signing with a co-signer's key does not make
// From report the co-signer.
func TestSignTxKeepsExplicitSender(t *testing.T) {
	senderKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed generating key: %v", err)
	}

	paymasterKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed generating key: %v", err)
	}

	sender := crypto.PubkeyToAddress(senderKey.PublicKey)
	paymaster := crypto.PubkeyToAddress(paymasterKey.PublicKey)

	tx := NewFrameTx(testChainID, sender, 0, FrameFees{GasFeeCap: uint256.NewInt(1e9)},
		[]*Frame{OnlyVerifyFrame(FrameLimits{Execution: 5_000}), PayFrame(paymaster, nil, FrameLimits{Execution: 5_000})},
		[]*FrameSignature{SenderSignature(), SignerSignature(paymaster)})

	signed, err := SignTx(NewTx(tx), testChainID.ToBig(), paymasterKey)
	if err != nil {
		t.Fatalf("failed signing: %v", err)
	}

	from, err := signed.From(testChainID.ToBig())
	if err != nil {
		t.Fatalf("From: %v", err)
	}

	if from != sender {
		t.Fatalf("From = %s, want the sender %s", from, sender)
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

	// Sixteen cold tuples cost 3119 + 16*2396 = 41455 by the bytecode; the budget must
	// cover that with room, without crowding the account validation out of the cap.
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

// TestFrameReceiptRLPRoundTrip checks the consensus receipt encoding.
func TestFrameReceiptRLPRoundTrip(t *testing.T) {
	extra := &FrameReceiptExtra{
		Payer: common.HexToAddress("0xdddd000000000000000000000000000000000004"),
		Frames: []*FrameReceipt{
			{Status: FrameStatusSuccess, ExecutionGas: 2_600, StateGas: 0},
			{Status: FrameStatusSuccess, ExecutionGas: 30_000, StateGas: 1_000, Logs: []*Log{{
				Address: testTarget,
				Topics:  []common.Hash{common.HexToHash("0x01")},
				Data:    []byte{0xaa},
			}}},
			{Status: FrameStatusSkipped},
		},
	}

	encoded, err := EncodeFrameReceiptRLP(0x1234, extra)
	if err != nil {
		t.Fatalf("failed encoding: %v", err)
	}

	if encoded[0] != FrameTxType {
		t.Fatalf("type byte 0x%02x", encoded[0])
	}

	for name, data := range map[string][]byte{"typed": encoded, "payload": encoded[1:]} {
		cumulative, back, err := DecodeFrameReceiptRLP(data)
		if err != nil {
			t.Fatalf("%s: failed decoding: %v", name, err)
		}

		if cumulative != 0x1234 || back.Payer != extra.Payer || len(back.Frames) != 3 {
			t.Fatalf("%s: decoded %d, %s, %d frames", name, cumulative, back.Payer, len(back.Frames))
		}

		if back.Frames[1].ExecutionGas != 30_000 || back.Frames[1].StateGas != 1_000 || len(back.Frames[1].Logs) != 1 {
			t.Fatalf("%s: frame 1 decoded as %+v", name, back.Frames[1])
		}

		if back.Frames[1].Logs[0].Address != testTarget || back.Frames[1].Logs[0].Data[0] != 0xaa {
			t.Fatalf("%s: log did not survive", name)
		}

		if !back.Frames[2].Skipped() {
			t.Fatalf("%s: frame 2 should be skipped", name)
		}
	}
}
