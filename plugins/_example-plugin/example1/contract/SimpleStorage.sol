// SPDX-License-Identifier: MIT
pragma solidity ^0.8.22;

/**
 * @title SimpleStorage
 * @dev A simple contract for demonstration purposes
 * Stores and retrieves a value, emits events, and allows incrementing
 */
contract SimpleStorage {
    uint256 private value;
    address public owner;
    
    event ValueSet(address indexed setter, uint256 oldValue, uint256 newValue);
    event ValueIncremented(address indexed incrementer, uint256 newValue);
    
    constructor(uint256 _initialValue) {
        value = _initialValue;
        owner = msg.sender;
        emit ValueSet(msg.sender, 0, _initialValue);
    }
    
    function setValue(uint256 _value) public {
        uint256 oldValue = value;
        value = _value;
        emit ValueSet(msg.sender, oldValue, _value);
    }
    
    function getValue() public view returns (uint256) {
        return value;
    }
    
    function increment() public {
        value += 1;
        emit ValueIncremented(msg.sender, value);
    }
    
    function incrementBy(uint256 _amount) public {
        value += _amount;
        emit ValueIncremented(msg.sender, value);
    }
}