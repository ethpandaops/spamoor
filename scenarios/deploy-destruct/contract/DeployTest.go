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

// DeployTestMetaData contains all meta data concerning the DeployTest contract.
var DeployTestMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"}],\"name\":\"TestSeed\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"childAddresses\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"childCode\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"childIdx\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"count\",\"type\":\"uint256\"}],\"name\":\"clean\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"destroy\",\"type\":\"bool\"}],\"name\":\"notify\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"}],\"name\":\"test\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x608060405260015f5560018055348015610017575f80fd5b50611a18806100255f395ff3fe60806040526004361062000076575f3560e01c806361bc221a116200005257806361bc221a14620000e3578063674ef0fa14620000f95780639c47ed9f146200011d578063d52c3bd2146200016e575f80fd5b8062fd4ee9146200007a5780630e7bfac514620000a457806329e99f0714620000ca575b5f80fd5b34801562000086575f80fd5b506200009160015481565b6040519081526020015b60405180910390f35b348015620000b0575f80fd5b50620000bb62000192565b6040516200009b919062000753565b620000e1620000db36600462000787565b620001be565b005b348015620000ef575f80fd5b50620000915f5481565b34801562000105575f80fd5b50620000e16200011736600462000787565b62000482565b34801562000129575f80fd5b50620001556200013b36600462000787565b60026020525f90815260409020546001600160a01b031681565b6040516001600160a01b0390911681526020016200009b565b3480156200017a575f80fd5b50620000e16200018c3660046200079f565b62000540565b606060405180602001620001a69062000721565b601f1982820381018352601f90910116604052919050565b805f03620001e357620186a0620001d46200062e565b620001e09190620007fc565b90505b6040518181527f47ff30cfcf5248fde6dac21a1944ca5982927073500951ddde02f1d68d63c3da9060200160405180910390a15f62000225825f60016200067f565b9050601e62000236603283620007fc565b1015620002b6575f6200024862000192565b60408051602081018590525f9181019190915230606082015260800160408051601f198184030181529082905262000284929160200162000812565b60405160208183030381529060405290505f620002a24783620006b6565b9050620002b0815f62000540565b50505050565b6028620002c5603283620007fc565b1015620003aa575f5b6032811015620003a55760018054118015620002ec575062061a805a115b156200039c5760015f81905260026020527fe90b7bceb6e7df5418fb78d8ee546e97c83a08bbccc01a0644d599ccd2a7c2e0546001600160a01b0316908190639d118770906200033f908790866200067f565b6040518263ffffffff1660e01b81526004016200035e91815260200190565b5f604051808303815f87803b15801562000376575f80fd5b505af115801562000389573d5f803e3d5ffd5b505050506200039a81600162000540565b505b600101620002ce565b505050565b4760015b600154811015620002b0575f818152600260205260408120546001600160a01b031690620003de8682856200067f565b90505f620003ee605083620007fc565b620003fb60648762000858565b6200040791906200086e565b905062000415818662000888565b604051631584b49360e11b8152600481018490529095506001600160a01b03841690632b0969269083906024015f604051808303818588803b1580156200045a575f80fd5b505af11580156200046d573d5f803e3d5ffd5b505060019096019550620003ae945050505050565b80156200053d57600180541115620005285760015f5260026020527fe90b7bceb6e7df5418fb78d8ee546e97c83a08bbccc01a0644d599ccd2a7c2e0546040516309d1187760e41b8152600560048201526001600160a01b03909116908190639d118770906024015f604051808303815f87803b15801562000502575f80fd5b505af115801562000515573d5f803e3d5ffd5b505050506200052681600162000540565b505b8062000534816200089e565b91505062000482565b50565b8015620005d9576001600160a01b0382165f908152600360205260409020548015620003a5576001600160a01b0383165f908152600360205260408120819055600180549162000590836200089e565b90915550506001545f9081526002602090815260408083205484845281842080546001600160a01b0319166001600160a01b039092169182179055835260039091529020555050565b600180545f90815260026020908152604080832080546001600160a01b0319166001600160a01b03881690811790915584549084526003909252822081905591906200062583620008b6565b91905055505050565b5f805481806200063e83620008b6565b90915550505f546040805144602082015242918101919091526060810191909152608001604051602081830303815290604052805190602001205f1c905090565b6040805160208082019590955280820193909352606080840192909252805180840390920182526080909201909152805191012090565b5f8082516020840185f09050803b620006cd575f80fd5b6001600160a01b038116620007185760405162461bcd60e51b815260206004820152600d60248201526c18dc99585d194819985a5b1959609a1b604482015260640160405180910390fd5b90505b92915050565b61111180620008d283390190565b5f5b838110156200074b57818101518382015260200162000731565b50505f910152565b602081525f8251806020840152620007738160408501602087016200072f565b601f01601f19169190910160400192915050565b5f6020828403121562000798575f80fd5b5035919050565b5f8060408385031215620007b1575f80fd5b82356001600160a01b0381168114620007c8575f80fd5b915060208301358015158114620007dd575f80fd5b809150509250929050565b634e487b7160e01b5f52601260045260245ffd5b5f826200080d576200080d620007e8565b500690565b5f8351620008258184602088016200072f565b8351908301906200083b8183602088016200072f565b01949350505050565b634e487b7160e01b5f52601160045260245ffd5b5f82620008695762000869620007e8565b500490565b80820281158282048414176200071b576200071b62000844565b818103818111156200071b576200071b62000844565b5f81620008af57620008af62000844565b505f190190565b5f60018201620008ca57620008ca62000844565b506001019056fe608060405260405162001111380380620011118339810160408190526200002691620005dc565b60408051848152602081018490527f06ed9ff5e25ad09aa503577a10190fb63ca2ae4e7ea00bc39a59d3daa8bcbec3910160405180910390a15f80546001600160a01b0319166001600160a01b0383161790556001839055620000898262000092565b5050506200082d565b5f620000a8600154835f620004d160201b60201c565b905047600483108015620000c85750603c620000c660648462000634565b105b15620004cc575f805f9054906101000a90046001600160a01b03166001600160a01b0316630e7bfac56040518163ffffffff1660e01b81526004015f60405180830381865afa1580156200011e573d5f803e3d5ffd5b505050506040513d5f823e601f3d908101601f1916820160405262000147919081019062000682565b90505f5b6200015860038562000634565b6200016590600162000749565b811015620004c957621e84805a10620004c9575f6200018e6001548784620004d160201b60201c565b90505f8382620001a089600162000749565b5f546040805160208101949094528301919091526001600160a01b0316606082015260800160408051601f1981840301815290829052620001e592916020016200075f565b60405160208183030381529060405290506200024060015460016200020b919062000749565b604080516020808201939093528082018b90526060808201889052825180830390910181526080909101909152805191012090565b91505f6200025060508462000634565b6200025d60648862000791565b620002699190620007a7565b9050620002778187620007c1565b95505f60326200028960648662000634565b1015620002a4576200029c828462000508565b9050620002c0565b620002bd82620002b687600162000749565b8562000572565b90505b6001600160a01b038116620002d95750505050620004c9565b620003236001546002620002ee919062000749565b604080516020808201939093528082018d905260608082018a9052825180830390910181526080909101909152805191012090565b935060286200033460648662000634565b1015620003fe576040516309d1187760e41b8152600481018590526001600160a01b03821690639d118770906024015f604051808303815f87803b1580156200037b575f80fd5b505af19250505080156200038d575060015b620003f8573d808015620003bd576040519150601f19603f3d011682016040523d82523d5f602084013e620003c2565b606091505b505f80516020620010f1833981519152600154600183604051620003e99392919062000804565b60405180910390a150620004b8565b620004b8565b5f8054604051636a961de960e11b81526001600160a01b038481166004830152602482019390935291169063d52c3bd2906044015f604051808303815f87803b1580156200044a575f80fd5b505af19250505080156200045c575060015b620004b8573d8080156200048c576040519150601f19603f3d011682016040523d82523d5f602084013e62000491565b606091505b505f80516020620010f1833981519152600154600283604051620003e99392919062000804565b5050600190920191506200014b9050565b50505b505050565b6040805160208082019590955280820193909352606080840192909252805180840390920182526080909201909152805191012090565b5f8082516020840185f09050803b6200051f575f80fd5b6001600160a01b0381166200056957600154604080515f815260208101918290525f80516020620010f1833981519152926200056092909160039162000804565b60405180910390a15b90505b92915050565b5f808383516020850187f59050803b6200058a575f80fd5b6001600160a01b038116620005d457600154604080515f815260208101918290525f80516020620010f183398151915292620005cb92909160039162000804565b60405180910390a15b949350505050565b5f805f60608486031215620005ef575f80fd5b83516020850151604086015191945092506001600160a01b038116811462000615575f80fd5b809150509250925092565b634e487b7160e01b5f52601260045260245ffd5b5f8262000645576200064562000620565b500690565b634e487b7160e01b5f52604160045260245ffd5b5f5b838110156200067a57818101518382015260200162000660565b50505f910152565b5f6020828403121562000693575f80fd5b81516001600160401b0380821115620006aa575f80fd5b818401915084601f830112620006be575f80fd5b815181811115620006d357620006d36200064a565b604051601f8201601f19908116603f01168101908382118183101715620006fe57620006fe6200064a565b8160405282815287602084870101111562000717575f80fd5b6200072a8360208301602088016200065e565b979650505050505050565b634e487b7160e01b5f52601160045260245ffd5b808201808211156200056c576200056c62000735565b5f8351620007728184602088016200065e565b835190830190620007888183602088016200065e565b01949350505050565b5f82620007a257620007a262000620565b500490565b80820281158282048414176200056c576200056c62000735565b818103818111156200056c576200056c62000735565b5f8151808452620007f08160208601602086016200065e565b601f01601f19169290920160200192915050565b838152826020820152606060408201525f620008246060830184620007d7565b95945050505050565b6108b6806200083b5f395ff3fe60806040526004361061003e575f3560e01c80632b096926146100425780635a34d356146100575780637469d0681461007f5780639d118770146100b5575b5f80fd5b610055610050366004610660565b6100d4565b005b348015610062575f80fd5b5061006c60015481565b6040519081526020015b60405180910390f35b34801561008a575f80fd5b505f5461009d906001600160a01b031681565b6040516001600160a01b039091168152602001610076565b3480156100c0575f80fd5b506100556100cf366004610660565b610466565b5f6100e2600154835f61055f565b9050476004831080156100ff5750603c6100fd60648461068b565b105b15610461575f805f9054906101000a90046001600160a01b03166001600160a01b0316630e7bfac56040518163ffffffff1660e01b81526004015f60405180830381865afa158015610153573d5f803e3d5ffd5b505050506040513d5f823e601f3d908101601f1916820160405261017a91908101906106d4565b90505f5b61018960038561068b565b610194906001610790565b81101561045e57621e84805a1061045e575f6101b3600154878461055f565b90505f83826101c3896001610790565b5f546040805160208101949094528301919091526001600160a01b0316606082015260800160408051601f198184030181529082905261020692916020016107a3565b604051602081830303815290604052905061023060015460016102299190610790565b888561055f565b91505f61023e60508461068b565b6102496064886107d1565b61025391906107e4565b905061025f81876107fb565b95505f603261026f60648661068b565b10156102865761027f8284610596565b905061029e565b61029b82610295876001610790565b856105fb565b90505b6001600160a01b0381166102b5575050505061045e565b6102ce60015460026102c79190610790565b8a8761055f565b935060286102dd60648661068b565b101561039c576040516309d1187760e41b8152600481018590526001600160a01b03821690639d118770906024015f604051808303815f87803b158015610322575f80fd5b505af1925050508015610333575060015b610397573d808015610360576040519150601f19603f3d011682016040523d82523d5f602084013e610365565b606091505b505f8051602061086183398151915260015460018360405161038993929190610839565b60405180910390a15061044e565b61044e565b5f8054604051636a961de960e11b81526001600160a01b038481166004830152602482019390935291169063d52c3bd2906044015f604051808303815f87803b1580156103e7575f80fd5b505af19250505080156103f8575060015b61044e573d808015610425576040519150601f19603f3d011682016040523d82523d5f602084013e61042a565b606091505b505f8051602061086183398151915260015460028360405161038993929190610839565b50506001909201915061017e9050565b50505b505050565b5f8061047360068461068b565b9050805f0361048457339150610512565b806001036104b8576040805160208101859052016040516020818303038152906040528051906020012060601c9150610512565b806002036104c857329150610512565b806003036104d857309150610512565b806004036104fc577349e0fd3800c117357057534e30c5b5115c6734889150610512565b80600503610512575f546001600160a01b031691505b604080518481526001600160a01b03841660208201527f77b29af0d4f525b395d176ebc2772a5ffe882cd17c4b934ee2fd3b773ba41040910160405180910390a1816001600160a01b0316ff5b6040805160208082019590955280820193909352606080840192909252805180840390920182526080909201909152805191012090565b5f8082516020840185f09050803b6105ac575f80fd5b6001600160a01b0381166105f257600154604080515f815260208101918290525f80516020610861833981519152926105e9929091600391610839565b60405180910390a15b90505b92915050565b5f808383516020850187f59050803b610612575f80fd5b6001600160a01b03811661065857600154604080515f815260208101918290525f805160206108618339815191529261064f929091600391610839565b60405180910390a15b949350505050565b5f60208284031215610670575f80fd5b5035919050565b634e487b7160e01b5f52601260045260245ffd5b5f8261069957610699610677565b500690565b634e487b7160e01b5f52604160045260245ffd5b5f5b838110156106cc5781810151838201526020016106b4565b50505f910152565b5f602082840312156106e4575f80fd5b815167ffffffffffffffff808211156106fb575f80fd5b818401915084601f83011261070e575f80fd5b8151818111156107205761072061069e565b604051601f8201601f19908116603f011681019083821181831017156107485761074861069e565b81604052828152876020848701011115610760575f80fd5b6107718360208301602088016106b2565b979650505050505050565b634e487b7160e01b5f52601160045260245ffd5b808201808211156105f5576105f561077c565b5f83516107b48184602088016106b2565b8351908301906107c88183602088016106b2565b01949350505050565b5f826107df576107df610677565b500490565b80820281158282048414176105f5576105f561077c565b818103818111156105f5576105f561077c565b5f81518084526108258160208601602086016106b2565b601f01601f19169290920160200192915050565b838152826020820152606060408201525f610857606083018461080e565b9594505050505056fe1cfe9a531f435de63e4684efae7f811234d1432eff9eb41d85b62ce30c477b8fa26469706673582212204b8eddc23a720cba6c8a95b8129042142410d0fa848777d86d6f0552cc468bb864736f6c634300081600331cfe9a531f435de63e4684efae7f811234d1432eff9eb41d85b62ce30c477b8fa2646970667358221220d4cad9e5595c01d4e9195704f02c698ec29b8fa0ebee2e57951d7736cb4d45f964736f6c63430008160033",
}

// DeployTestABI is the input ABI used to generate the binding from.
// Deprecated: Use DeployTestMetaData.ABI instead.
var DeployTestABI = DeployTestMetaData.ABI

// DeployTestBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DeployTestMetaData.Bin instead.
var DeployTestBin = DeployTestMetaData.Bin

// DeployDeployTest deploys a new Ethereum contract, binding an instance of DeployTest to it.
func DeployDeployTest(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *DeployTest, error) {
	parsed, err := DeployTestMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DeployTestBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DeployTest{DeployTestCaller: DeployTestCaller{contract: contract}, DeployTestTransactor: DeployTestTransactor{contract: contract}, DeployTestFilterer: DeployTestFilterer{contract: contract}}, nil
}

// DeployTest is an auto generated Go binding around an Ethereum contract.
type DeployTest struct {
	DeployTestCaller     // Read-only binding to the contract
	DeployTestTransactor // Write-only binding to the contract
	DeployTestFilterer   // Log filterer for contract events
}

// DeployTestCaller is an auto generated read-only Go binding around an Ethereum contract.
type DeployTestCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DeployTestTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DeployTestTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DeployTestFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DeployTestFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DeployTestSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DeployTestSession struct {
	Contract     *DeployTest       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DeployTestCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DeployTestCallerSession struct {
	Contract *DeployTestCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// DeployTestTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DeployTestTransactorSession struct {
	Contract     *DeployTestTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// DeployTestRaw is an auto generated low-level Go binding around an Ethereum contract.
type DeployTestRaw struct {
	Contract *DeployTest // Generic contract binding to access the raw methods on
}

// DeployTestCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DeployTestCallerRaw struct {
	Contract *DeployTestCaller // Generic read-only contract binding to access the raw methods on
}

// DeployTestTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DeployTestTransactorRaw struct {
	Contract *DeployTestTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDeployTest creates a new instance of DeployTest, bound to a specific deployed contract.
func NewDeployTest(address common.Address, backend bind.ContractBackend) (*DeployTest, error) {
	contract, err := bindDeployTest(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DeployTest{DeployTestCaller: DeployTestCaller{contract: contract}, DeployTestTransactor: DeployTestTransactor{contract: contract}, DeployTestFilterer: DeployTestFilterer{contract: contract}}, nil
}

// NewDeployTestCaller creates a new read-only instance of DeployTest, bound to a specific deployed contract.
func NewDeployTestCaller(address common.Address, caller bind.ContractCaller) (*DeployTestCaller, error) {
	contract, err := bindDeployTest(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DeployTestCaller{contract: contract}, nil
}

// NewDeployTestTransactor creates a new write-only instance of DeployTest, bound to a specific deployed contract.
func NewDeployTestTransactor(address common.Address, transactor bind.ContractTransactor) (*DeployTestTransactor, error) {
	contract, err := bindDeployTest(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DeployTestTransactor{contract: contract}, nil
}

// NewDeployTestFilterer creates a new log filterer instance of DeployTest, bound to a specific deployed contract.
func NewDeployTestFilterer(address common.Address, filterer bind.ContractFilterer) (*DeployTestFilterer, error) {
	contract, err := bindDeployTest(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DeployTestFilterer{contract: contract}, nil
}

// bindDeployTest binds a generic wrapper to an already deployed contract.
func bindDeployTest(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DeployTestMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DeployTest *DeployTestRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DeployTest.Contract.DeployTestCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DeployTest *DeployTestRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DeployTest.Contract.DeployTestTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DeployTest *DeployTestRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DeployTest.Contract.DeployTestTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DeployTest *DeployTestCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DeployTest.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DeployTest *DeployTestTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DeployTest.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DeployTest *DeployTestTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DeployTest.Contract.contract.Transact(opts, method, params...)
}

// ChildAddresses is a free data retrieval call binding the contract method 0x9c47ed9f.
//
// Solidity: function childAddresses(uint256 ) view returns(address)
func (_DeployTest *DeployTestCaller) ChildAddresses(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _DeployTest.contract.Call(opts, &out, "childAddresses", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ChildAddresses is a free data retrieval call binding the contract method 0x9c47ed9f.
//
// Solidity: function childAddresses(uint256 ) view returns(address)
func (_DeployTest *DeployTestSession) ChildAddresses(arg0 *big.Int) (common.Address, error) {
	return _DeployTest.Contract.ChildAddresses(&_DeployTest.CallOpts, arg0)
}

// ChildAddresses is a free data retrieval call binding the contract method 0x9c47ed9f.
//
// Solidity: function childAddresses(uint256 ) view returns(address)
func (_DeployTest *DeployTestCallerSession) ChildAddresses(arg0 *big.Int) (common.Address, error) {
	return _DeployTest.Contract.ChildAddresses(&_DeployTest.CallOpts, arg0)
}

// ChildCode is a free data retrieval call binding the contract method 0x0e7bfac5.
//
// Solidity: function childCode() pure returns(bytes)
func (_DeployTest *DeployTestCaller) ChildCode(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _DeployTest.contract.Call(opts, &out, "childCode")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// ChildCode is a free data retrieval call binding the contract method 0x0e7bfac5.
//
// Solidity: function childCode() pure returns(bytes)
func (_DeployTest *DeployTestSession) ChildCode() ([]byte, error) {
	return _DeployTest.Contract.ChildCode(&_DeployTest.CallOpts)
}

// ChildCode is a free data retrieval call binding the contract method 0x0e7bfac5.
//
// Solidity: function childCode() pure returns(bytes)
func (_DeployTest *DeployTestCallerSession) ChildCode() ([]byte, error) {
	return _DeployTest.Contract.ChildCode(&_DeployTest.CallOpts)
}

// ChildIdx is a free data retrieval call binding the contract method 0x00fd4ee9.
//
// Solidity: function childIdx() view returns(uint256)
func (_DeployTest *DeployTestCaller) ChildIdx(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DeployTest.contract.Call(opts, &out, "childIdx")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ChildIdx is a free data retrieval call binding the contract method 0x00fd4ee9.
//
// Solidity: function childIdx() view returns(uint256)
func (_DeployTest *DeployTestSession) ChildIdx() (*big.Int, error) {
	return _DeployTest.Contract.ChildIdx(&_DeployTest.CallOpts)
}

// ChildIdx is a free data retrieval call binding the contract method 0x00fd4ee9.
//
// Solidity: function childIdx() view returns(uint256)
func (_DeployTest *DeployTestCallerSession) ChildIdx() (*big.Int, error) {
	return _DeployTest.Contract.ChildIdx(&_DeployTest.CallOpts)
}

// Counter is a free data retrieval call binding the contract method 0x61bc221a.
//
// Solidity: function counter() view returns(uint256)
func (_DeployTest *DeployTestCaller) Counter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DeployTest.contract.Call(opts, &out, "counter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Counter is a free data retrieval call binding the contract method 0x61bc221a.
//
// Solidity: function counter() view returns(uint256)
func (_DeployTest *DeployTestSession) Counter() (*big.Int, error) {
	return _DeployTest.Contract.Counter(&_DeployTest.CallOpts)
}

// Counter is a free data retrieval call binding the contract method 0x61bc221a.
//
// Solidity: function counter() view returns(uint256)
func (_DeployTest *DeployTestCallerSession) Counter() (*big.Int, error) {
	return _DeployTest.Contract.Counter(&_DeployTest.CallOpts)
}

// Clean is a paid mutator transaction binding the contract method 0x674ef0fa.
//
// Solidity: function clean(uint256 count) returns()
func (_DeployTest *DeployTestTransactor) Clean(opts *bind.TransactOpts, count *big.Int) (*types.Transaction, error) {
	return _DeployTest.contract.Transact(opts, "clean", count)
}

// Clean is a paid mutator transaction binding the contract method 0x674ef0fa.
//
// Solidity: function clean(uint256 count) returns()
func (_DeployTest *DeployTestSession) Clean(count *big.Int) (*types.Transaction, error) {
	return _DeployTest.Contract.Clean(&_DeployTest.TransactOpts, count)
}

// Clean is a paid mutator transaction binding the contract method 0x674ef0fa.
//
// Solidity: function clean(uint256 count) returns()
func (_DeployTest *DeployTestTransactorSession) Clean(count *big.Int) (*types.Transaction, error) {
	return _DeployTest.Contract.Clean(&_DeployTest.TransactOpts, count)
}

// Notify is a paid mutator transaction binding the contract method 0xd52c3bd2.
//
// Solidity: function notify(address addr, bool destroy) returns()
func (_DeployTest *DeployTestTransactor) Notify(opts *bind.TransactOpts, addr common.Address, destroy bool) (*types.Transaction, error) {
	return _DeployTest.contract.Transact(opts, "notify", addr, destroy)
}

// Notify is a paid mutator transaction binding the contract method 0xd52c3bd2.
//
// Solidity: function notify(address addr, bool destroy) returns()
func (_DeployTest *DeployTestSession) Notify(addr common.Address, destroy bool) (*types.Transaction, error) {
	return _DeployTest.Contract.Notify(&_DeployTest.TransactOpts, addr, destroy)
}

// Notify is a paid mutator transaction binding the contract method 0xd52c3bd2.
//
// Solidity: function notify(address addr, bool destroy) returns()
func (_DeployTest *DeployTestTransactorSession) Notify(addr common.Address, destroy bool) (*types.Transaction, error) {
	return _DeployTest.Contract.Notify(&_DeployTest.TransactOpts, addr, destroy)
}

// Test is a paid mutator transaction binding the contract method 0x29e99f07.
//
// Solidity: function test(uint256 seed) payable returns()
func (_DeployTest *DeployTestTransactor) Test(opts *bind.TransactOpts, seed *big.Int) (*types.Transaction, error) {
	return _DeployTest.contract.Transact(opts, "test", seed)
}

// Test is a paid mutator transaction binding the contract method 0x29e99f07.
//
// Solidity: function test(uint256 seed) payable returns()
func (_DeployTest *DeployTestSession) Test(seed *big.Int) (*types.Transaction, error) {
	return _DeployTest.Contract.Test(&_DeployTest.TransactOpts, seed)
}

// Test is a paid mutator transaction binding the contract method 0x29e99f07.
//
// Solidity: function test(uint256 seed) payable returns()
func (_DeployTest *DeployTestTransactorSession) Test(seed *big.Int) (*types.Transaction, error) {
	return _DeployTest.Contract.Test(&_DeployTest.TransactOpts, seed)
}

// DeployTestTestSeedIterator is returned from FilterTestSeed and is used to iterate over the raw logs and unpacked data for TestSeed events raised by the DeployTest contract.
type DeployTestTestSeedIterator struct {
	Event *DeployTestTestSeed // Event containing the contract specifics and raw log

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
func (it *DeployTestTestSeedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployTestTestSeed)
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
		it.Event = new(DeployTestTestSeed)
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
func (it *DeployTestTestSeedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployTestTestSeedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployTestTestSeed represents a TestSeed event raised by the DeployTest contract.
type DeployTestTestSeed struct {
	Seed *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterTestSeed is a free log retrieval operation binding the contract event 0x47ff30cfcf5248fde6dac21a1944ca5982927073500951ddde02f1d68d63c3da.
//
// Solidity: event TestSeed(uint256 seed)
func (_DeployTest *DeployTestFilterer) FilterTestSeed(opts *bind.FilterOpts) (*DeployTestTestSeedIterator, error) {

	logs, sub, err := _DeployTest.contract.FilterLogs(opts, "TestSeed")
	if err != nil {
		return nil, err
	}
	return &DeployTestTestSeedIterator{contract: _DeployTest.contract, event: "TestSeed", logs: logs, sub: sub}, nil
}

// WatchTestSeed is a free log subscription operation binding the contract event 0x47ff30cfcf5248fde6dac21a1944ca5982927073500951ddde02f1d68d63c3da.
//
// Solidity: event TestSeed(uint256 seed)
func (_DeployTest *DeployTestFilterer) WatchTestSeed(opts *bind.WatchOpts, sink chan<- *DeployTestTestSeed) (event.Subscription, error) {

	logs, sub, err := _DeployTest.contract.WatchLogs(opts, "TestSeed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployTestTestSeed)
				if err := _DeployTest.contract.UnpackLog(event, "TestSeed", log); err != nil {
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
func (_DeployTest *DeployTestFilterer) ParseTestSeed(log types.Log) (*DeployTestTestSeed, error) {
	event := new(DeployTestTestSeed)
	if err := _DeployTest.contract.UnpackLog(event, "TestSeed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
