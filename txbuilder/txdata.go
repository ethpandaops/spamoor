package txbuilder

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
)

// TxMetadata contains the common transaction parameters used across different transaction types.
// It provides a unified interface for specifying transaction details including gas parameters,
// recipient address, value transfer, transaction data, and various EIP extensions like
// access lists (EIP-2930), blob fees (EIP-4844), and authorization lists (EIP-7702).
type TxMetadata struct {
	GasTipCap  *uint256.Int // a.k.a. maxPriorityFeePerGas
	GasFeeCap  *uint256.Int // a.k.a. maxFeePerGas
	BlobFeeCap *uint256.Int // a.k.a. maxFeePerBlobGas
	Gas        uint64
	To         *common.Address
	Value      *uint256.Int
	Data       []byte
	AccessList types.AccessList
	AuthList   []types.SetCodeAuthorization
}
