package mevswaps

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethpandaops/spamoor/scenarios/mev-swaps/contract"
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

func (s *Scenario) deployUniswapPairs(ctx context.Context, redeploy bool, ethLiquidityPerPair *uint256.Int, daiLiquidityFactor uint64) (*DeploymentInfo, *txbuilder.Client, error) {
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0)
	deployerWallet := s.walletPool.GetWellKnownWallet("deployer")
	ownerWallet := s.walletPool.GetWellKnownWallet("owner")

	var feeCap *big.Int
	var tipCap *big.Int

	if s.options.BaseFee > 0 {
		feeCap = new(big.Int).Mul(big.NewInt(int64(s.options.BaseFee)), big.NewInt(1000000000))
	}
	if s.options.TipFee > 0 {
		tipCap = new(big.Int).Mul(big.NewInt(int64(s.options.TipFee)), big.NewInt(1000000000))
	}

	if feeCap == nil || tipCap == nil {
		var err error
		feeCap, tipCap, err = client.GetSuggestedFee(s.walletPool.GetContext())
		if err != nil {
			return nil, client, err
		}
	}

	if feeCap.Cmp(big.NewInt(1000000000)) < 0 {
		feeCap = big.NewInt(1000000000)
	}
	if tipCap.Cmp(big.NewInt(1000000000)) < 0 {
		tipCap = big.NewInt(1000000000)
	}

	deploymentTxs := []*types.Transaction{}
	deploymentInfo := &DeploymentInfo{}
	deployerNonce := deployerWallet.GetNonce()
	contractNonce := uint64(0)
	usedNonce := uint64(0)
	var err error

	// deploy WETH9
	if redeploy || deployerNonce <= contractNonce {
		tx, err := deployerWallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Gas:       600000,
			Value:     uint256.NewInt(0),
		}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
			_, deployTx, _, err := contract.DeployWETH9(transactOpts, client.GetEthClient())
			return deployTx, err
		})
		if err != nil {
			return nil, client, fmt.Errorf("could not deploy WETH9: %w", err)
		}
		deploymentTxs = append(deploymentTxs, tx)
		usedNonce = tx.Nonce()
	} else {
		usedNonce = contractNonce
	}
	contractNonce++

	deploymentInfo.Weth9Addr = crypto.CreateAddress(deployerWallet.GetAddress(), usedNonce)
	deploymentInfo.Weth9, err = contract.NewWETH9(deploymentInfo.Weth9Addr, client.GetEthClient())
	if err != nil {
		return nil, client, fmt.Errorf("could not create instance of WETH9: %w", err)
	}

	// deploy two uniswap factories
	if redeploy || deployerNonce <= contractNonce {
		tx, err := deployerWallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Gas:       3100000,
			Value:     uint256.NewInt(0),
		}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
			_, deployTx, _, err := contract.DeployUniswapV2Factory(transactOpts, client.GetEthClient(), ownerWallet.GetAddress())
			return deployTx, err
		})
		if err != nil {
			return nil, client, fmt.Errorf("could not deploy uniswap v2 factory A: %w", err)
		}
		deploymentTxs = append(deploymentTxs, tx)
		usedNonce = tx.Nonce()
	} else {
		usedNonce = contractNonce
	}
	contractNonce++

	deploymentInfo.UniswapFactoryAAddr = crypto.CreateAddress(deployerWallet.GetAddress(), usedNonce)
	deploymentInfo.UniswapFactoryA, err = contract.NewUniswapV2Factory(deploymentInfo.UniswapFactoryAAddr, client.GetEthClient())
	if err != nil {
		return nil, client, fmt.Errorf("could not create instance of uniswap v2 factory A: %w", err)
	}

	if redeploy || deployerNonce <= contractNonce {
		tx, err := deployerWallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Gas:       3100000,
			Value:     uint256.NewInt(0),
		}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
			_, deployTx, _, err := contract.DeployUniswapV2Factory(transactOpts, client.GetEthClient(), ownerWallet.GetAddress())
			return deployTx, err
		})
		if err != nil {
			return nil, client, fmt.Errorf("could not deploy uniswap v2 factory B: %w", err)
		}
		deploymentTxs = append(deploymentTxs, tx)
		usedNonce = tx.Nonce()
	} else {
		usedNonce = contractNonce
	}
	contractNonce++

	deploymentInfo.UniswapFactoryBAddr = crypto.CreateAddress(deployerWallet.GetAddress(), usedNonce)
	deploymentInfo.UniswapFactoryB, err = contract.NewUniswapV2Factory(deploymentInfo.UniswapFactoryBAddr, client.GetEthClient())
	if err != nil {
		return nil, client, fmt.Errorf("could not create instance of uniswap v2 factory B: %w", err)
	}

	// deploy two uniswap routers
	if redeploy || deployerNonce <= contractNonce {
		tx, err := deployerWallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Gas:       5000000,
			Value:     uint256.NewInt(0),
		}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
			_, deployTx, _, err := contract.DeployUniswapV2Router02(transactOpts, client.GetEthClient(), deploymentInfo.UniswapFactoryAAddr, deploymentInfo.Weth9Addr)
			return deployTx, err
		})
		if err != nil {
			return nil, client, fmt.Errorf("could not deploy uniswap v2 router A: %w", err)
		}
		deploymentTxs = append(deploymentTxs, tx)
		usedNonce = tx.Nonce()
	} else {
		usedNonce = contractNonce
	}
	contractNonce++

	deploymentInfo.UniswapRouterAAddr = crypto.CreateAddress(deployerWallet.GetAddress(), usedNonce)
	deploymentInfo.UniswapRouterA, err = contract.NewUniswapV2Router02(deploymentInfo.UniswapRouterAAddr, client.GetEthClient())
	if err != nil {
		return nil, client, fmt.Errorf("could not create instance of uniswap v2 router A: %w", err)
	}

	if redeploy || deployerNonce <= contractNonce {
		tx, err := deployerWallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Gas:       5000000,
			Value:     uint256.NewInt(0),
		}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
			_, deployTx, _, err := contract.DeployUniswapV2Router02(transactOpts, client.GetEthClient(), deploymentInfo.UniswapFactoryBAddr, deploymentInfo.Weth9Addr)
			return deployTx, err
		})
		if err != nil {
			return nil, client, fmt.Errorf("could not deploy uniswap v2 router B: %w", err)
		}
		deploymentTxs = append(deploymentTxs, tx)
		usedNonce = tx.Nonce()
	} else {
		usedNonce = contractNonce
	}
	contractNonce++

	deploymentInfo.UniswapRouterBAddr = crypto.CreateAddress(deployerWallet.GetAddress(), usedNonce)
	deploymentInfo.UniswapRouterB, err = contract.NewUniswapV2Router02(deploymentInfo.UniswapRouterBAddr, client.GetEthClient())
	if err != nil {
		return nil, client, fmt.Errorf("could not create instance of uniswap v2 router B: %w", err)
	}

	// deploy pair liquidity provider
	if redeploy || deployerNonce <= contractNonce {
		tx, err := deployerWallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Gas:       800000,
			Value:     uint256.NewInt(0),
		}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
			_, deployTx, _, err := contract.DeployPairLiquidityProvider(transactOpts, client.GetEthClient(), ownerWallet.GetAddress(), s.walletPool.GetRootWallet().GetAddress(), deploymentInfo.UniswapRouterAAddr, deploymentInfo.UniswapRouterBAddr, deploymentInfo.Weth9Addr)
			return deployTx, err
		})
		if err != nil {
			return nil, client, fmt.Errorf("could not deploy pair liquidity provider: %w", err)
		}
		deploymentTxs = append(deploymentTxs, tx)
		usedNonce = tx.Nonce()
	} else {
		usedNonce = contractNonce
	}
	contractNonce++

	deploymentInfo.LiquidityProviderAddr = crypto.CreateAddress(deployerWallet.GetAddress(), usedNonce)
	deploymentInfo.LiquidityProvider, err = contract.NewPairLiquidityProvider(deploymentInfo.LiquidityProviderAddr, client.GetEthClient())
	if err != nil {
		return nil, client, fmt.Errorf("could not create instance of pair liquidity provider: %w", err)
	}

	// deploy tokens and uniswap pairs
	pairInitCode := common.FromHex(contract.UniswapV2PairBin)
	pairInitHash := crypto.Keccak256(pairInitCode)
	var pairSalt [32]byte

	for i := uint64(0); i < s.options.PairCount; i++ {
		pairInfo := &PairDeploymentInfo{}

		// deploy Dai
		if redeploy || deployerNonce <= contractNonce {
			tx, err := deployerWallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
				GasFeeCap: uint256.MustFromBig(feeCap),
				GasTipCap: uint256.MustFromBig(tipCap),
				Gas:       2000000,
				Value:     uint256.NewInt(0),
			}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
				_, deployTx, _, err := contract.DeployDai(transactOpts, client.GetEthClient(), deployerWallet.GetChainId())
				return deployTx, err
			})
			if err != nil {
				return nil, client, fmt.Errorf("could not deploy Dai: %w", err)
			}
			deploymentTxs = append(deploymentTxs, tx)
			usedNonce = tx.Nonce()
		} else {
			usedNonce = contractNonce
		}
		contractNonce++

		pairInfo.DaiAddr = crypto.CreateAddress(deployerWallet.GetAddress(), usedNonce)
		pairInfo.Dai, err = contract.NewDai(pairInfo.DaiAddr, client.GetEthClient())
		if err != nil {
			return nil, client, fmt.Errorf("could not create instance of Dai: %w", err)
		}

		// make owner wallet a minter for the Dai
		if redeploy || deployerNonce <= contractNonce {
			tx, err := deployerWallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
				GasFeeCap: uint256.MustFromBig(feeCap),
				GasTipCap: uint256.MustFromBig(tipCap),
				Gas:       200000,
				Value:     uint256.NewInt(0),
			}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
				return pairInfo.Dai.Rely(transactOpts, s.walletPool.GetRootWallet().GetAddress())
			})
			if err != nil {
				return nil, client, fmt.Errorf("could not make owner wallet a minter for the Dai: %w", err)
			}
			deploymentTxs = append(deploymentTxs, tx)
		}
		contractNonce++

		// make liquidity provider a minter for the Dai
		if redeploy || deployerNonce <= contractNonce {
			tx, err := deployerWallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
				GasFeeCap: uint256.MustFromBig(feeCap),
				GasTipCap: uint256.MustFromBig(tipCap),
				Gas:       200000,
				Value:     uint256.NewInt(0),
			}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
				return pairInfo.Dai.Rely(transactOpts, deploymentInfo.LiquidityProviderAddr)
			})
			if err != nil {
				return nil, client, fmt.Errorf("could not make liquidity provider a minter for the Dai: %w", err)
			}
			deploymentTxs = append(deploymentTxs, tx)
		}
		contractNonce++

		// get pair on factory A
		if pairInfo.DaiAddr.Big().Cmp(deploymentInfo.Weth9Addr.Big()) < 0 {
			copy(pairSalt[:], crypto.Keccak256(pairInfo.DaiAddr.Bytes(), deploymentInfo.Weth9Addr.Bytes()))
		} else {
			copy(pairSalt[:], crypto.Keccak256(deploymentInfo.Weth9Addr.Bytes(), pairInfo.DaiAddr.Bytes()))
		}
		pairInfo.PairAddrA = crypto.CreateAddress2(deploymentInfo.UniswapFactoryAAddr, pairSalt, pairInitHash)
		pairInfo.PairA, err = contract.NewUniswapV2Pair(pairInfo.PairAddrA, client.GetEthClient())
		if err != nil {
			return nil, client, fmt.Errorf("could not create instance of uniswap v2 pair A: %w", err)
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
			return nil, client, fmt.Errorf("could not create instance of uniswap v2 pair B: %w", err)
		}

		deploymentInfo.Pairs = append(deploymentInfo.Pairs, *pairInfo)
	}

	// submit & await all deployment transactions
	if len(deploymentTxs) > 0 {
		s.logger.Infof("deploying contracts... (0/%v)", len(deploymentTxs))
		for txIdx := 0; txIdx < len(deploymentTxs); txIdx += 10 {
			endIdx := txIdx + 10
			if txIdx > 0 {
				s.logger.Infof("deploying contracts... (%v/%v)", txIdx, len(deploymentTxs))
			}
			if endIdx > len(deploymentTxs) {
				endIdx = len(deploymentTxs)
			}
			err := s.walletPool.SendTxRange(deploymentTxs[txIdx:endIdx], client, deployerWallet, func(tx *types.Transaction, receipt *types.Receipt, err error) {
				if err != nil {
					s.logger.Warnf("could not send deployment tx %v: %v", tx.Hash().String(), err)
				}
			})
			if err != nil {
				return nil, client, fmt.Errorf("could not send deployment txs: %w", err)
			}
		}
		s.logger.Infof("contract deployment complete. (%v/%v)", len(deploymentTxs), len(deploymentTxs))
	}

	// provide liquidity to the pairs
	liquidityTxs := []*types.Transaction{}
	rootWallet := s.walletPool.GetRootWallet()

	daiLiquidity := new(big.Int).Mul(ethLiquidityPerPair.ToBig(), big.NewInt(int64(daiLiquidityFactor)))

	for _, pairInfo := range deploymentInfo.Pairs {
		tx, err := rootWallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Gas:       6000000,
			Value:     ethLiquidityPerPair,
		}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
			return deploymentInfo.LiquidityProvider.ProvidePairLiquidity(transactOpts, pairInfo.DaiAddr, daiLiquidity)
		})
		if err != nil {
			return nil, client, fmt.Errorf("could not provide liquidity for dai %v: %w", pairInfo.DaiAddr.String(), err)
		}
		liquidityTxs = append(liquidityTxs, tx)
	}

	// submit & await all liquidity txs
	if len(liquidityTxs) > 0 {
		s.logger.Infof("providing liquidity... (0/%v)", len(liquidityTxs))
		for txIdx := 0; txIdx < len(liquidityTxs); txIdx += 50 {
			endIdx := txIdx + 50
			if txIdx > 0 {
				s.logger.Infof("providing liquidity... (%v/%v)", txIdx, len(liquidityTxs))
			}
			if endIdx > len(liquidityTxs) {
				endIdx = len(liquidityTxs)
			}
			err := s.walletPool.SendTxRange(liquidityTxs[txIdx:endIdx], client, rootWallet, func(tx *types.Transaction, receipt *types.Receipt, err error) {
				if err != nil {
					s.logger.Warnf("could not send liquidity tx %v: %v", tx.Hash().String(), err)
				}
			})
			if err != nil {
				return nil, client, fmt.Errorf("could not send liquidity txs: %w", err)
			}
		}
		s.logger.Infof("liquidity provision complete. (%v/%v)", len(liquidityTxs), len(liquidityTxs))
	}

	return deploymentInfo, client, nil
}
