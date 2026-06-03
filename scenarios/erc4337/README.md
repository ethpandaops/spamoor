# ERC-4337 Account Abstraction

Generate ERC-4337 v0.7 account abstraction load by submitting `EntryPoint.handleOps` bundles. This scenario deploys the necessary contracts and sends UserOperations that call a target contract through a smart account, with gas sponsored by a built-in paymaster. Each child wallet maintains a pool of reusable smart accounts and only deploys a new account every Nth UserOperation, so most operations reuse an existing account instead of creating a new contract on every call.

## Usage

```bash
spamoor erc4337 [flags]
```

## Configuration

### Base Settings (required)
- `--privkey` - Private key of the sending wallet
- `--rpchost` - RPC endpoint(s) to send transactions to

### Volume Control (either -c or -t required)
- `-c, --count` - Total number of handleOps transactions to send
- `-t, --throughput` - handleOps transactions to send per slot
- `--max-pending` - Maximum number of pending transactions

### Transaction Settings
- `--basefee` - Max fee per gas in gwei (default: 20)
- `--tipfee` - Max tip per gas in gwei (default: 2)
- `--basefee-wei` - Max fee per gas in wei (overrides --basefee for L2 sub-gwei fees)
- `--tipfee-wei` - Max tip per gas in wei (overrides --tipfee for L2 sub-gwei fees)
- `--rebroadcast` - Enable reliable rebroadcast system (default: 1)

### Account Abstraction Settings
- `--bundle-size` - Number of UserOperations to bundle into each handleOps transaction (default: 1)
- `--new-account-interval` - Deploy a new smart account every Nth UserOperation, reusing existing accounts in between (default: 1000; `1` = a new account on every op, `0` = only ever the first account per wallet)
- `--paymaster-deposit` - EntryPoint deposit in ETH to keep funded on the sponsoring paymaster (default: 10)
- `--timeout` - Timeout for the scenario (e.g. '1h', '30m', '5s') - empty means no timeout

### Wallet Management
- `--max-wallets` - Maximum number of child wallets to use
- If not specified, automatically calculated based on count/throughput (max 1000)

### Client Settings
- `--client-group` - Client group to use for sending transactions

### Debug Options
- `-v, --verbose` - Enable verbose output
- `--log-txs` - Log all submitted transactions
- `--trace` - Enable tracing output

## How It Works

1. **Contract Deployment**: Automatically deploys the ERC-4337 v0.7 EntryPoint, a SimpleAccountFactory, an accept-all paymaster, and a Counter target contract (reuses them on restart).
2. **Paymaster Funding**: Funds the paymaster's EntryPoint deposit from the root wallet and keeps it topped up while the scenario runs, so individual operations need no pre-funding.
3. **Bundling**: Each transaction calls `handleOps` with `--bundle-size` UserOperations. The sending child wallet acts as both the bundler (gas refund beneficiary) and the owner that signs and owns the accounts used in the bundle.
4. **Account Pooling**: Each child wallet keeps a pool of SimpleAccounts (deterministic CREATE2 salts). A new account is deployed (via `initCode`) on the first op and then every `--new-account-interval` ops; in between, existing accounts are reused with sequential EntryPoint nonces, calling `Counter.increment()` through them. Reusing accounts avoids creating a new contract on every call, dramatically reducing state growth. Because each account is only ever used by its owning wallet, its nonce stays ordered with the bundler's EOA nonce.

## Smart Contract Details

The scenario uses four contracts:

- **EntryPoint**: The canonical ERC-4337 v0.7 singleton that validates and executes UserOperations.
- **SimpleAccountFactory**: Deploys SimpleAccount proxies via CREATE2.
- **AcceptAllPaymaster**: Sponsors every UserOperation from its EntryPoint deposit (testnet-only; performs no validation).
- **Counter**: Trivial state-writing target invoked by each account.

## Example Usage

Deploy and send 1000 handleOps transactions:
```bash
spamoor erc4337 -p "<PRIVKEY>" -h http://rpc-host:8545 -c 1000
```

Send 10 handleOps per slot with 4 UserOperations each:
```bash
spamoor erc4337 -p "<PRIVKEY>" -h http://rpc-host:8545 -t 10 --bundle-size 4
```

## Important Notes

- **Account Reuse**: By default accounts are reused and a new one is deployed only every 1000 ops, so state growth per op is minimal. Set `--new-account-interval 1` to deploy a fresh account on every op (heavier state churn, useful for state-bloat stress testing), or a large value / `0` to maximize reuse.
- **Bundle Size**: Increasing `--bundle-size` stresses the EntryPoint's batched validation/execution loop harder, validating and executing more UserOperations per transaction.
- **Paymaster**: The built-in accept-all paymaster is testnet-only - it sponsors any operation and must never be deployed on a network with real value.
- **Gas Limits**: UserOperation gas limits are static. Account-deploying ops are sized to cover CREATE2 under the Amsterdam state-creation fee schedule; account-reusing ops use a smaller verification budget.
- **Rebroadcast**: Enabled by default to ensure transaction inclusion in congested networks.
