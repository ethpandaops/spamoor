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

// PairLiquidityProviderMetaData contains all meta data concerning the PairLiquidityProvider contract.
var PairLiquidityProviderMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"router1\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"router2\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenA\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenB\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountA\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountB\",\"type\":\"uint256\"}],\"name\":\"providePairLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506040516106a43803806106a483398101604081905261002f9161008d565b600080546001600160a01b039485166001600160a01b0319918216179091556001805493851693821693909317909255600280549190931691161790556100d0565b80516001600160a01b038116811461008857600080fd5b919050565b6000806000606084860312156100a257600080fd5b6100ab84610071565b92506100b960208501610071565b91506100c760408501610071565b90509250925092565b6105c5806100df6000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c80630d98c18b14610030575b600080fd5b61004361003e3660046104a9565b610045565b005b6000546001600160a01b031633146100905760405162461bcd60e51b81526020600482015260096024820152683737ba1037bbb732b960b91b60448201526064015b60405180910390fd5b6001600160a01b0380851660009081526003602090815260408083209387168352929052205460ff16156101065760405162461bcd60e51b815260206004820152601a60248201527f6c697175696469747920616c7265616479206465706c6f7965640000000000006044820152606401610087565b6001600160a01b0380851660009081526003602090815260408083209387168352929052908120805460ff191660011790556101436002846104eb565b905060006101526002846104eb565b90506001600160a01b0386166340c10f193061016f85600261050d565b6040516001600160e01b031960e085901b1681526001600160a01b0390921660048301526024820152604401600060405180830381600087803b1580156101b557600080fd5b505af11580156101c9573d6000803e3d6000fd5b50505050846001600160a01b03166340c10f19308360026101ea919061050d565b6040516001600160e01b031960e085901b1681526001600160a01b0390921660048301526024820152604401600060405180830381600087803b15801561023057600080fd5b505af1158015610244573d6000803e3d6000fd5b505060015461026292506001600160a01b0316905087878585610283565b60025461027b906001600160a01b031687878585610283565b505050505050565b60405163095ea7b360e01b81526001600160a01b0386811660048301526024820184905285169063095ea7b3906044016020604051808303816000875af11580156102d2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102f69190610538565b6103355760405162461bcd60e51b815260206004820152601060248201526f185c1c1c9bdd9948104819985a5b195960821b6044820152606401610087565b60405163095ea7b360e01b81526001600160a01b0386811660048301526024820183905284169063095ea7b3906044016020604051808303816000875af1158015610384573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103a89190610538565b6103e75760405162461bcd60e51b815260206004820152601060248201526f185c1c1c9bdd9948108819985a5b195960821b6044820152606401610087565b60405162e8e33760e81b81526001600160a01b0385811660048301528481166024830152604482018490526064820183905260006084830181905260a48301523060c48301524260e483015286169063e8e3370090610104016060604051808303816000875af115801561045f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104839190610561565b5050505050505050565b80356001600160a01b03811681146104a457600080fd5b919050565b600080600080608085870312156104bf57600080fd5b6104c88561048d565b93506104d66020860161048d565b93969395505050506040820135916060013590565b60008261050857634e487b7160e01b600052601260045260246000fd5b500490565b808202811582820484141761053257634e487b7160e01b600052601160045260246000fd5b92915050565b60006020828403121561054a57600080fd5b8151801515811461055a57600080fd5b9392505050565b60008060006060848603121561057657600080fd5b835192506020840151915060408401519050925092509256fea2646970667358221220c6593ef838bbcceb101eae738f5aad6f487097f88fc6dda5652a2b93a5429d5764736f6c63430008110033",
}

// PairLiquidityProviderABI is the input ABI used to generate the binding from.
// Deprecated: Use PairLiquidityProviderMetaData.ABI instead.
var PairLiquidityProviderABI = PairLiquidityProviderMetaData.ABI

// PairLiquidityProviderBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use PairLiquidityProviderMetaData.Bin instead.
var PairLiquidityProviderBin = PairLiquidityProviderMetaData.Bin

// DeployPairLiquidityProvider deploys a new Ethereum contract, binding an instance of PairLiquidityProvider to it.
func DeployPairLiquidityProvider(auth *bind.TransactOpts, backend bind.ContractBackend, owner common.Address, router1 common.Address, router2 common.Address) (common.Address, *types.Transaction, *PairLiquidityProvider, error) {
	parsed, err := PairLiquidityProviderMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PairLiquidityProviderBin), backend, owner, router1, router2)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &PairLiquidityProvider{PairLiquidityProviderCaller: PairLiquidityProviderCaller{contract: contract}, PairLiquidityProviderTransactor: PairLiquidityProviderTransactor{contract: contract}, PairLiquidityProviderFilterer: PairLiquidityProviderFilterer{contract: contract}}, nil
}

// PairLiquidityProvider is an auto generated Go binding around an Ethereum contract.
type PairLiquidityProvider struct {
	PairLiquidityProviderCaller     // Read-only binding to the contract
	PairLiquidityProviderTransactor // Write-only binding to the contract
	PairLiquidityProviderFilterer   // Log filterer for contract events
}

// PairLiquidityProviderCaller is an auto generated read-only Go binding around an Ethereum contract.
type PairLiquidityProviderCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PairLiquidityProviderTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PairLiquidityProviderTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PairLiquidityProviderFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PairLiquidityProviderFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PairLiquidityProviderSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PairLiquidityProviderSession struct {
	Contract     *PairLiquidityProvider // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// PairLiquidityProviderCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PairLiquidityProviderCallerSession struct {
	Contract *PairLiquidityProviderCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// PairLiquidityProviderTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PairLiquidityProviderTransactorSession struct {
	Contract     *PairLiquidityProviderTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// PairLiquidityProviderRaw is an auto generated low-level Go binding around an Ethereum contract.
type PairLiquidityProviderRaw struct {
	Contract *PairLiquidityProvider // Generic contract binding to access the raw methods on
}

// PairLiquidityProviderCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PairLiquidityProviderCallerRaw struct {
	Contract *PairLiquidityProviderCaller // Generic read-only contract binding to access the raw methods on
}

// PairLiquidityProviderTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PairLiquidityProviderTransactorRaw struct {
	Contract *PairLiquidityProviderTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPairLiquidityProvider creates a new instance of PairLiquidityProvider, bound to a specific deployed contract.
func NewPairLiquidityProvider(address common.Address, backend bind.ContractBackend) (*PairLiquidityProvider, error) {
	contract, err := bindPairLiquidityProvider(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PairLiquidityProvider{PairLiquidityProviderCaller: PairLiquidityProviderCaller{contract: contract}, PairLiquidityProviderTransactor: PairLiquidityProviderTransactor{contract: contract}, PairLiquidityProviderFilterer: PairLiquidityProviderFilterer{contract: contract}}, nil
}

// NewPairLiquidityProviderCaller creates a new read-only instance of PairLiquidityProvider, bound to a specific deployed contract.
func NewPairLiquidityProviderCaller(address common.Address, caller bind.ContractCaller) (*PairLiquidityProviderCaller, error) {
	contract, err := bindPairLiquidityProvider(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PairLiquidityProviderCaller{contract: contract}, nil
}

// NewPairLiquidityProviderTransactor creates a new write-only instance of PairLiquidityProvider, bound to a specific deployed contract.
func NewPairLiquidityProviderTransactor(address common.Address, transactor bind.ContractTransactor) (*PairLiquidityProviderTransactor, error) {
	contract, err := bindPairLiquidityProvider(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PairLiquidityProviderTransactor{contract: contract}, nil
}

// NewPairLiquidityProviderFilterer creates a new log filterer instance of PairLiquidityProvider, bound to a specific deployed contract.
func NewPairLiquidityProviderFilterer(address common.Address, filterer bind.ContractFilterer) (*PairLiquidityProviderFilterer, error) {
	contract, err := bindPairLiquidityProvider(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PairLiquidityProviderFilterer{contract: contract}, nil
}

// bindPairLiquidityProvider binds a generic wrapper to an already deployed contract.
func bindPairLiquidityProvider(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PairLiquidityProviderMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PairLiquidityProvider *PairLiquidityProviderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PairLiquidityProvider.Contract.PairLiquidityProviderCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PairLiquidityProvider *PairLiquidityProviderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PairLiquidityProvider.Contract.PairLiquidityProviderTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PairLiquidityProvider *PairLiquidityProviderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PairLiquidityProvider.Contract.PairLiquidityProviderTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PairLiquidityProvider *PairLiquidityProviderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PairLiquidityProvider.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PairLiquidityProvider *PairLiquidityProviderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PairLiquidityProvider.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PairLiquidityProvider *PairLiquidityProviderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PairLiquidityProvider.Contract.contract.Transact(opts, method, params...)
}

// ProvidePairLiquidity is a paid mutator transaction binding the contract method 0x0d98c18b.
//
// Solidity: function providePairLiquidity(address tokenA, address tokenB, uint256 amountA, uint256 amountB) returns()
func (_PairLiquidityProvider *PairLiquidityProviderTransactor) ProvidePairLiquidity(opts *bind.TransactOpts, tokenA common.Address, tokenB common.Address, amountA *big.Int, amountB *big.Int) (*types.Transaction, error) {
	return _PairLiquidityProvider.contract.Transact(opts, "providePairLiquidity", tokenA, tokenB, amountA, amountB)
}

// ProvidePairLiquidity is a paid mutator transaction binding the contract method 0x0d98c18b.
//
// Solidity: function providePairLiquidity(address tokenA, address tokenB, uint256 amountA, uint256 amountB) returns()
func (_PairLiquidityProvider *PairLiquidityProviderSession) ProvidePairLiquidity(tokenA common.Address, tokenB common.Address, amountA *big.Int, amountB *big.Int) (*types.Transaction, error) {
	return _PairLiquidityProvider.Contract.ProvidePairLiquidity(&_PairLiquidityProvider.TransactOpts, tokenA, tokenB, amountA, amountB)
}

// ProvidePairLiquidity is a paid mutator transaction binding the contract method 0x0d98c18b.
//
// Solidity: function providePairLiquidity(address tokenA, address tokenB, uint256 amountA, uint256 amountB) returns()
func (_PairLiquidityProvider *PairLiquidityProviderTransactorSession) ProvidePairLiquidity(tokenA common.Address, tokenB common.Address, amountA *big.Int, amountB *big.Int) (*types.Transaction, error) {
	return _PairLiquidityProvider.Contract.ProvidePairLiquidity(&_PairLiquidityProvider.TransactOpts, tokenA, tokenB, amountA, amountB)
}
