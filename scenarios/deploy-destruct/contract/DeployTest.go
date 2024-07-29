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

// ContractMetaData contains all meta data concerning the Contract contract.
var ContractMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"}],\"name\":\"TestSeed\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"childAddresses\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"childCode\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"childIdx\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"count\",\"type\":\"uint256\"}],\"name\":\"clean\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"destroy\",\"type\":\"bool\"}],\"name\":\"notify\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"}],\"name\":\"test\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x608060405260015f55600180553480156016575f80fd5b506118bb806100245f395ff3fe60806040526004361061006d575f3560e01c806361bc221a1161004c57806361bc221a146100cf578063674ef0fa146100e35780639c47ed9f14610102578063d52c3bd21461014e575f80fd5b8062fd4ee9146100715780630e7bfac51461009957806329e99f07146100ba575b5f80fd5b34801561007c575f80fd5b5061008660015481565b6040519081526020015b60405180910390f35b3480156100a4575f80fd5b506100ad61016d565b60405161009091906106b8565b6100cd6100c83660046106ed565b610197565b005b3480156100da575f80fd5b506100865f5481565b3480156100ee575f80fd5b506100cd6100fd3660046106ed565b61041f565b34801561010d575f80fd5b5061013661011c3660046106ed565b60026020525f90815260409020546001600160a01b031681565b6040516001600160a01b039091168152602001610090565b348015610159575f80fd5b506100cd610168366004610704565b6104d4565b60606040518060200161017f906106ab565b601f1982820381018352601f90910116604052919050565b805f036101b757620186a06101aa6105bc565b6101b4919061075e565b90505b6040518181527f47ff30cfcf5248fde6dac21a1944ca5982927073500951ddde02f1d68d63c3da9060200160405180910390a15f6101f7825f600161060b565b9050601e61020660328361075e565b101561027d575f61021561016d565b60408051602081018590525f9181019190915230606082015260800160408051601f198184030181529082905261024f9291602001610788565b60405160208183030381529060405290505f61026b4783610642565b9050610277815f6104d4565b50505050565b602861028a60328361075e565b1015610355575f5b600a811015610350576001805411156103485760015f81905260026020527fe90b7bceb6e7df5418fb78d8ee546e97c83a08bbccc01a0644d599ccd2a7c2e0546001600160a01b0316908190639d118770906102f09087908661060b565b6040518263ffffffff1660e01b815260040161030e91815260200190565b5f604051808303815f87803b158015610325575f80fd5b505af1158015610337573d5f803e3d5ffd5b505050506103468160016104d4565b505b600101610292565b505050565b4760015b600154811015610277575f818152600260205260408120546001600160a01b03169061038686828561060b565b90505f61039460508361075e565b61039f6064876107b8565b6103a991906107cb565b90506103b581866107e2565b604051631584b49360e11b8152600481018490529095506001600160a01b03841690632b0969269083906024015f604051808303818588803b1580156103f9575f80fd5b505af115801561040b573d5f803e3d5ffd5b505060019096019550610359945050505050565b80156104d1576001805411156104bf5760015f5260026020527fe90b7bceb6e7df5418fb78d8ee546e97c83a08bbccc01a0644d599ccd2a7c2e0546040516309d1187760e41b8152600560048201526001600160a01b03909116908190639d118770906024015f604051808303815f87803b15801561049c575f80fd5b505af11580156104ae573d5f803e3d5ffd5b505050506104bd8160016104d4565b505b806104c9816107f5565b91505061041f565b50565b8015610569576001600160a01b0382165f908152600360205260409020548015610350576001600160a01b0383165f9081526003602052604081208190556001805491610520836107f5565b90915550506001545f9081526002602090815260408083205484845281842080546001600160a01b0319166001600160a01b039092169182179055835260039091529020555050565b600180545f90815260026020908152604080832080546001600160a01b0319166001600160a01b03881690811790915584549084526003909252822081905591906105b38361080a565b91905055505050565b5f805481806105ca8361080a565b90915550505f546040805144602082015242918101919091526060810191909152608001604051602081830303815290604052805190602001205f1c905090565b6040805160208082019590955280820193909352606080840192909252805180840390920182526080909201909152805191012090565b5f8082516020840185f09050803b610658575f80fd5b6001600160a01b0381166106a25760405162461bcd60e51b815260206004820152600d60248201526c18dc99585d194819985a5b1959609a1b604482015260640160405180910390fd5b90505b92915050565b6110638061082383390190565b602081525f82518060208401528060208501604085015e5f604082850101526040601f19601f83011684010191505092915050565b5f602082840312156106fd575f80fd5b5035919050565b5f8060408385031215610715575f80fd5b82356001600160a01b038116811461072b575f80fd5b91506020830135801515811461073f575f80fd5b809150509250929050565b634e487b7160e01b5f52601260045260245ffd5b5f8261076c5761076c61074a565b500690565b5f81518060208401855e5f93019283525090919050565b5f61079c6107968386610771565b84610771565b949350505050565b634e487b7160e01b5f52601160045260245ffd5b5f826107c6576107c661074a565b500490565b80820281158282048414176106a5576106a56107a4565b818103818111156106a5576106a56107a4565b5f81610803576108036107a4565b505f190190565b5f6001820161081b5761081b6107a4565b506001019056fe608060405260405161106338038061106383398101604081905261002291610586565b60408051848152602081018490527f06ed9ff5e25ad09aa503577a10190fb63ca2ae4e7ea00bc39a59d3daa8bcbec3910160405180910390a15f80546001600160a01b0319166001600160a01b03831617905560018390556100838261008b565b505050610797565b5f61009f600154835f61048560201b60201c565b9050476004831080156100bc5750603c6100ba6064846105dc565b105b15610480575f805f9054906101000a90046001600160a01b03166001600160a01b0316630e7bfac56040518163ffffffff1660e01b81526004015f60405180830381865afa158015610110573d5f803e3d5ffd5b505050506040513d5f823e601f3d908101601f191682016040526101379190810190610603565b90505f5b6101466003856105dc565b6101519060016106c7565b81101561047d57621e84805a1061047d575f610176600154878461048560201b60201c565b90505f83826101868960016106c7565b5f546040805160208101949094528301919091526001600160a01b0316606082015260800160408051601f19818403018152908290526101c992916020016106f1565b604051602081830303815290604052905061022160015460016101ec91906106c7565b604080516020808201939093528082018b90526060808201889052825180830390910181526080909101909152805191012090565b91505f61022f6050846105dc565b61023a606488610705565b6102449190610718565b9050610250818761072f565b95505f60326102606064866105dc565b10156102775761027082846104bc565b905061028f565b61028c826102868760016106c7565b85610521565b90505b6001600160a01b0381166102a6575050505061047d565b6102ed60015460026102b891906106c7565b604080516020808201939093528082018d905260608082018a9052825180830390910181526080909101909152805191012090565b935060286102fc6064866105dc565b10156103bb576040516309d1187760e41b8152600481018590526001600160a01b03821690639d118770906024015f604051808303815f87803b158015610341575f80fd5b505af1925050508015610352575060015b6103b6573d80801561037f576040519150601f19603f3d011682016040523d82523d5f602084013e610384565b606091505b505f805160206110438339815191526001546001836040516103a893929190610770565b60405180910390a15061046d565b61046d565b5f8054604051636a961de960e11b81526001600160a01b038481166004830152602482019390935291169063d52c3bd2906044015f604051808303815f87803b158015610406575f80fd5b505af1925050508015610417575060015b61046d573d808015610444576040519150601f19603f3d011682016040523d82523d5f602084013e610449565b606091505b505f805160206110438339815191526001546002836040516103a893929190610770565b50506001909201915061013b9050565b50505b505050565b6040805160208082019590955280820193909352606080840192909252805180840390920182526080909201909152805191012090565b5f8082516020840185f09050803b6104d2575f80fd5b6001600160a01b03811661051857600154604080515f815260208101918290525f805160206110438339815191529261050f929091600391610770565b60405180910390a15b90505b92915050565b5f808383516020850187f59050803b610538575f80fd5b6001600160a01b03811661057e57600154604080515f815260208101918290525f8051602061104383398151915292610575929091600391610770565b60405180910390a15b949350505050565b5f805f60608486031215610598575f80fd5b83516020850151604086015191945092506001600160a01b03811681146105bd575f80fd5b809150509250925092565b634e487b7160e01b5f52601260045260245ffd5b5f826105ea576105ea6105c8565b500690565b634e487b7160e01b5f52604160045260245ffd5b5f60208284031215610613575f80fd5b81516001600160401b03811115610628575f80fd5b8201601f81018413610638575f80fd5b80516001600160401b03811115610651576106516105ef565b604051601f8201601f19908116603f011681016001600160401b038111828210171561067f5761067f6105ef565b604052818152828201602001861015610696575f80fd5b8160208401602083015e5f91810160200191909152949350505050565b634e487b7160e01b5f52601160045260245ffd5b8082018082111561051b5761051b6106b3565b5f81518060208401855e5f93019283525090919050565b5f61057e6106ff83866106da565b846106da565b5f82610713576107136105c8565b500490565b808202811582820484141761051b5761051b6106b3565b8181038181111561051b5761051b6106b3565b5f81518084528060208401602086015e5f602082860101526020601f19601f83011685010191505092915050565b838152826020820152606060408201525f61078e6060830184610742565b95945050505050565b61089f806107a45f395ff3fe60806040526004361061003e575f3560e01c80632b096926146100425780635a34d356146100575780637469d0681461007f5780639d118770146100b5575b5f80fd5b610055610050366004610660565b6100d4565b005b348015610062575f80fd5b5061006c60015481565b6040519081526020015b60405180910390f35b34801561008a575f80fd5b505f5461009d906001600160a01b031681565b6040516001600160a01b039091168152602001610076565b3480156100c0575f80fd5b506100556100cf366004610660565b610466565b5f6100e2600154835f61055f565b9050476004831080156100ff5750603c6100fd60648461068b565b105b15610461575f805f9054906101000a90046001600160a01b03166001600160a01b0316630e7bfac56040518163ffffffff1660e01b81526004015f60405180830381865afa158015610153573d5f803e3d5ffd5b505050506040513d5f823e601f3d908101601f1916820160405261017a91908101906106b2565b90505f5b61018960038561068b565b610194906001610779565b81101561045e57621e84805a1061045e575f6101b3600154878461055f565b90505f83826101c3896001610779565b5f546040805160208101949094528301919091526001600160a01b0316606082015260800160408051601f198184030181529082905261020692916020016107a3565b604051602081830303815290604052905061023060015460016102299190610779565b888561055f565b91505f61023e60508461068b565b6102496064886107b7565b61025391906107ca565b905061025f81876107e1565b95505f603261026f60648661068b565b10156102865761027f8284610596565b905061029e565b61029b82610295876001610779565b856105fb565b90505b6001600160a01b0381166102b5575050505061045e565b6102ce60015460026102c79190610779565b8a8761055f565b935060286102dd60648661068b565b101561039c576040516309d1187760e41b8152600481018590526001600160a01b03821690639d118770906024015f604051808303815f87803b158015610322575f80fd5b505af1925050508015610333575060015b610397573d808015610360576040519150601f19603f3d011682016040523d82523d5f602084013e610365565b606091505b505f8051602061084a83398151915260015460018360405161038993929190610822565b60405180910390a15061044e565b61044e565b5f8054604051636a961de960e11b81526001600160a01b038481166004830152602482019390935291169063d52c3bd2906044015f604051808303815f87803b1580156103e7575f80fd5b505af19250505080156103f8575060015b61044e573d808015610425576040519150601f19603f3d011682016040523d82523d5f602084013e61042a565b606091505b505f8051602061084a83398151915260015460028360405161038993929190610822565b50506001909201915061017e9050565b50505b505050565b5f8061047360068461068b565b9050805f0361048457339150610512565b806001036104b8576040805160208101859052016040516020818303038152906040528051906020012060601c9150610512565b806002036104c857329150610512565b806003036104d857309150610512565b806004036104fc577349e0fd3800c117357057534e30c5b5115c6734889150610512565b80600503610512575f546001600160a01b031691505b604080518481526001600160a01b03841660208201527f77b29af0d4f525b395d176ebc2772a5ffe882cd17c4b934ee2fd3b773ba41040910160405180910390a1816001600160a01b0316ff5b6040805160208082019590955280820193909352606080840192909252805180840390920182526080909201909152805191012090565b5f8082516020840185f09050803b6105ac575f80fd5b6001600160a01b0381166105f257600154604080515f815260208101918290525f8051602061084a833981519152926105e9929091600391610822565b60405180910390a15b90505b92915050565b5f808383516020850187f59050803b610612575f80fd5b6001600160a01b03811661065857600154604080515f815260208101918290525f8051602061084a8339815191529261064f929091600391610822565b60405180910390a15b949350505050565b5f60208284031215610670575f80fd5b5035919050565b634e487b7160e01b5f52601260045260245ffd5b5f8261069957610699610677565b500690565b634e487b7160e01b5f52604160045260245ffd5b5f602082840312156106c2575f80fd5b815167ffffffffffffffff8111156106d8575f80fd5b8201601f810184136106e8575f80fd5b805167ffffffffffffffff8111156107025761070261069e565b604051601f8201601f19908116603f0116810167ffffffffffffffff811182821017156107315761073161069e565b604052818152828201602001861015610748575f80fd5b8160208401602083015e5f91810160200191909152949350505050565b634e487b7160e01b5f52601160045260245ffd5b808201808211156105f5576105f5610765565b5f81518060208401855e5f93019283525090919050565b5f6106586107b1838661078c565b8461078c565b5f826107c5576107c5610677565b500490565b80820281158282048414176105f5576105f5610765565b818103818111156105f5576105f5610765565b5f81518084528060208401602086015e5f602082860101526020601f19601f83011685010191505092915050565b838152826020820152606060408201525f61084060608301846107f4565b9594505050505056fe1cfe9a531f435de63e4684efae7f811234d1432eff9eb41d85b62ce30c477b8fa2646970667358221220df2332fb28ea6df56e3d90e10c50075d7e33f2484f9e4adfbe43ff1e2a1a338764736f6c634300081a00331cfe9a531f435de63e4684efae7f811234d1432eff9eb41d85b62ce30c477b8fa2646970667358221220b31ccedb0903b82c4cba4a8f25552fbcc395fa4a24b410728362f06b84760e2764736f6c634300081a0033",
}

// ContractABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractMetaData.ABI instead.
var ContractABI = ContractMetaData.ABI

// ContractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractMetaData.Bin instead.
var ContractBin = ContractMetaData.Bin

// DeployContract deploys a new Ethereum contract, binding an instance of Contract to it.
func DeployContract(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Contract, error) {
	parsed, err := ContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
}

// Contract is an auto generated Go binding around an Ethereum contract.
type Contract struct {
	ContractCaller     // Read-only binding to the contract
	ContractTransactor // Write-only binding to the contract
	ContractFilterer   // Log filterer for contract events
}

// ContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractSession struct {
	Contract     *Contract         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractCallerSession struct {
	Contract *ContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractTransactorSession struct {
	Contract     *ContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractRaw struct {
	Contract *Contract // Generic contract binding to access the raw methods on
}

// ContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractCallerRaw struct {
	Contract *ContractCaller // Generic read-only contract binding to access the raw methods on
}

// ContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractTransactorRaw struct {
	Contract *ContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContract creates a new instance of Contract, bound to a specific deployed contract.
func NewContract(address common.Address, backend bind.ContractBackend) (*Contract, error) {
	contract, err := bindContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
}

// NewContractCaller creates a new read-only instance of Contract, bound to a specific deployed contract.
func NewContractCaller(address common.Address, caller bind.ContractCaller) (*ContractCaller, error) {
	contract, err := bindContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractCaller{contract: contract}, nil
}

// NewContractTransactor creates a new write-only instance of Contract, bound to a specific deployed contract.
func NewContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractTransactor, error) {
	contract, err := bindContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractTransactor{contract: contract}, nil
}

// NewContractFilterer creates a new log filterer instance of Contract, bound to a specific deployed contract.
func NewContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractFilterer, error) {
	contract, err := bindContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractFilterer{contract: contract}, nil
}

// bindContract binds a generic wrapper to an already deployed contract.
func bindContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.ContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transact(opts, method, params...)
}

// ChildAddresses is a free data retrieval call binding the contract method 0x9c47ed9f.
//
// Solidity: function childAddresses(uint256 ) view returns(address)
func (_Contract *ContractCaller) ChildAddresses(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "childAddresses", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ChildAddresses is a free data retrieval call binding the contract method 0x9c47ed9f.
//
// Solidity: function childAddresses(uint256 ) view returns(address)
func (_Contract *ContractSession) ChildAddresses(arg0 *big.Int) (common.Address, error) {
	return _Contract.Contract.ChildAddresses(&_Contract.CallOpts, arg0)
}

// ChildAddresses is a free data retrieval call binding the contract method 0x9c47ed9f.
//
// Solidity: function childAddresses(uint256 ) view returns(address)
func (_Contract *ContractCallerSession) ChildAddresses(arg0 *big.Int) (common.Address, error) {
	return _Contract.Contract.ChildAddresses(&_Contract.CallOpts, arg0)
}

// ChildCode is a free data retrieval call binding the contract method 0x0e7bfac5.
//
// Solidity: function childCode() pure returns(bytes)
func (_Contract *ContractCaller) ChildCode(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "childCode")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// ChildCode is a free data retrieval call binding the contract method 0x0e7bfac5.
//
// Solidity: function childCode() pure returns(bytes)
func (_Contract *ContractSession) ChildCode() ([]byte, error) {
	return _Contract.Contract.ChildCode(&_Contract.CallOpts)
}

// ChildCode is a free data retrieval call binding the contract method 0x0e7bfac5.
//
// Solidity: function childCode() pure returns(bytes)
func (_Contract *ContractCallerSession) ChildCode() ([]byte, error) {
	return _Contract.Contract.ChildCode(&_Contract.CallOpts)
}

// ChildIdx is a free data retrieval call binding the contract method 0x00fd4ee9.
//
// Solidity: function childIdx() view returns(uint256)
func (_Contract *ContractCaller) ChildIdx(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "childIdx")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ChildIdx is a free data retrieval call binding the contract method 0x00fd4ee9.
//
// Solidity: function childIdx() view returns(uint256)
func (_Contract *ContractSession) ChildIdx() (*big.Int, error) {
	return _Contract.Contract.ChildIdx(&_Contract.CallOpts)
}

// ChildIdx is a free data retrieval call binding the contract method 0x00fd4ee9.
//
// Solidity: function childIdx() view returns(uint256)
func (_Contract *ContractCallerSession) ChildIdx() (*big.Int, error) {
	return _Contract.Contract.ChildIdx(&_Contract.CallOpts)
}

// Counter is a free data retrieval call binding the contract method 0x61bc221a.
//
// Solidity: function counter() view returns(uint256)
func (_Contract *ContractCaller) Counter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "counter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Counter is a free data retrieval call binding the contract method 0x61bc221a.
//
// Solidity: function counter() view returns(uint256)
func (_Contract *ContractSession) Counter() (*big.Int, error) {
	return _Contract.Contract.Counter(&_Contract.CallOpts)
}

// Counter is a free data retrieval call binding the contract method 0x61bc221a.
//
// Solidity: function counter() view returns(uint256)
func (_Contract *ContractCallerSession) Counter() (*big.Int, error) {
	return _Contract.Contract.Counter(&_Contract.CallOpts)
}

// Clean is a paid mutator transaction binding the contract method 0x674ef0fa.
//
// Solidity: function clean(uint256 count) returns()
func (_Contract *ContractTransactor) Clean(opts *bind.TransactOpts, count *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "clean", count)
}

// Clean is a paid mutator transaction binding the contract method 0x674ef0fa.
//
// Solidity: function clean(uint256 count) returns()
func (_Contract *ContractSession) Clean(count *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.Clean(&_Contract.TransactOpts, count)
}

// Clean is a paid mutator transaction binding the contract method 0x674ef0fa.
//
// Solidity: function clean(uint256 count) returns()
func (_Contract *ContractTransactorSession) Clean(count *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.Clean(&_Contract.TransactOpts, count)
}

// Notify is a paid mutator transaction binding the contract method 0xd52c3bd2.
//
// Solidity: function notify(address addr, bool destroy) returns()
func (_Contract *ContractTransactor) Notify(opts *bind.TransactOpts, addr common.Address, destroy bool) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "notify", addr, destroy)
}

// Notify is a paid mutator transaction binding the contract method 0xd52c3bd2.
//
// Solidity: function notify(address addr, bool destroy) returns()
func (_Contract *ContractSession) Notify(addr common.Address, destroy bool) (*types.Transaction, error) {
	return _Contract.Contract.Notify(&_Contract.TransactOpts, addr, destroy)
}

// Notify is a paid mutator transaction binding the contract method 0xd52c3bd2.
//
// Solidity: function notify(address addr, bool destroy) returns()
func (_Contract *ContractTransactorSession) Notify(addr common.Address, destroy bool) (*types.Transaction, error) {
	return _Contract.Contract.Notify(&_Contract.TransactOpts, addr, destroy)
}

// Test is a paid mutator transaction binding the contract method 0x29e99f07.
//
// Solidity: function test(uint256 seed) payable returns()
func (_Contract *ContractTransactor) Test(opts *bind.TransactOpts, seed *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "test", seed)
}

// Test is a paid mutator transaction binding the contract method 0x29e99f07.
//
// Solidity: function test(uint256 seed) payable returns()
func (_Contract *ContractSession) Test(seed *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.Test(&_Contract.TransactOpts, seed)
}

// Test is a paid mutator transaction binding the contract method 0x29e99f07.
//
// Solidity: function test(uint256 seed) payable returns()
func (_Contract *ContractTransactorSession) Test(seed *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.Test(&_Contract.TransactOpts, seed)
}

// ContractTestSeedIterator is returned from FilterTestSeed and is used to iterate over the raw logs and unpacked data for TestSeed events raised by the Contract contract.
type ContractTestSeedIterator struct {
	Event *ContractTestSeed // Event containing the contract specifics and raw log

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
func (it *ContractTestSeedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractTestSeed)
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
		it.Event = new(ContractTestSeed)
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
func (it *ContractTestSeedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractTestSeedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractTestSeed represents a TestSeed event raised by the Contract contract.
type ContractTestSeed struct {
	Seed *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterTestSeed is a free log retrieval operation binding the contract event 0x47ff30cfcf5248fde6dac21a1944ca5982927073500951ddde02f1d68d63c3da.
//
// Solidity: event TestSeed(uint256 seed)
func (_Contract *ContractFilterer) FilterTestSeed(opts *bind.FilterOpts) (*ContractTestSeedIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "TestSeed")
	if err != nil {
		return nil, err
	}
	return &ContractTestSeedIterator{contract: _Contract.contract, event: "TestSeed", logs: logs, sub: sub}, nil
}

// WatchTestSeed is a free log subscription operation binding the contract event 0x47ff30cfcf5248fde6dac21a1944ca5982927073500951ddde02f1d68d63c3da.
//
// Solidity: event TestSeed(uint256 seed)
func (_Contract *ContractFilterer) WatchTestSeed(opts *bind.WatchOpts, sink chan<- *ContractTestSeed) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "TestSeed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractTestSeed)
				if err := _Contract.contract.UnpackLog(event, "TestSeed", log); err != nil {
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

// ParseTestSeed is a log parse operation binding the contract event 0x47ff30cfcf5248fde6dac21a1944ca5982927073500951ddde02f1d68d63c3da.
//
// Solidity: event TestSeed(uint256 seed)
func (_Contract *ContractFilterer) ParseTestSeed(log types.Log) (*ContractTestSeed, error) {
	event := new(ContractTestSeed)
	if err := _Contract.contract.UnpackLog(event, "TestSeed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
