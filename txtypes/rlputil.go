package txtypes

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

// hasherPool reuses keccak states across hashing calls, mirroring go-ethereum's
// approach. Transaction hashing happens once per transaction on submission and once
// per transaction in every processed block, so the allocation matters.
var hasherPool = sync.Pool{
	New: func() any {
		return crypto.NewKeccakState()
	},
}

// rlpHash returns keccak256(rlp(x)).
func rlpHash(x any) (h common.Hash) {
	sha, _ := hasherPool.Get().(crypto.KeccakState)
	defer hasherPool.Put(sha)

	sha.Reset()

	if err := rlp.Encode(sha, x); err != nil {
		return common.Hash{}
	}

	_, _ = sha.Read(h[:])

	return h
}

// prefixedRlpHash returns keccak256(prefix || rlp(x)), the digest form used by all
// EIP-2718 typed transactions.
func prefixedRlpHash(prefix byte, x any) (h common.Hash) {
	sha, _ := hasherPool.Get().(crypto.KeccakState)
	defer hasherPool.Put(sha)

	sha.Reset()
	sha.Write([]byte{prefix})

	if err := rlp.Encode(sha, x); err != nil {
		return common.Hash{}
	}

	_, _ = sha.Read(h[:])

	return h
}

// copyAddressPtr copies an address pointer.
func copyAddressPtr(a *common.Address) *common.Address {
	if a == nil {
		return nil
	}

	cpy := *a

	return &cpy
}

// setBig copies src into dst, treating a nil src as zero. Decoders and copies use it
// so that every big.Int field on a transaction is non-nil, which keeps RLP encoding
// deterministic and accessor call sites free of nil checks.
func setBig(dst, src *big.Int) {
	if src != nil {
		dst.Set(src)
	}
}

// bigOrZero returns v, or a zero big.Int when v is nil.
func bigOrZero(v *big.Int) *big.Int {
	if v == nil {
		return new(big.Int)
	}

	return v
}
