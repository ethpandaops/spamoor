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

The scenario uses block-aware batching to maximize block filling:
- Monitors actual block production in real-time (no timing assumptions)
- Queries network for actual block gas limit on startup
- Gas per contract: 4,949,468 gas (measured from actual deployments)
- Dynamically calculates transactions per block based on gas limit
- Sends transaction batches immediately when new blocks are detected
- Updates MaxPending to match block capacity for optimal throughput
- Example: With 90M gas limit → ~18 contracts per block

This ensures every block is filled to capacity without skipping any blocks. The scenario adapts to the actual block production rate of the network.

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
- `--max-transactions` - Total number of contracts to deploy (0 = infinite, default: 0)
- `--max-wallets` - Max child wallets to use (0 = root wallet only, default: 0)
- `--basefee` - Base fee per gas in gwei (default: 20)
- `--tipfee` - Tip fee per gas in gwei (default: 2)

Note: The scenario uses only the root wallet by default to simplify nonce management and ensure reliable transaction ordering. It automatically calculates the optimal number of concurrent transactions based on the network's block gas limit and monitors actual block production to send batches that fill every block to capacity.

#### Example with Anvil node
```bash
./bin/spamoor --privkey ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 \
  --rpchost http://localhost:8545 contract-deploy \
  --max-transactions 100 --basefee 25 --tipfee 5
``` 