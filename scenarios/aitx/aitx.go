package aitx

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sort"
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

	// Persistence options
	MaxPayloads     int    `yaml:"max_payloads"`
	PersistenceFile string `yaml:"persistence_file"`
	SavePersistence bool   `yaml:"save_persistence"`
	LoadPersistence bool   `yaml:"load_persistence"`

	// Payload management options
	SuccessThreshold int `yaml:"success_threshold"`
}

type PayloadState struct {
	Template        PayloadTemplate
	IsDeployed      bool
	IsDeploying     bool
	ContractAddress common.Address
	SuccessCount    int
	FailCount       int
	LastUsed        time.Time
	BatchID         int        // AI batch ID this payload belongs to
	mutex           sync.Mutex // Protects individual payload state
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

	// Async payload management
	payloadStates         []*PayloadState
	payloadMutex          sync.RWMutex  // Protects payload states slice
	payloadRoundRobin     int           // Round-robin index
	aiRequestChan         chan struct{} // Signals need for more payloads
	aiReadyChan           chan struct{} // Signals AI has returned payloads
	shutdownChan          chan struct{} // Signals shutdown
	conversationHistory   []Message     // Persisted conversation history
	conversationResponses int           // Number of AI responses in current conversation
	aiWorkerRunning       bool          // Tracks if AI worker is running
	aiWorkerMutex         sync.Mutex    // Protects AI worker state
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

	// Persistence defaults
	MaxPayloads:     100,
	PersistenceFile: "",
	SavePersistence: true,
	LoadPersistence: true,

	// Payload management defaults
	SuccessThreshold: 20,
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

	// Persistence flags
	flags.IntVar(&s.options.MaxPayloads, "max-payloads", ScenarioDefaultOptions.MaxPayloads, "Maximum number of payloads to keep in memory")
	flags.StringVar(&s.options.PersistenceFile, "persistence-file", ScenarioDefaultOptions.PersistenceFile, "File to save/load payloads for persistence")
	flags.BoolVar(&s.options.SavePersistence, "save-persistence", ScenarioDefaultOptions.SavePersistence, "Save payloads to persistence file on shutdown")
	flags.BoolVar(&s.options.LoadPersistence, "load-persistence", ScenarioDefaultOptions.LoadPersistence, "Load payloads from persistence file on startup")

	// Payload management flags
	flags.IntVar(&s.options.SuccessThreshold, "success-threshold", ScenarioDefaultOptions.SuccessThreshold, "Number of successful calls before requesting new payloads")

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
	s.placeholderSubstituter = NewPlaceholderSubstituter(s.walletPool, s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, ""), s.logger)
	s.processor = NewPayloadProcessor(s.logger, s.placeholderSubstituter)
	s.feedbackCollector = NewFeedbackCollector(s.options.FeedbackBatchSize, s.logger)
	s.geasProcessor = NewGeasProcessor(s.logger)

	// Set AI base prompt based on generation mode
	basePrompt := s.aiService.buildBasePrompt(s.options.GenerationMode)
	s.aiService.SetBasePrompt(basePrompt)

	// Initialize async payload management
	s.payloadStates = make([]*PayloadState, 0, s.options.MaxPayloads)
	s.aiRequestChan = make(chan struct{}, 1)
	s.aiReadyChan = make(chan struct{}, 1)
	s.shutdownChan = make(chan struct{})

	// Load payloads from persistence file if enabled
	if s.options.LoadPersistence && s.options.PersistenceFile != "" {
		if err := s.loadPayloadsFromFile(); err != nil {
			s.logger.Warnf("failed to load payloads from persistence file: %v", err)
		} else {
			s.logger.Infof("loaded %d payloads from persistence file", len(s.payloadStates))
			// Verify contract deployments
			s.verifyDeployedContracts()
		}
	}

	return nil
}

func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Infof("starting AI transaction generator scenario")
	defer s.logger.Infof("AI transaction generator scenario finished")

	// Start background AI worker
	go s.aiWorker(ctx)

	if len(s.payloadStates) == 0 {
		// Initial AI request to get started
		select {
		case s.aiRequestChan <- struct{}{}:
		default:
		}

		// Wait for AI to be ready
		s.logger.Infof("waiting for AI payloads to be ready")
		select {
		case <-s.aiReadyChan:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	s.logger.Infof("AI payloads ready, starting transaction generation")

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

	// Signal shutdown to AI worker
	close(s.shutdownChan)

	// Save payloads to persistence file if enabled
	if s.options.SavePersistence && s.options.PersistenceFile != "" {
		if saveErr := s.savePayloadsToFile(); saveErr != nil {
			s.logger.Errorf("failed to save payloads to persistence file: %v", saveErr)
		} else {
			s.logger.Infof("saved %d payloads to persistence file", len(s.payloadStates))
		}
	}

	return err
}

func (s *Scenario) sendAITransaction(ctx context.Context, txIdx uint64, onComplete func()) (*types.Transaction, *spamoor.Client, *spamoor.Wallet, error) {
	// Send single call transaction using round-robin payload selection
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, int(txIdx), s.options.ClientGroup)
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByPendingTxCount, int(txIdx))

	txSubmitted := false
	defer func() {
		if !txSubmitted {
			onComplete()
		}
	}()

	if client == nil {
		return nil, client, wallet, fmt.Errorf("no client available")
	}

	// Get next payload using round-robin selection
	payloadState, err := s.getNextPayload(ctx)
	if err != nil {
		s.logger.Errorf("failed to get payload: %v", err)
		return nil, client, wallet, err
	}

	// Substitute placeholders
	payload, err := payloadState.Template.Substitute(s.placeholderSubstituter)
	if err != nil {
		s.logger.Errorf("failed to substitute placeholders: %v", err)
		s.recordPayloadFailure(payloadState, "placeholder_substitution_failed", err.Error())
		return nil, client, wallet, err
	}

	// Handle deployment if needed
	contractAddress, deployTx, err := s.ensureContractDeployed(ctx, payloadState, payload, wallet, client, txIdx)
	if err != nil {
		s.logger.Errorf("failed to deploy contract: %v", err)
		s.recordPayloadFailure(payloadState, "deployment_failed", err.Error())
		return deployTx, client, wallet, err
	}

	// Build call transaction
	callTx, err := s.callGeasContract(ctx, wallet, client, contractAddress, payload)
	if err != nil {
		s.logger.Errorf("failed to build call transaction: %v", err)
		s.recordPayloadFailure(payloadState, "call_build_failed", err.Error())
		return deployTx, client, wallet, err
	}

	// Send call transaction
	txSubmitted = true
	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, callTx, &spamoor.SendTransactionOptions{
		Client:      client,
		Rebroadcast: s.options.Rebroadcast > 0,
		OnConfirm: func(tx *types.Transaction, receipt *types.Receipt) {
			// Record success
			s.recordPayloadSuccess(payloadState, payload, tx, receipt)
		},
		OnComplete: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			onComplete()
			if err != nil {
				s.recordPayloadFailure(payloadState, "call_failed", err.Error())
			}
		},
		LogFn: spamoor.GetDefaultLogFn(s.logger, "call", fmt.Sprintf("%6d", txIdx+1), callTx),
	})

	if err != nil {
		s.logger.Errorf("failed to send call transaction: %v", err)
		s.recordPayloadFailure(payloadState, "call_send_failed", err.Error())
		return callTx, client, wallet, err
	}

	return callTx, client, wallet, nil
}

// getNextPayload selects the next payload using round-robin, preferring payloads with < 20 successes
func (s *Scenario) getNextPayload(ctx context.Context) (*PayloadState, error) {
	s.payloadMutex.Lock()
	defer s.payloadMutex.Unlock()

	// Check for payloads that haven't reached success threshold and haven't failed too much
	for i := 0; i < len(s.payloadStates); i++ {
		idx := (s.payloadRoundRobin + i) % len(s.payloadStates)
		payloadState := s.payloadStates[idx]

		// Skip if currently deploying
		if payloadState.IsDeploying {
			continue
		}

		// Use payload if it hasn't reached success threshold AND hasn't failed excessively
		// A payload is considered "exhausted" if it has reached either:
		// - Success threshold (working well)
		// - Failure threshold (not working, give up)
		totalCalls := payloadState.SuccessCount + payloadState.FailCount
		hasReachedSuccessThreshold := payloadState.SuccessCount >= s.options.SuccessThreshold
		hasReachedFailureThreshold := payloadState.FailCount >= s.options.SuccessThreshold

		if !hasReachedSuccessThreshold && !hasReachedFailureThreshold {
			s.payloadRoundRobin = (idx + 1) % len(s.payloadStates)
			payloadState.LastUsed = time.Now()
			s.logger.Debugf("selected payload: %s (success: %d, fail: %d, total: %d)",
				payloadState.Template.Description, payloadState.SuccessCount, payloadState.FailCount, totalCalls)
			return payloadState, nil
		}
	}

	// No payload available that hasn't reached threshold, request more payloads
	s.requestMorePayloads()

	// If we have any payloads, return the next one
	if len(s.payloadStates) > 0 {
		idx := s.payloadRoundRobin % len(s.payloadStates)
		payloadState := s.payloadStates[idx]
		s.payloadRoundRobin = (idx + 1) % len(s.payloadStates)
		payloadState.LastUsed = time.Now()
		return payloadState, nil
	}

	return nil, fmt.Errorf("no payloads available")
}

// requestMorePayloads signals the AI worker to generate more payloads
func (s *Scenario) requestMorePayloads() {
	select {
	case s.aiRequestChan <- struct{}{}:
		s.logger.Debugf("requested more payloads from AI worker")
	default:
		// Channel full, request already pending
	}
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

func (s *Scenario) recordPayloadSuccess(payloadState *PayloadState, payload *PayloadInstance, tx *types.Transaction, receipt *types.Receipt) {
	payloadState.mutex.Lock()
	batchID := payloadState.BatchID
	payloadState.SuccessCount++
	payloadState.mutex.Unlock()

	if receipt == nil {
		s.feedbackCollector.RecordFailure(payload, "receipt_nil", "receipt was nil", batchID)
		return
	}

	// Determine status
	status := "success"
	errorMsg := ""
	if receipt.Status == 0 {
		status = "reverted"
		// Record as failure instead of success
		payloadState.mutex.Lock()
		payloadState.SuccessCount--
		payloadState.FailCount++
		payloadState.mutex.Unlock()
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

	s.logger.Debugf("transaction confirmed: %s (%s) - %s, gas: %d, fees: %s, logs: %d, success: %d, fail: %d, batch: %d",
		payload.Description, payload.Type, status, receipt.GasUsed, txFees.TotalFeeGweiString(), len(receipt.Logs),
		payloadState.SuccessCount, payloadState.FailCount, batchID)

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

	s.feedbackCollector.RecordResult(result, batchID)
}

// ensureContractDeployed ensures the contract is deployed for the payload
func (s *Scenario) ensureContractDeployed(ctx context.Context, payloadState *PayloadState, payload *PayloadInstance, wallet *spamoor.Wallet, client *spamoor.Client, txIdx uint64) (common.Address, *types.Transaction, error) {
	payloadState.mutex.Lock()
	defer payloadState.mutex.Unlock()

	// Check if already deployed
	if payloadState.IsDeployed {
		return payloadState.ContractAddress, nil, nil
	}

	payloadState.IsDeploying = true
	defer func() {
		payloadState.IsDeploying = false
	}()

	// Deploy the contract
	s.logger.Infof("deploying contract for payload: %s", payload.Description)
	deployTx, contractAddress, err := s.deployGeasContract(ctx, wallet, client, payload)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to build deployment transaction: %w", err)
	}

	// Deploy contract and wait for confirmation
	deployReceipt, err := s.walletPool.GetTxPool().SendAndAwaitTransaction(ctx, wallet, deployTx, &spamoor.SendTransactionOptions{
		Client:      client,
		Rebroadcast: s.options.Rebroadcast > 0,
		LogFn:       spamoor.GetDefaultLogFn(s.logger, "deploy", fmt.Sprintf("%6d", txIdx+1), deployTx),
	})

	if err != nil {
		return common.Address{}, deployTx, fmt.Errorf("failed to deploy contract: %w", err)
	}

	if deployReceipt.Status != 1 {
		return common.Address{}, deployTx, fmt.Errorf("contract deployment failed (status: %d)", deployReceipt.Status)
	}

	// Mark as deployed
	payloadState.IsDeployed = true
	payloadState.ContractAddress = contractAddress

	s.logger.Infof("contract deployed successfully at %s for payload: %s", contractAddress.Hex(), payload.Description)
	return contractAddress, deployTx, nil
}

// aiWorker runs in background to generate payloads asynchronously
func (s *Scenario) aiWorker(ctx context.Context) {
	s.aiWorkerMutex.Lock()
	if s.aiWorkerRunning {
		s.aiWorkerMutex.Unlock()
		return
	}
	s.aiWorkerRunning = true
	s.aiWorkerMutex.Unlock()

	s.logger.Infof("starting AI worker for background payload generation")
	defer s.logger.Infof("AI worker stopped")

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.shutdownChan:
			return
		case <-s.aiRequestChan:
			// Generate new payloads
			if err := s.generatePayloads(ctx); err != nil {
				s.logger.Errorf("failed to generate payloads: %v", err)
				// Wait before retrying
				select {
				case <-time.After(5 * time.Second):
				case <-ctx.Done():
					return
				case <-s.shutdownChan:
					return
				}
			} else {
				select {
				case s.aiReadyChan <- struct{}{}:
				default:
				}
				select {
				case <-s.aiRequestChan:
				default:
				}
			}
		}
	}
}

// generatePayloads generates new payloads from AI and manages the payload pool
func (s *Scenario) generatePayloads(ctx context.Context) error {
	// Check AI limits
	if s.aiService.GetCallCount() >= s.options.MaxAICalls {
		s.logger.Warnf("maximum AI calls limit reached (%d)", s.options.MaxAICalls)
		return fmt.Errorf("maximum AI calls limit reached")
	}

	if s.aiService.GetTokenCount() >= s.options.MaxTokens {
		s.logger.Warnf("maximum token limit reached (%d)", s.options.MaxTokens)
		return fmt.Errorf("maximum token limit reached")
	}

	// Check if we need to start a new conversation (after 10 responses)
	if s.conversationResponses >= 10 {
		s.logger.Infof("resetting conversation after %d responses", s.conversationResponses)
		s.conversationHistory = nil
		s.conversationResponses = 0
	}

	if len(s.conversationHistory) == 0 {
		s.logger.Infof("making AI call #%d - starting new conversation", s.aiService.GetCallCount()+1)
	} else {
		s.logger.Infof("making AI call #%d - continuing conversation with %d messages",
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
		return fmt.Errorf("AI payload generation failed: %w", err)
	}

	// Increment conversation response count
	s.conversationResponses++

	// Start a new feedback batch for the new payloads
	s.feedbackCollector.StartNewBatch()
	batchID := s.feedbackCollector.GetCurrentBatchID()

	// Add new payloads to the pool with the current batch ID
	s.addPayloadsToPool(response.Payloads, batchID)

	s.logger.Infof("AI call completed, generated %d payloads (conversation: %d responses, batch: %d)",
		len(response.Payloads), s.conversationResponses, batchID)

	return nil
}

// addPayloadsToPool adds new payloads to the pool and manages the max payload limit
func (s *Scenario) addPayloadsToPool(templates []PayloadTemplate, batchID int) {
	s.payloadMutex.Lock()
	defer s.payloadMutex.Unlock()

	// Add new payloads with batch ID
	for _, template := range templates {
		payloadState := &PayloadState{
			Template:     template,
			IsDeployed:   false,
			IsDeploying:  false,
			SuccessCount: 0,
			FailCount:    0,
			LastUsed:     time.Now(),
			BatchID:      batchID,
		}
		s.payloadStates = append(s.payloadStates, payloadState)
	}

	// Clean up if we exceed max payloads
	if len(s.payloadStates) > s.options.MaxPayloads {
		s.cleanupPayloads()
	}

	s.logger.Infof("added %d payloads to pool (batch %d), total: %d", len(templates), batchID, len(s.payloadStates))
}

// cleanupPayloads removes failing payloads first, then payloads with highest success count
func (s *Scenario) cleanupPayloads() {
	// Sort by fail count (descending), then by success count (descending)
	sort.Slice(s.payloadStates, func(i, j int) bool {
		if s.payloadStates[i].FailCount != s.payloadStates[j].FailCount {
			return s.payloadStates[i].FailCount > s.payloadStates[j].FailCount
		}
		return s.payloadStates[i].SuccessCount > s.payloadStates[j].SuccessCount
	})

	// Remove the worst 25% to get back to 75% of max
	targetSize := (s.options.MaxPayloads * 3) / 4 // 75% of max
	if len(s.payloadStates) > targetSize {
		removedCount := len(s.payloadStates) - targetSize
		s.payloadStates = s.payloadStates[removedCount:]
		s.logger.Infof("cleaned up %d payloads, remaining: %d", removedCount, len(s.payloadStates))
	}

	// Reset round-robin index if needed
	if s.payloadRoundRobin >= len(s.payloadStates) {
		s.payloadRoundRobin = 0
	}
}

// PayloadStatePersistence represents the data to persist for a payload state
type PayloadStatePersistence struct {
	Template        PayloadTemplate `json:"template"`
	IsDeployed      bool            `json:"is_deployed"`
	ContractAddress string          `json:"contract_address,omitempty"`
	SuccessCount    int             `json:"success_count"`
	FailCount       int             `json:"fail_count"`
	LastUsed        time.Time       `json:"last_used"`
	BatchID         int             `json:"batch_id"`
}

// PayloadsPersistenceData represents the complete persistence data
type PayloadsPersistenceData struct {
	Payloads              []PayloadStatePersistence `json:"payloads"`
	ConversationHistory   []Message                 `json:"conversation_history,omitempty"`
	ConversationResponses int                       `json:"conversation_responses"`
	SavedAt               time.Time                 `json:"saved_at"`
}

// savePayloadsToFile saves the current payload states to a JSON file
func (s *Scenario) savePayloadsToFile() error {
	s.payloadMutex.RLock()
	defer s.payloadMutex.RUnlock()

	var persistenceData PayloadsPersistenceData
	persistenceData.ConversationHistory = s.conversationHistory
	persistenceData.ConversationResponses = s.conversationResponses
	persistenceData.SavedAt = time.Now()

	// Convert payload states to persistence format
	persistenceData.Payloads = make([]PayloadStatePersistence, len(s.payloadStates))
	for i, state := range s.payloadStates {
		state.mutex.Lock()
		persistenceData.Payloads[i] = PayloadStatePersistence{
			Template:        state.Template,
			IsDeployed:      state.IsDeployed,
			ContractAddress: state.ContractAddress.Hex(),
			SuccessCount:    state.SuccessCount,
			FailCount:       state.FailCount,
			LastUsed:        state.LastUsed,
			BatchID:         state.BatchID,
		}
		state.mutex.Unlock()
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(persistenceData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal persistence data: %w", err)
	}

	// Write to file
	err = os.WriteFile(s.options.PersistenceFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write persistence file: %w", err)
	}

	s.logger.Infof("saved %d payloads to persistence file: %s", len(s.payloadStates), s.options.PersistenceFile)
	return nil
}

// loadPayloadsFromFile loads payload states from a JSON file
func (s *Scenario) loadPayloadsFromFile() error {
	// Check if file exists
	if _, err := os.Stat(s.options.PersistenceFile); os.IsNotExist(err) {
		s.logger.Infof("persistence file does not exist: %s", s.options.PersistenceFile)
		return nil
	}

	// Read file
	data, err := os.ReadFile(s.options.PersistenceFile)
	if err != nil {
		return fmt.Errorf("failed to read persistence file: %w", err)
	}

	// Unmarshal JSON
	var persistenceData PayloadsPersistenceData
	err = json.Unmarshal(data, &persistenceData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal persistence data: %w", err)
	}

	s.payloadMutex.Lock()
	defer s.payloadMutex.Unlock()

	// Restore conversation state
	s.conversationHistory = persistenceData.ConversationHistory
	s.conversationResponses = persistenceData.ConversationResponses

	// Convert persistence format to payload states
	s.payloadStates = make([]*PayloadState, len(persistenceData.Payloads))
	for i, persistedState := range persistenceData.Payloads {
		contractAddr := common.Address{}
		if persistedState.ContractAddress != "" && persistedState.ContractAddress != "0x0000000000000000000000000000000000000000" {
			contractAddr = common.HexToAddress(persistedState.ContractAddress)
		}

		s.payloadStates[i] = &PayloadState{
			Template:        persistedState.Template,
			IsDeployed:      persistedState.IsDeployed,
			IsDeploying:     false, // Always start with false
			ContractAddress: contractAddr,
			SuccessCount:    persistedState.SuccessCount,
			FailCount:       persistedState.FailCount,
			LastUsed:        persistedState.LastUsed,
			BatchID:         persistedState.BatchID,
		}
	}

	s.logger.Infof("loaded %d payloads from persistence file: %s (saved at: %s)",
		len(s.payloadStates), s.options.PersistenceFile, persistenceData.SavedAt.Format(time.RFC3339))

	return nil
}

// verifyDeployedContracts checks if contracts are actually deployed at stored addresses
func (s *Scenario) verifyDeployedContracts() {
	if len(s.payloadStates) == 0 {
		return
	}

	// Get a client for verification
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, s.options.ClientGroup)
	if client == nil {
		s.logger.Warnf("no client available for contract verification")
		return
	}

	s.payloadMutex.Lock()
	defer s.payloadMutex.Unlock()

	verifiedCount := 0
	invalidatedCount := 0

	for _, payloadState := range s.payloadStates {
		payloadState.mutex.Lock()

		if payloadState.IsDeployed && payloadState.ContractAddress != (common.Address{}) {
			// Check if code exists at the address
			code, err := client.GetEthClient().CodeAt(context.Background(), payloadState.ContractAddress, nil)
			if err != nil {
				s.logger.Warnf("failed to check contract code at %s: %v", payloadState.ContractAddress.Hex(), err)
				// On error, assume contract is not deployed to be safe
				payloadState.IsDeployed = false
				payloadState.ContractAddress = common.Address{}
				invalidatedCount++
			} else if len(code) == 0 {
				// No code at address, contract not deployed
				s.logger.Debugf("no code found at %s, marking payload as not deployed: %s",
					payloadState.ContractAddress.Hex(), payloadState.Template.Description)
				payloadState.IsDeployed = false
				payloadState.ContractAddress = common.Address{}
				invalidatedCount++
			} else {
				// Code exists, contract is deployed
				s.logger.Debugf("verified contract at %s for payload: %s",
					payloadState.ContractAddress.Hex(), payloadState.Template.Description)
				verifiedCount++
			}
		}

		payloadState.mutex.Unlock()
	}

	if verifiedCount > 0 || invalidatedCount > 0 {
		s.logger.Infof("contract verification complete: %d verified, %d invalidated",
			verifiedCount, invalidatedCount)
	}
}

func (s *Scenario) recordPayloadFailure(payloadState *PayloadState, errorType string, errorMsg string) {
	payloadState.mutex.Lock()
	batchID := payloadState.BatchID
	payloadState.FailCount++
	payloadState.mutex.Unlock()

	// Create dummy payload for feedback
	dummyPayload := &PayloadInstance{
		Type:        payloadState.Template.Type,
		Description: payloadState.Template.Description,
	}
	s.feedbackCollector.RecordFailure(dummyPayload, errorType, errorMsg, batchID)

	s.logger.Debugf("payload failure: %s - %s: %s, success: %d, fail: %d, batch: %d",
		payloadState.Template.Description, errorType, errorMsg,
		payloadState.SuccessCount, payloadState.FailCount, batchID)
}
