package frametxfuzz

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"

	"github.com/ethpandaops/spamoor/txtypes"
)

// The payload wants keys in increasing order, but the ledger tracks them by slot.
// Sorting must move the slots with their keys, or a trimmed selection would advance
// the wrong slots once it lands.
func TestSortKeyedSlotsKeepsPairs(t *testing.T) {
	keys := []*uint256.Int{
		uint256.NewInt(30),
		uint256.NewInt(10),
		uint256.NewInt(40),
		uint256.NewInt(20),
	}
	slots := []int{0, 1, 2, 3}

	sortKeyedSlots(keys, slots)

	wantKeys := []uint64{10, 20, 30, 40}
	wantSlots := []int{1, 3, 0, 2}

	for i := range keys {
		if keys[i].Uint64() != wantKeys[i] {
			t.Fatalf("key %d = %d, want %d", i, keys[i].Uint64(), wantKeys[i])
		}

		if slots[i] != wantSlots[i] {
			t.Fatalf("slot %d = %d, want %d", i, slots[i], wantSlots[i])
		}
	}
}

// Two draws for one sender must not hand out the same slots while the first is in
// flight, and a released or consumed selection must make its slots available again.
func TestSelectKeysReservesSlotsUntilSettled(t *testing.T) {
	sender := common.HexToAddress("0x1234")
	ledger := newNonceLedger()

	// Pre-seed every slot the draw can look at so no chain read is needed.
	ledger.sequences[sender] = map[uint64]uint64{}
	ledger.known[sender] = map[uint64]bool{}
	ledger.reserved[sender] = map[uint64]bool{}

	for slot := uint64(0); slot < 4; slot++ {
		ledger.sequences[sender][slot] = 3
		ledger.known[sender][slot] = true
	}

	first, err := ledger.selectKeys(context.Background(), nil, sender, 2, 2)
	if err != nil || first == nil || len(first.slots) != 2 {
		t.Fatalf("first draw = %v, %v", first, err)
	}

	second, err := ledger.selectKeys(context.Background(), nil, sender, 2, 2)
	if err != nil || second == nil || len(second.slots) != 2 {
		t.Fatalf("second draw = %v, %v", second, err)
	}

	for _, a := range first.slots {
		for _, b := range second.slots {
			if a == b {
				t.Fatalf("slot %d handed out twice while in flight", a)
			}
		}
	}

	if third, err := ledger.selectKeys(context.Background(), nil, sender, 2, 2); err != nil || third != nil {
		t.Fatalf("third draw with every slot reserved = %v, %v, want nil", third, err)
	}

	ledger.consumed(sender, first)

	for _, slot := range first.slots {
		if got := ledger.sequences[sender][uint64(slot)]; got != 4 {
			t.Fatalf("consumed slot %d at sequence %d, want 4", slot, got)
		}
	}

	ledger.release(sender, second)

	for _, slot := range second.slots {
		if ledger.reserved[sender][uint64(slot)] {
			t.Fatalf("released slot %d still reserved", slot)
		}

		if ledger.known[sender][uint64(slot)] {
			t.Fatalf("released slot %d still cached, want a fresh chain read", slot)
		}
	}

	// Stand in for the chain read the released slots now need, then check that every
	// slot is free again and the consumed ones moved on.
	for _, slot := range second.slots {
		ledger.sequences[sender][uint64(slot)] = 4
		ledger.known[sender][uint64(slot)] = true
	}

	again, err := ledger.selectKeys(context.Background(), nil, sender, 2, 2)
	if err != nil || again == nil || len(again.slots) != 2 || again.sequence != 4 {
		t.Fatalf("draw after settling = %v, %v, want 2 slots at sequence 4", again, err)
	}
}

// A slot read at MaxNonceSeq is an exhausted key: a sequential use would have to carry
// the single-use sequence, which needs an unused key, so the ledger must never offer it.
func TestSelectKeysSkipsExhaustedSlots(t *testing.T) {
	sender := common.HexToAddress("0x1234")
	ledger := newNonceLedger()

	ledger.sequences[sender] = map[uint64]uint64{}
	ledger.known[sender] = map[uint64]bool{}
	ledger.reserved[sender] = map[uint64]bool{}

	for slot := uint64(0); slot < 4; slot++ {
		ledger.sequences[sender][slot] = txtypes.MaxNonceSeq
		ledger.known[sender][slot] = true
	}

	ledger.sequences[sender][1] = 2

	sel, err := ledger.selectKeys(context.Background(), nil, sender, 2, 2)
	if err != nil || sel == nil {
		t.Fatalf("draw = %v, %v", sel, err)
	}

	if len(sel.slots) != 1 || sel.slots[0] != 1 || sel.sequence != 2 {
		t.Fatalf("draw picked slots %v at sequence %d, want only slot 1 at 2", sel.slots, sel.sequence)
	}
}

// A single-use selection lives outside the slot window, so settling or releasing it must
// leave the sequential ledger alone, and its keys must be distinct from sequential keys.
func TestSingleUseSelectionLeavesSlotsAlone(t *testing.T) {
	sender := common.HexToAddress("0x1234")
	ledger := newNonceLedger()

	ledger.sequences[sender] = map[uint64]uint64{0: 3}
	ledger.known[sender] = map[uint64]bool{0: true}
	ledger.reserved[sender] = map[uint64]bool{}

	sel := &selection{
		keys:      []*uint256.Int{singleUseKey(sender, 42, 0)},
		slots:     []int{-1},
		sequence:  txtypes.MaxNonceSeq,
		firstUses: 1,
	}

	ledger.consumed(sender, sel)
	ledger.release(sender, sel)

	if got := ledger.sequences[sender][0]; got != 3 {
		t.Fatalf("sequential slot 0 moved to %d", got)
	}

	if !ledger.known[sender][0] {
		t.Fatal("sequential slot 0 was forgotten")
	}

	for slot := 0; slot < 8; slot++ {
		if singleUseKey(sender, 42, 0).Eq(nonceKey(sender, slot)) {
			t.Fatalf("single-use key collides with sequential slot %d", slot)
		}
	}

	if singleUseKey(sender, 42, 0).Eq(singleUseKey(sender, 42, 1)) || singleUseKey(sender, 42, 0).Eq(singleUseKey(sender, 43, 0)) {
		t.Fatal("single-use keys are not distinct per ordinal and index")
	}
}
