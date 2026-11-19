package txtypes

import (
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
)

// daisugiDependencyTx is transaction 0x576438a5...477f38 of daisugi-devnet-0, the only
// transaction of execution block 675: a self-relayed frame transaction declaring one
// leanSTARK dependency over its own signature hash.
const daisugiDependencyTx = "0x06f9012b8429d103d8c584ccbd53ec809405a20b6083bd71c58c10cbfa9db8036bf5a8e206f89bce01" +
	"0380c8830186a083017e808080f86b048080c48275308080b8600000000000000000000000000000000000000000000000000000" +
	"00000000001100000000000000000000000000000000000000000000000000000000000000004f0922a927365fec2cf13867145a" +
	"3de76b3d26385681cd100c9b3d23c3c03d46de028094000000000000000000000000000000000000deadc482c350807b80f85cf8" +
	"5a019405a20b6083bd71c58c10cbfa9db8036bf5a8e20680b84101f9e4f70b43dcfddbd1ce7062d64373a45945d6772d68e28c35" +
	"ca070c0fef77137420ebe89837279324ab8195934169a2b9f04542de1ec4eff441471a78c959c6cd8762bcad3f9900008477359" +
	"400c0"

// TestDependencyFrameDevnetVector checks the dependency rules against what the devnet
// clients agreed on for a real block: the transaction is valid, its intrinsic cost is
// the gas the receipt charged outside the frames, and the dependency digest is the
// block_deps_hash of the block's header.
func TestDependencyFrameDevnetVector(t *testing.T) {
	decoded, err := DecodeTx(hexutil.MustDecode(daisugiDependencyTx))
	if err != nil {
		t.Fatalf("failed decoding: %v", err)
	}

	if want := common.HexToHash("0x576438a5eff9cac236779606800c9ed4b172faf3e760557eef527c3739477f38"); decoded.Hash() != want {
		t.Fatalf("hash = %s, want %s", decoded.Hash(), want)
	}

	tx, ok := decoded.Inner().(*FrameTx)
	if !ok {
		t.Fatalf("decoded a %T, want a frame transaction", decoded.Inner())
	}

	if err := tx.ValidatePayload(); err != nil {
		t.Fatalf("payload rejected: %v", err)
	}

	if err := tx.VerifySignatures(); err != nil {
		t.Fatalf("signatures rejected: %v", err)
	}

	if species := tx.Frames[1].Species(tx.Sender); species != SpeciesDepVerify {
		t.Fatalf("frame 1 species = %s, want %s", species, SpeciesDepVerify)
	}

	declared := tx.Frames[1].Dependencies()
	if len(declared) != 1 || !declared[0].SignsTransaction() {
		t.Fatalf("frame 1 declares %+v, want one dependency over the signature hash", declared)
	}

	deps, err := tx.Dependencies()
	if err != nil {
		t.Fatalf("failed resolving dependencies: %v", err)
	}

	want := FrameDependency{
		Scheme:          DepSchemeLeanStark,
		DataHash:        common.HexToHash("0xf1b421aad5d69d2c014c87d8a29e95f8ab93dba1089f6054ea25947de28197b3"),
		VerificationKey: common.HexToHash("0x4f0922a927365fec2cf13867145a3de76b3d26385681cd100c9b3d23c3c03d46"),
	}

	if len(deps) != 1 || deps[0] != want {
		t.Fatalf("dependencies = %+v, want %+v", deps, want)
	}

	// The receipt charged 191,197 gas, of which the frames used 130,620. The devnet
	// runs without EIP-2780, so it did not charge the value transfer of frame 2.
	if gas := tx.IntrinsicGas() - TxValueCost; gas != 191_197-130_620 {
		t.Errorf("intrinsic gas = %d, want %d", gas, 191_197-130_620)
	}

	if gas := tx.RecursiveStarkGas(); gas != LeanStarkVerificationGas {
		t.Errorf("recursive stark gas = %d, want %d", gas, LeanStarkVerificationGas)
	}

	blockDepsHash := common.HexToHash("0x803ea13cfe84fb8a2e3d777f67d71f5c2205f06e001c5e1ac0cc83b6842d2d08")
	if hash := DependenciesHash(deps); hash != blockDepsHash {
		t.Errorf("dependencies hash = %s, want the block_deps_hash %s", hash, blockDepsHash)
	}
}

// TestDependencyFrameValidity checks the validity rules of a dependency frame.
func TestDependencyFrameValidity(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed generating key: %v", err)
	}

	sender := crypto.PubkeyToAddress(key.PublicKey)
	sphincs := FrameDependency{Scheme: DepSchemeLeanSphincs, VerificationKey: common.HexToHash("0x01")}
	stark := FrameDependency{Scheme: DepSchemeLeanStark, DataHash: common.HexToHash("0x02"), VerificationKey: common.HexToHash("0x03")}

	verify := func() *Frame { return SelfVerifyFrame(FrameLimits{Execution: 5_000}) }
	userOp := func() *Frame { return UserOpFrame(&testTarget, nil, nil, FrameLimits{Execution: 30_000}) }
	assert := func() *Frame { return AssertFrame(testTarget, nil, FrameLimits{Execution: 20_000}) }
	dep := func(mutate func(*Frame)) *Frame {
		frame := DependencyFrame(sphincs, stark)
		if mutate != nil {
			mutate(frame)
		}

		return frame
	}

	many := make([]FrameDependency, MaxDependenciesPerFrame+1)
	for i := range many {
		many[i] = sphincs
	}

	for name, tc := range map[string]struct {
		frames []*Frame
		valid  bool
	}{
		"after the validation prefix":  {[]*Frame{verify(), dep(nil), userOp()}, true},
		"before the validation prefix": {[]*Frame{dep(nil), verify(), userOp()}, true},
		"as the last frame":            {[]*Frame{verify(), userOp(), dep(nil)}, true},
		"full frame":                   {[]*Frame{verify(), DependencyFrame(many[:MaxDependenciesPerFrame]...)}, true},
		"with a target":                {[]*Frame{verify(), dep(func(f *Frame) { f.Target = &testTarget })}, false},
		"targeting the sender":         {[]*Frame{verify(), dep(func(f *Frame) { f.Target = &sender })}, false},
		"with value":                   {[]*Frame{verify(), dep(func(f *Frame) { f.Value = uint256.NewInt(1) })}, false},
		"with an approval scope":       {[]*Frame{verify(), dep(func(f *Frame) { f.Flags = ApprovePayment })}, false},
		"opening an atomic batch":      {[]*Frame{verify(), dep(func(f *Frame) { f.Flags = AtomicBatchFlag }), userOp()}, false},
		"closing an atomic batch":      {[]*Frame{verify(), userOp().WithAtomicBatch(), dep(nil)}, false},
		"after an assert frame":        {[]*Frame{verify(), userOp(), assert(), dep(nil)}, false},
		"without dependencies":         {[]*Frame{verify(), DependencyFrame()}, false},
		"too many dependencies":        {[]*Frame{verify(), DependencyFrame(many...)}, false},
		"partial triple":               {[]*Frame{verify(), dep(func(f *Frame) { f.Data = f.Data[:len(f.Data)-1] })}, false},
		"non-zero scheme padding":      {[]*Frame{verify(), dep(func(f *Frame) { f.Data[0] = 1 })}, false},
		"unknown scheme":               {[]*Frame{verify(), dep(func(f *Frame) { f.Data[31] = 0x12 })}, false},
		"state gas budget":             {[]*Frame{verify(), dep(func(f *Frame) { f.Limits.State = 1 })}, false},
		"execution budget too low":     {[]*Frame{verify(), dep(func(f *Frame) { f.Limits.Execution-- })}, false},
		"execution budget too high":    {[]*Frame{verify(), dep(func(f *Frame) { f.Limits.Execution++ })}, false},
		"unknown mode":                 {[]*Frame{verify(), dep(func(f *Frame) { f.Mode = FrameModeDepVerify + 1 })}, false},
	} {
		t.Run(name, func(t *testing.T) {
			tx := NewFrameTx(testChainID, sender, 0, FrameFees{GasFeeCap: uint256.NewInt(1e9)},
				tc.frames, []*FrameSignature{SenderSignature()})

			if err := tx.SignPayload(testChainID.ToBig(), key); err != nil {
				t.Fatalf("failed signing: %v", err)
			}

			err := tx.ValidatePayload()
			if tc.valid && err != nil {
				t.Fatalf("rejected: %v", err)
			}

			if !tc.valid && !errors.Is(err, ErrInvalidFrameTx) {
				t.Fatalf("got %v, want %v", err, ErrInvalidFrameTx)
			}
		})
	}
}

// TestDependencyGas checks what a dependency costs: the frame's budget per scheme, and
// the recursive STARK share in the intrinsic cost and the calldata floor per declared
// triple, duplicates included.
func TestDependencyGas(t *testing.T) {
	sphincs := FrameDependency{Scheme: DepSchemeLeanSphincs, VerificationKey: common.HexToHash("0x01")}
	stark := FrameDependency{Scheme: DepSchemeLeanStark, VerificationKey: common.HexToHash("0x02")}

	frame := DependencyFrame(sphincs, sphincs, stark)
	if want := uint64(2*LeanSphincsVerificationGas + LeanStarkVerificationGas); frame.Limits.Execution != want {
		t.Fatalf("frame budget = %d, want %d", frame.Limits.Execution, want)
	}

	build := func(frames ...*Frame) *FrameTx {
		return NewFrameTx(testChainID, testTarget, 0, FrameFees{GasFeeCap: uint256.NewInt(1e9)},
			frames, []*FrameSignature{SenderSignature()})
	}

	// The same frame as ordinary calldata prices everything but the recursive STARK share.
	plain := *frame
	plain.Mode = FrameModeDefault

	with := build(SelfVerifyFrame(FrameLimits{Execution: 5_000}), frame)
	without := build(SelfVerifyFrame(FrameLimits{Execution: 5_000}), &plain)

	if got := with.RecursiveStarkGas(); got != 3*LeanStarkVerificationGas {
		t.Fatalf("recursive stark gas = %d, want %d", got, 3*LeanStarkVerificationGas)
	}

	if got, want := with.IntrinsicGas(), without.IntrinsicGas()+3*LeanStarkVerificationGas; got != want {
		t.Errorf("intrinsic gas = %d, want %d", got, want)
	}

	if got, want := with.CalldataFloorGas(), without.CalldataFloorGas()+3*LeanStarkVerificationGas; got != want {
		t.Errorf("calldata floor = %d, want %d", got, want)
	}

	if got := without.RecursiveStarkGas(); got != 0 {
		t.Errorf("a transaction without dependency frames pays %d recursive stark gas", got)
	}
}

// TestDependenciesCanonicalForm checks the list a transaction's dependencies resolve
// to, and the digest taken over it.
func TestDependenciesCanonicalForm(t *testing.T) {
	a := FrameDependency{Scheme: DepSchemeLeanSphincs, DataHash: common.HexToHash("0x0a"), VerificationKey: common.HexToHash("0x01")}
	b := FrameDependency{Scheme: DepSchemeLeanSphincs, DataHash: common.HexToHash("0x0b"), VerificationKey: common.HexToHash("0x01")}
	c := FrameDependency{Scheme: DepSchemeLeanStark, DataHash: common.HexToHash("0x01"), VerificationKey: common.HexToHash("0x01")}
	own := FrameDependency{Scheme: DepSchemeLeanSphincs, VerificationKey: common.HexToHash("0x02")}

	tx := NewFrameTx(testChainID, testTarget, 0, FrameFees{GasFeeCap: uint256.NewInt(1e9)},
		[]*Frame{
			DependencyFrame(c, b, own),
			SelfVerifyFrame(FrameLimits{Execution: 5_000}),
			DependencyFrame(a, b, own),
		},
		[]*FrameSignature{SenderSignature()})

	if count := tx.DependencyCount(); count != 6 {
		t.Fatalf("dependency count = %d, want the 6 declared triples", count)
	}

	sigHash, err := tx.SigHash()
	if err != nil {
		t.Fatalf("failed hashing: %v", err)
	}

	resolved := own
	resolved.DataHash = sigHash

	// Sorted by scheme, then data hash: the signature hash sorts wherever its bytes put it.
	want := CanonicalDependencies([]FrameDependency{a, b, resolved, c})
	if want[len(want)-1] != c {
		t.Fatalf("the leanSTARK dependency sorts at %v, want it last", want)
	}

	deps, err := tx.Dependencies()
	if err != nil {
		t.Fatalf("failed resolving dependencies: %v", err)
	}

	if len(deps) != len(want) {
		t.Fatalf("dependencies = %+v, want %+v", deps, want)
	}

	for i := range want {
		if deps[i] != want[i] {
			t.Fatalf("dependency %d = %+v, want %+v", i, deps[i], want[i])
		}
	}

	// The frames still carry the zero data hash.
	if declared := tx.Frames[0].Dependencies(); !declared[2].SignsTransaction() {
		t.Error("resolving the dependencies rewrote the frame's data")
	}

	// A leanSPHINCS list hashes as the EIP specifies, in canonical order whatever order
	// it is given in.
	sphincsOnly := crypto.Keccak256Hash(a.DataHash[:], a.VerificationKey[:], b.DataHash[:], b.VerificationKey[:])
	if hash := DependenciesHash([]FrameDependency{b, a, b}); hash != sphincsOnly {
		t.Errorf("leanSPHINCS digest = %s, want %s", hash, sphincsOnly)
	}

	if hash := DependenciesHash(nil); hash != crypto.Keccak256Hash() {
		t.Errorf("empty digest = %s, want the hash of nothing", hash)
	}

	if DependenciesHash([]FrameDependency{a, c}) == DependenciesHash([]FrameDependency{a}) {
		t.Error("a leanSTARK dependency does not change the digest")
	}
}

// TestDependencyMempoolPolicy checks that the public mempool rules look past dependency
// frames when matching the validation prefix, and enforce the per-transaction limits.
func TestDependencyMempoolPolicy(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed generating key: %v", err)
	}

	sender := crypto.PubkeyToAddress(key.PublicKey)
	paymaster := common.HexToAddress("0xeeee000000000000000000000000000000000007")

	sphincs := func(n byte) FrameDependency {
		return FrameDependency{Scheme: DepSchemeLeanSphincs, VerificationKey: common.BytesToHash([]byte{n})}
	}
	stark := func(n byte) FrameDependency {
		return FrameDependency{Scheme: DepSchemeLeanStark, VerificationKey: common.BytesToHash([]byte{n})}
	}

	sigs := make([]FrameDependency, MaxSigsPerTx+1)
	for i := range sigs {
		sigs[i] = sphincs(byte(i + 1))
	}

	selfVerify := func() *Frame { return SelfVerifyFrame(FrameLimits{Execution: 5_000}) }
	onlyVerify := func() *Frame { return OnlyVerifyFrame(FrameLimits{Execution: 5_000}) }
	pay := func() *Frame { return PayFrame(paymaster, nil, FrameLimits{Execution: 5_000}) }
	userOp := func() *Frame { return UserOpFrame(&testTarget, nil, nil, FrameLimits{Execution: 30_000}) }
	expiry := func() *Frame { return ExpiryFrame(1, 5_000) }
	dep := func(deps ...FrameDependency) *Frame { return DependencyFrame(deps...) }

	// A full frame of leanSTARK dependencies budgets far more than the verification gas
	// cap, which it must not count against.
	repeated := make([]FrameDependency, MaxDependenciesPerFrame)
	for i := range repeated {
		repeated[i] = stark(1)
	}

	for name, tc := range map[string]struct {
		frames []*Frame
		valid  bool
	}{
		"no dependency frame":                   {[]*Frame{selfVerify(), userOp()}, true},
		"dependencies lead the prefix":          {[]*Frame{dep(sphincs(1)), selfVerify(), userOp()}, true},
		"dependencies follow the prefix":        {[]*Frame{selfVerify(), dep(sphincs(1)), userOp()}, true},
		"dependencies between verify and pay":   {[]*Frame{onlyVerify(), dep(sphincs(1)), pay(), userOp()}, true},
		"dependencies around the expiry frame":  {[]*Frame{dep(sphincs(1)), expiry(), dep(sphincs(2)), selfVerify()}, true},
		"expiry frame behind a verify frame":    {[]*Frame{dep(sphincs(1)), selfVerify(), expiry()}, false},
		"dependency budget outside verify gas":  {[]*Frame{dep(repeated...), selfVerify(), userOp()}, true},
		"signature limit":                       {[]*Frame{selfVerify(), dep(sigs[:MaxSigsPerTx]...)}, true},
		"over the signature limit":              {[]*Frame{selfVerify(), dep(sigs...)}, false},
		"over the signature limit, two frames":  {[]*Frame{dep(sigs[:1]...), selfVerify(), dep(sigs[1:]...)}, false},
		"one stark":                             {[]*Frame{selfVerify(), dep(stark(1), sphincs(1))}, true},
		"two starks":                            {[]*Frame{selfVerify(), dep(stark(1), stark(2))}, false},
		"repeated stark counts once":            {[]*Frame{selfVerify(), dep(stark(1), stark(1))}, true},
		"unrecognized prefix with dependencies": {[]*Frame{dep(sphincs(1)), pay(), userOp()}, false},
	} {
		t.Run(name, func(t *testing.T) {
			tx := NewFrameTx(testChainID, sender, 0, FrameFees{GasFeeCap: uint256.NewInt(1e9)},
				tc.frames, []*FrameSignature{SenderSignature()})

			if err := tx.SignPayload(testChainID.ToBig(), key); err != nil {
				t.Fatalf("failed signing: %v", err)
			}

			if err := tx.ValidatePayload(); err != nil {
				t.Fatalf("payload rejected: %v", err)
			}

			err := tx.ValidateMempoolPrefix()
			if tc.valid && err != nil {
				t.Fatalf("rejected: %v", err)
			}

			if !tc.valid && !errors.Is(err, ErrMempoolPolicy) {
				t.Fatalf("got %v, want %v", err, ErrMempoolPolicy)
			}
		})
	}
}
