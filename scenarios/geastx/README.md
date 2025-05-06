# Geas Transaction Scenario (`geastx`)

Send **EVM transactions** that first _deploy_ a contract containing custom [geas](https://github.com/fjl/geas) byte-code and then repeatedly _call_ that contract.

The scenario gives fine-grained control over throughput, gas limits, fee caps and wallet usage. Either a geas **file** or an **inline code string** must be supplied.

## Usage

```bash
spamoor geastx [flags]
```

## Configuration

### Base settings (required)
- `--privkey`  Private key of the root funding wallet
- `--rpchost`  RPC endpoint(s) used to broadcast the transactions

### Contract code ( **either** `--geasfile` **or** `--geascode` is required)
- `--geasfile`  Path to a file containing geas opcodes / assembly
- `--geascode`  Inline geas opcodes / assembly string (overrides `--geasfile`)


### Volume control ( **either** `-c` **or** `-t` is required )
- `-c, --count`  Total number of geas transactions to send
- `-t, --throughput` Transactions to send _per slot_
- `--max-pending` Maximum number of in-flight (pending) transactions

### Transaction parameters
- `--amount`  ETH amount (in **gwei**) to attach to **each** geas call (default `0`)
- `--basefee`  Max fee per gas (gwei)   (default `20`)
- `--tipfee`  Max priority fee (gwei)  (default `2`)
- `--gaslimit`  Gas limit for the **call** transactions (default `1,000,000`)
- `--deploy-gaslimit` Gas limit for the **deployment** transaction (default `1,000,000`)
- `--rebroadcast` Seconds to wait before re-broadcasting an unconfirmed tx (default `120`)

### Wallet management
- `--max-wallets` Maximum number of child wallets used in parallel

### Client selection
- `--client-group` Only use RPC clients that belong to this group label

### Debug options
- `-v, --verbose` Verbose log output
- `--trace`     Full trace log output

## Examples

Deploy the geas contract from `burner.geas` and send **100** calls with 1 M gas each:

```bash
spamoor geastx -p "<PRIVKEY>" -h http://rpc:8545 \
  --geasfile burner.geas \
  -c 100 --gaslimit 1000000
```

Send **5** calls per slot using inline geas byte-code, tip 3 gwei and rebroadcast after 30 s:

```bash
spamoor geastx -p "<PRIVKEY>" -h http://rpc:8545 \
  --geascode "PUSH1 0x60 PUSH1 0x40 MSTORE ..." \
  -t 5 --tipfee 3 --rebroadcast 30
```
