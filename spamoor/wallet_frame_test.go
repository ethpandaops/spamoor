package spamoor

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"

	"github.com/ethpandaops/spamoor/txtypes"
)

// A replacement has to raise max_fee, which a budget carried over from the replaced
// transaction would not, and has to be signed over the new signature hash.
func TestReplaceFrameTxRederivesMaxFee(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed generating key: %v", err)
	}

	wallet := NewWallet(key, crypto.PubkeyToAddress(key.PublicKey))
	wallet.SetChainId(big.NewInt(1))

	txData := txtypes.NewFrameTx(nil, wallet.GetAddress(), 0,
		txtypes.FrameFees{GasTipCap: uint256.NewInt(1e9), GasFeeCap: uint256.NewInt(10e9)},
		[]*txtypes.Frame{txtypes.SelfVerifyFrame(txtypes.FrameLimits{Execution: 5_000})},
		[]*txtypes.FrameSignature{txtypes.SenderSignature()})

	first, err := wallet.ReplaceFrameTx(txData, 3)
	if err != nil {
		t.Fatalf("failed building: %v", err)
	}

	signed := first.Inner().(*txtypes.FrameTx)

	txData = signed.CopyTx().(*txtypes.FrameTx)
	txData.Fees.GasFeeCap = uint256.NewInt(20e9)

	replacement, err := wallet.ReplaceFrameTx(txData, 3)
	if err != nil {
		t.Fatalf("failed replacing: %v", err)
	}

	replaced := replacement.Inner().(*txtypes.FrameTx)

	if replaced.GetMaxFee().Cmp(new(big.Int).Mul(signed.GetMaxFee(), big.NewInt(2))) != 0 {
		t.Fatalf("replacement max fee = %s, want twice %s", replaced.GetMaxFee(), signed.GetMaxFee())
	}

	if err := replaced.VerifySignatures(); err != nil {
		t.Fatalf("replacement is not validly signed: %v", err)
	}
}
