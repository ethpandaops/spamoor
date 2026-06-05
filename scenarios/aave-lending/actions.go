package aavelending

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"

	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
)

// healthFactorLiquidate (0.95) is the HF below which the engine attempts a
// liquidation. The margin below 1.0 is wider than one price step so an upward
// random-walk tick between building and executing the tx cannot push the target
// back above 1 (which would revert with HEALTH_FACTOR_NOT_BELOW_THRESHOLD).
var healthFactorLiquidate = big.NewInt(950000000000000000)

// tokenUnit is 10^18, the decimals of the mock reserve tokens. Used to convert
// between Aave's base-currency amounts and token amounts.
var tokenUnit = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)

// buildActionTx selects and builds the next lending transaction for the wallet
// chosen for txIdx. It returns the built transaction and a short label naming the
// action (for logging). Actions are chosen from what is actually feasible for the
// wallet's current on-chain position so the stream is varied but rarely reverts.
func (s *Scenario) buildActionTx(ctx context.Context, wallet *spamoor.Wallet, txIdx uint64, feeCap, tipCap *big.Int) (*types.Transaction, string, error) {
	rng := rand.New(rand.NewSource(int64(txIdx)*0x9e3779b1 + 1))
	walletIdx := txIdx % s.walletCount

	// periodic oracle price move (drives rate dynamics and, on dips, liquidations)
	if s.options.PriceTickInterval > 0 && txIdx > 0 && txIdx%s.options.PriceTickInterval == 0 {
		return s.buildPriceTick(ctx, wallet, feeCap, tipCap, rng)
	}

	// opportunistic liquidation of an underwater risky position (only after a dip)
	if s.options.Liquidations && s.options.RiskyRatio > 0 {
		if tx, label, ok, err := s.tryLiquidation(ctx, wallet, walletIdx, feeCap, tipCap); ok || err != nil {
			return tx, label, err
		}
	}

	if s.isRisky(walletIdx) {
		return s.buildRiskyAction(ctx, wallet, feeCap, tipCap, rng)
	}
	return s.buildNormalAction(ctx, wallet, feeCap, tipCap, rng)
}

// isRisky reports whether the wallet at walletIdx runs an aggressive near-max-LTV
// position. These wallets are the liquidation targets when prices dip.
func (s *Scenario) isRisky(walletIdx uint64) bool {
	return s.options.RiskyRatio > 0 && walletIdx%s.options.RiskyRatio == 0
}

func (s *Scenario) price(i int) *big.Int {
	s.pricesMu.RLock()
	defer s.pricesMu.RUnlock()
	return new(big.Int).Set(s.prices[i])
}

func (s *Scenario) setPrice(i int, v *big.Int) {
	s.pricesMu.Lock()
	defer s.pricesMu.Unlock()
	s.prices[i] = v
}

func (s *Scenario) txMeta(feeCap, tipCap *big.Int, gas uint64) *txbuilder.TxMetadata {
	return &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       gas,
		Value:     uint256.NewInt(0),
	}
}

// priceStepFraction is the maximum per-tick move of the persistent price random
// walk, as a fraction of the current price. Small steps make dips persist across
// several blocks so a liquidation built during a dip still executes while the
// position is underwater (independent jumps around the base price would rebound
// before the liquidation lands).
const priceStepFraction = 0.04

// priceShockChance is the probability that a price tick is a downward "crash"
// shock rather than a recovery step. Crashes create the occasional deep dip that
// makes a near-max position liquidatable.
const priceShockChance = 0.2

// priceRevertFraction is how far a non-shock tick pulls the price back toward the
// $1 base. Recovery is gradual (no upward jumps), so a dip persists for several
// blocks — long enough for a liquidation built during it to still execute while
// the position is underwater (avoiding HEALTH_FACTOR_NOT_BELOW_THRESHOLD).
const priceRevertFraction = 0.15

// buildPriceTick moves one reserve's oracle answer. Most ticks gently mean-revert
// toward the $1 base with small noise; occasionally a downward crash shock drops
// the price deep into the lower band. The gradual recovery (rather than symmetric
// jumps) keeps dips coherent so liquidations land. setAnswer is permissionless on
// the mock aggregator, so any wallet can post the move.
func (s *Scenario) buildPriceTick(ctx context.Context, wallet *spamoor.Wallet, feeCap, tipCap *big.Int, rng *rand.Rand) (*types.Transaction, string, error) {
	t := rng.Intn(len(s.deployment.Tokens))
	baseF := float64(oraclePriceAnswer)
	vol := float64(s.options.PriceVolatility) / 10000.0
	cur := float64(s.price(t).Int64())

	var np float64
	if rng.Float64() < priceShockChance {
		// downward crash into the lower half of the band: [1-vol, 1-vol/2]
		np = baseF * ((1.0 - vol) + rng.Float64()*(vol*0.5))
	} else {
		// gradual mean reversion toward base plus small symmetric noise
		np = cur + (baseF-cur)*priceRevertFraction + cur*(rng.Float64()*2-1)*priceStepFraction
	}
	// clamp to [base*(1-vol), base*(1+vol)]
	if lo := baseF * (1.0 - vol); np < lo {
		np = lo
	}
	if hi := baseF * (1.0 + vol); np > hi {
		np = hi
	}
	if np < 1 {
		np = 1
	}
	newPrice := big.NewInt(int64(np))
	s.setPrice(t, newPrice)

	tx, err := wallet.BuildBoundTx(ctx, s.txMeta(feeCap, tipCap, gasSetter), func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return s.deployment.Tokens[t].Agg.SetAnswer(opts, newPrice)
	})
	return tx, fmt.Sprintf("price-tick t%d->%s", t, newPrice.String()), err
}

// tryLiquidation looks for an underwater risky position and builds a
// liquidationCall against it. It only scans when a collateral price has dipped
// (so it costs nothing in the common case) and returns ok=false when there is
// nothing to liquidate.
func (s *Scenario) tryLiquidation(ctx context.Context, wallet *spamoor.Wallet, walletIdx uint64, feeCap, tipCap *big.Int) (*types.Transaction, string, bool, error) {
	base := big.NewInt(oraclePriceAnswer)
	dip := scaleByFloat(base, 0.97)
	if s.price(0).Cmp(dip) >= 0 && s.price(1).Cmp(dip) >= 0 {
		return nil, "", false, nil
	}

	victimIdx := s.nextRiskyVictim()
	if victimIdx < 0 || uint64(victimIdx) == walletIdx {
		return nil, "", false, nil
	}
	victim := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, victimIdx)
	if victim == nil || victim.GetAddress() == wallet.GetAddress() {
		return nil, "", false, nil
	}

	callOpts := &bind.CallOpts{Context: ctx}
	acct, err := s.deployment.Pool.GetUserAccountData(callOpts, victim.GetAddress())
	if err != nil {
		return nil, "", false, nil
	}
	// require the position to be comfortably underwater (HF < ~0.99) so a small
	// upward price move between building and executing this tx does not push it
	// back above 1 and revert with HEALTH_FACTOR_NOT_BELOW_THRESHOLD.
	if acct.TotalDebtBase.Sign() == 0 || acct.HealthFactor.Cmp(healthFactorLiquidate) >= 0 {
		return nil, "", false, nil
	}

	// pick the debt reserve the victim actually owes, and the opposite as collateral
	collIdx, debtIdx := 0, 1
	debt, err := s.deployment.Tokens[1].VarDebt.BalanceOf(callOpts, victim.GetAddress())
	if err != nil {
		return nil, "", false, nil
	}
	if debt.Sign() == 0 {
		collIdx, debtIdx = 1, 0
		debt, err = s.deployment.Tokens[0].VarDebt.BalanceOf(callOpts, victim.GetAddress())
		if err != nil || debt.Sign() == 0 {
			return nil, "", false, nil
		}
	}

	cover := scaleByFloat(debt, 0.5) // close factor: repay up to 50% of the debt
	bal, err := s.deployment.Tokens[debtIdx].Token.BalanceOf(callOpts, wallet.GetAddress())
	if err != nil {
		return nil, "", false, nil
	}
	if bal.Cmp(cover) < 0 {
		cover = bal
	}
	if cover.Sign() <= 0 {
		return nil, "", false, nil
	}

	collAsset := s.deployment.Tokens[collIdx].Addr
	debtAsset := s.deployment.Tokens[debtIdx].Addr
	victimAddr := victim.GetAddress()
	tx, err := wallet.BuildBoundTx(ctx, s.txMeta(feeCap, tipCap, s.actionGasLimit()), func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return s.deployment.Pool.LiquidationCall(opts, collAsset, debtAsset, victimAddr, cover, false)
	})
	return tx, fmt.Sprintf("liquidate w%d coll%d debt%d", victimIdx, collIdx, debtIdx), true, err
}

// nextRiskyVictim returns the index of the next risky wallet to inspect for
// liquidation, round-robining across all risky wallets, or -1 if there are none.
func (s *Scenario) nextRiskyVictim() int {
	if s.options.RiskyRatio == 0 || s.walletCount == 0 {
		return -1
	}
	riskyCount := (s.walletCount + s.options.RiskyRatio - 1) / s.options.RiskyRatio
	if riskyCount == 0 {
		return -1
	}
	k := atomic.AddUint64(&s.liqCursor, 1) % riskyCount
	return int(k * s.options.RiskyRatio)
}

// buildRiskyAction keeps a wallet's position at a high LTV: it grows token0
// collateral and borrows token1 up to its borrowable headroom (which leaves a
// small collateral buffer, settling the position at HF ~1.15). Such a position
// becomes liquidatable once cumulative price drift exceeds that buffer.
func (s *Scenario) buildRiskyAction(ctx context.Context, wallet *spamoor.Wallet, feeCap, tipCap *big.Int, rng *rand.Rand) (*types.Transaction, string, error) {
	callOpts := &bind.CallOpts{Context: ctx}
	acct, err := s.deployment.Pool.GetUserAccountData(callOpts, wallet.GetAddress())
	if err != nil {
		return nil, "", err
	}

	maxBorrow := s.borrowableTokens(acct.AvailableBorrowsBase, acct.TotalCollateralBase, s.price(1))
	if maxBorrow.Cmp(s.minAmount) >= 0 {
		amt := scaleByFloat(maxBorrow, 0.7+rng.Float64()*0.2)
		if amt.Cmp(s.minAmount) >= 0 {
			asset := s.deployment.Tokens[1].Addr
			onBehalf := wallet.GetAddress()
			tx, err := wallet.BuildBoundTx(ctx, s.txMeta(feeCap, tipCap, s.actionGasLimit()), func(opts *bind.TransactOpts) (*types.Transaction, error) {
				return s.deployment.Pool.Borrow(opts, asset, amt, big.NewInt(variableInterestRateMode), 0, onBehalf)
			})
			return tx, "risky-borrow t1", err
		}
	}

	bal0, err := s.deployment.Tokens[0].Token.BalanceOf(callOpts, wallet.GetAddress())
	if err == nil && bal0.Cmp(s.minAmount) >= 0 {
		amt := randAmount(rng, s.minAmount, minBig(s.maxAmount, bal0))
		asset := s.deployment.Tokens[0].Addr
		onBehalf := wallet.GetAddress()
		tx, err := wallet.BuildBoundTx(ctx, s.txMeta(feeCap, tipCap, s.actionGasLimit()), func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return s.deployment.Pool.Supply(opts, asset, amt, onBehalf, 0)
		})
		return tx, "risky-supply t0", err
	}

	// out of collateral headroom and out of token0 to add: fall back to a normal
	// (de-risking) action so the wallet keeps doing something useful.
	return s.buildNormalAction(ctx, wallet, feeCap, tipCap, rng)
}

// actionCandidate is one feasible action with a selection weight.
type actionCandidate struct {
	label  string
	weight int
	build  func(opts *bind.TransactOpts) (*types.Transaction, error)
}

// buildNormalAction picks a feasible action on a random reserve, weighted across
// supply/borrow/repay/withdraw, with amounts varied within the configured range
// and bounded by what the wallet can actually do.
func (s *Scenario) buildNormalAction(ctx context.Context, wallet *spamoor.Wallet, feeCap, tipCap *big.Int, rng *rand.Rand) (*types.Transaction, string, error) {
	callOpts := &bind.CallOpts{Context: ctx}
	addr := wallet.GetAddress()
	t := rng.Intn(len(s.deployment.Tokens))
	tok := s.deployment.Tokens[t]
	priceT := s.price(t)

	acct, err := s.deployment.Pool.GetUserAccountData(callOpts, addr)
	if err != nil {
		return nil, "", err
	}
	balT, _ := tok.Token.BalanceOf(callOpts, addr)
	aT, _ := tok.AToken.BalanceOf(callOpts, addr)
	dT, _ := tok.VarDebt.BalanceOf(callOpts, addr)
	if balT == nil {
		balT = big.NewInt(0)
	}
	if aT == nil {
		aT = big.NewInt(0)
	}
	if dT == nil {
		dT = big.NewInt(0)
	}

	asset := tok.Addr
	var cands []actionCandidate

	// supply: needs spare token balance
	if balT.Cmp(s.minAmount) >= 0 {
		hi := clampHi(minBig(s.maxAmount, halve(balT)), s.minAmount)
		amt := randAmount(rng, s.minAmount, hi)
		cands = append(cands, actionCandidate{"supply t" + itoa(t), 3, func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return s.deployment.Pool.Supply(opts, asset, amt, addr, 0)
		}})
	}

	// borrow: needs borrowing power (kept inside a collateral buffer)
	if acct.AvailableBorrowsBase.Sign() > 0 {
		maxBorrow := s.borrowableTokens(acct.AvailableBorrowsBase, acct.TotalCollateralBase, priceT)
		if maxBorrow.Cmp(s.minAmount) >= 0 {
			hi := clampHi(minBig(s.maxAmount, scaleByFloat(maxBorrow, 0.6)), s.minAmount)
			amt := randAmount(rng, s.minAmount, hi)
			cands = append(cands, actionCandidate{"borrow t" + itoa(t), 3, func(opts *bind.TransactOpts) (*types.Transaction, error) {
				return s.deployment.Pool.Borrow(opts, asset, amt, big.NewInt(variableInterestRateMode), 0, addr)
			}})
		}
	}

	// repay: needs outstanding debt and tokens to pay with
	if dT.Sign() > 0 && balT.Sign() > 0 {
		hi := minBig(dT, balT)
		amt := randAmount(rng, minBig(s.minAmount, hi), hi)
		if amt.Sign() > 0 {
			cands = append(cands, actionCandidate{"repay t" + itoa(t), 2, func(opts *bind.TransactOpts) (*types.Transaction, error) {
				return s.deployment.Pool.Repay(opts, asset, amt, big.NewInt(variableInterestRateMode), addr)
			}})
		}
	}

	// withdraw: needs collateral, and must keep the position healthy
	if aT.Sign() > 0 {
		safe := s.safeWithdraw(aT, priceT, acct)
		if safe.Cmp(s.minAmount) >= 0 {
			hi := clampHi(minBig(s.maxAmount, safe), s.minAmount)
			amt := randAmount(rng, s.minAmount, hi)
			cands = append(cands, actionCandidate{"withdraw t" + itoa(t), 2, func(opts *bind.TransactOpts) (*types.Transaction, error) {
				return s.deployment.Pool.Withdraw(opts, asset, amt, addr)
			}})
		}
	}

	chosen := pickCandidate(cands, rng)
	if chosen == nil {
		return s.buildFallback(ctx, wallet, feeCap, tipCap, rng, t)
	}

	tx, err := wallet.BuildBoundTx(ctx, s.txMeta(feeCap, tipCap, s.actionGasLimit()), chosen.build)
	return tx, chosen.label, err
}

// buildFallback handles the rare case where the randomly chosen reserve offered
// no feasible action: it supplies whichever reserve the wallet still holds tokens
// of, else withdraws a reserve it has collateral in. As a last resort it supplies
// the minimum on the original reserve.
func (s *Scenario) buildFallback(ctx context.Context, wallet *spamoor.Wallet, feeCap, tipCap *big.Int, rng *rand.Rand, t int) (*types.Transaction, string, error) {
	callOpts := &bind.CallOpts{Context: ctx}
	addr := wallet.GetAddress()

	for i := range s.deployment.Tokens {
		tok := s.deployment.Tokens[i]
		if bal, err := tok.Token.BalanceOf(callOpts, addr); err == nil && bal.Cmp(s.minAmount) >= 0 {
			hi := clampHi(minBig(s.maxAmount, halve(bal)), s.minAmount)
			amt := randAmount(rng, s.minAmount, hi)
			asset := tok.Addr
			tx, err := wallet.BuildBoundTx(ctx, s.txMeta(feeCap, tipCap, s.actionGasLimit()), func(opts *bind.TransactOpts) (*types.Transaction, error) {
				return s.deployment.Pool.Supply(opts, asset, amt, addr, 0)
			})
			return tx, "supply t" + itoa(i), err
		}
	}

	acct, _ := s.deployment.Pool.GetUserAccountData(callOpts, addr)
	for i := range s.deployment.Tokens {
		tok := s.deployment.Tokens[i]
		aT, err := tok.AToken.BalanceOf(callOpts, addr)
		if err != nil || aT.Sign() == 0 {
			continue
		}
		safe := s.safeWithdraw(aT, s.price(i), acct)
		if safe.Sign() <= 0 {
			continue
		}
		asset := tok.Addr
		amt := minBig(safe, s.maxAmount)
		tx, err := wallet.BuildBoundTx(ctx, s.txMeta(feeCap, tipCap, s.actionGasLimit()), func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return s.deployment.Pool.Withdraw(opts, asset, amt, addr)
		})
		return tx, "withdraw t" + itoa(i), err
	}

	// last resort: minimum supply on the original reserve (may revert if the
	// wallet is fully deployed, which is acceptable and rare).
	asset := s.deployment.Tokens[t].Addr
	amt := new(big.Int).Set(s.minAmount)
	tx, err := wallet.BuildBoundTx(ctx, s.txMeta(feeCap, tipCap, s.actionGasLimit()), func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return s.deployment.Pool.Supply(opts, asset, amt, addr, 0)
	})
	return tx, "supply t" + itoa(t), err
}

// safeWithdraw returns the largest amount of a collateral reserve the wallet can
// withdraw while keeping the position healthy. With no debt the whole aToken
// balance is withdrawable; otherwise it is bounded by the free collateral implied
// by the available borrows and LTV, with a safety margin.
func (s *Scenario) safeWithdraw(aTokenBal, priceT *big.Int, acct struct {
	TotalCollateralBase         *big.Int
	TotalDebtBase               *big.Int
	AvailableBorrowsBase        *big.Int
	CurrentLiquidationThreshold *big.Int
	Ltv                         *big.Int
	HealthFactor                *big.Int
}) *big.Int {
	if acct.TotalDebtBase.Sign() == 0 {
		return new(big.Int).Set(aTokenBal)
	}
	if acct.Ltv.Sign() == 0 || acct.AvailableBorrowsBase.Sign() <= 0 {
		return big.NewInt(0)
	}
	// freeable collateral value (base) = availableBorrows / ltv
	freeBase := new(big.Int).Div(new(big.Int).Mul(acct.AvailableBorrowsBase, big.NewInt(10000)), acct.Ltv)
	freeTokens := baseToToken(freeBase, priceT)
	safe := minBig(aTokenBal, freeTokens)
	return scaleByFloat(safe, 0.7) // margin to stay clear of the liquidation threshold
}

// pickCandidate selects one candidate weighted by its weight, or nil if empty.
func pickCandidate(cands []actionCandidate, rng *rand.Rand) *actionCandidate {
	total := 0
	for i := range cands {
		total += cands[i].weight
	}
	if total == 0 {
		return nil
	}
	r := rng.Intn(total)
	for i := range cands {
		r -= cands[i].weight
		if r < 0 {
			return &cands[i]
		}
	}
	return &cands[len(cands)-1]
}

// borrowBufferBps is the share of collateral value (in bps) the engine keeps as
// unused borrowing headroom. It absorbs the price move that can occur between
// building a borrow and it executing (a price tick lands in between), so borrows
// stay within the limit and rarely revert with COLLATERAL_CANNOT_COVER_NEW_BORROW.
// It also caps how aggressive a risky position gets: ~6% leaves HF ~1.15.
const borrowBufferBps = 600

// borrowableTokens returns how much of a reserve can be borrowed right now while
// keeping borrowBufferBps of collateral value unused, converted to token units at
// the given price. Returns 0 when there is no headroom left.
func (s *Scenario) borrowableTokens(availableBase, collateralBase, price *big.Int) *big.Int {
	buffer := new(big.Int).Div(new(big.Int).Mul(collateralBase, big.NewInt(borrowBufferBps)), big.NewInt(10000))
	spare := new(big.Int).Sub(availableBase, buffer)
	if spare.Sign() <= 0 {
		return big.NewInt(0)
	}
	return baseToToken(spare, price)
}

// baseToToken converts an Aave base-currency amount (8 decimals) to a token
// amount (18 decimals) at the given oracle price: amount = base * 1e18 / price.
func baseToToken(base, price *big.Int) *big.Int {
	if price == nil || price.Sign() == 0 {
		return big.NewInt(0)
	}
	return new(big.Int).Div(new(big.Int).Mul(base, tokenUnit), price)
}

// scaleByFloat returns x * f using 1e6 fixed-point precision for f.
func scaleByFloat(x *big.Int, f float64) *big.Int {
	if f <= 0 {
		return big.NewInt(0)
	}
	n := big.NewInt(int64(f * 1e6))
	r := new(big.Int).Mul(x, n)
	return r.Div(r, big.NewInt(1e6))
}

// randAmount returns a random amount in [lo, hi]; if hi <= lo it returns lo.
func randAmount(rng *rand.Rand, lo, hi *big.Int) *big.Int {
	if hi.Cmp(lo) <= 0 {
		return new(big.Int).Set(lo)
	}
	delta := new(big.Int).Sub(hi, lo)
	return new(big.Int).Add(lo, scaleByFloat(delta, rng.Float64()))
}

// clampHi returns hi raised to lo when it would otherwise fall below it.
func clampHi(hi, lo *big.Int) *big.Int {
	if hi.Cmp(lo) < 0 {
		return new(big.Int).Set(lo)
	}
	return hi
}

func minBig(a, b *big.Int) *big.Int {
	if a.Cmp(b) <= 0 {
		return new(big.Int).Set(a)
	}
	return new(big.Int).Set(b)
}

func halve(a *big.Int) *big.Int {
	return new(big.Int).Rsh(a, 1)
}

func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}
