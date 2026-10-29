# Code Adoption (EIP-8298, EIP-8151, EIP-4758)

Exercise the bogota-devnet-2 EIPs that change what an account's code and keys mean, and check the
state each case leaves behind:

- [EIP-8298](https://eips.ethereum.org/EIPS/eip-8298) `SETCODEFROM` (`0xa7`): an account adopts the
  code of an existing contract by code hash.
- [EIP-8151](https://eips.ethereum.org/EIPS/eip-8151): `ecRecover` returns nothing for an account
  whose code is neither empty nor an EIP-7702 delegation.
- [EIP-4758](https://eips.ethereum.org/EIPS/eip-4758) `SENDALL` (`0xff`): moves the balance and
  leaves code, storage and nonce in place.

Together, the first two retire an EOA's key: once it adopts code, its ECDSA transactions are
invalid (EIP-3607) and `ecRecover` no longer vouches for it, while frame transactions its code
approves keep working.

## Usage

```bash
spamoor codeadopt [flags]
```

## What it does

On start it deploys four fixtures from a well-known wallet: a source returning `0x42`, an approver
source whose code reads an approval scope from calldata and `APPROVE`s it, an adopter that runs
`SETCODEFROM(calldata[0:32])` and returns the result, and a contract that calls the adopter through
`STATICCALL`. It then waits until `SETCODEFROM` succeeds in an `eth_call`, so it can be started
before the fork, and checks the source rules once:

| Source | Expected |
|---|---|
| contract with code | `1` |
| the adopter itself (same code hash) | `1` |
| empty account, EOA, precompile | `0` |
| any source, under `STATICCALL` | the call fails |

Each case then runs as one unit of throughput, reading state back from the node that confirmed its
transactions:

| Case | Transactions | Checks |
|---|---|---|
| `migrate` | 7702 delegation of a fresh key to the adopter; a call that makes it adopt the approver; a frame transaction from it | `ecRecover` returns the key with no code and with a delegation; a delegated account is not a source; the adopted code is installed; `ecRecover` then returns nothing; an ECDSA transaction from the key is refused; a frame transaction approved by the adopted code lands |
| `adopt-runtime` | deploy an adopter, call it with the source | code hash matches the source; calls run the adopted code |
| `adopt-initcode` | creation whose initcode adopts the source and returns `0xEF` | the adopted code is installed and the return data ignored |
| `adopt-same-tx-source` | creation that `CREATE`s a child and adopts it, then adopts the source | the child is refused as a source; the source is accepted |
| `sendall-runtime` | deploy a funded contract with storage, call its `SENDALL` | balance moved; code, storage and nonce kept |
| `sendall-initcode` | creation that writes storage and `SENDALL`s from initcode | balance moved; the account survives with nonce 1 and its storage, where EIP-6780 would have deleted it |

A failed check is logged as a warning with what was read, and a summary of passed and failed checks
is logged every 50 cases and at the end.

## Configuration

### Base Settings (required)
- `--privkey` - Private key of the sending wallet
- `--rpchost` - RPC endpoint(s) to send transactions to

### Volume Control (either -c or -t required)
- `-c, --count` - Total number of cases to run
- `-t, --throughput` - Cases to start per slot (default: 2)
- `--max-pending` - Maximum number of cases in flight

### Case Settings
- `--cases` - Comma-separated cases to run, or `all` (default)
- `--fund-amount` - Gwei given to each migrated account to pay for its frame transaction (default: 10000000)
- `--sendall-value` - Gwei each `SENDALL` case moves (default: 1000)

### Transaction Settings
- `--basefee` / `--basefee-wei` - Max fee per gas
- `--tipfee` / `--tipfee-wei` - Max tip per gas
- `--client-group` - Client group to send through
- `--log-txs` - Log every case result
- `--timeout` - Stop after this duration
