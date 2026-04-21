# Storage Refund

Send transactions that write new storage slots and clear old ones to trigger gas refunds. Each transaction writes N new storage slots (zero to non-zero) and clears N previously set slots (non-zero to zero), generating SSTORE refunds. Only slots written in a previous block are cleared, so multiple calls within the same block target different slots.

This scenario is designed for testing EIP-7778 (Block Gas Accounting without Refunds) edge cases, where the gap between gas-before-refund and gas-after-refund can cause payload builder failures on full blocks.

## Usage

```bash
spamoor storagerefundtx [flags]
```

## Configuration

### Base Settings (required)
- `--privkey` - Private key of the sending wallet
- `--rpchost` - RPC endpoint(s) to send transactions to

### Volume Control (either -c or -t required)
- `-c, --count` - Total number of transactions to send
- `-t, --throughput` - Transactions to send per slot (default: 10)
- `--max-pending` - Maximum number of pending transactions

### Transaction Settings
- `--basefee` - Max fee per gas in gwei (default: 20)
- `--tipfee` - Max tip per gas in gwei (default: 2)
- `--slots-per-call` - Number of storage slots to write and clear per transaction (default: 500)
- `--gaslimit` - Gas limit per transaction (default: 0 = auto-estimate based on slots-per-call)
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

## How It Works

1. A `StorageRefund` contract is deployed that tracks write/clear pointers and clearable ranges
2. Each transaction calls `execute(slotsPerCall)` which:
   - Updates the clearable boundary on new blocks (slots written before this block become clearable)
   - Clears up to `slotsPerCall` old storage slots (non-zero to zero, generating gas refunds)
   - Writes `slotsPerCall` new storage slots (zero to non-zero)
3. The first block's transactions only write new slots (no old slots to clear yet)
4. Subsequent blocks both write and clear, generating significant gas refunds

### Gas Refund Impact

With 500 slots per call:
- Writing 500 new slots: ~10,000,000 gas (500 * 20,000)
- Clearing 500 old slots: ~2,500,000 gas before refund (500 * 5,000)
- Gas refund: ~2,400,000 (500 * 4,800)
- Gap between gas-before-refund and gas-after-refund: ~2,400,000 per transaction

This gap accumulates across transactions in a block and can trigger the EIP-7778 accounting mismatch in payload builders.

## Example

Send 10 transactions per slot with 500 storage slots each:
```bash
spamoor storagerefundtx -p "<PRIVKEY>" -h http://rpc-host:8545 -t 10 --slots-per-call 500
```

Send 100 transactions with 200 storage slots each (more txs per block, smaller refund per tx):
```bash
spamoor storagerefundtx -p "<PRIVKEY>" -h http://rpc-host:8545 -c 100 --slots-per-call 200
```

Run continuously with a 1-hour timeout:
```bash
spamoor storagerefundtx -p "<PRIVKEY>" -h http://rpc-host:8545 -t 5 --timeout 1h
```
