# State Bloat Scenarios

This directory contains scenarios designed to test different vectors of state growth on Ethereum. Each scenario focuses on a specific method of increasing the state size while minimizing ETH cost.

## Available Scenarios

1. `contract-deploy` - Deploys 24kB contracts (EIP-170 limit)
2. `delegate-flag` - Adds delegate flags to funded EOAs (EIP-7702)
3. `fund-eoa` - Funds fresh EOAs with minimal ETH
4. `empty-auth` - Creates EIP-7702 authorizations for empty addresses
5. `storage-slots` - Fills new storage slots in contracts

## Testing

These scenarios can be tested using Anvil (Foundry's local Ethereum node) or any other EVM-compatible testnet. For local testing:

```bash
# Start Anvil
anvil

# Run a scenario (example)
spamoor statebloat/contract-deploy [flags]
```

Each scenario directory contains its own README with specific configuration options and testing instructions.

## Gas Efficiency Comparison

| Rank | Scenario        | Gas/Byte | Max Units in 30M Gas Block |
| ---- | --------------- | -------- | -------------------------- |
| 1    | Contract Deploy | ~202     | 6 deployments              |
| 2    | Delegate Flag   | ~232     | 960 tuples                 |
| 3    | Fund EOA        | ~267     | 1000 accounts              |
| 4    | Empty Auth      | ~289     | 767 tuples                 |
| 5    | Storage Slots   | 625      | 1500 slots                 |