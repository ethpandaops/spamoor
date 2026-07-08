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

// PackedUserOperation is defined in the EntryPoint binding within this same
// package; the duplicate generated here is removed by compile.sh to avoid a
// redeclaration. See the dedupe step in contract/compile.sh.

// AcceptAllPaymasterMetaData contains all meta data concerning the AcceptAllPaymaster contract.
var AcceptAllPaymasterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_entryPoint\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"entryPoint\",\"outputs\":[{\"internalType\":\"contractIStakeManager\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDeposit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"postOp\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"accountGasLimits\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"preVerificationGas\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"gasFees\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"paymasterAndData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structPackedUserOperation\",\"name\":\"\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"validatePaymasterUserOp\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"context\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"validationData\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60a060405234801561000f575f80fd5b506040516105cf3803806105cf83398101604081905261002e9161003f565b6001600160a01b031660805261006c565b5f6020828403121561004f575f80fd5b81516001600160a01b0381168114610065575f80fd5b9392505050565b6080516105296100a65f395f818160660152818161012f015281816101a101528181610237015281816102c2015261034d01526105295ff3fe60806040526004361061004c575f3560e01c806352b7512c146100c95780637c627b21146100ff578063b0d691fe1461011e578063c399ec8814610169578063d0e30db01461018b575f80fd5b366100c55760405163b760faf960e01b81523060048201527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063b760faf99034906024015f604051808303818588803b1580156100b1575f80fd5b505af11580156100c3573d5f803e3d5ffd5b005b5f80fd5b3480156100d4575f80fd5b506100e86100e33660046103aa565b610193565b6040516100f69291906103f9565b60405180910390f35b34801561010a575f80fd5b506100c361011936600461044b565b61022c565b348015610129575f80fd5b506101517f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b0390911681526020016100f6565b348015610174575f80fd5b5061017d6102ab565b6040519081526020016100f6565b6100c3610338565b60605f336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146102135760405162461bcd60e51b815260206004820152601e60248201527f5061796d61737465723a206e6f742066726f6d20456e747279506f696e74000060448201526064015b60405180910390fd5b505060408051602081019091525f808252935093915050565b336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146102a45760405162461bcd60e51b815260206004820152601e60248201527f5061796d61737465723a206e6f742066726f6d20456e747279506f696e740000604482015260640161020a565b5050505050565b6040516370a0823160e01b81523060048201525f907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316906370a0823190602401602060405180830381865afa15801561030f573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061033391906104dc565b905090565b60405163b760faf960e01b81523060048201527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063b760faf99034906024015f604051808303818588803b158015610398575f80fd5b505af11580156102a4573d5f803e3d5ffd5b5f805f606084860312156103bc575f80fd5b833567ffffffffffffffff8111156103d2575f80fd5b840161012081870312156103e4575f80fd5b95602085013595506040909401359392505050565b604081525f83518060408401525f5b818110156104255760208187018101516060868401015201610408565b505f606082850101526060601f19601f8301168401019150508260208301529392505050565b5f805f805f6080868803121561045f575f80fd5b853560ff8116811461046f575f80fd5b9450602086013567ffffffffffffffff8082111561048b575f80fd5b818801915088601f83011261049e575f80fd5b8135818111156104ac575f80fd5b8960208285010111156104bd575f80fd5b9699602092909201985095966040810135965060600135945092505050565b5f602082840312156104ec575f80fd5b505191905056fea2646970667358221220f7a71d13dd3b4a2870ef1137e4a0b6a3553ff3d03ff05024919d4d0b5e87a37164736f6c63430008170033",
}

// AcceptAllPaymasterABI is the input ABI used to generate the binding from.
// Deprecated: Use AcceptAllPaymasterMetaData.ABI instead.
var AcceptAllPaymasterABI = AcceptAllPaymasterMetaData.ABI

// AcceptAllPaymasterBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use AcceptAllPaymasterMetaData.Bin instead.
var AcceptAllPaymasterBin = AcceptAllPaymasterMetaData.Bin

// DeployAcceptAllPaymaster deploys a new Ethereum contract, binding an instance of AcceptAllPaymaster to it.
func DeployAcceptAllPaymaster(auth *bind.TransactOpts, backend bind.ContractBackend, _entryPoint common.Address) (common.Address, *types.Transaction, *AcceptAllPaymaster, error) {
	parsed, err := AcceptAllPaymasterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AcceptAllPaymasterBin), backend, _entryPoint)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AcceptAllPaymaster{AcceptAllPaymasterCaller: AcceptAllPaymasterCaller{contract: contract}, AcceptAllPaymasterTransactor: AcceptAllPaymasterTransactor{contract: contract}, AcceptAllPaymasterFilterer: AcceptAllPaymasterFilterer{contract: contract}}, nil
}

// AcceptAllPaymaster is an auto generated Go binding around an Ethereum contract.
type AcceptAllPaymaster struct {
	AcceptAllPaymasterCaller     // Read-only binding to the contract
	AcceptAllPaymasterTransactor // Write-only binding to the contract
	AcceptAllPaymasterFilterer   // Log filterer for contract events
}

// AcceptAllPaymasterCaller is an auto generated read-only Go binding around an Ethereum contract.
type AcceptAllPaymasterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AcceptAllPaymasterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AcceptAllPaymasterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AcceptAllPaymasterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AcceptAllPaymasterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AcceptAllPaymasterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AcceptAllPaymasterSession struct {
	Contract     *AcceptAllPaymaster // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// AcceptAllPaymasterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AcceptAllPaymasterCallerSession struct {
	Contract *AcceptAllPaymasterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// AcceptAllPaymasterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AcceptAllPaymasterTransactorSession struct {
	Contract     *AcceptAllPaymasterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// AcceptAllPaymasterRaw is an auto generated low-level Go binding around an Ethereum contract.
type AcceptAllPaymasterRaw struct {
	Contract *AcceptAllPaymaster // Generic contract binding to access the raw methods on
}

// AcceptAllPaymasterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AcceptAllPaymasterCallerRaw struct {
	Contract *AcceptAllPaymasterCaller // Generic read-only contract binding to access the raw methods on
}

// AcceptAllPaymasterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AcceptAllPaymasterTransactorRaw struct {
	Contract *AcceptAllPaymasterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAcceptAllPaymaster creates a new instance of AcceptAllPaymaster, bound to a specific deployed contract.
func NewAcceptAllPaymaster(address common.Address, backend bind.ContractBackend) (*AcceptAllPaymaster, error) {
	contract, err := bindAcceptAllPaymaster(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AcceptAllPaymaster{AcceptAllPaymasterCaller: AcceptAllPaymasterCaller{contract: contract}, AcceptAllPaymasterTransactor: AcceptAllPaymasterTransactor{contract: contract}, AcceptAllPaymasterFilterer: AcceptAllPaymasterFilterer{contract: contract}}, nil
}

// NewAcceptAllPaymasterCaller creates a new read-only instance of AcceptAllPaymaster, bound to a specific deployed contract.
func NewAcceptAllPaymasterCaller(address common.Address, caller bind.ContractCaller) (*AcceptAllPaymasterCaller, error) {
	contract, err := bindAcceptAllPaymaster(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AcceptAllPaymasterCaller{contract: contract}, nil
}

// NewAcceptAllPaymasterTransactor creates a new write-only instance of AcceptAllPaymaster, bound to a specific deployed contract.
func NewAcceptAllPaymasterTransactor(address common.Address, transactor bind.ContractTransactor) (*AcceptAllPaymasterTransactor, error) {
	contract, err := bindAcceptAllPaymaster(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AcceptAllPaymasterTransactor{contract: contract}, nil
}

// NewAcceptAllPaymasterFilterer creates a new log filterer instance of AcceptAllPaymaster, bound to a specific deployed contract.
func NewAcceptAllPaymasterFilterer(address common.Address, filterer bind.ContractFilterer) (*AcceptAllPaymasterFilterer, error) {
	contract, err := bindAcceptAllPaymaster(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AcceptAllPaymasterFilterer{contract: contract}, nil
}

// bindAcceptAllPaymaster binds a generic wrapper to an already deployed contract.
func bindAcceptAllPaymaster(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AcceptAllPaymasterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AcceptAllPaymaster *AcceptAllPaymasterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AcceptAllPaymaster.Contract.AcceptAllPaymasterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AcceptAllPaymaster *AcceptAllPaymasterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AcceptAllPaymaster.Contract.AcceptAllPaymasterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AcceptAllPaymaster *AcceptAllPaymasterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AcceptAllPaymaster.Contract.AcceptAllPaymasterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AcceptAllPaymaster *AcceptAllPaymasterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AcceptAllPaymaster.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AcceptAllPaymaster *AcceptAllPaymasterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AcceptAllPaymaster.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AcceptAllPaymaster *AcceptAllPaymasterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AcceptAllPaymaster.Contract.contract.Transact(opts, method, params...)
}

// EntryPoint is a free data retrieval call binding the contract method 0xb0d691fe.
//
// Solidity: function entryPoint() view returns(address)
func (_AcceptAllPaymaster *AcceptAllPaymasterCaller) EntryPoint(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AcceptAllPaymaster.contract.Call(opts, &out, "entryPoint")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EntryPoint is a free data retrieval call binding the contract method 0xb0d691fe.
//
// Solidity: function entryPoint() view returns(address)
func (_AcceptAllPaymaster *AcceptAllPaymasterSession) EntryPoint() (common.Address, error) {
	return _AcceptAllPaymaster.Contract.EntryPoint(&_AcceptAllPaymaster.CallOpts)
}

// EntryPoint is a free data retrieval call binding the contract method 0xb0d691fe.
//
// Solidity: function entryPoint() view returns(address)
func (_AcceptAllPaymaster *AcceptAllPaymasterCallerSession) EntryPoint() (common.Address, error) {
	return _AcceptAllPaymaster.Contract.EntryPoint(&_AcceptAllPaymaster.CallOpts)
}

// GetDeposit is a free data retrieval call binding the contract method 0xc399ec88.
//
// Solidity: function getDeposit() view returns(uint256)
func (_AcceptAllPaymaster *AcceptAllPaymasterCaller) GetDeposit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AcceptAllPaymaster.contract.Call(opts, &out, "getDeposit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetDeposit is a free data retrieval call binding the contract method 0xc399ec88.
//
// Solidity: function getDeposit() view returns(uint256)
func (_AcceptAllPaymaster *AcceptAllPaymasterSession) GetDeposit() (*big.Int, error) {
	return _AcceptAllPaymaster.Contract.GetDeposit(&_AcceptAllPaymaster.CallOpts)
}

// GetDeposit is a free data retrieval call binding the contract method 0xc399ec88.
//
// Solidity: function getDeposit() view returns(uint256)
func (_AcceptAllPaymaster *AcceptAllPaymasterCallerSession) GetDeposit() (*big.Int, error) {
	return _AcceptAllPaymaster.Contract.GetDeposit(&_AcceptAllPaymaster.CallOpts)
}

// PostOp is a free data retrieval call binding the contract method 0x7c627b21.
//
// Solidity: function postOp(uint8 , bytes , uint256 , uint256 ) view returns()
func (_AcceptAllPaymaster *AcceptAllPaymasterCaller) PostOp(opts *bind.CallOpts, arg0 uint8, arg1 []byte, arg2 *big.Int, arg3 *big.Int) error {
	var out []interface{}
	err := _AcceptAllPaymaster.contract.Call(opts, &out, "postOp", arg0, arg1, arg2, arg3)

	if err != nil {
		return err
	}

	return err

}

// PostOp is a free data retrieval call binding the contract method 0x7c627b21.
//
// Solidity: function postOp(uint8 , bytes , uint256 , uint256 ) view returns()
func (_AcceptAllPaymaster *AcceptAllPaymasterSession) PostOp(arg0 uint8, arg1 []byte, arg2 *big.Int, arg3 *big.Int) error {
	return _AcceptAllPaymaster.Contract.PostOp(&_AcceptAllPaymaster.CallOpts, arg0, arg1, arg2, arg3)
}

// PostOp is a free data retrieval call binding the contract method 0x7c627b21.
//
// Solidity: function postOp(uint8 , bytes , uint256 , uint256 ) view returns()
func (_AcceptAllPaymaster *AcceptAllPaymasterCallerSession) PostOp(arg0 uint8, arg1 []byte, arg2 *big.Int, arg3 *big.Int) error {
	return _AcceptAllPaymaster.Contract.PostOp(&_AcceptAllPaymaster.CallOpts, arg0, arg1, arg2, arg3)
}

// ValidatePaymasterUserOp is a free data retrieval call binding the contract method 0x52b7512c.
//
// Solidity: function validatePaymasterUserOp((address,uint256,bytes,bytes,bytes32,uint256,bytes32,bytes,bytes) , bytes32 , uint256 ) view returns(bytes context, uint256 validationData)
func (_AcceptAllPaymaster *AcceptAllPaymasterCaller) ValidatePaymasterUserOp(opts *bind.CallOpts, arg0 PackedUserOperation, arg1 [32]byte, arg2 *big.Int) (struct {
	Context        []byte
	ValidationData *big.Int
}, error) {
	var out []interface{}
	err := _AcceptAllPaymaster.contract.Call(opts, &out, "validatePaymasterUserOp", arg0, arg1, arg2)

	outstruct := new(struct {
		Context        []byte
		ValidationData *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Context = *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	outstruct.ValidationData = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// ValidatePaymasterUserOp is a free data retrieval call binding the contract method 0x52b7512c.
//
// Solidity: function validatePaymasterUserOp((address,uint256,bytes,bytes,bytes32,uint256,bytes32,bytes,bytes) , bytes32 , uint256 ) view returns(bytes context, uint256 validationData)
func (_AcceptAllPaymaster *AcceptAllPaymasterSession) ValidatePaymasterUserOp(arg0 PackedUserOperation, arg1 [32]byte, arg2 *big.Int) (struct {
	Context        []byte
	ValidationData *big.Int
}, error) {
	return _AcceptAllPaymaster.Contract.ValidatePaymasterUserOp(&_AcceptAllPaymaster.CallOpts, arg0, arg1, arg2)
}

// ValidatePaymasterUserOp is a free data retrieval call binding the contract method 0x52b7512c.
//
// Solidity: function validatePaymasterUserOp((address,uint256,bytes,bytes,bytes32,uint256,bytes32,bytes,bytes) , bytes32 , uint256 ) view returns(bytes context, uint256 validationData)
func (_AcceptAllPaymaster *AcceptAllPaymasterCallerSession) ValidatePaymasterUserOp(arg0 PackedUserOperation, arg1 [32]byte, arg2 *big.Int) (struct {
	Context        []byte
	ValidationData *big.Int
}, error) {
	return _AcceptAllPaymaster.Contract.ValidatePaymasterUserOp(&_AcceptAllPaymaster.CallOpts, arg0, arg1, arg2)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_AcceptAllPaymaster *AcceptAllPaymasterTransactor) Deposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AcceptAllPaymaster.contract.Transact(opts, "deposit")
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_AcceptAllPaymaster *AcceptAllPaymasterSession) Deposit() (*types.Transaction, error) {
	return _AcceptAllPaymaster.Contract.Deposit(&_AcceptAllPaymaster.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_AcceptAllPaymaster *AcceptAllPaymasterTransactorSession) Deposit() (*types.Transaction, error) {
	return _AcceptAllPaymaster.Contract.Deposit(&_AcceptAllPaymaster.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_AcceptAllPaymaster *AcceptAllPaymasterTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AcceptAllPaymaster.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_AcceptAllPaymaster *AcceptAllPaymasterSession) Receive() (*types.Transaction, error) {
	return _AcceptAllPaymaster.Contract.Receive(&_AcceptAllPaymaster.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_AcceptAllPaymaster *AcceptAllPaymasterTransactorSession) Receive() (*types.Transaction, error) {
	return _AcceptAllPaymaster.Contract.Receive(&_AcceptAllPaymaster.TransactOpts)
}
