package aitx

import (
	"fmt"
	"strings"

	geas "github.com/fjl/geas/asm"
	"github.com/sirupsen/logrus"
)

type GeasProcessor struct {
	logger logrus.FieldLogger
}

func NewGeasProcessor(logger logrus.FieldLogger) *GeasProcessor {
	return &GeasProcessor{
		logger: logger.WithField("component", "geas_processor"),
	}
}

func (gp *GeasProcessor) CompileGeasPayload(payload *PayloadInstance) ([]byte, error) {
	compiler := geas.NewCompiler(nil)
	return gp.compileInitRunGeas(payload, compiler)
}

func (gp *GeasProcessor) compileInitRunGeas(payload *PayloadInstance, compiler *geas.Compiler) ([]byte, error) {

	// Build init code that deploys the contract
	initcodeGeas := `
	;; Init code
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

	// Build the contract template with init, run, and post code
	contractGeasTpl := `
	%s
	gas                   ;; [gas, custom]
	push 0                ;; [loop_counter, gas, custom]
	jump @loop

	exit:
		;; Execute post code once at the end
		%s
        stop              ;; [custom]

	loop:
		push %d           ;; [gas_remainder, loop_counter, gas, custom]
		gas               ;; [gas, gas_remainder, loop_counter, gas, custom]
		lt                ;; [gas < gas_remainder, loop_counter, gas, custom]
		jumpi @exit       ;; [loop_counter, gas, custom]

		;; increase loop_counter
		push 1            ;; [1, loop_counter, gas, custom]
		add               ;; [loop_counter+1, gas, custom]

		;; run the performance test code
		%s

		jump @loop
	`

	gp.logger.Debugf("compiling init_run geas - init: %s, run: %s, post: %s",
		strings.ReplaceAll(payload.InitCode, "\n", "\\n"),
		strings.ReplaceAll(payload.RunCode, "\n", "\\n"),
		strings.ReplaceAll(payload.PostCode, "\n", "\\n"))

	// Compile init code
	initcode := compiler.CompileString(initcodeGeas)
	if initcode == nil {
		return nil, fmt.Errorf("failed to compile geas init code: %v", compiler.Errors())
	}

	// Compile the contract code with init, run, and post parts
	contractCode := compiler.CompileString(fmt.Sprintf(contractGeasTpl, payload.InitCode, payload.PostCode, payload.GasRemainder, payload.RunCode))
	if contractCode == nil {
		return nil, fmt.Errorf("failed to compile geas contract code: %v", compiler.Errors())
	}

	// Combine init code and contract code
	combinedCode := append(initcode, contractCode...)

	gp.logger.Debugf("compiled init_run geas code to %d bytes (init: %d, contract: %d)",
		len(combinedCode), len(initcode), len(contractCode))

	return combinedCode, nil
}
