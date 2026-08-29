package frametxfuzz

import (
	"encoding/json"
	"fmt"

	"github.com/ethpandaops/spamoor/utils"
)

// Recipe is the abstract description of one generated transaction.
//
// It is drawn in a single pass and holds no chain state: the sending wallet, the nonce
// and the probe contract's address are applied afterwards. That is what makes a seed and
// index reproduce the same transaction on any chain.
type Recipe struct {
	// Index is the transaction index the recipe was drawn for.
	Index uint64 `json:"index"`

	// Prefix is the validation prefix shape.
	Prefix PrefixShape `json:"prefix"`

	// Expiry adds a leading expiry verifier frame, which the mempool rules skip when
	// matching the prefix against the recognized shapes.
	Expiry bool `json:"expiry"`

	// Sender selects what code the sender's validation frame runs.
	Sender SenderKind `json:"sender"`

	// FuzzedPaymaster sponsors through an account delegated to generated code instead
	// of the fixed one.
	FuzzedPaymaster bool `json:"fuzzedPaymaster,omitempty"`

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

	// Invalid names a deliberate violation the transaction carries, if any. Such a
	// transaction is sent from a burner wallet outside the managed pool, since one that
	// never lands would stall its sender's nonce.
	Invalid string `json:"invalid,omitempty"`
}

// PrefixShape is one of the validation prefixes the public mempool recognizes.
type PrefixShape string

// The recognized prefixes. The two deploy-led shapes are absent: a CREATE2 account
// deployment costs far more than the prefix's execution cap allows, and such a
// transaction's sender has no key for the pool to track.
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
	// deployed code rather than from the protocol's default code.
	SenderContract SenderKind = "contract"

	// SenderFuzzedContract is a generated account contract an earlier transaction
	// deployed, so the validation frame runs arbitrary code before it can approve
	// anything and the transaction carries no signature at all. Drawn rarely: a fuzzed
	// prologue often halts before it reaches the approval.
	SenderFuzzedContract SenderKind = "fuzzed_contract"
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

	// KindDeployCode is a SENDER frame that deploys generated contract code through
	// the CREATE2 factory, so a later frame in the same transaction can call it.
	KindDeployCode FrameKind = "deploy_code"
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

	// TargetDeployed is a contract a frame deployed: one from this transaction when it
	// has one, otherwise one an earlier transaction left behind.
	TargetDeployed TargetKind = "deployed"
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

	// NoApprove leaves the APPROVE off a deployed account's code, so it can never play
	// the sender or paymaster role successfully.
	NoApprove bool `json:"noApprove,omitempty"`
}

// recentRootEdgeNames are the awkward reference cases worth reaching.
//
// All but "duplicate" are refused by design, so they are drawn rarely: a recipe with no
// edge names a root the run committed a slot or more ago, which is the case that has to
// keep landing for the rest of the reference machinery to be exercised at all.
var recentRootEdgeNames = []string{
	"same_slot",    // current_slot - slot == 0
	"future_slot",  // slot ahead of the current one
	"unwritten",    // a slot the source never wrote
	"wrong_source", // a correct root under a source that did not write it
	"duplicate",    // the same reference twice, which is valid and charged twice
	"outside_window",
}

// recentRootEdgeChance is how often a recent-root recipe carries one of the awkward
// cases rather than a plain, landable reference.
const recentRootEdgeChance = 0.2

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
	axisCode      axis = "code"
)

// axisNames lists the axes in a stable order for help text and validation.
var axisNames = []axis{
	axisPrefix, axisBatches, axisFailures, axisSignature,
	axisNonces, axisRoots, axisPostTx, axisProbe, axisCode,
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
	Axes               axisWeights
	MaxBodyFrames      int
	AllowPostTx        bool
	AllowContract      bool
	AllowRoots         bool
	AllowKeyed         bool
	AllowProbe         bool
	AllowCode          bool
	AllowFuzzedAccount bool
	// InvalidChance is how often a drawn recipe carries a deliberate violation.
	InvalidChance float64

	// Violations is the catalog a violation is drawn from.
	Violations []string
}

// Draw produces the recipe for one transaction index.
//
// Every field is drawn unconditionally and only then gated by the options, so the stream
// position does not depend on the configuration: two runs differing in one flag disagree
// about that axis rather than about every later transaction.
func Draw(rng *utils.DeterministicRNG, index uint64, opts DrawOptions) *Recipe {
	recipe := &Recipe{Index: index}

	prefixDraw := rng.Float64()
	paymasterDraw := rng.Float64()
	expiryDraw := rng.Float64()
	senderDraw := rng.Float64()
	witnessDraw := rng.Float64()
	p256Draw := rng.Float64()
	assertDraw := rng.Float64()

	nonceKeyDraw := rng.Intn(4) + 1
	nonceDraw := rng.Float64()

	rootCountDraw := rng.Intn(3) + 1
	rootEdgePick := rng.Intn(len(recentRootEdgeNames))
	rootEdgeDraw := rng.Float64()
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

	// One draw, split into bands: the fixed contract sender is common and the fuzzed one
	// is rare, since a validation frame running generated code rarely approves.
	recipe.Sender = SenderDefaultCode

	switch {
	case opts.AllowFuzzedAccount && senderDraw < 0.05:
		recipe.Sender = SenderFuzzedContract
	case opts.AllowContract && senderDraw < 0.25:
		recipe.Sender = SenderContract
	}

	recipe.FuzzedPaymaster = opts.AllowFuzzedAccount && recipe.Prefix == PrefixPaymaster && paymasterDraw < 0.1

	recipe.Witness = opts.Axes.enabled(axisSignature) && witnessDraw < 0.4*opts.Axes.chance(axisSignature)
	recipe.P256 = opts.Axes.enabled(axisSignature) && p256Draw < 0.2*opts.Axes.chance(axisSignature)

	if opts.AllowKeyed && opts.Axes.draws(rng, axisNonces) && nonceDraw < 0.6 {
		recipe.NonceKeys = nonceKeyDraw
	}

	if opts.AllowRoots && opts.Axes.draws(rng, axisRoots) && rootDraw < 0.6 {
		recipe.RecentRoots = rootCountDraw

		if rootEdgeDraw < recentRootEdgeChance {
			recipe.RecentRootEdge = recentRootEdgeNames[rootEdgePick]
		}
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
	case kindDraw < 0.24:
		frame.Kind = KindTransfer
	case kindDraw < 0.48 && opts.AllowProbe && opts.Axes.enabled(axisProbe):
		frame.Kind = KindProbe
		frame.Target = TargetProbe
	case kindDraw < 0.62 && opts.AllowCode && opts.Axes.enabled(axisCode):
		frame.Kind = KindDeployCode
	case kindDraw < 0.70:
		frame.Kind = KindPostOp
	}

	if frame.Kind != KindProbe && frame.Kind != KindDeployCode {
		switch {
		case targetDraw < 0.15:
			frame.Target = TargetSender
		case targetDraw < 0.30:
			frame.Target = TargetEmpty
			frame.StateGas = true
		case targetDraw < 0.38:
			frame.Target = TargetPrecompile
		case targetDraw < 0.62 && opts.AllowCode && opts.Axes.enabled(axisCode):
			// Calling generated code is the point of deploying it, so it takes a
			// substantial share of the targets rather than a token one.
			frame.Target = TargetDeployed
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

	if !opts.AllowFuzzedAccount {
		if r.Sender == SenderFuzzedContract {
			r.Sender = SenderDefaultCode
		}

		r.FuzzedPaymaster = false
	}

	// An account contract is funded for a transaction of ordinary size, and it is
	// charged the maximum cost up front. Deploy frames are what make a transaction
	// expensive, so a recipe leaning on one keeps its body cheap.
	if r.Sender == SenderFuzzedContract || r.FuzzedPaymaster {
		for i := range r.Body {
			if r.Body[i].Kind == KindDeployCode {
				r.Body[i].Kind = KindCall
				r.Body[i].Target = TargetWallet
			}
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

	if !opts.AllowCode {
		for i := range r.Body {
			if r.Body[i].Kind == KindDeployCode {
				r.Body[i].Kind = KindCall
			}

			if r.Body[i].Target == TargetDeployed {
				r.Body[i].Target = TargetWallet
			}
		}
	}

	// A deploy frame writes code, so it cannot be a static POST_TX frame, and a frame
	// that deploys has nothing to say about a value transfer.
	for i := range r.Body {
		if r.Body[i].Kind == KindDeployCode {
			r.Body[i].Target = TargetWallet
			r.Body[i].Script = ScriptNone
			r.Body[i].StateGas = false
		}
	}

	if opts.AllowCode {
		r.pairDeployWithCall()
	}
}

// refusalKey names the aspect of a recipe a refusal is most likely about, so recorded
// reasons are grouped rather than lumped together. It is a label, not a claim.
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
	case r.Sender == SenderFuzzedContract:
		return "fuzzed-sender"
	case r.FuzzedPaymaster:
		return "fuzzed-paymaster"
	default:
		return "well-formed"
	}
}

// pairDeployWithCall promotes the first frame to a deployment when a later frame calls
// generated code without one, so both halves land in the same frame list.
//
// A call in the very first frame is left alone: there is nowhere to put a deployment
// ahead of it, and it falls back to an earlier transaction's contract.
func (r *Recipe) pairDeployWithCall() {
	deployIndex := -1

	for i := range r.Body {
		if r.Body[i].Kind == KindDeployCode {
			deployIndex = i

			break
		}
	}

	for i := range r.Body {
		if r.Body[i].Target != TargetDeployed || i == 0 {
			continue
		}

		if deployIndex >= 0 && deployIndex < i {
			return
		}

		r.Body[0].Kind = KindDeployCode
		r.Body[0].Target = TargetWallet
		r.Body[0].Script = ScriptNone
		r.Body[0].StateGas = false
		r.Body[0].Batch = false

		return
	}
}
