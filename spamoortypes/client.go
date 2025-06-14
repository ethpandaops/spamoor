package spamoortypes

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Client defines the interface for a spamoor client. It provides necessary methods to interact with a ethereum execution client from spamoor.
type Client interface {
	GetName() string
	GetClientGroup() string
	GetEthClient() bind.ContractBackend
	GetRPCHost() string
	GetTimeout() time.Duration
	SetTimeout(timeout time.Duration)
	UpdateWallet(ctx context.Context, wallet Wallet) error
	SetClientGroup(group string)
	IsEnabled() bool
	SetEnabled(enabled bool)
	GetChainId(ctx context.Context) (*big.Int, error)
	GetNonceAt(ctx context.Context, wallet common.Address, blockNumber *big.Int) (uint64, error)
	GetPendingNonceAt(ctx context.Context, wallet common.Address) (uint64, error)
	GetBalanceAt(ctx context.Context, wallet common.Address) (*big.Int, error)
	GetSuggestedFee(ctx context.Context) (*big.Int, *big.Int, error)
	SendTransaction(ctx context.Context, tx *types.Transaction) error
	SendRawTransaction(ctx context.Context, tx []byte) error
	GetTransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	GetBlockHeight(ctx context.Context) (uint64, error)
	GetLastBlockHeight() (uint64, time.Time)
	GetClientVersion(ctx context.Context) (string, error)
	GetBlock(ctx context.Context, blockNumber uint64) (*types.Block, error)
	GetBlockReceipts(ctx context.Context, blockNumber uint64) ([]*types.Receipt, error)
	GetLatestGasLimit(ctx context.Context) (uint64, error)
}
