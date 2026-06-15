# Invalid Transaction Fuzzer

Fuzzes **transaction validation** by firing deliberately-invalid transactions at
the network and recording whether the node correctly rejects them. A node that
*accepts* a transaction this scenario built to be invalid is a potential
validation gap — the signal this fuzzer hunts for.

This is the third fuzzing layer alongside `evm-fuzz` (EVM execution) and
`tx-fuzz` (valid tx envelopes).

## How it works (out-of-pool, fire-and-forget)

Invalid transactions fundamentally conflict with spamoor's managed transaction
pool, which assumes every submitted tx is well-formed and will eventually confirm
at its assigned nonce. So this scenario deliberately **bypasses the managed
pool**:

- Transactions are signed manually and submitted via `eth_sendRawTransaction`.
- They are **fire-and-forget** — the RPC accept/reject response is the result;
  no receipt is awaited, no rebroadcast, no nonce tracking through the pool.
- Submissions come from a **single reused burner wallet** (configurable via
  `--max-wallets` for sender diversity). Because invalid txs are never mined they
  don't consume a nonce, so one wallet can fire them indefinitely — and this
  keeps tx load on the shared root wallet to a single funding transaction, which
  is what otherwise breaks in combination with other scenarios. The burner
  barely spends (most invalid txs are rejected before execution), so it's funded
  minimally and rarely.

### Unexpectedly-valid txs

Invalidity can't be guaranteed upfront when fuzzing — a mutated tx may turn out
valid. New txs are based on the burner's **pending** nonce, so if such a tx is
accepted into the mempool the nonce advances and follow-ups don't all burn on a
consumed nonce. Only **structurally-invalid** categories (marked below) raise a
finding when accepted; state-dependent categories (future/low nonce, underpriced)
can be legitimately accepted and are not flagged.

### Unstuck helper

Most invalid txs are rejected at submission and consume no on-chain nonce, so the
burner can be reused indefinitely. The exception is a tx a node accepts into its
mempool at the current nonce (e.g. underpriced). When that happens the wallet is
**unstuck before reuse**: mirroring tx-fuzz's `tryUnstuck`, it replaces the stuck
nonces with fee-bumped self-transfers until the pending nonce catches up to the
confirmed nonce.

## Invalid categories

"Finding?" marks categories that are structurally invalid regardless of state —
acceptance of one is a genuine validation gap. The others are state-dependent and
their acceptance is expected/inconclusive.

| Category | What's wrong | Finding if accepted? |
|---|---|---|
| `chainid` | signed for a different chain id | ✅ |
| `lowgas` | gas limit below the intrinsic minimum | ✅ |
| `gasoverflow` | gas limit far above the block gas limit | ✅ |
| `emptyauth` | EIP-7702 set-code tx with an empty authorization list | ✅ |
| `badblob` | EIP-4844 blob tx with no blob hashes | ✅ |
| `malformed` | corrupted RLP bytes | ✅ |
| `truncated` | truncated RLP bytes | ✅ |
| `underpriced` | fee cap of 0 or 1 wei (can stick in mempool) | — (ok if base fee ~0) |
| `futurenonce` | nonce far ahead of the account nonce | — (queues normally) |
| `noncetoolow` | nonce below the account nonce | — (equals current on fresh acct) |
| `nofunds` | value far exceeding the wallet balance | — (balance is state-dependent) |

## Usage

```bash
spamoor tx-fuzz-invalid [flags]
```

### Volume Control (either -c or -t required)
- `-c, --count` - Total number of invalid transactions to send
- `-t, --throughput` - Invalid transactions to send per slot (default: 20)
- `--max-pending` - Maximum concurrent in-flight submissions (default: 50)

### Fuzzing Configuration
- `--categories` - Comma list of categories to fuzz (or `all`, default). See table above.
- `--payload-seed` - Custom hex seed for reproducible fuzzing (e.g. 0x1234abcd)

### Wallet / Client
- `--max-wallets` - Maximum number of burner wallets to use
- `--client-group` - Client group to use for sending transactions

### Debug
- `--log-txs` - Log every submission and its rejection reason

## Examples

Fire all invalid categories:
```bash
spamoor tx-fuzz-invalid -p "<PRIVKEY>" -h http://rpc-host:8545 -t 20
```

Only nonce and fee edge cases, logging each rejection:
```bash
spamoor tx-fuzz-invalid -p "<PRIVKEY>" -h http://rpc-host:8545 -t 10 --categories noncetoolow,futurenonce,underpriced --log-txs
```

## Output

The scenario logs a periodic summary of `accepted/sent` per category. **Any
non-zero `accepted` count is a potential finding** and is also logged loudly at
warning level when it happens, including the client name and category.
