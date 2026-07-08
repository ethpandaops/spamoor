#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../../../scripts/compile-contract.sh"
cd $SCRIPT_DIR

# MintableToken (mock 18-decimal stablecoin)
compile_contract "$(pwd)" 0.8.17 "--optimize --optimize-runs 200" MintableToken

# CurveLiquidityProvider (liquidity seeding helper)
compile_contract "$(pwd)" 0.8.17 "--optimize --optimize-runs 200" CurveLiquidityProvider

# Curve StableSwap 3pool + LP token (canonical Curve Vyper sources, see *.vy headers
# for the two minimal spamoor edits: 18-decimal rates and a parameterized minter).
compile_vyper 0.2.4 StableSwap StableSwap
compile_vyper 0.2.12 CurveToken CurveToken
