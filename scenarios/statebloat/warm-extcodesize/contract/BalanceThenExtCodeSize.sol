// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

contract BalanceThenExtCodeSize{
    function pokeBatch(address[] calldata addresses) external {
        for (uint i = 0; i < addresses.length; i++) {
            address a = addresses[i];
            assembly {
                let b := balance(a)
                let s := extcodesize(a)
                mstore(0x00, b)
            }
        }
    }
}