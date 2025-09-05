# Spamoor Scenario Developer Guide

This guide provides comprehensive documentation for developers who want to implement custom transaction scenarios for Spamoor. Scenarios are the core extensibility mechanism that allows you to define specific transaction patterns and behaviors.


## Table of Contents

- [Scenario Architecture](#scenario-architecture)
- [Scenario Lifecycle](#scenario-lifecycle-overview)
- [Core Interfaces](#core-interfaces)
- [Critical Development Rules](#critical-development-rules)
- [Project Structure](#project-structure)
- [Common Scenario Options](#common-scenario-options)
- [Client Pool Management](#client-pool-management)
- [Wallet Management](#wallet-management)
- [Transaction Building](#transaction-building)
- [Contract Deployments](#contract-deployments)
- [Common Patterns](#common-patterns)
- [RunTransactionScenario Helper](#runtransactionscenario-helper)
- [Best Practices](#best-practices)
- [Testing Scenarios](#testing-scenarios)
- [Example Implementation](#example-implementation)

## Scenario Architecture

### Core Concepts

Spamoor scenarios are self-contained Go packages that implement the `Scenario` interface. Each scenario:

- **Defines transaction patterns**: Specific types and sequences of transactions
- **Manages configuration**: YAML-based configuration with CLI flag support  
- **Controls execution flow**: Rate limiting, concurrency, and completion logic
- **Handles wallet management**: Child wallet selection and funding
- **Provides logging**: Structured logging with transaction tracking

### Scenario Lifecycle Overview

The scenario lifecycle consists of distinct phases that Spamoor manages:

1. **Registration**: Scenario registered in `scenarios/scenarios.go`
2. **Instantiation**: New scenario instance created via factory function
3. **Configuration**: CLI flags registered and parsed
4. **Initialization**: Scenario configured with wallet pool and options
5. **Wallet Funding**: Spamoor prepares and funds wallets automatically
6. **Execution**: Scenario runs transaction generation logic
7. **Cleanup**: Resources cleaned up on completion or cancellation

#### Integration with Web UI

Once registered, your scenario will automatically:
- Appear in the "Create Spammer" dialog (as shown in the [Create Spammer screenshot](./../.github/resources/create-spammer.png))
- Support YAML configuration in the web editor
- Display in the scenario dropdown
- Show custom metrics on the dashboard

### Detailed Lifecycle Phases

#### 1. Registration Phase

Every scenario must be registered in `scenarios/scenarios.go`:

```go
var ScenarioDescriptors = []*scenario.Descriptor{
    &simpleeoa.ScenarioDescriptor,
    &contractcalls.ScenarioDescriptor,
    &blobscenario.ScenarioDescriptor,
    // ... your scenario here
}
```

The descriptor provides metadata and a factory function:

```go
var ScenarioDescriptor = scenario.Descriptor{
    Name:           "my-scenario",
    Description:    "Description of what this scenario does",
    DefaultOptions: ScenarioDefaultOptions,
    NewScenario:    newScenario,  // Factory function
}
```

#### 2. Instantiation and Configuration

When a user runs a scenario, Spamoor:

1. **Creates a new instance** using the factory function:
```go
scenarioInstance := descriptor.NewScenario(logger)
```

2. **Registers CLI flags** by calling `Flags()`:
```go
func (s *Scenario) Flags(flags *pflag.FlagSet) error {
    flags.Uint64VarP(&s.options.TotalCount, "count", "c", 0, "Total transactions")
    flags.Uint64VarP(&s.options.Throughput, "throughput", "t", 10, "Transactions per slot")
    // Register all scenario-specific flags
    return nil
}
```

3. **Parses command-line arguments** to populate the options

#### 3. Initialization Phase

After configuration, Spamoor calls `Init()` with runtime options:

```go
func (s *Scenario) Init(options *scenario.Options) error {
    // 1. Store wallet pool reference
    s.walletPool = options.WalletPool
    
    // 2. Parse YAML config if provided
    if options.Config != "" {
        err := yaml.Unmarshal([]byte(options.Config), &s.options)
        if err != nil {
            return fmt.Errorf("failed to unmarshal config: %w", err)
        }
    }
    
    // 3. Configure wallet pool
    s.walletPool.SetWalletCount(100)  // Number of child wallets
    s.walletPool.SetRefillAmount(utils.EtherToWei(uint256.NewInt(5)))   // 5 ETH per refill
    s.walletPool.SetRefillBalance(utils.EtherToWei(uint256.NewInt(1)))  // Refill when < 1 ETH
    
    // 4. Add well-known wallets if needed
    s.walletPool.AddWellKnownWallet(&spamoor.WellKnownWalletConfig{
        Name:          "deployer",
        RefillAmount:  utils.EtherToWei(uint256.NewInt(50)),
        RefillBalance: utils.EtherToWei(uint256.NewInt(10)),
    })
    
    // 5. Validate configuration
    if s.options.TotalCount == 0 && s.options.Throughput == 0 {
        return fmt.Errorf("either total_count or throughput must be specified")
    }
    
    return nil
}
```

**Important**: During `Init()`, you MUST configure all wallets that your scenario will use. This includes:
- Setting the number of child wallets via `SetWalletCount()`
- Adding any well-known wallets via `AddWellKnownWallet()`
- Configuring refill amounts and thresholds

#### 4. Wallet Funding Phase (Automatic)

**After `Init()` completes and before `Run()` is called**, Spamoor automatically:

1. **Creates all configured wallets** based on your settings
2. **Funds wallets** that are below the refill threshold
3. **Waits for funding transactions** to be confirmed
4. **Starts background refill monitoring** to maintain balances

This automatic funding ensures all wallets have sufficient ETH before your scenario begins executing transactions.

#### 5. Execution Phase

Spamoor calls `Run()` with a cancellable context:

```go
func (s *Scenario) Run(ctx context.Context) error {
    s.logger.Infof("starting scenario: %s", s.options)
    
    // Most scenarios use the RunTransactionScenario helper
    return scenario.RunTransactionScenario(ctx, scenario.TransactionScenarioOptions{
        TotalCount:      s.options.TotalCount,
        Throughput:      s.options.Throughput,
        MaxPending:      s.options.MaxPending,
        Timeout:         timeout,
        WalletPool:      s.walletPool,
        Logger:          s.logger,
        ProcessNextTxFn: s.sendNextTransaction,
    })
}
```

The `Run()` method should:
- Respect context cancellation for graceful shutdown
- Execute the scenario's transaction logic
- Return when complete or cancelled
- Return an error if the scenario fails

#### 6. Context Cancellation and Cleanup

Scenarios must handle context cancellation properly:

```go
select {
case <-ctx.Done():
    s.logger.Info("scenario cancelled")
    return ctx.Err()
default:
    // Continue processing
}
```

When the context is cancelled (user presses Ctrl+C or timeout reached):
- Stop all transaction generation
- Allow pending transactions to complete if possible
- Clean up any resources
- Return promptly

## Core Interfaces

### Scenario Interface

Every scenario must implement this interface:

```go
type Scenario interface {
    // Flags registers CLI flags for the scenario
    Flags(flags *pflag.FlagSet) error
    
    // Init initializes the scenario with runtime options
    Init(options *Options) error
    
    // Run executes the scenario until completion or cancellation
    Run(ctx context.Context) error
}
```

### Scenario Descriptor

Each scenario provides metadata via a descriptor:

```go
type Descriptor struct {
    Name           string      // Unique scenario name
    Description    string      // Human-readable description
    DefaultOptions any         // Default configuration struct
    NewScenario    func(logger logrus.FieldLogger) Scenario
}
```

### Initialization Options

Scenarios receive these options during initialization:

```go
type Options struct {
    WalletPool *spamoor.WalletPool  // Managed wallet pool
    Config     string               // YAML configuration
    GlobalCfg  map[string]any       // Global daemon configuration
}
```

## Project Structure

### Directory Layout

```
scenarios/
├── my-scenario/              # Your scenario package
│   ├── my_scenario.go       # Main implementation
│   ├── README.md            # Scenario documentation
│   └── contract/            # Contract bindings (if needed)
│       ├── MyContract.go
│       ├── MyContract.sol
│       └── compile.sh
└── scenarios.go             # Registration file
```

### Registration

Add your scenario to `scenarios/scenarios.go`:

```go
import "github.com/ethpandaops/spamoor/scenarios/my-scenario"

var ScenarioDescriptors = []*scenario.Descriptor{
    // ... existing scenarios
    &myscenario.ScenarioDescriptor,
}
```

## Critical Development Rules

🚨 **NEVER use the root wallet directly in scenarios** - it's shared across all running scenarios and direct usage will cause nonce conflicts.

🚨 **NEVER use go-ethereum's bound contracts for transactions** - they manage nonces independently and will conflict with Spamoor's nonce tracking.

🚨 **ALWAYS spread transactions across multiple wallets** - Ethereum clients have limits on pending transactions per sender (typically 64-1000), so high-throughput scenarios must use multiple wallets.

🚨 **ALWAYS respect context cancellation** - scenarios must stop all operations when context is cancelled.

🚨 **ALWAYS call onComplete() in ProcessNextTxFn** - required for scenario.RunTransactionScenario transaction counting.

🚨 **NEVER assume receipt is non-nil in OnComplete** - handle cancellation and replacement transaction scenarios properly.

## Common Scenario Options

All scenarios should support a standard set of options to ensure consistent behavior across the tool. These options control transaction execution, rate limiting, and resource management.

### Standard Options

```go
type ScenarioOptions struct {
    // Transaction Control
    TotalCount  uint64 `yaml:"total_count"`  // Total number of transactions to send (0 = unlimited)
    Throughput  uint64 `yaml:"throughput"`   // Transactions per slot (12 seconds on mainnet)
    MaxPending  uint64 `yaml:"max_pending"`  // Maximum concurrent pending transactions
    MaxWallets  uint64 `yaml:"max_wallets"`  // Maximum number of child wallets to use
    
    // Gas Configuration
    BaseFee     uint64 `yaml:"base_fee"`     // Base fee in gwei for EIP-1559 transactions
    TipFee      uint64 `yaml:"tip_fee"`      // Priority fee (tip) in gwei
    GasLimit    uint64 `yaml:"gas_limit"`    // Gas limit for transactions
    
    // Client Control
    ClientGroup string `yaml:"client_group"` // Preferred client group (e.g., "validators", "archive")
    
    // Logging
    LogTxs      bool   `yaml:"log_txs"`      // Log individual transactions (vs just summary)
    
    // Scenario-specific options...
}
```

### Option Descriptions

#### Transaction Control Options

**`total_count`** (default: 0)
- Total number of transactions to send before scenario completes
- Set to 0 for unlimited (scenario runs until cancelled or timeout)
- Takes precedence over timeout if both are specified

**`throughput`** (default: varies by scenario)
- Target transactions per slot (12 seconds on mainnet, 4 seconds on some testnets)
- Set to 0 for maximum speed (limited only by max_pending)
- Automatically adjusted based on network congestion

**`max_pending`** (default: 0)
- Maximum number of transactions pending confirmation at any time
- Set to 0 for no limit (use with caution)
- Helps prevent overwhelming clients with too many pending transactions

**`max_wallets`** (default: 0)
- Maximum number of child wallets to create and use
- Set to 0 to let scenario determine based on throughput
- More wallets allow higher throughput but increase funding costs

#### Gas Configuration Options

**`base_fee`** (default: 20 gwei)
- Base fee multiplier for EIP-1559 transactions
- Actual fee = network base fee × multiplier
- Higher values increase transaction priority

**`tip_fee`** (default: 2 gwei)
- Priority fee (miner tip) for EIP-1559 transactions
- Added to base fee for total gas price
- Higher tips can improve inclusion speed

**`gas_limit`** (default: varies by transaction type)
- Gas limit for each transaction
- Must be sufficient for transaction type
- Common values: 21000 (simple transfer), 100000+ (contract calls)

#### Client Control Options

**`client_group`** (default: "")
- Preferred group of RPC clients to use
- Common groups: "validators", "archive", "local"
- Empty string uses any available client

#### Logging Options

**`log_txs`** (default: false)
- When true, logs each individual transaction
- When false, only logs summary statistics
- Useful for debugging but can be verbose

### Implementation Example

```go
func (s *Scenario) Flags(flags *pflag.FlagSet) error {
    flags.Uint64VarP(&s.options.TotalCount, "count", "c", 0, 
        "Total number of transactions to send (0 = unlimited)")
    flags.Uint64VarP(&s.options.Throughput, "throughput", "t", 10, 
        "Transactions per slot")
    flags.Uint64Var(&s.options.MaxPending, "max-pending", 0, 
        "Maximum pending transactions (0 = no limit)")
    flags.Uint64Var(&s.options.MaxWallets, "max-wallets", 0, 
        "Maximum number of wallets (0 = auto)")
    flags.Uint64Var(&s.options.BaseFee, "basefee", 20, 
        "Max base fee in gwei")
    flags.Uint64Var(&s.options.TipFee, "tipfee", 2, 
        "Max tip fee in gwei")
    flags.Uint64Var(&s.options.GasLimit, "gaslimit", 50000, 
        "Gas limit for transactions")
    flags.StringVar(&s.options.ClientGroup, "clientgroup", "", 
        "Client group to use for transactions")
    flags.BoolVar(&s.options.LogTxs, "log-txs", false, 
        "Log all transactions")
    
    // Scenario-specific flags...
    
    return nil
}
```

### Configuration File Example

Scenarios can also be configured via YAML:

```yaml
# Standard options
total_count: 1000
throughput: 50
max_pending: 100
max_wallets: 20

# Gas configuration
base_fee: 30
tip_fee: 3
gas_limit: 100000

# Execution control
timeout: "30m"
client_group: "validators"

# Logging
log_txs: true

# Scenario-specific options
amount: 1000000000000000000  # 1 ETH in wei
contract_address: "0x..."
```

### Best Practices

1. **Always support the standard options** - Users expect consistent behavior
2. **Provide sensible defaults** - Most users won't customize every option
3. **Validate option combinations** - Warn about conflicting settings
4. **Document scenario-specific options** - Explain any custom parameters
5. **Use standard flag names** - Maintain consistency with other scenarios

## Client Pool Management

The client pool is a critical component that manages RPC endpoint connections and ensures reliable transaction submission. Understanding how to properly use the client pool is essential for building robust scenarios.

### Overview

The client pool provides:
- **Health monitoring**: Continuously checks that clients are alive and following the chain head
- **Client groups**: Logical grouping of clients by capability (validators, archive nodes, etc.)
- **Load distribution**: Spreads transactions across multiple endpoints
- **Automatic failover**: Routes around unhealthy clients

### Client Health Monitoring

Spamoor automatically monitors all configured RPC clients:

1. **Liveness checks**: Regular polling to ensure clients respond to requests
2. **Chain head tracking**: Verifies clients are synced and following the canonical chain
3. **Automatic marking**: Unhealthy clients are marked as unavailable
4. **Recovery detection**: Previously unhealthy clients are re-enabled when they recover

### Client Selection

#### Basic Client Selection

```go
// Get any available healthy client
client := s.walletPool.GetClient(spamoor.SelectClientRandom, 0, "")

// Get client by index (round-robin)
client := s.walletPool.GetClient(spamoor.SelectClientByIndex, txIdx, "")

// Get client from specific group
client := s.walletPool.GetClient(spamoor.SelectClientRandom, 0, "validators")
```

#### Selection Strategies

**`SelectClientRandom`**
- Randomly selects from available healthy clients
- Best for: Even load distribution
- Use when: Order doesn't matter

**`SelectClientByIndex`**
- Deterministic selection using modulo of index
- Best for: Reproducible client assignment
- Use when: Debugging or testing specific endpoints

**`SelectClientRoundRobin`**
- Sequential selection with automatic wraparound
- Best for: Fair distribution across all clients
- Use when: Want to ensure all clients get traffic

#### Client Groups

Client groups allow filtering by capability or purpose:

```go
// Common client groups
"validators"    // Validator nodes (may have mempool restrictions)
"archive"       // Archive nodes with full history
"local"         // Local development nodes
"light"         // Light clients
""              // Any available client (default)
```

Configure groups when starting Spamoor:
```bash
./spamoor my-scenario \
  --rpchost "http://validator1:8545#group=validators" \
  --rpchost "http://archive1:8545#group=archive" \
  --rpchost "http://local:8545#group=local"
```

### Transaction Distribution Strategies

#### Single Transaction Submission

For individual transactions, distribute across clients:

```go
func (s *Scenario) sendNextTransaction(ctx context.Context, txIdx uint64, onComplete func()) (func(), error) {
    // Rotate clients for each transaction
    client := s.walletPool.GetClient(spamoor.SelectClientRoundRobin, txIdx, s.options.ClientGroup)
    
    // Build and submit transaction
    err := s.walletPool.GetTxPool().SendTransaction(ctx, wallet, signedTx, &spamoor.SendTransactionOptions{
        Client: client,  // Use selected client
        OnComplete: onComplete,
    })
}
```

#### Batch Submission Considerations

When submitting large batches of transactions, be aware of transaction ordering requirements:

**Problem**: Different clients may have gaps in their view of pending transactions, leading to nonce gaps if transactions are distributed across multiple clients.

**Solution 1 - Single Client for Batches**:
```go
// Get a single client for the entire batch
client := s.walletPool.GetClient(spamoor.SelectClientRandom, 0, s.options.ClientGroup)

// Submit all batch transactions to the same client
receipts, err := s.walletPool.GetTxPool().SendTransactionBatch(ctx, wallet, transactions, &spamoor.BatchOptions{
    SendTransactionOptions: spamoor.SendTransactionOptions{
        Client: client,  // Same client for all transactions
    },
    PendingLimit: 50,
})
```

**Solution 2 - Use Multi-Wallet Batch Submission**:
```go
// When using multiple wallets, transactions can be distributed
// because each wallet maintains its own nonce sequence
walletTxs := make(map[*spamoor.Wallet][]*types.Transaction)

// Group transactions by wallet
for i, tx := range transactions {
    wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, i)
    walletTxs[wallet] = append(walletTxs[wallet], tx)
}

// Submit using multi-wallet batch - internally handles client distribution
receipts, err := s.walletPool.GetTxPool().SendMultiTransactionBatch(ctx, walletTxs, &spamoor.BatchOptions{
    ClientGroup: s.options.ClientGroup,
    PendingLimit: 50,
})
```

### Best Practices

1. **Spread individual transactions across clients** - Maximizes throughput and resilience
2. **Use client groups** - Target appropriate client types for your scenario
3. **Keep transaction batches on single clients** - Avoids nonce gap issues
4. **Monitor client health in logs** - Watch for warnings about unhealthy clients
5. **Use multi-wallet batching for mass submission** - Allows safe client distribution

### Troubleshooting Client Issues

Common issues and solutions:

**"No healthy clients available"**
- Check that RPC endpoints are accessible
- Verify clients are synced to chain head
- Look for connection errors in logs

**"Transaction nonce gaps"**
- Ensure batch transactions use single client
- Consider using multi-wallet batching
- Check for client mempool limits

**"Uneven client load"**
- Use round-robin selection instead of random
- Verify all clients are marked healthy
- Check client group configuration

## Wallet Management

### Wallet Architecture Overview

Spamoor implements a sophisticated wallet management system that provides isolation, automation, and flexibility for transaction scenarios.

#### Scenario-Specific Wallet Pools

Each scenario has its **own unique set of wallets** derived from the root private key and a scenario-specific seed:

- **Deterministic derivation**: Wallets are generated by hashing the root private key with the scenario seed and wallet identifier
- **Scenario isolation**: Each scenario's wallets are completely separate from other scenarios
- **Reproducible addresses**: Same scenario with same seed always generates identical wallet addresses
- **Conflict prevention**: Scenario seeds ensure no address conflicts between different running scenarios

#### Wallet Types

Spamoor supports two types of wallets for different use cases:

**Numbered Wallets (Child Wallets)**:
- **Purpose**: Mass operations that distribute transactions across multiple source wallets
- **Selection**: Accessed by index using selection strategies (random, round-robin, by-index, by-pending-count)
- **Scaling**: Pool size configurable based on transaction volume needs
- **Use cases**: High-throughput transaction sending, load distribution across multiple senders

**Named Wallets (Well-Known Wallets)**:
- **Purpose**: Special use cases requiring consistent addresses (deployments, admin roles, etc.)
- **Selection**: Accessed by name using `GetWellKnownWallet("wallet-name")`
- **Consistency**: Same name always returns the same wallet address within a scenario
- **Use cases**: Contract deployments, token ownership, admin operations, test scenarios

#### Automatic Wallet Funding

All wallet types are automatically funded and managed:

- **Continuous monitoring**: Background service monitors all wallet balances
- **Automatic refills**: Wallets are automatically funded when balance drops below threshold
- **Configurable amounts**: Both refill amount and threshold are configurable per scenario
- **Batch operations**: Funding transactions are batched for efficiency
- **Root wallet management**: Uses thread-safe locking for root wallet access

#### Fund Recovery

After scenario execution, leftover funds can be reclaimed to the root wallet.

### Child Wallets and Selection Strategies

Child wallets are derived using deterministic key derivation and provide the backbone for high-throughput transaction scenarios.

#### Wallet Selection Modes

```go
// SelectWalletByIndex - Deterministic selection by index (modulo pool size)
// Use for: Predictable wallet assignment, testing specific wallets
wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, txIdx)

// SelectWalletRandom - Random wallet selection
// Use for: Even distribution when order doesn't matter
wallet := s.walletPool.GetWallet(spamoor.SelectWalletRandom, 0)

// SelectWalletRoundRobin - Sequential round-robin selection
// Use for: Even distribution across all wallets in order
wallet := s.walletPool.GetWallet(spamoor.SelectWalletRoundRobin, 0)

// SelectWalletByPendingTxCount - Select wallet with lowest pending transactions
// Use for: Optimal distribution when client limits are a concern
wallet := s.walletPool.GetWallet(spamoor.SelectWalletByPendingTxCount, txIdx)
```

#### Client Pending Transaction Limits

Ethereum clients limit pending transactions per sender address:
- **Geth**: ~64 pending transactions per account by default
- **Reth**: ~1000 pending transactions per account  
- **Nethermind**: ~1024 pending transactions per account
- **Besu**: ~64 pending transactions per account by default

For high-throughput scenarios (>50 tx/slot), always use multiple wallets:

```go
// Configure adequate wallet count based on throughput
func (s *Scenario) Init(options *scenario.Options) error {
    // Rule of thumb: 1 wallet per 50 transactions for safety margins
    if s.options.TotalCount > 0 {
		maxWallets := s.options.TotalCount / 50
		if maxWallets < 10 {
			maxWallets = 10
		} else if maxWallets > 1000 {
			maxWallets = 1000
		}

		s.walletPool.SetWalletCount(maxWallets)
	} else {
		if s.options.Throughput*2 < 1000 {
			s.walletPool.SetWalletCount(s.options.Throughput * 2)
		} else {
			s.walletPool.SetWalletCount(1000)
		}
	}
}

// Use pending-count-based selection for optimal distribution
wallet := s.walletPool.GetWallet(spamoor.SelectWalletByPendingTxCount, txIdx)
```

### Well-Known Wallets

Well-known wallets have deterministic addresses and are used for contracts or roles that need consistent addresses across scenario runs.

#### Types of Well-Known Wallets

**Regular Well-Known Wallets** (scenario-specific):
- Scoped to the current scenario instance
- Addresses change between different scenario runs
- Used for scenario-specific contracts and deployments

**Very Well-Known Wallets** (application-wide):
- Consistent addresses across ALL scenarios and instances
- Used for shared infrastructure contracts
- Examples: CREATE2 factories, shared registries, common utilities

#### Configuring Well-Known Wallets

```go
func (s *Scenario) Init(options *scenario.Options) error {
    // Regular well-known wallet (scenario-specific)
    s.walletPool.AddWellKnownWallet(&spamoor.WellKnownWalletConfig{
        Name:          "deployer",
        RefillAmount:  utils.EtherToWei(uint256.NewInt(100)), // 100 ETH
        RefillBalance: utils.EtherToWei(uint256.NewInt(50)),  // 50 ETH
        VeryWellKnown: false, // Default: scenario-specific
    })
    
    // Very well-known wallet (application-wide)
    s.walletPool.AddWellKnownWallet(&spamoor.WellKnownWalletConfig{
        Name:          "create2-factory-deployer",
        RefillAmount:  utils.EtherToWei(uint256.NewInt(10)),  // 10 ETH
        RefillBalance: utils.EtherToWei(uint256.NewInt(5)),   // 5 ETH
        VeryWellKnown: true, // Same address across all scenarios
    })
    
    // Multi-scenario shared wallet (for shared infrastructure)
    s.walletPool.AddWellKnownWallet(&spamoor.WellKnownWalletConfig{
        Name:          "registry-owner",
        RefillAmount:  utils.EtherToWei(uint256.NewInt(50)),
        RefillBalance: utils.EtherToWei(uint256.NewInt(25)),
        VeryWellKnown: true,
    })
}

// Usage in scenarios
func (s *Scenario) deployContracts() error {
    // Use scenario-specific deployer for scenario contracts
    deployerWallet := s.walletPool.GetWellKnownWallet("deployer")
    
    // Use very well-known wallet for shared infrastructure
    factoryDeployer := s.walletPool.GetWellKnownWallet("create2-factory-deployer")
    
    // These addresses will be the same across all scenario instances
    if factoryDeployer != nil {
        s.logger.Infof("CREATE2 factory deployer: %s", factoryDeployer.GetAddress().Hex())
    }
}
```

#### Use Cases for Well-Known Wallets

**Regular Well-Known Wallets**:
- **Contract deployers**: Deploy scenario-specific contracts
- **Token owners**: Own tokens created for the scenario
- **Admin roles**: Manage scenario-specific permissions
- **Test accounts**: Specific roles in testing scenarios

**Very Well-Known Wallets**:
- **CREATE2 factories**: Deploy deterministic contract addresses
- **Shared registries**: Cross-scenario contract registries
- **Common utilities**: Shared helper contracts
- **Infrastructure**: Multi-scenario dependencies

### Automatic Wallet Funding

Spamoor automatically manages wallet funding in the background through a continuous monitoring system.

#### Funding Configuration

```go
// Configure wallet pool funding parameters
func (s *Scenario) Init(options *scenario.Options) error {
    // Set wallet count
    s.walletPool.SetWalletCount(100)
    
    // Configure funding amounts (in wei)
    s.walletPool.SetRefillAmount(utils.EtherToWei(uint256.NewInt(5)))   // 5 ETH per refill
    s.walletPool.SetRefillBalance(utils.EtherToWei(uint256.NewInt(2)))  // Refill when below 2 ETH
    s.walletPool.SetRefillInterval(300) // Check balances every 300 seconds (5 minutes)
}
```

#### Funding Process

The funding system operates continuously:

1. **Balance Monitoring**: Every `RefillInterval` seconds, checks all wallet balances
2. **Funding Detection**: Identifies wallets below `RefillBalance` threshold  
3. **Batch Funding**: Groups funding transactions for efficiency
4. **Root Wallet Locking**: Uses root wallet locking mechanism for thread safety
5. **Transaction Batching**: Submits funding transactions in optimal batches

#### Root Wallet Locking

For scenarios that need large funding transactions (like providing liquidity), use the root wallet locking mechanism:

```go
// Lock root wallet for exclusive use
rootWallet := s.walletPool.GetRootWallet()
clientPool := s.walletPool.GetClientPool()
err := rootWallet.WithWalletLock(ctx, expectedTxCount, expectedTotalAmount, clientPool, func(reason string) {
    s.logger.Infof("Root wallet is locked, %s", reason)
}, func() error {
    // Perform large funding operations here
    tx, err := rootWallet.GetWallet().BuildBoundTx(ctx, &txbuilder.TxMetadata{
        GasFeeCap: uint256.MustFromBig(feeCap),
        GasTipCap: uint256.MustFromBig(tipCap),
        Gas:       6000000,
        Value:     utils.EtherToWei(uint256.NewInt(1000)), // Large amount
    }, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
        return contract.ProvideLiquidity(transactOpts, params...)
    })
    
    if err != nil {
        return err
    }
    
    // Submit transaction
    return s.walletPool.GetTxPool().SendAndAwaitTransaction(ctx, rootWallet.GetWallet(), tx, nil)
})
```

#### Funding Best Practices

```go
// Conservative funding for standard scenarios
s.walletPool.SetRefillAmount(utils.EtherToWei(uint256.NewInt(5)))   // 5 ETH
s.walletPool.SetRefillBalance(utils.EtherToWei(uint256.NewInt(1)))  // 1 ETH threshold

// Aggressive funding for high-throughput scenarios  
s.walletPool.SetRefillAmount(utils.EtherToWei(uint256.NewInt(20)))  // 20 ETH
s.walletPool.SetRefillBalance(utils.EtherToWei(uint256.NewInt(10))) // 10 ETH threshold
s.walletPool.SetRefillInterval(120) // Check every 2 minutes

// Funding for contract-heavy scenarios
s.walletPool.SetRefillAmount(utils.EtherToWei(uint256.NewInt(10)))  // 10 ETH
s.walletPool.SetRefillBalance(utils.EtherToWei(uint256.NewInt(5))) // 5 ETH threshold
```

#### Monitoring Wallet Health

```go
// Check wallet funding status in scenario logs
func (s *Scenario) logWalletStats() {
    walletCount := s.walletPool.GetConfiguredWalletCount()
    s.logger.Infof("Configured wallets: %d", walletCount)
    
    // Wallet funding is logged automatically by the wallet pool
    // Look for logs like:
    // - "funded X wallets with Y ETH total"
    // - "wallet funding completed in Z seconds"
    // - "low balance detected on N wallets"
}
```

The funding system ensures wallets always have sufficient ETH for transactions without manual intervention, allowing scenarios to focus on transaction generation rather than balance management.

## Transaction Building

### Using TxBuilder

Always use the `txbuilder` package for transaction construction:

```go
import "github.com/ethpandaops/spamoor/txbuilder"

// Build transaction metadata
txData := &txbuilder.TxMetadata{
    GasTipCap: uint256.NewInt(tipCap),
    GasFeeCap: uint256.NewInt(feeCap),
    Gas:       gasLimit,
    To:        &targetAddr,
    Value:     uint256.NewInt(amount),
    Data:      callData,
}

// Create different transaction types
dynFeeTx, err := txbuilder.DynFeeTx(txData)  // EIP-1559
blobTx, err := txbuilder.BlobTx(txData, blobs) // EIP-4844
setCodeTx, err := txbuilder.SetCodeTx(txData)  // EIP-7702
```

### Contract Interactions

Use `BuildBoundTx` with abigen-generated contracts for proper nonce management:

```go
// Create contract instance (for ABI, not transactions)
testToken, err := contract.NewContract(contractAddr, client.GetEthClient())
if err != nil {
    return err
}

// Use BuildBoundTx with a function that calls the contract method
tx, err := wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
    GasTipCap: uint256.MustFromBig(tipCap),
    GasFeeCap: uint256.MustFromBig(feeCap),
    Gas:       gasLimit,
    Value:     uint256.NewInt(0),
}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
    // Use the generated contract method directly
    return testToken.Transfer(transactOpts, toAddr, amount)
})

// For contract deployment
tx, err := wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
    GasTipCap: uint256.MustFromBig(tipCap),
    GasFeeCap: uint256.MustFromBig(feeCap),
    Gas:       2000000,
    Value:     uint256.NewInt(0),
}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
    _, deployTx, _, err := contract.DeployContract(transactOpts, client.GetEthClient())
    return deployTx, err
})
```

This pattern:
- Uses abigen-generated contract bindings safely
- Ensures Spamoor manages nonces correctly  
- Provides proper transaction metadata control
- Works for both contract calls and deployments

## Transaction Submission

Spamoor provides multiple transaction submission methods with different levels of control and functionality. All functions are available through the `TxPool` interface.

### Core Submission Functions

#### SendTransaction - Asynchronous Submission

The primary method for asynchronous transaction processing. Returns immediately after submission without waiting for confirmation:

```go
txpool := s.walletPool.GetTxPool()

err := txpool.SendTransaction(ctx, wallet, signedTx, &spamoor.SendTransactionOptions{
    // Client selection options
    Client:             client,           // Use specific client (optional)
    ClientGroup:        "validators",     // Prefer client group (optional)
    ClientsStartOffset: 0,               // Offset for client selection
    
    // Rebroadcast configuration
    Rebroadcast: true, // Enable exponential backoff rebroadcast
    
    // Callback functions
    OnConfirm: func(tx *types.Transaction, receipt *types.Receipt) {
        // Called only when transaction is confirmed in a block
        // receipt is guaranteed to be non-nil here
        s.logger.Infof("Transaction confirmed: %s (block %d)", 
            tx.Hash().Hex(), receipt.BlockNumber.Uint64())
        
        // Handle successful confirmation
        s.handleConfirmedTransaction(tx, receipt)
    },
    OnComplete: func(tx *types.Transaction, receipt *types.Receipt, err error) {
        // Always called when processing completes (success or failure)
        // CRITICAL: receipt may be nil if scenario cancelled or replacement confirmed
        if err != nil {
            s.logger.Warnf("Transaction failed: %v", err)
            return
        }
        
        if receipt == nil {
            // Scenario cancelled or replacement transaction confirmed
            s.logger.Debugf("Transaction cancelled or replaced: %s", tx.Hash().Hex())
            return
        }
        
        if receipt.Status == types.ReceiptStatusSuccessful {
            s.logger.Infof("Transaction completed successfully: %s", tx.Hash().Hex())
        } else {
            s.logger.Warnf("Transaction reverted: %s", tx.Hash().Hex())
        }
        
        // Call scenario completion callback
        onComplete()
    },
    OnEncode: func(tx *types.Transaction) ([]byte, error) {
        // Custom transaction encoding (optional)
        // Useful for alternative serialization formats
        return tx.MarshalBinary()
    },
    LogFn: spamoor.GetDefaultLogFn(s.logger, "transfer", fmt.Sprintf("%d", txIdx), signedTx),
})
```

#### SendAndAwaitTransaction - Synchronous Submission

Submits a transaction and waits for confirmation. Returns the receipt or error:

```go
receipt, err := txpool.SendAndAwaitTransaction(ctx, wallet, signedTx, &spamoor.SendTransactionOptions{
    Client:      client,
    Rebroadcast: true, // Recommended for important transactions
    OnConfirm: func(tx *types.Transaction, receipt *types.Receipt) {
        // Called when confirmed (tx confirmation go subroutine)
        s.logger.Infof("Deployment confirmed: %s", receipt.ContractAddress.Hex())
    },
})

if err != nil {
    return fmt.Errorf("deployment failed: %w", err)
}

if receipt == nil {
    return fmt.Errorf("deployment cancelled or replaced")
}

if receipt.Status == types.ReceiptStatusSuccessful {
    contractAddr := receipt.ContractAddress
    s.logger.Infof("Contract deployed at: %s", contractAddr.Hex())
} else {
    return fmt.Errorf("deployment transaction reverted")
}
```

#### SendTransactionBatch - Single-Wallet Batch Submission

Efficiently submits multiple transactions from a single wallet with concurrency control and retry logic:

**Important**: All transactions in a batch must originate from the same sender wallet. For multiple wallets, use `SendMultiTransactionBatch` instead.

```go
// Prepare batch of transactions from the same wallet
deploymentTxs := []*types.Transaction{tx1, tx2, tx3}

receipts, err := txpool.SendTransactionBatch(ctx, wallet, deploymentTxs, &spamoor.BatchOptions{
    SendTransactionOptions: spamoor.SendTransactionOptions{
        Client:      client,
        Rebroadcast: true,
        OnConfirm: func(tx *types.Transaction, receipt *types.Receipt) {
            // Called for each confirmed transaction
            s.logger.Infof("Batch transaction confirmed: %s", tx.Hash().Hex())
        },
    },
    PendingLimit: 50, // Limit concurrent pending transactions
})

if err != nil {
    return fmt.Errorf("batch deployment failed: %w", err)
}

// Process results - receipts array matches input transaction order
for i, receipt := range receipts {
    if receipt != nil && receipt.Status == types.ReceiptStatusSuccessful {
        s.logger.Infof("Transaction %d confirmed in block %d", i, receipt.BlockNumber.Uint64())
    } else if receipt != nil {
        s.logger.Warnf("Transaction %d reverted", i)
    } else {
        s.logger.Warnf("Transaction %d failed or cancelled", i)
    }
}
```

#### SendMultiTransactionBatch - Multi-Wallet Batch Submission

Efficiently submits multiple transactions across multiple wallets with advanced concurrency control and retry logic:

```go
// Prepare transactions grouped by wallet
walletTxs := make(map[*spamoor.Wallet][]*types.Transaction)
for _, tx := range allTransactions {
    wallet := s.getWalletForTransaction(tx)
    walletTxs[wallet] = append(walletTxs[wallet], tx)
}

// Submit all transactions across all wallets simultaneously
receipts, err := txpool.SendMultiTransactionBatch(ctx, walletTxs, &spamoor.BatchOptions{
    SendTransactionOptions: spamoor.SendTransactionOptions{
        Client:      client,
        Rebroadcast: true,
        OnConfirm: func(tx *types.Transaction, receipt *types.Receipt) {
            // Called for each confirmed transaction across all wallets
            s.logger.Infof("Multi-batch transaction confirmed: %s", tx.Hash().Hex())
        },
    },
    PendingLimit: 50,    // Maximum pending transactions per wallet
    MaxRetries:   3,     // Retry failed submissions up to 3 times
    ClientPool:   clientPool, // Optional: assign different clients to different wallets
    LogFn: func(confirmedCount int, totalCount int) {
        // Progress logging callback
        s.logger.Infof("Multi-batch progress: %d/%d transactions confirmed", confirmedCount, totalCount)
    },
    LogInterval: 10,     // Call LogFn every 10 confirmed transactions
})

if err != nil {
    return fmt.Errorf("multi-wallet batch failed: %w", err)
}

// Process results - receipts map matches input wallet structure
for wallet, walletReceipts := range receipts {
    for i, receipt := range walletReceipts {
        if receipt != nil && receipt.Status == types.ReceiptStatusSuccessful {
            s.logger.Infof("Wallet %s transaction %d confirmed in block %d", 
                wallet.GetAddress().Hex(), i, receipt.BlockNumber.Uint64())
        } else if receipt != nil {
            s.logger.Warnf("Wallet %s transaction %d reverted", wallet.GetAddress().Hex(), i)
        } else {
            s.logger.Warnf("Wallet %s transaction %d failed or cancelled", wallet.GetAddress().Hex(), i)
        }
    }
}
```

**Advanced Features**:

- **Per-wallet concurrency control**: `PendingLimit` applies to each wallet individually
- **Automatic retry logic**: Failed submissions are retried up to `MaxRetries` times
- **Client pool assignment**: Optionally assign different RPC clients to different wallets for load distribution
- **Progress tracking**: `LogFn` provides real-time progress updates across all wallets
- **Sliding window submission**: Maintains optimal throughput while respecting limits

**Single-Wallet Batch (Alternative Approach)**:

```go
// Alternative: Process wallets individually with separate batch calls
var allReceipts []*types.Receipt
for wallet, txs := range walletTxs {
    receipts, err := txpool.SendTransactionBatch(ctx, wallet, txs, &spamoor.BatchOptions{
        SendTransactionOptions: spamoor.SendTransactionOptions{
            Rebroadcast: true,
        },
        PendingLimit: 50,
    })
    if err != nil {
        return fmt.Errorf("batch failed for wallet %s: %w", wallet.GetAddress().Hex(), err)
    }
    allReceipts = append(allReceipts, receipts...)
}
```

#### AwaitTransaction - Manual Confirmation Waiting

Wait for an already-submitted transaction to confirm:

```go
// For transactions submitted elsewhere
receipt, err := txpool.AwaitTransaction(ctx, wallet, tx)
if err != nil {
    return fmt.Errorf("transaction confirmation failed: %w", err)
}

if receipt.Status == types.ReceiptStatusSuccessful {
    s.logger.Infof("Transaction confirmed: %s", tx.Hash().Hex())
}
```

### Submission Lifecycle and Handling

#### Transaction States

1. **Submission**: Transaction sent to RPC endpoint
2. **Pending**: Transaction in mempool awaiting inclusion
3. **Confirmation**: Transaction included in a block
4. **Completion**: Final state (confirmed, failed, or cancelled)

#### Rebroadcast Mechanism

Spamoor includes an automatic rebroadcast system with exponential backoff to handle stuck transactions:

```go
SendTransactionOptions{
    Rebroadcast: true, // Enable automatic rebroadcast
}
```

**Rebroadcast Features**:
- **Automatic detection**: Monitors for transactions stuck in mempool
- **Exponential backoff**: Increases wait time between rebroadcast attempts (2s, 4s, 8s, 16s, ...)
- **Replacement handling**: Detects when replacement transactions are confirmed
- **Context cancellation**: Stops rebroadcasting when scenario is cancelled
- **Logging**: Provides detailed logs for debugging stuck transactions

#### Confirmation Handling Best Practices

**Critical Notes**:
- `OnConfirm` callback only called when transaction is successfully included
- `OnComplete` callback may receive `nil` receipt in these cases:
  - Scenario was cancelled before confirmation
  - A replacement transaction was confirmed instead
  - Transaction permanently failed

```go
OnComplete: func(tx *types.Transaction, receipt *types.Receipt, err error) {
    if err != nil {
        // Transaction submission or processing error
        s.logger.Warnf("Transaction error: %v", err)
        return
    }
    
    if receipt == nil {
        // Handle cancellation or replacement scenarios
        s.logger.Debugf("Transaction cancelled or replaced: %s", tx.Hash().Hex())
        // Clean up any scenario state if needed
        s.cleanupCancelledTransaction(tx)
        return
    }
    
    if receipt.Status == types.ReceiptStatusSuccessful {
        // Transaction successfully executed
        s.handleSuccessfulTransaction(tx, receipt)
    } else {
        // Transaction reverted (still included in block)
        s.logger.Warnf("Transaction reverted: %s", tx.Hash().Hex())
        s.handleRevertedTransaction(tx, receipt)
    }
}
```

### Balance Management

#### Automatic Balance Updates

Spamoor automatically manages wallet balances for **regular ETH transfers**:
- **Immediate deduction**: Balance decreases applied upon transaction submission
- **Confirmation updates**: Balance increases applied when transactions confirm
- **Gas fee handling**: Gas costs automatically deducted from sender balance
- **Failed transaction handling**: Reverts balance changes for failed submissions

#### Manual Balance Updates for Internal Transfers

Scenarios must manually update wallet balances for internal transfers (contract calls that move tokens/ETH between wallets):

```go
OnConfirm: func(tx *types.Transaction, receipt *types.Receipt) {
    if receipt != nil && receipt.Status == types.ReceiptStatusSuccessful {
        // Parse logs to determine actual transfers
        for _, log := range receipt.Logs {
            if log.Topics[0] == transferEventHash {
                // Decode transfer event
                from := common.BytesToAddress(log.Topics[1][:])
                to := common.BytesToAddress(log.Topics[2][:])
                value := new(big.Int).SetBytes(log.Data)
                
                // Update wallet balances
                fromWallet := s.walletPool.GetWalletByAddress(from)
                toWallet := s.walletPool.GetWalletByAddress(to)
                
                if fromWallet != nil {
                    fromWallet.SubBalance(value)
                }
                if toWallet != nil {
                    toWallet.AddBalance(value)
                }
            }
        }
    }
}
```

### Transaction Submission Best Practices

1. **Always use callbacks**: Handle confirmation and completion appropriately
2. **Enable rebroadcast**: For important transactions that must be confirmed
3. **Check receipt status**: Distinguish between confirmation and success
4. **Handle nil receipts**: Don't assume receipt is always available in callbacks
5. **Update balances manually**: For any internal transfers or token movements
6. **Use appropriate submission method**: Async for throughput, sync for dependencies

## Common Patterns

### Standard Scenario Structure

```go
package myscenario

import (
    "context"
    "fmt"
    "time"
    
    "gopkg.in/yaml.v3"
    "github.com/sirupsen/logrus"
    "github.com/spf13/pflag"
    
    "github.com/ethpandaops/spamoor/scenario"
    "github.com/ethpandaops/spamoor/spamoor"
    "github.com/ethpandaops/spamoor/txbuilder"
    "github.com/ethpandaops/spamoor/utils"
)

// Configuration structure
type ScenarioOptions struct {
    TotalCount  uint64 `yaml:"total_count"`
    Throughput  uint64 `yaml:"throughput"`
    MaxPending  uint64 `yaml:"max_pending"`
    MaxWallets  uint64 `yaml:"max_wallets"`
    BaseFee     uint64 `yaml:"base_fee"`
    TipFee      uint64 `yaml:"tip_fee"`
    GasLimit    uint64 `yaml:"gas_limit"`
    Timeout     string `yaml:"timeout"`
    ClientGroup string `yaml:"client_group"`
    LogTxs      bool   `yaml:"log_txs"`
    // Scenario-specific options...
}

// Main scenario struct
type Scenario struct {
    options    ScenarioOptions
    logger     *logrus.Entry
    walletPool *spamoor.WalletPool
    
    // Scenario-specific state...
}

// Scenario metadata
var ScenarioName = "my-scenario"
var ScenarioDefaultOptions = ScenarioOptions{
    TotalCount:  0,
    Throughput:  10,
    MaxPending:  0,
    MaxWallets:  0,
    BaseFee:     20,
    TipFee:      2,
    GasLimit:    21000,
    Timeout:     "",
    ClientGroup: "",
    LogTxs:      false,
}

var ScenarioDescriptor = scenario.Descriptor{
    Name:           ScenarioName,
    Description:    "Description of what this scenario does",
    DefaultOptions: ScenarioDefaultOptions,
    NewScenario:    newScenario,
}

func newScenario(logger logrus.FieldLogger) scenario.Scenario {
    return &Scenario{
        options: ScenarioDefaultOptions,
        logger:  logger.WithField("scenario", ScenarioName),
    }
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
    flags.Uint64VarP(&s.options.TotalCount, "count", "c", ScenarioDefaultOptions.TotalCount, "Total number of transactions")
    flags.Uint64VarP(&s.options.Throughput, "throughput", "t", ScenarioDefaultOptions.Throughput, "Transactions per slot")
    flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum pending transactions")
    // Add scenario-specific flags...
    return nil
}

func (s *Scenario) Init(options *scenario.Options) error {
    s.walletPool = options.WalletPool
    
    // Parse YAML configuration if provided
    if options.Config != "" {
        err := yaml.Unmarshal([]byte(options.Config), &s.options)
        if err != nil {
            return fmt.Errorf("failed to unmarshal config: %w", err)
        }
    }
    
    // Configure wallet pool
    if s.options.MaxWallets > 0 {
        s.walletPool.SetWalletCount(s.options.MaxWallets)
    }
    
    // Scenario-specific initialization...
    return nil
}

func (s *Scenario) Run(ctx context.Context) error {
    // Parse timeout if specified
    var timeout time.Duration
    if s.options.Timeout != "" {
        var err error
        timeout, err = time.ParseDuration(s.options.Timeout)
        if err != nil {
            return fmt.Errorf("invalid timeout duration: %w", err)
        }
    }
    
    // Run the transaction scenario
    return scenario.RunTransactionScenario(ctx, scenario.TransactionScenarioOptions{
        TotalCount:     s.options.TotalCount,
        Throughput:     s.options.Throughput,
        MaxPending:     s.options.MaxPending,
        Timeout:        timeout,
        WalletPool:     s.walletPool,
        Logger:         s.logger,
        ProcessNextTxFn: func(ctx context.Context, txIdx uint64, onComplete func()) (func(), error) {
			logger := s.logger
			tx, client, wallet, err := s.sendNextTransaction(ctx, txIdx, onComplete)
			if client != nil {
				logger = logger.WithField("rpc", client.GetName())
			}
			if tx != nil {
				logger = logger.WithField("nonce", tx.Nonce())
			}
			if wallet != nil {
				logger = logger.WithField("wallet", s.walletPool.GetWalletName(wallet.GetAddress()))
			}

			return func() {
				if err != nil {
					logger.Warnf("could not send transaction: %v", err)
				} else if s.options.LogTxs {
					logger.Infof("sent tx #%6d: %v", txIdx+1, tx.Hash().String())
				} else {
					logger.Debugf("sent tx #%6d: %v", txIdx+1, tx.Hash().String())
				}
			}, err
		},
    })
}

func (s *Scenario) sendNextTransaction(ctx context.Context, txIdx uint64, onComplete func()) (*types.Transaction, *spamoor.Client, *spamoor.Wallet, error) {
    // Standard wallet and client selection - see Wallet Management section
    wallet := s.walletPool.GetWallet(spamoor.SelectWalletRandom)
    client := s.walletPool.GetClient(spamoor.SelectClientRandom, 0, s.options.ClientGroup)
    
    // Get suggested fees and build transaction - see Transaction Building section
    feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(client, s.options.BaseFee, s.options.TipFee)
    // ... transaction building logic ...
    
    // Submit transaction - see Transaction Submission section for detailed options
    err = txpool.SendTransaction(ctx, wallet, signedTx, &spamoor.SendTransactionOptions{
        Client:      client,
        Rebroadcast: true,
        OnComplete: func(tx *types.Transaction, receipt *types.Receipt, err error) {
            onComplete()  // CRITICAL: Always call this when done
        },
    })
    
    return signedTx, client, wallet, err
}
```

## RunTransactionScenario Helper

The `scenario.RunTransactionScenario` helper function is the recommended way to implement transaction-based scenarios. It provides a complete execution engine with rate limiting, concurrency control, and lifecycle management.

### Function Signature

```go
func RunTransactionScenario(ctx context.Context, options TransactionScenarioOptions) error

type TransactionScenarioOptions struct {
    TotalCount     uint64                    // Total transactions to send (0 = unlimited)
    Throughput     uint64                    // Transactions per slot (0 = as fast as possible)
    MaxPending     uint64                    // Maximum pending transactions (0 = no limit)
    Timeout        time.Duration            // Scenario timeout (0 = no timeout)
    WalletPool     *spamoor.WalletPool      // Wallet pool for transaction submission
    Logger         *logrus.Entry            // Logger for the scenario
    ProcessNextTxFn func(ctx context.Context, txIdx uint64, onComplete func()) (func(), error)
}
```

### How It Works

1. **Initialization**
   - Sets up rate limiter based on throughput setting
   - Configures pending transaction counter
   - Subscribes to block updates for statistics
   - Starts timeout timer if specified

2. **Main Execution Loop**
   - Waits for rate limiter permit (respects throughput limit)
   - Checks pending transaction limit
   - Calls `ProcessNextTxFn` to build and submit transaction
   - Tracks transaction in pending counter
   - Chains logging callbacks for sequential output

3. **Transaction Processing**
   - Your `ProcessNextTxFn` builds and submits the transaction
   - You MUST call `onComplete()` when transaction processing finishes
   - The helper tracks completion and updates counters

4. **Completion Handling**
   - Waits for all pending transactions to complete
   - Cancels on context cancellation or timeout
   - Logs final statistics and throughput metrics

### Key Features

- **Rate Limiting**: Controls transaction throughput per slot (12 seconds on mainnet)
- **Concurrency Management**: Limits concurrent pending transactions
- **Lifecycle Management**: Handles scenario completion, cancellation, and timeouts
- **Progress Tracking**: Automatically tracks transaction counts and completion
- **Context Handling**: Respects context cancellation for clean shutdown
- **Throughput Monitoring**: Tracks transactions per block over different windows (5, 20, 60 blocks)
- **Dynamic Throughput**: Optional incremental throughput increases
- **Statistics Logging**: Periodic progress updates and final summary

### ProcessNextTxFn Callback

The `ProcessNextTxFn` is the core of your scenario implementation. This callback is called for each transaction that needs to be generated.

#### Callback Parameters

- **`ctx context.Context`**: The scenario context (check for cancellation)
- **`txIdx uint64`**: Zero-based index of the current transaction
- **`onComplete func()`**: Completion callback that MUST be called when transaction processing finishes

#### Return Values

- **`func()`**: A logging function that will be called after transaction submission
- **`error`**: Any error that occurred during transaction preparation/submission

#### Implementation Pattern

```go
func (s *Scenario) sendNextTransaction(ctx context.Context, txIdx uint64, onComplete func()) (func(), error) {
    // 1. Check context for cancellation
    select {
    case <-ctx.Done():
        onComplete() // Always call onComplete
        return nil, ctx.Err()
    default:
    }
    
    // 2. Select wallet and client
    wallet := s.walletPool.GetWallet(spamoor.SelectWalletRandom)
    client := s.walletPool.GetClient(spamoor.SelectClientRandom, 0, s.options.ClientGroup)
    
    // 3. Build transaction
    feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(client, s.options.BaseFee, s.options.TipFee)
    if err != nil {
        onComplete() // Call onComplete on early failure
        return nil, err
    }
    
    txData := &txbuilder.TxMetadata{
        GasTipCap: uint256.NewInt(tipCap),
        GasFeeCap: uint256.NewInt(feeCap),
        Gas:       s.options.GasLimit,
        // ... other transaction fields
    }
    
    signedTx, err := wallet.BuildBoundTx(ctx, txData)
    if err != nil {
        onComplete() // Call onComplete on build failure
        return nil, err
    }
    
    // 4. Submit transaction
    err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, signedTx, &spamoor.SendTransactionOptions{
        Client:      client,
        Rebroadcast: true,
        OnComplete: func(tx *types.Transaction, receipt *types.Receipt, err error) {
            // This is called when transaction is confirmed or fails
            onComplete() // CRITICAL: Must call this to signal completion
            
            // Optional: Handle transaction result
            if err != nil {
                s.logger.Warnf("transaction failed: %v", err)
            } else if receipt != nil && receipt.Status == types.ReceiptStatusSuccessful {
                s.logger.Debugf("transaction confirmed in block %d", receipt.BlockNumber.Uint64())
            }
        },
    })
    
    if err != nil {
        onComplete() // Call onComplete if submission fails
        return nil, err
    }
    
    // 5. Return logging function
    return func() {
        // This function is called after transaction is submitted
        // Use it for deferred logging to maintain output order
        if s.options.LogTxs {
            s.logger.Infof("sent tx #%d: %s", txIdx+1, signedTx.Hash().Hex())
        }
    }, nil
}
```

#### Critical onComplete() Callback Requirements

**The `onComplete()` callback MUST be called in exactly one of these scenarios:**

1. **In SendTransaction OnComplete callback** (most common):
```go
err := txpool.SendTransaction(ctx, wallet, signedTx, &spamoor.SendTransactionOptions{
    OnComplete: func(tx *types.Transaction, receipt *types.Receipt, err error) {
        onComplete() // Called when transaction processing completes
    },
})
```

2. **Directly if anything before transaction submission fails**:
```go
err := s.prepareTransaction(...)
if err != nil {
    onComplete() // Call immediately on preparation failure
    return nil, err
}
```

3. **In synchronous submission patterns**:
```go
receipt, err := txpool.SendAndAwaitTransaction(ctx, wallet, signedTx, options)
onComplete() // Call after synchronous completion
if err != nil {
    return nil, err
}
```

#### Usage Example

```go
func (s *Scenario) Run(ctx context.Context) error {
    // Parse timeout if specified
    var timeout time.Duration
    if s.options.Timeout != "" {
        timeout, _ = time.ParseDuration(s.options.Timeout)
    }
    
    // Use the helper function for standardized execution
    return scenario.RunTransactionScenario(ctx, scenario.TransactionScenarioOptions{
        TotalCount:      s.options.TotalCount,
        Throughput:      s.options.Throughput,
        MaxPending:      s.options.MaxPending,
        Timeout:         timeout,
        WalletPool:      s.walletPool,
        Logger:          s.logger,
        ProcessNextTxFn: s.sendNextTransaction, // Your transaction logic
    })
}
```

### Benefits of Using RunTransactionScenario

- **Consistent behavior**: All scenarios follow the same execution patterns
- **Built-in rate limiting**: Automatic slot-based throughput control
- **Proper lifecycle management**: Handles completion, cancellation, and timeouts
- **Progress logging**: Automatic progress updates and statistics
- **Error handling**: Standardized error handling and recovery
- **Context respect**: Proper context cancellation handling

**Always prefer this helper function over custom transaction loops** - it provides battle-tested transaction execution logic that handles edge cases and provides consistent behavior across all scenarios.

### Example Output

When running a scenario with RunTransactionScenario, you'll see per block output like:

```
INFO[11184] block 1587: submitted=100, pending=800, confirmed=120, throughput: 5B=160.00 tx/B, 20B=300.40 tx/B, 60B=320.72 tx/B  module=daemon scenario=eoatx spammer_id=1 wallets=120
```

### Advanced Features

#### Dynamic Throughput

You can implement dynamic throughput increases:

```go
return scenario.RunTransactionScenario(ctx, scenario.TransactionScenarioOptions{
    TotalCount:           s.options.TotalCount,
    Throughput:           s.options.Throughput,
    MaxPending:           s.options.MaxPending,
    ThroughputIncrement:  10,  // Increase by 10 tx/slot every period
    ThroughputIncrementPeriod: 60 * time.Second,  // Every minute
    // ...
})
```

#### Custom Completion Logic

For scenarios that need custom completion handling:

```go
ProcessNextTxFn: func(ctx context.Context, txIdx uint64, onComplete func()) (func(), error) {
    // Track custom metrics
    startTime := time.Now()
    
    // ... build and submit transaction ...
    
    err = txpool.SendTransaction(ctx, wallet, signedTx, &spamoor.SendTransactionOptions{
        OnComplete: func(tx *types.Transaction, receipt *types.Receipt, err error) {
            // Custom completion logic
            duration := time.Since(startTime)
            s.recordTransactionMetrics(tx, receipt, duration)
            
            // Always call onComplete
            onComplete()
        },
    })
    
    return logFn, err
}

### Standard Scenario Structure

```go
package myscenario

import (
    "context"
    "fmt"
    "time"
    
    "gopkg.in/yaml.v3"
    "github.com/sirupsen/logrus"
    "github.com/spf13/pflag"
    
    "github.com/ethpandaops/spamoor/scenario"
    "github.com/ethpandaops/spamoor/spamoor"
    "github.com/ethpandaops/spamoor/txbuilder"
    "github.com/ethpandaops/spamoor/utils"
)

// Configuration structure
type ScenarioOptions struct {
    TotalCount  uint64 `yaml:"total_count"`
    Throughput  uint64 `yaml:"throughput"`
    MaxPending  uint64 `yaml:"max_pending"`
    MaxWallets  uint64 `yaml:"max_wallets"`
    BaseFee     uint64 `yaml:"base_fee"`
    TipFee      uint64 `yaml:"tip_fee"`
    GasLimit    uint64 `yaml:"gas_limit"`
    Timeout     string `yaml:"timeout"`
    ClientGroup string `yaml:"client_group"`
    LogTxs      bool   `yaml:"log_txs"`
    // Scenario-specific options...
}

// Main scenario struct
type Scenario struct {
    options    ScenarioOptions
    logger     *logrus.Entry
    walletPool *spamoor.WalletPool
    
    // Scenario-specific state...
}

// Scenario metadata
var ScenarioName = "my-scenario"
var ScenarioDefaultOptions = ScenarioOptions{
    TotalCount:  0,
    Throughput:  10,
    MaxPending:  0,
    MaxWallets:  0,
    BaseFee:     20,
    TipFee:      2,
    GasLimit:    21000,
    Timeout:     "",
    ClientGroup: "",
    LogTxs:      false,
}

var ScenarioDescriptor = scenario.Descriptor{
    Name:           ScenarioName,
    Description:    "Description of what this scenario does",
    DefaultOptions: ScenarioDefaultOptions,
    NewScenario:    newScenario,
}

func newScenario(logger logrus.FieldLogger) scenario.Scenario {
    return &Scenario{
        options: ScenarioDefaultOptions,
        logger:  logger.WithField("scenario", ScenarioName),
    }
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
    flags.Uint64VarP(&s.options.TotalCount, "count", "c", ScenarioDefaultOptions.TotalCount, "Total number of transactions")
    flags.Uint64VarP(&s.options.Throughput, "throughput", "t", ScenarioDefaultOptions.Throughput, "Transactions per slot")
    flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum pending transactions")
    // Add scenario-specific flags...
    return nil
}

func (s *Scenario) Init(options *scenario.Options) error {
    s.walletPool = options.WalletPool
    
    // Parse YAML configuration if provided
    if options.Config != "" {
        err := yaml.Unmarshal([]byte(options.Config), &s.options)
        if err != nil {
            return fmt.Errorf("failed to unmarshal config: %w", err)
        }
    }
    
    // Configure wallet pool
    if s.options.MaxWallets > 0 {
        s.walletPool.SetWalletCount(s.options.MaxWallets)
    }
    
    // Scenario-specific initialization...
    return nil
}

func (s *Scenario) Run(ctx context.Context) error {
    // Parse timeout if specified
    var timeout time.Duration
    if s.options.Timeout != "" {
        var err error
        timeout, err = time.ParseDuration(s.options.Timeout)
        if err != nil {
            return fmt.Errorf("invalid timeout duration: %w", err)
        }
    }
    
    // Run the transaction scenario
    return scenario.RunTransactionScenario(ctx, scenario.TransactionScenarioOptions{
        TotalCount:     s.options.TotalCount,
        Throughput:     s.options.Throughput,
        MaxPending:     s.options.MaxPending,
        Timeout:        timeout,
        WalletPool:     s.walletPool,
        Logger:         s.logger,
        ProcessNextTxFn: func(ctx context.Context, txIdx uint64, onComplete func()) (func(), error) {
			logger := s.logger
			tx, client, wallet, err := s.sendNextTransaction(ctx, txIdx, onComplete)
			if client != nil {
				logger = logger.WithField("rpc", client.GetName())
			}
			if tx != nil {
				logger = logger.WithField("nonce", tx.Nonce())
			}
			if wallet != nil {
				logger = logger.WithField("wallet", s.walletPool.GetWalletName(wallet.GetAddress()))
			}

			return func() {
				if err != nil {
					logger.Warnf("could not send transaction: %v", err)
				} else if s.options.LogTxs {
					logger.Infof("sent tx #%6d: %v", txIdx+1, tx.Hash().String())
				} else {
					logger.Debugf("sent tx #%6d: %v", txIdx+1, tx.Hash().String())
				}
			}, err
		},
    })
}

func (s *Scenario) sendNextTransaction(ctx context.Context, txIdx uint64, onComplete func()) (*types.Transaction, *spamoor.Client, *spamoor.Wallet, error) {
    // Standard wallet and client selection - see Wallet Management section
    wallet := s.walletPool.GetWallet(spamoor.SelectWalletRandom)
    client := s.walletPool.GetClient(spamoor.SelectClientRandom, 0, s.options.ClientGroup)
    
    // Get suggested fees and build transaction - see Transaction Building section
    feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(client, s.options.BaseFee, s.options.TipFee)
    // ... transaction building logic ...
    
    // Submit transaction - see Transaction Submission section for detailed options
    err = txpool.SendTransaction(ctx, wallet, signedTx, &spamoor.SendTransactionOptions{
        Client:      client,
        Rebroadcast: true,
        OnComplete: func(tx *types.Transaction, receipt *types.Receipt, err error) {
            onComplete()  // CRITICAL: Always call this when done
        },
    })
    
    return signedTx, client, wallet, err
}
```

## Contract Deployments

### Deployment Patterns

For scenarios that require contract deployments, follow these patterns:

#### One-time Deployments

```go
func (s *Scenario) deployContract(ctx context.Context) (common.Address, error) {
    deployerWallet := s.walletPool.GetWellKnownWallet("deployer")
    if deployerWallet == nil {
        return common.Address{}, fmt.Errorf("deployer wallet not found")
    }
    
    client := s.walletPool.GetClient(spamoor.SelectClientRandom, 0, s.options.ClientGroup)
    if client == nil {
        return common.Address{}, fmt.Errorf("no client available")
    }
    
    // Get deployment bytecode
    bytecode := common.Hex2Bytes("608060405234801561001057600080fd5b50...")
    
    // Build deployment transaction
    feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(
        client, s.options.BaseFee, s.options.TipFee)
    if err != nil {
        return common.Address{}, err
    }
    
    txData := &txbuilder.TxMetadata{
        GasTipCap: uint256.NewInt(tipCap),
        GasFeeCap: uint256.NewInt(feeCap),
        Gas:       1000000, // Sufficient gas for deployment
        To:        nil,     // Contract creation
        Value:     uint256.NewInt(0),
        Data:      bytecode,
    }
    
    signedTx, err := deployerWallet.BuildBoundTx(ctx, txData)
    if err != nil {
        return common.Address{}, err
    }
    
    // Submit and wait for confirmation
    txpool := s.walletPool.GetTxPool()
    err = txpool.SendTransaction(ctx, signedTx, &spamoor.SendTransactionOptions{
        AwaitConfirmation: true,
        LogTx:            true,
    })
    if err != nil {
        return common.Address{}, err
    }
    
    // Calculate contract address
    contractAddr := crypto.CreateAddress(deployerWallet.GetAddress(), signedTx.Nonce())
    s.logger.Infof("deployed contract at %s", contractAddr.Hex())
    
    return contractAddr, nil
}
```

#### Reusable Deployments (Uniswap Pattern)

For complex deployments that should be reused across runs, use nonce checking:

```go
func (s *Scenario) deployContracts(redeploy bool) (*DeploymentInfo, error) {
    client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, s.options.ClientGroup)
    deployerWallet := s.walletPool.GetWellKnownWallet("deployer")
    
    feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(client, s.options.BaseFee, s.options.TipFee)
    if err != nil {
        return nil, fmt.Errorf("could not get tx fee: %w", err)
    }
    
    deploymentTxs := []*types.Transaction{}
    deploymentInfo := &DeploymentInfo{}
    
    // Get current deployer nonce to check if deployments already exist
    deployerNonce := deployerWallet.GetNonce()
    contractNonce := uint64(0)  // Track expected nonce for each contract
    usedNonce := uint64(0)
    
    // Deploy first contract only if nonce indicates it doesn't exist
    if redeploy || deployerNonce <= contractNonce {
        tx, err := deployerWallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
            GasFeeCap: uint256.MustFromBig(feeCap),
            GasTipCap: uint256.MustFromBig(tipCap),
            Gas:       2000000,
            Value:     uint256.NewInt(0),
        }, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
            _, deployTx, _, err := contract.DeployMyContract(transactOpts, client.GetEthClient(), param1)
            return deployTx, err
        })
        if err != nil {
            return nil, fmt.Errorf("could not deploy contract: %w", err)
        }
        deploymentTxs = append(deploymentTxs, tx)
        usedNonce = tx.Nonce()
    } else {
        usedNonce = contractNonce  // Contract already exists at this nonce
    }
    contractNonce++
    
    // Calculate deterministic contract address
    deploymentInfo.ContractAddr = crypto.CreateAddress(deployerWallet.GetAddress(), usedNonce)
    deploymentInfo.Contract, err = contract.NewMyContract(deploymentInfo.ContractAddr, client.GetEthClient())
    if err != nil {
        return nil, fmt.Errorf("could not create contract instance: %w", err)
    }
    
    // Deploy second contract (similar pattern)
    if redeploy || deployerNonce <= contractNonce {
        // ... deploy second contract
        usedNonce = tx.Nonce()
    } else {
        usedNonce = contractNonce
    }
    contractNonce++
    
    // Submit all deployment transactions if any were created
    if len(deploymentTxs) > 0 {
        s.logger.Infof("deploying %d contracts...", len(deploymentTxs))
        _, err := s.walletPool.GetTxPool().SendTransactionBatch(ctx, deployerWallet, deploymentTxs, &spamoor.BatchOptions{
            SendTransactionOptions: spamoor.SendTransactionOptions{
                Client: client,
            },
        })
        if err != nil {
            return nil, fmt.Errorf("could not send deployment txs: %w", err)
        }
        s.logger.Infof("contract deployment complete")
    } else {
        s.logger.Infof("contracts already deployed, skipping deployment")
    }
    
    return deploymentInfo, nil
}
```

This pattern:
- **Checks deployer nonce**: If current nonce > expected contract nonce, contract already exists
- **Calculates deterministic addresses**: Uses `crypto.CreateAddress(deployerAddr, nonce)`
- **Conditional deployment**: Only deploys if `redeploy` flag or nonce check indicates missing contracts
- **Batch submission**: Efficiently submits multiple deployment transactions
- **Reusable across runs**: Same deployer wallet will create same addresses

## Best Practices

### Context Handling

Always respect context cancellation - scenarios must stop all operations when context is cancelled:

```go
func (s *Scenario) Run(ctx context.Context) error {
    // Check context before expensive operations
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }
    
    // Use context in all operations and pass to transaction scenario
    return scenario.RunTransactionScenario(ctx, options)
}
```

### Logging Standards

Use structured logging with consistent fields and appropriate log levels:

```go
// Transaction logging with structured fields
s.logger.WithFields(logrus.Fields{
    "txHash": tx.Hash().Hex(),
    "from":   wallet.GetAddress().Hex(),
    "nonce":  tx.Nonce(),
    "status": receipt.Status,
}).Info("transaction confirmed")

// Log levels: Debug (detailed) → Info (important) → Warn (recoverable) → Error (critical)
s.logger.Debug("building transaction")
s.logger.Info("transaction submitted") 
s.logger.Warn("transaction failed")
s.logger.Error("scenario initialization failed")
```

### Configuration Validation

Validate configuration in the Init method:

```go
func (s *Scenario) Init(options *scenario.Options) error {
    // Validate required parameters
    if s.options.TotalCount == 0 && s.options.Throughput == 0 {
        return fmt.Errorf("either total_count or throughput must be specified")
    }
    
    // Warn about potentially problematic configurations
    if s.options.MaxPending > 0 && s.options.MaxPending < s.options.Throughput {
        s.logger.Warnf("max_pending (%d) < throughput (%d)", s.options.MaxPending, s.options.Throughput)
    }
    
    return nil
}
```

## Testing Scenarios

### Quick Testing with DevNet

The fastest way to test scenarios is using the built-in development environment:

```bash
# Start a complete Ethereum testnet with spamoor daemon
make devnet-run
```

This provides:
- **Full Ethereum testnet**: Geth, Reth, and Lighthouse clients
- **Pre-funded accounts**: Well-known private key with ETH ready to use
- **Block explorers**: Dora and Blockscout for transaction monitoring
- **Web interface**: Spamoor dashboard at http://localhost:8080
- **Multiple RPC endpoints**: Automatically configured in `.hack/devnet/generated-hosts.txt`

Test your scenario via the web interface or API:
```bash
# Test via API (daemon already running from devnet-run)
curl -X POST http://localhost:8080/api/spammer \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test My Scenario",
    "scenario": "my-scenario",
    "config": "total_count: 10\nthroughput: 2",
    "startImmediately": true
  }'
```

Clean up when done:
```bash
make devnet-clean
```

### Local Testing with Custom RPC

```bash
# Build scenario
make build

# Test with minimal configuration
./bin/spamoor my-scenario \
  --privkey "0x..." \
  --rpchost "http://localhost:8545" \
  --count 10 \
  --verbose

# Test with YAML configuration
echo "
total_count: 5
throughput: 2
gas_limit: 100000
" > test-config.yaml

./bin/spamoor-daemon \
  --privkey "0x..." \
  --rpchost "http://localhost:8545" \
  --startup-spammer test-config.yaml
```


## Example Implementation

Here's a complete minimal scenario implementation:

```go
package simpleeoa

import (
    "context"
    "fmt"
    "math/rand"
    "time"
    
    "gopkg.in/yaml.v3"
    "github.com/ethereum/go-ethereum/common"
    "github.com/holiman/uint256"
    "github.com/sirupsen/logrus"
    "github.com/spf13/pflag"
    
    "github.com/ethpandaops/spamoor/scenario"
    "github.com/ethpandaops/spamoor/spamoor"
    "github.com/ethpandaops/spamoor/txbuilder"
)

type ScenarioOptions struct {
    TotalCount  uint64 `yaml:"total_count"`
    Throughput  uint64 `yaml:"throughput"`
    MaxPending  uint64 `yaml:"max_pending"`
    BaseFee     uint64 `yaml:"base_fee"`
    TipFee      uint64 `yaml:"tip_fee"`
    Amount      uint64 `yaml:"amount"`
    ClientGroup string `yaml:"client_group"`
}

type Scenario struct {
    options    ScenarioOptions
    logger     *logrus.Entry
    walletPool *spamoor.WalletPool
}

var ScenarioName = "simple-eoa"
var ScenarioDefaultOptions = ScenarioOptions{
    TotalCount:  10,
    Throughput:  5,
    MaxPending:  20,
    BaseFee:     20,
    TipFee:      2,
    Amount:      1000000000000000000, // 1 ETH in wei
    ClientGroup: "",
}

var ScenarioDescriptor = scenario.Descriptor{
    Name:           ScenarioName,
    Description:    "Simple EOA to EOA transfers",
    DefaultOptions: ScenarioDefaultOptions,
    NewScenario:    newScenario,
}

func newScenario(logger logrus.FieldLogger) scenario.Scenario {
    return &Scenario{
        options: ScenarioDefaultOptions,
        logger:  logger.WithField("scenario", ScenarioName),
    }
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
    flags.Uint64VarP(&s.options.TotalCount, "count", "c", 
        ScenarioDefaultOptions.TotalCount, "Total transactions to send")
    flags.Uint64VarP(&s.options.Throughput, "throughput", "t", 
        ScenarioDefaultOptions.Throughput, "Transactions per slot")
    flags.Uint64Var(&s.options.Amount, "amount", 
        ScenarioDefaultOptions.Amount, "Transfer amount in wei")
    return nil
}

func (s *Scenario) Init(options *scenario.Options) error {
    s.walletPool = options.WalletPool
    
    if options.Config != "" {
        err := yaml.Unmarshal([]byte(options.Config), &s.options)
        if err != nil {
            return fmt.Errorf("failed to unmarshal config: %w", err)
        }
    }
    
    // Configure wallet pool
    s.walletPool.SetWalletCount(10)
    
    return nil
}

func (s *Scenario) Run(ctx context.Context) error {
    return scenario.RunTransactionScenario(ctx, scenario.TransactionScenarioOptions{
        TotalCount:      s.options.TotalCount,
        Throughput:      s.options.Throughput,
        MaxPending:      s.options.MaxPending,
        WalletPool:      s.walletPool,
        Logger:          s.logger,
        ProcessNextTxFn: s.sendNextTransaction,
    })
}

func (s *Scenario) sendNextTransaction(ctx context.Context, txIdx uint64, onComplete func()) (func(), error) {
    // Wallet and client selection (see Wallet Management section for selection strategies)
    wallet := s.walletPool.GetWallet(spamoor.SelectWalletRandom)
    targetWallet := s.walletPool.GetWallet(spamoor.SelectWalletRandom)
    client := s.walletPool.GetClient(spamoor.SelectClientRandom, 0, s.options.ClientGroup)
    
    // Fee calculation and transaction building (see Transaction Building section)
    feeCap, tipCap, _ := s.walletPool.GetTxPool().GetSuggestedFees(client, s.options.BaseFee, s.options.TipFee)
    
    txData := &txbuilder.TxMetadata{
        GasTipCap: uint256.NewInt(tipCap),
        GasFeeCap: uint256.NewInt(feeCap),
        Gas:       21000,
        To:        &targetWallet.GetAddress(),
        Value:     uint256.NewInt(s.options.Amount),
    }
    
    signedTx, _ := wallet.BuildBoundTx(ctx, txData)
    
    // Transaction submission (see Transaction Submission section for all options)
    err := s.walletPool.GetTxPool().SendTransaction(ctx, wallet, signedTx, &spamoor.SendTransactionOptions{
        Client:     client,
        OnComplete: onComplete, // CRITICAL: Always call this
    })
    
    // Logging callback (see Best Practices section for logging standards)
    return func() {
        s.logger.WithField("txHash", signedTx.Hash().Hex()).Info("sent EOA transfer")
    }, err
}
```

This example demonstrates all the key concepts:
- Proper scenario structure and registration
- Configuration handling with YAML and flags
- Wallet selection and management
- Transaction building with txbuilder
- Using RunTransactionScenario for execution
- Proper context handling and logging
- Following the onComplete callback pattern

Follow these patterns and best practices to create robust, efficient scenarios that integrate seamlessly with the Spamoor ecosystem.