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

// DefaultReverseRegistrarMetaData contains all meta data concerning the DefaultReverseRegistrar contract.
var DefaultReverseRegistrarMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"InvalidSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SignatureExpired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SignatureExpiryTooHigh\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"controller\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"name\":\"ControllerChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"NameForAddrChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"controllers\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"nameForAddr\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"controller\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"name\":\"setController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"setName\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"setNameForAddr\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"signatureExpiry\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"setNameForAddrWithSignature\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceID\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600f57600080fd5b50601733601b565b606d565b600180546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b6114458061007c6000396000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c8063c47f002711610076578063da8c229e1161005b578063da8c229e14610176578063e0dba60f14610199578063f2fde38b146101ac57600080fd5b8063c47f002714610150578063c91199411461016357600080fd5b80634ec3bd23116100a75780634ec3bd2314610100578063715018a6146101205780638da5cb5b1461012857600080fd5b8063012a67bc146100c357806301ffc9a7146100d8575b600080fd5b6100d66100d1366004610d5a565b6101bf565b005b6100eb6100e6366004610de9565b610252565b60405190151581526020015b60405180910390f35b61011361010e366004610e2b565b6102ae565b6040516100f79190610eb4565b6100d6610367565b60015460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100f7565b6100d661015e366004610ec7565b61037b565b6100d6610171366004610f09565b61038a565b6100eb610184366004610e2b565b60026020526000908152604090205460ff1681565b6100d66101a7366004610f6a565b61043e565b6100d66101ba366004610e2b565b6104d0565b600061022f3063012a67bc60e01b898989896040516020016101e696959493929190610fa1565b604051602081830303815290604052805190602001207f19457468657265756d205369676e6564204d6573736167653a0a3332000000006000908152601c91909152603c902090565b905061023e838389848a610587565b6102498786866107a0565b50505050505050565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167f0c44feda0000000000000000000000000000000000000000000000000000000014806102a857506102a882610826565b92915050565b73ffffffffffffffffffffffffffffffffffffffff811660009081526020819052604090208054606091906102e29061103c565b80601f016020809104026020016040519081016040528092919081815260200182805461030e9061103c565b801561035b5780601f106103305761010080835404028352916020019161035b565b820191906000526020600020905b81548152906001019060200180831161033e57829003601f168201915b50505050509050919050565b61036f6108bd565b610379600061093e565b565b6103863383836107a0565b5050565b3360009081526002602052604090205460ff1661042e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602860248201527f436f6e74726f6c6c61626c653a2043616c6c6572206973206e6f74206120636f60448201527f6e74726f6c6c657200000000000000000000000000000000000000000000000060648201526084015b60405180910390fd5b6104398383836107a0565b505050565b6104466108bd565b73ffffffffffffffffffffffffffffffffffffffff821660008181526002602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001685151590811790915591519182527f4c97694570a07277810af7e5669ffd5f6a2d6b74b6e9a274b8b870fd5114cf87910160405180910390a25050565b6104d86108bd565b73ffffffffffffffffffffffffffffffffffffffff811661057b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201527f64647265737300000000000000000000000000000000000000000000000000006064820152608401610425565b6105848161093e565b50565b7f649264926492649264926492649264926492649264926492649264926492649285856105b56020826110be565b6105c1928892906110d1565b6105ca916110fb565b036106a3576040517f98ef1ed800000000000000000000000000000000000000000000000000000000815273164af34faf9879394370c7f09064127c043a35e9906398ef1ed89061062590869086908a908a90600401611180565b6020604051808303816000875af1158015610644573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061066891906111b6565b61069e576040517f8baa579f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61071a565b6106e4838387878080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506109b592505050565b61071a576040517f8baa579f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b42811015610754576040517f0819bdcd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61076042610e106111d3565b811115610799576040517f5e4989ee00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050505050565b73ffffffffffffffffffffffffffffffffffffffff831660009081526020819052604090206107d082848361125c565b508273ffffffffffffffffffffffffffffffffffffffff167f8af7a4c7007a33f680904f3b64733396b730fef22d79555dee29801ca2e479a98383604051610819929190611376565b60405180910390a2505050565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167f4ec3bd230000000000000000000000000000000000000000000000000000000014806102a857507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008316146102a8565b60015473ffffffffffffffffffffffffffffffffffffffff163314610379576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610425565b6001805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff0000000000000000000000000000000000000000831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b60008373ffffffffffffffffffffffffffffffffffffffff163b600003610a3e576000806109e38585610a53565b50909250905060008160038111156109fd576109fd611392565b148015610a3557508573ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16145b92505050610a4c565b610a49848484610aa0565b90505b9392505050565b60008060008351604103610a8d5760208401516040850151606086015160001a610a7f88828585610bee565b955095509550505050610a99565b50508151600091506002905b9250925092565b60008060008573ffffffffffffffffffffffffffffffffffffffff168585604051602401610acf9291906113c1565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f1626ba7e0000000000000000000000000000000000000000000000000000000017905251610b5091906113da565b600060405180830381855afa9150503d8060008114610b8b576040519150601f19603f3d011682016040523d82523d6000602084013e610b90565b606091505b5091509150818015610ba457506020815110155b8015610be4575080517f1626ba7e0000000000000000000000000000000000000000000000000000000090610be290830160209081019084016113f6565b145b9695505050505050565b600080807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a0841115610c295750600091506003905082610cde565b604080516000808252602082018084528a905260ff891692820192909252606081018790526080810186905260019060a0016020604051602081039080840390855afa158015610c7d573d6000803e3d6000fd5b50506040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0015191505073ffffffffffffffffffffffffffffffffffffffff8116610cd457506000925060019150829050610cde565b9250600091508190505b9450945094915050565b803573ffffffffffffffffffffffffffffffffffffffff81168114610d0c57600080fd5b919050565b60008083601f840112610d2357600080fd5b50813567ffffffffffffffff811115610d3b57600080fd5b602083019150836020828501011115610d5357600080fd5b9250929050565b60008060008060008060808789031215610d7357600080fd5b610d7c87610ce8565b955060208701359450604087013567ffffffffffffffff811115610d9f57600080fd5b610dab89828a01610d11565b909550935050606087013567ffffffffffffffff811115610dcb57600080fd5b610dd789828a01610d11565b979a9699509497509295939492505050565b600060208284031215610dfb57600080fd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114610a4c57600080fd5b600060208284031215610e3d57600080fd5b610a4c82610ce8565b60005b83811015610e61578181015183820152602001610e49565b50506000910152565b60008151808452610e82816020860160208601610e46565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000610a4c6020830184610e6a565b60008060208385031215610eda57600080fd5b823567ffffffffffffffff811115610ef157600080fd5b610efd85828601610d11565b90969095509350505050565b600080600060408486031215610f1e57600080fd5b610f2784610ce8565b9250602084013567ffffffffffffffff811115610f4357600080fd5b610f4f86828701610d11565b9497909650939450505050565b801515811461058457600080fd5b60008060408385031215610f7d57600080fd5b610f8683610ce8565b91506020830135610f9681610f5c565b809150509250929050565b7fffffffffffffffffffffffffffffffffffffffff0000000000000000000000008760601b1681527fffffffff00000000000000000000000000000000000000000000000000000000861660148201527fffffffffffffffffffffffffffffffffffffffff0000000000000000000000008560601b16601882015283602c8201528183604c83013760009101604c0190815295945050505050565b600181811c9082168061105057607f821691505b602082108103611089577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b818103818111156102a8576102a861108f565b600080858511156110e157600080fd5b838611156110ee57600080fd5b5050820193919092039150565b803560208310156102a8577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff602084900360031b1b1692915050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b73ffffffffffffffffffffffffffffffffffffffff85168152836020820152606060408201526000610be4606083018486611137565b6000602082840312156111c857600080fd5b8151610a4c81610f5c565b808201808211156102a8576102a861108f565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b601f82111561043957806000526020600020601f840160051c8101602085101561123c5750805b601f840160051c820191505b818110156107995760008155600101611248565b67ffffffffffffffff831115611274576112746111e6565b61128883611282835461103c565b83611215565b6000601f8411600181146112da57600085156112a45750838201355b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b1c1916600186901b178355610799565b6000838152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0861690835b828110156113295786850135825560209485019460019092019101611309565b5086821015611364577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88860031b161c19848701351681555b505060018560011b0183555050505050565b60208152600061138a602083018486611137565b949350505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b82815260406020820152600061138a6040830184610e6a565b600082516113ec818460208701610e46565b9190910192915050565b60006020828403121561140857600080fd5b505191905056fea264697066735822122094e05458cc5cf82ca6a2fcf6df1ae7710181c0c199876fab70502c1a27c997f564736f6c634300081a0033",
}

// DefaultReverseRegistrarABI is the input ABI used to generate the binding from.
// Deprecated: Use DefaultReverseRegistrarMetaData.ABI instead.
var DefaultReverseRegistrarABI = DefaultReverseRegistrarMetaData.ABI

// DefaultReverseRegistrarBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DefaultReverseRegistrarMetaData.Bin instead.
var DefaultReverseRegistrarBin = DefaultReverseRegistrarMetaData.Bin

// DeployDefaultReverseRegistrar deploys a new Ethereum contract, binding an instance of DefaultReverseRegistrar to it.
func DeployDefaultReverseRegistrar(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *DefaultReverseRegistrar, error) {
	parsed, err := DefaultReverseRegistrarMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DefaultReverseRegistrarBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DefaultReverseRegistrar{DefaultReverseRegistrarCaller: DefaultReverseRegistrarCaller{contract: contract}, DefaultReverseRegistrarTransactor: DefaultReverseRegistrarTransactor{contract: contract}, DefaultReverseRegistrarFilterer: DefaultReverseRegistrarFilterer{contract: contract}}, nil
}

// DefaultReverseRegistrar is an auto generated Go binding around an Ethereum contract.
type DefaultReverseRegistrar struct {
	DefaultReverseRegistrarCaller     // Read-only binding to the contract
	DefaultReverseRegistrarTransactor // Write-only binding to the contract
	DefaultReverseRegistrarFilterer   // Log filterer for contract events
}

// DefaultReverseRegistrarCaller is an auto generated read-only Go binding around an Ethereum contract.
type DefaultReverseRegistrarCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DefaultReverseRegistrarTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DefaultReverseRegistrarTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DefaultReverseRegistrarFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DefaultReverseRegistrarFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DefaultReverseRegistrarSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DefaultReverseRegistrarSession struct {
	Contract     *DefaultReverseRegistrar // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// DefaultReverseRegistrarCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DefaultReverseRegistrarCallerSession struct {
	Contract *DefaultReverseRegistrarCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// DefaultReverseRegistrarTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DefaultReverseRegistrarTransactorSession struct {
	Contract     *DefaultReverseRegistrarTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// DefaultReverseRegistrarRaw is an auto generated low-level Go binding around an Ethereum contract.
type DefaultReverseRegistrarRaw struct {
	Contract *DefaultReverseRegistrar // Generic contract binding to access the raw methods on
}

// DefaultReverseRegistrarCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DefaultReverseRegistrarCallerRaw struct {
	Contract *DefaultReverseRegistrarCaller // Generic read-only contract binding to access the raw methods on
}

// DefaultReverseRegistrarTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DefaultReverseRegistrarTransactorRaw struct {
	Contract *DefaultReverseRegistrarTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDefaultReverseRegistrar creates a new instance of DefaultReverseRegistrar, bound to a specific deployed contract.
func NewDefaultReverseRegistrar(address common.Address, backend bind.ContractBackend) (*DefaultReverseRegistrar, error) {
	contract, err := bindDefaultReverseRegistrar(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DefaultReverseRegistrar{DefaultReverseRegistrarCaller: DefaultReverseRegistrarCaller{contract: contract}, DefaultReverseRegistrarTransactor: DefaultReverseRegistrarTransactor{contract: contract}, DefaultReverseRegistrarFilterer: DefaultReverseRegistrarFilterer{contract: contract}}, nil
}

// NewDefaultReverseRegistrarCaller creates a new read-only instance of DefaultReverseRegistrar, bound to a specific deployed contract.
func NewDefaultReverseRegistrarCaller(address common.Address, caller bind.ContractCaller) (*DefaultReverseRegistrarCaller, error) {
	contract, err := bindDefaultReverseRegistrar(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DefaultReverseRegistrarCaller{contract: contract}, nil
}

// NewDefaultReverseRegistrarTransactor creates a new write-only instance of DefaultReverseRegistrar, bound to a specific deployed contract.
func NewDefaultReverseRegistrarTransactor(address common.Address, transactor bind.ContractTransactor) (*DefaultReverseRegistrarTransactor, error) {
	contract, err := bindDefaultReverseRegistrar(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DefaultReverseRegistrarTransactor{contract: contract}, nil
}

// NewDefaultReverseRegistrarFilterer creates a new log filterer instance of DefaultReverseRegistrar, bound to a specific deployed contract.
func NewDefaultReverseRegistrarFilterer(address common.Address, filterer bind.ContractFilterer) (*DefaultReverseRegistrarFilterer, error) {
	contract, err := bindDefaultReverseRegistrar(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DefaultReverseRegistrarFilterer{contract: contract}, nil
}

// bindDefaultReverseRegistrar binds a generic wrapper to an already deployed contract.
func bindDefaultReverseRegistrar(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DefaultReverseRegistrarMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DefaultReverseRegistrar *DefaultReverseRegistrarRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DefaultReverseRegistrar.Contract.DefaultReverseRegistrarCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DefaultReverseRegistrar *DefaultReverseRegistrarRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DefaultReverseRegistrar.Contract.DefaultReverseRegistrarTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DefaultReverseRegistrar *DefaultReverseRegistrarRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DefaultReverseRegistrar.Contract.DefaultReverseRegistrarTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DefaultReverseRegistrar *DefaultReverseRegistrarCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DefaultReverseRegistrar.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DefaultReverseRegistrar *DefaultReverseRegistrarTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DefaultReverseRegistrar.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DefaultReverseRegistrar *DefaultReverseRegistrarTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DefaultReverseRegistrar.Contract.contract.Transact(opts, method, params...)
}

// Controllers is a free data retrieval call binding the contract method 0xda8c229e.
//
// Solidity: function controllers(address ) view returns(bool)
func (_DefaultReverseRegistrar *DefaultReverseRegistrarCaller) Controllers(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _DefaultReverseRegistrar.contract.Call(opts, &out, "controllers", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Controllers is a free data retrieval call binding the contract method 0xda8c229e.
//
// Solidity: function controllers(address ) view returns(bool)
func (_DefaultReverseRegistrar *DefaultReverseRegistrarSession) Controllers(arg0 common.Address) (bool, error) {
	return _DefaultReverseRegistrar.Contract.Controllers(&_DefaultReverseRegistrar.CallOpts, arg0)
}

// Controllers is a free data retrieval call binding the contract method 0xda8c229e.
//
// Solidity: function controllers(address ) view returns(bool)
func (_DefaultReverseRegistrar *DefaultReverseRegistrarCallerSession) Controllers(arg0 common.Address) (bool, error) {
	return _DefaultReverseRegistrar.Contract.Controllers(&_DefaultReverseRegistrar.CallOpts, arg0)
}

// NameForAddr is a free data retrieval call binding the contract method 0x4ec3bd23.
//
// Solidity: function nameForAddr(address addr) view returns(string name)
func (_DefaultReverseRegistrar *DefaultReverseRegistrarCaller) NameForAddr(opts *bind.CallOpts, addr common.Address) (string, error) {
	var out []interface{}
	err := _DefaultReverseRegistrar.contract.Call(opts, &out, "nameForAddr", addr)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// NameForAddr is a free data retrieval call binding the contract method 0x4ec3bd23.
//
// Solidity: function nameForAddr(address addr) view returns(string name)
func (_DefaultReverseRegistrar *DefaultReverseRegistrarSession) NameForAddr(addr common.Address) (string, error) {
	return _DefaultReverseRegistrar.Contract.NameForAddr(&_DefaultReverseRegistrar.CallOpts, addr)
}

// NameForAddr is a free data retrieval call binding the contract method 0x4ec3bd23.
//
// Solidity: function nameForAddr(address addr) view returns(string name)
func (_DefaultReverseRegistrar *DefaultReverseRegistrarCallerSession) NameForAddr(addr common.Address) (string, error) {
	return _DefaultReverseRegistrar.Contract.NameForAddr(&_DefaultReverseRegistrar.CallOpts, addr)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DefaultReverseRegistrar *DefaultReverseRegistrarCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DefaultReverseRegistrar.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DefaultReverseRegistrar *DefaultReverseRegistrarSession) Owner() (common.Address, error) {
	return _DefaultReverseRegistrar.Contract.Owner(&_DefaultReverseRegistrar.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DefaultReverseRegistrar *DefaultReverseRegistrarCallerSession) Owner() (common.Address, error) {
	return _DefaultReverseRegistrar.Contract.Owner(&_DefaultReverseRegistrar.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceID) view returns(bool)
func (_DefaultReverseRegistrar *DefaultReverseRegistrarCaller) SupportsInterface(opts *bind.CallOpts, interfaceID [4]byte) (bool, error) {
	var out []interface{}
	err := _DefaultReverseRegistrar.contract.Call(opts, &out, "supportsInterface", interfaceID)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceID) view returns(bool)
func (_DefaultReverseRegistrar *DefaultReverseRegistrarSession) SupportsInterface(interfaceID [4]byte) (bool, error) {
	return _DefaultReverseRegistrar.Contract.SupportsInterface(&_DefaultReverseRegistrar.CallOpts, interfaceID)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceID) view returns(bool)
func (_DefaultReverseRegistrar *DefaultReverseRegistrarCallerSession) SupportsInterface(interfaceID [4]byte) (bool, error) {
	return _DefaultReverseRegistrar.Contract.SupportsInterface(&_DefaultReverseRegistrar.CallOpts, interfaceID)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DefaultReverseRegistrar *DefaultReverseRegistrarTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DefaultReverseRegistrar.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DefaultReverseRegistrar *DefaultReverseRegistrarSession) RenounceOwnership() (*types.Transaction, error) {
	return _DefaultReverseRegistrar.Contract.RenounceOwnership(&_DefaultReverseRegistrar.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DefaultReverseRegistrar *DefaultReverseRegistrarTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _DefaultReverseRegistrar.Contract.RenounceOwnership(&_DefaultReverseRegistrar.TransactOpts)
}

// SetController is a paid mutator transaction binding the contract method 0xe0dba60f.
//
// Solidity: function setController(address controller, bool enabled) returns()
func (_DefaultReverseRegistrar *DefaultReverseRegistrarTransactor) SetController(opts *bind.TransactOpts, controller common.Address, enabled bool) (*types.Transaction, error) {
	return _DefaultReverseRegistrar.contract.Transact(opts, "setController", controller, enabled)
}

// SetController is a paid mutator transaction binding the contract method 0xe0dba60f.
//
// Solidity: function setController(address controller, bool enabled) returns()
func (_DefaultReverseRegistrar *DefaultReverseRegistrarSession) SetController(controller common.Address, enabled bool) (*types.Transaction, error) {
	return _DefaultReverseRegistrar.Contract.SetController(&_DefaultReverseRegistrar.TransactOpts, controller, enabled)
}

// SetController is a paid mutator transaction binding the contract method 0xe0dba60f.
//
// Solidity: function setController(address controller, bool enabled) returns()
func (_DefaultReverseRegistrar *DefaultReverseRegistrarTransactorSession) SetController(controller common.Address, enabled bool) (*types.Transaction, error) {
	return _DefaultReverseRegistrar.Contract.SetController(&_DefaultReverseRegistrar.TransactOpts, controller, enabled)
}

// SetName is a paid mutator transaction binding the contract method 0xc47f0027.
//
// Solidity: function setName(string name) returns()
func (_DefaultReverseRegistrar *DefaultReverseRegistrarTransactor) SetName(opts *bind.TransactOpts, name string) (*types.Transaction, error) {
	return _DefaultReverseRegistrar.contract.Transact(opts, "setName", name)
}

// SetName is a paid mutator transaction binding the contract method 0xc47f0027.
//
// Solidity: function setName(string name) returns()
func (_DefaultReverseRegistrar *DefaultReverseRegistrarSession) SetName(name string) (*types.Transaction, error) {
	return _DefaultReverseRegistrar.Contract.SetName(&_DefaultReverseRegistrar.TransactOpts, name)
}

// SetName is a paid mutator transaction binding the contract method 0xc47f0027.
//
// Solidity: function setName(string name) returns()
func (_DefaultReverseRegistrar *DefaultReverseRegistrarTransactorSession) SetName(name string) (*types.Transaction, error) {
	return _DefaultReverseRegistrar.Contract.SetName(&_DefaultReverseRegistrar.TransactOpts, name)
}

// SetNameForAddr is a paid mutator transaction binding the contract method 0xc9119941.
//
// Solidity: function setNameForAddr(address addr, string name) returns()
func (_DefaultReverseRegistrar *DefaultReverseRegistrarTransactor) SetNameForAddr(opts *bind.TransactOpts, addr common.Address, name string) (*types.Transaction, error) {
	return _DefaultReverseRegistrar.contract.Transact(opts, "setNameForAddr", addr, name)
}

// SetNameForAddr is a paid mutator transaction binding the contract method 0xc9119941.
//
// Solidity: function setNameForAddr(address addr, string name) returns()
func (_DefaultReverseRegistrar *DefaultReverseRegistrarSession) SetNameForAddr(addr common.Address, name string) (*types.Transaction, error) {
	return _DefaultReverseRegistrar.Contract.SetNameForAddr(&_DefaultReverseRegistrar.TransactOpts, addr, name)
}

// SetNameForAddr is a paid mutator transaction binding the contract method 0xc9119941.
//
// Solidity: function setNameForAddr(address addr, string name) returns()
func (_DefaultReverseRegistrar *DefaultReverseRegistrarTransactorSession) SetNameForAddr(addr common.Address, name string) (*types.Transaction, error) {
	return _DefaultReverseRegistrar.Contract.SetNameForAddr(&_DefaultReverseRegistrar.TransactOpts, addr, name)
}

// SetNameForAddrWithSignature is a paid mutator transaction binding the contract method 0x012a67bc.
//
// Solidity: function setNameForAddrWithSignature(address addr, uint256 signatureExpiry, string name, bytes signature) returns()
func (_DefaultReverseRegistrar *DefaultReverseRegistrarTransactor) SetNameForAddrWithSignature(opts *bind.TransactOpts, addr common.Address, signatureExpiry *big.Int, name string, signature []byte) (*types.Transaction, error) {
	return _DefaultReverseRegistrar.contract.Transact(opts, "setNameForAddrWithSignature", addr, signatureExpiry, name, signature)
}

// SetNameForAddrWithSignature is a paid mutator transaction binding the contract method 0x012a67bc.
//
// Solidity: function setNameForAddrWithSignature(address addr, uint256 signatureExpiry, string name, bytes signature) returns()
func (_DefaultReverseRegistrar *DefaultReverseRegistrarSession) SetNameForAddrWithSignature(addr common.Address, signatureExpiry *big.Int, name string, signature []byte) (*types.Transaction, error) {
	return _DefaultReverseRegistrar.Contract.SetNameForAddrWithSignature(&_DefaultReverseRegistrar.TransactOpts, addr, signatureExpiry, name, signature)
}

// SetNameForAddrWithSignature is a paid mutator transaction binding the contract method 0x012a67bc.
//
// Solidity: function setNameForAddrWithSignature(address addr, uint256 signatureExpiry, string name, bytes signature) returns()
func (_DefaultReverseRegistrar *DefaultReverseRegistrarTransactorSession) SetNameForAddrWithSignature(addr common.Address, signatureExpiry *big.Int, name string, signature []byte) (*types.Transaction, error) {
	return _DefaultReverseRegistrar.Contract.SetNameForAddrWithSignature(&_DefaultReverseRegistrar.TransactOpts, addr, signatureExpiry, name, signature)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DefaultReverseRegistrar *DefaultReverseRegistrarTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _DefaultReverseRegistrar.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DefaultReverseRegistrar *DefaultReverseRegistrarSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _DefaultReverseRegistrar.Contract.TransferOwnership(&_DefaultReverseRegistrar.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DefaultReverseRegistrar *DefaultReverseRegistrarTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _DefaultReverseRegistrar.Contract.TransferOwnership(&_DefaultReverseRegistrar.TransactOpts, newOwner)
}

// DefaultReverseRegistrarControllerChangedIterator is returned from FilterControllerChanged and is used to iterate over the raw logs and unpacked data for ControllerChanged events raised by the DefaultReverseRegistrar contract.
type DefaultReverseRegistrarControllerChangedIterator struct {
	Event *DefaultReverseRegistrarControllerChanged // Event containing the contract specifics and raw log

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
func (it *DefaultReverseRegistrarControllerChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DefaultReverseRegistrarControllerChanged)
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
		it.Event = new(DefaultReverseRegistrarControllerChanged)
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
func (it *DefaultReverseRegistrarControllerChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DefaultReverseRegistrarControllerChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DefaultReverseRegistrarControllerChanged represents a ControllerChanged event raised by the DefaultReverseRegistrar contract.
type DefaultReverseRegistrarControllerChanged struct {
	Controller common.Address
	Enabled    bool
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterControllerChanged is a free log retrieval operation binding the contract event 0x4c97694570a07277810af7e5669ffd5f6a2d6b74b6e9a274b8b870fd5114cf87.
//
// Solidity: event ControllerChanged(address indexed controller, bool enabled)
func (_DefaultReverseRegistrar *DefaultReverseRegistrarFilterer) FilterControllerChanged(opts *bind.FilterOpts, controller []common.Address) (*DefaultReverseRegistrarControllerChangedIterator, error) {

	var controllerRule []interface{}
	for _, controllerItem := range controller {
		controllerRule = append(controllerRule, controllerItem)
	}

	logs, sub, err := _DefaultReverseRegistrar.contract.FilterLogs(opts, "ControllerChanged", controllerRule)
	if err != nil {
		return nil, err
	}
	return &DefaultReverseRegistrarControllerChangedIterator{contract: _DefaultReverseRegistrar.contract, event: "ControllerChanged", logs: logs, sub: sub}, nil
}

// WatchControllerChanged is a free log subscription operation binding the contract event 0x4c97694570a07277810af7e5669ffd5f6a2d6b74b6e9a274b8b870fd5114cf87.
//
// Solidity: event ControllerChanged(address indexed controller, bool enabled)
func (_DefaultReverseRegistrar *DefaultReverseRegistrarFilterer) WatchControllerChanged(opts *bind.WatchOpts, sink chan<- *DefaultReverseRegistrarControllerChanged, controller []common.Address) (event.Subscription, error) {

	var controllerRule []interface{}
	for _, controllerItem := range controller {
		controllerRule = append(controllerRule, controllerItem)
	}

	logs, sub, err := _DefaultReverseRegistrar.contract.WatchLogs(opts, "ControllerChanged", controllerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DefaultReverseRegistrarControllerChanged)
				if err := _DefaultReverseRegistrar.contract.UnpackLog(event, "ControllerChanged", log); err != nil {
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

// ParseControllerChanged is a log parse operation binding the contract event 0x4c97694570a07277810af7e5669ffd5f6a2d6b74b6e9a274b8b870fd5114cf87.
//
// Solidity: event ControllerChanged(address indexed controller, bool enabled)
func (_DefaultReverseRegistrar *DefaultReverseRegistrarFilterer) ParseControllerChanged(log types.Log) (*DefaultReverseRegistrarControllerChanged, error) {
	event := new(DefaultReverseRegistrarControllerChanged)
	if err := _DefaultReverseRegistrar.contract.UnpackLog(event, "ControllerChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DefaultReverseRegistrarNameForAddrChangedIterator is returned from FilterNameForAddrChanged and is used to iterate over the raw logs and unpacked data for NameForAddrChanged events raised by the DefaultReverseRegistrar contract.
type DefaultReverseRegistrarNameForAddrChangedIterator struct {
	Event *DefaultReverseRegistrarNameForAddrChanged // Event containing the contract specifics and raw log

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
func (it *DefaultReverseRegistrarNameForAddrChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DefaultReverseRegistrarNameForAddrChanged)
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
		it.Event = new(DefaultReverseRegistrarNameForAddrChanged)
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
func (it *DefaultReverseRegistrarNameForAddrChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DefaultReverseRegistrarNameForAddrChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DefaultReverseRegistrarNameForAddrChanged represents a NameForAddrChanged event raised by the DefaultReverseRegistrar contract.
type DefaultReverseRegistrarNameForAddrChanged struct {
	Addr common.Address
	Name string
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterNameForAddrChanged is a free log retrieval operation binding the contract event 0x8af7a4c7007a33f680904f3b64733396b730fef22d79555dee29801ca2e479a9.
//
// Solidity: event NameForAddrChanged(address indexed addr, string name)
func (_DefaultReverseRegistrar *DefaultReverseRegistrarFilterer) FilterNameForAddrChanged(opts *bind.FilterOpts, addr []common.Address) (*DefaultReverseRegistrarNameForAddrChangedIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _DefaultReverseRegistrar.contract.FilterLogs(opts, "NameForAddrChanged", addrRule)
	if err != nil {
		return nil, err
	}
	return &DefaultReverseRegistrarNameForAddrChangedIterator{contract: _DefaultReverseRegistrar.contract, event: "NameForAddrChanged", logs: logs, sub: sub}, nil
}

// WatchNameForAddrChanged is a free log subscription operation binding the contract event 0x8af7a4c7007a33f680904f3b64733396b730fef22d79555dee29801ca2e479a9.
//
// Solidity: event NameForAddrChanged(address indexed addr, string name)
func (_DefaultReverseRegistrar *DefaultReverseRegistrarFilterer) WatchNameForAddrChanged(opts *bind.WatchOpts, sink chan<- *DefaultReverseRegistrarNameForAddrChanged, addr []common.Address) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _DefaultReverseRegistrar.contract.WatchLogs(opts, "NameForAddrChanged", addrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DefaultReverseRegistrarNameForAddrChanged)
				if err := _DefaultReverseRegistrar.contract.UnpackLog(event, "NameForAddrChanged", log); err != nil {
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

// ParseNameForAddrChanged is a log parse operation binding the contract event 0x8af7a4c7007a33f680904f3b64733396b730fef22d79555dee29801ca2e479a9.
//
// Solidity: event NameForAddrChanged(address indexed addr, string name)
func (_DefaultReverseRegistrar *DefaultReverseRegistrarFilterer) ParseNameForAddrChanged(log types.Log) (*DefaultReverseRegistrarNameForAddrChanged, error) {
	event := new(DefaultReverseRegistrarNameForAddrChanged)
	if err := _DefaultReverseRegistrar.contract.UnpackLog(event, "NameForAddrChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DefaultReverseRegistrarOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the DefaultReverseRegistrar contract.
type DefaultReverseRegistrarOwnershipTransferredIterator struct {
	Event *DefaultReverseRegistrarOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *DefaultReverseRegistrarOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DefaultReverseRegistrarOwnershipTransferred)
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
		it.Event = new(DefaultReverseRegistrarOwnershipTransferred)
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
func (it *DefaultReverseRegistrarOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DefaultReverseRegistrarOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DefaultReverseRegistrarOwnershipTransferred represents a OwnershipTransferred event raised by the DefaultReverseRegistrar contract.
type DefaultReverseRegistrarOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DefaultReverseRegistrar *DefaultReverseRegistrarFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*DefaultReverseRegistrarOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DefaultReverseRegistrar.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &DefaultReverseRegistrarOwnershipTransferredIterator{contract: _DefaultReverseRegistrar.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DefaultReverseRegistrar *DefaultReverseRegistrarFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DefaultReverseRegistrarOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DefaultReverseRegistrar.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DefaultReverseRegistrarOwnershipTransferred)
				if err := _DefaultReverseRegistrar.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DefaultReverseRegistrar *DefaultReverseRegistrarFilterer) ParseOwnershipTransferred(log types.Log) (*DefaultReverseRegistrarOwnershipTransferred, error) {
	event := new(DefaultReverseRegistrarOwnershipTransferred)
	if err := _DefaultReverseRegistrar.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
