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

// PoolAddressesProviderMetaData contains all meta data concerning the PoolAddressesProvider contract.
var PoolAddressesProviderMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"marketId\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newAddress\",\"type\":\"address\"}],\"name\":\"ACLAdminUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newAddress\",\"type\":\"address\"}],\"name\":\"ACLManagerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newAddress\",\"type\":\"address\"}],\"name\":\"AddressSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proxyAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldImplementationAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newImplementationAddress\",\"type\":\"address\"}],\"name\":\"AddressSetAsProxy\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"oldMarketId\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"string\",\"name\":\"newMarketId\",\"type\":\"string\"}],\"name\":\"MarketIdSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newAddress\",\"type\":\"address\"}],\"name\":\"PoolConfiguratorUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newAddress\",\"type\":\"address\"}],\"name\":\"PoolDataProviderUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newAddress\",\"type\":\"address\"}],\"name\":\"PoolUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newAddress\",\"type\":\"address\"}],\"name\":\"PriceOracleSentinelUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newAddress\",\"type\":\"address\"}],\"name\":\"PriceOracleUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proxyAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementationAddress\",\"type\":\"address\"}],\"name\":\"ProxyCreated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"getACLAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getACLManager\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"getAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMarketId\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPool\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPoolConfigurator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPoolDataProvider\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPriceOracle\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPriceOracleSentinel\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAclAdmin\",\"type\":\"address\"}],\"name\":\"setACLAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAclManager\",\"type\":\"address\"}],\"name\":\"setACLManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"newAddress\",\"type\":\"address\"}],\"name\":\"setAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"newImplementationAddress\",\"type\":\"address\"}],\"name\":\"setAddressAsProxy\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"newMarketId\",\"type\":\"string\"}],\"name\":\"setMarketId\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newPoolConfiguratorImpl\",\"type\":\"address\"}],\"name\":\"setPoolConfiguratorImpl\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newDataProvider\",\"type\":\"address\"}],\"name\":\"setPoolDataProvider\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newPoolImpl\",\"type\":\"address\"}],\"name\":\"setPoolImpl\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newPriceOracle\",\"type\":\"address\"}],\"name\":\"setPriceOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newPriceOracleSentinel\",\"type\":\"address\"}],\"name\":\"setPriceOracleSentinel\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162002b3538038062002b358339810160408190526200003491620003aa565b600080546001600160a01b0319163390811782556040519091829160008051602062002b15833981519152908290a3506200006f8262000082565b6200007a816200018d565b5050620004d2565b600060018054620000939062000477565b80601f0160208091040260200160405190810160405280929190818152602001828054620000c19062000477565b8015620001125780601f10620000e65761010080835404028352916020019162000112565b820191906000526020600020905b815481529060010190602001808311620000f457829003601f168201915b5050855193945062000130936001935060208701925090506200029e565b5081604051620001419190620004b4565b604051809103902081604051620001599190620004b4565b604051908190038120907fe685c8cdecc6030c45030fd54778812cb84ed8e4467c38294403d68ba786082390600090a35050565b6000546001600160a01b03163314620001ed5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064015b60405180910390fd5b6001600160a01b038116620002545760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b6064820152608401620001e4565b600080546040516001600160a01b038085169392169160008051602062002b1583398151915291a3600080546001600160a01b0319166001600160a01b0392909216919091179055565b828054620002ac9062000477565b90600052602060002090601f016020900481019282620002d057600085556200031b565b82601f10620002eb57805160ff19168380011785556200031b565b828001600101855582156200031b579182015b828111156200031b578251825591602001919060010190620002fe565b50620003299291506200032d565b5090565b5b808211156200032957600081556001016200032e565b634e487b7160e01b600052604160045260246000fd5b60005b83811015620003775781810151838201526020016200035d565b8381111562000387576000848401525b50505050565b80516001600160a01b0381168114620003a557600080fd5b919050565b60008060408385031215620003be57600080fd5b82516001600160401b0380821115620003d657600080fd5b818501915085601f830112620003eb57600080fd5b81518181111562000400576200040062000344565b604051601f8201601f19908116603f011681019083821181831017156200042b576200042b62000344565b816040528281528860208487010111156200044557600080fd5b620004588360208301602088016200035a565b80965050505050506200046e602084016200038d565b90509250929050565b600181811c908216806200048c57607f821691505b60208210811415620004ae57634e487b7160e01b600052602260045260246000fd5b50919050565b60008251620004c88184602087016200035a565b9190910192915050565b61263380620004e26000396000f3fe608060405234801561001057600080fd5b50600436106101825760003560e01c806376d84ffc116100d8578063e4ca28b71161008c578063f2fde38b11610066578063f2fde38b1461052f578063f67b184714610542578063fca513a81461055557600080fd5b8063e4ca28b7146104a3578063e860accb146104b6578063ed301ca91461051c57600080fd5b8063a1564406116100bd578063a15644061461046a578063ca446dd91461047d578063e44e9ed11461049057600080fd5b806376d84ffc146104395780638da5cb5b1461044c57600080fd5b80635dcc528c1161013a578063707cd71611610114578063707cd716146103b8578063715018a61461041e57806374944cec1461042657600080fd5b80635dcc528c146102d95780635eb88d3d146102ec578063631adfca1461035257600080fd5b806321f8a7211161016b57806321f8a72114610279578063530e784f146102af578063568ef470146102c457600080fd5b8063026b1d5f146101875780630e67178c14610213575b600080fd5b7f504f4f4c0000000000000000000000000000000000000000000000000000000060005260026020527f4fe005067814bb4b024d9515847377d15011b64593c006223b4a722952d2c05a5473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b7f41434c5f41444d494e000000000000000000000000000000000000000000000060005260026020527ffab167ad2009dcb80ee379700bb4bd029d97c1181ed9d961625632c8a6f051c65473ffffffffffffffffffffffffffffffffffffffff166101e9565b6101e9610287366004611962565b60009081526002602052604090205473ffffffffffffffffffffffffffffffffffffffff1690565b6102c26102bd36600461199d565b6105bb565b005b6102cc6106ff565b60405161020a9190611a3b565b6102c26102e7366004611a4e565b610791565b7f50524943455f4f5241434c455f53454e54494e454c000000000000000000000060005260026020527f0d2c1bcee56447b4f46248272f34207a580a5c40f666a31f4e2fbb470ea53ab85473ffffffffffffffffffffffffffffffffffffffff166101e9565b7f504f4f4c5f434f4e464947555241544f5200000000000000000000000000000060005260026020527f90c127ef1c12c03f5781afeca3079527ea5333738078bba6fea26825bf9bf2c55473ffffffffffffffffffffffffffffffffffffffff166101e9565b7f41434c5f4d414e4147455200000000000000000000000000000000000000000060005260026020527f9edef266ef35fd0c6e131df0f31a330f3dd4c4d19dd31ed615c21d005c68116b5473ffffffffffffffffffffffffffffffffffffffff166101e9565b6102c26108a7565b6102c261043436600461199d565b610997565b6102c261044736600461199d565b610ad6565b60005473ffffffffffffffffffffffffffffffffffffffff166101e9565b6102c261047836600461199d565b610c15565b6102c261048b366004611a4e565b610d4b565b6102c261049e36600461199d565b610e4f565b6102c26104b136600461199d565b610f8e565b7f444154415f50524f56494445520000000000000000000000000000000000000060005260026020527fcd7944601aaa5cd7ccdae1bebec659e98c6aac8f12486b30e59db0d39698051f5473ffffffffffffffffffffffffffffffffffffffff166101e9565b6102c261052a36600461199d565b6110c4565b6102c261053d36600461199d565b611203565b6102c2610550366004611aad565b6113b4565b7f50524943455f4f5241434c45000000000000000000000000000000000000000060005260026020527f740f710666bd7a12af42df98311e541e47f7fd33d382d11602457a6d540cbd635473ffffffffffffffffffffffffffffffffffffffff166101e9565b60005473ffffffffffffffffffffffffffffffffffffffff163314610641576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064015b60405180910390fd5b7f50524943455f4f5241434c450000000000000000000000000000000000000000600090815260026020527f740f710666bd7a12af42df98311e541e47f7fd33d382d11602457a6d540cbd63805473ffffffffffffffffffffffffffffffffffffffff8481167fffffffffffffffffffffffff00000000000000000000000000000000000000008316811790935560405191169283917f56b5f80d8cac1479698aa7d01605fd6111e90b15fc4d2b377417f46034876cbd9190a35050565b60606001805461070e90611b7c565b80601f016020809104026020016040519081016040528092919081815260200182805461073a90611b7c565b80156107875780601f1061075c57610100808354040283529160200191610787565b820191906000526020600020905b81548152906001019060200180831161076a57829003601f168201915b5050505050905090565b60005473ffffffffffffffffffffffffffffffffffffffff163314610812576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610638565b60008281526002602052604081205473ffffffffffffffffffffffffffffffffffffffff169061084184611441565b905061084d84846114f8565b60405173ffffffffffffffffffffffffffffffffffffffff8281168252808516919084169086907f3bbd45b5429b385e3fb37ad5cd1cd1435a3c8ec32196c7937597365a3fd3e99c9060200160405180910390a450505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610928576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610638565b6000805460405173ffffffffffffffffffffffffffffffffffffffff909116907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a3600080547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055565b60005473ffffffffffffffffffffffffffffffffffffffff163314610a18576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610638565b7f50524943455f4f5241434c455f53454e54494e454c0000000000000000000000600090815260026020527f0d2c1bcee56447b4f46248272f34207a580a5c40f666a31f4e2fbb470ea53ab8805473ffffffffffffffffffffffffffffffffffffffff8481167fffffffffffffffffffffffff00000000000000000000000000000000000000008316811790935560405191169283917f5326514eeca90494a14bedabcff812a0e683029ee85d1e23824d44fd14cd6ae79190a35050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610b57576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610638565b7f41434c5f41444d494e0000000000000000000000000000000000000000000000600090815260026020527ffab167ad2009dcb80ee379700bb4bd029d97c1181ed9d961625632c8a6f051c6805473ffffffffffffffffffffffffffffffffffffffff8481167fffffffffffffffffffffffff00000000000000000000000000000000000000008316811790935560405191169283917fe9cf53972264dc95304fd424458745019ddfca0e37ae8f703d74772c41ad115b9190a35050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610c96576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610638565b6000610cc17f504f4f4c00000000000000000000000000000000000000000000000000000000611441565b9050610ced7f504f4f4c00000000000000000000000000000000000000000000000000000000836114f8565b8173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f90affc163f1a2dfedcd36aa02ed992eeeba8100a4014f0b4cdc20ea265a6662760405160405180910390a35050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610dcc576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610638565b60008281526002602052604080822080547fffffffffffffffffffffffff0000000000000000000000000000000000000000811673ffffffffffffffffffffffffffffffffffffffff8681169182179093559251911692839186917f9ef0e8c8e52743bb38b83b17d9429141d494b8041ca6d616a6c77cebae9cd8b791a4505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610ed0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610638565b7f444154415f50524f564944455200000000000000000000000000000000000000600090815260026020527fcd7944601aaa5cd7ccdae1bebec659e98c6aac8f12486b30e59db0d39698051f805473ffffffffffffffffffffffffffffffffffffffff8481167fffffffffffffffffffffffff00000000000000000000000000000000000000008316811790935560405191169283917fc853974cfbf81487a14a23565917bee63f527853bcb5fa54f2ae1cdf8a38356d9190a35050565b60005473ffffffffffffffffffffffffffffffffffffffff16331461100f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610638565b600061103a7f504f4f4c5f434f4e464947555241544f52000000000000000000000000000000611441565b90506110667f504f4f4c5f434f4e464947555241544f52000000000000000000000000000000836114f8565b8173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8932892569eba59c8382a089d9b732d1f49272878775235761a2a6b0309cd46560405160405180910390a35050565b60005473ffffffffffffffffffffffffffffffffffffffff163314611145576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610638565b7f41434c5f4d414e41474552000000000000000000000000000000000000000000600090815260026020527f9edef266ef35fd0c6e131df0f31a330f3dd4c4d19dd31ed615c21d005c68116b805473ffffffffffffffffffffffffffffffffffffffff8481167fffffffffffffffffffffffff00000000000000000000000000000000000000008316811790935560405191169283917fb30efa04327bb8a537d61cc1e5c48095345ad18ef7cc04e6bacf7dfb6caaf5079190a35050565b60005473ffffffffffffffffffffffffffffffffffffffff163314611284576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610638565b73ffffffffffffffffffffffffffffffffffffffff8116611327576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201527f64647265737300000000000000000000000000000000000000000000000000006064820152608401610638565b6000805460405173ffffffffffffffffffffffffffffffffffffffff808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a3600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60005473ffffffffffffffffffffffffffffffffffffffff163314611435576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610638565b61143e816117bf565b50565b60008181526002602052604081205473ffffffffffffffffffffffffffffffffffffffff16806114745750600092915050565b60008190508073ffffffffffffffffffffffffffffffffffffffff16635c60da1b6040518163ffffffff1660e01b81526004016020604051808303816000875af11580156114c6573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114ea9190611bca565b949350505050565b50919050565b60008281526002602052604080822054905130602482015273ffffffffffffffffffffffffffffffffffffffff90911691908190604401604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fc4d66de800000000000000000000000000000000000000000000000000000000179052905073ffffffffffffffffffffffffffffffffffffffff831661172e57306040516115cf906118bc565b73ffffffffffffffffffffffffffffffffffffffff9091168152602001604051809103906000f080158015611608573d6000803e3d6000fd5b506000868152600260205260409081902080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff841690811790915590517fd1f578940000000000000000000000000000000000000000000000000000000081529194508493509063d1f578949061169c9087908590600401611be7565b600060405180830381600087803b1580156116b657600080fd5b505af11580156116ca573d6000803e3d6000fd5b505050508373ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16867f4a465a9bd819d9662563c1e11ae958f8109e437e7f4bf1c6ef0b9a7b3f35d47860405160405180910390a46117b8565b6040517f4f1ef28600000000000000000000000000000000000000000000000000000000815283925073ffffffffffffffffffffffffffffffffffffffff831690634f1ef286906117859087908590600401611be7565b600060405180830381600087803b15801561179f57600080fd5b505af11580156117b3573d6000803e3d6000fd5b505050505b5050505050565b6000600180546117ce90611b7c565b80601f01602080910402602001604051908101604052809291908181526020018280546117fa90611b7c565b80156118475780601f1061181c57610100808354040283529160200191611847565b820191906000526020600020905b81548152906001019060200180831161182a57829003601f168201915b50508551939450611863936001935060208701925090506118c9565b50816040516118729190611c16565b6040518091039020816040516118889190611c16565b604051908190038120907fe685c8cdecc6030c45030fd54778812cb84ed8e4467c38294403d68ba786082390600090a35050565b6109cb80611c3383390190565b8280546118d590611b7c565b90600052602060002090601f0160209004810192826118f7576000855561193d565b82601f1061191057805160ff191683800117855561193d565b8280016001018555821561193d579182015b8281111561193d578251825591602001919060010190611922565b5061194992915061194d565b5090565b5b80821115611949576000815560010161194e565b60006020828403121561197457600080fd5b5035919050565b73ffffffffffffffffffffffffffffffffffffffff8116811461143e57600080fd5b6000602082840312156119af57600080fd5b81356119ba8161197b565b9392505050565b60005b838110156119dc5781810151838201526020016119c4565b838111156119eb576000848401525b50505050565b60008151808452611a098160208601602086016119c1565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006119ba60208301846119f1565b60008060408385031215611a6157600080fd5b823591506020830135611a738161197b565b809150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600060208284031215611abf57600080fd5b813567ffffffffffffffff80821115611ad757600080fd5b818401915084601f830112611aeb57600080fd5b813581811115611afd57611afd611a7e565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f01168101908382118183101715611b4357611b43611a7e565b81604052828152876020848701011115611b5c57600080fd5b826020860160208301376000928101602001929092525095945050505050565b600181811c90821680611b9057607f821691505b602082108114156114f2577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b600060208284031215611bdc57600080fd5b81516119ba8161197b565b73ffffffffffffffffffffffffffffffffffffffff831681526040602082015260006114ea60408301846119f1565b60008251611c288184602087016119c1565b919091019291505056fe60a060405234801561001057600080fd5b506040516109cb3803806109cb83398101604081905261002f91610040565b6001600160a01b0316608052610070565b60006020828403121561005257600080fd5b81516001600160a01b038116811461006957600080fd5b9392505050565b60805161091d6100ae6000396000818161014f015281816101a101528181610274015281816104110152818161043a01526105a4015261091d6000f3fe60806040526004361061005a5760003560e01c80635c60da1b116100435780635c60da1b14610097578063d1f57894146100d5578063f851a440146100e85761005a565b80633659cfe6146100645780634f1ef28614610084575b6100626100fd565b005b34801561007057600080fd5b5061006261007f36600461067b565b610137565b61006261009236600461069d565b610189565b3480156100a357600080fd5b506100ac61025a565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390f35b6100626100e336600461074f565b6102cb565b3480156100f457600080fd5b506100ac6103f7565b61010561045c565b6101356101307f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5490565b610464565b565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614156101815761017e81610488565b50565b61017e6100fd565b3373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016141561024d576101d083610488565b60008373ffffffffffffffffffffffffffffffffffffffff1683836040516101f992919061082f565b600060405180830381855af49150503d8060008114610234576040519150601f19603f3d011682016040523d82523d6000602084013e610239565b606091505b505090508061024757600080fd5b50505050565b6102556100fd565b505050565b60003373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614156102c057507f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5490565b6102c86100fd565b90565b60006102f57f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5490565b73ffffffffffffffffffffffffffffffffffffffff161461031557600080fd5b61034060017f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbd61083f565b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc1461036e5761036e61087d565b610377826104d5565b8051156103f35760008273ffffffffffffffffffffffffffffffffffffffff16826040516103a591906108ac565b600060405180830381855af49150503d80600081146103e0576040519150601f19603f3d011682016040523d82523d6000602084013e6103e5565b606091505b505090508061025557600080fd5b5050565b60003373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614156102c057507f000000000000000000000000000000000000000000000000000000000000000090565b61013561058c565b3660008037600080366000845af43d6000803e808015610483573d6000f35b3d6000fd5b610491816104d5565b60405173ffffffffffffffffffffffffffffffffffffffff8216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b803b610568576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603b60248201527f43616e6e6f742073657420612070726f787920696d706c656d656e746174696f60448201527f6e20746f2061206e6f6e2d636f6e74726163742061646472657373000000000060648201526084015b60405180910390fd5b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc55565b3373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161415610135576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603260248201527f43616e6e6f742063616c6c2066616c6c6261636b2066756e6374696f6e20667260448201527f6f6d207468652070726f78792061646d696e0000000000000000000000000000606482015260840161055f565b803573ffffffffffffffffffffffffffffffffffffffff8116811461067657600080fd5b919050565b60006020828403121561068d57600080fd5b61069682610652565b9392505050565b6000806000604084860312156106b257600080fd5b6106bb84610652565b9250602084013567ffffffffffffffff808211156106d857600080fd5b818601915086601f8301126106ec57600080fd5b8135818111156106fb57600080fd5b87602082850101111561070d57600080fd5b6020830194508093505050509250925092565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6000806040838503121561076257600080fd5b61076b83610652565b9150602083013567ffffffffffffffff8082111561078857600080fd5b818501915085601f83011261079c57600080fd5b8135818111156107ae576107ae610720565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f011681019083821181831017156107f4576107f4610720565b8160405282815288602084870101111561080d57600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b8183823760009101908152919050565b600082821015610878577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b500390565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052600160045260246000fd5b6000825160005b818110156108cd57602081860181015185830152016108b3565b818111156108dc576000828501525b50919091019291505056fea2646970667358221220899ba9574e8c52c72539176723d8c74a8618334587150196e2029371e7486a8464736f6c634300080a0033a2646970667358221220e62d78f287b0c390f4963ebdc446c3b33cd341f5cc839e23b26c90067a0b899564736f6c634300080a00338be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0",
}

// PoolAddressesProviderABI is the input ABI used to generate the binding from.
// Deprecated: Use PoolAddressesProviderMetaData.ABI instead.
var PoolAddressesProviderABI = PoolAddressesProviderMetaData.ABI

// PoolAddressesProviderBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use PoolAddressesProviderMetaData.Bin instead.
var PoolAddressesProviderBin = PoolAddressesProviderMetaData.Bin

// DeployPoolAddressesProvider deploys a new Ethereum contract, binding an instance of PoolAddressesProvider to it.
func DeployPoolAddressesProvider(auth *bind.TransactOpts, backend bind.ContractBackend, marketId string, owner common.Address) (common.Address, *types.Transaction, *PoolAddressesProvider, error) {
	parsed, err := PoolAddressesProviderMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PoolAddressesProviderBin), backend, marketId, owner)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &PoolAddressesProvider{PoolAddressesProviderCaller: PoolAddressesProviderCaller{contract: contract}, PoolAddressesProviderTransactor: PoolAddressesProviderTransactor{contract: contract}, PoolAddressesProviderFilterer: PoolAddressesProviderFilterer{contract: contract}}, nil
}

// PoolAddressesProvider is an auto generated Go binding around an Ethereum contract.
type PoolAddressesProvider struct {
	PoolAddressesProviderCaller     // Read-only binding to the contract
	PoolAddressesProviderTransactor // Write-only binding to the contract
	PoolAddressesProviderFilterer   // Log filterer for contract events
}

// PoolAddressesProviderCaller is an auto generated read-only Go binding around an Ethereum contract.
type PoolAddressesProviderCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PoolAddressesProviderTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PoolAddressesProviderTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PoolAddressesProviderFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PoolAddressesProviderFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PoolAddressesProviderSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PoolAddressesProviderSession struct {
	Contract     *PoolAddressesProvider // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// PoolAddressesProviderCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PoolAddressesProviderCallerSession struct {
	Contract *PoolAddressesProviderCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// PoolAddressesProviderTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PoolAddressesProviderTransactorSession struct {
	Contract     *PoolAddressesProviderTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// PoolAddressesProviderRaw is an auto generated low-level Go binding around an Ethereum contract.
type PoolAddressesProviderRaw struct {
	Contract *PoolAddressesProvider // Generic contract binding to access the raw methods on
}

// PoolAddressesProviderCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PoolAddressesProviderCallerRaw struct {
	Contract *PoolAddressesProviderCaller // Generic read-only contract binding to access the raw methods on
}

// PoolAddressesProviderTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PoolAddressesProviderTransactorRaw struct {
	Contract *PoolAddressesProviderTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPoolAddressesProvider creates a new instance of PoolAddressesProvider, bound to a specific deployed contract.
func NewPoolAddressesProvider(address common.Address, backend bind.ContractBackend) (*PoolAddressesProvider, error) {
	contract, err := bindPoolAddressesProvider(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PoolAddressesProvider{PoolAddressesProviderCaller: PoolAddressesProviderCaller{contract: contract}, PoolAddressesProviderTransactor: PoolAddressesProviderTransactor{contract: contract}, PoolAddressesProviderFilterer: PoolAddressesProviderFilterer{contract: contract}}, nil
}

// NewPoolAddressesProviderCaller creates a new read-only instance of PoolAddressesProvider, bound to a specific deployed contract.
func NewPoolAddressesProviderCaller(address common.Address, caller bind.ContractCaller) (*PoolAddressesProviderCaller, error) {
	contract, err := bindPoolAddressesProvider(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PoolAddressesProviderCaller{contract: contract}, nil
}

// NewPoolAddressesProviderTransactor creates a new write-only instance of PoolAddressesProvider, bound to a specific deployed contract.
func NewPoolAddressesProviderTransactor(address common.Address, transactor bind.ContractTransactor) (*PoolAddressesProviderTransactor, error) {
	contract, err := bindPoolAddressesProvider(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PoolAddressesProviderTransactor{contract: contract}, nil
}

// NewPoolAddressesProviderFilterer creates a new log filterer instance of PoolAddressesProvider, bound to a specific deployed contract.
func NewPoolAddressesProviderFilterer(address common.Address, filterer bind.ContractFilterer) (*PoolAddressesProviderFilterer, error) {
	contract, err := bindPoolAddressesProvider(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PoolAddressesProviderFilterer{contract: contract}, nil
}

// bindPoolAddressesProvider binds a generic wrapper to an already deployed contract.
func bindPoolAddressesProvider(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PoolAddressesProviderMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PoolAddressesProvider *PoolAddressesProviderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PoolAddressesProvider.Contract.PoolAddressesProviderCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PoolAddressesProvider *PoolAddressesProviderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PoolAddressesProvider.Contract.PoolAddressesProviderTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PoolAddressesProvider *PoolAddressesProviderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PoolAddressesProvider.Contract.PoolAddressesProviderTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PoolAddressesProvider *PoolAddressesProviderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PoolAddressesProvider.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PoolAddressesProvider *PoolAddressesProviderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PoolAddressesProvider.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PoolAddressesProvider *PoolAddressesProviderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PoolAddressesProvider.Contract.contract.Transact(opts, method, params...)
}

// GetACLAdmin is a free data retrieval call binding the contract method 0x0e67178c.
//
// Solidity: function getACLAdmin() view returns(address)
func (_PoolAddressesProvider *PoolAddressesProviderCaller) GetACLAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PoolAddressesProvider.contract.Call(opts, &out, "getACLAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetACLAdmin is a free data retrieval call binding the contract method 0x0e67178c.
//
// Solidity: function getACLAdmin() view returns(address)
func (_PoolAddressesProvider *PoolAddressesProviderSession) GetACLAdmin() (common.Address, error) {
	return _PoolAddressesProvider.Contract.GetACLAdmin(&_PoolAddressesProvider.CallOpts)
}

// GetACLAdmin is a free data retrieval call binding the contract method 0x0e67178c.
//
// Solidity: function getACLAdmin() view returns(address)
func (_PoolAddressesProvider *PoolAddressesProviderCallerSession) GetACLAdmin() (common.Address, error) {
	return _PoolAddressesProvider.Contract.GetACLAdmin(&_PoolAddressesProvider.CallOpts)
}

// GetACLManager is a free data retrieval call binding the contract method 0x707cd716.
//
// Solidity: function getACLManager() view returns(address)
func (_PoolAddressesProvider *PoolAddressesProviderCaller) GetACLManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PoolAddressesProvider.contract.Call(opts, &out, "getACLManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetACLManager is a free data retrieval call binding the contract method 0x707cd716.
//
// Solidity: function getACLManager() view returns(address)
func (_PoolAddressesProvider *PoolAddressesProviderSession) GetACLManager() (common.Address, error) {
	return _PoolAddressesProvider.Contract.GetACLManager(&_PoolAddressesProvider.CallOpts)
}

// GetACLManager is a free data retrieval call binding the contract method 0x707cd716.
//
// Solidity: function getACLManager() view returns(address)
func (_PoolAddressesProvider *PoolAddressesProviderCallerSession) GetACLManager() (common.Address, error) {
	return _PoolAddressesProvider.Contract.GetACLManager(&_PoolAddressesProvider.CallOpts)
}

// GetAddress is a free data retrieval call binding the contract method 0x21f8a721.
//
// Solidity: function getAddress(bytes32 id) view returns(address)
func (_PoolAddressesProvider *PoolAddressesProviderCaller) GetAddress(opts *bind.CallOpts, id [32]byte) (common.Address, error) {
	var out []interface{}
	err := _PoolAddressesProvider.contract.Call(opts, &out, "getAddress", id)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAddress is a free data retrieval call binding the contract method 0x21f8a721.
//
// Solidity: function getAddress(bytes32 id) view returns(address)
func (_PoolAddressesProvider *PoolAddressesProviderSession) GetAddress(id [32]byte) (common.Address, error) {
	return _PoolAddressesProvider.Contract.GetAddress(&_PoolAddressesProvider.CallOpts, id)
}

// GetAddress is a free data retrieval call binding the contract method 0x21f8a721.
//
// Solidity: function getAddress(bytes32 id) view returns(address)
func (_PoolAddressesProvider *PoolAddressesProviderCallerSession) GetAddress(id [32]byte) (common.Address, error) {
	return _PoolAddressesProvider.Contract.GetAddress(&_PoolAddressesProvider.CallOpts, id)
}

// GetMarketId is a free data retrieval call binding the contract method 0x568ef470.
//
// Solidity: function getMarketId() view returns(string)
func (_PoolAddressesProvider *PoolAddressesProviderCaller) GetMarketId(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _PoolAddressesProvider.contract.Call(opts, &out, "getMarketId")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetMarketId is a free data retrieval call binding the contract method 0x568ef470.
//
// Solidity: function getMarketId() view returns(string)
func (_PoolAddressesProvider *PoolAddressesProviderSession) GetMarketId() (string, error) {
	return _PoolAddressesProvider.Contract.GetMarketId(&_PoolAddressesProvider.CallOpts)
}

// GetMarketId is a free data retrieval call binding the contract method 0x568ef470.
//
// Solidity: function getMarketId() view returns(string)
func (_PoolAddressesProvider *PoolAddressesProviderCallerSession) GetMarketId() (string, error) {
	return _PoolAddressesProvider.Contract.GetMarketId(&_PoolAddressesProvider.CallOpts)
}

// GetPool is a free data retrieval call binding the contract method 0x026b1d5f.
//
// Solidity: function getPool() view returns(address)
func (_PoolAddressesProvider *PoolAddressesProviderCaller) GetPool(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PoolAddressesProvider.contract.Call(opts, &out, "getPool")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetPool is a free data retrieval call binding the contract method 0x026b1d5f.
//
// Solidity: function getPool() view returns(address)
func (_PoolAddressesProvider *PoolAddressesProviderSession) GetPool() (common.Address, error) {
	return _PoolAddressesProvider.Contract.GetPool(&_PoolAddressesProvider.CallOpts)
}

// GetPool is a free data retrieval call binding the contract method 0x026b1d5f.
//
// Solidity: function getPool() view returns(address)
func (_PoolAddressesProvider *PoolAddressesProviderCallerSession) GetPool() (common.Address, error) {
	return _PoolAddressesProvider.Contract.GetPool(&_PoolAddressesProvider.CallOpts)
}

// GetPoolConfigurator is a free data retrieval call binding the contract method 0x631adfca.
//
// Solidity: function getPoolConfigurator() view returns(address)
func (_PoolAddressesProvider *PoolAddressesProviderCaller) GetPoolConfigurator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PoolAddressesProvider.contract.Call(opts, &out, "getPoolConfigurator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetPoolConfigurator is a free data retrieval call binding the contract method 0x631adfca.
//
// Solidity: function getPoolConfigurator() view returns(address)
func (_PoolAddressesProvider *PoolAddressesProviderSession) GetPoolConfigurator() (common.Address, error) {
	return _PoolAddressesProvider.Contract.GetPoolConfigurator(&_PoolAddressesProvider.CallOpts)
}

// GetPoolConfigurator is a free data retrieval call binding the contract method 0x631adfca.
//
// Solidity: function getPoolConfigurator() view returns(address)
func (_PoolAddressesProvider *PoolAddressesProviderCallerSession) GetPoolConfigurator() (common.Address, error) {
	return _PoolAddressesProvider.Contract.GetPoolConfigurator(&_PoolAddressesProvider.CallOpts)
}

// GetPoolDataProvider is a free data retrieval call binding the contract method 0xe860accb.
//
// Solidity: function getPoolDataProvider() view returns(address)
func (_PoolAddressesProvider *PoolAddressesProviderCaller) GetPoolDataProvider(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PoolAddressesProvider.contract.Call(opts, &out, "getPoolDataProvider")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetPoolDataProvider is a free data retrieval call binding the contract method 0xe860accb.
//
// Solidity: function getPoolDataProvider() view returns(address)
func (_PoolAddressesProvider *PoolAddressesProviderSession) GetPoolDataProvider() (common.Address, error) {
	return _PoolAddressesProvider.Contract.GetPoolDataProvider(&_PoolAddressesProvider.CallOpts)
}

// GetPoolDataProvider is a free data retrieval call binding the contract method 0xe860accb.
//
// Solidity: function getPoolDataProvider() view returns(address)
func (_PoolAddressesProvider *PoolAddressesProviderCallerSession) GetPoolDataProvider() (common.Address, error) {
	return _PoolAddressesProvider.Contract.GetPoolDataProvider(&_PoolAddressesProvider.CallOpts)
}

// GetPriceOracle is a free data retrieval call binding the contract method 0xfca513a8.
//
// Solidity: function getPriceOracle() view returns(address)
func (_PoolAddressesProvider *PoolAddressesProviderCaller) GetPriceOracle(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PoolAddressesProvider.contract.Call(opts, &out, "getPriceOracle")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetPriceOracle is a free data retrieval call binding the contract method 0xfca513a8.
//
// Solidity: function getPriceOracle() view returns(address)
func (_PoolAddressesProvider *PoolAddressesProviderSession) GetPriceOracle() (common.Address, error) {
	return _PoolAddressesProvider.Contract.GetPriceOracle(&_PoolAddressesProvider.CallOpts)
}

// GetPriceOracle is a free data retrieval call binding the contract method 0xfca513a8.
//
// Solidity: function getPriceOracle() view returns(address)
func (_PoolAddressesProvider *PoolAddressesProviderCallerSession) GetPriceOracle() (common.Address, error) {
	return _PoolAddressesProvider.Contract.GetPriceOracle(&_PoolAddressesProvider.CallOpts)
}

// GetPriceOracleSentinel is a free data retrieval call binding the contract method 0x5eb88d3d.
//
// Solidity: function getPriceOracleSentinel() view returns(address)
func (_PoolAddressesProvider *PoolAddressesProviderCaller) GetPriceOracleSentinel(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PoolAddressesProvider.contract.Call(opts, &out, "getPriceOracleSentinel")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetPriceOracleSentinel is a free data retrieval call binding the contract method 0x5eb88d3d.
//
// Solidity: function getPriceOracleSentinel() view returns(address)
func (_PoolAddressesProvider *PoolAddressesProviderSession) GetPriceOracleSentinel() (common.Address, error) {
	return _PoolAddressesProvider.Contract.GetPriceOracleSentinel(&_PoolAddressesProvider.CallOpts)
}

// GetPriceOracleSentinel is a free data retrieval call binding the contract method 0x5eb88d3d.
//
// Solidity: function getPriceOracleSentinel() view returns(address)
func (_PoolAddressesProvider *PoolAddressesProviderCallerSession) GetPriceOracleSentinel() (common.Address, error) {
	return _PoolAddressesProvider.Contract.GetPriceOracleSentinel(&_PoolAddressesProvider.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PoolAddressesProvider *PoolAddressesProviderCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PoolAddressesProvider.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PoolAddressesProvider *PoolAddressesProviderSession) Owner() (common.Address, error) {
	return _PoolAddressesProvider.Contract.Owner(&_PoolAddressesProvider.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PoolAddressesProvider *PoolAddressesProviderCallerSession) Owner() (common.Address, error) {
	return _PoolAddressesProvider.Contract.Owner(&_PoolAddressesProvider.CallOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_PoolAddressesProvider *PoolAddressesProviderTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PoolAddressesProvider.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_PoolAddressesProvider *PoolAddressesProviderSession) RenounceOwnership() (*types.Transaction, error) {
	return _PoolAddressesProvider.Contract.RenounceOwnership(&_PoolAddressesProvider.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_PoolAddressesProvider *PoolAddressesProviderTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _PoolAddressesProvider.Contract.RenounceOwnership(&_PoolAddressesProvider.TransactOpts)
}

// SetACLAdmin is a paid mutator transaction binding the contract method 0x76d84ffc.
//
// Solidity: function setACLAdmin(address newAclAdmin) returns()
func (_PoolAddressesProvider *PoolAddressesProviderTransactor) SetACLAdmin(opts *bind.TransactOpts, newAclAdmin common.Address) (*types.Transaction, error) {
	return _PoolAddressesProvider.contract.Transact(opts, "setACLAdmin", newAclAdmin)
}

// SetACLAdmin is a paid mutator transaction binding the contract method 0x76d84ffc.
//
// Solidity: function setACLAdmin(address newAclAdmin) returns()
func (_PoolAddressesProvider *PoolAddressesProviderSession) SetACLAdmin(newAclAdmin common.Address) (*types.Transaction, error) {
	return _PoolAddressesProvider.Contract.SetACLAdmin(&_PoolAddressesProvider.TransactOpts, newAclAdmin)
}

// SetACLAdmin is a paid mutator transaction binding the contract method 0x76d84ffc.
//
// Solidity: function setACLAdmin(address newAclAdmin) returns()
func (_PoolAddressesProvider *PoolAddressesProviderTransactorSession) SetACLAdmin(newAclAdmin common.Address) (*types.Transaction, error) {
	return _PoolAddressesProvider.Contract.SetACLAdmin(&_PoolAddressesProvider.TransactOpts, newAclAdmin)
}

// SetACLManager is a paid mutator transaction binding the contract method 0xed301ca9.
//
// Solidity: function setACLManager(address newAclManager) returns()
func (_PoolAddressesProvider *PoolAddressesProviderTransactor) SetACLManager(opts *bind.TransactOpts, newAclManager common.Address) (*types.Transaction, error) {
	return _PoolAddressesProvider.contract.Transact(opts, "setACLManager", newAclManager)
}

// SetACLManager is a paid mutator transaction binding the contract method 0xed301ca9.
//
// Solidity: function setACLManager(address newAclManager) returns()
func (_PoolAddressesProvider *PoolAddressesProviderSession) SetACLManager(newAclManager common.Address) (*types.Transaction, error) {
	return _PoolAddressesProvider.Contract.SetACLManager(&_PoolAddressesProvider.TransactOpts, newAclManager)
}

// SetACLManager is a paid mutator transaction binding the contract method 0xed301ca9.
//
// Solidity: function setACLManager(address newAclManager) returns()
func (_PoolAddressesProvider *PoolAddressesProviderTransactorSession) SetACLManager(newAclManager common.Address) (*types.Transaction, error) {
	return _PoolAddressesProvider.Contract.SetACLManager(&_PoolAddressesProvider.TransactOpts, newAclManager)
}

// SetAddress is a paid mutator transaction binding the contract method 0xca446dd9.
//
// Solidity: function setAddress(bytes32 id, address newAddress) returns()
func (_PoolAddressesProvider *PoolAddressesProviderTransactor) SetAddress(opts *bind.TransactOpts, id [32]byte, newAddress common.Address) (*types.Transaction, error) {
	return _PoolAddressesProvider.contract.Transact(opts, "setAddress", id, newAddress)
}

// SetAddress is a paid mutator transaction binding the contract method 0xca446dd9.
//
// Solidity: function setAddress(bytes32 id, address newAddress) returns()
func (_PoolAddressesProvider *PoolAddressesProviderSession) SetAddress(id [32]byte, newAddress common.Address) (*types.Transaction, error) {
	return _PoolAddressesProvider.Contract.SetAddress(&_PoolAddressesProvider.TransactOpts, id, newAddress)
}

// SetAddress is a paid mutator transaction binding the contract method 0xca446dd9.
//
// Solidity: function setAddress(bytes32 id, address newAddress) returns()
func (_PoolAddressesProvider *PoolAddressesProviderTransactorSession) SetAddress(id [32]byte, newAddress common.Address) (*types.Transaction, error) {
	return _PoolAddressesProvider.Contract.SetAddress(&_PoolAddressesProvider.TransactOpts, id, newAddress)
}

// SetAddressAsProxy is a paid mutator transaction binding the contract method 0x5dcc528c.
//
// Solidity: function setAddressAsProxy(bytes32 id, address newImplementationAddress) returns()
func (_PoolAddressesProvider *PoolAddressesProviderTransactor) SetAddressAsProxy(opts *bind.TransactOpts, id [32]byte, newImplementationAddress common.Address) (*types.Transaction, error) {
	return _PoolAddressesProvider.contract.Transact(opts, "setAddressAsProxy", id, newImplementationAddress)
}

// SetAddressAsProxy is a paid mutator transaction binding the contract method 0x5dcc528c.
//
// Solidity: function setAddressAsProxy(bytes32 id, address newImplementationAddress) returns()
func (_PoolAddressesProvider *PoolAddressesProviderSession) SetAddressAsProxy(id [32]byte, newImplementationAddress common.Address) (*types.Transaction, error) {
	return _PoolAddressesProvider.Contract.SetAddressAsProxy(&_PoolAddressesProvider.TransactOpts, id, newImplementationAddress)
}

// SetAddressAsProxy is a paid mutator transaction binding the contract method 0x5dcc528c.
//
// Solidity: function setAddressAsProxy(bytes32 id, address newImplementationAddress) returns()
func (_PoolAddressesProvider *PoolAddressesProviderTransactorSession) SetAddressAsProxy(id [32]byte, newImplementationAddress common.Address) (*types.Transaction, error) {
	return _PoolAddressesProvider.Contract.SetAddressAsProxy(&_PoolAddressesProvider.TransactOpts, id, newImplementationAddress)
}

// SetMarketId is a paid mutator transaction binding the contract method 0xf67b1847.
//
// Solidity: function setMarketId(string newMarketId) returns()
func (_PoolAddressesProvider *PoolAddressesProviderTransactor) SetMarketId(opts *bind.TransactOpts, newMarketId string) (*types.Transaction, error) {
	return _PoolAddressesProvider.contract.Transact(opts, "setMarketId", newMarketId)
}

// SetMarketId is a paid mutator transaction binding the contract method 0xf67b1847.
//
// Solidity: function setMarketId(string newMarketId) returns()
func (_PoolAddressesProvider *PoolAddressesProviderSession) SetMarketId(newMarketId string) (*types.Transaction, error) {
	return _PoolAddressesProvider.Contract.SetMarketId(&_PoolAddressesProvider.TransactOpts, newMarketId)
}

// SetMarketId is a paid mutator transaction binding the contract method 0xf67b1847.
//
// Solidity: function setMarketId(string newMarketId) returns()
func (_PoolAddressesProvider *PoolAddressesProviderTransactorSession) SetMarketId(newMarketId string) (*types.Transaction, error) {
	return _PoolAddressesProvider.Contract.SetMarketId(&_PoolAddressesProvider.TransactOpts, newMarketId)
}

// SetPoolConfiguratorImpl is a paid mutator transaction binding the contract method 0xe4ca28b7.
//
// Solidity: function setPoolConfiguratorImpl(address newPoolConfiguratorImpl) returns()
func (_PoolAddressesProvider *PoolAddressesProviderTransactor) SetPoolConfiguratorImpl(opts *bind.TransactOpts, newPoolConfiguratorImpl common.Address) (*types.Transaction, error) {
	return _PoolAddressesProvider.contract.Transact(opts, "setPoolConfiguratorImpl", newPoolConfiguratorImpl)
}

// SetPoolConfiguratorImpl is a paid mutator transaction binding the contract method 0xe4ca28b7.
//
// Solidity: function setPoolConfiguratorImpl(address newPoolConfiguratorImpl) returns()
func (_PoolAddressesProvider *PoolAddressesProviderSession) SetPoolConfiguratorImpl(newPoolConfiguratorImpl common.Address) (*types.Transaction, error) {
	return _PoolAddressesProvider.Contract.SetPoolConfiguratorImpl(&_PoolAddressesProvider.TransactOpts, newPoolConfiguratorImpl)
}

// SetPoolConfiguratorImpl is a paid mutator transaction binding the contract method 0xe4ca28b7.
//
// Solidity: function setPoolConfiguratorImpl(address newPoolConfiguratorImpl) returns()
func (_PoolAddressesProvider *PoolAddressesProviderTransactorSession) SetPoolConfiguratorImpl(newPoolConfiguratorImpl common.Address) (*types.Transaction, error) {
	return _PoolAddressesProvider.Contract.SetPoolConfiguratorImpl(&_PoolAddressesProvider.TransactOpts, newPoolConfiguratorImpl)
}

// SetPoolDataProvider is a paid mutator transaction binding the contract method 0xe44e9ed1.
//
// Solidity: function setPoolDataProvider(address newDataProvider) returns()
func (_PoolAddressesProvider *PoolAddressesProviderTransactor) SetPoolDataProvider(opts *bind.TransactOpts, newDataProvider common.Address) (*types.Transaction, error) {
	return _PoolAddressesProvider.contract.Transact(opts, "setPoolDataProvider", newDataProvider)
}

// SetPoolDataProvider is a paid mutator transaction binding the contract method 0xe44e9ed1.
//
// Solidity: function setPoolDataProvider(address newDataProvider) returns()
func (_PoolAddressesProvider *PoolAddressesProviderSession) SetPoolDataProvider(newDataProvider common.Address) (*types.Transaction, error) {
	return _PoolAddressesProvider.Contract.SetPoolDataProvider(&_PoolAddressesProvider.TransactOpts, newDataProvider)
}

// SetPoolDataProvider is a paid mutator transaction binding the contract method 0xe44e9ed1.
//
// Solidity: function setPoolDataProvider(address newDataProvider) returns()
func (_PoolAddressesProvider *PoolAddressesProviderTransactorSession) SetPoolDataProvider(newDataProvider common.Address) (*types.Transaction, error) {
	return _PoolAddressesProvider.Contract.SetPoolDataProvider(&_PoolAddressesProvider.TransactOpts, newDataProvider)
}

// SetPoolImpl is a paid mutator transaction binding the contract method 0xa1564406.
//
// Solidity: function setPoolImpl(address newPoolImpl) returns()
func (_PoolAddressesProvider *PoolAddressesProviderTransactor) SetPoolImpl(opts *bind.TransactOpts, newPoolImpl common.Address) (*types.Transaction, error) {
	return _PoolAddressesProvider.contract.Transact(opts, "setPoolImpl", newPoolImpl)
}

// SetPoolImpl is a paid mutator transaction binding the contract method 0xa1564406.
//
// Solidity: function setPoolImpl(address newPoolImpl) returns()
func (_PoolAddressesProvider *PoolAddressesProviderSession) SetPoolImpl(newPoolImpl common.Address) (*types.Transaction, error) {
	return _PoolAddressesProvider.Contract.SetPoolImpl(&_PoolAddressesProvider.TransactOpts, newPoolImpl)
}

// SetPoolImpl is a paid mutator transaction binding the contract method 0xa1564406.
//
// Solidity: function setPoolImpl(address newPoolImpl) returns()
func (_PoolAddressesProvider *PoolAddressesProviderTransactorSession) SetPoolImpl(newPoolImpl common.Address) (*types.Transaction, error) {
	return _PoolAddressesProvider.Contract.SetPoolImpl(&_PoolAddressesProvider.TransactOpts, newPoolImpl)
}

// SetPriceOracle is a paid mutator transaction binding the contract method 0x530e784f.
//
// Solidity: function setPriceOracle(address newPriceOracle) returns()
func (_PoolAddressesProvider *PoolAddressesProviderTransactor) SetPriceOracle(opts *bind.TransactOpts, newPriceOracle common.Address) (*types.Transaction, error) {
	return _PoolAddressesProvider.contract.Transact(opts, "setPriceOracle", newPriceOracle)
}

// SetPriceOracle is a paid mutator transaction binding the contract method 0x530e784f.
//
// Solidity: function setPriceOracle(address newPriceOracle) returns()
func (_PoolAddressesProvider *PoolAddressesProviderSession) SetPriceOracle(newPriceOracle common.Address) (*types.Transaction, error) {
	return _PoolAddressesProvider.Contract.SetPriceOracle(&_PoolAddressesProvider.TransactOpts, newPriceOracle)
}

// SetPriceOracle is a paid mutator transaction binding the contract method 0x530e784f.
//
// Solidity: function setPriceOracle(address newPriceOracle) returns()
func (_PoolAddressesProvider *PoolAddressesProviderTransactorSession) SetPriceOracle(newPriceOracle common.Address) (*types.Transaction, error) {
	return _PoolAddressesProvider.Contract.SetPriceOracle(&_PoolAddressesProvider.TransactOpts, newPriceOracle)
}

// SetPriceOracleSentinel is a paid mutator transaction binding the contract method 0x74944cec.
//
// Solidity: function setPriceOracleSentinel(address newPriceOracleSentinel) returns()
func (_PoolAddressesProvider *PoolAddressesProviderTransactor) SetPriceOracleSentinel(opts *bind.TransactOpts, newPriceOracleSentinel common.Address) (*types.Transaction, error) {
	return _PoolAddressesProvider.contract.Transact(opts, "setPriceOracleSentinel", newPriceOracleSentinel)
}

// SetPriceOracleSentinel is a paid mutator transaction binding the contract method 0x74944cec.
//
// Solidity: function setPriceOracleSentinel(address newPriceOracleSentinel) returns()
func (_PoolAddressesProvider *PoolAddressesProviderSession) SetPriceOracleSentinel(newPriceOracleSentinel common.Address) (*types.Transaction, error) {
	return _PoolAddressesProvider.Contract.SetPriceOracleSentinel(&_PoolAddressesProvider.TransactOpts, newPriceOracleSentinel)
}

// SetPriceOracleSentinel is a paid mutator transaction binding the contract method 0x74944cec.
//
// Solidity: function setPriceOracleSentinel(address newPriceOracleSentinel) returns()
func (_PoolAddressesProvider *PoolAddressesProviderTransactorSession) SetPriceOracleSentinel(newPriceOracleSentinel common.Address) (*types.Transaction, error) {
	return _PoolAddressesProvider.Contract.SetPriceOracleSentinel(&_PoolAddressesProvider.TransactOpts, newPriceOracleSentinel)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_PoolAddressesProvider *PoolAddressesProviderTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _PoolAddressesProvider.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_PoolAddressesProvider *PoolAddressesProviderSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _PoolAddressesProvider.Contract.TransferOwnership(&_PoolAddressesProvider.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_PoolAddressesProvider *PoolAddressesProviderTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _PoolAddressesProvider.Contract.TransferOwnership(&_PoolAddressesProvider.TransactOpts, newOwner)
}

// PoolAddressesProviderACLAdminUpdatedIterator is returned from FilterACLAdminUpdated and is used to iterate over the raw logs and unpacked data for ACLAdminUpdated events raised by the PoolAddressesProvider contract.
type PoolAddressesProviderACLAdminUpdatedIterator struct {
	Event *PoolAddressesProviderACLAdminUpdated // Event containing the contract specifics and raw log

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
func (it *PoolAddressesProviderACLAdminUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolAddressesProviderACLAdminUpdated)
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
		it.Event = new(PoolAddressesProviderACLAdminUpdated)
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
func (it *PoolAddressesProviderACLAdminUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolAddressesProviderACLAdminUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolAddressesProviderACLAdminUpdated represents a ACLAdminUpdated event raised by the PoolAddressesProvider contract.
type PoolAddressesProviderACLAdminUpdated struct {
	OldAddress common.Address
	NewAddress common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterACLAdminUpdated is a free log retrieval operation binding the contract event 0xe9cf53972264dc95304fd424458745019ddfca0e37ae8f703d74772c41ad115b.
//
// Solidity: event ACLAdminUpdated(address indexed oldAddress, address indexed newAddress)
func (_PoolAddressesProvider *PoolAddressesProviderFilterer) FilterACLAdminUpdated(opts *bind.FilterOpts, oldAddress []common.Address, newAddress []common.Address) (*PoolAddressesProviderACLAdminUpdatedIterator, error) {

	var oldAddressRule []interface{}
	for _, oldAddressItem := range oldAddress {
		oldAddressRule = append(oldAddressRule, oldAddressItem)
	}
	var newAddressRule []interface{}
	for _, newAddressItem := range newAddress {
		newAddressRule = append(newAddressRule, newAddressItem)
	}

	logs, sub, err := _PoolAddressesProvider.contract.FilterLogs(opts, "ACLAdminUpdated", oldAddressRule, newAddressRule)
	if err != nil {
		return nil, err
	}
	return &PoolAddressesProviderACLAdminUpdatedIterator{contract: _PoolAddressesProvider.contract, event: "ACLAdminUpdated", logs: logs, sub: sub}, nil
}

// WatchACLAdminUpdated is a free log subscription operation binding the contract event 0xe9cf53972264dc95304fd424458745019ddfca0e37ae8f703d74772c41ad115b.
//
// Solidity: event ACLAdminUpdated(address indexed oldAddress, address indexed newAddress)
func (_PoolAddressesProvider *PoolAddressesProviderFilterer) WatchACLAdminUpdated(opts *bind.WatchOpts, sink chan<- *PoolAddressesProviderACLAdminUpdated, oldAddress []common.Address, newAddress []common.Address) (event.Subscription, error) {

	var oldAddressRule []interface{}
	for _, oldAddressItem := range oldAddress {
		oldAddressRule = append(oldAddressRule, oldAddressItem)
	}
	var newAddressRule []interface{}
	for _, newAddressItem := range newAddress {
		newAddressRule = append(newAddressRule, newAddressItem)
	}

	logs, sub, err := _PoolAddressesProvider.contract.WatchLogs(opts, "ACLAdminUpdated", oldAddressRule, newAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolAddressesProviderACLAdminUpdated)
				if err := _PoolAddressesProvider.contract.UnpackLog(event, "ACLAdminUpdated", log); err != nil {
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

// ParseACLAdminUpdated is a log parse operation binding the contract event 0xe9cf53972264dc95304fd424458745019ddfca0e37ae8f703d74772c41ad115b.
//
// Solidity: event ACLAdminUpdated(address indexed oldAddress, address indexed newAddress)
func (_PoolAddressesProvider *PoolAddressesProviderFilterer) ParseACLAdminUpdated(log types.Log) (*PoolAddressesProviderACLAdminUpdated, error) {
	event := new(PoolAddressesProviderACLAdminUpdated)
	if err := _PoolAddressesProvider.contract.UnpackLog(event, "ACLAdminUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolAddressesProviderACLManagerUpdatedIterator is returned from FilterACLManagerUpdated and is used to iterate over the raw logs and unpacked data for ACLManagerUpdated events raised by the PoolAddressesProvider contract.
type PoolAddressesProviderACLManagerUpdatedIterator struct {
	Event *PoolAddressesProviderACLManagerUpdated // Event containing the contract specifics and raw log

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
func (it *PoolAddressesProviderACLManagerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolAddressesProviderACLManagerUpdated)
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
		it.Event = new(PoolAddressesProviderACLManagerUpdated)
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
func (it *PoolAddressesProviderACLManagerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolAddressesProviderACLManagerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolAddressesProviderACLManagerUpdated represents a ACLManagerUpdated event raised by the PoolAddressesProvider contract.
type PoolAddressesProviderACLManagerUpdated struct {
	OldAddress common.Address
	NewAddress common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterACLManagerUpdated is a free log retrieval operation binding the contract event 0xb30efa04327bb8a537d61cc1e5c48095345ad18ef7cc04e6bacf7dfb6caaf507.
//
// Solidity: event ACLManagerUpdated(address indexed oldAddress, address indexed newAddress)
func (_PoolAddressesProvider *PoolAddressesProviderFilterer) FilterACLManagerUpdated(opts *bind.FilterOpts, oldAddress []common.Address, newAddress []common.Address) (*PoolAddressesProviderACLManagerUpdatedIterator, error) {

	var oldAddressRule []interface{}
	for _, oldAddressItem := range oldAddress {
		oldAddressRule = append(oldAddressRule, oldAddressItem)
	}
	var newAddressRule []interface{}
	for _, newAddressItem := range newAddress {
		newAddressRule = append(newAddressRule, newAddressItem)
	}

	logs, sub, err := _PoolAddressesProvider.contract.FilterLogs(opts, "ACLManagerUpdated", oldAddressRule, newAddressRule)
	if err != nil {
		return nil, err
	}
	return &PoolAddressesProviderACLManagerUpdatedIterator{contract: _PoolAddressesProvider.contract, event: "ACLManagerUpdated", logs: logs, sub: sub}, nil
}

// WatchACLManagerUpdated is a free log subscription operation binding the contract event 0xb30efa04327bb8a537d61cc1e5c48095345ad18ef7cc04e6bacf7dfb6caaf507.
//
// Solidity: event ACLManagerUpdated(address indexed oldAddress, address indexed newAddress)
func (_PoolAddressesProvider *PoolAddressesProviderFilterer) WatchACLManagerUpdated(opts *bind.WatchOpts, sink chan<- *PoolAddressesProviderACLManagerUpdated, oldAddress []common.Address, newAddress []common.Address) (event.Subscription, error) {

	var oldAddressRule []interface{}
	for _, oldAddressItem := range oldAddress {
		oldAddressRule = append(oldAddressRule, oldAddressItem)
	}
	var newAddressRule []interface{}
	for _, newAddressItem := range newAddress {
		newAddressRule = append(newAddressRule, newAddressItem)
	}

	logs, sub, err := _PoolAddressesProvider.contract.WatchLogs(opts, "ACLManagerUpdated", oldAddressRule, newAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolAddressesProviderACLManagerUpdated)
				if err := _PoolAddressesProvider.contract.UnpackLog(event, "ACLManagerUpdated", log); err != nil {
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

// ParseACLManagerUpdated is a log parse operation binding the contract event 0xb30efa04327bb8a537d61cc1e5c48095345ad18ef7cc04e6bacf7dfb6caaf507.
//
// Solidity: event ACLManagerUpdated(address indexed oldAddress, address indexed newAddress)
func (_PoolAddressesProvider *PoolAddressesProviderFilterer) ParseACLManagerUpdated(log types.Log) (*PoolAddressesProviderACLManagerUpdated, error) {
	event := new(PoolAddressesProviderACLManagerUpdated)
	if err := _PoolAddressesProvider.contract.UnpackLog(event, "ACLManagerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolAddressesProviderAddressSetIterator is returned from FilterAddressSet and is used to iterate over the raw logs and unpacked data for AddressSet events raised by the PoolAddressesProvider contract.
type PoolAddressesProviderAddressSetIterator struct {
	Event *PoolAddressesProviderAddressSet // Event containing the contract specifics and raw log

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
func (it *PoolAddressesProviderAddressSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolAddressesProviderAddressSet)
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
		it.Event = new(PoolAddressesProviderAddressSet)
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
func (it *PoolAddressesProviderAddressSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolAddressesProviderAddressSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolAddressesProviderAddressSet represents a AddressSet event raised by the PoolAddressesProvider contract.
type PoolAddressesProviderAddressSet struct {
	Id         [32]byte
	OldAddress common.Address
	NewAddress common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterAddressSet is a free log retrieval operation binding the contract event 0x9ef0e8c8e52743bb38b83b17d9429141d494b8041ca6d616a6c77cebae9cd8b7.
//
// Solidity: event AddressSet(bytes32 indexed id, address indexed oldAddress, address indexed newAddress)
func (_PoolAddressesProvider *PoolAddressesProviderFilterer) FilterAddressSet(opts *bind.FilterOpts, id [][32]byte, oldAddress []common.Address, newAddress []common.Address) (*PoolAddressesProviderAddressSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var oldAddressRule []interface{}
	for _, oldAddressItem := range oldAddress {
		oldAddressRule = append(oldAddressRule, oldAddressItem)
	}
	var newAddressRule []interface{}
	for _, newAddressItem := range newAddress {
		newAddressRule = append(newAddressRule, newAddressItem)
	}

	logs, sub, err := _PoolAddressesProvider.contract.FilterLogs(opts, "AddressSet", idRule, oldAddressRule, newAddressRule)
	if err != nil {
		return nil, err
	}
	return &PoolAddressesProviderAddressSetIterator{contract: _PoolAddressesProvider.contract, event: "AddressSet", logs: logs, sub: sub}, nil
}

// WatchAddressSet is a free log subscription operation binding the contract event 0x9ef0e8c8e52743bb38b83b17d9429141d494b8041ca6d616a6c77cebae9cd8b7.
//
// Solidity: event AddressSet(bytes32 indexed id, address indexed oldAddress, address indexed newAddress)
func (_PoolAddressesProvider *PoolAddressesProviderFilterer) WatchAddressSet(opts *bind.WatchOpts, sink chan<- *PoolAddressesProviderAddressSet, id [][32]byte, oldAddress []common.Address, newAddress []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var oldAddressRule []interface{}
	for _, oldAddressItem := range oldAddress {
		oldAddressRule = append(oldAddressRule, oldAddressItem)
	}
	var newAddressRule []interface{}
	for _, newAddressItem := range newAddress {
		newAddressRule = append(newAddressRule, newAddressItem)
	}

	logs, sub, err := _PoolAddressesProvider.contract.WatchLogs(opts, "AddressSet", idRule, oldAddressRule, newAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolAddressesProviderAddressSet)
				if err := _PoolAddressesProvider.contract.UnpackLog(event, "AddressSet", log); err != nil {
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

// ParseAddressSet is a log parse operation binding the contract event 0x9ef0e8c8e52743bb38b83b17d9429141d494b8041ca6d616a6c77cebae9cd8b7.
//
// Solidity: event AddressSet(bytes32 indexed id, address indexed oldAddress, address indexed newAddress)
func (_PoolAddressesProvider *PoolAddressesProviderFilterer) ParseAddressSet(log types.Log) (*PoolAddressesProviderAddressSet, error) {
	event := new(PoolAddressesProviderAddressSet)
	if err := _PoolAddressesProvider.contract.UnpackLog(event, "AddressSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolAddressesProviderAddressSetAsProxyIterator is returned from FilterAddressSetAsProxy and is used to iterate over the raw logs and unpacked data for AddressSetAsProxy events raised by the PoolAddressesProvider contract.
type PoolAddressesProviderAddressSetAsProxyIterator struct {
	Event *PoolAddressesProviderAddressSetAsProxy // Event containing the contract specifics and raw log

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
func (it *PoolAddressesProviderAddressSetAsProxyIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolAddressesProviderAddressSetAsProxy)
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
		it.Event = new(PoolAddressesProviderAddressSetAsProxy)
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
func (it *PoolAddressesProviderAddressSetAsProxyIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolAddressesProviderAddressSetAsProxyIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolAddressesProviderAddressSetAsProxy represents a AddressSetAsProxy event raised by the PoolAddressesProvider contract.
type PoolAddressesProviderAddressSetAsProxy struct {
	Id                       [32]byte
	ProxyAddress             common.Address
	OldImplementationAddress common.Address
	NewImplementationAddress common.Address
	Raw                      types.Log // Blockchain specific contextual infos
}

// FilterAddressSetAsProxy is a free log retrieval operation binding the contract event 0x3bbd45b5429b385e3fb37ad5cd1cd1435a3c8ec32196c7937597365a3fd3e99c.
//
// Solidity: event AddressSetAsProxy(bytes32 indexed id, address indexed proxyAddress, address oldImplementationAddress, address indexed newImplementationAddress)
func (_PoolAddressesProvider *PoolAddressesProviderFilterer) FilterAddressSetAsProxy(opts *bind.FilterOpts, id [][32]byte, proxyAddress []common.Address, newImplementationAddress []common.Address) (*PoolAddressesProviderAddressSetAsProxyIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var proxyAddressRule []interface{}
	for _, proxyAddressItem := range proxyAddress {
		proxyAddressRule = append(proxyAddressRule, proxyAddressItem)
	}

	var newImplementationAddressRule []interface{}
	for _, newImplementationAddressItem := range newImplementationAddress {
		newImplementationAddressRule = append(newImplementationAddressRule, newImplementationAddressItem)
	}

	logs, sub, err := _PoolAddressesProvider.contract.FilterLogs(opts, "AddressSetAsProxy", idRule, proxyAddressRule, newImplementationAddressRule)
	if err != nil {
		return nil, err
	}
	return &PoolAddressesProviderAddressSetAsProxyIterator{contract: _PoolAddressesProvider.contract, event: "AddressSetAsProxy", logs: logs, sub: sub}, nil
}

// WatchAddressSetAsProxy is a free log subscription operation binding the contract event 0x3bbd45b5429b385e3fb37ad5cd1cd1435a3c8ec32196c7937597365a3fd3e99c.
//
// Solidity: event AddressSetAsProxy(bytes32 indexed id, address indexed proxyAddress, address oldImplementationAddress, address indexed newImplementationAddress)
func (_PoolAddressesProvider *PoolAddressesProviderFilterer) WatchAddressSetAsProxy(opts *bind.WatchOpts, sink chan<- *PoolAddressesProviderAddressSetAsProxy, id [][32]byte, proxyAddress []common.Address, newImplementationAddress []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var proxyAddressRule []interface{}
	for _, proxyAddressItem := range proxyAddress {
		proxyAddressRule = append(proxyAddressRule, proxyAddressItem)
	}

	var newImplementationAddressRule []interface{}
	for _, newImplementationAddressItem := range newImplementationAddress {
		newImplementationAddressRule = append(newImplementationAddressRule, newImplementationAddressItem)
	}

	logs, sub, err := _PoolAddressesProvider.contract.WatchLogs(opts, "AddressSetAsProxy", idRule, proxyAddressRule, newImplementationAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolAddressesProviderAddressSetAsProxy)
				if err := _PoolAddressesProvider.contract.UnpackLog(event, "AddressSetAsProxy", log); err != nil {
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

// ParseAddressSetAsProxy is a log parse operation binding the contract event 0x3bbd45b5429b385e3fb37ad5cd1cd1435a3c8ec32196c7937597365a3fd3e99c.
//
// Solidity: event AddressSetAsProxy(bytes32 indexed id, address indexed proxyAddress, address oldImplementationAddress, address indexed newImplementationAddress)
func (_PoolAddressesProvider *PoolAddressesProviderFilterer) ParseAddressSetAsProxy(log types.Log) (*PoolAddressesProviderAddressSetAsProxy, error) {
	event := new(PoolAddressesProviderAddressSetAsProxy)
	if err := _PoolAddressesProvider.contract.UnpackLog(event, "AddressSetAsProxy", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolAddressesProviderMarketIdSetIterator is returned from FilterMarketIdSet and is used to iterate over the raw logs and unpacked data for MarketIdSet events raised by the PoolAddressesProvider contract.
type PoolAddressesProviderMarketIdSetIterator struct {
	Event *PoolAddressesProviderMarketIdSet // Event containing the contract specifics and raw log

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
func (it *PoolAddressesProviderMarketIdSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolAddressesProviderMarketIdSet)
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
		it.Event = new(PoolAddressesProviderMarketIdSet)
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
func (it *PoolAddressesProviderMarketIdSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolAddressesProviderMarketIdSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolAddressesProviderMarketIdSet represents a MarketIdSet event raised by the PoolAddressesProvider contract.
type PoolAddressesProviderMarketIdSet struct {
	OldMarketId common.Hash
	NewMarketId common.Hash
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterMarketIdSet is a free log retrieval operation binding the contract event 0xe685c8cdecc6030c45030fd54778812cb84ed8e4467c38294403d68ba7860823.
//
// Solidity: event MarketIdSet(string indexed oldMarketId, string indexed newMarketId)
func (_PoolAddressesProvider *PoolAddressesProviderFilterer) FilterMarketIdSet(opts *bind.FilterOpts, oldMarketId []string, newMarketId []string) (*PoolAddressesProviderMarketIdSetIterator, error) {

	var oldMarketIdRule []interface{}
	for _, oldMarketIdItem := range oldMarketId {
		oldMarketIdRule = append(oldMarketIdRule, oldMarketIdItem)
	}
	var newMarketIdRule []interface{}
	for _, newMarketIdItem := range newMarketId {
		newMarketIdRule = append(newMarketIdRule, newMarketIdItem)
	}

	logs, sub, err := _PoolAddressesProvider.contract.FilterLogs(opts, "MarketIdSet", oldMarketIdRule, newMarketIdRule)
	if err != nil {
		return nil, err
	}
	return &PoolAddressesProviderMarketIdSetIterator{contract: _PoolAddressesProvider.contract, event: "MarketIdSet", logs: logs, sub: sub}, nil
}

// WatchMarketIdSet is a free log subscription operation binding the contract event 0xe685c8cdecc6030c45030fd54778812cb84ed8e4467c38294403d68ba7860823.
//
// Solidity: event MarketIdSet(string indexed oldMarketId, string indexed newMarketId)
func (_PoolAddressesProvider *PoolAddressesProviderFilterer) WatchMarketIdSet(opts *bind.WatchOpts, sink chan<- *PoolAddressesProviderMarketIdSet, oldMarketId []string, newMarketId []string) (event.Subscription, error) {

	var oldMarketIdRule []interface{}
	for _, oldMarketIdItem := range oldMarketId {
		oldMarketIdRule = append(oldMarketIdRule, oldMarketIdItem)
	}
	var newMarketIdRule []interface{}
	for _, newMarketIdItem := range newMarketId {
		newMarketIdRule = append(newMarketIdRule, newMarketIdItem)
	}

	logs, sub, err := _PoolAddressesProvider.contract.WatchLogs(opts, "MarketIdSet", oldMarketIdRule, newMarketIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolAddressesProviderMarketIdSet)
				if err := _PoolAddressesProvider.contract.UnpackLog(event, "MarketIdSet", log); err != nil {
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

// ParseMarketIdSet is a log parse operation binding the contract event 0xe685c8cdecc6030c45030fd54778812cb84ed8e4467c38294403d68ba7860823.
//
// Solidity: event MarketIdSet(string indexed oldMarketId, string indexed newMarketId)
func (_PoolAddressesProvider *PoolAddressesProviderFilterer) ParseMarketIdSet(log types.Log) (*PoolAddressesProviderMarketIdSet, error) {
	event := new(PoolAddressesProviderMarketIdSet)
	if err := _PoolAddressesProvider.contract.UnpackLog(event, "MarketIdSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolAddressesProviderOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the PoolAddressesProvider contract.
type PoolAddressesProviderOwnershipTransferredIterator struct {
	Event *PoolAddressesProviderOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *PoolAddressesProviderOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolAddressesProviderOwnershipTransferred)
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
		it.Event = new(PoolAddressesProviderOwnershipTransferred)
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
func (it *PoolAddressesProviderOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolAddressesProviderOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolAddressesProviderOwnershipTransferred represents a OwnershipTransferred event raised by the PoolAddressesProvider contract.
type PoolAddressesProviderOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_PoolAddressesProvider *PoolAddressesProviderFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*PoolAddressesProviderOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _PoolAddressesProvider.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &PoolAddressesProviderOwnershipTransferredIterator{contract: _PoolAddressesProvider.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_PoolAddressesProvider *PoolAddressesProviderFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PoolAddressesProviderOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _PoolAddressesProvider.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolAddressesProviderOwnershipTransferred)
				if err := _PoolAddressesProvider.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_PoolAddressesProvider *PoolAddressesProviderFilterer) ParseOwnershipTransferred(log types.Log) (*PoolAddressesProviderOwnershipTransferred, error) {
	event := new(PoolAddressesProviderOwnershipTransferred)
	if err := _PoolAddressesProvider.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolAddressesProviderPoolConfiguratorUpdatedIterator is returned from FilterPoolConfiguratorUpdated and is used to iterate over the raw logs and unpacked data for PoolConfiguratorUpdated events raised by the PoolAddressesProvider contract.
type PoolAddressesProviderPoolConfiguratorUpdatedIterator struct {
	Event *PoolAddressesProviderPoolConfiguratorUpdated // Event containing the contract specifics and raw log

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
func (it *PoolAddressesProviderPoolConfiguratorUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolAddressesProviderPoolConfiguratorUpdated)
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
		it.Event = new(PoolAddressesProviderPoolConfiguratorUpdated)
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
func (it *PoolAddressesProviderPoolConfiguratorUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolAddressesProviderPoolConfiguratorUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolAddressesProviderPoolConfiguratorUpdated represents a PoolConfiguratorUpdated event raised by the PoolAddressesProvider contract.
type PoolAddressesProviderPoolConfiguratorUpdated struct {
	OldAddress common.Address
	NewAddress common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterPoolConfiguratorUpdated is a free log retrieval operation binding the contract event 0x8932892569eba59c8382a089d9b732d1f49272878775235761a2a6b0309cd465.
//
// Solidity: event PoolConfiguratorUpdated(address indexed oldAddress, address indexed newAddress)
func (_PoolAddressesProvider *PoolAddressesProviderFilterer) FilterPoolConfiguratorUpdated(opts *bind.FilterOpts, oldAddress []common.Address, newAddress []common.Address) (*PoolAddressesProviderPoolConfiguratorUpdatedIterator, error) {

	var oldAddressRule []interface{}
	for _, oldAddressItem := range oldAddress {
		oldAddressRule = append(oldAddressRule, oldAddressItem)
	}
	var newAddressRule []interface{}
	for _, newAddressItem := range newAddress {
		newAddressRule = append(newAddressRule, newAddressItem)
	}

	logs, sub, err := _PoolAddressesProvider.contract.FilterLogs(opts, "PoolConfiguratorUpdated", oldAddressRule, newAddressRule)
	if err != nil {
		return nil, err
	}
	return &PoolAddressesProviderPoolConfiguratorUpdatedIterator{contract: _PoolAddressesProvider.contract, event: "PoolConfiguratorUpdated", logs: logs, sub: sub}, nil
}

// WatchPoolConfiguratorUpdated is a free log subscription operation binding the contract event 0x8932892569eba59c8382a089d9b732d1f49272878775235761a2a6b0309cd465.
//
// Solidity: event PoolConfiguratorUpdated(address indexed oldAddress, address indexed newAddress)
func (_PoolAddressesProvider *PoolAddressesProviderFilterer) WatchPoolConfiguratorUpdated(opts *bind.WatchOpts, sink chan<- *PoolAddressesProviderPoolConfiguratorUpdated, oldAddress []common.Address, newAddress []common.Address) (event.Subscription, error) {

	var oldAddressRule []interface{}
	for _, oldAddressItem := range oldAddress {
		oldAddressRule = append(oldAddressRule, oldAddressItem)
	}
	var newAddressRule []interface{}
	for _, newAddressItem := range newAddress {
		newAddressRule = append(newAddressRule, newAddressItem)
	}

	logs, sub, err := _PoolAddressesProvider.contract.WatchLogs(opts, "PoolConfiguratorUpdated", oldAddressRule, newAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolAddressesProviderPoolConfiguratorUpdated)
				if err := _PoolAddressesProvider.contract.UnpackLog(event, "PoolConfiguratorUpdated", log); err != nil {
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

// ParsePoolConfiguratorUpdated is a log parse operation binding the contract event 0x8932892569eba59c8382a089d9b732d1f49272878775235761a2a6b0309cd465.
//
// Solidity: event PoolConfiguratorUpdated(address indexed oldAddress, address indexed newAddress)
func (_PoolAddressesProvider *PoolAddressesProviderFilterer) ParsePoolConfiguratorUpdated(log types.Log) (*PoolAddressesProviderPoolConfiguratorUpdated, error) {
	event := new(PoolAddressesProviderPoolConfiguratorUpdated)
	if err := _PoolAddressesProvider.contract.UnpackLog(event, "PoolConfiguratorUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolAddressesProviderPoolDataProviderUpdatedIterator is returned from FilterPoolDataProviderUpdated and is used to iterate over the raw logs and unpacked data for PoolDataProviderUpdated events raised by the PoolAddressesProvider contract.
type PoolAddressesProviderPoolDataProviderUpdatedIterator struct {
	Event *PoolAddressesProviderPoolDataProviderUpdated // Event containing the contract specifics and raw log

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
func (it *PoolAddressesProviderPoolDataProviderUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolAddressesProviderPoolDataProviderUpdated)
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
		it.Event = new(PoolAddressesProviderPoolDataProviderUpdated)
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
func (it *PoolAddressesProviderPoolDataProviderUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolAddressesProviderPoolDataProviderUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolAddressesProviderPoolDataProviderUpdated represents a PoolDataProviderUpdated event raised by the PoolAddressesProvider contract.
type PoolAddressesProviderPoolDataProviderUpdated struct {
	OldAddress common.Address
	NewAddress common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterPoolDataProviderUpdated is a free log retrieval operation binding the contract event 0xc853974cfbf81487a14a23565917bee63f527853bcb5fa54f2ae1cdf8a38356d.
//
// Solidity: event PoolDataProviderUpdated(address indexed oldAddress, address indexed newAddress)
func (_PoolAddressesProvider *PoolAddressesProviderFilterer) FilterPoolDataProviderUpdated(opts *bind.FilterOpts, oldAddress []common.Address, newAddress []common.Address) (*PoolAddressesProviderPoolDataProviderUpdatedIterator, error) {

	var oldAddressRule []interface{}
	for _, oldAddressItem := range oldAddress {
		oldAddressRule = append(oldAddressRule, oldAddressItem)
	}
	var newAddressRule []interface{}
	for _, newAddressItem := range newAddress {
		newAddressRule = append(newAddressRule, newAddressItem)
	}

	logs, sub, err := _PoolAddressesProvider.contract.FilterLogs(opts, "PoolDataProviderUpdated", oldAddressRule, newAddressRule)
	if err != nil {
		return nil, err
	}
	return &PoolAddressesProviderPoolDataProviderUpdatedIterator{contract: _PoolAddressesProvider.contract, event: "PoolDataProviderUpdated", logs: logs, sub: sub}, nil
}

// WatchPoolDataProviderUpdated is a free log subscription operation binding the contract event 0xc853974cfbf81487a14a23565917bee63f527853bcb5fa54f2ae1cdf8a38356d.
//
// Solidity: event PoolDataProviderUpdated(address indexed oldAddress, address indexed newAddress)
func (_PoolAddressesProvider *PoolAddressesProviderFilterer) WatchPoolDataProviderUpdated(opts *bind.WatchOpts, sink chan<- *PoolAddressesProviderPoolDataProviderUpdated, oldAddress []common.Address, newAddress []common.Address) (event.Subscription, error) {

	var oldAddressRule []interface{}
	for _, oldAddressItem := range oldAddress {
		oldAddressRule = append(oldAddressRule, oldAddressItem)
	}
	var newAddressRule []interface{}
	for _, newAddressItem := range newAddress {
		newAddressRule = append(newAddressRule, newAddressItem)
	}

	logs, sub, err := _PoolAddressesProvider.contract.WatchLogs(opts, "PoolDataProviderUpdated", oldAddressRule, newAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolAddressesProviderPoolDataProviderUpdated)
				if err := _PoolAddressesProvider.contract.UnpackLog(event, "PoolDataProviderUpdated", log); err != nil {
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

// ParsePoolDataProviderUpdated is a log parse operation binding the contract event 0xc853974cfbf81487a14a23565917bee63f527853bcb5fa54f2ae1cdf8a38356d.
//
// Solidity: event PoolDataProviderUpdated(address indexed oldAddress, address indexed newAddress)
func (_PoolAddressesProvider *PoolAddressesProviderFilterer) ParsePoolDataProviderUpdated(log types.Log) (*PoolAddressesProviderPoolDataProviderUpdated, error) {
	event := new(PoolAddressesProviderPoolDataProviderUpdated)
	if err := _PoolAddressesProvider.contract.UnpackLog(event, "PoolDataProviderUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolAddressesProviderPoolUpdatedIterator is returned from FilterPoolUpdated and is used to iterate over the raw logs and unpacked data for PoolUpdated events raised by the PoolAddressesProvider contract.
type PoolAddressesProviderPoolUpdatedIterator struct {
	Event *PoolAddressesProviderPoolUpdated // Event containing the contract specifics and raw log

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
func (it *PoolAddressesProviderPoolUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolAddressesProviderPoolUpdated)
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
		it.Event = new(PoolAddressesProviderPoolUpdated)
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
func (it *PoolAddressesProviderPoolUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolAddressesProviderPoolUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolAddressesProviderPoolUpdated represents a PoolUpdated event raised by the PoolAddressesProvider contract.
type PoolAddressesProviderPoolUpdated struct {
	OldAddress common.Address
	NewAddress common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterPoolUpdated is a free log retrieval operation binding the contract event 0x90affc163f1a2dfedcd36aa02ed992eeeba8100a4014f0b4cdc20ea265a66627.
//
// Solidity: event PoolUpdated(address indexed oldAddress, address indexed newAddress)
func (_PoolAddressesProvider *PoolAddressesProviderFilterer) FilterPoolUpdated(opts *bind.FilterOpts, oldAddress []common.Address, newAddress []common.Address) (*PoolAddressesProviderPoolUpdatedIterator, error) {

	var oldAddressRule []interface{}
	for _, oldAddressItem := range oldAddress {
		oldAddressRule = append(oldAddressRule, oldAddressItem)
	}
	var newAddressRule []interface{}
	for _, newAddressItem := range newAddress {
		newAddressRule = append(newAddressRule, newAddressItem)
	}

	logs, sub, err := _PoolAddressesProvider.contract.FilterLogs(opts, "PoolUpdated", oldAddressRule, newAddressRule)
	if err != nil {
		return nil, err
	}
	return &PoolAddressesProviderPoolUpdatedIterator{contract: _PoolAddressesProvider.contract, event: "PoolUpdated", logs: logs, sub: sub}, nil
}

// WatchPoolUpdated is a free log subscription operation binding the contract event 0x90affc163f1a2dfedcd36aa02ed992eeeba8100a4014f0b4cdc20ea265a66627.
//
// Solidity: event PoolUpdated(address indexed oldAddress, address indexed newAddress)
func (_PoolAddressesProvider *PoolAddressesProviderFilterer) WatchPoolUpdated(opts *bind.WatchOpts, sink chan<- *PoolAddressesProviderPoolUpdated, oldAddress []common.Address, newAddress []common.Address) (event.Subscription, error) {

	var oldAddressRule []interface{}
	for _, oldAddressItem := range oldAddress {
		oldAddressRule = append(oldAddressRule, oldAddressItem)
	}
	var newAddressRule []interface{}
	for _, newAddressItem := range newAddress {
		newAddressRule = append(newAddressRule, newAddressItem)
	}

	logs, sub, err := _PoolAddressesProvider.contract.WatchLogs(opts, "PoolUpdated", oldAddressRule, newAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolAddressesProviderPoolUpdated)
				if err := _PoolAddressesProvider.contract.UnpackLog(event, "PoolUpdated", log); err != nil {
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

// ParsePoolUpdated is a log parse operation binding the contract event 0x90affc163f1a2dfedcd36aa02ed992eeeba8100a4014f0b4cdc20ea265a66627.
//
// Solidity: event PoolUpdated(address indexed oldAddress, address indexed newAddress)
func (_PoolAddressesProvider *PoolAddressesProviderFilterer) ParsePoolUpdated(log types.Log) (*PoolAddressesProviderPoolUpdated, error) {
	event := new(PoolAddressesProviderPoolUpdated)
	if err := _PoolAddressesProvider.contract.UnpackLog(event, "PoolUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolAddressesProviderPriceOracleSentinelUpdatedIterator is returned from FilterPriceOracleSentinelUpdated and is used to iterate over the raw logs and unpacked data for PriceOracleSentinelUpdated events raised by the PoolAddressesProvider contract.
type PoolAddressesProviderPriceOracleSentinelUpdatedIterator struct {
	Event *PoolAddressesProviderPriceOracleSentinelUpdated // Event containing the contract specifics and raw log

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
func (it *PoolAddressesProviderPriceOracleSentinelUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolAddressesProviderPriceOracleSentinelUpdated)
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
		it.Event = new(PoolAddressesProviderPriceOracleSentinelUpdated)
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
func (it *PoolAddressesProviderPriceOracleSentinelUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolAddressesProviderPriceOracleSentinelUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolAddressesProviderPriceOracleSentinelUpdated represents a PriceOracleSentinelUpdated event raised by the PoolAddressesProvider contract.
type PoolAddressesProviderPriceOracleSentinelUpdated struct {
	OldAddress common.Address
	NewAddress common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterPriceOracleSentinelUpdated is a free log retrieval operation binding the contract event 0x5326514eeca90494a14bedabcff812a0e683029ee85d1e23824d44fd14cd6ae7.
//
// Solidity: event PriceOracleSentinelUpdated(address indexed oldAddress, address indexed newAddress)
func (_PoolAddressesProvider *PoolAddressesProviderFilterer) FilterPriceOracleSentinelUpdated(opts *bind.FilterOpts, oldAddress []common.Address, newAddress []common.Address) (*PoolAddressesProviderPriceOracleSentinelUpdatedIterator, error) {

	var oldAddressRule []interface{}
	for _, oldAddressItem := range oldAddress {
		oldAddressRule = append(oldAddressRule, oldAddressItem)
	}
	var newAddressRule []interface{}
	for _, newAddressItem := range newAddress {
		newAddressRule = append(newAddressRule, newAddressItem)
	}

	logs, sub, err := _PoolAddressesProvider.contract.FilterLogs(opts, "PriceOracleSentinelUpdated", oldAddressRule, newAddressRule)
	if err != nil {
		return nil, err
	}
	return &PoolAddressesProviderPriceOracleSentinelUpdatedIterator{contract: _PoolAddressesProvider.contract, event: "PriceOracleSentinelUpdated", logs: logs, sub: sub}, nil
}

// WatchPriceOracleSentinelUpdated is a free log subscription operation binding the contract event 0x5326514eeca90494a14bedabcff812a0e683029ee85d1e23824d44fd14cd6ae7.
//
// Solidity: event PriceOracleSentinelUpdated(address indexed oldAddress, address indexed newAddress)
func (_PoolAddressesProvider *PoolAddressesProviderFilterer) WatchPriceOracleSentinelUpdated(opts *bind.WatchOpts, sink chan<- *PoolAddressesProviderPriceOracleSentinelUpdated, oldAddress []common.Address, newAddress []common.Address) (event.Subscription, error) {

	var oldAddressRule []interface{}
	for _, oldAddressItem := range oldAddress {
		oldAddressRule = append(oldAddressRule, oldAddressItem)
	}
	var newAddressRule []interface{}
	for _, newAddressItem := range newAddress {
		newAddressRule = append(newAddressRule, newAddressItem)
	}

	logs, sub, err := _PoolAddressesProvider.contract.WatchLogs(opts, "PriceOracleSentinelUpdated", oldAddressRule, newAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolAddressesProviderPriceOracleSentinelUpdated)
				if err := _PoolAddressesProvider.contract.UnpackLog(event, "PriceOracleSentinelUpdated", log); err != nil {
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

// ParsePriceOracleSentinelUpdated is a log parse operation binding the contract event 0x5326514eeca90494a14bedabcff812a0e683029ee85d1e23824d44fd14cd6ae7.
//
// Solidity: event PriceOracleSentinelUpdated(address indexed oldAddress, address indexed newAddress)
func (_PoolAddressesProvider *PoolAddressesProviderFilterer) ParsePriceOracleSentinelUpdated(log types.Log) (*PoolAddressesProviderPriceOracleSentinelUpdated, error) {
	event := new(PoolAddressesProviderPriceOracleSentinelUpdated)
	if err := _PoolAddressesProvider.contract.UnpackLog(event, "PriceOracleSentinelUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolAddressesProviderPriceOracleUpdatedIterator is returned from FilterPriceOracleUpdated and is used to iterate over the raw logs and unpacked data for PriceOracleUpdated events raised by the PoolAddressesProvider contract.
type PoolAddressesProviderPriceOracleUpdatedIterator struct {
	Event *PoolAddressesProviderPriceOracleUpdated // Event containing the contract specifics and raw log

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
func (it *PoolAddressesProviderPriceOracleUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolAddressesProviderPriceOracleUpdated)
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
		it.Event = new(PoolAddressesProviderPriceOracleUpdated)
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
func (it *PoolAddressesProviderPriceOracleUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolAddressesProviderPriceOracleUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolAddressesProviderPriceOracleUpdated represents a PriceOracleUpdated event raised by the PoolAddressesProvider contract.
type PoolAddressesProviderPriceOracleUpdated struct {
	OldAddress common.Address
	NewAddress common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterPriceOracleUpdated is a free log retrieval operation binding the contract event 0x56b5f80d8cac1479698aa7d01605fd6111e90b15fc4d2b377417f46034876cbd.
//
// Solidity: event PriceOracleUpdated(address indexed oldAddress, address indexed newAddress)
func (_PoolAddressesProvider *PoolAddressesProviderFilterer) FilterPriceOracleUpdated(opts *bind.FilterOpts, oldAddress []common.Address, newAddress []common.Address) (*PoolAddressesProviderPriceOracleUpdatedIterator, error) {

	var oldAddressRule []interface{}
	for _, oldAddressItem := range oldAddress {
		oldAddressRule = append(oldAddressRule, oldAddressItem)
	}
	var newAddressRule []interface{}
	for _, newAddressItem := range newAddress {
		newAddressRule = append(newAddressRule, newAddressItem)
	}

	logs, sub, err := _PoolAddressesProvider.contract.FilterLogs(opts, "PriceOracleUpdated", oldAddressRule, newAddressRule)
	if err != nil {
		return nil, err
	}
	return &PoolAddressesProviderPriceOracleUpdatedIterator{contract: _PoolAddressesProvider.contract, event: "PriceOracleUpdated", logs: logs, sub: sub}, nil
}

// WatchPriceOracleUpdated is a free log subscription operation binding the contract event 0x56b5f80d8cac1479698aa7d01605fd6111e90b15fc4d2b377417f46034876cbd.
//
// Solidity: event PriceOracleUpdated(address indexed oldAddress, address indexed newAddress)
func (_PoolAddressesProvider *PoolAddressesProviderFilterer) WatchPriceOracleUpdated(opts *bind.WatchOpts, sink chan<- *PoolAddressesProviderPriceOracleUpdated, oldAddress []common.Address, newAddress []common.Address) (event.Subscription, error) {

	var oldAddressRule []interface{}
	for _, oldAddressItem := range oldAddress {
		oldAddressRule = append(oldAddressRule, oldAddressItem)
	}
	var newAddressRule []interface{}
	for _, newAddressItem := range newAddress {
		newAddressRule = append(newAddressRule, newAddressItem)
	}

	logs, sub, err := _PoolAddressesProvider.contract.WatchLogs(opts, "PriceOracleUpdated", oldAddressRule, newAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolAddressesProviderPriceOracleUpdated)
				if err := _PoolAddressesProvider.contract.UnpackLog(event, "PriceOracleUpdated", log); err != nil {
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

// ParsePriceOracleUpdated is a log parse operation binding the contract event 0x56b5f80d8cac1479698aa7d01605fd6111e90b15fc4d2b377417f46034876cbd.
//
// Solidity: event PriceOracleUpdated(address indexed oldAddress, address indexed newAddress)
func (_PoolAddressesProvider *PoolAddressesProviderFilterer) ParsePriceOracleUpdated(log types.Log) (*PoolAddressesProviderPriceOracleUpdated, error) {
	event := new(PoolAddressesProviderPriceOracleUpdated)
	if err := _PoolAddressesProvider.contract.UnpackLog(event, "PriceOracleUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolAddressesProviderProxyCreatedIterator is returned from FilterProxyCreated and is used to iterate over the raw logs and unpacked data for ProxyCreated events raised by the PoolAddressesProvider contract.
type PoolAddressesProviderProxyCreatedIterator struct {
	Event *PoolAddressesProviderProxyCreated // Event containing the contract specifics and raw log

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
func (it *PoolAddressesProviderProxyCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolAddressesProviderProxyCreated)
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
		it.Event = new(PoolAddressesProviderProxyCreated)
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
func (it *PoolAddressesProviderProxyCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolAddressesProviderProxyCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolAddressesProviderProxyCreated represents a ProxyCreated event raised by the PoolAddressesProvider contract.
type PoolAddressesProviderProxyCreated struct {
	Id                    [32]byte
	ProxyAddress          common.Address
	ImplementationAddress common.Address
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterProxyCreated is a free log retrieval operation binding the contract event 0x4a465a9bd819d9662563c1e11ae958f8109e437e7f4bf1c6ef0b9a7b3f35d478.
//
// Solidity: event ProxyCreated(bytes32 indexed id, address indexed proxyAddress, address indexed implementationAddress)
func (_PoolAddressesProvider *PoolAddressesProviderFilterer) FilterProxyCreated(opts *bind.FilterOpts, id [][32]byte, proxyAddress []common.Address, implementationAddress []common.Address) (*PoolAddressesProviderProxyCreatedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var proxyAddressRule []interface{}
	for _, proxyAddressItem := range proxyAddress {
		proxyAddressRule = append(proxyAddressRule, proxyAddressItem)
	}
	var implementationAddressRule []interface{}
	for _, implementationAddressItem := range implementationAddress {
		implementationAddressRule = append(implementationAddressRule, implementationAddressItem)
	}

	logs, sub, err := _PoolAddressesProvider.contract.FilterLogs(opts, "ProxyCreated", idRule, proxyAddressRule, implementationAddressRule)
	if err != nil {
		return nil, err
	}
	return &PoolAddressesProviderProxyCreatedIterator{contract: _PoolAddressesProvider.contract, event: "ProxyCreated", logs: logs, sub: sub}, nil
}

// WatchProxyCreated is a free log subscription operation binding the contract event 0x4a465a9bd819d9662563c1e11ae958f8109e437e7f4bf1c6ef0b9a7b3f35d478.
//
// Solidity: event ProxyCreated(bytes32 indexed id, address indexed proxyAddress, address indexed implementationAddress)
func (_PoolAddressesProvider *PoolAddressesProviderFilterer) WatchProxyCreated(opts *bind.WatchOpts, sink chan<- *PoolAddressesProviderProxyCreated, id [][32]byte, proxyAddress []common.Address, implementationAddress []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var proxyAddressRule []interface{}
	for _, proxyAddressItem := range proxyAddress {
		proxyAddressRule = append(proxyAddressRule, proxyAddressItem)
	}
	var implementationAddressRule []interface{}
	for _, implementationAddressItem := range implementationAddress {
		implementationAddressRule = append(implementationAddressRule, implementationAddressItem)
	}

	logs, sub, err := _PoolAddressesProvider.contract.WatchLogs(opts, "ProxyCreated", idRule, proxyAddressRule, implementationAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolAddressesProviderProxyCreated)
				if err := _PoolAddressesProvider.contract.UnpackLog(event, "ProxyCreated", log); err != nil {
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

// ParseProxyCreated is a log parse operation binding the contract event 0x4a465a9bd819d9662563c1e11ae958f8109e437e7f4bf1c6ef0b9a7b3f35d478.
//
// Solidity: event ProxyCreated(bytes32 indexed id, address indexed proxyAddress, address indexed implementationAddress)
func (_PoolAddressesProvider *PoolAddressesProviderFilterer) ParseProxyCreated(log types.Log) (*PoolAddressesProviderProxyCreated, error) {
	event := new(PoolAddressesProviderProxyCreated)
	if err := _PoolAddressesProvider.contract.UnpackLog(event, "ProxyCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
