#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../../../scripts/compile-contract.sh"
cd "$SCRIPT_DIR"

# The canonical ENS stack is pulled as prebuilt Hardhat artifacts from the
# @ensdomains/ens-contracts npm package and turned into Go bindings via abigen
# (same approach as the seaport-trades / safe-multisig / erc4337 scenarios).
#
# v1.7.0 is the current release lineage (staging/main); the GitHub master
# branch is years stale and ships an incompatible NameWrapper-era
# ETHRegistrarController. Compiled with solc 0.8.26 (NameWrapper 0.8.17),
# evmVersion paris.
ENS_VERSION="1.7.0"
ENS_BASE="https://unpkg.com/@ensdomains/ens-contracts@${ENS_VERSION}/artifacts/contracts"

gen_from_artifact "${ENS_BASE}/registry/ENSRegistry.sol/ENSRegistry.json" ENSRegistry
gen_from_artifact "${ENS_BASE}/ethregistrar/BaseRegistrarImplementation.sol/BaseRegistrarImplementation.json" BaseRegistrarImplementation
gen_from_artifact "${ENS_BASE}/ethregistrar/ETHRegistrarController.sol/ETHRegistrarController.json" ETHRegistrarController
gen_from_artifact "${ENS_BASE}/ethregistrar/StablePriceOracle.sol/StablePriceOracle.json" StablePriceOracle
gen_from_artifact "${ENS_BASE}/ethregistrar/DummyOracle.sol/DummyOracle.json" DummyOracle
gen_from_artifact "${ENS_BASE}/resolvers/PublicResolver.sol/PublicResolver.json" PublicResolver
gen_from_artifact "${ENS_BASE}/reverseRegistrar/ReverseRegistrar.sol/ReverseRegistrar.json" ReverseRegistrar
gen_from_artifact "${ENS_BASE}/reverseRegistrar/DefaultReverseRegistrar.sol/DefaultReverseRegistrar.json" DefaultReverseRegistrar
gen_from_artifact "${ENS_BASE}/wrapper/NameWrapper.sol/NameWrapper.json" NameWrapper
gen_from_artifact "${ENS_BASE}/wrapper/StaticMetadataService.sol/StaticMetadataService.json" StaticMetadataService

# Both the ETHRegistrarController and StablePriceOracle bindings emit the
# IPriceOracle.Price tuple struct; drop the duplicate from the oracle binding
# to avoid a redeclaration (dockerized so the script only needs docker/curl/jq).
docker run --rm -u "$(id -u):$(id -g)" -v "$(pwd)":/workspace -w /workspace python:3-alpine python3 -c '
import re

path = "StablePriceOracle.go"
with open(path) as f:
    src = f.read()

pattern = re.compile(
    r"// IPriceOraclePrice is an auto generated[^\n]*\n"
    r"type IPriceOraclePrice struct \{.*?\n\}\n",
    re.S,
)
replacement = (
    "// IPriceOraclePrice is defined in the ETHRegistrarController binding within\n"
    "// this same package; the duplicate generated here is removed by compile.sh to\n"
    "// avoid a redeclaration.\n"
)

with open(path, "w") as f:
    f.write(pattern.sub(replacement, src, count=1))
'

# Local helper contracts:
# - EnsExecutor: CREATE2 deployment/admin proxy that owns the ENS stack
#   (the shared deployment factory cannot be the owner, see EnsExecutor.sol).
# - SpamRegistrarController: permissionless direct registrar controller for
#   short-lived churn names and one-tx wallet naming.
compile_contract "$(pwd)" 0.8.24 "--optimize --optimize-runs 200" EnsExecutor EnsExecutor
compile_contract "$(pwd)" 0.8.24 "--optimize --optimize-runs 200" SpamRegistrarController SpamRegistrarController
