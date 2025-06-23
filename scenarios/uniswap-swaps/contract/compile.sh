#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../../../scripts/compile-contract.sh"
cd $SCRIPT_DIR

# Uniswap V2 Router
compile_contract "$(pwd)" 0.6.6 "--optimize --optimize-runs 999999" UniswapV2Router02

# WETH9
compile_contract "$(pwd)" 0.4.18 "--optimize --optimize-runs 200" WETH9

# Dai
compile_contract "$(pwd)" 0.5.12 "" Dai

# PairLiquidityProvider
compile_contract "$(pwd)" 0.8.17 "--optimize --optimize-runs 200" PairLiquidityProvider
