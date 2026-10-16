package taskrunner

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethpandaops/spamoor/utils"
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

// ProcessBasicPlaceholders processes the shared transaction placeholders (see utils.TxPlaceholders) and {stepid}.
// If stripPrefix is true, removes 0x prefix from addresses (for use in bytecode/calldata)
func ProcessBasicPlaceholders(str string, txIdx uint64, stepIdx int, stripPrefix bool) (string, error) {
	placeholders := utils.TxPlaceholders(txIdx)
	placeholders["stepid"] = func(string) (string, error) {
		return strconv.Itoa(stepIdx), nil
	}
	if stripPrefix {
		for _, name := range []string{"randomaddr", "create2"} {
			resolve := placeholders[name]
			placeholders[name] = func(arg string) (string, error) {
				addr, err := resolve(arg)
				return strings.TrimPrefix(addr, "0x"), err
			}
		}
	}
	return placeholders.Resolve(str)
}
