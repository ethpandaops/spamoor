# Spammer Configurations

This directory contains pre-built spammer configurations for various testing scenarios. These configurations are automatically indexed and made available through the spamoor web UI's "Spammer Library" feature.

## File Structure

Each YAML file in this directory contains one or more spammer configurations that can be imported directly into spamoor. The files follow this structure:

```yaml
# Name: Human-readable name for this configuration set
# Description: Brief description of what this configuration tests
# Tags: comma, separated, tags, for, categorization
# Min_Version: v1.0.0

- scenario: scenario_name
  name: 'Individual Spammer Name'
  description: 'What this specific spammer does'
  config:
    # Spammer-specific configuration options
    throughput: 10
    max_pending: 20
    # ... other config options

- scenario: another_scenario
  name: 'Another Spammer Name'
  description: 'Another spammer description'
  config:
    # Different configuration options
```

## Required Headers

Each configuration file **must** include these header comments at the top:

### Required Headers:
- **`# Name:`** - A descriptive name for the configuration set (used in UI)
- **`# Description:`** - Brief explanation of what this configuration tests
- **`# Tags:`** - Comma-separated tags for categorization and filtering

### Optional Headers:
- **`# Min_Version:`** - Minimum spamoor version required (e.g., `v1.1.5`)

## Header Examples

```yaml
# Name: EIP-4844 Blob Transaction Testing
# Description: Comprehensive test suite for blob transactions with various sizes and patterns
# Tags: blobs, eip4844, performance, gas-optimization
# Min_Version: v1.2.0
```

```yaml
# Name: ERC20 Token Stress Test
# Description: High-throughput ERC20 token transfers and interactions
# Tags: erc20, tokens, stress-test, defi
```

## Index Generation

The `_index.yaml` file is automatically generated from these configurations using:

```bash
make generate-spammer-index
```

**Do not manually edit `_index.yaml`** - it will be overwritten during the next generation.

## Tags and Categorization

Tags are used for:
- **Filtering** in the web UI
- **Categorization** of similar configurations
- **Search functionality** by keywords

### Suggested Tag Categories:
- **Protocol/EIP**: `eip4844`, `eip7702`, `erc20`, `precompiles`
- **Purpose**: `stress-test`, `performance`, `edge-cases`, `gas-optimization`
- **Complexity**: `basic`, `intermediate`, `advanced`
- **Feature**: `blobs`, `contracts`, `defi`, `logs`, `storage`

## Version Compatibility

The `Min_Version` header ensures configurations are only shown to compatible spamoor versions:

- **Compatible configs** are shown normally with import buttons
- **Incompatible configs** are shown with warning badges and disabled import buttons
- **Development versions** (git-xxx) can see all configurations

## File Naming

Use descriptive filenames that indicate the configuration's purpose:
- `max-blob-transactions.yaml`
- `erc20-stress-test.yaml`
- `gas-limit-edge-cases.yaml`
- `storage-spam-attack.yaml`

## Configuration Best Practices

1. **Use descriptive names** for both the file header and individual spammers
2. **Add comprehensive descriptions** explaining the test purpose
3. **Tag appropriately** for easy discovery
4. **Set reasonable defaults** that work out of the box
5. **Document any special requirements** in the description
6. **Test configurations** before committing to ensure they work
7. **Use semantic versioning** for Min_Version (e.g., v1.0.0, v2.1.0)

## Adding New Configurations

1. Create a new YAML file with the required headers
2. Add one or more spammer configurations
3. Test the configuration manually
4. Run `make generate-spammer-index` to update the index
5. Commit both the new file and updated `_index.yaml`

## Integration with CI

The index generation is integrated into the build process and CI pipeline to ensure:
- The index is always up-to-date
- New configurations are automatically indexed
- YAML syntax is validated
- Header completeness is checked

## Examples

See the existing files in this directory for examples of properly formatted configurations:
- `max-eip7702-authorizations-per-block.yaml` - EIP-7702 testing
- `max-logs-per-tx-spammer.yaml` - Gas optimization testing

For more information about spammer scenarios and configuration options, see the main documentation in `/docs/scenario-developers.md`.