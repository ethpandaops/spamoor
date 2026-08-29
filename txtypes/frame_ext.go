package txtypes

import (
	"encoding/binary"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
)

// Protocol state derivations for the two EIPs that extend the frame transaction
// envelope. Reading the resulting slots is the caller's business.

// EIP-8250 keyed nonce constants.
const (
	// KeyedNonceFirstUseGas is charged per never-before-used nonce key, deducted from
	// the frame executing the payment-scoped APPROVE and counted in its gas_used. With
	// the whole validation prefix bounded by MaxVerifyGas, it is what limits how many
	// fresh keys fit in a mempool-legal transaction.
	KeyedNonceFirstUseGas = 20_000

	// MaxNonceSeq is EIP-8250's exclusive bound on nonce_seq.
	MaxNonceSeq = uint64(1<<64 - 1)
)

// NonceManagerCode is the runtime code installed at NonceManager on activation, a
// bare revert(0, 0). Ordinary calls to the address are meant to fail; only the protocol
// writes its storage.
var NonceManagerCode = []byte{0x60, 0x00, 0x60, 0x00, 0xfd}

// EIP-8272 recent root constants.
const (
	// RecentRootLength is the number of entries a root source can hold, so a source
	// overwrites itself once every RecentRootLength slots.
	RecentRootLength = 8192

	// RecentRootUsableWindow is the largest current_slot - slot a reference may have.
	RecentRootUsableWindow = 8191

	// RecentRootWriteLength is the exact calldata length the recent root contract
	// accepts: a 32-byte salt followed by a 32-byte root. Anything else reverts.
	RecentRootWriteLength = 64
)

var (
	// RecentRootEntryDomain separates committed entry hashes from storage keys.
	RecentRootEntryDomain = crypto.Keccak256Hash([]byte("RECENT_ROOT_ENTRY"))

	// RecentRootStorageDomain separates storage keys from committed entry hashes.
	RecentRootStorageDomain = crypto.Keccak256Hash([]byte("RECENT_ROOT_STORAGE"))
)

// NonceManagerSlot returns the NonceManager storage slot holding the sequence of one
// keyed nonce domain, keccak256(pad32(sender) || pad32(key)).
//
// The slot is only meaningful for a non-zero key: key zero aliases the sender's
// ordinary account nonce and is never stored here.
func NonceManagerSlot(sender common.Address, key *uint256.Int) common.Hash {
	var buf [64]byte

	copy(buf[12:32], sender.Bytes())

	if key != nil {
		key.WriteToSlice(buf[32:64])
	}

	return crypto.Keccak256Hash(buf[:])
}

// NonceKeysHash returns EIP-8250's nonce_keys_hash, the value TXPARAM 0x0E reports.
//
// Valid key sets are strictly increasing, so the hash has one canonical form per
// selected key set.
func NonceKeysHash(keys []*uint256.Int) common.Hash {
	buf := make([]byte, 0, 32*(len(keys)+1))

	var word [32]byte

	binary.BigEndian.PutUint64(word[24:], uint64(len(keys)))
	buf = append(buf, word[:]...)

	for _, key := range keys {
		word = [32]byte{}
		if key != nil {
			key.WriteToSlice(word[:])
		}

		buf = append(buf, word[:]...)
	}

	return crypto.Keccak256Hash(buf)
}

// WithNonceKeys selects an EIP-8250 key set and its shared sequence number.
//
// Every selected key must currently sit at seq, and a successful payment approval moves
// all of them to seq+1 together. Keys must be strictly increasing and the zero key may
// only appear alone.
func (tx *FrameTx) WithNonceKeys(keys []*uint256.Int, seq uint64) *FrameTx {
	tx.Extensions |= FrameExtKeyedNonces
	tx.NonceKeys = keys
	tx.NonceSeq = seq

	return tx
}

// RecentRootSourceID returns the identifier of a root source, keccak256(address || salt).
//
// A single address addresses as many independent sources as it has salts.
func RecentRootSourceID(source common.Address, salt common.Hash) common.Hash {
	buf := make([]byte, 0, 52)
	buf = append(buf, source.Bytes()...)
	buf = append(buf, salt.Bytes()...)

	return crypto.Keccak256Hash(buf)
}

// RecentRootEntryHash returns the value the recent root contract commits for a
// (source, slot, root) triple. A reference is valid only when the source's storage
// holds exactly this hash.
func RecentRootEntryHash(sourceID common.Hash, slot uint64, root common.Hash) common.Hash {
	buf := make([]byte, 0, 104)
	buf = append(buf, RecentRootEntryDomain.Bytes()...)
	buf = append(buf, sourceID.Bytes()...)
	buf = binary.BigEndian.AppendUint64(buf, slot)
	buf = append(buf, root.Bytes()...)

	return crypto.Keccak256Hash(buf)
}

// RecentRootStorageKey returns the RecentRootAddress storage key holding a source's
// entry for a ring index. Use RecentRootIndex to derive the index from a slot.
func RecentRootStorageKey(sourceID common.Hash, index uint64) common.Hash {
	buf := make([]byte, 0, 72)
	buf = append(buf, RecentRootStorageDomain.Bytes()...)
	buf = append(buf, sourceID.Bytes()...)
	buf = binary.BigEndian.AppendUint64(buf, index)

	return crypto.Keccak256Hash(buf)
}

// RecentRootIndex returns the ring index a slot's entry occupies.
func RecentRootIndex(slot uint64) uint64 { return slot % RecentRootLength }

// RecentRootSlotStorageKey is the common case of RecentRootStorageKey: the key holding
// whatever a source committed during a given slot.
func RecentRootSlotStorageKey(sourceID common.Hash, slot uint64) common.Hash {
	return RecentRootStorageKey(sourceID, RecentRootIndex(slot))
}

// RecentRootWriteCalldata builds the recent root contract's only call: salt || root,
// exactly RecentRootWriteLength bytes, sent with zero value and not through
// DELEGATECALL, or the call reverts.
func RecentRootWriteCalldata(salt, root common.Hash) []byte {
	data := make([]byte, 0, RecentRootWriteLength)
	data = append(data, salt.Bytes()...)
	data = append(data, root.Bytes()...)

	return data
}

// RecentRootReferenceUsable reports whether a reference to slot may be declared by a
// transaction executing in currentSlot. A root written during slot S becomes usable in
// S+1 and stays usable for RecentRootUsableWindow slots.
func RecentRootReferenceUsable(currentSlot, slot uint64) bool {
	if currentSlot <= slot {
		return false
	}

	return currentSlot-slot <= RecentRootUsableWindow
}
