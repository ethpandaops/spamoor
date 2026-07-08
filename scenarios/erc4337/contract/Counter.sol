// SPDX-License-Identifier: MIT
pragma solidity ^0.8.23;

/// @title Counter
/// @notice Trivial state-writing target for ERC-4337 UserOperations. Each
///         distinct caller (every fresh smart account) writes its own mapping
///         slot, so a stream of ops produces steady state growth - useful for
///         stress testing.
contract Counter {
    uint256 public total;
    mapping(address => uint256) public counts;

    function increment() external {
        total += 1;
        counts[msg.sender] += 1;
    }
}
