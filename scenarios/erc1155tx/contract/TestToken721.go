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

// TestToken721MetaData contains all meta data concerning the TestToken721 contract.
var TestToken721MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferMint\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801562000010575f80fd5b506040518060400160405280600a81526020016914dc185b5bdbdc93919560b21b8152506040518060400160405280600481526020016314d3919560e21b815250815f908162000061919062000117565b50600162000070828262000117565b505050620001e3565b634e487b7160e01b5f52604160045260245ffd5b600181811c90821680620000a257607f821691505b602082108103620000c157634e487b7160e01b5f52602260045260245ffd5b50919050565b601f8211156200011257805f5260205f20601f840160051c81016020851015620000ee5750805b601f840160051c820191505b818110156200010f575f8155600101620000fa565b50505b505050565b81516001600160401b0381111562000133576200013362000079565b6200014b816200014484546200008d565b84620000c7565b602080601f83116001811462000181575f8415620001695750858301515b5f19600386901b1c1916600185901b178555620001db565b5f85815260208120601f198616915b82811015620001b15788860151825594840194600190910190840162000190565b5085821015620001cf57878501515f19600388901b60f8161c191681555b505060018460011b0185555b505050505050565b6112a880620001f15f395ff3fe608060405234801561000f575f80fd5b50600436106100f0575f3560e01c806370a0823111610093578063a22cb46511610063578063a22cb465146101f9578063b88d4fde1461020c578063c87b56dd1461021f578063e985e9c514610232575f80fd5b806370a08231146101aa57806395d89b41146101cb5780639d0f7cba146101d3578063a0712d68146101e6575f80fd5b8063095ea7b3116100ce578063095ea7b31461015c57806323b872dd1461017157806342842e0e146101845780636352211e14610197575f80fd5b806301ffc9a7146100f457806306fdde031461011c578063081812fc14610131575b5f80fd5b610107610102366004610e58565b61026d565b60405190151581526020015b60405180910390f35b6101246102be565b6040516101139190610ec0565b61014461013f366004610ed2565b61034d565b6040516001600160a01b039091168152602001610113565b61016f61016a366004610f04565b610372565b005b61016f61017f366004610f2c565b61048b565b61016f610192366004610f2c565b6104bc565b6101446101a5366004610ed2565b6104d6565b6101bd6101b8366004610f65565b610535565b604051908152602001610113565b6101246105b9565b6101076101e1366004610f04565b6105c8565b61016f6101f4366004610ed2565b6105e9565b61016f610207366004610f7e565b6105f6565b61016f61021a366004610fcb565b610605565b61012461022d366004610ed2565b61063d565b6101076102403660046110a0565b6001600160a01b039182165f90815260056020908152604080832093909416825291909152205460ff1690565b5f6001600160e01b031982166380ac58cd60e01b148061029d57506001600160e01b03198216635b5e139f60e01b145b806102b857506301ffc9a760e01b6001600160e01b03198316145b92915050565b60605f80546102cc906110d1565b80601f01602080910402602001604051908101604052809291908181526020018280546102f8906110d1565b80156103435780601f1061031a57610100808354040283529160200191610343565b820191905f5260205f20905b81548152906001019060200180831161032657829003601f168201915b5050505050905090565b5f610357826106ad565b505f908152600460205260409020546001600160a01b031690565b5f61037c826104d6565b9050806001600160a01b0316836001600160a01b0316036103ee5760405162461bcd60e51b815260206004820152602160248201527f4552433732313a20617070726f76616c20746f2063757272656e74206f776e656044820152603960f91b60648201526084015b60405180910390fd5b336001600160a01b038216148061040a575061040a8133610240565b61047c5760405162461bcd60e51b815260206004820152603d60248201527f4552433732313a20617070726f76652063616c6c6572206973206e6f7420746f60448201527f6b656e206f776e6572206f7220617070726f76656420666f7220616c6c00000060648201526084016103e5565b610486838361070b565b505050565b6104953382610778565b6104b15760405162461bcd60e51b81526004016103e590611109565b6104868383836107f5565b61048683838360405180602001604052805f815250610605565b5f818152600260205260408120546001600160a01b0316806102b85760405162461bcd60e51b8152602060048201526018602482015277115490cdcc8c4e881a5b9d985b1a59081d1bdad95b88125160421b60448201526064016103e5565b5f6001600160a01b03821661059e5760405162461bcd60e51b815260206004820152602960248201527f4552433732313a2061646472657373207a65726f206973206e6f7420612076616044820152683634b21037bbb732b960b91b60648201526084016103e5565b506001600160a01b03165f9081526003602052604090205490565b6060600180546102cc906110d1565b5f336105d48184610957565b6105df8185856107f5565b5060019392505050565b6105f33382610957565b50565b610601338383610adf565b5050565b61060f3383610778565b61062b5760405162461bcd60e51b81526004016103e590611109565b61063784848484610bac565b50505050565b6060610648826106ad565b5f61065d60408051602081019091525f815290565b90505f81511161067b5760405180602001604052805f8152506106a6565b8061068584610bdf565b604051602001610696929190611156565b6040516020818303038152906040525b9392505050565b5f818152600260205260409020546001600160a01b03166105f35760405162461bcd60e51b8152602060048201526018602482015277115490cdcc8c4e881a5b9d985b1a59081d1bdad95b88125160421b60448201526064016103e5565b5f81815260046020526040902080546001600160a01b0319166001600160a01b038416908117909155819061073f826104d6565b6001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92560405160405180910390a45050565b5f80610783836104d6565b9050806001600160a01b0316846001600160a01b031614806107c957506001600160a01b038082165f9081526005602090815260408083209388168352929052205460ff165b806107ed5750836001600160a01b03166107e28461034d565b6001600160a01b0316145b949350505050565b826001600160a01b0316610808826104d6565b6001600160a01b03161461082e5760405162461bcd60e51b81526004016103e590611184565b6001600160a01b0382166108905760405162461bcd60e51b8152602060048201526024808201527f4552433732313a207472616e7366657220746f20746865207a65726f206164646044820152637265737360e01b60648201526084016103e5565b826001600160a01b03166108a3826104d6565b6001600160a01b0316146108c95760405162461bcd60e51b81526004016103e590611184565b5f81815260046020908152604080832080546001600160a01b03199081169091556001600160a01b038781168086526003855283862080545f1901905590871680865283862080546001019055868652600290945282852080549092168417909155905184937fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef91a4505050565b6001600160a01b0382166109ad5760405162461bcd60e51b815260206004820181905260248201527f4552433732313a206d696e7420746f20746865207a65726f206164647265737360448201526064016103e5565b5f818152600260205260409020546001600160a01b031615610a115760405162461bcd60e51b815260206004820152601c60248201527f4552433732313a20746f6b656e20616c7265616479206d696e7465640000000060448201526064016103e5565b5f818152600260205260409020546001600160a01b031615610a755760405162461bcd60e51b815260206004820152601c60248201527f4552433732313a20746f6b656e20616c7265616479206d696e7465640000000060448201526064016103e5565b6001600160a01b0382165f81815260036020908152604080832080546001019055848352600290915280822080546001600160a01b0319168417905551839291907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef908290a45050565b816001600160a01b0316836001600160a01b031603610b405760405162461bcd60e51b815260206004820152601960248201527f4552433732313a20617070726f766520746f2063616c6c65720000000000000060448201526064016103e5565b6001600160a01b038381165f81815260056020908152604080832094871680845294825291829020805460ff191686151590811790915591519182527f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31910160405180910390a3505050565b610bb78484846107f5565b610bc384848484610c6f565b6106375760405162461bcd60e51b81526004016103e5906111c9565b60605f610beb83610d6c565b60010190505f8167ffffffffffffffff811115610c0a57610c0a610fb7565b6040519080825280601f01601f191660200182016040528015610c34576020820181803683370190505b5090508181016020015b5f19016f181899199a1a9b1b9c1cb0b131b232b360811b600a86061a8153600a8504945084610c3e57509392505050565b5f6001600160a01b0384163b15610d6157604051630a85bd0160e11b81526001600160a01b0385169063150b7a0290610cb290339089908890889060040161121b565b6020604051808303815f875af1925050508015610cec575060408051601f3d908101601f19168201909252610ce991810190611257565b60015b610d47573d808015610d19576040519150601f19603f3d011682016040523d82523d5f602084013e610d1e565b606091505b5080515f03610d3f5760405162461bcd60e51b81526004016103e5906111c9565b805181602001fd5b6001600160e01b031916630a85bd0160e11b1490506107ed565b506001949350505050565b5f8072184f03e93ff9f4daa797ed6e38ed64bf6a1f0160401b8310610daa5772184f03e93ff9f4daa797ed6e38ed64bf6a1f0160401b830492506040015b6d04ee2d6d415b85acef81000000008310610dd6576d04ee2d6d415b85acef8100000000830492506020015b662386f26fc100008310610df457662386f26fc10000830492506010015b6305f5e1008310610e0c576305f5e100830492506008015b6127108310610e2057612710830492506004015b60648310610e32576064830492506002015b600a83106102b85760010192915050565b6001600160e01b0319811681146105f3575f80fd5b5f60208284031215610e68575f80fd5b81356106a681610e43565b5f5b83811015610e8d578181015183820152602001610e75565b50505f910152565b5f8151808452610eac816020860160208601610e73565b601f01601f19169290920160200192915050565b602081525f6106a66020830184610e95565b5f60208284031215610ee2575f80fd5b5035919050565b80356001600160a01b0381168114610eff575f80fd5b919050565b5f8060408385031215610f15575f80fd5b610f1e83610ee9565b946020939093013593505050565b5f805f60608486031215610f3e575f80fd5b610f4784610ee9565b9250610f5560208501610ee9565b9150604084013590509250925092565b5f60208284031215610f75575f80fd5b6106a682610ee9565b5f8060408385031215610f8f575f80fd5b610f9883610ee9565b915060208301358015158114610fac575f80fd5b809150509250929050565b634e487b7160e01b5f52604160045260245ffd5b5f805f8060808587031215610fde575f80fd5b610fe785610ee9565b9350610ff560208601610ee9565b925060408501359150606085013567ffffffffffffffff80821115611018575f80fd5b818701915087601f83011261102b575f80fd5b81358181111561103d5761103d610fb7565b604051601f8201601f19908116603f0116810190838211818310171561106557611065610fb7565b816040528281528a602084870101111561107d575f80fd5b826020860160208301375f60208483010152809550505050505092959194509250565b5f80604083850312156110b1575f80fd5b6110ba83610ee9565b91506110c860208401610ee9565b90509250929050565b600181811c908216806110e557607f821691505b60208210810361110357634e487b7160e01b5f52602260045260245ffd5b50919050565b6020808252602d908201527f4552433732313a2063616c6c6572206973206e6f7420746f6b656e206f776e6560408201526c1c881bdc88185c1c1c9bdd9959609a1b606082015260800190565b5f8351611167818460208801610e73565b83519083019061117b818360208801610e73565b01949350505050565b60208082526025908201527f4552433732313a207472616e736665722066726f6d20696e636f72726563742060408201526437bbb732b960d91b606082015260800190565b60208082526032908201527f4552433732313a207472616e7366657220746f206e6f6e20455243373231526560408201527131b2b4bb32b91034b6b83632b6b2b73a32b960711b606082015260800190565b6001600160a01b03858116825284166020820152604081018390526080606082018190525f9061124d90830184610e95565b9695505050505050565b5f60208284031215611267575f80fd5b81516106a681610e4356fea264697066735822122025b83d37100d56225d133e587ac9dc6f4765fbe9d30661b4d4177003476c329964736f6c63430008160033",
}

// TestToken721ABI is the input ABI used to generate the binding from.
// Deprecated: Use TestToken721MetaData.ABI instead.
var TestToken721ABI = TestToken721MetaData.ABI

// TestToken721Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TestToken721MetaData.Bin instead.
var TestToken721Bin = TestToken721MetaData.Bin

// DeployTestToken721 deploys a new Ethereum contract, binding an instance of TestToken721 to it.
func DeployTestToken721(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TestToken721, error) {
	parsed, err := TestToken721MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TestToken721Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TestToken721{TestToken721Caller: TestToken721Caller{contract: contract}, TestToken721Transactor: TestToken721Transactor{contract: contract}, TestToken721Filterer: TestToken721Filterer{contract: contract}}, nil
}

// TestToken721 is an auto generated Go binding around an Ethereum contract.
type TestToken721 struct {
	TestToken721Caller     // Read-only binding to the contract
	TestToken721Transactor // Write-only binding to the contract
	TestToken721Filterer   // Log filterer for contract events
}

// TestToken721Caller is an auto generated read-only Go binding around an Ethereum contract.
type TestToken721Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestToken721Transactor is an auto generated write-only Go binding around an Ethereum contract.
type TestToken721Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestToken721Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestToken721Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestToken721Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestToken721Session struct {
	Contract     *TestToken721     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestToken721CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestToken721CallerSession struct {
	Contract *TestToken721Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// TestToken721TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestToken721TransactorSession struct {
	Contract     *TestToken721Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// TestToken721Raw is an auto generated low-level Go binding around an Ethereum contract.
type TestToken721Raw struct {
	Contract *TestToken721 // Generic contract binding to access the raw methods on
}

// TestToken721CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestToken721CallerRaw struct {
	Contract *TestToken721Caller // Generic read-only contract binding to access the raw methods on
}

// TestToken721TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestToken721TransactorRaw struct {
	Contract *TestToken721Transactor // Generic write-only contract binding to access the raw methods on
}

// NewTestToken721 creates a new instance of TestToken721, bound to a specific deployed contract.
func NewTestToken721(address common.Address, backend bind.ContractBackend) (*TestToken721, error) {
	contract, err := bindTestToken721(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestToken721{TestToken721Caller: TestToken721Caller{contract: contract}, TestToken721Transactor: TestToken721Transactor{contract: contract}, TestToken721Filterer: TestToken721Filterer{contract: contract}}, nil
}

// NewTestToken721Caller creates a new read-only instance of TestToken721, bound to a specific deployed contract.
func NewTestToken721Caller(address common.Address, caller bind.ContractCaller) (*TestToken721Caller, error) {
	contract, err := bindTestToken721(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestToken721Caller{contract: contract}, nil
}

// NewTestToken721Transactor creates a new write-only instance of TestToken721, bound to a specific deployed contract.
func NewTestToken721Transactor(address common.Address, transactor bind.ContractTransactor) (*TestToken721Transactor, error) {
	contract, err := bindTestToken721(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestToken721Transactor{contract: contract}, nil
}

// NewTestToken721Filterer creates a new log filterer instance of TestToken721, bound to a specific deployed contract.
func NewTestToken721Filterer(address common.Address, filterer bind.ContractFilterer) (*TestToken721Filterer, error) {
	contract, err := bindTestToken721(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestToken721Filterer{contract: contract}, nil
}

// bindTestToken721 binds a generic wrapper to an already deployed contract.
func bindTestToken721(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TestToken721MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestToken721 *TestToken721Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestToken721.Contract.TestToken721Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestToken721 *TestToken721Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestToken721.Contract.TestToken721Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestToken721 *TestToken721Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestToken721.Contract.TestToken721Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestToken721 *TestToken721CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestToken721.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestToken721 *TestToken721TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestToken721.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestToken721 *TestToken721TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestToken721.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_TestToken721 *TestToken721Caller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TestToken721.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_TestToken721 *TestToken721Session) BalanceOf(owner common.Address) (*big.Int, error) {
	return _TestToken721.Contract.BalanceOf(&_TestToken721.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_TestToken721 *TestToken721CallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _TestToken721.Contract.BalanceOf(&_TestToken721.CallOpts, owner)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_TestToken721 *TestToken721Caller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _TestToken721.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_TestToken721 *TestToken721Session) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _TestToken721.Contract.GetApproved(&_TestToken721.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_TestToken721 *TestToken721CallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _TestToken721.Contract.GetApproved(&_TestToken721.CallOpts, tokenId)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_TestToken721 *TestToken721Caller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _TestToken721.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_TestToken721 *TestToken721Session) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _TestToken721.Contract.IsApprovedForAll(&_TestToken721.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_TestToken721 *TestToken721CallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _TestToken721.Contract.IsApprovedForAll(&_TestToken721.CallOpts, owner, operator)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TestToken721 *TestToken721Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TestToken721.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TestToken721 *TestToken721Session) Name() (string, error) {
	return _TestToken721.Contract.Name(&_TestToken721.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TestToken721 *TestToken721CallerSession) Name() (string, error) {
	return _TestToken721.Contract.Name(&_TestToken721.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_TestToken721 *TestToken721Caller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _TestToken721.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_TestToken721 *TestToken721Session) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _TestToken721.Contract.OwnerOf(&_TestToken721.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_TestToken721 *TestToken721CallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _TestToken721.Contract.OwnerOf(&_TestToken721.CallOpts, tokenId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_TestToken721 *TestToken721Caller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _TestToken721.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_TestToken721 *TestToken721Session) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _TestToken721.Contract.SupportsInterface(&_TestToken721.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_TestToken721 *TestToken721CallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _TestToken721.Contract.SupportsInterface(&_TestToken721.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TestToken721 *TestToken721Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TestToken721.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TestToken721 *TestToken721Session) Symbol() (string, error) {
	return _TestToken721.Contract.Symbol(&_TestToken721.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TestToken721 *TestToken721CallerSession) Symbol() (string, error) {
	return _TestToken721.Contract.Symbol(&_TestToken721.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_TestToken721 *TestToken721Caller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _TestToken721.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_TestToken721 *TestToken721Session) TokenURI(tokenId *big.Int) (string, error) {
	return _TestToken721.Contract.TokenURI(&_TestToken721.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_TestToken721 *TestToken721CallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _TestToken721.Contract.TokenURI(&_TestToken721.CallOpts, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_TestToken721 *TestToken721Transactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _TestToken721.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_TestToken721 *TestToken721Session) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _TestToken721.Contract.Approve(&_TestToken721.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_TestToken721 *TestToken721TransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _TestToken721.Contract.Approve(&_TestToken721.TransactOpts, to, tokenId)
}

// Mint is a paid mutator transaction binding the contract method 0xa0712d68.
//
// Solidity: function mint(uint256 amount) returns()
func (_TestToken721 *TestToken721Transactor) Mint(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _TestToken721.contract.Transact(opts, "mint", amount)
}

// Mint is a paid mutator transaction binding the contract method 0xa0712d68.
//
// Solidity: function mint(uint256 amount) returns()
func (_TestToken721 *TestToken721Session) Mint(amount *big.Int) (*types.Transaction, error) {
	return _TestToken721.Contract.Mint(&_TestToken721.TransactOpts, amount)
}

// Mint is a paid mutator transaction binding the contract method 0xa0712d68.
//
// Solidity: function mint(uint256 amount) returns()
func (_TestToken721 *TestToken721TransactorSession) Mint(amount *big.Int) (*types.Transaction, error) {
	return _TestToken721.Contract.Mint(&_TestToken721.TransactOpts, amount)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_TestToken721 *TestToken721Transactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _TestToken721.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_TestToken721 *TestToken721Session) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _TestToken721.Contract.SafeTransferFrom(&_TestToken721.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_TestToken721 *TestToken721TransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _TestToken721.Contract.SafeTransferFrom(&_TestToken721.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_TestToken721 *TestToken721Transactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _TestToken721.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_TestToken721 *TestToken721Session) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _TestToken721.Contract.SafeTransferFrom0(&_TestToken721.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_TestToken721 *TestToken721TransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _TestToken721.Contract.SafeTransferFrom0(&_TestToken721.TransactOpts, from, to, tokenId, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_TestToken721 *TestToken721Transactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _TestToken721.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_TestToken721 *TestToken721Session) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _TestToken721.Contract.SetApprovalForAll(&_TestToken721.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_TestToken721 *TestToken721TransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _TestToken721.Contract.SetApprovalForAll(&_TestToken721.TransactOpts, operator, approved)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_TestToken721 *TestToken721Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _TestToken721.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_TestToken721 *TestToken721Session) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _TestToken721.Contract.TransferFrom(&_TestToken721.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_TestToken721 *TestToken721TransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _TestToken721.Contract.TransferFrom(&_TestToken721.TransactOpts, from, to, tokenId)
}

// TransferMint is a paid mutator transaction binding the contract method 0x9d0f7cba.
//
// Solidity: function transferMint(address recipient, uint256 amount) returns(bool)
func (_TestToken721 *TestToken721Transactor) TransferMint(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestToken721.contract.Transact(opts, "transferMint", recipient, amount)
}

// TransferMint is a paid mutator transaction binding the contract method 0x9d0f7cba.
//
// Solidity: function transferMint(address recipient, uint256 amount) returns(bool)
func (_TestToken721 *TestToken721Session) TransferMint(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestToken721.Contract.TransferMint(&_TestToken721.TransactOpts, recipient, amount)
}

// TransferMint is a paid mutator transaction binding the contract method 0x9d0f7cba.
//
// Solidity: function transferMint(address recipient, uint256 amount) returns(bool)
func (_TestToken721 *TestToken721TransactorSession) TransferMint(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestToken721.Contract.TransferMint(&_TestToken721.TransactOpts, recipient, amount)
}

// TestToken721ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the TestToken721 contract.
type TestToken721ApprovalIterator struct {
	Event *TestToken721Approval // Event containing the contract specifics and raw log

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
func (it *TestToken721ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestToken721Approval)
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
		it.Event = new(TestToken721Approval)
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
func (it *TestToken721ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestToken721ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestToken721Approval represents a Approval event raised by the TestToken721 contract.
type TestToken721Approval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_TestToken721 *TestToken721Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*TestToken721ApprovalIterator, error) {

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

	logs, sub, err := _TestToken721.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &TestToken721ApprovalIterator{contract: _TestToken721.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_TestToken721 *TestToken721Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *TestToken721Approval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _TestToken721.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestToken721Approval)
				if err := _TestToken721.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_TestToken721 *TestToken721Filterer) ParseApproval(log types.Log) (*TestToken721Approval, error) {
	event := new(TestToken721Approval)
	if err := _TestToken721.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TestToken721ApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the TestToken721 contract.
type TestToken721ApprovalForAllIterator struct {
	Event *TestToken721ApprovalForAll // Event containing the contract specifics and raw log

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
func (it *TestToken721ApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestToken721ApprovalForAll)
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
		it.Event = new(TestToken721ApprovalForAll)
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
func (it *TestToken721ApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestToken721ApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestToken721ApprovalForAll represents a ApprovalForAll event raised by the TestToken721 contract.
type TestToken721ApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_TestToken721 *TestToken721Filterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*TestToken721ApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _TestToken721.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &TestToken721ApprovalForAllIterator{contract: _TestToken721.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_TestToken721 *TestToken721Filterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *TestToken721ApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _TestToken721.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestToken721ApprovalForAll)
				if err := _TestToken721.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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
func (_TestToken721 *TestToken721Filterer) ParseApprovalForAll(log types.Log) (*TestToken721ApprovalForAll, error) {
	event := new(TestToken721ApprovalForAll)
	if err := _TestToken721.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TestToken721TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the TestToken721 contract.
type TestToken721TransferIterator struct {
	Event *TestToken721Transfer // Event containing the contract specifics and raw log

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
func (it *TestToken721TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestToken721Transfer)
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
		it.Event = new(TestToken721Transfer)
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
func (it *TestToken721TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestToken721TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestToken721Transfer represents a Transfer event raised by the TestToken721 contract.
type TestToken721Transfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_TestToken721 *TestToken721Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*TestToken721TransferIterator, error) {

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

	logs, sub, err := _TestToken721.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &TestToken721TransferIterator{contract: _TestToken721.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_TestToken721 *TestToken721Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *TestToken721Transfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _TestToken721.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestToken721Transfer)
				if err := _TestToken721.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_TestToken721 *TestToken721Filterer) ParseTransfer(log types.Log) (*TestToken721Transfer, error) {
	event := new(TestToken721Transfer)
	if err := _TestToken721.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
