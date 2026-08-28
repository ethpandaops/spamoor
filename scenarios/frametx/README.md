# Frame Transactions (EIP-8141)

Send [EIP-8141](https://eips.ethereum.org/EIPS/eip-8141) frame transactions (type `0x06`) in a
mix of frame shapes, and check the per-frame receipts against what each shape should produce.

A frame transaction is not a single call but an ordered list of up to 64 frames that carry out
validation, gas payment approval and the user's operations. This scenario uses the protocol's
**default code**, so ordinary child wallets can send frame transactions with nothing deployed.

The scenario refuses to start when the chain does not implement EIP-8141: it probes with a single
frame transaction from a throwaway key before sending anything.

## Usage

```bash
spamoor frametx [flags]
```

## Frame shapes

`--shapes` takes a weighted list, e.g. `--shapes self-verify:10,atomic:2`. The default, `all`,
enables every shape that EIP-8141 defines on its own, with equal weight. The `post-tx` shapes need
EIP-7906 and are selected by name.

| Shape | Frames | What it exercises |
|---|---|---|
| `self-verify` | `[self_verify, user_op]` | The minimal mempool-legal shape; default code and the single-frame validation prefix |
| `transfer` | `[self_verify, user_op(value)]` | Value-bearing frames, `TX_VALUE_COST`, and the state gas an account creation needs |
| `batch` | `[self_verify, user_op × N]` | Per-frame cost and block packing across both gas dimensions |
| `atomic` | `[self_verify, user_op(atomic) × 2, user_op]` | Atomic batching where every frame succeeds |
| `atomic-fail` | as above, with the middle frame starved of gas | Batch rollback, and the `skipped` frame status that only a rollback produces |
| `expiry` | `[expiry_verify, self_verify, user_op]` | The expiry verifier predeploy at `0x8141` and mempool deadline handling |
| `post-tx` | `[self_verify, user_op, post_tx]` | EIP-7906 assertion frames that pass |
| `post-tx-revert` | as above, with the assertion failing | Whole-body revert: the user operation reports success but its effects are discarded |

The `post-tx` shapes use the expiry verifier predeploy as their assertion contract, so they need
nothing deployed: with an 8-byte future deadline it succeeds, and with any other calldata length it
reverts. EIP-7906 installs no predeploy, so there is no chain state saying whether it is active;
these shapes are therefore excluded from `all` and must be named explicitly. A chain without
EIP-7906 rejects frame mode 3.

Each shape declares the frame statuses a correct client must report. Confirmed transactions are
checked against that and any disagreement is logged as a warning, with a summary at the end of the
run. Use `--verify-frames=false` to send load without checking.

## Configuration

### Base Settings (required)
- `--privkey` - Private key of the sending wallet
- `--rpchost` - RPC endpoint(s) to send transactions to

### Volume Control (either -c or -t required)
- `-c, --count` - Total number of transactions to send
- `-t, --throughput` - Transactions to send per slot
- `--max-pending` - Maximum number of pending transactions

### Frame Settings
- `--shapes` - Weighted list of frame shapes to send (default: `all`)
- `--envelope` - Envelope shape to encode: `auto` (default), `full`, `keyed`, `roots`, `base`
- `--frames-per-tx` - Number of user operation frames for the `batch` shape (default: 4)
- `--user-op-gas` - Execution gas limit per user operation frame (default: 30000)
- `--verify-gas` - Execution gas limit for validation frames (default: 5000)
- `--state-gas` - Additional state gas budget per user operation frame (default: 0)
- `--expiry-offset` - Seconds ahead of now for the `expiry` shape's deadline (default: 600)
- `--verify-frames` - Check per-frame receipt statuses against each shape (default: true)
- `--skip-mempool-check` - Skip the local public-mempool policy check before sending

### Transaction Settings
- `--basefee` - Max fee per gas in gwei (default: 20)
- `--tipfee` - Max tip per gas in gwei (default: 2)
- `--amount` - Transfer amount per value-bearing frame in gwei (default: 20)
- `--random-amount` - Use random amounts, with `--amount` as the limit (default: true)
- `--random-target` - Send to random, non-existent target addresses
- `--data` - Call data for user operation frames
- `--rebroadcast` - Enable the reliable rebroadcast system (default: 1)

### Wallet Management
- `--max-wallets` - Maximum number of child wallets to use
- `--refill-amount` - ETH amount to fund each child wallet (default: 5)
- `--refill-balance` - Minimum ETH balance before refilling (default: 2)
- `--refill-interval` - Seconds between balance checks (default: 300)

### Client Settings
- `--client-group` - Client group to use for sending transactions

### Debug Options
- `-v, --verbose` - Enable verbose output
- `--log-txs` - Log all submitted transactions
- `--timeout` - Scenario timeout, e.g. `1h` or `30m`

## Envelope shapes

EIP-8141's payload is amended independently by two further EIPs, so a chain may run any of four
shapes and encoding the wrong one fails to decode entirely:

| `--envelope` | payload | fields |
|---|---|---|
| `base` | EIP-8141 alone, scalar nonce | 7 |
| `keyed` | + EIP-8250 keyed nonces | 8 |
| `roots` | + EIP-8272 recent roots | 8 |
| `full` | both | 9 |

`auto`, the default, probes the chain at startup and reports which shape it found. Set it
explicitly to exercise a shape the chain would not otherwise receive.

## Gas budgets

Frame transactions declare both gas dimensions explicitly, per frame, rather than deriving them
from one limit:

- `--verify-gas` bounds a validation frame. The default code draws no execution gas of its own, so
  this only has to cover the sender's EIP-2929 account access, which is warm. The public mempool
  caps the whole validation prefix plus signature verification at 100,000 gas.
- `--user-op-gas` bounds each user operation frame's computation.
- `--state-gas` adds state gas for frames that create state. `transfer` adds an account-creation
  budget by itself when `--random-target` is set.

## Examples

Send the default mix:

```bash
spamoor frametx -p "<PRIVKEY>" -h http://rpc-host:8545 -t 10
```

Batching only, 32 calls per transaction:

```bash
spamoor frametx -p "<PRIVKEY>" -h http://rpc-host:8545 -t 5 \
  --shapes batch --frames-per-tx 32
```

Concentrate on atomic batch rollback:

```bash
spamoor frametx -p "<PRIVKEY>" -h http://rpc-host:8545 -t 5 \
  --shapes atomic:1,atomic-fail:3 -v
```

## Troubleshooting

**"this chain does not implement frame transactions"** - the chain has not activated EIP-8141. The
probe reports which client rejected the type.

**Frame receipt mismatches** - the client's per-frame statuses differ from the shape's expectation.
This is the signal the scenario exists to produce; the warning names the frame index and both
statuses. Note that the JSON-RPC encoding of frame receipts is not specified by EIP-8141, so an
early client may report a shape the decoder does not recognize.

The expectations are deliberately strict: a mismatch keeps being reported until the client
implements the behavior.

Beyond the per-frame statuses, the scenario checks whether a `POST_TX` failure reverted the whole
execution body, which the statuses alone do not show — a frame can report success and still have
had its effects discarded, by an unrolled atomic batch or a failed `POST_TX` frame. The confirm log
reports both, as `frames:` and `durable:`.

**"validation prefix is not a recognized shape"** - the local mempool policy check rejected the
transaction before sending. Use `--skip-mempool-check` to submit anyway, which is useful against a
node with a permissive local mempool.
