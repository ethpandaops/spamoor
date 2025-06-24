# TaskRunner

Execute configurable task sequences with initialization and recurring execution phases.

**⚠️ IMPORTANT**: This scenario requires task configuration via `--tasks` or `--tasks-file`. It will not work without tasks defined.

## Usage

```bash
spamoor taskrunner [flags]
```

## Configuration

### Base Settings (required)
- `--privkey` - Private key of the sending wallet
- `--rpchost` - RPC endpoint(s) to send transactions to

### Task Configuration (required)
- `--tasks` - Inline task configuration (YAML/JSON string)
- `--tasks-file` - Path to task configuration file or URL (http/https)
- `--await-txs` - Send and await each transaction individually instead of batching

### Volume Control (either -c or -t required)
- `-c, --count` - Total number of task execution cycles
- `-t, --throughput` - Task execution cycles per slot
- `--max-pending` - Maximum number of running execution cycles with pending transactions

### Transaction Settings
- `--basefee` - Max fee per gas in gwei (default: 20) - inherited by all tasks
- `--tipfee` - Max tip per gas in gwei (default: 2) - inherited by all tasks  
- `--rebroadcast` - Seconds to wait before rebroadcasting (default: 120)

### Wallet Management
- `--max-wallets` - Maximum number of child wallets to use

### Client Settings
- `--client-group` - Client group to use for sending transactions

### Debug Options
- `--timeout` - Maximum scenario runtime (e.g. '1h', '30m', '5s')
- `--log-txs` - Log all submitted transactions

## Examples

Deploy contract and call function repeatedly:
```bash
spamoor taskrunner -p "<PRIVKEY>" -h http://rpc-host:8545 -c 100 \
  --tasks '{
    "execution": [
      {
        "type": "deploy",
        "name": "test", 
        "data": {"contract_code": "0x6080..."}
      },
      {
        "type": "call",
        "data": {
          "target": "{contract:test}",
          "call_data": "0x6057361d{random:32}"
        }
      }
    ]
  }'
```

Load tasks from file:
```bash
spamoor taskrunner -p "<PRIVKEY>" -h http://rpc-host:8545 -t 5 \
  --tasks-file ./tasks.yaml
```

Load tasks from URL:
```bash
spamoor taskrunner -p "<PRIVKEY>" -h http://rpc-host:8545 -c 1000 \
  --tasks-file https://example.com/test-tasks.yaml
```

## Advanced Usage

### Transaction Execution Modes

**Batch Mode (default)**: All tasks in an execution cycle are built first, then sent as a batch for maximum throughput. This creates temporary nonce gaps during building.

**Await Mode (`--await-txs`)**: Each task is built, sent, and awaited individually before proceeding to the next task. This ensures:
- No nonce gaps during execution (critical for concurrent operations)
- Sequential execution within each cycle
- Each task waits for confirmation before the next task is built
- Easier debugging and transaction tracking
- Lower throughput but guaranteed ordering and no nonce conflicts

Example with sequential execution:
```bash
spamoor taskrunner -p "<PRIVKEY>" -h http://rpc-host:8545 -c 100 \
  --await-txs --tasks-file ./sequential-tasks.yaml
```

### Task Configuration Format

Tasks are configured in YAML or JSON with two phases. All tasks automatically inherit the global transaction settings (`--basefee`, `--tipfee`) from the scenario options:

```yaml
# Initialization tasks (run once)
init:
  - type: deploy
    name: factory
    data:
      contract_file: "./factory.bin"
      gas_limit: 3000000

# Execution tasks (run repeatedly)  
execution:
  - type: deploy
    name: token
    data:
      contract_file: "./erc20.bin"
      contract_args: "0x000000000000000000000000000000000000000000000000000000000000001e"
      
  - type: call
    data:
      target: "{contract:token}"
      call_abi_file: "./erc20-abi.json" 
      call_fn_name: "transfer"
      call_args: ["{randomaddr}", "{random:1000000}"]
```

### Task Types

#### Deploy Task
Deploys a contract and optionally registers it in the contract registry.

```yaml
- type: deploy
  name: token        # Optional: register contract as "token"
  data:
    contract_code: "0x608060405234801561001057600080fd5b50..."  # Inline bytecode
    # OR
    contract_file: "./contracts/erc20.bin"                       # File path or URL (http/https)
    
    contract_args: "0x000000000000000000000000000000000000000000000000000000000000001e"  # Constructor args
    gas_limit: 2000000    # Gas limit for deployment
    amount: 1000          # ETH to send in gwei (optional)
```

#### Call Task
Calls a function on a deployed contract.

```yaml
- type: call
  name: transfer     # Optional task name
  data:
    target: "{contract:token}"        # Contract reference
    # OR
    target: "0x742d35Cc9B3Ed5F9..."  # Direct address
    
    # Method 1: Raw calldata
    call_data: "0xa9059cbb000000000000000000000000742d35cc9b3ed5f9..."
    
    # Method 2: ABI + function
    call_abi: '[{"inputs":[...],"name":"transfer",...}]'  # Inline ABI
    # OR
    call_abi_file: "./contracts/erc20-abi.json"           # ABI file or URL (http/https)
    call_fn_name: "transfer"                              # Function name
    call_args: ["{randomaddr}", "1000000"]                # Function arguments
    
    gas_limit: 100000
    amount: 0         # ETH to send in gwei (optional)
```

### Placeholders

TaskRunner supports dynamic placeholders in arguments:

- `{random}` - Random uint256 value
- `{random:N}` - Random value between 0 and N
- `{randomaddr}` - Random Ethereum address
- `{contract:name}` - Reference to deployed contract
- `{txid}` - Current transaction ID
- `{stepid}` - Current step index within execution cycle

Example:
```yaml
call_args: ["{randomaddr}", "{random:1000000}", "{contract:token}"]
```

### Contract Registry

The contract registry maintains deployed contract addresses across phases:

1. **Init Phase**: Contracts deployed here are available throughout execution
2. **Execution Phase**: Each cycle has its own registry that inherits from init
3. **Contract References**: Use `{contract:name}` to reference deployed contracts

### Wallet Management

TaskRunner uses different wallet strategies for each phase:

1. **Init Phase**: Uses a well-known wallet (`taskrunner-init`) for all initialization tasks
   - Ensures init contracts have predictable addresses
   - Avoids conflicts with numbered wallets used in execution
   - Sequential execution with consistent wallet state

2. **Execution Phase**: Uses numbered wallets distributed across cycles
   - Each execution cycle can use different wallets for parallel processing
   - Follows standard spamoor wallet selection patterns
   - Enables high throughput execution

**Registry Inheritance:**
```
Init Registry (Global)
├── Contract "factory" → 0x123...
└── Contract "helper" → 0x456...

Execution Registry (Per Cycle)  
├── Inherits: factory, helper
├── Contract "token" → 0x789...    # Deployed this cycle
└── Contract "nft" → 0xabc...      # Deployed this cycle
```

### Advanced Examples

#### Multi-Contract Deployment
```yaml
init:
  - type: deploy
    name: factory
    data:
      contract_file: "./factory.bin"
      gas_limit: 3000000

execution:
  - type: call
    data:
      target: "{contract:factory}"
      call_abi_file: "./factory-abi.json"
      call_fn_name: "createToken"
      call_args: ["{random}", "{randomaddr}"]
      
  - type: call  
    data:
      target: "{contract:factory}"
      call_fn_name: "mint"
      call_args: ["{randomaddr}", "{random:1000000}"]
```

#### Complex ERC20 Testing
```yaml
init:
  - type: deploy
    name: token
    data:
      contract_file: "./erc20.bin"
      contract_args: "0x0000000000000000000000000000000000000000000000056bc75e2d630eb20000"

execution:
  - type: call
    data:
      target: "{contract:token}"
      call_abi_file: "./erc20-abi.json"
      call_fn_name: "transfer" 
      call_args: ["{randomaddr}", "{random:1000000}"]
      
  - type: call
    data:
      target: "{contract:token}"
      call_abi_file: "./erc20-abi.json"
      call_fn_name: "approve"
      call_args: ["{randomaddr}", "{random:2000000}"]
```

## When to Use TaskRunner

**Use TaskRunner when:**
- Testing complex multi-step contract interaction patterns
- Need to maintain state between transaction cycles
- Require flexible, configurable transaction sequences
- Other scenarios are too rigid or specific

**Don't use TaskRunner for:**
- Simple EOA transfers (use `eoatx`)
- Basic contract deployments (use `deploytx`) 
- Standard token transfers (use `erctx`)
- Single-function contract calls (use `calltx`)

## Error Handling

- **Configuration Errors**: Invalid YAML/JSON or missing required fields cause immediate failure
- **Task Validation**: Each task is validated before execution starts
- **Transaction Failures**: Failed transactions are logged but don't stop the scenario
- **Contract References**: Unknown contract references cause immediate task failure

## Performance Tips

1. **Minimize Init Tasks**: Init phase runs sequentially using a single well-known wallet - keep it minimal
2. **Use Multiple Wallets**: Distribute execution load across `--max-wallets` for better throughput  
3. **Choose Execution Mode**: 
   - Use **batch mode** (default) for maximum throughput
   - Use **await mode** (`--await-txs`) when tasks must execute sequentially or for debugging
4. **Wallet Strategy**: Init phase isolation prevents wallet conflicts during high-throughput execution
5. **Optimize Gas Limits**: Set appropriate gas limits to avoid failures

## Troubleshooting

### "no tasks configuration provided"
- **Cause**: Neither `--tasks` nor `--tasks-file` was specified
- **Fix**: Provide task configuration via one of these flags

### "contract 'name' not found in registry"
- **Cause**: Referencing a contract that wasn't deployed or named incorrectly
- **Fix**: Ensure contract is deployed in init phase or earlier in execution phase

### "failed to parse tasks as YAML or JSON"
- **Cause**: Malformed configuration syntax
- **Fix**: Validate YAML/JSON syntax using online validators

### High transaction failure rates
- **Cause**: Insufficient gas limits or network congestion
- **Fix**: Increase gas limits or reduce throughput