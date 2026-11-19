package spamoor

import (
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"

	"github.com/ethpandaops/spamoor/txtypes"
)

// trackedTx is one submitted transaction the pool follows until it completes: confirmed,
// replaced by another transaction on one of its nonce lanes, rejected or cancelled.
type trackedTx struct {
	PendingTx

	hash       common.Hash
	wallet     *Wallet
	sender     *txSender
	walletPool *WalletPool

	// seq and keys locate the transaction on its sender's nonce lanes. keys is nil for
	// the account nonce lane, else the EIP-8250 key set.
	seq   uint64
	keys  []uint256.Int
	laned bool

	done    atomic.Bool
	stopCtx atomic.Pointer[func() bool]

	rebroadcasting atomic.Bool
}

// laneSlot holds the transactions competing for one sequence number of a lane.
type laneSlot struct {
	seq uint64
	txs []*trackedTx
}

// txLane is one nonce sequence of a sender, ordered by sequence number. Slots emptied in
// the middle stay in place and are dropped once they reach either end.
type txLane struct {
	slots []laneSlot
}

func (lane *txLane) find(seq uint64) (int, bool) {
	idx := sort.Search(len(lane.slots), func(i int) bool { return lane.slots[i].seq >= seq })

	return idx, idx < len(lane.slots) && lane.slots[idx].seq == seq
}

func (lane *txLane) add(seq uint64, t *trackedTx) {
	if n := len(lane.slots); n == 0 || lane.slots[n-1].seq < seq {
		lane.slots = append(lane.slots, laneSlot{seq: seq, txs: []*trackedTx{t}})

		return
	}

	idx, found := lane.find(seq)
	if found {
		lane.slots[idx].txs = append(lane.slots[idx].txs, t)

		return
	}

	lane.slots = append(lane.slots, laneSlot{})
	copy(lane.slots[idx+1:], lane.slots[idx:])
	lane.slots[idx] = laneSlot{seq: seq, txs: []*trackedTx{t}}
}

func (lane *txLane) remove(seq uint64, t *trackedTx) {
	idx, found := lane.find(seq)
	if !found {
		return
	}

	slot := &lane.slots[idx]

	for i, other := range slot.txs {
		if other == t {
			slot.txs = append(slot.txs[:i], slot.txs[i+1:]...)

			break
		}
	}

	lane.trim()
}

// trim drops empty slots from both ends.
func (lane *txLane) trim() {
	head := 0
	for head < len(lane.slots) && len(lane.slots[head].txs) == 0 {
		head++
	}

	tail := len(lane.slots)
	for tail > head && len(lane.slots[tail-1].txs) == 0 {
		tail--
	}

	if head == tail {
		lane.slots = nil

		return
	}

	if head > 0 || tail < len(lane.slots) {
		clear(lane.slots[:head])
		clear(lane.slots[tail:])
		lane.slots = lane.slots[head:tail]
	}
}

// retire removes and returns every transaction at or below seq.
func (lane *txLane) retire(seq uint64, retired []*trackedTx) []*trackedTx {
	idx := sort.Search(len(lane.slots), func(i int) bool { return lane.slots[i].seq > seq })
	if idx == 0 {
		return retired
	}

	for i := 0; i < idx; i++ {
		retired = append(retired, lane.slots[i].txs...)
		lane.slots[i].txs = nil
	}

	lane.trim()

	return retired
}

func (lane *txLane) empty() bool {
	return len(lane.slots) == 0
}

// txSender is the tracking state of one sender address: its account nonce lane and its
// EIP-8250 keyed lanes. It exists while the sender has tracked transactions.
type txSender struct {
	mu      sync.Mutex
	address common.Address
	wallet  *Wallet

	account txLane
	keyed   map[uint256.Int]*txLane
	pending int
	removed bool

	// lastProgress is the block of the last confirmation or stale check, staleChecks the
	// number of stale checks since the last confirmation.
	lastProgress uint64
	staleChecks  uint64
	suspect      bool
}

// lanes calls fn for each lane a transaction with the given keys occupies.
func (sender *txSender) lanes(keys []uint256.Int, create bool, fn func(lane *txLane)) {
	if keys == nil {
		fn(&sender.account)

		return
	}

	for _, key := range keys {
		lane := sender.keyed[key]
		if lane == nil {
			if !create {
				continue
			}

			if sender.keyed == nil {
				sender.keyed = make(map[uint256.Int]*txLane, len(keys))
			}

			lane = &txLane{}
			sender.keyed[key] = lane
		}

		fn(lane)

		if lane.empty() {
			delete(sender.keyed, key)
		}
	}
}

// retire removes every transaction a transaction included with the given keys and
// sequence number makes impossible, the included one among them. Must hold sender.mu.
func (sender *txSender) retire(keys []uint256.Int, seq uint64) []*trackedTx {
	var retired []*trackedTx

	sender.lanes(keys, false, func(lane *txLane) {
		retired = lane.retire(seq, retired)
	})

	if len(keys) < 2 {
		return retired
	}

	// A transaction on several lanes can be retired from more than one of them.
	seen := make(map[*trackedTx]struct{}, len(retired))
	unique := retired[:0]

	for _, t := range retired {
		if _, ok := seen[t]; !ok {
			seen[t] = struct{}{}
			unique = append(unique, t)
		}
	}

	return unique
}

// pendingTxs returns the tracked transactions in lane order: the account lane by nonce,
// then the keyed lanes.
func (sender *txSender) pendingTxs() []*PendingTx {
	sender.mu.Lock()
	defer sender.mu.Unlock()

	txs := make([]*PendingTx, 0, sender.pending)

	for _, slot := range sender.account.slots {
		for _, t := range slot.txs {
			txs = append(txs, &t.PendingTx)
		}
	}

	if len(sender.keyed) == 0 {
		return txs
	}

	seen := make(map[*trackedTx]struct{}, len(sender.keyed))

	for _, lane := range sender.keyed {
		for _, slot := range lane.slots {
			for _, t := range slot.txs {
				if _, ok := seen[t]; !ok {
					seen[t] = struct{}{}
					txs = append(txs, &t.PendingTx)
				}
			}
		}
	}

	return txs
}

// txLanePosition returns the lanes and sequence number a transaction occupies. keys is
// nil for the account nonce.
func txLanePosition(tx *txtypes.Transaction) ([]uint256.Int, uint64) {
	if tx.UsesAccountNonce() {
		return nil, tx.Nonce()
	}

	frameTx, ok := tx.Inner().(*txtypes.FrameTx)
	if !ok {
		return []uint256.Int{}, tx.Nonce()
	}

	keys := make([]uint256.Int, 0, len(frameTx.NonceKeys))

	for _, key := range frameTx.NonceKeys {
		if key != nil {
			keys = append(keys, *key)
		}
	}

	return keys, frameTx.NonceSeq
}

// txTracker indexes the tracked transactions by hash and by sender.
type txTracker struct {
	hashMutex sync.RWMutex
	hashes    map[common.Hash][]*trackedTx

	sendersMutex sync.RWMutex
	senders      map[common.Address]*txSender
}

func (tracker *txTracker) getSender(address common.Address) *txSender {
	tracker.sendersMutex.RLock()
	sender := tracker.senders[address]
	tracker.sendersMutex.RUnlock()

	return sender
}

func (tracker *txTracker) allSenders() []*txSender {
	tracker.sendersMutex.RLock()
	defer tracker.sendersMutex.RUnlock()

	senders := make([]*txSender, 0, len(tracker.senders))
	for _, sender := range tracker.senders {
		senders = append(senders, sender)
	}

	return senders
}

func (tracker *txTracker) byHash(hash common.Hash) []*trackedTx {
	tracker.hashMutex.RLock()
	defer tracker.hashMutex.RUnlock()

	return tracker.hashes[hash]
}

// add starts tracking a transaction. With lane set it also joins its sender's nonce
// lanes; a transaction whose sequence number is already used does not.
func (tracker *txTracker) add(t *trackedTx, lane bool, block uint64) {
	address := t.wallet.GetAddress()

	for {
		tracker.sendersMutex.Lock()
		if tracker.senders == nil {
			tracker.senders = make(map[common.Address]*txSender, 64)
		}

		sender := tracker.senders[address]
		if sender == nil {
			sender = &txSender{address: address, wallet: t.wallet, lastProgress: block}
			tracker.senders[address] = sender
		}
		tracker.sendersMutex.Unlock()

		sender.mu.Lock()
		if sender.removed {
			sender.mu.Unlock()

			continue
		}

		t.sender = sender
		t.laned = lane
		sender.pending++

		if lane {
			sender.lanes(t.keys, true, func(l *txLane) { l.add(t.seq, t) })
		}
		sender.mu.Unlock()

		t.wallet.sender.Store(sender)

		break
	}

	tracker.hashMutex.Lock()
	if tracker.hashes == nil {
		tracker.hashes = make(map[common.Hash][]*trackedTx, 1024)
	}

	tracker.hashes[t.hash] = append(tracker.hashes[t.hash], t)
	tracker.hashMutex.Unlock()
}

// remove stops tracking a transaction and drops its sender once nothing is left on it.
func (tracker *txTracker) remove(t *trackedTx) {
	tracker.hashMutex.Lock()
	txs := tracker.hashes[t.hash]

	for i, other := range txs {
		if other != t {
			continue
		}

		if len(txs) == 1 {
			delete(tracker.hashes, t.hash)
		} else {
			tracker.hashes[t.hash] = append(txs[:i:i], txs[i+1:]...)
		}

		break
	}
	tracker.hashMutex.Unlock()

	sender := t.sender
	if sender == nil {
		return
	}

	sender.mu.Lock()
	if t.laned {
		sender.lanes(t.keys, false, func(l *txLane) { l.remove(t.seq, t) })
	}

	sender.pending--
	drop := sender.pending == 0 && !sender.removed
	sender.mu.Unlock()

	if !drop {
		return
	}

	tracker.sendersMutex.Lock()
	sender.mu.Lock()
	if sender.pending == 0 && !sender.removed {
		sender.removed = true
		delete(tracker.senders, sender.address)
	}
	sender.mu.Unlock()
	tracker.sendersMutex.Unlock()
}

// count returns the number of tracked transactions and senders.
func (tracker *txTracker) count() (txs int, senders int) {
	tracker.hashMutex.RLock()
	for _, list := range tracker.hashes {
		txs += len(list)
	}
	tracker.hashMutex.RUnlock()

	tracker.sendersMutex.RLock()
	senders = len(tracker.senders)
	tracker.sendersMutex.RUnlock()

	return txs, senders
}

func newTrackedTx(wallet *Wallet, tx *txtypes.Transaction, options *SendTransactionOptions) *trackedTx {
	t := &trackedTx{
		PendingTx: PendingTx{
			Tx:        tx,
			Submitted: time.Now(),
			Options:   options,
		},
		hash:   tx.Hash(),
		wallet: wallet,
	}
	t.keys, t.seq = txLanePosition(tx)

	return t
}
