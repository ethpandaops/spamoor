// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

// Minimal helper for the spamoor uniswap-swaps scenario (v3 mode).
//
// It seeds a full-range liquidity position into a v3 pool in a single tx,
// minting both pool tokens on demand inside the mint callback. The scenario's
// tokens are mock tokens with a public mint, so the helper never needs to be
// funded: the caller only pays gas. This mirrors the PairLiquidityProvider
// pattern used by the v2 path. Swaps are NOT handled here - they are routed
// through the canonical Uniswap v3 SwapRouter.

interface IMintableERC20 {
    function mint(address to, uint256 amount) external;
}

interface IUniswapV3Pool {
    function token0() external view returns (address);
    function token1() external view returns (address);

    function mint(
        address recipient,
        int24 tickLower,
        int24 tickUpper,
        uint128 amount,
        bytes calldata data
    ) external returns (uint256 amount0, uint256 amount1);
}

contract V3LiquidityProvider {
    address private _owner;

    // Expected callback caller (the pool) for the in-flight mint. Set right
    // before calling into the pool and cleared after; the callback requires
    // msg.sender to match so it cannot be invoked directly.
    address private _expectedPool;

    constructor(address owner) {
        _owner = owner;
    }

    // provideLiquidity seeds a position into the given pool. Both tokens are
    // minted on demand inside the callback, so no prior funding is needed.
    function provideLiquidity(
        address pool,
        int24 tickLower,
        int24 tickUpper,
        uint128 liquidity
    ) external {
        require(msg.sender == _owner, "not owner");

        _expectedPool = pool;
        IUniswapV3Pool(pool).mint(address(this), tickLower, tickUpper, liquidity, "");
        _expectedPool = address(0);
    }

    function uniswapV3MintCallback(uint256 amount0Owed, uint256 amount1Owed, bytes calldata) external {
        require(msg.sender == _expectedPool, "unexpected caller");

        if (amount0Owed > 0) {
            IMintableERC20(IUniswapV3Pool(msg.sender).token0()).mint(msg.sender, amount0Owed);
        }
        if (amount1Owed > 0) {
            IMintableERC20(IUniswapV3Pool(msg.sender).token1()).mint(msg.sender, amount1Owed);
        }
    }
}
