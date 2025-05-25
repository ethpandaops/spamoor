# 🏭 Contract Deployment State Bloat

This scenario deploys contracts that are exactly 24kB in size (EIP-170 limit) to maximize state growth while minimizing gas cost.

## How it Works

1. Generates a contract with exactly 24,576 bytes of runtime code
2. Deploys the contract using CREATE with a salt that makes the bytecode unique
3. Uses batch-based deployment:
   - Calculates how many contracts fit in one block based on gas limits
   - Sends a batch of transactions that fit within the block gas limit
   - Waits for a new block to be mined
   - Repeats the process ("bombards" the RPC after each block)
4. Each deployment adds:
   - 24,576 bytes of runtime code
   - Account trie node
   - Total state growth: ~24.7kB per deployment

## ⛽ Gas Cost Breakdown

- 32,000 gas for CREATE
- 20,000 gas for new account
- 200 gas per byte for code deposit (24,576 bytes)
- **Total: 4,967,200 gas per deployment**

## Batch Strategy

The scenario automatically calculates how many contracts can fit in one block:
- Default block gas limit: 30,000,000 gas
- Gas per contract: 4,967,200 gas
- Contracts per batch: ~6 contracts per block

This ensures optimal utilization of block space while maintaining predictable transaction inclusion patterns.

## 🚀 Usage

### Build
```bash
go build -o bin/spamoor cmd/spamoor/main.go
```

### Run
```bash
./bin/spamoor --privkey <PRIVATE_KEY> --rpchost http://localhost:8545 contract-deploy [flags]
```

#### Key Flags
- `--max-transactions` - Number of contracts to deploy (0 = infinite, default: 0)
- `--max-pending` - Max concurrent pending transactions (default: 10)
- `--max-wallets` - Max child wallets to use (default: 1000)
- `--basefee` - Base fee per gas in gwei (default: 20)
- `--tipfee` - Tip fee per gas in gwei (default: 2)

#### Example with Anvil node
```bash
./bin/spamoor --privkey ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 \
  --rpchost http://localhost:8545 contract-deploy \
  --max-transactions 100 --max-pending 20 --basefee 25 --tipfee 5
``` 