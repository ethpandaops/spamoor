package aitx

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type AIService struct {
	client           *http.Client
	apiKey           string
	model            string
	baseURL          string
	basePrompt       string
	tokenCount       uint64
	callCount        uint64
	logger           logrus.FieldLogger
	logConversations bool
}

type GenerationRequest struct {
	BasePrompt          string
	TestDirection       string
	GenerationMode      string
	PayloadCount        uint64
	PreviousSummary     string
	TransactionFeedback *TransactionFeedback
}

type TransactionFeedback struct {
	TotalTransactions    uint64              `json:"total_transactions"`
	SuccessfulTxs        uint64              `json:"successful_txs"`
	FailedTxs            uint64              `json:"failed_txs"`
	AverageGasUsed       uint64              `json:"average_gas_used"`
	MedianGasUsed        uint64              `json:"median_gas_used"`
	AverageBlockExecTime string              `json:"average_block_exec_time"`
	RecentResults        []TransactionResult `json:"recent_results"`
	Summary              string              `json:"summary"`
}

type TransactionResult struct {
	PayloadType        string   `json:"payload_type"`
	PayloadDescription string   `json:"payload_description"`
	Status             string   `json:"status"`
	GasUsed            uint64   `json:"gas_used"`
	BlockExecTime      string   `json:"block_exec_time"`
	ErrorMessage       string   `json:"error_message,omitempty"`
	LogData            []string `json:"log_data,omitempty"`
}

type GenerationResponse struct {
	Payloads   []PayloadTemplate
	Summary    string
	TokensUsed uint64
}

type ConversationContinuation struct {
	History  []Message
	Feedback string
}

type OpenRouterRequest struct {
	Model     string           `json:"model"`
	Messages  []Message        `json:"messages"`
	MaxTokens int              `json:"max_tokens"`
	Stream    bool             `json:"stream,omitempty"`
	Reasoning *ReasoningConfig `json:"reasoning,omitempty"`
}

type ReasoningConfig struct {
	MaxTokens int `json:"max_tokens"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenRouterResponse struct {
	ID       string `json:"id"`
	Provider string `json:"provider,omitempty"`
	Object   string `json:"object"`
	Created  int64  `json:"created"`
	Model    string `json:"model"`
	Choices  []struct {
		Index              int      `json:"index"`
		Message            *Message `json:"message,omitempty"`
		Delta              *Delta   `json:"delta,omitempty"`
		FinishReason       *string  `json:"finish_reason,omitempty"`
		NativeFinishReason *string  `json:"native_finish_reason,omitempty"`
		Logprobs           *string  `json:"logprobs,omitempty"`
	} `json:"choices"`
	Usage *struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage,omitempty"`
}

type Delta struct {
	Role             string   `json:"role,omitempty"`
	Content          string   `json:"content,omitempty"`
	Reasoning        string   `json:"reasoning,omitempty"`
	ReasoningDetails []string `json:"reasoning_details,omitempty"`
}

type StreamingCallback interface {
	OnContent(content string) error
	OnComplete(fullContent string) error
}

func NewAIService(apiKey, model string, logConversations bool, logger logrus.FieldLogger) *AIService {
	if apiKey == "" {
		apiKey = os.Getenv("OPENROUTER_API_KEY")
	}

	return &AIService{
		client: &http.Client{
			Timeout: 30 * time.Minute, // Increased timeout for longer AI requests
		},
		apiKey:           apiKey,
		model:            model,
		baseURL:          "https://openrouter.ai/api/v1/chat/completions",
		logger:           logger.WithField("component", "ai_service"),
		logConversations: logConversations,
	}
}

func (ai *AIService) SetBasePrompt(prompt string) {
	ai.basePrompt = prompt
}

func (ai *AIService) GeneratePayloads(ctx context.Context, req GenerationRequest, processor *PayloadProcessor) (*GenerationResponse, error) {
	maxRetries := 3
	var lastError error
	var conversationHistory []Message

	// Build initial prompt
	req.BasePrompt = ai.basePrompt
	initialPrompt := ai.buildPrompt(req)

	conversationHistory = append(conversationHistory, Message{
		Role:    "user",
		Content: initialPrompt,
	})

	for attempt := 0; attempt < maxRetries; attempt++ {
		ai.callCount++
		ai.logger.Debugf("making AI request #%d (attempt %d/%d) for %d payloads",
			ai.callCount, attempt+1, maxRetries, req.PayloadCount)

		openRouterReq := OpenRouterRequest{
			Model:     ai.model,
			Messages:  conversationHistory,
			MaxTokens: 10000,
			Reasoning: &ReasoningConfig{
				MaxTokens: 5000,
			},
		}

		response, err := ai.callOpenRouter(ctx, openRouterReq)
		if err != nil {
			lastError = fmt.Errorf("AI API call failed: %w", err)
			continue
		}

		if response.Usage != nil {
			ai.tokenCount += uint64(response.Usage.TotalTokens)
			ai.logger.Infof("AI call #%d completed: %d tokens used, %d total tokens",
				ai.callCount, response.Usage.TotalTokens, ai.tokenCount)
		} else {
			ai.logger.Infof("AI call #%d completed (token usage not available)", ai.callCount)
		}

		// Try to parse the response
		result, parseErr := ai.parseResponse(response)
		if parseErr == nil {
			// Validate payloads (including geas compilation)
			validPayloads, validationErr := processor.ProcessPayloads(result.Payloads)
			if validationErr == nil {
				// Log AI response for debugging if enabled
				if ai.logConversations {
					ai.logConversation(conversationHistory, attempt+1)
				}

				// Success! Update result with validated payloads and return
				result.Payloads = validPayloads
				ai.logger.Infof("AI conversation #%d completed successfully after %d attempt(s)", ai.callCount, attempt+1)
				return result, nil
			}
			// Validation failed, treat as parsing error for retry
			parseErr = validationErr
		}

		// Parsing failed, add the AI response and error feedback to conversation
		lastError = parseErr

		// Add AI response to conversation history
		if len(response.Choices) > 0 {
			conversationHistory = append(conversationHistory, Message{
				Role:    "assistant",
				Content: response.Choices[0].Message.Content,
			})
		}

		// Add error feedback for retry
		errorFeedback := ai.buildErrorFeedback(parseErr, attempt+1, maxRetries)
		conversationHistory = append(conversationHistory, Message{
			Role:    "user",
			Content: errorFeedback,
		})

		ai.logger.Warnf("AI request #%d failed (attempt %d/%d): %v, retrying...",
			ai.callCount, attempt+1, maxRetries, parseErr)
	}

	return nil, fmt.Errorf("failed to generate valid payloads after %d attempts, last error: %w", maxRetries, lastError)
}

func (ai *AIService) GeneratePayloadsWithConversation(ctx context.Context, req GenerationRequest, processor *PayloadProcessor, continuation *ConversationContinuation) (*GenerationResponse, []Message, error) {
	maxRetries := 3
	var lastError error
	var conversationHistory []Message

	if continuation != nil {
		// Continue existing conversation
		conversationHistory = continuation.History
		conversationHistory = append(conversationHistory, Message{
			Role:    "user",
			Content: continuation.Feedback,
		})
	} else {
		// Start new conversation
		req.BasePrompt = ai.basePrompt
		initialPrompt := ai.buildPrompt(req)
		conversationHistory = append(conversationHistory, Message{
			Role:    "user",
			Content: initialPrompt,
		})
	}

	for attempt := 0; attempt < maxRetries; attempt++ {
		ai.callCount++
		ai.logger.Debugf("making AI request #%d (attempt %d/%d) for conversation",
			ai.callCount, attempt+1, maxRetries)

		openRouterReq := OpenRouterRequest{
			Model:     ai.model,
			Messages:  conversationHistory,
			MaxTokens: 50000,
			Stream:    true,
			Reasoning: &ReasoningConfig{
				MaxTokens: 20000,
			},
		}

		// Create streaming callback for real-time payload processing
		callback := &PayloadStreamingCallback{
			processor:     processor,
			logger:        ai.logger,
			payloadBuffer: &strings.Builder{},
		}

		response, fullContent, err := ai.callOpenRouterStreaming(ctx, openRouterReq, callback)
		if err != nil {
			ai.logger.Warnf("AI streaming call failed: %v", err)
			lastError = fmt.Errorf("AI streaming call failed: %w", err)
			continue
		}

		if response.Usage != nil {
			ai.tokenCount += uint64(response.Usage.TotalTokens)
			ai.logger.Infof("AI call #%d completed: %d tokens used, %d total tokens",
				ai.callCount, response.Usage.TotalTokens, ai.tokenCount)
		}

		// Try to parse the final response
		mockResponse := &OpenRouterResponse{
			Choices: []struct {
				Index              int      `json:"index"`
				Message            *Message `json:"message,omitempty"`
				Delta              *Delta   `json:"delta,omitempty"`
				FinishReason       *string  `json:"finish_reason,omitempty"`
				NativeFinishReason *string  `json:"native_finish_reason,omitempty"`
				Logprobs           *string  `json:"logprobs,omitempty"`
			}{{Message: &Message{Content: fullContent}}},
		}

		result, parseErr := ai.parseResponse(mockResponse)
		if parseErr == nil {
			// Validate payloads (including geas compilation)
			validPayloads, validationErr := processor.ProcessPayloads(result.Payloads)
			if validationErr == nil {
				// Success! Update result with validated payloads and add AI response to history
				result.Payloads = validPayloads
				conversationHistory = append(conversationHistory, Message{
					Role:    "assistant",
					Content: fullContent,
				})

				// Log AI response for debugging if enabled
				if ai.logConversations {
					ai.logConversation(conversationHistory, attempt+1)
				}

				ai.logger.Infof("AI conversation #%d completed successfully after %d attempt(s)", ai.callCount, attempt+1)
				return result, conversationHistory, nil
			}
			// Validation failed, treat as parsing error for retry
			parseErr = validationErr
		}

		// Parsing failed, add the AI response and error feedback to conversation
		lastError = parseErr

		// Add AI response to conversation history
		conversationHistory = append(conversationHistory, Message{
			Role:    "assistant",
			Content: fullContent,
		})

		// Add error feedback for retry
		errorFeedback := ai.buildErrorFeedback(parseErr, attempt+1, maxRetries)
		conversationHistory = append(conversationHistory, Message{
			Role:    "user",
			Content: errorFeedback,
		})

		ai.logger.Warnf("AI request #%d failed (attempt %d/%d): %v, retrying...",
			ai.callCount, attempt+1, maxRetries, parseErr)
	}

	return nil, conversationHistory, fmt.Errorf("failed to generate valid payloads after %d attempts, last error: %w", maxRetries, lastError)
}

func (ai *AIService) buildPrompt(req GenerationRequest) string {
	var promptBuilder strings.Builder

	promptBuilder.WriteString(req.BasePrompt)
	promptBuilder.WriteString("\n\n")

	if req.TestDirection != "" {
		promptBuilder.WriteString(fmt.Sprintf("TEST DIRECTION: %s\n\n", req.TestDirection))
	}

	promptBuilder.WriteString(fmt.Sprintf("Generate 10 transaction payload(s) with placeholder variations, so we can test at least %v different patterns.\n", req.PayloadCount))

	if req.TransactionFeedback != nil {
		promptBuilder.WriteString("FEEDBACK FROM PREVIOUS TRANSACTIONS:\n")
		promptBuilder.WriteString(fmt.Sprintf("Total executed: %d (Success: %d, Failed: %d)\n",
			req.TransactionFeedback.TotalTransactions,
			req.TransactionFeedback.SuccessfulTxs,
			req.TransactionFeedback.FailedTxs))
		promptBuilder.WriteString(fmt.Sprintf("Gas usage - Average: %d, Median: %d\n",
			req.TransactionFeedback.AverageGasUsed,
			req.TransactionFeedback.MedianGasUsed))
		promptBuilder.WriteString(fmt.Sprintf("Average block execution time: %s\n",
			req.TransactionFeedback.AverageBlockExecTime))

		if len(req.TransactionFeedback.RecentResults) > 0 {
			promptBuilder.WriteString("\nRecent transaction results:\n")
			for _, result := range req.TransactionFeedback.RecentResults {
				promptBuilder.WriteString(fmt.Sprintf("- %s: %s (gas: %d, block_time: %s)\n",
					result.PayloadDescription, result.Status,
					result.GasUsed, result.BlockExecTime))
				if result.ErrorMessage != "" {
					promptBuilder.WriteString(fmt.Sprintf("  Error: %s\n", result.ErrorMessage))
				}
				if len(result.LogData) > 0 {
					promptBuilder.WriteString(fmt.Sprintf("  Logs: %v\n", result.LogData))
				}
			}
		}

		if req.TransactionFeedback.Summary != "" {
			promptBuilder.WriteString(fmt.Sprintf("\nPrevious summary: %s\n", req.TransactionFeedback.Summary))
		}

		promptBuilder.WriteString("\nPlease generate NEW, DIFFERENT payloads that:\n")
		promptBuilder.WriteString("1. Avoid patterns that consistently failed\n")
		promptBuilder.WriteString("2. Explore different gas usage patterns\n")
		promptBuilder.WriteString("3. Consider block execution time impact\n")
		promptBuilder.WriteString("4. Build on successful patterns but with variations\n")
		promptBuilder.WriteString("5. Consider log data from successful transactions\n\n")
	}

	if req.PreviousSummary != "" {
		promptBuilder.WriteString(fmt.Sprintf("Previous generation summary: %s\n\n", req.PreviousSummary))
	}

	return promptBuilder.String()
}

func (ai *AIService) buildBasePrompt(generationMode string) string {
	var promptBuilder strings.Builder

	promptBuilder.WriteString("You are an Ethereum transaction generator for the Spamoor testing framework.\n")
	promptBuilder.WriteString("Your role is to create geas init/run contracts for comprehensive EVM testing.\n\n")

	promptBuilder.WriteString("GEAS INIT/RUN CONTRACT GENERATION:\n\n")

	promptBuilder.WriteString("CONCEPT:\n")
	promptBuilder.WriteString("The init/run pattern deploys a contract with two phases:\n")
	promptBuilder.WriteString("1. INIT PHASE: Executes ONCE during contract deployment (constructor)\n")
	promptBuilder.WriteString("2. RUN PHASE: Executes in a LOOP when the contract is called, consuming all available gas\n\n")

	promptBuilder.WriteString("EXECUTION MODEL:\n")
	promptBuilder.WriteString("1. Contract is deployed with init_code executing once\n")
	promptBuilder.WriteString("2. Contract is then CALLED with optional calldata\n")
	promptBuilder.WriteString("3. Run code executes repeatedly until gas is almost exhausted\n")
	promptBuilder.WriteString("4. Post code executes ONCE at the end when gas is low (for final LOGs/cleanup)\n")
	promptBuilder.WriteString("5. Each run iteration MUST maintain clean stack (no pollution)\n\n")

	promptBuilder.WriteString("CRITICAL REQUIREMENTS:\n")
	promptBuilder.WriteString("1. RUN CODE should reuse previous iteration results for subsequent operations to avoid intermediate result caching in the EVM\n")
	promptBuilder.WriteString("2. RUN CODE may modify stack to keep track of previous results - push empty value from init code, modify via SWAPn in loop\n")
	promptBuilder.WriteString("3. Stack must be same size at the end of each run iteration (but may contain different values)\n")
	promptBuilder.WriteString("4. POST CODE executes once at end when gas is low - ideal for LOG events to report final results\n")
	promptBuilder.WriteString("5. Init, run, and post code can access calldata using CALLDATALOAD, CALLDATASIZE, CALLDATACOPY\n")
	promptBuilder.WriteString("6. Avoid LOG events in run code (expensive) - use post code for final result logging\n\n")

	promptBuilder.WriteString("CALLDATA ACCESS:\n")
	promptBuilder.WriteString("- CALLDATASIZE: Get size of input data\n")
	promptBuilder.WriteString("- PUSH1 0x00 CALLDATALOAD: Load first 32 bytes of calldata\n")
	promptBuilder.WriteString("- PUSH1 0x20 CALLDATALOAD: Load second 32 bytes of calldata\n")
	promptBuilder.WriteString("- CALLDATACOPY: Copy calldata to memory\n\n")

	promptBuilder.WriteString("GEAS CODE FORMAT:\n")
	promptBuilder.WriteString("- ONE opcode per line, separated by \\n\n")
	promptBuilder.WriteString("- Uppercase opcodes only\n")
	promptBuilder.WriteString("- Hex values with 0x prefix\n")
	promptBuilder.WriteString("- Example: PUSH1 0x20\\nPUSH1 0x00\\nMSTORE\n\n")

	promptBuilder.WriteString("EXAMPLE PATTERNS:\n")
	promptBuilder.WriteString("1. Parameter processing: Load calldata, perform operations, store results\n")
	promptBuilder.WriteString("2. Computation loops: Mathematical operations with clean stack management\n")
	promptBuilder.WriteString("3. Storage patterns: Read/write with counters or mappings\n")
	promptBuilder.WriteString("4. Event emission: Log computation results or state changes\n")
	promptBuilder.WriteString("5. Memory operations: Expand memory, hash data, manipulate arrays\n")
	promptBuilder.WriteString("6. Precompile/contract calls: Use CALL opcode to interact with other contracts\n\n")

	promptBuilder.WriteString("PRECOMPILE/CONTRACT CALL PATTERN:\n")
	promptBuilder.WriteString("To call precompiles (addresses 1-9) or other contracts, use this pattern:\n")
	promptBuilder.WriteString("```\n")
	promptBuilder.WriteString("PUSH1 0x20 ; retSize\n")
	promptBuilder.WriteString("PUSH1 0x00 ; retOffset\n")
	promptBuilder.WriteString("PUSH1 0x20 ; argsSize\n")
	promptBuilder.WriteString("PUSH1 0x00 ; argsOffset\n")
	promptBuilder.WriteString("PUSH1 0x00 ; value\n")
	promptBuilder.WriteString("PUSH1 0x05 ; address (example: precompile 5 = modexp)\n")
	promptBuilder.WriteString("PUSH2 0xC350 ; gas (50000)\n")
	promptBuilder.WriteString("GAS\n")
	promptBuilder.WriteString("SUB\n")
	promptBuilder.WriteString("CALL\n")
	promptBuilder.WriteString("POP        ; remove success flag\n")
	promptBuilder.WriteString("```\n")
	promptBuilder.WriteString("Common precompiles: 1=ecrecover, 2=sha256, 3=ripemd160, 4=identity, 5=modexp, 6=ecadd, 7=ecmul, 8=ecpairing, 9=blake2f\n\n")

	promptBuilder.WriteString("AVAILABLE PLACEHOLDERS:\n")
	promptBuilder.WriteString("- ${RANDOM_UINT256}: Random 256-bit unsigned integer\n")
	promptBuilder.WriteString("- ${RANDOM_BYTES32}: Random 32-byte value\n")
	promptBuilder.WriteString("- ${CURRENT_BLOCK}: Current block number\n\n")

	promptBuilder.WriteString("RESPONSE FORMAT:\n")
	promptBuilder.WriteString("CRITICAL: Your response is parsed programmatically. Return ONLY JSON objects in ```json blocks with NO explanations.\n")
	promptBuilder.WriteString("Generate at least 20 separate JSON objects (do not stop before), each wrapped in ```json and ``` tags:\n\n")

	promptBuilder.WriteString(`{
  "id": "unique_payload_id_1",
  "type": "geas",
  "description": "Brief description of what this contract does",
  "init_code": "PUSH1 0x00\nSSTORE",
  "run_code": "PUSH1 0x00\nSLOAD\nPUSH1 0x01\nADD\nDUP1\nPUSH1 0x00\nSSTORE\nPOP",
  "post_code": "PUSH1 0x00\nSLOAD\nPUSH1 0x00\nMSTORE\nPUSH1 0x20\nPUSH1 0x00\nLOG0",
  "gas_remainder": "10000",
  "calldata": "0x1234567800000000000000000000000000000000000000000000000000000005"
}` + "\n\n")

	promptBuilder.WriteString("POST_CODE FIELD:\n")
	promptBuilder.WriteString("- Optional code that executes ONCE at the end when gas is low\n")
	promptBuilder.WriteString("- Ideal for LOG events to report final computation results\n")
	promptBuilder.WriteString("- Can access stack values accumulated during run iterations\n")
	promptBuilder.WriteString("- Example: LOG0 to emit final counter value or computation result\n\n")

	promptBuilder.WriteString("CALLDATA FIELD:\n")
	promptBuilder.WriteString("- Optional hex-encoded calldata for the contract call\n")
	promptBuilder.WriteString("- Can be used to pass parameters to the run code\n")
	promptBuilder.WriteString("- Access in run code via calldataload, calldatasize, etc.\n")
	promptBuilder.WriteString("- Example: \"0x\" + 32-byte parameter as hex\n\n")

	promptBuilder.WriteString("IMPORTANT:\n")
	promptBuilder.WriteString("- Generate ONLY geas init_run contracts (type=\\\"geas\\\")\n")
	promptBuilder.WriteString("- Each payload MUST have a unique 'id' field (e.g., 'payload_1', 'payload_2', etc.)\n")
	promptBuilder.WriteString("- Focus on diverse EVM testing patterns\n")
	promptBuilder.WriteString("- Reuse previous iteration results to avoid EVM caching\n")
	promptBuilder.WriteString("- Use SWAPn to manage persistent values on stack\n")
	promptBuilder.WriteString("- Use calldata for dynamic behavior\n")
	promptBuilder.WriteString("- NO explanatory text - ONLY JSON objects\n\n")

	return promptBuilder.String()
}

func (ai *AIService) callOpenRouter(ctx context.Context, req OpenRouterRequest) (*OpenRouterResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", ai.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+ai.apiKey)
	httpReq.Header.Set("HTTP-Referer", "https://github.com/ethpandaops/spamoor")
	httpReq.Header.Set("X-Title", "Spamoor AI Transaction Generator")

	resp, err := ai.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("OpenRouter API error %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var openRouterResp OpenRouterResponse
	if err := json.Unmarshal(body, &openRouterResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &openRouterResp, nil
}

func (ai *AIService) callOpenRouterStreaming(ctx context.Context, req OpenRouterRequest, callback StreamingCallback) (*OpenRouterResponse, string, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", ai.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+ai.apiKey)
	httpReq.Header.Set("HTTP-Referer", "https://github.com/ethpandaops/spamoor")
	httpReq.Header.Set("X-Title", "Spamoor AI Transaction Generator")
	httpReq.Header.Set("Accept", "text/event-stream")
	httpReq.Header.Set("Cache-Control", "no-cache")

	resp, err := ai.client.Do(httpReq)
	if err != nil {
		return nil, "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, "", fmt.Errorf("OpenRouter API error %d: %s", resp.StatusCode, string(body))
	}

	return ai.parseStreamingResponse(ctx, resp.Body, callback)
}

func (ai *AIService) parseStreamingResponse(ctx context.Context, body io.Reader, callback StreamingCallback) (*OpenRouterResponse, string, error) {
	scanner := bufio.NewScanner(body)
	var fullContent strings.Builder
	var reasoningBuffer strings.Builder
	var reasoningDetailsBuffer []string
	var lastResponse *OpenRouterResponse

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return nil, fullContent.String(), ctx.Err()
		default:
		}

		line := scanner.Text()

		//ai.logger.Debugf("streaming rsp: %s", line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, ":") {
			continue
		}

		// Parse SSE data line
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")

			// Check for stream end
			if data == "[DONE]" {
				break
			}

			var streamResp OpenRouterResponse
			if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
				ai.logger.Warnf("failed to parse streaming response chunk: %v", err)
				continue
			}

			lastResponse = &streamResp

			// Extract content from delta
			if len(streamResp.Choices) > 0 && streamResp.Choices[0].Delta != nil {
				delta := streamResp.Choices[0].Delta

				// Buffer reasoning and print complete lines only
				if delta.Reasoning != "" {
					reasoningBuffer.WriteString(delta.Reasoning)

					// Check for complete lines in the buffer
					bufferContent := reasoningBuffer.String()
					lines := strings.Split(bufferContent, "\n")

					// Print all complete lines (all but the last if it doesn't end with newline)
					for i := 0; i < len(lines)-1; i++ {
						if lines[i] != "" {
							ai.logger.Debugf("AI reasoning: %s", lines[i])
						}
					}

					// Keep the incomplete line in the buffer
					if len(lines) > 0 && !strings.HasSuffix(bufferContent, "\n") {
						reasoningBuffer.Reset()
						reasoningBuffer.WriteString(lines[len(lines)-1])
					} else {
						reasoningBuffer.Reset()
					}
				}

				// Accumulate reasoning details
				if len(delta.ReasoningDetails) > 0 {
					reasoningDetailsBuffer = append(reasoningDetailsBuffer, delta.ReasoningDetails...)
				}

				// Process content
				if delta.Content != "" {
					fullContent.WriteString(delta.Content)

					// Call streaming callback with new content
					if callback != nil {
						if err := callback.OnContent(delta.Content); err != nil {
							ai.logger.Warnf("streaming callback error: %v", err)
						}
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fullContent.String(), fmt.Errorf("error reading streaming response: %w", err)
	}

	// Print any remaining reasoning content
	if reasoningBuffer.Len() > 0 {
		ai.logger.Debugf("AI reasoning: %s", reasoningBuffer.String())
	}

	// Print accumulated reasoning details
	if len(reasoningDetailsBuffer) > 0 {
		// Join all details and split by lines for cleaner output
		allDetails := strings.Join(reasoningDetailsBuffer, "\n")
		lines := strings.Split(allDetails, "\n")

		for _, line := range lines {
			if line != "" {
				ai.logger.Debugf("AI reasoning detail: %s", line)
			}
		}
	}

	// Call completion callback
	if callback != nil {
		if err := callback.OnComplete(fullContent.String()); err != nil {
			ai.logger.Warnf("streaming completion callback error: %v", err)
		}
	}

	// Return the last response (which should contain usage info) or create a mock response
	if lastResponse == nil {
		lastResponse = &OpenRouterResponse{
			Choices: []struct {
				Index              int      `json:"index"`
				Message            *Message `json:"message,omitempty"`
				Delta              *Delta   `json:"delta,omitempty"`
				FinishReason       *string  `json:"finish_reason,omitempty"`
				NativeFinishReason *string  `json:"native_finish_reason,omitempty"`
				Logprobs           *string  `json:"logprobs,omitempty"`
			}{{Message: &Message{Content: fullContent.String()}}},
		}
	}

	return lastResponse, fullContent.String(), nil
}

func (ai *AIService) parseResponse(response *OpenRouterResponse) (*GenerationResponse, error) {
	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no choices in AI response")
	}

	var content string
	if response.Choices[0].Message != nil {
		content = response.Choices[0].Message.Content
	} else if response.Choices[0].Delta != nil {
		content = response.Choices[0].Delta.Content
	} else {
		return nil, fmt.Errorf("no message or delta content in AI response")
	}

	var payloads []PayloadTemplate

	// Try parsing as direct JSON array first
	err := json.Unmarshal([]byte(content), &payloads)
	if err != nil {
		ai.logger.Debugf("failed to parse as direct JSON array, extracting from conversational response: %v", err)

		// Extract individual JSON objects from conversational text
		payloads, err = ai.extractJSONObjectsFromText(content)
		if err != nil {
			// Fallback to old extraction method
			ai.logger.Debugf("failed to extract JSON objects, trying array extraction: %v", err)
			payloads, err = ai.extractJSONFromText(content)
			if err != nil {
				return nil, fmt.Errorf("failed to parse AI response as JSON: %w", err)
			}
		}
	}

	if len(payloads) == 0 {
		return nil, fmt.Errorf("no payloads found in AI response")
	}

	ai.logger.Infof("Successfully parsed %d payloads from AI response", len(payloads))
	summary := fmt.Sprintf("Generated %d payloads using %s", len(payloads), ai.model)

	var tokensUsed uint64
	if response.Usage != nil {
		tokensUsed = uint64(response.Usage.TotalTokens)
	}

	return &GenerationResponse{
		Payloads:   payloads,
		Summary:    summary,
		TokensUsed: tokensUsed,
	}, nil
}

func (ai *AIService) extractJSONObjectsFromText(content string) ([]PayloadTemplate, error) {
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
			ai.logger.Debugf("Attempting to parse JSON block: %s", jsonStr)

			var payload PayloadTemplate
			if err := json.Unmarshal([]byte(jsonStr), &payload); err == nil {
				payloads = append(payloads, payload)
				ai.logger.Debugf("Successfully parsed payload: %s", payload.Description)
			} else {
				ai.logger.Errorf("Failed to parse JSON block: %v", err)
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
	return ai.extractJSONObjectsWithRegex(content)
}

func (ai *AIService) extractJSONObjectsWithRegex(content string) ([]PayloadTemplate, error) {
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
				ai.logger.Infof("Attempting to parse JSON object: %s", objStr)

				var payload PayloadTemplate
				if err := json.Unmarshal([]byte(objStr), &payload); err == nil {
					payloads = append(payloads, payload)
					ai.logger.Infof("Successfully parsed payload: %s", payload.Description)
				} else {
					ai.logger.Errorf("Failed to parse JSON object at position %d: %v", i, err)
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

func (ai *AIService) extractJSONFromText(content string) ([]PayloadTemplate, error) {
	start := strings.Index(content, "[")
	end := strings.LastIndex(content, "]")

	if start == -1 || end == -1 || start >= end {
		return nil, fmt.Errorf("no JSON array found in response")
	}

	jsonStr := content[start : end+1]

	var payloads []PayloadTemplate
	err := json.Unmarshal([]byte(jsonStr), &payloads)
	if err != nil {
		return nil, fmt.Errorf("failed to parse extracted JSON: %w", err)
	}

	return payloads, nil
}

func (ai *AIService) GetTokenCount() uint64 {
	return ai.tokenCount
}

func (ai *AIService) buildErrorFeedback(parseErr error, attempt int, maxRetries int) string {
	var feedbackBuilder strings.Builder

	errorStr := parseErr.Error()

	// Check if this is a geas compilation error
	if strings.Contains(errorStr, "geas compilation failed") {
		feedbackBuilder.WriteString("GEAS COMPILATION ERROR DETECTED:\n")
		feedbackBuilder.WriteString(fmt.Sprintf("Error: %v\n\n", parseErr))

		feedbackBuilder.WriteString("Your geas assembly code failed to compile. Please fix the following issues:\n\n")

		feedbackBuilder.WriteString("GEAS CODE REQUIREMENTS:\n")
		feedbackBuilder.WriteString("1. Use VALID EVM opcodes only (e.g., PUSH1, ADD, MUL, SSTORE, SLOAD, etc.)\n")
		feedbackBuilder.WriteString("2. Format: ONE opcode per line, separated by \\n\n")
		feedbackBuilder.WriteString("3. Use correct syntax: 'PUSH1 0x20' with uppercase opcodes\n")
		feedbackBuilder.WriteString("4. Hexadecimal values must start with 0x\n")
		feedbackBuilder.WriteString("5. All opcodes are allowed including selfdestruct, delegatecall, create2\n")
		feedbackBuilder.WriteString("6. Ensure stack balance (don't leave extra items on stack)\n")
		feedbackBuilder.WriteString("7. CRITICAL: Run code MUST have clean stack after each iteration\n\n")

		feedbackBuilder.WriteString("COMMON FIXES:\n")
		feedbackBuilder.WriteString("- Check opcode spelling and case sensitivity\n")
		feedbackBuilder.WriteString("- Verify hex values format (0x prefix)\n")
		feedbackBuilder.WriteString("- Ensure proper stack management with 'pop'\n")
		feedbackBuilder.WriteString("- Use 'pop' to clean up ALL unused stack items\n")
		feedbackBuilder.WriteString("- Remember: sha3 not keccak256 for EVM opcode\n\n")

		feedbackBuilder.WriteString("EXAMPLE VALID GEAS CODE:\n")
		feedbackBuilder.WriteString("\"PUSH1 0x20\\nPUSH1 0x00\\nMSTORE\\nPUSH1 0x20\\nPUSH1 0x00\\nSHA3\\nPOP\"\n\n")
	} else {
		feedbackBuilder.WriteString("PARSING/VALIDATION ERROR DETECTED:\n")
		feedbackBuilder.WriteString(fmt.Sprintf("Error: %v\n\n", parseErr))

		feedbackBuilder.WriteString("Your previous response could not be parsed or validated correctly. ")
		feedbackBuilder.WriteString("Please ensure your response follows the exact JSON format specified.\n\n")

		feedbackBuilder.WriteString("REQUIREMENTS:\n")
		feedbackBuilder.WriteString("1. Wrap JSON payload in ```json and ``` code blocks\n")
		feedbackBuilder.WriteString("2. Return ONLY ONE payload object (not an array)\n")
		feedbackBuilder.WriteString("3. Include all required fields: type, description, init_code, run_code\n")
		feedbackBuilder.WriteString("4. Set type=\"geas\" (init_run method is implied)\n")
		feedbackBuilder.WriteString("5. Use proper JSON syntax with quotes around strings\n")
		feedbackBuilder.WriteString("6. GEAS CODE FORMAT: Use newlines (\\n) to separate opcodes - ONE opcode per line\n")
		feedbackBuilder.WriteString("7. Include optional 'calldata' and 'post_code' fields\n")
		feedbackBuilder.WriteString("8. Do NOT include geas_method or placeholders fields\n\n")
	}

	if attempt < maxRetries {
		feedbackBuilder.WriteString(fmt.Sprintf("This is attempt %d of %d. Please try again with the corrected code.\n", attempt, maxRetries))
	} else {
		feedbackBuilder.WriteString("This is the final attempt. Please ensure your response is properly formatted and valid.\n")
	}

	return feedbackBuilder.String()
}

func (ai *AIService) logConversation(conversationHistory []Message, attempt int) {
	ai.logger.Infof("=== AI Conversation #%d (Attempt %d) ===", ai.callCount, attempt)

	for i, message := range conversationHistory {
		role := strings.ToUpper(message.Role)
		content := message.Content

		ai.logger.Infof("--- Message %d: %s ---\n%s\n", i+1, role, content)
	}

	ai.logger.Infof("=== End Conversation #%d ===", ai.callCount)
}

func (ai *AIService) GetCallCount() uint64 {
	return ai.callCount
}
