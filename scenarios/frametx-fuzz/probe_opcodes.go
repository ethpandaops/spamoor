package frametxfuzz

// Instructions and parameter indices introduced by the frame transaction EIPs.
//
// They live here rather than in txtypes because they describe the EVM, not the
// transaction: only code written to run inside a frame encodes them. The probe contract
// emits these bytes, and the assertion helpers name these indices.
//
// The opcode bytes are those the EIPs assign: EIP-8141 takes 0xB0-0xB5 and EIP-7906
// 0xB6-0xB8. A byte above them that no enabled EIP allocates is an invalid instruction.

// Instruction opcodes.
const (
	// OpApprove is EIP-8141's APPROVE, which exits the call frame successfully and
	// updates the transaction's approval context.
	OpApprove = 0xaa

	// OpTxParam is EIP-8141's TXPARAM: transaction-scoped introspection.
	OpTxParam = 0xb0

	// OpFrameDataLoad is EIP-8141's FRAMEDATALOAD: one word of another frame's data.
	OpFrameDataLoad = 0xb1

	// OpFrameDataCopy is EIP-8141's FRAMEDATACOPY.
	OpFrameDataCopy = 0xb2

	// OpFrameParam is EIP-8141's FRAMEPARAM: frame-scoped introspection.
	OpFrameParam = 0xb3

	// OpSigParam is EIP-8141's SIGPARAM: signature-scoped introspection.
	OpSigParam = 0xb4

	// OpSigDataCopy is EIP-8141's SIGDATACOPY, valid only for ARBITRARY entries.
	OpSigDataCopy = 0xb5

	// OpTxTrace is EIP-7906's TXTRACE: the transaction's state diff, enumerated. Valid
	// only inside an ASSERT frame.
	OpTxTrace = 0xb6

	// OpTxDiff is EIP-7906's TXDIFF: the state diff looked up by address and slot.
	// Valid only inside an ASSERT frame.
	OpTxDiff = 0xb7

	// OpEventDataCopy is EIP-7906's EVENTDATACOPY: an emitted event's data, copied to
	// memory. Valid only inside an ASSERT frame.
	OpEventDataCopy = 0xb8
)

// TXPARAM parameter indices (EIP-8141).
const (
	TxParamTxType         = 0x00
	TxParamNonce          = 0x01
	TxParamSender         = 0x02
	TxParamGasTipCap      = 0x03
	TxParamGasFeeCap      = 0x04
	TxParamBlobFeeCap     = 0x05
	TxParamMaxCost        = 0x06
	TxParamBlobHashCount  = 0x07
	TxParamSigHash        = 0x08
	TxParamFrameCount     = 0x09
	TxParamFrameIndex     = 0x0a
	TxParamSignatureCount = 0x0b
	TxParamStateGasLeft   = 0x0c
)

// TXPARAM parameter indices added by EIP-8250, allocated from the family registry
// directly above EIP-8141's range.
const (
	// TxParamLegacyNonce is EIP-8250's pre-state sender account nonce.
	TxParamLegacyNonce = 0x0d

	// TxParamNonceKeyCount is EIP-8250's len(nonce_keys).
	TxParamNonceKeyCount = 0x0e

	// TxParamNonceKeysHash is EIP-8250's nonce_keys_hash.
	TxParamNonceKeysHash = 0x0f

	// TxParamNonceKey0 is EIP-8250's nonce_keys[0].
	TxParamNonceKey0 = 0x10
)

// FRAMEPARAM parameter indices (EIP-8141).
const (
	FrameParamTarget         = 0x00
	FrameParamExecutionLimit = 0x01
	FrameParamMode           = 0x02
	FrameParamFlags          = 0x03
	FrameParamDataLength     = 0x04
	FrameParamStatus         = 0x05
	FrameParamAllowedScope   = 0x06
	FrameParamAtomicBatch    = 0x07
	FrameParamValue          = 0x08
	FrameParamStateLimit     = 0x09
	FrameParamExecutionUsed  = 0x0a
	FrameParamStateUsed      = 0x0b
)

// SIGPARAM parameter indices (EIP-8141).
const (
	SigParamSigner          = 0x00
	SigParamScheme          = 0x01
	SigParamMsg             = 0x02
	SigParamSignatureLength = 0x03
)

// TXTRACE parameters (EIP-7906). The count parameters take a zero second operand; the
// others take an index into the table the count describes.
const (
	TxTraceBalancesChanged   = 0x00
	TxTraceSlotsChanged      = 0x01
	TxTraceContractsDeployed = 0x02
	TxTraceBalanceAddress    = 0x03
	TxTraceBalanceBefore     = 0x04
	TxTraceBalanceAfter      = 0x05
	TxTraceSlotAddress       = 0x06
	TxTraceSlotKey           = 0x07
	TxTraceSlotBefore        = 0x08
	TxTraceSlotAfter         = 0x09
	TxTraceDeployedAddress   = 0x0a
	TxTraceDeployedCodeHash  = 0x0b
	TxTraceEventCount        = 0x0c
	TxTraceEventAddress      = 0x0d
	TxTraceEventTopicCount   = 0x0e
	TxTraceEventTopic0       = 0x0f
	TxTraceEventTopic1       = 0x10
	TxTraceEventTopic2       = 0x11
	TxTraceEventTopic3       = 0x12
	TxTraceEventDataLength   = 0x13
	TxTraceGasPreCharge      = 0x14
	TxTraceGasPayer          = 0x15
)

// TXDIFF parameters (EIP-7906), keyed by address and, for the slot parameters, a slot
// key or a per-address index; the topic parameters are keyed by a topic value instead.
const (
	TxDiffSlotBefore      = 0x00
	TxDiffSlotAfter       = 0x01
	TxDiffBalanceBefore   = 0x02
	TxDiffBalanceAfter    = 0x03
	TxDiffCodeHashBefore  = 0x04
	TxDiffCodeHashAfter   = 0x05
	TxDiffSlotCount       = 0x06
	TxDiffSlotIndex       = 0x07
	TxDiffEventCount      = 0x08
	TxDiffEventIndex      = 0x09
	TxDiffChangeFlags     = 0x0a
	TxDiffTopicEventCount = 0x0b
	TxDiffTopicEventIndex = 0x0c
)
