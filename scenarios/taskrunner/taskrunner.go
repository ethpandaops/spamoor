package taskrunner

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/spamoor"
)

type ScenarioOptions struct {
	TotalCount  uint64  `yaml:"total_count"`
	Throughput  uint64  `yaml:"throughput"`
	MaxPending  uint64  `yaml:"max_pending"`
	MaxWallets  uint64  `yaml:"max_wallets"`
	Rebroadcast uint64  `yaml:"rebroadcast"`
	BaseFee     float64 `yaml:"base_fee"`
	TipFee      float64 `yaml:"tip_fee"`

	TasksConfig string `yaml:"tasks_config"` // Inline YAML/JSON task configuration
	TasksFile   string `yaml:"tasks_file"`   // Path to task configuration file or URL
	AwaitTxs    bool   `yaml:"await_txs"`    // Send and await each transaction individually

	Timeout     string `yaml:"timeout"`
	ClientGroup string `yaml:"client_group"`
	LogTxs      bool   `yaml:"log_txs"`
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	initTasks      []Task
	executionTasks []Task
	initRegistry   *ContractRegistry // Permanent registry for init phase contracts
}

var ScenarioName = "taskrunner"
var ScenarioDefaultOptions = ScenarioOptions{
	TotalCount:  0,
	Throughput:  10,
	MaxPending:  0,
	MaxWallets:  0,
	Rebroadcast: 1,
	BaseFee:     20,
	TipFee:      2,
	TasksConfig: "",
	TasksFile:   "",
	AwaitTxs:    false,
	Timeout:     "",
	ClientGroup: "",
	LogTxs:      false,
}

var ScenarioDescriptor = scenario.Descriptor{
	Name:           ScenarioName,
	Description:    "Execute configurable task sequences with initialization and recurring execution phases",
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
	flags.Uint64VarP(&s.options.TotalCount, "count", "c", ScenarioDefaultOptions.TotalCount, "Total number of task execution cycles to run")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t", ScenarioDefaultOptions.Throughput, "Number of task execution cycles per slot")
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", ScenarioDefaultOptions.Rebroadcast, "Enable reliable rebroadcast with unlimited retries and exponential backoff")
	flags.Float64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use in transactions (in gwei)")
	flags.Float64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use in transactions (in gwei)")

	flags.StringVar(&s.options.TasksConfig, "tasks", ScenarioDefaultOptions.TasksConfig, "Inline task configuration (YAML/JSON string)")
	flags.StringVar(&s.options.TasksFile, "tasks-file", ScenarioDefaultOptions.TasksFile, "Path to task configuration file or URL (http/https)")
	flags.BoolVar(&s.options.AwaitTxs, "await-txs", ScenarioDefaultOptions.AwaitTxs, "Send and await each transaction individually instead of batching")

	flags.StringVar(&s.options.ClientGroup, "client-group", ScenarioDefaultOptions.ClientGroup, "Client group to use for sending transactions")
	flags.StringVar(&s.options.Timeout, "timeout", ScenarioDefaultOptions.Timeout, "Timeout for the scenario (e.g. '1h', '30m', '5s') - empty means no timeout")
	flags.BoolVar(&s.options.LogTxs, "log-txs", ScenarioDefaultOptions.LogTxs, "Log all submitted transactions")

	return nil
}

func (s *Scenario) Init(options *scenario.Options) error {
	s.walletPool = options.WalletPool

	if options.Config != "" {
		// Use the generalized config validation and parsing helper
		err := scenario.ParseAndValidateConfig(&ScenarioDescriptor, options.Config, &s.options, s.logger)
		if err != nil {
			return err
		}
	}

	// Load task configuration
	tasksData, err := s.loadTasksConfig()
	if err != nil {
		return fmt.Errorf("failed to load tasks configuration: %w", err)
	}

	// Parse and validate task configuration
	config, err := ParseTasksConfig(tasksData)
	if err != nil {
		return fmt.Errorf("failed to parse tasks configuration: %w", err)
	}

	s.initTasks = config.InitTasks
	s.executionTasks = config.ExecutionTasks

	if len(s.executionTasks) == 0 {
		return fmt.Errorf("no execution tasks defined - taskrunner requires at least one execution task")
	}

	// Initialize contract registry for init phase
	s.initRegistry = NewContractRegistry()

	// Configure wallet pool
	if err := s.configureWalletPool(); err != nil {
		return fmt.Errorf("failed to configure wallet pool: %w", err)
	}

	// Register well-known wallet for init phase
	s.walletPool.AddWellKnownWallet(&spamoor.WellKnownWalletConfig{
		Name:          "taskrunner-init",
		RefillAmount:  uint256.NewInt(5000000000000000000), // 5 ETH
		RefillBalance: uint256.NewInt(2000000000000000000), // 2 ETH
	})

	s.logger.Infof("initialized taskrunner with %d init tasks and %d execution tasks",
		len(s.initTasks), len(s.executionTasks))

	return nil
}

func (s *Scenario) loadTasksConfig() ([]byte, error) {
	// Priority: inline config > file/URL
	if s.options.TasksConfig != "" {
		return []byte(s.options.TasksConfig), nil
	}

	if s.options.TasksFile != "" {
		// Check if it's a URL
		if strings.HasPrefix(s.options.TasksFile, "http://") || strings.HasPrefix(s.options.TasksFile, "https://") {
			return s.loadFromURL(s.options.TasksFile)
		}
		// Otherwise treat as local file
		return os.ReadFile(s.options.TasksFile)
	}

	return nil, fmt.Errorf("no tasks configuration provided - use --tasks for inline config or --tasks-file for file/URL")
}

func (s *Scenario) loadFromURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error %d when fetching from URL", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return data, nil
}

func (s *Scenario) configureWalletPool() error {
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

	return nil
}

func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Infof("starting scenario: %s", ScenarioName)
	defer s.logger.Infof("scenario %s finished.", ScenarioName)

	// Phase 1: Execute initialization tasks (if any)
	if len(s.initTasks) > 0 {
		s.logger.Infof("executing %d initialization tasks using well-known wallet", len(s.initTasks))

		if err := s.executeInitTasks(ctx); err != nil {
			return fmt.Errorf("initialization phase failed: %w", err)
		}

		s.logger.Infof("initialization phase completed successfully")

		// Log registered contracts
		contracts := s.initRegistry.ListContracts()
		if len(contracts) > 0 {
			s.logger.Infof("registered %d contracts during initialization:", len(contracts))
			for name, addr := range contracts {
				s.logger.Infof("  %s: %s", name, addr.Hex())
			}
		}
	}

	// Phase 2: Execute recurring tasks
	s.logger.Infof("starting recurring execution with %d tasks using numbered wallets", len(s.executionTasks))

	return s.executeRecurringTasks(ctx)
}

func (s *Scenario) executeInitTasks(ctx context.Context) error {
	// Use the registered well-known wallet for init tasks to avoid conflicts with numbered wallets
	wallet := s.walletPool.GetWellKnownWallet("taskrunner-init")
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, s.options.ClientGroup)
	if client == nil {
		return fmt.Errorf("no client available")
	}

	// Execute all init tasks using the unified sequence method
	return s.executeTaskSequence(ctx, s.initTasks, 0, wallet, client, s.initRegistry, 0)
}

// executeTaskSequence executes a sequence of tasks with unified processing logic
func (s *Scenario) executeTaskSequence(ctx context.Context, tasks []Task, baseTaskIndex int, wallet *spamoor.Wallet, client *spamoor.Client, registry *ContractRegistry, txIdx uint64) error {
	// Create execution context with scenario fee settings
	execCtx := &TaskExecutionContext{
		BaseFee: s.options.BaseFee,
		TipFee:  s.options.TipFee,
		TxPool:  s.walletPool.GetTxPool(),
	}

	// Handle transaction processing based on await-txs mode
	if s.options.AwaitTxs {
		// Build, send and await each transaction individually to avoid nonce gaps
		for i, task := range tasks {
			taskName := task.GetName()
			if taskName == "" {
				taskName = task.GetType()
			}

			// Build transaction with proper placeholder processing
			tx, err := s.buildTaskTransaction(ctx, task, baseTaskIndex+i, wallet, registry, execCtx, txIdx)
			if err != nil {
				return fmt.Errorf("failed to build transaction for task %d (%s): %w",
					baseTaskIndex+i+1, task.GetType(), err)
			}

			// Send and await this transaction immediately
			receipt, err := s.walletPool.GetTxPool().SendAndAwaitTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
				Client:      client,
				ClientGroup: s.options.ClientGroup,
				Rebroadcast: s.options.Rebroadcast > 0,
			})
			if err != nil {
				s.logger.Warnf("task %d (%s) failed: %v", baseTaskIndex+i+1, taskName, err)
				wallet.ResetPendingNonce(ctx, client)
				return err
			}

			if s.options.LogTxs {
				s.logger.Infof("task %d/%d (%s) confirmed: %v (block #%v)",
					baseTaskIndex+i+1, len(tasks), taskName, tx.Hash().String(), receipt.BlockNumber.String())
			} else {
				s.logger.Debugf("task %d/%d (%s) confirmed: %v (block #%v)",
					baseTaskIndex+i+1, len(tasks), taskName, tx.Hash().String(), receipt.BlockNumber.String())
			}
		}

		return nil

	} else {
		// Batch mode: Build all transactions first, then send as batch
		var transactions []*types.Transaction
		for i, task := range tasks {
			tx, err := s.buildTaskTransaction(ctx, task, baseTaskIndex+i, wallet, registry, execCtx, txIdx)
			if err != nil {
				return fmt.Errorf("failed to build transaction for task %d (%s): %w",
					baseTaskIndex+i+1, task.GetType(), err)
			}
			transactions = append(transactions, tx)
		}

		// Send transactions as batch
		if len(transactions) == 1 {
			// Single transaction
			tx := transactions[0]
			taskName := tasks[0].GetName()
			if taskName == "" {
				taskName = tasks[0].GetType()
			}

			receipt, err := s.walletPool.GetTxPool().SendAndAwaitTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
				Client:      client,
				ClientGroup: s.options.ClientGroup,
				Rebroadcast: s.options.Rebroadcast > 0,
			})
			if err != nil {
				wallet.ResetPendingNonce(ctx, client)
				return err
			}

			if receipt == nil {
				return fmt.Errorf("transaction receipt not received")
			}

			if s.options.LogTxs {
				s.logger.Infof("task %d/%d (%s) confirmed: %v (block #%v)",
					baseTaskIndex+1, len(tasks), taskName, tx.Hash().String(), receipt.BlockNumber.String())
			} else {
				s.logger.Debugf("task %d/%d (%s) confirmed: %v (block #%v)",
					baseTaskIndex+1, len(tasks), taskName, tx.Hash().String(), receipt.BlockNumber.String())
			}

			return nil

		} else {
			// Multiple transactions - use batch sending
			receipts, err := s.walletPool.GetTxPool().SendTransactionBatch(ctx, wallet, transactions, &spamoor.BatchOptions{
				SendTransactionOptions: spamoor.SendTransactionOptions{
					Client:      client,
					ClientGroup: s.options.ClientGroup,
					Rebroadcast: s.options.Rebroadcast > 0,
				},
			})
			if err != nil {
				wallet.ResetPendingNonce(ctx, client)
				return err
			}

			// Log batch completion
			if s.options.LogTxs {
				s.logger.Infof("batch %d tasks confirmed: %d transactions (block #%v)",
					len(tasks), len(transactions), receipts[0].BlockNumber.String())
			} else {
				s.logger.Debugf("batch %d tasks confirmed: %d transactions (block #%v)",
					len(tasks), len(transactions), receipts[0].BlockNumber.String())
			}

			return nil
		}
	}
}

func (s *Scenario) executeRecurringTasks(ctx context.Context) error {
	// Calculate max pending
	maxPending := s.options.MaxPending
	if maxPending == 0 {
		maxPending = s.options.Throughput * 10
		if maxPending == 0 {
			maxPending = 4000
		}
		if maxPending > s.walletPool.GetConfiguredWalletCount()*10 {
			maxPending = s.walletPool.GetConfiguredWalletCount() * 10
		}
	}

	// Parse timeout
	var timeout time.Duration
	if s.options.Timeout != "" {
		var err error
		timeout, err = time.ParseDuration(s.options.Timeout)
		if err != nil {
			return fmt.Errorf("invalid timeout value: %v", err)
		}
		s.logger.Infof("timeout set to %v", timeout)
	}

	// Run the transaction scenario
	return scenario.RunTransactionScenario(ctx, scenario.TransactionScenarioOptions{
		TotalCount:                  s.options.TotalCount,
		Throughput:                  s.options.Throughput,
		MaxPending:                  maxPending,
		ThroughputIncrementInterval: 0,
		Timeout:                     timeout,
		WalletPool:                  s.walletPool,
		Logger:                      s.logger,
		ProcessNextTxFn:             s.processExecutionTx,
	})
}

func (s *Scenario) processExecutionTx(ctx context.Context, txIdx uint64, onComplete func()) (func(), error) {
	// Create execution-scoped registry that inherits from init registry
	execRegistry := s.initRegistry.Clone()

	// Get wallet and client for this transaction
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(txIdx))
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, int(txIdx), s.options.ClientGroup)

	if client == nil {
		onComplete()
		return nil, fmt.Errorf("no client available")
	}

	// Execute all execution tasks using the unified method
	err := s.executeTaskSequence(ctx, s.executionTasks, 0, wallet, client, execRegistry, txIdx)
	if err != nil {
		onComplete()
		return nil, err
	}

	onComplete()
	return func() {}, nil // No additional callback needed since unified method handles all processing
}

// buildTaskTransaction builds a transaction for a task with proper placeholder processing
func (s *Scenario) buildTaskTransaction(ctx context.Context, task Task, taskIndex int, wallet *spamoor.Wallet, registry *ContractRegistry, execCtx *TaskExecutionContext, txIdx uint64) (*types.Transaction, error) {
	// For call tasks, we need to process placeholders with proper context
	if callTask, ok := task.(*CallTask); ok {
		// Create a copy of the task with processed placeholders
		processedTask := *callTask

		// Process target placeholders
		processedTarget, err := callTask.processPlaceholders(callTask.Target, registry, txIdx, taskIndex)
		if err != nil {
			return nil, fmt.Errorf("failed to process target placeholders: %w", err)
		}
		processedTask.Target = processedTarget

		// Process calldata placeholders if it exists
		if callTask.CallData != "" {
			processedCallData, err := callTask.processPlaceholders(callTask.CallData, registry, txIdx, taskIndex)
			if err != nil {
				return nil, fmt.Errorf("failed to process calldata placeholders: %w", err)
			}
			processedTask.CallData = processedCallData
		}

		// Process arguments if they exist
		if len(callTask.CallArgs) > 0 {
			processedArgs, err := callTask.processArguments(callTask.CallArgs, wallet, registry, txIdx, taskIndex)
			if err != nil {
				return nil, fmt.Errorf("failed to process arguments: %w", err)
			}
			processedTask.CallArgs = processedArgs
		}

		return processedTask.BuildTransaction(ctx, wallet, registry, execCtx)
	} else if deployTask, ok := task.(*DeployTask); ok {
		// For deploy tasks, we need to process placeholders with proper context
		// Create a copy of the task with processed placeholders
		processedTask := *deployTask

		// Process contract code placeholders if it exists
		if deployTask.ContractCode != "" {
			processedCode, err := deployTask.processPlaceholders(deployTask.ContractCode, registry, txIdx, taskIndex)
			if err != nil {
				return nil, fmt.Errorf("failed to process contract code placeholders: %w", err)
			}
			processedTask.ContractCode = processedCode
		}

		// Process contract args placeholders if it exists
		if deployTask.ContractArgs != "" {
			processedArgs, err := deployTask.processPlaceholders(deployTask.ContractArgs, registry, txIdx, taskIndex)
			if err != nil {
				return nil, fmt.Errorf("failed to process contract args placeholders: %w", err)
			}
			processedTask.ContractArgs = processedArgs
		}

		return processedTask.BuildTransaction(ctx, wallet, registry, execCtx)
	} else {
		// For other task types, build normally
		return task.BuildTransaction(ctx, wallet, registry, execCtx)
	}
}
