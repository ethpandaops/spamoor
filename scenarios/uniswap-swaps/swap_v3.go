package uniswapswaps

import (
	"context"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"

	"github.com/ethpandaops/spamoor/scenarios/uniswap-swaps/contract"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
)

// buildV3SwapTx builds a single exact-input swap against a randomly selected
// pool, routed through one of the two SwapRouters (which picks the matching
// factory's pool). Swap sizes are DAI-denominated to match the v2 path; the
// router has no on-chain quote, so amountOutMinimum is left at 0 and balances
// are tracked with a conservative estimate to avoid insufficient-input reverts.
func (s *Scenario) buildV3SwapTx(ctx context.Context, wallet *spamoor.Wallet, feeCap, tipCap *big.Int) (*types.Transaction, error) {
	info := s.uniswap.v3Deployment
	if info == nil || len(info.Pools) == 0 {
		return nil, fmt.Errorf("no v3 pools deployed")
	}

	poolInfo := info.Pools[mathrand.Intn(len(info.Pools))]
	daiAddr := poolInfo.DaiAddr
	wethAddr := info.Weth9Addr

	// alternate between the two routers (factory A vs B pool)
	router := info.RouterA
	if mathrand.Intn(100) < 50 {
		router = info.RouterB
	}

	minAmount, ok := new(big.Int).SetString(s.options.MinSwapAmount, 10)
	if !ok {
		return nil, fmt.Errorf("invalid min swap amount: %s", s.options.MinSwapAmount)
	}
	maxAmount, ok := new(big.Int).SetString(s.options.MaxSwapAmount, 10)
	if !ok {
		return nil, fmt.Errorf("invalid max swap amount: %s", s.options.MaxSwapAmount)
	}
	sellThreshold, ok := new(big.Int).SetString(s.options.SellThreshold, 10)
	if !ok {
		return nil, fmt.Errorf("invalid sell threshold: %s", s.options.SellThreshold)
	}

	diff := new(big.Int).Sub(maxAmount, minAmount)
	randomAmount := new(big.Int).Add(minAmount, new(big.Int).Rand(mathrand.New(mathrand.NewSource(time.Now().UnixNano())), diff))

	priceFactor := new(big.Int).SetUint64(s.uniswap.options.DaiLiquidityFactor)
	if priceFactor.Sign() == 0 {
		priceFactor = big.NewInt(1)
	}

	addr := wallet.GetAddress()
	daiBalance := s.uniswap.GetTokenBalance(addr, daiAddr)
	wethBalance := s.uniswap.GetTokenBalance(addr, wethAddr)
	ethBalance := wallet.GetBalance()

	isBuy := mathrand.Intn(100) < int(s.options.BuyRatio)
	if daiBalance.Cmp(sellThreshold) > 0 {
		// too much DAI accumulated, force a sell to avoid depleting balances
		isBuy = false
	}
	if !isBuy && daiBalance.Cmp(randomAmount) < 0 {
		// not enough DAI to sell, switch to buy
		isBuy = true
	}

	deadline := big.NewInt(time.Now().Unix() + 300)
	mkParams := func(tokenIn, tokenOut common.Address, amountIn *big.Int) contract.ISwapRouterExactInputSingleParams {
		return contract.ISwapRouterExactInputSingleParams{
			TokenIn:           tokenIn,
			TokenOut:          tokenOut,
			Fee:               info.Fee,
			Recipient:         addr,
			Deadline:          deadline,
			AmountIn:          amountIn,
			AmountOutMinimum:  big.NewInt(0),
			SqrtPriceLimitX96: big.NewInt(0),
		}
	}

	if isBuy {
		// WETH input needed to buy ~randomAmount DAI worth of tokens.
		wethIn := new(big.Int).Div(randomAmount, priceFactor)
		if wethIn.Sign() == 0 {
			wethIn = big.NewInt(1)
		}
		// conservative DAI output estimate (minus fee/slippage headroom).
		daiOutEst := new(big.Int).Div(new(big.Int).Mul(randomAmount, big.NewInt(95)), big.NewInt(100))
		params := mkParams(wethAddr, daiAddr, wethIn)

		// prefer spending held WETH when available, otherwise pay with ETH.
		if mathrand.Intn(100) < 60 && wethBalance.Cmp(wethIn) >= 0 {
			tx, err := wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
				GasFeeCap: uint256.MustFromBig(feeCap),
				GasTipCap: uint256.MustFromBig(tipCap),
				Gas:       swapGasLimit,
				Value:     uint256.NewInt(0),
			}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
				return router.ExactInputSingle(transactOpts, params)
			})
			if err != nil {
				return nil, err
			}
			s.uniswap.UpdateTokenBalance(addr, wethAddr, new(big.Int).Sub(wethBalance, wethIn))
			s.uniswap.UpdateTokenBalance(addr, daiAddr, new(big.Int).Add(daiBalance, daiOutEst))
			return tx, nil
		}

		// pay with raw ETH (SwapRouter wraps msg.value when tokenIn is WETH).
		if ethBalance.Cmp(wethIn) < 0 {
			return nil, fmt.Errorf("insufficient ETH balance for swap")
		}
		tx, err := wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Gas:       swapGasLimit,
			Value:     uint256.MustFromBig(wethIn),
		}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
			return router.ExactInputSingle(transactOpts, params)
		})
		if err != nil {
			return nil, err
		}
		wallet.SubBalance(wethIn)
		s.uniswap.UpdateTokenBalance(addr, daiAddr, new(big.Int).Add(daiBalance, daiOutEst))
		return tx, nil
	}

	// sell DAI for WETH
	if daiBalance.Cmp(randomAmount) < 0 {
		return nil, fmt.Errorf("insufficient DAI balance for swap")
	}
	wethOutEst := new(big.Int).Div(new(big.Int).Mul(new(big.Int).Div(randomAmount, priceFactor), big.NewInt(95)), big.NewInt(100))
	params := mkParams(daiAddr, wethAddr, randomAmount)
	tx, err := wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       swapGasLimit,
		Value:     uint256.NewInt(0),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		return router.ExactInputSingle(transactOpts, params)
	})
	if err != nil {
		return nil, err
	}
	s.uniswap.UpdateTokenBalance(addr, daiAddr, new(big.Int).Sub(daiBalance, randomAmount))
	s.uniswap.UpdateTokenBalance(addr, wethAddr, new(big.Int).Add(wethBalance, wethOutEst))
	return tx, nil
}
