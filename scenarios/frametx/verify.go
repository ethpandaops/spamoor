package frametx

import (
	"fmt"
	"strings"

	"github.com/ethpandaops/spamoor/txtypes"
)

// statusNames maps frame receipt statuses to readable names.
var statusNames = map[uint64]string{
	txtypes.FrameStatusFailed:  "fail",
	txtypes.FrameStatusSuccess: "ok",
	txtypes.FrameStatusSkipped: "skip",
}

// statusName returns a readable name for a frame status.
func statusName(status uint64) string {
	if name, ok := statusNames[status]; ok {
		return name
	}

	return fmt.Sprintf("0x%x", status)
}

// formatStatuses renders the per-frame statuses of a receipt, e.g. "ok,ok,fail,skip".
func formatStatuses(extra *txtypes.FrameReceiptExtra) string {
	names := make([]string, len(extra.Frames))
	for i, frame := range extra.Frames {
		names[i] = statusName(frame.Status)
	}

	return strings.Join(names, ",")
}

// formatDurable renders which frames' effects survived, e.g. "yes,no,no". A frame can
// report success and still have been discarded, by an unrolled atomic batch or a
// failed POST_TX frame reverting the execution body.
func formatDurable(extra *txtypes.FrameReceiptExtra, tx *txtypes.FrameTx) string {
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

// compareStatuses checks a receipt's frame statuses against what the shape should have
// produced. It returns an empty string when they agree, and a description of the
// disagreement otherwise.
func compareStatuses(expected []uint64, extra *txtypes.FrameReceiptExtra) string {
	if len(expected) != len(extra.Frames) {
		return fmt.Sprintf("expected %d frame receipts, got %d (%s)",
			len(expected), len(extra.Frames), formatStatuses(extra))
	}

	for i, want := range expected {
		if extra.Frames[i].Status != want {
			return fmt.Sprintf("frame %d has status %s, expected %s (got %s)",
				i, statusName(extra.Frames[i].Status), statusName(want), formatStatuses(extra))
		}
	}

	return ""
}
