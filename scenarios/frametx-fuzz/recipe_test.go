package frametxfuzz

import (
	"testing"

	"github.com/ethpandaops/spamoor/utils"
)

// testDrawOptions enables everything, so a draw exercises every axis.
func testDrawOptions() DrawOptions {
	axes := axisWeights{}
	for _, name := range axisNames {
		axes[name] = 1
	}

	return DrawOptions{
		Axes:          axes,
		MaxBodyFrames: 6,
		AllowPostTx:   true,
		AllowContract: true,
		AllowRoots:    true,
		AllowKeyed:    true,
		AllowProbe:    true,
	}
}

// TestDrawIsDeterministic checks the property the whole reproduction story rests on: the
// same seed and index must produce the same recipe, and neighbouring indices must not
// produce the same one.
func TestDrawIsDeterministic(t *testing.T) {
	const seed = "0xfeedface"

	for index := uint64(0); index < 64; index++ {
		first := Draw(utils.NewDeterministicRNGWithSeed(index, seed), index, testDrawOptions())
		second := Draw(utils.NewDeterministicRNGWithSeed(index, seed), index, testDrawOptions())

		if first.String() != second.String() {
			t.Fatalf("index %d drew two different recipes:\n%s\n%s", index, first, second)
		}
	}

	same := 0

	for index := uint64(0); index < 64; index++ {
		a := Draw(utils.NewDeterministicRNGWithSeed(index, seed), index, testDrawOptions())
		b := Draw(utils.NewDeterministicRNGWithSeed(index+1, seed), index+1, testDrawOptions())

		if a.String() == b.String() {
			same++
		}
	}

	// Some collisions are expected in a small space; a stream that never varies is not.
	if same > 32 {
		t.Errorf("%d of 64 neighbouring indices drew identical recipes", same)
	}
}

// TestDrawRoundTripsThroughJSON checks that a recipe reported with a finding replays to
// the same transaction description.
func TestDrawRoundTripsThroughJSON(t *testing.T) {
	recipe := Draw(utils.NewDeterministicRNGWithSeed(7, "0x01"), 7, testDrawOptions())

	parsed, err := ParseRecipe(recipe.String())
	if err != nil {
		t.Fatalf("could not parse the reported recipe: %v", err)
	}

	if parsed.String() != recipe.String() {
		t.Errorf("recipe changed through its JSON form:\n%s\n%s", recipe, parsed)
	}
}

// TestDrawRespectsGates checks that a disabled capability never appears in a recipe,
// which is what keeps the generator from producing transactions the chain must refuse
// for reasons that have nothing to do with the client.
func TestDrawRespectsGates(t *testing.T) {
	opts := testDrawOptions()
	opts.AllowPostTx = false
	opts.AllowRoots = false
	opts.AllowKeyed = false
	opts.AllowProbe = false
	opts.AllowContract = false

	for index := uint64(0); index < 128; index++ {
		recipe := Draw(utils.NewDeterministicRNGWithSeed(index, "0x02"), index, opts)

		if recipe.NonceKeys != 0 {
			t.Fatalf("index %d selected nonce keys with EIP-8250 disabled", index)
		}

		if recipe.RecentRoots != 0 {
			t.Fatalf("index %d declared recent roots with EIP-8272 disabled", index)
		}

		if recipe.Sender == SenderContract {
			t.Fatalf("index %d used a contract sender with no probe contract", index)
		}

		for _, frame := range recipe.Body {
			if frame.Kind == KindPostTx {
				t.Fatalf("index %d generated a POST_TX frame with EIP-7906 disabled", index)
			}

			if frame.Kind == KindProbe {
				t.Fatalf("index %d generated a probe frame with no probe contract", index)
			}
		}
	}
}

// TestDrawNormalizesFormatConstraints checks that a drawn recipe never describes a
// transaction the payload format forbids: the negative tier is the only source of
// invalid transactions, so an accidental one would be reported as a client finding.
func TestDrawNormalizesFormatConstraints(t *testing.T) {
	for index := uint64(0); index < 256; index++ {
		recipe := Draw(utils.NewDeterministicRNGWithSeed(index, "0x03"), index, testDrawOptions())

		seenPostTx := false

		for i, frame := range recipe.Body {
			if frame.Kind == KindPostTx {
				seenPostTx = true
			} else if seenPostTx {
				t.Fatalf("index %d places frame %d after a POST_TX frame", index, i)
			}

			if frame.Batch && i == len(recipe.Body)-1 {
				t.Fatalf("index %d batches its last frame", index)
			}

			if frame.Kind == KindPostTx && frame.Batch {
				t.Fatalf("index %d batches a POST_TX frame", index)
			}

			if frame.Kind == KindPostTx && frame.Script != ScriptNone {
				t.Fatalf("index %d gives a POST_TX frame a writing script", index)
			}
		}
	}
}
