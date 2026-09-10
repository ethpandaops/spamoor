package uniswapswaps

import (
	"context"
	"math/big"
	mathrand "math/rand"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"

	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/txtypes"
)

// buildV2SwapTx builds a single uniswap v2 swap against a randomly selected
// DAI/quote pair and router, deciding buy vs sell from the configured ratio and
// the wallet's tracked balances. Amounts are DAI-denominated and priced via the
// router's on-chain getAmountsIn/getAmountsOut helpers. A wallet that cannot
// afford a buy mints itself more quote tokens instead of swapping.
func (s *Scenario) buildV2SwapTx(ctx context.Context, wallet *spamoor.Wallet, feeCap, tipCap *big.Int) (*txtypes.Transaction, error) {
	info := s.deploymentInfo
	pair := info.Pairs[mathrand.Intn(len(info.Pairs))]
	daiAddr := pair.DaiAddr
	quoteAddr := info.QuoteAddr

	// alternate between the two routers (factory A vs B pair)
	router := s.uniswap.RouterA
	if mathrand.Intn(100) < 50 {
		router = s.uniswap.RouterB
	}

	amount := s.randomSwapAmount()
	slippage := s.perTradeSlippage()

	addr := wallet.GetAddress()
	daiBalance := s.uniswap.GetTokenBalance(addr, daiAddr)
	quoteBalance := s.uniswap.GetTokenBalance(addr, quoteAddr)
	isBuy := s.decideSwapSide(daiBalance, amount)

	callOpts := &bind.CallOpts{Context: ctx}
	deadline := big.NewInt(time.Now().Unix() + 300)

	buildSwap := func(amountIn, amountOutMin *big.Int, path []common.Address) (*txtypes.Transaction, error) {
		return wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Gas:       swapGasLimit,
			Value:     uint256.NewInt(0),
		}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
			return router.SwapExactTokensForTokens(transactOpts, amountIn, amountOutMin, path, addr, deadline)
		})
	}

	if isBuy {
		// Buying DAI with quote tokens: quote the input needed for the desired
		// DAI amount, then swap exact input with a slippage-adjusted floor.
		path := []common.Address{quoteAddr, daiAddr}
		amounts, err := router.GetAmountsIn(callOpts, amount, path)
		if err != nil {
			return nil, err
		}
		quoteIn := amounts[0]

		if quoteBalance.Cmp(quoteIn) < 0 {
			// out of quote tokens: top up instead of swapping this round
			s.logger.WithField("wallet", s.walletPool.GetWalletName(addr)).Debugf("quote balance %s below %s needed for buy, minting funding", quoteBalance, quoteIn)
			return s.uniswap.buildQuoteMintTx(ctx, wallet, feeCap, tipCap)
		}

		tx, err := buildSwap(quoteIn, applySlippage(amount, slippage), path)
		if err != nil {
			return nil, err
		}

		s.uniswap.UpdateTokenBalance(addr, quoteAddr, new(big.Int).Sub(quoteBalance, quoteIn))
		s.uniswap.UpdateTokenBalance(addr, daiAddr, new(big.Int).Add(daiBalance, amount))
		return tx, nil
	}

	// Selling DAI for quote tokens.
	path := []common.Address{daiAddr, quoteAddr}
	amounts, err := router.GetAmountsOut(callOpts, amount, path)
	if err != nil {
		return nil, err
	}
	quoteOut := amounts[1]

	tx, err := buildSwap(amount, applySlippage(quoteOut, slippage), path)
	if err != nil {
		return nil, err
	}

	s.uniswap.UpdateTokenBalance(addr, daiAddr, new(big.Int).Sub(daiBalance, amount))
	s.uniswap.UpdateTokenBalance(addr, quoteAddr, new(big.Int).Add(quoteBalance, quoteOut))
	return tx, nil
}
