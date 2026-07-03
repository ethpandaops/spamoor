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

	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
)

// buildV2SwapTx builds a single uniswap v2 swap against a randomly selected pair
// and router, deciding buy vs sell from the configured ratio and the wallet's
// tracked balances. Amounts are DAI-denominated and priced via the router's
// on-chain getAmountsIn/getAmountsOut helpers.
func (s *Scenario) buildV2SwapTx(ctx context.Context, wallet *spamoor.Wallet, feeCap, tipCap *big.Int) (*types.Transaction, error) {
	// Select random pair
	pairIdx := mathrand.Intn(len(s.deploymentInfo.Pairs))
	pair := s.deploymentInfo.Pairs[pairIdx]

	// Parse min and max swap amounts
	minAmount, ok := new(big.Int).SetString(s.options.MinSwapAmount, 10)
	if !ok {
		return nil, fmt.Errorf("invalid min swap amount: %s", s.options.MinSwapAmount)
	}

	maxAmount, ok := new(big.Int).SetString(s.options.MaxSwapAmount, 10)
	if !ok {
		return nil, fmt.Errorf("invalid max swap amount: %s", s.options.MaxSwapAmount)
	}

	// Calculate random swap amount
	diff := new(big.Int).Sub(maxAmount, minAmount)
	randomAmount := new(big.Int).Add(minAmount, new(big.Int).Rand(mathrand.New(mathrand.NewSource(time.Now().UnixNano())), diff))

	// Per-trade slippage tolerance (fixed --slippage or a random draw from
	// the configured [slippage_min, slippage_max] band).
	slippage := s.perTradeSlippage()

	// Get current token balance from cache
	tokenBalance := s.uniswap.GetTokenBalance(wallet.GetAddress(), pair.DaiAddr)

	// Get current ETH balance from wallet
	ethBalance := wallet.GetBalance()

	// Get current WETH balance
	wethBalance := s.uniswap.GetTokenBalance(wallet.GetAddress(), s.deploymentInfo.Weth9Addr)

	// Decide if we're buying or selling based on buy ratio and balances
	isBuy := mathrand.Intn(100) < int(s.options.BuyRatio)

	// Parse sell threshold
	sellThreshold, ok := new(big.Int).SetString(s.options.SellThreshold, 10)
	if !ok {
		return nil, fmt.Errorf("invalid sell threshold: %s", s.options.SellThreshold)
	}

	// If we have a lot of DAI, force a sell to avoid depleting the pool
	if tokenBalance.Cmp(sellThreshold) > 0 {
		isBuy = false
	}

	// If we don't have enough DAI to sell, switch to buy
	if !isBuy && tokenBalance.Cmp(randomAmount) < 0 {
		isBuy = true
	}

	// Alternate between routers based on transaction index
	router := s.uniswap.RouterA
	if mathrand.Intn(100) < 50 {
		router = s.uniswap.RouterB
	}

	var tx *types.Transaction

	if isBuy {
		// Decide whether to use ETH or WETH for buying
		useWeth := mathrand.Intn(100) < 60 // 60% chance to use WETH if available

		if useWeth {
			// Buying DAI with WETH
			// Calculate how much WETH we need to spend to get the desired amount of DAI
			amounts, err := router.GetAmountsIn(&bind.CallOpts{}, randomAmount, []common.Address{s.deploymentInfo.Weth9Addr, pair.DaiAddr})
			if err != nil {
				return nil, err
			}

			wethAmount := amounts[0]

			// Check if we have enough WETH
			if wethBalance.Cmp(wethAmount) < 0 {
				// Fall back to ETH if not enough WETH
				useWeth = false
			} else {
				// Calculate minimum DAI amount to receive (with slippage)
				minDaiAmount := new(big.Int).Mul(randomAmount, big.NewInt(10000-int64(slippage)))
				minDaiAmount = minDaiAmount.Div(minDaiAmount, big.NewInt(10000))

				// Build buy transaction with WETH
				tx, err = wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
					GasFeeCap: uint256.MustFromBig(feeCap),
					GasTipCap: uint256.MustFromBig(tipCap),
					Gas:       swapGasLimit,
					Value:     uint256.NewInt(0),
				}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
					return router.SwapExactTokensForTokens(transactOpts, wethAmount, minDaiAmount, []common.Address{s.deploymentInfo.Weth9Addr, pair.DaiAddr}, wallet.GetAddress(), big.NewInt(time.Now().Unix()+300))
				})
				if err != nil {
					return nil, err
				}

				// Update balances in local cache
				if tx != nil {
					// Subtract WETH amount
					newWethBalance := new(big.Int).Sub(wethBalance, wethAmount)
					s.uniswap.UpdateTokenBalance(wallet.GetAddress(), s.deploymentInfo.Weth9Addr, newWethBalance)

					// Add DAI amount
					newDaiBalance := new(big.Int).Add(tokenBalance, randomAmount)
					s.uniswap.UpdateTokenBalance(wallet.GetAddress(), pair.DaiAddr, newDaiBalance)
				}
			}
		}

		// If not using WETH or not enough WETH, use ETH
		if !useWeth {
			// Buying DAI with ETH
			// Calculate how much ETH we need to spend to get the desired amount of DAI
			amounts, err := router.GetAmountsIn(&bind.CallOpts{}, randomAmount, []common.Address{s.deploymentInfo.Weth9Addr, pair.DaiAddr})
			if err != nil {
				return nil, err
			}

			ethAmount := amounts[0]

			// Check if we have enough ETH
			if ethBalance.Cmp(ethAmount) < 0 {
				return nil, fmt.Errorf("insufficient ETH balance for swap")
			}

			// Calculate minimum DAI amount to receive (with slippage)
			minDaiAmount := new(big.Int).Mul(randomAmount, big.NewInt(10000-int64(slippage)))
			minDaiAmount = minDaiAmount.Div(minDaiAmount, big.NewInt(10000))

			// Build buy transaction
			tx, err = wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
				GasFeeCap: uint256.MustFromBig(feeCap),
				GasTipCap: uint256.MustFromBig(tipCap),
				Gas:       swapGasLimit,
				Value:     uint256.MustFromBig(ethAmount),
			}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
				return router.SwapExactETHForTokens(transactOpts, minDaiAmount, []common.Address{s.deploymentInfo.Weth9Addr, pair.DaiAddr}, wallet.GetAddress(), big.NewInt(time.Now().Unix()+300))
			})
			if err != nil {
				return nil, err
			}

			// Update balances in local cache
			if tx != nil {
				// Subtract ETH amount
				wallet.SubBalance(ethAmount)

				// Add DAI amount
				newDaiBalance := new(big.Int).Add(tokenBalance, randomAmount)
				s.uniswap.UpdateTokenBalance(wallet.GetAddress(), pair.DaiAddr, newDaiBalance)
			}
		}
	} else {
		// Decide whether to keep WETH or convert to ETH
		keepWeth := mathrand.Intn(100) < 30 // 30% chance to keep WETH

		if keepWeth {
			// Selling DAI for WETH
			if tokenBalance.Cmp(randomAmount) < 0 {
				return nil, fmt.Errorf("insufficient DAI balance for swap")
			}

			// Calculate minimum WETH amount to receive (with slippage)
			amounts, err := router.GetAmountsOut(&bind.CallOpts{}, randomAmount, []common.Address{pair.DaiAddr, s.deploymentInfo.Weth9Addr})
			if err != nil {
				return nil, err
			}
			minWethAmount := new(big.Int).Mul(amounts[1], big.NewInt(10000-int64(slippage)))
			minWethAmount = minWethAmount.Div(minWethAmount, big.NewInt(10000))

			// Build sell transaction for WETH
			tx, err = wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
				GasFeeCap: uint256.MustFromBig(feeCap),
				GasTipCap: uint256.MustFromBig(tipCap),
				Gas:       swapGasLimit,
				Value:     uint256.NewInt(0),
			}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
				return router.SwapExactTokensForTokens(transactOpts, randomAmount, minWethAmount, []common.Address{pair.DaiAddr, s.deploymentInfo.Weth9Addr}, wallet.GetAddress(), big.NewInt(time.Now().Unix()+300))
			})
			if err != nil {
				return nil, err
			}

			// Update balances in local cache
			if tx != nil {
				// Subtract DAI amount
				newDaiBalance := new(big.Int).Sub(tokenBalance, randomAmount)
				s.uniswap.UpdateTokenBalance(wallet.GetAddress(), pair.DaiAddr, newDaiBalance)

				// Add WETH amount
				newWethBalance := new(big.Int).Add(wethBalance, amounts[1])
				s.uniswap.UpdateTokenBalance(wallet.GetAddress(), s.deploymentInfo.Weth9Addr, newWethBalance)
			}
		} else {
			// Selling DAI for ETH
			if tokenBalance.Cmp(randomAmount) < 0 {
				return nil, fmt.Errorf("insufficient DAI balance for swap")
			}

			// Calculate minimum ETH amount to receive (with slippage)
			amounts, err := router.GetAmountsOut(&bind.CallOpts{}, randomAmount, []common.Address{pair.DaiAddr, s.deploymentInfo.Weth9Addr})
			if err != nil {
				return nil, err
			}
			minEthAmount := new(big.Int).Mul(amounts[1], big.NewInt(10000-int64(slippage)))
			minEthAmount = minEthAmount.Div(minEthAmount, big.NewInt(10000))

			// Build sell transaction
			tx, err = wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
				GasFeeCap: uint256.MustFromBig(feeCap),
				GasTipCap: uint256.MustFromBig(tipCap),
				Gas:       swapGasLimit,
				Value:     uint256.NewInt(0),
			}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
				return router.SwapExactTokensForETH(transactOpts, randomAmount, minEthAmount, []common.Address{pair.DaiAddr, s.deploymentInfo.Weth9Addr}, wallet.GetAddress(), big.NewInt(time.Now().Unix()+300))
			})
			if err != nil {
				return nil, err
			}

			// Update balances in local cache
			if tx != nil {
				// Subtract DAI amount
				newDaiBalance := new(big.Int).Sub(tokenBalance, randomAmount)
				s.uniswap.UpdateTokenBalance(wallet.GetAddress(), pair.DaiAddr, newDaiBalance)

				// Add ETH amount
				wallet.AddBalance(amounts[1])
			}
		}
	}

	return tx, nil
}
