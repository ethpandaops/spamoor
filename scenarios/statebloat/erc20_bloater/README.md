# ERC20 Bloater Scenario

State bloating scenario that creates sequential ERC20 token balances and allowances to generate precise storage growth.

## Overview

Deploys an ERC20 token contract and systematically creates storage slots by transferring tokens and setting allowances for sequential addresses. Each address generates exactly 2 storage slots (64 bytes total), enabling predictable state growth.

## Features

- **Automatic gas estimation** with `eth_estimateGas` + 5% buffer
- **Nonce-based contract detection** - automatically deploys or resumes based on wallet nonce
- **On-chain checkpointing** - queries contract's `nextStorageSlot` for resume progress
- **Deterministic addressing** - same seed always produces same contract address
- **Real-time progress tracking** with storage size calculations
- **Web UI compatible** - no file system dependencies

## Configuration

| Option              | Type   | Default | Description                                 |
| ------------------- | ------ | ------- | ------------------------------------------- |
| `target_storage_gb` | float  | 1.0     | Target storage size in GB                   |
| `target_gas_ratio`  | float  | 0.50    | Gas usage ratio of block limit (0.0-1.0)    |
| `base_fee`          | float  | 20.0    | Base fee in gwei                            |
| `tip_fee`           | float  | 2.0     | Priority fee in gwei                        |
| `existing_contract` | string | ""      | Manual contract address override (optional) |

## Usage

### Fresh Deployment
```bash
./bin/spamoor erc20_bloater \
  --rpchost=http://localhost:8545 \
  --privkey=<root-wallet-key> \
  --seed="my-bloat-project" \
  --target-gb=1.0 \
  --target-gas-ratio=0.50 \
  -v
```

**What happens:**
- Child wallet derived from seed (nonce=0)
- Deploys contract at deterministic address
- Starts bloating storage

### Automatic Resume
```bash
# Same command as above - just run it again!
./bin/spamoor erc20_bloater \
  --rpchost=http://localhost:8545 \
  --privkey=<root-wallet-key> \
  --seed="my-bloat-project" \
  --target-gb=1.0 \
  --target-gas-ratio=0.50 \
  -v
```

**What happens:**
- Child wallet detected (nonce>0)
- Contract address calculated via `crypto.CreateAddress(childWallet, 0)`
- Queries `nextStorageSlot` from contract
- Resumes from on-chain progress

### Multiple ERCs
```bash
# Project A
./bin/spamoor erc20_bloater --seed="project-alpha" --target-gb=10

# Project B (different contract)
./bin/spamoor erc20_bloater --seed="project-beta" --target-gb=5
```

Different seeds = different contracts at different addresses.

## YAML Configuration

```yaml
target_storage_gb: 1.0
target_gas_ratio: 0.50
base_fee: 20.0
tip_fee: 2.0
existing_contract: ""  # Optional manual override
```

## Storage Calculations

| Target  | Slots      | Addresses  | ~Transactions | ~Blocks |
| ------- | ---------- | ---------- | ------------- | ------- |
| 0.01 GB | 335,544    | 167,772    | 311           | 622     |
| 0.1 GB  | 3,355,443  | 1,677,721  | 3,107         | 6,214   |
| 1.0 GB  | 33,554,432 | 16,777,216 | 31,069        | 62,138  |

*Based on 540 addresses/tx @ 36M gas limit, 1 tx per 2 blocks*

## How It Works

### Nonce-Based Contract Detection

1. **Check child wallet nonce:**
   - `nonce == 0`: Deploy new contract
   - `nonce > 0`: Calculate existing contract address

2. **Contract address calculation:**
   ```go
   contractAddr = crypto.CreateAddress(childWallet, 0)
   ```
   Always returns the same address for the same seed.

3. **Resume from on-chain state:**
   ```go
   nextSlot, _ := contract.NextStorageSlot()
   startFrom = nextSlot
   ```

### No File Dependencies

- ✅ Contract address: Deterministically calculated from seed
- ✅ Progress: Queried from contract's `nextStorageSlot`
- ✅ Web UI: Fully compatible (no file system access needed)

## Contract Compilation

```bash
cd scenarios/statebloat/erc20_bloater/contract
./compile.sh
```

Requires Solidity 0.8.22+ and abigen.

## Troubleshooting

### "Wrong contract address calculated"

If the child wallet was used for other transactions before bloating:
```bash
./bin/spamoor erc20_bloater \
  --existing-contract=0xYourActualContractAddress \
  ...
```

Use `--existing-contract` to manually specify the address.

### "How do I find my contract address?"

Check the logs for:
- `deployed contract: 0x...` (on first run)
- `calculated contract address: 0x...` (on resume)

Or view the Config output in the web UI - the contract address is displayed there.
