# Blob Replacements

Send blob transactions and replace them with new versions until they are included in a block.

## Usage

```bash
spamoor blob-replacements [flags]
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
- `--max-replace` - Maximum number of replacement transactions (default: 4)
- `--replace` - Seconds to wait before replacing a transaction (default: 30)
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
- `--log-txs` - Log all submitted transactions
- `--trace` - Enable tracing output

## Example

Send 100 blob transactions with replacements:
```bash
spamoor blob-replacements -p "<PRIVKEY>" -h http://rpc-host:8545 -c 100
```

Send 2 blob transactions per slot with up to 6 replacements each:
```bash
spamoor blob-replacements -p "<PRIVKEY>" -h http://rpc-host:8545 -t 2 --max-replace 6
``` 