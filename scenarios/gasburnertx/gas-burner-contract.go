package gasburnertx

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

// GasBurnerMetaData contains all meta data concerning the GasBurner contract.
var GasBurnerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"burn1000k\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"burn100k\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"burn1500k\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"burn2000k\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"burn500k\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"burnGasUnits\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"}]",
	Bin: "0x608060405234801561001057600080fd5b506101ed806100206000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c80631ff94f73146100675780634158735914610071578063419f772e14610084578063a2cea6351461008c578063d0da1e3514610094578063e37937701461009c575b600080fd5b61006f6100a4565b005b61006f61007f366004610156565b6100c6565b61006f6100e6565b61006f6100f2565b61006f6100fe565b61006f61010a565b6100b06207a120610112565b6000805490806100bf83610185565b9190505550565b6100cf81610112565b6000805490806100de83610185565b919050555050565b6100b0620186a0610112565b6100b06216e360610112565b6100b0621e8480610112565b6100b0620f42405b60005a905081811161012757506103e8610134565b610131828261019e565b90505b60005b815a1115610151578061014981610185565b915050610137565b505050565b60006020828403121561016857600080fd5b5035919050565b634e487b7160e01b600052601160045260246000fd5b6000600182016101975761019761016f565b5060010190565b818103818111156101b1576101b161016f565b9291505056fea2646970667358221220270149b77953b8ff91111fb3cffd74ceeccffbc9594d46ada1a9a7c08b00575664736f6c63430008180033",
}

// GasBurnerABI is the input ABI used to generate the binding from.
// Deprecated: Use GasBurnerMetaData.ABI instead.
var GasBurnerABI = GasBurnerMetaData.ABI

// GasBurnerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use GasBurnerMetaData.Bin instead.
var GasBurnerBin = GasBurnerMetaData.Bin

// DeployGasBurner deploys a new Ethereum contract, binding an instance of GasBurner to it.
func DeployGasBurner(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *GasBurner, error) {
	parsed, err := GasBurnerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(GasBurnerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &GasBurner{GasBurnerCaller: GasBurnerCaller{contract: contract}, GasBurnerTransactor: GasBurnerTransactor{contract: contract}, GasBurnerFilterer: GasBurnerFilterer{contract: contract}}, nil
}

// GasBurner is an auto generated Go binding around an Ethereum contract.
type GasBurner struct {
	GasBurnerCaller     // Read-only binding to the contract
	GasBurnerTransactor // Write-only binding to the contract
	GasBurnerFilterer   // Log filterer for contract events
}

// GasBurnerCaller is an auto generated read-only Go binding around an Ethereum contract.
type GasBurnerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GasBurnerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GasBurnerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GasBurnerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GasBurnerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GasBurnerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GasBurnerSession struct {
	Contract     *GasBurner        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GasBurnerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GasBurnerCallerSession struct {
	Contract *GasBurnerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// GasBurnerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GasBurnerTransactorSession struct {
	Contract     *GasBurnerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// GasBurnerRaw is an auto generated low-level Go binding around an Ethereum contract.
type GasBurnerRaw struct {
	Contract *GasBurner // Generic contract binding to access the raw methods on
}

// GasBurnerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GasBurnerCallerRaw struct {
	Contract *GasBurnerCaller // Generic read-only contract binding to access the raw methods on
}

// GasBurnerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GasBurnerTransactorRaw struct {
	Contract *GasBurnerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGasBurner creates a new instance of GasBurner, bound to a specific deployed contract.
func NewGasBurner(address common.Address, backend bind.ContractBackend) (*GasBurner, error) {
	contract, err := bindGasBurner(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GasBurner{GasBurnerCaller: GasBurnerCaller{contract: contract}, GasBurnerTransactor: GasBurnerTransactor{contract: contract}, GasBurnerFilterer: GasBurnerFilterer{contract: contract}}, nil
}

// NewGasBurnerCaller creates a new read-only instance of GasBurner, bound to a specific deployed contract.
func NewGasBurnerCaller(address common.Address, caller bind.ContractCaller) (*GasBurnerCaller, error) {
	contract, err := bindGasBurner(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GasBurnerCaller{contract: contract}, nil
}

// NewGasBurnerTransactor creates a new write-only instance of GasBurner, bound to a specific deployed contract.
func NewGasBurnerTransactor(address common.Address, transactor bind.ContractTransactor) (*GasBurnerTransactor, error) {
	contract, err := bindGasBurner(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GasBurnerTransactor{contract: contract}, nil
}

// NewGasBurnerFilterer creates a new log filterer instance of GasBurner, bound to a specific deployed contract.
func NewGasBurnerFilterer(address common.Address, filterer bind.ContractFilterer) (*GasBurnerFilterer, error) {
	contract, err := bindGasBurner(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GasBurnerFilterer{contract: contract}, nil
}

// bindGasBurner binds a generic wrapper to an already deployed contract.
func bindGasBurner(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := GasBurnerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GasBurner *GasBurnerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GasBurner.Contract.GasBurnerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GasBurner *GasBurnerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GasBurner.Contract.GasBurnerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GasBurner *GasBurnerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GasBurner.Contract.GasBurnerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GasBurner *GasBurnerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GasBurner.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GasBurner *GasBurnerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GasBurner.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GasBurner *GasBurnerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GasBurner.Contract.contract.Transact(opts, method, params...)
}

// Burn1000k is a paid mutator transaction binding the contract method 0xe3793770.
//
// Solidity: function burn1000k() returns()
func (_GasBurner *GasBurnerTransactor) Burn1000k(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GasBurner.contract.Transact(opts, "burn1000k")
}

// Burn1000k is a paid mutator transaction binding the contract method 0xe3793770.
//
// Solidity: function burn1000k() returns()
func (_GasBurner *GasBurnerSession) Burn1000k() (*types.Transaction, error) {
	return _GasBurner.Contract.Burn1000k(&_GasBurner.TransactOpts)
}

// Burn1000k is a paid mutator transaction binding the contract method 0xe3793770.
//
// Solidity: function burn1000k() returns()
func (_GasBurner *GasBurnerTransactorSession) Burn1000k() (*types.Transaction, error) {
	return _GasBurner.Contract.Burn1000k(&_GasBurner.TransactOpts)
}

// Burn100k is a paid mutator transaction binding the contract method 0x419f772e.
//
// Solidity: function burn100k() returns()
func (_GasBurner *GasBurnerTransactor) Burn100k(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GasBurner.contract.Transact(opts, "burn100k")
}

// Burn100k is a paid mutator transaction binding the contract method 0x419f772e.
//
// Solidity: function burn100k() returns()
func (_GasBurner *GasBurnerSession) Burn100k() (*types.Transaction, error) {
	return _GasBurner.Contract.Burn100k(&_GasBurner.TransactOpts)
}

// Burn100k is a paid mutator transaction binding the contract method 0x419f772e.
//
// Solidity: function burn100k() returns()
func (_GasBurner *GasBurnerTransactorSession) Burn100k() (*types.Transaction, error) {
	return _GasBurner.Contract.Burn100k(&_GasBurner.TransactOpts)
}

// Burn1500k is a paid mutator transaction binding the contract method 0xa2cea635.
//
// Solidity: function burn1500k() returns()
func (_GasBurner *GasBurnerTransactor) Burn1500k(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GasBurner.contract.Transact(opts, "burn1500k")
}

// Burn1500k is a paid mutator transaction binding the contract method 0xa2cea635.
//
// Solidity: function burn1500k() returns()
func (_GasBurner *GasBurnerSession) Burn1500k() (*types.Transaction, error) {
	return _GasBurner.Contract.Burn1500k(&_GasBurner.TransactOpts)
}

// Burn1500k is a paid mutator transaction binding the contract method 0xa2cea635.
//
// Solidity: function burn1500k() returns()
func (_GasBurner *GasBurnerTransactorSession) Burn1500k() (*types.Transaction, error) {
	return _GasBurner.Contract.Burn1500k(&_GasBurner.TransactOpts)
}

// Burn2000k is a paid mutator transaction binding the contract method 0xd0da1e35.
//
// Solidity: function burn2000k() returns()
func (_GasBurner *GasBurnerTransactor) Burn2000k(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GasBurner.contract.Transact(opts, "burn2000k")
}

// Burn2000k is a paid mutator transaction binding the contract method 0xd0da1e35.
//
// Solidity: function burn2000k() returns()
func (_GasBurner *GasBurnerSession) Burn2000k() (*types.Transaction, error) {
	return _GasBurner.Contract.Burn2000k(&_GasBurner.TransactOpts)
}

// Burn2000k is a paid mutator transaction binding the contract method 0xd0da1e35.
//
// Solidity: function burn2000k() returns()
func (_GasBurner *GasBurnerTransactorSession) Burn2000k() (*types.Transaction, error) {
	return _GasBurner.Contract.Burn2000k(&_GasBurner.TransactOpts)
}

// Burn500k is a paid mutator transaction binding the contract method 0x1ff94f73.
//
// Solidity: function burn500k() returns()
func (_GasBurner *GasBurnerTransactor) Burn500k(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GasBurner.contract.Transact(opts, "burn500k")
}

// Burn500k is a paid mutator transaction binding the contract method 0x1ff94f73.
//
// Solidity: function burn500k() returns()
func (_GasBurner *GasBurnerSession) Burn500k() (*types.Transaction, error) {
	return _GasBurner.Contract.Burn500k(&_GasBurner.TransactOpts)
}

// Burn500k is a paid mutator transaction binding the contract method 0x1ff94f73.
//
// Solidity: function burn500k() returns()
func (_GasBurner *GasBurnerTransactorSession) Burn500k() (*types.Transaction, error) {
	return _GasBurner.Contract.Burn500k(&_GasBurner.TransactOpts)
}

// BurnGasUnits is a paid mutator transaction binding the contract method 0x41587359.
//
// Solidity: function burnGasUnits(uint256 amount) returns()
func (_GasBurner *GasBurnerTransactor) BurnGasUnits(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _GasBurner.contract.Transact(opts, "burnGasUnits", amount)
}

// BurnGasUnits is a paid mutator transaction binding the contract method 0x41587359.
//
// Solidity: function burnGasUnits(uint256 amount) returns()
func (_GasBurner *GasBurnerSession) BurnGasUnits(amount *big.Int) (*types.Transaction, error) {
	return _GasBurner.Contract.BurnGasUnits(&_GasBurner.TransactOpts, amount)
}

// BurnGasUnits is a paid mutator transaction binding the contract method 0x41587359.
//
// Solidity: function burnGasUnits(uint256 amount) returns()
func (_GasBurner *GasBurnerTransactorSession) BurnGasUnits(amount *big.Int) (*types.Transaction, error) {
	return _GasBurner.Contract.BurnGasUnits(&_GasBurner.TransactOpts, amount)
}
