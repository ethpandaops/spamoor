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

// TestToken1155MetaData contains all meta data concerning the TestToken1155 contract.
var TestToken1155MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"}],\"name\":\"TransferBatch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"TransferSingle\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"value\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"URI\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"}],\"name\":\"balanceOfBatch\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"id\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amount\",\"type\":\"uint256[]\"}],\"name\":\"mintBatch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeBatchTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferMint\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"uri\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801562000010575f80fd5b5060408051808201909152600881526714185b991853919560c21b60208201526200003b8162000042565b50620001be565b6002620000508282620000f2565b5050565b634e487b7160e01b5f52604160045260245ffd5b600181811c908216806200007d57607f821691505b6020821081036200009c57634e487b7160e01b5f52602260045260245ffd5b50919050565b601f821115620000ed57805f5260205f20601f840160051c81016020851015620000c95750805b601f840160051c820191505b81811015620000ea575f8155600101620000d5565b50505b505050565b81516001600160401b038111156200010e576200010e62000054565b62000126816200011f845462000068565b84620000a2565b602080601f8311600181146200015c575f8415620001445750858301515b5f19600386901b1c1916600185901b178555620001b6565b5f85815260208120601f198616915b828110156200018c578886015182559484019460019091019084016200016b565b5085821015620001aa57878501515f19600388901b60f8161c191681555b505060018460011b0185555b505050505050565b6116ca80620001cc5f395ff3fe608060405234801561000f575f80fd5b50600436106100a5575f3560e01c80634e1273f41161006e5780634e1273f41461013a578063a22cb4651461015a578063d81d0a151461016d578063e985e9c514610180578063f242432a146101bb578063f8541991146101ce575f80fd5b8062fdd58e146100a957806301ffc9a7146100cf5780630e89341c146100f25780631b2ef1ca146101125780632eb2c2d614610127575b5f80fd5b6100bc6100b7366004610da5565b6101e1565b6040519081526020015b60405180910390f35b6100e26100dd366004610de5565b610278565b60405190151581526020016100c6565b610105610100366004610e00565b6102c7565b6040516100c69190610e5a565b610125610120366004610e6c565b610359565b005b610125610135366004610fd3565b610377565b61014d610148366004611076565b6103c3565b6040516100c69190611175565b610125610168366004611187565b6104e3565b6100e261017b3660046111c0565b6104ee565b6100e261018e36600461122f565b6001600160a01b039182165f90815260016020908152604080832093909416825291909152205460ff1690565b6101256101c9366004611260565b610510565b6100e26101dc3660046112c0565b610555565b5f6001600160a01b0383166102505760405162461bcd60e51b815260206004820152602a60248201527f455243313135353a2061646472657373207a65726f206973206e6f742061207660448201526930b634b21037bbb732b960b11b60648201526084015b60405180910390fd5b505f818152602081815260408083206001600160a01b03861684529091529020545b92915050565b5f6001600160e01b03198216636cdb3d1360e11b14806102a857506001600160e01b031982166303a24d0760e21b145b8061027257506301ffc9a760e01b6001600160e01b0319831614610272565b6060600280546102d6906112f0565b80601f0160208091040260200160405190810160405280929190818152602001828054610302906112f0565b801561034d5780601f106103245761010080835404028352916020019161034d565b820191905f5260205f20905b81548152906001019060200180831161033057829003601f168201915b50505050509050919050565b61037333838360405180602001604052805f81525061058b565b5050565b6001600160a01b0385163314806103935750610393853361018e565b6103af5760405162461bcd60e51b815260040161024790611328565b6103bc8585858585610660565b5050505050565b606081518351146104285760405162461bcd60e51b815260206004820152602960248201527f455243313135353a206163636f756e747320616e6420696473206c656e677468604482015268040dad2e6dac2e8c6d60bb1b6064820152608401610247565b5f835167ffffffffffffffff81111561044357610443610e8c565b60405190808252806020026020018201604052801561046c578160200160208202803683370190505b5090505f5b84518110156104db576104b685828151811061048f5761048f611376565b60200260200101518583815181106104a9576104a9611376565b60200260200101516101e1565b8282815181106104c8576104c8611376565b6020908102919091010152600101610471565b509392505050565b6103733383836107f0565b5f61050984848460405180602001604052805f8152506108cf565b9392505050565b6001600160a01b03851633148061052c575061052c853361018e565b6105485760405162461bcd60e51b815260040161024790611328565b6103bc8585858585610a07565b5f61057033848460405180602001604052805f81525061058b565b6105093385858560405180602001604052805f815250610a07565b6001600160a01b0384166105b15760405162461bcd60e51b81526004016102479061138a565b335f6105bc85610b2d565b90505f6105c885610b2d565b90505f868152602081815260408083206001600160a01b038b168452909152812080548792906105f99084906113cb565b909155505060408051878152602081018790526001600160a01b03808a16925f92918716917fc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62910160405180910390a4610657835f89898989610b76565b50505050505050565b81518351146106815760405162461bcd60e51b8152600401610247906113ea565b6001600160a01b0384166106a75760405162461bcd60e51b815260040161024790611432565b335f5b8451811015610782575f8582815181106106c6576106c6611376565b602002602001015190505f8583815181106106e3576106e3611376565b6020908102919091018101515f84815280835260408082206001600160a01b038e1683529093529190912054909150818110156107325760405162461bcd60e51b815260040161024790611477565b5f838152602081815260408083206001600160a01b038e8116855292528083208585039055908b1682528120805484929061076e9084906113cb565b9091555050600190930192506106aa915050565b50846001600160a01b0316866001600160a01b0316826001600160a01b03167f4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb87876040516107d29291906114c1565b60405180910390a46107e8818787878787610cd0565b505050505050565b816001600160a01b0316836001600160a01b0316036108635760405162461bcd60e51b815260206004820152602960248201527f455243313135353a2073657474696e6720617070726f76616c20737461747573604482015268103337b91039b2b63360b91b6064820152608401610247565b6001600160a01b038381165f81815260016020908152604080832094871680845294825291829020805460ff191686151590811790915591519182527f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31910160405180910390a3505050565b6001600160a01b0384166108f55760405162461bcd60e51b81526004016102479061138a565b81518351146109165760405162461bcd60e51b8152600401610247906113ea565b335f5b84518110156109a15783818151811061093457610934611376565b60200260200101515f8087848151811061095057610950611376565b602002602001015181526020019081526020015f205f886001600160a01b03166001600160a01b031681526020019081526020015f205f82825461099491906113cb565b9091555050600101610919565b50846001600160a01b03165f6001600160a01b0316826001600160a01b03167f4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb87876040516109f19291906114c1565b60405180910390a46103bc815f87878787610cd0565b6001600160a01b038416610a2d5760405162461bcd60e51b815260040161024790611432565b335f610a3885610b2d565b90505f610a4485610b2d565b90505f868152602081815260408083206001600160a01b038c16845290915290205485811015610a865760405162461bcd60e51b815260040161024790611477565b5f878152602081815260408083206001600160a01b038d8116855292528083208985039055908a16825281208054889290610ac29084906113cb565b909155505060408051888152602081018890526001600160a01b03808b16928c821692918816917fc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62910160405180910390a4610b22848a8a8a8a8a610b76565b505050505050505050565b6040805160018082528183019092526060915f91906020808301908036833701905050905082815f81518110610b6557610b65611376565b602090810291909101015292915050565b6001600160a01b0384163b156107e85760405163f23a6e6160e01b81526001600160a01b0385169063f23a6e6190610bba90899089908890889088906004016114ee565b6020604051808303815f875af1925050508015610bf4575060408051601f3d908101601f19168201909252610bf191810190611532565b60015b610ca057610c0061154d565b806308c379a003610c395750610c14611566565b80610c1f5750610c3b565b8060405162461bcd60e51b81526004016102479190610e5a565b505b60405162461bcd60e51b815260206004820152603460248201527f455243313135353a207472616e7366657220746f206e6f6e2d455243313135356044820152732932b1b2b4bb32b91034b6b83632b6b2b73a32b960611b6064820152608401610247565b6001600160e01b0319811663f23a6e6160e01b146106575760405162461bcd60e51b8152600401610247906115ef565b6001600160a01b0384163b156107e85760405163bc197c8160e01b81526001600160a01b0385169063bc197c8190610d149089908990889088908890600401611637565b6020604051808303815f875af1925050508015610d4e575060408051601f3d908101601f19168201909252610d4b91810190611532565b60015b610d5a57610c0061154d565b6001600160e01b0319811663bc197c8160e01b146106575760405162461bcd60e51b8152600401610247906115ef565b80356001600160a01b0381168114610da0575f80fd5b919050565b5f8060408385031215610db6575f80fd5b610dbf83610d8a565b946020939093013593505050565b6001600160e01b031981168114610de2575f80fd5b50565b5f60208284031215610df5575f80fd5b813561050981610dcd565b5f60208284031215610e10575f80fd5b5035919050565b5f81518084525f5b81811015610e3b57602081850181015186830182015201610e1f565b505f602082860101526020601f19601f83011685010191505092915050565b602081525f6105096020830184610e17565b5f8060408385031215610e7d575f80fd5b50508035926020909101359150565b634e487b7160e01b5f52604160045260245ffd5b601f8201601f1916810167ffffffffffffffff81118282101715610ec657610ec6610e8c565b6040525050565b5f67ffffffffffffffff821115610ee657610ee6610e8c565b5060051b60200190565b5f82601f830112610eff575f80fd5b81356020610f0c82610ecd565b604051610f198282610ea0565b80915083815260208101915060208460051b870101935086841115610f3c575f80fd5b602086015b84811015610f585780358352918301918301610f41565b509695505050505050565b5f82601f830112610f72575f80fd5b813567ffffffffffffffff811115610f8c57610f8c610e8c565b604051610fa3601f8301601f191660200182610ea0565b818152846020838601011115610fb7575f80fd5b816020850160208301375f918101602001919091529392505050565b5f805f805f60a08688031215610fe7575f80fd5b610ff086610d8a565b9450610ffe60208701610d8a565b9350604086013567ffffffffffffffff8082111561101a575f80fd5b61102689838a01610ef0565b9450606088013591508082111561103b575f80fd5b61104789838a01610ef0565b9350608088013591508082111561105c575f80fd5b5061106988828901610f63565b9150509295509295909350565b5f8060408385031215611087575f80fd5b823567ffffffffffffffff8082111561109e575f80fd5b818501915085601f8301126110b1575f80fd5b813560206110be82610ecd565b6040516110cb8282610ea0565b83815260059390931b85018201928281019150898411156110ea575f80fd5b948201945b8386101561110f5761110086610d8a565b825294820194908201906110ef565b96505086013592505080821115611124575f80fd5b5061113185828601610ef0565b9150509250929050565b5f815180845260208085019450602084015f5b8381101561116a5781518752958201959082019060010161114e565b509495945050505050565b602081525f610509602083018461113b565b5f8060408385031215611198575f80fd5b6111a183610d8a565b9150602083013580151581146111b5575f80fd5b809150509250929050565b5f805f606084860312156111d2575f80fd5b6111db84610d8a565b9250602084013567ffffffffffffffff808211156111f7575f80fd5b61120387838801610ef0565b93506040860135915080821115611218575f80fd5b5061122586828701610ef0565b9150509250925092565b5f8060408385031215611240575f80fd5b61124983610d8a565b915061125760208401610d8a565b90509250929050565b5f805f805f60a08688031215611274575f80fd5b61127d86610d8a565b945061128b60208701610d8a565b93506040860135925060608601359150608086013567ffffffffffffffff8111156112b4575f80fd5b61106988828901610f63565b5f805f606084860312156112d2575f80fd5b6112db84610d8a565b95602085013595506040909401359392505050565b600181811c9082168061130457607f821691505b60208210810361132257634e487b7160e01b5f52602260045260245ffd5b50919050565b6020808252602e908201527f455243313135353a2063616c6c6572206973206e6f7420746f6b656e206f776e60408201526d195c881bdc88185c1c1c9bdd995960921b606082015260800190565b634e487b7160e01b5f52603260045260245ffd5b60208082526021908201527f455243313135353a206d696e7420746f20746865207a65726f206164647265736040820152607360f81b606082015260800190565b8082018082111561027257634e487b7160e01b5f52601160045260245ffd5b60208082526028908201527f455243313135353a2069647320616e6420616d6f756e7473206c656e677468206040820152670dad2e6dac2e8c6d60c31b606082015260800190565b60208082526025908201527f455243313135353a207472616e7366657220746f20746865207a65726f206164604082015264647265737360d81b606082015260800190565b6020808252602a908201527f455243313135353a20696e73756666696369656e742062616c616e636520666f60408201526939103a3930b739b332b960b11b606082015260800190565b604081525f6114d3604083018561113b565b82810360208401526114e5818561113b565b95945050505050565b6001600160a01b03868116825285166020820152604081018490526060810183905260a0608082018190525f9061152790830184610e17565b979650505050505050565b5f60208284031215611542575f80fd5b815161050981610dcd565b5f60033d11156115635760045f803e505f5160e01c5b90565b5f60443d10156115735790565b6040516003193d81016004833e81513d67ffffffffffffffff81602484011181841117156115a357505050505090565b82850191508151818111156115bb5750505050505090565b843d87010160208285010111156115d55750505050505090565b6115e460208286010187610ea0565b509095945050505050565b60208082526028908201527f455243313135353a204552433131353552656365697665722072656a656374656040820152676420746f6b656e7360c01b606082015260800190565b6001600160a01b0386811682528516602082015260a0604082018190525f906116629083018661113b565b8281036060840152611674818661113b565b905082810360808401526116888185610e17565b9897505050505050505056fea2646970667358221220b38aec1473f3c350543003e772b4417a024d54a8bbbf5d15001d4c618dfde1f864736f6c63430008160033",
}

// TestToken1155ABI is the input ABI used to generate the binding from.
// Deprecated: Use TestToken1155MetaData.ABI instead.
var TestToken1155ABI = TestToken1155MetaData.ABI

// TestToken1155Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TestToken1155MetaData.Bin instead.
var TestToken1155Bin = TestToken1155MetaData.Bin

// DeployTestToken1155 deploys a new Ethereum contract, binding an instance of TestToken1155 to it.
func DeployTestToken1155(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TestToken1155, error) {
	parsed, err := TestToken1155MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TestToken1155Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TestToken1155{TestToken1155Caller: TestToken1155Caller{contract: contract}, TestToken1155Transactor: TestToken1155Transactor{contract: contract}, TestToken1155Filterer: TestToken1155Filterer{contract: contract}}, nil
}

// TestToken1155 is an auto generated Go binding around an Ethereum contract.
type TestToken1155 struct {
	TestToken1155Caller     // Read-only binding to the contract
	TestToken1155Transactor // Write-only binding to the contract
	TestToken1155Filterer   // Log filterer for contract events
}

// TestToken1155Caller is an auto generated read-only Go binding around an Ethereum contract.
type TestToken1155Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestToken1155Transactor is an auto generated write-only Go binding around an Ethereum contract.
type TestToken1155Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestToken1155Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestToken1155Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestToken1155Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestToken1155Session struct {
	Contract     *TestToken1155    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestToken1155CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestToken1155CallerSession struct {
	Contract *TestToken1155Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// TestToken1155TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestToken1155TransactorSession struct {
	Contract     *TestToken1155Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// TestToken1155Raw is an auto generated low-level Go binding around an Ethereum contract.
type TestToken1155Raw struct {
	Contract *TestToken1155 // Generic contract binding to access the raw methods on
}

// TestToken1155CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestToken1155CallerRaw struct {
	Contract *TestToken1155Caller // Generic read-only contract binding to access the raw methods on
}

// TestToken1155TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestToken1155TransactorRaw struct {
	Contract *TestToken1155Transactor // Generic write-only contract binding to access the raw methods on
}

// NewTestToken1155 creates a new instance of TestToken1155, bound to a specific deployed contract.
func NewTestToken1155(address common.Address, backend bind.ContractBackend) (*TestToken1155, error) {
	contract, err := bindTestToken1155(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestToken1155{TestToken1155Caller: TestToken1155Caller{contract: contract}, TestToken1155Transactor: TestToken1155Transactor{contract: contract}, TestToken1155Filterer: TestToken1155Filterer{contract: contract}}, nil
}

// NewTestToken1155Caller creates a new read-only instance of TestToken1155, bound to a specific deployed contract.
func NewTestToken1155Caller(address common.Address, caller bind.ContractCaller) (*TestToken1155Caller, error) {
	contract, err := bindTestToken1155(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestToken1155Caller{contract: contract}, nil
}

// NewTestToken1155Transactor creates a new write-only instance of TestToken1155, bound to a specific deployed contract.
func NewTestToken1155Transactor(address common.Address, transactor bind.ContractTransactor) (*TestToken1155Transactor, error) {
	contract, err := bindTestToken1155(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestToken1155Transactor{contract: contract}, nil
}

// NewTestToken1155Filterer creates a new log filterer instance of TestToken1155, bound to a specific deployed contract.
func NewTestToken1155Filterer(address common.Address, filterer bind.ContractFilterer) (*TestToken1155Filterer, error) {
	contract, err := bindTestToken1155(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestToken1155Filterer{contract: contract}, nil
}

// bindTestToken1155 binds a generic wrapper to an already deployed contract.
func bindTestToken1155(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TestToken1155MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestToken1155 *TestToken1155Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestToken1155.Contract.TestToken1155Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestToken1155 *TestToken1155Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestToken1155.Contract.TestToken1155Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestToken1155 *TestToken1155Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestToken1155.Contract.TestToken1155Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestToken1155 *TestToken1155CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestToken1155.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestToken1155 *TestToken1155TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestToken1155.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestToken1155 *TestToken1155TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestToken1155.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x00fdd58e.
//
// Solidity: function balanceOf(address account, uint256 id) view returns(uint256)
func (_TestToken1155 *TestToken1155Caller) BalanceOf(opts *bind.CallOpts, account common.Address, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _TestToken1155.contract.Call(opts, &out, "balanceOf", account, id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x00fdd58e.
//
// Solidity: function balanceOf(address account, uint256 id) view returns(uint256)
func (_TestToken1155 *TestToken1155Session) BalanceOf(account common.Address, id *big.Int) (*big.Int, error) {
	return _TestToken1155.Contract.BalanceOf(&_TestToken1155.CallOpts, account, id)
}

// BalanceOf is a free data retrieval call binding the contract method 0x00fdd58e.
//
// Solidity: function balanceOf(address account, uint256 id) view returns(uint256)
func (_TestToken1155 *TestToken1155CallerSession) BalanceOf(account common.Address, id *big.Int) (*big.Int, error) {
	return _TestToken1155.Contract.BalanceOf(&_TestToken1155.CallOpts, account, id)
}

// BalanceOfBatch is a free data retrieval call binding the contract method 0x4e1273f4.
//
// Solidity: function balanceOfBatch(address[] accounts, uint256[] ids) view returns(uint256[])
func (_TestToken1155 *TestToken1155Caller) BalanceOfBatch(opts *bind.CallOpts, accounts []common.Address, ids []*big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _TestToken1155.contract.Call(opts, &out, "balanceOfBatch", accounts, ids)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// BalanceOfBatch is a free data retrieval call binding the contract method 0x4e1273f4.
//
// Solidity: function balanceOfBatch(address[] accounts, uint256[] ids) view returns(uint256[])
func (_TestToken1155 *TestToken1155Session) BalanceOfBatch(accounts []common.Address, ids []*big.Int) ([]*big.Int, error) {
	return _TestToken1155.Contract.BalanceOfBatch(&_TestToken1155.CallOpts, accounts, ids)
}

// BalanceOfBatch is a free data retrieval call binding the contract method 0x4e1273f4.
//
// Solidity: function balanceOfBatch(address[] accounts, uint256[] ids) view returns(uint256[])
func (_TestToken1155 *TestToken1155CallerSession) BalanceOfBatch(accounts []common.Address, ids []*big.Int) ([]*big.Int, error) {
	return _TestToken1155.Contract.BalanceOfBatch(&_TestToken1155.CallOpts, accounts, ids)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address account, address operator) view returns(bool)
func (_TestToken1155 *TestToken1155Caller) IsApprovedForAll(opts *bind.CallOpts, account common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _TestToken1155.contract.Call(opts, &out, "isApprovedForAll", account, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address account, address operator) view returns(bool)
func (_TestToken1155 *TestToken1155Session) IsApprovedForAll(account common.Address, operator common.Address) (bool, error) {
	return _TestToken1155.Contract.IsApprovedForAll(&_TestToken1155.CallOpts, account, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address account, address operator) view returns(bool)
func (_TestToken1155 *TestToken1155CallerSession) IsApprovedForAll(account common.Address, operator common.Address) (bool, error) {
	return _TestToken1155.Contract.IsApprovedForAll(&_TestToken1155.CallOpts, account, operator)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_TestToken1155 *TestToken1155Caller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _TestToken1155.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_TestToken1155 *TestToken1155Session) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _TestToken1155.Contract.SupportsInterface(&_TestToken1155.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_TestToken1155 *TestToken1155CallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _TestToken1155.Contract.SupportsInterface(&_TestToken1155.CallOpts, interfaceId)
}

// Uri is a free data retrieval call binding the contract method 0x0e89341c.
//
// Solidity: function uri(uint256 ) view returns(string)
func (_TestToken1155 *TestToken1155Caller) Uri(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var out []interface{}
	err := _TestToken1155.contract.Call(opts, &out, "uri", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Uri is a free data retrieval call binding the contract method 0x0e89341c.
//
// Solidity: function uri(uint256 ) view returns(string)
func (_TestToken1155 *TestToken1155Session) Uri(arg0 *big.Int) (string, error) {
	return _TestToken1155.Contract.Uri(&_TestToken1155.CallOpts, arg0)
}

// Uri is a free data retrieval call binding the contract method 0x0e89341c.
//
// Solidity: function uri(uint256 ) view returns(string)
func (_TestToken1155 *TestToken1155CallerSession) Uri(arg0 *big.Int) (string, error) {
	return _TestToken1155.Contract.Uri(&_TestToken1155.CallOpts, arg0)
}

// Mint is a paid mutator transaction binding the contract method 0x1b2ef1ca.
//
// Solidity: function mint(uint256 id, uint256 amount) returns()
func (_TestToken1155 *TestToken1155Transactor) Mint(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _TestToken1155.contract.Transact(opts, "mint", id, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x1b2ef1ca.
//
// Solidity: function mint(uint256 id, uint256 amount) returns()
func (_TestToken1155 *TestToken1155Session) Mint(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _TestToken1155.Contract.Mint(&_TestToken1155.TransactOpts, id, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x1b2ef1ca.
//
// Solidity: function mint(uint256 id, uint256 amount) returns()
func (_TestToken1155 *TestToken1155TransactorSession) Mint(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _TestToken1155.Contract.Mint(&_TestToken1155.TransactOpts, id, amount)
}

// MintBatch is a paid mutator transaction binding the contract method 0xd81d0a15.
//
// Solidity: function mintBatch(address recipient, uint256[] id, uint256[] amount) returns(bool)
func (_TestToken1155 *TestToken1155Transactor) MintBatch(opts *bind.TransactOpts, recipient common.Address, id []*big.Int, amount []*big.Int) (*types.Transaction, error) {
	return _TestToken1155.contract.Transact(opts, "mintBatch", recipient, id, amount)
}

// MintBatch is a paid mutator transaction binding the contract method 0xd81d0a15.
//
// Solidity: function mintBatch(address recipient, uint256[] id, uint256[] amount) returns(bool)
func (_TestToken1155 *TestToken1155Session) MintBatch(recipient common.Address, id []*big.Int, amount []*big.Int) (*types.Transaction, error) {
	return _TestToken1155.Contract.MintBatch(&_TestToken1155.TransactOpts, recipient, id, amount)
}

// MintBatch is a paid mutator transaction binding the contract method 0xd81d0a15.
//
// Solidity: function mintBatch(address recipient, uint256[] id, uint256[] amount) returns(bool)
func (_TestToken1155 *TestToken1155TransactorSession) MintBatch(recipient common.Address, id []*big.Int, amount []*big.Int) (*types.Transaction, error) {
	return _TestToken1155.Contract.MintBatch(&_TestToken1155.TransactOpts, recipient, id, amount)
}

// SafeBatchTransferFrom is a paid mutator transaction binding the contract method 0x2eb2c2d6.
//
// Solidity: function safeBatchTransferFrom(address from, address to, uint256[] ids, uint256[] amounts, bytes data) returns()
func (_TestToken1155 *TestToken1155Transactor) SafeBatchTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, ids []*big.Int, amounts []*big.Int, data []byte) (*types.Transaction, error) {
	return _TestToken1155.contract.Transact(opts, "safeBatchTransferFrom", from, to, ids, amounts, data)
}

// SafeBatchTransferFrom is a paid mutator transaction binding the contract method 0x2eb2c2d6.
//
// Solidity: function safeBatchTransferFrom(address from, address to, uint256[] ids, uint256[] amounts, bytes data) returns()
func (_TestToken1155 *TestToken1155Session) SafeBatchTransferFrom(from common.Address, to common.Address, ids []*big.Int, amounts []*big.Int, data []byte) (*types.Transaction, error) {
	return _TestToken1155.Contract.SafeBatchTransferFrom(&_TestToken1155.TransactOpts, from, to, ids, amounts, data)
}

// SafeBatchTransferFrom is a paid mutator transaction binding the contract method 0x2eb2c2d6.
//
// Solidity: function safeBatchTransferFrom(address from, address to, uint256[] ids, uint256[] amounts, bytes data) returns()
func (_TestToken1155 *TestToken1155TransactorSession) SafeBatchTransferFrom(from common.Address, to common.Address, ids []*big.Int, amounts []*big.Int, data []byte) (*types.Transaction, error) {
	return _TestToken1155.Contract.SafeBatchTransferFrom(&_TestToken1155.TransactOpts, from, to, ids, amounts, data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xf242432a.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 id, uint256 amount, bytes data) returns()
func (_TestToken1155 *TestToken1155Transactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, id *big.Int, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _TestToken1155.contract.Transact(opts, "safeTransferFrom", from, to, id, amount, data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xf242432a.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 id, uint256 amount, bytes data) returns()
func (_TestToken1155 *TestToken1155Session) SafeTransferFrom(from common.Address, to common.Address, id *big.Int, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _TestToken1155.Contract.SafeTransferFrom(&_TestToken1155.TransactOpts, from, to, id, amount, data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xf242432a.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 id, uint256 amount, bytes data) returns()
func (_TestToken1155 *TestToken1155TransactorSession) SafeTransferFrom(from common.Address, to common.Address, id *big.Int, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _TestToken1155.Contract.SafeTransferFrom(&_TestToken1155.TransactOpts, from, to, id, amount, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_TestToken1155 *TestToken1155Transactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _TestToken1155.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_TestToken1155 *TestToken1155Session) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _TestToken1155.Contract.SetApprovalForAll(&_TestToken1155.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_TestToken1155 *TestToken1155TransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _TestToken1155.Contract.SetApprovalForAll(&_TestToken1155.TransactOpts, operator, approved)
}

// TransferMint is a paid mutator transaction binding the contract method 0xf8541991.
//
// Solidity: function transferMint(address recipient, uint256 id, uint256 amount) returns(bool)
func (_TestToken1155 *TestToken1155Transactor) TransferMint(opts *bind.TransactOpts, recipient common.Address, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _TestToken1155.contract.Transact(opts, "transferMint", recipient, id, amount)
}

// TransferMint is a paid mutator transaction binding the contract method 0xf8541991.
//
// Solidity: function transferMint(address recipient, uint256 id, uint256 amount) returns(bool)
func (_TestToken1155 *TestToken1155Session) TransferMint(recipient common.Address, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _TestToken1155.Contract.TransferMint(&_TestToken1155.TransactOpts, recipient, id, amount)
}

// TransferMint is a paid mutator transaction binding the contract method 0xf8541991.
//
// Solidity: function transferMint(address recipient, uint256 id, uint256 amount) returns(bool)
func (_TestToken1155 *TestToken1155TransactorSession) TransferMint(recipient common.Address, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _TestToken1155.Contract.TransferMint(&_TestToken1155.TransactOpts, recipient, id, amount)
}

// TestToken1155ApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the TestToken1155 contract.
type TestToken1155ApprovalForAllIterator struct {
	Event *TestToken1155ApprovalForAll // Event containing the contract specifics and raw log

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
func (it *TestToken1155ApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestToken1155ApprovalForAll)
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
		it.Event = new(TestToken1155ApprovalForAll)
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
func (it *TestToken1155ApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestToken1155ApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestToken1155ApprovalForAll represents a ApprovalForAll event raised by the TestToken1155 contract.
type TestToken1155ApprovalForAll struct {
	Account  common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed account, address indexed operator, bool approved)
func (_TestToken1155 *TestToken1155Filterer) FilterApprovalForAll(opts *bind.FilterOpts, account []common.Address, operator []common.Address) (*TestToken1155ApprovalForAllIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _TestToken1155.contract.FilterLogs(opts, "ApprovalForAll", accountRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &TestToken1155ApprovalForAllIterator{contract: _TestToken1155.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed account, address indexed operator, bool approved)
func (_TestToken1155 *TestToken1155Filterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *TestToken1155ApprovalForAll, account []common.Address, operator []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _TestToken1155.contract.WatchLogs(opts, "ApprovalForAll", accountRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestToken1155ApprovalForAll)
				if err := _TestToken1155.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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
// Solidity: event ApprovalForAll(address indexed account, address indexed operator, bool approved)
func (_TestToken1155 *TestToken1155Filterer) ParseApprovalForAll(log types.Log) (*TestToken1155ApprovalForAll, error) {
	event := new(TestToken1155ApprovalForAll)
	if err := _TestToken1155.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TestToken1155TransferBatchIterator is returned from FilterTransferBatch and is used to iterate over the raw logs and unpacked data for TransferBatch events raised by the TestToken1155 contract.
type TestToken1155TransferBatchIterator struct {
	Event *TestToken1155TransferBatch // Event containing the contract specifics and raw log

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
func (it *TestToken1155TransferBatchIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestToken1155TransferBatch)
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
		it.Event = new(TestToken1155TransferBatch)
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
func (it *TestToken1155TransferBatchIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestToken1155TransferBatchIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestToken1155TransferBatch represents a TransferBatch event raised by the TestToken1155 contract.
type TestToken1155TransferBatch struct {
	Operator common.Address
	From     common.Address
	To       common.Address
	Ids      []*big.Int
	Values   []*big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTransferBatch is a free log retrieval operation binding the contract event 0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb.
//
// Solidity: event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] values)
func (_TestToken1155 *TestToken1155Filterer) FilterTransferBatch(opts *bind.FilterOpts, operator []common.Address, from []common.Address, to []common.Address) (*TestToken1155TransferBatchIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TestToken1155.contract.FilterLogs(opts, "TransferBatch", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &TestToken1155TransferBatchIterator{contract: _TestToken1155.contract, event: "TransferBatch", logs: logs, sub: sub}, nil
}

// WatchTransferBatch is a free log subscription operation binding the contract event 0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb.
//
// Solidity: event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] values)
func (_TestToken1155 *TestToken1155Filterer) WatchTransferBatch(opts *bind.WatchOpts, sink chan<- *TestToken1155TransferBatch, operator []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TestToken1155.contract.WatchLogs(opts, "TransferBatch", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestToken1155TransferBatch)
				if err := _TestToken1155.contract.UnpackLog(event, "TransferBatch", log); err != nil {
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

// ParseTransferBatch is a log parse operation binding the contract event 0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb.
//
// Solidity: event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] values)
func (_TestToken1155 *TestToken1155Filterer) ParseTransferBatch(log types.Log) (*TestToken1155TransferBatch, error) {
	event := new(TestToken1155TransferBatch)
	if err := _TestToken1155.contract.UnpackLog(event, "TransferBatch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TestToken1155TransferSingleIterator is returned from FilterTransferSingle and is used to iterate over the raw logs and unpacked data for TransferSingle events raised by the TestToken1155 contract.
type TestToken1155TransferSingleIterator struct {
	Event *TestToken1155TransferSingle // Event containing the contract specifics and raw log

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
func (it *TestToken1155TransferSingleIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestToken1155TransferSingle)
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
		it.Event = new(TestToken1155TransferSingle)
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
func (it *TestToken1155TransferSingleIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestToken1155TransferSingleIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestToken1155TransferSingle represents a TransferSingle event raised by the TestToken1155 contract.
type TestToken1155TransferSingle struct {
	Operator common.Address
	From     common.Address
	To       common.Address
	Id       *big.Int
	Value    *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTransferSingle is a free log retrieval operation binding the contract event 0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62.
//
// Solidity: event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value)
func (_TestToken1155 *TestToken1155Filterer) FilterTransferSingle(opts *bind.FilterOpts, operator []common.Address, from []common.Address, to []common.Address) (*TestToken1155TransferSingleIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TestToken1155.contract.FilterLogs(opts, "TransferSingle", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &TestToken1155TransferSingleIterator{contract: _TestToken1155.contract, event: "TransferSingle", logs: logs, sub: sub}, nil
}

// WatchTransferSingle is a free log subscription operation binding the contract event 0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62.
//
// Solidity: event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value)
func (_TestToken1155 *TestToken1155Filterer) WatchTransferSingle(opts *bind.WatchOpts, sink chan<- *TestToken1155TransferSingle, operator []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TestToken1155.contract.WatchLogs(opts, "TransferSingle", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestToken1155TransferSingle)
				if err := _TestToken1155.contract.UnpackLog(event, "TransferSingle", log); err != nil {
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

// ParseTransferSingle is a log parse operation binding the contract event 0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62.
//
// Solidity: event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value)
func (_TestToken1155 *TestToken1155Filterer) ParseTransferSingle(log types.Log) (*TestToken1155TransferSingle, error) {
	event := new(TestToken1155TransferSingle)
	if err := _TestToken1155.contract.UnpackLog(event, "TransferSingle", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TestToken1155URIIterator is returned from FilterURI and is used to iterate over the raw logs and unpacked data for URI events raised by the TestToken1155 contract.
type TestToken1155URIIterator struct {
	Event *TestToken1155URI // Event containing the contract specifics and raw log

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
func (it *TestToken1155URIIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestToken1155URI)
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
		it.Event = new(TestToken1155URI)
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
func (it *TestToken1155URIIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestToken1155URIIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestToken1155URI represents a URI event raised by the TestToken1155 contract.
type TestToken1155URI struct {
	Value string
	Id    *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterURI is a free log retrieval operation binding the contract event 0x6bb7ff708619ba0610cba295a58592e0451dee2622938c8755667688daf3529b.
//
// Solidity: event URI(string value, uint256 indexed id)
func (_TestToken1155 *TestToken1155Filterer) FilterURI(opts *bind.FilterOpts, id []*big.Int) (*TestToken1155URIIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _TestToken1155.contract.FilterLogs(opts, "URI", idRule)
	if err != nil {
		return nil, err
	}
	return &TestToken1155URIIterator{contract: _TestToken1155.contract, event: "URI", logs: logs, sub: sub}, nil
}

// WatchURI is a free log subscription operation binding the contract event 0x6bb7ff708619ba0610cba295a58592e0451dee2622938c8755667688daf3529b.
//
// Solidity: event URI(string value, uint256 indexed id)
func (_TestToken1155 *TestToken1155Filterer) WatchURI(opts *bind.WatchOpts, sink chan<- *TestToken1155URI, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _TestToken1155.contract.WatchLogs(opts, "URI", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestToken1155URI)
				if err := _TestToken1155.contract.UnpackLog(event, "URI", log); err != nil {
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

// ParseURI is a log parse operation binding the contract event 0x6bb7ff708619ba0610cba295a58592e0451dee2622938c8755667688daf3529b.
//
// Solidity: event URI(string value, uint256 indexed id)
func (_TestToken1155 *TestToken1155Filterer) ParseURI(log types.Log) (*TestToken1155URI, error) {
	event := new(TestToken1155URI)
	if err := _TestToken1155.contract.UnpackLog(event, "URI", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
