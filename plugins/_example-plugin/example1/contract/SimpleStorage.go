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

// SimpleStorageMetaData contains all meta data concerning the SimpleStorage contract.
var SimpleStorageMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_initialValue\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"incrementer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newValue\",\"type\":\"uint256\"}],\"name\":\"ValueIncremented\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"setter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldValue\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newValue\",\"type\":\"uint256\"}],\"name\":\"ValueSet\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"getValue\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"increment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"incrementBy\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"setValue\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561000f575f80fd5b506040516102cb3803806102cb83398101604081905261002e91610088565b5f818155600180546001600160a01b03191633908117909155604080519283526020830184905290917fc0e9036d619701c94569e1462d8120ef4a6d2b15a70d27cc942fd29ae2cc0e59910160405180910390a25061009f565b5f60208284031215610098575f80fd5b5051919050565b61021f806100ac5f395ff3fe608060405234801561000f575f80fd5b5060043610610055575f3560e01c806303df179c14610059578063209652551461006e57806355241077146100835780638da5cb5b14610096578063d09de08a146100c1575b5f80fd5b61006c6100673660046101ad565b6100c9565b005b5f546040519081526020015b60405180910390f35b61006c6100913660046101ad565b610118565b6001546100a9906001600160a01b031681565b6040516001600160a01b03909116815260200161007a565b61006c61015e565b805f808282546100d991906101c4565b90915550505f5460405190815233907fa21dda55b1348fedf9117958d187753c14c74e506bd79a71c7ba6b52ebe71c419060200160405180910390a250565b5f805490829055604080518281526020810184905233917fc0e9036d619701c94569e1462d8120ef4a6d2b15a70d27cc942fd29ae2cc0e59910160405180910390a25050565b60015f8082825461016f91906101c4565b90915550505f5460405190815233907fa21dda55b1348fedf9117958d187753c14c74e506bd79a71c7ba6b52ebe71c419060200160405180910390a2565b5f602082840312156101bd575f80fd5b5035919050565b808201808211156101e357634e487b7160e01b5f52601160045260245ffd5b9291505056fea2646970667358221220a1b8b6939f49320f465f0e3b759f6eeb0b4f930ee31fa22c05327b950c64794f64736f6c63430008160033",
}

// SimpleStorageABI is the input ABI used to generate the binding from.
// Deprecated: Use SimpleStorageMetaData.ABI instead.
var SimpleStorageABI = SimpleStorageMetaData.ABI

// SimpleStorageBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SimpleStorageMetaData.Bin instead.
var SimpleStorageBin = SimpleStorageMetaData.Bin

// DeploySimpleStorage deploys a new Ethereum contract, binding an instance of SimpleStorage to it.
func DeploySimpleStorage(auth *bind.TransactOpts, backend bind.ContractBackend, _initialValue *big.Int) (common.Address, *types.Transaction, *SimpleStorage, error) {
	parsed, err := SimpleStorageMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SimpleStorageBin), backend, _initialValue)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SimpleStorage{SimpleStorageCaller: SimpleStorageCaller{contract: contract}, SimpleStorageTransactor: SimpleStorageTransactor{contract: contract}, SimpleStorageFilterer: SimpleStorageFilterer{contract: contract}}, nil
}

// SimpleStorage is an auto generated Go binding around an Ethereum contract.
type SimpleStorage struct {
	SimpleStorageCaller     // Read-only binding to the contract
	SimpleStorageTransactor // Write-only binding to the contract
	SimpleStorageFilterer   // Log filterer for contract events
}

// SimpleStorageCaller is an auto generated read-only Go binding around an Ethereum contract.
type SimpleStorageCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimpleStorageTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SimpleStorageTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimpleStorageFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SimpleStorageFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimpleStorageSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SimpleStorageSession struct {
	Contract     *SimpleStorage    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SimpleStorageCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SimpleStorageCallerSession struct {
	Contract *SimpleStorageCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// SimpleStorageTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SimpleStorageTransactorSession struct {
	Contract     *SimpleStorageTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// SimpleStorageRaw is an auto generated low-level Go binding around an Ethereum contract.
type SimpleStorageRaw struct {
	Contract *SimpleStorage // Generic contract binding to access the raw methods on
}

// SimpleStorageCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SimpleStorageCallerRaw struct {
	Contract *SimpleStorageCaller // Generic read-only contract binding to access the raw methods on
}

// SimpleStorageTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SimpleStorageTransactorRaw struct {
	Contract *SimpleStorageTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSimpleStorage creates a new instance of SimpleStorage, bound to a specific deployed contract.
func NewSimpleStorage(address common.Address, backend bind.ContractBackend) (*SimpleStorage, error) {
	contract, err := bindSimpleStorage(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SimpleStorage{SimpleStorageCaller: SimpleStorageCaller{contract: contract}, SimpleStorageTransactor: SimpleStorageTransactor{contract: contract}, SimpleStorageFilterer: SimpleStorageFilterer{contract: contract}}, nil
}

// NewSimpleStorageCaller creates a new read-only instance of SimpleStorage, bound to a specific deployed contract.
func NewSimpleStorageCaller(address common.Address, caller bind.ContractCaller) (*SimpleStorageCaller, error) {
	contract, err := bindSimpleStorage(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleStorageCaller{contract: contract}, nil
}

// NewSimpleStorageTransactor creates a new write-only instance of SimpleStorage, bound to a specific deployed contract.
func NewSimpleStorageTransactor(address common.Address, transactor bind.ContractTransactor) (*SimpleStorageTransactor, error) {
	contract, err := bindSimpleStorage(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleStorageTransactor{contract: contract}, nil
}

// NewSimpleStorageFilterer creates a new log filterer instance of SimpleStorage, bound to a specific deployed contract.
func NewSimpleStorageFilterer(address common.Address, filterer bind.ContractFilterer) (*SimpleStorageFilterer, error) {
	contract, err := bindSimpleStorage(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SimpleStorageFilterer{contract: contract}, nil
}

// bindSimpleStorage binds a generic wrapper to an already deployed contract.
func bindSimpleStorage(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SimpleStorageMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SimpleStorage *SimpleStorageRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimpleStorage.Contract.SimpleStorageCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SimpleStorage *SimpleStorageRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleStorage.Contract.SimpleStorageTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SimpleStorage *SimpleStorageRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleStorage.Contract.SimpleStorageTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SimpleStorage *SimpleStorageCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimpleStorage.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SimpleStorage *SimpleStorageTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleStorage.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SimpleStorage *SimpleStorageTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleStorage.Contract.contract.Transact(opts, method, params...)
}

// GetValue is a free data retrieval call binding the contract method 0x20965255.
//
// Solidity: function getValue() view returns(uint256)
func (_SimpleStorage *SimpleStorageCaller) GetValue(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SimpleStorage.contract.Call(opts, &out, "getValue")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetValue is a free data retrieval call binding the contract method 0x20965255.
//
// Solidity: function getValue() view returns(uint256)
func (_SimpleStorage *SimpleStorageSession) GetValue() (*big.Int, error) {
	return _SimpleStorage.Contract.GetValue(&_SimpleStorage.CallOpts)
}

// GetValue is a free data retrieval call binding the contract method 0x20965255.
//
// Solidity: function getValue() view returns(uint256)
func (_SimpleStorage *SimpleStorageCallerSession) GetValue() (*big.Int, error) {
	return _SimpleStorage.Contract.GetValue(&_SimpleStorage.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SimpleStorage *SimpleStorageCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SimpleStorage.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SimpleStorage *SimpleStorageSession) Owner() (common.Address, error) {
	return _SimpleStorage.Contract.Owner(&_SimpleStorage.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SimpleStorage *SimpleStorageCallerSession) Owner() (common.Address, error) {
	return _SimpleStorage.Contract.Owner(&_SimpleStorage.CallOpts)
}

// Increment is a paid mutator transaction binding the contract method 0xd09de08a.
//
// Solidity: function increment() returns()
func (_SimpleStorage *SimpleStorageTransactor) Increment(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleStorage.contract.Transact(opts, "increment")
}

// Increment is a paid mutator transaction binding the contract method 0xd09de08a.
//
// Solidity: function increment() returns()
func (_SimpleStorage *SimpleStorageSession) Increment() (*types.Transaction, error) {
	return _SimpleStorage.Contract.Increment(&_SimpleStorage.TransactOpts)
}

// Increment is a paid mutator transaction binding the contract method 0xd09de08a.
//
// Solidity: function increment() returns()
func (_SimpleStorage *SimpleStorageTransactorSession) Increment() (*types.Transaction, error) {
	return _SimpleStorage.Contract.Increment(&_SimpleStorage.TransactOpts)
}

// IncrementBy is a paid mutator transaction binding the contract method 0x03df179c.
//
// Solidity: function incrementBy(uint256 _amount) returns()
func (_SimpleStorage *SimpleStorageTransactor) IncrementBy(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _SimpleStorage.contract.Transact(opts, "incrementBy", _amount)
}

// IncrementBy is a paid mutator transaction binding the contract method 0x03df179c.
//
// Solidity: function incrementBy(uint256 _amount) returns()
func (_SimpleStorage *SimpleStorageSession) IncrementBy(_amount *big.Int) (*types.Transaction, error) {
	return _SimpleStorage.Contract.IncrementBy(&_SimpleStorage.TransactOpts, _amount)
}

// IncrementBy is a paid mutator transaction binding the contract method 0x03df179c.
//
// Solidity: function incrementBy(uint256 _amount) returns()
func (_SimpleStorage *SimpleStorageTransactorSession) IncrementBy(_amount *big.Int) (*types.Transaction, error) {
	return _SimpleStorage.Contract.IncrementBy(&_SimpleStorage.TransactOpts, _amount)
}

// SetValue is a paid mutator transaction binding the contract method 0x55241077.
//
// Solidity: function setValue(uint256 _value) returns()
func (_SimpleStorage *SimpleStorageTransactor) SetValue(opts *bind.TransactOpts, _value *big.Int) (*types.Transaction, error) {
	return _SimpleStorage.contract.Transact(opts, "setValue", _value)
}

// SetValue is a paid mutator transaction binding the contract method 0x55241077.
//
// Solidity: function setValue(uint256 _value) returns()
func (_SimpleStorage *SimpleStorageSession) SetValue(_value *big.Int) (*types.Transaction, error) {
	return _SimpleStorage.Contract.SetValue(&_SimpleStorage.TransactOpts, _value)
}

// SetValue is a paid mutator transaction binding the contract method 0x55241077.
//
// Solidity: function setValue(uint256 _value) returns()
func (_SimpleStorage *SimpleStorageTransactorSession) SetValue(_value *big.Int) (*types.Transaction, error) {
	return _SimpleStorage.Contract.SetValue(&_SimpleStorage.TransactOpts, _value)
}

// SimpleStorageValueIncrementedIterator is returned from FilterValueIncremented and is used to iterate over the raw logs and unpacked data for ValueIncremented events raised by the SimpleStorage contract.
type SimpleStorageValueIncrementedIterator struct {
	Event *SimpleStorageValueIncremented // Event containing the contract specifics and raw log

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
func (it *SimpleStorageValueIncrementedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleStorageValueIncremented)
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
		it.Event = new(SimpleStorageValueIncremented)
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
func (it *SimpleStorageValueIncrementedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleStorageValueIncrementedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleStorageValueIncremented represents a ValueIncremented event raised by the SimpleStorage contract.
type SimpleStorageValueIncremented struct {
	Incrementer common.Address
	NewValue    *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterValueIncremented is a free log retrieval operation binding the contract event 0xa21dda55b1348fedf9117958d187753c14c74e506bd79a71c7ba6b52ebe71c41.
//
// Solidity: event ValueIncremented(address indexed incrementer, uint256 newValue)
func (_SimpleStorage *SimpleStorageFilterer) FilterValueIncremented(opts *bind.FilterOpts, incrementer []common.Address) (*SimpleStorageValueIncrementedIterator, error) {

	var incrementerRule []interface{}
	for _, incrementerItem := range incrementer {
		incrementerRule = append(incrementerRule, incrementerItem)
	}

	logs, sub, err := _SimpleStorage.contract.FilterLogs(opts, "ValueIncremented", incrementerRule)
	if err != nil {
		return nil, err
	}
	return &SimpleStorageValueIncrementedIterator{contract: _SimpleStorage.contract, event: "ValueIncremented", logs: logs, sub: sub}, nil
}

// WatchValueIncremented is a free log subscription operation binding the contract event 0xa21dda55b1348fedf9117958d187753c14c74e506bd79a71c7ba6b52ebe71c41.
//
// Solidity: event ValueIncremented(address indexed incrementer, uint256 newValue)
func (_SimpleStorage *SimpleStorageFilterer) WatchValueIncremented(opts *bind.WatchOpts, sink chan<- *SimpleStorageValueIncremented, incrementer []common.Address) (event.Subscription, error) {

	var incrementerRule []interface{}
	for _, incrementerItem := range incrementer {
		incrementerRule = append(incrementerRule, incrementerItem)
	}

	logs, sub, err := _SimpleStorage.contract.WatchLogs(opts, "ValueIncremented", incrementerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SimpleStorageValueIncremented)
				if err := _SimpleStorage.contract.UnpackLog(event, "ValueIncremented", log); err != nil {
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

// ParseValueIncremented is a log parse operation binding the contract event 0xa21dda55b1348fedf9117958d187753c14c74e506bd79a71c7ba6b52ebe71c41.
//
// Solidity: event ValueIncremented(address indexed incrementer, uint256 newValue)
func (_SimpleStorage *SimpleStorageFilterer) ParseValueIncremented(log types.Log) (*SimpleStorageValueIncremented, error) {
	event := new(SimpleStorageValueIncremented)
	if err := _SimpleStorage.contract.UnpackLog(event, "ValueIncremented", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SimpleStorageValueSetIterator is returned from FilterValueSet and is used to iterate over the raw logs and unpacked data for ValueSet events raised by the SimpleStorage contract.
type SimpleStorageValueSetIterator struct {
	Event *SimpleStorageValueSet // Event containing the contract specifics and raw log

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
func (it *SimpleStorageValueSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleStorageValueSet)
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
		it.Event = new(SimpleStorageValueSet)
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
func (it *SimpleStorageValueSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleStorageValueSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleStorageValueSet represents a ValueSet event raised by the SimpleStorage contract.
type SimpleStorageValueSet struct {
	Setter   common.Address
	OldValue *big.Int
	NewValue *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterValueSet is a free log retrieval operation binding the contract event 0xc0e9036d619701c94569e1462d8120ef4a6d2b15a70d27cc942fd29ae2cc0e59.
//
// Solidity: event ValueSet(address indexed setter, uint256 oldValue, uint256 newValue)
func (_SimpleStorage *SimpleStorageFilterer) FilterValueSet(opts *bind.FilterOpts, setter []common.Address) (*SimpleStorageValueSetIterator, error) {

	var setterRule []interface{}
	for _, setterItem := range setter {
		setterRule = append(setterRule, setterItem)
	}

	logs, sub, err := _SimpleStorage.contract.FilterLogs(opts, "ValueSet", setterRule)
	if err != nil {
		return nil, err
	}
	return &SimpleStorageValueSetIterator{contract: _SimpleStorage.contract, event: "ValueSet", logs: logs, sub: sub}, nil
}

// WatchValueSet is a free log subscription operation binding the contract event 0xc0e9036d619701c94569e1462d8120ef4a6d2b15a70d27cc942fd29ae2cc0e59.
//
// Solidity: event ValueSet(address indexed setter, uint256 oldValue, uint256 newValue)
func (_SimpleStorage *SimpleStorageFilterer) WatchValueSet(opts *bind.WatchOpts, sink chan<- *SimpleStorageValueSet, setter []common.Address) (event.Subscription, error) {

	var setterRule []interface{}
	for _, setterItem := range setter {
		setterRule = append(setterRule, setterItem)
	}

	logs, sub, err := _SimpleStorage.contract.WatchLogs(opts, "ValueSet", setterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SimpleStorageValueSet)
				if err := _SimpleStorage.contract.UnpackLog(event, "ValueSet", log); err != nil {
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

// ParseValueSet is a log parse operation binding the contract event 0xc0e9036d619701c94569e1462d8120ef4a6d2b15a70d27cc942fd29ae2cc0e59.
//
// Solidity: event ValueSet(address indexed setter, uint256 oldValue, uint256 newValue)
func (_SimpleStorage *SimpleStorageFilterer) ParseValueSet(log types.Log) (*SimpleStorageValueSet, error) {
	event := new(SimpleStorageValueSet)
	if err := _SimpleStorage.contract.UnpackLog(event, "ValueSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
