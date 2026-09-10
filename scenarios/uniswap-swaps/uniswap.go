package uniswapswaps

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"

	"github.com/ethpandaops/spamoor/scenarios/uniswap-swaps/contract"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/txtypes"
)

type UniswapOptions struct {
	Version    uint64
	BaseFee    float64
	TipFee     float64
	BaseFeeWei string
	TipFeeWei  string
	// PairCount is the number of DAI tokens to deploy; each one is paired with
	// the shared quote token on both factories.
	PairCount uint64
	// QuoteLiquidityPerPool is the quote token reserve seeded into every pair /
	// pool; the DAI side is QuoteLiquidityPerPool * TokensPerQuote.
	QuoteLiquidityPerPool *big.Int
	// TokensPerQuote is the initial price: DAI tokens per quote token.
	TokensPerQuote uint64
	// QuoteFunding is the quote token amount minted to each child wallet at
	// startup and whenever a wallet cannot afford a buy.
	QuoteFunding *big.Int
	FeeTier      uint64
	ClientGroup  string
}

// Uniswap owns the deployed contract set (v2 pairs or v3 pools) and the local
// per-wallet token balance cache used to decide swap directions without an RPC
// round trip per swap.
//
// Every pair trades a per-pair mock DAI token against one shared mock quote
// token. Both are ERC20s with a public mint, so no ETH capital beyond gas is
// needed: pools are seeded by minting and child wallets mint their own quote
// tokens.
type Uniswap struct {
	ctx            context.Context
	walletPool     *spamoor.WalletPool
	deploymentInfo *DeploymentInfo
	v3Deployment   *V3DeploymentInfo
	logger         *logrus.Entry
	options        UniswapOptions

	// local cache of token balances: wallet -> token -> balance
	tokenBalances      map[common.Address]map[common.Address]*big.Int
	tokenBalancesMutex sync.RWMutex

	// contract instances bound to the static call client
	RouterA *contract.UniswapV2Router02 // v2 only
	RouterB *contract.UniswapV2Router02 // v2 only
	Quote   *contract.Dai
}

// daiAddrs returns the per-pair DAI token addresses of the active deployment.
func (u *Uniswap) daiAddrs() []common.Address {
	if u.options.Version == 3 {
		addrs := make([]common.Address, 0, len(u.v3Deployment.Pools))
		for _, pool := range u.v3Deployment.Pools {
			addrs = append(addrs, pool.DaiAddr)
		}
		return addrs
	}
	addrs := make([]common.Address, 0, len(u.deploymentInfo.Pairs))
	for _, pair := range u.deploymentInfo.Pairs {
		addrs = append(addrs, pair.DaiAddr)
	}
	return addrs
}

// quoteAddr returns the shared quote token address of the active deployment.
func (u *Uniswap) quoteAddr() common.Address {
	if u.options.Version == 3 {
		return u.v3Deployment.QuoteAddr
	}
	return u.deploymentInfo.QuoteAddr
}

// allTokenAddrs returns every token a child wallet holds: all DAI tokens plus
// the quote token. Used by the generic balance/allowance setup phases.
func (u *Uniswap) allTokenAddrs() []common.Address {
	return append(u.daiAddrs(), u.quoteAddr())
}

// spenderAddrs returns the addresses child wallets must approve for token
// transfers: both v2 routers, or both v3 SwapRouters.
func (u *Uniswap) spenderAddrs() []common.Address {
	if u.options.Version == 3 {
		return []common.Address{u.v3Deployment.RouterAAddr, u.v3Deployment.RouterBAddr}
	}
	return []common.Address{u.deploymentInfo.UniswapRouterAAddr, u.deploymentInfo.UniswapRouterBAddr}
}

func NewUniswap(ctx context.Context, walletPool *spamoor.WalletPool, logger *logrus.Entry, options UniswapOptions) *Uniswap {
	return &Uniswap{
		ctx:           ctx,
		walletPool:    walletPool,
		logger:        logger,
		options:       options,
		tokenBalances: make(map[common.Address]map[common.Address]*big.Int),
	}
}

// staticCallClient returns the client used for eth_calls and for binding the
// reusable contract instances.
func (u *Uniswap) staticCallClient() (*spamoor.Client, error) {
	client := u.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0),
		spamoor.WithoutBuilder(), // avoid using builders for eth_calls
	)
	if client == nil {
		return nil, fmt.Errorf("no client available")
	}
	u.logger.Infof("Using client for static calls: %s", client.GetName())
	return client, nil
}

// InitializeContracts binds the deployed v2 contract instances to the static
// call client and stores the deployment for the swap phase.
func (u *Uniswap) InitializeContracts(deploymentInfo *DeploymentInfo) error {
	u.deploymentInfo = deploymentInfo

	client, err := u.staticCallClient()
	if err != nil {
		return err
	}

	u.RouterA, err = contract.NewUniswapV2Router02(deploymentInfo.UniswapRouterAAddr, client.GetEthClient())
	if err != nil {
		return fmt.Errorf("could not initialize router A: %w", err)
	}
	u.RouterB, err = contract.NewUniswapV2Router02(deploymentInfo.UniswapRouterBAddr, client.GetEthClient())
	if err != nil {
		return fmt.Errorf("could not initialize router B: %w", err)
	}
	u.Quote, err = contract.NewDai(deploymentInfo.QuoteAddr, client.GetEthClient())
	if err != nil {
		return fmt.Errorf("could not initialize quote token: %w", err)
	}

	return nil
}

// InitializeTokenBalances reads the DAI and quote token balances of all child
// wallets into the local cache.
func (u *Uniswap) InitializeTokenBalances() {
	// Initialize the 2D map
	u.tokenBalances = make(map[common.Address]map[common.Address]*big.Int)
	u.tokenBalancesMutex = sync.RWMutex{}

	// Get all wallets
	wallets := u.walletPool.GetAllWallets()
	tokenAddrs := u.allTokenAddrs()

	// Read balances for each wallet in parallel across clients. Doing this
	// serially over hundreds of wallets is hundreds of blocking RPC calls; the
	// context-aware CallOpts also let a UI stop cancel the in-flight reads.
	sem := make(chan struct{}, u.setupConcurrency())
	var wg sync.WaitGroup

	for idx, wallet := range wallets {
		if u.ctx.Err() != nil {
			break
		}
		wg.Add(1)
		go func(idx int, wallet *spamoor.Wallet) {
			defer wg.Done()
			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
			case <-u.ctx.Done():
				return
			}

			walletAddr := wallet.GetAddress()
			rclient := u.walletPool.GetClient(
				spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, idx),
				spamoor.WithClientGroup(u.options.ClientGroup),
				spamoor.WithoutBuilder(),
			)
			if rclient == nil {
				return
			}
			callOpts := &bind.CallOpts{Context: u.ctx}

			balances := make(map[common.Address]*big.Int, len(tokenAddrs))
			for _, tokenAddr := range tokenAddrs {
				token, err := contract.NewDai(tokenAddr, rclient.GetEthClient())
				if err != nil {
					u.logger.Errorf("could not bind token %v: %v", tokenAddr, err)
					continue
				}
				balance, err := token.BalanceOf(callOpts, walletAddr)
				if err != nil {
					u.logger.Errorf("could not get token balance for %v: %v", walletAddr, err)
					continue
				}
				balances[tokenAddr] = balance
			}

			u.tokenBalancesMutex.Lock()
			u.tokenBalances[walletAddr] = balances
			u.tokenBalancesMutex.Unlock()
		}(idx, wallet)
	}
	wg.Wait()
}

// GetTokenBalance returns the cached balance of a token for a wallet.
func (u *Uniswap) GetTokenBalance(walletAddr common.Address, tokenAddr common.Address) *big.Int {
	u.tokenBalancesMutex.RLock()
	defer u.tokenBalancesMutex.RUnlock()

	walletBalances, exists := u.tokenBalances[walletAddr]
	if !exists {
		return big.NewInt(0)
	}

	balance, exists := walletBalances[tokenAddr]
	if !exists {
		return big.NewInt(0)
	}
	return balance
}

// UpdateTokenBalance overwrites the cached balance of a token for a wallet.
func (u *Uniswap) UpdateTokenBalance(walletAddr common.Address, tokenAddr common.Address, newBalance *big.Int) {
	u.tokenBalancesMutex.Lock()
	defer u.tokenBalancesMutex.Unlock()

	// Ensure the wallet map exists
	if _, exists := u.tokenBalances[walletAddr]; !exists {
		u.tokenBalances[walletAddr] = make(map[common.Address]*big.Int)
	}

	u.tokenBalances[walletAddr][tokenAddr] = newBalance
}

// approvalGasLimit is the static gas limit for ERC20 approve txs. Under the
// Amsterdam fee schedule a fresh allowance slot makes approve cost ~128k; this
// keeps headroom. It is deliberately static (not estimated) so that setting
// allowances for hundreds of wallets needs no per-tx eth_estimateGas round trip.
const approvalGasLimit = 250000

// mintGasLimit is the static gas limit for quote token mint txs. A mint to a
// wallet without a balance yet creates one fresh storage slot (like approve)
// and updates the total supply, so the same headroom applies.
const mintGasLimit = 250000

// setupConcurrency bounds the parallel per-wallet RPC fan-out used by the setup
// phases (balance reads, allowance checks). Sized to the number of healthy
// clients so the load spreads across nodes, capped to avoid overwhelming them.
func (u *Uniswap) setupConcurrency() int {
	n := len(u.walletPool.GetClientPool().GetAllGoodClients())
	return min(max(n, 1), 50)
}

// buildQuoteMintTx builds a tx in which the wallet mints QuoteFunding quote
// tokens to itself, and bumps the cached balance accordingly. The quote token
// is a mock with a public mint, so this is how child wallets are capitalized
// instead of receiving ETH.
func (u *Uniswap) buildQuoteMintTx(ctx context.Context, wallet *spamoor.Wallet, feeCap, tipCap *big.Int) (*txtypes.Transaction, error) {
	quoteAddr := u.quoteAddr()
	walletAddr := wallet.GetAddress()

	tx, err := wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       mintGasLimit,
		Value:     uint256.NewInt(0),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		return u.Quote.Mint(transactOpts, walletAddr, u.options.QuoteFunding)
	})
	if err != nil {
		return nil, fmt.Errorf("could not build quote mint tx: %w", err)
	}

	balance := u.GetTokenBalance(walletAddr, quoteAddr)
	u.UpdateTokenBalance(walletAddr, quoteAddr, new(big.Int).Add(balance, u.options.QuoteFunding))
	return tx, nil
}

// PrepareWallets gets every child wallet ready for swapping: it sets unlimited
// allowances for all tokens to the router(s) and mints the initial quote token
// funding to wallets holding less than that amount. Must run after
// InitializeTokenBalances, whose cache decides which wallets need funding.
func (u *Uniswap) PrepareWallets() error {
	u.logger.Infof("Preparing wallets (allowances + quote token funding)...")

	// Get all wallets
	wallets := u.walletPool.GetAllWallets()

	// Maximum uint256 value for unlimited allowance
	maxAllowance := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))

	// Get a client for fee calculation
	client := u.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0),
		spamoor.WithClientGroup(u.options.ClientGroup),
	)
	if client == nil {
		return fmt.Errorf("no client available")
	}

	baseFeeWei, tipFeeWei := spamoor.ResolveFees(u.options.BaseFee, u.options.TipFee, u.options.BaseFeeWei, u.options.TipFeeWei)
	feeCap, tipCap, err := u.walletPool.GetSuggestedFees(client, baseFeeWei, tipFeeWei)
	if err != nil {
		return fmt.Errorf("could not get tx fee: %v", err)
	}

	routers := u.spenderAddrs()
	tokenAddrs := u.allTokenAddrs()
	quoteAddr := u.quoteAddr()

	// Track all setup transactions
	var (
		setupTxs     []*txtypes.Transaction
		setupWallets []*spamoor.Wallet
		mintCount    int
		mu           sync.Mutex
		wg           sync.WaitGroup
	)

	addSetupTx := func(wallet *spamoor.Wallet, tx *txtypes.Transaction) {
		mu.Lock()
		setupTxs = append(setupTxs, tx)
		setupWallets = append(setupWallets, wallet)
		mu.Unlock()
	}

	// Check allowances and build approval txs in parallel across clients. For N
	// wallets this is up to 2*(pairs+1)*N allowance reads (every token × router
	// A+B); doing them serially on one client blocks the scenario for minutes at
	// large wallet counts. The context-aware CallOpts also let a UI stop actually
	// cancel the in-flight reads.
	sem := make(chan struct{}, u.setupConcurrency())

	buildApproval := func(wallet *spamoor.Wallet, approve func(*bind.TransactOpts) (*types.Transaction, error)) {
		// Static gas: approve is uniform, so estimating each one would just add a
		// redundant round trip per wallet.
		approveTx, err := wallet.BuildBoundTx(u.ctx, &txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Gas:       approvalGasLimit,
			Value:     uint256.NewInt(0),
		}, approve)
		if err != nil {
			u.logger.Errorf("could not build approval tx for %v: %v", wallet.GetAddress(), err)
			return
		}
		addSetupTx(wallet, approveTx)
	}

	for idx, wallet := range wallets {
		if u.ctx.Err() != nil {
			break
		}
		wg.Add(1)
		go func(idx int, wallet *spamoor.Wallet) {
			defer wg.Done()
			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
			case <-u.ctx.Done():
				return
			}

			rclient := u.walletPool.GetClient(
				spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, idx),
				spamoor.WithClientGroup(u.options.ClientGroup),
				spamoor.WithoutBuilder(),
			)
			if rclient == nil {
				rclient = client
			}
			callOpts := &bind.CallOpts{Context: u.ctx}

			for _, tokenAddr := range tokenAddrs {
				token, err := contract.NewDai(tokenAddr, rclient.GetEthClient())
				if err != nil {
					u.logger.Errorf("could not bind token %v: %v", tokenAddr, err)
					continue
				}
				for _, router := range routers {
					allowance, err := token.Allowance(callOpts, wallet.GetAddress(), router)
					if err != nil {
						u.logger.Errorf("could not check allowance for %v: %v", wallet.GetAddress(), err)
						continue
					}
					if allowance.Cmp(maxAllowance) >= 0 {
						continue
					}
					buildApproval(wallet, func(opts *bind.TransactOpts) (*types.Transaction, error) {
						return token.Approve(opts, router, maxAllowance)
					})
				}
			}

			// Initial quote token funding for wallets below the funding amount.
			if u.GetTokenBalance(wallet.GetAddress(), quoteAddr).Cmp(u.options.QuoteFunding) < 0 {
				mintTx, err := u.buildQuoteMintTx(u.ctx, wallet, feeCap, tipCap)
				if err != nil {
					u.logger.Errorf("could not build quote mint tx for %v: %v", wallet.GetAddress(), err)
					return
				}
				mu.Lock()
				mintCount++
				mu.Unlock()
				addSetupTx(wallet, mintTx)
			}
		}(idx, wallet)
	}
	wg.Wait()

	if u.ctx.Err() != nil {
		return u.ctx.Err()
	}

	// Send all setup transactions in parallel
	if len(setupTxs) > 0 {
		u.logger.Infof("Sending %d wallet setup transactions (%d approvals, %d quote mints)...", len(setupTxs), len(setupTxs)-mintCount, mintCount)

		// Reuse the wait group (back to zero after the build phase) to track sends.
		// Send each transaction to a different client
		for i, tx := range setupTxs {
			// Get a different client for each transaction
			txClient := u.walletPool.GetClient(
				spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, i),
				spamoor.WithClientGroup(u.options.ClientGroup),
			)
			if txClient == nil {
				txClient = client
			}

			wg.Add(1)

			go func(tx *txtypes.Transaction, client *spamoor.Client, wallet *spamoor.Wallet) {
				u.walletPool.GetTxPool().SendTransaction(u.ctx, wallet, tx, &spamoor.SendTransactionOptions{
					Client:      client,
					ClientGroup: u.options.ClientGroup,
					Rebroadcast: true,
					OnComplete: func(tx *txtypes.Transaction, receipt *txtypes.Receipt, err error) {
						if err != nil {
							u.logger.Errorf("wallet setup tx failed: %v", err)
						}
						wg.Done()
					},
				})
			}(tx, txClient, setupWallets[i])
		}

		// Wait for all transactions to be sent
		wg.Wait()
		u.logger.Infof("All wallet setup transactions sent")
	} else {
		u.logger.Infof("No wallet setup transactions needed (allowances and funding already in place)")
	}

	return nil
}
