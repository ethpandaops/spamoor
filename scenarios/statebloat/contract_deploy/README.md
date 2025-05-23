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

## Running Against a Local Anvil Node (without Go tests)

You can run this scenario directly against your own Anvil node using the Spamoor CLI, without running Go tests.

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
  --contracts-per-block 1
```

- Replace `<YOUR_PRIVATE_KEY>` with the private key of a funded account in your Anvil instance (Anvil prints these on startup).
- Adjust other flags as needed (`--contracts-per-block`, `--max-wallets`, etc.).

### 4. Mining Notes

- If you started Anvil with `--no-mining`, you will need to manually mine blocks (e.g., by calling `evm_mine` via RPC or using Anvil's console) to include the transactions.
- If you use Anvil's default mining mode, blocks will be mined automatically as transactions are sent.

### 5. See All Flags

```bash
./spamoor statebloat/contract-deploy --help
``` 