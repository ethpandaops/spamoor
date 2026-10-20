package txtypes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// These are responses for the same confirmed transaction on bogota-devnet-1,
// block 28508. Check the reconstructed hash, not the cached RPC hash: a cached
// hash hides missing nonce, gas and signature fields.
func TestFrameRPCClientResponses(t *testing.T) {
	for _, client := range []string{"geth", "ethrex", "nethermind", "reth"} {
		t.Run(client, func(t *testing.T) {
			raw, err := os.ReadFile("testdata/frame-rpc/" + client + ".json")
			if err != nil {
				t.Fatal(err)
			}
			tx, err := UnmarshalJSONTx(raw)
			if err != nil {
				t.Fatal(err)
			}
			frame, ok := tx.Inner().(*FrameTx)
			if !ok {
				t.Fatalf("frame decoded as %T", tx.Inner())
			}
			if tx.Nonce() != 70 {
				t.Errorf("nonce = %d, want 70", tx.Nonce())
			}
			if !tx.UsesAccountNonce() {
				t.Error("key zero must track the account nonce on this devnet")
			}
			if frame.Frames[0].Limits.Execution != 5000 || frame.Frames[1].Limits.Execution != 30000 {
				t.Error("frame execution gas was lost")
			}
			if got := NewTx(frame).Hash(); got != tx.Hash() {
				t.Errorf("reconstructed hash = %s, RPC hash = %s", got, tx.Hash())
			}
		})
	}
}

func TestFrameRPCNonceSeqPrecedence(t *testing.T) {
	raw, err := os.ReadFile("testdata/frame-rpc/ethrex.json")
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]any
	if err := json.Unmarshal(raw, &fields); err != nil {
		t.Fatal(err)
	}
	// An explicit zero sequence must not fall back to the generic nonce.
	fields["nonceSeq"] = "0x0"
	raw, err = json.Marshal(fields)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := UnmarshalJSONTx(raw)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := tx.Inner().(*FrameTx); !ok {
		t.Fatalf("frame decoded as %T", tx.Inner())
	}
	if tx.Nonce() != 0 {
		t.Fatalf("explicit sequence ignored: %d", tx.Nonce())
	}
}

func TestFrameRPCSignerLengths(t *testing.T) {
	fixture, err := os.ReadFile("testdata/frame-rpc/geth.json")
	if err != nil {
		t.Fatal(err)
	}

	for _, size := range []int{0, 1, 19, 20, 21, 32} {
		t.Run(fmt.Sprintf("%d_bytes", size), func(t *testing.T) {
			var object map[string]any
			if err := json.Unmarshal(fixture, &object); err != nil {
				t.Fatal(err)
			}
			signer := bytes.Repeat([]byte{0x11}, size)
			object["signatures"].([]any)[0].(map[string]any)["signer"] = hexutil.Encode(signer)
			raw, err := json.Marshal(object)
			if err != nil {
				t.Fatal(err)
			}
			var fields JSONTxFields
			if err := json.Unmarshal(raw, &fields); err != nil {
				t.Fatal(err)
			}
			fields.Raw = raw
			var frame FrameTx
			err = frame.DecodeJSONTx(&fields)
			if size != 0 && size != 20 {
				if err == nil {
					t.Fatal("malformed signer length accepted by frame decoder")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(frame.Signatures[0].Signer, signer) {
				t.Fatal("signer bytes changed during decoding")
			}

			encoded := make(map[string]any)
			frame.EncodeJSONTx(encoded)
			raw, err = json.Marshal(encoded)
			if err != nil {
				t.Fatal(err)
			}
			var roundTripFields JSONTxFields
			if err := json.Unmarshal(raw, &roundTripFields); err != nil {
				t.Fatal(err)
			}
			roundTripFields.Raw = raw
			var roundTrip FrameTx
			if err := roundTrip.DecodeJSONTx(&roundTripFields); err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(roundTrip.Signatures[0].Signer, signer) {
				t.Fatal("signer bytes changed across JSON round-trip")
			}
			if NewTx(&frame).Hash() != NewTx(&roundTrip).Hash() {
				t.Fatal("transaction hash changed across JSON round-trip")
			}
		})
	}
}
