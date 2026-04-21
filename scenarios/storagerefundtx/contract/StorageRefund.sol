// SPDX-License-Identifier: MIT
pragma solidity ^0.8.22;

/// @title StorageRefund - Contract for triggering gas refunds via storage clearing
/// @notice Writes new storage slots and clears old ones to generate SSTORE refunds.
/// Used for testing EIP-7778 (Block Gas Accounting without Refunds) edge cases.
contract StorageRefund {
    /// @notice Next slot index to write to (monotonically increasing)
    uint256 public writePointer;

    /// @notice Next slot index to clear
    uint256 public clearPointer;

    /// @notice Slots before this index are clearable (written in a previous block)
    uint256 public clearableUpTo;

    /// @notice Last block number when clearable boundary was updated
    uint256 public lastBlockNumber;

    /// @notice Storage slots used for writing/clearing
    mapping(uint256 => uint256) public storageSlots;

    /// @notice Execute a round of storage writes and clears
    /// @param slotsPerCall Number of storage slots to write and clear per call
    function execute(uint256 slotsPerCall) external {
        // Update clearable boundary on new block
        if (block.number > lastBlockNumber) {
            clearableUpTo = writePointer;
            lastBlockNumber = block.number;
        }

        // Clear old slots (non-zero to zero = gas refund)
        uint256 toClear = clearableUpTo - clearPointer;
        if (toClear > slotsPerCall) {
            toClear = slotsPerCall;
        }
        uint256 cp = clearPointer;
        for (uint256 i = 0; i < toClear; i++) {
            delete storageSlots[cp + i];
        }
        clearPointer = cp + toClear;

        // Write new slots (zero to non-zero = expensive, no refund)
        uint256 wp = writePointer;
        for (uint256 i = 0; i < slotsPerCall; i++) {
            storageSlots[wp + i] = block.number;
        }
        writePointer = wp + slotsPerCall;
    }
}
