package txtypes

import (
	"encoding/binary"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

// Builders for the frame shapes EIP-8141 defines. They encode the mode and flag
// combinations the public mempool recognizes, so callers describe intent rather than
// bit patterns.

// SelfVerifyFrame builds a VERIFY frame that approves both execution and payment for
// the sender. It is the single-frame validation prefix used by self-relayed
// transactions.
func SelfVerifyFrame(limits FrameLimits) *Frame {
	return &Frame{
		Mode:   FrameModeVerify,
		Flags:  ApproveExecutionAndPayment,
		Target: nil,
		Limits: limits,
		Value:  new(uint256.Int),
	}
}

// OnlyVerifyFrame builds a VERIFY frame that approves execution but not payment,
// leaving payment to a following pay frame.
func OnlyVerifyFrame(limits FrameLimits) *Frame {
	return &Frame{
		Mode:   FrameModeVerify,
		Flags:  ApproveExecution,
		Target: nil,
		Limits: limits,
		Value:  new(uint256.Int),
	}
}

// PayFrame builds a VERIFY frame in which a paymaster approves payment.
func PayFrame(paymaster common.Address, data []byte, limits FrameLimits) *Frame {
	return &Frame{
		Mode:   FrameModeVerify,
		Flags:  ApprovePayment,
		Target: &paymaster,
		Limits: limits,
		Value:  new(uint256.Int),
		Data:   data,
	}
}

// ExpiryFrame builds a VERIFY frame calling the expiry verifier predeploy, which
// reverts once the deadline has passed. It may only be the first frame, and nodes drop
// the transaction from the mempool once the deadline is in the past.
func ExpiryFrame(deadline uint64, executionGas uint64) *Frame {
	data := make([]byte, ExpiryDataLength)
	binary.BigEndian.PutUint64(data, deadline)

	target := ExpiryVerifier

	return &Frame{
		Mode:   FrameModeVerify,
		Flags:  ApproveNone,
		Target: &target,
		Limits: FrameLimits{Execution: executionGas},
		Value:  new(uint256.Int),
		Data:   data,
	}
}

// UserOpFrame builds a SENDER frame: the user operation, executed with tx.sender as
// the caller. A nil target resolves to the sender itself.
func UserOpFrame(target *common.Address, value *uint256.Int, data []byte, limits FrameLimits) *Frame {
	if value == nil {
		value = new(uint256.Int)
	}

	return &Frame{
		Mode:   FrameModeSender,
		Flags:  ApproveNone,
		Target: copyAddressPtr(target),
		Limits: limits,
		Value:  value,
		Data:   data,
	}
}

// DeployFrame builds a DEFAULT frame that deploys the sender's account code, typically
// by calling a deterministic factory. It must be the first frame.
func DeployFrame(factory common.Address, data []byte, limits FrameLimits) *Frame {
	return &Frame{
		Mode:   FrameModeDefault,
		Flags:  ApproveNone,
		Target: &factory,
		Limits: limits,
		Value:  new(uint256.Int),
		Data:   data,
	}
}

// PostTxFrame builds an EIP-7906 POST_TX assertion frame.
//
// Not to be confused with PostOpFrame: that is a DEFAULT frame a paymaster uses to
// settle up, this is a STATICCALL-semantics assertion that may not APPROVE and whose
// failure reverts the whole execution body. POST_TX frames must be a trailing suffix
// of the frame list.
func PostTxFrame(target common.Address, data []byte, limits FrameLimits) *Frame {
	return &Frame{
		Mode:   FrameModePostTx,
		Flags:  ApproveNone,
		Target: &target,
		Limits: limits,
		Value:  new(uint256.Int),
		Data:   data,
	}
}

// PostOpFrame builds a DEFAULT frame that runs after the user operations, used by
// paymasters for settlement.
func PostOpFrame(target common.Address, data []byte, limits FrameLimits) *Frame {
	return &Frame{
		Mode:   FrameModeDefault,
		Flags:  ApproveNone,
		Target: &target,
		Limits: limits,
		Value:  new(uint256.Int),
		Data:   data,
	}
}

// WithAtomicBatch marks the frame as batched with the frames that follow it. All
// frames in a batch succeed or all of them revert.
func (f *Frame) WithAtomicBatch() *Frame {
	f.Flags |= AtomicBatchFlag

	return f
}

// SenderSignature builds an unsigned SECP256K1 entry authorizing the canonical
// signature hash on the sender's behalf. The default account code requires this entry
// at index 0 when the verify frame approves execution, and at index 1 otherwise.
func SenderSignature() *FrameSignature {
	return &FrameSignature{Scheme: SigSchemeSecp256k1}
}

// SignerSignature builds an unsigned SECP256K1 entry for an explicit signer.
func SignerSignature(signer common.Address) *FrameSignature {
	return &FrameSignature{Scheme: SigSchemeSecp256k1, Signer: signer.Bytes()}
}

// ArbitrarySignature builds an ARBITRARY witness entry. The protocol does not validate
// it; contracts read it with SIGDATACOPY.
func ArbitrarySignature(witness []byte) *FrameSignature {
	return &FrameSignature{Scheme: SigSchemeArbitrary, Signature: witness}
}

// DefaultCodeVerifyLimits returns a validation frame budget sized for an account with
// no deployed code.
//
// The default code draws no execution gas of its own: the frame's only execution
// charge is the resolved target's EIP-2929 account access, and the sender is warm from
// the start of the transaction. State gas is only needed when the sender account does
// not exist yet, in which case APPROVE charges for creating it.
func DefaultCodeVerifyLimits(senderExists bool) FrameLimits {
	limits := FrameLimits{Execution: 5_000}

	if !senderExists {
		limits.State = StateBytesPerNewAccount * CostPerStateByte
	}

	return limits
}

// NewFrameTx assembles a frame transaction from its parts, filling in the zero values
// the encoder requires.
//
// It builds the envelope shape current devnets run, with both extensions active and
// the nonce in EIP-8250's [0] key set, which is the sender's ordinary account nonce.
// Use NewFrameTxWithExtensions for a chain running a different combination.
func NewFrameTx(chainID *uint256.Int, sender common.Address, nonce uint64, fees FrameFees, frames []*Frame, signatures []*FrameSignature) *FrameTx {
	return NewFrameTxWithExtensions(FrameExtAll, chainID, sender, nonce, fees, frames, signatures)
}

// NewFrameTxWithExtensions assembles a frame transaction in a chosen envelope shape.
func NewFrameTxWithExtensions(extensions FrameExtensions, chainID *uint256.Int, sender common.Address, nonce uint64, fees FrameFees, frames []*Frame, signatures []*FrameSignature) *FrameTx {
	if fees.GasTipCap == nil {
		fees.GasTipCap = new(uint256.Int)
	}

	if fees.GasFeeCap == nil {
		fees.GasFeeCap = new(uint256.Int)
	}

	if fees.BlobFeeCap == nil {
		fees.BlobFeeCap = new(uint256.Int)
	}

	if chainID == nil {
		chainID = new(uint256.Int)
	}

	for _, frame := range frames {
		if frame.Value == nil {
			frame.Value = new(uint256.Int)
		}
	}

	tx := &FrameTx{
		ChainID:    chainID,
		NonceSeq:   nonce,
		Sender:     sender,
		Frames:     frames,
		Signatures: signatures,
		Fees:       fees,
		Extensions: extensions,
	}

	if extensions.Has(FrameExtKeyedNonces) {
		tx.NonceKeys = []*uint256.Int{new(uint256.Int)}
	}

	return tx
}
