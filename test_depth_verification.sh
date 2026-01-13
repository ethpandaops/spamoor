#!/bin/bash

echo "=== Storage Trie Branch Depth Verification Test ==="
echo ""
echo "This test demonstrates how to verify that branches are exactly 3 nibbles deep"
echo ""

# Test addresses that would create different depths
echo "Example 1: Testing a random address"
echo "----------------------------------------"
go run verify_branch_depth.go 0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb7

echo ""
echo "Example 2: Testing another address"
echo "----------------------------------------"
go run verify_branch_depth.go 0x5aAeb6053f3E94C9b9A09f33669435E7Ef1BeAed

echo ""
echo "How to verify depth 3:"
echo "----------------------"
echo "1. Look at the computed storage slot"
echo "2. Check if exactly the first 3 nibbles match the target prefix"
echo "3. The 4th nibble must be different (divergence point)"
echo ""
echo "For example:"
echo "  Target prefix: 0xABC..."
echo "  Storage slot:  0xABC7... → Depth 3 ✓ (first 3 match, 4th differs)"
echo "  Storage slot:  0xAB57... → Depth 2 ✗ (only 2 match)"
echo "  Storage slot:  0xABCD... → Depth 4 ✗ (too deep)"