package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"os"
)

// This utility verifies the depth of storage trie branches
// Usage: go run verify_branch_depth.go <address> [spender]

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verify_branch_depth.go <address> [spender]")
		fmt.Println("Example (balance): go run verify_branch_depth.go 0x123...")
		fmt.Println("Example (allowance): go run verify_branch_depth.go 0x123... 0x456...")
		os.Exit(1)
	}

	address := common.HexToAddress(os.Args[1])

	var storageSlot common.Hash
	if len(os.Args) == 2 {
		// Balance mapping: keccak256(address || uint256(0))
		storageSlot = crypto.Keccak256Hash(
			common.LeftPadBytes(address.Bytes(), 32),
			common.LeftPadBytes([]byte{0}, 32),
		)
		fmt.Println("Type: Balance mapping")
		fmt.Printf("Address: %s\n", address.Hex())
	} else {
		// Allowance mapping: keccak256(spender || keccak256(owner || uint256(1)))
		spender := common.HexToAddress(os.Args[2])
		innerHash := crypto.Keccak256Hash(
			common.LeftPadBytes(address.Bytes(), 32),
			common.LeftPadBytes([]byte{1}, 32),
		)
		storageSlot = crypto.Keccak256Hash(
			common.LeftPadBytes(spender.Bytes(), 32),
			innerHash.Bytes(),
		)
		fmt.Println("Type: Allowance mapping")
		fmt.Printf("Owner: %s\n", address.Hex())
		fmt.Printf("Spender: %s\n", spender.Hex())
	}

	fmt.Printf("\nComputed Storage Slot: %s\n", storageSlot.Hex())
	fmt.Println("\nNibble breakdown:")

	slotHex := storageSlot.Hex()[2:]              // Remove "0x"
	for i := 0; i < len(slotHex) && i < 16; i++ { // Show first 8 bytes (16 nibbles)
		if i%2 == 0 && i > 0 {
			fmt.Print(" ")
		}
		fmt.Printf("%c", slotHex[i])
	}
	fmt.Println("...")

	// Show how to verify depth against a target
	fmt.Println("\nTo verify depth against a target prefix:")
	fmt.Println("- Count matching nibbles from the start")
	fmt.Println("- Depth 3 = first 3 nibbles match")
	fmt.Printf("- Example: If target is 0xABC..., this slot has depth 3 if it starts with 0xABC\n")

	// Example depth check
	fmt.Println("\nDepth Analysis:")
	for depth := uint8(1); depth <= 8; depth++ {
		prefix := slotHex[:depth]
		fmt.Printf("  Depth %d: Prefix is %s\n", depth, prefix)
	}
}
