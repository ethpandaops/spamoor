package utils

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestPlaceholdersResolve(t *testing.T) {
	factory := common.HexToAddress("0x3333333333333333333333333333333333333333")
	initCodeHash := crypto.Keccak256Hash([]byte{0x60, 0x00, 0x60, 0x00, 0xf3})
	var salt [32]byte
	salt[31] = 7
	create2 := crypto.CreateAddress2(factory, salt, initCodeHash[:]).Hex()

	placeholders := TxPlaceholders(7)
	placeholders["factory_address"] = func(string) (string, error) { return factory.Hex(), nil }

	tests := []struct {
		name, in, want string
	}{
		{"txid", "tx-{txid}", "tx-7"},
		{"nested create2", "{create2:{factory_address}:" + initCodeHash.Hex() + ":{txid}}", create2},
		{"hex salt without 0x hash prefix", "{create2:" + factory.Hex() + ":" + initCodeHash.Hex()[2:] + ":0x7}", create2},
		{"unknown placeholder kept", "{contract:token} {txid}", "{contract:token} 7"},
		{"json object kept", `[{"a":1},"{txid}"]`, `[{"a":1},"7"]`},
		{"geas block kept", "#define %m() = { push 1 }", "#define %m() = { push 1 }"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := placeholders.Resolve(tt.in)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPlaceholdersRandom(t *testing.T) {
	got, err := TxPlaceholders(0).Resolve("{random:1} {randomaddr} {random}")
	if err != nil {
		t.Fatal(err)
	}
	fields := strings.Fields(got)
	if fields[0] != "0" || !common.IsHexAddress(fields[1]) || strings.Contains(fields[2], "{") {
		t.Fatalf("unexpected resolution %q", got)
	}
}

func TestPlaceholdersErrors(t *testing.T) {
	for _, in := range []string{
		"{random:abc}",
		"{create2:0x33:0x00:1}",
		"{create2:0x3333333333333333333333333333333333333333:0x1234:1}",
		"{create2:0x3333333333333333333333333333333333333333:0x1234}",
		"{create2:0x3333333333333333333333333333333333333333:" + common.Hash{}.Hex() + ":-1}",
	} {
		if _, err := TxPlaceholders(0).Resolve(in); err == nil {
			t.Errorf("%s: expected error", in)
		}
	}
}
