# Load Testing Framework - Technical Overview

## High-Level Architecture

The load testing framework is designed to simulate high-throughput Ethereum transaction scenarios with comprehensive validation. It uses a **producer-consumer pattern** with **concurrent transaction generation** and **real-time validation**.

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   LoadTestEngine │───▶│  MetricsCollector │───▶│ LoadTestResult  │
│                 │    │                  │    │                 │
│ ┌─────────────┐ │    │ ┌──────────────┐ │    │ • TPS Metrics   │
│ │ Transaction │ │    │ │ Performance  │ │    │ • Latency Stats │
│ │ Generators  │ │    │ │ Monitoring   │ │    │ • Error Rates   │
│ └─────────────┘ │    │ └──────────────┘ │    │ • Validation    │
└─────────────────┘    └──────────────────┘    └─────────────────┘
         │                        │
         ▼                        ▼
┌─────────────────┐    ┌──────────────────┐
│TransactionValidator│    │  ValidatingMock  │
│                 │    │     Clients      │
│ • Nonce Gaps    │◀───│                  │
│ • Ordering      │    │ • Tx Capture     │
│ • Real-time     │    │ • Mock Responses │
└─────────────────┘    └──────────────────┘
```

## Key Components You Need to Know

### 1. **LoadTestEngine** - The Orchestrator
```go
engine := framework.NewLoadTestEngine(config, clients, wallets, logger)
result, err := engine.RunLoadTest(ctx)
```

**What it does:**
- Orchestrates the entire load test lifecycle
- Manages concurrent transaction generation across multiple goroutines
- Coordinates warmup → main phase → validation → reporting
- Controls transaction rate limiting to achieve target TPS

**Key phases:**
1. **Warmup Phase** - Pre-generates transactions to reach steady state
2. **Main Phase** - Sustains target TPS for the configured duration
3. **Validation Phase** - Performs comprehensive transaction integrity checks

### 2. **LoadTestConfig** - Your Control Panel
```go
config := &framework.LoadTestConfig{
    Duration:        30 * time.Second,    // How long to run
    TargetTPS:       50,                  // Transactions per second
    ClientCount:     3,                   // Number of RPC clients
    WalletCount:     15,                  // Number of wallets (nonce sources)
    TransactionType: "dynamicfee",        // Transaction type to generate
    
    // Validation controls
    EnableNonceValidation:      true,     // Detect nonce gaps
    EnableOrderingValidation:   true,     // Detect ordering violations
    EnableThroughputValidation: false,    // Strict TPS validation
}
```

**Transaction Types Available:**
- `"dynamicfee"` - EIP-1559 transactions
- `"blob"` - EIP-4844 blob transactions  
- `"setcode"` - Contract deployment transactions
- `"mixed"` - Random mix of all types

### 3. **TransactionValidator** - The Quality Assurance
```go
validator := validators.NewTransactionValidator(logger)
```

**What it validates:**
- **Nonce Gaps** - Detects missing nonces in transaction sequences
- **Ordering Violations** - Ensures transactions maintain proper chronological order
- **Real-time Monitoring** - Catches issues during test execution, not just at the end

**How it works:**
- Captures every submitted transaction in real-time
- Maintains per-wallet nonce tracking
- Performs global ordering validation at test completion
- Reports violations with detailed error context

### 4. **MetricsCollector** - The Performance Monitor
```go
// Automatically embedded in LoadTestEngine
// Collects metrics throughout test execution
```

**Metrics Collected:**
- **Throughput**: Average TPS, Peak TPS, Transaction counts
- **Latency**: Average, P95, P99 percentiles
- **Resource Usage**: Memory consumption, Goroutine counts
- **Error Rates**: Success/failure ratios, Error categorization

### 5. **Mock Infrastructure** - The Test Environment

**ValidatingMockClient:**
```go
mockClient := testingutils.NewMockClient()
validatingClient := testingutils.NewValidatingMockClientFromExisting(
    mockClient, validator, logger)
```

**MockWallet:**
```go
wallet := testingutils.NewMockWallet()
wallet.SetChainId(big.NewInt(1337))
wallet.SetNonce(0)
```

**What the mocks provide:**
- Deterministic transaction responses
- Configurable error injection for testing edge cases
- Real transaction building and signing (using actual Ethereum transaction types)
- Predictable gas fee and balance management

## How to Use the Framework

### Basic Usage Pattern

```go
// 1. Set up test environment
suite, cleanup := SetupLoadTestSuite(t, clientCount, walletCount)
defer cleanup()

// 2. Configure your test
config := &framework.LoadTestConfig{
    Duration:        30 * time.Second,
    TargetTPS:       50,
    TransactionType: "dynamicfee",
    // ... other settings
}

// 3. Run the test
engine := framework.NewLoadTestEngine(config, suite.clients, suite.wallets, suite.logger)
result, err := engine.RunLoadTest(ctx)

// 4. Analyze results
assert.NoError(t, err)
assert.Equal(t, 0, result.NonceGaps)
assert.True(t, result.AverageTPS > expectedMinTPS)
```

### Running Load Tests

```bash
# Run all load tests (includes 15-minute timeout)
make test-load

# Run specific test
go test -tags=loadtest -v ./testing/load -run TestBasicLoadTest

# Run with custom timeout
go test -tags=loadtest -v ./testing/load -timeout 15m

# Run specific scenario
go test -tags=loadtest -v ./testing/load/scenarios -run TestLowLoadDynamicFee
```

## Available Test Scenarios

### Core Tests (`testing/load/load_test.go`)
- **TestBasicLoadTest** - Basic functionality validation
- **TestTransactionTypes** - Tests all transaction types (dynamicfee, blob, setcode, mixed)
- **TestValidationModes** - Different validation configurations
- **TestErrorHandling** - Error injection and recovery testing
- **TestConfigurationValidation** - Edge case configuration testing

### Throughput Tests (`testing/load/scenarios/throughput_test.go`)
- **TestLowLoadDynamicFee** - Low-intensity baseline testing
- **TestMediumLoadBlobTransactions** - Medium load with complex transactions
- **TestHighLoadMixedTransactions** - High load with mixed transaction types
- **TestSustainedLoad** - Extended duration testing (5 minutes)
- **BenchmarkThroughput** - Maximum performance benchmarking

### Scaling Tests (`testing/load/scenarios/scaling_test.go`)
- **TestClientScaling** - Performance with varying client counts
- **TestWalletScaling** - Performance with varying wallet counts
- **TestCombinedScaling** - Combined client and wallet scaling
- **TestResourceUtilizationScaling** - Memory and resource usage validation

### Example Usage (`testing/load/examples/simple_load_test.go`)
- **SimpleLoadTestExample** - Basic framework usage example
- **AdvancedLoadTestExample** - Advanced configuration examples

## Key Design Decisions

### **Why Build Tags?**
Load tests are computationally intensive and time-consuming. Build tags (`//go:build loadtest`) ensure they don't block regular CI/CD pipelines while remaining easily accessible for performance testing.

### **Why Mock Infrastructure?**
- **Deterministic Results** - Tests produce consistent, reproducible results
- **High Throughput** - No network latency or external dependencies
- **Error Injection** - Can simulate various failure scenarios
- **Transaction Integrity** - Uses real Ethereum transaction signing and validation

### **Why Real-Time Validation?**
Traditional load testing often only checks final results. This framework validates transaction integrity **during execution**, catching issues like:
- Race conditions in nonce management
- Out-of-order transaction processing
- Transaction loss or duplication

## Performance Characteristics

The framework focuses on **correctness validation** rather than specific performance targets:
- **Transaction Integrity**: 100% accuracy in detecting nonce gaps and ordering violations
- **Resource Efficiency**: Low memory usage and minimal goroutine overhead
- **Hardware Agnostic**: TPS performance is logged for comparison but not asserted
- **Consistent Validation**: Same correctness checks across all test scenarios

Performance will vary based on hardware and configuration. Tests log TPS, latency, and resource usage for comparison across different environments.

## Configuration Options

### Essential Settings
```go
type LoadTestConfig struct {
    // Test Duration and Scale
    Duration        time.Duration  // How long to run the test
    TargetTPS       int           // Target transactions per second
    ClientCount     int           // Number of RPC clients to use
    WalletCount     int           // Number of wallets (nonce sources)
    
    // Transaction Configuration
    TransactionType string        // "dynamicfee", "blob", "setcode", "mixed"
    GasLimit        uint64        // Gas limit per transaction
    BaseFee         uint64        // Base fee in gwei
    TipFee          uint64        // Priority fee in gwei
    Amount          uint64        // Transaction value in wei
    
    // Test Control
    MaxPending      uint64        // Maximum pending transactions
    Timeout         time.Duration // Test timeout
    WarmupDuration  time.Duration // Warmup phase duration
    
    // Validation Flags
    EnableNonceValidation      bool  // Check for nonce gaps
    EnableOrderingValidation   bool  // Check transaction ordering
    EnableThroughputValidation bool  // Strict TPS validation
}
```

### Validation Modes
- **"full"** - All validations enabled (recommended for correctness testing)
- **"minimal"** - Basic validations only (recommended for performance testing)
- **"custom"** - Use individual validation flags for fine control

## Extending the Framework

### Adding New Transaction Types
1. Add transaction type to `LoadTestEngine.generateTransaction()`
2. Implement transaction building logic in `MockWallet`
3. Add test scenarios in `scenarios/` directory

### Adding New Metrics
1. Extend `LoadTestResult` struct in `framework/config.go`
2. Add collection logic in `MetricsCollector`
3. Update report generation in `LoadTestEngine`

### Adding New Validation Rules
1. Extend `TransactionValidator` in `validators/transaction.go`
2. Add validation logic to real-time and final validation phases
3. Add corresponding test cases

## Troubleshooting

### Common Issues

**Tests timing out:**
- Reduce `Duration` or `TargetTPS` in config
- Use `make test-load` (includes 15-minute timeout) or increase test timeout: `go test -timeout 15m`
- Note: Some scaling tests may take several minutes to complete

**Want to reduce logging verbosity:**
- Individual transaction errors are logged at Debug level
- Pending transaction warnings are rate-limited to once per 5 seconds
- Use `logger.SetLevel(logrus.WarnLevel)` to reduce output

**Nonce gaps detected:**
- Usually indicates race conditions in wallet nonce management
- Check wallet initialization and nonce synchronization
- Review concurrent access patterns

**Memory usage growing:**
- Monitor `MaxPending` setting - high values can cause memory buildup
- Check for goroutine leaks in custom test code
- Review transaction cleanup logic

### Debug Logging
Enable detailed logging for troubleshooting:
```go
logger := logrus.NewEntry(logrus.New())
logger.Logger.SetLevel(logrus.DebugLevel)
```

## Future Enhancements

Potential areas for framework expansion:
- **Real Network Testing** - Support for testing against live Ethereum networks
- **Advanced Transaction Patterns** - Contract interactions, multi-sig transactions
- **Distributed Load Testing** - Coordination across multiple test runners
- **Performance Profiling** - Integration with Go profiling tools
- **Custom Metrics** - Plugin system for domain-specific metrics