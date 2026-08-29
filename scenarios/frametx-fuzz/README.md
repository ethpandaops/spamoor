# Frame Transaction Fuzzer (`frametx-fuzz`)

Fuzzes [EIP-8141](https://eips.ethereum.org/EIPS/eip-8141) frame transactions and everything
stacked on them — [EIP-8250](https://eips.ethereum.org/EIPS/eip-8250) keyed nonces,
[EIP-8272](https://eips.ethereum.org/EIPS/eip-8272) recent roots and
[EIP-7906](https://eips.ethereum.org/EIPS/eip-7906) POST_TX frames — by generating transactions
across every dimension they define and reporting what was reached.

This is the generated counterpart to `frametx`, which sends a fixed list of shapes and checks each
one against the statuses it expects. Where that scenario enumerates and asserts, this one draws from
a seeded generator and reports coverage.

## What it does, and what it deliberately does not

A frame transaction is not one call with parameters. It is a list of calls, a validation grammar in
front of that list, an envelope whose shape depends on which EIPs the chain runs, and a receipt
reporting a result per frame. The scenario's job is to **reach every combination of those**, not to
decide what each one should produce.

That distinction is deliberate. Whether a frame should have failed, what an instruction should have
returned, whether a shape ought to propagate — each is a reading of a specification that is still
moving, and a fuzzer that bakes in its author's reading turns a client disagreement into a false
alarm about the client that disagreed. On a network of more than one node a genuine disagreement
splits the chain by itself, which is a far stronger signal than any comparison made here.

So the run reports **coverage**: which dimensions each generated transaction exercised, how many
landed, and what the chain said about the ones it refused. A dimension with a zero count is the
actionable result — it means the generator never got there.

```
frametx-fuzz coverage: 20 generated, 10 confirmed, 10 refused, 3 invalid sent (0 accepted),
covered [arbitrary-witness=6 batch=10 contract-sender=4 expiry-frame=3 frame:probe=13
introspection-reads=6 keyed-nonce-first-use=10 keyed-nonces=11 p256-signature=4
prefix:only_verify+pay=7 recent-roots=8 root-edge:same_slot=2 starved-frame=6 …]
```

Rejection reasons are recorded per shape and compared across the run with numbers stripped, so a
chain that starts refusing the same shape for a *different* reason is visible without anyone having
decided which reason is right.

## Reproducibility

Every transaction is drawn from a per-index generator seeded by `--payload-seed`, exactly as
`evm-fuzz` does. A run with no seed generates one and logs it.

Draws happen in a single pass that produces an abstract **recipe** containing no chain state: which
wallet sends it, what nonce it gets and where the probe contract lives are applied afterwards. If a
draw consulted chain state, the same seed would produce different transactions on a chain in a
different state and the reproduction line would be a lie.

A finding therefore reports both:

```
--payload-seed 0x… --tx-id-offset 1234 -c 1     # regenerate the recipe on any chain
--recipe '{"index":1234,…}'                      # replay this exact description
```

## Dimensions

`--axes` takes a weighted list, e.g. `--axes nonces:5,roots:2`. The default, `all`, enables every
axis with equal weight. An axis whose EIP the chain does not run is disabled automatically.

| Axis | What it varies |
|---|---|
| `prefix` | The recognized validation prefixes, and the optional leading expiry frame |
| `batches` | Atomic batch runs, holding and unrolling |
| `failures` | Frames built to fail: a revert, or a budget too small for the frame-entry charge |
| `signatures` | ARBITRARY witnesses and P256 entries alongside the sender's own |
| `nonces` | EIP-8250 keyed nonce domains, including how many see their first use |
| `roots` | EIP-8272 references, including the window edges and the cases that must be refused |
| `posttx` | EIP-7906 assertion frames, passing and failing |
| `probe` | Calls into the fixed probe contract, including the introspection sweep |
| `code` | Deploying fuzzed contracts from one frame and calling them from another, and the fuzzed sender and paymaster roles |

### Generated contracts

Frames that only address wallets never reach code, so none of the instructions EIP-8141
introduces would ever execute. Two contracts do that work.

**A fixed probe contract**, deployed once through the CREATE2 factory, whose calldata is a
script: revert, write storage, emit a log, burn gas, `APPROVE`, and execute each
introspection instruction. The read operations discard their results — the point is to
reach the instruction, not to decide what it should have returned. One consequence worth
noting: a script that reads both `SIGDATACOPY` and `RECENTROOTREFLOAD` emits the byte
`0xb5` twice, because EIP-8141 and EIP-8272 both assign it.

The same code plays paymaster and contract sender through an EIP-7702 delegation, because
EIP-8141's own answer (a deploy frame) costs ~224k gas against the 100k execution cap on
the whole validation prefix. Since a delegated wallet is no longer an account without
code, the pool is split: the last quarter carries the delegation, the rest keep exercising
the protocol's default code. Sponsorship spreads across all of them, since the mempool
caps how many pending transactions one non-canonical paymaster may sponsor.

**Fuzzed contracts**, generated per transaction by the same stack-aware generator
`evm-fuzz` uses, with the frame instructions added to its table. A fixed contract can only
execute the sequences someone wrote into it; the interesting space is the instructions in
combination — with each other, with ordinary EVM, and at depth inside nested calls.

One frame deploys through the factory and a later frame calls the result. Because the
CREATE2 address is a function of the code, it is known before the transaction is sent, so
a frame can name a contract the transaction has not created yet. Frames also call
contracts earlier transactions deployed, which the run keeps in a bounded registry.

**Fuzzed account contracts**, for the sender and paymaster roles themselves. Body frames
deploy them — generated code ending in an `APPROVE` that reads its scope from calldata,
though one in ten is generated without one so the role can fail outright — and later
transactions use them. A small ring keeps only the most recent, so the code behind those
roles keeps turning over instead of being a fixed set installed once.

Deployment carries value in the same frame: the CREATE2 factory forwards `CALLVALUE`, and
a paymaster is charged the transaction's maximum cost up front, so it has to hold one. A
separate funding transfer would not do, because transferring to the contract *calls* it and
its trailing `APPROVE` reverts outside a validation frame.

This is the only thing that puts generated code inside the validation prefix, where
EIP-8141's banned-opcode and storage rules apply and a public mempool node has to simulate
what it finds. A contract sender also carries **no signature at all** — its code approves,
so there is nothing to sign — and its nonce, which starts at one for a created contract, is
tracked through a keyless wallet registered with the transaction pool.

Both are drawn rarely: a fuzzed prologue often halts before reaching the `APPROVE`, and
such a transaction never lands, so the fixed contract plays both roles most of the time.

> A *fresh* sender per transaction would be possible through EIP-8141's
> `[deploy, only_verify, pay]` prefix, where the deploy frame installs code at `tx.sender`
> because `tx.sender` is the CREATE2 address. It is not reachable on the client tested
> here: a deploy frame's code cost lands in the **execution** dimension (~1,530/byte, with
> `stateGasUsed` reported as zero), so even a small contract needs roughly 500k execution
> against a `MAX_VERIFY_GAS` of 100k. On a client charging it to state it would fit inside
> both caps.

## Invalid combinations

Malformed transactions are drawn from the same stream as well-formed ones rather than living in a
separate mode: the edge cases worth reaching are combinations, and a combination that happens to be
illegal is one of them. `--invalid-ratio` (default 0.05) sets how often a recipe carries a
deliberate violation — an unapproved verification frame, a batched prefix frame, a reserved mode, a
value outside SENDER mode, an out-of-range nonce key set, a truncated payload, a non-canonical
signature, and so on.

Those transactions are signed manually and fired with `eth_sendRawTransaction` from a burner wallet,
bypassing the managed pool: a transaction that never lands would otherwise stall its sender's nonce.
Structural violations are applied *before* signing so the signature covers them — applying one
afterwards would have every case refused for a bad signature instead of for the thing being
exercised.

What the chain did is recorded, not judged. An accepted violation is logged as an observation with
its reproduction line, because what a client must reject is part of what is still being settled.

The same applies to the awkward-but-legal cases the generator reaches on its own: recent root
references at the window edges, references to a slot that was never written, duplicate references,
and key sets whose sequences have to be read from chain state before they can be used.

## Usage

```bash
spamoor frametx-fuzz -p "<PRIVKEY>" -h http://rpc-host:8545 -t 10
```

### Volume control (either -c or -t required)
- `-c, --count` — Total number of transactions to generate
- `-t, --throughput` — Transactions per slot
- `--max-pending` — Maximum number of pending transactions

### Fuzzing
- `--payload-seed` — Hex seed for reproducible generation (empty generates and logs one)
- `--tx-id-offset` — Start generating from a specific transaction id
- `--recipe` — Replay a single recipe, as reported with a finding
- `--axes` — Weighted list of dimensions (default: `all`)
- `--max-frames` — Maximum body frames per transaction (default: 6)
- `--invalid-ratio` — Share of the stream carrying a deliberate violation (default: 0.05)
- `--log-frames` — Log the per-frame result of every landed transaction
- `--max-code-size` — Maximum size of a generated contract's runtime code (default: 256)
- `--code-gas` — Base execution gas for frames that deploy or call generated code

### Frame settings
- `--envelope` — Pin the payload shape: `auto` (default), `base`, `keyed`, `roots`, `full`
- `--post-tx` — EIP-7906 frames: `auto` (probe the chain), `on`, `off`
- `--user-op-gas`, `--verify-gas`, `--state-gas`, `--amount`, `--expiry-offset`, `--data`

### Wallets, fees and clients
- `--max-wallets`, `--basefee`, `--tipfee`, `--basefee-wei`, `--tipfee-wei`, `--rebroadcast`,
  `--client-group`, `--timeout`, `--log-txs`

## Throughput

The public mempool keeps **one pending frame transaction per sender**, so throughput is the wallet
count divided by confirmation latency rather than a function of pending depth. The scenario sizes
its pool at roughly twelve wallets per transaction per slot for that reason; a pool sized like an
ordinary scenario's starves the generator at a few transactions per block whatever `-t` says.

## Requirements

The chain must have an account at the EIP-8141 expiry verifier predeploy (`0x…8141`); the scenario
refuses to start otherwise rather than sending transactions every client will reject. The envelope
shape and the two extensions are read from their predeploys the same way, and `--pre-amsterdam-fee-model`
is incompatible: frame transactions need the EIP-8037 gas model.
