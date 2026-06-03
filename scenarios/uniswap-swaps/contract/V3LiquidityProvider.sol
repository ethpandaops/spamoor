// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

// Minimal helper for the spamoor uniswap-swaps scenario (v3 mode).
//
// It seeds a full-range liquidity position into a v3 pool in a single tx,
// minting DAI on demand and wrapping ETH from msg.value inside the mint
// callback. This mirrors the PairLiquidityProvider pattern used by the v2 path.
// Swaps are NOT handled here - they are routed through the canonical Uniswap v3
// SwapRouter.

interface IERC20 {
    function transfer(address to, uint256 value) external returns (bool);
    function balanceOf(address account) external view returns (uint256);
}

interface IDai is IERC20 {
    function mint(address to, uint256 amount) external;
}

interface IWETH9 is IERC20 {
    function deposit() external payable;
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
    address private _owner1;
    address private _owner2;
    address private _weth9;

    // Expected callback caller (the pool) for the in-flight mint. Set right
    // before calling into the pool and cleared after; the callback requires
    // msg.sender to match so it cannot be invoked directly.
    address private _expectedPool;

    constructor(address owner1, address owner2, address weth9) {
        _owner1 = owner1;
        _owner2 = owner2;
        _weth9 = weth9;
    }

    receive() external payable {}

    // provideLiquidity seeds a full-range position into the given pool. DAI is
    // minted on demand and WETH is wrapped from msg.value inside the callback;
    // any unused ETH is refunded to the original sender.
    function provideLiquidity(
        address pool,
        int24 tickLower,
        int24 tickUpper,
        uint128 liquidity
    ) external payable {
        require(msg.sender == _owner1 || msg.sender == _owner2, "not owner");

        _expectedPool = pool;
        IUniswapV3Pool(pool).mint(address(this), tickLower, tickUpper, liquidity, "");
        _expectedPool = address(0);

        uint256 bal = address(this).balance;
        if (bal > 0) {
            (bool sent, ) = payable(tx.origin).call{value: bal}("");
            require(sent, "refund failed");
        }
    }

    function uniswapV3MintCallback(uint256 amount0Owed, uint256 amount1Owed, bytes calldata) external {
        require(msg.sender == _expectedPool, "unexpected caller");

        address token0 = IUniswapV3Pool(msg.sender).token0();
        address token1 = IUniswapV3Pool(msg.sender).token1();

        if (amount0Owed > 0) {
            _payMint(token0, amount0Owed);
        }
        if (amount1Owed > 0) {
            _payMint(token1, amount1Owed);
        }
    }

    function _payMint(address token, uint256 amount) private {
        if (token == _weth9) {
            IWETH9(_weth9).deposit{value: amount}();
            IWETH9(_weth9).transfer(msg.sender, amount);
        } else {
            IDai(token).mint(msg.sender, amount);
        }
    }
}
