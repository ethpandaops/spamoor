# StorageSpam Contract Analyzer

This utility script crawls and analyzes all transactions to the StorageSpam contract, including both successful `RandomForGas` events and failed transactions.

## Usage

```bash
# Using the shell wrapper
./scripts/storagespam-events.sh <RPC_URL> <CONTRACT_ADDRESS> [BATCH_SIZE]

# Or directly with Go
go run scripts/storagespam-events.go -rpc <RPC_URL> -contract <CONTRACT_ADDRESS> [-batch <SIZE>]
```

## Parameters

- `RPC_URL`: Ethereum RPC endpoint URL (required)
- `CONTRACT_ADDRESS`: StorageSpam contract address (required)
- `BATCH_SIZE`: Number of blocks to query at once (default: 100)

## Example

```bash
./scripts/storagespam-events.sh \
  https://rpc.perf-devnet-2.ethpandaops.io/ \
  0xFa3CE7108b73FA44a798A3aa23523c974ed5a6dE \
  50
```

## Features

- **Complete transaction analysis**: Examines all transactions to the contract, not just successful logs
- **Failed transaction detection**: Identifies and reports failed transactions
- **Backwards crawling**: Starts from the latest block and walks backwards
- **Automatic stopping**: Stops after finding no transactions for 500+ blocks
- **Batched receipt fetching**: Efficiently fetches receipts per block for better performance
- **Comprehensive statistics**: Provides success/failure rates, gas usage, and storage analysis
- **Detailed output**: Optional detailed listing of both successful events and failed transactions

## Output

The script provides:

1. **Transaction summary**:
   - Total transactions to the contract
   - Number of successful vs failed transactions
   - Success and failure rates

2. **Failed transaction analysis**:
   - List of all failed transactions with block numbers, transaction hashes, gas used/limit, and failure reasons
   - Detailed breakdown for debugging contract issues

3. **Successful events summary**:
   - Total successful RandomForGas events
   - Block and time ranges
   - Gas limit breakdown with counts and totals

4. **Gas limit analysis**:
   - Count of events per gas limit
   - Total and average loops per gas limit
   - Total gas burned across all successful events

5. **Storage usage analysis**:
   - Total storage slots created (one per loop)
   - Total storage used in bytes/KB/MB/GB (64 bytes per loop: 32 byte key + 32 byte value)
   - Average storage per successful event

6. **Efficiency metrics**:
   - Average gas per storage slot
   - Average gas per KB of storage

7. **Optional detailed listings**:
   - Comprehensive details for successful events: block number, gas limit, loops, actual gas used, and transaction hash
   - Failed transactions are always displayed when present

## Implementation Details

The script uses:
- The generated contract bindings from `scenarios/storagespam/contract/`
- Ethereum JSON-RPC to query blocks and transaction receipts
- Batched receipt fetching per block for optimal performance
- Complete transaction analysis to catch both successful and failed calls
- Event parsing using the contract ABI for successful transactions

## Performance Considerations

- Large deployments may have thousands of transactions and events
- The script uses batched receipt fetching per block for optimal RPC performance
- Use smaller batch sizes if encountering timeouts or RPC rate limits
- The script processes all data in memory, so very large datasets may require modifications
- Failed transactions are tracked separately to provide comprehensive analysis