package spamoor

import (
	"context"
	"crypto/ecdsa"
	"fmt"
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
	// Convert private key to proper hex format (64 characters, padded with zeros)
	privHex := fmt.Sprintf("%064x", w1.GetPrivateKey().D)
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

func TestWallet_BuildSetCodeTx(t *testing.T) {
	w, _ := NewWallet("")
	w.SetChainId(big.NewInt(1))
	txData := &types.SetCodeTx{
		To:    common.Address{},
		Gas:   21000,
		Value: uint256.NewInt(1),
	}
	tx, err := w.BuildSetCodeTx(txData)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	assert.Equal(t, uint64(0), tx.Nonce())
}

func TestWallet_SetConfirmedNonce(t *testing.T) {
	w, _ := NewWallet("")
	w.SetConfirmedNonce(10)
	assert.Equal(t, uint64(10), w.GetConfirmedNonce())
}

func TestWallet_SetLastConfirmation(t *testing.T) {
	w, _ := NewWallet("")
	w.SetLastConfirmation(100)
	assert.Equal(t, uint64(100), w.GetLastConfirmation())
}

func TestWallet_GetTxNonceChan(t *testing.T) {
	w, _ := NewWallet("")
	w.SetConfirmedNonce(5)

	// Test getting nonce chan for future nonce
	nonceChan, isFirst := w.GetTxNonceChan(10)
	assert.NotNil(t, nonceChan)
	assert.True(t, isFirst)
	assert.NotNil(t, nonceChan.Channel)

	// Test getting same nonce chan again
	nonceChan2, isFirst2 := w.GetTxNonceChan(10)
	assert.Equal(t, nonceChan, nonceChan2)
	assert.False(t, isFirst2)

	// Test getting nonce chan for already confirmed nonce
	nonceChan3, isFirst3 := w.GetTxNonceChan(3)
	assert.Nil(t, nonceChan3)
	assert.False(t, isFirst3)
}

func TestWallet_ProcessTransactionInclusion(t *testing.T) {
	w, _ := NewWallet("")
	w.SetChainId(big.NewInt(1))

	// Create a transaction
	txData := &types.DynamicFeeTx{
		To:        &common.Address{},
		Gas:       21000,
		GasFeeCap: big.NewInt(100),
		GasTipCap: big.NewInt(2),
		Value:     big.NewInt(1),
		Nonce:     5,
	}
	tx, _ := w.ReplaceDynamicFeeTx(txData, 5)

	// Get nonce chan for this transaction
	nonceChan, _ := w.GetTxNonceChan(5)

	// Create a mock receipt
	receipt := &types.Receipt{
		Status: 1,
	}

	// Process transaction inclusion
	w.ProcessTransactionInclusion(100, tx, receipt)

	// Check that confirmed nonce was updated
	assert.Equal(t, uint64(6), w.GetConfirmedNonce())
	assert.Equal(t, uint64(100), w.GetLastConfirmation())

	// Check that nonce channel was closed and receipt set
	select {
	case <-nonceChan.Channel:
		// Channel should be closed
	default:
		t.Error("Expected nonce channel to be closed")
	}
	assert.Equal(t, receipt, nonceChan.Receipt)
}

func TestWallet_ProcessTransactionReceival(t *testing.T) {
	w, _ := NewWallet("")
	w.SetBalance(big.NewInt(100))
	w.SetChainId(big.NewInt(1))

	// Create a transaction with value
	txData := &types.DynamicFeeTx{
		To:        &common.Address{},
		Gas:       21000,
		GasFeeCap: big.NewInt(100),
		GasTipCap: big.NewInt(2),
		Value:     big.NewInt(50),
	}
	tx, _ := w.BuildDynamicFeeTx(txData)

	w.ProcessTransactionReceival(tx)
	assert.Equal(t, big.NewInt(150), w.GetBalance())
}

func TestWallet_RevertTransactionReceival(t *testing.T) {
	w, _ := NewWallet("")
	w.SetBalance(big.NewInt(100))
	w.SetChainId(big.NewInt(1))

	// Create a transaction with value
	txData := &types.DynamicFeeTx{
		To:        &common.Address{},
		Gas:       21000,
		GasFeeCap: big.NewInt(100),
		GasTipCap: big.NewInt(2),
		Value:     big.NewInt(30),
	}
	tx, _ := w.BuildDynamicFeeTx(txData)

	w.RevertTransactionReceival(tx)
	assert.Equal(t, big.NewInt(70), w.GetBalance())
}

func TestWallet_GetPendingTxCount(t *testing.T) {
	w, _ := NewWallet("")
	assert.Equal(t, 0, w.GetPendingTxCount())

	// Add some pending transactions
	w.GetTxNonceChan(10)
	w.GetTxNonceChan(11)
	w.GetTxNonceChan(12)

	assert.Equal(t, 3, w.GetPendingTxCount())
}

func TestWallet_ProcessStaleTransactions(t *testing.T) {
	w, _ := NewWallet("")

	// Add some pending transactions
	nonceChan1, _ := w.GetTxNonceChan(5)
	nonceChan2, _ := w.GetTxNonceChan(6)
	nonceChan3, _ := w.GetTxNonceChan(8)

	assert.Equal(t, 3, w.GetPendingTxCount())

	// Process stale transactions (nonce 7 means 5 and 6 are stale)
	w.ProcessStaleTransactions(100, 7)

	// Check that stale transactions were cleaned up
	assert.Equal(t, 1, w.GetPendingTxCount())

	// Check that channels were closed
	select {
	case <-nonceChan1.Channel:
		// Channel should be closed
	default:
		t.Error("Expected nonce channel 1 to be closed")
	}

	select {
	case <-nonceChan2.Channel:
		// Channel should be closed
	default:
		t.Error("Expected nonce channel 2 to be closed")
	}

	// Channel 3 should still be open
	select {
	case <-nonceChan3.Channel:
		t.Error("Expected nonce channel 3 to remain open")
	default:
		// Channel should still be open
	}
}

func TestWallet_HexKeyWithPrefix(t *testing.T) {
	w1, _ := NewWallet("")
	privHex := fmt.Sprintf("%064x", w1.GetPrivateKey().D)

	// Test with 0x prefix
	w2, err := NewWallet("0x" + privHex)
	require.NoError(t, err)
	assert.Equal(t, w1.GetAddress(), w2.GetAddress())

	// Test with 0X prefix (uppercase)
	w3, err := NewWallet("0X" + privHex)
	require.NoError(t, err)
	assert.Equal(t, w1.GetAddress(), w3.GetAddress())
}

func TestWallet_LoadPrivateKeyErrors(t *testing.T) {
	// Test invalid hex characters
	_, err := NewWallet("gggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg")
	assert.Error(t, err)

	// Test too short key
	_, err = NewWallet("1234")
	assert.Error(t, err)
}

func TestWallet_SignTxError(t *testing.T) {
	w, _ := NewWallet("")
	// Don't set chain ID to cause signing error

	txData := &types.DynamicFeeTx{
		To:        &common.Address{},
		Gas:       21000,
		GasFeeCap: big.NewInt(100),
		GasTipCap: big.NewInt(2),
		Value:     big.NewInt(1),
	}

	_, err := w.BuildDynamicFeeTx(txData)
	assert.Error(t, err)
}
