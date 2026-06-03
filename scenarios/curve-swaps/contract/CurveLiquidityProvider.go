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

// CurveLiquidityProviderMetaData contains all meta data concerning the CurveLiquidityProvider contract.
var CurveLiquidityProviderMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"},{\"internalType\":\"address[3]\",\"name\":\"coins\",\"type\":\"address[3]\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"seedLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610379806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c80632894eca414610030575b600080fd5b61004361003e366004610247565b610045565b005b61004d61020d565b60005b60038110156101a55783816003811061006b5761006b61028a565b60200201602081019061007e91906102a0565b6040516340c10f1960e01b8152306004820152602481018590526001600160a01b0391909116906340c10f1990604401600060405180830381600087803b1580156100c857600080fd5b505af11580156100dc573d6000803e3d6000fd5b505050508381600381106100f2576100f261028a565b60200201602081019061010591906102a0565b60405163095ea7b360e01b81526001600160a01b03878116600483015260248201869052919091169063095ea7b3906044016020604051808303816000875af1158015610156573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061017a91906102c2565b508282826003811061018e5761018e61028a565b60200201528061019d816102e4565b915050610050565b50604051634515cef360e01b81526001600160a01b03851690634515cef3906101d590849060009060040161030b565b600060405180830381600087803b1580156101ef57600080fd5b505af1158015610203573d6000803e3d6000fd5b5050505050505050565b60405180606001604052806003906020820280368337509192915050565b80356001600160a01b038116811461024257600080fd5b919050565b600080600060a0848603121561025c57600080fd5b6102658461022b565b9250608084018581111561027857600080fd5b60208501925080359150509250925092565b634e487b7160e01b600052603260045260246000fd5b6000602082840312156102b257600080fd5b6102bb8261022b565b9392505050565b6000602082840312156102d457600080fd5b815180151581146102bb57600080fd5b60006001820161030457634e487b7160e01b600052601160045260246000fd5b5060010190565b60808101818460005b6003811015610333578151835260209283019290910190600101610314565b505050826060830152939250505056fea2646970667358221220f3a73a0b3a2cebd5915178acb93f1953c52fae1ad7a1e9a068df19af63069bd164736f6c63430008110033",
}

// CurveLiquidityProviderABI is the input ABI used to generate the binding from.
// Deprecated: Use CurveLiquidityProviderMetaData.ABI instead.
var CurveLiquidityProviderABI = CurveLiquidityProviderMetaData.ABI

// CurveLiquidityProviderBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use CurveLiquidityProviderMetaData.Bin instead.
var CurveLiquidityProviderBin = CurveLiquidityProviderMetaData.Bin

// DeployCurveLiquidityProvider deploys a new Ethereum contract, binding an instance of CurveLiquidityProvider to it.
func DeployCurveLiquidityProvider(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *CurveLiquidityProvider, error) {
	parsed, err := CurveLiquidityProviderMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CurveLiquidityProviderBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CurveLiquidityProvider{CurveLiquidityProviderCaller: CurveLiquidityProviderCaller{contract: contract}, CurveLiquidityProviderTransactor: CurveLiquidityProviderTransactor{contract: contract}, CurveLiquidityProviderFilterer: CurveLiquidityProviderFilterer{contract: contract}}, nil
}

// CurveLiquidityProvider is an auto generated Go binding around an Ethereum contract.
type CurveLiquidityProvider struct {
	CurveLiquidityProviderCaller     // Read-only binding to the contract
	CurveLiquidityProviderTransactor // Write-only binding to the contract
	CurveLiquidityProviderFilterer   // Log filterer for contract events
}

// CurveLiquidityProviderCaller is an auto generated read-only Go binding around an Ethereum contract.
type CurveLiquidityProviderCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CurveLiquidityProviderTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CurveLiquidityProviderTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CurveLiquidityProviderFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CurveLiquidityProviderFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CurveLiquidityProviderSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CurveLiquidityProviderSession struct {
	Contract     *CurveLiquidityProvider // Generic contract binding to set the session for
	CallOpts     bind.CallOpts           // Call options to use throughout this session
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// CurveLiquidityProviderCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CurveLiquidityProviderCallerSession struct {
	Contract *CurveLiquidityProviderCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                 // Call options to use throughout this session
}

// CurveLiquidityProviderTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CurveLiquidityProviderTransactorSession struct {
	Contract     *CurveLiquidityProviderTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// CurveLiquidityProviderRaw is an auto generated low-level Go binding around an Ethereum contract.
type CurveLiquidityProviderRaw struct {
	Contract *CurveLiquidityProvider // Generic contract binding to access the raw methods on
}

// CurveLiquidityProviderCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CurveLiquidityProviderCallerRaw struct {
	Contract *CurveLiquidityProviderCaller // Generic read-only contract binding to access the raw methods on
}

// CurveLiquidityProviderTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CurveLiquidityProviderTransactorRaw struct {
	Contract *CurveLiquidityProviderTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCurveLiquidityProvider creates a new instance of CurveLiquidityProvider, bound to a specific deployed contract.
func NewCurveLiquidityProvider(address common.Address, backend bind.ContractBackend) (*CurveLiquidityProvider, error) {
	contract, err := bindCurveLiquidityProvider(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CurveLiquidityProvider{CurveLiquidityProviderCaller: CurveLiquidityProviderCaller{contract: contract}, CurveLiquidityProviderTransactor: CurveLiquidityProviderTransactor{contract: contract}, CurveLiquidityProviderFilterer: CurveLiquidityProviderFilterer{contract: contract}}, nil
}

// NewCurveLiquidityProviderCaller creates a new read-only instance of CurveLiquidityProvider, bound to a specific deployed contract.
func NewCurveLiquidityProviderCaller(address common.Address, caller bind.ContractCaller) (*CurveLiquidityProviderCaller, error) {
	contract, err := bindCurveLiquidityProvider(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CurveLiquidityProviderCaller{contract: contract}, nil
}

// NewCurveLiquidityProviderTransactor creates a new write-only instance of CurveLiquidityProvider, bound to a specific deployed contract.
func NewCurveLiquidityProviderTransactor(address common.Address, transactor bind.ContractTransactor) (*CurveLiquidityProviderTransactor, error) {
	contract, err := bindCurveLiquidityProvider(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CurveLiquidityProviderTransactor{contract: contract}, nil
}

// NewCurveLiquidityProviderFilterer creates a new log filterer instance of CurveLiquidityProvider, bound to a specific deployed contract.
func NewCurveLiquidityProviderFilterer(address common.Address, filterer bind.ContractFilterer) (*CurveLiquidityProviderFilterer, error) {
	contract, err := bindCurveLiquidityProvider(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CurveLiquidityProviderFilterer{contract: contract}, nil
}

// bindCurveLiquidityProvider binds a generic wrapper to an already deployed contract.
func bindCurveLiquidityProvider(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CurveLiquidityProviderMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CurveLiquidityProvider *CurveLiquidityProviderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CurveLiquidityProvider.Contract.CurveLiquidityProviderCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CurveLiquidityProvider *CurveLiquidityProviderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CurveLiquidityProvider.Contract.CurveLiquidityProviderTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CurveLiquidityProvider *CurveLiquidityProviderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CurveLiquidityProvider.Contract.CurveLiquidityProviderTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CurveLiquidityProvider *CurveLiquidityProviderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CurveLiquidityProvider.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CurveLiquidityProvider *CurveLiquidityProviderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CurveLiquidityProvider.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CurveLiquidityProvider *CurveLiquidityProviderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CurveLiquidityProvider.Contract.contract.Transact(opts, method, params...)
}

// SeedLiquidity is a paid mutator transaction binding the contract method 0x2894eca4.
//
// Solidity: function seedLiquidity(address pool, address[3] coins, uint256 amount) returns()
func (_CurveLiquidityProvider *CurveLiquidityProviderTransactor) SeedLiquidity(opts *bind.TransactOpts, pool common.Address, coins [3]common.Address, amount *big.Int) (*types.Transaction, error) {
	return _CurveLiquidityProvider.contract.Transact(opts, "seedLiquidity", pool, coins, amount)
}

// SeedLiquidity is a paid mutator transaction binding the contract method 0x2894eca4.
//
// Solidity: function seedLiquidity(address pool, address[3] coins, uint256 amount) returns()
func (_CurveLiquidityProvider *CurveLiquidityProviderSession) SeedLiquidity(pool common.Address, coins [3]common.Address, amount *big.Int) (*types.Transaction, error) {
	return _CurveLiquidityProvider.Contract.SeedLiquidity(&_CurveLiquidityProvider.TransactOpts, pool, coins, amount)
}

// SeedLiquidity is a paid mutator transaction binding the contract method 0x2894eca4.
//
// Solidity: function seedLiquidity(address pool, address[3] coins, uint256 amount) returns()
func (_CurveLiquidityProvider *CurveLiquidityProviderTransactorSession) SeedLiquidity(pool common.Address, coins [3]common.Address, amount *big.Int) (*types.Transaction, error) {
	return _CurveLiquidityProvider.Contract.SeedLiquidity(&_CurveLiquidityProvider.TransactOpts, pool, coins, amount)
}
