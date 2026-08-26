package txtypes

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/holiman/uint256"
)

func init() {
	RegisterTxType(BlobTxType, func() TxData { return &BlobTx{} })
}

// BlobTx is an EIP-4844 blob transaction.
//
// The sidecar is wire-only data: it is excluded from the canonical encoding and from
// the transaction hash, and travels only in the network encoding written by
// EncodeNetworkPayload.
type BlobTx struct {
	ChainID    *uint256.Int
	Nonce      uint64
	GasTipCap  *uint256.Int // a.k.a. maxPriorityFeePerGas
	GasFeeCap  *uint256.Int // a.k.a. maxFeePerGas
	Gas        uint64
	To         common.Address
	Value      *uint256.Int
	Data       []byte
	AccessList AccessList
	BlobFeeCap *uint256.Int // a.k.a. maxFeePerBlobGas
	BlobHashes []common.Hash

	// Sidecar is optional on the wire but must be set to create a transaction for
	// submission.
	Sidecar *BlobSidecar `rlp:"-"`

	V, R, S *uint256.Int // signature values
}

var (
	_ TxData           = (*BlobTx)(nil)
	_ ECDSASignedTx    = (*BlobTx)(nil)
	_ AccessListTxData = (*BlobTx)(nil)
	_ BlobTxData       = (*BlobTx)(nil)
	_ NetworkEncodedTx = (*BlobTx)(nil)
)

// blobTxWithSidecarV0 is the pre-EIP-7594 network encoding: [tx, blobs, commitments,
// proofs].
type blobTxWithSidecarV0 struct {
	Tx          *BlobTx
	Blobs       []kzg4844.Blob
	Commitments []kzg4844.Commitment
	Proofs      []kzg4844.Proof
}

// blobTxWithSidecarV1 is the versioned network encoding introduced by EIP-7594:
// [tx, version, blobs, commitments, cell_proofs].
type blobTxWithSidecarV1 struct {
	Tx          *BlobTx
	Version     byte
	Blobs       []kzg4844.Blob
	Commitments []kzg4844.Commitment
	Proofs      []kzg4844.Proof
}

// TxType returns the EIP-2718 type byte.
func (tx *BlobTx) TxType() byte { return BlobTxType }

// CopyTx returns a deep copy with all fields initialized. The sidecar is shared
// rather than cloned, matching go-ethereum: sidecars are large and treated as
// immutable once attached.
func (tx *BlobTx) CopyTx() TxData {
	cpy := &BlobTx{
		Nonce:      tx.Nonce,
		To:         tx.To,
		Data:       common.CopyBytes(tx.Data),
		Gas:        tx.Gas,
		AccessList: make(AccessList, len(tx.AccessList)),
		BlobHashes: make([]common.Hash, len(tx.BlobHashes)),
		Sidecar:    tx.Sidecar,
		ChainID:    new(uint256.Int),
		GasTipCap:  new(uint256.Int),
		GasFeeCap:  new(uint256.Int),
		Value:      new(uint256.Int),
		BlobFeeCap: new(uint256.Int),
		V:          new(uint256.Int),
		R:          new(uint256.Int),
		S:          new(uint256.Int),
	}

	copy(cpy.AccessList, tx.AccessList)
	copy(cpy.BlobHashes, tx.BlobHashes)
	setU256(cpy.ChainID, tx.ChainID)
	setU256(cpy.GasTipCap, tx.GasTipCap)
	setU256(cpy.GasFeeCap, tx.GasFeeCap)
	setU256(cpy.Value, tx.Value)
	setU256(cpy.BlobFeeCap, tx.BlobFeeCap)
	setU256(cpy.V, tx.V)
	setU256(cpy.R, tx.R)
	setU256(cpy.S, tx.S)

	return cpy
}

// GetChainID returns the destination chain id.
func (tx *BlobTx) GetChainID() *big.Int { return u256ToBig(tx.ChainID) }

// GetNonce returns the sender account nonce.
func (tx *BlobTx) GetNonce() uint64 { return tx.Nonce }

// GetGas returns the gas limit.
func (tx *BlobTx) GetGas() uint64 { return tx.Gas }

// GetGasPrice returns the fee cap; blob transactions have no single gas price.
func (tx *BlobTx) GetGasPrice() *big.Int { return u256ToBig(tx.GasFeeCap) }

// GetGasFeeCap returns the maximum fee per gas.
func (tx *BlobTx) GetGasFeeCap() *big.Int { return u256ToBig(tx.GasFeeCap) }

// GetGasTipCap returns the maximum priority fee per gas.
func (tx *BlobTx) GetGasTipCap() *big.Int { return u256ToBig(tx.GasTipCap) }

// GetTo returns the recipient. Blob transactions cannot create contracts, so this is
// never nil.
func (tx *BlobTx) GetTo() *common.Address {
	to := tx.To

	return &to
}

// GetValue returns the transferred amount in wei.
func (tx *BlobTx) GetValue() *big.Int { return u256ToBig(tx.Value) }

// GetData returns the transaction calldata.
func (tx *BlobTx) GetData() []byte { return tx.Data }

// GetAccessList returns the EIP-2930 access list.
func (tx *BlobTx) GetAccessList() AccessList { return tx.AccessList }

// GetBlobHashes returns the blob versioned hashes.
func (tx *BlobTx) GetBlobHashes() []common.Hash { return tx.BlobHashes }

// GetBlobGasFeeCap returns the maximum fee per blob gas.
func (tx *BlobTx) GetBlobGasFeeCap() *big.Int { return u256ToBig(tx.BlobFeeCap) }

// GetBlobSidecar returns the attached sidecar, or nil.
func (tx *BlobTx) GetBlobSidecar() *BlobSidecar { return tx.Sidecar }

// SetBlobSidecar attaches a sidecar.
func (tx *BlobTx) SetBlobSidecar(sidecar *BlobSidecar) { tx.Sidecar = sidecar }

// EncodePayload writes the canonical payload, without the sidecar. This is the form
// covered by the transaction hash and included in block bodies.
func (tx *BlobTx) EncodePayload(w *bytes.Buffer) error {
	return rlp.Encode(w, tx)
}

// EncodeNetworkPayload writes the payload including the blob sidecar, in the wire
// form matching the sidecar's version.
func (tx *BlobTx) EncodeNetworkPayload(w *bytes.Buffer) error {
	if tx.Sidecar == nil {
		return rlp.Encode(w, tx)
	}

	switch tx.Sidecar.Version {
	case BlobSidecarVersion0:
		return rlp.Encode(w, &blobTxWithSidecarV0{
			Tx:          tx,
			Blobs:       tx.Sidecar.Blobs,
			Commitments: tx.Sidecar.Commitments,
			Proofs:      tx.Sidecar.Proofs,
		})
	case BlobSidecarVersion1:
		return rlp.Encode(w, &blobTxWithSidecarV1{
			Tx:          tx,
			Version:     tx.Sidecar.Version,
			Blobs:       tx.Sidecar.Blobs,
			Commitments: tx.Sidecar.Commitments,
			Proofs:      tx.Sidecar.Proofs,
		})
	default:
		return fmt.Errorf("unsupported blob sidecar version %d", tx.Sidecar.Version)
	}
}

// DecodePayload parses either the canonical payload or one of the two network
// encodings.
//
// The three forms are distinguished structurally, as go-ethereum does. The canonical
// form is a list whose first element is the chain id, i.e. a number. The network forms
// wrap that list, so their first element is itself a list; of those, the v0 form has a
// list (the blobs) as its second element while the v1 form has the version byte.
func (tx *BlobTx) DecodePayload(b []byte) error {
	firstElem, _, err := rlp.SplitList(b)
	if err != nil {
		return err
	}

	firstElemKind, _, secondElem, err := rlp.Split(firstElem)
	if err != nil {
		return err
	}

	if firstElemKind != rlp.List {
		// Canonical encoding: no sidecar.
		return rlp.DecodeBytes(b, tx)
	}

	secondElemKind, _, _, err := rlp.Split(secondElem)
	if err != nil {
		return err
	}

	if secondElemKind == rlp.List {
		var payload blobTxWithSidecarV0
		if err := rlp.DecodeBytes(b, &payload); err != nil {
			return err
		}

		if payload.Tx == nil {
			return errors.New("blob transaction wrapper without transaction")
		}

		*tx = *payload.Tx
		tx.Sidecar = &BlobSidecar{
			Version:     BlobSidecarVersion0,
			Blobs:       payload.Blobs,
			Commitments: payload.Commitments,
			Proofs:      payload.Proofs,
		}

		return nil
	}

	var payload blobTxWithSidecarV1
	if err := rlp.DecodeBytes(b, &payload); err != nil {
		return err
	}

	if payload.Tx == nil {
		return errors.New("blob transaction wrapper without transaction")
	}

	if payload.Version != BlobSidecarVersion1 {
		return fmt.Errorf("unsupported blob sidecar version %d", payload.Version)
	}

	*tx = *payload.Tx
	tx.Sidecar = &BlobSidecar{
		Version:     payload.Version,
		Blobs:       payload.Blobs,
		Commitments: payload.Commitments,
		Proofs:      payload.Proofs,
	}

	return nil
}

// SigningHash returns the digest the sender signs.
func (tx *BlobTx) SigningHash(chainID *big.Int) common.Hash {
	return prefixedRlpHash(BlobTxType, []any{
		chainID,
		tx.Nonce,
		tx.GasTipCap,
		tx.GasFeeCap,
		tx.Gas,
		tx.To,
		tx.Value,
		tx.Data,
		tx.AccessList,
		tx.BlobFeeCap,
		tx.BlobHashes,
	})
}

// GetSignatureValues returns the signature values as encoded.
func (tx *BlobTx) GetSignatureValues() (v, r, s *big.Int) {
	return u256ToBig(tx.V), u256ToBig(tx.R), u256ToBig(tx.S)
}

// SetSignatureValues stores a signature. Typed transactions encode v as the raw
// y-parity bit.
func (tx *BlobTx) SetSignatureValues(chainID, v, r, s *big.Int) {
	tx.ChainID = uint256.MustFromBig(bigOrZero(chainID))
	tx.V = uint256.MustFromBig(bigOrZero(v))
	tx.R = uint256.MustFromBig(bigOrZero(r))
	tx.S = uint256.MustFromBig(bigOrZero(s))
}
