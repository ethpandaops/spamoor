package txtypes

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

func init() {
	RegisterTxType(AccessListTxType, func() TxData { return &AccessListTx{} })
}

// AccessListTx is an EIP-2930 transaction.
type AccessListTx struct {
	ChainID    *big.Int
	Nonce      uint64
	GasPrice   *big.Int
	Gas        uint64
	To         *common.Address `rlp:"nil"` // nil means contract creation
	Value      *big.Int
	Data       []byte
	AccessList AccessList
	V, R, S    *big.Int // signature values
}

var (
	_ TxData           = (*AccessListTx)(nil)
	_ ECDSASignedTx    = (*AccessListTx)(nil)
	_ AccessListTxData = (*AccessListTx)(nil)
)

// TxType returns the EIP-2718 type byte.
func (tx *AccessListTx) TxType() byte { return AccessListTxType }

// CopyTx returns a deep copy with all fields initialized.
func (tx *AccessListTx) CopyTx() TxData {
	cpy := &AccessListTx{
		Nonce:      tx.Nonce,
		To:         copyAddressPtr(tx.To),
		Data:       common.CopyBytes(tx.Data),
		Gas:        tx.Gas,
		AccessList: make(AccessList, len(tx.AccessList)),
		ChainID:    new(big.Int),
		GasPrice:   new(big.Int),
		Value:      new(big.Int),
		V:          new(big.Int),
		R:          new(big.Int),
		S:          new(big.Int),
	}

	copy(cpy.AccessList, tx.AccessList)
	setBig(cpy.ChainID, tx.ChainID)
	setBig(cpy.GasPrice, tx.GasPrice)
	setBig(cpy.Value, tx.Value)
	setBig(cpy.V, tx.V)
	setBig(cpy.R, tx.R)
	setBig(cpy.S, tx.S)

	return cpy
}

// GetChainID returns the destination chain id.
func (tx *AccessListTx) GetChainID() *big.Int { return bigOrZero(tx.ChainID) }

// GetNonce returns the sender account nonce.
func (tx *AccessListTx) GetNonce() uint64 { return tx.Nonce }

// GetGas returns the gas limit.
func (tx *AccessListTx) GetGas() uint64 { return tx.Gas }

// GetGasPrice returns the wei per gas the sender pays.
func (tx *AccessListTx) GetGasPrice() *big.Int { return bigOrZero(tx.GasPrice) }

// GetGasFeeCap returns the gas price; type 0x01 has no separate fee cap.
func (tx *AccessListTx) GetGasFeeCap() *big.Int { return bigOrZero(tx.GasPrice) }

// GetGasTipCap returns the gas price; type 0x01 has no separate tip.
func (tx *AccessListTx) GetGasTipCap() *big.Int { return bigOrZero(tx.GasPrice) }

// GetTo returns the recipient, or nil for contract creation.
func (tx *AccessListTx) GetTo() *common.Address { return tx.To }

// GetValue returns the transferred amount in wei.
func (tx *AccessListTx) GetValue() *big.Int { return bigOrZero(tx.Value) }

// GetData returns the transaction calldata.
func (tx *AccessListTx) GetData() []byte { return tx.Data }

// GetAccessList returns the EIP-2930 access list.
func (tx *AccessListTx) GetAccessList() AccessList { return tx.AccessList }

// EncodePayload writes the RLP payload following the type byte.
func (tx *AccessListTx) EncodePayload(w *bytes.Buffer) error {
	return rlp.Encode(w, tx)
}

// DecodePayload parses the RLP payload following the type byte.
func (tx *AccessListTx) DecodePayload(b []byte) error {
	return rlp.DecodeBytes(b, tx)
}

// SigningHash returns the digest the sender signs.
func (tx *AccessListTx) SigningHash(chainID *big.Int) common.Hash {
	return prefixedRlpHash(AccessListTxType, []any{
		chainID,
		tx.Nonce,
		tx.GasPrice,
		tx.Gas,
		tx.To,
		tx.Value,
		tx.Data,
		tx.AccessList,
	})
}

// GetSignatureValues returns the signature values as encoded.
func (tx *AccessListTx) GetSignatureValues() (v, r, s *big.Int) {
	return tx.V, tx.R, tx.S
}

// SetSignatureValues stores a signature. Typed transactions encode v as the raw
// y-parity bit.
func (tx *AccessListTx) SetSignatureValues(chainID, v, r, s *big.Int) {
	tx.ChainID, tx.V, tx.R, tx.S = chainID, v, r, s
}
