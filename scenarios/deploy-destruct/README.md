# Self-Destruct Deployments

Deploy contracts that immediately self-destruct after deployment.

## Usage

```bash
spamoor deploy-destruct [flags]
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
- `--amount` - Transfer amount per transaction in gwei (default: 20)
- `--basefee` - Max fee per gas in gwei (default: 20)
- `--tipfee` - Max tip per gas in gwei (default: 2)
- `--gaslimit` - Gas limit for each deployment in gwei (default: 10000000)

### Wallet Management
- `--max-wallets` - Maximum number of child wallets to use
- `--refill-amount` - ETH amount to fund each child wallet (default: 5)
- `--refill-balance` - Minimum ETH balance before refilling (default: 2)
- `--refill-interval` - Seconds between balance checks (default: 300)

### Transaction Behavior
- `--random-amount` - Use random amounts (with --amount as limit)
- `--rebroadcast` - Seconds to wait before rebroadcasting (default: 120)

### Debug Options
- `-v, --verbose` - Enable verbose output
- `--trace` - Enable tracing output

## Example

Deploy 100 self-destructing contracts:
```bash
spamoor deploy-destruct -p "<PRIVKEY>" -h http://rpc-host:8545 -c 100
```

Deploy 3 self-destructing contracts per slot:
```bash
spamoor deploy-destruct -p "<PRIVKEY>" -h http://rpc-host:8545 -t 3
``` 