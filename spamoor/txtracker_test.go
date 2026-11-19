package spamoor

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"

	"github.com/ethpandaops/spamoor/txtypes"
)

func laneSeqs(lane *txLane) []uint64 {
	seqs := make([]uint64, 0, len(lane.slots))
	for _, slot := range lane.slots {
		if len(slot.txs) > 0 {
			seqs = append(seqs, slot.seq)
		}
	}

	return seqs
}

func equalSeqs(a, b []uint64) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func TestTxLane(t *testing.T) {
	lane := &txLane{}
	txs := map[uint64]*trackedTx{}

	for _, seq := range []uint64{5, 6, 8, 7, 3} {
		txs[seq] = &trackedTx{seq: seq}
		lane.add(seq, txs[seq])
	}

	replacement := &trackedTx{seq: 6}
	lane.add(6, replacement)

	if got := laneSeqs(lane); !equalSeqs(got, []uint64{3, 5, 6, 7, 8}) {
		t.Fatalf("unexpected lane order: %v", got)
	}

	// an entry removed from the middle leaves its slot until it reaches an end
	lane.remove(7, txs[7])
	lane.remove(7, txs[7])

	if got := laneSeqs(lane); !equalSeqs(got, []uint64{3, 5, 6, 8}) || len(lane.slots) != 5 {
		t.Fatalf("unexpected lane after removal: %v (%d slots)", got, len(lane.slots))
	}

	lane.remove(8, txs[8])

	if len(lane.slots) != 3 {
		t.Fatalf("empty slots were not trimmed from the tail: %d slots", len(lane.slots))
	}

	retired := lane.retire(6, nil)
	if len(retired) != 4 || !lane.empty() {
		t.Fatalf("expected 4 retired transactions and an empty lane, got %d, empty=%v", len(retired), lane.empty())
	}

	if retired := lane.retire(100, nil); len(retired) != 0 {
		t.Fatalf("retired %d transactions from an empty lane", len(retired))
	}
}

func TestTxLanePosition(t *testing.T) {
	keys, seq := txLanePosition(txtypes.NewTx(&txtypes.LegacyTx{Nonce: 7}))
	if keys != nil || seq != 7 {
		t.Fatalf("legacy tx: keys=%v seq=%d", keys, seq)
	}

	frameTx := (&txtypes.FrameTx{}).WithNonceKeys([]*uint256.Int{uint256.NewInt(1), uint256.NewInt(9)}, 4)

	keys, seq = txLanePosition(txtypes.NewTx(frameTx))
	if len(keys) != 2 || keys[0] != *uint256.NewInt(1) || keys[1] != *uint256.NewInt(9) || seq != 4 {
		t.Fatalf("keyed frame tx: keys=%v seq=%d", keys, seq)
	}
}

// Nothing may stay behind in the tracker once its transactions are removed, whatever
// order they leave in.
func TestTxTrackerAddRemove(t *testing.T) {
	tracker := &txTracker{}
	wallets := []*Wallet{
		NewWallet(nil, common.HexToAddress("0x01")),
		NewWallet(nil, common.HexToAddress("0x02")),
	}

	var tracked []*trackedTx

	for _, wallet := range wallets {
		for nonce := uint64(0); nonce < 50; nonce++ {
			tx := newTrackedTx(wallet, txtypes.NewTx(&txtypes.LegacyTx{Nonce: nonce, To: &wallet.address}), &SendTransactionOptions{})
			tracker.add(tx, true, 0)
			tracked = append(tracked, tx)
		}

		keyed := (&txtypes.FrameTx{Sender: wallet.address}).WithNonceKeys([]*uint256.Int{uint256.NewInt(3), uint256.NewInt(4)}, 0)
		tx := newTrackedTx(wallet, txtypes.NewTx(keyed), &SendTransactionOptions{})
		tracker.add(tx, true, 0)
		tracked = append(tracked, tx)
	}

	if txs, senders := tracker.count(); txs != 102 || senders != 2 {
		t.Fatalf("expected 102 transactions of 2 senders, got %d of %d", txs, senders)
	}

	if pending := wallets[0].GetPendingTxs(); len(pending) != 51 || pending[0].Tx.Nonce() != 0 || pending[49].Tx.Nonce() != 49 {
		t.Fatalf("unexpected pending transactions: %d", len(pending))
	}

	if lowest, ok := wallets[0].GetLowestPendingNonce(); !ok || lowest != 0 {
		t.Fatalf("unexpected lowest pending nonce: %d %v", lowest, ok)
	}

	if wallets[0].GetPendingTx(tracked[7].Tx) != &tracked[7].PendingTx {
		t.Fatal("pending transaction not found")
	}

	// middle first, then alternating ends
	for i := len(tracked) / 2; i < len(tracked); i++ {
		tracker.remove(tracked[i], false)
	}

	for i := len(tracked)/2 - 1; i >= 0; i-- {
		tracker.remove(tracked[i], false)
	}

	if txs, senders := tracker.count(); txs != 0 || senders != 0 {
		t.Fatalf("tracker not empty: %d transactions of %d senders", txs, senders)
	}

	if pending := wallets[0].GetPendingTxs(); len(pending) != 0 {
		t.Fatalf("wallet still lists %d pending transactions", len(pending))
	}

	if _, ok := wallets[0].GetLowestPendingNonce(); ok {
		t.Fatal("wallet still reports a pending nonce")
	}
}

// A transaction taken off its lanes stays tracked by its hash only.
func TestTxTrackerUnlane(t *testing.T) {
	tracker := &txTracker{}
	wallet := NewWallet(nil, common.HexToAddress("0x01"))

	first := newTrackedTx(wallet, txtypes.NewTx(&txtypes.LegacyTx{Nonce: 3}), &SendTransactionOptions{})
	second := newTrackedTx(wallet, txtypes.NewTx(&txtypes.LegacyTx{Nonce: 4}), &SendTransactionOptions{})
	tracker.add(first, true, 0)
	tracker.add(second, true, 0)

	tracker.unlane(first)

	if lowest, _ := wallet.GetLowestPendingNonce(); lowest != 4 {
		t.Fatalf("unlaned transaction still holds nonce %d", lowest)
	}

	if len(tracker.byHash(first.hash)) != 1 {
		t.Fatal("unlaned transaction is not tracked by hash")
	}

	tracker.remove(first, false)
	tracker.remove(second, false)

	if txs, senders := tracker.count(); txs != 0 || senders != 0 {
		t.Fatalf("tracker not empty: %d transactions of %d senders", txs, senders)
	}
}

// A transaction that is finished while it is being added must not stay behind in the
// hash index or keep its sender alive.
func TestTxTrackerAddOfFinishedTx(t *testing.T) {
	tracker := &txTracker{}
	wallet := NewWallet(nil, common.HexToAddress("0x01"))

	tracked := newTrackedTx(wallet, txtypes.NewTx(&txtypes.LegacyTx{Nonce: 0}), &SendTransactionOptions{})

	// retired from its lane and removed before add got to the hash index
	tracked.done.Store(true)
	tracker.add(tracked, true, 0)
	tracker.remove(tracked, false)

	if txs, senders := tracker.count(); txs != 0 || senders != 0 {
		t.Fatalf("tracker not empty: %d transactions of %d senders", txs, senders)
	}

	// adding and finishing concurrently never leaves anything behind either
	for i := 0; i < 2000; i++ {
		racing := newTrackedTx(wallet, txtypes.NewTx(&txtypes.LegacyTx{Nonce: uint64(i)}), &SendTransactionOptions{})
		added := make(chan struct{})

		go func() {
			tracker.add(racing, true, 0)
			close(added)
		}()

		// what block processing does once it finds the transaction on its lane
		for {
			sender := tracker.getSender(wallet.GetAddress())
			if sender == nil {
				continue
			}

			sender.mu.Lock()
			retired := sender.retire(nil, uint64(i))
			sender.mu.Unlock()

			if len(retired) > 0 {
				racing.done.Store(true)
				tracker.remove(racing, false)

				break
			}
		}

		<-added

		if txs, senders := tracker.count(); txs != 0 || senders != 0 {
			t.Fatalf("round %d: tracker not empty: %d transactions of %d senders", i, txs, senders)
		}
	}
}
