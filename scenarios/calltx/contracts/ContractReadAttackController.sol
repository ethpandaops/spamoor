// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract AttackController {
    // Event to emit the collected bytes
    event BytesCollected(uint256 collectedBytes);
    // Event to emit log data with target address, size, and last byte
    event LogData(address indexed target, uint256 indexed size, uint256 indexed lastByte);
    
    function extcodehashAttack(
        address factory,
        bytes32 initCodeHash,
        uint256 contractCount,
        uint256 gasBuffer
    ) external {
        uint256 seed = block.timestamp ^ block.prevrandao ^ uint256(uint160(msg.sender));
        
        while (gasleft() > gasBuffer) {
            // Generate pseudo-random salt
            seed = uint256(keccak256(abi.encode(seed, gasleft())));
            uint256 targetIndex = seed % contractCount;
            
            // Calculate target contract address
            bytes32 salt = bytes32(targetIndex);
            address target = address(uint160(uint256(keccak256(abi.encodePacked(
                bytes1(0xff),
                factory,
                salt,
                initCodeHash
            )))));
            
            // Perform extcodehash operation
            bytes32 hash;
            assembly { hash := extcodehash(target) }
        }
    }
    
    function lastByteAttack(
        address factory,
        bytes32 initCodeHash,
        uint256 contractCount,
        uint256 gasBuffer
    ) external {
        uint256 seed = block.timestamp ^ block.prevrandao ^ uint256(uint160(msg.sender));
        uint256 collectedCount;
        
        while (gasleft() > gasBuffer) {
            // Generate pseudo-random salt
            seed = uint256(keccak256(abi.encode(seed, gasleft())));
            uint256 targetIndex = seed % contractCount;
            
            // Calculate target contract address
            bytes32 salt = bytes32(targetIndex);
            address target = address(uint160(uint256(keccak256(abi.encodePacked(
                bytes1(0xff),
                factory,
                salt,
                initCodeHash
            )))));
            
            // Read last byte of contract code
            uint256 size;
            assembly { size := extcodesize(target) }
            bytes32 lastBytes;
            if (size > 0) {
                assembly {
                    // Copy the last byte to memory
                    extcodecopy(target, 0x00, sub(size, 32), 32)
                    lastBytes := mload(0x00)
                }
            }

            collectedCount++;

            // Emit log with target address, code size, and last byte
            //emit LogData(target, size, uint256(lastBytes));
        }

        emit BytesCollected(collectedCount);
    }

    /// Port of the `opcode = CALL, overhead_baseline = False` arm of
    /// `test_account_access` from execution-specs
    /// (tests/benchmark/stateful/bloatnet/test_account_query.py).
    ///
    /// `callValue = 0` measures the cold account access alone, `callValue = 1`
    /// adds the value transfer. Unlike the two attacks above, the salt walk is
    /// monotonic - `startSalt`, `startSalt + 1`, ... with no wrap-around -
    /// matching `Create2PreimageLayout.increment_salt_op` upstream. Keeping
    /// the walk inside the salt range that was actually deployed is the
    /// caller's job: it advances `startSalt` per tx from the receiver address
    /// list.
    ///
    /// `payable` so the contract can be topped up by the attack tx itself
    /// (`calltx --amount`), which it needs when `callValue > 0`.
    function callAttack(
        address factory,
        bytes32 initCodeHash,
        uint256 startSalt,
        uint256 callValue,
        uint256 gasBuffer
    ) external payable {
        uint256 salt = startSalt;

        while (gasleft() > gasBuffer) {
            address target = address(uint160(uint256(keccak256(abi.encodePacked(
                bytes1(0xff),
                factory,
                bytes32(salt),
                initCodeHash
            )))));

            // The actual attack: cold account access, plus a value transfer
            // when callValue > 0. The return value is ignored on purpose - a
            // reverting target must not abort the walk.
            assembly {
                pop(call(gas(), target, callValue, 0, 0, 0, 0))
            }

            salt++;
        }
    }

    /// Accept plain top-ups so `callValue > 0` runs do not starve.
    receive() external payable {}
}
