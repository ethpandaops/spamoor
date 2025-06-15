package spamoor

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"

	"github.com/ethpandaops/spamoor/spamoortypes"
	"github.com/ethpandaops/spamoor/txbuilder"
)

// Wallet represents an Ethereum wallet with private key management, nonce tracking,
// and balance management. It provides thread-safe operations for transaction building,
// nonce management, and balance updates. The wallet automatically handles nonce
// sequencing and provides confirmation tracking for submitted transactions.
type Wallet struct {
	nonceMutex     sync.Mutex
	balanceMutex   sync.RWMutex
	privkey        *ecdsa.PrivateKey
	address        common.Address
	chainid        *big.Int
	pendingNonce   atomic.Uint64
	confirmedNonce uint64
	balance        *big.Int

	txNonceChans     map[uint64]*spamoortypes.NonceStatus
	txNonceMutex     sync.Mutex
	lastConfirmation uint64
}

// NewWallet creates a new wallet from a private key string.
// If privkey is empty, generates a new random private key.
// The privkey parameter accepts hex strings with or without "0x" prefix.
func NewWallet(privkey string) (spamoortypes.Wallet, error) {
	wallet := &Wallet{
		txNonceChans: map[uint64]*spamoortypes.NonceStatus{},
	}
	err := wallet.loadPrivateKey(privkey)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

// loadPrivateKey loads a private key from a hex string or generates a new random key.
// If privkey is empty, it generates a new random private key.
// If privkey is provided, it accepts hex strings with or without "0x" prefix.
// Also derives and sets the wallet's Ethereum address from the private key.
func (wallet *Wallet) loadPrivateKey(privkey string) error {
	var (
		privateKey *ecdsa.PrivateKey
		err        error
	)
	if privkey == "" {
		privateKey, err = crypto.GenerateKey()
	} else {
		if strings.HasPrefix(strings.ToLower(privkey), "0x") {
			privkey = privkey[2:]
		}
		privateKey, err = crypto.HexToECDSA(privkey)
	}
	if err != nil {
		return err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return errors.New("error casting public key to ECDSA")
	}

	wallet.privkey = privateKey
	wallet.address = crypto.PubkeyToAddress(*publicKeyECDSA)
	return nil
}

// GetAddress returns the Ethereum address associated with this wallet.
func (wallet *Wallet) GetAddress() common.Address {
	return wallet.address
}

// SetAddress updates the wallet's Ethereum address.
// This is typically only used for special cases or testing.
func (wallet *Wallet) SetAddress(address common.Address) {
	wallet.address = address
}

// GetPrivateKey returns the wallet's private key.
// Handle with care to avoid exposing sensitive data.
func (wallet *Wallet) GetPrivateKey() *ecdsa.PrivateKey {
	return wallet.privkey
}

// GetChainId returns the chain ID this wallet is configured for.
// Returns nil if no chain ID has been set.
func (wallet *Wallet) GetChainId() *big.Int {
	return wallet.chainid
}

// GetNonce returns the current pending nonce for this wallet.
// This nonce is the next nonce that should be used for the next transaction, but it is for informational purposes only.
// To actually use a nonce, you need to call GetNextNonce() which will increment the nonce and return the next nonce.
func (wallet *Wallet) GetNonce() uint64 {
	return wallet.pendingNonce.Load()
}

// GetConfirmedNonce returns the last confirmed nonce for this wallet.
// This represents the highest nonce that has been confirmed on-chain.
func (wallet *Wallet) GetConfirmedNonce() uint64 {
	return wallet.confirmedNonce
}

// SetConfirmedNonce updates the last confirmed nonce for this wallet.
// This is typically called when a transaction is confirmed to update the confirmed nonce.
func (wallet *Wallet) SetConfirmedNonce(nonce uint64) {
	wallet.confirmedNonce = nonce
}

// GetBalance returns the current balance of the wallet.
// The returned value is thread-safe to read.
func (wallet *Wallet) GetBalance() *big.Int {
	wallet.balanceMutex.RLock()
	defer wallet.balanceMutex.RUnlock()
	return wallet.balance
}

// SetChainId sets the chain ID for this wallet.
// This affects transaction signing and should match the target network.
func (wallet *Wallet) SetChainId(chainid *big.Int) {
	wallet.chainid = chainid
}

// SetNonce updates both the confirmed and pending nonce if the new nonce is higher.
// This is typically called when syncing wallet state with the blockchain.
func (wallet *Wallet) SetNonce(nonce uint64) {
	wallet.nonceMutex.Lock()
	defer wallet.nonceMutex.Unlock()

	pendingNonce := wallet.pendingNonce.Load()
	if nonce > pendingNonce {
		wallet.pendingNonce.Store(nonce)
	}

	wallet.confirmedNonce = nonce
}

// GetNextNonce atomically increments and returns the next available nonce.
// This is used when building transactions to ensure unique nonces.
func (wallet *Wallet) GetNextNonce() uint64 {
	wallet.nonceMutex.Lock()
	defer wallet.nonceMutex.Unlock()
	return wallet.pendingNonce.Add(1) - 1
}

// SetBalance sets the wallet's balance to the specified amount.
// This is typically called when syncing wallet state with the blockchain.
func (wallet *Wallet) SetBalance(balance *big.Int) {
	wallet.balanceMutex.Lock()
	defer wallet.balanceMutex.Unlock()
	wallet.balance = balance
}

// SubBalance subtracts the specified amount from the wallet's balance.
// This is typically called when a transaction is confirmed to update the balance.
func (wallet *Wallet) SubBalance(amount *big.Int) {
	wallet.balanceMutex.Lock()
	defer wallet.balanceMutex.Unlock()
	wallet.balance = wallet.balance.Sub(wallet.balance, amount)
}

// AddBalance adds the specified amount to the wallet's balance.
// This is typically called when the wallet receives funds.
func (wallet *Wallet) AddBalance(amount *big.Int) {
	wallet.balanceMutex.Lock()
	defer wallet.balanceMutex.Unlock()
	wallet.balance = wallet.balance.Add(wallet.balance, amount)
}

// BuildDynamicFeeTx builds and signs a dynamic fee (EIP-1559) transaction.
// It automatically assigns the next available nonce and signs the transaction.
func (wallet *Wallet) BuildDynamicFeeTx(txData *types.DynamicFeeTx) (*types.Transaction, error) {
	txData.ChainID = wallet.chainid
	txData.Nonce = wallet.GetNextNonce()
	return wallet.signTx(txData)
}

// BuildBlobTx builds and signs a blob transaction (EIP-4844).
// It automatically assigns the next available nonce and signs the transaction.
func (wallet *Wallet) BuildBlobTx(txData *types.BlobTx) (*types.Transaction, error) {
	txData.ChainID = uint256.MustFromBig(wallet.chainid)
	txData.Nonce = wallet.GetNextNonce()
	return wallet.signTx(txData)
}

// BuildSetCodeTx builds and signs a set code transaction (EIP-7702).
// It automatically assigns the next available nonce and signs the transaction.
func (wallet *Wallet) BuildSetCodeTx(txData *types.SetCodeTx) (*types.Transaction, error) {
	txData.ChainID = uint256.NewInt(wallet.chainid.Uint64())
	txData.Nonce = wallet.GetNextNonce()
	return wallet.signTx(txData)
}

// BuildBoundTx builds a transaction using the go-ethereum bind package.
// It sets up a TransactOpts with the wallet's credentials and calls the provided
// buildFn to construct the actual transaction. Useful for contract interactions.
func (wallet *Wallet) BuildBoundTx(ctx context.Context, txData *txbuilder.TxMetadata, buildFn func(transactOpts *bind.TransactOpts) (*types.Transaction, error)) (*types.Transaction, error) {
	transactor, err := bind.NewKeyedTransactorWithChainID(wallet.privkey, wallet.chainid)
	if err != nil {
		return nil, err
	}

	wallet.nonceMutex.Lock()
	defer wallet.nonceMutex.Unlock()

	transactor.Context = ctx
	transactor.From = wallet.address
	nonce := wallet.pendingNonce.Add(1) - 1
	transactor.Nonce = big.NewInt(0).SetUint64(nonce)

	transactor.GasTipCap = txData.GasTipCap.ToBig()
	transactor.GasFeeCap = txData.GasFeeCap.ToBig()
	transactor.GasLimit = txData.Gas
	transactor.Value = txData.Value.ToBig()
	transactor.NoSend = true

	tx, err := buildFn(transactor)
	if err != nil {
		wallet.pendingNonce.Store(nonce)
		return nil, err
	}

	return tx, nil
}

// ReplaceDynamicFeeTx builds a replacement dynamic fee transaction with a specific nonce.
// This is useful for replacing stuck transactions with higher gas prices.
func (wallet *Wallet) ReplaceDynamicFeeTx(txData *types.DynamicFeeTx, nonce uint64) (*types.Transaction, error) {
	txData.ChainID = wallet.chainid
	txData.Nonce = nonce
	return wallet.signTx(txData)
}

// ReplaceBlobTx builds a replacement blob transaction with a specific nonce.
// This is useful for replacing stuck blob transactions with higher gas prices.
func (wallet *Wallet) ReplaceBlobTx(txData *types.BlobTx, nonce uint64) (*types.Transaction, error) {
	txData.ChainID = uint256.MustFromBig(wallet.chainid)
	txData.Nonce = nonce
	return wallet.signTx(txData)
}

// ResetPendingNonce syncs the wallet's pending nonce with the blockchain.
// This is useful for recovering from nonce mismatches or wallet state corruption.
// It queries the pending nonce from the client and updates the wallet accordingly.
func (wallet *Wallet) ResetPendingNonce(ctx context.Context, client spamoortypes.Client) {
	wallet.nonceMutex.Lock()
	defer wallet.nonceMutex.Unlock()

	nonce, err := client.GetPendingNonceAt(ctx, wallet.address)
	if nonce < wallet.confirmedNonce {
		nonce = wallet.confirmedNonce
	}

	if err == nil && wallet.pendingNonce.Load() != nonce {
		logrus.Warnf("Resyncing pending nonce for %v from %d to %d", wallet.address.String(), wallet.pendingNonce.Load(), nonce)
		wallet.pendingNonce.Store(nonce)
	}
}

// signTx signs a transaction using the wallet's private key and chain ID.
// It creates a new transaction from the provided transaction data and signs it
// using the latest signer for the wallet's configured chain ID.
func (wallet *Wallet) signTx(txData types.TxData) (*types.Transaction, error) {
	tx := types.NewTx(txData)
	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(wallet.chainid), wallet.privkey)
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}

// GetTxNonceChan returns or creates a nonce status channel for tracking transaction confirmation.
// It manages a map of nonce channels used to wait for specific transaction confirmations.
// Returns the nonce status and a boolean indicating if this is the first pending transaction.
// If the target nonce is already confirmed, returns nil and false.
func (wallet *Wallet) GetTxNonceChan(targetNonce uint64) (*spamoortypes.NonceStatus, bool) {
	wallet.txNonceMutex.Lock()
	defer wallet.txNonceMutex.Unlock()

	if wallet.confirmedNonce > targetNonce {
		return nil, false
	}

	nonceChan := wallet.txNonceChans[targetNonce]
	if nonceChan != nil {
		return nonceChan, false
	}

	nonceChan = &spamoortypes.NonceStatus{
		Channel: make(chan bool),
	}
	wallet.txNonceChans[targetNonce] = nonceChan

	return nonceChan, len(wallet.txNonceChans) == 1
}

func (wallet *Wallet) GetLastConfirmation() uint64 {
	return wallet.lastConfirmation
}

func (wallet *Wallet) SetLastConfirmation(nonce uint64) {
	wallet.lastConfirmation = nonce
}

// processTransactionInclusion handles the confirmation of a transaction from a tracked wallet.
// It updates the wallet's nonce state, signals any waiting confirmation channels,
// and cleans up completed nonce channels. Updates the wallet's confirmation tracking.
func (wallet *Wallet) ProcessTransactionInclusion(blockNumber uint64, tx *types.Transaction, receipt *types.Receipt) {
	nonce := tx.Nonce()

	wallet.txNonceMutex.Lock()
	defer wallet.txNonceMutex.Unlock()

	if nonceChan := wallet.txNonceChans[nonce]; nonceChan != nil {
		nonceChan.Receipt = receipt
	}

	wallet.confirmedNonce = nonce + 1
	if nonce+1 > wallet.pendingNonce.Load() {
		wallet.pendingNonce.Store(nonce + 1)
	}
	if blockNumber > wallet.lastConfirmation {
		wallet.lastConfirmation = blockNumber
	}

	for n := range wallet.txNonceChans {
		if n <= nonce {
			close(wallet.txNonceChans[n].Channel)
			delete(wallet.txNonceChans, n)
		}
	}
}

// RevertTransactionReceival subtracts the value of the transaction from the wallet's balance.
// This is used to revert the balance when a transaction is reverted.
func (wallet *Wallet) RevertTransactionReceival(tx *types.Transaction) {
	wallet.balance = wallet.balance.Sub(wallet.balance, tx.Value())
}

// ProcessTransactionReceival adds the value of the transaction to the wallet's balance.
// This is used to update the balance when a transaction is received.
func (wallet *Wallet) ProcessTransactionReceival(tx *types.Transaction) {
	wallet.balance = wallet.balance.Add(wallet.balance, tx.Value())
}

// GetPendingTxCount returns the number of pending transactions for this wallet.
func (wallet *Wallet) GetPendingTxCount() int {
	return len(wallet.txNonceChans)
}

// ProcessStaleTransactions recovers stale transactions that may have been missed.
// It closes the confirmation channels for stale transactions and deletes them from the nonce channel map.
func (wallet *Wallet) ProcessStaleTransactions(blockNumber uint64, nonce uint64) {
	pendingNonce := 0
	for n := range wallet.txNonceChans {
		pendingNonce = int(n)
		break
	}

	logrus.Debugf("recovering %v stale transactions for %v (current nonce %v, cache nonce %v, first pending nonce: %v)", len(wallet.txNonceChans), wallet.address.String(), nonce, wallet.confirmedNonce, pendingNonce)

	wallet.txNonceMutex.Lock()
	defer wallet.txNonceMutex.Unlock()

	for n := range wallet.txNonceChans {
		if n < nonce {
			logrus.Warnf("recovering stale confirmed transactions for %v (nonce %v)", wallet.address.String(), n)
			close(wallet.txNonceChans[n].Channel)
			delete(wallet.txNonceChans, n)
		}
	}
}
