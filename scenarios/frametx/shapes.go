package frametx

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"

	"github.com/ethpandaops/spamoor/txtypes"
)

// A shape is one frame layout the scenario knows how to build, together with the frame
// statuses a correct client must report for it. Comparing the two turns the load
// generator into a conformance check.
type shapeName string

const (
	shapeSelfVerify shapeName = "self-verify"
	shapeTransfer   shapeName = "transfer"
	shapeBatch      shapeName = "batch"
	shapeAtomic     shapeName = "atomic"
	shapeAtomicFail shapeName = "atomic-fail"
	shapeExpiry     shapeName = "expiry"
)

var allShapes = []shapeName{
	shapeSelfVerify,
	shapeTransfer,
	shapeBatch,
	shapeAtomic,
	shapeAtomicFail,
	shapeExpiry,
}

// shapeParams carries everything a shape needs to build its frames.
type shapeParams struct {
	sender      common.Address
	target      common.Address
	amount      *uint256.Int
	data        []byte
	userOpGas   uint64
	verifyGas   uint64
	stateGas    uint64
	frameCount  uint64
	expiryAt    uint64
	targetEmpty bool
}

// builtShape is a shape's frame list together with its expected frame statuses.
type builtShape struct {
	name           shapeName
	frames         []*txtypes.Frame
	expectedStatus []uint64
}

// buildShape assembles the frames for a shape.
func buildShape(name shapeName, p shapeParams) (*builtShape, error) {
	verifyLimits := txtypes.FrameLimits{Execution: p.verifyGas}

	userOp := func(value *uint256.Int, stateGas uint64) *txtypes.Frame {
		target := p.target

		return txtypes.UserOpFrame(&target, value, p.data, txtypes.FrameLimits{
			Execution: p.userOpGas,
			State:     stateGas,
		})
	}

	switch name {
	case shapeSelfVerify:
		// The minimal mempool-legal shape: validate, then one call.
		return &builtShape{
			name: name,
			frames: []*txtypes.Frame{
				txtypes.SelfVerifyFrame(verifyLimits),
				userOp(nil, 0),
			},
			expectedStatus: []uint64{txtypes.FrameStatusSuccess, txtypes.FrameStatusSuccess},
		}, nil

	case shapeTransfer:
		// A value-bearing call. Creating the recipient costs state gas, which the
		// frame must budget for itself.
		stateGas := uint64(0)
		if p.targetEmpty {
			stateGas = txtypes.StateBytesPerNewAccount * txtypes.CostPerStateByte
		}

		return &builtShape{
			name: name,
			frames: []*txtypes.Frame{
				txtypes.SelfVerifyFrame(verifyLimits),
				userOp(p.amount, stateGas+p.stateGas),
			},
			expectedStatus: []uint64{txtypes.FrameStatusSuccess, txtypes.FrameStatusSuccess},
		}, nil

	case shapeBatch:
		// Many independent calls in one transaction, exercising the per-frame cost and
		// block packing across both gas dimensions.
		count := p.frameCount
		if count < 1 {
			count = 1
		}

		if count > txtypes.MaxFrames-1 {
			count = txtypes.MaxFrames - 1
		}

		frames := make([]*txtypes.Frame, 0, count+1)
		frames = append(frames, txtypes.SelfVerifyFrame(verifyLimits))

		status := []uint64{txtypes.FrameStatusSuccess}

		for i := uint64(0); i < count; i++ {
			frames = append(frames, userOp(nil, 0))
			status = append(status, txtypes.FrameStatusSuccess)
		}

		return &builtShape{name: name, frames: frames, expectedStatus: status}, nil

	case shapeAtomic:
		// Two calls that must both succeed or both revert, plus a trailing call that
		// closes the batch.
		return &builtShape{
			name: name,
			frames: []*txtypes.Frame{
				txtypes.SelfVerifyFrame(verifyLimits),
				userOp(nil, 0).WithAtomicBatch(),
				userOp(nil, 0).WithAtomicBatch(),
				userOp(nil, 0),
			},
			expectedStatus: []uint64{
				txtypes.FrameStatusSuccess,
				txtypes.FrameStatusSuccess,
				txtypes.FrameStatusSuccess,
				txtypes.FrameStatusSuccess,
			},
		}, nil

	case shapeAtomicFail:
		// The middle frame is starved of execution gas so it halts on the frame-entry
		// account access charge. The batch rolls back and the frame after the failure
		// is skipped, which is the only way to observe status 2.
		return &builtShape{
			name: name,
			frames: []*txtypes.Frame{
				txtypes.SelfVerifyFrame(verifyLimits),
				userOp(nil, 0).WithAtomicBatch(),
				txtypes.UserOpFrame(&p.target, nil, p.data,
					txtypes.FrameLimits{Execution: 1}).WithAtomicBatch(),
				userOp(nil, 0),
			},
			expectedStatus: []uint64{
				txtypes.FrameStatusSuccess,
				txtypes.FrameStatusSuccess, // executed, then rolled back
				txtypes.FrameStatusFailed,
				txtypes.FrameStatusSkipped,
			},
		}, nil

	case shapeExpiry:
		// The expiry verifier frame is skipped when matching mempool prefixes, so this
		// still counts as a self-relayed transaction.
		return &builtShape{
			name: name,
			frames: []*txtypes.Frame{
				txtypes.ExpiryFrame(p.expiryAt, p.verifyGas),
				txtypes.SelfVerifyFrame(verifyLimits),
				userOp(nil, 0),
			},
			expectedStatus: []uint64{
				txtypes.FrameStatusSuccess,
				txtypes.FrameStatusSuccess,
				txtypes.FrameStatusSuccess,
			},
		}, nil
	}

	return nil, fmt.Errorf("unknown frame shape %q", name)
}

// weightedShapes selects shapes according to configured weights.
type weightedShapes struct {
	names   []shapeName
	weights []uint64
	total   uint64
}

// parseShapes parses a "name:weight,name:weight" selection. An empty spec, or "all",
// enables every shape with equal weight.
func parseShapes(spec string) (*weightedShapes, error) {
	spec = strings.TrimSpace(spec)
	if spec == "" || spec == "all" {
		selection := &weightedShapes{}
		for _, name := range allShapes {
			selection.add(name, 1)
		}

		return selection, nil
	}

	selection := &weightedShapes{}

	for _, entry := range strings.Split(spec, ",") {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		name := entry
		weight := uint64(1)

		if idx := strings.LastIndex(entry, ":"); idx >= 0 {
			name = strings.TrimSpace(entry[:idx])

			parsed, err := strconv.ParseUint(strings.TrimSpace(entry[idx+1:]), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid weight in %q: %w", entry, err)
			}

			weight = parsed
		}

		if !isKnownShape(shapeName(name)) {
			return nil, fmt.Errorf("unknown frame shape %q, known shapes: %s", name, strings.Join(shapeNames(), ", "))
		}

		if weight > 0 {
			selection.add(shapeName(name), weight)
		}
	}

	if selection.total == 0 {
		return nil, fmt.Errorf("no frame shapes enabled")
	}

	return selection, nil
}

// add registers a shape with a weight.
func (w *weightedShapes) add(name shapeName, weight uint64) {
	w.names = append(w.names, name)
	w.weights = append(w.weights, weight)
	w.total += weight
}

// pick deterministically selects a shape for a transaction index.
func (w *weightedShapes) pick(txIdx uint64) shapeName {
	position := txIdx % w.total
	cursor := uint64(0)

	for i, weight := range w.weights {
		cursor += weight
		if position < cursor {
			return w.names[i]
		}
	}

	return w.names[0]
}

// isKnownShape reports whether name is a shape this scenario can build.
func isKnownShape(name shapeName) bool {
	for _, known := range allShapes {
		if known == name {
			return true
		}
	}

	return false
}

// shapeNames returns the known shape names.
func shapeNames() []string {
	names := make([]string, len(allShapes))
	for i, name := range allShapes {
		names[i] = string(name)
	}

	return names
}
