#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../../../scripts/compile-contract.sh"
cd "$SCRIPT_DIR"

# Safe (formerly Gnosis Safe) v1.4.1 canonical contracts are pulled as prebuilt
# hardhat artifacts from the @safe-global/safe-contracts package and turned into
# Go bindings via abigen. They are multi-file and can't be built with the
# single-file solc helper, so we consume the published creation bytecode directly
# (same approach the erc4337 scenario uses for the EntryPoint).
SAFE_VERSION="1.4.1"
SAFE_BASE_URL="https://unpkg.com/@safe-global/safe-contracts@${SAFE_VERSION}/build/artifacts/contracts"

# gen_from_artifact is provided by the shared compile-contract.sh sourced above.

# Safe singleton (master copy) and the proxy factory that creates the per-multisig
# proxies via createProxyWithNonce.
gen_from_artifact "${SAFE_BASE_URL}/Safe.sol/Safe.json" Safe
gen_from_artifact "${SAFE_BASE_URL}/proxies/SafeProxyFactory.sol/SafeProxyFactory.json" SafeProxyFactory

# Local single-file dummy gas-burner call target.
compile_contract "$(pwd)" 0.8.24 "--optimize --optimize-runs 200" GasBurner GasBurner
