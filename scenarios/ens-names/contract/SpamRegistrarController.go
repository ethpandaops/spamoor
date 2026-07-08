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

// SpamRegistrarControllerMetaData contains all meta data concerning the SpamRegistrarController contract.
var SpamRegistrarControllerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIBaseRegistrar\",\"name\":\"_base\",\"type\":\"address\"},{\"internalType\":\"contractIENS\",\"name\":\"_ens\",\"type\":\"address\"},{\"internalType\":\"contractIReverseRegistrar\",\"name\":\"_reverseRegistrar\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_resolver\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ETH_NODE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"base\",\"outputs\":[{\"internalType\":\"contractIBaseRegistrar\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ens\",\"outputs\":[{\"internalType\":\"contractIENS\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"label\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"duration\",\"type\":\"uint256\"}],\"name\":\"register\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"expires\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"label\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"duration\",\"type\":\"uint256\"}],\"name\":\"registerNamed\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"label\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"duration\",\"type\":\"uint256\"}],\"name\":\"renew\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"expires\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"resolver\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reverseRegistrar\",\"outputs\":[{\"internalType\":\"contractIReverseRegistrar\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x610100604052348015610010575f80fd5b5060405161099b38038061099b83398101604081905261002f91610068565b6001600160a01b0393841660805291831660a052821660c0521660e0526100c4565b6001600160a01b0381168114610065575f80fd5b50565b5f805f806080858703121561007b575f80fd5b845161008681610051565b602086015190945061009781610051565b60408601519093506100a881610051565b60608601519092506100b981610051565b939692955090935050565b60805160a05160c05160e0516108686101335f395f8181608e015281816102c70152818161036c01526103f801525f818161012001526103c701525f818160d2015281816102f201526104bf01525f818160f9015281816101ed0152818161052201526105d301526108685ff3fe608060405234801561000f575f80fd5b5060043610610085575f3560e01c80638ba778f0116100585780638ba778f014610142578063acf1a84114610163578063cc473be314610176578063d393c8711461019d575f80fd5b806304f3bcec146100895780633f15457f146100cd5780635001f3b5146100f4578063808698531461011b575b5f80fd5b6100b07f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b0390911681526020015b60405180910390f35b6100b07f000000000000000000000000000000000000000000000000000000000000000081565b6100b07f000000000000000000000000000000000000000000000000000000000000000081565b6100b07f000000000000000000000000000000000000000000000000000000000000000081565b6101556101503660046106d6565b6101b0565b6040519081526020016100c4565b61015561017136600461073a565b61051f565b6101557f93cdeb708b7545dc668eb9280176169d1c33cfd8ed6f04690a0bcc88a93fc4ae81565b6101556101ab3660046106d6565b6105d0565b5f8085856040516101c2929190610782565b604051908190038120633f2891eb60e21b8252600482018190523060248301526044820185905291507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063fca247ac906064016020604051808303815f875af115801561023b573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061025f9190610791565b50604080517f93cdeb708b7545dc668eb9280176169d1c33cfd8ed6f04690a0bcc88a93fc4ae602082015290810182905260600160408051601f19818403018152908290528051602090910120630c4b7b8560e11b8252600482018190526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000811660248401529093507f00000000000000000000000000000000000000000000000000000000000000001690631896f70a906044015f604051808303815f87803b158015610333575f80fd5b505af1158015610345573d5f803e3d5ffd5b505060405162d5fa2b60e81b8152600481018590526001600160a01b0387811660248301527f000000000000000000000000000000000000000000000000000000000000000016925063d5fa2b0091506044015f604051808303815f87803b1580156103af575f80fd5b505af11580156103c1573d5f803e3d5ffd5b505050507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316637a806d6b85867f00000000000000000000000000000000000000000000000000000000000000008a8a60405160200161042a9291906107a8565b6040516020818303038152906040526040518563ffffffff1660e01b815260040161045894939291906107c1565b6020604051808303815f875af1158015610474573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906104989190610791565b50604051635b0fc9c360e01b8152600481018390526001600160a01b0385811660248301527f00000000000000000000000000000000000000000000000000000000000000001690635b0fc9c3906044015f604051808303815f87803b158015610500575f80fd5b505af1158015610512573d5f803e3d5ffd5b5050505050949350505050565b5f7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663c475abff858560405161055f929190610782565b60405190819003812060e083901b6001600160e01b03191682526004820152602481018590526044016020604051808303815f875af11580156105a4573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906105c89190610791565b949350505050565b5f7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663fca247ac8686604051610610929190610782565b60405190819003812060e083901b6001600160e01b031916825260048201526001600160a01b0386166024820152604481018590526064016020604051808303815f875af1158015610664573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906106889190610791565b95945050505050565b5f8083601f8401126106a1575f80fd5b50813567ffffffffffffffff8111156106b8575f80fd5b6020830191508360208285010111156106cf575f80fd5b9250929050565b5f805f80606085870312156106e9575f80fd5b843567ffffffffffffffff8111156106ff575f80fd5b61070b87828801610691565b90955093505060208501356001600160a01b038116811461072a575f80fd5b9396929550929360400135925050565b5f805f6040848603121561074c575f80fd5b833567ffffffffffffffff811115610762575f80fd5b61076e86828701610691565b909790965060209590950135949350505050565b818382375f9101908152919050565b5f602082840312156107a1575f80fd5b5051919050565b818382376305ccae8d60e31b9101908152600401919050565b5f60018060a01b03808716835260208187166020850152818616604085015260806060850152845191508160808501525f5b8281101561080f5785810182015185820160a0015281016107f3565b50505f60a0828501015260a0601f19601f8301168401019150509594505050505056fea264697066735822122003461b17cbbb4633d62eeb338b81f16a91470e520c6a0fb084a02cc280f1758b64736f6c63430008180033",
}

// SpamRegistrarControllerABI is the input ABI used to generate the binding from.
// Deprecated: Use SpamRegistrarControllerMetaData.ABI instead.
var SpamRegistrarControllerABI = SpamRegistrarControllerMetaData.ABI

// SpamRegistrarControllerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SpamRegistrarControllerMetaData.Bin instead.
var SpamRegistrarControllerBin = SpamRegistrarControllerMetaData.Bin

// DeploySpamRegistrarController deploys a new Ethereum contract, binding an instance of SpamRegistrarController to it.
func DeploySpamRegistrarController(auth *bind.TransactOpts, backend bind.ContractBackend, _base common.Address, _ens common.Address, _reverseRegistrar common.Address, _resolver common.Address) (common.Address, *types.Transaction, *SpamRegistrarController, error) {
	parsed, err := SpamRegistrarControllerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SpamRegistrarControllerBin), backend, _base, _ens, _reverseRegistrar, _resolver)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SpamRegistrarController{SpamRegistrarControllerCaller: SpamRegistrarControllerCaller{contract: contract}, SpamRegistrarControllerTransactor: SpamRegistrarControllerTransactor{contract: contract}, SpamRegistrarControllerFilterer: SpamRegistrarControllerFilterer{contract: contract}}, nil
}

// SpamRegistrarController is an auto generated Go binding around an Ethereum contract.
type SpamRegistrarController struct {
	SpamRegistrarControllerCaller     // Read-only binding to the contract
	SpamRegistrarControllerTransactor // Write-only binding to the contract
	SpamRegistrarControllerFilterer   // Log filterer for contract events
}

// SpamRegistrarControllerCaller is an auto generated read-only Go binding around an Ethereum contract.
type SpamRegistrarControllerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SpamRegistrarControllerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SpamRegistrarControllerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SpamRegistrarControllerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SpamRegistrarControllerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SpamRegistrarControllerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SpamRegistrarControllerSession struct {
	Contract     *SpamRegistrarController // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// SpamRegistrarControllerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SpamRegistrarControllerCallerSession struct {
	Contract *SpamRegistrarControllerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// SpamRegistrarControllerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SpamRegistrarControllerTransactorSession struct {
	Contract     *SpamRegistrarControllerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// SpamRegistrarControllerRaw is an auto generated low-level Go binding around an Ethereum contract.
type SpamRegistrarControllerRaw struct {
	Contract *SpamRegistrarController // Generic contract binding to access the raw methods on
}

// SpamRegistrarControllerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SpamRegistrarControllerCallerRaw struct {
	Contract *SpamRegistrarControllerCaller // Generic read-only contract binding to access the raw methods on
}

// SpamRegistrarControllerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SpamRegistrarControllerTransactorRaw struct {
	Contract *SpamRegistrarControllerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSpamRegistrarController creates a new instance of SpamRegistrarController, bound to a specific deployed contract.
func NewSpamRegistrarController(address common.Address, backend bind.ContractBackend) (*SpamRegistrarController, error) {
	contract, err := bindSpamRegistrarController(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SpamRegistrarController{SpamRegistrarControllerCaller: SpamRegistrarControllerCaller{contract: contract}, SpamRegistrarControllerTransactor: SpamRegistrarControllerTransactor{contract: contract}, SpamRegistrarControllerFilterer: SpamRegistrarControllerFilterer{contract: contract}}, nil
}

// NewSpamRegistrarControllerCaller creates a new read-only instance of SpamRegistrarController, bound to a specific deployed contract.
func NewSpamRegistrarControllerCaller(address common.Address, caller bind.ContractCaller) (*SpamRegistrarControllerCaller, error) {
	contract, err := bindSpamRegistrarController(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SpamRegistrarControllerCaller{contract: contract}, nil
}

// NewSpamRegistrarControllerTransactor creates a new write-only instance of SpamRegistrarController, bound to a specific deployed contract.
func NewSpamRegistrarControllerTransactor(address common.Address, transactor bind.ContractTransactor) (*SpamRegistrarControllerTransactor, error) {
	contract, err := bindSpamRegistrarController(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SpamRegistrarControllerTransactor{contract: contract}, nil
}

// NewSpamRegistrarControllerFilterer creates a new log filterer instance of SpamRegistrarController, bound to a specific deployed contract.
func NewSpamRegistrarControllerFilterer(address common.Address, filterer bind.ContractFilterer) (*SpamRegistrarControllerFilterer, error) {
	contract, err := bindSpamRegistrarController(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SpamRegistrarControllerFilterer{contract: contract}, nil
}

// bindSpamRegistrarController binds a generic wrapper to an already deployed contract.
func bindSpamRegistrarController(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SpamRegistrarControllerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SpamRegistrarController *SpamRegistrarControllerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SpamRegistrarController.Contract.SpamRegistrarControllerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SpamRegistrarController *SpamRegistrarControllerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SpamRegistrarController.Contract.SpamRegistrarControllerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SpamRegistrarController *SpamRegistrarControllerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SpamRegistrarController.Contract.SpamRegistrarControllerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SpamRegistrarController *SpamRegistrarControllerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SpamRegistrarController.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SpamRegistrarController *SpamRegistrarControllerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SpamRegistrarController.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SpamRegistrarController *SpamRegistrarControllerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SpamRegistrarController.Contract.contract.Transact(opts, method, params...)
}

// ETHNODE is a free data retrieval call binding the contract method 0xcc473be3.
//
// Solidity: function ETH_NODE() view returns(bytes32)
func (_SpamRegistrarController *SpamRegistrarControllerCaller) ETHNODE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _SpamRegistrarController.contract.Call(opts, &out, "ETH_NODE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ETHNODE is a free data retrieval call binding the contract method 0xcc473be3.
//
// Solidity: function ETH_NODE() view returns(bytes32)
func (_SpamRegistrarController *SpamRegistrarControllerSession) ETHNODE() ([32]byte, error) {
	return _SpamRegistrarController.Contract.ETHNODE(&_SpamRegistrarController.CallOpts)
}

// ETHNODE is a free data retrieval call binding the contract method 0xcc473be3.
//
// Solidity: function ETH_NODE() view returns(bytes32)
func (_SpamRegistrarController *SpamRegistrarControllerCallerSession) ETHNODE() ([32]byte, error) {
	return _SpamRegistrarController.Contract.ETHNODE(&_SpamRegistrarController.CallOpts)
}

// Base is a free data retrieval call binding the contract method 0x5001f3b5.
//
// Solidity: function base() view returns(address)
func (_SpamRegistrarController *SpamRegistrarControllerCaller) Base(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SpamRegistrarController.contract.Call(opts, &out, "base")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Base is a free data retrieval call binding the contract method 0x5001f3b5.
//
// Solidity: function base() view returns(address)
func (_SpamRegistrarController *SpamRegistrarControllerSession) Base() (common.Address, error) {
	return _SpamRegistrarController.Contract.Base(&_SpamRegistrarController.CallOpts)
}

// Base is a free data retrieval call binding the contract method 0x5001f3b5.
//
// Solidity: function base() view returns(address)
func (_SpamRegistrarController *SpamRegistrarControllerCallerSession) Base() (common.Address, error) {
	return _SpamRegistrarController.Contract.Base(&_SpamRegistrarController.CallOpts)
}

// Ens is a free data retrieval call binding the contract method 0x3f15457f.
//
// Solidity: function ens() view returns(address)
func (_SpamRegistrarController *SpamRegistrarControllerCaller) Ens(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SpamRegistrarController.contract.Call(opts, &out, "ens")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Ens is a free data retrieval call binding the contract method 0x3f15457f.
//
// Solidity: function ens() view returns(address)
func (_SpamRegistrarController *SpamRegistrarControllerSession) Ens() (common.Address, error) {
	return _SpamRegistrarController.Contract.Ens(&_SpamRegistrarController.CallOpts)
}

// Ens is a free data retrieval call binding the contract method 0x3f15457f.
//
// Solidity: function ens() view returns(address)
func (_SpamRegistrarController *SpamRegistrarControllerCallerSession) Ens() (common.Address, error) {
	return _SpamRegistrarController.Contract.Ens(&_SpamRegistrarController.CallOpts)
}

// Resolver is a free data retrieval call binding the contract method 0x04f3bcec.
//
// Solidity: function resolver() view returns(address)
func (_SpamRegistrarController *SpamRegistrarControllerCaller) Resolver(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SpamRegistrarController.contract.Call(opts, &out, "resolver")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Resolver is a free data retrieval call binding the contract method 0x04f3bcec.
//
// Solidity: function resolver() view returns(address)
func (_SpamRegistrarController *SpamRegistrarControllerSession) Resolver() (common.Address, error) {
	return _SpamRegistrarController.Contract.Resolver(&_SpamRegistrarController.CallOpts)
}

// Resolver is a free data retrieval call binding the contract method 0x04f3bcec.
//
// Solidity: function resolver() view returns(address)
func (_SpamRegistrarController *SpamRegistrarControllerCallerSession) Resolver() (common.Address, error) {
	return _SpamRegistrarController.Contract.Resolver(&_SpamRegistrarController.CallOpts)
}

// ReverseRegistrar is a free data retrieval call binding the contract method 0x80869853.
//
// Solidity: function reverseRegistrar() view returns(address)
func (_SpamRegistrarController *SpamRegistrarControllerCaller) ReverseRegistrar(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SpamRegistrarController.contract.Call(opts, &out, "reverseRegistrar")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ReverseRegistrar is a free data retrieval call binding the contract method 0x80869853.
//
// Solidity: function reverseRegistrar() view returns(address)
func (_SpamRegistrarController *SpamRegistrarControllerSession) ReverseRegistrar() (common.Address, error) {
	return _SpamRegistrarController.Contract.ReverseRegistrar(&_SpamRegistrarController.CallOpts)
}

// ReverseRegistrar is a free data retrieval call binding the contract method 0x80869853.
//
// Solidity: function reverseRegistrar() view returns(address)
func (_SpamRegistrarController *SpamRegistrarControllerCallerSession) ReverseRegistrar() (common.Address, error) {
	return _SpamRegistrarController.Contract.ReverseRegistrar(&_SpamRegistrarController.CallOpts)
}

// Register is a paid mutator transaction binding the contract method 0xd393c871.
//
// Solidity: function register(string label, address owner, uint256 duration) returns(uint256 expires)
func (_SpamRegistrarController *SpamRegistrarControllerTransactor) Register(opts *bind.TransactOpts, label string, owner common.Address, duration *big.Int) (*types.Transaction, error) {
	return _SpamRegistrarController.contract.Transact(opts, "register", label, owner, duration)
}

// Register is a paid mutator transaction binding the contract method 0xd393c871.
//
// Solidity: function register(string label, address owner, uint256 duration) returns(uint256 expires)
func (_SpamRegistrarController *SpamRegistrarControllerSession) Register(label string, owner common.Address, duration *big.Int) (*types.Transaction, error) {
	return _SpamRegistrarController.Contract.Register(&_SpamRegistrarController.TransactOpts, label, owner, duration)
}

// Register is a paid mutator transaction binding the contract method 0xd393c871.
//
// Solidity: function register(string label, address owner, uint256 duration) returns(uint256 expires)
func (_SpamRegistrarController *SpamRegistrarControllerTransactorSession) Register(label string, owner common.Address, duration *big.Int) (*types.Transaction, error) {
	return _SpamRegistrarController.Contract.Register(&_SpamRegistrarController.TransactOpts, label, owner, duration)
}

// RegisterNamed is a paid mutator transaction binding the contract method 0x8ba778f0.
//
// Solidity: function registerNamed(string label, address addr, uint256 duration) returns(bytes32 node)
func (_SpamRegistrarController *SpamRegistrarControllerTransactor) RegisterNamed(opts *bind.TransactOpts, label string, addr common.Address, duration *big.Int) (*types.Transaction, error) {
	return _SpamRegistrarController.contract.Transact(opts, "registerNamed", label, addr, duration)
}

// RegisterNamed is a paid mutator transaction binding the contract method 0x8ba778f0.
//
// Solidity: function registerNamed(string label, address addr, uint256 duration) returns(bytes32 node)
func (_SpamRegistrarController *SpamRegistrarControllerSession) RegisterNamed(label string, addr common.Address, duration *big.Int) (*types.Transaction, error) {
	return _SpamRegistrarController.Contract.RegisterNamed(&_SpamRegistrarController.TransactOpts, label, addr, duration)
}

// RegisterNamed is a paid mutator transaction binding the contract method 0x8ba778f0.
//
// Solidity: function registerNamed(string label, address addr, uint256 duration) returns(bytes32 node)
func (_SpamRegistrarController *SpamRegistrarControllerTransactorSession) RegisterNamed(label string, addr common.Address, duration *big.Int) (*types.Transaction, error) {
	return _SpamRegistrarController.Contract.RegisterNamed(&_SpamRegistrarController.TransactOpts, label, addr, duration)
}

// Renew is a paid mutator transaction binding the contract method 0xacf1a841.
//
// Solidity: function renew(string label, uint256 duration) returns(uint256 expires)
func (_SpamRegistrarController *SpamRegistrarControllerTransactor) Renew(opts *bind.TransactOpts, label string, duration *big.Int) (*types.Transaction, error) {
	return _SpamRegistrarController.contract.Transact(opts, "renew", label, duration)
}

// Renew is a paid mutator transaction binding the contract method 0xacf1a841.
//
// Solidity: function renew(string label, uint256 duration) returns(uint256 expires)
func (_SpamRegistrarController *SpamRegistrarControllerSession) Renew(label string, duration *big.Int) (*types.Transaction, error) {
	return _SpamRegistrarController.Contract.Renew(&_SpamRegistrarController.TransactOpts, label, duration)
}

// Renew is a paid mutator transaction binding the contract method 0xacf1a841.
//
// Solidity: function renew(string label, uint256 duration) returns(uint256 expires)
func (_SpamRegistrarController *SpamRegistrarControllerTransactorSession) Renew(label string, duration *big.Int) (*types.Transaction, error) {
	return _SpamRegistrarController.Contract.Renew(&_SpamRegistrarController.TransactOpts, label, duration)
}
