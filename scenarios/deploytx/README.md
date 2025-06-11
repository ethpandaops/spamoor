# Contract Deployments

Deploy contracts with custom bytecode to the network.

## Usage

```bash
spamoor deploytx [flags]
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
- `--gaslimit` - Gas limit per deployment. Set to 0 to burn all available block gas. (default: 1000000)
- `--rebroadcast` - Seconds to wait before rebroadcasting (default: 120)

### Bytecode Configuration
- `--bytecodes` - Comma-separated list of hex bytecodes to deploy
- `--bytecodes-file` - File containing hex bytecodes to deploy

### Wallet Management
- `--max-wallets` - Maximum number of child wallets to use
- `--refill-amount` - ETH amount to fund each child wallet (default: 5)
- `--refill-balance` - Minimum ETH balance before refilling (default: 2)
- `--refill-interval` - Seconds between balance checks (default: 300)

### Client Settings
- `--client-group` - Client group to use for sending transactions

### Debug Options
- `-v, --verbose` - Enable verbose output
- `--trace` - Enable tracing output

## Example

Deploy 100 contracts using bytecode from a file:
```bash
spamoor deploytx -p "<PRIVKEY>" -h http://rpc-host:8545 -c 100 --bytecodes-file bytecodes.txt
```

Deploy 2 contracts per slot using specific bytecode:
```bash
spamoor deploytx -p "<PRIVKEY>" -h http://rpc-host:8545 -t 2 --bytecodes "0x1234,0x5678" 