package main

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
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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

// EEST fixture types

type EESTFixture map[string]EESTTestCase

type EESTTestCase struct {
	Network            string                 `json:"network"`
	GenesisBlockHeader EESTBlockHeader        `json:"genesisBlockHeader"`
	Pre                map[string]EESTAccount `json:"pre"`
	PostState          map[string]EESTAccount `json:"postState"`
	Blocks             []EESTBlock            `json:"blocks"`
	Config             EESTConfig             `json:"config"`
	Info               map[string]any         `json:"_info"`
}

type EESTBlockHeader struct {
	ParentHash    string `json:"parentHash"`
	Coinbase      string `json:"coinbase"`
	StateRoot     string `json:"stateRoot"`
	Number        string `json:"number"`
	GasLimit      string `json:"gasLimit"`
	GasUsed       string `json:"gasUsed"`
	Timestamp     string `json:"timestamp"`
	BaseFeePerGas string `json:"baseFeePerGas"`
}

type EESTAccount struct {
	Nonce   string            `json:"nonce"`
	Balance string            `json:"balance"`
	Code    string            `json:"code"`
	Storage map[string]string `json:"storage"`
}

type EESTBlock struct {
	BlockHeader  EESTBlockHeader   `json:"blockHeader"`
	Transactions []EESTTransaction `json:"transactions"`
}

type EESTTransaction struct {
	Type     string `json:"type"`
	ChainID  string `json:"chainId"`
	Nonce    string `json:"nonce"`
	GasPrice string `json:"gasPrice"`
	GasLimit string `json:"gasLimit"`
	To       string `json:"to"`
	Value    string `json:"value"`
	Data     string `json:"data"`
	V        string `json:"v"`
	R        string `json:"r"`
	S        string `json:"s"`
	Sender   string `json:"sender"`
	// EIP-1559 fields
	MaxFeePerGas         string `json:"maxFeePerGas"`
	MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas"`
	// EIP-2930 access list
	AccessList []EESTAccessListEntry `json:"accessList"`
	// EIP-4844 blob tx fields
	MaxFeePerBlobGas    string   `json:"maxFeePerBlobGas"`
	BlobVersionedHashes []string `json:"blobVersionedHashes"`
	// EIP-7702 set code tx fields
	AuthorizationList []EESTAuthorizationEntry `json:"authorizationList"`
}

type EESTAccessListEntry struct {
	Address     string   `json:"address"`
	StorageKeys []string `json:"storageKeys"`
}

type EESTAuthorizationEntry struct {
	ChainID string `json:"chainId"`
	Address string `json:"address"`
	Nonce   string `json:"nonce"`
	V       string `json:"v"`
	R       string `json:"r"`
	S       string `json:"s"`
	Signer  string `json:"signer"`
	YParity string `json:"yParity"`
}

type EESTConfig struct {
	Network string `json:"network"`
	ChainID string `json:"chainid"`
}

// Intermediate representation types

type ConvertedPayload struct {
	Name          string                    `yaml:"name"`
	Prerequisites map[string]string         `yaml:"prerequisites"`
	Txs           []ConvertedTx             `yaml:"txs"`
	PostCheck     map[string]PostCheckEntry `yaml:"postcheck"`
}

type ConvertedTx struct {
	From                 string                       `yaml:"from"`
	Type                 int                          `yaml:"type"`
	To                   string                       `yaml:"to,omitempty"`
	Data                 string                       `yaml:"data,omitempty"`
	Gas                  uint64                       `yaml:"gas"`
	GasPrice             uint64                       `yaml:"gasPrice,omitempty"`
	MaxFeePerGas         uint64                       `yaml:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas uint64                       `yaml:"maxPriorityFeePerGas,omitempty"`
	Value                string                       `yaml:"value,omitempty"`
	AccessList           []ConvertedAccessListItem    `yaml:"accessList,omitempty"`
	BlobCount            int                          `yaml:"blobCount,omitempty"`         // Number of blobs for type 3 tx
	AuthorizationList    []ConvertedAuthorizationItem `yaml:"authorizationList,omitempty"` // EIP-7702 authorizations for type 4 tx
	FixtureBaseFee       uint64                       `yaml:"fixtureBaseFee,omitempty"`    // Block base fee from fixture for balance scaling
}

type ConvertedAccessListItem struct {
	Address     string   `yaml:"address"`
	StorageKeys []string `yaml:"storageKeys,omitempty"`
}

type ConvertedAuthorizationItem struct {
	ChainID uint64 `yaml:"chainId"`
	Address string `yaml:"address"` // Address to delegate to (can be $contract[i])
	Nonce   uint64 `yaml:"nonce"`
	Signer  string `yaml:"signer,omitempty"` // Original signer address or $sender[i] if it's a sender
	V       string `yaml:"v,omitempty"`      // Original signature V (for non-sender signers)
	R       string `yaml:"r,omitempty"`      // Original signature R (for non-sender signers)
	S       string `yaml:"s,omitempty"`      // Original signature S (for non-sender signers)
}

type PostCheckEntry struct {
	Balance string            `yaml:"balance,omitempty"`
	Storage map[string]string `yaml:"storage,omitempty"`
}

type ConvertedOutput struct {
	Payloads []ConvertedPayload `yaml:"payloads"`
}

// AddressMapper tracks address to placeholder mappings
type AddressMapper struct {
	contractAddresses map[string]int // original address -> contract index
	senderAddresses   map[string]int // original address -> sender index
	contractIndex     int
	senderIndex       int

	// Reverse mappings for output
	indexToContract map[int]string
	indexToSender   map[int]string
}

func NewAddressMapper() *AddressMapper {
	return &AddressMapper{
		contractAddresses: make(map[string]int),
		senderAddresses:   make(map[string]int),
		contractIndex:     0,
		senderIndex:       0,
		indexToContract:   make(map[int]string),
		indexToSender:     make(map[int]string),
	}
}

func (m *AddressMapper) RegisterContract(addr string) int {
	addrLower := strings.ToLower(addr)
	if idx, exists := m.contractAddresses[addrLower]; exists {
		return idx
	}
	m.contractIndex++
	m.contractAddresses[addrLower] = m.contractIndex
	m.indexToContract[m.contractIndex] = addrLower
	return m.contractIndex
}

func (m *AddressMapper) RegisterSender(addr string) int {
	addrLower := strings.ToLower(addr)
	if idx, exists := m.senderAddresses[addrLower]; exists {
		return idx
	}
	m.senderIndex++
	m.senderAddresses[addrLower] = m.senderIndex
	m.indexToSender[m.senderIndex] = addrLower
	return m.senderIndex
}

func (m *AddressMapper) GetContractIndex(addr string) (int, bool) {
	addrLower := strings.ToLower(addr)
	idx, exists := m.contractAddresses[addrLower]
	return idx, exists
}

func (m *AddressMapper) GetSenderIndex(addr string) (int, bool) {
	addrLower := strings.ToLower(addr)
	idx, exists := m.senderAddresses[addrLower]
	return idx, exists
}

// GetBytecodeKey returns the placeholder for use in bytecode/addresses: $contract[i] or $sender[i]
func (m *AddressMapper) GetBytecodeKey(addr string) string {
	addrLower := strings.ToLower(addr)
	if idx, exists := m.contractAddresses[addrLower]; exists {
		return fmt.Sprintf("$contract[%d]", idx)
	}
	if idx, exists := m.senderAddresses[addrLower]; exists {
		return fmt.Sprintf("$sender[%d]", idx)
	}
	return addr
}

// GetYAMLKey returns the key for use in YAML maps (prerequisites, postcheck, from): contract[1] or sender[1]
func (m *AddressMapper) GetYAMLKey(addr string) string {
	addrLower := strings.ToLower(addr)
	if idx, exists := m.contractAddresses[addrLower]; exists {
		return fmt.Sprintf("contract[%d]", idx)
	}
	if idx, exists := m.senderAddresses[addrLower]; exists {
		return fmt.Sprintf("sender[%d]", idx)
	}
	return addr
}

func (m *AddressMapper) ReplaceAddresses(data string) string {
	result := data

	// Replace all contract addresses with $contract[i] format
	for addr, idx := range m.contractAddresses {
		addrNoPrefix := strings.TrimPrefix(addr, "0x")
		placeholder := fmt.Sprintf("$contract[%d]", idx)

		// Replace both with and without 0x prefix (case insensitive)
		// Use ReplaceAllLiteralString because $ has special meaning in replacement strings
		re := regexp.MustCompile("(?i)" + regexp.QuoteMeta(addrNoPrefix))
		result = re.ReplaceAllLiteralString(result, placeholder)
	}

	// Replace all sender addresses with $sender[i] format
	for addr, idx := range m.senderAddresses {
		addrNoPrefix := strings.TrimPrefix(addr, "0x")
		placeholder := fmt.Sprintf("$sender[%d]", idx)

		// Use ReplaceAllLiteralString because $ has special meaning in replacement strings
		re := regexp.MustCompile("(?i)" + regexp.QuoteMeta(addrNoPrefix))
		result = re.ReplaceAllLiteralString(result, placeholder)
	}

	return result
}

// NewConvertEESTCmd creates the convert-eest command
func NewConvertEESTCmd(logger logrus.FieldLogger) *cobra.Command {
	var outputFile string
	var verbosePayloads bool

	cmd := &cobra.Command{
		Use:   "convert-eest <path>",
		Short: "Convert EEST fixtures to intermediate representation",
		Long: `Converts Ethereum Execution Spec Test (EEST) fixtures to an intermediate
representation that can be replayed on a normal network.

The intermediate representation includes:
- Unsigned transactions with placeholder addresses ($contract[1], $sender[1], etc.)
- Prerequisites (sender funding amounts)
- Post-execution checks (storage values, balances)`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inputPath := args[0]
			return runConvertEEST(logger, inputPath, outputFile, verbosePayloads)
		},
	}

	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (default: stdout)")
	cmd.Flags().BoolVar(&verbosePayloads, "verbose-payloads", false, "Print each converted payload name")

	return cmd
}

func runConvertEEST(logger logrus.FieldLogger, inputPath, outputFile string, verbose bool) error {
	logger.WithField("path", inputPath).Info("scanning for EEST fixtures")

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

		logger.WithField("file", relPath).Debug("processing fixture file")

		payloads, err := convertFixtureFile(logger, path, relPath, verbose)
		if err != nil {
			logger.WithError(err).WithField("file", relPath).Warn("failed to convert fixture")
			failedFileCount++
			return nil
		}

		allPayloads = append(allPayloads, payloads...)
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to walk path: %w", err)
	}

	output := ConvertedOutput{
		Payloads: allPayloads,
	}

	yamlData, err := yaml.Marshal(output)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	if outputFile != "" {
		err = os.WriteFile(outputFile, yamlData, 0644)
		if err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		logger.WithField("file", outputFile).Info("wrote output file")
	} else {
		fmt.Println(string(yamlData))
	}

	// Print summary
	logger.Info("--- Conversion Summary ---")
	logger.WithField("count", jsonFileCount).Info("JSON files processed")
	if failedFileCount > 0 {
		logger.WithField("count", failedFileCount).Warn("JSON files failed")
	}
	logger.WithField("count", len(allPayloads)).Info("test payloads converted")

	return nil
}

func convertFixtureFile(logger logrus.FieldLogger, filePath, relPath string, verbose bool) ([]ConvertedPayload, error) {
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

		payload, err := convertTestCase(logger, testCase, fullName)
		if err != nil {
			logger.WithError(err).WithField("test", testName).Warn("failed to convert test case")
			continue
		}

		if verbose {
			logger.WithField("payload", fullName).Info("converted payload")
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

func convertTestCase(logger logrus.FieldLogger, testCase EESTTestCase, name string) (ConvertedPayload, error) {
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
	for _, c := range contracts {
		mapper.RegisterContract(c.addr)
	}

	// Collect all contract bytecodes to check for sender address references
	var allBytecodes []string
	for _, c := range contracts {
		allBytecodes = append(allBytecodes, strings.ToLower(c.account.Code))
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
	for _, c := range contracts {
		yamlKey := mapper.GetYAMLKey(c.addr)
		code := strings.TrimPrefix(c.account.Code, "0x")

		// Create init code: prefix + runtime code
		initCode := "0x" + initCodePrefix + code

		// Replace any known addresses in the init code
		initCode = mapper.ReplaceAddresses(initCode)

		value := "0x0"
		if c.account.Balance != "" && c.account.Balance != "0x00" && c.account.Balance != "0x0" {
			value = c.account.Balance
		}

		tx := ConvertedTx{
			From:                 "deployer",
			Type:                 2,
			To:                   "", // contract creation
			Data:                 initCode,
			Gas:                  1000000,     // default gas
			MaxFeePerGas:         10000000000, // 10 gwei
			MaxPriorityFeePerGas: 1000000000,  // 1 gwei
			FixtureBaseFee:       genesisBaseFee,
		}

		if value != "0x0" {
			tx.Value = value
		}

		payload.Txs = append(payload.Txs, tx)

		logger.WithFields(logrus.Fields{
			"address": c.addr,
			"yamlKey": yamlKey,
		}).Debug("added deployment transaction")
	}

	// Convert block transactions
	for _, block := range testCase.Blocks {
		blockBaseFee := parseHexUint64(block.BlockHeader.BaseFeePerGas)
		for _, tx := range block.Transactions {
			convertedTx, err := convertTransaction(mapper, tx, blockBaseFee)
			if err != nil {
				logger.WithError(err).Warn("failed to convert transaction")
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
