# Example1 Scenario

This is a comprehensive example scenario that demonstrates the key patterns and best practices for developing Spamoor scenarios. It showcases contract deployment, bound transactions, proper wallet management, and all essential scenario components.

## Features Demonstrated

### Contract Deployment
- **Automatic deployment** of a SimpleStorage contract during scenario initialization
- **Well-known wallet usage** for consistent deployment addresses
- **BuildBoundTx pattern** with abigen-generated contract bindings
- **Synchronous deployment** with confirmation waiting

### Contract Interactions  
- **Multiple operation types**: setValue, increment, and incrementBy calls
- **Proper BuildBoundTx usage** for all contract interactions
- **Transaction cycling** between different operations
- **Randomizable parameters** for varied testing

### Wallet Management
- **Configurable wallet pool** based on transaction volume
- **Well-known deployer wallet** for consistent deployment addresses  
- **Optimal wallet selection** using pending transaction count strategy
- **Automatic wallet funding** with configurable amounts

### Transaction Submission
- **RunTransactionScenario helper** for standardized execution
- **Proper onComplete handling** for transaction counting
- **Rebroadcast support** for reliability
- **Error handling and nonce management**

### Configuration
- **Comprehensive CLI flags** for all scenario options
- **YAML configuration support** for complex setups
- **Sensible defaults** for easy testing
- **Validation** of required parameters

## Usage Examples

### Basic Usage
```bash
# Deploy contract and send 50 transactions
spamoor _example1 \
  --privkey "0x1234567890abcdef..." \
  --rpchost "http://localhost:8545" \
  --count 50 \
  --initial-value 100

# Continuous operation with 5 tx/slot
spamoor _example1 \
  --privkey "0x1234567890abcdef..." \
  --rpchost "http://localhost:8545" \
  --throughput 5 \
  --randomize-values
```

### Advanced Configuration
```bash
# High-throughput with custom gas settings
spamoor _example1 \
  --privkey "0x1234567890abcdef..." \
  --rpchost "http://localhost:8545" \
  --throughput 20 \
  --max-wallets 50 \
  --basefee 30 \
  --tipfee 5 \
  --max-increment 100 \
  --randomize-values \
  --log-txs
```

### YAML Configuration
```yaml
# example1-config.yaml
total_count: 100
throughput: 10
initial_value: 42
max_increment: 50
randomize_values: true
base_fee: 25
tip_fee: 3
log_txs: true
```

```bash
# Use with daemon mode
spamoor-daemon \
  --privkey "0x1234567890abcdef..." \
  --rpchost "http://localhost:8545" \
  --startup-spammer example1-config.yaml
```

## Configuration Options

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--count`, `-c` | uint64 | 0 | Total number of transactions to send |
| `--throughput`, `-t` | uint64 | 10 | Number of transactions to send per slot |
| `--max-pending` | uint64 | 0 | Maximum number of pending transactions |
| `--max-wallets` | uint64 | 0 | Maximum number of child wallets to use |
| `--rebroadcast` | uint64 | 1 | Enable reliable rebroadcast system |
| `--basefee` | uint64 | 20 | Max fee per gas (gwei) |
| `--tipfee` | uint64 | 2 | Max tip per gas (gwei) |
| `--initial-value` | uint64 | 100 | Initial value for the storage contract |
| `--max-increment` | uint64 | 10 | Maximum increment value for random increments |
| `--randomize-values` | bool | true | Use random values for contract interactions |
| `--timeout` | string | "" | Timeout for the scenario (e.g. '1h', '30m', '5s') |
| `--client-group` | string | "" | Client group to use for sending transactions |
| `--log-txs` | bool | false | Log all submitted transactions |

## Contract Details

The scenario deploys a `SimpleStorage` contract with the following functionality:

- **setValue(uint256)**: Set the stored value
- **getValue()**: Get the current stored value (view function)
- **increment()**: Increment the stored value by 1
- **incrementBy(uint256)**: Increment the stored value by specified amount
- **owner**: Address of the contract deployer

Events emitted:
- **ValueSet(address indexed setter, uint256 oldValue, uint256 newValue)**
- **ValueIncremented(address indexed incrementer, uint256 newValue)**

## Development Patterns Demonstrated

### 1. Proper Scenario Structure
```go
// Standard scenario structure with all required components
type Scenario struct {
    options    ScenarioOptions  // Configuration
    logger     *logrus.Entry    // Structured logging
    walletPool *spamoor.WalletPool // Wallet management
    contractAddr common.Address  // Deployed contract
}
```

### 2. Contract Deployment Pattern
```go
// Well-known wallet for deployment
s.walletPool.AddWellKnownWallet(&spamoor.WellKnownWalletConfig{
    Name: "deployer",
    RefillAmount: utils.EtherToWei(uint256.NewInt(10)),
    RefillBalance: utils.EtherToWei(uint256.NewInt(5)),
})

// BuildBoundTx with abigen deployment
deploymentTx, err := deployerWallet.BuildBoundTx(ctx, txMetadata, 
    func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
        _, deployTx, _, err := contract.DeploySimpleStorage(
            transactOpts, client.GetEthClient(), initialValue)
        return deployTx, err
    })
```

### 3. Contract Interaction Pattern
```go
// BuildBoundTx with abigen contract calls
tx, err := wallet.BuildBoundTx(ctx, txMetadata,
    func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
        return storageContract.SetValue(transactOpts, value)
    })
```

### 4. Proper Transaction Submission
```go
err = txpool.SendTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
    Client: client,
    Rebroadcast: true,
    OnComplete: func(tx *types.Transaction, receipt *types.Receipt, err error) {
        onComplete() // CRITICAL: Must call for scenario counting
    },
})
```

## Learning Objectives

This example scenario teaches:

1. **Scenario Architecture**: Complete scenario structure and lifecycle
2. **Wallet Management**: Well-known wallets, pool configuration, selection strategies
3. **Contract Deployment**: Proper deployment patterns with BuildBoundTx
4. **Contract Interactions**: Safe contract calls using abigen bindings
5. **Transaction Submission**: All submission patterns and callback handling
6. **Configuration Management**: CLI flags, YAML config, validation
7. **Error Handling**: Proper error handling and nonce management
8. **Logging**: Structured logging with transaction tracking
9. **Best Practices**: Following all Spamoor development best practices

## Testing

```bash
# Test with DevNet
make devnet-run

# In another terminal, test the scenario
spamoor _example1 \
  --privkey "3fd98b5187bf6526734efaa644ffbb4e3670d66f5d0268ce0323ec09124bff61" \
  --rpchost "http://localhost:8545" \
  --count 10 \
  --log-txs

# View results in block explorer at http://localhost:port
```

This example provides a complete foundation for developing new Spamoor scenarios with all the essential patterns and best practices.