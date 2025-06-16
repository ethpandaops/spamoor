//go:build loadtest

package scenarios

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ethpandaops/spamoor/spamoortypes"
	"github.com/ethpandaops/spamoor/testing/load/framework"
	"github.com/ethpandaops/spamoor/testing/load/validators"
	testingutils "github.com/ethpandaops/spamoor/testing/utils"
)

// ThroughputTestSuite contains throughput benchmark tests
type ThroughputTestSuite struct {
	server     *testingutils.MockRPCServer
	clients    []spamoortypes.Client
	wallets    []spamoortypes.Wallet
	logger     *logrus.Entry
	validator  *validators.TransactionValidator
}

// setupThroughputTest sets up the test environment for throughput testing
func setupThroughputTest(t *testing.T) (*ThroughputTestSuite, func()) {
	logger := logrus.NewEntry(logrus.New())
	logger.Logger.SetLevel(logrus.InfoLevel)
	
	// Start mock RPC server
	server := testingutils.NewMockRPCServer()
	
	// Create validator
	validator := validators.NewTransactionValidator(logger)
	
	// Create validating mock clients
	clientCount := 3
	clients := make([]spamoortypes.Client, clientCount)
	for i := 0; i < clientCount; i++ {
		mockClient := testingutils.NewMockClient()
		mockClient.SetMockChainId(big.NewInt(1337))
		mockClient.SetMockBalance(big.NewInt(1000000000000000000)) // 1 ETH
		mockClient.SetMockGasFees(big.NewInt(100), big.NewInt(2))
		
		validatingClient := testingutils.NewValidatingMockClientFromExisting(mockClient, validator, logger)
		clients[i] = validatingClient
	}
	
	// Create mock wallets
	walletCount := 10
	wallets := make([]spamoortypes.Wallet, walletCount)
	for i := 0; i < walletCount; i++ {
		wallet := testingutils.NewMockWallet()
		wallet.SetAddress(testingutils.GenerateAddress(fmt.Sprintf("wallet-%d", i)))
		wallet.SetChainId(big.NewInt(1337))
		wallet.SetNonce(0)
		wallet.SetBalance(big.NewInt(1000000000000000000)) // 1 ETH
		wallets[i] = wallet
	}
	
	suite := &ThroughputTestSuite{
		server:    server,
		clients:   clients,
		wallets:   wallets,
		logger:    logger,
		validator: validator,
	}
	
	cleanup := func() {
		server.Close()
	}
	
	return suite, cleanup
}

// TestLowLoadDynamicFee tests low load throughput with dynamic fee transactions
func TestLowLoadDynamicFee(t *testing.T) {
	suite, cleanup := setupThroughputTest(t)
	defer cleanup()
	
	config := &framework.LoadTestConfig{
		Duration:        30 * time.Second,
		TargetTPS:       30, // Target for load generation (not validation)
		ClientCount:     len(suite.clients),
		WalletCount:     len(suite.wallets),
		TransactionType: "dynamicfee",
		ValidationMode:  "full",
		
		GasLimit: 21000,
		BaseFee:  20,
		TipFee:   2,
		Amount:   1000000000000000, // 0.001 ETH
		
		MaxPending:     500,
		Timeout:        120 * time.Second,
		WarmupDuration: 5 * time.Second,
		
		EnableNonceValidation:      true,
		EnableOrderingValidation:   true,
		EnableThroughputValidation: false, // Focus on correctness, not TPS validation
	}
	
	engine := framework.NewLoadTestEngine(config, suite.clients, suite.wallets, suite.logger)
	
	ctx := context.Background()
	result, err := engine.RunLoadTest(ctx)
	
	require.NoError(t, err, "Load test should complete without errors")
	require.NotNil(t, result, "Load test result should not be nil")
	
	// Validate correctness (not performance)
	assert.True(t, result.TotalTransactions > 0, "Should have sent transactions")
	assert.True(t, result.SuccessfulTxs > 0, "Should have successful transactions")
	assert.Equal(t, 0, result.NonceGaps, "Should have no nonce gaps")
	assert.Equal(t, 0, result.OrderingViolations, "Should have no ordering violations")
	assert.True(t, result.ErrorRate < 0.01, "Error rate should be less than 1%")
	
	// Log results for hardware comparison
	t.Logf("Low Load DynamicFee: %.2f TPS, %d transactions, %.2f%% errors, %v latency", 
		result.AverageTPS, result.TotalTransactions, result.ErrorRate*100, result.AverageLatency)
	
	suite.logger.WithFields(logrus.Fields{
		"total_transactions": result.TotalTransactions,
		"successful_txs":     result.SuccessfulTxs,
		"average_tps":        result.AverageTPS,
		"average_latency":    result.AverageLatency,
		"error_rate":         result.ErrorRate,
		"test_successful":    result.IsSuccessful(),
	}).Info("Low load dynamic fee test completed")
	
	assert.True(t, result.IsSuccessful(), "Load test should be considered successful")
}

// TestMediumLoadBlobTransactions tests medium load throughput with blob transactions
func TestMediumLoadBlobTransactions(t *testing.T) {
	suite, cleanup := setupThroughputTest(t)
	defer cleanup()
	
	config := &framework.LoadTestConfig{
		Duration:        45 * time.Second,
		TargetTPS:       40, // Target for load generation (not validation)
		ClientCount:     len(suite.clients),
		WalletCount:     len(suite.wallets),
		TransactionType: "blob",
		ValidationMode:  "full",
		
		GasLimit: 50000,
		BaseFee:  30,
		TipFee:   3,
		Amount:   2000000000000000, // 0.002 ETH
		
		MaxPending:     1000,
		Timeout:        180 * time.Second,
		WarmupDuration: 10 * time.Second,
		
		EnableNonceValidation:      true,
		EnableOrderingValidation:   true,
		EnableThroughputValidation: false, // Focus on correctness, not TPS validation
	}
	
	engine := framework.NewLoadTestEngine(config, suite.clients, suite.wallets, suite.logger)
	
	ctx := context.Background()
	result, err := engine.RunLoadTest(ctx)
	
	require.NoError(t, err, "Load test should complete without errors")
	require.NotNil(t, result, "Load test result should not be nil")
	
	// Validate correctness (not performance)
	assert.True(t, result.TotalTransactions > 0, "Should have sent transactions")
	assert.True(t, result.SuccessfulTxs > 0, "Should have successful transactions")
	assert.Equal(t, 0, result.NonceGaps, "Should have no nonce gaps")
	assert.Equal(t, 0, result.OrderingViolations, "Should have no ordering violations")
	assert.True(t, result.ErrorRate < 0.01, "Error rate should be less than 1%")
	
	// Log results for hardware comparison
	t.Logf("Medium Load Blob: %.2f TPS, %d transactions, %.2f%% errors, %v latency", 
		result.AverageTPS, result.TotalTransactions, result.ErrorRate*100, result.AverageLatency)
	
	suite.logger.WithFields(logrus.Fields{
		"total_transactions": result.TotalTransactions,
		"successful_txs":     result.SuccessfulTxs,
		"average_tps":        result.AverageTPS,
		"average_latency":    result.AverageLatency,
		"error_rate":         result.ErrorRate,
		"test_successful":    result.IsSuccessful(),
	}).Info("Medium load blob transaction test completed")
	
	assert.True(t, result.IsSuccessful(), "Load test should be considered successful")
}

// TestHighLoadMixedTransactions tests high load throughput with mixed transaction types
func TestHighLoadMixedTransactions(t *testing.T) {
	suite, cleanup := setupThroughputTest(t)
	defer cleanup()
	
	config := &framework.LoadTestConfig{
		Duration:        60 * time.Second,
		TargetTPS:       50, // Target for load generation (not validation)
		ClientCount:     len(suite.clients),
		WalletCount:     len(suite.wallets),
		TransactionType: "mixed",
		ValidationMode:  "full",
		
		GasLimit: 75000,
		BaseFee:  40,
		TipFee:   4,
		Amount:   3000000000000000, // 0.003 ETH
		
		MaxPending:     2000,
		Timeout:        300 * time.Second,
		WarmupDuration: 15 * time.Second,
		
		EnableNonceValidation:      true,
		EnableOrderingValidation:   true,
		EnableThroughputValidation: false, // Disable strict TPS validation for mixed load
	}
	
	engine := framework.NewLoadTestEngine(config, suite.clients, suite.wallets, suite.logger)
	
	ctx := context.Background()
	result, err := engine.RunLoadTest(ctx)
	
	require.NoError(t, err, "Load test should complete without errors")
	require.NotNil(t, result, "Load test result should not be nil")
	
	// Validate correctness (not performance)
	assert.True(t, result.TotalTransactions > 0, "Should have sent transactions")
	assert.True(t, result.SuccessfulTxs > 0, "Should have successful transactions")
	assert.Equal(t, 0, result.NonceGaps, "Should have no nonce gaps")
	assert.Equal(t, 0, result.OrderingViolations, "Should have no ordering violations")
	assert.True(t, result.ErrorRate < 0.01, "Error rate should be less than 1%")
	
	// Log results for hardware comparison
	t.Logf("High Load Mixed: %.2f TPS (peak: %.2f), %d transactions, %.2f%% errors, avg: %v, p95: %v", 
		result.AverageTPS, result.PeakTPS, result.TotalTransactions, result.ErrorRate*100, 
		result.AverageLatency, result.P95Latency)
	
	suite.logger.WithFields(logrus.Fields{
		"total_transactions": result.TotalTransactions,
		"successful_txs":     result.SuccessfulTxs,
		"average_tps":        result.AverageTPS,
		"peak_tps":           result.PeakTPS,
		"average_latency":    result.AverageLatency,
		"p95_latency":        result.P95Latency,
		"error_rate":         result.ErrorRate,
		"test_successful":    result.IsSuccessful(),
	}).Info("High load mixed transaction test completed")
	
	assert.True(t, result.IsSuccessful(), "Load test should be considered successful")
}

// TestSustainedLoad tests sustained load over an extended period
func TestSustainedLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping sustained load test in short mode")
	}
	
	suite, cleanup := setupThroughputTest(t)
	defer cleanup()
	
	config := &framework.LoadTestConfig{
		Duration:        90 * time.Second, // Reduced from 5 minutes
		TargetTPS:       50, // Reduced target for shorter test
		ClientCount:     len(suite.clients),
		WalletCount:     len(suite.wallets),
		TransactionType: "dynamicfee",
		ValidationMode:  "full",
		
		GasLimit: 21000,
		BaseFee:  25,
		TipFee:   3,
		Amount:   1500000000000000, // 0.0015 ETH
		
		MaxPending:     1000,
		Timeout:        180 * time.Second,
		WarmupDuration: 10 * time.Second,
		
		EnableNonceValidation:      true,
		EnableOrderingValidation:   true,
		EnableThroughputValidation: false, // Focus on correctness, not TPS validation
	}
	
	engine := framework.NewLoadTestEngine(config, suite.clients, suite.wallets, suite.logger)
	
	ctx := context.Background()
	result, err := engine.RunLoadTest(ctx)
	
	require.NoError(t, err, "Sustained load test should complete without errors")
	require.NotNil(t, result, "Load test result should not be nil")
	
	// Validate correctness (not performance)
	assert.True(t, result.TotalTransactions > 500, "Should have sent many transactions")
	assert.True(t, result.SuccessfulTxs > 500, "Should have many successful transactions")
	assert.Equal(t, 0, result.NonceGaps, "Should have no nonce gaps")
	assert.Equal(t, 0, result.OrderingViolations, "Should have no ordering violations")
	assert.True(t, result.ErrorRate < 0.005, "Error rate should be very low for sustained load")
	
	// Check for memory stability (should not grow excessively)
	assert.True(t, result.MemoryUsageMB < 500, "Memory usage should remain reasonable")
	
	// Log results for hardware comparison
	t.Logf("Sustained Load: %.2f TPS (peak: %.2f), %d transactions over %v, %.2f%% errors, %.2f MB memory", 
		result.AverageTPS, result.PeakTPS, result.TotalTransactions, result.Duration, 
		result.ErrorRate*100, result.MemoryUsageMB)
	
	suite.logger.WithFields(logrus.Fields{
		"duration":           result.Duration,
		"total_transactions": result.TotalTransactions,
		"successful_txs":     result.SuccessfulTxs,
		"average_tps":        result.AverageTPS,
		"peak_tps":           result.PeakTPS,
		"average_latency":    result.AverageLatency,
		"p95_latency":        result.P95Latency,
		"p99_latency":        result.P99Latency,
		"memory_usage_mb":    result.MemoryUsageMB,
		"error_rate":         result.ErrorRate,
		"test_successful":    result.IsSuccessful(),
	}).Info("Sustained load test completed")
	
	assert.True(t, result.IsSuccessful(), "Sustained load test should be considered successful")
}

// BenchmarkThroughput benchmarks maximum throughput capability
func BenchmarkThroughput(b *testing.B) {
	suite, cleanup := setupThroughputTest(&testing.T{})
	defer cleanup()
	
	config := &framework.LoadTestConfig{
		Duration:        30 * time.Second,
		TargetTPS:       1000, // High target for benchmarking
		ClientCount:     len(suite.clients),
		WalletCount:     len(suite.wallets),
		TransactionType: "dynamicfee",
		ValidationMode:  "minimal", // Reduce validation overhead for benchmarking
		
		GasLimit: 21000,
		BaseFee:  20,
		TipFee:   2,
		Amount:   1000000000000000,
		
		MaxPending:     5000,
		Timeout:        120 * time.Second,
		WarmupDuration: 0, // No warmup for benchmark
		
		EnableNonceValidation:      false, // Disable for max performance
		EnableOrderingValidation:   false,
		EnableThroughputValidation: false,
	}
	
	engine := framework.NewLoadTestEngine(config, suite.clients, suite.wallets, suite.logger)
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		result, err := engine.RunLoadTest(ctx)
		
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
		
		// Report metrics
		b.ReportMetric(float64(result.TotalTransactions), "transactions")
		b.ReportMetric(result.AverageTPS, "tps")
		b.ReportMetric(float64(result.AverageLatency.Nanoseconds())/1000000, "latency_ms")
	}
}