package uniswapswaps

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethpandaops/spamoor/scenarios/uniswap-swaps/contract"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
)

type UniswapOptions struct {
	Version             uint64
	BaseFee             float64
	TipFee              float64
	BaseFeeWei          string
	TipFeeWei           string
	DaiPairs            uint64
	EthLiquidityPerPair *uint256.Int
	DaiLiquidityFactor  uint64
	FeeTier             uint64
	ClientGroup         string
}

type Uniswap struct {
	ctx            context.Context
	walletPool     *spamoor.WalletPool
	deploymentInfo *DeploymentInfo
	v3Deployment   *V3DeploymentInfo
	logger         *logrus.Entry
	options        UniswapOptions

	// local cache of token balances
	tokenBalances      map[common.Address]map[common.Address]*big.Int
	tokenBalancesMutex sync.RWMutex

	// v2 contract instances
	RouterA *contract.UniswapV2Router02
	RouterB *contract.UniswapV2Router02
	Weth    *contract.WETH9
	Tokens  map[common.Address]*contract.Dai
}

// tokenAddrs returns the list of DAI token addresses across all deployed pairs
// or pools, used by the generic balance/allowance setup phases.
func (u *Uniswap) tokenAddrs() []common.Address {
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

// wethAddr returns the WETH9 address of the active deployment.
func (u *Uniswap) wethAddr() common.Address {
	if u.options.Version == 3 {
		return u.v3Deployment.Weth9Addr
	}
	return u.deploymentInfo.Weth9Addr
}

// spenderAddrs returns the addresses child wallets must approve for token
// transfers: both v2 routers, or the single v3 SwapRouter.
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

// Initialize contract instances to reuse
func (u *Uniswap) InitializeContracts(deploymentInfo *DeploymentInfo) error {
	u.deploymentInfo = deploymentInfo

	client := u.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0),
		spamoor.WithoutBuilder(), // avoid using builders for eth_calls
	)
	if client == nil {
		return fmt.Errorf("no client available")
	}

	u.logger.Infof("Using client for static calls: %s", client.GetName())

	// Initialize router A
	routerA, err := contract.NewUniswapV2Router02(u.deploymentInfo.UniswapRouterAAddr, client.GetEthClient())
	if err != nil {
		return fmt.Errorf("could not initialize router A: %w", err)
	}
	u.RouterA = routerA

	// Initialize router B
	routerB, err := contract.NewUniswapV2Router02(u.deploymentInfo.UniswapRouterBAddr, client.GetEthClient())
	if err != nil {
		return fmt.Errorf("could not initialize router B: %w", err)
	}
	u.RouterB = routerB

	// Initialize WETH9
	weth, err := contract.NewWETH9(u.deploymentInfo.Weth9Addr, client.GetEthClient())
	if err != nil {
		return fmt.Errorf("could not initialize WETH9: %w", err)
	}
	u.Weth = weth

	// Initialize token contracts
	u.Tokens = make(map[common.Address]*contract.Dai)
	for _, pair := range u.deploymentInfo.Pairs {
		token, err := contract.NewDai(pair.DaiAddr, client.GetEthClient())
		if err != nil {
			return fmt.Errorf("could not initialize token %v: %w", pair.DaiAddr, err)
		}
		u.Tokens[pair.DaiAddr] = token
	}

	return nil
}

// Initialize token balances for all wallets
func (u *Uniswap) InitializeTokenBalances() {
	// Initialize the 2D map
	u.tokenBalances = make(map[common.Address]map[common.Address]*big.Int)
	u.tokenBalancesMutex = sync.RWMutex{}

	// Get all wallets
	wallets := u.walletPool.GetAllWallets()

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

			tokenAddrs := u.tokenAddrs()
			wethAddr := u.wethAddr()
			balances := make(map[common.Address]*big.Int, len(tokenAddrs)+1)
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

			if weth, err := contract.NewWETH9(wethAddr, rclient.GetEthClient()); err != nil {
				u.logger.Errorf("could not bind WETH9: %v", err)
			} else if wethBalance, err := weth.BalanceOf(callOpts, walletAddr); err != nil {
				u.logger.Errorf("could not get WETH balance for %v: %v", walletAddr, err)
			} else {
				balances[wethAddr] = wethBalance
			}

			u.tokenBalancesMutex.Lock()
			u.tokenBalances[walletAddr] = balances
			u.tokenBalancesMutex.Unlock()
		}(idx, wallet)
	}
	wg.Wait()
}

// Get DAI balance from local cache
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

// Update DAI balance in local cache
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

// setupConcurrency bounds the parallel per-wallet RPC fan-out used by the setup
// phases (balance reads, allowance checks). Sized to the number of healthy
// clients so the load spreads across nodes, capped to avoid overwhelming them.
func (u *Uniswap) setupConcurrency() int {
	n := len(u.walletPool.GetClientPool().GetAllGoodClients())
	return min(max(n, 1), 50)
}

// Set unlimited allowances for all wallets to both routers
func (u *Uniswap) SetUnlimitedAllowances() error {
	u.logger.Infof("Setting unlimited allowances for all wallets...")

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
	tokenAddrs := u.tokenAddrs()
	wethAddr := u.wethAddr()

	// Track all approval transactions
	var (
		approvalTxs     []*types.Transaction
		approvalWallets []*spamoor.Wallet
		mu              sync.Mutex
		wg              sync.WaitGroup
	)

	// Check allowances and build approval txs in parallel across clients. For N
	// wallets this is up to 4*N allowance reads (DAI+WETH × router A+B); doing
	// them serially on one client blocks the scenario for minutes at large wallet
	// counts. The context-aware CallOpts also let a UI stop actually cancel the
	// in-flight reads.
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
		mu.Lock()
		approvalTxs = append(approvalTxs, approveTx)
		approvalWallets = append(approvalWallets, wallet)
		mu.Unlock()
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

			// DAI tokens
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

			// WETH
			weth, err := contract.NewWETH9(wethAddr, rclient.GetEthClient())
			if err != nil {
				u.logger.Errorf("could not bind WETH9: %v", err)
				return
			}
			for _, router := range routers {
				allowance, err := weth.Allowance(callOpts, wallet.GetAddress(), router)
				if err != nil {
					u.logger.Errorf("could not check WETH allowance for %v: %v", wallet.GetAddress(), err)
					continue
				}
				if allowance.Cmp(maxAllowance) >= 0 {
					continue
				}
				buildApproval(wallet, func(opts *bind.TransactOpts) (*types.Transaction, error) {
					return weth.Approve(opts, router, maxAllowance)
				})
			}
		}(idx, wallet)
	}
	wg.Wait()

	if u.ctx.Err() != nil {
		return u.ctx.Err()
	}

	// Send all approval transactions in parallel
	if len(approvalTxs) > 0 {
		u.logger.Infof("Sending %d approval transactions...", len(approvalTxs))

		// Reuse the wait group (back to zero after the build phase) to track sends.
		// Send each transaction to a different client
		for i, tx := range approvalTxs {
			// Get a different client for each transaction
			txClient := u.walletPool.GetClient(
				spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, i),
				spamoor.WithClientGroup(u.options.ClientGroup),
			)
			if txClient == nil {
				txClient = client
			}

			wg.Add(1)

			go func(tx *types.Transaction, client *spamoor.Client, wallet *spamoor.Wallet) {
				u.walletPool.GetTxPool().SendTransaction(u.ctx, wallet, tx, &spamoor.SendTransactionOptions{
					Client:      client,
					ClientGroup: u.options.ClientGroup,
					Rebroadcast: true,
					OnComplete: func(tx *types.Transaction, receipt *types.Receipt, err error) {
						if err != nil {
							u.logger.Errorf("approval tx failed: %v", err)
						}
						wg.Done()
					},
				})
			}(tx, txClient, approvalWallets[i])
		}

		// Wait for all transactions to be sent
		wg.Wait()
		u.logger.Infof("All approval transactions sent")
	} else {
		u.logger.Infof("No approval transactions needed (allowances already set)")
	}

	return nil
}
