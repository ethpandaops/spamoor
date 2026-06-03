#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../../../scripts/compile-contract.sh"
cd "$SCRIPT_DIR"

# ERC-4337 v0.7 reference contracts (EntryPoint, SimpleAccountFactory) are pulled
# as prebuilt artifacts from the canonical @account-abstraction/contracts package
# and turned into Go bindings via abigen. They are too large / multi-file to
# compile with the single-file solc helper, so we consume the published bytecode
# directly. The factory artifact embeds the SimpleAccount implementation bytecode.
AA_VERSION="0.7.0"
AA_BASE_URL="https://unpkg.com/@account-abstraction/contracts@${AA_VERSION}/artifacts"

gen_from_artifact() {
    local name=$1 # artifact basename, also the abigen type and output .go basename
    local json

    json=$(curl -fsSL "${AA_BASE_URL}/${name}.json")
    echo "$json" | jq -r '.abi' >"${name}.abi"
    echo "$json" | jq -r '.bytecode' | sed 's/^0x//' >"${name}.bin"

    docker run --rm -u "$(id -u):$(id -g)" -v "$(pwd)":/workspace ethereum/client-go:alltools-latest \
        abigen --bin=/workspace/"${name}.bin" --abi=/workspace/"${name}.abi" --pkg=contract --out=/workspace/"${name}.go" --type "${name}"

    echo "$json" | jq >"${name}.output.json"
    rm "${name}.abi" "${name}.bin"
}

# Only EntryPoint needs the PackedUserOperation tuple; SimpleAccount is omitted on
# purpose to avoid a duplicate struct definition in the shared contract package
# (its execute() calldata is ABI-encoded directly in Go instead).
gen_from_artifact EntryPoint
gen_from_artifact SimpleAccountFactory

# Local single-file contracts.
compile_contract "$(pwd)" 0.8.23 "--optimize --optimize-runs 200" Paymaster AcceptAllPaymaster
compile_contract "$(pwd)" 0.8.23 "--optimize --optimize-runs 200" Counter Counter

# The Paymaster binding re-declares the PackedUserOperation struct (its ABI
# references the tuple), which collides with the EntryPoint binding in this same
# package. The paymaster's struct-referencing methods are never called from Go,
# so strip the duplicate type declaration.
perl -0pi -e 's{// PackedUserOperation is an auto generated[^\n]*\ntype PackedUserOperation struct \{.*?\n\}\n}{// PackedUserOperation is defined in the EntryPoint binding within this same\n// package; the duplicate generated here is removed by compile.sh to avoid a\n// redeclaration.\n}s' Paymaster.go
