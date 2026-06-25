package evmfuzz

// deploy-size mode: structured generator that emits init code returning runtime
// of an exact target size, swept across EIP-7954 (max code size 24KiB->64KiB)
// and EIP-3860 (initcode size limit) boundaries. Bypasses the random Generate()
// loop entirely; every byte is fixed by (seed, txID) for full determinism.

// deploySizeInitOverhead is the fixed byte length of the init-code prologue that
// precedes the filler (the deployed runtime). All operands use fixed-width PUSH3
// so this is constant for every target size. Keep in sync with generateDeploySizeInit.
const deploySizeInitOverhead = 84

// deploySizeBoundaries are the EIP-7954/EIP-3860 critical runtime sizes (bytes):
// +/-1 around the old 24KiB limit, a midpoint, and the new 64KiB limit.
// Sizes above the cap (65536, 65537) deploy-fail on purpose; rejection must be
// consistent across clients, which is itself a valid differential target.
var deploySizeBoundaries = []int{
	24575, 24576, 24577, // old 24KiB limit -1 / limit / +1
	49152,               // ~48KiB midpoint
	65535, 65536, 65537, // new 64KiB (EIP-7954) limit -1 / limit / +1
}

// selectDeploySize deterministically picks the exact runtime size N for this
// (seed, txID): a boundary +/- a small delta, occasionally probing the EIP-3860
// initcode-size edge instead. Uses ONLY g.rng.
func (g *OpcodeGenerator) selectDeploySize() int {
	// 15% chance: land the initcode just below/at/above the EIP-3860 limit
	// (2*MAX_CODE_SIZE = 49152 pre-7954). initcode = overhead + runtime.
	if g.rng.Float64() < 0.15 {
		const initcodeLimit = 49152
		delta := g.rng.Intn(3) - 1 // -1, 0, +1
		return initcodeLimit - deploySizeInitOverhead + delta
	}

	base := deploySizeBoundaries[g.rng.Intn(len(deploySizeBoundaries))]
	switch g.rng.Intn(5) {
	case 0:
		return base - 1
	case 1:
		return base + 1
	default:
		return base // bias toward the exact boundary
	}
}

// generateDeploySizeInit emits init code that RETURNs exactly N bytes of benign
// runtime (0xfe INVALID filler) via CODECOPY of this code's own tail. The init
// code is fingerprinted with PUSH32 seed; PUSH32 txID (then POPped). Deterministic
// and bounded: only g.rng (via selectDeploySize) drives the single chosen size.
func (g *OpcodeGenerator) generateDeploySizeInit() []byte {
	n := g.selectDeploySize()
	if n < 1 {
		n = 1
	}
	g.deploySizeN = n

	const fillerStart = deploySizeInitOverhead // fixed prologue length; keep in sync

	g.bytecode = g.bytecode[:0]

	// fingerprint: PUSH32 seed; PUSH32 txID; then POP both to leave a clean stack
	g.pushSeedAndTxID()
	g.bytecode = append(g.bytecode, 0x50, 0x50) // POP POP

	push3 := func(v int) []byte { return []byte{0x62, byte(v >> 16), byte(v >> 8), byte(v)} }

	g.bytecode = append(g.bytecode, push3(n)...)           // PUSH3 N           (copy length)
	g.bytecode = append(g.bytecode, push3(fillerStart)...) // PUSH3 fillerStart (src offset)
	g.bytecode = append(g.bytecode, 0x5f)                  // PUSH0             (dest offset 0)
	g.bytecode = append(g.bytecode, 0x39)                  // CODECOPY
	g.bytecode = append(g.bytecode, push3(n)...)           // PUSH3 N           (return size)
	g.bytecode = append(g.bytecode, 0x5f)                  // PUSH0             (return offset 0)
	g.bytecode = append(g.bytecode, 0xf3)                  // RETURN

	// filler == deployed runtime: exactly N bytes of INVALID (0xfe)
	filler := make([]byte, n)
	for i := range filler {
		filler[i] = 0xfe
	}
	g.bytecode = append(g.bytecode, filler...)

	g.stackSize = 0
	return g.bytecode
}
