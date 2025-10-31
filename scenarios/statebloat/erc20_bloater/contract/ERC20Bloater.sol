// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title ERC20Bloater
 * @dev Minimal ERC20 contract optimized for state bloat benchmarking
 * Removes all requires and events from transfer/approve to benchmark pure SSTORE/SLOAD
 */
contract ERC20Bloater {
    string public constant name = "BloatToken";
    string public constant symbol = "BLOAT";
    uint8 public constant decimals = 18;
    uint256 public totalSupply;
    uint256 public nextStorageSlot;

    mapping(address => uint256) public balanceOf;
    mapping(address => mapping(address => uint256)) public allowance;

    event StorageBloated(uint256 startSlot, uint256 endSlot, uint256 slotsWritten);

    constructor(uint256 initialSupply) {
        totalSupply = initialSupply;
        balanceOf[msg.sender] = initialSupply;
        nextStorageSlot = 1; // Start from address 0x0000...0001
    }

    /**
     * @dev Transfer tokens - NO requires, NO events for pure SSTORE benchmarking
     */
    function transfer(address to, uint256 amount) external returns (bool) {
        unchecked {
            balanceOf[msg.sender] -= amount;
            balanceOf[to] += amount;
        }
        return true;
    }

    /**
     * @dev Approve spender - NO requires, NO events for pure SSTORE benchmarking
     */
    function approve(address spender, uint256 amount) external returns (bool) {
        unchecked {
            allowance[msg.sender][spender] = amount;
        }
        return true;
    }

    /**
     * @dev TransferFrom - included for completeness but not used in bloating
     */
    function transferFrom(address from, address to, uint256 amount) external returns (bool) {
        unchecked {
            allowance[from][msg.sender] -= amount;
            balanceOf[from] -= amount;
            balanceOf[to] += amount;
        }
        return true;
    }

    /**
     * @dev Bloat storage by transferring tokens and approving sequential addresses
     * This creates new storage slots in both balanceOf and allowance mappings
     * @param startSlot The address index to start from (allows resuming after errors)
     * @param numAddresses Number of addresses to bloat (each creates 2 storage slots)
     */
    function bloatStorage(uint256 startSlot, uint256 numAddresses) external {
        uint256 endSlot = startSlot + numAddresses;

        unchecked {
            for (uint256 i = startSlot; i < endSlot; i++) {
                // Generate sequential address: 0x0000...0001, 0x0000...0002, etc.
                address targetAddr = address(uint160(i));

                // 1. Transfer i tokens to the address (creates balanceOf[targetAddr] slot)
                balanceOf[msg.sender] -= i;
                balanceOf[targetAddr] += i;

                // 2. Approve the address (creates allowance[msg.sender][targetAddr] slot)
                allowance[msg.sender][targetAddr] = i;
            }
        }

        nextStorageSlot = endSlot;
        emit StorageBloated(startSlot, endSlot, numAddresses * 2); // 2 slots per address
    }

    /**
     * @dev Get the current bloating progress
     */
    function getBloatProgress() external view returns (uint256) {
        return nextStorageSlot;
    }

    /**
     * @dev Emergency mint function to refill supply if needed
     */
    function mint(address to, uint256 amount) external {
        unchecked {
            totalSupply += amount;
            balanceOf[to] += amount;
        }
    }
}
