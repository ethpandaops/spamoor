package spamoor

import "math/bits"

// EIP-8037 / EIP-2780 constants.
//
// These mirror the values in github.com/ethereum/go-ethereum/params used by the
// Amsterdam (Glamsterdam EL) implementation. They are duplicated here so that
// spamoor can compute gas limits without depending on the typed gas APIs in
// core/vm (which change between devnets).
const (
	// targetStateGrowthPerYear targets 100 GiB annual state growth per EIP-8037.
	targetStateGrowthPerYear = 100 * 1024 * 1024 * 1024
	// cpsbSignificantBits caps the precision of cost_per_state_byte quantization.
	cpsbSignificantBits = 5
	// cpsbOffset is an additive bias applied before quantization (EIP-8037).
	cpsbOffset = 9578
	// cpsbFloor is the minimum cpsb value spamoor will use. It matches the
	// value that go-ethereum's glamsterdam-devnet-0 hardcodes in
	// core/evm.go (CostPerStateByte returns 1174 with the dynamic formula
	// commented out per the devnet-3 TODO). The dynamic formula yields 662
	// at 60M gas limit and 1174 at 100M, so on smaller chains the floor is
	// what keeps spamoor's gas estimates in sync with what the node actually
	// charges. Once geth activates the dynamic formula, natural cpsb at
	// 100M+ gas limits equals or exceeds this floor, making the clamp a
	// no-op.
	cpsbFloor = 1174

	// AccountCreationSize is the EIP-8037 charge unit count for creating a new
	// account (AccountCreationSize). State gas = AccountCreationSize * cpsb.
	AccountCreationSize = 112
	// AuthorizationCreationSize is the EIP-8037 charge unit count applied to
	// each EIP-7702 SetCode authorization that creates a delegator account.
	AuthorizationCreationSize = 23

	// callRegularGas is a generous estimate of the regular-gas cost a CALL
	// with value incurs inside the batcher contract: CallGasEIP150 (700) +
	// ColdAccountAccess (2600) + CallValueTransfer (9000) + loop opcodes + slack.
	callRegularGas = 12_500

	// txpoolBufferNum/Denom mirrors the 10/9 (≈111.1%) buffer the Amsterdam
	// txpool enforces on the regular-gas intrinsic (see
	// core/txpool/validation.go: `tx.Gas() < (intrGas.RegularGas*10)/9`).
	// We apply the same factor so tx.Gas == (base*10)/9 exactly meets the
	// pool's `<` check and passes.
	txpoolBufferNum   = 10
	txpoolBufferDenom = 9
)

// computeCostPerStateByte implements the EIP-8037 cost_per_state_byte formula.
// Returns 0 when gasLimit is 0 (unknown / pre-Amsterdam). Results are deterministic.
func computeCostPerStateByte(gasLimit uint64) uint64 {
	if gasLimit == 0 {
		return 0
	}
	// raw = ceil((gasLimit * 2_628_000) / (2 * TARGET_STATE_GROWTH_PER_YEAR))
	num := gasLimit*2_628_000 + (2*targetStateGrowthPerYear - 1)
	raw := num / (2 * targetStateGrowthPerYear)

	shifted := raw + cpsbOffset
	shift := max(bits.Len64(shifted)-cpsbSignificantBits, 0)
	quantized := (shifted >> shift) << shift
	var formula uint64 = 1
	if quantized > cpsbOffset {
		formula = quantized - cpsbOffset
	}
	// Clamp to the current hardcoded value in geth so spamoor's gas budgets
	// match the node's charging on devnets that haven't yet enabled the
	// dynamic formula. See cpsbFloor for background.
	return max(formula, cpsbFloor)
}
