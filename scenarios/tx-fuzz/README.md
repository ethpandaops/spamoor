# Transaction Fuzzer

Fuzzes the **transaction layer** by sending well-formed transactions across all
EVM transaction types with randomized envelope fields. Where `evm-fuzz` fuzzes
EVM *execution* (random deployed bytecode), `tx-fuzz` fuzzes the *transaction
envelope* itself — type, calldata, access lists, authorizations, blobs, and
targets.

All transactions are valid/well-formed and flow through spamoor's managed wallet
pool (nonces, rebroadcast, receipt tracking). Transactions may revert or have
non-applying sub-components (e.g. invalid EIP-7702 authorizations), which is
expected — they still mine. Genuinely *invalid* submissions (bad nonce, malformed
RLP, etc.) are out of scope for this scenario; they require an out-of-pool
fire-and-forget harness.

## Usage

```bash
spamoor tx-fuzz [flags]
```

## What gets fuzzed

- **Transaction type** — legacy (0), access list (2930), dynamic fee (1559),
  blob (4844), set code (7702), chosen uniformly from the enabled set per tx.
- **Calldata / initcode** — mixed empty / small / large random payloads, with an
  occasional near-`--max-call-data` payload.
- **Recipient** — recoverable child wallets, contracts deployed by earlier
  fuzzed creations (so calldata executes against real bytecode), the zero
  address, precompile addresses, system contracts (beacon roots,
  withdrawal/consolidation queues, history storage), and fully random
  addresses. Legacy/2930/1559 are sometimes contract creations (`to == nil`).
- **Value** — zero, 1 wei, or a small random amount (up to ~0.001 ETH).
  Meaningful value is only sent to pool wallets so funds stay reclaimable.
- **Gas limit** — fuzzed across buckets: the exact intrinsic-gas floor
  (out-of-gas on the first unit of work), floor plus small slack, mid-range, and
  the configured `--gaslimit` cap. Always at or above the intrinsic floor so
  every tx stays includable.
- **Access lists** — random EIP-2930 tuples (often empty), with an occasional
  near-limit list (max entries, each fully loaded with storage keys).
- **Authorizations** — fuzzed EIP-7702 auth lists signed with throwaway keys
  (decoupled from pool nonce accounting), with mixed valid/garbage chain IDs and
  nonces. A fraction of setcode txs additionally carry an authorization that
  actually *applies* (fresh ephemeral authority with nonce 0, real chain id,
  delegate with code) and call the authority so the delegated code executes.
- **Blobs** — random blob sidecars plus known edge cases (all-zero, repeated,
  duplicate commitments), occasionally the full `--max-blobs` budget.

## Configuration

### Volume Control (either -c or -t required)
- `-c, --count` - Total number of transactions to send
- `-t, --throughput` - Transactions to send per slot (default: 50)
- `--max-pending` - Maximum number of concurrent pending submissions (default: 0
  = auto-sized from throughput and wallet count). Workers no longer hold a
  pending slot until confirmation, so throughput is governed by `-t` alone.

### Transaction Settings
- `--basefee` - Max fee per gas in gwei (default: 20)
- `--tipfee` - Max tip per gas in gwei (default: 2)
- `--gaslimit` - Gas limit per transaction (default: 500000)
- `--rebroadcast` - Seconds to wait before rebroadcasting (default: 30)
- `--unstuck-time` - Seconds to wait for a fuzzed tx before replacing it with a cancel tx to free the nonce (default: 60, 0 disables)
- `--timeout` - Maximum duration to run (e.g. '1h', '30m', '5s')

### Fuzzing Configuration
- `--tx-types` - Comma list of tx types to fuzz: `legacy,accesslist,dynfee,blob,setcode` (or `all`, default)
- `--payload-seed` - Custom hex seed for reproducible fuzzing (e.g. 0x1234abcd)
- `--tx-id-offset` - Start fuzzing from a specific transaction ID (default: 0)
- `--max-call-data` - Maximum calldata/initcode size in bytes (default: 1024)
- `--max-access-list` - Maximum access list entries and storage keys (default: 5)
- `--max-auth-list` - Maximum EIP-7702 authorizations per setcode tx (default: 5)
- `--max-blobs` - Maximum blob sidecars per blob tx (default: 3)
- `--max-blob-tx-per-slot` - Maximum blob txs submitted per slot (default: 4,
  0 = unlimited). Excess blob txs are re-fuzzed as non-blob txs instead of
  hoarding pending slots on blob capacity.

### Blob Format (EIP-4844 / EIP-7594)
Blob txs use the v0 KZG-proof sidecar before Fulu and the v1 (cell-proof) wrapper
after, so they are always submitted in the format the network expects. The Fulu
activation timestamp is taken from the daemon's global config (`fulu_activation`)
when available, otherwise from this flag.
- `--fulu-activation` - Unix timestamp of the Fulu activation. **Default: 0** (Fulu
  active since genesis → v1 blobs), which is correct for current post-Fusaka
  networks. Set this to a future timestamp to send v0 blobs on a pre-Fulu chain.
- `--blob-v1-percent` - Percentage of blob txs sent with the v1 wrapper after Fulu (default: 100)

### Wallet Management
- `--max-wallets` - Maximum number of child wallets to use

The wallet pool is partitioned into a blob and a non-blob class: each wallet
only ever sends one class of transactions. go-ethereum reserves every account
to a single subpool, so mixing blob and non-blob txs on one wallet would cause
"address already reserved" rejections.

### Client Settings
- `--client-group` - Client group to use for sending transactions

### Debug Options
- `--log-txs` - Log all submitted transactions

## Examples

Fuzz all transaction types:
```bash
spamoor tx-fuzz -p "<PRIVKEY>" -h http://rpc-host:8545 -t 50
```

Fuzz only blob and setcode transactions for an hour:
```bash
spamoor tx-fuzz -p "<PRIVKEY>" -h http://rpc-host:8545 -t 20 --tx-types blob,setcode --timeout 1h
```

Reproducible run with a fixed seed:
```bash
spamoor tx-fuzz -p "<PRIVKEY>" -h http://rpc-host:8545 -c 1000 --payload-seed 0x1234abcd
```
