package txfuzzinvalid

import (
	"crypto/ecdsa"
	cryptorand "crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/holiman/uint256"
)

// invalidCategory identifies a class of deliberately-malformed transaction.
type invalidCategory struct {
	name string
	// canStick is true when a node might accept the tx into its mempool at the
	// account's current nonce (occupying it) rather than rejecting it outright.
	// Those are the only ones that require an unstuck pass before wallet reuse.
	canStick bool
	// alwaysInvalid is true for categories that are structurally invalid
	// regardless of account/network state, so a node accepting one is a genuine
	// finding. State-dependent categories (future/low nonce, underpriced) can be
	// legitimately accepted - e.g. a future nonce queues normally, a too-low
	// nonce equals the current nonce on a fresh account, an underpriced tx is
	// fine when the base fee is ~0 - so their acceptance is NOT flagged.
	alwaysInvalid bool
}

var categories = []invalidCategory{
	{name: "chainid", alwaysInvalid: true},     // wrong chain id in signature
	{name: "lowgas", alwaysInvalid: true},      // gas below intrinsic minimum
	{name: "underpriced", canStick: true},      // fee cap below base fee (state-dependent)
	{name: "futurenonce"},                      // nonce far ahead (queued, dangling - expected)
	{name: "noncetoolow"},                      // nonce below account nonce (state-dependent)
	{name: "nofunds"},                          // value exceeds balance (state-dependent)
	{name: "gasoverflow", alwaysInvalid: true}, // gas limit far above block limit
	{name: "emptyauth", alwaysInvalid: true},   // EIP-7702 setcode tx with empty auth list
	{name: "badblob", alwaysInvalid: true},     // EIP-4844 blob tx with no blob hashes
	{name: "malformed", alwaysInvalid: true},   // corrupted RLP bytes
	{name: "truncated", alwaysInvalid: true},   // truncated RLP bytes
}

// categoriesByName indexes the category list for --categories filtering.
var categoriesByName = func() map[string]invalidCategory {
	m := make(map[string]invalidCategory, len(categories))
	for _, c := range categories {
		m[c.name] = c
	}
	return m
}()

// genInput carries the per-transaction context the generator needs to craft an
// invalid transaction for a specific burner wallet.
type genInput struct {
	key           *ecdsa.PrivateKey
	from          common.Address
	chainID       uint64
	baseNonce     uint64
	balance       *big.Int
	blockGasLimit uint64
}

// genResult is a single fuzzed, invalid transaction ready for raw submission.
type genResult struct {
	category invalidCategory
	raw      []byte
}

// rngFor returns a deterministic RNG for a transaction index, derived from the
// run seed so a failing submission can be reproduced with --payload-seed.
func rngFor(seed []byte, txIdx uint64) *rand.Rand {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], txIdx)
	h := sha256.Sum256(append(append([]byte(nil), seed...), buf[:]...))
	return rand.New(rand.NewSource(int64(binary.BigEndian.Uint64(h[:8]))))
}

// generateInvalid crafts one invalid transaction from the enabled category set.
func generateInvalid(rng *rand.Rand, enabled []invalidCategory, in genInput) (*genResult, error) {
	cat := enabled[rng.Intn(len(enabled))]

	const (
		gwei     = 1_000_000_000
		baseFee  = 20 * gwei
		baseTip  = 2 * gwei
		validGas = 21000
	)

	chainID := in.chainID
	nonce := in.baseNonce
	gasFee := big.NewInt(baseFee)
	gasTip := big.NewInt(baseTip)
	gas := uint64(validGas)
	to := randomTarget(rng)
	value := uint256.NewInt(0)

	switch cat.name {
	case "chainid":
		// pick a chain id that is definitely not the real one
		chainID = rng.Uint64() | 1
		if chainID == in.chainID {
			chainID = in.chainID + 1
		}
	case "lowgas":
		gas = uint64(rng.Intn(validGas)) // 0 .. 20999, below intrinsic
	case "underpriced":
		gasFee = big.NewInt(int64(rng.Intn(2))) // 0 or 1 wei
		gasTip = big.NewInt(0)
	case "futurenonce":
		nonce = in.baseNonce + uint64(1_000_000) + uint64(rng.Intn(1_000_000))
	case "noncetoolow":
		if in.baseNonce > 0 {
			nonce = in.baseNonce - 1
		} else {
			nonce = 0
		}
	case "nofunds":
		// value far above any plausible balance
		v := new(big.Int).Add(in.balance, new(big.Int).Mul(big.NewInt(1000), big.NewInt(params.Ether)))
		value = uint256.MustFromBig(v)
	case "gasoverflow":
		gas = in.blockGasLimit
		if gas == 0 {
			gas = 36_000_000
		}
		gas = gas*100 + uint64(rng.Intn(1_000_000))
	case "emptyauth":
		return signEnvelope(in.key, chainID, &types.SetCodeTx{
			ChainID:   uint256.NewInt(chainID),
			Nonce:     nonce,
			GasTipCap: uint256.MustFromBig(gasTip),
			GasFeeCap: uint256.MustFromBig(gasFee),
			Gas:       gas,
			To:        to,
			Value:     value,
			AuthList:  []types.SetCodeAuthorization{}, // empty -> invalid 7702
		}, cat)
	case "badblob":
		return signEnvelope(in.key, chainID, &types.BlobTx{
			ChainID:    uint256.NewInt(chainID),
			Nonce:      nonce,
			GasTipCap:  uint256.MustFromBig(gasTip),
			GasFeeCap:  uint256.MustFromBig(gasFee),
			Gas:        gas,
			To:         to,
			Value:      value,
			BlobFeeCap: uint256.NewInt(gwei),
			BlobHashes: []common.Hash{}, // no blobs -> invalid 4844
		}, cat)
	case "malformed", "truncated":
		// build a well-formed dynfee tx, then corrupt or truncate its bytes
		res, err := signEnvelope(in.key, chainID, &types.DynamicFeeTx{
			ChainID:   big.NewInt(int64(chainID)),
			Nonce:     nonce,
			GasTipCap: gasTip,
			GasFeeCap: gasFee,
			Gas:       gas,
			To:        &to,
			Value:     value.ToBig(),
		}, cat)
		if err != nil {
			return nil, err
		}
		if cat.name == "truncated" {
			res.raw = truncateBytes(rng, res.raw)
		} else {
			res.raw = corruptBytes(rng, res.raw)
		}
		return res, nil
	}

	// default path: a properly-signed dynfee tx whose fields make it invalid
	return signEnvelope(in.key, chainID, &types.DynamicFeeTx{
		ChainID:   big.NewInt(int64(chainID)),
		Nonce:     nonce,
		GasTipCap: gasTip,
		GasFeeCap: gasFee,
		Gas:       gas,
		To:        &to,
		Value:     value.ToBig(),
	}, cat)
}

// signEnvelope signs the given tx data with the burner key (using a signer for
// the provided chainID, which may be deliberately wrong) and returns its raw
// network encoding.
func signEnvelope(key *ecdsa.PrivateKey, chainID uint64, data types.TxData, cat invalidCategory) (*genResult, error) {
	tx := types.NewTx(data)
	signer := types.LatestSignerForChainID(new(big.Int).SetUint64(chainID))
	signed, err := types.SignTx(tx, signer, key)
	if err != nil {
		return nil, err
	}
	raw, err := signed.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return &genResult{category: cat, raw: raw}, nil
}

// corruptBytes flips a handful of random bytes in the encoding so it fails to
// decode or recover a valid sender.
func corruptBytes(rng *rand.Rand, b []byte) []byte {
	if len(b) == 0 {
		return b
	}
	out := append([]byte(nil), b...)
	flips := 1 + rng.Intn(8)
	for i := 0; i < flips; i++ {
		out[rng.Intn(len(out))] ^= byte(1 + rng.Intn(255))
	}
	return out
}

// truncateBytes cuts off the tail of the encoding.
func truncateBytes(rng *rand.Rand, b []byte) []byte {
	if len(b) <= 2 {
		return b
	}
	keep := 1 + rng.Intn(len(b)-1)
	return append([]byte(nil), b[:keep]...)
}

// randomTarget returns a random recipient across interesting buckets.
func randomTarget(rng *rand.Rand) common.Address {
	switch rng.Intn(6) {
	case 0:
		return common.Address{} // zero address
	case 1:
		var a common.Address
		a[19] = byte(1 + rng.Intn(0x11)) // precompile range
		return a
	case 2:
		return params.BeaconRootsAddress
	default:
		var a common.Address
		b := make([]byte, 20)
		if _, err := cryptorand.Read(b); err != nil {
			for i := range b {
				b[i] = byte(rng.Intn(256))
			}
		}
		copy(a[:], b)
		return a
	}
}
