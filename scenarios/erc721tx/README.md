# ERC721 Transactions

Deploy an ERC721 NFT contract and perform token transfers between accounts.

## Usage

```bash
spamoor erc721tx [flags]
```

## Configuration

### Base Settings (required)
- `--privkey` - Private key of the sending wallet
- `--rpchost` - RPC endpoint(s) to send transactions to

### Volume Control (either -c or -t required)
- `-c, --count` - Total number of transactions to send (default: 0)
- `-t, --throughput` - Transactions to send per slot (default: 200)
- `--max-pending` - Maximum number of pending transactions (default: 0)

### Transaction Settings
- `--basefee` - Max fee per gas in gwei (default: 20)
- `--tipfee` - Max tip per gas in gwei (default: 2)
- `--rebroadcast` - Enable reliable rebroadcast system (default: 1)
- `--timeout` - Timeout for the scenario (e.g. '1h', '30m', '5s') - empty means no timeout

### Token Settings
- `--max-index` - Maximum token index to mint (default: 0, unlimited)
- `--random-index` - Use random token index to mint
- `--random-target` - Use random destination addresses

### Wallet Management
- `--max-wallets` - Maximum number of child wallets to use (default: 0, auto-calculated)

### Client Settings
- `--client-group` - Client group to use for sending transactions
- `--deploy-client-group` - Client group to use for deployment transactions (same as --client-group if empty)

### Debug Options
- `-v, --verbose` - Enable verbose output
- `--log-txs` - Log all submitted transactions
- `--trace` - Enable tracing output

## Example

Deploy ERC721 and send 1000 NFT transfers:
```bash
spamoor erc721tx -p "<PRIVKEY>" -h http://rpc-host:8545 -c 1000
```

Send 5 NFT transfers per slot continuously:
```bash
spamoor erc721tx -p "<PRIVKEY>" -h http://rpc-host:8545 -t 5
```

Mint NFTs with random token IDs up to index 1000:
```bash
spamoor erc721tx -p "<PRIVKEY>" -h http://rpc-host:8545 -t 10 --random-index --max-index 1000
```
