package eoatx

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	mathrand "math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"gopkg.in/yaml.v3"
)

// TargetPoolOptions describes a combined pool of explicit and deterministic
// CREATE2 receiver addresses.
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

type targetEntry struct {
	Address common.Address
	Source  string
}

type targetPool struct {
	targets []targetEntry
	repeat  bool
}

func (o TargetPoolOptions) configured() bool {
	return len(o.Addresses) > 0 || len(o.Create2Patterns) > 0 || o.Order != "" || o.Repeat || o.Seed != ""
}

func (p *targetPool) len() uint64 {
	return uint64(len(p.targets))
}

func (p *targetPool) targetFor(index uint64) (targetEntry, bool) {
	if p == nil || len(p.targets) == 0 {
		return targetEntry{}, false
	}

	if index >= uint64(len(p.targets)) {
		if !p.repeat {
			return targetEntry{}, false
		}
		index %= uint64(len(p.targets))
	}

	return p.targets[index], true
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

func buildTargetPool(options TargetPoolOptions, baseDir string) (*targetPool, map[string]uint64, error) {
	order := strings.ToLower(strings.TrimSpace(options.Order))
	if order == "" {
		order = "sequential"
	}
	if order != "sequential" && order != "shuffle" {
		return nil, nil, fmt.Errorf("invalid target order %q (must be sequential or shuffle)", options.Order)
	}

	totalCount := uint64(len(options.Addresses))
	for i, pattern := range options.Create2Patterns {
		if pattern.Count == 0 {
			return nil, nil, fmt.Errorf("create2 pattern %d (%q) has zero count", i+1, pattern.Name)
		}
		if pattern.Count-1 > math.MaxUint64-pattern.StartSalt {
			return nil, nil, fmt.Errorf("create2 pattern %d (%q) salt range overflows uint64", i+1, pattern.Name)
		}
		if pattern.Count > math.MaxUint64-totalCount {
			return nil, nil, fmt.Errorf("target count overflows uint64")
		}
		totalCount += pattern.Count
	}
	if totalCount == 0 {
		return nil, nil, fmt.Errorf("targets must contain at least one address or CREATE2 pattern")
	}
	if totalCount > uint64(int(^uint(0)>>1)) {
		return nil, nil, fmt.Errorf("target count %d exceeds platform capacity", totalCount)
	}

	targets := make([]targetEntry, 0, int(totalCount))
	sourceCounts := make(map[string]uint64)
	seen := make(map[common.Address]string, int(totalCount))
	appendTarget := func(entry targetEntry) error {
		if previousSource, exists := seen[entry.Address]; exists {
			return fmt.Errorf("duplicate target %s from %q (already provided by %q)", entry.Address.Hex(), entry.Source, previousSource)
		}
		seen[entry.Address] = entry.Source
		targets = append(targets, entry)
		sourceCounts[entry.Source]++
		return nil
	}

	for i, addressText := range options.Addresses {
		if !common.IsHexAddress(addressText) {
			return nil, nil, fmt.Errorf("invalid explicit target address at index %d: %q", i, addressText)
		}
		if err := appendTarget(targetEntry{Address: common.HexToAddress(addressText), Source: "addresses"}); err != nil {
			return nil, nil, err
		}
	}

	for i, pattern := range options.Create2Patterns {
		name := strings.TrimSpace(pattern.Name)
		if name == "" {
			name = fmt.Sprintf("create2_%d", i+1)
		}
		if !common.IsHexAddress(pattern.Factory) {
			return nil, nil, fmt.Errorf("create2 pattern %q has invalid factory address %q", name, pattern.Factory)
		}

		initCodeHash, err := pattern.initCodeHash(baseDir)
		if err != nil {
			return nil, nil, fmt.Errorf("create2 pattern %q: %w", name, err)
		}

		factory := common.HexToAddress(pattern.Factory)
		for offset := uint64(0); offset < pattern.Count; offset++ {
			var salt [32]byte
			binary.BigEndian.PutUint64(salt[24:], pattern.StartSalt+offset)
			address := crypto.CreateAddress2(factory, salt, initCodeHash[:])
			if err := appendTarget(targetEntry{Address: address, Source: name}); err != nil {
				return nil, nil, err
			}
		}
	}

	if order == "shuffle" {
		seed := options.Seed
		if seed == "" {
			seed = "eoatx-targets"
		}
		hash := sha256.Sum256([]byte(seed))
		rng := mathrand.New(mathrand.NewSource(int64(binary.BigEndian.Uint64(hash[:8]))))
		rng.Shuffle(len(targets), func(i, j int) {
			targets[i], targets[j] = targets[j], targets[i]
		})
	}

	return &targetPool{targets: targets, repeat: options.Repeat}, sourceCounts, nil
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
