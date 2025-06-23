package daemon

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/ethpandaops/spamoor/daemon/db"
	"github.com/ethpandaops/spamoor/daemon/logscope"
	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/scenarios"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/holiman/uint256"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// SpammerStatus represents the execution state of a spammer instance.
// It tracks whether the spammer is paused, running, finished, or failed.
type SpammerStatus int

const (
	// SpammerStatusPaused indicates the spammer is stopped and can be resumed
	SpammerStatusPaused SpammerStatus = iota
	// SpammerStatusRunning indicates the spammer is actively executing its scenario
	SpammerStatusRunning
	// SpammerStatusFinished indicates the spammer completed its scenario successfully
	SpammerStatusFinished
	// SpammerStatusFailed indicates the spammer encountered an error and stopped
	SpammerStatusFailed
)

// Spammer represents a single scenario execution instance with database persistence.
// It manages the lifecycle of a scenario including wallet pool creation, configuration loading,
// and execution monitoring with panic recovery.
type Spammer struct {
	daemon         *Daemon
	logger         logrus.FieldLogger
	logscope       *logscope.LogScope
	dbEntity       *db.Spammer
	walletPool     *spamoor.WalletPool
	scenarioCtx    context.Context
	scenarioCancel context.CancelFunc
	running        bool
	runningChan    chan struct{}
}

// SpammerConfig defines the wallet configuration for a spammer instance.
// This includes the seed for deterministic wallet generation and funding parameters.
type SpammerConfig struct {
	Seed           string       `yaml:"seed"`
	RefillAmount   *uint256.Int `yaml:"refill_amount"`
	RefillBalance  *uint256.Int `yaml:"refill_balance"`
	RefillInterval uint64       `yaml:"refill_interval"`
	WalletCount    int          `yaml:"wallet_count"`
}

// restoreSpammer creates a spammer instance from database entity data.
// It initializes logging with buffered scope, adds the spammer to the daemon map,
// and automatically starts it if it was previously running.
func (d *Daemon) restoreSpammer(dbEntity *db.Spammer) (*Spammer, error) {
	logger := logscope.NewLogger(&logscope.ScopeOptions{
		Parent:     d.logger.WithField("spammer_id", dbEntity.ID),
		BufferSize: 1000,
	})
	logger.GetLogger().SetLevel(logrus.GetLevel())

	spammer := &Spammer{
		daemon:   d,
		dbEntity: dbEntity,
		logscope: logger,
		logger:   logger.GetLogger(),
	}

	d.spammerMapMtx.Lock()
	d.spammerMap[spammer.dbEntity.ID] = spammer
	d.spammerMapMtx.Unlock()

	if spammer.dbEntity.Status == int(SpammerStatusRunning) {
		spammer.Start()
	}

	return spammer, nil
}

// NewSpammer creates a new spammer instance with the specified configuration.
// It validates the config, generates a unique ID (starting from 100), persists to database,
// and optionally starts execution immediately. The ID counter is stored in database state.
func (d *Daemon) NewSpammer(scenarioName string, config string, name string, description string, startImmediately bool) (*Spammer, error) {
	// parse config for sanity checks
	yaml.Unmarshal([]byte(config), &SpammerConfig{})

	// get next id
	d.spammerIdMtx.Lock()
	scenarioCounter := 0
	d.db.GetSpamoorState("scenario_counter", &scenarioCounter)
	if scenarioCounter < 100 {
		scenarioCounter = 100
	} else {
		scenarioCounter++
	}
	d.db.SetSpamoorState(nil, "scenario_counter", scenarioCounter)
	d.spammerIdMtx.Unlock()

	dbEntity := &db.Spammer{
		ID:          int64(scenarioCounter),
		Scenario:    scenarioName,
		Name:        name,
		Description: description,
		Config:      config,
		CreatedAt:   time.Now().Unix(),
	}

	logger := logscope.NewLogger(&logscope.ScopeOptions{
		Parent:     d.logger.WithField("spammer_id", dbEntity.ID),
		BufferSize: 1000,
	})
	logger.GetLogger().SetLevel(logrus.GetLevel())

	spammer := &Spammer{
		daemon:   d,
		dbEntity: dbEntity,
		logscope: logger,
		logger:   logger.GetLogger(),
	}

	err := d.db.RunDBTransaction(func(tx *sqlx.Tx) error {
		return d.db.InsertSpammer(tx, dbEntity)
	})
	if err != nil {
		return nil, err
	}

	d.spammerMapMtx.Lock()
	d.spammerMap[spammer.dbEntity.ID] = spammer
	d.spammerMapMtx.Unlock()

	if startImmediately {
		err = spammer.Start()
		if err != nil {
			return nil, err
		}
	}

	return spammer, nil
}

// Start begins execution of the spammer's scenario in a separate goroutine.
// It updates the status to running in the database and launches the scenario runner.
// Returns an error if the database update fails.
func (s *Spammer) Start() error {
	s.dbEntity.Status = int(SpammerStatusRunning)
	err := s.daemon.db.RunDBTransaction(func(tx *sqlx.Tx) error {
		return s.daemon.db.UpdateSpammer(tx, s.dbEntity)
	})
	if err != nil {
		return fmt.Errorf("failed to update spammer: %w", err)
	}

	// Track spammer start for metrics
	s.daemon.TrackSpammerStatusChange(s.dbEntity.ID, true)

	go s.runScenario()

	return nil
}

// Pause stops the running scenario by canceling its context.
// It waits up to 10 seconds for the scenario to stop gracefully.
// Returns an error if the scenario is not running or fails to stop within the timeout.
func (s *Spammer) Pause() error {
	scenarioCancel := s.scenarioCancel
	if scenarioCancel == nil {
		return fmt.Errorf("scenario is not running")
	}

	scenarioCancel()

	select {
	case <-s.runningChan:
		return nil
	case <-time.After(10 * time.Second):
		return fmt.Errorf("failed to pause spammer: timeout")
	}
}

// runScenario executes the spammer's scenario with comprehensive error handling.
// It manages wallet pool creation, scenario initialization, configuration loading,
// and handles panics with stack trace logging. Updates status in database on completion.
func (s *Spammer) runScenario() {
	if s.running {
		return
	}
	s.running = true
	s.daemon.spammerWg.Add(1)

	if s.scenarioCancel != nil {
		s.scenarioCancel()
	}
	runningChan := make(chan struct{})
	s.runningChan = runningChan
	s.scenarioCtx, s.scenarioCancel = context.WithCancel(s.daemon.ctx)

	var scenarioErr error
	defer func() {
		if err := recover(); err != nil {
			s.logger.Errorf("uncaught panic in spammer subroutine: %v, stack: %v", s.dbEntity.Scenario, err, string(debug.Stack()))
			if err2, ok := err.(error); ok {
				scenarioErr = err2
			} else {
				scenarioErr = fmt.Errorf("unknown panic: %v", err)
			}
		}

		s.daemon.spammerWg.Done()
		close(s.runningChan)

		if s.daemon.ctx.Err() != nil {
			return
		} else if s.scenarioCtx.Err() != nil {
			s.dbEntity.Status = int(SpammerStatusPaused)
		} else if scenarioErr != nil {
			s.logger.Errorf("failed to run scenario: %w", scenarioErr)
			s.dbEntity.Status = int(SpammerStatusFailed)
		} else {
			s.logger.Info("scenario finished successfully")
			s.dbEntity.Status = int(SpammerStatusFinished)
		}

		// Track spammer stop for metrics (not running anymore)
		s.daemon.TrackSpammerStatusChange(s.dbEntity.ID, false)

		err := s.daemon.db.RunDBTransaction(func(tx *sqlx.Tx) error {
			return s.daemon.db.UpdateSpammer(tx, s.dbEntity)
		})
		if err != nil {
			s.logger.Errorf("failed to update spammer: %w", err)
		}

		s.running = false
		s.scenarioCancel()
		s.scenarioCancel = nil
	}()

	scenarioDescriptor := scenarios.GetScenario(s.dbEntity.Scenario)
	if scenarioDescriptor == nil {
		s.logger.Errorf("scenario %s not found", s.dbEntity.Scenario)
		return
	}

	s.walletPool = spamoor.NewWalletPool(s.scenarioCtx, s.logger, s.daemon.rootWallet, s.daemon.clientPool, s.daemon.txpool)
	s.walletPool.SetSpammerID(uint64(s.dbEntity.ID))
	s.walletPool.SetTransactionTracker(s.TrackTransactionResult)
	options := &scenario.Options{
		WalletPool: s.walletPool,
		Config:     s.dbEntity.Config,
		GlobalCfg:  s.daemon.GetGlobalCfg(),
	}
	scenario := scenarioDescriptor.NewScenario(s.logger)

	err := scenario.Init(options)
	if err != nil {
		scenarioErr = fmt.Errorf("failed to init scenario: %w", err)
		return
	}

	if s.dbEntity.Config != "" {
		scenarioErr = s.walletPool.LoadConfig(s.dbEntity.Config)
		if scenarioErr != nil {
			return
		}
	}

	err = s.walletPool.PrepareWallets()
	if err != nil {
		scenarioErr = fmt.Errorf("failed to prepare wallets: %w", err)
		return
	}

	scenarioErr = scenario.Run(s.scenarioCtx)
}

// GetID returns the unique identifier of the spammer instance.
func (s *Spammer) GetID() int64 {
	return s.dbEntity.ID
}

// GetName returns the human-readable name of the spammer instance.
func (s *Spammer) GetName() string {
	return s.dbEntity.Name
}

// GetDescription returns the description text of the spammer instance.
func (s *Spammer) GetDescription() string {
	return s.dbEntity.Description
}

// GetScenario returns the name of the scenario being executed by this spammer.
func (s *Spammer) GetScenario() string {
	return s.dbEntity.Scenario
}

// GetStatus returns the current execution status of the spammer as an integer.
// Use SpammerStatus constants to interpret the returned value.
func (s *Spammer) GetStatus() int {
	return s.dbEntity.Status
}

// GetConfig returns the YAML configuration string used by this spammer.
func (s *Spammer) GetConfig() string {
	return s.dbEntity.Config
}

// GetCreatedAt returns the Unix timestamp when this spammer was created.
func (s *Spammer) GetCreatedAt() int64 {
	return s.dbEntity.CreatedAt
}

// GetLogScope returns the buffered log scope for this spammer.
// This provides access to captured log messages for debugging and monitoring.
func (s *Spammer) GetLogScope() *logscope.LogScope {
	return s.logscope
}

// GetWalletPool returns the wallet pool used by this spammer for transaction submission.
// Returns nil if the spammer has not been started or the wallet pool is not initialized.
func (s *Spammer) GetWalletPool() *spamoor.WalletPool {
	return s.walletPool
}

// TrackTransactionResult records transaction success or failure for metrics
func (s *Spammer) TrackTransactionResult(err error) {
	if err != nil {
		s.daemon.TrackTransactionFailure(s.dbEntity.ID)
	} else {
		s.daemon.TrackTransactionSent(s.dbEntity.ID)
	}
}
