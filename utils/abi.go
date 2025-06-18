package utils

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type ABICallDataBuilder struct {
	contractABI    *abi.ABI
	functionName   string
	functionSig    string
	callArgs       string
	parsedFunction *abi.Method
}

func NewABICallDataBuilder(callABI, callFnName, callFnSig, callArgs string) (*ABICallDataBuilder, error) {
	builder := &ABICallDataBuilder{
		functionName: callFnName,
		functionSig:  callFnSig,
		callArgs:     callArgs,
	}

	// Parse ABI if provided
	if callABI != "" {
		contractABI, err := abi.JSON(strings.NewReader(callABI))
		if err != nil {
			return nil, fmt.Errorf("failed to parse contract ABI: %w", err)
		}
		builder.contractABI = &contractABI
	}

	// Parse function
	if err := builder.parseFunction(); err != nil {
		return nil, err
	}

	return builder, nil
}

func (b *ABICallDataBuilder) parseFunction() error {
	if b.contractABI != nil && b.functionName != "" {
		// Use ABI + function name
		method, exists := b.contractABI.Methods[b.functionName]
		if !exists {
			return fmt.Errorf("function %s not found in ABI", b.functionName)
		}
		b.parsedFunction = &method
	} else if b.functionSig != "" {
		// Parse function signature directly
		functionABI, err := parseFunctionSignature(b.functionSig)
		if err != nil {
			return fmt.Errorf("failed to parse function signature: %w", err)
		}

		parsedABI, err := abi.JSON(strings.NewReader(fmt.Sprintf("[%s]", functionABI)))
		if err != nil {
			return fmt.Errorf("failed to parse constructed ABI: %w", err)
		}

		// Get the function from the parsed ABI
		for _, method := range parsedABI.Methods {
			b.parsedFunction = &method
			break
		}
	} else {
		return fmt.Errorf("either call-abi+call-fn-name or call-fn-signature must be provided")
	}

	return nil
}

// parseFunctionSignature parses a function signature like "transfer(address,uint256)"
// and returns the corresponding ABI JSON string
func parseFunctionSignature(signature string) (string, error) {
	// Remove any whitespace
	signature = strings.ReplaceAll(signature, " ", "")

	// Find the opening parenthesis
	parenIndex := strings.Index(signature, "(")
	if parenIndex == -1 {
		return "", fmt.Errorf("invalid function signature: missing opening parenthesis")
	}

	functionName := signature[:parenIndex]
	if functionName == "" {
		return "", fmt.Errorf("invalid function signature: missing function name")
	}

	// Extract the parameters part
	if !strings.HasSuffix(signature, ")") {
		return "", fmt.Errorf("invalid function signature: missing closing parenthesis")
	}

	paramsPart := signature[parenIndex+1 : len(signature)-1]

	// Parse parameters
	var inputs []string
	if paramsPart != "" {
		paramTypes := strings.Split(paramsPart, ",")
		for i, paramType := range paramTypes {
			paramType = strings.TrimSpace(paramType)
			if paramType == "" {
				return "", fmt.Errorf("invalid function signature: empty parameter type at position %d", i)
			}

			// Validate the parameter type
			if err := validateSolidityType(paramType); err != nil {
				return "", fmt.Errorf("invalid parameter type '%s': %w", paramType, err)
			}

			inputs = append(inputs, fmt.Sprintf(`{"name":"param%d","type":"%s"}`, i, paramType))
		}
	}

	// Construct ABI JSON
	abiJSON := fmt.Sprintf(`{
		"type": "function",
		"name": "%s",
		"inputs": [%s],
		"outputs": [],
		"stateMutability": "payable"
	}`, functionName, strings.Join(inputs, ","))

	return abiJSON, nil
}

// validateSolidityType validates if a string is a valid Solidity type
func validateSolidityType(typeName string) error {
	// Basic types
	basicTypes := map[string]bool{
		"address": true, "bool": true, "string": true, "bytes": true,
	}

	if basicTypes[typeName] {
		return nil
	}

	// Check uint types (uint8, uint16, ..., uint256)
	if strings.HasPrefix(typeName, "uint") {
		suffix := typeName[4:]
		if suffix == "" {
			return nil // uint is valid (defaults to uint256)
		}
		size, err := strconv.Atoi(suffix)
		if err != nil || size < 8 || size > 256 || size%8 != 0 {
			return fmt.Errorf("invalid uint size: must be between 8 and 256 and divisible by 8")
		}
		return nil
	}

	// Check int types (int8, int16, ..., int256)
	if strings.HasPrefix(typeName, "int") {
		suffix := typeName[3:]
		if suffix == "" {
			return nil // int is valid (defaults to int256)
		}
		size, err := strconv.Atoi(suffix)
		if err != nil || size < 8 || size > 256 || size%8 != 0 {
			return fmt.Errorf("invalid int size: must be between 8 and 256 and divisible by 8")
		}
		return nil
	}

	// Check fixed bytes types (bytes1, bytes2, ..., bytes32)
	if strings.HasPrefix(typeName, "bytes") {
		suffix := typeName[5:]
		if suffix == "" {
			return nil // bytes is valid (dynamic bytes)
		}
		size, err := strconv.Atoi(suffix)
		if err != nil || size < 1 || size > 32 {
			return fmt.Errorf("invalid fixed bytes size: must be between 1 and 32")
		}
		return nil
	}

	// Check array types
	if strings.HasSuffix(typeName, "[]") {
		baseType := typeName[:len(typeName)-2]
		return validateSolidityType(baseType)
	}

	// Check fixed array types like uint256[10]
	if strings.Contains(typeName, "[") && strings.HasSuffix(typeName, "]") {
		bracketIndex := strings.LastIndex(typeName, "[")
		baseType := typeName[:bracketIndex]
		sizeStr := typeName[bracketIndex+1 : len(typeName)-1]

		if sizeStr == "" {
			return fmt.Errorf("invalid array type: missing size")
		}

		size, err := strconv.Atoi(sizeStr)
		if err != nil || size <= 0 {
			return fmt.Errorf("invalid array size: must be positive integer")
		}

		return validateSolidityType(baseType)
	}

	return fmt.Errorf("unknown type '%s'", typeName)
}

func (b *ABICallDataBuilder) BuildCallData(txIdx uint64) ([]byte, error) {
	if b.parsedFunction == nil {
		return nil, fmt.Errorf("function not parsed")
	}

	// Parse and substitute placeholders in args
	processedArgs, err := b.processArgs(b.callArgs, txIdx)
	if err != nil {
		return nil, fmt.Errorf("failed to process args: %w", err)
	}

	// Parse JSON args
	var argValues []interface{}
	if processedArgs != "" && processedArgs != "[]" {
		if err := json.Unmarshal([]byte(processedArgs), &argValues); err != nil {
			return nil, fmt.Errorf("failed to parse call args as JSON: %w", err)
		}
	}

	// Validate arguments against ABI
	if err := b.validateArgs(argValues); err != nil {
		return nil, fmt.Errorf("argument validation failed: %w", err)
	}

	// Convert args to proper types
	convertedArgs, err := b.convertArgs(argValues)
	if err != nil {
		return nil, fmt.Errorf("failed to convert args: %w", err)
	}

	// Pack the function call
	callData, err := b.parsedFunction.Inputs.Pack(convertedArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to pack function arguments: %w", err)
	}

	// Prepend function selector
	fullCallData := append(b.parsedFunction.ID, callData...)
	return fullCallData, nil
}

func (b *ABICallDataBuilder) processArgs(args string, txIdx uint64) (string, error) {
	if args == "" {
		return "[]", nil
	}

	// Replace placeholders
	processed := args

	// Replace {txid} with transaction index
	processed = strings.ReplaceAll(processed, "{txid}", fmt.Sprintf("%d", txIdx))

	// Replace {random} with random uint256
	randomRegex := regexp.MustCompile(`\{random\}`)
	randomMatches := randomRegex.FindAllStringSubmatch(processed, -1)
	for _, match := range randomMatches {
		randomVal, err := generateRandomUint256()
		if err != nil {
			return "", fmt.Errorf("failed to generate random value: %w", err)
		}
		processed = strings.Replace(processed, match[0], randomVal, 1)
	}

	// Replace {random:N} with random number between 0 and N
	randomRangeRegex := regexp.MustCompile(`\{random:(\d+)\}`)
	randomRangeMatches := randomRangeRegex.FindAllStringSubmatch(processed, -1)
	for _, match := range randomRangeMatches {
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
	randomAddrMatches := randomAddrRegex.FindAllStringSubmatch(processed, -1)
	for _, match := range randomAddrMatches {
		randomAddr, err := generateRandomAddress()
		if err != nil {
			return "", fmt.Errorf("failed to generate random address: %w", err)
		}
		processed = strings.Replace(processed, match[0], randomAddr, 1)
	}

	return processed, nil
}

func (b *ABICallDataBuilder) convertArgs(args []interface{}) ([]interface{}, error) {
	if len(args) != len(b.parsedFunction.Inputs) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(b.parsedFunction.Inputs), len(args))
	}

	converted := make([]interface{}, len(args))
	for i, arg := range args {
		expectedType := b.parsedFunction.Inputs[i].Type
		convertedArg, err := convertArgument(arg, expectedType)
		if err != nil {
			return nil, fmt.Errorf("failed to convert argument %d: %w", i, err)
		}
		converted[i] = convertedArg
	}

	return converted, nil
}

func convertArgument(arg interface{}, expectedType abi.Type) (interface{}, error) {
	switch expectedType.T {
	case abi.UintTy:
		return convertToUint(arg, expectedType.Size)
	case abi.IntTy:
		return convertToInt(arg, expectedType.Size)
	case abi.AddressTy:
		return convertToAddress(arg)
	case abi.BoolTy:
		return convertToBool(arg)
	case abi.StringTy:
		return convertToString(arg)
	case abi.BytesTy:
		return convertToBytes(arg)
	case abi.FixedBytesTy:
		return convertToFixedBytes(arg, expectedType.Size)
	case abi.SliceTy:
		return convertToSlice(arg, expectedType)
	default:
		return nil, fmt.Errorf("unsupported type: %s", expectedType.String())
	}
}

func convertToUint(arg interface{}, size int) (*big.Int, error) {
	switch v := arg.(type) {
	case string:
		if strings.HasPrefix(v, "0x") {
			result, ok := new(big.Int).SetString(v[2:], 16)
			if !ok {
				return nil, fmt.Errorf("invalid hex number: %s", v)
			}
			return result, nil
		}
		result, ok := new(big.Int).SetString(v, 10)
		if !ok {
			return nil, fmt.Errorf("invalid number: %s", v)
		}
		return result, nil
	case float64:
		return big.NewInt(int64(v)), nil
	case int:
		return big.NewInt(int64(v)), nil
	default:
		return nil, fmt.Errorf("cannot convert %T to uint%d", arg, size*8)
	}
}

func convertToInt(arg interface{}, size int) (*big.Int, error) {
	return convertToUint(arg, size) // Same conversion logic
}

func convertToAddress(arg interface{}) (common.Address, error) {
	switch v := arg.(type) {
	case string:
		if !common.IsHexAddress(v) {
			return common.Address{}, fmt.Errorf("invalid address format: %s", v)
		}
		return common.HexToAddress(v), nil
	default:
		return common.Address{}, fmt.Errorf("cannot convert %T to address", arg)
	}
}

func convertToBool(arg interface{}) (bool, error) {
	switch v := arg.(type) {
	case bool:
		return v, nil
	case string:
		return strconv.ParseBool(v)
	default:
		return false, fmt.Errorf("cannot convert %T to bool", arg)
	}
}

func convertToString(arg interface{}) (string, error) {
	switch v := arg.(type) {
	case string:
		return v, nil
	default:
		return fmt.Sprintf("%v", arg), nil
	}
}

func convertToBytes(arg interface{}) ([]byte, error) {
	switch v := arg.(type) {
	case string:
		if strings.HasPrefix(v, "0x") {
			return common.FromHex(v), nil
		}
		return []byte(v), nil
	default:
		return nil, fmt.Errorf("cannot convert %T to bytes", arg)
	}
}

func convertToFixedBytes(arg interface{}, size int) (interface{}, error) {
	bytes, err := convertToBytes(arg)
	if err != nil {
		return nil, err
	}

	if len(bytes) > size {
		return nil, fmt.Errorf("bytes too long for fixed bytes%d", size)
	}

	// Pad to the right size
	padded := make([]byte, size)
	copy(padded, bytes)

	// Return as fixed-size array
	switch size {
	case 32:
		var result [32]byte
		copy(result[:], padded)
		return result, nil
	default:
		return padded, nil
	}
}

func convertToSlice(arg interface{}, expectedType abi.Type) (interface{}, error) {
	slice, ok := arg.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected array/slice for %s", expectedType.String())
	}

	converted := make([]interface{}, len(slice))
	for i, elem := range slice {
		convertedElem, err := convertArgument(elem, *expectedType.Elem)
		if err != nil {
			return nil, fmt.Errorf("failed to convert slice element %d: %w", i, err)
		}
		converted[i] = convertedElem
	}

	return converted, nil
}

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

// validateArgs validates that the provided arguments are compatible with the function signature
func (b *ABICallDataBuilder) validateArgs(args []interface{}) error {
	expectedCount := len(b.parsedFunction.Inputs)
	actualCount := len(args)

	if actualCount != expectedCount {
		return fmt.Errorf("function %s expects %d arguments, got %d",
			b.parsedFunction.Name, expectedCount, actualCount)
	}

	// Validate each argument against its expected type
	for i, arg := range args {
		expectedType := b.parsedFunction.Inputs[i].Type
		paramName := b.parsedFunction.Inputs[i].Name
		if paramName == "" {
			paramName = fmt.Sprintf("param%d", i)
		}

		if err := b.validateArgType(arg, expectedType, paramName, i); err != nil {
			return err
		}
	}

	return nil
}

// validateArgType validates a single argument against its expected ABI type
func (b *ABICallDataBuilder) validateArgType(arg interface{}, expectedType abi.Type, paramName string, paramIndex int) error {
	switch expectedType.T {
	case abi.UintTy, abi.IntTy:
		return b.validateIntegerArg(arg, expectedType, paramName, paramIndex)
	case abi.AddressTy:
		return b.validateAddressArg(arg, paramName, paramIndex)
	case abi.BoolTy:
		return b.validateBoolArg(arg, paramName, paramIndex)
	case abi.StringTy:
		return b.validateStringArg(arg, paramName, paramIndex)
	case abi.BytesTy:
		return b.validateBytesArg(arg, paramName, paramIndex)
	case abi.FixedBytesTy:
		return b.validateFixedBytesArg(arg, expectedType, paramName, paramIndex)
	case abi.SliceTy:
		return b.validateSliceArg(arg, expectedType, paramName, paramIndex)
	case abi.ArrayTy:
		return b.validateArrayArg(arg, expectedType, paramName, paramIndex)
	default:
		return fmt.Errorf("parameter %s (index %d): unsupported type %s", paramName, paramIndex, expectedType.String())
	}
}

func (b *ABICallDataBuilder) validateIntegerArg(arg interface{}, expectedType abi.Type, paramName string, paramIndex int) error {
	switch v := arg.(type) {
	case string:
		if v == "" {
			return fmt.Errorf("parameter %s (index %d): empty string not valid for %s", paramName, paramIndex, expectedType.String())
		}
		// Try to parse as number to validate format
		if strings.HasPrefix(v, "0x") {
			if _, ok := new(big.Int).SetString(v[2:], 16); !ok {
				return fmt.Errorf("parameter %s (index %d): invalid hex number format '%s'", paramName, paramIndex, v)
			}
		} else {
			if _, ok := new(big.Int).SetString(v, 10); !ok {
				return fmt.Errorf("parameter %s (index %d): invalid number format '%s'", paramName, paramIndex, v)
			}
		}
	case float64:
		if v < 0 && expectedType.T == abi.UintTy {
			return fmt.Errorf("parameter %s (index %d): negative number not valid for %s", paramName, paramIndex, expectedType.String())
		}
	case int:
		if v < 0 && expectedType.T == abi.UintTy {
			return fmt.Errorf("parameter %s (index %d): negative number not valid for %s", paramName, paramIndex, expectedType.String())
		}
	default:
		return fmt.Errorf("parameter %s (index %d): expected string or number for %s, got %T", paramName, paramIndex, expectedType.String(), arg)
	}
	return nil
}

func (b *ABICallDataBuilder) validateAddressArg(arg interface{}, paramName string, paramIndex int) error {
	switch v := arg.(type) {
	case string:
		if v == "" {
			return fmt.Errorf("parameter %s (index %d): empty string not valid for address", paramName, paramIndex)
		}
		if !common.IsHexAddress(v) {
			return fmt.Errorf("parameter %s (index %d): invalid address format '%s'", paramName, paramIndex, v)
		}
	default:
		return fmt.Errorf("parameter %s (index %d): expected string for address, got %T", paramName, paramIndex, arg)
	}
	return nil
}

func (b *ABICallDataBuilder) validateBoolArg(arg interface{}, paramName string, paramIndex int) error {
	switch v := arg.(type) {
	case bool:
		// Valid
	case string:
		if _, err := strconv.ParseBool(v); err != nil {
			return fmt.Errorf("parameter %s (index %d): invalid bool format '%s'", paramName, paramIndex, v)
		}
	default:
		return fmt.Errorf("parameter %s (index %d): expected bool or string for bool, got %T", paramName, paramIndex, arg)
	}
	return nil
}

func (b *ABICallDataBuilder) validateStringArg(arg interface{}, paramName string, paramIndex int) error {
	// Almost anything can be converted to string, so this is quite permissive
	return nil
}

func (b *ABICallDataBuilder) validateBytesArg(arg interface{}, paramName string, paramIndex int) error {
	switch v := arg.(type) {
	case string:
		if strings.HasPrefix(v, "0x") {
			if len(v)%2 != 0 {
				return fmt.Errorf("parameter %s (index %d): hex string must have even length", paramName, paramIndex)
			}
		}
		// Non-hex strings are valid (will be converted to bytes)
	default:
		return fmt.Errorf("parameter %s (index %d): expected string for bytes, got %T", paramName, paramIndex, arg)
	}
	return nil
}

func (b *ABICallDataBuilder) validateFixedBytesArg(arg interface{}, expectedType abi.Type, paramName string, paramIndex int) error {
	switch v := arg.(type) {
	case string:
		var bytes []byte
		if strings.HasPrefix(v, "0x") {
			if len(v)%2 != 0 {
				return fmt.Errorf("parameter %s (index %d): hex string must have even length", paramName, paramIndex)
			}
			bytes = common.FromHex(v)
		} else {
			bytes = []byte(v)
		}

		if len(bytes) > expectedType.Size {
			return fmt.Errorf("parameter %s (index %d): bytes too long for %s (max %d bytes, got %d)",
				paramName, paramIndex, expectedType.String(), expectedType.Size, len(bytes))
		}
	default:
		return fmt.Errorf("parameter %s (index %d): expected string for %s, got %T", paramName, paramIndex, expectedType.String(), arg)
	}
	return nil
}

func (b *ABICallDataBuilder) validateSliceArg(arg interface{}, expectedType abi.Type, paramName string, paramIndex int) error {
	slice, ok := arg.([]interface{})
	if !ok {
		return fmt.Errorf("parameter %s (index %d): expected array for %s, got %T", paramName, paramIndex, expectedType.String(), arg)
	}

	// Validate each element
	for i, elem := range slice {
		elemParamName := fmt.Sprintf("%s[%d]", paramName, i)
		if err := b.validateArgType(elem, *expectedType.Elem, elemParamName, paramIndex); err != nil {
			return err
		}
	}

	return nil
}

func (b *ABICallDataBuilder) validateArrayArg(arg interface{}, expectedType abi.Type, paramName string, paramIndex int) error {
	slice, ok := arg.([]interface{})
	if !ok {
		return fmt.Errorf("parameter %s (index %d): expected array for %s, got %T", paramName, paramIndex, expectedType.String(), arg)
	}

	// Check fixed array size
	if len(slice) != expectedType.Size {
		return fmt.Errorf("parameter %s (index %d): expected array of size %d for %s, got %d elements",
			paramName, paramIndex, expectedType.Size, expectedType.String(), len(slice))
	}

	// Validate each element
	for i, elem := range slice {
		elemParamName := fmt.Sprintf("%s[%d]", paramName, i)
		if err := b.validateArgType(elem, *expectedType.Elem, elemParamName, paramIndex); err != nil {
			return err
		}
	}

	return nil
}
