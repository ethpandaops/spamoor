# Blob Average

Send blob transactions to maintain a network-wide average blob count per block. This scenario monitors all blob transactions on the network and dynamically sends blobs only when the average falls below a configurable target.

## Usage

```bash
spamoor blob-average [flags]
```

## Configuration

### Base Settings (required)
- `--privkey` - Private key of the sending wallet
- `--rpchost` - RPC endpoint(s) to send transactions to

### Target Control
- `-a, --target-average` - Target average blob count per block (default: 3)
- `--tracking-seconds` - Time window in seconds for tracking blob averages (default: 3600 / 1 hour)

### Volume Control
- `--max-pending` - Maximum number of pending blob transactions (default: 10)

### Transaction Settings
- `--basefee` - Max fee per gas in gwei (default: 20)
- `--tipfee` - Max tip per gas in gwei (default: 2)
- `--blobfee` - Max blob fee in gwei (default: 20)
- `--rebroadcast` - Seconds to wait before rebroadcasting (default: 1)
- `--submit-count` - Number of clients to submit transaction to (default: 3)

### Blob Configuration
- `-b, --sidecars` - Number of blob sidecars per transaction (default: 1)
- `--blob-v1-percent` - Percentage of transactions using v1 wrapper format (default: 100)
- `--fulu-activation` - Unix timestamp of Fulu activation

### Wallet Management
- `--max-wallets` - Maximum number of child wallets to use (default: 20)
- `--refill-amount` - ETH amount to fund each child wallet (default: 5)
- `--refill-balance` - Minimum ETH balance before refilling (default: 2)
- `--refill-interval` - Seconds between balance checks (default: 300)

### Client Settings
- `--client-group` - Client group to use for sending transactions

### Debug Options
- `-v, --verbose` - Enable verbose output
- `--log-txs` - Log all submitted transactions
- `--trace` - Enable tracing output

## How It Works

1. **Block Monitoring**: Subscribes to block updates and counts all blobs from type 3 (blob) transactions across the network
2. **Average Tracking**: Maintains a sliding window of blob counts over the configured time period (default 1 hour)
3. **Deficit Calculation**: Calculates the difference between target average and current average
4. **Dynamic Sending**: When below target, sends enough blob transactions to compensate for the deficit
5. **Automatic Throttling**: When at or above target, pauses blob sending to avoid oversaturation

## Example

Maintain an average of 5 blobs per block:
```bash
spamoor blob-average -p "<PRIVKEY>" -h http://rpc-host:8545 -a 5
```

Maintain 3 blobs average with 2 hour tracking window and 3 sidecars per transaction:
```bash
spamoor blob-average -p "<PRIVKEY>" -h http://rpc-host:8545 \
  -a 3 --tracking-seconds 7200 -b 3
```

High-throughput configuration with more pending transactions:
```bash
spamoor blob-average -p "<PRIVKEY>" -h http://rpc-host:8545 \
  -a 6 --max-pending 20 --max-wallets 50
```
