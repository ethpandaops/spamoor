package daemon

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/ethpandaops/spamoor/daemon/db"
	"github.com/ethpandaops/spamoor/daemon/logscope"
	"github.com/ethpandaops/spamoor/scenarios"
	"github.com/ethpandaops/spamoor/scenariotypes"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/holiman/uint256"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type SpammerStatus int

const (
	SpammerStatusPaused SpammerStatus = iota
	SpammerStatusRunning
	SpammerStatusFinished
	SpammerStatusFailed
)

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

type SpammerConfig struct {
	Seed           string       `yaml:"seed"`
	RefillAmount   *uint256.Int `yaml:"refill_amount"`
	RefillBalance  *uint256.Int `yaml:"refill_balance"`
	RefillInterval uint64       `yaml:"refill_interval"`
	WalletCount    int          `yaml:"wallet_count"`
}

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

func (s *Spammer) Start() error {
	s.dbEntity.Status = int(SpammerStatusRunning)
	err := s.daemon.db.RunDBTransaction(func(tx *sqlx.Tx) error {
		return s.daemon.db.UpdateSpammer(tx, s.dbEntity)
	})
	if err != nil {
		return fmt.Errorf("failed to update spammer: %w", err)
	}

	go s.runScenario()

	return nil
}

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
	scenarioOptions := &scenariotypes.ScenarioOptions{
		WalletPool: s.walletPool,
		Config:     s.dbEntity.Config,
		GlobalCfg:  s.daemon.GetGlobalCfg(),
	}
	scenario := scenarioDescriptor.NewScenario(s.logger)

	err := scenario.Init(scenarioOptions)
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

	err = s.walletPool.PrepareWallets(true)
	if err != nil {
		scenarioErr = fmt.Errorf("failed to prepare wallets: %w", err)
		return
	}

	scenarioErr = scenario.Run(s.scenarioCtx)
}

func (s *Spammer) GetID() int64 {
	return s.dbEntity.ID
}

func (s *Spammer) GetName() string {
	return s.dbEntity.Name
}

func (s *Spammer) GetDescription() string {
	return s.dbEntity.Description
}

func (s *Spammer) GetScenario() string {
	return s.dbEntity.Scenario
}

func (s *Spammer) GetStatus() int {
	return s.dbEntity.Status
}

func (s *Spammer) GetConfig() string {
	return s.dbEntity.Config
}

func (s *Spammer) GetCreatedAt() int64 {
	return s.dbEntity.CreatedAt
}

func (s *Spammer) GetLogScope() *logscope.LogScope {
	return s.logscope
}

func (s *Spammer) GetWalletPool() *spamoor.WalletPool {
	return s.walletPool
}
