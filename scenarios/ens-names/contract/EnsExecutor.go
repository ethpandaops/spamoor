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

// EnsExecutorMetaData contains all meta data concerning the EnsExecutor contract.
var EnsExecutorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_admin\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"Deployed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"name\":\"OperatorChanged\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"initcode\",\"type\":\"bytes\"}],\"name\":\"deploy\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"execute\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"operators\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_admin\",\"type\":\"address\"}],\"name\":\"setAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"name\":\"setOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x608060405234801561000f575f80fd5b5060405161078538038061078583398101604081905261002e91610052565b5f80546001600160a01b0319166001600160a01b039290921691909117905561007f565b5f60208284031215610062575f80fd5b81516001600160a01b0381168114610078575f80fd5b9392505050565b6106f98061008c5f395ff3fe608060405260043610610057575f3560e01c806313e7c9d814610062578063558a7297146100a5578063704b6c02146100c6578063b61d27f6146100e5578063cdcb760a14610105578063f851a4401461013c575f80fd5b3661005e57005b5f80fd5b34801561006d575f80fd5b5061009061007c3660046104c6565b60016020525f908152604090205460ff1681565b60405190151581526020015b60405180910390f35b3480156100b0575f80fd5b506100c46100bf3660046104e6565b61015a565b005b3480156100d1575f80fd5b506100c46100e03660046104c6565b61020f565b6100f86100f336600461051f565b6102a8565b60405161009c919061059f565b348015610110575f80fd5b5061012461011f3660046105ff565b610395565b6040516001600160a01b03909116815260200161009c565b348015610147575f80fd5b505f54610124906001600160a01b031681565b5f546001600160a01b031633146101b15760405162461bcd60e51b815260206004820152601660248201527522b739a2bc32b1baba37b91d103737ba1030b236b4b760511b60448201526064015b60405180910390fd5b6001600160a01b0382165f81815260016020908152604091829020805460ff191685151590811790915591519182527f193de8d500b5cb7b720089b258a39e9c1d0b840019a73ae7c51c3f9101732b02910160405180910390a25050565b5f546001600160a01b031633146102615760405162461bcd60e51b815260206004820152601660248201527522b739a2bc32b1baba37b91d103737ba1030b236b4b760511b60448201526064016101a8565b5f80546001600160a01b0319166001600160a01b038316908117825560405190917f7ce7ec0b50378fb6c0186ffb5f48325f6593fcb4ca4386f21861af3129188f5c91a250565b5f546060906001600160a01b03163314806102d15750335f9081526001602052604090205460ff165b61031d5760405162461bcd60e51b815260206004820152601b60248201527f456e734578656375746f723a206e6f7420617574686f72697a6564000000000060448201526064016101a8565b5f80866001600160a01b031686868660405161033a9291906106b4565b5f6040518083038185875af1925050503d805f8114610374576040519150601f19603f3d011682016040523d82523d5f602084013e610379565b606091505b50915091508161038b57805160208201fd5b9695505050505050565b5f80546001600160a01b03163314806103bc5750335f9081526001602052604090205460ff165b6104085760405162461bcd60e51b815260206004820152601b60248201527f456e734578656375746f723a206e6f7420617574686f72697a6564000000000060448201526064016101a8565b828251602084015ff590506001600160a01b0381166104695760405162461bcd60e51b815260206004820152601b60248201527f456e734578656375746f723a2063726561746532206661696c6564000000000060448201526064016101a8565b6040516001600160a01b03821681527ff40fcec21964ffb566044d083b4073f29f7f7929110ea19e1b3ebe375d89055e9060200160405180910390a192915050565b80356001600160a01b03811681146104c1575f80fd5b919050565b5f602082840312156104d6575f80fd5b6104df826104ab565b9392505050565b5f80604083850312156104f7575f80fd5b610500836104ab565b915060208301358015158114610514575f80fd5b809150509250929050565b5f805f8060608587031215610532575f80fd5b61053b856104ab565b935060208501359250604085013567ffffffffffffffff8082111561055e575f80fd5b818701915087601f830112610571575f80fd5b81358181111561057f575f80fd5b886020828501011115610590575f80fd5b95989497505060200194505050565b5f602080835283518060208501525f5b818110156105cb578581018301518582016040015282016105af565b505f604082860101526040601f19601f8301168501019250505092915050565b634e487b7160e01b5f52604160045260245ffd5b5f8060408385031215610610575f80fd5b82359150602083013567ffffffffffffffff8082111561062e575f80fd5b818501915085601f830112610641575f80fd5b813581811115610653576106536105eb565b604051601f8201601f19908116603f0116810190838211818310171561067b5761067b6105eb565b81604052828152886020848701011115610693575f80fd5b826020860160208301375f6020848301015280955050505050509250929050565b818382375f910190815291905056fea26469706673582212201e1ce10b91c370642bd97d78a3073461b9340071f8fe31364f030e3b978dd15464736f6c63430008180033",
}

// EnsExecutorABI is the input ABI used to generate the binding from.
// Deprecated: Use EnsExecutorMetaData.ABI instead.
var EnsExecutorABI = EnsExecutorMetaData.ABI

// EnsExecutorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use EnsExecutorMetaData.Bin instead.
var EnsExecutorBin = EnsExecutorMetaData.Bin

// DeployEnsExecutor deploys a new Ethereum contract, binding an instance of EnsExecutor to it.
func DeployEnsExecutor(auth *bind.TransactOpts, backend bind.ContractBackend, _admin common.Address) (common.Address, *types.Transaction, *EnsExecutor, error) {
	parsed, err := EnsExecutorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(EnsExecutorBin), backend, _admin)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &EnsExecutor{EnsExecutorCaller: EnsExecutorCaller{contract: contract}, EnsExecutorTransactor: EnsExecutorTransactor{contract: contract}, EnsExecutorFilterer: EnsExecutorFilterer{contract: contract}}, nil
}

// EnsExecutor is an auto generated Go binding around an Ethereum contract.
type EnsExecutor struct {
	EnsExecutorCaller     // Read-only binding to the contract
	EnsExecutorTransactor // Write-only binding to the contract
	EnsExecutorFilterer   // Log filterer for contract events
}

// EnsExecutorCaller is an auto generated read-only Go binding around an Ethereum contract.
type EnsExecutorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EnsExecutorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EnsExecutorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EnsExecutorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EnsExecutorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EnsExecutorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EnsExecutorSession struct {
	Contract     *EnsExecutor      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EnsExecutorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EnsExecutorCallerSession struct {
	Contract *EnsExecutorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// EnsExecutorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EnsExecutorTransactorSession struct {
	Contract     *EnsExecutorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// EnsExecutorRaw is an auto generated low-level Go binding around an Ethereum contract.
type EnsExecutorRaw struct {
	Contract *EnsExecutor // Generic contract binding to access the raw methods on
}

// EnsExecutorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EnsExecutorCallerRaw struct {
	Contract *EnsExecutorCaller // Generic read-only contract binding to access the raw methods on
}

// EnsExecutorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EnsExecutorTransactorRaw struct {
	Contract *EnsExecutorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEnsExecutor creates a new instance of EnsExecutor, bound to a specific deployed contract.
func NewEnsExecutor(address common.Address, backend bind.ContractBackend) (*EnsExecutor, error) {
	contract, err := bindEnsExecutor(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EnsExecutor{EnsExecutorCaller: EnsExecutorCaller{contract: contract}, EnsExecutorTransactor: EnsExecutorTransactor{contract: contract}, EnsExecutorFilterer: EnsExecutorFilterer{contract: contract}}, nil
}

// NewEnsExecutorCaller creates a new read-only instance of EnsExecutor, bound to a specific deployed contract.
func NewEnsExecutorCaller(address common.Address, caller bind.ContractCaller) (*EnsExecutorCaller, error) {
	contract, err := bindEnsExecutor(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EnsExecutorCaller{contract: contract}, nil
}

// NewEnsExecutorTransactor creates a new write-only instance of EnsExecutor, bound to a specific deployed contract.
func NewEnsExecutorTransactor(address common.Address, transactor bind.ContractTransactor) (*EnsExecutorTransactor, error) {
	contract, err := bindEnsExecutor(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EnsExecutorTransactor{contract: contract}, nil
}

// NewEnsExecutorFilterer creates a new log filterer instance of EnsExecutor, bound to a specific deployed contract.
func NewEnsExecutorFilterer(address common.Address, filterer bind.ContractFilterer) (*EnsExecutorFilterer, error) {
	contract, err := bindEnsExecutor(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EnsExecutorFilterer{contract: contract}, nil
}

// bindEnsExecutor binds a generic wrapper to an already deployed contract.
func bindEnsExecutor(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := EnsExecutorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EnsExecutor *EnsExecutorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EnsExecutor.Contract.EnsExecutorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EnsExecutor *EnsExecutorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EnsExecutor.Contract.EnsExecutorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EnsExecutor *EnsExecutorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EnsExecutor.Contract.EnsExecutorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EnsExecutor *EnsExecutorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EnsExecutor.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EnsExecutor *EnsExecutorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EnsExecutor.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EnsExecutor *EnsExecutorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EnsExecutor.Contract.contract.Transact(opts, method, params...)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_EnsExecutor *EnsExecutorCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EnsExecutor.contract.Call(opts, &out, "admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_EnsExecutor *EnsExecutorSession) Admin() (common.Address, error) {
	return _EnsExecutor.Contract.Admin(&_EnsExecutor.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_EnsExecutor *EnsExecutorCallerSession) Admin() (common.Address, error) {
	return _EnsExecutor.Contract.Admin(&_EnsExecutor.CallOpts)
}

// Operators is a free data retrieval call binding the contract method 0x13e7c9d8.
//
// Solidity: function operators(address ) view returns(bool)
func (_EnsExecutor *EnsExecutorCaller) Operators(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _EnsExecutor.contract.Call(opts, &out, "operators", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Operators is a free data retrieval call binding the contract method 0x13e7c9d8.
//
// Solidity: function operators(address ) view returns(bool)
func (_EnsExecutor *EnsExecutorSession) Operators(arg0 common.Address) (bool, error) {
	return _EnsExecutor.Contract.Operators(&_EnsExecutor.CallOpts, arg0)
}

// Operators is a free data retrieval call binding the contract method 0x13e7c9d8.
//
// Solidity: function operators(address ) view returns(bool)
func (_EnsExecutor *EnsExecutorCallerSession) Operators(arg0 common.Address) (bool, error) {
	return _EnsExecutor.Contract.Operators(&_EnsExecutor.CallOpts, arg0)
}

// Deploy is a paid mutator transaction binding the contract method 0xcdcb760a.
//
// Solidity: function deploy(bytes32 salt, bytes initcode) returns(address addr)
func (_EnsExecutor *EnsExecutorTransactor) Deploy(opts *bind.TransactOpts, salt [32]byte, initcode []byte) (*types.Transaction, error) {
	return _EnsExecutor.contract.Transact(opts, "deploy", salt, initcode)
}

// Deploy is a paid mutator transaction binding the contract method 0xcdcb760a.
//
// Solidity: function deploy(bytes32 salt, bytes initcode) returns(address addr)
func (_EnsExecutor *EnsExecutorSession) Deploy(salt [32]byte, initcode []byte) (*types.Transaction, error) {
	return _EnsExecutor.Contract.Deploy(&_EnsExecutor.TransactOpts, salt, initcode)
}

// Deploy is a paid mutator transaction binding the contract method 0xcdcb760a.
//
// Solidity: function deploy(bytes32 salt, bytes initcode) returns(address addr)
func (_EnsExecutor *EnsExecutorTransactorSession) Deploy(salt [32]byte, initcode []byte) (*types.Transaction, error) {
	return _EnsExecutor.Contract.Deploy(&_EnsExecutor.TransactOpts, salt, initcode)
}

// Execute is a paid mutator transaction binding the contract method 0xb61d27f6.
//
// Solidity: function execute(address target, uint256 value, bytes data) payable returns(bytes)
func (_EnsExecutor *EnsExecutorTransactor) Execute(opts *bind.TransactOpts, target common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _EnsExecutor.contract.Transact(opts, "execute", target, value, data)
}

// Execute is a paid mutator transaction binding the contract method 0xb61d27f6.
//
// Solidity: function execute(address target, uint256 value, bytes data) payable returns(bytes)
func (_EnsExecutor *EnsExecutorSession) Execute(target common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _EnsExecutor.Contract.Execute(&_EnsExecutor.TransactOpts, target, value, data)
}

// Execute is a paid mutator transaction binding the contract method 0xb61d27f6.
//
// Solidity: function execute(address target, uint256 value, bytes data) payable returns(bytes)
func (_EnsExecutor *EnsExecutorTransactorSession) Execute(target common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _EnsExecutor.Contract.Execute(&_EnsExecutor.TransactOpts, target, value, data)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address _admin) returns()
func (_EnsExecutor *EnsExecutorTransactor) SetAdmin(opts *bind.TransactOpts, _admin common.Address) (*types.Transaction, error) {
	return _EnsExecutor.contract.Transact(opts, "setAdmin", _admin)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address _admin) returns()
func (_EnsExecutor *EnsExecutorSession) SetAdmin(_admin common.Address) (*types.Transaction, error) {
	return _EnsExecutor.Contract.SetAdmin(&_EnsExecutor.TransactOpts, _admin)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address _admin) returns()
func (_EnsExecutor *EnsExecutorTransactorSession) SetAdmin(_admin common.Address) (*types.Transaction, error) {
	return _EnsExecutor.Contract.SetAdmin(&_EnsExecutor.TransactOpts, _admin)
}

// SetOperator is a paid mutator transaction binding the contract method 0x558a7297.
//
// Solidity: function setOperator(address operator, bool enabled) returns()
func (_EnsExecutor *EnsExecutorTransactor) SetOperator(opts *bind.TransactOpts, operator common.Address, enabled bool) (*types.Transaction, error) {
	return _EnsExecutor.contract.Transact(opts, "setOperator", operator, enabled)
}

// SetOperator is a paid mutator transaction binding the contract method 0x558a7297.
//
// Solidity: function setOperator(address operator, bool enabled) returns()
func (_EnsExecutor *EnsExecutorSession) SetOperator(operator common.Address, enabled bool) (*types.Transaction, error) {
	return _EnsExecutor.Contract.SetOperator(&_EnsExecutor.TransactOpts, operator, enabled)
}

// SetOperator is a paid mutator transaction binding the contract method 0x558a7297.
//
// Solidity: function setOperator(address operator, bool enabled) returns()
func (_EnsExecutor *EnsExecutorTransactorSession) SetOperator(operator common.Address, enabled bool) (*types.Transaction, error) {
	return _EnsExecutor.Contract.SetOperator(&_EnsExecutor.TransactOpts, operator, enabled)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_EnsExecutor *EnsExecutorTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EnsExecutor.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_EnsExecutor *EnsExecutorSession) Receive() (*types.Transaction, error) {
	return _EnsExecutor.Contract.Receive(&_EnsExecutor.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_EnsExecutor *EnsExecutorTransactorSession) Receive() (*types.Transaction, error) {
	return _EnsExecutor.Contract.Receive(&_EnsExecutor.TransactOpts)
}

// EnsExecutorAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the EnsExecutor contract.
type EnsExecutorAdminChangedIterator struct {
	Event *EnsExecutorAdminChanged // Event containing the contract specifics and raw log

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
func (it *EnsExecutorAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EnsExecutorAdminChanged)
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
		it.Event = new(EnsExecutorAdminChanged)
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
func (it *EnsExecutorAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EnsExecutorAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EnsExecutorAdminChanged represents a AdminChanged event raised by the EnsExecutor contract.
type EnsExecutorAdminChanged struct {
	NewAdmin common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7ce7ec0b50378fb6c0186ffb5f48325f6593fcb4ca4386f21861af3129188f5c.
//
// Solidity: event AdminChanged(address indexed newAdmin)
func (_EnsExecutor *EnsExecutorFilterer) FilterAdminChanged(opts *bind.FilterOpts, newAdmin []common.Address) (*EnsExecutorAdminChangedIterator, error) {

	var newAdminRule []interface{}
	for _, newAdminItem := range newAdmin {
		newAdminRule = append(newAdminRule, newAdminItem)
	}

	logs, sub, err := _EnsExecutor.contract.FilterLogs(opts, "AdminChanged", newAdminRule)
	if err != nil {
		return nil, err
	}
	return &EnsExecutorAdminChangedIterator{contract: _EnsExecutor.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7ce7ec0b50378fb6c0186ffb5f48325f6593fcb4ca4386f21861af3129188f5c.
//
// Solidity: event AdminChanged(address indexed newAdmin)
func (_EnsExecutor *EnsExecutorFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *EnsExecutorAdminChanged, newAdmin []common.Address) (event.Subscription, error) {

	var newAdminRule []interface{}
	for _, newAdminItem := range newAdmin {
		newAdminRule = append(newAdminRule, newAdminItem)
	}

	logs, sub, err := _EnsExecutor.contract.WatchLogs(opts, "AdminChanged", newAdminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EnsExecutorAdminChanged)
				if err := _EnsExecutor.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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

// ParseAdminChanged is a log parse operation binding the contract event 0x7ce7ec0b50378fb6c0186ffb5f48325f6593fcb4ca4386f21861af3129188f5c.
//
// Solidity: event AdminChanged(address indexed newAdmin)
func (_EnsExecutor *EnsExecutorFilterer) ParseAdminChanged(log types.Log) (*EnsExecutorAdminChanged, error) {
	event := new(EnsExecutorAdminChanged)
	if err := _EnsExecutor.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EnsExecutorDeployedIterator is returned from FilterDeployed and is used to iterate over the raw logs and unpacked data for Deployed events raised by the EnsExecutor contract.
type EnsExecutorDeployedIterator struct {
	Event *EnsExecutorDeployed // Event containing the contract specifics and raw log

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
func (it *EnsExecutorDeployedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EnsExecutorDeployed)
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
		it.Event = new(EnsExecutorDeployed)
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
func (it *EnsExecutorDeployedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EnsExecutorDeployedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EnsExecutorDeployed represents a Deployed event raised by the EnsExecutor contract.
type EnsExecutorDeployed struct {
	Addr common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterDeployed is a free log retrieval operation binding the contract event 0xf40fcec21964ffb566044d083b4073f29f7f7929110ea19e1b3ebe375d89055e.
//
// Solidity: event Deployed(address addr)
func (_EnsExecutor *EnsExecutorFilterer) FilterDeployed(opts *bind.FilterOpts) (*EnsExecutorDeployedIterator, error) {

	logs, sub, err := _EnsExecutor.contract.FilterLogs(opts, "Deployed")
	if err != nil {
		return nil, err
	}
	return &EnsExecutorDeployedIterator{contract: _EnsExecutor.contract, event: "Deployed", logs: logs, sub: sub}, nil
}

// WatchDeployed is a free log subscription operation binding the contract event 0xf40fcec21964ffb566044d083b4073f29f7f7929110ea19e1b3ebe375d89055e.
//
// Solidity: event Deployed(address addr)
func (_EnsExecutor *EnsExecutorFilterer) WatchDeployed(opts *bind.WatchOpts, sink chan<- *EnsExecutorDeployed) (event.Subscription, error) {

	logs, sub, err := _EnsExecutor.contract.WatchLogs(opts, "Deployed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EnsExecutorDeployed)
				if err := _EnsExecutor.contract.UnpackLog(event, "Deployed", log); err != nil {
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

// ParseDeployed is a log parse operation binding the contract event 0xf40fcec21964ffb566044d083b4073f29f7f7929110ea19e1b3ebe375d89055e.
//
// Solidity: event Deployed(address addr)
func (_EnsExecutor *EnsExecutorFilterer) ParseDeployed(log types.Log) (*EnsExecutorDeployed, error) {
	event := new(EnsExecutorDeployed)
	if err := _EnsExecutor.contract.UnpackLog(event, "Deployed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EnsExecutorOperatorChangedIterator is returned from FilterOperatorChanged and is used to iterate over the raw logs and unpacked data for OperatorChanged events raised by the EnsExecutor contract.
type EnsExecutorOperatorChangedIterator struct {
	Event *EnsExecutorOperatorChanged // Event containing the contract specifics and raw log

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
func (it *EnsExecutorOperatorChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EnsExecutorOperatorChanged)
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
		it.Event = new(EnsExecutorOperatorChanged)
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
func (it *EnsExecutorOperatorChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EnsExecutorOperatorChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EnsExecutorOperatorChanged represents a OperatorChanged event raised by the EnsExecutor contract.
type EnsExecutorOperatorChanged struct {
	Operator common.Address
	Enabled  bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOperatorChanged is a free log retrieval operation binding the contract event 0x193de8d500b5cb7b720089b258a39e9c1d0b840019a73ae7c51c3f9101732b02.
//
// Solidity: event OperatorChanged(address indexed operator, bool enabled)
func (_EnsExecutor *EnsExecutorFilterer) FilterOperatorChanged(opts *bind.FilterOpts, operator []common.Address) (*EnsExecutorOperatorChangedIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _EnsExecutor.contract.FilterLogs(opts, "OperatorChanged", operatorRule)
	if err != nil {
		return nil, err
	}
	return &EnsExecutorOperatorChangedIterator{contract: _EnsExecutor.contract, event: "OperatorChanged", logs: logs, sub: sub}, nil
}

// WatchOperatorChanged is a free log subscription operation binding the contract event 0x193de8d500b5cb7b720089b258a39e9c1d0b840019a73ae7c51c3f9101732b02.
//
// Solidity: event OperatorChanged(address indexed operator, bool enabled)
func (_EnsExecutor *EnsExecutorFilterer) WatchOperatorChanged(opts *bind.WatchOpts, sink chan<- *EnsExecutorOperatorChanged, operator []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _EnsExecutor.contract.WatchLogs(opts, "OperatorChanged", operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EnsExecutorOperatorChanged)
				if err := _EnsExecutor.contract.UnpackLog(event, "OperatorChanged", log); err != nil {
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

// ParseOperatorChanged is a log parse operation binding the contract event 0x193de8d500b5cb7b720089b258a39e9c1d0b840019a73ae7c51c3f9101732b02.
//
// Solidity: event OperatorChanged(address indexed operator, bool enabled)
func (_EnsExecutor *EnsExecutorFilterer) ParseOperatorChanged(log types.Log) (*EnsExecutorOperatorChanged, error) {
	event := new(EnsExecutorOperatorChanged)
	if err := _EnsExecutor.contract.UnpackLog(event, "OperatorChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
