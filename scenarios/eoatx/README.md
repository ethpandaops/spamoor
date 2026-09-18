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
- `--to` - Target address to send transactions to (overrides other target options when specified)
- `--random-target` - Use random destination addresses
- `--self-tx-only` - Use sender wallet as destination address
- `--targets-file` - Load a combined pool of explicit and CREATE2-derived receiver addresses from YAML
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
- `--log-txs` - Log all submitted transactions
- `--trace` - Enable tracing output

## Examples

Send 1000 transactions with random amounts up to 50 gwei:
```bash
spamoor eoatx -p "<PRIVKEY>" -h http://rpc-host:8545 -c 1000 --amount 50 --random-amount
```

Send 3 transactions per slot continuously:
```bash
spamoor eoatx -p "<PRIVKEY>" -h http://rpc-host:8545 -t 3
```

Send 100 transactions to a specific address:
```bash
spamoor eoatx -p "<PRIVKEY>" -h http://rpc-host:8545 -c 100 --to "0x1234567890123456789012345678901234567890"
```

## Receiver target pools

Daemon/config-file users can combine explicit addresses and one or more
deterministic CREATE2 ranges under `targets`:

```yaml
- scenario: eoatx
  name: mixed-existing-receivers
  config:
    throughput: 200
    amount: 1
    targets:
      addresses:
        - "0x1111111111111111111111111111111111111111"
        - "0x2222222222222222222222222222222222222222"
      create2_patterns:
        - name: minimal-contracts
          factory: "0x3333333333333333333333333333333333333333"
          initcode: "0x60006000f3"
          start_salt: 0
          count: 2000
        - name: max-code-contracts
          factory: "0x3333333333333333333333333333333333333333"
          initcode_hash: "0x0000000000000000000000000000000000000000000000000000000000000000"
          start_salt: 0
          count: 2000
      order: shuffle
      repeat: false
      seed: glamsterdam-p0
```

Each CREATE2 pattern must set exactly one of `initcode`, `initcode_file`, or
`initcode_hash`. Relative `initcode_file` paths are resolved relative to the
targets file when using `--targets-file`. `start_salt` uses the same uint64,
big-endian bytes32 encoding as the `factorydeploytx` scenario.

`order` can be `sequential` (the default) or `shuffle`. Shuffling is
deterministic for a given `seed`. With `repeat: false`, omitted `total_count`
is automatically limited to one transaction per target; an explicit
`total_count` cannot exceed the target pool size. With `repeat: true`, targets
wrap in the selected order.

The `targets` block is mutually exclusive with `--to`, `--random-target`, and
`--self-tx-only`. For direct CLI use, put either the `targets` block above or
its contents in a separate file:

```bash
spamoor eoatx -p "<PRIVKEY>" -h http://rpc-host:8545 \
  --throughput 200 --targets-file ./receivers.yaml
```
