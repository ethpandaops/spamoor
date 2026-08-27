package txtypes

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

func init() {
	RegisterTxType(LegacyTxType, func() TxData { return &LegacyTx{} })
}

// LegacyTx is a pre-EIP-2718 transaction.
//
// Field order and RLP tags mirror the consensus encoding exactly; the struct is
// encoded and decoded as a plain RLP list with no type prefix.
type LegacyTx struct {
	Nonce    uint64          // nonce of sender account
	GasPrice *big.Int        // wei per gas
	Gas      uint64          // gas limit
	To       *common.Address `rlp:"nil"` // nil means contract creation
	Value    *big.Int        // wei amount
	Data     []byte          // contract invocation input data
	V, R, S  *big.Int        // signature values
}

var (
	_ TxData        = (*LegacyTx)(nil)
	_ ECDSASignedTx = (*LegacyTx)(nil)
)

// TxType returns the EIP-2718 type byte.
func (tx *LegacyTx) TxType() byte { return LegacyTxType }

// CopyTx returns a deep copy with all fields initialized.
func (tx *LegacyTx) CopyTx() TxData {
	cpy := &LegacyTx{
		Nonce:    tx.Nonce,
		To:       copyAddressPtr(tx.To),
		Data:     common.CopyBytes(tx.Data),
		Gas:      tx.Gas,
		GasPrice: new(big.Int),
		Value:    new(big.Int),
		V:        new(big.Int),
		R:        new(big.Int),
		S:        new(big.Int),
	}

	setBig(cpy.GasPrice, tx.GasPrice)
	setBig(cpy.Value, tx.Value)
	setBig(cpy.V, tx.V)
	setBig(cpy.R, tx.R)
	setBig(cpy.S, tx.S)

	return cpy
}

// GetChainID returns the chain id encoded in the EIP-155 signature, or zero for
// unprotected transactions.
func (tx *LegacyTx) GetChainID() *big.Int { return deriveLegacyChainID(tx.V) }

// GetNonce returns the sender account nonce.
func (tx *LegacyTx) GetNonce() uint64 { return tx.Nonce }

// GetGas returns the gas limit.
func (tx *LegacyTx) GetGas() uint64 { return tx.Gas }

// GetGasPrice returns the wei per gas the sender pays.
func (tx *LegacyTx) GetGasPrice() *big.Int { return bigOrZero(tx.GasPrice) }

// GetGasFeeCap returns the gas price; legacy transactions have no separate fee cap.
func (tx *LegacyTx) GetGasFeeCap() *big.Int { return bigOrZero(tx.GasPrice) }

// GetGasTipCap returns the gas price; legacy transactions have no separate tip.
func (tx *LegacyTx) GetGasTipCap() *big.Int { return bigOrZero(tx.GasPrice) }

// GetTo returns the recipient, or nil for contract creation.
func (tx *LegacyTx) GetTo() *common.Address { return tx.To }

// GetValue returns the transferred amount in wei.
func (tx *LegacyTx) GetValue() *big.Int { return bigOrZero(tx.Value) }

// GetData returns the transaction calldata.
func (tx *LegacyTx) GetData() []byte { return tx.Data }

// EncodePayload writes the plain RLP list; legacy transactions carry no type prefix.
func (tx *LegacyTx) EncodePayload(w *bytes.Buffer) error {
	return rlp.Encode(w, tx)
}

// DecodePayload parses the plain RLP list.
func (tx *LegacyTx) DecodePayload(b []byte) error {
	return rlp.DecodeBytes(b, tx)
}

// SigningHash returns the EIP-155 signing digest, or the pre-155 digest when chainID
// is nil or zero.
func (tx *LegacyTx) SigningHash(chainID *big.Int) common.Hash {
	if chainID == nil || chainID.Sign() == 0 {
		return rlpHash([]any{
			tx.Nonce,
			tx.GasPrice,
			tx.Gas,
			tx.To,
			tx.Value,
			tx.Data,
		})
	}

	return rlpHash([]any{
		tx.Nonce,
		tx.GasPrice,
		tx.Gas,
		tx.To,
		tx.Value,
		tx.Data,
		chainID, uint(0), uint(0),
	})
}

// GetSignatureValues returns the signature values as encoded.
func (tx *LegacyTx) GetSignatureValues() (v, r, s *big.Int) {
	return tx.V, tx.R, tx.S
}

// SetSignatureValues stores a signature, applying the EIP-155 v offset when chainID
// is set and the pre-155 offset otherwise.
func (tx *LegacyTx) SetSignatureValues(chainID, v, r, s *big.Int) {
	if chainID != nil && chainID.Sign() != 0 {
		// v = yParity + 35 + 2*chainID
		v = new(big.Int).Add(v, big.NewInt(35))
		v.Add(v, new(big.Int).Lsh(chainID, 1))
	} else {
		v = new(big.Int).Add(v, big.NewInt(27))
	}

	tx.V, tx.R, tx.S = v, r, s
}

// deriveLegacyChainID recovers the chain id from an EIP-155 v value. Unprotected
// transactions (v of 27 or 28) return zero.
func deriveLegacyChainID(v *big.Int) *big.Int {
	if v == nil {
		return new(big.Int)
	}

	if v.BitLen() <= 64 {
		value := v.Uint64()
		if value == 27 || value == 28 {
			return new(big.Int)
		}

		if value < 35 {
			return new(big.Int)
		}

		return new(big.Int).SetUint64((value - 35) / 2)
	}

	chainID := new(big.Int).Sub(v, big.NewInt(35))

	return chainID.Rsh(chainID, 1)
}

// legacyYParity extracts the y-parity bit from an encoded legacy v value.
func legacyYParity(v *big.Int) (byte, bool) {
	if v == nil {
		return 0, false
	}

	if v.BitLen() <= 64 {
		value := v.Uint64()
		if value == 27 || value == 28 {
			return byte(value - 27), true
		}

		if value < 35 {
			return 0, false
		}
	}

	// EIP-155: v = yParity + 35 + 2*chainID, so the parity is the low bit of v-35.
	parity := new(big.Int).Sub(v, big.NewInt(35))

	return byte(parity.Bit(0)), true
}
