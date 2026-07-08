package safemultisig

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// EIP-712 type hashes used by the Safe contract (v1.3.0+ / v1.4.1). They match
// the constants compiled into the singleton, so a SafeTx hash computed here is
// identical to the one Safe.getTransactionHash returns on-chain.
var (
	// keccak256("EIP712Domain(uint256 chainId,address verifyingContract)")
	domainSeparatorTypehash = crypto.Keccak256Hash([]byte("EIP712Domain(uint256 chainId,address verifyingContract)"))
	// keccak256("SafeTx(address to,uint256 value,bytes data,uint8 operation,uint256 safeTxGas,uint256 baseGas,uint256 gasPrice,address gasToken,address refundReceiver,uint256 nonce)")
	safeTxTypehash = crypto.Keccak256Hash([]byte("SafeTx(address to,uint256 value,bytes data,uint8 operation,uint256 safeTxGas,uint256 baseGas,uint256 gasPrice,address gasToken,address refundReceiver,uint256 nonce)"))
)

// ABI primitive types reused when ABI-encoding the EIP-712 preimages. Parsed
// once - the inputs are constant so the errors can never trigger.
var (
	abiAddressTy = mustNewABIType("address")
	abiUint256Ty = mustNewABIType("uint256")
	abiUint8Ty   = mustNewABIType("uint8")
	abiBytes32Ty = mustNewABIType("bytes32")
)

func mustNewABIType(t string) abi.Type {
	ty, err := abi.NewType(t, "", nil)
	if err != nil {
		panic(fmt.Sprintf("safemultisig: invalid abi type %q: %v", t, err))
	}
	return ty
}

// safeTxParams holds the execTransaction parameters that are hashed and signed.
// The scenario always uses the zero-refund path (safeTxGas/baseGas/gasPrice = 0,
// gasToken/refundReceiver = zero address), so only to/value/data/operation/nonce
// vary per transaction.
type safeTxParams struct {
	To        common.Address
	Value     *big.Int
	Data      []byte
	Operation uint8
	Nonce     *big.Int
}

// computeDomainSeparator reproduces Safe.domainSeparator() for the given chain
// and safe address: keccak256(abi.encode(DOMAIN_SEPARATOR_TYPEHASH, chainId,
// safe)). The scenario reads the value from-chain at safe creation instead of
// trusting this for the hot path; this is kept for the startup self-check.
func computeDomainSeparator(chainID *big.Int, safe common.Address) ([32]byte, error) {
	args := abi.Arguments{{Type: abiBytes32Ty}, {Type: abiUint256Ty}, {Type: abiAddressTy}}
	packed, err := args.Pack([32]byte(domainSeparatorTypehash), chainID, safe)
	if err != nil {
		return [32]byte{}, fmt.Errorf("could not pack domain separator: %w", err)
	}
	return crypto.Keccak256Hash(packed), nil
}

// computeSafeTxHash reproduces Safe.getTransactionHash for the zero-refund path:
//
//	structHash = keccak256(abi.encode(SAFE_TX_TYPEHASH, to, value, keccak(data),
//	                                  operation, 0, 0, 0, address(0), address(0), nonce))
//	hash       = keccak256(0x19 || 0x01 || domainSeparator || structHash)
//
// domainSeparator is passed in (cached per safe from the on-chain value) so this
// has no RPC dependency.
func computeSafeTxHash(domainSeparator [32]byte, p *safeTxParams) ([32]byte, error) {
	zero := big.NewInt(0)
	structArgs := abi.Arguments{
		{Type: abiBytes32Ty}, // SAFE_TX_TYPEHASH
		{Type: abiAddressTy}, // to
		{Type: abiUint256Ty}, // value
		{Type: abiBytes32Ty}, // keccak(data)
		{Type: abiUint8Ty},   // operation
		{Type: abiUint256Ty}, // safeTxGas
		{Type: abiUint256Ty}, // baseGas
		{Type: abiUint256Ty}, // gasPrice
		{Type: abiAddressTy}, // gasToken
		{Type: abiAddressTy}, // refundReceiver
		{Type: abiUint256Ty}, // nonce
	}
	structPacked, err := structArgs.Pack(
		[32]byte(safeTxTypehash),
		p.To,
		p.Value,
		[32]byte(crypto.Keccak256Hash(p.Data)),
		p.Operation,
		zero, zero, zero,
		common.Address{},
		common.Address{},
		p.Nonce,
	)
	if err != nil {
		return [32]byte{}, fmt.Errorf("could not pack SafeTx struct: %w", err)
	}
	structHash := crypto.Keccak256Hash(structPacked)

	preimage := make([]byte, 0, 2+32+32)
	preimage = append(preimage, 0x19, 0x01)
	preimage = append(preimage, domainSeparator[:]...)
	preimage = append(preimage, structHash[:]...)
	return crypto.Keccak256Hash(preimage), nil
}

// signSafeTx produces the packed Safe signature blob for safeTxHash. Each signer
// contributes a 65-byte ECDSA signature (r || s || v) over the hash directly,
// with v normalized to the 27/28 convention Safe.checkNSignatures expects for
// the ecrecover path. Signatures must be ordered by ascending signer address;
// callers must pass signers already sorted that way (sortWalletsByAddress).
func signSafeTx(safeTxHash [32]byte, signers []*signer) ([]byte, error) {
	sig := make([]byte, 0, len(signers)*65)
	for _, s := range signers {
		if s.key == nil {
			return nil, fmt.Errorf("owner %s has no private key", s.addr.Hex())
		}
		part, err := crypto.Sign(safeTxHash[:], s.key)
		if err != nil {
			return nil, fmt.Errorf("could not sign SafeTx for owner %s: %w", s.addr.Hex(), err)
		}
		part[64] += 27
		sig = append(sig, part...)
	}
	return sig, nil
}

// computeBurnRounds mirrors GasBurner.burn's round derivation so the scenario
// can size the execTransaction gas budget to the exact work a given call will
// do, rather than always reserving the worst case:
//
//	rounds = 1 + (uint256(keccak256(abi.encodePacked(seed))) % maxRounds)
func computeBurnRounds(seed *big.Int, maxRounds uint64) uint64 {
	if maxRounds == 0 {
		maxRounds = 1
	}
	digest := crypto.Keccak256(common.LeftPadBytes(seed.Bytes(), 32))
	d := new(big.Int).SetBytes(digest)
	return 1 + new(big.Int).Mod(d, new(big.Int).SetUint64(maxRounds)).Uint64()
}

// signer pairs an owner address with its signing key. Kept tiny so the scenario
// can sort owners by address without dragging the full wallet type through the
// signing path.
type signer struct {
	addr common.Address
	key  *ecdsa.PrivateKey
}

// sortSignersByAddress sorts signers in ascending address order, as required by
// Safe.checkNSignatures (it enforces strictly increasing recovered owners).
func sortSignersByAddress(signers []*signer) {
	sort.Slice(signers, func(i, j int) bool {
		return signers[i].addr.Cmp(signers[j].addr) < 0
	})
}
