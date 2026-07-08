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

// EModeLogicMetaData contains all meta data concerning the EModeLogic contract.
var EModeLogicMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"categoryId\",\"type\":\"uint8\"}],\"name\":\"UserEModeSet\",\"type\":\"event\"}]",
	Bin: "0x61146e61003a600b82828239805160001a60731461002d57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600436106100355760003560e01c80635d5dc3131461003a575b600080fd5b81801561004657600080fd5b5061005a610055366004611192565b61005c565b005b60408051602081018252835481528251918301516100809289928992899290610145565b336000908152602084905260409081902080549183015160ff9081167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff008416179091551680156100fe576100fb87878786604051806020016040529081600082015481525050338760400151886000015189602001516102e0565b50505b604080830151905160ff909116815233907fd728da875fc88944cbf17638bcbe4af0eedaef63becd1d1c57cc097eb4608d849060200160405180910390a250505050505050565b60ff81161580610170575060ff811660009081526020859052604090205462010000900461ffff1615155b6040518060400160405280600281526020017f3538000000000000000000000000000000000000000000000000000000000000815250906101e7576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016101de91906112a7565b60405180910390fd5b5082516101f3576102d8565b60ff8116156102d85760005b828110156102d65761021184826103db565b156102ce576000818152602087815260408083205473ffffffffffffffffffffffffffffffffffffffff168352898252918290208251918201909252905480825260ff8481169160a81c16146040518060400160405280600281526020017f3538000000000000000000000000000000000000000000000000000000000000815250906102cb576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016101de91906112a7565b50505b6001016101ff565b505b505050505050565b6000806000806103478c8c8c6040518060a001604052808e81526020018b81526020018d73ffffffffffffffffffffffffffffffffffffffff1681526020018a73ffffffffffffffffffffffffffffffffffffffff1681526020018c60ff1681525061045d565b9550955050505050670de0b6b3a76400008210156040518060400160405280600281526020017f3335000000000000000000000000000000000000000000000000000000000000815250906103c9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016101de91906112a7565b50909b909a5098505050505050505050565b60408051808201909152600281527f373400000000000000000000000000000000000000000000000000000000000060208201526000906080831061044d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016101de91906112a7565b50509051600191821b1c16151590565b6000806000806000806104738760000151511590565b156104af5750600094508493508392508291507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9050816109ba565b61055e60405180610260016040528060008152602001600081526020016000815260200160008152602001600081526020016000815260200160008152602001600081526020016000815260200160008152602001600081526020016000815260200160008152602001600081526020016000815260200160008152602001600073ffffffffffffffffffffffffffffffffffffffff1681526020016000151581526020016000151581525090565b608088015160ff16156105a357608088015160ff16600090815260208a905260409020606089015161059091906109c7565b6101808401526101c08301526101a08201525b87602001518160c0015110156108c25760c081015188516105c391610aa6565b6105d75760c08101805160010190526105a3565b60c0810151600090815260208b9052604090205473ffffffffffffffffffffffffffffffffffffffff16610200820181905261061d5760c08101805160010190526105a3565b61020081015173ffffffffffffffffffffffffffffffffffffffff16600090815260208c8152604091829020825180830190935280549283905260ff60a884901c81166101e0860152603084901c166060850181905261ffff601085901c811660a08701529093166080850152600a9290920a90830152610180820151158015906106b35750816101e00151896080015160ff16145b6107575760608901516102008301516040517fb3596f0700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff918216600482015291169063b3596f0790602401602060405180830381865afa15801561072e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610752919061131a565b61075e565b8161018001515b825260a08201511580159061077e575060c0820151895161077e91610b2b565b1561086e5761079b89604001518284600001518560200151610baf565b60408301819052610100830180516107b4908390611362565b90525060808901516101e08301516107cf9160ff1690610c8e565b1515610240830152608082015115610825578161024001516107f55781608001516107fc565b816101a001515b826040015161080b919061137a565b826101400181815161081d9190611362565b90525061082e565b60016102208301525b816102400151610842578160a00151610849565b816101c001515b8260400151610858919061137a565b826101600181815161086a9190611362565b9052505b60c0820151895161087e916103db565b156108b15761089b89604001518284600001518560200151610ca5565b82610120018181516108ad9190611362565b9052505b5060c08101805160010190526105a3565b6101008101516108d35760006108ee565b806101000151816101400151816108ec576108ec6113b7565b045b610140820152610100810151610905576000610920565b8061010001518161016001518161091e5761091e6113b7565b045b610160820152610120810151156109625761095d816101200151610957836101600151846101000151610e2590919063ffffffff16565b90610e68565b610984565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff5b60e0820181905261010082015161012083015161014084015161016085015161022090950151929a509098509650919450925090505b9499939850945094509450565b81546000908190819081906601000000000000900473ffffffffffffffffffffffffffffffffffffffff168015610a8b576040517fb3596f0700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff828116600483015287169063b3596f0790602401602060405180830381865afa158015610a64573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a88919061131a565b91505b50945461ffff80821697620100009092041695945092505050565b60408051808201909152600281527f3734000000000000000000000000000000000000000000000000000000000000602082015260009060808310610b18576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016101de91906112a7565b5050905160019190911b1c600316151590565b60408051808201909152600281527f3734000000000000000000000000000000000000000000000000000000000000602082015260009060808310610b9d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016101de91906112a7565b50509051600191821b82011c16151590565b600080610bbb85610e9f565b6004868101546040517f1da24f3e00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8a8116938201939093529293506000928792610c67928692911690631da24f3e90602401602060405180830381865afa158015610c3d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c61919061131a565b90610f23565b610c71919061137a565b9050838181610c8257610c826113b7565b04979650505050505050565b60008215801590610c9e57508282145b9392505050565b60068301546040517f1da24f3e00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff86811660048301526000928392911690631da24f3e90602401602060405180830381865afa158015610d1b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d3f919061131a565b90508015610d5d57610d5a610d5386610f7a565b8290610f23565b90505b60058501546040517f70a0823100000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8881166004830152909116906370a0823190602401602060405180830381865afa158015610dcf573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610df3919061131a565b610dfd9082611362565b9050610e09818561137a565b9050828181610e1a57610e1a6113b7565b049695505050505050565b600081157fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffec7783900484111517610e5a57600080fd5b506127109102611388010490565b60008115670de0b6b3a764000060028404190484111715610e8857600080fd5b50670de0b6b3a76400009190910260028204010490565b6003810154600090700100000000000000000000000000000000900464ffffffffff1642811415610ee5575050600101546fffffffffffffffffffffffffffffffff1690565b6001830154610c9e906fffffffffffffffffffffffffffffffff80821691610c61917001000000000000000000000000000000009091041684610ffe565b600081157ffffffffffffffffffffffffffffffffffffffffffe6268e1b017bfe18bffffff83900484111517610f5857600080fd5b506b033b2e3c9fd0803ce800000091026b019d971e4fe8401e74000000010490565b6003810154600090700100000000000000000000000000000000900464ffffffffff1642811415610fc0575050600201546fffffffffffffffffffffffffffffffff1690565b6002830154610c9e906fffffffffffffffffffffffffffffffff80821691610c61917001000000000000000000000000000000009091041684611043565b60008061101264ffffffffff8416426113e6565b61101c908561137a565b6301e133809004905061103b816b033b2e3c9fd0803ce8000000611362565b949350505050565b6000610c9e83834260008061105f64ffffffffff8516846113e6565b90508061107b576b033b2e3c9fd0803ce8000000915050610c9e565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff810160008080600285116110b15760006110b6565b600285035b925066038882915c40006110ca8a80610f23565b816110d7576110d76113b7565b0491506301e133806110e9838b610f23565b816110f6576110f66113b7565b049050600082611106868861137a565b611110919061137a565b60029004905060008285611124888a61137a565b61112e919061137a565b611138919061137a565b60069004905080826301e1338061114f8a8f61137a565b61115991906113fd565b61116f906b033b2e3c9fd0803ce8000000611362565b6111799190611362565b6111839190611362565b9b9a5050505050505050505050565b6000806000806000808688036101008112156111ad57600080fd5b873596506020880135955060408801359450606088013593506080880135925060607fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60820112156111fd57600080fd5b506040516060810181811067ffffffffffffffff82111715611248577f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405260a0880135815260c088013573ffffffffffffffffffffffffffffffffffffffff8116811461127957600080fd5b602082015260e088013560ff8116811461129257600080fd5b80604083015250809150509295509295509295565b600060208083528351808285015260005b818110156112d4578581018301518582016040015282016112b8565b818111156112e6576000604083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016929092016040019392505050565b60006020828403121561132c57600080fd5b5051919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000821982111561137557611375611333565b500190565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156113b2576113b2611333565b500290565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000828210156113f8576113f8611333565b500390565b600082611433577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b50049056fea26469706673582212202c62f8da245632d62b473e996662e92326dfe4623b6e5bc5fe41492f6033e41764736f6c634300080a0033",
}

// EModeLogicABI is the input ABI used to generate the binding from.
// Deprecated: Use EModeLogicMetaData.ABI instead.
var EModeLogicABI = EModeLogicMetaData.ABI

// EModeLogicBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use EModeLogicMetaData.Bin instead.
var EModeLogicBin = EModeLogicMetaData.Bin

// DeployEModeLogic deploys a new Ethereum contract, binding an instance of EModeLogic to it.
func DeployEModeLogic(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *EModeLogic, error) {
	parsed, err := EModeLogicMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(EModeLogicBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &EModeLogic{EModeLogicCaller: EModeLogicCaller{contract: contract}, EModeLogicTransactor: EModeLogicTransactor{contract: contract}, EModeLogicFilterer: EModeLogicFilterer{contract: contract}}, nil
}

// EModeLogic is an auto generated Go binding around an Ethereum contract.
type EModeLogic struct {
	EModeLogicCaller     // Read-only binding to the contract
	EModeLogicTransactor // Write-only binding to the contract
	EModeLogicFilterer   // Log filterer for contract events
}

// EModeLogicCaller is an auto generated read-only Go binding around an Ethereum contract.
type EModeLogicCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EModeLogicTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EModeLogicTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EModeLogicFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EModeLogicFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EModeLogicSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EModeLogicSession struct {
	Contract     *EModeLogic       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EModeLogicCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EModeLogicCallerSession struct {
	Contract *EModeLogicCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// EModeLogicTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EModeLogicTransactorSession struct {
	Contract     *EModeLogicTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// EModeLogicRaw is an auto generated low-level Go binding around an Ethereum contract.
type EModeLogicRaw struct {
	Contract *EModeLogic // Generic contract binding to access the raw methods on
}

// EModeLogicCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EModeLogicCallerRaw struct {
	Contract *EModeLogicCaller // Generic read-only contract binding to access the raw methods on
}

// EModeLogicTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EModeLogicTransactorRaw struct {
	Contract *EModeLogicTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEModeLogic creates a new instance of EModeLogic, bound to a specific deployed contract.
func NewEModeLogic(address common.Address, backend bind.ContractBackend) (*EModeLogic, error) {
	contract, err := bindEModeLogic(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EModeLogic{EModeLogicCaller: EModeLogicCaller{contract: contract}, EModeLogicTransactor: EModeLogicTransactor{contract: contract}, EModeLogicFilterer: EModeLogicFilterer{contract: contract}}, nil
}

// NewEModeLogicCaller creates a new read-only instance of EModeLogic, bound to a specific deployed contract.
func NewEModeLogicCaller(address common.Address, caller bind.ContractCaller) (*EModeLogicCaller, error) {
	contract, err := bindEModeLogic(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EModeLogicCaller{contract: contract}, nil
}

// NewEModeLogicTransactor creates a new write-only instance of EModeLogic, bound to a specific deployed contract.
func NewEModeLogicTransactor(address common.Address, transactor bind.ContractTransactor) (*EModeLogicTransactor, error) {
	contract, err := bindEModeLogic(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EModeLogicTransactor{contract: contract}, nil
}

// NewEModeLogicFilterer creates a new log filterer instance of EModeLogic, bound to a specific deployed contract.
func NewEModeLogicFilterer(address common.Address, filterer bind.ContractFilterer) (*EModeLogicFilterer, error) {
	contract, err := bindEModeLogic(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EModeLogicFilterer{contract: contract}, nil
}

// bindEModeLogic binds a generic wrapper to an already deployed contract.
func bindEModeLogic(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := EModeLogicMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EModeLogic *EModeLogicRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EModeLogic.Contract.EModeLogicCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EModeLogic *EModeLogicRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EModeLogic.Contract.EModeLogicTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EModeLogic *EModeLogicRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EModeLogic.Contract.EModeLogicTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EModeLogic *EModeLogicCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EModeLogic.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EModeLogic *EModeLogicTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EModeLogic.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EModeLogic *EModeLogicTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EModeLogic.Contract.contract.Transact(opts, method, params...)
}

// EModeLogicUserEModeSetIterator is returned from FilterUserEModeSet and is used to iterate over the raw logs and unpacked data for UserEModeSet events raised by the EModeLogic contract.
type EModeLogicUserEModeSetIterator struct {
	Event *EModeLogicUserEModeSet // Event containing the contract specifics and raw log

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
func (it *EModeLogicUserEModeSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EModeLogicUserEModeSet)
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
		it.Event = new(EModeLogicUserEModeSet)
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
func (it *EModeLogicUserEModeSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EModeLogicUserEModeSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EModeLogicUserEModeSet represents a UserEModeSet event raised by the EModeLogic contract.
type EModeLogicUserEModeSet struct {
	User       common.Address
	CategoryId uint8
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterUserEModeSet is a free log retrieval operation binding the contract event 0xd728da875fc88944cbf17638bcbe4af0eedaef63becd1d1c57cc097eb4608d84.
//
// Solidity: event UserEModeSet(address indexed user, uint8 categoryId)
func (_EModeLogic *EModeLogicFilterer) FilterUserEModeSet(opts *bind.FilterOpts, user []common.Address) (*EModeLogicUserEModeSetIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _EModeLogic.contract.FilterLogs(opts, "UserEModeSet", userRule)
	if err != nil {
		return nil, err
	}
	return &EModeLogicUserEModeSetIterator{contract: _EModeLogic.contract, event: "UserEModeSet", logs: logs, sub: sub}, nil
}

// WatchUserEModeSet is a free log subscription operation binding the contract event 0xd728da875fc88944cbf17638bcbe4af0eedaef63becd1d1c57cc097eb4608d84.
//
// Solidity: event UserEModeSet(address indexed user, uint8 categoryId)
func (_EModeLogic *EModeLogicFilterer) WatchUserEModeSet(opts *bind.WatchOpts, sink chan<- *EModeLogicUserEModeSet, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _EModeLogic.contract.WatchLogs(opts, "UserEModeSet", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EModeLogicUserEModeSet)
				if err := _EModeLogic.contract.UnpackLog(event, "UserEModeSet", log); err != nil {
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

// ParseUserEModeSet is a log parse operation binding the contract event 0xd728da875fc88944cbf17638bcbe4af0eedaef63becd1d1c57cc097eb4608d84.
//
// Solidity: event UserEModeSet(address indexed user, uint8 categoryId)
func (_EModeLogic *EModeLogicFilterer) ParseUserEModeSet(log types.Log) (*EModeLogicUserEModeSet, error) {
	event := new(EModeLogicUserEModeSet)
	if err := _EModeLogic.contract.UnpackLog(event, "UserEModeSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
