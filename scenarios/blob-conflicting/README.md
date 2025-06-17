# Conflicting Blobs

Send blob transactions and conflicting normal transactions simultaneously or with a small delay.

## Usage

```bash
spamoor blob-conflicting [flags]
```

## Configuration

### Base Settings (required)
- `--privkey` - Private key of the sending wallet
- `--rpchost` - RPC endpoint(s) to send transactions to

### Volume Control (either -c or -t required)
- `-c, --count` - Total number of transactions to send
- `-t, --throughput` - Transactions to send per slot
- `--max-pending` - Maximum number of pending transactions
- `--throughput-increment-interval` - Increment the throughput and pending limit every interval (in seconds). Useful for gradually increasing load over time.

### Transaction Settings
- `--basefee` - Max fee per gas in gwei (default: 20)
- `--tipfee` - Max tip per gas in gwei (default: 2)
- `--blobfee` - Max blob fee in gwei (default: 20)

### Blob Configuration
- `-b, --sidecars` - Maximum number of blob sidecars per transaction (default: 3)
- `--conflict-delay` - Milliseconds to wait before sending conflicting transaction (default: 0)

### Wallet Management
- `--max-wallets` - Maximum number of child wallets to use
- `--refill-amount` - ETH amount to fund each child wallet (default: 5)
- `--refill-balance` - Minimum ETH balance before refilling (default: 2)
- `--refill-interval` - Seconds between balance checks (default: 300)

### Client Settings
- `--client-group` - Client group to use for sending transactions

### Debug Options
- `-v, --verbose` - Enable verbose output
- `--log-txs` - Log all submitted transactions
- `--trace` - Enable tracing output

## Example

Send 100 blob transactions with immediate conflicts:
```bash
spamoor blob-conflicting -p "<PRIVKEY>" -h http://rpc-host:8545 -c 100
```

Send 2 blob transactions per slot with 100ms delay before conflicts:
```bash
spamoor blob-conflicting -p "<PRIVKEY>" -h http://rpc-host:8545 -t 2 --conflict-delay 100
``` 