# XEN Token Sybil Attack

Execute sybil attacks against XEN Crypto token contracts by creating multiple proxy contracts to claim ranks from different addresses. This scenario deploys the necessary contracts and performs automated sybil attacks.

## Usage

```bash
spamoor xentoken [flags]
```

## Configuration

### Base Settings (required)
- `--privkey` - Private key of the sending wallet
- `--rpchost` - RPC endpoint(s) to send transactions to

### Volume Control (either -c or -t required)
- `-c, --count` - Total number of sybil attack transactions to send
- `-t, --throughput` - Sybil attack transactions to send per slot
- `--max-pending` - Maximum number of pending transactions

### Transaction Settings
- `--basefee` - Max fee per gas in gwei (default: 20)
- `--tipfee` - Max tip per gas in gwei (default: 2)
- `--gas-limit` - Gas limit for sybil attack transactions (default: 30000000)
- `--rebroadcast` - Enable reliable rebroadcast system (default: 1)

### XEN Attack Settings
- `--xen-address` - XEN token contract address (if empty, will deploy new contract)
- `--claim-term` - XEN claim term in days (default: 1)
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

1. **Contract Deployment**: Automatically deploys XENCrypto and XENSybilAttacker contracts if not provided
2. **Sybil Attack**: Each transaction creates multiple minimal proxy contracts using CREATE2
3. **Rank Claiming**: Each proxy contract calls `claimRank()` on the XEN contract as a unique address
4. **Gas Optimization**: Uses EIP-1167 minimal proxy pattern for maximum gas efficiency

## Smart Contract Details

The scenario uses two main contracts:

- **XENCrypto**: The main XEN token contract with ranking mechanism
- **XENSybilAttacker**: Creates proxy contracts that delegate calls to claim ranks

Each sybil attack transaction will create as many proxy contracts as gas allows, with each proxy claiming a rank from the XEN contract.

## Example Usage

Deploy and attack with 100 transactions:
```bash
spamoor xentoken -p "<PRIVKEY>" -h http://rpc-host:8545 -c 100
```

Attack existing XEN contract with 5 transactions per slot:
```bash
spamoor xentoken -p "<PRIVKEY>" -h http://rpc-host:8545 -t 5 --xen-address 0x123...
```

Use custom claim term and gas limit:
```bash
spamoor xentoken -p "<PRIVKEY>" -h http://rpc-host:8545 -c 50 --claim-term 7 --gas-limit 50000000
```

## Important Notes

- **Claim Term**: Default is 1 day (minimum valid term for XEN). Higher terms may provide better rewards but have longer maturity periods.
- **Gas Limit**: Default 30M gas allows for many proxy contract creations per transaction.
- **Fresh Contracts**: On fresh XEN deployments, maximum term is 100 days.
- **Rebroadcast**: Enabled by default to ensure transaction inclusion in congested networks. 