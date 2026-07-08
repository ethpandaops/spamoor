package txfuzz

import (
	cryptorand "crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"math"
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
	gas        uint64
	accessList types.AccessList
	authList   []types.SetCodeAuthorization
	blobRefs   [][]string
}

// fuzzer produces deterministic-per-seed, fuzzed transaction envelopes.
type fuzzer struct {
	seed         []byte
	chainID      uint64
	enabledKinds []txKind
	nonBlobKinds []txKind // enabledKinds without kindBlob (fallback when blob rate limited)
	maxCallData  int
	maxAccessLen int
	maxAuthList  int
	maxBlobs     int
	gasLimit     uint64 // configured per-tx gas cap; fuzzed gas never exceeds it
	// poolAddrs are recoverable target addresses (child wallets) so value-bearing
	// or simple call txs don't permanently burn funds to unspendable addresses.
	poolAddrs func() []common.Address
	// deployedAddrs are contracts successfully deployed by earlier fuzzed creation
	// txs; targeting them makes fuzzed calldata execute against real bytecode.
	deployedAddrs func() []common.Address
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
	return f.generateFromKinds(txIdx, f.enabledKinds)
}

// generateNonBlob regenerates the envelope for txIdx restricted to the enabled
// non-blob kinds. Used when the blob rate limiter rejects a blob tx.
func (f *fuzzer) generateNonBlob(txIdx uint64) *fuzzedTx {
	if len(f.nonBlobKinds) == 0 {
		return f.generate(txIdx)
	}
	return f.generateFromKinds(txIdx, f.nonBlobKinds)
}

// generateFromKinds builds a fuzzed transaction envelope for the given index,
// picking the envelope type from the given kind set.
func (f *fuzzer) generateFromKinds(txIdx uint64, kinds []txKind) *fuzzedTx {
	rng := f.rngFor(txIdx)

	kind := kinds[rng.Intn(len(kinds))]
	tx := &fuzzedTx{
		kind:  kind,
		value: uint256.NewInt(0),
		data:  f.fuzzCallData(rng),
	}

	// Blob and setcode txs require a non-nil recipient; the others may be
	// contract creations (~20% of the time) to also fuzz the creation path.
	requiresTo := kind == kindBlob || kind == kindSetCode
	recoverable := false
	if requiresTo || rng.Intn(5) != 0 {
		tx.to, recoverable = f.fuzzTarget(rng)
	}

	// Value: only send meaningful value when the target is a pool wallet, so the
	// funds move between child wallets and stay reclaimable. Non-recoverable
	// targets occasionally get 1 wei to still exercise value-transfer paths.
	if recoverable {
		switch rng.Intn(4) {
		case 0:
			// keep zero
		case 1:
			tx.value = uint256.NewInt(1)
		default:
			tx.value = uint256.NewInt(uint64(rng.Int63n(1_000_000_000_000_000)) + 1) // up to ~0.001 ETH
		}
	} else if tx.to != nil && rng.Intn(20) == 0 {
		tx.value = uint256.NewInt(1)
	}

	// Access lists apply to all typed txs (2930/1559/4844/7702). Legacy can't carry one.
	if kind != kindLegacy {
		tx.accessList = f.fuzzAccessList(rng)
	}

	switch kind {
	case kindSetCode:
		tx.authList = f.fuzzAuthList(rng)
		// ~1/3 of setcode txs additionally carry an authorization that actually
		// applies (fresh ephemeral authority, nonce 0, real chain id, delegate
		// with code) and call the authority itself so the delegated code
		// executes within this very tx.
		if rng.Intn(3) == 0 {
			if auth, authority := f.applyingAuth(rng); auth != nil {
				tx.authList = append(tx.authList, *auth)
				tx.to = authority
				tx.value = uint256.NewInt(0) // authority key is discarded; value would be burned
			}
		}
	case kindBlob:
		tx.blobRefs = f.fuzzBlobRefs(rng)
	}

	tx.gas = f.fuzzGas(rng, tx)

	return tx
}

// intrinsicGasFloor approximates the minimum gas a tx of this shape needs to be
// accepted by the pool (intrinsic gas). Staying at or above it keeps every
// fuzzed tx includable; the "exact floor" gas bucket then exercises
// out-of-gas-on-first-op paths without producing invalid transactions.
func intrinsicGasFloor(tx *fuzzedTx) uint64 {
	gas := uint64(21000)
	gas += 16 * uint64(len(tx.data))
	for _, tuple := range tx.accessList {
		gas += 2400
		gas += 1900 * uint64(len(tuple.StorageKeys))
	}
	if tx.to == nil {
		// contract creation: base creation cost + EIP-3860 per-word initcode cost
		gas += 32000
		gas += 2 * ((uint64(len(tx.data)) + 31) / 32)
	}
	// EIP-7702 charges PER_EMPTY_ACCOUNT_COST per authorization tuple
	gas += 25000 * uint64(len(tx.authList))
	return gas
}

// fuzzGas picks a gas limit across buckets: the exact intrinsic floor, floor
// plus a little slack, a mid-range value, and the configured cap. The result is
// always >= the intrinsic floor (so the tx stays includable) and <= the
// configured gas limit cap whenever the cap itself is above the floor.
func (f *fuzzer) fuzzGas(rng *rand.Rand, tx *fuzzedTx) uint64 {
	floor := intrinsicGasFloor(tx)
	gasCap := f.gasLimit
	if gasCap <= floor {
		return floor
	}

	var gas uint64
	switch rng.Intn(8) {
	case 0:
		gas = floor // exact intrinsic floor: any execution work runs out of gas
	case 1, 2:
		gas = floor + uint64(rng.Int63n(50000)) + 1 // floor plus small slack
	case 3, 4:
		span := gasCap - floor
		if span > math.MaxInt64 {
			span = math.MaxInt64
		}
		gas = floor + uint64(rng.Int63n(int64(span))) // anywhere in range
	default:
		gas = gasCap
	}
	if gas > gasCap {
		gas = gasCap
	}
	return gas
}

// fuzzCallData produces a calldata/initcode payload with a mix of shapes:
// empty, small, near-limit, and larger random blobs up to the configured cap.
func (f *fuzzer) fuzzCallData(rng *rand.Rand) []byte {
	switch rng.Intn(10) {
	case 0:
		return nil // empty calldata
	case 1, 2:
		// small payload (e.g. 4-byte selector + a word or two)
		return randomBytes(rng, rng.Intn(36))
	case 3:
		// near the configured maximum (low probability so it doesn't dominate)
		if f.maxCallData <= 0 {
			return nil
		}
		slack := rng.Intn(33)
		if slack >= f.maxCallData {
			slack = 0
		}
		return randomBytes(rng, f.maxCallData-slack)
	default:
		if f.maxCallData <= 0 {
			return nil
		}
		return randomBytes(rng, rng.Intn(f.maxCallData))
	}
}

// fuzzTarget picks a recipient address across a spread of interesting buckets:
// recoverable pool wallets, contracts deployed by earlier fuzzed creations, the
// zero address, low precompile addresses, system contracts, and fully random
// addresses. The second return value is true only for pool-wallet targets, i.e.
// when value sent to the address stays reclaimable.
func (f *fuzzer) fuzzTarget(rng *rand.Rand) (*common.Address, bool) {
	switch rng.Intn(10) {
	case 0, 1, 2:
		// recoverable child wallet (keeps any value transfers reclaimable)
		if addrs := f.poolAddrs(); len(addrs) > 0 {
			a := addrs[rng.Intn(len(addrs))]
			return &a, true
		}
		fallthrough
	case 3:
		a := common.Address{}
		return &a, false // zero address
	case 4, 5:
		// a contract deployed by an earlier fuzzed creation tx, so the fuzzed
		// calldata executes against real bytecode
		if f.deployedAddrs != nil {
			if addrs := f.deployedAddrs(); len(addrs) > 0 {
				a := addrs[rng.Intn(len(addrs))]
				return &a, false
			}
		}
		fallthrough
	case 6:
		// a precompile address (0x01 .. 0x11)
		a := common.Address{}
		a[19] = byte(1 + rng.Intn(0x11))
		return &a, false
	case 7:
		a := systemAddresses[rng.Intn(len(systemAddresses))]
		return &a, false
	default:
		var a common.Address
		copy(a[:], randomBytes(rng, 20))
		return &a, false
	}
}

// fuzzAccessList builds a random EIP-2930 access list (often empty), with a
// handful of random addresses each carrying a few random storage keys.
func (f *fuzzer) fuzzAccessList(rng *rand.Rand) types.AccessList {
	if f.maxAccessLen <= 0 || rng.Intn(2) == 0 {
		return nil // ~50% empty
	}

	// occasionally emit the near-limit shape: max entries, each fully loaded
	nearLimit := rng.Intn(10) == 0

	n := rng.Intn(f.maxAccessLen) + 1
	if nearLimit {
		n = f.maxAccessLen
	}
	al := make(types.AccessList, 0, n)
	for i := 0; i < n; i++ {
		var addr common.Address
		copy(addr[:], randomBytes(rng, 20))

		keyCount := rng.Intn(f.maxAccessLen + 1)
		if nearLimit {
			keyCount = f.maxAccessLen
		}
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
		delegate, _ := f.fuzzTarget(rng)

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

// applyingAuth builds an EIP-7702 authorization that actually applies on chain:
// the authority is a fresh ephemeral key (on-chain nonce 0, never part of the
// managed wallet pool, so the nonce bump from applying it cannot corrupt pool
// nonce tracking), the chain id is the real one, and the delegate is an address
// known to hold code (a fuzz-deployed contract if available, else a system
// contract). Returns the signed authorization plus the authority address so the
// caller can target it and execute the delegated code in the same tx.
func (f *fuzzer) applyingAuth(rng *rand.Rand) (*types.SetCodeAuthorization, *common.Address) {
	var delegate common.Address
	if f.deployedAddrs != nil {
		if addrs := f.deployedAddrs(); len(addrs) > 0 {
			delegate = addrs[rng.Intn(len(addrs))]
		}
	}
	if delegate == (common.Address{}) {
		// system contracts with code (skip the last entry, params.SystemAddress,
		// which has no code deployed)
		delegate = systemAddresses[rng.Intn(len(systemAddresses)-1)]
	}

	sk, err := crypto.GenerateKey()
	if err != nil {
		return nil, nil
	}
	authority := crypto.PubkeyToAddress(sk.PublicKey)
	signed, err := types.SignSetCode(sk, types.SetCodeAuthorization{
		ChainID: *uint256.NewInt(f.chainID),
		Address: delegate,
		Nonce:   0, // fresh key => on-chain nonce 0 => the authorization applies
	})
	if err != nil {
		return nil, nil
	}
	return &signed, &authority
}

// fuzzBlobRefs produces 1..maxBlobs blob references using spamoor's blob ref DSL,
// mixing fully-random blobs with a few known edge cases (all-zero, repeated,
// duplicate commitments).
func (f *fuzzer) fuzzBlobRefs(rng *rand.Rand) [][]string {
	count := 1
	if f.maxBlobs > 1 {
		count = rng.Intn(f.maxBlobs) + 1
		if rng.Intn(10) == 0 {
			count = f.maxBlobs // occasionally use the full blob budget
		}
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
