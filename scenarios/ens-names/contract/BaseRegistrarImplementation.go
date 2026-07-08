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

// BaseRegistrarImplementationMetaData contains all meta data concerning the BaseRegistrarImplementation contract.
var BaseRegistrarImplementationMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractENS\",\"name\":\"_ens\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_baseNode\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"controller\",\"type\":\"address\"}],\"name\":\"ControllerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"controller\",\"type\":\"address\"}],\"name\":\"ControllerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"expires\",\"type\":\"uint256\"}],\"name\":\"NameMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"expires\",\"type\":\"uint256\"}],\"name\":\"NameRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"expires\",\"type\":\"uint256\"}],\"name\":\"NameRenewed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"GRACE_PERIOD\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"controller\",\"type\":\"address\"}],\"name\":\"addController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"available\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"baseNode\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"controllers\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ens\",\"outputs\":[{\"internalType\":\"contractENS\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"nameExpires\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"reclaim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"duration\",\"type\":\"uint256\"}],\"name\":\"register\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"duration\",\"type\":\"uint256\"}],\"name\":\"registerOnly\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"controller\",\"type\":\"address\"}],\"name\":\"removeController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"duration\",\"type\":\"uint256\"}],\"name\":\"renew\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"resolver\",\"type\":\"address\"}],\"name\":\"setResolver\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceID\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506040516129d53803806129d583398101604081905261002f916100fb565b604080516020808201835260008083528351918201909352828152909161005683826101d4565b50600161006382826101d4565b50505061007c6100776100a560201b60201c565b6100a9565b600880546001600160a01b0319166001600160a01b039390931692909217909155600955610292565b3390565b600680546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b6000806040838503121561010e57600080fd5b82516001600160a01b038116811461012557600080fd5b6020939093015192949293505050565b634e487b7160e01b600052604160045260246000fd5b600181811c9082168061015f57607f821691505b60208210810361017f57634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156101cf57806000526020600020601f840160051c810160208510156101ac5750805b601f840160051c820191505b818110156101cc57600081556001016101b8565b50505b505050565b81516001600160401b038111156101ed576101ed610135565b610201816101fb845461014b565b84610185565b6020601f821160018114610235576000831561021d5750848201515b600019600385901b1c1916600184901b1784556101cc565b600084815260208120601f198516915b828110156102655787850151825560209485019460019092019101610245565b50848210156102835786840151600019600387901b60f8161c191681555b50505050600190811b01905550565b612734806102a16000396000f3fe608060405234801561001057600080fd5b50600436106101cf5760003560e01c806395d89b4111610104578063c87b56dd116100a2578063e985e9c511610071578063e985e9c514610407578063f2fde38b14610450578063f6a74ed714610463578063fca247ac1461047657600080fd5b8063c87b56dd146103a8578063d6e4fa86146103bb578063da8c229e146103db578063ddf7fcb0146103fe57600080fd5b8063a7fc7a07116100de578063a7fc7a0714610365578063b88d4fde14610378578063c1a287e21461038b578063c475abff1461039557600080fd5b806395d89b411461033757806396e494e81461033f578063a22cb4651461035257600080fd5b80633f15457f116101715780636352211e1161014b5780636352211e146102eb57806370a08231146102fe578063715018a6146103115780638da5cb5b1461031957600080fd5b80633f15457f146102a557806342842e0e146102c55780634e543b26146102d857600080fd5b8063095ea7b3116101ad578063095ea7b3146102495780630e297b451461025e57806323b872dd1461027f57806328ed4f6c1461029257600080fd5b806301ffc9a7146101d457806306fdde03146101fc578063081812fc14610211575b600080fd5b6101e76101e2366004612203565b610489565b60405190151581526020015b60405180910390f35b61020461056e565b6040516101f3919061228e565b61022461021f3660046122a1565b610600565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101f3565b61025c6102573660046122dc565b610634565b005b61027161026c366004612308565b6107c5565b6040519081526020016101f3565b61025c61028d366004612340565b6107dc565b61025c6102a0366004612370565b61087d565b6008546102249073ffffffffffffffffffffffffffffffffffffffff1681565b61025c6102d3366004612340565b6109ef565b61025c6102e63660046123a0565b610a0a565b6102246102f93660046122a1565b610aa5565b61027161030c3660046123a0565b610ac8565b61025c610b96565b60065473ffffffffffffffffffffffffffffffffffffffff16610224565b610204610baa565b6101e761034d3660046122a1565b610bb9565b61025c6103603660046123bd565b610bdf565b61025c6103733660046123a0565b610bee565b61025c61038636600461241f565b610c6d565b6102716276a70081565b6102716103a3366004612544565b610d15565b6102046103b63660046122a1565b610ed9565b6102716103c93660046122a1565b60009081526007602052604090205490565b6101e76103e93660046123a0565b600a6020526000908152604090205460ff1681565b61027160095481565b6101e7610415366004612566565b73ffffffffffffffffffffffffffffffffffffffff918216600090815260056020908152604080832093909416825291909152205460ff1690565b61025c61045e3660046123a0565b610f4d565b61025c6104713660046123a0565b611004565b610271610484366004612308565b611080565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a700000000000000000000000000000000000000000000000000000000148061051c57507fffffffff0000000000000000000000000000000000000000000000000000000082167f80ac58cd00000000000000000000000000000000000000000000000000000000145b8061056857507fffffffff0000000000000000000000000000000000000000000000000000000082167f28ed4f6c00000000000000000000000000000000000000000000000000000000145b92915050565b60606000805461057d90612594565b80601f01602080910402602001604051908101604052809291908181526020018280546105a990612594565b80156105f65780601f106105cb576101008083540402835291602001916105f6565b820191906000526020600020905b8154815290600101906020018083116105d957829003601f168201915b5050505050905090565b600061060b8261108f565b5060009081526004602052604090205473ffffffffffffffffffffffffffffffffffffffff1690565b600061063f8261111a565b90508073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610701576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602160248201527f4552433732313a20617070726f76616c20746f2063757272656e74206f776e6560448201527f720000000000000000000000000000000000000000000000000000000000000060648201526084015b60405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff8216148061072a575061072a8133610415565b6107b6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603d60248201527f4552433732313a20617070726f76652063616c6c6572206973206e6f7420746f60448201527f6b656e206f776e6572206f7220617070726f76656420666f7220616c6c00000060648201526084016106f8565b6107c083836111a6565b505050565b60006107d48484846000611246565b949350505050565b6107e633826114c9565b610872576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602d60248201527f4552433732313a2063616c6c6572206973206e6f7420746f6b656e206f776e6560448201527f72206f7220617070726f7665640000000000000000000000000000000000000060648201526084016106f8565b6107c0838383611585565b6008546009546040517f02571be30000000000000000000000000000000000000000000000000000000081526004810191909152309173ffffffffffffffffffffffffffffffffffffffff16906302571be390602401602060405180830381865afa1580156108f0573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061091491906125e7565b73ffffffffffffffffffffffffffffffffffffffff161461093457600080fd5b61093e33836114c9565b61094757600080fd5b6008546009546040517f06ab592300000000000000000000000000000000000000000000000000000000815260048101919091526024810184905273ffffffffffffffffffffffffffffffffffffffff8381166044830152909116906306ab5923906064016020604051808303816000875af11580156109cb573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107c09190612604565b6107c083838360405180602001604052806000815250610c6d565b610a12611880565b6008546009546040517f1896f70a000000000000000000000000000000000000000000000000000000008152600481019190915273ffffffffffffffffffffffffffffffffffffffff838116602483015290911690631896f70a90604401600060405180830381600087803b158015610a8a57600080fd5b505af1158015610a9e573d6000803e3d6000fd5b5050505050565b6000818152600760205260408120544210610abf57600080fd5b6105688261111a565b600073ffffffffffffffffffffffffffffffffffffffff8216610b6d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602960248201527f4552433732313a2061646472657373207a65726f206973206e6f74206120766160448201527f6c6964206f776e6572000000000000000000000000000000000000000000000060648201526084016106f8565b5073ffffffffffffffffffffffffffffffffffffffff1660009081526003602052604090205490565b610b9e611880565b610ba86000611901565b565b60606001805461057d90612594565b6000818152600760205260408120544290610bd8906276a7009061261d565b1092915050565b610bea338383611978565b5050565b610bf6611880565b73ffffffffffffffffffffffffffffffffffffffff81166000818152600a602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055517f0a8bb31534c0ed46f380cb867bd5c803a189ced9a764e30b3a4991a9901d74749190a250565b610c7733836114c9565b610d03576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602d60248201527f4552433732313a2063616c6c6572206973206e6f7420746f6b656e206f776e6560448201527f72206f7220617070726f7665640000000000000000000000000000000000000060648201526084016106f8565b610d0f84848484611aa5565b50505050565b6008546009546040517f02571be30000000000000000000000000000000000000000000000000000000081526004810191909152600091309173ffffffffffffffffffffffffffffffffffffffff909116906302571be390602401602060405180830381865afa158015610d8d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610db191906125e7565b73ffffffffffffffffffffffffffffffffffffffff1614610dd157600080fd5b336000908152600a602052604090205460ff16610ded57600080fd5b6000838152600760205260409020544290610e0c906276a7009061261d565b1015610e1757600080fd5b610e246276a7008361261d565b6000848152600760205260409020546276a70090610e4390859061261d565b610e4d919061261d565b11610e5757600080fd5b60008381526007602052604081208054849290610e7590849061261d565b90915550506000838152600760205260409081902054905184917f9b87a00e30f1ac65d898f070f8a3488fe60517182d0a2098e1b4b93a54aa9bd691610ebd91815260200190565b60405180910390a2505060009081526007602052604090205490565b6060610ee48261108f565b6000610efb60408051602081019091526000815290565b90506000815111610f1b5760405180602001604052806000815250610f46565b80610f2584611b48565b604051602001610f36929190612657565b6040516020818303038152906040525b9392505050565b610f55611880565b73ffffffffffffffffffffffffffffffffffffffff8116610ff8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201527f646472657373000000000000000000000000000000000000000000000000000060648201526084016106f8565b61100181611901565b50565b61100c611880565b73ffffffffffffffffffffffffffffffffffffffff81166000818152600a602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055517f33d83959be2573f5453b12eb9d43b3499bc57d96bd2f067ba44803c859e811139190a250565b60006107d48484846001611246565b60008181526002602052604090205473ffffffffffffffffffffffffffffffffffffffff16611001576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f4552433732313a20696e76616c696420746f6b656e204944000000000000000060448201526064016106f8565b60008181526002602052604081205473ffffffffffffffffffffffffffffffffffffffff1680610568576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f4552433732313a20696e76616c696420746f6b656e204944000000000000000060448201526064016106f8565b600081815260046020526040902080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff841690811790915581906112008261111a565b73ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92560405160405180910390a45050565b6008546009546040517f02571be30000000000000000000000000000000000000000000000000000000081526004810191909152600091309173ffffffffffffffffffffffffffffffffffffffff909116906302571be390602401602060405180830381865afa1580156112be573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112e291906125e7565b73ffffffffffffffffffffffffffffffffffffffff161461130257600080fd5b336000908152600a602052604090205460ff1661131e57600080fd5b61132785610bb9565b61133057600080fd5b61133d6276a7004261261d565b6276a70061134b854261261d565b611355919061261d565b1161135f57600080fd5b611369834261261d565b60008681526007602090815260408083209390935560029052205473ffffffffffffffffffffffffffffffffffffffff16156113a8576113a885611c06565b6113b28486611cde565b8115611462576008546009546040517f06ab592300000000000000000000000000000000000000000000000000000000815260048101919091526024810187905273ffffffffffffffffffffffffffffffffffffffff8681166044830152909116906306ab5923906064016020604051808303816000875af115801561143c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114609190612604565b505b73ffffffffffffffffffffffffffffffffffffffff8416857fb3d987963d01b2f68493b4bdb130988f157ea43070d4ad840fee0466ed9370d96114a5864261261d565b60405190815260200160405180910390a36114c0834261261d565b95945050505050565b6000806114d583610aa5565b90508073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff16148061154457508373ffffffffffffffffffffffffffffffffffffffff1661152c84610600565b73ffffffffffffffffffffffffffffffffffffffff16145b806107d4575073ffffffffffffffffffffffffffffffffffffffff80821660009081526005602090815260408083209388168352929052205460ff166107d4565b8273ffffffffffffffffffffffffffffffffffffffff166115a58261111a565b73ffffffffffffffffffffffffffffffffffffffff1614611648576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602560248201527f4552433732313a207472616e736665722066726f6d20696e636f72726563742060448201527f6f776e657200000000000000000000000000000000000000000000000000000060648201526084016106f8565b73ffffffffffffffffffffffffffffffffffffffff82166116ea576040517f08c379a0000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f4552433732313a207472616e7366657220746f20746865207a65726f2061646460448201527f726573730000000000000000000000000000000000000000000000000000000060648201526084016106f8565b8273ffffffffffffffffffffffffffffffffffffffff1661170a8261111a565b73ffffffffffffffffffffffffffffffffffffffff16146117ad576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602560248201527f4552433732313a207472616e736665722066726f6d20696e636f72726563742060448201527f6f776e657200000000000000000000000000000000000000000000000000000060648201526084016106f8565b600081815260046020908152604080832080547fffffffffffffffffffffffff000000000000000000000000000000000000000090811690915573ffffffffffffffffffffffffffffffffffffffff8781168086526003855283862080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01905590871680865283862080546001019055868652600290945282852080549092168417909155905184937fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef91a4505050565b60065473ffffffffffffffffffffffffffffffffffffffff163314610ba8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016106f8565b6006805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff0000000000000000000000000000000000000000831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b8173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603611a0d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f4552433732313a20617070726f766520746f2063616c6c65720000000000000060448201526064016106f8565b73ffffffffffffffffffffffffffffffffffffffff83811660008181526005602090815260408083209487168084529482529182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001686151590811790915591519182527f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31910160405180910390a3505050565b611ab0848484611585565b611abc84848484611f03565b610d0f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603260248201527f4552433732313a207472616e7366657220746f206e6f6e20455243373231526560448201527f63656976657220696d706c656d656e746572000000000000000000000000000060648201526084016106f8565b60606000611b55836120f3565b600101905060008167ffffffffffffffff811115611b7557611b756123f0565b6040519080825280601f01601f191660200182016040528015611b9f576020820181803683370190505b5090508181016020015b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff017f3031323334353637383961626364656600000000000000000000000000000000600a86061a8153600a8504945084611ba957509392505050565b6000611c118261111a565b9050611c1c8261111a565b600083815260046020908152604080832080547fffffffffffffffffffffffff000000000000000000000000000000000000000090811690915573ffffffffffffffffffffffffffffffffffffffff85168085526003845282852080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190558785526002909352818420805490911690555192935084927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef908390a45050565b73ffffffffffffffffffffffffffffffffffffffff8216611d5b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4552433732313a206d696e7420746f20746865207a65726f206164647265737360448201526064016106f8565b60008181526002602052604090205473ffffffffffffffffffffffffffffffffffffffff1615611de7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f4552433732313a20746f6b656e20616c7265616479206d696e7465640000000060448201526064016106f8565b60008181526002602052604090205473ffffffffffffffffffffffffffffffffffffffff1615611e73576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f4552433732313a20746f6b656e20616c7265616479206d696e7465640000000060448201526064016106f8565b73ffffffffffffffffffffffffffffffffffffffff8216600081815260036020908152604080832080546001019055848352600290915280822080547fffffffffffffffffffffffff0000000000000000000000000000000000000000168417905551839291907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef908290a45050565b600073ffffffffffffffffffffffffffffffffffffffff84163b156120eb576040517f150b7a0200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff85169063150b7a0290611f7a903390899088908890600401612686565b6020604051808303816000875af1925050508015611fd3575060408051601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201909252611fd0918101906126e1565b60015b6120a0573d808015612001576040519150601f19603f3d011682016040523d82523d6000602084013e612006565b606091505b508051600003612098576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603260248201527f4552433732313a207472616e7366657220746f206e6f6e20455243373231526560448201527f63656976657220696d706c656d656e746572000000000000000000000000000060648201526084016106f8565b805181602001fd5b7fffffffff00000000000000000000000000000000000000000000000000000000167f150b7a02000000000000000000000000000000000000000000000000000000001490506107d4565b5060016107d4565b6000807a184f03e93ff9f4daa797ed6e38ed64bf6a1f010000000000000000831061213c577a184f03e93ff9f4daa797ed6e38ed64bf6a1f010000000000000000830492506040015b6d04ee2d6d415b85acef81000000008310612168576d04ee2d6d415b85acef8100000000830492506020015b662386f26fc10000831061218657662386f26fc10000830492506010015b6305f5e100831061219e576305f5e100830492506008015b61271083106121b257612710830492506004015b606483106121c4576064830492506002015b600a83106105685760010192915050565b7fffffffff000000000000000000000000000000000000000000000000000000008116811461100157600080fd5b60006020828403121561221557600080fd5b8135610f46816121d5565b60005b8381101561223b578181015183820152602001612223565b50506000910152565b6000815180845261225c816020860160208601612220565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000610f466020830184612244565b6000602082840312156122b357600080fd5b5035919050565b73ffffffffffffffffffffffffffffffffffffffff8116811461100157600080fd5b600080604083850312156122ef57600080fd5b82356122fa816122ba565b946020939093013593505050565b60008060006060848603121561231d57600080fd5b83359250602084013561232f816122ba565b929592945050506040919091013590565b60008060006060848603121561235557600080fd5b8335612360816122ba565b9250602084013561232f816122ba565b6000806040838503121561238357600080fd5b823591506020830135612395816122ba565b809150509250929050565b6000602082840312156123b257600080fd5b8135610f46816122ba565b600080604083850312156123d057600080fd5b82356123db816122ba565b91506020830135801515811461239557600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6000806000806080858703121561243557600080fd5b8435612440816122ba565b93506020850135612450816122ba565b925060408501359150606085013567ffffffffffffffff81111561247357600080fd5b8501601f8101871361248457600080fd5b803567ffffffffffffffff81111561249e5761249e6123f0565b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8501160116810181811067ffffffffffffffff8211171561250a5761250a6123f0565b60405281815282820160200189101561252257600080fd5b8160208401602083013760006020838301015280935050505092959194509250565b6000806040838503121561255757600080fd5b50508035926020909101359150565b6000806040838503121561257957600080fd5b8235612584816122ba565b91506020830135612395816122ba565b600181811c908216806125a857607f821691505b6020821081036125e1577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b6000602082840312156125f957600080fd5b8151610f46816122ba565b60006020828403121561261657600080fd5b5051919050565b80820180821115610568577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60008351612669818460208801612220565b83519083019061267d818360208801612220565b01949350505050565b73ffffffffffffffffffffffffffffffffffffffff8516815273ffffffffffffffffffffffffffffffffffffffff841660208201528260408201526080606082015260006126d76080830184612244565b9695505050505050565b6000602082840312156126f357600080fd5b8151610f46816121d556fea264697066735822122087690f6a5ae035e7e9697abb37ef5b1005c3edfdf0dc9ad3a01755ad7b6967a564736f6c634300081a0033",
}

// BaseRegistrarImplementationABI is the input ABI used to generate the binding from.
// Deprecated: Use BaseRegistrarImplementationMetaData.ABI instead.
var BaseRegistrarImplementationABI = BaseRegistrarImplementationMetaData.ABI

// BaseRegistrarImplementationBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BaseRegistrarImplementationMetaData.Bin instead.
var BaseRegistrarImplementationBin = BaseRegistrarImplementationMetaData.Bin

// DeployBaseRegistrarImplementation deploys a new Ethereum contract, binding an instance of BaseRegistrarImplementation to it.
func DeployBaseRegistrarImplementation(auth *bind.TransactOpts, backend bind.ContractBackend, _ens common.Address, _baseNode [32]byte) (common.Address, *types.Transaction, *BaseRegistrarImplementation, error) {
	parsed, err := BaseRegistrarImplementationMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BaseRegistrarImplementationBin), backend, _ens, _baseNode)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BaseRegistrarImplementation{BaseRegistrarImplementationCaller: BaseRegistrarImplementationCaller{contract: contract}, BaseRegistrarImplementationTransactor: BaseRegistrarImplementationTransactor{contract: contract}, BaseRegistrarImplementationFilterer: BaseRegistrarImplementationFilterer{contract: contract}}, nil
}

// BaseRegistrarImplementation is an auto generated Go binding around an Ethereum contract.
type BaseRegistrarImplementation struct {
	BaseRegistrarImplementationCaller     // Read-only binding to the contract
	BaseRegistrarImplementationTransactor // Write-only binding to the contract
	BaseRegistrarImplementationFilterer   // Log filterer for contract events
}

// BaseRegistrarImplementationCaller is an auto generated read-only Go binding around an Ethereum contract.
type BaseRegistrarImplementationCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BaseRegistrarImplementationTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BaseRegistrarImplementationTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BaseRegistrarImplementationFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BaseRegistrarImplementationFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BaseRegistrarImplementationSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BaseRegistrarImplementationSession struct {
	Contract     *BaseRegistrarImplementation // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                // Call options to use throughout this session
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// BaseRegistrarImplementationCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BaseRegistrarImplementationCallerSession struct {
	Contract *BaseRegistrarImplementationCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                      // Call options to use throughout this session
}

// BaseRegistrarImplementationTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BaseRegistrarImplementationTransactorSession struct {
	Contract     *BaseRegistrarImplementationTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                      // Transaction auth options to use throughout this session
}

// BaseRegistrarImplementationRaw is an auto generated low-level Go binding around an Ethereum contract.
type BaseRegistrarImplementationRaw struct {
	Contract *BaseRegistrarImplementation // Generic contract binding to access the raw methods on
}

// BaseRegistrarImplementationCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BaseRegistrarImplementationCallerRaw struct {
	Contract *BaseRegistrarImplementationCaller // Generic read-only contract binding to access the raw methods on
}

// BaseRegistrarImplementationTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BaseRegistrarImplementationTransactorRaw struct {
	Contract *BaseRegistrarImplementationTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBaseRegistrarImplementation creates a new instance of BaseRegistrarImplementation, bound to a specific deployed contract.
func NewBaseRegistrarImplementation(address common.Address, backend bind.ContractBackend) (*BaseRegistrarImplementation, error) {
	contract, err := bindBaseRegistrarImplementation(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BaseRegistrarImplementation{BaseRegistrarImplementationCaller: BaseRegistrarImplementationCaller{contract: contract}, BaseRegistrarImplementationTransactor: BaseRegistrarImplementationTransactor{contract: contract}, BaseRegistrarImplementationFilterer: BaseRegistrarImplementationFilterer{contract: contract}}, nil
}

// NewBaseRegistrarImplementationCaller creates a new read-only instance of BaseRegistrarImplementation, bound to a specific deployed contract.
func NewBaseRegistrarImplementationCaller(address common.Address, caller bind.ContractCaller) (*BaseRegistrarImplementationCaller, error) {
	contract, err := bindBaseRegistrarImplementation(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BaseRegistrarImplementationCaller{contract: contract}, nil
}

// NewBaseRegistrarImplementationTransactor creates a new write-only instance of BaseRegistrarImplementation, bound to a specific deployed contract.
func NewBaseRegistrarImplementationTransactor(address common.Address, transactor bind.ContractTransactor) (*BaseRegistrarImplementationTransactor, error) {
	contract, err := bindBaseRegistrarImplementation(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BaseRegistrarImplementationTransactor{contract: contract}, nil
}

// NewBaseRegistrarImplementationFilterer creates a new log filterer instance of BaseRegistrarImplementation, bound to a specific deployed contract.
func NewBaseRegistrarImplementationFilterer(address common.Address, filterer bind.ContractFilterer) (*BaseRegistrarImplementationFilterer, error) {
	contract, err := bindBaseRegistrarImplementation(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BaseRegistrarImplementationFilterer{contract: contract}, nil
}

// bindBaseRegistrarImplementation binds a generic wrapper to an already deployed contract.
func bindBaseRegistrarImplementation(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BaseRegistrarImplementationMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BaseRegistrarImplementation *BaseRegistrarImplementationRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BaseRegistrarImplementation.Contract.BaseRegistrarImplementationCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BaseRegistrarImplementation *BaseRegistrarImplementationRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.Contract.BaseRegistrarImplementationTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BaseRegistrarImplementation *BaseRegistrarImplementationRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.Contract.BaseRegistrarImplementationTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BaseRegistrarImplementation *BaseRegistrarImplementationCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BaseRegistrarImplementation.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BaseRegistrarImplementation *BaseRegistrarImplementationTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BaseRegistrarImplementation *BaseRegistrarImplementationTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.Contract.contract.Transact(opts, method, params...)
}

// GRACEPERIOD is a free data retrieval call binding the contract method 0xc1a287e2.
//
// Solidity: function GRACE_PERIOD() view returns(uint256)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationCaller) GRACEPERIOD(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BaseRegistrarImplementation.contract.Call(opts, &out, "GRACE_PERIOD")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GRACEPERIOD is a free data retrieval call binding the contract method 0xc1a287e2.
//
// Solidity: function GRACE_PERIOD() view returns(uint256)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationSession) GRACEPERIOD() (*big.Int, error) {
	return _BaseRegistrarImplementation.Contract.GRACEPERIOD(&_BaseRegistrarImplementation.CallOpts)
}

// GRACEPERIOD is a free data retrieval call binding the contract method 0xc1a287e2.
//
// Solidity: function GRACE_PERIOD() view returns(uint256)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationCallerSession) GRACEPERIOD() (*big.Int, error) {
	return _BaseRegistrarImplementation.Contract.GRACEPERIOD(&_BaseRegistrarImplementation.CallOpts)
}

// Available is a free data retrieval call binding the contract method 0x96e494e8.
//
// Solidity: function available(uint256 id) view returns(bool)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationCaller) Available(opts *bind.CallOpts, id *big.Int) (bool, error) {
	var out []interface{}
	err := _BaseRegistrarImplementation.contract.Call(opts, &out, "available", id)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Available is a free data retrieval call binding the contract method 0x96e494e8.
//
// Solidity: function available(uint256 id) view returns(bool)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationSession) Available(id *big.Int) (bool, error) {
	return _BaseRegistrarImplementation.Contract.Available(&_BaseRegistrarImplementation.CallOpts, id)
}

// Available is a free data retrieval call binding the contract method 0x96e494e8.
//
// Solidity: function available(uint256 id) view returns(bool)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationCallerSession) Available(id *big.Int) (bool, error) {
	return _BaseRegistrarImplementation.Contract.Available(&_BaseRegistrarImplementation.CallOpts, id)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _BaseRegistrarImplementation.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _BaseRegistrarImplementation.Contract.BalanceOf(&_BaseRegistrarImplementation.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _BaseRegistrarImplementation.Contract.BalanceOf(&_BaseRegistrarImplementation.CallOpts, owner)
}

// BaseNode is a free data retrieval call binding the contract method 0xddf7fcb0.
//
// Solidity: function baseNode() view returns(bytes32)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationCaller) BaseNode(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _BaseRegistrarImplementation.contract.Call(opts, &out, "baseNode")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// BaseNode is a free data retrieval call binding the contract method 0xddf7fcb0.
//
// Solidity: function baseNode() view returns(bytes32)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationSession) BaseNode() ([32]byte, error) {
	return _BaseRegistrarImplementation.Contract.BaseNode(&_BaseRegistrarImplementation.CallOpts)
}

// BaseNode is a free data retrieval call binding the contract method 0xddf7fcb0.
//
// Solidity: function baseNode() view returns(bytes32)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationCallerSession) BaseNode() ([32]byte, error) {
	return _BaseRegistrarImplementation.Contract.BaseNode(&_BaseRegistrarImplementation.CallOpts)
}

// Controllers is a free data retrieval call binding the contract method 0xda8c229e.
//
// Solidity: function controllers(address ) view returns(bool)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationCaller) Controllers(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _BaseRegistrarImplementation.contract.Call(opts, &out, "controllers", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Controllers is a free data retrieval call binding the contract method 0xda8c229e.
//
// Solidity: function controllers(address ) view returns(bool)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationSession) Controllers(arg0 common.Address) (bool, error) {
	return _BaseRegistrarImplementation.Contract.Controllers(&_BaseRegistrarImplementation.CallOpts, arg0)
}

// Controllers is a free data retrieval call binding the contract method 0xda8c229e.
//
// Solidity: function controllers(address ) view returns(bool)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationCallerSession) Controllers(arg0 common.Address) (bool, error) {
	return _BaseRegistrarImplementation.Contract.Controllers(&_BaseRegistrarImplementation.CallOpts, arg0)
}

// Ens is a free data retrieval call binding the contract method 0x3f15457f.
//
// Solidity: function ens() view returns(address)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationCaller) Ens(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BaseRegistrarImplementation.contract.Call(opts, &out, "ens")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Ens is a free data retrieval call binding the contract method 0x3f15457f.
//
// Solidity: function ens() view returns(address)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationSession) Ens() (common.Address, error) {
	return _BaseRegistrarImplementation.Contract.Ens(&_BaseRegistrarImplementation.CallOpts)
}

// Ens is a free data retrieval call binding the contract method 0x3f15457f.
//
// Solidity: function ens() view returns(address)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationCallerSession) Ens() (common.Address, error) {
	return _BaseRegistrarImplementation.Contract.Ens(&_BaseRegistrarImplementation.CallOpts)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _BaseRegistrarImplementation.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _BaseRegistrarImplementation.Contract.GetApproved(&_BaseRegistrarImplementation.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _BaseRegistrarImplementation.Contract.GetApproved(&_BaseRegistrarImplementation.CallOpts, tokenId)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _BaseRegistrarImplementation.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _BaseRegistrarImplementation.Contract.IsApprovedForAll(&_BaseRegistrarImplementation.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _BaseRegistrarImplementation.Contract.IsApprovedForAll(&_BaseRegistrarImplementation.CallOpts, owner, operator)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _BaseRegistrarImplementation.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationSession) Name() (string, error) {
	return _BaseRegistrarImplementation.Contract.Name(&_BaseRegistrarImplementation.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationCallerSession) Name() (string, error) {
	return _BaseRegistrarImplementation.Contract.Name(&_BaseRegistrarImplementation.CallOpts)
}

// NameExpires is a free data retrieval call binding the contract method 0xd6e4fa86.
//
// Solidity: function nameExpires(uint256 id) view returns(uint256)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationCaller) NameExpires(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _BaseRegistrarImplementation.contract.Call(opts, &out, "nameExpires", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NameExpires is a free data retrieval call binding the contract method 0xd6e4fa86.
//
// Solidity: function nameExpires(uint256 id) view returns(uint256)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationSession) NameExpires(id *big.Int) (*big.Int, error) {
	return _BaseRegistrarImplementation.Contract.NameExpires(&_BaseRegistrarImplementation.CallOpts, id)
}

// NameExpires is a free data retrieval call binding the contract method 0xd6e4fa86.
//
// Solidity: function nameExpires(uint256 id) view returns(uint256)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationCallerSession) NameExpires(id *big.Int) (*big.Int, error) {
	return _BaseRegistrarImplementation.Contract.NameExpires(&_BaseRegistrarImplementation.CallOpts, id)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BaseRegistrarImplementation.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationSession) Owner() (common.Address, error) {
	return _BaseRegistrarImplementation.Contract.Owner(&_BaseRegistrarImplementation.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationCallerSession) Owner() (common.Address, error) {
	return _BaseRegistrarImplementation.Contract.Owner(&_BaseRegistrarImplementation.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _BaseRegistrarImplementation.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _BaseRegistrarImplementation.Contract.OwnerOf(&_BaseRegistrarImplementation.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _BaseRegistrarImplementation.Contract.OwnerOf(&_BaseRegistrarImplementation.CallOpts, tokenId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceID) view returns(bool)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationCaller) SupportsInterface(opts *bind.CallOpts, interfaceID [4]byte) (bool, error) {
	var out []interface{}
	err := _BaseRegistrarImplementation.contract.Call(opts, &out, "supportsInterface", interfaceID)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceID) view returns(bool)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationSession) SupportsInterface(interfaceID [4]byte) (bool, error) {
	return _BaseRegistrarImplementation.Contract.SupportsInterface(&_BaseRegistrarImplementation.CallOpts, interfaceID)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceID) view returns(bool)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationCallerSession) SupportsInterface(interfaceID [4]byte) (bool, error) {
	return _BaseRegistrarImplementation.Contract.SupportsInterface(&_BaseRegistrarImplementation.CallOpts, interfaceID)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _BaseRegistrarImplementation.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationSession) Symbol() (string, error) {
	return _BaseRegistrarImplementation.Contract.Symbol(&_BaseRegistrarImplementation.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationCallerSession) Symbol() (string, error) {
	return _BaseRegistrarImplementation.Contract.Symbol(&_BaseRegistrarImplementation.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationCaller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _BaseRegistrarImplementation.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationSession) TokenURI(tokenId *big.Int) (string, error) {
	return _BaseRegistrarImplementation.Contract.TokenURI(&_BaseRegistrarImplementation.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationCallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _BaseRegistrarImplementation.Contract.TokenURI(&_BaseRegistrarImplementation.CallOpts, tokenId)
}

// AddController is a paid mutator transaction binding the contract method 0xa7fc7a07.
//
// Solidity: function addController(address controller) returns()
func (_BaseRegistrarImplementation *BaseRegistrarImplementationTransactor) AddController(opts *bind.TransactOpts, controller common.Address) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.contract.Transact(opts, "addController", controller)
}

// AddController is a paid mutator transaction binding the contract method 0xa7fc7a07.
//
// Solidity: function addController(address controller) returns()
func (_BaseRegistrarImplementation *BaseRegistrarImplementationSession) AddController(controller common.Address) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.Contract.AddController(&_BaseRegistrarImplementation.TransactOpts, controller)
}

// AddController is a paid mutator transaction binding the contract method 0xa7fc7a07.
//
// Solidity: function addController(address controller) returns()
func (_BaseRegistrarImplementation *BaseRegistrarImplementationTransactorSession) AddController(controller common.Address) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.Contract.AddController(&_BaseRegistrarImplementation.TransactOpts, controller)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_BaseRegistrarImplementation *BaseRegistrarImplementationTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_BaseRegistrarImplementation *BaseRegistrarImplementationSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.Contract.Approve(&_BaseRegistrarImplementation.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_BaseRegistrarImplementation *BaseRegistrarImplementationTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.Contract.Approve(&_BaseRegistrarImplementation.TransactOpts, to, tokenId)
}

// Reclaim is a paid mutator transaction binding the contract method 0x28ed4f6c.
//
// Solidity: function reclaim(uint256 id, address owner) returns()
func (_BaseRegistrarImplementation *BaseRegistrarImplementationTransactor) Reclaim(opts *bind.TransactOpts, id *big.Int, owner common.Address) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.contract.Transact(opts, "reclaim", id, owner)
}

// Reclaim is a paid mutator transaction binding the contract method 0x28ed4f6c.
//
// Solidity: function reclaim(uint256 id, address owner) returns()
func (_BaseRegistrarImplementation *BaseRegistrarImplementationSession) Reclaim(id *big.Int, owner common.Address) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.Contract.Reclaim(&_BaseRegistrarImplementation.TransactOpts, id, owner)
}

// Reclaim is a paid mutator transaction binding the contract method 0x28ed4f6c.
//
// Solidity: function reclaim(uint256 id, address owner) returns()
func (_BaseRegistrarImplementation *BaseRegistrarImplementationTransactorSession) Reclaim(id *big.Int, owner common.Address) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.Contract.Reclaim(&_BaseRegistrarImplementation.TransactOpts, id, owner)
}

// Register is a paid mutator transaction binding the contract method 0xfca247ac.
//
// Solidity: function register(uint256 id, address owner, uint256 duration) returns(uint256)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationTransactor) Register(opts *bind.TransactOpts, id *big.Int, owner common.Address, duration *big.Int) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.contract.Transact(opts, "register", id, owner, duration)
}

// Register is a paid mutator transaction binding the contract method 0xfca247ac.
//
// Solidity: function register(uint256 id, address owner, uint256 duration) returns(uint256)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationSession) Register(id *big.Int, owner common.Address, duration *big.Int) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.Contract.Register(&_BaseRegistrarImplementation.TransactOpts, id, owner, duration)
}

// Register is a paid mutator transaction binding the contract method 0xfca247ac.
//
// Solidity: function register(uint256 id, address owner, uint256 duration) returns(uint256)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationTransactorSession) Register(id *big.Int, owner common.Address, duration *big.Int) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.Contract.Register(&_BaseRegistrarImplementation.TransactOpts, id, owner, duration)
}

// RegisterOnly is a paid mutator transaction binding the contract method 0x0e297b45.
//
// Solidity: function registerOnly(uint256 id, address owner, uint256 duration) returns(uint256)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationTransactor) RegisterOnly(opts *bind.TransactOpts, id *big.Int, owner common.Address, duration *big.Int) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.contract.Transact(opts, "registerOnly", id, owner, duration)
}

// RegisterOnly is a paid mutator transaction binding the contract method 0x0e297b45.
//
// Solidity: function registerOnly(uint256 id, address owner, uint256 duration) returns(uint256)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationSession) RegisterOnly(id *big.Int, owner common.Address, duration *big.Int) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.Contract.RegisterOnly(&_BaseRegistrarImplementation.TransactOpts, id, owner, duration)
}

// RegisterOnly is a paid mutator transaction binding the contract method 0x0e297b45.
//
// Solidity: function registerOnly(uint256 id, address owner, uint256 duration) returns(uint256)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationTransactorSession) RegisterOnly(id *big.Int, owner common.Address, duration *big.Int) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.Contract.RegisterOnly(&_BaseRegistrarImplementation.TransactOpts, id, owner, duration)
}

// RemoveController is a paid mutator transaction binding the contract method 0xf6a74ed7.
//
// Solidity: function removeController(address controller) returns()
func (_BaseRegistrarImplementation *BaseRegistrarImplementationTransactor) RemoveController(opts *bind.TransactOpts, controller common.Address) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.contract.Transact(opts, "removeController", controller)
}

// RemoveController is a paid mutator transaction binding the contract method 0xf6a74ed7.
//
// Solidity: function removeController(address controller) returns()
func (_BaseRegistrarImplementation *BaseRegistrarImplementationSession) RemoveController(controller common.Address) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.Contract.RemoveController(&_BaseRegistrarImplementation.TransactOpts, controller)
}

// RemoveController is a paid mutator transaction binding the contract method 0xf6a74ed7.
//
// Solidity: function removeController(address controller) returns()
func (_BaseRegistrarImplementation *BaseRegistrarImplementationTransactorSession) RemoveController(controller common.Address) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.Contract.RemoveController(&_BaseRegistrarImplementation.TransactOpts, controller)
}

// Renew is a paid mutator transaction binding the contract method 0xc475abff.
//
// Solidity: function renew(uint256 id, uint256 duration) returns(uint256)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationTransactor) Renew(opts *bind.TransactOpts, id *big.Int, duration *big.Int) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.contract.Transact(opts, "renew", id, duration)
}

// Renew is a paid mutator transaction binding the contract method 0xc475abff.
//
// Solidity: function renew(uint256 id, uint256 duration) returns(uint256)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationSession) Renew(id *big.Int, duration *big.Int) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.Contract.Renew(&_BaseRegistrarImplementation.TransactOpts, id, duration)
}

// Renew is a paid mutator transaction binding the contract method 0xc475abff.
//
// Solidity: function renew(uint256 id, uint256 duration) returns(uint256)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationTransactorSession) Renew(id *big.Int, duration *big.Int) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.Contract.Renew(&_BaseRegistrarImplementation.TransactOpts, id, duration)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BaseRegistrarImplementation *BaseRegistrarImplementationTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BaseRegistrarImplementation *BaseRegistrarImplementationSession) RenounceOwnership() (*types.Transaction, error) {
	return _BaseRegistrarImplementation.Contract.RenounceOwnership(&_BaseRegistrarImplementation.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BaseRegistrarImplementation *BaseRegistrarImplementationTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _BaseRegistrarImplementation.Contract.RenounceOwnership(&_BaseRegistrarImplementation.TransactOpts)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_BaseRegistrarImplementation *BaseRegistrarImplementationTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_BaseRegistrarImplementation *BaseRegistrarImplementationSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.Contract.SafeTransferFrom(&_BaseRegistrarImplementation.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_BaseRegistrarImplementation *BaseRegistrarImplementationTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.Contract.SafeTransferFrom(&_BaseRegistrarImplementation.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_BaseRegistrarImplementation *BaseRegistrarImplementationTransactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_BaseRegistrarImplementation *BaseRegistrarImplementationSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.Contract.SafeTransferFrom0(&_BaseRegistrarImplementation.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_BaseRegistrarImplementation *BaseRegistrarImplementationTransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.Contract.SafeTransferFrom0(&_BaseRegistrarImplementation.TransactOpts, from, to, tokenId, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_BaseRegistrarImplementation *BaseRegistrarImplementationTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_BaseRegistrarImplementation *BaseRegistrarImplementationSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.Contract.SetApprovalForAll(&_BaseRegistrarImplementation.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_BaseRegistrarImplementation *BaseRegistrarImplementationTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.Contract.SetApprovalForAll(&_BaseRegistrarImplementation.TransactOpts, operator, approved)
}

// SetResolver is a paid mutator transaction binding the contract method 0x4e543b26.
//
// Solidity: function setResolver(address resolver) returns()
func (_BaseRegistrarImplementation *BaseRegistrarImplementationTransactor) SetResolver(opts *bind.TransactOpts, resolver common.Address) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.contract.Transact(opts, "setResolver", resolver)
}

// SetResolver is a paid mutator transaction binding the contract method 0x4e543b26.
//
// Solidity: function setResolver(address resolver) returns()
func (_BaseRegistrarImplementation *BaseRegistrarImplementationSession) SetResolver(resolver common.Address) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.Contract.SetResolver(&_BaseRegistrarImplementation.TransactOpts, resolver)
}

// SetResolver is a paid mutator transaction binding the contract method 0x4e543b26.
//
// Solidity: function setResolver(address resolver) returns()
func (_BaseRegistrarImplementation *BaseRegistrarImplementationTransactorSession) SetResolver(resolver common.Address) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.Contract.SetResolver(&_BaseRegistrarImplementation.TransactOpts, resolver)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_BaseRegistrarImplementation *BaseRegistrarImplementationTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_BaseRegistrarImplementation *BaseRegistrarImplementationSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.Contract.TransferFrom(&_BaseRegistrarImplementation.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_BaseRegistrarImplementation *BaseRegistrarImplementationTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.Contract.TransferFrom(&_BaseRegistrarImplementation.TransactOpts, from, to, tokenId)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BaseRegistrarImplementation *BaseRegistrarImplementationTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BaseRegistrarImplementation *BaseRegistrarImplementationSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.Contract.TransferOwnership(&_BaseRegistrarImplementation.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BaseRegistrarImplementation *BaseRegistrarImplementationTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _BaseRegistrarImplementation.Contract.TransferOwnership(&_BaseRegistrarImplementation.TransactOpts, newOwner)
}

// BaseRegistrarImplementationApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the BaseRegistrarImplementation contract.
type BaseRegistrarImplementationApprovalIterator struct {
	Event *BaseRegistrarImplementationApproval // Event containing the contract specifics and raw log

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
func (it *BaseRegistrarImplementationApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BaseRegistrarImplementationApproval)
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
		it.Event = new(BaseRegistrarImplementationApproval)
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
func (it *BaseRegistrarImplementationApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BaseRegistrarImplementationApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BaseRegistrarImplementationApproval represents a Approval event raised by the BaseRegistrarImplementation contract.
type BaseRegistrarImplementationApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*BaseRegistrarImplementationApprovalIterator, error) {

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

	logs, sub, err := _BaseRegistrarImplementation.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &BaseRegistrarImplementationApprovalIterator{contract: _BaseRegistrarImplementation.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *BaseRegistrarImplementationApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _BaseRegistrarImplementation.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BaseRegistrarImplementationApproval)
				if err := _BaseRegistrarImplementation.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_BaseRegistrarImplementation *BaseRegistrarImplementationFilterer) ParseApproval(log types.Log) (*BaseRegistrarImplementationApproval, error) {
	event := new(BaseRegistrarImplementationApproval)
	if err := _BaseRegistrarImplementation.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BaseRegistrarImplementationApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the BaseRegistrarImplementation contract.
type BaseRegistrarImplementationApprovalForAllIterator struct {
	Event *BaseRegistrarImplementationApprovalForAll // Event containing the contract specifics and raw log

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
func (it *BaseRegistrarImplementationApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BaseRegistrarImplementationApprovalForAll)
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
		it.Event = new(BaseRegistrarImplementationApprovalForAll)
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
func (it *BaseRegistrarImplementationApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BaseRegistrarImplementationApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BaseRegistrarImplementationApprovalForAll represents a ApprovalForAll event raised by the BaseRegistrarImplementation contract.
type BaseRegistrarImplementationApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*BaseRegistrarImplementationApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _BaseRegistrarImplementation.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &BaseRegistrarImplementationApprovalForAllIterator{contract: _BaseRegistrarImplementation.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *BaseRegistrarImplementationApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _BaseRegistrarImplementation.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BaseRegistrarImplementationApprovalForAll)
				if err := _BaseRegistrarImplementation.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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
func (_BaseRegistrarImplementation *BaseRegistrarImplementationFilterer) ParseApprovalForAll(log types.Log) (*BaseRegistrarImplementationApprovalForAll, error) {
	event := new(BaseRegistrarImplementationApprovalForAll)
	if err := _BaseRegistrarImplementation.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BaseRegistrarImplementationControllerAddedIterator is returned from FilterControllerAdded and is used to iterate over the raw logs and unpacked data for ControllerAdded events raised by the BaseRegistrarImplementation contract.
type BaseRegistrarImplementationControllerAddedIterator struct {
	Event *BaseRegistrarImplementationControllerAdded // Event containing the contract specifics and raw log

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
func (it *BaseRegistrarImplementationControllerAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BaseRegistrarImplementationControllerAdded)
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
		it.Event = new(BaseRegistrarImplementationControllerAdded)
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
func (it *BaseRegistrarImplementationControllerAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BaseRegistrarImplementationControllerAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BaseRegistrarImplementationControllerAdded represents a ControllerAdded event raised by the BaseRegistrarImplementation contract.
type BaseRegistrarImplementationControllerAdded struct {
	Controller common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterControllerAdded is a free log retrieval operation binding the contract event 0x0a8bb31534c0ed46f380cb867bd5c803a189ced9a764e30b3a4991a9901d7474.
//
// Solidity: event ControllerAdded(address indexed controller)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationFilterer) FilterControllerAdded(opts *bind.FilterOpts, controller []common.Address) (*BaseRegistrarImplementationControllerAddedIterator, error) {

	var controllerRule []interface{}
	for _, controllerItem := range controller {
		controllerRule = append(controllerRule, controllerItem)
	}

	logs, sub, err := _BaseRegistrarImplementation.contract.FilterLogs(opts, "ControllerAdded", controllerRule)
	if err != nil {
		return nil, err
	}
	return &BaseRegistrarImplementationControllerAddedIterator{contract: _BaseRegistrarImplementation.contract, event: "ControllerAdded", logs: logs, sub: sub}, nil
}

// WatchControllerAdded is a free log subscription operation binding the contract event 0x0a8bb31534c0ed46f380cb867bd5c803a189ced9a764e30b3a4991a9901d7474.
//
// Solidity: event ControllerAdded(address indexed controller)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationFilterer) WatchControllerAdded(opts *bind.WatchOpts, sink chan<- *BaseRegistrarImplementationControllerAdded, controller []common.Address) (event.Subscription, error) {

	var controllerRule []interface{}
	for _, controllerItem := range controller {
		controllerRule = append(controllerRule, controllerItem)
	}

	logs, sub, err := _BaseRegistrarImplementation.contract.WatchLogs(opts, "ControllerAdded", controllerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BaseRegistrarImplementationControllerAdded)
				if err := _BaseRegistrarImplementation.contract.UnpackLog(event, "ControllerAdded", log); err != nil {
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

// ParseControllerAdded is a log parse operation binding the contract event 0x0a8bb31534c0ed46f380cb867bd5c803a189ced9a764e30b3a4991a9901d7474.
//
// Solidity: event ControllerAdded(address indexed controller)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationFilterer) ParseControllerAdded(log types.Log) (*BaseRegistrarImplementationControllerAdded, error) {
	event := new(BaseRegistrarImplementationControllerAdded)
	if err := _BaseRegistrarImplementation.contract.UnpackLog(event, "ControllerAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BaseRegistrarImplementationControllerRemovedIterator is returned from FilterControllerRemoved and is used to iterate over the raw logs and unpacked data for ControllerRemoved events raised by the BaseRegistrarImplementation contract.
type BaseRegistrarImplementationControllerRemovedIterator struct {
	Event *BaseRegistrarImplementationControllerRemoved // Event containing the contract specifics and raw log

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
func (it *BaseRegistrarImplementationControllerRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BaseRegistrarImplementationControllerRemoved)
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
		it.Event = new(BaseRegistrarImplementationControllerRemoved)
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
func (it *BaseRegistrarImplementationControllerRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BaseRegistrarImplementationControllerRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BaseRegistrarImplementationControllerRemoved represents a ControllerRemoved event raised by the BaseRegistrarImplementation contract.
type BaseRegistrarImplementationControllerRemoved struct {
	Controller common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterControllerRemoved is a free log retrieval operation binding the contract event 0x33d83959be2573f5453b12eb9d43b3499bc57d96bd2f067ba44803c859e81113.
//
// Solidity: event ControllerRemoved(address indexed controller)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationFilterer) FilterControllerRemoved(opts *bind.FilterOpts, controller []common.Address) (*BaseRegistrarImplementationControllerRemovedIterator, error) {

	var controllerRule []interface{}
	for _, controllerItem := range controller {
		controllerRule = append(controllerRule, controllerItem)
	}

	logs, sub, err := _BaseRegistrarImplementation.contract.FilterLogs(opts, "ControllerRemoved", controllerRule)
	if err != nil {
		return nil, err
	}
	return &BaseRegistrarImplementationControllerRemovedIterator{contract: _BaseRegistrarImplementation.contract, event: "ControllerRemoved", logs: logs, sub: sub}, nil
}

// WatchControllerRemoved is a free log subscription operation binding the contract event 0x33d83959be2573f5453b12eb9d43b3499bc57d96bd2f067ba44803c859e81113.
//
// Solidity: event ControllerRemoved(address indexed controller)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationFilterer) WatchControllerRemoved(opts *bind.WatchOpts, sink chan<- *BaseRegistrarImplementationControllerRemoved, controller []common.Address) (event.Subscription, error) {

	var controllerRule []interface{}
	for _, controllerItem := range controller {
		controllerRule = append(controllerRule, controllerItem)
	}

	logs, sub, err := _BaseRegistrarImplementation.contract.WatchLogs(opts, "ControllerRemoved", controllerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BaseRegistrarImplementationControllerRemoved)
				if err := _BaseRegistrarImplementation.contract.UnpackLog(event, "ControllerRemoved", log); err != nil {
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

// ParseControllerRemoved is a log parse operation binding the contract event 0x33d83959be2573f5453b12eb9d43b3499bc57d96bd2f067ba44803c859e81113.
//
// Solidity: event ControllerRemoved(address indexed controller)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationFilterer) ParseControllerRemoved(log types.Log) (*BaseRegistrarImplementationControllerRemoved, error) {
	event := new(BaseRegistrarImplementationControllerRemoved)
	if err := _BaseRegistrarImplementation.contract.UnpackLog(event, "ControllerRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BaseRegistrarImplementationNameMigratedIterator is returned from FilterNameMigrated and is used to iterate over the raw logs and unpacked data for NameMigrated events raised by the BaseRegistrarImplementation contract.
type BaseRegistrarImplementationNameMigratedIterator struct {
	Event *BaseRegistrarImplementationNameMigrated // Event containing the contract specifics and raw log

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
func (it *BaseRegistrarImplementationNameMigratedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BaseRegistrarImplementationNameMigrated)
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
		it.Event = new(BaseRegistrarImplementationNameMigrated)
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
func (it *BaseRegistrarImplementationNameMigratedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BaseRegistrarImplementationNameMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BaseRegistrarImplementationNameMigrated represents a NameMigrated event raised by the BaseRegistrarImplementation contract.
type BaseRegistrarImplementationNameMigrated struct {
	Id      *big.Int
	Owner   common.Address
	Expires *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterNameMigrated is a free log retrieval operation binding the contract event 0xea3d7e1195a15d2ddcd859b01abd4c6b960fa9f9264e499a70a90c7f0c64b717.
//
// Solidity: event NameMigrated(uint256 indexed id, address indexed owner, uint256 expires)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationFilterer) FilterNameMigrated(opts *bind.FilterOpts, id []*big.Int, owner []common.Address) (*BaseRegistrarImplementationNameMigratedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _BaseRegistrarImplementation.contract.FilterLogs(opts, "NameMigrated", idRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return &BaseRegistrarImplementationNameMigratedIterator{contract: _BaseRegistrarImplementation.contract, event: "NameMigrated", logs: logs, sub: sub}, nil
}

// WatchNameMigrated is a free log subscription operation binding the contract event 0xea3d7e1195a15d2ddcd859b01abd4c6b960fa9f9264e499a70a90c7f0c64b717.
//
// Solidity: event NameMigrated(uint256 indexed id, address indexed owner, uint256 expires)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationFilterer) WatchNameMigrated(opts *bind.WatchOpts, sink chan<- *BaseRegistrarImplementationNameMigrated, id []*big.Int, owner []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _BaseRegistrarImplementation.contract.WatchLogs(opts, "NameMigrated", idRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BaseRegistrarImplementationNameMigrated)
				if err := _BaseRegistrarImplementation.contract.UnpackLog(event, "NameMigrated", log); err != nil {
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

// ParseNameMigrated is a log parse operation binding the contract event 0xea3d7e1195a15d2ddcd859b01abd4c6b960fa9f9264e499a70a90c7f0c64b717.
//
// Solidity: event NameMigrated(uint256 indexed id, address indexed owner, uint256 expires)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationFilterer) ParseNameMigrated(log types.Log) (*BaseRegistrarImplementationNameMigrated, error) {
	event := new(BaseRegistrarImplementationNameMigrated)
	if err := _BaseRegistrarImplementation.contract.UnpackLog(event, "NameMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BaseRegistrarImplementationNameRegisteredIterator is returned from FilterNameRegistered and is used to iterate over the raw logs and unpacked data for NameRegistered events raised by the BaseRegistrarImplementation contract.
type BaseRegistrarImplementationNameRegisteredIterator struct {
	Event *BaseRegistrarImplementationNameRegistered // Event containing the contract specifics and raw log

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
func (it *BaseRegistrarImplementationNameRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BaseRegistrarImplementationNameRegistered)
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
		it.Event = new(BaseRegistrarImplementationNameRegistered)
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
func (it *BaseRegistrarImplementationNameRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BaseRegistrarImplementationNameRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BaseRegistrarImplementationNameRegistered represents a NameRegistered event raised by the BaseRegistrarImplementation contract.
type BaseRegistrarImplementationNameRegistered struct {
	Id      *big.Int
	Owner   common.Address
	Expires *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterNameRegistered is a free log retrieval operation binding the contract event 0xb3d987963d01b2f68493b4bdb130988f157ea43070d4ad840fee0466ed9370d9.
//
// Solidity: event NameRegistered(uint256 indexed id, address indexed owner, uint256 expires)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationFilterer) FilterNameRegistered(opts *bind.FilterOpts, id []*big.Int, owner []common.Address) (*BaseRegistrarImplementationNameRegisteredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _BaseRegistrarImplementation.contract.FilterLogs(opts, "NameRegistered", idRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return &BaseRegistrarImplementationNameRegisteredIterator{contract: _BaseRegistrarImplementation.contract, event: "NameRegistered", logs: logs, sub: sub}, nil
}

// WatchNameRegistered is a free log subscription operation binding the contract event 0xb3d987963d01b2f68493b4bdb130988f157ea43070d4ad840fee0466ed9370d9.
//
// Solidity: event NameRegistered(uint256 indexed id, address indexed owner, uint256 expires)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationFilterer) WatchNameRegistered(opts *bind.WatchOpts, sink chan<- *BaseRegistrarImplementationNameRegistered, id []*big.Int, owner []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _BaseRegistrarImplementation.contract.WatchLogs(opts, "NameRegistered", idRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BaseRegistrarImplementationNameRegistered)
				if err := _BaseRegistrarImplementation.contract.UnpackLog(event, "NameRegistered", log); err != nil {
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

// ParseNameRegistered is a log parse operation binding the contract event 0xb3d987963d01b2f68493b4bdb130988f157ea43070d4ad840fee0466ed9370d9.
//
// Solidity: event NameRegistered(uint256 indexed id, address indexed owner, uint256 expires)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationFilterer) ParseNameRegistered(log types.Log) (*BaseRegistrarImplementationNameRegistered, error) {
	event := new(BaseRegistrarImplementationNameRegistered)
	if err := _BaseRegistrarImplementation.contract.UnpackLog(event, "NameRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BaseRegistrarImplementationNameRenewedIterator is returned from FilterNameRenewed and is used to iterate over the raw logs and unpacked data for NameRenewed events raised by the BaseRegistrarImplementation contract.
type BaseRegistrarImplementationNameRenewedIterator struct {
	Event *BaseRegistrarImplementationNameRenewed // Event containing the contract specifics and raw log

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
func (it *BaseRegistrarImplementationNameRenewedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BaseRegistrarImplementationNameRenewed)
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
		it.Event = new(BaseRegistrarImplementationNameRenewed)
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
func (it *BaseRegistrarImplementationNameRenewedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BaseRegistrarImplementationNameRenewedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BaseRegistrarImplementationNameRenewed represents a NameRenewed event raised by the BaseRegistrarImplementation contract.
type BaseRegistrarImplementationNameRenewed struct {
	Id      *big.Int
	Expires *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterNameRenewed is a free log retrieval operation binding the contract event 0x9b87a00e30f1ac65d898f070f8a3488fe60517182d0a2098e1b4b93a54aa9bd6.
//
// Solidity: event NameRenewed(uint256 indexed id, uint256 expires)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationFilterer) FilterNameRenewed(opts *bind.FilterOpts, id []*big.Int) (*BaseRegistrarImplementationNameRenewedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _BaseRegistrarImplementation.contract.FilterLogs(opts, "NameRenewed", idRule)
	if err != nil {
		return nil, err
	}
	return &BaseRegistrarImplementationNameRenewedIterator{contract: _BaseRegistrarImplementation.contract, event: "NameRenewed", logs: logs, sub: sub}, nil
}

// WatchNameRenewed is a free log subscription operation binding the contract event 0x9b87a00e30f1ac65d898f070f8a3488fe60517182d0a2098e1b4b93a54aa9bd6.
//
// Solidity: event NameRenewed(uint256 indexed id, uint256 expires)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationFilterer) WatchNameRenewed(opts *bind.WatchOpts, sink chan<- *BaseRegistrarImplementationNameRenewed, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _BaseRegistrarImplementation.contract.WatchLogs(opts, "NameRenewed", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BaseRegistrarImplementationNameRenewed)
				if err := _BaseRegistrarImplementation.contract.UnpackLog(event, "NameRenewed", log); err != nil {
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

// ParseNameRenewed is a log parse operation binding the contract event 0x9b87a00e30f1ac65d898f070f8a3488fe60517182d0a2098e1b4b93a54aa9bd6.
//
// Solidity: event NameRenewed(uint256 indexed id, uint256 expires)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationFilterer) ParseNameRenewed(log types.Log) (*BaseRegistrarImplementationNameRenewed, error) {
	event := new(BaseRegistrarImplementationNameRenewed)
	if err := _BaseRegistrarImplementation.contract.UnpackLog(event, "NameRenewed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BaseRegistrarImplementationOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the BaseRegistrarImplementation contract.
type BaseRegistrarImplementationOwnershipTransferredIterator struct {
	Event *BaseRegistrarImplementationOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *BaseRegistrarImplementationOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BaseRegistrarImplementationOwnershipTransferred)
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
		it.Event = new(BaseRegistrarImplementationOwnershipTransferred)
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
func (it *BaseRegistrarImplementationOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BaseRegistrarImplementationOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BaseRegistrarImplementationOwnershipTransferred represents a OwnershipTransferred event raised by the BaseRegistrarImplementation contract.
type BaseRegistrarImplementationOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*BaseRegistrarImplementationOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _BaseRegistrarImplementation.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BaseRegistrarImplementationOwnershipTransferredIterator{contract: _BaseRegistrarImplementation.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BaseRegistrarImplementationOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _BaseRegistrarImplementation.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BaseRegistrarImplementationOwnershipTransferred)
				if err := _BaseRegistrarImplementation.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_BaseRegistrarImplementation *BaseRegistrarImplementationFilterer) ParseOwnershipTransferred(log types.Log) (*BaseRegistrarImplementationOwnershipTransferred, error) {
	event := new(BaseRegistrarImplementationOwnershipTransferred)
	if err := _BaseRegistrarImplementation.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BaseRegistrarImplementationTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the BaseRegistrarImplementation contract.
type BaseRegistrarImplementationTransferIterator struct {
	Event *BaseRegistrarImplementationTransfer // Event containing the contract specifics and raw log

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
func (it *BaseRegistrarImplementationTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BaseRegistrarImplementationTransfer)
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
		it.Event = new(BaseRegistrarImplementationTransfer)
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
func (it *BaseRegistrarImplementationTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BaseRegistrarImplementationTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BaseRegistrarImplementationTransfer represents a Transfer event raised by the BaseRegistrarImplementation contract.
type BaseRegistrarImplementationTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*BaseRegistrarImplementationTransferIterator, error) {

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

	logs, sub, err := _BaseRegistrarImplementation.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &BaseRegistrarImplementationTransferIterator{contract: _BaseRegistrarImplementation.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_BaseRegistrarImplementation *BaseRegistrarImplementationFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *BaseRegistrarImplementationTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _BaseRegistrarImplementation.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BaseRegistrarImplementationTransfer)
				if err := _BaseRegistrarImplementation.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_BaseRegistrarImplementation *BaseRegistrarImplementationFilterer) ParseTransfer(log types.Log) (*BaseRegistrarImplementationTransfer, error) {
	event := new(BaseRegistrarImplementationTransfer)
	if err := _BaseRegistrarImplementation.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
