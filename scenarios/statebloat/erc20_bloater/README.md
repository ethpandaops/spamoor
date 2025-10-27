# ERC20 Bloater Scenario

State bloating scenario that creates sequential ERC20 token balances and allowances to generate precise storage growth.

## Overview

Deploys an ERC20 token contract and systematically creates storage slots by transferring tokens and setting allowances for sequential addresses. Each address generates exactly 2 storage slots (64 bytes total), enabling predictable state growth.

## Features

- **Automatic gas estimation** with `eth_estimateGas` + 5% buffer
- **Checkpoint system** for automatic save/resume (`.erc20_bloater_checkpoint.json`)
- **Resume capability** to continue bloating existing contracts
- **Real-time progress tracking** with storage size calculations
- **Error handling** with exponential backoff retry

## Configuration

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `target_storage_gb` | float | 1.0 | Target storage size in GB |
| `target_gas_ratio` | float | 0.75 | Gas usage ratio of block limit (0.0-1.0) |
| `base_fee` | uint64 | 20 | Base fee in gwei |
| `tip_fee` | uint64 | 2 | Priority fee in gwei |
| `existing_contract` | string | "" | Existing contract address (optional) |

## Usage

```bash
# Command line
./bin/spamoor erc20_bloater \
  --rpchost=http://localhost:8545 \
  --privkey=<key> \
  --target-gb=0.1 \
  --target-gas-ratio=0.75 \
  -v

# YAML config
./bin/spamoor erc20_bloater \
  --rpchost=http://localhost:8545 \
  --privkey=<key> \
  --scenario-file=config.yaml
```

Example YAML:
```yaml
target_storage_gb: 0.1
target_gas_ratio: 0.75
base_fee: 20
tip_fee: 2
existing_contract: ""
```

## Storage Calculations

| Target | Slots | Addresses | ~Transactions | ~Blocks |
|--------|-------|-----------|---------------|---------|
| 0.01 GB | 335,544 | 167,772 | 311 | 622 |
| 0.1 GB | 3,355,443 | 1,677,721 | 3,107 | 6,214 |
| 1.0 GB | 33,554,432 | 16,777,216 | 31,069 | 62,138 |

*Based on 540 addresses/tx @ 36M gas limit, 1 tx per 2 blocks*

## Contract Compilation

```bash
cd scenarios/statebloat/erc20_bloater/contract
./compile.sh
```

Requires Solidity 0.8.22+ and abigen.
