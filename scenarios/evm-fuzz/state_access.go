package evmfuzz

import (
	"crypto/sha256"
	"encoding/binary"
)

// slotModel mirrors the virtual-stack model for storage: tracks the value the
// generator believes is committed (original) and pending (current) for a slot.
type slotModel struct {
	original uint64 // value at start of this tx (modeled, low 64 bits)
	current  uint64 // value after the last modeled SSTORE
}

// initStatePool derives the shared (seed-only) and private (seed,txID) slot/address
// pools deterministically. Shared keys are identical across txIDs so contracts in
// one block collide on storage -> exercises EIP-7928 BAL merge & cross-tx RAW.
// Pools are derived via sha256 (NOT g.rng) so they never perturb the rng stream
// used by the other fuzz modes.
func (g *OpcodeGenerator) initStatePool() {
	g.slotState = make(map[uint64]slotModel)

	seedBytes, _ := parseHexSeed(g.baseSeed)

	sh := sha256.Sum256(append([]byte("spamoor/state-access/shared"), seedBytes...))
	for i := 0; i < 4; i++ {
		g.sharedSlots[i] = binary.BigEndian.Uint64(sh[i*8:]) | 1 // nonzero key
	}
	for i := 0; i < 3; i++ {
		copy(g.sharedAddrs[i][:], sh[(i*8)%12:])
		copy(g.sharedAddrs[i][8:], sh[16+i*4:])
	}

	ph := sha256.Sum256(binary.BigEndian.AppendUint64(
		append([]byte("spamoor/state-access/private"), seedBytes...), g.txID))
	for i := 0; i < 4; i++ {
		g.privateSlots[i] = binary.BigEndian.Uint64(ph[i*8:]) | 1
	}
}

// pushU64 emits PUSH8 <v> (1+8 bytes); one opcode, +1 stack item.
func pushU64(v uint64) []byte {
	b := make([]byte, 9)
	b[0] = 0x67 // PUSH8
	binary.BigEndian.PutUint64(b[1:], v)
	return b
}

// pushAddr emits PUSH20 <addr> (1+20 bytes); one opcode, +1 stack item.
func pushAddr(a [20]byte) []byte {
	b := make([]byte, 21)
	b[0] = 0x73 // PUSH20
	copy(b[1:], a[:])
	return b
}

func (g *OpcodeGenerator) pickSharedSlot() uint64   { return g.sharedSlots[g.rng.Intn(4)] }
func (g *OpcodeGenerator) pickPrivateSlot() uint64  { return g.privateSlots[g.rng.Intn(4)] }
func (g *OpcodeGenerator) pickSharedAddr() [20]byte { return g.sharedAddrs[g.rng.Intn(3)] }

// modelSlot returns the current model for key, initializing original=current=0
// the first time the slot is touched in this tx (matches a cold, never-written slot).
func (g *OpcodeGenerator) modelSlot(key uint64) slotModel {
	m, ok := g.slotState[key]
	if !ok {
		m = slotModel{original: 0, current: 0}
		g.slotState[key] = m
	}
	return m
}

func (g *OpcodeGenerator) setSlotCurrent(key, v uint64) {
	m := g.modelSlot(key)
	m.current = v
	g.slotState[key] = m
}

// isStateAccessOpcode reports whether opcode is one of the virtual state-access
// generator selectors (0x120/0x121). These are never emitted as real bytes.
func (g *OpcodeGenerator) isStateAccessOpcode(opcode uint16) bool {
	return opcode == 0x120 || opcode == 0x121
}

// generateStorageRefundSeq drives SSTORE through a chosen (original,current,new)
// transition from the refund truth table (EIP-8037 source-based refunds,
// EIP-7778 no block-level refunds). Net stack effect: 0.
func (g *OpcodeGenerator) generateStorageRefundSeq() []byte {
	var bc []byte
	useShared := g.rng.Float64() < 0.6 // bias shared -> cross-contract collisions
	key := g.pickPrivateSlot()
	if useShared {
		key = g.pickSharedSlot()
	}
	m := g.modelSlot(key)

	// sstore emits: PUSH8 val ; PUSH8 key ; SSTORE   (consumes both -> net 0)
	sstore := func(val uint64) {
		bc = append(bc, pushU64(val)...) // value
		bc = append(bc, pushU64(key)...) // slot key
		bc = append(bc, 0x55)            // SSTORE
	}

	switch g.rng.Intn(5) {
	case 0: // 0 -> x  (fresh set, no refund)
		x := uint64(1) + g.rng.Uint64()
		sstore(x)
		g.setSlotCurrent(key, x)
	case 1: // x -> 0  (clear refund). Ensure current!=0 first.
		if m.current == 0 {
			sstore(7) // establish nonzero
		}
		sstore(0) // clear -> refund
		g.setSlotCurrent(key, 0)
	case 2: // x -> y  (dirty update, both nonzero)
		if m.current == 0 {
			sstore(3)
		}
		sstore(9)
		g.setSlotCurrent(key, 9)
	case 3: // x -> y -> x  (reset to original within same tx: refund rollback path)
		orig := m.original
		sstore(orig + 100) // dirty
		sstore(orig)       // reset back to original
		g.setSlotCurrent(key, orig)
	case 4: // set-then-REVERT in a child CREATE frame (refund must be rolled back)
		bc = append(bc, g.generateRevertingSstore(key)...)
		// model unchanged: revert discards the write
	}
	return bc
}

// generateRevertingSstore deploys tiny init-code via CREATE that does
// PUSH key; SSTORE; PUSH0 PUSH0 REVERT. The CREATE fails (revert in init),
// so the storage write to the *created* address is rolled back -> tests refund
// rollback (EIP-8037) without aborting this contract's deployment. Net stack: 0.
func (g *OpcodeGenerator) generateRevertingSstore(key uint64) []byte {
	// init code: PUSH8 v(0x01) ; PUSH8 key ; SSTORE ; PUSH0 ; PUSH0 ; REVERT
	var init []byte
	init = append(init, pushU64(1)...)
	init = append(init, pushU64(key)...)
	init = append(init, 0x55)             // SSTORE (in new account's storage)
	init = append(init, 0x5f, 0x5f, 0xfd) // PUSH0 PUSH0 REVERT

	// store init code left-aligned in one 32-byte word at mem[0]
	word := make([]byte, 32)
	copy(word, init)

	var bc []byte
	bc = append(bc, 0x7f)                  // PUSH32
	bc = append(bc, word...)               // init-code word
	bc = append(bc, 0x5f, 0x52)            // PUSH0 ; MSTORE  (mem[0]=word)
	bc = append(bc, 0x60, byte(len(init))) // PUSH1 size
	bc = append(bc, 0x5f)                  // PUSH0  offset
	bc = append(bc, 0x5f)                  // PUSH0  value
	bc = append(bc, 0xf0)                  // CREATE -> 0 (reverted); leaves 1 word
	bc = append(bc, 0x50)                  // POP the create result -> net 0
	return bc
}

// generateAccessWalk touches one pool slot/address twice to exercise cold->warm
// repricing (EIP-8038) and access-list pre-warming (EIP-7981/7928).
// Net stack effect: 0.
func (g *OpcodeGenerator) generateAccessWalk() []byte {
	var bc []byte
	sload := func(key uint64) { // PUSH8 key ; SLOAD ; POP   (net 0)
		bc = append(bc, pushU64(key)...)
		bc = append(bc, 0x54, 0x50) // SLOAD ; POP
	}
	sstore := func(key, v uint64) {
		bc = append(bc, pushU64(v)...)
		bc = append(bc, pushU64(key)...)
		bc = append(bc, 0x55) // SSTORE
	}

	switch g.rng.Intn(4) {
	case 0: // SLOAD then SLOAD: cold -> warm on the SAME slot
		key := g.pickSharedSlot()
		sload(key)
		sload(key)
	case 1: // SLOAD then SSTORE: warm read feeds dirty write pricing
		key := g.pickPrivateSlot()
		sload(key)
		sstore(key, 1)
		g.setSlotCurrent(key, 1)
	case 2: // account-access opcodes on the SAME address repeatedly (cold->warm)
		addr := g.pickSharedAddr()
		bc = append(bc, pushAddr(addr)...)
		bc = append(bc, 0x3b, 0x50) // EXTCODESIZE ; POP   (cold)
		bc = append(bc, pushAddr(addr)...)
		bc = append(bc, 0x31, 0x50) // BALANCE ; POP        (warm)
		bc = append(bc, pushAddr(addr)...)
		bc = append(bc, 0x3f, 0x50) // EXTCODEHASH ; POP
	case 3: // EXTCODECOPY on the SAME address twice (cold->warm, mem-bounded)
		addr := g.pickSharedAddr()
		copyOnce := func() { // EXTCODECOPY(addr, destOff=0, off=0, size=32)
			bc = append(bc, 0x60, 0x20)        // PUSH1 32   size
			bc = append(bc, 0x5f)              // PUSH0      codeOffset
			bc = append(bc, 0x5f)              // PUSH0      destOffset
			bc = append(bc, pushAddr(addr)...) // address
			bc = append(bc, 0x3c)              // EXTCODECOPY (consumes 4 -> net 0)
		}
		copyOnce()
		copyOnce()
	}
	return bc
}
