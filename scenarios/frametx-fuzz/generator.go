package frametxfuzz

import (
	"context"
	"crypto/ecdsa"
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"

	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txtypes"
)

// identityPrecompile is a target every chain has.
var identityPrecompile = common.BytesToAddress([]byte{0x04})

// deployGasPerCodeByte is the execution gas budgeted per byte of init code, covering the
// copy, the CREATE2 and the code deposit.
const deployGasPerCodeByte = 1_600

// build is a recipe turned into a transaction, together with the dimensions it covers.
type build struct {
	recipe *Recipe

	tx     *txtypes.FrameTx
	sender *spamoor.Wallet
	p256   *ecdsa.PrivateKey
	frames []*txtypes.Frame

	// prefixLen is how many leading frames make up the validation prefix.
	prefixLen int

	// nonces is the keyed nonce selection this transaction consumes on success.
	nonces *selection

	// mempoolLegal records whether the transaction is expected to propagate.
	mempoolLegal bool

	// coverage names the dimensions this transaction exercises, which is what the run
	// reports: the point is to trigger combinations, not to predict their outcome.
	coverage []string

	// deployed are the generated contracts this transaction's frames create, in frame
	// order. A CREATE2 address is known before the transaction is sent, so a later
	// frame can call code the transaction has not deployed yet.
	deployed []common.Address
}

// buildRecipe assembles the transaction a recipe describes.
func (s *Scenario) buildRecipe(ctx context.Context, client *spamoor.Client, env *environment, recipe *Recipe, feeCap, tipCap *big.Int) (*build, error) {
	sender := env.senderFor(recipe, int(recipe.Index))
	if sender == nil {
		return nil, fmt.Errorf("no wallet available")
	}

	result := &build{recipe: recipe, sender: sender, mempoolLegal: true}

	if recipe.Reads {
		result.cover("introspection-reads")
	}

	if recipe.Invalid != "" {
		result.cover("invalid:" + recipe.Invalid)
	}

	if err := s.buildPrefix(env, recipe, result); err != nil {
		return nil, err
	}

	if err := s.buildBody(env, recipe, result); err != nil {
		return nil, err
	}

	signatures, err := s.buildSignatures(recipe, result)
	if err != nil {
		return nil, err
	}

	result.tx = txtypes.NewFrameTxWithExtensions(env.extensions, nil, sender.GetAddress(), 0,
		txtypes.FrameFees{
			GasTipCap: uint256.MustFromBig(tipCap),
			GasFeeCap: uint256.MustFromBig(feeCap),
		},
		result.frames,
		signatures,
	)

	if err := s.applyNonceKeys(ctx, client, env, recipe, result); err != nil {
		return nil, err
	}

	if err := s.applyRecentRoots(env, recipe, result); err != nil {
		return nil, err
	}

	return result, nil
}

// buildPrefix assembles the validation prefix.
//
// The prefix is the part the public mempool simulates, so everything in it is sized to
// stay inside the verification gas caps and to match one of the recognized shapes.
func (s *Scenario) buildPrefix(env *environment, recipe *Recipe, result *build) error {
	if recipe.Expiry {
		// The expiry frame may lead the frame list and nothing else, and is skipped
		// when the prefix is matched against the recognized shapes.
		result.append(txtypes.ExpiryFrame(s.deadline(), s.options.VerifyGas))
	}

	contractSender := recipe.Sender == SenderContract && env.probe != nil && env.contractCount > 0

	result.cover("prefix:" + string(recipe.Prefix))

	if recipe.Expiry {
		result.cover("expiry-frame")
	}

	if contractSender {
		result.cover("contract-sender")
	}

	switch recipe.Prefix {
	case PrefixPaymaster:
		if env.probe == nil {
			// Without the probe contract there is nothing to play the paymaster, so
			// the recipe falls back to the self-relayed shape rather than building a
			// transaction that cannot validate.
			return s.buildSelfVerifyPrefix(env, contractSender, result)
		}

		verify := txtypes.OnlyVerifyFrame(txtypes.FrameLimits{Execution: s.options.VerifyGas})
		if contractSender {
			script := NewProbeScript().Approve(txtypes.ApproveExecution)
			verify.Data = script.Bytes()
			verify.Limits.Execution = script.ExecutionGas() + probeEntryGas
		}

		payScript := NewProbeScript().Approve(txtypes.ApprovePayment)

		result.append(verify)
		result.append(txtypes.PayFrame(env.paymasterFor(int(recipe.Index), result.sender.GetAddress()),
			payScript.Bytes(), txtypes.FrameLimits{
				Execution: payScript.ExecutionGas() + probeEntryGas,
			}))

	default:
		return s.buildSelfVerifyPrefix(env, contractSender, result)
	}

	result.prefixLen = len(result.frames)

	return nil
}

// buildSelfVerifyPrefix appends the single-frame validation prefix.
func (s *Scenario) buildSelfVerifyPrefix(env *environment, contractSender bool, result *build) error {
	verify := txtypes.SelfVerifyFrame(txtypes.FrameLimits{Execution: s.options.VerifyGas})

	if contractSender {
		script := NewProbeScript().Approve(txtypes.ApproveExecutionAndPayment)
		verify.Data = script.Bytes()
		verify.Limits.Execution = script.ExecutionGas() + probeEntryGas
	}

	result.append(verify)
	result.prefixLen = len(result.frames)

	return nil
}

// buildBody assembles the frames after the validation prefix.
func (s *Scenario) buildBody(env *environment, recipe *Recipe, result *build) error {
	for i, spec := range recipe.Body {
		frame, err := s.buildBodyFrame(env, recipe, result, spec, i)
		if err != nil {
			return err
		}

		if spec.Batch {
			frame.WithAtomicBatch()
			result.cover("batch")
		}

		result.cover("frame:" + string(spec.Kind))
		result.cover("target:" + string(spec.Target))

		if spec.Budget == BudgetStarved {
			result.cover("starved-frame")
		}

		if spec.Script != ScriptNone {
			result.cover("script:" + string(spec.Script))
		}

		result.append(frame)
	}

	if len(recipe.Body) == 0 {
		result.append(txtypes.UserOpFrame(nil, nil, nil, txtypes.FrameLimits{
			Execution: s.options.UserOpGas,
		}))
	}

	return nil
}

// buildBodyFrame assembles one body frame and says what it is built to do.
func (s *Scenario) buildBodyFrame(env *environment, recipe *Recipe, result *build, spec BodyFrame, index int) (*txtypes.Frame, error) {
	if spec.Kind == KindDeployCode {
		return s.buildDeployFrame(env, recipe, result, index)
	}

	target, targetEmpty := s.resolveTarget(env, recipe, result, spec, index)

	script, err := s.buildScript(env, recipe, spec, index)
	if err != nil {
		return nil, err
	}

	limits := txtypes.FrameLimits{Execution: s.options.UserOpGas}
	if spec.Target == TargetDeployed {
		// Generated code is worth giving room to run: at the ordinary per-frame budget
		// almost every call would end on gas before it reached anything interesting.
		limits.Execution = s.options.CodeGas
	}

	if script != nil {
		limits.Execution = script.ExecutionGas() + probeEntryGas
		limits.State = script.StateGas()
	}

	if spec.StateGas && targetEmpty {
		limits.State += txtypes.StateBytesPerNewAccount * txtypes.CostPerStateByte
	}

	limits.State += s.options.StateGas

	// One gas cannot cover the frame-entry account access charge, which the EIP takes
	// before the target's code runs. Whether a client honours that is exactly the kind
	// of thing worth putting in front of a network rather than asserting here.
	if spec.Budget == BudgetStarved {
		limits.Execution = 1
	}

	var data []byte
	if script != nil {
		data = script.Bytes()
	} else {
		data = s.callData
	}

	switch spec.Kind {
	case KindTransfer:
		return txtypes.UserOpFrame(&target, s.amount(recipe), data, limits), nil
	case KindPostOp:
		return txtypes.PostOpFrame(target, data, limits), nil
	case KindPostTx:
		if !env.allowPostTx {
			// The chain refused a POST_TX probe at startup, so the frame is downgraded
			// rather than generating transactions that cannot land.
			return txtypes.UserOpFrame(&target, nil, data, limits), nil
		}

		return txtypes.PostTxFrame(target, data, limits), nil
	default:
		return txtypes.UserOpFrame(&target, nil, data, limits), nil
	}
}

// buildDeployFrame assembles a frame that deploys generated contract code.
//
// The frame calls the CREATE2 factory, whose calldata is a salt followed by init code, so
// the address is a function of the code and is known here rather than after the fact. That
// is what lets a later frame in the same transaction call a contract this one has not
// created yet.
func (s *Scenario) buildDeployFrame(env *environment, recipe *Recipe, result *build, index int) (*txtypes.Frame, error) {
	runtime := GenerateContractCode(s.seed, recipe.Index, index, int(s.options.MaxCodeSize), s.options.CodeGas)

	initcode, err := DeploymentInitCode(runtime)
	if err != nil {
		return nil, err
	}

	salt := deploySalt(recipe.Index, index)

	data := make([]byte, 0, len(salt)+len(initcode))
	data = append(data, salt[:]...)
	data = append(data, initcode...)

	result.deployed = append(result.deployed, crypto.CreateAddress2(env.factory, salt, crypto.Keccak256(initcode)))
	result.cover("deploy-code")

	factory := env.factory

	return txtypes.UserOpFrame(&factory, nil, data, txtypes.FrameLimits{
		Execution: deployExecutionGas(s.options.CodeGas, len(initcode)),
		State:     deployStateGas(len(initcode)),
	}), nil
}

// deployExecutionGas budgets a deployment's execution gas, which scales with the code
// being deployed rather than being flat.
func deployExecutionGas(base uint64, codeSize int) uint64 {
	return base + uint64(codeSize)*deployGasPerCodeByte
}

// deployStateGas budgets a deployment's EIP-8037 state gas: an account leaf plus a byte
// per byte of code.
//
// Sized against the init code, which is longer than the deployed runtime, with headroom.
// A deploy frame that cannot cover the state charge fails invisibly -- CREATE2 returns
// zero and the factory reverts -- taking everything downstream of the deployment with it.
func deployStateGas(codeSize int) uint64 {
	return uint64(txtypes.StateBytesPerNewAccount+codeSize) * txtypes.CostPerStateByte * 5 / 4
}

// deploySalt derives the CREATE2 salt for one deploy frame, so a replayed recipe deploys
// to the same address.
func deploySalt(txIndex uint64, frameIndex int) [32]byte {
	var salt [32]byte

	copy(salt[:], []byte("frametx-fuzz-code"))
	binary.BigEndian.PutUint64(salt[20:28], txIndex)
	binary.BigEndian.PutUint32(salt[28:32], uint32(frameIndex))

	return salt
}

// resolveTarget turns a recipe's target kind into an address, reporting whether it is
// an account that does not exist yet.
func (s *Scenario) resolveTarget(env *environment, recipe *Recipe, result *build, spec BodyFrame, index int) (common.Address, bool) {
	switch spec.Target {
	case TargetDeployed:
		// Prefer code this transaction is deploying, so the call and the deployment sit
		// in the same frame list; otherwise reach for something an earlier transaction
		// left behind.
		if len(result.deployed) > 0 {
			result.cover("call-own-deployed")

			return result.deployed[index%len(result.deployed)], false
		}

		if address, ok := env.pickContract(index); ok {
			result.cover("call-foreign-deployed")

			return address, false
		}
	case TargetSender:
		return env.senderFor(recipe, int(recipe.Index)).GetAddress(), false
	case TargetProbe:
		if env.probe != nil {
			return env.probe.Address, false
		}
	case TargetPrecompile:
		return identityPrecompile, false
	case TargetEmpty:
		// Derived from the recipe rather than drawn at random, so a replayed recipe
		// addresses the same account and costs the same state gas.
		word := uint256.NewInt(recipe.Index).Bytes32()

		seed := append([]byte("frametx-fuzz-empty"), byte(index))
		seed = append(seed, word[:]...)

		return common.BytesToAddress(crypto.Keccak256(seed)[12:]), true
	}

	return env.targetWallet(int(recipe.Index) + index), false
}

// buildScript assembles a probe frame's script, or nil when the frame is not a probe
// frame.
//
// A script does two things: whatever the recipe asked for, and, when the recipe enables
// them, the read operations that make the instructions EIP-8141 introduces execute
// inside a frame. The reads discard their results -- the point is to reach the
// instruction, not to decide what it should have returned.
func (s *Scenario) buildScript(env *environment, recipe *Recipe, spec BodyFrame, index int) (*ProbeScript, error) {
	if spec.Kind != KindProbe || env.probe == nil {
		return nil, nil
	}

	script := NewProbeScript()

	if recipe.Reads {
		appendReads(script, recipe, index)
	}

	switch spec.Script {
	case ScriptRevert:
		script.Revert()
	case ScriptLog:
		script.Log(common.HexToHash("0x10c0"), 64)
	case ScriptStore:
		script.SStore(common.HexToHash("0x01"), common.HexToHash("0x02"))
	}

	return script, nil
}

// buildSignatures assembles the signature list.
func (s *Scenario) buildSignatures(recipe *Recipe, result *build) ([]*txtypes.FrameSignature, error) {
	signatures := []*txtypes.FrameSignature{txtypes.SenderSignature()}

	if recipe.Witness {
		signatures = append(signatures, txtypes.ArbitrarySignature([]byte("spamoor frame fuzz witness")))
		result.cover("arbitrary-witness")
	}

	if recipe.P256 {
		p256Key, err := generateP256Key()
		if err != nil {
			return nil, err
		}

		result.p256 = p256Key
		signatures = append(signatures, &txtypes.FrameSignature{Scheme: txtypes.SigSchemeP256})
		result.cover("p256-signature")
	}

	return signatures, nil
}

// applyNonceKeys selects an EIP-8250 key set for the transaction.
func (s *Scenario) applyNonceKeys(ctx context.Context, client *spamoor.Client, env *environment, recipe *Recipe, result *build) error {
	if recipe.NonceKeys == 0 || env.nonces == nil || !env.extensions.Has(txtypes.FrameExtKeyedNonces) {
		return nil
	}

	// A key's first use is charged out of the frame executing APPROVE, and that frame
	// is inside the validation prefix, so what fits depends on what the rest of the
	// prefix already costs. Deriving the headroom rather than assuming a fixed budget
	// is what keeps a transaction that also carries a P256 entry or a sponsored prefix
	// from tipping over the cap.
	sel, err := env.nonces.selectKeys(ctx, client, result.sender.GetAddress(), recipe.NonceKeys,
		s.firstUseHeadroom(result))
	if err != nil {
		return err
	}

	if sel == nil || len(sel.keys) == 0 {
		return nil
	}

	result.nonces = sel
	result.tx.WithNonceKeys(sel.keys, sel.sequence)
	result.cover("keyed-nonces")

	if sel.firstUses > 0 {
		result.cover("keyed-nonce-first-use")
	}

	// The frame executing APPROVE has to budget for the surcharge, or the approval
	// halts out of gas.
	if sel.firstUses > 0 {
		result.frames[result.prefixLen-1].Limits.Execution += uint64(sel.firstUses) * txtypes.KeyedNonceFirstUseGas
	}

	return nil
}

// firstUseHeadroom returns how many never-before-used nonce keys still fit inside the
// public mempool's verification gas cap, given what this transaction's validation prefix
// and signature list already cost.
func (s *Scenario) firstUseHeadroom(result *build) int {
	used, err := result.tx.SignatureVerificationGas()
	if err != nil {
		return 0
	}

	for i := 0; i < result.prefixLen && i < len(result.frames); i++ {
		used += result.frames[i].Limits.Execution
	}

	if used >= txtypes.MaxVerifyGas {
		return 0
	}

	return int((txtypes.MaxVerifyGas - used) / txtypes.KeyedNonceFirstUseGas)
}

// applyRecentRoots declares EIP-8272 references on the transaction.
func (s *Scenario) applyRecentRoots(env *environment, recipe *Recipe, result *build) error {
	if recipe.RecentRoots == 0 || env.roots == nil || !env.extensions.Has(txtypes.FrameExtRecentRoots) {
		return nil
	}

	if !env.roots.calibratedClock() {
		return nil
	}

	// One slot of margin: the reference has to be usable in the block the transaction
	// lands in, not in the one the clock says we are in now, and a root committed in
	// slot S only becomes referenceable in S+1.
	current := s.currentSlot()
	if current == 0 {
		return nil
	}

	references, legal := env.roots.references(recipe, current-1)
	if len(references) == 0 {
		return nil
	}

	result.tx.RecentRoots = references
	result.mempoolLegal = result.mempoolLegal && legal
	result.cover("recent-roots")

	if recipe.RecentRootEdge != "" {
		result.cover("root-edge:" + recipe.RecentRootEdge)
	}

	return nil
}

// append adds a frame.
func (b *build) append(frame *txtypes.Frame) {
	b.frames = append(b.frames, frame)
}

// cover records that the transaction exercises a dimension.
func (b *build) cover(name string) {
	b.coverage = append(b.coverage, name)
}
