# Uniswap Swaps

Execute Uniswap V2 or V3 swaps with configurable parameters. This scenario allows you to test Uniswap DEX interactions by performing buy and sell transactions.

The scenario self-deploys everything it needs on a fresh network: two factories (each with its own router), a shared mock quote token, one mock DAI token per pair, and seeded liquidity. Every pair trades a DAI token against the quote token; both are plain ERC20s with a public `mint`, so **no ETH capital is needed beyond gas**: liquidity is minted into the pools by the liquidity-provider helper, and each child wallet mints its own quote tokens (see `--quote-funding`) both at startup and whenever it runs low. Select the version with `--uniswap-version` (default: 2). In v3 mode each DAI gets a pool on both factories at the configured `--fee-tier`, and swaps are routed through the canonical Uniswap v3 `SwapRouter`.

Every pool is seeded with 2000 quote tokens and 20,000,000 DAI, so the initial price is 10,000 DAI per quote token. Swap sizes (`--min-swap`, `--max-swap`, `--sell-threshold`) are denominated in DAI.

## Usage

```bash
spamoor uniswap-swaps [flags]
```

## Configuration

### Base Settings (required)
- `--privkey` - Private key of the sending wallet
- `--rpchost` - RPC endpoint(s) to send transactions to

### Volume Control (either -c or -t required)
- `-c, --count` - Total number of transactions to send
- `-t, --throughput` - Transactions to send per slot
- `--max-pending` - Maximum number of pending transactions

### Transaction Settings
- `--basefee` - Max fee per gas in gwei (default: 20)
- `--tipfee` - Max tip per gas in gwei (default: 2)
- `--gaslimit` - Gas limit for each transaction in gwei (default: 200000)
- `--rebroadcast` - Seconds to wait before rebroadcasting (default: 120)

### Swap Settings
- `--uniswap-version` - Uniswap version to use, 2 or 3 (default: 2)
- `--fee-tier` - Uniswap v3 fee tier in hundredths of a bip: 500, 3000 or 10000 (default: 3000)
- `--pair-count` - Number of uniswap pairs (DAI tokens) to deploy (default: 1)
- `--min-swap` - Minimum swap amount in DAI wei (default: 100000000000000000)
- `--max-swap` - Maximum swap amount in DAI wei (default: 1000000000000000000000)
- `--buy-ratio` - Ratio of buy vs sell swaps (0-100, default: 40)
- `--slippage` - Slippage tolerance in basis points (default: 50)
- `--slippage-min` - Min per-trade slippage in basis points; 0 disables and uses the fixed `--slippage` (default: 0)
- `--slippage-max` - Max per-trade slippage in basis points; when greater than `--slippage-min`, each trade draws a uniform-random tolerance in `[min, max]` (default: 0)
- `--sell-threshold` - DAI balance threshold to force sell in wei (default: 50000000000000000000000)
- `--quote-funding` - Quote token amount in wei minted to each child wallet at startup and whenever a wallet cannot afford a buy (default: 5000000000000000000, i.e. 5 quote tokens = 50,000 DAI at the seeded price)

### Wallet Management
- `--max-wallets` - Maximum number of child wallets to use
- `--refill-amount` - ETH amount to fund each child wallet (default: 5)
- `--refill-balance` - Minimum ETH balance before refilling (default: 2)
- `--refill-interval` - Seconds between balance checks (default: 300)

Child wallets only spend ETH on gas, so the refill settings can be kept low.

### Client Settings
- `--client-group` - Client group to use for sending transactions
- `--deploy-client-group` - Client group to use for deployment transactions (same as --client-group if empty)

### Debug Options
- `-v, --verbose` - Enable verbose output
- `--log-txs` - Log all submitted transactions
- `--trace` - Enable tracing output

## Example

Send 100 buy transactions:
```bash
spamoor uniswap-swaps -p "<PRIVKEY>" -h http://rpc-host:8545 -c 100 --buy-ratio 100
```

Send 2 sell transactions per slot:
```bash
spamoor uniswap-swaps -p "<PRIVKEY>" -h http://rpc-host:8545 -t 2 --buy-ratio 0
```

Send Uniswap v3 swaps at the 0.3% fee tier:
```bash
spamoor uniswap-swaps -p "<PRIVKEY>" -h http://rpc-host:8545 -t 2 --uniswap-version 3 --fee-tier 3000
```

## Troubleshooting

- **Wallets keep minting instead of swapping**: `--quote-funding` is too small relative to `--max-swap`. A buy of N DAI costs roughly N / 10,000 quote tokens, so make sure the funding covers many buys.
- **Swaps revert with insufficient output**: the pool price moved more than the slippage tolerance allows between quoting and execution. Raise `--slippage` or lower `--max-swap` so single trades have less price impact.
