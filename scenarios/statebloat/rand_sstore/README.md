# 🔥 Random SSTORE State Bloater

This scenario maximizes state growth by performing the maximum number of SSTORE operations per block using random key distribution.

## 🛠️ Contract Compilation

### Prerequisites
- Solidity compiler (solc) version 0.8.30 or compatible
- Go 1.16+ (for go:embed directive)

### Compiling the Contract

To compile the SSTOREStorageBloater contract:

```bash
cd scenarios/statebloat/rand_sstore_bloater/contract
solc --optimize --optimize-runs 200 --combined-json abi,bin SSTOREStorageBloater.sol
```

### Extracting ABI and Bytecode

The compilation output is in JSON format. Extract the ABI and bytecode:

```bash
# Extract ABI (already done)
jq -r '.contracts["SSTOREStorageBloater.sol:SSTOREStorageBloater"].abi' < output.json > SSTOREStorageBloater.abi

# Extract bytecode (already done)
jq -r '.contracts["SSTOREStorageBloater.sol:SSTOREStorageBloater"].bin' < output.json > SSTOREStorageBloater.bin
```

### Regenerating Go Bindings (Optional)

If you need to regenerate the Go bindings:

```bash
abigen --abi SSTOREStorageBloater.abi --bin SSTOREStorageBloater.bin --pkg contract --out SSTOREStorageBloater.go
```

**Note**: The Go scenario code uses `go:embed` directives to automatically include the ABI and bytecode files at compile time.

## How it Works

1. **Contract Deployment**: Deploys an optimized `SSTOREStorageBloater` contract that uses assembly for minimal overhead
2. **Two-Stage Process**:
   - **Stage 1**: Creates new storage slots (0 → non-zero transitions)
   - **Stage 2**: Updates existing storage slots (non-zero → non-zero transitions)
3. **Key Distribution**: Uses curve25519 prime multiplication to distribute keys across the entire storage space, maximizing trie node creation
4. **Adaptive Gas Estimation**: Dynamically adjusts gas estimates based on actual usage to maximize slots per transaction

## ⛽ Gas Cost Breakdown

### Actual Gas Costs (Measured)
The actual gas cost per SSTORE operation is higher than the base opcode cost due to additional overhead:

**For New Slots (0 → non-zero):**
- Base SSTORE cost: 20,000 gas
- Assembly loop overhead per iteration:
  - MULMOD for key calculation: ~8 gas
  - TIMESTAMP calls: ~2 gas
  - AND operation: ~3 gas
  - Loop control (JUMPI, LT, ADD): ~10 gas
- **Total: ~22,000 gas per slot**

**For Updates (non-zero → non-zero):**
- Base SSTORE cost: 5,000 gas
- Same assembly overhead: ~2,000 gas
- **Total: ~7,000 gas per slot**

**Transaction Overhead:**
- Base transaction cost: 21,000 gas
- Function selector matching: ~100 gas
- ABI decoding (uint256 parameter): ~1,000 gas
- Contract code loading: ~2,600 gas
- Memory allocation: ~1,000 gas
- Function dispatch: ~300 gas
- Return handling: ~1,000 gas
- Safety margin: ~73,000 gas
- **Total: ~100,000 gas overhead**

Example with 30M gas limit block (97% utilization):
- Stage 1: ~1,300 new slots per block
- Stage 2: ~4,100 slot updates per block

## 🚀 Usage

### Build
```bash
go build -o bin/spamoor cmd/spamoor/main.go
```

### Run
```bash
./bin/spamoor --privkey <PRIVATE_KEY> --rpchost http://localhost:8545 rand_sstore_bloater [flags]
```

#### Flags
- `--basefee` - Base fee per gas in gwei (default: 10)
- `--tipfee` - Tip fee per gas in gwei (default: 2)

### Example
```bash
./bin/spamoor --privkey ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 \
  --rpchost http://localhost:8545 rand_sstore_bloater
```

## 📊 Deployment Tracking

The scenario tracks all contract deployments and storage operations in `deployments_sstore_bloating.json`. This file enables future scenarios to perform targeted SLOAD operations on known storage slots.

### File Format
```json
{
  "0xContractAddress": {
    "storage_rounds": [
      {
        "block_number": 123,
        "timestamp": 1234567890
      },
      ...
    ]
  }
}
```