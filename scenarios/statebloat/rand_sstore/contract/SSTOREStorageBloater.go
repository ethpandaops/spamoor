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

// SSTOREStorageBloaterMetaData contains all meta data concerning the SSTOREStorageBloater contract.
var SSTOREStorageBloaterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"count\",\"type\":\"uint256\"}],\"name\":\"createSlots\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060f48061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063e3b393a414602d575b600080fd5b603c603836600460a7565b603e565b005b7f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffed816001430340421860005b8281101560a0577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff84838301098055600101606a565b5050505050565b60006020828403121560b7578081fd5b503591905056fea264697066735822122079e4ae597ac68792390217bfee59415f78d7e370132ab6590720d9f7898a50be64736f6c63430008000033",
}

// SSTOREStorageBloaterABI is the input ABI used to generate the binding from.
// Deprecated: Use SSTOREStorageBloaterMetaData.ABI instead.
var SSTOREStorageBloaterABI = SSTOREStorageBloaterMetaData.ABI

// SSTOREStorageBloaterBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SSTOREStorageBloaterMetaData.Bin instead.
var SSTOREStorageBloaterBin = SSTOREStorageBloaterMetaData.Bin

// DeploySSTOREStorageBloater deploys a new Ethereum contract, binding an instance of SSTOREStorageBloater to it.
func DeploySSTOREStorageBloater(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SSTOREStorageBloater, error) {
	parsed, err := SSTOREStorageBloaterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SSTOREStorageBloaterBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SSTOREStorageBloater{SSTOREStorageBloaterCaller: SSTOREStorageBloaterCaller{contract: contract}, SSTOREStorageBloaterTransactor: SSTOREStorageBloaterTransactor{contract: contract}, SSTOREStorageBloaterFilterer: SSTOREStorageBloaterFilterer{contract: contract}}, nil
}

// SSTOREStorageBloater is an auto generated Go binding around an Ethereum contract.
type SSTOREStorageBloater struct {
	SSTOREStorageBloaterCaller     // Read-only binding to the contract
	SSTOREStorageBloaterTransactor // Write-only binding to the contract
	SSTOREStorageBloaterFilterer   // Log filterer for contract events
}

// SSTOREStorageBloaterCaller is an auto generated read-only Go binding around an Ethereum contract.
type SSTOREStorageBloaterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SSTOREStorageBloaterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SSTOREStorageBloaterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SSTOREStorageBloaterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SSTOREStorageBloaterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SSTOREStorageBloaterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SSTOREStorageBloaterSession struct {
	Contract     *SSTOREStorageBloater // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// SSTOREStorageBloaterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SSTOREStorageBloaterCallerSession struct {
	Contract *SSTOREStorageBloaterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// SSTOREStorageBloaterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SSTOREStorageBloaterTransactorSession struct {
	Contract     *SSTOREStorageBloaterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// SSTOREStorageBloaterRaw is an auto generated low-level Go binding around an Ethereum contract.
type SSTOREStorageBloaterRaw struct {
	Contract *SSTOREStorageBloater // Generic contract binding to access the raw methods on
}

// SSTOREStorageBloaterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SSTOREStorageBloaterCallerRaw struct {
	Contract *SSTOREStorageBloaterCaller // Generic read-only contract binding to access the raw methods on
}

// SSTOREStorageBloaterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SSTOREStorageBloaterTransactorRaw struct {
	Contract *SSTOREStorageBloaterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSSTOREStorageBloater creates a new instance of SSTOREStorageBloater, bound to a specific deployed contract.
func NewSSTOREStorageBloater(address common.Address, backend bind.ContractBackend) (*SSTOREStorageBloater, error) {
	contract, err := bindSSTOREStorageBloater(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SSTOREStorageBloater{SSTOREStorageBloaterCaller: SSTOREStorageBloaterCaller{contract: contract}, SSTOREStorageBloaterTransactor: SSTOREStorageBloaterTransactor{contract: contract}, SSTOREStorageBloaterFilterer: SSTOREStorageBloaterFilterer{contract: contract}}, nil
}

// NewSSTOREStorageBloaterCaller creates a new read-only instance of SSTOREStorageBloater, bound to a specific deployed contract.
func NewSSTOREStorageBloaterCaller(address common.Address, caller bind.ContractCaller) (*SSTOREStorageBloaterCaller, error) {
	contract, err := bindSSTOREStorageBloater(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SSTOREStorageBloaterCaller{contract: contract}, nil
}

// NewSSTOREStorageBloaterTransactor creates a new write-only instance of SSTOREStorageBloater, bound to a specific deployed contract.
func NewSSTOREStorageBloaterTransactor(address common.Address, transactor bind.ContractTransactor) (*SSTOREStorageBloaterTransactor, error) {
	contract, err := bindSSTOREStorageBloater(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SSTOREStorageBloaterTransactor{contract: contract}, nil
}

// NewSSTOREStorageBloaterFilterer creates a new log filterer instance of SSTOREStorageBloater, bound to a specific deployed contract.
func NewSSTOREStorageBloaterFilterer(address common.Address, filterer bind.ContractFilterer) (*SSTOREStorageBloaterFilterer, error) {
	contract, err := bindSSTOREStorageBloater(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SSTOREStorageBloaterFilterer{contract: contract}, nil
}

// bindSSTOREStorageBloater binds a generic wrapper to an already deployed contract.
func bindSSTOREStorageBloater(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SSTOREStorageBloaterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SSTOREStorageBloater *SSTOREStorageBloaterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SSTOREStorageBloater.Contract.SSTOREStorageBloaterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SSTOREStorageBloater *SSTOREStorageBloaterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SSTOREStorageBloater.Contract.SSTOREStorageBloaterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SSTOREStorageBloater *SSTOREStorageBloaterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SSTOREStorageBloater.Contract.SSTOREStorageBloaterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SSTOREStorageBloater *SSTOREStorageBloaterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SSTOREStorageBloater.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SSTOREStorageBloater *SSTOREStorageBloaterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SSTOREStorageBloater.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SSTOREStorageBloater *SSTOREStorageBloaterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SSTOREStorageBloater.Contract.contract.Transact(opts, method, params...)
}

// CreateSlots is a paid mutator transaction binding the contract method 0xe3b393a4.
//
// Solidity: function createSlots(uint256 count) returns()
func (_SSTOREStorageBloater *SSTOREStorageBloaterTransactor) CreateSlots(opts *bind.TransactOpts, count *big.Int) (*types.Transaction, error) {
	return _SSTOREStorageBloater.contract.Transact(opts, "createSlots", count)
}

// CreateSlots is a paid mutator transaction binding the contract method 0xe3b393a4.
//
// Solidity: function createSlots(uint256 count) returns()
func (_SSTOREStorageBloater *SSTOREStorageBloaterSession) CreateSlots(count *big.Int) (*types.Transaction, error) {
	return _SSTOREStorageBloater.Contract.CreateSlots(&_SSTOREStorageBloater.TransactOpts, count)
}

// CreateSlots is a paid mutator transaction binding the contract method 0xe3b393a4.
//
// Solidity: function createSlots(uint256 count) returns()
func (_SSTOREStorageBloater *SSTOREStorageBloaterTransactorSession) CreateSlots(count *big.Int) (*types.Transaction, error) {
	return _SSTOREStorageBloater.Contract.CreateSlots(&_SSTOREStorageBloater.TransactOpts, count)
}
