package frametxfuzz

import (
	"bytes"
	"testing"

	"github.com/ethereum/go-ethereum/core/vm/runtime"
)

// TestGeneratedCodeIsDeterministic checks the property the reproduction line rests on:
// the same seed, transaction index and frame position deploy byte-identical code.
func TestGeneratedCodeIsDeterministic(t *testing.T) {
	first := GenerateContractCode("0xfeed", 7, 2, 256, 400_000)
	second := GenerateContractCode("0xfeed", 7, 2, 256, 400_000)

	if !bytes.Equal(first, second) {
		t.Fatal("the same seed and position generated different code")
	}

	if bytes.Equal(first, GenerateContractCode("0xfeed", 7, 3, 256, 400_000)) {
		t.Error("two frames of the same transaction generated identical code")
	}

	if bytes.Equal(first, GenerateContractCode("0xbeef", 7, 2, 256, 400_000)) {
		t.Error("two seeds generated identical code")
	}
}

// TestDeploymentInitCodeReturnsTheRuntime checks that the constructor deploys arbitrary
// generated code unchanged.
//
// This is the load-bearing part of deploying from inside a frame: if the constructor
// mangled the code, or depended on what the code does, no contract would ever appear at
// the address a later frame is already pointed at.
func TestDeploymentInitCodeReturnsTheRuntime(t *testing.T) {
	for _, size := range []int{64, 256, 512} {
		generated := GenerateContractCode("0x01", 1, 0, size, 400_000)

		initcode, err := DeploymentInitCode(generated)
		if err != nil {
			t.Fatalf("could not build init code: %v", err)
		}

		deployed, _, _, err := runtime.Create(initcode, &runtime.Config{GasLimit: 20_000_000})
		if err != nil {
			t.Fatalf("deploying %d bytes of generated code failed: %v", len(generated), err)
		}

		if !bytes.Equal(deployed, generated) {
			t.Errorf("deployed code differs from the generated runtime (%d vs %d bytes)",
				len(deployed), len(generated))
		}
	}
}

// TestGeneratedCodeCarriesFrameInstructions checks that the frame instructions actually
// reach the generated contracts.
//
// Without this the generator would happily produce ordinary EVM code and the whole point
// of adding them to its table would be lost silently.
func TestGeneratedCodeCarriesFrameInstructions(t *testing.T) {
	wanted := map[byte]string{
		opcodeTxParam:       "TXPARAM",
		opcodeFrameDataLoad: "FRAMEDATALOAD",
		opcodeFrameDataCopy: "FRAMEDATACOPY",
		opcodeFrameParam:    "FRAMEPARAM",
		opcodeSigParam:      "SIGPARAM",
		opcodeSigDataCopy:   "SIGDATACOPY / RECENTROOTREFLOAD",
		opcodeApprove:       "APPROVE",
	}

	seen := map[byte]bool{}

	for index := uint64(0); index < 64; index++ {
		for _, opcode := range instructionsIn(GenerateContractCode("0x02", index, 0, 512, 400_000)) {
			if _, ok := wanted[opcode]; ok {
				seen[opcode] = true
			}
		}
	}

	for opcode, name := range wanted {
		if !seen[opcode] {
			t.Errorf("%s (0x%02x) never appeared in 64 generated contracts", name, opcode)
		}
	}
}

// instructionsIn returns the opcodes of a bytecode stream, skipping push data so that a
// byte inside a PUSH is not mistaken for an instruction.
func instructionsIn(code []byte) []byte {
	opcodes := make([]byte, 0, len(code))

	for i := 0; i < len(code); i++ {
		op := code[i]
		opcodes = append(opcodes, op)

		if op >= 0x60 && op <= 0x7f {
			i += int(op-0x60) + 1
		}
	}

	return opcodes
}
