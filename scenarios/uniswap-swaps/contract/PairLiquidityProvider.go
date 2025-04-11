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

// PairLiquidityProviderMetaData contains all meta data concerning the PairLiquidityProvider contract.
var PairLiquidityProviderMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner1\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner2\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"router1\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"router2\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"weth9\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"call\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"onERC721Received\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"dai\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"daiDesired\",\"type\":\"uint256\"}],\"name\":\"providePairLiquidity\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x608060405234801561001057600080fd5b506040516109c43803806109c483398101604081905261002f916100ad565b600080546001600160a01b03199081166001600160a01b03978816179091556001805482169587169590951790945560028054851693861693909317909255600380548416918516919091179055600480549092169216919091179055610112565b80516001600160a01b03811681146100a857600080fd5b919050565b600080600080600060a086880312156100c557600080fd5b6100ce86610091565b94506100dc60208701610091565b93506100ea60408701610091565b92506100f860608701610091565b915061010660808701610091565b90509295509295909350565b6108a3806101216000396000f3fe6080604052600436106100385760003560e01c8063150b7a02146100445780631b8b921d1461008d5780633f0e37c3146100a257600080fd5b3661003f57005b600080fd5b34801561005057600080fd5b5061007061005f366004610663565b630a85bd0160e11b95945050505050565b6040516001600160e01b0319909116815260200160405180910390f35b6100a061009b3660046106d2565b6100b5565b005b6100a06100b0366004610725565b6101b9565b6000546001600160a01b03163314806100d857506001546001600160a01b031633145b6101155760405162461bcd60e51b81526020600482015260096024820152683737ba1037bbb732b960b91b60448201526064015b60405180910390fd5b6000836001600160a01b031634848460405161013292919061074f565b60006040518083038185875af1925050503d806000811461016f576040519150601f19603f3d011682016040523d82523d6000602084013e610174565b606091505b50509050806101b35760405162461bcd60e51b815260206004820152600b60248201526a18d85b1b0819985a5b195960aa1b604482015260640161010c565b50505050565b6000546001600160a01b03163314806101dc57506001546001600160a01b031633145b6102145760405162461bcd60e51b81526020600482015260096024820152683737ba1037bbb732b960b91b604482015260640161010c565b6001600160a01b03821660009081526005602052604090205460ff161561027d5760405162461bcd60e51b815260206004820152601a60248201527f6c697175696469747920616c7265616479206465706c6f796564000000000000604482015260640161010c565b6001600160a01b0382166000908152600560205260408120805460ff191660011790556102ab60028361075f565b90506102b8816002610781565b6040516340c10f1960e01b8152306004820152602481018290529092506001600160a01b038416906340c10f1990604401600060405180830381600087803b15801561030357600080fd5b505af1158015610317573d6000803e3d6000fd5b505060025461033392508591506001600160a01b0316836104fe565b60035461034b9084906001600160a01b0316836104fe565b600280546001600160a01b03169063f305d71990610369903461075f565b858460008030426040518863ffffffff1660e01b8152600401610391969594939291906107ac565b60606040518083038185885af11580156103af573d6000803e3d6000fd5b50505050506040513d601f19601f820116820180604052508101906103d491906107e7565b50506003546001600160a01b0316905063f305d7196103f460023461075f565b858460008030426040518863ffffffff1660e01b815260040161041c969594939291906107ac565b60606040518083038185885af115801561043a573d6000803e3d6000fd5b50505050506040513d601f19601f8201168201806040525081019061045f91906107e7565b5050471590506104f957604051600090329047908381818185875af1925050503d80600081146104ab576040519150601f19603f3d011682016040523d82523d6000602084013e6104b0565b606091505b50509050806101b35760405162461bcd60e51b815260206004820152601560248201527419985a5b1959081d1bc81cd95b99081c99599d5b99605a1b604482015260640161010c565b505050565b604080516001600160a01b038481166024830152604480830185905283518084039091018152606490920183526020820180516001600160e01b031663095ea7b360e01b179052915160009283929087169161055a9190610815565b6000604051808303816000865af19150503d8060008114610597576040519150601f19603f3d011682016040523d82523d6000602084013e61059c565b606091505b50915091508180156105c65750805115806105c65750808060200190518101906105c69190610844565b6105f75760405162461bcd60e51b8152602060048201526002602482015261534160f01b604482015260640161010c565b5050505050565b80356001600160a01b038116811461061557600080fd5b919050565b60008083601f84011261062c57600080fd5b50813567ffffffffffffffff81111561064457600080fd5b60208301915083602082850101111561065c57600080fd5b9250929050565b60008060008060006080868803121561067b57600080fd5b610684866105fe565b9450610692602087016105fe565b935060408601359250606086013567ffffffffffffffff8111156106b557600080fd5b6106c18882890161061a565b969995985093965092949392505050565b6000806000604084860312156106e757600080fd5b6106f0846105fe565b9250602084013567ffffffffffffffff81111561070c57600080fd5b6107188682870161061a565b9497909650939450505050565b6000806040838503121561073857600080fd5b610741836105fe565b946020939093013593505050565b8183823760009101908152919050565b60008261077c57634e487b7160e01b600052601260045260246000fd5b500490565b80820281158282048414176107a657634e487b7160e01b600052601160045260246000fd5b92915050565b6001600160a01b039687168152602081019590955260408501939093526060840191909152909216608082015260a081019190915260c00190565b6000806000606084860312156107fc57600080fd5b8351925060208401519150604084015190509250925092565b6000825160005b81811015610836576020818601810151858301520161081c565b506000920191825250919050565b60006020828403121561085657600080fd5b8151801515811461086657600080fd5b939250505056fea264697066735822122023a6ff2127f25edf381e1a5ce2fe98d438a91ce501aafe376227c0e576b603af64736f6c63430008110033",
}

// PairLiquidityProviderABI is the input ABI used to generate the binding from.
// Deprecated: Use PairLiquidityProviderMetaData.ABI instead.
var PairLiquidityProviderABI = PairLiquidityProviderMetaData.ABI

// PairLiquidityProviderBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use PairLiquidityProviderMetaData.Bin instead.
var PairLiquidityProviderBin = PairLiquidityProviderMetaData.Bin

// DeployPairLiquidityProvider deploys a new Ethereum contract, binding an instance of PairLiquidityProvider to it.
func DeployPairLiquidityProvider(auth *bind.TransactOpts, backend bind.ContractBackend, owner1 common.Address, owner2 common.Address, router1 common.Address, router2 common.Address, weth9 common.Address) (common.Address, *types.Transaction, *PairLiquidityProvider, error) {
	parsed, err := PairLiquidityProviderMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PairLiquidityProviderBin), backend, owner1, owner2, router1, router2, weth9)
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

// OnERC721Received is a free data retrieval call binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address , uint256 , bytes ) pure returns(bytes4)
func (_PairLiquidityProvider *PairLiquidityProviderCaller) OnERC721Received(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) ([4]byte, error) {
	var out []interface{}
	err := _PairLiquidityProvider.contract.Call(opts, &out, "onERC721Received", arg0, arg1, arg2, arg3)

	if err != nil {
		return *new([4]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([4]byte)).(*[4]byte)

	return out0, err

}

// OnERC721Received is a free data retrieval call binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address , uint256 , bytes ) pure returns(bytes4)
func (_PairLiquidityProvider *PairLiquidityProviderSession) OnERC721Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) ([4]byte, error) {
	return _PairLiquidityProvider.Contract.OnERC721Received(&_PairLiquidityProvider.CallOpts, arg0, arg1, arg2, arg3)
}

// OnERC721Received is a free data retrieval call binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address , uint256 , bytes ) pure returns(bytes4)
func (_PairLiquidityProvider *PairLiquidityProviderCallerSession) OnERC721Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) ([4]byte, error) {
	return _PairLiquidityProvider.Contract.OnERC721Received(&_PairLiquidityProvider.CallOpts, arg0, arg1, arg2, arg3)
}

// Call is a paid mutator transaction binding the contract method 0x1b8b921d.
//
// Solidity: function call(address addr, bytes data) payable returns()
func (_PairLiquidityProvider *PairLiquidityProviderTransactor) Call(opts *bind.TransactOpts, addr common.Address, data []byte) (*types.Transaction, error) {
	return _PairLiquidityProvider.contract.Transact(opts, "call", addr, data)
}

// Call is a paid mutator transaction binding the contract method 0x1b8b921d.
//
// Solidity: function call(address addr, bytes data) payable returns()
func (_PairLiquidityProvider *PairLiquidityProviderSession) Call(addr common.Address, data []byte) (*types.Transaction, error) {
	return _PairLiquidityProvider.Contract.Call(&_PairLiquidityProvider.TransactOpts, addr, data)
}

// Call is a paid mutator transaction binding the contract method 0x1b8b921d.
//
// Solidity: function call(address addr, bytes data) payable returns()
func (_PairLiquidityProvider *PairLiquidityProviderTransactorSession) Call(addr common.Address, data []byte) (*types.Transaction, error) {
	return _PairLiquidityProvider.Contract.Call(&_PairLiquidityProvider.TransactOpts, addr, data)
}

// ProvidePairLiquidity is a paid mutator transaction binding the contract method 0x3f0e37c3.
//
// Solidity: function providePairLiquidity(address dai, uint256 daiDesired) payable returns()
func (_PairLiquidityProvider *PairLiquidityProviderTransactor) ProvidePairLiquidity(opts *bind.TransactOpts, dai common.Address, daiDesired *big.Int) (*types.Transaction, error) {
	return _PairLiquidityProvider.contract.Transact(opts, "providePairLiquidity", dai, daiDesired)
}

// ProvidePairLiquidity is a paid mutator transaction binding the contract method 0x3f0e37c3.
//
// Solidity: function providePairLiquidity(address dai, uint256 daiDesired) payable returns()
func (_PairLiquidityProvider *PairLiquidityProviderSession) ProvidePairLiquidity(dai common.Address, daiDesired *big.Int) (*types.Transaction, error) {
	return _PairLiquidityProvider.Contract.ProvidePairLiquidity(&_PairLiquidityProvider.TransactOpts, dai, daiDesired)
}

// ProvidePairLiquidity is a paid mutator transaction binding the contract method 0x3f0e37c3.
//
// Solidity: function providePairLiquidity(address dai, uint256 daiDesired) payable returns()
func (_PairLiquidityProvider *PairLiquidityProviderTransactorSession) ProvidePairLiquidity(dai common.Address, daiDesired *big.Int) (*types.Transaction, error) {
	return _PairLiquidityProvider.Contract.ProvidePairLiquidity(&_PairLiquidityProvider.TransactOpts, dai, daiDesired)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_PairLiquidityProvider *PairLiquidityProviderTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PairLiquidityProvider.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_PairLiquidityProvider *PairLiquidityProviderSession) Receive() (*types.Transaction, error) {
	return _PairLiquidityProvider.Contract.Receive(&_PairLiquidityProvider.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_PairLiquidityProvider *PairLiquidityProviderTransactorSession) Receive() (*types.Transaction, error) {
	return _PairLiquidityProvider.Contract.Receive(&_PairLiquidityProvider.TransactOpts)
}
