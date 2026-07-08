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

// ACLManagerMetaData contains all meta data concerning the ACLManager contract.
var ACLManagerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIPoolAddressesProvider\",\"name\":\"provider\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADDRESSES_PROVIDER\",\"outputs\":[{\"internalType\":\"contractIPoolAddressesProvider\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ASSET_LISTING_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"BRIDGE_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"EMERGENCY_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FLASH_BORROWER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"POOL_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"RISK_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"addAssetListingAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"bridge\",\"type\":\"address\"}],\"name\":\"addBridge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"addEmergencyAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"borrower\",\"type\":\"address\"}],\"name\":\"addFlashBorrower\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"addPoolAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"addRiskAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"isAssetListingAdmin\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"bridge\",\"type\":\"address\"}],\"name\":\"isBridge\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"isEmergencyAdmin\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"borrower\",\"type\":\"address\"}],\"name\":\"isFlashBorrower\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"isPoolAdmin\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"isRiskAdmin\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"removeAssetListingAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"bridge\",\"type\":\"address\"}],\"name\":\"removeBridge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"removeEmergencyAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"borrower\",\"type\":\"address\"}],\"name\":\"removeFlashBorrower\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"removePoolAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"removeRiskAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"adminRole\",\"type\":\"bytes32\"}],\"name\":\"setRoleAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b50604051620015cd380380620015cd8339810160408190526200003491620001e3565b806001600160a01b03166080816001600160a01b0316815250506000816001600160a01b0316630e67178c6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200008f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620000b59190620001e3565b604080518082019091526002815261373560f01b60208201529091506001600160a01b038216620001045760405162461bcd60e51b8152600401620000fb91906200020a565b60405180910390fd5b50620001126000826200011a565b505062000262565b6200012682826200012a565b5050565b6000828152602081815260408083206001600160a01b038516845290915290205460ff1662000126576000828152602081815260408083206001600160a01b03851684529091529020805460ff19166001179055620001863390565b6001600160a01b0316816001600160a01b0316837f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45050565b6001600160a01b0381168114620001e057600080fd5b50565b600060208284031215620001f657600080fd5b81516200020381620001ca565b9392505050565b600060208083528351808285015260005b8181101562000239578581018301518582016040015282016200021b565b818111156200024c576000604083870101525b50601f01601f1916929092016040019392505050565b60805161134f6200027e6000396000610252015261134f6000f3fe608060405234801561001057600080fd5b506004361061020b5760003560e01c8063674b5e4d1161012a5780639a2b96f7116100bd578063b5bfddea1161008c578063d547741f11610071578063d547741f1461059e578063f83695cb146105b1578063fa50f297146105c457600080fd5b8063b5bfddea14610550578063b8f6dba71461057757600080fd5b80639a2b96f71461050f5780639ac9d80b14610522578063a217fddf14610535578063a21bce151461053d57600080fd5b80637a9a93f4116100f95780637a9a93f41461044a5780637be53ca11461045d57806391d14854146104b85780639712fdf8146104fc57600080fd5b8063674b5e4d146103d65780636e76fc8f146103e9578063726600ce1461041057806378bb0a431461042357600080fd5b80632500f2b6116101a25780633c5a08e5116101715780633c5a08e5146103625780634f16b425146103755780635577b7a91461039c5780635b9a94e4146103c357600080fd5b80632500f2b614610316578063253cf980146103295780632f2ff15d1461033c57806336568abe1461034f57600080fd5b8063179efb09116101de578063179efb09146102ac5780631e4e0091146102bf57806322650caf146102d2578063248a9ca3146102e557600080fd5b806301ffc9a71461021057806304df017d146102385780630542975c1461024d57806313ee32e014610299575b600080fd5b61022361021e366004611013565b6105d7565b60405190151581526020015b60405180910390f35b61024b61024636600461107e565b610670565b005b6102747f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161022f565b6102236102a736600461107e565b61069d565b61024b6102ba36600461107e565b6106ea565b61024b6102cd366004611099565b610714565b61024b6102e036600461107e565b61072f565b6103086102f33660046110bb565b60009081526020819052604090206001015490565b60405190815260200161022f565b61022361032436600461107e565b610759565b61024b61033736600461107e565b6107a6565b61024b61034a3660046110d4565b6107d0565b61024b61035d3660046110d4565b6107f6565b61024b61037036600461107e565b6108ae565b6103087f8aa855a911518ecfbe5bc3088c8f3dda7badf130faaf8ace33fdc33828e1816781565b6103087f939b8dfb57ecef2aea54a93a15e86768b9d4089f1ba61c245e6ec980695f4ca481565b61024b6103d136600461107e565b6108d8565b6102236103e436600461107e565b610902565b6103087f5c91514091af31f62f596a314af7d5be40146b2f2355969392f055e12e0982fb81565b61022361041e36600461107e565b61094f565b6103087f19c860a63258efbd0ecb7d55c626237bf5c2044c26c073390b74f0c13c85743381565b61024b61045836600461107e565b61099c565b61022361046b36600461107e565b73ffffffffffffffffffffffffffffffffffffffff811660009081527fd21b659ff028ba5860060da0a2ef0b8b1b13b1f79963511fcee160c2e54d2f22602052604081205460ff1661066a565b6102236104c63660046110d4565b60009182526020828152604080842073ffffffffffffffffffffffffffffffffffffffff93909316845291905290205460ff1690565b61024b61050a36600461107e565b6109c6565b61024b61051d36600461107e565b6109f0565b61024b61053036600461107e565b610a1a565b610308600081565b61024b61054b36600461107e565b610a44565b6103087f08fb31c3e81624356c3314088aa971b73bcc82d22bc3e3b184b4593077ae327881565b6103087f12ad05bde78c5ab75238ce885307f96ecd482bb402ef831f99e7018a0f169b7b81565b61024b6105ac3660046110d4565b610a6a565b61024b6105bf36600461107e565b610a90565b6102236105d236600461107e565b610aba565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167f7965db0b00000000000000000000000000000000000000000000000000000000148061066a57507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008316145b92915050565b61069a7f08fb31c3e81624356c3314088aa971b73bcc82d22bc3e3b184b4593077ae327882610a6a565b50565b73ffffffffffffffffffffffffffffffffffffffff811660009081527fcba084d2e26105260e9ae84b007967d64af085c681345e4941eeba502738cf44602052604081205460ff1661066a565b61069a7f5c91514091af31f62f596a314af7d5be40146b2f2355969392f055e12e0982fb826107d0565b60006107208133610b07565b61072a8383610bd7565b505050565b61069a7f12ad05bde78c5ab75238ce885307f96ecd482bb402ef831f99e7018a0f169b7b826107d0565b73ffffffffffffffffffffffffffffffffffffffff811660009081527fac55d60145c2b1e72232130507b090ddd2cd26daa31eeab1e3e64b89140e668d602052604081205460ff1661066a565b61069a7f939b8dfb57ecef2aea54a93a15e86768b9d4089f1ba61c245e6ec980695f4ca482610a6a565b6000828152602081905260409020600101546107ec8133610b07565b61072a8383610c22565b73ffffffffffffffffffffffffffffffffffffffff811633146108a0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602f60248201527f416363657373436f6e74726f6c3a2063616e206f6e6c792072656e6f756e636560448201527f20726f6c657320666f722073656c66000000000000000000000000000000000060648201526084015b60405180910390fd5b6108aa8282610d12565b5050565b61069a7f8aa855a911518ecfbe5bc3088c8f3dda7badf130faaf8ace33fdc33828e1816782610a6a565b61069a7f8aa855a911518ecfbe5bc3088c8f3dda7badf130faaf8ace33fdc33828e18167826107d0565b73ffffffffffffffffffffffffffffffffffffffff811660009081527fa2630211c42039a24e17727bf18ec344681c4916090d2a50e04b9b6e50b7fea9602052604081205460ff1661066a565b73ffffffffffffffffffffffffffffffffffffffff811660009081527f9e350b38c6d0090a0631963682975411c4e88e66bd66d7f4ffcc296b4c83bf93602052604081205460ff1661066a565b61069a7f5c91514091af31f62f596a314af7d5be40146b2f2355969392f055e12e0982fb82610a6a565b61069a7f08fb31c3e81624356c3314088aa971b73bcc82d22bc3e3b184b4593077ae3278826107d0565b61069a7f19c860a63258efbd0ecb7d55c626237bf5c2044c26c073390b74f0c13c857433826107d0565b61069a7f939b8dfb57ecef2aea54a93a15e86768b9d4089f1ba61c245e6ec980695f4ca4826107d0565b61069a7f19c860a63258efbd0ecb7d55c626237bf5c2044c26c073390b74f0c13c857433825b600082815260208190526040902060010154610a868133610b07565b61072a8383610d12565b61069a7f12ad05bde78c5ab75238ce885307f96ecd482bb402ef831f99e7018a0f169b7b82610a6a565b73ffffffffffffffffffffffffffffffffffffffff811660009081527f2eadd72b6698cc7bfac8abf613f53107771ac2a3e4a3221cda0a8e2b1b91b0b4602052604081205460ff1661066a565b60008281526020818152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915290205460ff166108aa57610b5d8173ffffffffffffffffffffffffffffffffffffffff166014610dc9565b610b68836020610dc9565b604051602001610b79929190611130565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290527f08c379a0000000000000000000000000000000000000000000000000000000008252610897916004016111b1565b600082815260208190526040808220600101805490849055905190918391839186917fbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff9190a4505050565b60008281526020818152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915290205460ff166108aa5760008281526020818152604080832073ffffffffffffffffffffffffffffffffffffffff85168452909152902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055610cb43390565b73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16837f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45050565b60008281526020818152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915290205460ff16156108aa5760008281526020818152604080832073ffffffffffffffffffffffffffffffffffffffff8516808552925280832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905551339285917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a45050565b60606000610dd8836002611231565b610de390600261126e565b67ffffffffffffffff811115610dfb57610dfb611286565b6040519080825280601f01601f191660200182016040528015610e25576020820181803683370190505b5090507f300000000000000000000000000000000000000000000000000000000000000081600081518110610e5c57610e5c6112b5565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053507f780000000000000000000000000000000000000000000000000000000000000081600181518110610ebf57610ebf6112b5565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053506000610efb846002611231565b610f0690600161126e565b90505b6001811115610fa3577f303132333435363738396162636465660000000000000000000000000000000085600f1660108110610f4757610f476112b5565b1a60f81b828281518110610f5d57610f5d6112b5565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535060049490941c93610f9c816112e4565b9050610f09565b50831561100c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f537472696e67733a20686578206c656e67746820696e73756666696369656e746044820152606401610897565b9392505050565b60006020828403121561102557600080fd5b81357fffffffff000000000000000000000000000000000000000000000000000000008116811461100c57600080fd5b803573ffffffffffffffffffffffffffffffffffffffff8116811461107957600080fd5b919050565b60006020828403121561109057600080fd5b61100c82611055565b600080604083850312156110ac57600080fd5b50508035926020909101359150565b6000602082840312156110cd57600080fd5b5035919050565b600080604083850312156110e757600080fd5b823591506110f760208401611055565b90509250929050565b60005b8381101561111b578181015183820152602001611103565b8381111561112a576000848401525b50505050565b7f416363657373436f6e74726f6c3a206163636f756e7420000000000000000000815260008351611168816017850160208801611100565b7f206973206d697373696e6720726f6c652000000000000000000000000000000060179184019182015283516111a5816028840160208801611100565b01602801949350505050565b60208152600082518060208401526111d0816040850160208701611100565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169190910160400192915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048311821515161561126957611269611202565b500290565b6000821982111561128157611281611202565b500190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b6000816112f3576112f3611202565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff019056fea2646970667358221220d36d6d2e7df54059c4f97367cdc649fdae9ca664fd0c0f3b09dac9d6b21f8a7564736f6c634300080a0033",
}

// ACLManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use ACLManagerMetaData.ABI instead.
var ACLManagerABI = ACLManagerMetaData.ABI

// ACLManagerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ACLManagerMetaData.Bin instead.
var ACLManagerBin = ACLManagerMetaData.Bin

// DeployACLManager deploys a new Ethereum contract, binding an instance of ACLManager to it.
func DeployACLManager(auth *bind.TransactOpts, backend bind.ContractBackend, provider common.Address) (common.Address, *types.Transaction, *ACLManager, error) {
	parsed, err := ACLManagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ACLManagerBin), backend, provider)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ACLManager{ACLManagerCaller: ACLManagerCaller{contract: contract}, ACLManagerTransactor: ACLManagerTransactor{contract: contract}, ACLManagerFilterer: ACLManagerFilterer{contract: contract}}, nil
}

// ACLManager is an auto generated Go binding around an Ethereum contract.
type ACLManager struct {
	ACLManagerCaller     // Read-only binding to the contract
	ACLManagerTransactor // Write-only binding to the contract
	ACLManagerFilterer   // Log filterer for contract events
}

// ACLManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type ACLManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ACLManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ACLManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ACLManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ACLManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ACLManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ACLManagerSession struct {
	Contract     *ACLManager       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ACLManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ACLManagerCallerSession struct {
	Contract *ACLManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// ACLManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ACLManagerTransactorSession struct {
	Contract     *ACLManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// ACLManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type ACLManagerRaw struct {
	Contract *ACLManager // Generic contract binding to access the raw methods on
}

// ACLManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ACLManagerCallerRaw struct {
	Contract *ACLManagerCaller // Generic read-only contract binding to access the raw methods on
}

// ACLManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ACLManagerTransactorRaw struct {
	Contract *ACLManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewACLManager creates a new instance of ACLManager, bound to a specific deployed contract.
func NewACLManager(address common.Address, backend bind.ContractBackend) (*ACLManager, error) {
	contract, err := bindACLManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ACLManager{ACLManagerCaller: ACLManagerCaller{contract: contract}, ACLManagerTransactor: ACLManagerTransactor{contract: contract}, ACLManagerFilterer: ACLManagerFilterer{contract: contract}}, nil
}

// NewACLManagerCaller creates a new read-only instance of ACLManager, bound to a specific deployed contract.
func NewACLManagerCaller(address common.Address, caller bind.ContractCaller) (*ACLManagerCaller, error) {
	contract, err := bindACLManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ACLManagerCaller{contract: contract}, nil
}

// NewACLManagerTransactor creates a new write-only instance of ACLManager, bound to a specific deployed contract.
func NewACLManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*ACLManagerTransactor, error) {
	contract, err := bindACLManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ACLManagerTransactor{contract: contract}, nil
}

// NewACLManagerFilterer creates a new log filterer instance of ACLManager, bound to a specific deployed contract.
func NewACLManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*ACLManagerFilterer, error) {
	contract, err := bindACLManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ACLManagerFilterer{contract: contract}, nil
}

// bindACLManager binds a generic wrapper to an already deployed contract.
func bindACLManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ACLManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ACLManager *ACLManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ACLManager.Contract.ACLManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ACLManager *ACLManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ACLManager.Contract.ACLManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ACLManager *ACLManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ACLManager.Contract.ACLManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ACLManager *ACLManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ACLManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ACLManager *ACLManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ACLManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ACLManager *ACLManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ACLManager.Contract.contract.Transact(opts, method, params...)
}

// ADDRESSESPROVIDER is a free data retrieval call binding the contract method 0x0542975c.
//
// Solidity: function ADDRESSES_PROVIDER() view returns(address)
func (_ACLManager *ACLManagerCaller) ADDRESSESPROVIDER(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ACLManager.contract.Call(opts, &out, "ADDRESSES_PROVIDER")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ADDRESSESPROVIDER is a free data retrieval call binding the contract method 0x0542975c.
//
// Solidity: function ADDRESSES_PROVIDER() view returns(address)
func (_ACLManager *ACLManagerSession) ADDRESSESPROVIDER() (common.Address, error) {
	return _ACLManager.Contract.ADDRESSESPROVIDER(&_ACLManager.CallOpts)
}

// ADDRESSESPROVIDER is a free data retrieval call binding the contract method 0x0542975c.
//
// Solidity: function ADDRESSES_PROVIDER() view returns(address)
func (_ACLManager *ACLManagerCallerSession) ADDRESSESPROVIDER() (common.Address, error) {
	return _ACLManager.Contract.ADDRESSESPROVIDER(&_ACLManager.CallOpts)
}

// ASSETLISTINGADMINROLE is a free data retrieval call binding the contract method 0x78bb0a43.
//
// Solidity: function ASSET_LISTING_ADMIN_ROLE() view returns(bytes32)
func (_ACLManager *ACLManagerCaller) ASSETLISTINGADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ACLManager.contract.Call(opts, &out, "ASSET_LISTING_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ASSETLISTINGADMINROLE is a free data retrieval call binding the contract method 0x78bb0a43.
//
// Solidity: function ASSET_LISTING_ADMIN_ROLE() view returns(bytes32)
func (_ACLManager *ACLManagerSession) ASSETLISTINGADMINROLE() ([32]byte, error) {
	return _ACLManager.Contract.ASSETLISTINGADMINROLE(&_ACLManager.CallOpts)
}

// ASSETLISTINGADMINROLE is a free data retrieval call binding the contract method 0x78bb0a43.
//
// Solidity: function ASSET_LISTING_ADMIN_ROLE() view returns(bytes32)
func (_ACLManager *ACLManagerCallerSession) ASSETLISTINGADMINROLE() ([32]byte, error) {
	return _ACLManager.Contract.ASSETLISTINGADMINROLE(&_ACLManager.CallOpts)
}

// BRIDGEROLE is a free data retrieval call binding the contract method 0xb5bfddea.
//
// Solidity: function BRIDGE_ROLE() view returns(bytes32)
func (_ACLManager *ACLManagerCaller) BRIDGEROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ACLManager.contract.Call(opts, &out, "BRIDGE_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// BRIDGEROLE is a free data retrieval call binding the contract method 0xb5bfddea.
//
// Solidity: function BRIDGE_ROLE() view returns(bytes32)
func (_ACLManager *ACLManagerSession) BRIDGEROLE() ([32]byte, error) {
	return _ACLManager.Contract.BRIDGEROLE(&_ACLManager.CallOpts)
}

// BRIDGEROLE is a free data retrieval call binding the contract method 0xb5bfddea.
//
// Solidity: function BRIDGE_ROLE() view returns(bytes32)
func (_ACLManager *ACLManagerCallerSession) BRIDGEROLE() ([32]byte, error) {
	return _ACLManager.Contract.BRIDGEROLE(&_ACLManager.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_ACLManager *ACLManagerCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ACLManager.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_ACLManager *ACLManagerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _ACLManager.Contract.DEFAULTADMINROLE(&_ACLManager.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_ACLManager *ACLManagerCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _ACLManager.Contract.DEFAULTADMINROLE(&_ACLManager.CallOpts)
}

// EMERGENCYADMINROLE is a free data retrieval call binding the contract method 0x6e76fc8f.
//
// Solidity: function EMERGENCY_ADMIN_ROLE() view returns(bytes32)
func (_ACLManager *ACLManagerCaller) EMERGENCYADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ACLManager.contract.Call(opts, &out, "EMERGENCY_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// EMERGENCYADMINROLE is a free data retrieval call binding the contract method 0x6e76fc8f.
//
// Solidity: function EMERGENCY_ADMIN_ROLE() view returns(bytes32)
func (_ACLManager *ACLManagerSession) EMERGENCYADMINROLE() ([32]byte, error) {
	return _ACLManager.Contract.EMERGENCYADMINROLE(&_ACLManager.CallOpts)
}

// EMERGENCYADMINROLE is a free data retrieval call binding the contract method 0x6e76fc8f.
//
// Solidity: function EMERGENCY_ADMIN_ROLE() view returns(bytes32)
func (_ACLManager *ACLManagerCallerSession) EMERGENCYADMINROLE() ([32]byte, error) {
	return _ACLManager.Contract.EMERGENCYADMINROLE(&_ACLManager.CallOpts)
}

// FLASHBORROWERROLE is a free data retrieval call binding the contract method 0x5577b7a9.
//
// Solidity: function FLASH_BORROWER_ROLE() view returns(bytes32)
func (_ACLManager *ACLManagerCaller) FLASHBORROWERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ACLManager.contract.Call(opts, &out, "FLASH_BORROWER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// FLASHBORROWERROLE is a free data retrieval call binding the contract method 0x5577b7a9.
//
// Solidity: function FLASH_BORROWER_ROLE() view returns(bytes32)
func (_ACLManager *ACLManagerSession) FLASHBORROWERROLE() ([32]byte, error) {
	return _ACLManager.Contract.FLASHBORROWERROLE(&_ACLManager.CallOpts)
}

// FLASHBORROWERROLE is a free data retrieval call binding the contract method 0x5577b7a9.
//
// Solidity: function FLASH_BORROWER_ROLE() view returns(bytes32)
func (_ACLManager *ACLManagerCallerSession) FLASHBORROWERROLE() ([32]byte, error) {
	return _ACLManager.Contract.FLASHBORROWERROLE(&_ACLManager.CallOpts)
}

// POOLADMINROLE is a free data retrieval call binding the contract method 0xb8f6dba7.
//
// Solidity: function POOL_ADMIN_ROLE() view returns(bytes32)
func (_ACLManager *ACLManagerCaller) POOLADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ACLManager.contract.Call(opts, &out, "POOL_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// POOLADMINROLE is a free data retrieval call binding the contract method 0xb8f6dba7.
//
// Solidity: function POOL_ADMIN_ROLE() view returns(bytes32)
func (_ACLManager *ACLManagerSession) POOLADMINROLE() ([32]byte, error) {
	return _ACLManager.Contract.POOLADMINROLE(&_ACLManager.CallOpts)
}

// POOLADMINROLE is a free data retrieval call binding the contract method 0xb8f6dba7.
//
// Solidity: function POOL_ADMIN_ROLE() view returns(bytes32)
func (_ACLManager *ACLManagerCallerSession) POOLADMINROLE() ([32]byte, error) {
	return _ACLManager.Contract.POOLADMINROLE(&_ACLManager.CallOpts)
}

// RISKADMINROLE is a free data retrieval call binding the contract method 0x4f16b425.
//
// Solidity: function RISK_ADMIN_ROLE() view returns(bytes32)
func (_ACLManager *ACLManagerCaller) RISKADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ACLManager.contract.Call(opts, &out, "RISK_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// RISKADMINROLE is a free data retrieval call binding the contract method 0x4f16b425.
//
// Solidity: function RISK_ADMIN_ROLE() view returns(bytes32)
func (_ACLManager *ACLManagerSession) RISKADMINROLE() ([32]byte, error) {
	return _ACLManager.Contract.RISKADMINROLE(&_ACLManager.CallOpts)
}

// RISKADMINROLE is a free data retrieval call binding the contract method 0x4f16b425.
//
// Solidity: function RISK_ADMIN_ROLE() view returns(bytes32)
func (_ACLManager *ACLManagerCallerSession) RISKADMINROLE() ([32]byte, error) {
	return _ACLManager.Contract.RISKADMINROLE(&_ACLManager.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_ACLManager *ACLManagerCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _ACLManager.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_ACLManager *ACLManagerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _ACLManager.Contract.GetRoleAdmin(&_ACLManager.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_ACLManager *ACLManagerCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _ACLManager.Contract.GetRoleAdmin(&_ACLManager.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_ACLManager *ACLManagerCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _ACLManager.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_ACLManager *ACLManagerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _ACLManager.Contract.HasRole(&_ACLManager.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_ACLManager *ACLManagerCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _ACLManager.Contract.HasRole(&_ACLManager.CallOpts, role, account)
}

// IsAssetListingAdmin is a free data retrieval call binding the contract method 0x13ee32e0.
//
// Solidity: function isAssetListingAdmin(address admin) view returns(bool)
func (_ACLManager *ACLManagerCaller) IsAssetListingAdmin(opts *bind.CallOpts, admin common.Address) (bool, error) {
	var out []interface{}
	err := _ACLManager.contract.Call(opts, &out, "isAssetListingAdmin", admin)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsAssetListingAdmin is a free data retrieval call binding the contract method 0x13ee32e0.
//
// Solidity: function isAssetListingAdmin(address admin) view returns(bool)
func (_ACLManager *ACLManagerSession) IsAssetListingAdmin(admin common.Address) (bool, error) {
	return _ACLManager.Contract.IsAssetListingAdmin(&_ACLManager.CallOpts, admin)
}

// IsAssetListingAdmin is a free data retrieval call binding the contract method 0x13ee32e0.
//
// Solidity: function isAssetListingAdmin(address admin) view returns(bool)
func (_ACLManager *ACLManagerCallerSession) IsAssetListingAdmin(admin common.Address) (bool, error) {
	return _ACLManager.Contract.IsAssetListingAdmin(&_ACLManager.CallOpts, admin)
}

// IsBridge is a free data retrieval call binding the contract method 0x726600ce.
//
// Solidity: function isBridge(address bridge) view returns(bool)
func (_ACLManager *ACLManagerCaller) IsBridge(opts *bind.CallOpts, bridge common.Address) (bool, error) {
	var out []interface{}
	err := _ACLManager.contract.Call(opts, &out, "isBridge", bridge)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsBridge is a free data retrieval call binding the contract method 0x726600ce.
//
// Solidity: function isBridge(address bridge) view returns(bool)
func (_ACLManager *ACLManagerSession) IsBridge(bridge common.Address) (bool, error) {
	return _ACLManager.Contract.IsBridge(&_ACLManager.CallOpts, bridge)
}

// IsBridge is a free data retrieval call binding the contract method 0x726600ce.
//
// Solidity: function isBridge(address bridge) view returns(bool)
func (_ACLManager *ACLManagerCallerSession) IsBridge(bridge common.Address) (bool, error) {
	return _ACLManager.Contract.IsBridge(&_ACLManager.CallOpts, bridge)
}

// IsEmergencyAdmin is a free data retrieval call binding the contract method 0x2500f2b6.
//
// Solidity: function isEmergencyAdmin(address admin) view returns(bool)
func (_ACLManager *ACLManagerCaller) IsEmergencyAdmin(opts *bind.CallOpts, admin common.Address) (bool, error) {
	var out []interface{}
	err := _ACLManager.contract.Call(opts, &out, "isEmergencyAdmin", admin)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsEmergencyAdmin is a free data retrieval call binding the contract method 0x2500f2b6.
//
// Solidity: function isEmergencyAdmin(address admin) view returns(bool)
func (_ACLManager *ACLManagerSession) IsEmergencyAdmin(admin common.Address) (bool, error) {
	return _ACLManager.Contract.IsEmergencyAdmin(&_ACLManager.CallOpts, admin)
}

// IsEmergencyAdmin is a free data retrieval call binding the contract method 0x2500f2b6.
//
// Solidity: function isEmergencyAdmin(address admin) view returns(bool)
func (_ACLManager *ACLManagerCallerSession) IsEmergencyAdmin(admin common.Address) (bool, error) {
	return _ACLManager.Contract.IsEmergencyAdmin(&_ACLManager.CallOpts, admin)
}

// IsFlashBorrower is a free data retrieval call binding the contract method 0xfa50f297.
//
// Solidity: function isFlashBorrower(address borrower) view returns(bool)
func (_ACLManager *ACLManagerCaller) IsFlashBorrower(opts *bind.CallOpts, borrower common.Address) (bool, error) {
	var out []interface{}
	err := _ACLManager.contract.Call(opts, &out, "isFlashBorrower", borrower)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsFlashBorrower is a free data retrieval call binding the contract method 0xfa50f297.
//
// Solidity: function isFlashBorrower(address borrower) view returns(bool)
func (_ACLManager *ACLManagerSession) IsFlashBorrower(borrower common.Address) (bool, error) {
	return _ACLManager.Contract.IsFlashBorrower(&_ACLManager.CallOpts, borrower)
}

// IsFlashBorrower is a free data retrieval call binding the contract method 0xfa50f297.
//
// Solidity: function isFlashBorrower(address borrower) view returns(bool)
func (_ACLManager *ACLManagerCallerSession) IsFlashBorrower(borrower common.Address) (bool, error) {
	return _ACLManager.Contract.IsFlashBorrower(&_ACLManager.CallOpts, borrower)
}

// IsPoolAdmin is a free data retrieval call binding the contract method 0x7be53ca1.
//
// Solidity: function isPoolAdmin(address admin) view returns(bool)
func (_ACLManager *ACLManagerCaller) IsPoolAdmin(opts *bind.CallOpts, admin common.Address) (bool, error) {
	var out []interface{}
	err := _ACLManager.contract.Call(opts, &out, "isPoolAdmin", admin)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsPoolAdmin is a free data retrieval call binding the contract method 0x7be53ca1.
//
// Solidity: function isPoolAdmin(address admin) view returns(bool)
func (_ACLManager *ACLManagerSession) IsPoolAdmin(admin common.Address) (bool, error) {
	return _ACLManager.Contract.IsPoolAdmin(&_ACLManager.CallOpts, admin)
}

// IsPoolAdmin is a free data retrieval call binding the contract method 0x7be53ca1.
//
// Solidity: function isPoolAdmin(address admin) view returns(bool)
func (_ACLManager *ACLManagerCallerSession) IsPoolAdmin(admin common.Address) (bool, error) {
	return _ACLManager.Contract.IsPoolAdmin(&_ACLManager.CallOpts, admin)
}

// IsRiskAdmin is a free data retrieval call binding the contract method 0x674b5e4d.
//
// Solidity: function isRiskAdmin(address admin) view returns(bool)
func (_ACLManager *ACLManagerCaller) IsRiskAdmin(opts *bind.CallOpts, admin common.Address) (bool, error) {
	var out []interface{}
	err := _ACLManager.contract.Call(opts, &out, "isRiskAdmin", admin)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsRiskAdmin is a free data retrieval call binding the contract method 0x674b5e4d.
//
// Solidity: function isRiskAdmin(address admin) view returns(bool)
func (_ACLManager *ACLManagerSession) IsRiskAdmin(admin common.Address) (bool, error) {
	return _ACLManager.Contract.IsRiskAdmin(&_ACLManager.CallOpts, admin)
}

// IsRiskAdmin is a free data retrieval call binding the contract method 0x674b5e4d.
//
// Solidity: function isRiskAdmin(address admin) view returns(bool)
func (_ACLManager *ACLManagerCallerSession) IsRiskAdmin(admin common.Address) (bool, error) {
	return _ACLManager.Contract.IsRiskAdmin(&_ACLManager.CallOpts, admin)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ACLManager *ACLManagerCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _ACLManager.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ACLManager *ACLManagerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ACLManager.Contract.SupportsInterface(&_ACLManager.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ACLManager *ACLManagerCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ACLManager.Contract.SupportsInterface(&_ACLManager.CallOpts, interfaceId)
}

// AddAssetListingAdmin is a paid mutator transaction binding the contract method 0x9a2b96f7.
//
// Solidity: function addAssetListingAdmin(address admin) returns()
func (_ACLManager *ACLManagerTransactor) AddAssetListingAdmin(opts *bind.TransactOpts, admin common.Address) (*types.Transaction, error) {
	return _ACLManager.contract.Transact(opts, "addAssetListingAdmin", admin)
}

// AddAssetListingAdmin is a paid mutator transaction binding the contract method 0x9a2b96f7.
//
// Solidity: function addAssetListingAdmin(address admin) returns()
func (_ACLManager *ACLManagerSession) AddAssetListingAdmin(admin common.Address) (*types.Transaction, error) {
	return _ACLManager.Contract.AddAssetListingAdmin(&_ACLManager.TransactOpts, admin)
}

// AddAssetListingAdmin is a paid mutator transaction binding the contract method 0x9a2b96f7.
//
// Solidity: function addAssetListingAdmin(address admin) returns()
func (_ACLManager *ACLManagerTransactorSession) AddAssetListingAdmin(admin common.Address) (*types.Transaction, error) {
	return _ACLManager.Contract.AddAssetListingAdmin(&_ACLManager.TransactOpts, admin)
}

// AddBridge is a paid mutator transaction binding the contract method 0x9712fdf8.
//
// Solidity: function addBridge(address bridge) returns()
func (_ACLManager *ACLManagerTransactor) AddBridge(opts *bind.TransactOpts, bridge common.Address) (*types.Transaction, error) {
	return _ACLManager.contract.Transact(opts, "addBridge", bridge)
}

// AddBridge is a paid mutator transaction binding the contract method 0x9712fdf8.
//
// Solidity: function addBridge(address bridge) returns()
func (_ACLManager *ACLManagerSession) AddBridge(bridge common.Address) (*types.Transaction, error) {
	return _ACLManager.Contract.AddBridge(&_ACLManager.TransactOpts, bridge)
}

// AddBridge is a paid mutator transaction binding the contract method 0x9712fdf8.
//
// Solidity: function addBridge(address bridge) returns()
func (_ACLManager *ACLManagerTransactorSession) AddBridge(bridge common.Address) (*types.Transaction, error) {
	return _ACLManager.Contract.AddBridge(&_ACLManager.TransactOpts, bridge)
}

// AddEmergencyAdmin is a paid mutator transaction binding the contract method 0x179efb09.
//
// Solidity: function addEmergencyAdmin(address admin) returns()
func (_ACLManager *ACLManagerTransactor) AddEmergencyAdmin(opts *bind.TransactOpts, admin common.Address) (*types.Transaction, error) {
	return _ACLManager.contract.Transact(opts, "addEmergencyAdmin", admin)
}

// AddEmergencyAdmin is a paid mutator transaction binding the contract method 0x179efb09.
//
// Solidity: function addEmergencyAdmin(address admin) returns()
func (_ACLManager *ACLManagerSession) AddEmergencyAdmin(admin common.Address) (*types.Transaction, error) {
	return _ACLManager.Contract.AddEmergencyAdmin(&_ACLManager.TransactOpts, admin)
}

// AddEmergencyAdmin is a paid mutator transaction binding the contract method 0x179efb09.
//
// Solidity: function addEmergencyAdmin(address admin) returns()
func (_ACLManager *ACLManagerTransactorSession) AddEmergencyAdmin(admin common.Address) (*types.Transaction, error) {
	return _ACLManager.Contract.AddEmergencyAdmin(&_ACLManager.TransactOpts, admin)
}

// AddFlashBorrower is a paid mutator transaction binding the contract method 0x9ac9d80b.
//
// Solidity: function addFlashBorrower(address borrower) returns()
func (_ACLManager *ACLManagerTransactor) AddFlashBorrower(opts *bind.TransactOpts, borrower common.Address) (*types.Transaction, error) {
	return _ACLManager.contract.Transact(opts, "addFlashBorrower", borrower)
}

// AddFlashBorrower is a paid mutator transaction binding the contract method 0x9ac9d80b.
//
// Solidity: function addFlashBorrower(address borrower) returns()
func (_ACLManager *ACLManagerSession) AddFlashBorrower(borrower common.Address) (*types.Transaction, error) {
	return _ACLManager.Contract.AddFlashBorrower(&_ACLManager.TransactOpts, borrower)
}

// AddFlashBorrower is a paid mutator transaction binding the contract method 0x9ac9d80b.
//
// Solidity: function addFlashBorrower(address borrower) returns()
func (_ACLManager *ACLManagerTransactorSession) AddFlashBorrower(borrower common.Address) (*types.Transaction, error) {
	return _ACLManager.Contract.AddFlashBorrower(&_ACLManager.TransactOpts, borrower)
}

// AddPoolAdmin is a paid mutator transaction binding the contract method 0x22650caf.
//
// Solidity: function addPoolAdmin(address admin) returns()
func (_ACLManager *ACLManagerTransactor) AddPoolAdmin(opts *bind.TransactOpts, admin common.Address) (*types.Transaction, error) {
	return _ACLManager.contract.Transact(opts, "addPoolAdmin", admin)
}

// AddPoolAdmin is a paid mutator transaction binding the contract method 0x22650caf.
//
// Solidity: function addPoolAdmin(address admin) returns()
func (_ACLManager *ACLManagerSession) AddPoolAdmin(admin common.Address) (*types.Transaction, error) {
	return _ACLManager.Contract.AddPoolAdmin(&_ACLManager.TransactOpts, admin)
}

// AddPoolAdmin is a paid mutator transaction binding the contract method 0x22650caf.
//
// Solidity: function addPoolAdmin(address admin) returns()
func (_ACLManager *ACLManagerTransactorSession) AddPoolAdmin(admin common.Address) (*types.Transaction, error) {
	return _ACLManager.Contract.AddPoolAdmin(&_ACLManager.TransactOpts, admin)
}

// AddRiskAdmin is a paid mutator transaction binding the contract method 0x5b9a94e4.
//
// Solidity: function addRiskAdmin(address admin) returns()
func (_ACLManager *ACLManagerTransactor) AddRiskAdmin(opts *bind.TransactOpts, admin common.Address) (*types.Transaction, error) {
	return _ACLManager.contract.Transact(opts, "addRiskAdmin", admin)
}

// AddRiskAdmin is a paid mutator transaction binding the contract method 0x5b9a94e4.
//
// Solidity: function addRiskAdmin(address admin) returns()
func (_ACLManager *ACLManagerSession) AddRiskAdmin(admin common.Address) (*types.Transaction, error) {
	return _ACLManager.Contract.AddRiskAdmin(&_ACLManager.TransactOpts, admin)
}

// AddRiskAdmin is a paid mutator transaction binding the contract method 0x5b9a94e4.
//
// Solidity: function addRiskAdmin(address admin) returns()
func (_ACLManager *ACLManagerTransactorSession) AddRiskAdmin(admin common.Address) (*types.Transaction, error) {
	return _ACLManager.Contract.AddRiskAdmin(&_ACLManager.TransactOpts, admin)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_ACLManager *ACLManagerTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ACLManager.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_ACLManager *ACLManagerSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ACLManager.Contract.GrantRole(&_ACLManager.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_ACLManager *ACLManagerTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ACLManager.Contract.GrantRole(&_ACLManager.TransactOpts, role, account)
}

// RemoveAssetListingAdmin is a paid mutator transaction binding the contract method 0xa21bce15.
//
// Solidity: function removeAssetListingAdmin(address admin) returns()
func (_ACLManager *ACLManagerTransactor) RemoveAssetListingAdmin(opts *bind.TransactOpts, admin common.Address) (*types.Transaction, error) {
	return _ACLManager.contract.Transact(opts, "removeAssetListingAdmin", admin)
}

// RemoveAssetListingAdmin is a paid mutator transaction binding the contract method 0xa21bce15.
//
// Solidity: function removeAssetListingAdmin(address admin) returns()
func (_ACLManager *ACLManagerSession) RemoveAssetListingAdmin(admin common.Address) (*types.Transaction, error) {
	return _ACLManager.Contract.RemoveAssetListingAdmin(&_ACLManager.TransactOpts, admin)
}

// RemoveAssetListingAdmin is a paid mutator transaction binding the contract method 0xa21bce15.
//
// Solidity: function removeAssetListingAdmin(address admin) returns()
func (_ACLManager *ACLManagerTransactorSession) RemoveAssetListingAdmin(admin common.Address) (*types.Transaction, error) {
	return _ACLManager.Contract.RemoveAssetListingAdmin(&_ACLManager.TransactOpts, admin)
}

// RemoveBridge is a paid mutator transaction binding the contract method 0x04df017d.
//
// Solidity: function removeBridge(address bridge) returns()
func (_ACLManager *ACLManagerTransactor) RemoveBridge(opts *bind.TransactOpts, bridge common.Address) (*types.Transaction, error) {
	return _ACLManager.contract.Transact(opts, "removeBridge", bridge)
}

// RemoveBridge is a paid mutator transaction binding the contract method 0x04df017d.
//
// Solidity: function removeBridge(address bridge) returns()
func (_ACLManager *ACLManagerSession) RemoveBridge(bridge common.Address) (*types.Transaction, error) {
	return _ACLManager.Contract.RemoveBridge(&_ACLManager.TransactOpts, bridge)
}

// RemoveBridge is a paid mutator transaction binding the contract method 0x04df017d.
//
// Solidity: function removeBridge(address bridge) returns()
func (_ACLManager *ACLManagerTransactorSession) RemoveBridge(bridge common.Address) (*types.Transaction, error) {
	return _ACLManager.Contract.RemoveBridge(&_ACLManager.TransactOpts, bridge)
}

// RemoveEmergencyAdmin is a paid mutator transaction binding the contract method 0x7a9a93f4.
//
// Solidity: function removeEmergencyAdmin(address admin) returns()
func (_ACLManager *ACLManagerTransactor) RemoveEmergencyAdmin(opts *bind.TransactOpts, admin common.Address) (*types.Transaction, error) {
	return _ACLManager.contract.Transact(opts, "removeEmergencyAdmin", admin)
}

// RemoveEmergencyAdmin is a paid mutator transaction binding the contract method 0x7a9a93f4.
//
// Solidity: function removeEmergencyAdmin(address admin) returns()
func (_ACLManager *ACLManagerSession) RemoveEmergencyAdmin(admin common.Address) (*types.Transaction, error) {
	return _ACLManager.Contract.RemoveEmergencyAdmin(&_ACLManager.TransactOpts, admin)
}

// RemoveEmergencyAdmin is a paid mutator transaction binding the contract method 0x7a9a93f4.
//
// Solidity: function removeEmergencyAdmin(address admin) returns()
func (_ACLManager *ACLManagerTransactorSession) RemoveEmergencyAdmin(admin common.Address) (*types.Transaction, error) {
	return _ACLManager.Contract.RemoveEmergencyAdmin(&_ACLManager.TransactOpts, admin)
}

// RemoveFlashBorrower is a paid mutator transaction binding the contract method 0x253cf980.
//
// Solidity: function removeFlashBorrower(address borrower) returns()
func (_ACLManager *ACLManagerTransactor) RemoveFlashBorrower(opts *bind.TransactOpts, borrower common.Address) (*types.Transaction, error) {
	return _ACLManager.contract.Transact(opts, "removeFlashBorrower", borrower)
}

// RemoveFlashBorrower is a paid mutator transaction binding the contract method 0x253cf980.
//
// Solidity: function removeFlashBorrower(address borrower) returns()
func (_ACLManager *ACLManagerSession) RemoveFlashBorrower(borrower common.Address) (*types.Transaction, error) {
	return _ACLManager.Contract.RemoveFlashBorrower(&_ACLManager.TransactOpts, borrower)
}

// RemoveFlashBorrower is a paid mutator transaction binding the contract method 0x253cf980.
//
// Solidity: function removeFlashBorrower(address borrower) returns()
func (_ACLManager *ACLManagerTransactorSession) RemoveFlashBorrower(borrower common.Address) (*types.Transaction, error) {
	return _ACLManager.Contract.RemoveFlashBorrower(&_ACLManager.TransactOpts, borrower)
}

// RemovePoolAdmin is a paid mutator transaction binding the contract method 0xf83695cb.
//
// Solidity: function removePoolAdmin(address admin) returns()
func (_ACLManager *ACLManagerTransactor) RemovePoolAdmin(opts *bind.TransactOpts, admin common.Address) (*types.Transaction, error) {
	return _ACLManager.contract.Transact(opts, "removePoolAdmin", admin)
}

// RemovePoolAdmin is a paid mutator transaction binding the contract method 0xf83695cb.
//
// Solidity: function removePoolAdmin(address admin) returns()
func (_ACLManager *ACLManagerSession) RemovePoolAdmin(admin common.Address) (*types.Transaction, error) {
	return _ACLManager.Contract.RemovePoolAdmin(&_ACLManager.TransactOpts, admin)
}

// RemovePoolAdmin is a paid mutator transaction binding the contract method 0xf83695cb.
//
// Solidity: function removePoolAdmin(address admin) returns()
func (_ACLManager *ACLManagerTransactorSession) RemovePoolAdmin(admin common.Address) (*types.Transaction, error) {
	return _ACLManager.Contract.RemovePoolAdmin(&_ACLManager.TransactOpts, admin)
}

// RemoveRiskAdmin is a paid mutator transaction binding the contract method 0x3c5a08e5.
//
// Solidity: function removeRiskAdmin(address admin) returns()
func (_ACLManager *ACLManagerTransactor) RemoveRiskAdmin(opts *bind.TransactOpts, admin common.Address) (*types.Transaction, error) {
	return _ACLManager.contract.Transact(opts, "removeRiskAdmin", admin)
}

// RemoveRiskAdmin is a paid mutator transaction binding the contract method 0x3c5a08e5.
//
// Solidity: function removeRiskAdmin(address admin) returns()
func (_ACLManager *ACLManagerSession) RemoveRiskAdmin(admin common.Address) (*types.Transaction, error) {
	return _ACLManager.Contract.RemoveRiskAdmin(&_ACLManager.TransactOpts, admin)
}

// RemoveRiskAdmin is a paid mutator transaction binding the contract method 0x3c5a08e5.
//
// Solidity: function removeRiskAdmin(address admin) returns()
func (_ACLManager *ACLManagerTransactorSession) RemoveRiskAdmin(admin common.Address) (*types.Transaction, error) {
	return _ACLManager.Contract.RemoveRiskAdmin(&_ACLManager.TransactOpts, admin)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_ACLManager *ACLManagerTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ACLManager.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_ACLManager *ACLManagerSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ACLManager.Contract.RenounceRole(&_ACLManager.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_ACLManager *ACLManagerTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ACLManager.Contract.RenounceRole(&_ACLManager.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_ACLManager *ACLManagerTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ACLManager.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_ACLManager *ACLManagerSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ACLManager.Contract.RevokeRole(&_ACLManager.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_ACLManager *ACLManagerTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ACLManager.Contract.RevokeRole(&_ACLManager.TransactOpts, role, account)
}

// SetRoleAdmin is a paid mutator transaction binding the contract method 0x1e4e0091.
//
// Solidity: function setRoleAdmin(bytes32 role, bytes32 adminRole) returns()
func (_ACLManager *ACLManagerTransactor) SetRoleAdmin(opts *bind.TransactOpts, role [32]byte, adminRole [32]byte) (*types.Transaction, error) {
	return _ACLManager.contract.Transact(opts, "setRoleAdmin", role, adminRole)
}

// SetRoleAdmin is a paid mutator transaction binding the contract method 0x1e4e0091.
//
// Solidity: function setRoleAdmin(bytes32 role, bytes32 adminRole) returns()
func (_ACLManager *ACLManagerSession) SetRoleAdmin(role [32]byte, adminRole [32]byte) (*types.Transaction, error) {
	return _ACLManager.Contract.SetRoleAdmin(&_ACLManager.TransactOpts, role, adminRole)
}

// SetRoleAdmin is a paid mutator transaction binding the contract method 0x1e4e0091.
//
// Solidity: function setRoleAdmin(bytes32 role, bytes32 adminRole) returns()
func (_ACLManager *ACLManagerTransactorSession) SetRoleAdmin(role [32]byte, adminRole [32]byte) (*types.Transaction, error) {
	return _ACLManager.Contract.SetRoleAdmin(&_ACLManager.TransactOpts, role, adminRole)
}

// ACLManagerRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the ACLManager contract.
type ACLManagerRoleAdminChangedIterator struct {
	Event *ACLManagerRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *ACLManagerRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ACLManagerRoleAdminChanged)
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
		it.Event = new(ACLManagerRoleAdminChanged)
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
func (it *ACLManagerRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ACLManagerRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ACLManagerRoleAdminChanged represents a RoleAdminChanged event raised by the ACLManager contract.
type ACLManagerRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_ACLManager *ACLManagerFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*ACLManagerRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _ACLManager.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &ACLManagerRoleAdminChangedIterator{contract: _ACLManager.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_ACLManager *ACLManagerFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *ACLManagerRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _ACLManager.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ACLManagerRoleAdminChanged)
				if err := _ACLManager.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_ACLManager *ACLManagerFilterer) ParseRoleAdminChanged(log types.Log) (*ACLManagerRoleAdminChanged, error) {
	event := new(ACLManagerRoleAdminChanged)
	if err := _ACLManager.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ACLManagerRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the ACLManager contract.
type ACLManagerRoleGrantedIterator struct {
	Event *ACLManagerRoleGranted // Event containing the contract specifics and raw log

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
func (it *ACLManagerRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ACLManagerRoleGranted)
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
		it.Event = new(ACLManagerRoleGranted)
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
func (it *ACLManagerRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ACLManagerRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ACLManagerRoleGranted represents a RoleGranted event raised by the ACLManager contract.
type ACLManagerRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_ACLManager *ACLManagerFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*ACLManagerRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _ACLManager.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &ACLManagerRoleGrantedIterator{contract: _ACLManager.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_ACLManager *ACLManagerFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *ACLManagerRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _ACLManager.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ACLManagerRoleGranted)
				if err := _ACLManager.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_ACLManager *ACLManagerFilterer) ParseRoleGranted(log types.Log) (*ACLManagerRoleGranted, error) {
	event := new(ACLManagerRoleGranted)
	if err := _ACLManager.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ACLManagerRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the ACLManager contract.
type ACLManagerRoleRevokedIterator struct {
	Event *ACLManagerRoleRevoked // Event containing the contract specifics and raw log

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
func (it *ACLManagerRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ACLManagerRoleRevoked)
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
		it.Event = new(ACLManagerRoleRevoked)
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
func (it *ACLManagerRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ACLManagerRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ACLManagerRoleRevoked represents a RoleRevoked event raised by the ACLManager contract.
type ACLManagerRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_ACLManager *ACLManagerFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*ACLManagerRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _ACLManager.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &ACLManagerRoleRevokedIterator{contract: _ACLManager.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_ACLManager *ACLManagerFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *ACLManagerRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _ACLManager.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ACLManagerRoleRevoked)
				if err := _ACLManager.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_ACLManager *ACLManagerFilterer) ParseRoleRevoked(log types.Log) (*ACLManagerRoleRevoked, error) {
	event := new(ACLManagerRoleRevoked)
	if err := _ACLManager.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
