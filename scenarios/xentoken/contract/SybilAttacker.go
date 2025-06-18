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

// XENSybilAttackerMetaData contains all meta data concerning the XENSybilAttacker contract.
var XENSybilAttackerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"xenContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"term\",\"type\":\"uint256\"}],\"name\":\"claimRank\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"predictProxyAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"xenContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"term\",\"type\":\"uint256\"}],\"name\":\"sybilAttack\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561000f575f80fd5b506105778061001d5f395ff3fe608060405234801561000f575f80fd5b506004361061003f575f3560e01c806351333b1914610043578063e2a4d3ad1461006c578063e8b1769f14610081575b5f80fd5b610056610051366004610397565b610094565b60405161006391906103b7565b60405180910390f35b61007f61007a3660046103e2565b610107565b005b61007f61008f366004610414565b6102d7565b5f608083901b8217816100a630610330565b8051602091820120604080516001600160f81b0319818501523060601b6001600160601b0319166021820152603581019590955260558086019290925280518086039092018252607590940190935282519201919091209150505b92915050565b620186a05f805b60148110156101bc575f610123600143610452565b4060c01c6001600160401b0316604083901b605a88901b17179050866001600160a01b031663df282331826001600160a01b03166040518263ffffffff1660e01b815260040161017391906103b7565b60c060405180830381865afa15801561018e573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906101b29190610475565b505060010161010e565b505b815a11156102d0575f6101d2600143610452565b4060c01c604083901b60a087901b171790505f6101ee30610330565b90505f828251602084015ff590506001600160a01b038116610212575050506102d0565b6040516001600160a01b038981166024830152604482018890525f919083169060640160408051601f198184030181529181526020820180516001600160e01b031663e8b1769f60e01b1790525161026a91906104fd565b5f604051808303815f865af19150503d805f81146102a3576040519150601f19603f3d011682016040523d82523d5f602084013e6102a8565b606091505b50509050806102ba57505050506102d0565b846102c481610529565b955050505050506101be565b5050505050565b604051639ff054df60e01b8152600481018290526001600160a01b03831690639ff054df906024015f604051808303815f87803b158015610316575f80fd5b505af1158015610328573d5f803e3d5ffd5b505050505050565b60405174600b380380600b5f395ff3363d3d373d3d3d363d7360581b6020820152606082811b6001600160601b03191660358301526e5af43d82803e903d91602b57fd5bf360881b6049830152906058016040516020818303038152906040529050919050565b5f80604083850312156103a8575f80fd5b50508035926020909101359150565b6001600160a01b0391909116815260200190565b6001600160a01b03811681146103df575f80fd5b50565b5f805f606084860312156103f4575f80fd5b83356103ff816103cb565b95602085013595506040909401359392505050565b5f8060408385031215610425575f80fd5b8235610430816103cb565b946020939093013593505050565b634e487b7160e01b5f52601160045260245ffd5b818103818111156101015761010161043e565b8051610470816103cb565b919050565b5f60c08284031215610485575f80fd5b60405160c081018181106001600160401b03821117156104b357634e487b7160e01b5f52604160045260245ffd5b6040526104bf83610465565b81526020830151602082015260408301516040820152606083015160608201526080830151608082015260a083015160a08201528091505092915050565b5f82515f5b8181101561051c5760208186018101518583015201610502565b505f920191825250919050565b5f6001820161053a5761053a61043e565b506001019056fea26469706673582212204da5bb1487ad2c137c4c10696e98081a765846fd4843b9e0c1cb9e23a2cea81164736f6c63430008160033",
}

// XENSybilAttackerABI is the input ABI used to generate the binding from.
// Deprecated: Use XENSybilAttackerMetaData.ABI instead.
var XENSybilAttackerABI = XENSybilAttackerMetaData.ABI

// XENSybilAttackerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use XENSybilAttackerMetaData.Bin instead.
var XENSybilAttackerBin = XENSybilAttackerMetaData.Bin

// DeployXENSybilAttacker deploys a new Ethereum contract, binding an instance of XENSybilAttacker to it.
func DeployXENSybilAttacker(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *XENSybilAttacker, error) {
	parsed, err := XENSybilAttackerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(XENSybilAttackerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &XENSybilAttacker{XENSybilAttackerCaller: XENSybilAttackerCaller{contract: contract}, XENSybilAttackerTransactor: XENSybilAttackerTransactor{contract: contract}, XENSybilAttackerFilterer: XENSybilAttackerFilterer{contract: contract}}, nil
}

// XENSybilAttacker is an auto generated Go binding around an Ethereum contract.
type XENSybilAttacker struct {
	XENSybilAttackerCaller     // Read-only binding to the contract
	XENSybilAttackerTransactor // Write-only binding to the contract
	XENSybilAttackerFilterer   // Log filterer for contract events
}

// XENSybilAttackerCaller is an auto generated read-only Go binding around an Ethereum contract.
type XENSybilAttackerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// XENSybilAttackerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type XENSybilAttackerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// XENSybilAttackerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type XENSybilAttackerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// XENSybilAttackerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type XENSybilAttackerSession struct {
	Contract     *XENSybilAttacker // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// XENSybilAttackerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type XENSybilAttackerCallerSession struct {
	Contract *XENSybilAttackerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// XENSybilAttackerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type XENSybilAttackerTransactorSession struct {
	Contract     *XENSybilAttackerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// XENSybilAttackerRaw is an auto generated low-level Go binding around an Ethereum contract.
type XENSybilAttackerRaw struct {
	Contract *XENSybilAttacker // Generic contract binding to access the raw methods on
}

// XENSybilAttackerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type XENSybilAttackerCallerRaw struct {
	Contract *XENSybilAttackerCaller // Generic read-only contract binding to access the raw methods on
}

// XENSybilAttackerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type XENSybilAttackerTransactorRaw struct {
	Contract *XENSybilAttackerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewXENSybilAttacker creates a new instance of XENSybilAttacker, bound to a specific deployed contract.
func NewXENSybilAttacker(address common.Address, backend bind.ContractBackend) (*XENSybilAttacker, error) {
	contract, err := bindXENSybilAttacker(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &XENSybilAttacker{XENSybilAttackerCaller: XENSybilAttackerCaller{contract: contract}, XENSybilAttackerTransactor: XENSybilAttackerTransactor{contract: contract}, XENSybilAttackerFilterer: XENSybilAttackerFilterer{contract: contract}}, nil
}

// NewXENSybilAttackerCaller creates a new read-only instance of XENSybilAttacker, bound to a specific deployed contract.
func NewXENSybilAttackerCaller(address common.Address, caller bind.ContractCaller) (*XENSybilAttackerCaller, error) {
	contract, err := bindXENSybilAttacker(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &XENSybilAttackerCaller{contract: contract}, nil
}

// NewXENSybilAttackerTransactor creates a new write-only instance of XENSybilAttacker, bound to a specific deployed contract.
func NewXENSybilAttackerTransactor(address common.Address, transactor bind.ContractTransactor) (*XENSybilAttackerTransactor, error) {
	contract, err := bindXENSybilAttacker(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &XENSybilAttackerTransactor{contract: contract}, nil
}

// NewXENSybilAttackerFilterer creates a new log filterer instance of XENSybilAttacker, bound to a specific deployed contract.
func NewXENSybilAttackerFilterer(address common.Address, filterer bind.ContractFilterer) (*XENSybilAttackerFilterer, error) {
	contract, err := bindXENSybilAttacker(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &XENSybilAttackerFilterer{contract: contract}, nil
}

// bindXENSybilAttacker binds a generic wrapper to an already deployed contract.
func bindXENSybilAttacker(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := XENSybilAttackerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_XENSybilAttacker *XENSybilAttackerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _XENSybilAttacker.Contract.XENSybilAttackerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_XENSybilAttacker *XENSybilAttackerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _XENSybilAttacker.Contract.XENSybilAttackerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_XENSybilAttacker *XENSybilAttackerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _XENSybilAttacker.Contract.XENSybilAttackerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_XENSybilAttacker *XENSybilAttackerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _XENSybilAttacker.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_XENSybilAttacker *XENSybilAttackerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _XENSybilAttacker.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_XENSybilAttacker *XENSybilAttackerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _XENSybilAttacker.Contract.contract.Transact(opts, method, params...)
}

// PredictProxyAddress is a free data retrieval call binding the contract method 0x51333b19.
//
// Solidity: function predictProxyAddress(uint256 seed, uint256 index) view returns(address)
func (_XENSybilAttacker *XENSybilAttackerCaller) PredictProxyAddress(opts *bind.CallOpts, seed *big.Int, index *big.Int) (common.Address, error) {
	var out []interface{}
	err := _XENSybilAttacker.contract.Call(opts, &out, "predictProxyAddress", seed, index)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PredictProxyAddress is a free data retrieval call binding the contract method 0x51333b19.
//
// Solidity: function predictProxyAddress(uint256 seed, uint256 index) view returns(address)
func (_XENSybilAttacker *XENSybilAttackerSession) PredictProxyAddress(seed *big.Int, index *big.Int) (common.Address, error) {
	return _XENSybilAttacker.Contract.PredictProxyAddress(&_XENSybilAttacker.CallOpts, seed, index)
}

// PredictProxyAddress is a free data retrieval call binding the contract method 0x51333b19.
//
// Solidity: function predictProxyAddress(uint256 seed, uint256 index) view returns(address)
func (_XENSybilAttacker *XENSybilAttackerCallerSession) PredictProxyAddress(seed *big.Int, index *big.Int) (common.Address, error) {
	return _XENSybilAttacker.Contract.PredictProxyAddress(&_XENSybilAttacker.CallOpts, seed, index)
}

// ClaimRank is a paid mutator transaction binding the contract method 0xe8b1769f.
//
// Solidity: function claimRank(address xenContract, uint256 term) returns()
func (_XENSybilAttacker *XENSybilAttackerTransactor) ClaimRank(opts *bind.TransactOpts, xenContract common.Address, term *big.Int) (*types.Transaction, error) {
	return _XENSybilAttacker.contract.Transact(opts, "claimRank", xenContract, term)
}

// ClaimRank is a paid mutator transaction binding the contract method 0xe8b1769f.
//
// Solidity: function claimRank(address xenContract, uint256 term) returns()
func (_XENSybilAttacker *XENSybilAttackerSession) ClaimRank(xenContract common.Address, term *big.Int) (*types.Transaction, error) {
	return _XENSybilAttacker.Contract.ClaimRank(&_XENSybilAttacker.TransactOpts, xenContract, term)
}

// ClaimRank is a paid mutator transaction binding the contract method 0xe8b1769f.
//
// Solidity: function claimRank(address xenContract, uint256 term) returns()
func (_XENSybilAttacker *XENSybilAttackerTransactorSession) ClaimRank(xenContract common.Address, term *big.Int) (*types.Transaction, error) {
	return _XENSybilAttacker.Contract.ClaimRank(&_XENSybilAttacker.TransactOpts, xenContract, term)
}

// SybilAttack is a paid mutator transaction binding the contract method 0xe2a4d3ad.
//
// Solidity: function sybilAttack(address xenContract, uint256 seed, uint256 term) returns()
func (_XENSybilAttacker *XENSybilAttackerTransactor) SybilAttack(opts *bind.TransactOpts, xenContract common.Address, seed *big.Int, term *big.Int) (*types.Transaction, error) {
	return _XENSybilAttacker.contract.Transact(opts, "sybilAttack", xenContract, seed, term)
}

// SybilAttack is a paid mutator transaction binding the contract method 0xe2a4d3ad.
//
// Solidity: function sybilAttack(address xenContract, uint256 seed, uint256 term) returns()
func (_XENSybilAttacker *XENSybilAttackerSession) SybilAttack(xenContract common.Address, seed *big.Int, term *big.Int) (*types.Transaction, error) {
	return _XENSybilAttacker.Contract.SybilAttack(&_XENSybilAttacker.TransactOpts, xenContract, seed, term)
}

// SybilAttack is a paid mutator transaction binding the contract method 0xe2a4d3ad.
//
// Solidity: function sybilAttack(address xenContract, uint256 seed, uint256 term) returns()
func (_XENSybilAttacker *XENSybilAttackerTransactorSession) SybilAttack(xenContract common.Address, seed *big.Int, term *big.Int) (*types.Transaction, error) {
	return _XENSybilAttacker.Contract.SybilAttack(&_XENSybilAttacker.TransactOpts, xenContract, seed, term)
}
