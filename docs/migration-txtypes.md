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

## Frame transaction changes, October 2026

`txtypes` was aligned with the EIP texts as of 2026-10-18 (EIP-8141, 8250, 7906, 8272 at
EIPs commit `d4b89184`). Behaviour changes for library consumers:

**EIP-8250 keyed nonces.** The account nonce is selected by an *empty* `NonceKeys`; every
key, including zero, is an independent domain in `NONCE_MANAGER` and pays
`KeyedNonceAccessCost`. `UsesLegacyNonce`/`UsesAccountNonce` are true only for an empty
key set (or no keyed-nonce extension). `[0]`, `[0, 5]` and up to sixteen strictly
increasing keys are valid. A `NonceSeq` of `MaxNonceSeq` (inclusive bound) with a non-empty
key set is the single-use form (`UsesSingleUseNonce`, `WithSingleUseNonceKeys`); with an
empty key set it is invalid. JSON decoding keeps an empty `nonceKeys` empty. The nonce
calldata's floor cost is flat per byte, as the EIP prices it. `LegacyNonceKey` is gone.

The revision before 2026-10-16, in which `[0]` selected the account nonce, is selectable
with `FrameExtLegacyNonceKeyZero` for building and validating transactions for chains whose
clients still implement it. It is not on the wire: decoders never set it, and a decoded
`[0]` is key zero. `WithLegacyNonce` spells the account nonce for either revision.
`spamoor.KeyedNonceRevision` is what spamoor's own probing assumes.

**EIP-7906.** The mode is `ASSERT`: `FrameModeAssert`, `SpeciesAssert` (`"assert"`),
`AssertFrame`, `FrameTx.AssertIndex`, `FrameReceiptExtra.AssertReverted`. The `PostTx` names
remain as deprecated aliases. `ValidatePayload` now rejects the atomic batch flag on an
ASSERT frame and on the frame directly before one.

**EIP-8272.** `ValidateMempoolPrefix` counts the recent root verifier frame's execution
budget toward `MaxVerifyGas`. `RecentRootVerifyGas` is sized from the now-fixed
`RecentRootCode` (about 3.4k plus 2.5k per tuple). `ExpiryVerifierCode`, `NonceManagerCode`
and `RecentRootCode` are exported.

**Validation and signing.** Gas totals are summed with overflow checks: `ValidatePayload`
rejects limits that wrap, and `IntrinsicGas`, `CalldataFloorGas`, `ExecutionGas`, `MaxGas`
and `GetStateGas` saturate at `math.MaxUint64` instead of wrapping. `VerifySignatures`
verifies P256 entries (previously only their encoding was checked). `SignPayload` signs
explicit-msg entries before canonical-hash entries; `SignEntryP256` refuses to run after a
canonical-hash entry has been signed, since it changes that entry's digest. `SignTx` sets
`From` to the explicit sender for `ExplicitSenderTx` types, even when a co-signer's key
signs.

**Gas and cost accessors.** `Gas()`/`GetGas()` return `MaxGas`, the payer's charge basis
across both dimensions; block packing uses `ExecutionReservation()` and `StateReservation()`
separately. `Cost()` is `MaxGas × fee cap + Σ value`, which spans two accounts for frame
transactions: `FrameTx.MaxCost(blobBaseFee)` is the payer's up-front charge and
`FrameTx.SenderValue()` the sender's own outlay. `EncodeFrameReceiptRLP` and
`DecodeFrameReceiptRLP` handle the consensus receipt encoding.
