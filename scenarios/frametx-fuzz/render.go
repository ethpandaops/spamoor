package frametxfuzz

import (
	"strings"

	"github.com/ethpandaops/spamoor/txtypes"
)

// Rendering for the per-frame result of a landed transaction. These describe what a
// receipt said; nothing here decides whether it should have said it.

// statusNames maps a frame receipt status onto a readable name.
var statusNames = map[uint64]string{
	txtypes.FrameStatusFailed:  "fail",
	txtypes.FrameStatusSuccess: "ok",
	txtypes.FrameStatusSkipped: "skip",
}

// renderStatuses renders the reported per-frame statuses, e.g. "ok,ok,fail,skip".
func renderStatuses(extra *txtypes.FrameReceiptExtra) string {
	names := make([]string, len(extra.Frames))

	for i, frame := range extra.Frames {
		name, ok := statusNames[frame.Status]
		if !ok {
			name = "?"
		}

		names[i] = name
	}

	return strings.Join(names, ",")
}

// renderDurable renders which frames' effects survived.
//
// A frame can report success and still have been discarded, by an unrolled atomic batch
// or by a failing POST_TX frame reverting the execution body, so the two are worth
// showing side by side.
func renderDurable(extra *txtypes.FrameReceiptExtra, tx *txtypes.FrameTx) string {
	if tx == nil {
		return "?"
	}

	durable := extra.DurableFrames(tx)
	names := make([]string, len(durable))

	for i, ok := range durable {
		names[i] = "no"
		if ok {
			names[i] = "yes"
		}
	}

	return strings.Join(names, ",")
}
