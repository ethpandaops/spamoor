package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/daemon/configs"
	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/scenarios"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/utils"
)

// RunSpammer represents a spammer instance for the run command
type RunSpammer struct {
	config     configs.SpammerConfig
	scenario   scenario.Scenario
	walletPool *spamoor.WalletPool
	logger     *logrus.Entry
}

// RunCommand handles the run subcommand
func RunCommand(args []string) {
	cliArgs := CliArgs{}
	var selectedSpammers []int

	flags := pflag.NewFlagSet("run", pflag.ExitOnError)
	flags.BoolVarP(&cliArgs.verbose, "verbose", "v", false, "Run with verbose output")
	flags.BoolVar(&cliArgs.trace, "trace", false, "Run with tracing output")
	flags.StringArrayVarP(&cliArgs.rpchosts, "rpchost", "h", []string{}, "The RPC host to send transactions to.")
	flags.StringVar(&cliArgs.rpchostsFile, "rpchost-file", "", "File with a list of RPC hosts to send transactions to.")
	flags.StringVarP(&cliArgs.privkey, "privkey", "p", "", "The private key of the wallet to send funds from.")
	flags.IntSliceVarP(&selectedSpammers, "spammers", "s", []int{}, "Indexes of spammers to run (0-based). If not specified, runs all spammers.")

	flags.Parse(args)

	if flags.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Error: missing YAML file path\n")
		fmt.Fprintf(os.Stderr, "Usage: spamoor run <yaml-file> [options]\n")
		os.Exit(1)
	}

	yamlFile := flags.Args()[0]

	// Set log level
	if cliArgs.trace {
		logrus.SetLevel(logrus.TraceLevel)
	} else if cliArgs.verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}

	// Initialize logger
	logger := logrus.StandardLogger()
	logger.WithFields(logrus.Fields{
		"version":   utils.GetBuildVersion(),
		"buildtime": utils.BuildTime,
	}).Infof("starting spamoor run command")

	// Create context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		<-sigChan
		logger.Info("Received interrupt signal, shutting down...")
		cancel()
	}()

	// Initialize client pool
	rpcHosts := []string{}
	for _, rpcHost := range cliArgs.rpchosts {
		if rpcHost != "" {
			rpcHosts = append(rpcHosts, rpcHost)
		}
	}

	if cliArgs.rpchostsFile != "" {
		fileLines, err := utils.ReadFileLinesTrimmed(cliArgs.rpchostsFile)
		if err != nil {
			logger.WithError(err).Fatal("Failed to read RPC hosts file")
		}
		rpcHosts = append(rpcHosts, fileLines...)
	}

	if len(rpcHosts) == 0 {
		logger.Fatal("No RPC hosts specified")
	}

	clientPool := spamoor.NewClientPool(ctx, rpcHosts, logger.WithField("module", "clientpool"))
	err := clientPool.PrepareClients()
	if err != nil {
		logger.WithError(err).Fatal("Failed to prepare clients")
	}

	// Initialize root wallet
	client := clientPool.GetClient(spamoor.SelectClientRandom, 0, "")
	if client == nil {
		logger.Fatal("No client available")
	}

	rootWallet, err := spamoor.InitRootWallet(ctx, cliArgs.privkey, client, logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to init root wallet")
	}

	// Load and parse YAML file using scenario config logic
	spammerConfigs, err := configs.ResolveConfigImports(yamlFile, "", make(map[string]bool))
	if err != nil {
		logger.WithError(err).Fatal("Failed to load spammer configurations")
	}

	if len(spammerConfigs) == 0 {
		logger.Fatal("No spammer configurations found in YAML file")
	}

	// Validate configurations
	for i, config := range spammerConfigs {
		if config.Scenario == "" {
			logger.Fatalf("Spammer %d: missing scenario", i)
		}

		scenario := scenarios.GetScenario(config.Scenario)
		if scenario == nil {
			logger.Fatalf("Spammer %d: unknown scenario '%s'", i, config.Scenario)
		}

		if config.Name == "" {
			config.Name = fmt.Sprintf("Spammer %d", i+1)
			spammerConfigs[i] = config
		}
	}

	// Filter spammers if specified
	var configsToRun []configs.SpammerConfig
	if len(selectedSpammers) > 0 {
		for _, idx := range selectedSpammers {
			if idx < 0 || idx >= len(spammerConfigs) {
				logger.Fatalf("Invalid spammer index %d (valid range: 0-%d)", idx, len(spammerConfigs)-1)
			}
			configsToRun = append(configsToRun, spammerConfigs[idx])
		}
	} else {
		configsToRun = spammerConfigs
	}

	logger.Infof("Preparing to run %d spammer(s)", len(configsToRun))

	// Create wallet pools slice for txpool
	walletPools := []*spamoor.WalletPool{}
	walletPoolMutex := &sync.RWMutex{}

	// Initialize transaction pool
	txpool := spamoor.NewTxPool(&spamoor.TxPoolOptions{
		Context:    ctx,
		ClientPool: clientPool,
		GetActiveWalletPools: func() []*spamoor.WalletPool {
			walletPoolMutex.RLock()
			defer walletPoolMutex.RUnlock()
			return walletPools
		},
	})

	// Create and initialize spammers
	spammers := make([]*RunSpammer, len(configsToRun))
	for i, config := range configsToRun {
		spammer, err := createSpammer(ctx, config, rootWallet, clientPool, txpool, logger, i)
		if err != nil {
			logger.WithError(err).Fatalf("Failed to create spammer %d (%s)", i, config.Name)
		}
		spammers[i] = spammer

		// Add wallet pool to the slice
		walletPoolMutex.Lock()
		walletPools = append(walletPools, spammer.walletPool)
		walletPoolMutex.Unlock()
	}

	// Prepare all wallet pools
	logger.Info("Preparing wallets for all spammers...")
	var wg sync.WaitGroup
	for i, spammer := range spammers {
		wg.Add(1)
		go func(idx int, s *RunSpammer) {
			defer wg.Done()
			err := s.walletPool.PrepareWallets()
			if err != nil {
				logger.WithError(err).Fatalf("Failed to prepare wallets for spammer %d (%s)", idx, s.config.Name)
			}
			logger.Infof("Wallets prepared for spammer %d (%s)", idx, s.config.Name)
		}(i, spammer)
	}
	wg.Wait()

	// Run all spammers concurrently
	logger.Info("Starting all spammers...")
	errChan := make(chan error, len(spammers))

	for i, spammer := range spammers {
		wg.Add(1)
		go func(idx int, s *RunSpammer) {
			defer wg.Done()
			logger.Infof("Starting spammer %d: %s (%s)", idx, s.config.Name, s.config.Scenario)
			err := s.scenario.Run(ctx)
			if err != nil {
				errChan <- fmt.Errorf("spammer %d (%s) failed: %w", idx, s.config.Name, err)
			}
		}(i, spammer)
	}

	// Wait for context cancellation or errors
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Monitor for errors
	for err := range errChan {
		if err != nil {
			logger.WithError(err).Error("Spammer error")
		}
	}

	logger.Info("All spammers completed")
}

// createSpammer creates and initializes a spammer instance
func createSpammer(ctx context.Context, config configs.SpammerConfig, rootWallet *spamoor.RootWallet,
	clientPool *spamoor.ClientPool, txpool *spamoor.TxPool, logger *logrus.Logger, index int) (*RunSpammer, error) {

	// Get scenario descriptor
	descriptor := scenarios.GetScenario(config.Scenario)
	if descriptor == nil {
		return nil, fmt.Errorf("unknown scenario: %s", config.Scenario)
	}

	// Create logger for this spammer
	spammerLogger := logger.WithFields(logrus.Fields{
		"spammer":  index,
		"scenario": config.Scenario,
		"name":     config.Name,
	})

	// Create scenario instance
	scenarioInstance := descriptor.NewScenario(spammerLogger)
	if scenarioInstance == nil {
		return nil, fmt.Errorf("failed to create scenario instance")
	}

	// Create wallet pool for this spammer
	walletPool := spamoor.NewWalletPool(ctx, spammerLogger.WithField("module", "walletpool"),
		rootWallet, clientPool, txpool)

	// Merge configuration with defaults using scenario's method
	configYAML, err := configs.MergeScenarioConfiguration(descriptor, config.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to merge configuration: %w", err)
	}

	// Initialize scenario
	if err := scenarioInstance.Init(&scenario.Options{
		WalletPool: walletPool,
		Config:     configYAML,
		GlobalCfg:  config.Config,
	}); err != nil {
		return nil, fmt.Errorf("failed to initialize scenario: %w", err)
	}

	// Load config to wallet pool
	if err := walletPool.LoadConfig(configYAML); err != nil {
		return nil, fmt.Errorf("failed to load wallet config: %w", err)
	}

	return &RunSpammer{
		config:     config,
		scenario:   scenarioInstance,
		walletPool: walletPool,
		logger:     spammerLogger,
	}, nil
}
