package txtypes

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

func init() {
	RegisterTxType(DynamicFeeTxType, func() TxData { return &DynamicFeeTx{} })
}

// DynamicFeeTx is an EIP-1559 transaction.
type DynamicFeeTx struct {
	ChainID    *big.Int
	Nonce      uint64
	GasTipCap  *big.Int // a.k.a. maxPriorityFeePerGas
	GasFeeCap  *big.Int // a.k.a. maxFeePerGas
	Gas        uint64
	To         *common.Address `rlp:"nil"` // nil means contract creation
	Value      *big.Int
	Data       []byte
	AccessList AccessList
	V, R, S    *big.Int // signature values
}

var (
	_ TxData           = (*DynamicFeeTx)(nil)
	_ ECDSASignedTx    = (*DynamicFeeTx)(nil)
	_ AccessListTxData = (*DynamicFeeTx)(nil)
)

// TxType returns the EIP-2718 type byte.
func (tx *DynamicFeeTx) TxType() byte { return DynamicFeeTxType }

// CopyTx returns a deep copy with all fields initialized.
func (tx *DynamicFeeTx) CopyTx() TxData {
	cpy := &DynamicFeeTx{
		Nonce:      tx.Nonce,
		To:         copyAddressPtr(tx.To),
		Data:       common.CopyBytes(tx.Data),
		Gas:        tx.Gas,
		AccessList: make(AccessList, len(tx.AccessList)),
		ChainID:    new(big.Int),
		GasTipCap:  new(big.Int),
		GasFeeCap:  new(big.Int),
		Value:      new(big.Int),
		V:          new(big.Int),
		R:          new(big.Int),
		S:          new(big.Int),
	}

	copy(cpy.AccessList, tx.AccessList)
	setBig(cpy.ChainID, tx.ChainID)
	setBig(cpy.GasTipCap, tx.GasTipCap)
	setBig(cpy.GasFeeCap, tx.GasFeeCap)
	setBig(cpy.Value, tx.Value)
	setBig(cpy.V, tx.V)
	setBig(cpy.R, tx.R)
	setBig(cpy.S, tx.S)

	return cpy
}

// GetChainID returns the destination chain id.
func (tx *DynamicFeeTx) GetChainID() *big.Int { return bigOrZero(tx.ChainID) }

// GetNonce returns the sender account nonce.
func (tx *DynamicFeeTx) GetNonce() uint64 { return tx.Nonce }

// GetGas returns the gas limit.
func (tx *DynamicFeeTx) GetGas() uint64 { return tx.Gas }

// GetGasPrice returns the fee cap. EIP-1559 transactions have no single gas price;
// the fee cap is the value a sender is charged against in the worst case.
func (tx *DynamicFeeTx) GetGasPrice() *big.Int { return bigOrZero(tx.GasFeeCap) }

// GetGasFeeCap returns the maximum fee per gas.
func (tx *DynamicFeeTx) GetGasFeeCap() *big.Int { return bigOrZero(tx.GasFeeCap) }

// GetGasTipCap returns the maximum priority fee per gas.
func (tx *DynamicFeeTx) GetGasTipCap() *big.Int { return bigOrZero(tx.GasTipCap) }

// GetTo returns the recipient, or nil for contract creation.
func (tx *DynamicFeeTx) GetTo() *common.Address { return tx.To }

// GetValue returns the transferred amount in wei.
func (tx *DynamicFeeTx) GetValue() *big.Int { return bigOrZero(tx.Value) }

// GetData returns the transaction calldata.
func (tx *DynamicFeeTx) GetData() []byte { return tx.Data }

// GetAccessList returns the EIP-2930 access list.
func (tx *DynamicFeeTx) GetAccessList() AccessList { return tx.AccessList }

// EncodePayload writes the RLP payload following the type byte.
func (tx *DynamicFeeTx) EncodePayload(w *bytes.Buffer) error {
	return rlp.Encode(w, tx)
}

// DecodePayload parses the RLP payload following the type byte.
func (tx *DynamicFeeTx) DecodePayload(b []byte) error {
	return rlp.DecodeBytes(b, tx)
}

// SigningHash returns the digest the sender signs.
func (tx *DynamicFeeTx) SigningHash(chainID *big.Int) common.Hash {
	return prefixedRlpHash(DynamicFeeTxType, []any{
		chainID,
		tx.Nonce,
		tx.GasTipCap,
		tx.GasFeeCap,
		tx.Gas,
		tx.To,
		tx.Value,
		tx.Data,
		tx.AccessList,
	})
}

// GetSignatureValues returns the signature values as encoded.
func (tx *DynamicFeeTx) GetSignatureValues() (v, r, s *big.Int) {
	return tx.V, tx.R, tx.S
}

// SetSignatureValues stores a signature. Typed transactions encode v as the raw
// y-parity bit.
func (tx *DynamicFeeTx) SetSignatureValues(chainID, v, r, s *big.Int) {
	tx.ChainID, tx.V, tx.R, tx.S = chainID, v, r, s
}
