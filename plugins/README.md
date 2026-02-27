# Spamoor Plugins

This directory contains dynamically-loaded plugins that extend spamoor with custom transaction scenarios. Plugins are interpreted at runtime using the [Yaegi](https://github.com/traefik/yaegi) Go interpreter.

## Plugin Structure

A plugin is a Go package that exports a `PluginDescriptor` variable and contains one or more scenario implementations.

### Directory Layout

```
my-plugin/
├── plugin.go              # Main plugin file with PluginDescriptor
├── scenario1/
│   ├── scenario1.go       # First scenario implementation
│   └── contract/          # Optional: Contracts for this scenario
│       ├── MyContract.sol
│       ├── MyContract.go  # Generated with abigen
│       └── compile.sh
└── scenario2/
    └── scenario2.go       # Second scenario implementation
```

### Required Files

#### `plugin.go`

The main plugin file must export a `PluginDescriptor` variable:

```go
package plugin

import (
    "github.com/ethpandaops/spamoor/plugins/my-plugin/scenario1"
    "github.com/ethpandaops/spamoor/scenario"
)

// PluginDescriptor defines the plugin metadata and scenarios.
// This variable MUST be exported and named exactly "PluginDescriptor".
var PluginDescriptor = scenario.PluginDescriptor{
    Name:        "my-plugin",
    Description: "Description of what this plugin does",
    Scenarios: []*scenario.Descriptor{
        &scenario1.ScenarioDescriptor,
    },
}
```

#### `plugin.yaml` (auto-generated for archives)

When building plugin archives with `make plugins`, a `plugin.yaml` metadata file is automatically generated and included in the archive:

```yaml
name: my-plugin
build_time: "2024-01-01T00:00:00Z"
git_version: "v1.0.0"
```

> **Note:** You don't need to create `plugin.yaml` manually - it's auto-generated during `make plugins`. Local path plugins also don't require it - the directory name is used as the plugin name.

## Creating a Scenario

Scenarios in plugins follow the same structure as native scenarios. Each scenario must:

1. Export a `ScenarioDescriptor` variable of type `scenario.Descriptor`
2. Implement the `scenario.Scenario` interface (`Flags`, `Init`, `Run` methods)

See `docs/scenario-developers.md` for the complete scenario development guide, including:
- Wallet management patterns
- Transaction building with `BuildBoundTx`
- The `onComplete` callback pattern
- Configuration via flags and YAML

## Loading Plugins

### CLI Usage

```bash
# Load from local directory (development)
spamoor --plugin ./plugins/my-plugin eoatx -h http://localhost:8545

# Load from tar.gz file
spamoor --plugin ./my-plugin.tar.gz eoatx -h http://localhost:8545

# Load from URL
spamoor --plugin https://example.com/my-plugin.tar.gz eoatx -h http://localhost:8545
```

### Daemon Mode

Plugins can be registered via the WebUI (`/plugins` page) or API:

```bash
# Register from URL
curl -X POST http://localhost:8080/api/plugins \
  -F "type=url" \
  -F "path=https://example.com/my-plugin.tar.gz"

# Register from local path
curl -X POST http://localhost:8080/api/plugins \
  -F "type=local" \
  -F "path=/path/to/my-plugin"

# Upload archive
curl -X POST http://localhost:8080/api/plugins \
  -F "type=upload" \
  -F "file=@my-plugin.tar.gz"
```

## Creating a Plugin Archive

Use `make plugins` to build all plugin archives:

```bash
# Build all plugins in the plugins/ directory
make plugins

# Archives are created as plugins/<plugin-name>.tar.gz
```

The build process automatically:
- Generates `plugin.yaml` with build metadata (name, build time, git version)
- Creates the proper directory structure inside the archive
- Compresses the archive with gzip

The resulting archive structure:
```
.
├── plugin.yaml    # Auto-generated
├── plugin.go
└── scenario1/
    └── scenario1.go
```

## Validating Plugins

Use the validation tool to check plugin structure:

```bash
spamoor-utils validate-plugin ./plugins/my-plugin
spamoor-utils validate-plugin ./my-plugin.tar.gz
```

## Yaegi Limitations

The Yaegi interpreter has some limitations compared to compiled Go:

1. **No CGO support** - Cannot use packages that require CGO
2. **Symbol extraction required** - External packages must have symbols extracted via `yaegi extract`
3. **Some language features unsupported** - Complex type assertions, certain channel patterns in closures
4. **Performance** - Interpreted code is slower than compiled code

### Available Packages

The following packages are pre-extracted and available for plugins:

- Standard library (`fmt`, `context`, `time`, etc.)
- `github.com/ethereum/go-ethereum/*` (accounts, common, core, crypto, etc.)
- `github.com/ethpandaops/spamoor/scenario`
- `github.com/ethpandaops/spamoor/spamoor`
- `github.com/ethpandaops/spamoor/txbuilder`
- `github.com/ethpandaops/spamoor/utils`
- `github.com/sirupsen/logrus`
- `github.com/spf13/pflag`
- `github.com/holiman/uint256`
- `gopkg.in/yaml.v3`

See `plugin/symbols/` for the full list of extracted symbols.

## Example Plugin

See `_example-plugin/` for a complete working example that demonstrates:

- Plugin structure and descriptor
- Scenario implementation with contract deployment
- Using `BuildBoundTx` for contract interactions
- Proper `onComplete` handling
- Well-known wallets for deployment
- Configuration via flags and YAML

## Troubleshooting

### "undefined: package.Symbol"

The package needs symbols extracted. Add it to the extraction in `plugin/symbols/generate.go` and run `go generate ./plugin/symbols/...`.

### "interpreter panic"

Usually caused by unsupported Go patterns. Try simplifying:
- Use explicit function types instead of closures where possible
- Avoid complex type assertions
- Use simpler channel patterns

### "plugin.yaml not found"

For `.tar.gz` archives, `plugin.yaml` must be present. Use `make plugins` to build archives properly - it auto-generates this file. If building manually, ensure `plugin.yaml` is at the root of the archive (inside the proper path structure).

### "PluginDescriptor not found"

Ensure your `plugin.go` exports a variable named exactly `PluginDescriptor` of type `scenario.PluginDescriptor`.
