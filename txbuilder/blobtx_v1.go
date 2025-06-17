//go:build with_blob_v1

package txbuilder

import (
	"bytes"

	gokzg4844 "github.com/crate-crypto/go-eth-kzg"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/holiman/uint256"
)

// init initializes the blob v1 function pointers when this package is built with blob v1 support.
// This enables the experimental blob v1 functionality in the main blobtx.go file by setting
// the function pointers that were declared as variables.
func init() {
	blobV1Marshaller = marshalBlobV1Tx
	blobV1GenerateCellProof = generateCellProofs
}

// blobV1TxWithBlobs represents a blob v1 transaction with additional cell proof data.
// It extends the standard blob transaction format with cell-level proofs for enhanced
// verification capabilities. The Version field indicates the blob transaction version (1).
type blobV1TxWithBlobs struct {
	BlobTx      *types.BlobTx
	Version     uint8
	Blobs       []kzg4844.Blob
	Commitments []kzg4844.Commitment
	CellProofs  []kzg4844.Proof
}

// generateCellProofs generates cell-level KZG proofs for all blobs in a transaction sidecar.
// It computes cells for each blob and creates proofs for each cell, enabling more granular
// verification of blob data. Returns a slice of proofs covering all cells across all blobs.
func generateCellProofs(sidecar *types.BlobTxSidecar) ([]kzg4844.Proof, error) {
	cellProofs := make([]kzg4844.Proof, 0, len(sidecar.Blobs)*gokzg4844.CellsPerExtBlob)
	for _, blobs := range sidecar.Blobs {
		cellProof, err := kzg4844.ComputeCells(&blobs)
		if err != nil {
			return nil, err
		}
		cellProofs = append(cellProofs, cellProof...)
	}
	return cellProofs, nil
}

// marshalBlobV1Tx marshals a transaction into the experimental blob v1 format with cell proofs.
// It extracts the blob transaction data, signature values, and combines them with the provided
// cell proofs into a blob v1 transaction structure. The result is RLP-encoded with the
// transaction type prefix. For non-blob transactions, it falls back to standard marshaling.
func marshalBlobV1Tx(tx *types.Transaction, cellProofs []kzg4844.Proof) ([]byte, error) {
	blobTxSidecar := tx.BlobTxSidecar()
	if tx.Type() != types.BlobTxType || len(blobTxSidecar.Blobs) == 0 {
		return tx.MarshalBinary()
	}

	blobTx := &types.BlobTx{
		ChainID:    uint256.MustFromBig(tx.ChainId()),
		Nonce:      tx.Nonce(),
		GasTipCap:  uint256.MustFromBig(tx.GasTipCap()),
		GasFeeCap:  uint256.MustFromBig(tx.GasFeeCap()),
		BlobFeeCap: uint256.MustFromBig(tx.BlobGasFeeCap()),
		Gas:        tx.Gas(),
		To:         *tx.To(),
		Value:      uint256.MustFromBig(tx.Value()),
		Data:       tx.Data(),
		AccessList: tx.AccessList(),
		BlobHashes: tx.BlobHashes(),
	}

	v, r, s := tx.RawSignatureValues()
	blobTx.R = uint256.MustFromBig(r)
	blobTx.S = uint256.MustFromBig(s)
	blobTx.V = uint256.MustFromBig(v)

	blobV1Tx := blobV1TxWithBlobs{
		BlobTx:      blobTx,
		Version:     1,
		Blobs:       blobTxSidecar.Blobs,
		Commitments: blobTxSidecar.Commitments,
		CellProofs:  cellProofs,
	}

	var buf bytes.Buffer
	buf.WriteByte(tx.Type())
	err := rlp.Encode(&buf, blobV1Tx)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
