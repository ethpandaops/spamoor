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
`spamoor` is a tool for sending mass blob transactions.

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

The tool provides multiple scenarios, that focus on different aspects of blob transactions. One of the scenarios must be selected to run the tool:

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

