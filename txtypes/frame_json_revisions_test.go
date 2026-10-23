package txtypes

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// nethermindFrameTx returns nethermind's devnet-1 frame transaction fixture: bare JSON
// numbers for mode, flags and scheme, and executionGasLimit for the execution dimension.
func nethermindFrameTx(t *testing.T) string {
	t.Helper()

	for _, tc := range clientFrameTransactions {
		if tc.client == "nethermind" {
			return tc.raw
		}
	}

	t.Fatal("no nethermind fixture")

	return ""
}

// TestNethermindDevnet1FrameTxDecodes checks that nethermind's devnet-1 object decodes to
// the per-gas layout and re-encodes to the hash the node reported.
func TestNethermindDevnet1FrameTxDecodes(t *testing.T) {
	raw := nethermindFrameTx(t)

	tx, err := UnmarshalJSONTx([]byte(raw))
	if err != nil {
		t.Fatalf("failed decoding: %v", err)
	}

	frameTx, ok := tx.Inner().(*FrameTx)
	if !ok {
		t.Fatalf("decoded as %T, want a frame transaction", tx.Inner())
	}

	if !frameTx.LegacyFees || frameTx.Frames[0].Mode != FrameModeVerify || frameTx.Frames[0].Flags != 3 ||
		frameTx.Frames[1].Limits.Execution != 0x7530 || frameTx.Signatures[0].Scheme != SigSchemeSecp256k1 {
		t.Fatalf("decoded %+v", frameTx)
	}

	var reported struct{ Hash common.Hash }
	if err := json.Unmarshal([]byte(raw), &reported); err != nil {
		t.Fatalf("failed reading the reported hash: %v", err)
	}

	encoded, err := tx.MarshalBinary()
	if err != nil {
		t.Fatalf("failed encoding: %v", err)
	}

	if got := crypto.Keccak256Hash(encoded); got != reported.Hash {
		t.Fatalf("re-encoded hash = %s, want the reported %s", got, reported.Hash)
	}
}

// TestNethermindDevnet2FrameTxDecodes checks the devnet-2 object in nethermind's
// spelling: maxFee in place of the per-gas caps, and an empty key set. The object still
// reports the devnet-1 hash, so the encodings are compared by their own hashes.
func TestNethermindDevnet2FrameTxDecodes(t *testing.T) {
	raw := strings.NewReplacer(
		`"maxFeePerBlobGas": "0x0",`, `"nonceKeys": [],`,
		`"maxFeePerGas": "0x4a817c800",`, `"maxFee": "0x2386f26fc10000",`,
		`"gasPrice": "0x77359407",`, ``,
	).Replace(nethermindFrameTx(t))

	tx, err := UnmarshalJSONTx([]byte(raw))
	if err != nil {
		t.Fatalf("failed decoding: %v", err)
	}

	frameTx, ok := tx.Inner().(*FrameTx)
	if !ok {
		t.Fatalf("decoded as %T, want a frame transaction", tx.Inner())
	}

	if frameTx.LegacyFees || !frameTx.HasKeyedNonces() || !frameTx.UsesLegacyNonce() ||
		frameTx.GetMaxFee().Uint64() != 0x2386f26fc10000 || frameTx.NonceSeq != 0xd33 ||
		frameTx.Frames[1].Limits.Execution != 0x7530 || frameTx.Frames[0].Flags != 3 {
		t.Fatalf("decoded %+v", frameTx)
	}

	if feeListLength(t, frameTx) != 2 {
		t.Fatal("devnet-2 transaction does not encode a two-field fee list")
	}

	// Its own JSON and binary encodings carry the same transaction.
	encoded, err := tx.MarshalBinary()
	if err != nil {
		t.Fatalf("failed encoding: %v", err)
	}

	asJSON, err := tx.MarshalJSON()
	if err != nil {
		t.Fatalf("failed JSON encoding: %v", err)
	}

	back, err := UnmarshalJSONTx(asJSON)
	if err != nil {
		t.Fatalf("failed JSON decoding: %v", err)
	}

	reencoded, err := back.MarshalBinary()
	if err != nil {
		t.Fatalf("failed re-encoding: %v", err)
	}

	decoded, err := DecodeTx(encoded)
	if err != nil {
		t.Fatalf("failed decoding: %v", err)
	}

	if want := crypto.Keccak256Hash(encoded); crypto.Keccak256Hash(reencoded) != want || decoded.Hash() != want {
		t.Fatal("hash changed across encodings")
	}
}

// A frame transaction object with neither fee shape is refused rather than guessed at,
// since the guess decides the hash.
func TestFrameTxJSONWithoutFeesIsRefused(t *testing.T) {
	raw := strings.NewReplacer(`"maxFeePerGas": "0x4a817c800",`, ``).Replace(nethermindFrameTx(t))

	var fields JSONTxFields
	if err := json.Unmarshal([]byte(raw), &fields); err != nil {
		t.Fatalf("failed reading fields: %v", err)
	}

	fields.Raw = json.RawMessage(raw)

	if err := new(FrameTx).DecodeJSONTx(&fields); err == nil {
		t.Fatal("decoded a frame transaction without fees")
	}
}
