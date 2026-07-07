// SPDX-License-Identifier: MIT
pragma solidity ^0.8.22;

contract EIP2780Helper {
    function nop() external payable {}

    function forwardValue(address payable target) external payable {
        (bool success,) = target.call{value: msg.value}("");
        require(success, "forward failed");
    }

    function createAndDestroy() external payable {
        SelfDestructor child = new SelfDestructor{value: msg.value}();
        child.destroy(payable(msg.sender));
    }
}

contract SelfDestructor {
    constructor() payable {}

    function destroy(address payable beneficiary) external {
        selfdestruct(beneficiary);
    }
}
