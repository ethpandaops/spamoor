// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

/// @title GasBurner
/// @notice Dummy call target for the safe-multisig scenario. burn() consumes a
///         deterministic, seed-derived pseudo-random amount of gas and emits an
///         event. The number of hashing rounds is derived from the seed (bounded
///         by maxRounds), so the gas cost is reproducible for a given
///         (seed, maxRounds) yet varies across different seeds - producing a
///         realistic mix of contract-call gas usage when driven by a stream of
///         random seeds.
contract GasBurner {
    event GasBurned(address indexed caller, uint256 seed, uint256 rounds, bytes32 digest);

    /// @notice Burn a deterministic, seed-derived amount of gas and emit GasBurned.
    /// @param seed Caller-provided seed; the same seed always burns the same gas.
    /// @param maxRounds Upper bound (exclusive offset) on the hashing rounds; the
    ///        actual round count is 1 + (keccak(seed) % maxRounds).
    /// @return rounds The number of hashing rounds performed.
    function burn(uint256 seed, uint256 maxRounds) external returns (uint256 rounds) {
        if (maxRounds == 0) {
            maxRounds = 1;
        }

        bytes32 digest = keccak256(abi.encodePacked(seed));
        rounds = 1 + (uint256(digest) % maxRounds);

        for (uint256 i = 0; i < rounds; i++) {
            digest = keccak256(abi.encodePacked(digest, i));
        }

        emit GasBurned(msg.sender, seed, rounds, digest);
    }
}
