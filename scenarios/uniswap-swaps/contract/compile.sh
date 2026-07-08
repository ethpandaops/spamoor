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

# Uniswap V3 core + periphery (canonical precompiled artifacts). Using the
# published bytecode keeps the pool init code hash consistent with the
# SwapRouter's hardcoded POOL_INIT_CODE_HASH (0xe34f199b...).
gen_from_artifact "https://unpkg.com/@uniswap/v3-core@1.0.1/artifacts/contracts/UniswapV3Factory.sol/UniswapV3Factory.json" UniswapV3Factory
gen_from_artifact "https://unpkg.com/@uniswap/v3-core@1.0.1/artifacts/contracts/UniswapV3Pool.sol/UniswapV3Pool.json" UniswapV3Pool
gen_from_artifact "https://unpkg.com/@uniswap/v3-periphery@1.4.4/artifacts/contracts/SwapRouter.sol/SwapRouter.json" SwapRouter

# V3LiquidityProvider (liquidity seeding + swap helper for v3 mode)
compile_contract "$(pwd)" 0.8.17 "--optimize --optimize-runs 200" V3LiquidityProvider
