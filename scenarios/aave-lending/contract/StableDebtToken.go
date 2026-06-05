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

// StableDebtTokenMetaData contains all meta data concerning the StableDebtToken contract.
var StableDebtTokenMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIPool\",\"name\":\"pool\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"fromUser\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"toUser\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BorrowAllowanceDelegated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"currentBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"balanceIncrease\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"avgStableRate\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newTotalSupply\",\"type\":\"uint256\"}],\"name\":\"Burn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"underlyingAsset\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"incentivesController\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"debtTokenDecimals\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"debtTokenName\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"debtTokenSymbol\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"params\",\"type\":\"bytes\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"onBehalfOf\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"currentBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"balanceIncrease\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newRate\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"avgStableRate\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newTotalSupply\",\"type\":\"uint256\"}],\"name\":\"Mint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEBT_TOKEN_REVISION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DELEGATION_WITH_SIG_TYPEHASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DOMAIN_SEPARATOR\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"EIP712_REVISION\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"POOL\",\"outputs\":[{\"internalType\":\"contractIPool\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"UNDERLYING_ASSET_ADDRESS\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegatee\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approveDelegation\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"fromUser\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"toUser\",\"type\":\"address\"}],\"name\":\"borrowAllowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegatee\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"delegationWithSig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAverageStableRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getIncentivesController\",\"outputs\":[{\"internalType\":\"contractIAaveIncentivesController\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSupplyData\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint40\",\"name\":\"\",\"type\":\"uint40\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalSupplyAndAvgRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalSupplyLastUpdated\",\"outputs\":[{\"internalType\":\"uint40\",\"name\":\"\",\"type\":\"uint40\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getUserLastUpdated\",\"outputs\":[{\"internalType\":\"uint40\",\"name\":\"\",\"type\":\"uint40\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getUserStableRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIPool\",\"name\":\"initializingPool\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"underlyingAsset\",\"type\":\"address\"},{\"internalType\":\"contractIAaveIncentivesController\",\"name\":\"incentivesController\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"debtTokenDecimals\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"debtTokenName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"debtTokenSymbol\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"params\",\"type\":\"bytes\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"onBehalfOf\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"rate\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"principalBalanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIAaveIncentivesController\",\"name\":\"controller\",\"type\":\"address\"}],\"name\":\"setIncentivesController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60e0604052600080553480156200001557600080fd5b5060405162002caf38038062002caf833981016040819052620000389162000234565b806040518060400160405280601681526020017f535441424c455f444542545f544f4b454e5f494d504c000000000000000000008152506040518060400160405280601681526020017f535441424c455f444542545f544f4b454e5f494d504c0000000000000000000081525060004660808181525050836001600160a01b0316630542975c6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620000ee573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000114919062000234565b6001600160a01b031660a05282516200013590603b90602086019062000175565b5081516200014b90603c90602085019062000175565b50603d805460ff191660ff9290921691909117905550506001600160a01b031660c0525062000298565b82805462000183906200025b565b90600052602060002090601f016020900481019282620001a75760008555620001f2565b82601f10620001c257805160ff1916838001178555620001f2565b82800160010185558215620001f2579182015b82811115620001f2578251825591602001919060010190620001d5565b506200020092915062000204565b5090565b5b8082111562000200576000815560010162000205565b6001600160a01b03811681146200023157600080fd5b50565b6000602082840312156200024757600080fd5b815162000254816200021b565b9392505050565b600181811c908216806200027057607f821691505b602082108114156200029257634e487b7160e01b600052602260045260246000fd5b50919050565b60805160a05160c0516129cb620002e46000396000818161030501528181610c4501528181611144015281816116a201526117fd015260006118c901526000610abd01526129cb6000f3fe608060405234801561001057600080fd5b506004361061020b5760003560e01c806390f6fcf21161012a578063c04a8a10116100bd578063e655dbd81161008c578063e78c9b3b11610071578063e78c9b3b146105b5578063f3bfc73814610611578063f731e9be1461063857600080fd5b8063e655dbd81461057f578063e74848901461059257600080fd5b8063c04a8a1014610503578063c222ec8a14610516578063c634dfaa14610529578063dd62ed3e1461057157600080fd5b8063a9059cbb116100f9578063a9059cbb1461022e578063b16a19de146104ad578063b3f1c93d146104cb578063b9a7b622146104fb57600080fd5b806390f6fcf21461046357806395d89b411461047d5780639dc29fac14610485578063a457c2d71461022e57600080fd5b80636bd76d24116101a25780637816037611610171578063781603761461036f57806379774338146103ab57806379ce6b8c146103da5780637ecebe001461042d57600080fd5b80636bd76d24146102a757806370a08231146102ed5780637535d2461461030057806375d264131461034c57600080fd5b806323b872dd116101de57806323b872dd1461027c578063313ce5671461028a5780633644e5151461029f578063395093511461022e57600080fd5b806306fdde0314610210578063095ea7b31461022e5780630b52d5581461025157806318160ddd14610266575b600080fd5b610218610640565b604051610225919061233c565b60405180910390f35b61024161023c36600461237f565b6106d2565b6040519015158152602001610225565b61026461025f3660046123bc565b610742565b005b61026e610a93565b604051908152602001610225565b61024161023c36600461242a565b603d5460405160ff9091168152602001610225565b61026e610ab9565b61026e6102b536600461246b565b73ffffffffffffffffffffffffffffffffffffffff918216600090815260366020908152604080832093909416825291909152205490565b61026e6102fb3660046124a4565b610af2565b6103277f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610225565b603d54610100900473ffffffffffffffffffffffffffffffffffffffff16610327565b6102186040518060400160405280600181526020017f310000000000000000000000000000000000000000000000000000000000000081525081565b6103b3610b9e565b6040805194855260208501939093529183015264ffffffffff166060820152608001610225565b6104176103e83660046124a4565b73ffffffffffffffffffffffffffffffffffffffff166000908152603e602052604090205464ffffffffff1690565b60405164ffffffffff9091168152602001610225565b61026e61043b3660046124a4565b73ffffffffffffffffffffffffffffffffffffffff1660009081526034602052604090205490565b603f546fffffffffffffffffffffffffffffffff1661026e565b610218610bfa565b61049861049336600461237f565b610c09565b60408051928352602083019190915201610225565b60375473ffffffffffffffffffffffffffffffffffffffff16610327565b6104de6104d93660046124c1565b611129565b604080519315158452602084019290925290820152606001610225565b61026e600181565b61026461051136600461237f565b6115ab565b610264610524366004612623565b6115ba565b61026e6105373660046124a4565b73ffffffffffffffffffffffffffffffffffffffff166000908152603860205260409020546fffffffffffffffffffffffffffffffff1690565b61026e61023c36600461246b565b61026461058d3660046124a4565b6118c5565b603f54700100000000000000000000000000000000900464ffffffffff16610417565b61026e6105c33660046124a4565b73ffffffffffffffffffffffffffffffffffffffff1660009081526038602052604090205470010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff1690565b61026e7f323db0410fecc107e39e2af5908671f4c8d106123b35a51501bb805c5fa36aa081565b610498611aa3565b6060603b805461064f906126f8565b80601f016020809104026020016040519081016040528092919081815260200182805461067b906126f8565b80156106c85780601f1061069d576101008083540402835291602001916106c8565b820191906000526020600020905b8154815290600101906020018083116106ab57829003601f168201915b5050505050905090565b604080518082018252600281527f3830000000000000000000000000000000000000000000000000000000000000602082015290517f08c379a00000000000000000000000000000000000000000000000000000000081526000916107399160040161233c565b60405180910390fd5b60408051808201909152600281527f3737000000000000000000000000000000000000000000000000000000000000602082015273ffffffffffffffffffffffffffffffffffffffff88166107c4576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610739919061233c565b50834211156040518060400160405280600281526020017f373800000000000000000000000000000000000000000000000000000000000081525090610837576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610739919061233c565b5073ffffffffffffffffffffffffffffffffffffffff871660009081526034602052604081205490610867610ab9565b604080517f323db0410fecc107e39e2af5908671f4c8d106123b35a51501bb805c5fa36aa0602082015273ffffffffffffffffffffffffffffffffffffffff8b1691810191909152606081018990526080810184905260a0810188905260c0016040516020818303038152906040528051906020012060405160200161091f9291907f190100000000000000000000000000000000000000000000000000000000000081526002810192909252602282015260420190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201206000845290830180835281905260ff8816918301919091526060820186905260808201859052915060019060a0016020604051602081039080840390855afa1580156109a5573d6000803e3d6000fd5b5050506020604051035173ffffffffffffffffffffffffffffffffffffffff168973ffffffffffffffffffffffffffffffffffffffff16146040518060400160405280600281526020017f373900000000000000000000000000000000000000000000000000000000000081525090610a4b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610739919061233c565b50610a5782600161277b565b73ffffffffffffffffffffffffffffffffffffffff8a16600090815260346020526040902055610a88898989611ace565b505050505050505050565b603f54600090610ab4906fffffffffffffffffffffffffffffffff16611b45565b905090565b60007f0000000000000000000000000000000000000000000000000000000000000000461415610aea575060355490565b610ab4611b94565b73ffffffffffffffffffffffffffffffffffffffff81166000908152603860205260408120546fffffffffffffffffffffffffffffffff8082169170010000000000000000000000000000000090041681610b51575060009392505050565b73ffffffffffffffffffffffffffffffffffffffff84166000908152603e6020526040812054610b8990839064ffffffffff16611c59565b9050610b958382611c6d565b95945050505050565b603f546000908190819081906fffffffffffffffffffffffffffffffff16610bc5603a5490565b610bce82611b45565b603f549197909650919450700100000000000000000000000000000000900464ffffffffff1692509050565b6060603c805461064f906126f8565b60408051808201909152600281527f323300000000000000000000000000000000000000000000000000000000000060208201526000908190337f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1614610cb2576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610739919061233c565b50600080610cbf86611cc4565b92509250506000610cce610a93565b73ffffffffffffffffffffffffffffffffffffffff881660009081526038602052604081205491925090819070010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff16888411610d5957603f80547fffffffffffffffffffffffffffffffff000000000000000000000000000000001690556000603a55610e53565b610d638985612793565b603a81905591506000610d93610d7886611d49565b603f546fffffffffffffffffffffffffffffffff1690611c6d565b90506000610daa610da38c611d49565b8490611c6d565b9050818110610de957603f80547fffffffffffffffffffffffffffffffff000000000000000000000000000000001690556000603a8190559450610e50565b610e0d610e08610df886611d49565b610e028486612793565b90611d64565b611da3565b603f80547fffffffffffffffffffffffffffffffff00000000000000000000000000000000166fffffffffffffffffffffffffffffffff92909216918217905594505b50505b85891415610ecb5773ffffffffffffffffffffffffffffffffffffffff8a16600090815260386020908152604080832080546fffffffffffffffffffffffffffffffff169055603e909152902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000000169055610f20565b73ffffffffffffffffffffffffffffffffffffffff8a166000908152603e6020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000000164264ffffffffff161790555b603f80547fffffffffffffffffffffff0000000000ffffffffffffffffffffffffffffffff167001000000000000000000000000000000004264ffffffffff160217905588851115611049576000610f788a87612793565b9050610f858b8287611e49565b60405181815273ffffffffffffffffffffffffffffffffffffffff8c16906000907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a36040805182815260208101899052908101879052606081018390526080810185905260a0810184905273ffffffffffffffffffffffffffffffffffffffff8c169081907fc16f4e4ca34d790de4c656c72fd015c667d688f20be64eea360618545c4c530f9060c00160405180910390a350611119565b6000611055868b612793565b90506110628b8287611fba565b60405181815260009073ffffffffffffffffffffffffffffffffffffffff8d16907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a36040805182815260208101899052908101879052606081018590526080810184905273ffffffffffffffffffffffffffffffffffffffff8c16907f44bd20a79e993bdcc7cbedf54a3b4d19fb78490124b6b90d04fe3242eea579e89060a00160405180910390a2505b50955093505050505b9250929050565b6000808073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000163373ffffffffffffffffffffffffffffffffffffffff16146040518060400160405280600281526020017f3233000000000000000000000000000000000000000000000000000000000000815250906111ea576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610739919061233c565b506112246040518060c001604052806000815260200160008152602001600081526020016000815260200160008152602001600081525090565b8673ffffffffffffffffffffffffffffffffffffffff168873ffffffffffffffffffffffffffffffffffffffff16146112625761126287898861200a565b60008061126e89611cc4565b925092505061127b610a93565b808452603f546fffffffffffffffffffffffffffffffff1660a08501526112a390899061277b565b603a81905560208401526112b688611d49565b60408481019190915273ffffffffffffffffffffffffffffffffffffffff8a1660009081526038602052205470010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff16606084015261135261132261131d8a8561277b565b611d49565b6040850151611331908a611c6d565b61134861133d86611d49565b606088015190611c6d565b610e02919061277b565b6080840181905261136290611da3565b73ffffffffffffffffffffffffffffffffffffffff8a16600090815260386020908152604080832080546fffffffffffffffffffffffffffffffff908116700100000000000000000000000000000000969091168602179055603e825290912080547fffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000000164264ffffffffff16908117909155603f80547fffffffffffffffffffffff0000000000ffffffffffffffffffffffffffffffff16919093021790915583015161146190610e089061143690611d49565b6040860151611446908b90611c6d565b6113486114568860000151611d49565b60a089015190611c6d565b603f80547fffffffffffffffffffffffffffffffff00000000000000000000000000000000166fffffffffffffffffffffffffffffffff92909216918217905560a084015260006114b2828a61277b565b90506114c38a828660000151611e49565b60405181815273ffffffffffffffffffffffffffffffffffffffff8b16906000907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a360808085015160a080870151602080890151604080518881529283018a9052820188905260608201949094529384015282015273ffffffffffffffffffffffffffffffffffffffff808c1691908d16907fc16f4e4ca34d790de4c656c72fd015c667d688f20be64eea360618545c4c530f9060c00160405180910390a35050602082015160a0909201519015999198509650945050505050565b6115b6338383611ace565b5050565b6001805460ff16806115cb5750303b155b806115d7575060005481115b611663576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f436f6e747261637420696e7374616e63652068617320616c726561647920626560448201527f656e20696e697469616c697a65640000000000000000000000000000000000006064820152608401610739565b60015460ff161580156116a057600180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00168117905560008290555b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168a73ffffffffffffffffffffffffffffffffffffffff16146040518060400160405280600281526020017f38370000000000000000000000000000000000000000000000000000000000008152509061175d576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610739919061233c565b50611767866120ca565b611770856120dd565b603d80546037805473ffffffffffffffffffffffffffffffffffffffff8d81167fffffffffffffffffffffffff0000000000000000000000000000000000000000909216919091179091558a16610100027fffffffffffffffffffffff00000000000000000000000000000000000000000090911660ff8a16171790556117f5611b94565b6035819055507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168973ffffffffffffffffffffffffffffffffffffffff167f40251fbfb6656cfa65a00d7879029fec1fad21d28fdcff2f4f68f52795b74f2c8a8a8a8a8a8a604051611882969594939291906127aa565b60405180910390a380156118b957600180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690555b50505050505050505050565b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663707cd7166040518163ffffffff1660e01b8152600401602060405180830381865afa158015611932573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611956919061284a565b6040517f7be53ca100000000000000000000000000000000000000000000000000000000815233600482015290915073ffffffffffffffffffffffffffffffffffffffff821690637be53ca190602401602060405180830381865afa1580156119c3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906119e79190612867565b6040518060400160405280600181526020017f310000000000000000000000000000000000000000000000000000000000000081525090611a55576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610739919061233c565b5050603d805473ffffffffffffffffffffffffffffffffffffffff909216610100027fffffffffffffffffffffff0000000000000000000000000000000000000000ff909216919091179055565b603f5460009081906fffffffffffffffffffffffffffffffff16611ac681611b45565b939092509050565b73ffffffffffffffffffffffffffffffffffffffff83811660008181526036602090815260408083208786168085529083529281902086905560375490518681529416939192917fda919360433220e13b51e8c211e490d148e61a3bd53de8c097194e458b97f3e1910160405180910390a4505050565b600080611b51603a5490565b905080611b615750600092915050565b6000611b8084603f60109054906101000a900464ffffffffff16611c59565b9050611b8c8282611c6d565b949350505050565b60007f8b73c3c69bb8fe3d512ecc4cf759cc79239f7b179b0ffacaa9a75d522b39400f611bbf6120f0565b8051602091820120604080518082018252600181527f310000000000000000000000000000000000000000000000000000000000000090840152805192830193909352918101919091527fc89efdaa54c0f20c7adf612882df0950f5a951637e0307cdcb4c672f298b8bc660608201524660808201523060a082015260c00160405160208183030381529060405280519060200120905090565b6000611c668383426120fa565b9392505050565b600081157ffffffffffffffffffffffffffffffffffffffffffe6268e1b017bfe18bffffff83900484111517611ca257600080fd5b506b033b2e3c9fd0803ce800000091026b019d971e4fe8401e74000000010490565b600080600080611d088573ffffffffffffffffffffffffffffffffffffffff166000908152603860205260409020546fffffffffffffffffffffffffffffffff1690565b905080611d2057600080600093509350935050611d42565b6000611d2b86610af2565b90508181611d398282612793565b94509450945050505b9193909250565b633b9aca008181029081048214611d5f57600080fd5b919050565b600081156b033b2e3c9fd0803ce800000060028404190484111715611d8857600080fd5b506b033b2e3c9fd0803ce80000009190910260028204010490565b60006fffffffffffffffffffffffffffffffff821115611e45576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602760248201527f53616665436173743a2076616c756520646f65736e27742066697420696e203160448201527f32382062697473000000000000000000000000000000000000000000000000006064820152608401610739565b5090565b6000611e5483611da3565b73ffffffffffffffffffffffffffffffffffffffff85166000908152603860205260409020549091506fffffffffffffffffffffffffffffffff16611e998282612889565b73ffffffffffffffffffffffffffffffffffffffff868116600090815260386020526040902080547fffffffffffffffffffffffffffffffff00000000000000000000000000000000166fffffffffffffffffffffffffffffffff9390931692909217909155603d5461010090041615611fb357603d546040517f31873e2e00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8781166004830152602482018690526fffffffffffffffffffffffffffffffff84166044830152610100909204909116906331873e2e90606401600060405180830381600087803b158015611f9f57600080fd5b505af1158015610a88573d6000803e3d6000fd5b5050505050565b6000611fc583611da3565b73ffffffffffffffffffffffffffffffffffffffff85166000908152603860205260409020549091506fffffffffffffffffffffffffffffffff16611e9982826128bd565b73ffffffffffffffffffffffffffffffffffffffff808416600090815260366020908152604080832093861683529290529081205461204a908390612793565b73ffffffffffffffffffffffffffffffffffffffff808616600081815260366020908152604080832089861680855292529182902085905560375491519495509216927fda919360433220e13b51e8c211e490d148e61a3bd53de8c097194e458b97f3e1906120bc9086815260200190565b60405180910390a450505050565b80516115b690603b906020840190612241565b80516115b690603c906020840190612241565b6060610ab4610640565b60008061210e64ffffffffff851684612793565b90508061212a576b033b2e3c9fd0803ce8000000915050611c66565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81016000808060028511612160576000612165565b600285035b925066038882915c40006121798a80611c6d565b81612186576121866128ee565b0491506301e13380612198838b611c6d565b816121a5576121a56128ee565b0490506000826121b5868861291d565b6121bf919061291d565b600290049050600082856121d3888a61291d565b6121dd919061291d565b6121e7919061291d565b60069004905080826301e133806121fe8a8f61291d565b612208919061295a565b61221e906b033b2e3c9fd0803ce800000061277b565b612228919061277b565b612232919061277b565b9b9a5050505050505050505050565b82805461224d906126f8565b90600052602060002090601f01602090048101928261226f57600085556122b5565b82601f1061228857805160ff19168380011785556122b5565b828001600101855582156122b5579182015b828111156122b557825182559160200191906001019061229a565b50611e459291505b80821115611e4557600081556001016122bd565b6000815180845260005b818110156122f7576020818501810151868301820152016122db565b81811115612309576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000611c6660208301846122d1565b73ffffffffffffffffffffffffffffffffffffffff8116811461237157600080fd5b50565b8035611d5f8161234f565b6000806040838503121561239257600080fd5b823561239d8161234f565b946020939093013593505050565b803560ff81168114611d5f57600080fd5b600080600080600080600060e0888a0312156123d757600080fd5b87356123e28161234f565b965060208801356123f28161234f565b9550604088013594506060880135935061240e608089016123ab565b925060a0880135915060c0880135905092959891949750929550565b60008060006060848603121561243f57600080fd5b833561244a8161234f565b9250602084013561245a8161234f565b929592945050506040919091013590565b6000806040838503121561247e57600080fd5b82356124898161234f565b915060208301356124998161234f565b809150509250929050565b6000602082840312156124b657600080fd5b8135611c668161234f565b600080600080608085870312156124d757600080fd5b84356124e28161234f565b935060208501356124f28161234f565b93969395505050506040820135916060013590565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600082601f83011261254757600080fd5b813567ffffffffffffffff8082111561256257612562612507565b604051601f83017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f011681019082821181831017156125a8576125a8612507565b816040528381528660208588010111156125c157600080fd5b836020870160208301376000602085830101528094505050505092915050565b60008083601f8401126125f357600080fd5b50813567ffffffffffffffff81111561260b57600080fd5b60208301915083602082850101111561112257600080fd5b60008060008060008060008060e0898b03121561263f57600080fd5b883561264a8161234f565b9750602089013561265a8161234f565b965061266860408a01612374565b955061267660608a016123ab565b9450608089013567ffffffffffffffff8082111561269357600080fd5b61269f8c838d01612536565b955060a08b01359150808211156126b557600080fd5b6126c18c838d01612536565b945060c08b01359150808211156126d757600080fd5b506126e48b828c016125e1565b999c989b5096995094979396929594505050565b600181811c9082168061270c57607f821691505b60208210811415612746577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000821982111561278e5761278e61274c565b500190565b6000828210156127a5576127a561274c565b500390565b73ffffffffffffffffffffffffffffffffffffffff8716815260ff8616602082015260a0604082015260006127e260a08301876122d1565b82810360608401526127f481876122d1565b905082810360808401528381528385602083013760006020858301015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f860116820101915050979650505050505050565b60006020828403121561285c57600080fd5b8151611c668161234f565b60006020828403121561287957600080fd5b81518015158114611c6657600080fd5b60006fffffffffffffffffffffffffffffffff8083168185168083038211156128b4576128b461274c565b01949350505050565b60006fffffffffffffffffffffffffffffffff838116908316818110156128e6576128e661274c565b039392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156129555761295561274c565b500290565b600082612990577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b50049056fea26469706673582212208a96cc31172a645845e36b7b1797b3d710a42ebd5d093ccc7a8accfee126cfa664736f6c634300080a0033",
}

// StableDebtTokenABI is the input ABI used to generate the binding from.
// Deprecated: Use StableDebtTokenMetaData.ABI instead.
var StableDebtTokenABI = StableDebtTokenMetaData.ABI

// StableDebtTokenBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use StableDebtTokenMetaData.Bin instead.
var StableDebtTokenBin = StableDebtTokenMetaData.Bin

// DeployStableDebtToken deploys a new Ethereum contract, binding an instance of StableDebtToken to it.
func DeployStableDebtToken(auth *bind.TransactOpts, backend bind.ContractBackend, pool common.Address) (common.Address, *types.Transaction, *StableDebtToken, error) {
	parsed, err := StableDebtTokenMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StableDebtTokenBin), backend, pool)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &StableDebtToken{StableDebtTokenCaller: StableDebtTokenCaller{contract: contract}, StableDebtTokenTransactor: StableDebtTokenTransactor{contract: contract}, StableDebtTokenFilterer: StableDebtTokenFilterer{contract: contract}}, nil
}

// StableDebtToken is an auto generated Go binding around an Ethereum contract.
type StableDebtToken struct {
	StableDebtTokenCaller     // Read-only binding to the contract
	StableDebtTokenTransactor // Write-only binding to the contract
	StableDebtTokenFilterer   // Log filterer for contract events
}

// StableDebtTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type StableDebtTokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StableDebtTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StableDebtTokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StableDebtTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StableDebtTokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StableDebtTokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StableDebtTokenSession struct {
	Contract     *StableDebtToken  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StableDebtTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StableDebtTokenCallerSession struct {
	Contract *StableDebtTokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// StableDebtTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StableDebtTokenTransactorSession struct {
	Contract     *StableDebtTokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// StableDebtTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type StableDebtTokenRaw struct {
	Contract *StableDebtToken // Generic contract binding to access the raw methods on
}

// StableDebtTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StableDebtTokenCallerRaw struct {
	Contract *StableDebtTokenCaller // Generic read-only contract binding to access the raw methods on
}

// StableDebtTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StableDebtTokenTransactorRaw struct {
	Contract *StableDebtTokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStableDebtToken creates a new instance of StableDebtToken, bound to a specific deployed contract.
func NewStableDebtToken(address common.Address, backend bind.ContractBackend) (*StableDebtToken, error) {
	contract, err := bindStableDebtToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StableDebtToken{StableDebtTokenCaller: StableDebtTokenCaller{contract: contract}, StableDebtTokenTransactor: StableDebtTokenTransactor{contract: contract}, StableDebtTokenFilterer: StableDebtTokenFilterer{contract: contract}}, nil
}

// NewStableDebtTokenCaller creates a new read-only instance of StableDebtToken, bound to a specific deployed contract.
func NewStableDebtTokenCaller(address common.Address, caller bind.ContractCaller) (*StableDebtTokenCaller, error) {
	contract, err := bindStableDebtToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StableDebtTokenCaller{contract: contract}, nil
}

// NewStableDebtTokenTransactor creates a new write-only instance of StableDebtToken, bound to a specific deployed contract.
func NewStableDebtTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*StableDebtTokenTransactor, error) {
	contract, err := bindStableDebtToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StableDebtTokenTransactor{contract: contract}, nil
}

// NewStableDebtTokenFilterer creates a new log filterer instance of StableDebtToken, bound to a specific deployed contract.
func NewStableDebtTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*StableDebtTokenFilterer, error) {
	contract, err := bindStableDebtToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StableDebtTokenFilterer{contract: contract}, nil
}

// bindStableDebtToken binds a generic wrapper to an already deployed contract.
func bindStableDebtToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StableDebtTokenMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StableDebtToken *StableDebtTokenRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StableDebtToken.Contract.StableDebtTokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StableDebtToken *StableDebtTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StableDebtToken.Contract.StableDebtTokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StableDebtToken *StableDebtTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StableDebtToken.Contract.StableDebtTokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StableDebtToken *StableDebtTokenCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StableDebtToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StableDebtToken *StableDebtTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StableDebtToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StableDebtToken *StableDebtTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StableDebtToken.Contract.contract.Transact(opts, method, params...)
}

// DEBTTOKENREVISION is a free data retrieval call binding the contract method 0xb9a7b622.
//
// Solidity: function DEBT_TOKEN_REVISION() view returns(uint256)
func (_StableDebtToken *StableDebtTokenCaller) DEBTTOKENREVISION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StableDebtToken.contract.Call(opts, &out, "DEBT_TOKEN_REVISION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DEBTTOKENREVISION is a free data retrieval call binding the contract method 0xb9a7b622.
//
// Solidity: function DEBT_TOKEN_REVISION() view returns(uint256)
func (_StableDebtToken *StableDebtTokenSession) DEBTTOKENREVISION() (*big.Int, error) {
	return _StableDebtToken.Contract.DEBTTOKENREVISION(&_StableDebtToken.CallOpts)
}

// DEBTTOKENREVISION is a free data retrieval call binding the contract method 0xb9a7b622.
//
// Solidity: function DEBT_TOKEN_REVISION() view returns(uint256)
func (_StableDebtToken *StableDebtTokenCallerSession) DEBTTOKENREVISION() (*big.Int, error) {
	return _StableDebtToken.Contract.DEBTTOKENREVISION(&_StableDebtToken.CallOpts)
}

// DELEGATIONWITHSIGTYPEHASH is a free data retrieval call binding the contract method 0xf3bfc738.
//
// Solidity: function DELEGATION_WITH_SIG_TYPEHASH() view returns(bytes32)
func (_StableDebtToken *StableDebtTokenCaller) DELEGATIONWITHSIGTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _StableDebtToken.contract.Call(opts, &out, "DELEGATION_WITH_SIG_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DELEGATIONWITHSIGTYPEHASH is a free data retrieval call binding the contract method 0xf3bfc738.
//
// Solidity: function DELEGATION_WITH_SIG_TYPEHASH() view returns(bytes32)
func (_StableDebtToken *StableDebtTokenSession) DELEGATIONWITHSIGTYPEHASH() ([32]byte, error) {
	return _StableDebtToken.Contract.DELEGATIONWITHSIGTYPEHASH(&_StableDebtToken.CallOpts)
}

// DELEGATIONWITHSIGTYPEHASH is a free data retrieval call binding the contract method 0xf3bfc738.
//
// Solidity: function DELEGATION_WITH_SIG_TYPEHASH() view returns(bytes32)
func (_StableDebtToken *StableDebtTokenCallerSession) DELEGATIONWITHSIGTYPEHASH() ([32]byte, error) {
	return _StableDebtToken.Contract.DELEGATIONWITHSIGTYPEHASH(&_StableDebtToken.CallOpts)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_StableDebtToken *StableDebtTokenCaller) DOMAINSEPARATOR(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _StableDebtToken.contract.Call(opts, &out, "DOMAIN_SEPARATOR")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_StableDebtToken *StableDebtTokenSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _StableDebtToken.Contract.DOMAINSEPARATOR(&_StableDebtToken.CallOpts)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_StableDebtToken *StableDebtTokenCallerSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _StableDebtToken.Contract.DOMAINSEPARATOR(&_StableDebtToken.CallOpts)
}

// EIP712REVISION is a free data retrieval call binding the contract method 0x78160376.
//
// Solidity: function EIP712_REVISION() view returns(bytes)
func (_StableDebtToken *StableDebtTokenCaller) EIP712REVISION(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _StableDebtToken.contract.Call(opts, &out, "EIP712_REVISION")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// EIP712REVISION is a free data retrieval call binding the contract method 0x78160376.
//
// Solidity: function EIP712_REVISION() view returns(bytes)
func (_StableDebtToken *StableDebtTokenSession) EIP712REVISION() ([]byte, error) {
	return _StableDebtToken.Contract.EIP712REVISION(&_StableDebtToken.CallOpts)
}

// EIP712REVISION is a free data retrieval call binding the contract method 0x78160376.
//
// Solidity: function EIP712_REVISION() view returns(bytes)
func (_StableDebtToken *StableDebtTokenCallerSession) EIP712REVISION() ([]byte, error) {
	return _StableDebtToken.Contract.EIP712REVISION(&_StableDebtToken.CallOpts)
}

// POOL is a free data retrieval call binding the contract method 0x7535d246.
//
// Solidity: function POOL() view returns(address)
func (_StableDebtToken *StableDebtTokenCaller) POOL(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StableDebtToken.contract.Call(opts, &out, "POOL")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// POOL is a free data retrieval call binding the contract method 0x7535d246.
//
// Solidity: function POOL() view returns(address)
func (_StableDebtToken *StableDebtTokenSession) POOL() (common.Address, error) {
	return _StableDebtToken.Contract.POOL(&_StableDebtToken.CallOpts)
}

// POOL is a free data retrieval call binding the contract method 0x7535d246.
//
// Solidity: function POOL() view returns(address)
func (_StableDebtToken *StableDebtTokenCallerSession) POOL() (common.Address, error) {
	return _StableDebtToken.Contract.POOL(&_StableDebtToken.CallOpts)
}

// UNDERLYINGASSETADDRESS is a free data retrieval call binding the contract method 0xb16a19de.
//
// Solidity: function UNDERLYING_ASSET_ADDRESS() view returns(address)
func (_StableDebtToken *StableDebtTokenCaller) UNDERLYINGASSETADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StableDebtToken.contract.Call(opts, &out, "UNDERLYING_ASSET_ADDRESS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UNDERLYINGASSETADDRESS is a free data retrieval call binding the contract method 0xb16a19de.
//
// Solidity: function UNDERLYING_ASSET_ADDRESS() view returns(address)
func (_StableDebtToken *StableDebtTokenSession) UNDERLYINGASSETADDRESS() (common.Address, error) {
	return _StableDebtToken.Contract.UNDERLYINGASSETADDRESS(&_StableDebtToken.CallOpts)
}

// UNDERLYINGASSETADDRESS is a free data retrieval call binding the contract method 0xb16a19de.
//
// Solidity: function UNDERLYING_ASSET_ADDRESS() view returns(address)
func (_StableDebtToken *StableDebtTokenCallerSession) UNDERLYINGASSETADDRESS() (common.Address, error) {
	return _StableDebtToken.Contract.UNDERLYINGASSETADDRESS(&_StableDebtToken.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_StableDebtToken *StableDebtTokenCaller) Allowance(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _StableDebtToken.contract.Call(opts, &out, "allowance", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_StableDebtToken *StableDebtTokenSession) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _StableDebtToken.Contract.Allowance(&_StableDebtToken.CallOpts, arg0, arg1)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_StableDebtToken *StableDebtTokenCallerSession) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _StableDebtToken.Contract.Allowance(&_StableDebtToken.CallOpts, arg0, arg1)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_StableDebtToken *StableDebtTokenCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _StableDebtToken.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_StableDebtToken *StableDebtTokenSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _StableDebtToken.Contract.BalanceOf(&_StableDebtToken.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_StableDebtToken *StableDebtTokenCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _StableDebtToken.Contract.BalanceOf(&_StableDebtToken.CallOpts, account)
}

// BorrowAllowance is a free data retrieval call binding the contract method 0x6bd76d24.
//
// Solidity: function borrowAllowance(address fromUser, address toUser) view returns(uint256)
func (_StableDebtToken *StableDebtTokenCaller) BorrowAllowance(opts *bind.CallOpts, fromUser common.Address, toUser common.Address) (*big.Int, error) {
	var out []interface{}
	err := _StableDebtToken.contract.Call(opts, &out, "borrowAllowance", fromUser, toUser)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BorrowAllowance is a free data retrieval call binding the contract method 0x6bd76d24.
//
// Solidity: function borrowAllowance(address fromUser, address toUser) view returns(uint256)
func (_StableDebtToken *StableDebtTokenSession) BorrowAllowance(fromUser common.Address, toUser common.Address) (*big.Int, error) {
	return _StableDebtToken.Contract.BorrowAllowance(&_StableDebtToken.CallOpts, fromUser, toUser)
}

// BorrowAllowance is a free data retrieval call binding the contract method 0x6bd76d24.
//
// Solidity: function borrowAllowance(address fromUser, address toUser) view returns(uint256)
func (_StableDebtToken *StableDebtTokenCallerSession) BorrowAllowance(fromUser common.Address, toUser common.Address) (*big.Int, error) {
	return _StableDebtToken.Contract.BorrowAllowance(&_StableDebtToken.CallOpts, fromUser, toUser)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_StableDebtToken *StableDebtTokenCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _StableDebtToken.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_StableDebtToken *StableDebtTokenSession) Decimals() (uint8, error) {
	return _StableDebtToken.Contract.Decimals(&_StableDebtToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_StableDebtToken *StableDebtTokenCallerSession) Decimals() (uint8, error) {
	return _StableDebtToken.Contract.Decimals(&_StableDebtToken.CallOpts)
}

// GetAverageStableRate is a free data retrieval call binding the contract method 0x90f6fcf2.
//
// Solidity: function getAverageStableRate() view returns(uint256)
func (_StableDebtToken *StableDebtTokenCaller) GetAverageStableRate(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StableDebtToken.contract.Call(opts, &out, "getAverageStableRate")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAverageStableRate is a free data retrieval call binding the contract method 0x90f6fcf2.
//
// Solidity: function getAverageStableRate() view returns(uint256)
func (_StableDebtToken *StableDebtTokenSession) GetAverageStableRate() (*big.Int, error) {
	return _StableDebtToken.Contract.GetAverageStableRate(&_StableDebtToken.CallOpts)
}

// GetAverageStableRate is a free data retrieval call binding the contract method 0x90f6fcf2.
//
// Solidity: function getAverageStableRate() view returns(uint256)
func (_StableDebtToken *StableDebtTokenCallerSession) GetAverageStableRate() (*big.Int, error) {
	return _StableDebtToken.Contract.GetAverageStableRate(&_StableDebtToken.CallOpts)
}

// GetIncentivesController is a free data retrieval call binding the contract method 0x75d26413.
//
// Solidity: function getIncentivesController() view returns(address)
func (_StableDebtToken *StableDebtTokenCaller) GetIncentivesController(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StableDebtToken.contract.Call(opts, &out, "getIncentivesController")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetIncentivesController is a free data retrieval call binding the contract method 0x75d26413.
//
// Solidity: function getIncentivesController() view returns(address)
func (_StableDebtToken *StableDebtTokenSession) GetIncentivesController() (common.Address, error) {
	return _StableDebtToken.Contract.GetIncentivesController(&_StableDebtToken.CallOpts)
}

// GetIncentivesController is a free data retrieval call binding the contract method 0x75d26413.
//
// Solidity: function getIncentivesController() view returns(address)
func (_StableDebtToken *StableDebtTokenCallerSession) GetIncentivesController() (common.Address, error) {
	return _StableDebtToken.Contract.GetIncentivesController(&_StableDebtToken.CallOpts)
}

// GetSupplyData is a free data retrieval call binding the contract method 0x79774338.
//
// Solidity: function getSupplyData() view returns(uint256, uint256, uint256, uint40)
func (_StableDebtToken *StableDebtTokenCaller) GetSupplyData(opts *bind.CallOpts) (*big.Int, *big.Int, *big.Int, *big.Int, error) {
	var out []interface{}
	err := _StableDebtToken.contract.Call(opts, &out, "getSupplyData")

	if err != nil {
		return *new(*big.Int), *new(*big.Int), *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	out2 := *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	out3 := *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return out0, out1, out2, out3, err

}

// GetSupplyData is a free data retrieval call binding the contract method 0x79774338.
//
// Solidity: function getSupplyData() view returns(uint256, uint256, uint256, uint40)
func (_StableDebtToken *StableDebtTokenSession) GetSupplyData() (*big.Int, *big.Int, *big.Int, *big.Int, error) {
	return _StableDebtToken.Contract.GetSupplyData(&_StableDebtToken.CallOpts)
}

// GetSupplyData is a free data retrieval call binding the contract method 0x79774338.
//
// Solidity: function getSupplyData() view returns(uint256, uint256, uint256, uint40)
func (_StableDebtToken *StableDebtTokenCallerSession) GetSupplyData() (*big.Int, *big.Int, *big.Int, *big.Int, error) {
	return _StableDebtToken.Contract.GetSupplyData(&_StableDebtToken.CallOpts)
}

// GetTotalSupplyAndAvgRate is a free data retrieval call binding the contract method 0xf731e9be.
//
// Solidity: function getTotalSupplyAndAvgRate() view returns(uint256, uint256)
func (_StableDebtToken *StableDebtTokenCaller) GetTotalSupplyAndAvgRate(opts *bind.CallOpts) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _StableDebtToken.contract.Call(opts, &out, "getTotalSupplyAndAvgRate")

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// GetTotalSupplyAndAvgRate is a free data retrieval call binding the contract method 0xf731e9be.
//
// Solidity: function getTotalSupplyAndAvgRate() view returns(uint256, uint256)
func (_StableDebtToken *StableDebtTokenSession) GetTotalSupplyAndAvgRate() (*big.Int, *big.Int, error) {
	return _StableDebtToken.Contract.GetTotalSupplyAndAvgRate(&_StableDebtToken.CallOpts)
}

// GetTotalSupplyAndAvgRate is a free data retrieval call binding the contract method 0xf731e9be.
//
// Solidity: function getTotalSupplyAndAvgRate() view returns(uint256, uint256)
func (_StableDebtToken *StableDebtTokenCallerSession) GetTotalSupplyAndAvgRate() (*big.Int, *big.Int, error) {
	return _StableDebtToken.Contract.GetTotalSupplyAndAvgRate(&_StableDebtToken.CallOpts)
}

// GetTotalSupplyLastUpdated is a free data retrieval call binding the contract method 0xe7484890.
//
// Solidity: function getTotalSupplyLastUpdated() view returns(uint40)
func (_StableDebtToken *StableDebtTokenCaller) GetTotalSupplyLastUpdated(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StableDebtToken.contract.Call(opts, &out, "getTotalSupplyLastUpdated")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTotalSupplyLastUpdated is a free data retrieval call binding the contract method 0xe7484890.
//
// Solidity: function getTotalSupplyLastUpdated() view returns(uint40)
func (_StableDebtToken *StableDebtTokenSession) GetTotalSupplyLastUpdated() (*big.Int, error) {
	return _StableDebtToken.Contract.GetTotalSupplyLastUpdated(&_StableDebtToken.CallOpts)
}

// GetTotalSupplyLastUpdated is a free data retrieval call binding the contract method 0xe7484890.
//
// Solidity: function getTotalSupplyLastUpdated() view returns(uint40)
func (_StableDebtToken *StableDebtTokenCallerSession) GetTotalSupplyLastUpdated() (*big.Int, error) {
	return _StableDebtToken.Contract.GetTotalSupplyLastUpdated(&_StableDebtToken.CallOpts)
}

// GetUserLastUpdated is a free data retrieval call binding the contract method 0x79ce6b8c.
//
// Solidity: function getUserLastUpdated(address user) view returns(uint40)
func (_StableDebtToken *StableDebtTokenCaller) GetUserLastUpdated(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _StableDebtToken.contract.Call(opts, &out, "getUserLastUpdated", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetUserLastUpdated is a free data retrieval call binding the contract method 0x79ce6b8c.
//
// Solidity: function getUserLastUpdated(address user) view returns(uint40)
func (_StableDebtToken *StableDebtTokenSession) GetUserLastUpdated(user common.Address) (*big.Int, error) {
	return _StableDebtToken.Contract.GetUserLastUpdated(&_StableDebtToken.CallOpts, user)
}

// GetUserLastUpdated is a free data retrieval call binding the contract method 0x79ce6b8c.
//
// Solidity: function getUserLastUpdated(address user) view returns(uint40)
func (_StableDebtToken *StableDebtTokenCallerSession) GetUserLastUpdated(user common.Address) (*big.Int, error) {
	return _StableDebtToken.Contract.GetUserLastUpdated(&_StableDebtToken.CallOpts, user)
}

// GetUserStableRate is a free data retrieval call binding the contract method 0xe78c9b3b.
//
// Solidity: function getUserStableRate(address user) view returns(uint256)
func (_StableDebtToken *StableDebtTokenCaller) GetUserStableRate(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _StableDebtToken.contract.Call(opts, &out, "getUserStableRate", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetUserStableRate is a free data retrieval call binding the contract method 0xe78c9b3b.
//
// Solidity: function getUserStableRate(address user) view returns(uint256)
func (_StableDebtToken *StableDebtTokenSession) GetUserStableRate(user common.Address) (*big.Int, error) {
	return _StableDebtToken.Contract.GetUserStableRate(&_StableDebtToken.CallOpts, user)
}

// GetUserStableRate is a free data retrieval call binding the contract method 0xe78c9b3b.
//
// Solidity: function getUserStableRate(address user) view returns(uint256)
func (_StableDebtToken *StableDebtTokenCallerSession) GetUserStableRate(user common.Address) (*big.Int, error) {
	return _StableDebtToken.Contract.GetUserStableRate(&_StableDebtToken.CallOpts, user)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_StableDebtToken *StableDebtTokenCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _StableDebtToken.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_StableDebtToken *StableDebtTokenSession) Name() (string, error) {
	return _StableDebtToken.Contract.Name(&_StableDebtToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_StableDebtToken *StableDebtTokenCallerSession) Name() (string, error) {
	return _StableDebtToken.Contract.Name(&_StableDebtToken.CallOpts)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_StableDebtToken *StableDebtTokenCaller) Nonces(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _StableDebtToken.contract.Call(opts, &out, "nonces", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_StableDebtToken *StableDebtTokenSession) Nonces(owner common.Address) (*big.Int, error) {
	return _StableDebtToken.Contract.Nonces(&_StableDebtToken.CallOpts, owner)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_StableDebtToken *StableDebtTokenCallerSession) Nonces(owner common.Address) (*big.Int, error) {
	return _StableDebtToken.Contract.Nonces(&_StableDebtToken.CallOpts, owner)
}

// PrincipalBalanceOf is a free data retrieval call binding the contract method 0xc634dfaa.
//
// Solidity: function principalBalanceOf(address user) view returns(uint256)
func (_StableDebtToken *StableDebtTokenCaller) PrincipalBalanceOf(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _StableDebtToken.contract.Call(opts, &out, "principalBalanceOf", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PrincipalBalanceOf is a free data retrieval call binding the contract method 0xc634dfaa.
//
// Solidity: function principalBalanceOf(address user) view returns(uint256)
func (_StableDebtToken *StableDebtTokenSession) PrincipalBalanceOf(user common.Address) (*big.Int, error) {
	return _StableDebtToken.Contract.PrincipalBalanceOf(&_StableDebtToken.CallOpts, user)
}

// PrincipalBalanceOf is a free data retrieval call binding the contract method 0xc634dfaa.
//
// Solidity: function principalBalanceOf(address user) view returns(uint256)
func (_StableDebtToken *StableDebtTokenCallerSession) PrincipalBalanceOf(user common.Address) (*big.Int, error) {
	return _StableDebtToken.Contract.PrincipalBalanceOf(&_StableDebtToken.CallOpts, user)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_StableDebtToken *StableDebtTokenCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _StableDebtToken.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_StableDebtToken *StableDebtTokenSession) Symbol() (string, error) {
	return _StableDebtToken.Contract.Symbol(&_StableDebtToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_StableDebtToken *StableDebtTokenCallerSession) Symbol() (string, error) {
	return _StableDebtToken.Contract.Symbol(&_StableDebtToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_StableDebtToken *StableDebtTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StableDebtToken.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_StableDebtToken *StableDebtTokenSession) TotalSupply() (*big.Int, error) {
	return _StableDebtToken.Contract.TotalSupply(&_StableDebtToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_StableDebtToken *StableDebtTokenCallerSession) TotalSupply() (*big.Int, error) {
	return _StableDebtToken.Contract.TotalSupply(&_StableDebtToken.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address , uint256 ) returns(bool)
func (_StableDebtToken *StableDebtTokenTransactor) Approve(opts *bind.TransactOpts, arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _StableDebtToken.contract.Transact(opts, "approve", arg0, arg1)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address , uint256 ) returns(bool)
func (_StableDebtToken *StableDebtTokenSession) Approve(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _StableDebtToken.Contract.Approve(&_StableDebtToken.TransactOpts, arg0, arg1)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address , uint256 ) returns(bool)
func (_StableDebtToken *StableDebtTokenTransactorSession) Approve(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _StableDebtToken.Contract.Approve(&_StableDebtToken.TransactOpts, arg0, arg1)
}

// ApproveDelegation is a paid mutator transaction binding the contract method 0xc04a8a10.
//
// Solidity: function approveDelegation(address delegatee, uint256 amount) returns()
func (_StableDebtToken *StableDebtTokenTransactor) ApproveDelegation(opts *bind.TransactOpts, delegatee common.Address, amount *big.Int) (*types.Transaction, error) {
	return _StableDebtToken.contract.Transact(opts, "approveDelegation", delegatee, amount)
}

// ApproveDelegation is a paid mutator transaction binding the contract method 0xc04a8a10.
//
// Solidity: function approveDelegation(address delegatee, uint256 amount) returns()
func (_StableDebtToken *StableDebtTokenSession) ApproveDelegation(delegatee common.Address, amount *big.Int) (*types.Transaction, error) {
	return _StableDebtToken.Contract.ApproveDelegation(&_StableDebtToken.TransactOpts, delegatee, amount)
}

// ApproveDelegation is a paid mutator transaction binding the contract method 0xc04a8a10.
//
// Solidity: function approveDelegation(address delegatee, uint256 amount) returns()
func (_StableDebtToken *StableDebtTokenTransactorSession) ApproveDelegation(delegatee common.Address, amount *big.Int) (*types.Transaction, error) {
	return _StableDebtToken.Contract.ApproveDelegation(&_StableDebtToken.TransactOpts, delegatee, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address from, uint256 amount) returns(uint256, uint256)
func (_StableDebtToken *StableDebtTokenTransactor) Burn(opts *bind.TransactOpts, from common.Address, amount *big.Int) (*types.Transaction, error) {
	return _StableDebtToken.contract.Transact(opts, "burn", from, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address from, uint256 amount) returns(uint256, uint256)
func (_StableDebtToken *StableDebtTokenSession) Burn(from common.Address, amount *big.Int) (*types.Transaction, error) {
	return _StableDebtToken.Contract.Burn(&_StableDebtToken.TransactOpts, from, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address from, uint256 amount) returns(uint256, uint256)
func (_StableDebtToken *StableDebtTokenTransactorSession) Burn(from common.Address, amount *big.Int) (*types.Transaction, error) {
	return _StableDebtToken.Contract.Burn(&_StableDebtToken.TransactOpts, from, amount)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address , uint256 ) returns(bool)
func (_StableDebtToken *StableDebtTokenTransactor) DecreaseAllowance(opts *bind.TransactOpts, arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _StableDebtToken.contract.Transact(opts, "decreaseAllowance", arg0, arg1)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address , uint256 ) returns(bool)
func (_StableDebtToken *StableDebtTokenSession) DecreaseAllowance(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _StableDebtToken.Contract.DecreaseAllowance(&_StableDebtToken.TransactOpts, arg0, arg1)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address , uint256 ) returns(bool)
func (_StableDebtToken *StableDebtTokenTransactorSession) DecreaseAllowance(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _StableDebtToken.Contract.DecreaseAllowance(&_StableDebtToken.TransactOpts, arg0, arg1)
}

// DelegationWithSig is a paid mutator transaction binding the contract method 0x0b52d558.
//
// Solidity: function delegationWithSig(address delegator, address delegatee, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_StableDebtToken *StableDebtTokenTransactor) DelegationWithSig(opts *bind.TransactOpts, delegator common.Address, delegatee common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _StableDebtToken.contract.Transact(opts, "delegationWithSig", delegator, delegatee, value, deadline, v, r, s)
}

// DelegationWithSig is a paid mutator transaction binding the contract method 0x0b52d558.
//
// Solidity: function delegationWithSig(address delegator, address delegatee, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_StableDebtToken *StableDebtTokenSession) DelegationWithSig(delegator common.Address, delegatee common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _StableDebtToken.Contract.DelegationWithSig(&_StableDebtToken.TransactOpts, delegator, delegatee, value, deadline, v, r, s)
}

// DelegationWithSig is a paid mutator transaction binding the contract method 0x0b52d558.
//
// Solidity: function delegationWithSig(address delegator, address delegatee, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_StableDebtToken *StableDebtTokenTransactorSession) DelegationWithSig(delegator common.Address, delegatee common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _StableDebtToken.Contract.DelegationWithSig(&_StableDebtToken.TransactOpts, delegator, delegatee, value, deadline, v, r, s)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address , uint256 ) returns(bool)
func (_StableDebtToken *StableDebtTokenTransactor) IncreaseAllowance(opts *bind.TransactOpts, arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _StableDebtToken.contract.Transact(opts, "increaseAllowance", arg0, arg1)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address , uint256 ) returns(bool)
func (_StableDebtToken *StableDebtTokenSession) IncreaseAllowance(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _StableDebtToken.Contract.IncreaseAllowance(&_StableDebtToken.TransactOpts, arg0, arg1)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address , uint256 ) returns(bool)
func (_StableDebtToken *StableDebtTokenTransactorSession) IncreaseAllowance(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _StableDebtToken.Contract.IncreaseAllowance(&_StableDebtToken.TransactOpts, arg0, arg1)
}

// Initialize is a paid mutator transaction binding the contract method 0xc222ec8a.
//
// Solidity: function initialize(address initializingPool, address underlyingAsset, address incentivesController, uint8 debtTokenDecimals, string debtTokenName, string debtTokenSymbol, bytes params) returns()
func (_StableDebtToken *StableDebtTokenTransactor) Initialize(opts *bind.TransactOpts, initializingPool common.Address, underlyingAsset common.Address, incentivesController common.Address, debtTokenDecimals uint8, debtTokenName string, debtTokenSymbol string, params []byte) (*types.Transaction, error) {
	return _StableDebtToken.contract.Transact(opts, "initialize", initializingPool, underlyingAsset, incentivesController, debtTokenDecimals, debtTokenName, debtTokenSymbol, params)
}

// Initialize is a paid mutator transaction binding the contract method 0xc222ec8a.
//
// Solidity: function initialize(address initializingPool, address underlyingAsset, address incentivesController, uint8 debtTokenDecimals, string debtTokenName, string debtTokenSymbol, bytes params) returns()
func (_StableDebtToken *StableDebtTokenSession) Initialize(initializingPool common.Address, underlyingAsset common.Address, incentivesController common.Address, debtTokenDecimals uint8, debtTokenName string, debtTokenSymbol string, params []byte) (*types.Transaction, error) {
	return _StableDebtToken.Contract.Initialize(&_StableDebtToken.TransactOpts, initializingPool, underlyingAsset, incentivesController, debtTokenDecimals, debtTokenName, debtTokenSymbol, params)
}

// Initialize is a paid mutator transaction binding the contract method 0xc222ec8a.
//
// Solidity: function initialize(address initializingPool, address underlyingAsset, address incentivesController, uint8 debtTokenDecimals, string debtTokenName, string debtTokenSymbol, bytes params) returns()
func (_StableDebtToken *StableDebtTokenTransactorSession) Initialize(initializingPool common.Address, underlyingAsset common.Address, incentivesController common.Address, debtTokenDecimals uint8, debtTokenName string, debtTokenSymbol string, params []byte) (*types.Transaction, error) {
	return _StableDebtToken.Contract.Initialize(&_StableDebtToken.TransactOpts, initializingPool, underlyingAsset, incentivesController, debtTokenDecimals, debtTokenName, debtTokenSymbol, params)
}

// Mint is a paid mutator transaction binding the contract method 0xb3f1c93d.
//
// Solidity: function mint(address user, address onBehalfOf, uint256 amount, uint256 rate) returns(bool, uint256, uint256)
func (_StableDebtToken *StableDebtTokenTransactor) Mint(opts *bind.TransactOpts, user common.Address, onBehalfOf common.Address, amount *big.Int, rate *big.Int) (*types.Transaction, error) {
	return _StableDebtToken.contract.Transact(opts, "mint", user, onBehalfOf, amount, rate)
}

// Mint is a paid mutator transaction binding the contract method 0xb3f1c93d.
//
// Solidity: function mint(address user, address onBehalfOf, uint256 amount, uint256 rate) returns(bool, uint256, uint256)
func (_StableDebtToken *StableDebtTokenSession) Mint(user common.Address, onBehalfOf common.Address, amount *big.Int, rate *big.Int) (*types.Transaction, error) {
	return _StableDebtToken.Contract.Mint(&_StableDebtToken.TransactOpts, user, onBehalfOf, amount, rate)
}

// Mint is a paid mutator transaction binding the contract method 0xb3f1c93d.
//
// Solidity: function mint(address user, address onBehalfOf, uint256 amount, uint256 rate) returns(bool, uint256, uint256)
func (_StableDebtToken *StableDebtTokenTransactorSession) Mint(user common.Address, onBehalfOf common.Address, amount *big.Int, rate *big.Int) (*types.Transaction, error) {
	return _StableDebtToken.Contract.Mint(&_StableDebtToken.TransactOpts, user, onBehalfOf, amount, rate)
}

// SetIncentivesController is a paid mutator transaction binding the contract method 0xe655dbd8.
//
// Solidity: function setIncentivesController(address controller) returns()
func (_StableDebtToken *StableDebtTokenTransactor) SetIncentivesController(opts *bind.TransactOpts, controller common.Address) (*types.Transaction, error) {
	return _StableDebtToken.contract.Transact(opts, "setIncentivesController", controller)
}

// SetIncentivesController is a paid mutator transaction binding the contract method 0xe655dbd8.
//
// Solidity: function setIncentivesController(address controller) returns()
func (_StableDebtToken *StableDebtTokenSession) SetIncentivesController(controller common.Address) (*types.Transaction, error) {
	return _StableDebtToken.Contract.SetIncentivesController(&_StableDebtToken.TransactOpts, controller)
}

// SetIncentivesController is a paid mutator transaction binding the contract method 0xe655dbd8.
//
// Solidity: function setIncentivesController(address controller) returns()
func (_StableDebtToken *StableDebtTokenTransactorSession) SetIncentivesController(controller common.Address) (*types.Transaction, error) {
	return _StableDebtToken.Contract.SetIncentivesController(&_StableDebtToken.TransactOpts, controller)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address , uint256 ) returns(bool)
func (_StableDebtToken *StableDebtTokenTransactor) Transfer(opts *bind.TransactOpts, arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _StableDebtToken.contract.Transact(opts, "transfer", arg0, arg1)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address , uint256 ) returns(bool)
func (_StableDebtToken *StableDebtTokenSession) Transfer(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _StableDebtToken.Contract.Transfer(&_StableDebtToken.TransactOpts, arg0, arg1)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address , uint256 ) returns(bool)
func (_StableDebtToken *StableDebtTokenTransactorSession) Transfer(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _StableDebtToken.Contract.Transfer(&_StableDebtToken.TransactOpts, arg0, arg1)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address , address , uint256 ) returns(bool)
func (_StableDebtToken *StableDebtTokenTransactor) TransferFrom(opts *bind.TransactOpts, arg0 common.Address, arg1 common.Address, arg2 *big.Int) (*types.Transaction, error) {
	return _StableDebtToken.contract.Transact(opts, "transferFrom", arg0, arg1, arg2)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address , address , uint256 ) returns(bool)
func (_StableDebtToken *StableDebtTokenSession) TransferFrom(arg0 common.Address, arg1 common.Address, arg2 *big.Int) (*types.Transaction, error) {
	return _StableDebtToken.Contract.TransferFrom(&_StableDebtToken.TransactOpts, arg0, arg1, arg2)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address , address , uint256 ) returns(bool)
func (_StableDebtToken *StableDebtTokenTransactorSession) TransferFrom(arg0 common.Address, arg1 common.Address, arg2 *big.Int) (*types.Transaction, error) {
	return _StableDebtToken.Contract.TransferFrom(&_StableDebtToken.TransactOpts, arg0, arg1, arg2)
}

// StableDebtTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the StableDebtToken contract.
type StableDebtTokenApprovalIterator struct {
	Event *StableDebtTokenApproval // Event containing the contract specifics and raw log

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
func (it *StableDebtTokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StableDebtTokenApproval)
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
		it.Event = new(StableDebtTokenApproval)
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
func (it *StableDebtTokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StableDebtTokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StableDebtTokenApproval represents a Approval event raised by the StableDebtToken contract.
type StableDebtTokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_StableDebtToken *StableDebtTokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*StableDebtTokenApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _StableDebtToken.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &StableDebtTokenApprovalIterator{contract: _StableDebtToken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_StableDebtToken *StableDebtTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *StableDebtTokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _StableDebtToken.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StableDebtTokenApproval)
				if err := _StableDebtToken.contract.UnpackLog(event, "Approval", log); err != nil {
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
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_StableDebtToken *StableDebtTokenFilterer) ParseApproval(log types.Log) (*StableDebtTokenApproval, error) {
	event := new(StableDebtTokenApproval)
	if err := _StableDebtToken.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StableDebtTokenBorrowAllowanceDelegatedIterator is returned from FilterBorrowAllowanceDelegated and is used to iterate over the raw logs and unpacked data for BorrowAllowanceDelegated events raised by the StableDebtToken contract.
type StableDebtTokenBorrowAllowanceDelegatedIterator struct {
	Event *StableDebtTokenBorrowAllowanceDelegated // Event containing the contract specifics and raw log

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
func (it *StableDebtTokenBorrowAllowanceDelegatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StableDebtTokenBorrowAllowanceDelegated)
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
		it.Event = new(StableDebtTokenBorrowAllowanceDelegated)
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
func (it *StableDebtTokenBorrowAllowanceDelegatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StableDebtTokenBorrowAllowanceDelegatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StableDebtTokenBorrowAllowanceDelegated represents a BorrowAllowanceDelegated event raised by the StableDebtToken contract.
type StableDebtTokenBorrowAllowanceDelegated struct {
	FromUser common.Address
	ToUser   common.Address
	Asset    common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterBorrowAllowanceDelegated is a free log retrieval operation binding the contract event 0xda919360433220e13b51e8c211e490d148e61a3bd53de8c097194e458b97f3e1.
//
// Solidity: event BorrowAllowanceDelegated(address indexed fromUser, address indexed toUser, address indexed asset, uint256 amount)
func (_StableDebtToken *StableDebtTokenFilterer) FilterBorrowAllowanceDelegated(opts *bind.FilterOpts, fromUser []common.Address, toUser []common.Address, asset []common.Address) (*StableDebtTokenBorrowAllowanceDelegatedIterator, error) {

	var fromUserRule []interface{}
	for _, fromUserItem := range fromUser {
		fromUserRule = append(fromUserRule, fromUserItem)
	}
	var toUserRule []interface{}
	for _, toUserItem := range toUser {
		toUserRule = append(toUserRule, toUserItem)
	}
	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _StableDebtToken.contract.FilterLogs(opts, "BorrowAllowanceDelegated", fromUserRule, toUserRule, assetRule)
	if err != nil {
		return nil, err
	}
	return &StableDebtTokenBorrowAllowanceDelegatedIterator{contract: _StableDebtToken.contract, event: "BorrowAllowanceDelegated", logs: logs, sub: sub}, nil
}

// WatchBorrowAllowanceDelegated is a free log subscription operation binding the contract event 0xda919360433220e13b51e8c211e490d148e61a3bd53de8c097194e458b97f3e1.
//
// Solidity: event BorrowAllowanceDelegated(address indexed fromUser, address indexed toUser, address indexed asset, uint256 amount)
func (_StableDebtToken *StableDebtTokenFilterer) WatchBorrowAllowanceDelegated(opts *bind.WatchOpts, sink chan<- *StableDebtTokenBorrowAllowanceDelegated, fromUser []common.Address, toUser []common.Address, asset []common.Address) (event.Subscription, error) {

	var fromUserRule []interface{}
	for _, fromUserItem := range fromUser {
		fromUserRule = append(fromUserRule, fromUserItem)
	}
	var toUserRule []interface{}
	for _, toUserItem := range toUser {
		toUserRule = append(toUserRule, toUserItem)
	}
	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _StableDebtToken.contract.WatchLogs(opts, "BorrowAllowanceDelegated", fromUserRule, toUserRule, assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StableDebtTokenBorrowAllowanceDelegated)
				if err := _StableDebtToken.contract.UnpackLog(event, "BorrowAllowanceDelegated", log); err != nil {
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

// ParseBorrowAllowanceDelegated is a log parse operation binding the contract event 0xda919360433220e13b51e8c211e490d148e61a3bd53de8c097194e458b97f3e1.
//
// Solidity: event BorrowAllowanceDelegated(address indexed fromUser, address indexed toUser, address indexed asset, uint256 amount)
func (_StableDebtToken *StableDebtTokenFilterer) ParseBorrowAllowanceDelegated(log types.Log) (*StableDebtTokenBorrowAllowanceDelegated, error) {
	event := new(StableDebtTokenBorrowAllowanceDelegated)
	if err := _StableDebtToken.contract.UnpackLog(event, "BorrowAllowanceDelegated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StableDebtTokenBurnIterator is returned from FilterBurn and is used to iterate over the raw logs and unpacked data for Burn events raised by the StableDebtToken contract.
type StableDebtTokenBurnIterator struct {
	Event *StableDebtTokenBurn // Event containing the contract specifics and raw log

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
func (it *StableDebtTokenBurnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StableDebtTokenBurn)
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
		it.Event = new(StableDebtTokenBurn)
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
func (it *StableDebtTokenBurnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StableDebtTokenBurnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StableDebtTokenBurn represents a Burn event raised by the StableDebtToken contract.
type StableDebtTokenBurn struct {
	From            common.Address
	Amount          *big.Int
	CurrentBalance  *big.Int
	BalanceIncrease *big.Int
	AvgStableRate   *big.Int
	NewTotalSupply  *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterBurn is a free log retrieval operation binding the contract event 0x44bd20a79e993bdcc7cbedf54a3b4d19fb78490124b6b90d04fe3242eea579e8.
//
// Solidity: event Burn(address indexed from, uint256 amount, uint256 currentBalance, uint256 balanceIncrease, uint256 avgStableRate, uint256 newTotalSupply)
func (_StableDebtToken *StableDebtTokenFilterer) FilterBurn(opts *bind.FilterOpts, from []common.Address) (*StableDebtTokenBurnIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _StableDebtToken.contract.FilterLogs(opts, "Burn", fromRule)
	if err != nil {
		return nil, err
	}
	return &StableDebtTokenBurnIterator{contract: _StableDebtToken.contract, event: "Burn", logs: logs, sub: sub}, nil
}

// WatchBurn is a free log subscription operation binding the contract event 0x44bd20a79e993bdcc7cbedf54a3b4d19fb78490124b6b90d04fe3242eea579e8.
//
// Solidity: event Burn(address indexed from, uint256 amount, uint256 currentBalance, uint256 balanceIncrease, uint256 avgStableRate, uint256 newTotalSupply)
func (_StableDebtToken *StableDebtTokenFilterer) WatchBurn(opts *bind.WatchOpts, sink chan<- *StableDebtTokenBurn, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _StableDebtToken.contract.WatchLogs(opts, "Burn", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StableDebtTokenBurn)
				if err := _StableDebtToken.contract.UnpackLog(event, "Burn", log); err != nil {
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

// ParseBurn is a log parse operation binding the contract event 0x44bd20a79e993bdcc7cbedf54a3b4d19fb78490124b6b90d04fe3242eea579e8.
//
// Solidity: event Burn(address indexed from, uint256 amount, uint256 currentBalance, uint256 balanceIncrease, uint256 avgStableRate, uint256 newTotalSupply)
func (_StableDebtToken *StableDebtTokenFilterer) ParseBurn(log types.Log) (*StableDebtTokenBurn, error) {
	event := new(StableDebtTokenBurn)
	if err := _StableDebtToken.contract.UnpackLog(event, "Burn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StableDebtTokenInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the StableDebtToken contract.
type StableDebtTokenInitializedIterator struct {
	Event *StableDebtTokenInitialized // Event containing the contract specifics and raw log

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
func (it *StableDebtTokenInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StableDebtTokenInitialized)
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
		it.Event = new(StableDebtTokenInitialized)
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
func (it *StableDebtTokenInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StableDebtTokenInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StableDebtTokenInitialized represents a Initialized event raised by the StableDebtToken contract.
type StableDebtTokenInitialized struct {
	UnderlyingAsset      common.Address
	Pool                 common.Address
	IncentivesController common.Address
	DebtTokenDecimals    uint8
	DebtTokenName        string
	DebtTokenSymbol      string
	Params               []byte
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x40251fbfb6656cfa65a00d7879029fec1fad21d28fdcff2f4f68f52795b74f2c.
//
// Solidity: event Initialized(address indexed underlyingAsset, address indexed pool, address incentivesController, uint8 debtTokenDecimals, string debtTokenName, string debtTokenSymbol, bytes params)
func (_StableDebtToken *StableDebtTokenFilterer) FilterInitialized(opts *bind.FilterOpts, underlyingAsset []common.Address, pool []common.Address) (*StableDebtTokenInitializedIterator, error) {

	var underlyingAssetRule []interface{}
	for _, underlyingAssetItem := range underlyingAsset {
		underlyingAssetRule = append(underlyingAssetRule, underlyingAssetItem)
	}
	var poolRule []interface{}
	for _, poolItem := range pool {
		poolRule = append(poolRule, poolItem)
	}

	logs, sub, err := _StableDebtToken.contract.FilterLogs(opts, "Initialized", underlyingAssetRule, poolRule)
	if err != nil {
		return nil, err
	}
	return &StableDebtTokenInitializedIterator{contract: _StableDebtToken.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x40251fbfb6656cfa65a00d7879029fec1fad21d28fdcff2f4f68f52795b74f2c.
//
// Solidity: event Initialized(address indexed underlyingAsset, address indexed pool, address incentivesController, uint8 debtTokenDecimals, string debtTokenName, string debtTokenSymbol, bytes params)
func (_StableDebtToken *StableDebtTokenFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *StableDebtTokenInitialized, underlyingAsset []common.Address, pool []common.Address) (event.Subscription, error) {

	var underlyingAssetRule []interface{}
	for _, underlyingAssetItem := range underlyingAsset {
		underlyingAssetRule = append(underlyingAssetRule, underlyingAssetItem)
	}
	var poolRule []interface{}
	for _, poolItem := range pool {
		poolRule = append(poolRule, poolItem)
	}

	logs, sub, err := _StableDebtToken.contract.WatchLogs(opts, "Initialized", underlyingAssetRule, poolRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StableDebtTokenInitialized)
				if err := _StableDebtToken.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x40251fbfb6656cfa65a00d7879029fec1fad21d28fdcff2f4f68f52795b74f2c.
//
// Solidity: event Initialized(address indexed underlyingAsset, address indexed pool, address incentivesController, uint8 debtTokenDecimals, string debtTokenName, string debtTokenSymbol, bytes params)
func (_StableDebtToken *StableDebtTokenFilterer) ParseInitialized(log types.Log) (*StableDebtTokenInitialized, error) {
	event := new(StableDebtTokenInitialized)
	if err := _StableDebtToken.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StableDebtTokenMintIterator is returned from FilterMint and is used to iterate over the raw logs and unpacked data for Mint events raised by the StableDebtToken contract.
type StableDebtTokenMintIterator struct {
	Event *StableDebtTokenMint // Event containing the contract specifics and raw log

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
func (it *StableDebtTokenMintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StableDebtTokenMint)
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
		it.Event = new(StableDebtTokenMint)
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
func (it *StableDebtTokenMintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StableDebtTokenMintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StableDebtTokenMint represents a Mint event raised by the StableDebtToken contract.
type StableDebtTokenMint struct {
	User            common.Address
	OnBehalfOf      common.Address
	Amount          *big.Int
	CurrentBalance  *big.Int
	BalanceIncrease *big.Int
	NewRate         *big.Int
	AvgStableRate   *big.Int
	NewTotalSupply  *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterMint is a free log retrieval operation binding the contract event 0xc16f4e4ca34d790de4c656c72fd015c667d688f20be64eea360618545c4c530f.
//
// Solidity: event Mint(address indexed user, address indexed onBehalfOf, uint256 amount, uint256 currentBalance, uint256 balanceIncrease, uint256 newRate, uint256 avgStableRate, uint256 newTotalSupply)
func (_StableDebtToken *StableDebtTokenFilterer) FilterMint(opts *bind.FilterOpts, user []common.Address, onBehalfOf []common.Address) (*StableDebtTokenMintIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var onBehalfOfRule []interface{}
	for _, onBehalfOfItem := range onBehalfOf {
		onBehalfOfRule = append(onBehalfOfRule, onBehalfOfItem)
	}

	logs, sub, err := _StableDebtToken.contract.FilterLogs(opts, "Mint", userRule, onBehalfOfRule)
	if err != nil {
		return nil, err
	}
	return &StableDebtTokenMintIterator{contract: _StableDebtToken.contract, event: "Mint", logs: logs, sub: sub}, nil
}

// WatchMint is a free log subscription operation binding the contract event 0xc16f4e4ca34d790de4c656c72fd015c667d688f20be64eea360618545c4c530f.
//
// Solidity: event Mint(address indexed user, address indexed onBehalfOf, uint256 amount, uint256 currentBalance, uint256 balanceIncrease, uint256 newRate, uint256 avgStableRate, uint256 newTotalSupply)
func (_StableDebtToken *StableDebtTokenFilterer) WatchMint(opts *bind.WatchOpts, sink chan<- *StableDebtTokenMint, user []common.Address, onBehalfOf []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var onBehalfOfRule []interface{}
	for _, onBehalfOfItem := range onBehalfOf {
		onBehalfOfRule = append(onBehalfOfRule, onBehalfOfItem)
	}

	logs, sub, err := _StableDebtToken.contract.WatchLogs(opts, "Mint", userRule, onBehalfOfRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StableDebtTokenMint)
				if err := _StableDebtToken.contract.UnpackLog(event, "Mint", log); err != nil {
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

// ParseMint is a log parse operation binding the contract event 0xc16f4e4ca34d790de4c656c72fd015c667d688f20be64eea360618545c4c530f.
//
// Solidity: event Mint(address indexed user, address indexed onBehalfOf, uint256 amount, uint256 currentBalance, uint256 balanceIncrease, uint256 newRate, uint256 avgStableRate, uint256 newTotalSupply)
func (_StableDebtToken *StableDebtTokenFilterer) ParseMint(log types.Log) (*StableDebtTokenMint, error) {
	event := new(StableDebtTokenMint)
	if err := _StableDebtToken.contract.UnpackLog(event, "Mint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StableDebtTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the StableDebtToken contract.
type StableDebtTokenTransferIterator struct {
	Event *StableDebtTokenTransfer // Event containing the contract specifics and raw log

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
func (it *StableDebtTokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StableDebtTokenTransfer)
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
		it.Event = new(StableDebtTokenTransfer)
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
func (it *StableDebtTokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StableDebtTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StableDebtTokenTransfer represents a Transfer event raised by the StableDebtToken contract.
type StableDebtTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_StableDebtToken *StableDebtTokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*StableDebtTokenTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _StableDebtToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &StableDebtTokenTransferIterator{contract: _StableDebtToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_StableDebtToken *StableDebtTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *StableDebtTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _StableDebtToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StableDebtTokenTransfer)
				if err := _StableDebtToken.contract.UnpackLog(event, "Transfer", log); err != nil {
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
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_StableDebtToken *StableDebtTokenFilterer) ParseTransfer(log types.Log) (*StableDebtTokenTransfer, error) {
	event := new(StableDebtTokenTransfer)
	if err := _StableDebtToken.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
