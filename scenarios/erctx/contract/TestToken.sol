// SPDX-License-Identifier: MIT
pragma solidity ^0.8.22;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";

contract TestToken is ERC20, ERC20Burnable {

    constructor() ERC20("TestToken", "TT") {
    }

    function mint(uint256 amount) public {
        _mint(_msgSender(), amount);
    }

    function transferMint(address recipient, uint256 amount) public returns (bool) {
        address owner = _msgSender();
        _mint(owner, amount);
        _transfer(owner, recipient, amount);
        return true;
    }

}
