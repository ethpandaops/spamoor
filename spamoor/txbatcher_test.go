package spamoor

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
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
	txpool := NewTxPool(&TxPoolOptions{
		Context:    context.Background(),
		ClientPool: &ClientPool{},
	})
	batcher := NewTxBatcher(txpool)

	wallet1, _ := NewWallet("")
	wallet2, _ := NewWallet("")

	requests := []*FundingRequest{
		{
			Wallet: wallet1,
			Amount: uint256.NewInt(1000),
		},
		{
			Wallet: wallet2,
			Amount: uint256.NewInt(2000),
		},
	}

	calldata, err := batcher.GetRequestCalldata(requests)
	assert.NoError(t, err)
	assert.Len(t, calldata, 64) // 2 requests * 32 bytes each

	// Verify first request encoding
	assert.Equal(t, wallet1.GetAddress().Bytes(), calldata[0:20])
	// Amount should be in the last 12 bytes of the 32-byte slot
	expectedAmount1 := uint256.NewInt(1000).Bytes32()
	assert.Equal(t, expectedAmount1[20:], calldata[20:32])

	// Verify second request encoding
	assert.Equal(t, wallet2.GetAddress().Bytes(), calldata[32:52])
	expectedAmount2 := uint256.NewInt(2000).Bytes32()
	assert.Equal(t, expectedAmount2[20:], calldata[52:64])
}

func TestTxBatcher_Deploy(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}
	clientPool := NewClientPool(ctx, []string{}, logger)

	mock := testingutils.NewMockClient()
	mock.SetMockChainId(big.NewInt(1))
	mock.SetMockGasFees(big.NewInt(500000000000), big.NewInt(300000000000))
	clientPool.allClients = []spamoortypes.Client{mock}
	clientPool.goodClients = []spamoortypes.Client{mock}

	txpool := NewTxPool(&TxPoolOptions{
		Context:    ctx,
		ClientPool: clientPool,
	})

	batcher := NewTxBatcher(txpool)
	wallet, _ := NewWallet("")
	wallet.SetChainId(big.NewInt(1))

	// Test initial state
	assert.Equal(t, common.Address{}, batcher.GetAddress())

	// Test deployment
	err := batcher.Deploy(ctx, wallet, mock)
	assert.NoError(t, err)

	// Test that address is set after deployment
	assert.NotEqual(t, common.Address{}, batcher.GetAddress())

	// Test that second deployment is a no-op
	oldAddress := batcher.GetAddress()
	err = batcher.Deploy(ctx, wallet, mock)
	assert.NoError(t, err)
	assert.Equal(t, oldAddress, batcher.GetAddress())
}

func TestTxBatcher_DeployWithLowGasFees(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}
	clientPool := NewClientPool(ctx, []string{}, logger)

	mock := testingutils.NewMockClient()
	mock.SetMockChainId(big.NewInt(1))
	// Set low gas fees to test the minimum fee logic
	mock.SetMockGasFees(big.NewInt(100000000000), big.NewInt(50000000000))
	clientPool.allClients = []spamoortypes.Client{mock}
	clientPool.goodClients = []spamoortypes.Client{mock}

	txpool := NewTxPool(&TxPoolOptions{
		Context:    ctx,
		ClientPool: clientPool,
	})

	batcher := NewTxBatcher(txpool)
	wallet, _ := NewWallet("")
	wallet.SetChainId(big.NewInt(1))

	err := batcher.Deploy(ctx, wallet, mock)
	assert.NoError(t, err)
	assert.NotEqual(t, common.Address{}, batcher.GetAddress())
}

func TestTxBatcher_DeployWithNilClient(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}
	clientPool := NewClientPool(ctx, []string{}, logger)

	mock := testingutils.NewMockClient()
	mock.SetMockChainId(big.NewInt(1))
	mock.SetMockGasFees(big.NewInt(500000000000), big.NewInt(300000000000))
	clientPool.allClients = []spamoortypes.Client{mock}
	clientPool.goodClients = []spamoortypes.Client{mock}

	txpool := NewTxPool(&TxPoolOptions{
		Context:    ctx,
		ClientPool: clientPool,
	})

	batcher := NewTxBatcher(txpool)
	wallet, _ := NewWallet("")
	wallet.SetChainId(big.NewInt(1))

	// Test deployment with nil client (should use pool's client)
	err := batcher.Deploy(ctx, wallet, nil)
	assert.NoError(t, err)
	assert.NotEqual(t, common.Address{}, batcher.GetAddress())
}

func TestTxBatcher_DeployErrors(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}
	clientPool := NewClientPool(ctx, []string{}, logger)

	// Test with no clients available
	txpool := NewTxPool(&TxPoolOptions{
		Context:    ctx,
		ClientPool: clientPool,
	})

	batcher := NewTxBatcher(txpool)
	wallet, _ := NewWallet("")
	wallet.SetChainId(big.NewInt(1))

	err := batcher.Deploy(ctx, wallet, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no client available")
}

func TestTxBatcher_DeployGetSuggestedFeeError(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}
	clientPool := NewClientPool(ctx, []string{}, logger)

	mock := testingutils.NewMockClient()
	mock.SetMockError(errors.New("fee error"))
	clientPool.allClients = []spamoortypes.Client{mock}
	clientPool.goodClients = []spamoortypes.Client{mock}

	txpool := NewTxPool(&TxPoolOptions{
		Context:    ctx,
		ClientPool: clientPool,
	})

	batcher := NewTxBatcher(txpool)
	wallet, _ := NewWallet("")
	wallet.SetChainId(big.NewInt(1))

	err := batcher.Deploy(ctx, wallet, mock)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "fee error")
}
