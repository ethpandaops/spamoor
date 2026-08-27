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

// JSONTxEncoder is implemented by transaction types that can render themselves as a
// JSON-RPC transaction object. Types that do not implement it marshal with the common
// fields only.
type JSONTxEncoder interface {
	// EncodeJSONTx adds the type's fields to the object.
	EncodeJSONTx(fields map[string]any)
}

// UnmarshalJSONTx decodes a transaction from a JSON-RPC transaction object. The hash
// and sender reported by the node are adopted verbatim, which keeps block processing
// correct for types this build decodes imperfectly and avoids a signature recovery per
// transaction.
func UnmarshalJSONTx(raw json.RawMessage) (*Transaction, error) {
	tx := &Transaction{}
	if err := tx.UnmarshalJSON(raw); err != nil {
		return nil, err
	}

	return tx, nil
}

// UnmarshalJSON decodes a transaction from a JSON-RPC transaction object.
//
// Transaction carries only unexported fields, so without this the standard decoder
// would quietly leave the receiver empty and every later method call would panic.
func (tx *Transaction) UnmarshalJSON(input []byte) error {
	var fields JSONTxFields
	if err := json.Unmarshal(input, &fields); err != nil {
		return err
	}

	fields.Raw = input

	if fields.Hash == nil {
		return errors.New("transaction object is missing hash")
	}

	txType := byte(LegacyTxType)
	if fields.Type != nil {
		txType = byte(*fields.Type)
	}

	tx.setDecoded(decodeJSONInner(txType, &fields))
	tx.hash.Store(fields.Hash)

	if fields.From != nil {
		tx.SetFrom(*fields.From)
	}

	return nil
}

// MarshalJSON renders the transaction as a JSON-RPC transaction object.
//
// The common fields are filled from the accessor set; the type contributes the rest
// through JSONTxEncoder, mirroring how DecodeJSONTx reads them back.
func (tx *Transaction) MarshalJSON() ([]byte, error) {
	if tx.inner == nil {
		return nil, errors.New("cannot marshal an empty transaction")
	}

	fields := map[string]any{
		"type":  hexutil.Uint64(tx.Type()),
		"hash":  tx.Hash(),
		"nonce": hexutil.Uint64(tx.Nonce()),
		"gas":   hexutil.Uint64(tx.Gas()),
		"value": (*hexutil.Big)(tx.Value()),
		"input": hexutil.Bytes(tx.Data()),
		"to":    tx.To(),
	}

	if from := tx.from.Load(); from != nil {
		fields["from"] = *from
	}

	if encoder, ok := tx.inner.(JSONTxEncoder); ok {
		encoder.EncodeJSONTx(fields)
	}

	return json.Marshal(fields)
}

// jsonSignatureFields adds the signature values of an ECDSA-signed transaction.
// Typed transactions report the y-parity bit alongside v, as nodes do.
func jsonSignatureFields(fields map[string]any, typed bool, v, r, s *big.Int) {
	fields["v"] = (*hexutil.Big)(v)
	fields["r"] = (*hexutil.Big)(r)
	fields["s"] = (*hexutil.Big)(s)

	if typed && v != nil {
		fields["yParity"] = hexutil.Uint64(v.Uint64())
	}
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

// EncodeJSONTx adds the legacy transaction's fields.
func (tx *LegacyTx) EncodeJSONTx(fields map[string]any) {
	fields["gasPrice"] = (*hexutil.Big)(tx.GetGasPrice())
	jsonSignatureFields(fields, false, tx.V, tx.R, tx.S)
}

// EncodeJSONTx adds the access list transaction's fields.
func (tx *AccessListTx) EncodeJSONTx(fields map[string]any) {
	fields["chainId"] = (*hexutil.Big)(tx.GetChainID())
	fields["gasPrice"] = (*hexutil.Big)(tx.GetGasPrice())
	fields["accessList"] = tx.AccessList
	jsonSignatureFields(fields, true, tx.V, tx.R, tx.S)
}

// EncodeJSONTx adds the dynamic fee transaction's fields.
func (tx *DynamicFeeTx) EncodeJSONTx(fields map[string]any) {
	fields["chainId"] = (*hexutil.Big)(tx.GetChainID())
	fields["maxFeePerGas"] = (*hexutil.Big)(tx.GetGasFeeCap())
	fields["maxPriorityFeePerGas"] = (*hexutil.Big)(tx.GetGasTipCap())
	fields["accessList"] = tx.AccessList
	jsonSignatureFields(fields, true, tx.V, tx.R, tx.S)
}

// EncodeJSONTx adds the blob transaction's fields. The sidecar is wire-only and has no
// JSON-RPC representation.
func (tx *BlobTx) EncodeJSONTx(fields map[string]any) {
	v, r, s := tx.GetSignatureValues()

	fields["chainId"] = (*hexutil.Big)(tx.GetChainID())
	fields["maxFeePerGas"] = (*hexutil.Big)(tx.GetGasFeeCap())
	fields["maxPriorityFeePerGas"] = (*hexutil.Big)(tx.GetGasTipCap())
	fields["maxFeePerBlobGas"] = (*hexutil.Big)(tx.GetBlobGasFeeCap())
	fields["blobVersionedHashes"] = tx.BlobHashes
	fields["accessList"] = tx.AccessList
	jsonSignatureFields(fields, true, v, r, s)
}

// EncodeJSONTx adds the set code transaction's fields.
func (tx *SetCodeTx) EncodeJSONTx(fields map[string]any) {
	v, r, s := tx.GetSignatureValues()

	fields["chainId"] = (*hexutil.Big)(tx.GetChainID())
	fields["maxFeePerGas"] = (*hexutil.Big)(tx.GetGasFeeCap())
	fields["maxPriorityFeePerGas"] = (*hexutil.Big)(tx.GetGasTipCap())
	fields["accessList"] = tx.AccessList
	fields["authorizationList"] = tx.AuthList
	jsonSignatureFields(fields, true, v, r, s)
}

// EncodeJSONTx adds the fee fields an unknown type reported. Its type-specific content
// was never decoded, so nothing else can be rendered.
func (tx *UnknownTx) EncodeJSONTx(fields map[string]any) {
	fields["chainId"] = (*hexutil.Big)(tx.GetChainID())
	fields["gasPrice"] = (*hexutil.Big)(tx.GetGasPrice())
	fields["maxFeePerGas"] = (*hexutil.Big)(tx.GetGasFeeCap())
	fields["maxPriorityFeePerGas"] = (*hexutil.Big)(tx.GetGasTipCap())
}
