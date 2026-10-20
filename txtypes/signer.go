package txtypes

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
)

// secp256k1halfN is the curve order halved. Signatures with s above it are malleable
// and rejected since EIP-2.
var secp256k1halfN = new(big.Int).Rsh(crypto.S256().Params().N, 1)

var (
	// ErrUnsignableTx is returned when a transaction type provides no way to sign it.
	ErrUnsignableTx = errors.New("transaction type cannot be signed")

	// ErrUnrecoverableSender is returned when a transaction type provides neither an
	// explicit sender nor a recoverable signature.
	ErrUnrecoverableSender = errors.New("transaction type has no recoverable sender")
)

// SignTx signs a transaction with key and returns a new signed transaction.
//
// The signing scheme is chosen from the capability the transaction type implements:
// ECDSASignedTx types get a secp256k1 signature over their signing hash, while
// ExplicitSenderTx types sign their own internal signature material.
func SignTx(tx *Transaction, chainID *big.Int, key *ecdsa.PrivateKey) (*Transaction, error) {
	if key == nil {
		return nil, errors.New("no private key provided")
	}

	inner := tx.inner.CopyTx()

	switch signable := inner.(type) {
	case ECDSASignedTx:
		sigHash := signable.SigningHash(chainID)

		sig, err := crypto.Sign(sigHash[:], key)
		if err != nil {
			return nil, fmt.Errorf("failed signing transaction: %w", err)
		}

		r, s, v, err := decodeSignature(sig)
		if err != nil {
			return nil, err
		}

		signable.SetSignatureValues(chainID, v, r, s)

	case ExplicitSenderTx:
		if err := signable.SignPayload(chainID, key); err != nil {
			return nil, fmt.Errorf("failed signing transaction: %w", err)
		}

	default:
		return nil, fmt.Errorf("%w: 0x%02x", ErrUnsignableTx, inner.TxType())
	}

	signed := NewTx(inner)
	signed.SetFrom(crypto.PubkeyToAddress(key.PublicKey))

	return signed, nil
}

// recoverSender derives the sender address from a transaction's signature.
func recoverSender(tx *Transaction, chainID *big.Int) (common.Address, error) {
	signed, ok := tx.inner.(ECDSASignedTx)
	if !ok {
		return common.Address{}, fmt.Errorf("%w: 0x%02x", ErrUnrecoverableSender, tx.Type())
	}

	v, r, s := signed.GetSignatureValues()
	if v == nil || r == nil || s == nil {
		return common.Address{}, ErrInvalidSig
	}

	// Legacy transactions encode the chain id into v; typed ones carry the raw parity
	// bit and take the chain id from their own field.
	var (
		sigHashChainID = chainID
		yParity        byte
	)

	if tx.Type() == LegacyTxType {
		derived := deriveLegacyChainID(v)
		if derived.Sign() != 0 {
			sigHashChainID = derived
		} else {
			sigHashChainID = nil
		}

		parity, ok := legacyYParity(v)
		if !ok {
			return common.Address{}, ErrInvalidSig
		}

		yParity = parity
	} else {
		if v.BitLen() > 8 || (v.Uint64() != 0 && v.Uint64() != 1) {
			return common.Address{}, ErrInvalidSig
		}

		yParity = byte(v.Uint64())

		if txChainID := tx.inner.GetChainID(); txChainID.Sign() != 0 {
			sigHashChainID = txChainID
		}
	}

	if !crypto.ValidateSignatureValues(yParity, r, s, true) {
		return common.Address{}, ErrInvalidSig
	}

	sigHash := signed.SigningHash(sigHashChainID)

	sig := make([]byte, crypto.SignatureLength)
	r.FillBytes(sig[:32])
	s.FillBytes(sig[32:64])
	sig[64] = yParity

	pubKey, err := crypto.Ecrecover(sigHash[:], sig)
	if err != nil {
		return common.Address{}, err
	}

	if len(pubKey) == 0 || pubKey[0] != 4 {
		return common.Address{}, errors.New("invalid public key recovered from signature")
	}

	var addr common.Address

	copy(addr[:], crypto.Keccak256(pubKey[1:])[12:])

	return addr, nil
}

// decodeSignature splits a 65-byte secp256k1 signature into r, s and the y-parity bit.
func decodeSignature(sig []byte) (r, s, v *big.Int, err error) {
	if len(sig) != crypto.SignatureLength {
		return nil, nil, nil, fmt.Errorf("wrong signature size: %d", len(sig))
	}

	r = new(big.Int).SetBytes(sig[:32])
	s = new(big.Int).SetBytes(sig[32:64])
	v = new(big.Int).SetBytes([]byte{sig[64]})

	if s.Cmp(secp256k1halfN) > 0 {
		return nil, nil, nil, errors.New("signature s value is not canonical")
	}

	return r, s, v, nil
}

// setU256 copies src into dst, treating a nil src as zero.
func setU256(dst, src *uint256.Int) {
	if src != nil {
		dst.Set(src)
	}
}

// u256ToBig converts a uint256 to a big.Int, mapping nil to zero.
func u256ToBig(v *uint256.Int) *big.Int {
	if v == nil {
		return new(big.Int)
	}

	return v.ToBig()
}
