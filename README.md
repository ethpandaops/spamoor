<img align="left" src="./.github/resources/goomy.png" width="75">
<h1>Spamoor the Transaction Spammer</h1>

spamoor is a simple tool that can be used to generate various types of random transactions for ethereum testnets.

Spamoor can be used for stress testing (flooding the network with thousands of transactions) or to have a continuous amount of transactions over long time for testing purposes.

## Build

You can use this tool via pre-build docker images: [ethpandaops/spamoor](https://hub.docker.com/r/ethpandaops/spamoor)

Or build it yourself:

```
git clone https://github.com/ethpandaops/spamoor.git
cd spamoor
make
./bin/spamoor
```

## Usage

### `spamoor`
`spamoor` is a tool for sending mass transactions.

```
Usage of spamoor:
Required:
  -p, --privkey string        The private key of the wallet to send funds from.
  
  -h, --rpchost string        The RPC host to send transactions to (multiple allowed).
      --rpchost-file string   File with a list of RPC hosts to send transactions to.
      
Optional:
  -s, --seed string           The child wallet seed.
  -v, --verbose               Run the tool with verbose output.
```

The tool provides multiple scenarios, that focus on different aspects of transactions. One of the scenarios must be selected to run the tool:

#### `spamoor eoatx`

The `eoatx` scenario sends normal dynamic fee transactions.

```
Usage of ./bin/spamoor eoatx:
      --amount uint            Transfer amount per transaction (in gwei) (default 20)
      --basefee uint           Max fee per gas to use in transfer transactions (in gwei) (default 20)
  -c, --count uint             Total number of transfer transactions to send
      --data string            Transaction call data to send
      --gaslimit uint          Gas limit to use in transactions (default 21000)
      --max-pending uint       Maximum number of pending transactions
      --max-wallets uint       Maximum number of child wallets to use
  -p, --privkey string         The private key of the wallet to send funds from.
      --random-amount          Use random amounts for transactions (with --amount as limit)
      --random-target          Use random to addresses for transactions
      --rebroadcast uint       Number of seconds to wait before re-broadcasting a transaction (default 120)
      --refill-amount uint     Amount of ETH to fund/refill each child wallet with. (default 5)
      --refill-balance uint    Min amount of ETH each child wallet should hold before refilling. (default 2)
      --refill-interval uint   Interval for child wallet rbalance check and refilling if needed (in sec). (default 300)
  -h, --rpchost stringArray    The RPC host to send transactions to.
      --rpchost-file string    File with a list of RPC hosts to send transactions to.
  -s, --seed string            The child wallet seed.
  -t, --throughput uint        Number of transfer transactions to send per slot
      --tipfee uint            Max tip per gas to use in transfer transactions (in gwei) (default 2)
      --trace                  Run the script with tracing output
  -v, --verbose                Run the script with verbose output
```

#### `spamoor erctx`

The `erctx` scenario deploys an ERC20 contract and performs token transfers.

```
Usage of ./bin/spamoor erctx:
      --amount uint            Transfer amount per transaction (in gwei) (default 20)
      --basefee uint           Max fee per gas to use in transfer transactions (in gwei) (default 20)
  -c, --count uint             Total number of transfer transactions to send
      --max-pending uint       Maximum number of pending transactions
      --max-wallets uint       Maximum number of child wallets to use
  -p, --privkey string         The private key of the wallet to send funds from.
      --random-amount          Use random amounts for transactions (with --amount as limit)
      --random-target          Use random to addresses for transactions
      --rebroadcast uint       Number of seconds to wait before re-broadcasting a transaction (default 120)
      --refill-amount uint     Amount of ETH to fund/refill each child wallet with. (default 5)
      --refill-balance uint    Min amount of ETH each child wallet should hold before refilling. (default 2)
      --refill-interval uint   Interval for child wallet rbalance check and refilling if needed (in sec). (default 300)
  -h, --rpchost stringArray    The RPC host to send transactions to.
      --rpchost-file string    File with a list of RPC hosts to send transactions to.
  -s, --seed string            The child wallet seed.
  -t, --throughput uint        Number of transfer transactions to send per slot
      --tipfee uint            Max tip per gas to use in transfer transactions (in gwei) (default 2)
      --trace                  Run the script with tracing output
  -v, --verbose                Run the script with verbose output
```

#### `spamoor deploytx`

The `deploytx` scenario sends contract deployment transactions.

```
Usage of ./bin/spamoor deploytx:
      --basefee uint            Max fee per gas to use in deployment transactions (in gwei) (default 20)
      --bytecodes string        Bytecodes to deploy (, separated list of hex bytecodes)
      --bytecodes-file string   File with bytecodes to deploy (list with hex bytecodes)
  -c, --count uint              Total number of deployment transactions to send
      --gaslimit uint           Gas limit to use in deployment transactions (in gwei) (default 1000000)
      --max-pending uint        Maximum number of pending transactions
      --max-wallets uint        Maximum number of child wallets to use
  -p, --privkey string          The private key of the wallet to send funds from.
      --rebroadcast uint        Number of seconds to wait before re-broadcasting a transaction (default 120)
      --refill-amount uint      Amount of ETH to fund/refill each child wallet with. (default 5)
      --refill-balance uint     Min amount of ETH each child wallet should hold before refilling. (default 2)
      --refill-interval uint    Interval for child wallet rbalance check and refilling if needed (in sec). (default 300)
  -h, --rpchost stringArray     The RPC host to send transactions to.
      --rpchost-file string     File with a list of RPC hosts to send transactions to.
  -s, --seed string             The child wallet seed.
  -t, --throughput uint         Number of deployment transactions to send per slot
      --tipfee uint             Max tip per gas to use in deployment transactions (in gwei) (default 2)
      --trace                   Run the script with tracing output
  -v, --verbose                 Run the script with verbose output
```

#### `spamoor deploy-destruct`

The `deploy-destruct` scenario deploys contracts that self-destruct.

```
Usage of ./bin/spamoor deploy-destruct:
      --amount uint            Transfer amount per transaction (in gwei) (default 20)
      --basefee uint           Max fee per gas to use in transfer transactions (in gwei) (default 20)
  -c, --count uint             Total number of transfer transactions to send
      --gaslimit uint          The gas limit for each deployment test tx (default 10000000)
      --max-pending uint       Maximum number of pending transactions
      --max-wallets uint       Maximum number of child wallets to use
  -p, --privkey string         The private key of the wallet to send funds from.
      --random-amount          Use random amounts for transactions (with --amount as limit)
      --rebroadcast uint       Number of seconds to wait before re-broadcasting a transaction (default 120)
      --refill-amount uint     Amount of ETH to fund/refill each child wallet with. (default 5)
      --refill-balance uint    Min amount of ETH each child wallet should hold before refilling. (default 2)
      --refill-interval uint   Interval for child wallet rbalance check and refilling if needed (in sec). (default 300)
  -h, --rpchost stringArray    The RPC host to send transactions to.
      --rpchost-file string    File with a list of RPC hosts to send transactions to.
  -s, --seed string            The child wallet seed.
  -t, --throughput uint        Number of transfer transactions to send per slot
      --tipfee uint            Max tip per gas to use in transfer transactions (in gwei) (default 2)
      --trace                  Run the script with tracing output
  -v, --verbose                Run the script with verbose output
```

#### `spamoor blobs`

The `blobs` scenario sends out normal blobs with random data only.\
No replacement or cancellation transactions are being send.

#### `spamoor blob-replacements`

The `blob-replacements` scenario sends out blobs and always tries to replace these blobs with replacement blob transactions a few seconds later, further replacement transactions are being sent until inclusion in a block.

#### `spamoor blob-conflicting`

The `blob-conflicting` scenario sends out blob transactions and conflicting normal transactions at the same time or with a small delay.

#### `spamoor blob-combined`

For general testing, there is a `blob-combined` scenario, which combines parts of all other blob scenarios in a randomized way.

```
Usage of spamoor blob-combined:
Required (at least one of):
  -c, --count uint            Total number of blob transactions to send
  -t, --throughput uint       Number of blob transactions to send per slot
  
Optional:
      --basefee uint          Max fee per gas to use in blob transactions (in gwei) (default 20)
      --blobfee uint          Max blob fee to use in blob transactions (in gwei) (default 20)
      --max-pending uint      Maximum number of pending transactions
      --max-replace uint      Maximum number of replacement transactions (default 4)
      --max-wallets uint      Maximum number of child wallets to use
      --rebroadcast uint      Number of seconds to wait before re-broadcasting a transaction (default 30)
      --replace uint          Number of seconds to wait before replace a transaction (default 30)
  -b, --sidecars uint         Maximum number of blob sidecars per blob transactions (default 3)
      --tipfee uint           Max tip per gas to use in blob transactions (in gwei) (default 2)
```

### `spamoor gasburnertx`

The `gasburnertx` scenario sends out transactions with a configurable amount of gas units. Note that the estimated gas units is not 100% accurate.

```
Usage of spamoor gasburnertx:
Required (at least one of):
  -c, --count uint            Total number of gasburner transactions to send
  -t, --throughput uint       Number of gasburner transactions to send per slot
  
Optional:
      --basefee uint             Max fee per gas to use in gasburner transactions (in gwei) (default 20)
      --gas-units-to-burn uint   The number of gas units for each tx to cost (default 2000000)
      --max-pending uint         Maximum number of pending transactions
      --max-wallets uint         Maximum number of child wallets to use
  -p, --privkey string           The private key of the wallet to send funds from.
      --rebroadcast uint         Number of seconds to wait before re-broadcasting a transaction (default 120)
  -h, --rpchost stringArray      The RPC host to send transactions to.
      --rpchost-file string      File with a list of RPC hosts to send transactions to.
  -s, --seed string              The child wallet seed.
      --tipfee uint              Max tip per gas to use in gasburner transactions (in gwei) (default 2)
      --trace                    Run the script with tracing output
  -v, --verbose                  Run the script with verbose output
```


##### Examples:

Continuous random blob spamming (~2-4 sidecars per block):
```
spamoor blob-combined -p "<PRIVKEY>" -h http://rpc-host1:8545 -b 2 -t 3 --max-pending 3
```

flood the network with 1000 blobs (+some replacements) via 2 rpc hosts:
```
spamoor blob-combined -p "<PRIVKEY>" -h http://rpc-host1:8545 -h http://rpc-host2:8545 -c 1000
```

#### `spamoor wallets`

The `wallets` scenario prepares & prints the list of child wallets that are used to send blob transactions from.\
It's more intended for debugging. The tool takes care of these wallets internally, so there is nothing to do with them ;)

## Spamoor Daemon

The daemon provides a web interface and API for managing multiple spammers. It allows you to create, monitor, and control spammers through a user interface or programmatically via HTTP endpoints.

### Usage
```bash
spamoor-daemon [flags]
```

### Flags
```
-d, --db string         The file to store the database in (default "spamoor.db")
    --debug             Run the tool in debug mode
-h, --rpchost strings   The RPC host to send transactions to
    --rpchost-file      File with a list of RPC hosts to send transactions to
-p, --privkey string    The private key of the wallet to send funds from
-P, --port int          The port to run the webui on (default 8080)
-v, --verbose           Run the tool with verbose output
    --trace             Run the tool with tracing output
```

### Web Interface
The web interface runs on `http://localhost:8080` by default and provides:
- Dashboard for managing spammers
- Real-time log streaming
- Configuration management
- Start/pause/delete functionality

### API
The daemon exposes a REST API for programmatic control.
See the API Documentation in the spamoor web interface for details.
