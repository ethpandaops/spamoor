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

## Account-access attack

`contracts/ContractCallAttack.geas` is a port of the
`opcode = CALL, overhead_baseline = False` arm of `test_account_access` from
[execution-specs](https://github.com/ethereum/execution-specs/blob/master/tests/benchmark/stateful/bloatnet/test_account_query.py).
It derives the CREATE2 addresses of contracts deployed by the `factorydeploytx`
scenario and hits each one with a cold `CALL`, so a run only has to carry the
three CREATE2 inputs - factory, initcode hash and salt range - rather than a
list of addresses.

The memory layout is byte-for-byte the `Create2PreimageLayout(offset=0)` of
execution-specs, and the salt walk is monotonic - `startSalt`, `startSalt + 1`,
... with no wrap-around - matching its `increment_salt_op`. Keeping the walk
inside the salt range that was actually deployed is therefore the caller's job.

The call data is ABI compatible with the `callAttack` function of
`ContractReadAttackController.sol`, so the same arguments drive either
implementation:

```
callAttack(address factory, bytes32 initCodeHash, uint256 startSalt, uint256 callValue, uint256 gasBuffer)
```

`callValue` selects the two arms the benchmark parametrises: `0` measures the
cold account access alone, `1` adds the value transfer. `gasBuffer` is the gas
the loop leaves unspent before it returns - it has to exceed one iteration, or
the loop runs out of gas instead of stopping.

Build the deploy bytecode with [geas](https://github.com/fjl/geas):

```bash
geas scenarios/calltx/contracts/ContractCallAttackCtor.geas
```

### Step 1: deploy the targets

A cold access is only interesting against an account that exists, so deploy the
targets first. Keep the `--init-code` and `--start-salt` values - the attack has
to be pointed at the same salt range:

```bash
spamoor factorydeploytx -p "<PRIVKEY>" -h http://rpc-host:8545 \
  -c 20000 -t 50 \
  --init-code "0x60006001f3" \
  --start-salt 0
```

`factorydeploytx` logs the factory it used as `using CREATE2 factory at: 0x...`.

### Step 2: run the attack

`initCodeHash` is `keccak256` of the init code from step 1 - for the
`0x60006001f3` above that is
`0x69bbd4b361c1909d4a86161461f44312fb9ec4112b37041b9411ad86310b6b48`. The
`{factory_address}` placeholder fills in the well-known factory of
`factorydeploytx`, which is derived from the root wallet rather than being a
global constant:

```bash
spamoor calltx -p "<PRIVKEY>" -h http://rpc-host:8545 -t 10 \
  --contract-code "$(geas scenarios/calltx/contracts/ContractCallAttackCtor.geas)" \
  --gas-limit 30000000 \
  --call-fn-sig "callAttack(address,bytes32,uint256,uint256,uint256)" \
  --call-args '["{factory_address}","0x69bbd4b361c1909d4a86161461f44312fb9ec4112b37041b9411ad86310b6b48","0","0","50000"]'
```

The `value_sent = 1` arm sets `callValue` to `1`. Each iteration then sends
1 wei, so the attack contract needs a balance: `--amount` tops it up with every
call tx (in gwei), and one gwei covers a million iterations.

```bash
spamoor calltx -p "<PRIVKEY>" -h http://rpc-host:8545 -t 10 \
  --contract-code "$(geas scenarios/calltx/contracts/ContractCallAttackCtor.geas)" \
  --gas-limit 30000000 --amount 1 \
  --call-fn-sig "callAttack(address,bytes32,uint256,uint256,uint256)" \
  --call-args '["{factory_address}","0x69bbd4b361c1909d4a86161461f44312fb9ec4112b37041b9411ad86310b6b48","0","1","50000"]'
```

Because `--call-args` is fixed for the whole run, every transaction restarts its
walk at the same `startSalt`. EIP-2929 warmth is per-tx, so each transaction
still pays the full cold-access cost, while the client's own state cache stays
warm across them - that trades some disk-level pressure for a fully predictable
walk. Use `{txid}` in `--call-args` to give each transaction its own slice.

### Sizing the salt range

Measured per-iteration cost (geth `evm run`, shanghai, targets deployed):

| `callValue` | gas/iteration | attack share | iterations per 30M gas |
| --- | --- | --- | --- |
| 0 | 2725 | 2600 cold account access (95.4%) | ~11000 |
| 1 | 9425 | 2600 cold + 9000 value transfer, less the 2300 stipend the `STOP` target returns unused (98.7%) | ~3180 |

The remaining 125 gas per iteration is the loop itself: 48 for `keccak256` over
the 85-byte CREATE2 preimage and 77 for stack and memory bookkeeping.

Deploy at least that many targets per transaction, and keep the walk inside the
deployed range: at `callValue = 1`, a salt that was never deployed is charged
the 25000 gas account-creation cost instead of a plain cold access, which
measures a different benchmark (`account_mode = NON_EXISTING_ACCOUNT`).
