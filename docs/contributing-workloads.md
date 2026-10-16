# Contributing Workloads

Spamoor core is the transaction machinery (wallets, nonces, fees, clients, rate limiting, receipts) plus a set of **generic, programmable scenarios**. Benchmarks, attack patterns and EIP tests are **workloads** built on top of it. Contribute a workload at the highest layer that can express it, and change core only to add a primitive that a workload is missing.

| Layer | What you contribute | Lives in | Review bar |
|---|---|---|---|
| 1. Spammer config | YAML composing existing scenarios | `spammer-configs/*.yaml` (Spammer Library) | Runs on a devnet; headers filled in |
| 2. Plugin | Go scenario with custom logic | Your own repo, loaded with `--plugin <url>` | Yours; core does not review it |
| 3. Core primitive | Generic capability that layers 1–2 need | This repo | Proven by a layer 1/2 consumer, tested, documented |

## Choosing a layer

**Start with a spammer config.** Most workloads are a parameterisation of a generic scenario:

| Need | Scenario |
|---|---|
| Run custom bytecode in every tx | `geastx` (inline geas or raw bytecode) |
| Deploy a contract and call it with ABI-encoded args | `calltx` |
| Plain value transfers to a fixed or derived address | `eoatx` |
| Deploy many contracts at deterministic CREATE2 addresses | `factorydeploytx` |
| Multi-step deploy/call sequences per cycle | `taskrunner` |
| Replay execution-spec test fixtures | `replay-eest` |

A config file can hold several spammers, e.g. a `factorydeploytx` preparation spammer and an attack spammer that targets the deployed contracts. See [`max-contract-read.yaml`](../spammer-configs/max-contract-read.yaml) and [`create2-receiver-transfers.yaml`](../spammer-configs/create2-receiver-transfers.yaml).

**Write a plugin** when the workload needs logic that config cannot express: multi-phase setup that depends on receipts, EIP-7702 delegation flows, per-tx state tracking, custom tx types. Plugins are Go source interpreted at runtime, have access to everything exported in [`plugin/symbols`](../plugin/symbols), and are distributed as `tar.gz` archives that the CLI and daemon load from a URL. See the [Plugin System Guide](./plugin-system.md).

**Propose a core change** only when a workload cannot be built at layer 1 or 2 because a primitive is missing, for example:
- a placeholder, option or helper that generic scenarios or plugins should share
- a symbol that plugins need but `plugin/symbols` does not export
- a new transaction type or a fix to wallet, nonce, fee or client handling

Scenario-specific options on a generic scenario (target pools, salt strides, benchmark presets) are not primitives. If a feature only makes sense for one workload, it belongs in that workload's config or plugin.

## Examples

One reference contribution per layer.

### Spammer config: CREATE2 receiver transfers

[`spammer-configs/create2-receiver-transfers.yaml`](../spammer-configs/create2-receiver-transfers.yaml) sends ETH transfers to a set of contracts at CREATE2 addresses. Two generic scenarios, no Go code:
- `factorydeploytx` deploys 1000 contracts with distinct 1 KiB code through the shared factory.
- `eoatx` sends each tx to `{create2:{factory_address}:<initcodehash>:{random:1000}}`.

```bash
spamoor run spammer-configs/create2-receiver-transfers.yaml -h <RPC> -p <PRIVKEY> -s 0   # deploy receivers
spamoor run spammer-configs/create2-receiver-transfers.yaml -h <RPC> -p <PRIVKEY> -s 1   # transfers
```

In the daemon the same file appears in the Spammer Library after `make generate-spammer-index`.

### Plugin: delegated counter

[`plugins/_example-delegation`](../plugins/_example-delegation) deploys a counter contract, delegates an EOA to it with an EIP-7702 set-code tx, then sends txs to the EOA so each one increments its storage. The delegation needs the deployment receipt and the spam targets the delegated EOA, so no spammer config can express it.

```bash
spamoor --plugin ./plugins/_example-delegation delegated-counter -h <RPC> -p <PRIVKEY> -c 40 -t 20
```

The [plugin README](../plugins/_example-delegation/README.md) covers validation, packaging and loading from a URL.

### Core primitive: CREATE2 placeholders

The `{create2:...}` and `{factory_address}` placeholders ([#297](https://github.com/ethpandaops/spamoor/pull/297)) were added because the config above could not otherwise address a CREATE2 contract set. Against the checklist below:

| Checklist item | In the change |
|---|---|
| Consumer | `spammer-configs/create2-receiver-transfers.yaml` |
| Generic | One resolver (`utils/placeholders.go`, `scenario/placeholders.go`) replacing three copies, used by `eoatx`, `calltx`, `geastx`, `taskrunner` and plugins |
| Tests and symbols | `utils/placeholders_test.go`, `scenarios/taskrunner/placeholders_test.go`, regenerated `plugin/symbols` |
| Docs | Scenario READMEs, this guide, the Plugin System Guide |

## Placeholders

Generic scenarios resolve these placeholders in user-supplied values. Placeholders nest and resolve innermost-first.

| Placeholder | Value |
|---|---|
| `{txid}` | Transaction index |
| `{random}` | Random uint256 |
| `{random:N}` | Random integer in `[0, N)` |
| `{randomaddr}` | Random address |
| `{factory_address}` | Shared CREATE2 factory deployed by `factorydeploytx` (`well_known_factory: true`) |
| `{create2:<factory>:<initcodehash>:<salt>}` | CREATE2 address; `salt` is a decimal or `0x` uint256 |

Where they apply:

| Scenario | Field | Placeholders |
|---|---|---|
| `eoatx` | `to` (per tx) | all |
| `calltx` | `call_args` (per tx) | all |
| `calltx` | `call_data` | `{factory_address}` |
| `geastx` | `geas_code` / `geas_file` | `{factory_address}` |
| `taskrunner` | task data | all except `{factory_address}`, plus `{stepid}`, `{contract:name}`, `{contract:name:nonce}` |

Example: send to a random contract from a 1000-contract CREATE2 set.

```yaml
to: "{create2:{factory_address}:0xd0923e8cfd0a2bfc5da9c3ba0e9c938e9bc40774659544ab488bf241b918a07a:{random:1000}}"
```

Plugins use the same resolver: `scenario.TxPlaceholders(walletPool, txIdx).Resolve(s)`, extended with their own entries if needed.

## Core primitive checklist

A core PR must:
1. Ship with, or link to, the spammer config or plugin that needs it.
2. Be generic: usable by more than one scenario or workload, with no workload-specific options.
3. Include tests, and regenerate symbols (`make generate-symbols`) when the exported API changes.
4. Document the primitive where users look for it (scenario README, this guide, the Plugin System Guide).
5. Pass `go fmt ./...`, `go vet ./...` and `staticcheck ./...`.

## Worked example: account-access benchmarks

Porting the execution-specs `test_account_access` benchmarks maps onto the layers like this:

| Workload | Layer | How |
|---|---|---|
| ETH transfers to a CREATE2 contract set | Config | `factorydeploytx` + `eoatx` with `to: "{create2:...}"` |
| CALL/EXTCODEHASH loops over a CREATE2 contract set | Config | `factorydeploytx` + `calltx` (attack bytecode, `{factory_address}` in `call_args`) or `geastx` (`{factory_address}` in the geas code) |
| SLOAD/SSTORE on a storage-bloated EIP-7702 EOA | Plugin | Delegation and slot-range bookkeeping are custom logic |
| Deriving CREATE2 addresses in config | Core primitive | `{create2:...}` and `{factory_address}` placeholders, shared by all scenarios |
