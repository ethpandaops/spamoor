package curveswaps

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/scenarios/curve-swaps/contract"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
)

// coinsPerPool is the number of coins in each StableSwap pool. It is fixed at 3
// to match the deployed Curve 3pool StableSwap contract.
const coinsPerPool = 3

// approvalGasLimit is the static gas limit for the ERC20 approve and mint setup
// txs. A fresh allowance/balance slot under the Amsterdam fee schedule costs
// ~128k; this keeps headroom while avoiding a per-tx eth_estimateGas round trip.
const approvalGasLimit = 250000

type CurveOptions struct {
	BaseFee       float64
	TipFee        float64
	BaseFeeWei    string
	TipFeeWei     string
	PoolCount     uint64
	Amplification uint64
	Fee           uint64
	SeedAmount    *uint256.Int
	WalletFunding *uint256.Int
	ClientGroup   string
}

// CurvePoolInfo holds a single deployed StableSwap pool, its LP token, and its
// three coins.
type CurvePoolInfo struct {
	PoolAddr      common.Address
	Pool          *contract.StableSwap
	LpAddr        common.Address
	Lp            *contract.CurveToken
	Coins         [coinsPerPool]common.Address
	CoinContracts [coinsPerPool]*contract.MintableToken
}

// CurveDeploymentInfo holds the deployed StableSwap pools and the shared
// liquidity provider helper used to seed them.
type CurveDeploymentInfo struct {
	LiquidityProviderAddr common.Address
	LiquidityProvider     *contract.CurveLiquidityProvider
	Pools                 []CurvePoolInfo
}

type Curve struct {
	ctx        context.Context
	walletPool *spamoor.WalletPool
	logger     *logrus.Entry
	options    CurveOptions
	deployment *CurveDeploymentInfo

	// local cache of per-wallet token balances (wallet -> coin -> balance)
	tokenBalances      map[common.Address]map[common.Address]*big.Int
	tokenBalancesMutex sync.RWMutex
}

func NewCurve(ctx context.Context, walletPool *spamoor.WalletPool, logger *logrus.Entry, options CurveOptions) *Curve {
	return &Curve{
		ctx:           ctx,
		walletPool:    walletPool,
		logger:        logger,
		options:       options,
		tokenBalances: make(map[common.Address]map[common.Address]*big.Int),
	}
}

// allCoins returns every coin address across all deployed pools.
func (c *Curve) allCoins() []common.Address {
	addrs := make([]common.Address, 0, len(c.deployment.Pools)*coinsPerPool)
	for _, pool := range c.deployment.Pools {
		addrs = append(addrs, pool.Coins[:]...)
	}
	return addrs
}

// setupConcurrency bounds the parallel per-wallet RPC fan-out used by the setup
// phases (balance reads, allowance checks). Sized to the number of healthy
// clients so the load spreads across nodes, capped to avoid overwhelming them.
func (c *Curve) setupConcurrency() int {
	n := len(c.walletPool.GetClientPool().GetAllGoodClients())
	return min(max(n, 1), 50)
}

// DeployCurvePools deploys the liquidity provider helper, one LP token and three
// mintable coins per pool, and the StableSwap pools themselves. It then hands LP
// minting rights to each pool and seeds balanced liquidity.
func (c *Curve) DeployCurvePools() (*CurveDeploymentInfo, error) {
	client := c.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0),
		spamoor.WithClientGroup(c.options.ClientGroup),
	)
	if client == nil {
		return nil, scenario.ErrNoClients
	}

	deployerWallet := c.walletPool.GetWellKnownWallet("deployer")
	if deployerWallet == nil {
		return nil, scenario.ErrNoWallet
	}
	deployerAddr := deployerWallet.GetAddress()

	deployerSeed := [32]byte{}
	copy(deployerSeed[:], deployerAddr.Bytes())

	baseFeeWei, tipFeeWei := spamoor.ResolveFees(c.options.BaseFee, c.options.TipFee, c.options.BaseFeeWei, c.options.TipFeeWei)
	feeCap, tipCap, err := c.walletPool.GetSuggestedFees(client, baseFeeWei, tipFeeWei)
	if err != nil {
		return nil, fmt.Errorf("could not get tx fee: %w", err)
	}

	deploymentTxs := []*types.Transaction{}
	deployContract := func(metadata *bind.MetaData, global bool, salt uint32, params ...interface{}) (common.Address, error) {
		parsed, err := metadata.GetAbi()
		if err != nil {
			return common.Address{}, err
		}
		if parsed == nil {
			return common.Address{}, fmt.Errorf("GetABI returned nil")
		}

		initCodeBytes := common.FromHex(metadata.Bin)
		packed, err := parsed.Pack("", params...)
		if err != nil {
			return common.Address{}, err
		}
		initCodeBytes = append(initCodeBytes, packed...)

		seed := [32]byte{}
		if !global {
			copy(seed[:], deployerSeed[:])
		}
		if salt != 0 {
			binary.BigEndian.PutUint32(seed[28:], salt)
		}
		addr, tx, err := c.walletPool.GetDeploymentFactory().GetContractDeployment(c.ctx, initCodeBytes, seed, client, deployerWallet, feeCap, tipCap, false)
		if err != nil {
			return common.Address{}, err
		}
		if tx != nil {
			deploymentTxs = append(deploymentTxs, tx)
		}
		return addr, nil
	}

	info := &CurveDeploymentInfo{}

	// deploy the liquidity provider helper
	info.LiquidityProviderAddr, err = deployContract(contract.CurveLiquidityProviderMetaData, false, 0)
	if err != nil {
		return nil, fmt.Errorf("could not deploy curve liquidity provider: %w", err)
	}
	info.LiquidityProvider, err = contract.NewCurveLiquidityProvider(info.LiquidityProviderAddr, client.GetEthClient())
	if err != nil {
		return nil, fmt.Errorf("could not create instance of curve liquidity provider: %w", err)
	}

	amplification := new(big.Int).SetUint64(c.options.Amplification)
	swapFee := new(big.Int).SetUint64(c.options.Fee)
	adminFee := big.NewInt(0)

	// deploy coins, LP token and a StableSwap pool per configured pool
	for i := uint64(0); i < c.options.PoolCount; i++ {
		poolInfo := CurvePoolInfo{}

		for j := 0; j < coinsPerPool; j++ {
			name := fmt.Sprintf("Curve Mock Stablecoin %d-%d", i, j)
			symbol := fmt.Sprintf("USD%d%d", i, j)
			salt := uint32(i*coinsPerPool + uint64(j) + 1)
			poolInfo.Coins[j], err = deployContract(contract.MintableTokenMetaData, true, salt, name, symbol)
			if err != nil {
				return nil, fmt.Errorf("could not deploy coin %d-%d: %w", i, j, err)
			}
			poolInfo.CoinContracts[j], err = contract.NewMintableToken(poolInfo.Coins[j], client.GetEthClient())
			if err != nil {
				return nil, fmt.Errorf("could not create instance of coin %d-%d: %w", i, j, err)
			}
		}

		// LP token: minter is set to the deployer so it can later hand minting
		// rights to the pool (the pool address is not known when this is packed).
		lpName := fmt.Sprintf("Curve Mock LP %d", i)
		lpSymbol := fmt.Sprintf("crvMock%d", i)
		poolInfo.LpAddr, err = deployContract(contract.CurveTokenMetaData, false, uint32(1000+i), lpName, lpSymbol, deployerAddr)
		if err != nil {
			return nil, fmt.Errorf("could not deploy LP token %d: %w", i, err)
		}
		poolInfo.Lp, err = contract.NewCurveToken(poolInfo.LpAddr, client.GetEthClient())
		if err != nil {
			return nil, fmt.Errorf("could not create instance of LP token %d: %w", i, err)
		}

		poolInfo.PoolAddr, err = deployContract(contract.StableSwapMetaData, false, uint32(i+1), deployerAddr, poolInfo.Coins, poolInfo.LpAddr, amplification, swapFee, adminFee)
		if err != nil {
			return nil, fmt.Errorf("could not deploy stableswap pool %d: %w", i, err)
		}
		poolInfo.Pool, err = contract.NewStableSwap(poolInfo.PoolAddr, client.GetEthClient())
		if err != nil {
			return nil, fmt.Errorf("could not create instance of stableswap pool %d: %w", i, err)
		}

		info.Pools = append(info.Pools, poolInfo)
	}

	// submit & await deployment batch
	if err := c.sendBatch(deployerWallet, client, deploymentTxs, "deploying contracts"); err != nil {
		return nil, err
	}

	c.deployment = info

	// hand LP minting rights to each pool (deployer is the current minter).
	minterTxs := []*types.Transaction{}
	for _, poolInfo := range info.Pools {
		poolInfo := poolInfo
		tx, err := deployerWallet.BuildBoundTxWithEstimate(c.ctx, client, c.walletPool.GetTxPool(), &txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Value:     uint256.NewInt(0),
		}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
			return poolInfo.Lp.SetMinter(transactOpts, poolInfo.PoolAddr)
		})
		if err != nil {
			return nil, fmt.Errorf("could not build set_minter tx for pool %v: %w", poolInfo.PoolAddr.Hex(), err)
		}
		minterTxs = append(minterTxs, tx)
	}
	if err := c.sendBatch(deployerWallet, client, minterTxs, "setting LP minters"); err != nil {
		return nil, err
	}

	// seed balanced liquidity into every pool (one estimable tx per pool: the
	// helper mints, approves and adds liquidity atomically).
	seedTxs := []*types.Transaction{}
	seedAmount := c.options.SeedAmount.ToBig()
	for _, poolInfo := range info.Pools {
		poolInfo := poolInfo
		tx, err := deployerWallet.BuildBoundTxWithEstimate(c.ctx, client, c.walletPool.GetTxPool(), &txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Value:     uint256.NewInt(0),
		}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
			return info.LiquidityProvider.SeedLiquidity(transactOpts, poolInfo.PoolAddr, poolInfo.Coins, seedAmount)
		})
		if err != nil {
			return nil, fmt.Errorf("could not build seed liquidity tx for pool %v: %w", poolInfo.PoolAddr.Hex(), err)
		}
		seedTxs = append(seedTxs, tx)
	}
	if err := c.sendBatch(deployerWallet, client, seedTxs, "seeding liquidity"); err != nil {
		return nil, err
	}

	return info, nil
}

// InitializeContracts binds the deployed contract instances to the static call
// client and stores the deployment for the swap phase.
func (c *Curve) InitializeContracts(info *CurveDeploymentInfo) error {
	client := c.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0),
		spamoor.WithoutBuilder(), // avoid using builders for eth_calls
	)
	if client == nil {
		return scenario.ErrNoClients
	}

	c.logger.Infof("Using client for static calls: %s", client.GetName())

	lp, err := contract.NewCurveLiquidityProvider(info.LiquidityProviderAddr, client.GetEthClient())
	if err != nil {
		return fmt.Errorf("could not initialize curve liquidity provider: %w", err)
	}
	info.LiquidityProvider = lp

	for i := range info.Pools {
		info.Pools[i].Pool, err = contract.NewStableSwap(info.Pools[i].PoolAddr, client.GetEthClient())
		if err != nil {
			return fmt.Errorf("could not initialize stableswap pool %d: %w", i, err)
		}
		info.Pools[i].Lp, err = contract.NewCurveToken(info.Pools[i].LpAddr, client.GetEthClient())
		if err != nil {
			return fmt.Errorf("could not initialize LP token %d: %w", i, err)
		}
		for j := 0; j < coinsPerPool; j++ {
			info.Pools[i].CoinContracts[j], err = contract.NewMintableToken(info.Pools[i].Coins[j], client.GetEthClient())
			if err != nil {
				return fmt.Errorf("could not initialize coin %d-%d: %w", i, j, err)
			}
		}
	}

	c.deployment = info
	return nil
}

// FundAndApproveWallets mints the initial coin balance to every child wallet and
// sets unlimited allowances from each wallet to the pools so swaps can pull the
// input coin. Each wallet is funded with the configured amount of the first coin
// of every pool; the balancing logic in the swap builder spreads that into the
// other coins over time.
func (c *Curve) FundAndApproveWallets() error {
	c.logger.Infof("Funding and approving wallets...")

	wallets := c.walletPool.GetAllWallets()
	maxAllowance := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))
	fundingAmount := c.options.WalletFunding.ToBig()

	client := c.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0),
		spamoor.WithClientGroup(c.options.ClientGroup),
	)
	if client == nil {
		return scenario.ErrNoClients
	}

	baseFeeWei, tipFeeWei := spamoor.ResolveFees(c.options.BaseFee, c.options.TipFee, c.options.BaseFeeWei, c.options.TipFeeWei)
	feeCap, tipCap, err := c.walletPool.GetSuggestedFees(client, baseFeeWei, tipFeeWei)
	if err != nil {
		return fmt.Errorf("could not get tx fee: %w", err)
	}

	var (
		setupTxs     []*types.Transaction
		setupWallets []*spamoor.Wallet
		mu           sync.Mutex
		wg           sync.WaitGroup
	)
	sem := make(chan struct{}, c.setupConcurrency())

	buildTx := func(wallet *spamoor.Wallet, build func(*bind.TransactOpts) (*types.Transaction, error)) {
		tx, err := wallet.BuildBoundTx(c.ctx, &txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Gas:       approvalGasLimit,
			Value:     uint256.NewInt(0),
		}, build)
		if err != nil {
			c.logger.Errorf("could not build setup tx for %v: %v", wallet.GetAddress(), err)
			return
		}
		mu.Lock()
		setupTxs = append(setupTxs, tx)
		setupWallets = append(setupWallets, wallet)
		mu.Unlock()
	}

	for idx, wallet := range wallets {
		if c.ctx.Err() != nil {
			break
		}
		wg.Add(1)
		go func(idx int, wallet *spamoor.Wallet) {
			defer wg.Done()
			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
			case <-c.ctx.Done():
				return
			}

			rclient := c.walletPool.GetClient(
				spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, idx),
				spamoor.WithClientGroup(c.options.ClientGroup),
				spamoor.WithoutBuilder(),
			)
			if rclient == nil {
				rclient = client
			}
			callOpts := &bind.CallOpts{Context: c.ctx}

			for _, poolInfo := range c.deployment.Pools {
				for j := 0; j < coinsPerPool; j++ {
					coinAddr := poolInfo.Coins[j]
					token, err := contract.NewMintableToken(coinAddr, rclient.GetEthClient())
					if err != nil {
						c.logger.Errorf("could not bind coin %v: %v", coinAddr, err)
						continue
					}

					// mint the initial balance into the first coin of each pool
					if j == 0 {
						balance, err := token.BalanceOf(callOpts, wallet.GetAddress())
						if err != nil {
							c.logger.Errorf("could not read coin balance for %v: %v", wallet.GetAddress(), err)
						} else if balance.Cmp(fundingAmount) < 0 {
							buildTx(wallet, func(opts *bind.TransactOpts) (*types.Transaction, error) {
								return token.Mint(opts, wallet.GetAddress(), fundingAmount)
							})
						}
					}

					// approve the pool to spend this coin
					allowance, err := token.Allowance(callOpts, wallet.GetAddress(), poolInfo.PoolAddr)
					if err != nil {
						c.logger.Errorf("could not check allowance for %v: %v", wallet.GetAddress(), err)
						continue
					}
					if allowance.Cmp(maxAllowance) >= 0 {
						continue
					}
					buildTx(wallet, func(opts *bind.TransactOpts) (*types.Transaction, error) {
						return token.Approve(opts, poolInfo.PoolAddr, maxAllowance)
					})
				}
			}
		}(idx, wallet)
	}
	wg.Wait()

	if c.ctx.Err() != nil {
		return c.ctx.Err()
	}

	if len(setupTxs) == 0 {
		c.logger.Infof("No funding/approval transactions needed")
		return nil
	}

	c.logger.Infof("Sending %d funding/approval transactions...", len(setupTxs))
	for i, tx := range setupTxs {
		txClient := c.walletPool.GetClient(
			spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, i),
			spamoor.WithClientGroup(c.options.ClientGroup),
		)
		if txClient == nil {
			txClient = client
		}

		wg.Add(1)
		go func(tx *types.Transaction, client *spamoor.Client, wallet *spamoor.Wallet) {
			c.walletPool.GetTxPool().SendTransaction(c.ctx, wallet, tx, &spamoor.SendTransactionOptions{
				Client:      client,
				ClientGroup: c.options.ClientGroup,
				Rebroadcast: true,
				OnComplete: func(tx *types.Transaction, receipt *types.Receipt, err error) {
					if err != nil {
						c.logger.Errorf("funding/approval tx failed: %v", err)
					}
					wg.Done()
				},
			})
		}(tx, txClient, setupWallets[i])
	}
	wg.Wait()
	c.logger.Infof("All funding/approval transactions sent")

	return nil
}

// InitializeTokenBalances reads each wallet's balance of every coin into the
// local cache. Run after funding so the swap builder starts from real balances.
func (c *Curve) InitializeTokenBalances() {
	c.tokenBalances = make(map[common.Address]map[common.Address]*big.Int)
	c.tokenBalancesMutex = sync.RWMutex{}

	wallets := c.walletPool.GetAllWallets()
	coinAddrs := c.allCoins()

	sem := make(chan struct{}, c.setupConcurrency())
	var wg sync.WaitGroup

	for idx, wallet := range wallets {
		if c.ctx.Err() != nil {
			break
		}
		wg.Add(1)
		go func(idx int, wallet *spamoor.Wallet) {
			defer wg.Done()
			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
			case <-c.ctx.Done():
				return
			}

			walletAddr := wallet.GetAddress()
			rclient := c.walletPool.GetClient(
				spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, idx),
				spamoor.WithClientGroup(c.options.ClientGroup),
				spamoor.WithoutBuilder(),
			)
			if rclient == nil {
				return
			}
			callOpts := &bind.CallOpts{Context: c.ctx}

			balances := make(map[common.Address]*big.Int, len(coinAddrs))
			for _, coinAddr := range coinAddrs {
				token, err := contract.NewMintableToken(coinAddr, rclient.GetEthClient())
				if err != nil {
					c.logger.Errorf("could not bind coin %v: %v", coinAddr, err)
					continue
				}
				balance, err := token.BalanceOf(callOpts, walletAddr)
				if err != nil {
					c.logger.Errorf("could not get coin balance for %v: %v", walletAddr, err)
					continue
				}
				balances[coinAddr] = balance
			}

			c.tokenBalancesMutex.Lock()
			c.tokenBalances[walletAddr] = balances
			c.tokenBalancesMutex.Unlock()
		}(idx, wallet)
	}
	wg.Wait()
}

// GetTokenBalance returns the cached balance of a coin for a wallet.
func (c *Curve) GetTokenBalance(walletAddr common.Address, coinAddr common.Address) *big.Int {
	c.tokenBalancesMutex.RLock()
	defer c.tokenBalancesMutex.RUnlock()

	walletBalances, exists := c.tokenBalances[walletAddr]
	if !exists {
		return big.NewInt(0)
	}
	balance, exists := walletBalances[coinAddr]
	if !exists {
		return big.NewInt(0)
	}
	return balance
}

// UpdateTokenBalance updates the cached balance of a coin for a wallet.
func (c *Curve) UpdateTokenBalance(walletAddr common.Address, coinAddr common.Address, newBalance *big.Int) {
	c.tokenBalancesMutex.Lock()
	defer c.tokenBalancesMutex.Unlock()

	if _, exists := c.tokenBalances[walletAddr]; !exists {
		c.tokenBalances[walletAddr] = make(map[common.Address]*big.Int)
	}
	c.tokenBalances[walletAddr][coinAddr] = newBalance
}

// sendBatch submits a batch of transactions from a single wallet and waits for
// them to confirm, logging progress. It is a no-op for an empty batch.
func (c *Curve) sendBatch(wallet *spamoor.Wallet, client *spamoor.Client, txs []*types.Transaction, action string) error {
	if len(txs) == 0 {
		return nil
	}

	_, err := c.walletPool.GetTxPool().SendTransactionBatch(c.ctx, wallet, txs, &spamoor.BatchOptions{
		SendTransactionOptions: spamoor.SendTransactionOptions{
			Client:      client,
			ClientGroup: c.options.ClientGroup,
		},
		MaxRetries:   3,
		PendingLimit: 10,
		LogFn: func(confirmedCount int, totalCount int) {
			c.logger.Infof("%s... (%v/%v)", action, confirmedCount, totalCount)
		},
		LogInterval: 10,
	})
	if err != nil {
		return fmt.Errorf("could not %s: %w", action, err)
	}
	c.logger.Infof("%s complete. (%v/%v)", action, len(txs), len(txs))
	return nil
}
