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

// ContractMetaData contains all meta data concerning the Contract contract.
var ContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"depositContract\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"_depositContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"pubkey\",\"type\":\"bytes\"}],\"name\":\"topup\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"pubkeys\",\"type\":\"bytes\"}],\"name\":\"topupEqual\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x6080604052604051610a37380380610a378339810160408190526020916043565b5f80546001600160a01b0319166001600160a01b0392909216919091179055606e565b5f602082840312156052575f80fd5b81516001600160a01b03811681146067575f80fd5b9392505050565b6109bc8061007b5f395ff3fe608060405260043610610033575f3560e01c8063089036e0146100375780632d1272121461004c578063c55dc8fa1461005f575b5f80fd5b61004a61004536600461074f565b610099565b005b61004a61005a36600461074f565b6100f8565b34801561006a575f80fd5b505f5461007d906001600160a01b031681565b6040516001600160a01b03909116815260200160405180910390f35b34670de0b6b3a76400008110156100e85760405162461bcd60e51b815260206004820152600e60248201526d616d6f756e7420746f6f206c6f7760901b60448201526064015b60405180910390fd5b6100f38383836101b5565b505050565b805f63ffffffff821661010c3460306107d1565b61011691906107ee565b90505f670de0b6b3a76400008210156101625760405162461bcd60e51b815260206004820152600e60248201526d616d6f756e7420746f6f206c6f7760901b60448201526064016100df565b8263ffffffff168163ffffffff1610156101ae576101a68563ffffffff83168661018d85603061080d565b63ffffffff16926101a093929190610829565b846101b5565b603001610162565b5050505050565b5f806101cd6101c8633b9aca00856107ee565b6105a5565b604080515f60208201819052818301819052606080830182905283518084039091018152608083019093529293509091906002906102139089908990859060a001610850565b60408051601f198184030181529082905261022d9161088e565b602060405180830381855afa158015610248573d5f803e3d5ffd5b5050506040513d601f19601f8201168201806040525081019061026b91906108a0565b604080515f6020820181905291810182905291925090600290819060600160408051601f19818403018152908290526102a39161088e565b602060405180830381855afa1580156102be573d5f803e3d5ffd5b5050506040513d601f19601f820116820180604052508101906102e191906108a0565b604051600290610300905f908190602001918252602082015260400190565b60408051601f198184030181529082905261031a9161088e565b602060405180830381855afa158015610335573d5f803e3d5ffd5b5050506040513d601f19601f8201168201806040525081019061035891906108a0565b60408051602081019390935282015260600160408051601f19818403018152908290526103849161088e565b602060405180830381855afa15801561039f573d5f803e3d5ffd5b5050506040513d601f19601f820116820180604052508101906103c291906108a0565b90505f60028084886040516020016103e4929190918252602082015260400190565b60408051601f19818403018152908290526103fe9161088e565b602060405180830381855afa158015610419573d5f803e3d5ffd5b5050506040513d601f19601f8201168201806040525081019061043c91906108a0565b6040516002906104549089905f9088906020016108b7565b60408051601f198184030181529082905261046e9161088e565b602060405180830381855afa158015610489573d5f803e3d5ffd5b5050506040513d601f19601f820116820180604052508101906104ac91906108a0565b60408051602081019390935282015260600160408051601f19818403018152908290526104d89161088e565b602060405180830381855afa1580156104f3573d5f803e3d5ffd5b5050506040513d601f19601f8201168201806040525081019061051691906108a0565b5f546040805160208082018b905282518083039091018152818301928390526304512a2360e31b9092529293506001600160a01b03909116916322895118918a9161056c918e918e91908b90899060440161090f565b5f604051808303818588803b158015610583575f80fd5b505af1158015610595573d5f803e3d5ffd5b5050505050505050505050505050565b60408051600880825281830190925260609160208201818036833701905050905060c082901b8060071a60f81b825f815181106105e4576105e4610972565b60200101906001600160f81b03191690815f1a9053508060061a60f81b8260018151811061061457610614610972565b60200101906001600160f81b03191690815f1a9053508060051a60f81b8260028151811061064457610644610972565b60200101906001600160f81b03191690815f1a9053508060041a60f81b8260038151811061067457610674610972565b60200101906001600160f81b03191690815f1a9053508060031a60f81b826004815181106106a4576106a4610972565b60200101906001600160f81b03191690815f1a9053508060021a60f81b826005815181106106d4576106d4610972565b60200101906001600160f81b03191690815f1a9053508060011a60f81b8260068151811061070457610704610972565b60200101906001600160f81b03191690815f1a905350805f1a60f81b8260078151811061073357610733610972565b60200101906001600160f81b03191690815f1a90535050919050565b5f8060208385031215610760575f80fd5b823567ffffffffffffffff811115610776575f80fd5b8301601f81018513610786575f80fd5b803567ffffffffffffffff81111561079c575f80fd5b8560208284010111156107ad575f80fd5b6020919091019590945092505050565b634e487b7160e01b5f52601160045260245ffd5b80820281158282048414176107e8576107e86107bd565b92915050565b5f8261080857634e487b7160e01b5f52601260045260245ffd5b500490565b63ffffffff81811683821601908111156107e8576107e86107bd565b5f8085851115610837575f80fd5b83861115610843575f80fd5b5050820193919092039150565b828482376fffffffffffffffffffffffffffffffff19919091169101908152601001919050565b5f81518060208401855e5f93019283525090919050565b5f6108998284610877565b9392505050565b5f602082840312156108b0575f80fd5b5051919050565b5f6108c28286610877565b67ffffffffffffffff1994909416845250506018820152603801919050565b5f81518084528060208401602086015e5f602082860101526020601f19601f83011685010191505092915050565b60808152846080820152848660a08301375f60a086830101525f601f19601f870116820160a083820301602084015261094b60a08201876108e1565b9050828103604084015261095f81866108e1565b9150508260608301529695505050505050565b634e487b7160e01b5f52603260045260245ffdfea2646970667358221220f475e77759100582ace3313dc45b15d88d0eefc572a2264cb1b9b5b54be7a92764736f6c634300081a0033",
}

// ContractABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractMetaData.ABI instead.
var ContractABI = ContractMetaData.ABI

// ContractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractMetaData.Bin instead.
var ContractBin = ContractMetaData.Bin

// DeployContract deploys a new Ethereum contract, binding an instance of Contract to it.
func DeployContract(auth *bind.TransactOpts, backend bind.ContractBackend, depositContract common.Address) (common.Address, *types.Transaction, *Contract, error) {
	parsed, err := ContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractBin), backend, depositContract)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
}

// Contract is an auto generated Go binding around an Ethereum contract.
type Contract struct {
	ContractCaller     // Read-only binding to the contract
	ContractTransactor // Write-only binding to the contract
	ContractFilterer   // Log filterer for contract events
}

// ContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractSession struct {
	Contract     *Contract         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractCallerSession struct {
	Contract *ContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractTransactorSession struct {
	Contract     *ContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractRaw struct {
	Contract *Contract // Generic contract binding to access the raw methods on
}

// ContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractCallerRaw struct {
	Contract *ContractCaller // Generic read-only contract binding to access the raw methods on
}

// ContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractTransactorRaw struct {
	Contract *ContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContract creates a new instance of Contract, bound to a specific deployed contract.
func NewContract(address common.Address, backend bind.ContractBackend) (*Contract, error) {
	contract, err := bindContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
}

// NewContractCaller creates a new read-only instance of Contract, bound to a specific deployed contract.
func NewContractCaller(address common.Address, caller bind.ContractCaller) (*ContractCaller, error) {
	contract, err := bindContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractCaller{contract: contract}, nil
}

// NewContractTransactor creates a new write-only instance of Contract, bound to a specific deployed contract.
func NewContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractTransactor, error) {
	contract, err := bindContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractTransactor{contract: contract}, nil
}

// NewContractFilterer creates a new log filterer instance of Contract, bound to a specific deployed contract.
func NewContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractFilterer, error) {
	contract, err := bindContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractFilterer{contract: contract}, nil
}

// bindContract binds a generic wrapper to an already deployed contract.
func bindContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.ContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transact(opts, method, params...)
}

// DepositContract is a free data retrieval call binding the contract method 0xc55dc8fa.
//
// Solidity: function _depositContract() view returns(address)
func (_Contract *ContractCaller) DepositContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "_depositContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DepositContract is a free data retrieval call binding the contract method 0xc55dc8fa.
//
// Solidity: function _depositContract() view returns(address)
func (_Contract *ContractSession) DepositContract() (common.Address, error) {
	return _Contract.Contract.DepositContract(&_Contract.CallOpts)
}

// DepositContract is a free data retrieval call binding the contract method 0xc55dc8fa.
//
// Solidity: function _depositContract() view returns(address)
func (_Contract *ContractCallerSession) DepositContract() (common.Address, error) {
	return _Contract.Contract.DepositContract(&_Contract.CallOpts)
}

// Topup is a paid mutator transaction binding the contract method 0x089036e0.
//
// Solidity: function topup(bytes pubkey) payable returns()
func (_Contract *ContractTransactor) Topup(opts *bind.TransactOpts, pubkey []byte) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "topup", pubkey)
}

// Topup is a paid mutator transaction binding the contract method 0x089036e0.
//
// Solidity: function topup(bytes pubkey) payable returns()
func (_Contract *ContractSession) Topup(pubkey []byte) (*types.Transaction, error) {
	return _Contract.Contract.Topup(&_Contract.TransactOpts, pubkey)
}

// Topup is a paid mutator transaction binding the contract method 0x089036e0.
//
// Solidity: function topup(bytes pubkey) payable returns()
func (_Contract *ContractTransactorSession) Topup(pubkey []byte) (*types.Transaction, error) {
	return _Contract.Contract.Topup(&_Contract.TransactOpts, pubkey)
}

// TopupEqual is a paid mutator transaction binding the contract method 0x2d127212.
//
// Solidity: function topupEqual(bytes pubkeys) payable returns()
func (_Contract *ContractTransactor) TopupEqual(opts *bind.TransactOpts, pubkeys []byte) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "topupEqual", pubkeys)
}

// TopupEqual is a paid mutator transaction binding the contract method 0x2d127212.
//
// Solidity: function topupEqual(bytes pubkeys) payable returns()
func (_Contract *ContractSession) TopupEqual(pubkeys []byte) (*types.Transaction, error) {
	return _Contract.Contract.TopupEqual(&_Contract.TransactOpts, pubkeys)
}

// TopupEqual is a paid mutator transaction binding the contract method 0x2d127212.
//
// Solidity: function topupEqual(bytes pubkeys) payable returns()
func (_Contract *ContractTransactorSession) TopupEqual(pubkeys []byte) (*types.Transaction, error) {
	return _Contract.Contract.TopupEqual(&_Contract.TransactOpts, pubkeys)
}
