package txtypes

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/holiman/uint256"
)

// EIP-8141 frame transaction constants.
const (
	FrameTxType = 0x06

	FrameTxIntrinsicCost = 12000
	FrameTxPerFrameCost  = 475
	MaxFrames            = 64
	ExpiryDataLength     = 8

	// Public mempool limits.
	MaxVerifyGas      = 100_000
	MaxVerifyStateGas = 500_000

	// EIP-8250 keyed nonces.
	MaxNonceKeys = 16

	// EIP-8272 recent root references.
	MaxRecentRootReferences = 16
)

// Constants defined by the EIPs frame transactions build on.
const (
	TxMaxGasLimit            = 16_777_216 // EIP-7825
	TxValueCost              = 6_000      // EIP-2780
	CostPerStateByte         = 1_530      // EIP-8037
	StateBytesPerNewAccount  = 120        // EIP-8037
	StandardTokenCost        = 4          // EIP-7976
	TotalCostFloorPerToken   = 16         // EIP-7976
	GasPerBlob               = 131_072    // EIP-4844
	VersionedHashVersionKZG  = 0x01       // EIP-4844
	FrameBlobSidecarVersion1 = 1          // EIP-7594 wrapper version
)

// Protocol-defined addresses.
var (
	// EntryPoint is the caller of DEFAULT and VERIFY mode frames.
	EntryPoint = common.HexToAddress("0x00000000000000000000000000000000000000aa")

	// ExpiryVerifier is the built-in deadline verifier predeploy.
	ExpiryVerifier = common.HexToAddress("0x0000000000000000000000000000000000008141")

	// NonceManager stores the non-zero keyed nonce sequences of EIP-8250.
	NonceManager = common.HexToAddress("0x0000000000000000000000000000000000008250")

	// RecentRootAddress stores the recent roots of EIP-8272.
	RecentRootAddress = common.HexToAddress("0x0000000000000000000000000000000000008272")

	// LegacyNonceKey selects the sender's ordinary account nonce sequence.
	LegacyNonceKey = uint256.NewInt(0)
)

// FrameMode is a frame's execution context.
type FrameMode uint8

// Frame modes.
const (
	FrameModeDefault FrameMode = 0 // caller is ENTRY_POINT
	FrameModeVerify  FrameMode = 1 // STATICCALL semantics, must not revert
	FrameModeSender  FrameMode = 2 // caller is tx.sender, may carry value
)

// Frame flags. Bits 0-1 are the approval scope, bit 2 marks an atomic batch.
const (
	ApproveNone                = 0x0
	ApprovePayment             = 0x1
	ApproveExecution           = 0x2
	ApproveExecutionAndPayment = 0x3
	ApproveScopeMask           = 0x3
	AtomicBatchFlag            = 0x4
	FrameFlagsMask             = 0x7
)

// FrameSigScheme identifies how a signature entry's raw bytes are interpreted.
type FrameSigScheme uint8

// Signature schemes and their verification costs.
const (
	SigSchemeArbitrary FrameSigScheme = 0x0
	SigSchemeSecp256k1 FrameSigScheme = 0x1
	SigSchemeP256      FrameSigScheme = 0x2

	sigGasArbitrary = 100
	sigGasSecp256k1 = 2800
	sigGasP256      = 6700
)

var (
	// ErrInvalidFrameTx is the base error for static frame transaction validation.
	ErrInvalidFrameTx = errors.New("invalid frame transaction")

	// ErrMempoolPolicy is returned when a transaction is well-formed but would not
	// propagate through the public mempool.
	ErrMempoolPolicy = errors.New("frame transaction violates public mempool policy")
)

func init() {
	RegisterTxType(FrameTxType, func() TxData { return &FrameTx{} })
}

// FrameLimits is a frame's two-dimensional gas budget.
type FrameLimits struct {
	Execution uint64
	State     uint64
}

// Frame is a single call in a frame transaction.
type Frame struct {
	Mode   FrameMode
	Flags  uint8
	Target *common.Address `rlp:"nil"` // nil resolves to tx.sender
	Limits FrameLimits
	Value  *uint256.Int
	Data   []byte
}

// FrameSignature is one entry of a frame transaction's signature list.
type FrameSignature struct {
	Scheme    FrameSigScheme
	Signer    []byte // 20-byte address, or empty for tx.sender
	Msg       []byte // empty for the canonical signature hash, or a 32-byte digest
	Signature []byte
}

// FrameFees holds a frame transaction's EIP-1559 and blob fee parameters.
type FrameFees struct {
	GasTipCap  *uint256.Int
	GasFeeCap  *uint256.Int
	BlobFeeCap *uint256.Int
}

// RecentRootReference is an EIP-8272 declared recent root, identified by its source
// and slot.
type RecentRootReference struct {
	SourceID common.Hash
	Slot     uint64
	Root     common.Hash
}

// Copy returns a copy of the reference.
func (r *RecentRootReference) Copy() *RecentRootReference {
	cpy := *r

	return &cpy
}

// FrameTx is an EIP-8141 frame transaction, with the envelope extensions of EIP-8250
// (keyed nonces) and EIP-8272 (recent root references).
//
// Its sender is an explicit field rather than something recovered from a signature,
// and its signature list is validated by the protocol before any frame executes.
//
// EIP-8250 replaces EIP-8141's single nonce with NonceKeys and NonceSeq. NonceKeys of
// exactly [0] selects the sender's ordinary account nonce, which is what a transaction
// that does not use independent nonce domains carries.
type FrameTx struct {
	ChainID     *uint256.Int
	NonceKeys   []*uint256.Int
	NonceSeq    uint64
	Sender      common.Address
	Frames      []*Frame
	Signatures  []*FrameSignature
	Fees        FrameFees
	BlobHashes  []common.Hash
	RecentRoots []*RecentRootReference

	// Sidecar is wire-only data, excluded from the canonical encoding and the hash.
	Sidecar *BlobSidecar `rlp:"-"`
}

var (
	_ TxData           = (*FrameTx)(nil)
	_ ExplicitSenderTx = (*FrameTx)(nil)
	_ NetworkEncodedTx = (*FrameTx)(nil)
	_ BlobTxData       = (*FrameTx)(nil)
	_ StateGasTxData   = (*FrameTx)(nil)
)

// frameTxWithSidecar is the EIP-7594 wrapper used for blob-carrying frame
// transactions on the wire.
type frameTxWithSidecar struct {
	Tx          *FrameTx
	Version     byte
	Blobs       []kzg4844.Blob
	Commitments []kzg4844.Commitment
	Proofs      []kzg4844.Proof
}

// ResolvedTarget returns the frame's target, resolving a nil target to sender.
func (f *Frame) ResolvedTarget(sender common.Address) common.Address {
	if f.Target == nil {
		return sender
	}

	return *f.Target
}

// ApprovalScope returns the approval scope bits of the frame's flags.
func (f *Frame) ApprovalScope() uint8 {
	return f.Flags & ApproveScopeMask
}

// IsAtomicBatch reports whether the frame is batched with the frames that follow.
func (f *Frame) IsAtomicBatch() bool {
	return f.Flags&AtomicBatchFlag != 0
}

// IsExpiryVerifier reports whether the frame is an expiry verifier frame. Any VERIFY
// frame targeting EXPIRY_VERIFIER is one, whatever its flags; the flags are then
// constrained by the frame's own validity rules.
func (f *Frame) IsExpiryVerifier() bool {
	return f.Mode == FrameModeVerify && f.Target != nil && *f.Target == ExpiryVerifier
}

// ExpiryDeadline returns the deadline encoded in an expiry verifier frame's data.
func (f *Frame) ExpiryDeadline() (uint64, bool) {
	if !f.IsExpiryVerifier() || len(f.Data) != ExpiryDataLength {
		return 0, false
	}

	return binary.BigEndian.Uint64(f.Data), true
}

// Copy returns a deep copy of the frame.
func (f *Frame) Copy() *Frame {
	cpy := &Frame{
		Mode:   f.Mode,
		Flags:  f.Flags,
		Target: copyAddressPtr(f.Target),
		Limits: f.Limits,
		Value:  new(uint256.Int),
		Data:   common.CopyBytes(f.Data),
	}

	setU256(cpy.Value, f.Value)

	return cpy
}

// Copy returns a deep copy of the signature entry.
func (s *FrameSignature) Copy() *FrameSignature {
	return &FrameSignature{
		Scheme:    s.Scheme,
		Signer:    common.CopyBytes(s.Signer),
		Msg:       common.CopyBytes(s.Msg),
		Signature: common.CopyBytes(s.Signature),
	}
}

// ResolvedSigner returns the entry's signer, resolving an empty signer to sender.
// It reports false for ARBITRARY entries, which have no protocol-assigned signer.
func (s *FrameSignature) ResolvedSigner(sender common.Address) (common.Address, bool) {
	if s.Scheme == SigSchemeArbitrary {
		return common.Address{}, false
	}

	if len(s.Signer) == 0 {
		return sender, true
	}

	return common.BytesToAddress(s.Signer), true
}

// VerificationGas returns the entry's protocol signature verification cost.
func (s *FrameSignature) VerificationGas() (uint64, error) {
	switch s.Scheme {
	case SigSchemeArbitrary:
		return sigGasArbitrary, nil
	case SigSchemeSecp256k1:
		return sigGasSecp256k1, nil
	case SigSchemeP256:
		return sigGasP256, nil
	default:
		return 0, fmt.Errorf("%w: unknown signature scheme 0x%x", ErrInvalidFrameTx, s.Scheme)
	}
}

// TxType returns the EIP-2718 type byte.
func (tx *FrameTx) TxType() byte { return FrameTxType }

// CopyTx returns a deep copy with all fields initialized. The sidecar is shared, as
// for blob transactions.
func (tx *FrameTx) CopyTx() TxData {
	cpy := &FrameTx{
		NonceSeq:    tx.NonceSeq,
		Sender:      tx.Sender,
		NonceKeys:   make([]*uint256.Int, len(tx.NonceKeys)),
		Frames:      make([]*Frame, len(tx.Frames)),
		Signatures:  make([]*FrameSignature, len(tx.Signatures)),
		BlobHashes:  make([]common.Hash, len(tx.BlobHashes)),
		RecentRoots: make([]*RecentRootReference, len(tx.RecentRoots)),
		Sidecar:     tx.Sidecar,
		ChainID:     new(uint256.Int),
		Fees: FrameFees{
			GasTipCap:  new(uint256.Int),
			GasFeeCap:  new(uint256.Int),
			BlobFeeCap: new(uint256.Int),
		},
	}

	for i, key := range tx.NonceKeys {
		cpy.NonceKeys[i] = new(uint256.Int)
		setU256(cpy.NonceKeys[i], key)
	}

	for i, frame := range tx.Frames {
		cpy.Frames[i] = frame.Copy()
	}

	for i, root := range tx.RecentRoots {
		cpy.RecentRoots[i] = root.Copy()
	}

	for i, sig := range tx.Signatures {
		cpy.Signatures[i] = sig.Copy()
	}

	copy(cpy.BlobHashes, tx.BlobHashes)
	setU256(cpy.ChainID, tx.ChainID)
	setU256(cpy.Fees.GasTipCap, tx.Fees.GasTipCap)
	setU256(cpy.Fees.GasFeeCap, tx.Fees.GasFeeCap)
	setU256(cpy.Fees.BlobFeeCap, tx.Fees.BlobFeeCap)

	return cpy
}

// GetChainID returns the destination chain id.
func (tx *FrameTx) GetChainID() *big.Int { return u256ToBig(tx.ChainID) }

// GetNonce returns the nonce sequence number. For the ordinary [0] key set this is
// the sender's account nonce, which is what the engine tracks.
func (tx *FrameTx) GetNonce() uint64 { return tx.NonceSeq }

// UsesLegacyNonce reports whether the transaction selects only the sender's ordinary
// account nonce sequence.
func (tx *FrameTx) UsesLegacyNonce() bool {
	return len(tx.NonceKeys) == 1 && tx.NonceKeys[0] != nil && tx.NonceKeys[0].IsZero()
}

// GetGas returns the maximum gas the transaction can consume, which is the budget the
// payer is charged against.
func (tx *FrameTx) GetGas() uint64 { return tx.MaxGas() }

// GetGasPrice returns the fee cap.
func (tx *FrameTx) GetGasPrice() *big.Int { return u256ToBig(tx.Fees.GasFeeCap) }

// GetGasFeeCap returns the maximum fee per gas.
func (tx *FrameTx) GetGasFeeCap() *big.Int { return u256ToBig(tx.Fees.GasFeeCap) }

// GetGasTipCap returns the maximum priority fee per gas.
func (tx *FrameTx) GetGasTipCap() *big.Int { return u256ToBig(tx.Fees.GasTipCap) }

// GetTo returns the resolved target of the first SENDER frame, or nil when the
// transaction has none. A frame transaction has no single recipient; this is the
// closest equivalent for engine-level accounting and logging.
func (tx *FrameTx) GetTo() *common.Address {
	for _, frame := range tx.Frames {
		if frame.Mode == FrameModeSender {
			target := frame.ResolvedTarget(tx.Sender)

			return &target
		}
	}

	return nil
}

// GetValue returns the sum of all frame values.
func (tx *FrameTx) GetValue() *big.Int {
	total := new(big.Int)

	for _, frame := range tx.Frames {
		if frame.Value != nil {
			total.Add(total, frame.Value.ToBig())
		}
	}

	return total
}

// GetData returns the calldata of the first SENDER frame, or nil.
func (tx *FrameTx) GetData() []byte {
	for _, frame := range tx.Frames {
		if frame.Mode == FrameModeSender {
			return frame.Data
		}
	}

	return nil
}

// GetSender returns the transaction's declared sender.
func (tx *FrameTx) GetSender() common.Address { return tx.Sender }

// GetBlobHashes returns the blob versioned hashes.
func (tx *FrameTx) GetBlobHashes() []common.Hash { return tx.BlobHashes }

// GetBlobGasFeeCap returns the maximum fee per blob gas.
func (tx *FrameTx) GetBlobGasFeeCap() *big.Int { return u256ToBig(tx.Fees.BlobFeeCap) }

// GetBlobSidecar returns the attached sidecar, or nil.
func (tx *FrameTx) GetBlobSidecar() *BlobSidecar { return tx.Sidecar }

// SetBlobSidecar attaches a sidecar.
func (tx *FrameTx) SetBlobSidecar(sidecar *BlobSidecar) { tx.Sidecar = sidecar }

// GetStateGas returns the sum of the declared frame state gas budgets, which is the
// transaction's reservation against the EIP-8037 block state gas dimension.
func (tx *FrameTx) GetStateGas() uint64 {
	total := uint64(0)

	for _, frame := range tx.Frames {
		total += frame.Limits.State
	}

	return total
}

// EncodePayload writes the canonical payload, without the blob sidecar.
func (tx *FrameTx) EncodePayload(w *bytes.Buffer) error {
	return rlp.Encode(w, tx)
}

// EncodeNetworkPayload writes the payload including the blob sidecar, using the
// EIP-7594 wrapper. Frame transactions have no unversioned wrapper form.
func (tx *FrameTx) EncodeNetworkPayload(w *bytes.Buffer) error {
	if tx.Sidecar == nil {
		return rlp.Encode(w, tx)
	}

	return rlp.Encode(w, &frameTxWithSidecar{
		Tx:          tx,
		Version:     FrameBlobSidecarVersion1,
		Blobs:       tx.Sidecar.Blobs,
		Commitments: tx.Sidecar.Commitments,
		Proofs:      tx.Sidecar.Proofs,
	})
}

// DecodePayload parses either the canonical payload or the EIP-7594 wrapper. The two
// are told apart by whether the first element is itself a list, as for blob
// transactions.
func (tx *FrameTx) DecodePayload(b []byte) error {
	firstElem, _, err := rlp.SplitList(b)
	if err != nil {
		return err
	}

	firstElemKind, _, _, err := rlp.Split(firstElem)
	if err != nil {
		return err
	}

	if firstElemKind != rlp.List {
		return rlp.DecodeBytes(b, tx)
	}

	var payload frameTxWithSidecar
	if err := rlp.DecodeBytes(b, &payload); err != nil {
		return err
	}

	if payload.Tx == nil {
		return errors.New("frame transaction wrapper without transaction")
	}

	if payload.Version != FrameBlobSidecarVersion1 {
		return fmt.Errorf("unsupported blob sidecar version %d", payload.Version)
	}

	*tx = *payload.Tx
	tx.Sidecar = &BlobSidecar{
		Version:     payload.Version,
		Blobs:       payload.Blobs,
		Commitments: payload.Commitments,
		Proofs:      payload.Proofs,
	}

	return nil
}

// SigHash returns the canonical signature hash.
//
// Signature entries that sign the canonical digest have their own raw bytes elided
// from the preimage; entries carrying an explicit digest keep theirs.
func (tx *FrameTx) SigHash() common.Hash {
	elided := &FrameTx{
		ChainID:     tx.ChainID,
		NonceKeys:   tx.NonceKeys,
		NonceSeq:    tx.NonceSeq,
		Sender:      tx.Sender,
		Frames:      tx.Frames,
		Signatures:  make([]*FrameSignature, len(tx.Signatures)),
		Fees:        tx.Fees,
		BlobHashes:  tx.BlobHashes,
		RecentRoots: tx.RecentRoots,
	}

	for i, sig := range tx.Signatures {
		if len(sig.Msg) == 0 {
			elided.Signatures[i] = &FrameSignature{
				Scheme: sig.Scheme,
				Signer: sig.Signer,
				Msg:    sig.Msg,
			}
		} else {
			elided.Signatures[i] = sig
		}
	}

	return prefixedRlpHash(FrameTxType, elided)
}

// SignPayload signs every unsigned SECP256K1 entry that the key is the signer for.
//
// An entry qualifies when it carries no signature yet and its signer is either empty
// (resolving to the sender) or the key's own address. Entries for other parties, such
// as a paymaster's, are signed with SignEntry.
func (tx *FrameTx) SignPayload(chainID *big.Int, key *ecdsa.PrivateKey) error {
	if chainID != nil && chainID.Sign() != 0 {
		converted, overflow := uint256.FromBig(chainID)
		if overflow {
			return errors.New("chain id overflows 256 bits")
		}

		tx.ChainID = converted
	}

	address := crypto.PubkeyToAddress(key.PublicKey)
	signed := 0

	for i, sig := range tx.Signatures {
		if sig.Scheme != SigSchemeSecp256k1 || len(sig.Signature) != 0 {
			continue
		}

		if len(sig.Signer) != 0 && common.BytesToAddress(sig.Signer) != address {
			continue
		}

		if err := tx.SignEntry(i, key); err != nil {
			return err
		}

		signed++
	}

	if signed == 0 {
		return errors.New("frame transaction has no secp256k1 entry for this key")
	}

	return nil
}

// SignEntry signs a single signature entry with key.
//
// The digest is the entry's explicit msg when it carries one, and the canonical
// signature hash otherwise. Because the canonical hash covers every entry's metadata,
// entries must be fully populated before any of them is signed.
func (tx *FrameTx) SignEntry(index int, key *ecdsa.PrivateKey) error {
	if index < 0 || index >= len(tx.Signatures) {
		return fmt.Errorf("signature index %d out of range", index)
	}

	sig := tx.Signatures[index]

	digest, err := tx.entryDigest(sig)
	if err != nil {
		return err
	}

	switch sig.Scheme {
	case SigSchemeSecp256k1:
		raw, err := crypto.Sign(digest[:], key)
		if err != nil {
			return fmt.Errorf("failed signing frame signature entry: %w", err)
		}

		// The entry encoding is v || r || s, not go-ethereum's r || s || v.
		sig.Signature = append([]byte{raw[64]}, raw[:64]...)

	case SigSchemeP256:
		return errors.New("use SignEntryP256 for P256 entries")

	default:
		return fmt.Errorf("%w: scheme 0x%x is not protocol-signed", ErrInvalidFrameTx, sig.Scheme)
	}

	return nil
}

// SignEntryP256 signs a signature entry with a NIST P-256 key, filling in both the
// signature and the derived signer address.
func (tx *FrameTx) SignEntryP256(index int, key *ecdsa.PrivateKey) error {
	if index < 0 || index >= len(tx.Signatures) {
		return fmt.Errorf("signature index %d out of range", index)
	}

	sig := tx.Signatures[index]
	if sig.Scheme != SigSchemeP256 {
		return fmt.Errorf("signature entry %d is not a P256 entry", index)
	}

	if key.Curve != elliptic.P256() {
		return errors.New("key is not on the P-256 curve")
	}

	digest, err := tx.entryDigest(sig)
	if err != nil {
		return err
	}

	r, s, err := ecdsa.Sign(rand.Reader, key, digest[:])
	if err != nil {
		return fmt.Errorf("failed signing P256 entry: %w", err)
	}

	// Low-s normalization, as required for secp256r1 entries.
	order := elliptic.P256().Params().N
	if s.Cmp(new(big.Int).Rsh(order, 1)) > 0 {
		s = new(big.Int).Sub(order, s)
	}

	raw := make([]byte, 128)
	r.FillBytes(raw[0:32])
	s.FillBytes(raw[32:64])
	key.X.FillBytes(raw[64:96])
	key.Y.FillBytes(raw[96:128])

	sig.Signature = raw
	sig.Signer = P256Signer(key.X, key.Y).Bytes()

	return nil
}

// P256Signer returns the address a P256 public key resolves to.
func P256Signer(qx, qy *big.Int) common.Address {
	buf := make([]byte, 64)
	qx.FillBytes(buf[:32])
	qy.FillBytes(buf[32:])

	return common.BytesToAddress(crypto.Keccak256(buf)[12:])
}

// entryDigest returns the digest a signature entry authorizes.
func (tx *FrameTx) entryDigest(sig *FrameSignature) (common.Hash, error) {
	switch len(sig.Msg) {
	case 0:
		return tx.SigHash(), nil
	case 32:
		return common.BytesToHash(sig.Msg), nil
	default:
		return common.Hash{}, fmt.Errorf("%w: signature msg must be empty or 32 bytes", ErrInvalidFrameTx)
	}
}

// SignatureVerificationGas returns the total protocol signature verification cost.
func (tx *FrameTx) SignatureVerificationGas() (uint64, error) {
	total := uint64(0)

	for _, sig := range tx.Signatures {
		gas, err := sig.VerificationGas()
		if err != nil {
			return 0, err
		}

		total += gas
	}

	return total, nil
}

// IntrinsicGas returns the transaction's intrinsic cost, which is derivable from the
// transaction fields alone and charged entirely in the execution dimension.
func (tx *FrameTx) IntrinsicGas() uint64 {
	sigGas, err := tx.SignatureVerificationGas()
	if err != nil {
		sigGas = 0
	}

	gas := uint64(FrameTxIntrinsicCost) + uint64(len(tx.Frames))*FrameTxPerFrameCost + sigGas

	for _, frame := range tx.Frames {
		gas += calldataCost(frame.Data)
		gas += tx.valueCost(frame)
	}

	for _, sig := range tx.Signatures {
		gas += calldataCost(sig.Signer) + calldataCost(sig.Msg) + calldataCost(sig.Signature)
	}

	return gas + calldataCost(tx.nonceCalldata())
}

// nonceCalldata returns the encoding EIP-8250 prices as transaction data:
// rlp(nonce_keys) || rlp(nonce_seq).
func (tx *FrameTx) nonceCalldata() []byte {
	var buf bytes.Buffer

	if err := rlp.Encode(&buf, tx.NonceKeys); err != nil {
		return nil
	}

	if err := rlp.Encode(&buf, tx.NonceSeq); err != nil {
		return nil
	}

	return buf.Bytes()
}

// CalldataFloorGas returns the EIP-7976 floor cost of the transaction's data.
func (tx *FrameTx) CalldataFloorGas() uint64 {
	sigGas, err := tx.SignatureVerificationGas()
	if err != nil {
		sigGas = 0
	}

	tokens := uint64(0)

	for _, frame := range tx.Frames {
		tokens += uint64(len(frame.Data)) * StandardTokenCost
	}

	for _, sig := range tx.Signatures {
		tokens += uint64(len(sig.Signer)+len(sig.Msg)+len(sig.Signature)) * StandardTokenCost
	}

	// EIP-8250 adds the nonce encoding to the token count. It contributes its
	// weighted count, as the EIP prices it alongside the other transaction data.
	tokens += weightedTokens(tx.nonceCalldata())

	gas := uint64(FrameTxIntrinsicCost) + uint64(len(tx.Frames))*FrameTxPerFrameCost + sigGas

	for _, frame := range tx.Frames {
		gas += tx.valueCost(frame)
	}

	return gas + TotalCostFloorPerToken*tokens
}

// ExecutionGas returns the intrinsic cost plus the declared frame execution budgets,
// which is the transaction's reservation against the block execution gas dimension.
func (tx *FrameTx) ExecutionGas() uint64 {
	gas := tx.IntrinsicGas()

	for _, frame := range tx.Frames {
		gas += frame.Limits.Execution
	}

	if floor := tx.CalldataFloorGas(); floor > gas {
		return floor
	}

	return gas
}

// MaxGas returns the maximum gas the transaction can consume across both dimensions.
func (tx *FrameTx) MaxGas() uint64 {
	standard := tx.IntrinsicGas()
	stateGas := tx.GetStateGas()

	for _, frame := range tx.Frames {
		standard += frame.Limits.Execution
	}

	standard += stateGas

	if floor := tx.CalldataFloorGas() + stateGas; floor > standard {
		return floor
	}

	return standard
}

// MaxCost returns the amount the payer is charged up front when payment is approved.
func (tx *FrameTx) MaxCost(blobBaseFee *big.Int) *big.Int {
	cost := new(big.Int).Mul(new(big.Int).SetUint64(tx.MaxGas()), u256ToBig(tx.Fees.GasFeeCap))

	return cost.Add(cost, tx.blobGasCost(blobBaseFee))
}

// valueCost returns the static charge for a value-bearing frame.
func (tx *FrameTx) valueCost(frame *Frame) uint64 {
	if frame.Value == nil || frame.Value.IsZero() {
		return 0
	}

	if frame.Target == nil || *frame.Target == tx.Sender {
		return 0
	}

	return TxValueCost
}

// calldataCost returns the EIP-7976 weighted token cost of a byte string.
func calldataCost(data []byte) uint64 {
	return StandardTokenCost * weightedTokens(data)
}

// weightedTokens returns the EIP-7976 weighted token count of a byte string.
func weightedTokens(data []byte) uint64 {
	tokens := uint64(0)

	for _, b := range data {
		if b == 0 {
			tokens++
		} else {
			tokens += 4
		}
	}

	return tokens
}

// jsonFrame is a frame as reported by JSON-RPC. EIP-8141 does not specify a JSON
// encoding; these key names follow ethrex, the first client to ship the type.
type jsonFrame struct {
	Mode       *hexutil.Uint64 `json:"mode"`
	Flags      *hexutil.Uint64 `json:"flags"`
	To         *common.Address `json:"to"`
	Target     *common.Address `json:"target"`
	GasLimit   *hexutil.Uint64 `json:"gasLimit"`
	StateLimit *hexutil.Uint64 `json:"stateLimit"`
	Value      *hexutil.Big    `json:"value"`
	Data       *hexutil.Bytes  `json:"data"`
}

// jsonFrameSignature is a signature entry as reported by JSON-RPC.
type jsonFrameSignature struct {
	Scheme    *hexutil.Uint64 `json:"scheme"`
	Signer    *common.Address `json:"signer"`
	Msg       *hexutil.Bytes  `json:"msg"`
	Signature *hexutil.Bytes  `json:"signature"`
}

// jsonRecentRootReference is a recent root reference as reported by JSON-RPC.
type jsonRecentRootReference struct {
	SourceID common.Hash     `json:"sourceId"`
	Slot     *hexutil.Uint64 `json:"slot"`
	Root     common.Hash     `json:"root"`
}

// jsonFrameTx holds the frame-specific fields of a JSON-RPC transaction object.
type jsonFrameTx struct {
	Sender     *common.Address            `json:"sender"`
	NonceKeys  []*hexutil.Big             `json:"nonceKeys"`
	NonceSeq   *hexutil.Uint64            `json:"nonceSeq"`
	Frames     []*jsonFrame               `json:"frames"`
	Signatures []*jsonFrameSignature      `json:"signatures"`
	Roots      []*jsonRecentRootReference `json:"recentRootReferences"`
}

// DecodeJSONTx populates a frame transaction from a JSON-RPC transaction object.
//
// Blob sidecars are not part of the JSON-RPC representation.
func (tx *FrameTx) DecodeJSONTx(fields *JSONTxFields) error {
	var dec jsonFrameTx
	if err := json.Unmarshal(fields.Raw, &dec); err != nil {
		return err
	}

	if len(dec.Frames) == 0 {
		return errors.New("frame transaction object carries no frames")
	}

	tx.ChainID = jsonU256(fields.ChainID)
	tx.NonceSeq = jsonUint64(dec.NonceSeq)
	tx.Fees = FrameFees{
		GasTipCap:  jsonU256(fields.MaxPriorityFeePerGas),
		GasFeeCap:  jsonU256(fields.MaxFeePerGas),
		BlobFeeCap: jsonU256(fields.MaxFeePerBlobGas),
	}
	tx.BlobHashes = fields.BlobVersionedHashes

	// The sender is an explicit field; fall back to the generic "from" when a client
	// reports only that.
	switch {
	case dec.Sender != nil:
		tx.Sender = *dec.Sender
	case fields.From != nil:
		tx.Sender = *fields.From
	default:
		return errors.New("frame transaction object carries no sender")
	}

	tx.NonceKeys = make([]*uint256.Int, 0, len(dec.NonceKeys))
	for _, key := range dec.NonceKeys {
		tx.NonceKeys = append(tx.NonceKeys, jsonU256(key))
	}

	if len(tx.NonceKeys) == 0 {
		tx.NonceKeys = []*uint256.Int{new(uint256.Int)}
	}

	tx.Frames = make([]*Frame, 0, len(dec.Frames))

	for _, frame := range dec.Frames {
		target := frame.Target
		if target == nil {
			target = frame.To
		}

		tx.Frames = append(tx.Frames, &Frame{
			Mode:   FrameMode(jsonUint64(frame.Mode)),
			Flags:  uint8(jsonUint64(frame.Flags)),
			Target: copyAddressPtr(target),
			Limits: FrameLimits{
				Execution: jsonUint64(frame.GasLimit),
				State:     jsonUint64(frame.StateLimit),
			},
			Value: jsonU256(frame.Value),
			Data:  jsonBytes(frame.Data),
		})
	}

	tx.Signatures = make([]*FrameSignature, 0, len(dec.Signatures))

	for _, sig := range dec.Signatures {
		entry := &FrameSignature{
			Scheme:    FrameSigScheme(jsonUint64(sig.Scheme)),
			Msg:       jsonBytes(sig.Msg),
			Signature: jsonBytes(sig.Signature),
		}

		if sig.Signer != nil {
			entry.Signer = sig.Signer.Bytes()
		}

		tx.Signatures = append(tx.Signatures, entry)
	}

	tx.RecentRoots = make([]*RecentRootReference, 0, len(dec.Roots))
	for _, root := range dec.Roots {
		tx.RecentRoots = append(tx.RecentRoots, &RecentRootReference{
			SourceID: root.SourceID,
			Slot:     jsonUint64(root.Slot),
			Root:     root.Root,
		})
	}

	return nil
}
