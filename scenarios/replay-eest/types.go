package replayeest

import (
	"regexp"
)

// PayloadFile is the top-level structure of the YAML file
type PayloadFile struct {
	Payloads []Payload `yaml:"payloads"`
}

// Payload represents a single test case
type Payload struct {
	Name          string                    `yaml:"name"`
	Prerequisites map[string]string         `yaml:"prerequisites"` // sender[1] -> balance
	Txs           []Tx                      `yaml:"txs"`
	PostCheck     map[string]PostCheckEntry `yaml:"postcheck"` // contract[1]/sender[1] -> checks
}

// Tx represents a transaction in the intermediate format
type Tx struct {
	From                 string              `yaml:"from"` // "deployer" or "sender[1]", etc.
	Type                 int                 `yaml:"type"`
	To                   string              `yaml:"to,omitempty"`   // "" for contract creation, or "$contract[1]", etc.
	Data                 string              `yaml:"data,omitempty"` // hex string with placeholders
	Gas                  uint64              `yaml:"gas"`
	GasPrice             uint64              `yaml:"gasPrice,omitempty"`
	MaxFeePerGas         uint64              `yaml:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas uint64              `yaml:"maxPriorityFeePerGas,omitempty"`
	Value                string              `yaml:"value,omitempty"`
	AccessList           []AccessListItem    `yaml:"accessList,omitempty"`
	BlobCount            int                 `yaml:"blobCount,omitempty"`         // Number of blobs for type 3 tx
	AuthorizationList    []AuthorizationItem `yaml:"authorizationList,omitempty"` // EIP-7702 authorizations for type 4 tx
	FixtureBaseFee       uint64              `yaml:"fixtureBaseFee,omitempty"`    // Block base fee from fixture for balance scaling
}

// AccessListItem represents an EIP-2930 access list entry
type AccessListItem struct {
	Address     string   `yaml:"address"`
	StorageKeys []string `yaml:"storageKeys,omitempty"`
}

// AuthorizationItem represents an EIP-7702 authorization entry
type AuthorizationItem struct {
	ChainID uint64 `yaml:"chainId"`
	Address string `yaml:"address"` // Address to delegate to (can be $contract[i])
	Nonce   uint64 `yaml:"nonce"`
	Signer  string `yaml:"signer,omitempty"` // Original signer address or $sender[i] if it's a sender
	V       string `yaml:"v,omitempty"`      // Original signature V (for non-sender signers)
	R       string `yaml:"r,omitempty"`      // Original signature R (for non-sender signers)
	S       string `yaml:"s,omitempty"`      // Original signature S (for non-sender signers)
}

// PostCheckEntry represents checks to perform after execution
type PostCheckEntry struct {
	Balance string            `yaml:"balance,omitempty"`
	Storage map[string]string `yaml:"storage,omitempty"` // slot -> expected value
}

// GetSenderCount returns the number of distinct senders in this payload
func (p *Payload) GetSenderCount() int {
	senderSet := make(map[string]bool)
	senderRegex := regexp.MustCompile(`sender\[\d+\]`)

	for _, tx := range p.Txs {
		if senderRegex.MatchString(tx.From) {
			senderSet[tx.From] = true
		}
	}

	// Also count senders from prerequisites
	for key := range p.Prerequisites {
		if senderRegex.MatchString(key) {
			senderSet[key] = true
		}
	}

	return len(senderSet)
}

// GetContractCount returns the number of contract deployments in this payload
func (p *Payload) GetContractCount() int {
	count := 0
	for _, tx := range p.Txs {
		if tx.From == "deployer" && tx.To == "" {
			count++
		}
	}
	return count
}
