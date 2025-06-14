package spamoor

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	testingutils "github.com/ethpandaops/spamoor/testing/utils"
)

func TestNewWallet_GeneratesNewKey(t *testing.T) {
	w, err := NewWallet("")
	require.NoError(t, err)
	assert.NotNil(t, w)
	assert.IsType(t, &ecdsa.PrivateKey{}, w.GetPrivateKey())
	assert.NotEqual(t, common.Address{}, w.GetAddress())
}

func TestNewWallet_FromHexKey(t *testing.T) {
	w1, _ := NewWallet("")
	privHex := w1.GetPrivateKey().D.Text(16)
	w2, err := NewWallet(privHex)
	require.NoError(t, err)
	assert.Equal(t, w1.GetAddress(), w2.GetAddress())
}

func TestNewWallet_InvalidKey(t *testing.T) {
	_, err := NewWallet("notakey")
	assert.Error(t, err)
}

func TestWallet_NonceManagement(t *testing.T) {
	w, _ := NewWallet("")
	assert.Equal(t, uint64(0), w.GetNonce())
	w.SetNonce(5)
	assert.Equal(t, uint64(5), w.GetNonce())
	w.SetNonce(3) // Should not decrease
	assert.Equal(t, uint64(5), w.GetNonce())
	n := w.GetNextNonce()
	assert.Equal(t, uint64(5), n)
	assert.Equal(t, uint64(6), w.GetNonce())
}

func TestWallet_ConfirmedNonce(t *testing.T) {
	w, _ := NewWallet("")
	w.SetNonce(7)
	assert.Equal(t, uint64(7), w.GetConfirmedNonce())
}

func TestWallet_BalanceManagement(t *testing.T) {
	w, _ := NewWallet("")
	assert.Nil(t, w.GetBalance())
	w.SetBalance(big.NewInt(100))
	assert.Equal(t, big.NewInt(100), w.GetBalance())
	w.AddBalance(big.NewInt(50))
	assert.Equal(t, big.NewInt(150), w.GetBalance())
	w.SubBalance(big.NewInt(20))
	assert.Equal(t, big.NewInt(130), w.GetBalance())
}

func TestWallet_ChainId(t *testing.T) {
	w, _ := NewWallet("")
	assert.Nil(t, w.GetChainId())
	w.SetChainId(big.NewInt(42))
	assert.Equal(t, big.NewInt(42), w.GetChainId())
}

func TestWallet_AddressSetter(t *testing.T) {
	w, _ := NewWallet("")
	addr := common.HexToAddress("0x1234")
	w.SetAddress(addr)
	assert.Equal(t, addr, w.GetAddress())
}

func TestWallet_ConcurrentNonceAndBalance(t *testing.T) {
	w, _ := NewWallet("")
	w.SetBalance(big.NewInt(0))
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func() {
			w.GetNextNonce()
			wg.Done()
		}()
		go func() {
			w.AddBalance(big.NewInt(1))
			wg.Done()
		}()
	}
	wg.Wait()
	assert.Equal(t, uint64(100), w.GetNonce())
	assert.Equal(t, big.NewInt(100), w.GetBalance())
}

func TestWallet_ResetPendingNonce(t *testing.T) {
	w, _ := NewWallet("")
	w.SetNonce(5)
	mock := testingutils.NewMockClient()
	mock.SetMockNonce(10)
	w.ResetPendingNonce(context.Background(), mock)
	assert.Equal(t, uint64(10), w.GetNonce())
}

func TestWallet_BuildDynamicFeeTx(t *testing.T) {
	w, _ := NewWallet("")
	w.SetChainId(big.NewInt(1))
	txData := &types.DynamicFeeTx{
		To:        &common.Address{},
		Gas:       21000,
		GasFeeCap: big.NewInt(100),
		GasTipCap: big.NewInt(2),
		Value:     big.NewInt(1),
	}
	tx, err := w.BuildDynamicFeeTx(txData)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	assert.Equal(t, uint64(0), tx.Nonce())
}

func TestWallet_BuildBlobTx(t *testing.T) {
	w, _ := NewWallet("")
	w.SetChainId(big.NewInt(1))
	txData := &types.BlobTx{
		To:    common.Address{},
		Gas:   21000,
		Value: uint256.NewInt(1),
	}
	tx, err := w.BuildBlobTx(txData)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	assert.Equal(t, uint64(0), tx.Nonce())
}

func TestWallet_ReplaceDynamicFeeTx(t *testing.T) {
	w, _ := NewWallet("")
	w.SetChainId(big.NewInt(1))
	txData := &types.DynamicFeeTx{
		To:        &common.Address{},
		Gas:       21000,
		GasFeeCap: big.NewInt(100),
		GasTipCap: big.NewInt(2),
		Value:     big.NewInt(1),
	}
	tx, err := w.ReplaceDynamicFeeTx(txData, 42)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	assert.Equal(t, uint64(42), tx.Nonce())
}

func TestWallet_ReplaceBlobTx(t *testing.T) {
	w, _ := NewWallet("")
	w.SetChainId(big.NewInt(1))
	txData := &types.BlobTx{
		To:    common.Address{},
		Gas:   21000,
		Value: uint256.NewInt(1),
	}
	tx, err := w.ReplaceBlobTx(txData, 42)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	assert.Equal(t, uint64(42), tx.Nonce())
}
