package txtypes

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"math/big"
	"math/bits"

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
	SpeciesAssert       FrameSpecies = "assert"        // ASSERT (EIP-7906)
	SpeciesOther        FrameSpecies = "other"

	// SpeciesRecentRootVerify is EIP-8272's canonical frame: VERIFY, flags 0x0, target
	// RECENT_ROOT_ADDRESS, whole tuples as data.
	SpeciesRecentRootVerify FrameSpecies = "recent_root_verify"

	// SpeciesPostTx is the former name of SpeciesAssert.
	//
	// Deprecated: EIP-7906 renamed the mode to ASSERT; use SpeciesAssert.
	SpeciesPostTx = SpeciesAssert
)

// Species classifies the frame for mempool prefix matching.
func (f *Frame) Species(sender common.Address) FrameSpecies {
	switch f.Mode {
	case FrameModeVerify:
		if f.IsExpiryVerifier() {
			return SpeciesExpiryVerify
		}

		if f.IsRecentRootVerifier() {
			return SpeciesRecentRootVerify
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

	case FrameModeAssert:
		return SpeciesAssert
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

// validateNonce checks the EIP-8250 keyed nonce constraints: at most MaxNonceKeys keys,
// strictly increasing, and the single-use sequence only with a non-empty key set. Under
// FrameExtLegacyNonceKeyZero the older rules apply instead: at least one key, the zero
// key alone, and a sequence below MaxNonceSeq.
func (tx *FrameTx) validateNonce() error {
	if !tx.HasKeyedNonces() {
		// EIP-8141's scalar nonce: there are no keys to check, and carrying any
		// would not survive encoding.
		if len(tx.NonceKeys) > 0 {
			return fmt.Errorf("%w: nonce keys set without the EIP-8250 extension", ErrInvalidFrameTx)
		}

		return nil
	}

	if len(tx.NonceKeys) > MaxNonceKeys {
		return fmt.Errorf("%w: nonce key count %d exceeds %d", ErrInvalidFrameTx, len(tx.NonceKeys), MaxNonceKeys)
	}

	for i, key := range tx.NonceKeys {
		if key == nil {
			return fmt.Errorf("%w: nonce key %d is nil", ErrInvalidFrameTx, i)
		}

		if i > 0 && key.Cmp(tx.NonceKeys[i-1]) <= 0 {
			return fmt.Errorf("%w: nonce keys must be strictly increasing", ErrInvalidFrameTx)
		}
	}

	if tx.Extensions.Has(FrameExtLegacyNonceKeyZero) {
		if len(tx.NonceKeys) == 0 {
			return fmt.Errorf("%w: nonce key set is empty; this revision selects the account nonce with [0]", ErrInvalidFrameTx)
		}

		if tx.NonceKeys[0].IsZero() && len(tx.NonceKeys) != 1 {
			return fmt.Errorf("%w: the zero nonce key must be used alone in this revision", ErrInvalidFrameTx)
		}

		if tx.NonceSeq == MaxNonceSeq {
			return fmt.Errorf("%w: nonce sequence %d is out of range in this revision", ErrInvalidFrameTx, tx.NonceSeq)
		}

		return nil
	}

	if len(tx.NonceKeys) == 0 && tx.NonceSeq == MaxNonceSeq {
		return fmt.Errorf("%w: the single-use nonce sequence needs a non-empty key set", ErrInvalidFrameTx)
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
	// total_frame_gas is accumulated with explicit carry checks: the limits are
	// attacker-chosen 64-bit values, and a wrapped total would pass every cap below.
	totalGas := uint64(0)

	for i, frame := range tx.Frames {
		if frame.Mode > FrameModeAssert {
			return fmt.Errorf("%w: frame %d has unknown mode %d", ErrInvalidFrameTx, i, frame.Mode)
		}

		if frame.Flags&^FrameFlagsMask != 0 {
			return fmt.Errorf("%w: frame %d sets reserved flag bits 0x%02x", ErrInvalidFrameTx, i, frame.Flags)
		}

		if frame.Mode != FrameModeSender && frame.Value != nil && !frame.Value.IsZero() {
			return fmt.Errorf("%w: frame %d carries value outside SENDER mode", ErrInvalidFrameTx, i)
		}

		frameGas, carry := bits.Add64(frame.Limits.Execution, frame.Limits.State, 0)
		if carry != 0 {
			return fmt.Errorf("%w: frame %d gas limits overflow 64 bits", ErrInvalidFrameTx, i)
		}

		totalGas, carry = bits.Add64(totalGas, frameGas, 0)
		if carry != 0 {
			return fmt.Errorf("%w: total frame gas overflows 64 bits", ErrInvalidFrameTx)
		}

		// Approving execution commits the sender's account, so only the sender may
		// be the target.
		if frame.Flags&ApproveExecution != 0 && frame.Target != nil && *frame.Target != tx.Sender {
			return fmt.Errorf("%w: frame %d approves execution for a foreign target", ErrInvalidFrameTx, i)
		}

		if frame.IsAtomicBatch() {
			// The flag is valid with DEFAULT and SENDER only: a VERIFY frame's revert
			// invalidates the transaction and an ASSERT frame's revert discards the
			// whole body, so neither has anything a batch could roll back. The frame
			// it batches with has the same constraint, which also keeps a batch from
			// running into the ASSERT suffix.
			if !batchable(frame.Mode) {
				return fmt.Errorf("%w: frame %d batches in mode %d", ErrInvalidFrameTx, i, frame.Mode)
			}

			if i+1 >= len(tx.Frames) {
				return fmt.Errorf("%w: frame %d batches but is the last frame", ErrInvalidFrameTx, i)
			}

			if next := tx.Frames[i+1].Mode; !batchable(next) {
				return fmt.Errorf("%w: frame %d batches with a mode %d frame", ErrInvalidFrameTx, i, next)
			}
		}

		inBatch := frame.IsAtomicBatch() || (i > 0 && tx.Frames[i-1].IsAtomicBatch())
		if inBatch && frame.ApprovalScope() != ApproveNone {
			return fmt.Errorf("%w: frame %d combines batching with an approval scope", ErrInvalidFrameTx, i)
		}
	}

	if err := tx.validateAssertSuffix(); err != nil {
		return err
	}

	// ExecutionGas saturates rather than wraps, so an overflowing sum lands here.
	if executionGas := tx.ExecutionGas(); executionGas > TxMaxGasLimit {
		return fmt.Errorf("%w: execution gas %d exceeds the per-transaction cap %d",
			ErrInvalidFrameTx, executionGas, TxMaxGasLimit)
	}

	return nil
}

// batchable reports whether the atomic batch flag may sit on, or point at, a frame of
// the given mode.
func batchable(mode FrameMode) bool {
	return mode == FrameModeDefault || mode == FrameModeSender
}

// validateAssertSuffix checks EIP-7906's placement rule: once a frame has mode ASSERT,
// every later frame must too.
func (tx *FrameTx) validateAssertSuffix() error {
	seen := false

	for i, frame := range tx.Frames {
		switch {
		case frame.Mode == FrameModeAssert:
			seen = true
		case seen:
			return fmt.Errorf("%w: frame %d follows an ASSERT frame but is mode %d; ASSERT frames must be a trailing suffix",
				ErrInvalidFrameTx, i, frame.Mode)
		}
	}

	return nil
}

// AssertIndex returns the index of the first ASSERT frame, or -1 when the transaction
// has none.
func (tx *FrameTx) AssertIndex() int {
	for i, frame := range tx.Frames {
		if frame.Mode == FrameModeAssert {
			return i
		}
	}

	return -1
}

// PostTxIndex is the former name of AssertIndex.
//
// Deprecated: EIP-7906 renamed the mode to ASSERT; use AssertIndex.
func (tx *FrameTx) PostTxIndex() int { return tx.AssertIndex() }

// ValidationPrefixLength returns the number of leading frames that make up the
// validation prefix: the shortest prefix whose successful execution sets the payer.
// It returns 0 when no frame in the transaction approves payment.
//
// This is a static approximation. EIP-8141 defines the prefix at runtime as the frames
// up to the one after whose execution the payer is set, and a frame may carry a payment
// scope without approving, or approve and then revert, in which case a later frame ends
// the prefix. The first payment-scoped frame is where a well-formed transaction's prefix
// ends and is what the public mempool shapes assume.
func (tx *FrameTx) ValidationPrefixLength() int {
	for i, frame := range tx.Frames {
		if frame.ApprovalScope()&ApprovePayment != 0 {
			return i + 1
		}
	}

	return 0
}

// ValidateMempoolPrefix checks the public mempool policy of EIP-8141 and EIP-8272: the
// validation prefix, after the optional leading protocol verifier frames, must match
// one of the four recognized shapes and stay within the verification gas caps.
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

	// The protocol verifier frames lead the frame list in a fixed order, an expiry
	// frame then a recent root frame, each at most once, and are skipped when the
	// prefix is matched against the recognized shapes.
	if err := tx.validateProtocolVerifierPlacement(); err != nil {
		return err
	}

	shape := species
	if len(shape) > 0 && shape[0] == SpeciesExpiryVerify {
		shape = shape[1:]
	}

	if len(shape) > 0 && shape[0] == SpeciesRecentRootVerify {
		shape = shape[1:]
	}

	// A deploy frame leads what remains.
	deployIndex := prefixLen - len(shape)

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

// validateProtocolVerifierPlacement checks where the protocol-defined verifier frames
// may sit for public mempool eligibility: an expiry frame only at index 0, a recent
// root frame only directly after it or at index 0, and no more than one of each.
func (tx *FrameTx) validateProtocolVerifierPlacement() error {
	recentRootIndex := -1

	for i, frame := range tx.Frames {
		if frame.IsExpiryVerifier() && i != 0 {
			return fmt.Errorf("%w: expiry verifier frame must be the first frame", ErrMempoolPolicy)
		}

		if !frame.IsRecentRootVerifier() {
			continue
		}

		if recentRootIndex >= 0 {
			return fmt.Errorf("%w: more than one recent root verifier frame", ErrMempoolPolicy)
		}

		recentRootIndex = i

		allowed := 0
		if tx.Frames[0].IsExpiryVerifier() {
			allowed = 1
		}

		if i != allowed {
			return fmt.Errorf("%w: recent root verifier frame must directly follow the expiry frame or lead the frame list", ErrMempoolPolicy)
		}
	}

	return nil
}

// validatePrefixGas checks the verification gas caps over the validation prefix: the
// signature verification cost plus every prefix frame's execution budget against
// MaxVerifyGas, and the summed state budgets against MaxVerifyStateGas.
//
// The protocol verifier frames count like any other. EIP-8272 says so explicitly for the
// recent root verifier frame, and the expiry frame is simulated prefix work as well.
func (tx *FrameTx) validatePrefixGas(prefixLen int) error {
	sigGas, err := tx.SignatureVerificationGas()
	if err != nil {
		return err
	}

	executionGas := sigGas
	stateGas := uint64(0)

	for i := 0; i < prefixLen; i++ {
		executionGas = saturatingAdd(executionGas, tx.Frames[i].Limits.Execution)
		stateGas = saturatingAdd(stateGas, tx.Frames[i].Limits.State)
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
// it authorizes, which is what a client does before executing any frame: SECP256K1
// entries by recovery against their resolved signer, P256 entries by P256VERIFY over
// the public key they carry. ARBITRARY entries are not protocol-validated.
//
// It is kept separate from ValidatePayload because it does curve arithmetic, while
// ValidatePayload only checks what is decidable from the encoded fields.
func (tx *FrameTx) VerifySignatures() error {
	sigHash := tx.SigHash()

	for i, sig := range tx.Signatures {
		if sig.Scheme == SigSchemeArbitrary {
			continue
		}

		if err := tx.validateSignatureBytes(i, sig); err != nil {
			return err
		}

		digest := sigHash
		if len(sig.Msg) == 32 {
			digest = common.BytesToHash(sig.Msg)
		}

		switch sig.Scheme {
		case SigSchemeSecp256k1:
			if err := verifySecp256k1Entry(i, sig, digest, tx.Sender); err != nil {
				return err
			}

		case SigSchemeP256:
			if err := verifyP256Entry(i, sig, digest); err != nil {
				return err
			}
		}
	}

	return nil
}

// verifySecp256k1Entry recovers the entry's signer and compares it with the resolved one.
func verifySecp256k1Entry(index int, sig *FrameSignature, digest common.Hash, sender common.Address) error {
	// go-ethereum's recovery wants r || s || v, the entry carries v || r || s.
	raw := make([]byte, 65)
	copy(raw[:64], sig.Signature[1:])
	raw[64] = sig.Signature[0]

	pub, err := crypto.Ecrecover(digest[:], raw)
	if err != nil {
		return fmt.Errorf("%w: signature %d could not be recovered: %v", ErrInvalidFrameTx, index, err)
	}

	var recovered common.Address

	copy(recovered[:], crypto.Keccak256(pub[1:])[12:])

	expected, _ := sig.ResolvedSigner(sender)
	if recovered != expected {
		return fmt.Errorf("%w: signature %d recovered %s, expected %s",
			ErrInvalidFrameTx, index, recovered, expected)
	}

	return nil
}

// verifyP256Entry checks the entry's signature over the public key it carries, which is
// what P256VERIFY does. The signer address was already matched against that key by
// validateSignatureBytes, so a passing signature binds the resolved signer.
func verifyP256Entry(index int, sig *FrameSignature, digest common.Hash) error {
	// The uncompressed point encoding is 0x04 || qx || qy; parsing it performs the
	// on-curve check.
	point := append([]byte{0x04}, sig.Signature[64:128]...)

	pub, err := ecdsa.ParseUncompressedPublicKey(elliptic.P256(), point)
	if err != nil {
		return fmt.Errorf("%w: signature %d carries an invalid P-256 public key: %v", ErrInvalidFrameTx, index, err)
	}

	r := new(big.Int).SetBytes(sig.Signature[0:32])
	s := new(big.Int).SetBytes(sig.Signature[32:64])

	if !ecdsa.Verify(pub, digest[:], r, s) {
		return fmt.Errorf("%w: signature %d does not verify under P-256", ErrInvalidFrameTx, index)
	}

	return nil
}
