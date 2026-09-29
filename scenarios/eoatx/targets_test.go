package eoatx

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func TestBuildTargetPoolCombinesExplicitAndCreate2Targets(t *testing.T) {
	factory := common.HexToAddress("0x1000000000000000000000000000000000000001")
	initCode := []byte{0x60, 0x00, 0x60, 0x00, 0xf3}
	options := TargetPoolOptions{
		Addresses: []string{
			"0x2000000000000000000000000000000000000002",
			"0x3000000000000000000000000000000000000003",
		},
		Create2Patterns: []Create2Pattern{{
			Name:      "minimal-contracts",
			Factory:   factory.Hex(),
			InitCode:  "0x60006000f3",
			StartSalt: 7,
			Count:     2,
		}},
		Order: "sequential",
	}

	pool, sourceCounts, err := buildTargetPool(options, ".")
	if err != nil {
		t.Fatalf("buildTargetPool returned error: %v", err)
	}
	if got, want := pool.len(), uint64(4); got != want {
		t.Fatalf("pool length = %d, want %d", got, want)
	}
	if got, want := sourceCounts["addresses"], uint64(2); got != want {
		t.Errorf("explicit address count = %d, want %d", got, want)
	}
	if got, want := sourceCounts["minimal-contracts"], uint64(2); got != want {
		t.Errorf("CREATE2 address count = %d, want %d", got, want)
	}

	for i, saltValue := range []uint64{7, 8} {
		var salt [32]byte
		binary.BigEndian.PutUint64(salt[24:], saltValue)
		want := crypto.CreateAddress2(factory, salt, crypto.Keccak256(initCode))
		got, ok := pool.targetFor(uint64(i + 2))
		if !ok {
			t.Fatalf("target %d not found", i+2)
		}
		if got.Address != want {
			t.Errorf("target %d = %s, want %s", i+2, got.Address, want)
		}
		if got.Source != "minimal-contracts" {
			t.Errorf("target %d source = %q, want minimal-contracts", i+2, got.Source)
		}
	}

	if _, ok := pool.targetFor(4); ok {
		t.Fatal("non-repeating pool returned a target after exhaustion")
	}
}

func TestBuildTargetPoolShuffleIsDeterministicAndRepeatable(t *testing.T) {
	options := TargetPoolOptions{
		Addresses: []string{
			"0x1000000000000000000000000000000000000001",
			"0x2000000000000000000000000000000000000002",
			"0x3000000000000000000000000000000000000003",
			"0x4000000000000000000000000000000000000004",
		},
		Order:  "shuffle",
		Repeat: true,
		Seed:   "test-seed",
	}

	first, _, err := buildTargetPool(options, ".")
	if err != nil {
		t.Fatalf("first buildTargetPool returned error: %v", err)
	}
	second, _, err := buildTargetPool(options, ".")
	if err != nil {
		t.Fatalf("second buildTargetPool returned error: %v", err)
	}
	if !reflect.DeepEqual(first.targets, second.targets) {
		t.Fatal("same shuffle seed produced different target order")
	}

	for i := uint64(0); i < first.len(); i++ {
		initial, _ := first.targetFor(i)
		repeated, ok := first.targetFor(i + first.len())
		if !ok {
			t.Fatalf("repeating pool did not return target %d", i+first.len())
		}
		if initial != repeated {
			t.Errorf("repeat target %d = %+v, want %+v", i, repeated, initial)
		}
	}
}

func TestBuildTargetPoolReadsInitCodeFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "initcode.hex")
	if err := os.WriteFile(path, []byte("  0x60006000f3\n"), 0o600); err != nil {
		t.Fatalf("write initcode fixture: %v", err)
	}

	options := TargetPoolOptions{Create2Patterns: []Create2Pattern{{
		Factory:      "0x1000000000000000000000000000000000000001",
		InitCodeFile: "initcode.hex",
		Count:        1,
	}}}
	if _, _, err := buildTargetPool(options, dir); err != nil {
		t.Fatalf("buildTargetPool returned error: %v", err)
	}
}

func TestLoadTargetPoolOptionsAcceptsWrappedAndDirectFiles(t *testing.T) {
	for name, body := range map[string]string{
		"wrapped": "targets:\n  addresses:\n    - '0x1000000000000000000000000000000000000001'\n",
		"direct":  "addresses:\n  - '0x1000000000000000000000000000000000000001'\n",
	} {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "targets.yaml")
			if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
				t.Fatalf("write targets fixture: %v", err)
			}
			options, baseDir, err := loadTargetPoolOptions(path)
			if err != nil {
				t.Fatalf("loadTargetPoolOptions returned error: %v", err)
			}
			if got, want := len(options.Addresses), 1; got != want {
				t.Fatalf("address count = %d, want %d", got, want)
			}
			if baseDir != dir {
				t.Fatalf("base directory = %q, want %q", baseDir, dir)
			}
		})
	}
}

func TestBuildTargetPoolValidation(t *testing.T) {
	validFactory := "0x1000000000000000000000000000000000000001"
	tests := []struct {
		name    string
		options TargetPoolOptions
		wantErr string
	}{
		{
			name:    "empty",
			options: TargetPoolOptions{},
			wantErr: "at least one",
		},
		{
			name:    "invalid explicit address",
			options: TargetPoolOptions{Addresses: []string{"not-an-address"}},
			wantErr: "invalid explicit target",
		},
		{
			name: "zero pattern count",
			options: TargetPoolOptions{Create2Patterns: []Create2Pattern{{
				Factory: validFactory, InitCode: "0x00",
			}}},
			wantErr: "zero count",
		},
		{
			name: "invalid factory",
			options: TargetPoolOptions{Create2Patterns: []Create2Pattern{{
				Factory: "invalid", InitCode: "0x00", Count: 1,
			}}},
			wantErr: "invalid factory",
		},
		{
			name: "multiple initcode inputs",
			options: TargetPoolOptions{Create2Patterns: []Create2Pattern{{
				Factory: validFactory, InitCode: "0x00", InitCodeHash: common.Hash{}.Hex(), Count: 1,
			}}},
			wantErr: "exactly one",
		},
		{
			name: "bad initcode hash",
			options: TargetPoolOptions{Create2Patterns: []Create2Pattern{{
				Factory: validFactory, InitCodeHash: "0x1234", Count: 1,
			}}},
			wantErr: "32-byte",
		},
		{
			name: "duplicate",
			options: TargetPoolOptions{Addresses: []string{
				"0x1000000000000000000000000000000000000001",
				"0x1000000000000000000000000000000000000001",
			}},
			wantErr: "duplicate target",
		},
		{
			name:    "invalid order",
			options: TargetPoolOptions{Addresses: []string{validFactory}, Order: "random"},
			wantErr: "invalid target order",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, _, err := buildTargetPool(test.options, ".")
			if err == nil || !strings.Contains(err.Error(), test.wantErr) {
				t.Fatalf("buildTargetPool error = %v, want error containing %q", err, test.wantErr)
			}
		})
	}
}

func TestTargetPoolYAMLRejectsUnknownFields(t *testing.T) {
	for name, body := range map[string]string{
		"pool":    "address:\n  - '0x1000000000000000000000000000000000000001'\n",
		"pattern": "create2_patterns:\n  - factory: '0x1000000000000000000000000000000000000001'\n    init_code: '0x00'\n    count: 1\n",
	} {
		t.Run(name, func(t *testing.T) {
			var options TargetPoolOptions
			if err := yaml.Unmarshal([]byte(body), &options); err == nil || !strings.Contains(err.Error(), "unknown field") {
				t.Fatalf("yaml.Unmarshal error = %v, want unknown field error", err)
			}
		})
	}
}

func TestInitTargetPoolAppliesNonRepeatingCount(t *testing.T) {
	s := &Scenario{
		options: ScenarioOptions{
			Targets: TargetPoolOptions{Addresses: []string{
				"0x1000000000000000000000000000000000000001",
				"0x2000000000000000000000000000000000000002",
			}},
		},
		logger: logrus.New().WithField("test", t.Name()),
	}
	if err := s.initTargetPool(); err != nil {
		t.Fatalf("initTargetPool returned error: %v", err)
	}
	if got, want := s.options.TotalCount, uint64(2); got != want {
		t.Fatalf("total count = %d, want %d", got, want)
	}
	if s.targetPool == nil || s.targetPool.len() != 2 {
		t.Fatalf("target pool was not initialized: %+v", s.targetPool)
	}
}

func TestInitTargetPoolRejectsConflictsAndOversubscription(t *testing.T) {
	address := "0x1000000000000000000000000000000000000001"
	tests := []struct {
		name    string
		options ScenarioOptions
		wantErr string
	}{
		{
			name: "legacy target conflict",
			options: ScenarioOptions{
				To:      address,
				Targets: TargetPoolOptions{Addresses: []string{address}},
			},
			wantErr: "cannot be combined",
		},
		{
			name: "too many transactions",
			options: ScenarioOptions{
				TotalCount: 2,
				Targets:    TargetPoolOptions{Addresses: []string{address}},
			},
			wantErr: "exceeds non-repeating",
		},
		{
			name: "inline and file",
			options: ScenarioOptions{
				Targets:     TargetPoolOptions{Addresses: []string{address}},
				TargetsFile: "targets.yaml",
			},
			wantErr: "cannot be used together",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &Scenario{options: test.options, logger: logrus.New().WithField("test", t.Name())}
			err := s.initTargetPool()
			if err == nil || !strings.Contains(err.Error(), test.wantErr) {
				t.Fatalf("initTargetPool error = %v, want error containing %q", err, test.wantErr)
			}
		})
	}
}
