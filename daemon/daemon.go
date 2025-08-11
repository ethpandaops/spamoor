package daemon

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"
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

	// Audit logger for tracking actions
	auditLogger *AuditLogger
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

	// load and apply client configs from database
	err := d.loadAndApplyClientConfigs()
	if err != nil {
		d.logger.Warnf("failed to load client configs: %v", err)
	}

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
func (d *Daemon) DeleteSpammer(id int64, userEmail string) error {
	d.spammerMapMtx.Lock()
	defer d.spammerMapMtx.Unlock()

	spammer := d.spammerMap[id]
	if spammer == nil {
		return fmt.Errorf("spammer not found")
	}

	// Capture name for audit log
	spammerName := spammer.GetName()

	// Stop if running
	if spammer.scenarioCancel != nil {
		spammer.Pause()
	}

	// Delete from DB
	err := d.db.RunDBTransaction(func(tx *sqlx.Tx) error {
		if err := d.db.DeleteSpammer(tx, id); err != nil {
			return err
		}

		// Audit log the deletion
		if d.auditLogger != nil {
			return d.auditLogger.LogSpammerDelete(tx, userEmail, id, spammerName)
		}

		return nil
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
func (d *Daemon) UpdateSpammer(id int64, name string, description string, config string, userEmail string) error {
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

	// Capture old values for audit logging
	oldName := spammer.dbEntity.Name
	oldDescription := spammer.dbEntity.Description
	oldConfig := spammer.dbEntity.Config

	// Update DB
	spammer.dbEntity.Name = name
	spammer.dbEntity.Description = description
	spammer.dbEntity.Config = config

	err := d.db.RunDBTransaction(func(tx *sqlx.Tx) error {
		if err := d.db.UpdateSpammer(tx, spammer.dbEntity); err != nil {
			return err
		}

		// Audit log the update
		if d.auditLogger != nil {
			return d.auditLogger.LogSpammerUpdate(tx, userEmail, id, oldName, name, oldDescription, description, oldConfig, config)
		}

		return nil
	})
	if err != nil {
		// Revert changes on error
		spammer.dbEntity.Name = oldName
		spammer.dbEntity.Description = oldDescription
		spammer.dbEntity.Config = oldConfig
		return fmt.Errorf("failed to update spammer: %w", err)
	}

	return nil
}

// StartSpammer starts a spammer and logs the action
func (d *Daemon) StartSpammer(id int64, userEmail string) error {
	spammer := d.GetSpammer(id)
	if spammer == nil {
		return fmt.Errorf("spammer not found")
	}

	err := spammer.Start()
	if err != nil {
		return err
	}

	// Audit log the start action
	if d.auditLogger != nil && userEmail != "" {
		return d.auditLogger.LogSpammerAction(userEmail, db.AuditActionSpammerStart, id, spammer.GetName(), nil)
	}

	return nil
}

// PauseSpammer pauses a spammer and logs the action
func (d *Daemon) PauseSpammer(id int64, userEmail string) error {
	spammer := d.GetSpammer(id)
	if spammer == nil {
		return fmt.Errorf("spammer not found")
	}

	err := spammer.Pause()
	if err != nil {
		return err
	}

	// Audit log the pause action
	if d.auditLogger != nil && userEmail != "" {
		return d.auditLogger.LogSpammerAction(userEmail, db.AuditActionSpammerPause, id, spammer.GetName(), nil)
	}

	return nil
}

// ReclaimSpammer reclaims funds from all wallets in the spammer's wallet pool.
// This transfers remaining ETH back to the root wallet for reuse.
// Returns an error if the spammer is not found or if fund reclamation fails.
func (d *Daemon) ReclaimSpammer(id int64, userEmail string) error {
	spammer := d.GetSpammer(id)
	if spammer == nil {
		return fmt.Errorf("spammer not found")
	}

	if spammer.walletPool == nil {
		return nil
	}

	err := spammer.walletPool.ReclaimFunds(d.ctx, nil)
	if err != nil {
		return err
	}

	// Audit log the reclaim action
	if d.auditLogger != nil && userEmail != "" {
		return d.auditLogger.LogSpammerAction(userEmail, db.AuditActionSpammerReclaim, id, spammer.GetName(), nil)
	}

	return nil
}

// GetRootWallet returns the root wallet used for funding spammer wallets.
// This provides access to the main wallet that distributes funds to child wallets.
func (d *Daemon) GetRootWallet() *spamoor.RootWallet {
	return d.rootWallet
}

// SetAuditLogger sets the audit logger for the daemon
func (d *Daemon) SetAuditLogger(logger *AuditLogger) {
	d.auditLogger = logger
}

// GetAuditLogger returns the audit logger
func (d *Daemon) GetAuditLogger() *AuditLogger {
	return d.auditLogger
}

// GetDatabase returns the database instance
func (d *Daemon) GetDatabase() *db.Database {
	return d.db
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

// loadAndApplyClientConfigs retrieves client configurations from the database
// and applies them to the corresponding clients in the client pool.
// This merges database settings with flag-provided settings according to the rules:
// - Name and enabled state from DB take precedence over flags
// - Tags are combined (flags + database tags)
func (d *Daemon) loadAndApplyClientConfigs() error {
	// Get all client configs from database
	configs, err := d.db.GetClientConfigs()
	if err != nil {
		return fmt.Errorf("failed to get client configs: %w", err)
	}

	// Create a map for quick lookup
	configMap := make(map[string]*db.ClientConfig)
	for _, config := range configs {
		configMap[config.RpcUrl] = config
	}

	// Apply configs to clients in the pool
	allClients := d.clientPool.GetAllClients()
	for _, client := range allClients {
		rpcUrl := client.GetRPCHost()
		config, exists := configMap[rpcUrl]

		if exists {
			// Apply enabled state from database
			client.SetEnabled(config.Enabled)

			// Apply name override from database
			if config.Name != "" {
				client.SetNameOverride(config.Name)
			}

			// Merge tags (combine flag groups with database tags)
			currentGroups := client.GetClientGroups()
			dbTags := config.GetTagsAsSlice()

			// Create combined groups (current groups + database tags)
			combinedGroups := make([]string, 0)
			groupSet := make(map[string]bool)

			// Add current groups first
			for _, group := range currentGroups {
				if !groupSet[group] {
					combinedGroups = append(combinedGroups, group)
					groupSet[group] = true
				}
			}

			// Add database tags
			for _, tag := range dbTags {
				if !groupSet[tag] {
					combinedGroups = append(combinedGroups, tag)
					groupSet[tag] = true
				}
			}

			// Set combined groups if we have any
			if len(combinedGroups) > 0 {
				client.SetClientGroups(combinedGroups)
			}

			// Apply client type override from database
			if config.ClientType != "" {
				client.SetClientTypeOverride(config.ClientType)
			}

			d.logger.Debugf("Applied config for client %s: enabled=%v, name=%s, type=%s, groups=%v",
				rpcUrl, config.Enabled, config.Name, config.ClientType, client.GetClientGroups())
		}
	}

	return nil
}

// UpdateClientConfig updates the configuration for a specific client
// and persists the changes to the database.
func (d *Daemon) UpdateClientConfig(rpcUrl, name, tags, clientType string, enabled bool, userEmail string) error {
	// Find the client in the pool
	allClients := d.clientPool.GetAllClients()
	var targetClient *spamoor.Client
	for _, client := range allClients {
		if client.GetRPCHost() == rpcUrl {
			targetClient = client
			break
		}
	}

	if targetClient == nil {
		return fmt.Errorf("client with RPC URL %s not found", rpcUrl)
	}

	// Update the client configuration
	targetClient.SetEnabled(enabled)

	// Apply name override
	targetClient.SetNameOverride(name)

	// Parse tags and set as client groups
	tagSlice := strings.Split(tags, ",")
	for i, tag := range tagSlice {
		tagSlice[i] = strings.TrimSpace(tag)
	}

	// Remove empty tags
	filteredTags := make([]string, 0)
	for _, tag := range tagSlice {
		if tag != "" {
			filteredTags = append(filteredTags, tag)
		}
	}

	if len(filteredTags) > 0 {
		targetClient.SetClientGroups(filteredTags)
	}

	// Apply client type override
	if clientType != "" {
		targetClient.SetClientTypeOverride(clientType)
	}

	// Persist to database
	config := &db.ClientConfig{
		RpcUrl:     rpcUrl,
		Name:       name,
		Tags:       tags,
		ClientType: clientType,
		Enabled:    enabled,
	}

	// Get existing config for audit logging
	existingConfig, _ := d.db.GetClientConfig(rpcUrl)

	err := d.db.RunDBTransaction(func(tx *sqlx.Tx) error {
		if err := d.db.UpsertClientConfig(tx, config); err != nil {
			return err
		}

		// Audit log the client update
		if d.auditLogger != nil && userEmail != "" {
			changes := make(map[string]interface{})

			if existingConfig != nil {
				if existingConfig.Name != name {
					changes["name"] = map[string]interface{}{"old": existingConfig.Name, "new": name}
				}
				if existingConfig.Tags != tags {
					changes["tags"] = map[string]interface{}{"old": existingConfig.Tags, "new": tags}
				}
				if existingConfig.ClientType != clientType {
					changes["client_type"] = map[string]interface{}{"old": existingConfig.ClientType, "new": clientType}
				}
				if existingConfig.Enabled != enabled {
					changes["enabled"] = map[string]interface{}{"old": existingConfig.Enabled, "new": enabled}
				}
			} else {
				// New client config - only log fields that differ from defaults
				if name != "" {
					changes["name"] = map[string]interface{}{"old": "", "new": name}
				}
				if tags != "" {
					changes["tags"] = map[string]interface{}{"old": "", "new": tags}
				}
				if clientType != "" {
					changes["client_type"] = map[string]interface{}{"old": "", "new": clientType}
				}
				if enabled != true {
					changes["enabled"] = map[string]interface{}{"old": true, "new": enabled}
				}
			}

			if len(changes) > 0 {
				return d.auditLogger.LogClientUpdate(tx, userEmail, rpcUrl, targetClient.GetName(), changes)
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to update client config in database: %w", err)
	}

	d.logger.Infof("Updated client config: %s (enabled=%v, name=%s, type=%s, tags=%s)",
		rpcUrl, enabled, name, clientType, tags)

	return nil
}

// GetClientConfig retrieves the configuration for a specific client from the database.
func (d *Daemon) GetClientConfig(rpcUrl string) (*db.ClientConfig, error) {
	config, err := d.db.GetClientConfig(rpcUrl)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return default config if not found in database
			return &db.ClientConfig{
				RpcUrl:     rpcUrl,
				Name:       "",
				Tags:       "",
				ClientType: "",
				Enabled:    true,
			}, nil
		}
		return nil, fmt.Errorf("failed to get client config: %w", err)
	}
	return config, nil
}
