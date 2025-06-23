#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../../../scripts/compile-contract.sh"
cd $SCRIPT_DIR

# CREATE2Factory
compile_contract "$(pwd)" 0.8.0 "--optimize --optimize-runs 999999" CREATE2Factory
