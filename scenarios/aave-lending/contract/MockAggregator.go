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

// MockAggregatorMetaData contains all meta data concerning the MockAggregator contract.
var MockAggregatorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"_answer\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"answer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"_answer\",\"type\":\"int256\"}],\"name\":\"setAnswer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060405161011d38038061011d83398101604081905261002f91610037565b600055610050565b60006020828403121561004957600080fd5b5051919050565b60bf8061005e6000396000f3fe6080604052348015600f57600080fd5b5060043610603c5760003560e01c806350d25bcd14604157806385bb7d6914605757806399213cd814605f575b600080fd5b6000545b60405190815260200160405180910390f35b604560005481565b606f606a3660046071565b600055565b005b600060208284031215608257600080fd5b503591905056fea26469706673582212203334cde6ec31eb07a8557a3ba707d284e76bebdb610a665d75f8b9ba6b4af8c364736f6c63430008110033",
}

// MockAggregatorABI is the input ABI used to generate the binding from.
// Deprecated: Use MockAggregatorMetaData.ABI instead.
var MockAggregatorABI = MockAggregatorMetaData.ABI

// MockAggregatorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MockAggregatorMetaData.Bin instead.
var MockAggregatorBin = MockAggregatorMetaData.Bin

// DeployMockAggregator deploys a new Ethereum contract, binding an instance of MockAggregator to it.
func DeployMockAggregator(auth *bind.TransactOpts, backend bind.ContractBackend, _answer *big.Int) (common.Address, *types.Transaction, *MockAggregator, error) {
	parsed, err := MockAggregatorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MockAggregatorBin), backend, _answer)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MockAggregator{MockAggregatorCaller: MockAggregatorCaller{contract: contract}, MockAggregatorTransactor: MockAggregatorTransactor{contract: contract}, MockAggregatorFilterer: MockAggregatorFilterer{contract: contract}}, nil
}

// MockAggregator is an auto generated Go binding around an Ethereum contract.
type MockAggregator struct {
	MockAggregatorCaller     // Read-only binding to the contract
	MockAggregatorTransactor // Write-only binding to the contract
	MockAggregatorFilterer   // Log filterer for contract events
}

// MockAggregatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type MockAggregatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockAggregatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MockAggregatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockAggregatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MockAggregatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockAggregatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MockAggregatorSession struct {
	Contract     *MockAggregator   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MockAggregatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MockAggregatorCallerSession struct {
	Contract *MockAggregatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// MockAggregatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MockAggregatorTransactorSession struct {
	Contract     *MockAggregatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// MockAggregatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type MockAggregatorRaw struct {
	Contract *MockAggregator // Generic contract binding to access the raw methods on
}

// MockAggregatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MockAggregatorCallerRaw struct {
	Contract *MockAggregatorCaller // Generic read-only contract binding to access the raw methods on
}

// MockAggregatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MockAggregatorTransactorRaw struct {
	Contract *MockAggregatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMockAggregator creates a new instance of MockAggregator, bound to a specific deployed contract.
func NewMockAggregator(address common.Address, backend bind.ContractBackend) (*MockAggregator, error) {
	contract, err := bindMockAggregator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MockAggregator{MockAggregatorCaller: MockAggregatorCaller{contract: contract}, MockAggregatorTransactor: MockAggregatorTransactor{contract: contract}, MockAggregatorFilterer: MockAggregatorFilterer{contract: contract}}, nil
}

// NewMockAggregatorCaller creates a new read-only instance of MockAggregator, bound to a specific deployed contract.
func NewMockAggregatorCaller(address common.Address, caller bind.ContractCaller) (*MockAggregatorCaller, error) {
	contract, err := bindMockAggregator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MockAggregatorCaller{contract: contract}, nil
}

// NewMockAggregatorTransactor creates a new write-only instance of MockAggregator, bound to a specific deployed contract.
func NewMockAggregatorTransactor(address common.Address, transactor bind.ContractTransactor) (*MockAggregatorTransactor, error) {
	contract, err := bindMockAggregator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MockAggregatorTransactor{contract: contract}, nil
}

// NewMockAggregatorFilterer creates a new log filterer instance of MockAggregator, bound to a specific deployed contract.
func NewMockAggregatorFilterer(address common.Address, filterer bind.ContractFilterer) (*MockAggregatorFilterer, error) {
	contract, err := bindMockAggregator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MockAggregatorFilterer{contract: contract}, nil
}

// bindMockAggregator binds a generic wrapper to an already deployed contract.
func bindMockAggregator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MockAggregatorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MockAggregator *MockAggregatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockAggregator.Contract.MockAggregatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MockAggregator *MockAggregatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockAggregator.Contract.MockAggregatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MockAggregator *MockAggregatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockAggregator.Contract.MockAggregatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MockAggregator *MockAggregatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockAggregator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MockAggregator *MockAggregatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockAggregator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MockAggregator *MockAggregatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockAggregator.Contract.contract.Transact(opts, method, params...)
}

// Answer is a free data retrieval call binding the contract method 0x85bb7d69.
//
// Solidity: function answer() view returns(int256)
func (_MockAggregator *MockAggregatorCaller) Answer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MockAggregator.contract.Call(opts, &out, "answer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Answer is a free data retrieval call binding the contract method 0x85bb7d69.
//
// Solidity: function answer() view returns(int256)
func (_MockAggregator *MockAggregatorSession) Answer() (*big.Int, error) {
	return _MockAggregator.Contract.Answer(&_MockAggregator.CallOpts)
}

// Answer is a free data retrieval call binding the contract method 0x85bb7d69.
//
// Solidity: function answer() view returns(int256)
func (_MockAggregator *MockAggregatorCallerSession) Answer() (*big.Int, error) {
	return _MockAggregator.Contract.Answer(&_MockAggregator.CallOpts)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_MockAggregator *MockAggregatorCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MockAggregator.contract.Call(opts, &out, "latestAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_MockAggregator *MockAggregatorSession) LatestAnswer() (*big.Int, error) {
	return _MockAggregator.Contract.LatestAnswer(&_MockAggregator.CallOpts)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_MockAggregator *MockAggregatorCallerSession) LatestAnswer() (*big.Int, error) {
	return _MockAggregator.Contract.LatestAnswer(&_MockAggregator.CallOpts)
}

// SetAnswer is a paid mutator transaction binding the contract method 0x99213cd8.
//
// Solidity: function setAnswer(int256 _answer) returns()
func (_MockAggregator *MockAggregatorTransactor) SetAnswer(opts *bind.TransactOpts, _answer *big.Int) (*types.Transaction, error) {
	return _MockAggregator.contract.Transact(opts, "setAnswer", _answer)
}

// SetAnswer is a paid mutator transaction binding the contract method 0x99213cd8.
//
// Solidity: function setAnswer(int256 _answer) returns()
func (_MockAggregator *MockAggregatorSession) SetAnswer(_answer *big.Int) (*types.Transaction, error) {
	return _MockAggregator.Contract.SetAnswer(&_MockAggregator.TransactOpts, _answer)
}

// SetAnswer is a paid mutator transaction binding the contract method 0x99213cd8.
//
// Solidity: function setAnswer(int256 _answer) returns()
func (_MockAggregator *MockAggregatorTransactorSession) SetAnswer(_answer *big.Int) (*types.Transaction, error) {
	return _MockAggregator.Contract.SetAnswer(&_MockAggregator.TransactOpts, _answer)
}
