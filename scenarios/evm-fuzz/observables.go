package evmfuzz

import "encoding/binary"

// Gas observability probes: stack-balanced LOG sequences that surface in-EVM gas
// consumption so cross-client receipt/log diffs localize a repricing divergence to
// a single operation. Targets gas-repricing EIPs (EIP-2929, EIP-7778, EIP-7976,
// EIP-7981, EIP-8037, EIP-8038). Off by default; composes with any --fuzz-mode.

const (
	probeKindCheckpoint = 0xC0
	probeKindDelta      = 0xDE

	// Modeled gas charged per probe so the generator's gas budget tracks the real
	// EVM cost and the deploy stays within --gaslimit. LOG1 of one 32-byte word with
	// a single topic dominates: 375 base + 375 topic + 8*32 data = 1006.
	gasLOG1Word           = 375 + 375 + 8*32
	gasCheckpointProbe    = gasLOG1Word + 24 // + GAS/PUSH/MSTORE scaffold (~18), incl. one 32B mem expansion
	gasDeltaProbeScaffold = gasLOG1Word + 32 // + second GAS, SWAP1, SUB scaffold (~26)
	// Conservative EIP-2929 cold-access cost for the BALANCE measured op so the
	// access-cost probe is never undercounted (BALANCE(self) is warm in practice).
	gasColdAccountAccess = 2600
)

// SetGasProbes enables/disables in-EVM gas observability probes.
func (g *OpcodeGenerator) SetGasProbes(enabled bool) {
	g.gasProbes = enabled
}

// makeProbeTopic builds a self-identifying 32-byte LOG topic:
// [0:8]=txID, [8]=kind, [28:32]=seq. Pure function of generator state, so deterministic.
func (g *OpcodeGenerator) makeProbeTopic(kind byte, seq uint32) []byte {
	t := make([]byte, 32)
	binary.BigEndian.PutUint64(t[0:8], g.txID)
	t[8] = kind
	binary.BigEndian.PutUint32(t[28:32], seq)
	return t
}

// emitGasCheckpoint emits a stack-balanced gas checkpoint: the GAS value is stored
// at mem[0:32] and LOG1'd under a self-identifying topic. Net stack effect: 0.
func (g *OpcodeGenerator) emitGasCheckpoint(seq uint32) []byte {
	bytecode := []byte{
		0x5a, // GAS    -> [.. gas]
		0x5f, // PUSH0  -> [.. gas 0]
		0x52, // MSTORE -> [..]   mem[0:32]=gas
		0x7f, // PUSH32 topic
	}
	bytecode = append(bytecode, g.makeProbeTopic(probeKindCheckpoint, seq)...)
	bytecode = append(bytecode,
		0x60, 0x20, // PUSH1 0x20 (size)
		0x5f, // PUSH0 (offset)
		0xa1, // LOG1  -> [..]
	)
	return bytecode
}

// emitGasDelta wraps a net-zero-stack measured sub-sequence with GAS reads, computes
// gasBefore-gasAfter, stores it at mem[0:32] and LOG1's it under a delta topic.
// The measured op MUST be net-zero on the stack; net stack effect of the result: 0.
// Note: the reported delta includes the trailing GAS opcode's 2 gas (a constant
// identical across clients, so it does not pollute cross-client diffs).
func (g *OpcodeGenerator) emitGasDelta(measured []byte, seq uint32) []byte {
	bytecode := []byte{0x5a} // GAS -> [.. gB]
	bytecode = append(bytecode, measured...)
	bytecode = append(bytecode,
		0x5a, // GAS    -> [.. gB gA]
		0x90, // SWAP1  -> [.. gA gB]
		0x03, // SUB    -> [.. gB-gA]
		0x5f, // PUSH0  -> [.. delta 0]
		0x52, // MSTORE -> [..]   mem[0:32]=delta
		0x7f, // PUSH32 topic
	)
	bytecode = append(bytecode, g.makeProbeTopic(probeKindDelta, seq)...)
	bytecode = append(bytecode,
		0x60, 0x20, // PUSH1 0x20 (size)
		0x5f, // PUSH0 (offset)
		0xa1, // LOG1  -> [..]
	)
	return bytecode
}

// pickNetZeroMeasuredOp returns a single instruction with net-zero stack effect and
// zero stack-input requirement (always safe to splice), chosen deterministically,
// along with its conservative modeled gas cost. These touch opcodes whose pricing the
// target EIPs change, so measuring each in isolation localizes a mispricing to one opcode.
func (g *OpcodeGenerator) pickNetZeroMeasuredOp() ([]byte, uint64) {
	switch g.rng.Intn(4) {
	case 0:
		return []byte{0x3a, 0x50}, 2 + 2 // GASPRICE, POP
	case 1:
		return []byte{0x46, 0x50}, 2 + 2 // CHAINID, POP
	case 2:
		return []byte{0x30, 0x31, 0x50}, 2 + gasColdAccountAccess + 2 // ADDRESS, BALANCE(self), POP -> EIP-2929 access cost
	default:
		return []byte{0x47, 0x50}, 5 + 2 // SELFBALANCE, POP
	}
}

// tryEmitGasProbe emits a gas observability probe if budgets allow (70% cheap
// checkpoint, 30% delta wrapping one net-zero op). Returns true on emit. All emitted
// sequences are net-zero on the stack, so g.stackSize is intentionally untouched.
// Fails fast: on any tight budget it returns false without emitting a partial probe.
func (g *OpcodeGenerator) tryEmitGasProbe() bool {
	// Cheapest probe is a checkpoint; bail before consuming rng if even that won't fit.
	if g.currentGas+gasCheckpointProbe > g.maxGas {
		return false
	}

	var probe []byte
	var cost uint64
	if g.rng.Float64() < 0.7 {
		probe = g.emitGasCheckpoint(g.probeSeq)
		cost = gasCheckpointProbe
	} else {
		measured, measuredGas := g.pickNetZeroMeasuredOp()
		probe = g.emitGasDelta(measured, g.probeSeq)
		cost = gasDeltaProbeScaffold + measuredGas
	}

	if g.currentGas+cost > g.maxGas ||
		len(g.bytecode)+len(probe) > g.maxSize-32 ||
		g.opcodeCount+g.countOpcodesInBytecode(probe) > g.maxOpcodeCount-10 {
		return false
	}

	g.bytecode = append(g.bytecode, probe...)
	g.currentGas += cost
	g.opcodeCount += g.countOpcodesInBytecode(probe)
	g.probeSeq++
	return true
}
