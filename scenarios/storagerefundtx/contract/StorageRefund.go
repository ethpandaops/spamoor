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

// StorageRefundMetaData contains all meta data concerning the StorageRefund contract.
var StorageRefundMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"clearPointer\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"clearableUpTo\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"slotsPerCall\",\"type\":\"uint256\"}],\"name\":\"execute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastBlockNumber\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"storageSlots\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"writePointer\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561000f575f80fd5b506102148061001d5f395ff3fe608060405234801561000f575f80fd5b5060043610610060575f3560e01c80632552317c1461006457806351bbd3581461007f5780635387694b1461008857806375635015146100a75780638c60b2ee146100af578063fe0d94c1146100b8575b5f80fd5b61006d60035481565b60405190815260200160405180910390f35b61006d60025481565b61006d610096366004610187565b60046020525f908152604090205481565b61006d5f5481565b61006d60015481565b6100cb6100c6366004610187565b6100cd565b005b6003544311156100e1575f54600255436003555b5f6001546002546100f291906101b2565b9050818111156100ff5750805b6001545f5b828110156101335760045f61011983856101cb565b815260208101919091526040015f90812055600101610104565b5061013e82826101cb565b6001555f8054905b84811015610174574360045f61015c84866101cb565b815260208101919091526040015f2055600101610146565b5061017f84826101cb565b5f5550505050565b5f60208284031215610197575f80fd5b5035919050565b634e487b7160e01b5f52601160045260245ffd5b818103818111156101c5576101c561019e565b92915050565b808201808211156101c5576101c561019e56fea2646970667358221220c918dd9b1ba607b3ceaac63b561025c13f1fe54a55ec0e6c3401467ce28df6b664736f6c63430008160033",
}

// StorageRefundABI is the input ABI used to generate the binding from.
// Deprecated: Use StorageRefundMetaData.ABI instead.
var StorageRefundABI = StorageRefundMetaData.ABI

// StorageRefundBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use StorageRefundMetaData.Bin instead.
var StorageRefundBin = StorageRefundMetaData.Bin

// DeployStorageRefund deploys a new Ethereum contract, binding an instance of StorageRefund to it.
func DeployStorageRefund(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *StorageRefund, error) {
	parsed, err := StorageRefundMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StorageRefundBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &StorageRefund{StorageRefundCaller: StorageRefundCaller{contract: contract}, StorageRefundTransactor: StorageRefundTransactor{contract: contract}, StorageRefundFilterer: StorageRefundFilterer{contract: contract}}, nil
}

// StorageRefund is an auto generated Go binding around an Ethereum contract.
type StorageRefund struct {
	StorageRefundCaller     // Read-only binding to the contract
	StorageRefundTransactor // Write-only binding to the contract
	StorageRefundFilterer   // Log filterer for contract events
}

// StorageRefundCaller is an auto generated read-only Go binding around an Ethereum contract.
type StorageRefundCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageRefundTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StorageRefundTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageRefundFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StorageRefundFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageRefundSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StorageRefundSession struct {
	Contract     *StorageRefund    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StorageRefundCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StorageRefundCallerSession struct {
	Contract *StorageRefundCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// StorageRefundTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StorageRefundTransactorSession struct {
	Contract     *StorageRefundTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// StorageRefundRaw is an auto generated low-level Go binding around an Ethereum contract.
type StorageRefundRaw struct {
	Contract *StorageRefund // Generic contract binding to access the raw methods on
}

// StorageRefundCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StorageRefundCallerRaw struct {
	Contract *StorageRefundCaller // Generic read-only contract binding to access the raw methods on
}

// StorageRefundTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StorageRefundTransactorRaw struct {
	Contract *StorageRefundTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStorageRefund creates a new instance of StorageRefund, bound to a specific deployed contract.
func NewStorageRefund(address common.Address, backend bind.ContractBackend) (*StorageRefund, error) {
	contract, err := bindStorageRefund(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StorageRefund{StorageRefundCaller: StorageRefundCaller{contract: contract}, StorageRefundTransactor: StorageRefundTransactor{contract: contract}, StorageRefundFilterer: StorageRefundFilterer{contract: contract}}, nil
}

// NewStorageRefundCaller creates a new read-only instance of StorageRefund, bound to a specific deployed contract.
func NewStorageRefundCaller(address common.Address, caller bind.ContractCaller) (*StorageRefundCaller, error) {
	contract, err := bindStorageRefund(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StorageRefundCaller{contract: contract}, nil
}

// NewStorageRefundTransactor creates a new write-only instance of StorageRefund, bound to a specific deployed contract.
func NewStorageRefundTransactor(address common.Address, transactor bind.ContractTransactor) (*StorageRefundTransactor, error) {
	contract, err := bindStorageRefund(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StorageRefundTransactor{contract: contract}, nil
}

// NewStorageRefundFilterer creates a new log filterer instance of StorageRefund, bound to a specific deployed contract.
func NewStorageRefundFilterer(address common.Address, filterer bind.ContractFilterer) (*StorageRefundFilterer, error) {
	contract, err := bindStorageRefund(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StorageRefundFilterer{contract: contract}, nil
}

// bindStorageRefund binds a generic wrapper to an already deployed contract.
func bindStorageRefund(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StorageRefundMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StorageRefund *StorageRefundRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StorageRefund.Contract.StorageRefundCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StorageRefund *StorageRefundRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StorageRefund.Contract.StorageRefundTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StorageRefund *StorageRefundRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StorageRefund.Contract.StorageRefundTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StorageRefund *StorageRefundCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StorageRefund.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StorageRefund *StorageRefundTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StorageRefund.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StorageRefund *StorageRefundTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StorageRefund.Contract.contract.Transact(opts, method, params...)
}

// ClearPointer is a free data retrieval call binding the contract method 0x8c60b2ee.
//
// Solidity: function clearPointer() view returns(uint256)
func (_StorageRefund *StorageRefundCaller) ClearPointer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StorageRefund.contract.Call(opts, &out, "clearPointer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ClearPointer is a free data retrieval call binding the contract method 0x8c60b2ee.
//
// Solidity: function clearPointer() view returns(uint256)
func (_StorageRefund *StorageRefundSession) ClearPointer() (*big.Int, error) {
	return _StorageRefund.Contract.ClearPointer(&_StorageRefund.CallOpts)
}

// ClearPointer is a free data retrieval call binding the contract method 0x8c60b2ee.
//
// Solidity: function clearPointer() view returns(uint256)
func (_StorageRefund *StorageRefundCallerSession) ClearPointer() (*big.Int, error) {
	return _StorageRefund.Contract.ClearPointer(&_StorageRefund.CallOpts)
}

// ClearableUpTo is a free data retrieval call binding the contract method 0x51bbd358.
//
// Solidity: function clearableUpTo() view returns(uint256)
func (_StorageRefund *StorageRefundCaller) ClearableUpTo(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StorageRefund.contract.Call(opts, &out, "clearableUpTo")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ClearableUpTo is a free data retrieval call binding the contract method 0x51bbd358.
//
// Solidity: function clearableUpTo() view returns(uint256)
func (_StorageRefund *StorageRefundSession) ClearableUpTo() (*big.Int, error) {
	return _StorageRefund.Contract.ClearableUpTo(&_StorageRefund.CallOpts)
}

// ClearableUpTo is a free data retrieval call binding the contract method 0x51bbd358.
//
// Solidity: function clearableUpTo() view returns(uint256)
func (_StorageRefund *StorageRefundCallerSession) ClearableUpTo() (*big.Int, error) {
	return _StorageRefund.Contract.ClearableUpTo(&_StorageRefund.CallOpts)
}

// LastBlockNumber is a free data retrieval call binding the contract method 0x2552317c.
//
// Solidity: function lastBlockNumber() view returns(uint256)
func (_StorageRefund *StorageRefundCaller) LastBlockNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StorageRefund.contract.Call(opts, &out, "lastBlockNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastBlockNumber is a free data retrieval call binding the contract method 0x2552317c.
//
// Solidity: function lastBlockNumber() view returns(uint256)
func (_StorageRefund *StorageRefundSession) LastBlockNumber() (*big.Int, error) {
	return _StorageRefund.Contract.LastBlockNumber(&_StorageRefund.CallOpts)
}

// LastBlockNumber is a free data retrieval call binding the contract method 0x2552317c.
//
// Solidity: function lastBlockNumber() view returns(uint256)
func (_StorageRefund *StorageRefundCallerSession) LastBlockNumber() (*big.Int, error) {
	return _StorageRefund.Contract.LastBlockNumber(&_StorageRefund.CallOpts)
}

// StorageSlots is a free data retrieval call binding the contract method 0x5387694b.
//
// Solidity: function storageSlots(uint256 ) view returns(uint256)
func (_StorageRefund *StorageRefundCaller) StorageSlots(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _StorageRefund.contract.Call(opts, &out, "storageSlots", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StorageSlots is a free data retrieval call binding the contract method 0x5387694b.
//
// Solidity: function storageSlots(uint256 ) view returns(uint256)
func (_StorageRefund *StorageRefundSession) StorageSlots(arg0 *big.Int) (*big.Int, error) {
	return _StorageRefund.Contract.StorageSlots(&_StorageRefund.CallOpts, arg0)
}

// StorageSlots is a free data retrieval call binding the contract method 0x5387694b.
//
// Solidity: function storageSlots(uint256 ) view returns(uint256)
func (_StorageRefund *StorageRefundCallerSession) StorageSlots(arg0 *big.Int) (*big.Int, error) {
	return _StorageRefund.Contract.StorageSlots(&_StorageRefund.CallOpts, arg0)
}

// WritePointer is a free data retrieval call binding the contract method 0x75635015.
//
// Solidity: function writePointer() view returns(uint256)
func (_StorageRefund *StorageRefundCaller) WritePointer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StorageRefund.contract.Call(opts, &out, "writePointer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WritePointer is a free data retrieval call binding the contract method 0x75635015.
//
// Solidity: function writePointer() view returns(uint256)
func (_StorageRefund *StorageRefundSession) WritePointer() (*big.Int, error) {
	return _StorageRefund.Contract.WritePointer(&_StorageRefund.CallOpts)
}

// WritePointer is a free data retrieval call binding the contract method 0x75635015.
//
// Solidity: function writePointer() view returns(uint256)
func (_StorageRefund *StorageRefundCallerSession) WritePointer() (*big.Int, error) {
	return _StorageRefund.Contract.WritePointer(&_StorageRefund.CallOpts)
}

// Execute is a paid mutator transaction binding the contract method 0xfe0d94c1.
//
// Solidity: function execute(uint256 slotsPerCall) returns()
func (_StorageRefund *StorageRefundTransactor) Execute(opts *bind.TransactOpts, slotsPerCall *big.Int) (*types.Transaction, error) {
	return _StorageRefund.contract.Transact(opts, "execute", slotsPerCall)
}

// Execute is a paid mutator transaction binding the contract method 0xfe0d94c1.
//
// Solidity: function execute(uint256 slotsPerCall) returns()
func (_StorageRefund *StorageRefundSession) Execute(slotsPerCall *big.Int) (*types.Transaction, error) {
	return _StorageRefund.Contract.Execute(&_StorageRefund.TransactOpts, slotsPerCall)
}

// Execute is a paid mutator transaction binding the contract method 0xfe0d94c1.
//
// Solidity: function execute(uint256 slotsPerCall) returns()
func (_StorageRefund *StorageRefundTransactorSession) Execute(slotsPerCall *big.Int) (*types.Transaction, error) {
	return _StorageRefund.Contract.Execute(&_StorageRefund.TransactOpts, slotsPerCall)
}
