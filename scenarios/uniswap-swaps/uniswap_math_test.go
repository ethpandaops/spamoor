package uniswapswaps

import (
	"math/big"
	"testing"
)

func TestEncodeSqrtRatioX96(t *testing.T) {
	// price = amount1/amount0 = 10000 -> sqrt = 100 -> sqrtPriceX96 = 100 * 2^96
	weth := new(big.Int).Mul(big.NewInt(2000), big.NewInt(1e18))
	dai := new(big.Int).Mul(weth, big.NewInt(10000))

	got := encodeSqrtRatioX96(dai, weth)
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

func TestFullRangeLiquidityForWeth(t *testing.T) {
	weth := new(big.Int).Mul(big.NewInt(2000), big.NewInt(1e18))
	dai := new(big.Int).Mul(weth, big.NewInt(10000))

	// weth as token0
	sqrtP0 := encodeSqrtRatioX96(dai, weth)
	l0 := fullRangeLiquidityForWeth(sqrtP0, true, weth)
	if l0.Sign() <= 0 || l0.Cmp(maxUint128) > 0 {
		t.Fatalf("liquidity (weth token0) out of range: %s", l0)
	}

	// weth as token1
	sqrtP1 := encodeSqrtRatioX96(weth, dai)
	l1 := fullRangeLiquidityForWeth(sqrtP1, false, weth)
	if l1.Sign() <= 0 || l1.Cmp(maxUint128) > 0 {
		t.Fatalf("liquidity (weth token1) out of range: %s", l1)
	}
}
