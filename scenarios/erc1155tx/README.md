# ERC1155 Transactions

Deploy an ERC1155 multi-token contract and perform token transfers between accounts.

## Usage

```bash
spamoor erc1155tx [flags]
```

## Configuration

### Base Settings (required)
- `--privkey` - Private key of the sending wallet
- `--rpchost` - RPC endpoint(s) to send transactions to

### Volume Control (either -c or -t required)
- `-c, --count` - Total number of transactions to send (default: 0)
- `-t, --throughput` - Transactions to send per slot (default: 200)
- `--max-pending` - Maximum number of pending transactions (default: 0)

### Transaction Settings
- `--basefee` - Max fee per gas in gwei (default: 20)
- `--tipfee` - Max tip per gas in gwei (default: 2)
- `--rebroadcast` - Enable reliable rebroadcast system (default: 1)
- `--timeout` - Timeout for the scenario (e.g. '1h', '30m', '5s') - empty means no timeout

### Token Settings
- `--amount` - Transfer amount per transaction (default: 20)
- `--max-index` - Maximum token index to mint (default: 0, unlimited)
- `--batch-size` - Batch size for transactions (default: 1)
- `--random-amount` - Use random amounts for transactions (with --amount as limit)
- `--random-index` - Use random token indexes for transactions
- `--random-target` - Use random destination addresses
- `--random-batch-size` - Use random batch sizes for transactions (with --batch-size as limit)

### Wallet Management
- `--max-wallets` - Maximum number of child wallets to use (default: 0, auto-calculated)

### Client Settings
- `--client-group` - Client group to use for sending transactions
- `--deploy-client-group` - Client group to use for deployment transactions (same as --client-group if empty)

### Debug Options
- `-v, --verbose` - Enable verbose output
- `--log-txs` - Log all submitted transactions
- `--trace` - Enable tracing output

## Example

Deploy ERC1155 and send 1000 token transfers:
```bash
spamoor erc1155tx -p "<PRIVKEY>" -h http://rpc-host:8545 -c 1000
```

Send 5 token transfers per slot continuously:
```bash
spamoor erc1155tx -p "<PRIVKEY>" -h http://rpc-host:8545 -t 5
```

Send batch mints with random amounts and indexes:
```bash
spamoor erc1155tx -p "<PRIVKEY>" -h http://rpc-host:8545 -t 10 --batch-size 5 --random-amount --random-index
```

Use random batch sizes for variable transaction complexity:
```bash
spamoor erc1155tx -p "<PRIVKEY>" -h http://rpc-host:8545 -t 10 --batch-size 10 --random-batch-size
```
