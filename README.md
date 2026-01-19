<img align="left" src="./.github/resources/goomy.png" width="75">
<h1>Spamoor the Transaction Spammer</h1>

**A powerful Ethereum transaction generator for testnets** 🚀

[![Go Report Card](https://goreportcard.com/badge/github.com/ethpandaops/spamoor)](https://goreportcard.com/report/github.com/ethpandaops/spamoor)
[![License](https://img.shields.io/github/license/ethpandaops/spamoor)](LICENSE)
[![Docker](https://img.shields.io/docker/pulls/ethpandaops/spamoor)](https://hub.docker.com/r/ethpandaops/spamoor)

Spamoor is a robust transaction spamming tool designed for stress testing, network validation, and continuous transaction testing on Ethereum testnets. With 14+ different transaction scenarios and a powerful web-based daemon mode, it's the ultimate tool for Ethereum network testing.

## 🚀 Quick Start

### CLI Output Example
![CLI Running](.github/resources/cli-running-example.png)

```bash
# 🐳 Using Docker
docker run ethpandaops/spamoor

# 🔧 Building from source
git clone https://github.com/ethpandaops/spamoor.git
cd spamoor
make
./bin/spamoor
```

## 📘 Usage

### Basic Command Structure
```bash
# Run a single scenario
spamoor <scenario> [flags]

# Run multiple scenarios from YAML configuration
spamoor run <yaml-file> [flags]
```

### 🔑 Required Parameters
| Parameter | Description |
|-----------|-------------|
| `--privkey` | Private key for the root wallet (funds child wallets) |
| `--rpchost` | RPC endpoint(s) to send transactions to |

### 🔧 Advanced Configuration

```bash
# Basic usage
spamoor eoatx --privkey="0x..." --rpchost="http://localhost:8545"

# Multiple RPC endpoints
spamoor eoatx --privkey="0x..." \
  --rpchost="http://node1:8545" \
  --rpchost="http://node2:8545"

# With authentication
spamoor eoatx --privkey="0x..." \
  --rpchost="headers(Authorization:Bearer token)http://node:8545"
```

💡 **See the [App User Guide](./docs/app-users.md) for advanced RPC configuration options**

## 🎯 Transaction Scenarios

Spamoor provides a comprehensive suite of transaction scenarios for different testing needs:

| Scenario | Description |
|----------|-------------|
| [`eoatx`](./scenarios/eoatx/README.md) | **EOA Transactions** - Send standard ETH transfers with configurable amounts |
| [`erc20tx`](./scenarios/erc20tx/README.md) | **ERC20 Transactions** - Deploy ERC20 tokens and perform transfers |
| [`erc721tx`](./scenarios/erc721tx/README.md) | **ERC721 Transactions** - Deploy ERC721 NFTs and perform transfers |
| [`erc1155tx`](./scenarios/erc1155tx/README.md) | **ERC1155 Transactions** - Deploy ERC1155 NFTs and perform transfers |
| [`calltx`](./scenarios/calltx/README.md) | **Contract Calls** - Deploy contracts and repeatedly call functions |
| [`deploytx`](./scenarios/deploytx/README.md) | **Contract Deployments** - Deploy contracts with custom bytecode |
| [`deploy-destruct`](./scenarios/deploy-destruct/README.md) | **Self-Destruct Deployments** - Deploy self-destructing contracts |
| [`setcodetx`](./scenarios/setcodetx/README.md) | **Set Code Transactions** - EIP-7702 setcode transactions |
| [`uniswap-swaps`](./scenarios/uniswap-swaps/README.md) | **Uniswap Swaps** - Deploy and test Uniswap V2 token swaps |
| [`blobs`](./scenarios/blobs/README.md) | **Blob Transactions** - Send blob transactions with random data |
| [`blob-average`](./scenarios/blob-average/README.md) | **Blob Average** - Maintain network-wide average blob count per block |
| [`blob-replacements`](./scenarios/blob-replacements/README.md) | **Blob Replacements** - Test blob transaction replacement |
| [`blob-conflicting`](./scenarios/blob-conflicting/README.md) | **Conflicting Blobs** - Test conflicting blob/normal transactions |
| [`blob-combined`](./scenarios/blob-combined/README.md) | **Combined Blob Testing** - Randomized blob scenario combinations |
| [`gasburnertx`](./scenarios/gasburnertx/README.md) | **Gas Burner** - Burn specific amounts of gas |
| [`storagespam`](./scenarios/storagespam/README.md) | **Storage Spam** - Stress test EVM storage |
| [`geastx`](./scenarios/geastx/README.md) | **Geas Transactions** - Execute custom geas bytecode |
| [`xentoken`](./scenarios/xentoken/README.md) | **XEN Sybil Attack** - Simulate XEN token sybil attacks |
| [`taskrunner`](./scenarios/taskrunner/README.md) | **Task Runner** - Execute configurable task sequences with init and execution phases |

## 🚄 Run Command - Execute Multiple Scenarios

The `spamoor run` command allows you to execute multiple scenarios concurrently from a YAML configuration file, without needing the full daemon infrastructure.

### Usage
```bash
spamoor run <yaml-file> [flags]
```

### Flags
| Flag | Short | Description |
|------|-------|-------------|
| `--spammers` | `-s` | Indexes of spammers to run (e.g., `-s 0,2`). If not specified, runs all |
| `--privkey` | `-p` | Private key for the root wallet |
| `--rpchost` | `-h` | RPC endpoints (multiple allowed) |
| `--rpchost-file` | - | File containing RPC endpoints |
| `--verbose` | `-v` | Enable verbose logging |
| `--trace` | - | Enable trace logging |

### Example Configuration
```yaml
# example-run-config.yaml
- scenario: eoatx
  name: "High Volume ETH Transfers"
  description: "Basic ETH transfers with 400 tx/slot throughput"
  config:
    seed: eoatx-demo
    refill_amount: 1000000000000000000  # 1 ETH in wei
    refill_balance: 500000000000000000  # 0.5 ETH in wei
    refill_interval: 600  # 10 minutes
    throughput: 400
    max_pending: 800
    max_wallets: 200
    base_fee: 20  # gwei
    tip_fee: 2    # gwei
    amount: 100   # 100 wei per transfer
    random_amount: true
    random_target: true

- scenario: erctx
  name: "ERC20 Token Spammer"
  description: "Deploy and transfer ERC20 tokens"
  config:
    seed: erctx-demo
    refill_amount: 2000000000000000000  # 2 ETH
    refill_balance: 1000000000000000000 # 1 ETH
    throughput: 100
    max_pending: 200
    max_wallets: 100
    amount: 50
    random_amount: true
    
# Include other configuration files
- include: ./high-load-configs.yaml
```

### Examples
```bash
# Run all spammers from configuration
spamoor run config.yaml -h http://localhost:8545 -p 0x...

# Run only specific spammers (0-based indexes)
spamoor run config.yaml -h http://localhost:8545 -p 0x... -s 0,2

# With multiple RPC endpoints
spamoor run config.yaml -p 0x... \
  -h http://node1:8545 \
  -h http://node2:8545
```

## 🖥️ Daemon Mode

Run Spamoor as a daemon with a powerful web interface for managing multiple concurrent spammers.

![Spamoor Dashboard](.github/resources/dashboard-overview.png)

### 🚀 Starting the Daemon

```bash
spamoor-daemon [flags]
```

### ⚙️ Configuration Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--db` | `-d` | Database file location | `spamoor.db` |
| `--rpchost` | `-h` | RPC endpoints (multiple allowed) | - |
| `--rpchost-file` | - | File containing RPC endpoints | - |
| `--privkey` | `-p` | Root wallet private key | - |
| `--port` | `-P` | Web UI port | `8080` |
| `--startup-spammer` | - | YAML file or URL with startup spammers | - |
| `--startup-delay` | - | Delay in seconds before starting spammers on startup | `30` |
| `--fulu-activation` | - | Unix timestamp of Fulu activation | `0` |
| `--without-batcher` | - | Disable transaction batching | `false` |
| `--disable-tx-metrics` | - | Disable transaction metrics collection | `false` |
| `--verbose` | `-v` | Enable verbose logging | `false` |
| `--debug` | - | Enable debug mode | `false` |
| `--trace` | - | Enable trace logging | `false` |

### 🌐 Web Interface Features

Access the web UI at `http://localhost:8080` for:

#### 📊 Dashboard - Real-time Spammer Management
![Dashboard with Spammers](.github/resources/dashboard-interface.png)
- View all spammers at a glance with their status, scenario type, and actions
- Mass actions for bulk operations on multiple spammers
- Quick access to start, stop, edit, and delete spammers
- Expandable live logs for each spammer

#### 💰 Wallets Overview
![Wallets Page](.github/resources/wallets-overview.png)
- Monitor wallet balances and transaction counts
- Track child wallets spawned from the root wallet
- Real-time balance updates

#### 🔧 Client Management
![Clients Configuration](.github/resources/clients-overview.png)
- View and manage RPC client connections
- Group clients for load balancing
- Enable/disable clients dynamically
- Monitor client status and block height

#### 📈 Performance Graphs
![Performance Metrics](.github/resources/graphs-overview.png)
- Visual monitoring of transaction throughput
- Success/failure rates per spammer
- Network performance metrics

### 🔌 REST API

The daemon exposes a comprehensive REST API for automation:

#### Key Endpoints
- **📡 Client Management**
  - `GET /api/clients` - List all RPC clients
  - `PUT /api/client/{index}/groups` - Update client groups
  - `PUT /api/client/{index}/name` - Set custom display name
  
- **🤖 Spammer Control**
  - `GET /api/spammers` - List all spammers
  - `POST /api/spammer` - Create new spammer
  - `PUT /api/spammer/{id}/start` - Start spammer
  - `PUT /api/spammer/{id}/pause` - Pause spammer
  
- **💾 Import/Export**
  - `GET /api/export` - Export configurations
  - `POST /api/import` - Import configurations

📚 **Full API documentation available at** `/docs` **when daemon is running**

### 📦 Import/Export System

#### 🎯 Create New Spammer Dialog
![Create Spammer Dialog](.github/resources/create-spammer.png)
- Select from 13+ pre-configured scenarios
- Customize configuration with YAML editor
- Real-time YAML validation
- Access to pre-built configuration library

#### Export Configurations
```bash
# Export all spammers to YAML
curl http://localhost:8080/api/export > my-spammers.yaml
```

#### Import Configurations

![Import Spammer Dialog](.github/resources/import-spammer.png)
*Import dialog showing YAML validation and configuration options*

![Spammer Library](.github/resources/spammer-import-library.png)
*Pre-built configuration library for common test scenarios*

```bash
# Import on startup
spamoor-daemon --startup-spammer="spammer-configs.yaml"

# Import via API
curl -X POST http://localhost:8080/api/import \
  -H "Content-Type: application/yaml" \
  --data-binary @spammer-configs.yaml
```

#### YAML Configuration Format
```yaml
# Direct spammer definition
- scenario: "eoatx"
  name: "Basic ETH Transfers"
  config:
    wallet_count: 10
    transaction_count: 1000

# Include other files
- include: "common-spammers.yaml"
- include: "https://example.com/stress-tests.yaml"
```

### 📊 Prometheus Metrics

Access metrics at `/metrics` for monitoring with Prometheus/Grafana:

- **🔢 Transaction Metrics**: Success rates, gas usage, confirmation times
- **📈 Performance Metrics**: TPS, latency, queue sizes
- **🎯 Scenario Metrics**: Custom metrics per scenario type

## 📚 Documentation

Comprehensive documentation is available for different user types:

### 📖 [App User Guide](./docs/app-users.md)
Complete guide for CLI usage, daemon mode, configuration, and troubleshooting. Covers:
- Installation and setup options
- Detailed CLI usage with examples
- Web interface walkthrough
- Configuration management
- Troubleshooting common issues

### 🔌 [API Consumer Guide](./docs/api-consumers.md)
REST API documentation with bash/curl examples for all endpoints. Includes:
- Complete API reference with examples
- Spammer management endpoints
- Client configuration and monitoring
- Real-time data streaming
- Import/export functionality

### 🛠️ [Scenario Developer Guide](./docs/scenario-developers.md)
Comprehensive guide for implementing custom transaction scenarios. Covers:
- Scenario architecture and lifecycle
- Critical development rules and best practices
- Wallet management and transaction building
- Contract interaction patterns
- Testing and debugging scenarios

## ✨ Key Features

- **🎯 14+ Transaction Scenarios**: From basic EOA transfers to complex DeFi interactions
- **🖥️ Web-Based Daemon Mode**: Manage multiple spammers through an intuitive UI
- **🔌 REST API**: Full programmatic control for automation
- **📊 Prometheus Metrics**: Built-in monitoring and observability
- **🐳 Docker Support**: Easy deployment with official images
- **🔧 Highly Configurable**: YAML configs, CLI flags, and runtime adjustments
- **🚀 High Performance**: Optimized for maximum transaction throughput
- **🛡️ Production Ready**: Battle-tested on Ethereum testnets

## 🏗️ Development

### Prerequisites
- Go 1.24+
- CGO enabled (required for SQLite and cryptographic operations)
- Build tags: `with_blob_v1,ckzg` for blob transaction support

### Quick Development Setup
```bash
# Clone the repository
git clone https://github.com/ethpandaops/spamoor.git
cd spamoor

# Run tests
make test

# Build binaries
make build

# Run development environment with local testnet
make devnet-run
```

### Project Structure
```
spamoor/
├── .cursor/rules/          # Development standards and guidelines
├── cmd/                    # CLI applications (spamoor & spamoor-daemon)
├── daemon/                 # Daemon implementation
├── docs/                   # User documentation
├── scenario/               # Core scenario framework
├── scenarios/              # Transaction scenario implementations
├── scripts/                # Build and utility scripts
├── spammer-configs/        # Pre-built spammer configuration library
├── spamoor/                # Core spamoor logic
├── txbuilder/              # Transaction building utilities
├── utils/                  # Shared utilities
└── webui/                  # Web interface frontend
```

## 🤝 Contributing

We welcome contributions! Here's how you can help:

1. **🐛 Report Issues**: Found a bug? [Open an issue](https://github.com/ethpandaops/spamoor/issues)
2. **💡 Suggest Features**: Have an idea? Share it in the discussions
3. **📝 Improve Documentation**: Help us make the docs better
4. **🧑‍💻 Submit Code**: Check our [development guidelines](.cursor/rules/development_workflow.mdc)

### Development Process
1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Standards
- Follow the [code standards](.cursor/rules/code_standards.mdc)
- Add tests for new features
- Update documentation as needed
- Ensure all tests pass

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

Built with ❤️ by the [Ethpandaops](https://github.com/ethpandaops) team.

Special thanks to all [contributors](https://github.com/ethpandaops/spamoor/graphs/contributors) who have helped make Spamoor better!
