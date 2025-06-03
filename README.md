<img align="left" src="./.github/resources/goomy.png" width="75">
<h1>Spamoor the Transaction Spammer</h1>

Spamoor is a powerful tool for generating various types of random transactions on Ethereum testnets. Perfect for stress testing, network validation, or continuous transaction testing.

## Quick Start

```bash
# Using Docker
docker run ethpandaops/spamoor

# Building from source
git clone https://github.com/ethpandaops/spamoor.git
cd spamoor
make
./bin/spamoor
```

### Usage
```bash
spamoor <scenario> [flags]
```

All scenarios require:
- `--privkey` - Private key for the sending wallet
- `--rpchost` - RPC endpoint(s) to send transactions to

## Scenarios

Spamoor provides multiple scenarios for different transaction types:

| Scenario | Description |
|----------|-------------|
| [`eoatx`](./scenarios/eoatx/README.md) | **EOA Transactions**<br>Send standard EOA transactions with configurable amounts and targets |
| [`erctx`](./scenarios/erctx/README.md) | **ERC20 Transactions**<br>Deploy a ERC20 contract and transfer tokens |
| [`calltx`](./scenarios/calltx/README.md) | **Contract Calls**<br>Deploy a contract and repeatedly call a function on it |
| [`deploytx`](./scenarios/deploytx/README.md) | **Contract Deployments**<br>Deploy contracts with custom bytecode |
| [`deploy-destruct`](./scenarios/deploy-destruct/README.md) | **Self-Destruct Deployments**<br>Deploy contracts that self-destruct |
| [`setcodetx`](./scenarios/setcodetx/README.md) | **Set Code Transactions**<br>Send EIP-7702 setcode-transactions with various settings |
| [`uniswap-swaps`](./scenarios/uniswap-swaps/README.md) | **Uniswap Swaps**<br>Deploy and perform token swaps on Uniswap V2 pools |
| [`blobs`](./scenarios/blobs/README.md) | **Blob Transactions**<br>Send blob transactions with random data |
| [`blob-replacements`](./scenarios/blob-replacements/README.md) | **Blob Replacements**<br>Send and replace blob transactions |
| [`blob-conflicting`](./scenarios/blob-conflicting/README.md) | **Conflicting Blobs**<br>Send conflicting blob and normal transactions |
| [`blob-combined`](./scenarios/blob-combined/README.md) | **Combined Blob Testing**<br>Randomized combination of all blob scenarios |
| [`gasburnertx`](./scenarios/gasburnertx/README.md) | **Gas Burner**<br>Send transactions that burn specific amounts of gas |
| [`geastx`](./scenarios/geastx/README.md) | **Geas Transactions**<br>Send transactions that execute custom geas bytecode |
| [`storagespam`](./scenarios/storagespam/README.md) | **Storage Spam**<br>Send transactions that spam the persistent EVM storage |

## Daemon Mode

Spamoor also includes a daemon mode with web UI for managing multiple spammers. It allows you to create, monitor, and control spammers through a user interface or programmatically via HTTP endpoints.

### Usage
```bash
spamoor-daemon [flags]
```

### Flags
```
-d, --db string         The file to store the database in (default "spamoor.db")
    --debug             Run the tool in debug mode
-h, --rpchost strings   The RPC host to send transactions to
    --rpchost-file      File with a list of RPC hosts to send transactions to
-p, --privkey string    The private key of the wallet to send funds from
-P, --port int          The port to run the webui on (default 8080)
-v, --verbose           Run the tool with verbose output
    --trace             Run the tool with tracing output
```

### Web Interface
The web interface runs on `http://localhost:8080` by default and provides:
- Dashboard for managing spammers
- Real-time log streaming
- Configuration management
- Start/pause/delete functionality

### API
The daemon exposes a REST API for programmatic control.
See the API Documentation in the spamoor web interface for details.


## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
