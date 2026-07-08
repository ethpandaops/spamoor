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

// StaticMetadataServiceMetaData contains all meta data concerning the StaticMetadataService contract.
var StaticMetadataServiceMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_metaDataUri\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"uri\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060405161047538038061047583398101604081905261002f91610058565b600061003b82826101ad565b505061026b565b634e487b7160e01b600052604160045260246000fd5b60006020828403121561006a57600080fd5b81516001600160401b0381111561008057600080fd5b8201601f8101841361009157600080fd5b80516001600160401b038111156100aa576100aa610042565b604051601f8201601f19908116603f011681016001600160401b03811182821017156100d8576100d8610042565b6040528181528282016020018610156100f057600080fd5b60005b8281101561010f576020818501810151838301820152016100f3565b50600091810160200191909152949350505050565b600181811c9082168061013857607f821691505b60208210810361015857634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156101a857806000526020600020601f840160051c810160208510156101855750805b601f840160051c820191505b818110156101a55760008155600101610191565b50505b505050565b81516001600160401b038111156101c6576101c6610042565b6101da816101d48454610124565b8461015e565b6020601f82116001811461020e57600083156101f65750848201515b600019600385901b1c1916600184901b1784556101a5565b600084815260208120601f198516915b8281101561023e578785015182556020948501946001909201910161021e565b508482101561025c5786840151600019600387901b60f8161c191681555b50505050600190811b01905550565b6101fb8061027a6000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c80630e89341c14610030575b600080fd5b61004361003e3660046100ed565b610059565b6040516100509190610106565b60405180910390f35b60606000805461006890610172565b80601f016020809104026020016040519081016040528092919081815260200182805461009490610172565b80156100e15780601f106100b6576101008083540402835291602001916100e1565b820191906000526020600020905b8154815290600101906020018083116100c457829003601f168201915b50505050509050919050565b6000602082840312156100ff57600080fd5b5035919050565b602081526000825180602084015260005b818110156101345760208186018101516040868401015201610117565b5060006040828501015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011684010191505092915050565b600181811c9082168061018657607f821691505b6020821081036101bf577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b5091905056fea26469706673582212203c3ed2817f858b86979cd0f2b514be4da40fa9e27af4ad303d7957c72d0cee5b64736f6c634300081a0033",
}

// StaticMetadataServiceABI is the input ABI used to generate the binding from.
// Deprecated: Use StaticMetadataServiceMetaData.ABI instead.
var StaticMetadataServiceABI = StaticMetadataServiceMetaData.ABI

// StaticMetadataServiceBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use StaticMetadataServiceMetaData.Bin instead.
var StaticMetadataServiceBin = StaticMetadataServiceMetaData.Bin

// DeployStaticMetadataService deploys a new Ethereum contract, binding an instance of StaticMetadataService to it.
func DeployStaticMetadataService(auth *bind.TransactOpts, backend bind.ContractBackend, _metaDataUri string) (common.Address, *types.Transaction, *StaticMetadataService, error) {
	parsed, err := StaticMetadataServiceMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StaticMetadataServiceBin), backend, _metaDataUri)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &StaticMetadataService{StaticMetadataServiceCaller: StaticMetadataServiceCaller{contract: contract}, StaticMetadataServiceTransactor: StaticMetadataServiceTransactor{contract: contract}, StaticMetadataServiceFilterer: StaticMetadataServiceFilterer{contract: contract}}, nil
}

// StaticMetadataService is an auto generated Go binding around an Ethereum contract.
type StaticMetadataService struct {
	StaticMetadataServiceCaller     // Read-only binding to the contract
	StaticMetadataServiceTransactor // Write-only binding to the contract
	StaticMetadataServiceFilterer   // Log filterer for contract events
}

// StaticMetadataServiceCaller is an auto generated read-only Go binding around an Ethereum contract.
type StaticMetadataServiceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StaticMetadataServiceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StaticMetadataServiceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StaticMetadataServiceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StaticMetadataServiceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StaticMetadataServiceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StaticMetadataServiceSession struct {
	Contract     *StaticMetadataService // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// StaticMetadataServiceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StaticMetadataServiceCallerSession struct {
	Contract *StaticMetadataServiceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// StaticMetadataServiceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StaticMetadataServiceTransactorSession struct {
	Contract     *StaticMetadataServiceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// StaticMetadataServiceRaw is an auto generated low-level Go binding around an Ethereum contract.
type StaticMetadataServiceRaw struct {
	Contract *StaticMetadataService // Generic contract binding to access the raw methods on
}

// StaticMetadataServiceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StaticMetadataServiceCallerRaw struct {
	Contract *StaticMetadataServiceCaller // Generic read-only contract binding to access the raw methods on
}

// StaticMetadataServiceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StaticMetadataServiceTransactorRaw struct {
	Contract *StaticMetadataServiceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStaticMetadataService creates a new instance of StaticMetadataService, bound to a specific deployed contract.
func NewStaticMetadataService(address common.Address, backend bind.ContractBackend) (*StaticMetadataService, error) {
	contract, err := bindStaticMetadataService(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StaticMetadataService{StaticMetadataServiceCaller: StaticMetadataServiceCaller{contract: contract}, StaticMetadataServiceTransactor: StaticMetadataServiceTransactor{contract: contract}, StaticMetadataServiceFilterer: StaticMetadataServiceFilterer{contract: contract}}, nil
}

// NewStaticMetadataServiceCaller creates a new read-only instance of StaticMetadataService, bound to a specific deployed contract.
func NewStaticMetadataServiceCaller(address common.Address, caller bind.ContractCaller) (*StaticMetadataServiceCaller, error) {
	contract, err := bindStaticMetadataService(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StaticMetadataServiceCaller{contract: contract}, nil
}

// NewStaticMetadataServiceTransactor creates a new write-only instance of StaticMetadataService, bound to a specific deployed contract.
func NewStaticMetadataServiceTransactor(address common.Address, transactor bind.ContractTransactor) (*StaticMetadataServiceTransactor, error) {
	contract, err := bindStaticMetadataService(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StaticMetadataServiceTransactor{contract: contract}, nil
}

// NewStaticMetadataServiceFilterer creates a new log filterer instance of StaticMetadataService, bound to a specific deployed contract.
func NewStaticMetadataServiceFilterer(address common.Address, filterer bind.ContractFilterer) (*StaticMetadataServiceFilterer, error) {
	contract, err := bindStaticMetadataService(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StaticMetadataServiceFilterer{contract: contract}, nil
}

// bindStaticMetadataService binds a generic wrapper to an already deployed contract.
func bindStaticMetadataService(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StaticMetadataServiceMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StaticMetadataService *StaticMetadataServiceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StaticMetadataService.Contract.StaticMetadataServiceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StaticMetadataService *StaticMetadataServiceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StaticMetadataService.Contract.StaticMetadataServiceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StaticMetadataService *StaticMetadataServiceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StaticMetadataService.Contract.StaticMetadataServiceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StaticMetadataService *StaticMetadataServiceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StaticMetadataService.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StaticMetadataService *StaticMetadataServiceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StaticMetadataService.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StaticMetadataService *StaticMetadataServiceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StaticMetadataService.Contract.contract.Transact(opts, method, params...)
}

// Uri is a free data retrieval call binding the contract method 0x0e89341c.
//
// Solidity: function uri(uint256 ) view returns(string)
func (_StaticMetadataService *StaticMetadataServiceCaller) Uri(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var out []interface{}
	err := _StaticMetadataService.contract.Call(opts, &out, "uri", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Uri is a free data retrieval call binding the contract method 0x0e89341c.
//
// Solidity: function uri(uint256 ) view returns(string)
func (_StaticMetadataService *StaticMetadataServiceSession) Uri(arg0 *big.Int) (string, error) {
	return _StaticMetadataService.Contract.Uri(&_StaticMetadataService.CallOpts, arg0)
}

// Uri is a free data retrieval call binding the contract method 0x0e89341c.
//
// Solidity: function uri(uint256 ) view returns(string)
func (_StaticMetadataService *StaticMetadataServiceCallerSession) Uri(arg0 *big.Int) (string, error) {
	return _StaticMetadataService.Contract.Uri(&_StaticMetadataService.CallOpts, arg0)
}
