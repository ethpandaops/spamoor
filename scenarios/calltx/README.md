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
- `--deploy-gas-limit` - Gas limit for deployment transaction (default: 2000000)
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