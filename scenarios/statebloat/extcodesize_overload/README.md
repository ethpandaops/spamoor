# 🔍 EXTCODESIZE Overload Attack

This scenario performs an attack that maximizes EXTCODESIZE calls within a single transaction to stress-test Ethereum clients and measure gas efficiency of state access operations.

## How it Works

1. **Gas Calculation**: Calculates the maximum number of EXTCODESIZE calls that fit in a block:
   - Formula: `(blockGasLimit - 21000) / 2600`
   - 21,000 gas: Base transaction cost
   - 2,600 gas: Cold storage access cost per EXTCODESIZE call

2. **Contract Loading**: Loads contract addresses from `deployments.json`
   - Uses the first private key's contracts from the deployment file

3. **Bytecode Generation**: Creates optimized bytecode that performs EXTCODESIZE calls:
   ```
   For each contract address:
   - PUSH20 <address>  (0x73 + 20 bytes)
   - EXTCODESIZE       (0x3B)
   - POP               (0x50)
   ```

## 🚀 Usage

### Prerequisites
Ensure you have a `deployments.json` file with contract addresses:
```json
{
  "private_key_hash_1": ["0x1234...", "0x5678..."],
  "private_key_hash_2": ["0xabcd...", "0xefgh..."]
}
```

### Basic Usage
```bash
./bin/spamoor --privkey <PRIVATE_KEY> --rpchost http://localhost:8545 extcodesize-overload
```

#### Flags
- `--basefee` - Base fee per gas in gwei (default: 20)
- `--tipfee` - Tip fee per gas in gwei (default: 2)

### Example with Anvil
```bash
# Start Anvil
anvil

# Run the attack
./bin/spamoor --privkey ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 \
  --rpchost http://localhost:8545 extcodesize-overload 
```
