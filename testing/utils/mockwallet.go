package testingutils

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"hash/fnv"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"

	"github.com/ethpandaops/spamoor/spamoortypes"
	"github.com/ethpandaops/spamoor/txbuilder"
)

// MockWallet is a mock implementation of Wallet for testing
type MockWallet struct {
	mu sync.RWMutex
	
	// Wallet identity
	address    common.Address
	privateKey *ecdsa.PrivateKey
	
	// Chain configuration
	chainId *big.Int
	
	// Transaction state
	nonce   uint64
	balance *big.Int
	
	// Mock errors
	err error
}

// NewMockWallet creates a new mock wallet with a generated private key
func NewMockWallet() spamoortypes.Wallet {
	privateKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		// Fallback to a deterministic key for testing
		privateKey, _ = crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	}
	
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	
	return &MockWallet{
		address:    address,
		privateKey: privateKey,
		chainId:    big.NewInt(1),
		nonce:      0,
		balance:    big.NewInt(0),
	}
}

// NewMockWalletWithKey creates a mock wallet with a specific private key
func NewMockWalletWithKey(privateKeyHex string) spamoortypes.Wallet {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		// Fallback to default key
		privateKey, _ = crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	}
	
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	
	return &MockWallet{
		address:    address,
		privateKey: privateKey,
		chainId:    big.NewInt(1),
		nonce:      0,
		balance:    big.NewInt(0),
	}
}

// GenerateAddress generates a deterministic address from a seed string
func GenerateAddress(seed string) common.Address {
	hasher := fnv.New64a()
	hasher.Write([]byte(seed))
	hash := hasher.Sum64()
	
	// Create a 20-byte address from the hash
	addr := common.Address{}
	for i := 0; i < 20; i++ {
		addr[i] = byte(hash >> (i * 8))
	}
	
	return addr
}

// Wallet interface implementation for MockWallet

func (w *MockWallet) GetAddress() common.Address {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.address
}

func (w *MockWallet) SetAddress(address common.Address) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.address = address
}

func (w *MockWallet) GetPrivateKey() *ecdsa.PrivateKey {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.privateKey
}

func (w *MockWallet) GetChainId() *big.Int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.chainId
}

func (w *MockWallet) SetChainId(chainId *big.Int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.chainId = chainId
}

func (w *MockWallet) GetNonce() uint64 {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.nonce
}

func (w *MockWallet) SetNonce(nonce uint64) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.nonce = nonce
}

func (w *MockWallet) GetNextNonce() uint64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.nonce++
	return w.nonce - 1
}

func (w *MockWallet) GetBalance() *big.Int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.balance
}

func (w *MockWallet) SetBalance(balance *big.Int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.balance = balance
}

func (w *MockWallet) AddBalance(amount *big.Int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.balance = new(big.Int).Add(w.balance, amount)
}

func (w *MockWallet) SubBalance(amount *big.Int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.balance = new(big.Int).Sub(w.balance, amount)
}

// Transaction building methods - mock implementations
func (w *MockWallet) BuildDynamicFeeTx(txData *types.DynamicFeeTx) (*types.Transaction, error) {
	if w.err != nil {
		return nil, w.err
	}
	
	// Set chain ID and nonce like the real wallet
	txData.ChainID = w.GetChainId()
	txData.Nonce = w.GetNextNonce()
	
	return w.signTx(txData)
}

func (w *MockWallet) BuildBlobTx(txData *types.BlobTx) (*types.Transaction, error) {
	if w.err != nil {
		return nil, w.err
	}
	
	// Set chain ID and nonce like the real wallet
	txData.ChainID = uint256.MustFromBig(w.GetChainId())
	txData.Nonce = w.GetNextNonce()
	
	return w.signTx(txData)
}

func (w *MockWallet) BuildSetCodeTx(txData *types.SetCodeTx) (*types.Transaction, error) {
	if w.err != nil {
		return nil, w.err
	}
	
	// Set chain ID and nonce like the real wallet
	txData.ChainID = uint256.NewInt(w.GetChainId().Uint64())
	txData.Nonce = w.GetNextNonce()
	
	return w.signTx(txData)
}

// signTx signs a transaction like the real wallet
func (w *MockWallet) signTx(txData types.TxData) (*types.Transaction, error) {
	tx := types.NewTx(txData)
	signer := types.LatestSignerForChainID(w.GetChainId())
	signedTx, err := types.SignTx(tx, signer, w.GetPrivateKey())
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}

func (w *MockWallet) BuildBoundTx(ctx context.Context, txData *txbuilder.TxMetadata, buildFn func(transactOpts *bind.TransactOpts) (*types.Transaction, error)) (*types.Transaction, error) {
	if w.err != nil {
		return nil, w.err
	}
	
	// Create mock transact opts
	transactOpts := &bind.TransactOpts{
		From:     w.GetAddress(),
		Nonce:    big.NewInt(int64(w.GetNonce())),
		Signer:   nil, // Mock implementation
		Value:    txData.Value.ToBig(),
		GasPrice: nil,
		GasFeeCap: txData.GasFeeCap.ToBig(),
		GasTipCap: txData.GasTipCap.ToBig(),
		GasLimit: txData.Gas,
		Context:  ctx,
	}
	
	return buildFn(transactOpts)
}

// Mock error setting
func (w *MockWallet) SetMockError(err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.err = err
}

// Additional Wallet interface methods
func (w *MockWallet) GetConfirmedNonce() uint64 {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.nonce
}

func (w *MockWallet) SetConfirmedNonce(nonce uint64) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.nonce = nonce
}

func (w *MockWallet) ResetPendingNonce(ctx context.Context, client spamoortypes.Client) {
	// Mock implementation - no-op
}

func (w *MockWallet) ReplaceDynamicFeeTx(txData *types.DynamicFeeTx, nonce uint64) (*types.Transaction, error) {
	// Set the specific nonce for replacement
	originalNonce := txData.Nonce
	txData.Nonce = nonce
	
	result, err := w.BuildDynamicFeeTx(txData)
	
	// Restore original nonce
	txData.Nonce = originalNonce
	
	return result, err
}

func (w *MockWallet) ReplaceBlobTx(txData *types.BlobTx, nonce uint64) (*types.Transaction, error) {
	// Set the specific nonce for replacement
	originalNonce := txData.Nonce
	txData.Nonce = nonce
	
	result, err := w.BuildBlobTx(txData)
	
	// Restore original nonce
	txData.Nonce = originalNonce
	
	return result, err
}

func (w *MockWallet) ProcessTransactionInclusion(blockNumber uint64, tx *types.Transaction, receipt *types.Receipt) {
	// Mock implementation - update nonce if this is our transaction
	if receipt != nil && receipt.Status == types.ReceiptStatusSuccessful {
		w.mu.Lock()
		if tx.Nonce() >= w.nonce {
			w.nonce = tx.Nonce() + 1
		}
		w.mu.Unlock()
	}
}

func (w *MockWallet) RevertTransactionReceival(tx *types.Transaction) {
	// Mock implementation - no-op
}

func (w *MockWallet) ProcessTransactionReceival(tx *types.Transaction) {
	// Mock implementation - no-op
}

func (w *MockWallet) ProcessStaleTransactions(blockNumber uint64, nonce uint64) {
	// Mock implementation - no-op
}

func (w *MockWallet) GetTxNonceChan(targetNonce uint64) (*spamoortypes.NonceStatus, bool) {
	// Mock implementation - return false (not found)
	return nil, false
}

func (w *MockWallet) GetLastConfirmation() uint64 {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.nonce
}

func (w *MockWallet) SetLastConfirmation(nonce uint64) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if nonce > w.nonce {
		w.nonce = nonce
	}
}

func (w *MockWallet) GetPendingTxCount() int {
	// Mock implementation - return 0
	return 0
}