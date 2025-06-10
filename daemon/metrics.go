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
	}

	// Register metrics with default registry
	prometheus.MustRegister(mc.spammerRunning)
	prometheus.MustRegister(mc.transactionsSent)
	prometheus.MustRegister(mc.transactionFailures)

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

func (d *Daemon) InitializeMetrics() error {
	// Initialize metrics collector
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
