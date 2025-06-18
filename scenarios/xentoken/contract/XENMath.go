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

// XENMathMetaData contains all meta data concerning the XENMath contract.
var XENMathMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"logX64\",\"outputs\":[{\"internalType\":\"int128\",\"name\":\"\",\"type\":\"int128\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"a\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"b\",\"type\":\"uint256\"}],\"name\":\"max\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"a\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"b\",\"type\":\"uint256\"}],\"name\":\"min\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x61026061003a600b82828239805160001a60731461002d57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe730000000000000000000000000000000000000000301460806040526004361061004b5760003560e01c80634793dbab146100505780636d5433e61461007b5780637ae2b5c71461009c575b600080fd5b61006361005e3660046101ef565b6100af565b604051600f9190910b81526020015b60405180910390f35b61008e610089366004610208565b6100c8565b604051908152602001610072565b61008e6100aa366004610208565b6100df565b60006100c26100bd836100f7565b610115565b92915050565b6000818311156100d95750816100c2565b50919050565b6000818311156100f05750806100c2565b5090919050565b6000677fffffffffffffff82111561010e57600080fd5b5060401b90565b60008082600f0b1361012657600080fd5b6000600f83900b600160401b8112610140576040918201911d5b600160201b8112610153576020918201911d5b620100008112610165576010918201911d5b6101008112610176576008918201911d5b60108112610186576004918201911d5b60048112610196576002918201911d5b600281126101a5576001820191505b603f19820160401b600f85900b607f8490031b6001603f1b5b60008113156101e45790800260ff81901c8281029390930192607f011c9060011d6101be565b509095945050505050565b60006020828403121561020157600080fd5b5035919050565b6000806040838503121561021b57600080fd5b5050803592602090910135915056fea2646970667358221220f63a9e5d63fb093bb712d5e405f4e9ca1e008bdd42abe7430bca039ca91172e664736f6c63430008110033",
}

// XENMathABI is the input ABI used to generate the binding from.
// Deprecated: Use XENMathMetaData.ABI instead.
var XENMathABI = XENMathMetaData.ABI

// XENMathBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use XENMathMetaData.Bin instead.
var XENMathBin = XENMathMetaData.Bin

// DeployXENMath deploys a new Ethereum contract, binding an instance of XENMath to it.
func DeployXENMath(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *XENMath, error) {
	parsed, err := XENMathMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(XENMathBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &XENMath{XENMathCaller: XENMathCaller{contract: contract}, XENMathTransactor: XENMathTransactor{contract: contract}, XENMathFilterer: XENMathFilterer{contract: contract}}, nil
}

// XENMath is an auto generated Go binding around an Ethereum contract.
type XENMath struct {
	XENMathCaller     // Read-only binding to the contract
	XENMathTransactor // Write-only binding to the contract
	XENMathFilterer   // Log filterer for contract events
}

// XENMathCaller is an auto generated read-only Go binding around an Ethereum contract.
type XENMathCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// XENMathTransactor is an auto generated write-only Go binding around an Ethereum contract.
type XENMathTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// XENMathFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type XENMathFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// XENMathSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type XENMathSession struct {
	Contract     *XENMath          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// XENMathCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type XENMathCallerSession struct {
	Contract *XENMathCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// XENMathTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type XENMathTransactorSession struct {
	Contract     *XENMathTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// XENMathRaw is an auto generated low-level Go binding around an Ethereum contract.
type XENMathRaw struct {
	Contract *XENMath // Generic contract binding to access the raw methods on
}

// XENMathCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type XENMathCallerRaw struct {
	Contract *XENMathCaller // Generic read-only contract binding to access the raw methods on
}

// XENMathTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type XENMathTransactorRaw struct {
	Contract *XENMathTransactor // Generic write-only contract binding to access the raw methods on
}

// NewXENMath creates a new instance of XENMath, bound to a specific deployed contract.
func NewXENMath(address common.Address, backend bind.ContractBackend) (*XENMath, error) {
	contract, err := bindXENMath(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &XENMath{XENMathCaller: XENMathCaller{contract: contract}, XENMathTransactor: XENMathTransactor{contract: contract}, XENMathFilterer: XENMathFilterer{contract: contract}}, nil
}

// NewXENMathCaller creates a new read-only instance of XENMath, bound to a specific deployed contract.
func NewXENMathCaller(address common.Address, caller bind.ContractCaller) (*XENMathCaller, error) {
	contract, err := bindXENMath(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &XENMathCaller{contract: contract}, nil
}

// NewXENMathTransactor creates a new write-only instance of XENMath, bound to a specific deployed contract.
func NewXENMathTransactor(address common.Address, transactor bind.ContractTransactor) (*XENMathTransactor, error) {
	contract, err := bindXENMath(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &XENMathTransactor{contract: contract}, nil
}

// NewXENMathFilterer creates a new log filterer instance of XENMath, bound to a specific deployed contract.
func NewXENMathFilterer(address common.Address, filterer bind.ContractFilterer) (*XENMathFilterer, error) {
	contract, err := bindXENMath(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &XENMathFilterer{contract: contract}, nil
}

// bindXENMath binds a generic wrapper to an already deployed contract.
func bindXENMath(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := XENMathMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_XENMath *XENMathRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _XENMath.Contract.XENMathCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_XENMath *XENMathRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _XENMath.Contract.XENMathTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_XENMath *XENMathRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _XENMath.Contract.XENMathTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_XENMath *XENMathCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _XENMath.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_XENMath *XENMathTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _XENMath.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_XENMath *XENMathTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _XENMath.Contract.contract.Transact(opts, method, params...)
}

// LogX64 is a free data retrieval call binding the contract method 0x4793dbab.
//
// Solidity: function logX64(uint256 x) pure returns(int128)
func (_XENMath *XENMathCaller) LogX64(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _XENMath.contract.Call(opts, &out, "logX64", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LogX64 is a free data retrieval call binding the contract method 0x4793dbab.
//
// Solidity: function logX64(uint256 x) pure returns(int128)
func (_XENMath *XENMathSession) LogX64(x *big.Int) (*big.Int, error) {
	return _XENMath.Contract.LogX64(&_XENMath.CallOpts, x)
}

// LogX64 is a free data retrieval call binding the contract method 0x4793dbab.
//
// Solidity: function logX64(uint256 x) pure returns(int128)
func (_XENMath *XENMathCallerSession) LogX64(x *big.Int) (*big.Int, error) {
	return _XENMath.Contract.LogX64(&_XENMath.CallOpts, x)
}

// Max is a free data retrieval call binding the contract method 0x6d5433e6.
//
// Solidity: function max(uint256 a, uint256 b) pure returns(uint256)
func (_XENMath *XENMathCaller) Max(opts *bind.CallOpts, a *big.Int, b *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _XENMath.contract.Call(opts, &out, "max", a, b)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Max is a free data retrieval call binding the contract method 0x6d5433e6.
//
// Solidity: function max(uint256 a, uint256 b) pure returns(uint256)
func (_XENMath *XENMathSession) Max(a *big.Int, b *big.Int) (*big.Int, error) {
	return _XENMath.Contract.Max(&_XENMath.CallOpts, a, b)
}

// Max is a free data retrieval call binding the contract method 0x6d5433e6.
//
// Solidity: function max(uint256 a, uint256 b) pure returns(uint256)
func (_XENMath *XENMathCallerSession) Max(a *big.Int, b *big.Int) (*big.Int, error) {
	return _XENMath.Contract.Max(&_XENMath.CallOpts, a, b)
}

// Min is a free data retrieval call binding the contract method 0x7ae2b5c7.
//
// Solidity: function min(uint256 a, uint256 b) pure returns(uint256)
func (_XENMath *XENMathCaller) Min(opts *bind.CallOpts, a *big.Int, b *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _XENMath.contract.Call(opts, &out, "min", a, b)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Min is a free data retrieval call binding the contract method 0x7ae2b5c7.
//
// Solidity: function min(uint256 a, uint256 b) pure returns(uint256)
func (_XENMath *XENMathSession) Min(a *big.Int, b *big.Int) (*big.Int, error) {
	return _XENMath.Contract.Min(&_XENMath.CallOpts, a, b)
}

// Min is a free data retrieval call binding the contract method 0x7ae2b5c7.
//
// Solidity: function min(uint256 a, uint256 b) pure returns(uint256)
func (_XENMath *XENMathCallerSession) Min(a *big.Int, b *big.Int) (*big.Int, error) {
	return _XENMath.Contract.Min(&_XENMath.CallOpts, a, b)
}
