// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract CREATE2Factory {
    event ContractDeployed(address indexed deployedAddress, bytes32 indexed salt);
    
    function deploy(bytes32 salt, bytes calldata initCode) external payable returns (address) {
        address deployedAddress;
        bytes memory code = initCode;
        
        // create2(value, code, salt)
        assembly {
            deployedAddress := create2(callvalue(), add(code, 0x20), mload(code), salt)
        }
        
        require(deployedAddress != address(0), "Deployment failed");
        emit ContractDeployed(deployedAddress, salt);
        return deployedAddress;
    }
    
    function predictAddress(bytes32 salt, bytes32 initCodeHash) external view returns (address) {
        return address(uint160(uint256(keccak256(abi.encodePacked(
            bytes1(0xff),
            address(this),
            salt,
            initCodeHash
        )))));
    }
}