package txtypes

import (
	"bytes"
	"crypto/ecdsa"
	"math/big"
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/holiman/uint256"
)

// Differential tests against go-ethereum: for every transaction type it implements,
// our encoding must be byte-identical and hashes, signing digests and signatures must
// match.

// txPair is one generated transaction in both representations.
type txPair struct {
	name string
	ours TxData
	geth types.TxData
}

// generateTxPairs builds equivalent transactions of every supported type.
func generateTxPairs(t *testing.T, rng *rand.Rand, chainID *big.Int) []txPair {
	t.Helper()

	var (
		to         = randAddress(rng)
		value      = randBig(rng, 32)
		data       = randBytes(rng, rng.Intn(64))
		accessList = randAccessList(rng)
		authList   = randAuthList(rng)
		blobHashes = randBlobHashes(rng)
		nonce      = rng.Uint64()
		gas        = rng.Uint64() % 30_000_000
		gasPrice   = randBig(rng, 8)
		tipCap     = randBig(rng, 8)
		feeCap     = randBig(rng, 8)
		blobCap    = randBig(rng, 8)
	)

	// Half the legacy-style transactions are contract creations.
	toPtr := &to
	if rng.Intn(2) == 0 {
		toPtr = nil
	}

	pairs := []txPair{
		{
			name: "legacy",
			ours: &LegacyTx{
				Nonce: nonce, GasPrice: gasPrice, Gas: gas,
				To: toPtr, Value: value, Data: data,
			},
			geth: &types.LegacyTx{
				Nonce: nonce, GasPrice: gasPrice, Gas: gas,
				To: toPtr, Value: value, Data: data,
			},
		},
		{
			name: "accesslist",
			ours: &AccessListTx{
				ChainID: chainID, Nonce: nonce, GasPrice: gasPrice, Gas: gas,
				To: toPtr, Value: value, Data: data, AccessList: accessList,
			},
			geth: &types.AccessListTx{
				ChainID: chainID, Nonce: nonce, GasPrice: gasPrice, Gas: gas,
				To: toPtr, Value: value, Data: data, AccessList: accessList,
			},
		},
		{
			name: "dynfee",
			ours: &DynamicFeeTx{
				ChainID: chainID, Nonce: nonce, GasTipCap: tipCap, GasFeeCap: feeCap,
				Gas: gas, To: toPtr, Value: value, Data: data, AccessList: accessList,
			},
			geth: &types.DynamicFeeTx{
				ChainID: chainID, Nonce: nonce, GasTipCap: tipCap, GasFeeCap: feeCap,
				Gas: gas, To: toPtr, Value: value, Data: data, AccessList: accessList,
			},
		},
		{
			name: "setcode",
			ours: &SetCodeTx{
				ChainID: uint256.MustFromBig(chainID), Nonce: nonce,
				GasTipCap: uint256.MustFromBig(tipCap), GasFeeCap: uint256.MustFromBig(feeCap),
				Gas: gas, To: to, Value: uint256.MustFromBig(value), Data: data,
				AccessList: accessList, AuthList: authList,
			},
			geth: &types.SetCodeTx{
				ChainID: uint256.MustFromBig(chainID), Nonce: nonce,
				GasTipCap: uint256.MustFromBig(tipCap), GasFeeCap: uint256.MustFromBig(feeCap),
				Gas: gas, To: to, Value: uint256.MustFromBig(value), Data: data,
				AccessList: accessList, AuthList: authList,
			},
		},
	}

	// Blob transactions in all three sidecar states: absent, v0 and v1.
	for _, sidecarVersion := range []int{-1, int(BlobSidecarVersion0), int(BlobSidecarVersion1)} {
		var ourSidecar, gethSidecar *types.BlobTxSidecar

		name := "blob/no-sidecar"

		if sidecarVersion >= 0 {
			sidecar := randSidecar(rng, byte(sidecarVersion), len(blobHashes))
			ourSidecar, gethSidecar = sidecar, sidecar
			name = "blob/sidecar-v" + string(rune('0'+sidecarVersion))
		}

		pairs = append(pairs, txPair{
			name: name,
			ours: &BlobTx{
				ChainID: uint256.MustFromBig(chainID), Nonce: nonce,
				GasTipCap: uint256.MustFromBig(tipCap), GasFeeCap: uint256.MustFromBig(feeCap),
				Gas: gas, To: to, Value: uint256.MustFromBig(value), Data: data,
				AccessList: accessList, BlobFeeCap: uint256.MustFromBig(blobCap),
				BlobHashes: blobHashes, Sidecar: ourSidecar,
			},
			geth: &types.BlobTx{
				ChainID: uint256.MustFromBig(chainID), Nonce: nonce,
				GasTipCap: uint256.MustFromBig(tipCap), GasFeeCap: uint256.MustFromBig(feeCap),
				Gas: gas, To: to, Value: uint256.MustFromBig(value), Data: data,
				AccessList: accessList, BlobFeeCap: uint256.MustFromBig(blobCap),
				BlobHashes: blobHashes, Sidecar: gethSidecar,
			},
		})
	}

	return pairs
}

// TestCodecMatchesGeth checks encoding, hashing and signing against go-ethereum over
// a deterministic corpus covering every supported transaction type.
func TestCodecMatchesGeth(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed generating key: %v", err)
	}

	sender := crypto.PubkeyToAddress(key.PublicKey)

	for _, chainID := range []*big.Int{big.NewInt(1), big.NewInt(7088110746), new(big.Int).Lsh(big.NewInt(1), 200)} {
		rng := rand.New(rand.NewSource(int64(chainID.Uint64()) + 1))

		for round := 0; round < 64; round++ {
			for _, pair := range generateTxPairs(t, rng, chainID) {
				t.Run(pair.name, func(t *testing.T) {
					assertPairMatches(t, pair, chainID, key, sender)
				})
			}
		}
	}
}

// assertPairMatches runs every equivalence check on a single generated transaction.
func assertPairMatches(t *testing.T, pair txPair, chainID *big.Int, key *ecdsa.PrivateKey, sender common.Address) {
	t.Helper()

	var (
		ourTx    = NewTx(pair.ours)
		gethTx   = types.NewTx(pair.geth)
		gethSign = types.LatestSignerForChainID(chainID)
	)

	// go-ethereum's MarshalBinary emits the network encoding.
	ourBytes, err := ourTx.MarshalNetwork()
	if err != nil {
		t.Fatalf("failed encoding transaction: %v", err)
	}

	gethBytes, err := gethTx.MarshalBinary()
	if err != nil {
		t.Fatalf("failed encoding go-ethereum transaction: %v", err)
	}

	if !bytes.Equal(ourBytes, gethBytes) {
		t.Fatalf("network encoding mismatch:\n ours: %x\n geth: %x", ourBytes, gethBytes)
	}

	if ourTx.Hash() != gethTx.Hash() {
		t.Fatalf("hash mismatch: ours %s, geth %s", ourTx.Hash(), gethTx.Hash())
	}

	// The signing digest must match for signatures to be interchangeable.
	signable, ok := pair.ours.(ECDSASignedTx)
	if !ok {
		t.Fatalf("type 0x%02x does not implement ECDSASignedTx", pair.ours.TxType())
	}

	if got, want := signable.SigningHash(chainID), gethSign.Hash(gethTx); got != want {
		t.Fatalf("signing hash mismatch: ours %s, geth %s", got, want)
	}

	// Sign with both implementations and compare the results byte for byte.
	ourSigned, err := SignTx(ourTx, chainID, key)
	if err != nil {
		t.Fatalf("failed signing transaction: %v", err)
	}

	gethSigned, err := types.SignTx(gethTx, gethSign, key)
	if err != nil {
		t.Fatalf("failed signing go-ethereum transaction: %v", err)
	}

	ourSignedBytes, err := ourSigned.MarshalNetwork()
	if err != nil {
		t.Fatalf("failed encoding signed transaction: %v", err)
	}

	gethSignedBytes, err := gethSigned.MarshalBinary()
	if err != nil {
		t.Fatalf("failed encoding signed go-ethereum transaction: %v", err)
	}

	if !bytes.Equal(ourSignedBytes, gethSignedBytes) {
		t.Fatalf("signed encoding mismatch:\n ours: %x\n geth: %x", ourSignedBytes, gethSignedBytes)
	}

	if ourSigned.Hash() != gethSigned.Hash() {
		t.Fatalf("signed hash mismatch: ours %s, geth %s", ourSigned.Hash(), gethSigned.Hash())
	}

	// Recovery must work without the cached sender that signing installs.
	decoded, err := DecodeTx(ourSignedBytes)
	if err != nil {
		t.Fatalf("failed decoding signed transaction: %v", err)
	}

	if decoded.Hash() != ourSigned.Hash() {
		t.Fatalf("hash changed across decode: %s != %s", decoded.Hash(), ourSigned.Hash())
	}

	recovered, err := decoded.From(chainID)
	if err != nil {
		t.Fatalf("failed recovering sender: %v", err)
	}

	if recovered != sender {
		t.Fatalf("recovered wrong sender: %s != %s", recovered, sender)
	}

	// Sidecars must survive a decode round trip.
	assertSidecarsEqual(t, ourSigned.BlobTxSidecar(), decoded.BlobTxSidecar())

	// Both conversion directions must preserve the transaction.
	fromGeth, err := FromGethTx(gethSigned)
	if err != nil {
		t.Fatalf("failed converting from go-ethereum: %v", err)
	}

	if fromGeth.Hash() != ourSigned.Hash() {
		t.Fatalf("FromGethTx hash mismatch: %s != %s", fromGeth.Hash(), ourSigned.Hash())
	}

	assertSidecarsEqual(t, ourSigned.BlobTxSidecar(), fromGeth.BlobTxSidecar())

	toGeth, err := ourSigned.ToGethTx()
	if err != nil {
		t.Fatalf("failed converting to go-ethereum: %v", err)
	}

	if toGeth.Hash() != gethSigned.Hash() {
		t.Fatalf("ToGethTx hash mismatch: %s != %s", toGeth.Hash(), gethSigned.Hash())
	}

	// The canonical encoding must never contain the sidecar.
	canonical, err := ourSigned.MarshalBinary()
	if err != nil {
		t.Fatalf("failed encoding canonical form: %v", err)
	}

	if ourSigned.BlobTxSidecar() != nil && bytes.Equal(canonical, ourSignedBytes) {
		t.Fatal("canonical encoding unexpectedly includes the blob sidecar")
	}

	if crypto.Keccak256Hash(canonical) != ourSigned.Hash() {
		t.Fatal("transaction hash does not cover the canonical encoding")
	}
}

// assertSidecarsEqual compares two blob sidecars, allowing both to be absent.
func assertSidecarsEqual(t *testing.T, want, got *BlobSidecar) {
	t.Helper()

	if want == nil {
		if got != nil {
			t.Fatal("sidecar appeared where none was expected")
		}

		return
	}

	if got == nil {
		t.Fatal("sidecar lost")
	}

	if want.Version != got.Version {
		t.Fatalf("sidecar version mismatch: %d != %d", want.Version, got.Version)
	}

	if len(want.Blobs) != len(got.Blobs) || len(want.Commitments) != len(got.Commitments) || len(want.Proofs) != len(got.Proofs) {
		t.Fatal("sidecar contents changed length")
	}

	for i := range want.Blobs {
		if want.Blobs[i] != got.Blobs[i] {
			t.Fatalf("blob %d changed", i)
		}
	}

	for i := range want.Commitments {
		if want.Commitments[i] != got.Commitments[i] {
			t.Fatalf("commitment %d changed", i)
		}
	}

	for i := range want.Proofs {
		if want.Proofs[i] != got.Proofs[i] {
			t.Fatalf("proof %d changed", i)
		}
	}
}

// TestUnsupportedTxTypeDecode checks that an unregistered type byte fails cleanly.
func TestUnsupportedTxTypeDecode(t *testing.T) {
	if _, err := DecodeTx([]byte{0x7f, 0xc0}); err == nil {
		t.Fatal("expected an error decoding an unregistered transaction type")
	}
}

func randBytes(rng *rand.Rand, n int) []byte {
	if n == 0 {
		return nil
	}

	b := make([]byte, n)
	rng.Read(b)

	return b
}

func randAddress(rng *rand.Rand) common.Address {
	var addr common.Address

	rng.Read(addr[:])

	return addr
}

func randHash(rng *rand.Rand) common.Hash {
	var hash common.Hash

	rng.Read(hash[:])

	return hash
}

// randBig returns a random value of at most maxBytes, covering every width so RLP's
// minimal-length integer encoding is exercised.
func randBig(rng *rand.Rand, maxBytes int) *big.Int {
	n := rng.Intn(maxBytes + 1)
	if n == 0 {
		return new(big.Int)
	}

	return new(big.Int).SetBytes(randBytes(rng, n))
}

func randAccessList(rng *rand.Rand) AccessList {
	count := rng.Intn(3)
	if count == 0 {
		return nil
	}

	list := make(AccessList, count)
	for i := range list {
		keys := make([]common.Hash, rng.Intn(3))
		for j := range keys {
			keys[j] = randHash(rng)
		}

		list[i] = AccessTuple{Address: randAddress(rng), StorageKeys: keys}
	}

	return list
}

func randAuthList(rng *rand.Rand) []SetCodeAuthorization {
	count := rng.Intn(3)
	if count == 0 {
		return nil
	}

	list := make([]SetCodeAuthorization, count)
	for i := range list {
		list[i] = SetCodeAuthorization{
			ChainID: *uint256.MustFromBig(randBig(rng, 8)),
			Address: randAddress(rng),
			Nonce:   rng.Uint64(),
			V:       byte(rng.Intn(2)),
			R:       *uint256.MustFromBig(randBig(rng, 32)),
			S:       *uint256.MustFromBig(randBig(rng, 32)),
		}
	}

	return list
}

func randBlobHashes(rng *rand.Rand) []common.Hash {
	hashes := make([]common.Hash, 1+rng.Intn(3))
	for i := range hashes {
		hashes[i] = randHash(rng)
		hashes[i][0] = 0x01
	}

	return hashes
}

// randSidecar builds a sidecar with arbitrary contents. The values are not valid KZG
// material; this test covers the wire format only.
func randSidecar(rng *rand.Rand, version byte, blobCount int) *BlobSidecar {
	proofCount := blobCount
	if version == BlobSidecarVersion1 {
		proofCount = blobCount * kzg4844.CellProofsPerBlob
	}

	sidecar := &BlobSidecar{
		Version:     version,
		Blobs:       make([]kzg4844.Blob, blobCount),
		Commitments: make([]kzg4844.Commitment, blobCount),
		Proofs:      make([]kzg4844.Proof, proofCount),
	}

	for i := range sidecar.Blobs {
		rng.Read(sidecar.Blobs[i][:])
		rng.Read(sidecar.Commitments[i][:])
	}

	for i := range sidecar.Proofs {
		rng.Read(sidecar.Proofs[i][:])
	}

	return sidecar
}
