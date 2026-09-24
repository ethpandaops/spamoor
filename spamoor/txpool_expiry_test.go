package spamoor

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"

	"github.com/ethpandaops/spamoor/txtypes"
)

// expiringFrameTx builds a frame transaction whose expiry verifier frame carries deadline.
func expiringFrameTx(deadline uint64, feeCap, tipCap uint64) *txtypes.Transaction {
	return txtypes.NewTx(&txtypes.FrameTx{
		ChainID: uint256.NewInt(1),
		Sender:  common.HexToAddress("0x1111111111111111111111111111111111111111"),
		Fees: txtypes.FrameFees{
			GasFeeCap: uint256.NewInt(feeCap),
			GasTipCap: uint256.NewInt(tipCap),
		},
		Frames: []*txtypes.Frame{
			txtypes.ExpiryFrame(deadline, 5000),
			txtypes.SelfVerifyFrame(txtypes.FrameLimits{Execution: 5000}),
		},
	})
}

// A transaction whose expiry deadline the chain has already reached can never be included
// again, so it has to be recognised as expired rather than waited on. The head block's
// timestamp is the reference: a later block can only carry a later one.
func TestIsExpiredTx(t *testing.T) {
	tests := []struct {
		name      string
		tx        *txtypes.Transaction
		chainTime uint64
		want      bool
	}{
		{"deadline ahead of the chain", expiringFrameTx(2000, 10, 1), 1999, false},
		{"deadline reached", expiringFrameTx(2000, 10, 1), 2000, true},
		{"deadline behind the chain", expiringFrameTx(2000, 10, 1), 2001, true},
		{"chain time unknown", expiringFrameTx(2000, 10, 1), 0, false},
		{
			"frame transaction without an expiry frame",
			txtypes.NewTx(&txtypes.FrameTx{
				ChainID: uint256.NewInt(1),
				Fees:    txtypes.FrameFees{GasFeeCap: uint256.NewInt(10), GasTipCap: uint256.NewInt(1)},
				Frames:  []*txtypes.Frame{txtypes.SelfVerifyFrame(txtypes.FrameLimits{Execution: 5000})},
			}),
			5000,
			false,
		},
		{
			"a type that cannot expire",
			txtypes.NewTx(&txtypes.DynamicFeeTx{
				ChainID:   big.NewInt(1),
				GasFeeCap: big.NewInt(10),
				GasTipCap: big.NewInt(1),
			}),
			5000,
			false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := isExpiredTx(test.tx, test.chainTime); got != test.want {
				t.Fatalf("isExpiredTx = %v, want %v", got, test.want)
			}
		})
	}
}

// A replacement has to outprice the transaction it replaces, or the pools refuse it and
// the nonce stays blocked. The bump grows with the attempt, so a client asking for more
// than the first offer is eventually satisfied.
func TestReplacementFeesOutpriceTheReplacedTx(t *testing.T) {
	pool := &TxPool{}
	pool.currentBaseFee = big.NewInt(1)

	replaced := expiringFrameTx(2000, 100_000_000_000, 4_000_000_000)

	var lastFeeCap, lastTipCap *big.Int

	for attempt := uint64(1); attempt <= 3; attempt++ {
		feeCap, tipCap := pool.replacementFees(replaced, attempt)

		if feeCap.Cmp(replaced.GasFeeCap()) <= 0 {
			t.Fatalf("attempt %d: fee cap %v does not exceed the replaced %v", attempt, feeCap, replaced.GasFeeCap())
		}

		if tipCap.Cmp(replaced.GasTipCap()) <= 0 {
			t.Fatalf("attempt %d: tip cap %v does not exceed the replaced %v", attempt, tipCap, replaced.GasTipCap())
		}

		if tipCap.Cmp(feeCap) > 0 {
			t.Fatalf("attempt %d: tip cap %v above fee cap %v", attempt, tipCap, feeCap)
		}

		if lastFeeCap != nil && feeCap.Cmp(lastFeeCap) <= 0 {
			t.Fatalf("attempt %d: fee cap %v did not grow over %v", attempt, feeCap, lastFeeCap)
		}

		if lastTipCap != nil && tipCap.Cmp(lastTipCap) <= 0 {
			t.Fatalf("attempt %d: tip cap %v did not grow over %v", attempt, tipCap, lastTipCap)
		}

		lastFeeCap, lastTipCap = feeCap, tipCap
	}
}

// A transaction replacing one that was priced below the current base fee must still be
// worth including, so the current chain price is a floor rather than the old price.
func TestReplacementFeesRespectTheCurrentBaseFee(t *testing.T) {
	pool := &TxPool{}
	pool.currentBaseFee = big.NewInt(50_000_000_000)

	replaced := expiringFrameTx(2000, 1, 1)

	feeCap, tipCap := pool.replacementFees(replaced, 1)

	if feeCap.Cmp(big.NewInt(100_000_000_000)) != 0 {
		t.Fatalf("fee cap %v, want twice the base fee", feeCap)
	}

	if tipCap.Cmp(big.NewInt(1_000_000_000)) != 0 {
		t.Fatalf("tip cap %v, want the 1 gwei floor", tipCap)
	}
}

// A rebroadcast that starts at the client the previous attempt reached is wasted: that
// client already holds the transaction and takes it again, so the submit loop stops there
// and the transaction never reaches a node that would include it. Every retry has to start
// one client further along.
func TestRebroadcastStartOffsetAdvancesWithRetries(t *testing.T) {
	hash := common.HexToHash("0x0102000000000000000000000000000000000000000000000000000000000000")

	const clientCount = 7

	seen := map[int]bool{}

	for retry := uint64(0); retry < clientCount; retry++ {
		offset := rebroadcastStartOffset(0, hash, retry)

		client := offset % clientCount
		if seen[client] {
			t.Fatalf("retry %d starts at client %d again", retry, client)
		}

		seen[client] = true
	}

	if len(seen) != clientCount {
		t.Fatalf("retries reached %d of %d clients", len(seen), clientCount)
	}

	// A caller that pins the starting client keeps it, still advanced per retry.
	if got := rebroadcastStartOffset(3, hash, 0); got != 3 {
		t.Fatalf("configured offset became %d, want 3", got)
	}

	if got := rebroadcastStartOffset(3, hash, 2); got != 5 {
		t.Fatalf("configured offset on retry 2 became %d, want 5", got)
	}
}

// Expiry is judged against chain time, so the pool has to keep answering with a block
// timestamp even once the block it was asked about has been pruned from its window.
func TestBlockTimeFallsBackToTheNewestKnownBlock(t *testing.T) {
	pool := &TxPool{blocks: map[uint64]*BlockInfo{}}

	if got := pool.blockTime(10); got != 0 {
		t.Fatalf("blockTime with no blocks = %d, want 0", got)
	}

	pool.blocks[10] = &BlockInfo{Number: 10, Timestamp: 1000}
	pool.blocks[11] = &BlockInfo{Number: 11, Timestamp: 1012}

	if got := pool.blockTime(10); got != 1000 {
		t.Fatalf("blockTime(10) = %d, want 1000", got)
	}

	if got := pool.blockTime(9); got != 1012 {
		t.Fatalf("blockTime of a pruned block = %d, want the newest known 1012", got)
	}
}

// Pricing must not depend on block statistics being initialised: before the first block
// is processed there is no base fee, and the helper has to cope rather than crash. The tip
// floor never lifts the tip above the fee cap, since that pair would be invalid.
func TestReplacementFeesWithoutABaseFee(t *testing.T) {
	pool := &TxPool{}

	feeCap, tipCap := pool.replacementFees(expiringFrameTx(2000, 100, 8), 1)

	if feeCap.Cmp(big.NewInt(125)) != 0 || tipCap.Cmp(big.NewInt(125)) != 0 {
		t.Fatalf("fees without a base fee = %v/%v, want 125/125 (tip clamped to the fee cap)", feeCap, tipCap)
	}
}
