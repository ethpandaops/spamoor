package taskrunner

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// ProcessContractPlaceholders processes {contract:name} and {contract:name:nonce} placeholders
// If stripPrefix is true, removes the 0x prefix from addresses (for use in bytecode/calldata)
func ProcessContractPlaceholders(str string, registry *ContractRegistry, stripPrefix bool) (string, error) {
	contractRegex := regexp.MustCompile(`\{contract:([^}:]+)(?::(\d+))?\}`)

	// Use ReplaceAllStringFunc to replace each match as we find it
	processed := contractRegex.ReplaceAllStringFunc(str, func(match string) string {
		// Re-parse the match to extract parts
		submatch := contractRegex.FindStringSubmatch(match)
		if len(submatch) < 2 {
			return match // Return original if parsing fails
		}

		contractName := submatch[1]
		contractAddr, exists := registry.Get(contractName)
		if !exists {
			// Return original match - this will cause an error to be returned later
			return match
		}

		var addressStr string
		// Check if child address calculation is requested
		if len(submatch) > 2 && submatch[2] != "" {
			// Parse nonce for child address calculation
			nonce, err := strconv.ParseUint(submatch[2], 10, 64)
			if err != nil {
				// Return original match - this will cause an error to be returned later
				return match
			}
			// Calculate child contract address
			childAddr := crypto.CreateAddress(contractAddr, nonce)
			addressStr = childAddr.Hex()
		} else {
			// Use original contract address
			addressStr = contractAddr.Hex()
		}

		// Strip 0x prefix if requested (for bytecode/calldata usage)
		if stripPrefix && strings.HasPrefix(addressStr, "0x") {
			addressStr = addressStr[2:]
		}

		return addressStr
	})

	// Check if any placeholders were left unreplaced (indicating errors)
	if contractRegex.MatchString(processed) {
		// Find the first unresolved placeholder to report error
		matches := contractRegex.FindAllStringSubmatch(processed, -1)
		for _, match := range matches {
			if len(match) >= 2 {
				contractName := match[1]
				if _, exists := registry.Get(contractName); !exists {
					return "", fmt.Errorf("contract '%s' not found in registry", contractName)
				}
				// If contract exists but placeholder is still there, it's a nonce parsing error
				if len(match) > 2 && match[2] != "" {
					if _, err := strconv.ParseUint(match[2], 10, 64); err != nil {
						return "", fmt.Errorf("invalid nonce value for contract '%s': %s", contractName, match[2])
					}
				}
			}
		}
		return "", fmt.Errorf("failed to process some contract placeholders")
	}

	return processed, nil
}

// ProcessBasicPlaceholders processes {txid}, {stepid}, {random}, {random:N}, {randomaddr} placeholders
// If stripPrefix is true, removes 0x prefix from addresses (for use in bytecode/calldata)
func ProcessBasicPlaceholders(str string, txIdx uint64, stepIdx int, stripPrefix bool) (string, error) {
	processed := str

	// Replace {txid} with transaction index
	processed = strings.ReplaceAll(processed, "{txid}", fmt.Sprintf("%d", txIdx))

	// Replace {stepid} with step index
	processed = strings.ReplaceAll(processed, "{stepid}", fmt.Sprintf("%d", stepIdx))

	// Replace {random} with random uint256
	randomRegex := regexp.MustCompile(`\{random\}`)
	processed = randomRegex.ReplaceAllStringFunc(processed, func(match string) string {
		randomVal, err := generateRandomUint256()
		if err != nil {
			return match // Will cause error to be returned later
		}
		return randomVal
	})

	// Replace {random:N} with random number between 0 and N
	randomRangeRegex := regexp.MustCompile(`\{random:(\d+)\}`)
	processed = randomRangeRegex.ReplaceAllStringFunc(processed, func(match string) string {
		submatch := randomRangeRegex.FindStringSubmatch(match)
		if len(submatch) != 2 {
			return match
		}
		maxVal, err := strconv.ParseUint(submatch[1], 10, 64)
		if err != nil {
			return match // Will cause error to be returned later
		}
		randomVal, err := generateRandomInRange(maxVal)
		if err != nil {
			return match // Will cause error to be returned later
		}
		return fmt.Sprintf("%d", randomVal)
	})

	// Check for any unresolved {random:N} placeholders (indicating errors)
	if randomRangeRegex.MatchString(processed) {
		matches := randomRangeRegex.FindAllStringSubmatch(processed, -1)
		for _, match := range matches {
			if len(match) == 2 {
				if _, err := strconv.ParseUint(match[1], 10, 64); err != nil {
					return "", fmt.Errorf("invalid random range: %s", match[1])
				}
			}
		}
		return "", fmt.Errorf("failed to process some random range placeholders")
	}

	// Replace {randomaddr} with random address
	randomAddrRegex := regexp.MustCompile(`\{randomaddr\}`)
	processed = randomAddrRegex.ReplaceAllStringFunc(processed, func(match string) string {
		randomAddr, err := generateRandomAddress()
		if err != nil {
			return match // Will cause error to be returned later
		}
		// Strip 0x prefix if requested (for bytecode/calldata usage)
		if stripPrefix && strings.HasPrefix(randomAddr, "0x") {
			randomAddr = randomAddr[2:]
		}
		return randomAddr
	})

	return processed, nil
}

// Helper functions for placeholder generation
func generateRandomUint256() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return new(big.Int).SetBytes(bytes).String(), nil
}

func generateRandomInRange(max uint64) (uint64, error) {
	if max == 0 {
		return 0, nil
	}
	bytes := make([]byte, 8)
	_, err := rand.Read(bytes)
	if err != nil {
		return 0, err
	}
	randomVal := new(big.Int).SetBytes(bytes).Uint64()
	return randomVal % max, nil
}

func generateRandomAddress() (string, error) {
	bytes := make([]byte, 20)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return common.BytesToAddress(bytes).Hex(), nil
}
