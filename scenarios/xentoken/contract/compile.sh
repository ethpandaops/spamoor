#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../../../scripts/compile-contract.sh"
cd $SCRIPT_DIR

# XENMath
compile_contract "$(pwd)" 0.8.17 "--optimize --optimize-runs 20" XENMath

# XENCrypto
compile_contract "$(pwd)" 0.8.17 "--optimize --optimize-runs 20" XENCrypto

# XENSybilAttacker
compile_contract "$(pwd)" 0.8.22 "--optimize --optimize-runs 20" SybilAttacker XENSybilAttacker
