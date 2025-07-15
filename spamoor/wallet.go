package spamoor

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
)

// Wallet represents an Ethereum wallet with private key management, nonce tracking,
// and balance management. It provides thread-safe operations for transaction building,
// nonce management, and balance updates. The wallet automatically handles nonce
// sequencing and provides confirmation tracking for submitted transactions.
type Wallet struct {
	nonceMutex       sync.Mutex
	balanceMutex     sync.RWMutex
	privkey          *ecdsa.PrivateKey
	address          common.Address
	chainid          *big.Int
	pendingTxCount   atomic.Uint64
	submittedTxCount atomic.Uint64
	confirmedTxCount uint64
	balance          *big.Int

	txNonceChans     map[uint64]*nonceStatus
	txNonceMutex     sync.Mutex
	lastConfirmation uint64

	lowBalanceNotifyChan chan<- struct{}
	lowBalanceThreshold  *big.Int
}

// nonceStatus tracks the confirmation status of a transaction with a specific nonce
type nonceStatus struct {
	txs     []*PendingTx
	receipt *types.Receipt
	channel chan bool
}

type PendingTx struct {
	Tx               *types.Transaction
	Submitted        time.Time
	LastRebroadcast  time.Time
	RebroadcastCount uint64
}

// NewWallet creates a new wallet from a private key string.
// If privkey is empty, generates a new random private key.
// The privkey parameter accepts hex strings with or without "0x" prefix.
func NewWallet(privkey string) (*Wallet, error) {
	wallet := &Wallet{
		txNonceChans: map[uint64]*nonceStatus{},
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
	return wallet.pendingTxCount.Load()
}

// GetConfirmedNonce returns the last confirmed nonce for this wallet.
// This represents the highest nonce that has been confirmed on-chain.
func (wallet *Wallet) GetConfirmedNonce() uint64 {
	return wallet.confirmedTxCount
}

// GetSubmittedTxCount returns the total number of transactions submitted by this wallet.
// This represents the cumulative count of all transactions that have been submitted to the network.
func (wallet *Wallet) GetSubmittedTxCount() uint64 {
	return wallet.submittedTxCount.Load()
}

// IncrementSubmittedTxCount increments the submitted transaction counter.
// This should be called when a transaction is successfully submitted to the network.
func (wallet *Wallet) IncrementSubmittedTxCount() {
	wallet.submittedTxCount.Add(1)
}

// GetBalance returns the current balance of the wallet.
// The returned value is thread-safe to read.
func (wallet *Wallet) GetBalance() *big.Int {
	wallet.balanceMutex.RLock()
	defer wallet.balanceMutex.RUnlock()
	return wallet.balance
}

// setLowBalanceNotification sets up low balance notification for this wallet.
// When the balance falls below the threshold, the wallet will send a notification
// to the channel (non-blocking).
func (wallet *Wallet) setLowBalanceNotification(notifyChan chan<- struct{}, threshold *big.Int) {
	wallet.balanceMutex.Lock()
	defer wallet.balanceMutex.Unlock()
	wallet.lowBalanceNotifyChan = notifyChan
	wallet.lowBalanceThreshold = threshold
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

	pendingNonce := wallet.pendingTxCount.Load()
	if nonce > pendingNonce {
		wallet.pendingTxCount.Store(nonce)
	}

	wallet.confirmedTxCount = nonce
}

// GetNextNonce atomically increments and returns the next available nonce.
// This is used when building transactions to ensure unique nonces.
func (wallet *Wallet) GetNextNonce() uint64 {
	wallet.nonceMutex.Lock()
	defer wallet.nonceMutex.Unlock()
	return wallet.pendingTxCount.Add(1) - 1
}

// SetBalance sets the wallet's balance to the specified amount.
// This is typically called when syncing wallet state with the blockchain.
func (wallet *Wallet) SetBalance(balance *big.Int) {
	wallet.balanceMutex.Lock()
	defer wallet.balanceMutex.Unlock()
	wallet.balance = balance
	wallet.checkLowBalance()
}

// SubBalance subtracts the specified amount from the wallet's balance.
// This is typically called when a transaction is confirmed to update the balance.
func (wallet *Wallet) SubBalance(amount *big.Int) {
	wallet.balanceMutex.Lock()
	defer wallet.balanceMutex.Unlock()
	wallet.balance = wallet.balance.Sub(wallet.balance, amount)
	wallet.checkLowBalance()
}

// AddBalance adds the specified amount to the wallet's balance.
// This is typically called when the wallet receives funds.
func (wallet *Wallet) AddBalance(amount *big.Int) {
	wallet.balanceMutex.Lock()
	defer wallet.balanceMutex.Unlock()
	wallet.balance = wallet.balance.Add(wallet.balance, amount)
}

// checkLowBalance checks if balance has fallen below threshold and sends notification.
// Must be called with balanceMutex held. Non-blocking - drops notification if channel is full.
func (wallet *Wallet) checkLowBalance() {
	if wallet.lowBalanceNotifyChan == nil || wallet.lowBalanceThreshold == nil || wallet.balance == nil {
		return
	}

	if wallet.balance.Cmp(wallet.lowBalanceThreshold) < 0 {
		// Non-blocking send
		select {
		case wallet.lowBalanceNotifyChan <- struct{}{}:
		default:
			// Channel full, drop notification
		}
	}
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
	nonce := wallet.pendingTxCount.Add(1) - 1
	transactor.Nonce = big.NewInt(0).SetUint64(nonce)

	transactor.GasTipCap = txData.GasTipCap.ToBig()
	transactor.GasFeeCap = txData.GasFeeCap.ToBig()
	transactor.GasLimit = txData.Gas
	transactor.Value = txData.Value.ToBig()
	transactor.NoSend = true

	tx, err := buildFn(transactor)
	if err != nil {
		wallet.pendingTxCount.Store(nonce)
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
func (wallet *Wallet) ResetPendingNonce(ctx context.Context, client *Client) {
	wallet.nonceMutex.Lock()
	defer wallet.nonceMutex.Unlock()

	nonce, err := client.GetPendingNonceAt(ctx, wallet.address)
	if err == nil && nonce < wallet.confirmedTxCount {
		logrus.Errorf("Resyncing confirmed nonce for %v from %d to %d (this should never happen)", wallet.address.String(), wallet.confirmedTxCount, nonce)
		wallet.confirmedTxCount = nonce
	}

	if err == nil && wallet.pendingTxCount.Load() != nonce {
		logrus.Warnf("Resyncing pending nonce for %v from %d to %d", wallet.address.String(), wallet.pendingTxCount.Load(), nonce)
		wallet.pendingTxCount.Store(nonce)
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

// getTxNonceChan returns or creates a nonce status channel for tracking transaction confirmation.
// It manages a map of nonce channels used to wait for specific transaction confirmations.
// Returns the nonce status and a boolean indicating if this is the first pending transaction.
// If the target nonce is already confirmed, returns nil and false.
func (wallet *Wallet) getTxNonceChan(tx *types.Transaction) (*nonceStatus, bool) {
	wallet.txNonceMutex.Lock()
	defer wallet.txNonceMutex.Unlock()

	targetNonce := tx.Nonce()

	if wallet.confirmedTxCount > targetNonce {
		return nil, false
	}

	nonceChan := wallet.txNonceChans[targetNonce]
	if nonceChan != nil {
		for _, existingTx := range nonceChan.txs {
			if existingTx.Tx.Hash() == tx.Hash() {
				return nonceChan, false
			}
		}

		nonceChan.txs = append(nonceChan.txs, &PendingTx{
			Tx:        tx,
			Submitted: time.Now(),
		})
		return nonceChan, false
	}

	nonceChan = &nonceStatus{
		txs: []*PendingTx{
			{
				Tx:        tx,
				Submitted: time.Now(),
			},
		},
		channel: make(chan bool),
	}
	wallet.txNonceChans[targetNonce] = nonceChan

	return nonceChan, len(wallet.txNonceChans) == 1
}

func (wallet *Wallet) GetPendingTx(tx *types.Transaction) *PendingTx {
	wallet.txNonceMutex.Lock()
	defer wallet.txNonceMutex.Unlock()

	nonceChan := wallet.txNonceChans[tx.Nonce()]
	if nonceChan == nil {
		return nil
	}

	for _, pendingTx := range nonceChan.txs {
		if pendingTx.Tx.Hash() == tx.Hash() {
			return pendingTx
		}
	}

	return nil
}

func (wallet *Wallet) GetPendingTxs() []*PendingTx {
	wallet.txNonceMutex.Lock()
	defer wallet.txNonceMutex.Unlock()

	pendingNonces := []uint64{}
	for nonce := range wallet.txNonceChans {
		pendingNonces = append(pendingNonces, nonce)
	}

	sort.Slice(pendingNonces, func(i, j int) bool {
		return pendingNonces[i] < pendingNonces[j]
	})

	pendingTxs := []*PendingTx{}
	for _, nonce := range pendingNonces {
		pendingTxs = append(pendingTxs, wallet.txNonceChans[nonce].txs...)
	}

	return pendingTxs
}
