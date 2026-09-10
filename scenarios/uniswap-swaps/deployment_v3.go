package uniswapswaps

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/scenarios/uniswap-swaps/contract"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/txtypes"
)

// V3DeploymentInfo holds the deployed Uniswap v3 contract set for the scenario.
// Two factories (each with its own SwapRouter) are deployed so that every DAI
// token gets a separate pool with the shared quote token per factory at the
// same fee tier, mirroring the two-factory layout of the v2 path.
type V3DeploymentInfo struct {
	// Weth9Addr is only needed because the canonical SwapRouter requires a WETH
	// address at construction; the scenario never trades ETH/WETH.
	Weth9Addr             common.Address
	QuoteAddr             common.Address
	Quote                 *contract.Dai
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
	DaiAddr       common.Address
	Dai           *contract.Dai
	QuoteIsToken0 bool
	PoolAAddr     common.Address
	PoolA         *contract.UniswapV3Pool
	PoolBAddr     common.Address
	PoolB         *contract.UniswapV3Pool
}

// DeployUniswapV3 deploys two canonical Uniswap v3 factories + SwapRouters, the
// custom liquidity provider, the shared quote token and one DAI token per
// configured pair. Each DAI gets a pool on both factories, which are then
// initialized and seeded with a full-range position.
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

	deployerSeed := [32]byte{}
	copy(deployerSeed[:], deployerWallet.GetAddress().Bytes())

	baseFeeWei, tipFeeWei := spamoor.ResolveFees(u.options.BaseFee, u.options.TipFee, u.options.BaseFeeWei, u.options.TipFeeWei)
	feeCap, tipCap, err := u.walletPool.GetSuggestedFees(client, baseFeeWei, tipFeeWei)
	if err != nil {
		return nil, fmt.Errorf("could not get tx fee: %w", err)
	}

	deploymentTxs := []*txtypes.Transaction{}
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

	// deploy WETH9 (SwapRouter constructor dependency only)
	info.Weth9Addr, err = deployContract(contract.WETH9MetaData, true, 0)
	if err != nil {
		return nil, fmt.Errorf("could not deploy WETH9: %w", err)
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

	// deploy liquidity provider helper (owner-only, mints both tokens on demand)
	info.LiquidityProviderAddr, err = deployContract(contract.V3LiquidityProviderMetaData, false, 0, ownerWallet.GetAddress())
	if err != nil {
		return nil, fmt.Errorf("could not deploy v3 liquidity provider: %w", err)
	}
	info.LiquidityProvider, err = contract.NewV3LiquidityProvider(info.LiquidityProviderAddr, client.GetEthClient())
	if err != nil {
		return nil, fmt.Errorf("could not create instance of v3 liquidity provider: %w", err)
	}

	// deploy the shared quote token
	info.QuoteAddr, err = deployContract(contract.DaiMetaData, true, quoteTokenSalt, deployerWallet.GetChainId())
	if err != nil {
		return nil, fmt.Errorf("could not deploy quote token: %w", err)
	}
	info.Quote, err = contract.NewDai(info.QuoteAddr, client.GetEthClient())
	if err != nil {
		return nil, fmt.Errorf("could not create instance of quote token: %w", err)
	}

	// deploy DAI tokens (one per pair)
	for i := uint64(0); i < u.options.PairCount; i++ {
		poolInfo := V3PoolDeploymentInfo{}
		poolInfo.DaiAddr, err = deployContract(contract.DaiMetaData, true, daiTokenSalt+uint32(i), deployerWallet.GetChainId())
		if err != nil {
			return nil, fmt.Errorf("could not deploy Dai: %w", err)
		}
		poolInfo.Dai, err = contract.NewDai(poolInfo.DaiAddr, client.GetEthClient())
		if err != nil {
			return nil, fmt.Errorf("could not create instance of Dai: %w", err)
		}
		poolInfo.QuoteIsToken0 = info.QuoteAddr.Big().Cmp(poolInfo.DaiAddr.Big()) < 0
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
	createTxs := []*txtypes.Transaction{}
	for i := range info.Pools {
		dai := info.Pools[i].DaiAddr
		for _, factory := range []*contract.UniswapV3Factory{info.FactoryA, info.FactoryB} {
			poolAddr, err := factory.GetPool(callOpts, dai, info.QuoteAddr, info.Fee)
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
				return factory.CreatePool(transactOpts, dai, info.QuoteAddr, info.Fee)
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

		poolAAddr, err := info.FactoryA.GetPool(callOpts, dai, info.QuoteAddr, info.Fee)
		if err != nil {
			return nil, fmt.Errorf("could not read pool A address: %w", err)
		}
		poolBAddr, err := info.FactoryB.GetPool(callOpts, dai, info.QuoteAddr, info.Fee)
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

	// Phase 3: initialize pools that have no price yet.
	setupTxs := []*txtypes.Transaction{}
	for i := range info.Pools {
		poolInfo := info.Pools[i]
		sqrtPriceX96 := u.v3SqrtPriceX96(poolInfo.QuoteIsToken0)

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
	if err := u.provideV3Liquidity(info, client, ownerWallet, feeCap, tipCap); err != nil {
		return nil, err
	}

	return info, nil
}

// InitializeContractsV3 binds the deployed v3 contract instances to the static
// call client and stores the deployment for the swap phase.
func (u *Uniswap) InitializeContractsV3(info *V3DeploymentInfo) error {
	client, err := u.staticCallClient()
	if err != nil {
		return err
	}

	info.RouterA, err = contract.NewSwapRouter(info.RouterAAddr, client.GetEthClient())
	if err != nil {
		return fmt.Errorf("could not initialize swap router A: %w", err)
	}
	info.RouterB, err = contract.NewSwapRouter(info.RouterBAddr, client.GetEthClient())
	if err != nil {
		return fmt.Errorf("could not initialize swap router B: %w", err)
	}
	for i := range info.Pools {
		info.Pools[i].PoolA, err = contract.NewUniswapV3Pool(info.Pools[i].PoolAAddr, client.GetEthClient())
		if err != nil {
			return fmt.Errorf("could not initialize pool A: %w", err)
		}
		info.Pools[i].PoolB, err = contract.NewUniswapV3Pool(info.Pools[i].PoolBAddr, client.GetEthClient())
		if err != nil {
			return fmt.Errorf("could not initialize pool B: %w", err)
		}
	}
	u.Quote, err = contract.NewDai(info.QuoteAddr, client.GetEthClient())
	if err != nil {
		return fmt.Errorf("could not initialize quote token: %w", err)
	}

	u.v3Deployment = info
	return nil
}

// v3SqrtPriceX96 returns the starting price for a pool, derived from the
// seeded DAI/quote reserve ratio (TokensPerQuote DAI per quote token).
func (u *Uniswap) v3SqrtPriceX96(quoteIsToken0 bool) *big.Int {
	quoteReserve := u.options.QuoteLiquidityPerPool
	daiReserve := new(big.Int).Mul(quoteReserve, new(big.Int).SetUint64(u.options.TokensPerQuote))

	if quoteIsToken0 {
		// token0 = quote, token1 = DAI -> price = DAI/quote
		return encodeSqrtRatioX96(daiReserve, quoteReserve)
	}
	// token0 = DAI, token1 = quote -> price = quote/DAI
	return encodeSqrtRatioX96(quoteReserve, daiReserve)
}

// provideV3Liquidity seeds a full-range position into every pool that has no
// liquidity yet. Both tokens are minted on demand by the liquidity provider, so
// the owner wallet only pays gas.
func (u *Uniswap) provideV3Liquidity(info *V3DeploymentInfo, client *spamoor.Client, ownerWallet *spamoor.Wallet, feeCap, tipCap *big.Int) error {
	tickLower, tickUpper := fullRangeTicks(info.TickSpacing)
	callOpts := &bind.CallOpts{Context: u.ctx}
	liquidityTxs := []*txtypes.Transaction{}

	for i := range info.Pools {
		poolInfo := info.Pools[i]
		sqrtPriceX96 := u.v3SqrtPriceX96(poolInfo.QuoteIsToken0)
		liquidity := fullRangeLiquidityForToken(sqrtPriceX96, poolInfo.QuoteIsToken0, u.options.QuoteLiquidityPerPool)

		pools := []struct {
			addr common.Address
			pool *contract.UniswapV3Pool
		}{
			{poolInfo.PoolAAddr, poolInfo.PoolA},
			{poolInfo.PoolBAddr, poolInfo.PoolB},
		}
		for _, p := range pools {
			existing, err := p.pool.Liquidity(callOpts)
			if err != nil {
				return fmt.Errorf("could not read pool liquidity: %w", err)
			}
			if existing.Sign() > 0 {
				continue
			}

			poolAddr := p.addr
			tx, err := ownerWallet.BuildBoundTxWithEstimate(u.ctx, client, u.walletPool.GetTxPool(), &txbuilder.TxMetadata{
				GasFeeCap: uint256.MustFromBig(feeCap),
				GasTipCap: uint256.MustFromBig(tipCap),
				Value:     uint256.NewInt(0),
			}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
				return info.LiquidityProvider.ProvideLiquidity(transactOpts, poolAddr, tickLower, tickUpper, liquidity)
			})
			if err != nil {
				return fmt.Errorf("could not provide liquidity for pool %v: %w", poolAddr.Hex(), err)
			}
			liquidityTxs = append(liquidityTxs, tx)
		}
	}

	return u.sendBatch(ownerWallet, client, liquidityTxs, "providing liquidity")
}

// sendBatch submits a batch of transactions from a single wallet and waits for
// them to confirm, logging progress. It is a no-op for an empty batch.
func (u *Uniswap) sendBatch(wallet *spamoor.Wallet, client *spamoor.Client, txs []*txtypes.Transaction, action string) error {
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
