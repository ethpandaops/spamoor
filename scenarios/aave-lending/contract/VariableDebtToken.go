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

// VariableDebtTokenMetaData contains all meta data concerning the VariableDebtToken contract.
var VariableDebtTokenMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIPool\",\"name\":\"pool\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"fromUser\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"toUser\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BorrowAllowanceDelegated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"balanceIncrease\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"Burn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"underlyingAsset\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"incentivesController\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"debtTokenDecimals\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"debtTokenName\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"debtTokenSymbol\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"params\",\"type\":\"bytes\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"onBehalfOf\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"balanceIncrease\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"Mint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEBT_TOKEN_REVISION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DELEGATION_WITH_SIG_TYPEHASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DOMAIN_SEPARATOR\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"EIP712_REVISION\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"POOL\",\"outputs\":[{\"internalType\":\"contractIPool\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"UNDERLYING_ASSET_ADDRESS\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegatee\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approveDelegation\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"fromUser\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"toUser\",\"type\":\"address\"}],\"name\":\"borrowAllowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegatee\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"delegationWithSig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getIncentivesController\",\"outputs\":[{\"internalType\":\"contractIAaveIncentivesController\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getPreviousIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getScaledUserBalanceAndSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIPool\",\"name\":\"initializingPool\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"underlyingAsset\",\"type\":\"address\"},{\"internalType\":\"contractIAaveIncentivesController\",\"name\":\"incentivesController\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"debtTokenDecimals\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"debtTokenName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"debtTokenSymbol\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"params\",\"type\":\"bytes\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"onBehalfOf\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"scaledBalanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"scaledTotalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIAaveIncentivesController\",\"name\":\"controller\",\"type\":\"address\"}],\"name\":\"setIncentivesController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60e0604052600080553480156200001557600080fd5b50604051620027be380380620027be833981016040819052620000389162000245565b806040518060400160405280601881526020017f5641524941424c455f444542545f544f4b454e5f494d504c00000000000000008152506040518060400160405280601881526020017f5641524941424c455f444542545f544f4b454e5f494d504c0000000000000000815250600083838383838383834660808181525050836001600160a01b0316630542975c6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620000f6573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200011c919062000245565b6001600160a01b031660a05282516200013d90603b90602086019062000186565b5081516200015390603c90602085019062000186565b50603d805460ff191660ff9290921691909117905550506001600160a01b031660c05250620002a9975050505050505050565b82805462000194906200026c565b90600052602060002090601f016020900481019282620001b8576000855562000203565b82601f10620001d357805160ff191683800117855562000203565b8280016001018555821562000203579182015b8281111562000203578251825591602001919060010190620001e6565b506200021192915062000215565b5090565b5b8082111562000211576000815560010162000216565b6001600160a01b03811681146200024257600080fd5b50565b6000602082840312156200025857600080fd5b815162000265816200022c565b9392505050565b600181811c908216806200028157607f821691505b60208210811415620002a357634e487b7160e01b600052602260045260246000fd5b50919050565b60805160a05160c0516124bb620003036000396000818161037e01528181610a3901528181610b7f01528181610c4e01528181610e1201528181610f6d015261124d0152600061103901526000610ab801526124bb6000f3fe608060405234801561001057600080fd5b50600436106101da5760003560e01c80637ecebe0011610104578063b9a7b622116100a2578063e075398611610071578063e0753986146104ee578063e655dbd81461054a578063f3bfc7381461055d578063f5298aca1461058457600080fd5b8063b9a7b622146104b2578063c04a8a10146104ba578063c222ec8a146104cd578063dd62ed3e146104e057600080fd5b8063a9059cbb116100de578063a9059cbb146101fd578063b16a19de14610462578063b1bf962d14610480578063b3f1c93d1461048857600080fd5b80637ecebe001461042457806395d89b411461045a578063a457c2d7146101fd57600080fd5b8063313ce5671161017c57806370a082311161014b57806370a08231146103665780637535d2461461037957806375d26413146103c557806378160376146103e857600080fd5b8063313ce567146103035780633644e5151461031857806339509351146101fd5780636bd76d241461032057600080fd5b80630b52d558116101b85780630b52d5581461028257806318160ddd146102975780631da24f3e146102ad57806323b872dd146102f557600080fd5b806306fdde03146101df578063095ea7b3146101fd5780630afbcdc914610220575b600080fd5b6101e7610597565b6040516101f49190611e79565b60405180910390f35b61021061020b366004611ec1565b610629565b60405190151581526020016101f4565b61026d61022e366004611eed565b73ffffffffffffffffffffffffffffffffffffffff16600090815260386020526040902054603a546fffffffffffffffffffffffffffffffff90911691565b604080519283526020830191909152016101f4565b610295610290366004611f1b565b610699565b005b61029f6109ea565b6040519081526020016101f4565b61029f6102bb366004611eed565b73ffffffffffffffffffffffffffffffffffffffff166000908152603860205260409020546fffffffffffffffffffffffffffffffff1690565b61021061020b366004611f89565b603d5460405160ff90911681526020016101f4565b61029f610ab4565b61029f61032e366004611fca565b73ffffffffffffffffffffffffffffffffffffffff918216600090815260366020908152604080832093909416825291909152205490565b61029f610374366004611eed565b610aed565b6103a07f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101f4565b603d54610100900473ffffffffffffffffffffffffffffffffffffffff166103a0565b6101e76040518060400160405280600181526020017f310000000000000000000000000000000000000000000000000000000000000081525081565b61029f610432366004611eed565b73ffffffffffffffffffffffffffffffffffffffff1660009081526034602052604090205490565b6101e7610bf8565b60375473ffffffffffffffffffffffffffffffffffffffff166103a0565b61029f610c07565b61049b610496366004612003565b610c12565b6040805192151583526020830191909152016101f4565b61029f600181565b6102956104c8366004611ec1565b610d1b565b6102956104db36600461216c565b610d2a565b61029f61020b366004611fca565b61029f6104fc366004611eed565b73ffffffffffffffffffffffffffffffffffffffff1660009081526038602052604090205470010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff1690565b610295610558366004611eed565b611035565b61029f7f323db0410fecc107e39e2af5908671f4c8d106123b35a51501bb805c5fa36aa081565b61029f610592366004612241565b611213565b6060603b80546105a690612276565b80601f01602080910402602001604051908101604052809291908181526020018280546105d290612276565b801561061f5780601f106105f45761010080835404028352916020019161061f565b820191906000526020600020905b81548152906001019060200180831161060257829003601f168201915b5050505050905090565b604080518082018252600281527f3830000000000000000000000000000000000000000000000000000000000000602082015290517f08c379a000000000000000000000000000000000000000000000000000000000815260009161069091600401611e79565b60405180910390fd5b60408051808201909152600281527f3737000000000000000000000000000000000000000000000000000000000000602082015273ffffffffffffffffffffffffffffffffffffffff881661071b576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106909190611e79565b50834211156040518060400160405280600281526020017f37380000000000000000000000000000000000000000000000000000000000008152509061078e576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106909190611e79565b5073ffffffffffffffffffffffffffffffffffffffff8716600090815260346020526040812054906107be610ab4565b604080517f323db0410fecc107e39e2af5908671f4c8d106123b35a51501bb805c5fa36aa0602082015273ffffffffffffffffffffffffffffffffffffffff8b1691810191909152606081018990526080810184905260a0810188905260c001604051602081830303815290604052805190602001206040516020016108769291907f190100000000000000000000000000000000000000000000000000000000000081526002810192909252602282015260420190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201206000845290830180835281905260ff8816918301919091526060820186905260808201859052915060019060a0016020604051602081039080840390855afa1580156108fc573d6000803e3d6000fd5b5050506020604051035173ffffffffffffffffffffffffffffffffffffffff168973ffffffffffffffffffffffffffffffffffffffff16146040518060400160405280600281526020017f3739000000000000000000000000000000000000000000000000000000000000815250906109a2576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106909190611e79565b506109ae8260016122f9565b73ffffffffffffffffffffffffffffffffffffffff8a166000908152603460205260409020556109df8989896112d8565b505050505050505050565b6037546040517f386497fd00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9182166004820152600091610aaf917f00000000000000000000000000000000000000000000000000000000000000009091169063386497fd90602401602060405180830381865afa158015610a82573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610aa69190612311565b603a549061134f565b905090565b60007f0000000000000000000000000000000000000000000000000000000000000000461415610ae5575060355490565b610aaf6113a6565b73ffffffffffffffffffffffffffffffffffffffff81166000908152603860205260408120546fffffffffffffffffffffffffffffffff1680610b335750600092915050565b6037546040517f386497fd00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9182166004820152610bf1917f0000000000000000000000000000000000000000000000000000000000000000169063386497fd90602401602060405180830381865afa158015610bc6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610bea9190612311565b829061134f565b9392505050565b6060603c80546105a690612276565b6000610aaf603a5490565b60408051808201909152600281527f323300000000000000000000000000000000000000000000000000000000000060208201526000908190337f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1614610cbb576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106909190611e79565b508473ffffffffffffffffffffffffffffffffffffffff168673ffffffffffffffffffffffffffffffffffffffff1614610cfa57610cfa85878661146b565b610d068686868661152b565b610d0e610c07565b9150915094509492505050565b610d263383836112d8565b5050565b6001805460ff1680610d3b5750303b155b80610d47575060005481115b610dd3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f436f6e747261637420696e7374616e63652068617320616c726561647920626560448201527f656e20696e697469616c697a65640000000000000000000000000000000000006064820152608401610690565b60015460ff16158015610e1057600180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00168117905560008290555b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168a73ffffffffffffffffffffffffffffffffffffffff16146040518060400160405280600281526020017f383700000000000000000000000000000000000000000000000000000000000081525090610ecd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106909190611e79565b50610ed78661176c565b610ee08561177f565b603d80546037805473ffffffffffffffffffffffffffffffffffffffff8d81167fffffffffffffffffffffffff0000000000000000000000000000000000000000909216919091179091558a16610100027fffffffffffffffffffffff00000000000000000000000000000000000000000090911660ff8a1617179055610f656113a6565b6035819055507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168973ffffffffffffffffffffffffffffffffffffffff167f40251fbfb6656cfa65a00d7879029fec1fad21d28fdcff2f4f68f52795b74f2c8a8a8a8a8a8a604051610ff29695949392919061232a565b60405180910390a3801561102957600180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690555b50505050505050505050565b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663707cd7166040518163ffffffff1660e01b8152600401602060405180830381865afa1580156110a2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906110c691906123ca565b6040517f7be53ca100000000000000000000000000000000000000000000000000000000815233600482015290915073ffffffffffffffffffffffffffffffffffffffff821690637be53ca190602401602060405180830381865afa158015611133573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061115791906123e7565b6040518060400160405280600181526020017f3100000000000000000000000000000000000000000000000000000000000000815250906111c5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106909190611e79565b5050603d805473ffffffffffffffffffffffffffffffffffffffff909216610100027fffffffffffffffffffffff0000000000000000000000000000000000000000ff909216919091179055565b60408051808201909152600281527f32330000000000000000000000000000000000000000000000000000000000006020820152600090337f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16146112ba576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106909190611e79565b506112c88460008585611792565b6112d0610c07565b949350505050565b73ffffffffffffffffffffffffffffffffffffffff83811660008181526036602090815260408083208786168085529083529281902086905560375490518681529416939192917fda919360433220e13b51e8c211e490d148e61a3bd53de8c097194e458b97f3e1910160405180910390a4505050565b600081157ffffffffffffffffffffffffffffffffffffffffffe6268e1b017bfe18bffffff8390048411151761138457600080fd5b506b033b2e3c9fd0803ce800000091026b019d971e4fe8401e74000000010490565b60007f8b73c3c69bb8fe3d512ecc4cf759cc79239f7b179b0ffacaa9a75d522b39400f6113d1611aaf565b8051602091820120604080518082018252600181527f310000000000000000000000000000000000000000000000000000000000000090840152805192830193909352918101919091527fc89efdaa54c0f20c7adf612882df0950f5a951637e0307cdcb4c672f298b8bc660608201524660808201523060a082015260c00160405160208183030381529060405280519060200120905090565b73ffffffffffffffffffffffffffffffffffffffff80841660009081526036602090815260408083209386168352929052908120546114ab908390612409565b73ffffffffffffffffffffffffffffffffffffffff808616600081815260366020908152604080832089861680855292529182902085905560375491519495509216927fda919360433220e13b51e8c211e490d148e61a3bd53de8c097194e458b97f3e19061151d9086815260200190565b60405180910390a450505050565b6000806115388484611ab9565b60408051808201909152600281527f32340000000000000000000000000000000000000000000000000000000000006020820152909150816115a7576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106909190611e79565b5073ffffffffffffffffffffffffffffffffffffffff85166000908152603860205260408120546fffffffffffffffffffffffffffffffff808216929161160491849170010000000000000000000000000000000090041661134f565b61160e838761134f565b6116189190612409565b905061162385611af8565b73ffffffffffffffffffffffffffffffffffffffff8816600090815260386020526040902080546fffffffffffffffffffffffffffffffff92831670010000000000000000000000000000000002921691909117905561168b8761168685611af8565b611b9e565b600061169782886122f9565b90508773ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040516116f991815260200190565b60405180910390a3604080518281526020810184905290810187905273ffffffffffffffffffffffffffffffffffffffff808a1691908b16907f458f5fa412d0f69b08dd84872b0215675cc67bc1d5b6fd93300a1c3878b861969060600160405180910390a35050159695505050505050565b8051610d2690603b906020840190611d7e565b8051610d2690603c906020840190611d7e565b600061179e8383611ab9565b60408051808201909152600281527f323500000000000000000000000000000000000000000000000000000000000060208201529091508161180d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106909190611e79565b5073ffffffffffffffffffffffffffffffffffffffff85166000908152603860205260408120546fffffffffffffffffffffffffffffffff808216929161186a91849170010000000000000000000000000000000090041661134f565b611874838661134f565b61187e9190612409565b905061188984611af8565b73ffffffffffffffffffffffffffffffffffffffff8816600090815260386020526040902080546fffffffffffffffffffffffffffffffff9283167001000000000000000000000000000000000292169190911790556118f1876118ec85611af8565b611d1a565b848111156119d05760006119058683612409565b90508773ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8360405161196791815260200190565b60405180910390a3604080518281526020810184905290810186905273ffffffffffffffffffffffffffffffffffffffff89169081907f458f5fa412d0f69b08dd84872b0215675cc67bc1d5b6fd93300a1c3878b861969060600160405180910390a350611aa6565b60006119dc8287612409565b9050600073ffffffffffffffffffffffffffffffffffffffff168873ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef83604051611a3e91815260200190565b60405180910390a3604080518281526020810184905290810186905273ffffffffffffffffffffffffffffffffffffffff80891691908a16907f4cf25bc1d991c17529c25213d3cc0cda295eeaad5f13f361969b12ea48015f909060600160405180910390a3505b50505050505050565b6060610aaf610597565b600081156b033b2e3c9fd0803ce800000060028404190484111715611add57600080fd5b506b033b2e3c9fd0803ce80000009190910260028204010490565b60006fffffffffffffffffffffffffffffffff821115611b9a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602760248201527f53616665436173743a2076616c756520646f65736e27742066697420696e203160448201527f32382062697473000000000000000000000000000000000000000000000000006064820152608401610690565b5090565b603a54611bbd6fffffffffffffffffffffffffffffffff8316826122f9565b603a5573ffffffffffffffffffffffffffffffffffffffff83166000908152603860205260409020546fffffffffffffffffffffffffffffffff16611c028382612420565b73ffffffffffffffffffffffffffffffffffffffff858116600090815260386020526040902080547fffffffffffffffffffffffffffffffff00000000000000000000000000000000166fffffffffffffffffffffffffffffffff9390931692909217909155603d546101009004168015611d13576040517f31873e2e00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8681166004830152602482018590526fffffffffffffffffffffffffffffffff841660448301528216906331873e2e90606401600060405180830381600087803b158015611cff57600080fd5b505af11580156109df573d6000803e3d6000fd5b5050505050565b603a54611d396fffffffffffffffffffffffffffffffff831682612409565b603a5573ffffffffffffffffffffffffffffffffffffffff83166000908152603860205260409020546fffffffffffffffffffffffffffffffff16611c028382612454565b828054611d8a90612276565b90600052602060002090601f016020900481019282611dac5760008555611df2565b82601f10611dc557805160ff1916838001178555611df2565b82800160010185558215611df2579182015b82811115611df2578251825591602001919060010190611dd7565b50611b9a9291505b80821115611b9a5760008155600101611dfa565b6000815180845260005b81811015611e3457602081850181015186830182015201611e18565b81811115611e46576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000610bf16020830184611e0e565b73ffffffffffffffffffffffffffffffffffffffff81168114611eae57600080fd5b50565b8035611ebc81611e8c565b919050565b60008060408385031215611ed457600080fd5b8235611edf81611e8c565b946020939093013593505050565b600060208284031215611eff57600080fd5b8135610bf181611e8c565b803560ff81168114611ebc57600080fd5b600080600080600080600060e0888a031215611f3657600080fd5b8735611f4181611e8c565b96506020880135611f5181611e8c565b95506040880135945060608801359350611f6d60808901611f0a565b925060a0880135915060c0880135905092959891949750929550565b600080600060608486031215611f9e57600080fd5b8335611fa981611e8c565b92506020840135611fb981611e8c565b929592945050506040919091013590565b60008060408385031215611fdd57600080fd5b8235611fe881611e8c565b91506020830135611ff881611e8c565b809150509250929050565b6000806000806080858703121561201957600080fd5b843561202481611e8c565b9350602085013561203481611e8c565b93969395505050506040820135916060013590565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600082601f83011261208957600080fd5b813567ffffffffffffffff808211156120a4576120a4612049565b604051601f83017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f011681019082821181831017156120ea576120ea612049565b8160405283815286602085880101111561210357600080fd5b836020870160208301376000602085830101528094505050505092915050565b60008083601f84011261213557600080fd5b50813567ffffffffffffffff81111561214d57600080fd5b60208301915083602082850101111561216557600080fd5b9250929050565b60008060008060008060008060e0898b03121561218857600080fd5b883561219381611e8c565b975060208901356121a381611e8c565b96506121b160408a01611eb1565b95506121bf60608a01611f0a565b9450608089013567ffffffffffffffff808211156121dc57600080fd5b6121e88c838d01612078565b955060a08b01359150808211156121fe57600080fd5b61220a8c838d01612078565b945060c08b013591508082111561222057600080fd5b5061222d8b828c01612123565b999c989b5096995094979396929594505050565b60008060006060848603121561225657600080fd5b833561226181611e8c565b95602085013595506040909401359392505050565b600181811c9082168061228a57607f821691505b602082108114156122c4577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000821982111561230c5761230c6122ca565b500190565b60006020828403121561232357600080fd5b5051919050565b73ffffffffffffffffffffffffffffffffffffffff8716815260ff8616602082015260a06040820152600061236260a0830187611e0e565b82810360608401526123748187611e0e565b905082810360808401528381528385602083013760006020858301015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f860116820101915050979650505050505050565b6000602082840312156123dc57600080fd5b8151610bf181611e8c565b6000602082840312156123f957600080fd5b81518015158114610bf157600080fd5b60008282101561241b5761241b6122ca565b500390565b60006fffffffffffffffffffffffffffffffff80831681851680830382111561244b5761244b6122ca565b01949350505050565b60006fffffffffffffffffffffffffffffffff8381169083168181101561247d5761247d6122ca565b03939250505056fea26469706673582212208c26f709abec806abf23a180dc9b2d4a72a1949f588c15cbbd49912617d13a4a64736f6c634300080a0033",
}

// VariableDebtTokenABI is the input ABI used to generate the binding from.
// Deprecated: Use VariableDebtTokenMetaData.ABI instead.
var VariableDebtTokenABI = VariableDebtTokenMetaData.ABI

// VariableDebtTokenBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use VariableDebtTokenMetaData.Bin instead.
var VariableDebtTokenBin = VariableDebtTokenMetaData.Bin

// DeployVariableDebtToken deploys a new Ethereum contract, binding an instance of VariableDebtToken to it.
func DeployVariableDebtToken(auth *bind.TransactOpts, backend bind.ContractBackend, pool common.Address) (common.Address, *types.Transaction, *VariableDebtToken, error) {
	parsed, err := VariableDebtTokenMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VariableDebtTokenBin), backend, pool)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VariableDebtToken{VariableDebtTokenCaller: VariableDebtTokenCaller{contract: contract}, VariableDebtTokenTransactor: VariableDebtTokenTransactor{contract: contract}, VariableDebtTokenFilterer: VariableDebtTokenFilterer{contract: contract}}, nil
}

// VariableDebtToken is an auto generated Go binding around an Ethereum contract.
type VariableDebtToken struct {
	VariableDebtTokenCaller     // Read-only binding to the contract
	VariableDebtTokenTransactor // Write-only binding to the contract
	VariableDebtTokenFilterer   // Log filterer for contract events
}

// VariableDebtTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type VariableDebtTokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VariableDebtTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VariableDebtTokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VariableDebtTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VariableDebtTokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VariableDebtTokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VariableDebtTokenSession struct {
	Contract     *VariableDebtToken // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// VariableDebtTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VariableDebtTokenCallerSession struct {
	Contract *VariableDebtTokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// VariableDebtTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VariableDebtTokenTransactorSession struct {
	Contract     *VariableDebtTokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// VariableDebtTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type VariableDebtTokenRaw struct {
	Contract *VariableDebtToken // Generic contract binding to access the raw methods on
}

// VariableDebtTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VariableDebtTokenCallerRaw struct {
	Contract *VariableDebtTokenCaller // Generic read-only contract binding to access the raw methods on
}

// VariableDebtTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VariableDebtTokenTransactorRaw struct {
	Contract *VariableDebtTokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVariableDebtToken creates a new instance of VariableDebtToken, bound to a specific deployed contract.
func NewVariableDebtToken(address common.Address, backend bind.ContractBackend) (*VariableDebtToken, error) {
	contract, err := bindVariableDebtToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VariableDebtToken{VariableDebtTokenCaller: VariableDebtTokenCaller{contract: contract}, VariableDebtTokenTransactor: VariableDebtTokenTransactor{contract: contract}, VariableDebtTokenFilterer: VariableDebtTokenFilterer{contract: contract}}, nil
}

// NewVariableDebtTokenCaller creates a new read-only instance of VariableDebtToken, bound to a specific deployed contract.
func NewVariableDebtTokenCaller(address common.Address, caller bind.ContractCaller) (*VariableDebtTokenCaller, error) {
	contract, err := bindVariableDebtToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VariableDebtTokenCaller{contract: contract}, nil
}

// NewVariableDebtTokenTransactor creates a new write-only instance of VariableDebtToken, bound to a specific deployed contract.
func NewVariableDebtTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*VariableDebtTokenTransactor, error) {
	contract, err := bindVariableDebtToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VariableDebtTokenTransactor{contract: contract}, nil
}

// NewVariableDebtTokenFilterer creates a new log filterer instance of VariableDebtToken, bound to a specific deployed contract.
func NewVariableDebtTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*VariableDebtTokenFilterer, error) {
	contract, err := bindVariableDebtToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VariableDebtTokenFilterer{contract: contract}, nil
}

// bindVariableDebtToken binds a generic wrapper to an already deployed contract.
func bindVariableDebtToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VariableDebtTokenMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VariableDebtToken *VariableDebtTokenRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VariableDebtToken.Contract.VariableDebtTokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VariableDebtToken *VariableDebtTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VariableDebtToken.Contract.VariableDebtTokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VariableDebtToken *VariableDebtTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VariableDebtToken.Contract.VariableDebtTokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VariableDebtToken *VariableDebtTokenCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VariableDebtToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VariableDebtToken *VariableDebtTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VariableDebtToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VariableDebtToken *VariableDebtTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VariableDebtToken.Contract.contract.Transact(opts, method, params...)
}

// DEBTTOKENREVISION is a free data retrieval call binding the contract method 0xb9a7b622.
//
// Solidity: function DEBT_TOKEN_REVISION() view returns(uint256)
func (_VariableDebtToken *VariableDebtTokenCaller) DEBTTOKENREVISION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VariableDebtToken.contract.Call(opts, &out, "DEBT_TOKEN_REVISION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DEBTTOKENREVISION is a free data retrieval call binding the contract method 0xb9a7b622.
//
// Solidity: function DEBT_TOKEN_REVISION() view returns(uint256)
func (_VariableDebtToken *VariableDebtTokenSession) DEBTTOKENREVISION() (*big.Int, error) {
	return _VariableDebtToken.Contract.DEBTTOKENREVISION(&_VariableDebtToken.CallOpts)
}

// DEBTTOKENREVISION is a free data retrieval call binding the contract method 0xb9a7b622.
//
// Solidity: function DEBT_TOKEN_REVISION() view returns(uint256)
func (_VariableDebtToken *VariableDebtTokenCallerSession) DEBTTOKENREVISION() (*big.Int, error) {
	return _VariableDebtToken.Contract.DEBTTOKENREVISION(&_VariableDebtToken.CallOpts)
}

// DELEGATIONWITHSIGTYPEHASH is a free data retrieval call binding the contract method 0xf3bfc738.
//
// Solidity: function DELEGATION_WITH_SIG_TYPEHASH() view returns(bytes32)
func (_VariableDebtToken *VariableDebtTokenCaller) DELEGATIONWITHSIGTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _VariableDebtToken.contract.Call(opts, &out, "DELEGATION_WITH_SIG_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DELEGATIONWITHSIGTYPEHASH is a free data retrieval call binding the contract method 0xf3bfc738.
//
// Solidity: function DELEGATION_WITH_SIG_TYPEHASH() view returns(bytes32)
func (_VariableDebtToken *VariableDebtTokenSession) DELEGATIONWITHSIGTYPEHASH() ([32]byte, error) {
	return _VariableDebtToken.Contract.DELEGATIONWITHSIGTYPEHASH(&_VariableDebtToken.CallOpts)
}

// DELEGATIONWITHSIGTYPEHASH is a free data retrieval call binding the contract method 0xf3bfc738.
//
// Solidity: function DELEGATION_WITH_SIG_TYPEHASH() view returns(bytes32)
func (_VariableDebtToken *VariableDebtTokenCallerSession) DELEGATIONWITHSIGTYPEHASH() ([32]byte, error) {
	return _VariableDebtToken.Contract.DELEGATIONWITHSIGTYPEHASH(&_VariableDebtToken.CallOpts)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_VariableDebtToken *VariableDebtTokenCaller) DOMAINSEPARATOR(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _VariableDebtToken.contract.Call(opts, &out, "DOMAIN_SEPARATOR")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_VariableDebtToken *VariableDebtTokenSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _VariableDebtToken.Contract.DOMAINSEPARATOR(&_VariableDebtToken.CallOpts)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_VariableDebtToken *VariableDebtTokenCallerSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _VariableDebtToken.Contract.DOMAINSEPARATOR(&_VariableDebtToken.CallOpts)
}

// EIP712REVISION is a free data retrieval call binding the contract method 0x78160376.
//
// Solidity: function EIP712_REVISION() view returns(bytes)
func (_VariableDebtToken *VariableDebtTokenCaller) EIP712REVISION(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _VariableDebtToken.contract.Call(opts, &out, "EIP712_REVISION")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// EIP712REVISION is a free data retrieval call binding the contract method 0x78160376.
//
// Solidity: function EIP712_REVISION() view returns(bytes)
func (_VariableDebtToken *VariableDebtTokenSession) EIP712REVISION() ([]byte, error) {
	return _VariableDebtToken.Contract.EIP712REVISION(&_VariableDebtToken.CallOpts)
}

// EIP712REVISION is a free data retrieval call binding the contract method 0x78160376.
//
// Solidity: function EIP712_REVISION() view returns(bytes)
func (_VariableDebtToken *VariableDebtTokenCallerSession) EIP712REVISION() ([]byte, error) {
	return _VariableDebtToken.Contract.EIP712REVISION(&_VariableDebtToken.CallOpts)
}

// POOL is a free data retrieval call binding the contract method 0x7535d246.
//
// Solidity: function POOL() view returns(address)
func (_VariableDebtToken *VariableDebtTokenCaller) POOL(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VariableDebtToken.contract.Call(opts, &out, "POOL")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// POOL is a free data retrieval call binding the contract method 0x7535d246.
//
// Solidity: function POOL() view returns(address)
func (_VariableDebtToken *VariableDebtTokenSession) POOL() (common.Address, error) {
	return _VariableDebtToken.Contract.POOL(&_VariableDebtToken.CallOpts)
}

// POOL is a free data retrieval call binding the contract method 0x7535d246.
//
// Solidity: function POOL() view returns(address)
func (_VariableDebtToken *VariableDebtTokenCallerSession) POOL() (common.Address, error) {
	return _VariableDebtToken.Contract.POOL(&_VariableDebtToken.CallOpts)
}

// UNDERLYINGASSETADDRESS is a free data retrieval call binding the contract method 0xb16a19de.
//
// Solidity: function UNDERLYING_ASSET_ADDRESS() view returns(address)
func (_VariableDebtToken *VariableDebtTokenCaller) UNDERLYINGASSETADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VariableDebtToken.contract.Call(opts, &out, "UNDERLYING_ASSET_ADDRESS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UNDERLYINGASSETADDRESS is a free data retrieval call binding the contract method 0xb16a19de.
//
// Solidity: function UNDERLYING_ASSET_ADDRESS() view returns(address)
func (_VariableDebtToken *VariableDebtTokenSession) UNDERLYINGASSETADDRESS() (common.Address, error) {
	return _VariableDebtToken.Contract.UNDERLYINGASSETADDRESS(&_VariableDebtToken.CallOpts)
}

// UNDERLYINGASSETADDRESS is a free data retrieval call binding the contract method 0xb16a19de.
//
// Solidity: function UNDERLYING_ASSET_ADDRESS() view returns(address)
func (_VariableDebtToken *VariableDebtTokenCallerSession) UNDERLYINGASSETADDRESS() (common.Address, error) {
	return _VariableDebtToken.Contract.UNDERLYINGASSETADDRESS(&_VariableDebtToken.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_VariableDebtToken *VariableDebtTokenCaller) Allowance(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _VariableDebtToken.contract.Call(opts, &out, "allowance", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_VariableDebtToken *VariableDebtTokenSession) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _VariableDebtToken.Contract.Allowance(&_VariableDebtToken.CallOpts, arg0, arg1)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_VariableDebtToken *VariableDebtTokenCallerSession) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _VariableDebtToken.Contract.Allowance(&_VariableDebtToken.CallOpts, arg0, arg1)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address user) view returns(uint256)
func (_VariableDebtToken *VariableDebtTokenCaller) BalanceOf(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _VariableDebtToken.contract.Call(opts, &out, "balanceOf", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address user) view returns(uint256)
func (_VariableDebtToken *VariableDebtTokenSession) BalanceOf(user common.Address) (*big.Int, error) {
	return _VariableDebtToken.Contract.BalanceOf(&_VariableDebtToken.CallOpts, user)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address user) view returns(uint256)
func (_VariableDebtToken *VariableDebtTokenCallerSession) BalanceOf(user common.Address) (*big.Int, error) {
	return _VariableDebtToken.Contract.BalanceOf(&_VariableDebtToken.CallOpts, user)
}

// BorrowAllowance is a free data retrieval call binding the contract method 0x6bd76d24.
//
// Solidity: function borrowAllowance(address fromUser, address toUser) view returns(uint256)
func (_VariableDebtToken *VariableDebtTokenCaller) BorrowAllowance(opts *bind.CallOpts, fromUser common.Address, toUser common.Address) (*big.Int, error) {
	var out []interface{}
	err := _VariableDebtToken.contract.Call(opts, &out, "borrowAllowance", fromUser, toUser)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BorrowAllowance is a free data retrieval call binding the contract method 0x6bd76d24.
//
// Solidity: function borrowAllowance(address fromUser, address toUser) view returns(uint256)
func (_VariableDebtToken *VariableDebtTokenSession) BorrowAllowance(fromUser common.Address, toUser common.Address) (*big.Int, error) {
	return _VariableDebtToken.Contract.BorrowAllowance(&_VariableDebtToken.CallOpts, fromUser, toUser)
}

// BorrowAllowance is a free data retrieval call binding the contract method 0x6bd76d24.
//
// Solidity: function borrowAllowance(address fromUser, address toUser) view returns(uint256)
func (_VariableDebtToken *VariableDebtTokenCallerSession) BorrowAllowance(fromUser common.Address, toUser common.Address) (*big.Int, error) {
	return _VariableDebtToken.Contract.BorrowAllowance(&_VariableDebtToken.CallOpts, fromUser, toUser)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_VariableDebtToken *VariableDebtTokenCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _VariableDebtToken.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_VariableDebtToken *VariableDebtTokenSession) Decimals() (uint8, error) {
	return _VariableDebtToken.Contract.Decimals(&_VariableDebtToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_VariableDebtToken *VariableDebtTokenCallerSession) Decimals() (uint8, error) {
	return _VariableDebtToken.Contract.Decimals(&_VariableDebtToken.CallOpts)
}

// GetIncentivesController is a free data retrieval call binding the contract method 0x75d26413.
//
// Solidity: function getIncentivesController() view returns(address)
func (_VariableDebtToken *VariableDebtTokenCaller) GetIncentivesController(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VariableDebtToken.contract.Call(opts, &out, "getIncentivesController")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetIncentivesController is a free data retrieval call binding the contract method 0x75d26413.
//
// Solidity: function getIncentivesController() view returns(address)
func (_VariableDebtToken *VariableDebtTokenSession) GetIncentivesController() (common.Address, error) {
	return _VariableDebtToken.Contract.GetIncentivesController(&_VariableDebtToken.CallOpts)
}

// GetIncentivesController is a free data retrieval call binding the contract method 0x75d26413.
//
// Solidity: function getIncentivesController() view returns(address)
func (_VariableDebtToken *VariableDebtTokenCallerSession) GetIncentivesController() (common.Address, error) {
	return _VariableDebtToken.Contract.GetIncentivesController(&_VariableDebtToken.CallOpts)
}

// GetPreviousIndex is a free data retrieval call binding the contract method 0xe0753986.
//
// Solidity: function getPreviousIndex(address user) view returns(uint256)
func (_VariableDebtToken *VariableDebtTokenCaller) GetPreviousIndex(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _VariableDebtToken.contract.Call(opts, &out, "getPreviousIndex", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPreviousIndex is a free data retrieval call binding the contract method 0xe0753986.
//
// Solidity: function getPreviousIndex(address user) view returns(uint256)
func (_VariableDebtToken *VariableDebtTokenSession) GetPreviousIndex(user common.Address) (*big.Int, error) {
	return _VariableDebtToken.Contract.GetPreviousIndex(&_VariableDebtToken.CallOpts, user)
}

// GetPreviousIndex is a free data retrieval call binding the contract method 0xe0753986.
//
// Solidity: function getPreviousIndex(address user) view returns(uint256)
func (_VariableDebtToken *VariableDebtTokenCallerSession) GetPreviousIndex(user common.Address) (*big.Int, error) {
	return _VariableDebtToken.Contract.GetPreviousIndex(&_VariableDebtToken.CallOpts, user)
}

// GetScaledUserBalanceAndSupply is a free data retrieval call binding the contract method 0x0afbcdc9.
//
// Solidity: function getScaledUserBalanceAndSupply(address user) view returns(uint256, uint256)
func (_VariableDebtToken *VariableDebtTokenCaller) GetScaledUserBalanceAndSupply(opts *bind.CallOpts, user common.Address) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _VariableDebtToken.contract.Call(opts, &out, "getScaledUserBalanceAndSupply", user)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// GetScaledUserBalanceAndSupply is a free data retrieval call binding the contract method 0x0afbcdc9.
//
// Solidity: function getScaledUserBalanceAndSupply(address user) view returns(uint256, uint256)
func (_VariableDebtToken *VariableDebtTokenSession) GetScaledUserBalanceAndSupply(user common.Address) (*big.Int, *big.Int, error) {
	return _VariableDebtToken.Contract.GetScaledUserBalanceAndSupply(&_VariableDebtToken.CallOpts, user)
}

// GetScaledUserBalanceAndSupply is a free data retrieval call binding the contract method 0x0afbcdc9.
//
// Solidity: function getScaledUserBalanceAndSupply(address user) view returns(uint256, uint256)
func (_VariableDebtToken *VariableDebtTokenCallerSession) GetScaledUserBalanceAndSupply(user common.Address) (*big.Int, *big.Int, error) {
	return _VariableDebtToken.Contract.GetScaledUserBalanceAndSupply(&_VariableDebtToken.CallOpts, user)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_VariableDebtToken *VariableDebtTokenCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _VariableDebtToken.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_VariableDebtToken *VariableDebtTokenSession) Name() (string, error) {
	return _VariableDebtToken.Contract.Name(&_VariableDebtToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_VariableDebtToken *VariableDebtTokenCallerSession) Name() (string, error) {
	return _VariableDebtToken.Contract.Name(&_VariableDebtToken.CallOpts)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_VariableDebtToken *VariableDebtTokenCaller) Nonces(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _VariableDebtToken.contract.Call(opts, &out, "nonces", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_VariableDebtToken *VariableDebtTokenSession) Nonces(owner common.Address) (*big.Int, error) {
	return _VariableDebtToken.Contract.Nonces(&_VariableDebtToken.CallOpts, owner)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_VariableDebtToken *VariableDebtTokenCallerSession) Nonces(owner common.Address) (*big.Int, error) {
	return _VariableDebtToken.Contract.Nonces(&_VariableDebtToken.CallOpts, owner)
}

// ScaledBalanceOf is a free data retrieval call binding the contract method 0x1da24f3e.
//
// Solidity: function scaledBalanceOf(address user) view returns(uint256)
func (_VariableDebtToken *VariableDebtTokenCaller) ScaledBalanceOf(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _VariableDebtToken.contract.Call(opts, &out, "scaledBalanceOf", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ScaledBalanceOf is a free data retrieval call binding the contract method 0x1da24f3e.
//
// Solidity: function scaledBalanceOf(address user) view returns(uint256)
func (_VariableDebtToken *VariableDebtTokenSession) ScaledBalanceOf(user common.Address) (*big.Int, error) {
	return _VariableDebtToken.Contract.ScaledBalanceOf(&_VariableDebtToken.CallOpts, user)
}

// ScaledBalanceOf is a free data retrieval call binding the contract method 0x1da24f3e.
//
// Solidity: function scaledBalanceOf(address user) view returns(uint256)
func (_VariableDebtToken *VariableDebtTokenCallerSession) ScaledBalanceOf(user common.Address) (*big.Int, error) {
	return _VariableDebtToken.Contract.ScaledBalanceOf(&_VariableDebtToken.CallOpts, user)
}

// ScaledTotalSupply is a free data retrieval call binding the contract method 0xb1bf962d.
//
// Solidity: function scaledTotalSupply() view returns(uint256)
func (_VariableDebtToken *VariableDebtTokenCaller) ScaledTotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VariableDebtToken.contract.Call(opts, &out, "scaledTotalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ScaledTotalSupply is a free data retrieval call binding the contract method 0xb1bf962d.
//
// Solidity: function scaledTotalSupply() view returns(uint256)
func (_VariableDebtToken *VariableDebtTokenSession) ScaledTotalSupply() (*big.Int, error) {
	return _VariableDebtToken.Contract.ScaledTotalSupply(&_VariableDebtToken.CallOpts)
}

// ScaledTotalSupply is a free data retrieval call binding the contract method 0xb1bf962d.
//
// Solidity: function scaledTotalSupply() view returns(uint256)
func (_VariableDebtToken *VariableDebtTokenCallerSession) ScaledTotalSupply() (*big.Int, error) {
	return _VariableDebtToken.Contract.ScaledTotalSupply(&_VariableDebtToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_VariableDebtToken *VariableDebtTokenCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _VariableDebtToken.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_VariableDebtToken *VariableDebtTokenSession) Symbol() (string, error) {
	return _VariableDebtToken.Contract.Symbol(&_VariableDebtToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_VariableDebtToken *VariableDebtTokenCallerSession) Symbol() (string, error) {
	return _VariableDebtToken.Contract.Symbol(&_VariableDebtToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_VariableDebtToken *VariableDebtTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VariableDebtToken.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_VariableDebtToken *VariableDebtTokenSession) TotalSupply() (*big.Int, error) {
	return _VariableDebtToken.Contract.TotalSupply(&_VariableDebtToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_VariableDebtToken *VariableDebtTokenCallerSession) TotalSupply() (*big.Int, error) {
	return _VariableDebtToken.Contract.TotalSupply(&_VariableDebtToken.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address , uint256 ) returns(bool)
func (_VariableDebtToken *VariableDebtTokenTransactor) Approve(opts *bind.TransactOpts, arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _VariableDebtToken.contract.Transact(opts, "approve", arg0, arg1)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address , uint256 ) returns(bool)
func (_VariableDebtToken *VariableDebtTokenSession) Approve(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _VariableDebtToken.Contract.Approve(&_VariableDebtToken.TransactOpts, arg0, arg1)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address , uint256 ) returns(bool)
func (_VariableDebtToken *VariableDebtTokenTransactorSession) Approve(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _VariableDebtToken.Contract.Approve(&_VariableDebtToken.TransactOpts, arg0, arg1)
}

// ApproveDelegation is a paid mutator transaction binding the contract method 0xc04a8a10.
//
// Solidity: function approveDelegation(address delegatee, uint256 amount) returns()
func (_VariableDebtToken *VariableDebtTokenTransactor) ApproveDelegation(opts *bind.TransactOpts, delegatee common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VariableDebtToken.contract.Transact(opts, "approveDelegation", delegatee, amount)
}

// ApproveDelegation is a paid mutator transaction binding the contract method 0xc04a8a10.
//
// Solidity: function approveDelegation(address delegatee, uint256 amount) returns()
func (_VariableDebtToken *VariableDebtTokenSession) ApproveDelegation(delegatee common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VariableDebtToken.Contract.ApproveDelegation(&_VariableDebtToken.TransactOpts, delegatee, amount)
}

// ApproveDelegation is a paid mutator transaction binding the contract method 0xc04a8a10.
//
// Solidity: function approveDelegation(address delegatee, uint256 amount) returns()
func (_VariableDebtToken *VariableDebtTokenTransactorSession) ApproveDelegation(delegatee common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VariableDebtToken.Contract.ApproveDelegation(&_VariableDebtToken.TransactOpts, delegatee, amount)
}

// Burn is a paid mutator transaction binding the contract method 0xf5298aca.
//
// Solidity: function burn(address from, uint256 amount, uint256 index) returns(uint256)
func (_VariableDebtToken *VariableDebtTokenTransactor) Burn(opts *bind.TransactOpts, from common.Address, amount *big.Int, index *big.Int) (*types.Transaction, error) {
	return _VariableDebtToken.contract.Transact(opts, "burn", from, amount, index)
}

// Burn is a paid mutator transaction binding the contract method 0xf5298aca.
//
// Solidity: function burn(address from, uint256 amount, uint256 index) returns(uint256)
func (_VariableDebtToken *VariableDebtTokenSession) Burn(from common.Address, amount *big.Int, index *big.Int) (*types.Transaction, error) {
	return _VariableDebtToken.Contract.Burn(&_VariableDebtToken.TransactOpts, from, amount, index)
}

// Burn is a paid mutator transaction binding the contract method 0xf5298aca.
//
// Solidity: function burn(address from, uint256 amount, uint256 index) returns(uint256)
func (_VariableDebtToken *VariableDebtTokenTransactorSession) Burn(from common.Address, amount *big.Int, index *big.Int) (*types.Transaction, error) {
	return _VariableDebtToken.Contract.Burn(&_VariableDebtToken.TransactOpts, from, amount, index)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address , uint256 ) returns(bool)
func (_VariableDebtToken *VariableDebtTokenTransactor) DecreaseAllowance(opts *bind.TransactOpts, arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _VariableDebtToken.contract.Transact(opts, "decreaseAllowance", arg0, arg1)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address , uint256 ) returns(bool)
func (_VariableDebtToken *VariableDebtTokenSession) DecreaseAllowance(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _VariableDebtToken.Contract.DecreaseAllowance(&_VariableDebtToken.TransactOpts, arg0, arg1)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address , uint256 ) returns(bool)
func (_VariableDebtToken *VariableDebtTokenTransactorSession) DecreaseAllowance(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _VariableDebtToken.Contract.DecreaseAllowance(&_VariableDebtToken.TransactOpts, arg0, arg1)
}

// DelegationWithSig is a paid mutator transaction binding the contract method 0x0b52d558.
//
// Solidity: function delegationWithSig(address delegator, address delegatee, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_VariableDebtToken *VariableDebtTokenTransactor) DelegationWithSig(opts *bind.TransactOpts, delegator common.Address, delegatee common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _VariableDebtToken.contract.Transact(opts, "delegationWithSig", delegator, delegatee, value, deadline, v, r, s)
}

// DelegationWithSig is a paid mutator transaction binding the contract method 0x0b52d558.
//
// Solidity: function delegationWithSig(address delegator, address delegatee, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_VariableDebtToken *VariableDebtTokenSession) DelegationWithSig(delegator common.Address, delegatee common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _VariableDebtToken.Contract.DelegationWithSig(&_VariableDebtToken.TransactOpts, delegator, delegatee, value, deadline, v, r, s)
}

// DelegationWithSig is a paid mutator transaction binding the contract method 0x0b52d558.
//
// Solidity: function delegationWithSig(address delegator, address delegatee, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_VariableDebtToken *VariableDebtTokenTransactorSession) DelegationWithSig(delegator common.Address, delegatee common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _VariableDebtToken.Contract.DelegationWithSig(&_VariableDebtToken.TransactOpts, delegator, delegatee, value, deadline, v, r, s)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address , uint256 ) returns(bool)
func (_VariableDebtToken *VariableDebtTokenTransactor) IncreaseAllowance(opts *bind.TransactOpts, arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _VariableDebtToken.contract.Transact(opts, "increaseAllowance", arg0, arg1)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address , uint256 ) returns(bool)
func (_VariableDebtToken *VariableDebtTokenSession) IncreaseAllowance(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _VariableDebtToken.Contract.IncreaseAllowance(&_VariableDebtToken.TransactOpts, arg0, arg1)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address , uint256 ) returns(bool)
func (_VariableDebtToken *VariableDebtTokenTransactorSession) IncreaseAllowance(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _VariableDebtToken.Contract.IncreaseAllowance(&_VariableDebtToken.TransactOpts, arg0, arg1)
}

// Initialize is a paid mutator transaction binding the contract method 0xc222ec8a.
//
// Solidity: function initialize(address initializingPool, address underlyingAsset, address incentivesController, uint8 debtTokenDecimals, string debtTokenName, string debtTokenSymbol, bytes params) returns()
func (_VariableDebtToken *VariableDebtTokenTransactor) Initialize(opts *bind.TransactOpts, initializingPool common.Address, underlyingAsset common.Address, incentivesController common.Address, debtTokenDecimals uint8, debtTokenName string, debtTokenSymbol string, params []byte) (*types.Transaction, error) {
	return _VariableDebtToken.contract.Transact(opts, "initialize", initializingPool, underlyingAsset, incentivesController, debtTokenDecimals, debtTokenName, debtTokenSymbol, params)
}

// Initialize is a paid mutator transaction binding the contract method 0xc222ec8a.
//
// Solidity: function initialize(address initializingPool, address underlyingAsset, address incentivesController, uint8 debtTokenDecimals, string debtTokenName, string debtTokenSymbol, bytes params) returns()
func (_VariableDebtToken *VariableDebtTokenSession) Initialize(initializingPool common.Address, underlyingAsset common.Address, incentivesController common.Address, debtTokenDecimals uint8, debtTokenName string, debtTokenSymbol string, params []byte) (*types.Transaction, error) {
	return _VariableDebtToken.Contract.Initialize(&_VariableDebtToken.TransactOpts, initializingPool, underlyingAsset, incentivesController, debtTokenDecimals, debtTokenName, debtTokenSymbol, params)
}

// Initialize is a paid mutator transaction binding the contract method 0xc222ec8a.
//
// Solidity: function initialize(address initializingPool, address underlyingAsset, address incentivesController, uint8 debtTokenDecimals, string debtTokenName, string debtTokenSymbol, bytes params) returns()
func (_VariableDebtToken *VariableDebtTokenTransactorSession) Initialize(initializingPool common.Address, underlyingAsset common.Address, incentivesController common.Address, debtTokenDecimals uint8, debtTokenName string, debtTokenSymbol string, params []byte) (*types.Transaction, error) {
	return _VariableDebtToken.Contract.Initialize(&_VariableDebtToken.TransactOpts, initializingPool, underlyingAsset, incentivesController, debtTokenDecimals, debtTokenName, debtTokenSymbol, params)
}

// Mint is a paid mutator transaction binding the contract method 0xb3f1c93d.
//
// Solidity: function mint(address user, address onBehalfOf, uint256 amount, uint256 index) returns(bool, uint256)
func (_VariableDebtToken *VariableDebtTokenTransactor) Mint(opts *bind.TransactOpts, user common.Address, onBehalfOf common.Address, amount *big.Int, index *big.Int) (*types.Transaction, error) {
	return _VariableDebtToken.contract.Transact(opts, "mint", user, onBehalfOf, amount, index)
}

// Mint is a paid mutator transaction binding the contract method 0xb3f1c93d.
//
// Solidity: function mint(address user, address onBehalfOf, uint256 amount, uint256 index) returns(bool, uint256)
func (_VariableDebtToken *VariableDebtTokenSession) Mint(user common.Address, onBehalfOf common.Address, amount *big.Int, index *big.Int) (*types.Transaction, error) {
	return _VariableDebtToken.Contract.Mint(&_VariableDebtToken.TransactOpts, user, onBehalfOf, amount, index)
}

// Mint is a paid mutator transaction binding the contract method 0xb3f1c93d.
//
// Solidity: function mint(address user, address onBehalfOf, uint256 amount, uint256 index) returns(bool, uint256)
func (_VariableDebtToken *VariableDebtTokenTransactorSession) Mint(user common.Address, onBehalfOf common.Address, amount *big.Int, index *big.Int) (*types.Transaction, error) {
	return _VariableDebtToken.Contract.Mint(&_VariableDebtToken.TransactOpts, user, onBehalfOf, amount, index)
}

// SetIncentivesController is a paid mutator transaction binding the contract method 0xe655dbd8.
//
// Solidity: function setIncentivesController(address controller) returns()
func (_VariableDebtToken *VariableDebtTokenTransactor) SetIncentivesController(opts *bind.TransactOpts, controller common.Address) (*types.Transaction, error) {
	return _VariableDebtToken.contract.Transact(opts, "setIncentivesController", controller)
}

// SetIncentivesController is a paid mutator transaction binding the contract method 0xe655dbd8.
//
// Solidity: function setIncentivesController(address controller) returns()
func (_VariableDebtToken *VariableDebtTokenSession) SetIncentivesController(controller common.Address) (*types.Transaction, error) {
	return _VariableDebtToken.Contract.SetIncentivesController(&_VariableDebtToken.TransactOpts, controller)
}

// SetIncentivesController is a paid mutator transaction binding the contract method 0xe655dbd8.
//
// Solidity: function setIncentivesController(address controller) returns()
func (_VariableDebtToken *VariableDebtTokenTransactorSession) SetIncentivesController(controller common.Address) (*types.Transaction, error) {
	return _VariableDebtToken.Contract.SetIncentivesController(&_VariableDebtToken.TransactOpts, controller)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address , uint256 ) returns(bool)
func (_VariableDebtToken *VariableDebtTokenTransactor) Transfer(opts *bind.TransactOpts, arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _VariableDebtToken.contract.Transact(opts, "transfer", arg0, arg1)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address , uint256 ) returns(bool)
func (_VariableDebtToken *VariableDebtTokenSession) Transfer(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _VariableDebtToken.Contract.Transfer(&_VariableDebtToken.TransactOpts, arg0, arg1)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address , uint256 ) returns(bool)
func (_VariableDebtToken *VariableDebtTokenTransactorSession) Transfer(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _VariableDebtToken.Contract.Transfer(&_VariableDebtToken.TransactOpts, arg0, arg1)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address , address , uint256 ) returns(bool)
func (_VariableDebtToken *VariableDebtTokenTransactor) TransferFrom(opts *bind.TransactOpts, arg0 common.Address, arg1 common.Address, arg2 *big.Int) (*types.Transaction, error) {
	return _VariableDebtToken.contract.Transact(opts, "transferFrom", arg0, arg1, arg2)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address , address , uint256 ) returns(bool)
func (_VariableDebtToken *VariableDebtTokenSession) TransferFrom(arg0 common.Address, arg1 common.Address, arg2 *big.Int) (*types.Transaction, error) {
	return _VariableDebtToken.Contract.TransferFrom(&_VariableDebtToken.TransactOpts, arg0, arg1, arg2)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address , address , uint256 ) returns(bool)
func (_VariableDebtToken *VariableDebtTokenTransactorSession) TransferFrom(arg0 common.Address, arg1 common.Address, arg2 *big.Int) (*types.Transaction, error) {
	return _VariableDebtToken.Contract.TransferFrom(&_VariableDebtToken.TransactOpts, arg0, arg1, arg2)
}

// VariableDebtTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the VariableDebtToken contract.
type VariableDebtTokenApprovalIterator struct {
	Event *VariableDebtTokenApproval // Event containing the contract specifics and raw log

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
func (it *VariableDebtTokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VariableDebtTokenApproval)
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
		it.Event = new(VariableDebtTokenApproval)
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
func (it *VariableDebtTokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VariableDebtTokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VariableDebtTokenApproval represents a Approval event raised by the VariableDebtToken contract.
type VariableDebtTokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_VariableDebtToken *VariableDebtTokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*VariableDebtTokenApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _VariableDebtToken.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &VariableDebtTokenApprovalIterator{contract: _VariableDebtToken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_VariableDebtToken *VariableDebtTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *VariableDebtTokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _VariableDebtToken.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VariableDebtTokenApproval)
				if err := _VariableDebtToken.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_VariableDebtToken *VariableDebtTokenFilterer) ParseApproval(log types.Log) (*VariableDebtTokenApproval, error) {
	event := new(VariableDebtTokenApproval)
	if err := _VariableDebtToken.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VariableDebtTokenBorrowAllowanceDelegatedIterator is returned from FilterBorrowAllowanceDelegated and is used to iterate over the raw logs and unpacked data for BorrowAllowanceDelegated events raised by the VariableDebtToken contract.
type VariableDebtTokenBorrowAllowanceDelegatedIterator struct {
	Event *VariableDebtTokenBorrowAllowanceDelegated // Event containing the contract specifics and raw log

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
func (it *VariableDebtTokenBorrowAllowanceDelegatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VariableDebtTokenBorrowAllowanceDelegated)
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
		it.Event = new(VariableDebtTokenBorrowAllowanceDelegated)
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
func (it *VariableDebtTokenBorrowAllowanceDelegatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VariableDebtTokenBorrowAllowanceDelegatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VariableDebtTokenBorrowAllowanceDelegated represents a BorrowAllowanceDelegated event raised by the VariableDebtToken contract.
type VariableDebtTokenBorrowAllowanceDelegated struct {
	FromUser common.Address
	ToUser   common.Address
	Asset    common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterBorrowAllowanceDelegated is a free log retrieval operation binding the contract event 0xda919360433220e13b51e8c211e490d148e61a3bd53de8c097194e458b97f3e1.
//
// Solidity: event BorrowAllowanceDelegated(address indexed fromUser, address indexed toUser, address indexed asset, uint256 amount)
func (_VariableDebtToken *VariableDebtTokenFilterer) FilterBorrowAllowanceDelegated(opts *bind.FilterOpts, fromUser []common.Address, toUser []common.Address, asset []common.Address) (*VariableDebtTokenBorrowAllowanceDelegatedIterator, error) {

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

	logs, sub, err := _VariableDebtToken.contract.FilterLogs(opts, "BorrowAllowanceDelegated", fromUserRule, toUserRule, assetRule)
	if err != nil {
		return nil, err
	}
	return &VariableDebtTokenBorrowAllowanceDelegatedIterator{contract: _VariableDebtToken.contract, event: "BorrowAllowanceDelegated", logs: logs, sub: sub}, nil
}

// WatchBorrowAllowanceDelegated is a free log subscription operation binding the contract event 0xda919360433220e13b51e8c211e490d148e61a3bd53de8c097194e458b97f3e1.
//
// Solidity: event BorrowAllowanceDelegated(address indexed fromUser, address indexed toUser, address indexed asset, uint256 amount)
func (_VariableDebtToken *VariableDebtTokenFilterer) WatchBorrowAllowanceDelegated(opts *bind.WatchOpts, sink chan<- *VariableDebtTokenBorrowAllowanceDelegated, fromUser []common.Address, toUser []common.Address, asset []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _VariableDebtToken.contract.WatchLogs(opts, "BorrowAllowanceDelegated", fromUserRule, toUserRule, assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VariableDebtTokenBorrowAllowanceDelegated)
				if err := _VariableDebtToken.contract.UnpackLog(event, "BorrowAllowanceDelegated", log); err != nil {
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
func (_VariableDebtToken *VariableDebtTokenFilterer) ParseBorrowAllowanceDelegated(log types.Log) (*VariableDebtTokenBorrowAllowanceDelegated, error) {
	event := new(VariableDebtTokenBorrowAllowanceDelegated)
	if err := _VariableDebtToken.contract.UnpackLog(event, "BorrowAllowanceDelegated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VariableDebtTokenBurnIterator is returned from FilterBurn and is used to iterate over the raw logs and unpacked data for Burn events raised by the VariableDebtToken contract.
type VariableDebtTokenBurnIterator struct {
	Event *VariableDebtTokenBurn // Event containing the contract specifics and raw log

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
func (it *VariableDebtTokenBurnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VariableDebtTokenBurn)
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
		it.Event = new(VariableDebtTokenBurn)
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
func (it *VariableDebtTokenBurnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VariableDebtTokenBurnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VariableDebtTokenBurn represents a Burn event raised by the VariableDebtToken contract.
type VariableDebtTokenBurn struct {
	From            common.Address
	Target          common.Address
	Value           *big.Int
	BalanceIncrease *big.Int
	Index           *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterBurn is a free log retrieval operation binding the contract event 0x4cf25bc1d991c17529c25213d3cc0cda295eeaad5f13f361969b12ea48015f90.
//
// Solidity: event Burn(address indexed from, address indexed target, uint256 value, uint256 balanceIncrease, uint256 index)
func (_VariableDebtToken *VariableDebtTokenFilterer) FilterBurn(opts *bind.FilterOpts, from []common.Address, target []common.Address) (*VariableDebtTokenBurnIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var targetRule []interface{}
	for _, targetItem := range target {
		targetRule = append(targetRule, targetItem)
	}

	logs, sub, err := _VariableDebtToken.contract.FilterLogs(opts, "Burn", fromRule, targetRule)
	if err != nil {
		return nil, err
	}
	return &VariableDebtTokenBurnIterator{contract: _VariableDebtToken.contract, event: "Burn", logs: logs, sub: sub}, nil
}

// WatchBurn is a free log subscription operation binding the contract event 0x4cf25bc1d991c17529c25213d3cc0cda295eeaad5f13f361969b12ea48015f90.
//
// Solidity: event Burn(address indexed from, address indexed target, uint256 value, uint256 balanceIncrease, uint256 index)
func (_VariableDebtToken *VariableDebtTokenFilterer) WatchBurn(opts *bind.WatchOpts, sink chan<- *VariableDebtTokenBurn, from []common.Address, target []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var targetRule []interface{}
	for _, targetItem := range target {
		targetRule = append(targetRule, targetItem)
	}

	logs, sub, err := _VariableDebtToken.contract.WatchLogs(opts, "Burn", fromRule, targetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VariableDebtTokenBurn)
				if err := _VariableDebtToken.contract.UnpackLog(event, "Burn", log); err != nil {
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

// ParseBurn is a log parse operation binding the contract event 0x4cf25bc1d991c17529c25213d3cc0cda295eeaad5f13f361969b12ea48015f90.
//
// Solidity: event Burn(address indexed from, address indexed target, uint256 value, uint256 balanceIncrease, uint256 index)
func (_VariableDebtToken *VariableDebtTokenFilterer) ParseBurn(log types.Log) (*VariableDebtTokenBurn, error) {
	event := new(VariableDebtTokenBurn)
	if err := _VariableDebtToken.contract.UnpackLog(event, "Burn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VariableDebtTokenInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the VariableDebtToken contract.
type VariableDebtTokenInitializedIterator struct {
	Event *VariableDebtTokenInitialized // Event containing the contract specifics and raw log

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
func (it *VariableDebtTokenInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VariableDebtTokenInitialized)
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
		it.Event = new(VariableDebtTokenInitialized)
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
func (it *VariableDebtTokenInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VariableDebtTokenInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VariableDebtTokenInitialized represents a Initialized event raised by the VariableDebtToken contract.
type VariableDebtTokenInitialized struct {
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
func (_VariableDebtToken *VariableDebtTokenFilterer) FilterInitialized(opts *bind.FilterOpts, underlyingAsset []common.Address, pool []common.Address) (*VariableDebtTokenInitializedIterator, error) {

	var underlyingAssetRule []interface{}
	for _, underlyingAssetItem := range underlyingAsset {
		underlyingAssetRule = append(underlyingAssetRule, underlyingAssetItem)
	}
	var poolRule []interface{}
	for _, poolItem := range pool {
		poolRule = append(poolRule, poolItem)
	}

	logs, sub, err := _VariableDebtToken.contract.FilterLogs(opts, "Initialized", underlyingAssetRule, poolRule)
	if err != nil {
		return nil, err
	}
	return &VariableDebtTokenInitializedIterator{contract: _VariableDebtToken.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x40251fbfb6656cfa65a00d7879029fec1fad21d28fdcff2f4f68f52795b74f2c.
//
// Solidity: event Initialized(address indexed underlyingAsset, address indexed pool, address incentivesController, uint8 debtTokenDecimals, string debtTokenName, string debtTokenSymbol, bytes params)
func (_VariableDebtToken *VariableDebtTokenFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *VariableDebtTokenInitialized, underlyingAsset []common.Address, pool []common.Address) (event.Subscription, error) {

	var underlyingAssetRule []interface{}
	for _, underlyingAssetItem := range underlyingAsset {
		underlyingAssetRule = append(underlyingAssetRule, underlyingAssetItem)
	}
	var poolRule []interface{}
	for _, poolItem := range pool {
		poolRule = append(poolRule, poolItem)
	}

	logs, sub, err := _VariableDebtToken.contract.WatchLogs(opts, "Initialized", underlyingAssetRule, poolRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VariableDebtTokenInitialized)
				if err := _VariableDebtToken.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_VariableDebtToken *VariableDebtTokenFilterer) ParseInitialized(log types.Log) (*VariableDebtTokenInitialized, error) {
	event := new(VariableDebtTokenInitialized)
	if err := _VariableDebtToken.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VariableDebtTokenMintIterator is returned from FilterMint and is used to iterate over the raw logs and unpacked data for Mint events raised by the VariableDebtToken contract.
type VariableDebtTokenMintIterator struct {
	Event *VariableDebtTokenMint // Event containing the contract specifics and raw log

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
func (it *VariableDebtTokenMintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VariableDebtTokenMint)
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
		it.Event = new(VariableDebtTokenMint)
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
func (it *VariableDebtTokenMintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VariableDebtTokenMintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VariableDebtTokenMint represents a Mint event raised by the VariableDebtToken contract.
type VariableDebtTokenMint struct {
	Caller          common.Address
	OnBehalfOf      common.Address
	Value           *big.Int
	BalanceIncrease *big.Int
	Index           *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterMint is a free log retrieval operation binding the contract event 0x458f5fa412d0f69b08dd84872b0215675cc67bc1d5b6fd93300a1c3878b86196.
//
// Solidity: event Mint(address indexed caller, address indexed onBehalfOf, uint256 value, uint256 balanceIncrease, uint256 index)
func (_VariableDebtToken *VariableDebtTokenFilterer) FilterMint(opts *bind.FilterOpts, caller []common.Address, onBehalfOf []common.Address) (*VariableDebtTokenMintIterator, error) {

	var callerRule []interface{}
	for _, callerItem := range caller {
		callerRule = append(callerRule, callerItem)
	}
	var onBehalfOfRule []interface{}
	for _, onBehalfOfItem := range onBehalfOf {
		onBehalfOfRule = append(onBehalfOfRule, onBehalfOfItem)
	}

	logs, sub, err := _VariableDebtToken.contract.FilterLogs(opts, "Mint", callerRule, onBehalfOfRule)
	if err != nil {
		return nil, err
	}
	return &VariableDebtTokenMintIterator{contract: _VariableDebtToken.contract, event: "Mint", logs: logs, sub: sub}, nil
}

// WatchMint is a free log subscription operation binding the contract event 0x458f5fa412d0f69b08dd84872b0215675cc67bc1d5b6fd93300a1c3878b86196.
//
// Solidity: event Mint(address indexed caller, address indexed onBehalfOf, uint256 value, uint256 balanceIncrease, uint256 index)
func (_VariableDebtToken *VariableDebtTokenFilterer) WatchMint(opts *bind.WatchOpts, sink chan<- *VariableDebtTokenMint, caller []common.Address, onBehalfOf []common.Address) (event.Subscription, error) {

	var callerRule []interface{}
	for _, callerItem := range caller {
		callerRule = append(callerRule, callerItem)
	}
	var onBehalfOfRule []interface{}
	for _, onBehalfOfItem := range onBehalfOf {
		onBehalfOfRule = append(onBehalfOfRule, onBehalfOfItem)
	}

	logs, sub, err := _VariableDebtToken.contract.WatchLogs(opts, "Mint", callerRule, onBehalfOfRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VariableDebtTokenMint)
				if err := _VariableDebtToken.contract.UnpackLog(event, "Mint", log); err != nil {
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

// ParseMint is a log parse operation binding the contract event 0x458f5fa412d0f69b08dd84872b0215675cc67bc1d5b6fd93300a1c3878b86196.
//
// Solidity: event Mint(address indexed caller, address indexed onBehalfOf, uint256 value, uint256 balanceIncrease, uint256 index)
func (_VariableDebtToken *VariableDebtTokenFilterer) ParseMint(log types.Log) (*VariableDebtTokenMint, error) {
	event := new(VariableDebtTokenMint)
	if err := _VariableDebtToken.contract.UnpackLog(event, "Mint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VariableDebtTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the VariableDebtToken contract.
type VariableDebtTokenTransferIterator struct {
	Event *VariableDebtTokenTransfer // Event containing the contract specifics and raw log

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
func (it *VariableDebtTokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VariableDebtTokenTransfer)
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
		it.Event = new(VariableDebtTokenTransfer)
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
func (it *VariableDebtTokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VariableDebtTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VariableDebtTokenTransfer represents a Transfer event raised by the VariableDebtToken contract.
type VariableDebtTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_VariableDebtToken *VariableDebtTokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VariableDebtTokenTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VariableDebtToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VariableDebtTokenTransferIterator{contract: _VariableDebtToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_VariableDebtToken *VariableDebtTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *VariableDebtTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VariableDebtToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VariableDebtTokenTransfer)
				if err := _VariableDebtToken.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_VariableDebtToken *VariableDebtTokenFilterer) ParseTransfer(log types.Log) (*VariableDebtTokenTransfer, error) {
	event := new(VariableDebtTokenTransfer)
	if err := _VariableDebtToken.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
