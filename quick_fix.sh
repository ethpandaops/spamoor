#!/bin/bash

# Fix remaining compilation errors

# Fix GetChainId() - it already returns *big.Int, no need to wrap
perl -pi -e 's/big\.NewInt\(wallet\.GetChainId\(\)\)/wallet.GetChainId()/g' scenarios/statebloat/storage_trie_brancher/*.go

# Fix string multiplication
perl -pi -e 's/"=" \* 60/strings.Repeat("=", 60)/g' scenarios/statebloat/storage_trie_brancher/verify_depth.go

# Remove unused os import
perl -pi -e 's/^\s*"os"\n//g' scenarios/statebloat/storage_trie_brancher/storage_trie_brancher.go

# Fix GetRPCHost to GetRPCEndpoint
perl -pi -e 's/GetRPCHost\(\)/GetRPCEndpoint()/g' scenarios/statebloat/storage_trie_brancher/storage_trie_brancher.go

# Add strings import to verify_depth.go
perl -pi -e 's/^import \(/import (\n\t"strings"/g if /^import \(/' scenarios/statebloat/storage_trie_brancher/verify_depth.go

echo "Fixes applied"