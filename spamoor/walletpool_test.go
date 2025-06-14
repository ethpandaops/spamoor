package spamoor

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ethpandaops/spamoor/spamoortypes"
	testingutils "github.com/ethpandaops/spamoor/testing/utils"
)

func TestWalletPool_Basic(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}

	// Create dependencies
	mockClient := testingutils.NewMockClient()
	rootWallet, err := InitRootWallet(ctx, "", mockClient, logger)
	require.NoError(t, err)

	clientPool := NewClientPool(ctx, []string{"mock://localhost:8545"}, logger)
	clientPool.clientFactory = func(rpchost string) (spamoortypes.Client, error) {
		return testingutils.NewMockClient(), nil
	}
	err = clientPool.PrepareClients()
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

	// Create wallet pool
	walletPool := NewWalletPool(ctx, logger, rootWallet, clientPool, txPool)
	assert.NotNil(t, walletPool)
	assert.Equal(t, ctx, walletPool.GetContext())
	assert.Equal(t, txPool, walletPool.GetTxPool())
	assert.Equal(t, clientPool, walletPool.GetClientPool())
	assert.Equal(t, rootWallet, walletPool.GetRootWallet())
}

func TestWalletPool_Configuration(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}

	mockClient := testingutils.NewMockClient()
	rootWallet, err := InitRootWallet(ctx, "", mockClient, logger)
	require.NoError(t, err)

	clientPool := NewClientPool(ctx, []string{"mock://localhost:8545"}, logger)
	clientPool.clientFactory = func(rpchost string) (spamoortypes.Client, error) {
		return testingutils.NewMockClient(), nil
	}
	err = clientPool.PrepareClients()
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

	walletPool := NewWalletPool(ctx, logger, rootWallet, clientPool, txPool)

	// Test configuration methods
	walletPool.SetWalletCount(10)
	assert.Equal(t, uint64(10), walletPool.GetConfiguredWalletCount())

	walletPool.SetRefillAmount(uint256.NewInt(5000000000000000000))
	walletPool.SetRefillBalance(uint256.NewInt(1000000000000000000))
	walletPool.SetWalletSeed("test-seed")
	walletPool.SetRefillInterval(300)
	walletPool.SetRunFundings(false)

	// Test well-known wallet configuration
	wellKnownConfig := &spamoortypes.WellKnownWalletConfig{
		Name:          "test-wallet",
		RefillAmount:  uint256.NewInt(2000000000000000000),
		RefillBalance: uint256.NewInt(500000000000000000),
		VeryWellKnown: false,
	}
	walletPool.AddWellKnownWallet(wellKnownConfig)
}

func TestGetDefaultWalletConfig(t *testing.T) {
	config := GetDefaultWalletConfig("test-scenario")
	assert.NotNil(t, config)
	assert.Contains(t, config.WalletSeed, "test-scenario")
	assert.Equal(t, uint64(0), config.WalletCount)
	assert.Equal(t, uint256.NewInt(5000000000000000000), config.RefillAmount)
	assert.Equal(t, uint256.NewInt(1000000000000000000), config.RefillBalance)
	assert.Equal(t, uint64(600), config.RefillInterval)
}

// TestWalletPool_WalletSelection tests wallet selection functionality
func TestWalletPool_WalletSelection(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}

	mockClient := testingutils.NewMockClient()
	rootWallet, err := InitRootWallet(ctx, "", mockClient, logger)
	require.NoError(t, err)

	clientPool := NewClientPool(ctx, []string{"mock://localhost:8545"}, logger)
	clientPool.clientFactory = func(rpchost string) (spamoortypes.Client, error) {
		return testingutils.NewMockClient(), nil
	}
	err = clientPool.PrepareClients()
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

	walletPool := NewWalletPool(ctx, logger, rootWallet, clientPool, txPool)

	// Test with no wallets
	wallet := walletPool.GetWallet(spamoortypes.SelectWalletByIndex, 0)
	assert.Nil(t, wallet)

	// Test wallet count
	assert.Equal(t, uint64(0), walletPool.GetWalletCount())
	assert.Equal(t, uint64(0), walletPool.GetConfiguredWalletCount())

	// Test all wallets
	allWallets := walletPool.GetAllWallets()
	assert.Equal(t, 0, len(allWallets))
}

// TestWalletPool_WellKnownWallets tests well-known wallet functionality
func TestWalletPool_WellKnownWallets(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}

	mockClient := testingutils.NewMockClient()
	rootWallet, err := InitRootWallet(ctx, "", mockClient, logger)
	require.NoError(t, err)

	clientPool := NewClientPool(ctx, []string{"mock://localhost:8545"}, logger)
	clientPool.clientFactory = func(rpchost string) (spamoortypes.Client, error) {
		return testingutils.NewMockClient(), nil
	}
	err = clientPool.PrepareClients()
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

	walletPool := NewWalletPool(ctx, logger, rootWallet, clientPool, txPool)

	// Test getting non-existent well-known wallet
	wallet := walletPool.GetWellKnownWallet("non-existent")
	assert.Nil(t, wallet)

	// Test very well known wallet address derivation
	address := walletPool.GetVeryWellKnownWalletAddress("test-wallet")
	assert.NotEqual(t, common.Address{}, address)

	// Test wallet name for root wallet
	name := walletPool.GetWalletName(rootWallet.GetWallet().GetAddress())
	assert.Equal(t, "root", name)

	// Test wallet name for unknown address
	unknownName := walletPool.GetWalletName(common.Address{})
	assert.Equal(t, "unknown", unknownName)
}

// TestWalletPool_TransactionTracker tests transaction tracking functionality
func TestWalletPool_TransactionTracker(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}

	mockClient := testingutils.NewMockClient()
	rootWallet, err := InitRootWallet(ctx, "", mockClient, logger)
	require.NoError(t, err)

	clientPool := NewClientPool(ctx, []string{"mock://localhost:8545"}, logger)
	clientPool.clientFactory = func(rpchost string) (spamoortypes.Client, error) {
		return testingutils.NewMockClient(), nil
	}
	err = clientPool.PrepareClients()
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

	walletPool := NewWalletPool(ctx, logger, rootWallet, clientPool, txPool)

	// Test initially no tracker
	tracker := walletPool.GetTransactionTracker()
	assert.Nil(t, tracker)

	// Test setting tracker
	called := false
	testTracker := func(err error) {
		called = true
	}
	walletPool.SetTransactionTracker(testTracker)

	gotTracker := walletPool.GetTransactionTracker()
	assert.NotNil(t, gotTracker)

	// Test calling tracker
	gotTracker(nil)
	assert.True(t, called)
}
