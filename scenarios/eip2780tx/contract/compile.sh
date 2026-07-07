#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../../../scripts/compile-contract.sh"
cd $SCRIPT_DIR

# EIP2780Helper
compile_contract "$(pwd)" 0.8.22 "--optimize --optimize-runs 200" EIP2780Helper
