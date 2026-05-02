package calltxfuzz

import (
	"encoding/binary"
	"fmt"
	"sort"

	"github.com/ethereum/go-ethereum/common"

	evmfuzz "github.com/ethpandaops/spamoor/scenarios/evm-fuzz"
)

// CallFuzzGenerator generates runtime bytecode for deployed contracts
// with delegation-focused patterns including cross-contract calls,
// storage/transient storage patterns, and identity checks.
// When poolAddresses is non-empty, templates embed actual pool contract
// addresses for cross-contract calls. Otherwise, falls back to
// ADDRESS/CALLER/ORIGIN as call targets.
type CallFuzzGenerator struct {
	rng              *evmfuzz.DeterministicRNG
	transformer      *evmfuzz.InputTransformer
	stackSize        int
	bytecode         []byte
	jumpTargets      []int
	jumpPlaceholders []int
	maxGas           uint64
	currentGas       uint64
	maxSize          int
	opcodeCount      int
	maxOpcodeCount   int
	opcodeInfos      map[uint16]*evmfuzz.OpcodeInfo
	validOpcodes     []*evmfuzz.OpcodeInfo
	invalidOpcodes   []byte
	stackBuilders    []*evmfuzz.OpcodeInfo
	txID             uint64
	baseSeed         string
	poolAddresses    []common.Address
}

// NewCallFuzzGenerator creates a new generator for runtime bytecode
// with delegation-focused opcode weights.
func NewCallFuzzGenerator(
	txID uint64,
	baseSeed string,
	maxSize int,
	maxGas uint64,
	poolAddresses []common.Address,
) *CallFuzzGenerator {
	rng := evmfuzz.NewDeterministicRNGWithSeed(txID, baseSeed)

	g := &CallFuzzGenerator{
		rng:              rng,
		transformer:      evmfuzz.NewInputTransformer(rng),
		stackSize:        0,
		bytecode:         make([]byte, 0, maxSize),
		jumpTargets:      make([]int, 0, 16),
		jumpPlaceholders: make([]int, 0, 16),
		maxGas:           maxGas,
		currentGas:       0,
		maxSize:          maxSize,
		opcodeCount:      0,
		maxOpcodeCount:   maxSize * 10,
		opcodeInfos:      make(map[uint16]*evmfuzz.OpcodeInfo, 128),
		txID:             txID,
		baseSeed:         baseSeed,
		poolAddresses:    poolAddresses,
	}

	g.initializeOpcodes()
	g.buildValidOpcodeList()
	g.buildInvalidOpcodeList()
	g.buildStackBuilderList()

	return g
}

// initializeOpcodes sets up opcode definitions with delegation-focused weights
func (g *CallFuzzGenerator) initializeOpcodes() {
	opcodes := getCalltxFuzzOpcodeDefinitions()

	// Add DUP1-DUP16 and SWAP1-SWAP16
	opcodes = append(opcodes, evmfuzz.GetDupOpcodeDefinitions()...)
	opcodes = append(opcodes, evmfuzz.GetSwapOpcodeDefinitions()...)

	// Add generator-specific opcodes
	generatorOpcodes := []*evmfuzz.OpcodeInfo{
		// JUMP/JUMPI need generator for target tracking
		op("JUMP", 0x56, 0, 0, 8, g.generateJump, 1.0),
		op("JUMPI", 0x57, 1, 0, 10, g.generateJumpi, 1.0),
	}
	opcodes = append(opcodes, generatorOpcodes...)

	// Add PUSH1-PUSH32
	for i := 1; i <= 32; i++ {
		pushOpcode := uint16(0x5f + i)
		pushSize := i
		opcodes = append(opcodes, op(
			fmt.Sprintf("PUSH%d", i), pushOpcode, 0, 1, 3,
			g.makePushTemplate(pushOpcode, pushSize), 1.0,
		))
	}

	// Store in map
	for _, op := range opcodes {
		// Skip opcodes with nil templates (precompile placeholders)
		if op.Template == nil {
			continue
		}
		g.opcodeInfos[op.Opcode] = op
	}
}

// buildValidOpcodeList builds a sorted list of usable opcodes
func (g *CallFuzzGenerator) buildValidOpcodeList() {
	g.validOpcodes = g.validOpcodes[:0]
	for _, op := range g.opcodeInfos {
		g.validOpcodes = append(g.validOpcodes, op)
	}
	sort.Slice(g.validOpcodes, func(i, j int) bool {
		return g.validOpcodes[i].Opcode < g.validOpcodes[j].Opcode
	})
}

// buildInvalidOpcodeList finds opcodes not in the valid set
func (g *CallFuzzGenerator) buildInvalidOpcodeList() {
	validSet := make(map[byte]bool, len(g.opcodeInfos))
	for _, op := range g.opcodeInfos {
		if op.Opcode <= 0xFF {
			validSet[byte(op.Opcode)] = true
		}
	}
	g.invalidOpcodes = g.invalidOpcodes[:0]
	for opcode := 0; opcode <= 0xFF; opcode++ {
		if !validSet[byte(opcode)] {
			g.invalidOpcodes = append(g.invalidOpcodes, byte(opcode))
		}
	}
}

// buildStackBuilderList caches opcodes that push one item without consuming any
func (g *CallFuzzGenerator) buildStackBuilderList() {
	g.stackBuilders = g.stackBuilders[:0]
	for _, op := range g.opcodeInfos {
		if op.StackInput != 0 || op.StackOutput != 1 {
			continue
		}
		if op.Opcode >= 0x100 { // skip precompiles
			continue
		}
		if op.Opcode >= 0x60 && op.Opcode <= 0x7f { // skip PUSH1-PUSH32
			continue
		}
		g.stackBuilders = append(g.stackBuilders, op)
	}
	sort.Slice(g.stackBuilders, func(i, j int) bool {
		return g.stackBuilders[i].Opcode < g.stackBuilders[j].Opcode
	})
}

// Generate produces runtime bytecode (not init code).
// The output is meant to be deployed via an init code wrapper and
// then called through Type 2/4/6 transactions.
func (g *CallFuzzGenerator) Generate() []byte {
	g.bytecode = g.bytecode[:0]
	g.stackSize = 0
	g.jumpTargets = g.jumpTargets[:0]
	g.jumpPlaceholders = g.jumpPlaceholders[:0]
	g.currentGas = 0
	g.opcodeCount = 0

	// Push seed and txID for deterministic initial stack state
	g.pushSeedAndTxID()
	g.opcodeCount += 2

	// 95% chance: emit an initial gas bailout that returns calldata.
	// When a contract is entered with low gas (e.g. deep in a call chain),
	// this immediately copies calldata to memory and RETURNs it,
	// preventing an out-of-gas revert.
	if g.rng.Float64() < 0.95 {
		g.emitCalldataBailout()
	}

	// Main generation loop
	instrSinceBailout := 0
	for len(g.bytecode) < g.maxSize-32 &&
		g.currentGas < g.maxGas-1000 &&
		g.opcodeCount < g.maxOpcodeCount-10 {

		// 80% chance to inject a gas bailout check every ~8-15 instructions.
		// If GAS < 15000, store top-of-stack to memory and RETURN it.
		// This prevents deep call chains from running out of gas.
		instrSinceBailout++
		if instrSinceBailout > 8 && g.rng.Float64() < 0.80 {
			g.emitGasBailout()
			instrSinceBailout = 0
		}

		// 20% chance to place JUMPDESTs when we have few targets
		if len(g.jumpTargets) < 10 && g.rng.Float64() < 0.2 {
			pc := len(g.bytecode)
			g.bytecode = append(g.bytecode, 0x5b) // JUMPDEST
			g.jumpTargets = append(g.jumpTargets, pc)
			g.currentGas++
			g.opcodeCount++
			// Reset stack assumption — jumps may arrive with any stack depth.
			// Conservative reset ensures addStackItems fills as needed.
			g.stackSize = 0
			continue
		}

		// 40% chance for delegation-specific template
		if g.rng.Float64() < 0.40 {
			if g.generateDelegationTemplate() {
				continue
			}
		}

		if !g.generateNextInstruction() {
			break
		}
	}

	// Final gas bailout before STOP
	g.emitGasBailout()

	// Terminate safely
	if len(g.bytecode) < g.maxSize-1 && g.opcodeCount < g.maxOpcodeCount {
		g.bytecode = append(g.bytecode, 0x00) // STOP
		g.opcodeCount++
	}

	g.fixJumpTargets()
	return g.bytecode
}

// generateDelegationTemplate injects a delegation-specific bytecode pattern.
// When poolAddresses is non-empty, templates embed actual pool contract
// addresses via PUSH20 for cross-contract calls. Otherwise, falls back
// to ADDRESS/CALLER/ORIGIN as call targets.
func (g *CallFuzzGenerator) generateDelegationTemplate() bool {
	// Weighted distribution: CALL and DELEGATECALL are preferred over
	// STATICCALL to reduce write-protection errors.
	choice := g.rng.Intn(12)
	switch choice {
	case 0, 1:
		return g.generatePoolCall()
	case 2, 3:
		return g.generatePoolDelegateCall()
	case 4:
		return g.generatePoolStaticCall()
	case 5:
		return g.generateCrossCallSequence()
	case 6:
		return g.generateStorageCallStorage()
	case 7:
		return g.generateTransientCallPattern()
	case 8:
		return g.generateCallWithReturnData()
	case 9:
		return g.generateDelegateCallChain()
	case 10:
		return g.generateStoragePattern()
	case 11:
		return g.generateIdentityCheck()
	}
	return false
}

// randomSubcallGas returns PUSH3 bytecode for a random gas limit (200000-800000).
// Minimum 200k ensures enough gas for the reentrancy sentry (2300 gas)
// plus the callee's initial gas bailout check and some useful execution.
func (g *CallFuzzGenerator) randomSubcallGas() []byte {
	gas := uint32(200000 + g.rng.Intn(600000))
	return []byte{0x62, byte(gas >> 16), byte(gas >> 8), byte(gas)}
}

// pushCallTarget emits bytecode to push a call target address onto the stack.
// When pool addresses are available, embeds a random peer address via PUSH20.
// Otherwise, falls back to ADDRESS (0x30), CALLER (0x33), or ORIGIN (0x32).
func (g *CallFuzzGenerator) pushCallTarget() []byte {
	if len(g.poolAddresses) > 0 {
		addr := g.poolAddresses[g.rng.Intn(len(g.poolAddresses))]
		bc := make([]byte, 21)
		bc[0] = 0x73 // PUSH20
		copy(bc[1:], addr.Bytes())
		return bc
	}
	// Fallback: use dynamic address opcode
	targets := []byte{0x30, 0x33, 0x32} // ADDRESS, CALLER, ORIGIN
	return []byte{targets[g.rng.Intn(len(targets))]}
}

// pushDistinctTargets returns bytecode snippets for n distinct call targets.
// Uses Fisher-Yates shuffle on pool addresses when available.
func (g *CallFuzzGenerator) pushDistinctTargets(n int) [][]byte {
	if len(g.poolAddresses) == 0 || n <= 0 {
		result := make([][]byte, n)
		for i := 0; i < n; i++ {
			result[i] = g.pushCallTarget()
		}
		return result
	}

	// Fisher-Yates partial shuffle for n distinct addresses
	count := n
	if count > len(g.poolAddresses) {
		count = len(g.poolAddresses)
	}

	indices := make([]int, len(g.poolAddresses))
	for i := range indices {
		indices[i] = i
	}
	for i := 0; i < count; i++ {
		j := i + g.rng.Intn(len(indices)-i)
		indices[i], indices[j] = indices[j], indices[i]
	}

	result := make([][]byte, n)
	for i := 0; i < n; i++ {
		idx := indices[i%count]
		addr := g.poolAddresses[idx]
		bc := make([]byte, 21)
		bc[0] = 0x73 // PUSH20
		copy(bc[1:], addr.Bytes())
		result[i] = bc
	}
	return result
}

// generatePoolCall emits CALL to a pool contract address.
// Cross-contract call with limited gas to prevent deep recursion.
// Stack: +1 (CALL result)
func (g *CallFuzzGenerator) generatePoolCall() bool {
	gasBytes := g.randomSubcallGas()
	targetBytes := g.pushCallTarget()

	bc := []byte{
		0x60, 0x20, // PUSH1 retLength=32
		0x5f, // PUSH0 retOffset=0
		0x5f, // PUSH0 argsLength=0
		0x5f, // PUSH0 argsOffset=0
		0x5f, // PUSH0 value=0
	}
	bc = append(bc, targetBytes...) // PUSH20 <addr> or ADDRESS/CALLER/ORIGIN
	bc = append(bc, gasBytes...)    // PUSH2 gas
	bc = append(bc, 0xf1)           // CALL

	if len(g.bytecode)+len(bc) > g.maxSize-32 {
		return false
	}
	g.bytecode = append(g.bytecode, bc...)
	g.stackSize++
	g.currentGas += 100
	g.opcodeCount += countOpcodesInBytecode(bc)
	return true
}

// generatePoolDelegateCall emits DELEGATECALL to a pool contract.
// Tests storage delegation where the callee executes in our storage context.
// Stack: +1 (DELEGATECALL result)
func (g *CallFuzzGenerator) generatePoolDelegateCall() bool {
	gasBytes := g.randomSubcallGas()
	targetBytes := g.pushCallTarget()

	// DELEGATECALL: 6 args (gas, addr, argsOff, argsLen, retOff, retLen)
	bc := []byte{
		0x60, 0x20, // PUSH1 retLength=32
		0x5f, // PUSH0 retOffset=0
		0x5f, // PUSH0 argsLength=0
		0x5f, // PUSH0 argsOffset=0
	}
	bc = append(bc, targetBytes...) // PUSH20 <addr> or ADDRESS/CALLER/ORIGIN
	bc = append(bc, gasBytes...)    // PUSH2 gas
	bc = append(bc, 0xf4)           // DELEGATECALL

	if len(g.bytecode)+len(bc) > g.maxSize-32 {
		return false
	}
	g.bytecode = append(g.bytecode, bc...)
	g.stackSize++
	g.currentGas += 100
	g.opcodeCount += countOpcodesInBytecode(bc)
	return true
}

// generatePoolStaticCall emits STATICCALL to a pool contract.
// Tests read-only cross-contract calls in delegation context.
// Stack: +1 (STATICCALL result)
func (g *CallFuzzGenerator) generatePoolStaticCall() bool {
	gasBytes := g.randomSubcallGas()
	targetBytes := g.pushCallTarget()

	// STATICCALL: 6 args (gas, addr, argsOff, argsLen, retOff, retLen)
	bc := []byte{
		0x60, 0x20, // PUSH1 retLength=32
		0x5f, // PUSH0 retOffset=0
		0x5f, // PUSH0 argsLength=0
		0x5f, // PUSH0 argsOffset=0
	}
	bc = append(bc, targetBytes...) // PUSH20 <addr> or ADDRESS/CALLER/ORIGIN
	bc = append(bc, gasBytes...)    // PUSH2 gas
	bc = append(bc, 0xfa)           // STATICCALL

	if len(g.bytecode)+len(bc) > g.maxSize-32 {
		return false
	}
	g.bytecode = append(g.bytecode, bc...)
	g.stackSize++
	g.currentGas += 100
	g.opcodeCount += countOpcodesInBytecode(bc)
	return true
}

// generateCrossCallSequence emits 2-3 calls to distinct pool contracts
// in sequence. Each call targets a different contract using CALL,
// STATICCALL, or DELEGATECALL. No self-calls, breaking recursion loops.
// Stack: +numCalls (one result per call)
func (g *CallFuzzGenerator) generateCrossCallSequence() bool {
	numCalls := 2 + g.rng.Intn(2) // 2 or 3
	targets := g.pushDistinctTargets(numCalls)

	// Call opcodes: prefer CALL/DELEGATECALL over STATICCALL to reduce
	// write-protection errors (callees often execute SSTORE/TSTORE/LOG).
	callOps := []byte{0xf1, 0xf4, 0xf1, 0xf4, 0xfa}

	var bc []byte
	for i := 0; i < numCalls; i++ {
		callOp := callOps[g.rng.Intn(len(callOps))]
		gasBytes := g.randomSubcallGas()

		if callOp == 0xf1 { // CALL (7 args)
			bc = append(bc,
				0x5f, // PUSH0 retLength
				0x5f, // PUSH0 retOffset
				0x5f, // PUSH0 argsLength
				0x5f, // PUSH0 argsOffset
				0x5f, // PUSH0 value
			)
			bc = append(bc, targets[i]...)
			bc = append(bc, gasBytes...)
			bc = append(bc, callOp)
		} else { // DELEGATECALL/STATICCALL (6 args)
			bc = append(bc,
				0x5f, // PUSH0 retLength
				0x5f, // PUSH0 retOffset
				0x5f, // PUSH0 argsLength
				0x5f, // PUSH0 argsOffset
			)
			bc = append(bc, targets[i]...)
			bc = append(bc, gasBytes...)
			bc = append(bc, callOp)
		}
	}

	if len(g.bytecode)+len(bc) > g.maxSize-32 {
		return false
	}
	g.bytecode = append(g.bytecode, bc...)
	g.stackSize += numCalls
	g.currentGas += uint64(numCalls) * 100
	g.opcodeCount += countOpcodesInBytecode(bc)
	return true
}

// generateStorageCallStorage emits SSTORE, CALL to a pool contract,
// then SLOAD with the same key. Tests whether cross-contract calls
// affect storage in the calling contract's context.
// Stack: +1 (SLOAD result)
func (g *CallFuzzGenerator) generateStorageCallStorage() bool {
	key := byte(g.rng.Intn(32))
	val := byte(g.rng.Intn(256))
	gasBytes := g.randomSubcallGas()
	targetBytes := g.pushCallTarget()

	bc := []byte{
		// SSTORE: store value at key
		0x60, val, // PUSH1 value
		0x60, key, // PUSH1 key
		0x55, // SSTORE (consumes 2)
		// CALL to pool contract
		0x5f, // PUSH0 retLength
		0x5f, // PUSH0 retOffset
		0x5f, // PUSH0 argsLength
		0x5f, // PUSH0 argsOffset
		0x5f, // PUSH0 value
	}
	bc = append(bc, targetBytes...) // PUSH20 <addr>
	bc = append(bc, gasBytes...)    // PUSH2 gas
	bc = append(bc,
		0xf1,      // CALL (consumes 7, pushes 1)
		0x50,      // POP (discard call result)
		0x60, key, // PUSH1 same key
		0x54, // SLOAD (consumes 1, pushes 1)
	)

	if len(g.bytecode)+len(bc) > g.maxSize-32 {
		return false
	}
	g.bytecode = append(g.bytecode, bc...)
	g.stackSize++ // net: +1
	g.currentGas += 306
	g.opcodeCount += countOpcodesInBytecode(bc)
	return true
}

// generateTransientCallPattern emits TSTORE, CALL to a pool contract,
// then TLOAD with the same key. Tests transient storage isolation
// across cross-contract calls within the same transaction.
// Stack: +1 (TLOAD result)
func (g *CallFuzzGenerator) generateTransientCallPattern() bool {
	key := byte(g.rng.Intn(32))
	val := byte(g.rng.Intn(256))
	gasBytes := g.randomSubcallGas()
	targetBytes := g.pushCallTarget()

	bc := []byte{
		// TSTORE: store value at key
		0x60, val, // PUSH1 value
		0x60, key, // PUSH1 key
		0x5d, // TSTORE (consumes 2)
		// CALL to pool contract
		0x5f, // PUSH0 retLength
		0x5f, // PUSH0 retOffset
		0x5f, // PUSH0 argsLength
		0x5f, // PUSH0 argsOffset
		0x5f, // PUSH0 value
	}
	bc = append(bc, targetBytes...) // PUSH20 <addr>
	bc = append(bc, gasBytes...)    // PUSH2 gas
	bc = append(bc,
		0xf1,      // CALL (consumes 7, pushes 1)
		0x50,      // POP (discard call result)
		0x60, key, // PUSH1 same key
		0x5c, // TLOAD (consumes 1, pushes 1)
	)

	if len(g.bytecode)+len(bc) > g.maxSize-32 {
		return false
	}
	g.bytecode = append(g.bytecode, bc...)
	g.stackSize++ // net: +1
	g.currentGas += 306
	g.opcodeCount += countOpcodesInBytecode(bc)
	return true
}

// generateCallWithReturnData emits CALL to a pool contract, then
// copies and loads the return data. Tests RETURNDATASIZE/RETURNDATACOPY
// across cross-contract boundaries.
// Stack: +1 (MLOAD result from return data)
func (g *CallFuzzGenerator) generateCallWithReturnData() bool {
	gasBytes := g.randomSubcallGas()
	targetBytes := g.pushCallTarget()

	bc := []byte{
		// CALL to pool contract with return buffer
		0x60, 0x20, // PUSH1 retLength=32
		0x5f, // PUSH0 retOffset=0
		0x5f, // PUSH0 argsLength=0
		0x5f, // PUSH0 argsOffset=0
		0x5f, // PUSH0 value=0
	}
	bc = append(bc, targetBytes...) // PUSH20 <addr>
	bc = append(bc, gasBytes...)    // PUSH2 gas
	bc = append(bc,
		0xf1, // CALL (consumes 7, pushes 1)
		0x50, // POP (discard success flag)
		// Copy return data to memory
		0x3d, // RETURNDATASIZE (push size)
		0x5f, // PUSH0 offset=0
		0x5f, // PUSH0 destOffset=0
		0x3e, // RETURNDATACOPY (consumes 3)
		// Load from memory
		0x5f, // PUSH0 offset=0
		0x51, // MLOAD (consumes 1, pushes 1)
	)

	if len(g.bytecode)+len(bc) > g.maxSize-32 {
		return false
	}
	g.bytecode = append(g.bytecode, bc...)
	g.stackSize++ // net: +1
	g.currentGas += 112
	g.opcodeCount += countOpcodesInBytecode(bc)
	return true
}

// generateDelegateCallChain emits DELEGATECALL or STATICCALL to a
// pool contract. Tests delegation chains where the callee executes
// in the caller's storage context.
// Stack: +1 (DELEGATECALL/STATICCALL result)
func (g *CallFuzzGenerator) generateDelegateCallChain() bool {
	gasBytes := g.randomSubcallGas()
	targetBytes := g.pushCallTarget()

	// Strongly prefer DELEGATECALL over STATICCALL to reduce
	// write-protection errors (callees execute SSTORE/TSTORE/LOG).
	useDelegatecall := g.rng.Float64() < 0.9

	// Both use 6 args: gas, addr, argsOff, argsLen, retOff, retLen
	bc := []byte{
		0x60, 0x20, // PUSH1 retLength=32
		0x5f, // PUSH0 retOffset=0
		0x5f, // PUSH0 argsLength=0
		0x5f, // PUSH0 argsOffset=0
	}
	bc = append(bc, targetBytes...) // PUSH20 <addr>
	bc = append(bc, gasBytes...)    // PUSH2 gas
	if useDelegatecall {
		bc = append(bc, 0xf4) // DELEGATECALL
	} else {
		bc = append(bc, 0xfa) // STATICCALL
	}

	if len(g.bytecode)+len(bc) > g.maxSize-32 {
		return false
	}
	g.bytecode = append(g.bytecode, bc...)
	g.stackSize++
	g.currentGas += 100
	g.opcodeCount += countOpcodesInBytecode(bc)
	return true
}

// generateStoragePattern emits SSTORE key, SLOAD key
// using deterministic keys for cross-delegation verification.
// Stack: +1 (SLOAD result)
func (g *CallFuzzGenerator) generateStoragePattern() bool {
	key := byte(g.rng.Intn(32))
	val := byte(g.rng.Intn(256))
	bc := []byte{
		0x60, val, // PUSH1 value
		0x60, key, // PUSH1 key
		0x55,      // SSTORE
		0x60, key, // PUSH1 key
		0x54, // SLOAD
	}
	if len(g.bytecode)+len(bc) > g.maxSize-32 {
		return false
	}
	g.bytecode = append(g.bytecode, bc...)
	g.stackSize++ // net: 2-2 + 1-1+1 = +1
	g.currentGas += 206
	g.opcodeCount += 5
	return true
}

// generateIdentityCheck emits ADDRESS, CALLER, EQ for delegation context
// inspection. In delegated calls, ADDRESS != CALLER reveals the delegation.
// Stack: +1 (EQ result)
func (g *CallFuzzGenerator) generateIdentityCheck() bool {
	bc := []byte{
		0x30, // ADDRESS
		0x33, // CALLER
		0x14, // EQ
	}
	if len(g.bytecode)+len(bc) > g.maxSize-32 {
		return false
	}
	g.bytecode = append(g.bytecode, bc...)
	g.stackSize++ // net: +2 -2 +1 = +1
	g.currentGas += 7
	g.opcodeCount += 3
	return true
}

// emitCalldataBailout injects a gas check at the start of contract execution:
// if GAS < 15000, copy the full calldata to memory and RETURN it.
// This is placed near the bytecode start so contracts entered with very
// little gas (e.g. deep in a call chain) immediately return calldata
// rather than reverting with out-of-gas.
//
// Bytecode emitted (17 bytes):
//
//	GAS                  ; push remaining gas
//	PUSH2 0x3a98         ; push 15000
//	LT                   ; 15000 < gas → 1 when gas is sufficient
//	PUSH2 <skip>         ; jump target past RETURN
//	JUMPI                ; skip if gas >= 15000
//	--- bailout path ---
//	CALLDATASIZE         ; push calldata length
//	PUSH0                ; source offset (0)
//	PUSH0                ; dest offset (0)
//	CALLDATACOPY         ; copy calldata to memory[0..]
//	CALLDATASIZE         ; push calldata length again
//	PUSH0                ; memory offset (0)
//	RETURN               ; return calldata
//	JUMPDEST             ; <skip> landing
func (g *CallFuzzGenerator) emitCalldataBailout() {
	skipPC := len(g.bytecode) + 16 // offset of JUMPDEST after RETURN

	bc := []byte{
		0x5a,             // GAS
		0x61, 0x3a, 0x98, // PUSH2 15000
		0x10,                                  // LT (15000 < gas → 1 when gas high)
		0x61, byte(skipPC >> 8), byte(skipPC), // PUSH2 <skip>
		0x57, // JUMPI (jump if gas >= 15000)
		0x36, // CALLDATASIZE
		0x5f, // PUSH0 (srcOffset=0)
		0x5f, // PUSH0 (destOffset=0)
		0x37, // CALLDATACOPY
		0x36, // CALLDATASIZE
		0x5f, // PUSH0 (memOffset=0)
		0xf3, // RETURN
		0x5b, // JUMPDEST <skip>
	}

	if len(g.bytecode)+len(bc) > g.maxSize-32 {
		return
	}

	g.bytecode = append(g.bytecode, bc...)
	// Stack effect: GAS+PUSH2+LT+ISZERO+PUSH2+JUMPI all consumed.
	// On the non-bailout path, stack is unchanged.
	g.currentGas += 20
	g.opcodeCount += countOpcodesInBytecode(bc)
}

// emitGasBailout injects a gas check: if GAS < 15000, store the current
// top-of-stack value to memory and RETURN it. This prevents execution from
// running out of gas in deep call chains by gracefully returning early.
//
// Bytecode emitted (16 bytes):
//
//	GAS                  ; push remaining gas
//	PUSH2 0x3a98         ; push 15000
//	LT                   ; 15000 < gas → 1 when gas is sufficient
//	PUSH2 <skip>         ; jump target past RETURN
//	JUMPI                ; skip if gas >= 15000
//	--- bailout path ---
//	PUSH0                ; memory offset
//	MSTORE               ; store TOS at memory[0] (consumes 1 stack item)
//	PUSH1 0x20           ; return 32 bytes
//	PUSH0                ; from memory offset 0
//	RETURN
//	JUMPDEST             ; <skip> landing
func (g *CallFuzzGenerator) emitGasBailout() {
	// Need at least 1 stack item for MSTORE; ensure it.
	if g.stackSize < 1 {
		g.bytecode = append(g.bytecode, 0x5f) // PUSH0
		g.stackSize++
		g.currentGas += 2
		g.opcodeCount++
	}

	skipPC := len(g.bytecode) + 15 // offset of JUMPDEST after RETURN

	bc := []byte{
		0x5a,             // GAS
		0x61, 0x3a, 0x98, // PUSH2 15000
		0x10,                                  // LT (15000 < gas → 1 when gas high)
		0x61, byte(skipPC >> 8), byte(skipPC), // PUSH2 <skip>
		0x57,       // JUMPI (jump if gas >= 15000)
		0x5f,       // PUSH0 (memory offset)
		0x52,       // MSTORE (store TOS to memory[0])
		0x60, 0x20, // PUSH1 32
		0x5f, // PUSH0
		0xf3, // RETURN
		0x5b, // JUMPDEST <skip>
	}

	if len(g.bytecode)+len(bc) > g.maxSize-32 {
		return
	}

	g.bytecode = append(g.bytecode, bc...)
	// Stack effect: GAS+PUSH2+LT+ISZERO+PUSH2+JUMPI = net 0 (all consumed by JUMPI)
	// On the non-bailout path, stack is unchanged.
	// On the bailout path, MSTORE consumes 1 item (memory offset was pushed by us).
	g.currentGas += 20
	g.opcodeCount += countOpcodesInBytecode(bc)
}

// --- Stack management ---

func (g *CallFuzzGenerator) addStackItems(n int) bool {
	for i := 0; i < n; i++ {
		if len(g.bytecode)+2 > g.maxSize-32 || g.opcodeCount+1 > g.maxOpcodeCount-10 {
			return false
		}
		if g.rng.Float64() < 0.7 {
			// PUSH with random data
			pushSize := g.selectPushSize()
			if len(g.bytecode)+1+pushSize > g.maxSize-32 {
				if len(g.bytecode)+1 > g.maxSize-32 {
					return false
				}
				g.bytecode = append(g.bytecode, 0x5f)
				g.stackSize++
				g.currentGas += 2
				g.opcodeCount++
				continue
			}
			pushBytes := make([]byte, 1+pushSize)
			pushBytes[0] = byte(0x5f + pushSize)
			copy(pushBytes[1:], g.rng.Bytes(pushSize))
			g.bytecode = append(g.bytecode, pushBytes...)
			g.stackSize++
			g.currentGas += 3
			g.opcodeCount++
		} else {
			if len(g.stackBuilders) == 0 {
				g.bytecode = append(g.bytecode, 0x5f)
				g.stackSize++
				g.currentGas += 2
				g.opcodeCount++
				continue
			}
			op := g.stackBuilders[g.rng.Intn(len(g.stackBuilders))]
			if g.currentGas+op.GasCost > g.maxGas {
				g.bytecode = append(g.bytecode, 0x5f)
				g.stackSize++
				g.currentGas += 2
				g.opcodeCount++
				continue
			}
			seq := op.Template()
			if len(g.bytecode)+len(seq) > g.maxSize-32 {
				g.bytecode = append(g.bytecode, 0x5f)
				g.stackSize++
				g.currentGas += 2
				g.opcodeCount++
				continue
			}
			g.bytecode = append(g.bytecode, seq...)
			g.stackSize++
			g.currentGas += op.GasCost
			g.opcodeCount++
		}
	}
	return true
}

func (g *CallFuzzGenerator) selectPushSize() int {
	choice := g.rng.Float64()
	switch {
	case choice < 0.40:
		return 1
	case choice < 0.60:
		return 2
	case choice < 0.75:
		return 3
	case choice < 0.90:
		return 4 + g.rng.Intn(5)
	default:
		return 9 + g.rng.Intn(24)
	}
}

func (g *CallFuzzGenerator) removeStackItems(n int) bool {
	for n > 0 {
		if len(g.bytecode)+1 > g.maxSize-32 || g.opcodeCount+1 > g.maxOpcodeCount-10 {
			return false
		}
		if g.currentGas+2 > g.maxGas {
			return false
		}
		g.bytecode = append(g.bytecode, 0x50) // POP
		g.stackSize--
		g.currentGas += 2
		g.opcodeCount++
		n--
	}
	return true
}

// --- Instruction generation ---

func (g *CallFuzzGenerator) generateNextInstruction() bool {
	choice := g.rng.Float64()

	// Much lower error rates than evm-fuzz for better valid execution coverage
	if choice < 0.0001 { // 0.01% invalid opcode
		if len(g.invalidOpcodes) > 0 {
			g.bytecode = append(g.bytecode, g.invalidOpcodes[g.rng.Intn(len(g.invalidOpcodes))])
			g.currentGas += 3
			g.opcodeCount++
			return true
		}
	}
	if choice < 0.0004 { // 0.03% random byte
		g.bytecode = append(g.bytecode, byte(g.rng.Intn(256)))
		g.currentGas += 3
		g.opcodeCount++
		return true
	}

	return g.generateValidInstruction()
}

func (g *CallFuzzGenerator) generateValidInstruction() bool {
	var candidates []*evmfuzz.OpcodeInfo
	for _, op := range g.validOpcodes {
		if g.currentGas+op.GasCost <= g.maxGas {
			candidates = append(candidates, op)
		}
	}
	if len(candidates) == 0 {
		return false
	}

	op := g.selectWeightedOpcode(candidates)

	// Fulfill stack requirements
	needed := op.StackInput - g.stackSize
	if needed > 0 {
		if !g.addStackItems(needed) {
			return g.generateFallbackInstruction()
		}
	}

	// Prevent stack overflow
	result := g.stackSize - op.StackInput + op.StackOutput
	if result > 1024 {
		toRemove := result - 1024
		if !g.removeStackItems(toRemove) {
			return g.generateFallbackInstruction()
		}
	}

	seq := op.Template()
	seqOps := countOpcodesInBytecode(seq)

	if len(g.bytecode)+len(seq) > g.maxSize || g.opcodeCount+seqOps > g.maxOpcodeCount {
		return g.generateFallbackInstruction()
	}

	g.bytecode = append(g.bytecode, seq...)
	g.currentGas += op.GasCost
	g.opcodeCount += seqOps

	// Update stack
	if g.stackSize >= op.StackInput {
		g.stackSize -= op.StackInput
	} else {
		g.stackSize = 0
	}
	g.stackSize += op.StackOutput

	return true
}

func (g *CallFuzzGenerator) generateFallbackInstruction() bool {
	var candidates []*evmfuzz.OpcodeInfo
	for _, op := range g.validOpcodes {
		if g.stackSize >= op.StackInput &&
			g.stackSize-op.StackInput+op.StackOutput <= 1024 &&
			g.currentGas+op.GasCost <= g.maxGas {
			candidates = append(candidates, op)
		}
	}
	if len(candidates) == 0 {
		if g.stackSize < 1020 {
			g.bytecode = append(g.bytecode, 0x60, byte(g.rng.Intn(256)))
			g.stackSize++
			g.currentGas += 3
			g.opcodeCount++
			return true
		}
		return false
	}

	op := candidates[g.rng.Intn(len(candidates))]
	seq := op.Template()
	seqOps := countOpcodesInBytecode(seq)
	if len(g.bytecode)+len(seq) > g.maxSize || g.opcodeCount+seqOps > g.maxOpcodeCount {
		return false
	}
	g.bytecode = append(g.bytecode, seq...)
	g.currentGas += op.GasCost
	g.opcodeCount += seqOps

	if g.stackSize >= op.StackInput {
		g.stackSize -= op.StackInput
	} else {
		g.stackSize = 0
	}
	g.stackSize += op.StackOutput

	return true
}

func (g *CallFuzzGenerator) selectWeightedOpcode(candidates []*evmfuzz.OpcodeInfo) *evmfuzz.OpcodeInfo {
	totalWeight := 0.0
	for _, op := range candidates {
		totalWeight += op.Probability
	}
	if totalWeight == 0 {
		return candidates[g.rng.Intn(len(candidates))]
	}
	rv := g.rng.Float64() * totalWeight
	cw := 0.0
	for _, op := range candidates {
		cw += op.Probability
		if rv <= cw {
			return op
		}
	}
	return candidates[len(candidates)-1]
}

// --- PUSH, JUMP templates ---

func (g *CallFuzzGenerator) makePushTemplate(opcode uint16, size int) func() []byte {
	return func() []byte {
		result := make([]byte, 1+size)
		result[0] = byte(opcode)
		copy(result[1:], g.rng.Bytes(size))
		return result
	}
}

func (g *CallFuzzGenerator) generateJump() []byte {
	g.jumpPlaceholders = append(g.jumpPlaceholders, len(g.bytecode)+1)
	return []byte{0x61, 0x00, 0x00, 0x56} // PUSH2 0x0000 JUMP
}

func (g *CallFuzzGenerator) generateJumpi() []byte {
	g.jumpPlaceholders = append(g.jumpPlaceholders, len(g.bytecode)+1)
	return []byte{0x61, 0x00, 0x00, 0x57} // PUSH2 0x0000 JUMPI
}

func (g *CallFuzzGenerator) fixJumpTargets() {
	for _, pos := range g.jumpPlaceholders {
		var target int
		// 2% chance of invalid jump (lower than evm-fuzz's 10%)
		if g.rng.Float64() < 0.02 {
			target = g.rng.Intn(max(len(g.bytecode), 1))
		} else if len(g.jumpTargets) > 0 {
			target = g.jumpTargets[g.rng.Intn(len(g.jumpTargets))]
		} else {
			target = g.rng.Intn(max(len(g.bytecode), 1))
		}
		if pos+1 < len(g.bytecode) {
			g.bytecode[pos] = byte(target >> 8)
			g.bytecode[pos+1] = byte(target)
		}
	}
}

func (g *CallFuzzGenerator) pushSeedAndTxID() {
	seedBytes := make([]byte, 32)
	if g.baseSeed != "" {
		if parsed, err := evmfuzz.ParseHexSeed(g.baseSeed); err == nil {
			if len(parsed) >= 32 {
				copy(seedBytes, parsed[:32])
			} else {
				copy(seedBytes[32-len(parsed):], parsed)
			}
		} else {
			s := []byte(g.baseSeed)
			if len(s) >= 32 {
				copy(seedBytes, s[:32])
			} else {
				copy(seedBytes[32-len(s):], s)
			}
		}
	}
	seedPush := make([]byte, 33)
	seedPush[0] = 0x7f // PUSH32
	copy(seedPush[1:], seedBytes)
	g.bytecode = append(g.bytecode, seedPush...)
	g.stackSize++
	g.currentGas += 3

	txIDBytes := make([]byte, 32)
	binary.BigEndian.PutUint64(txIDBytes[24:], g.txID)
	txIDPush := make([]byte, 33)
	txIDPush[0] = 0x7f // PUSH32
	copy(txIDPush[1:], txIDBytes)
	g.bytecode = append(g.bytecode, txIDPush...)
	g.stackSize++
	g.currentGas += 3
}

// countOpcodesInBytecode counts opcodes properly handling PUSH data bytes
func countOpcodesInBytecode(bytecode []byte) int {
	count := 0
	pc := 0
	for pc < len(bytecode) {
		op := bytecode[pc]
		count++
		if op >= 0x60 && op <= 0x7f {
			pc += int(op-0x5f) + 1
		} else {
			pc++
		}
	}
	return count
}
