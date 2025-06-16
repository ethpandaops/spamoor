//go:build loadtest

package examples

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/ethpandaops/spamoor/spamoortypes"
	"github.com/ethpandaops/spamoor/testing/load/framework"
	"github.com/ethpandaops/spamoor/testing/load/validators"
	testingutils "github.com/ethpandaops/spamoor/testing/utils"
)

// SimpleLoadTestExample demonstrates basic usage of the load testing framework
func SimpleLoadTestExample() {
	// Set up logging
	logger := logrus.NewEntry(logrus.New())
	logger.Logger.SetLevel(logrus.InfoLevel)
	
	// Create mock RPC server
	server := testingutils.NewMockRPCServer()
	defer server.Close()
	
	// Create transaction validator
	validator := validators.NewTransactionValidator(logger)
	
	// Create mock clients
	clientCount := 2
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
		wallet.SetAddress(testingutils.GenerateAddress(fmt.Sprintf("example-wallet-%d", i)))
		wallet.SetChainId(big.NewInt(1337))
		wallet.SetNonce(0)
		wallet.SetBalance(big.NewInt(1000000000000000000)) // 1 ETH
		wallets[i] = wallet
	}
	
	// Configure load test
	config := &framework.LoadTestConfig{
		Duration:        30 * time.Second,
		TargetTPS:       100,
		ClientCount:     clientCount,
		WalletCount:     walletCount,
		TransactionType: "dynamicfee",
		ValidationMode:  "full",
		
		GasLimit: 21000,
		BaseFee:  20,
		TipFee:   2,
		Amount:   1000000000000000, // 0.001 ETH
		
		MaxPending:     500,
		Timeout:        120 * time.Second,
		WarmupDuration: 10 * time.Second,
		
		EnableNonceValidation:      true,
		EnableOrderingValidation:   true,
		EnableThroughputValidation: true,
	}
	
	// Create and run load test
	engine := framework.NewLoadTestEngine(config, clients, wallets, logger)
	
	ctx := context.Background()
	result, err := engine.RunLoadTest(ctx)
	
	if err != nil {
		log.Fatalf("Load test failed: %v", err)
	}
	
	// Display results
	fmt.Printf("\n=== Load Test Results ===\n")
	fmt.Printf("Duration: %v\n", result.Duration)
	fmt.Printf("Total Transactions: %d\n", result.TotalTransactions)
	fmt.Printf("Successful Transactions: %d\n", result.SuccessfulTxs)
	fmt.Printf("Failed Transactions: %d\n", result.FailedTxs)
	fmt.Printf("Average TPS: %.2f\n", result.AverageTPS)
	fmt.Printf("Peak TPS: %.2f\n", result.PeakTPS)
	fmt.Printf("Average Latency: %v\n", result.AverageLatency)
	fmt.Printf("P95 Latency: %v\n", result.P95Latency)
	fmt.Printf("P99 Latency: %v\n", result.P99Latency)
	fmt.Printf("Error Rate: %.4f%%\n", result.ErrorRate*100)
	fmt.Printf("Memory Usage: %.2f MB\n", result.MemoryUsageMB)
	fmt.Printf("Goroutine Count: %d\n", result.GoroutineCount)
	fmt.Printf("Nonce Gaps: %d\n", result.NonceGaps)
	fmt.Printf("Ordering Violations: %d\n", result.OrderingViolations)
	fmt.Printf("Test Successful: %t\n", result.IsSuccessful())
	
	if len(result.ValidationErrors) > 0 {
		fmt.Printf("\nValidation Errors:\n")
		for _, err := range result.ValidationErrors {
			fmt.Printf("  - %s\n", err)
		}
	}
	
	if len(result.CriticalErrors) > 0 {
		fmt.Printf("\nCritical Errors:\n")
		for _, err := range result.CriticalErrors {
			fmt.Printf("  - %s\n", err)
		}
	}
}

// AdvancedLoadTestExample demonstrates advanced features
func AdvancedLoadTestExample() {
	logger := logrus.NewEntry(logrus.New())
	logger.Logger.SetLevel(logrus.InfoLevel)
	
	server := testingutils.NewMockRPCServer()
	defer server.Close()
	
	validator := validators.NewTransactionValidator(logger)
	
	// Create multiple client configurations
	configs := []struct {
		name        string
		clientCount int
		walletCount int
		txType      string
		targetTPS   int
		duration    time.Duration
	}{
		{"Small Load", 1, 5, "dynamicfee", 25, 20 * time.Second},
		{"Medium Load", 3, 15, "blob", 75, 30 * time.Second},
		{"Large Load", 5, 30, "mixed", 150, 45 * time.Second},
	}
	
	for _, cfg := range configs {
		fmt.Printf("\n=== Running %s Test ===\n", cfg.name)
		
		// Create clients and wallets for this configuration
		clients := make([]spamoortypes.Client, cfg.clientCount)
		for i := 0; i < cfg.clientCount; i++ {
			mockClient := testingutils.NewMockClient()
			mockClient.SetMockChainId(big.NewInt(1337))
			mockClient.SetMockBalance(big.NewInt(1000000000000000000))
			mockClient.SetMockGasFees(big.NewInt(100), big.NewInt(2))
			
			validatingClient := testingutils.NewValidatingMockClientFromExisting(mockClient, validator, logger)
			clients[i] = validatingClient
		}
		
		wallets := make([]spamoortypes.Wallet, cfg.walletCount)
		for i := 0; i < cfg.walletCount; i++ {
			wallet := testingutils.NewMockWallet()
			wallet.SetAddress(testingutils.GenerateAddress(fmt.Sprintf("%s-wallet-%d", cfg.name, i)))
			wallet.SetChainId(big.NewInt(1337))
			wallet.SetNonce(0)
			wallet.SetBalance(big.NewInt(1000000000000000000))
			wallets[i] = wallet
		}
		
		// Configure test
		config := &framework.LoadTestConfig{
			Duration:        cfg.duration,
			TargetTPS:       cfg.targetTPS,
			ClientCount:     cfg.clientCount,
			WalletCount:     cfg.walletCount,
			TransactionType: cfg.txType,
			ValidationMode:  "full",
			
			GasLimit: 50000, // Higher for complex transactions
			BaseFee:  30,
			TipFee:   3,
			Amount:   2000000000000000, // 0.002 ETH
			
			MaxPending:     uint64(cfg.targetTPS * 2),
			Timeout:        cfg.duration + 60*time.Second,
			WarmupDuration: 10 * time.Second,
			
			EnableNonceValidation:      true,
			EnableOrderingValidation:   true,
			EnableThroughputValidation: true,
		}
		
		// Reset validator for each test
		validator.Reset()
		
		// Run test
		engine := framework.NewLoadTestEngine(config, clients, wallets, logger)
		
		ctx := context.Background()
		result, err := engine.RunLoadTest(ctx)
		
		if err != nil {
			fmt.Printf("Test failed: %v\n", err)
			continue
		}
		
		// Display summary results
		fmt.Printf("Clients: %d, Wallets: %d, Type: %s\n", cfg.clientCount, cfg.walletCount, cfg.txType)
		fmt.Printf("Target TPS: %d, Actual TPS: %.2f\n", cfg.targetTPS, result.AverageTPS)
		fmt.Printf("Transactions: %d successful, %d failed\n", result.SuccessfulTxs, result.FailedTxs)
		fmt.Printf("Latency: avg=%v, p95=%v\n", result.AverageLatency, result.P95Latency)
		fmt.Printf("Validation: %d nonce gaps, %d ordering violations\n", result.NonceGaps, result.OrderingViolations)
		fmt.Printf("Success: %t\n", result.IsSuccessful())
	}
}