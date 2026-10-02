# Bloated EOA Storage

Delegate an EOA with deep (bloated) storage to an attack contract via EIP-7702 and hammer its storage with SLOAD or SSTORE loops. This is a port of the EELS bloatnet benchmarks `test_sload_bloated` and `test_sstore_bloated` (`tests/benchmark/stateful/bloatnet/`).

## How it works

1. **Deploy** the attack contract (compiled from the geas sources in [`contract/`](./contract)) from a child wallet.
2. **Delegate** the bloated EOA to the attack contract with a setcode transaction signed by `--authority-privkey`.
3. **Attack**: child wallets send max-gas transactions to the bloated EOA, which now runs the attack contract's code.

The storage of the bloated EOA is assumed to be filled linearly with `s[i] = i` (as in EELS).

The two modes mirror the EELS runtimes exactly and differ in how the slot range is driven:

- **`sload`** (`test_sload_bloated`): the loop counter lives in slot 0. The delegation tx seeds it from calldata; each attack tx (no calldata) loads it, does `SLOAD(s[counter])`, `SLOAD(s[counter+1])`, … while `gas > --loop-gas-threshold`, then writes the next counter back to slot 0 so consecutive transactions hit untouched slots. Slot 0 is overwritten by the counter.
- **`sstore`** (`test_sstore_bloated`): there is no on-chain counter. Each attack tx carries its own disjoint `[start, end)` range as calldata (`start ‖ end`, 32 bytes each) and loops `SSTORE` while `counter+1 < end`. Ranges are derived from the transaction index so they never overlap, and slot 0 is never touched. `--slots-per-tx` sets the range width (default: derived from the gas limit). The delegation tx is sent with empty calldata and is a no-op.

## Contract variants

| File | Mode | Per-slot operation |
|------|------|--------------------|
| `sload.geas` | `--mode sload` | `SLOAD(s[i])` |
| `sstore_same.geas` | `--mode sstore --existing-slots` (no `--write-new-value`) | `s[i] = i` (rewrite same value) |
| `sstore_increment.geas` | `--mode sstore --write-new-value` | `s[i] = i + 1` |
| `sstore_zero.geas` | `--mode sstore --existing-slots=false` (no `--write-new-value`) | `s[i] = 0` (no-op on empty slot) |

`loop.geas` is the sload counter loop, `range_loop.geas` the sstore calldata-range loop, and `initcode.geas` the constructor. The per-slot `SSTORE`/`SLOAD` sequences are byte-identical to the EELS runtimes.

## Usage

```bash
spamoor bloated-eoa [flags]
```

## Configuration

### Base Settings (required)
- `--privkey` - Private key of the sending wallet
- `--rpchost` - RPC endpoint(s) to send transactions to
- `--authority-privkey` - Private key of the storage-bloated EOA

### Volume Control (either -c or -t required)
- `-c, --count` - Total number of attack transactions to send
- `-t, --throughput` - Attack transactions to send per slot (default: 1)
- `--max-pending` - Maximum number of pending transactions

### Attack Settings
- `--mode` - `sload` or `sstore` (default: `sload`)
- `--existing-slots` - Target existing slots starting at slot 1 (default: true). With `--existing-slots=false` the attack starts at `keccak256("random") % 2^160`, a range assumed to be empty
- `--write-new-value` - sstore mode: write `slot + 1` instead of the current value
- `--start-slot` - Override the starting slot (decimal or `0x` hex)
- `--slots-per-tx` - sstore mode: slots to write per transaction (default: 0 = derive from the gas limit)
- `--keep-counter` - sload mode: continue from the counter stored in slot 0 instead of resetting it (e.g. to resume a previous run)
- `--loop-gas-threshold` - sload mode: remaining gas at which the loop stops and persists the counter (default: 65535)

### Transaction Settings
- `--gaslimit` - Gas limit per attack transaction (default: 0 = min(block gas limit, 2^24))
- `--basefee` - Max fee per gas in gwei (default: 20)
- `--tipfee` - Max tip per gas in gwei (default: 2)
- `--rebroadcast` - Enable reliable rebroadcast system (default: 1)

### Client Settings
- `--client-group` - Client group to use for sending transactions
- `--deploy-client-group` - Client group to use for the deployment and delegation transactions (same as `--client-group` if empty)

### Debug Options
- `--log-txs` - Log all submitted transactions

## Examples

SLOAD existing slots, one max-gas transaction per slot:
```bash
spamoor bloated-eoa -p "<PRIVKEY>" -h http://rpc-host:8545 --authority-privkey "<BLOATED_EOA_PRIVKEY>" -t 1 --mode sload
```

SSTORE new values into existing slots:
```bash
spamoor bloated-eoa -p "<PRIVKEY>" -h http://rpc-host:8545 --authority-privkey "<BLOATED_EOA_PRIVKEY>" -t 1 --mode sstore --write-new-value
```

## Notes

- The bloated EOA's storage is assumed to be `s[i] = i`. If your account uses a different layout, `--existing-slots` and the sstore write values will not match its real state — set `--start-slot` and pick the mode/value accordingly.
- **sload** transactions share the counter in slot 0, so transactions in the same block execute serially against each other (matching EELS `test_sload_bloated`). Anyone can reset the counter by calling the bloated EOA with 32 bytes of calldata.
- **sstore** transactions use independent calldata ranges (no shared slot), so they can execute in parallel, matching EELS `test_sstore_bloated`. `--slots-per-tx` is sized conservatively from the gas limit; if transactions revert (e.g. because EIP-8037 state gas makes new-slot SSTOREs more expensive than estimated), lower it.
