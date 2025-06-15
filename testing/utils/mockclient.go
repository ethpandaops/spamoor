package testingutils

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"

	"github.com/ethpandaops/spamoor/spamoortypes"
)

// MockClient is a mock implementation of Client for testing
type MockClient struct {
	mu          sync.RWMutex
	name        string
	clientGroup string
	enabled     bool
	timeout     time.Duration

	// Mock data
	chainId         *big.Int
	nonce           uint64
	balance         *big.Int
	blockHeight     uint64
	blockHeightTime time.Time
	gasCap          *big.Int
	tipCap          *big.Int
	clientVersion   string
	receipt         *types.Receipt
	block           *types.Block
	blockReceipts   []*types.Receipt
	gasLimit        uint64

	// Mock errors
	err error
}

// NewMockClient creates a new mock client with default values
func NewMockClient() *MockClient {
	return &MockClient{
		name:        "mock",
		clientGroup: "default",
		enabled:     true,
		chainId:     big.NewInt(1337),
		nonce:       0,
		balance:     big.NewInt(0),
		blockHeight: 0,
		gasCap:      big.NewInt(1000000000),
		tipCap:      big.NewInt(1000000000),
		gasLimit:    30000000,
	}
}

// GetName returns the mock client name
func (m *MockClient) GetName() string {
	return m.name
}

// GetClientGroup returns the mock client group
func (m *MockClient) GetClientGroup() string {
	return m.clientGroup
}

// GetEthClient returns a mock ethclient
func (m *MockClient) GetEthClient() bind.ContractBackend {
	return nil // Mock implementation
}

// GetRPCHost returns a mock RPC host
func (m *MockClient) GetRPCHost() string {
	return "mock://localhost:8545"
}

// GetTimeout returns the mock client timeout
func (m *MockClient) GetTimeout() time.Duration {
	return m.timeout
}

// SetTimeout sets the mock client timeout
func (m *MockClient) SetTimeout(timeout time.Duration) {
	m.timeout = timeout
}

// UpdateWallet updates the mock wallet with mock data
func (m *MockClient) UpdateWallet(ctx context.Context, wallet spamoortypes.Wallet) error {
	if m.err != nil {
		return m.err
	}
	if wallet.GetChainId() == nil {
		wallet.SetChainId(m.chainId)
	}
	wallet.SetNonce(m.nonce)
	wallet.SetBalance(m.balance)
	return nil
}

// SetClientGroup sets the mock client group
func (m *MockClient) SetClientGroup(group string) {
	m.clientGroup = group
}

// IsEnabled returns whether the mock client is enabled
func (m *MockClient) IsEnabled() bool {
	return m.enabled
}

// SetEnabled sets the mock client enabled state
func (m *MockClient) SetEnabled(enabled bool) {
	m.enabled = enabled
}

// GetChainId returns the mock chain ID
func (m *MockClient) GetChainId(ctx context.Context) (*big.Int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.err != nil {
		return nil, m.err
	}
	return m.chainId, nil
}

// GetNonce returns the mock nonce
func (m *MockClient) GetNonce() uint64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.nonce
}

// GetNonceAt returns the mock nonce
func (m *MockClient) GetNonceAt(ctx context.Context, wallet common.Address, blockNumber *big.Int) (uint64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.err != nil {
		return 0, m.err
	}
	return m.nonce, nil
}

// GetPendingNonceAt returns the mock pending nonce
func (m *MockClient) GetPendingNonceAt(ctx context.Context, wallet common.Address) (uint64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.err != nil {
		return 0, m.err
	}
	return m.nonce, nil
}

// GetBalanceAt returns the mock balance
func (m *MockClient) GetBalanceAt(ctx context.Context, wallet common.Address) (*big.Int, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.balance, nil
}

// GetSuggestedFee returns mock gas price and tip cap
func (m *MockClient) GetSuggestedFee(ctx context.Context) (*big.Int, *big.Int, error) {
	if m.err != nil {
		return nil, nil, m.err
	}
	return m.gasCap, m.tipCap, nil
}

// SendTransaction is a mock implementation
func (m *MockClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.err
}

// SendRawTransaction is a mock implementation
func (m *MockClient) SendRawTransaction(ctx context.Context, tx []byte) error {
	return m.err
}

// GetTransactionReceipt returns the mock receipt
func (m *MockClient) GetTransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.receipt, nil
}

// GetBlockHeight returns the mock block height
func (m *MockClient) GetBlockHeight(ctx context.Context) (uint64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.err != nil {
		return 0, m.err
	}
	return m.blockHeight, nil
}

// GetLastBlockHeight returns the mock last block height and time
func (m *MockClient) GetLastBlockHeight() (uint64, time.Time) {
	return m.blockHeight, m.blockHeightTime
}

// GetClientVersion returns the mock client version
func (m *MockClient) GetClientVersion(ctx context.Context) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.clientVersion, nil
}

// GetBlock returns the mock block
func (m *MockClient) GetBlock(ctx context.Context, blockNumber uint64) (*types.Block, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.block, nil
}

// GetBlockReceipts returns the mock block receipts
func (m *MockClient) GetBlockReceipts(ctx context.Context, blockNumber uint64) ([]*types.Receipt, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.blockReceipts, nil
}

// GetLatestGasLimit returns the mock gas limit
func (m *MockClient) GetLatestGasLimit(ctx context.Context) (uint64, error) {
	if m.err != nil {
		return 0, m.err
	}
	return m.gasLimit, nil
}

// SetMockError sets an error to be returned by mock methods
func (m *MockClient) SetMockError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.err = err
}

// SetMockChainId sets the mock chain ID
func (m *MockClient) SetMockChainId(chainId *big.Int) {
	m.chainId = chainId
}

// SetMockNonce sets the mock nonce
func (m *MockClient) SetMockNonce(nonce uint64) {
	m.nonce = nonce
}

// SetMockBalance sets the mock balance
func (m *MockClient) SetMockBalance(balance *big.Int) {
	m.balance = balance
}

// SetMockBlockHeight sets the mock block height
func (m *MockClient) SetMockBlockHeight(height uint64) {
	m.blockHeight = height
	m.blockHeightTime = time.Now()
}

// SetMockGasFees sets the mock gas fees
func (m *MockClient) SetMockGasFees(gasCap, tipCap *big.Int) {
	m.gasCap = gasCap
	m.tipCap = tipCap
}

// SetMockReceipt sets the mock receipt
func (m *MockClient) SetMockReceipt(receipt *types.Receipt) {
	m.receipt = receipt
}

// SetMockBlock sets the mock block
func (m *MockClient) SetMockBlock(block *types.Block) {
	m.block = block
}

// SetMockBlockReceipts sets the mock block receipts
func (m *MockClient) SetMockBlockReceipts(receipts []*types.Receipt) {
	m.blockReceipts = receipts
}

// SetMockGasLimit sets the mock gas limit
func (m *MockClient) SetMockGasLimit(gasLimit uint64) {
	m.gasLimit = gasLimit
}

// SetMockClientVersion sets the mock client version
func (m *MockClient) SetMockClientVersion(version string) {
	m.clientVersion = version
}

// MockWalletPool is a mock implementation of WalletPool for testing
type MockWalletPool struct {
	wallets map[common.Address]spamoortypes.Wallet
}

// NewMockWalletPool creates a new mock wallet pool
func NewMockWalletPool() *MockWalletPool {
	return &MockWalletPool{
		wallets: make(map[common.Address]spamoortypes.Wallet),
	}
}

// SetWallets sets the wallets in the mock pool
func (m *MockWalletPool) SetWallets(wallets map[common.Address]spamoortypes.Wallet) {
	m.wallets = wallets
}

// CollectPoolWallets implements the WalletPool interface
func (m *MockWalletPool) CollectPoolWallets(walletMap map[common.Address]spamoortypes.Wallet) {
	for addr, wallet := range m.wallets {
		walletMap[addr] = wallet
	}
}

// GetChainId implements the WalletPool interface
func (m *MockWalletPool) GetChainId() *big.Int {
	return big.NewInt(1)
}

// Minimal WalletPool interface implementation - just the methods needed for testing
func (m *MockWalletPool) GetContext() context.Context            { return context.Background() }
func (m *MockWalletPool) GetTxPool() spamoortypes.TxPool         { return nil }
func (m *MockWalletPool) GetClientPool() spamoortypes.ClientPool { return nil }
func (m *MockWalletPool) GetRootWallet() spamoortypes.Wallet     { return nil }
func (m *MockWalletPool) WithRootWalletLock(ctx context.Context, txCount int, lockedLogFn func(), lockedFn func() error) error {
	return nil
}
func (m *MockWalletPool) LoadConfig(configYaml string) error                            { return nil }
func (m *MockWalletPool) MarshalConfig() (string, error)                                { return "", nil }
func (m *MockWalletPool) SetWalletCount(count uint64)                                   {}
func (m *MockWalletPool) SetRunFundings(runFundings bool)                               {}
func (m *MockWalletPool) AddWellKnownWallet(config *spamoortypes.WellKnownWalletConfig) {}
func (m *MockWalletPool) SetRefillAmount(amount *uint256.Int)                           {}
func (m *MockWalletPool) SetRefillBalance(balance *uint256.Int)                         {}
func (m *MockWalletPool) SetWalletSeed(seed string)                                     {}
func (m *MockWalletPool) SetRefillInterval(interval uint64)                             {}
func (m *MockWalletPool) SetTransactionTracker(tracker func(err error))                 {}
func (m *MockWalletPool) GetTransactionTracker() func(err error)                        { return nil }
func (m *MockWalletPool) GetClient(mode spamoortypes.ClientSelectionMode, input int, group string) spamoortypes.Client {
	return nil
}
func (m *MockWalletPool) GetWallet(mode spamoortypes.WalletSelectionMode, input int) spamoortypes.Wallet {
	return nil
}
func (m *MockWalletPool) GetWellKnownWallet(name string) spamoortypes.Wallet { return nil }
func (m *MockWalletPool) GetVeryWellKnownWalletAddress(name string) common.Address {
	return common.Address{}
}
func (m *MockWalletPool) GetWalletName(address common.Address) string                   { return "" }
func (m *MockWalletPool) GetAllWallets() []spamoortypes.Wallet                          { return nil }
func (m *MockWalletPool) GetConfiguredWalletCount() uint64                              { return 0 }
func (m *MockWalletPool) GetWalletCount() uint64                                        { return 0 }
func (m *MockWalletPool) PrepareWallets() error                                         { return nil }
func (m *MockWalletPool) CheckChildWalletBalance(childWallet spamoortypes.Wallet) error { return nil }
func (m *MockWalletPool) ReclaimFunds(ctx context.Context, client spamoortypes.Client) error {
	return nil
}
