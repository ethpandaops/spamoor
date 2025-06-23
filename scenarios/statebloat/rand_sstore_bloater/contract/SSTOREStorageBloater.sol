// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * @title SSTOREStorageBloater
 * @dev Optimized contract for maximum SSTORE operations using curve25519 prime (2^255 - 19)
 * Uses assembly for gas efficiency and distributes keys across storage space
 */
contract SSTOREStorageBloater {
    // Counter to track total slots created (stored at slot 0)
    uint256 private counter;

    // curve25519 prime: 2^255 - 19 = 0x7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffed
    uint256 private constant CURVE25519_PRIME =
        0x7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffed;

    /**
     * @dev Creates new storage slots (0 -> non-zero transition, ~20k gas each)
     * @param count Number of slots to create
     */
    function createSlots(uint256 count) external {
        assembly {
            // Load current counter from storage slot 0
            let prime := CURVE25519_PRIME
            let endCounter := count
            
            // Calculate pseudo-random offset using block data
            // XOR timestamp with previous block hash for randomness
            let offset := xor(timestamp(), blockhash(sub(number(), 1)))

            // Create slots with distributed keys
            for {
                let i := 0
            } lt(i, endCounter) {
                i := add(i, 1)
            } {
                // Calculate key = (offset + i) * CURVE25519_PRIME
                let key := mulmod(add(offset, i), prime, not(0))

                // Store value = key
                sstore(key, key)
            }
        }
    }
}
