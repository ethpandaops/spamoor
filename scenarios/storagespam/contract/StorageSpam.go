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

// StorageSpamMetaData contains all meta data concerning the StorageSpam contract.
var StorageSpamMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"gas\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"loops\",\"type\":\"uint256\"}],\"name\":\"RandomForGas\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"key\",\"type\":\"uint256\"}],\"name\":\"getStorage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"txid\",\"type\":\"uint256\"}],\"name\":\"setRandomForGas\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"key\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setStorage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"storageMap\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561000f575f80fd5b5061049f8061001d5f395ff3fe608060405234801561000f575f80fd5b506004361061004a575f3560e01c80635e666e4a1461004e57806365fd47721461007e578063936ad72f146100ae578063fed72935146100ca575b5f80fd5b61006860048036038101906100639190610289565b6100e6565b60405161007591906102c3565b60405180910390f35b61009860048036038101906100939190610289565b6100fa565b6040516100a591906102c3565b60405180910390f35b6100c860048036038101906100c391906102dc565b610113565b005b6100e460048036038101906100df91906102dc565b61012c565b005b5f602052805f5260405f205f915090505481565b5f805f8381526020019081526020015f20549050919050565b805f808481526020019081526020015f20819055505050565b5f5a90505f8311801561013e57508281115b1561015657828161014f9190610347565b905061015b565b606490505b5f80438460405160200161017092919061039a565b604051602081830303815290604052805190602001205f1c90505b825a1115610213575f60c8836101a191906103f2565b036101d45780826040516020016101b992919061039a565b604051602081830303815290604052805190602001205f1c90505b608081901c608082901b175f808381526020019081526020015f208190555060ff81901c600182901b179050818061020b90610422565b92505061018b565b847f6951887835703bee48a4b856b283f6bb242d53dcd20c2a257a95e9326e594cd18360405161024391906102c3565b60405180910390a25050505050565b5f80fd5b5f819050919050565b61026881610256565b8114610272575f80fd5b50565b5f813590506102838161025f565b92915050565b5f6020828403121561029e5761029d610252565b5b5f6102ab84828501610275565b91505092915050565b6102bd81610256565b82525050565b5f6020820190506102d65f8301846102b4565b92915050565b5f80604083850312156102f2576102f1610252565b5b5f6102ff85828601610275565b925050602061031085828601610275565b9150509250929050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f61035182610256565b915061035c83610256565b92508282039050818111156103745761037361031a565b5b92915050565b5f819050919050565b61039461038f82610256565b61037a565b82525050565b5f6103a58285610383565b6020820191506103b58284610383565b6020820191508190509392505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601260045260245ffd5b5f6103fc82610256565b915061040783610256565b925082610417576104166103c5565b5b828206905092915050565b5f61042c82610256565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361045e5761045d61031a565b5b60018201905091905056fea2646970667358221220056c9be65d1ccef6406a26a83e27e531c21ddf89b8f988013315c799e3f98b0164736f6c63430008160033",
}

// StorageSpamABI is the input ABI used to generate the binding from.
// Deprecated: Use StorageSpamMetaData.ABI instead.
var StorageSpamABI = StorageSpamMetaData.ABI

// StorageSpamBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use StorageSpamMetaData.Bin instead.
var StorageSpamBin = StorageSpamMetaData.Bin

// DeployStorageSpam deploys a new Ethereum contract, binding an instance of StorageSpam to it.
func DeployStorageSpam(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *StorageSpam, error) {
	parsed, err := StorageSpamMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StorageSpamBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &StorageSpam{StorageSpamCaller: StorageSpamCaller{contract: contract}, StorageSpamTransactor: StorageSpamTransactor{contract: contract}, StorageSpamFilterer: StorageSpamFilterer{contract: contract}}, nil
}

// StorageSpam is an auto generated Go binding around an Ethereum contract.
type StorageSpam struct {
	StorageSpamCaller     // Read-only binding to the contract
	StorageSpamTransactor // Write-only binding to the contract
	StorageSpamFilterer   // Log filterer for contract events
}

// StorageSpamCaller is an auto generated read-only Go binding around an Ethereum contract.
type StorageSpamCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageSpamTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StorageSpamTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageSpamFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StorageSpamFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageSpamSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StorageSpamSession struct {
	Contract     *StorageSpam      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StorageSpamCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StorageSpamCallerSession struct {
	Contract *StorageSpamCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// StorageSpamTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StorageSpamTransactorSession struct {
	Contract     *StorageSpamTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// StorageSpamRaw is an auto generated low-level Go binding around an Ethereum contract.
type StorageSpamRaw struct {
	Contract *StorageSpam // Generic contract binding to access the raw methods on
}

// StorageSpamCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StorageSpamCallerRaw struct {
	Contract *StorageSpamCaller // Generic read-only contract binding to access the raw methods on
}

// StorageSpamTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StorageSpamTransactorRaw struct {
	Contract *StorageSpamTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStorageSpam creates a new instance of StorageSpam, bound to a specific deployed contract.
func NewStorageSpam(address common.Address, backend bind.ContractBackend) (*StorageSpam, error) {
	contract, err := bindStorageSpam(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StorageSpam{StorageSpamCaller: StorageSpamCaller{contract: contract}, StorageSpamTransactor: StorageSpamTransactor{contract: contract}, StorageSpamFilterer: StorageSpamFilterer{contract: contract}}, nil
}

// NewStorageSpamCaller creates a new read-only instance of StorageSpam, bound to a specific deployed contract.
func NewStorageSpamCaller(address common.Address, caller bind.ContractCaller) (*StorageSpamCaller, error) {
	contract, err := bindStorageSpam(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StorageSpamCaller{contract: contract}, nil
}

// NewStorageSpamTransactor creates a new write-only instance of StorageSpam, bound to a specific deployed contract.
func NewStorageSpamTransactor(address common.Address, transactor bind.ContractTransactor) (*StorageSpamTransactor, error) {
	contract, err := bindStorageSpam(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StorageSpamTransactor{contract: contract}, nil
}

// NewStorageSpamFilterer creates a new log filterer instance of StorageSpam, bound to a specific deployed contract.
func NewStorageSpamFilterer(address common.Address, filterer bind.ContractFilterer) (*StorageSpamFilterer, error) {
	contract, err := bindStorageSpam(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StorageSpamFilterer{contract: contract}, nil
}

// bindStorageSpam binds a generic wrapper to an already deployed contract.
func bindStorageSpam(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StorageSpamMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StorageSpam *StorageSpamRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StorageSpam.Contract.StorageSpamCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StorageSpam *StorageSpamRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StorageSpam.Contract.StorageSpamTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StorageSpam *StorageSpamRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StorageSpam.Contract.StorageSpamTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StorageSpam *StorageSpamCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StorageSpam.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StorageSpam *StorageSpamTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StorageSpam.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StorageSpam *StorageSpamTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StorageSpam.Contract.contract.Transact(opts, method, params...)
}

// GetStorage is a free data retrieval call binding the contract method 0x65fd4772.
//
// Solidity: function getStorage(uint256 key) view returns(uint256)
func (_StorageSpam *StorageSpamCaller) GetStorage(opts *bind.CallOpts, key *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _StorageSpam.contract.Call(opts, &out, "getStorage", key)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetStorage is a free data retrieval call binding the contract method 0x65fd4772.
//
// Solidity: function getStorage(uint256 key) view returns(uint256)
func (_StorageSpam *StorageSpamSession) GetStorage(key *big.Int) (*big.Int, error) {
	return _StorageSpam.Contract.GetStorage(&_StorageSpam.CallOpts, key)
}

// GetStorage is a free data retrieval call binding the contract method 0x65fd4772.
//
// Solidity: function getStorage(uint256 key) view returns(uint256)
func (_StorageSpam *StorageSpamCallerSession) GetStorage(key *big.Int) (*big.Int, error) {
	return _StorageSpam.Contract.GetStorage(&_StorageSpam.CallOpts, key)
}

// StorageMap is a free data retrieval call binding the contract method 0x5e666e4a.
//
// Solidity: function storageMap(uint256 ) view returns(uint256)
func (_StorageSpam *StorageSpamCaller) StorageMap(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _StorageSpam.contract.Call(opts, &out, "storageMap", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StorageMap is a free data retrieval call binding the contract method 0x5e666e4a.
//
// Solidity: function storageMap(uint256 ) view returns(uint256)
func (_StorageSpam *StorageSpamSession) StorageMap(arg0 *big.Int) (*big.Int, error) {
	return _StorageSpam.Contract.StorageMap(&_StorageSpam.CallOpts, arg0)
}

// StorageMap is a free data retrieval call binding the contract method 0x5e666e4a.
//
// Solidity: function storageMap(uint256 ) view returns(uint256)
func (_StorageSpam *StorageSpamCallerSession) StorageMap(arg0 *big.Int) (*big.Int, error) {
	return _StorageSpam.Contract.StorageMap(&_StorageSpam.CallOpts, arg0)
}

// SetRandomForGas is a paid mutator transaction binding the contract method 0xfed72935.
//
// Solidity: function setRandomForGas(uint256 gasLimit, uint256 txid) returns()
func (_StorageSpam *StorageSpamTransactor) SetRandomForGas(opts *bind.TransactOpts, gasLimit *big.Int, txid *big.Int) (*types.Transaction, error) {
	return _StorageSpam.contract.Transact(opts, "setRandomForGas", gasLimit, txid)
}

// SetRandomForGas is a paid mutator transaction binding the contract method 0xfed72935.
//
// Solidity: function setRandomForGas(uint256 gasLimit, uint256 txid) returns()
func (_StorageSpam *StorageSpamSession) SetRandomForGas(gasLimit *big.Int, txid *big.Int) (*types.Transaction, error) {
	return _StorageSpam.Contract.SetRandomForGas(&_StorageSpam.TransactOpts, gasLimit, txid)
}

// SetRandomForGas is a paid mutator transaction binding the contract method 0xfed72935.
//
// Solidity: function setRandomForGas(uint256 gasLimit, uint256 txid) returns()
func (_StorageSpam *StorageSpamTransactorSession) SetRandomForGas(gasLimit *big.Int, txid *big.Int) (*types.Transaction, error) {
	return _StorageSpam.Contract.SetRandomForGas(&_StorageSpam.TransactOpts, gasLimit, txid)
}

// SetStorage is a paid mutator transaction binding the contract method 0x936ad72f.
//
// Solidity: function setStorage(uint256 key, uint256 value) returns()
func (_StorageSpam *StorageSpamTransactor) SetStorage(opts *bind.TransactOpts, key *big.Int, value *big.Int) (*types.Transaction, error) {
	return _StorageSpam.contract.Transact(opts, "setStorage", key, value)
}

// SetStorage is a paid mutator transaction binding the contract method 0x936ad72f.
//
// Solidity: function setStorage(uint256 key, uint256 value) returns()
func (_StorageSpam *StorageSpamSession) SetStorage(key *big.Int, value *big.Int) (*types.Transaction, error) {
	return _StorageSpam.Contract.SetStorage(&_StorageSpam.TransactOpts, key, value)
}

// SetStorage is a paid mutator transaction binding the contract method 0x936ad72f.
//
// Solidity: function setStorage(uint256 key, uint256 value) returns()
func (_StorageSpam *StorageSpamTransactorSession) SetStorage(key *big.Int, value *big.Int) (*types.Transaction, error) {
	return _StorageSpam.Contract.SetStorage(&_StorageSpam.TransactOpts, key, value)
}

// StorageSpamRandomForGasIterator is returned from FilterRandomForGas and is used to iterate over the raw logs and unpacked data for RandomForGas events raised by the StorageSpam contract.
type StorageSpamRandomForGasIterator struct {
	Event *StorageSpamRandomForGas // Event containing the contract specifics and raw log

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
func (it *StorageSpamRandomForGasIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StorageSpamRandomForGas)
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
		it.Event = new(StorageSpamRandomForGas)
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
func (it *StorageSpamRandomForGasIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StorageSpamRandomForGasIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StorageSpamRandomForGas represents a RandomForGas event raised by the StorageSpam contract.
type StorageSpamRandomForGas struct {
	Gas   *big.Int
	Loops *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterRandomForGas is a free log retrieval operation binding the contract event 0x6951887835703bee48a4b856b283f6bb242d53dcd20c2a257a95e9326e594cd1.
//
// Solidity: event RandomForGas(uint256 indexed gas, uint256 loops)
func (_StorageSpam *StorageSpamFilterer) FilterRandomForGas(opts *bind.FilterOpts, gas []*big.Int) (*StorageSpamRandomForGasIterator, error) {

	var gasRule []interface{}
	for _, gasItem := range gas {
		gasRule = append(gasRule, gasItem)
	}

	logs, sub, err := _StorageSpam.contract.FilterLogs(opts, "RandomForGas", gasRule)
	if err != nil {
		return nil, err
	}
	return &StorageSpamRandomForGasIterator{contract: _StorageSpam.contract, event: "RandomForGas", logs: logs, sub: sub}, nil
}

// WatchRandomForGas is a free log subscription operation binding the contract event 0x6951887835703bee48a4b856b283f6bb242d53dcd20c2a257a95e9326e594cd1.
//
// Solidity: event RandomForGas(uint256 indexed gas, uint256 loops)
func (_StorageSpam *StorageSpamFilterer) WatchRandomForGas(opts *bind.WatchOpts, sink chan<- *StorageSpamRandomForGas, gas []*big.Int) (event.Subscription, error) {

	var gasRule []interface{}
	for _, gasItem := range gas {
		gasRule = append(gasRule, gasItem)
	}

	logs, sub, err := _StorageSpam.contract.WatchLogs(opts, "RandomForGas", gasRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StorageSpamRandomForGas)
				if err := _StorageSpam.contract.UnpackLog(event, "RandomForGas", log); err != nil {
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

// ParseRandomForGas is a log parse operation binding the contract event 0x6951887835703bee48a4b856b283f6bb242d53dcd20c2a257a95e9326e594cd1.
//
// Solidity: event RandomForGas(uint256 indexed gas, uint256 loops)
func (_StorageSpam *StorageSpamFilterer) ParseRandomForGas(log types.Log) (*StorageSpamRandomForGas, error) {
	event := new(StorageSpamRandomForGas)
	if err := _StorageSpam.contract.UnpackLog(event, "RandomForGas", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
