package eestconv

// EEST fixture types

// EESTFixture represents a collection of test cases
type EESTFixture map[string]EESTTestCase

// EESTTestCase represents a single test case in the fixture
type EESTTestCase struct {
	Network            string                 `json:"network"`
	GenesisBlockHeader EESTBlockHeader        `json:"genesisBlockHeader"`
	Pre                map[string]EESTAccount `json:"pre"`
	PostState          map[string]EESTAccount `json:"postState"`
	Blocks             []EESTBlock            `json:"blocks"`
	Config             EESTConfig             `json:"config"`
	Info               map[string]any         `json:"_info"`
}

// EESTBlockHeader represents a block header
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

// EESTAccount represents an account state
type EESTAccount struct {
	Nonce   string            `json:"nonce"`
	Balance string            `json:"balance"`
	Code    string            `json:"code"`
	Storage map[string]string `json:"storage"`
}

// EESTBlock represents a block
type EESTBlock struct {
	BlockHeader  EESTBlockHeader   `json:"blockHeader"`
	Transactions []EESTTransaction `json:"transactions"`
}

// EESTTransaction represents a transaction
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

// EESTAccessListEntry represents an EIP-2930 access list entry
type EESTAccessListEntry struct {
	Address     string   `json:"address"`
	StorageKeys []string `json:"storageKeys"`
}

// EESTAuthorizationEntry represents an EIP-7702 authorization entry
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

// EESTConfig represents the test configuration
type EESTConfig struct {
	Network string `json:"network"`
	ChainID string `json:"chainid"`
}

// Converted output types

// ConvertedOutput is the top-level output structure
type ConvertedOutput struct {
	Payloads []ConvertedPayload `yaml:"payloads"`
}

// ConvertedPayload represents a single converted test case
type ConvertedPayload struct {
	Name          string                    `yaml:"name"`
	Prerequisites map[string]string         `yaml:"prerequisites"`
	Txs           []ConvertedTx             `yaml:"txs"`
	PostCheck     map[string]PostCheckEntry `yaml:"postcheck"`
}

// ConvertedTx represents a converted transaction
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

// ConvertedAccessListItem represents a converted access list entry
type ConvertedAccessListItem struct {
	Address     string   `yaml:"address"`
	StorageKeys []string `yaml:"storageKeys,omitempty"`
}

// ConvertedAuthorizationItem represents a converted authorization entry
type ConvertedAuthorizationItem struct {
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
	Storage map[string]string `yaml:"storage,omitempty"`
}
