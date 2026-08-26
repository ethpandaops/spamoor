# Migrating to `txtypes`

spamoor's engine no longer represents transactions, receipts and blocks with
go-ethereum's `core/types`. It uses the `github.com/ethpandaops/spamoor/txtypes`
package instead, which implements its own encoding, signing and JSON-RPC decoding.

This is a breaking change for code that uses spamoor as a library. For most consumers
it is an import swap.

## Why

`types.TxData` declares unexported methods, so no package outside go-ethereum's
`core/types` can implement it. That made spamoor unable to send any transaction type
go-ethereum had not implemented yet — a problem for a tool whose job is to test client
implementations, since go-ethereum is often not the first client to ship an EIP.

Alongside that, block scanning kept a hardcoded list of known transaction types and
skipped everything else, and `types.Header` only exposes the fields go-ethereum's
current fork knows.

## What changed

| Before | After |
|---|---|
| `types.Transaction` | `txtypes.Transaction` |
| `types.Receipt` | `txtypes.Receipt` |
| `types.Block`, `types.Header` | `txtypes.Block`, `txtypes.Header` |
| `types.LegacyTx`, `types.AccessListTx`, `types.DynamicFeeTx`, `types.BlobTx`, `types.SetCodeTx` | same names under `txtypes` |
| `types.BlobTxSidecar` | `txtypes.BlobSidecar` |
| `types.Log`, `types.AccessList`, `types.AccessTuple`, `types.SetCodeAuthorization` | aliases in `txtypes`, unchanged |
| `types.SignSetCode(key, auth)` | `txtypes.SignAuthorization(auth, key)` |
| `types.LegacyTxType`, `types.ReceiptStatusSuccessful`, ... | same names under `txtypes` |

Accessor and field names are unchanged, so code reading `tx.Hash()`, `tx.Nonce()`,
`receipt.Status`, `receipt.Logs` and so on needs no edits.

## Migrating

Most files need only:

```diff
- "github.com/ethereum/go-ethereum/core/types"
+ "github.com/ethpandaops/spamoor/txtypes"
```

plus renaming `types.X` to `txtypes.X` for the types in the table above.

Two cases need attention:

**abigen callbacks keep go-ethereum's type.** `BuildBoundTx` and
`BuildBoundTxWithEstimate` take a callback that returns whatever the generated binding
builds, which is a `*types.Transaction`. Leave that signature alone; the returned
transaction is converted for you.

```go
tx, err := wallet.BuildBoundTx(ctx, txData, func(opts *bind.TransactOpts) (*types.Transaction, error) {
    return contract.Transfer(opts, toAddr, amount)
})
// tx is a *txtypes.Transaction
```

**Passing values into go-ethereum APIs** requires an explicit conversion:

```go
gethTx, err := tx.ToGethTx()          // errors for types go-ethereum cannot represent
gethReceipt := receipt.ToGethReceipt()
tx, err := txtypes.FromGethTx(gethTx)
```

## Compatibility shims

`spamoor/compat.go` keeps go-ethereum-typed entry points for consumers that cannot
migrate immediately. They are deprecated and convert internally:

- `TxPool.SendGethTransaction(ctx, wallet, *types.Transaction, opts)`
- `Wallet.BuildGethDynamicFeeTx`, `BuildGethLegacyTx`, `BuildGethBlobTx`, `BuildGethSetCodeTx`
- `SendTransactionOptions.OnGethConfirm`, `OnGethComplete`, `OnGethEncode`

The `OnGeth*` callbacks are skipped for transaction types go-ethereum cannot represent,
so scenarios that must handle every type should use the `txtypes` callbacks.

## Adding a transaction type

Implement `txtypes.TxData` and register it:

```go
func init() {
    txtypes.RegisterTxType(0x06, func() txtypes.TxData { return &FrameTx{} })
}
```

Optional capability interfaces cover the rest: `ECDSASignedTx` or `ExplicitSenderTx`
for signing, `NetworkEncodedTx` when the wire encoding differs from the canonical one,
`BlobTxData`, `AccessListTxData`, `AuthListTxData`, `StateGasTxData`, and `JSONTxData`
for decoding from JSON-RPC. Receipts with type-specific content register a decoder with
`txtypes.RegisterReceiptDecoder`.

Types with no registered decoder still appear in block processing as `txtypes.UnknownTx`
with the generic fields the node reports, so they take part in confirmation tracking and
block accounting.
