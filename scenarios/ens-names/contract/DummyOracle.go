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

// DummyOracleMetaData contains all meta data concerning the DummyOracle contract.
var DummyOracleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"_value\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"_value\",\"type\":\"int256\"}],\"name\":\"set\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600f57600080fd5b5060405161010b38038061010b833981016040819052602c916039565b603481600055565b506051565b600060208284031215604a57600080fd5b5051919050565b60ac8061005f6000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c806350d25bcd146037578063e5c19b2d14604c575b600080fd5b60005460405190815260200160405180910390f35b605c6057366004605e565b600055565b005b600060208284031215606f57600080fd5b503591905056fea26469706673582212209978ff9a8fe4b98ec4229c64d331a421fefeb2a9250de66692c36b0c33a0d91364736f6c634300081a0033",
}

// DummyOracleABI is the input ABI used to generate the binding from.
// Deprecated: Use DummyOracleMetaData.ABI instead.
var DummyOracleABI = DummyOracleMetaData.ABI

// DummyOracleBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DummyOracleMetaData.Bin instead.
var DummyOracleBin = DummyOracleMetaData.Bin

// DeployDummyOracle deploys a new Ethereum contract, binding an instance of DummyOracle to it.
func DeployDummyOracle(auth *bind.TransactOpts, backend bind.ContractBackend, _value *big.Int) (common.Address, *types.Transaction, *DummyOracle, error) {
	parsed, err := DummyOracleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DummyOracleBin), backend, _value)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DummyOracle{DummyOracleCaller: DummyOracleCaller{contract: contract}, DummyOracleTransactor: DummyOracleTransactor{contract: contract}, DummyOracleFilterer: DummyOracleFilterer{contract: contract}}, nil
}

// DummyOracle is an auto generated Go binding around an Ethereum contract.
type DummyOracle struct {
	DummyOracleCaller     // Read-only binding to the contract
	DummyOracleTransactor // Write-only binding to the contract
	DummyOracleFilterer   // Log filterer for contract events
}

// DummyOracleCaller is an auto generated read-only Go binding around an Ethereum contract.
type DummyOracleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DummyOracleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DummyOracleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DummyOracleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DummyOracleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DummyOracleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DummyOracleSession struct {
	Contract     *DummyOracle      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DummyOracleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DummyOracleCallerSession struct {
	Contract *DummyOracleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// DummyOracleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DummyOracleTransactorSession struct {
	Contract     *DummyOracleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// DummyOracleRaw is an auto generated low-level Go binding around an Ethereum contract.
type DummyOracleRaw struct {
	Contract *DummyOracle // Generic contract binding to access the raw methods on
}

// DummyOracleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DummyOracleCallerRaw struct {
	Contract *DummyOracleCaller // Generic read-only contract binding to access the raw methods on
}

// DummyOracleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DummyOracleTransactorRaw struct {
	Contract *DummyOracleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDummyOracle creates a new instance of DummyOracle, bound to a specific deployed contract.
func NewDummyOracle(address common.Address, backend bind.ContractBackend) (*DummyOracle, error) {
	contract, err := bindDummyOracle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DummyOracle{DummyOracleCaller: DummyOracleCaller{contract: contract}, DummyOracleTransactor: DummyOracleTransactor{contract: contract}, DummyOracleFilterer: DummyOracleFilterer{contract: contract}}, nil
}

// NewDummyOracleCaller creates a new read-only instance of DummyOracle, bound to a specific deployed contract.
func NewDummyOracleCaller(address common.Address, caller bind.ContractCaller) (*DummyOracleCaller, error) {
	contract, err := bindDummyOracle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DummyOracleCaller{contract: contract}, nil
}

// NewDummyOracleTransactor creates a new write-only instance of DummyOracle, bound to a specific deployed contract.
func NewDummyOracleTransactor(address common.Address, transactor bind.ContractTransactor) (*DummyOracleTransactor, error) {
	contract, err := bindDummyOracle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DummyOracleTransactor{contract: contract}, nil
}

// NewDummyOracleFilterer creates a new log filterer instance of DummyOracle, bound to a specific deployed contract.
func NewDummyOracleFilterer(address common.Address, filterer bind.ContractFilterer) (*DummyOracleFilterer, error) {
	contract, err := bindDummyOracle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DummyOracleFilterer{contract: contract}, nil
}

// bindDummyOracle binds a generic wrapper to an already deployed contract.
func bindDummyOracle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DummyOracleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DummyOracle *DummyOracleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DummyOracle.Contract.DummyOracleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DummyOracle *DummyOracleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DummyOracle.Contract.DummyOracleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DummyOracle *DummyOracleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DummyOracle.Contract.DummyOracleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DummyOracle *DummyOracleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DummyOracle.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DummyOracle *DummyOracleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DummyOracle.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DummyOracle *DummyOracleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DummyOracle.Contract.contract.Transact(opts, method, params...)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_DummyOracle *DummyOracleCaller) LatestAnswer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DummyOracle.contract.Call(opts, &out, "latestAnswer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_DummyOracle *DummyOracleSession) LatestAnswer() (*big.Int, error) {
	return _DummyOracle.Contract.LatestAnswer(&_DummyOracle.CallOpts)
}

// LatestAnswer is a free data retrieval call binding the contract method 0x50d25bcd.
//
// Solidity: function latestAnswer() view returns(int256)
func (_DummyOracle *DummyOracleCallerSession) LatestAnswer() (*big.Int, error) {
	return _DummyOracle.Contract.LatestAnswer(&_DummyOracle.CallOpts)
}

// Set is a paid mutator transaction binding the contract method 0xe5c19b2d.
//
// Solidity: function set(int256 _value) returns()
func (_DummyOracle *DummyOracleTransactor) Set(opts *bind.TransactOpts, _value *big.Int) (*types.Transaction, error) {
	return _DummyOracle.contract.Transact(opts, "set", _value)
}

// Set is a paid mutator transaction binding the contract method 0xe5c19b2d.
//
// Solidity: function set(int256 _value) returns()
func (_DummyOracle *DummyOracleSession) Set(_value *big.Int) (*types.Transaction, error) {
	return _DummyOracle.Contract.Set(&_DummyOracle.TransactOpts, _value)
}

// Set is a paid mutator transaction binding the contract method 0xe5c19b2d.
//
// Solidity: function set(int256 _value) returns()
func (_DummyOracle *DummyOracleTransactorSession) Set(_value *big.Int) (*types.Transaction, error) {
	return _DummyOracle.Contract.Set(&_DummyOracle.TransactOpts, _value)
}
