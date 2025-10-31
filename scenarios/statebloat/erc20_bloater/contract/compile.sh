#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../../../../scripts/compile-contract.sh"
cd $SCRIPT_DIR

# ERC20Bloater
compile_contract "$(pwd)" 0.8.22 "--optimize --optimize-runs 200" ERC20Bloater
