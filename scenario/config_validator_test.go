package scenario

import (
	"reflect"
	"strings"
	"testing"
)

// validatorTestOptions mimics a typical scenario options struct with yaml tags.
type validatorTestOptions struct {
	Throughput uint64 `yaml:"throughput" usage:"transactions per slot"`
	MaxPending uint64 `yaml:"max_pending,omitempty" usage:"max pending transactions"`
	DashField  string `yaml:"dash-field"`
	Hidden     string `yaml:"-"`
	NoTag      string
	Fallback   string `yaml:",omitempty"`
}

func testDescriptor() *Descriptor {
	return &Descriptor{
		Name:           "testscenario",
		Description:    "scenario for validator tests",
		DefaultOptions: &validatorTestOptions{Throughput: 10},
	}
}

func TestGetScenarioValidFields(t *testing.T) {
	fields := GetScenarioValidFields(testDescriptor())

	throughput, ok := fields["throughput"]
	if !ok {
		t.Fatal("expected 'throughput' field to be extracted")
	}
	if throughput.Type != reflect.TypeOf(uint64(0)) {
		t.Fatalf("expected uint64 type for 'throughput', got %v", throughput.Type)
	}
	if throughput.DefaultValue != uint64(10) {
		t.Fatalf("expected default value 10 for 'throughput', got %v", throughput.DefaultValue)
	}
	if throughput.Description != "transactions per slot" {
		t.Fatalf("unexpected description: %q", throughput.Description)
	}

	// yaml tag options like ",omitempty" are stripped from the field name
	if _, ok := fields["max_pending"]; !ok {
		t.Fatal("expected 'max_pending' field to be extracted")
	}

	// a yaml tag without a name falls back to the lowercased Go field name
	if _, ok := fields["fallback"]; !ok {
		t.Fatal("expected 'fallback' field to be extracted")
	}

	// fields without a yaml tag or explicitly excluded are not extracted
	if _, ok := fields["hidden"]; ok {
		t.Fatal("expected 'hidden' (yaml:\"-\") field to be excluded")
	}
	if _, ok := fields["notag"]; ok {
		t.Fatal("expected untagged field to be excluded")
	}

	// wallet pool configuration fields are always valid
	for _, walletField := range []string{"wallet_count", "seed", "refill_amount"} {
		if _, ok := fields[walletField]; !ok {
			t.Fatalf("expected wallet pool field %q to be extracted", walletField)
		}
	}
}

func TestGetScenarioValidFieldsNonStructOptions(t *testing.T) {
	fields := GetScenarioValidFields(&Descriptor{
		Name:           "badoptions",
		DefaultOptions: "not a struct",
	})

	// only the wallet pool fields remain
	if _, ok := fields["wallet_count"]; !ok {
		t.Fatal("expected wallet pool fields to be extracted")
	}
	if _, ok := fields["throughput"]; ok {
		t.Fatal("did not expect scenario fields from a non-struct options value")
	}
}

func TestValidateConfig(t *testing.T) {
	fields := GetScenarioValidFields(testDescriptor())
	validator := NewConfigValidator("testscenario", fields, discardLogger())

	tests := []struct {
		name        string
		config      string
		valid       bool
		errContains string
	}{
		{name: "empty config", config: "", valid: true},
		{name: "known fields", config: "throughput: 10\nmax_pending: 5", valid: true},
		{name: "invalid yaml", config: "throughput: [unclosed", valid: false, errContains: "Invalid YAML"},
		{name: "unknown field", config: "totally_unknown_xyz: 1", valid: false, errContains: "invalid configuration fields"},
		{name: "dash instead of underscore", config: "max-pending: 5", valid: false, errContains: "'max-pending' -> 'max_pending'"},
		{name: "underscore instead of dash", config: "dash_field: x", valid: false, errContains: "'dash_field' -> 'dash-field'"},
		{name: "wrong case", config: "THROUGHPUT: 5", valid: false, errContains: "'THROUGHPUT' -> 'throughput'"},
		{name: "partial field name", config: "pending: 5", valid: false, errContains: "-> 'max_pending'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateConfig(tt.config)
			if result.Valid != tt.valid {
				t.Fatalf("expected valid=%v, got %v (errors: %v)", tt.valid, result.Valid, result.Errors)
			}
			if tt.valid {
				if len(result.Errors) != 0 {
					t.Fatalf("expected no errors, got %v", result.Errors)
				}
				return
			}
			if len(result.Errors) == 0 {
				t.Fatal("expected validation errors")
			}
			if tt.errContains != "" && !strings.Contains(result.Errors[0], tt.errContains) {
				t.Fatalf("expected error to contain %q, got %q", tt.errContains, result.Errors[0])
			}
		})
	}
}

func TestValidateScenarioConfig(t *testing.T) {
	if err := ValidateScenarioConfig(nil, "", discardLogger()); err == nil {
		t.Fatal("expected error for nil descriptor")
	}

	if err := ValidateScenarioConfig(testDescriptor(), "throughput: 20", discardLogger()); err != nil {
		t.Fatalf("expected valid config to pass, got %v", err)
	}

	if err := ValidateScenarioConfig(testDescriptor(), "totally_unknown_xyz: 1", discardLogger()); err == nil {
		t.Fatal("expected error for unknown config field")
	}
}

func TestParseAndValidateConfig(t *testing.T) {
	target := &validatorTestOptions{}
	err := ParseAndValidateConfig(testDescriptor(), "throughput: 42\nmax_pending: 7", target, discardLogger())
	if err != nil {
		t.Fatalf("expected config to parse, got %v", err)
	}
	if target.Throughput != 42 || target.MaxPending != 7 {
		t.Fatalf("expected parsed values 42/7, got %d/%d", target.Throughput, target.MaxPending)
	}

	// an empty config leaves the target untouched
	target = &validatorTestOptions{Throughput: 5}
	if err := ParseAndValidateConfig(testDescriptor(), "", target, discardLogger()); err != nil {
		t.Fatalf("expected empty config to be valid, got %v", err)
	}
	if target.Throughput != 5 {
		t.Fatalf("expected target to be unchanged, got %d", target.Throughput)
	}

	// unknown fields are rejected before parsing
	if err := ParseAndValidateConfig(testDescriptor(), "totally_unknown_xyz: 1", target, discardLogger()); err == nil {
		t.Fatal("expected error for unknown config field")
	}

	// type mismatches surface as unmarshal errors
	err = ParseAndValidateConfig(testDescriptor(), "throughput: notanumber", target, discardLogger())
	if err == nil || !strings.Contains(err.Error(), "failed to unmarshal config") {
		t.Fatalf("expected unmarshal error, got %v", err)
	}
}
