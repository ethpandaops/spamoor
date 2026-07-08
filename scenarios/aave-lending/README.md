# Aave Lending

Deploy a complete **Aave V3** lending market and have child wallets repeatedly supply, borrow, repay and withdraw two reserve tokens. This stress-tests the full Aave interaction surface — collateral accounting, the oracle-backed health-factor checks, variable-debt index accrual and the aToken/debt-token mint/burn paths — which together make lending one of Ethereum's heavier state- and compute-bound workloads.

The scenario deploys the **canonical Aave V3 contracts (real `@aave/core-v3@1.19.3` creation bytecode)**, so the on-chain behavior and gas profile match mainnet Aave rather than a reimplementation. On a fresh network it self-deploys everything it needs and wires the market end to end:

- the seven logic libraries (`SupplyLogic`, `BorrowLogic`, `LiquidationLogic`, `EModeLogic`, `BridgeLogic`, `FlashLoanLogic`, `PoolLogic`), linked into the `Pool` at deploy time, plus `ConfiguratorLogic` linked into the `PoolConfigurator`;
- the `PoolAddressesProvider`, `ACLManager`, `Pool` + `PoolConfigurator` (implementations behind provider-deployed proxies), `AaveProtocolDataProvider` and a shared `DefaultReserveInterestRateStrategy`;
- the `AaveOracle` backed by a minimal fixed-price aggregator per token (USD, 8 decimals);
- two 18-decimal mock ERC20 reserve tokens, registered as reserves via `initReserves`, enabled as collateral (80% LTV) and for variable borrowing, then seeded with borrowable liquidity by the deployer.

The published `Pool`/`PoolConfigurator`/`FlashLoanLogic` bytecode ships with unlinked `__$...$__` library placeholders; the scenario resolves them to the deployed library addresses at deploy time (the placeholder strings are stable for the pinned package version). Because the market is deployed across several dependent steps (libraries → provider/proxies → reserves → liquidity), deployment runs as a sequence of batches, each mined before the next.

The stateless, shareable contracts (the logic libraries and the mock tokens) deploy at deployer-independent (global) CREATE2 addresses so they are deployed once and reused, while the market itself (addresses provider, pool, configurator, reserves) is deployer-key specific — every deployer key gets its own isolated Aave state rather than contending on a shared market.

### Organic activity

Each child wallet is funded with both tokens and grants the pool an unlimited allowance. Wallets then act **organically** rather than in lockstep: on each turn a wallet inspects its current on-chain position and picks a feasible action — supply, borrow, repay or withdraw, on either reserve — weighted at random, with amounts varied within a configurable range and bounded by what it can actually do (token balance, borrowing power, debt, collateral). Because actions are gated by feasibility, the stream is varied but rarely reverts.

A configurable fraction of wallets (`--risky-ratio`) instead run aggressive near-maximum-LTV positions (high collateral, borrowed up to their headroom). These are the liquidation targets.

### Price moves & liquidations

Every `--price-tick-interval` transactions the scenario moves one reserve's oracle price: most ticks gently mean-revert toward the $1 base with small noise, but occasionally a downward **crash shock** drops the price deep into the volatility band (`--price-volatility`). The varying utilization from the organic activity makes the borrow/supply **interest rates** move; the crashes occasionally push the risky near-max positions underwater, and any wallet then opportunistically **liquidates** them via `liquidationCall`. Liquidations are rare by design (they need a deep enough dip); raise `--price-volatility` / lower `--price-tick-interval` to make them more frequent, or set `--liquidations=false` / `--price-tick-interval 0` to disable the price/liquidation dynamics entirely.

## Usage

```bash
spamoor aave-lending [flags]
```

## Configuration

### Base Settings (required)
- `--privkey` - Private key of the sending wallet
- `--rpchost` - RPC endpoint(s) to send transactions to

### Volume Control (either -c or -t required)
- `-c, --count` - Total number of lending transactions to send
- `-t, --throughput` - Transactions to send per slot
- `--max-pending` - Maximum number of pending transactions
- `--max-wallets` - Maximum number of child wallets to use

### Transaction Settings
- `--basefee` - Max fee per gas in gwei (default: 20)
- `--tipfee` - Max tip per gas in gwei (default: 2)
- `--basefee-wei` / `--tipfee-wei` - Fee/tip per gas in wei (overrides the gwei flags for L2 sub-gwei fees)
- `--rebroadcast` - Enable the reliable rebroadcast system (default: 1)
- `--gas-limit` - Gas limit per lending transaction (default: 0 = built-in default of 1,200,000)

### Lending Settings
- `--min-amount` - Minimum amount per lending action, in wei (default: 1e18)
- `--max-amount` - Maximum amount per lending action, in wei (default: 2000e18)
- `--seed-amount` - Initial borrowable liquidity supplied to each reserve by the deployer, in wei (default: 1,000,000e18)
- `--wallet-funding` - Amount of each token minted to every child wallet, in wei (default: 100,000e18)

### Market Dynamics
- `--risky-ratio` - 1 in N wallets runs a near-max-LTV position (liquidation targets); 0 disables risky positions (default: 6)
- `--liquidations` - Enable opportunistic liquidation of underwater positions (default: true)
- `--price-tick-interval` - Move an oracle price every N transactions; 0 disables price moves (default: 40)
- `--price-volatility` - Oracle price band in basis points around the $1 base, e.g. 2000 = ±20% (default: 2000)

### Client Settings
- `--client-group` - Client group to use for sending transactions
- `--deploy-client-group` - Client group to use for deployments

### Timeout & Debug
- `--timeout` - Stop the scenario after a duration (e.g. `1h`, `30m`, `5s`); empty means no timeout
- `--log-txs` - Log all submitted transactions

## Example

```bash
# Steady 10 lending txs/slot against a local devnet
spamoor aave-lending -p "<PRIVKEY>" -h http://localhost:8545 -t 10

# Burst test: 5000 lending txs across more wallets, larger positions
spamoor aave-lending -p "<PRIVKEY>" -h http://localhost:8545 \
  -c 5000 --max-wallets 200 --max-amount 5000000000000000000000

# Volatile market with frequent liquidations
spamoor aave-lending -p "<PRIVKEY>" -h http://localhost:8545 -t 8 \
  --risky-ratio 3 --price-tick-interval 8 --price-volatility 3000

# Pure supply/borrow load, no price moves or liquidations
spamoor aave-lending -p "<PRIVKEY>" -h http://localhost:8545 -t 10 \
  --price-tick-interval 0 --liquidations=false
```

## Notes

- The deployer well-known wallet funds the heavy contract deployment, seeds reserve liquidity and acts as the reserve treasury, so it is given a larger ETH refill than the child wallets.
- The market is deployment-idempotent: a re-run reuses the already-deployed contracts and skips wiring, reserve initialization and seeding when it detects them already in place.
- Contract bindings are generated from canonical artifacts via `contract/compile.sh` (which fetches the pinned `@aave/core-v3` artifacts). Re-run it only to regenerate the `*.go` bindings; the library link placeholders hardcoded in `deployment.go` track the pinned package version.
