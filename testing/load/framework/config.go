package framework

import (
	"time"
)

// LoadTestConfig defines the configuration for load tests
type LoadTestConfig struct {
	// Test execution parameters
	Duration        time.Duration `yaml:"duration" json:"duration"`
	TargetTPS       int           `yaml:"target_tps" json:"target_tps"`
	ClientCount     int           `yaml:"client_count" json:"client_count"`
	WalletCount     int           `yaml:"wallet_count" json:"wallet_count"`
	TransactionType string        `yaml:"transaction_type" json:"transaction_type"`
	ValidationMode  string        `yaml:"validation_mode" json:"validation_mode"`
	
	// Transaction parameters
	GasLimit        uint64 `yaml:"gas_limit" json:"gas_limit"`
	BaseFee         uint64 `yaml:"base_fee" json:"base_fee"`
	TipFee          uint64 `yaml:"tip_fee" json:"tip_fee"`
	Amount          uint64 `yaml:"amount" json:"amount"`
	
	// Test behavior
	MaxPending      uint64        `yaml:"max_pending" json:"max_pending"`
	Timeout         time.Duration `yaml:"timeout" json:"timeout"`
	WarmupDuration  time.Duration `yaml:"warmup_duration" json:"warmup_duration"`
	
	// Validation settings
	EnableNonceValidation     bool `yaml:"enable_nonce_validation" json:"enable_nonce_validation"`
	EnableOrderingValidation  bool `yaml:"enable_ordering_validation" json:"enable_ordering_validation"`
	EnableThroughputValidation bool `yaml:"enable_throughput_validation" json:"enable_throughput_validation"`
}

// DefaultLoadTestConfig returns a default configuration for load tests
func DefaultLoadTestConfig() *LoadTestConfig {
	return &LoadTestConfig{
		Duration:        60 * time.Second,
		TargetTPS:       100,
		ClientCount:     1,
		WalletCount:     10,
		TransactionType: "dynamicfee",
		ValidationMode:  "full",
		
		GasLimit: 21000,
		BaseFee:  20,
		TipFee:   2,
		Amount:   1000000000000000, // 0.001 ETH in wei
		
		MaxPending:     1000,
		Timeout:        300 * time.Second,
		WarmupDuration: 10 * time.Second,
		
		EnableNonceValidation:      true,
		EnableOrderingValidation:   true,
		EnableThroughputValidation: true,
	}
}

// LoadTestResult contains the results of a load test execution
type LoadTestResult struct {
	// Test metadata
	Config    *LoadTestConfig `json:"config"`
	StartTime time.Time       `json:"start_time"`
	EndTime   time.Time       `json:"end_time"`
	Duration  time.Duration   `json:"duration"`
	
	// Transaction metrics
	TotalTransactions   int64   `json:"total_transactions"`
	SuccessfulTxs       int64   `json:"successful_txs"`
	FailedTxs          int64   `json:"failed_txs"`
	AverageTPS         float64 `json:"average_tps"`
	PeakTPS            float64 `json:"peak_tps"`
	
	// Latency metrics
	AverageLatency  time.Duration `json:"average_latency"`
	P50Latency      time.Duration `json:"p50_latency"`
	P95Latency      time.Duration `json:"p95_latency"`
	P99Latency      time.Duration `json:"p99_latency"`
	MaxLatency      time.Duration `json:"max_latency"`
	
	// Validation results
	NonceGaps           int    `json:"nonce_gaps"`
	OrderingViolations  int    `json:"ordering_violations"`
	ValidationErrors    []string `json:"validation_errors"`
	
	// Resource metrics
	MemoryUsageMB      float64 `json:"memory_usage_mb"`
	CPUUsagePercent    float64 `json:"cpu_usage_percent"`
	GoroutineCount     int     `json:"goroutine_count"`
	
	// Error summary
	ErrorRate      float64           `json:"error_rate"`
	ErrorsByType   map[string]int    `json:"errors_by_type"`
	CriticalErrors []string          `json:"critical_errors"`
}

// IsSuccessful returns true if the load test passed all validation criteria
func (r *LoadTestResult) IsSuccessful() bool {
	// Check for critical failures
	if len(r.CriticalErrors) > 0 {
		return false
	}
	
	// Check nonce validation
	if r.Config.EnableNonceValidation && r.NonceGaps > 0 {
		return false
	}
	
	// Check ordering validation
	if r.Config.EnableOrderingValidation && r.OrderingViolations > 0 {
		return false
	}
	
	// Check throughput validation (within 10% of target)
	if r.Config.EnableThroughputValidation {
		targetTPS := float64(r.Config.TargetTPS)
		if r.AverageTPS < targetTPS*0.9 {
			return false
		}
	}
	
	// Check error rate (should be less than 1%)
	if r.ErrorRate > 0.01 {
		return false
	}
	
	return true
}