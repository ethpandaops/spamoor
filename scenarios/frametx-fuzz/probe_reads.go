package frametxfuzz

import "github.com/ethereum/go-ethereum/common"

// The introspection sweep: operations that make each instruction EIP-8141 and its
// extensions introduce execute inside a frame. Returned values are discarded.
//
// What the sweep respects is where an instruction halts by definition, since a frame that
// halts stops exercising anything after it:
//
//   - FRAMEPARAM's status and gas-used parameters halt for the current or a later frame;
//   - SIGPARAM's resolved signer halts on an ARBITRARY entry and its raw length halts on
//     every other scheme;
//   - SIGDATACOPY is defined for ARBITRARY entries only;
//   - TXTRACE, TXDIFF and EVENTDATACOPY exist only inside a POST_TX frame, and their
//     indexed parameters halt past the end of the table they index.

// appendReads adds the introspection sweep for a frame at the given index.
func appendReads(script *ProbeScript, recipe *Recipe, frameIndex int) {
	appendTxParamReads(script, recipe)
	appendFrameParamReads(script, frameIndex)
	appendSigParamReads(script, recipe)
}

// appendTxParamReads sweeps the transaction-scoped parameters.
func appendTxParamReads(script *ProbeScript, recipe *Recipe) {
	for _, param := range []uint8{
		TxParamTxType, TxParamNonce, TxParamSender,
		TxParamGasTipCap, TxParamGasFeeCap, TxParamBlobFeeCap,
		TxParamMaxCost, TxParamBlobHashCount, TxParamSigHash,
		TxParamFrameCount, TxParamFrameIndex, TxParamSignatureCount,
		TxParamStateGasLeft,
	} {
		script.ReadTxParam(param)
	}

	if recipe.NonceKeys > 0 {
		// EIP-8250's indices, defined whatever key set the transaction selects.
		for _, param := range []uint8{TxParamLegacyNonce, TxParamNonceKeyCount, TxParamNonceKeysHash, TxParamNonceKey0} {
			script.ReadTxParam(param)
		}
	}
}

// appendFrameParamReads sweeps the frame-scoped parameters.
func appendFrameParamReads(script *ProbeScript, frameIndex int) {
	for _, param := range []uint8{
		FrameParamTarget, FrameParamExecutionLimit, FrameParamMode, FrameParamFlags,
		FrameParamDataLength, FrameParamAllowedScope, FrameParamAtomicBatch,
		FrameParamValue, FrameParamStateLimit,
	} {
		script.ReadFrameParam(param, frameIndex)
	}

	// Status and gas used are defined only for frames that have completed, so they are
	// read against the frame before this one.
	if frameIndex > 0 {
		for _, param := range []uint8{FrameParamStatus, FrameParamExecutionUsed, FrameParamStateUsed} {
			script.ReadFrameParam(param, frameIndex-1)
		}
	}

	// One word of another frame's data, which is what FRAMEDATALOAD is for.
	if frameIndex > 0 {
		script.ReadFrameData(frameIndex-1, 0)
	}
}

// appendSigParamReads sweeps the signature-scoped parameters. The sender's entry is at
// index 0; an ARBITRARY witness, when present, follows it and is the only entry whose raw
// bytes the EVM may read.
func appendSigParamReads(script *ProbeScript, recipe *Recipe) {
	script.ReadSigParam(SigParamScheme, 0)
	script.ReadSigParam(SigParamMsg, 0)
	script.ReadSigParam(SigParamSigner, 0)

	if !recipe.Witness {
		return
	}

	const witnessIndex = 1

	script.ReadSigParam(SigParamScheme, witnessIndex)
	script.ReadSigParam(SigParamSignatureLength, witnessIndex)
	script.ReadSigData(witnessIndex, 0)
}

// appendPostTxReads sweeps EIP-7906's assertion instructions from inside a POST_TX
// frame, against the transaction's own sender.
//
// The count parameters and the per-address lookups are defined for any transaction. The
// balance table always has an entry, since the payer's pre-charge is a balance change,
// so its index 0 is safe to read. The event table is only read when an earlier probe
// frame emitted a log, and even then a batch unroll can have discarded it: a halt there
// is the whole-body revert EIP-7906 specifies, which is as much a case as the read.
func appendPostTxReads(script *ProbeScript, recipe *Recipe, sender common.Address) {
	for _, param := range []uint8{
		TxTraceBalancesChanged, TxTraceSlotsChanged, TxTraceContractsDeployed,
		TxTraceEventCount, TxTraceGasPreCharge, TxTraceGasPayer,
	} {
		script.ReadTxTrace(param, 0)
	}

	for _, param := range []uint8{TxTraceBalanceAddress, TxTraceBalanceBefore, TxTraceBalanceAfter} {
		script.ReadTxTrace(param, 0)
	}

	for _, param := range []uint8{
		TxDiffBalanceBefore, TxDiffBalanceAfter, TxDiffCodeHashBefore, TxDiffCodeHashAfter,
		TxDiffSlotCount, TxDiffEventCount, TxDiffChangeFlags,
	} {
		script.ReadTxDiff(param, sender, common.Hash{})
	}

	script.ReadTxDiff(TxDiffSlotBefore, sender, common.HexToHash("0x01"))
	script.ReadTxDiff(TxDiffSlotAfter, sender, common.HexToHash("0x01"))

	if !recipe.emitsLog() {
		return
	}

	for _, param := range []uint8{TxTraceEventAddress, TxTraceEventTopicCount, TxTraceEventTopic0, TxTraceEventDataLength} {
		script.ReadTxTrace(param, 0)
	}

	script.ReadEventData(0, 0, 32)
}
