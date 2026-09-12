package txtypes

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
)

// keccakHex hashes a hex string, giving the tests a construction independent of the
// byte-slice building the helpers do.
func keccakHex(t *testing.T, s string) common.Hash {
	t.Helper()

	raw, err := hex.DecodeString(s)
	if err != nil {
		t.Fatalf("bad test vector %q: %v", s, err)
	}

	return crypto.Keccak256Hash(raw)
}

func TestNonceManagerSlot(t *testing.T) {
	sender := common.HexToAddress("0x8943545177806ed17b9f23f0a21ee5948ecaa776")

	// keccak256(pad32(sender) || pad32(key)), with both operands written out.
	want := keccakHex(t,
		"0000000000000000000000008943545177806ed17b9f23f0a21ee5948ecaa776"+
			"0000000000000000000000000000000000000000000000000000000000000007")

	if got := NonceManagerSlot(sender, uint256.NewInt(7)); got != want {
		t.Errorf("slot for key 7 = %s, want %s", got, want)
	}

	// Distinct keys and distinct senders must not collide.
	other := NonceManagerSlot(sender, uint256.NewInt(9))
	if other == want {
		t.Error("keys 7 and 9 map to the same slot")
	}

	if NonceManagerSlot(common.Address{}, uint256.NewInt(7)) == want {
		t.Error("different senders map to the same slot")
	}
}

func TestNonceKeysHash(t *testing.T) {
	// keccak256(uint256(len) || uint256(k) for k in keys)
	want := keccakHex(t,
		"0000000000000000000000000000000000000000000000000000000000000002"+
			"0000000000000000000000000000000000000000000000000000000000000007"+
			"0000000000000000000000000000000000000000000000000000000000000009")

	got := NonceKeysHash([]*uint256.Int{uint256.NewInt(7), uint256.NewInt(9)})
	if got != want {
		t.Errorf("nonce keys hash = %s, want %s", got, want)
	}

	// The hash must distinguish key sets, not just their length.
	if NonceKeysHash([]*uint256.Int{uint256.NewInt(9), uint256.NewInt(7)}) == want {
		t.Error("key order does not affect the hash")
	}
}

func TestRecentRootDerivations(t *testing.T) {
	source := common.HexToAddress("0x8943545177806ed17b9f23f0a21ee5948ecaa776")
	salt := common.HexToHash("0x01")
	root := common.HexToHash("0x02")

	sourceID := RecentRootSourceID(source, salt)

	wantSource := keccakHex(t,
		"8943545177806ed17b9f23f0a21ee5948ecaa776"+
			"0000000000000000000000000000000000000000000000000000000000000001")
	if sourceID != wantSource {
		t.Errorf("source id = %s, want %s", sourceID, wantSource)
	}

	// The write calldata is the concatenation the contract parses positionally.
	calldata := RecentRootWriteCalldata(salt, root)
	if len(calldata) != RecentRootWriteLength {
		t.Fatalf("write calldata is %d bytes, want %d", len(calldata), RecentRootWriteLength)
	}

	if common.BytesToHash(calldata[:32]) != salt || common.BytesToHash(calldata[32:]) != root {
		t.Error("write calldata does not carry salt followed by root")
	}

	// entry_hash and storage_key must be domain-separated: the same source and slot
	// must not produce the same 32 bytes for both.
	slot := uint64(1234)
	if RecentRootEntryHash(sourceID, slot, root) == RecentRootStorageKey(sourceID, slot) {
		t.Error("entry hash and storage key are not domain separated")
	}

	// The ring wraps every RecentRootLength slots, so two slots one period apart share
	// a storage key while their entry hashes differ.
	if RecentRootSlotStorageKey(sourceID, slot) != RecentRootSlotStorageKey(sourceID, slot+RecentRootLength) {
		t.Error("storage key does not wrap at RecentRootLength")
	}

	if RecentRootEntryHash(sourceID, slot, root) == RecentRootEntryHash(sourceID, slot+RecentRootLength, root) {
		t.Error("entry hash does not bind the slot")
	}
}

func TestRecentRootVerifyData(t *testing.T) {
	// The EIP's reference vector: source 0x…01 with a zero salt, slot 1, root 2.
	source := common.HexToAddress("0x0000000000000000000000000000000000000001")
	sourceID := RecentRootSourceID(source, common.Hash{})

	if sourceID != common.HexToHash("0xb9382d35273c75a50631a3e84d3c75ec9266e2b18c35a627e16cdbf26a18ca85") {
		t.Fatalf("source id = %s does not match the EIP's vector", sourceID)
	}

	ref := &RecentRootReference{SourceID: sourceID, Slot: 1, Root: common.HexToHash("0x02")}

	if got := RecentRootEntryHash(ref.SourceID, ref.Slot, ref.Root); got != common.HexToHash("0x0a0d1254c851be5a133b4c9a9e300f5602fc0f43dbe65aa6a66930d4ca0a51b8") {
		t.Errorf("entry hash = %s does not match the EIP's vector", got)
	}

	if got := RecentRootSlotStorageKey(ref.SourceID, ref.Slot); got != common.HexToHash("0x5f027aa1cbe2df279bf6518edd4b44ea5409fd800189ec35224e10ab05e574c3") {
		t.Errorf("storage key = %s does not match the EIP's vector", got)
	}

	want := common.FromHex("0xb9382d35273c75a50631a3e84d3c75ec9266e2b18c35a627e16cdbf26a18ca85" +
		"0000000000000001" +
		"0000000000000000000000000000000000000000000000000000000000000002")

	data := RecentRootVerifyData([]*RecentRootReference{ref})
	if !bytes.Equal(data, want) {
		t.Fatalf("verify data = %x, want %x", data, want)
	}

	back, err := ParseRecentRootVerifyData(data)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if len(back) != 1 || *back[0] != *ref {
		t.Fatalf("parsed %+v, want %+v", back, ref)
	}

	// The contract rejects empty data, misaligned data and more than the cap.
	for name, bad := range map[string][]byte{
		"empty":     nil,
		"71 bytes":  data[:71],
		"73 bytes":  append(append([]byte{}, data...), 0),
		"17 tuples": bytes.Repeat(data, MaxRecentRootReferences+1),
	} {
		if _, err := ParseRecentRootVerifyData(bad); err == nil {
			t.Errorf("%s: parsed without error", name)
		}
	}
}

func TestRecentRootVerifierFrame(t *testing.T) {
	sender := common.HexToAddress("0x8943545177806ed17b9f23f0a21ee5948ecaa776")
	ref := &RecentRootReference{SourceID: common.HexToHash("0x01"), Slot: 5, Root: common.HexToHash("0x02")}

	frame := RecentRootVerifyFrame([]*RecentRootReference{ref}, RecentRootVerifyGas(1))
	if !frame.IsRecentRootVerifier() {
		t.Fatal("the built frame is not recognized as a recent root verifier")
	}

	if frame.Species(sender) != SpeciesRecentRootVerify {
		t.Fatalf("species = %s, want %s", frame.Species(sender), SpeciesRecentRootVerify)
	}

	// Each defining property, when missed, makes it an ordinary VERIFY frame that the
	// mempool does not treat as a protocol verifier.
	for name, mutate := range map[string]func(*Frame){
		"flags":        func(f *Frame) { f.Flags = ApprovePayment },
		"state limit":  func(f *Frame) { f.Limits.State = 1 },
		"value":        func(f *Frame) { f.Value = uint256.NewInt(1) },
		"data length":  func(f *Frame) { f.Data = f.Data[:71] },
		"mode":         func(f *Frame) { f.Mode = FrameModeSender },
		"other target": func(f *Frame) { f.Target = &sender },
	} {
		cpy := frame.Copy()
		mutate(cpy)

		if cpy.IsRecentRootVerifier() {
			t.Errorf("%s: still classified as a recent root verifier", name)
		}
	}

	build := func(frames ...*Frame) *FrameTx {
		return NewFrameTxWithExtensions(0, uint256.NewInt(1), sender, 0,
			FrameFees{GasFeeCap: uint256.NewInt(1)}, frames,
			[]*FrameSignature{{Scheme: SigSchemeSecp256k1, Signature: append([]byte{0}, bytes.Repeat([]byte{0x11}, 64)...)}})
	}

	roots := func() *Frame { return RecentRootVerifyFrame([]*RecentRootReference{ref}, RecentRootVerifyGas(1)) }
	verify := func() *Frame { return SelfVerifyFrame(FrameLimits{Execution: 5_000}) }
	userOp := func() *Frame { return UserOpFrame(&sender, nil, nil, FrameLimits{Execution: 30_000}) }

	// The frame's execution budget is excluded from the verification gas cap: a frame
	// that alone would exceed it stays mempool-legal.
	oversized := roots()
	oversized.Limits.Execution = MaxVerifyGas + 1

	for name, tc := range map[string]struct {
		tx    *FrameTx
		legal bool
	}{
		"leading":          {build(roots(), verify(), userOp()), true},
		"after expiry":     {build(ExpiryFrame(1, 5_000), roots(), verify(), userOp()), true},
		"over verify gas":  {build(oversized, verify(), userOp()), true},
		"before expiry":    {build(roots(), ExpiryFrame(1, 5_000), verify(), userOp()), false},
		"after validation": {build(verify(), roots(), userOp()), false},
		"twice":            {build(roots(), roots(), verify(), userOp()), false},
	} {
		err := tc.tx.ValidateMempoolPrefix()
		if tc.legal && err != nil {
			t.Errorf("%s: rejected: %v", name, err)
		}

		if !tc.legal && err == nil {
			t.Errorf("%s: accepted", name)
		}
	}
}

func TestRecentRootReferenceUsable(t *testing.T) {
	const current = 10_000

	for _, tc := range []struct {
		name string
		slot uint64
		want bool
	}{
		{"same slot is too new", current, false},
		{"future slot", current + 1, false},
		{"previous slot", current - 1, true},
		{"window edge", current - RecentRootUsableWindow, true},
		{"just outside the window", current - RecentRootUsableWindow - 1, false},
	} {
		if got := RecentRootReferenceUsable(current, tc.slot); got != tc.want {
			t.Errorf("%s: usable = %v, want %v", tc.name, got, tc.want)
		}
	}
}

func TestFrameTxUsesAccountNonce(t *testing.T) {
	for _, tc := range []struct {
		name string
		tx   *FrameTx
		want bool
	}{
		{
			name: "scalar nonce",
			tx:   &FrameTx{},
			want: true,
		},
		{
			name: "key zero aliases the account nonce",
			tx:   (&FrameTx{}).WithNonceKeys([]*uint256.Int{uint256.NewInt(0)}, 3),
			want: true,
		},
		{
			name: "non-zero key is an independent domain",
			tx:   (&FrameTx{}).WithNonceKeys([]*uint256.Int{uint256.NewInt(7)}, 3),
			want: false,
		},
		{
			name: "multiple keys are an independent domain",
			tx:   (&FrameTx{}).WithNonceKeys([]*uint256.Int{uint256.NewInt(7), uint256.NewInt(9)}, 3),
			want: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.tx.UsesAccountNonce(); got != tc.want {
				t.Errorf("UsesAccountNonce = %v, want %v", got, tc.want)
			}

			if got := NewTx(tc.tx).UsesAccountNonce(); got != tc.want {
				t.Errorf("Transaction.UsesAccountNonce = %v, want %v", got, tc.want)
			}
		})
	}

	// Every other type is sequenced by the account nonce.
	if !NewTx(&DynamicFeeTx{}).UsesAccountNonce() {
		t.Error("dynamic fee transactions must use the account nonce")
	}
}
