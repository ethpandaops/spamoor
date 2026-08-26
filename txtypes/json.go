package txtypes

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/holiman/uint256"
)

// JSONTxFields holds the transaction fields returned by the standard JSON-RPC methods.
// Fields are pointers so absent can be told from zero; Raw carries the original object
// for types with fields not covered here.
type JSONTxFields struct {
	Raw json.RawMessage `json:"-"`

	Type                 *hexutil.Uint64        `json:"type"`
	ChainID              *hexutil.Big           `json:"chainId"`
	Nonce                *hexutil.Uint64        `json:"nonce"`
	GasPrice             *hexutil.Big           `json:"gasPrice"`
	MaxPriorityFeePerGas *hexutil.Big           `json:"maxPriorityFeePerGas"`
	MaxFeePerGas         *hexutil.Big           `json:"maxFeePerGas"`
	MaxFeePerBlobGas     *hexutil.Big           `json:"maxFeePerBlobGas"`
	Gas                  *hexutil.Uint64        `json:"gas"`
	To                   *common.Address        `json:"to"`
	Value                *hexutil.Big           `json:"value"`
	Input                *hexutil.Bytes         `json:"input"`
	AccessList           *AccessList            `json:"accessList"`
	BlobVersionedHashes  []common.Hash          `json:"blobVersionedHashes"`
	AuthorizationList    []SetCodeAuthorization `json:"authorizationList"`
	V                    *hexutil.Big           `json:"v"`
	YParity              *hexutil.Uint64        `json:"yParity"`
	R                    *hexutil.Big           `json:"r"`
	S                    *hexutil.Big           `json:"s"`

	From *common.Address `json:"from"`
	Hash *common.Hash    `json:"hash"`
}

// JSONTxData is implemented by transaction types that can be reconstructed from a
// JSON-RPC transaction object. Types that do not implement it decode as UnknownTx.
type JSONTxData interface {
	// DecodeJSONTx populates the transaction from parsed JSON-RPC fields.
	DecodeJSONTx(fields *JSONTxFields) error
}

// UnmarshalJSONTx decodes a transaction from a JSON-RPC transaction object. The hash
// and sender reported by the node are adopted verbatim, which keeps block processing
// correct for types this build decodes imperfectly and avoids a signature recovery per
// transaction.
func UnmarshalJSONTx(raw json.RawMessage) (*Transaction, error) {
	var fields JSONTxFields
	if err := json.Unmarshal(raw, &fields); err != nil {
		return nil, err
	}

	fields.Raw = raw

	if fields.Hash == nil {
		return nil, errors.New("transaction object is missing hash")
	}

	txType := byte(LegacyTxType)
	if fields.Type != nil {
		txType = byte(*fields.Type)
	}

	inner := decodeJSONInner(txType, &fields)

	tx := NewTx(inner)
	tx.hash.Store(fields.Hash)

	if fields.From != nil {
		tx.SetFrom(*fields.From)
	}

	return tx, nil
}

// decodeJSONInner builds the type-specific content, falling back to UnknownTx when
// the type is unregistered or cannot be built from JSON.
func decodeJSONInner(txType byte, fields *JSONTxFields) TxData {
	unknown := newUnknownTx(txType, fields)

	if !IsTxTypeSupported(txType) {
		return unknown
	}

	inner, err := newTxData(txType)
	if err != nil {
		return unknown
	}

	decodable, ok := inner.(JSONTxData)
	if !ok {
		return unknown
	}

	if err := decodable.DecodeJSONTx(fields); err != nil {
		return unknown
	}

	return inner
}

// jsonSignatureValues returns the signature values, preferring yParity over v for
// typed transactions.
func (f *JSONTxFields) jsonSignatureValues(typed bool) (v, r, s *big.Int) {
	r = jsonBig(f.R)
	s = jsonBig(f.S)

	switch {
	case typed && f.YParity != nil:
		v = new(big.Int).SetUint64(uint64(*f.YParity))
	case f.V != nil:
		v = (*big.Int)(f.V)
	default:
		v = new(big.Int)
	}

	return v, r, s
}

// jsonBig converts an optional hexutil.Big to a big.Int, mapping absent to zero.
func jsonBig(v *hexutil.Big) *big.Int {
	if v == nil {
		return new(big.Int)
	}

	return (*big.Int)(v)
}

// jsonU256 converts an optional hexutil.Big to a uint256.Int, mapping absent to zero.
func jsonU256(v *hexutil.Big) *uint256.Int {
	if v == nil {
		return new(uint256.Int)
	}

	converted, overflow := uint256.FromBig((*big.Int)(v))
	if overflow {
		return new(uint256.Int)
	}

	return converted
}

// jsonUint64 converts an optional hexutil.Uint64, mapping absent to zero.
func jsonUint64(v *hexutil.Uint64) uint64 {
	if v == nil {
		return 0
	}

	return uint64(*v)
}

// jsonBytes converts optional hex bytes, mapping absent to nil.
func jsonBytes(v *hexutil.Bytes) []byte {
	if v == nil {
		return nil
	}

	return *v
}

// jsonAccessList converts an optional access list, mapping absent to nil.
func jsonAccessList(v *AccessList) AccessList {
	if v == nil {
		return nil
	}

	return *v
}

// DecodeJSONTx populates a legacy transaction from JSON-RPC fields.
func (tx *LegacyTx) DecodeJSONTx(fields *JSONTxFields) error {
	v, r, s := fields.jsonSignatureValues(false)

	tx.Nonce = jsonUint64(fields.Nonce)
	tx.GasPrice = jsonBig(fields.GasPrice)
	tx.Gas = jsonUint64(fields.Gas)
	tx.To = copyAddressPtr(fields.To)
	tx.Value = jsonBig(fields.Value)
	tx.Data = jsonBytes(fields.Input)
	tx.V, tx.R, tx.S = v, r, s

	return nil
}

// DecodeJSONTx populates an access list transaction from JSON-RPC fields.
func (tx *AccessListTx) DecodeJSONTx(fields *JSONTxFields) error {
	v, r, s := fields.jsonSignatureValues(true)

	tx.ChainID = jsonBig(fields.ChainID)
	tx.Nonce = jsonUint64(fields.Nonce)
	tx.GasPrice = jsonBig(fields.GasPrice)
	tx.Gas = jsonUint64(fields.Gas)
	tx.To = copyAddressPtr(fields.To)
	tx.Value = jsonBig(fields.Value)
	tx.Data = jsonBytes(fields.Input)
	tx.AccessList = jsonAccessList(fields.AccessList)
	tx.V, tx.R, tx.S = v, r, s

	return nil
}

// DecodeJSONTx populates a dynamic fee transaction from JSON-RPC fields.
func (tx *DynamicFeeTx) DecodeJSONTx(fields *JSONTxFields) error {
	v, r, s := fields.jsonSignatureValues(true)

	tx.ChainID = jsonBig(fields.ChainID)
	tx.Nonce = jsonUint64(fields.Nonce)
	tx.GasTipCap = jsonBig(fields.MaxPriorityFeePerGas)
	tx.GasFeeCap = jsonBig(fields.MaxFeePerGas)
	tx.Gas = jsonUint64(fields.Gas)
	tx.To = copyAddressPtr(fields.To)
	tx.Value = jsonBig(fields.Value)
	tx.Data = jsonBytes(fields.Input)
	tx.AccessList = jsonAccessList(fields.AccessList)
	tx.V, tx.R, tx.S = v, r, s

	return nil
}

// DecodeJSONTx populates a blob transaction from JSON-RPC fields. Sidecars are not
// part of the JSON-RPC representation.
func (tx *BlobTx) DecodeJSONTx(fields *JSONTxFields) error {
	if fields.To == nil {
		return errors.New("blob transaction without recipient")
	}

	v, r, s := fields.jsonSignatureValues(true)

	tx.ChainID = jsonU256(fields.ChainID)
	tx.Nonce = jsonUint64(fields.Nonce)
	tx.GasTipCap = jsonU256(fields.MaxPriorityFeePerGas)
	tx.GasFeeCap = jsonU256(fields.MaxFeePerGas)
	tx.Gas = jsonUint64(fields.Gas)
	tx.To = *fields.To
	tx.Value = jsonU256(fields.Value)
	tx.Data = jsonBytes(fields.Input)
	tx.AccessList = jsonAccessList(fields.AccessList)
	tx.BlobFeeCap = jsonU256(fields.MaxFeePerBlobGas)
	tx.BlobHashes = fields.BlobVersionedHashes
	tx.V = uint256.MustFromBig(v)
	tx.R = uint256.MustFromBig(r)
	tx.S = uint256.MustFromBig(s)

	return nil
}

// DecodeJSONTx populates a set code transaction from JSON-RPC fields.
func (tx *SetCodeTx) DecodeJSONTx(fields *JSONTxFields) error {
	if fields.To == nil {
		return errors.New("set code transaction without recipient")
	}

	v, r, s := fields.jsonSignatureValues(true)

	tx.ChainID = jsonU256(fields.ChainID)
	tx.Nonce = jsonUint64(fields.Nonce)
	tx.GasTipCap = jsonU256(fields.MaxPriorityFeePerGas)
	tx.GasFeeCap = jsonU256(fields.MaxFeePerGas)
	tx.Gas = jsonUint64(fields.Gas)
	tx.To = *fields.To
	tx.Value = jsonU256(fields.Value)
	tx.Data = jsonBytes(fields.Input)
	tx.AccessList = jsonAccessList(fields.AccessList)
	tx.AuthList = fields.AuthorizationList
	tx.V = uint256.MustFromBig(v)
	tx.R = uint256.MustFromBig(r)
	tx.S = uint256.MustFromBig(s)

	return nil
}
