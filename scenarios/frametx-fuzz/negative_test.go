package frametxfuzz

import (
	"crypto/ecdsa"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"

	"github.com/ethpandaops/spamoor/txtypes"
)

// testKey signs the transactions the negative tests build, and its address is their
// sender so the default sender entry resolves to it.
var testKey, _ = crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")

// testChainID is the chain the negative tests sign for.
var testChainID = big.NewInt(1)

// validFrameTx builds the well-formed transaction every violation starts from.
func validFrameTx() *txtypes.FrameTx {
	sender := crypto.PubkeyToAddress(testKey.PublicKey)

	return txtypes.NewFrameTxWithExtensions(txtypes.FrameExtAll, uint256.NewInt(1), sender, 0,
		txtypes.FrameFees{GasFeeCap: uint256.NewInt(1e9)},
		[]*txtypes.Frame{
			txtypes.SelfVerifyFrame(txtypes.FrameLimits{Execution: 5_000}),
			txtypes.UserOpFrame(nil, nil, nil, txtypes.FrameLimits{Execution: 21_000}),
		},
		[]*txtypes.FrameSignature{txtypes.SenderSignature()},
	)
}

// signTestTx signs a transaction the way the scenario does, after any violation has been
// applied: a structural violation has to be covered by the signature, or every case would
// be refused for a bad signature instead of for the thing being tested.
func signTestTx(t *testing.T, tx *txtypes.FrameTx, key *ecdsa.PrivateKey) {
	t.Helper()

	if err := tx.SignPayload(testChainID, key); err != nil {
		t.Fatalf("could not sign the transaction: %v", err)
	}
}

// TestBaselineIsValid guards the premise of the whole negative tier: the transaction the
// violations start from must itself be acceptable, or every case would be "rejected" for
// the wrong reason.
func TestBaselineIsValid(t *testing.T) {
	tx := validFrameTx()
	signTestTx(t, tx, testKey)

	if err := tx.ValidatePayload(); err != nil {
		t.Fatalf("the unmutated transaction is invalid: %v", err)
	}

	if err := tx.VerifySignatures(); err != nil {
		t.Fatalf("the unmutated transaction does not verify: %v", err)
	}

	if err := tx.ValidateMempoolPrefix(); err != nil {
		t.Fatalf("the unmutated transaction would not propagate: %v", err)
	}
}

// TestViolationsProduceInvalidTransactions checks that every structural violation actually
// produces something a correct node must refuse.
//
// This is the guard that keeps the fuzzer honest: a violation that quietly stops being
// invalid would have its acceptance reported as a client finding, which is worse than
// having no case at all.
func TestViolationsProduceInvalidTransactions(t *testing.T) {
	for i := range violations {
		violator := &violations[i]

		if violator.apply == nil {
			continue
		}

		t.Run(violator.name, func(t *testing.T) {
			tx := validFrameTx()

			if err := violator.apply(tx); err != nil {
				if errors.Is(err, errViolationNotApplicable) {
					t.Skip("not applicable to the full envelope")
				}

				t.Fatalf("violation failed to apply: %v", err)
			}

			signTestTx(t, tx, testKey)

			payloadErr := tx.ValidatePayload()
			mempoolErr := tx.ValidateMempoolPrefix()

			if violator.alwaysInvalid {
				if payloadErr == nil {
					t.Errorf("violation is marked as always invalid but the payload validates")
				}

				return
			}

			// The remaining violations are mempool policy rather than block validity, so
			// they must be refused for propagation without being malformed.
			if mempoolErr == nil {
				t.Errorf("violation is marked as a policy violation but the transaction would propagate")
			}
		})
	}
}

// TestByteViolationsChangeTheEncoding checks that the byte-level cases actually alter the
// transaction, so a no-op corruption is not mistaken for a client accepting a malformed
// payload.
func TestByteViolationsChangeTheEncoding(t *testing.T) {
	raw := []byte{0x06, 0xf8, 0x40, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}

	for i := range violations {
		violator := &violations[i]

		if violator.corrupt == nil {
			continue
		}

		t.Run(violator.name, func(t *testing.T) {
			corrupted := violator.corrupt(raw)

			if string(corrupted) == string(raw) {
				t.Error("corruption left the encoding unchanged")
			}
		})
	}
}
