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

func init() {
	blobV1Marshaller = marshalBlobV1Tx
}

func marshalBlobV1Tx(tx *types.Transaction) ([]byte, error) {
	blobTxSidecar := tx.BlobTxSidecar()
	if tx.Type() != types.BlobTxType || len(blobTxSidecar.Blobs) == 0 {
		return tx.MarshalBinary()
	}

	// compute the cell proof
	cellProofs := make([]kzg4844.Proof, 0, len(blobTxSidecar.Blobs)*gokzg4844.CellsPerExtBlob)
	for _, blobs := range blobTxSidecar.Blobs {
		cellProof, err := kzg4844.ComputeCells(&blobs)
		if err != nil {
			return nil, err
		}
		cellProofs = append(cellProofs, cellProof...)
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
