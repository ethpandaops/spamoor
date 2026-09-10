// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"context"
	"errors"
	"math/big"
	"strings"
	"time"

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
	_ = time.Tick
	_ = context.Background
)

// V3LiquidityProviderMetaData contains all meta data concerning the V3LiquidityProvider contract.
var V3LiquidityProviderMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"},{\"internalType\":\"int24\",\"name\":\"tickLower\",\"type\":\"int24\"},{\"internalType\":\"int24\",\"name\":\"tickUpper\",\"type\":\"int24\"},{\"internalType\":\"uint128\",\"name\":\"liquidity\",\"type\":\"uint128\"}],\"name\":\"provideLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount0Owed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount1Owed\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"uniswapV3MintCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060405161057f38038061057f83398101604081905261002f91610054565b600080546001600160a01b0319166001600160a01b0392909216919091179055610084565b60006020828403121561006657600080fd5b81516001600160a01b038116811461007d57600080fd5b9392505050565b6104ec806100936000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80632b2785081461003b578063d348799714610050575b600080fd5b61004e61004936600461038a565b610063565b005b61004e61005e3660046103ee565b610171565b6000546001600160a01b031633146100ae5760405162461bcd60e51b81526020600482015260096024820152683737ba1037bbb732b960b91b60448201526064015b60405180910390fd5b600180546001600160a01b0319166001600160a01b038616908117909155604051633c8a7d8d60e01b8152306004820152600285810b602483015284900b60448201526001600160801b038316606482015260a06084820152600060a4820152633c8a7d8d9060c40160408051808303816000875af1158015610135573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610159919061046e565b5050600180546001600160a01b031916905550505050565b6001546001600160a01b031633146101bf5760405162461bcd60e51b81526020600482015260116024820152703ab732bc3832b1ba32b21031b0b63632b960791b60448201526064016100a5565b831561028a57336001600160a01b0316630dfe16816040518163ffffffff1660e01b8152600401602060405180830381865afa158015610203573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102279190610492565b6040516340c10f1960e01b8152336004820152602481018690526001600160a01b0391909116906340c10f1990604401600060405180830381600087803b15801561027157600080fd5b505af1158015610285573d6000803e3d6000fd5b505050505b821561035557336001600160a01b031663d21220a76040518163ffffffff1660e01b8152600401602060405180830381865afa1580156102ce573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102f29190610492565b6040516340c10f1960e01b8152336004820152602481018590526001600160a01b0391909116906340c10f1990604401600060405180830381600087803b15801561033c57600080fd5b505af1158015610350573d6000803e3d6000fd5b505050505b50505050565b6001600160a01b038116811461037057600080fd5b50565b8035600281900b811461038557600080fd5b919050565b600080600080608085870312156103a057600080fd5b84356103ab8161035b565b93506103b960208601610373565b92506103c760408601610373565b915060608501356001600160801b03811681146103e357600080fd5b939692955090935050565b6000806000806060858703121561040457600080fd5b8435935060208501359250604085013567ffffffffffffffff8082111561042a57600080fd5b818701915087601f83011261043e57600080fd5b81358181111561044d57600080fd5b88602082850101111561045f57600080fd5b95989497505060200194505050565b6000806040838503121561048157600080fd5b505080516020909101519092909150565b6000602082840312156104a457600080fd5b81516104af8161035b565b939250505056fea26469706673582212208c4b3f6b983a1b5ab98728fca626459dcbfa23a216091deaf779d11d17ef132a64736f6c63430008110033",
}

// V3LiquidityProviderABI is the input ABI used to generate the binding from.
// Deprecated: Use V3LiquidityProviderMetaData.ABI instead.
var V3LiquidityProviderABI = V3LiquidityProviderMetaData.ABI

// V3LiquidityProviderBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use V3LiquidityProviderMetaData.Bin instead.
var V3LiquidityProviderBin = V3LiquidityProviderMetaData.Bin

// DeployV3LiquidityProvider deploys a new Ethereum contract, binding an instance of V3LiquidityProvider to it.
func DeployV3LiquidityProvider(auth *bind.TransactOpts, backend bind.ContractBackend, owner common.Address) (common.Address, *types.Transaction, *V3LiquidityProvider, error) {
	parsed, err := V3LiquidityProviderMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(V3LiquidityProviderBin), backend, owner)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &V3LiquidityProvider{V3LiquidityProviderCaller: V3LiquidityProviderCaller{contract: contract}, V3LiquidityProviderTransactor: V3LiquidityProviderTransactor{contract: contract}, V3LiquidityProviderFilterer: V3LiquidityProviderFilterer{contract: contract}}, nil
}

// V3LiquidityProvider is an auto generated Go binding around an Ethereum contract.
type V3LiquidityProvider struct {
	V3LiquidityProviderCaller     // Read-only binding to the contract
	V3LiquidityProviderTransactor // Write-only binding to the contract
	V3LiquidityProviderFilterer   // Log filterer for contract events
}

// V3LiquidityProviderCaller is an auto generated read-only Go binding around an Ethereum contract.
type V3LiquidityProviderCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// V3LiquidityProviderTransactor is an auto generated write-only Go binding around an Ethereum contract.
type V3LiquidityProviderTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// V3LiquidityProviderFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type V3LiquidityProviderFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// V3LiquidityProviderSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type V3LiquidityProviderSession struct {
	Contract     *V3LiquidityProvider // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// V3LiquidityProviderCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type V3LiquidityProviderCallerSession struct {
	Contract *V3LiquidityProviderCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// V3LiquidityProviderTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type V3LiquidityProviderTransactorSession struct {
	Contract     *V3LiquidityProviderTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// V3LiquidityProviderRaw is an auto generated low-level Go binding around an Ethereum contract.
type V3LiquidityProviderRaw struct {
	Contract *V3LiquidityProvider // Generic contract binding to access the raw methods on
}

// V3LiquidityProviderCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type V3LiquidityProviderCallerRaw struct {
	Contract *V3LiquidityProviderCaller // Generic read-only contract binding to access the raw methods on
}

// V3LiquidityProviderTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type V3LiquidityProviderTransactorRaw struct {
	Contract *V3LiquidityProviderTransactor // Generic write-only contract binding to access the raw methods on
}

// NewV3LiquidityProvider creates a new instance of V3LiquidityProvider, bound to a specific deployed contract.
func NewV3LiquidityProvider(address common.Address, backend bind.ContractBackend) (*V3LiquidityProvider, error) {
	contract, err := bindV3LiquidityProvider(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &V3LiquidityProvider{V3LiquidityProviderCaller: V3LiquidityProviderCaller{contract: contract}, V3LiquidityProviderTransactor: V3LiquidityProviderTransactor{contract: contract}, V3LiquidityProviderFilterer: V3LiquidityProviderFilterer{contract: contract}}, nil
}

// NewV3LiquidityProviderCaller creates a new read-only instance of V3LiquidityProvider, bound to a specific deployed contract.
func NewV3LiquidityProviderCaller(address common.Address, caller bind.ContractCaller) (*V3LiquidityProviderCaller, error) {
	contract, err := bindV3LiquidityProvider(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &V3LiquidityProviderCaller{contract: contract}, nil
}

// NewV3LiquidityProviderTransactor creates a new write-only instance of V3LiquidityProvider, bound to a specific deployed contract.
func NewV3LiquidityProviderTransactor(address common.Address, transactor bind.ContractTransactor) (*V3LiquidityProviderTransactor, error) {
	contract, err := bindV3LiquidityProvider(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &V3LiquidityProviderTransactor{contract: contract}, nil
}

// NewV3LiquidityProviderFilterer creates a new log filterer instance of V3LiquidityProvider, bound to a specific deployed contract.
func NewV3LiquidityProviderFilterer(address common.Address, filterer bind.ContractFilterer) (*V3LiquidityProviderFilterer, error) {
	contract, err := bindV3LiquidityProvider(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &V3LiquidityProviderFilterer{contract: contract}, nil
}

// bindV3LiquidityProvider binds a generic wrapper to an already deployed contract.
func bindV3LiquidityProvider(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := V3LiquidityProviderMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_V3LiquidityProvider *V3LiquidityProviderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _V3LiquidityProvider.Contract.V3LiquidityProviderCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_V3LiquidityProvider *V3LiquidityProviderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _V3LiquidityProvider.Contract.V3LiquidityProviderTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_V3LiquidityProvider *V3LiquidityProviderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _V3LiquidityProvider.Contract.V3LiquidityProviderTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_V3LiquidityProvider *V3LiquidityProviderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _V3LiquidityProvider.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_V3LiquidityProvider *V3LiquidityProviderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _V3LiquidityProvider.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_V3LiquidityProvider *V3LiquidityProviderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _V3LiquidityProvider.Contract.contract.Transact(opts, method, params...)
}

// ProvideLiquidity is a paid mutator transaction binding the contract method 0x2b278508.
//
// Solidity: function provideLiquidity(address pool, int24 tickLower, int24 tickUpper, uint128 liquidity) returns()
func (_V3LiquidityProvider *V3LiquidityProviderTransactor) ProvideLiquidity(opts *bind.TransactOpts, pool common.Address, tickLower *big.Int, tickUpper *big.Int, liquidity *big.Int) (*types.Transaction, error) {
	return _V3LiquidityProvider.contract.Transact(opts, "provideLiquidity", pool, tickLower, tickUpper, liquidity)
}

// ProvideLiquidity is a paid mutator transaction binding the contract method 0x2b278508.
//
// Solidity: function provideLiquidity(address pool, int24 tickLower, int24 tickUpper, uint128 liquidity) returns()
func (_V3LiquidityProvider *V3LiquidityProviderSession) ProvideLiquidity(pool common.Address, tickLower *big.Int, tickUpper *big.Int, liquidity *big.Int) (*types.Transaction, error) {
	return _V3LiquidityProvider.Contract.ProvideLiquidity(&_V3LiquidityProvider.TransactOpts, pool, tickLower, tickUpper, liquidity)
}

// ProvideLiquidity is a paid mutator transaction binding the contract method 0x2b278508.
//
// Solidity: function provideLiquidity(address pool, int24 tickLower, int24 tickUpper, uint128 liquidity) returns()
func (_V3LiquidityProvider *V3LiquidityProviderTransactorSession) ProvideLiquidity(pool common.Address, tickLower *big.Int, tickUpper *big.Int, liquidity *big.Int) (*types.Transaction, error) {
	return _V3LiquidityProvider.Contract.ProvideLiquidity(&_V3LiquidityProvider.TransactOpts, pool, tickLower, tickUpper, liquidity)
}

// UniswapV3MintCallback is a paid mutator transaction binding the contract method 0xd3487997.
//
// Solidity: function uniswapV3MintCallback(uint256 amount0Owed, uint256 amount1Owed, bytes ) returns()
func (_V3LiquidityProvider *V3LiquidityProviderTransactor) UniswapV3MintCallback(opts *bind.TransactOpts, amount0Owed *big.Int, amount1Owed *big.Int, arg2 []byte) (*types.Transaction, error) {
	return _V3LiquidityProvider.contract.Transact(opts, "uniswapV3MintCallback", amount0Owed, amount1Owed, arg2)
}

// UniswapV3MintCallback is a paid mutator transaction binding the contract method 0xd3487997.
//
// Solidity: function uniswapV3MintCallback(uint256 amount0Owed, uint256 amount1Owed, bytes ) returns()
func (_V3LiquidityProvider *V3LiquidityProviderSession) UniswapV3MintCallback(amount0Owed *big.Int, amount1Owed *big.Int, arg2 []byte) (*types.Transaction, error) {
	return _V3LiquidityProvider.Contract.UniswapV3MintCallback(&_V3LiquidityProvider.TransactOpts, amount0Owed, amount1Owed, arg2)
}

// UniswapV3MintCallback is a paid mutator transaction binding the contract method 0xd3487997.
//
// Solidity: function uniswapV3MintCallback(uint256 amount0Owed, uint256 amount1Owed, bytes ) returns()
func (_V3LiquidityProvider *V3LiquidityProviderTransactorSession) UniswapV3MintCallback(amount0Owed *big.Int, amount1Owed *big.Int, arg2 []byte) (*types.Transaction, error) {
	return _V3LiquidityProvider.Contract.UniswapV3MintCallback(&_V3LiquidityProvider.TransactOpts, amount0Owed, amount1Owed, arg2)
}
