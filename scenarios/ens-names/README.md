# ENS Names

Deploys the full **real ENS stack** (`@ensdomains/ens-contracts@1.7.0`) at
deterministic addresses and drives organic ENS usage from child wallets:
commit-reveal `.eth` registrations, renewals, resolver record updates, name
transfers, NameWrapper wrap/unwrap and short-lived churn registrations.

A background **wallet naming service** (opt-out) additionally registers a
forward + reverse resolving `.eth` name for every spamoor wallet on the host:
the root wallet (`spamoor-root.eth`), this scenario's wallets and — in daemon
mode — every other running spammer's wallets
(`spamoor-<scenario><spammerid>-<name/index>.eth`).

## Usage

```bash
spamoor ens-names [flags]
```

### Configuration

#### Base Settings (required)
- `--privkey` - Private key of the sending wallet
- `--rpchost` - RPC endpoint(s) to send transactions to

#### Volume Control (either -c or -t required)
- `-c, --count` - Total number of ENS transactions to send
- `-t, --throughput` - ENS transactions to send per slot
- `--max-pending` - Maximum number of pending transactions

#### Deployment Settings
- `--deployment-seed` - Seed for the deterministic stack addresses (empty = one shared stack per root key)
- `--min-commitment-age` - Controller min commitment age in seconds (default: 10)
- `--max-commitment-age` - Controller max commitment age in seconds (default: 86400)
- `--deploy-client-group` - Client group to use for deployments

Note: the commitment ages are ETHRegistrarController constructor args, so
changing them yields a different (still deterministic) set of stack addresses.

#### Name Engine Settings
- `--names-per-wallet` - Names each child wallet registers & maintains (default: 3)
- `--registration-duration` - Registration duration in seconds (default: 90 days, min 28 days)
- `--renewal-duration` - Renewal duration in seconds (default: 30 days)
- `--rotation-weight` - Weight of rotation registrations (default: 10): keep registering fresh names via commit-reveal, retiring the oldest from active management
- `--renew-weight` - Weight of renewals (default: 20)
- `--record-update-weight` - Weight of resolver record updates (default: 25)
- `--transfer-weight` - Weight of name transfers to sibling wallets (default: 10)
- `--abandon-weight` - Weight of name abandons (default: 5)
- `--reverse-weight` - Weight of default reverse record updates (default: 10)
- `--wrap-weight` - Weight of NameWrapper wrap/unwrap ops (default: 15)
- `--churn-weight` - Weight of short-lived churn registrations (default: 15)

Weights set to 0 disable an operation. Churn registrations remain the fallback
when no other operation is feasible.

#### Wallet Naming Service
- `--wallet-naming` - Scope: `all` (every spamoor wallet on the host, default), `pool` (own wallets only), `off`
- `--naming-per-slot` - Max wallets named per slot (default: 2)

#### Transaction Settings
- `--basefee` - Max fee per gas in gwei (default: 20)
- `--tipfee` - Max tip per gas in gwei (default: 2)
- `--basefee-wei` / `--tipfee-wei` - Wei-denominated overrides for L2s
- `--rebroadcast` - Enable reliable rebroadcast system (default: 1)
- `--client-group` - Client group to use for sending transactions

#### Wallet Settings
- `--max-wallets` - Maximum number of child wallets to use

#### Debug Options
- `-v, --verbose` - Enable verbose output
- `--log-txs` - Log all submitted transactions
- `--trace` - Enable tracing output
- `--timeout` - Timeout for the scenario (e.g. '1h', '30m')

### Example Usage

```bash
# Steady ENS activity: registrations, renewals, record churn
spamoor ens-names -p "<PRIVKEY>" -h http://rpc-host:8545 -t 5

# 500 ENS txs without the wallet naming service
spamoor ens-names -p "<PRIVKEY>" -h http://rpc-host:8545 -c 500 --wallet-naming off

# Registration-heavy run with an isolated stack
spamoor ens-names -p "<PRIVKEY>" -h http://rpc-host:8545 -t 10 \
  --deployment-seed my-test --names-per-wallet 10 --churn-weight 50
```

## How It Works

### Deployment & permission model

All ENS contracts derive ownership from `msg.sender` at construction time, so
they cannot be deployed directly through the shared CREATE2 factory (the
factory proxy would become the permanent, unusable owner). Instead:

1. A small **EnsExecutor** contract is deployed through the CREATE2 factory
   with the **root wallet** as admin (constructor arg). Its address is
   deterministic per (root key, deployment seed).
2. The scenario's `deployer` wallet is authorized as executor **operator**
   with a single root-wallet `setOperator()` tx (skipped when the stack is
   already deployed and wired). Each scenario instance uses its own operator
   wallet, so multiple ENS spammers can run concurrently on the same stack
   without nonce conflicts.
3. The executor CREATE2-deploys the ENS contracts (making itself their owner)
   and performs all owner-gated wiring via `execute()`: registry ownership of
   `eth`/`reverse`/`addr.reverse`, registrar controller grants and the default
   reverse resolver.

Deployed stack: ENSRegistry, BaseRegistrarImplementation, DummyOracle +
StablePriceOracle (official mainnet rent prices, $1600/ETH), ReverseRegistrar,
DefaultReverseRegistrar, StaticMetadataService + NameWrapper,
ETHRegistrarController, PublicResolver and a local **SpamRegistrarController**.
Every deploy and wiring step is checked against chain state first, so restarts
and concurrent instances reuse the existing stack.

The SpamRegistrarController is a permissionless auxiliary controller that
allows direct registrations without commit-reveal and without the 28-day
minimum duration - used for the short-lived churn names and the one-tx
`registerNamed()` naming service path. Free direct registration is exactly the
kind of traffic this testnet scenario exists to generate; commit-reveal names
are unaffected.

### Name engine

Each child wallet works toward `--names-per-wallet` names via the full
commit-reveal flow (`commit()` → wait `min-commitment-age` → `register()` with
resolver records, paying oracle rent). The first registration also sets the
wallet's default (chain-agnostic) reverse record via the controller's
`reverseRecord` bitmask. Once at target, wallets perform weighted-random
maintenance: renewals, `setAddr`/`setText` multicalls, ERC721 transfers to
sibling wallets (followed by `reclaim()` from the receiver), abandons
(zeroing the registry owner), default reverse re-points, NameWrapper
wrap/unwrap and short-lived churn registrations that visibly expire during
the run.

Registration never stops long-term: the rotation op (`--rotation-weight`)
keeps registering fresh names through the full commit-reveal flow at a low
chance, and each time a rotation name lands, the wallet's oldest name is
retired from active management (it stays registered on-chain but no longer
receives maintenance). The actively managed set is always the newest
`--names-per-wallet` names.

Labels embed a per-run random id (`sp<runid>-<wallet>-<seq>`), so runs and
concurrent instances never collide. Name state is kept in memory per run;
nothing is reconstructed from previous runs.

### Wallet naming service

Every slot, the service polls the host wallet registry (all running spammers'
wallet pools + root wallet in daemon mode; the own pool in CLI mode or `pool`
mode) and names unnamed wallets via
`SpamRegistrarController.registerNamed()` — one tx per wallet that registers
`spamoor-<scenario><spammerid>-<name/index>.eth` (`spamoor-root.eth` for the
root wallet), sets the forward `addr` record and the `addr.reverse` record.
All txs are sent from the scenario's `nameservice` wallet, so wallets of other
scenarios are never touched (no nonce interference). On-chain reverse-record
checks make restarts and concurrent naming services converge without
coordination.

The `addr.reverse` records are exclusively owned by the naming service; the
name engine's reverse operations use the (registry-detached) default reverse
registrar, so the two never fight.
