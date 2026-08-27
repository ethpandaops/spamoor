package txtypes

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Header carries the block fields spamoor uses. It is a subset of the consensus
// header, decoded permissively so unknown fields from a new fork cannot break block
// processing.
type Header struct {
	Hash       common.Hash
	ParentHash common.Hash
	Number     uint64
	Timestamp  uint64
	GasLimit   uint64
	GasUsed    uint64
	BaseFee    *big.Int

	// BlobGasUsed and ExcessBlobGas are EIP-4844 fields, zero before Cancun.
	BlobGasUsed   uint64
	ExcessBlobGas uint64

	// StateGasUsed is the EIP-8037 state gas dimension, zero on chains that do not
	// report it.
	StateGasUsed uint64
}

// Block is a block header together with its transactions.
type Block struct {
	Header

	Transactions []*Transaction
}

// NumberU64 returns the block number.
func (b *Block) NumberU64() uint64 { return b.Number }

// jsonHeader is the JSON-RPC representation of the header fields spamoor reads.
type jsonHeader struct {
	Hash          *common.Hash    `json:"hash"`
	ParentHash    *common.Hash    `json:"parentHash"`
	Number        *hexutil.Big    `json:"number"`
	Timestamp     *hexutil.Uint64 `json:"timestamp"`
	GasLimit      *hexutil.Uint64 `json:"gasLimit"`
	GasUsed       *hexutil.Uint64 `json:"gasUsed"`
	BaseFee       *hexutil.Big    `json:"baseFeePerGas"`
	BlobGasUsed   *hexutil.Uint64 `json:"blobGasUsed"`
	ExcessBlobGas *hexutil.Uint64 `json:"excessBlobGas"`
	StateGasUsed  *hexutil.Uint64 `json:"stateGasUsed"`
}

// UnmarshalJSON decodes a header from a JSON-RPC block response.
func (h *Header) UnmarshalJSON(input []byte) error {
	var dec jsonHeader
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}

	if dec.Number == nil {
		return errors.New("block is missing number")
	}

	h.Number = (*big.Int)(dec.Number).Uint64()

	if dec.Hash != nil {
		h.Hash = *dec.Hash
	}

	if dec.ParentHash != nil {
		h.ParentHash = *dec.ParentHash
	}

	if dec.Timestamp != nil {
		h.Timestamp = uint64(*dec.Timestamp)
	}

	if dec.GasLimit != nil {
		h.GasLimit = uint64(*dec.GasLimit)
	}

	if dec.GasUsed != nil {
		h.GasUsed = uint64(*dec.GasUsed)
	}

	if dec.BaseFee != nil {
		h.BaseFee = (*big.Int)(dec.BaseFee)
	} else {
		h.BaseFee = new(big.Int)
	}

	if dec.BlobGasUsed != nil {
		h.BlobGasUsed = uint64(*dec.BlobGasUsed)
	}

	if dec.ExcessBlobGas != nil {
		h.ExcessBlobGas = uint64(*dec.ExcessBlobGas)
	}

	if dec.StateGasUsed != nil {
		h.StateGasUsed = uint64(*dec.StateGasUsed)
	}

	return nil
}

// UnmarshalJSON decodes a block including its full transaction objects. Transactions
// of unregistered types are kept as UnknownTx rather than dropped, so transaction
// indices stay aligned with the block's receipts.
func (b *Block) UnmarshalJSON(input []byte) error {
	if err := b.Header.UnmarshalJSON(input); err != nil {
		return err
	}

	var body struct {
		Transactions []json.RawMessage `json:"transactions"`
	}

	if err := json.Unmarshal(input, &body); err != nil {
		return err
	}

	b.Transactions = make([]*Transaction, 0, len(body.Transactions))

	for _, rawTx := range body.Transactions {
		tx, err := UnmarshalJSONTx(rawTx)
		if err != nil {
			return err
		}

		b.Transactions = append(b.Transactions, tx)
	}

	return nil
}
