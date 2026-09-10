package uniswapswaps

import (
	"math/big"
	"testing"
)

func TestEncodeSqrtRatioX96(t *testing.T) {
	// price = amount1/amount0 = 10000 -> sqrt = 100 -> sqrtPriceX96 = 100 * 2^96
	quote := new(big.Int).Mul(big.NewInt(2000), big.NewInt(1e18))
	dai := new(big.Int).Mul(quote, big.NewInt(10000))

	got := encodeSqrtRatioX96(dai, quote)
	want := new(big.Int).Mul(big.NewInt(100), q96)
	if got.Cmp(want) != 0 {
		t.Fatalf("encodeSqrtRatioX96 = %s, want %s", got, want)
	}
}

func TestFullRangeTicks(t *testing.T) {
	lower, upper := fullRangeTicks(60)
	if lower.Int64() != -887220 || upper.Int64() != 887220 {
		t.Fatalf("fullRangeTicks(60) = (%d, %d), want (-887220, 887220)", lower.Int64(), upper.Int64())
	}
}

func TestSpotAmountOut(t *testing.T) {
	// price = token1/token0 = 10000
	quote := new(big.Int).Mul(big.NewInt(2000), big.NewInt(1e18))
	dai := new(big.Int).Mul(quote, big.NewInt(10000))
	sqrtP := encodeSqrtRatioX96(dai, quote)

	// token0 -> token1: 1 token0 yields 10000 token1
	in := big.NewInt(1e18)
	out := spotAmountOut(sqrtP, in, true)
	want := new(big.Int).Mul(big.NewInt(10000), big.NewInt(1e18))
	if out.Cmp(want) != 0 {
		t.Fatalf("spotAmountOut(zeroForOne) = %s, want %s", out, want)
	}

	// token1 -> token0: 10000 token1 yield 1 token0
	out = spotAmountOut(sqrtP, want, false)
	if out.Cmp(in) != 0 {
		t.Fatalf("spotAmountOut(oneForZero) = %s, want %s", out, in)
	}
}

func TestSpotAmountInRoundTrip(t *testing.T) {
	// price = token1/token0 = 10000, fee = 0.3%
	quote := new(big.Int).Mul(big.NewInt(2000), big.NewInt(1e18))
	dai := new(big.Int).Mul(quote, big.NewInt(10000))
	sqrtP := encodeSqrtRatioX96(dai, quote)
	fee := big.NewInt(3000)

	// buying 10000 token1 costs 1 token0 plus the 0.3% fee on the input
	wantOut := new(big.Int).Mul(big.NewInt(10000), big.NewInt(1e18))
	in := spotAmountIn(sqrtP, fee, wantOut, true)
	if in.Cmp(big.NewInt(1e18)) <= 0 {
		t.Fatalf("spotAmountIn = %s, want > 1e18 (fee must be included)", in)
	}

	// swapping that input back through the fee model yields at least the target
	out := spotAmountOutAfterFee(sqrtP, fee, in, true)
	if out.Cmp(wantOut) < 0 {
		t.Fatalf("spotAmountOutAfterFee(spotAmountIn(x)) = %s, want >= %s", out, wantOut)
	}

	// the rounding slack is at most one output unit per input unit of price
	slack := new(big.Int).Sub(out, wantOut)
	if slack.Cmp(big.NewInt(10000)) > 0 {
		t.Fatalf("spotAmountIn overshoots by %s, want <= 10000", slack)
	}

	// tiny outputs never round down to a zero input
	if spotAmountIn(sqrtP, fee, big.NewInt(1), true).Sign() <= 0 {
		t.Fatalf("spotAmountIn must be at least 1")
	}
}

func TestFullRangeLiquidityForToken(t *testing.T) {
	quote := new(big.Int).Mul(big.NewInt(2000), big.NewInt(1e18))
	dai := new(big.Int).Mul(quote, big.NewInt(10000))

	// quote as token0
	sqrtP0 := encodeSqrtRatioX96(dai, quote)
	l0 := fullRangeLiquidityForToken(sqrtP0, true, quote)
	if l0.Sign() <= 0 || l0.Cmp(maxUint128) > 0 {
		t.Fatalf("liquidity (quote token0) out of range: %s", l0)
	}

	// quote as token1
	sqrtP1 := encodeSqrtRatioX96(quote, dai)
	l1 := fullRangeLiquidityForToken(sqrtP1, false, quote)
	if l1.Sign() <= 0 || l1.Cmp(maxUint128) > 0 {
		t.Fatalf("liquidity (quote token1) out of range: %s", l1)
	}
}
