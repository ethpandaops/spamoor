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
	compiledPrefix []byte
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
	runtime, err := ProbeRuntimeCode()
	if err != nil {
		return nil, err
	}

	return DeploymentInitCode(runtime)
}

// DeploymentInitCode wraps runtime code in the constructor that returns it.
//
// The constructor is generic -- it copies everything after itself and returns it -- so
// generated contracts reuse the same wrapper as the probe. Deployment then succeeds
// whatever the runtime code does, which is what puts fuzzed code in a frame's call
// rather than in a constructor that may never finish.
func DeploymentInitCode(runtime []byte) ([]byte, error) {
	compileOnce.Do(compileContract)

	if compileFailure != nil {
		return nil, compileFailure
	}

	initcode := make([]byte, 0, len(compiledPrefix)+len(runtime))
	initcode = append(initcode, compiledPrefix...)
	initcode = append(initcode, runtime...)

	return initcode, nil
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
	compiledPrefix = prefix
}
