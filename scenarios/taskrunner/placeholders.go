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
	contractMatches := contractRegex.FindAllStringSubmatch(str, -1)

	processed := str
	for _, match := range contractMatches {
		if len(match) < 2 {
			continue
		}
		contractName := match[1]
		contractAddr, exists := registry.Get(contractName)
		if !exists {
			return "", fmt.Errorf("contract '%s' not found in registry", contractName)
		}

		var addressStr string
		// Check if child address calculation is requested
		if len(match) > 2 && match[2] != "" {
			// Parse nonce for child address calculation
			nonce, err := strconv.ParseUint(match[2], 10, 64)
			if err != nil {
				return "", fmt.Errorf("invalid nonce value for contract '%s': %s", contractName, match[2])
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

		processed = strings.Replace(processed, match[0], addressStr, 1)
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
	for randomRegex.MatchString(processed) {
		randomVal, err := generateRandomUint256()
		if err != nil {
			return "", fmt.Errorf("failed to generate random value: %w", err)
		}
		processed = randomRegex.ReplaceAllString(processed, randomVal)
	}

	// Replace {random:N} with random number between 0 and N
	randomRangeRegex := regexp.MustCompile(`\{random:(\d+)\}`)
	matches := randomRangeRegex.FindAllStringSubmatch(processed, -1)
	for _, match := range matches {
		if len(match) != 2 {
			continue
		}
		maxVal, err := strconv.ParseUint(match[1], 10, 64)
		if err != nil {
			return "", fmt.Errorf("invalid random range: %s", match[1])
		}
		randomVal, err := generateRandomInRange(maxVal)
		if err != nil {
			return "", fmt.Errorf("failed to generate random value in range: %w", err)
		}
		processed = strings.Replace(processed, match[0], fmt.Sprintf("%d", randomVal), 1)
	}

	// Replace {randomaddr} with random address
	randomAddrRegex := regexp.MustCompile(`\{randomaddr\}`)
	for randomAddrRegex.MatchString(processed) {
		randomAddr, err := generateRandomAddress()
		if err != nil {
			return "", fmt.Errorf("failed to generate random address: %w", err)
		}
		// Strip 0x prefix if requested (for bytecode/calldata usage)
		if stripPrefix && strings.HasPrefix(randomAddr, "0x") {
			randomAddr = randomAddr[2:]
		}
		processed = randomAddrRegex.ReplaceAllString(processed, randomAddr)
	}

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
