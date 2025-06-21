# EOA Delegation Scenario

## Overview

The EOA Delegation scenario is designed for maximum state bloating through EIP-7702 SetCode transactions. It creates the largest possible state growth by delegating thousands of Externally Owned Accounts (EOAs) to existing contracts in a single block.

This scenario is a specialized stress-testing tool that:
- Automatically adjusts to fill 99.5% of block gas limit
- Funds EOAs with 1 wei each before delegation
- Tracks all funded EOAs in `EOAs.json` for future reference
- Handles transaction size limits by batching when necessary
- Continuously operates in a loop, creating maximum state growth

## How It Works

### Three-Phase Operation

1. **Funding Phase**: Pre-funds delegator EOAs with 1 wei each
2. **Bloating Phase**: Sends SetCode transaction(s) with maximum authorizations
3. **Analysis Phase**: Measures performance and adjusts parameters

### Key Features

- **Self-Adjusting**: Dynamically adjusts authorization count based on actual gas usage
- **Network-Aware**: Queries actual block gas limit from the network
- **Size-Aware**: Automatically splits large transactions that exceed 128KiB limit
- **Persistent Storage**: Saves all funded EOAs to `EOAs.json` for reuse

### Technical Details

- Each EOA delegation creates ~135 bytes of new state
- Uses ecrecover precompile (0x1) as default delegation target
- Targets 99.5% block utilization for maximum impact
- Handles up to ~1300 authorizations per 128KiB RLP-encoded transaction
- Transaction size limit: 128KiB (131,072 bytes) applies to the RLP-encoded transaction
- Actual authorization size in transaction: ~94 bytes per authorization (RLP-encoded)

## Usage

### Basic Usage

```bash
# Run with default settings (auto-adjusts to network) 
# Address is set to Identity precompile by default.
eoa-delegation -h <HOSTNAME:PORT> -p <PRIVATE KEY>

# Specify custom delegation target
eoa-delegation eoa-delegation --code-addr 0x1234567890123456789012345678901234567890

# Control gas prices
eoa-delegation eoa-delegation --basefee 30 --tipfee 3
```

### Command Line Options

- `--code-addr`: Contract address to delegate to (default: ecrecover precompile)
- `--basefee`: Base fee in gwei (default: 20)
- `--tipfee`: Priority fee in gwei (default: 2)
- `--client-group`: Specific client group to use
- `--rebroadcast`: Seconds between transaction rebroadcasts (default: 120)

### Output

The scenario logs detailed metrics for each iteration:

```
STATE BLOATING METRICS - Total bytes written: 130.5 KiB, Gas used: 25.2M, Block utilization: 98.7%, Authorizations: 990, Gas/auth: 25454.5, Gas/byte: 188.6, Total fee: 0.0504 ETH
```

## State Impact

This scenario creates maximum state growth by:
1. Creating new EOA accounts (funded with 1 wei)
2. Adding delegation records for each EOA
3. Updating nonce and balance for each account

## Files Created

- `EOAs.json`: Contains addresses and private keys of all funded EOAs
