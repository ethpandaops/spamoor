// SPDX-License-Identifier: MIT
pragma solidity ^0.8.23;

// Minimal mirror of the ERC-4337 v0.7 PackedUserOperation struct so this file
// stays self-contained (no external imports) and compiles with the single-file
// solc helper used across spamoor scenarios.
struct PackedUserOperation {
    address sender;
    uint256 nonce;
    bytes initCode;
    bytes callData;
    bytes32 accountGasLimits;
    uint256 preVerificationGas;
    bytes32 gasFees;
    bytes paymasterAndData;
    bytes signature;
}

interface IStakeManager {
    function depositTo(address account) external payable;
    function balanceOf(address account) external view returns (uint256);
}

/// @title AcceptAllPaymaster
/// @notice Minimal ERC-4337 v0.7 paymaster that sponsors every UserOperation
///         routed through it. It performs NO validation and pays for any op out
///         of its EntryPoint deposit. Intended solely for testnet load
///         generation - never deploy this on a network with real value.
contract AcceptAllPaymaster {
    IStakeManager public immutable entryPoint;

    constructor(address _entryPoint) {
        entryPoint = IStakeManager(_entryPoint);
    }

    /// @notice Forwards msg.value to the EntryPoint as this paymaster's deposit.
    function deposit() external payable {
        entryPoint.depositTo{value: msg.value}(address(this));
    }

    /// @notice Returns this paymaster's current EntryPoint deposit balance.
    function getDeposit() external view returns (uint256) {
        return entryPoint.balanceOf(address(this));
    }

    /// @notice Accepts every UserOperation without checks. validationData == 0
    ///         signals "valid forever, no aggregator". Returns an empty context
    ///         so the EntryPoint skips the postOp callback.
    function validatePaymasterUserOp(PackedUserOperation calldata, bytes32, uint256)
        external
        view
        returns (bytes memory context, uint256 validationData)
    {
        require(msg.sender == address(entryPoint), "Paymaster: not from EntryPoint");
        return ("", 0);
    }

    /// @notice No-op postOp. Only reachable if a non-empty context is returned,
    ///         which this paymaster never does; kept for interface completeness.
    function postOp(uint8, bytes calldata, uint256, uint256) external view {
        require(msg.sender == address(entryPoint), "Paymaster: not from EntryPoint");
    }

    /// @notice Plain ETH transfers top up the EntryPoint deposit too.
    receive() external payable {
        entryPoint.depositTo{value: msg.value}(address(this));
    }
}
