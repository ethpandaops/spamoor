# Spamoor App User Guide

Spamoor is a powerful Ethereum transaction spamming tool designed for testing and stress-testing Ethereum networks. This guide covers how to use Spamoor as an end user, both as a CLI utility and through the web-based daemon interface.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [CLI Usage](#cli-usage)
- [Daemon Mode](#daemon-mode)
- [Available Scenarios](#available-scenarios)
- [Configuration](#configuration)
- [Examples](#examples)
- [Troubleshooting](#troubleshooting)

## Prerequisites

Before using Spamoor, you need:

1. **Prefunded Root Private Key**: A private key with sufficient ETH to fund child wallets and pay for transaction fees
2. **Execution Layer RPC Endpoints**: At least one Ethereum RPC endpoint to send transactions to
3. **System Requirements**: 
   - Go 1.24+ (if building from source)
   - Docker (for containerized deployment)

## Installation

### Option 1: Using Docker (Recommended)

Pull the latest Docker image:
```bash
docker pull ethpandaops/spamoor:latest
```

### Option 2: Using Pre-built Releases

Download the latest release from the [GitHub releases page](https://github.com/ethpandaops/spamoor/releases).

### Option 3: Building from Source

```bash
git clone https://github.com/ethpandaops/spamoor.git
cd spamoor
make build
```

The built binaries will be available in the `bin/` directory.

## Quick Start

### CLI Mode - Send 100 EOA Transactions

```bash
# Using binary
./spamoor eoatx \
  --privkey "0x1234567890abcdef..." \
  --rpchost "http://localhost:8545" \
  --count 100

# Using Docker
docker run --rm ethpandaops/spamoor eoatx \
  --privkey "0x1234567890abcdef..." \
  --rpchost "http://localhost:8545" \
  --count 100
```

### Daemon Mode - Web Interface

```bash
# Using binary
./spamoor-daemon \
  --privkey "0x1234567890abcdef..." \
  --rpchost "http://localhost:8545" \
  --port 8080

# Using Docker
docker run -p 8080:8080 ethpandaops/spamoor spamoor-daemon \
  --privkey "0x1234567890abcdef..." \
  --rpchost "http://localhost:8545" \
  --port 8080
```

Then open http://localhost:8080 in your browser.

## CLI Usage

### Basic Command Structure

```bash
spamoor [scenario] [flags]
```

### Global Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--privkey`, `-p` | Private key of the root wallet | Required |
| `--rpchost`, `-h` | RPC endpoint(s) to send transactions to | Required |
| `--rpchost-file` | File containing list of RPC hosts | |
| `--seed`, `-s` | Seed for deterministic child wallet generation | |
| `--refill-amount` | ETH amount to fund each child wallet | 5 |
| `--refill-balance` | Minimum ETH balance before refilling | 2 |
| `--refill-interval` | Interval for balance checks (seconds) | 300 |
| `--verbose`, `-v` | Enable verbose logging | false |
| `--trace` | Enable trace logging | false |

### RPC Host Configuration

You can specify RPC hosts in multiple ways:

1. **Single host**: `--rpchost http://localhost:8545`
2. **Multiple hosts**: `--rpchost http://host1:8545 --rpchost http://host2:8545`
3. **Comma-separated**: `--rpchost http://host1:8545,http://host2:8545`
4. **From file**: `--rpchost-file hosts.txt`

Example `hosts.txt` file:
```
http://host1:8545
http://host2:8545
```

### Wallet Management

Spamoor uses a hierarchical wallet system:

- **Root Wallet**: Funded wallet specified by `--privkey`
- **Child Wallets**: Automatically generated and funded from the root wallet
- **Wallet Pool**: Manages multiple child wallets for concurrent transaction sending

The tool automatically:
- Generates child wallets using derived private keys
- Funds child wallets from the root wallet
- Monitors balances and refills when needed
- Distributes transactions across wallets

## Daemon Mode

The daemon mode provides a web-based interface for managing multiple concurrent spammers.

### Starting the Daemon

```bash
spamoor-daemon \
  --privkey "0x1234567890abcdef..." \
  --rpchost "http://localhost:8545" \
  --port 8080 \
  --db spamoor.db
```

### Daemon-Specific Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--port`, `-P` | Web interface port | 8080 |
| `--db`, `-d` | SQLite database file | spamoor.db |
| `--startup-spammer` | YAML file with startup spammer configurations | |
| `--without-batcher` | Disable transaction batching | false |
| `--debug` | Enable debug mode | false |

### Web Interface Features

Once the daemon is running, access the web interface at `http://localhost:8080`:

1. **Dashboard**: Overview of active spammers and system status
2. **Spammer Management**: 
   - Create new spammers with custom configurations
   - Start/stop individual spammers
   - View real-time logs and metrics
3. **Wallet Management**: Monitor wallet balances and funding status
4. **Client Pool**: View RPC endpoint status and health
5. **API Documentation**: Interactive Swagger UI at `/docs`

### Creating Spammers via Web Interface

1. Navigate to the dashboard
2. Click "Create New Spammer"
3. Select a scenario type
4. Configure parameters:
   - Transaction count or throughput
   - Gas settings
   - Scenario-specific options
5. Click "Create" to add the spammer
6. Use "Start" to begin transaction generation

### Startup Spammers

You can configure spammers to start automatically using a YAML configuration file:

```yaml
# startup-spammers.yaml
- scenario: eoatx
  config:
    throughput: 10
    amount: 100
    random_amount: true

- scenario: blobs
  config:
    throughput: 5
    sidecars: 3
```

Start daemon with startup spammers:
```bash
spamoor-daemon \
  --privkey "0x1234567890abcdef..." \
  --rpchost "http://localhost:8545" \
  --startup-spammer startup-spammers.yaml
```

## Available Scenarios

Spamoor supports multiple transaction scenarios for different testing needs. For a complete list of available scenarios with detailed descriptions and links to their documentation, see the [main README](../README.md#scenarios).

### Common Scenario Flags

Most scenarios support these common flags:

| Flag | Description |
|------|-------------|
| `--count`, `-c` | Total number of transactions to send |
| `--throughput`, `-t` | Transactions per slot (continuous mode) |
| `--max-pending` | Maximum pending transactions |
| `--basefee` | Max fee per gas (gwei) |
| `--tipfee` | Max tip per gas (gwei) |
| `--gaslimit` | Gas limit per transaction |
| `--max-wallets` | Maximum child wallets to use |

## Configuration

### Transaction Fees

Configure gas pricing:
```bash
spamoor eoatx \
  --basefee 30 \      # Max fee per gas (gwei)
  --tipfee 5 \        # Max tip per gas (gwei)
  --gaslimit 21000   # Gas limit
```

### Throughput Control

Two modes for controlling transaction volume:

1. **Fixed Count**: Send exactly N transactions
```bash
spamoor eoatx --count 1000
```

2. **Continuous Throughput**: Send N transactions per slot indefinitely
```bash
spamoor eoatx --throughput 10
```

### Concurrency Management

Control transaction concurrency:
```bash
spamoor eoatx \
  --throughput 10 \
  --max-pending 50 \    # Max unconfirmed transactions
  --max-wallets 20      # Max concurrent wallets
```

## Examples

### Example 1: Basic EOA Transactions

Send 500 EOA transactions with random amounts:
```bash
spamoor eoatx \
  --privkey "0x1234567890abcdef..." \
  --rpchost "http://localhost:8545" \
  --count 500 \
  --amount 100 \
  --random-amount \
  --verbose
```

### Example 2: Continuous Blob Transactions

Send 5 blob transactions per slot continuously:
```bash
spamoor blobs \
  --privkey "0x1234567890abcdef..." \
  --rpchost "http://localhost:8545" \
  --throughput 5 \
  --sidecars 4 \
  --blobfee 50
```

### Example 3: Contract Deployments with Multiple RPC Hosts

Deploy contracts using multiple RPC endpoints:
```bash
spamoor deploytx \
  --privkey "0x1234567890abcdef..." \
  --rpchost "http://host1:8545,http://host2:8545" \
  --count 100 \
  --basefee 50 \
  --max-wallets 10
```

### Example 4: ERC-20 Token Transfers

Send ERC-20 transfers to deployed token contracts:
```bash
spamoor erctx \
  --privkey "0x1234567890abcdef..." \
  --rpchost "http://localhost:8545" \
  --throughput 8 \
  --contract-address "0x742d35cc6639c0532fea001e77e0e4d64e7dd8a7"
```

### Example 5: Docker with Custom Configuration

Using Docker with mounted configuration:
```bash
# Create hosts file
echo -e "http://localhost:8545\nhttp://localhost:8546" > hosts.txt

# Run with mounted file
docker run --rm \
  -v $(pwd)/hosts.txt:/app/hosts.txt \
  ethpandaops/spamoor eoatx \
  --privkey "0x1234567890abcdef..." \
  --rpchost-file /app/hosts.txt \
  --count 1000
```

### Example 6: Daemon with Startup Configuration

```bash
# Create startup configuration
cat > startup.yaml << EOF
- scenario: eoatx
  config:
    throughput: 10
    amount: 50
    random_amount: true

- scenario: blobs  
  config:
    throughput: 3
    sidecars: 2
EOF

# Start daemon with startup spammers
docker run -p 8080:8080 \
  -v $(pwd)/startup.yaml:/app/startup.yaml \
  ethpandaops/spamoor spamoor-daemon \
  --privkey "0x1234567890abcdef..." \
  --rpchost "http://localhost:8545" \
  --startup-spammer /app/startup.yaml
```

## Troubleshooting

### Common Issues

1. **"No client available" Error**
   - Check RPC endpoint connectivity
   - Verify RPC hosts are reachable
   - Try with `--verbose` for detailed logs

2. **"Failed to prepare wallets" Error**
   - Ensure root wallet has sufficient ETH
   - Check private key format
   - Verify network connectivity

3. **High Gas Prices/Failed Transactions**
   - Adjust `--basefee` and `--tipfee` values
   - Monitor network congestion
   - Reduce `--throughput` if network is congested

4. **Database Errors (Daemon Mode)**
   - Check file permissions on database location
   - Ensure sufficient disk space
   - Try removing existing database file

5. **Build Errors**
   - Ensure CGO build is available (needs proper c compiler)
   - Check Go version: requires 1.24+
   - Install required build tags: `make build` uses proper tags

### Debug Options

Enable detailed logging:
```bash
# Verbose logging
spamoor eoatx --verbose ...

# Trace logging (very detailed)
spamoor eoatx --trace ...

# Log all transactions
spamoor eoatx --log-txs ...
```

### Getting Help

- View available scenarios: `spamoor` (without arguments)
- Scenario-specific help: `spamoor [scenario] --help`
- Check version: `spamoor --version`
- Web API documentation: `http://localhost:8080/docs` (daemon mode)

### Performance Tips

1. **Optimize Throughput**:
   - Use multiple RPC endpoints for better distribution
   - Increase `--max-wallets` for higher concurrency
   - Adjust `--max-pending` based on network capacity

2. **Resource Management**:
   - Monitor system memory usage during long runs
   - Use `--refill-interval` to control wallet management overhead
   - Consider using Docker for resource isolation

3. **Network Considerations**:
   - Test with lower throughput first
   - Monitor RPC endpoint response times
   - Use local or dedicated RPC endpoints for best performance

## Development Environment

### Quick Start with DevNet

Spamoor includes a complete development environment using Kurtosis that spins up a full Ethereum testnet with multiple clients:

```bash
# Start a full Ethereum devnet and launch spamoor daemon
make devnet-run
```

This command will:
1. **Start Ethereum testnet**: Uses Kurtosis to launch Geth, Reth, and Lighthouse clients
2. **Deploy additional services**: Includes Dora explorer and Blockscout block explorer
3. **Generate configuration**: Creates RPC host list and chain configuration automatically
4. **Launch spamoor daemon**: Starts spamoor-daemon connected to the testnet
5. **Pre-fund accounts**: Uses a well-known private key with pre-funded accounts

### Accessing the DevNet

Once running, you can access:
- **Spamoor Web UI**: http://localhost:8080 - Main spamoor dashboard
- **Dora Explorer**: Available on auto-assigned port (check Kurtosis output)
- **Blockscout Explorer**: Available on auto-assigned port (check Kurtosis output)
- **RPC Endpoints**: Listed in `.hack/devnet/generated-hosts.txt`

### DevNet Configuration

The devnet uses:
- **Well-known private key**: `3fd98b5187bf6526734efaa644ffbb4e3670d66f5d0268ce0323ec09124bff61`
- **Multiple EL clients**: Geth and Reth for execution
- **CL client**: Lighthouse for consensus
- **High gas limits**: 100M gas limit for testing large transactions
- **Fast finality**: Electra fork enabled for quick testing

### Customizing DevNet

Create `.hack/devnet/custom-kurtosis.devnet.config.yaml` to override the default configuration:

```yaml
participants_matrix:
  el:
    - el_type: geth
    - el_type: reth
    - el_type: nethermind  # Add more clients
  cl:
    - cl_type: lighthouse
network_params:
  preset: mainnet
  gas_limit: 200000000     # Even higher gas limit
additional_services:
  - spamoor
  - dora
  - blockscout
```

### Cleaning Up

Stop and remove the entire devnet:

```bash
make devnet-clean
```

This removes all containers, networks, and generated files.

### Benefits for Testing

The DevNet environment provides:
- **Realistic multi-client setup**: Test against multiple Ethereum implementations
- **Fast block times**: Quick feedback for transaction testing
- **Pre-configured tools**: Block explorers and monitoring included
- **Isolated environment**: No impact on mainnet or public testnets
- **Reproducible setup**: Same configuration every time
- **Complete ecosystem**: Full Ethereum stack in one command

For additional support, check the [GitHub repository](https://github.com/ethpandaops/spamoor) or create an issue with detailed logs and configuration.