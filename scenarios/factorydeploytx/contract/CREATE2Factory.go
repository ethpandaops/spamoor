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

// CREATE2FactoryMetaData contains all meta data concerning the CREATE2Factory contract.
var CREATE2FactoryMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"deployedAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"}],\"name\":\"ContractDeployed\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"}],\"name\":\"deploy\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"initCodeHash\",\"type\":\"bytes32\"}],\"name\":\"predictAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610365806100206000396000f3fe6080604052600436106100295760003560e01c806310a935281461002e578063cdcb760a14610064575b600080fd5b34801561003a57600080fd5b5061004e6100493660046101db565b610077565b60405161005b91906102d7565b60405180910390f35b61004e6100723660046101fc565b6100ee565b6040516000906100b1907fff0000000000000000000000000000000000000000000000000000000000000090309086908690602001610273565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905280516020909101209392505050565b600080600084848080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050825192935088929150506020830134f5915073ffffffffffffffffffffffffffffffffffffffff821661018f576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610186906102f8565b60405180910390fd5b604051869073ffffffffffffffffffffffffffffffffffffffff8416907fb085ff794f342ed78acc7791d067e28a931e614b52476c0305795e1ff0a154bc90600090a350949350505050565b600080604083850312156101ed578182fd5b50508035926020909101359150565b600080600060408486031215610210578081fd5b83359250602084013567ffffffffffffffff8082111561022e578283fd5b818601915086601f830112610241578283fd5b81358181111561024f578384fd5b876020828501011115610260578384fd5b6020830194508093505050509250925092565b7fff0000000000000000000000000000000000000000000000000000000000000094909416845260609290921b7fffffffffffffffffffffffffffffffffffffffff0000000000000000000000001660018401526015830152603582015260550190565b73ffffffffffffffffffffffffffffffffffffffff91909116815260200190565b60208082526011908201527f4465706c6f796d656e74206661696c656400000000000000000000000000000060408201526060019056fea26469706673582212202d3e87dd998c22df28ccb2c934734610461c1e6888114d8003aa51583d65054c64736f6c63430008000033",
}

// CREATE2FactoryABI is the input ABI used to generate the binding from.
// Deprecated: Use CREATE2FactoryMetaData.ABI instead.
var CREATE2FactoryABI = CREATE2FactoryMetaData.ABI

// CREATE2FactoryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use CREATE2FactoryMetaData.Bin instead.
var CREATE2FactoryBin = CREATE2FactoryMetaData.Bin

// DeployCREATE2Factory deploys a new Ethereum contract, binding an instance of CREATE2Factory to it.
func DeployCREATE2Factory(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *CREATE2Factory, error) {
	parsed, err := CREATE2FactoryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CREATE2FactoryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CREATE2Factory{CREATE2FactoryCaller: CREATE2FactoryCaller{contract: contract}, CREATE2FactoryTransactor: CREATE2FactoryTransactor{contract: contract}, CREATE2FactoryFilterer: CREATE2FactoryFilterer{contract: contract}}, nil
}

// CREATE2Factory is an auto generated Go binding around an Ethereum contract.
type CREATE2Factory struct {
	CREATE2FactoryCaller     // Read-only binding to the contract
	CREATE2FactoryTransactor // Write-only binding to the contract
	CREATE2FactoryFilterer   // Log filterer for contract events
}

// CREATE2FactoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type CREATE2FactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CREATE2FactoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CREATE2FactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CREATE2FactoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CREATE2FactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CREATE2FactorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CREATE2FactorySession struct {
	Contract     *CREATE2Factory   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CREATE2FactoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CREATE2FactoryCallerSession struct {
	Contract *CREATE2FactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// CREATE2FactoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CREATE2FactoryTransactorSession struct {
	Contract     *CREATE2FactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// CREATE2FactoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type CREATE2FactoryRaw struct {
	Contract *CREATE2Factory // Generic contract binding to access the raw methods on
}

// CREATE2FactoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CREATE2FactoryCallerRaw struct {
	Contract *CREATE2FactoryCaller // Generic read-only contract binding to access the raw methods on
}

// CREATE2FactoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CREATE2FactoryTransactorRaw struct {
	Contract *CREATE2FactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCREATE2Factory creates a new instance of CREATE2Factory, bound to a specific deployed contract.
func NewCREATE2Factory(address common.Address, backend bind.ContractBackend) (*CREATE2Factory, error) {
	contract, err := bindCREATE2Factory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CREATE2Factory{CREATE2FactoryCaller: CREATE2FactoryCaller{contract: contract}, CREATE2FactoryTransactor: CREATE2FactoryTransactor{contract: contract}, CREATE2FactoryFilterer: CREATE2FactoryFilterer{contract: contract}}, nil
}

// NewCREATE2FactoryCaller creates a new read-only instance of CREATE2Factory, bound to a specific deployed contract.
func NewCREATE2FactoryCaller(address common.Address, caller bind.ContractCaller) (*CREATE2FactoryCaller, error) {
	contract, err := bindCREATE2Factory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CREATE2FactoryCaller{contract: contract}, nil
}

// NewCREATE2FactoryTransactor creates a new write-only instance of CREATE2Factory, bound to a specific deployed contract.
func NewCREATE2FactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*CREATE2FactoryTransactor, error) {
	contract, err := bindCREATE2Factory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CREATE2FactoryTransactor{contract: contract}, nil
}

// NewCREATE2FactoryFilterer creates a new log filterer instance of CREATE2Factory, bound to a specific deployed contract.
func NewCREATE2FactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*CREATE2FactoryFilterer, error) {
	contract, err := bindCREATE2Factory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CREATE2FactoryFilterer{contract: contract}, nil
}

// bindCREATE2Factory binds a generic wrapper to an already deployed contract.
func bindCREATE2Factory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CREATE2FactoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CREATE2Factory *CREATE2FactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CREATE2Factory.Contract.CREATE2FactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CREATE2Factory *CREATE2FactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CREATE2Factory.Contract.CREATE2FactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CREATE2Factory *CREATE2FactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CREATE2Factory.Contract.CREATE2FactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CREATE2Factory *CREATE2FactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CREATE2Factory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CREATE2Factory *CREATE2FactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CREATE2Factory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CREATE2Factory *CREATE2FactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CREATE2Factory.Contract.contract.Transact(opts, method, params...)
}

// PredictAddress is a free data retrieval call binding the contract method 0x10a93528.
//
// Solidity: function predictAddress(bytes32 salt, bytes32 initCodeHash) view returns(address)
func (_CREATE2Factory *CREATE2FactoryCaller) PredictAddress(opts *bind.CallOpts, salt [32]byte, initCodeHash [32]byte) (common.Address, error) {
	var out []interface{}
	err := _CREATE2Factory.contract.Call(opts, &out, "predictAddress", salt, initCodeHash)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PredictAddress is a free data retrieval call binding the contract method 0x10a93528.
//
// Solidity: function predictAddress(bytes32 salt, bytes32 initCodeHash) view returns(address)
func (_CREATE2Factory *CREATE2FactorySession) PredictAddress(salt [32]byte, initCodeHash [32]byte) (common.Address, error) {
	return _CREATE2Factory.Contract.PredictAddress(&_CREATE2Factory.CallOpts, salt, initCodeHash)
}

// PredictAddress is a free data retrieval call binding the contract method 0x10a93528.
//
// Solidity: function predictAddress(bytes32 salt, bytes32 initCodeHash) view returns(address)
func (_CREATE2Factory *CREATE2FactoryCallerSession) PredictAddress(salt [32]byte, initCodeHash [32]byte) (common.Address, error) {
	return _CREATE2Factory.Contract.PredictAddress(&_CREATE2Factory.CallOpts, salt, initCodeHash)
}

// Deploy is a paid mutator transaction binding the contract method 0xcdcb760a.
//
// Solidity: function deploy(bytes32 salt, bytes initCode) payable returns(address)
func (_CREATE2Factory *CREATE2FactoryTransactor) Deploy(opts *bind.TransactOpts, salt [32]byte, initCode []byte) (*types.Transaction, error) {
	return _CREATE2Factory.contract.Transact(opts, "deploy", salt, initCode)
}

// Deploy is a paid mutator transaction binding the contract method 0xcdcb760a.
//
// Solidity: function deploy(bytes32 salt, bytes initCode) payable returns(address)
func (_CREATE2Factory *CREATE2FactorySession) Deploy(salt [32]byte, initCode []byte) (*types.Transaction, error) {
	return _CREATE2Factory.Contract.Deploy(&_CREATE2Factory.TransactOpts, salt, initCode)
}

// Deploy is a paid mutator transaction binding the contract method 0xcdcb760a.
//
// Solidity: function deploy(bytes32 salt, bytes initCode) payable returns(address)
func (_CREATE2Factory *CREATE2FactoryTransactorSession) Deploy(salt [32]byte, initCode []byte) (*types.Transaction, error) {
	return _CREATE2Factory.Contract.Deploy(&_CREATE2Factory.TransactOpts, salt, initCode)
}

// CREATE2FactoryContractDeployedIterator is returned from FilterContractDeployed and is used to iterate over the raw logs and unpacked data for ContractDeployed events raised by the CREATE2Factory contract.
type CREATE2FactoryContractDeployedIterator struct {
	Event *CREATE2FactoryContractDeployed // Event containing the contract specifics and raw log

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
func (it *CREATE2FactoryContractDeployedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CREATE2FactoryContractDeployed)
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
		it.Event = new(CREATE2FactoryContractDeployed)
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
func (it *CREATE2FactoryContractDeployedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CREATE2FactoryContractDeployedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CREATE2FactoryContractDeployed represents a ContractDeployed event raised by the CREATE2Factory contract.
type CREATE2FactoryContractDeployed struct {
	DeployedAddress common.Address
	Salt            [32]byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterContractDeployed is a free log retrieval operation binding the contract event 0xb085ff794f342ed78acc7791d067e28a931e614b52476c0305795e1ff0a154bc.
//
// Solidity: event ContractDeployed(address indexed deployedAddress, bytes32 indexed salt)
func (_CREATE2Factory *CREATE2FactoryFilterer) FilterContractDeployed(opts *bind.FilterOpts, deployedAddress []common.Address, salt [][32]byte) (*CREATE2FactoryContractDeployedIterator, error) {

	var deployedAddressRule []interface{}
	for _, deployedAddressItem := range deployedAddress {
		deployedAddressRule = append(deployedAddressRule, deployedAddressItem)
	}
	var saltRule []interface{}
	for _, saltItem := range salt {
		saltRule = append(saltRule, saltItem)
	}

	logs, sub, err := _CREATE2Factory.contract.FilterLogs(opts, "ContractDeployed", deployedAddressRule, saltRule)
	if err != nil {
		return nil, err
	}
	return &CREATE2FactoryContractDeployedIterator{contract: _CREATE2Factory.contract, event: "ContractDeployed", logs: logs, sub: sub}, nil
}

// WatchContractDeployed is a free log subscription operation binding the contract event 0xb085ff794f342ed78acc7791d067e28a931e614b52476c0305795e1ff0a154bc.
//
// Solidity: event ContractDeployed(address indexed deployedAddress, bytes32 indexed salt)
func (_CREATE2Factory *CREATE2FactoryFilterer) WatchContractDeployed(opts *bind.WatchOpts, sink chan<- *CREATE2FactoryContractDeployed, deployedAddress []common.Address, salt [][32]byte) (event.Subscription, error) {

	var deployedAddressRule []interface{}
	for _, deployedAddressItem := range deployedAddress {
		deployedAddressRule = append(deployedAddressRule, deployedAddressItem)
	}
	var saltRule []interface{}
	for _, saltItem := range salt {
		saltRule = append(saltRule, saltItem)
	}

	logs, sub, err := _CREATE2Factory.contract.WatchLogs(opts, "ContractDeployed", deployedAddressRule, saltRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CREATE2FactoryContractDeployed)
				if err := _CREATE2Factory.contract.UnpackLog(event, "ContractDeployed", log); err != nil {
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

// ParseContractDeployed is a log parse operation binding the contract event 0xb085ff794f342ed78acc7791d067e28a931e614b52476c0305795e1ff0a154bc.
//
// Solidity: event ContractDeployed(address indexed deployedAddress, bytes32 indexed salt)
func (_CREATE2Factory *CREATE2FactoryFilterer) ParseContractDeployed(log types.Log) (*CREATE2FactoryContractDeployed, error) {
	event := new(CREATE2FactoryContractDeployed)
	if err := _CREATE2Factory.contract.UnpackLog(event, "ContractDeployed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
