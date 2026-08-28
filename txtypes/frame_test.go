package txtypes

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/holiman/uint256"
)

var (
	testChainID = uint256.NewInt(3151908)
	testTarget  = common.HexToAddress("0x1234567890123456789012345678901234567890")
)

// buildTransferTx builds the canonical self-relayed transfer: a self_verify frame
// followed by a value-bearing user operation.
func buildTransferTx(t *testing.T, sender common.Address) *FrameTx {
	t.Helper()

	return NewFrameTx(testChainID, sender, 7,
		FrameFees{
			GasTipCap: uint256.NewInt(1e9),
			GasFeeCap: uint256.NewInt(20e9),
		},
		[]*Frame{
			SelfVerifyFrame(DefaultCodeVerifyLimits(true)),
			UserOpFrame(&testTarget, uint256.NewInt(1e15), nil, FrameLimits{Execution: 30_000}),
		},
		[]*FrameSignature{SenderSignature()},
	)
}

// TestFrameTxRoundTrip checks that a signed frame transaction survives encoding.
func TestFrameTxRoundTrip(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed generating key: %v", err)
	}

	sender := crypto.PubkeyToAddress(key.PublicKey)

	signed, err := SignTx(NewTx(buildTransferTx(t, sender)), testChainID.ToBig(), key)
	if err != nil {
		t.Fatalf("failed signing frame transaction: %v", err)
	}

	if signed.Type() != FrameTxType {
		t.Fatalf("wrong type: 0x%02x", signed.Type())
	}

	encoded, err := signed.MarshalNetwork()
	if err != nil {
		t.Fatalf("failed encoding: %v", err)
	}

	if encoded[0] != FrameTxType {
		t.Fatalf("wrong type byte: 0x%02x", encoded[0])
	}

	decoded, err := DecodeTx(encoded)
	if err != nil {
		t.Fatalf("failed decoding: %v", err)
	}

	if decoded.Hash() != signed.Hash() {
		t.Fatalf("hash changed across encoding: %s != %s", decoded.Hash(), signed.Hash())
	}

	// The sender is an explicit field, so it needs no recovery.
	from, err := decoded.From(testChainID.ToBig())
	if err != nil {
		t.Fatalf("failed reading sender: %v", err)
	}

	if from != sender {
		t.Fatalf("wrong sender: %s != %s", from, sender)
	}

	inner, ok := decoded.Inner().(*FrameTx)
	if !ok {
		t.Fatalf("decoded to %T", decoded.Inner())
	}

	if len(inner.Frames) != 2 || len(inner.Signatures) != 1 {
		t.Fatal("frame or signature list changed across encoding")
	}

	if inner.Frames[0].Target != nil {
		t.Fatal("nil target did not survive encoding")
	}

	if inner.Frames[1].Target == nil || *inner.Frames[1].Target != testTarget {
		t.Fatal("user op target did not survive encoding")
	}

	if inner.Frames[1].Limits.Execution != 30_000 {
		t.Fatal("frame limits did not survive encoding")
	}

	if decoded.To() == nil || *decoded.To() != testTarget {
		t.Fatal("To should resolve to the first SENDER frame target")
	}

	if decoded.Value().Cmp(big.NewInt(1e15)) != 0 {
		t.Fatalf("Value should sum the frame values, got %s", decoded.Value())
	}
}

// TestFrameSigHashElision checks that entries signing the canonical digest have their
// own signature bytes elided from the preimage, and that entries with an explicit
// digest do not.
func TestFrameSigHashElision(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed generating key: %v", err)
	}

	sender := crypto.PubkeyToAddress(key.PublicKey)
	tx := buildTransferTx(t, sender)

	before := tx.SigHash()

	if err := tx.SignPayload(testChainID.ToBig(), key); err != nil {
		t.Fatalf("failed signing: %v", err)
	}

	if len(tx.Signatures[0].Signature) != 65 {
		t.Fatalf("secp256k1 entry should be 65 bytes, got %d", len(tx.Signatures[0].Signature))
	}

	if after := tx.SigHash(); after != before {
		t.Fatal("signing an empty-msg entry changed the signature hash")
	}

	// An entry with an explicit digest keeps its bytes in the preimage, so filling it
	// must change the hash.
	digest := common.HexToHash("0x1122334455667788112233445566778811223344556677881122334455667788")
	tx.Signatures = append(tx.Signatures, &FrameSignature{
		Scheme: SigSchemeSecp256k1,
		Signer: testTarget.Bytes(),
		Msg:    digest.Bytes(),
	})

	withEmpty := tx.SigHash()

	tx.Signatures[1].Signature = bytes.Repeat([]byte{0xaa}, 65)

	if withFilled := tx.SigHash(); withFilled == withEmpty {
		t.Fatal("an explicit-digest entry's signature bytes should be covered by the hash")
	}
}

// TestFrameSignatureRecovery checks that the entry encoding is v || r || s, which
// differs from go-ethereum's r || s || v.
func TestFrameSignatureRecovery(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed generating key: %v", err)
	}

	sender := crypto.PubkeyToAddress(key.PublicKey)
	tx := buildTransferTx(t, sender)
	digest := tx.SigHash()

	if err := tx.SignPayload(testChainID.ToBig(), key); err != nil {
		t.Fatalf("failed signing: %v", err)
	}

	entry := tx.Signatures[0].Signature

	// Rearrange into go-ethereum's layout to recover with its primitives.
	gethSig := make([]byte, 65)
	copy(gethSig[:64], entry[1:])
	gethSig[64] = entry[0]

	pub, err := crypto.Ecrecover(digest[:], gethSig)
	if err != nil {
		t.Fatalf("failed recovering signature: %v", err)
	}

	var recovered common.Address

	copy(recovered[:], crypto.Keccak256(pub[1:])[12:])

	if recovered != sender {
		t.Fatalf("entry recovered to %s, want %s", recovered, sender)
	}
}

// TestFrameP256Signature checks P256 entry signing and the derived signer address.
func TestFrameP256Signature(t *testing.T) {
	key, err := generateP256Key()
	if err != nil {
		t.Fatalf("failed generating P256 key: %v", err)
	}

	sender := common.HexToAddress("0xaaaa000000000000000000000000000000000001")

	tx := NewFrameTx(testChainID, sender, 0,
		FrameFees{GasFeeCap: uint256.NewInt(1e9)},
		[]*Frame{SelfVerifyFrame(FrameLimits{Execution: 5_000})},
		[]*FrameSignature{SenderSignature(), {Scheme: SigSchemeP256}},
	)

	if err := tx.SignEntryP256(1, key); err != nil {
		t.Fatalf("failed signing P256 entry: %v", err)
	}

	if len(tx.Signatures[1].Signature) != 128 {
		t.Fatalf("P256 signature should be 128 bytes, got %d", len(tx.Signatures[1].Signature))
	}

	signer, ok := tx.Signatures[1].ResolvedSigner(sender)
	if !ok {
		t.Fatal("P256 entry should resolve a signer")
	}

	if signer != P256Signer(key.X, key.Y) {
		t.Fatal("P256 signer address does not match the public key")
	}
}

// TestFrameIntrinsicGas checks the intrinsic and floor gas formulas against values
// computed directly from the EIP's definitions.
func TestFrameIntrinsicGas(t *testing.T) {
	sender := common.HexToAddress("0xaaaa000000000000000000000000000000000001")
	data := []byte{0x00, 0x01, 0x02, 0x00}

	tx := NewFrameTx(testChainID, sender, 0,
		FrameFees{GasFeeCap: uint256.NewInt(1e9)},
		[]*Frame{
			SelfVerifyFrame(FrameLimits{Execution: 5_000}),
			UserOpFrame(&testTarget, uint256.NewInt(1), data, FrameLimits{Execution: 30_000, State: 1_000}),
		},
		[]*FrameSignature{{Scheme: SigSchemeSecp256k1, Signature: bytes.Repeat([]byte{0x11}, 65)}},
	)

	// 12000 base + 2*475 per frame + 2800 signature verification.
	want := uint64(12_000 + 2*475 + 2800)
	// Frame data: two zero bytes (1 token each) and two non-zero (4 tokens each).
	want += StandardTokenCost * (1 + 4 + 4 + 1)
	// Signature bytes: 65 non-zero bytes, empty signer and msg.
	want += StandardTokenCost * (65 * 4)
	// The user op carries value to a foreign target.
	want += TxValueCost
	// Both envelope extensions price their own encoding as transaction data. The
	// default [0] key set encodes as rlp([0]) || rlp(0) = c1 80 80, and the empty
	// recent root list as c0: four non-zero bytes in total.
	extensionTokens := uint64(4 * 4)
	want += StandardTokenCost * extensionTokens

	if got := tx.IntrinsicGas(); got != want {
		t.Fatalf("intrinsic gas = %d, want %d", got, want)
	}

	floorTokens := uint64(len(data)+65)*StandardTokenCost + extensionTokens
	wantFloor := uint64(12_000+2*475+2800) + TxValueCost + TotalCostFloorPerToken*floorTokens

	if got := tx.CalldataFloorGas(); got != wantFloor {
		t.Fatalf("calldata floor gas = %d, want %d", got, wantFloor)
	}

	wantMax := want + 5_000 + 30_000 + 1_000
	if got := tx.MaxGas(); got != wantMax {
		t.Fatalf("max gas = %d, want %d", got, wantMax)
	}

	if got := tx.GetStateGas(); got != 1_000 {
		t.Fatalf("state gas = %d, want 1000", got)
	}

	// MaxCost is what APPROVE collects from the payer up front.
	wantCost := new(big.Int).Mul(big.NewInt(int64(wantMax)), big.NewInt(1e9))
	if got := tx.MaxCost(nil); got.Cmp(wantCost) != 0 {
		t.Fatalf("max cost = %s, want %s", got, wantCost)
	}
}

// TestFrameValidatePayload checks the static validity rules.
func TestFrameValidatePayload(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed generating key: %v", err)
	}

	sender := crypto.PubkeyToAddress(key.PublicKey)

	// ValidatePayload checks the encoded signature bytes, so build signed.
	valid := func() *FrameTx {
		tx := buildTransferTx(t, sender)
		if err := tx.SignPayload(testChainID.ToBig(), key); err != nil {
			t.Fatalf("failed signing: %v", err)
		}

		return tx
	}

	if err := valid().ValidatePayload(); err != nil {
		t.Fatalf("canonical transfer should be valid: %v", err)
	}

	if err := valid().VerifySignatures(); err != nil {
		t.Fatalf("canonical transfer signatures should verify: %v", err)
	}

	cases := []struct {
		name   string
		mutate func(tx *FrameTx)
	}{
		{"no frames", func(tx *FrameTx) { tx.Frames = nil }},
		{"too many frames", func(tx *FrameTx) {
			tx.Frames = make([]*Frame, MaxFrames+1)
			for i := range tx.Frames {
				tx.Frames[i] = SelfVerifyFrame(FrameLimits{})
			}
		}},
		{"reserved flag bits", func(tx *FrameTx) { tx.Frames[1].Flags = 0x8 }},
		{"value outside SENDER mode", func(tx *FrameTx) { tx.Frames[0].Value = uint256.NewInt(1) }},
		{"unknown mode", func(tx *FrameTx) { tx.Frames[1].Mode = 3 }},
		{"unknown signature scheme", func(tx *FrameTx) { tx.Signatures[0].Scheme = 0x7 }},
		{"arbitrary entry with signer", func(tx *FrameTx) {
			tx.Signatures[0].Scheme = SigSchemeArbitrary
			tx.Signatures[0].Signer = testTarget.Bytes()
		}},
		{"short signer", func(tx *FrameTx) { tx.Signatures[0].Signer = []byte{0x01} }},
		{"bad msg length", func(tx *FrameTx) { tx.Signatures[0].Msg = []byte{0x01} }},
		{"zero explicit digest", func(tx *FrameTx) { tx.Signatures[0].Msg = make([]byte, 32) }},
		{"blob fee without blobs", func(tx *FrameTx) { tx.Fees.BlobFeeCap = uint256.NewInt(1) }},
		{"bad blob hash version", func(tx *FrameTx) {
			tx.BlobHashes = []common.Hash{{0x02}}
			tx.Fees.BlobFeeCap = uint256.NewInt(1)
		}},
		{"approve execution for foreign target", func(tx *FrameTx) { tx.Frames[0].Target = &testTarget }},
		{"batch flag on last frame", func(tx *FrameTx) { tx.Frames[1].Flags |= AtomicBatchFlag }},
		{"batch flag in VERIFY mode", func(tx *FrameTx) { tx.Frames[0].Flags |= AtomicBatchFlag }},
		{"execution gas over the cap", func(tx *FrameTx) {
			tx.Frames[1].Limits.Execution = TxMaxGasLimit
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tx := valid()
			tc.mutate(tx)

			if err := tx.ValidatePayload(); err == nil {
				t.Fatal("expected a validation error")
			}
		})
	}
}

// TestFrameMempoolPrefix checks the four recognized validation prefixes and the
// verification gas caps.
func TestFrameMempoolPrefix(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed generating key: %v", err)
	}

	sender := crypto.PubkeyToAddress(key.PublicKey)
	paymaster := common.HexToAddress("0xbbbb000000000000000000000000000000000002")
	factory := common.HexToAddress("0xcccc000000000000000000000000000000000003")
	verifyLimits := FrameLimits{Execution: 5_000}

	newTx := func(frames ...*Frame) *FrameTx {
		tx := NewFrameTx(testChainID, sender, 0,
			FrameFees{GasFeeCap: uint256.NewInt(1e9)}, frames,
			[]*FrameSignature{SenderSignature()})

		if err := tx.SignPayload(testChainID.ToBig(), key); err != nil {
			t.Fatalf("failed signing: %v", err)
		}

		return tx
	}

	userOp := func() *Frame {
		return UserOpFrame(&testTarget, nil, nil, FrameLimits{Execution: 30_000})
	}

	valid := []struct {
		name string
		tx   *FrameTx
	}{
		{"self relay", newTx(SelfVerifyFrame(verifyLimits), userOp())},
		{"deploy then self relay", newTx(
			DeployFrame(factory, []byte{0x01}, FrameLimits{Execution: 20_000}),
			SelfVerifyFrame(verifyLimits), userOp())},
		{"canonical paymaster", newTx(
			OnlyVerifyFrame(verifyLimits),
			PayFrame(paymaster, nil, verifyLimits), userOp())},
		{"deploy then paymaster", newTx(
			DeployFrame(factory, []byte{0x01}, FrameLimits{Execution: 20_000}),
			OnlyVerifyFrame(verifyLimits),
			PayFrame(paymaster, nil, verifyLimits), userOp())},
		{"expiry then self relay", newTx(
			ExpiryFrame(1<<40, 5_000),
			SelfVerifyFrame(verifyLimits), userOp())},
		{"atomic batch after the prefix", newTx(
			SelfVerifyFrame(verifyLimits),
			userOp().WithAtomicBatch(), userOp())},
	}

	for _, tc := range valid {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.tx.ValidatePayload(); err != nil {
				t.Fatalf("payload should be valid: %v", err)
			}

			if err := tc.tx.ValidateMempoolPrefix(); err != nil {
				t.Fatalf("prefix should be accepted: %v", err)
			}
		})
	}

	invalid := []struct {
		name string
		tx   *FrameTx
	}{
		{"no payment approval", newTx(OnlyVerifyFrame(verifyLimits), userOp())},
		{"pay without only_verify", newTx(PayFrame(paymaster, nil, verifyLimits), userOp())},
		{"deploy after verify", newTx(
			SelfVerifyFrame(verifyLimits),
			DeployFrame(factory, nil, FrameLimits{Execution: 20_000}),
			PayFrame(paymaster, nil, verifyLimits))},
		{"verify frame after the prefix", newTx(
			SelfVerifyFrame(verifyLimits), userOp(), SelfVerifyFrame(verifyLimits))},
		{"expiry frame not first", newTx(
			SelfVerifyFrame(verifyLimits), ExpiryFrame(1<<40, 5_000), userOp())},
		{"prefix over the verification gas cap", newTx(
			SelfVerifyFrame(FrameLimits{Execution: MaxVerifyGas}), userOp())},
		{"prefix over the state gas cap", newTx(
			SelfVerifyFrame(FrameLimits{Execution: 5_000, State: MaxVerifyStateGas + 1}), userOp())},
	}

	for _, tc := range invalid {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.tx.ValidateMempoolPrefix(); err == nil {
				t.Fatal("expected the prefix to be rejected")
			}
		})
	}
}

// TestFrameReceiptDecode checks the frame receipt decoder, including the derived
// transaction status the consensus receipt does not carry.
func TestFrameReceiptDecode(t *testing.T) {
	raw := `{
		"type": "0x6",
		"cumulativeGasUsed": "0x1234",
		"transactionHash": "0xeeee000000000000000000000000000000000000000000000000000000000009",
		"payer": "0xdddd000000000000000000000000000000000004",
		"frames": [
			{"status": "0x1", "gasUsed": {"execution": "0x2710", "state": "0x0"}, "logs": []},
			{"status": "0x1", "gasUsed": ["0x7530", "0x3e8"], "logs": []},
			{"status": "0x2", "gasUsed": {"execution": "0x0", "state": "0x0"}, "logs": []}
		]
	}`

	var receipt Receipt
	if err := json.Unmarshal([]byte(raw), &receipt); err != nil {
		t.Fatalf("failed decoding frame receipt: %v", err)
	}

	extra := receipt.FrameExtra()
	if extra == nil {
		t.Fatal("frame receipt content missing")
	}

	if extra.Payer != common.HexToAddress("0xdddd000000000000000000000000000000000004") {
		t.Fatalf("wrong payer: %s", extra.Payer)
	}

	if len(extra.Frames) != 3 {
		t.Fatalf("expected 3 frame receipts, got %d", len(extra.Frames))
	}

	if extra.Frames[0].ExecutionGas != 10_000 || extra.Frames[1].ExecutionGas != 30_000 {
		t.Fatal("execution gas was not decoded from both shapes")
	}

	if extra.Frames[1].StateGas != 1_000 {
		t.Fatalf("state gas = %d, want 1000", extra.Frames[1].StateGas)
	}

	if !extra.Frames[2].Skipped() {
		t.Fatal("third frame should be marked skipped")
	}

	if extra.FailedFrame() != -1 {
		t.Fatal("no frame failed")
	}

	// The consensus receipt carries no transaction status, so it is derived.
	if !receipt.Successful() {
		t.Fatal("receipt status should be derived as successful")
	}

	if receipt.GasUsed != 41_000 {
		t.Fatalf("gas used = %d, want 41000", receipt.GasUsed)
	}

	// A failed frame must make the derived status a failure.
	failing := `{"type":"0x6","transactionHash":"0xeeee00000000000000000000000000000000000000000000000000000000000a",
		"payer":"0xdddd000000000000000000000000000000000004",
		"frames":[{"status":"0x1","gasUsed":["0x10","0x0"]},{"status":"0x0","gasUsed":["0x20","0x0"]}]}`

	var failed Receipt
	if err := json.Unmarshal([]byte(failing), &failed); err != nil {
		t.Fatalf("failed decoding: %v", err)
	}

	if failed.Successful() {
		t.Fatal("a failed frame should make the transaction status a failure")
	}

	if failed.FrameExtra().FailedFrame() != 1 {
		t.Fatal("wrong failed frame index")
	}
}

// TestFrameReceiptWithoutFrameFields checks that a node reporting a frame transaction
// without the frame-specific fields still yields a usable receipt.
func TestFrameReceiptWithoutFrameFields(t *testing.T) {
	raw := `{"type":"0x6","status":"0x1","gasUsed":"0x5208","cumulativeGasUsed":"0x5208",
		"transactionHash":"0xeeee00000000000000000000000000000000000000000000000000000000000b","logs":[]}`

	var receipt Receipt
	if err := json.Unmarshal([]byte(raw), &receipt); err != nil {
		t.Fatalf("failed decoding: %v", err)
	}

	if receipt.FrameExtra() != nil {
		t.Fatal("no frame content should have been decoded")
	}

	if !receipt.Successful() || receipt.GasUsed != 21_000 {
		t.Fatal("generic receipt fields should still be usable")
	}
}

// TestFrameTxNotRepresentableInGeth checks that conversion fails loudly rather than
// producing a transaction with a different hash.
func TestFrameTxNotRepresentableInGeth(t *testing.T) {
	sender := common.HexToAddress("0xaaaa000000000000000000000000000000000001")
	tx := NewTx(buildTransferTx(t, sender))

	if _, err := tx.ToGethTx(); err == nil {
		t.Fatal("expected conversion to go-ethereum to fail")
	}
}

// generateP256Key returns a fresh key on the P-256 curve.
func generateP256Key() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

// TestFrameExpiryValidity checks the expiry verifier frame's own validity rules.
func TestFrameExpiryValidity(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed generating key: %v", err)
	}

	sender := crypto.PubkeyToAddress(key.PublicKey)

	build := func(mutate func(expiry *Frame) []*Frame) *FrameTx {
		expiry := ExpiryFrame(1<<40, 5_000)

		frames := []*Frame{expiry, SelfVerifyFrame(FrameLimits{Execution: 5_000}),
			UserOpFrame(&testTarget, nil, nil, FrameLimits{Execution: 30_000})}

		if mutate != nil {
			frames = mutate(expiry)
		}

		tx := NewFrameTx(testChainID, sender, 0,
			FrameFees{GasFeeCap: uint256.NewInt(1e9)}, frames,
			[]*FrameSignature{SenderSignature()})

		if err := tx.SignPayload(testChainID.ToBig(), key); err != nil {
			t.Fatalf("failed signing: %v", err)
		}

		return tx
	}

	if err := build(nil).ValidatePayload(); err != nil {
		t.Fatalf("a well-formed expiry frame should be valid: %v", err)
	}

	if deadline, ok := build(nil).Frames[0].ExpiryDeadline(); !ok || deadline != 1<<40 {
		t.Fatalf("expiry deadline did not round-trip, got %d (ok=%v)", deadline, ok)
	}

	cases := []struct {
		name   string
		mutate func(expiry *Frame) []*Frame
	}{
		{"non-zero flags", func(e *Frame) []*Frame {
			e.Flags = ApprovePayment

			return []*Frame{e, SelfVerifyFrame(FrameLimits{Execution: 5_000})}
		}},
		{"carries value", func(e *Frame) []*Frame {
			e.Value = uint256.NewInt(1)

			return []*Frame{e, SelfVerifyFrame(FrameLimits{Execution: 5_000})}
		}},
		{"budgets state gas", func(e *Frame) []*Frame {
			e.Limits.State = 1

			return []*Frame{e, SelfVerifyFrame(FrameLimits{Execution: 5_000})}
		}},
		{"wrong data length", func(e *Frame) []*Frame {
			e.Data = []byte{0x01}

			return []*Frame{e, SelfVerifyFrame(FrameLimits{Execution: 5_000})}
		}},
		{"two expiry frames", func(e *Frame) []*Frame {
			return []*Frame{e, ExpiryFrame(1<<40, 5_000), SelfVerifyFrame(FrameLimits{Execution: 5_000})}
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := build(tc.mutate).ValidatePayload(); err == nil {
				t.Fatal("expected a validation error")
			}
		})
	}

	// A malformed expiry frame must not be mistaken for a self_verify frame.
	malformed := &Frame{Mode: FrameModeVerify, Flags: ApproveExecutionAndPayment,
		Target: &ExpiryVerifier, Value: new(uint256.Int)}
	if species := malformed.Species(sender); species != SpeciesExpiryVerify {
		t.Fatalf("a VERIFY frame targeting the expiry verifier classified as %s", species)
	}
}

// TestFrameExpiryPlacement checks that the expiry frame may lead the frame list and
// nothing else, and that a deploy frame may follow it.
func TestFrameExpiryPlacement(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed generating key: %v", err)
	}

	sender := crypto.PubkeyToAddress(key.PublicKey)
	factory := common.HexToAddress("0xcccc000000000000000000000000000000000003")
	verifyLimits := FrameLimits{Execution: 5_000}

	newTx := func(frames ...*Frame) *FrameTx {
		tx := NewFrameTx(testChainID, sender, 0,
			FrameFees{GasFeeCap: uint256.NewInt(1e9)}, frames,
			[]*FrameSignature{SenderSignature()})

		if err := tx.SignPayload(testChainID.ToBig(), key); err != nil {
			t.Fatalf("failed signing: %v", err)
		}

		return tx
	}

	userOp := UserOpFrame(&testTarget, nil, nil, FrameLimits{Execution: 30_000})

	// [expiry, deploy, self_verify] matches [deploy, self_verify] once the expiry
	// frame is skipped.
	withDeploy := newTx(
		ExpiryFrame(1<<40, 5_000),
		DeployFrame(factory, []byte{0x01}, FrameLimits{Execution: 20_000}),
		SelfVerifyFrame(verifyLimits), userOp)

	if err := withDeploy.ValidateMempoolPrefix(); err != nil {
		t.Fatalf("expiry followed by deploy should be accepted: %v", err)
	}

	// An expiry frame anywhere but the front is rejected, even outside the prefix.
	trailing := newTx(SelfVerifyFrame(verifyLimits), userOp, ExpiryFrame(1<<40, 5_000))
	if err := trailing.ValidateMempoolPrefix(); err == nil {
		t.Fatal("an expiry frame after the first position should be rejected")
	}
}

// TestFrameSignatureBytesValidation checks the encoding rules that keep each signature
// to a single valid form.
func TestFrameSignatureBytesValidation(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed generating key: %v", err)
	}

	sender := crypto.PubkeyToAddress(key.PublicKey)

	build := func(mutate func(sig *FrameSignature)) *FrameTx {
		tx := buildTransferTx(t, sender)
		if err := tx.SignPayload(testChainID.ToBig(), key); err != nil {
			t.Fatalf("failed signing: %v", err)
		}

		mutate(tx.Signatures[0])

		return tx
	}

	cases := []struct {
		name   string
		mutate func(sig *FrameSignature)
	}{
		{"short signature", func(s *FrameSignature) { s.Signature = s.Signature[:64] }},
		{"recovery id out of range", func(s *FrameSignature) { s.Signature[0] = 27 }},
		{"zero r", func(s *FrameSignature) { copy(s.Signature[1:33], make([]byte, 32)) }},
		{"zero s", func(s *FrameSignature) { copy(s.Signature[33:65], make([]byte, 32)) }},
		{"high s", func(s *FrameSignature) {
			// N - s is the malleable counterpart of a canonical low-s value.
			order := crypto.S256().Params().N
			s2 := new(big.Int).Sub(order, new(big.Int).SetBytes(s.Signature[33:65]))
			s2.FillBytes(s.Signature[33:65])
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := build(tc.mutate).ValidatePayload(); err == nil {
				t.Fatal("expected a validation error")
			}
		})
	}

	// Tampering with a covered field leaves the encoding valid but the signature wrong.
	tampered := build(func(*FrameSignature) {})
	tampered.Frames[1].Value = uint256.NewInt(999)

	if err := tampered.ValidatePayload(); err != nil {
		t.Fatalf("a tampered transaction should still be structurally valid: %v", err)
	}

	if err := tampered.VerifySignatures(); err == nil {
		t.Fatal("VerifySignatures should reject a transaction signed over different content")
	}
}

// TestFrameEnvelopeShapes checks that each of the four payload shapes encodes, decodes
// back to itself, and is told apart from the others.
//
// EIP-8250 and EIP-8272 amend EIP-8141's envelope independently, so a chain may run
// any combination and the shape has to be read off the payload rather than assumed.
func TestFrameEnvelopeShapes(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed generating key: %v", err)
	}

	sender := crypto.PubkeyToAddress(key.PublicKey)

	shapes := []struct {
		name       string
		extensions FrameExtensions
		fields     int
	}{
		{"8141", 0, 7},
		{"8141+8250", FrameExtKeyedNonces, 8},
		{"8141+8272", FrameExtRecentRoots, 8},
		{"8141+8250+8272", FrameExtAll, 9},
	}

	for _, shape := range shapes {
		t.Run(shape.name, func(t *testing.T) {
			tx := NewFrameTxWithExtensions(shape.extensions, testChainID, sender, 7,
				FrameFees{GasTipCap: uint256.NewInt(1e9), GasFeeCap: uint256.NewInt(20e9)},
				[]*Frame{
					SelfVerifyFrame(FrameLimits{Execution: 5_000}),
					UserOpFrame(&testTarget, nil, nil, FrameLimits{Execution: 30_000}),
				},
				[]*FrameSignature{SenderSignature()})

			if shape.extensions.Has(FrameExtRecentRoots) {
				tx.RecentRoots = []*RecentRootReference{
					{SourceID: common.HexToHash("0x01"), Slot: 9, Root: common.HexToHash("0x02")},
				}
			}

			if tx.Extensions.String() != shape.name {
				t.Fatalf("extension name = %s, want %s", tx.Extensions, shape.name)
			}

			signed, err := SignTx(NewTx(tx), testChainID.ToBig(), key)
			if err != nil {
				t.Fatalf("failed signing: %v", err)
			}

			inner, _ := signed.Inner().(*FrameTx)
			if err := inner.ValidatePayload(); err != nil {
				t.Fatalf("payload should be valid: %v", err)
			}

			encoded, err := signed.MarshalBinary()
			if err != nil {
				t.Fatalf("failed encoding: %v", err)
			}

			// The payload must carry exactly the field count the shape implies.
			count, err := countRLPFields(encoded[1:])
			if err != nil {
				t.Fatalf("failed counting fields: %v", err)
			}

			if count != shape.fields {
				t.Fatalf("payload has %d fields, want %d", count, shape.fields)
			}

			decoded, err := DecodeTx(encoded)
			if err != nil {
				t.Fatalf("failed decoding: %v", err)
			}

			back, ok := decoded.Inner().(*FrameTx)
			if !ok {
				t.Fatalf("decoded to %T", decoded.Inner())
			}

			if back.Extensions != shape.extensions {
				t.Fatalf("extensions changed: %s != %s", back.Extensions, shape.extensions)
			}

			if decoded.Hash() != signed.Hash() {
				t.Fatalf("hash changed: %s != %s", decoded.Hash(), signed.Hash())
			}

			if back.NonceSeq != 7 {
				t.Fatalf("nonce did not survive: %d", back.NonceSeq)
			}

			if back.HasKeyedNonces() != shape.extensions.Has(FrameExtKeyedNonces) {
				t.Fatal("keyed nonce flag did not survive")
			}

			if len(back.RecentRoots) != len(tx.RecentRoots) {
				t.Fatalf("recent roots did not survive: %d != %d", len(back.RecentRoots), len(tx.RecentRoots))
			}

			// A scalar-nonce transaction must not claim keys, and vice versa.
			if !shape.extensions.Has(FrameExtKeyedNonces) && len(back.NonceKeys) != 0 {
				t.Fatal("scalar nonce shape decoded nonce keys")
			}
		})
	}
}

// countRLPFields returns the number of top-level elements in an RLP list.
func countRLPFields(payload []byte) (int, error) {
	content, _, err := rlp.SplitList(payload)
	if err != nil {
		return 0, err
	}

	count := 0

	for len(content) > 0 {
		_, _, rest, err := rlp.Split(content)
		if err != nil {
			return 0, err
		}

		count++
		content = rest
	}

	return count, nil
}

// TestFrameExtensionMismatch checks that fields belonging to an inactive extension are
// rejected rather than silently dropped at encoding time.
func TestFrameExtensionMismatch(t *testing.T) {
	sender := common.HexToAddress("0xaaaa000000000000000000000000000000000001")

	base := func() *FrameTx {
		return NewFrameTxWithExtensions(0, testChainID, sender, 1,
			FrameFees{GasFeeCap: uint256.NewInt(1e9)},
			[]*Frame{SelfVerifyFrame(FrameLimits{Execution: 5_000})},
			[]*FrameSignature{{Scheme: SigSchemeSecp256k1, Signature: bytes.Repeat([]byte{0x11}, 65)}})
	}

	withKeys := base()
	withKeys.NonceKeys = []*uint256.Int{new(uint256.Int)}

	if err := withKeys.ValidatePayload(); err == nil {
		t.Fatal("nonce keys without the EIP-8250 extension should be rejected")
	}

	withRoots := base()
	withRoots.RecentRoots = []*RecentRootReference{{}}

	if err := withRoots.ValidatePayload(); err == nil {
		t.Fatal("recent roots without the EIP-8272 extension should be rejected")
	}
}

// TestFramePostTxMode checks EIP-7906's POST_TX mode: it is a valid mode, it must form
// a trailing suffix, and it is classified for display.
func TestFramePostTxMode(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed generating key: %v", err)
	}

	sender := crypto.PubkeyToAddress(key.PublicKey)
	assertions := common.HexToAddress("0xeeee000000000000000000000000000000000006")

	build := func(frames ...*Frame) *FrameTx {
		tx := NewFrameTx(testChainID, sender, 0,
			FrameFees{GasFeeCap: uint256.NewInt(1e9)}, frames,
			[]*FrameSignature{SenderSignature()})

		if err := tx.SignPayload(testChainID.ToBig(), key); err != nil {
			t.Fatalf("failed signing: %v", err)
		}

		return tx
	}

	userOp := func() *Frame {
		return UserOpFrame(&testTarget, nil, nil, FrameLimits{Execution: 30_000})
	}

	postTx := func() *Frame {
		return PostTxFrame(assertions, []byte{0x01}, FrameLimits{Execution: 20_000})
	}

	// A trailing suffix of POST_TX frames is valid.
	valid := build(SelfVerifyFrame(FrameLimits{Execution: 5_000}), userOp(), postTx(), postTx())
	if err := valid.ValidatePayload(); err != nil {
		t.Fatalf("a trailing POST_TX suffix should be valid: %v", err)
	}

	if err := valid.ValidateMempoolPrefix(); err != nil {
		t.Fatalf("POST_TX frames sit outside the validation prefix: %v", err)
	}

	if species := valid.Frames[2].Species(sender); species != SpeciesPostTx {
		t.Fatalf("POST_TX frame classified as %s", species)
	}

	if idx := valid.PostTxIndex(); idx != 2 {
		t.Fatalf("PostTxIndex = %d, want 2", idx)
	}

	// A non-POST_TX frame after a POST_TX frame is invalid.
	interleaved := build(SelfVerifyFrame(FrameLimits{Execution: 5_000}), postTx(), userOp())
	if err := interleaved.ValidatePayload(); err == nil {
		t.Fatal("a frame following a POST_TX frame should be rejected")
	}

	// Mode 4 remains unknown.
	bad := build(SelfVerifyFrame(FrameLimits{Execution: 5_000}), userOp())
	bad.Frames[1].Mode = 4

	if err := bad.ValidatePayload(); err == nil {
		t.Fatal("mode 4 should be rejected")
	}

	// A POST_TX frame carrying value is rejected, as for any non-SENDER frame.
	valued := build(SelfVerifyFrame(FrameLimits{Execution: 5_000}), postTx())
	valued.Frames[1].Value = uint256.NewInt(1)

	if err := valued.ValidatePayload(); err == nil {
		t.Fatal("a POST_TX frame carrying value should be rejected")
	}
}

// TestDurableFrames checks which frames' effects survive, which a frame's own status
// does not answer: an unrolled atomic batch and a failed POST_TX both discard the
// effects of frames that report success.
func TestDurableFrames(t *testing.T) {
	sender := common.HexToAddress("0xaaaa000000000000000000000000000000000001")
	assertions := common.HexToAddress("0xeeee000000000000000000000000000000000006")

	userOp := func() *Frame {
		return UserOpFrame(&testTarget, nil, nil, FrameLimits{Execution: 30_000})
	}

	build := func(frames ...*Frame) *FrameTx {
		return NewFrameTx(testChainID, sender, 0,
			FrameFees{GasFeeCap: uint256.NewInt(1e9)}, frames,
			[]*FrameSignature{SenderSignature()})
	}

	statuses := func(codes ...uint64) *FrameReceiptExtra {
		extra := &FrameReceiptExtra{}
		for _, code := range codes {
			extra.Frames = append(extra.Frames, &FrameReceipt{Status: code})
		}

		return extra
	}

	t.Run("all succeed", func(t *testing.T) {
		tx := build(SelfVerifyFrame(FrameLimits{Execution: 5_000}), userOp(), userOp())
		got := statuses(1, 1, 1).DurableFrames(tx)

		for i, ok := range got {
			if !ok {
				t.Fatalf("frame %d should be durable", i)
			}
		}
	})

	t.Run("atomic batch unrolled", func(t *testing.T) {
		// Frames 1 and 2 form a batch with frame 3 closing it; frame 2 fails, so all
		// three lose their effects while frame 1 still reports success.
		tx := build(SelfVerifyFrame(FrameLimits{Execution: 5_000}),
			userOp().WithAtomicBatch(), userOp().WithAtomicBatch(), userOp())
		got := statuses(1, 1, 0, 2).DurableFrames(tx)

		want := []bool{true, false, false, false}
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("frame %d durable = %v, want %v (got %v)", i, got[i], want[i], got)
			}
		}
	})

	t.Run("post-tx reverts the body", func(t *testing.T) {
		// Every frame reports success except the POST_TX one, yet nothing after the
		// validation prefix survives.
		tx := build(SelfVerifyFrame(FrameLimits{Execution: 5_000}), userOp(), userOp(),
			PostTxFrame(assertions, nil, FrameLimits{Execution: 20_000}))
		extra := statuses(1, 1, 1, 0)

		if !extra.PostTxReverted(tx) {
			t.Fatal("a failed POST_TX frame should revert the body")
		}

		got := extra.DurableFrames(tx)

		want := []bool{true, false, false, false}
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("frame %d durable = %v, want %v (got %v)", i, got[i], want[i], got)
			}
		}
	})

	t.Run("post-tx succeeds", func(t *testing.T) {
		tx := build(SelfVerifyFrame(FrameLimits{Execution: 5_000}), userOp(),
			PostTxFrame(assertions, nil, FrameLimits{Execution: 20_000}))
		extra := statuses(1, 1, 1)

		if extra.PostTxReverted(tx) {
			t.Fatal("a successful POST_TX frame should not revert the body")
		}

		for i, ok := range extra.DurableFrames(tx) {
			if !ok {
				t.Fatalf("frame %d should be durable", i)
			}
		}
	})
}
