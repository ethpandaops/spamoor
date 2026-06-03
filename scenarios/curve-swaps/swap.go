package curveswaps

import (
	"context"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"

	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
)

// feeDenominator matches the StableSwap contract's FEE_DENOMINATOR (1e10) and is
// used to size min_dy / track balances net of the swap fee.
var feeDenominator = big.NewInt(1e10)

// buildSwapTx builds a single StableSwap exchange against a randomly selected
// pool. To keep both wallets and pools balanced it always sells the coin the
// wallet holds most of, swapping it into a random other coin. Stablecoin output
// is ~1:1 with input (minus fee), so min_dy is derived locally from the input
// amount and the configured slippage rather than via an on-chain quote.
func (s *Scenario) buildSwapTx(ctx context.Context, wallet *spamoor.Wallet, feeCap, tipCap *big.Int) (*types.Transaction, error) {
	pools := s.deployment.Pools
	if len(pools) == 0 {
		return nil, fmt.Errorf("no pools deployed")
	}
	pool := pools[mathrand.Intn(len(pools))]

	minAmount, ok := new(big.Int).SetString(s.options.MinSwapAmount, 10)
	if !ok {
		return nil, fmt.Errorf("invalid min swap amount: %s", s.options.MinSwapAmount)
	}
	maxAmount, ok := new(big.Int).SetString(s.options.MaxSwapAmount, 10)
	if !ok {
		return nil, fmt.Errorf("invalid max swap amount: %s", s.options.MaxSwapAmount)
	}

	addr := wallet.GetAddress()

	// pick the input coin as the one the wallet holds most of, so swaps drain the
	// accumulated coin and naturally rebalance both the wallet and the pool.
	balances := make([]*big.Int, coinsPerPool)
	maxIdx := 0
	for k := 0; k < coinsPerPool; k++ {
		balances[k] = s.curve.GetTokenBalance(addr, pool.Coins[k])
		if balances[k].Cmp(balances[maxIdx]) > 0 {
			maxIdx = k
		}
	}
	if balances[maxIdx].Sign() == 0 {
		return nil, fmt.Errorf("wallet has no coin balance to swap")
	}
	i := maxIdx

	// output coin: a random coin different from the input.
	j := mathrand.Intn(coinsPerPool)
	for j == i {
		j = mathrand.Intn(coinsPerPool)
	}

	// random swap size in [min, max], capped to the available input balance.
	diff := new(big.Int).Sub(maxAmount, minAmount)
	dx := new(big.Int).Set(minAmount)
	if diff.Sign() > 0 {
		dx.Add(dx, new(big.Int).Rand(mathrand.New(mathrand.NewSource(time.Now().UnixNano())), diff))
	}
	if dx.Cmp(balances[i]) > 0 {
		dx.Set(balances[i])
	}
	if dx.Sign() == 0 {
		return nil, fmt.Errorf("swap amount resolved to zero")
	}

	// expected output ~= dx minus the swap fee (stablecoin ~1:1).
	swapFee := new(big.Int).SetUint64(s.options.Fee)
	dyAfterFee := new(big.Int).Sub(dx, new(big.Int).Div(new(big.Int).Mul(dx, swapFee), feeDenominator))

	// min_dy applies the configured slippage tolerance on top of the fee.
	minDy := new(big.Int).Div(
		new(big.Int).Mul(dyAfterFee, big.NewInt(int64(10000-s.options.Slippage))),
		big.NewInt(10000),
	)

	iIdx := big.NewInt(int64(i))
	jIdx := big.NewInt(int64(j))

	tx, err := wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       swapGasLimit,
		Value:     uint256.NewInt(0),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		return pool.Pool.Exchange(transactOpts, iIdx, jIdx, dx, minDy)
	})
	if err != nil {
		return nil, err
	}

	// update the local balance cache with the conservative estimate.
	s.curve.UpdateTokenBalance(addr, pool.Coins[i], new(big.Int).Sub(balances[i], dx))
	s.curve.UpdateTokenBalance(addr, pool.Coins[j], new(big.Int).Add(balances[j], dyAfterFee))

	return tx, nil
}
