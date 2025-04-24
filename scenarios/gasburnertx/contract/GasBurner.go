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

// GasBurnerMetaData contains all meta data concerning the GasBurner contract.
var GasBurnerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"workerCode\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"burn1000k\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"burn100k\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"burn1500k\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"burn2000k\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"burn500k\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burnGasUnits\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"worker\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801562000010575f80fd5b5060405162000918380380620009188339818101604052810190620000369190620002a2565b62000047816200008c60201b60201c565b5f806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550506200036f565b5f808251602084015ff09050803b620000a3575f80fd5b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff160362000114576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016200010b906200034f565b60405180910390fd5b80915050919050565b5f604051905090565b5f80fd5b5f80fd5b5f80fd5b5f80fd5b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b6200017e8262000136565b810181811067ffffffffffffffff82111715620001a0576200019f62000146565b5b80604052505050565b5f620001b46200011d565b9050620001c2828262000173565b919050565b5f67ffffffffffffffff821115620001e457620001e362000146565b5b620001ef8262000136565b9050602081019050919050565b5f5b838110156200021b578082015181840152602081019050620001fe565b5f8484015250505050565b5f6200023c6200023684620001c7565b620001a9565b9050828152602081018484840111156200025b576200025a62000132565b5b62000268848285620001fc565b509392505050565b5f82601f8301126200028757620002866200012e565b5b81516200029984826020860162000226565b91505092915050565b5f60208284031215620002ba57620002b962000126565b5b5f82015167ffffffffffffffff811115620002da57620002d96200012a565b5b620002e88482850162000270565b91505092915050565b5f82825260208201905092915050565b7f637265617465206661696c6564000000000000000000000000000000000000005f82015250565b5f62000337600d83620002f1565b9150620003448262000301565b602082019050919050565b5f6020820190508181035f830152620003688162000329565b9050919050565b61059b806200037d5f395ff3fe608060405234801561000f575f80fd5b506004361061007b575f3560e01c80634d547ada116100595780634d547ada146100af578063a2cea635146100cd578063d0da1e35146100d7578063e3793770146100e15761007b565b80631ff94f731461007f5780634158735914610089578063419f772e146100a5575b5f80fd5b6100876100eb565b005b6100a3600480360381019061009e9190610319565b610110565b005b6100ad610133565b005b6100b7610158565b6040516100c49190610383565b60405180910390f35b6100d561017b565b005b6100df6101a0565b005b6100e96101c5565b005b6100f76207a1206101ea565b60015f815480929190610109906103c9565b9190505550565b610119816101ea565b60015f81548092919061012b906103c9565b919050555050565b61013f620186a06101ea565b60015f815480929190610151906103c9565b9190505550565b5f8054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6101876216e3606101ea565b60015f815480929190610199906103c9565b9190505550565b6101ac621e84806101ea565b60015f8154809291906101be906103c9565b9190505550565b6101d1620f42406101ea565b60015f8154809291906101e3906103c9565b9190505550565b5f805f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1661c1f2836102309190610410565b6040516020016102409190610452565b60405160208183030381529060405260405161025c91906104d7565b5f604051808303815f865af19150503d805f8114610295576040519150601f19603f3d011682016040523d82523d5f602084013e61029a565b606091505b50509050806102de576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102d590610547565b60405180910390fd5b5050565b5f80fd5b5f819050919050565b6102f8816102e6565b8114610302575f80fd5b50565b5f81359050610313816102ef565b92915050565b5f6020828403121561032e5761032d6102e2565b5b5f61033b84828501610305565b91505092915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f61036d82610344565b9050919050565b61037d81610363565b82525050565b5f6020820190506103965f830184610374565b92915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f6103d3826102e6565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036104055761040461039c565b5b600182019050919050565b5f61041a826102e6565b9150610425836102e6565b925082820390508181111561043d5761043c61039c565b5b92915050565b61044c816102e6565b82525050565b5f6020820190506104655f830184610443565b92915050565b5f81519050919050565b5f81905092915050565b5f5b8381101561049c578082015181840152602081019050610481565b5f8484015250505050565b5f6104b18261046b565b6104bb8185610475565b93506104cb81856020860161047f565b80840191505092915050565b5f6104e282846104a7565b915081905092915050565b5f82825260208201905092915050565b7f776f726b65722063616c6c206661696c656400000000000000000000000000005f82015250565b5f6105316012836104ed565b915061053c826104fd565b602082019050919050565b5f6020820190508181035f83015261055e81610525565b905091905056fea2646970667358221220e19ebcd32d93051a272fba090cab5511cbbf7c4ebf917a54ba7bc5eaad9f545c64736f6c63430008160033",
}

// GasBurnerABI is the input ABI used to generate the binding from.
// Deprecated: Use GasBurnerMetaData.ABI instead.
var GasBurnerABI = GasBurnerMetaData.ABI

// GasBurnerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use GasBurnerMetaData.Bin instead.
var GasBurnerBin = GasBurnerMetaData.Bin

// DeployGasBurner deploys a new Ethereum contract, binding an instance of GasBurner to it.
func DeployGasBurner(auth *bind.TransactOpts, backend bind.ContractBackend, workerCode []byte) (common.Address, *types.Transaction, *GasBurner, error) {
	parsed, err := GasBurnerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(GasBurnerBin), backend, workerCode)
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

// Worker is a free data retrieval call binding the contract method 0x4d547ada.
//
// Solidity: function worker() view returns(address)
func (_GasBurner *GasBurnerCaller) Worker(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GasBurner.contract.Call(opts, &out, "worker")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Worker is a free data retrieval call binding the contract method 0x4d547ada.
//
// Solidity: function worker() view returns(address)
func (_GasBurner *GasBurnerSession) Worker() (common.Address, error) {
	return _GasBurner.Contract.Worker(&_GasBurner.CallOpts)
}

// Worker is a free data retrieval call binding the contract method 0x4d547ada.
//
// Solidity: function worker() view returns(address)
func (_GasBurner *GasBurnerCallerSession) Worker() (common.Address, error) {
	return _GasBurner.Contract.Worker(&_GasBurner.CallOpts)
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
