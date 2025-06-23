package daemon

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"github.com/ethpandaops/spamoor/daemon/db"
	"github.com/ethpandaops/spamoor/spamoor"
	"gopkg.in/yaml.v3"
)

// Daemon manages multiple spammer instances with database persistence and graceful shutdown.
// It orchestrates the lifecycle of spammers including restoration from database, global configuration
// management, and coordinated shutdown with a 10-second timeout.
type Daemon struct {
	ctx        context.Context
	cancel     context.CancelFunc
	logger     logrus.FieldLogger
	clientPool *spamoor.ClientPool
	rootWallet *spamoor.RootWallet
	txpool     *spamoor.TxPool
	db         *db.Database

	spammerIdMtx  sync.Mutex
	spammerMap    map[int64]*Spammer
	spammerMapMtx sync.RWMutex
	spammerWg     sync.WaitGroup

	globalCfg map[string]interface{}

	// Metrics collector for Prometheus metrics
	metricsCollector *MetricsCollector

	// TxPool metrics collector for advanced transaction metrics
	txPoolMetricsCollector *TxPoolMetricsCollector
}

// NewDaemon creates a new daemon instance with the provided components.
// It initializes a cancellable context, spammer map, and global configuration map.
// The daemon manages spammer instances and handles their lifecycle.
func NewDaemon(parentCtx context.Context, logger logrus.FieldLogger, clientPool *spamoor.ClientPool, rootWallet *spamoor.RootWallet, txpool *spamoor.TxPool, db *db.Database) *Daemon {
	ctx, cancel := context.WithCancel(parentCtx)
	return &Daemon{
		ctx:        ctx,
		cancel:     cancel,
		logger:     logger,
		clientPool: clientPool,
		rootWallet: rootWallet,
		txpool:     txpool,
		db:         db,
		spammerMap: make(map[int64]*Spammer),
		globalCfg:  make(map[string]interface{}),
	}
}

// SetGlobalCfg sets a global configuration value that can be accessed by all spammers.
// This allows sharing configuration data across multiple spammer instances.
func (d *Daemon) SetGlobalCfg(name string, value interface{}) {
	d.globalCfg[name] = value
}

// GetGlobalCfg returns the global configuration map that contains shared configuration
// values accessible to all spammer instances.
func (d *Daemon) GetGlobalCfg() map[string]interface{} {
	return d.globalCfg
}

// GetClientPool returns the RPC client pool used by all spammers.
// This provides access to the Ethereum client connections for transaction submission.
func (d *Daemon) GetClientPool() *spamoor.ClientPool {
	return d.clientPool
}

// Run initializes the daemon by restoring spammers from the database.
// Returns true if this is the first launch, false if spammers were restored.
// Marks the first launch state in the database to track initialization status.
func (d *Daemon) Run() (bool, error) {
	// check if this is the first launch
	var notFirstLaunch bool
	d.db.GetSpamoorState("first_launch", &notFirstLaunch)

	// restore all spammer from db
	spammerList, err := d.db.GetSpammers()
	if err != nil {
		return false, fmt.Errorf("failed to get all spammer: %w", err)
	}

	for _, spammer := range spammerList {
		_, err := d.restoreSpammer(spammer)
		if err != nil {
			return false, fmt.Errorf("failed to restore spammer: %w", err)
		}
	}

	// Mark that we've launched
	if !notFirstLaunch {
		err = d.db.SetSpamoorState(nil, "first_launch", true)
		if err != nil {
			return false, fmt.Errorf("failed to mark first launch: %w", err)
		}
	}

	return !notFirstLaunch, nil
}

// GetSpammer retrieves a spammer instance by ID from the internal map.
// Returns nil if the spammer with the given ID does not exist.
// This method is thread-safe using read locks.
func (d *Daemon) GetSpammer(id int64) *Spammer {
	d.spammerMapMtx.RLock()
	defer d.spammerMapMtx.RUnlock()

	return d.spammerMap[id]
}

// GetAllSpammers returns all spammer instances sorted by ID in descending order.
// This provides a snapshot of all active spammer instances managed by the daemon.
// The returned slice is safe to modify as it's a copy.
func (d *Daemon) GetAllSpammers() []*Spammer {
	d.spammerMapMtx.RLock()
	defer d.spammerMapMtx.RUnlock()

	spammerList := make([]*Spammer, 0, len(d.spammerMap))
	for _, spammer := range d.spammerMap {
		spammerList = append(spammerList, spammer)
	}

	sort.Slice(spammerList, func(i, j int) bool {
		return spammerList[i].GetID() > spammerList[j].GetID()
	})

	return spammerList
}

// DeleteSpammer removes a spammer from both the daemon and database.
// If the spammer is running, it will be paused first before deletion.
// Returns an error if the spammer is not found or if database deletion fails.
func (d *Daemon) DeleteSpammer(id int64) error {
	d.spammerMapMtx.Lock()
	defer d.spammerMapMtx.Unlock()

	spammer := d.spammerMap[id]
	if spammer == nil {
		return fmt.Errorf("spammer not found")
	}

	// Stop if running
	if spammer.scenarioCancel != nil {
		spammer.Pause()
	}

	// Delete from DB
	err := d.db.RunDBTransaction(func(tx *sqlx.Tx) error {
		return d.db.DeleteSpammer(tx, id)
	})
	if err != nil {
		return fmt.Errorf("failed to delete spammer: %w", err)
	}

	// Remove from map
	delete(d.spammerMap, id)
	return nil
}

// UpdateSpammer modifies the name, description, and configuration of an existing spammer.
// The configuration is validated by attempting to unmarshal it into SpammerConfig.
// Returns an error if the spammer is not found, config is invalid, or database update fails.
func (d *Daemon) UpdateSpammer(id int64, name string, description string, config string) error {
	d.spammerMapMtx.Lock()
	defer d.spammerMapMtx.Unlock()

	spammer := d.spammerMap[id]
	if spammer == nil {
		return fmt.Errorf("spammer not found")
	}

	// Validate config
	if err := yaml.Unmarshal([]byte(config), &SpammerConfig{}); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	// Update DB
	spammer.dbEntity.Name = name
	spammer.dbEntity.Description = description
	spammer.dbEntity.Config = config

	err := d.db.RunDBTransaction(func(tx *sqlx.Tx) error {
		return d.db.UpdateSpammer(tx, spammer.dbEntity)
	})
	if err != nil {
		return fmt.Errorf("failed to update spammer: %w", err)
	}

	return nil
}

// ReclaimSpammer reclaims funds from all wallets in the spammer's wallet pool.
// This transfers remaining ETH back to the root wallet for reuse.
// Returns an error if the spammer is not found or if fund reclamation fails.
func (d *Daemon) ReclaimSpammer(id int64) error {
	spammer := d.GetSpammer(id)
	if spammer == nil {
		return fmt.Errorf("spammer not found")
	}

	if spammer.walletPool == nil {
		return nil
	}

	return spammer.walletPool.ReclaimFunds(d.ctx, nil)
}

// GetRootWallet returns the root wallet used for funding spammer wallets.
// This provides access to the main wallet that distributes funds to child wallets.
func (d *Daemon) GetRootWallet() *spamoor.RootWallet {
	return d.rootWallet
}

// TrackTransactionSent records a successful transaction send for metrics
func (d *Daemon) TrackTransactionSent(spammerID int64) {
	if d.metricsCollector == nil {
		return
	}

	spammer := d.GetSpammer(spammerID)
	if spammer == nil {
		return
	}

	d.metricsCollector.IncrementTransactionsSent(
		spammer.GetID(),
		spammer.GetName(),
		spammer.GetScenario(),
	)
}

// TrackTransactionFailure records a failed transaction for metrics
func (d *Daemon) TrackTransactionFailure(spammerID int64) {
	if d.metricsCollector == nil {
		return
	}

	spammer := d.GetSpammer(spammerID)
	if spammer == nil {
		return
	}

	d.metricsCollector.IncrementTransactionFailures(
		spammer.GetID(),
		spammer.GetName(),
		spammer.GetScenario(),
	)
}

// TrackSpammerStatusChange updates the spammer running status metric
func (d *Daemon) TrackSpammerStatusChange(spammerID int64, running bool) {
	if d.metricsCollector == nil {
		return
	}

	spammer := d.GetSpammer(spammerID)
	if spammer == nil {
		return
	}

	d.metricsCollector.SetSpammerRunning(
		spammer.GetID(),
		spammer.GetName(),
		spammer.GetScenario(),
		running,
	)
}

// GetShortWindowMetrics returns the 30-minute per-block metrics for the dashboard
func (d *Daemon) GetShortWindowMetrics() *MultiGranularityMetrics {
	if d.txPoolMetricsCollector == nil {
		return nil
	}
	return d.txPoolMetricsCollector.GetShortWindowMetrics()
}

// GetLongWindowMetrics returns the 6-hour per-32-block metrics for the dashboard
func (d *Daemon) GetLongWindowMetrics() *MultiGranularityMetrics {
	if d.txPoolMetricsCollector == nil {
		return nil
	}
	return d.txPoolMetricsCollector.GetLongWindowMetrics()
}

// GetMetricsCollector returns the TxPool metrics collector for real-time subscriptions
func (d *Daemon) GetMetricsCollector() *TxPoolMetricsCollector {
	return d.txPoolMetricsCollector
}

// No longer need RegisterSpammerForMetrics/UnregisterSpammerFromMetrics
// Spammer metrics are automatically tracked via WalletPool.GetSpammerID()

// GetSpammerName returns the name of a spammer by ID for metrics dashboard
func (d *Daemon) GetSpammerName(spammerID uint64) string {
	spammer := d.GetSpammer(int64(spammerID))
	if spammer == nil {
		return "unknown"
	}
	return spammer.GetName()
}

// Shutdown performs a graceful shutdown of the daemon and all running spammers.
// It cancels the context to stop all spammers, waits up to 10 seconds for them to finish,
// and closes the database connection. This ensures clean resource cleanup.
func (d *Daemon) Shutdown() {
	d.logger.Info("initiating graceful shutdown")

	// Cancel context to stop all spammers
	d.cancel()

	// Get all running spammers
	d.spammerMapMtx.RLock()
	spammers := make([]*Spammer, 0, len(d.spammerMap))
	for _, s := range d.spammerMap {
		if s.running {
			spammers = append(spammers, s)
		}
	}
	d.spammerMapMtx.RUnlock()

	// Wait for all spammers to finish with a timeout
	if len(spammers) > 0 {
		d.logger.Infof("waiting for %d spammers to stop", len(spammers))

		done := make(chan struct{})
		go func() {
			d.spammerWg.Wait()
			close(done)
		}()

		select {
		case <-done:
			d.logger.Info("all spammers stopped successfully")
		case <-time.After(10 * time.Second):
			d.logger.Warn("timeout waiting for spammers to stop")
		}
	}

	// Shutdown TxPool metrics collector
	if d.txPoolMetricsCollector != nil {
		d.txPoolMetricsCollector.Shutdown()
		d.logger.Info("TxPool metrics collector shutdown")
	}

	// Close database connection
	if err := d.db.Close(); err != nil {
		d.logger.Errorf("error closing database: %v", err)
	} else {
		d.logger.Info("database closed successfully")
	}

	d.logger.Info("shutdown complete")
}
