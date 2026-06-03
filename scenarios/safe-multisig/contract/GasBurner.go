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
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rounds\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"digest\",\"type\":\"bytes32\"}],\"name\":\"GasBurned\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxRounds\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"rounds\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561000f575f80fd5b506101c98061001d5f395ff3fe608060405234801561000f575f80fd5b5060043610610029575f3560e01c8063b390c0ab1461002d575b5f80fd5b61004061003b36600461012f565b610052565b60405190815260200160405180910390f35b5f815f0361005f57600191505b5f8360405160200161007391815260200190565b60408051601f1981840301815291905280516020909101209050610097838261014f565b6100a290600161016e565b91505f5b828110156100e457604080516020810184905290810182905260600160408051601f19818403018152919052805160209091012091506001016100a6565b50604080518581526020810184905290810182905233907f55eb00fd9d99665c565ffd1235b5b5df1de64efc57d0ed57d197fd7256e8d0549060600160405180910390a25092915050565b5f8060408385031215610140575f80fd5b50508035926020909101359150565b5f8261016957634e487b7160e01b5f52601260045260245ffd5b500690565b8082018082111561018d57634e487b7160e01b5f52601160045260245ffd5b9291505056fea2646970667358221220c80f7ef215f80149137a9c5503a2783ef8487faa8d6b6ea64f57e6e4fc26fe6764736f6c63430008180033",
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

// Burn is a paid mutator transaction binding the contract method 0xb390c0ab.
//
// Solidity: function burn(uint256 seed, uint256 maxRounds) returns(uint256 rounds)
func (_GasBurner *GasBurnerTransactor) Burn(opts *bind.TransactOpts, seed *big.Int, maxRounds *big.Int) (*types.Transaction, error) {
	return _GasBurner.contract.Transact(opts, "burn", seed, maxRounds)
}

// Burn is a paid mutator transaction binding the contract method 0xb390c0ab.
//
// Solidity: function burn(uint256 seed, uint256 maxRounds) returns(uint256 rounds)
func (_GasBurner *GasBurnerSession) Burn(seed *big.Int, maxRounds *big.Int) (*types.Transaction, error) {
	return _GasBurner.Contract.Burn(&_GasBurner.TransactOpts, seed, maxRounds)
}

// Burn is a paid mutator transaction binding the contract method 0xb390c0ab.
//
// Solidity: function burn(uint256 seed, uint256 maxRounds) returns(uint256 rounds)
func (_GasBurner *GasBurnerTransactorSession) Burn(seed *big.Int, maxRounds *big.Int) (*types.Transaction, error) {
	return _GasBurner.Contract.Burn(&_GasBurner.TransactOpts, seed, maxRounds)
}

// GasBurnerGasBurnedIterator is returned from FilterGasBurned and is used to iterate over the raw logs and unpacked data for GasBurned events raised by the GasBurner contract.
type GasBurnerGasBurnedIterator struct {
	Event *GasBurnerGasBurned // Event containing the contract specifics and raw log

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
func (it *GasBurnerGasBurnedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GasBurnerGasBurned)
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
		it.Event = new(GasBurnerGasBurned)
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
func (it *GasBurnerGasBurnedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GasBurnerGasBurnedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GasBurnerGasBurned represents a GasBurned event raised by the GasBurner contract.
type GasBurnerGasBurned struct {
	Caller common.Address
	Seed   *big.Int
	Rounds *big.Int
	Digest [32]byte
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterGasBurned is a free log retrieval operation binding the contract event 0x55eb00fd9d99665c565ffd1235b5b5df1de64efc57d0ed57d197fd7256e8d054.
//
// Solidity: event GasBurned(address indexed caller, uint256 seed, uint256 rounds, bytes32 digest)
func (_GasBurner *GasBurnerFilterer) FilterGasBurned(opts *bind.FilterOpts, caller []common.Address) (*GasBurnerGasBurnedIterator, error) {

	var callerRule []interface{}
	for _, callerItem := range caller {
		callerRule = append(callerRule, callerItem)
	}

	logs, sub, err := _GasBurner.contract.FilterLogs(opts, "GasBurned", callerRule)
	if err != nil {
		return nil, err
	}
	return &GasBurnerGasBurnedIterator{contract: _GasBurner.contract, event: "GasBurned", logs: logs, sub: sub}, nil
}

// WatchGasBurned is a free log subscription operation binding the contract event 0x55eb00fd9d99665c565ffd1235b5b5df1de64efc57d0ed57d197fd7256e8d054.
//
// Solidity: event GasBurned(address indexed caller, uint256 seed, uint256 rounds, bytes32 digest)
func (_GasBurner *GasBurnerFilterer) WatchGasBurned(opts *bind.WatchOpts, sink chan<- *GasBurnerGasBurned, caller []common.Address) (event.Subscription, error) {

	var callerRule []interface{}
	for _, callerItem := range caller {
		callerRule = append(callerRule, callerItem)
	}

	logs, sub, err := _GasBurner.contract.WatchLogs(opts, "GasBurned", callerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GasBurnerGasBurned)
				if err := _GasBurner.contract.UnpackLog(event, "GasBurned", log); err != nil {
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

// ParseGasBurned is a log parse operation binding the contract event 0x55eb00fd9d99665c565ffd1235b5b5df1de64efc57d0ed57d197fd7256e8d054.
//
// Solidity: event GasBurned(address indexed caller, uint256 seed, uint256 rounds, bytes32 digest)
func (_GasBurner *GasBurnerFilterer) ParseGasBurned(log types.Log) (*GasBurnerGasBurned, error) {
	event := new(GasBurnerGasBurned)
	if err := _GasBurner.contract.UnpackLog(event, "GasBurned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
