package bloatedeoa

import (
	"embed"
	"errors"
	"fmt"
	"math/big"

	geas "github.com/fjl/geas/asm"
)

//go:embed contract/*.geas
var contractFS embed.FS

// attackContractFile returns the geas source of the attack contract variant
// matching the given options.
func attackContractFile(options *ScenarioOptions) string {
	switch {
	case options.Mode == ModeSload:
		return "contract/sload.geas"
	case options.WriteNewValue:
		return "contract/sstore_increment.geas"
	case options.ExistingSlots:
		return "contract/sstore_same.geas"
	default:
		return "contract/sstore_zero.geas"
	}
}

// compileAttackContract compiles the attack contract variant matching the
// given options and returns its runtime code.
func compileAttackContract(options *ScenarioOptions) ([]byte, error) {
	compiler := geas.NewCompiler(contractFS)
	compiler.SetGlobal("LoopGasThreshold", new(big.Int).SetUint64(options.LoopGasThreshold))

	return compileGeasFile(compiler, attackContractFile(options))
}

// compileInitcode compiles the constructor that deploys the runtime code
// appended after it.
func compileInitcode() ([]byte, error) {
	return compileGeasFile(geas.NewCompiler(contractFS), "contract/initcode.geas")
}

func compileGeasFile(compiler *geas.Compiler, file string) ([]byte, error) {
	code := compiler.CompileFile(file)
	if code == nil {
		return nil, fmt.Errorf("failed to compile %v: %w", file, errors.Join(compiler.Errors()...))
	}

	return code, nil
}
