package txbuilder

import (
	"crypto/rand"
	"fmt"
	"io"
	mathRand "math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/params"
)

func BuildBlobTx(txData *TxMetadata, blobRefs [][]string) (*types.BlobTx, error) {
	if txData.To == nil {
		return nil, fmt.Errorf("to cannot be nil for blob transaction")
	}
	tx := types.BlobTx{
		GasTipCap:  txData.GasTipCap,
		GasFeeCap:  txData.GasFeeCap,
		BlobFeeCap: txData.BlobFeeCap,
		Gas:        txData.Gas,
		To:         *txData.To,
		Value:      txData.Value,
		Data:       txData.Data,
		AccessList: txData.AccessList,
		BlobHashes: make([]common.Hash, 0),
		Sidecar: &types.BlobTxSidecar{
			Blobs:       make([]kzg4844.Blob, 0),
			Commitments: make([]kzg4844.Commitment, 0),
			Proofs:      make([]kzg4844.Proof, 0),
		},
	}

	for _, blobRef := range blobRefs {
		err := parseBlobRefs(&tx, blobRef)
		if err != nil {
			return nil, err
		}
	}

	return &tx, nil
}

func ParseBlobRefsBytes(blobRefs []string, tx *types.BlobTx) ([]byte, error) {
	var err error
	var blobBytes []byte

	for _, blobRef := range blobRefs {
		var blobRefBytes []byte
		if strings.HasPrefix(blobRef, "0x") {
			blobRefBytes = common.FromHex(blobRef)
		} else {
			refParts := strings.Split(blobRef, ":")
			switch refParts[0] {
			case "file":
				blobRefBytes, err = os.ReadFile(strings.Join(refParts[1:], ":"))
				if err != nil {
					return nil, err
				}
			case "url":
				blobRefBytes, err = loadUrlRef(strings.Join(refParts[1:], ":"))
				if err != nil {
					return nil, err
				}
			case "repeat":
				if len(refParts) != 3 {
					return nil, fmt.Errorf("invalid repeat ref format: %v", blobRef)
				}
				repeatCount, err := strconv.Atoi(refParts[2])
				if err != nil {
					return nil, fmt.Errorf("invalid repeat count: %v", refParts[2])
				}
				repeatBytes := common.FromHex(refParts[1])
				repeatBytesLen := len(repeatBytes)
				blobRefBytes = make([]byte, repeatCount*repeatBytesLen)
				for i := 0; i < repeatCount; i++ {
					copy(blobRefBytes[(i*repeatBytesLen):], repeatBytes)
				}
			case "random":
				blobLen := -1
				if len(refParts) > 1 {
					var err error
					blobLen, err = strconv.Atoi(refParts[1])
					if err != nil {
						return nil, fmt.Errorf("invalid random count: %v", refParts[1])
					}
				} else {
					blobLen = mathRand.Intn((params.BlobTxFieldElementsPerBlob * (params.BlobTxBytesPerFieldElement - 1)) - len(blobBytes))
				}
				blobRefBytes, err = randomBlobData(blobLen)
				if err != nil {
					return nil, err
				}
			case "copy":
				if tx == nil {
					return nil, fmt.Errorf("copy ref not supported for non blob transactions: %v", blobRef)
				}
				if len(refParts) != 2 {
					return nil, fmt.Errorf("invalid copy ref format: %v", blobRef)
				}
				copyIdx, err := strconv.Atoi(refParts[1])
				if err != nil {
					return nil, fmt.Errorf("invalid copy index: %v", refParts[1])
				}
				if copyIdx >= len(tx.Sidecar.Blobs) {
					return nil, fmt.Errorf("invalid copy index: %v must be smaller than current blob index", refParts[1])
				}
				blobLen := mathRand.Intn((params.BlobTxFieldElementsPerBlob * (params.BlobTxBytesPerFieldElement - 1)) - len(blobBytes))
				if blobLen > len(tx.Sidecar.Blobs[copyIdx]) {
					blobLen = len(tx.Sidecar.Blobs[copyIdx])
				}
				blobRefBytes = tx.Sidecar.Blobs[copyIdx][:blobLen]
			}
		}

		if blobRefBytes == nil {
			return nil, fmt.Errorf("unknown blob ref: %v", blobRef)
		}
		blobBytes = append(blobBytes, blobRefBytes...)
	}

	return blobBytes, nil
}

func parseBlobRefs(tx *types.BlobTx, blobRefs []string) error {
	blobBytes, err := ParseBlobRefsBytes(blobRefs, tx)
	if err != nil {
		return err
	}

	blobCommitment, err := EncodeBlob(blobBytes)
	if err != nil {
		return fmt.Errorf("invalid blob: %w", err)
	}

	tx.BlobHashes = append(tx.BlobHashes, blobCommitment.VersionedHash)
	tx.Sidecar.Blobs = append(tx.Sidecar.Blobs, *blobCommitment.Blob)
	tx.Sidecar.Commitments = append(tx.Sidecar.Commitments, blobCommitment.Commitment)
	tx.Sidecar.Proofs = append(tx.Sidecar.Proofs, blobCommitment.Proof)
	return nil
}

func loadUrlRef(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("received http error: %v", response.Status)
	}
	return io.ReadAll(response.Body)
}

func randomBlobData(size int) ([]byte, error) {
	data := make([]byte, size)
	n, err := rand.Read(data)
	if err != nil {
		return nil, err
	}
	if n != size {
		return nil, fmt.Errorf("could not create random blob data with size %d: %v", size, err)
	}
	return data, nil
}

var blobV1GenerateCellProof func(tx *types.BlobTxSidecar) ([]kzg4844.Proof, error)
var blobV1Marshaller func(tx *types.Transaction, cellProofs []kzg4844.Proof) ([]byte, error)

func GenerateCellProofs(sidecar *types.BlobTxSidecar) ([]kzg4844.Proof, error) {
	if blobV1GenerateCellProof == nil {
		return nil, fmt.Errorf("blob-v1 not supported when using spamoor as library")
	}
	return blobV1GenerateCellProof(sidecar)
}

func MarshalBlobV1Tx(tx *types.Transaction, cellProofs []kzg4844.Proof) ([]byte, error) {
	if blobV1Marshaller == nil {
		return nil, fmt.Errorf("blob-v1 not supported when using spamoor as library")
	}

	return blobV1Marshaller(tx, cellProofs)
}
