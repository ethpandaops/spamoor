# Self-Destruct Deployments

Deploy contracts that immediately self-destruct after deployment.

On Amsterdam (EIP-8246), `SELFDESTRUCT` no longer burns ETH: if a contract created in the same transaction calls `SELFDESTRUCT` with itself as beneficiary (or receives ETH after the `SELFDESTRUCT`), the balance is preserved rather than destroyed. The account's code, storage, and nonce are still cleared; if the final balance is zero the account is deleted by EIP-161. This scenario exercises CREATE/CREATE2 + SELFDESTRUCT patterns extensively, including self-targeting and cross-contract destructs, and is therefore a good stress test for EIP-8246 correctness.

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
- `--basefee` - Max fee per gas in gwei (default: 20)
- `--tipfee` - Max tip per gas in gwei (default: 2)
- `--gaslimit` - Gas limit for each deployment in gwei (default: 10000000)
- `--amount` - Transfer amount per transaction in gwei (default: 20)
- `--random-amount` - Use random amounts (with --amount as limit)
- `--rebroadcast` - Seconds to wait before rebroadcasting (default: 120)

### Wallet Management
- `--max-wallets` - Maximum number of child wallets to use
- `--refill-amount` - ETH amount to fund each child wallet (default: 5)
- `--refill-balance` - Minimum ETH balance before refilling (default: 2)
- `--refill-interval` - Seconds between balance checks (default: 300)

### Client Settings
- `--client-group` - Client group to use for sending transactions
- `--deploy-client-group` - Client group to use for deployment transactions (same as --client-group if empty)

### Debug Options
- `-v, --verbose` - Enable verbose output
- `--log-txs` - Log all submitted transactions
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