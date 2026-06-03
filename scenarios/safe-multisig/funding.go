package safemultisig

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/holiman/uint256"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/spamoor"
)

// Safe balances are kept intentionally low and topped up incrementally rather
// than pre-funded with a large lump sum. A safe is refilled when its balance
// drops below safeRefillThresholdCalls worth of max-value transfers, back up to
// safeRefillTargetCalls worth. With the default sub-gwei eoa-value this keeps
// each safe holding only a tiny float.
const (
	safeRefillThresholdCalls = 64
	safeRefillTargetCalls    = 256
	// safeFundingGas is the per-recipient gas budget for funding a safe. A Safe
	// proxy receiving ETH delegatecalls its singleton and emits a SafeReceived
	// event (~28k under Amsterdam), far more than the plain-EOA transfer the
	// pool's default funding budget assumes - so it must be passed explicitly.
	safeFundingGas = uint64(60000)
)

// runSafeTopUp periodically tops up the safe balances until the context is
// cancelled. The interval is --funding-interval slots (converted via the global
// slot duration). The initial (synchronous) top-up is done once by the caller
// before this loop starts.
func (s *Scenario) runSafeTopUp(ctx context.Context) {
	slots := s.options.FundingInterval
	if slots == 0 {
		slots = 1
	}
	interval := time.Duration(slots) * scenario.GlobalSlotDuration
	if interval <= 0 {
		interval = 32 * 12 * time.Second
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := s.topUpSafes(ctx); err != nil && ctx.Err() == nil {
				s.logger.Warnf("could not top up safe balances: %v", err)
			}
		}
	}
}

// topUpSafes reads every tracked safe's balance (refreshing the per-safe balance
// used to gate EOA transfers) and refills the ones that have fallen below the
// low-balance threshold. The transfers are delegated to the wallet pool's
// FundAddresses, so they reuse the batcher contract (when enabled) and the
// root-wallet-locked funding path used for child wallets. It is a no-op when EOA
// value transfers are disabled.
func (s *Scenario) topUpSafes(ctx context.Context) error {
	if s.options.EoaValue == 0 {
		return nil
	}

	refs := s.snapshotSafes()
	if len(refs) == 0 {
		return nil
	}

	client := s.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0),
		spamoor.WithClientGroup(s.options.ClientGroup),
	)
	if client == nil {
		return scenario.ErrNoClients
	}
	ethClient := client.GetEthClient()

	perCall := gweiToWei(s.options.EoaValue)
	threshold := new(big.Int).Mul(perCall, big.NewInt(safeRefillThresholdCalls))
	target := new(big.Int).Mul(perCall, big.NewInt(safeRefillTargetCalls))

	reqs := make([]*spamoor.FundingRequest, 0)
	for _, ref := range refs {
		balance, err := ethClient.BalanceAt(ctx, ref.addr, nil)
		if err != nil {
			return fmt.Errorf("could not read safe balance %s: %w", ref.addr.Hex(), err)
		}

		// Refresh the authoritative balance used by the exec path to decide
		// whether a value transfer can be sourced from this safe.
		ref.pool.mu.Lock()
		ref.entry.balance = balance
		ref.pool.mu.Unlock()

		if balance.Cmp(threshold) >= 0 {
			continue
		}

		targetWallet := spamoor.NewWallet(nil, ref.addr)
		targetWallet.SetBalance(balance)
		reqs = append(reqs, &spamoor.FundingRequest{
			Wallet: targetWallet,
			Amount: uint256.MustFromBig(new(big.Int).Sub(target, balance)),
			// The safe proxy already exists (has code), so funding it does not
			// create a new account - no extra state-creation gas needed.
			IsEmpty: false,
			// But its receive() does real work (delegatecall + event), so the
			// default EOA-sized funding gas is insufficient.
			Gas: safeFundingGas,
		})
	}
	if len(reqs) == 0 {
		return nil
	}

	if err := s.walletPool.FundAddresses(reqs); err != nil {
		return fmt.Errorf("could not fund safes: %w", err)
	}

	s.logger.Infof("topped up %d safe balances", len(reqs))
	return nil
}
