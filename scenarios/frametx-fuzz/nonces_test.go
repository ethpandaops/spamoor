package frametxfuzz

import (
	"testing"

	"github.com/holiman/uint256"
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
