package calltx

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const (
	testFactory  = "0x4e59b44847b379578588920cA78FbF26c0B4956C"
	testInitCode = "0x60006001f3"
)

func testOptions(patterns ...Create2Pattern) TargetPoolOptions {
	return TargetPoolOptions{Create2Patterns: patterns}
}

func TestBuildAttackTargetsResolvesPatterns(t *testing.T) {
	options := testOptions(Create2Pattern{
		Name:      "minimal-contracts",
		Factory:   testFactory,
		InitCode:  testInitCode,
		StartSalt: 7,
		Count:     2000,
	})

	targets, ignored, err := buildAttackTargets(options, ".", 0)
	if err != nil {
		t.Fatalf("buildAttackTargets returned error: %v", err)
	}
	if len(ignored) != 0 {
		t.Errorf("ignored fields = %v, want none", ignored)
	}
	if got, want := len(targets.patterns), 1; got != want {
		t.Fatalf("pattern count = %d, want %d", got, want)
	}

	pattern := targets.patterns[0]
	if got, want := pattern.Name, "minimal-contracts"; got != want {
		t.Errorf("name = %q, want %q", got, want)
	}
	if got, want := pattern.Factory, common.HexToAddress(testFactory); got != want {
		t.Errorf("factory = %s, want %s", got, want)
	}
	if got, want := pattern.InitCodeHash, crypto.Keccak256Hash(common.FromHex(testInitCode)); got != want {
		t.Errorf("initCodeHash = %s, want %s", got, want)
	}
	if pattern.StartSalt != 7 || pattern.Count != 2000 {
		t.Errorf("salt range = %d..%d, want 7..2006", pattern.StartSalt, pattern.StartSalt+pattern.Count-1)
	}
}

func TestBuildAttackTargetsNamesPatternsByIndex(t *testing.T) {
	options := testOptions(
		Create2Pattern{Factory: testFactory, InitCode: testInitCode, Count: 1},
		Create2Pattern{Factory: testFactory, InitCodeHash: crypto.Keccak256Hash([]byte("other")).Hex(), Count: 1},
	)

	targets, _, err := buildAttackTargets(options, ".", 0)
	if err != nil {
		t.Fatalf("buildAttackTargets returned error: %v", err)
	}
	for i, want := range []string{"create2_1", "create2_2"} {
		if got := targets.patterns[i].Name; got != want {
			t.Errorf("pattern %d name = %q, want %q", i, got, want)
		}
	}
}

// The calldata layout is the contract ABI: the offsets asserted here are the
// ones ContractCallAttack.geas reads with CALLDATALOAD.
func TestBuildCallDataMatchesContractLayout(t *testing.T) {
	options := testOptions(Create2Pattern{
		Name:      "minimal-contracts",
		Factory:   testFactory,
		InitCode:  testInitCode,
		StartSalt: 11,
		Count:     100,
	})
	targets, _, err := buildAttackTargets(options, ".", 0)
	if err != nil {
		t.Fatalf("buildAttackTargets returned error: %v", err)
	}

	data, pattern, startSalt := targets.buildCallData(0, 1, 50000)
	if got, want := len(data), 4+5*32; got != want {
		t.Fatalf("calldata length = %d, want %d", got, want)
	}
	if got, want := data[:4], crypto.Keccak256([]byte(CallAttackFnSig))[:4]; !bytes.Equal(got, want) {
		t.Errorf("selector = 0x%x, want 0x%x", got, want)
	}
	if got, want := common.BytesToAddress(data[0x04:0x24]), common.HexToAddress(testFactory); got != want {
		t.Errorf("factory at 0x04 = %s, want %s", got, want)
	}
	if got, want := common.BytesToHash(data[0x24:0x44]), pattern.InitCodeHash; got != want {
		t.Errorf("initCodeHash at 0x24 = %s, want %s", got, want)
	}
	for _, field := range []struct {
		name   string
		offset int
		want   uint64
	}{
		{"startSalt at 0x44", 0x44, 11},
		{"callValue at 0x64", 0x64, 1},
		{"gasBuffer at 0x84", 0x84, 50000},
	} {
		if got := binary.BigEndian.Uint64(data[field.offset+24 : field.offset+32]); got != field.want {
			t.Errorf("%s = %d, want %d", field.name, got, field.want)
		}
	}
	if startSalt != 11 {
		t.Errorf("startSalt = %d, want 11", startSalt)
	}

	// The address the contract derives from these fields must be the one the
	// factory actually deployed.
	var salt [32]byte
	binary.BigEndian.PutUint64(salt[24:], startSalt)
	want := crypto.CreateAddress2(pattern.Factory, salt, pattern.InitCodeHash[:])
	got := crypto.CreateAddress2(
		common.BytesToAddress(data[0x04:0x24]),
		common.BytesToHash(data[0x44:0x64]),
		common.BytesToHash(data[0x24:0x44]).Bytes(),
	)
	if got != want {
		t.Errorf("derived target = %s, want %s", got, want)
	}
}

func TestTargetForRoundRobinsPatterns(t *testing.T) {
	options := testOptions(
		Create2Pattern{Name: "minimal", Factory: testFactory, InitCode: testInitCode, StartSalt: 0, Count: 10},
		Create2Pattern{Name: "max-code", Factory: testFactory, InitCodeHash: crypto.Keccak256Hash([]byte("max")).Hex(), StartSalt: 500, Count: 10},
	)
	targets, _, err := buildAttackTargets(options, ".", 0)
	if err != nil {
		t.Fatalf("buildAttackTargets returned error: %v", err)
	}

	wantNames := []string{"minimal", "max-code", "minimal", "max-code", "minimal"}
	wantSalts := []uint64{0, 500, 0, 500, 0}
	for txIdx, wantName := range wantNames {
		pattern, startSalt := targets.targetFor(uint64(txIdx))
		if pattern.Name != wantName {
			t.Errorf("tx %d pattern = %q, want %q", txIdx, pattern.Name, wantName)
		}
		if startSalt != wantSalts[txIdx] {
			t.Errorf("tx %d startSalt = %d, want %d", txIdx, startSalt, wantSalts[txIdx])
		}
	}
}

func TestTargetForSaltStrideSlicesAndWraps(t *testing.T) {
	options := testOptions(Create2Pattern{
		Name: "minimal", Factory: testFactory, InitCode: testInitCode, StartSalt: 100, Count: 10,
	})
	targets, _, err := buildAttackTargets(options, ".", 4)
	if err != nil {
		t.Fatalf("buildAttackTargets returned error: %v", err)
	}

	// count 10 / stride 4 = 2 whole slices; the trailing 2 salts are skipped
	// rather than walked past the end of the range.
	want := []uint64{100, 104, 100, 104, 100}
	for txIdx, wantSalt := range want {
		_, startSalt := targets.targetFor(uint64(txIdx))
		if startSalt != wantSalt {
			t.Errorf("tx %d startSalt = %d, want %d", txIdx, startSalt, wantSalt)
		}
	}
}

func TestTargetForSaltStrideEqualToCount(t *testing.T) {
	options := testOptions(Create2Pattern{
		Name: "minimal", Factory: testFactory, InitCode: testInitCode, StartSalt: 0, Count: 8,
	})
	targets, _, err := buildAttackTargets(options, ".", 8)
	if err != nil {
		t.Fatalf("buildAttackTargets returned error: %v", err)
	}
	for txIdx := uint64(0); txIdx < 3; txIdx++ {
		if _, startSalt := targets.targetFor(txIdx); startSalt != 0 {
			t.Errorf("tx %d startSalt = %d, want 0", txIdx, startSalt)
		}
	}
}

func TestBuildAttackTargetsReportsIgnoredFields(t *testing.T) {
	options := TargetPoolOptions{
		Addresses:       []string{"0x2000000000000000000000000000000000000002"},
		Create2Patterns: []Create2Pattern{{Factory: testFactory, InitCode: testInitCode, Count: 1}},
		Order:           "shuffle",
		Repeat:          true,
		Seed:            "glamsterdam-p0",
	}

	_, ignored, err := buildAttackTargets(options, ".", 0)
	if err != nil {
		t.Fatalf("buildAttackTargets returned error: %v", err)
	}
	joined := strings.Join(ignored, " ")
	for _, want := range []string{"addresses", "order", "repeat", "seed"} {
		if !strings.Contains(joined, want) {
			t.Errorf("ignored fields %q do not mention %q", joined, want)
		}
	}
}

func TestBuildAttackTargetsReadsInitCodeFile(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "initcode.hex"), []byte("  "+testInitCode+"\n"), 0o644); err != nil {
		t.Fatalf("could not write initcode file: %v", err)
	}

	options := testOptions(Create2Pattern{
		Factory: testFactory, InitCodeFile: "initcode.hex", Count: 1,
	})
	targets, _, err := buildAttackTargets(options, dir, 0)
	if err != nil {
		t.Fatalf("buildAttackTargets returned error: %v", err)
	}
	if got, want := targets.patterns[0].InitCodeHash, crypto.Keccak256Hash(common.FromHex(testInitCode)); got != want {
		t.Errorf("initCodeHash = %s, want %s", got, want)
	}
}

func TestLoadTargetPoolOptionsAcceptsWrappedAndDirectFiles(t *testing.T) {
	body := `create2_patterns:
  - name: minimal
    factory: "` + testFactory + `"
    initcode: "` + testInitCode + `"
    start_salt: 3
    count: 9
`
	tests := map[string]string{
		"direct":  body,
		"wrapped": "targets:\n  " + strings.ReplaceAll(body, "\n", "\n  "),
	}

	for name, content := range tests {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "targets.yaml")
			if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
				t.Fatalf("could not write targets file: %v", err)
			}

			options, baseDir, err := loadTargetPoolOptions(path)
			if err != nil {
				t.Fatalf("loadTargetPoolOptions returned error: %v", err)
			}
			if baseDir != filepath.Dir(path) {
				t.Errorf("baseDir = %q, want %q", baseDir, filepath.Dir(path))
			}
			if len(options.Create2Patterns) != 1 {
				t.Fatalf("pattern count = %d, want 1", len(options.Create2Patterns))
			}
			if got := options.Create2Patterns[0]; got.StartSalt != 3 || got.Count != 9 {
				t.Errorf("salt range = %d count %d, want 3 count 9", got.StartSalt, got.Count)
			}
		})
	}
}

func TestTargetPoolYAMLRejectsUnknownFields(t *testing.T) {
	tests := map[string]string{
		"pool":    "unexpected: 1\n",
		"pattern": "create2_patterns:\n  - factory: \"" + testFactory + "\"\n    unexpected: 1\n",
	}

	for name, content := range tests {
		t.Run(name, func(t *testing.T) {
			var options TargetPoolOptions
			err := yaml.Unmarshal([]byte(content), &options)
			if err == nil || !strings.Contains(err.Error(), "unknown field") {
				t.Fatalf("yaml.Unmarshal error = %v, want error containing \"unknown field\"", err)
			}
		})
	}
}

func TestBuildAttackTargetsValidation(t *testing.T) {
	tests := []struct {
		name       string
		options    TargetPoolOptions
		saltStride uint64
		wantErr    string
	}{
		{
			name:    "no patterns",
			options: TargetPoolOptions{Addresses: []string{"0x2000000000000000000000000000000000000002"}},
			wantErr: "at least one CREATE2 pattern",
		},
		{
			name:    "zero count",
			options: testOptions(Create2Pattern{Factory: testFactory, InitCode: testInitCode, Count: 0}),
			wantErr: "zero count",
		},
		{
			name:    "invalid factory",
			options: testOptions(Create2Pattern{Factory: "not-an-address", InitCode: testInitCode, Count: 1}),
			wantErr: "invalid factory address",
		},
		{
			name:    "multiple initcode inputs",
			options: testOptions(Create2Pattern{Factory: testFactory, InitCode: testInitCode, InitCodeHash: crypto.Keccak256Hash([]byte("x")).Hex(), Count: 1}),
			wantErr: "exactly one of initcode",
		},
		{
			name:    "no initcode input",
			options: testOptions(Create2Pattern{Factory: testFactory, Count: 1}),
			wantErr: "exactly one of initcode",
		},
		{
			name:    "bad initcode hash",
			options: testOptions(Create2Pattern{Factory: testFactory, InitCodeHash: "0x1234", Count: 1}),
			wantErr: "32-byte hex value",
		},
		{
			name:    "salt range overflow",
			options: testOptions(Create2Pattern{Factory: testFactory, InitCode: testInitCode, StartSalt: ^uint64(0), Count: 2}),
			wantErr: "overflows uint64",
		},
		{
			name: "duplicate names",
			options: testOptions(
				Create2Pattern{Name: "dup", Factory: testFactory, InitCode: testInitCode, Count: 1},
				Create2Pattern{Name: "dup", Factory: testFactory, InitCode: testInitCode, Count: 1},
			),
			wantErr: "duplicate CREATE2 pattern name",
		},
		{
			name:       "stride larger than range",
			options:    testOptions(Create2Pattern{Factory: testFactory, InitCode: testInitCode, Count: 4}),
			saltStride: 5,
			wantErr:    "exceeds the 4 salts",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, _, err := buildAttackTargets(test.options, ".", test.saltStride)
			if err == nil || !strings.Contains(err.Error(), test.wantErr) {
				t.Fatalf("buildAttackTargets error = %v, want error containing %q", err, test.wantErr)
			}
		})
	}
}

func TestInitAttackTargetsBuildsPool(t *testing.T) {
	s := &Scenario{
		options: ScenarioOptions{
			Targets:   testOptions(Create2Pattern{Name: "minimal", Factory: testFactory, InitCode: testInitCode, Count: 5}),
			GasBuffer: 50000,
		},
		logger: logrus.New().WithField("test", t.Name()),
	}
	if err := s.initAttackTargets(); err != nil {
		t.Fatalf("initAttackTargets returned error: %v", err)
	}
	if s.attackTargets == nil || len(s.attackTargets.patterns) != 1 {
		t.Fatalf("attack targets were not initialized: %+v", s.attackTargets)
	}
}

func TestInitAttackTargetsStaysNilWithoutTargets(t *testing.T) {
	s := &Scenario{
		options: ScenarioOptions{CallFnSig: "foo()", GasBuffer: 50000},
		logger:  logrus.New().WithField("test", t.Name()),
	}
	if err := s.initAttackTargets(); err != nil {
		t.Fatalf("initAttackTargets returned error: %v", err)
	}
	if s.attackTargets != nil {
		t.Fatalf("attack targets = %+v, want nil", s.attackTargets)
	}
}

func TestInitAttackTargetsRejectsConflicts(t *testing.T) {
	pool := testOptions(Create2Pattern{Factory: testFactory, InitCode: testInitCode, Count: 5})
	tests := []struct {
		name    string
		options ScenarioOptions
		wantErr string
	}{
		{
			name:    "inline and file",
			options: ScenarioOptions{Targets: pool, TargetsFile: "targets.yaml", GasBuffer: 50000},
			wantErr: "cannot be used together",
		},
		{
			name:    "manual call data",
			options: ScenarioOptions{Targets: pool, CallData: "0x1234", GasBuffer: 50000},
			wantErr: "cannot be combined with call_data",
		},
		{
			name:    "manual fn sig",
			options: ScenarioOptions{Targets: pool, CallFnSig: "foo()", CallArgs: "[]", GasBuffer: 50000},
			wantErr: "call_fn_sig, call_args",
		},
		{
			name:    "zero gas buffer",
			options: ScenarioOptions{Targets: pool, GasBuffer: 0},
			wantErr: "gas_buffer must be greater than zero",
		},
		{
			name:    "missing file",
			options: ScenarioOptions{TargetsFile: filepath.Join(t.TempDir(), "missing.yaml"), GasBuffer: 50000},
			wantErr: "could not read targets file",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &Scenario{options: test.options, logger: logrus.New().WithField("test", t.Name())}
			err := s.initAttackTargets()
			if err == nil || !strings.Contains(err.Error(), test.wantErr) {
				t.Fatalf("initAttackTargets error = %v, want error containing %q", err, test.wantErr)
			}
		})
	}
}

func TestResolveFactoryPlaceholdersRequiresWalletPool(t *testing.T) {
	s := &Scenario{
		options: ScenarioOptions{GasBuffer: 50000},
		logger:  logrus.New().WithField("test", t.Name()),
	}
	options := testOptions(Create2Pattern{
		Factory: factoryAddressPlaceholder, InitCode: testInitCode, Count: 1,
	})

	err := s.resolveFactoryPlaceholders(&options)
	if err == nil || !strings.Contains(err.Error(), "requires a wallet pool") {
		t.Fatalf("resolveFactoryPlaceholders error = %v, want error containing %q", err, "requires a wallet pool")
	}
}

func TestResolveFactoryPlaceholdersLeavesLiteralAddresses(t *testing.T) {
	s := &Scenario{
		options: ScenarioOptions{GasBuffer: 50000},
		logger:  logrus.New().WithField("test", t.Name()),
	}
	options := testOptions(Create2Pattern{
		Factory: testFactory, InitCode: testInitCode, Count: 1,
	})

	// No placeholder, so the nil wallet pool must not be touched.
	if err := s.resolveFactoryPlaceholders(&options); err != nil {
		t.Fatalf("resolveFactoryPlaceholders returned error: %v", err)
	}
	if got := options.Create2Patterns[0].Factory; got != testFactory {
		t.Errorf("factory = %q, want %q", got, testFactory)
	}
}
