# Seaport Trades

Generate Seaport (OpenSea) NFT marketplace order fulfillments with configurable parameters. This scenario stress-tests the dominant on-chain NFT-trading interaction by performing buy and sell `fulfillOrder` transactions against a self-deployed marketplace.

The scenario self-deploys everything it needs on a fresh network: the canonical **Seaport 1.6** marketplace and its **ConduitController**, plus a mock NFT collection and a mock ERC20 stablecoin (both with permissionless mint). A single **market** wallet acts as the standing counterparty - it signs every listing and bid off-chain (EIP-712), while the child wallets are the on-chain fulfillers (takers) that pay the gas.

Every transaction is a real Seaport `fulfillOrder`:
- **Buy** - the market lists one of its NFTs for the stablecoin; the taker pays the price and receives the NFT.
- **Sell** - the market bids the stablecoin for one of the taker's NFTs; the taker hands over the NFT and receives the price.

The buy/sell mix is controlled by `--buy-ratio` (like the uniswap-swaps scenario) and nudged by each wallet's current inventory: a wallet with no NFTs buys, and one holding more than `--sell-threshold` NFTs is pushed to sell, so NFTs and coins keep recirculating. NFTs and stablecoins are minted up front and **self-minted on demand** when the market's inventory or coin float runs low, so a run never stalls for lack of something to trade.

Transfers use a zero conduit key (they route through Seaport directly), so the only setup per participant is an NFT operator approval and an ERC20 allowance to the Seaport contract, both granted during seeding.

## Usage

```bash
spamoor seaport-trades [flags]
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
- `--basefee-wei` - Max fee per gas in wei (overrides --basefee for L2 sub-gwei fees)
- `--tipfee-wei` - Max tip per gas in wei (overrides --tipfee for L2 sub-gwei fees)
- `--rebroadcast` - Enable reliable rebroadcast system (default: 1)

### Trade Settings
- `--buy-ratio` - Ratio of buy vs sell fulfillments (0-100, default: 50)
- `--min-price` - Minimum trade price in stablecoin wei (default: 10000000000000000)
- `--max-price` - Maximum trade price in stablecoin wei (default: 100000000000000000000)
- `--sell-threshold` - Force a sell when a wallet holds more than this many NFTs (default: 20)

### Inventory Settings (self-seeding)
- `--market-inventory` - NFTs minted to the market counterparty at start (default: 50)
- `--wallet-inventory` - NFTs minted to each trader wallet at start (default: 5)
- `--replenish-threshold` - Self-mint more market NFTs when its inventory drops below this (default: 10)
- `--replenish-batch` - Number of NFTs to self-mint per replenish (default: 50)

### Wallet Management
- `--max-wallets` - Maximum number of child wallets to use

### Client Settings
- `--client-group` - Client group to use for sending transactions
- `--deploy-client-group` - Client group to use for deployment transactions (same as --client-group if empty)

### Debug Options
- `-v, --verbose` - Enable verbose output
- `--log-txs` - Log all submitted transactions
- `--trace` - Enable tracing output

## Example

Send 100 fulfillments with an even buy/sell mix:
```bash
spamoor seaport-trades -p "<PRIVKEY>" -h http://rpc-host:8545 -c 100
```

Send 5 buy fulfillments per slot:
```bash
spamoor seaport-trades -p "<PRIVKEY>" -h http://rpc-host:8545 -t 5 --buy-ratio 100
```

Send sells only, with a larger starting inventory per wallet:
```bash
spamoor seaport-trades -p "<PRIVKEY>" -h http://rpc-host:8545 -t 2 --buy-ratio 0 --wallet-inventory 20
```
