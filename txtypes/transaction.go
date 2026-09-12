package txtypes

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	// ErrTxTypeNotSupported is returned when a transaction type has no registered
	// implementation.
	ErrTxTypeNotSupported = errors.New("transaction type not supported")

	// ErrInvalidSig is returned when a transaction carries an unrecoverable signature.
	ErrInvalidSig = errors.New("invalid transaction v, r, s values")

	errShortTypedTx = errors.New("typed transaction too short")
	errEmptyTypedTx = errors.New("empty typed transaction bytes")
)

// TxData is the consensus content of a single transaction.
//
// Every method is exported so transaction types can be implemented outside this
// package and registered with RegisterTxType. The accessors cover the shape common to
// all types; anything beyond it is expressed through the optional interfaces below.
// Types with no single natural value for an accessor return a derived one.
type TxData interface {
	// TxType returns the EIP-2718 type byte.
	TxType() byte

	// CopyTx returns a deep copy with all fields initialized.
	CopyTx() TxData

	// Get-prefixed because implementations carry these names as exported fields.
	GetChainID() *big.Int
	GetNonce() uint64
	GetGas() uint64
	GetGasPrice() *big.Int
	GetGasFeeCap() *big.Int
	GetGasTipCap() *big.Int
	GetTo() *common.Address
	GetValue() *big.Int
	GetData() []byte

	// EncodePayload writes the canonical EIP-2718 payload, i.e. everything after the
	// type byte. Legacy transactions write their plain RLP list.
	EncodePayload(w *bytes.Buffer) error

	// DecodePayload parses a canonical payload produced by EncodePayload. Types whose
	// network encoding differs (blob sidecars, EIP-7594 wrappers) should accept both
	// forms here so that UnmarshalBinary works for either.
	DecodePayload(b []byte) error
}

// ECDSASignedTx is implemented by transaction types authenticated by a single
// secp256k1 signature carried in v/r/s fields (types 0x00 - 0x04).
type ECDSASignedTx interface {
	// SigningHash returns the digest the sender signs.
	SigningHash(chainID *big.Int) common.Hash

	// GetSignatureValues returns the signature values as encoded.
	GetSignatureValues() (v, r, s *big.Int)

	// SetSignatureValues stores a signature on the transaction.
	SetSignatureValues(chainID, v, r, s *big.Int)
}

// ExplicitSenderTx is implemented by transaction types that carry their sender as an
// explicit field rather than deriving it from a signature. Their signature material is
// internal to the payload, so they sign themselves.
type ExplicitSenderTx interface {
	// GetSender returns the transaction's declared sender.
	GetSender() common.Address

	// SignPayload signs the transaction's internal signature material with key.
	SignPayload(chainID *big.Int, key *ecdsa.PrivateKey) error
}

// NetworkEncodedTx is implemented by transaction types whose wire encoding differs
// from the canonical (block/hash) encoding, such as blob transactions carrying a
// sidecar.
type NetworkEncodedTx interface {
	// EncodeNetworkPayload writes the payload as it is sent to eth_sendRawTransaction
	// and gossiped in PooledTransactions.
	EncodeNetworkPayload(w *bytes.Buffer) error
}

// AccessListTxData is implemented by transaction types carrying an EIP-2930 access list.
type AccessListTxData interface {
	GetAccessList() AccessList
}

// BlobTxData is implemented by transaction types that can carry EIP-4844 blobs.
type BlobTxData interface {
	GetBlobHashes() []common.Hash
	GetBlobGasFeeCap() *big.Int
	GetBlobSidecar() *BlobSidecar
	SetBlobSidecar(sidecar *BlobSidecar)
}

// AuthListTxData is implemented by transaction types carrying an EIP-7702
// authorization list.
type AuthListTxData interface {
	GetAuthList() []SetCodeAuthorization
}

// StateGasTxData is implemented by transaction types that declare an EIP-8037 state
// gas budget separately from their execution gas budget.
type StateGasTxData interface {
	GetStateGas() uint64
}

// IndependentNonceTx is implemented by transaction types whose nonce does not address the
// sender's account nonce sequence, such as a frame transaction on a non-zero EIP-8250 key
// set. Anything tracking transactions by account nonce has to ask before assuming.
type IndependentNonceTx interface {
	// UsesAccountNonce reports whether GetNonce addresses the sender's account nonce.
	UsesAccountNonce() bool
}

// txRegistry maps EIP-2718 type bytes to constructors. Guarded by a mutex because
// plugins may register types after startup.
var (
	txRegistryMutex sync.RWMutex
	txRegistry      = make(map[byte]func() TxData, 8)
)

// RegisterTxType makes a transaction type known to the decoder, normally from the
// type's init function. Registering the same type byte twice panics.
func RegisterTxType(txType byte, newFn func() TxData) {
	txRegistryMutex.Lock()
	defer txRegistryMutex.Unlock()

	if _, exists := txRegistry[txType]; exists {
		panic(fmt.Sprintf("txtypes: transaction type 0x%02x already registered", txType))
	}

	txRegistry[txType] = newFn
}

// IsTxTypeSupported reports whether a decoder is registered for txType.
func IsTxTypeSupported(txType byte) bool {
	txRegistryMutex.RLock()
	defer txRegistryMutex.RUnlock()

	_, exists := txRegistry[txType]

	return exists
}

// RegisteredTxTypes returns the type bytes with a registered decoder.
func RegisteredTxTypes() []byte {
	txRegistryMutex.RLock()
	defer txRegistryMutex.RUnlock()

	txTypes := make([]byte, 0, len(txRegistry))
	for txType := range txRegistry {
		txTypes = append(txTypes, txType)
	}

	return txTypes
}

// newTxData constructs an empty TxData for the given type.
func newTxData(txType byte) (TxData, error) {
	txRegistryMutex.RLock()
	newFn, exists := txRegistry[txType]
	txRegistryMutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("%w: 0x%02x", ErrTxTypeNotSupported, txType)
	}

	return newFn(), nil
}

// Transaction wraps a TxData with the caches and derived values the engine needs.
//
// The zero value is not usable; construct with NewTx or one of the decoders.
type Transaction struct {
	inner TxData

	hash atomic.Pointer[common.Hash]
	from atomic.Pointer[common.Address]
	size atomic.Uint64
}

// NewTx creates a transaction from its consensus content.
func NewTx(inner TxData) *Transaction {
	return &Transaction{inner: inner}
}

// Inner returns the underlying type-specific data. Callers type-assert it to reach
// fields that are not part of the common accessor set.
func (tx *Transaction) Inner() TxData {
	return tx.inner
}

// Type returns the EIP-2718 transaction type.
func (tx *Transaction) Type() uint8 {
	return tx.inner.TxType()
}

// ChainId returns the transaction's chain id. Named to match go-ethereum so that
// call sites migrating from core/types do not need to change.
func (tx *Transaction) ChainId() *big.Int {
	return tx.inner.GetChainID()
}

// Nonce returns the sender account nonce.
func (tx *Transaction) Nonce() uint64 { return tx.inner.GetNonce() }

// Gas returns the gas budget the sender is charged against.
func (tx *Transaction) Gas() uint64 { return tx.inner.GetGas() }

// GasPrice returns the gas price for legacy transactions, or the fee cap otherwise.
func (tx *Transaction) GasPrice() *big.Int { return tx.inner.GetGasPrice() }

// GasFeeCap returns the maximum fee per gas.
func (tx *Transaction) GasFeeCap() *big.Int { return tx.inner.GetGasFeeCap() }

// GasTipCap returns the maximum priority fee per gas.
func (tx *Transaction) GasTipCap() *big.Int { return tx.inner.GetGasTipCap() }

// To returns the recipient, or nil for contract creation.
func (tx *Transaction) To() *common.Address { return tx.inner.GetTo() }

// Value returns the transferred amount in wei.
func (tx *Transaction) Value() *big.Int { return tx.inner.GetValue() }

// Data returns the transaction calldata.
func (tx *Transaction) Data() []byte { return tx.inner.GetData() }

// AccessList returns the EIP-2930 access list, or nil if the type has none.
func (tx *Transaction) AccessList() AccessList {
	if inner, ok := tx.inner.(AccessListTxData); ok {
		return inner.GetAccessList()
	}

	return nil
}

// BlobHashes returns the EIP-4844 blob versioned hashes, or nil if the type has none.
func (tx *Transaction) BlobHashes() []common.Hash {
	if inner, ok := tx.inner.(BlobTxData); ok {
		return inner.GetBlobHashes()
	}

	return nil
}

// BlobGasFeeCap returns the maximum fee per blob gas, or nil if the type has none.
func (tx *Transaction) BlobGasFeeCap() *big.Int {
	if inner, ok := tx.inner.(BlobTxData); ok {
		return inner.GetBlobGasFeeCap()
	}

	return nil
}

// BlobTxSidecar returns the attached blob sidecar, or nil when the transaction
// carries no blobs or was decoded from its canonical (sidecar-less) form.
func (tx *Transaction) BlobTxSidecar() *BlobSidecar {
	if inner, ok := tx.inner.(BlobTxData); ok {
		return inner.GetBlobSidecar()
	}

	return nil
}

// SetBlobTxSidecar attaches a blob sidecar. It resets the size cache but not the hash
// cache, as the sidecar is not covered by the transaction hash.
func (tx *Transaction) SetBlobTxSidecar(sidecar *BlobSidecar) {
	if inner, ok := tx.inner.(BlobTxData); ok {
		inner.SetBlobSidecar(sidecar)
		tx.size.Store(0)
	}
}

// AuthList returns the EIP-7702 authorization list, or nil if the type has none.
func (tx *Transaction) AuthList() []SetCodeAuthorization {
	if inner, ok := tx.inner.(AuthListTxData); ok {
		return inner.GetAuthList()
	}

	return nil
}

// StateGas returns the EIP-8037 state gas budget declared by the transaction, or 0
// for types that do not declare one separately.
func (tx *Transaction) StateGas() uint64 {
	if inner, ok := tx.inner.(StateGasTxData); ok {
		return inner.GetStateGas()
	}

	return 0
}

// UsesAccountNonce reports whether the transaction is sequenced by its sender's account
// nonce. When it is false, Nonce returns a sequence number from an unrelated domain and
// comparing it against the account nonce is meaningless.
func (tx *Transaction) UsesAccountNonce() bool {
	if inner, ok := tx.inner.(IndependentNonceTx); ok {
		return inner.UsesAccountNonce()
	}

	return true
}

// RawSignatureValues returns the signature values of an ECDSA-signed transaction.
// Types that authenticate differently return nil values.
func (tx *Transaction) RawSignatureValues() (v, r, s *big.Int) {
	if inner, ok := tx.inner.(ECDSASignedTx); ok {
		return inner.GetSignatureValues()
	}

	return nil, nil, nil
}

// Cost returns the maximum wei the transaction can consume: gas budget times fee cap
// plus value. Blob fees are charged from a separate market and are not included.
func (tx *Transaction) Cost() *big.Int {
	total := new(big.Int).Mul(tx.GasFeeCap(), new(big.Int).SetUint64(tx.Gas()))

	return total.Add(total, tx.Value())
}

// Hash returns the transaction hash, which covers the canonical encoding only. It is
// computed on first use and cached.
func (tx *Transaction) Hash() common.Hash {
	if hash := tx.hash.Load(); hash != nil {
		return *hash
	}

	encoded, err := tx.MarshalBinary()
	if err != nil {
		return common.Hash{}
	}

	hash := crypto.Keccak256Hash(encoded)
	tx.hash.Store(&hash)

	return hash
}

// Size returns the length of the network encoding in bytes.
func (tx *Transaction) Size() uint64 {
	if size := tx.size.Load(); size > 0 {
		return size
	}

	encoded, err := tx.MarshalNetwork()
	if err != nil {
		return 0
	}

	size := uint64(len(encoded))
	tx.size.Store(size)

	return size
}

// SetFrom records the transaction's sender, letting decoders supply the sender
// reported by the node instead of recovering it.
func (tx *Transaction) SetFrom(from common.Address) {
	tx.from.Store(&from)
}

// From returns the transaction sender: one supplied by SetFrom, the explicit sender
// field of types that carry one, or secp256k1 recovery from the signature. Cached.
func (tx *Transaction) From(chainID *big.Int) (common.Address, error) {
	if from := tx.from.Load(); from != nil {
		return *from, nil
	}

	if inner, ok := tx.inner.(ExplicitSenderTx); ok {
		from := inner.GetSender()
		tx.from.Store(&from)

		return from, nil
	}

	from, err := recoverSender(tx, chainID)
	if err != nil {
		return common.Address{}, err
	}

	tx.from.Store(&from)

	return from, nil
}

// Copy returns a deep copy of the transaction without its caches.
func (tx *Transaction) Copy() *Transaction {
	return NewTx(tx.inner.CopyTx())
}

// MarshalBinary returns the canonical EIP-2718 encoding, covered by the transaction
// hash and included in blocks. Blob sidecars are excluded; use MarshalNetwork to
// submit or gossip a transaction.
func (tx *Transaction) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	if tx.Type() != LegacyTxType {
		buf.WriteByte(tx.Type())
	}

	if err := tx.inner.EncodePayload(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// MarshalNetwork returns the encoding used for eth_sendRawTransaction and transaction
// gossip. For types with no separate network form it is identical to MarshalBinary.
func (tx *Transaction) MarshalNetwork() ([]byte, error) {
	inner, ok := tx.inner.(NetworkEncodedTx)
	if !ok {
		return tx.MarshalBinary()
	}

	var buf bytes.Buffer

	if tx.Type() != LegacyTxType {
		buf.WriteByte(tx.Type())
	}

	if err := inner.EncodeNetworkPayload(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary parses a transaction from its canonical or network encoding.
func (tx *Transaction) UnmarshalBinary(b []byte) error {
	if len(b) == 0 {
		return errEmptyTypedTx
	}

	// A leading byte >= 0xc0 is an RLP list header, i.e. a legacy transaction.
	if b[0] >= 0xc0 {
		inner, err := newTxData(LegacyTxType)
		if err != nil {
			return err
		}

		if err := inner.DecodePayload(b); err != nil {
			return err
		}

		tx.setDecoded(inner)

		return nil
	}

	if len(b) <= 1 {
		return errShortTypedTx
	}

	inner, err := newTxData(b[0])
	if err != nil {
		return err
	}

	if err := inner.DecodePayload(b[1:]); err != nil {
		return err
	}

	tx.setDecoded(inner)

	return nil
}

// setDecoded installs freshly decoded content and drops any cached values.
func (tx *Transaction) setDecoded(inner TxData) {
	tx.inner = inner
	tx.hash.Store(nil)
	tx.from.Store(nil)
	tx.size.Store(0)
}

// DecodeTx parses a transaction from its canonical or network encoding.
func DecodeTx(b []byte) (*Transaction, error) {
	tx := &Transaction{}
	if err := tx.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	return tx, nil
}
