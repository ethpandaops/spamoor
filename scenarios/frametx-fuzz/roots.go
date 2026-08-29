package frametxfuzz

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"

	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/txtypes"
)

// rootRing tracks the EIP-8272 roots a run has committed, and the clock that maps a block
// to the consensus slot a reference names.
//
// spamoor has no consensus client to ask for the slot, so it is derived from block
// timestamps and then confirmed against the chain: a committed entry's storage key and
// value are both functions of the slot, so reading the storage back proves which slot the
// write landed in. The confirmed offset is cached.
type rootRing struct {
	mutex sync.Mutex

	// genesisTime and slotSeconds map a block timestamp to a slot.
	genesisTime uint64
	slotSeconds uint64

	// calibrated marks that a write has confirmed the mapping against chain storage.
	calibrated bool

	// entries are the roots this run has committed, newest last.
	entries []rootEntry
}

// rootEntry is one committed root, with everything a reference to it needs.
type rootEntry struct {
	SourceAddress common.Address
	Salt          common.Hash
	SourceID      common.Hash
	Slot          uint64
	Root          common.Hash
}

// reference renders the entry as a declared reference.
func (e rootEntry) reference() *txtypes.RecentRootReference {
	return &txtypes.RecentRootReference{SourceID: e.SourceID, Slot: e.Slot, Root: e.Root}
}

// newRootRing establishes the slot clock from the chain's own blocks.
func newRootRing(ctx context.Context, logger logrus.FieldLogger, client *spamoor.Client) (*rootRing, error) {
	genesis, err := client.GetHeaderByNumber(ctx, big0)
	if err != nil {
		return nil, fmt.Errorf("could not read the genesis header: %w", err)
	}

	head, err := client.GetHeaderByNumber(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("could not read the head header: %w", err)
	}

	if head.Number == 0 {
		return nil, fmt.Errorf("the chain has only a genesis block, so no slot can be derived yet")
	}

	// The slot length is the smallest gap between consecutive blocks in a short window:
	// a missed slot widens a gap to a multiple of the slot time but never shortens one.
	slotSeconds, err := deriveSlotSeconds(ctx, client, head.Number)
	if err != nil {
		return nil, err
	}

	logger.Infof("recent root slot clock: genesis %v, %v seconds per slot", genesis.Timestamp, slotSeconds)

	return &rootRing{genesisTime: genesis.Timestamp, slotSeconds: slotSeconds}, nil
}

// deriveSlotSeconds infers the slot length from recent block timestamps.
func deriveSlotSeconds(ctx context.Context, client *spamoor.Client, head uint64) (uint64, error) {
	const window = 16

	first := uint64(1)
	if head > window {
		first = head - window
	}

	previous := uint64(0)
	smallest := uint64(0)

	for number := first; number <= head; number++ {
		header, err := client.GetHeaderByNumber(ctx, bigFromUint64(number))
		if err != nil {
			return 0, fmt.Errorf("could not read header %d: %w", number, err)
		}

		if previous != 0 && header.Timestamp > previous {
			gap := header.Timestamp - previous
			if smallest == 0 || gap < smallest {
				smallest = gap
			}
		}

		previous = header.Timestamp
	}

	if smallest == 0 {
		return 0, fmt.Errorf("could not derive a slot length from block timestamps")
	}

	return smallest, nil
}

// slotFor returns the consensus slot a block timestamp belongs to.
func (r *rootRing) slotFor(timestamp uint64) uint64 {
	if timestamp < r.genesisTime || r.slotSeconds == 0 {
		return 0
	}

	return (timestamp - r.genesisTime) / r.slotSeconds
}

// record confirms which slot a committed root landed in and remembers it. A small scan
// around the timestamp-derived candidate absorbs an off-by-one in the mapping.
func (r *rootRing) record(ctx context.Context, client *spamoor.Client, source common.Address, salt, root common.Hash, blockTime uint64) error {
	sourceID := txtypes.RecentRootSourceID(source, salt)
	candidate := r.slotFor(blockTime)

	const scan = 4

	for offset := -scan; offset <= scan; offset++ {
		slot := uint64(int64(candidate) + int64(offset))
		if int64(candidate)+int64(offset) < 0 {
			continue
		}

		stored, err := client.GetStorageAt(ctx, txtypes.RecentRootAddress, txtypes.RecentRootSlotStorageKey(sourceID, slot))
		if err != nil {
			return err
		}

		if stored != txtypes.RecentRootEntryHash(sourceID, slot, root) {
			continue
		}

		r.mutex.Lock()
		defer r.mutex.Unlock()

		if !r.calibrated && offset != 0 {
			// The timestamp mapping is off by a fixed amount; fold it into the genesis
			// time so every later derivation lands on the first try.
			r.genesisTime -= uint64(int64(offset) * int64(r.slotSeconds))
		}

		r.calibrated = true
		r.entries = append(r.entries, rootEntry{
			SourceAddress: source,
			Salt:          salt,
			SourceID:      sourceID,
			Slot:          slot,
			Root:          root,
		})

		return nil
	}

	return fmt.Errorf("no committed entry found for the root written by %v around slot %d", source, candidate)
}

// usable returns the committed entries a transaction in currentSlot may reference.
func (r *rootRing) usable(currentSlot uint64) []rootEntry {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	usable := make([]rootEntry, 0, len(r.entries))

	for _, entry := range r.entries {
		if txtypes.RecentRootReferenceUsable(currentSlot, entry.Slot) {
			usable = append(usable, entry)
		}
	}

	return usable
}

// calibratedClock reports whether a write has confirmed the slot mapping. Until it has,
// references are not declared: an unconfirmed slot would make every transaction invalid
// for a reason that has nothing to do with the client under test.
func (r *rootRing) calibratedClock() bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	return r.calibrated
}

// references builds the declared references a recipe asks for, reporting whether the
// result is a plain reference or one of the edge cases that is refused by design.
func (r *rootRing) references(recipe *Recipe, currentSlot uint64) ([]*txtypes.RecentRootReference, bool) {
	entries := r.usable(currentSlot)
	if len(entries) == 0 {
		return nil, true
	}

	count := recipe.RecentRoots
	if count > len(entries) {
		count = len(entries)
	}

	references := make([]*txtypes.RecentRootReference, 0, count+1)
	for i := 0; i < count; i++ {
		references = append(references, entries[len(entries)-1-i].reference())
	}

	base := entries[len(entries)-1]

	switch recipe.RecentRootEdge {
	case "same_slot":
		// A root written during slot S becomes referenceable in S+1, so naming the
		// current slot must be refused.
		reference := base.reference()
		reference.Slot = currentSlot

		return append(references, reference), false

	case "future_slot":
		reference := base.reference()
		reference.Slot = currentSlot + 1

		return append(references, reference), false

	case "unwritten":
		reference := base.reference()
		if reference.Slot == 0 {
			return references, true
		}

		reference.Slot--

		return append(references, reference), false

	case "wrong_source":
		reference := base.reference()
		reference.SourceID[0] ^= 0xff

		return append(references, reference), false

	case "outside_window":
		if currentSlot <= txtypes.RecentRootUsableWindow {
			// The chain is too young for the far edge of the window to exist.
			return references, true
		}

		reference := base.reference()
		reference.Slot = currentSlot - txtypes.RecentRootUsableWindow - 1

		return append(references, reference), false

	case "duplicate":
		// Duplicates are valid, and are checked, charged and preserved independently.
		return append(references, base.reference()), true
	}

	return references, true
}

// RootSourceWalletName is the wallet whose address identifies the run's root source. A
// source is keyed by the address that wrote it, so a later run can reference roots an
// earlier one committed.
const RootSourceWalletName = "frametx-fuzz-roots"

// rootWriteGas covers the recent root contract's call and the state gas a new entry
// costs. Most writes are cheaper, since the ring reuses a source's slots.
const rootWriteGas = 60_000 + txtypes.StateBytesPerStorageSet*txtypes.CostPerStateByte

// writeRoot commits one root and records which slot it landed in. The recent root
// contract takes exactly 64 bytes of calldata and zero value, and refuses anything else.
func (s *Scenario) writeRoot(ctx context.Context, client *spamoor.Client, feeCap, tipCap *big.Int) error {
	wallet := s.walletPool.GetWellKnownWallet(RootSourceWalletName)
	if wallet == nil {
		return fmt.Errorf("the %q well-known wallet is not registered", RootSourceWalletName)
	}

	counter := s.rootCounter.Add(1)

	salt := common.BigToHash(new(big.Int).SetUint64(counter))
	root := crypto.Keccak256Hash(wallet.GetAddress().Bytes(), salt.Bytes())

	target := txtypes.RecentRootAddress

	txData, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       rootWriteGas,
		To:        &target,
		Value:     uint256.NewInt(0),
		Data:      txtypes.RecentRootWriteCalldata(salt, root),
	})
	if err != nil {
		return err
	}

	tx, err := wallet.BuildDynamicFeeTx(txData)
	if err != nil {
		return err
	}

	receipt, err := s.walletPool.GetTxPool().SendAndAwaitTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		ClientGroup: s.options.ClientGroup,
		Rebroadcast: true,
	})
	if err != nil {
		wallet.MarkSkippedNonce(tx.Nonce())

		return fmt.Errorf("could not commit a recent root: %w", err)
	}

	if receipt == nil || receipt.Status != txtypes.ReceiptStatusSuccessful {
		return fmt.Errorf("the recent root write did not succeed")
	}

	header, err := client.GetHeaderByNumber(ctx, receipt.BlockNumber)
	if err != nil {
		return fmt.Errorf("could not read the block the root was written in: %w", err)
	}

	return s.env.roots.record(ctx, client, wallet.GetAddress(), salt, root, header.Timestamp)
}

// maintainRoots keeps committed roots available for reference. The first write
// calibrates the slot clock; the trickle afterwards keeps a long run from going stale.
func (s *Scenario) maintainRoots(ctx context.Context) {
	client := s.walletPool.GetClient(spamoor.WithClientGroup(s.options.ClientGroup))
	if client == nil {
		return
	}

	for {
		feeCap, tipCap, err := s.fees(client)
		if err == nil {
			if err := s.writeRoot(ctx, client, feeCap, tipCap); err != nil {
				s.logger.Warnf("recent root write failed: %v", err)
			}
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(rootWriteInterval):
		}
	}
}

// rootWriteInterval is how often a fresh root is committed.
const rootWriteInterval = 2 * time.Minute
