# Safe Multisig

Generate Safe (formerly Gnosis Safe) multisig load by creating multisigs with varying owner counts and driving `execTransaction` calls through them. This scenario deploys the canonical Safe v1.4.1 contracts, creates a mix of multisigs (different owner counts and thresholds) owned by the child wallets, and then continuously executes both EOA calls and contract calls through those safes, signing each `SafeTx` off-chain with the required number of owner keys.

## Usage

```bash
spamoor safe-multisig [flags]
```

## Configuration

### Base Settings (required)
- `--privkey` - Private key of the sending wallet
- `--rpchost` - RPC endpoint(s) to send transactions to

### Volume Control (either -c or -t required)
- `-c, --count` - Total number of execTransaction transactions to send
- `-t, --throughput` - execTransaction transactions to send per slot
- `--max-pending` - Maximum number of pending transactions

### Transaction Settings
- `--basefee` - Max fee per gas in gwei (default: 20)
- `--tipfee` - Max tip per gas in gwei (default: 2)
- `--basefee-wei` - Max fee per gas in wei (overrides --basefee for L2 sub-gwei fees)
- `--tipfee-wei` - Max tip per gas in wei (overrides --tipfee for L2 sub-gwei fees)
- `--rebroadcast` - Enable reliable rebroadcast system (default: 1)

### Multisig Settings
- `--min-owners` - Minimum number of owners per safe (default: 1)
- `--max-owners` - Maximum number of owners per safe, clamped to the wallet pool size (default: 5). The owner count is randomized per safe within `[min-owners, max-owners]` to produce a good mix.
- `--threshold` - Signing threshold per safe (default: 0 = n-of-n; otherwise clamped to the owner count)
- `--safes-per-wallet` - Number of safes to create per child wallet (default: 3)

### Action Settings
- `--contract-ratio` - Fraction of execTransactions that call the gas-burner contract; the rest are EOA value transfers (default: 0.5)
- `--burn-rounds` - Upper bound on hashing rounds for the gas-burner contract; the per-call gas is seed-derived up to this bound, producing a deterministic but varied amount of burned gas (default: 1000)
- `--eoa-value` - Value in gwei sent in EOA-call execTransactions (default: 0). Transfers target another safe so funds stay within the pool (or the root wallet when a safe is over-funded), and are only used when the source safe is funded (otherwise the iteration falls back to a contract call). When > 0, safes are kept funded with small incremental top-ups.
- `--funding-interval` - Interval, in slots, between safe balance top-up checks (default: 32).
- `--recreate-rate` - Probability (0..1) that an iteration re-creates/churns a safe instead of executing a transaction (default: 0). See *Recreate Mode* below; `1` makes it a creation-only, state-bloat workload.
- `--gas-limit` - Gas limit for execTransaction txs (default: 0 = auto-compute per transaction)
- `--timeout` - Timeout for the scenario (e.g. '1h', '30m', '5s') - empty means no timeout

### Wallet Management
- `--max-wallets` - Maximum number of child wallets to use
- If not specified, automatically calculated based on count/throughput (max 1000)

### Client Settings
- `--client-group` - Client group to use for sending transactions
- `--deploy-client-group` - Client group to use for deployments

### Debug Options
- `-v, --verbose` - Enable verbose output
- `--log-txs` - Log all submitted transactions
- `--trace` - Enable tracing output

## How It Works

1. **Contract Deployment**: Automatically deploys the canonical Safe v1.4.1 singleton (master copy), the `SafeProxyFactory`, and a `GasBurner` call target via the CREATE2 deployment factory (reused on restart).
2. **Safe Creation**: Each child wallet gets `--safes-per-wallet` multisigs. For each safe, the owner count, owner set (a random subset of the child wallets), threshold, and CREATE2 salt are all derived deterministically from the executor wallet and slot index, so a restarted scenario reproduces the same safes instead of creating duplicates. Safes are created with `SafeProxyFactory.createProxyWithNonce`; existing ones are detected by probing the counterfactual address and reused.
3. **Hashing Self-Check**: Once at startup, the off-chain EIP-712 `SafeTx` hashing (domain separator + struct hash) is validated against the deployed Safe's `domainSeparator()` and `getTransactionHash()`, failing fast if they ever diverge.
4. **Execution**: Each transaction picks one of the executor wallet's safes round-robin, builds a `SafeTx` (a gas-burner contract call or an EOA call, per `--contract-ratio`), signs it off-chain with the lowest-addressed `threshold` owners (sorted ascending as Safe requires), and submits `execTransaction` from the executor wallet. Because each safe is only ever exercised by its owning executor, its internal nonce stays ordered with the executor's EOA nonce.
5. **Nonce Self-Healing**: If an `execTransaction` reverts, the on-chain Safe nonce does not advance; the safe is flagged to re-read its nonce from chain before the next use, so a single revert does not cascade.

## Recreate Mode

When `--recreate-rate > 0`, only **one safe per child wallet** is created upfront (and funded); the rest are grown and churned lazily in the transaction loop. Each iteration for a wallet:

- With probability `--recreate-rate`:
  - if the wallet has fewer than `--safes-per-wallet` safes, **create** one (grow);
  - if it is at the limit, **tear down the oldest** safe:
    - empty safe → drop it and create the replacement directly (no teardown tx);
    - funded safe → forward its balance to the lowest-balance sibling via one `execTransaction` (funds stay in the pool); the next iteration creates the replacement. If the safe is **over-funded** (> 10× the funding target), only a small balance goes to the sibling and the remainder drains back to the **root wallet** on a follow-up iteration.
- Otherwise execute a normal transaction through an existing safe.

At `--recreate-rate 1` this is a continuous safe-creation/churn workload (each abandoned safe stays on-chain), useful for state-bloat testing. (When `--recreate-rate 0`, all safes are created upfront and every iteration executes through one.)

## Smart Contract Details

The scenario uses three contracts:

- **Safe**: The canonical Safe v1.4.1 singleton that holds the owner/threshold logic and the `execTransaction` entrypoint. Each multisig is a minimal proxy delegating to this singleton.
- **SafeProxyFactory**: Deploys Safe proxies via CREATE2 (`createProxyWithNonce`).
- **GasBurner**: A dummy call target whose `burn(seed, maxRounds)` consumes a deterministic, seed-derived amount of gas (`1 + keccak(seed) % maxRounds` hashing rounds) and emits a `GasBurned` event.

## Example Usage

Deploy and send 1000 execTransactions with a mix of 1-to-5-owner safes:
```bash
spamoor safe-multisig -p "<PRIVKEY>" -h http://rpc-host:8545 -c 1000
```

Send 20 per slot, all contract calls, 3-of-5 multisigs:
```bash
spamoor safe-multisig -p "<PRIVKEY>" -h http://rpc-host:8545 -t 20 \
  --min-owners 5 --max-owners 5 --threshold 3 --contract-ratio 1
```

Drive ETH transfers between safes (auto-funded with low incremental top-ups):
```bash
spamoor safe-multisig -p "<PRIVKEY>" -h http://rpc-host:8545 -t 10 \
  --contract-ratio 0 --eoa-value 1000000
```

Continuously create safes for state-bloat testing:
```bash
spamoor safe-multisig -p "<PRIVKEY>" -h http://rpc-host:8545 -t 20 --recreate-rate 1
```

Mostly execute, occasionally churn safes (grow/teardown):
```bash
spamoor safe-multisig -p "<PRIVKEY>" -h http://rpc-host:8545 -t 10 --recreate-rate 0.2
```

## Important Notes

- **Variable Owner Counts**: Owner counts are randomized per safe within `[--min-owners, --max-owners]`, giving a realistic mix of 1-of-1 up to n-of-n multisigs. With the default `--threshold 0`, every safe is n-of-n, so the number of signatures verified per transaction varies with the owner count.
- **Off-chain Signing**: Each `SafeTx` is signed with the owners' child-wallet keys (standard ECDSA, `v` in {27,28}, owners ordered ascending), exactly as `Safe.checkNSignatures` expects.
- **Funding**: By default EOA calls move zero value, so safes need no funds. When `--eoa-value` is set, safes are funded once synchronously after deployment and then topped up incrementally in the background via the wallet pool's funding path (using the batcher contract when enabled) - kept intentionally low rather than pre-funded with a large lump sum.
- **Gas Limits**: execTransaction gas is auto-sized per transaction from the signing threshold and (for gas-burner calls) the actual round count, with headroom for the Amsterdam state-access fee schedule. Override with `--gas-limit` for a fixed value.
- **Rebroadcast**: Enabled by default to ensure transaction inclusion in congested networks.
