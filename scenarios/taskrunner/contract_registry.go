package taskrunner

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

// ContractRegistry manages deployed contract addresses for task references
type ContractRegistry struct {
	contracts map[string]common.Address
	parent    *ContractRegistry // Reference to parent registry (for inheritance)
}

// NewContractRegistry creates a new empty contract registry
func NewContractRegistry() *ContractRegistry {
	return &ContractRegistry{
		contracts: make(map[string]common.Address),
	}
}

// Clone creates a new registry that inherits from this one
// This is used to create execution-scoped registries that can see init contracts
func (r *ContractRegistry) Clone() *ContractRegistry {
	return &ContractRegistry{
		contracts: make(map[string]common.Address),
		parent:    r,
	}
}

// Set stores a contract address with the given name
func (r *ContractRegistry) Set(name string, address common.Address) {
	if name != "" {
		r.contracts[name] = address
	}
}

// Get retrieves a contract address by name
// It first checks the local registry, then the parent registry
func (r *ContractRegistry) Get(name string) (common.Address, bool) {
	// Check local contracts first
	if addr, ok := r.contracts[name]; ok {
		return addr, true
	}

	// Fall back to parent registry (init contracts)
	if r.parent != nil {
		return r.parent.Get(name)
	}

	return common.Address{}, false
}

// Has checks if a contract with the given name exists
func (r *ContractRegistry) Has(name string) bool {
	_, exists := r.Get(name)
	return exists
}

// ResolveReference resolves a contract reference string to an address
// Supports both direct addresses and contract references like {contract:name}
func (r *ContractRegistry) ResolveReference(ref string) (common.Address, error) {
	if ref == "" {
		return common.Address{}, fmt.Errorf("empty reference")
	}

	// Check if it's a contract reference: {contract:name}
	if strings.HasPrefix(ref, "{contract:") && strings.HasSuffix(ref, "}") {
		name := ref[10 : len(ref)-1] // Extract name between {contract: and }
		addr, ok := r.Get(name)
		if !ok {
			return common.Address{}, fmt.Errorf("contract '%s' not found in registry", name)
		}
		return addr, nil
	}

	// Otherwise treat as direct address
	if !common.IsHexAddress(ref) {
		return common.Address{}, fmt.Errorf("invalid address format: %s", ref)
	}

	return common.HexToAddress(ref), nil
}

// ListContracts returns all contract names and addresses in this registry
// (excluding parent registry)
func (r *ContractRegistry) ListContracts() map[string]common.Address {
	result := make(map[string]common.Address)

	// Add all contracts
	for name, addr := range r.contracts {
		result[name] = addr
	}

	return result
}

// ListAllContracts returns all contracts including from parent registries
func (r *ContractRegistry) ListAllContracts() map[string]common.Address {
	result := make(map[string]common.Address)

	// Start with parent contracts (if any)
	if r.parent != nil {
		parentContracts := r.parent.ListAllContracts()
		for name, addr := range parentContracts {
			result[name] = addr
		}
	}

	// Override with local contracts
	localContracts := r.ListContracts()
	for name, addr := range localContracts {
		result[name] = addr
	}

	return result
}
