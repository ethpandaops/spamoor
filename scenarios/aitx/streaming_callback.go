package aitx

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

type PayloadStreamingCallback struct {
	processor     *PayloadProcessor
	logger        logrus.FieldLogger
	payloadBuffer *strings.Builder

	// Track parsing state
	mutex          sync.Mutex
	parsedPayloads []PayloadTemplate
	totalContent   strings.Builder
}

func (c *PayloadStreamingCallback) OnContent(content string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Simply accumulate all content
	c.totalContent.WriteString(content)
	c.payloadBuffer.WriteString(content)

	// Check for complete JSON blocks in the accumulated content
	c.processCompleteJSONBlocks()

	return nil
}

// processCompleteJSONBlocks searches for and processes complete JSON blocks
func (c *PayloadStreamingCallback) processCompleteJSONBlocks() {
	content := c.payloadBuffer.String()

	//c.logger.Debugf("streaming rsp: %s", content)

	// Look for complete JSON blocks (```json ... ```)
	for {
		startIdx := strings.Index(content, "```json")
		if startIdx == -1 {
			break
		}

		// Find the end of this JSON block
		endMarker := "```"
		endIdx := strings.Index(content[startIdx+7:], endMarker) // Skip past "```json"
		if endIdx == -1 {
			// Incomplete block, wait for more content
			break
		}

		// Extract the JSON content
		jsonStart := startIdx + 7 // Skip "```json"
		jsonEnd := startIdx + 7 + endIdx
		jsonStr := strings.TrimSpace(content[jsonStart:jsonEnd])

		if jsonStr != "" {
			c.logger.Debugf("Found complete JSON block: %s", jsonStr)

			// Try to parse and validate immediately
			var payload PayloadTemplate
			if err := json.Unmarshal([]byte(jsonStr), &payload); err == nil {
				c.logger.Infof("Streaming: Parsed payload '%s', validating...", payload.Description)

				// Validate payload in real-time
				if validatedPayloads, err := c.processor.ProcessPayloads([]PayloadTemplate{payload}); err == nil {
					c.parsedPayloads = append(c.parsedPayloads, validatedPayloads...)
					c.logger.Infof("Streaming: Validated and ready payload '%s'", payload.Description)
				} else {
					c.logger.Warnf("Streaming: Payload validation failed for '%s': %v", payload.Description, err)
				}
			} else {
				c.logger.Warnf("Streaming: Failed to parse JSON block: %v", err)
			}
		}

		// Remove the processed block from buffer and continue searching
		content = content[jsonEnd+3:] // Skip past the closing ```
	}

	// Update the buffer with remaining unprocessed content
	c.payloadBuffer.Reset()
	c.payloadBuffer.WriteString(content)
}

func (c *PayloadStreamingCallback) OnComplete(fullContent string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.logger.Infof("Streaming completed. Total payloads parsed in real-time: %d", len(c.parsedPayloads))

	// Always try to extract additional payloads from the full content in case streaming missed any
	// This handles cases where JSON blocks span multiple content chunks or have formatting issues
	additionalPayloads, err := c.extractJSONObjectsFromText(fullContent)
	if err == nil && len(additionalPayloads) > len(c.parsedPayloads) {
		c.logger.Infof("Found %d additional payloads in complete content (streaming found %d)",
			len(additionalPayloads)-len(c.parsedPayloads), len(c.parsedPayloads))

		// Validate the additional payloads
		if validatedPayloads, err := c.processor.ProcessPayloads(additionalPayloads); err == nil {
			c.parsedPayloads = validatedPayloads
			c.logger.Infof("Final: Validated %d total payloads", len(c.parsedPayloads))
		} else {
			c.logger.Warnf("Failed to validate additional payloads: %v", err)
		}
	} else if len(c.parsedPayloads) == 0 {
		c.logger.Debugf("No payloads parsed during streaming, attempting full content parsing")

		// Try to extract payloads from the complete content as fallback
		payloads, err := c.extractJSONObjectsFromText(fullContent)
		if err != nil {
			return fmt.Errorf("failed to parse any payloads from complete content: %w", err)
		}

		// Validate all payloads at once as fallback
		if validatedPayloads, err := c.processor.ProcessPayloads(payloads); err == nil {
			c.parsedPayloads = validatedPayloads
			c.logger.Infof("Fallback: Validated %d payloads from complete content", len(c.parsedPayloads))
		} else {
			return fmt.Errorf("failed to validate fallback payloads: %w", err)
		}
	}

	return nil
}

func (c *PayloadStreamingCallback) GetParsedPayloads() []PayloadTemplate {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.parsedPayloads
}

// extractJSONObjectsFromText is a fallback method copied from ai_service.go
func (c *PayloadStreamingCallback) extractJSONObjectsFromText(content string) ([]PayloadTemplate, error) {
	var payloads []PayloadTemplate

	// Look for JSON code blocks marked with ```json
	lines := strings.Split(content, "\n")
	var jsonBlock strings.Builder
	inJSONBlock := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "```json") {
			inJSONBlock = true
			jsonBlock.Reset()
			continue
		}

		if strings.HasPrefix(line, "```") && inJSONBlock {
			// End of JSON block, try to parse it
			jsonStr := jsonBlock.String()
			c.logger.Debugf("Attempting to parse JSON block: %s", jsonStr)

			var payload PayloadTemplate
			if err := json.Unmarshal([]byte(jsonStr), &payload); err == nil {
				payloads = append(payloads, payload)
				c.logger.Debugf("Successfully parsed payload: %s", payload.Description)
			} else {
				c.logger.Errorf("Failed to parse JSON block: %v", err)
			}

			inJSONBlock = false
			continue
		}

		if inJSONBlock {
			jsonBlock.WriteString(line)
			jsonBlock.WriteString("\n")
		}
	}

	// If we found payloads, return them
	if len(payloads) > 0 {
		return payloads, nil
	}

	// Fallback: look for individual JSON objects using regex-like approach
	return c.extractJSONObjectsWithRegex(content)
}

func (c *PayloadStreamingCallback) extractJSONObjectsWithRegex(content string) ([]PayloadTemplate, error) {
	var payloads []PayloadTemplate

	// Look for patterns like { ... } that might be JSON objects
	braceLevel := 0
	var currentObj strings.Builder
	inObject := false

	for i, r := range content {
		if r == '{' {
			if braceLevel == 0 {
				inObject = true
				currentObj.Reset()
			}
			braceLevel++
			currentObj.WriteRune(r)
		} else if r == '}' {
			braceLevel--
			currentObj.WriteRune(r)

			if braceLevel == 0 && inObject {
				// Try to parse this object
				objStr := strings.TrimSpace(currentObj.String())
				c.logger.Infof("Attempting to parse JSON object: %s", objStr)

				var payload PayloadTemplate
				if err := json.Unmarshal([]byte(objStr), &payload); err == nil {
					payloads = append(payloads, payload)
					c.logger.Infof("Successfully parsed payload: %s", payload.Description)
				} else {
					c.logger.Errorf("Failed to parse JSON object at position %d: %v", i, err)
				}

				inObject = false
			}
		} else if inObject {
			currentObj.WriteRune(r)
		}
	}

	if len(payloads) == 0 {
		return nil, fmt.Errorf("no valid JSON objects found in response")
	}

	return payloads, nil
}
