package uniswapswaps

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/scenarios/uniswap-swaps/contract"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/holiman/uint256"
)

// V3DeploymentInfo holds the deployed Uniswap v3 contract set for the scenario.
// Two factories (each with its own SwapRouter) are deployed so that every DAI
// instance gets a separate pool per factory at the same fee tier, mirroring the
// two-factory layout of the v2 path.
type V3DeploymentInfo struct {
	Weth9Addr             common.Address
	Weth9                 *contract.WETH9
	FactoryAAddr          common.Address
	FactoryA              *contract.UniswapV3Factory
	FactoryBAddr          common.Address
	FactoryB              *contract.UniswapV3Factory
	RouterAAddr           common.Address
	RouterA               *contract.SwapRouter
	RouterBAddr           common.Address
	RouterB               *contract.SwapRouter
	LiquidityProviderAddr common.Address
	LiquidityProvider     *contract.V3LiquidityProvider
	Fee                   *big.Int
	TickSpacing           int64
	Pools                 []V3PoolDeploymentInfo
}

type V3PoolDeploymentInfo struct {
	DaiAddr      common.Address
	Dai          *contract.Dai
	WethIsToken0 bool
	PoolAAddr    common.Address
	PoolA        *contract.UniswapV3Pool
	PoolBAddr    common.Address
	PoolB        *contract.UniswapV3Pool
}

// liquidityBudgetBps applies a small safety margin (0.1%) to the WETH liquidity
// budget when sizing the position, so the pool's round-up of owed amounts can
// never exceed the ETH value forwarded to the liquidity provider.
const liquidityBudgetBps = 9990

// DeployUniswapV3 deploys two canonical Uniswap v3 factories + SwapRouters, the
// custom liquidity provider, and one DAI token per configured pair. Each DAI
// gets a pool on both factories, which are then initialized and seeded with a
// full-range position.
func (u *Uniswap) DeployUniswapV3() (*V3DeploymentInfo, error) {
	client := u.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0),
		spamoor.WithClientGroup(u.options.ClientGroup),
	)
	if client == nil {
		return nil, scenario.ErrNoClients
	}

	deployerWallet := u.walletPool.GetWellKnownWallet("deployer")
	ownerWallet := u.walletPool.GetWellKnownWallet("owner")
	if deployerWallet == nil || ownerWallet == nil {
		return nil, scenario.ErrNoWallet
	}
	rootAddr := u.walletPool.GetRootWallet().GetWallet().GetAddress()

	deployerSeed := [32]byte{}
	copy(deployerSeed[:], deployerWallet.GetAddress().Bytes())

	baseFeeWei, tipFeeWei := spamoor.ResolveFees(u.options.BaseFee, u.options.TipFee, u.options.BaseFeeWei, u.options.TipFeeWei)
	feeCap, tipCap, err := u.walletPool.GetSuggestedFees(client, baseFeeWei, tipFeeWei)
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
		addr, tx, err := u.walletPool.GetDeploymentFactory().GetContractDeployment(u.ctx, initCodeBytes, seed, client, deployerWallet, feeCap, tipCap, false)
		if err != nil {
			return common.Address{}, err
		}
		if tx != nil {
			deploymentTxs = append(deploymentTxs, tx)
		}
		return addr, nil
	}

	info := &V3DeploymentInfo{
		Fee: new(big.Int).SetUint64(u.options.FeeTier),
	}

	// deploy WETH9
	info.Weth9Addr, err = deployContract(contract.WETH9MetaData, true, 0)
	if err != nil {
		return nil, fmt.Errorf("could not deploy WETH9: %w", err)
	}
	info.Weth9, err = contract.NewWETH9(info.Weth9Addr, client.GetEthClient())
	if err != nil {
		return nil, fmt.Errorf("could not create instance of WETH9: %w", err)
	}

	// deploy two v3 factories (identical bytecode -> distinct salts)
	info.FactoryAAddr, err = deployContract(contract.UniswapV3FactoryMetaData, false, 0)
	if err != nil {
		return nil, fmt.Errorf("could not deploy uniswap v3 factory A: %w", err)
	}
	info.FactoryA, err = contract.NewUniswapV3Factory(info.FactoryAAddr, client.GetEthClient())
	if err != nil {
		return nil, fmt.Errorf("could not create instance of uniswap v3 factory A: %w", err)
	}

	info.FactoryBAddr, err = deployContract(contract.UniswapV3FactoryMetaData, false, 1)
	if err != nil {
		return nil, fmt.Errorf("could not deploy uniswap v3 factory B: %w", err)
	}
	info.FactoryB, err = contract.NewUniswapV3Factory(info.FactoryBAddr, client.GetEthClient())
	if err != nil {
		return nil, fmt.Errorf("could not create instance of uniswap v3 factory B: %w", err)
	}

	// deploy a swap router per factory
	info.RouterAAddr, err = deployContract(contract.SwapRouterMetaData, false, 0, info.FactoryAAddr, info.Weth9Addr)
	if err != nil {
		return nil, fmt.Errorf("could not deploy swap router A: %w", err)
	}
	info.RouterA, err = contract.NewSwapRouter(info.RouterAAddr, client.GetEthClient())
	if err != nil {
		return nil, fmt.Errorf("could not create instance of swap router A: %w", err)
	}

	info.RouterBAddr, err = deployContract(contract.SwapRouterMetaData, false, 1, info.FactoryBAddr, info.Weth9Addr)
	if err != nil {
		return nil, fmt.Errorf("could not deploy swap router B: %w", err)
	}
	info.RouterB, err = contract.NewSwapRouter(info.RouterBAddr, client.GetEthClient())
	if err != nil {
		return nil, fmt.Errorf("could not create instance of swap router B: %w", err)
	}

	// deploy liquidity provider helper
	info.LiquidityProviderAddr, err = deployContract(contract.V3LiquidityProviderMetaData, false, 0, ownerWallet.GetAddress(), rootAddr, info.Weth9Addr)
	if err != nil {
		return nil, fmt.Errorf("could not deploy v3 liquidity provider: %w", err)
	}
	info.LiquidityProvider, err = contract.NewV3LiquidityProvider(info.LiquidityProviderAddr, client.GetEthClient())
	if err != nil {
		return nil, fmt.Errorf("could not create instance of v3 liquidity provider: %w", err)
	}

	// deploy DAI tokens (one per pair)
	for i := uint64(0); i < u.options.DaiPairs; i++ {
		poolInfo := V3PoolDeploymentInfo{}
		poolInfo.DaiAddr, err = deployContract(contract.DaiMetaData, true, uint32(i), deployerWallet.GetChainId())
		if err != nil {
			return nil, fmt.Errorf("could not deploy Dai: %w", err)
		}
		poolInfo.Dai, err = contract.NewDai(poolInfo.DaiAddr, client.GetEthClient())
		if err != nil {
			return nil, fmt.Errorf("could not create instance of Dai: %w", err)
		}
		poolInfo.WethIsToken0 = info.Weth9Addr.Big().Cmp(poolInfo.DaiAddr.Big()) < 0
		info.Pools = append(info.Pools, poolInfo)
	}

	// submit & await deployment batch
	if err := u.sendBatch(deployerWallet, client, deploymentTxs, "deploying contracts v3"); err != nil {
		return nil, err
	}

	// read the tick spacing for the configured fee tier
	callOpts := &bind.CallOpts{Context: u.ctx}
	tickSpacing, err := info.FactoryA.FeeAmountTickSpacing(callOpts, info.Fee)
	if err != nil {
		return nil, fmt.Errorf("could not read tick spacing: %w", err)
	}
	if tickSpacing == nil || tickSpacing.Sign() == 0 {
		return nil, fmt.Errorf("unsupported fee tier %d (no tick spacing)", u.options.FeeTier)
	}
	info.TickSpacing = tickSpacing.Int64()

	// Phase 2: create the per-factory pools that don't exist yet.
	createTxs := []*types.Transaction{}
	for i := range info.Pools {
		dai := info.Pools[i].DaiAddr
		for _, factory := range []*contract.UniswapV3Factory{info.FactoryA, info.FactoryB} {
			poolAddr, err := factory.GetPool(callOpts, dai, info.Weth9Addr, info.Fee)
			if err != nil {
				return nil, fmt.Errorf("could not check pool existence: %w", err)
			}
			if poolAddr != (common.Address{}) {
				continue
			}
			factory := factory
			tx, err := deployerWallet.BuildBoundTxWithEstimate(u.ctx, client, u.walletPool.GetTxPool(), &txbuilder.TxMetadata{
				GasFeeCap: uint256.MustFromBig(feeCap),
				GasTipCap: uint256.MustFromBig(tipCap),
				Value:     uint256.NewInt(0),
			}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
				return factory.CreatePool(transactOpts, dai, info.Weth9Addr, info.Fee)
			})
			if err != nil {
				return nil, fmt.Errorf("could not create pool: %w", err)
			}
			createTxs = append(createTxs, tx)
		}
	}
	if err := u.sendBatch(deployerWallet, client, createTxs, "creating pools"); err != nil {
		return nil, err
	}

	// resolve pool addresses and bind instances
	for i := range info.Pools {
		dai := info.Pools[i].DaiAddr

		poolAAddr, err := info.FactoryA.GetPool(callOpts, dai, info.Weth9Addr, info.Fee)
		if err != nil {
			return nil, fmt.Errorf("could not read pool A address: %w", err)
		}
		poolBAddr, err := info.FactoryB.GetPool(callOpts, dai, info.Weth9Addr, info.Fee)
		if err != nil {
			return nil, fmt.Errorf("could not read pool B address: %w", err)
		}
		if poolAAddr == (common.Address{}) || poolBAddr == (common.Address{}) {
			return nil, fmt.Errorf("pool for dai %v was not created", dai.Hex())
		}

		info.Pools[i].PoolAAddr = poolAAddr
		info.Pools[i].PoolA, err = contract.NewUniswapV3Pool(poolAAddr, client.GetEthClient())
		if err != nil {
			return nil, fmt.Errorf("could not create instance of uniswap v3 pool A: %w", err)
		}
		info.Pools[i].PoolBAddr = poolBAddr
		info.Pools[i].PoolB, err = contract.NewUniswapV3Pool(poolBAddr, client.GetEthClient())
		if err != nil {
			return nil, fmt.Errorf("could not create instance of uniswap v3 pool B: %w", err)
		}
	}

	// Phase 3: initialize pools and grant the liquidity provider mint rights.
	setupTxs := []*types.Transaction{}
	for i := range info.Pools {
		poolInfo := info.Pools[i]
		sqrtPriceX96 := u.v3SqrtPriceX96(poolInfo.WethIsToken0)

		for _, pool := range []*contract.UniswapV3Pool{poolInfo.PoolA, poolInfo.PoolB} {
			slot0, err := pool.Slot0(callOpts)
			if err != nil {
				return nil, fmt.Errorf("could not read pool slot0: %w", err)
			}
			if slot0.SqrtPriceX96.Sign() != 0 {
				continue
			}
			pool := pool
			tx, err := ownerWallet.BuildBoundTxWithEstimate(u.ctx, client, u.walletPool.GetTxPool(), &txbuilder.TxMetadata{
				GasFeeCap: uint256.MustFromBig(feeCap),
				GasTipCap: uint256.MustFromBig(tipCap),
				Value:     uint256.NewInt(0),
			}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
				return pool.Initialize(transactOpts, sqrtPriceX96)
			})
			if err != nil {
				return nil, fmt.Errorf("could not initialize pool: %w", err)
			}
			setupTxs = append(setupTxs, tx)
		}
	}
	if err := u.sendBatch(ownerWallet, client, setupTxs, "initializing pools"); err != nil {
		return nil, err
	}

	// Phase 4: seed full-range liquidity into every pool.
	if err := u.provideV3Liquidity(info, client, feeCap, tipCap); err != nil {
		return nil, err
	}

	return info, nil
}

// InitializeContractsV3 binds the deployed v3 contract instances to the static
// call client and stores the deployment for the swap phase.
func (u *Uniswap) InitializeContractsV3(info *V3DeploymentInfo) error {
	client := u.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0),
		spamoor.WithoutBuilder(), // avoid using builders for eth_calls
	)
	if client == nil {
		return scenario.ErrNoClients
	}

	u.logger.Infof("Using client for static calls: %s", client.GetName())

	weth, err := contract.NewWETH9(info.Weth9Addr, client.GetEthClient())
	if err != nil {
		return fmt.Errorf("could not initialize WETH9: %w", err)
	}
	u.Weth = weth

	info.RouterA, err = contract.NewSwapRouter(info.RouterAAddr, client.GetEthClient())
	if err != nil {
		return fmt.Errorf("could not initialize swap router A: %w", err)
	}
	info.RouterB, err = contract.NewSwapRouter(info.RouterBAddr, client.GetEthClient())
	if err != nil {
		return fmt.Errorf("could not initialize swap router B: %w", err)
	}

	u.Tokens = make(map[common.Address]*contract.Dai, len(info.Pools))
	for _, poolInfo := range info.Pools {
		token, err := contract.NewDai(poolInfo.DaiAddr, client.GetEthClient())
		if err != nil {
			return fmt.Errorf("could not initialize token %v: %w", poolInfo.DaiAddr, err)
		}
		u.Tokens[poolInfo.DaiAddr] = token
	}

	u.v3Deployment = info
	return nil
}

// v3SqrtPriceX96 returns the starting price for a pool, derived from the desired
// DAI/WETH reserve ratio (the same ratio the v2 path uses for liquidity depth).
func (u *Uniswap) v3SqrtPriceX96(wethIsToken0 bool) *big.Int {
	wethReserve := u.options.EthLiquidityPerPair.ToBig()
	daiReserve := new(big.Int).Mul(wethReserve, new(big.Int).SetUint64(u.options.DaiLiquidityFactor))

	if wethIsToken0 {
		// token0 = WETH, token1 = DAI -> price = DAI/WETH
		return encodeSqrtRatioX96(daiReserve, wethReserve)
	}
	// token0 = DAI, token1 = WETH -> price = WETH/DAI
	return encodeSqrtRatioX96(wethReserve, daiReserve)
}

// provideV3Liquidity seeds a full-range position into every pool from the root
// wallet, forwarding ETH for the WETH side while DAI is minted on demand.
func (u *Uniswap) provideV3Liquidity(info *V3DeploymentInfo, client *spamoor.Client, feeCap, tipCap *big.Int) error {
	tickLower, tickUpper := fullRangeTicks(info.TickSpacing)

	// each DAI has a pool on both factories -> two liquidity txs per DAI.
	poolCount := len(info.Pools) * 2

	pairFundingAmount := uint256.NewInt(0)
	for i := 0; i < poolCount; i++ {
		pairFundingAmount = pairFundingAmount.Add(pairFundingAmount, u.options.EthLiquidityPerPair)
		fundingFees := uint256.NewInt(6000000)
		fundingFees = fundingFees.Mul(fundingFees, uint256.MustFromBig(feeCap))
		pairFundingAmount = pairFundingAmount.Add(pairFundingAmount, fundingFees)
	}

	rootWallet := u.walletPool.GetRootWallet()
	return rootWallet.WithWalletLock(u.ctx, poolCount, pairFundingAmount, u.walletPool.GetClientPool(), func(reason string) {
		u.logger.Infof("root wallet is locked, %s", reason)
	}, func() error {
		liquidityTxs := []*types.Transaction{}

		// WETH budget bounds the seeded liquidity; DAI is minted on demand.
		wethBudget := new(big.Int).Div(
			new(big.Int).Mul(u.options.EthLiquidityPerPair.ToBig(), big.NewInt(liquidityBudgetBps)),
			big.NewInt(10000),
		)

		for i := range info.Pools {
			poolInfo := info.Pools[i]
			sqrtPriceX96 := u.v3SqrtPriceX96(poolInfo.WethIsToken0)
			liquidity := fullRangeLiquidityForWeth(sqrtPriceX96, poolInfo.WethIsToken0, wethBudget)

			for _, poolAddr := range []common.Address{poolInfo.PoolAAddr, poolInfo.PoolBAddr} {
				poolAddr := poolAddr
				tx, err := rootWallet.GetWallet().BuildBoundTxWithEstimate(u.ctx, client, u.walletPool.GetTxPool(), &txbuilder.TxMetadata{
					GasFeeCap: uint256.MustFromBig(feeCap),
					GasTipCap: uint256.MustFromBig(tipCap),
					Value:     u.options.EthLiquidityPerPair,
				}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
					return info.LiquidityProvider.ProvideLiquidity(transactOpts, poolAddr, tickLower, tickUpper, liquidity)
				})
				if err != nil {
					return fmt.Errorf("could not provide liquidity for pool %v: %w", poolAddr.Hex(), err)
				}
				liquidityTxs = append(liquidityTxs, tx)
			}
		}

		return u.sendBatch(rootWallet.GetWallet(), client, liquidityTxs, "providing liquidity")
	})
}

// sendBatch submits a batch of transactions from a single wallet and waits for
// them to confirm, logging progress. It is a no-op for an empty batch.
func (u *Uniswap) sendBatch(wallet *spamoor.Wallet, client *spamoor.Client, txs []*types.Transaction, action string) error {
	if len(txs) == 0 {
		return nil
	}

	_, err := u.walletPool.GetTxPool().SendTransactionBatch(u.ctx, wallet, txs, &spamoor.BatchOptions{
		SendTransactionOptions: spamoor.SendTransactionOptions{
			Client:      client,
			ClientGroup: u.options.ClientGroup,
		},
		MaxRetries:   3,
		PendingLimit: 10,
		LogFn: func(confirmedCount int, totalCount int) {
			u.logger.Infof("%s... (%v/%v)", action, confirmedCount, totalCount)
		},
		LogInterval: 10,
	})
	if err != nil {
		return fmt.Errorf("could not %s: %w", action, err)
	}
	u.logger.Infof("%s complete. (%v/%v)", action, len(txs), len(txs))
	return nil
}
