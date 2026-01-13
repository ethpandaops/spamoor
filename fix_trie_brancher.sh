#!/bin/bash

# Fix the compilation issues in storage_trie_brancher

# Fix wallet selection constant
find scenarios/statebloat/storage_trie_brancher -name "*.go" -exec perl -pi -e 's/WalletSelectionByIndex/SelectWalletByIndex/g' {} \;

# Fix GetChainID to GetChainId
find scenarios/statebloat/storage_trie_brancher -name "*.go" -exec perl -pi -e 's/GetChainID\(\)/GetChainId()/g' {} \;

# Remove unused json import
perl -pi -e 's/"encoding\/json"\n//g' scenarios/statebloat/storage_trie_brancher/storage_trie_brancher.go

# Fix GetRPCClient
perl -pi -e 's/s\.walletPool\.GetRPCClient\(\)/s.walletPool.GetRPCHost()/g' scenarios/statebloat/storage_trie_brancher/storage_trie_brancher.go

echo "Fixes applied"