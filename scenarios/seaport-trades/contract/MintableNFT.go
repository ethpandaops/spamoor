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

// MintableNFTMetaData contains all meta data concerning the MintableNFT contract.
var MintableNFTMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_symbol\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"startId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"count\",\"type\":\"uint256\"}],\"name\":\"mintBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801562000010575f80fd5b5060405162000d9938038062000d99833981016040819052620000339162000119565b5f62000040838262000209565b5060016200004f828262000209565b505050620002d5565b634e487b7160e01b5f52604160045260245ffd5b5f82601f8301126200007c575f80fd5b81516001600160401b038082111562000099576200009962000058565b604051601f8301601f19908116603f01168101908282118183101715620000c457620000c462000058565b8160405283815260209250866020858801011115620000e1575f80fd5b5f91505b83821015620001045785820183015181830184015290820190620000e5565b5f602085830101528094505050505092915050565b5f80604083850312156200012b575f80fd5b82516001600160401b038082111562000142575f80fd5b62000150868387016200006c565b9350602085015191508082111562000166575f80fd5b5062000175858286016200006c565b9150509250929050565b600181811c908216806200019457607f821691505b602082108103620001b357634e487b7160e01b5f52602260045260245ffd5b50919050565b601f8211156200020457805f5260205f20601f840160051c81016020851015620001e05750805b601f840160051c820191505b8181101562000201575f8155600101620001ec565b50505b505050565b81516001600160401b0381111562000225576200022562000058565b6200023d816200023684546200017f565b84620001b9565b602080601f83116001811462000273575f84156200025b5750858301515b5f19600386901b1c1916600185901b178555620002cd565b5f85815260208120601f198616915b82811015620002a35788860151825594840194600190910190840162000282565b5085821015620002c157878501515f19600388901b60f8161c191681555b505060018460011b0185555b505050505050565b610ab680620002e35f395ff3fe608060405234801561000f575f80fd5b50600436106100b1575f3560e01c806340c10f191161006e57806340c10f19146101585780636352211e1461016b57806370a082311461017e57806395d89b411461019f578063a22cb465146101a7578063e985e9c5146101ba575f80fd5b806301ffc9a7146100b557806306fdde03146100dd578063081812fc146100f2578063095ea7b31461011d57806323b872dd146101325780632e81aaea14610145575b5f80fd5b6100c86100c3366004610828565b6101f5565b60405190151581526020015b60405180910390f35b6100e561022b565b6040516100d49190610856565b6101056101003660046108a2565b6102b6565b6040516001600160a01b0390911681526020016100d4565b61013061012b3660046108d4565b61032d565b005b6101306101403660046108fc565b61040a565b610130610153366004610935565b6105f0565b6101306101663660046108d4565b610616565b6101056101793660046108a2565b610716565b61019161018c366004610965565b61076e565b6040519081526020016100d4565b6100e56107b0565b6101306101b536600461097e565b6107bd565b6100c86101c83660046109b7565b6001600160a01b039182165f90815260056020908152604080832093909416825291909152205460ff1690565b5f6301ffc9a760e01b6001600160e01b03198316148061022557506380ac58cd60e01b6001600160e01b03198316145b92915050565b5f8054610237906109e8565b80601f0160208091040260200160405190810160405280929190818152602001828054610263906109e8565b80156102ae5780601f10610285576101008083540402835291602001916102ae565b820191905f5260205f20905b81548152906001019060200180831161029157829003601f168201915b505050505081565b5f818152600260205260408120546001600160a01b03166103125760405162461bcd60e51b81526020600482015260116024820152703737b732bc34b9ba32b73a103a37b5b2b760791b60448201526064015b60405180910390fd5b505f908152600460205260409020546001600160a01b031690565b5f61033782610716565b9050336001600160a01b038216148061037257506001600160a01b0381165f90815260056020908152604080832033845290915290205460ff165b6103af5760405162461bcd60e51b815260206004820152600e60248201526d1b9bdd08185d5d1a1bdc9a5e995960921b6044820152606401610309565b5f8281526004602052604080822080546001600160a01b0319166001600160a01b0387811691821790925591518593918516917f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92591a4505050565b5f61041482610716565b9050836001600160a01b0316816001600160a01b0316146104645760405162461bcd60e51b815260206004820152600a60248201526977726f6e672066726f6d60b01b6044820152606401610309565b6001600160a01b03831661048a5760405162461bcd60e51b815260040161030990610a20565b336001600160a01b03821614806104b657505f828152600460205260409020546001600160a01b031633145b806104e357506001600160a01b0381165f90815260056020908152604080832033845290915290205460ff165b6105205760405162461bcd60e51b815260206004820152600e60248201526d1b9bdd08185d5d1a1bdc9a5e995960921b6044820152606401610309565b5f82815260046020908152604080832080546001600160a01b03191690556001600160a01b038716835260039091528120805460019290610562908490610a5a565b90915550506001600160a01b0383165f90815260036020526040812080546001929061058f908490610a6d565b90915550505f8281526002602052604080822080546001600160a01b0319166001600160a01b0387811691821790925591518593918816917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef91a450505050565b5f5b8181101561061057610608846101668386610a6d565b6001016105f2565b50505050565b6001600160a01b03821661063c5760405162461bcd60e51b815260040161030990610a20565b5f818152600260205260409020546001600160a01b0316156106915760405162461bcd60e51b815260206004820152600e60248201526d185b1c9958591e481b5a5b9d195960921b6044820152606401610309565b6001600160a01b0382165f9081526003602052604081208054600192906106b9908490610a6d565b90915550505f8181526002602052604080822080546001600160a01b0319166001600160a01b03861690811790915590518392907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef908290a45050565b5f818152600260205260408120546001600160a01b0316806102255760405162461bcd60e51b81526020600482015260116024820152703737b732bc34b9ba32b73a103a37b5b2b760791b6044820152606401610309565b5f6001600160a01b0382166107955760405162461bcd60e51b815260040161030990610a20565b506001600160a01b03165f9081526003602052604090205490565b60018054610237906109e8565b335f8181526005602090815260408083206001600160a01b03871680855290835292819020805460ff191686151590811790915590519081529192917f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31910160405180910390a35050565b5f60208284031215610838575f80fd5b81356001600160e01b03198116811461084f575f80fd5b9392505050565b5f602080835283518060208501525f5b8181101561088257858101830151858201604001528201610866565b505f604082860101526040601f19601f8301168501019250505092915050565b5f602082840312156108b2575f80fd5b5035919050565b80356001600160a01b03811681146108cf575f80fd5b919050565b5f80604083850312156108e5575f80fd5b6108ee836108b9565b946020939093013593505050565b5f805f6060848603121561090e575f80fd5b610917846108b9565b9250610925602085016108b9565b9150604084013590509250925092565b5f805f60608486031215610947575f80fd5b610950846108b9565b95602085013595506040909401359392505050565b5f60208284031215610975575f80fd5b61084f826108b9565b5f806040838503121561098f575f80fd5b610998836108b9565b9150602083013580151581146109ac575f80fd5b809150509250929050565b5f80604083850312156109c8575f80fd5b6109d1836108b9565b91506109df602084016108b9565b90509250929050565b600181811c908216806109fc57607f821691505b602082108103610a1a57634e487b7160e01b5f52602260045260245ffd5b50919050565b6020808252600c908201526b7a65726f206164647265737360a01b604082015260600190565b634e487b7160e01b5f52601160045260245ffd5b8181038181111561022557610225610a46565b8082018082111561022557610225610a4656fea2646970667358221220792dbc282367976908aa3b6a88b402e4b92e40b586fec582338033b57e7b281e64736f6c63430008180033",
}

// MintableNFTABI is the input ABI used to generate the binding from.
// Deprecated: Use MintableNFTMetaData.ABI instead.
var MintableNFTABI = MintableNFTMetaData.ABI

// MintableNFTBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MintableNFTMetaData.Bin instead.
var MintableNFTBin = MintableNFTMetaData.Bin

// DeployMintableNFT deploys a new Ethereum contract, binding an instance of MintableNFT to it.
func DeployMintableNFT(auth *bind.TransactOpts, backend bind.ContractBackend, _name string, _symbol string) (common.Address, *types.Transaction, *MintableNFT, error) {
	parsed, err := MintableNFTMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MintableNFTBin), backend, _name, _symbol)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MintableNFT{MintableNFTCaller: MintableNFTCaller{contract: contract}, MintableNFTTransactor: MintableNFTTransactor{contract: contract}, MintableNFTFilterer: MintableNFTFilterer{contract: contract}}, nil
}

// MintableNFT is an auto generated Go binding around an Ethereum contract.
type MintableNFT struct {
	MintableNFTCaller     // Read-only binding to the contract
	MintableNFTTransactor // Write-only binding to the contract
	MintableNFTFilterer   // Log filterer for contract events
}

// MintableNFTCaller is an auto generated read-only Go binding around an Ethereum contract.
type MintableNFTCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MintableNFTTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MintableNFTTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MintableNFTFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MintableNFTFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MintableNFTSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MintableNFTSession struct {
	Contract     *MintableNFT      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MintableNFTCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MintableNFTCallerSession struct {
	Contract *MintableNFTCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// MintableNFTTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MintableNFTTransactorSession struct {
	Contract     *MintableNFTTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// MintableNFTRaw is an auto generated low-level Go binding around an Ethereum contract.
type MintableNFTRaw struct {
	Contract *MintableNFT // Generic contract binding to access the raw methods on
}

// MintableNFTCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MintableNFTCallerRaw struct {
	Contract *MintableNFTCaller // Generic read-only contract binding to access the raw methods on
}

// MintableNFTTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MintableNFTTransactorRaw struct {
	Contract *MintableNFTTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMintableNFT creates a new instance of MintableNFT, bound to a specific deployed contract.
func NewMintableNFT(address common.Address, backend bind.ContractBackend) (*MintableNFT, error) {
	contract, err := bindMintableNFT(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MintableNFT{MintableNFTCaller: MintableNFTCaller{contract: contract}, MintableNFTTransactor: MintableNFTTransactor{contract: contract}, MintableNFTFilterer: MintableNFTFilterer{contract: contract}}, nil
}

// NewMintableNFTCaller creates a new read-only instance of MintableNFT, bound to a specific deployed contract.
func NewMintableNFTCaller(address common.Address, caller bind.ContractCaller) (*MintableNFTCaller, error) {
	contract, err := bindMintableNFT(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MintableNFTCaller{contract: contract}, nil
}

// NewMintableNFTTransactor creates a new write-only instance of MintableNFT, bound to a specific deployed contract.
func NewMintableNFTTransactor(address common.Address, transactor bind.ContractTransactor) (*MintableNFTTransactor, error) {
	contract, err := bindMintableNFT(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MintableNFTTransactor{contract: contract}, nil
}

// NewMintableNFTFilterer creates a new log filterer instance of MintableNFT, bound to a specific deployed contract.
func NewMintableNFTFilterer(address common.Address, filterer bind.ContractFilterer) (*MintableNFTFilterer, error) {
	contract, err := bindMintableNFT(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MintableNFTFilterer{contract: contract}, nil
}

// bindMintableNFT binds a generic wrapper to an already deployed contract.
func bindMintableNFT(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MintableNFTMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MintableNFT *MintableNFTRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MintableNFT.Contract.MintableNFTCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MintableNFT *MintableNFTRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MintableNFT.Contract.MintableNFTTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MintableNFT *MintableNFTRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MintableNFT.Contract.MintableNFTTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MintableNFT *MintableNFTCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MintableNFT.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MintableNFT *MintableNFTTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MintableNFT.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MintableNFT *MintableNFTTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MintableNFT.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_MintableNFT *MintableNFTCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _MintableNFT.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_MintableNFT *MintableNFTSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _MintableNFT.Contract.BalanceOf(&_MintableNFT.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_MintableNFT *MintableNFTCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _MintableNFT.Contract.BalanceOf(&_MintableNFT.CallOpts, owner)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_MintableNFT *MintableNFTCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _MintableNFT.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_MintableNFT *MintableNFTSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _MintableNFT.Contract.GetApproved(&_MintableNFT.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_MintableNFT *MintableNFTCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _MintableNFT.Contract.GetApproved(&_MintableNFT.CallOpts, tokenId)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_MintableNFT *MintableNFTCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _MintableNFT.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_MintableNFT *MintableNFTSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _MintableNFT.Contract.IsApprovedForAll(&_MintableNFT.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_MintableNFT *MintableNFTCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _MintableNFT.Contract.IsApprovedForAll(&_MintableNFT.CallOpts, owner, operator)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_MintableNFT *MintableNFTCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MintableNFT.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_MintableNFT *MintableNFTSession) Name() (string, error) {
	return _MintableNFT.Contract.Name(&_MintableNFT.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_MintableNFT *MintableNFTCallerSession) Name() (string, error) {
	return _MintableNFT.Contract.Name(&_MintableNFT.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_MintableNFT *MintableNFTCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _MintableNFT.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_MintableNFT *MintableNFTSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _MintableNFT.Contract.OwnerOf(&_MintableNFT.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_MintableNFT *MintableNFTCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _MintableNFT.Contract.OwnerOf(&_MintableNFT.CallOpts, tokenId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) pure returns(bool)
func (_MintableNFT *MintableNFTCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _MintableNFT.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) pure returns(bool)
func (_MintableNFT *MintableNFTSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _MintableNFT.Contract.SupportsInterface(&_MintableNFT.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) pure returns(bool)
func (_MintableNFT *MintableNFTCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _MintableNFT.Contract.SupportsInterface(&_MintableNFT.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_MintableNFT *MintableNFTCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MintableNFT.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_MintableNFT *MintableNFTSession) Symbol() (string, error) {
	return _MintableNFT.Contract.Symbol(&_MintableNFT.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_MintableNFT *MintableNFTCallerSession) Symbol() (string, error) {
	return _MintableNFT.Contract.Symbol(&_MintableNFT.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_MintableNFT *MintableNFTTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _MintableNFT.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_MintableNFT *MintableNFTSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _MintableNFT.Contract.Approve(&_MintableNFT.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_MintableNFT *MintableNFTTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _MintableNFT.Contract.Approve(&_MintableNFT.TransactOpts, to, tokenId)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 tokenId) returns()
func (_MintableNFT *MintableNFTTransactor) Mint(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _MintableNFT.contract.Transact(opts, "mint", to, tokenId)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 tokenId) returns()
func (_MintableNFT *MintableNFTSession) Mint(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _MintableNFT.Contract.Mint(&_MintableNFT.TransactOpts, to, tokenId)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 tokenId) returns()
func (_MintableNFT *MintableNFTTransactorSession) Mint(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _MintableNFT.Contract.Mint(&_MintableNFT.TransactOpts, to, tokenId)
}

// MintBatch is a paid mutator transaction binding the contract method 0x2e81aaea.
//
// Solidity: function mintBatch(address to, uint256 startId, uint256 count) returns()
func (_MintableNFT *MintableNFTTransactor) MintBatch(opts *bind.TransactOpts, to common.Address, startId *big.Int, count *big.Int) (*types.Transaction, error) {
	return _MintableNFT.contract.Transact(opts, "mintBatch", to, startId, count)
}

// MintBatch is a paid mutator transaction binding the contract method 0x2e81aaea.
//
// Solidity: function mintBatch(address to, uint256 startId, uint256 count) returns()
func (_MintableNFT *MintableNFTSession) MintBatch(to common.Address, startId *big.Int, count *big.Int) (*types.Transaction, error) {
	return _MintableNFT.Contract.MintBatch(&_MintableNFT.TransactOpts, to, startId, count)
}

// MintBatch is a paid mutator transaction binding the contract method 0x2e81aaea.
//
// Solidity: function mintBatch(address to, uint256 startId, uint256 count) returns()
func (_MintableNFT *MintableNFTTransactorSession) MintBatch(to common.Address, startId *big.Int, count *big.Int) (*types.Transaction, error) {
	return _MintableNFT.Contract.MintBatch(&_MintableNFT.TransactOpts, to, startId, count)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_MintableNFT *MintableNFTTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _MintableNFT.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_MintableNFT *MintableNFTSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _MintableNFT.Contract.SetApprovalForAll(&_MintableNFT.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_MintableNFT *MintableNFTTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _MintableNFT.Contract.SetApprovalForAll(&_MintableNFT.TransactOpts, operator, approved)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_MintableNFT *MintableNFTTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _MintableNFT.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_MintableNFT *MintableNFTSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _MintableNFT.Contract.TransferFrom(&_MintableNFT.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_MintableNFT *MintableNFTTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _MintableNFT.Contract.TransferFrom(&_MintableNFT.TransactOpts, from, to, tokenId)
}

// MintableNFTApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the MintableNFT contract.
type MintableNFTApprovalIterator struct {
	Event *MintableNFTApproval // Event containing the contract specifics and raw log

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
func (it *MintableNFTApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MintableNFTApproval)
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
		it.Event = new(MintableNFTApproval)
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
func (it *MintableNFTApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MintableNFTApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MintableNFTApproval represents a Approval event raised by the MintableNFT contract.
type MintableNFTApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_MintableNFT *MintableNFTFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*MintableNFTApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _MintableNFT.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &MintableNFTApprovalIterator{contract: _MintableNFT.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_MintableNFT *MintableNFTFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *MintableNFTApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _MintableNFT.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MintableNFTApproval)
				if err := _MintableNFT.contract.UnpackLog(event, "Approval", log); err != nil {
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
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_MintableNFT *MintableNFTFilterer) ParseApproval(log types.Log) (*MintableNFTApproval, error) {
	event := new(MintableNFTApproval)
	if err := _MintableNFT.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MintableNFTApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the MintableNFT contract.
type MintableNFTApprovalForAllIterator struct {
	Event *MintableNFTApprovalForAll // Event containing the contract specifics and raw log

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
func (it *MintableNFTApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MintableNFTApprovalForAll)
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
		it.Event = new(MintableNFTApprovalForAll)
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
func (it *MintableNFTApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MintableNFTApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MintableNFTApprovalForAll represents a ApprovalForAll event raised by the MintableNFT contract.
type MintableNFTApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_MintableNFT *MintableNFTFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*MintableNFTApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _MintableNFT.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &MintableNFTApprovalForAllIterator{contract: _MintableNFT.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_MintableNFT *MintableNFTFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *MintableNFTApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _MintableNFT.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MintableNFTApprovalForAll)
				if err := _MintableNFT.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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

// ParseApprovalForAll is a log parse operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_MintableNFT *MintableNFTFilterer) ParseApprovalForAll(log types.Log) (*MintableNFTApprovalForAll, error) {
	event := new(MintableNFTApprovalForAll)
	if err := _MintableNFT.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MintableNFTTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the MintableNFT contract.
type MintableNFTTransferIterator struct {
	Event *MintableNFTTransfer // Event containing the contract specifics and raw log

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
func (it *MintableNFTTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MintableNFTTransfer)
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
		it.Event = new(MintableNFTTransfer)
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
func (it *MintableNFTTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MintableNFTTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MintableNFTTransfer represents a Transfer event raised by the MintableNFT contract.
type MintableNFTTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_MintableNFT *MintableNFTFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*MintableNFTTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _MintableNFT.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &MintableNFTTransferIterator{contract: _MintableNFT.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_MintableNFT *MintableNFTFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *MintableNFTTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _MintableNFT.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MintableNFTTransfer)
				if err := _MintableNFT.contract.UnpackLog(event, "Transfer", log); err != nil {
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
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_MintableNFT *MintableNFTFilterer) ParseTransfer(log types.Log) (*MintableNFTTransfer, error) {
	event := new(MintableNFTTransfer)
	if err := _MintableNFT.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
