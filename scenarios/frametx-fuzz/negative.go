package frametxfuzz

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"

	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txtypes"
)

// BurnerWalletName is the wallet deliberately invalid transactions are fired from. They
// are never mined, so they consume no nonce and one wallet can send them indefinitely.
const BurnerWalletName = "frametx-fuzz-burner"

// violation is a deliberate corruption of a well-formed transaction.
//
// Structural violations are applied before signing so the signature covers them, or every
// case would fail on the signature rather than on the thing being exercised. Byte-level
// corruption is applied to the encoded transaction.
type violation struct {
	name string

	// alwaysInvalid marks a mutation that is invalid regardless of chain state, so
	// acceptance is a genuine finding rather than an expected outcome.
	alwaysInvalid bool

	// apply corrupts the transaction before it is signed.
	apply func(tx *txtypes.FrameTx) error

	// corrupt corrupts the encoded bytes after signing.
	corrupt func(raw []byte) []byte

	// rewrite re-encodes a signed transaction after altering something the signature
	// does not cover, such as a signature's own bytes.
	rewrite func(tx *txtypes.Transaction, raw []byte) []byte
}

// violations is the catalog a recipe draws from. Each is expected to be refused; the
// ones marked alwaysInvalid are refused by the payload format itself rather than by
// chain state or by mempool policy.
var violations = []violation{
	{
		// A VERIFY frame that approves nothing leaves the transaction without a
		// validation prefix, so no frame ever sets the payer.
		name:          "no-payment-approval",
		alwaysInvalid: false,
		apply: func(tx *txtypes.FrameTx) error {
			for _, frame := range tx.Frames {
				if frame.Mode == txtypes.FrameModeVerify {
					frame.Flags &^= txtypes.ApproveScopeMask
				}
			}

			return nil
		},
	},
	{
		// Approving execution commits the sender's account, so only the sender may be
		// the target of a frame that does it.
		name:          "approve-execution-for-foreign-target",
		alwaysInvalid: true,
		apply: func(tx *txtypes.FrameTx) error {
			foreign := common.HexToAddress("0x00000000000000000000000000000000deadbeef")

			for _, frame := range tx.Frames {
				if frame.Flags&txtypes.ApproveExecution != 0 {
					frame.Target = &foreign

					return nil
				}
			}

			return errViolationNotApplicable
		},
	},
	{
		// No frame in the validation prefix may be batched.
		name:          "batched-prefix-frame",
		alwaysInvalid: false,
		apply: func(tx *txtypes.FrameTx) error {
			tx.Frames[0].Flags |= txtypes.AtomicBatchFlag

			return nil
		},
	},
	{
		name:          "truncated",
		alwaysInvalid: true,
		corrupt: func(raw []byte) []byte {
			if len(raw) < 8 {
				return raw
			}

			return raw[:len(raw)-4]
		},
	},
	{
		name:          "corrupt-rlp",
		alwaysInvalid: true,
		corrupt: func(raw []byte) []byte {
			corrupted := append([]byte{}, raw...)
			corrupted[len(corrupted)/2] ^= 0xff

			return corrupted
		},
	},
	{
		name:          "too-many-frames",
		alwaysInvalid: true,
		apply: func(tx *txtypes.FrameTx) error {
			for len(tx.Frames) <= txtypes.MaxFrames {
				tx.Frames = append(tx.Frames, txtypes.UserOpFrame(nil, nil, nil, txtypes.FrameLimits{Execution: 1_000}))
			}

			return nil
		},
	},
	{
		name:          "value-outside-sender",
		alwaysInvalid: true,
		apply: func(tx *txtypes.FrameTx) error {
			frame := txtypes.PostOpFrame(tx.Sender, nil, txtypes.FrameLimits{Execution: 21_000})
			frame.Value = uint256.NewInt(1)
			tx.Frames = append(tx.Frames, frame)

			return nil
		},
	},
	{
		name:          "reserved-mode",
		alwaysInvalid: true,
		apply: func(tx *txtypes.FrameTx) error {
			frame := txtypes.UserOpFrame(nil, nil, nil, txtypes.FrameLimits{Execution: 21_000})
			frame.Mode = txtypes.FrameModePostTx + 1
			tx.Frames = append(tx.Frames, frame)

			return nil
		},
	},
	{
		name:          "reserved-flags",
		alwaysInvalid: true,
		apply: func(tx *txtypes.FrameTx) error {
			frame := txtypes.UserOpFrame(nil, nil, nil, txtypes.FrameLimits{Execution: 21_000})
			frame.Flags = txtypes.FrameFlagsMask + 1
			tx.Frames = append(tx.Frames, frame)

			return nil
		},
	},
	{
		name:          "batch-on-last-frame",
		alwaysInvalid: true,
		apply: func(tx *txtypes.FrameTx) error {
			tx.Frames[len(tx.Frames)-1].WithAtomicBatch()

			return nil
		},
	},
	{
		// The expiry frame's placement is a mempool rule rather than a payload
		// constraint: such a transaction is refused for propagation but is perfectly
		// valid inside a block.
		name:          "expiry-not-first",
		alwaysInvalid: false,
		apply: func(tx *txtypes.FrameTx) error {
			tx.Frames = append(tx.Frames, txtypes.ExpiryFrame(1<<40, 5_000))

			return nil
		},
	},
	{
		// Also mempool policy: a VERIFY frame after the prefix could invalidate the
		// whole transaction on revert, which the public mempool refuses to carry, but
		// nothing in the payload format forbids it.
		name:          "verify-after-prefix",
		alwaysInvalid: false,
		apply: func(tx *txtypes.FrameTx) error {
			tx.Frames = append(tx.Frames, txtypes.SelfVerifyFrame(txtypes.FrameLimits{Execution: 5_000}))

			return nil
		},
	},
	{
		name:          "post-tx-not-suffix",
		alwaysInvalid: true,
		apply: func(tx *txtypes.FrameTx) error {
			tx.Frames = append(tx.Frames,
				txtypes.PostTxFrame(txtypes.ExpiryVerifier, expiryData(1<<40), txtypes.FrameLimits{Execution: 5_000}),
				txtypes.UserOpFrame(nil, nil, nil, txtypes.FrameLimits{Execution: 21_000}),
			)

			return nil
		},
	},
	{
		name:          "prefix-over-verify-gas",
		alwaysInvalid: false, // mempool policy, not block validity
		apply: func(tx *txtypes.FrameTx) error {
			tx.Frames[0].Limits.Execution = txtypes.MaxVerifyGas + 1

			return nil
		},
	},
	{
		name:          "too-many-nonce-keys",
		alwaysInvalid: true,
		apply: func(tx *txtypes.FrameTx) error {
			if !tx.HasKeyedNonces() {
				return errViolationNotApplicable
			}

			keys := make([]*uint256.Int, 0, txtypes.MaxNonceKeys+1)
			for i := 0; i <= txtypes.MaxNonceKeys; i++ {
				keys = append(keys, uint256.NewInt(uint64(i)+1))
			}

			tx.NonceKeys = keys

			return nil
		},
	},
	{
		name:          "unsorted-nonce-keys",
		alwaysInvalid: true,
		apply: func(tx *txtypes.FrameTx) error {
			if !tx.HasKeyedNonces() {
				return errViolationNotApplicable
			}

			tx.NonceKeys = []*uint256.Int{uint256.NewInt(9), uint256.NewInt(7)}

			return nil
		},
	},
	{
		name:          "zero-key-mixed",
		alwaysInvalid: true,
		apply: func(tx *txtypes.FrameTx) error {
			if !tx.HasKeyedNonces() {
				return errViolationNotApplicable
			}

			tx.NonceKeys = []*uint256.Int{uint256.NewInt(0), uint256.NewInt(7)}

			return nil
		},
	},
	{
		name:          "empty-nonce-keys",
		alwaysInvalid: true,
		apply: func(tx *txtypes.FrameTx) error {
			if !tx.HasKeyedNonces() {
				return errViolationNotApplicable
			}

			tx.NonceKeys = nil

			return nil
		},
	},
	{
		name:          "too-many-recent-roots",
		alwaysInvalid: true,
		apply: func(tx *txtypes.FrameTx) error {
			if !tx.Extensions.Has(txtypes.FrameExtRecentRoots) {
				return errViolationNotApplicable
			}

			references := make([]*txtypes.RecentRootReference, 0, txtypes.MaxRecentRootReferences+1)
			for i := 0; i <= txtypes.MaxRecentRootReferences; i++ {
				references = append(references, &txtypes.RecentRootReference{Slot: uint64(i)})
			}

			tx.RecentRoots = references

			return nil
		},
	},
	{
		name:          "blob-fee-without-blobs",
		alwaysInvalid: true,
		apply: func(tx *txtypes.FrameTx) error {
			tx.Fees.BlobFeeCap = uint256.NewInt(1)

			return nil
		},
	},
	{
		name:          "signer-on-arbitrary-entry",
		alwaysInvalid: true,
		apply: func(tx *txtypes.FrameTx) error {
			entry := txtypes.ArbitrarySignature([]byte("witness"))
			entry.Signer = tx.Sender.Bytes()
			tx.Signatures = append(tx.Signatures, entry)

			return nil
		},
	},
	{
		name:          "zero-explicit-digest",
		alwaysInvalid: true,
		apply: func(tx *txtypes.FrameTx) error {
			entry := txtypes.SenderSignature()
			entry.Msg = make([]byte, 32)
			tx.Signatures = append(tx.Signatures, entry)

			return nil
		},
	},
	{
		// Flipping s to its high form keeps the signature recoverable but
		// non-canonical, which the protocol excludes so that each signature has
		// exactly one valid encoding.
		name:          "high-s-signature",
		alwaysInvalid: true,
		rewrite:       flipSignatureS,
	},
	{
		// A P256 entry whose public key is left intact but whose signature no longer
		// verifies. The distinction matters: an entry with a corrupted key fails the
		// format check that the signer matches its public key, which says nothing
		// about whether the signature itself was ever checked.
		name:          "p256-signature-does-not-verify",
		alwaysInvalid: true,
		rewrite:       corruptP256Signature,
	},
}

// errViolationNotApplicable reports that a mutation cannot be expressed on this chain's
// envelope, so the transaction is sent unmutated rather than sent wrong.
var errViolationNotApplicable = fmt.Errorf("violation not applicable to this envelope")

// violationNames lists the catalog for the draw and for help text.
func violationNames() []string {
	names := make([]string, len(violations))
	for i := range violations {
		names[i] = violations[i].name
	}

	return names
}

// findViolation returns a violation by name.
func findViolation(name string) *violation {
	for i := range violations {
		if violations[i].name == name {
			return &violations[i]
		}
	}

	return nil
}

// sendInvalid submits a transaction carrying a deliberate violation, out of pool: the
// RPC's accept or reject is the whole result. What happened is recorded, not judged.
func (s *Scenario) sendInvalid(ctx context.Context, client *spamoor.Client, result *build) error {
	violator := findViolation(result.recipe.Invalid)
	if violator == nil {
		return fmt.Errorf("unknown violation %q", result.recipe.Invalid)
	}

	// The nonce comes from the chain, since an invalid transaction never lands and the
	// pool's sequence would run away from the account's.
	nonce, err := client.GetPendingNonceAt(ctx, result.sender.GetAddress())
	if err != nil {
		return err
	}

	result.tx.NonceSeq = nonce

	if violator.apply != nil {
		if err := violator.apply(result.tx); err != nil {
			if errors.Is(err, errViolationNotApplicable) {
				return nil
			}

			return err
		}
	}

	tx, err := result.sender.SignFrameTx(result.tx)
	if err != nil {
		return err
	}

	raw, err := tx.MarshalNetwork()
	if err != nil {
		return err
	}

	if violator.rewrite != nil {
		raw = violator.rewrite(tx, raw)
	}

	if violator.corrupt != nil {
		raw = violator.corrupt(raw)
	}

	if err := client.SendRawTransaction(ctx, raw); err != nil {
		s.coverage.refusedOne(result.recipe, err.Error())
		s.coverage.invalidSubmitted(result.recipe, false)

		return nil
	}

	s.coverage.invalidSubmitted(result.recipe, true)

	return nil
}

// flipSignatureS rewrites the sender entry's s value to its high form, which keeps the
// signature recoverable but non-canonical.
func flipSignatureS(tx *txtypes.Transaction, raw []byte) []byte {
	frameTx, ok := tx.Inner().(*txtypes.FrameTx)
	if !ok || len(frameTx.Signatures) == 0 || len(frameTx.Signatures[0].Signature) != 65 {
		return raw
	}

	original := common.CopyBytes(frameTx.Signatures[0].Signature)

	order := new(big.Int).SetBytes(common.FromHex("0xfffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141"))
	s := new(big.Int).SetBytes(original[33:65])
	high := new(big.Int).Sub(order, s)

	flipped := common.CopyBytes(original)
	high.FillBytes(flipped[33:65])

	frameTx.Signatures[0].Signature = flipped

	mutated, err := tx.MarshalNetwork()
	frameTx.Signatures[0].Signature = original

	if err != nil {
		return raw
	}

	return mutated
}

// corruptP256Signature flips a bit of a P256 entry's r value, leaving the public key in
// place so the entry still resolves to the signer it names.
func corruptP256Signature(tx *txtypes.Transaction, raw []byte) []byte {
	frameTx, ok := tx.Inner().(*txtypes.FrameTx)
	if !ok {
		return raw
	}

	for _, sig := range frameTx.Signatures {
		if sig.Scheme != txtypes.SigSchemeP256 || len(sig.Signature) != 128 {
			continue
		}

		original := common.CopyBytes(sig.Signature)

		corrupted := common.CopyBytes(original)
		corrupted[31] ^= 0x01

		sig.Signature = corrupted

		mutated, err := tx.MarshalNetwork()
		sig.Signature = original

		if err != nil {
			return raw
		}

		return mutated
	}

	return raw
}
