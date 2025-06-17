# Set Code Transactions

Send transactions that update contract code using EIP-7702 set code authorizations. This scenario is useful for testing contract code delegation and state changes.

## Usage

```bash
spamoor setcodetx [flags]
```

## Configuration

### Base Settings (required)
- `--privkey` - Private key of the sending wallet
- `--rpchost` - RPC endpoint(s) to send transactions to

### Volume Control (either -c or -t required, unless using --max-bloating)
- `-c, --count` - Total number of transactions to send
- `-t, --throughput` - Transactions to send per slot
- `--max-pending` - Maximum number of pending transactions

### Transaction Settings
- `--basefee` - Max fee per gas in gwei (default: 20)
- `--tipfee` - Max tip per gas in gwei (default: 2)
- `--gaslimit` - Gas limit for each transaction in gwei (default: 200000)
- `--amount` - Amount to send with each transaction in gwei (default: 20)
- `--random-amount` - Use random amounts for transactions (with --amount as limit)
- `--random-target` - Use random to addresses for transactions
- `--data` - Transaction call data to send
- `--code-addr` - Code delegation target address to use for transactions
- `--random-code-addr` - Use random delegation target for transactions
- `--rebroadcast` - Seconds to wait before rebroadcasting (default: 120)

### Wallet Management
- `--max-wallets` - Maximum number of child wallets to use
- `--min-authorizations` - Minimum number of authorizations to send per transaction (default: 1)
- `--max-authorizations` - Maximum number of authorizations to send per transaction (default: 10)
- `--max-delegators` - Maximum number of random delegators to use (0 = no delegator gets reused)

### Max Bloating Mode
- `--max-bloating` - **Enable maximum state bloating mode**: Creates ~960 EOA delegations in a single block-filling transaction for maximum state growth testing

### Client Settings
- `--client-group` - Client group to use for sending transactions

### Debug Options
- `-v, --verbose` - Enable verbose output
- `--log-txs` - Log all submitted transactions
- `--trace` - Enable tracing output

## Max Bloating Mode

The `--max-bloating` flag enables a special operation mode designed for maximum blockchain state growth testing. This mode:

### Key Features
- **Self-adjusting gas targeting**: Automatically targets BLOCK_LIMIT_GAS for maximum block utilization
- **Continuous operation**: Runs indefinitely, creating new state bloat each block
- **EOA delegation focus**: Creates fresh EOA accounts and delegates them to maximize state growth
- **Dynamic optimization**: Adjusts authorization count based on actual gas usage to stay within block limits

### How It Works
1. **Funding Phase**: Creates and funds new EOA accounts (1 wei each) for delegation
2. **Bloating Phase**: Sends a single transaction with maximum EIP-7702 authorizations
3. **Analysis Phase**: Analyzes gas usage and adjusts parameters for next iteration
4. **EOA Export**: Saves all funded EOA accounts to `EOAs.json` for potential reuse

### Configuration for Max Bloating
- Uses root wallet only (`--max-wallets` is automatically set to 1)
- Count and throughput limits are ignored (continuous operation)
- Default delegate target: `0x0000000000000000000000000000000000000001` (ecrecover precompile) 
This is because this is mostly 0's (cheaper calldata) and already has code so gives a discount.
- Can override delegate with `--code-addr` for specific testing scenarios

### Performance Metrics
The mode provides detailed analytics including:
- Gas used per block
- Number of authorizations processed
- Gas efficiency per authorization
- Gas efficiency per byte of state change
- Total transaction fees

### Output Files
- `EOAs.json` - Contains all funded EOA accounts with private keys for potential reuse

## Examples

### Standard Usage
Send 100 set code transactions:
```bash
spamoor setcodetx -p "<PRIVKEY>" -h http://rpc-host:8545 -c 100
```

Send 2 set code transactions per slot with random amounts:
```bash
spamoor setcodetx -p "<PRIVKEY>" -h http://rpc-host:8545 -t 2 --random-amount
```

### Max Bloating Mode
Run continuous state bloating with self-adjusting parameters:
```bash
spamoor setcodetx -p "<PRIVKEY>" -h http://rpc-host:8545 --max-bloating
```

Run max bloating with custom delegate target:
```bash
spamoor setcodetx -p "<PRIVKEY>" -h http://rpc-host:8545 --max-bloating --code-addr 0x1234567890123456789012345678901234567890
```

## Warning

⚠️ **Max Bloating Mode**: This mode is designed for testnets and development environments. It creates significant blockchain state growth and should not be used on production networks. It will continuously consume ETH for transaction fees. 