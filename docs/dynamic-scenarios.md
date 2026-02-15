# Plugin Development Guide

This guide covers how to create and use plugins to extend Spamoor with custom transaction scenarios. Plugins are Go source files loaded at runtime using the [Yaegi](https://github.com/traefik/yaegi) interpreter, allowing you to add new scenarios without recompiling Spamoor.

## Table of Contents

- [Overview](#overview)
- [Quick Start](#quick-start)
- [Plugin Structure](#plugin-structure)
- [Creating a Scenario](#creating-a-scenario)
- [Available Packages](#available-packages)
- [CLI Usage](#cli-usage)
- [Daemon Mode](#daemon-mode)
- [Validating Plugins](#validating-plugins)
- [Yaegi Limitations](#yaegi-limitations)
- [Troubleshooting](#troubleshooting)

## Overview

Plugins provide a way to extend Spamoor without modifying the core codebase or recompiling. They are useful for:

- **Rapid prototyping**: Test new transaction patterns quickly
- **Custom deployments**: Add organization-specific scenarios
- **External contributions**: Share scenarios without core integration

Plugins are interpreted at runtime by Yaegi and have access to Spamoor's APIs through pre-extracted symbols.

## Quick Start

1. Create a plugin directory:
```bash
mkdir -p my-plugin/my-scenario
```

2. Create `my-plugin/plugin.go`:
```go
package plugin

import (
    "github.com/ethpandaops/spamoor/plugins/my-plugin/my-scenario"
    "github.com/ethpandaops/spamoor/scenario"
)

var PluginDescriptor = scenario.PluginDescriptor{
    Name:        "my-plugin",
    Description: "My custom plugin",
    Categories: []*scenario.Category{
        {
            Name:        "Custom",
            Description: "Custom scenarios",
            Descriptors: []*scenario.Descriptor{
                &myscenario.ScenarioDescriptor,
            },
        },
    },
}
```

3. Create `my-plugin/my-scenario/scenario.go`:
```go
package myscenario

import (
    "context"

    "github.com/ethpandaops/spamoor/scenario"
    "github.com/ethpandaops/spamoor/spamoor"
    "github.com/sirupsen/logrus"
    "github.com/spf13/pflag"
)

type Options struct {
    Throughput uint64 `yaml:"throughput"`
}

var DefaultOptions = Options{Throughput: 10}

var ScenarioDescriptor = scenario.Descriptor{
    Name:           "my-scenario",
    Description:    "My custom scenario",
    DefaultOptions: DefaultOptions,
    NewScenario:    newScenario,
}

type Scenario struct {
    options    Options
    logger     *logrus.Entry
    walletPool *spamoor.WalletPool
}

func newScenario(logger logrus.FieldLogger) scenario.Scenario {
    return &Scenario{
        options: DefaultOptions,
        logger:  logger.WithField("scenario", "my-scenario"),
    }
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
    flags.Uint64VarP(&s.options.Throughput, "throughput", "t", 10, "Transactions per slot")
    return nil
}

func (s *Scenario) Init(opts *scenario.Options) error {
    s.walletPool = opts.WalletPool
    return nil
}

func (s *Scenario) Run(ctx context.Context) error {
    s.logger.Info("Running my-scenario")
    <-ctx.Done()
    return nil
}
```

4. Validate and run:
```bash
# Validate the plugin
spamoor-utils validate-plugin ./my-plugin

# Run a scenario from the plugin
spamoor --plugin ./plugins/my-plugin my-scenario -h http://localhost:8545 -p YOUR_PRIVATE_KEY
```

## Plugin Structure

A plugin consists of a directory with Go source files:

```
my-plugin/
├── plugin.go              # Required: exports PluginDescriptor
├── scenario1/
│   └── scenario1.go       # Scenario implementation
├── scenario2/
│   └── scenario2.go       # Another scenario
└── contract/              # Optional: generated contract bindings
    ├── MyContract.go
    └── MyContract.sol
```

### PluginDescriptor

Every plugin must export a `PluginDescriptor` variable in `plugin.go`:

```go
var PluginDescriptor = scenario.PluginDescriptor{
    Name:        "plugin-name",        // Unique identifier (required)
    Description: "What this plugin does",
    Categories: []*scenario.Category{  // Group scenarios by category
        {
            Name:        "Category Name",
            Description: "Category description",
            Descriptors: []*scenario.Descriptor{
                &scenario1.ScenarioDescriptor,
            },
        },
    },
}
```

## Creating a Scenario

Each scenario must:
1. Export a `ScenarioDescriptor` variable of type `scenario.Descriptor`
2. Implement the `scenario.Scenario` interface

### ScenarioDescriptor

```go
var ScenarioDescriptor = scenario.Descriptor{
    Name:           "scenario-name",     // Unique identifier (required)
    Description:    "What it does",      // Shown in CLI help
    Aliases:        []string{"alias1"},  // Alternative names (optional)
    DefaultOptions: DefaultOptions,      // Default config struct
    NewScenario:    newScenario,         // Factory function (required)
}
```

### Scenario Interface

```go
type Scenario interface {
    Flags(flags *pflag.FlagSet) error  // Register CLI flags
    Init(options *Options) error        // Initialize with wallet pool
    Run(ctx context.Context) error      // Execute the scenario
}
```

### Configuration with YAML

Options are automatically parsed from YAML configuration. Use struct tags:

```go
type Options struct {
    Throughput  uint64  `yaml:"throughput"`
    TotalCount  uint64  `yaml:"total_count"`
    MaxPending  uint64  `yaml:"max_pending"`
    BaseFee     float64 `yaml:"base_fee"`
}
```

### Transaction Pattern

For sending transactions, use `scenario.RunTransactionScenario`:

```go
func (s *Scenario) Run(ctx context.Context) error {
    return scenario.RunTransactionScenario(ctx, s.logger, &scenario.TransactionScenarioOptions{
        WalletPool:  s.walletPool,
        Throughput:  s.options.Throughput,
        TotalCount:  s.options.TotalCount,
        MaxPending:  s.options.MaxPending,
        Rebroadcast: s.options.Rebroadcast,
    }, s.buildTransaction)
}

func (s *Scenario) buildTransaction(
    ctx context.Context,
    txIdx uint64,
    wallet *spamoor.Wallet,
    client *spamoor.Client,
    txNonce uint64,
    onComplete scenario.TransactionCompleteCallback,
) (*types.Transaction, error) {
    defer onComplete()  // ALWAYS call onComplete!

    // Build and return your transaction
    tx, err := txbuilder.BuildDynamicFeeTx(&txbuilder.TxMetadata{
        // ...
    })
    return tx, err
}
```

## Available Packages

The following packages have pre-extracted symbols for use in plugins:

### Spamoor Packages

| Package   | Import Path                                | Purpose                                    |
| --------- | ------------------------------------------ | ------------------------------------------ |
| scenario  | `github.com/ethpandaops/spamoor/scenario`  | Core scenario interfaces and helpers       |
| spamoor   | `github.com/ethpandaops/spamoor/spamoor`   | Client pool, wallet pool, transaction pool |
| txbuilder | `github.com/ethpandaops/spamoor/txbuilder` | Transaction building utilities             |
| utils     | `github.com/ethpandaops/spamoor/utils`     | Common utilities                           |

### Third-party Packages

| Package | Import Path                                         | Purpose                        |
| ------- | --------------------------------------------------- | ------------------------------ |
| logrus  | `github.com/sirupsen/logrus`                        | Structured logging             |
| pflag   | `github.com/spf13/pflag`                            | CLI flag parsing               |
| uint256 | `github.com/holiman/uint256`                        | Big integer operations         |
| common  | `github.com/ethereum/go-ethereum/common`            | Ethereum types (Address, Hash) |
| types   | `github.com/ethereum/go-ethereum/core/types`        | Transaction types              |
| abi     | `github.com/ethereum/go-ethereum/accounts/abi`      | ABI encoding/decoding          |
| bind    | `github.com/ethereum/go-ethereum/accounts/abi/bind` | Contract bindings              |
| crypto  | `github.com/ethereum/go-ethereum/crypto`            | Cryptographic functions        |
| yaml    | `gopkg.in/yaml.v3`                                  | YAML parsing                   |

See `plugin/symbols/` for the complete list of extracted symbols.

## CLI Usage

### Load from Local Directory (Development)

```bash
spamoor --plugin ./my-plugin my-scenario -h http://localhost:8545
```

### Load from Archive

```bash
spamoor --plugin ./my-plugin.tar.gz my-scenario -h http://localhost:8545
```

### Load from URL

```bash
spamoor --plugin https://example.com/my-plugin.tar.gz my-scenario -h http://localhost:8545
```

### Multiple Plugins

```bash
spamoor --plugin ./plugin1 --plugin ./plugin2 some-scenario -h http://localhost:8545
```

### With YAML Configuration

```bash
spamoor run config.yaml --plugin ./my-plugin
```

## Daemon Mode

In daemon mode, plugins can be managed via the web UI (`/plugins` page) or REST API.

### Register via API

```bash
# From URL
curl -X POST http://localhost:8080/api/plugins \
  -F "type=url" \
  -F "path=https://example.com/my-plugin.tar.gz"

# From local path
curl -X POST http://localhost:8080/api/plugins \
  -F "type=local" \
  -F "path=/path/to/my-plugin"

# Upload archive
curl -X POST http://localhost:8080/api/plugins \
  -F "type=upload" \
  -F "file=@my-plugin.tar.gz"
```

### Building Plugin Archives

Use `make plugins` to build distributable archives:

```bash
make plugins
# Creates plugins/<plugin-name>.tar.gz with auto-generated plugin.yaml
```

## Validating Plugins

Validate plugins before deployment:

```bash
# Validate directory
spamoor-utils validate-plugin ./my-plugin

# Validate archive
spamoor-utils validate-plugin ./my-plugin.tar.gz
```

This checks:
1. Plugin loads successfully via Yaegi
2. `PluginDescriptor` is properly defined
3. All scenarios can be instantiated
4. Flags are registered correctly

## Yaegi Limitations

Yaegi is a Go interpreter with some limitations:

### Unsupported Features

- **CGO**: Cannot use packages requiring CGO
- **go:embed**: Embed directives not supported
- **Complex generics**: Keep generics usage simple
- **Reflection on unexported fields**: May fail

### Known Issues

- Sending to channels with type aliases inside closures may panic
- Some complex type assertions may fail
- Pointer vs value receiver types must match exactly

### Performance

Interpreted code is slower than compiled code. For high-throughput scenarios, consider contributing to the core codebase instead.

## Troubleshooting

### "undefined: package.Symbol"

The package needs symbols extracted. Add it to `plugin/symbols/generate.go` and run:
```bash
go generate ./plugin/symbols/...
```

### "PluginDescriptor not found"

Ensure your `plugin.go` exports a variable named exactly `PluginDescriptor` of type `scenario.PluginDescriptor`.

### "interpreter panic"

Simplify your code:
- Use explicit function types instead of closures where possible
- Avoid complex type assertions
- Use simpler channel patterns

### Scenario Not Found

Verify:
1. The `ScenarioDescriptor` is referenced in `PluginDescriptor.Categories`
2. The `Name` field matches what you're trying to run
3. The plugin loaded without errors

### Flags Not Working

Ensure your `Flags()` method:
- Returns `nil` on success
- Registers flags on the provided `*pflag.FlagSet`
- Uses unique flag names (no conflicts with global flags)

### Transactions Not Counted

Always call `onComplete()` in your `ProcessNextTxFn`:
```go
func (s *Scenario) buildTx(..., onComplete scenario.TransactionCompleteCallback) (*types.Transaction, error) {
    defer onComplete()  // Required!
    // ...
}
```

## Example Plugin

See `plugins/_example-plugin/` for a complete working example demonstrating:

- Plugin and scenario structure
- Contract deployment with well-known wallets
- Contract interactions using `BuildBoundTx`
- Proper `onComplete` handling
- Configuration via flags and YAML
