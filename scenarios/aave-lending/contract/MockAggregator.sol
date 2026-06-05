// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

// Minimal Chainlink-style price source for the spamoor aave-lending scenario.
// AaveOracle only calls latestAnswer() on its sources (it falls back to the
// fallback oracle when the answer is <= 0), so a fixed positive answer is all
// that is required to give a reserve a stable price. Prices use 8 decimals to
// match the oracle's USD base currency unit (1e8 == $1.00).
contract MockAggregator {
    int256 public answer;

    constructor(int256 _answer) {
        answer = _answer;
    }

    function latestAnswer() external view returns (int256) {
        return answer;
    }

    function setAnswer(int256 _answer) external {
        answer = _answer;
    }
}
