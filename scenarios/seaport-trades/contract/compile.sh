#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../../../scripts/compile-contract.sh"
cd "$SCRIPT_DIR"

# Seaport 1.6 (OpenSea's marketplace settlement contract) and its
# ConduitController are pulled as prebuilt Hardhat artifacts from the
# @opensea/seaport-js package and turned into Go bindings via abigen. They are
# multi-file + solc-viaIR builds that can't be produced with the single-file
# solc helper, so we consume the published creation bytecode directly (same
# approach the safe-multisig and erc4337 scenarios use). seaport-js 4.1.3 ships
# Seaport 1.6 (solc 0.8.24, viaIR).
SEAPORT_JS_VERSION="4.1.3"
SEAPORT_BASE="https://unpkg.com/@opensea/seaport-js@${SEAPORT_JS_VERSION}/src/artifacts/seaport/contracts"

# gen_from_artifact is provided by the shared compile-contract.sh sourced above.

# ConduitController must be deployed first: Seaport's constructor reads the
# conduit code hashes from it. LocalConduitController is the canonical
# ConduitController logic with a local-deploy wrapper (no-arg constructor).
gen_from_artifact "${SEAPORT_BASE}/conduit/ConduitController.sol/LocalConduitController.json" ConduitController

# Seaport's constructor takes the ConduitController address; abigen appends it at
# deploy time (the creation bytecode has no pre-appended constructor arg).
gen_from_artifact "${SEAPORT_BASE}/Seaport.sol/Seaport.json" Seaport

# Local single-file mock collection (NFT) and mock stablecoin (ERC20) traded by
# the scenario. Both have permissionless mint so the scenario self-seeds.
compile_contract "$(pwd)" 0.8.24 "--optimize --optimize-runs 200" MintableNFT MintableNFT
compile_contract "$(pwd)" 0.8.24 "--optimize --optimize-runs 200" MintableToken MintableToken
