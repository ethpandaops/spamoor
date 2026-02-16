# Replay EEST

Replay Ethereum Execution Spec Test (EEST) fixtures from an intermediate representation format.

This scenario executes test cases converted from EEST fixtures, deploying contracts and sending transactions while validating post-execution state.

## Usage

```bash
spamoor replay-eest [flags]
```

## Configuration

### Base Settings (required)
- `--privkey` - Private key of the sending wallet
- `--rpchost` - RPC endpoint(s) to send transactions to
- `-f, --payload` - Path or URL to the YAML payload file

### Volume Control
- `-c, --count` - Total number of test cases to run (default: all)
- `-t, --throughput` - Test cases to run per slot (default: 1)
- `--max-pending` - Maximum number of pending test cases (default: 10)

### Transaction Settings
- `--basefee` - Max fee per gas in gwei (default: 20)
- `--tipfee` - Max tip per gas in gwei (default: 2)
- `--rebroadcast` - Enable rebroadcasting (default: 1)

### Wallet Management
- `--max-wallets` - Maximum number of child wallets to use (default: 100)
- `--refill-amount` - ETH amount to fund each child wallet (default: 50)
- `--refill-balance` - Minimum ETH balance before refilling (default: 10)

### Client Settings
- `--client-group` - Client group to use for sending transactions

### Debug Options
- `-v, --verbose` - Enable verbose output
- `--log-txs` - Log individual transactions
- `--trace` - Enable tracing output
- `--timeout` - Timeout for the scenario (e.g., '1h', '30m')

## Payload Format

The payload file is a YAML file with the following structure:

```yaml
payloads:
  - name: test_case_name
    prerequisites:
      sender[1]: "0x3635c9adc5dea00000"  # Initial balance for sender 1
    txs:
      - from: deployer
        type: 2
        data: "0x600b380380600b5f395ff3..."
        gas: 1000000
        maxFeePerGas: 10000000000
        maxPriorityFeePerGas: 1000000000
        fixtureBaseFee: 7
      - from: sender[1]
        type: 0
        to: "$contract[1]"
        data: "0x..."
        gas: 200000
        gasPrice: 10
        fixtureBaseFee: 7
    postcheck:
      contract[1]:
        storage:
          "0x00": "0xff"
          "0x01": "0xbf"
      sender[1]:
        balance: "0x3635c9adc5de8e3762"
```

### Placeholders

- `$contract[N]` - Address of the Nth contract deployed by the deployer
- `$sender[N]` - Address of the Nth sender wallet

Placeholders can appear in:
- Transaction `to` field
- Transaction `data` field (bytecode)
- Post-check storage values

### Transaction Types

- `from: deployer` - Transactions sent by the deployer wallet (contract deployments)
- `from: sender[N]` - Transactions sent by sender wallet N

### Gas Cost Scaling

Each transaction includes a `fixtureBaseFee` field that stores the block base fee from the original fixture. During balance verification:
- The expected balance change is adjusted based on the difference between fixture and actual gas costs
- Formula: `adjustedExpected = expectedChange - (gas * fixtureBaseFee) + actualGasCost`
- This allows accurate balance verification even when network base fees differ from fixture values

### Post-checks

- `storage` - Verify storage slot values after execution
- `balance` - Verify balance (relative comparison for senders, absolute for contracts)

## Examples

Run all test cases from a local file:
```bash
spamoor replay-eest -p "<PRIVKEY>" -h http://rpc-host:8545 -f ./converted.yaml
```

Run first 10 test cases from a URL:
```bash
spamoor replay-eest -p "<PRIVKEY>" -h http://rpc-host:8545 -f https://example.com/tests.yaml -c 10
```

Run 2 test cases per slot with verbose logging:
```bash
spamoor replay-eest -p "<PRIVKEY>" -h http://rpc-host:8545 -f ./converted.yaml -t 2 --log-txs -v
```

## Converting EEST Fixtures

Use the `spamoor-utils convert-eest` command to convert EEST fixtures to the intermediate format:

```bash
spamoor-utils convert-eest /path/to/fixtures -o output.yaml
```

See the [spamoor-utils README](../../cmd/spamoor-utils/README.md) for more details.
