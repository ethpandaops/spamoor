package frametxfuzz

// Instructions and parameter indices introduced by the frame transaction EIPs.
//
// They live here rather than in txtypes because they describe the EVM, not the
// transaction: only code written to run inside a frame encodes them. The probe contract
// emits these bytes, and the assertion helpers name these indices.

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

	// OpRecentRootRefLoad is EIP-8272's RECENTROOTREFLOAD.
	//
	// It collides with OpSigDataCopy. EIP-8272 justifies 0xB5 with "EIP-8141 assigns
	// opcode 0xB4 to SIGPARAM. RECENTROOTREFLOAD uses 0xB5 to avoid that collision",
	// which was written before EIP-8141 grew SIGDATACOPY at 0xB5. On a chain running
	// both EIPs the opcode is assigned twice; the constants are kept distinct so a
	// consumer can say which one it meant.
	OpRecentRootRefLoad = 0xb5
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

// TXPARAM parameter indices added by the extension EIPs.
//
// TxParamLegacyNonce collides with EIP-8141's TxParamStateGasLeft: EIP-8250 claims 0x0C
// on the premise that "EIP-8141 assigns TXPARAM indices through 0x0B", which stopped
// being true when EIP-8141 added state_gas_left at 0x0C. Both names are defined so a
// consumer can be explicit about which reading it is testing.
const (
	// TxParamLegacyNonce is EIP-8250's pre-state sender account nonce.
	TxParamLegacyNonce = 0x0c

	// TxParamNonceKeyCount is EIP-8250's len(nonce_keys).
	TxParamNonceKeyCount = 0x0d

	// TxParamNonceKeysHash is EIP-8250's nonce_keys_hash.
	TxParamNonceKeysHash = 0x0e

	// TxParamRecentRootReferenceCount is EIP-8272's len(recent_root_references).
	TxParamRecentRootReferenceCount = 0x0f

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

// RECENTROOTREFLOAD field selectors (EIP-8272).
const (
	RecentRootFieldSourceID = 0
	RecentRootFieldSlot     = 1
	RecentRootFieldRoot     = 2
)
