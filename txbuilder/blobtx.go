package txbuilder

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	mathRand "math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"

	gokzg4844 "github.com/crate-crypto/go-eth-kzg"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/holiman/uint256"
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

type blobV1TxWithBlobs struct {
	BlobTx      *types.BlobTx
	Version     uint8
	Blobs       []kzg4844.Blob
	Commitments []kzg4844.Commitment
	CellProofs  []kzg4844.Proof
}

func MarshalBlobV1Tx(tx *types.Transaction) ([]byte, error) {
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

	r, s, v := tx.RawSignatureValues()
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

	/*
		fmt.Printf("to: %v\n", blobV1Tx.BlobTx.To)
		fmt.Printf("nonce: %v\n", blobV1Tx.BlobTx.Nonce)
		fmt.Printf("gasTipCap: %v\n", blobV1Tx.BlobTx.GasTipCap)
		fmt.Printf("gasFeeCap: %v\n", blobV1Tx.BlobTx.GasFeeCap)
		fmt.Printf("blobFeeCap: %v\n", blobV1Tx.BlobTx.BlobFeeCap)
		fmt.Printf("gas: %v\n", blobV1Tx.BlobTx.Gas)
		fmt.Printf("value: %v\n", blobV1Tx.BlobTx.Value)
		fmt.Printf("data: %v\n", blobV1Tx.BlobTx.Data)
		fmt.Printf("accessList: %v\n", blobV1Tx.BlobTx.AccessList)

		blobBytes := buf.Bytes()

		os.WriteFile("/home/pk910/Downloads/tx_rlp/2.txt", []byte(common.Bytes2Hex(blobBytes)), 0644)

		testHex, _ := os.ReadFile("/home/pk910/Downloads/tx_rlp/1.txt")
		testBytes := common.FromHex(string(testHex))

		err = rlp.DecodeBytes(testBytes[1:], &blobV1Tx)
		fmt.Printf("test decoded: %v\n", err)

		fmt.Printf("to: %v\n", blobV1Tx.BlobTx.To)
		fmt.Printf("nonce: %v\n", blobV1Tx.BlobTx.Nonce)
		fmt.Printf("gasTipCap: %v\n", blobV1Tx.BlobTx.GasTipCap)
		fmt.Printf("gasFeeCap: %v\n", blobV1Tx.BlobTx.GasFeeCap)
		fmt.Printf("blobFeeCap: %v\n", blobV1Tx.BlobTx.BlobFeeCap)
		fmt.Printf("gas: %v\n", blobV1Tx.BlobTx.Gas)
		fmt.Printf("value: %v\n", blobV1Tx.BlobTx.Value)
		fmt.Printf("data: %v\n", blobV1Tx.BlobTx.Data)
		fmt.Printf("accessList: %v\n", blobV1Tx.BlobTx.AccessList)

		//return blobBytes, nil
	*/
}
