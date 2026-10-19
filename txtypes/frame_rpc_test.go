package txtypes

import (
	"encoding/json"
	"os"
	"testing"
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
