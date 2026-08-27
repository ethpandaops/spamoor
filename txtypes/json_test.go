package txtypes

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/holiman/uint256"
)

// marshalRPCTx renders a signed transaction the way a node's JSON-RPC does, so the
// decoder is exercised against the field shape it sees in production.
func marshalRPCTx(t *testing.T, tx *Transaction, from common.Address) json.RawMessage {
	t.Helper()

	fields := map[string]any{
		"type":  hexutil.Uint64(tx.Type()),
		"nonce": hexutil.Uint64(tx.Nonce()),
		"gas":   hexutil.Uint64(tx.Gas()),
		"value": (*hexutil.Big)(tx.Value()),
		"input": hexutil.Bytes(tx.Data()),
		"hash":  tx.Hash(),
		"from":  from,
		"to":    tx.To(),
	}

	v, r, s := tx.RawSignatureValues()
	fields["r"] = (*hexutil.Big)(r)
	fields["s"] = (*hexutil.Big)(s)

	if tx.Type() == LegacyTxType {
		fields["v"] = (*hexutil.Big)(v)
		fields["gasPrice"] = (*hexutil.Big)(tx.GasPrice())
	} else {
		fields["yParity"] = hexutil.Uint64(v.Uint64())
		fields["chainId"] = (*hexutil.Big)(tx.ChainId())
		fields["accessList"] = tx.AccessList()

		if tx.Type() == AccessListTxType {
			fields["gasPrice"] = (*hexutil.Big)(tx.GasPrice())
		} else {
			fields["maxFeePerGas"] = (*hexutil.Big)(tx.GasFeeCap())
			fields["maxPriorityFeePerGas"] = (*hexutil.Big)(tx.GasTipCap())
		}
	}

	if hashes := tx.BlobHashes(); hashes != nil {
		fields["blobVersionedHashes"] = hashes
		fields["maxFeePerBlobGas"] = (*hexutil.Big)(tx.BlobGasFeeCap())
	}

	if auths := tx.AuthList(); auths != nil {
		fields["authorizationList"] = auths
	}

	raw, err := json.Marshal(fields)
	if err != nil {
		t.Fatalf("failed marshalling transaction fields: %v", err)
	}

	return raw
}

// TestJSONTxRoundTrip checks that transactions reconstructed from JSON-RPC fields hash
// to the same value as the originals.
func TestJSONTxRoundTrip(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed generating key: %v", err)
	}

	var (
		sender  = crypto.PubkeyToAddress(key.PublicKey)
		chainID = big.NewInt(7088110746)
		rng     = rand.New(rand.NewSource(99))
	)

	for round := 0; round < 16; round++ {
		for _, pair := range generateTxPairs(t, rng, chainID) {
			t.Run(pair.name, func(t *testing.T) {
				signed, err := SignTx(NewTx(pair.ours), chainID, key)
				if err != nil {
					t.Fatalf("failed signing transaction: %v", err)
				}

				decoded, err := UnmarshalJSONTx(marshalRPCTx(t, signed, sender))
				if err != nil {
					t.Fatalf("failed decoding transaction: %v", err)
				}

				if decoded.Hash() != signed.Hash() {
					t.Fatalf("hash mismatch: %s != %s", decoded.Hash(), signed.Hash())
				}

				if _, ok := decoded.Inner().(*UnknownTx); ok {
					t.Fatalf("type 0x%02x decoded as UnknownTx", signed.Type())
				}

				// The adopted hash is cached, so recompute from the decoded content to
				// confirm the fields themselves were reconstructed correctly. Blob
				// sidecars are not part of the JSON-RPC representation.
				if rebuilt := NewTx(decoded.Inner()); rebuilt.Hash() != signed.Hash() {
					t.Fatalf("reconstructed content hashes to %s, want %s", rebuilt.Hash(), signed.Hash())
				}

				from, err := decoded.From(chainID)
				if err != nil {
					t.Fatalf("failed reading sender: %v", err)
				}

				if from != sender {
					t.Fatalf("wrong sender: %s != %s", from, sender)
				}
			})
		}
	}
}

// TestJSONUnknownTxType checks that an unregistered transaction type stays trackable
// instead of failing the whole block decode.
func TestJSONUnknownTxType(t *testing.T) {
	var (
		hash = common.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111")
		from = common.HexToAddress("0x2222222222222222222222222222222222222222")
		to   = common.HexToAddress("0x3333333333333333333333333333333333333333")
	)

	raw := fmt.Sprintf(`{
		"type": "0x63",
		"nonce": "0x2a",
		"gas": "0x5208",
		"maxFeePerGas": "0x3b9aca00",
		"maxPriorityFeePerGas": "0x77359400",
		"value": "0xde0b6b3a7640000",
		"input": "0xdeadbeef",
		"to": "%s",
		"from": "%s",
		"hash": "%s"
	}`, to.Hex(), from.Hex(), hash.Hex())

	tx, err := UnmarshalJSONTx(json.RawMessage(raw))
	if err != nil {
		t.Fatalf("failed decoding unknown transaction type: %v", err)
	}

	if _, ok := tx.Inner().(*UnknownTx); !ok {
		t.Fatalf("expected UnknownTx, got %T", tx.Inner())
	}

	if tx.Type() != 0x63 {
		t.Fatalf("wrong type: 0x%02x", tx.Type())
	}

	if tx.Hash() != hash {
		t.Fatalf("wrong hash: %s", tx.Hash())
	}

	sender, err := tx.From(big.NewInt(1))
	if err != nil {
		t.Fatalf("failed reading sender: %v", err)
	}

	if sender != from {
		t.Fatalf("wrong sender: %s", sender)
	}

	if tx.Nonce() != 42 || tx.Gas() != 21000 || *tx.To() != to {
		t.Fatal("generic fields were not preserved")
	}

	if tx.Value().Cmp(big.NewInt(1e18)) != 0 {
		t.Fatalf("wrong value: %s", tx.Value())
	}

	// An unknown type cannot be re-encoded; it must fail rather than emit garbage.
	if _, err := tx.MarshalBinary(); err == nil {
		t.Fatal("expected an error re-encoding an unknown transaction type")
	}
}

// TestBlockDecode checks block decoding, including that an unknown transaction type
// keeps the transaction list aligned with the block's receipts.
func TestBlockDecode(t *testing.T) {
	raw := `{
		"number": "0x10",
		"hash": "0xaaaa000000000000000000000000000000000000000000000000000000000001",
		"parentHash": "0xbbbb000000000000000000000000000000000000000000000000000000000002",
		"timestamp": "0x64",
		"gasLimit": "0x1c9c380",
		"gasUsed": "0xa410",
		"baseFeePerGas": "0x7",
		"blobGasUsed": "0x20000",
		"excessBlobGas": "0x40000",
		"stateGasUsed": "0x1e8480",
		"transactions": [
			{"type":"0x0","nonce":"0x1","gas":"0x5208","gasPrice":"0x7","value":"0x0","input":"0x","hash":"0xcccc000000000000000000000000000000000000000000000000000000000003","from":"0x1111111111111111111111111111111111111111","to":"0x2222222222222222222222222222222222222222","v":"0x1b","r":"0x1","s":"0x2"},
			{"type":"0x63","nonce":"0x2","gas":"0x5208","value":"0x0","input":"0x","hash":"0xdddd000000000000000000000000000000000000000000000000000000000004","from":"0x3333333333333333333333333333333333333333"}
		]
	}`

	var block Block
	if err := json.Unmarshal([]byte(raw), &block); err != nil {
		t.Fatalf("failed decoding block: %v", err)
	}

	if block.Number != 16 || block.Timestamp != 100 || block.GasLimit != 30_000_000 || block.GasUsed != 42_000 {
		t.Fatal("header fields were not decoded correctly")
	}

	if block.BaseFee.Cmp(big.NewInt(7)) != 0 {
		t.Fatalf("wrong base fee: %s", block.BaseFee)
	}

	if block.StateGasUsed != 2_000_000 {
		t.Fatalf("wrong state gas: %d", block.StateGasUsed)
	}

	if len(block.Transactions) != 2 {
		t.Fatalf("expected 2 transactions, got %d", len(block.Transactions))
	}

	if _, ok := block.Transactions[1].Inner().(*UnknownTx); !ok {
		t.Fatal("unknown transaction type was not preserved in the block")
	}
}

// TestHeaderMissingNumber checks that a malformed block is rejected rather than
// silently decoded as block zero.
func TestHeaderMissingNumber(t *testing.T) {
	var header Header
	if err := json.Unmarshal([]byte(`{"hash":"0x00"}`), &header); err == nil {
		t.Fatal("expected an error decoding a header without a number")
	}
}

// TestReceiptDecode checks receipt decoding and the type-specific decoder hook.
func TestReceiptDecode(t *testing.T) {
	raw := `{
		"type": "0x2",
		"status": "0x1",
		"cumulativeGasUsed": "0x1234",
		"gasUsed": "0x5208",
		"effectiveGasPrice": "0x9",
		"blobGasUsed": "0x20000",
		"blobGasPrice": "0x3",
		"transactionHash": "0xeeee000000000000000000000000000000000000000000000000000000000005",
		"blockHash": "0xffff000000000000000000000000000000000000000000000000000000000006",
		"blockNumber": "0x10",
		"transactionIndex": "0x2",
		"contractAddress": null,
		"logs": []
	}`

	var receipt Receipt
	if err := json.Unmarshal([]byte(raw), &receipt); err != nil {
		t.Fatalf("failed decoding receipt: %v", err)
	}

	if !receipt.Successful() {
		t.Fatal("receipt should be successful")
	}

	if receipt.GasUsed != 21000 || receipt.CumulativeGasUsed != 0x1234 || receipt.TransactionIndex != 2 {
		t.Fatal("receipt fields were not decoded correctly")
	}

	if receipt.BlobGasUsed != 131072 || receipt.BlobGasPrice.Cmp(big.NewInt(3)) != 0 {
		t.Fatal("blob fields were not decoded correctly")
	}

	if receipt.BlockNumber.Cmp(big.NewInt(16)) != 0 {
		t.Fatalf("wrong block number: %s", receipt.BlockNumber)
	}

	if err := json.Unmarshal([]byte(`{"status":"0x1"}`), &receipt); err == nil {
		t.Fatal("expected an error decoding a receipt without a transaction hash")
	}
}

// TestReceiptDecoderRegistry checks that a type-specific decoder is invoked and can
// fill in fields the generic shape does not carry.
func TestReceiptDecoderRegistry(t *testing.T) {
	const testTxType = 0x7e

	RegisterReceiptDecoder(testTxType, func(receipt *Receipt, raw json.RawMessage) error {
		var payload struct {
			Payer common.Address `json:"payer"`
		}

		if err := json.Unmarshal(raw, &payload); err != nil {
			return err
		}

		receipt.Status = ReceiptStatusSuccessful
		receipt.Extra = &testReceiptExtra{Payer: payload.Payer}

		return nil
	})

	raw := `{
		"type": "0x7e",
		"cumulativeGasUsed": "0x10",
		"gasUsed": "0x10",
		"transactionHash": "0xeeee000000000000000000000000000000000000000000000000000000000007",
		"payer": "0x4444444444444444444444444444444444444444",
		"logs": []
	}`

	var receipt Receipt
	if err := json.Unmarshal([]byte(raw), &receipt); err != nil {
		t.Fatalf("failed decoding receipt: %v", err)
	}

	extra, ok := receipt.Extra.(*testReceiptExtra)
	if !ok {
		t.Fatalf("expected testReceiptExtra, got %T", receipt.Extra)
	}

	if extra.Payer != common.HexToAddress("0x4444444444444444444444444444444444444444") {
		t.Fatalf("wrong payer: %s", extra.Payer)
	}

	if !receipt.Successful() {
		t.Fatal("decoder-derived status was not applied")
	}
}

type testReceiptExtra struct {
	Payer common.Address
}

func (e *testReceiptExtra) ReceiptTxType() byte { return 0x7e }

// TestRegisterTxTypeDuplicate checks that competing implementations of the same wire
// type are rejected.
func TestRegisterTxTypeDuplicate(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected a panic registering a duplicate transaction type")
		}
	}()

	RegisterTxType(DynamicFeeTxType, func() TxData { return &DynamicFeeTx{} })
}

// TestCustomTxTypeRegistration checks that a transaction type defined outside the
// built-in set encodes, decodes and signs through the generic paths.
func TestCustomTxTypeRegistration(t *testing.T) {
	const customType = 0x55

	RegisterTxType(customType, func() TxData { return &customTx{} })

	if !IsTxTypeSupported(customType) {
		t.Fatal("custom type was not registered")
	}

	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed generating key: %v", err)
	}

	chainID := big.NewInt(1337)

	signed, err := SignTx(NewTx(&customTx{Nonce: 7, Note: "hello"}), chainID, key)
	if err != nil {
		t.Fatalf("failed signing custom transaction: %v", err)
	}

	encoded, err := signed.MarshalNetwork()
	if err != nil {
		t.Fatalf("failed encoding custom transaction: %v", err)
	}

	if encoded[0] != customType {
		t.Fatalf("wrong type byte: 0x%02x", encoded[0])
	}

	decoded, err := DecodeTx(encoded)
	if err != nil {
		t.Fatalf("failed decoding custom transaction: %v", err)
	}

	if decoded.Hash() != signed.Hash() {
		t.Fatalf("hash mismatch: %s != %s", decoded.Hash(), signed.Hash())
	}

	sender, err := decoded.From(chainID)
	if err != nil {
		t.Fatalf("failed reading sender: %v", err)
	}

	if sender != crypto.PubkeyToAddress(key.PublicKey) {
		t.Fatalf("wrong sender: %s", sender)
	}

	// go-ethereum cannot represent this type, so conversion must fail rather than
	// produce a transaction with a different hash.
	if _, err := signed.ToGethTx(); err == nil {
		t.Fatal("expected conversion to go-ethereum to fail")
	}

	if _, err := signed.ToGethTx(); err != nil && !strings.Contains(err.Error(), "0x55") {
		t.Fatalf("error should name the transaction type: %v", err)
	}
}

// customTx is a minimal out-of-package transaction type used to verify that the
// registry, encoder and signer work without changes to this package.
type customTx struct {
	Sender common.Address
	Nonce  uint64
	Note   string
	Sig    []byte
}

var (
	_ TxData           = (*customTx)(nil)
	_ ExplicitSenderTx = (*customTx)(nil)
)

func (tx *customTx) TxType() byte { return 0x55 }

func (tx *customTx) CopyTx() TxData {
	cpy := *tx
	cpy.Sig = common.CopyBytes(tx.Sig)

	return &cpy
}

func (tx *customTx) GetChainID() *big.Int      { return new(big.Int) }
func (tx *customTx) GetNonce() uint64          { return tx.Nonce }
func (tx *customTx) GetGas() uint64            { return 0 }
func (tx *customTx) GetGasPrice() *big.Int     { return new(big.Int) }
func (tx *customTx) GetGasFeeCap() *big.Int    { return new(big.Int) }
func (tx *customTx) GetGasTipCap() *big.Int    { return new(big.Int) }
func (tx *customTx) GetTo() *common.Address    { return nil }
func (tx *customTx) GetValue() *big.Int        { return new(big.Int) }
func (tx *customTx) GetData() []byte           { return []byte(tx.Note) }
func (tx *customTx) GetSender() common.Address { return tx.Sender }

func (tx *customTx) EncodePayload(w *bytes.Buffer) error { return rlp.Encode(w, tx) }
func (tx *customTx) DecodePayload(b []byte) error        { return rlp.DecodeBytes(b, tx) }

func (tx *customTx) SignPayload(chainID *big.Int, key *ecdsa.PrivateKey) error {
	tx.Sender = crypto.PubkeyToAddress(key.PublicKey)

	digest := prefixedRlpHash(tx.TxType(), []any{chainID, tx.Sender, tx.Nonce, tx.Note})

	sig, err := crypto.Sign(digest[:], key)
	if err != nil {
		return err
	}

	tx.Sig = sig

	return nil
}

// TestSignAuthorization checks that an EIP-7702 authorization is signed by the key it
// is given, guarding the argument order against the go-ethereum helper it wraps.
func TestSignAuthorization(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed generating key: %v", err)
	}

	auth := SetCodeAuthorization{
		ChainID: *uint256.NewInt(1337),
		Address: common.HexToAddress("0x5555555555555555555555555555555555555555"),
		Nonce:   9,
	}

	signed, err := SignAuthorization(auth, key)
	if err != nil {
		t.Fatalf("failed signing authorization: %v", err)
	}

	if signed.ChainID != auth.ChainID || signed.Address != auth.Address || signed.Nonce != auth.Nonce {
		t.Fatal("signing altered the authorization content")
	}

	authority, err := signed.Authority()
	if err != nil {
		t.Fatalf("failed recovering authority: %v", err)
	}

	if authority != crypto.PubkeyToAddress(key.PublicKey) {
		t.Fatalf("authorization signed by the wrong key: %s", authority)
	}
}

// TestReceiptBloomAndPostState checks the two fields an explorer needs that are not
// part of the gas/status set: the log bloom and the pre-Byzantium state root.
func TestReceiptBloomAndPostState(t *testing.T) {
	// A bloom with a couple of bits set, so a zero value cannot pass by accident.
	var bloom Bloom

	bloom[0] = 0x01
	bloom[128] = 0xff
	bloom[len(bloom)-1] = 0x80

	root := common.HexToHash("0x1234567890123456789012345678901234567890123456789012345678901234")

	raw := fmt.Sprintf(`{
		"type": "0x0",
		"root": "%s",
		"cumulativeGasUsed": "0x5208",
		"gasUsed": "0x5208",
		"logsBloom": "0x%x",
		"transactionHash": "0xeeee00000000000000000000000000000000000000000000000000000000000c",
		"logs": []
	}`, root.Hex(), bloom)

	var receipt Receipt
	if err := json.Unmarshal([]byte(raw), &receipt); err != nil {
		t.Fatalf("failed decoding receipt: %v", err)
	}

	if receipt.Bloom != bloom {
		t.Fatalf("bloom did not decode: got %x", receipt.Bloom)
	}

	if !bytes.Equal(receipt.PostState, root.Bytes()) {
		t.Fatalf("post state did not decode: got %x", receipt.PostState)
	}

	// Both directions of the go-ethereum conversion must carry them.
	gethReceipt := receipt.ToGethReceipt()
	if gethReceipt.Bloom != bloom || !bytes.Equal(gethReceipt.PostState, root.Bytes()) {
		t.Fatal("ToGethReceipt dropped the bloom or post state")
	}

	back := FromGethReceipt(gethReceipt)
	if back.Bloom != bloom || !bytes.Equal(back.PostState, root.Bytes()) {
		t.Fatal("FromGethReceipt dropped the bloom or post state")
	}

	// A receipt without either field stays usable.
	var bare Receipt
	if err := json.Unmarshal([]byte(`{"transactionHash":"0xeeee00000000000000000000000000000000000000000000000000000000000d","logs":[]}`), &bare); err != nil {
		t.Fatalf("failed decoding bare receipt: %v", err)
	}

	if bare.Bloom != (Bloom{}) || bare.PostState != nil {
		t.Fatal("absent fields should stay zero")
	}
}

// TestTransactionJSONRoundTrip checks that a transaction survives MarshalJSON followed
// by UnmarshalJSON with its content intact, for every registered type.
func TestTransactionJSONRoundTrip(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed generating key: %v", err)
	}

	var (
		sender  = crypto.PubkeyToAddress(key.PublicKey)
		chainID = big.NewInt(7088110746)
		rng     = rand.New(rand.NewSource(7))
	)

	pairs := generateTxPairs(t, rng, chainID)
	pairs = append(pairs, txPair{name: "frame", ours: buildTransferTx(t, sender)})

	for _, pair := range pairs {
		t.Run(pair.name, func(t *testing.T) {
			signed, err := SignTx(NewTx(pair.ours), chainID, key)
			if err != nil {
				t.Fatalf("failed signing: %v", err)
			}

			encoded, err := json.Marshal(signed)
			if err != nil {
				t.Fatalf("failed marshalling: %v", err)
			}

			var decoded Transaction
			if err := json.Unmarshal(encoded, &decoded); err != nil {
				t.Fatalf("failed unmarshalling: %v", err)
			}

			if decoded.Type() != signed.Type() {
				t.Fatalf("type changed: 0x%02x != 0x%02x", decoded.Type(), signed.Type())
			}

			if decoded.Hash() != signed.Hash() {
				t.Fatalf("hash changed: %s != %s", decoded.Hash(), signed.Hash())
			}

			if _, isUnknown := decoded.Inner().(*UnknownTx); isUnknown {
				t.Fatal("round trip lost the transaction type")
			}

			// The adopted hash is cached, so recompute from the decoded content to
			// confirm every covered field survived.
			if rebuilt := NewTx(decoded.Inner()); rebuilt.Hash() != signed.Hash() {
				t.Fatalf("reconstructed content hashes to %s, want %s", rebuilt.Hash(), signed.Hash())
			}

			from, err := decoded.From(chainID)
			if err != nil {
				t.Fatalf("failed reading sender: %v", err)
			}

			if from != sender {
				t.Fatalf("wrong sender: %s", from)
			}
		})
	}
}

// TestTransactionUnmarshalJSONEmpty checks that decoding a malformed object fails
// rather than leaving a transaction that panics on first use.
func TestTransactionUnmarshalJSONEmpty(t *testing.T) {
	var tx Transaction
	if err := json.Unmarshal([]byte(`{"type":"0x2"}`), &tx); err == nil {
		t.Fatal("expected an error decoding a transaction object without a hash")
	}

	if _, err := json.Marshal(&Transaction{}); err == nil {
		t.Fatal("expected an error marshalling an empty transaction")
	}
}

// TestReceiptJSONRoundTrip checks that MarshalJSON and UnmarshalJSON are inverses,
// including the frame-specific content.
func TestReceiptJSONRoundTrip(t *testing.T) {
	raw := `{
		"type": "0x6",
		"status": "0x1",
		"cumulativeGasUsed": "0x5261",
		"gasUsed": "0x5261",
		"effectiveGasPrice": "0x77359407",
		"blobGasUsed": "0x20000",
		"blobGasPrice": "0x3",
		"transactionHash": "0x6c7a2d7d1eb9a754781378632ca05f528c626ea0928efbdbdaf226bef960e172",
		"transactionIndex": "0x2",
		"blockHash": "0xd3c53d7b84a70b9398e414eef14e778ae29975e1c126a37f13a22171a1f6ad96",
		"blockNumber": "0x10b",
		"payer": "0x6df35438a4dfcdbd25c7a364ab77e3cfdce87fc5",
		"frameReceipts": [
			{"status":"0x1","gasUsed":"0x33","stateGasUsed":"0x0","logs":[]},
			{"status":"0x2","gasUsed":"0x0","stateGasUsed":"0x7","logs":[]}
		],
		"logs": [],
		"logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004"
	}`

	var first Receipt
	if err := json.Unmarshal([]byte(raw), &first); err != nil {
		t.Fatalf("failed decoding: %v", err)
	}

	encoded, err := json.Marshal(&first)
	if err != nil {
		t.Fatalf("failed marshalling: %v", err)
	}

	var second Receipt
	if err := json.Unmarshal(encoded, &second); err != nil {
		t.Fatalf("failed decoding the re-encoded receipt: %v", err)
	}

	if second.Type != first.Type || second.Status != first.Status ||
		second.GasUsed != first.GasUsed || second.CumulativeGasUsed != first.CumulativeGasUsed ||
		second.TxHash != first.TxHash || second.BlockHash != first.BlockHash ||
		second.TransactionIndex != first.TransactionIndex ||
		second.BlobGasUsed != first.BlobGasUsed || second.Bloom != first.Bloom {
		t.Fatal("generic receipt fields did not survive the round trip")
	}

	if second.BlockNumber.Cmp(first.BlockNumber) != 0 ||
		second.EffectiveGasPrice.Cmp(first.EffectiveGasPrice) != 0 ||
		second.BlobGasPrice.Cmp(first.BlobGasPrice) != 0 {
		t.Fatal("receipt fee fields did not survive the round trip")
	}

	before, after := first.FrameExtra(), second.FrameExtra()
	if before == nil || after == nil {
		t.Fatal("frame content missing")
	}

	if before.Payer != after.Payer || len(before.Frames) != len(after.Frames) {
		t.Fatal("frame receipt content did not survive the round trip")
	}

	for i := range before.Frames {
		if before.Frames[i].Status != after.Frames[i].Status ||
			before.Frames[i].ExecutionGas != after.Frames[i].ExecutionGas ||
			before.Frames[i].StateGas != after.Frames[i].StateGas {
			t.Fatalf("frame receipt %d changed", i)
		}
	}

	if !after.Frames[1].Skipped() {
		t.Fatal("the skipped status did not survive the round trip")
	}
}

// TestReceiptLogsWithoutPosition checks that logs nested in a frame receipt decode
// even though they carry none of the position fields go-ethereum's Log requires, and
// that they inherit the receipt's position.
func TestReceiptLogsWithoutPosition(t *testing.T) {
	raw := `{
		"type": "0x6",
		"status": "0x1",
		"cumulativeGasUsed": "0x10",
		"gasUsed": "0x10",
		"transactionHash": "0xaaaa00000000000000000000000000000000000000000000000000000000000e",
		"blockHash": "0xbbbb00000000000000000000000000000000000000000000000000000000000f",
		"blockNumber": "0x2a",
		"transactionIndex": "0x3",
		"payer": "0xdddd000000000000000000000000000000000004",
		"logs": [],
		"frameReceipts": [
			{"status":"0x1","gasUsed":"0x0","stateGasUsed":"0x0","logs":[
				{"address":"0xcccc000000000000000000000000000000000005",
				 "topics":["0x1111111111111111111111111111111111111111111111111111111111111111"],
				 "data":"0xdeadbeef"}
			]}
		]
	}`

	var receipt Receipt
	if err := json.Unmarshal([]byte(raw), &receipt); err != nil {
		t.Fatalf("a frame log without position fields should decode: %v", err)
	}

	extra := receipt.FrameExtra()
	if extra == nil || len(extra.Frames) != 1 || len(extra.Frames[0].Logs) != 1 {
		t.Fatal("frame log was not decoded")
	}

	log := extra.Frames[0].Logs[0]
	if log.TxHash != receipt.TxHash || log.BlockHash != receipt.BlockHash ||
		log.BlockNumber != 42 || log.TxIndex != 3 {
		t.Fatalf("log did not inherit the receipt position: %+v", log)
	}

	if log.Address != common.HexToAddress("0xcccc000000000000000000000000000000000005") ||
		len(log.Topics) != 1 || !bytes.Equal(log.Data, []byte{0xde, 0xad, 0xbe, 0xef}) {
		t.Fatal("log content was not decoded")
	}

	// A receipt-level log keeps its own position when it reports one.
	withPos := `{"transactionHash":"0xaaaa000000000000000000000000000000000000000000000000000000000010",
		"blockNumber":"0x1","transactionIndex":"0x0","logs":[
		{"address":"0xcccc000000000000000000000000000000000005","topics":[],"data":"0x",
		 "logIndex":"0x7","transactionIndex":"0x9"}]}`

	var second Receipt
	if err := json.Unmarshal([]byte(withPos), &second); err != nil {
		t.Fatalf("failed decoding: %v", err)
	}

	if second.Logs[0].Index != 7 || second.Logs[0].TxIndex != 9 {
		t.Fatal("log's own position fields should win")
	}
}
