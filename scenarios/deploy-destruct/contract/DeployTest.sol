// SPDX-License-Identifier: MIT
pragma solidity  ^0.8.22;

contract DeployTest {
    event TestSeed(uint seed);

    uint public counter = 1;
    uint public childIdx = 1;
    mapping(uint => address) public childAddresses;
    mapping(address => uint) childAddressIndex;

    function random() private returns (uint) {
        counter++;
        return uint(keccak256(abi.encodePacked(block.difficulty, block.timestamp, counter)));
    }

    function random2(uint seed, uint id1, uint id2) private pure returns (uint) {
        return uint(keccak256(abi.encodePacked(seed, id1, id2)));
    }

    function childCode() public pure returns (bytes memory) {
        return type(DeployTestChild).creationCode;
    }

    function create(uint256 amount, bytes memory bytecode) private returns (address) {
        address addr;
        assembly {
            addr := create(amount, add(bytecode, 0x20), mload(bytecode))

            if iszero(extcodesize(addr)) {
                revert(0, 0)
            }
        }
        require(addr != address(0), "create failed");
        return addr;
    }

    function notify(address addr, bool destroy) public {
        if(destroy) {
            uint index = childAddressIndex[addr];
            if(index > 0) {
                childAddressIndex[addr] = 0;
                childIdx--;
                childAddresses[index] = childAddresses[childIdx];
                childAddressIndex[childAddresses[index]] = index;
            }
        } else {
            childAddresses[childIdx] = addr;
            childAddressIndex[addr] = childIdx;
            childIdx++;
        }
    }

    function test(uint seed) public payable {
        if(seed == 0) {
            seed = random() % 100000;
        }

        emit TestSeed(seed);

        uint s1 = random2(seed, 0, 1);
        if(s1 % 50 < 30) {
            // deploy new contract
            bytes memory code = abi.encodePacked(childCode(), abi.encode(s1, 0, address(this)));
            address addr = create(address(this).balance, code);
            notify(addr, false);
        } else if(s1 % 50 < 40) {
            // destruct child contracts
            for(uint i = 0; i < 50; i++) {
                if(childIdx > 1 && gasleft() > 400000) {
                    address child = childAddresses[1];
                    DeployTestChild(child).destroy(random2(seed, 1, i));
                    notify(child, true);
                }
            }
        } else {
            // call child contracts
            uint256 value = address(this).balance;
            for(uint i = 1; i < childIdx; i++) {
                address child = childAddresses[i];
                uint s3 = random2(seed, 0, i);
                uint256 childValue = value / 100 * (s3 % 80);
                value -= childValue;
                DeployTestChild(child).call{value: childValue}(s3);
            }
        }
    }

    function clean(uint count) public {
        while(count > 0) {
            if(childIdx > 1) {
                address child = childAddresses[1];
                DeployTestChild(child).destroy(5);
                notify(child, true);
            }
            count--;
        }
    }

}

contract DeployTestChild {
    address public _main;
    uint public _seed;

    event ChildCreated(uint seed, uint depth);
    event ChildDestroyed(uint seed, address target);
    event CatchedRevert(uint seed, uint situation, bytes data);

    constructor(uint seed, uint depth, address main) payable {
        emit ChildCreated(seed, depth);

        _main = main;
        _seed = seed;
        call(depth);
    }

    function call(uint depth) public payable {
        uint s2 = random2(_seed, depth, 0);
        uint256 value = address(this).balance;
        if(depth < 4 && s2 % 100 < 60) {
            // create up to 3 nested contracts (chance: 60%)
            bytes memory initCode = DeployTest(_main).childCode();
            for(uint i = 0; i < (s2 % 3)+1; i++) {
                if(gasleft() < 2000000) {
                    break;
                }

                uint s3 = random2(_seed, depth, i);
                bytes memory code = abi.encodePacked(initCode, abi.encode(s3, depth + 1, _main));

                s3 = random2(_seed+1, depth, i);
                uint256 childValue = value / 100 * (s3 % 80);
                value -= childValue;

                // CREATE / CREATE2 (chance: 50%)
                address childAddr;
                if(s3 % 100 < 50) {
                    childAddr = create(childValue, code);
                } else {
                    childAddr = create2(childValue, i + 1, code);
                }

                if(childAddr == address(0)) {
                    break;
                }

                // selfdestruct on creation (chance: 50%)
                s3 = random2(_seed+2, depth, i);
                if(s3 % 100 < 40) {
                    try DeployTestChild(childAddr).destroy(s3) {
                    } catch (bytes memory _err) {
                        emit CatchedRevert(_seed, 1, _err);
                    }
                } else {
                    // notify main contract about this instance, so we can selfdestruct / call it later on
                    try DeployTest(_main).notify(childAddr, false) {
                    } catch (bytes memory _err) {
                        emit CatchedRevert(_seed, 2, _err);
                    }
                }
            }
        }
    }

    function random2(uint seed, uint id1, uint id2) private pure returns (uint) {
        return uint(keccak256(abi.encodePacked(seed, id1, id2)));
    }

    function create(uint256 amount, bytes memory bytecode) private returns (address) {
        address addr;
        assembly {
            addr := create(amount, add(bytecode, 0x20), mload(bytecode))

            if iszero(extcodesize(addr)) {
                revert(0, 0)
            }
        }

        if(addr == address(0)) {
            emit CatchedRevert(_seed, 3, abi.encode());
        }
        return addr;
    }

    function create2(uint256 amount, uint salt, bytes memory bytecode) private returns (address) {
        address addr;
        assembly {
            addr := create2(amount, add(bytecode, 0x20), mload(bytecode), salt)

            if iszero(extcodesize(addr)) {
                revert(0, 0)
            }
        }
        if(addr == address(0)) {
            emit CatchedRevert(_seed, 3, abi.encode());
        }
        return addr;
    }

    function destroy(uint seed) public {
        address target;

        uint action = seed % 6;
        // selfdestruct beneficiary (1/6 chance each)
        if(action == 0) {
            // target: sender address
            target = msg.sender;
        } else if(action == 1) {
            // target: random new address
            target = address(bytes20(keccak256(abi.encodePacked(seed))));
        } else if(action == 2) {
            // target: origin address
            target = tx.origin;
        } else if(action == 3) {
            // target: self
            target = address(this);
        } else if(action == 4) {
            // target: static address
            target = address(0x49e0fd3800C117357057534E30c5B5115C673488);
        } else if(action == 5) {
            // target: main contract
            target = _main;
        }

        emit ChildDestroyed(seed, target);

        selfdestruct(payable(target));
    }

}
