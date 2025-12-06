// Package eestconv provides functionality to convert Ethereum Execution Spec Test (EEST)
// fixtures to an intermediate representation that can be replayed on a normal network.
package eestconv

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
)

// System contracts to filter out
var systemContracts = map[string]bool{
	"0x00000000219ab540356cbb839cbe05303d7705fa": true,
	"0x00000961ef480eb55e80d19ad83579a64c007002": true,
	"0x0000bbddc7ce488642fb579f8b00f3a590007251": true,
	"0x0000f90827f1c53a10cb7a02335b175320002935": true,
	"0x000f3df6d732807ef1319fb7b8bb8522d0beac02": true,
}

// Init code prefix for deploying contracts
const initCodePrefix = "600b380380600b5f395ff3"

// ConvertOptions holds options for the conversion process
type ConvertOptions struct {
	// TestPattern is a regex pattern to include tests by path/name
	TestPattern string
	// ExcludePattern is a regex pattern to exclude tests by path/name
	ExcludePattern string
	// Verbose enables verbose logging of each converted payload
	Verbose bool
}

// Converter handles the conversion of EEST fixtures
type Converter struct {
	logger         logrus.FieldLogger
	options        ConvertOptions
	pattern        *regexp.Regexp
	excludePattern *regexp.Regexp
}

// NewConverter creates a new Converter instance
func NewConverter(logger logrus.FieldLogger, options ConvertOptions) (*Converter, error) {
	c := &Converter{
		logger:  logger,
		options: options,
	}

	if options.TestPattern != "" {
		pattern, err := regexp.Compile(options.TestPattern)
		if err != nil {
			return nil, fmt.Errorf("invalid test pattern: %w", err)
		}
		c.pattern = pattern
	}

	if options.ExcludePattern != "" {
		excludePattern, err := regexp.Compile(options.ExcludePattern)
		if err != nil {
			return nil, fmt.Errorf("invalid exclude pattern: %w", err)
		}
		c.excludePattern = excludePattern
	}

	return c, nil
}

// ConvertDirectory converts all EEST fixtures in a directory
func (c *Converter) ConvertDirectory(inputPath string) (*ConvertedOutput, error) {
	c.logger.WithField("path", inputPath).Info("scanning for EEST fixtures")

	var allPayloads []ConvertedPayload
	jsonFileCount := 0
	failedFileCount := 0

	err := filepath.Walk(inputPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, ".json") {
			return nil
		}

		jsonFileCount++

		relPath, err := filepath.Rel(inputPath, path)
		if err != nil {
			relPath = path
		}

		c.logger.WithField("file", relPath).Debug("processing fixture file")

		payloads, err := c.convertFixtureFile(path, relPath)
		if err != nil {
			c.logger.WithError(err).WithField("file", relPath).Warn("failed to convert fixture")
			failedFileCount++
			return nil
		}

		allPayloads = append(allPayloads, payloads...)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk path: %w", err)
	}

	// Print summary
	c.logger.Info("--- Conversion Summary ---")
	c.logger.WithField("count", jsonFileCount).Info("JSON files processed")
	if failedFileCount > 0 {
		c.logger.WithField("count", failedFileCount).Warn("JSON files failed")
	}
	c.logger.WithField("count", len(allPayloads)).Info("test payloads converted")

	return &ConvertedOutput{Payloads: allPayloads}, nil
}

func (c *Converter) convertFixtureFile(filePath, relPath string) ([]ConvertedPayload, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var fixture EESTFixture
	if err := json.Unmarshal(data, &fixture); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	var payloads []ConvertedPayload

	// Sort test case names for deterministic output
	testNames := make([]string, 0, len(fixture))
	for name := range fixture {
		testNames = append(testNames, name)
	}
	sort.Strings(testNames)

	for _, testName := range testNames {
		testCase := fixture[testName]

		// Create a shortened name by extracting the test function name
		shortName := extractShortName(testName)
		fullName := fmt.Sprintf("%s/%s", strings.TrimSuffix(relPath, ".json"), shortName)

		// Apply test pattern filter if specified
		if c.pattern != nil && !c.pattern.MatchString(fullName) {
			continue
		}

		// Apply exclude pattern filter if specified
		if c.excludePattern != nil && c.excludePattern.MatchString(fullName) {
			continue
		}

		payload, err := c.convertTestCase(testCase, fullName)
		if err != nil {
			c.logger.WithError(err).WithField("test", testName).Warn("failed to convert test case")
			continue
		}

		if c.options.Verbose {
			c.logger.WithField("payload", fullName).Info("converted payload")
		}

		payloads = append(payloads, payload)
	}

	return payloads, nil
}

func extractShortName(fullTestName string) string {
	// Extract test name from format like:
	// "tests/osaka/eip7939_count_leading_zeros/test_count_leading_zeros.py::test_clz_call_operation[fork_Amsterdam-blockchain_test_from_state_test-call]"
	parts := strings.Split(fullTestName, "::")
	if len(parts) >= 2 {
		return parts[len(parts)-1]
	}
	return fullTestName
}

func (c *Converter) convertTestCase(testCase EESTTestCase, name string) (ConvertedPayload, error) {
	mapper := NewAddressMapper()
	payload := ConvertedPayload{
		Name:          name,
		Prerequisites: make(map[string]string),
		Txs:           []ConvertedTx{},
		PostCheck:     make(map[string]PostCheckEntry),
	}

	// Get genesis base fee for deployment transactions
	genesisBaseFee := parseHexUint64(testCase.GenesisBlockHeader.BaseFeePerGas)

	// First pass: identify all addresses
	// 1. Contracts: pre-state accounts with code (excluding system contracts)
	// 2. Potential senders: pre-state accounts with only balance (no code)

	type preAccount struct {
		addr    string
		account EESTAccount
	}

	var contracts []preAccount
	potentialSenders := make(map[string]EESTAccount) // addr -> account

	// Sort addresses for deterministic processing
	preAddrs := make([]string, 0, len(testCase.Pre))
	for addr := range testCase.Pre {
		preAddrs = append(preAddrs, addr)
	}
	sort.Strings(preAddrs)

	for _, addr := range preAddrs {
		account := testCase.Pre[addr]
		addrLower := strings.ToLower(addr)

		// Skip system contracts
		if systemContracts[addrLower] {
			continue
		}

		code := strings.TrimPrefix(account.Code, "0x")
		if code != "" {
			contracts = append(contracts, preAccount{addr: addr, account: account})
		} else {
			potentialSenders[addrLower] = account
		}
	}

	// Register contracts first (they get $contract[1], $contract[2], etc.)
	for _, contract := range contracts {
		mapper.RegisterContract(contract.addr)
	}

	// Collect all contract bytecodes to check for sender address references
	var allBytecodes []string
	for _, contract := range contracts {
		allBytecodes = append(allBytecodes, strings.ToLower(contract.account.Code))
	}

	// Collect all transaction senders and to addresses
	usedSenderAddrs := make(map[string]bool)
	for _, block := range testCase.Blocks {
		for _, tx := range block.Transactions {
			if tx.Sender != "" {
				usedSenderAddrs[strings.ToLower(tx.Sender)] = true
			}
			if tx.To != "" {
				toLower := strings.ToLower(tx.To)
				// Check if 'to' is a potential sender (not a contract)
				if _, isPotentialSender := potentialSenders[toLower]; isPotentialSender {
					usedSenderAddrs[toLower] = true
				}
			}
			// Check if sender addresses appear in tx data
			dataLower := strings.ToLower(tx.Data)
			for senderAddr := range potentialSenders {
				addrNoPrefix := strings.TrimPrefix(senderAddr, "0x")
				if strings.Contains(dataLower, addrNoPrefix) {
					usedSenderAddrs[senderAddr] = true
				}
			}
		}
	}

	// Check if sender addresses appear in contract bytecodes
	for senderAddr := range potentialSenders {
		addrNoPrefix := strings.TrimPrefix(senderAddr, "0x")
		for _, bytecode := range allBytecodes {
			if strings.Contains(bytecode, addrNoPrefix) {
				usedSenderAddrs[senderAddr] = true
				break
			}
		}
	}

	// Register only the actually used senders (sorted for deterministic order)
	usedSenderList := make([]string, 0, len(usedSenderAddrs))
	for addr := range usedSenderAddrs {
		usedSenderList = append(usedSenderList, addr)
	}
	sort.Strings(usedSenderList)

	for _, addrLower := range usedSenderList {
		account := potentialSenders[addrLower]
		idx := mapper.RegisterSender(addrLower)

		// Add to prerequisites if they have balance
		balance := parseHexBigInt(account.Balance)
		if balance != nil && balance.Sign() > 0 {
			key := fmt.Sprintf("sender[%d]", idx)
			payload.Prerequisites[key] = account.Balance
		}
	}

	// Generate deployment transactions for contracts
	for _, contract := range contracts {
		yamlKey := mapper.GetYAMLKey(contract.addr)
		code := strings.TrimPrefix(contract.account.Code, "0x")

		// Create init code: prefix + runtime code
		initCode := "0x" + initCodePrefix + code

		// Replace any known addresses in the init code
		initCode = mapper.ReplaceAddresses(initCode)

		value := "0x0"
		if contract.account.Balance != "" && contract.account.Balance != "0x00" && contract.account.Balance != "0x0" {
			value = contract.account.Balance
		}

		// Estimate gas for deployment based on initcode size
		gasLimit := estimateDeploymentGas(initCode)

		tx := ConvertedTx{
			From:                 "deployer",
			Type:                 2,
			To:                   "", // contract creation
			Data:                 initCode,
			Gas:                  gasLimit,
			MaxFeePerGas:         10000000000, // 10 gwei
			MaxPriorityFeePerGas: 1000000000,  // 1 gwei
			FixtureBaseFee:       genesisBaseFee,
		}

		if value != "0x0" {
			tx.Value = value
		}

		payload.Txs = append(payload.Txs, tx)

		c.logger.WithFields(logrus.Fields{
			"address": contract.addr,
			"yamlKey": yamlKey,
		}).Debug("added deployment transaction")
	}

	// Convert block transactions
	for _, block := range testCase.Blocks {
		blockBaseFee := parseHexUint64(block.BlockHeader.BaseFeePerGas)
		for _, tx := range block.Transactions {
			convertedTx, err := convertTransaction(mapper, tx, blockBaseFee)
			if err != nil {
				c.logger.WithError(err).Warn("failed to convert transaction")
				continue
			}
			payload.Txs = append(payload.Txs, convertedTx)
		}
	}

	// Convert post-state to postchecks
	postAddrs := make([]string, 0, len(testCase.PostState))
	for addr := range testCase.PostState {
		postAddrs = append(postAddrs, addr)
	}
	sort.Strings(postAddrs)

	for _, addr := range postAddrs {
		account := testCase.PostState[addr]
		addrLower := strings.ToLower(addr)

		// Skip system contracts
		if systemContracts[addrLower] {
			continue
		}

		// Get YAML key for this address (contract1, sender1, etc.)
		yamlKey := mapper.GetYAMLKey(addr)

		// If not mapped, skip (it might be a new address created during execution)
		if yamlKey == addr {
			continue
		}

		entry := PostCheckEntry{
			Storage: make(map[string]string),
		}

		// Add balance check if non-zero
		balance := parseHexBigInt(account.Balance)
		if balance != nil && balance.Sign() > 0 {
			entry.Balance = account.Balance
		}

		// Add storage checks
		for slot, value := range account.Storage {
			// Replace addresses in storage values
			replacedValue := mapper.ReplaceAddresses(value)
			entry.Storage[slot] = replacedValue
		}

		// Only add if there's something to check
		if entry.Balance != "" || len(entry.Storage) > 0 {
			payload.PostCheck[yamlKey] = entry
		}
	}

	return payload, nil
}

func convertTransaction(mapper *AddressMapper, tx EESTTransaction, blockBaseFee uint64) (ConvertedTx, error) {
	// Determine sender key (sender1, sender2, etc. - no $ prefix for YAML)
	senderKey := mapper.GetYAMLKey(tx.Sender)

	// Determine target address (use $contract[i] format for bytecode replacement)
	toAddr := ""
	if tx.To != "" {
		toAddr = mapper.GetBytecodeKey(tx.To)
		// If it still looks like an address, replace addresses in it
		if strings.HasPrefix(toAddr, "0x") {
			toAddr = mapper.ReplaceAddresses(toAddr)
		}
	}

	// Replace addresses in call data
	data := tx.Data
	if data != "" {
		data = mapper.ReplaceAddresses(data)
	}

	// Parse transaction type
	txType := 0
	if tx.Type != "" {
		txType = int(parseHexUint64(tx.Type))
	}

	// Parse gas
	gas := parseHexUint64(tx.GasLimit)
	if gas == 0 {
		gas = 21000
	}

	// Parse value
	value := ""
	valueBig := parseHexBigInt(tx.Value)
	if valueBig != nil && valueBig.Sign() > 0 {
		value = tx.Value
	}

	converted := ConvertedTx{
		From:           senderKey,
		Type:           txType,
		Data:           data,
		Gas:            gas,
		FixtureBaseFee: blockBaseFee,
	}

	if toAddr != "" {
		converted.To = toAddr
	}

	if value != "" {
		converted.Value = value
	}

	// Handle gas pricing based on transaction type
	switch txType {
	case 0, 1: // Legacy or Access List
		gasPrice := parseHexUint64(tx.GasPrice)
		if gasPrice > 0 {
			converted.GasPrice = gasPrice
		}
	case 2, 3, 4: // EIP-1559, Blob, SetCode
		maxFee := parseHexUint64(tx.MaxFeePerGas)
		maxPriority := parseHexUint64(tx.MaxPriorityFeePerGas)
		if maxFee > 0 {
			converted.MaxFeePerGas = maxFee
		}
		if maxPriority > 0 {
			converted.MaxPriorityFeePerGas = maxPriority
		}
	}

	// Convert access list for type 1, 2, 3, and 4 transactions
	if len(tx.AccessList) > 0 && txType >= 1 && txType <= 4 {
		accessList := make([]ConvertedAccessListItem, len(tx.AccessList))
		for i, entry := range tx.AccessList {
			// Replace addresses in access list entries
			address := mapper.ReplaceAddresses(entry.Address)
			storageKeys := make([]string, len(entry.StorageKeys))
			copy(storageKeys, entry.StorageKeys)
			accessList[i] = ConvertedAccessListItem{
				Address:     address,
				StorageKeys: storageKeys,
			}
		}
		converted.AccessList = accessList
	}

	// Convert blob count for type 3 transactions
	if txType == 3 && len(tx.BlobVersionedHashes) > 0 {
		converted.BlobCount = len(tx.BlobVersionedHashes)
	}

	// Convert authorization list for type 4 transactions
	if txType == 4 && len(tx.AuthorizationList) > 0 {
		authList := make([]ConvertedAuthorizationItem, len(tx.AuthorizationList))
		for i, entry := range tx.AuthorizationList {
			// Check if the signer is a known sender - if so, use placeholder
			signerAddr := entry.Signer
			_, isSender := mapper.GetSenderIndex(signerAddr)
			if isSender {
				signerAddr = mapper.GetYAMLKey(signerAddr)
			}

			// Replace addresses in the target address
			address := mapper.ReplaceAddresses(entry.Address)

			authItem := ConvertedAuthorizationItem{
				ChainID: parseHexUint64(entry.ChainID),
				Address: address,
				Nonce:   parseHexUint64(entry.Nonce),
				Signer:  signerAddr,
			}

			// For non-sender signers, preserve the original signature
			if !isSender {
				authItem.V = entry.V
				authItem.R = entry.R
				authItem.S = entry.S
			}

			authList[i] = authItem
		}
		converted.AuthorizationList = authList
	}

	return converted, nil
}

func parseHexUint64(s string) uint64 {
	if s == "" {
		return 0
	}
	s = strings.TrimPrefix(s, "0x")
	if s == "" {
		return 0
	}

	val, err := hex.DecodeString(padHexString(s))
	if err != nil {
		return 0
	}

	result := uint64(0)
	for _, b := range val {
		result = result<<8 | uint64(b)
	}
	return result
}

func parseHexBigInt(s string) *big.Int {
	if s == "" {
		return nil
	}
	s = strings.TrimPrefix(s, "0x")
	if s == "" {
		return big.NewInt(0)
	}

	val := new(big.Int)
	val.SetString(s, 16)
	return val
}

func padHexString(s string) string {
	if len(s)%2 != 0 {
		return "0" + s
	}
	return s
}

// estimateDeploymentGas estimates the gas needed for contract deployment
// based on the initcode. Our initcode is a simple copy routine that just
// copies the runtime code to memory and returns it.
//
// Gas costs breakdown:
// - Transaction base cost: 21000
// - Contract creation cost: 32000
// - Calldata cost: 16 gas per non-zero byte, 4 gas per zero byte
// - Code deposit cost: 200 gas per byte of runtime code (EIP-3860 style)
// - Execution overhead for the init code: ~100 gas for our simple copier
// - Memory expansion: roughly 3 gas per word for runtime code
//
// Adds a static 20k gas buffer, capped at 16M gas limit.
func estimateDeploymentGas(initCode string) uint64 {
	const maxGasLimit = uint64(16_000_000)
	const staticBuffer = uint64(20_000)

	// Remove 0x prefix if present
	code := strings.TrimPrefix(initCode, "0x")

	// Calculate byte length
	byteLen := len(code) / 2
	if byteLen == 0 {
		return 100000 // minimum gas for empty deployment
	}

	// Decode to count zero vs non-zero bytes for calldata cost
	decoded, err := hex.DecodeString(code)
	if err != nil {
		// Fallback: assume worst case (all non-zero)
		decoded = make([]byte, byteLen)
		for i := range decoded {
			decoded[i] = 0xff
		}
	}

	// Count zero and non-zero bytes for calldata cost
	var zeroBytes, nonZeroBytes uint64
	for _, b := range decoded {
		if b == 0 {
			zeroBytes++
		} else {
			nonZeroBytes++
		}
	}

	// Runtime code length (initcode minus the 11-byte prefix)
	// Our prefix is "600b380380600b5f395ff3" = 11 bytes
	runtimeCodeLen := uint64(byteLen)
	if runtimeCodeLen > 11 {
		runtimeCodeLen -= 11
	}

	// Calculate gas components
	baseCost := uint64(21000)                         // Transaction base
	createCost := uint64(32000)                       // Contract creation
	calldataCost := nonZeroBytes*16 + zeroBytes*4     // Calldata gas
	codeDepositCost := runtimeCodeLen * 200           // Code deposit (200 per byte)
	initCodeExecution := uint64(100)                  // Simple copy routine overhead
	memoryExpansion := (runtimeCodeLen + 31) / 32 * 3 // Memory expansion (~3 per word)

	totalGas := baseCost + createCost + calldataCost + codeDepositCost + initCodeExecution + memoryExpansion

	// Add static 20k gas buffer
	totalGas += staticBuffer

	// Cap at 16M gas limit
	if totalGas > maxGasLimit {
		totalGas = maxGasLimit
	}

	return totalGas
}
