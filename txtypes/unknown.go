package txtypes

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// UnknownTx represents a transaction of a type this build cannot decode. It keeps the
// generic fields the node reports so such transactions still take part in block
// accounting; only type-specific introspection is unavailable.
//
// It cannot be encoded or signed and exists only as the result of decoding chain data.
type UnknownTx struct {
	Type      byte
	ChainID   *big.Int
	Nonce     uint64
	GasPrice  *big.Int
	GasFeeCap *big.Int
	GasTipCap *big.Int
	Gas       uint64
	To        *common.Address
	Value     *big.Int
	Data      []byte
}

var _ TxData = (*UnknownTx)(nil)

// newUnknownTx builds an UnknownTx from parsed JSON-RPC fields.
func newUnknownTx(txType byte, fields *JSONTxFields) *UnknownTx {
	tx := &UnknownTx{
		Type:      txType,
		ChainID:   jsonBig(fields.ChainID),
		Nonce:     jsonUint64(fields.Nonce),
		GasPrice:  jsonBig(fields.GasPrice),
		GasFeeCap: jsonBig(fields.MaxFeePerGas),
		GasTipCap: jsonBig(fields.MaxPriorityFeePerGas),
		Gas:       jsonUint64(fields.Gas),
		To:        copyAddressPtr(fields.To),
		Value:     jsonBig(fields.Value),
		Data:      jsonBytes(fields.Input),
	}

	// Pre-1559 types report gasPrice only, later ones the caps only. Fill in whichever
	// side is missing so fee accounting has a usable value.
	if tx.GasFeeCap.Sign() == 0 && tx.GasPrice.Sign() != 0 {
		tx.GasFeeCap = new(big.Int).Set(tx.GasPrice)
		tx.GasTipCap = new(big.Int).Set(tx.GasPrice)
	}

	if tx.GasPrice.Sign() == 0 && tx.GasFeeCap.Sign() != 0 {
		tx.GasPrice = new(big.Int).Set(tx.GasFeeCap)
	}

	return tx
}

// TxType returns the EIP-2718 type byte reported by the node.
func (tx *UnknownTx) TxType() byte { return tx.Type }

// CopyTx returns a deep copy with all fields initialized.
func (tx *UnknownTx) CopyTx() TxData {
	cpy := &UnknownTx{
		Type:      tx.Type,
		Nonce:     tx.Nonce,
		Gas:       tx.Gas,
		To:        copyAddressPtr(tx.To),
		Data:      common.CopyBytes(tx.Data),
		ChainID:   new(big.Int),
		GasPrice:  new(big.Int),
		GasFeeCap: new(big.Int),
		GasTipCap: new(big.Int),
		Value:     new(big.Int),
	}

	setBig(cpy.ChainID, tx.ChainID)
	setBig(cpy.GasPrice, tx.GasPrice)
	setBig(cpy.GasFeeCap, tx.GasFeeCap)
	setBig(cpy.GasTipCap, tx.GasTipCap)
	setBig(cpy.Value, tx.Value)

	return cpy
}

// GetChainID returns the chain id reported by the node.
func (tx *UnknownTx) GetChainID() *big.Int { return bigOrZero(tx.ChainID) }

// GetNonce returns the sender account nonce.
func (tx *UnknownTx) GetNonce() uint64 { return tx.Nonce }

// GetGas returns the gas limit.
func (tx *UnknownTx) GetGas() uint64 { return tx.Gas }

// GetGasPrice returns the gas price.
func (tx *UnknownTx) GetGasPrice() *big.Int { return bigOrZero(tx.GasPrice) }

// GetGasFeeCap returns the maximum fee per gas.
func (tx *UnknownTx) GetGasFeeCap() *big.Int { return bigOrZero(tx.GasFeeCap) }

// GetGasTipCap returns the maximum priority fee per gas.
func (tx *UnknownTx) GetGasTipCap() *big.Int { return bigOrZero(tx.GasTipCap) }

// GetTo returns the recipient, or nil.
func (tx *UnknownTx) GetTo() *common.Address { return tx.To }

// GetValue returns the transferred amount in wei.
func (tx *UnknownTx) GetValue() *big.Int { return bigOrZero(tx.Value) }

// GetData returns the transaction calldata.
func (tx *UnknownTx) GetData() []byte { return tx.Data }

// EncodePayload always fails: the wire bytes are not recoverable from the JSON-RPC
// representation of an unknown type.
func (tx *UnknownTx) EncodePayload(_ *bytes.Buffer) error {
	return fmt.Errorf("%w: 0x%02x cannot be re-encoded", ErrTxTypeNotSupported, tx.Type)
}

// DecodePayload always fails: an unknown type has no payload structure to parse.
func (tx *UnknownTx) DecodePayload(_ []byte) error {
	return fmt.Errorf("%w: 0x%02x", ErrTxTypeNotSupported, tx.Type)
}
