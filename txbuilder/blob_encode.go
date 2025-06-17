package txbuilder

import (
	"crypto/sha256"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/params"
)

// BlobCommitment contains all the cryptographic components needed for an EIP-4844 blob.
// It includes the blob data itself, the KZG commitment, proof, and the versioned hash
// that will be included in the transaction. This structure represents a complete
// blob ready for inclusion in a blob transaction.
type BlobCommitment struct {
	Blob          *kzg4844.Blob
	Commitment    kzg4844.Commitment
	Proof         kzg4844.Proof
	VersionedHash common.Hash
}

// encodeBlobData encodes raw byte data into the EIP-4844 blob format.
// It packs the data into 32-byte field elements, using only 31 bytes per field
// (the first byte is always zero to ensure the field element is valid).
// The data is distributed across field elements up to the blob size limit.
func encodeBlobData(data []byte) *kzg4844.Blob {
	blob := kzg4844.Blob{}
	fieldIndex := -1
	for i := 0; i < len(data); i += 31 {
		fieldIndex++
		if fieldIndex == params.BlobTxFieldElementsPerBlob {
			break
		}
		max := i + 31
		if max > len(data) {
			max = len(data)
		}
		copy(blob[fieldIndex*32+1:], data[i:max])
	}
	return &blob
}

// EncodeBlob encodes arbitrary byte data into a complete blob commitment structure.
// It validates the data size against EIP-4844 limits, encodes the data into blob format,
// generates the KZG commitment and proof, and creates the versioned hash.
// Returns an error if the data exceeds the maximum blob size or if cryptographic
// operations fail.
func EncodeBlob(data []byte) (*BlobCommitment, error) {
	dataLen := len(data)
	if dataLen > params.BlobTxFieldElementsPerBlob*(params.BlobTxBytesPerFieldElement-1) {
		return nil, fmt.Errorf("blob data longer than allowed (length: %v, limit: %v)", dataLen, params.BlobTxFieldElementsPerBlob*(params.BlobTxBytesPerFieldElement-1))
	}
	blobCommitment := BlobCommitment{
		Blob: encodeBlobData(data),
	}
	var err error

	// generate blob commitment
	blobCommitment.Commitment, err = kzg4844.BlobToCommitment(blobCommitment.Blob)
	if err != nil {
		return nil, fmt.Errorf("failed generating blob commitment: %w", err)
	}

	// generate blob proof
	blobCommitment.Proof, err = kzg4844.ComputeBlobProof(blobCommitment.Blob, blobCommitment.Commitment)
	if err != nil {
		return nil, fmt.Errorf("failed generating blob proof: %w", err)
	}

	// build versioned hash
	blobCommitment.VersionedHash = sha256.Sum256(blobCommitment.Commitment[:])
	blobCommitment.VersionedHash[0] = 0x01
	return &blobCommitment, nil
}
