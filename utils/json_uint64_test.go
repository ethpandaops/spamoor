package utils

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

// sentinel is a non-zero default so we can tell "left unset" apart from "parsed 0".
const sentinel = FlexibleJsonUInt64(1 << 62)

func TestFlexibleJsonUInt64_JSON(t *testing.T) {
	cases := []struct {
		in   string
		want FlexibleJsonUInt64
	}{
		{`123`, 123},          // bare int
		{`"123"`, 123},        // quoted int
		{`"null"`, sentinel},  // unset -> keep default
		{`""`, sentinel},      // unset -> keep default
		{`null`, sentinel},    // JSON null literal -> keep default
	}
	for _, c := range cases {
		v := sentinel
		if err := json.Unmarshal([]byte(c.in), &v); err != nil {
			t.Fatalf("JSON %q: unexpected error: %v", c.in, err)
		}
		if v != c.want {
			t.Errorf("JSON %q: got %d, want %d", c.in, v, c.want)
		}
	}
}

func TestFlexibleJsonUInt64_YAML(t *testing.T) {
	cases := []struct {
		in   string
		want FlexibleJsonUInt64
	}{
		{`"123"`, 123},
		{`"null"`, sentinel},
		{`""`, sentinel},
	}
	for _, c := range cases {
		v := sentinel
		if err := yaml.Unmarshal([]byte(c.in), &v); err != nil {
			t.Fatalf("YAML %q: unexpected error: %v", c.in, err)
		}
		if v != c.want {
			t.Errorf("YAML %q: got %d, want %d", c.in, v, c.want)
		}
	}
}
