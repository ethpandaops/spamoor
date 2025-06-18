# ERC20 Max Transfers Scenario

This scenario maximizes the number of ERC20 token transfers per block to unique recipient addresses, creating state bloat through new account storage entries.

## Overview

The scenario uses deployed StateBloatToken contracts from `deployments.json` to send the maximum possible number of ERC20 transfers per block. Each transfer sends 1 token to a unique, never-before-used address, maximizing state growth.

## Features

- **Dynamic Block Gas Limit**: Fetches the actual network block gas limit before starting
- **Self-Adjusting Transfer Count**: Automatically adjusts the number of transfers based on actual gas usage
- **Unique Recipients**: Generates deterministic unique addresses for each transfer
- **Minimum Gas Fees**: Uses configured minimum gas fees (default: 10 gwei base, 5 gwei tip)
- **Round-Robin Contract Usage**: Distributes transfers across multiple deployed contracts
- **Recipient Tracking**: Saves all recipient addresses to `recipients.json` for analysis

## Configuration

### Command Line Flags

- `--basefee`: Max fee per gas in gwei (default: 10)
- `--tipfee`: Max tip per gas in gwei (default: 5)
- `--contract`: Specific contract address to use (default: rotate through all)

### YAML Configuration

```yaml
basefee: 10
tipfee: 5
contract: ""  # Empty string means use all contracts
```

## How It Works

1. **Initialization**:
   - Loads deployed contracts and private key from `deployments.json`
   - Sets up the deployer wallet (which holds all tokens)
   - Fetches network block gas limit

2. **Transfer Phase**:
   - Calculates optimal transfer count based on gas limit
   - Generates unique recipient addresses deterministically
   - Sends transfers in batches with minimal delays
   - Uses round-robin contract selection

3. **Analysis Phase**:
   - Tracks confirmed transfers and gas usage
   - Calculates actual gas per transfer
   - Adjusts transfer count for next iteration
   - Saves recipient data to file

4. **Self-Adjustment**:
   - If under target gas usage: increases transfers
   - If over target gas usage: decreases transfers
   - Aims for 99.5% block utilization

## State Growth Impact

Each successful transfer creates:
- New account entry for the recipient (~100 bytes)
- Token balance storage slot for the recipient
- Estimated state growth: 100 bytes per transfer

## Output

The scenario logs detailed metrics for each block:
- Number of transfers sent and confirmed
- Unique recipients created
- Gas usage and block utilization
- Estimated state growth
- Cumulative totals

Recipient addresses are saved to `recipients.json` with:
- Address
- Block number
- Tokens sent

## Requirements

- Deployed StateBloatToken contracts (via contract_deploy scenario)
- Deployer private key with full token supply
- Sufficient ETH for gas fees

## Example Usage

```bash
# Use default settings
./spamoor scenario --scenario erc20-max-transfers

# Custom gas fees
./spamoor scenario --scenario erc20-max-transfers --basefee 20 --tipfee 10

# Use specific contract only
./spamoor scenario --scenario erc20-max-transfers --contract 0xa513E6E4b8f2a923D98304ec87F64353C4D5C853
```
