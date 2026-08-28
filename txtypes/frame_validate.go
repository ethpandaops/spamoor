package txtypes

import (
	"crypto/elliptic"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// FrameSpecies classifies a frame for public mempool validation.
type FrameSpecies string

// Frame species recognized by the public mempool rules.
const (
	SpeciesSelfVerify   FrameSpecies = "self_verify"   // VERIFY, flags 0x3
	SpeciesOnlyVerify   FrameSpecies = "only_verify"   // VERIFY, flags 0x2
	SpeciesPay          FrameSpecies = "pay"           // VERIFY, flags 0x1
	SpeciesExpiryVerify FrameSpecies = "expiry_verify" // VERIFY, flags 0x0, target EXPIRY_VERIFIER
	SpeciesDeploy       FrameSpecies = "deploy"        // DEFAULT, flags 0x0
	SpeciesUserOp       FrameSpecies = "user_op"       // SENDER
	SpeciesPostOp       FrameSpecies = "post_op"       // DEFAULT with flags
	SpeciesPostTx       FrameSpecies = "post_tx"       // POST_TX (EIP-7906)
	SpeciesOther        FrameSpecies = "other"
)

// Species classifies the frame for mempool prefix matching.
func (f *Frame) Species(sender common.Address) FrameSpecies {
	switch f.Mode {
	case FrameModeVerify:
		if f.IsExpiryVerifier() {
			return SpeciesExpiryVerify
		}

		switch f.Flags {
		case ApproveExecutionAndPayment:
			return SpeciesSelfVerify
		case ApproveExecution:
			return SpeciesOnlyVerify
		case ApprovePayment:
			return SpeciesPay
		}

		return SpeciesOther

	case FrameModeDefault:
		if f.Flags == 0 {
			return SpeciesDeploy
		}

		return SpeciesPostOp

	case FrameModeSender:
		return SpeciesUserOp

	case FrameModePostTx:
		return SpeciesPostTx
	}

	return SpeciesOther
}

// ValidatePayload checks the constraints that are decidable from the transaction
// fields alone, mirroring the static assertions in EIP-8141.
func (tx *FrameTx) ValidatePayload() error {
	if len(tx.Frames) == 0 || len(tx.Frames) > MaxFrames {
		return fmt.Errorf("%w: frame count %d must be between 1 and %d", ErrInvalidFrameTx, len(tx.Frames), MaxFrames)
	}

	if len(tx.BlobHashes) == 0 && tx.Fees.BlobFeeCap != nil && !tx.Fees.BlobFeeCap.IsZero() {
		return fmt.Errorf("%w: blob fee cap must be zero without blobs", ErrInvalidFrameTx)
	}

	for i, hash := range tx.BlobHashes {
		if hash[0] != VersionedHashVersionKZG {
			return fmt.Errorf("%w: blob hash %d has wrong version byte 0x%02x", ErrInvalidFrameTx, i, hash[0])
		}
	}

	if err := tx.validateNonce(); err != nil {
		return err
	}

	if !tx.Extensions.Has(FrameExtRecentRoots) && len(tx.RecentRoots) > 0 {
		return fmt.Errorf("%w: recent root references set without the EIP-8272 extension", ErrInvalidFrameTx)
	}

	if len(tx.RecentRoots) > MaxRecentRootReferences {
		return fmt.Errorf("%w: %d recent root references exceeds the cap of %d",
			ErrInvalidFrameTx, len(tx.RecentRoots), MaxRecentRootReferences)
	}

	if err := tx.validateSignatures(); err != nil {
		return err
	}

	if err := tx.validateExpiryFrames(); err != nil {
		return err
	}

	return tx.validateFrames()
}

// validateExpiryFrames checks the expiry verifier frame's own validity rules. A
// transaction may carry at most one, and it must carry nothing but the deadline.
func (tx *FrameTx) validateExpiryFrames() error {
	seen := false

	for i, frame := range tx.Frames {
		if !frame.IsExpiryVerifier() {
			continue
		}

		if seen {
			return fmt.Errorf("%w: more than one expiry verifier frame", ErrInvalidFrameTx)
		}

		seen = true

		if frame.Flags != 0 {
			return fmt.Errorf("%w: expiry verifier frame %d must have zero flags", ErrInvalidFrameTx, i)
		}

		if frame.Value != nil && !frame.Value.IsZero() {
			return fmt.Errorf("%w: expiry verifier frame %d must carry no value", ErrInvalidFrameTx, i)
		}

		if frame.Limits.State != 0 {
			return fmt.Errorf("%w: expiry verifier frame %d must budget no state gas", ErrInvalidFrameTx, i)
		}

		if len(frame.Data) != ExpiryDataLength {
			return fmt.Errorf("%w: expiry verifier frame %d data must be %d bytes, got %d",
				ErrInvalidFrameTx, i, ExpiryDataLength, len(frame.Data))
		}
	}

	return nil
}

// validateNonce checks the EIP-8250 keyed nonce constraints. The key set must be
// non-empty, bounded, strictly increasing, and the zero key may only appear alone
// because it aliases the sender's ordinary account nonce.
func (tx *FrameTx) validateNonce() error {
	if !tx.HasKeyedNonces() {
		// EIP-8141's scalar nonce: there are no keys to check, and carrying any
		// would not survive encoding.
		if len(tx.NonceKeys) > 0 {
			return fmt.Errorf("%w: nonce keys set without the EIP-8250 extension", ErrInvalidFrameTx)
		}

		return nil
	}

	if len(tx.NonceKeys) < 1 || len(tx.NonceKeys) > MaxNonceKeys {
		return fmt.Errorf("%w: nonce key count %d must be between 1 and %d",
			ErrInvalidFrameTx, len(tx.NonceKeys), MaxNonceKeys)
	}

	for i, key := range tx.NonceKeys {
		if key == nil {
			return fmt.Errorf("%w: nonce key %d is nil", ErrInvalidFrameTx, i)
		}

		if i > 0 && key.Cmp(tx.NonceKeys[i-1]) <= 0 {
			return fmt.Errorf("%w: nonce keys must be strictly increasing", ErrInvalidFrameTx)
		}

		if key.IsZero() && len(tx.NonceKeys) != 1 {
			return fmt.Errorf("%w: the zero nonce key must be used alone", ErrInvalidFrameTx)
		}
	}

	return nil
}

// validateSignatures checks the signature list's structural constraints.
func (tx *FrameTx) validateSignatures() error {
	for i, sig := range tx.Signatures {
		if _, err := sig.VerificationGas(); err != nil {
			return err
		}

		switch sig.Scheme {
		case SigSchemeSecp256k1, SigSchemeP256:
			if len(sig.Signer) != 0 && len(sig.Signer) != 20 {
				return fmt.Errorf("%w: signature %d signer must be empty or 20 bytes", ErrInvalidFrameTx, i)
			}
		case SigSchemeArbitrary:
			if len(sig.Signer) != 0 {
				return fmt.Errorf("%w: signature %d is ARBITRARY and must have an empty signer", ErrInvalidFrameTx, i)
			}
		}

		switch len(sig.Msg) {
		case 0:
		case 32:
			if common.BytesToHash(sig.Msg) == (common.Hash{}) {
				return fmt.Errorf("%w: signature %d has a zero explicit digest", ErrInvalidFrameTx, i)
			}
		default:
			return fmt.Errorf("%w: signature %d msg must be empty or 32 bytes", ErrInvalidFrameTx, i)
		}

		if err := tx.validateSignatureBytes(i, sig); err != nil {
			return err
		}
	}

	return nil
}

// validateSignatureBytes checks a protocol-validated entry's raw signature encoding.
// The curve checks reject the malleable encodings the protocol excludes, so each
// signature has exactly one valid form.
func (tx *FrameTx) validateSignatureBytes(index int, sig *FrameSignature) error {
	switch sig.Scheme {
	case SigSchemeSecp256k1:
		if len(sig.Signature) != 65 {
			return fmt.Errorf("%w: signature %d must be 65 bytes, got %d", ErrInvalidFrameTx, index, len(sig.Signature))
		}

		if sig.Signature[0] > 1 {
			return fmt.Errorf("%w: signature %d has recovery id %d, must be 0 or 1", ErrInvalidFrameTx, index, sig.Signature[0])
		}

		r := new(big.Int).SetBytes(sig.Signature[1:33])
		s := new(big.Int).SetBytes(sig.Signature[33:65])

		return validateCurveValues(index, r, s, crypto.S256().Params().N)

	case SigSchemeP256:
		if len(sig.Signature) != 128 {
			return fmt.Errorf("%w: signature %d must be 128 bytes, got %d", ErrInvalidFrameTx, index, len(sig.Signature))
		}

		r := new(big.Int).SetBytes(sig.Signature[0:32])
		s := new(big.Int).SetBytes(sig.Signature[32:64])

		if err := validateCurveValues(index, r, s, elliptic.P256().Params().N); err != nil {
			return err
		}

		qx := new(big.Int).SetBytes(sig.Signature[64:96])
		qy := new(big.Int).SetBytes(sig.Signature[96:128])

		signer, _ := sig.ResolvedSigner(tx.Sender)
		if signer != P256Signer(qx, qy) {
			return fmt.Errorf("%w: signature %d signer does not match its public key", ErrInvalidFrameTx, index)
		}
	}

	return nil
}

// validateCurveValues checks that r and s are canonical with a low s.
func validateCurveValues(index int, r, s, order *big.Int) error {
	if r.Sign() <= 0 || r.Cmp(order) >= 0 {
		return fmt.Errorf("%w: signature %d has a non-canonical r value", ErrInvalidFrameTx, index)
	}

	halfOrder := new(big.Int).Rsh(order, 1)
	if s.Sign() <= 0 || s.Cmp(halfOrder) > 0 {
		return fmt.Errorf("%w: signature %d has a non-canonical or high s value", ErrInvalidFrameTx, index)
	}

	return nil
}

// validateFrames checks per-frame constraints and the aggregate gas limits.
func (tx *FrameTx) validateFrames() error {
	totalGas := uint64(0)

	for i, frame := range tx.Frames {
		if frame.Mode > FrameModePostTx {
			return fmt.Errorf("%w: frame %d has unknown mode %d", ErrInvalidFrameTx, i, frame.Mode)
		}

		if frame.Flags&^FrameFlagsMask != 0 {
			return fmt.Errorf("%w: frame %d sets reserved flag bits 0x%02x", ErrInvalidFrameTx, i, frame.Flags)
		}

		if frame.Mode != FrameModeSender && frame.Value != nil && !frame.Value.IsZero() {
			return fmt.Errorf("%w: frame %d carries value outside SENDER mode", ErrInvalidFrameTx, i)
		}

		nextGas := totalGas + frame.Limits.Execution + frame.Limits.State
		if nextGas < totalGas {
			return fmt.Errorf("%w: total frame gas overflows 64 bits", ErrInvalidFrameTx)
		}

		totalGas = nextGas

		// Approving execution commits the sender's account, so only the sender may
		// be the target.
		if frame.Flags&ApproveExecution != 0 && frame.Target != nil && *frame.Target != tx.Sender {
			return fmt.Errorf("%w: frame %d approves execution for a foreign target", ErrInvalidFrameTx, i)
		}

		if frame.IsAtomicBatch() {
			if frame.Mode == FrameModeVerify {
				return fmt.Errorf("%w: frame %d batches in VERIFY mode", ErrInvalidFrameTx, i)
			}

			if i+1 >= len(tx.Frames) {
				return fmt.Errorf("%w: frame %d batches but is the last frame", ErrInvalidFrameTx, i)
			}

			if tx.Frames[i+1].Mode == FrameModeVerify {
				return fmt.Errorf("%w: frame %d batches with a VERIFY frame", ErrInvalidFrameTx, i)
			}
		}

		inBatch := frame.IsAtomicBatch() || (i > 0 && tx.Frames[i-1].IsAtomicBatch())
		if inBatch && frame.ApprovalScope() != ApproveNone {
			return fmt.Errorf("%w: frame %d combines batching with an approval scope", ErrInvalidFrameTx, i)
		}
	}

	if err := tx.validatePostTxSuffix(); err != nil {
		return err
	}

	if executionGas := tx.ExecutionGas(); executionGas > TxMaxGasLimit {
		return fmt.Errorf("%w: execution gas %d exceeds the per-transaction cap %d",
			ErrInvalidFrameTx, executionGas, TxMaxGasLimit)
	}

	return nil
}

// validatePostTxSuffix checks EIP-7906's placement rule: once a frame has mode
// POST_TX, every later frame must too.
func (tx *FrameTx) validatePostTxSuffix() error {
	seen := false

	for i, frame := range tx.Frames {
		switch {
		case frame.Mode == FrameModePostTx:
			seen = true
		case seen:
			return fmt.Errorf("%w: frame %d follows a POST_TX frame but is mode %d; POST_TX frames must be a trailing suffix",
				ErrInvalidFrameTx, i, frame.Mode)
		}
	}

	return nil
}

// PostTxIndex returns the index of the first POST_TX frame, or -1 when the
// transaction has none.
func (tx *FrameTx) PostTxIndex() int {
	for i, frame := range tx.Frames {
		if frame.Mode == FrameModePostTx {
			return i
		}
	}

	return -1
}

// ValidationPrefixLength returns the number of leading frames that make up the
// validation prefix: the shortest prefix whose successful execution sets the payer.
// It returns 0 when no frame in the transaction approves payment.
func (tx *FrameTx) ValidationPrefixLength() int {
	for i, frame := range tx.Frames {
		if frame.ApprovalScope()&ApprovePayment != 0 {
			return i + 1
		}
	}

	return 0
}

// ValidateMempoolPrefix checks the public mempool policy of EIP-8141: the validation
// prefix must match one of the four recognized shapes and stay within the verification
// gas caps.
//
// A transaction failing this check may still be valid in a block; it just will not
// propagate through the public mempool.
func (tx *FrameTx) ValidateMempoolPrefix() error {
	prefixLen := tx.ValidationPrefixLength()
	if prefixLen == 0 {
		return fmt.Errorf("%w: no frame approves payment", ErrMempoolPolicy)
	}

	species := make([]FrameSpecies, prefixLen)
	for i := 0; i < prefixLen; i++ {
		species[i] = tx.Frames[i].Species(tx.Sender)
	}

	// An expiry verifier frame may lead the prefix and is skipped when matching.
	shape := species
	if shape[0] == SpeciesExpiryVerify {
		shape = shape[1:]
	}

	// The expiry frame may lead the frame list and nothing else.
	for i, frame := range tx.Frames {
		if frame.IsExpiryVerifier() && i != 0 {
			return fmt.Errorf("%w: expiry verifier frame must be the first frame", ErrMempoolPolicy)
		}
	}

	// A deploy frame leads the prefix, after a leading expiry frame if there is one.
	deployIndex := 0
	if len(species) > 0 && species[0] == SpeciesExpiryVerify {
		deployIndex = 1
	}

	if !matchesRecognizedPrefix(shape) {
		return fmt.Errorf("%w: validation prefix %v is not a recognized shape", ErrMempoolPolicy, shape)
	}

	for i := 0; i < prefixLen; i++ {
		frame := tx.Frames[i]

		if frame.IsAtomicBatch() {
			return fmt.Errorf("%w: frame %d in the validation prefix is batched", ErrMempoolPolicy, i)
		}

		switch species[i] {
		case SpeciesSelfVerify, SpeciesOnlyVerify:
			if frame.Target != nil && *frame.Target != tx.Sender {
				return fmt.Errorf("%w: frame %d must verify the sender", ErrMempoolPolicy, i)
			}
		case SpeciesDeploy:
			if i != deployIndex {
				return fmt.Errorf("%w: deploy frame must lead the validation prefix", ErrMempoolPolicy)
			}
		}
	}

	// No VERIFY frame may follow the prefix: its revert would invalidate the whole
	// transaction after payment was already approved.
	for i := prefixLen; i < len(tx.Frames); i++ {
		if tx.Frames[i].Mode == FrameModeVerify {
			return fmt.Errorf("%w: frame %d is a VERIFY frame after the validation prefix", ErrMempoolPolicy, i)
		}
	}

	return tx.validatePrefixGas(prefixLen)
}

// validatePrefixGas checks the verification gas caps over the validation prefix.
func (tx *FrameTx) validatePrefixGas(prefixLen int) error {
	sigGas, err := tx.SignatureVerificationGas()
	if err != nil {
		return err
	}

	executionGas := sigGas
	stateGas := uint64(0)

	for i := 0; i < prefixLen; i++ {
		executionGas += tx.Frames[i].Limits.Execution
		stateGas += tx.Frames[i].Limits.State
	}

	if executionGas > MaxVerifyGas {
		return fmt.Errorf("%w: validation prefix uses %d verification gas, cap is %d",
			ErrMempoolPolicy, executionGas, MaxVerifyGas)
	}

	if stateGas > MaxVerifyStateGas {
		return fmt.Errorf("%w: validation prefix budgets %d state gas, cap is %d",
			ErrMempoolPolicy, stateGas, MaxVerifyStateGas)
	}

	return nil
}

// recognizedPrefixes are the four validation prefix shapes the public mempool accepts.
var recognizedPrefixes = [][]FrameSpecies{
	{SpeciesSelfVerify},
	{SpeciesDeploy, SpeciesSelfVerify},
	{SpeciesOnlyVerify, SpeciesPay},
	{SpeciesDeploy, SpeciesOnlyVerify, SpeciesPay},
}

// matchesRecognizedPrefix reports whether shape is one of the recognized prefixes.
func matchesRecognizedPrefix(shape []FrameSpecies) bool {
	for _, candidate := range recognizedPrefixes {
		if len(candidate) != len(shape) {
			continue
		}

		match := true

		for i := range candidate {
			if candidate[i] != shape[i] {
				match = false

				break
			}
		}

		if match {
			return true
		}
	}

	return false
}

// VerifySignatures checks every protocol-validated signature entry against the digest
// it authorizes, which is what a client does before executing any frame.
//
// It is kept separate from ValidatePayload because it recovers public keys, while
// ValidatePayload only checks what is decidable from the encoded fields.
func (tx *FrameTx) VerifySignatures() error {
	sigHash := tx.SigHash()

	for i, sig := range tx.Signatures {
		if sig.Scheme != SigSchemeSecp256k1 {
			// P256 entries need the P256VERIFY precompile, and ARBITRARY entries are
			// not protocol-validated at all.
			continue
		}

		if err := tx.validateSignatureBytes(i, sig); err != nil {
			return err
		}

		digest := sigHash
		if len(sig.Msg) == 32 {
			digest = common.BytesToHash(sig.Msg)
		}

		// go-ethereum's recovery wants r || s || v, the entry carries v || r || s.
		raw := make([]byte, 65)
		copy(raw[:64], sig.Signature[1:])
		raw[64] = sig.Signature[0]

		pub, err := crypto.Ecrecover(digest[:], raw)
		if err != nil {
			return fmt.Errorf("%w: signature %d could not be recovered: %v", ErrInvalidFrameTx, i, err)
		}

		var recovered common.Address

		copy(recovered[:], crypto.Keccak256(pub[1:])[12:])

		expected, _ := sig.ResolvedSigner(tx.Sender)
		if recovered != expected {
			return fmt.Errorf("%w: signature %d recovered %s, expected %s",
				ErrInvalidFrameTx, i, recovered, expected)
		}
	}

	return nil
}
