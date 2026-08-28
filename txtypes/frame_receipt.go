package txtypes

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Frame receipt status codes. Skipped is new in EIP-8141 and marks a frame that was
// not executed because an earlier frame in its atomic batch failed.
const (
	FrameStatusFailed  = uint64(0)
	FrameStatusSuccess = uint64(1)
	FrameStatusSkipped = uint64(2)
)

func init() {
	RegisterReceiptDecoder(FrameTxType, decodeFrameReceipt)
}

// FrameReceipt is the result of executing one frame.
type FrameReceipt struct {
	Status       uint64
	ExecutionGas uint64
	StateGas     uint64
	Logs         []*Log
}

// Successful reports whether the frame executed successfully.
func (r *FrameReceipt) Successful() bool { return r.Status == FrameStatusSuccess }

// Skipped reports whether the frame was skipped by an atomic batch rollback.
func (r *FrameReceipt) Skipped() bool { return r.Status == FrameStatusSkipped }

// FrameReceiptExtra is the frame-transaction-specific part of a receipt.
//
// EIP-8141 receipts carry no transaction-level status: the only statuses they hold are
// the per-frame ones, and a single status has to be derived from them.
type FrameReceiptExtra struct {
	Payer  common.Address
	Frames []*FrameReceipt
}

// ReceiptTxType returns the transaction type this content belongs to.
func (e *FrameReceiptExtra) ReceiptTxType() byte { return FrameTxType }

// TotalExecutionGas returns the execution gas used across all frames.
func (e *FrameReceiptExtra) TotalExecutionGas() uint64 {
	total := uint64(0)

	for _, frame := range e.Frames {
		total += frame.ExecutionGas
	}

	return total
}

// TotalStateGas returns the state gas attributed across all frames after refills.
func (e *FrameReceiptExtra) TotalStateGas() uint64 {
	total := uint64(0)

	for _, frame := range e.Frames {
		total += frame.StateGas
	}

	return total
}

// FailedFrame returns the index of the first frame that failed, or -1.
func (e *FrameReceiptExtra) FailedFrame() int {
	for i, frame := range e.Frames {
		if frame.Status == FrameStatusFailed {
			return i
		}
	}

	return -1
}

// DerivedStatus returns the transaction status derived from the frame statuses: a
// transaction succeeded when no executed frame failed.
func (e *FrameReceiptExtra) DerivedStatus() uint64 {
	if e.FailedFrame() >= 0 {
		return ReceiptStatusFailed
	}

	return ReceiptStatusSuccessful
}

// FrameExtra returns the frame-specific content of a receipt, or nil.
func (r *Receipt) FrameExtra() *FrameReceiptExtra {
	extra, _ := r.Extra.(*FrameReceiptExtra)

	return extra
}

// jsonFrameReceipt is the JSON-RPC representation of one frame receipt.
//
// EIP-8141 does not specify a JSON-RPC encoding, so the decoder accepts the shapes
// clients are likely to use: the two-dimensional gas as an object, as a two-element
// array mirroring the consensus encoding, or as flat sibling fields.
type jsonFrameReceipt struct {
	Status  *hexutil.Uint64 `json:"status"`
	GasUsed json.RawMessage `json:"gasUsed"`
	Logs    []*jsonLog      `json:"logs"`

	ExecutionGasUsed *hexutil.Uint64 `json:"executionGasUsed"`
	StateGasUsed     *hexutil.Uint64 `json:"stateGasUsed"`
}

// jsonFrameReceiptExtra is the frame-specific part of a JSON-RPC receipt. Clients
// differ on the key: ethrex reports "frameReceipts".
type jsonFrameReceiptExtra struct {
	Payer         *common.Address     `json:"payer"`
	Frames        []*jsonFrameReceipt `json:"frames"`
	FrameReceipts []*jsonFrameReceipt `json:"frameReceipts"`
}

// frameList returns whichever key the client used.
func (e *jsonFrameReceiptExtra) frameList() []*jsonFrameReceipt {
	if e.Frames != nil {
		return e.Frames
	}

	return e.FrameReceipts
}

// decodeFrameReceipt parses the frame-specific receipt fields and derives the generic
// status and gas totals the rest of the engine reads.
func decodeFrameReceipt(receipt *Receipt, raw json.RawMessage) error {
	var dec jsonFrameReceiptExtra
	if err := json.Unmarshal(raw, &dec); err != nil {
		return err
	}

	if dec.Payer == nil && dec.frameList() == nil {
		// The node reports a frame transaction without the frame fields. Keep the
		// generic receipt rather than failing the whole block's decoding.
		return nil
	}

	frames := dec.frameList()

	extra := &FrameReceiptExtra{
		Frames: make([]*FrameReceipt, 0, len(frames)),
	}

	if dec.Payer != nil {
		extra.Payer = *dec.Payer
	}

	for _, frame := range frames {
		decoded := &FrameReceipt{Logs: toLogs(frame.Logs, receipt)}

		if frame.Status != nil {
			decoded.Status = uint64(*frame.Status)
		}

		decoded.ExecutionGas, decoded.StateGas = frameGasUsed(frame)
		extra.Frames = append(extra.Frames, decoded)
	}

	receipt.Extra = extra

	// The consensus receipt has no transaction-level status, so derive one unless the
	// node already reported it.
	if !json.Valid(raw) || !hasJSONField(raw, "status") {
		receipt.Status = extra.DerivedStatus()
	}

	if receipt.GasUsed == 0 {
		receipt.GasUsed = extra.TotalExecutionGas() + extra.TotalStateGas()
	}

	if len(receipt.Logs) == 0 {
		for _, frame := range extra.Frames {
			receipt.Logs = append(receipt.Logs, frame.Logs...)
		}
	}

	return nil
}

// frameGasUsed extracts a frame's two gas dimensions from whichever shape the node used.
func frameGasUsed(frame *jsonFrameReceipt) (execution, state uint64) {
	if frame.ExecutionGasUsed != nil {
		execution = uint64(*frame.ExecutionGasUsed)
	}

	if frame.StateGasUsed != nil {
		state = uint64(*frame.StateGasUsed)
	}

	if len(frame.GasUsed) == 0 {
		return execution, state
	}

	var pair []hexutil.Uint64
	if err := json.Unmarshal(frame.GasUsed, &pair); err == nil {
		if len(pair) > 0 {
			execution = uint64(pair[0])
		}

		if len(pair) > 1 {
			state = uint64(pair[1])
		}

		return execution, state
	}

	var object struct {
		Execution *hexutil.Uint64 `json:"execution"`
		State     *hexutil.Uint64 `json:"state"`
	}

	if err := json.Unmarshal(frame.GasUsed, &object); err == nil {
		if object.Execution != nil {
			execution = uint64(*object.Execution)
		}

		if object.State != nil {
			state = uint64(*object.State)
		}

		return execution, state
	}

	// A plain number means the node reports a single combined value.
	var combined hexutil.Uint64
	if err := json.Unmarshal(frame.GasUsed, &combined); err == nil {
		execution = uint64(combined)
	}

	return execution, state
}

// hasJSONField reports whether a JSON object carries a non-null field.
func hasJSONField(raw json.RawMessage, field string) bool {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		return false
	}

	value, ok := fields[field]

	return ok && string(value) != "null"
}

// blobGasCost returns the blob fee a frame transaction pays at the given base fee.
func (tx *FrameTx) blobGasCost(blobBaseFee *big.Int) *big.Int {
	if len(tx.BlobHashes) == 0 || blobBaseFee == nil {
		return new(big.Int)
	}

	blobGas := new(big.Int).SetUint64(uint64(len(tx.BlobHashes)) * GasPerBlob)

	return blobGas.Mul(blobGas, blobBaseFee)
}

// EncodeReceiptJSON adds the frame-specific receipt fields, using the key names the
// decoder reads back.
func (e *FrameReceiptExtra) EncodeReceiptJSON(fields map[string]any) {
	fields["payer"] = e.Payer

	frames := make([]map[string]any, 0, len(e.Frames))

	for _, frame := range e.Frames {
		logs := frame.Logs
		if logs == nil {
			logs = []*Log{}
		}

		frames = append(frames, map[string]any{
			"status":       hexutil.Uint64(frame.Status),
			"gasUsed":      hexutil.Uint64(frame.ExecutionGas),
			"stateGasUsed": hexutil.Uint64(frame.StateGas),
			"logs":         logs,
		})
	}

	fields["frameReceipts"] = frames
}

// PostTxReverted reports whether a POST_TX frame failed, which under EIP-7906 reverts
// the whole execution body rather than one atomic batch.
//
// The receipt alone cannot answer this: frame modes live on the transaction, so the
// transaction has to be supplied.
func (e *FrameReceiptExtra) PostTxReverted(tx *FrameTx) bool {
	postTx := tx.PostTxIndex()
	if postTx < 0 {
		return false
	}

	for i := postTx; i < len(e.Frames) && i < len(tx.Frames); i++ {
		if e.Frames[i].Status == FrameStatusFailed {
			return true
		}
	}

	return false
}

// DurableFrames reports, per frame, whether that frame's state changes survived the
// transaction.
//
// A frame receipt's status says whether the frame executed successfully, not whether
// its effects lasted. Two rules discard the effects of frames that report success:
//
//   - an atomic batch that fails is unrolled, so the frames before the failure inside
//     it lose their state changes and logs while keeping their status;
//   - a failing POST_TX frame (EIP-7906) reverts the entire execution body, so every
//     frame after the validation prefix loses its effects.
//
// Frames of the validation prefix are always durable: their changes are committed even
// when the body reverts.
func (e *FrameReceiptExtra) DurableFrames(tx *FrameTx) []bool {
	durable := make([]bool, len(e.Frames))

	prefix := tx.ValidationPrefixLength()
	bodyReverted := e.PostTxReverted(tx)

	for i := range e.Frames {
		switch {
		case i < prefix:
			// Committed regardless of what the body does.
			durable[i] = e.Frames[i].Status == FrameStatusSuccess
		case bodyReverted:
			durable[i] = false
		default:
			durable[i] = e.Frames[i].Status == FrameStatusSuccess
		}
	}

	if bodyReverted {
		return durable
	}

	// Unroll failed atomic batches. A batch is a maximal run of frames where all but
	// the last carry the flag.
	for start := prefix; start < len(tx.Frames); {
		end := start
		for end < len(tx.Frames) && tx.Frames[end].IsAtomicBatch() {
			end++
		}

		if end == start {
			start++

			continue
		}

		// end is the frame that closes the batch.
		failed := false

		for i := start; i <= end && i < len(e.Frames); i++ {
			if e.Frames[i].Status != FrameStatusSuccess {
				failed = true

				break
			}
		}

		if failed {
			for i := start; i <= end && i < len(durable); i++ {
				durable[i] = false
			}
		}

		start = end + 1
	}

	return durable
}
