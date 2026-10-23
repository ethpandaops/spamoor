package frametxfuzz

// Instructions and parameter indices introduced by the frame transaction EIPs.
//
// They live here rather than in txtypes because they describe the EVM, not the
// transaction: only code written to run inside a frame encodes them. The probe contract
// emits these bytes, and the assertion helpers name these indices.
//
// The opcode bytes are those the EIPs assign: EIP-8141 takes 0xB0-0xB5 and EIP-7906
// 0xB6-0xB7. A byte above them that no enabled EIP allocates is an invalid instruction.

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

	// OpEventDataCopy is EIP-7906's EVENTDATACOPY: an emitted event's data, copied to
	// memory. Valid only inside an ASSERT frame.
	OpEventDataCopy = 0xb7
)

// TXPARAM parameter indices (EIP-8141). 0x05 and 0x06 went with the per-gas fee caps.
const (
	TxParamTxType         = 0x00
	TxParamNonce          = 0x01
	TxParamSender         = 0x02
	TxParamGasTipCap      = 0x03
	TxParamMaxFee         = 0x04
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
	TxTraceBalancesChanged = 0x00
	TxTraceSlotsChanged    = 0x01
	TxTraceCodesChanged    = 0x02
	TxTraceNoncesChanged   = 0x03
	TxTraceBalanceAddress  = 0x04
	TxTraceBalanceBefore   = 0x05
	TxTraceBalanceAfter    = 0x06
	TxTraceSlotAddress     = 0x07
	TxTraceSlotBefore      = 0x08
	TxTraceSlotAfter       = 0x09
	TxTraceSlotKey         = 0x0a
	TxTraceCodeAddress     = 0x0b
	TxTraceCodeBefore      = 0x0c
	TxTraceCodeAfter       = 0x0d
	TxTraceNonceAddress    = 0x0e
	TxTraceNonceBefore     = 0x0f
	TxTraceNonceAfter      = 0x10
	TxTraceEventCount      = 0x11
	TxTraceEventAddress    = 0x12
	TxTraceEventTopicCount = 0x13
	TxTraceEventTopic0     = 0x14
	TxTraceEventTopic1     = 0x15
	TxTraceEventTopic2     = 0x16
	TxTraceEventTopic3     = 0x17
	TxTraceEventDataLength = 0x18
	TxTraceGasPreCharge    = 0x19
	TxTraceGasPayer        = 0x1a
)
