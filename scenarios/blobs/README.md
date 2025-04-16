# Blob Transactions

Send blob transactions with random data. This scenario focuses on basic blob transaction functionality without replacements or conflicts.

## Usage

```bash
spamoor blobs [flags]
```

## Configuration

### Base Settings (required)
- `--privkey` - Private key of the sending wallet
- `--rpchost` - RPC endpoint(s) to send transactions to

### Volume Control (either -c or -t required)
- `-c, --count` - Total number of transactions to send
- `-t, --throughput` - Transactions to send per slot
- `--max-pending` - Maximum number of pending transactions

### Transaction Settings
- `--basefee` - Max fee per gas in gwei (default: 20)
- `--tipfee` - Max tip per gas in gwei (default: 2)
- `--blobfee` - Max blob fee in gwei (default: 20)
- `--rebroadcast` - Seconds to wait before rebroadcasting (default: 30)

### Blob Configuration
- `-b, --sidecars` - Maximum number of blob sidecars per transaction (default: 3)

### Wallet Management
- `--max-wallets` - Maximum number of child wallets to use
- `--refill-amount` - ETH amount to fund each child wallet (default: 5)
- `--refill-balance` - Minimum ETH balance before refilling (default: 2)
- `--refill-interval` - Seconds between balance checks (default: 300)

### Client Settings
- `--client-group` - Client group to use for sending transactions

### Debug Options
- `-v, --verbose` - Enable verbose output
- `--trace` - Enable tracing output

## Example

Send 1000 blob transactions:
```bash
spamoor blobs -p "<PRIVKEY>" -h http://rpc-host:8545 -c 1000
```

Send 2 blob transactions per slot with 4 sidecars each:
```bash
spamoor blobs -p "<PRIVKEY>" -h http://rpc-host:8545 -t 2 -b 4
``` 