// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

interface IBaseRegistrar {
    function register(uint256 id, address owner, uint256 duration) external returns (uint256);

    function renew(uint256 id, uint256 duration) external returns (uint256);

    function available(uint256 id) external view returns (bool);
}

interface IENS {
    function setResolver(bytes32 node, address resolver) external;

    function setOwner(bytes32 node, address owner) external;
}

interface IReverseRegistrar {
    function setNameForAddr(
        address addr,
        address owner,
        address resolver,
        string memory name
    ) external returns (bytes32);
}

interface IAddrResolver {
    function setAddr(bytes32 node, address a) external;
}

/// @title SpamRegistrarController
/// @notice Permissionless auxiliary .eth registrar controller for spamoor
/// testnets. It is granted controller rights ONCE on the base registrar and
/// the reverse registrar (by the EnsExecutor that owns them) and then lets
/// any wallet:
///
/// - register short-lived names directly (no commit-reveal, no 28-day
///   minimum duration) so name expiry is observable within a test run, and
/// - register a fully wired name (forward + reverse resolution) for an
///   arbitrary address in a single transaction - used by the wallet naming
///   service to label all spamoor wallets without per-instance grants.
///
/// Deliberately unauthenticated: on a testnet, free direct registrations are
/// exactly what this scenario exists to generate. Names claimed through the
/// commit-reveal ETHRegistrarController are unaffected (already-registered
/// names cannot be re-registered here).
contract SpamRegistrarController {
    // namehash("eth")
    bytes32 public constant ETH_NODE = 0x93cdeb708b7545dc668eb9280176169d1c33cfd8ed6f04690a0bcc88a93fc4ae;

    IBaseRegistrar public immutable base;
    IENS public immutable ens;
    IReverseRegistrar public immutable reverseRegistrar;
    address public immutable resolver;

    constructor(IBaseRegistrar _base, IENS _ens, IReverseRegistrar _reverseRegistrar, address _resolver) {
        base = _base;
        ens = _ens;
        reverseRegistrar = _reverseRegistrar;
        resolver = _resolver;
    }

    /// @notice Register `label`.eth for `owner` with an arbitrary duration
    /// (seconds granularity, no minimum). Updates the registry so the owner
    /// controls the node.
    function register(string calldata label, address owner, uint256 duration) external returns (uint256 expires) {
        return base.register(uint256(keccak256(bytes(label))), owner, duration);
    }

    /// @notice Renew `label`.eth for an arbitrary duration.
    function renew(string calldata label, uint256 duration) external returns (uint256 expires) {
        return base.renew(uint256(keccak256(bytes(label))), duration);
    }

    /// @notice Register `label`.eth pointing at `addr` with forward + reverse
    /// resolution in one transaction:
    /// forward: label.eth -> addr (public resolver addr record)
    /// reverse: addr.addr.reverse -> label.eth
    /// The registry node ownership is handed to `addr` afterwards.
    function registerNamed(string calldata label, address addr, uint256 duration) external returns (bytes32 node) {
        bytes32 labelHash = keccak256(bytes(label));
        base.register(uint256(labelHash), address(this), duration);

        node = keccak256(abi.encodePacked(ETH_NODE, labelHash));
        ens.setResolver(node, resolver);
        IAddrResolver(resolver).setAddr(node, addr);
        reverseRegistrar.setNameForAddr(addr, addr, resolver, string.concat(label, ".eth"));
        ens.setOwner(node, addr);
    }
}
