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

// DataTypesCalculateInterestRatesParams is an auto generated low-level Go binding around an user-defined struct.
type DataTypesCalculateInterestRatesParams struct {
	Unbacked                *big.Int
	LiquidityAdded          *big.Int
	LiquidityTaken          *big.Int
	TotalStableDebt         *big.Int
	TotalVariableDebt       *big.Int
	AverageStableBorrowRate *big.Int
	ReserveFactor           *big.Int
	Reserve                 common.Address
	AToken                  common.Address
}

// DefaultReserveInterestRateStrategyMetaData contains all meta data concerning the DefaultReserveInterestRateStrategy contract.
var DefaultReserveInterestRateStrategyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIPoolAddressesProvider\",\"name\":\"provider\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"optimalUsageRatio\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseVariableBorrowRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"variableRateSlope1\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"variableRateSlope2\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stableRateSlope1\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stableRateSlope2\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseStableRateOffset\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stableRateExcessOffset\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"optimalStableToTotalDebtRatio\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ADDRESSES_PROVIDER\",\"outputs\":[{\"internalType\":\"contractIPoolAddressesProvider\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_EXCESS_STABLE_TO_TOTAL_DEBT_RATIO\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_EXCESS_USAGE_RATIO\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"OPTIMAL_STABLE_TO_TOTAL_DEBT_RATIO\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"OPTIMAL_USAGE_RATIO\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"unbacked\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"liquidityAdded\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"liquidityTaken\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalStableDebt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalVariableDebt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"averageStableBorrowRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserveFactor\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"reserve\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"aToken\",\"type\":\"address\"}],\"internalType\":\"structDataTypes.CalculateInterestRatesParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"calculateInterestRates\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBaseStableBorrowRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBaseVariableBorrowRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMaxVariableBorrowRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStableRateExcessOffset\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStableRateSlope1\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStableRateSlope2\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVariableRateSlope1\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVariableRateSlope2\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x61020060405234801561001157600080fd5b5060405162000f7538038062000f7583398101604081905261003291610146565b886b033b2e3c9fd0803ce8000000101560405180604001604052806002815260200161383360f01b815250906100845760405162461bcd60e51b815260040161007b91906101d1565b60405180910390fd5b50806b033b2e3c9fd0803ce80000001015604051806040016040528060028152602001610e0d60f21b815250906100ce5760405162461bcd60e51b815260040161007b91906101d1565b5060808990526100ea896b033b2e3c9fd0803ce8000000610226565b60c05260a0819052610108816b033b2e3c9fd0803ce8000000610226565b60e052506001600160a01b0390981661010052610120959095526101409390935261016091909152610180526101a0526101c052506101e05261024b565b6000806000806000806000806000806101408b8d03121561016657600080fd5b8a516001600160a01b038116811461017d57600080fd5b809a505060208b0151985060408b0151975060608b0151965060808b0151955060a08b0151945060c08b0151935060e08b015192506101008b015191506101208b015190509295989b9194979a5092959850565b600060208083528351808285015260005b818110156101fe578581018301518582016040015282016101e2565b81811115610210576000604083870101525b50601f01601f1916929092016040019392505050565b60008282101561024657634e487b7160e01b600052601160045260246000fd5b500390565b60805160a05160c05160e05161010051610120516101405161016051610180516101a0516101c0516101e051610c0d62000368600039600081816102710152610821015260006108c601526000818161017201526105ec0152600081816102970152818161061701526106ec0152600081816102bd0152818161030c0152610654015260008181610142015281816103300152818161067f0152818161075e01526108e70152600081816101980152818161035101526103fa0152600060f40152600081816102e601526107cb01526000818161024501526105900152600081816101e80152818161079a01526107ec0152600081816101c10152818161055f015281816105b1015281816106c301526107380152610c0d6000f3fe608060405234801561001057600080fd5b50600436106100ea5760003560e01c8063a58987091161008c578063bc62690811610066578063bc6269081461026f578063d5cd739114610295578063f4202409146102bb578063fe5fd698146102e157600080fd5b8063a589870914610212578063a9c622f814610240578063acd786861461026757600080fd5b806334762ca5116100c857806334762ca51461019657806354c365c6146101bc5780636fb92589146101e357806380031e371461020a57600080fd5b80630542975c146100ef5780630b3429a21461014057806314e32da414610170575b600080fd5b6101167f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b7f00000000000000000000000000000000000000000000000000000000000000005b604051908152602001610137565b7f0000000000000000000000000000000000000000000000000000000000000000610162565b7f0000000000000000000000000000000000000000000000000000000000000000610162565b6101627f000000000000000000000000000000000000000000000000000000000000000081565b6101627f000000000000000000000000000000000000000000000000000000000000000081565b610162610308565b610225610220366004610adb565b610384565b60408051938452602084019290925290820152606001610137565b6101627f000000000000000000000000000000000000000000000000000000000000000081565b6101626108bf565b7f0000000000000000000000000000000000000000000000000000000000000000610162565b7f0000000000000000000000000000000000000000000000000000000000000000610162565b7f0000000000000000000000000000000000000000000000000000000000000000610162565b6101627f000000000000000000000000000000000000000000000000000000000000000081565b60007f00000000000000000000000000000000000000000000000000000000000000006103757f00000000000000000000000000000000000000000000000000000000000000007f0000000000000000000000000000000000000000000000000000000000000000610b8f565b61037f9190610b8f565b905090565b60008060006103d86040518061012001604052806000815260200160008152602001600081526020016000815260200160008152602001600081526020016000815260200160008152602001600081525090565b846080015185606001516103ec9190610b8f565b6020820152600060808201527f000000000000000000000000000000000000000000000000000000000000000060408201526104266108bf565b606082015260208101511561055d57602081015160608601516104489161090b565b60e08083019190915260408087015160208801519288015161010089015192517f70a0823100000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff938416600482015291939216906370a0823190602401602060405180830381865afa1580156104d3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104f79190610ba7565b6105019190610b8f565b61050b9190610bc0565b808252602082015161051c91610b8f565b610100820181905260208201516105329161090b565b60a082015284516101008201516105579161054c91610b8f565b60208301519061090b565b60c08201525b7f00000000000000000000000000000000000000000000000000000000000000008160a0015111156106be5760006105e57f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000008460a001516105df9190610bc0565b9061090b565b90506106117f00000000000000000000000000000000000000000000000000000000000000008261094a565b61063b907f0000000000000000000000000000000000000000000000000000000000000000610b8f565b8260600181815161064c9190610b8f565b9052506106797f00000000000000000000000000000000000000000000000000000000000000008261094a565b6106a3907f0000000000000000000000000000000000000000000000000000000000000000610b8f565b826040018181516106b49190610b8f565b9052506107989050565b6107197f00000000000000000000000000000000000000000000000000000000000000006105df8360a001517f000000000000000000000000000000000000000000000000000000000000000061094a90919063ffffffff16565b8160600181815161072a9190610b8f565b90525060a0810151610783907f0000000000000000000000000000000000000000000000000000000000000000906105df907f00000000000000000000000000000000000000000000000000000000000000009061094a565b816040018181516107949190610b8f565b9052505b7f00000000000000000000000000000000000000000000000000000000000000008160e00151111561085c57600061081a7f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000008460e001516105df9190610bc0565b90506108467f00000000000000000000000000000000000000000000000000000000000000008261094a565b826060018181516108579190610b8f565b905250505b6108a18560c001516127106108719190610bc0565b61089b8360c0015161089589606001518a6080015187604001518c60a001516109a1565b9061094a565b90610a08565b60808201819052606082015160409092015190969195509350915050565b600061037f7f00000000000000000000000000000000000000000000000000000000000000007f0000000000000000000000000000000000000000000000000000000000000000610b8f565b600081156b033b2e3c9fd0803ce80000006002840419048411171561092f57600080fd5b506b033b2e3c9fd0803ce80000009190910260028204010490565b600081157ffffffffffffffffffffffffffffffffffffffffffe6268e1b017bfe18bffffff8390048411151761097f57600080fd5b506b033b2e3c9fd0803ce800000091026b019d971e4fe8401e74000000010490565b6000806109ae8587610b8f565b9050806109bf576000915050610a00565b60006109ce8561089588610a4b565b905060006109df856108958a610a4b565b905060006109f96109ef85610a4b565b6105df8486610b8f565b9450505050505b949350505050565b600081157fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffec7783900484111517610a3d57600080fd5b506127109102611388010490565b633b9aca008181029081048214610a6157600080fd5b919050565b604051610120810167ffffffffffffffff81118282101715610ab1577f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405290565b803573ffffffffffffffffffffffffffffffffffffffff81168114610a6157600080fd5b60006101208284031215610aee57600080fd5b610af6610a66565b823581526020830135602082015260408301356040820152606083013560608201526080830135608082015260a083013560a082015260c083013560c0820152610b4260e08401610ab7565b60e0820152610100610b55818501610ab7565b908201529392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60008219821115610ba257610ba2610b60565b500190565b600060208284031215610bb957600080fd5b5051919050565b600082821015610bd257610bd2610b60565b50039056fea2646970667358221220f1e92c2f40dc3cef47809ecfaf576df226cdd354a99c1d82b228c4c71b0e0dc164736f6c634300080a0033",
}

// DefaultReserveInterestRateStrategyABI is the input ABI used to generate the binding from.
// Deprecated: Use DefaultReserveInterestRateStrategyMetaData.ABI instead.
var DefaultReserveInterestRateStrategyABI = DefaultReserveInterestRateStrategyMetaData.ABI

// DefaultReserveInterestRateStrategyBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DefaultReserveInterestRateStrategyMetaData.Bin instead.
var DefaultReserveInterestRateStrategyBin = DefaultReserveInterestRateStrategyMetaData.Bin

// DeployDefaultReserveInterestRateStrategy deploys a new Ethereum contract, binding an instance of DefaultReserveInterestRateStrategy to it.
func DeployDefaultReserveInterestRateStrategy(auth *bind.TransactOpts, backend bind.ContractBackend, provider common.Address, optimalUsageRatio *big.Int, baseVariableBorrowRate *big.Int, variableRateSlope1 *big.Int, variableRateSlope2 *big.Int, stableRateSlope1 *big.Int, stableRateSlope2 *big.Int, baseStableRateOffset *big.Int, stableRateExcessOffset *big.Int, optimalStableToTotalDebtRatio *big.Int) (common.Address, *types.Transaction, *DefaultReserveInterestRateStrategy, error) {
	parsed, err := DefaultReserveInterestRateStrategyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DefaultReserveInterestRateStrategyBin), backend, provider, optimalUsageRatio, baseVariableBorrowRate, variableRateSlope1, variableRateSlope2, stableRateSlope1, stableRateSlope2, baseStableRateOffset, stableRateExcessOffset, optimalStableToTotalDebtRatio)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DefaultReserveInterestRateStrategy{DefaultReserveInterestRateStrategyCaller: DefaultReserveInterestRateStrategyCaller{contract: contract}, DefaultReserveInterestRateStrategyTransactor: DefaultReserveInterestRateStrategyTransactor{contract: contract}, DefaultReserveInterestRateStrategyFilterer: DefaultReserveInterestRateStrategyFilterer{contract: contract}}, nil
}

// DefaultReserveInterestRateStrategy is an auto generated Go binding around an Ethereum contract.
type DefaultReserveInterestRateStrategy struct {
	DefaultReserveInterestRateStrategyCaller     // Read-only binding to the contract
	DefaultReserveInterestRateStrategyTransactor // Write-only binding to the contract
	DefaultReserveInterestRateStrategyFilterer   // Log filterer for contract events
}

// DefaultReserveInterestRateStrategyCaller is an auto generated read-only Go binding around an Ethereum contract.
type DefaultReserveInterestRateStrategyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DefaultReserveInterestRateStrategyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DefaultReserveInterestRateStrategyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DefaultReserveInterestRateStrategyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DefaultReserveInterestRateStrategyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DefaultReserveInterestRateStrategySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DefaultReserveInterestRateStrategySession struct {
	Contract     *DefaultReserveInterestRateStrategy // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                       // Call options to use throughout this session
	TransactOpts bind.TransactOpts                   // Transaction auth options to use throughout this session
}

// DefaultReserveInterestRateStrategyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DefaultReserveInterestRateStrategyCallerSession struct {
	Contract *DefaultReserveInterestRateStrategyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                             // Call options to use throughout this session
}

// DefaultReserveInterestRateStrategyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DefaultReserveInterestRateStrategyTransactorSession struct {
	Contract     *DefaultReserveInterestRateStrategyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                             // Transaction auth options to use throughout this session
}

// DefaultReserveInterestRateStrategyRaw is an auto generated low-level Go binding around an Ethereum contract.
type DefaultReserveInterestRateStrategyRaw struct {
	Contract *DefaultReserveInterestRateStrategy // Generic contract binding to access the raw methods on
}

// DefaultReserveInterestRateStrategyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DefaultReserveInterestRateStrategyCallerRaw struct {
	Contract *DefaultReserveInterestRateStrategyCaller // Generic read-only contract binding to access the raw methods on
}

// DefaultReserveInterestRateStrategyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DefaultReserveInterestRateStrategyTransactorRaw struct {
	Contract *DefaultReserveInterestRateStrategyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDefaultReserveInterestRateStrategy creates a new instance of DefaultReserveInterestRateStrategy, bound to a specific deployed contract.
func NewDefaultReserveInterestRateStrategy(address common.Address, backend bind.ContractBackend) (*DefaultReserveInterestRateStrategy, error) {
	contract, err := bindDefaultReserveInterestRateStrategy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DefaultReserveInterestRateStrategy{DefaultReserveInterestRateStrategyCaller: DefaultReserveInterestRateStrategyCaller{contract: contract}, DefaultReserveInterestRateStrategyTransactor: DefaultReserveInterestRateStrategyTransactor{contract: contract}, DefaultReserveInterestRateStrategyFilterer: DefaultReserveInterestRateStrategyFilterer{contract: contract}}, nil
}

// NewDefaultReserveInterestRateStrategyCaller creates a new read-only instance of DefaultReserveInterestRateStrategy, bound to a specific deployed contract.
func NewDefaultReserveInterestRateStrategyCaller(address common.Address, caller bind.ContractCaller) (*DefaultReserveInterestRateStrategyCaller, error) {
	contract, err := bindDefaultReserveInterestRateStrategy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DefaultReserveInterestRateStrategyCaller{contract: contract}, nil
}

// NewDefaultReserveInterestRateStrategyTransactor creates a new write-only instance of DefaultReserveInterestRateStrategy, bound to a specific deployed contract.
func NewDefaultReserveInterestRateStrategyTransactor(address common.Address, transactor bind.ContractTransactor) (*DefaultReserveInterestRateStrategyTransactor, error) {
	contract, err := bindDefaultReserveInterestRateStrategy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DefaultReserveInterestRateStrategyTransactor{contract: contract}, nil
}

// NewDefaultReserveInterestRateStrategyFilterer creates a new log filterer instance of DefaultReserveInterestRateStrategy, bound to a specific deployed contract.
func NewDefaultReserveInterestRateStrategyFilterer(address common.Address, filterer bind.ContractFilterer) (*DefaultReserveInterestRateStrategyFilterer, error) {
	contract, err := bindDefaultReserveInterestRateStrategy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DefaultReserveInterestRateStrategyFilterer{contract: contract}, nil
}

// bindDefaultReserveInterestRateStrategy binds a generic wrapper to an already deployed contract.
func bindDefaultReserveInterestRateStrategy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DefaultReserveInterestRateStrategyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DefaultReserveInterestRateStrategy.Contract.DefaultReserveInterestRateStrategyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DefaultReserveInterestRateStrategy.Contract.DefaultReserveInterestRateStrategyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DefaultReserveInterestRateStrategy.Contract.DefaultReserveInterestRateStrategyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DefaultReserveInterestRateStrategy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DefaultReserveInterestRateStrategy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DefaultReserveInterestRateStrategy.Contract.contract.Transact(opts, method, params...)
}

// ADDRESSESPROVIDER is a free data retrieval call binding the contract method 0x0542975c.
//
// Solidity: function ADDRESSES_PROVIDER() view returns(address)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategyCaller) ADDRESSESPROVIDER(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DefaultReserveInterestRateStrategy.contract.Call(opts, &out, "ADDRESSES_PROVIDER")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ADDRESSESPROVIDER is a free data retrieval call binding the contract method 0x0542975c.
//
// Solidity: function ADDRESSES_PROVIDER() view returns(address)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategySession) ADDRESSESPROVIDER() (common.Address, error) {
	return _DefaultReserveInterestRateStrategy.Contract.ADDRESSESPROVIDER(&_DefaultReserveInterestRateStrategy.CallOpts)
}

// ADDRESSESPROVIDER is a free data retrieval call binding the contract method 0x0542975c.
//
// Solidity: function ADDRESSES_PROVIDER() view returns(address)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategyCallerSession) ADDRESSESPROVIDER() (common.Address, error) {
	return _DefaultReserveInterestRateStrategy.Contract.ADDRESSESPROVIDER(&_DefaultReserveInterestRateStrategy.CallOpts)
}

// MAXEXCESSSTABLETOTOTALDEBTRATIO is a free data retrieval call binding the contract method 0xfe5fd698.
//
// Solidity: function MAX_EXCESS_STABLE_TO_TOTAL_DEBT_RATIO() view returns(uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategyCaller) MAXEXCESSSTABLETOTOTALDEBTRATIO(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DefaultReserveInterestRateStrategy.contract.Call(opts, &out, "MAX_EXCESS_STABLE_TO_TOTAL_DEBT_RATIO")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXEXCESSSTABLETOTOTALDEBTRATIO is a free data retrieval call binding the contract method 0xfe5fd698.
//
// Solidity: function MAX_EXCESS_STABLE_TO_TOTAL_DEBT_RATIO() view returns(uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategySession) MAXEXCESSSTABLETOTOTALDEBTRATIO() (*big.Int, error) {
	return _DefaultReserveInterestRateStrategy.Contract.MAXEXCESSSTABLETOTOTALDEBTRATIO(&_DefaultReserveInterestRateStrategy.CallOpts)
}

// MAXEXCESSSTABLETOTOTALDEBTRATIO is a free data retrieval call binding the contract method 0xfe5fd698.
//
// Solidity: function MAX_EXCESS_STABLE_TO_TOTAL_DEBT_RATIO() view returns(uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategyCallerSession) MAXEXCESSSTABLETOTOTALDEBTRATIO() (*big.Int, error) {
	return _DefaultReserveInterestRateStrategy.Contract.MAXEXCESSSTABLETOTOTALDEBTRATIO(&_DefaultReserveInterestRateStrategy.CallOpts)
}

// MAXEXCESSUSAGERATIO is a free data retrieval call binding the contract method 0xa9c622f8.
//
// Solidity: function MAX_EXCESS_USAGE_RATIO() view returns(uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategyCaller) MAXEXCESSUSAGERATIO(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DefaultReserveInterestRateStrategy.contract.Call(opts, &out, "MAX_EXCESS_USAGE_RATIO")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXEXCESSUSAGERATIO is a free data retrieval call binding the contract method 0xa9c622f8.
//
// Solidity: function MAX_EXCESS_USAGE_RATIO() view returns(uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategySession) MAXEXCESSUSAGERATIO() (*big.Int, error) {
	return _DefaultReserveInterestRateStrategy.Contract.MAXEXCESSUSAGERATIO(&_DefaultReserveInterestRateStrategy.CallOpts)
}

// MAXEXCESSUSAGERATIO is a free data retrieval call binding the contract method 0xa9c622f8.
//
// Solidity: function MAX_EXCESS_USAGE_RATIO() view returns(uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategyCallerSession) MAXEXCESSUSAGERATIO() (*big.Int, error) {
	return _DefaultReserveInterestRateStrategy.Contract.MAXEXCESSUSAGERATIO(&_DefaultReserveInterestRateStrategy.CallOpts)
}

// OPTIMALSTABLETOTOTALDEBTRATIO is a free data retrieval call binding the contract method 0x6fb92589.
//
// Solidity: function OPTIMAL_STABLE_TO_TOTAL_DEBT_RATIO() view returns(uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategyCaller) OPTIMALSTABLETOTOTALDEBTRATIO(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DefaultReserveInterestRateStrategy.contract.Call(opts, &out, "OPTIMAL_STABLE_TO_TOTAL_DEBT_RATIO")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// OPTIMALSTABLETOTOTALDEBTRATIO is a free data retrieval call binding the contract method 0x6fb92589.
//
// Solidity: function OPTIMAL_STABLE_TO_TOTAL_DEBT_RATIO() view returns(uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategySession) OPTIMALSTABLETOTOTALDEBTRATIO() (*big.Int, error) {
	return _DefaultReserveInterestRateStrategy.Contract.OPTIMALSTABLETOTOTALDEBTRATIO(&_DefaultReserveInterestRateStrategy.CallOpts)
}

// OPTIMALSTABLETOTOTALDEBTRATIO is a free data retrieval call binding the contract method 0x6fb92589.
//
// Solidity: function OPTIMAL_STABLE_TO_TOTAL_DEBT_RATIO() view returns(uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategyCallerSession) OPTIMALSTABLETOTOTALDEBTRATIO() (*big.Int, error) {
	return _DefaultReserveInterestRateStrategy.Contract.OPTIMALSTABLETOTOTALDEBTRATIO(&_DefaultReserveInterestRateStrategy.CallOpts)
}

// OPTIMALUSAGERATIO is a free data retrieval call binding the contract method 0x54c365c6.
//
// Solidity: function OPTIMAL_USAGE_RATIO() view returns(uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategyCaller) OPTIMALUSAGERATIO(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DefaultReserveInterestRateStrategy.contract.Call(opts, &out, "OPTIMAL_USAGE_RATIO")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// OPTIMALUSAGERATIO is a free data retrieval call binding the contract method 0x54c365c6.
//
// Solidity: function OPTIMAL_USAGE_RATIO() view returns(uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategySession) OPTIMALUSAGERATIO() (*big.Int, error) {
	return _DefaultReserveInterestRateStrategy.Contract.OPTIMALUSAGERATIO(&_DefaultReserveInterestRateStrategy.CallOpts)
}

// OPTIMALUSAGERATIO is a free data retrieval call binding the contract method 0x54c365c6.
//
// Solidity: function OPTIMAL_USAGE_RATIO() view returns(uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategyCallerSession) OPTIMALUSAGERATIO() (*big.Int, error) {
	return _DefaultReserveInterestRateStrategy.Contract.OPTIMALUSAGERATIO(&_DefaultReserveInterestRateStrategy.CallOpts)
}

// CalculateInterestRates is a free data retrieval call binding the contract method 0xa5898709.
//
// Solidity: function calculateInterestRates((uint256,uint256,uint256,uint256,uint256,uint256,uint256,address,address) params) view returns(uint256, uint256, uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategyCaller) CalculateInterestRates(opts *bind.CallOpts, params DataTypesCalculateInterestRatesParams) (*big.Int, *big.Int, *big.Int, error) {
	var out []interface{}
	err := _DefaultReserveInterestRateStrategy.contract.Call(opts, &out, "calculateInterestRates", params)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	out2 := *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return out0, out1, out2, err

}

// CalculateInterestRates is a free data retrieval call binding the contract method 0xa5898709.
//
// Solidity: function calculateInterestRates((uint256,uint256,uint256,uint256,uint256,uint256,uint256,address,address) params) view returns(uint256, uint256, uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategySession) CalculateInterestRates(params DataTypesCalculateInterestRatesParams) (*big.Int, *big.Int, *big.Int, error) {
	return _DefaultReserveInterestRateStrategy.Contract.CalculateInterestRates(&_DefaultReserveInterestRateStrategy.CallOpts, params)
}

// CalculateInterestRates is a free data retrieval call binding the contract method 0xa5898709.
//
// Solidity: function calculateInterestRates((uint256,uint256,uint256,uint256,uint256,uint256,uint256,address,address) params) view returns(uint256, uint256, uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategyCallerSession) CalculateInterestRates(params DataTypesCalculateInterestRatesParams) (*big.Int, *big.Int, *big.Int, error) {
	return _DefaultReserveInterestRateStrategy.Contract.CalculateInterestRates(&_DefaultReserveInterestRateStrategy.CallOpts, params)
}

// GetBaseStableBorrowRate is a free data retrieval call binding the contract method 0xacd78686.
//
// Solidity: function getBaseStableBorrowRate() view returns(uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategyCaller) GetBaseStableBorrowRate(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DefaultReserveInterestRateStrategy.contract.Call(opts, &out, "getBaseStableBorrowRate")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBaseStableBorrowRate is a free data retrieval call binding the contract method 0xacd78686.
//
// Solidity: function getBaseStableBorrowRate() view returns(uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategySession) GetBaseStableBorrowRate() (*big.Int, error) {
	return _DefaultReserveInterestRateStrategy.Contract.GetBaseStableBorrowRate(&_DefaultReserveInterestRateStrategy.CallOpts)
}

// GetBaseStableBorrowRate is a free data retrieval call binding the contract method 0xacd78686.
//
// Solidity: function getBaseStableBorrowRate() view returns(uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategyCallerSession) GetBaseStableBorrowRate() (*big.Int, error) {
	return _DefaultReserveInterestRateStrategy.Contract.GetBaseStableBorrowRate(&_DefaultReserveInterestRateStrategy.CallOpts)
}

// GetBaseVariableBorrowRate is a free data retrieval call binding the contract method 0x34762ca5.
//
// Solidity: function getBaseVariableBorrowRate() view returns(uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategyCaller) GetBaseVariableBorrowRate(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DefaultReserveInterestRateStrategy.contract.Call(opts, &out, "getBaseVariableBorrowRate")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBaseVariableBorrowRate is a free data retrieval call binding the contract method 0x34762ca5.
//
// Solidity: function getBaseVariableBorrowRate() view returns(uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategySession) GetBaseVariableBorrowRate() (*big.Int, error) {
	return _DefaultReserveInterestRateStrategy.Contract.GetBaseVariableBorrowRate(&_DefaultReserveInterestRateStrategy.CallOpts)
}

// GetBaseVariableBorrowRate is a free data retrieval call binding the contract method 0x34762ca5.
//
// Solidity: function getBaseVariableBorrowRate() view returns(uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategyCallerSession) GetBaseVariableBorrowRate() (*big.Int, error) {
	return _DefaultReserveInterestRateStrategy.Contract.GetBaseVariableBorrowRate(&_DefaultReserveInterestRateStrategy.CallOpts)
}

// GetMaxVariableBorrowRate is a free data retrieval call binding the contract method 0x80031e37.
//
// Solidity: function getMaxVariableBorrowRate() view returns(uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategyCaller) GetMaxVariableBorrowRate(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DefaultReserveInterestRateStrategy.contract.Call(opts, &out, "getMaxVariableBorrowRate")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMaxVariableBorrowRate is a free data retrieval call binding the contract method 0x80031e37.
//
// Solidity: function getMaxVariableBorrowRate() view returns(uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategySession) GetMaxVariableBorrowRate() (*big.Int, error) {
	return _DefaultReserveInterestRateStrategy.Contract.GetMaxVariableBorrowRate(&_DefaultReserveInterestRateStrategy.CallOpts)
}

// GetMaxVariableBorrowRate is a free data retrieval call binding the contract method 0x80031e37.
//
// Solidity: function getMaxVariableBorrowRate() view returns(uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategyCallerSession) GetMaxVariableBorrowRate() (*big.Int, error) {
	return _DefaultReserveInterestRateStrategy.Contract.GetMaxVariableBorrowRate(&_DefaultReserveInterestRateStrategy.CallOpts)
}

// GetStableRateExcessOffset is a free data retrieval call binding the contract method 0xbc626908.
//
// Solidity: function getStableRateExcessOffset() view returns(uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategyCaller) GetStableRateExcessOffset(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DefaultReserveInterestRateStrategy.contract.Call(opts, &out, "getStableRateExcessOffset")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetStableRateExcessOffset is a free data retrieval call binding the contract method 0xbc626908.
//
// Solidity: function getStableRateExcessOffset() view returns(uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategySession) GetStableRateExcessOffset() (*big.Int, error) {
	return _DefaultReserveInterestRateStrategy.Contract.GetStableRateExcessOffset(&_DefaultReserveInterestRateStrategy.CallOpts)
}

// GetStableRateExcessOffset is a free data retrieval call binding the contract method 0xbc626908.
//
// Solidity: function getStableRateExcessOffset() view returns(uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategyCallerSession) GetStableRateExcessOffset() (*big.Int, error) {
	return _DefaultReserveInterestRateStrategy.Contract.GetStableRateExcessOffset(&_DefaultReserveInterestRateStrategy.CallOpts)
}

// GetStableRateSlope1 is a free data retrieval call binding the contract method 0xd5cd7391.
//
// Solidity: function getStableRateSlope1() view returns(uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategyCaller) GetStableRateSlope1(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DefaultReserveInterestRateStrategy.contract.Call(opts, &out, "getStableRateSlope1")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetStableRateSlope1 is a free data retrieval call binding the contract method 0xd5cd7391.
//
// Solidity: function getStableRateSlope1() view returns(uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategySession) GetStableRateSlope1() (*big.Int, error) {
	return _DefaultReserveInterestRateStrategy.Contract.GetStableRateSlope1(&_DefaultReserveInterestRateStrategy.CallOpts)
}

// GetStableRateSlope1 is a free data retrieval call binding the contract method 0xd5cd7391.
//
// Solidity: function getStableRateSlope1() view returns(uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategyCallerSession) GetStableRateSlope1() (*big.Int, error) {
	return _DefaultReserveInterestRateStrategy.Contract.GetStableRateSlope1(&_DefaultReserveInterestRateStrategy.CallOpts)
}

// GetStableRateSlope2 is a free data retrieval call binding the contract method 0x14e32da4.
//
// Solidity: function getStableRateSlope2() view returns(uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategyCaller) GetStableRateSlope2(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DefaultReserveInterestRateStrategy.contract.Call(opts, &out, "getStableRateSlope2")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetStableRateSlope2 is a free data retrieval call binding the contract method 0x14e32da4.
//
// Solidity: function getStableRateSlope2() view returns(uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategySession) GetStableRateSlope2() (*big.Int, error) {
	return _DefaultReserveInterestRateStrategy.Contract.GetStableRateSlope2(&_DefaultReserveInterestRateStrategy.CallOpts)
}

// GetStableRateSlope2 is a free data retrieval call binding the contract method 0x14e32da4.
//
// Solidity: function getStableRateSlope2() view returns(uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategyCallerSession) GetStableRateSlope2() (*big.Int, error) {
	return _DefaultReserveInterestRateStrategy.Contract.GetStableRateSlope2(&_DefaultReserveInterestRateStrategy.CallOpts)
}

// GetVariableRateSlope1 is a free data retrieval call binding the contract method 0x0b3429a2.
//
// Solidity: function getVariableRateSlope1() view returns(uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategyCaller) GetVariableRateSlope1(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DefaultReserveInterestRateStrategy.contract.Call(opts, &out, "getVariableRateSlope1")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetVariableRateSlope1 is a free data retrieval call binding the contract method 0x0b3429a2.
//
// Solidity: function getVariableRateSlope1() view returns(uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategySession) GetVariableRateSlope1() (*big.Int, error) {
	return _DefaultReserveInterestRateStrategy.Contract.GetVariableRateSlope1(&_DefaultReserveInterestRateStrategy.CallOpts)
}

// GetVariableRateSlope1 is a free data retrieval call binding the contract method 0x0b3429a2.
//
// Solidity: function getVariableRateSlope1() view returns(uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategyCallerSession) GetVariableRateSlope1() (*big.Int, error) {
	return _DefaultReserveInterestRateStrategy.Contract.GetVariableRateSlope1(&_DefaultReserveInterestRateStrategy.CallOpts)
}

// GetVariableRateSlope2 is a free data retrieval call binding the contract method 0xf4202409.
//
// Solidity: function getVariableRateSlope2() view returns(uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategyCaller) GetVariableRateSlope2(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DefaultReserveInterestRateStrategy.contract.Call(opts, &out, "getVariableRateSlope2")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetVariableRateSlope2 is a free data retrieval call binding the contract method 0xf4202409.
//
// Solidity: function getVariableRateSlope2() view returns(uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategySession) GetVariableRateSlope2() (*big.Int, error) {
	return _DefaultReserveInterestRateStrategy.Contract.GetVariableRateSlope2(&_DefaultReserveInterestRateStrategy.CallOpts)
}

// GetVariableRateSlope2 is a free data retrieval call binding the contract method 0xf4202409.
//
// Solidity: function getVariableRateSlope2() view returns(uint256)
func (_DefaultReserveInterestRateStrategy *DefaultReserveInterestRateStrategyCallerSession) GetVariableRateSlope2() (*big.Int, error) {
	return _DefaultReserveInterestRateStrategy.Contract.GetVariableRateSlope2(&_DefaultReserveInterestRateStrategy.CallOpts)
}
