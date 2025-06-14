package spamoor

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ethpandaops/spamoor/spamoortypes"
	testingutils "github.com/ethpandaops/spamoor/testing/utils"
)

// TestClientPoolCreation tests ClientPool creation and initialization
func TestClientPoolCreation(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}
	pool := NewClientPool(ctx, []string{}, logger)
	require.NotNil(t, pool)
	assert.Equal(t, ctx, pool.ctx)
	assert.Equal(t, []string{}, pool.rpcHosts)
	assert.Equal(t, logger, pool.logger)
}

// TestClientPoolPrepareClients tests client preparation with various scenarios
func TestClientPoolPrepareClients(t *testing.T) {
	tests := []struct {
		name      string
		rpcHosts  []string
		mockError error
		wantErr   bool
	}{
		{
			name:     "no hosts",
			rpcHosts: []string{},
			wantErr:  true,
		},
		{
			name:     "valid host",
			rpcHosts: []string{"mock://localhost:8545"},
			wantErr:  false,
		},
		{
			name:      "client error",
			rpcHosts:  []string{"mock://localhost:8545"},
			mockError: errors.New("client error"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			logger := &testingutils.MockLogger{}
			pool := NewClientPool(ctx, tt.rpcHosts, logger)

			// Set up mock client factory
			pool.clientFactory = func(rpchost string) (spamoortypes.Client, error) {
				mock := testingutils.NewMockClient()
				if tt.mockError != nil {
					mock.SetMockError(tt.mockError)
				}
				return mock, nil
			}

			err := pool.PrepareClients()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, pool.allClients)
				assert.NotNil(t, pool.chainId)
			}
		})
	}
}

// TestClientPoolClientSelection tests different client selection modes
func TestClientPoolClientSelection(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}
	pool := NewClientPool(ctx, []string{"mock://localhost:8545"}, logger)

	// Create and add mock clients with different groups
	mock1 := testingutils.NewMockClient()
	mock1.SetClientGroup("group1")
	mock2 := testingutils.NewMockClient()
	mock2.SetClientGroup("group2")
	mock3 := testingutils.NewMockClient()
	mock3.SetClientGroup("default")

	pool.allClients = []spamoortypes.Client{mock1, mock2, mock3}
	pool.goodClients = []spamoortypes.Client{mock1, mock2, mock3}

	tests := []struct {
		name  string
		mode  spamoortypes.ClientSelectionMode
		input int
		group string
		want  string // Expected client group
	}{
		{
			name:  "select by index",
			mode:  spamoortypes.SelectClientByIndex,
			input: 0,
			group: "",
			want:  "default", // First client in default group
		},
		{
			name:  "select random",
			mode:  spamoortypes.SelectClientRandom,
			input: 0,
			group: "",
			want:  "default", // Random selection from default group
		},
		{
			name:  "select round robin",
			mode:  spamoortypes.SelectClientRoundRobin,
			input: 0,
			group: "",
			want:  "default", // First client in default group
		},
		{
			name:  "select by group",
			mode:  spamoortypes.SelectClientByIndex,
			input: 0,
			group: "group1",
			want:  "group1",
		},
		{
			name:  "select any group",
			mode:  spamoortypes.SelectClientByIndex,
			input: 0,
			group: "*",
			want:  "group1", // First client in any group
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := pool.GetClient(tt.mode, tt.input, tt.group)
			if tt.want != "" {
				assert.NotNil(t, client)
				assert.Equal(t, tt.want, client.GetClientGroup())
			} else {
				assert.Nil(t, client)
			}
		})
	}
}

// TestClientPoolClientStatus tests client status monitoring
func TestClientPoolClientStatus(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}
	pool := NewClientPool(ctx, []string{"mock://localhost:8545"}, logger)

	// Create mock clients with different block heights
	mock1 := testingutils.NewMockClient()
	mock1.SetMockBlockHeight(100)
	mock2 := testingutils.NewMockClient()
	mock2.SetMockBlockHeight(98)
	mock3 := testingutils.NewMockClient()
	mock3.SetMockBlockHeight(99)

	pool.allClients = []spamoortypes.Client{mock1, mock2, mock3}

	// Test status monitoring
	err := pool.watchClientStatus()
	assert.NoError(t, err)

	// Verify good clients list (should include all clients within 2 blocks)
	goodClients := pool.GetAllGoodClients()
	assert.Equal(t, 3, len(goodClients))

	// Test with one client falling behind
	mock2.SetMockBlockHeight(95)
	err = pool.watchClientStatus()
	assert.NoError(t, err)
	goodClients = pool.GetAllGoodClients()
	assert.Equal(t, 2, len(goodClients))
}

// TestClientPoolConcurrentOperations tests concurrent operations on the client pool
func TestClientPoolConcurrentOperations(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}
	pool := NewClientPool(ctx, []string{"mock://localhost:8545"}, logger)

	// Create mock clients
	mock1 := testingutils.NewMockClient()
	mock2 := testingutils.NewMockClient()
	mock3 := testingutils.NewMockClient()

	pool.allClients = []spamoortypes.Client{mock1, mock2, mock3}
	pool.goodClients = []spamoortypes.Client{mock1, mock2, mock3}

	// Test concurrent client selections
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			client := pool.GetClient(spamoortypes.SelectClientRandom, 0, "")
			assert.NotNil(t, client)
			done <- true
		}()
	}

	// Wait for all operations to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

// TestClientPoolErrorHandling tests error handling in various scenarios
func TestClientPoolErrorHandling(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}
	pool := NewClientPool(ctx, []string{"mock://localhost:8545"}, logger)

	// Create mock client with error
	mock := testingutils.NewMockClient()
	mock.SetMockError(errors.New("test error"))
	pool.allClients = []spamoortypes.Client{mock}
	pool.goodClients = []spamoortypes.Client{mock}

	// Test error handling in status monitoring
	err := pool.watchClientStatus()
	assert.NoError(t, err) // Should not return error, just log warning
	goodClients := pool.GetAllGoodClients()
	assert.Equal(t, 1, len(goodClients)) // Client should remain in good clients until next status check

	// Test error handling in client selection
	client := pool.GetClient(spamoortypes.SelectClientByIndex, 0, "")
	assert.NotNil(t, client) // Client should still be available
}

// TestClientPoolClientGroups tests client group functionality
func TestClientPoolClientGroups(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}
	pool := NewClientPool(ctx, []string{"mock://localhost:8545"}, logger)

	// Create mock clients in different groups
	mock1 := testingutils.NewMockClient()
	mock1.SetClientGroup("group1")
	mock2 := testingutils.NewMockClient()
	mock2.SetClientGroup("group2")
	mock3 := testingutils.NewMockClient()
	mock3.SetClientGroup("default")

	pool.allClients = []spamoortypes.Client{mock1, mock2, mock3}
	pool.goodClients = []spamoortypes.Client{mock1, mock2, mock3}

	tests := []struct {
		name  string
		group string
		want  string // Expected client group
	}{
		{
			name:  "default group",
			group: "",
			want:  "default", // First client in default group
		},
		{
			name:  "specific group",
			group: "group1",
			want:  "group1",
		},
		{
			name:  "any group",
			group: "*",
			want:  "group1", // First client in any group
		},
		{
			name:  "non-existent group",
			group: "nonexistent",
			want:  "", // No client should be returned
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := pool.GetClient(spamoortypes.SelectClientByIndex, 0, tt.group)
			if tt.want != "" {
				assert.NotNil(t, client)
				assert.Equal(t, tt.want, client.GetClientGroup())
			} else {
				assert.Nil(t, client)
			}
		})
	}
}

// TestClientPoolClientPreparation tests client preparation with various scenarios
func TestClientPoolClientPreparation(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}
	pool := NewClientPool(ctx, []string{}, logger)

	// Test preparing clients with no RPC hosts
	err := pool.PrepareClients()
	assert.Error(t, err)
	assert.Equal(t, "no rpc hosts provided", err.Error())

	// Test preparing clients with valid host
	pool.rpcHosts = []string{"mock://localhost:8545"}
	pool.clientFactory = func(rpchost string) (spamoortypes.Client, error) {
		mock := testingutils.NewMockClient()
		return mock, nil
	}

	err = pool.PrepareClients()
	assert.NoError(t, err)
	assert.NotEmpty(t, pool.allClients)
	assert.NotNil(t, pool.chainId)

	// Test preparing clients with client error
	pool.allClients = nil // Reset clients
	pool.goodClients = nil
	pool.chainId = nil
	pool.clientFactory = func(rpchost string) (spamoortypes.Client, error) {
		return nil, errors.New("client error")
	}

	err = pool.PrepareClients()
	assert.Error(t, err)
	assert.Equal(t, "no useable clients", err.Error())
}

// TestClientPoolClientSelectionModes tests all client selection modes
func TestClientPoolClientSelectionModes(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}
	pool := NewClientPool(ctx, []string{"mock://localhost:8545"}, logger)

	// Create mock clients
	mock1 := testingutils.NewMockClient()
	mock1.SetClientGroup("default")
	mock1.SetMockNonce(1)
	mock2 := testingutils.NewMockClient()
	mock2.SetClientGroup("default")
	mock2.SetMockNonce(2)
	mock3 := testingutils.NewMockClient()
	mock3.SetClientGroup("default")
	mock3.SetMockNonce(3)

	clients := []*testingutils.MockClient{mock1, mock2, mock3}
	pool.allClients = []spamoortypes.Client{mock1, mock2, mock3}
	pool.goodClients = []spamoortypes.Client{mock1, mock2, mock3}

	// Test all selection modes
	modes := []spamoortypes.ClientSelectionMode{
		spamoortypes.SelectClientByIndex,
		spamoortypes.SelectClientRandom,
		spamoortypes.SelectClientRoundRobin,
	}

	for _, mode := range modes {
		client := pool.GetClient(mode, 0, "")
		assert.NotNil(t, client)
	}

	// Test round robin selection specifically
	startIdx := pool.rrClientIdx
	clientA := pool.GetClient(spamoortypes.SelectClientRoundRobin, 0, "").(*testingutils.MockClient)
	clientB := pool.GetClient(spamoortypes.SelectClientRoundRobin, 0, "").(*testingutils.MockClient)
	clientC := pool.GetClient(spamoortypes.SelectClientRoundRobin, 0, "").(*testingutils.MockClient)
	clientD := pool.GetClient(spamoortypes.SelectClientRoundRobin, 0, "").(*testingutils.MockClient)

	// Should cycle through all three, then back to the first
	expected := []uint64{
		clients[(startIdx+0)%3].GetNonce(),
		clients[(startIdx+1)%3].GetNonce(),
		clients[(startIdx+2)%3].GetNonce(),
		clients[(startIdx+3)%3].GetNonce(),
	}
	actual := []uint64{
		clientA.GetNonce(),
		clientB.GetNonce(),
		clientC.GetNonce(),
		clientD.GetNonce(),
	}
	assert.Equal(t, expected, actual)
}

// TestClientPoolClientStatusMonitoring tests the client status monitoring loop
func TestClientPoolClientStatusMonitoring(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	logger := &testingutils.MockLogger{}
	pool := NewClientPool(ctx, []string{"mock://localhost:8545"}, logger)

	// Create mock clients
	mock1 := testingutils.NewMockClient()
	mock1.SetMockBlockHeight(100)
	mock2 := testingutils.NewMockClient()
	mock2.SetMockBlockHeight(98)

	pool.allClients = []spamoortypes.Client{mock1, mock2}
	pool.goodClients = []spamoortypes.Client{mock1, mock2}

	// Start monitoring loop
	go pool.watchClientStatusLoop()

	// Wait for context cancellation
	<-ctx.Done()

	// Verify good clients list
	goodClients := pool.GetAllGoodClients()
	assert.Equal(t, 2, len(goodClients))
}

// TestClientPoolClientTimeout tests client timeout handling
func TestClientPoolClientTimeout(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}
	pool := NewClientPool(ctx, []string{"mock://localhost:8545"}, logger)

	// Create mock client with timeout
	mock := testingutils.NewMockClient()
	mock.SetTimeout(5 * time.Second)
	pool.allClients = []spamoortypes.Client{mock}
	pool.goodClients = []spamoortypes.Client{mock}

	// Test timeout handling
	err := pool.watchClientStatus()
	assert.NoError(t, err)
	assert.Equal(t, 5*time.Second, mock.GetTimeout())
}

// TestClientPoolClientVersion tests client version handling
func TestClientPoolClientVersion(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}
	pool := NewClientPool(ctx, []string{"mock://localhost:8545"}, logger)

	// Create mock client with version
	mock := testingutils.NewMockClient()
	pool.allClients = []spamoortypes.Client{mock}
	pool.goodClients = []spamoortypes.Client{mock}

	// Test version handling
	version, err := mock.GetClientVersion(ctx)
	assert.NoError(t, err)
	assert.Empty(t, version) // Mock client returns empty version by default
}

// TestClientPoolClientReceipts tests client receipt handling
func TestClientPoolClientReceipts(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}
	pool := NewClientPool(ctx, []string{"mock://localhost:8545"}, logger)

	// Create mock client with receipt
	mock := testingutils.NewMockClient()
	receipt := &types.Receipt{
		Status: 1,
	}
	mock.SetMockReceipt(receipt)
	pool.allClients = []spamoortypes.Client{mock}
	pool.goodClients = []spamoortypes.Client{mock}

	// Test receipt handling
	gotReceipt, err := mock.GetTransactionReceipt(ctx, common.Hash{})
	assert.NoError(t, err)
	assert.Equal(t, receipt, gotReceipt)
}

// TestClientPoolClientBlockReceipts tests client block receipts handling
func TestClientPoolClientBlockReceipts(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}
	pool := NewClientPool(ctx, []string{"mock://localhost:8545"}, logger)

	// Create mock client with block receipts
	mock := testingutils.NewMockClient()
	receipts := []*types.Receipt{
		{Status: 1},
		{Status: 1},
	}
	mock.SetMockBlockReceipts(receipts)
	pool.allClients = []spamoortypes.Client{mock}
	pool.goodClients = []spamoortypes.Client{mock}

	// Test block receipts handling
	gotReceipts, err := mock.GetBlockReceipts(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, receipts, gotReceipts)
}

// TestClientPoolClientGasLimit tests client gas limit handling
func TestClientPoolClientGasLimit(t *testing.T) {
	ctx := context.Background()
	logger := &testingutils.MockLogger{}
	pool := NewClientPool(ctx, []string{"mock://localhost:8545"}, logger)

	// Create mock client with gas limit
	mock := testingutils.NewMockClient()
	mock.SetMockGasLimit(30000000)
	pool.allClients = []spamoortypes.Client{mock}
	pool.goodClients = []spamoortypes.Client{mock}

	// Test gas limit handling
	gasLimit, err := mock.GetLatestGasLimit(ctx)
	assert.NoError(t, err)
	assert.Equal(t, uint64(30000000), gasLimit)
}
