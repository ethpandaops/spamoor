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

// nonceLedger tracks the EIP-8250 keyed nonce sequences a run has consumed.
//
// A keyed nonce lives in protocol storage under NONCE_MANAGER rather than in the account
// nonce, so nothing in the engine tracks it. Two properties of the EIP shape this:
//
//   - every key a transaction selects must currently sit at the same nonce_seq, and a
//     successful payment approval moves all of them to nonce_seq + 1 together, so keys
//     can only be combined while their sequences agree;
//   - a key's first use costs KEYED_NONCE_FIRST_USE_GAS out of the frame executing
//     APPROVE, and the whole validation prefix must stay under MaxVerifyGas, so at most
//     four keys can see their first use in one publicly propagated transaction.
//
// The ledger therefore groups keys by sequence and hands out subsets that share one, and
// the caller budgets the first-use surcharge from the ledger's own view of which keys are
// new.
type nonceLedger struct {
	mutex sync.Mutex

	// sequences maps a sender's key to the sequence its next transaction must carry.
	sequences map[common.Address]map[uint64]uint64

	// known records which keys have been read from the chain, so an unread key is
	// looked up once rather than assumed to be at zero.
	known map[common.Address]map[uint64]bool
}

// maxFirstUseKeys is how many never-before-used keys fit in a mempool-legal transaction.
//
// MaxVerifyGas covers the whole validation prefix including signature verification, and
// each first use costs KeyedNonceFirstUseGas out of the approving frame.
const maxFirstUseKeys = (txtypes.MaxVerifyGas - 10_000) / txtypes.KeyedNonceFirstUseGas

// newNonceLedger returns an empty ledger.
func newNonceLedger() *nonceLedger {
	return &nonceLedger{
		sequences: map[common.Address]map[uint64]uint64{},
		known:     map[common.Address]map[uint64]bool{},
	}
}

// nonceKey derives the key a sender uses for one of its key slots.
//
// Keys are derived rather than random so that a rerun addresses the same domains, and
// they are salted with the sender so that two wallets never share one: overlapping key
// sets across senders would make the sequences depend on each other's ordering.
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

// selectKeys returns up to count keys for a sender that all sit at the same sequence.
//
// Reading the current sequence from the chain rather than assuming zero is what lets a
// second run against the same chain work at all: the protocol never writes zero, so an
// absent slot and a first use are the same observation.
func (l *nonceLedger) selectKeys(ctx context.Context, client *spamoor.Client, sender common.Address, count int) (*selection, error) {
	if count > txtypes.MaxNonceKeys {
		count = txtypes.MaxNonceKeys
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	sequences := l.sequences[sender]
	if sequences == nil {
		sequences = map[uint64]uint64{}
		l.sequences[sender] = sequences
		l.known[sender] = map[uint64]bool{}
	}

	known := l.known[sender]

	// Group the candidate slots by the sequence they are at, reading any slot this run
	// has not seen before.
	bySequence := map[uint64][]int{}

	for slot := 0; slot < count*2 && len(bySequence) <= count; slot++ {
		if !known[uint64(slot)] {
			key := nonceKey(sender, slot)

			value, err := client.GetStorageAt(ctx, txtypes.NonceManager, txtypes.NonceManagerSlot(sender, key))
			if err != nil {
				return nil, err
			}

			sequences[uint64(slot)] = new(uint256.Int).SetBytes(value.Bytes()).Uint64()
			known[uint64(slot)] = true
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

		if best == 0 {
			result.firstUses++
		}
	}

	// Keys must be strictly increasing by value, which the slot order does not
	// guarantee once they are hashed into the key space.
	sortKeys(result.keys)

	if result.firstUses > maxFirstUseKeys {
		// Trimming keeps the transaction inside the public mempool's verification gas
		// cap. Generating one that exceeds it is the negative tier's job, not an
		// accident of how many fresh keys a wallet happened to have.
		result.keys = result.keys[:maxFirstUseKeys]
		result.firstUses = maxFirstUseKeys
	}

	result.slots = bestSlots[:len(result.keys)]

	return result, nil
}

// consumed records that a selection landed, moving every key in it to the next sequence.
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
	}
}

// forget drops a sender's cached sequences so the next selection reads them again. It is
// used when a transaction is rejected for a reason that may mean the ledger is stale.
func (l *nonceLedger) forget(sender common.Address) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	delete(l.sequences, sender)
	delete(l.known, sender)
}

// sortKeys orders keys by numeric value, which is what the payload requires.
func sortKeys(keys []*uint256.Int) {
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j].Cmp(keys[j-1]) < 0; j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
		}
	}
}
