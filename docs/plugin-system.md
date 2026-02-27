# Spamoor Plugin System Guide

This guide covers the complete plugin system for Spamoor: how to develop plugins, distribute them, use them from the CLI and daemon, and manage them via the REST API. Plugins are Go source files interpreted at runtime by [Yaegi](https://github.com/traefik/yaegi), allowing you to add custom transaction scenarios without recompiling Spamoor.

## Table of Contents

- [Overview](#overview)
- [Quick Start](#quick-start)
- [Plugin Development](#plugin-development)
  - [Plugin Structure](#plugin-structure)
  - [PluginDescriptor](#plugindescriptor)
  - [Creating a Scenario](#creating-a-scenario)
  - [Scenario Interface](#scenario-interface)
  - [Configuration with YAML and Flags](#configuration-with-yaml-and-flags)
  - [Wallet Management](#wallet-management)
  - [Transaction Patterns](#transaction-patterns)
  - [Contract Interactions](#contract-interactions)
  - [Available Packages](#available-packages)
- [Building and Distributing Plugins](#building-and-distributing-plugins)
  - [Plugin Archives](#plugin-archives)
  - [Plugin Metadata (plugin.yaml)](#plugin-metadata-pluginyaml)
  - [Building Archives with Make](#building-archives-with-make)
- [Validating Plugins](#validating-plugins)
- [CLI Usage](#cli-usage)
  - [Single Scenario Mode](#single-scenario-mode)
  - [Multi-Scenario Mode (run command)](#multi-scenario-mode-run-command)
  - [Plugin Source Types](#plugin-source-types)
- [Daemon Mode](#daemon-mode)
  - [Loading Plugins on Startup](#loading-plugins-on-startup)
  - [Plugin Persistence](#plugin-persistence)
  - [Plugin Lifecycle](#plugin-lifecycle)
  - [Web UI](#web-ui)
  - [Disabling the Plugin API](#disabling-the-plugin-api)
- [REST API Reference](#rest-api-reference)
  - [List All Plugins](#list-all-plugins)
  - [Get Plugin Details](#get-plugin-details)
  - [Register a Plugin](#register-a-plugin)
  - [Delete a Plugin](#delete-a-plugin)
  - [Reload a Plugin](#reload-a-plugin)
- [Architecture](#architecture)
  - [Plugin Loading Pipeline](#plugin-loading-pipeline)
  - [Scenario Registry Integration](#scenario-registry-integration)
  - [Category Merging](#category-merging)
  - [Reference Counting and Cleanup](#reference-counting-and-cleanup)
- [Yaegi Limitations](#yaegi-limitations)
- [Troubleshooting](#troubleshooting)


## Overview

Plugins provide a way to extend Spamoor without modifying the core codebase or recompiling. They are useful for:

- **Rapid prototyping**: Test new transaction patterns quickly without a full build cycle
- **Custom deployments**: Add organization-specific scenarios
- **External contributions**: Share scenarios without core integration
- **Runtime management**: Add, reload, and remove plugins while the daemon is running

Plugins are interpreted at runtime by Yaegi and have access to Spamoor's APIs through pre-extracted symbols. Each plugin can contain one or more transaction scenarios organized into categories.


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
    "fmt"

    "github.com/ethpandaops/spamoor/scenario"
    "github.com/ethpandaops/spamoor/spamoor"
    "github.com/sirupsen/logrus"
    "github.com/spf13/pflag"
)

type Options struct {
    TotalCount uint64 `yaml:"total_count"`
    Throughput uint64 `yaml:"throughput"`
    MaxPending uint64 `yaml:"max_pending"`
}

var DefaultOptions = Options{
    Throughput: 10,
}

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
    flags.Uint64VarP(&s.options.TotalCount, "count", "c", DefaultOptions.TotalCount, "Total transactions to send")
    flags.Uint64VarP(&s.options.Throughput, "throughput", "t", DefaultOptions.Throughput, "Transactions per slot")
    flags.Uint64Var(&s.options.MaxPending, "max-pending", DefaultOptions.MaxPending, "Max pending transactions")
    return nil
}

func (s *Scenario) Init(opts *scenario.Options) error {
    s.walletPool = opts.WalletPool

    if opts.Config != "" {
        err := scenario.ParseAndValidateConfig(&ScenarioDescriptor, opts.Config, &s.options, s.logger)
        if err != nil {
            return err
        }
    }

    s.walletPool.SetWalletCount(s.options.Throughput * 10)

    if s.options.TotalCount == 0 && s.options.Throughput == 0 {
        return fmt.Errorf("neither total count nor throughput limit set")
    }

    return nil
}

func (s *Scenario) Run(ctx context.Context) error {
    s.logger.Info("starting my-scenario")
    // Your transaction logic here
    return nil
}
```

4. Validate and run:
```bash
# Validate the plugin
spamoor-utils validate-plugin ./my-plugin

# Run a scenario from the plugin
spamoor --plugin ./my-plugin my-scenario \
  -h http://localhost:8545 \
  -p YOUR_PRIVATE_KEY \
  -t 5
```


## Plugin Development

### Plugin Structure

A plugin is a directory of Go source files with a required entry point (`plugin.go`) and one or more scenario sub-packages:

```
my-plugin/
├── plugin.go              # Required: exports PluginDescriptor
├── scenario1/
│   └── scenario1.go       # Scenario implementation
├── scenario2/
│   └── scenario2.go       # Another scenario
└── contract/              # Optional: generated contract bindings
    ├── MyContract.sol
    ├── MyContract.go      # Generated with abigen
    └── compile.sh
```

**Key rules:**
- The root `plugin.go` must be in a package named `plugin`
- The root `plugin.go` must export a variable named exactly `PluginDescriptor`
- Each scenario should be in its own sub-package
- Import paths use `github.com/ethpandaops/spamoor/plugins/<plugin-name>/...`

### PluginDescriptor

Every plugin must export a `PluginDescriptor` variable of type `scenario.PluginDescriptor` in `plugin.go`:

```go
var PluginDescriptor = scenario.PluginDescriptor{
    Name:        "plugin-name",        // Unique identifier (required)
    Description: "What this plugin does",
    Categories: []*scenario.Category{  // Organize scenarios by category
        {
            Name:        "Category Name",
            Description: "Category description",
            Descriptors: []*scenario.Descriptor{
                &scenario1.ScenarioDescriptor,
            },
            Children: []*scenario.Category{  // Optional nested categories
                // ...
            },
        },
    },
}
```

Plugin categories are merged with Spamoor's native categories. If a plugin category name matches an existing native category (e.g. "Simple"), the plugin's scenarios are added to that category. New category names create new top-level categories.

### Creating a Scenario

Each scenario must:
1. Export a `ScenarioDescriptor` variable of type `scenario.Descriptor`
2. Implement the `scenario.Scenario` interface

```go
var ScenarioDescriptor = scenario.Descriptor{
    Name:           "scenario-name",     // Unique name (required)
    Description:    "What it does",      // Shown in CLI help
    Aliases:        []string{"alias1"},  // Alternative names (optional)
    DefaultOptions: DefaultOptions,      // Default config struct (for YAML merging)
    NewScenario:    newScenario,         // Factory function (required)
}
```

**Important:** Plugin scenarios cannot override built-in (native) scenarios. If a plugin scenario has the same name as a native scenario, registration will fail.

### Scenario Interface

```go
type Scenario interface {
    // Flags registers CLI flags for this scenario.
    Flags(flags *pflag.FlagSet) error

    // Init initializes the scenario with wallet pool and configuration.
    Init(options *Options) error

    // Run executes the scenario. Must respect context cancellation.
    Run(ctx context.Context) error
}
```

The `Options` struct passed to `Init` contains:
- `WalletPool`: The wallet pool for managing child wallets
- `Config`: YAML configuration string (from daemon or `run` command)
- `GlobalCfg`: Global configuration map
- `PluginPath`: Path to the plugin's resources directory (empty for native scenarios)

### Configuration with YAML and Flags

Scenarios support dual configuration: CLI flags for single-scenario mode and YAML for daemon/multi-scenario mode.

**Define options with YAML struct tags:**
```go
type ScenarioOptions struct {
    TotalCount  uint64  `yaml:"total_count"`
    Throughput  uint64  `yaml:"throughput"`
    MaxPending  uint64  `yaml:"max_pending"`
    MaxWallets  uint64  `yaml:"max_wallets"`
    Rebroadcast uint64  `yaml:"rebroadcast"`
    BaseFee     float64 `yaml:"base_fee"`
    TipFee      float64 `yaml:"tip_fee"`
    BaseFeeWei  string  `yaml:"base_fee_wei"`
    TipFeeWei   string  `yaml:"tip_fee_wei"`
    ClientGroup string  `yaml:"client_group"`
    Timeout     string  `yaml:"timeout"`
    LogTxs      bool    `yaml:"log_txs"`
}
```

**Register corresponding CLI flags:**
```go
func (s *Scenario) Flags(flags *pflag.FlagSet) error {
    flags.Uint64VarP(&s.options.TotalCount, "count", "c", DefaultOptions.TotalCount, "Total transactions")
    flags.Uint64VarP(&s.options.Throughput, "throughput", "t", DefaultOptions.Throughput, "Transactions per slot")
    flags.Uint64Var(&s.options.MaxPending, "max-pending", DefaultOptions.MaxPending, "Max pending txs")
    flags.Float64Var(&s.options.BaseFee, "basefee", DefaultOptions.BaseFee, "Max fee per gas (gwei)")
    flags.Float64Var(&s.options.TipFee, "tipfee", DefaultOptions.TipFee, "Max tip per gas (gwei)")
    flags.StringVar(&s.options.BaseFeeWei, "basefee-wei", "", "Max fee per gas in wei (overrides --basefee)")
    flags.StringVar(&s.options.TipFeeWei, "tipfee-wei", "", "Max tip per gas in wei (overrides --tipfee)")
    return nil
}
```

**Parse YAML config in Init:**
```go
func (s *Scenario) Init(opts *scenario.Options) error {
    s.walletPool = opts.WalletPool

    if opts.Config != "" {
        err := scenario.ParseAndValidateConfig(&ScenarioDescriptor, opts.Config, &s.options, s.logger)
        if err != nil {
            return err
        }
    }

    // Configure wallet pool, validate options, etc.
    return nil
}
```

### Wallet Management

Scenarios receive a `WalletPool` via `Init`. Configure it based on your transaction volume:

```go
func (s *Scenario) Init(opts *scenario.Options) error {
    s.walletPool = opts.WalletPool

    // Set wallet count based on throughput
    if s.options.MaxWallets > 0 {
        s.walletPool.SetWalletCount(s.options.MaxWallets)
    } else if s.options.Throughput*10 < 1000 {
        s.walletPool.SetWalletCount(s.options.Throughput * 10)
    } else {
        s.walletPool.SetWalletCount(1000)
    }

    // Add well-known wallets for special purposes (e.g. contract deployment)
    s.walletPool.AddWellKnownWallet(&spamoor.WellKnownWalletConfig{
        Name:          "deployer",
        RefillAmount:  utils.EtherToWei(uint256.NewInt(10)),
        RefillBalance: utils.EtherToWei(uint256.NewInt(5)),
        VeryWellKnown: false,  // false = scenario-specific, true = shared across scenarios
    })

    return nil
}
```

**Critical rules:**
- Never use the root wallet directly in scenarios — it causes nonce conflicts
- Always spread transactions across multiple wallets
- Use well-known wallets for contract deployments

### Transaction Patterns

Use `scenario.RunTransactionScenario` for standardized transaction loops with rate limiting and progress tracking:

```go
func (s *Scenario) Run(ctx context.Context) error {
    maxPending := s.options.MaxPending
    if maxPending == 0 {
        maxPending = s.options.Throughput * 10
    }

    return scenario.RunTransactionScenario(ctx, scenario.TransactionScenarioOptions{
        TotalCount: s.options.TotalCount,
        Throughput: s.options.Throughput,
        MaxPending: maxPending,
        WalletPool: s.walletPool,
        Logger:     s.logger,
        ProcessNextTxFn: func(ctx context.Context, params *scenario.ProcessNextTxParams) error {
            // Build and send your transaction
            receiptChan, tx, client, wallet, err := s.sendTx(ctx, params.TxIdx)

            params.NotifySubmitted()
            params.OrderedLogCb(func() {
                if err != nil {
                    s.logger.Warnf("tx failed: %v", err)
                } else {
                    s.logger.Debugf("sent tx #%d: %v", params.TxIdx+1, tx.Hash().String())
                }
            })

            if _, err := receiptChan.Wait(ctx); err != nil {
                return err
            }
            return err
        },
    })
}
```

### Contract Interactions

Use the `BuildBoundTx` pattern for contract interactions. This is required because go-ethereum's standard bound contract methods must not be used directly for transaction submission in plugins.

**Deploy a contract:**
```go
func (s *Scenario) deployContract(ctx context.Context) (*types.Receipt, error) {
    deployerWallet := s.walletPool.GetWellKnownWallet("deployer")
    client := s.walletPool.GetClient(
        spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0),
    )

    baseFeeWei, tipFeeWei := spamoor.ResolveFees(s.options.BaseFee, s.options.TipFee, s.options.BaseFeeWei, s.options.TipFeeWei)
    feeCap, tipCap, _ := s.walletPool.GetTxPool().GetSuggestedFees(client, baseFeeWei, tipFeeWei)

    deploymentTx, err := deployerWallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
        GasFeeCap: uint256.MustFromBig(feeCap),
        GasTipCap: uint256.MustFromBig(tipCap),
        Gas:       2000000,
        Value:     uint256.NewInt(0),
    }, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
        _, deployTx, _, err := contract.DeployMyContract(transactOpts, client.GetEthClient())
        return deployTx, err
    })
    if err != nil {
        return nil, err
    }

    receipt, err := s.walletPool.GetTxPool().SendAndAwaitTransaction(ctx, deployerWallet, deploymentTx, &spamoor.SendTransactionOptions{
        Client:      client,
        Rebroadcast: true,
    })
    return receipt, err
}
```

**Call a contract method:**
```go
contractInstance, _ := contract.NewMyContract(contractAddr, client.GetEthClient())

tx, err := wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
    GasFeeCap: uint256.MustFromBig(feeCap),
    GasTipCap: uint256.MustFromBig(tipCap),
    Gas:       100000,
    Value:     uint256.NewInt(0),
}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
    return contractInstance.MyMethod(transactOpts, arg1, arg2)
})
```

### Available Packages

Plugins can import the following pre-extracted packages:

**Spamoor packages:**

| Package   | Import Path                                | Purpose                                    |
| --------- | ------------------------------------------ | ------------------------------------------ |
| scenario  | `github.com/ethpandaops/spamoor/scenario`  | Core scenario interfaces and helpers       |
| spamoor   | `github.com/ethpandaops/spamoor/spamoor`   | Client pool, wallet pool, transaction pool |
| txbuilder | `github.com/ethpandaops/spamoor/txbuilder` | Transaction building utilities             |
| utils     | `github.com/ethpandaops/spamoor/utils`     | Common utilities (e.g. EtherToWei)         |

**Ethereum packages:**

| Package | Import Path                                              | Purpose                        |
| ------- | -------------------------------------------------------- | ------------------------------ |
| common  | `github.com/ethereum/go-ethereum/common`                 | Ethereum types (Address, Hash) |
| types   | `github.com/ethereum/go-ethereum/core/types`             | Transaction and receipt types  |
| abi     | `github.com/ethereum/go-ethereum/accounts/abi`           | ABI encoding/decoding          |
| bind    | `github.com/ethereum/go-ethereum/accounts/abi/bind`      | Contract bindings              |
| bind/v2 | `github.com/ethereum/go-ethereum/accounts/abi/bind/v2`   | Contract bindings v2           |
| crypto  | `github.com/ethereum/go-ethereum/crypto`                 | Cryptographic functions        |
| event   | `github.com/ethereum/go-ethereum/event`                  | Event subscriptions            |
| -       | `github.com/ethereum/go-ethereum`                        | Core ethereum interfaces       |

**Third-party packages:**

| Package | Import Path                  | Purpose                |
| ------- | ---------------------------- | ---------------------- |
| logrus  | `github.com/sirupsen/logrus` | Structured logging     |
| pflag   | `github.com/spf13/pflag`    | CLI flag parsing       |
| uint256 | `github.com/holiman/uint256` | Big integer operations |
| yaml    | `gopkg.in/yaml.v3`          | YAML parsing           |

**Standard library:** All Go standard library packages are available (e.g. `fmt`, `context`, `time`, `math/big`, `math/rand`).

If you need a package that is not listed here, symbols must be extracted first. See [Adding New Package Symbols](#adding-new-package-symbols) in the troubleshooting section.


## Building and Distributing Plugins

### Plugin Archives

Plugins can be distributed as `.tar.gz` archives. The archive must contain:
- `plugin.yaml` at the root (auto-generated by `make plugins`)
- `plugin.go` and all Go source files
- Any additional resources (contract ABIs, etc.)

Archive structure:
```
.
├── plugin.yaml        # Metadata (auto-generated)
├── plugin.go          # Plugin entry point
└── scenario1/
    ├── scenario1.go
    └── contract/
        └── MyContract.go
```

### Plugin Metadata (plugin.yaml)

The `plugin.yaml` file contains build metadata:

```yaml
name: my-plugin
build_time: "2024-01-01T00:00:00Z"
git_version: "abc1234"
```

- **Required for archives** (`.tar.gz` files): Must be present at the archive root
- **Not required for local directories**: The directory name is used as the plugin name
- **Auto-generated** by `make plugins`: You don't need to create it manually

### Building Archives with Make

Build all plugins in the `plugins/` directory:

```bash
make plugins
```

This process:
1. Iterates over every subdirectory in `plugins/`
2. Generates `plugin.yaml` with current build time and git version
3. Creates a `.tar.gz` archive for each plugin
4. Archives are saved as `plugins/<plugin-name>.tar.gz`

To build a single plugin manually:

```bash
cd plugins/my-plugin
# Create plugin.yaml
echo "name: my-plugin" > plugin.yaml
echo "build_time: $(date -u '+%Y-%m-%dT%H:%M:%SZ')" >> plugin.yaml
echo "git_version: $(git rev-parse --short HEAD)" >> plugin.yaml
# Create archive
tar -czf ../my-plugin.tar.gz .
# Clean up
rm plugin.yaml
```


## Validating Plugins

Use the `spamoor-utils validate-plugin` command to verify a plugin before deployment:

```bash
# Validate a local directory
spamoor-utils validate-plugin ./my-plugin

# Validate an archive
spamoor-utils validate-plugin ./my-plugin.tar.gz
```

The validation checks:
1. The plugin loads successfully via the Yaegi interpreter
2. `PluginDescriptor` has a non-empty `Name`
3. All scenarios in the descriptor:
   - Have a non-empty `Name`
   - Have a non-nil `NewScenario` factory function
   - Can be instantiated successfully
   - Register CLI flags without errors

Example output:
```
Validating plugin: ./my-plugin

✓ Plugin loaded successfully

Plugin Descriptor:
  Name:        my-plugin
  Description: My custom plugin
  Categories:  1
  Scenarios:   2

Scenarios:
  [0] my-scenario-1
      ✓ NewScenario defined
      ✓ Instance created successfully
      ✓ Flags registered: 8
      Description: First scenario

  [1] my-scenario-2
      ✓ NewScenario defined
      ✓ Instance created successfully
      ✓ Flags registered: 5
      Description: Second scenario

Result: ✓ Plugin 'my-plugin' is valid with 2 scenario(s)
```


## CLI Usage

### Single Scenario Mode

Run a specific scenario from a plugin:

```bash
# Load from local directory (best for development)
spamoor --plugin ./my-plugin my-scenario \
  -h http://localhost:8545 \
  -p YOUR_PRIVATE_KEY \
  -t 10 -c 100

# Load from tar.gz archive
spamoor --plugin ./my-plugin.tar.gz my-scenario \
  -h http://localhost:8545 \
  -p YOUR_PRIVATE_KEY

# Load from URL
spamoor --plugin https://example.com/my-plugin.tar.gz my-scenario \
  -h http://localhost:8545 \
  -p YOUR_PRIVATE_KEY

# Load multiple plugins
spamoor --plugin ./plugin1 --plugin ./plugin2 some-scenario \
  -h http://localhost:8545 \
  -p YOUR_PRIVATE_KEY
```

The source type is auto-detected:
- URLs (starting with `http://` or `https://`) are downloaded
- Existing directories are loaded as local plugins
- Everything else is treated as a file path (`.tar.gz` archive)

### Multi-Scenario Mode (run command)

Run multiple scenarios from a YAML configuration file:

```bash
spamoor run config.yaml --plugin ./my-plugin \
  -h http://localhost:8545 \
  -p YOUR_PRIVATE_KEY
```

The YAML config can reference plugin scenarios just like native ones:

```yaml
- name: "Custom Plugin Spammer"
  scenario: my-scenario
  config:
    throughput: 20
    total_count: 500
    max_pending: 100

- name: "Native EOA Spammer"
  scenario: eoatx
  config:
    throughput: 10
```

### Plugin Source Types

| Source | CLI Flag | Description |
| --- | --- | --- |
| Local directory | `--plugin ./path/to/dir` | Symlinked into temp GOPATH (best for development) |
| Archive file | `--plugin ./plugin.tar.gz` | Extracted to temp directory (auto-detects gzip) |
| URL | `--plugin https://example.com/plugin.tar.gz` | Downloaded and extracted |


## Daemon Mode

In daemon mode (`spamoor-daemon`), plugins are persisted in the database and managed at runtime.

### Loading Plugins on Startup

Plugins specified via `--plugin` flags are loaded and persisted to the database:

```bash
spamoor-daemon \
  --plugin ./my-plugin \
  --plugin https://example.com/another-plugin.tar.gz \
  -h http://localhost:8545 \
  -p YOUR_PRIVATE_KEY
```

On subsequent startups, previously registered plugins are automatically restored from the database. CLI-specified plugins are also re-registered (updating if already present).

### Plugin Persistence

The daemon stores plugins in its SQLite database with the following behavior:

| Source Type | What's Stored | On Restore | On Reload |
| --- | --- | --- | --- |
| **URL** | Base64-encoded archive + source URL | Loaded from stored archive | Re-downloaded from URL |
| **File upload** | Base64-encoded archive | Loaded from stored archive | Not supported |
| **File path** | Base64-encoded archive | Loaded from stored archive | Not supported |
| **Local directory** | Directory path only (no archive) | Re-read from original path | Re-read from path |

This means:
- URL and uploaded plugins work even if the original source is unavailable (archive is stored)
- Local directory plugins always pick up the latest source changes on restart or reload
- Reloading is only supported for URL and local directory plugins

### Plugin Lifecycle

1. **Register**: Plugin is loaded, scenarios are registered, and metadata is saved to the database
2. **Active**: Plugin scenarios are available for creating spammers
3. **Update/Reload**: A new version replaces the old; the old version is marked as "deprecated" if it has running spammers
4. **Deprecated**: Old plugin version kept alive until all spammers using it finish
5. **Cleanup**: Once no spammers reference the plugin, its temporary files are removed
6. **Delete**: Plugin is removed from the database and all its scenarios are unregistered (fails if spammers are still running)

### Web UI

The daemon provides a plugin management page at `/plugins` in the web UI, accessible from the navigation header. From this page you can:
- View all registered plugins and their status
- Register new plugins from URL, local path, or file upload
- Reload plugins from their original source
- Delete plugins
- See which scenarios each plugin provides
- Monitor running spammer counts per plugin

### Disabling the Plugin API

To prevent runtime plugin management while still allowing `--plugin` loading at startup:

```bash
spamoor-daemon --disable-plugin-api \
  --plugin ./my-plugin \
  -h http://localhost:8545 \
  -p YOUR_PRIVATE_KEY
```

This disables the `POST /api/plugins`, `DELETE /api/plugins/{name}`, and `POST /api/plugins/{name}/reload` endpoints. The `GET` endpoints remain available.


## REST API Reference

All plugin API endpoints require authentication when `--enable-auth` is set. Pass the auth token via the `Authorization` header.

### List All Plugins

```http
GET /api/plugins
```

Returns all registered plugins including deprecated ones.

**Response:**
```json
[
  {
    "name": "my-plugin",
    "source_type": "url",
    "source_path": "https://example.com/my-plugin.tar.gz",
    "metadata_name": "my-plugin",
    "metadata_build_time": "2024-01-01T00:00:00Z",
    "metadata_git_version": "abc1234",
    "scenarios": ["my-scenario-1", "my-scenario-2"],
    "enabled": true,
    "load_error": "",
    "running_count": 1,
    "is_loaded": true,
    "deprecated": false,
    "created_at": 1704067200,
    "updated_at": 1704067200
  }
]
```

**curl example:**
```bash
curl http://localhost:8080/api/plugins
```

### Get Plugin Details

```http
GET /api/plugins/{name}
```

Returns details for a single plugin.

**curl example:**
```bash
curl http://localhost:8080/api/plugins/my-plugin
```

**Response:** Same structure as a single entry from the list endpoint.

### Register a Plugin

```http
POST /api/plugins
```

Accepts both `multipart/form-data` and `application/json` content types.

**From URL (multipart):**
```bash
curl -X POST http://localhost:8080/api/plugins \
  -F "type=url" \
  -F "path=https://example.com/my-plugin.tar.gz"
```

**From URL (JSON):**
```bash
curl -X POST http://localhost:8080/api/plugins \
  -H "Content-Type: application/json" \
  -d '{"type": "url", "path": "https://example.com/my-plugin.tar.gz"}'
```

**From local path (multipart):**
```bash
curl -X POST http://localhost:8080/api/plugins \
  -F "type=local" \
  -F "path=/path/to/my-plugin"
```

**From local path (JSON):**
```bash
curl -X POST http://localhost:8080/api/plugins \
  -H "Content-Type: application/json" \
  -d '{"type": "local", "path": "/path/to/my-plugin"}'
```

**Upload archive (multipart only):**
```bash
curl -X POST http://localhost:8080/api/plugins \
  -F "type=upload" \
  -F "file=@my-plugin.tar.gz"
```

**Response:**
```json
{
  "name": "my-plugin",
  "scenarios": ["my-scenario-1", "my-scenario-2"]
}
```

If a plugin with the same name already exists, it is replaced. The old version is marked as deprecated if it has running spammers.

### Delete a Plugin

```http
DELETE /api/plugins/{name}
```

Removes a plugin and unregisters all its scenarios.

**curl example:**
```bash
curl -X DELETE http://localhost:8080/api/plugins/my-plugin
```

**Fails with 400** if the plugin has running spammers. Stop all spammers using the plugin's scenarios before deleting.

### Reload a Plugin

```http
POST /api/plugins/{name}/reload
```

Re-loads a plugin from its original source.

**curl example:**
```bash
curl -X POST http://localhost:8080/api/plugins/my-plugin/reload
```

**Response:**
```json
{
  "name": "my-plugin",
  "scenarios": ["my-scenario-1", "my-scenario-2"]
}
```

**Constraints:**
- Only supported for `url` and `local` source types (not `upload` or `file`)
- Fails if the plugin has running spammers


## Architecture

### Plugin Loading Pipeline

When a plugin is loaded (from any source), the following steps occur:

1. **Source resolution**: Download URL, read file, or symlink directory
2. **Metadata extraction**: Parse `plugin.yaml` from the archive (or derive from directory name for local plugins)
3. **Temp directory setup**: Create a temporary GOPATH structure:
   ```
   <tmpdir>/src/github.com/ethpandaops/spamoor/plugins/<plugin-name>/
   ```
4. **Yaegi interpreter**: Create a new interpreter instance with:
   - The temp directory as GOPATH
   - Pre-extracted symbols for all available packages
   - A `SymlinkFS` that follows symlinks (needed for local plugins)
5. **Plugin import**: Import the plugin package and extract the `PluginDescriptor` variable
6. **LoadedPlugin creation**: Track the plugin with its descriptor, metadata, temp directory, and source info

### Scenario Registry Integration

Spamoor maintains two registries:
- **PluginRegistry**: Tracks loaded plugins and their lifecycle
- **ScenarioRegistry**: Maps scenario names to descriptors (both native and plugin)

When plugin scenarios are registered:
- Each scenario is checked against native scenarios — **plugin scenarios cannot override native ones**
- If a plugin scenario replaces one from a different plugin, the old plugin's scenario is unregistered
- The scenario entry tracks its source (`Native` or `Plugin`) and a reference to the `LoadedPlugin`

### Category Merging

Plugin categories are merged with native categories when listing available scenarios:
- If a plugin category name matches a native category name, plugin scenarios are added to that category
- New category names create new top-level entries
- Category descriptions from plugins update existing category descriptions
- Nested categories (children) are merged recursively

### Reference Counting and Cleanup

Plugins use reference counting to ensure safe cleanup:

1. **Running counter**: Each `LoadedPlugin` tracks how many spammers are actively using its scenarios via `AddRunning()` / `RemoveRunning()`
2. **Deprecation**: When a plugin is replaced by a new version, the old version is marked as deprecated if it still has running spammers
3. **Cleanup guard**: A plugin's temporary files are only removed when:
   - All its scenarios have been unregistered, AND
   - No spammers are running with its scenarios
4. **Automatic cleanup**: The daemon registers a cleanup callback that removes temp directories when plugins become eligible


## Yaegi Limitations

Yaegi is a Go interpreter with some limitations compared to compiled Go:

### Unsupported Features

- **CGO**: Cannot use packages that require CGO (e.g. SQLite bindings)
- **go:embed**: Embed directives are not supported
- **Complex generics**: Keep generics usage simple
- **Reflection on unexported fields**: May fail or behave unexpectedly

### Known Issues

- Sending to channels with type aliases inside closures may panic
- Some complex type assertions may fail — use type switches instead
- Pointer vs value receiver types must match exactly
- Complex channel patterns in closures may not work correctly

### Performance

Interpreted code is slower than compiled code. For very high-throughput scenarios where performance is critical, consider contributing the scenario to the core codebase instead.


## Troubleshooting

### "undefined: package.Symbol"

The package needs its symbols extracted for Yaegi. Add the package to `plugin/symbols/generate.go` and run:

```bash
go generate ./plugin/symbols/...
```

<a id="adding-new-package-symbols"></a>
To add a new package, add a `//go:generate` directive in `plugin/symbols/generate.go`:

```go
//go:generate yaegi extract github.com/some/package
```

Then run `go generate ./plugin/symbols/...` and commit the generated symbol file.

### "PluginDescriptor not found"

Ensure your `plugin.go`:
- Is in a package named `plugin`
- Exports a variable named exactly `PluginDescriptor`
- The variable is of type `scenario.PluginDescriptor`

### "plugin.yaml not found"

This only applies to `.tar.gz` archives. Either:
- Use `make plugins` to build archives (auto-generates `plugin.yaml`)
- Create `plugin.yaml` manually at the archive root

Local directory plugins do not need `plugin.yaml`.

### "interpreter panic"

Simplify your code:
- Use explicit function types instead of closures where possible
- Avoid complex type assertions — use type switches
- Use simpler channel patterns
- Avoid deeply nested generic types

The error message may include a hint about missing symbols with a suggestion to run `yaegi extract`.

### Scenario Not Found After Loading Plugin

Verify:
1. The `ScenarioDescriptor` is referenced in a category within `PluginDescriptor.Categories`
2. The scenario `Name` matches what you're trying to run
3. The plugin loaded without errors (check logs)
4. The scenario name doesn't conflict with a native scenario

### Flags Not Working

Ensure your `Flags()` method:
- Returns `nil` on success
- Registers flags on the provided `*pflag.FlagSet` (not a new one)
- Uses unique flag names that don't conflict with global flags (`--verbose`, `--trace`, `--rpchost`, `--privkey`, `--plugin`, etc.)

### Cannot Delete or Reload Plugin

Both operations fail if the plugin has running spammers. Stop all spammers using the plugin's scenarios first, then retry the operation.

### Plugin Works Locally but Not as Archive

Check that:
- All Go source files are included in the archive
- The archive structure starts at the plugin root (not a parent directory)
- `plugin.yaml` is at the root of the archive
- Generated contract bindings (`.go` files) are included (not just `.sol` files)

### Example Plugin

See `plugins/_example-plugin/` for a complete working example that demonstrates:

- Plugin structure with `PluginDescriptor` and categories
- Scenario implementation with full configuration
- Solidity contract compilation and abigen binding
- Contract deployment using well-known wallets and `BuildBoundTx`
- Contract method calls with `BuildBoundTx`
- Proper transaction submission with `OnComplete` and `OnConfirm` callbacks
- Fee resolution supporting both gwei and wei
- Nonce management with `MarkSkippedNonce`
- Wallet pool sizing based on transaction volume
