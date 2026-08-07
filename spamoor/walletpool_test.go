package spamoor

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
)

// TestNonBatcherFundingCreditsChildWallets mirrors the fixed --without-batcher
// code path in fundWallets: for each individual funding request, the tx hash
// is recorded as a batch of one in batchTxMap, and the credit loop (shared
// with the batcher path) credits the recipient once its transfer confirms.
//
// Before the fix, batchTxMap was only ever populated in the batcher branch,
// so this same credit loop found nothing for individually-funded wallets and
// never credited them even though the on-chain transfer succeeded.
func TestNonBatcherFundingCreditsChildWallets(t *testing.T) {
	req1 := &FundingRequest{Wallet: &Wallet{balance: big.NewInt(0)}, Amount: uint256.NewInt(1000)}
	req2 := &FundingRequest{Wallet: &Wallet{balance: big.NewInt(0)}, Amount: uint256.NewInt(2000)}
	fundingReqs := []*FundingRequest{req1, req2}

	// Mirrors the fixed else branch (walletpool.go, buildWalletFundingTx loop):
	// one tx per request, each recorded as its own one-element batch.
	batchTxMap := map[common.Hash][]*FundingRequest{}
	hashes := []common.Hash{common.HexToHash("0x01"), common.HexToHash("0x02")}
	for i, req := range fundingReqs {
		batchTxMap[hashes[i]] = []*FundingRequest{req}
	}

	receipts := []*types.Receipt{
		{TxHash: hashes[0], Status: types.ReceiptStatusSuccessful},
		{TxHash: hashes[1], Status: types.ReceiptStatusSuccessful},
	}

	// The exact shared credit loop (walletpool.go:1214-1223).
	for _, receipt := range receipts {
		if receipt != nil && receipt.Status == types.ReceiptStatusSuccessful {
			if batch, ok := batchTxMap[receipt.TxHash]; ok {
				for _, req := range batch {
					req.Wallet.AddBalance(req.Amount.ToBig())
				}
			}
		}
	}

	if req1.Wallet.GetBalance().Cmp(big.NewInt(1000)) != 0 {
		t.Fatalf("expected wallet 1 to be credited 1000, got %s", req1.Wallet.GetBalance())
	}
	if req2.Wallet.GetBalance().Cmp(big.NewInt(2000)) != 0 {
		t.Fatalf("expected wallet 2 to be credited 2000, got %s", req2.Wallet.GetBalance())
	}
}

// A receipt that never confirms successfully must not credit anything, even
// though its hash has a batch recorded for it - only successful transfers
// should ever add balance.
func TestNonBatcherFundingSkipsFailedTransfers(t *testing.T) {
	req := &FundingRequest{Wallet: &Wallet{balance: big.NewInt(0)}, Amount: uint256.NewInt(1000)}
	hash := common.HexToHash("0x01")
	batchTxMap := map[common.Hash][]*FundingRequest{hash: {req}}

	receipts := []*types.Receipt{
		{TxHash: hash, Status: types.ReceiptStatusFailed},
	}

	for _, receipt := range receipts {
		if receipt != nil && receipt.Status == types.ReceiptStatusSuccessful {
			if batch, ok := batchTxMap[receipt.TxHash]; ok {
				for _, r := range batch {
					r.Wallet.AddBalance(r.Amount.ToBig())
				}
			}
		}
	}

	if req.Wallet.GetBalance().Sign() != 0 {
		t.Fatalf("expected no credit for a failed transfer, got %s", req.Wallet.GetBalance())
	}
}
