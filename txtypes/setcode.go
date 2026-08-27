package txtypes

import (
	"bytes"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/holiman/uint256"
)

func init() {
	RegisterTxType(SetCodeTxType, func() TxData { return &SetCodeTx{} })
}

// DelegationPrefix is the EIP-7702 code prefix marking a delegation indicator.
var DelegationPrefix = types.DelegationPrefix

// SetCodeTx is an EIP-7702 transaction.
type SetCodeTx struct {
	ChainID    *uint256.Int
	Nonce      uint64
	GasTipCap  *uint256.Int // a.k.a. maxPriorityFeePerGas
	GasFeeCap  *uint256.Int // a.k.a. maxFeePerGas
	Gas        uint64
	To         common.Address
	Value      *uint256.Int
	Data       []byte
	AccessList AccessList
	AuthList   []SetCodeAuthorization

	V, R, S *uint256.Int // signature values
}

var (
	_ TxData           = (*SetCodeTx)(nil)
	_ ECDSASignedTx    = (*SetCodeTx)(nil)
	_ AccessListTxData = (*SetCodeTx)(nil)
	_ AuthListTxData   = (*SetCodeTx)(nil)
)

// TxType returns the EIP-2718 type byte.
func (tx *SetCodeTx) TxType() byte { return SetCodeTxType }

// CopyTx returns a deep copy with all fields initialized.
func (tx *SetCodeTx) CopyTx() TxData {
	cpy := &SetCodeTx{
		Nonce:      tx.Nonce,
		To:         tx.To,
		Data:       common.CopyBytes(tx.Data),
		Gas:        tx.Gas,
		AccessList: make(AccessList, len(tx.AccessList)),
		AuthList:   make([]SetCodeAuthorization, len(tx.AuthList)),
		ChainID:    new(uint256.Int),
		GasTipCap:  new(uint256.Int),
		GasFeeCap:  new(uint256.Int),
		Value:      new(uint256.Int),
		V:          new(uint256.Int),
		R:          new(uint256.Int),
		S:          new(uint256.Int),
	}

	copy(cpy.AccessList, tx.AccessList)
	copy(cpy.AuthList, tx.AuthList)
	setU256(cpy.ChainID, tx.ChainID)
	setU256(cpy.GasTipCap, tx.GasTipCap)
	setU256(cpy.GasFeeCap, tx.GasFeeCap)
	setU256(cpy.Value, tx.Value)
	setU256(cpy.V, tx.V)
	setU256(cpy.R, tx.R)
	setU256(cpy.S, tx.S)

	return cpy
}

// GetChainID returns the destination chain id.
func (tx *SetCodeTx) GetChainID() *big.Int { return u256ToBig(tx.ChainID) }

// GetNonce returns the sender account nonce.
func (tx *SetCodeTx) GetNonce() uint64 { return tx.Nonce }

// GetGas returns the gas limit.
func (tx *SetCodeTx) GetGas() uint64 { return tx.Gas }

// GetGasPrice returns the fee cap; type 0x04 has no single gas price.
func (tx *SetCodeTx) GetGasPrice() *big.Int { return u256ToBig(tx.GasFeeCap) }

// GetGasFeeCap returns the maximum fee per gas.
func (tx *SetCodeTx) GetGasFeeCap() *big.Int { return u256ToBig(tx.GasFeeCap) }

// GetGasTipCap returns the maximum priority fee per gas.
func (tx *SetCodeTx) GetGasTipCap() *big.Int { return u256ToBig(tx.GasTipCap) }

// GetTo returns the recipient. Set code transactions cannot create contracts, so this
// is never nil.
func (tx *SetCodeTx) GetTo() *common.Address {
	to := tx.To

	return &to
}

// GetValue returns the transferred amount in wei.
func (tx *SetCodeTx) GetValue() *big.Int { return u256ToBig(tx.Value) }

// GetData returns the transaction calldata.
func (tx *SetCodeTx) GetData() []byte { return tx.Data }

// GetAccessList returns the EIP-2930 access list.
func (tx *SetCodeTx) GetAccessList() AccessList { return tx.AccessList }

// GetAuthList returns the EIP-7702 authorization list.
func (tx *SetCodeTx) GetAuthList() []SetCodeAuthorization { return tx.AuthList }

// EncodePayload writes the RLP payload following the type byte.
func (tx *SetCodeTx) EncodePayload(w *bytes.Buffer) error {
	return rlp.Encode(w, tx)
}

// DecodePayload parses the RLP payload following the type byte.
func (tx *SetCodeTx) DecodePayload(b []byte) error {
	return rlp.DecodeBytes(b, tx)
}

// SigningHash returns the digest the sender signs.
func (tx *SetCodeTx) SigningHash(chainID *big.Int) common.Hash {
	return prefixedRlpHash(SetCodeTxType, []any{
		chainID,
		tx.Nonce,
		tx.GasTipCap,
		tx.GasFeeCap,
		tx.Gas,
		tx.To,
		tx.Value,
		tx.Data,
		tx.AccessList,
		tx.AuthList,
	})
}

// GetSignatureValues returns the signature values as encoded.
func (tx *SetCodeTx) GetSignatureValues() (v, r, s *big.Int) {
	return u256ToBig(tx.V), u256ToBig(tx.R), u256ToBig(tx.S)
}

// SetSignatureValues stores a signature. Typed transactions encode v as the raw
// y-parity bit.
func (tx *SetCodeTx) SetSignatureValues(chainID, v, r, s *big.Int) {
	tx.ChainID = uint256.MustFromBig(bigOrZero(chainID))
	tx.V = uint256.MustFromBig(bigOrZero(v))
	tx.R = uint256.MustFromBig(bigOrZero(r))
	tx.S = uint256.MustFromBig(bigOrZero(s))
}

// SignAuthorization signs an EIP-7702 authorization tuple.
func SignAuthorization(auth SetCodeAuthorization, key *ecdsa.PrivateKey) (SetCodeAuthorization, error) {
	return types.SignSetCode(key, auth)
}

// ParseDelegation extracts the delegate address from an EIP-7702 delegation
// indicator, reporting whether the code is a delegation indicator at all.
func ParseDelegation(code []byte) (common.Address, bool) {
	return types.ParseDelegation(code)
}

// AddressToDelegation builds the EIP-7702 delegation indicator code for an address.
func AddressToDelegation(addr common.Address) []byte {
	return types.AddressToDelegation(addr)
}
