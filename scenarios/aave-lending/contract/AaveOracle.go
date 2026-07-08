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

// AaveOracleMetaData contains all meta data concerning the AaveOracle contract.
var AaveOracleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIPoolAddressesProvider\",\"name\":\"provider\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"assets\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"sources\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"fallbackOracle\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"baseCurrency\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"baseCurrencyUnit\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"source\",\"type\":\"address\"}],\"name\":\"AssetSourceUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"baseCurrency\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"baseCurrencyUnit\",\"type\":\"uint256\"}],\"name\":\"BaseCurrencySet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"fallbackOracle\",\"type\":\"address\"}],\"name\":\"FallbackOracleUpdated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADDRESSES_PROVIDER\",\"outputs\":[{\"internalType\":\"contractIPoolAddressesProvider\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"BASE_CURRENCY\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"BASE_CURRENCY_UNIT\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"}],\"name\":\"getAssetPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"assets\",\"type\":\"address[]\"}],\"name\":\"getAssetsPrices\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFallbackOracle\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"}],\"name\":\"getSourceOfAsset\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"assets\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"sources\",\"type\":\"address[]\"}],\"name\":\"setAssetSources\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"fallbackOracle\",\"type\":\"address\"}],\"name\":\"setFallbackOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60e06040523480156200001157600080fd5b506040516200122b3803806200122b83398101604081905262000034916200034e565b6001600160a01b0386166080526200004c83620000ab565b620000588585620000f5565b6001600160a01b03821660a081905260c08290526040518281527fe27c4c1372396a3d15a9922f74f9dfc7c72b1ad6d63868470787249c356454c19060200160405180910390a25050505050506200049a565b600180546001600160a01b0319166001600160a01b0383169081179091556040517fce7a780d33665b1ea097af5f155e3821b809ecbaa839d3b33aa83ba28168cefb90600090a250565b8051825114604051806040016040528060028152602001611b9b60f11b815250906200013f5760405162461bcd60e51b815260040162000136919062000402565b60405180910390fd5b5060005b82518110156200025b578181815181106200016257620001626200045a565b60200260200101516000808584815181106200018257620001826200045a565b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002060006101000a8154816001600160a01b0302191690836001600160a01b03160217905550818181518110620001e357620001e36200045a565b60200260200101516001600160a01b03168382815181106200020957620002096200045a565b60200260200101516001600160a01b03167f22c5b7b2d8561d39f7f210b6b326a1aa69f15311163082308ac4877db6339dc160405160405180910390a380620002528162000470565b91505062000143565b505050565b6001600160a01b03811681146200027657600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b80516200029c8162000260565b919050565b600082601f830112620002b357600080fd5b815160206001600160401b0380831115620002d257620002d262000279565b8260051b604051601f19603f83011681018181108482111715620002fa57620002fa62000279565b6040529384528581018301938381019250878511156200031957600080fd5b83870191505b84821015620003435762000333826200028f565b835291830191908301906200031f565b979650505050505050565b60008060008060008060c087890312156200036857600080fd5b8651620003758162000260565b60208801519096506001600160401b03808211156200039357600080fd5b620003a18a838b01620002a1565b96506040890151915080821115620003b857600080fd5b50620003c789828a01620002a1565b9450506060870151620003da8162000260565b6080880151909350620003ed8162000260565b8092505060a087015190509295509295509295565b600060208083528351808285015260005b81811015620004315785810183015185820160400152820162000413565b8181111562000444576000604083870101525b50601f01601f1916929092016040019392505050565b634e487b7160e01b600052603260045260246000fd5b60006000198214156200049357634e487b7160e01b600052601160045260246000fd5b5060010190565b60805160a05160c051610d4d620004de6000396000818161013101526103a50152600081816101e5015261037a01526000818160ad01526105a30152610d4d6000f3fe608060405234801561001057600080fd5b50600436106100a35760003560e01c806392bf2be011610076578063abfd53101161005b578063abfd5310146101ba578063b3596f07146101cd578063e19f4700146101e057600080fd5b806392bf2be0146101615780639d23d9f21461019a57600080fd5b80630542975c146100a8578063170aee73146100f95780636210308c1461010e5780638c89b64f1461012c575b600080fd5b6100cf7f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b61010c610107366004610a33565b610207565b005b60015473ffffffffffffffffffffffffffffffffffffffff166100cf565b6101537f000000000000000000000000000000000000000000000000000000000000000081565b6040519081526020016100f0565b6100cf61016f366004610a33565b73ffffffffffffffffffffffffffffffffffffffff9081166000908152602081905260409020541690565b6101ad6101a8366004610a9c565b61021b565b6040516100f09190610ade565b61010c6101c8366004610b22565b6102d0565b6101536101db366004610a33565b61034b565b6100cf7f000000000000000000000000000000000000000000000000000000000000000081565b61020f61059f565b610218816107d0565b50565b606060008267ffffffffffffffff81111561023857610238610b8e565b604051908082528060200260200182016040528015610261578160200160208202803683370190505b50905060005b838110156102c85761029985858381811061028457610284610bbd565b90506020020160208101906101db9190610a33565b8282815181106102ab576102ab610bbd565b6020908102919091010152806102c081610bec565b915050610267565b509392505050565b6102d861059f565b6103458484808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152505060408051602080880282810182019093528782529093508792508691829185019084908082843760009201919091525061083f92505050565b50505050565b73ffffffffffffffffffffffffffffffffffffffff8082166000818152602081905260408120549092908116917f000000000000000000000000000000000000000000000000000000000000000090911614156103ca57507f000000000000000000000000000000000000000000000000000000000000000092915050565b73ffffffffffffffffffffffffffffffffffffffff8116610480576001546040517fb3596f0700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff85811660048301529091169063b3596f0790602401602060405180830381865afa158015610455573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104799190610c4c565b9392505050565b60008173ffffffffffffffffffffffffffffffffffffffff166350d25bcd6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156104cd573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104f19190610c4c565b90506000811315610503579392505050565b6001546040517fb3596f0700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff86811660048301529091169063b3596f0790602401602060405180830381865afa158015610573573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105979190610c4c565b949350505050565b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663707cd7166040518163ffffffff1660e01b8152600401602060405180830381865afa15801561060c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106309190610c65565b6040517f13ee32e000000000000000000000000000000000000000000000000000000000815233600482015290915073ffffffffffffffffffffffffffffffffffffffff8216906313ee32e090602401602060405180830381865afa15801561069d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106c19190610c82565b8061075557506040517f7be53ca100000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff821690637be53ca190602401602060405180830381865afa158015610731573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107559190610c82565b6040518060400160405280600181526020017f3500000000000000000000000000000000000000000000000000000000000000815250906107cc576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107c39190610ca4565b60405180910390fd5b5050565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040517fce7a780d33665b1ea097af5f155e3821b809ecbaa839d3b33aa83ba28168cefb90600090a250565b80518251146040518060400160405280600281526020017f3736000000000000000000000000000000000000000000000000000000000000815250906108b2576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107c39190610ca4565b5060005b8251811015610a0c578181815181106108d1576108d1610bbd565b60200260200101516000808584815181106108ee576108ee610bbd565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555081818151811061098057610980610bbd565b602002602001015173ffffffffffffffffffffffffffffffffffffffff168382815181106109b0576109b0610bbd565b602002602001015173ffffffffffffffffffffffffffffffffffffffff167f22c5b7b2d8561d39f7f210b6b326a1aa69f15311163082308ac4877db6339dc160405160405180910390a380610a0481610bec565b9150506108b6565b505050565b73ffffffffffffffffffffffffffffffffffffffff8116811461021857600080fd5b600060208284031215610a4557600080fd5b813561047981610a11565b60008083601f840112610a6257600080fd5b50813567ffffffffffffffff811115610a7a57600080fd5b6020830191508360208260051b8501011115610a9557600080fd5b9250929050565b60008060208385031215610aaf57600080fd5b823567ffffffffffffffff811115610ac657600080fd5b610ad285828601610a50565b90969095509350505050565b6020808252825182820181905260009190848201906040850190845b81811015610b1657835183529284019291840191600101610afa565b50909695505050505050565b60008060008060408587031215610b3857600080fd5b843567ffffffffffffffff80821115610b5057600080fd5b610b5c88838901610a50565b90965094506020870135915080821115610b7557600080fd5b50610b8287828801610a50565b95989497509550505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415610c45577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5060010190565b600060208284031215610c5e57600080fd5b5051919050565b600060208284031215610c7757600080fd5b815161047981610a11565b600060208284031215610c9457600080fd5b8151801515811461047957600080fd5b600060208083528351808285015260005b81811015610cd157858101830151858201604001528201610cb5565b81811115610ce3576000604083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01692909201604001939250505056fea26469706673582212206d11581d90e2c001849c900cb0a1bbec73d0959df8b4bbde17089057d8d95f5f64736f6c634300080a0033",
}

// AaveOracleABI is the input ABI used to generate the binding from.
// Deprecated: Use AaveOracleMetaData.ABI instead.
var AaveOracleABI = AaveOracleMetaData.ABI

// AaveOracleBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use AaveOracleMetaData.Bin instead.
var AaveOracleBin = AaveOracleMetaData.Bin

// DeployAaveOracle deploys a new Ethereum contract, binding an instance of AaveOracle to it.
func DeployAaveOracle(auth *bind.TransactOpts, backend bind.ContractBackend, provider common.Address, assets []common.Address, sources []common.Address, fallbackOracle common.Address, baseCurrency common.Address, baseCurrencyUnit *big.Int) (common.Address, *types.Transaction, *AaveOracle, error) {
	parsed, err := AaveOracleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AaveOracleBin), backend, provider, assets, sources, fallbackOracle, baseCurrency, baseCurrencyUnit)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AaveOracle{AaveOracleCaller: AaveOracleCaller{contract: contract}, AaveOracleTransactor: AaveOracleTransactor{contract: contract}, AaveOracleFilterer: AaveOracleFilterer{contract: contract}}, nil
}

// AaveOracle is an auto generated Go binding around an Ethereum contract.
type AaveOracle struct {
	AaveOracleCaller     // Read-only binding to the contract
	AaveOracleTransactor // Write-only binding to the contract
	AaveOracleFilterer   // Log filterer for contract events
}

// AaveOracleCaller is an auto generated read-only Go binding around an Ethereum contract.
type AaveOracleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AaveOracleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AaveOracleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AaveOracleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AaveOracleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AaveOracleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AaveOracleSession struct {
	Contract     *AaveOracle       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AaveOracleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AaveOracleCallerSession struct {
	Contract *AaveOracleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// AaveOracleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AaveOracleTransactorSession struct {
	Contract     *AaveOracleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// AaveOracleRaw is an auto generated low-level Go binding around an Ethereum contract.
type AaveOracleRaw struct {
	Contract *AaveOracle // Generic contract binding to access the raw methods on
}

// AaveOracleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AaveOracleCallerRaw struct {
	Contract *AaveOracleCaller // Generic read-only contract binding to access the raw methods on
}

// AaveOracleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AaveOracleTransactorRaw struct {
	Contract *AaveOracleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAaveOracle creates a new instance of AaveOracle, bound to a specific deployed contract.
func NewAaveOracle(address common.Address, backend bind.ContractBackend) (*AaveOracle, error) {
	contract, err := bindAaveOracle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AaveOracle{AaveOracleCaller: AaveOracleCaller{contract: contract}, AaveOracleTransactor: AaveOracleTransactor{contract: contract}, AaveOracleFilterer: AaveOracleFilterer{contract: contract}}, nil
}

// NewAaveOracleCaller creates a new read-only instance of AaveOracle, bound to a specific deployed contract.
func NewAaveOracleCaller(address common.Address, caller bind.ContractCaller) (*AaveOracleCaller, error) {
	contract, err := bindAaveOracle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AaveOracleCaller{contract: contract}, nil
}

// NewAaveOracleTransactor creates a new write-only instance of AaveOracle, bound to a specific deployed contract.
func NewAaveOracleTransactor(address common.Address, transactor bind.ContractTransactor) (*AaveOracleTransactor, error) {
	contract, err := bindAaveOracle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AaveOracleTransactor{contract: contract}, nil
}

// NewAaveOracleFilterer creates a new log filterer instance of AaveOracle, bound to a specific deployed contract.
func NewAaveOracleFilterer(address common.Address, filterer bind.ContractFilterer) (*AaveOracleFilterer, error) {
	contract, err := bindAaveOracle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AaveOracleFilterer{contract: contract}, nil
}

// bindAaveOracle binds a generic wrapper to an already deployed contract.
func bindAaveOracle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AaveOracleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AaveOracle *AaveOracleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AaveOracle.Contract.AaveOracleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AaveOracle *AaveOracleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AaveOracle.Contract.AaveOracleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AaveOracle *AaveOracleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AaveOracle.Contract.AaveOracleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AaveOracle *AaveOracleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AaveOracle.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AaveOracle *AaveOracleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AaveOracle.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AaveOracle *AaveOracleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AaveOracle.Contract.contract.Transact(opts, method, params...)
}

// ADDRESSESPROVIDER is a free data retrieval call binding the contract method 0x0542975c.
//
// Solidity: function ADDRESSES_PROVIDER() view returns(address)
func (_AaveOracle *AaveOracleCaller) ADDRESSESPROVIDER(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AaveOracle.contract.Call(opts, &out, "ADDRESSES_PROVIDER")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ADDRESSESPROVIDER is a free data retrieval call binding the contract method 0x0542975c.
//
// Solidity: function ADDRESSES_PROVIDER() view returns(address)
func (_AaveOracle *AaveOracleSession) ADDRESSESPROVIDER() (common.Address, error) {
	return _AaveOracle.Contract.ADDRESSESPROVIDER(&_AaveOracle.CallOpts)
}

// ADDRESSESPROVIDER is a free data retrieval call binding the contract method 0x0542975c.
//
// Solidity: function ADDRESSES_PROVIDER() view returns(address)
func (_AaveOracle *AaveOracleCallerSession) ADDRESSESPROVIDER() (common.Address, error) {
	return _AaveOracle.Contract.ADDRESSESPROVIDER(&_AaveOracle.CallOpts)
}

// BASECURRENCY is a free data retrieval call binding the contract method 0xe19f4700.
//
// Solidity: function BASE_CURRENCY() view returns(address)
func (_AaveOracle *AaveOracleCaller) BASECURRENCY(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AaveOracle.contract.Call(opts, &out, "BASE_CURRENCY")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BASECURRENCY is a free data retrieval call binding the contract method 0xe19f4700.
//
// Solidity: function BASE_CURRENCY() view returns(address)
func (_AaveOracle *AaveOracleSession) BASECURRENCY() (common.Address, error) {
	return _AaveOracle.Contract.BASECURRENCY(&_AaveOracle.CallOpts)
}

// BASECURRENCY is a free data retrieval call binding the contract method 0xe19f4700.
//
// Solidity: function BASE_CURRENCY() view returns(address)
func (_AaveOracle *AaveOracleCallerSession) BASECURRENCY() (common.Address, error) {
	return _AaveOracle.Contract.BASECURRENCY(&_AaveOracle.CallOpts)
}

// BASECURRENCYUNIT is a free data retrieval call binding the contract method 0x8c89b64f.
//
// Solidity: function BASE_CURRENCY_UNIT() view returns(uint256)
func (_AaveOracle *AaveOracleCaller) BASECURRENCYUNIT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AaveOracle.contract.Call(opts, &out, "BASE_CURRENCY_UNIT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BASECURRENCYUNIT is a free data retrieval call binding the contract method 0x8c89b64f.
//
// Solidity: function BASE_CURRENCY_UNIT() view returns(uint256)
func (_AaveOracle *AaveOracleSession) BASECURRENCYUNIT() (*big.Int, error) {
	return _AaveOracle.Contract.BASECURRENCYUNIT(&_AaveOracle.CallOpts)
}

// BASECURRENCYUNIT is a free data retrieval call binding the contract method 0x8c89b64f.
//
// Solidity: function BASE_CURRENCY_UNIT() view returns(uint256)
func (_AaveOracle *AaveOracleCallerSession) BASECURRENCYUNIT() (*big.Int, error) {
	return _AaveOracle.Contract.BASECURRENCYUNIT(&_AaveOracle.CallOpts)
}

// GetAssetPrice is a free data retrieval call binding the contract method 0xb3596f07.
//
// Solidity: function getAssetPrice(address asset) view returns(uint256)
func (_AaveOracle *AaveOracleCaller) GetAssetPrice(opts *bind.CallOpts, asset common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AaveOracle.contract.Call(opts, &out, "getAssetPrice", asset)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAssetPrice is a free data retrieval call binding the contract method 0xb3596f07.
//
// Solidity: function getAssetPrice(address asset) view returns(uint256)
func (_AaveOracle *AaveOracleSession) GetAssetPrice(asset common.Address) (*big.Int, error) {
	return _AaveOracle.Contract.GetAssetPrice(&_AaveOracle.CallOpts, asset)
}

// GetAssetPrice is a free data retrieval call binding the contract method 0xb3596f07.
//
// Solidity: function getAssetPrice(address asset) view returns(uint256)
func (_AaveOracle *AaveOracleCallerSession) GetAssetPrice(asset common.Address) (*big.Int, error) {
	return _AaveOracle.Contract.GetAssetPrice(&_AaveOracle.CallOpts, asset)
}

// GetAssetsPrices is a free data retrieval call binding the contract method 0x9d23d9f2.
//
// Solidity: function getAssetsPrices(address[] assets) view returns(uint256[])
func (_AaveOracle *AaveOracleCaller) GetAssetsPrices(opts *bind.CallOpts, assets []common.Address) ([]*big.Int, error) {
	var out []interface{}
	err := _AaveOracle.contract.Call(opts, &out, "getAssetsPrices", assets)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetAssetsPrices is a free data retrieval call binding the contract method 0x9d23d9f2.
//
// Solidity: function getAssetsPrices(address[] assets) view returns(uint256[])
func (_AaveOracle *AaveOracleSession) GetAssetsPrices(assets []common.Address) ([]*big.Int, error) {
	return _AaveOracle.Contract.GetAssetsPrices(&_AaveOracle.CallOpts, assets)
}

// GetAssetsPrices is a free data retrieval call binding the contract method 0x9d23d9f2.
//
// Solidity: function getAssetsPrices(address[] assets) view returns(uint256[])
func (_AaveOracle *AaveOracleCallerSession) GetAssetsPrices(assets []common.Address) ([]*big.Int, error) {
	return _AaveOracle.Contract.GetAssetsPrices(&_AaveOracle.CallOpts, assets)
}

// GetFallbackOracle is a free data retrieval call binding the contract method 0x6210308c.
//
// Solidity: function getFallbackOracle() view returns(address)
func (_AaveOracle *AaveOracleCaller) GetFallbackOracle(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AaveOracle.contract.Call(opts, &out, "getFallbackOracle")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetFallbackOracle is a free data retrieval call binding the contract method 0x6210308c.
//
// Solidity: function getFallbackOracle() view returns(address)
func (_AaveOracle *AaveOracleSession) GetFallbackOracle() (common.Address, error) {
	return _AaveOracle.Contract.GetFallbackOracle(&_AaveOracle.CallOpts)
}

// GetFallbackOracle is a free data retrieval call binding the contract method 0x6210308c.
//
// Solidity: function getFallbackOracle() view returns(address)
func (_AaveOracle *AaveOracleCallerSession) GetFallbackOracle() (common.Address, error) {
	return _AaveOracle.Contract.GetFallbackOracle(&_AaveOracle.CallOpts)
}

// GetSourceOfAsset is a free data retrieval call binding the contract method 0x92bf2be0.
//
// Solidity: function getSourceOfAsset(address asset) view returns(address)
func (_AaveOracle *AaveOracleCaller) GetSourceOfAsset(opts *bind.CallOpts, asset common.Address) (common.Address, error) {
	var out []interface{}
	err := _AaveOracle.contract.Call(opts, &out, "getSourceOfAsset", asset)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetSourceOfAsset is a free data retrieval call binding the contract method 0x92bf2be0.
//
// Solidity: function getSourceOfAsset(address asset) view returns(address)
func (_AaveOracle *AaveOracleSession) GetSourceOfAsset(asset common.Address) (common.Address, error) {
	return _AaveOracle.Contract.GetSourceOfAsset(&_AaveOracle.CallOpts, asset)
}

// GetSourceOfAsset is a free data retrieval call binding the contract method 0x92bf2be0.
//
// Solidity: function getSourceOfAsset(address asset) view returns(address)
func (_AaveOracle *AaveOracleCallerSession) GetSourceOfAsset(asset common.Address) (common.Address, error) {
	return _AaveOracle.Contract.GetSourceOfAsset(&_AaveOracle.CallOpts, asset)
}

// SetAssetSources is a paid mutator transaction binding the contract method 0xabfd5310.
//
// Solidity: function setAssetSources(address[] assets, address[] sources) returns()
func (_AaveOracle *AaveOracleTransactor) SetAssetSources(opts *bind.TransactOpts, assets []common.Address, sources []common.Address) (*types.Transaction, error) {
	return _AaveOracle.contract.Transact(opts, "setAssetSources", assets, sources)
}

// SetAssetSources is a paid mutator transaction binding the contract method 0xabfd5310.
//
// Solidity: function setAssetSources(address[] assets, address[] sources) returns()
func (_AaveOracle *AaveOracleSession) SetAssetSources(assets []common.Address, sources []common.Address) (*types.Transaction, error) {
	return _AaveOracle.Contract.SetAssetSources(&_AaveOracle.TransactOpts, assets, sources)
}

// SetAssetSources is a paid mutator transaction binding the contract method 0xabfd5310.
//
// Solidity: function setAssetSources(address[] assets, address[] sources) returns()
func (_AaveOracle *AaveOracleTransactorSession) SetAssetSources(assets []common.Address, sources []common.Address) (*types.Transaction, error) {
	return _AaveOracle.Contract.SetAssetSources(&_AaveOracle.TransactOpts, assets, sources)
}

// SetFallbackOracle is a paid mutator transaction binding the contract method 0x170aee73.
//
// Solidity: function setFallbackOracle(address fallbackOracle) returns()
func (_AaveOracle *AaveOracleTransactor) SetFallbackOracle(opts *bind.TransactOpts, fallbackOracle common.Address) (*types.Transaction, error) {
	return _AaveOracle.contract.Transact(opts, "setFallbackOracle", fallbackOracle)
}

// SetFallbackOracle is a paid mutator transaction binding the contract method 0x170aee73.
//
// Solidity: function setFallbackOracle(address fallbackOracle) returns()
func (_AaveOracle *AaveOracleSession) SetFallbackOracle(fallbackOracle common.Address) (*types.Transaction, error) {
	return _AaveOracle.Contract.SetFallbackOracle(&_AaveOracle.TransactOpts, fallbackOracle)
}

// SetFallbackOracle is a paid mutator transaction binding the contract method 0x170aee73.
//
// Solidity: function setFallbackOracle(address fallbackOracle) returns()
func (_AaveOracle *AaveOracleTransactorSession) SetFallbackOracle(fallbackOracle common.Address) (*types.Transaction, error) {
	return _AaveOracle.Contract.SetFallbackOracle(&_AaveOracle.TransactOpts, fallbackOracle)
}

// AaveOracleAssetSourceUpdatedIterator is returned from FilterAssetSourceUpdated and is used to iterate over the raw logs and unpacked data for AssetSourceUpdated events raised by the AaveOracle contract.
type AaveOracleAssetSourceUpdatedIterator struct {
	Event *AaveOracleAssetSourceUpdated // Event containing the contract specifics and raw log

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
func (it *AaveOracleAssetSourceUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AaveOracleAssetSourceUpdated)
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
		it.Event = new(AaveOracleAssetSourceUpdated)
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
func (it *AaveOracleAssetSourceUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AaveOracleAssetSourceUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AaveOracleAssetSourceUpdated represents a AssetSourceUpdated event raised by the AaveOracle contract.
type AaveOracleAssetSourceUpdated struct {
	Asset  common.Address
	Source common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterAssetSourceUpdated is a free log retrieval operation binding the contract event 0x22c5b7b2d8561d39f7f210b6b326a1aa69f15311163082308ac4877db6339dc1.
//
// Solidity: event AssetSourceUpdated(address indexed asset, address indexed source)
func (_AaveOracle *AaveOracleFilterer) FilterAssetSourceUpdated(opts *bind.FilterOpts, asset []common.Address, source []common.Address) (*AaveOracleAssetSourceUpdatedIterator, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}
	var sourceRule []interface{}
	for _, sourceItem := range source {
		sourceRule = append(sourceRule, sourceItem)
	}

	logs, sub, err := _AaveOracle.contract.FilterLogs(opts, "AssetSourceUpdated", assetRule, sourceRule)
	if err != nil {
		return nil, err
	}
	return &AaveOracleAssetSourceUpdatedIterator{contract: _AaveOracle.contract, event: "AssetSourceUpdated", logs: logs, sub: sub}, nil
}

// WatchAssetSourceUpdated is a free log subscription operation binding the contract event 0x22c5b7b2d8561d39f7f210b6b326a1aa69f15311163082308ac4877db6339dc1.
//
// Solidity: event AssetSourceUpdated(address indexed asset, address indexed source)
func (_AaveOracle *AaveOracleFilterer) WatchAssetSourceUpdated(opts *bind.WatchOpts, sink chan<- *AaveOracleAssetSourceUpdated, asset []common.Address, source []common.Address) (event.Subscription, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}
	var sourceRule []interface{}
	for _, sourceItem := range source {
		sourceRule = append(sourceRule, sourceItem)
	}

	logs, sub, err := _AaveOracle.contract.WatchLogs(opts, "AssetSourceUpdated", assetRule, sourceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AaveOracleAssetSourceUpdated)
				if err := _AaveOracle.contract.UnpackLog(event, "AssetSourceUpdated", log); err != nil {
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

// ParseAssetSourceUpdated is a log parse operation binding the contract event 0x22c5b7b2d8561d39f7f210b6b326a1aa69f15311163082308ac4877db6339dc1.
//
// Solidity: event AssetSourceUpdated(address indexed asset, address indexed source)
func (_AaveOracle *AaveOracleFilterer) ParseAssetSourceUpdated(log types.Log) (*AaveOracleAssetSourceUpdated, error) {
	event := new(AaveOracleAssetSourceUpdated)
	if err := _AaveOracle.contract.UnpackLog(event, "AssetSourceUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AaveOracleBaseCurrencySetIterator is returned from FilterBaseCurrencySet and is used to iterate over the raw logs and unpacked data for BaseCurrencySet events raised by the AaveOracle contract.
type AaveOracleBaseCurrencySetIterator struct {
	Event *AaveOracleBaseCurrencySet // Event containing the contract specifics and raw log

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
func (it *AaveOracleBaseCurrencySetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AaveOracleBaseCurrencySet)
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
		it.Event = new(AaveOracleBaseCurrencySet)
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
func (it *AaveOracleBaseCurrencySetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AaveOracleBaseCurrencySetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AaveOracleBaseCurrencySet represents a BaseCurrencySet event raised by the AaveOracle contract.
type AaveOracleBaseCurrencySet struct {
	BaseCurrency     common.Address
	BaseCurrencyUnit *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterBaseCurrencySet is a free log retrieval operation binding the contract event 0xe27c4c1372396a3d15a9922f74f9dfc7c72b1ad6d63868470787249c356454c1.
//
// Solidity: event BaseCurrencySet(address indexed baseCurrency, uint256 baseCurrencyUnit)
func (_AaveOracle *AaveOracleFilterer) FilterBaseCurrencySet(opts *bind.FilterOpts, baseCurrency []common.Address) (*AaveOracleBaseCurrencySetIterator, error) {

	var baseCurrencyRule []interface{}
	for _, baseCurrencyItem := range baseCurrency {
		baseCurrencyRule = append(baseCurrencyRule, baseCurrencyItem)
	}

	logs, sub, err := _AaveOracle.contract.FilterLogs(opts, "BaseCurrencySet", baseCurrencyRule)
	if err != nil {
		return nil, err
	}
	return &AaveOracleBaseCurrencySetIterator{contract: _AaveOracle.contract, event: "BaseCurrencySet", logs: logs, sub: sub}, nil
}

// WatchBaseCurrencySet is a free log subscription operation binding the contract event 0xe27c4c1372396a3d15a9922f74f9dfc7c72b1ad6d63868470787249c356454c1.
//
// Solidity: event BaseCurrencySet(address indexed baseCurrency, uint256 baseCurrencyUnit)
func (_AaveOracle *AaveOracleFilterer) WatchBaseCurrencySet(opts *bind.WatchOpts, sink chan<- *AaveOracleBaseCurrencySet, baseCurrency []common.Address) (event.Subscription, error) {

	var baseCurrencyRule []interface{}
	for _, baseCurrencyItem := range baseCurrency {
		baseCurrencyRule = append(baseCurrencyRule, baseCurrencyItem)
	}

	logs, sub, err := _AaveOracle.contract.WatchLogs(opts, "BaseCurrencySet", baseCurrencyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AaveOracleBaseCurrencySet)
				if err := _AaveOracle.contract.UnpackLog(event, "BaseCurrencySet", log); err != nil {
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

// ParseBaseCurrencySet is a log parse operation binding the contract event 0xe27c4c1372396a3d15a9922f74f9dfc7c72b1ad6d63868470787249c356454c1.
//
// Solidity: event BaseCurrencySet(address indexed baseCurrency, uint256 baseCurrencyUnit)
func (_AaveOracle *AaveOracleFilterer) ParseBaseCurrencySet(log types.Log) (*AaveOracleBaseCurrencySet, error) {
	event := new(AaveOracleBaseCurrencySet)
	if err := _AaveOracle.contract.UnpackLog(event, "BaseCurrencySet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AaveOracleFallbackOracleUpdatedIterator is returned from FilterFallbackOracleUpdated and is used to iterate over the raw logs and unpacked data for FallbackOracleUpdated events raised by the AaveOracle contract.
type AaveOracleFallbackOracleUpdatedIterator struct {
	Event *AaveOracleFallbackOracleUpdated // Event containing the contract specifics and raw log

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
func (it *AaveOracleFallbackOracleUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AaveOracleFallbackOracleUpdated)
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
		it.Event = new(AaveOracleFallbackOracleUpdated)
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
func (it *AaveOracleFallbackOracleUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AaveOracleFallbackOracleUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AaveOracleFallbackOracleUpdated represents a FallbackOracleUpdated event raised by the AaveOracle contract.
type AaveOracleFallbackOracleUpdated struct {
	FallbackOracle common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterFallbackOracleUpdated is a free log retrieval operation binding the contract event 0xce7a780d33665b1ea097af5f155e3821b809ecbaa839d3b33aa83ba28168cefb.
//
// Solidity: event FallbackOracleUpdated(address indexed fallbackOracle)
func (_AaveOracle *AaveOracleFilterer) FilterFallbackOracleUpdated(opts *bind.FilterOpts, fallbackOracle []common.Address) (*AaveOracleFallbackOracleUpdatedIterator, error) {

	var fallbackOracleRule []interface{}
	for _, fallbackOracleItem := range fallbackOracle {
		fallbackOracleRule = append(fallbackOracleRule, fallbackOracleItem)
	}

	logs, sub, err := _AaveOracle.contract.FilterLogs(opts, "FallbackOracleUpdated", fallbackOracleRule)
	if err != nil {
		return nil, err
	}
	return &AaveOracleFallbackOracleUpdatedIterator{contract: _AaveOracle.contract, event: "FallbackOracleUpdated", logs: logs, sub: sub}, nil
}

// WatchFallbackOracleUpdated is a free log subscription operation binding the contract event 0xce7a780d33665b1ea097af5f155e3821b809ecbaa839d3b33aa83ba28168cefb.
//
// Solidity: event FallbackOracleUpdated(address indexed fallbackOracle)
func (_AaveOracle *AaveOracleFilterer) WatchFallbackOracleUpdated(opts *bind.WatchOpts, sink chan<- *AaveOracleFallbackOracleUpdated, fallbackOracle []common.Address) (event.Subscription, error) {

	var fallbackOracleRule []interface{}
	for _, fallbackOracleItem := range fallbackOracle {
		fallbackOracleRule = append(fallbackOracleRule, fallbackOracleItem)
	}

	logs, sub, err := _AaveOracle.contract.WatchLogs(opts, "FallbackOracleUpdated", fallbackOracleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AaveOracleFallbackOracleUpdated)
				if err := _AaveOracle.contract.UnpackLog(event, "FallbackOracleUpdated", log); err != nil {
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

// ParseFallbackOracleUpdated is a log parse operation binding the contract event 0xce7a780d33665b1ea097af5f155e3821b809ecbaa839d3b33aa83ba28168cefb.
//
// Solidity: event FallbackOracleUpdated(address indexed fallbackOracle)
func (_AaveOracle *AaveOracleFilterer) ParseFallbackOracleUpdated(log types.Log) (*AaveOracleFallbackOracleUpdated, error) {
	event := new(AaveOracleFallbackOracleUpdated)
	if err := _AaveOracle.contract.UnpackLog(event, "FallbackOracleUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
