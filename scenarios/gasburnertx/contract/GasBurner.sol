// SPDX-License-Identifier: MIT
pragma solidity  ^0.8.22;

contract GasBurner {
    address public worker;

    event BurnedGas(uint256 indexed gas, uint256 loops);

    constructor(bytes memory workerCode) {
        worker = create(workerCode);
    }

    function create(bytes memory bytecode) private returns (address) {
        address addr;
        assembly {
            addr := create(0, add(bytecode, 0x20), mload(bytecode))

            if iszero(extcodesize(addr)) {
                revert(0, 0)
            }
        }
        require(addr != address(0), "create failed");
        return addr;
    }

    function wasteEther(uint256 amount) internal {
        (bool success, bytes memory result) = worker.call{gas: amount - 18800}("");
        require(success, "worker call failed");
        emit BurnedGas(amount, uint256(bytes32(result)));
    }

    function burn2000k() public {
        wasteEther(2000000);
    }

    function burn1500k() public {
        wasteEther(1500000);
    }

    function burn1000k() public {
        wasteEther(1000000);
    }

    function burn500k() public {
        wasteEther(500000);
    }

    function burn100k() public {
        wasteEther(100000);
    }

    function burnGasUnits(uint256 amount) public {
        wasteEther(amount);
    }
}

