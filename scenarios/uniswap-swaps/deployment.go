package uniswapswaps

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/scenarios/uniswap-swaps/contract"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/holiman/uint256"
)

type DeploymentInfo struct {
	Weth9Addr             common.Address
	Weth9                 *contract.WETH9
	UniswapFactoryAAddr   common.Address
	UniswapFactoryA       *contract.UniswapV2Factory
	UniswapRouterAAddr    common.Address
	UniswapRouterA        *contract.UniswapV2Router02
	UniswapFactoryBAddr   common.Address
	UniswapFactoryB       *contract.UniswapV2Factory
	UniswapRouterBAddr    common.Address
	UniswapRouterB        *contract.UniswapV2Router02
	LiquidityProviderAddr common.Address
	LiquidityProvider     *contract.PairLiquidityProvider
	Pairs                 []PairDeploymentInfo
}

type PairDeploymentInfo struct {
	DaiAddr   common.Address
	Dai       *contract.Dai
	PairAddrA common.Address
	PairA     *contract.UniswapV2Pair
	PairAddrB common.Address
	PairB     *contract.UniswapV2Pair
}

func (u *Uniswap) DeployUniswapPairs(redeploy bool) (*DeploymentInfo, error) {
	client := u.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0),
		spamoor.WithClientGroup(u.options.ClientGroup),
	)
	if client == nil {
		return nil, scenario.ErrNoClients
	}

	deployerWallet := u.walletPool.GetWellKnownWallet("deployer")
	deployerSeed := [32]byte{}
	copy(deployerSeed[:], deployerWallet.GetAddress().Bytes())

	if redeploy {
		copy(deployerSeed[20:], []byte(fmt.Sprintf("%x", deployerWallet.GetNonce()+1)))
	}

	ownerWallet := u.walletPool.GetWellKnownWallet("owner")
	if deployerWallet == nil {
		return nil, scenario.ErrNoWallet
	}
	if ownerWallet == nil {
		return nil, scenario.ErrNoWallet
	}

	baseFeeWei, tipFeeWei := spamoor.ResolveFees(u.options.BaseFee, u.options.TipFee, u.options.BaseFeeWei, u.options.TipFeeWei)
	feeCap, tipCap, err := u.walletPool.GetSuggestedFees(client, baseFeeWei, tipFeeWei)
	if err != nil {
		return nil, fmt.Errorf("could not get tx fee: %w", err)
	}

	deploymentTxs := []*types.Transaction{}
	deploymentInfo := &DeploymentInfo{}
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
			binary.BigEndian.PutUint32(deployerSeed[28:], salt)
		}
		addr, tx, err := u.walletPool.GetDeploymentFactory().GetContractDeployment(u.ctx, initCodeBytes, deployerSeed, client, deployerWallet, feeCap, tipCap, false)
		if err != nil {
			return common.Address{}, err
		}

		if tx != nil {
			deploymentTxs = append(deploymentTxs, tx)
		}

		return addr, nil
	}

	// deploy WETH9
	deploymentInfo.Weth9Addr, err = deployContract(contract.WETH9MetaData, true, 0)
	if err != nil {
		return nil, fmt.Errorf("could not deploy WETH9: %w", err)
	}
	deploymentInfo.Weth9, err = contract.NewWETH9(deploymentInfo.Weth9Addr, client.GetEthClient())
	if err != nil {
		return nil, fmt.Errorf("could not create instance of WETH9: %w", err)
	}

	// deploy uniswap factory A
	deploymentInfo.UniswapFactoryAAddr, err = deployContract(contract.UniswapV2FactoryMetaData, false, 0, ownerWallet.GetAddress())
	if err != nil {
		return nil, fmt.Errorf("could not deploy uniswap v2 factory A: %w", err)
	}
	deploymentInfo.UniswapFactoryA, err = contract.NewUniswapV2Factory(deploymentInfo.UniswapFactoryAAddr, client.GetEthClient())
	if err != nil {
		return nil, fmt.Errorf("could not create instance of uniswap v2 factory A: %w", err)
	}

	// deploy uniswap factory B
	deploymentInfo.UniswapFactoryBAddr, err = deployContract(contract.UniswapV2FactoryMetaData, false, 1, ownerWallet.GetAddress())
	if err != nil {
		return nil, fmt.Errorf("could not deploy uniswap v2 factory B: %w", err)
	}
	deploymentInfo.UniswapFactoryB, err = contract.NewUniswapV2Factory(deploymentInfo.UniswapFactoryBAddr, client.GetEthClient())
	if err != nil {
		return nil, fmt.Errorf("could not create instance of uniswap v2 factory B: %w", err)
	}

	// deploy uniswap router A
	deploymentInfo.UniswapRouterAAddr, err = deployContract(contract.UniswapV2Router02MetaData, false, 0, deploymentInfo.UniswapFactoryAAddr, deploymentInfo.Weth9Addr)
	if err != nil {
		return nil, fmt.Errorf("could not deploy uniswap v2 router A: %w", err)
	}
	deploymentInfo.UniswapRouterA, err = contract.NewUniswapV2Router02(deploymentInfo.UniswapRouterAAddr, client.GetEthClient())
	if err != nil {
		return nil, fmt.Errorf("could not create instance of uniswap v2 router A: %w", err)
	}

	// deploy uniswap router B
	deploymentInfo.UniswapRouterBAddr, err = deployContract(contract.UniswapV2Router02MetaData, false, 1, deploymentInfo.UniswapFactoryBAddr, deploymentInfo.Weth9Addr)
	if err != nil {
		return nil, fmt.Errorf("could not deploy uniswap v2 router B: %w", err)
	}
	deploymentInfo.UniswapRouterB, err = contract.NewUniswapV2Router02(deploymentInfo.UniswapRouterBAddr, client.GetEthClient())
	if err != nil {
		return nil, fmt.Errorf("could not create instance of uniswap v2 router B: %w", err)
	}

	// deploy pair liquidity provider
	deploymentInfo.LiquidityProviderAddr, err = deployContract(
		contract.PairLiquidityProviderMetaData, false, 0,
		ownerWallet.GetAddress(),
		u.walletPool.GetRootWallet().GetWallet().GetAddress(),
		deploymentInfo.UniswapRouterAAddr,
		deploymentInfo.UniswapRouterBAddr,
		deploymentInfo.Weth9Addr,
	)
	if err != nil {
		return nil, fmt.Errorf("could not deploy pair liquidity provider: %w", err)
	}
	deploymentInfo.LiquidityProvider, err = contract.NewPairLiquidityProvider(deploymentInfo.LiquidityProviderAddr, client.GetEthClient())
	if err != nil {
		return nil, fmt.Errorf("could not create instance of pair liquidity provider: %w", err)
	}

	// deploy tokens and uniswap pairs
	pairInitCode := common.FromHex(contract.UniswapV2PairBin)
	pairInitHash := crypto.Keccak256(pairInitCode)
	pairFundingAmount := uint256.NewInt(0)
	var pairSalt [32]byte

	for i := uint64(0); i < u.options.DaiPairs; i++ {
		pairInfo := &PairDeploymentInfo{}

		// deploy Dai
		pairInfo.DaiAddr, err = deployContract(contract.DaiMetaData, true, uint32(i), deployerWallet.GetChainId(), ownerWallet.GetAddress())
		if err != nil {
			return nil, fmt.Errorf("could not deploy Dai: %w", err)
		}
		pairInfo.Dai, err = contract.NewDai(pairInfo.DaiAddr, client.GetEthClient())
		if err != nil {
			return nil, fmt.Errorf("could not create instance of Dai: %w", err)
		}

		// get pair on factory A
		if pairInfo.DaiAddr.Big().Cmp(deploymentInfo.Weth9Addr.Big()) < 0 {
			copy(pairSalt[:], crypto.Keccak256(pairInfo.DaiAddr.Bytes(), deploymentInfo.Weth9Addr.Bytes()))
		} else {
			copy(pairSalt[:], crypto.Keccak256(deploymentInfo.Weth9Addr.Bytes(), pairInfo.DaiAddr.Bytes()))
		}
		pairInfo.PairAddrA = crypto.CreateAddress2(deploymentInfo.UniswapFactoryAAddr, pairSalt, pairInitHash)
		pairInfo.PairA, err = contract.NewUniswapV2Pair(pairInfo.PairAddrA, client.GetEthClient())
		if err != nil {
			return nil, fmt.Errorf("could not create instance of uniswap v2 pair A: %w", err)
		}

		// get pair on factory B
		if pairInfo.DaiAddr.Big().Cmp(deploymentInfo.Weth9Addr.Big()) < 0 {
			copy(pairSalt[:], crypto.Keccak256(pairInfo.DaiAddr.Bytes(), deploymentInfo.Weth9Addr.Bytes()))
		} else {
			copy(pairSalt[:], crypto.Keccak256(deploymentInfo.Weth9Addr.Bytes(), pairInfo.DaiAddr.Bytes()))
		}
		pairInfo.PairAddrB = crypto.CreateAddress2(deploymentInfo.UniswapFactoryBAddr, pairSalt, pairInitHash)
		pairInfo.PairB, err = contract.NewUniswapV2Pair(pairInfo.PairAddrB, client.GetEthClient())
		if err != nil {
			return nil, fmt.Errorf("could not create instance of uniswap v2 pair B: %w", err)
		}

		deploymentInfo.Pairs = append(deploymentInfo.Pairs, *pairInfo)

		pairFundingAmount = pairFundingAmount.Add(pairFundingAmount, u.options.EthLiquidityPerPair)
		fundingFees := uint256.NewInt(6000000)
		fundingFees = fundingFees.Mul(fundingFees, uint256.MustFromBig(feeCap))
		pairFundingAmount = pairFundingAmount.Add(pairFundingAmount, fundingFees)
	}

	// submit & await all deployment transactions
	if len(deploymentTxs) > 0 {
		_, err := u.walletPool.GetTxPool().SendTransactionBatch(u.ctx, deployerWallet, deploymentTxs, &spamoor.BatchOptions{
			SendTransactionOptions: spamoor.SendTransactionOptions{
				Client:      client,
				ClientGroup: u.options.ClientGroup,
			},
			MaxRetries:   3,
			PendingLimit: 10,
			LogFn: func(confirmedCount int, totalCount int) {
				u.logger.Infof("deploying contracts... (%v/%v)", confirmedCount, totalCount)
			},
			LogInterval: 10,
		})
		if err != nil {
			return nil, fmt.Errorf("could not send deployment txs: %w", err)
		}
		u.logger.Infof("contract deployment complete. (%v/%v)", len(deploymentTxs), len(deploymentTxs))
	}

	// Phase 2: post-deployment setup calls. Built only after the deployment
	// batch has been mined so eth_estimateGas dispatches into the real
	// contract code instead of treating the target as an EOA.
	setupTxs := []*types.Transaction{}
	callOpts := &bind.CallOpts{Context: u.ctx}

	for _, pairInfo := range deploymentInfo.Pairs {
		// make liquidity provider a minter for the Dai
		lpIsWard, err := pairInfo.Dai.Wards(callOpts, deploymentInfo.LiquidityProviderAddr)
		if err != nil {
			return nil, fmt.Errorf("could not check if liquidity provider is a ward for the Dai: %w", err)
		}
		if lpIsWard.Cmp(big.NewInt(0)) == 0 {
			tx, err := ownerWallet.BuildBoundTxWithEstimate(u.ctx, client, u.walletPool.GetTxPool(), &txbuilder.TxMetadata{
				GasFeeCap: uint256.MustFromBig(feeCap),
				GasTipCap: uint256.MustFromBig(tipCap),
				Value:     uint256.NewInt(0),
			}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
				return pairInfo.Dai.Rely(transactOpts, deploymentInfo.LiquidityProviderAddr)
			})
			if err != nil {
				return nil, fmt.Errorf("could not make liquidity provider a minter for the Dai: %w", err)
			}
			setupTxs = append(setupTxs, tx)
		}
	}

	if len(setupTxs) > 0 {
		_, err := u.walletPool.GetTxPool().SendTransactionBatch(u.ctx, ownerWallet, setupTxs, &spamoor.BatchOptions{
			SendTransactionOptions: spamoor.SendTransactionOptions{
				Client:      client,
				ClientGroup: u.options.ClientGroup,
			},
			MaxRetries:   3,
			PendingLimit: 10,
			LogFn: func(confirmedCount int, totalCount int) {
				u.logger.Infof("running post-deployment setup... (%v/%v)", confirmedCount, totalCount)
			},
			LogInterval: 10,
		})
		if err != nil {
			return nil, fmt.Errorf("could not send post-deployment setup txs: %w", err)
		}
		u.logger.Infof("post-deployment setup complete. (%v/%v)", len(setupTxs), len(setupTxs))
	}

	// provide liquidity to the pairs
	rootWallet := u.walletPool.GetRootWallet()
	err = rootWallet.WithWalletLock(u.ctx, len(deploymentInfo.Pairs), pairFundingAmount, u.walletPool.GetClientPool(), func(reason string) {
		u.logger.Infof("root wallet is locked, %s", reason)
	}, func() error {
		liquidityTxs := []*types.Transaction{}
		daiLiquidity := new(big.Int).Mul(u.options.EthLiquidityPerPair.ToBig(), big.NewInt(int64(u.options.DaiLiquidityFactor)))

		for _, pairInfo := range deploymentInfo.Pairs {
			tx, err := rootWallet.GetWallet().BuildBoundTxWithEstimate(u.ctx, client, u.walletPool.GetTxPool(), &txbuilder.TxMetadata{
				GasFeeCap: uint256.MustFromBig(feeCap),
				GasTipCap: uint256.MustFromBig(tipCap),
				Value:     u.options.EthLiquidityPerPair,
			}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
				return deploymentInfo.LiquidityProvider.ProvidePairLiquidity(transactOpts, pairInfo.DaiAddr, daiLiquidity)
			})
			if err != nil {
				return fmt.Errorf("could not provide liquidity for dai %v: %w", pairInfo.DaiAddr.String(), err)
			}
			liquidityTxs = append(liquidityTxs, tx)
		}

		// submit & await all liquidity txs
		if len(liquidityTxs) > 0 {
			_, err := u.walletPool.GetTxPool().SendTransactionBatch(u.ctx, rootWallet.GetWallet(), liquidityTxs, &spamoor.BatchOptions{
				SendTransactionOptions: spamoor.SendTransactionOptions{
					Client:      client,
					ClientGroup: u.options.ClientGroup,
				},
				MaxRetries:   3,
				PendingLimit: 10,
				LogFn: func(confirmedCount int, totalCount int) {
					u.logger.Infof("providing liquidity... (%v/%v)", confirmedCount, totalCount)
				},
				LogInterval: 10,
			})
			if err != nil {
				return fmt.Errorf("could not send liquidity txs: %w", err)
			}

			u.logger.Infof("liquidity provision complete. (%v/%v)", len(liquidityTxs), len(liquidityTxs))
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not provide liquidity: %w", err)
	}

	return deploymentInfo, nil
}
