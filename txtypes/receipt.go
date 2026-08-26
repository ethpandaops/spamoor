package txtypes

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Receipt is the result of executing a transaction. Field names and types mirror
// go-ethereum's types.Receipt. Type-specific content lives in Extra, decoded by the
// receipt decoder registered for the transaction type.
type Receipt struct {
	Type              uint8
	Status            uint64
	CumulativeGasUsed uint64
	Logs              []*Log
	TxHash            common.Hash
	ContractAddress   common.Address
	GasUsed           uint64
	EffectiveGasPrice *big.Int
	BlobGasUsed       uint64
	BlobGasPrice      *big.Int
	BlockHash         common.Hash
	BlockNumber       *big.Int
	TransactionIndex  uint

	// Extra carries type-specific receipt content, or nil when the type has none.
	Extra ReceiptExtra
}

// ReceiptExtra is type-specific receipt content.
type ReceiptExtra interface {
	// ReceiptTxType returns the transaction type this content belongs to.
	ReceiptTxType() byte
}

// ReceiptDecoder parses the type-specific part of a receipt from the raw JSON-RPC
// response. It may also fill in generic fields the type does not report directly,
// such as deriving a status.
type ReceiptDecoder func(receipt *Receipt, raw json.RawMessage) error

var (
	receiptRegistryMutex sync.RWMutex
	receiptRegistry      = make(map[byte]ReceiptDecoder, 4)
)

// RegisterReceiptDecoder installs a receipt decoder for a transaction type.
func RegisterReceiptDecoder(txType byte, decoder ReceiptDecoder) {
	receiptRegistryMutex.Lock()
	defer receiptRegistryMutex.Unlock()

	if _, exists := receiptRegistry[txType]; exists {
		panic(fmt.Sprintf("txtypes: receipt decoder for type 0x%02x already registered", txType))
	}

	receiptRegistry[txType] = decoder
}

// receiptDecoderFor returns the decoder for a transaction type, or nil.
func receiptDecoderFor(txType byte) ReceiptDecoder {
	receiptRegistryMutex.RLock()
	defer receiptRegistryMutex.RUnlock()

	return receiptRegistry[txType]
}

// jsonReceipt is the JSON-RPC representation of a receipt. Every field is optional:
// clients differ in what they report and unknown transaction types may omit fields.
// Only the transaction hash is required.
type jsonReceipt struct {
	Type              *hexutil.Uint64 `json:"type"`
	Root              *hexutil.Bytes  `json:"root"`
	Status            *hexutil.Uint64 `json:"status"`
	CumulativeGasUsed *hexutil.Uint64 `json:"cumulativeGasUsed"`
	Logs              []*Log          `json:"logs"`
	TxHash            *common.Hash    `json:"transactionHash"`
	ContractAddress   *common.Address `json:"contractAddress"`
	GasUsed           *hexutil.Uint64 `json:"gasUsed"`
	EffectiveGasPrice *hexutil.Big    `json:"effectiveGasPrice"`
	BlobGasUsed       *hexutil.Uint64 `json:"blobGasUsed"`
	BlobGasPrice      *hexutil.Big    `json:"blobGasPrice"`
	BlockHash         *common.Hash    `json:"blockHash"`
	BlockNumber       *hexutil.Big    `json:"blockNumber"`
	TransactionIndex  *hexutil.Uint   `json:"transactionIndex"`
}

// UnmarshalJSON decodes a receipt from a JSON-RPC response.
func (r *Receipt) UnmarshalJSON(input []byte) error {
	var dec jsonReceipt
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}

	if dec.TxHash == nil {
		return errors.New("receipt is missing transactionHash")
	}

	r.TxHash = *dec.TxHash
	r.Logs = dec.Logs

	if dec.Type != nil {
		r.Type = uint8(*dec.Type)
	}

	if dec.Status != nil {
		r.Status = uint64(*dec.Status)
	}

	if dec.CumulativeGasUsed != nil {
		r.CumulativeGasUsed = uint64(*dec.CumulativeGasUsed)
	}

	if dec.ContractAddress != nil {
		r.ContractAddress = *dec.ContractAddress
	}

	if dec.GasUsed != nil {
		r.GasUsed = uint64(*dec.GasUsed)
	}

	if dec.EffectiveGasPrice != nil {
		r.EffectiveGasPrice = (*big.Int)(dec.EffectiveGasPrice)
	}

	if dec.BlobGasUsed != nil {
		r.BlobGasUsed = uint64(*dec.BlobGasUsed)
	}

	if dec.BlobGasPrice != nil {
		r.BlobGasPrice = (*big.Int)(dec.BlobGasPrice)
	}

	if dec.BlockHash != nil {
		r.BlockHash = *dec.BlockHash
	}

	if dec.BlockNumber != nil {
		r.BlockNumber = (*big.Int)(dec.BlockNumber)
	}

	if dec.TransactionIndex != nil {
		r.TransactionIndex = uint(*dec.TransactionIndex)
	}

	if decoder := receiptDecoderFor(r.Type); decoder != nil {
		if err := decoder(r, input); err != nil {
			return fmt.Errorf("failed decoding type 0x%02x receipt: %w", r.Type, err)
		}
	}

	return nil
}

// Successful reports whether the transaction executed successfully.
func (r *Receipt) Successful() bool {
	return r != nil && r.Status == ReceiptStatusSuccessful
}
