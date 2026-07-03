package uniswapswaps

import (
	"math/big"
)

// Uniswap v3 fixed-point math helpers. These mirror the relevant parts of the
// Uniswap v3 SDK / TickMath / LiquidityAmounts libraries, just enough to seed a
// full-range position at a chosen starting price.

const (
	// v3 tick bounds as defined by TickMath.
	minTick int64 = -887272
	maxTick int64 = 887272
)

var (
	// q96 = 2^96, the fixed-point scaling factor for sqrtPriceX96 values.
	q96 = new(big.Int).Lsh(big.NewInt(1), 96)

	// minSqrtRatio / maxSqrtRatio are the sqrt price bounds from TickMath.
	minSqrtRatio    = big.NewInt(4295128739)
	maxSqrtRatio, _ = new(big.Int).SetString("1461446703485210103287273052203988822378723970342", 10)

	// maxUint128 bounds the liquidity value the pool can accept.
	maxUint128 = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 128), big.NewInt(1))
)

// encodeSqrtRatioX96 computes the sqrtPriceX96 representing the price
// amount1/amount0, i.e. floor(sqrt(amount1/amount0) * 2^96). amount0/amount1 are
// the reserve amounts of token0/token1 (ordered by address).
func encodeSqrtRatioX96(amount1, amount0 *big.Int) *big.Int {
	// ratioX192 = (amount1 << 192) / amount0
	numerator := new(big.Int).Lsh(amount1, 192)
	ratioX192 := new(big.Int).Div(numerator, amount0)
	return new(big.Int).Sqrt(ratioX192)
}

// fullRangeTicks returns the widest tick range aligned to tickSpacing.
func fullRangeTicks(tickSpacing int64) (lower, upper *big.Int) {
	lo := (minTick / tickSpacing) * tickSpacing
	hi := (maxTick / tickSpacing) * tickSpacing
	return big.NewInt(lo), big.NewInt(hi)
}

// getLiquidityForAmount0 returns the liquidity for a given amount of token0
// across the price range [sqrtA, sqrtB] (X96).
func getLiquidityForAmount0(sqrtA, sqrtB, amount0 *big.Int) *big.Int {
	if sqrtA.Cmp(sqrtB) > 0 {
		sqrtA, sqrtB = sqrtB, sqrtA
	}
	// intermediate = (sqrtA * sqrtB) / Q96
	intermediate := new(big.Int).Div(new(big.Int).Mul(sqrtA, sqrtB), q96)
	// liquidity = amount0 * intermediate / (sqrtB - sqrtA)
	num := new(big.Int).Mul(amount0, intermediate)
	return new(big.Int).Div(num, new(big.Int).Sub(sqrtB, sqrtA))
}

// getLiquidityForAmount1 returns the liquidity for a given amount of token1
// across the price range [sqrtA, sqrtB] (X96).
func getLiquidityForAmount1(sqrtA, sqrtB, amount1 *big.Int) *big.Int {
	if sqrtA.Cmp(sqrtB) > 0 {
		sqrtA, sqrtB = sqrtB, sqrtA
	}
	// liquidity = amount1 * Q96 / (sqrtB - sqrtA)
	num := new(big.Int).Mul(amount1, q96)
	return new(big.Int).Div(num, new(big.Int).Sub(sqrtB, sqrtA))
}

// spotAmountOut returns the output amount implied by the pool's current
// sqrtPriceX96 for an exact input swap, ignoring fee and price impact.
// zeroForOne indicates the swap direction (token0 in, token1 out).
func spotAmountOut(sqrtPriceX96, amountIn *big.Int, zeroForOne bool) *big.Int {
	// priceX192 = sqrtPriceX96^2 is the token1/token0 price scaled by 2^192.
	priceX192 := new(big.Int).Mul(sqrtPriceX96, sqrtPriceX96)
	q192 := new(big.Int).Mul(q96, q96)
	if zeroForOne {
		return new(big.Int).Div(new(big.Int).Mul(amountIn, priceX192), q192)
	}
	return new(big.Int).Div(new(big.Int).Mul(amountIn, q192), priceX192)
}

// fullRangeLiquidityForWeth computes the full-range liquidity bounded by the
// available WETH budget. DAI is minted on demand by the liquidity provider, so
// only the WETH side constrains how much liquidity can be seeded. sqrtPriceX96
// is the pool's current price and wethIsToken0 indicates the token ordering.
func fullRangeLiquidityForWeth(sqrtPriceX96 *big.Int, wethIsToken0 bool, wethBudget *big.Int) *big.Int {
	var liquidity *big.Int
	if wethIsToken0 {
		// token0 is WETH: amount0 is supplied over [current, max].
		liquidity = getLiquidityForAmount0(sqrtPriceX96, maxSqrtRatio, wethBudget)
	} else {
		// token1 is WETH: amount1 is supplied over [min, current].
		liquidity = getLiquidityForAmount1(minSqrtRatio, sqrtPriceX96, wethBudget)
	}
	if liquidity.Cmp(maxUint128) > 0 {
		liquidity = new(big.Int).Set(maxUint128)
	}
	return liquidity
}
