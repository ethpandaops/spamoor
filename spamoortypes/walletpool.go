package spamoortypes

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

// WalletSelectionMode defines how wallets are selected from the pool.
type WalletSelectionMode uint8

var (
	// SelectWalletByIndex selects a wallet by index (modulo pool size).
	SelectWalletByIndex WalletSelectionMode = 0
	// SelectWalletRandom selects a random wallet from the pool.
	SelectWalletRandom WalletSelectionMode = 1
	// SelectWalletRoundRobin selects wallets in round-robin fashion.
	SelectWalletRoundRobin WalletSelectionMode = 2
)

// WellKnownWalletConfig defines configuration for a named wallet with custom funding settings.
type WellKnownWalletConfig struct {
	Name          string
	RefillAmount  *uint256.Int
	RefillBalance *uint256.Int
	VeryWellKnown bool
}

// WalletPool defines the interface for managing a pool of child wallets derived from a root wallet
// with automatic funding and balance monitoring. It provides wallet selection strategies, automatic
// refills when balances drop below thresholds, and batch funding operations for efficiency.
type WalletPool interface {
	// Context and Dependencies
	GetContext() context.Context
	GetTxPool() TxPool
	GetClientPool() ClientPool
	GetRootWallet() Wallet
	WithRootWalletLock(ctx context.Context, txCount int, lockedLogFn func(), lockedFn func() error) error
	GetChainId() *big.Int

	// Configuration Management
	LoadConfig(configYaml string) error
	MarshalConfig() (string, error)
	SetWalletCount(count uint64)
	SetRunFundings(runFundings bool)
	AddWellKnownWallet(config *WellKnownWalletConfig)
	SetRefillAmount(amount *uint256.Int)
	SetRefillBalance(balance *uint256.Int)
	SetWalletSeed(seed string)
	SetRefillInterval(interval uint64)
	SetTransactionTracker(tracker func(err error))
	GetTransactionTracker() func(err error)

	// Client and Wallet Selection
	GetClient(mode ClientSelectionMode, input int, group string) Client
	GetWallet(mode WalletSelectionMode, input int) Wallet
	GetWellKnownWallet(name string) Wallet
	GetVeryWellKnownWalletAddress(name string) common.Address
	GetWalletName(address common.Address) string
	GetAllWallets() []Wallet
	GetConfiguredWalletCount() uint64
	GetWalletCount() uint64

	// Wallet Management
	PrepareWallets() error
	CollectPoolWallets(walletMap map[common.Address]Wallet)
	CheckChildWalletBalance(childWallet Wallet) error
	ReclaimFunds(ctx context.Context, client Client) error
}
