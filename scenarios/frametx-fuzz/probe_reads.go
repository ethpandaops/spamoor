package frametxfuzz

// The introspection sweep: operations that make each instruction EIP-8141 and its
// extensions introduce execute inside a frame. Returned values are discarded.
//
// What the sweep respects is where an instruction halts by definition, since a frame that
// halts stops exercising anything after it:
//
//   - FRAMEPARAM's status and gas-used parameters halt for the current or a later frame;
//   - SIGPARAM's resolved signer halts on an ARBITRARY entry and its raw length halts on
//     every other scheme;
//   - SIGDATACOPY is defined for ARBITRARY entries only.

// appendReads adds the introspection sweep for a frame at the given index.
func appendReads(script *ProbeScript, recipe *Recipe, frameIndex int) {
	appendTxParamReads(script, recipe)
	appendFrameParamReads(script, frameIndex)
	appendSigParamReads(script, recipe)
	appendRootRefReads(script, recipe)
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
		// EIP-8250's indices. TxParamLegacyNonce shares 0x0C with EIP-8141's
		// state_gas_left, which both EIPs claim.
		for _, param := range []uint8{TxParamNonceKeyCount, TxParamNonceKeysHash, TxParamNonceKey0} {
			script.ReadTxParam(param)
		}
	}

	if recipe.RecentRoots > 0 {
		script.ReadTxParam(TxParamRecentRootReferenceCount)
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

// appendRootRefReads sweeps the declared recent root references.
func appendRootRefReads(script *ProbeScript, recipe *Recipe) {
	if recipe.RecentRoots == 0 {
		return
	}

	for _, field := range []uint8{RecentRootFieldSourceID, RecentRootFieldSlot, RecentRootFieldRoot} {
		script.ReadRootRef(0, field)
	}
}
