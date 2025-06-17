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

// ScalingTestSuite contains scaling validation tests
type ScalingTestSuite struct {
	server    *testingutils.MockRPCServer
	logger    *logrus.Entry
	validator *validators.TransactionValidator
}

// setupScalingTest sets up the test environment for scaling tests
func setupScalingTest(t *testing.T) (*ScalingTestSuite, func()) {
	logger := logrus.NewEntry(logrus.New())
	logger.Logger.SetLevel(logrus.InfoLevel)
	
	// Start mock RPC server
	server := testingutils.NewMockRPCServer()
	
	// Create validator
	validator := validators.NewTransactionValidator(logger)
	
	suite := &ScalingTestSuite{
		server:    server,
		logger:    logger,
		validator: validator,
	}
	
	cleanup := func() {
		server.Close()
	}
	
	return suite, cleanup
}

// createClients creates a specified number of validating mock clients
func (s *ScalingTestSuite) createClients(count int) []spamoortypes.Client {
	clients := make([]spamoortypes.Client, count)
	for i := 0; i < count; i++ {
		mockClient := testingutils.NewMockClient()
		mockClient.SetMockChainId(big.NewInt(1337))
		mockClient.SetMockBalance(big.NewInt(1000000000000000000)) // 1 ETH
		mockClient.SetMockGasFees(big.NewInt(100), big.NewInt(2))
		
		validatingClient := testingutils.NewValidatingMockClientFromExisting(mockClient, s.validator, s.logger)
		clients[i] = validatingClient
	}
	return clients
}

// createWallets creates a specified number of mock wallets
func (s *ScalingTestSuite) createWallets(count int) []spamoortypes.Wallet {
	wallets := make([]spamoortypes.Wallet, count)
	for i := 0; i < count; i++ {
		wallet := testingutils.NewMockWallet()
		wallet.SetAddress(testingutils.GenerateAddress(fmt.Sprintf("scaling-wallet-%d", i)))
		wallet.SetChainId(big.NewInt(1337))
		wallet.SetNonce(0)
		wallet.SetBalance(big.NewInt(1000000000000000000)) // 1 ETH
		wallets[i] = wallet
	}
	return wallets
}

// TestClientScaling tests performance with varying numbers of clients
func TestClientScaling(t *testing.T) {
	suite, cleanup := setupScalingTest(t)
	defer cleanup()
	
	clientCounts := []int{1, 2, 5, 10}
	baseWalletCount := 20
	baseTPS := 100
	
	results := make([]*framework.LoadTestResult, len(clientCounts))
	
	for i, clientCount := range clientCounts {
		t.Run(fmt.Sprintf("Clients_%d", clientCount), func(t *testing.T) {
			// Reset validator for each test
			suite.validator.Reset()
			
			clients := suite.createClients(clientCount)
			wallets := suite.createWallets(baseWalletCount)
			
			config := &framework.LoadTestConfig{
				Duration:        30 * time.Second,
				TargetTPS:       baseTPS,
				ClientCount:     clientCount,
				WalletCount:     baseWalletCount,
				TransactionType: "dynamicfee",
				ValidationMode:  "full",
				
				GasLimit: 21000,
				BaseFee:  20,
				TipFee:   2,
				Amount:   1000000000000000, // 0.001 ETH
				
				MaxPending:     1000,
				Timeout:        120 * time.Second,
				WarmupDuration: 5 * time.Second,
				
				EnableNonceValidation:      true,
				EnableOrderingValidation:   true,
				EnableThroughputValidation: false, // Focus on scaling behavior
			}
			
			engine := framework.NewLoadTestEngine(config, clients, wallets, suite.logger)
			
			ctx := context.Background()
			result, err := engine.RunLoadTest(ctx)
			
			require.NoError(t, err, "Load test should complete without errors")
			require.NotNil(t, result, "Load test result should not be nil")
			
			results[i] = result
			
			// Validate results
			assert.True(t, result.TotalTransactions > 0, "Should have sent transactions")
			assert.True(t, result.SuccessfulTxs > 0, "Should have successful transactions")
			assert.True(t, result.AverageTPS > 0, "Should have positive average TPS")
			assert.Equal(t, 0, result.NonceGaps, "Should have no nonce gaps")
			assert.Equal(t, 0, result.OrderingViolations, "Should have no ordering violations")
			
			suite.logger.WithFields(logrus.Fields{
				"client_count":       clientCount,
				"total_transactions": result.TotalTransactions,
				"successful_txs":     result.SuccessfulTxs,
				"average_tps":        result.AverageTPS,
				"average_latency":    result.AverageLatency,
				"error_rate":         result.ErrorRate,
			}).Info("Client scaling test completed")
		})
	}
	
	// Analyze correctness across scaling scenarios
	for i := 0; i < len(results); i++ {
		result := results[i]
		clientCount := clientCounts[i]
		
		// Validate transaction correctness regardless of performance
		assert.Equal(t, 0, result.NonceGaps, "Client scaling should maintain nonce consistency with %d clients", clientCount)
		assert.Equal(t, 0, result.OrderingViolations, "Client scaling should maintain transaction ordering with %d clients", clientCount)
		assert.True(t, result.ErrorRate < 0.01, "Error rate should remain low with %d clients", clientCount)
		
		// Log performance for comparison (no assertions)
		t.Logf("Client scaling result [%d clients]: %.2f TPS, %d transactions, %.2f%% errors", 
			clientCount, result.AverageTPS, result.TotalTransactions, result.ErrorRate*100)
	}
}

// TestWalletScaling tests performance with varying numbers of wallets
func TestWalletScaling(t *testing.T) {
	suite, cleanup := setupScalingTest(t)
	defer cleanup()
	
	walletCounts := []int{5, 10, 25, 50, 100}
	baseClientCount := 3
	baseTPS := 150
	
	results := make([]*framework.LoadTestResult, len(walletCounts))
	
	for i, walletCount := range walletCounts {
		t.Run(fmt.Sprintf("Wallets_%d", walletCount), func(t *testing.T) {
			// Reset validator for each test
			suite.validator.Reset()
			
			clients := suite.createClients(baseClientCount)
			wallets := suite.createWallets(walletCount)
			
			config := &framework.LoadTestConfig{
				Duration:        45 * time.Second,
				TargetTPS:       baseTPS,
				ClientCount:     baseClientCount,
				WalletCount:     walletCount,
				TransactionType: "dynamicfee",
				ValidationMode:  "full",
				
				GasLimit: 21000,
				BaseFee:  25,
				TipFee:   3,
				Amount:   1500000000000000, // 0.0015 ETH
				
				MaxPending:     1500,
				Timeout:        180 * time.Second,
				WarmupDuration: 10 * time.Second,
				
				EnableNonceValidation:      true,
				EnableOrderingValidation:   true,
				EnableThroughputValidation: false, // Focus on scaling behavior
			}
			
			engine := framework.NewLoadTestEngine(config, clients, wallets, suite.logger)
			
			ctx := context.Background()
			result, err := engine.RunLoadTest(ctx)
			
			require.NoError(t, err, "Load test should complete without errors")
			require.NotNil(t, result, "Load test result should not be nil")
			
			results[i] = result
			
			// Validate results
			assert.True(t, result.TotalTransactions > 0, "Should have sent transactions")
			assert.True(t, result.SuccessfulTxs > 0, "Should have successful transactions")
			assert.True(t, result.AverageTPS > 0, "Should have positive average TPS")
			assert.Equal(t, 0, result.NonceGaps, "Should have no nonce gaps")
			assert.Equal(t, 0, result.OrderingViolations, "Should have no ordering violations")
			
			suite.logger.WithFields(logrus.Fields{
				"wallet_count":       walletCount,
				"total_transactions": result.TotalTransactions,
				"successful_txs":     result.SuccessfulTxs,
				"average_tps":        result.AverageTPS,
				"average_latency":    result.AverageLatency,
				"memory_usage_mb":    result.MemoryUsageMB,
				"error_rate":         result.ErrorRate,
			}).Info("Wallet scaling test completed")
		})
	}
	
	// Analyze correctness across wallet scaling scenarios
	for i := 0; i < len(results); i++ {
		result := results[i]
		walletCount := walletCounts[i]
		
		// Validate transaction correctness regardless of performance
		assert.Equal(t, 0, result.NonceGaps, "Wallet scaling should maintain nonce consistency with %d wallets", walletCount)
		assert.Equal(t, 0, result.OrderingViolations, "Wallet scaling should maintain transaction ordering with %d wallets", walletCount)
		assert.True(t, result.ErrorRate < 0.01, "Error rate should remain low with %d wallets", walletCount)
		
		// Check memory usage doesn't grow excessively
		if result.MemoryUsageMB > 0 {
			assert.True(t, result.MemoryUsageMB < 500, "Memory usage should remain reasonable with %d wallets", walletCount)
		}
		
		// Log performance for comparison (no assertions)
		t.Logf("Wallet scaling result [%d wallets]: %.2f TPS, %d transactions, %.2f MB memory, %.2f%% errors", 
			walletCount, result.AverageTPS, result.TotalTransactions, result.MemoryUsageMB, result.ErrorRate*100)
	}
}

// TestCombinedScaling tests performance with varying both clients and wallets
func TestCombinedScaling(t *testing.T) {
	t.Log("📈 Starting combined scaling test...")
	suite, cleanup := setupScalingTest(t)
	defer cleanup()
	
	scenarios := []struct {
		name         string
		clientCount  int
		walletCount  int
		targetTPS    int
	}{
		{"Small_1c_5w", 1, 5, 30},
		{"Medium_3c_20w", 3, 20, 50},
		{"Large_5c_50w", 5, 50, 70},
		// Removed XLarge scenario to reduce test time
	}
	
	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			t.Logf("  🔄 Testing scenario: %s (%d clients, %d wallets, %d TPS)", 
				scenario.name, scenario.clientCount, scenario.walletCount, scenario.targetTPS)
			// Reset validator for each test
			suite.validator.Reset()
			
			clients := suite.createClients(scenario.clientCount)
			wallets := suite.createWallets(scenario.walletCount)
			
			config := &framework.LoadTestConfig{
				Duration:        30 * time.Second,
				TargetTPS:       scenario.targetTPS,
				ClientCount:     scenario.clientCount,
				WalletCount:     scenario.walletCount,
				TransactionType: "mixed",
				ValidationMode:  "full", // Use full validation for correctness testing
				
				GasLimit: 50000,
				BaseFee:  30,
				TipFee:   4,
				Amount:   2000000000000000, // 0.002 ETH
				
				MaxPending:     uint64(scenario.targetTPS * 3), // Increased buffer
				Timeout:        120 * time.Second,
				WarmupDuration: 5 * time.Second,
				
				EnableNonceValidation:      true,  // Enable for correctness validation
				EnableOrderingValidation:   true,  // Enable for correctness validation
				EnableThroughputValidation: false, // Disable TPS validation
			}
			
			engine := framework.NewLoadTestEngine(config, clients, wallets, suite.logger)
			
			ctx := context.Background()
			result, err := engine.RunLoadTest(ctx)
			
			require.NoError(t, err, "Load test should complete without errors")
			require.NotNil(t, result, "Load test result should not be nil")
			
			// Validate correctness (not performance)
			assert.True(t, result.TotalTransactions > 0, "Should have sent transactions")
			assert.True(t, result.SuccessfulTxs > 0, "Should have successful transactions")
			assert.Equal(t, 0, result.NonceGaps, "Should have no nonce gaps")
			assert.Equal(t, 0, result.OrderingViolations, "Should have no ordering violations")
			assert.True(t, result.ErrorRate < 0.01, "Error rate should be very low")
			
			suite.logger.WithFields(logrus.Fields{
				"scenario":           scenario.name,
				"client_count":       scenario.clientCount,
				"wallet_count":       scenario.walletCount,
				"target_tps":         scenario.targetTPS,
				"actual_tps":         result.AverageTPS,
				"peak_tps":           result.PeakTPS,
				"total_transactions": result.TotalTransactions,
				"successful_txs":     result.SuccessfulTxs,
				"average_latency":    result.AverageLatency,
				"p95_latency":        result.P95Latency,
				"memory_usage_mb":    result.MemoryUsageMB,
				"goroutine_count":    result.GoroutineCount,
				"error_rate":         result.ErrorRate,
				"test_successful":    result.IsSuccessful(),
			}).Info("Combined scaling test completed")
			
			assert.True(t, result.IsSuccessful(), "Combined scaling test should be successful")
		})
	}
}

// TestResourceUtilizationScaling tests how resource usage scales with load
func TestResourceUtilizationScaling(t *testing.T) {
	suite, cleanup := setupScalingTest(t)
	defer cleanup()
	
	scenarios := []struct {
		name        string
		clientCount int
		walletCount int
		targetTPS   int
		duration    time.Duration
	}{
		{"Baseline", 2, 10, 50, 30 * time.Second},
		{"2x_Load", 4, 20, 100, 30 * time.Second},
		{"4x_Load", 8, 40, 200, 30 * time.Second},
	}
	
	baselineMemory := 0.0
	baselineGoroutines := 0
	
	for i, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			// Reset validator for each test
			suite.validator.Reset()
			
			clients := suite.createClients(scenario.clientCount)
			wallets := suite.createWallets(scenario.walletCount)
			
			config := &framework.LoadTestConfig{
				Duration:        scenario.duration,
				TargetTPS:       scenario.targetTPS,
				ClientCount:     scenario.clientCount,
				WalletCount:     scenario.walletCount,
				TransactionType: "dynamicfee",
				ValidationMode:  "full",
				
				GasLimit: 21000,
				BaseFee:  20,
				TipFee:   2,
				Amount:   1000000000000000,
				
				MaxPending:     uint64(scenario.targetTPS * 3),
				Timeout:        120 * time.Second,
				WarmupDuration: 5 * time.Second,
				
				EnableNonceValidation:      true,
				EnableOrderingValidation:   true,
				EnableThroughputValidation: false,
			}
			
			engine := framework.NewLoadTestEngine(config, clients, wallets, suite.logger)
			
			ctx := context.Background()
			result, err := engine.RunLoadTest(ctx)
			
			require.NoError(t, err, "Resource utilization test should complete without errors")
			require.NotNil(t, result, "Load test result should not be nil")
			
			// Store baseline metrics
			if i == 0 {
				baselineMemory = result.MemoryUsageMB
				baselineGoroutines = result.GoroutineCount
			}
			
			// Validate basic functionality
			assert.True(t, result.TotalTransactions > 0, "Should have sent transactions")
			assert.True(t, result.SuccessfulTxs > 0, "Should have successful transactions")
			assert.Equal(t, 0, result.NonceGaps, "Should have no nonce gaps")
			assert.Equal(t, 0, result.OrderingViolations, "Should have no ordering violations")
			
			// Resource usage validation
			if i > 0 && baselineMemory > 0 {
				memoryGrowthRatio := result.MemoryUsageMB / baselineMemory
				loadMultiplier := float64(scenario.targetTPS) / float64(scenarios[0].targetTPS)
				
				// Memory should not grow more than 3x the load multiplier
				maxExpectedMemoryRatio := loadMultiplier * 3
				assert.True(t, memoryGrowthRatio <= maxExpectedMemoryRatio,
					"Memory growth should be reasonable: %.2fx growth for %.2fx load", 
					memoryGrowthRatio, loadMultiplier)
			}
			
			if i > 0 && baselineGoroutines > 0 {
				goroutineGrowthRatio := float64(result.GoroutineCount) / float64(baselineGoroutines)
				loadMultiplier := float64(scenario.targetTPS) / float64(scenarios[0].targetTPS)
				
				// Goroutine count should not grow excessively
				maxExpectedGoroutineRatio := loadMultiplier * 2
				assert.True(t, goroutineGrowthRatio <= maxExpectedGoroutineRatio,
					"Goroutine growth should be reasonable: %.2fx growth for %.2fx load", 
					goroutineGrowthRatio, loadMultiplier)
			}
			
			suite.logger.WithFields(logrus.Fields{
				"scenario":           scenario.name,
				"client_count":       scenario.clientCount,
				"wallet_count":       scenario.walletCount,
				"target_tps":         scenario.targetTPS,
				"actual_tps":         result.AverageTPS,
				"total_transactions": result.TotalTransactions,
				"memory_usage_mb":    result.MemoryUsageMB,
				"goroutine_count":    result.GoroutineCount,
				"error_rate":         result.ErrorRate,
			}).Info("Resource utilization scaling test completed")
		})
	}
}