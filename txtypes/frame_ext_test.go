package txtypes

import (
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
