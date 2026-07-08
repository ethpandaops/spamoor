// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

/// @title EnsExecutor
/// @notice Deployment & administration proxy for the spamoor ENS stack.
///
/// The ENS contracts derive ownership from msg.sender at construction time
/// (ENSRegistry root node, Ownable on the registrar/controller/reverse
/// registrar). Deploying them through the shared CREATE2 deployment factory
/// would make the factory proxy the immutable owner and permanently lock all
/// owner-gated wiring. This executor is deployed through the factory instead
/// (with the spamoor root wallet as admin constructor arg) and then deploys
/// the ENS contracts itself via CREATE2, becoming their owner. Owner-gated
/// calls are performed via execute().
///
/// The admin (root wallet) authorizes per-scenario operator wallets, so
/// multiple concurrently running scenario instances can share the stack
/// without sharing a wallet.
contract EnsExecutor {
    address public admin;
    mapping(address => bool) public operators;

    event AdminChanged(address indexed newAdmin);
    event OperatorChanged(address indexed operator, bool enabled);
    event Deployed(address addr);

    modifier onlyAdmin() {
        require(msg.sender == admin, "EnsExecutor: not admin");
        _;
    }

    modifier onlyAuthorized() {
        require(msg.sender == admin || operators[msg.sender], "EnsExecutor: not authorized");
        _;
    }

    constructor(address _admin) {
        admin = _admin;
    }

    /// @notice Accept ETH (e.g. rent swept out of owned ENS controllers).
    receive() external payable {}

    function setAdmin(address _admin) external onlyAdmin {
        admin = _admin;
        emit AdminChanged(_admin);
    }

    function setOperator(address operator, bool enabled) external onlyAdmin {
        operators[operator] = enabled;
        emit OperatorChanged(operator, enabled);
    }

    /// @notice Deploy a contract via CREATE2 so this executor is its
    /// msg.sender (and thus owner for ownership-from-deployer contracts).
    /// Addresses are computed off-chain from (executor, salt, initcode hash).
    function deploy(bytes32 salt, bytes memory initcode) external onlyAuthorized returns (address addr) {
        assembly {
            addr := create2(0, add(initcode, 0x20), mload(initcode), salt)
        }
        require(addr != address(0), "EnsExecutor: create2 failed");
        emit Deployed(addr);
    }

    /// @notice Perform an owner-gated call on a contract owned by this executor.
    function execute(address target, uint256 value, bytes calldata data)
        external
        payable
        onlyAuthorized
        returns (bytes memory)
    {
        (bool success, bytes memory result) = target.call{value: value}(data);
        if (!success) {
            // bubble up the inner revert reason
            assembly {
                revert(add(result, 0x20), mload(result))
            }
        }
        return result;
    }
}
