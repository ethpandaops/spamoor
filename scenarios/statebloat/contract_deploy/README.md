# Contract Deployment State Bloat

This scenario deploys contracts that are exactly 24kB in size (EIP-170 limit) to maximize state growth while minimizing gas cost.

## How it Works

1. Generates a contract with exactly 24,576 bytes of runtime code
2. Deploys the contract using CREATE
3. Each deployment adds:
   - 24,576 bytes of runtime code
   - Account trie node
   - Total state growth: ~24.7kB per deployment

## Gas Cost Breakdown

- 32,000 gas for CREATE
- 20,000 gas for new account
- 200 gas per byte for code deposit (24,576 bytes)
- Total: 4,967,200 gas per deployment

## Usage

```bash
spamoor statebloat/contract-deploy [flags]
```

### Configuration

#### Base Settings (required)
- `--privkey` - Private key of the sending wallet
- `--rpchost` - RPC endpoint(s) to send transactions to

#### Volume Control (either -c or -t required)
- `-c, --count` - Total number of contracts to deploy
- `-t, --throughput` - Contracts to deploy per slot
- `--max-pending` - Maximum number of pending deployments

#### Transaction Settings
- `--basefee` - Max fee per gas in gwei (default: 20)
- `--tipfee` - Max tip per gas in gwei (default: 2)
- `--rebroadcast` - Seconds to wait before rebroadcasting (default: 120)

#### Wallet Management
- `--max-wallets` - Maximum number of child wallets to use
- `--refill-amount` - ETH amount to fund each child wallet (default: 5)
- `--refill-balance` - Minimum ETH balance before refilling (default: 2)
- `--refill-interval` - Seconds between balance checks (default: 300)

## Example

Deploy 10 contracts:
```bash
spamoor statebloat/contract-deploy -p "<PRIVKEY>" -h http://localhost:8545 -c 10
```

Deploy 1 contract per slot:
```bash
spamoor statebloat/contract-deploy -p "<PRIVKEY>" -h http://localhost:8545 -t 1
```

## Testing with Anvil

1. Start Anvil:
```bash
anvil
```

2. Run the scenario:
```bash
spamoor statebloat/contract-deploy -p "<PRIVKEY>" -h http://localhost:8545 -c 1
```

3. Verify state growth:
```bash
# Get the contract address from the deployment receipt
# Then check the code size
cast code <CONTRACT_ADDRESS> --rpc-url http://localhost:8545
```

The code size should be exactly 24,576 bytes. 