package txfuzz

import (
	cryptorand "crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/holiman/uint256"
)

// txKind enumerates the transaction envelope types the fuzzer can emit.
// These map onto the five tx types spamoor's txbuilder/wallet can build.
type txKind int

const (
	kindLegacy     txKind = iota // type 0
	kindAccessList               // type 1 (EIP-2930)
	kindDynFee                   // type 2 (EIP-1559)
	kindBlob                     // type 3 (EIP-4844)
	kindSetCode                  // type 4 (EIP-7702)
)

// allKinds is the ordered list of selectable kinds, used to filter against the
// enabled set and to pick uniformly at random per transaction.
var allKinds = []txKind{kindLegacy, kindAccessList, kindDynFee, kindBlob, kindSetCode}

func (k txKind) String() string {
	switch k {
	case kindLegacy:
		return "legacy"
	case kindAccessList:
		return "accesslist"
	case kindDynFee:
		return "dynfee"
	case kindBlob:
		return "blob"
	case kindSetCode:
		return "setcode"
	default:
		return "unknown"
	}
}

// systemAddresses are interesting precompile-adjacent and system contract targets.
// Calling these with fuzzed data exercises edge paths in the EL that random EOAs don't.
var systemAddresses = []common.Address{
	params.BeaconRootsAddress,
	params.HistoryStorageAddress,
	params.WithdrawalQueueAddress,
	params.ConsolidationQueueAddress,
	params.SystemAddress,
}

// fuzzedTx is the envelope-level description produced by the fuzzer for a single
// transaction. The scenario's send path resolves fees/gas and turns this into a
// concrete signed transaction of the matching type.
type fuzzedTx struct {
	kind       txKind
	to         *common.Address // nil => contract creation (only for legacy/accesslist/dynfee)
	data       []byte
	value      *uint256.Int
	accessList types.AccessList
	authList   []types.SetCodeAuthorization
	blobRefs   [][]string
}

// fuzzer produces deterministic-per-seed, fuzzed transaction envelopes.
type fuzzer struct {
	seed         []byte
	chainID      uint64
	enabledKinds []txKind
	maxCallData  int
	maxAccessLen int
	maxAuthList  int
	maxBlobs     int
	// poolAddrs are recoverable target addresses (child wallets) so value-bearing
	// or simple call txs don't permanently burn funds to unspendable addresses.
	poolAddrs func() []common.Address
}

// rngFor returns a deterministic RNG for a given transaction index, derived from
// the run seed. Same seed + same index always yields the same transaction, which
// makes failures reproducible via --payload-seed.
func (f *fuzzer) rngFor(txIdx uint64) *rand.Rand {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], txIdx)
	h := sha256.Sum256(append(append([]byte(nil), f.seed...), buf[:]...))
	return rand.New(rand.NewSource(int64(binary.BigEndian.Uint64(h[:8]))))
}

// generate builds a fuzzed transaction envelope for the given index.
func (f *fuzzer) generate(txIdx uint64) *fuzzedTx {
	rng := f.rngFor(txIdx)

	kind := f.enabledKinds[rng.Intn(len(f.enabledKinds))]
	tx := &fuzzedTx{
		kind:  kind,
		value: uint256.NewInt(0),
		data:  f.fuzzCallData(rng),
	}

	// Blob and setcode txs require a non-nil recipient; the others may be
	// contract creations (~20% of the time) to also fuzz the creation path.
	requiresTo := kind == kindBlob || kind == kindSetCode
	if requiresTo || rng.Intn(5) != 0 {
		tx.to = f.fuzzTarget(rng)
	}

	// Access lists apply to all typed txs (2930/1559/4844/7702). Legacy can't carry one.
	if kind != kindLegacy {
		tx.accessList = f.fuzzAccessList(rng)
	}

	switch kind {
	case kindSetCode:
		tx.authList = f.fuzzAuthList(rng)
	case kindBlob:
		tx.blobRefs = f.fuzzBlobRefs(rng)
	}

	return tx
}

// fuzzCallData produces a calldata/initcode payload with a mix of shapes:
// empty, small, and larger random blobs up to the configured cap.
func (f *fuzzer) fuzzCallData(rng *rand.Rand) []byte {
	switch rng.Intn(8) {
	case 0:
		return nil // empty calldata
	case 1, 2:
		// small payload (e.g. 4-byte selector + a word or two)
		return randomBytes(rng, rng.Intn(36))
	default:
		if f.maxCallData <= 0 {
			return nil
		}
		return randomBytes(rng, rng.Intn(f.maxCallData))
	}
}

// fuzzTarget picks a recipient address across a spread of interesting buckets:
// recoverable pool wallets, the zero address, low precompile addresses, system
// contracts, and fully random addresses.
func (f *fuzzer) fuzzTarget(rng *rand.Rand) *common.Address {
	switch rng.Intn(10) {
	case 0, 1, 2, 3:
		// recoverable child wallet (keeps any value transfers reclaimable)
		if addrs := f.poolAddrs(); len(addrs) > 0 {
			a := addrs[rng.Intn(len(addrs))]
			return &a
		}
		fallthrough
	case 4:
		a := common.Address{}
		return &a // zero address
	case 5, 6:
		// a precompile address (0x01 .. 0x11)
		a := common.Address{}
		a[19] = byte(1 + rng.Intn(0x11))
		return &a
	case 7:
		a := systemAddresses[rng.Intn(len(systemAddresses))]
		return &a
	default:
		var a common.Address
		copy(a[:], randomBytes(rng, 20))
		return &a
	}
}

// fuzzAccessList builds a random EIP-2930 access list (often empty), with a
// handful of random addresses each carrying a few random storage keys.
func (f *fuzzer) fuzzAccessList(rng *rand.Rand) types.AccessList {
	if f.maxAccessLen <= 0 || rng.Intn(2) == 0 {
		return nil // ~50% empty
	}

	n := rng.Intn(f.maxAccessLen) + 1
	al := make(types.AccessList, 0, n)
	for i := 0; i < n; i++ {
		var addr common.Address
		copy(addr[:], randomBytes(rng, 20))

		keyCount := rng.Intn(f.maxAccessLen + 1)
		keys := make([]common.Hash, 0, keyCount)
		for j := 0; j < keyCount; j++ {
			var key common.Hash
			copy(key[:], randomBytes(rng, 32))
			keys = append(keys, key)
		}

		al = append(al, types.AccessTuple{Address: addr, StorageKeys: keys})
	}
	return al
}

// fuzzAuthList builds a fuzzed EIP-7702 authorization list. Each entry is signed
// by a fresh ephemeral key so it never touches managed wallet nonce state -
// invalid/non-applying authorizations are expected and fine (the tx still mines).
func (f *fuzzer) fuzzAuthList(rng *rand.Rand) []types.SetCodeAuthorization {
	if f.maxAuthList <= 0 {
		return nil
	}

	n := rng.Intn(f.maxAuthList) + 1
	auths := make([]types.SetCodeAuthorization, 0, n)
	for i := 0; i < n; i++ {
		// delegate target address (random, system, or zero)
		delegate := f.fuzzTarget(rng)

		// chainID: mostly the real chain (so some auths can actually apply),
		// occasionally 0 (valid for any chain) or a random garbage value.
		chainID := f.chainID
		switch rng.Intn(4) {
		case 0:
			chainID = 0
		case 1:
			chainID = rng.Uint64()
		}

		auth := types.SetCodeAuthorization{
			ChainID: *uint256.NewInt(chainID),
			Address: *delegate,
			Nonce:   rng.Uint64(),
		}

		// sign with a throwaway key - decoupled from pool nonce accounting
		sk, err := crypto.GenerateKey()
		if err != nil {
			continue
		}
		signed, err := types.SignSetCode(sk, auth)
		if err != nil {
			continue
		}
		auths = append(auths, signed)
	}
	return auths
}

// fuzzBlobRefs produces 1..maxBlobs blob references using spamoor's blob ref DSL,
// mixing fully-random blobs with a few known edge cases (all-zero, repeated,
// duplicate commitments).
func (f *fuzzer) fuzzBlobRefs(rng *rand.Rand) [][]string {
	count := 1
	if f.maxBlobs > 1 {
		count = rng.Intn(f.maxBlobs) + 1
	}

	refs := make([][]string, count)
	for i := 0; i < count; i++ {
		switch rng.Intn(10) {
		case 0:
			refs[i] = []string{"0x0"} // all-zero blob commitment edge case
		case 1:
			refs[i] = []string{"repeat:0x42:1337"} // well-known repeated pattern
		case 2:
			if i > 0 {
				refs[i] = []string{"copy:0"} // duplicate commitment
				continue
			}
			fallthrough
		default:
			refs[i] = []string{"random:full"}
		}
	}
	return refs
}

// randomBytes returns n cryptographically-random bytes. The deterministic RNG
// only drives the structural choices (sizes, buckets); payload bytes use crypto
// randomness so each run still explores fresh inputs at a fixed structure.
func randomBytes(rng *rand.Rand, n int) []byte {
	if n <= 0 {
		return nil
	}
	b := make([]byte, n)
	if _, err := cryptorand.Read(b); err != nil {
		// fall back to the deterministic rng if the system source fails
		for i := range b {
			b[i] = byte(rng.Intn(256))
		}
	}
	return b
}
