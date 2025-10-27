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

// ERC20BloaterMetaData contains all meta data concerning the ERC20Bloater contract.
var ERC20BloaterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"initialSupply\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startSlot\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"endSlot\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"slotsWritten\",\"type\":\"uint256\"}],\"name\":\"StorageBloated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startSlot\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numAddresses\",\"type\":\"uint256\"}],\"name\":\"bloatStorage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBloatProgress\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nextStorageSlot\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b50604051610697380380610697833981016040819052602b916059565b600480546001600160a01b031916339081179091555f8281559081526002602052604090205560018055606f565b5f602082840312156068575f5ffd5b5051919050565b61061b8061007c5f395ff3fe608060405234801561000f575f5ffd5b50600436106100e5575f3560e01c806370a0823111610088578063a9059cbb11610063578063a9059cbb14610261578063b330b8e9146102a1578063c1926de5146102a9578063dd62ed3e146102bc575f5ffd5b806370a08231146101f35780638da5cb5b1461021257806395d89b411461023d575f5ffd5b806318160ddd116100c357806318160ddd1461015e57806323b872dd14610166578063313ce567146101c457806340c10f19146101de575f5ffd5b806305b3b2b1146100e957806306fdde0314610105578063095ea7b31461013b575b5f5ffd5b6100f260015481565b6040519081526020015b60405180910390f35b61012e6040518060400160405280600a815260200169213637b0ba2a37b5b2b760b11b81525081565b6040516100fc9190610484565b61014e6101493660046104d4565b6102e6565b60405190151581526020016100fc565b6100f25f5481565b61014e6101743660046104fc565b6001600160a01b039283165f818152600360209081526040808320338452825280832080548690039055928252600290528181208054849003905592909316825291902080549091019055600190565b6101cc601281565b60405160ff90911681526020016100fc565b6101f16101ec3660046104d4565b610314565b005b6100f2610201366004610536565b60026020525f908152604090205481565b600454610225906001600160a01b031681565b6040516001600160a01b0390911681526020016100fc565b61012e60405180604001604052806005815260200164109313d05560da1b81525081565b61014e61026f3660046104d4565b335f90815260026020526040808220805484900390556001600160a01b03841682529020805482019055600192915050565b6001546100f2565b6101f16102b7366004610556565b610386565b6100f26102ca366004610576565b600360209081525f928352604080842090915290825290205481565b335f9081526003602090815260408083206001600160a01b0386168452909152902081905560015b92915050565b6004546001600160a01b0316331461035f5760405162461bcd60e51b81526020600482015260096024820152682737ba1037bbb732b960b91b60448201526064015b60405180910390fd5b5f8054820181556001600160a01b0390921682526002602052604090912080549091019055565b6004546001600160a01b031633146103cc5760405162461bcd60e51b81526020600482015260096024820152682737ba1037bbb732b960b91b6044820152606401610356565b5f6103d782846105bb565b9050825b8181101561042d57335f818152600260209081526040808320805486900390556001600160a01b03851680845281842080548701905593835260038252808320938352929052208190556001016103db565b5060018190557f81da33eb4626fe0493d92fb5d273452c79e4c61c8defa0504d1deb856bfba73783826104618560026105ce565b6040805193845260208401929092529082015260600160405180910390a1505050565b602081525f82518060208401528060208501604085015e5f604082850101526040601f19601f83011684010191505092915050565b80356001600160a01b03811681146104cf575f5ffd5b919050565b5f5f604083850312156104e5575f5ffd5b6104ee836104b9565b946020939093013593505050565b5f5f5f6060848603121561050e575f5ffd5b610517846104b9565b9250610525602085016104b9565b929592945050506040919091013590565b5f60208284031215610546575f5ffd5b61054f826104b9565b9392505050565b5f5f60408385031215610567575f5ffd5b50508035926020909101359150565b5f5f60408385031215610587575f5ffd5b610590836104b9565b915061059e602084016104b9565b90509250929050565b634e487b7160e01b5f52601160045260245ffd5b8082018082111561030e5761030e6105a7565b808202811582820484141761030e5761030e6105a756fea26469706673582212200e7ffa6479f21bb6d5341a7c0261b57e5dbb0866f6a603bfa70bb8ac7570b94f64736f6c634300081e0033",
}

// ERC20BloaterABI is the input ABI used to generate the binding from.
// Deprecated: Use ERC20BloaterMetaData.ABI instead.
var ERC20BloaterABI = ERC20BloaterMetaData.ABI

// ERC20BloaterBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ERC20BloaterMetaData.Bin instead.
var ERC20BloaterBin = ERC20BloaterMetaData.Bin

// DeployERC20Bloater deploys a new Ethereum contract, binding an instance of ERC20Bloater to it.
func DeployERC20Bloater(auth *bind.TransactOpts, backend bind.ContractBackend, initialSupply *big.Int) (common.Address, *types.Transaction, *ERC20Bloater, error) {
	parsed, err := ERC20BloaterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ERC20BloaterBin), backend, initialSupply)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ERC20Bloater{ERC20BloaterCaller: ERC20BloaterCaller{contract: contract}, ERC20BloaterTransactor: ERC20BloaterTransactor{contract: contract}, ERC20BloaterFilterer: ERC20BloaterFilterer{contract: contract}}, nil
}

// ERC20Bloater is an auto generated Go binding around an Ethereum contract.
type ERC20Bloater struct {
	ERC20BloaterCaller     // Read-only binding to the contract
	ERC20BloaterTransactor // Write-only binding to the contract
	ERC20BloaterFilterer   // Log filterer for contract events
}

// ERC20BloaterCaller is an auto generated read-only Go binding around an Ethereum contract.
type ERC20BloaterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20BloaterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ERC20BloaterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20BloaterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ERC20BloaterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20BloaterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ERC20BloaterSession struct {
	Contract     *ERC20Bloater     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC20BloaterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ERC20BloaterCallerSession struct {
	Contract *ERC20BloaterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// ERC20BloaterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ERC20BloaterTransactorSession struct {
	Contract     *ERC20BloaterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// ERC20BloaterRaw is an auto generated low-level Go binding around an Ethereum contract.
type ERC20BloaterRaw struct {
	Contract *ERC20Bloater // Generic contract binding to access the raw methods on
}

// ERC20BloaterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ERC20BloaterCallerRaw struct {
	Contract *ERC20BloaterCaller // Generic read-only contract binding to access the raw methods on
}

// ERC20BloaterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ERC20BloaterTransactorRaw struct {
	Contract *ERC20BloaterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC20Bloater creates a new instance of ERC20Bloater, bound to a specific deployed contract.
func NewERC20Bloater(address common.Address, backend bind.ContractBackend) (*ERC20Bloater, error) {
	contract, err := bindERC20Bloater(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC20Bloater{ERC20BloaterCaller: ERC20BloaterCaller{contract: contract}, ERC20BloaterTransactor: ERC20BloaterTransactor{contract: contract}, ERC20BloaterFilterer: ERC20BloaterFilterer{contract: contract}}, nil
}

// NewERC20BloaterCaller creates a new read-only instance of ERC20Bloater, bound to a specific deployed contract.
func NewERC20BloaterCaller(address common.Address, caller bind.ContractCaller) (*ERC20BloaterCaller, error) {
	contract, err := bindERC20Bloater(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20BloaterCaller{contract: contract}, nil
}

// NewERC20BloaterTransactor creates a new write-only instance of ERC20Bloater, bound to a specific deployed contract.
func NewERC20BloaterTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC20BloaterTransactor, error) {
	contract, err := bindERC20Bloater(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20BloaterTransactor{contract: contract}, nil
}

// NewERC20BloaterFilterer creates a new log filterer instance of ERC20Bloater, bound to a specific deployed contract.
func NewERC20BloaterFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC20BloaterFilterer, error) {
	contract, err := bindERC20Bloater(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC20BloaterFilterer{contract: contract}, nil
}

// bindERC20Bloater binds a generic wrapper to an already deployed contract.
func bindERC20Bloater(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ERC20BloaterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20Bloater *ERC20BloaterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC20Bloater.Contract.ERC20BloaterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20Bloater *ERC20BloaterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20Bloater.Contract.ERC20BloaterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20Bloater *ERC20BloaterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20Bloater.Contract.ERC20BloaterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20Bloater *ERC20BloaterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC20Bloater.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20Bloater *ERC20BloaterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20Bloater.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20Bloater *ERC20BloaterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20Bloater.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_ERC20Bloater *ERC20BloaterCaller) Allowance(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ERC20Bloater.contract.Call(opts, &out, "allowance", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_ERC20Bloater *ERC20BloaterSession) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _ERC20Bloater.Contract.Allowance(&_ERC20Bloater.CallOpts, arg0, arg1)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_ERC20Bloater *ERC20BloaterCallerSession) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _ERC20Bloater.Contract.Allowance(&_ERC20Bloater.CallOpts, arg0, arg1)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_ERC20Bloater *ERC20BloaterCaller) BalanceOf(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ERC20Bloater.contract.Call(opts, &out, "balanceOf", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_ERC20Bloater *ERC20BloaterSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _ERC20Bloater.Contract.BalanceOf(&_ERC20Bloater.CallOpts, arg0)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_ERC20Bloater *ERC20BloaterCallerSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _ERC20Bloater.Contract.BalanceOf(&_ERC20Bloater.CallOpts, arg0)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ERC20Bloater *ERC20BloaterCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ERC20Bloater.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ERC20Bloater *ERC20BloaterSession) Decimals() (uint8, error) {
	return _ERC20Bloater.Contract.Decimals(&_ERC20Bloater.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ERC20Bloater *ERC20BloaterCallerSession) Decimals() (uint8, error) {
	return _ERC20Bloater.Contract.Decimals(&_ERC20Bloater.CallOpts)
}

// GetBloatProgress is a free data retrieval call binding the contract method 0xb330b8e9.
//
// Solidity: function getBloatProgress() view returns(uint256)
func (_ERC20Bloater *ERC20BloaterCaller) GetBloatProgress(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ERC20Bloater.contract.Call(opts, &out, "getBloatProgress")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBloatProgress is a free data retrieval call binding the contract method 0xb330b8e9.
//
// Solidity: function getBloatProgress() view returns(uint256)
func (_ERC20Bloater *ERC20BloaterSession) GetBloatProgress() (*big.Int, error) {
	return _ERC20Bloater.Contract.GetBloatProgress(&_ERC20Bloater.CallOpts)
}

// GetBloatProgress is a free data retrieval call binding the contract method 0xb330b8e9.
//
// Solidity: function getBloatProgress() view returns(uint256)
func (_ERC20Bloater *ERC20BloaterCallerSession) GetBloatProgress() (*big.Int, error) {
	return _ERC20Bloater.Contract.GetBloatProgress(&_ERC20Bloater.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ERC20Bloater *ERC20BloaterCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ERC20Bloater.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ERC20Bloater *ERC20BloaterSession) Name() (string, error) {
	return _ERC20Bloater.Contract.Name(&_ERC20Bloater.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ERC20Bloater *ERC20BloaterCallerSession) Name() (string, error) {
	return _ERC20Bloater.Contract.Name(&_ERC20Bloater.CallOpts)
}

// NextStorageSlot is a free data retrieval call binding the contract method 0x05b3b2b1.
//
// Solidity: function nextStorageSlot() view returns(uint256)
func (_ERC20Bloater *ERC20BloaterCaller) NextStorageSlot(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ERC20Bloater.contract.Call(opts, &out, "nextStorageSlot")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextStorageSlot is a free data retrieval call binding the contract method 0x05b3b2b1.
//
// Solidity: function nextStorageSlot() view returns(uint256)
func (_ERC20Bloater *ERC20BloaterSession) NextStorageSlot() (*big.Int, error) {
	return _ERC20Bloater.Contract.NextStorageSlot(&_ERC20Bloater.CallOpts)
}

// NextStorageSlot is a free data retrieval call binding the contract method 0x05b3b2b1.
//
// Solidity: function nextStorageSlot() view returns(uint256)
func (_ERC20Bloater *ERC20BloaterCallerSession) NextStorageSlot() (*big.Int, error) {
	return _ERC20Bloater.Contract.NextStorageSlot(&_ERC20Bloater.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ERC20Bloater *ERC20BloaterCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ERC20Bloater.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ERC20Bloater *ERC20BloaterSession) Owner() (common.Address, error) {
	return _ERC20Bloater.Contract.Owner(&_ERC20Bloater.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ERC20Bloater *ERC20BloaterCallerSession) Owner() (common.Address, error) {
	return _ERC20Bloater.Contract.Owner(&_ERC20Bloater.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ERC20Bloater *ERC20BloaterCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ERC20Bloater.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ERC20Bloater *ERC20BloaterSession) Symbol() (string, error) {
	return _ERC20Bloater.Contract.Symbol(&_ERC20Bloater.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ERC20Bloater *ERC20BloaterCallerSession) Symbol() (string, error) {
	return _ERC20Bloater.Contract.Symbol(&_ERC20Bloater.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ERC20Bloater *ERC20BloaterCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ERC20Bloater.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ERC20Bloater *ERC20BloaterSession) TotalSupply() (*big.Int, error) {
	return _ERC20Bloater.Contract.TotalSupply(&_ERC20Bloater.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ERC20Bloater *ERC20BloaterCallerSession) TotalSupply() (*big.Int, error) {
	return _ERC20Bloater.Contract.TotalSupply(&_ERC20Bloater.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_ERC20Bloater *ERC20BloaterTransactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Bloater.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_ERC20Bloater *ERC20BloaterSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Bloater.Contract.Approve(&_ERC20Bloater.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_ERC20Bloater *ERC20BloaterTransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Bloater.Contract.Approve(&_ERC20Bloater.TransactOpts, spender, amount)
}

// BloatStorage is a paid mutator transaction binding the contract method 0xc1926de5.
//
// Solidity: function bloatStorage(uint256 startSlot, uint256 numAddresses) returns()
func (_ERC20Bloater *ERC20BloaterTransactor) BloatStorage(opts *bind.TransactOpts, startSlot *big.Int, numAddresses *big.Int) (*types.Transaction, error) {
	return _ERC20Bloater.contract.Transact(opts, "bloatStorage", startSlot, numAddresses)
}

// BloatStorage is a paid mutator transaction binding the contract method 0xc1926de5.
//
// Solidity: function bloatStorage(uint256 startSlot, uint256 numAddresses) returns()
func (_ERC20Bloater *ERC20BloaterSession) BloatStorage(startSlot *big.Int, numAddresses *big.Int) (*types.Transaction, error) {
	return _ERC20Bloater.Contract.BloatStorage(&_ERC20Bloater.TransactOpts, startSlot, numAddresses)
}

// BloatStorage is a paid mutator transaction binding the contract method 0xc1926de5.
//
// Solidity: function bloatStorage(uint256 startSlot, uint256 numAddresses) returns()
func (_ERC20Bloater *ERC20BloaterTransactorSession) BloatStorage(startSlot *big.Int, numAddresses *big.Int) (*types.Transaction, error) {
	return _ERC20Bloater.Contract.BloatStorage(&_ERC20Bloater.TransactOpts, startSlot, numAddresses)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_ERC20Bloater *ERC20BloaterTransactor) Mint(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Bloater.contract.Transact(opts, "mint", to, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_ERC20Bloater *ERC20BloaterSession) Mint(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Bloater.Contract.Mint(&_ERC20Bloater.TransactOpts, to, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_ERC20Bloater *ERC20BloaterTransactorSession) Mint(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Bloater.Contract.Mint(&_ERC20Bloater.TransactOpts, to, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_ERC20Bloater *ERC20BloaterTransactor) Transfer(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Bloater.contract.Transact(opts, "transfer", to, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_ERC20Bloater *ERC20BloaterSession) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Bloater.Contract.Transfer(&_ERC20Bloater.TransactOpts, to, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_ERC20Bloater *ERC20BloaterTransactorSession) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Bloater.Contract.Transfer(&_ERC20Bloater.TransactOpts, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_ERC20Bloater *ERC20BloaterTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Bloater.contract.Transact(opts, "transferFrom", from, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_ERC20Bloater *ERC20BloaterSession) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Bloater.Contract.TransferFrom(&_ERC20Bloater.TransactOpts, from, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_ERC20Bloater *ERC20BloaterTransactorSession) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Bloater.Contract.TransferFrom(&_ERC20Bloater.TransactOpts, from, to, amount)
}

// ERC20BloaterStorageBloatedIterator is returned from FilterStorageBloated and is used to iterate over the raw logs and unpacked data for StorageBloated events raised by the ERC20Bloater contract.
type ERC20BloaterStorageBloatedIterator struct {
	Event *ERC20BloaterStorageBloated // Event containing the contract specifics and raw log

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
func (it *ERC20BloaterStorageBloatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20BloaterStorageBloated)
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
		it.Event = new(ERC20BloaterStorageBloated)
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
func (it *ERC20BloaterStorageBloatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20BloaterStorageBloatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20BloaterStorageBloated represents a StorageBloated event raised by the ERC20Bloater contract.
type ERC20BloaterStorageBloated struct {
	StartSlot    *big.Int
	EndSlot      *big.Int
	SlotsWritten *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterStorageBloated is a free log retrieval operation binding the contract event 0x81da33eb4626fe0493d92fb5d273452c79e4c61c8defa0504d1deb856bfba737.
//
// Solidity: event StorageBloated(uint256 startSlot, uint256 endSlot, uint256 slotsWritten)
func (_ERC20Bloater *ERC20BloaterFilterer) FilterStorageBloated(opts *bind.FilterOpts) (*ERC20BloaterStorageBloatedIterator, error) {

	logs, sub, err := _ERC20Bloater.contract.FilterLogs(opts, "StorageBloated")
	if err != nil {
		return nil, err
	}
	return &ERC20BloaterStorageBloatedIterator{contract: _ERC20Bloater.contract, event: "StorageBloated", logs: logs, sub: sub}, nil
}

// WatchStorageBloated is a free log subscription operation binding the contract event 0x81da33eb4626fe0493d92fb5d273452c79e4c61c8defa0504d1deb856bfba737.
//
// Solidity: event StorageBloated(uint256 startSlot, uint256 endSlot, uint256 slotsWritten)
func (_ERC20Bloater *ERC20BloaterFilterer) WatchStorageBloated(opts *bind.WatchOpts, sink chan<- *ERC20BloaterStorageBloated) (event.Subscription, error) {

	logs, sub, err := _ERC20Bloater.contract.WatchLogs(opts, "StorageBloated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20BloaterStorageBloated)
				if err := _ERC20Bloater.contract.UnpackLog(event, "StorageBloated", log); err != nil {
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

// ParseStorageBloated is a log parse operation binding the contract event 0x81da33eb4626fe0493d92fb5d273452c79e4c61c8defa0504d1deb856bfba737.
//
// Solidity: event StorageBloated(uint256 startSlot, uint256 endSlot, uint256 slotsWritten)
func (_ERC20Bloater *ERC20BloaterFilterer) ParseStorageBloated(log types.Log) (*ERC20BloaterStorageBloated, error) {
	event := new(ERC20BloaterStorageBloated)
	if err := _ERC20Bloater.contract.UnpackLog(event, "StorageBloated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
