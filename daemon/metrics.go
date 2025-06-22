package daemon

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

// MetricsCollector holds all Prometheus metrics for the spamoor daemon
type MetricsCollector struct {
	// Gauge to track if a spammer is running (1) or not (0)
	spammerRunning *prometheus.GaugeVec

	// Counter for transactions sent per spammer
	transactionsSent *prometheus.CounterVec

	// Counter for transaction failures per spammer
	transactionFailures *prometheus.CounterVec

	// New metrics for txpool analytics
	pendingTransactionsGauge *prometheus.GaugeVec
	blockGasUsageGauge       *prometheus.GaugeVec
	totalBlockGasGauge       *prometheus.GaugeVec
}

// NewMetricsCollector creates and registers all metrics
func (d *Daemon) NewMetricsCollector() *MetricsCollector {
	mc := &MetricsCollector{
		spammerRunning: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "spamoor_spammer_running",
				Help: "Whether a spammer is currently running (1) or not (0)",
			},
			[]string{"spammer_id", "spammer_name", "scenario"},
		),

		transactionsSent: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "spamoor_transactions_sent_total",
				Help: "Total number of transactions sent by each spammer",
			},
			[]string{"spammer_id", "spammer_name", "scenario"},
		),

		transactionFailures: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "spamoor_transaction_failures_total",
				Help: "Total number of transaction failures by each spammer",
			},
			[]string{"spammer_id", "spammer_name", "scenario"},
		),

		pendingTransactionsGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "spamoor_pending_transactions",
				Help: "Number of pending transactions per spammer",
			},
			[]string{"spammer_id", "spammer_name"},
		),

		blockGasUsageGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "spamoor_block_gas_usage",
				Help: "Gas used per spammer in the latest block",
			},
			[]string{"spammer_id", "spammer_name", "block_number"},
		),

		totalBlockGasGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "spamoor_total_block_gas",
				Help: "Total gas used in block with spammer/others breakdown",
			},
			[]string{"category"}, // "spammer" or "others"
		),
	}

	// Register metrics with default registry
	prometheus.MustRegister(mc.spammerRunning)
	prometheus.MustRegister(mc.transactionsSent)
	prometheus.MustRegister(mc.transactionFailures)
	prometheus.MustRegister(mc.pendingTransactionsGauge)
	prometheus.MustRegister(mc.blockGasUsageGauge)
	prometheus.MustRegister(mc.totalBlockGasGauge)

	return mc
}

// SetSpammerRunning sets the running status for a spammer
func (mc *MetricsCollector) SetSpammerRunning(spammerID int64, spammerName, scenario string, running bool) {
	labels := prometheus.Labels{
		"spammer_id":   strconv.FormatInt(spammerID, 10),
		"spammer_name": spammerName,
		"scenario":     scenario,
	}

	if running {
		mc.spammerRunning.With(labels).Set(1)
	} else {
		mc.spammerRunning.With(labels).Set(0)
	}
}

// IncrementTransactionsSent increments the sent transaction counter
func (mc *MetricsCollector) IncrementTransactionsSent(spammerID int64, spammerName, scenario string) {
	labels := prometheus.Labels{
		"spammer_id":   strconv.FormatInt(spammerID, 10),
		"spammer_name": spammerName,
		"scenario":     scenario,
	}
	mc.transactionsSent.With(labels).Inc()
}

// IncrementTransactionFailures increments the failed transaction counter
func (mc *MetricsCollector) IncrementTransactionFailures(spammerID int64, spammerName, scenario string) {
	labels := prometheus.Labels{
		"spammer_id":   strconv.FormatInt(spammerID, 10),
		"spammer_name": spammerName,
		"scenario":     scenario,
	}
	mc.transactionFailures.With(labels).Inc()
}

// UpdatePendingTransactions updates the pending transactions gauge for a spammer
func (mc *MetricsCollector) UpdatePendingTransactions(spammerID int64, spammerName string, count uint64) {
	labels := prometheus.Labels{
		"spammer_id":   strconv.FormatInt(spammerID, 10),
		"spammer_name": spammerName,
	}
	mc.pendingTransactionsGauge.With(labels).Set(float64(count))
}

// UpdateBlockGasUsage updates the block gas usage gauge for a spammer
func (mc *MetricsCollector) UpdateBlockGasUsage(spammerID int64, spammerName string, blockNumber uint64, gasUsed uint64) {
	labels := prometheus.Labels{
		"spammer_id":   strconv.FormatInt(spammerID, 10),
		"spammer_name": spammerName,
		"block_number": strconv.FormatUint(blockNumber, 10),
	}
	mc.blockGasUsageGauge.With(labels).Set(float64(gasUsed))
}

// UpdateTotalBlockGas updates the total block gas gauge with category breakdown
func (mc *MetricsCollector) UpdateTotalBlockGas(category string, gasUsed uint64) {
	labels := prometheus.Labels{
		"category": category,
	}
	mc.totalBlockGasGauge.With(labels).Set(float64(gasUsed))
}

func (d *Daemon) InitializeTxPoolMetrics() error {
	// Initialize TxPool metrics collector with multi-granularity windows
	d.txPoolMetricsCollector = NewTxPoolMetricsCollector(d.txpool)

	return nil
}

func (d *Daemon) InitializeMetrics() error {
	// Initialize metrics collector (Prometheus metrics)
	d.metricsCollector = d.NewMetricsCollector()

	// Initialize metrics for existing spammers
	d.initializeExistingSpammerMetrics()

	return nil
}

// initializeExistingSpammerMetrics sets initial metric values for all existing spammers
func (d *Daemon) initializeExistingSpammerMetrics() {
	spammers := d.GetAllSpammers()
	for _, spammer := range spammers {
		// Set initial running status based on current spammer status
		running := spammer.GetStatus() == int(SpammerStatusRunning)
		d.metricsCollector.SetSpammerRunning(
			spammer.GetID(),
			spammer.GetName(),
			spammer.GetScenario(),
			running,
		)
	}
}
