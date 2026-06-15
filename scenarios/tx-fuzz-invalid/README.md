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
- Submissions come from a **small, reused pool of burner wallets**, not one
  funded wallet per transaction. Per-tx funding would flood the shared root
  wallet with funding transactions and break in combination with other
  scenarios. The burners barely spend (most invalid txs are rejected before
  execution and never touch balance), so they are funded minimally and rarely.

### Unstuck helper

Most invalid txs are rejected at submission and consume no on-chain nonce, so a
burner can be reused indefinitely. The exceptions are txs a node might accept
into its mempool at the current nonce (e.g. underpriced). When that happens the
wallet is flagged and **unstuck before reuse**: mirroring tx-fuzz's `tryUnstuck`,
it replaces the stuck nonces with fee-bumped self-transfers until the pending
nonce catches up to the confirmed nonce. Unstuck *is* the recycle step — burners
are never abandoned unless unstuck times out.

## Invalid categories

| Category | What's wrong |
|---|---|
| `chainid` | signed for a different chain id |
| `lowgas` | gas limit below the intrinsic minimum |
| `underpriced` | fee cap of 0 or 1 wei (can stick in mempool) |
| `futurenonce` | nonce far ahead of the account nonce |
| `noncetoolow` | nonce below the account nonce |
| `nofunds` | value far exceeding the wallet balance |
| `gasoverflow` | gas limit far above the block gas limit |
| `emptyauth` | EIP-7702 set-code tx with an empty authorization list |
| `badblob` | EIP-4844 blob tx with no blob hashes |
| `malformed` | corrupted RLP bytes |
| `truncated` | truncated RLP bytes |

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
