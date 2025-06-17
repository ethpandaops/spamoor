//go:build loadtest

package load

import (
	"context"
	"fmt"
	"math/big"
	"os"
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

// LoadTestSuite provides a complete load testing environment
type LoadTestSuite struct {
	server    *testingutils.MockRPCServer
	clients   []spamoortypes.Client
	wallets   []spamoortypes.Wallet
	logger    *logrus.Entry
	validator *validators.TransactionValidator
}

// SetupLoadTestSuite creates a comprehensive load test environment
func SetupLoadTestSuite(t *testing.T, clientCount, walletCount int) (*LoadTestSuite, func()) {
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.SetOutput(os.Stderr) // Force output to stderr for immediate display
	logger := logrus.NewEntry(log)
	
	// Start mock RPC server
	server := testingutils.NewMockRPCServer()
	
	// Create validator
	validator := validators.NewTransactionValidator(logger)
	
	// Create validating mock clients
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
	wallets := make([]spamoortypes.Wallet, walletCount)
	for i := 0; i < walletCount; i++ {
		wallet := testingutils.NewMockWallet()
		wallet.SetAddress(testingutils.GenerateAddress(fmt.Sprintf("load-wallet-%d", i)))
		wallet.SetChainId(big.NewInt(1337))
		wallet.SetNonce(0)
		wallet.SetBalance(big.NewInt(1000000000000000000)) // 1 ETH
		wallets[i] = wallet
	}
	
	suite := &LoadTestSuite{
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

// TestBasicLoadTest tests basic load testing functionality
func TestBasicLoadTest(t *testing.T) {
	t.Log("🚀 Starting basic load test...")
	suite, cleanup := SetupLoadTestSuite(t, 2, 10)
	defer cleanup()
	
	config := framework.DefaultLoadTestConfig()
	config.Duration = 20 * time.Second
	config.TargetTPS = 30 // Lower target for basic test
	config.WarmupDuration = 5 * time.Second
	config.EnableThroughputValidation = false // Focus on functionality, not performance
	
	engine := framework.NewLoadTestEngine(config, suite.clients, suite.wallets, suite.logger)
	
	ctx := context.Background()
	result, err := engine.RunLoadTest(ctx)
	
	require.NoError(t, err, "Basic load test should complete without errors")
	require.NotNil(t, result, "Load test result should not be nil")
	
	// Validate basic functionality (correctness, not performance)
	assert.True(t, result.TotalTransactions > 0, "Should have sent transactions")
	assert.True(t, result.SuccessfulTxs > 0, "Should have successful transactions")
	assert.Equal(t, 0, result.NonceGaps, "Should have no nonce gaps")
	assert.Equal(t, 0, result.OrderingViolations, "Should have no ordering violations")
	assert.True(t, result.ErrorRate < 0.05, "Error rate should be reasonable")
	
	// Log results for hardware comparison
	t.Logf("Basic Load Test: %.2f TPS, %d transactions, %.2f%% errors, %v latency", 
		result.AverageTPS, result.TotalTransactions, result.ErrorRate*100, result.AverageLatency)
	
	suite.logger.WithFields(logrus.Fields{
		"total_transactions": result.TotalTransactions,
		"successful_txs":     result.SuccessfulTxs,
		"average_tps":        result.AverageTPS,
		"average_latency":    result.AverageLatency,
		"error_rate":         result.ErrorRate,
		"test_successful":    result.IsSuccessful(),
	}).Info("Basic load test completed")
	
	assert.True(t, result.IsSuccessful(), "Basic load test should be considered successful")
}

// TestTransactionTypes tests different transaction types under load
func TestTransactionTypes(t *testing.T) {
	t.Log("🔧 Starting transaction types test...")
	suite, cleanup := SetupLoadTestSuite(t, 3, 15)
	defer cleanup()
	
	transactionTypes := []string{"dynamicfee", "blob", "setcode", "mixed"}
	
	for _, txType := range transactionTypes {
		t.Run(fmt.Sprintf("TransactionType_%s", txType), func(t *testing.T) {
			t.Logf("  📊 Testing transaction type: %s", txType)
			// Reset validator for each test
			suite.validator.Reset()
			
			config := &framework.LoadTestConfig{
				Duration:        25 * time.Second,
				TargetTPS:       75,
				ClientCount:     len(suite.clients),
				WalletCount:     len(suite.wallets),
				TransactionType: txType,
				ValidationMode:  "full",
				
				GasLimit: 50000, // Higher gas limit for complex transactions
				BaseFee:  25,
				TipFee:   3,
				Amount:   1500000000000000, // 0.0015 ETH
				
				MaxPending:     800,
				Timeout:        120 * time.Second,
				WarmupDuration: 8 * time.Second,
				
				EnableNonceValidation:      true,
				EnableOrderingValidation:   true,
				EnableThroughputValidation: false, // Different tx types may have different perf
			}
			
			engine := framework.NewLoadTestEngine(config, suite.clients, suite.wallets, suite.logger)
			
			ctx := context.Background()
			result, err := engine.RunLoadTest(ctx)
			
			require.NoError(t, err, "Transaction type test should complete without errors")
			require.NotNil(t, result, "Load test result should not be nil")
			
			// Validate correctness (not performance)
			assert.True(t, result.TotalTransactions > 0, "Should have sent transactions")
			assert.True(t, result.SuccessfulTxs > 0, "Should have successful transactions")
			assert.Equal(t, 0, result.NonceGaps, "Should have no nonce gaps")
			assert.Equal(t, 0, result.OrderingViolations, "Should have no ordering violations")
			assert.True(t, result.ErrorRate < 0.01, "Error rate should be very low")
			
			// Log results for hardware comparison
			t.Logf("Transaction Type %s: %.2f TPS, %d transactions, %.2f%% errors, %v latency", 
				txType, result.AverageTPS, result.TotalTransactions, result.ErrorRate*100, result.AverageLatency)
			
			suite.logger.WithFields(logrus.Fields{
				"transaction_type":   txType,
				"total_transactions": result.TotalTransactions,
				"successful_txs":     result.SuccessfulTxs,
				"average_tps":        result.AverageTPS,
				"average_latency":    result.AverageLatency,
				"error_rate":         result.ErrorRate,
				"test_successful":    result.IsSuccessful(),
			}).Info("Transaction type load test completed")
		})
	}
}

// TestValidationModes tests different validation modes
func TestValidationModes(t *testing.T) {
	suite, cleanup := SetupLoadTestSuite(t, 2, 10)
	defer cleanup()
	
	validationScenarios := []struct {
		name                       string
		enableNonceValidation      bool
		enableOrderingValidation   bool
		enableThroughputValidation bool
	}{
		{"Full_Validation", true, true, false},
		{"Nonce_Only", true, false, false},
		{"Ordering_Only", false, true, false},
		{"Minimal_Validation", false, false, false},
	}
	
	for _, scenario := range validationScenarios {
		t.Run(scenario.name, func(t *testing.T) {
			// Reset validator for each test
			suite.validator.Reset()
			
			config := &framework.LoadTestConfig{
				Duration:        30 * time.Second,
				TargetTPS:       80,
				ClientCount:     len(suite.clients),
				WalletCount:     len(suite.wallets),
				TransactionType: "dynamicfee",
				ValidationMode:  "custom",
				
				GasLimit: 21000,
				BaseFee:  20,
				TipFee:   2,
				Amount:   1000000000000000,
				
				MaxPending:     1000,
				Timeout:        120 * time.Second,
				WarmupDuration: 10 * time.Second,
				
				EnableNonceValidation:      scenario.enableNonceValidation,
				EnableOrderingValidation:   scenario.enableOrderingValidation,
				EnableThroughputValidation: scenario.enableThroughputValidation,
			}
			
			engine := framework.NewLoadTestEngine(config, suite.clients, suite.wallets, suite.logger)
			
			ctx := context.Background()
			result, err := engine.RunLoadTest(ctx)
			
			require.NoError(t, err, "Validation mode test should complete without errors")
			require.NotNil(t, result, "Load test result should not be nil")
			
			// Validate results (no TPS assertions)
			assert.True(t, result.TotalTransactions > 0, "Should have sent transactions")
			assert.True(t, result.SuccessfulTxs > 0, "Should have successful transactions")
			
			// Validate based on enabled validations
			if scenario.enableNonceValidation {
				assert.Equal(t, 0, result.NonceGaps, "Should have no nonce gaps when validation enabled")
			}
			
			if scenario.enableOrderingValidation {
				assert.Equal(t, 0, result.OrderingViolations, "Should have no ordering violations when validation enabled")
			}
			
			// Log results for hardware comparison
			t.Logf("Validation Mode %s: %.2f TPS, %d transactions, %.2f%% errors, %v latency", 
				scenario.name, result.AverageTPS, result.TotalTransactions, result.ErrorRate*100, result.AverageLatency)
			
			suite.logger.WithFields(logrus.Fields{
				"validation_mode":    scenario.name,
				"nonce_validation":   scenario.enableNonceValidation,
				"ordering_validation": scenario.enableOrderingValidation,
				"throughput_validation": scenario.enableThroughputValidation,
				"total_transactions": result.TotalTransactions,
				"successful_txs":     result.SuccessfulTxs,
				"average_tps":        result.AverageTPS,
				"average_latency":    result.AverageLatency,
				"error_rate":         result.ErrorRate,
				"test_successful":    result.IsSuccessful(),
			}).Info("Validation mode test completed")
		})
	}
}

// TestErrorHandling tests load test behavior under error conditions
func TestErrorHandling(t *testing.T) {
	suite, cleanup := SetupLoadTestSuite(t, 2, 5)
	defer cleanup()
	
	// Inject errors into one of the clients
	if mockClient, ok := suite.clients[0].(*testingutils.ValidatingMockClient); ok {
		// Simulate intermittent failures (20% error rate)
		mockClient.MockClient.SetMockError(fmt.Errorf("simulated network error"))
	}
	
	config := &framework.LoadTestConfig{
		Duration:        20 * time.Second,
		TargetTPS:       50,
		ClientCount:     len(suite.clients),
		WalletCount:     len(suite.wallets),
		TransactionType: "dynamicfee",
		ValidationMode:  "full",
		
		GasLimit: 21000,
		BaseFee:  20,
		TipFee:   2,
		Amount:   1000000000000000,
		
		MaxPending:     500,
		Timeout:        120 * time.Second,
		WarmupDuration: 5 * time.Second,
		
		EnableNonceValidation:      true,
		EnableOrderingValidation:   true,
		EnableThroughputValidation: false, // Don't validate TPS with errors
	}
	
	engine := framework.NewLoadTestEngine(config, suite.clients, suite.wallets, suite.logger)
	
	ctx := context.Background()
	result, err := engine.RunLoadTest(ctx)
	
	require.NoError(t, err, "Error handling test should complete without errors")
	require.NotNil(t, result, "Load test result should not be nil")
	
	// With errors injected, we expect some failures
	assert.True(t, result.TotalTransactions > 0, "Should have attempted transactions")
	assert.True(t, result.FailedTxs > 0, "Should have some failed transactions due to injected errors")
	assert.True(t, result.ErrorRate > 0, "Should have non-zero error rate")
	assert.True(t, result.ErrorRate < 0.6, "Error rate should not be excessive")
	
	// Check that we still have some successful transactions
	assert.True(t, result.SuccessfulTxs > 0, "Should still have some successful transactions")
	
	suite.logger.WithFields(logrus.Fields{
		"total_transactions": result.TotalTransactions,
		"successful_txs":     result.SuccessfulTxs,
		"failed_txs":         result.FailedTxs,
		"error_rate":         result.ErrorRate,
		"average_tps":        result.AverageTPS,
		"errors_by_type":     result.ErrorsByType,
	}).Info("Error handling test completed")
}

// TestConfigurationValidation tests various configuration scenarios
func TestConfigurationValidation(t *testing.T) {
	suite, cleanup := SetupLoadTestSuite(t, 1, 5)
	defer cleanup()
	
	// Test with very low TPS
	t.Run("Low_TPS", func(t *testing.T) {
		config := framework.DefaultLoadTestConfig()
		config.Duration = 15 * time.Second
		config.TargetTPS = 5 // Very low TPS
		config.WarmupDuration = 2 * time.Second
		
		engine := framework.NewLoadTestEngine(config, suite.clients, suite.wallets, suite.logger)
		
		ctx := context.Background()
		result, err := engine.RunLoadTest(ctx)
		
		require.NoError(t, err)
		assert.True(t, result.TotalTransactions > 0)
		assert.True(t, result.AverageTPS <= 10) // Should respect low TPS setting
	})
	
	// Test with very short duration
	t.Run("Short_Duration", func(t *testing.T) {
		config := framework.DefaultLoadTestConfig()
		config.Duration = 5 * time.Second
		config.TargetTPS = 100
		config.WarmupDuration = 1 * time.Second
		
		engine := framework.NewLoadTestEngine(config, suite.clients, suite.wallets, suite.logger)
		
		ctx := context.Background()
		result, err := engine.RunLoadTest(ctx)
		
		require.NoError(t, err)
		assert.True(t, result.Duration <= 10*time.Second) // Should respect short duration
	})
	
	// Test with high pending limit
	t.Run("High_Pending_Limit", func(t *testing.T) {
		config := framework.DefaultLoadTestConfig()
		config.Duration = 15 * time.Second
		config.TargetTPS = 200
		config.MaxPending = 5000 // Very high pending limit
		config.WarmupDuration = 3 * time.Second
		
		engine := framework.NewLoadTestEngine(config, suite.clients, suite.wallets, suite.logger)
		
		ctx := context.Background()
		result, err := engine.RunLoadTest(ctx)
		
		require.NoError(t, err)
		assert.True(t, result.TotalTransactions > 0)
		// Should handle high pending limit without issues
	})
}

// BenchmarkLoadTestFramework benchmarks the load testing framework itself
func BenchmarkLoadTestFramework(b *testing.B) {
	suite, cleanup := SetupLoadTestSuite(&testing.T{}, 2, 10)
	defer cleanup()
	
	config := &framework.LoadTestConfig{
		Duration:        10 * time.Second,
		TargetTPS:       200,
		ClientCount:     len(suite.clients),
		WalletCount:     len(suite.wallets),
		TransactionType: "dynamicfee",
		ValidationMode:  "minimal", // Minimal validation for benchmarking
		
		GasLimit: 21000,
		BaseFee:  20,
		TipFee:   2,
		Amount:   1000000000000000,
		
		MaxPending:     2000,
		Timeout:        60 * time.Second,
		WarmupDuration: 0, // No warmup for benchmark
		
		EnableNonceValidation:      false,
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
		b.ReportMetric(result.MemoryUsageMB, "memory_mb")
		b.ReportMetric(float64(result.GoroutineCount), "goroutines")
	}
}