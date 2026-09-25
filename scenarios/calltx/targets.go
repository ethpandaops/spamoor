package calltx

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"gopkg.in/yaml.v3"
)

// CallAttackFnSig is the ABI signature the attack calldata is built for. It is
// implemented by both ContractCallAttack.geas and the `callAttack` function of
// ContractReadAttackController.sol.
const CallAttackFnSig = "callAttack(address,bytes32,uint256,uint256,uint256)"

// factoryAddressPlaceholder expands to the well-known CREATE2 factory address
// of the factorydeploytx scenario, in --call-args as well as in the `factory`
// field of a CREATE2 pattern.
const factoryAddressPlaceholder = "{factory_address}"

// TargetPoolOptions describes a pool of deterministic CREATE2 receiver
// addresses. It intentionally mirrors the `targets` block of the eoatx
// scenario so the same file can be pointed at either scenario, but calltx only
// consumes `create2_patterns`: the attack contract derives addresses on-chain
// from (factory, initCodeHash, startSalt), so explicit addresses and the
// ordering controls have nothing to map onto and are ignored with a warning.
type TargetPoolOptions struct {
	Addresses       []string         `yaml:"addresses"`
	Create2Patterns []Create2Pattern `yaml:"create2_patterns"`
	Order           string           `yaml:"order"`
	Repeat          bool             `yaml:"repeat"`
	Seed            string           `yaml:"seed"`
}

func (o *TargetPoolOptions) UnmarshalYAML(node *yaml.Node) error {
	if err := rejectUnknownYAMLFields(node, map[string]struct{}{
		"addresses": {}, "create2_patterns": {}, "order": {}, "repeat": {}, "seed": {},
	}); err != nil {
		return fmt.Errorf("invalid targets configuration: %w", err)
	}
	type plain TargetPoolOptions
	return node.Decode((*plain)(o))
}

// Create2Pattern describes a contiguous range of contracts deployed by one
// CREATE2 factory. Exactly one initcode input must be set.
type Create2Pattern struct {
	Name         string `yaml:"name"`
	Factory      string `yaml:"factory"`
	InitCode     string `yaml:"initcode"`
	InitCodeFile string `yaml:"initcode_file"`
	InitCodeHash string `yaml:"initcode_hash"`
	StartSalt    uint64 `yaml:"start_salt"`
	Count        uint64 `yaml:"count"`
}

func (p *Create2Pattern) UnmarshalYAML(node *yaml.Node) error {
	if err := rejectUnknownYAMLFields(node, map[string]struct{}{
		"name": {}, "factory": {}, "initcode": {}, "initcode_file": {}, "initcode_hash": {}, "start_salt": {}, "count": {},
	}); err != nil {
		return fmt.Errorf("invalid CREATE2 pattern: %w", err)
	}
	type plain Create2Pattern
	return node.Decode((*plain)(p))
}

func rejectUnknownYAMLFields(node *yaml.Node, allowed map[string]struct{}) error {
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("expected a YAML mapping")
	}
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i].Value
		if _, ok := allowed[key]; !ok {
			return fmt.Errorf("unknown field %q", key)
		}
	}
	return nil
}

func (o TargetPoolOptions) configured() bool {
	return len(o.Addresses) > 0 || len(o.Create2Patterns) > 0 || o.Order != "" || o.Repeat || o.Seed != ""
}

// create2Target is one resolved CREATE2 pattern: everything the attack
// contract needs to derive addresses on-chain.
type create2Target struct {
	Name         string
	Factory      common.Address
	InitCodeHash common.Hash
	StartSalt    uint64
	Count        uint64
}

// attackTargets is the round-robin pool of resolved patterns.
type attackTargets struct {
	patterns   []create2Target
	saltStride uint64
}

func loadTargetPoolOptions(path string) (TargetPoolOptions, string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return TargetPoolOptions{}, "", fmt.Errorf("could not read targets file %q: %w", path, err)
	}

	// Accept both a standalone target-pool document and a document wrapped in
	// a top-level `targets` key, so the same block can be copied between a
	// scenario config and --targets-file.
	var wrapped struct {
		Targets *TargetPoolOptions `yaml:"targets"`
	}
	if err := yaml.Unmarshal(data, &wrapped); err != nil {
		return TargetPoolOptions{}, "", fmt.Errorf("could not parse targets file %q: %w", path, err)
	}
	if wrapped.Targets != nil {
		return *wrapped.Targets, filepath.Dir(path), nil
	}

	var options TargetPoolOptions
	if err := yaml.Unmarshal(data, &options); err != nil {
		return TargetPoolOptions{}, "", fmt.Errorf("could not parse targets file %q: %w", path, err)
	}

	return options, filepath.Dir(path), nil
}

// buildAttackTargets resolves every CREATE2 pattern into the factory address,
// initcode hash and salt range that the attack calldata carries.
func buildAttackTargets(options TargetPoolOptions, baseDir string, saltStride uint64) (*attackTargets, []string, error) {
	if len(options.Create2Patterns) == 0 {
		return nil, nil, fmt.Errorf("targets must contain at least one CREATE2 pattern")
	}

	// Fields that only make sense when the caller enumerates addresses itself.
	var ignored []string
	if len(options.Addresses) > 0 {
		ignored = append(ignored, fmt.Sprintf("addresses (%d)", len(options.Addresses)))
	}
	if strings.TrimSpace(options.Order) != "" {
		ignored = append(ignored, fmt.Sprintf("order (%q)", options.Order))
	}
	if options.Repeat {
		ignored = append(ignored, "repeat")
	}
	if strings.TrimSpace(options.Seed) != "" {
		ignored = append(ignored, fmt.Sprintf("seed (%q)", options.Seed))
	}

	patterns := make([]create2Target, 0, len(options.Create2Patterns))
	seen := make(map[string]struct{}, len(options.Create2Patterns))

	for i, pattern := range options.Create2Patterns {
		name := strings.TrimSpace(pattern.Name)
		if name == "" {
			name = fmt.Sprintf("create2_%d", i+1)
		}
		if _, exists := seen[name]; exists {
			return nil, nil, fmt.Errorf("duplicate CREATE2 pattern name %q", name)
		}
		seen[name] = struct{}{}

		if pattern.Count == 0 {
			return nil, nil, fmt.Errorf("create2 pattern %q has zero count", name)
		}
		if pattern.Count-1 > math.MaxUint64-pattern.StartSalt {
			return nil, nil, fmt.Errorf("create2 pattern %q salt range overflows uint64", name)
		}
		if !common.IsHexAddress(pattern.Factory) {
			return nil, nil, fmt.Errorf("create2 pattern %q has invalid factory address %q", name, pattern.Factory)
		}

		initCodeHash, err := pattern.initCodeHash(baseDir)
		if err != nil {
			return nil, nil, fmt.Errorf("create2 pattern %q: %w", name, err)
		}

		if saltStride > pattern.Count {
			return nil, nil, fmt.Errorf("salt_stride %d exceeds the %d salts of create2 pattern %q", saltStride, pattern.Count, name)
		}

		patterns = append(patterns, create2Target{
			Name:         name,
			Factory:      common.HexToAddress(pattern.Factory),
			InitCodeHash: initCodeHash,
			StartSalt:    pattern.StartSalt,
			Count:        pattern.Count,
		})
	}

	return &attackTargets{patterns: patterns, saltStride: saltStride}, ignored, nil
}

// targetFor round-robins over the patterns and returns the pattern for txIdx
// together with the salt that tx starts its walk at.
//
// With saltStride == 0 every tx restarts at the pattern's start_salt, so all
// txs walk the same prefix of the range. With a non-zero stride, consecutive
// txs of the same pattern cover disjoint slices, wrapping back to start_salt
// once the range is exhausted.
func (t *attackTargets) targetFor(txIdx uint64) (create2Target, uint64) {
	pattern := t.patterns[txIdx%uint64(len(t.patterns))]
	if t.saltStride == 0 {
		return pattern, pattern.StartSalt
	}

	slice := txIdx / uint64(len(t.patterns))
	// slices per range, rounded down - a trailing partial slice would walk
	// past the deployed range, so it is skipped rather than truncated.
	slices := pattern.Count / t.saltStride
	if slices == 0 {
		return pattern, pattern.StartSalt
	}

	return pattern, pattern.StartSalt + (slice%slices)*t.saltStride
}

// buildCallData encodes the attack calldata for txIdx:
//
//	callAttack(factory, initCodeHash, startSalt, callValue, gasBuffer)
func (t *attackTargets) buildCallData(txIdx uint64, callValue uint64, gasBuffer uint64) ([]byte, create2Target, uint64) {
	pattern, startSalt := t.targetFor(txIdx)

	data := make([]byte, 0, 4+5*32)
	data = append(data, crypto.Keccak256([]byte(CallAttackFnSig))[:4]...)
	data = append(data, common.LeftPadBytes(pattern.Factory.Bytes(), 32)...)
	data = append(data, pattern.InitCodeHash.Bytes()...)
	data = append(data, uint64Word(startSalt)...)
	data = append(data, uint64Word(callValue)...)
	data = append(data, uint64Word(gasBuffer)...)

	return data, pattern, startSalt
}

func uint64Word(value uint64) []byte {
	word := make([]byte, 32)
	binary.BigEndian.PutUint64(word[24:], value)
	return word
}

func (p Create2Pattern) initCodeHash(baseDir string) (common.Hash, error) {
	inputs := 0
	if strings.TrimSpace(p.InitCode) != "" {
		inputs++
	}
	if strings.TrimSpace(p.InitCodeFile) != "" {
		inputs++
	}
	if strings.TrimSpace(p.InitCodeHash) != "" {
		inputs++
	}
	if inputs != 1 {
		return common.Hash{}, fmt.Errorf("exactly one of initcode, initcode_file, or initcode_hash must be set")
	}

	if p.InitCodeHash != "" {
		decoded, err := decodeHex(p.InitCodeHash)
		if err != nil || len(decoded) != common.HashLength {
			return common.Hash{}, fmt.Errorf("initcode_hash must be a 32-byte hex value")
		}
		return common.BytesToHash(decoded), nil
	}

	initCodeText := p.InitCode
	if p.InitCodeFile != "" {
		path := p.InitCodeFile
		if !filepath.IsAbs(path) {
			path = filepath.Join(baseDir, path)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return common.Hash{}, fmt.Errorf("could not read initcode file %q: %w", path, err)
		}
		initCodeText = strings.TrimSpace(string(data))
	}

	initCode, err := decodeHex(initCodeText)
	if err != nil || len(initCode) == 0 {
		return common.Hash{}, fmt.Errorf("initcode must be non-empty hex data")
	}
	return crypto.Keccak256Hash(initCode), nil
}

func decodeHex(value string) ([]byte, error) {
	value = strings.TrimSpace(value)
	if !strings.HasPrefix(value, "0x") && !strings.HasPrefix(value, "0X") {
		value = "0x" + value
	}
	return hexutil.Decode(value)
}
