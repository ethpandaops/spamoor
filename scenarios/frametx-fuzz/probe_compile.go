// The probe contract: what frames call to reach the instructions EIP-8141 and its
// extension EIPs introduce, together with the deploy-and-use lifecycle around it.
//
// A generator that only sends frames at plain wallet addresses never executes a single
// one of the new instructions, and three of the four recognized validation prefixes
// need a contract to play the sender or the paymaster.
package frametxfuzz

import (
	_ "embed"
	"fmt"
	"sync"

	geas "github.com/fjl/geas/asm"
)

//go:embed probe.geas
var contractSource string

// initcodeSource returns the runtime code that follows it, which is the standard
// constructor for a contract with no constructor arguments.
const initcodeSource = `
	push @.start
	codesize
	sub
	dup1
	push @.start
	push0
	codecopy
	push0
	return

	.start:
`

var (
	compileOnce    sync.Once
	compiledCode   []byte
	compiledInit   []byte
	compileFailure error
)

// ProbeRuntimeCode returns the probe contract's runtime code.
func ProbeRuntimeCode() ([]byte, error) {
	compileOnce.Do(compileContract)

	return compiledCode, compileFailure
}

// ProbeInitCode returns the deployment code producing ProbeRuntimeCode.
//
// It is deterministic, which matters: the contract is deployed through the CREATE2
// factory, so the same source always yields the same address and a run that finds the
// contract already there can skip deployment.
func ProbeInitCode() ([]byte, error) {
	compileOnce.Do(compileContract)

	return compiledInit, compileFailure
}

// compileContract assembles the embedded source once.
func compileContract() {
	compiler := geas.NewCompiler(nil)

	runtime := compiler.CompileString(contractSource)
	if runtime == nil || compiler.Failed() {
		compileFailure = fmt.Errorf("failed to compile the frame probe contract: %v", compiler.Errors())

		return
	}

	initCompiler := geas.NewCompiler(nil)

	prefix := initCompiler.CompileString(initcodeSource)
	if prefix == nil || initCompiler.Failed() {
		compileFailure = fmt.Errorf("failed to compile the frame probe initcode: %v", initCompiler.Errors())

		return
	}

	compiledCode = runtime
	compiledInit = append(append([]byte{}, prefix...), runtime...)
}
