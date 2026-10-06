package taskrunner

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestNestedCreate2Placeholder(t *testing.T) {
	factory := common.HexToAddress("0x3333333333333333333333333333333333333333")
	initCodeHash := crypto.Keccak256Hash([]byte{0x00})
	registry := NewContractRegistry()
	registry.Set("factory", factory)

	var salt [32]byte
	salt[31] = 3
	want := crypto.CreateAddress2(factory, salt, initCodeHash[:]).Hex()

	for _, stripPrefix := range []bool{false, true} {
		in := "{create2:{contract:factory}:" + initCodeHash.Hex() + ":{txid}}-{stepid}"
		processed, err := ProcessContractPlaceholders(in, registry, stripPrefix)
		if err != nil {
			t.Fatal(err)
		}
		got, err := ProcessBasicPlaceholders(processed, 3, 1, stripPrefix)
		if err != nil {
			t.Fatal(err)
		}
		expected := want + "-1"
		if stripPrefix {
			expected = strings.TrimPrefix(expected, "0x")
		}
		if got != expected {
			t.Fatalf("stripPrefix=%v: got %q, want %q", stripPrefix, got, expected)
		}
	}
}
