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

// IPriceOraclePrice is defined in the ETHRegistrarController binding within
// this same package; the duplicate generated here is removed by compile.sh to
// avoid a redeclaration.

// StablePriceOracleMetaData contains all meta data concerning the StablePriceOracle contract.
var StablePriceOracleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractAggregatorInterface\",\"name\":\"_usdOracle\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"_rentPrices\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"prices\",\"type\":\"uint256[]\"}],\"name\":\"RentPriceChanged\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"expires\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"duration\",\"type\":\"uint256\"}],\"name\":\"premium\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"expires\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"duration\",\"type\":\"uint256\"}],\"name\":\"price\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"base\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"premium\",\"type\":\"uint256\"}],\"internalType\":\"structIPriceOracle.Price\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"price1Letter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"price2Letter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"price3Letter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"price4Letter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"price5Letter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceID\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"usdOracle\",\"outputs\":[{\"internalType\":\"contractAggregatorInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x61014060405234801561001157600080fd5b50604051610c31380380610c3183398101604081905261003091610106565b6001600160a01b0382166101205280518190600090610051576100516101f4565b60200260200101516080818152505080600181518110610073576100736101f4565b602002602001015160a0818152505080600281518110610095576100956101f4565b602002602001015160c08181525050806003815181106100b7576100b76101f4565b602002602001015160e08181525050806004815181106100d9576100d96101f4565b60200260200101516101008181525050505061020a565b634e487b7160e01b600052604160045260246000fd5b6000806040838503121561011957600080fd5b82516001600160a01b038116811461013057600080fd5b60208401519092506001600160401b0381111561014c57600080fd5b8301601f8101851361015d57600080fd5b80516001600160401b03811115610176576101766100f0565b604051600582901b90603f8201601f191681016001600160401b03811182821017156101a4576101a46100f0565b6040529182526020818401810192908101888411156101c257600080fd5b6020850194505b838510156101e5578451808252602095860195909350016101c9565b50809450505050509250929050565b634e487b7160e01b600052603260045260246000fd5b60805160a05160c05160e05161010051610120516109af6102826000396000818161019901526106ea015260008181610138015261032e01526000818161020c015261036701526000818161015f01526103990152600081816101e501526103cb01526000818160d501526103f501526109af6000f3fe608060405234801561001057600080fd5b50600436106100a35760003560e01c8063a200e15311610076578063c8a4271f1161005b578063c8a4271f14610194578063cd5d2c74146101e0578063d820ed421461020757600080fd5b8063a200e1531461015a578063a34e35961461018157600080fd5b806301ffc9a7146100a85780632c0fd74c146100d057806350e9a7151461010557806359b6b86c14610133575b600080fd5b6100bb6100b63660046107a2565b61022e565b60405190151581526020015b60405180910390f35b6100f77f000000000000000000000000000000000000000000000000000000000000000081565b6040519081526020016100c7565b6101186101133660046107e4565b6102c7565b604080518251815260209283015192810192909252016100c7565b6100f77f000000000000000000000000000000000000000000000000000000000000000081565b6100f77f000000000000000000000000000000000000000000000000000000000000000081565b6100f761018f3660046107e4565b61048d565b6101bb7f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100c7565b6100f77f000000000000000000000000000000000000000000000000000000000000000081565b6100f77f000000000000000000000000000000000000000000000000000000000000000081565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a70000000000000000000000000000000000000000000000000000000014806102c157507fffffffff0000000000000000000000000000000000000000000000000000000082167f50e9a71500000000000000000000000000000000000000000000000000000000145b92915050565b6040805180820190915260008082526020820152600061031c86868080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506104de92505050565b905060006005821061035957610352847f0000000000000000000000000000000000000000000000000000000000000000610894565b905061041c565b8160040361038b57610352847f0000000000000000000000000000000000000000000000000000000000000000610894565b816003036103bd57610352847f0000000000000000000000000000000000000000000000000000000000000000610894565b816002036103ef57610352847f0000000000000000000000000000000000000000000000000000000000000000610894565b610419847f0000000000000000000000000000000000000000000000000000000000000000610894565b90505b6040518060400160405280610430836106e5565b815260200161048061047b8a8a8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508c92508b91506107999050565b6106e5565b9052979650505050505050565b60006104d561047b86868080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508892508791506107999050565b95945050505050565b8051600090819081905b808210156106dc576000858381518110610504576105046108ab565b01602001517fff000000000000000000000000000000000000000000000000000000000000001690507f8000000000000000000000000000000000000000000000000000000000000000811015610567576105606001846108da565b92506106c9565b7fe0000000000000000000000000000000000000000000000000000000000000007fff00000000000000000000000000000000000000000000000000000000000000821610156105bc576105606002846108da565b7ff0000000000000000000000000000000000000000000000000000000000000007fff0000000000000000000000000000000000000000000000000000000000000082161015610611576105606003846108da565b7ff8000000000000000000000000000000000000000000000000000000000000007fff0000000000000000000000000000000000000000000000000000000000000082161015610666576105606004846108da565b7ffc000000000000000000000000000000000000000000000000000000000000007fff00000000000000000000000000000000000000000000000000000000000000821610156106bb576105606005846108da565b6106c66006846108da565b92505b50826106d4816108ed565b9350506104e8565b50909392505050565b6000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166350d25bcd6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610753573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107779190610925565b905080610788846305f5e100610894565b610792919061093e565b9392505050565b60009392505050565b6000602082840312156107b457600080fd5b81357fffffffff000000000000000000000000000000000000000000000000000000008116811461079257600080fd5b600080600080606085870312156107fa57600080fd5b843567ffffffffffffffff81111561081157600080fd5b8501601f8101871361082257600080fd5b803567ffffffffffffffff81111561083957600080fd5b87602082840101111561084b57600080fd5b602091820198909750908601359560400135945092505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b80820281158282048414176102c1576102c1610865565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b808201808211156102c1576102c1610865565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361091e5761091e610865565b5060010190565b60006020828403121561093757600080fd5b5051919050565b600082610974577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b50049056fea26469706673582212209e5ca84b8bc64502e2d50a9222cb93c3eb80e1b41c9c3f2b5c1d4736f811a01464736f6c634300081a0033",
}

// StablePriceOracleABI is the input ABI used to generate the binding from.
// Deprecated: Use StablePriceOracleMetaData.ABI instead.
var StablePriceOracleABI = StablePriceOracleMetaData.ABI

// StablePriceOracleBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use StablePriceOracleMetaData.Bin instead.
var StablePriceOracleBin = StablePriceOracleMetaData.Bin

// DeployStablePriceOracle deploys a new Ethereum contract, binding an instance of StablePriceOracle to it.
func DeployStablePriceOracle(auth *bind.TransactOpts, backend bind.ContractBackend, _usdOracle common.Address, _rentPrices []*big.Int) (common.Address, *types.Transaction, *StablePriceOracle, error) {
	parsed, err := StablePriceOracleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StablePriceOracleBin), backend, _usdOracle, _rentPrices)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &StablePriceOracle{StablePriceOracleCaller: StablePriceOracleCaller{contract: contract}, StablePriceOracleTransactor: StablePriceOracleTransactor{contract: contract}, StablePriceOracleFilterer: StablePriceOracleFilterer{contract: contract}}, nil
}

// StablePriceOracle is an auto generated Go binding around an Ethereum contract.
type StablePriceOracle struct {
	StablePriceOracleCaller     // Read-only binding to the contract
	StablePriceOracleTransactor // Write-only binding to the contract
	StablePriceOracleFilterer   // Log filterer for contract events
}

// StablePriceOracleCaller is an auto generated read-only Go binding around an Ethereum contract.
type StablePriceOracleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StablePriceOracleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StablePriceOracleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StablePriceOracleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StablePriceOracleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StablePriceOracleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StablePriceOracleSession struct {
	Contract     *StablePriceOracle // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// StablePriceOracleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StablePriceOracleCallerSession struct {
	Contract *StablePriceOracleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// StablePriceOracleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StablePriceOracleTransactorSession struct {
	Contract     *StablePriceOracleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// StablePriceOracleRaw is an auto generated low-level Go binding around an Ethereum contract.
type StablePriceOracleRaw struct {
	Contract *StablePriceOracle // Generic contract binding to access the raw methods on
}

// StablePriceOracleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StablePriceOracleCallerRaw struct {
	Contract *StablePriceOracleCaller // Generic read-only contract binding to access the raw methods on
}

// StablePriceOracleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StablePriceOracleTransactorRaw struct {
	Contract *StablePriceOracleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStablePriceOracle creates a new instance of StablePriceOracle, bound to a specific deployed contract.
func NewStablePriceOracle(address common.Address, backend bind.ContractBackend) (*StablePriceOracle, error) {
	contract, err := bindStablePriceOracle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StablePriceOracle{StablePriceOracleCaller: StablePriceOracleCaller{contract: contract}, StablePriceOracleTransactor: StablePriceOracleTransactor{contract: contract}, StablePriceOracleFilterer: StablePriceOracleFilterer{contract: contract}}, nil
}

// NewStablePriceOracleCaller creates a new read-only instance of StablePriceOracle, bound to a specific deployed contract.
func NewStablePriceOracleCaller(address common.Address, caller bind.ContractCaller) (*StablePriceOracleCaller, error) {
	contract, err := bindStablePriceOracle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StablePriceOracleCaller{contract: contract}, nil
}

// NewStablePriceOracleTransactor creates a new write-only instance of StablePriceOracle, bound to a specific deployed contract.
func NewStablePriceOracleTransactor(address common.Address, transactor bind.ContractTransactor) (*StablePriceOracleTransactor, error) {
	contract, err := bindStablePriceOracle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StablePriceOracleTransactor{contract: contract}, nil
}

// NewStablePriceOracleFilterer creates a new log filterer instance of StablePriceOracle, bound to a specific deployed contract.
func NewStablePriceOracleFilterer(address common.Address, filterer bind.ContractFilterer) (*StablePriceOracleFilterer, error) {
	contract, err := bindStablePriceOracle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StablePriceOracleFilterer{contract: contract}, nil
}

// bindStablePriceOracle binds a generic wrapper to an already deployed contract.
func bindStablePriceOracle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StablePriceOracleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StablePriceOracle *StablePriceOracleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StablePriceOracle.Contract.StablePriceOracleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StablePriceOracle *StablePriceOracleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StablePriceOracle.Contract.StablePriceOracleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StablePriceOracle *StablePriceOracleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StablePriceOracle.Contract.StablePriceOracleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StablePriceOracle *StablePriceOracleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StablePriceOracle.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StablePriceOracle *StablePriceOracleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StablePriceOracle.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StablePriceOracle *StablePriceOracleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StablePriceOracle.Contract.contract.Transact(opts, method, params...)
}

// Premium is a free data retrieval call binding the contract method 0xa34e3596.
//
// Solidity: function premium(string name, uint256 expires, uint256 duration) view returns(uint256)
func (_StablePriceOracle *StablePriceOracleCaller) Premium(opts *bind.CallOpts, name string, expires *big.Int, duration *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _StablePriceOracle.contract.Call(opts, &out, "premium", name, expires, duration)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Premium is a free data retrieval call binding the contract method 0xa34e3596.
//
// Solidity: function premium(string name, uint256 expires, uint256 duration) view returns(uint256)
func (_StablePriceOracle *StablePriceOracleSession) Premium(name string, expires *big.Int, duration *big.Int) (*big.Int, error) {
	return _StablePriceOracle.Contract.Premium(&_StablePriceOracle.CallOpts, name, expires, duration)
}

// Premium is a free data retrieval call binding the contract method 0xa34e3596.
//
// Solidity: function premium(string name, uint256 expires, uint256 duration) view returns(uint256)
func (_StablePriceOracle *StablePriceOracleCallerSession) Premium(name string, expires *big.Int, duration *big.Int) (*big.Int, error) {
	return _StablePriceOracle.Contract.Premium(&_StablePriceOracle.CallOpts, name, expires, duration)
}

// Price is a free data retrieval call binding the contract method 0x50e9a715.
//
// Solidity: function price(string name, uint256 expires, uint256 duration) view returns((uint256,uint256))
func (_StablePriceOracle *StablePriceOracleCaller) Price(opts *bind.CallOpts, name string, expires *big.Int, duration *big.Int) (IPriceOraclePrice, error) {
	var out []interface{}
	err := _StablePriceOracle.contract.Call(opts, &out, "price", name, expires, duration)

	if err != nil {
		return *new(IPriceOraclePrice), err
	}

	out0 := *abi.ConvertType(out[0], new(IPriceOraclePrice)).(*IPriceOraclePrice)

	return out0, err

}

// Price is a free data retrieval call binding the contract method 0x50e9a715.
//
// Solidity: function price(string name, uint256 expires, uint256 duration) view returns((uint256,uint256))
func (_StablePriceOracle *StablePriceOracleSession) Price(name string, expires *big.Int, duration *big.Int) (IPriceOraclePrice, error) {
	return _StablePriceOracle.Contract.Price(&_StablePriceOracle.CallOpts, name, expires, duration)
}

// Price is a free data retrieval call binding the contract method 0x50e9a715.
//
// Solidity: function price(string name, uint256 expires, uint256 duration) view returns((uint256,uint256))
func (_StablePriceOracle *StablePriceOracleCallerSession) Price(name string, expires *big.Int, duration *big.Int) (IPriceOraclePrice, error) {
	return _StablePriceOracle.Contract.Price(&_StablePriceOracle.CallOpts, name, expires, duration)
}

// Price1Letter is a free data retrieval call binding the contract method 0x2c0fd74c.
//
// Solidity: function price1Letter() view returns(uint256)
func (_StablePriceOracle *StablePriceOracleCaller) Price1Letter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StablePriceOracle.contract.Call(opts, &out, "price1Letter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Price1Letter is a free data retrieval call binding the contract method 0x2c0fd74c.
//
// Solidity: function price1Letter() view returns(uint256)
func (_StablePriceOracle *StablePriceOracleSession) Price1Letter() (*big.Int, error) {
	return _StablePriceOracle.Contract.Price1Letter(&_StablePriceOracle.CallOpts)
}

// Price1Letter is a free data retrieval call binding the contract method 0x2c0fd74c.
//
// Solidity: function price1Letter() view returns(uint256)
func (_StablePriceOracle *StablePriceOracleCallerSession) Price1Letter() (*big.Int, error) {
	return _StablePriceOracle.Contract.Price1Letter(&_StablePriceOracle.CallOpts)
}

// Price2Letter is a free data retrieval call binding the contract method 0xcd5d2c74.
//
// Solidity: function price2Letter() view returns(uint256)
func (_StablePriceOracle *StablePriceOracleCaller) Price2Letter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StablePriceOracle.contract.Call(opts, &out, "price2Letter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Price2Letter is a free data retrieval call binding the contract method 0xcd5d2c74.
//
// Solidity: function price2Letter() view returns(uint256)
func (_StablePriceOracle *StablePriceOracleSession) Price2Letter() (*big.Int, error) {
	return _StablePriceOracle.Contract.Price2Letter(&_StablePriceOracle.CallOpts)
}

// Price2Letter is a free data retrieval call binding the contract method 0xcd5d2c74.
//
// Solidity: function price2Letter() view returns(uint256)
func (_StablePriceOracle *StablePriceOracleCallerSession) Price2Letter() (*big.Int, error) {
	return _StablePriceOracle.Contract.Price2Letter(&_StablePriceOracle.CallOpts)
}

// Price3Letter is a free data retrieval call binding the contract method 0xa200e153.
//
// Solidity: function price3Letter() view returns(uint256)
func (_StablePriceOracle *StablePriceOracleCaller) Price3Letter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StablePriceOracle.contract.Call(opts, &out, "price3Letter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Price3Letter is a free data retrieval call binding the contract method 0xa200e153.
//
// Solidity: function price3Letter() view returns(uint256)
func (_StablePriceOracle *StablePriceOracleSession) Price3Letter() (*big.Int, error) {
	return _StablePriceOracle.Contract.Price3Letter(&_StablePriceOracle.CallOpts)
}

// Price3Letter is a free data retrieval call binding the contract method 0xa200e153.
//
// Solidity: function price3Letter() view returns(uint256)
func (_StablePriceOracle *StablePriceOracleCallerSession) Price3Letter() (*big.Int, error) {
	return _StablePriceOracle.Contract.Price3Letter(&_StablePriceOracle.CallOpts)
}

// Price4Letter is a free data retrieval call binding the contract method 0xd820ed42.
//
// Solidity: function price4Letter() view returns(uint256)
func (_StablePriceOracle *StablePriceOracleCaller) Price4Letter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StablePriceOracle.contract.Call(opts, &out, "price4Letter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Price4Letter is a free data retrieval call binding the contract method 0xd820ed42.
//
// Solidity: function price4Letter() view returns(uint256)
func (_StablePriceOracle *StablePriceOracleSession) Price4Letter() (*big.Int, error) {
	return _StablePriceOracle.Contract.Price4Letter(&_StablePriceOracle.CallOpts)
}

// Price4Letter is a free data retrieval call binding the contract method 0xd820ed42.
//
// Solidity: function price4Letter() view returns(uint256)
func (_StablePriceOracle *StablePriceOracleCallerSession) Price4Letter() (*big.Int, error) {
	return _StablePriceOracle.Contract.Price4Letter(&_StablePriceOracle.CallOpts)
}

// Price5Letter is a free data retrieval call binding the contract method 0x59b6b86c.
//
// Solidity: function price5Letter() view returns(uint256)
func (_StablePriceOracle *StablePriceOracleCaller) Price5Letter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StablePriceOracle.contract.Call(opts, &out, "price5Letter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Price5Letter is a free data retrieval call binding the contract method 0x59b6b86c.
//
// Solidity: function price5Letter() view returns(uint256)
func (_StablePriceOracle *StablePriceOracleSession) Price5Letter() (*big.Int, error) {
	return _StablePriceOracle.Contract.Price5Letter(&_StablePriceOracle.CallOpts)
}

// Price5Letter is a free data retrieval call binding the contract method 0x59b6b86c.
//
// Solidity: function price5Letter() view returns(uint256)
func (_StablePriceOracle *StablePriceOracleCallerSession) Price5Letter() (*big.Int, error) {
	return _StablePriceOracle.Contract.Price5Letter(&_StablePriceOracle.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceID) view returns(bool)
func (_StablePriceOracle *StablePriceOracleCaller) SupportsInterface(opts *bind.CallOpts, interfaceID [4]byte) (bool, error) {
	var out []interface{}
	err := _StablePriceOracle.contract.Call(opts, &out, "supportsInterface", interfaceID)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceID) view returns(bool)
func (_StablePriceOracle *StablePriceOracleSession) SupportsInterface(interfaceID [4]byte) (bool, error) {
	return _StablePriceOracle.Contract.SupportsInterface(&_StablePriceOracle.CallOpts, interfaceID)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceID) view returns(bool)
func (_StablePriceOracle *StablePriceOracleCallerSession) SupportsInterface(interfaceID [4]byte) (bool, error) {
	return _StablePriceOracle.Contract.SupportsInterface(&_StablePriceOracle.CallOpts, interfaceID)
}

// UsdOracle is a free data retrieval call binding the contract method 0xc8a4271f.
//
// Solidity: function usdOracle() view returns(address)
func (_StablePriceOracle *StablePriceOracleCaller) UsdOracle(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StablePriceOracle.contract.Call(opts, &out, "usdOracle")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UsdOracle is a free data retrieval call binding the contract method 0xc8a4271f.
//
// Solidity: function usdOracle() view returns(address)
func (_StablePriceOracle *StablePriceOracleSession) UsdOracle() (common.Address, error) {
	return _StablePriceOracle.Contract.UsdOracle(&_StablePriceOracle.CallOpts)
}

// UsdOracle is a free data retrieval call binding the contract method 0xc8a4271f.
//
// Solidity: function usdOracle() view returns(address)
func (_StablePriceOracle *StablePriceOracleCallerSession) UsdOracle() (common.Address, error) {
	return _StablePriceOracle.Contract.UsdOracle(&_StablePriceOracle.CallOpts)
}

// StablePriceOracleRentPriceChangedIterator is returned from FilterRentPriceChanged and is used to iterate over the raw logs and unpacked data for RentPriceChanged events raised by the StablePriceOracle contract.
type StablePriceOracleRentPriceChangedIterator struct {
	Event *StablePriceOracleRentPriceChanged // Event containing the contract specifics and raw log

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
func (it *StablePriceOracleRentPriceChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StablePriceOracleRentPriceChanged)
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
		it.Event = new(StablePriceOracleRentPriceChanged)
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
func (it *StablePriceOracleRentPriceChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StablePriceOracleRentPriceChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StablePriceOracleRentPriceChanged represents a RentPriceChanged event raised by the StablePriceOracle contract.
type StablePriceOracleRentPriceChanged struct {
	Prices []*big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRentPriceChanged is a free log retrieval operation binding the contract event 0x73422d94aedd596c2d4d39f27a01033adc390a9054efaf259afefd95ef7331df.
//
// Solidity: event RentPriceChanged(uint256[] prices)
func (_StablePriceOracle *StablePriceOracleFilterer) FilterRentPriceChanged(opts *bind.FilterOpts) (*StablePriceOracleRentPriceChangedIterator, error) {

	logs, sub, err := _StablePriceOracle.contract.FilterLogs(opts, "RentPriceChanged")
	if err != nil {
		return nil, err
	}
	return &StablePriceOracleRentPriceChangedIterator{contract: _StablePriceOracle.contract, event: "RentPriceChanged", logs: logs, sub: sub}, nil
}

// WatchRentPriceChanged is a free log subscription operation binding the contract event 0x73422d94aedd596c2d4d39f27a01033adc390a9054efaf259afefd95ef7331df.
//
// Solidity: event RentPriceChanged(uint256[] prices)
func (_StablePriceOracle *StablePriceOracleFilterer) WatchRentPriceChanged(opts *bind.WatchOpts, sink chan<- *StablePriceOracleRentPriceChanged) (event.Subscription, error) {

	logs, sub, err := _StablePriceOracle.contract.WatchLogs(opts, "RentPriceChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StablePriceOracleRentPriceChanged)
				if err := _StablePriceOracle.contract.UnpackLog(event, "RentPriceChanged", log); err != nil {
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

// ParseRentPriceChanged is a log parse operation binding the contract event 0x73422d94aedd596c2d4d39f27a01033adc390a9054efaf259afefd95ef7331df.
//
// Solidity: event RentPriceChanged(uint256[] prices)
func (_StablePriceOracle *StablePriceOracleFilterer) ParseRentPriceChanged(log types.Log) (*StablePriceOracleRentPriceChanged, error) {
	event := new(StablePriceOracleRentPriceChanged)
	if err := _StablePriceOracle.contract.UnpackLog(event, "RentPriceChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
