package frametxfuzz

import (
	"context"
	"encoding/binary"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"

	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txtypes"
)

// nonceLedger tracks the EIP-8250 keyed nonce sequences a run has consumed. They live in
// protocol storage under NONCE_MANAGER, so nothing in the engine tracks them.
//
// Two rules shape it: every key a transaction selects must sit at the same nonce_seq, and
// a key's first use costs KeyedNonceFirstUseGas out of the frame executing APPROVE. The
// ledger therefore hands out subsets that share a sequence, and reports how many of them
// are new so the caller can budget the surcharge.
type nonceLedger struct {
	mutex sync.Mutex

	// sequences maps a sender's key to the sequence its next transaction must carry.
	sequences map[common.Address]map[uint64]uint64

	// known records which keys have been read from the chain, so an unread key is
	// looked up once rather than assumed to be at zero.
	known map[common.Address]map[uint64]bool

	// reserved records the slots an in-flight transaction is using. The chain only
	// advances a sequence once the transaction lands, so without this two concurrent
	// draws for one sender would pick the same keys at the same sequence and the
	// second would be refused for a nonce the run itself burned.
	reserved map[common.Address]map[uint64]bool
}

// newNonceLedger returns an empty ledger.
func newNonceLedger() *nonceLedger {
	return &nonceLedger{
		sequences: map[common.Address]map[uint64]uint64{},
		known:     map[common.Address]map[uint64]bool{},
		reserved:  map[common.Address]map[uint64]bool{},
	}
}

// nonceKey derives the key a sender uses for one of its slots. Salting with the sender
// keeps two wallets from sharing a key, whose sequences would then depend on each other's
// ordering.
func nonceKey(sender common.Address, slot int) *uint256.Int {
	var buf [32]byte

	copy(buf[:20], sender.Bytes())
	binary.BigEndian.PutUint64(buf[24:], uint64(slot)+1)

	return new(uint256.Int).SetBytes(buf[:])
}

// selection is a key set that shares a sequence, together with what using it costs.
type selection struct {
	keys      []*uint256.Int
	slots     []int
	sequence  uint64
	firstUses int
}

// selectKeys returns up to count keys for a sender that all sit at the same sequence,
// with no more than maxFirstUses of them seeing their first use. The slots it returns
// are reserved until the caller reports them consumed or releases them.
//
// Reading the current sequence from the chain rather than assuming zero is what lets a
// second run against the same chain work at all: the protocol never writes zero, so an
// absent slot and a first use are the same observation.
func (l *nonceLedger) selectKeys(ctx context.Context, client *spamoor.Client, sender common.Address, count, maxFirstUses int) (*selection, error) {
	if count > txtypes.MaxNonceKeys {
		count = txtypes.MaxNonceKeys
	}

	// The chain reads happen outside the lock, so a slow node holds up this draw and
	// not every sender's. What they return is merged afterwards, and a slot that was
	// consumed in the meantime keeps the newer value the ledger already holds.
	unread := l.unreadSlots(sender, count*2)
	if len(unread) > 0 {
		values := make(map[int]uint64, len(unread))

		for _, slot := range unread {
			key := nonceKey(sender, slot)

			value, err := client.GetStorageAt(ctx, txtypes.NonceManager, txtypes.NonceManagerSlot(sender, key))
			if err != nil {
				return nil, err
			}

			values[slot] = new(uint256.Int).SetBytes(value.Bytes()).Uint64()
		}

		l.learn(sender, values)
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	sequences := l.sequences[sender]
	known := l.known[sender]
	reserved := l.reserved[sender]

	// Group the candidate slots by the sequence they are at. A slot another transaction
	// is using is not a candidate, and neither is one released since it was read, which
	// the next draw will read again.
	bySequence := map[uint64][]int{}

	for slot := 0; slot < count*2; slot++ {
		if reserved[uint64(slot)] || !known[uint64(slot)] {
			continue
		}

		sequence := sequences[uint64(slot)]
		bySequence[sequence] = append(bySequence[sequence], slot)
	}

	// Prefer the largest group, so that a run converges on using many keys together
	// rather than one at a time.
	best := uint64(0)
	bestSlots := []int(nil)

	for sequence, slots := range bySequence {
		if len(slots) > len(bestSlots) || (len(slots) == len(bestSlots) && sequence < best) {
			best, bestSlots = sequence, slots
		}
	}

	if len(bestSlots) == 0 {
		return nil, nil
	}

	if len(bestSlots) > count {
		bestSlots = bestSlots[:count]
	}

	result := &selection{sequence: best}

	for _, slot := range bestSlots {
		result.keys = append(result.keys, nonceKey(sender, slot))
		result.slots = append(result.slots, slot)

		if best == 0 {
			result.firstUses++
		}
	}

	// Keys must be strictly increasing by value, which the slot order does not
	// guarantee once they are hashed into the key space. The slots move with their
	// keys so that consumed() advances exactly the slots the transaction used.
	sortKeyedSlots(result.keys, result.slots)

	if result.firstUses > maxFirstUses {
		// Trimming keeps the transaction inside the public mempool's verification gas
		// cap. Generating one that exceeds it is a job for a deliberate violation, not
		// an accident of how many fresh keys a wallet happened to have.
		if maxFirstUses <= 0 {
			return nil, nil
		}

		result.keys = result.keys[:maxFirstUses]
		result.slots = result.slots[:maxFirstUses]
		result.firstUses = maxFirstUses
	}

	for _, slot := range result.slots {
		reserved[uint64(slot)] = true
	}

	return result, nil
}

// unreadSlots lists the slots below limit this run has not read from the chain yet,
// creating the sender's ledger entries on first sight.
func (l *nonceLedger) unreadSlots(sender common.Address, limit int) []int {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.sequences[sender] == nil {
		l.sequences[sender] = map[uint64]uint64{}
		l.known[sender] = map[uint64]bool{}
		l.reserved[sender] = map[uint64]bool{}
	}

	known := l.known[sender]
	unread := []int(nil)

	for slot := 0; slot < limit; slot++ {
		if !known[uint64(slot)] {
			unread = append(unread, slot)
		}
	}

	return unread
}

// learn records sequences read from the chain for slots the ledger did not know. A slot
// that became known while the read was in flight was consumed by a transaction that
// landed after the read, so the ledger's value is the newer one and is kept.
func (l *nonceLedger) learn(sender common.Address, values map[int]uint64) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	sequences := l.sequences[sender]
	known := l.known[sender]

	for slot, sequence := range values {
		if known[uint64(slot)] {
			continue
		}

		sequences[uint64(slot)] = sequence
		known[uint64(slot)] = true
	}
}

// consumed records that a selection landed, moving every key in it to the next sequence
// and freeing its slots for the next draw.
func (l *nonceLedger) consumed(sender common.Address, sel *selection) {
	if sel == nil {
		return
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	sequences := l.sequences[sender]
	if sequences == nil {
		return
	}

	for _, slot := range sel.slots {
		sequences[uint64(slot)] = sel.sequence + 1
		l.known[sender][uint64(slot)] = true
		delete(l.reserved[sender], uint64(slot))
	}
}

// release frees a selection's slots without advancing them, for a transaction that was
// not submitted or was refused. The slots are also forgotten so the next draw reads
// them from the chain again: a refusal may mean the ledger's view of them is stale, and
// re-reading is what brings it back in step.
func (l *nonceLedger) release(sender common.Address, sel *selection) {
	if sel == nil {
		return
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.sequences[sender] == nil {
		return
	}

	for _, slot := range sel.slots {
		delete(l.reserved[sender], uint64(slot))
		delete(l.known[sender], uint64(slot))
	}
}

// sortKeyedSlots orders keys by numeric value, which is what the payload requires,
// carrying each key's slot along so the two stay paired.
func sortKeyedSlots(keys []*uint256.Int, slots []int) {
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j].Cmp(keys[j-1]) < 0; j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
			slots[j], slots[j-1] = slots[j-1], slots[j]
		}
	}
}
