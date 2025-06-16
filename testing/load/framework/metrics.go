package framework

import (
	"fmt"
	"runtime"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// MetricsCollector collects and analyzes performance metrics during load tests
type MetricsCollector struct {
	startTime time.Time
	endTime   time.Time
	
	// Transaction tracking
	transactionTimes    map[common.Hash]time.Time
	completionTimes     map[common.Hash]time.Time
	transactionErrors   map[common.Hash]error
	
	// Real-time metrics
	tpsHistory          []float64
	latencyHistory      []time.Duration
	
	// Error tracking
	errorsByType        map[string]int
	criticalErrors      []string
	
	// Validation results
	nonceGaps           int
	orderingViolations  int
	validationErrors    []string
	
	// Resource metrics
	memoryUsage         []float64
	cpuUsage           []float64
	goroutineCount     []int
	
	// Synchronization
	mutex sync.RWMutex
	
	// Configuration
	metricsInterval time.Duration
	stopChan       chan struct{}
	stopped        bool
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		transactionTimes:  make(map[common.Hash]time.Time),
		completionTimes:   make(map[common.Hash]time.Time),
		transactionErrors: make(map[common.Hash]error),
		errorsByType:      make(map[string]int),
		criticalErrors:    make([]string, 0),
		validationErrors:  make([]string, 0),
		metricsInterval:   time.Second,
		stopChan:         make(chan struct{}),
	}
}

// Start begins metrics collection
func (m *MetricsCollector) Start() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	if m.stopped {
		return
	}
	
	m.startTime = time.Now()
	
	// Start resource monitoring goroutine
	go m.monitorResources()
}

// Stop ends metrics collection
func (m *MetricsCollector) Stop() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	if m.stopped {
		return
	}
	
	m.endTime = time.Now()
	m.stopped = true
	close(m.stopChan)
}

// RecordTransactionSent records when a transaction was sent
func (m *MetricsCollector) RecordTransactionSent(txHash common.Hash) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.transactionTimes[txHash] = time.Now()
}

// RecordTransactionCompleted records when a transaction was completed
func (m *MetricsCollector) RecordTransactionCompleted(txHash common.Hash) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	completionTime := time.Now()
	m.completionTimes[txHash] = completionTime
	
	// Calculate latency if we have the start time
	if startTime, exists := m.transactionTimes[txHash]; exists {
		latency := completionTime.Sub(startTime)
		m.latencyHistory = append(m.latencyHistory, latency)
	}
}

// RecordTransactionError records a transaction error
func (m *MetricsCollector) RecordTransactionError(txHash common.Hash, err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.transactionErrors[txHash] = err
	
	// Categorize error
	errorType := "unknown"
	if err != nil {
		errorType = err.Error()
	}
	m.errorsByType[errorType]++
}

// RecordCriticalError records a critical error that affects test validity
func (m *MetricsCollector) RecordCriticalError(errorMsg string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.criticalErrors = append(m.criticalErrors, errorMsg)
}

// RecordNonceGap records a nonce gap detection
func (m *MetricsCollector) RecordNonceGap(address common.Address, expectedNonce, actualNonce uint64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.nonceGaps++
	m.validationErrors = append(m.validationErrors, 
		fmt.Sprintf("Nonce gap for %s: expected %d, got %d", address.Hex(), expectedNonce, actualNonce))
}

// RecordOrderingViolation records a transaction ordering violation
func (m *MetricsCollector) RecordOrderingViolation(address common.Address, details string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.orderingViolations++
	m.validationErrors = append(m.validationErrors, 
		fmt.Sprintf("Ordering violation for %s: %s", address.Hex(), details))
}

// RecordValidationError records a general validation error
func (m *MetricsCollector) RecordValidationError(errorMsg string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.validationErrors = append(m.validationErrors, errorMsg)
}

// monitorResources monitors system resource usage
func (m *MetricsCollector) monitorResources() {
	ticker := time.NewTicker(m.metricsInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			m.collectResourceMetrics()
		case <-m.stopChan:
			return
		}
	}
}

// collectResourceMetrics collects current resource usage
func (m *MetricsCollector) collectResourceMetrics() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	// Memory usage
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	memoryMB := float64(memStats.Alloc) / 1024 / 1024
	m.memoryUsage = append(m.memoryUsage, memoryMB)
	
	// Goroutine count
	goroutines := runtime.NumGoroutine()
	m.goroutineCount = append(m.goroutineCount, goroutines)
	
	// Calculate current TPS
	if len(m.completionTimes) > 0 {
		currentTime := time.Now()
		recentCompletions := 0
		cutoff := currentTime.Add(-m.metricsInterval)
		
		for _, completionTime := range m.completionTimes {
			if completionTime.After(cutoff) {
				recentCompletions++
			}
		}
		
		tps := float64(recentCompletions) / m.metricsInterval.Seconds()
		m.tpsHistory = append(m.tpsHistory, tps)
	}
}

// GenerateReport generates a comprehensive load test report
func (m *MetricsCollector) GenerateReport(config *LoadTestConfig) *LoadTestResult {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	result := &LoadTestResult{
		Config:    config,
		StartTime: m.startTime,
		EndTime:   m.endTime,
		Duration:  m.endTime.Sub(m.startTime),
		
		NonceGaps:          m.nonceGaps,
		OrderingViolations: m.orderingViolations,
		ValidationErrors:   make([]string, len(m.validationErrors)),
		ErrorsByType:       make(map[string]int),
		CriticalErrors:     make([]string, len(m.criticalErrors)),
	}
	
	// Copy validation errors and critical errors
	copy(result.ValidationErrors, m.validationErrors)
	copy(result.CriticalErrors, m.criticalErrors)
	
	// Copy errors by type
	for k, v := range m.errorsByType {
		result.ErrorsByType[k] = v
	}
	
	// Calculate transaction metrics
	result.TotalTransactions = int64(len(m.transactionTimes))
	result.SuccessfulTxs = int64(len(m.completionTimes))
	result.FailedTxs = int64(len(m.transactionErrors))
	
	// Calculate TPS metrics
	if result.Duration > 0 {
		result.AverageTPS = float64(result.SuccessfulTxs) / result.Duration.Seconds()
	}
	
	if len(m.tpsHistory) > 0 {
		result.PeakTPS = maxFloat64(m.tpsHistory)
	}
	
	// Calculate latency metrics
	if len(m.latencyHistory) > 0 {
		result.AverageLatency = m.calculateAverageLatency()
		result.P50Latency = m.calculatePercentileLatency(50)
		result.P95Latency = m.calculatePercentileLatency(95)
		result.P99Latency = m.calculatePercentileLatency(99)
		result.MaxLatency = maxDuration(m.latencyHistory)
	}
	
	// Calculate resource metrics
	if len(m.memoryUsage) > 0 {
		result.MemoryUsageMB = averageFloat64(m.memoryUsage)
	}
	
	if len(m.goroutineCount) > 0 {
		result.GoroutineCount = int(averageFloat64(convertIntSliceToFloat64(m.goroutineCount)))
	}
	
	// Calculate error rate
	if result.TotalTransactions > 0 {
		result.ErrorRate = float64(result.FailedTxs) / float64(result.TotalTransactions)
	}
	
	return result
}

// calculateAverageLatency calculates the average latency
func (m *MetricsCollector) calculateAverageLatency() time.Duration {
	if len(m.latencyHistory) == 0 {
		return 0
	}
	
	var total time.Duration
	for _, latency := range m.latencyHistory {
		total += latency
	}
	
	return total / time.Duration(len(m.latencyHistory))
}

// calculatePercentileLatency calculates the nth percentile latency
func (m *MetricsCollector) calculatePercentileLatency(percentile float64) time.Duration {
	if len(m.latencyHistory) == 0 {
		return 0
	}
	
	// Create a copy and sort it
	sorted := make([]time.Duration, len(m.latencyHistory))
	copy(sorted, m.latencyHistory)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})
	
	// Calculate percentile index
	index := int(percentile/100.0*float64(len(sorted))) - 1
	if index < 0 {
		index = 0
	}
	if index >= len(sorted) {
		index = len(sorted) - 1
	}
	
	return sorted[index]
}

// Helper functions
func maxFloat64(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	
	max := values[0]
	for _, v := range values[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

func maxDuration(values []time.Duration) time.Duration {
	if len(values) == 0 {
		return 0
	}
	
	max := values[0]
	for _, v := range values[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

func averageFloat64(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	
	var sum float64
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func convertIntSliceToFloat64(values []int) []float64 {
	result := make([]float64, len(values))
	for i, v := range values {
		result[i] = float64(v)
	}
	return result
}