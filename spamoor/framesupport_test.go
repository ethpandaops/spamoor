package spamoor

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"

	"github.com/ethpandaops/spamoor/txtypes"
)

// TestApplyKeyedNonceRevision checks that a frame transaction decoded from a block is read
// in the chain's EIP-8250 revision: on a key-zero chain a decoded [0] is the account nonce
// and must take the nonce-tracked confirmation path, and the hash must not change.
func TestApplyKeyedNonceRevision(t *testing.T) {
	sender := common.HexToAddress("0x8943545177806ed17b9f23f0a21ee5948ecaa776")

	built := txtypes.NewFrameTxWithExtensions(txtypes.FrameExtKeyedNonces|txtypes.FrameExtLegacyNonceKeyZero,
		uint256.NewInt(1), sender, 7, txtypes.FrameFees{GasFeeCap: uint256.NewInt(1)},
		[]*txtypes.Frame{txtypes.SelfVerifyFrame(txtypes.FrameLimits{Execution: 5_000})},
		[]*txtypes.FrameSignature{{Scheme: txtypes.SigSchemeSecp256k1, Signature: make([]byte, 65)}})

	if !built.UsesAccountNonce() {
		t.Fatal("the built transaction should use the account nonce")
	}

	encoded, err := txtypes.NewTx(built).MarshalBinary()
	if err != nil {
		t.Fatalf("failed encoding: %v", err)
	}

	decoded, err := txtypes.DecodeTx(encoded)
	if err != nil {
		t.Fatalf("failed decoding: %v", err)
	}

	if decoded.UsesAccountNonce() {
		t.Fatal("a decoded [0] should read as key zero before the revision is applied")
	}

	previous := KeyedNonceRevision
	defer func() { KeyedNonceRevision = previous }()

	KeyedNonceRevision = txtypes.FrameExtLegacyNonceKeyZero
	ApplyKeyedNonceRevision(decoded)

	if !decoded.UsesAccountNonce() {
		t.Fatal("the key-zero revision should make a decoded [0] the account nonce")
	}

	if decoded.Hash() != txtypes.NewTx(built).Hash() {
		t.Fatal("applying the revision changed the hash")
	}

	// Under the current revision nothing is stamped and [0] stays key zero.
	spec, err := txtypes.DecodeTx(encoded)
	if err != nil {
		t.Fatalf("failed decoding: %v", err)
	}

	KeyedNonceRevision = 0
	ApplyKeyedNonceRevision(spec)

	if spec.UsesAccountNonce() {
		t.Fatal("the current revision reads [0] as key zero")
	}
}
