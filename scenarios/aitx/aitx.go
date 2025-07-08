package aitx

import (
	"context"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/utils"
)

type ScenarioOptions struct {
	TotalCount  uint64  `yaml:"total_count"`
	Throughput  uint64  `yaml:"throughput"`
	MaxPending  uint64  `yaml:"max_pending"`
	MaxWallets  uint64  `yaml:"max_wallets"`
	Rebroadcast uint64  `yaml:"rebroadcast"`
	BaseFee     float64 `yaml:"base_fee"`
	TipFee      float64 `yaml:"tip_fee"`
	GasLimit    uint64  `yaml:"gas_limit"`
	Timeout     string  `yaml:"timeout"`
	ClientGroup string  `yaml:"client_group"`
	LogTxs      bool    `yaml:"log_txs"`

	// AI-specific options
	OpenRouterAPIKey   string `yaml:"openrouter_api_key"`
	Model              string `yaml:"model"`
	TestDirection      string `yaml:"test_direction"`
	PayloadsPerRequest uint64 `yaml:"payloads_per_request"`
	MaxAICalls         uint64 `yaml:"max_ai_calls"`
	MaxTokens          uint64 `yaml:"max_tokens"`

	// Generation options
	GenerationMode string `yaml:"generation_mode"` // "geas", "calldata", "transfer", "mixed"

	// Feedback options
	FeedbackBatchSize  uint64 `yaml:"feedback_batch_size"`
	EnableFeedbackLoop bool   `yaml:"enable_feedback_loop"`

	// Debug options
	LogAIConversations bool `yaml:"log_ai_conversations"`
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	aiService              *AIService
	processor              *PayloadProcessor
	placeholderSubstituter *PlaceholderSubstituter
	feedbackCollector      *FeedbackCollector
	geasProcessor          *GeasProcessor

	payloadCache          []PayloadTemplate
	cacheIndex            int
	aiMutex               sync.Mutex // Protects AI calls and payload cache
	conversationHistory   []Message  // Persisted conversation history
	conversationResponses int        // Number of AI responses in current conversation
}

var ScenarioName = "aitx"
var ScenarioDefaultOptions = ScenarioOptions{
	TotalCount:  0,
	Throughput:  10,
	MaxPending:  0,
	MaxWallets:  0,
	Rebroadcast: 1,
	BaseFee:     20,
	TipFee:      2,
	GasLimit:    5000000,
	Timeout:     "",
	ClientGroup: "",
	LogTxs:      false,

	// AI defaults
	OpenRouterAPIKey:   "",
	Model:              "anthropic/claude-3.5-sonnet",
	TestDirection:      "",
	PayloadsPerRequest: 50,
	MaxAICalls:         10,
	MaxTokens:          100000,

	// Generation defaults
	GenerationMode: "geas",

	// Feedback defaults
	FeedbackBatchSize:  20,
	EnableFeedbackLoop: true,

	// Debug defaults
	LogAIConversations: false,
}

var ScenarioDescriptor = scenario.Descriptor{
	Name:           ScenarioName,
	Description:    "AI-powered transaction generator using OpenRouter for diverse test payloads",
	DefaultOptions: ScenarioDefaultOptions,
	NewScenario:    newScenario,
}

func newScenario(logger logrus.FieldLogger) scenario.Scenario {
	return &Scenario{
		options: ScenarioDefaultOptions,
		logger:  logger.WithField("scenario", ScenarioName),
	}
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
	flags.Uint64VarP(&s.options.TotalCount, "count", "c", ScenarioDefaultOptions.TotalCount, "Total number of AI transactions to send")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t", ScenarioDefaultOptions.Throughput, "Number of AI transactions to send per slot")
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", ScenarioDefaultOptions.Rebroadcast, "Enable reliable rebroadcast system")
	flags.Float64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use in transactions (in gwei)")
	flags.Float64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use in transactions (in gwei)")
	flags.Uint64Var(&s.options.GasLimit, "gaslimit", ScenarioDefaultOptions.GasLimit, "Gas limit to use in transactions")
	flags.StringVar(&s.options.Timeout, "timeout", ScenarioDefaultOptions.Timeout, "Timeout for the scenario (e.g. '1h', '30m', '5s') - empty means no timeout")
	flags.StringVar(&s.options.ClientGroup, "client-group", ScenarioDefaultOptions.ClientGroup, "Client group to use for sending transactions")
	flags.BoolVar(&s.options.LogTxs, "log-txs", ScenarioDefaultOptions.LogTxs, "Log all submitted transactions")

	// AI-specific flags
	flags.StringVar(&s.options.OpenRouterAPIKey, "openrouter-api-key", ScenarioDefaultOptions.OpenRouterAPIKey, "OpenRouter API key (can also use OPENROUTER_API_KEY env var)")
	flags.StringVar(&s.options.Model, "model", ScenarioDefaultOptions.Model, "AI model to use for generation")
	flags.StringVar(&s.options.TestDirection, "test-direction", ScenarioDefaultOptions.TestDirection, "Directional guidance for AI test generation")
	flags.Uint64Var(&s.options.PayloadsPerRequest, "payloads-per-request", ScenarioDefaultOptions.PayloadsPerRequest, "Number of payload templates to generate per AI request")
	flags.Uint64Var(&s.options.MaxAICalls, "max-ai-calls", ScenarioDefaultOptions.MaxAICalls, "Maximum number of AI API calls to make")
	flags.Uint64Var(&s.options.MaxTokens, "max-tokens", ScenarioDefaultOptions.MaxTokens, "Maximum total tokens to consume")

	// Generation flags (removed - only geas init_run is supported)

	// Feedback flags
	flags.Uint64Var(&s.options.FeedbackBatchSize, "feedback-batch-size", ScenarioDefaultOptions.FeedbackBatchSize, "Number of transaction results to include in feedback")
	flags.BoolVar(&s.options.EnableFeedbackLoop, "enable-feedback-loop", ScenarioDefaultOptions.EnableFeedbackLoop, "Enable result feedback to AI for learning")

	// Debug flags
	flags.BoolVar(&s.options.LogAIConversations, "log-ai-conversations", ScenarioDefaultOptions.LogAIConversations, "Enable detailed logging of AI conversations for debugging")

	return nil
}

func (s *Scenario) Init(options *scenario.Options) error {
	s.walletPool = options.WalletPool

	if options.Config != "" {
		err := scenario.ParseAndValidateConfig(&ScenarioDescriptor, options.Config, &s.options, s.logger)
		if err != nil {
			return err
		}
	}

	// Validate AI configuration
	if s.options.OpenRouterAPIKey == "" {
		return fmt.Errorf("OpenRouter API key is required (use --openrouter-api-key flag or OPENROUTER_API_KEY env var)")
	}

	// Generation mode is fixed to "geas" for init_run only

	// Configure wallet count
	if s.options.MaxWallets > 0 {
		s.walletPool.SetWalletCount(s.options.MaxWallets)
	} else if s.options.TotalCount > 0 {
		maxWallets := s.options.TotalCount / 50
		if maxWallets < 10 {
			maxWallets = 10
		} else if maxWallets > 1000 {
			maxWallets = 1000
		}
		s.walletPool.SetWalletCount(maxWallets)
	} else {
		if s.options.Throughput*10 < 1000 {
			s.walletPool.SetWalletCount(s.options.Throughput * 10)
		} else {
			s.walletPool.SetWalletCount(1000)
		}
	}

	if s.options.TotalCount == 0 && s.options.Throughput == 0 {
		return fmt.Errorf("neither total count nor throughput limit set, must define at least one of them")
	}

	// Initialize AI components
	s.aiService = NewAIService(s.options.OpenRouterAPIKey, s.options.Model, s.options.LogAIConversations, s.logger)
	s.processor = NewPayloadProcessor(s.logger)
	s.placeholderSubstituter = NewPlaceholderSubstituter(s.walletPool, s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, ""), s.logger)
	s.feedbackCollector = NewFeedbackCollector(s.options.FeedbackBatchSize, s.logger)
	s.geasProcessor = NewGeasProcessor(s.logger)

	// Set AI base prompt based on generation mode
	basePrompt := s.aiService.buildBasePrompt(s.options.GenerationMode)
	s.aiService.SetBasePrompt(basePrompt)

	return nil
}

func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Infof("starting AI transaction generator scenario")
	defer s.logger.Infof("AI transaction generator scenario finished")

	maxPending := s.options.MaxPending
	if maxPending == 0 {
		maxPending = s.options.Throughput * 10
		if maxPending == 0 {
			maxPending = 1000
		}
		if maxPending > s.walletPool.GetConfiguredWalletCount()*10 {
			maxPending = s.walletPool.GetConfiguredWalletCount() * 10
		}
	}

	// Parse timeout duration
	var timeout time.Duration
	if s.options.Timeout != "" {
		var err error
		timeout, err = time.ParseDuration(s.options.Timeout)
		if err != nil {
			return fmt.Errorf("invalid timeout format '%s': %w", s.options.Timeout, err)
		}
	}

	err := scenario.RunTransactionScenario(ctx, scenario.TransactionScenarioOptions{
		TotalCount:                  s.options.TotalCount,
		Throughput:                  s.options.Throughput,
		MaxPending:                  maxPending,
		ThroughputIncrementInterval: 0,
		Timeout:                     timeout,
		WalletPool:                  s.walletPool,

		Logger: s.logger,
		ProcessNextTxFn: func(ctx context.Context, txIdx uint64, onComplete func()) (func(), error) {
			logger := s.logger
			tx, client, wallet, err := s.sendAITransaction(ctx, txIdx, onComplete)
			if client != nil {
				logger = logger.WithField("rpc", client.GetName())
			}
			if tx != nil {
				logger = logger.WithField("nonce", tx.Nonce())
			}
			if wallet != nil {
				logger = logger.WithField("wallet", s.walletPool.GetWalletName(wallet.GetAddress()))
			}

			return func() {
				if err != nil {
					logger.Warnf("could not send AI transaction: %v", err)
				} else if s.options.LogTxs {
					logger.Infof("sent AI tx #%6d: %v", txIdx+1, tx.Hash().String())
				} else {
					logger.Debugf("sent AI tx #%6d: %v", txIdx+1, tx.Hash().String())
				}
			}, err
		},
	})

	return err
}

func (s *Scenario) sendAITransaction(ctx context.Context, txIdx uint64, onComplete func()) (*types.Transaction, *spamoor.Client, *spamoor.Wallet, error) {
	// Deploy a contract and send 10 call transactions using batch sending
	defer onComplete()

	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, int(txIdx), s.options.ClientGroup)
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByPendingTxCount, int(txIdx))

	if client == nil {
		return nil, client, wallet, fmt.Errorf("no client available")
	}

	// Get next payload template from AI or cache
	template, err := s.getNextPayloadTemplate(ctx)
	if err != nil {
		s.logger.Errorf("failed to get AI payload template: %v", err)
		dummyPayload := &PayloadInstance{Type: "geas", Description: "failed_generation"}
		s.feedbackCollector.RecordFailure(dummyPayload, "payload_generation_failed", err.Error())
		return nil, client, wallet, err
	}

	// Substitute placeholders
	payload, err := template.Substitute(s.placeholderSubstituter)
	if err != nil {
		s.logger.Errorf("failed to substitute placeholders: %v", err)
		dummyPayload := &PayloadInstance{Type: template.Type, Description: template.Description}
		s.feedbackCollector.RecordFailure(dummyPayload, "placeholder_substitution_failed", err.Error())
		return nil, client, wallet, err
	}

	// Build deployment transaction
	s.logger.Infof("deploying contract for payload: %s", payload.Description)
	deployTx, contractAddress, err := s.deployGeasContract(ctx, wallet, client, payload)
	if err != nil {
		s.logger.Errorf("failed to deploy contract: %v", err)
		s.feedbackCollector.RecordFailure(payload, "deployment_failed", err.Error())
		return nil, client, wallet, err
	}

	// Deploy contract and wait for confirmation using SendAndAwaitTransaction
	deployReceipt, err := s.walletPool.GetTxPool().SendAndAwaitTransaction(ctx, wallet, deployTx, &spamoor.SendTransactionOptions{
		Client:      client,
		Rebroadcast: s.options.Rebroadcast > 0,
		LogFn:       spamoor.GetDefaultLogFn(s.logger, "deploy", fmt.Sprintf("%6d", txIdx+1), deployTx),
	})

	if err != nil {
		s.logger.Errorf("failed to deploy contract: %v", err)
		s.feedbackCollector.RecordFailure(payload, "deployment_failed", err.Error())
		return nil, client, wallet, err
	}

	if deployReceipt.Status != 1 {
		s.logger.Errorf("contract deployment failed (status: %d)", deployReceipt.Status)
		s.feedbackCollector.RecordFailure(payload, "deployment_reverted", "deployment transaction reverted")
		return deployTx, client, wallet, nil
	}

	s.logger.Infof("contract deployed successfully at %s for payload: %s", contractAddress.Hex(), payload.Description)

	// Build 10 call transactions
	var callTxs []*types.Transaction
	for i := 0; i < 10; i++ {
		callTx, err := s.callGeasContract(ctx, wallet, client, contractAddress, payload)
		if err != nil {
			s.logger.Errorf("failed to build call transaction %d: %v", i+1, err)
			s.feedbackCollector.RecordFailure(payload, "call_build_failed", err.Error())
			return deployTx, client, wallet, err
		}
		callTxs = append(callTxs, callTx)
	}

	// Send all call transactions as a batch from same wallet
	_, err = s.walletPool.GetTxPool().SendTransactionBatch(ctx, wallet, callTxs, &spamoor.BatchOptions{
		SendTransactionOptions: spamoor.SendTransactionOptions{
			Client:      client,
			Rebroadcast: s.options.Rebroadcast > 0,
			OnConfirm: func(tx *types.Transaction, receipt *types.Receipt) {
				// Collect execution results for feedback from call transactions
				s.collectTransactionResult(payload, tx, receipt)
			},
		},
	})

	if err != nil {
		s.logger.Errorf("failed to send call transaction batch: %v", err)
		s.feedbackCollector.RecordFailure(payload, "batch_send_failed", err.Error())
		return deployTx, client, wallet, err
	}

	return deployTx, client, wallet, nil
}

func (s *Scenario) getNextPayloadTemplate(ctx context.Context) (*PayloadTemplate, error) {
	// Lock to ensure only one AI call happens at a time
	s.aiMutex.Lock()
	defer s.aiMutex.Unlock()

	// Check if we have cached payloads
	if s.cacheIndex < len(s.payloadCache) {
		template := s.payloadCache[s.cacheIndex]
		s.cacheIndex++
		return &template, nil
	}

	// Generate new batch of payloads
	if s.aiService.GetCallCount() >= s.options.MaxAICalls {
		return nil, fmt.Errorf("maximum AI calls limit reached (%d)", s.options.MaxAICalls)
	}

	if s.aiService.GetTokenCount() >= s.options.MaxTokens {
		return nil, fmt.Errorf("maximum token limit reached (%d)", s.options.MaxTokens)
	}

	// Check if we need to start a new conversation (after 10 responses)
	if s.conversationResponses >= 10 {
		s.logger.Infof("resetting conversation after %d responses", s.conversationResponses)
		s.conversationHistory = nil
		s.conversationResponses = 0
	}

	if len(s.conversationHistory) == 0 {
		s.logger.Infof("making AI call #%d - starting new conversation (other transactions waiting)", s.aiService.GetCallCount()+1)
	} else {
		s.logger.Infof("making AI call #%d - continuing conversation with %d messages (other transactions waiting)",
			s.aiService.GetCallCount()+1, len(s.conversationHistory))
	}

	// Generate payloads using conversation continuation
	var response *GenerationResponse
	var err error

	if len(s.conversationHistory) == 0 {
		// Start new conversation
		req := GenerationRequest{
			TestDirection:       s.options.TestDirection,
			GenerationMode:      s.options.GenerationMode,
			PayloadCount:        s.options.PayloadsPerRequest,
			PreviousSummary:     "",
			TransactionFeedback: nil,
		}

		// Add feedback if enabled
		if s.options.EnableFeedbackLoop {
			req.TransactionFeedback = s.feedbackCollector.GenerateFeedback()
		}

		response, s.conversationHistory, err = s.aiService.GeneratePayloadsWithConversation(ctx, req, s.processor, nil)
	} else {
		// Continue existing conversation
		feedback := ""
		if s.options.EnableFeedbackLoop {
			txFeedback := s.feedbackCollector.GenerateFeedback()
			if txFeedback != nil {
				feedback = fmt.Sprintf("Transaction feedback: %d total (%d success, %d failed), avg gas: %d. Generate more diverse patterns based on this data.",
					txFeedback.TotalTransactions, txFeedback.SuccessfulTxs, txFeedback.FailedTxs, txFeedback.AverageGasUsed)
			}
		}

		if feedback == "" {
			feedback = fmt.Sprintf("Generate %d more unique geas init_run contracts with different patterns and behaviors.", s.options.PayloadsPerRequest)
		}

		response, s.conversationHistory, err = s.aiService.GeneratePayloadsWithConversation(ctx, GenerationRequest{}, s.processor, &ConversationContinuation{
			History:  s.conversationHistory,
			Feedback: feedback,
		})
	}

	if err != nil {
		return nil, fmt.Errorf("AI payload generation failed: %w", err)
	}

	// Increment conversation response count
	s.conversationResponses++

	// Payloads are already validated by the AI service
	validPayloads := response.Payloads

	// Update cache
	s.payloadCache = validPayloads
	s.cacheIndex = 1 // Return first, set index to second

	if len(validPayloads) == 0 {
		return nil, fmt.Errorf("no valid payloads generated")
	}

	s.logger.Infof("AI call completed, generated %d payloads (conversation: %d responses, cache refilled)",
		len(validPayloads), s.conversationResponses)

	return &validPayloads[0], nil
}

func (s *Scenario) deployGeasContract(ctx context.Context, wallet *spamoor.Wallet, client *spamoor.Client, payload *PayloadInstance) (*types.Transaction, common.Address, error) {
	// Compile geas code
	bytecode, err := s.geasProcessor.CompileGeasPayload(payload)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("geas compilation failed: %w", err)
	}

	// Get suggested fees
	feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(client, s.options.BaseFee, s.options.TipFee)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("failed to get suggested fees: %w", err)
	}

	// Build deployment transaction
	txData, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       s.options.GasLimit,
		To:        nil,               // Contract creation
		Value:     uint256.NewInt(0), // No value for contract deployment
		Data:      bytecode,
	})
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("failed to build transaction data: %w", err)
	}

	tx, err := wallet.BuildDynamicFeeTx(txData)
	if err != nil {
		return nil, common.Address{}, err
	}

	// Calculate contract address
	contractAddr := crypto.CreateAddress(wallet.GetAddress(), tx.Nonce())

	return tx, contractAddr, nil
}

func (s *Scenario) callGeasContract(ctx context.Context, wallet *spamoor.Wallet, client *spamoor.Client, contractAddr common.Address, payload *PayloadInstance) (*types.Transaction, error) {
	// Get suggested fees
	feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(client, s.options.BaseFee, s.options.TipFee)
	if err != nil {
		return nil, fmt.Errorf("failed to get suggested fees: %w", err)
	}

	// Build call transaction with calldata
	txData, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       s.options.GasLimit,
		To:        &contractAddr,
		Value:     uint256.NewInt(0),
		Data:      payload.Calldata,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to build transaction data: %w", err)
	}

	return wallet.BuildDynamicFeeTx(txData)
}

func (s *Scenario) collectTransactionResult(payload *PayloadInstance, tx *types.Transaction, receipt *types.Receipt) {
	if receipt == nil {
		s.feedbackCollector.RecordFailure(payload, "receipt_nil", "receipt was nil")
		return
	}

	// Determine status
	status := "success"
	errorMsg := ""
	if receipt.Status == 0 {
		status = "reverted"
	}

	// Calculate transaction fees
	txFees := utils.GetTransactionFees(tx, receipt)

	// Extract log data
	var logData []string
	for _, log := range receipt.Logs {
		// Convert log data to hex for analysis
		logData = append(logData, fmt.Sprintf("addr:%s topics:%d data:%s",
			log.Address.Hex(),
			len(log.Topics),
			hex.EncodeToString(log.Data)))
	}

	s.logger.Debugf("transaction confirmed: %s (%s) - %s, gas: %d, fees: %s, logs: %d",
		payload.Description, payload.Type, status, receipt.GasUsed, txFees.TotalFeeGweiString(), len(receipt.Logs))

	// Record result for feedback
	result := TransactionResult{
		PayloadType:        payload.Type,
		PayloadDescription: payload.Description,
		Status:             status,
		GasUsed:            receipt.GasUsed,
		BlockExecTime:      "N/A", // Placeholder for external system
		ErrorMessage:       errorMsg,
		LogData:            logData,
	}

	s.feedbackCollector.RecordResult(result)
}
