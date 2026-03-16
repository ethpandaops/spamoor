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

// EIP2780HelperMetaData contains all meta data concerning the EIP2780Helper contract.
var EIP2780HelperMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"createAndDestroy\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"forwardValue\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nop\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b506102938061001c5f395ff3fe608060405260043610610033575f3560e01c80632e9c58bb1461003757806339859c06146100395780637aa397b91461004c575b5f5ffd5b005b61003761004736600461017a565b610054565b6100376100ec565b5f816001600160a01b0316346040515f6040518083038185875af1925050503d805f811461009d576040519150601f19603f3d011682016040523d82523d5f602084013e6100a2565b606091505b50509050806100e85760405162461bcd60e51b815260206004820152600e60248201526d199bdc9dd85c990819985a5b195960921b604482015260640160405180910390fd5b5050565b5f346040516100fa9061016e565b6040518091039082f0905080158015610115573d5f5f3e3d5ffd5b5060405162f55d9d60e01b81523360048201529091506001600160a01b0382169062f55d9d906024015f604051808303815f87803b158015610155575f5ffd5b505af1158015610167573d5f5f3e3d5ffd5b5050505050565b60b6806101a883390190565b5f6020828403121561018a575f5ffd5b81356001600160a01b03811681146101a0575f5ffd5b939250505056fe608060405260a780600f5f395ff3fe6080604052348015600e575f5ffd5b50600436106025575f3560e01c8062f55d9d146029575b5f5ffd5b603860343660046046565b603a565b005b806001600160a01b0316ff5b5f602082840312156055575f5ffd5b81356001600160a01b0381168114606a575f5ffd5b939250505056fea26469706673582212203b08093cc5ee33423820b46b59c88ad68c6824d3e21efbe0967d532a49da8a7064736f6c634300081d0033a2646970667358221220a8e3fa438451c54f8fde6d17de9da0396eec075c830484a301aca453c2bb241064736f6c634300081d0033",
}

// EIP2780HelperABI is the input ABI used to generate the binding from.
// Deprecated: Use EIP2780HelperMetaData.ABI instead.
var EIP2780HelperABI = EIP2780HelperMetaData.ABI

// EIP2780HelperBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use EIP2780HelperMetaData.Bin instead.
var EIP2780HelperBin = EIP2780HelperMetaData.Bin

// DeployEIP2780Helper deploys a new Ethereum contract, binding an instance of EIP2780Helper to it.
func DeployEIP2780Helper(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *EIP2780Helper, error) {
	parsed, err := EIP2780HelperMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(EIP2780HelperBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &EIP2780Helper{EIP2780HelperCaller: EIP2780HelperCaller{contract: contract}, EIP2780HelperTransactor: EIP2780HelperTransactor{contract: contract}, EIP2780HelperFilterer: EIP2780HelperFilterer{contract: contract}}, nil
}

// EIP2780Helper is an auto generated Go binding around an Ethereum contract.
type EIP2780Helper struct {
	EIP2780HelperCaller     // Read-only binding to the contract
	EIP2780HelperTransactor // Write-only binding to the contract
	EIP2780HelperFilterer   // Log filterer for contract events
}

// EIP2780HelperCaller is an auto generated read-only Go binding around an Ethereum contract.
type EIP2780HelperCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EIP2780HelperTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EIP2780HelperTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EIP2780HelperFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EIP2780HelperFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EIP2780HelperSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EIP2780HelperSession struct {
	Contract     *EIP2780Helper    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EIP2780HelperCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EIP2780HelperCallerSession struct {
	Contract *EIP2780HelperCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// EIP2780HelperTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EIP2780HelperTransactorSession struct {
	Contract     *EIP2780HelperTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// EIP2780HelperRaw is an auto generated low-level Go binding around an Ethereum contract.
type EIP2780HelperRaw struct {
	Contract *EIP2780Helper // Generic contract binding to access the raw methods on
}

// EIP2780HelperCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EIP2780HelperCallerRaw struct {
	Contract *EIP2780HelperCaller // Generic read-only contract binding to access the raw methods on
}

// EIP2780HelperTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EIP2780HelperTransactorRaw struct {
	Contract *EIP2780HelperTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEIP2780Helper creates a new instance of EIP2780Helper, bound to a specific deployed contract.
func NewEIP2780Helper(address common.Address, backend bind.ContractBackend) (*EIP2780Helper, error) {
	contract, err := bindEIP2780Helper(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EIP2780Helper{EIP2780HelperCaller: EIP2780HelperCaller{contract: contract}, EIP2780HelperTransactor: EIP2780HelperTransactor{contract: contract}, EIP2780HelperFilterer: EIP2780HelperFilterer{contract: contract}}, nil
}

// NewEIP2780HelperCaller creates a new read-only instance of EIP2780Helper, bound to a specific deployed contract.
func NewEIP2780HelperCaller(address common.Address, caller bind.ContractCaller) (*EIP2780HelperCaller, error) {
	contract, err := bindEIP2780Helper(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EIP2780HelperCaller{contract: contract}, nil
}

// NewEIP2780HelperTransactor creates a new write-only instance of EIP2780Helper, bound to a specific deployed contract.
func NewEIP2780HelperTransactor(address common.Address, transactor bind.ContractTransactor) (*EIP2780HelperTransactor, error) {
	contract, err := bindEIP2780Helper(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EIP2780HelperTransactor{contract: contract}, nil
}

// NewEIP2780HelperFilterer creates a new log filterer instance of EIP2780Helper, bound to a specific deployed contract.
func NewEIP2780HelperFilterer(address common.Address, filterer bind.ContractFilterer) (*EIP2780HelperFilterer, error) {
	contract, err := bindEIP2780Helper(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EIP2780HelperFilterer{contract: contract}, nil
}

// bindEIP2780Helper binds a generic wrapper to an already deployed contract.
func bindEIP2780Helper(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := EIP2780HelperMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EIP2780Helper *EIP2780HelperRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EIP2780Helper.Contract.EIP2780HelperCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EIP2780Helper *EIP2780HelperRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EIP2780Helper.Contract.EIP2780HelperTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EIP2780Helper *EIP2780HelperRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EIP2780Helper.Contract.EIP2780HelperTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EIP2780Helper *EIP2780HelperCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EIP2780Helper.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EIP2780Helper *EIP2780HelperTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EIP2780Helper.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EIP2780Helper *EIP2780HelperTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EIP2780Helper.Contract.contract.Transact(opts, method, params...)
}

// CreateAndDestroy is a paid mutator transaction binding the contract method 0x7aa397b9.
//
// Solidity: function createAndDestroy() payable returns()
func (_EIP2780Helper *EIP2780HelperTransactor) CreateAndDestroy(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EIP2780Helper.contract.Transact(opts, "createAndDestroy")
}

// CreateAndDestroy is a paid mutator transaction binding the contract method 0x7aa397b9.
//
// Solidity: function createAndDestroy() payable returns()
func (_EIP2780Helper *EIP2780HelperSession) CreateAndDestroy() (*types.Transaction, error) {
	return _EIP2780Helper.Contract.CreateAndDestroy(&_EIP2780Helper.TransactOpts)
}

// CreateAndDestroy is a paid mutator transaction binding the contract method 0x7aa397b9.
//
// Solidity: function createAndDestroy() payable returns()
func (_EIP2780Helper *EIP2780HelperTransactorSession) CreateAndDestroy() (*types.Transaction, error) {
	return _EIP2780Helper.Contract.CreateAndDestroy(&_EIP2780Helper.TransactOpts)
}

// ForwardValue is a paid mutator transaction binding the contract method 0x39859c06.
//
// Solidity: function forwardValue(address target) payable returns()
func (_EIP2780Helper *EIP2780HelperTransactor) ForwardValue(opts *bind.TransactOpts, target common.Address) (*types.Transaction, error) {
	return _EIP2780Helper.contract.Transact(opts, "forwardValue", target)
}

// ForwardValue is a paid mutator transaction binding the contract method 0x39859c06.
//
// Solidity: function forwardValue(address target) payable returns()
func (_EIP2780Helper *EIP2780HelperSession) ForwardValue(target common.Address) (*types.Transaction, error) {
	return _EIP2780Helper.Contract.ForwardValue(&_EIP2780Helper.TransactOpts, target)
}

// ForwardValue is a paid mutator transaction binding the contract method 0x39859c06.
//
// Solidity: function forwardValue(address target) payable returns()
func (_EIP2780Helper *EIP2780HelperTransactorSession) ForwardValue(target common.Address) (*types.Transaction, error) {
	return _EIP2780Helper.Contract.ForwardValue(&_EIP2780Helper.TransactOpts, target)
}

// Nop is a paid mutator transaction binding the contract method 0x2e9c58bb.
//
// Solidity: function nop() payable returns()
func (_EIP2780Helper *EIP2780HelperTransactor) Nop(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EIP2780Helper.contract.Transact(opts, "nop")
}

// Nop is a paid mutator transaction binding the contract method 0x2e9c58bb.
//
// Solidity: function nop() payable returns()
func (_EIP2780Helper *EIP2780HelperSession) Nop() (*types.Transaction, error) {
	return _EIP2780Helper.Contract.Nop(&_EIP2780Helper.TransactOpts)
}

// Nop is a paid mutator transaction binding the contract method 0x2e9c58bb.
//
// Solidity: function nop() payable returns()
func (_EIP2780Helper *EIP2780HelperTransactorSession) Nop() (*types.Transaction, error) {
	return _EIP2780Helper.Contract.Nop(&_EIP2780Helper.TransactOpts)
}
