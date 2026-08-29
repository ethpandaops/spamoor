package frametxfuzz

import (
	"encoding/json"
	"fmt"

	"github.com/ethpandaops/spamoor/utils"
)

// Recipe is the abstract description of one generated transaction.
//
// It is drawn from the seeded generator in a single pass and contains no chain state:
// which wallet ends up sending it, what nonce it gets and where the probe contract
// lives are all applied afterwards. That separation is what makes a run reproducible.
// If a draw depended on chain state, the same seed would produce different transactions
// on a chain in a different state and the reproduction line printed with a finding
// would be a lie.
type Recipe struct {
	// Index is the transaction index the recipe was drawn for.
	Index uint64 `json:"index"`

	// Prefix is the validation prefix shape.
	Prefix PrefixShape `json:"prefix"`

	// Expiry adds a leading expiry verifier frame, which the mempool rules skip when
	// matching the prefix against the recognized shapes.
	Expiry bool `json:"expiry"`

	// Sender selects whether the sender is an account without code, running the
	// protocol's default code, or one carrying the probe delegation.
	Sender SenderKind `json:"sender"`

	// Body is the frame list after the validation prefix.
	Body []BodyFrame `json:"body"`

	// Signatures describes the signature list beyond the sender's own entry.
	Witness bool `json:"witness"`

	// P256 adds a P256 entry, whose verification gas is more than twice a secp256k1
	// entry's and counts against the same validation prefix cap.
	P256 bool `json:"p256"`

	// NonceKeys is how many EIP-8250 keys the transaction selects. Zero means the
	// legacy account nonce.
	NonceKeys int `json:"nonceKeys"`

	// RecentRoots is how many EIP-8272 references the transaction declares.
	RecentRoots int `json:"recentRoots"`

	// RecentRootEdge picks a deliberately awkward reference, named by
	// recentRootEdgeNames.
	RecentRootEdge string `json:"recentRootEdge,omitempty"`

	// Reads turns on the introspection sweep: the operations that make every
	// instruction the frame transaction EIPs introduce execute inside a probe frame.
	Reads bool `json:"reads"`

	// Invalid names a deliberate violation the transaction carries, if any.
	//
	// Malformed transactions are drawn from the same stream as well-formed ones rather
	// than living in a separate mode: the edge cases worth reaching are combinations,
	// and a combination that happens to be illegal is one of them. They are sent from
	// a burner wallet outside the managed pool, because a transaction that never lands
	// would otherwise stall its sender's nonce.
	Invalid string `json:"invalid,omitempty"`
}

// PrefixShape is one of the validation prefixes the public mempool recognizes.
type PrefixShape string

// The recognized prefixes. The two deploy-led shapes are absent on purpose: a CREATE2
// account deployment costs far more than the prefix's execution cap allows, and the
// sender of such a transaction has no key, which the managed transaction pool cannot
// represent. They belong to the raw negative tier.
const (
	PrefixSelfVerify PrefixShape = "self_verify"
	PrefixPaymaster  PrefixShape = "only_verify+pay"
)

// SenderKind selects what code runs in the sender's validation frame.
type SenderKind string

// Sender kinds.
const (
	// SenderDefaultCode is an account with no code, validated by the protocol's
	// default code against the signature list.
	SenderDefaultCode SenderKind = "default_code"

	// SenderContract is an account carrying the probe delegation, so APPROVE runs from
	// deployed code. One client build accepts this and another reverts the prefix, so
	// it is tracked as a known divergence rather than a finding.
	SenderContract SenderKind = "contract"
)

// FrameKind is what a body frame does.
type FrameKind string

// Body frame kinds.
const (
	// KindTransfer is a SENDER frame moving value, the only mode that may carry any.
	KindTransfer FrameKind = "transfer"

	// KindCall is a SENDER frame calling a wallet with data.
	KindCall FrameKind = "call"

	// KindProbe is a SENDER frame calling the probe contract with a script.
	KindProbe FrameKind = "probe"

	// KindPostOp is a DEFAULT frame with flags, the settlement species.
	KindPostOp FrameKind = "post_op"

	// KindPostTx is an EIP-7906 assertion frame, restricted to a trailing suffix.
	KindPostTx FrameKind = "post_tx"
)

// FrameBudget selects how much execution gas a frame is given relative to what it needs.
type FrameBudget string

// Frame budgets.
const (
	// BudgetAmple is comfortably more than the frame needs.
	BudgetAmple FrameBudget = "ample"

	// BudgetStarved is one gas, which cannot even cover the frame-entry account access
	// charge, so the frame halts before its code runs. This is the only budget whose
	// outcome is certain, which is what makes it usable as an oracle.
	BudgetStarved FrameBudget = "starved"
)

// ScriptKind is what a probe frame's script does beyond its assertions.
type ScriptKind string

// Script kinds.
const (
	ScriptNone   ScriptKind = "none"
	ScriptRevert ScriptKind = "revert"
	ScriptLog    ScriptKind = "log"
	ScriptStore  ScriptKind = "store"
)

// TargetKind is what a frame addresses.
type TargetKind string

// Target kinds.
const (
	TargetWallet     TargetKind = "wallet"
	TargetSender     TargetKind = "sender"
	TargetProbe      TargetKind = "probe"
	TargetEmpty      TargetKind = "empty"
	TargetPrecompile TargetKind = "precompile"
)

// BodyFrame describes one frame after the validation prefix.
type BodyFrame struct {
	Kind   FrameKind   `json:"kind"`
	Target TargetKind  `json:"target"`
	Budget FrameBudget `json:"budget"`
	Script ScriptKind  `json:"script"`

	// Batch marks the frame as continuing into the next one. A batch is the maximal
	// run of frames where all but the last carry the flag, so the closing frame does
	// not.
	Batch bool `json:"batch"`

	// StateGas budgets state gas for a frame that creates state.
	StateGas bool `json:"stateGas"`
}

// recentRootEdgeNames are the awkward reference cases worth generating. Each is
// expected to be refused except the plain one, which must be accepted.
var recentRootEdgeNames = []string{
	"",             // a reference to a root written a slot ago: must be accepted
	"same_slot",    // current_slot - slot == 0
	"future_slot",  // slot ahead of the current one
	"unwritten",    // a slot the source never wrote
	"wrong_source", // a correct root under a source that did not write it
	"duplicate",    // the same reference twice, which is valid and charged twice
	"outside_window",
}

// String renders the recipe as the compact JSON a finding reports.
func (r *Recipe) String() string {
	encoded, err := json.Marshal(r)
	if err != nil {
		return fmt.Sprintf("recipe %d (unrenderable: %v)", r.Index, err)
	}

	return string(encoded)
}

// ParseRecipe reads a recipe back from the JSON a finding reported, which is how a
// single generated transaction is replayed on a chain that has moved on.
func ParseRecipe(encoded string) (*Recipe, error) {
	recipe := &Recipe{}
	if err := json.Unmarshal([]byte(encoded), recipe); err != nil {
		return nil, fmt.Errorf("could not parse the recipe: %w", err)
	}

	return recipe, nil
}

// axis names the dimensions a run can weight, so an operator can point a run at one EIP.
type axis string

// The weightable axes.
const (
	axisPrefix    axis = "prefix"
	axisBatches   axis = "batches"
	axisFailures  axis = "failures"
	axisSignature axis = "signatures"
	axisNonces    axis = "nonces"
	axisRoots     axis = "roots"
	axisPostTx    axis = "posttx"
	axisProbe     axis = "probe"
)

// axisNames lists the axes in a stable order for help text and validation.
var axisNames = []axis{
	axisPrefix, axisBatches, axisFailures, axisSignature,
	axisNonces, axisRoots, axisPostTx, axisProbe,
}

// axisWeights is a validated axis selection.
type axisWeights map[axis]uint64

// enabled reports whether an axis may be drawn.
func (w axisWeights) enabled(name axis) bool { return w[name] > 0 }

// chance returns the probability an axis contributes to a transaction, normalised
// against the highest weight in the selection so that "roots:2,nonces:1" means roots
// twice as often as nonces rather than a fixed rate.
func (w axisWeights) chance(name axis) float64 {
	weight := w[name]
	if weight == 0 {
		return 0
	}

	highest := uint64(0)
	for _, value := range w {
		if value > highest {
			highest = value
		}
	}

	return float64(weight) / float64(highest)
}

// draws returns true with the axis's probability.
func (w axisWeights) draws(rng *utils.DeterministicRNG, name axis) bool {
	chance := w.chance(name)
	if chance == 0 {
		return false
	}

	return rng.Float64() < chance
}

// DrawOptions bounds what a draw may produce.
type DrawOptions struct {
	Axes          axisWeights
	MaxBodyFrames int
	AllowPostTx   bool
	AllowContract bool
	AllowRoots    bool
	AllowKeyed    bool
	AllowProbe    bool
	// InvalidChance is how often a drawn recipe carries a deliberate violation.
	InvalidChance float64

	// Violations is the catalog a violation is drawn from.
	Violations []string
}

// Draw produces the recipe for one transaction index.
//
// Every field is drawn unconditionally and only then gated by the options. Drawing
// inside a conditional would make the stream position depend on the configuration, so
// two runs that differ in one flag would disagree about every later transaction rather
// than about the one axis that changed.
func Draw(rng *utils.DeterministicRNG, index uint64, opts DrawOptions) *Recipe {
	recipe := &Recipe{Index: index}

	prefixDraw := rng.Float64()
	expiryDraw := rng.Float64()
	senderDraw := rng.Float64()
	witnessDraw := rng.Float64()
	p256Draw := rng.Float64()
	assertDraw := rng.Float64()

	nonceKeyDraw := rng.Intn(4) + 1
	nonceDraw := rng.Float64()

	rootCountDraw := rng.Intn(3) + 1
	rootEdgeDraw := rng.Intn(len(recentRootEdgeNames))
	rootDraw := rng.Float64()

	bodyCount := rng.Intn(opts.MaxBodyFrames) + 1
	body := make([]BodyFrame, 0, bodyCount)

	for i := 0; i < bodyCount; i++ {
		body = append(body, drawBodyFrame(rng, opts))
	}

	postTxCount := rng.Intn(2)
	postTxDraw := rng.Float64()
	invalidDraw := rng.Float64()
	invalidPick := rng.Intn(max(len(opts.Violations), 1))

	// Prefix: the sponsored shape is the only one whose payer is not the sender, so it
	// is worth a substantial share of the stream rather than a rare case.
	recipe.Prefix = PrefixSelfVerify
	if opts.Axes.enabled(axisPrefix) && prefixDraw < 0.35*opts.Axes.chance(axisPrefix) {
		recipe.Prefix = PrefixPaymaster
	}

	recipe.Expiry = opts.Axes.enabled(axisPrefix) && expiryDraw < 0.25

	recipe.Sender = SenderDefaultCode
	if opts.AllowContract && senderDraw < 0.25 {
		recipe.Sender = SenderContract
	}

	recipe.Witness = opts.Axes.enabled(axisSignature) && witnessDraw < 0.4*opts.Axes.chance(axisSignature)
	recipe.P256 = opts.Axes.enabled(axisSignature) && p256Draw < 0.2*opts.Axes.chance(axisSignature)

	if opts.AllowKeyed && opts.Axes.draws(rng, axisNonces) && nonceDraw < 0.6 {
		recipe.NonceKeys = nonceKeyDraw
	}

	if opts.AllowRoots && opts.Axes.draws(rng, axisRoots) && rootDraw < 0.6 {
		recipe.RecentRoots = rootCountDraw
		recipe.RecentRootEdge = recentRootEdgeNames[rootEdgeDraw]
	}

	recipe.Body = body

	if opts.AllowPostTx && opts.Axes.enabled(axisPostTx) && postTxDraw < 0.3*opts.Axes.chance(axisPostTx) {
		for i := 0; i <= postTxCount; i++ {
			recipe.Body = append(recipe.Body, BodyFrame{
				Kind:   KindPostTx,
				Target: TargetProbe,
				Budget: BudgetAmple,
				Script: ScriptNone,
			})
		}
	}

	recipe.Reads = opts.AllowProbe && opts.Axes.enabled(axisProbe) && assertDraw < 0.5*opts.Axes.chance(axisProbe)

	if len(opts.Violations) > 0 && invalidDraw < opts.InvalidChance {
		recipe.Invalid = opts.Violations[invalidPick]
	}

	recipe.normalize(opts)

	return recipe
}

// drawBodyFrame draws one frame after the validation prefix.
func drawBodyFrame(rng *utils.DeterministicRNG, opts DrawOptions) BodyFrame {
	kindDraw := rng.Float64()
	targetDraw := rng.Float64()
	budgetDraw := rng.Float64()
	scriptDraw := rng.Float64()
	batchDraw := rng.Float64()

	frame := BodyFrame{Kind: KindCall, Target: TargetWallet, Budget: BudgetAmple, Script: ScriptNone}

	switch {
	case kindDraw < 0.30:
		frame.Kind = KindTransfer
	case kindDraw < 0.60 && opts.AllowProbe && opts.Axes.enabled(axisProbe):
		frame.Kind = KindProbe
		frame.Target = TargetProbe
	case kindDraw < 0.70:
		frame.Kind = KindPostOp
	}

	if frame.Kind != KindProbe {
		switch {
		case targetDraw < 0.15:
			frame.Target = TargetSender
		case targetDraw < 0.30:
			frame.Target = TargetEmpty
			frame.StateGas = true
		case targetDraw < 0.38:
			frame.Target = TargetPrecompile
		}
	}

	// A starved frame is the reliable way to produce a failure: the frame-entry account
	// access is charged before anything runs, so one gas cannot cover it whatever the
	// target does.
	if opts.Axes.enabled(axisFailures) && budgetDraw < 0.15*opts.Axes.chance(axisFailures) {
		frame.Budget = BudgetStarved
	}

	if frame.Kind == KindProbe && opts.Axes.enabled(axisFailures) && scriptDraw < 0.15 {
		frame.Script = ScriptRevert
	} else if frame.Kind == KindProbe {
		switch {
		case scriptDraw < 0.45:
			frame.Script = ScriptLog
		case scriptDraw < 0.60:
			frame.Script = ScriptStore
			frame.StateGas = true
		}
	}

	frame.Batch = opts.Axes.enabled(axisBatches) && batchDraw < 0.35*opts.Axes.chance(axisBatches)

	return frame
}

// normalize repairs the combinations the drawn fields can produce but the transaction
// format forbids, so that a recipe always describes a well-formed transaction and the
// negative tier stays the only source of invalid ones.
func (r *Recipe) normalize(opts DrawOptions) {
	// POST_TX frames must be a trailing suffix, so anything drawn after one moves in
	// front of it.
	ordinary := make([]BodyFrame, 0, len(r.Body))
	postTx := make([]BodyFrame, 0, len(r.Body))

	for _, frame := range r.Body {
		if frame.Kind == KindPostTx {
			postTx = append(postTx, frame)

			continue
		}

		ordinary = append(ordinary, frame)
	}

	r.Body = append(ordinary, postTx...)

	for i := range r.Body {
		frame := &r.Body[i]

		// Only SENDER frames may carry value.
		if frame.Kind != KindTransfer && frame.Kind != KindCall && frame.Kind != KindProbe {
			frame.StateGas = frame.StateGas && frame.Kind != KindPostTx
		}

		// A batched frame may not be the last one, and a POST_TX frame's failure
		// overrides batching anyway, so the flag is meaningless on one.
		if i == len(r.Body)-1 || frame.Kind == KindPostTx {
			frame.Batch = false
		}
	}

	// A POST_TX frame executes as a static call, so it cannot write.
	for i := range r.Body {
		if r.Body[i].Kind == KindPostTx {
			r.Body[i].Script = ScriptNone
			r.Body[i].StateGas = false
		}
	}

	if !opts.AllowProbe {
		for i := range r.Body {
			if r.Body[i].Kind == KindProbe {
				r.Body[i].Kind = KindCall
				r.Body[i].Target = TargetWallet
				r.Body[i].Script = ScriptNone
			}
		}

		r.Reads = false
	}
}

// refusalKey names the aspect of a recipe a refusal is most likely about, so the
// recorded reasons are grouped by something meaningful rather than lumped together.
//
// It is a label, not a claim: several of these shapes are refused by design, and which
// ones a chain ought to refuse is exactly the question this scenario declines to answer
// on its own.
func (r *Recipe) refusalKey() string {
	switch {
	case r.Invalid != "":
		return "invalid:" + r.Invalid
	case r.RecentRootEdge != "":
		return "root-edge:" + r.RecentRootEdge
	case r.RecentRoots > 0:
		return "recent-roots"
	case r.Prefix == PrefixPaymaster:
		return "prefix:" + string(r.Prefix)
	case r.Sender == SenderContract:
		return "contract-sender"
	default:
		return "well-formed"
	}
}
