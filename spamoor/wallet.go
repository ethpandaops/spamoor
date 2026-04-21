package spamoor

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
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
	needNonceResync  bool
	isSynced         bool
	chainid          *big.Int
	pendingTxCount   atomic.Uint64
	submittedTxCount atomic.Uint64
	confirmedTxCount uint64
	skippedNonces    []uint64
	balance          *big.Int
	protected        bool // When true, GetPrivateKey() returns nil to prevent key extraction

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
	Options          *SendTransactionOptions
}

// NewWallet creates a new wallet from a private key string.
// If privkey is empty, generates a new random private key.
// The privkey parameter accepts hex strings with or without "0x" prefix.
func NewWallet(privkey *ecdsa.PrivateKey, address common.Address) *Wallet {
	wallet := &Wallet{
		privkey:      privkey,
		address:      address,
		txNonceChans: map[uint64]*nonceStatus{},
	}

	return wallet
}

// loadPrivateKey loads a private key from a hex string or generates a new random key.
// If privkey is empty, it generates a new random private key.
// If privkey is provided, it accepts hex strings with or without "0x" prefix.
// Also derives and sets the wallet's Ethereum address from the private key.
func LoadPrivateKey(privkey string) (*ecdsa.PrivateKey, common.Address, error) {
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
		return nil, common.Address{}, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, common.Address{}, errors.New("error casting public key to ECDSA")
	}

	return privateKey, crypto.PubkeyToAddress(*publicKeyECDSA), nil
}

func (wallet *Wallet) UpdateWallet(ctx context.Context, client *Client, refresh bool) error {
	if wallet.isSynced && !refresh {
		return nil
	}

	if wallet.GetChainId() == nil {
		chainId, err := client.GetChainId(ctx)
		if err != nil {
			return err
		}
		wallet.SetChainId(chainId)
	}

	nonce, err := client.GetNonceAt(ctx, wallet.GetAddress(), nil)
	if err != nil {
		return err
	}
	wallet.SetNonce(nonce)

	balance, err := client.GetBalanceAt(ctx, wallet.GetAddress())
	if err != nil {
		return err
	}
	wallet.SetBalance(balance)

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
// Returns nil if the wallet is protected (e.g. the root wallet).
// Handle with care to avoid exposing sensitive data.
func (wallet *Wallet) GetPrivateKey() *ecdsa.PrivateKey {
	if wallet.protected {
		return nil
	}
	return wallet.privkey
}

// setProtected sets the wallet's protection flag.
// When protected, GetPrivateKey() returns nil to prevent key extraction.
func (wallet *Wallet) setProtected(protected bool) {
	wallet.protected = protected
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
	if wallet.balance == nil {
		return new(big.Int)
	}
	return wallet.balance
}

func (wallet *Wallet) GetReadableBalance(unitDigits, maxPreCommaDigitsBeforeTrim, digits int, addPositiveSign, trimAmount bool) string {
	// Initialize trimmedAmount and postComma variables to "0"
	fullAmount := ""
	trimmedAmount := "0"
	postComma := "0"
	proceed := ""
	amount := wallet.GetBalance()

	if amount != nil {
		s := amount.String()

		if amount.Sign() > 0 && addPositiveSign {
			proceed = "+"
		} else if amount.Sign() < 0 {
			proceed = "-"
			s = strings.Replace(s, "-", "", 1)
		}

		l := len(s)

		// Check if there is a part of the amount before the decimal point
		switch {
		case l > unitDigits:
			// Calculate length of preComma part
			l -= unitDigits
			// Set preComma to part of the string before the decimal point
			trimmedAmount = s[:l]
			// Set postComma to part of the string after the decimal point, after removing trailing zeros
			postComma = strings.TrimRight(s[l:], "0")

			// Check if the preComma part exceeds the maximum number of digits before the decimal point
			if maxPreCommaDigitsBeforeTrim > 0 && l > maxPreCommaDigitsBeforeTrim {
				// Reduce the number of digits after the decimal point by the excess number of digits in the preComma part
				l -= maxPreCommaDigitsBeforeTrim
				if digits < l {
					digits = 0
				} else {
					digits -= l
				}
			}
			// Check if there is only a part of the amount after the decimal point, and no leading zeros need to be added
		case l == unitDigits:
			// Set postComma to part of the string after the decimal point, after removing trailing zeros
			postComma = strings.TrimRight(s, "0")
			// Check if there is only a part of the amount after the decimal point, and leading zeros need to be added
		case l != 0:
			// Use fmt package to add leading zeros to the string
			d := fmt.Sprintf("%%0%dd", unitDigits-l)
			// Set postComma to resulting string, after removing trailing zeros
			postComma = strings.TrimRight(fmt.Sprintf(d, 0)+s, "0")
		}

		fullAmount = trimmedAmount
		if postComma != "" {
			fullAmount += "." + postComma
		}

		// limit floating part
		if len(postComma) > digits {
			postComma = postComma[:digits]
		}

		// set floating point
		if postComma != "" {
			trimmedAmount += "." + postComma
		}
	}

	if trimAmount {
		return proceed + trimmedAmount
	}

	return proceed + fullAmount
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
// It first checks for any skipped nonces that can be reused.
func (wallet *Wallet) GetNextNonce() uint64 {
	wallet.nonceMutex.Lock()
	defer wallet.nonceMutex.Unlock()

	if len(wallet.skippedNonces) > 0 {
		// Sort skipped nonces to use the lowest one first
		sort.Slice(wallet.skippedNonces, func(i, j int) bool {
			return wallet.skippedNonces[i] < wallet.skippedNonces[j]
		})

		lowestIndex := -1

		for idx, skipped := range wallet.skippedNonces {
			if skipped < wallet.confirmedTxCount {
				continue
			}

			lowestIndex = idx
			break
		}

		if lowestIndex > -1 {
			// Take the lowest skipped nonce
			nonce := wallet.skippedNonces[lowestIndex]
			wallet.skippedNonces = wallet.skippedNonces[lowestIndex+1:]

			logrus.Infof("Reusing skipped nonce %d for wallet %s", nonce, wallet.address.Hex())
			return nonce
		} else {
			wallet.skippedNonces = wallet.skippedNonces[:0]
		}
	}

	// No skipped nonces, use the next sequential nonce
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

// BuildLegacyTx builds and signs a legacy transaction.
// It automatically assigns the next available nonce and signs the transaction.
func (wallet *Wallet) BuildLegacyTx(txData *types.LegacyTx) (*types.Transaction, error) {
	txData.Nonce = wallet.GetNextNonce()
	return wallet.signTx(txData)
}

// BuildAccessListTx builds and signs an access list transaction.
// It automatically assigns the next available nonce and signs the transaction.
func (wallet *Wallet) BuildAccessListTx(txData *types.AccessListTx) (*types.Transaction, error) {
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
	if wallet.privkey == nil {
		return nil, errors.New("wallet has no private key")
	}

	transactor, err := bind.NewKeyedTransactorWithChainID(wallet.privkey, wallet.chainid)
	if err != nil {
		return nil, err
	}

	transactor.Context = ctx
	transactor.From = wallet.address
	nonce := wallet.GetNextNonce()
	transactor.Nonce = big.NewInt(0).SetUint64(nonce)

	transactor.GasTipCap = txData.GasTipCap.ToBig()
	transactor.GasFeeCap = txData.GasFeeCap.ToBig()
	transactor.GasLimit = txData.Gas
	transactor.Value = txData.Value.ToBig()
	transactor.NoSend = true

	tx, err := buildFn(transactor)
	if err != nil {
		wallet.MarkSkippedNonce(nonce)
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

// ReplaceLegacyTx builds a replacement legacy transaction with a specific nonce.
// This is useful for replacing stuck transactions with higher gas prices.
func (wallet *Wallet) ReplaceLegacyTx(txData *types.LegacyTx, nonce uint64) (*types.Transaction, error) {
	txData.Nonce = nonce
	return wallet.signTx(txData)
}

// ReplaceAccessListTx builds a replacement access list transaction with a specific nonce.
// This is useful for replacing stuck transactions with higher gas prices.
func (wallet *Wallet) ReplaceAccessListTx(txData *types.AccessListTx, nonce uint64) (*types.Transaction, error) {
	txData.ChainID = wallet.chainid
	txData.Nonce = nonce
	return wallet.signTx(txData)
}

// ReplaceSetCodeTx builds a replacement set code transaction with a specific nonce.
// This is useful for replacing stuck set code transactions with higher gas prices.
func (wallet *Wallet) ReplaceSetCodeTx(txData *types.SetCodeTx, nonce uint64) (*types.Transaction, error) {
	txData.ChainID = uint256.NewInt(wallet.chainid.Uint64())
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

func (wallet *Wallet) ResetNoncesIfNeeded(ctx context.Context, client *Client) error {
	if !wallet.needNonceResync {
		return nil
	}

	err := client.UpdateWallet(ctx, wallet)
	if err != nil {
		return fmt.Errorf("failed to refresh wallet state: %v", err)
	}

	wallet.needNonceResync = false

	return nil
}

// MarkSkippedNonce marks a nonce as skipped/failed so it can be reused later.
// This should be called when a transaction fails to submit.
func (wallet *Wallet) MarkSkippedNonce(nonce uint64) {
	wallet.nonceMutex.Lock()
	defer wallet.nonceMutex.Unlock()

	if nonce >= wallet.confirmedTxCount {
		return
	}

	for _, skipped := range wallet.skippedNonces {
		if skipped == nonce {
			return
		}
	}

	wallet.skippedNonces = append(wallet.skippedNonces, nonce)
}

func (wallet *Wallet) MarkNeedResync() {
	wallet.needNonceResync = true
}

// signTx signs a transaction using the wallet's private key and chain ID.
// It creates a new transaction from the provided transaction data and signs it
// using the latest signer for the wallet's configured chain ID.
func (wallet *Wallet) signTx(txData types.TxData) (*types.Transaction, error) {
	if wallet.privkey == nil {
		return nil, errors.New("wallet has no private key")
	}

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
func (wallet *Wallet) getTxNonceChan(tx *types.Transaction, options *SendTransactionOptions) (*nonceStatus, bool) {
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
				if existingTx.Options == nil && options != nil {
					existingTx.Options = options
				}
				return nonceChan, false
			}
		}

		nonceChan.txs = append(nonceChan.txs, &PendingTx{
			Tx:        tx,
			Submitted: time.Now(),
			Options:   options,
		})
		return nonceChan, false
	}

	nonceChan = &nonceStatus{
		txs: []*PendingTx{
			{
				Tx:        tx,
				Submitted: time.Now(),
				Options:   options,
			},
		},
		channel: make(chan bool),
	}
	wallet.txNonceChans[targetNonce] = nonceChan

	return nonceChan, len(wallet.txNonceChans) == 1
}

func (wallet *Wallet) dropPendingTx(tx *types.Transaction) {
	wallet.txNonceMutex.Lock()
	defer wallet.txNonceMutex.Unlock()

	nonceChan := wallet.txNonceChans[tx.Nonce()]
	if nonceChan == nil {
		return
	}

	txs := make([]*PendingTx, 0, len(nonceChan.txs))
	for _, pendingTx := range nonceChan.txs {
		if pendingTx.Tx != tx {
			txs = append(txs, pendingTx)
		}
	}
	nonceChan.txs = txs
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

// GetLowestPendingNonce returns the lowest pending nonce for this wallet.
// Returns 0, false if there are no pending transactions.
func (wallet *Wallet) GetLowestPendingNonce() (uint64, bool) {
	wallet.txNonceMutex.Lock()
	defer wallet.txNonceMutex.Unlock()

	if len(wallet.txNonceChans) == 0 {
		return 0, false
	}

	lowestNonce := uint64(0)
	first := true
	for nonce := range wallet.txNonceChans {
		if first || nonce < lowestNonce {
			lowestNonce = nonce
			first = false
		}
	}

	return lowestNonce, true
}

// GetNonceGaps returns missing nonces between confirmedTxCount and the lowest pending nonce.
// This helps detect nonce gaps that need to be filled with dummy transactions.
func (wallet *Wallet) GetNonceGaps() []uint64 {
	wallet.txNonceMutex.Lock()
	defer wallet.txNonceMutex.Unlock()

	if len(wallet.txNonceChans) == 0 {
		return nil
	}

	// Find the lowest pending nonce
	lowestPendingNonce := uint64(0)
	first := true
	for nonce := range wallet.txNonceChans {
		if first || nonce < lowestPendingNonce {
			lowestPendingNonce = nonce
			first = false
		}
	}

	// No gap if lowest pending nonce equals confirmed count
	if lowestPendingNonce <= wallet.confirmedTxCount {
		return nil
	}

	// Build list of missing nonces
	gaps := make([]uint64, 0, lowestPendingNonce-wallet.confirmedTxCount)
	for nonce := wallet.confirmedTxCount; nonce < lowestPendingNonce; nonce++ {
		// Check if this nonce already has a pending tx
		if _, exists := wallet.txNonceChans[nonce]; !exists {
			gaps = append(gaps, nonce)
		}
	}

	return gaps
}

// BuildFillerTx creates a simple self-transfer transaction to fill a nonce gap.
// The transaction sends 0 value to the wallet's own address with minimal gas.
func (wallet *Wallet) BuildFillerTx(nonce uint64, gasTipCap, gasFeeCap *big.Int) (*types.Transaction, error) {
	txData := &types.DynamicFeeTx{
		ChainID:   wallet.chainid,
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       21000, // Minimum gas for simple transfer
		To:        &wallet.address,
		Value:     big.NewInt(0),
		Data:      nil,
	}
	return wallet.signTx(txData)
}
