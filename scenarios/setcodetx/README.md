# Set Code Transactions

Send transactions that update contract code using the `SELFDESTRUCT` opcode. This scenario is useful for testing contract code updates and state changes.

## Usage

```bash
spamoor setcodetx [flags]
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
- `--amount` - Amount to send with each transaction in gwei (default: 20)
- `--data` - Transaction call data to send
- `--code-addr` - Code delegation target address to use for transactions

### Wallet Management
- `--max-wallets` - Maximum number of child wallets to use
- `--min-authorizations` - Minimum number of authorizations to send per transaction (default: 1)
- `--max-authorizations` - Maximum number of authorizations to send per transaction (default: 10)
- `--max-delegators` - Maximum number of random delegators to use (0 = no delegator gets reused)

### Transaction Behavior
- `--rebroadcast` - Seconds to wait before rebroadcasting (default: 120)
- `--random-amount` - Use random amounts for transactions (with --amount as limit)
- `--random-target` - Use random to addresses for transactions
- `--random-code-addr` - Use random delegation target for transactions

### Client Settings
- `--client-group` - Client group to use for sending transactions

## Example

Send 100 set code transactions:
```bash
spamoor setcodetx -p "<PRIVKEY>" -h http://rpc-host:8545 -c 100
```

Send 2 set code transactions per slot with random amounts:
```bash
spamoor setcodetx -p "<PRIVKEY>" -h http://rpc-host:8545 -t 2 --random-amount
``` 