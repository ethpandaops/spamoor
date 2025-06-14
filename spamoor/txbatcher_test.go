package spamoor

import (
	"context"
	"testing"

	"github.com/holiman/uint256"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ethpandaops/spamoor/spamoortypes"
	testingutils "github.com/ethpandaops/spamoor/testing/utils"
)

func TestTxBatcher_Basic(t *testing.T) {
	ctx := context.Background()
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

	batcher := NewTxBatcher(txPool)
	assert.NotNil(t, batcher)
	assert.Equal(t, txPool, batcher.txpool)
	assert.False(t, batcher.isDeployed)
}

func TestTxBatcher_GetRequestCalldata(t *testing.T) {
	ctx := context.Background()
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
	batcher := NewTxBatcher(txPool)

	wallet1, _ := NewWallet("")
	wallet2, _ := NewWallet("")

	requests := []*FundingRequest{
		{Wallet: wallet1, Amount: uint256.NewInt(1000)},
		{Wallet: wallet2, Amount: uint256.NewInt(2000)},
	}

	calldata, err := batcher.GetRequestCalldata(requests)
	assert.NoError(t, err)
	assert.Equal(t, 64, len(calldata)) // 2 requests * 32 bytes each
}
