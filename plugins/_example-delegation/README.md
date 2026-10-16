# Example plugin: delegated counter

Reference for the **plugin** layer of [Contributing Workloads](../../docs/contributing-workloads.md).

The `delegated-counter` scenario:
1. deploys a counter contract (`push0 sload push1 1 add push0 sstore stop`),
2. delegates an unfunded EOA (the `authority` well-known wallet) to it with an EIP-7702 set-code tx,
3. sends `--count`/`--throughput` txs from child wallets to the authority, each incrementing slot 0 of the authority's storage.

It is a plugin rather than a spammer config because step 2 depends on the receipt of step 1, and the target of step 3 is the delegated EOA rather than the deployed contract. No generic scenario chains those steps.

## Run

```bash
spamoor --plugin ./plugins/_example-delegation delegated-counter \
  -h http://localhost:8545 -p <PRIVKEY> -c 40 -t 20
```

Afterwards, slot 0 of the authority address (logged at startup) equals the number of confirmed txs:

```bash
cast storage <authority> 0 --rpc-url http://localhost:8545
```

## Ship it

Validate, package, and host the archive anywhere reachable over HTTP (e.g. a GitHub release of your repo):

```bash
spamoor-utils validate-plugin ./plugins/_example-delegation
make plugins   # builds plugins/_example-delegation.tar.gz
spamoor --plugin https://example.com/example-delegation.tar.gz delegated-counter ...
```
