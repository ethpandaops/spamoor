package txbuilder

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
)

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
