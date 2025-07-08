package aitx

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type PayloadProcessor struct {
	logger        logrus.FieldLogger
	geasProcessor *GeasProcessor
}

func NewPayloadProcessor(logger logrus.FieldLogger) *PayloadProcessor {
	return &PayloadProcessor{
		logger:        logger.WithField("component", "payload_processor"),
		geasProcessor: NewGeasProcessor(logger),
	}
}

func (pp *PayloadProcessor) ValidatePayload(payload *PayloadTemplate) error {
	if payload.Type != "geas" {
		return fmt.Errorf("only 'geas' type is supported, got: %s", payload.Type)
	}

	if payload.Description == "" {
		return fmt.Errorf("payload description is required")
	}

	if payload.InitCode == "" {
		return fmt.Errorf("geas payload requires init_code")
	}

	if payload.RunCode == "" {
		return fmt.Errorf("geas payload requires run_code")
	}

	// Validate geas compilation
	if err := pp.validateGeasCompilation(payload); err != nil {
		return fmt.Errorf("geas compilation failed: %w", err)
	}

	return nil
}

func (pp *PayloadProcessor) ProcessPayloads(templates []PayloadTemplate) ([]PayloadTemplate, error) {
	var validPayloads []PayloadTemplate

	for i, template := range templates {
		err := pp.ValidatePayload(&template)
		if err != nil {
			pp.logger.Errorf("invalid payload #%d (%s): %v", i+1, template.Description, err)
			continue
		}

		pp.logger.Infof("payload #%d (%s) validated successfully", i+1, template.Description)
		validPayloads = append(validPayloads, template)
	}

	if len(validPayloads) == 0 {
		return nil, fmt.Errorf("no valid payloads found")
	}

	pp.logger.Infof("processed %d payloads, %d valid", len(templates), len(validPayloads))
	return validPayloads, nil
}

func (pp *PayloadProcessor) validateGeasCompilation(payload *PayloadTemplate) error {
	// Create a temporary payload instance for compilation testing
	tempPayload := &PayloadInstance{
		Type:         payload.Type,
		Description:  payload.Description,
		InitCode:     payload.InitCode,
		RunCode:      payload.RunCode,
		PostCode:     payload.PostCode,
		GasRemainder: 10000, // Default value for validation
	}

	// Attempt to compile the geas code
	_, err := pp.geasProcessor.CompileGeasPayload(tempPayload)
	if err != nil {
		pp.logger.Debugf("geas compilation validation failed for payload '%s': %v", payload.Description, err)
		return err
	}

	pp.logger.Debugf("geas compilation validation passed for payload '%s'", payload.Description)
	return nil
}
