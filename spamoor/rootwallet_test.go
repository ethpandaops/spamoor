package spamoor

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ethpandaops/spamoor/spamoortypes"
	testingutils "github.com/ethpandaops/spamoor/testing/utils"
)

func TestRootWallet_Basic(t *testing.T) {
	ctx := context.Background()
	mockClient := testingutils.NewMockClient()

	rootWallet, err := InitRootWallet(ctx, "", mockClient, nil)
	require.NoError(t, err)
	assert.NotNil(t, rootWallet)
	assert.NotNil(t, rootWallet.GetWallet())
}

// TestRootWallet_WithWalletLock tests the semaphore locking mechanism
func TestRootWallet_WithWalletLock(t *testing.T) {
	ctx := context.Background()
	mockClient := testingutils.NewMockClient()

	rootWallet, err := InitRootWallet(ctx, "", mockClient, nil)
	require.NoError(t, err)

	t.Run("single lock acquisition", func(t *testing.T) {
		executed := false

		err := rootWallet.WithWalletLock(ctx, 1, nil, func() error {
			executed = true
			return nil
		})

		assert.NoError(t, err)
		assert.True(t, executed)
	})

	t.Run("function returns error", func(t *testing.T) {
		testErr := errors.New("test error")

		err := rootWallet.WithWalletLock(ctx, 1, nil, func() error {
			return testErr
		})

		assert.Equal(t, testErr, err)
	})

	t.Run("context cancellation", func(t *testing.T) {
		cancelCtx, cancel := context.WithCancel(ctx)

		// Fill up the semaphore to force waiting
		for i := 0; i < 200; i++ {
			rootWallet.txSemaphore <- struct{}{}
		}

		// Cancel context while waiting for locks
		go func() {
			time.Sleep(10 * time.Millisecond)
			cancel()
		}()

		err := rootWallet.WithWalletLock(cancelCtx, 1, nil, func() error {
			return nil
		})

		assert.Equal(t, context.Canceled, err)

		// Clean up semaphore
		for i := 0; i < 200; i++ {
			<-rootWallet.txSemaphore
		}
	})
}

// TestRootWallet_TxBatcher tests transaction batcher functionality
func TestRootWallet_TxBatcher(t *testing.T) {
	ctx := context.Background()
	mockClient := testingutils.NewMockClient()

	rootWallet, err := InitRootWallet(ctx, "", mockClient, nil)
	require.NoError(t, err)

	t.Run("initially no batcher", func(t *testing.T) {
		batcher := rootWallet.GetTxBatcher()
		assert.Nil(t, batcher)
	})

	t.Run("initialize batcher", func(t *testing.T) {
		clientPool := NewClientPool(ctx, []string{"mock://localhost:8545"}, &testingutils.MockLogger{})
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

		rootWallet.InitTxBatcher(ctx, txPool)

		batcher := rootWallet.GetTxBatcher()
		assert.NotNil(t, batcher)
		assert.Equal(t, txPool, batcher.txpool)
	})
}
