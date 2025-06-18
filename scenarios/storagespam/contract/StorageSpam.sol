// SPDX-License-Identifier: MIT
pragma solidity  ^0.8.22;

contract StorageSpam {
    mapping(uint256 => uint256) public storageMap;

    event RandomForGas(uint256 indexed gas, uint256 loops);

    constructor() {
        
    }

    function setStorage(uint256 key, uint256 value) public {
        storageMap[key] = value;
    }

    function getStorage(uint256 key) public view returns (uint256) {
        return storageMap[key];
    }

    function setRandomForGas(uint256 gasLimit, uint256 txid) public {
        uint256 gaslimit = gasleft();
        if (gasLimit > 0 && gaslimit > gasLimit) {
            gaslimit = gaslimit - gasLimit;
        } else {
            gaslimit = 100;
        }

        uint256 loops = 0;
        uint256 randomKey = uint256(keccak256(abi.encodePacked(block.number, txid)));
        while (gasleft() > gaslimit) {
            if (loops % 200 == 0) {
                randomKey = uint256(keccak256(abi.encodePacked(randomKey, loops)));
            }
            storageMap[randomKey] = randomKey << 128 | randomKey >> 128;
            randomKey = randomKey << 1 | randomKey >> 255;
            loops++;
        }
        emit RandomForGas(gasLimit, loops);
    }

}

