// Package txtypes implements spamoor's Ethereum transaction, receipt and block types.
//
// Transaction types are registered at runtime through RegisterTxType, so support for a
// new EIP is a new file in this package or in a plugin. Encoding, signing, sender
// recovery and JSON-RPC decoding are implemented here rather than taken from
// go-ethereum, whose TxData interface cannot be implemented outside its own package.
//
// go-ethereum is still used for crypto primitives, the RLP codec, common value types
// and the plain data structs aliased below.
package txtypes

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
)

// Transaction types implemented by this package.
const (
	LegacyTxType     = 0x00
	AccessListTxType = 0x01
	DynamicFeeTxType = 0x02
	BlobTxType       = 0x03
	SetCodeTxType    = 0x04
)

// Receipt status codes.
const (
	ReceiptStatusFailed     = uint64(0)
	ReceiptStatusSuccessful = uint64(1)
)

// Blob sidecar versions (EIP-4844 / EIP-7594).
const (
	BlobSidecarVersion0 = types.BlobSidecarVersion0
	BlobSidecarVersion1 = types.BlobSidecarVersion1
)

// Plain data structs re-exported from go-ethereum. They carry no transaction-type
// coupling, and aliasing keeps abigen bindings and event parsing source compatible.
type (
	// Log is an EVM log record as returned in receipts.
	Log = types.Log

	// AccessList is an EIP-2930 access list.
	AccessList = types.AccessList

	// AccessTuple is a single EIP-2930 access list entry.
	AccessTuple = types.AccessTuple

	// BlobSidecar carries the blobs, commitments and proofs of a blob transaction.
	// It is wire-only data and is not covered by the transaction hash.
	BlobSidecar = types.BlobTxSidecar

	// SetCodeAuthorization is an EIP-7702 authorization tuple.
	SetCodeAuthorization = types.SetCodeAuthorization
)

// NewBlobSidecar builds a blob sidecar from its components.
func NewBlobSidecar(version byte, blobs []kzg4844.Blob, commitments []kzg4844.Commitment, proofs []kzg4844.Proof) *BlobSidecar {
	return types.NewBlobTxSidecar(version, blobs, commitments, proofs)
}
