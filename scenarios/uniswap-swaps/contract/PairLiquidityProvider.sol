// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

// Minimal helper for the spamoor uniswap-swaps scenario (v2 mode).
//
// It seeds an ERC20/ERC20 pair on both routers in a single tx, minting both
// tokens on demand. The scenario's tokens are mock tokens with a public mint,
// so the helper never needs to be funded: the caller only pays gas.

interface IMintableERC20 {
    function approve(address spender, uint256 value) external returns (bool);
    function mint(address to, uint256 amount) external;
}

interface IUniswapV2Router02 {
    function addLiquidity(
        address tokenA,
        address tokenB,
        uint amountADesired,
        uint amountBDesired,
        uint amountAMin,
        uint amountBMin,
        address to,
        uint deadline
    ) external returns (uint amountA, uint amountB, uint liquidity);
}

contract PairLiquidityProvider {
    address private _owner;
    address private _router1;
    address private _router2;
    mapping(address => mapping(address => bool)) private _deployedLiquidity;

    constructor(address owner, address router1, address router2) {
        _owner = owner;
        _router1 = router1;
        _router2 = router2;
    }

    // providePairLiquidity mints amountA of tokenA and amountB of tokenB and
    // splits them evenly between the tokenA/tokenB pair on router1 and the one
    // on router2 (creating the pairs if needed). The LP tokens stay locked in
    // this contract.
    function providePairLiquidity(address tokenA, address tokenB, uint256 amountA, uint256 amountB) external {
        require(msg.sender == _owner, "not owner");
        require(!_deployedLiquidity[tokenA][tokenB], "liquidity already deployed");
        _deployedLiquidity[tokenA][tokenB] = true;

        uint256 amountAPerPair = amountA / 2;
        uint256 amountBPerPair = amountB / 2;

        IMintableERC20(tokenA).mint(address(this), amountAPerPair * 2);
        IMintableERC20(tokenB).mint(address(this), amountBPerPair * 2);

        _addLiquidity(_router1, tokenA, tokenB, amountAPerPair, amountBPerPair);
        _addLiquidity(_router2, tokenA, tokenB, amountAPerPair, amountBPerPair);
    }

    function _addLiquidity(address router, address tokenA, address tokenB, uint256 amountA, uint256 amountB) private {
        require(IMintableERC20(tokenA).approve(router, amountA), "approve A failed");
        require(IMintableERC20(tokenB).approve(router, amountB), "approve B failed");
        IUniswapV2Router02(router).addLiquidity(tokenA, tokenB, amountA, amountB, 0, 0, address(this), block.timestamp);
    }
}
