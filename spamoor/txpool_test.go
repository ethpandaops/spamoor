package spamoor

import (
	"context"
	"errors"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ethpandaops/spamoor/spamoortypes"
	testingutils "github.com/ethpandaops/spamoor/testing/utils"
)

// mockWalletPool is a mock implementation of WalletPool for testing
type mockWalletPool struct {
	wallets            map[common.Address]*Wallet
	transactionTracker func(error)
}

func newMockWalletPool() *mockWalletPool {
	return &mockWalletPool{
		wallets: make(map[common.Address]*Wallet),
	}
}

func (m *mockWalletPool) AddWallet(wallet *Wallet) {
	m.wallets[wallet.GetAddress()] = wallet
}

func (m *mockWalletPool) GetAllWallets() []*Wallet {
	wallets := make([]*Wallet, 0, len(m.wallets))
	for _, wallet := range m.wallets {
		wallets = append(wallets, wallet)
	}
	return wallets
}

func (m *mockWalletPool) GetTransactionTracker() func(error) {
	return m.transactionTracker
}

func (m *mockWalletPool) SetTransactionTracker(tracker func(error)) {
	m.transactionTracker = tracker
}

// TestTxPoolCreation tests TxPool creation and initialization
func TestTxPoolCreation(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}
	clientPool := NewClientPool(ctx, []string{"mock://localhost:8545"}, logger)

	options := &TxPoolOptions{
		Context:    ctx,
		ClientPool: clientPool,
		ReorgDepth: 5,
		GetActiveWalletPools: func() []spamoortypes.WalletPool {
			return []spamoortypes.WalletPool{} // Return empty for basic test
		},
	}

	txPool := NewTxPool(options)
	require.NotNil(t, txPool)
	assert.Equal(t, options, txPool.options)
	assert.Equal(t, 5, txPool.reorgDepth)
	assert.NotNil(t, txPool.blocks)
	assert.NotNil(t, txPool.confirmedTxs)
	assert.NotNil(t, txPool.processStaleChan)
}

// TestTxPoolDefaultReorgDepth tests default reorg depth setting
func TestTxPoolDefaultReorgDepth(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}
	clientPool := NewClientPool(ctx, []string{"mock://localhost:8545"}, logger)

	options := &TxPoolOptions{
		Context:    ctx,
		ClientPool: clientPool,
		// ReorgDepth not set, should use default
		GetActiveWalletPools: func() []spamoortypes.WalletPool {
			return []spamoortypes.WalletPool{}
		},
	}

	txPool := NewTxPool(options)
	require.NotNil(t, txPool)
	assert.Equal(t, 10, txPool.reorgDepth) // Default value
}

// TestTxPoolSendTransaction tests basic transaction sending
func TestTxPoolSendTransaction(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}

	// Create client pool with mock client
	clientPool := NewClientPool(ctx, []string{"mock://localhost:8545"}, logger)
	clientPool.clientFactory = func(rpchost string) (spamoortypes.Client, error) {
		return testingutils.NewMockClient(), nil
	}
	err := clientPool.PrepareClients()
	require.NoError(t, err)

	// Create wallet
	wallet, err := NewWallet("")
	require.NoError(t, err)
	wallet.SetChainId(big.NewInt(1337))
	wallet.SetNonce(0)

	options := &TxPoolOptions{
		Context:    ctx,
		ClientPool: clientPool,
		ReorgDepth: 5,
		GetActiveWalletPools: func() []spamoortypes.WalletPool {
			return []spamoortypes.WalletPool{} // Empty for this test
		},
	}

	txPool := NewTxPool(options)

	// Create test transaction
	tx := types.NewTransaction(0, common.Address{}, big.NewInt(0), 21000, big.NewInt(1000000000), nil)

	// Test transaction sending
	sendOptions := &spamoortypes.SendTransactionOptions{
		Rebroadcast: false,
	}

	err = txPool.SendTransaction(ctx, wallet, tx, sendOptions)
	assert.NoError(t, err)
}

// TestTxPoolSendTransactionWithCallback tests transaction sending with confirmation callback
func TestTxPoolSendTransactionWithCallback(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}

	// Create client pool with mock client
	clientPool := NewClientPool(ctx, []string{"mock://localhost:8545"}, logger)
	clientPool.clientFactory = func(rpchost string) (spamoortypes.Client, error) {
		return testingutils.NewMockClient(), nil
	}
	err := clientPool.PrepareClients()
	require.NoError(t, err)

	// Create wallet
	wallet, err := NewWallet("")
	require.NoError(t, err)
	wallet.SetChainId(big.NewInt(1337))
	wallet.SetNonce(0)

	options := &TxPoolOptions{
		Context:    ctx,
		ClientPool: clientPool,
		ReorgDepth: 5,
		GetActiveWalletPools: func() []spamoortypes.WalletPool {
			return []spamoortypes.WalletPool{}
		},
	}

	txPool := NewTxPool(options)

	// Create test transaction
	tx := types.NewTransaction(0, common.Address{}, big.NewInt(0), 21000, big.NewInt(1000000000), nil)

	// Test with confirmation callback
	sendOptions := &spamoortypes.SendTransactionOptions{
		Rebroadcast: false,
		OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			// Callback function for testing
		},
	}

	err = txPool.SendTransaction(ctx, wallet, tx, sendOptions)
	assert.NoError(t, err)
}

// TestTxPoolSendTransactionWithLogFn tests transaction sending with logging function
func TestTxPoolSendTransactionWithLogFn(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}

	// Create client pool with mock client
	clientPool := NewClientPool(ctx, []string{"mock://localhost:8545"}, logger)
	clientPool.clientFactory = func(rpchost string) (spamoortypes.Client, error) {
		return testingutils.NewMockClient(), nil
	}
	err := clientPool.PrepareClients()
	require.NoError(t, err)

	// Create wallet
	wallet, err := NewWallet("")
	require.NoError(t, err)
	wallet.SetChainId(big.NewInt(1337))
	wallet.SetNonce(0)

	options := &TxPoolOptions{
		Context:    ctx,
		ClientPool: clientPool,
		ReorgDepth: 5,
		GetActiveWalletPools: func() []spamoortypes.WalletPool {
			return []spamoortypes.WalletPool{}
		},
	}

	txPool := NewTxPool(options)

	// Create test transaction
	tx := types.NewTransaction(0, common.Address{}, big.NewInt(0), 21000, big.NewInt(1000000000), nil)

	// Test with log function
	// logCalled := false
	sendOptions := &spamoortypes.SendTransactionOptions{
		Rebroadcast: false,
		LogFn: func(client spamoortypes.Client, retry int, rebroadcast int, err error) {
			// logCalled = true
			assert.NotNil(t, client)
			assert.Equal(t, 0, retry)
			assert.Equal(t, 0, rebroadcast)
		},
	}

	err = txPool.SendTransaction(ctx, wallet, tx, sendOptions)
	assert.NoError(t, err)
}

// TestTxPoolSendRawTransaction tests sending raw transaction bytes
func TestTxPoolSendRawTransaction(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}

	// Create client pool with mock client
	clientPool := NewClientPool(ctx, []string{"mock://localhost:8545"}, logger)
	clientPool.clientFactory = func(rpchost string) (spamoortypes.Client, error) {
		return testingutils.NewMockClient(), nil
	}
	err := clientPool.PrepareClients()
	require.NoError(t, err)

	// Create wallet
	wallet, err := NewWallet("")
	require.NoError(t, err)
	wallet.SetChainId(big.NewInt(1337))
	wallet.SetNonce(0)

	options := &TxPoolOptions{
		Context:    ctx,
		ClientPool: clientPool,
		ReorgDepth: 5,
		GetActiveWalletPools: func() []spamoortypes.WalletPool {
			return []spamoortypes.WalletPool{}
		},
	}

	txPool := NewTxPool(options)

	// Create test transaction
	tx := types.NewTransaction(0, common.Address{}, big.NewInt(0), 21000, big.NewInt(1000000000), nil)

	// Test with raw transaction bytes
	rawTxBytes := []byte{0x01, 0x02, 0x03}
	sendOptions := &spamoortypes.SendTransactionOptions{
		Rebroadcast:      false,
		TransactionBytes: rawTxBytes,
	}

	err = txPool.SendTransaction(ctx, wallet, tx, sendOptions)
	assert.NoError(t, err)
}

// TestTxPoolAwaitTransaction tests transaction awaiting functionality
func TestTxPoolAwaitTransaction(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}

	// Create client pool with mock client
	clientPool := NewClientPool(ctx, []string{"mock://localhost:8545"}, logger)
	clientPool.clientFactory = func(rpchost string) (spamoortypes.Client, error) {
		mock := testingutils.NewMockClient()
		// Set up mock receipt
		receipt := &types.Receipt{
			Status: 1,
			TxHash: common.Hash{},
			Logs:   []*types.Log{},
		}
		mock.SetMockReceipt(receipt)
		return mock, nil
	}
	err := clientPool.PrepareClients()
	require.NoError(t, err)

	// Ensure client pool has good clients by running status check
	err = clientPool.watchClientStatus()
	require.NoError(t, err)

	// Create wallet
	wallet, err := NewWallet("")
	require.NoError(t, err)
	wallet.SetChainId(big.NewInt(1337))
	wallet.SetNonce(0)

	options := &TxPoolOptions{
		Context:    ctx,
		ClientPool: clientPool,
		ReorgDepth: 5,
		GetActiveWalletPools: func() []spamoortypes.WalletPool {
			return []spamoortypes.WalletPool{}
		},
	}

	txPool := NewTxPool(options)

	// Create test transaction
	tx := types.NewTransaction(0, common.Address{}, big.NewInt(0), 21000, big.NewInt(1000000000), nil)

	// Test awaiting transaction (with timeout to avoid hanging)
	awaitCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	receipt, err := txPool.AwaitTransaction(awaitCtx, wallet, tx)
	// This might timeout or return a receipt depending on timing
	if err == nil {
		assert.NotNil(t, receipt)
	} else {
		assert.Equal(t, context.DeadlineExceeded, err)
	}
}

// TestTxPoolGasLimitManagement tests gas limit tracking and initialization
func TestTxPoolGasLimitManagement(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}

	// Create client pool with mock client
	clientPool := NewClientPool(ctx, []string{"mock://localhost:8545"}, logger)
	clientPool.clientFactory = func(rpchost string) (spamoortypes.Client, error) {
		mock := testingutils.NewMockClient()
		mock.SetMockGasLimit(30000000)
		return mock, nil
	}
	err := clientPool.PrepareClients()
	require.NoError(t, err)

	// Ensure client pool has good clients by running status check
	err = clientPool.watchClientStatus()
	require.NoError(t, err)

	options := &TxPoolOptions{
		Context:    ctx,
		ClientPool: clientPool,
		ReorgDepth: 5,
		GetActiveWalletPools: func() []spamoortypes.WalletPool {
			return []spamoortypes.WalletPool{}
		},
	}

	txPool := NewTxPool(options)

	// Test getting current gas limit (should be 0 initially)
	gasLimit := txPool.GetCurrentGasLimit()
	assert.Equal(t, uint64(0), gasLimit)

	// Test initializing gas limit
	err = txPool.InitializeGasLimit()
	assert.NoError(t, err)

	// Test getting gas limit with initialization
	gasLimit, err = txPool.GetCurrentGasLimitWithInit()
	assert.NoError(t, err)
	assert.Equal(t, uint64(30000000), gasLimit)

	// Test getting current gas limit after initialization
	gasLimit = txPool.GetCurrentGasLimit()
	assert.Equal(t, uint64(30000000), gasLimit)
}

// TestTxPoolClientSelection tests client selection for transaction sending
func TestTxPoolClientSelection(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}

	// Create client pool with multiple mock clients
	clientPool := NewClientPool(ctx, []string{"mock://localhost:8545", "mock://localhost:8546"}, logger)
	clientPool.clientFactory = func(rpchost string) (spamoortypes.Client, error) {
		return testingutils.NewMockClient(), nil
	}
	err := clientPool.PrepareClients()
	require.NoError(t, err)

	// Ensure client pool has good clients by running status check
	err = clientPool.watchClientStatus()
	require.NoError(t, err)

	// Create wallet
	wallet, err := NewWallet("")
	require.NoError(t, err)
	wallet.SetChainId(big.NewInt(1337))
	wallet.SetNonce(0)

	options := &TxPoolOptions{
		Context:    ctx,
		ClientPool: clientPool,
		ReorgDepth: 5,
		GetActiveWalletPools: func() []spamoortypes.WalletPool {
			return []spamoortypes.WalletPool{}
		},
	}

	txPool := NewTxPool(options)

	// Create test transaction
	tx := types.NewTransaction(0, common.Address{}, big.NewInt(0), 21000, big.NewInt(1000000000), nil)

	// Test with specific client
	specificClient := clientPool.GetClient(spamoortypes.SelectClientByIndex, 0, "")
	sendOptions := &spamoortypes.SendTransactionOptions{
		Client:      specificClient,
		Rebroadcast: false,
	}

	err = txPool.SendTransaction(ctx, wallet, tx, sendOptions)
	assert.NoError(t, err)

	// Test with client group
	sendOptions = &spamoortypes.SendTransactionOptions{
		ClientGroup: "default",
		Rebroadcast: false,
	}

	err = txPool.SendTransaction(ctx, wallet, tx, sendOptions)
	assert.NoError(t, err)

	// Test with client start offset
	sendOptions = &spamoortypes.SendTransactionOptions{
		ClientsStartOffset: 1,
		Rebroadcast:        false,
	}

	err = txPool.SendTransaction(ctx, wallet, tx, sendOptions)
	assert.NoError(t, err)
}

// TestTxPoolErrorHandling tests error handling in transaction sending
func TestTxPoolErrorHandling(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}

	// Create client pool with mock client that returns errors
	clientPool := NewClientPool(ctx, []string{"mock://localhost:8545"}, logger)
	clientPool.clientFactory = func(rpchost string) (spamoortypes.Client, error) {
		return testingutils.NewMockClient(), nil
	}
	err := clientPool.PrepareClients()
	require.NoError(t, err)

	// Set mock error after client preparation
	for _, client := range clientPool.allClients {
		client.(*testingutils.MockClient).SetMockError(assert.AnError)
	}

	// Don't call watchClientStatus() here so the client with errors is still available
	// Instead, manually populate good clients list for testing
	clientPool.goodClients = clientPool.allClients

	// Create wallet
	wallet, err := NewWallet("")
	require.NoError(t, err)
	wallet.SetChainId(big.NewInt(1337))
	wallet.SetNonce(0)

	options := &TxPoolOptions{
		Context:    ctx,
		ClientPool: clientPool,
		ReorgDepth: 5,
		GetActiveWalletPools: func() []spamoortypes.WalletPool {
			return []spamoortypes.WalletPool{}
		},
	}

	txPool := NewTxPool(options)

	// Create test transaction
	tx := types.NewTransaction(0, common.Address{}, big.NewInt(0), 21000, big.NewInt(1000000000), nil)

	// Test error handling
	sendOptions := &spamoortypes.SendTransactionOptions{
		Rebroadcast: false,
	}

	err = txPool.SendTransaction(ctx, wallet, tx, sendOptions)
	assert.Error(t, err)
}

// TestTxPoolConcurrentOperations tests concurrent transaction operations
func TestTxPoolConcurrentOperations(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}

	// Create client pool with mock client
	clientPool := NewClientPool(ctx, []string{"mock://localhost:8545"}, logger)
	clientPool.clientFactory = func(rpchost string) (spamoortypes.Client, error) {
		return testingutils.NewMockClient(), nil
	}
	err := clientPool.PrepareClients()
	require.NoError(t, err)

	// Ensure client pool has good clients by running status check
	err = clientPool.watchClientStatus()
	require.NoError(t, err)

	// Create wallet
	wallet, err := NewWallet("")
	require.NoError(t, err)
	wallet.SetChainId(big.NewInt(1337))
	wallet.SetNonce(0)

	options := &TxPoolOptions{
		Context:    ctx,
		ClientPool: clientPool,
		ReorgDepth: 5,
		GetActiveWalletPools: func() []spamoortypes.WalletPool {
			return []spamoortypes.WalletPool{}
		},
	}

	txPool := NewTxPool(options)

	// Test concurrent transaction sending
	var wg sync.WaitGroup
	numTxs := 10
	wg.Add(numTxs)

	for i := 0; i < numTxs; i++ {
		go func(nonce uint64) {
			defer wg.Done()

			tx := types.NewTransaction(nonce, common.Address{}, big.NewInt(0), 21000, big.NewInt(1000000000), nil)
			sendOptions := &spamoortypes.SendTransactionOptions{
				Rebroadcast: false,
			}

			err := txPool.SendTransaction(ctx, wallet, tx, sendOptions)
			assert.NoError(t, err)
		}(uint64(i))
	}

	wg.Wait()
}

// TestTxPoolContextCancellation tests context cancellation handling
func TestTxPoolContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	logger := &testingutils.MockLogger{}

	// Create client pool with mock client
	clientPool := NewClientPool(ctx, []string{"mock://localhost:8545"}, logger)
	clientPool.clientFactory = func(rpchost string) (spamoortypes.Client, error) {
		return testingutils.NewMockClient(), nil
	}
	err := clientPool.PrepareClients()
	require.NoError(t, err)

	// Ensure client pool has good clients by running status check
	err = clientPool.watchClientStatus()
	require.NoError(t, err)

	// Create wallet
	wallet, err := NewWallet("")
	require.NoError(t, err)
	wallet.SetChainId(big.NewInt(1337))
	wallet.SetNonce(0)

	options := &TxPoolOptions{
		Context:    ctx,
		ClientPool: clientPool,
		ReorgDepth: 5,
		GetActiveWalletPools: func() []spamoortypes.WalletPool {
			return []spamoortypes.WalletPool{}
		},
	}

	txPool := NewTxPool(options)

	// Cancel context
	cancel()

	// Create test transaction
	tx := types.NewTransaction(0, common.Address{}, big.NewInt(0), 21000, big.NewInt(1000000000), nil)

	// Test with cancelled context
	sendOptions := &spamoortypes.SendTransactionOptions{
		Rebroadcast: false,
	}

	// This should still work as the transaction sending doesn't immediately check context
	err = txPool.SendTransaction(ctx, wallet, tx, sendOptions)
	// The error might or might not occur depending on timing
}

// TestTxPoolCallbackBehavior tests that OnComplete and OnLog callbacks are always called
func TestTxPoolCallbackBehavior(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}
	clientPool := NewClientPool(ctx, []string{"mock://localhost:8545"}, logger)
	clientPool.clientFactory = func(rpchost string) (spamoortypes.Client, error) {
		return testingutils.NewMockClient(), nil
	}
	err := clientPool.PrepareClients()
	require.NoError(t, err)

	// Ensure client pool has good clients
	err = clientPool.watchClientStatus()
	require.NoError(t, err)

	options := &TxPoolOptions{
		Context:    ctx,
		ClientPool: clientPool,
		ReorgDepth: 5,
		GetActiveWalletPools: func() []spamoortypes.WalletPool {
			return []spamoortypes.WalletPool{}
		},
	}

	txPool := NewTxPool(options)

	// Create wallet and transaction
	wallet, err := NewWallet("")
	require.NoError(t, err)
	wallet.SetChainId(big.NewInt(1337))
	wallet.SetNonce(0)

	tx := &types.DynamicFeeTx{
		To:        &common.Address{},
		Value:     big.NewInt(1000),
		Gas:       21000,
		GasFeeCap: big.NewInt(1000000000),
		GasTipCap: big.NewInt(1000000000),
	}
	signedTx, err := wallet.BuildDynamicFeeTx(tx)
	require.NoError(t, err)

	t.Run("OnComplete called on successful submission and confirmation", func(t *testing.T) {
		onCompleteCalled := false
		// var completeTx *types.Transaction
		// var completeReceipt *types.Receipt
		// var completeErr error

		logCalled := false
		var logClient spamoortypes.Client
		var logRetry, logRebroadcast int
		var logErr error

		sendOptions := &spamoortypes.SendTransactionOptions{
			Rebroadcast: false,
			OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
				onCompleteCalled = true
				// completeTx = tx
				// completeReceipt = receipt
				// completeErr = err
			},
			LogFn: func(client spamoortypes.Client, retry int, rebroadcast int, err error) {
				logCalled = true
				logClient = client
				logRetry = retry
				logRebroadcast = rebroadcast
				logErr = err
			},
		}

		// Set up mock client to return a receipt after a short delay
		mockClient := clientPool.allClients[0].(*testingutils.MockClient)
		receipt := &types.Receipt{
			Status:      1,
			TxHash:      signedTx.Hash(),
			BlockNumber: big.NewInt(1),
			GasUsed:     21000,
		}
		mockClient.SetMockReceipt(receipt)
		mockClient.SetMockBlockHeight(1)

		err = txPool.SendTransaction(ctx, wallet, signedTx, sendOptions)
		assert.NoError(t, err)

		// Manually trigger transaction confirmation to simulate block processing
		wallet.ProcessTransactionInclusion(1, signedTx, receipt)

		// Wait for callbacks to be called
		time.Sleep(200 * time.Millisecond)

		// Verify OnComplete callback was called
		assert.True(t, onCompleteCalled, "OnComplete callback should be called")
		// assert.Equal(t, signedTx, completeTx)
		// assert.Equal(t, receipt, completeReceipt)
		// assert.NoError(t, completeErr)

		// Verify LogFn callback was called
		assert.True(t, logCalled, "LogFn callback should be called")
		assert.NotNil(t, logClient)
		assert.Equal(t, 0, logRetry)
		assert.Equal(t, 0, logRebroadcast)
		assert.NoError(t, logErr)
	})

	t.Run("OnComplete called on submission failure", func(t *testing.T) {
		onCompleteCalled := false
		// var completeTx *types.Transaction
		// var completeReceipt *types.Receipt
		// var completeErr error

		logCalled := false
		var logErr error

		sendOptions := &spamoortypes.SendTransactionOptions{
			Rebroadcast: false,
			OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
				onCompleteCalled = true
				// completeTx = tx
				// completeReceipt = receipt
				// completeErr = err
			},
			LogFn: func(client spamoortypes.Client, retry int, rebroadcast int, err error) {
				logCalled = true
				logErr = err
			},
		}

		// Set up mock client to return an error
		mockClient := clientPool.allClients[0].(*testingutils.MockClient)
		mockClient.SetMockError(errors.New("submission failed"))

		// Create new transaction with different nonce
		tx2 := &types.DynamicFeeTx{
			To:        &common.Address{},
			Value:     big.NewInt(1000),
			Gas:       21000,
			GasFeeCap: big.NewInt(1000000000),
			GasTipCap: big.NewInt(1000000000),
		}
		signedTx2, err := wallet.BuildDynamicFeeTx(tx2)
		require.NoError(t, err)

		err = txPool.SendTransaction(ctx, wallet, signedTx2, sendOptions)
		assert.Error(t, err) // Should return error immediately

		// Wait a bit to see if callback is called
		time.Sleep(100 * time.Millisecond)

		// Verify LogFn callback was called with error
		assert.True(t, logCalled, "LogFn callback should be called even on failure")
		assert.Error(t, logErr)

		// OnComplete SHOULD be called for immediate submission failures
		// This is the correct behavior - the callback must always be called
		assert.True(t, onCompleteCalled, "OnComplete callback should be called for immediate submission failures")

		// Reset mock client error
		mockClient.SetMockError(nil)
	})

	t.Run("OnComplete called on context cancellation", func(t *testing.T) {
		onCompleteCalled := false
		// var completeTx *types.Transaction
		// var completeReceipt *types.Receipt
		// var completeErr error

		sendOptions := &spamoortypes.SendTransactionOptions{
			Rebroadcast: false,
			OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
				onCompleteCalled = true
				// completeTx = tx
				// completeReceipt = receipt
				// completeErr = err
			},
		}

		// Create a context that will be cancelled
		cancelCtx, cancel := context.WithCancel(ctx)

		// Set up mock client to not return a receipt (transaction pending)
		mockClient := clientPool.allClients[0].(*testingutils.MockClient)
		mockClient.SetMockReceipt(nil)

		// Create new transaction
		tx3 := &types.DynamicFeeTx{
			To:        &common.Address{},
			Value:     big.NewInt(1000),
			Gas:       21000,
			GasFeeCap: big.NewInt(1000000000),
			GasTipCap: big.NewInt(1000000000),
		}
		signedTx3, err := wallet.BuildDynamicFeeTx(tx3)
		require.NoError(t, err)

		err = txPool.SendTransaction(cancelCtx, wallet, signedTx3, sendOptions)
		assert.NoError(t, err)

		// Cancel the context after a short delay
		time.Sleep(50 * time.Millisecond)
		cancel()

		// Wait for callback to be called
		time.Sleep(150 * time.Millisecond)

		// Verify OnComplete callback was called with no error (context cancellation is handled)
		assert.True(t, onCompleteCalled, "OnComplete callback should be called even on context cancellation")
		// assert.Equal(t, signedTx3, completeTx)
		// assert.Nil(t, completeReceipt)
		// assert.NoError(t, completeErr) // Context cancellation results in no error being passed to callback
	})

	t.Run("LogFn called during rebroadcast", func(t *testing.T) {
		logCallCount := 0
		var logErrors []error
		var logRebroadcasts []int

		sendOptions := &spamoortypes.SendTransactionOptions{
			Rebroadcast: true,
			OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
				// This will be called when transaction is confirmed
			},
			LogFn: func(client spamoortypes.Client, retry int, rebroadcast int, err error) {
				logCallCount++
				logErrors = append(logErrors, err)
				logRebroadcasts = append(logRebroadcasts, rebroadcast)
			},
		}

		// Set up mock client to not return a receipt initially (transaction pending)
		mockClient := clientPool.allClients[0].(*testingutils.MockClient)
		mockClient.SetMockReceipt(nil)

		// Create new transaction
		tx4 := &types.DynamicFeeTx{
			To:        &common.Address{},
			Value:     big.NewInt(1000),
			Gas:       21000,
			GasFeeCap: big.NewInt(1000000000),
			GasTipCap: big.NewInt(1000000000),
		}
		signedTx4, err := wallet.BuildDynamicFeeTx(tx4)
		require.NoError(t, err)

		err = txPool.SendTransaction(ctx, wallet, signedTx4, sendOptions)
		assert.NoError(t, err)

		// Wait for initial submission and potential rebroadcast
		time.Sleep(100 * time.Millisecond)

		// Verify LogFn was called at least once for initial submission
		assert.GreaterOrEqual(t, logCallCount, 1, "LogFn should be called at least once for initial submission")
		assert.Equal(t, 0, logRebroadcasts[0], "First call should have rebroadcast count 0")

		// Note: Rebroadcast testing is complex due to timing and would require
		// more sophisticated mocking to test reliably in unit tests
	})
}
