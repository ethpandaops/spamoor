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

// ReverseRegistrarMetaData contains all meta data concerning the ReverseRegistrar contract.
var ReverseRegistrarMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractENS\",\"name\":\"ensAddr\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"controller\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"name\":\"ControllerChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractNameResolver\",\"name\":\"resolver\",\"type\":\"address\"}],\"name\":\"DefaultResolverChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"}],\"name\":\"ReverseClaimed\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"claim\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"resolver\",\"type\":\"address\"}],\"name\":\"claimForAddr\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"resolver\",\"type\":\"address\"}],\"name\":\"claimWithResolver\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"controllers\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"defaultResolver\",\"outputs\":[{\"internalType\":\"contractNameResolver\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ens\",\"outputs\":[{\"internalType\":\"contractENS\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"node\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"controller\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"name\":\"setController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"resolver\",\"type\":\"address\"}],\"name\":\"setDefaultResolver\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"setName\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"resolver\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"setNameForAddr\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b5060405161119938038061119983398101604081905261002f916101b4565b6100383361014c565b6001600160a01b03811660808190526040516302571be360e01b81527f91d1777781884d03a6757a803996e38de2a42967fb37eeaca72729271025a9e26004820152600091906302571be390602401602060405180830381865afa1580156100a4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906100c891906101b4565b90506001600160a01b0381161561014557604051630f41a04d60e11b81523360048201526001600160a01b03821690631e83409a906024016020604051808303816000875af115801561011f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061014391906101d8565b505b50506101f1565b600080546001600160a01b038381166001600160a01b0319831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b6001600160a01b03811681146101b157600080fd5b50565b6000602082840312156101c657600080fd5b81516101d18161019c565b9392505050565b6000602082840312156101ea57600080fd5b5051919050565b608051610f7f61021a6000396000818161012d0152818161033e01526105890152610f7f6000f3fe608060405234801561001057600080fd5b50600436106100ea5760003560e01c80638da5cb5b1161008c578063c66485b211610066578063c66485b214610208578063da8c229e1461021b578063e0dba60f1461024e578063f2fde38b1461026157600080fd5b80638da5cb5b146101c4578063bffbe61c146101e2578063c47f0027146101f557600080fd5b806365669631116100c85780636566963114610174578063715018a6146101875780637a806d6b14610191578063828eab0e146101a457600080fd5b80630f5a5466146100ef5780631e83409a146101155780633f15457f14610128575b600080fd5b6101026100fd366004610c12565b610274565b6040519081526020015b60405180910390f35b610102610123366004610c4b565b610288565b61014f7f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161010c565b610102610182366004610c68565b6102b7565b61018f6105f0565b005b61010261019f366004610dae565b610604565b60025461014f9073ffffffffffffffffffffffffffffffffffffffff1681565b60005473ffffffffffffffffffffffffffffffffffffffff1661014f565b6101026101f0366004610c4b565b6106a5565b610102610203366004610e23565b610700565b61018f610216366004610c4b565b61072a565b61023e610229366004610c4b565b60016020526000908152604090205460ff1681565b604051901515815260200161010c565b61018f61025c366004610e6e565b610844565b61018f61026f366004610c4b565b6108d6565b60006102813384846102b7565b9392505050565b6002546000906102b1903390849073ffffffffffffffffffffffffffffffffffffffff166102b7565b92915050565b60008373ffffffffffffffffffffffffffffffffffffffff81163314806102ed57503360009081526001602052604090205460ff165b806103a957506040517fe985e9c500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82811660048301523360248301527f0000000000000000000000000000000000000000000000000000000000000000169063e985e9c590604401602060405180830381865afa158015610385573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103a99190610e9c565b806103b857506103b88161098d565b61046f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152605b60248201527f526576657273655265676973747261723a2043616c6c6572206973206e6f742060448201527f6120636f6e74726f6c6c6572206f7220617574686f726973656420627920616460648201527f6472657373206f7220746865206164647265737320697473656c660000000000608482015260a4015b60405180910390fd5b600061047a86610a3e565b604080517f91d1777781884d03a6757a803996e38de2a42967fb37eeaca72729271025a9e26020808301919091528183018490528251808303840181526060909201928390528151910120919250819073ffffffffffffffffffffffffffffffffffffffff8916907f6ada868dd3058cf77a48a74489fd7963688e5464b2b0fa957ace976243270e9290600090a36040517f5ef2c7f00000000000000000000000000000000000000000000000000000000081527f91d1777781884d03a6757a803996e38de2a42967fb37eeaca72729271025a9e260048201526024810183905273ffffffffffffffffffffffffffffffffffffffff87811660448301528681166064830152600060848301527f00000000000000000000000000000000000000000000000000000000000000001690635ef2c7f09060a401600060405180830381600087803b1580156105cd57600080fd5b505af11580156105e1573d6000803e3d6000fd5b50929998505050505050505050565b6105f8610afa565b6106026000610b7b565b565b6000806106128686866102b7565b6040517f7737221300000000000000000000000000000000000000000000000000000000815290915073ffffffffffffffffffffffffffffffffffffffff8516906377372213906106699084908790600401610eb9565b600060405180830381600087803b15801561068357600080fd5b505af1158015610697573d6000803e3d6000fd5b509298975050505050505050565b60007f91d1777781884d03a6757a803996e38de2a42967fb37eeaca72729271025a9e26106d183610a3e565b604080516020810193909352820152606001604051602081830303815290604052805190602001209050919050565b6002546000906102b1903390819073ffffffffffffffffffffffffffffffffffffffff1685610604565b610732610afa565b73ffffffffffffffffffffffffffffffffffffffff81166107d5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603060248201527f526576657273655265676973747261723a205265736f6c76657220616464726560448201527f7373206d757374206e6f742062652030000000000000000000000000000000006064820152608401610466565b600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040517feae17a84d9eb83d8c8eb317f9e7d64857bc363fa51674d996c023f4340c577cf90600090a250565b61084c610afa565b73ffffffffffffffffffffffffffffffffffffffff821660008181526001602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001685151590811790915591519182527f4c97694570a07277810af7e5669ffd5f6a2d6b74b6e9a274b8b870fd5114cf87910160405180910390a25050565b6108de610afa565b73ffffffffffffffffffffffffffffffffffffffff8116610981576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201527f64647265737300000000000000000000000000000000000000000000000000006064820152608401610466565b61098a81610b7b565b50565b60008173ffffffffffffffffffffffffffffffffffffffff16638da5cb5b6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015610a14575060408051601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201909252610a1191810190610f2c565b60015b610a2057506000919050565b73ffffffffffffffffffffffffffffffffffffffff16331492915050565b600060285b8015610aee577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff017f3031323334353637383961626364656600000000000000000000000000000000600f84161a81536010909204917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff017f3031323334353637383961626364656600000000000000000000000000000000600f84161a8153601083049250610a43565b50506028600020919050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610602576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610466565b6000805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff0000000000000000000000000000000000000000831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b73ffffffffffffffffffffffffffffffffffffffff8116811461098a57600080fd5b60008060408385031215610c2557600080fd5b8235610c3081610bf0565b91506020830135610c4081610bf0565b809150509250929050565b600060208284031215610c5d57600080fd5b813561028181610bf0565b600080600060608486031215610c7d57600080fd5b8335610c8881610bf0565b92506020840135610c9881610bf0565b91506040840135610ca881610bf0565b809150509250925092565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600082601f830112610cf357600080fd5b813567ffffffffffffffff811115610d0d57610d0d610cb3565b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8501160116810181811067ffffffffffffffff82111715610d7957610d79610cb3565b604052818152838201602001851015610d9157600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008060808587031215610dc457600080fd5b8435610dcf81610bf0565b93506020850135610ddf81610bf0565b92506040850135610def81610bf0565b9150606085013567ffffffffffffffff811115610e0b57600080fd5b610e1787828801610ce2565b91505092959194509250565b600060208284031215610e3557600080fd5b813567ffffffffffffffff811115610e4c57600080fd5b610e5884828501610ce2565b949350505050565b801515811461098a57600080fd5b60008060408385031215610e8157600080fd5b8235610e8c81610bf0565b91506020830135610c4081610e60565b600060208284031215610eae57600080fd5b815161028181610e60565b828152604060208201526000825180604084015260005b81811015610eed5760208186018101516060868401015201610ed0565b5060006060828501015260607fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168401019150509392505050565b600060208284031215610f3e57600080fd5b815161028181610bf056fea2646970667358221220ffa2da93664a067cb12464d34f551bffbe591fedba1524d94ecde426edaf6b6e64736f6c634300081a0033",
}

// ReverseRegistrarABI is the input ABI used to generate the binding from.
// Deprecated: Use ReverseRegistrarMetaData.ABI instead.
var ReverseRegistrarABI = ReverseRegistrarMetaData.ABI

// ReverseRegistrarBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ReverseRegistrarMetaData.Bin instead.
var ReverseRegistrarBin = ReverseRegistrarMetaData.Bin

// DeployReverseRegistrar deploys a new Ethereum contract, binding an instance of ReverseRegistrar to it.
func DeployReverseRegistrar(auth *bind.TransactOpts, backend bind.ContractBackend, ensAddr common.Address) (common.Address, *types.Transaction, *ReverseRegistrar, error) {
	parsed, err := ReverseRegistrarMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ReverseRegistrarBin), backend, ensAddr)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ReverseRegistrar{ReverseRegistrarCaller: ReverseRegistrarCaller{contract: contract}, ReverseRegistrarTransactor: ReverseRegistrarTransactor{contract: contract}, ReverseRegistrarFilterer: ReverseRegistrarFilterer{contract: contract}}, nil
}

// ReverseRegistrar is an auto generated Go binding around an Ethereum contract.
type ReverseRegistrar struct {
	ReverseRegistrarCaller     // Read-only binding to the contract
	ReverseRegistrarTransactor // Write-only binding to the contract
	ReverseRegistrarFilterer   // Log filterer for contract events
}

// ReverseRegistrarCaller is an auto generated read-only Go binding around an Ethereum contract.
type ReverseRegistrarCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReverseRegistrarTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ReverseRegistrarTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReverseRegistrarFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ReverseRegistrarFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReverseRegistrarSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ReverseRegistrarSession struct {
	Contract     *ReverseRegistrar // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ReverseRegistrarCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ReverseRegistrarCallerSession struct {
	Contract *ReverseRegistrarCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// ReverseRegistrarTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ReverseRegistrarTransactorSession struct {
	Contract     *ReverseRegistrarTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// ReverseRegistrarRaw is an auto generated low-level Go binding around an Ethereum contract.
type ReverseRegistrarRaw struct {
	Contract *ReverseRegistrar // Generic contract binding to access the raw methods on
}

// ReverseRegistrarCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ReverseRegistrarCallerRaw struct {
	Contract *ReverseRegistrarCaller // Generic read-only contract binding to access the raw methods on
}

// ReverseRegistrarTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ReverseRegistrarTransactorRaw struct {
	Contract *ReverseRegistrarTransactor // Generic write-only contract binding to access the raw methods on
}

// NewReverseRegistrar creates a new instance of ReverseRegistrar, bound to a specific deployed contract.
func NewReverseRegistrar(address common.Address, backend bind.ContractBackend) (*ReverseRegistrar, error) {
	contract, err := bindReverseRegistrar(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ReverseRegistrar{ReverseRegistrarCaller: ReverseRegistrarCaller{contract: contract}, ReverseRegistrarTransactor: ReverseRegistrarTransactor{contract: contract}, ReverseRegistrarFilterer: ReverseRegistrarFilterer{contract: contract}}, nil
}

// NewReverseRegistrarCaller creates a new read-only instance of ReverseRegistrar, bound to a specific deployed contract.
func NewReverseRegistrarCaller(address common.Address, caller bind.ContractCaller) (*ReverseRegistrarCaller, error) {
	contract, err := bindReverseRegistrar(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ReverseRegistrarCaller{contract: contract}, nil
}

// NewReverseRegistrarTransactor creates a new write-only instance of ReverseRegistrar, bound to a specific deployed contract.
func NewReverseRegistrarTransactor(address common.Address, transactor bind.ContractTransactor) (*ReverseRegistrarTransactor, error) {
	contract, err := bindReverseRegistrar(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ReverseRegistrarTransactor{contract: contract}, nil
}

// NewReverseRegistrarFilterer creates a new log filterer instance of ReverseRegistrar, bound to a specific deployed contract.
func NewReverseRegistrarFilterer(address common.Address, filterer bind.ContractFilterer) (*ReverseRegistrarFilterer, error) {
	contract, err := bindReverseRegistrar(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ReverseRegistrarFilterer{contract: contract}, nil
}

// bindReverseRegistrar binds a generic wrapper to an already deployed contract.
func bindReverseRegistrar(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ReverseRegistrarMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ReverseRegistrar *ReverseRegistrarRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ReverseRegistrar.Contract.ReverseRegistrarCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ReverseRegistrar *ReverseRegistrarRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReverseRegistrar.Contract.ReverseRegistrarTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ReverseRegistrar *ReverseRegistrarRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ReverseRegistrar.Contract.ReverseRegistrarTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ReverseRegistrar *ReverseRegistrarCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ReverseRegistrar.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ReverseRegistrar *ReverseRegistrarTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReverseRegistrar.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ReverseRegistrar *ReverseRegistrarTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ReverseRegistrar.Contract.contract.Transact(opts, method, params...)
}

// Controllers is a free data retrieval call binding the contract method 0xda8c229e.
//
// Solidity: function controllers(address ) view returns(bool)
func (_ReverseRegistrar *ReverseRegistrarCaller) Controllers(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _ReverseRegistrar.contract.Call(opts, &out, "controllers", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Controllers is a free data retrieval call binding the contract method 0xda8c229e.
//
// Solidity: function controllers(address ) view returns(bool)
func (_ReverseRegistrar *ReverseRegistrarSession) Controllers(arg0 common.Address) (bool, error) {
	return _ReverseRegistrar.Contract.Controllers(&_ReverseRegistrar.CallOpts, arg0)
}

// Controllers is a free data retrieval call binding the contract method 0xda8c229e.
//
// Solidity: function controllers(address ) view returns(bool)
func (_ReverseRegistrar *ReverseRegistrarCallerSession) Controllers(arg0 common.Address) (bool, error) {
	return _ReverseRegistrar.Contract.Controllers(&_ReverseRegistrar.CallOpts, arg0)
}

// DefaultResolver is a free data retrieval call binding the contract method 0x828eab0e.
//
// Solidity: function defaultResolver() view returns(address)
func (_ReverseRegistrar *ReverseRegistrarCaller) DefaultResolver(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ReverseRegistrar.contract.Call(opts, &out, "defaultResolver")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DefaultResolver is a free data retrieval call binding the contract method 0x828eab0e.
//
// Solidity: function defaultResolver() view returns(address)
func (_ReverseRegistrar *ReverseRegistrarSession) DefaultResolver() (common.Address, error) {
	return _ReverseRegistrar.Contract.DefaultResolver(&_ReverseRegistrar.CallOpts)
}

// DefaultResolver is a free data retrieval call binding the contract method 0x828eab0e.
//
// Solidity: function defaultResolver() view returns(address)
func (_ReverseRegistrar *ReverseRegistrarCallerSession) DefaultResolver() (common.Address, error) {
	return _ReverseRegistrar.Contract.DefaultResolver(&_ReverseRegistrar.CallOpts)
}

// Ens is a free data retrieval call binding the contract method 0x3f15457f.
//
// Solidity: function ens() view returns(address)
func (_ReverseRegistrar *ReverseRegistrarCaller) Ens(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ReverseRegistrar.contract.Call(opts, &out, "ens")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Ens is a free data retrieval call binding the contract method 0x3f15457f.
//
// Solidity: function ens() view returns(address)
func (_ReverseRegistrar *ReverseRegistrarSession) Ens() (common.Address, error) {
	return _ReverseRegistrar.Contract.Ens(&_ReverseRegistrar.CallOpts)
}

// Ens is a free data retrieval call binding the contract method 0x3f15457f.
//
// Solidity: function ens() view returns(address)
func (_ReverseRegistrar *ReverseRegistrarCallerSession) Ens() (common.Address, error) {
	return _ReverseRegistrar.Contract.Ens(&_ReverseRegistrar.CallOpts)
}

// Node is a free data retrieval call binding the contract method 0xbffbe61c.
//
// Solidity: function node(address addr) pure returns(bytes32)
func (_ReverseRegistrar *ReverseRegistrarCaller) Node(opts *bind.CallOpts, addr common.Address) ([32]byte, error) {
	var out []interface{}
	err := _ReverseRegistrar.contract.Call(opts, &out, "node", addr)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Node is a free data retrieval call binding the contract method 0xbffbe61c.
//
// Solidity: function node(address addr) pure returns(bytes32)
func (_ReverseRegistrar *ReverseRegistrarSession) Node(addr common.Address) ([32]byte, error) {
	return _ReverseRegistrar.Contract.Node(&_ReverseRegistrar.CallOpts, addr)
}

// Node is a free data retrieval call binding the contract method 0xbffbe61c.
//
// Solidity: function node(address addr) pure returns(bytes32)
func (_ReverseRegistrar *ReverseRegistrarCallerSession) Node(addr common.Address) ([32]byte, error) {
	return _ReverseRegistrar.Contract.Node(&_ReverseRegistrar.CallOpts, addr)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ReverseRegistrar *ReverseRegistrarCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ReverseRegistrar.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ReverseRegistrar *ReverseRegistrarSession) Owner() (common.Address, error) {
	return _ReverseRegistrar.Contract.Owner(&_ReverseRegistrar.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ReverseRegistrar *ReverseRegistrarCallerSession) Owner() (common.Address, error) {
	return _ReverseRegistrar.Contract.Owner(&_ReverseRegistrar.CallOpts)
}

// Claim is a paid mutator transaction binding the contract method 0x1e83409a.
//
// Solidity: function claim(address owner) returns(bytes32)
func (_ReverseRegistrar *ReverseRegistrarTransactor) Claim(opts *bind.TransactOpts, owner common.Address) (*types.Transaction, error) {
	return _ReverseRegistrar.contract.Transact(opts, "claim", owner)
}

// Claim is a paid mutator transaction binding the contract method 0x1e83409a.
//
// Solidity: function claim(address owner) returns(bytes32)
func (_ReverseRegistrar *ReverseRegistrarSession) Claim(owner common.Address) (*types.Transaction, error) {
	return _ReverseRegistrar.Contract.Claim(&_ReverseRegistrar.TransactOpts, owner)
}

// Claim is a paid mutator transaction binding the contract method 0x1e83409a.
//
// Solidity: function claim(address owner) returns(bytes32)
func (_ReverseRegistrar *ReverseRegistrarTransactorSession) Claim(owner common.Address) (*types.Transaction, error) {
	return _ReverseRegistrar.Contract.Claim(&_ReverseRegistrar.TransactOpts, owner)
}

// ClaimForAddr is a paid mutator transaction binding the contract method 0x65669631.
//
// Solidity: function claimForAddr(address addr, address owner, address resolver) returns(bytes32)
func (_ReverseRegistrar *ReverseRegistrarTransactor) ClaimForAddr(opts *bind.TransactOpts, addr common.Address, owner common.Address, resolver common.Address) (*types.Transaction, error) {
	return _ReverseRegistrar.contract.Transact(opts, "claimForAddr", addr, owner, resolver)
}

// ClaimForAddr is a paid mutator transaction binding the contract method 0x65669631.
//
// Solidity: function claimForAddr(address addr, address owner, address resolver) returns(bytes32)
func (_ReverseRegistrar *ReverseRegistrarSession) ClaimForAddr(addr common.Address, owner common.Address, resolver common.Address) (*types.Transaction, error) {
	return _ReverseRegistrar.Contract.ClaimForAddr(&_ReverseRegistrar.TransactOpts, addr, owner, resolver)
}

// ClaimForAddr is a paid mutator transaction binding the contract method 0x65669631.
//
// Solidity: function claimForAddr(address addr, address owner, address resolver) returns(bytes32)
func (_ReverseRegistrar *ReverseRegistrarTransactorSession) ClaimForAddr(addr common.Address, owner common.Address, resolver common.Address) (*types.Transaction, error) {
	return _ReverseRegistrar.Contract.ClaimForAddr(&_ReverseRegistrar.TransactOpts, addr, owner, resolver)
}

// ClaimWithResolver is a paid mutator transaction binding the contract method 0x0f5a5466.
//
// Solidity: function claimWithResolver(address owner, address resolver) returns(bytes32)
func (_ReverseRegistrar *ReverseRegistrarTransactor) ClaimWithResolver(opts *bind.TransactOpts, owner common.Address, resolver common.Address) (*types.Transaction, error) {
	return _ReverseRegistrar.contract.Transact(opts, "claimWithResolver", owner, resolver)
}

// ClaimWithResolver is a paid mutator transaction binding the contract method 0x0f5a5466.
//
// Solidity: function claimWithResolver(address owner, address resolver) returns(bytes32)
func (_ReverseRegistrar *ReverseRegistrarSession) ClaimWithResolver(owner common.Address, resolver common.Address) (*types.Transaction, error) {
	return _ReverseRegistrar.Contract.ClaimWithResolver(&_ReverseRegistrar.TransactOpts, owner, resolver)
}

// ClaimWithResolver is a paid mutator transaction binding the contract method 0x0f5a5466.
//
// Solidity: function claimWithResolver(address owner, address resolver) returns(bytes32)
func (_ReverseRegistrar *ReverseRegistrarTransactorSession) ClaimWithResolver(owner common.Address, resolver common.Address) (*types.Transaction, error) {
	return _ReverseRegistrar.Contract.ClaimWithResolver(&_ReverseRegistrar.TransactOpts, owner, resolver)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ReverseRegistrar *ReverseRegistrarTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReverseRegistrar.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ReverseRegistrar *ReverseRegistrarSession) RenounceOwnership() (*types.Transaction, error) {
	return _ReverseRegistrar.Contract.RenounceOwnership(&_ReverseRegistrar.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ReverseRegistrar *ReverseRegistrarTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ReverseRegistrar.Contract.RenounceOwnership(&_ReverseRegistrar.TransactOpts)
}

// SetController is a paid mutator transaction binding the contract method 0xe0dba60f.
//
// Solidity: function setController(address controller, bool enabled) returns()
func (_ReverseRegistrar *ReverseRegistrarTransactor) SetController(opts *bind.TransactOpts, controller common.Address, enabled bool) (*types.Transaction, error) {
	return _ReverseRegistrar.contract.Transact(opts, "setController", controller, enabled)
}

// SetController is a paid mutator transaction binding the contract method 0xe0dba60f.
//
// Solidity: function setController(address controller, bool enabled) returns()
func (_ReverseRegistrar *ReverseRegistrarSession) SetController(controller common.Address, enabled bool) (*types.Transaction, error) {
	return _ReverseRegistrar.Contract.SetController(&_ReverseRegistrar.TransactOpts, controller, enabled)
}

// SetController is a paid mutator transaction binding the contract method 0xe0dba60f.
//
// Solidity: function setController(address controller, bool enabled) returns()
func (_ReverseRegistrar *ReverseRegistrarTransactorSession) SetController(controller common.Address, enabled bool) (*types.Transaction, error) {
	return _ReverseRegistrar.Contract.SetController(&_ReverseRegistrar.TransactOpts, controller, enabled)
}

// SetDefaultResolver is a paid mutator transaction binding the contract method 0xc66485b2.
//
// Solidity: function setDefaultResolver(address resolver) returns()
func (_ReverseRegistrar *ReverseRegistrarTransactor) SetDefaultResolver(opts *bind.TransactOpts, resolver common.Address) (*types.Transaction, error) {
	return _ReverseRegistrar.contract.Transact(opts, "setDefaultResolver", resolver)
}

// SetDefaultResolver is a paid mutator transaction binding the contract method 0xc66485b2.
//
// Solidity: function setDefaultResolver(address resolver) returns()
func (_ReverseRegistrar *ReverseRegistrarSession) SetDefaultResolver(resolver common.Address) (*types.Transaction, error) {
	return _ReverseRegistrar.Contract.SetDefaultResolver(&_ReverseRegistrar.TransactOpts, resolver)
}

// SetDefaultResolver is a paid mutator transaction binding the contract method 0xc66485b2.
//
// Solidity: function setDefaultResolver(address resolver) returns()
func (_ReverseRegistrar *ReverseRegistrarTransactorSession) SetDefaultResolver(resolver common.Address) (*types.Transaction, error) {
	return _ReverseRegistrar.Contract.SetDefaultResolver(&_ReverseRegistrar.TransactOpts, resolver)
}

// SetName is a paid mutator transaction binding the contract method 0xc47f0027.
//
// Solidity: function setName(string name) returns(bytes32)
func (_ReverseRegistrar *ReverseRegistrarTransactor) SetName(opts *bind.TransactOpts, name string) (*types.Transaction, error) {
	return _ReverseRegistrar.contract.Transact(opts, "setName", name)
}

// SetName is a paid mutator transaction binding the contract method 0xc47f0027.
//
// Solidity: function setName(string name) returns(bytes32)
func (_ReverseRegistrar *ReverseRegistrarSession) SetName(name string) (*types.Transaction, error) {
	return _ReverseRegistrar.Contract.SetName(&_ReverseRegistrar.TransactOpts, name)
}

// SetName is a paid mutator transaction binding the contract method 0xc47f0027.
//
// Solidity: function setName(string name) returns(bytes32)
func (_ReverseRegistrar *ReverseRegistrarTransactorSession) SetName(name string) (*types.Transaction, error) {
	return _ReverseRegistrar.Contract.SetName(&_ReverseRegistrar.TransactOpts, name)
}

// SetNameForAddr is a paid mutator transaction binding the contract method 0x7a806d6b.
//
// Solidity: function setNameForAddr(address addr, address owner, address resolver, string name) returns(bytes32)
func (_ReverseRegistrar *ReverseRegistrarTransactor) SetNameForAddr(opts *bind.TransactOpts, addr common.Address, owner common.Address, resolver common.Address, name string) (*types.Transaction, error) {
	return _ReverseRegistrar.contract.Transact(opts, "setNameForAddr", addr, owner, resolver, name)
}

// SetNameForAddr is a paid mutator transaction binding the contract method 0x7a806d6b.
//
// Solidity: function setNameForAddr(address addr, address owner, address resolver, string name) returns(bytes32)
func (_ReverseRegistrar *ReverseRegistrarSession) SetNameForAddr(addr common.Address, owner common.Address, resolver common.Address, name string) (*types.Transaction, error) {
	return _ReverseRegistrar.Contract.SetNameForAddr(&_ReverseRegistrar.TransactOpts, addr, owner, resolver, name)
}

// SetNameForAddr is a paid mutator transaction binding the contract method 0x7a806d6b.
//
// Solidity: function setNameForAddr(address addr, address owner, address resolver, string name) returns(bytes32)
func (_ReverseRegistrar *ReverseRegistrarTransactorSession) SetNameForAddr(addr common.Address, owner common.Address, resolver common.Address, name string) (*types.Transaction, error) {
	return _ReverseRegistrar.Contract.SetNameForAddr(&_ReverseRegistrar.TransactOpts, addr, owner, resolver, name)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ReverseRegistrar *ReverseRegistrarTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ReverseRegistrar.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ReverseRegistrar *ReverseRegistrarSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ReverseRegistrar.Contract.TransferOwnership(&_ReverseRegistrar.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ReverseRegistrar *ReverseRegistrarTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ReverseRegistrar.Contract.TransferOwnership(&_ReverseRegistrar.TransactOpts, newOwner)
}

// ReverseRegistrarControllerChangedIterator is returned from FilterControllerChanged and is used to iterate over the raw logs and unpacked data for ControllerChanged events raised by the ReverseRegistrar contract.
type ReverseRegistrarControllerChangedIterator struct {
	Event *ReverseRegistrarControllerChanged // Event containing the contract specifics and raw log

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
func (it *ReverseRegistrarControllerChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReverseRegistrarControllerChanged)
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
		it.Event = new(ReverseRegistrarControllerChanged)
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
func (it *ReverseRegistrarControllerChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReverseRegistrarControllerChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReverseRegistrarControllerChanged represents a ControllerChanged event raised by the ReverseRegistrar contract.
type ReverseRegistrarControllerChanged struct {
	Controller common.Address
	Enabled    bool
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterControllerChanged is a free log retrieval operation binding the contract event 0x4c97694570a07277810af7e5669ffd5f6a2d6b74b6e9a274b8b870fd5114cf87.
//
// Solidity: event ControllerChanged(address indexed controller, bool enabled)
func (_ReverseRegistrar *ReverseRegistrarFilterer) FilterControllerChanged(opts *bind.FilterOpts, controller []common.Address) (*ReverseRegistrarControllerChangedIterator, error) {

	var controllerRule []interface{}
	for _, controllerItem := range controller {
		controllerRule = append(controllerRule, controllerItem)
	}

	logs, sub, err := _ReverseRegistrar.contract.FilterLogs(opts, "ControllerChanged", controllerRule)
	if err != nil {
		return nil, err
	}
	return &ReverseRegistrarControllerChangedIterator{contract: _ReverseRegistrar.contract, event: "ControllerChanged", logs: logs, sub: sub}, nil
}

// WatchControllerChanged is a free log subscription operation binding the contract event 0x4c97694570a07277810af7e5669ffd5f6a2d6b74b6e9a274b8b870fd5114cf87.
//
// Solidity: event ControllerChanged(address indexed controller, bool enabled)
func (_ReverseRegistrar *ReverseRegistrarFilterer) WatchControllerChanged(opts *bind.WatchOpts, sink chan<- *ReverseRegistrarControllerChanged, controller []common.Address) (event.Subscription, error) {

	var controllerRule []interface{}
	for _, controllerItem := range controller {
		controllerRule = append(controllerRule, controllerItem)
	}

	logs, sub, err := _ReverseRegistrar.contract.WatchLogs(opts, "ControllerChanged", controllerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReverseRegistrarControllerChanged)
				if err := _ReverseRegistrar.contract.UnpackLog(event, "ControllerChanged", log); err != nil {
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
func (_ReverseRegistrar *ReverseRegistrarFilterer) ParseControllerChanged(log types.Log) (*ReverseRegistrarControllerChanged, error) {
	event := new(ReverseRegistrarControllerChanged)
	if err := _ReverseRegistrar.contract.UnpackLog(event, "ControllerChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReverseRegistrarDefaultResolverChangedIterator is returned from FilterDefaultResolverChanged and is used to iterate over the raw logs and unpacked data for DefaultResolverChanged events raised by the ReverseRegistrar contract.
type ReverseRegistrarDefaultResolverChangedIterator struct {
	Event *ReverseRegistrarDefaultResolverChanged // Event containing the contract specifics and raw log

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
func (it *ReverseRegistrarDefaultResolverChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReverseRegistrarDefaultResolverChanged)
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
		it.Event = new(ReverseRegistrarDefaultResolverChanged)
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
func (it *ReverseRegistrarDefaultResolverChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReverseRegistrarDefaultResolverChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReverseRegistrarDefaultResolverChanged represents a DefaultResolverChanged event raised by the ReverseRegistrar contract.
type ReverseRegistrarDefaultResolverChanged struct {
	Resolver common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterDefaultResolverChanged is a free log retrieval operation binding the contract event 0xeae17a84d9eb83d8c8eb317f9e7d64857bc363fa51674d996c023f4340c577cf.
//
// Solidity: event DefaultResolverChanged(address indexed resolver)
func (_ReverseRegistrar *ReverseRegistrarFilterer) FilterDefaultResolverChanged(opts *bind.FilterOpts, resolver []common.Address) (*ReverseRegistrarDefaultResolverChangedIterator, error) {

	var resolverRule []interface{}
	for _, resolverItem := range resolver {
		resolverRule = append(resolverRule, resolverItem)
	}

	logs, sub, err := _ReverseRegistrar.contract.FilterLogs(opts, "DefaultResolverChanged", resolverRule)
	if err != nil {
		return nil, err
	}
	return &ReverseRegistrarDefaultResolverChangedIterator{contract: _ReverseRegistrar.contract, event: "DefaultResolverChanged", logs: logs, sub: sub}, nil
}

// WatchDefaultResolverChanged is a free log subscription operation binding the contract event 0xeae17a84d9eb83d8c8eb317f9e7d64857bc363fa51674d996c023f4340c577cf.
//
// Solidity: event DefaultResolverChanged(address indexed resolver)
func (_ReverseRegistrar *ReverseRegistrarFilterer) WatchDefaultResolverChanged(opts *bind.WatchOpts, sink chan<- *ReverseRegistrarDefaultResolverChanged, resolver []common.Address) (event.Subscription, error) {

	var resolverRule []interface{}
	for _, resolverItem := range resolver {
		resolverRule = append(resolverRule, resolverItem)
	}

	logs, sub, err := _ReverseRegistrar.contract.WatchLogs(opts, "DefaultResolverChanged", resolverRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReverseRegistrarDefaultResolverChanged)
				if err := _ReverseRegistrar.contract.UnpackLog(event, "DefaultResolverChanged", log); err != nil {
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

// ParseDefaultResolverChanged is a log parse operation binding the contract event 0xeae17a84d9eb83d8c8eb317f9e7d64857bc363fa51674d996c023f4340c577cf.
//
// Solidity: event DefaultResolverChanged(address indexed resolver)
func (_ReverseRegistrar *ReverseRegistrarFilterer) ParseDefaultResolverChanged(log types.Log) (*ReverseRegistrarDefaultResolverChanged, error) {
	event := new(ReverseRegistrarDefaultResolverChanged)
	if err := _ReverseRegistrar.contract.UnpackLog(event, "DefaultResolverChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReverseRegistrarOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ReverseRegistrar contract.
type ReverseRegistrarOwnershipTransferredIterator struct {
	Event *ReverseRegistrarOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ReverseRegistrarOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReverseRegistrarOwnershipTransferred)
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
		it.Event = new(ReverseRegistrarOwnershipTransferred)
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
func (it *ReverseRegistrarOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReverseRegistrarOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReverseRegistrarOwnershipTransferred represents a OwnershipTransferred event raised by the ReverseRegistrar contract.
type ReverseRegistrarOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ReverseRegistrar *ReverseRegistrarFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ReverseRegistrarOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ReverseRegistrar.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ReverseRegistrarOwnershipTransferredIterator{contract: _ReverseRegistrar.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ReverseRegistrar *ReverseRegistrarFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ReverseRegistrarOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ReverseRegistrar.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReverseRegistrarOwnershipTransferred)
				if err := _ReverseRegistrar.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ReverseRegistrar *ReverseRegistrarFilterer) ParseOwnershipTransferred(log types.Log) (*ReverseRegistrarOwnershipTransferred, error) {
	event := new(ReverseRegistrarOwnershipTransferred)
	if err := _ReverseRegistrar.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReverseRegistrarReverseClaimedIterator is returned from FilterReverseClaimed and is used to iterate over the raw logs and unpacked data for ReverseClaimed events raised by the ReverseRegistrar contract.
type ReverseRegistrarReverseClaimedIterator struct {
	Event *ReverseRegistrarReverseClaimed // Event containing the contract specifics and raw log

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
func (it *ReverseRegistrarReverseClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReverseRegistrarReverseClaimed)
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
		it.Event = new(ReverseRegistrarReverseClaimed)
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
func (it *ReverseRegistrarReverseClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReverseRegistrarReverseClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReverseRegistrarReverseClaimed represents a ReverseClaimed event raised by the ReverseRegistrar contract.
type ReverseRegistrarReverseClaimed struct {
	Addr common.Address
	Node [32]byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterReverseClaimed is a free log retrieval operation binding the contract event 0x6ada868dd3058cf77a48a74489fd7963688e5464b2b0fa957ace976243270e92.
//
// Solidity: event ReverseClaimed(address indexed addr, bytes32 indexed node)
func (_ReverseRegistrar *ReverseRegistrarFilterer) FilterReverseClaimed(opts *bind.FilterOpts, addr []common.Address, node [][32]byte) (*ReverseRegistrarReverseClaimedIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _ReverseRegistrar.contract.FilterLogs(opts, "ReverseClaimed", addrRule, nodeRule)
	if err != nil {
		return nil, err
	}
	return &ReverseRegistrarReverseClaimedIterator{contract: _ReverseRegistrar.contract, event: "ReverseClaimed", logs: logs, sub: sub}, nil
}

// WatchReverseClaimed is a free log subscription operation binding the contract event 0x6ada868dd3058cf77a48a74489fd7963688e5464b2b0fa957ace976243270e92.
//
// Solidity: event ReverseClaimed(address indexed addr, bytes32 indexed node)
func (_ReverseRegistrar *ReverseRegistrarFilterer) WatchReverseClaimed(opts *bind.WatchOpts, sink chan<- *ReverseRegistrarReverseClaimed, addr []common.Address, node [][32]byte) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _ReverseRegistrar.contract.WatchLogs(opts, "ReverseClaimed", addrRule, nodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReverseRegistrarReverseClaimed)
				if err := _ReverseRegistrar.contract.UnpackLog(event, "ReverseClaimed", log); err != nil {
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

// ParseReverseClaimed is a log parse operation binding the contract event 0x6ada868dd3058cf77a48a74489fd7963688e5464b2b0fa957ace976243270e92.
//
// Solidity: event ReverseClaimed(address indexed addr, bytes32 indexed node)
func (_ReverseRegistrar *ReverseRegistrarFilterer) ParseReverseClaimed(log types.Log) (*ReverseRegistrarReverseClaimed, error) {
	event := new(ReverseRegistrarReverseClaimed)
	if err := _ReverseRegistrar.contract.UnpackLog(event, "ReverseClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
