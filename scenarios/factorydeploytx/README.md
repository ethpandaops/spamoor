# FactoryDeployTx Scenario

Deploy contracts using CREATE2 factory for deterministic address calculation.

## Description

This scenario deploys a CREATE2 factory contract and then uses it to deploy multiple contracts with deterministic addresses. It's designed to work with the [AttackController](../calltx/contract/AttackController.sol) to enable efficient random read testing.

## Configuration

### Basic Options
- `total_count`: Total number of contracts to deploy
- `throughput`: Number of deployment transactions per slot
- `max_pending`: Maximum pending transactions
- `max_wallets`: Maximum child wallets to use
- `gas_limit`: Gas limit for transactions (default: 2000000)

### Factory Options
- `factory_address`: Address of existing CREATE2 factory (optional)
- `init_code`: Hex-encoded init code of contract to deploy (required)
- `start_salt`: Starting salt value for deployments (default: 0)
- `well_known_factory`: Use well-known factory deployer wallet (default: true)

### Example Configuration

```yaml
- scenario: factorydeploytx
  name: "Deploy Large Contracts"
  config:
    seed: factory-deploy-001
    total_count: 1000
    throughput: 50
    gas_limit: 2000000
    well_known_factory: true
    init_code: "0x608060405234801561001057600080fd5b50......" # Large contract bytecode
    start_salt: 0
```

## Usage with AttackController

1. First run this scenario to deploy contracts via CREATE2 factory
2. Then use the `calltx` scenario with AttackController to perform random reads
3. The factory address and init code hash are used to calculate deployed contract addresses

## Contract Address Calculation

Deployed contracts can be predicted using:
```
address = keccak256(0xff, factory_address, salt, keccak256(init_code))
```

Where `salt` starts from `start_salt` and increments for each deployment.