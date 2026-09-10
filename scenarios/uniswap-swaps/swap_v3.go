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
	"github.com/ethpandaops/spamoor/txtypes"
)

// buildV3SwapTx builds a single exact-input swap against a randomly selected
// DAI/quote pool, routed through one of the two SwapRouters (which picks the
// matching factory's pool). Swap sizes are DAI-denominated to match the v2
// path. There is no quoter deployed, so both the required input and the
// amountOutMinimum are derived from the pool's current spot price with the pool
// fee and the configured slippage tolerance applied. A wallet that cannot
// afford a buy mints itself more quote tokens instead of swapping.
func (s *Scenario) buildV3SwapTx(ctx context.Context, wallet *spamoor.Wallet, feeCap, tipCap *big.Int) (*txtypes.Transaction, error) {
	info := s.uniswap.v3Deployment
	if info == nil || len(info.Pools) == 0 {
		return nil, fmt.Errorf("no v3 pools deployed")
	}

	poolInfo := info.Pools[mathrand.Intn(len(info.Pools))]
	daiAddr := poolInfo.DaiAddr
	quoteAddr := info.QuoteAddr

	// alternate between the two routers (factory A vs B pool)
	router := info.RouterA
	pool := poolInfo.PoolA
	if mathrand.Intn(100) < 50 {
		router = info.RouterB
		pool = poolInfo.PoolB
	}

	amount := s.randomSwapAmount()
	slippage := s.perTradeSlippage()

	addr := wallet.GetAddress()
	daiBalance := s.uniswap.GetTokenBalance(addr, daiAddr)
	quoteBalance := s.uniswap.GetTokenBalance(addr, quoteAddr)
	isBuy := s.decideSwapSide(daiBalance, amount)

	// Spot price of the routed pool. Price impact and price movement until
	// execution must fit into the slippage tolerance, matching how the v2 path
	// quotes before applying it.
	slot0, err := pool.Slot0(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, fmt.Errorf("could not read pool slot0: %w", err)
	}
	// zeroForOne for a swap that spends quote and receives DAI
	quoteForDai := poolInfo.QuoteIsToken0

	deadline := big.NewInt(time.Now().Unix() + 300)
	buildSwap := func(tokenIn, tokenOut common.Address, amountIn, amountOutMinimum *big.Int) (*txtypes.Transaction, error) {
		params := contract.ISwapRouterExactInputSingleParams{
			TokenIn:           tokenIn,
			TokenOut:          tokenOut,
			Fee:               info.Fee,
			Recipient:         addr,
			Deadline:          deadline,
			AmountIn:          amountIn,
			AmountOutMinimum:  amountOutMinimum,
			SqrtPriceLimitX96: big.NewInt(0),
		}
		return wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Gas:       swapGasLimit,
			Value:     uint256.NewInt(0),
		}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
			return router.ExactInputSingle(transactOpts, params)
		})
	}

	if isBuy {
		// Buying DAI with quote tokens: quote input needed for the desired DAI
		// amount at spot (grossed up for the pool fee), floor = spot output of
		// that input minus slippage.
		quoteIn := spotAmountIn(slot0.SqrtPriceX96, info.Fee, amount, quoteForDai)
		if quoteBalance.Cmp(quoteIn) < 0 {
			// out of quote tokens: top up instead of swapping this round
			s.logger.WithField("wallet", s.walletPool.GetWalletName(addr)).Debugf("quote balance %s below %s needed for buy, minting funding", quoteBalance, quoteIn)
			return s.uniswap.buildQuoteMintTx(ctx, wallet, feeCap, tipCap)
		}

		minDaiOut := applySlippage(spotAmountOutAfterFee(slot0.SqrtPriceX96, info.Fee, quoteIn, quoteForDai), slippage)
		tx, err := buildSwap(quoteAddr, daiAddr, quoteIn, minDaiOut)
		if err != nil {
			return nil, err
		}

		// track the guaranteed floor so the cache never overstates holdings
		s.uniswap.UpdateTokenBalance(addr, quoteAddr, new(big.Int).Sub(quoteBalance, quoteIn))
		s.uniswap.UpdateTokenBalance(addr, daiAddr, new(big.Int).Add(daiBalance, minDaiOut))
		return tx, nil
	}

	// Selling DAI for quote tokens.
	minQuoteOut := applySlippage(spotAmountOutAfterFee(slot0.SqrtPriceX96, info.Fee, amount, !quoteForDai), slippage)
	tx, err := buildSwap(daiAddr, quoteAddr, amount, minQuoteOut)
	if err != nil {
		return nil, err
	}

	s.uniswap.UpdateTokenBalance(addr, daiAddr, new(big.Int).Sub(daiBalance, amount))
	s.uniswap.UpdateTokenBalance(addr, quoteAddr, new(big.Int).Add(quoteBalance, minQuoteOut))
	return tx, nil
}
