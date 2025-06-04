# State Bloat Scenarios

This directory contains scenarios designed to test different vectors of state growth on Ethereum. Each scenario focuses on a specific method of increasing the state size while minimizing ETH cost.

## Available Scenarios

1. `contract-deploy` - Deploys 24kB contracts (EIP-170 limit)
2. `delegate-flag` - Adds delegate flags to funded EOAs (EIP-7702) (tirggered via setcodetx scenario using `--max_bloating` flag.)
3. `fund-eoa` - Funds fresh EOAs with minimal ETH (tirggered alongside setcodetx scenario using `--max_bloating` flag.)
4. `storage-slots` - Fills new storage slots in contracts
5. `extcodesize-overload` - Maximizes EXTCODESIZE calls in a single transaction

## Testing

These scenarios can be tested using Anvil (Foundry's local Ethereum node) or any other EVM-compatible testnet. For local testing:

```bash
# Start Anvil
anvil

# Run a scenario (example)
spamoor statebloat/contract-deploy [flags]
```

Each scenario directory contains its own README with specific configuration options and testing instructions.

