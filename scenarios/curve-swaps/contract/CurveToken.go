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

// CurveTokenMetaData contains all meta data concerning the CurveToken contract.
var CurveTokenMetaData = &bind.MetaData{
	ABI: "[{\"name\":\"Transfer\",\"inputs\":[{\"name\":\"_from\",\"type\":\"address\",\"indexed\":true},{\"name\":\"_to\",\"type\":\"address\",\"indexed\":true},{\"name\":\"_value\",\"type\":\"uint256\",\"indexed\":false}],\"anonymous\":false,\"type\":\"event\"},{\"name\":\"Approval\",\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\",\"indexed\":true},{\"name\":\"_spender\",\"type\":\"address\",\"indexed\":true},{\"name\":\"_value\",\"type\":\"uint256\",\"indexed\":false}],\"anonymous\":false,\"type\":\"event\"},{\"stateMutability\":\"nonpayable\",\"type\":\"constructor\",\"inputs\":[{\"name\":\"_name\",\"type\":\"string\"},{\"name\":\"_symbol\",\"type\":\"string\"},{\"name\":\"_minter\",\"type\":\"address\"}],\"outputs\":[]},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":288},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"gas\":77340},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"gas\":115282},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"gas\":37821},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"increaseAllowance\",\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_added_value\",\"type\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"gas\":40365},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"decreaseAllowance\",\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_subtracted_value\",\"type\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"gas\":40389},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"mint\",\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"gas\":79579},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"burnFrom\",\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"gas\":79597},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"set_minter\",\"inputs\":[{\"name\":\"_minter\",\"type\":\"address\"}],\"outputs\":[],\"gas\":37785},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"set_name\",\"inputs\":[{\"name\":\"_name\",\"type\":\"string\"},{\"name\":\"_symbol\",\"type\":\"string\"}],\"outputs\":[],\"gas\":181606},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"gas\":12990},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"gas\":10743},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"arg0\",\"type\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":2963},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"arg0\",\"type\":\"address\"},{\"name\":\"arg1\",\"type\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":3208},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":2808},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"minter\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"gas\":2838}]",
	Bin: "0x6060610a8f6101403960606020610a8f60c03960c051610a8f016101a03960406020610a8f60c03960c05160040135111561003957600080fd5b604060206020610a8f0160c03960c051610a8f0161022039602060206020610a8f0160c03960c05160040135111561007057600080fd5b60206040610a8f0160c03960c05160a01c1561008b57600080fd5b6101a080600060c052602060c020602082510161012060006003818352015b826101205160200211156100bd576100df565b61012051602002850151610120518501555b81516001018083528114156100aa575b50505050505061022080600160c052602060c020602082510161012060006002818352015b8261012051602002111561011757610139565b61012051602002850151610120518501555b8151600101808352811415610104575b505050505050610180516005556000610280523360007fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef6020610280a3610a7756600436101561000d576108f6565b600035601c52600051341561002157600080fd5b63313ce56781141561003857601260005260206000f35b63a9059cbb8114156100e95760043560a01c1561005457600080fd5b60023360e05260c052604060c02080546024358082101561007457600080fd5b80820390509050815550600260043560e05260c052604060c02080546024358181830110156100a257600080fd5b8082019050905081555060243561014052600435337fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef6020610140a3600160005260206000f35b6323b872dd8114156102335760043560a01c1561010557600080fd5b60243560a01c1561011557600080fd5b600260043560e05260c052604060c02080546044358082101561013757600080fd5b80820390509050815550600260243560e05260c052604060c020805460443581818301101561016557600080fd5b80820190509050815550600360043560e05260c052604060c0203360e05260c052604060c02054610140527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6101405118156101f45761014051604435808210156101cf57600080fd5b80820390509050600360043560e05260c052604060c0203360e05260c052604060c020555b604435610160526024356004357fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef6020610160a3600160005260206000f35b63095ea7b38114156102ac5760043560a01c1561024f57600080fd5b60243560033360e05260c052604060c02060043560e05260c052604060c0205560243561014052600435337f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9256020610140a3600160005260206000f35b63395093518114156103615760043560a01c156102c857600080fd5b60033360e05260c052604060c02060043560e05260c052604060c020546024358181830110156102f757600080fd5b80820190509050610140526101405160033360e05260c052604060c02060043560e05260c052604060c020556101405161016052600435337f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9256020610160a3600160005260206000f35b63a457c2d78114156104145760043560a01c1561037d57600080fd5b60033360e05260c052604060c02060043560e05260c052604060c02054602435808210156103aa57600080fd5b80820390509050610140526101405160033360e05260c052604060c02060043560e05260c052604060c020556101405161016052600435337f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9256020610160a3600160005260206000f35b6340c10f198114156104ca5760043560a01c1561043057600080fd5b600554331461043e57600080fd5b6004805460243581818301101561045457600080fd5b80820190509050815550600260043560e05260c052604060c020805460243581818301101561048257600080fd5b808201905090508155506024356101405260043560007fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef6020610140a3600160005260206000f35b6379cc679081141561057c5760043560a01c156104e657600080fd5b60055433146104f457600080fd5b600480546024358082101561050857600080fd5b80820390509050815550600260043560e05260c052604060c02080546024358082101561053457600080fd5b808203905090508155506024356101405260006004357fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef6020610140a3600160005260206000f35b631652e9fc8114156105ae5760043560a01c1561059857600080fd5b60055433146105a657600080fd5b600435600555005b63e1430e068114156106f05760606004356004016101403760406004356004013511156105da57600080fd5b60406024356004016101c03760206024356004013511156105fa57600080fd5b3360206102806004638da5cb5b6102205261023c6005545afa61061c57600080fd5b601f3d1161062957600080fd5b600050610280511461063a57600080fd5b61014080600060c052602060c020602082510161012060006003818352015b8261012051602002111561066c5761068e565b61012051602002850151610120518501555b8151600101808352811415610659575b5050505050506101c080600160c052602060c020602082510161012060006002818352015b826101205160200211156106c6576106e8565b61012051602002850151610120518501555b81516001018083528114156106b3575b505050505050005b6306fdde038114156107955760008060c052602060c020610180602082540161012060006003818352015b8261012051602002111561072e57610750565b61012051850154610120516020028501525b815160010180835281141561071b575b50505050505061018051806101a001818260206001820306601f82010390500336823750506020610160526040610180510160206001820306601f8201039050610160f35b6395d89b4181141561083a5760018060c052602060c020610180602082540161012060006002818352015b826101205160200211156107d3576107f5565b61012051850154610120516020028501525b81516001018083528114156107c0575b50505050505061018051806101a001818260206001820306601f82010390500336823750506020610160526040610180510160206001820306601f8201039050610160f35b6370a082318114156108705760043560a01c1561085657600080fd5b600260043560e05260c052604060c0205460005260206000f35b63dd62ed3e8114156108c45760043560a01c1561088c57600080fd5b60243560a01c1561089c57600080fd5b600360043560e05260c052604060c02060243560e05260c052604060c0205460005260206000f35b6318160ddd8114156108dc5760045460005260206000f35b63075461728114156108f45760055460005260206000f35b505b60006000fd5b61017b610a770361017b60003961017b610a77036000f3",
}

// CurveTokenABI is the input ABI used to generate the binding from.
// Deprecated: Use CurveTokenMetaData.ABI instead.
var CurveTokenABI = CurveTokenMetaData.ABI

// CurveTokenBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use CurveTokenMetaData.Bin instead.
var CurveTokenBin = CurveTokenMetaData.Bin

// DeployCurveToken deploys a new Ethereum contract, binding an instance of CurveToken to it.
func DeployCurveToken(auth *bind.TransactOpts, backend bind.ContractBackend, _name string, _symbol string, _minter common.Address) (common.Address, *types.Transaction, *CurveToken, error) {
	parsed, err := CurveTokenMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CurveTokenBin), backend, _name, _symbol, _minter)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CurveToken{CurveTokenCaller: CurveTokenCaller{contract: contract}, CurveTokenTransactor: CurveTokenTransactor{contract: contract}, CurveTokenFilterer: CurveTokenFilterer{contract: contract}}, nil
}

// CurveToken is an auto generated Go binding around an Ethereum contract.
type CurveToken struct {
	CurveTokenCaller     // Read-only binding to the contract
	CurveTokenTransactor // Write-only binding to the contract
	CurveTokenFilterer   // Log filterer for contract events
}

// CurveTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type CurveTokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CurveTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CurveTokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CurveTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CurveTokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CurveTokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CurveTokenSession struct {
	Contract     *CurveToken       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CurveTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CurveTokenCallerSession struct {
	Contract *CurveTokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// CurveTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CurveTokenTransactorSession struct {
	Contract     *CurveTokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// CurveTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type CurveTokenRaw struct {
	Contract *CurveToken // Generic contract binding to access the raw methods on
}

// CurveTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CurveTokenCallerRaw struct {
	Contract *CurveTokenCaller // Generic read-only contract binding to access the raw methods on
}

// CurveTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CurveTokenTransactorRaw struct {
	Contract *CurveTokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCurveToken creates a new instance of CurveToken, bound to a specific deployed contract.
func NewCurveToken(address common.Address, backend bind.ContractBackend) (*CurveToken, error) {
	contract, err := bindCurveToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CurveToken{CurveTokenCaller: CurveTokenCaller{contract: contract}, CurveTokenTransactor: CurveTokenTransactor{contract: contract}, CurveTokenFilterer: CurveTokenFilterer{contract: contract}}, nil
}

// NewCurveTokenCaller creates a new read-only instance of CurveToken, bound to a specific deployed contract.
func NewCurveTokenCaller(address common.Address, caller bind.ContractCaller) (*CurveTokenCaller, error) {
	contract, err := bindCurveToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CurveTokenCaller{contract: contract}, nil
}

// NewCurveTokenTransactor creates a new write-only instance of CurveToken, bound to a specific deployed contract.
func NewCurveTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*CurveTokenTransactor, error) {
	contract, err := bindCurveToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CurveTokenTransactor{contract: contract}, nil
}

// NewCurveTokenFilterer creates a new log filterer instance of CurveToken, bound to a specific deployed contract.
func NewCurveTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*CurveTokenFilterer, error) {
	contract, err := bindCurveToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CurveTokenFilterer{contract: contract}, nil
}

// bindCurveToken binds a generic wrapper to an already deployed contract.
func bindCurveToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CurveTokenMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CurveToken *CurveTokenRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CurveToken.Contract.CurveTokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CurveToken *CurveTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CurveToken.Contract.CurveTokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CurveToken *CurveTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CurveToken.Contract.CurveTokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CurveToken *CurveTokenCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CurveToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CurveToken *CurveTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CurveToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CurveToken *CurveTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CurveToken.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address arg0, address arg1) view returns(uint256)
func (_CurveToken *CurveTokenCaller) Allowance(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _CurveToken.contract.Call(opts, &out, "allowance", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address arg0, address arg1) view returns(uint256)
func (_CurveToken *CurveTokenSession) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _CurveToken.Contract.Allowance(&_CurveToken.CallOpts, arg0, arg1)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address arg0, address arg1) view returns(uint256)
func (_CurveToken *CurveTokenCallerSession) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _CurveToken.Contract.Allowance(&_CurveToken.CallOpts, arg0, arg1)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address arg0) view returns(uint256)
func (_CurveToken *CurveTokenCaller) BalanceOf(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _CurveToken.contract.Call(opts, &out, "balanceOf", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address arg0) view returns(uint256)
func (_CurveToken *CurveTokenSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _CurveToken.Contract.BalanceOf(&_CurveToken.CallOpts, arg0)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address arg0) view returns(uint256)
func (_CurveToken *CurveTokenCallerSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _CurveToken.Contract.BalanceOf(&_CurveToken.CallOpts, arg0)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint256)
func (_CurveToken *CurveTokenCaller) Decimals(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CurveToken.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint256)
func (_CurveToken *CurveTokenSession) Decimals() (*big.Int, error) {
	return _CurveToken.Contract.Decimals(&_CurveToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint256)
func (_CurveToken *CurveTokenCallerSession) Decimals() (*big.Int, error) {
	return _CurveToken.Contract.Decimals(&_CurveToken.CallOpts)
}

// Minter is a free data retrieval call binding the contract method 0x07546172.
//
// Solidity: function minter() view returns(address)
func (_CurveToken *CurveTokenCaller) Minter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CurveToken.contract.Call(opts, &out, "minter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Minter is a free data retrieval call binding the contract method 0x07546172.
//
// Solidity: function minter() view returns(address)
func (_CurveToken *CurveTokenSession) Minter() (common.Address, error) {
	return _CurveToken.Contract.Minter(&_CurveToken.CallOpts)
}

// Minter is a free data retrieval call binding the contract method 0x07546172.
//
// Solidity: function minter() view returns(address)
func (_CurveToken *CurveTokenCallerSession) Minter() (common.Address, error) {
	return _CurveToken.Contract.Minter(&_CurveToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_CurveToken *CurveTokenCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _CurveToken.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_CurveToken *CurveTokenSession) Name() (string, error) {
	return _CurveToken.Contract.Name(&_CurveToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_CurveToken *CurveTokenCallerSession) Name() (string, error) {
	return _CurveToken.Contract.Name(&_CurveToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_CurveToken *CurveTokenCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _CurveToken.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_CurveToken *CurveTokenSession) Symbol() (string, error) {
	return _CurveToken.Contract.Symbol(&_CurveToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_CurveToken *CurveTokenCallerSession) Symbol() (string, error) {
	return _CurveToken.Contract.Symbol(&_CurveToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_CurveToken *CurveTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CurveToken.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_CurveToken *CurveTokenSession) TotalSupply() (*big.Int, error) {
	return _CurveToken.Contract.TotalSupply(&_CurveToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_CurveToken *CurveTokenCallerSession) TotalSupply() (*big.Int, error) {
	return _CurveToken.Contract.TotalSupply(&_CurveToken.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _spender, uint256 _value) returns(bool)
func (_CurveToken *CurveTokenTransactor) Approve(opts *bind.TransactOpts, _spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _CurveToken.contract.Transact(opts, "approve", _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _spender, uint256 _value) returns(bool)
func (_CurveToken *CurveTokenSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _CurveToken.Contract.Approve(&_CurveToken.TransactOpts, _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _spender, uint256 _value) returns(bool)
func (_CurveToken *CurveTokenTransactorSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _CurveToken.Contract.Approve(&_CurveToken.TransactOpts, _spender, _value)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address _to, uint256 _value) returns(bool)
func (_CurveToken *CurveTokenTransactor) BurnFrom(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _CurveToken.contract.Transact(opts, "burnFrom", _to, _value)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address _to, uint256 _value) returns(bool)
func (_CurveToken *CurveTokenSession) BurnFrom(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _CurveToken.Contract.BurnFrom(&_CurveToken.TransactOpts, _to, _value)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address _to, uint256 _value) returns(bool)
func (_CurveToken *CurveTokenTransactorSession) BurnFrom(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _CurveToken.Contract.BurnFrom(&_CurveToken.TransactOpts, _to, _value)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address _spender, uint256 _subtracted_value) returns(bool)
func (_CurveToken *CurveTokenTransactor) DecreaseAllowance(opts *bind.TransactOpts, _spender common.Address, _subtracted_value *big.Int) (*types.Transaction, error) {
	return _CurveToken.contract.Transact(opts, "decreaseAllowance", _spender, _subtracted_value)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address _spender, uint256 _subtracted_value) returns(bool)
func (_CurveToken *CurveTokenSession) DecreaseAllowance(_spender common.Address, _subtracted_value *big.Int) (*types.Transaction, error) {
	return _CurveToken.Contract.DecreaseAllowance(&_CurveToken.TransactOpts, _spender, _subtracted_value)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address _spender, uint256 _subtracted_value) returns(bool)
func (_CurveToken *CurveTokenTransactorSession) DecreaseAllowance(_spender common.Address, _subtracted_value *big.Int) (*types.Transaction, error) {
	return _CurveToken.Contract.DecreaseAllowance(&_CurveToken.TransactOpts, _spender, _subtracted_value)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address _spender, uint256 _added_value) returns(bool)
func (_CurveToken *CurveTokenTransactor) IncreaseAllowance(opts *bind.TransactOpts, _spender common.Address, _added_value *big.Int) (*types.Transaction, error) {
	return _CurveToken.contract.Transact(opts, "increaseAllowance", _spender, _added_value)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address _spender, uint256 _added_value) returns(bool)
func (_CurveToken *CurveTokenSession) IncreaseAllowance(_spender common.Address, _added_value *big.Int) (*types.Transaction, error) {
	return _CurveToken.Contract.IncreaseAllowance(&_CurveToken.TransactOpts, _spender, _added_value)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address _spender, uint256 _added_value) returns(bool)
func (_CurveToken *CurveTokenTransactorSession) IncreaseAllowance(_spender common.Address, _added_value *big.Int) (*types.Transaction, error) {
	return _CurveToken.Contract.IncreaseAllowance(&_CurveToken.TransactOpts, _spender, _added_value)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _to, uint256 _value) returns(bool)
func (_CurveToken *CurveTokenTransactor) Mint(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _CurveToken.contract.Transact(opts, "mint", _to, _value)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _to, uint256 _value) returns(bool)
func (_CurveToken *CurveTokenSession) Mint(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _CurveToken.Contract.Mint(&_CurveToken.TransactOpts, _to, _value)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _to, uint256 _value) returns(bool)
func (_CurveToken *CurveTokenTransactorSession) Mint(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _CurveToken.Contract.Mint(&_CurveToken.TransactOpts, _to, _value)
}

// SetMinter is a paid mutator transaction binding the contract method 0x1652e9fc.
//
// Solidity: function set_minter(address _minter) returns()
func (_CurveToken *CurveTokenTransactor) SetMinter(opts *bind.TransactOpts, _minter common.Address) (*types.Transaction, error) {
	return _CurveToken.contract.Transact(opts, "set_minter", _minter)
}

// SetMinter is a paid mutator transaction binding the contract method 0x1652e9fc.
//
// Solidity: function set_minter(address _minter) returns()
func (_CurveToken *CurveTokenSession) SetMinter(_minter common.Address) (*types.Transaction, error) {
	return _CurveToken.Contract.SetMinter(&_CurveToken.TransactOpts, _minter)
}

// SetMinter is a paid mutator transaction binding the contract method 0x1652e9fc.
//
// Solidity: function set_minter(address _minter) returns()
func (_CurveToken *CurveTokenTransactorSession) SetMinter(_minter common.Address) (*types.Transaction, error) {
	return _CurveToken.Contract.SetMinter(&_CurveToken.TransactOpts, _minter)
}

// SetName is a paid mutator transaction binding the contract method 0xe1430e06.
//
// Solidity: function set_name(string _name, string _symbol) returns()
func (_CurveToken *CurveTokenTransactor) SetName(opts *bind.TransactOpts, _name string, _symbol string) (*types.Transaction, error) {
	return _CurveToken.contract.Transact(opts, "set_name", _name, _symbol)
}

// SetName is a paid mutator transaction binding the contract method 0xe1430e06.
//
// Solidity: function set_name(string _name, string _symbol) returns()
func (_CurveToken *CurveTokenSession) SetName(_name string, _symbol string) (*types.Transaction, error) {
	return _CurveToken.Contract.SetName(&_CurveToken.TransactOpts, _name, _symbol)
}

// SetName is a paid mutator transaction binding the contract method 0xe1430e06.
//
// Solidity: function set_name(string _name, string _symbol) returns()
func (_CurveToken *CurveTokenTransactorSession) SetName(_name string, _symbol string) (*types.Transaction, error) {
	return _CurveToken.Contract.SetName(&_CurveToken.TransactOpts, _name, _symbol)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _value) returns(bool)
func (_CurveToken *CurveTokenTransactor) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _CurveToken.contract.Transact(opts, "transfer", _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _value) returns(bool)
func (_CurveToken *CurveTokenSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _CurveToken.Contract.Transfer(&_CurveToken.TransactOpts, _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _value) returns(bool)
func (_CurveToken *CurveTokenTransactorSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _CurveToken.Contract.Transfer(&_CurveToken.TransactOpts, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _from, address _to, uint256 _value) returns(bool)
func (_CurveToken *CurveTokenTransactor) TransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _CurveToken.contract.Transact(opts, "transferFrom", _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _from, address _to, uint256 _value) returns(bool)
func (_CurveToken *CurveTokenSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _CurveToken.Contract.TransferFrom(&_CurveToken.TransactOpts, _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _from, address _to, uint256 _value) returns(bool)
func (_CurveToken *CurveTokenTransactorSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _CurveToken.Contract.TransferFrom(&_CurveToken.TransactOpts, _from, _to, _value)
}

// CurveTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the CurveToken contract.
type CurveTokenApprovalIterator struct {
	Event *CurveTokenApproval // Event containing the contract specifics and raw log

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
func (it *CurveTokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CurveTokenApproval)
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
		it.Event = new(CurveTokenApproval)
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
func (it *CurveTokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CurveTokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CurveTokenApproval represents a Approval event raised by the CurveToken contract.
type CurveTokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed _owner, address indexed _spender, uint256 _value)
func (_CurveToken *CurveTokenFilterer) FilterApproval(opts *bind.FilterOpts, _owner []common.Address, _spender []common.Address) (*CurveTokenApprovalIterator, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _spenderRule []interface{}
	for _, _spenderItem := range _spender {
		_spenderRule = append(_spenderRule, _spenderItem)
	}

	logs, sub, err := _CurveToken.contract.FilterLogs(opts, "Approval", _ownerRule, _spenderRule)
	if err != nil {
		return nil, err
	}
	return &CurveTokenApprovalIterator{contract: _CurveToken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed _owner, address indexed _spender, uint256 _value)
func (_CurveToken *CurveTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *CurveTokenApproval, _owner []common.Address, _spender []common.Address) (event.Subscription, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _spenderRule []interface{}
	for _, _spenderItem := range _spender {
		_spenderRule = append(_spenderRule, _spenderItem)
	}

	logs, sub, err := _CurveToken.contract.WatchLogs(opts, "Approval", _ownerRule, _spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CurveTokenApproval)
				if err := _CurveToken.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed _owner, address indexed _spender, uint256 _value)
func (_CurveToken *CurveTokenFilterer) ParseApproval(log types.Log) (*CurveTokenApproval, error) {
	event := new(CurveTokenApproval)
	if err := _CurveToken.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CurveTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the CurveToken contract.
type CurveTokenTransferIterator struct {
	Event *CurveTokenTransfer // Event containing the contract specifics and raw log

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
func (it *CurveTokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CurveTokenTransfer)
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
		it.Event = new(CurveTokenTransfer)
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
func (it *CurveTokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CurveTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CurveTokenTransfer represents a Transfer event raised by the CurveToken contract.
type CurveTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed _from, address indexed _to, uint256 _value)
func (_CurveToken *CurveTokenFilterer) FilterTransfer(opts *bind.FilterOpts, _from []common.Address, _to []common.Address) (*CurveTokenTransferIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _CurveToken.contract.FilterLogs(opts, "Transfer", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return &CurveTokenTransferIterator{contract: _CurveToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed _from, address indexed _to, uint256 _value)
func (_CurveToken *CurveTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *CurveTokenTransfer, _from []common.Address, _to []common.Address) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _CurveToken.contract.WatchLogs(opts, "Transfer", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CurveTokenTransfer)
				if err := _CurveToken.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed _from, address indexed _to, uint256 _value)
func (_CurveToken *CurveTokenFilterer) ParseTransfer(log types.Log) (*CurveTokenTransfer, error) {
	event := new(CurveTokenTransfer)
	if err := _CurveToken.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
