package seaporttrades

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"

	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
)

// The reserve/commit inventory protocol. A trade builder pops the asset(s) it
// needs "in flight" (under the mutex) while building the order, then either
// commits them to the destination once the fulfillment confirms or rolls them
// back on failure. Popped assets live in neither side's books in between, so a
// not-yet-confirmed NFT is never offered in a second order. All of these hold
// m.mu for the duration; none does I/O.

// reserveMarketNFT pops a market-owned tokenId to list. ok=false when the market
// has no inventory to sell.
func (m *Market) reserveMarketNFT() (*big.Int, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	n := len(m.marketNFTs)
	if n == 0 {
		return nil, false
	}
	id := m.marketNFTs[n-1]
	m.marketNFTs = m.marketNFTs[:n-1]
	return id, true
}

// reserveWalletNFT pops a tokenId currently held by the given trader. ok=false
// when the trader holds nothing to sell.
func (m *Market) reserveWalletNFT(owner common.Address) (*big.Int, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	ids := m.walletNFTs[owner]
	n := len(ids)
	if n == 0 {
		return nil, false
	}
	id := ids[n-1]
	m.walletNFTs[owner] = ids[:n-1]
	return id, true
}

func (m *Market) giveMarketNFT(id *big.Int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.marketNFTs = append(m.marketNFTs, id)
}

func (m *Market) giveWalletNFT(owner common.Address, id *big.Int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.walletNFTs[owner] = append(m.walletNFTs[owner], id)
}

// reserveCoin optimistically debits the payer's cached balance. ok=false when
// the payer cannot afford the amount, leaving the balance untouched.
func (m *Market) reserveCoin(payer common.Address, amount *big.Int) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	bal := m.coinBal[payer]
	if bal == nil || bal.Cmp(amount) < 0 {
		return false
	}
	m.coinBal[payer] = new(big.Int).Sub(bal, amount)
	return true
}

// creditCoin adds amount to an address's cached balance (commit of a reservation
// to the payee, or refund to the payer on rollback).
func (m *Market) creditCoin(addr common.Address, amount *big.Int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	bal := m.coinBal[addr]
	if bal == nil {
		bal = big.NewInt(0)
	}
	m.coinBal[addr] = new(big.Int).Add(bal, amount)
}

// marketInventorySize reports how many NFTs the market currently has available
// to list (used by the replenisher's low-water check).
func (m *Market) marketInventorySize() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.marketNFTs)
}

// walletNFTCount reports how many NFTs a trader currently holds (used to bias
// the buy/sell decision toward shedding a large inventory).
func (m *Market) walletNFTCount(owner common.Address) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.walletNFTs[owner])
}

// maybeReplenish self-mints additional NFTs to the market when its listing
// inventory runs low, and tops up the market's stablecoin float when it dips -
// this is what lets a long run keep generating both buys and sells without ever
// needing external setup. It is safe to call concurrently: a single replenish
// runs at a time (guarded by m.replenishing); the rest return immediately. The
// mint/top-up transactions are submitted from the market wallet and awaited, so
// the freshly minted ids are only added to the in-memory pool once confirmed.
func (m *Market) maybeReplenish(ctx context.Context, client *spamoor.Client, feeCap, tipCap *big.Int) {
	lowNFT := m.marketInventorySize() < int(m.options.ReplenishThreshold)
	lowCoin := m.coinBalanceBelow(m.marketAddr, m.coinSeedAmount())
	if !lowNFT && !lowCoin {
		return
	}

	m.mu.Lock()
	if m.replenishing {
		m.mu.Unlock()
		return
	}
	m.replenishing = true
	m.mu.Unlock()
	defer func() {
		m.mu.Lock()
		m.replenishing = false
		m.mu.Unlock()
	}()

	if lowNFT {
		if err := m.replenishNFTs(ctx, client, feeCap, tipCap); err != nil {
			m.logger.Warnf("could not replenish market NFTs: %v", err)
		}
	}
	if lowCoin {
		if err := m.replenishCoin(ctx, client, feeCap, tipCap); err != nil {
			m.logger.Warnf("could not replenish market coin: %v", err)
		}
	}
}

func (m *Market) coinBalanceBelow(addr common.Address, threshold *big.Int) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	bal := m.coinBal[addr]
	return bal == nil || bal.Cmp(threshold) < 0
}

// replenishNFTs mints a fresh contiguous id range to the market and adds the ids
// to the listing pool after confirmation.
func (m *Market) replenishNFTs(ctx context.Context, client *spamoor.Client, feeCap, tipCap *big.Int) error {
	count := m.options.ReplenishBatch
	if count == 0 {
		return nil
	}

	m.mu.Lock()
	startID := m.nextTokenID
	m.nextTokenID += count
	m.mu.Unlock()

	mintGas := mintBaseGas + count*mintPerTokenGas
	tx, err := m.marketWallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       mintGas,
		Value:     uint256.NewInt(0),
	}, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return m.deployment.NFT.MintBatch(opts, m.marketAddr, new(big.Int).SetUint64(startID), new(big.Int).SetUint64(count))
	})
	if err != nil {
		return fmt.Errorf("could not build replenish mint: %w", err)
	}

	if _, err := m.walletPool.GetTxPool().SendAndAwaitTransaction(ctx, m.marketWallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		ClientGroup: m.options.ClientGroup,
		Rebroadcast: true,
	}); err != nil {
		return fmt.Errorf("could not send replenish mint: %w", err)
	}

	for _, id := range idRange(startID, count) {
		m.giveMarketNFT(id)
	}
	m.logger.Infof("replenished market inventory with %d nfts (ids %d..%d)", count, startID, startID+count-1)
	return nil
}

// replenishCoin mints additional stablecoin to the market so its bids never run
// dry, and updates the cached balance after confirmation.
func (m *Market) replenishCoin(ctx context.Context, client *spamoor.Client, feeCap, tipCap *big.Int) error {
	amount := new(big.Int).Mul(m.coinSeedAmount(), big.NewInt(int64(m.walletPool.GetConfiguredWalletCount()+10)))
	tx, err := m.marketWallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       coinMintGasLimit,
		Value:     uint256.NewInt(0),
	}, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return m.deployment.Coin.Mint(opts, m.marketAddr, amount)
	})
	if err != nil {
		return fmt.Errorf("could not build replenish coin mint: %w", err)
	}

	if _, err := m.walletPool.GetTxPool().SendAndAwaitTransaction(ctx, m.marketWallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		ClientGroup: m.options.ClientGroup,
		Rebroadcast: true,
	}); err != nil {
		return fmt.Errorf("could not send replenish coin mint: %w", err)
	}

	m.creditCoin(m.marketAddr, amount)
	m.logger.Infof("replenished market coin float by %s", amount.String())
	return nil
}
