# EOA Transactions

Send standard EOA (Externally Owned Account) transactions with configurable amounts and targets.

## Usage

```bash
spamoor eoatx [flags]
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
- `--gaslimit` - Gas limit per transaction (default: 21000)
- `--amount` - Transfer amount per transaction in gwei (default: 20)
- `--random-amount` - Use random amounts (with --amount as limit)
- `--random-target` - Use random destination addresses
- `--self-tx-only` - Use sender wallet as destination address
- `--data` - Custom transaction call data to send
- `--rebroadcast` - Seconds to wait before rebroadcasting (default: 120)

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

Send 1000 transactions with random amounts up to 50 gwei:
```bash
spamoor eoatx -p "<PRIVKEY>" -h http://rpc-host:8545 -c 1000 --amount 50 --random-amount
```

Send 3 transactions per slot continuously:
```bash
spamoor eoatx -p "<PRIVKEY>" -h http://rpc-host:8545 -t 3
``` 