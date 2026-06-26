package seaporttrades

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"

	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
)

// Static gas limits for the setup/maintenance txs. These deliberately avoid
// per-tx estimation (hundreds of wallets would mean hundreds of extra RPC round
// trips) and keep generous headroom for the Amsterdam state-creation surcharge
// on fresh storage slots.
const (
	// approvalGasLimit covers setApprovalForAll and ERC20 approve. Under the
	// Amsterdam fee schedule a fresh approval/allowance slot costs ~128k (the
	// uniswap-swaps scenario uses the same 250k budget), well above the plain
	// ~46k a warm slot would imply - a tighter limit OOG-reverts the approval.
	approvalGasLimit = 250000
	// coinMintGasLimit covers the ERC20 mint, which writes the recipient balance
	// slot and totalSupply (both fresh on first mint) under Amsterdam.
	coinMintGasLimit = 300000
	mintBaseGas      = 80000  // base for an NFT mintBatch tx
	mintPerTokenGas  = 130000 // per fresh NFT slot under Amsterdam (measured ~126k)
)

// maxAllowance is the unlimited ERC20 allowance granted to Seaport.
var maxAllowance = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))

// Market is the single standing counterparty for every trade. It is the offerer
// (maker) that signs both listings (sell side) and bids (buy side) off-chain;
// the child wallets are the on-chain fulfillers (takers). It owns a pool of NFTs
// available to list, tracks which NFTs each trader currently holds, and caches
// stablecoin balances - all so the trade builder can pick a fulfillable
// direction and size without per-tx RPC reads. NFT ownership churns as trades
// confirm, so the inventory is mutated under a mutex with a reserve/commit
// protocol: assets are popped "in flight" when an order is built and only handed
// to the destination once the fulfillment confirms (rolled back on failure), so
// an unconfirmed NFT is never offered in a second order.
type Market struct {
	walletPool *spamoor.WalletPool
	logger     *logrus.Entry
	options    *ScenarioOptions
	deployment *DeploymentInfo

	marketWallet *spamoor.Wallet
	marketAddr   common.Address
	marketKey    *ecdsa.PrivateKey

	seaportAddr     common.Address
	domainSeparator [32]byte
	counter         *big.Int // market's Seaport counter (constant for the run)

	mu           sync.Mutex
	idBase       uint64                        // first tokenId this run mints (above any prior run's range)
	marketNFTs   []*big.Int                    // tokenIds owned by the market, available to list
	walletNFTs   map[common.Address][]*big.Int // tokenIds currently owned by each trader
	coinBal      map[common.Address]*big.Int   // cached stablecoin balance per address
	nextTokenID  uint64                        // next free tokenId for replenishment mints
	replenishing bool                          // guards a single in-flight market-NFT replenish
}

// NewMarket binds the deployed contracts to a market counterparty backed by the
// "market" well-known wallet.
func NewMarket(walletPool *spamoor.WalletPool, logger *logrus.Entry, options *ScenarioOptions, deployment *DeploymentInfo) (*Market, error) {
	marketWallet := walletPool.GetWellKnownWallet("market")
	if marketWallet == nil {
		return nil, fmt.Errorf("market wallet not available")
	}
	key := marketWallet.GetPrivateKey()
	if key == nil {
		return nil, fmt.Errorf("market wallet has no private key")
	}

	return &Market{
		walletPool:      walletPool,
		logger:          logger,
		options:         options,
		deployment:      deployment,
		marketWallet:    marketWallet,
		marketAddr:      marketWallet.GetAddress(),
		marketKey:       key,
		seaportAddr:     deployment.SeaportAddr,
		domainSeparator: deployment.domainSeparator,
		walletNFTs:      make(map[common.Address][]*big.Int),
		coinBal:         make(map[common.Address]*big.Int),
	}, nil
}

// childWallets returns the trader wallets in index order, excluding the
// well-known wallets that GetAllWallets() also includes. The index lines up with
// walletBaseID so the deterministic id layout is reproducible.
func (m *Market) childWallets() []*spamoor.Wallet {
	count := m.walletPool.GetConfiguredWalletCount()
	wallets := make([]*spamoor.Wallet, 0, count)
	for i := uint64(0); i < count; i++ {
		if w := m.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(i)); w != nil {
			wallets = append(wallets, w)
		}
	}
	return wallets
}

// tokenId layout: the NFT collection is global (shared across runs), but each run
// mints its own fresh id range starting at idBase (the chain's current high-water
// mark), so runs never collide on token ids. Within the run the market owns
// [idBase, idBase+MarketInventory) and trader i owns
// [idBase+MarketInventory + i*WalletInventory, +WalletInventory). Replenishment
// mints continue past the end of the run's seeded range.
func (m *Market) walletBaseID(idx uint64) uint64 {
	return m.idBase + m.options.MarketInventory + idx*m.options.WalletInventory
}

func (m *Market) seededTokenCount() uint64 {
	return m.options.MarketInventory + m.walletPool.GetConfiguredWalletCount()*m.options.WalletInventory
}

// Seed prepares all trading state: it reads the market's Seaport counter, mints
// the NFT inventories and stablecoin balances, and grants Seaport the NFT
// operator approval and ERC20 allowance for every participant. On a restart
// (contracts already deployed and seeded) it skips minting and reconstructs the
// in-memory inventory from chain instead.
func (m *Market) Seed(ctx context.Context) error {
	client := m.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0),
		spamoor.WithClientGroup(m.options.ClientGroup),
		spamoor.WithoutBuilder(),
	)
	if client == nil {
		return fmt.Errorf("no client available")
	}

	counter, err := m.deployment.Seaport.GetCounter(&bind.CallOpts{Context: ctx}, m.marketAddr)
	if err != nil {
		return fmt.Errorf("could not read market counter: %w", err)
	}
	m.counter = counter

	// The NFT collection is global and may already hold tokens from earlier runs.
	// Rather than reuse and reconcile that stale inventory, every run mints a fresh
	// id range starting above the current high-water mark, so its token ids never
	// collide with another run's.
	highWater, err := m.findHighWaterMark(ctx)
	if err != nil {
		return fmt.Errorf("could not determine nft high-water mark: %w", err)
	}
	m.idBase = highWater
	m.nextTokenID = m.idBase + m.seededTokenCount()
	if highWater > 0 {
		m.logger.Infof("nft collection already holds %d tokens, minting fresh range starting at id %d", highWater, m.idBase)
	}

	return m.seedFresh(ctx, client)
}

// findHighWaterMark returns the number of tokens already minted in the global NFT
// collection (the first tokenId whose owner read reverts). Minting is contiguous
// from id 0 across runs, so this is found with an exponential probe followed by a
// binary search - O(log n) reads rather than scanning every id.
func (m *Market) findHighWaterMark(ctx context.Context) (uint64, error) {
	exists := func(id uint64) (bool, error) {
		_, err := m.deployment.NFT.OwnerOf(&bind.CallOpts{Context: ctx}, new(big.Int).SetUint64(id))
		if err == nil {
			return true, nil
		}
		if isNonexistentToken(err) {
			return false, nil
		}
		return false, err
	}

	if ok, err := exists(0); err != nil {
		return 0, err
	} else if !ok {
		return 0, nil
	}

	// exponential search for an upper bound that does not exist
	lo, hi := uint64(0), uint64(1)
	for {
		ok, err := exists(hi)
		if err != nil {
			return 0, err
		}
		if !ok {
			break
		}
		lo = hi
		hi *= 2
	}
	// binary search in (lo, hi]: lo exists, hi does not
	for hi-lo > 1 {
		mid := lo + (hi-lo)/2
		ok, err := exists(mid)
		if err != nil {
			return 0, err
		}
		if ok {
			lo = mid
		} else {
			hi = mid
		}
	}
	return hi, nil
}

// isNonexistentToken reports whether an ownerOf error is the expected
// "nonexistent token" revert (as opposed to an RPC/transport error).
func isNonexistentToken(err error) bool {
	return err != nil && strings.Contains(strings.ToLower(err.Error()), "nonexistent token")
}

// seedFresh mints inventories + balances and sets approvals across the market and
// all trader wallets in a single multi-wallet batch, then populates the
// in-memory caches from the deterministic id layout.
func (m *Market) seedFresh(ctx context.Context, client *spamoor.Client) error {
	baseFeeWei, tipFeeWei := spamoor.ResolveFees(m.options.BaseFee, m.options.TipFee, m.options.BaseFeeWei, m.options.TipFeeWei)
	feeCap, tipCap, err := m.walletPool.GetSuggestedFees(client, baseFeeWei, tipFeeWei)
	if err != nil {
		return fmt.Errorf("could not get tx fee: %w", err)
	}

	coinSeed := m.coinSeedAmount()
	marketCoinSeed := new(big.Int).Mul(coinSeed, big.NewInt(int64(m.walletPool.GetConfiguredWalletCount()+10)))

	walletTxs := make(map[*spamoor.Wallet][]*types.Transaction)

	// Market wallet: mint its listing inventory + a large coin float, then approve.
	marketTxs, err := m.buildSeedTxs(ctx, m.marketWallet, m.idBase, m.options.MarketInventory, marketCoinSeed, feeCap, tipCap)
	if err != nil {
		return err
	}
	walletTxs[m.marketWallet] = marketTxs
	m.marketNFTs = idRange(m.idBase, m.options.MarketInventory)
	m.coinBal[m.marketAddr] = new(big.Int).Set(marketCoinSeed)

	// Each trader (child) wallet self-mints its own NFT range + coin float, then
	// approves. Iterate the child wallets explicitly by index: GetAllWallets() also
	// returns the well-known wallets (deployer, market), and seeding the market
	// here would overwrite its [0, MarketInventory) listing mint with a
	// trader-style range, leaving the listing tokens unminted.
	for idx, wallet := range m.childWallets() {
		base := m.walletBaseID(uint64(idx))
		txs, err := m.buildSeedTxs(ctx, wallet, base, m.options.WalletInventory, coinSeed, feeCap, tipCap)
		if err != nil {
			return err
		}
		walletTxs[wallet] = txs
		m.walletNFTs[wallet.GetAddress()] = idRange(base, m.options.WalletInventory)
		m.coinBal[wallet.GetAddress()] = new(big.Int).Set(coinSeed)
	}

	total := 0
	for _, txs := range walletTxs {
		total += len(txs)
	}
	m.logger.Infof("seeding seaport market: %d txs across %d wallets (mint+approve)", total, len(walletTxs))

	_, err = m.walletPool.GetTxPool().SendMultiTransactionBatch(ctx, walletTxs, &spamoor.BatchOptions{
		SendTransactionOptions: spamoor.SendTransactionOptions{
			ClientGroup: m.options.ClientGroup,
			Rebroadcast: true,
		},
		MaxRetries:   3,
		PendingLimit: 50,
		LogFn: func(confirmedCount int, totalCount int) {
			m.logger.Infof("seeding market... (%v/%v)", confirmedCount, totalCount)
		},
		LogInterval: 50,
	})
	if err != nil {
		return fmt.Errorf("could not send seed txs: %w", err)
	}
	m.logger.Infof("seaport market seeded")
	return nil
}

// buildSeedTxs builds the per-wallet setup transactions: mint an NFT id range,
// mint a coin balance, grant Seaport the NFT operator approval and the ERC20
// allowance. count==0 skips the NFT mint (a wallet with no starting inventory).
func (m *Market) buildSeedTxs(ctx context.Context, wallet *spamoor.Wallet, startID, count uint64, coinAmount *big.Int, feeCap, tipCap *big.Int) ([]*types.Transaction, error) {
	txs := make([]*types.Transaction, 0, 4)
	build := func(gas uint64, fn func(*bind.TransactOpts) (*types.Transaction, error)) error {
		tx, err := wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Gas:       gas,
			Value:     uint256.NewInt(0),
		}, fn)
		if err != nil {
			return err
		}
		txs = append(txs, tx)
		return nil
	}

	if count > 0 {
		mintGas := mintBaseGas + count*mintPerTokenGas
		if err := build(mintGas, func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return m.deployment.NFT.MintBatch(opts, wallet.GetAddress(), new(big.Int).SetUint64(startID), new(big.Int).SetUint64(count))
		}); err != nil {
			return nil, fmt.Errorf("could not build mintBatch for %s: %w", wallet.GetAddress().Hex(), err)
		}
	}
	if err := build(coinMintGasLimit, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return m.deployment.Coin.Mint(opts, wallet.GetAddress(), coinAmount)
	}); err != nil {
		return nil, fmt.Errorf("could not build coin mint for %s: %w", wallet.GetAddress().Hex(), err)
	}
	if err := build(approvalGasLimit, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return m.deployment.NFT.SetApprovalForAll(opts, m.seaportAddr, true)
	}); err != nil {
		return nil, fmt.Errorf("could not build setApprovalForAll for %s: %w", wallet.GetAddress().Hex(), err)
	}
	if err := build(approvalGasLimit, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return m.deployment.Coin.Approve(opts, m.seaportAddr, maxAllowance)
	}); err != nil {
		return nil, fmt.Errorf("could not build coin approve for %s: %w", wallet.GetAddress().Hex(), err)
	}
	return txs, nil
}

// coinSeedAmount returns the per-trader starting stablecoin balance: enough to
// cover many max-price buys before relying on sell proceeds to recirculate.
func (m *Market) coinSeedAmount() *big.Int {
	maxPrice, ok := new(big.Int).SetString(m.options.MaxPrice, 10)
	if !ok || maxPrice.Sign() == 0 {
		maxPrice = big.NewInt(1)
	}
	return new(big.Int).Mul(maxPrice, big.NewInt(1000))
}

// idRange returns [start, start+count) as big.Int tokenIds.
func idRange(start, count uint64) []*big.Int {
	ids := make([]*big.Int, 0, count)
	for i := uint64(0); i < count; i++ {
		ids = append(ids, new(big.Int).SetUint64(start+i))
	}
	return ids
}
