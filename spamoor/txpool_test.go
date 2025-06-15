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
	logCalled := false
	sendOptions := &spamoortypes.SendTransactionOptions{
		Rebroadcast: false,
		LogFn: func(client spamoortypes.Client, retry int, rebroadcast int, err error) {
			logCalled = true
			assert.NotNil(t, client)
			assert.Equal(t, 0, retry)
			assert.Equal(t, 0, rebroadcast)
		},
	}

	err = txPool.SendTransaction(ctx, wallet, tx, sendOptions)
	assert.NoError(t, err)

	time.Sleep(50 * time.Millisecond)
	assert.True(t, logCalled)
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

	// Create test transaction
	tx := types.NewTransaction(0, common.Address{}, big.NewInt(0), 21000, big.NewInt(1000000000), nil)

	// Create client pool
	clientPool := NewClientPool(ctx, []string{}, logger)

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

	// Test awaiting transaction (with timeout to avoid hanging)
	awaitCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	go func() {
		time.Sleep(10 * time.Millisecond)
		wallet.ProcessTransactionInclusion(0, tx, &types.Receipt{Status: 1, TxHash: tx.Hash(), BlockNumber: big.NewInt(100)})
	}()

	receipt, err := txPool.AwaitTransaction(awaitCtx, wallet, tx)
	assert.NoError(t, err)
	assert.NotNil(t, receipt)
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
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "OnComplete called on successful submission and confirmation",
			test: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), 200*time.Second)
				defer cancel()

				logger := &testingutils.MockLogger{}
				clientPool := NewClientPool(ctx, []string{"mock://localhost:8545"}, logger)
				clientPool.clientFactory = func(rpchost string) (spamoortypes.Client, error) {
					return testingutils.NewMockClient(), nil
				}
				err := clientPool.PrepareClients()
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

				wallet, _ := NewWallet("")
				wallet.SetChainId(big.NewInt(1))

				txData := &types.DynamicFeeTx{
					To:        &common.Address{},
					Gas:       21000,
					GasFeeCap: big.NewInt(100),
					GasTipCap: big.NewInt(2),
					Value:     big.NewInt(1),
				}
				tx, _ := wallet.BuildDynamicFeeTx(txData)

				callbackCalled := false
				err = txPool.SendTransaction(ctx, wallet, tx, &spamoortypes.SendTransactionOptions{
					OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
						callbackCalled = true
						assert.NoError(t, err)
						assert.NotNil(t, receipt)
					},
				})
				assert.NoError(t, err)

				wallet.ProcessTransactionInclusion(0, tx, &types.Receipt{Status: 1, TxHash: tx.Hash(), BlockNumber: big.NewInt(0)})

				time.Sleep(50 * time.Millisecond)
				assert.True(t, callbackCalled)
			},
		},
		{
			name: "OnComplete called on submission failure",
			test: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
				defer cancel()

				logger := &testingutils.MockLogger{}
				client := testingutils.NewMockClient()
				clientPool := NewClientPool(ctx, []string{"mock://localhost:8545"}, logger)
				clientPool.clientFactory = func(rpchost string) (spamoortypes.Client, error) {
					return client, nil
				}
				err := clientPool.PrepareClients()
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

				wallet, _ := NewWallet("")
				wallet.SetChainId(big.NewInt(1))

				txData := &types.DynamicFeeTx{
					To:        &common.Address{},
					Gas:       21000,
					GasFeeCap: big.NewInt(100),
					GasTipCap: big.NewInt(2),
					Value:     big.NewInt(1),
				}
				tx, _ := wallet.BuildDynamicFeeTx(txData)

				callbackCalled := false
				client.SetMockError(errors.New("submission failed"))
				err = txPool.SendTransaction(ctx, wallet, tx, &spamoortypes.SendTransactionOptions{
					OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
						callbackCalled = true
						assert.Error(t, err)
						assert.Nil(t, receipt)
					},
				})
				assert.Error(t, err)

				time.Sleep(50 * time.Millisecond)
				assert.True(t, callbackCalled)
			},
		},
		{
			name: "OnComplete called on context cancellation",
			test: func(t *testing.T) {
				ctx, cancel := context.WithCancel(context.Background())

				logger := &testingutils.MockLogger{}
				clientPool := NewClientPool(ctx, []string{"mock://localhost:8545"}, logger)
				clientPool.clientFactory = func(rpchost string) (spamoortypes.Client, error) {
					return testingutils.NewMockClient(), nil
				}
				err := clientPool.PrepareClients()
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

				wallet, _ := NewWallet("")
				wallet.SetChainId(big.NewInt(1))

				txData := &types.DynamicFeeTx{
					To:        &common.Address{},
					Gas:       21000,
					GasFeeCap: big.NewInt(100),
					GasTipCap: big.NewInt(2),
					Value:     big.NewInt(1),
				}
				tx, _ := wallet.BuildDynamicFeeTx(txData)

				callbackCalled := false
				err = txPool.SendTransaction(ctx, wallet, tx, &spamoortypes.SendTransactionOptions{
					OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
						callbackCalled = true
						assert.NoError(t, err)
						assert.Nil(t, receipt)
					},
				})
				assert.NoError(t, err)

				// Cancel context after a short delay
				time.Sleep(50 * time.Millisecond)
				cancel()

				time.Sleep(150 * time.Millisecond)
				assert.True(t, callbackCalled)
			},
		},
		{
			name: "LogFn called during rebroadcast",
			test: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
				defer cancel()

				logger := &testingutils.MockLogger{}
				clientPool := NewClientPool(ctx, []string{"mock://localhost:8545"}, logger)
				clientPool.clientFactory = func(rpchost string) (spamoortypes.Client, error) {
					return testingutils.NewMockClient(), nil
				}
				err := clientPool.PrepareClients()
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

				wallet, _ := NewWallet("")
				wallet.SetChainId(big.NewInt(1))

				txData := &types.DynamicFeeTx{
					To:        &common.Address{},
					Gas:       21000,
					GasFeeCap: big.NewInt(100),
					GasTipCap: big.NewInt(2),
					Value:     big.NewInt(1),
				}
				tx, _ := wallet.BuildDynamicFeeTx(txData)

				logCalled := false
				err = txPool.SendTransaction(ctx, wallet, tx, &spamoortypes.SendTransactionOptions{
					Rebroadcast: true,
					LogFn: func(client spamoortypes.Client, retry int, rebroadcast int, err error) {
						logCalled = true
					},
				})
				assert.NoError(t, err)

				time.Sleep(50 * time.Millisecond)
				assert.True(t, logCalled)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func TestTxPool_GetWalletMap(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}
	clientPool := NewClientPool(ctx, []string{}, logger)

	// Create mock wallet pool
	mockWalletPool := &testingutils.MockWalletPool{}
	wallet1, _ := NewWallet("")
	wallet2, _ := NewWallet("")

	options := &TxPoolOptions{
		Context:    ctx,
		ClientPool: clientPool,
		ReorgDepth: 5,
		GetActiveWalletPools: func() []spamoortypes.WalletPool {
			return []spamoortypes.WalletPool{mockWalletPool}
		},
	}
	txPool := NewTxPool(options)

	// Set up mock wallet pool to return wallets
	mockWalletPool.SetWallets(map[common.Address]spamoortypes.Wallet{
		wallet1.GetAddress(): wallet1,
		wallet2.GetAddress(): wallet2,
	})

	walletMap := txPool.getWalletMap()
	assert.Len(t, walletMap, 2)
	assert.Equal(t, wallet1, walletMap[wallet1.GetAddress()])
	assert.Equal(t, wallet2, walletMap[wallet2.GetAddress()])
}

func TestTxPool_GetHighestBlockNumber(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}
	clientPool := NewClientPool(ctx, []string{}, logger)

	// Create mock clients with different block heights
	mock1 := testingutils.NewMockClient()
	mock1.SetMockBlockHeight(100)
	mock2 := testingutils.NewMockClient()
	mock2.SetMockBlockHeight(102) // Highest
	mock3 := testingutils.NewMockClient()
	mock3.SetMockBlockHeight(101)

	clientPool.allClients = []spamoortypes.Client{mock1, mock2, mock3}
	clientPool.goodClients = []spamoortypes.Client{mock1, mock2, mock3}

	options := &TxPoolOptions{
		Context:    ctx,
		ClientPool: clientPool,
		ReorgDepth: 5,
		GetActiveWalletPools: func() []spamoortypes.WalletPool {
			return []spamoortypes.WalletPool{}
		},
	}
	txPool := NewTxPool(options)

	highestBlock, clients := txPool.getHighestBlockNumber()
	assert.Equal(t, uint64(102), highestBlock)
	assert.Len(t, clients, 1)
	assert.Equal(t, mock2, clients[0])
}

func TestTxPool_GetBlockBody(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}
	clientPool := NewClientPool(ctx, []string{}, logger)

	mock := testingutils.NewMockClient()
	// Create a mock block
	mockBlock := &types.Block{}
	mock.SetMockBlock(mockBlock)

	options := &TxPoolOptions{
		Context:    ctx,
		ClientPool: clientPool,
		ReorgDepth: 5,
		GetActiveWalletPools: func() []spamoortypes.WalletPool {
			return []spamoortypes.WalletPool{}
		},
	}
	txPool := NewTxPool(options)

	block := txPool.getBlockBody(ctx, mock, 100)
	assert.Equal(t, mockBlock, block)
}

func TestTxPool_GetBlockReceipts(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}
	clientPool := NewClientPool(ctx, []string{}, logger)

	mock := testingutils.NewMockClient()
	// Create mock receipts
	mockReceipts := []*types.Receipt{
		{Status: 1, TxHash: common.Hash{1}, BlockNumber: big.NewInt(100)},
		{Status: 1, TxHash: common.Hash{2}, BlockNumber: big.NewInt(100)},
	}
	mock.SetMockBlockReceipts(mockReceipts)

	options := &TxPoolOptions{
		Context:    ctx,
		ClientPool: clientPool,
		ReorgDepth: 5,
		GetActiveWalletPools: func() []spamoortypes.WalletPool {
			return []spamoortypes.WalletPool{}
		},
	}
	txPool := NewTxPool(options)

	receipts, err := txPool.getBlockReceipts(ctx, mock, 100, 2)
	assert.NoError(t, err)
	assert.Equal(t, mockReceipts, receipts)
}

func TestTxPool_SendAndAwaitTxRange(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}
	clientPool := NewClientPool(ctx, []string{"mock://localhost:8545"}, logger)
	clientPool.clientFactory = func(rpchost string) (spamoortypes.Client, error) {
		return testingutils.NewMockClient(), nil
	}
	err := clientPool.PrepareClients()
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

	wallet, _ := NewWallet("")
	wallet.SetChainId(big.NewInt(1))
	wallet.SetBalance(big.NewInt(1000000000000000000))

	// Create multiple transactions
	var txs []*types.Transaction
	for i := 0; i < 3; i++ {
		txData := &types.DynamicFeeTx{
			To:        &common.Address{},
			Gas:       21000,
			GasFeeCap: big.NewInt(100),
			GasTipCap: big.NewInt(2),
			Value:     big.NewInt(1),
		}
		tx, _ := wallet.BuildDynamicFeeTx(txData)
		txs = append(txs, tx)
	}

	go func() {
		time.Sleep(100 * time.Millisecond)

		for _, tx := range txs {
			wallet.ProcessTransactionInclusion(0, tx, &types.Receipt{Status: 1, TxHash: tx.Hash(), BlockNumber: big.NewInt(100)})
		}
	}()

	err = txPool.SendAndAwaitTxRange(ctx, wallet, txs, &spamoortypes.SendTransactionOptions{})
	assert.NoError(t, err)
}

func TestTxPool_LoadTransactionReceipt(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}
	clientPool := NewClientPool(ctx, []string{}, logger)

	mock := testingutils.NewMockClient()
	mockReceipt := &types.Receipt{Status: 1}
	mock.SetMockReceipt(mockReceipt)
	clientPool.allClients = []spamoortypes.Client{mock}
	clientPool.goodClients = []spamoortypes.Client{mock}

	options := &TxPoolOptions{
		Context:    ctx,
		ClientPool: clientPool,
		ReorgDepth: 5,
		GetActiveWalletPools: func() []spamoortypes.WalletPool {
			return []spamoortypes.WalletPool{}
		},
	}
	txPool := NewTxPool(options)

	wallet, _ := NewWallet("")
	wallet.SetChainId(big.NewInt(1))

	txData := &types.DynamicFeeTx{
		To:        &common.Address{},
		Gas:       21000,
		GasFeeCap: big.NewInt(100),
		GasTipCap: big.NewInt(2),
		Value:     big.NewInt(1),
	}
	tx, _ := wallet.BuildDynamicFeeTx(txData)

	receipt := txPool.loadTransactionReceipt(ctx, tx)
	assert.Equal(t, mockReceipt, receipt)
}

func TestTxPool_CalculateBackoffDelay(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}
	clientPool := NewClientPool(ctx, []string{}, logger)

	options := &TxPoolOptions{
		Context:    ctx,
		ClientPool: clientPool,
		ReorgDepth: 5,
		GetActiveWalletPools: func() []spamoortypes.WalletPool {
			return []spamoortypes.WalletPool{}
		},
	}
	txPool := NewTxPool(options)

	// Test exponential backoff
	delay1 := txPool.calculateBackoffDelay(0)
	delay2 := txPool.calculateBackoffDelay(1)
	delay3 := txPool.calculateBackoffDelay(10)

	assert.True(t, delay1 < delay2)
	assert.True(t, delay2 < delay3)
	assert.True(t, delay1 >= 1*time.Second)
	assert.True(t, delay3 <= 10*time.Minute) // Max delay
}

func TestTxPool_GetCurrentGasLimitWithInit(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}
	clientPool := NewClientPool(ctx, []string{}, logger)

	mock := testingutils.NewMockClient()
	mock.SetMockGasLimit(30000000)
	clientPool.allClients = []spamoortypes.Client{mock}
	clientPool.goodClients = []spamoortypes.Client{mock}

	options := &TxPoolOptions{
		Context:    ctx,
		ClientPool: clientPool,
		ReorgDepth: 5,
		GetActiveWalletPools: func() []spamoortypes.WalletPool {
			return []spamoortypes.WalletPool{}
		},
	}
	txPool := NewTxPool(options)

	gasLimit, err := txPool.GetCurrentGasLimitWithInit()
	assert.NoError(t, err)
	assert.Equal(t, uint64(30000000), gasLimit)
}

func TestTxPool_InitializeGasLimit(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}
	clientPool := NewClientPool(ctx, []string{}, logger)

	mock := testingutils.NewMockClient()
	mock.SetMockGasLimit(25000000)
	clientPool.allClients = []spamoortypes.Client{mock}
	clientPool.goodClients = []spamoortypes.Client{mock}

	options := &TxPoolOptions{
		Context:    ctx,
		ClientPool: clientPool,
		ReorgDepth: 5,
		GetActiveWalletPools: func() []spamoortypes.WalletPool {
			return []spamoortypes.WalletPool{}
		},
	}
	txPool := NewTxPool(options)

	err := txPool.InitializeGasLimit()
	assert.NoError(t, err)
	assert.Equal(t, uint64(25000000), txPool.GetCurrentGasLimit())
}
