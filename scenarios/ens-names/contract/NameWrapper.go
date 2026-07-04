// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// NameWrapperMetaData contains all meta data concerning the NameWrapper contract.
var NameWrapperMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractENS\",\"name\":\"_ens\",\"type\":\"address\"},{\"internalType\":\"contractIBaseRegistrar\",\"name\":\"_registrar\",\"type\":\"address\"},{\"internalType\":\"contractIMetadataService\",\"name\":\"_metadataService\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"CannotUpgrade\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncompatibleParent\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"IncorrectTargetOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectTokenType\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"labelHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"expectedLabelhash\",\"type\":\"bytes32\"}],\"name\":\"LabelMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"label\",\"type\":\"string\"}],\"name\":\"LabelTooLong\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LabelTooShort\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NameIsNotWrapped\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"offset\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"OffsetOutOfBoundsError\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"}],\"name\":\"OperationProhibited\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"Unauthorised\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"controller\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"}],\"name\":\"ControllerChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"expiry\",\"type\":\"uint64\"}],\"name\":\"ExpiryExtended\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"fuses\",\"type\":\"uint32\"}],\"name\":\"FusesSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"NameUnwrapped\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"name\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"fuses\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"expiry\",\"type\":\"uint64\"}],\"name\":\"NameWrapped\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"}],\"name\":\"TransferBatch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"TransferSingle\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"value\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"URI\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"_tokens\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"fuseMask\",\"type\":\"uint32\"}],\"name\":\"allFusesBurned\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"}],\"name\":\"balanceOfBatch\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"canExtendSubnames\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"canModifyName\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"controllers\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ens\",\"outputs\":[{\"internalType\":\"contractENS\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"parentNode\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"labelhash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"expiry\",\"type\":\"uint64\"}],\"name\":\"extendExpiry\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getData\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"fuses\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"expiry\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"parentNode\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"labelhash\",\"type\":\"bytes32\"}],\"name\":\"isWrapped\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"}],\"name\":\"isWrapped\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"metadataService\",\"outputs\":[{\"internalType\":\"contractIMetadataService\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"names\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onERC721Received\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"label\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"wrappedOwner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"duration\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"resolver\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"ownerControlledFuses\",\"type\":\"uint16\"}],\"name\":\"registerAndWrapETH2LD\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"registrarExpiry\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registrar\",\"outputs\":[{\"internalType\":\"contractIBaseRegistrar\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"duration\",\"type\":\"uint256\"}],\"name\":\"renew\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"expires\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeBatchTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"parentNode\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"labelhash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"fuses\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"expiry\",\"type\":\"uint64\"}],\"name\":\"setChildFuses\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"controller\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"}],\"name\":\"setController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"},{\"internalType\":\"uint16\",\"name\":\"ownerControlledFuses\",\"type\":\"uint16\"}],\"name\":\"setFuses\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIMetadataService\",\"name\":\"_metadataService\",\"type\":\"address\"}],\"name\":\"setMetadataService\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"resolver\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"ttl\",\"type\":\"uint64\"}],\"name\":\"setRecord\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"resolver\",\"type\":\"address\"}],\"name\":\"setResolver\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"parentNode\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"label\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"fuses\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"expiry\",\"type\":\"uint64\"}],\"name\":\"setSubnodeOwner\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"parentNode\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"label\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"resolver\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"ttl\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"fuses\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"expiry\",\"type\":\"uint64\"}],\"name\":\"setSubnodeRecord\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"ttl\",\"type\":\"uint64\"}],\"name\":\"setTTL\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractINameWrapperUpgrade\",\"name\":\"_upgradeAddress\",\"type\":\"address\"}],\"name\":\"setUpgradeContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"parentNode\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"labelhash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"controller\",\"type\":\"address\"}],\"name\":\"unwrap\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"labelhash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"registrant\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"controller\",\"type\":\"address\"}],\"name\":\"unwrapETH2LD\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"name\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"upgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upgradeContract\",\"outputs\":[{\"internalType\":\"contractINameWrapperUpgrade\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"uri\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"name\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"wrappedOwner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"resolver\",\"type\":\"address\"}],\"name\":\"wrap\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"label\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"wrappedOwner\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"ownerControlledFuses\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"resolver\",\"type\":\"address\"}],\"name\":\"wrapETH2LD\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"expiry\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60c06040523480156200001157600080fd5b5060405162006545380380620065458339810160408190526200003491620002f8565b823362000041816200028f565b6040516302571be360e01b81527f91d1777781884d03a6757a803996e38de2a42967fb37eeaca72729271025a9e260048201526000906001600160a01b038416906302571be390602401602060405180830381865afa158015620000a9573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620000cf91906200034c565b604051630f41a04d60e11b81526001600160a01b03848116600483015291925090821690631e83409a906024016020604051808303816000875af11580156200011c573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000142919062000373565b505050506001600160a01b0383811660805282811660a052600580546001600160a01b031916918316919091179055600163fffeffff60a01b03197fafa26c20e8b3d9a2853d642cfe1021dae26242ffedfac91c97aab212c1a4b93b8190557fa6eef7e35abe7026729641147f7915573c7e97b47efa546f5f6e3230263bcb4955604080518082019091526001815260006020808301829052908052600690527f54cdd369e4e8a8515e52ca72ec816c2101831ad1f18bf44102ed171459c9b4f89062000210908262000432565b506040805180820190915260058152626cae8d60e31b6020808301919091527f93cdeb708b7545dc668eb9280176169d1c33cfd8ed6f04690a0bcc88a93fc4ae600052600690527ffb9e8e321b8a5ec48f12a7b41f22c6e595d761285c9eb19d8dda7c99edf1b54f9062000285908262000432565b50505050620004fe565b600080546001600160a01b038381166001600160a01b0319831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b6001600160a01b0381168114620002f557600080fd5b50565b6000806000606084860312156200030e57600080fd5b83516200031b81620002df565b60208501519093506200032e81620002df565b60408501519092506200034181620002df565b809150509250925092565b6000602082840312156200035f57600080fd5b81516200036c81620002df565b9392505050565b6000602082840312156200038657600080fd5b5051919050565b634e487b7160e01b600052604160045260246000fd5b600181811c90821680620003b857607f821691505b602082108103620003d957634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200042d57600081815260208120601f850160051c81016020861015620004085750805b601f850160051c820191505b81811015620004295782815560010162000414565b5050505b505050565b81516001600160401b038111156200044e576200044e6200038d565b62000466816200045f8454620003a3565b84620003df565b602080601f8311600181146200049e5760008415620004855750858301515b600019600386901b1c1916600185901b17855562000429565b600085815260208120601f198616915b82811015620004cf57888601518255948401946001909101908401620004ae565b5085821015620004ee5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b60805160a051615f3a6200060b6000396000818161050601528181610c1501528181610cef01528181610d7901528181611c6601528181611cfc01528181611daa01528181611ecc01528181611f4201528181611fc20152818161224401528181612380015281816124b2015281816126970152818161271d0152612f5201526000818161055301528181610b9b01528181610ed70152818161108b0152818161113d015281816115550152818161240501528181612537015281816127c8015281816129bf01528181612ccd0152818161317d0152818161322b015281816132f40152818161336d015281816139ca01528181613ae501528181613d4d015261438d0152615f3a6000f3fe608060405234801561001057600080fd5b506004361061031f5760003560e01c80636352211e116101a7578063c93ab3fd116100ee578063e985e9c511610097578063f242432a11610071578063f242432a146107d7578063f2fde38b146107ea578063fd0cd0d9146107fd57600080fd5b8063e985e9c514610768578063eb8ae530146107a4578063ed70554d146107b757600080fd5b8063d9a50c12116100c8578063d9a50c121461071f578063da8c229e14610732578063e0dba60f1461075557600080fd5b8063c93ab3fd146106e6578063cf408823146106f9578063d8c9921a1461070c57600080fd5b8063a22cb46511610150578063b6bcad261161012a578063b6bcad26146106ad578063c475abff146106c0578063c658e086146106d357600080fd5b8063a22cb46514610674578063a401498214610687578063adf4960a1461069a57600080fd5b80638b4dfa75116101815780638b4dfa751461063d5780638cf8b41e146106505780638da5cb5b1461066357600080fd5b80636352211e146105f65780636e5d6ad214610609578063715018a61461063557600080fd5b80631f4e15041161026b5780633f15457f116102145780634e1273f4116101ee5780634e1273f4146105b057806353095467146105d05780635d3590d5146105e357600080fd5b80633f15457f1461054e578063402906fc1461057557806341415eab1461059d57600080fd5b80632b20e397116102455780632b20e397146105015780632eb2c2d61461052857806333c69ea91461053b57600080fd5b80631f4e1504146104c857806320c38e2b146104db57806324c1af44146104ee57600080fd5b80630e4cd725116102cd578063150b7a02116102a7578063150b7a02146104765780631534e177146104a25780631896f70a146104b557600080fd5b80630e4cd7251461043d5780630e89341c1461045057806314ab90381461046357600080fd5b806306fdde03116102fe57806306fdde03146103b4578063081812fc146103fd578063095ea7b31461042857600080fd5b8062fdd58e146103245780630178fe3f1461034a57806301ffc9a714610391575b600080fd5b610337610332366004614d69565b610810565b6040519081526020015b60405180910390f35b61035d610358366004614d95565b6108cf565b604080516001600160a01b03909416845263ffffffff909216602084015267ffffffffffffffff1690820152606001610341565b6103a461039f366004614dc4565b6108ff565b6040519015158152602001610341565b6103f06040518060400160405280600b81526020017f4e616d655772617070657200000000000000000000000000000000000000000081525081565b6040516103419190614e31565b61041061040b366004614d95565b610958565b6040516001600160a01b039091168152602001610341565b61043b610436366004614d69565b61099d565b005b6103a461044b366004614e44565b6109e3565b6103f061045e366004614d95565b610a7d565b61043b610471366004614e91565b610aef565b610489610484366004614f06565b610c08565b6040516001600160e01b03199091168152602001610341565b61043b6104b0366004614f79565b610e1a565b61043b6104c3366004614e44565b610e44565b600754610410906001600160a01b031681565b6103f06104e9366004614d95565b610f06565b6103376104fc366004615071565b610fa0565b6104107f000000000000000000000000000000000000000000000000000000000000000081565b61043b610536366004615199565b6111b4565b61043b610549366004615247565b6114de565b6104107f000000000000000000000000000000000000000000000000000000000000000081565b61058861058336600461529f565b6116d3565b60405163ffffffff9091168152602001610341565b6103a46105ab366004614e44565b611775565b6105c36105be3660046152c2565b6117d2565b60405161034191906153c0565b600554610410906001600160a01b031681565b61043b6105f13660046153d3565b611910565b610410610604366004614d95565b6119aa565b61061c610617366004615414565b6119b5565b60405167ffffffffffffffff9091168152602001610341565b61043b611b0a565b61043b61064b366004615449565b611b1e565b61061c61065e36600461548b565b611cc8565b6000546001600160a01b0316610410565b61043b610682366004615514565b612094565b610337610695366004615542565b61217e565b6103a46106a83660046155c3565b612319565b61043b6106bb366004614f79565b61233e565b6103376106ce3660046155e6565b612596565b6103376106e1366004615608565b61288d565b61043b6106f436600461567b565b612a9a565b61043b6107073660046156e7565b612c0b565b61043b61071a36600461571f565b612dc4565b6103a461072d3660046155e6565b612ed4565b6103a4610740366004614f79565b60046020526000908152604090205460ff1681565b61043b610763366004615514565b612fe1565b6103a461077636600461574d565b6001600160a01b03918216600090815260026020908152604080832093909416825291909152205460ff1690565b61043b6107b236600461577b565b613049565b6103376107c5366004614d95565b60016020526000908152604090205481565b61043b6107e53660046157e3565b613414565b61043b6107f8366004614f79565b613531565b6103a461080b366004614d95565b6135be565b60006001600160a01b0383166108935760405162461bcd60e51b815260206004820152602b60248201527f455243313135353a2062616c616e636520717565727920666f7220746865207a60448201527f65726f206164647265737300000000000000000000000000000000000000000060648201526084015b60405180910390fd5b600061089e836119aa565b9050836001600160a01b0316816001600160a01b0316036108c35760019150506108c9565b60009150505b92915050565b60008181526001602052604090205460a081901c60c082901c6108f3838383613696565b90959094509092509050565b60006001600160e01b031982167fd82c42d800000000000000000000000000000000000000000000000000000000148061094957506001600160e01b03198216630a85bd0160e11b145b806108c957506108c9826136cd565b600080610964836119aa565b90506001600160a01b03811661097d5750600092915050565b6000838152600360205260409020546001600160a01b03165b9392505050565b60006109a8826108cf565b50915050603f1960408216016109d45760405163a2a7201360e01b81526004810183905260240161088a565b6109de838361374f565b505050565b60008080806109f1866108cf565b925092509250846001600160a01b0316836001600160a01b03161480610a3c57506001600160a01b0380841660009081526002602090815260408083209389168352929052205460ff165b80610a6057506001600160a01b038516610a5587610958565b6001600160a01b0316145b8015610a735750610a718282613899565b155b9695505050505050565b6005546040516303a24d0760e21b8152600481018390526060916001600160a01b031690630e89341c90602401600060405180830381865afa158015610ac7573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526108c9919081019061584c565b81610afa8133611775565b610b205760405163168ab55d60e31b81526004810182905233602482015260440161088a565b8260106000610b2e836108cf565b5091505063ffffffff8282161615610b5c5760405163a2a7201360e01b81526004810184905260240161088a565b6040517f14ab90380000000000000000000000000000000000000000000000000000000081526004810187905267ffffffffffffffff861660248201527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316906314ab9038906044015b600060405180830381600087803b158015610be857600080fd5b505af1158015610bfc573d6000803e3d6000fd5b50505050505050505050565b6000336001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614610c6c576040517f1931a53800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000808080610c7d868801886158c4565b83516020850120939750919550935091508890808214610cd3576040517fc65c3ccc000000000000000000000000000000000000000000000000000000008152600481018290526024810183905260440161088a565b604051630a3b53db60e21b8152600481018390523060248201527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316906328ed4f6c90604401600060405180830381600087803b158015610d3b57600080fd5b505af1158015610d4f573d6000803e3d6000fd5b5050604051636b727d4360e11b8152600481018d9052600092506276a70091506001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169063d6e4fa8690602401602060405180830381865afa158015610dc0573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610de4919061592c565b610dee919061595b565b9050610e0187878761ffff1684886138ca565b50630a85bd0160e11b9c9b505050505050505050505050565b610e22613a30565b600580546001600160a01b0319166001600160a01b0392909216919091179055565b81610e4f8133611775565b610e755760405163168ab55d60e31b81526004810182905233602482015260440161088a565b8260086000610e83836108cf565b5091505063ffffffff8282161615610eb15760405163a2a7201360e01b81526004810184905260240161088a565b604051630c4b7b8560e11b8152600481018790526001600160a01b0386811660248301527f00000000000000000000000000000000000000000000000000000000000000001690631896f70a90604401610bce565b60066020526000908152604090208054610f1f90615983565b80601f0160208091040260200160405190810160405280929190818152602001828054610f4b90615983565b8015610f985780601f10610f6d57610100808354040283529160200191610f98565b820191906000526020600020905b815481529060010190602001808311610f7b57829003601f168201915b505050505081565b600087610fad8133611775565b610fd35760405163168ab55d60e31b81526004810182905233602482015260440161088a565b8751602089012061100b8a82604080516020808201949094528082019290925280518083038201815260609092019052805191012090565b92506110178a84613a8a565b6110218386613bc9565b61102c8a848b613bfc565b506110398a848787613cc9565b935061104483613d0f565b6110fa576040516305ef2c7f60e41b8152600481018b9052602481018290523060448201526001600160a01b03888116606483015267ffffffffffffffff881660848301527f00000000000000000000000000000000000000000000000000000000000000001690635ef2c7f09060a401600060405180830381600087803b1580156110cf57600080fd5b505af11580156110e3573d6000803e3d6000fd5b505050506110f58a848b8b8989613dc8565b6111a7565b6040516305ef2c7f60e41b8152600481018b9052602481018290523060448201526001600160a01b03888116606483015267ffffffffffffffff881660848301527f00000000000000000000000000000000000000000000000000000000000000001690635ef2c7f09060a401600060405180830381600087803b15801561118157600080fd5b505af1158015611195573d6000803e3d6000fd5b505050506111a78a848b8b8989613dff565b5050979650505050505050565b815183511461122b5760405162461bcd60e51b815260206004820152602860248201527f455243313135353a2069647320616e6420616d6f756e7473206c656e6774682060448201527f6d69736d61746368000000000000000000000000000000000000000000000000606482015260840161088a565b6001600160a01b03841661128f5760405162461bcd60e51b815260206004820152602560248201527f455243313135353a207472616e7366657220746f20746865207a65726f206164604482015264647265737360d81b606482015260840161088a565b6001600160a01b0385163314806112c957506001600160a01b038516600090815260026020908152604080832033845290915290205460ff165b61133b5760405162461bcd60e51b815260206004820152603260248201527f455243313135353a207472616e736665722063616c6c6572206973206e6f742060448201527f6f776e6572206e6f7220617070726f7665640000000000000000000000000000606482015260840161088a565b60005b835181101561147157600084828151811061135b5761135b6159bd565b602002602001015190506000848381518110611379576113796159bd565b602002602001015190506000806000611391856108cf565b9250925092506113a2858383613ec3565b8360011480156113c357508a6001600160a01b0316836001600160a01b0316145b6114225760405162461bcd60e51b815260206004820152602a60248201527f455243313135353a20696e73756666696369656e742062616c616e636520666f60448201526939103a3930b739b332b960b11b606482015260840161088a565b60008581526001602052604090206001600160a01b038b1663ffffffff60a01b60a085901b16176001600160c01b031960c084901b1617905550505050508061146a906159d3565b905061133e565b50836001600160a01b0316856001600160a01b0316336001600160a01b03167f4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb86866040516114c19291906159ec565b60405180910390a46114d7338686868686613fb0565b5050505050565b604080516020808201879052818301869052825180830384018152606090920190925280519101206115108184613bc9565b6000808061151d846108cf565b919450925090506001600160a01b03831615806115cc57506040516302571be360e01b81526004810185905230906001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016906302571be390602401602060405180830381865afa15801561159c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906115c09190615a1a565b6001600160a01b031614155b156115ea57604051635374b59960e01b815260040160405180910390fd5b6000806115f68a6108cf565b90935091508a90506116375761160c8633611775565b6116325760405163168ab55d60e31b81526004810187905233602482015260440161088a565b611667565b6116418a33611775565b6116675760405163168ab55d60e31b8152600481018b905233602482015260440161088a565b611672868984614155565b61167d878483614190565b9650620100008416158015906116a157508363ffffffff1688851763ffffffff1614155b156116c25760405163a2a7201360e01b81526004810187905260240161088a565b96831796610bfc86868a868b6141da565b6000826116e08133611775565b6117065760405163168ab55d60e31b81526004810182905233602482015260440161088a565b8360026000611714836108cf565b5091505063ffffffff82821616156117425760405163a2a7201360e01b81526004810184905260240161088a565b6000808061174f8a6108cf565b9250925092506117688a84848c61ffff161784856141da565b5098975050505050505050565b6000808080611783866108cf565b925092509250846001600160a01b0316836001600160a01b03161480610a6057506001600160a01b0380841660009081526002602090815260408083209389168352929052205460ff16610a60565b6060815183511461184b5760405162461bcd60e51b815260206004820152602960248201527f455243313135353a206163636f756e747320616e6420696473206c656e67746860448201527f206d69736d617463680000000000000000000000000000000000000000000000606482015260840161088a565b6000835167ffffffffffffffff81111561186757611867614f96565b604051908082528060200260200182016040528015611890578160200160208202803683370190505b50905060005b8451811015611908576118db8582815181106118b4576118b46159bd565b60200260200101518583815181106118ce576118ce6159bd565b6020026020010151610810565b8282815181106118ed576118ed6159bd565b6020908102919091010152611901816159d3565b9050611896565b509392505050565b611918613a30565b6040517fa9059cbb0000000000000000000000000000000000000000000000000000000081526001600160a01b0383811660048301526024820183905284169063a9059cbb906044016020604051808303816000875af1158015611980573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906119a49190615a37565b50505050565b60006108c982614284565b604080516020808201869052818301859052825180830384018152606090920190925280519101206000906119e981613d0f565b611a0657604051635374b59960e01b815260040160405180910390fd5b6000611a1286336109e3565b905080158015611a295750611a278233611775565b155b15611a505760405163168ab55d60e31b81526004810183905233602482015260440161088a565b60008080611a5d856108cf565b92509250925083158015611a745750620400008216155b15611a955760405163a2a7201360e01b81526004810186905260240161088a565b6000611aa08a6108cf565b92505050611aaf888383614190565b9750611abd8685858b61429a565b60405167ffffffffffffffff8916815286907ff675815a0817338f93a7da433f6bd5f5542f1029b11b455191ac96c7f6a9b1329060200160405180910390a2509598975050505050505050565b611b12613a30565b611b1c60006142e2565b565b604080517f93cdeb708b7545dc668eb9280176169d1c33cfd8ed6f04690a0bcc88a93fc4ae60208083019190915281830186905282518083038401815260609092019092528051910120611b728133611775565b611b985760405163168ab55d60e31b81526004810182905233602482015260440161088a565b306001600160a01b03841603611bcc57604051632ca49b0d60e11b81526001600160a01b038416600482015260240161088a565b604080517f93cdeb708b7545dc668eb9280176169d1c33cfd8ed6f04690a0bcc88a93fc4ae60208083019190915281830187905282518083038401815260609092019092528051910120611c21905b83614332565b6040517f42842e0e0000000000000000000000000000000000000000000000000000000081523060048201526001600160a01b038481166024830152604482018690527f000000000000000000000000000000000000000000000000000000000000000016906342842e0e90606401600060405180830381600087803b158015611caa57600080fd5b505af1158015611cbe573d6000803e3d6000fd5b5050505050505050565b6000808686604051611cdb929190615a54565b6040519081900381206331a9108f60e11b82526004820181905291506000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031690636352211e90602401602060405180830381865afa158015611d4b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611d6f9190615a1a565b90506001600160a01b0381163314801590611e17575060405163e985e9c560e01b81526001600160a01b0382811660048301523360248301527f0000000000000000000000000000000000000000000000000000000000000000169063e985e9c590604401602060405180830381865afa158015611df1573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611e159190615a37565b155b15611e8757604080517f93cdeb708b7545dc668eb9280176169d1c33cfd8ed6f04690a0bcc88a93fc4ae6020808301919091528183018590528251808303840181526060830193849052805191012063168ab55d60e31b909252606481019190915233608482015260a40161088a565b6040517f23b872dd0000000000000000000000000000000000000000000000000000000081526001600160a01b038281166004830152306024830152604482018490527f000000000000000000000000000000000000000000000000000000000000000016906323b872dd90606401600060405180830381600087803b158015611f1057600080fd5b505af1158015611f24573d6000803e3d6000fd5b5050604051630a3b53db60e21b8152600481018590523060248201527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031692506328ed4f6c9150604401600060405180830381600087803b158015611f9057600080fd5b505af1158015611fa4573d6000803e3d6000fd5b5050604051636b727d4360e11b8152600481018590526276a70092507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316915063d6e4fa8690602401602060405180830381865afa158015612012573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612036919061592c565b612040919061595b565b925061208988888080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508a9250505061ffff881686886138ca565b505095945050505050565b6001600160a01b03821633036121125760405162461bcd60e51b815260206004820152602960248201527f455243313135353a2073657474696e6720617070726f76616c2073746174757360448201527f20666f722073656c660000000000000000000000000000000000000000000000606482015260840161088a565b3360008181526002602090815260408083206001600160a01b03871680855290835292819020805460ff191686151590811790915590519081529192917f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31910160405180910390a35050565b3360009081526004602052604081205460ff166121ee5760405162461bcd60e51b815260206004820152602860248201527f436f6e74726f6c6c61626c653a2043616c6c6572206973206e6f74206120636f604482015267373a3937b63632b960c11b606482015260840161088a565b60008787604051612200929190615a54565b6040519081900381207ffca247ac000000000000000000000000000000000000000000000000000000008252600482018190523060248301526044820187905291507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063fca247ac906064016020604051808303816000875af1158015612295573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906122b9919061592c565b915061230e88888080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508a9250505061ffff86166123086276a7008761595b565b886138ca565b509695505050505050565b600080612325846108cf565b50841663ffffffff908116908516149250505092915050565b612346613a30565b6007546001600160a01b0316156124665760075460405163a22cb46560e01b81526001600160a01b039182166004820152600060248201527f00000000000000000000000000000000000000000000000000000000000000009091169063a22cb46590604401600060405180830381600087803b1580156123c657600080fd5b505af11580156123da573d6000803e3d6000fd5b505060075460405163a22cb46560e01b81526001600160a01b039182166004820152600060248201527f0000000000000000000000000000000000000000000000000000000000000000909116925063a22cb4659150604401600060405180830381600087803b15801561244d57600080fd5b505af1158015612461573d6000803e3d6000fd5b505050505b600780546001600160a01b0319166001600160a01b038316908117909155156125935760075460405163a22cb46560e01b81526001600160a01b039182166004820152600160248201527f00000000000000000000000000000000000000000000000000000000000000009091169063a22cb46590604401600060405180830381600087803b1580156124f857600080fd5b505af115801561250c573d6000803e3d6000fd5b505060075460405163a22cb46560e01b81526001600160a01b039182166004820152600160248201527f0000000000000000000000000000000000000000000000000000000000000000909116925063a22cb4659150604401600060405180830381600087803b15801561257f57600080fd5b505af11580156114d7573d6000803e3d6000fd5b50565b3360009081526004602052604081205460ff166126065760405162461bcd60e51b815260206004820152602860248201527f436f6e74726f6c6c61626c653a2043616c6c6572206973206e6f74206120636f604482015267373a3937b63632b960c11b606482015260840161088a565b604080517f93cdeb708b7545dc668eb9280176169d1c33cfd8ed6f04690a0bcc88a93fc4ae602080830191909152818301869052825180830384018152606090920190925280519101206000906040517fc475abff00000000000000000000000000000000000000000000000000000000815260048101869052602481018590529091506000906001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169063c475abff906044016020604051808303816000875af11580156126e0573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612704919061592c565b6040516331a9108f60e11b8152600481018790529091507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031690636352211e90602401602060405180830381865afa925050508015612788575060408051601f3d908101601f1916820190925261278591810190615a1a565b60015b6127955791506108c99050565b6001600160a01b0381163014158061283f57506040516302571be360e01b81526004810184905230906001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016906302571be390602401602060405180830381865afa15801561280f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906128339190615a1a565b6001600160a01b031614155b1561284e575091506108c99050565b50600061285e6276a7008361595b565b60008481526001602052604090205490915060a081901c6128818583838661429a565b50919695505050505050565b60008661289a8133611775565b6128c05760405163168ab55d60e31b81526004810182905233602482015260440161088a565b600087876040516128d2929190615a54565b6040518091039020905061290d8982604080516020808201949094528082019290925280518083038201815260609092019052805191012090565b92506129198984613a8a565b6129238386613bc9565b60006129668a858b8b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250613bfc92505050565b90506129748a858888613cc9565b945061297f84613d0f565b612a47576040517f06ab5923000000000000000000000000000000000000000000000000000000008152600481018b9052602481018390523060448201527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316906306ab5923906064016020604051808303816000875af1158015612a10573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612a34919061592c565b50612a428482898989614424565b612a8d565b612a8d8a858b8b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508d92508c91508b9050613dff565b5050509695505050505050565b6000612ae0600086868080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525092939250506144669050565b6007549091506001600160a01b0316612b25576040517f24c1d6d400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b612b2f8133611775565b612b555760405163168ab55d60e31b81526004810182905233602482015260440161088a565b60008080612b62846108cf565b919450925090506000612b7485610958565b9050612b7f85614525565b600760009054906101000a90046001600160a01b03166001600160a01b0316639198c2768a8a878787878e8e6040518963ffffffff1660e01b8152600401612bce989796959493929190615a8d565b600060405180830381600087803b158015612be857600080fd5b505af1158015612bfc573d6000803e3d6000fd5b50505050505050505050505050565b83612c168133611775565b612c3c5760405163168ab55d60e31b81526004810182905233602482015260440161088a565b84601c6000612c4a836108cf565b5091505063ffffffff8282161615612c785760405163a2a7201360e01b81526004810184905260240161088a565b6040517fcf408823000000000000000000000000000000000000000000000000000000008152600481018990523060248201526001600160a01b03878116604483015267ffffffffffffffff871660648301527f0000000000000000000000000000000000000000000000000000000000000000169063cf40882390608401600060405180830381600087803b158015612d1157600080fd5b505af1158015612d25573d6000803e3d6000fd5b5050506001600160a01b0388169050612d8c576000612d43896108cf565b509150506201ffff1962020000821601612d7b57604051632ca49b0d60e11b81526001600160a01b038916600482015260240161088a565b612d86896000614332565b50611cbe565b6000612d97896119aa565b9050612db981898b60001c6001604051806020016040528060008152506145e7565b505050505050505050565b60408051602080820186905281830185905282518083038401815260609092019092528051910120612df68133611775565b612e1c5760405163168ab55d60e31b81526004810182905233602482015260440161088a565b7f6c32148f748aba23997146d7fe89e962e3cc30271290fb96f5f4337756c03b528401612e5c5760405163615a470360e01b815260040160405180910390fd5b6001600160a01b0382161580612e7a57506001600160a01b03821630145b15612ea357604051632ca49b0d60e11b81526001600160a01b038316600482015260240161088a565b604080516020808201879052818301869052825180830384018152606090920190925280519101206119a490611c1b565b604080516020808201859052818301849052825180830384018152606090920190925280519101206000906000612f0a82613d0f565b90507f93cdeb708b7545dc668eb9280176169d1c33cfd8ed6f04690a0bcc88a93fc4ae8514612f3c5791506108c99050565b6040516331a9108f60e11b8152600481018590527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031690636352211e90602401602060405180830381865afa925050508015612fbd575060408051601f3d908101601f19168201909252612fba91810190615a1a565b60015b612fcc576000925050506108c9565b6001600160a01b0316301492506108c9915050565b612fe9613a30565b6001600160a01b038216600081815260046020908152604091829020805460ff191685151590811790915591519182527f4c97694570a07277810af7e5669ffd5f6a2d6b74b6e9a274b8b870fd5114cf8791015b60405180910390a25050565b600080613090600087878080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525092939250506147399050565b9150915060006130d98288888080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525092939250506144669050565b604080516020808201849052818301879052825180830384018152606090920190925280519101209091506000906000818152600660205260409020909150613123888a83615b3c565b507f6c32148f748aba23997146d7fe89e962e3cc30271290fb96f5f4337756c03b5282016131645760405163615a470360e01b815260040160405180910390fd5b6040516302571be360e01b8152600481018290526000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316906302571be390602401602060405180830381865afa1580156131cc573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906131f09190615a1a565b90506001600160a01b0381163314801590613298575060405163e985e9c560e01b81526001600160a01b0382811660048301523360248301527f0000000000000000000000000000000000000000000000000000000000000000169063e985e9c590604401602060405180830381865afa158015613272573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906132969190615a37565b155b156132bf5760405163168ab55d60e31b81526004810183905233602482015260440161088a565b6001600160a01b0386161561335157604051630c4b7b8560e11b8152600481018390526001600160a01b0387811660248301527f00000000000000000000000000000000000000000000000000000000000000001690631896f70a90604401600060405180830381600087803b15801561333857600080fd5b505af115801561334c573d6000803e3d6000fd5b505050505b604051635b0fc9c360e01b8152600481018390523060248201527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031690635b0fc9c390604401600060405180830381600087803b1580156133b957600080fd5b505af11580156133cd573d6000803e3d6000fd5b50505050612db9828a8a8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201829052508d93509150819050614424565b6001600160a01b0384166134785760405162461bcd60e51b815260206004820152602560248201527f455243313135353a207472616e7366657220746f20746865207a65726f206164604482015264647265737360d81b606482015260840161088a565b6001600160a01b0385163314806134b257506001600160a01b038516600090815260026020908152604080832033845290915290205460ff165b6135245760405162461bcd60e51b815260206004820152602960248201527f455243313135353a2063616c6c6572206973206e6f74206f776e6572206e6f7260448201527f20617070726f7665640000000000000000000000000000000000000000000000606482015260840161088a565b6114d785858585856145e7565b613539613a30565b6001600160a01b0381166135b55760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201527f6464726573730000000000000000000000000000000000000000000000000000606482015260840161088a565b612593816142e2565b600081815260066020526040812080548291906135da90615983565b80601f016020809104026020016040519081016040528092919081815260200182805461360690615983565b80156136535780601f1061362857610100808354040283529160200191613653565b820191906000526020600020905b81548152906001019060200180831161363657829003601f168201915b50505050509050805160000361366c5750600092915050565b6000806136798382614739565b9092509050600061368a8483614466565b9050610a738184612ed4565b600080428367ffffffffffffffff1610156136c45761ffff19620100008516016136bf57600094505b600093505b50929391925050565b60006001600160e01b031982167fd9b67a2600000000000000000000000000000000000000000000000000000000148061371757506001600160e01b031982166303a24d0760e21b145b806108c957507f01ffc9a7000000000000000000000000000000000000000000000000000000006001600160e01b03198316146108c9565b600061375a826119aa565b9050806001600160a01b0316836001600160a01b0316036137e35760405162461bcd60e51b815260206004820152602160248201527f4552433732313a20617070726f76616c20746f2063757272656e74206f776e6560448201527f7200000000000000000000000000000000000000000000000000000000000000606482015260840161088a565b336001600160a01b038216148061381d57506001600160a01b038116600090815260026020908152604080832033845290915290205460ff165b61388f5760405162461bcd60e51b815260206004820152603d60248201527f4552433732313a20617070726f76652063616c6c6572206973206e6f7420746f60448201527f6b656e206f776e6572206f7220617070726f76656420666f7220616c6c000000606482015260840161088a565b6109de83836147f0565b6000620200008381161480156109965750426138b86276a70084615bfc565b67ffffffffffffffff16109392505050565b8451602086012060006139247f93cdeb708b7545dc668eb9280176169d1c33cfd8ed6f04690a0bcc88a93fc4ae83604080516020808201949094528082019290925280518083038201815260609092019052805191012090565b90506000613967886040518060400160405280600581526020017f036574680000000000000000000000000000000000000000000000000000000081525061485e565b60008381526006602052604090209091506139828282615c1d565b50613995828289620300008a1789614424565b6001600160a01b03841615611cbe57604051630c4b7b8560e11b8152600481018390526001600160a01b0385811660248301527f00000000000000000000000000000000000000000000000000000000000000001690631896f70a90604401600060405180830381600087803b158015613a0e57600080fd5b505af1158015613a22573d6000803e3d6000fd5b505050505050505050505050565b6000546001600160a01b03163314611b1c5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015260640161088a565b60008080613a97846108cf565b919450925090504267ffffffffffffffff821610808015613b5b57506001600160a01b0384161580613b5b57506040516302571be360e01b8152600481018690526000906001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016906302571be390602401602060405180830381865afa158015613b2c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613b509190615a1a565b6001600160a01b0316145b15613b9a576000613b6b876108cf565b509150506020811615613b945760405163a2a7201360e01b81526004810187905260240161088a565b50613bc1565b62010000831615613bc15760405163a2a7201360e01b81526004810186905260240161088a565b505050505050565b63fffdffff81811763ffffffff1614613bf85760405163a2a7201360e01b81526004810183905260240161088a565b5050565b60606000613ca583600660008881526020019081526020016000208054613c2290615983565b80601f0160208091040260200160405190810160405280929190818152602001828054613c4e90615983565b8015613c9b5780601f10613c7057610100808354040283529160200191613c9b565b820191906000526020600020905b815481529060010190602001808311613c7e57829003601f168201915b505050505061485e565b6000858152600660205260409020909150613cc08282615c1d565b50949350505050565b600080613cd5856108cf565b92505050600080613ce88860001c6108cf565b9250925050613cf8878784614155565b613d03858483614190565b98975050505050505050565b600080613d1b836119aa565b6001600160a01b0316141580156108c957506040516302571be360e01b81526004810183905230906001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016906302571be390602401602060405180830381865afa158015613d94573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613db89190615a1a565b6001600160a01b03161492915050565b60008681526006602052604081208054613de7918791613c2290615983565b9050613df68682868686614424565b50505050505050565b60008080613e0c886108cf565b9250925092506000613e3688600660008d81526020019081526020016000208054613c2290615983565b60008a8152600660205260409020805491925090613e5390615983565b9050600003613e76576000898152600660205260409020613e748282615c1d565b505b613e85898588861785896141da565b6001600160a01b038716613ea357613e9e896000614332565b610bfc565b610bfc84888b60001c6001604051806020016040528060008152506145e7565b6201ffff1962020000831601613ee357613ee06276a70082615bfc565b90505b428167ffffffffffffffff161015613f605762010000821615613f5b5760405162461bcd60e51b815260206004820152602a60248201527f455243313135353a20696e73756666696369656e742062616c616e636520666f60448201526939103a3930b739b332b960b11b606482015260840161088a565b613f85565b6004821615613f855760405163a2a7201360e01b81526004810184905260240161088a565b604082166000036109de575050600090815260036020526040902080546001600160a01b0319169055565b6001600160a01b0384163b15613bc15760405163bc197c8160e01b81526001600160a01b0385169063bc197c8190613ff49089908990889088908890600401615cdd565b6020604051808303816000875af192505050801561402f575060408051601f3d908101601f1916820190925261402c91810190615d2f565b60015b6140e45761403b615d4c565b806308c379a003614074575061404f615d68565b8061405a5750614076565b8060405162461bcd60e51b815260040161088a9190614e31565b505b60405162461bcd60e51b815260206004820152603460248201527f455243313135353a207472616e7366657220746f206e6f6e204552433131353560448201527f526563656976657220696d706c656d656e746572000000000000000000000000606482015260840161088a565b6001600160e01b0319811663bc197c8160e01b14613df65760405162461bcd60e51b815260206004820152602860248201527f455243313135353a204552433131353552656365697665722072656a656374656044820152676420746f6b656e7360c01b606482015260840161088a565b63ffff0000821615801590600183161590829061416f5750805b156114d75760405163a2a7201360e01b81526004810186905260240161088a565b60008167ffffffffffffffff168467ffffffffffffffff1611156141b2578193505b8267ffffffffffffffff168467ffffffffffffffff1610156141d2578293505b509192915050565b6141e68585858461429a565b60405163ffffffff8416815285907f39873f00c80f4f94b7bd1594aebcf650f003545b74824d57ddf4939e3ff3a34b9060200160405180910390a28167ffffffffffffffff168167ffffffffffffffff1611156114d75760405167ffffffffffffffff8216815285907ff675815a0817338f93a7da433f6bd5f5542f1029b11b455191ac96c7f6a9b132906020015b60405180910390a25050505050565b600080614290836108cf565b5090949350505050565b6142a48483614907565b60008481526001602052604090206001600160a01b03841663ffffffff60a01b60a085901b16176001600160c01b031960c084901b161790556119a4565b600080546001600160a01b038381166001600160a01b0319831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b61433d826001612319565b1561435e5760405163a2a7201360e01b81526004810183905260240161088a565b61436782614525565b604051635b0fc9c360e01b8152600481018390526001600160a01b0382811660248301527f00000000000000000000000000000000000000000000000000000000000000001690635b0fc9c390604401600060405180830381600087803b1580156143d157600080fd5b505af11580156143e5573d6000803e3d6000fd5b50506040516001600160a01b03841681528492507fee2ba1195c65bcf218a83d874335c6bf9d9067b4c672f3c3bf16cf40de7586c4915060200161303d565b61443085848484614940565b847f8ce7013e8abebc55c3890a68f5a27c67c3f7efa64e584de5fb22363c606fd340858585856040516142759493929190615df2565b60008060006144758585614739565b9092509050816144e7576001855161448d9190615e3a565b84146144db5760405162461bcd60e51b815260206004820152601d60248201527f6e616d65686173683a204a756e6b20617420656e64206f66206e616d65000000604482015260640161088a565b50600091506108c99050565b6144f18582614466565b6040805160208101929092528101839052606001604051602081830303815290604052805190602001209250505092915050565b60008181526001602052604090205460a081901c60c082901c614549838383613696565b600086815260036020908152604080832080546001600160a01b03191690556001909152902063ffffffff60a01b60a083901b166001600160c01b031960c086901b1617905592506145989050565b60408051858152600160208201526000916001600160a01b0386169133917fc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62910160405180910390a450505050565b60008060006145f5866108cf565b925092509250614606868383613ec3565b8460011480156146275750876001600160a01b0316836001600160a01b0316145b6146865760405162461bcd60e51b815260206004820152602a60248201527f455243313135353a20696e73756666696369656e742062616c616e636520666f60448201526939103a3930b739b332b960b11b606482015260840161088a565b866001600160a01b0316836001600160a01b0316036146a7575050506114d7565b60008681526001602052604090206001600160a01b03881663ffffffff60a01b60a085901b16176001600160c01b031960c084901b1617905560408051878152602081018790526001600160a01b03808a1692908b169133917fc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62910160405180910390a4611cbe3389898989896149b4565b6000808351831061478c5760405162461bcd60e51b815260206004820152601e60248201527f726561644c6162656c3a20496e646578206f7574206f6620626f756e64730000604482015260640161088a565b60008484815181106147a0576147a06159bd565b016020015160f81c905080156147cc576147c5856147bf866001615e4d565b83614ab0565b92506147d1565b600092505b6147db8185615e4d565b6147e6906001615e4d565b9150509250929050565b600081815260036020526040902080546001600160a01b0319166001600160a01b0384169081179091558190614825826119aa565b6001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92560405160405180910390a45050565b606060018351101561489c576040517f280dacb600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60ff835111156148da57826040517fe3ba295f00000000000000000000000000000000000000000000000000000000815260040161088a9190614e31565b825183836040516020016148f093929190615e60565b604051602081830303815290604052905092915050565b61ffff81161580159061491f57506201000181811614155b15613bf85760405163a2a7201360e01b81526004810183905260240161088a565b61494a8483614907565b6000848152600160205260409020546001600160a01b038116156149a85761497185614525565b6040516000815285907fee2ba1195c65bcf218a83d874335c6bf9d9067b4c672f3c3bf16cf40de7586c49060200160405180910390a25b6114d785858585614acf565b6001600160a01b0384163b15613bc15760405163f23a6e6160e01b81526001600160a01b0385169063f23a6e61906149f89089908990889088908890600401615ec1565b6020604051808303816000875af1925050508015614a33575060408051601f3d908101601f19168201909252614a3091810190615d2f565b60015b614a3f5761403b615d4c565b6001600160e01b0319811663f23a6e6160e01b14613df65760405162461bcd60e51b815260206004820152602860248201527f455243313135353a204552433131353552656365697665722072656a656374656044820152676420746f6b656e7360c01b606482015260840161088a565b6000614ac584614ac08486615e4d565b614d0c565b5091016020012090565b8360008080614add846108cf565b9194509250905063ffff0000821667ffffffffffffffff8087169083161115614b04578195505b428267ffffffffffffffff1610614b1a57958617955b6001600160a01b03841615614b715760405162461bcd60e51b815260206004820152601f60248201527f455243313135353a206d696e74206f66206578697374696e6720746f6b656e00604482015260640161088a565b6001600160a01b038816614bed5760405162461bcd60e51b815260206004820152602160248201527f455243313135353a206d696e7420746f20746865207a65726f2061646472657360448201527f7300000000000000000000000000000000000000000000000000000000000000606482015260840161088a565b306001600160a01b03891603614c6b5760405162461bcd60e51b815260206004820152603460248201527f455243313135353a206e65774f776e65722063616e6e6f74206265207468652060448201527f4e616d655772617070657220636f6e7472616374000000000000000000000000606482015260840161088a565b60008581526001602052604090206001600160a01b03891663ffffffff60a01b60a08a901b16176001600160c01b031960c089901b1617905560408051868152600160208201526001600160a01b038a169160009133917fc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62910160405180910390a4612db93360008a886001604051806020016040528060008152506149b4565b8151811115613bf85781516040517f8a3c1cfb00000000000000000000000000000000000000000000000000000000815261088a918391600401918252602082015260400190565b6001600160a01b038116811461259357600080fd5b60008060408385031215614d7c57600080fd5b8235614d8781614d54565b946020939093013593505050565b600060208284031215614da757600080fd5b5035919050565b6001600160e01b03198116811461259357600080fd5b600060208284031215614dd657600080fd5b813561099681614dae565b60005b83811015614dfc578181015183820152602001614de4565b50506000910152565b60008151808452614e1d816020860160208601614de1565b601f01601f19169290920160200192915050565b6020815260006109966020830184614e05565b60008060408385031215614e5757600080fd5b823591506020830135614e6981614d54565b809150509250929050565b803567ffffffffffffffff81168114614e8c57600080fd5b919050565b60008060408385031215614ea457600080fd5b82359150614eb460208401614e74565b90509250929050565b60008083601f840112614ecf57600080fd5b50813567ffffffffffffffff811115614ee757600080fd5b602083019150836020828501011115614eff57600080fd5b9250929050565b600080600080600060808688031215614f1e57600080fd5b8535614f2981614d54565b94506020860135614f3981614d54565b935060408601359250606086013567ffffffffffffffff811115614f5c57600080fd5b614f6888828901614ebd565b969995985093965092949392505050565b600060208284031215614f8b57600080fd5b813561099681614d54565b634e487b7160e01b600052604160045260246000fd5b601f8201601f1916810167ffffffffffffffff81118282101715614fd257614fd2614f96565b6040525050565b600067ffffffffffffffff821115614ff357614ff3614f96565b50601f01601f191660200190565b600082601f83011261501257600080fd5b813561501d81614fd9565b60405161502a8282614fac565b82815285602084870101111561503f57600080fd5b82602086016020830137600092810160200192909252509392505050565b803563ffffffff81168114614e8c57600080fd5b600080600080600080600060e0888a03121561508c57600080fd5b87359650602088013567ffffffffffffffff8111156150aa57600080fd5b6150b68a828b01615001565b96505060408801356150c781614d54565b945060608801356150d781614d54565b93506150e560808901614e74565b92506150f360a0890161505d565b915061510160c08901614e74565b905092959891949750929550565b600067ffffffffffffffff82111561512957615129614f96565b5060051b60200190565b600082601f83011261514457600080fd5b813560206151518261510f565b60405161515e8282614fac565b83815260059390931b850182019282810191508684111561517e57600080fd5b8286015b8481101561230e5780358352918301918301615182565b600080600080600060a086880312156151b157600080fd5b85356151bc81614d54565b945060208601356151cc81614d54565b9350604086013567ffffffffffffffff808211156151e957600080fd5b6151f589838a01615133565b9450606088013591508082111561520b57600080fd5b61521789838a01615133565b9350608088013591508082111561522d57600080fd5b5061523a88828901615001565b9150509295509295909350565b6000806000806080858703121561525d57600080fd5b84359350602085013592506152746040860161505d565b915061528260608601614e74565b905092959194509250565b803561ffff81168114614e8c57600080fd5b600080604083850312156152b257600080fd5b82359150614eb46020840161528d565b600080604083850312156152d557600080fd5b823567ffffffffffffffff808211156152ed57600080fd5b818501915085601f83011261530157600080fd5b8135602061530e8261510f565b60405161531b8282614fac565b83815260059390931b850182019282810191508984111561533b57600080fd5b948201945b8386101561536257853561535381614d54565b82529482019490820190615340565b9650508601359250508082111561537857600080fd5b506147e685828601615133565b600081518084526020808501945080840160005b838110156153b557815187529582019590820190600101615399565b509495945050505050565b6020815260006109966020830184615385565b6000806000606084860312156153e857600080fd5b83356153f381614d54565b9250602084013561540381614d54565b929592945050506040919091013590565b60008060006060848603121561542957600080fd5b833592506020840135915061544060408501614e74565b90509250925092565b60008060006060848603121561545e57600080fd5b83359250602084013561547081614d54565b9150604084013561548081614d54565b809150509250925092565b6000806000806000608086880312156154a357600080fd5b853567ffffffffffffffff8111156154ba57600080fd5b6154c688828901614ebd565b90965094505060208601356154da81614d54565b92506154e86040870161528d565b915060608601356154f881614d54565b809150509295509295909350565b801515811461259357600080fd5b6000806040838503121561552757600080fd5b823561553281614d54565b91506020830135614e6981615506565b60008060008060008060a0878903121561555b57600080fd5b863567ffffffffffffffff81111561557257600080fd5b61557e89828a01614ebd565b909750955050602087013561559281614d54565b93506040870135925060608701356155a981614d54565b91506155b76080880161528d565b90509295509295509295565b600080604083850312156155d657600080fd5b82359150614eb46020840161505d565b600080604083850312156155f957600080fd5b50508035926020909101359150565b60008060008060008060a0878903121561562157600080fd5b86359550602087013567ffffffffffffffff81111561563f57600080fd5b61564b89828a01614ebd565b909650945050604087013561565f81614d54565b925061566d6060880161505d565b91506155b760808801614e74565b6000806000806040858703121561569157600080fd5b843567ffffffffffffffff808211156156a957600080fd5b6156b588838901614ebd565b909650945060208701359150808211156156ce57600080fd5b506156db87828801614ebd565b95989497509550505050565b600080600080608085870312156156fd57600080fd5b84359350602085013561570f81614d54565b9250604085013561527481614d54565b60008060006060848603121561573457600080fd5b8335925060208401359150604084013561548081614d54565b6000806040838503121561576057600080fd5b823561576b81614d54565b91506020830135614e6981614d54565b6000806000806060858703121561579157600080fd5b843567ffffffffffffffff8111156157a857600080fd5b6157b487828801614ebd565b90955093505060208501356157c881614d54565b915060408501356157d881614d54565b939692955090935050565b600080600080600060a086880312156157fb57600080fd5b853561580681614d54565b9450602086013561581681614d54565b93506040860135925060608601359150608086013567ffffffffffffffff81111561584057600080fd5b61523a88828901615001565b60006020828403121561585e57600080fd5b815167ffffffffffffffff81111561587557600080fd5b8201601f8101841361588657600080fd5b805161589181614fd9565b60405161589e8282614fac565b8281528660208486010111156158b357600080fd5b610a73836020830160208701614de1565b600080600080608085870312156158da57600080fd5b843567ffffffffffffffff8111156158f157600080fd5b6158fd87828801615001565b945050602085013561590e81614d54565b925061591c6040860161528d565b915060608501356157d881614d54565b60006020828403121561593e57600080fd5b5051919050565b634e487b7160e01b600052601160045260246000fd5b67ffffffffffffffff81811683821601908082111561597c5761597c615945565b5092915050565b600181811c9082168061599757607f821691505b6020821081036159b757634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052603260045260246000fd5b6000600182016159e5576159e5615945565b5060010190565b6040815260006159ff6040830185615385565b8281036020840152615a118185615385565b95945050505050565b600060208284031215615a2c57600080fd5b815161099681614d54565b600060208284031215615a4957600080fd5b815161099681615506565b8183823760009101908152919050565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b60c081526000615aa160c083018a8c615a64565b6001600160a01b03898116602085015263ffffffff8916604085015267ffffffffffffffff881660608501528616608084015282810360a0840152615ae7818587615a64565b9b9a5050505050505050505050565b601f8211156109de57600081815260208120601f850160051c81016020861015615b1d5750805b601f850160051c820191505b81811015613bc157828155600101615b29565b67ffffffffffffffff831115615b5457615b54614f96565b615b6883615b628354615983565b83615af6565b6000601f841160018114615b9c5760008515615b845750838201355b600019600387901b1c1916600186901b1783556114d7565b600083815260209020601f19861690835b82811015615bcd5786850135825560209485019460019092019101615bad565b5086821015615bea5760001960f88860031b161c19848701351681555b505060018560011b0183555050505050565b67ffffffffffffffff82811682821603908082111561597c5761597c615945565b815167ffffffffffffffff811115615c3757615c37614f96565b615c4b81615c458454615983565b84615af6565b602080601f831160018114615c805760008415615c685750858301515b600019600386901b1c1916600185901b178555613bc1565b600085815260208120601f198616915b82811015615caf57888601518255948401946001909101908401615c90565b5085821015615ccd5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b60006001600160a01b03808816835280871660208401525060a06040830152615d0960a0830186615385565b8281036060840152615d1b8186615385565b90508281036080840152613d038185614e05565b600060208284031215615d4157600080fd5b815161099681614dae565b600060033d1115615d655760046000803e5060005160e01c5b90565b600060443d1015615d765790565b6040516003193d81016004833e81513d67ffffffffffffffff8160248401118184111715615da657505050505090565b8285019150815181811115615dbe5750505050505090565b843d8701016020828501011115615dd85750505050505090565b615de760208286010187614fac565b509095945050505050565b608081526000615e056080830187614e05565b6001600160a01b039590951660208301525063ffffffff92909216604083015267ffffffffffffffff16606090910152919050565b818103818111156108c9576108c9615945565b808201808211156108c9576108c9615945565b7fff000000000000000000000000000000000000000000000000000000000000008460f81b16815260008351615e9d816001850160208801614de1565b835190830190615eb4816001840160208801614de1565b0160010195945050505050565b60006001600160a01b03808816835280871660208401525084604083015283606083015260a06080830152615ef960a0830184614e05565b97965050505050505056fea2646970667358221220bea32fee24824cf7ac7438e28659c91906f7454ecc272004d70d9057099cb26f64736f6c63430008110033",
}

// NameWrapperABI is the input ABI used to generate the binding from.
// Deprecated: Use NameWrapperMetaData.ABI instead.
var NameWrapperABI = NameWrapperMetaData.ABI

// NameWrapperBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use NameWrapperMetaData.Bin instead.
var NameWrapperBin = NameWrapperMetaData.Bin

// DeployNameWrapper deploys a new Ethereum contract, binding an instance of NameWrapper to it.
func DeployNameWrapper(auth *bind.TransactOpts, backend bind.ContractBackend, _ens common.Address, _registrar common.Address, _metadataService common.Address) (common.Address, *types.Transaction, *NameWrapper, error) {
	parsed, err := NameWrapperMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(NameWrapperBin), backend, _ens, _registrar, _metadataService)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &NameWrapper{NameWrapperCaller: NameWrapperCaller{contract: contract}, NameWrapperTransactor: NameWrapperTransactor{contract: contract}, NameWrapperFilterer: NameWrapperFilterer{contract: contract}}, nil
}

// NameWrapper is an auto generated Go binding around an Ethereum contract.
type NameWrapper struct {
	NameWrapperCaller     // Read-only binding to the contract
	NameWrapperTransactor // Write-only binding to the contract
	NameWrapperFilterer   // Log filterer for contract events
}

// NameWrapperCaller is an auto generated read-only Go binding around an Ethereum contract.
type NameWrapperCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NameWrapperTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NameWrapperTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NameWrapperFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NameWrapperFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NameWrapperSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NameWrapperSession struct {
	Contract     *NameWrapper      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NameWrapperCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NameWrapperCallerSession struct {
	Contract *NameWrapperCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// NameWrapperTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NameWrapperTransactorSession struct {
	Contract     *NameWrapperTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// NameWrapperRaw is an auto generated low-level Go binding around an Ethereum contract.
type NameWrapperRaw struct {
	Contract *NameWrapper // Generic contract binding to access the raw methods on
}

// NameWrapperCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NameWrapperCallerRaw struct {
	Contract *NameWrapperCaller // Generic read-only contract binding to access the raw methods on
}

// NameWrapperTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NameWrapperTransactorRaw struct {
	Contract *NameWrapperTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNameWrapper creates a new instance of NameWrapper, bound to a specific deployed contract.
func NewNameWrapper(address common.Address, backend bind.ContractBackend) (*NameWrapper, error) {
	contract, err := bindNameWrapper(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NameWrapper{NameWrapperCaller: NameWrapperCaller{contract: contract}, NameWrapperTransactor: NameWrapperTransactor{contract: contract}, NameWrapperFilterer: NameWrapperFilterer{contract: contract}}, nil
}

// NewNameWrapperCaller creates a new read-only instance of NameWrapper, bound to a specific deployed contract.
func NewNameWrapperCaller(address common.Address, caller bind.ContractCaller) (*NameWrapperCaller, error) {
	contract, err := bindNameWrapper(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NameWrapperCaller{contract: contract}, nil
}

// NewNameWrapperTransactor creates a new write-only instance of NameWrapper, bound to a specific deployed contract.
func NewNameWrapperTransactor(address common.Address, transactor bind.ContractTransactor) (*NameWrapperTransactor, error) {
	contract, err := bindNameWrapper(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NameWrapperTransactor{contract: contract}, nil
}

// NewNameWrapperFilterer creates a new log filterer instance of NameWrapper, bound to a specific deployed contract.
func NewNameWrapperFilterer(address common.Address, filterer bind.ContractFilterer) (*NameWrapperFilterer, error) {
	contract, err := bindNameWrapper(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NameWrapperFilterer{contract: contract}, nil
}

// bindNameWrapper binds a generic wrapper to an already deployed contract.
func bindNameWrapper(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := NameWrapperMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NameWrapper *NameWrapperRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NameWrapper.Contract.NameWrapperCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NameWrapper *NameWrapperRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NameWrapper.Contract.NameWrapperTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NameWrapper *NameWrapperRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NameWrapper.Contract.NameWrapperTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NameWrapper *NameWrapperCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NameWrapper.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NameWrapper *NameWrapperTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NameWrapper.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NameWrapper *NameWrapperTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NameWrapper.Contract.contract.Transact(opts, method, params...)
}

// Tokens is a free data retrieval call binding the contract method 0xed70554d.
//
// Solidity: function _tokens(uint256 ) view returns(uint256)
func (_NameWrapper *NameWrapperCaller) Tokens(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _NameWrapper.contract.Call(opts, &out, "_tokens", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Tokens is a free data retrieval call binding the contract method 0xed70554d.
//
// Solidity: function _tokens(uint256 ) view returns(uint256)
func (_NameWrapper *NameWrapperSession) Tokens(arg0 *big.Int) (*big.Int, error) {
	return _NameWrapper.Contract.Tokens(&_NameWrapper.CallOpts, arg0)
}

// Tokens is a free data retrieval call binding the contract method 0xed70554d.
//
// Solidity: function _tokens(uint256 ) view returns(uint256)
func (_NameWrapper *NameWrapperCallerSession) Tokens(arg0 *big.Int) (*big.Int, error) {
	return _NameWrapper.Contract.Tokens(&_NameWrapper.CallOpts, arg0)
}

// AllFusesBurned is a free data retrieval call binding the contract method 0xadf4960a.
//
// Solidity: function allFusesBurned(bytes32 node, uint32 fuseMask) view returns(bool)
func (_NameWrapper *NameWrapperCaller) AllFusesBurned(opts *bind.CallOpts, node [32]byte, fuseMask uint32) (bool, error) {
	var out []interface{}
	err := _NameWrapper.contract.Call(opts, &out, "allFusesBurned", node, fuseMask)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AllFusesBurned is a free data retrieval call binding the contract method 0xadf4960a.
//
// Solidity: function allFusesBurned(bytes32 node, uint32 fuseMask) view returns(bool)
func (_NameWrapper *NameWrapperSession) AllFusesBurned(node [32]byte, fuseMask uint32) (bool, error) {
	return _NameWrapper.Contract.AllFusesBurned(&_NameWrapper.CallOpts, node, fuseMask)
}

// AllFusesBurned is a free data retrieval call binding the contract method 0xadf4960a.
//
// Solidity: function allFusesBurned(bytes32 node, uint32 fuseMask) view returns(bool)
func (_NameWrapper *NameWrapperCallerSession) AllFusesBurned(node [32]byte, fuseMask uint32) (bool, error) {
	return _NameWrapper.Contract.AllFusesBurned(&_NameWrapper.CallOpts, node, fuseMask)
}

// BalanceOf is a free data retrieval call binding the contract method 0x00fdd58e.
//
// Solidity: function balanceOf(address account, uint256 id) view returns(uint256)
func (_NameWrapper *NameWrapperCaller) BalanceOf(opts *bind.CallOpts, account common.Address, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _NameWrapper.contract.Call(opts, &out, "balanceOf", account, id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x00fdd58e.
//
// Solidity: function balanceOf(address account, uint256 id) view returns(uint256)
func (_NameWrapper *NameWrapperSession) BalanceOf(account common.Address, id *big.Int) (*big.Int, error) {
	return _NameWrapper.Contract.BalanceOf(&_NameWrapper.CallOpts, account, id)
}

// BalanceOf is a free data retrieval call binding the contract method 0x00fdd58e.
//
// Solidity: function balanceOf(address account, uint256 id) view returns(uint256)
func (_NameWrapper *NameWrapperCallerSession) BalanceOf(account common.Address, id *big.Int) (*big.Int, error) {
	return _NameWrapper.Contract.BalanceOf(&_NameWrapper.CallOpts, account, id)
}

// BalanceOfBatch is a free data retrieval call binding the contract method 0x4e1273f4.
//
// Solidity: function balanceOfBatch(address[] accounts, uint256[] ids) view returns(uint256[])
func (_NameWrapper *NameWrapperCaller) BalanceOfBatch(opts *bind.CallOpts, accounts []common.Address, ids []*big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _NameWrapper.contract.Call(opts, &out, "balanceOfBatch", accounts, ids)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// BalanceOfBatch is a free data retrieval call binding the contract method 0x4e1273f4.
//
// Solidity: function balanceOfBatch(address[] accounts, uint256[] ids) view returns(uint256[])
func (_NameWrapper *NameWrapperSession) BalanceOfBatch(accounts []common.Address, ids []*big.Int) ([]*big.Int, error) {
	return _NameWrapper.Contract.BalanceOfBatch(&_NameWrapper.CallOpts, accounts, ids)
}

// BalanceOfBatch is a free data retrieval call binding the contract method 0x4e1273f4.
//
// Solidity: function balanceOfBatch(address[] accounts, uint256[] ids) view returns(uint256[])
func (_NameWrapper *NameWrapperCallerSession) BalanceOfBatch(accounts []common.Address, ids []*big.Int) ([]*big.Int, error) {
	return _NameWrapper.Contract.BalanceOfBatch(&_NameWrapper.CallOpts, accounts, ids)
}

// CanExtendSubnames is a free data retrieval call binding the contract method 0x0e4cd725.
//
// Solidity: function canExtendSubnames(bytes32 node, address addr) view returns(bool)
func (_NameWrapper *NameWrapperCaller) CanExtendSubnames(opts *bind.CallOpts, node [32]byte, addr common.Address) (bool, error) {
	var out []interface{}
	err := _NameWrapper.contract.Call(opts, &out, "canExtendSubnames", node, addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CanExtendSubnames is a free data retrieval call binding the contract method 0x0e4cd725.
//
// Solidity: function canExtendSubnames(bytes32 node, address addr) view returns(bool)
func (_NameWrapper *NameWrapperSession) CanExtendSubnames(node [32]byte, addr common.Address) (bool, error) {
	return _NameWrapper.Contract.CanExtendSubnames(&_NameWrapper.CallOpts, node, addr)
}

// CanExtendSubnames is a free data retrieval call binding the contract method 0x0e4cd725.
//
// Solidity: function canExtendSubnames(bytes32 node, address addr) view returns(bool)
func (_NameWrapper *NameWrapperCallerSession) CanExtendSubnames(node [32]byte, addr common.Address) (bool, error) {
	return _NameWrapper.Contract.CanExtendSubnames(&_NameWrapper.CallOpts, node, addr)
}

// CanModifyName is a free data retrieval call binding the contract method 0x41415eab.
//
// Solidity: function canModifyName(bytes32 node, address addr) view returns(bool)
func (_NameWrapper *NameWrapperCaller) CanModifyName(opts *bind.CallOpts, node [32]byte, addr common.Address) (bool, error) {
	var out []interface{}
	err := _NameWrapper.contract.Call(opts, &out, "canModifyName", node, addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CanModifyName is a free data retrieval call binding the contract method 0x41415eab.
//
// Solidity: function canModifyName(bytes32 node, address addr) view returns(bool)
func (_NameWrapper *NameWrapperSession) CanModifyName(node [32]byte, addr common.Address) (bool, error) {
	return _NameWrapper.Contract.CanModifyName(&_NameWrapper.CallOpts, node, addr)
}

// CanModifyName is a free data retrieval call binding the contract method 0x41415eab.
//
// Solidity: function canModifyName(bytes32 node, address addr) view returns(bool)
func (_NameWrapper *NameWrapperCallerSession) CanModifyName(node [32]byte, addr common.Address) (bool, error) {
	return _NameWrapper.Contract.CanModifyName(&_NameWrapper.CallOpts, node, addr)
}

// Controllers is a free data retrieval call binding the contract method 0xda8c229e.
//
// Solidity: function controllers(address ) view returns(bool)
func (_NameWrapper *NameWrapperCaller) Controllers(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _NameWrapper.contract.Call(opts, &out, "controllers", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Controllers is a free data retrieval call binding the contract method 0xda8c229e.
//
// Solidity: function controllers(address ) view returns(bool)
func (_NameWrapper *NameWrapperSession) Controllers(arg0 common.Address) (bool, error) {
	return _NameWrapper.Contract.Controllers(&_NameWrapper.CallOpts, arg0)
}

// Controllers is a free data retrieval call binding the contract method 0xda8c229e.
//
// Solidity: function controllers(address ) view returns(bool)
func (_NameWrapper *NameWrapperCallerSession) Controllers(arg0 common.Address) (bool, error) {
	return _NameWrapper.Contract.Controllers(&_NameWrapper.CallOpts, arg0)
}

// Ens is a free data retrieval call binding the contract method 0x3f15457f.
//
// Solidity: function ens() view returns(address)
func (_NameWrapper *NameWrapperCaller) Ens(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NameWrapper.contract.Call(opts, &out, "ens")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Ens is a free data retrieval call binding the contract method 0x3f15457f.
//
// Solidity: function ens() view returns(address)
func (_NameWrapper *NameWrapperSession) Ens() (common.Address, error) {
	return _NameWrapper.Contract.Ens(&_NameWrapper.CallOpts)
}

// Ens is a free data retrieval call binding the contract method 0x3f15457f.
//
// Solidity: function ens() view returns(address)
func (_NameWrapper *NameWrapperCallerSession) Ens() (common.Address, error) {
	return _NameWrapper.Contract.Ens(&_NameWrapper.CallOpts)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 id) view returns(address operator)
func (_NameWrapper *NameWrapperCaller) GetApproved(opts *bind.CallOpts, id *big.Int) (common.Address, error) {
	var out []interface{}
	err := _NameWrapper.contract.Call(opts, &out, "getApproved", id)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 id) view returns(address operator)
func (_NameWrapper *NameWrapperSession) GetApproved(id *big.Int) (common.Address, error) {
	return _NameWrapper.Contract.GetApproved(&_NameWrapper.CallOpts, id)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 id) view returns(address operator)
func (_NameWrapper *NameWrapperCallerSession) GetApproved(id *big.Int) (common.Address, error) {
	return _NameWrapper.Contract.GetApproved(&_NameWrapper.CallOpts, id)
}

// GetData is a free data retrieval call binding the contract method 0x0178fe3f.
//
// Solidity: function getData(uint256 id) view returns(address owner, uint32 fuses, uint64 expiry)
func (_NameWrapper *NameWrapperCaller) GetData(opts *bind.CallOpts, id *big.Int) (struct {
	Owner  common.Address
	Fuses  uint32
	Expiry uint64
}, error) {
	var out []interface{}
	err := _NameWrapper.contract.Call(opts, &out, "getData", id)

	outstruct := new(struct {
		Owner  common.Address
		Fuses  uint32
		Expiry uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Owner = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Fuses = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.Expiry = *abi.ConvertType(out[2], new(uint64)).(*uint64)

	return *outstruct, err

}

// GetData is a free data retrieval call binding the contract method 0x0178fe3f.
//
// Solidity: function getData(uint256 id) view returns(address owner, uint32 fuses, uint64 expiry)
func (_NameWrapper *NameWrapperSession) GetData(id *big.Int) (struct {
	Owner  common.Address
	Fuses  uint32
	Expiry uint64
}, error) {
	return _NameWrapper.Contract.GetData(&_NameWrapper.CallOpts, id)
}

// GetData is a free data retrieval call binding the contract method 0x0178fe3f.
//
// Solidity: function getData(uint256 id) view returns(address owner, uint32 fuses, uint64 expiry)
func (_NameWrapper *NameWrapperCallerSession) GetData(id *big.Int) (struct {
	Owner  common.Address
	Fuses  uint32
	Expiry uint64
}, error) {
	return _NameWrapper.Contract.GetData(&_NameWrapper.CallOpts, id)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address account, address operator) view returns(bool)
func (_NameWrapper *NameWrapperCaller) IsApprovedForAll(opts *bind.CallOpts, account common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _NameWrapper.contract.Call(opts, &out, "isApprovedForAll", account, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address account, address operator) view returns(bool)
func (_NameWrapper *NameWrapperSession) IsApprovedForAll(account common.Address, operator common.Address) (bool, error) {
	return _NameWrapper.Contract.IsApprovedForAll(&_NameWrapper.CallOpts, account, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address account, address operator) view returns(bool)
func (_NameWrapper *NameWrapperCallerSession) IsApprovedForAll(account common.Address, operator common.Address) (bool, error) {
	return _NameWrapper.Contract.IsApprovedForAll(&_NameWrapper.CallOpts, account, operator)
}

// IsWrapped is a free data retrieval call binding the contract method 0xd9a50c12.
//
// Solidity: function isWrapped(bytes32 parentNode, bytes32 labelhash) view returns(bool)
func (_NameWrapper *NameWrapperCaller) IsWrapped(opts *bind.CallOpts, parentNode [32]byte, labelhash [32]byte) (bool, error) {
	var out []interface{}
	err := _NameWrapper.contract.Call(opts, &out, "isWrapped", parentNode, labelhash)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsWrapped is a free data retrieval call binding the contract method 0xd9a50c12.
//
// Solidity: function isWrapped(bytes32 parentNode, bytes32 labelhash) view returns(bool)
func (_NameWrapper *NameWrapperSession) IsWrapped(parentNode [32]byte, labelhash [32]byte) (bool, error) {
	return _NameWrapper.Contract.IsWrapped(&_NameWrapper.CallOpts, parentNode, labelhash)
}

// IsWrapped is a free data retrieval call binding the contract method 0xd9a50c12.
//
// Solidity: function isWrapped(bytes32 parentNode, bytes32 labelhash) view returns(bool)
func (_NameWrapper *NameWrapperCallerSession) IsWrapped(parentNode [32]byte, labelhash [32]byte) (bool, error) {
	return _NameWrapper.Contract.IsWrapped(&_NameWrapper.CallOpts, parentNode, labelhash)
}

// IsWrapped0 is a free data retrieval call binding the contract method 0xfd0cd0d9.
//
// Solidity: function isWrapped(bytes32 node) view returns(bool)
func (_NameWrapper *NameWrapperCaller) IsWrapped0(opts *bind.CallOpts, node [32]byte) (bool, error) {
	var out []interface{}
	err := _NameWrapper.contract.Call(opts, &out, "isWrapped0", node)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsWrapped0 is a free data retrieval call binding the contract method 0xfd0cd0d9.
//
// Solidity: function isWrapped(bytes32 node) view returns(bool)
func (_NameWrapper *NameWrapperSession) IsWrapped0(node [32]byte) (bool, error) {
	return _NameWrapper.Contract.IsWrapped0(&_NameWrapper.CallOpts, node)
}

// IsWrapped0 is a free data retrieval call binding the contract method 0xfd0cd0d9.
//
// Solidity: function isWrapped(bytes32 node) view returns(bool)
func (_NameWrapper *NameWrapperCallerSession) IsWrapped0(node [32]byte) (bool, error) {
	return _NameWrapper.Contract.IsWrapped0(&_NameWrapper.CallOpts, node)
}

// MetadataService is a free data retrieval call binding the contract method 0x53095467.
//
// Solidity: function metadataService() view returns(address)
func (_NameWrapper *NameWrapperCaller) MetadataService(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NameWrapper.contract.Call(opts, &out, "metadataService")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MetadataService is a free data retrieval call binding the contract method 0x53095467.
//
// Solidity: function metadataService() view returns(address)
func (_NameWrapper *NameWrapperSession) MetadataService() (common.Address, error) {
	return _NameWrapper.Contract.MetadataService(&_NameWrapper.CallOpts)
}

// MetadataService is a free data retrieval call binding the contract method 0x53095467.
//
// Solidity: function metadataService() view returns(address)
func (_NameWrapper *NameWrapperCallerSession) MetadataService() (common.Address, error) {
	return _NameWrapper.Contract.MetadataService(&_NameWrapper.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_NameWrapper *NameWrapperCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _NameWrapper.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_NameWrapper *NameWrapperSession) Name() (string, error) {
	return _NameWrapper.Contract.Name(&_NameWrapper.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_NameWrapper *NameWrapperCallerSession) Name() (string, error) {
	return _NameWrapper.Contract.Name(&_NameWrapper.CallOpts)
}

// Names is a free data retrieval call binding the contract method 0x20c38e2b.
//
// Solidity: function names(bytes32 ) view returns(bytes)
func (_NameWrapper *NameWrapperCaller) Names(opts *bind.CallOpts, arg0 [32]byte) ([]byte, error) {
	var out []interface{}
	err := _NameWrapper.contract.Call(opts, &out, "names", arg0)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Names is a free data retrieval call binding the contract method 0x20c38e2b.
//
// Solidity: function names(bytes32 ) view returns(bytes)
func (_NameWrapper *NameWrapperSession) Names(arg0 [32]byte) ([]byte, error) {
	return _NameWrapper.Contract.Names(&_NameWrapper.CallOpts, arg0)
}

// Names is a free data retrieval call binding the contract method 0x20c38e2b.
//
// Solidity: function names(bytes32 ) view returns(bytes)
func (_NameWrapper *NameWrapperCallerSession) Names(arg0 [32]byte) ([]byte, error) {
	return _NameWrapper.Contract.Names(&_NameWrapper.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_NameWrapper *NameWrapperCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NameWrapper.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_NameWrapper *NameWrapperSession) Owner() (common.Address, error) {
	return _NameWrapper.Contract.Owner(&_NameWrapper.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_NameWrapper *NameWrapperCallerSession) Owner() (common.Address, error) {
	return _NameWrapper.Contract.Owner(&_NameWrapper.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 id) view returns(address owner)
func (_NameWrapper *NameWrapperCaller) OwnerOf(opts *bind.CallOpts, id *big.Int) (common.Address, error) {
	var out []interface{}
	err := _NameWrapper.contract.Call(opts, &out, "ownerOf", id)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 id) view returns(address owner)
func (_NameWrapper *NameWrapperSession) OwnerOf(id *big.Int) (common.Address, error) {
	return _NameWrapper.Contract.OwnerOf(&_NameWrapper.CallOpts, id)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 id) view returns(address owner)
func (_NameWrapper *NameWrapperCallerSession) OwnerOf(id *big.Int) (common.Address, error) {
	return _NameWrapper.Contract.OwnerOf(&_NameWrapper.CallOpts, id)
}

// Registrar is a free data retrieval call binding the contract method 0x2b20e397.
//
// Solidity: function registrar() view returns(address)
func (_NameWrapper *NameWrapperCaller) Registrar(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NameWrapper.contract.Call(opts, &out, "registrar")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Registrar is a free data retrieval call binding the contract method 0x2b20e397.
//
// Solidity: function registrar() view returns(address)
func (_NameWrapper *NameWrapperSession) Registrar() (common.Address, error) {
	return _NameWrapper.Contract.Registrar(&_NameWrapper.CallOpts)
}

// Registrar is a free data retrieval call binding the contract method 0x2b20e397.
//
// Solidity: function registrar() view returns(address)
func (_NameWrapper *NameWrapperCallerSession) Registrar() (common.Address, error) {
	return _NameWrapper.Contract.Registrar(&_NameWrapper.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_NameWrapper *NameWrapperCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _NameWrapper.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_NameWrapper *NameWrapperSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _NameWrapper.Contract.SupportsInterface(&_NameWrapper.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_NameWrapper *NameWrapperCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _NameWrapper.Contract.SupportsInterface(&_NameWrapper.CallOpts, interfaceId)
}

// UpgradeContract is a free data retrieval call binding the contract method 0x1f4e1504.
//
// Solidity: function upgradeContract() view returns(address)
func (_NameWrapper *NameWrapperCaller) UpgradeContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NameWrapper.contract.Call(opts, &out, "upgradeContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UpgradeContract is a free data retrieval call binding the contract method 0x1f4e1504.
//
// Solidity: function upgradeContract() view returns(address)
func (_NameWrapper *NameWrapperSession) UpgradeContract() (common.Address, error) {
	return _NameWrapper.Contract.UpgradeContract(&_NameWrapper.CallOpts)
}

// UpgradeContract is a free data retrieval call binding the contract method 0x1f4e1504.
//
// Solidity: function upgradeContract() view returns(address)
func (_NameWrapper *NameWrapperCallerSession) UpgradeContract() (common.Address, error) {
	return _NameWrapper.Contract.UpgradeContract(&_NameWrapper.CallOpts)
}

// Uri is a free data retrieval call binding the contract method 0x0e89341c.
//
// Solidity: function uri(uint256 tokenId) view returns(string)
func (_NameWrapper *NameWrapperCaller) Uri(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _NameWrapper.contract.Call(opts, &out, "uri", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Uri is a free data retrieval call binding the contract method 0x0e89341c.
//
// Solidity: function uri(uint256 tokenId) view returns(string)
func (_NameWrapper *NameWrapperSession) Uri(tokenId *big.Int) (string, error) {
	return _NameWrapper.Contract.Uri(&_NameWrapper.CallOpts, tokenId)
}

// Uri is a free data retrieval call binding the contract method 0x0e89341c.
//
// Solidity: function uri(uint256 tokenId) view returns(string)
func (_NameWrapper *NameWrapperCallerSession) Uri(tokenId *big.Int) (string, error) {
	return _NameWrapper.Contract.Uri(&_NameWrapper.CallOpts, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_NameWrapper *NameWrapperTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NameWrapper.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_NameWrapper *NameWrapperSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NameWrapper.Contract.Approve(&_NameWrapper.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_NameWrapper *NameWrapperTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NameWrapper.Contract.Approve(&_NameWrapper.TransactOpts, to, tokenId)
}

// ExtendExpiry is a paid mutator transaction binding the contract method 0x6e5d6ad2.
//
// Solidity: function extendExpiry(bytes32 parentNode, bytes32 labelhash, uint64 expiry) returns(uint64)
func (_NameWrapper *NameWrapperTransactor) ExtendExpiry(opts *bind.TransactOpts, parentNode [32]byte, labelhash [32]byte, expiry uint64) (*types.Transaction, error) {
	return _NameWrapper.contract.Transact(opts, "extendExpiry", parentNode, labelhash, expiry)
}

// ExtendExpiry is a paid mutator transaction binding the contract method 0x6e5d6ad2.
//
// Solidity: function extendExpiry(bytes32 parentNode, bytes32 labelhash, uint64 expiry) returns(uint64)
func (_NameWrapper *NameWrapperSession) ExtendExpiry(parentNode [32]byte, labelhash [32]byte, expiry uint64) (*types.Transaction, error) {
	return _NameWrapper.Contract.ExtendExpiry(&_NameWrapper.TransactOpts, parentNode, labelhash, expiry)
}

// ExtendExpiry is a paid mutator transaction binding the contract method 0x6e5d6ad2.
//
// Solidity: function extendExpiry(bytes32 parentNode, bytes32 labelhash, uint64 expiry) returns(uint64)
func (_NameWrapper *NameWrapperTransactorSession) ExtendExpiry(parentNode [32]byte, labelhash [32]byte, expiry uint64) (*types.Transaction, error) {
	return _NameWrapper.Contract.ExtendExpiry(&_NameWrapper.TransactOpts, parentNode, labelhash, expiry)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address to, address , uint256 tokenId, bytes data) returns(bytes4)
func (_NameWrapper *NameWrapperTransactor) OnERC721Received(opts *bind.TransactOpts, to common.Address, arg1 common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _NameWrapper.contract.Transact(opts, "onERC721Received", to, arg1, tokenId, data)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address to, address , uint256 tokenId, bytes data) returns(bytes4)
func (_NameWrapper *NameWrapperSession) OnERC721Received(to common.Address, arg1 common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _NameWrapper.Contract.OnERC721Received(&_NameWrapper.TransactOpts, to, arg1, tokenId, data)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address to, address , uint256 tokenId, bytes data) returns(bytes4)
func (_NameWrapper *NameWrapperTransactorSession) OnERC721Received(to common.Address, arg1 common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _NameWrapper.Contract.OnERC721Received(&_NameWrapper.TransactOpts, to, arg1, tokenId, data)
}

// RecoverFunds is a paid mutator transaction binding the contract method 0x5d3590d5.
//
// Solidity: function recoverFunds(address _token, address _to, uint256 _amount) returns()
func (_NameWrapper *NameWrapperTransactor) RecoverFunds(opts *bind.TransactOpts, _token common.Address, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _NameWrapper.contract.Transact(opts, "recoverFunds", _token, _to, _amount)
}

// RecoverFunds is a paid mutator transaction binding the contract method 0x5d3590d5.
//
// Solidity: function recoverFunds(address _token, address _to, uint256 _amount) returns()
func (_NameWrapper *NameWrapperSession) RecoverFunds(_token common.Address, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _NameWrapper.Contract.RecoverFunds(&_NameWrapper.TransactOpts, _token, _to, _amount)
}

// RecoverFunds is a paid mutator transaction binding the contract method 0x5d3590d5.
//
// Solidity: function recoverFunds(address _token, address _to, uint256 _amount) returns()
func (_NameWrapper *NameWrapperTransactorSession) RecoverFunds(_token common.Address, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _NameWrapper.Contract.RecoverFunds(&_NameWrapper.TransactOpts, _token, _to, _amount)
}

// RegisterAndWrapETH2LD is a paid mutator transaction binding the contract method 0xa4014982.
//
// Solidity: function registerAndWrapETH2LD(string label, address wrappedOwner, uint256 duration, address resolver, uint16 ownerControlledFuses) returns(uint256 registrarExpiry)
func (_NameWrapper *NameWrapperTransactor) RegisterAndWrapETH2LD(opts *bind.TransactOpts, label string, wrappedOwner common.Address, duration *big.Int, resolver common.Address, ownerControlledFuses uint16) (*types.Transaction, error) {
	return _NameWrapper.contract.Transact(opts, "registerAndWrapETH2LD", label, wrappedOwner, duration, resolver, ownerControlledFuses)
}

// RegisterAndWrapETH2LD is a paid mutator transaction binding the contract method 0xa4014982.
//
// Solidity: function registerAndWrapETH2LD(string label, address wrappedOwner, uint256 duration, address resolver, uint16 ownerControlledFuses) returns(uint256 registrarExpiry)
func (_NameWrapper *NameWrapperSession) RegisterAndWrapETH2LD(label string, wrappedOwner common.Address, duration *big.Int, resolver common.Address, ownerControlledFuses uint16) (*types.Transaction, error) {
	return _NameWrapper.Contract.RegisterAndWrapETH2LD(&_NameWrapper.TransactOpts, label, wrappedOwner, duration, resolver, ownerControlledFuses)
}

// RegisterAndWrapETH2LD is a paid mutator transaction binding the contract method 0xa4014982.
//
// Solidity: function registerAndWrapETH2LD(string label, address wrappedOwner, uint256 duration, address resolver, uint16 ownerControlledFuses) returns(uint256 registrarExpiry)
func (_NameWrapper *NameWrapperTransactorSession) RegisterAndWrapETH2LD(label string, wrappedOwner common.Address, duration *big.Int, resolver common.Address, ownerControlledFuses uint16) (*types.Transaction, error) {
	return _NameWrapper.Contract.RegisterAndWrapETH2LD(&_NameWrapper.TransactOpts, label, wrappedOwner, duration, resolver, ownerControlledFuses)
}

// Renew is a paid mutator transaction binding the contract method 0xc475abff.
//
// Solidity: function renew(uint256 tokenId, uint256 duration) returns(uint256 expires)
func (_NameWrapper *NameWrapperTransactor) Renew(opts *bind.TransactOpts, tokenId *big.Int, duration *big.Int) (*types.Transaction, error) {
	return _NameWrapper.contract.Transact(opts, "renew", tokenId, duration)
}

// Renew is a paid mutator transaction binding the contract method 0xc475abff.
//
// Solidity: function renew(uint256 tokenId, uint256 duration) returns(uint256 expires)
func (_NameWrapper *NameWrapperSession) Renew(tokenId *big.Int, duration *big.Int) (*types.Transaction, error) {
	return _NameWrapper.Contract.Renew(&_NameWrapper.TransactOpts, tokenId, duration)
}

// Renew is a paid mutator transaction binding the contract method 0xc475abff.
//
// Solidity: function renew(uint256 tokenId, uint256 duration) returns(uint256 expires)
func (_NameWrapper *NameWrapperTransactorSession) Renew(tokenId *big.Int, duration *big.Int) (*types.Transaction, error) {
	return _NameWrapper.Contract.Renew(&_NameWrapper.TransactOpts, tokenId, duration)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_NameWrapper *NameWrapperTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NameWrapper.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_NameWrapper *NameWrapperSession) RenounceOwnership() (*types.Transaction, error) {
	return _NameWrapper.Contract.RenounceOwnership(&_NameWrapper.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_NameWrapper *NameWrapperTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _NameWrapper.Contract.RenounceOwnership(&_NameWrapper.TransactOpts)
}

// SafeBatchTransferFrom is a paid mutator transaction binding the contract method 0x2eb2c2d6.
//
// Solidity: function safeBatchTransferFrom(address from, address to, uint256[] ids, uint256[] amounts, bytes data) returns()
func (_NameWrapper *NameWrapperTransactor) SafeBatchTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, ids []*big.Int, amounts []*big.Int, data []byte) (*types.Transaction, error) {
	return _NameWrapper.contract.Transact(opts, "safeBatchTransferFrom", from, to, ids, amounts, data)
}

// SafeBatchTransferFrom is a paid mutator transaction binding the contract method 0x2eb2c2d6.
//
// Solidity: function safeBatchTransferFrom(address from, address to, uint256[] ids, uint256[] amounts, bytes data) returns()
func (_NameWrapper *NameWrapperSession) SafeBatchTransferFrom(from common.Address, to common.Address, ids []*big.Int, amounts []*big.Int, data []byte) (*types.Transaction, error) {
	return _NameWrapper.Contract.SafeBatchTransferFrom(&_NameWrapper.TransactOpts, from, to, ids, amounts, data)
}

// SafeBatchTransferFrom is a paid mutator transaction binding the contract method 0x2eb2c2d6.
//
// Solidity: function safeBatchTransferFrom(address from, address to, uint256[] ids, uint256[] amounts, bytes data) returns()
func (_NameWrapper *NameWrapperTransactorSession) SafeBatchTransferFrom(from common.Address, to common.Address, ids []*big.Int, amounts []*big.Int, data []byte) (*types.Transaction, error) {
	return _NameWrapper.Contract.SafeBatchTransferFrom(&_NameWrapper.TransactOpts, from, to, ids, amounts, data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xf242432a.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 id, uint256 amount, bytes data) returns()
func (_NameWrapper *NameWrapperTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, id *big.Int, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _NameWrapper.contract.Transact(opts, "safeTransferFrom", from, to, id, amount, data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xf242432a.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 id, uint256 amount, bytes data) returns()
func (_NameWrapper *NameWrapperSession) SafeTransferFrom(from common.Address, to common.Address, id *big.Int, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _NameWrapper.Contract.SafeTransferFrom(&_NameWrapper.TransactOpts, from, to, id, amount, data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xf242432a.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 id, uint256 amount, bytes data) returns()
func (_NameWrapper *NameWrapperTransactorSession) SafeTransferFrom(from common.Address, to common.Address, id *big.Int, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _NameWrapper.Contract.SafeTransferFrom(&_NameWrapper.TransactOpts, from, to, id, amount, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_NameWrapper *NameWrapperTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _NameWrapper.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_NameWrapper *NameWrapperSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _NameWrapper.Contract.SetApprovalForAll(&_NameWrapper.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_NameWrapper *NameWrapperTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _NameWrapper.Contract.SetApprovalForAll(&_NameWrapper.TransactOpts, operator, approved)
}

// SetChildFuses is a paid mutator transaction binding the contract method 0x33c69ea9.
//
// Solidity: function setChildFuses(bytes32 parentNode, bytes32 labelhash, uint32 fuses, uint64 expiry) returns()
func (_NameWrapper *NameWrapperTransactor) SetChildFuses(opts *bind.TransactOpts, parentNode [32]byte, labelhash [32]byte, fuses uint32, expiry uint64) (*types.Transaction, error) {
	return _NameWrapper.contract.Transact(opts, "setChildFuses", parentNode, labelhash, fuses, expiry)
}

// SetChildFuses is a paid mutator transaction binding the contract method 0x33c69ea9.
//
// Solidity: function setChildFuses(bytes32 parentNode, bytes32 labelhash, uint32 fuses, uint64 expiry) returns()
func (_NameWrapper *NameWrapperSession) SetChildFuses(parentNode [32]byte, labelhash [32]byte, fuses uint32, expiry uint64) (*types.Transaction, error) {
	return _NameWrapper.Contract.SetChildFuses(&_NameWrapper.TransactOpts, parentNode, labelhash, fuses, expiry)
}

// SetChildFuses is a paid mutator transaction binding the contract method 0x33c69ea9.
//
// Solidity: function setChildFuses(bytes32 parentNode, bytes32 labelhash, uint32 fuses, uint64 expiry) returns()
func (_NameWrapper *NameWrapperTransactorSession) SetChildFuses(parentNode [32]byte, labelhash [32]byte, fuses uint32, expiry uint64) (*types.Transaction, error) {
	return _NameWrapper.Contract.SetChildFuses(&_NameWrapper.TransactOpts, parentNode, labelhash, fuses, expiry)
}

// SetController is a paid mutator transaction binding the contract method 0xe0dba60f.
//
// Solidity: function setController(address controller, bool active) returns()
func (_NameWrapper *NameWrapperTransactor) SetController(opts *bind.TransactOpts, controller common.Address, active bool) (*types.Transaction, error) {
	return _NameWrapper.contract.Transact(opts, "setController", controller, active)
}

// SetController is a paid mutator transaction binding the contract method 0xe0dba60f.
//
// Solidity: function setController(address controller, bool active) returns()
func (_NameWrapper *NameWrapperSession) SetController(controller common.Address, active bool) (*types.Transaction, error) {
	return _NameWrapper.Contract.SetController(&_NameWrapper.TransactOpts, controller, active)
}

// SetController is a paid mutator transaction binding the contract method 0xe0dba60f.
//
// Solidity: function setController(address controller, bool active) returns()
func (_NameWrapper *NameWrapperTransactorSession) SetController(controller common.Address, active bool) (*types.Transaction, error) {
	return _NameWrapper.Contract.SetController(&_NameWrapper.TransactOpts, controller, active)
}

// SetFuses is a paid mutator transaction binding the contract method 0x402906fc.
//
// Solidity: function setFuses(bytes32 node, uint16 ownerControlledFuses) returns(uint32)
func (_NameWrapper *NameWrapperTransactor) SetFuses(opts *bind.TransactOpts, node [32]byte, ownerControlledFuses uint16) (*types.Transaction, error) {
	return _NameWrapper.contract.Transact(opts, "setFuses", node, ownerControlledFuses)
}

// SetFuses is a paid mutator transaction binding the contract method 0x402906fc.
//
// Solidity: function setFuses(bytes32 node, uint16 ownerControlledFuses) returns(uint32)
func (_NameWrapper *NameWrapperSession) SetFuses(node [32]byte, ownerControlledFuses uint16) (*types.Transaction, error) {
	return _NameWrapper.Contract.SetFuses(&_NameWrapper.TransactOpts, node, ownerControlledFuses)
}

// SetFuses is a paid mutator transaction binding the contract method 0x402906fc.
//
// Solidity: function setFuses(bytes32 node, uint16 ownerControlledFuses) returns(uint32)
func (_NameWrapper *NameWrapperTransactorSession) SetFuses(node [32]byte, ownerControlledFuses uint16) (*types.Transaction, error) {
	return _NameWrapper.Contract.SetFuses(&_NameWrapper.TransactOpts, node, ownerControlledFuses)
}

// SetMetadataService is a paid mutator transaction binding the contract method 0x1534e177.
//
// Solidity: function setMetadataService(address _metadataService) returns()
func (_NameWrapper *NameWrapperTransactor) SetMetadataService(opts *bind.TransactOpts, _metadataService common.Address) (*types.Transaction, error) {
	return _NameWrapper.contract.Transact(opts, "setMetadataService", _metadataService)
}

// SetMetadataService is a paid mutator transaction binding the contract method 0x1534e177.
//
// Solidity: function setMetadataService(address _metadataService) returns()
func (_NameWrapper *NameWrapperSession) SetMetadataService(_metadataService common.Address) (*types.Transaction, error) {
	return _NameWrapper.Contract.SetMetadataService(&_NameWrapper.TransactOpts, _metadataService)
}

// SetMetadataService is a paid mutator transaction binding the contract method 0x1534e177.
//
// Solidity: function setMetadataService(address _metadataService) returns()
func (_NameWrapper *NameWrapperTransactorSession) SetMetadataService(_metadataService common.Address) (*types.Transaction, error) {
	return _NameWrapper.Contract.SetMetadataService(&_NameWrapper.TransactOpts, _metadataService)
}

// SetRecord is a paid mutator transaction binding the contract method 0xcf408823.
//
// Solidity: function setRecord(bytes32 node, address owner, address resolver, uint64 ttl) returns()
func (_NameWrapper *NameWrapperTransactor) SetRecord(opts *bind.TransactOpts, node [32]byte, owner common.Address, resolver common.Address, ttl uint64) (*types.Transaction, error) {
	return _NameWrapper.contract.Transact(opts, "setRecord", node, owner, resolver, ttl)
}

// SetRecord is a paid mutator transaction binding the contract method 0xcf408823.
//
// Solidity: function setRecord(bytes32 node, address owner, address resolver, uint64 ttl) returns()
func (_NameWrapper *NameWrapperSession) SetRecord(node [32]byte, owner common.Address, resolver common.Address, ttl uint64) (*types.Transaction, error) {
	return _NameWrapper.Contract.SetRecord(&_NameWrapper.TransactOpts, node, owner, resolver, ttl)
}

// SetRecord is a paid mutator transaction binding the contract method 0xcf408823.
//
// Solidity: function setRecord(bytes32 node, address owner, address resolver, uint64 ttl) returns()
func (_NameWrapper *NameWrapperTransactorSession) SetRecord(node [32]byte, owner common.Address, resolver common.Address, ttl uint64) (*types.Transaction, error) {
	return _NameWrapper.Contract.SetRecord(&_NameWrapper.TransactOpts, node, owner, resolver, ttl)
}

// SetResolver is a paid mutator transaction binding the contract method 0x1896f70a.
//
// Solidity: function setResolver(bytes32 node, address resolver) returns()
func (_NameWrapper *NameWrapperTransactor) SetResolver(opts *bind.TransactOpts, node [32]byte, resolver common.Address) (*types.Transaction, error) {
	return _NameWrapper.contract.Transact(opts, "setResolver", node, resolver)
}

// SetResolver is a paid mutator transaction binding the contract method 0x1896f70a.
//
// Solidity: function setResolver(bytes32 node, address resolver) returns()
func (_NameWrapper *NameWrapperSession) SetResolver(node [32]byte, resolver common.Address) (*types.Transaction, error) {
	return _NameWrapper.Contract.SetResolver(&_NameWrapper.TransactOpts, node, resolver)
}

// SetResolver is a paid mutator transaction binding the contract method 0x1896f70a.
//
// Solidity: function setResolver(bytes32 node, address resolver) returns()
func (_NameWrapper *NameWrapperTransactorSession) SetResolver(node [32]byte, resolver common.Address) (*types.Transaction, error) {
	return _NameWrapper.Contract.SetResolver(&_NameWrapper.TransactOpts, node, resolver)
}

// SetSubnodeOwner is a paid mutator transaction binding the contract method 0xc658e086.
//
// Solidity: function setSubnodeOwner(bytes32 parentNode, string label, address owner, uint32 fuses, uint64 expiry) returns(bytes32 node)
func (_NameWrapper *NameWrapperTransactor) SetSubnodeOwner(opts *bind.TransactOpts, parentNode [32]byte, label string, owner common.Address, fuses uint32, expiry uint64) (*types.Transaction, error) {
	return _NameWrapper.contract.Transact(opts, "setSubnodeOwner", parentNode, label, owner, fuses, expiry)
}

// SetSubnodeOwner is a paid mutator transaction binding the contract method 0xc658e086.
//
// Solidity: function setSubnodeOwner(bytes32 parentNode, string label, address owner, uint32 fuses, uint64 expiry) returns(bytes32 node)
func (_NameWrapper *NameWrapperSession) SetSubnodeOwner(parentNode [32]byte, label string, owner common.Address, fuses uint32, expiry uint64) (*types.Transaction, error) {
	return _NameWrapper.Contract.SetSubnodeOwner(&_NameWrapper.TransactOpts, parentNode, label, owner, fuses, expiry)
}

// SetSubnodeOwner is a paid mutator transaction binding the contract method 0xc658e086.
//
// Solidity: function setSubnodeOwner(bytes32 parentNode, string label, address owner, uint32 fuses, uint64 expiry) returns(bytes32 node)
func (_NameWrapper *NameWrapperTransactorSession) SetSubnodeOwner(parentNode [32]byte, label string, owner common.Address, fuses uint32, expiry uint64) (*types.Transaction, error) {
	return _NameWrapper.Contract.SetSubnodeOwner(&_NameWrapper.TransactOpts, parentNode, label, owner, fuses, expiry)
}

// SetSubnodeRecord is a paid mutator transaction binding the contract method 0x24c1af44.
//
// Solidity: function setSubnodeRecord(bytes32 parentNode, string label, address owner, address resolver, uint64 ttl, uint32 fuses, uint64 expiry) returns(bytes32 node)
func (_NameWrapper *NameWrapperTransactor) SetSubnodeRecord(opts *bind.TransactOpts, parentNode [32]byte, label string, owner common.Address, resolver common.Address, ttl uint64, fuses uint32, expiry uint64) (*types.Transaction, error) {
	return _NameWrapper.contract.Transact(opts, "setSubnodeRecord", parentNode, label, owner, resolver, ttl, fuses, expiry)
}

// SetSubnodeRecord is a paid mutator transaction binding the contract method 0x24c1af44.
//
// Solidity: function setSubnodeRecord(bytes32 parentNode, string label, address owner, address resolver, uint64 ttl, uint32 fuses, uint64 expiry) returns(bytes32 node)
func (_NameWrapper *NameWrapperSession) SetSubnodeRecord(parentNode [32]byte, label string, owner common.Address, resolver common.Address, ttl uint64, fuses uint32, expiry uint64) (*types.Transaction, error) {
	return _NameWrapper.Contract.SetSubnodeRecord(&_NameWrapper.TransactOpts, parentNode, label, owner, resolver, ttl, fuses, expiry)
}

// SetSubnodeRecord is a paid mutator transaction binding the contract method 0x24c1af44.
//
// Solidity: function setSubnodeRecord(bytes32 parentNode, string label, address owner, address resolver, uint64 ttl, uint32 fuses, uint64 expiry) returns(bytes32 node)
func (_NameWrapper *NameWrapperTransactorSession) SetSubnodeRecord(parentNode [32]byte, label string, owner common.Address, resolver common.Address, ttl uint64, fuses uint32, expiry uint64) (*types.Transaction, error) {
	return _NameWrapper.Contract.SetSubnodeRecord(&_NameWrapper.TransactOpts, parentNode, label, owner, resolver, ttl, fuses, expiry)
}

// SetTTL is a paid mutator transaction binding the contract method 0x14ab9038.
//
// Solidity: function setTTL(bytes32 node, uint64 ttl) returns()
func (_NameWrapper *NameWrapperTransactor) SetTTL(opts *bind.TransactOpts, node [32]byte, ttl uint64) (*types.Transaction, error) {
	return _NameWrapper.contract.Transact(opts, "setTTL", node, ttl)
}

// SetTTL is a paid mutator transaction binding the contract method 0x14ab9038.
//
// Solidity: function setTTL(bytes32 node, uint64 ttl) returns()
func (_NameWrapper *NameWrapperSession) SetTTL(node [32]byte, ttl uint64) (*types.Transaction, error) {
	return _NameWrapper.Contract.SetTTL(&_NameWrapper.TransactOpts, node, ttl)
}

// SetTTL is a paid mutator transaction binding the contract method 0x14ab9038.
//
// Solidity: function setTTL(bytes32 node, uint64 ttl) returns()
func (_NameWrapper *NameWrapperTransactorSession) SetTTL(node [32]byte, ttl uint64) (*types.Transaction, error) {
	return _NameWrapper.Contract.SetTTL(&_NameWrapper.TransactOpts, node, ttl)
}

// SetUpgradeContract is a paid mutator transaction binding the contract method 0xb6bcad26.
//
// Solidity: function setUpgradeContract(address _upgradeAddress) returns()
func (_NameWrapper *NameWrapperTransactor) SetUpgradeContract(opts *bind.TransactOpts, _upgradeAddress common.Address) (*types.Transaction, error) {
	return _NameWrapper.contract.Transact(opts, "setUpgradeContract", _upgradeAddress)
}

// SetUpgradeContract is a paid mutator transaction binding the contract method 0xb6bcad26.
//
// Solidity: function setUpgradeContract(address _upgradeAddress) returns()
func (_NameWrapper *NameWrapperSession) SetUpgradeContract(_upgradeAddress common.Address) (*types.Transaction, error) {
	return _NameWrapper.Contract.SetUpgradeContract(&_NameWrapper.TransactOpts, _upgradeAddress)
}

// SetUpgradeContract is a paid mutator transaction binding the contract method 0xb6bcad26.
//
// Solidity: function setUpgradeContract(address _upgradeAddress) returns()
func (_NameWrapper *NameWrapperTransactorSession) SetUpgradeContract(_upgradeAddress common.Address) (*types.Transaction, error) {
	return _NameWrapper.Contract.SetUpgradeContract(&_NameWrapper.TransactOpts, _upgradeAddress)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_NameWrapper *NameWrapperTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _NameWrapper.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_NameWrapper *NameWrapperSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _NameWrapper.Contract.TransferOwnership(&_NameWrapper.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_NameWrapper *NameWrapperTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _NameWrapper.Contract.TransferOwnership(&_NameWrapper.TransactOpts, newOwner)
}

// Unwrap is a paid mutator transaction binding the contract method 0xd8c9921a.
//
// Solidity: function unwrap(bytes32 parentNode, bytes32 labelhash, address controller) returns()
func (_NameWrapper *NameWrapperTransactor) Unwrap(opts *bind.TransactOpts, parentNode [32]byte, labelhash [32]byte, controller common.Address) (*types.Transaction, error) {
	return _NameWrapper.contract.Transact(opts, "unwrap", parentNode, labelhash, controller)
}

// Unwrap is a paid mutator transaction binding the contract method 0xd8c9921a.
//
// Solidity: function unwrap(bytes32 parentNode, bytes32 labelhash, address controller) returns()
func (_NameWrapper *NameWrapperSession) Unwrap(parentNode [32]byte, labelhash [32]byte, controller common.Address) (*types.Transaction, error) {
	return _NameWrapper.Contract.Unwrap(&_NameWrapper.TransactOpts, parentNode, labelhash, controller)
}

// Unwrap is a paid mutator transaction binding the contract method 0xd8c9921a.
//
// Solidity: function unwrap(bytes32 parentNode, bytes32 labelhash, address controller) returns()
func (_NameWrapper *NameWrapperTransactorSession) Unwrap(parentNode [32]byte, labelhash [32]byte, controller common.Address) (*types.Transaction, error) {
	return _NameWrapper.Contract.Unwrap(&_NameWrapper.TransactOpts, parentNode, labelhash, controller)
}

// UnwrapETH2LD is a paid mutator transaction binding the contract method 0x8b4dfa75.
//
// Solidity: function unwrapETH2LD(bytes32 labelhash, address registrant, address controller) returns()
func (_NameWrapper *NameWrapperTransactor) UnwrapETH2LD(opts *bind.TransactOpts, labelhash [32]byte, registrant common.Address, controller common.Address) (*types.Transaction, error) {
	return _NameWrapper.contract.Transact(opts, "unwrapETH2LD", labelhash, registrant, controller)
}

// UnwrapETH2LD is a paid mutator transaction binding the contract method 0x8b4dfa75.
//
// Solidity: function unwrapETH2LD(bytes32 labelhash, address registrant, address controller) returns()
func (_NameWrapper *NameWrapperSession) UnwrapETH2LD(labelhash [32]byte, registrant common.Address, controller common.Address) (*types.Transaction, error) {
	return _NameWrapper.Contract.UnwrapETH2LD(&_NameWrapper.TransactOpts, labelhash, registrant, controller)
}

// UnwrapETH2LD is a paid mutator transaction binding the contract method 0x8b4dfa75.
//
// Solidity: function unwrapETH2LD(bytes32 labelhash, address registrant, address controller) returns()
func (_NameWrapper *NameWrapperTransactorSession) UnwrapETH2LD(labelhash [32]byte, registrant common.Address, controller common.Address) (*types.Transaction, error) {
	return _NameWrapper.Contract.UnwrapETH2LD(&_NameWrapper.TransactOpts, labelhash, registrant, controller)
}

// Upgrade is a paid mutator transaction binding the contract method 0xc93ab3fd.
//
// Solidity: function upgrade(bytes name, bytes extraData) returns()
func (_NameWrapper *NameWrapperTransactor) Upgrade(opts *bind.TransactOpts, name []byte, extraData []byte) (*types.Transaction, error) {
	return _NameWrapper.contract.Transact(opts, "upgrade", name, extraData)
}

// Upgrade is a paid mutator transaction binding the contract method 0xc93ab3fd.
//
// Solidity: function upgrade(bytes name, bytes extraData) returns()
func (_NameWrapper *NameWrapperSession) Upgrade(name []byte, extraData []byte) (*types.Transaction, error) {
	return _NameWrapper.Contract.Upgrade(&_NameWrapper.TransactOpts, name, extraData)
}

// Upgrade is a paid mutator transaction binding the contract method 0xc93ab3fd.
//
// Solidity: function upgrade(bytes name, bytes extraData) returns()
func (_NameWrapper *NameWrapperTransactorSession) Upgrade(name []byte, extraData []byte) (*types.Transaction, error) {
	return _NameWrapper.Contract.Upgrade(&_NameWrapper.TransactOpts, name, extraData)
}

// Wrap is a paid mutator transaction binding the contract method 0xeb8ae530.
//
// Solidity: function wrap(bytes name, address wrappedOwner, address resolver) returns()
func (_NameWrapper *NameWrapperTransactor) Wrap(opts *bind.TransactOpts, name []byte, wrappedOwner common.Address, resolver common.Address) (*types.Transaction, error) {
	return _NameWrapper.contract.Transact(opts, "wrap", name, wrappedOwner, resolver)
}

// Wrap is a paid mutator transaction binding the contract method 0xeb8ae530.
//
// Solidity: function wrap(bytes name, address wrappedOwner, address resolver) returns()
func (_NameWrapper *NameWrapperSession) Wrap(name []byte, wrappedOwner common.Address, resolver common.Address) (*types.Transaction, error) {
	return _NameWrapper.Contract.Wrap(&_NameWrapper.TransactOpts, name, wrappedOwner, resolver)
}

// Wrap is a paid mutator transaction binding the contract method 0xeb8ae530.
//
// Solidity: function wrap(bytes name, address wrappedOwner, address resolver) returns()
func (_NameWrapper *NameWrapperTransactorSession) Wrap(name []byte, wrappedOwner common.Address, resolver common.Address) (*types.Transaction, error) {
	return _NameWrapper.Contract.Wrap(&_NameWrapper.TransactOpts, name, wrappedOwner, resolver)
}

// WrapETH2LD is a paid mutator transaction binding the contract method 0x8cf8b41e.
//
// Solidity: function wrapETH2LD(string label, address wrappedOwner, uint16 ownerControlledFuses, address resolver) returns(uint64 expiry)
func (_NameWrapper *NameWrapperTransactor) WrapETH2LD(opts *bind.TransactOpts, label string, wrappedOwner common.Address, ownerControlledFuses uint16, resolver common.Address) (*types.Transaction, error) {
	return _NameWrapper.contract.Transact(opts, "wrapETH2LD", label, wrappedOwner, ownerControlledFuses, resolver)
}

// WrapETH2LD is a paid mutator transaction binding the contract method 0x8cf8b41e.
//
// Solidity: function wrapETH2LD(string label, address wrappedOwner, uint16 ownerControlledFuses, address resolver) returns(uint64 expiry)
func (_NameWrapper *NameWrapperSession) WrapETH2LD(label string, wrappedOwner common.Address, ownerControlledFuses uint16, resolver common.Address) (*types.Transaction, error) {
	return _NameWrapper.Contract.WrapETH2LD(&_NameWrapper.TransactOpts, label, wrappedOwner, ownerControlledFuses, resolver)
}

// WrapETH2LD is a paid mutator transaction binding the contract method 0x8cf8b41e.
//
// Solidity: function wrapETH2LD(string label, address wrappedOwner, uint16 ownerControlledFuses, address resolver) returns(uint64 expiry)
func (_NameWrapper *NameWrapperTransactorSession) WrapETH2LD(label string, wrappedOwner common.Address, ownerControlledFuses uint16, resolver common.Address) (*types.Transaction, error) {
	return _NameWrapper.Contract.WrapETH2LD(&_NameWrapper.TransactOpts, label, wrappedOwner, ownerControlledFuses, resolver)
}

// NameWrapperApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the NameWrapper contract.
type NameWrapperApprovalIterator struct {
	Event *NameWrapperApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NameWrapperApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NameWrapperApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NameWrapperApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NameWrapperApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NameWrapperApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NameWrapperApproval represents a Approval event raised by the NameWrapper contract.
type NameWrapperApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_NameWrapper *NameWrapperFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*NameWrapperApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _NameWrapper.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &NameWrapperApprovalIterator{contract: _NameWrapper.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_NameWrapper *NameWrapperFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *NameWrapperApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _NameWrapper.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NameWrapperApproval)
				if err := _NameWrapper.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_NameWrapper *NameWrapperFilterer) ParseApproval(log types.Log) (*NameWrapperApproval, error) {
	event := new(NameWrapperApproval)
	if err := _NameWrapper.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NameWrapperApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the NameWrapper contract.
type NameWrapperApprovalForAllIterator struct {
	Event *NameWrapperApprovalForAll // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NameWrapperApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NameWrapperApprovalForAll)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NameWrapperApprovalForAll)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NameWrapperApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NameWrapperApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NameWrapperApprovalForAll represents a ApprovalForAll event raised by the NameWrapper contract.
type NameWrapperApprovalForAll struct {
	Account  common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed account, address indexed operator, bool approved)
func (_NameWrapper *NameWrapperFilterer) FilterApprovalForAll(opts *bind.FilterOpts, account []common.Address, operator []common.Address) (*NameWrapperApprovalForAllIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _NameWrapper.contract.FilterLogs(opts, "ApprovalForAll", accountRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &NameWrapperApprovalForAllIterator{contract: _NameWrapper.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed account, address indexed operator, bool approved)
func (_NameWrapper *NameWrapperFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *NameWrapperApprovalForAll, account []common.Address, operator []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _NameWrapper.contract.WatchLogs(opts, "ApprovalForAll", accountRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NameWrapperApprovalForAll)
				if err := _NameWrapper.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApprovalForAll is a log parse operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed account, address indexed operator, bool approved)
func (_NameWrapper *NameWrapperFilterer) ParseApprovalForAll(log types.Log) (*NameWrapperApprovalForAll, error) {
	event := new(NameWrapperApprovalForAll)
	if err := _NameWrapper.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NameWrapperControllerChangedIterator is returned from FilterControllerChanged and is used to iterate over the raw logs and unpacked data for ControllerChanged events raised by the NameWrapper contract.
type NameWrapperControllerChangedIterator struct {
	Event *NameWrapperControllerChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NameWrapperControllerChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NameWrapperControllerChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NameWrapperControllerChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NameWrapperControllerChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NameWrapperControllerChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NameWrapperControllerChanged represents a ControllerChanged event raised by the NameWrapper contract.
type NameWrapperControllerChanged struct {
	Controller common.Address
	Active     bool
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterControllerChanged is a free log retrieval operation binding the contract event 0x4c97694570a07277810af7e5669ffd5f6a2d6b74b6e9a274b8b870fd5114cf87.
//
// Solidity: event ControllerChanged(address indexed controller, bool active)
func (_NameWrapper *NameWrapperFilterer) FilterControllerChanged(opts *bind.FilterOpts, controller []common.Address) (*NameWrapperControllerChangedIterator, error) {

	var controllerRule []interface{}
	for _, controllerItem := range controller {
		controllerRule = append(controllerRule, controllerItem)
	}

	logs, sub, err := _NameWrapper.contract.FilterLogs(opts, "ControllerChanged", controllerRule)
	if err != nil {
		return nil, err
	}
	return &NameWrapperControllerChangedIterator{contract: _NameWrapper.contract, event: "ControllerChanged", logs: logs, sub: sub}, nil
}

// WatchControllerChanged is a free log subscription operation binding the contract event 0x4c97694570a07277810af7e5669ffd5f6a2d6b74b6e9a274b8b870fd5114cf87.
//
// Solidity: event ControllerChanged(address indexed controller, bool active)
func (_NameWrapper *NameWrapperFilterer) WatchControllerChanged(opts *bind.WatchOpts, sink chan<- *NameWrapperControllerChanged, controller []common.Address) (event.Subscription, error) {

	var controllerRule []interface{}
	for _, controllerItem := range controller {
		controllerRule = append(controllerRule, controllerItem)
	}

	logs, sub, err := _NameWrapper.contract.WatchLogs(opts, "ControllerChanged", controllerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NameWrapperControllerChanged)
				if err := _NameWrapper.contract.UnpackLog(event, "ControllerChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseControllerChanged is a log parse operation binding the contract event 0x4c97694570a07277810af7e5669ffd5f6a2d6b74b6e9a274b8b870fd5114cf87.
//
// Solidity: event ControllerChanged(address indexed controller, bool active)
func (_NameWrapper *NameWrapperFilterer) ParseControllerChanged(log types.Log) (*NameWrapperControllerChanged, error) {
	event := new(NameWrapperControllerChanged)
	if err := _NameWrapper.contract.UnpackLog(event, "ControllerChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NameWrapperExpiryExtendedIterator is returned from FilterExpiryExtended and is used to iterate over the raw logs and unpacked data for ExpiryExtended events raised by the NameWrapper contract.
type NameWrapperExpiryExtendedIterator struct {
	Event *NameWrapperExpiryExtended // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NameWrapperExpiryExtendedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NameWrapperExpiryExtended)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NameWrapperExpiryExtended)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NameWrapperExpiryExtendedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NameWrapperExpiryExtendedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NameWrapperExpiryExtended represents a ExpiryExtended event raised by the NameWrapper contract.
type NameWrapperExpiryExtended struct {
	Node   [32]byte
	Expiry uint64
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterExpiryExtended is a free log retrieval operation binding the contract event 0xf675815a0817338f93a7da433f6bd5f5542f1029b11b455191ac96c7f6a9b132.
//
// Solidity: event ExpiryExtended(bytes32 indexed node, uint64 expiry)
func (_NameWrapper *NameWrapperFilterer) FilterExpiryExtended(opts *bind.FilterOpts, node [][32]byte) (*NameWrapperExpiryExtendedIterator, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _NameWrapper.contract.FilterLogs(opts, "ExpiryExtended", nodeRule)
	if err != nil {
		return nil, err
	}
	return &NameWrapperExpiryExtendedIterator{contract: _NameWrapper.contract, event: "ExpiryExtended", logs: logs, sub: sub}, nil
}

// WatchExpiryExtended is a free log subscription operation binding the contract event 0xf675815a0817338f93a7da433f6bd5f5542f1029b11b455191ac96c7f6a9b132.
//
// Solidity: event ExpiryExtended(bytes32 indexed node, uint64 expiry)
func (_NameWrapper *NameWrapperFilterer) WatchExpiryExtended(opts *bind.WatchOpts, sink chan<- *NameWrapperExpiryExtended, node [][32]byte) (event.Subscription, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _NameWrapper.contract.WatchLogs(opts, "ExpiryExtended", nodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NameWrapperExpiryExtended)
				if err := _NameWrapper.contract.UnpackLog(event, "ExpiryExtended", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseExpiryExtended is a log parse operation binding the contract event 0xf675815a0817338f93a7da433f6bd5f5542f1029b11b455191ac96c7f6a9b132.
//
// Solidity: event ExpiryExtended(bytes32 indexed node, uint64 expiry)
func (_NameWrapper *NameWrapperFilterer) ParseExpiryExtended(log types.Log) (*NameWrapperExpiryExtended, error) {
	event := new(NameWrapperExpiryExtended)
	if err := _NameWrapper.contract.UnpackLog(event, "ExpiryExtended", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NameWrapperFusesSetIterator is returned from FilterFusesSet and is used to iterate over the raw logs and unpacked data for FusesSet events raised by the NameWrapper contract.
type NameWrapperFusesSetIterator struct {
	Event *NameWrapperFusesSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NameWrapperFusesSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NameWrapperFusesSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NameWrapperFusesSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NameWrapperFusesSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NameWrapperFusesSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NameWrapperFusesSet represents a FusesSet event raised by the NameWrapper contract.
type NameWrapperFusesSet struct {
	Node  [32]byte
	Fuses uint32
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterFusesSet is a free log retrieval operation binding the contract event 0x39873f00c80f4f94b7bd1594aebcf650f003545b74824d57ddf4939e3ff3a34b.
//
// Solidity: event FusesSet(bytes32 indexed node, uint32 fuses)
func (_NameWrapper *NameWrapperFilterer) FilterFusesSet(opts *bind.FilterOpts, node [][32]byte) (*NameWrapperFusesSetIterator, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _NameWrapper.contract.FilterLogs(opts, "FusesSet", nodeRule)
	if err != nil {
		return nil, err
	}
	return &NameWrapperFusesSetIterator{contract: _NameWrapper.contract, event: "FusesSet", logs: logs, sub: sub}, nil
}

// WatchFusesSet is a free log subscription operation binding the contract event 0x39873f00c80f4f94b7bd1594aebcf650f003545b74824d57ddf4939e3ff3a34b.
//
// Solidity: event FusesSet(bytes32 indexed node, uint32 fuses)
func (_NameWrapper *NameWrapperFilterer) WatchFusesSet(opts *bind.WatchOpts, sink chan<- *NameWrapperFusesSet, node [][32]byte) (event.Subscription, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _NameWrapper.contract.WatchLogs(opts, "FusesSet", nodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NameWrapperFusesSet)
				if err := _NameWrapper.contract.UnpackLog(event, "FusesSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseFusesSet is a log parse operation binding the contract event 0x39873f00c80f4f94b7bd1594aebcf650f003545b74824d57ddf4939e3ff3a34b.
//
// Solidity: event FusesSet(bytes32 indexed node, uint32 fuses)
func (_NameWrapper *NameWrapperFilterer) ParseFusesSet(log types.Log) (*NameWrapperFusesSet, error) {
	event := new(NameWrapperFusesSet)
	if err := _NameWrapper.contract.UnpackLog(event, "FusesSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NameWrapperNameUnwrappedIterator is returned from FilterNameUnwrapped and is used to iterate over the raw logs and unpacked data for NameUnwrapped events raised by the NameWrapper contract.
type NameWrapperNameUnwrappedIterator struct {
	Event *NameWrapperNameUnwrapped // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NameWrapperNameUnwrappedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NameWrapperNameUnwrapped)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NameWrapperNameUnwrapped)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NameWrapperNameUnwrappedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NameWrapperNameUnwrappedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NameWrapperNameUnwrapped represents a NameUnwrapped event raised by the NameWrapper contract.
type NameWrapperNameUnwrapped struct {
	Node  [32]byte
	Owner common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterNameUnwrapped is a free log retrieval operation binding the contract event 0xee2ba1195c65bcf218a83d874335c6bf9d9067b4c672f3c3bf16cf40de7586c4.
//
// Solidity: event NameUnwrapped(bytes32 indexed node, address owner)
func (_NameWrapper *NameWrapperFilterer) FilterNameUnwrapped(opts *bind.FilterOpts, node [][32]byte) (*NameWrapperNameUnwrappedIterator, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _NameWrapper.contract.FilterLogs(opts, "NameUnwrapped", nodeRule)
	if err != nil {
		return nil, err
	}
	return &NameWrapperNameUnwrappedIterator{contract: _NameWrapper.contract, event: "NameUnwrapped", logs: logs, sub: sub}, nil
}

// WatchNameUnwrapped is a free log subscription operation binding the contract event 0xee2ba1195c65bcf218a83d874335c6bf9d9067b4c672f3c3bf16cf40de7586c4.
//
// Solidity: event NameUnwrapped(bytes32 indexed node, address owner)
func (_NameWrapper *NameWrapperFilterer) WatchNameUnwrapped(opts *bind.WatchOpts, sink chan<- *NameWrapperNameUnwrapped, node [][32]byte) (event.Subscription, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _NameWrapper.contract.WatchLogs(opts, "NameUnwrapped", nodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NameWrapperNameUnwrapped)
				if err := _NameWrapper.contract.UnpackLog(event, "NameUnwrapped", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNameUnwrapped is a log parse operation binding the contract event 0xee2ba1195c65bcf218a83d874335c6bf9d9067b4c672f3c3bf16cf40de7586c4.
//
// Solidity: event NameUnwrapped(bytes32 indexed node, address owner)
func (_NameWrapper *NameWrapperFilterer) ParseNameUnwrapped(log types.Log) (*NameWrapperNameUnwrapped, error) {
	event := new(NameWrapperNameUnwrapped)
	if err := _NameWrapper.contract.UnpackLog(event, "NameUnwrapped", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NameWrapperNameWrappedIterator is returned from FilterNameWrapped and is used to iterate over the raw logs and unpacked data for NameWrapped events raised by the NameWrapper contract.
type NameWrapperNameWrappedIterator struct {
	Event *NameWrapperNameWrapped // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NameWrapperNameWrappedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NameWrapperNameWrapped)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NameWrapperNameWrapped)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NameWrapperNameWrappedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NameWrapperNameWrappedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NameWrapperNameWrapped represents a NameWrapped event raised by the NameWrapper contract.
type NameWrapperNameWrapped struct {
	Node   [32]byte
	Name   []byte
	Owner  common.Address
	Fuses  uint32
	Expiry uint64
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterNameWrapped is a free log retrieval operation binding the contract event 0x8ce7013e8abebc55c3890a68f5a27c67c3f7efa64e584de5fb22363c606fd340.
//
// Solidity: event NameWrapped(bytes32 indexed node, bytes name, address owner, uint32 fuses, uint64 expiry)
func (_NameWrapper *NameWrapperFilterer) FilterNameWrapped(opts *bind.FilterOpts, node [][32]byte) (*NameWrapperNameWrappedIterator, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _NameWrapper.contract.FilterLogs(opts, "NameWrapped", nodeRule)
	if err != nil {
		return nil, err
	}
	return &NameWrapperNameWrappedIterator{contract: _NameWrapper.contract, event: "NameWrapped", logs: logs, sub: sub}, nil
}

// WatchNameWrapped is a free log subscription operation binding the contract event 0x8ce7013e8abebc55c3890a68f5a27c67c3f7efa64e584de5fb22363c606fd340.
//
// Solidity: event NameWrapped(bytes32 indexed node, bytes name, address owner, uint32 fuses, uint64 expiry)
func (_NameWrapper *NameWrapperFilterer) WatchNameWrapped(opts *bind.WatchOpts, sink chan<- *NameWrapperNameWrapped, node [][32]byte) (event.Subscription, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _NameWrapper.contract.WatchLogs(opts, "NameWrapped", nodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NameWrapperNameWrapped)
				if err := _NameWrapper.contract.UnpackLog(event, "NameWrapped", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNameWrapped is a log parse operation binding the contract event 0x8ce7013e8abebc55c3890a68f5a27c67c3f7efa64e584de5fb22363c606fd340.
//
// Solidity: event NameWrapped(bytes32 indexed node, bytes name, address owner, uint32 fuses, uint64 expiry)
func (_NameWrapper *NameWrapperFilterer) ParseNameWrapped(log types.Log) (*NameWrapperNameWrapped, error) {
	event := new(NameWrapperNameWrapped)
	if err := _NameWrapper.contract.UnpackLog(event, "NameWrapped", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NameWrapperOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the NameWrapper contract.
type NameWrapperOwnershipTransferredIterator struct {
	Event *NameWrapperOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NameWrapperOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NameWrapperOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NameWrapperOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NameWrapperOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NameWrapperOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NameWrapperOwnershipTransferred represents a OwnershipTransferred event raised by the NameWrapper contract.
type NameWrapperOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_NameWrapper *NameWrapperFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*NameWrapperOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _NameWrapper.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &NameWrapperOwnershipTransferredIterator{contract: _NameWrapper.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_NameWrapper *NameWrapperFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *NameWrapperOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _NameWrapper.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NameWrapperOwnershipTransferred)
				if err := _NameWrapper.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_NameWrapper *NameWrapperFilterer) ParseOwnershipTransferred(log types.Log) (*NameWrapperOwnershipTransferred, error) {
	event := new(NameWrapperOwnershipTransferred)
	if err := _NameWrapper.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NameWrapperTransferBatchIterator is returned from FilterTransferBatch and is used to iterate over the raw logs and unpacked data for TransferBatch events raised by the NameWrapper contract.
type NameWrapperTransferBatchIterator struct {
	Event *NameWrapperTransferBatch // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NameWrapperTransferBatchIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NameWrapperTransferBatch)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NameWrapperTransferBatch)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NameWrapperTransferBatchIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NameWrapperTransferBatchIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NameWrapperTransferBatch represents a TransferBatch event raised by the NameWrapper contract.
type NameWrapperTransferBatch struct {
	Operator common.Address
	From     common.Address
	To       common.Address
	Ids      []*big.Int
	Values   []*big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTransferBatch is a free log retrieval operation binding the contract event 0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb.
//
// Solidity: event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] values)
func (_NameWrapper *NameWrapperFilterer) FilterTransferBatch(opts *bind.FilterOpts, operator []common.Address, from []common.Address, to []common.Address) (*NameWrapperTransferBatchIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _NameWrapper.contract.FilterLogs(opts, "TransferBatch", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &NameWrapperTransferBatchIterator{contract: _NameWrapper.contract, event: "TransferBatch", logs: logs, sub: sub}, nil
}

// WatchTransferBatch is a free log subscription operation binding the contract event 0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb.
//
// Solidity: event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] values)
func (_NameWrapper *NameWrapperFilterer) WatchTransferBatch(opts *bind.WatchOpts, sink chan<- *NameWrapperTransferBatch, operator []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _NameWrapper.contract.WatchLogs(opts, "TransferBatch", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NameWrapperTransferBatch)
				if err := _NameWrapper.contract.UnpackLog(event, "TransferBatch", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransferBatch is a log parse operation binding the contract event 0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb.
//
// Solidity: event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] values)
func (_NameWrapper *NameWrapperFilterer) ParseTransferBatch(log types.Log) (*NameWrapperTransferBatch, error) {
	event := new(NameWrapperTransferBatch)
	if err := _NameWrapper.contract.UnpackLog(event, "TransferBatch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NameWrapperTransferSingleIterator is returned from FilterTransferSingle and is used to iterate over the raw logs and unpacked data for TransferSingle events raised by the NameWrapper contract.
type NameWrapperTransferSingleIterator struct {
	Event *NameWrapperTransferSingle // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NameWrapperTransferSingleIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NameWrapperTransferSingle)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NameWrapperTransferSingle)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NameWrapperTransferSingleIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NameWrapperTransferSingleIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NameWrapperTransferSingle represents a TransferSingle event raised by the NameWrapper contract.
type NameWrapperTransferSingle struct {
	Operator common.Address
	From     common.Address
	To       common.Address
	Id       *big.Int
	Value    *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTransferSingle is a free log retrieval operation binding the contract event 0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62.
//
// Solidity: event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value)
func (_NameWrapper *NameWrapperFilterer) FilterTransferSingle(opts *bind.FilterOpts, operator []common.Address, from []common.Address, to []common.Address) (*NameWrapperTransferSingleIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _NameWrapper.contract.FilterLogs(opts, "TransferSingle", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &NameWrapperTransferSingleIterator{contract: _NameWrapper.contract, event: "TransferSingle", logs: logs, sub: sub}, nil
}

// WatchTransferSingle is a free log subscription operation binding the contract event 0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62.
//
// Solidity: event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value)
func (_NameWrapper *NameWrapperFilterer) WatchTransferSingle(opts *bind.WatchOpts, sink chan<- *NameWrapperTransferSingle, operator []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _NameWrapper.contract.WatchLogs(opts, "TransferSingle", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NameWrapperTransferSingle)
				if err := _NameWrapper.contract.UnpackLog(event, "TransferSingle", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransferSingle is a log parse operation binding the contract event 0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62.
//
// Solidity: event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value)
func (_NameWrapper *NameWrapperFilterer) ParseTransferSingle(log types.Log) (*NameWrapperTransferSingle, error) {
	event := new(NameWrapperTransferSingle)
	if err := _NameWrapper.contract.UnpackLog(event, "TransferSingle", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NameWrapperURIIterator is returned from FilterURI and is used to iterate over the raw logs and unpacked data for URI events raised by the NameWrapper contract.
type NameWrapperURIIterator struct {
	Event *NameWrapperURI // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NameWrapperURIIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NameWrapperURI)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NameWrapperURI)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NameWrapperURIIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NameWrapperURIIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NameWrapperURI represents a URI event raised by the NameWrapper contract.
type NameWrapperURI struct {
	Value string
	Id    *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterURI is a free log retrieval operation binding the contract event 0x6bb7ff708619ba0610cba295a58592e0451dee2622938c8755667688daf3529b.
//
// Solidity: event URI(string value, uint256 indexed id)
func (_NameWrapper *NameWrapperFilterer) FilterURI(opts *bind.FilterOpts, id []*big.Int) (*NameWrapperURIIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _NameWrapper.contract.FilterLogs(opts, "URI", idRule)
	if err != nil {
		return nil, err
	}
	return &NameWrapperURIIterator{contract: _NameWrapper.contract, event: "URI", logs: logs, sub: sub}, nil
}

// WatchURI is a free log subscription operation binding the contract event 0x6bb7ff708619ba0610cba295a58592e0451dee2622938c8755667688daf3529b.
//
// Solidity: event URI(string value, uint256 indexed id)
func (_NameWrapper *NameWrapperFilterer) WatchURI(opts *bind.WatchOpts, sink chan<- *NameWrapperURI, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _NameWrapper.contract.WatchLogs(opts, "URI", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NameWrapperURI)
				if err := _NameWrapper.contract.UnpackLog(event, "URI", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseURI is a log parse operation binding the contract event 0x6bb7ff708619ba0610cba295a58592e0451dee2622938c8755667688daf3529b.
//
// Solidity: event URI(string value, uint256 indexed id)
func (_NameWrapper *NameWrapperFilterer) ParseURI(log types.Log) (*NameWrapperURI, error) {
	event := new(NameWrapperURI)
	if err := _NameWrapper.contract.UnpackLog(event, "URI", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
