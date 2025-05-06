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
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"workerCode\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"gas\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"loops\",\"type\":\"uint256\"}],\"name\":\"BurnedGas\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"burn1000k\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"burn100k\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"burn1500k\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"burn2000k\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"burn500k\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burnGasUnits\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"worker\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801562000010575f80fd5b50604051620008d4380380620008d48339818101604052810190620000369190620002a2565b62000047816200008c60201b60201c565b5f806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550506200036f565b5f808251602084015ff09050803b620000a3575f80fd5b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff160362000114576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016200010b906200034f565b60405180910390fd5b80915050919050565b5f604051905090565b5f80fd5b5f80fd5b5f80fd5b5f80fd5b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b6200017e8262000136565b810181811067ffffffffffffffff82111715620001a0576200019f62000146565b5b80604052505050565b5f620001b46200011d565b9050620001c2828262000173565b919050565b5f67ffffffffffffffff821115620001e457620001e362000146565b5b620001ef8262000136565b9050602081019050919050565b5f5b838110156200021b578082015181840152602081019050620001fe565b5f8484015250505050565b5f6200023c6200023684620001c7565b620001a9565b9050828152602081018484840111156200025b576200025a62000132565b5b62000268848285620001fc565b509392505050565b5f82601f8301126200028757620002866200012e565b5b81516200029984826020860162000226565b91505092915050565b5f60208284031215620002ba57620002b962000126565b5b5f82015167ffffffffffffffff811115620002da57620002d96200012a565b5b620002e88482850162000270565b91505092915050565b5f82825260208201905092915050565b7f637265617465206661696c6564000000000000000000000000000000000000005f82015250565b5f62000337600d83620002f1565b9150620003448262000301565b602082019050919050565b5f6020820190508181035f830152620003688162000329565b9050919050565b610557806200037d5f395ff3fe608060405234801561000f575f80fd5b506004361061007b575f3560e01c80634d547ada116100595780634d547ada146100af578063a2cea635146100cd578063d0da1e35146100d7578063e3793770146100e15761007b565b80631ff94f731461007f5780634158735914610089578063419f772e146100a5575b5f80fd5b6100876100eb565b005b6100a3600480360381019061009e91906102b5565b6100f9565b005b6100ad610105565b005b6100b7610113565b6040516100c4919061031f565b60405180910390f35b6100d5610136565b005b6100df610144565b005b6100e9610152565b005b6100f76207a120610160565b565b61010281610160565b50565b610111620186a0610160565b565b5f8054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6101426216e360610160565b565b610150621e8480610160565b565b61015e620f4240610160565b565b5f805f8054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16614970846101a69190610365565b6040516101b2906103c5565b5f604051808303815f8787f1925050503d805f81146101ec576040519150601f19603f3d011682016040523d82523d5f602084013e6101f1565b606091505b509150915081610236576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161022d90610433565b60405180910390fd5b827fe26f5b262c5e61409339d774c1d1b647a193255d9172f54300225389190e05ad8261026290610493565b5f1c6040516102719190610508565b60405180910390a2505050565b5f80fd5b5f819050919050565b61029481610282565b811461029e575f80fd5b50565b5f813590506102af8161028b565b92915050565b5f602082840312156102ca576102c961027e565b5b5f6102d7848285016102a1565b91505092915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f610309826102e0565b9050919050565b610319816102ff565b82525050565b5f6020820190506103325f830184610310565b92915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f61036f82610282565b915061037a83610282565b925082820390508181111561039257610391610338565b5b92915050565b5f81905092915050565b50565b5f6103b05f83610398565b91506103bb826103a2565b5f82019050919050565b5f6103cf826103a5565b9150819050919050565b5f82825260208201905092915050565b7f776f726b65722063616c6c206661696c656400000000000000000000000000005f82015250565b5f61041d6012836103d9565b9150610428826103e9565b602082019050919050565b5f6020820190508181035f83015261044a81610411565b9050919050565b5f81519050919050565b5f819050602082019050919050565b5f819050919050565b5f61047e825161046a565b80915050919050565b5f82821b905092915050565b5f61049d82610451565b826104a78461045b565b90506104b281610473565b925060208210156104f2576104ed7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff83602003600802610487565b831692505b5050919050565b61050281610282565b82525050565b5f60208201905061051b5f8301846104f9565b9291505056fea264697066735822122058d6bfaf47ccb485962334ed95b0313cbf88eb0b98ff65e9d247da605a580e7c64736f6c63430008160033",
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

// GasBurnerBurnedGasIterator is returned from FilterBurnedGas and is used to iterate over the raw logs and unpacked data for BurnedGas events raised by the GasBurner contract.
type GasBurnerBurnedGasIterator struct {
	Event *GasBurnerBurnedGas // Event containing the contract specifics and raw log

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
func (it *GasBurnerBurnedGasIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GasBurnerBurnedGas)
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
		it.Event = new(GasBurnerBurnedGas)
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
func (it *GasBurnerBurnedGasIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GasBurnerBurnedGasIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GasBurnerBurnedGas represents a BurnedGas event raised by the GasBurner contract.
type GasBurnerBurnedGas struct {
	Gas   *big.Int
	Loops *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterBurnedGas is a free log retrieval operation binding the contract event 0xe26f5b262c5e61409339d774c1d1b647a193255d9172f54300225389190e05ad.
//
// Solidity: event BurnedGas(uint256 indexed gas, uint256 loops)
func (_GasBurner *GasBurnerFilterer) FilterBurnedGas(opts *bind.FilterOpts, gas []*big.Int) (*GasBurnerBurnedGasIterator, error) {

	var gasRule []interface{}
	for _, gasItem := range gas {
		gasRule = append(gasRule, gasItem)
	}

	logs, sub, err := _GasBurner.contract.FilterLogs(opts, "BurnedGas", gasRule)
	if err != nil {
		return nil, err
	}
	return &GasBurnerBurnedGasIterator{contract: _GasBurner.contract, event: "BurnedGas", logs: logs, sub: sub}, nil
}

// WatchBurnedGas is a free log subscription operation binding the contract event 0xe26f5b262c5e61409339d774c1d1b647a193255d9172f54300225389190e05ad.
//
// Solidity: event BurnedGas(uint256 indexed gas, uint256 loops)
func (_GasBurner *GasBurnerFilterer) WatchBurnedGas(opts *bind.WatchOpts, sink chan<- *GasBurnerBurnedGas, gas []*big.Int) (event.Subscription, error) {

	var gasRule []interface{}
	for _, gasItem := range gas {
		gasRule = append(gasRule, gasItem)
	}

	logs, sub, err := _GasBurner.contract.WatchLogs(opts, "BurnedGas", gasRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GasBurnerBurnedGas)
				if err := _GasBurner.contract.UnpackLog(event, "BurnedGas", log); err != nil {
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

// ParseBurnedGas is a log parse operation binding the contract event 0xe26f5b262c5e61409339d774c1d1b647a193255d9172f54300225389190e05ad.
//
// Solidity: event BurnedGas(uint256 indexed gas, uint256 loops)
func (_GasBurner *GasBurnerFilterer) ParseBurnedGas(log types.Log) (*GasBurnerBurnedGas, error) {
	event := new(GasBurnerBurnedGas)
	if err := _GasBurner.contract.UnpackLog(event, "BurnedGas", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
