package main

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// NewWalletKeyCmd creates the wallet-key subcommand for deriving child wallet
// private keys from a root key, seed, and wallet index or well-known name.
func NewWalletKeyCmd(logger logrus.FieldLogger) *cobra.Command {
	var rootKey string
	var seed string
	var walletIndex int64
	var walletName string
	var veryWellKnown bool
	var count uint64

	cmd := &cobra.Command{
		Use:   "wallet-key",
		Short: "Derive a child wallet private key",
		Long: `Derive the private key and address of a spamoor child wallet.

The derivation uses the same logic as spamoor's WalletPool:
  SHA256(root_private_key || identifier || seed)

For indexed wallets, the identifier is the 8-byte big-endian wallet index.
For well-known wallets, the identifier is the wallet name as bytes.

The seed is appended to the identifier unless --very-well-known is set.

Use --count to generate multiple consecutive indexed wallets starting from --index.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWalletKey(rootKey, seed, walletIndex, walletName, veryWellKnown, count)
		},
	}

	cmd.Flags().StringVar(&rootKey, "root-key", "", "Root wallet private key (hex, with or without 0x prefix)")
	cmd.Flags().StringVar(&seed, "seed", "", "Wallet seed (scenario seed)")
	cmd.Flags().Int64Var(&walletIndex, "index", -1, "Child wallet index (0-based)")
	cmd.Flags().StringVar(&walletName, "name", "", "Well-known wallet name")
	cmd.Flags().BoolVar(&veryWellKnown, "very-well-known", false, "Skip seed for well-known wallet derivation")
	cmd.Flags().Uint64Var(&count, "count", 0, "Number of consecutive indexed wallets to derive (starting from --index)")

	if err := cmd.MarkFlagRequired("root-key"); err != nil {
		logger.WithError(err).Error("failed to mark root-key flag as required")
	}

	return cmd
}

// deriveChildKey derives a child private key and address from a parent key and identifier bytes.
func deriveChildKey(parentKeyBytes, idxBytes []byte) (string, common.Address, error) {
	childKeyHash := sha256.Sum256(append(parentKeyBytes, idxBytes...))
	childKeyHex := fmt.Sprintf("%x", childKeyHash)

	childPrivKey, err := crypto.HexToECDSA(childKeyHex)
	if err != nil {
		return "", common.Address{}, fmt.Errorf("failed to derive child key: %w", err)
	}

	publicKey := childPrivKey.Public()

	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", common.Address{}, fmt.Errorf("failed to cast public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	return childKeyHex, address, nil
}

func runWalletKey(rootKey, seed string, walletIndex int64, walletName string, veryWellKnown bool, count uint64) error {
	if walletIndex < 0 && walletName == "" {
		return fmt.Errorf("either --index or --name must be specified")
	}

	if walletIndex >= 0 && walletName != "" {
		return fmt.Errorf("--index and --name are mutually exclusive")
	}

	if count > 0 && walletName != "" {
		return fmt.Errorf("--count can only be used with --index")
	}

	// Parse root private key.
	rootKey = strings.TrimPrefix(rootKey, "0x")

	rootPrivKey, err := crypto.HexToECDSA(rootKey)
	if err != nil {
		return fmt.Errorf("invalid root key: %w", err)
	}

	parentKeyBytes := crypto.FromECDSA(rootPrivKey)

	// Batch mode: generate multiple consecutive indexed wallets.
	if count > 0 {
		seedBytes := []byte(seed)

		for i := range count {
			idx := uint64(walletIndex) + i

			idxBytes := make([]byte, 8)
			binary.BigEndian.PutUint64(idxBytes, idx)

			if seed != "" {
				idxBytes = append(idxBytes, seedBytes...)
			}

			keyHex, address, err := deriveChildKey(parentKeyBytes, idxBytes)
			if err != nil {
				return fmt.Errorf("failed to derive wallet %d: %w", idx, err)
			}

			fmt.Printf("%s 0x%s\n", address.Hex(), keyHex)
		}

		return nil
	}

	// Single wallet mode.
	var idxBytes []byte

	if walletName != "" {
		// Well-known wallet: identifier is the name bytes.
		idxBytes = make([]byte, len(walletName))
		copy(idxBytes, walletName)

		if seed != "" && !veryWellKnown {
			idxBytes = append(idxBytes, []byte(seed)...)
		}
	} else {
		// Indexed wallet: identifier is 8-byte big-endian index.
		idxBytes = make([]byte, 8)
		binary.BigEndian.PutUint64(idxBytes, uint64(walletIndex))

		if seed != "" {
			idxBytes = append(idxBytes, []byte(seed)...)
		}
	}

	keyHex, address, err := deriveChildKey(parentKeyBytes, idxBytes)
	if err != nil {
		return err
	}

	fmt.Printf("Private Key: 0x%s\n", keyHex)
	fmt.Printf("Address:     %s\n", address.Hex())

	return nil
}
