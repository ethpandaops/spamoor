# Contract Call Transactions

Deploy a contract and repeatedly call a function on it.

## Usage

```bash
spamoor calltx [flags]
```

## Configuration

### Base Settings (required)
- `--privkey` - Private key of the sending wallet
- `--rpchost` - RPC endpoint(s) to send transactions to

### Volume Control (either -c or -t required)
- `-c, --count` - Total number of call transactions to send
- `-t, --throughput` - Call transactions to send per slot
- `--max-pending` - Maximum number of pending transactions

### Contract Settings (required)
- `--contract-code` - Contract bytecode to deploy (hex string)
- `--contract-file` - Contract file to deploy (local file or HTTP URL)
- `--contract-address` - Address of already deployed contract (skips deployment)
- `--contract-args` - Constructor arguments for the contract (hex string)
- `--contract-addr-path` - Path to child contract created during deployment (e.g. '.0.1' for nonce 1 of nonce 0)
- `--call-data` - Data to pass to the function calls (hex string)

### ABI-Based Call Data (alternative to --call-data)
- `--call-abi` - JSON ABI of the contract for function calls
- `--call-abi-file` - JSON ABI file of the contract for function calls (local file or HTTP URL)
- `--call-fn-name` - Function name to call (requires --call-abi or --call-abi-file)
- `--call-fn-sig` - Function signature to call (alternative to --call-abi/--call-abi-file)
- `--call-args` - JSON array of arguments to pass to the function

#### ABI Call Arguments Placeholders
The `--call-args` parameter supports the following placeholders:
- `{txid}` - Transaction index/ID
- `{random}` - Random uint256 value
- `{random:N}` - Random number between 0 and N
- `{randomaddr}` - Random Ethereum address

### CREATE2 Receiver Address List (alternative to --call-data / --call-abi)
- `--targets-file` - YAML file with CREATE2 receiver patterns; calltx encodes them as `callAttack` call data
- `--call-value` - Wei sent by each CALL the target contract makes (`0` or `1`)
- `--gas-buffer` - Gas the target contract leaves unspent before it stops looping (default: 50000)
- `--salt-stride` - Salts to advance between consecutive txs of one pattern (default: 0 = every tx restarts at `start_salt`)

These are mutually exclusive with `--call-data`, `--call-abi`, `--call-abi-file`,
`--call-fn-name`, `--call-fn-sig` and `--call-args`: when a target list is set,
calltx encodes the call data itself. See
[Target list format](#target-list-format).

### Transaction Settings
- `--basefee` - Max fee per gas in gwei (default: 20)
- `--tipfee` - Max tip per gas in gwei (default: 2)
- `--deploy-gas-limit` - Gas limit for deployment transaction (default: 4000000)
- `--gas-limit` - Gas limit for call transactions. Set to 0 to burn all available block gas. (default: 1000000)
- `--amount` - ETH amount to send with each call in gwei (default: 20)
- `--random-amount` - Use random amounts (with --amount as limit)
- `--random-target` - Use random destination addresses
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

Deploy contract from bytecode and send 1000 function calls:
```bash
spamoor calltx -p "<PRIVKEY>" -h http://rpc-host:8545 -c 1000 \
  --contract-code "608060405234801561001057600080fd5b50..." \
  --call-data "0xa9059cbb000000000000000000000000..."
```

Deploy contract from file and send 5 calls per slot continuously:
```bash
spamoor calltx -p "<PRIVKEY>" -h http://rpc-host:8545 -t 5 \
  --contract-file "./contract.bin" \
  --contract-args "0x000000000000000000000000..." \
  --call-data "0x06fdde03"
```

Deploy contract from remote URL:
```bash
spamoor calltx -p "<PRIVKEY>" -h http://rpc-host:8545 -c 100 \
  --contract-file "https://example.com/contract.bin" \
  --call-data "0x70a08231000000000000000000000000..."
```

## ABI-Based Examples

Call transfer function using ABI:
```bash
spamoor calltx -p "<PRIVKEY>" -h http://rpc-host:8545 -c 1000 \
  --contract-file "./erc20.bin" \
  --call-abi '[{"type":"function","name":"transfer","inputs":[{"name":"to","type":"address"},{"name":"amount","type":"uint256"}]}]' \
  --call-fn-name "transfer" \
  --call-args '["{randomaddr}", "{random:1000000}"]'
```

Call function using signature only:
```bash
spamoor calltx -p "<PRIVKEY>" -h http://rpc-host:8545 -t 5 \
  --contract-code "608060405234801561001057600080fd5b50..." \
  --call-fn-sig "setValue(uint256)" \
  --call-args '["{txid}"]'
```

Complex function call with multiple argument types:
```bash
spamoor calltx -p "<PRIVKEY>" -h http://rpc-host:8545 -c 500 \
  --contract-file "./contract.bin" \
  --call-fn-sig "complexFunction(address,uint256,bool,string)" \
  --call-args '["{randomaddr}", "{random}", true, "test-{txid}"]'
```

Call function using ABI from local file:
```bash
spamoor calltx -p "<PRIVKEY>" -h http://rpc-host:8545 -c 1000 \
  --contract-file "./erc20.bin" \
  --call-abi-file "./erc20-abi.json" \
  --call-fn-name "transfer" \
  --call-args '["{randomaddr}", "{random:1000000}"]'
```

Call function using ABI from remote URL:
```bash
spamoor calltx -p "<PRIVKEY>" -h http://rpc-host:8545 -t 5 \
  --contract-code "608060405234801561001057600080fd5b50..." \
  --call-abi-file "https://example.com/contract-abi.json" \
  --call-fn-name "setValue" \
  --call-args '["{txid}"]'
```

Call existing contract using address:
```bash
spamoor calltx -p "<PRIVKEY>" -h http://rpc-host:8545 -c 1000 \
  --contract-address "0x1234567890123456789012345678901234567890" \
  --call-fn-sig "transfer(address,uint256)" \
  --call-args '["{randomaddr}", "{random:1000000}"]'
```

Call existing contract with ABI file:
```bash
spamoor calltx -p "<PRIVKEY>" -h http://rpc-host:8545 -t 10 \
  --contract-address "0xA0b86a33E6441e047c84E7c19Ff8e4Ca6c2B5B2F" \
  --call-abi-file "./usdc-abi.json" \
  --call-fn-name "balanceOf" \
  --call-args '["{randomaddr}"]'
```

## Child Contract Examples

Call a child contract created during deployment:
```bash
spamoor calltx -p "<PRIVKEY>" -h http://rpc-host:8545 -c 1000 \
  --contract-file "./factory.bin" \
  --contract-addr-path ".0" \
  --call-fn-sig "setValue(uint256)" \
  --call-args '["{txid}"]'
```

Call a nested child contract (second level):
```bash
spamoor calltx -p "<PRIVKEY>" -h http://rpc-host:8545 -t 5 \
  --contract-file "./factory.bin" \
  --contract-addr-path ".0.1" \
  --call-data "0x06fdde03"
```

## Target list format

`--targets-file` points at a list of deterministic CREATE2 receiver patterns.
calltx never enumerates the addresses itself: it passes the three CREATE2 inputs
- factory, initcode hash and salt range - to the contract being called, which
derives the addresses on-chain. The call data it builds is

```
callAttack(address factory, bytes32 initCodeHash, uint256 startSalt, uint256 callValue, uint256 gasBuffer)
```

so the contract at `--contract-code` / `--contract-address` has to implement that
signature. `startSalt` is where that transaction begins its walk; the contract is
expected to increment the salt from there.

The file uses the same YAML shape as the `targets` block of the `eoatx`
scenario, so one file can be pointed at either scenario:

```yaml
create2_patterns:
  - name: minimal-contracts
    factory: "{factory_address}"
    initcode: "0x60006001f3"
    start_salt: 0
    count: 20000
  - name: max-code-contracts
    factory: "{factory_address}"
    initcode_file: "./bloat-initcode.hex"
    start_salt: 0
    count: 20000
```

It can also be wrapped in a top-level `targets` key, so the same block works as
a scenario config and as a `--targets-file`. In a config file or in daemon mode
the block goes inline under `targets:` instead, which is mutually exclusive with
`targets_file`.

### Fields

- `factory` accepts the `{factory_address}` placeholder, which expands to the
  well-known CREATE2 factory of the `factorydeploytx` scenario. Prefer it over a
  literal address: that factory is derived from the root wallet rather than being
  a global constant, so a hardcoded address only works for one deployment. A
  literal address still works when the factory was deployed elsewhere -
  `factorydeploytx` logs it as `using CREATE2 factory at: 0x...`.
- Exactly one of `initcode`, `initcode_file` or `initcode_hash` must be set.
  Relative `initcode_file` paths resolve relative to the targets file. Only the
  hash is needed to derive addresses, so `initcode_hash` alone is enough.
- `start_salt` and `count` describe the deployed salt range, and must match what
  `factorydeploytx` was run with. `start_salt` uses the same uint64 big-endian
  bytes32 encoding as that scenario.
- `addresses`, `order`, `repeat` and `seed` exist in the eoatx block but have
  nothing to map onto here, because the addresses are derived on-chain from a
  salt range rather than listed. They are ignored with a warning, which keeps a
  file shareable between the two scenarios.

### Multiple patterns

Consecutive transactions round-robin over the patterns, so one run can mix
target code sizes: tx 0 uses the first pattern, tx 1 the second, and so on.

### Choosing a salt stride

`--salt-stride` decides how far `start_salt` moves between consecutive txs of
one pattern:

- `0` (the default) restarts every tx at `start_salt`, so all txs re-walk the
  same prefix of the range. EIP-2929 warmth is per-tx, so each tx still pays the
  full cold-access cost, while the client's own state cache stays warm across
  them.
- A non-zero stride gives each tx a disjoint slice and wraps back to
  `start_salt` once the range is exhausted, which spreads the walk over the
  whole range.

Set it to the number of iterations a single transaction can pay for at
`--gas-limit`. `count` must be at least the stride, and should be well above it
so the walk covers a useful span before wrapping. When `count` is not a multiple
of the stride, the trailing partial slice is skipped rather than truncated, so a
transaction never walks past the deployed range - which matters with
`--call-value 1`, where a salt that was never deployed is charged the 25000 gas
account-creation cost instead of a plain cold access.

`--call-value 1` also needs the called contract to hold balance, since it sends
wei per iteration. `--amount` tops it up with every call tx (in gwei); calltx
warns when `--call-value` is set and `--amount` is 0.
