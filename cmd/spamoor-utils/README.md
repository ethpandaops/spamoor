# spamoor-utils

A collection of utility commands for spamoor.

## Usage

```bash
spamoor-utils [command] [flags]
```

## Commands

### convert-eest

Convert Ethereum Execution Spec Test (EEST) fixtures to an intermediate representation that can be replayed on a normal network using the `replay-eest` scenario.

```bash
spamoor-utils convert-eest <path> [flags]
```

#### Flags

- `-o, --output` - Output file path (default: stdout)
- `-v, --verbose` - Enable verbose output
- `--trace` - Enable tracing output

#### Input Format

The command accepts a path to a directory containing EEST fixture JSON files in "blockchain_tests" format. It recursively scans all `.json` files and converts each test case.

EEST fixtures are the standard test format from the [ethereum/execution-spec-tests](https://github.com/ethereum/execution-spec-tests) repository.

#### Output Format

The output is a YAML file containing all converted test cases:

```yaml
payloads:
  - name: path/to/test/test_name[fork_variant]
    prerequisites:
      sender[1]: "0x3635c9adc5dea00000"
    txs:
      - from: deployer
        type: 2
        to: ""
        data: "0x600b380380600b5f395ff3..."
        gas: 1000000
        maxFeePerGas: 10000000000
        maxPriorityFeePerGas: 1000000000
        fixtureBaseFee: 7
      - from: sender[1]
        type: 0
        to: "$contract[1]"
        data: "0x"
        gas: 200000
        gasPrice: 10
        fixtureBaseFee: 7
    postcheck:
      contract[1]:
        storage:
          "0x00": "0xff"
      sender[1]:
        balance: "0x3635c9adc5de8e3762"
```

#### Conversion Details

**Pre-state Processing:**
- Contracts (accounts with code) are converted to deployment transactions
- Init code prefix `0x600b380380600b5f395ff3` is prepended to runtime bytecode
- Sender accounts (balance-only) are tracked for prerequisite funding

**Address Placeholders:**
- `$contract[N]` - Replaced with the Nth deployed contract address
- `$sender[N]` - Replaced with the Nth sender wallet address
- Placeholders appear in bytecode, `to` fields, and storage values

**System Contracts Filtered:**
- `0x00000000219ab540356cbb839cbe05303d7705fa` (Deposit Contract)
- `0x00000961ef480eb55e80d19ad83579a64c007002`
- `0x0000bbddc7ce488642fb579f8b00f3a590007251`
- `0x0000f90827f1c53a10cb7a02335b175320002935`
- `0x000f3df6d732807ef1319fb7b8bb8522d0beac02` (Beacon Roots)

**Gas Cost Tracking:**
- `fixtureBaseFee` is captured per transaction from the fixture's block header
- Used during replay to adjust balance checks for different network base fees
- Genesis base fee is used for deployment transactions

**Post-state Checks:**
- Storage slot values are captured for verification
- Balance checks are included for modified accounts
- Only accounts referenced as contracts or senders are checked

#### Examples

Convert all fixtures in a directory:
```bash
spamoor-utils convert-eest /path/to/fixtures -o output.yaml
```

Convert with verbose logging:
```bash
spamoor-utils convert-eest /path/to/fixtures -o output.yaml -v
```

Convert and pipe to stdout:
```bash
spamoor-utils convert-eest /path/to/fixtures > output.yaml
```

## Workflow

1. **Download EEST fixtures** from ethereum/execution-spec-tests releases
2. **Convert fixtures** using `spamoor-utils convert-eest`
3. **Replay tests** using `spamoor replay-eest`

```bash
# Download fixtures
wget https://github.com/ethereum/execution-spec-tests/releases/download/v1.0.0/fixtures.tar.gz
tar -xzf fixtures.tar.gz

# Convert to intermediate format
spamoor-utils convert-eest ./fixtures/blockchain_tests -o tests.yaml

# Replay on testnet
spamoor replay-eest -p "<PRIVKEY>" -h http://rpc:8545 -f tests.yaml
```
