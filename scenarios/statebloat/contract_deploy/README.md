# Contract Deployment State Bloat

This scenario deploys contracts that are exactly 24kB in size (EIP-170 limit) to maximize state growth while minimizing gas cost.

## How it Works

1. Generates a contract with exactly 24,576 bytes of runtime code
2. Deploys the contract using CREATE
3. Uses batch-based deployment:
   - Calculates how many contracts fit in one block based on gas limits
   - Sends a batch of transactions that fit within the block gas limit
   - Waits for a new block to be mined
   - Repeats the process ("bombards" the RPC after each block)
4. Each deployment adds:
   - 24,576 bytes of runtime code
   - Account trie node
   - Total state growth: ~24.7kB per deployment

## Gas Cost Breakdown

- 32,000 gas for CREATE
- 20,000 gas for new account
- 200 gas per byte for code deposit (24,576 bytes)
- Total: 4,967,200 gas per deployment

## Batch Strategy

The scenario automatically calculates how many contracts can fit in one block:
- Default block gas limit: 30,000,000 gas
- Gas per contract: 4,967,200 gas
- Contracts per batch: ~6 contracts per block

This ensures optimal utilization of block space while maintaining predictable transaction inclusion patterns.

## Usage

```bash
spamoor statebloat/contract-deploy [flags]
```

### Configuration

#### Base Settings (required)
- `--privkey` - Private key of the sending wallet
- `--rpchost` - RPC endpoint(s) to send transactions to

#### Volume Control
- `--max-transactions` - Maximum number of transactions to send (0 for unlimited)
- `--max-pending` - Maximum number of pending deployments
- `--block-gas-limit` - Block gas limit for batching (default: 30,000,000, 0 = use default)

#### Transaction Settings
- `--basefee` - Max fee per gas in gwei (default: 20)
- `--tipfee` - Max tip per gas in gwei (default: 2)
- `--rebroadcast` - Seconds to wait before rebroadcasting (default: 1)

#### Wallet Management
- `--max-wallets` - Maximum number of child wallets to use
- `--client-group` - Client group to use for sending transactions (default: "default")

## Examples

Deploy 30 contracts in batches:
```bash
spamoor statebloat/contract-deploy \
  --privkey "<PRIVKEY>" \
  --rpchost http://localhost:8545 \
  --max-transactions 30
```

Deploy with custom block gas limit (for testnets):
```bash
spamoor statebloat/contract-deploy \
  --privkey "<PRIVKEY>" \
  --rpchost http://localhost:8545 \
  --max-transactions 10 \
  --block-gas-limit 15000000
```

## Running Against a Local Anvil Node

### 1. Start Anvil

```bash
anvil --hardfork pectra
```

You can add other flags as needed (e.g., `--no-mining` if you want to control mining manually).

### 2. Build Spamoor (if not already built)

```bash
go build -o spamoor ./cmd/spamoor
```

### 3. Run the Scenario

```bash
./spamoor statebloat/contract-deploy \
  --privkey <YOUR_PRIVATE_KEY> \
  --rpchost http://localhost:8545 \
  --max-transactions 30
```

- Replace `<YOUR_PRIVATE_KEY>` with the private key of a funded account in your Anvil instance (Anvil prints these on startup).
- The scenario will automatically calculate the optimal batch size based on gas limits.

### 4. Mining Notes

- If you started Anvil with `--no-mining`, you will need to manually mine blocks (e.g., by calling `evm_mine` via RPC or using Anvil's console) to include the transactions.
- If you use Anvil's default mining mode, blocks will be mined automatically as transactions are sent.
- The scenario waits for each block to be mined before sending the next batch, ensuring predictable transaction ordering.

### 5. See All Flags

```bash
./spamoor statebloat/contract-deploy --help
``` 