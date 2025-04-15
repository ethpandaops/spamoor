# Uniswap Swaps

Execute Uniswap V2 swaps with configurable parameters. This scenario allows you to test Uniswap DEX interactions by performing buy and sell transactions.

## Usage

```bash
spamoor uniswap-swaps [flags]
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
- `--gaslimit` - Gas limit for each transaction in gwei (default: 200000)
- `--rebroadcast` - Seconds to wait before rebroadcasting (default: 120)

### Swap Settings
- `--pair-count` - Number of uniswap pairs to deploy (default: 1)
- `--min-swap` - Minimum swap amount in wei (default: 100000000000000000)
- `--max-swap` - Maximum swap amount in wei (default: 1000000000000000000000)
- `--buy-ratio` - Ratio of buy vs sell swaps (0-100, default: 50)
- `--slippage` - Slippage tolerance in basis points (default: 50)
- `--sell-threshold` - DAI balance threshold to force sell in wei (default: 100000000000000000000000)

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

Send 100 buy transactions:
```bash
spamoor uniswap-swaps -p "<PRIVKEY>" -h http://rpc-host:8545 -c 100 --buy-ratio 100
```

Send 2 sell transactions per slot:
```bash
spamoor uniswap-swaps -p "<PRIVKEY>" -h http://rpc-host:8545 -t 2 --buy-ratio 0
``` 