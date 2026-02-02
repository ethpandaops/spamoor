# Dynamic Scenario Loading Guide

This guide covers how to create and use dynamic scenarios in Spamoor. Dynamic scenarios are Go source files that are loaded at runtime using the Yaegi interpreter, allowing you to add new scenarios without recompiling Spamoor.

## Table of Contents

- [Overview](#overview)
- [Quick Start](#quick-start)
- [Creating a Dynamic Scenario](#creating-a-dynamic-scenario)
- [Available Packages](#available-packages)
- [CLI Usage](#cli-usage)
- [Auto-loading and Hot-reload](#auto-loading-and-hot-reload)
- [Validating Scenarios](#validating-scenarios)
- [Adding New Symbol Packages](#adding-new-symbol-packages)
- [Yaegi Limitations and Workarounds](#yaegi-limitations-and-workarounds)
- [Troubleshooting](#troubleshooting)

## Overview

Dynamic scenarios provide a way to extend Spamoor without modifying the core codebase or recompiling. They are useful for:

- **Rapid prototyping**: Test new transaction patterns quickly
- **Custom deployments**: Add organization-specific scenarios
- **External contributions**: Share scenarios without core integration

Dynamic scenarios are interpreted at runtime by Yaegi and have access to a subset of Spamoor's APIs through extracted symbols.

## Quick Start

1. Create a directory for your scenario:
```bash
mkdir -p scenarios/external/my-scenario
```

2. Create a Go source file (e.g., `my-scenario.go`):
```go
package main

import (
    "context"

    "github.com/ethpandaops/spamoor/scenario"
    "github.com/ethpandaops/spamoor/spamoor"
    "github.com/sirupsen/logrus"
    "github.com/spf13/pflag"
)

var ScenarioDescriptor = scenario.Descriptor{
    Name:        "my-scenario",
    Description: "My custom transaction scenario",
    NewScenario: NewScenario,
}

type Scenario struct {
    logger  *logrus.Entry
    options *Options
}

type Options struct {
    Throughput uint64 `yaml:"throughput"`
}

var DefaultOptions = &Options{
    Throughput: 10,
}

func NewScenario(logger logrus.FieldLogger) scenario.Scenario {
    return &Scenario{
        logger:  logger.WithField("scenario", "my-scenario"),
        options: &Options{},
    }
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
    flags.Uint64VarP(&s.options.Throughput, "throughput", "t", 10, "Transactions per slot")
    return nil
}

func (s *Scenario) Init(opts *scenario.Options) error {
    // Initialize scenario with wallet pool
    return nil
}

func (s *Scenario) Run(ctx context.Context) error {
    s.logger.Info("Running my-scenario")
    // Implement transaction logic here
    <-ctx.Done()
    return nil
}

func main() {}
```

3. Validate your scenario:
```bash
spamoor validate-scenario scenarios/external/my-scenario/my-scenario.go
```

4. Run with the scenario:
```bash
spamoor my-scenario -h localhost:8545 -p YOUR_PRIVATE_KEY
```

## Creating a Dynamic Scenario

### Required Structure

Every dynamic scenario must have:

1. **Package declaration**: Must be `package main`
2. **ScenarioDescriptor variable**: A `scenario.Descriptor` struct exported at package level
3. **Empty main function**: Required for valid Go syntax

### Scenario Descriptor Fields

```go
var ScenarioDescriptor = scenario.Descriptor{
    Name:           "scenario-name",      // Unique identifier (required)
    Description:    "What it does",       // Shown in CLI help (recommended)
    Aliases:        []string{"alias1"},   // Alternative names (optional)
    DefaultOptions: DefaultOptions,       // Default config (optional)
    NewScenario:    NewScenario,          // Factory function (required)
}
```

### Implementing the Scenario Interface

Your scenario struct must implement the `scenario.Scenario` interface:

```go
type Scenario interface {
    Flags(flags *pflag.FlagSet) error    // Register CLI flags
    Init(options *Options) error          // Initialize with wallet pool
    Run(ctx context.Context) error        // Execute the scenario
}
```

### Configuration with YAML

Options are automatically parsed from YAML configuration. Use struct tags:

```go
type Options struct {
    Throughput  uint64 `yaml:"throughput"`
    TotalCount  uint64 `yaml:"total_count"`
    MaxPending  uint64 `yaml:"max_pending"`
}
```

## Available Packages

The following packages have extracted symbols and are available for use in dynamic scenarios:

### Spamoor Packages

| Package | Import Path | Purpose |
|---------|-------------|---------|
| scenario | `github.com/ethpandaops/spamoor/scenario` | Core scenario interfaces and helpers |
| spamoor | `github.com/ethpandaops/spamoor/spamoor` | Client pool, wallet pool, transaction pool |
| txbuilder | `github.com/ethpandaops/spamoor/txbuilder` | Transaction building utilities |
| utils | `github.com/ethpandaops/spamoor/utils` | Common utilities |

### Third-party Packages

| Package | Import Path | Purpose |
|---------|-------------|---------|
| logrus | `github.com/sirupsen/logrus` | Structured logging |
| pflag | `github.com/spf13/pflag` | CLI flag parsing |
| uint256 | `github.com/holiman/uint256` | Big integer operations |
| common | `github.com/ethereum/go-ethereum/common` | Ethereum types (Address, Hash) |
| types | `github.com/ethereum/go-ethereum/core/types` | Transaction types |
| abi | `github.com/ethereum/go-ethereum/accounts/abi` | ABI encoding/decoding |
| bind | `github.com/ethereum/go-ethereum/accounts/abi/bind` | Contract bindings |
| crypto | `github.com/ethereum/go-ethereum/crypto` | Cryptographic functions |
| yaml | `gopkg.in/yaml.v3` | YAML parsing |

## CLI Usage

### Loading a Single Scenario File

```bash
spamoor --scenario-file /path/to/scenario.go <scenario-name> [options]
```

### Loading from a Directory

```bash
spamoor --scenario-dir /path/to/scenarios/ <scenario-name> [options]
```

### Using External Directory

Place scenarios in `scenarios/external/<scenario-name>/` and they will be auto-loaded at startup.

## Auto-loading and Hot-reload

### Auto-loading at Startup

Spamoor automatically loads scenarios from `scenarios/external/` on startup:

```
scenarios/
  external/
    eoatx/
      eoatx.go
    erc20_bloater/
      erc20_bloater.go
```

Each subdirectory should contain `.go` files with scenario definitions.

### Hot-reload in Daemon Mode

In daemon mode, you can reload scenarios without restarting:

**Via API:**
```bash
curl -X POST http://localhost:8080/api/scenarios/reload
```

**Via GUI:**
Click the "Reload Scenarios" button on the dashboard.

**Specifying a different directory:**
```bash
curl -X POST "http://localhost:8080/api/scenarios/reload?dir=/custom/path"
```

## Validating Scenarios

Use the validation command to check scenarios before running:

```bash
spamoor validate-scenario /path/to/scenario.go
```

This will:
1. Load the scenario using Yaegi
2. Verify the `ScenarioDescriptor` fields
3. Instantiate the scenario
4. Check registered flags
5. Report any issues

Example output:
```
Validating scenario: my-scenario.go

✓ Scenario loaded successfully

Descriptor:
  Name:        my-scenario
  Description: My custom scenario
  NewScenario: ✓ defined
  DefaultOpts: ✓ defined (*main.Options)

Instantiation:
  ✓ Instance created successfully
  ✓ Flags registered: 3

Result: ✓ Scenario 'my-scenario' is valid
```

## Adding New Symbol Packages

If your scenario needs a package that isn't available, you can extract its symbols:

1. Install yaegi:
```bash
go install github.com/traefik/yaegi/cmd/yaegi@v0.16.1
```

2. Extract symbols:
```bash
cd scenarios/loader
yaegi extract github.com/your/package
```

3. Fix the package declaration:
```bash
perl -i -pe 's/^package \w+$/package loader/' symbols_*.go
```

4. Rebuild spamoor:
```bash
go build ./cmd/spamoor
```

Alternatively, use the Makefile target:
```bash
make generate-symbols
```

## Yaegi Limitations and Workarounds

Yaegi is a Go interpreter with some limitations. Here are common issues and solutions:

### Channel Type Aliases in Closures

**Problem**: Sending to channels with type aliases inside closures may cause panics.

**Workaround**: Use raw channel types instead of aliases, or extract the send operation outside the closure.

### Complex Type Assertions

**Problem**: Some complex type assertions may fail.

**Workaround**: Use intermediate interface conversions or restructure the code.

### Unsupported Language Features

Yaegi doesn't support all Go features. If you encounter issues:
- Avoid `go:embed` directives
- Avoid CGO
- Keep generics usage simple
- Avoid reflection on unexported fields

### Interface Satisfaction

**Problem**: Types defined in dynamic scenarios may not satisfy interfaces defined in extracted symbols.

**Workaround**: Ensure your types exactly match the expected interface, including method receiver types (pointer vs value).

## Troubleshooting

### Error: "undefined: package.Symbol"

The package doesn't have extracted symbols. Extract them:
```bash
cd scenarios/loader
yaegi extract github.com/the/package
perl -i -pe 's/^package \w+$/package loader/' symbols_*.go
```

### Error: "ScenarioDescriptor is not of type scenario.Descriptor"

Ensure you're using the correct type from extracted symbols, not a locally defined type.

### Error: "CFG post-order panic"

This is usually a Yaegi limitation. Try:
- Simplifying the code structure
- Moving complex operations outside closures
- Using explicit type conversions

### Scenario Not Found

Verify:
1. The scenario file is in the correct directory
2. The `ScenarioDescriptor` variable is exported (capitalized)
3. The `Name` field matches what you're trying to run

### Flags Not Working

Ensure your `Flags()` method returns `nil` on success and registers flags on the provided `*pflag.FlagSet`.

### Scenario Loads But Doesn't Run

Check:
1. The `Run()` method is implemented correctly
2. Context cancellation is handled
3. Errors are returned, not just logged

## Example: EOA Transfer Scenario

Here's a complete example of a dynamic EOA transfer scenario:

```go
package main

import (
    "context"
    "fmt"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethpandaops/spamoor/scenario"
    "github.com/ethpandaops/spamoor/spamoor"
    "github.com/holiman/uint256"
    "github.com/sirupsen/logrus"
    "github.com/spf13/pflag"
)

var ScenarioDescriptor = scenario.Descriptor{
    Name:           "dynamic-eoatx",
    Description:    "Simple EOA transfer scenario (dynamic)",
    DefaultOptions: DefaultOptions,
    NewScenario:    NewScenario,
}

type Scenario struct {
    logger     *logrus.Entry
    walletPool *spamoor.WalletPool
    options    *Options
}

type Options struct {
    TotalCount uint64 `yaml:"total_count"`
    Throughput uint64 `yaml:"throughput"`
    MaxPending uint64 `yaml:"max_pending"`
    MaxWallets uint64 `yaml:"max_wallets"`
    Amount     uint64 `yaml:"amount"`
}

var DefaultOptions = &Options{
    TotalCount: 0,
    Throughput: 10,
    MaxPending: 100,
    MaxWallets: 50,
    Amount:     1000000000000, // 0.000001 ETH
}

func NewScenario(logger logrus.FieldLogger) scenario.Scenario {
    return &Scenario{
        logger:  logger.WithField("scenario", "dynamic-eoatx"),
        options: &Options{},
    }
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
    flags.Uint64VarP(&s.options.TotalCount, "count", "c", DefaultOptions.TotalCount, "Total transaction count")
    flags.Uint64VarP(&s.options.Throughput, "throughput", "t", DefaultOptions.Throughput, "Transactions per slot")
    flags.Uint64VarP(&s.options.MaxPending, "max-pending", "p", DefaultOptions.MaxPending, "Max pending transactions")
    flags.Uint64VarP(&s.options.MaxWallets, "max-wallets", "w", DefaultOptions.MaxWallets, "Max wallets to use")
    flags.Uint64Var(&s.options.Amount, "amount", DefaultOptions.Amount, "Transfer amount in wei")
    return nil
}

func (s *Scenario) Init(opts *scenario.Options) error {
    s.walletPool = opts.WalletPool
    s.walletPool.SetWalletCount(s.options.MaxWallets)
    return nil
}

func (s *Scenario) Run(ctx context.Context) error {
    return scenario.RunTransactionScenario(ctx, &scenario.TransactionScenarioOptions{
        WalletPool:    s.walletPool,
        Logger:        s.logger,
        Throughput:    s.options.Throughput,
        MaxPending:    s.options.MaxPending,
        TotalCount:    s.options.TotalCount,
        ProcessNextTx: s.processNextTx,
    })
}

func (s *Scenario) processNextTx(ctx context.Context, txIdx uint64, opts *scenario.TransactionScenarioTxOptions) (*types.Transaction, *spamoor.Wallet, func(), error) {
    wallet := s.walletPool.GetWallet(spamoor.SelectByIndex, int(txIdx))
    if wallet == nil {
        return nil, nil, nil, fmt.Errorf("no wallet available")
    }

    toAddr := s.walletPool.GetWallet(spamoor.SelectRandom, 0).GetAddress()
    amount := uint256.NewInt(s.options.Amount)

    tx, err := wallet.BuildTransaction(ctx, func(txData *types.DynamicFeeTx) (*types.DynamicFeeTx, error) {
        txData.To = (*common.Address)(&toAddr)
        txData.Value = amount.ToBig()
        txData.Gas = 21000
        return txData, nil
    })
    if err != nil {
        return nil, wallet, nil, err
    }

    return tx, wallet, nil, nil
}

func main() {}
```

Save this as `scenarios/external/dynamic-eoatx/dynamic-eoatx.go` and it will be auto-loaded.
