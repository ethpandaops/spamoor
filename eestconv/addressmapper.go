package eestconv

import (
	"fmt"
	"regexp"
	"strings"
)

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

// NewAddressMapper creates a new AddressMapper instance
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

// RegisterContract registers a contract address and returns its index
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

// RegisterSender registers a sender address and returns its index
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

// GetContractIndex returns the contract index for an address
func (m *AddressMapper) GetContractIndex(addr string) (int, bool) {
	addrLower := strings.ToLower(addr)
	idx, exists := m.contractAddresses[addrLower]
	return idx, exists
}

// GetSenderIndex returns the sender index for an address
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

// ReplaceAddresses replaces all known addresses in a string with their placeholders
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
