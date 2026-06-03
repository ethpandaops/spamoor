// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

interface IMintableToken {
    function mint(address to, uint256 amount) external;
    function approve(address spender, uint256 amount) external returns (bool);
}

interface IStableSwap {
    function add_liquidity(uint256[3] memory amounts, uint256 min_mint_amount) external;
}

// Seeds balanced liquidity into a StableSwap pool in a single transaction: it
// mints `amount` of each coin to itself, approves the pool, and adds liquidity.
// Doing it atomically lets the deployer seed pools with one estimable tx per
// pool, mirroring the PairLiquidityProvider pattern used by uniswap-swaps.
contract CurveLiquidityProvider {
    function seedLiquidity(address pool, address[3] calldata coins, uint256 amount) external {
        uint256[3] memory amounts;
        for (uint256 i = 0; i < 3; i++) {
            IMintableToken(coins[i]).mint(address(this), amount);
            IMintableToken(coins[i]).approve(pool, amount);
            amounts[i] = amount;
        }
        IStableSwap(pool).add_liquidity(amounts, 0);
    }
}
