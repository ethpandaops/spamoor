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

// IETHRegistrarControllerRegistration is an auto generated low-level Go binding around an user-defined struct.
type IETHRegistrarControllerRegistration struct {
	Label         string
	Owner         common.Address
	Duration      *big.Int
	Secret        [32]byte
	Resolver      common.Address
	Data          [][]byte
	ReverseRecord uint8
	Referrer      [32]byte
}

// IPriceOraclePrice is an auto generated low-level Go binding around an user-defined struct.
type IPriceOraclePrice struct {
	Base    *big.Int
	Premium *big.Int
}

// ETHRegistrarControllerMetaData contains all meta data concerning the ETHRegistrarController contract.
var ETHRegistrarControllerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractBaseRegistrarImplementation\",\"name\":\"_base\",\"type\":\"address\"},{\"internalType\":\"contractIPriceOracle\",\"name\":\"_prices\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_minCommitmentAge\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_maxCommitmentAge\",\"type\":\"uint256\"},{\"internalType\":\"contractIReverseRegistrar\",\"name\":\"_reverseRegistrar\",\"type\":\"address\"},{\"internalType\":\"contractIDefaultReverseRegistrar\",\"name\":\"_defaultReverseRegistrar\",\"type\":\"address\"},{\"internalType\":\"contractENS\",\"name\":\"_ens\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"commitment\",\"type\":\"bytes32\"}],\"name\":\"CommitmentNotFound\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"commitment\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"minimumCommitmentTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"currentTimestamp\",\"type\":\"uint256\"}],\"name\":\"CommitmentTooNew\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"commitment\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"maximumCommitmentTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"currentTimestamp\",\"type\":\"uint256\"}],\"name\":\"CommitmentTooOld\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"duration\",\"type\":\"uint256\"}],\"name\":\"DurationTooShort\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientValue\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxCommitmentAgeTooHigh\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxCommitmentAgeTooLow\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"NameNotAvailable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ResolverRequiredForReverseRecord\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ResolverRequiredWhenDataSupplied\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"commitment\",\"type\":\"bytes32\"}],\"name\":\"UnexpiredCommitmentExists\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"label\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"labelhash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"baseCost\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"premium\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"expires\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"referrer\",\"type\":\"bytes32\"}],\"name\":\"NameRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"label\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"labelhash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"cost\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"expires\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"referrer\",\"type\":\"bytes32\"}],\"name\":\"NameRenewed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"MIN_REGISTRATION_DURATION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"label\",\"type\":\"string\"}],\"name\":\"available\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"commitment\",\"type\":\"bytes32\"}],\"name\":\"commit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"commitments\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"defaultReverseRegistrar\",\"outputs\":[{\"internalType\":\"contractIDefaultReverseRegistrar\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ens\",\"outputs\":[{\"internalType\":\"contractENS\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"label\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"duration\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"secret\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"resolver\",\"type\":\"address\"},{\"internalType\":\"bytes[]\",\"name\":\"data\",\"type\":\"bytes[]\"},{\"internalType\":\"uint8\",\"name\":\"reverseRecord\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"referrer\",\"type\":\"bytes32\"}],\"internalType\":\"structIETHRegistrarController.Registration\",\"name\":\"registration\",\"type\":\"tuple\"}],\"name\":\"makeCommitment\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"commitment\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxCommitmentAge\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minCommitmentAge\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"prices\",\"outputs\":[{\"internalType\":\"contractIPriceOracle\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"recoverFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"label\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"duration\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"secret\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"resolver\",\"type\":\"address\"},{\"internalType\":\"bytes[]\",\"name\":\"data\",\"type\":\"bytes[]\"},{\"internalType\":\"uint8\",\"name\":\"reverseRecord\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"referrer\",\"type\":\"bytes32\"}],\"internalType\":\"structIETHRegistrarController.Registration\",\"name\":\"registration\",\"type\":\"tuple\"}],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"label\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"duration\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"referrer\",\"type\":\"bytes32\"}],\"name\":\"renew\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"label\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"duration\",\"type\":\"uint256\"}],\"name\":\"rentPrice\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"base\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"premium\",\"type\":\"uint256\"}],\"internalType\":\"structIPriceOracle.Price\",\"name\":\"price\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reverseRegistrar\",\"outputs\":[{\"internalType\":\"contractIReverseRegistrar\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceID\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"label\",\"type\":\"string\"}],\"name\":\"valid\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x61016060405234801561001157600080fd5b5060405161286638038061286683398101604081905261003091610116565b610039336100ae565b848411610059576040516307cb550760e31b815260040160405180910390fd5b4284111561007a57604051630b4319e560e21b815260040160405180910390fd5b6001600160a01b0390811660805295861660a0529385166101405260c09290925260e052821661010052166101205261019f565b600080546001600160a01b038381166001600160a01b0319831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b6001600160a01b038116811461011357600080fd5b50565b600080600080600080600060e0888a03121561013157600080fd5b875161013c816100fe565b602089015190975061014d816100fe565b604089015160608a015160808b0151929850909650945061016d816100fe565b60a089015190935061017e816100fe565b60c089015190925061018f816100fe565b8091505092959891949750929550565b60805160a05160c05160e05161010051610120516101405161260b61025b6000396000818161045601526115ac01526000818161023901526112e30152600081816102a201526111ea01526000818161040201528181610bfb01528181610c68015261142e01526000818161036301528181610b5e0152610b8f01526000818161060501528181610d2201528181610e3e015281816110f2015281816115da01526119e40152600081816101e00152610f13015261260b6000f3fe60806040526004361061016a5760003560e01c80638a95b09f116100cb578063ce1e09c01161007f578063ef9c880511610059578063ef9c880514610478578063f14fcbc81461048b578063f2fde38b146104ab57600080fd5b8063ce1e09c0146103f0578063cf7d6e0114610424578063d3419bf31461044457600080fd5b80638da5cb5b116100b05780638da5cb5b146103855780639791c097146103b0578063aeb8ce9b146103d057600080fd5b80638a95b09f1461033a5780638d839ffe1461035157600080fd5b80635d3590d51161012257806380869853116101075780638086985314610290578063839df945146102c457806383e7f6ff146102ff57600080fd5b80635d3590d51461025b578063715018a61461027b57600080fd5b80633ccfd60b116101535780633ccfd60b146101b95780633f15457f146101ce578063469bf4411461022757600080fd5b806301ffc9a71461016f57806318026ad1146101a4575b600080fd5b34801561017b57600080fd5b5061018f61018a366004611a64565b6104cb565b60405190151581526020015b60405180910390f35b6101b76101b2366004611af6565b610564565b005b3480156101c557600080fd5b506101b761071b565b3480156101da57600080fd5b506102027f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161019b565b34801561023357600080fd5b506102027f000000000000000000000000000000000000000000000000000000000000000081565b34801561026757600080fd5b506101b7610276366004611b70565b610765565b34801561028757600080fd5b506101b761080c565b34801561029c57600080fd5b506102027f000000000000000000000000000000000000000000000000000000000000000081565b3480156102d057600080fd5b506102f16102df366004611bad565b60016020526000908152604090205481565b60405190815260200161019b565b34801561030b57600080fd5b5061031f61031a366004611bc6565b610820565b6040805182518152602092830151928101929092520161019b565b34801561034657600080fd5b506102f16224ea0081565b34801561035d57600080fd5b506102f17f000000000000000000000000000000000000000000000000000000000000000081565b34801561039157600080fd5b5060005473ffffffffffffffffffffffffffffffffffffffff16610202565b3480156103bc57600080fd5b5061018f6103cb366004611c12565b610865565b3480156103dc57600080fd5b5061018f6103eb366004611c12565b6108b1565b3480156103fc57600080fd5b506102f17f000000000000000000000000000000000000000000000000000000000000000081565b34801561043057600080fd5b506102f161043f366004611c54565b6108e1565b34801561045057600080fd5b506102027f000000000000000000000000000000000000000000000000000000000000000081565b6101b7610486366004611c54565b610a53565b34801561049757600080fd5b506101b76104a6366004611bad565b611417565b3480156104b757600080fd5b506101b76104c6366004611c90565b6114a0565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167fe4f37f7900000000000000000000000000000000000000000000000000000000148061055e57507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008316145b92915050565b60008484604051610576929190611cab565b60405180910390209050600061058e86868487611554565b80519091503410156105cc576040517f1101129400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040517fc475abff00000000000000000000000000000000000000000000000000000000815260048101839052602481018590526000907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff169063c475abff906044016020604051808303816000875af1158015610663573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106879190611cbb565b9050827ffa956c3bce4cb4b01166868ecaf0620566bc7e33fc70b0b9c6aef61e37e50b948888856000015185896040516106c5959493929190611d1d565b60405180910390a2815134111561071257815133906108fc906106e89034611d7d565b6040518115909202916000818181858888f19350505050158015610710573d6000803e3d6000fd5b505b50505050505050565b6000805460405173ffffffffffffffffffffffffffffffffffffffff909116914780156108fc02929091818181858888f19350505050158015610762573d6000803e3d6000fd5b50565b61076d6116a5565b6040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff83811660048301526024820183905284169063a9059cbb906044016020604051808303816000875af11580156107e2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108069190611d90565b50505050565b6108146116a5565b61081e6000611726565b565b604080518082019091526000808252602082015260008484604051610846929190611cab565b6040518091039020905061085c85858386611554565b95945050505050565b600060036108a884848080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061179b92505050565b10159392505050565b60008083836040516108c4929190611cab565b604051809103902090506108d98484836119a2565b949350505050565b6000806108f160a0840184611db2565b90501180156109255750600061090d60a0840160808501611c90565b73ffffffffffffffffffffffffffffffffffffffff16145b1561095c576040517fd3f605c400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61096c60e0830160c08401611e2b565b60ff16158015906109a25750600061098a60a0840160808501611c90565b73ffffffffffffffffffffffffffffffffffffffff16145b156109d9576040517f7d4a034a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6224ea0082604001351015610a2557604080517f9a71997b0000000000000000000000000000000000000000000000000000000081529083013560048201526024015b60405180910390fd5b81604051602001610a369190611f92565b604051602081830303815290604052805190602001209050919050565b6000610a5f828061209e565b604051610a6d929190611cab565b60405190819003902090506000610a92610a87848061209e565b848660400135611554565b9050600081602001518260000151610aaa9190612103565b905080341015610ae6576040517f1101129400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610af9610af3858061209e565b856119a2565b610b3b57610b07848061209e565b6040517f477707e8000000000000000000000000000000000000000000000000000000008152600401610a1c929190612116565b6000610b46856108e1565b60008181526001602052604090205490915042610b837f000000000000000000000000000000000000000000000000000000000000000083612103565b1115610bf55781610bb47f000000000000000000000000000000000000000000000000000000000000000083612103565b6040517f74480cc900000000000000000000000000000000000000000000000000000000815260048101929092526024820152426044820152606401610a1c565b42610c207f000000000000000000000000000000000000000000000000000000000000000083612103565b11610cce5780600003610c62576040517f836588c900000000000000000000000000000000000000000000000000000000815260048101839052602401610a1c565b81610c8d7f000000000000000000000000000000000000000000000000000000000000000083612103565b6040517f256e221600000000000000000000000000000000000000000000000000000000815260048101929092526024820152426044820152606401610a1c565b600082815260016020526040812081905580610cf060a0890160808a01611c90565b73ffffffffffffffffffffffffffffffffffffffff1603610dff5773ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001663fca247ac87610d5860408b0160208c01611c90565b604080517fffffffff0000000000000000000000000000000000000000000000000000000060e086901b168152600481019390935273ffffffffffffffffffffffffffffffffffffffff90911660248301528a013560448201526064016020604051808303816000875af1158015610dd4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610df89190611cbb565b9050611385565b604080517ffca247ac000000000000000000000000000000000000000000000000000000008152600481018890523060248201529088013560448201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff169063fca247ac906064016020604051808303816000875af1158015610e9c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ec09190611cbb565b604080517f93cdeb708b7545dc668eb9280176169d1c33cfd8ed6f04690a0bcc88a93fc4ae60208201529081018890529091506000906060016040516020818303038152906040528051906020012090507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663cf408823828a6020016020810190610f619190611c90565b610f7160a08d0160808e01611c90565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e086901b168152600481019390935273ffffffffffffffffffffffffffffffffffffffff918216602484015216604482015260006064820152608401600060405180830381600087803b158015610fec57600080fd5b505af1158015611000573d6000803e3d6000fd5b506000925061101591505060a08a018a611db2565b905011156110db5761102d60a0890160808a01611c90565b73ffffffffffffffffffffffffffffffffffffffff1663e32954eb8261105660a08c018c611db2565b6040518463ffffffff1660e01b81526004016110749392919061212a565b6000604051808303816000875af1158015611093573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526110d99190810190612258565b505b73ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000166323b872dd3061112860408c0160208d01611c90565b60405160e084901b7fffffffff0000000000000000000000000000000000000000000000000000000016815273ffffffffffffffffffffffffffffffffffffffff928316600482015291166024820152604481018a9052606401600060405180830381600087803b15801561119c57600080fd5b505af11580156111b0573d6000803e3d6000fd5b50600192506111c891505060e08a0160c08b01611e2b565b1660ff166000146112af5773ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016637a806d6b338061122160a08d0160808e01611c90565b61122b8d8061209e565b60405160200161123c9291906123a1565b6040516020818303038152906040526040518563ffffffff1660e01b815260040161126a949392919061241d565b6020604051808303816000875af1158015611289573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112ad9190611cbb565b505b60026112c160e08a0160c08b01611e2b565b1660ff166000146113835773ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001663c9119941336113138b8061209e565b6040516020016113249291906123a1565b6040516020818303038152906040526040518363ffffffff1660e01b815260040161135092919061248e565b600060405180830381600087803b15801561136a57600080fd5b505af115801561137e573d6000803e3d6000fd5b505050505b505b6113956040880160208901611c90565b73ffffffffffffffffffffffffffffffffffffffff16867fc2240194853531f1ae318dcef227de79c6ad0fd9d1b0e4fe08568415be2e08a56113d78a8061209e565b89600001518a60200151878e60e001356040516113f9969594939291906124bd565b60405180910390a38334111561071257336108fc6106e88634611d7d565b6000818152600160205260409020544290611453907f000000000000000000000000000000000000000000000000000000000000000090612103565b1061148d576040517f0a059d7100000000000000000000000000000000000000000000000000000000815260048101829052602401610a1c565b6000908152600160205260409020429055565b6114a86116a5565b73ffffffffffffffffffffffffffffffffffffffff811661154b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201527f64647265737300000000000000000000000000000000000000000000000000006064820152608401610a1c565b61076281611726565b60408051808201909152600080825260208201526040517fd6e4fa860000000000000000000000000000000000000000000000000000000081526004810184905273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000008116916350e9a71591889188917f0000000000000000000000000000000000000000000000000000000000000000169063d6e4fa8690602401602060405180830381865afa158015611621573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906116459190611cbb565b866040518563ffffffff1660e01b815260040161166594939291906124f6565b6040805180830381865afa158015611681573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061085c919061251d565b60005473ffffffffffffffffffffffffffffffffffffffff16331461081e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610a1c565b6000805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff0000000000000000000000000000000000000000831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b8051600090819081905b808210156119995760008583815181106117c1576117c161256e565b01602001517fff000000000000000000000000000000000000000000000000000000000000001690507f80000000000000000000000000000000000000000000000000000000000000008110156118245761181d600184612103565b9250611986565b7fe0000000000000000000000000000000000000000000000000000000000000007fff00000000000000000000000000000000000000000000000000000000000000821610156118795761181d600284612103565b7ff0000000000000000000000000000000000000000000000000000000000000007fff00000000000000000000000000000000000000000000000000000000000000821610156118ce5761181d600384612103565b7ff8000000000000000000000000000000000000000000000000000000000000007fff00000000000000000000000000000000000000000000000000000000000000821610156119235761181d600484612103565b7ffc000000000000000000000000000000000000000000000000000000000000007fff00000000000000000000000000000000000000000000000000000000000000821610156119785761181d600584612103565b611983600684612103565b92505b50826119918161259d565b9350506117a5565b50909392505050565b60006119ae8484610865565b80156108d957506040517f96e494e8000000000000000000000000000000000000000000000000000000008152600481018390527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906396e494e890602401602060405180830381865afa158015611a40573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108d99190611d90565b600060208284031215611a7657600080fd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114611aa657600080fd5b9392505050565b60008083601f840112611abf57600080fd5b50813567ffffffffffffffff811115611ad757600080fd5b602083019150836020828501011115611aef57600080fd5b9250929050565b60008060008060608587031215611b0c57600080fd5b843567ffffffffffffffff811115611b2357600080fd5b611b2f87828801611aad565b90989097506020870135966040013595509350505050565b803573ffffffffffffffffffffffffffffffffffffffff81168114611b6b57600080fd5b919050565b600080600060608486031215611b8557600080fd5b611b8e84611b47565b9250611b9c60208501611b47565b929592945050506040919091013590565b600060208284031215611bbf57600080fd5b5035919050565b600080600060408486031215611bdb57600080fd5b833567ffffffffffffffff811115611bf257600080fd5b611bfe86828701611aad565b909790965060209590950135949350505050565b60008060208385031215611c2557600080fd5b823567ffffffffffffffff811115611c3c57600080fd5b611c4885828601611aad565b90969095509350505050565b600060208284031215611c6657600080fd5b813567ffffffffffffffff811115611c7d57600080fd5b82016101008185031215611aa657600080fd5b600060208284031215611ca257600080fd5b611aa682611b47565b8183823760009101908152919050565b600060208284031215611ccd57600080fd5b5051919050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b608081526000611d31608083018789611cd4565b602083019590955250604081019290925260609091015292915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8181038181111561055e5761055e611d4e565b600060208284031215611da257600080fd5b81518015158114611aa657600080fd5b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112611de757600080fd5b83018035915067ffffffffffffffff821115611e0257600080fd5b6020019150600581901b3603821315611aef57600080fd5b803560ff81168114611b6b57600080fd5b600060208284031215611e3d57600080fd5b611aa682611e1a565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112611e7b57600080fd5b830160208101925035905067ffffffffffffffff811115611e9b57600080fd5b803603821315611aef57600080fd5b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112611edf57600080fd5b830160208101925035905067ffffffffffffffff811115611eff57600080fd5b8060051b3603821315611aef57600080fd5b60008383855260208501945060208460051b8201018360005b86811015611f86577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0848403018852611f638287611e46565b611f6e858284611cd4565b60209a8b019a90955093909301925050600101611f2a565b50909695505050505050565b602081526000611fa28384611e46565b6101006020850152611fb961012085018284611cd4565b91505073ffffffffffffffffffffffffffffffffffffffff611fdd60208601611b47565b166040840152600060408501359050806060850152506000606085013590508060808501525061200f60808501611b47565b73ffffffffffffffffffffffffffffffffffffffff811660a08501525061203960a0850185611eaa565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08584030160c086015261206e838284611f11565b9250505061207e60c08501611e1a565b60ff811660e08501525060e0939093013561010092909201919091525090565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18436030181126120d357600080fd5b83018035915067ffffffffffffffff8211156120ee57600080fd5b602001915036819003821315611aef57600080fd5b8082018082111561055e5761055e611d4e565b6020815260006108d9602083018486611cd4565b838152604060208201819052810182905260006060600584901b8301810190830185835b868110156121a9577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa08685030183526121878289611e46565b612192868284611cd4565b95505050602092830192919091019060010161214e565b5091979650505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561222c5761222c6121b6565b604052919050565b60005b8381101561224f578181015183820152602001612237565b50506000910152565b60006020828403121561226a57600080fd5b815167ffffffffffffffff81111561228157600080fd5b8201601f8101841361229257600080fd5b805167ffffffffffffffff8111156122ac576122ac6121b6565b8060051b6122bc602082016121e5565b918252602081840181019290810190878411156122d857600080fd5b6020850192505b8383101561239657825167ffffffffffffffff8111156122fe57600080fd5b8501603f8101891361230f57600080fd5b602081015167ffffffffffffffff81111561232c5761232c6121b6565b61235d60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116016121e5565b8181526040838301018b101561237257600080fd5b612383826020830160408601612234565b84525050602092830192909101906122df565b979650505050505050565b818382377f2e657468000000000000000000000000000000000000000000000000000000009101908152600401919050565b600081518084526123eb816020860160208601612234565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b73ffffffffffffffffffffffffffffffffffffffff8516815273ffffffffffffffffffffffffffffffffffffffff8416602082015273ffffffffffffffffffffffffffffffffffffffff8316604082015260806060820152600061248460808301846123d3565b9695505050505050565b73ffffffffffffffffffffffffffffffffffffffff831681526040602082015260006108d960408301846123d3565b60a0815260006124d160a08301888a611cd4565b9050856020830152846040830152836060830152826080830152979650505050505050565b60608152600061250a606083018688611cd4565b6020830194909452506040015292915050565b6000604082840312801561253057600080fd5b506040805190810167ffffffffffffffff81118282101715612554576125546121b6565b604052825181526020928301519281019290925250919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036125ce576125ce611d4e565b506001019056fea26469706673582212208623dbfc34fd664dc71ace35b12f78220e7039b5ce82fcaa82211441bb9265e364736f6c634300081a0033",
}

// ETHRegistrarControllerABI is the input ABI used to generate the binding from.
// Deprecated: Use ETHRegistrarControllerMetaData.ABI instead.
var ETHRegistrarControllerABI = ETHRegistrarControllerMetaData.ABI

// ETHRegistrarControllerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ETHRegistrarControllerMetaData.Bin instead.
var ETHRegistrarControllerBin = ETHRegistrarControllerMetaData.Bin

// DeployETHRegistrarController deploys a new Ethereum contract, binding an instance of ETHRegistrarController to it.
func DeployETHRegistrarController(auth *bind.TransactOpts, backend bind.ContractBackend, _base common.Address, _prices common.Address, _minCommitmentAge *big.Int, _maxCommitmentAge *big.Int, _reverseRegistrar common.Address, _defaultReverseRegistrar common.Address, _ens common.Address) (common.Address, *types.Transaction, *ETHRegistrarController, error) {
	parsed, err := ETHRegistrarControllerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ETHRegistrarControllerBin), backend, _base, _prices, _minCommitmentAge, _maxCommitmentAge, _reverseRegistrar, _defaultReverseRegistrar, _ens)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ETHRegistrarController{ETHRegistrarControllerCaller: ETHRegistrarControllerCaller{contract: contract}, ETHRegistrarControllerTransactor: ETHRegistrarControllerTransactor{contract: contract}, ETHRegistrarControllerFilterer: ETHRegistrarControllerFilterer{contract: contract}}, nil
}

// ETHRegistrarController is an auto generated Go binding around an Ethereum contract.
type ETHRegistrarController struct {
	ETHRegistrarControllerCaller     // Read-only binding to the contract
	ETHRegistrarControllerTransactor // Write-only binding to the contract
	ETHRegistrarControllerFilterer   // Log filterer for contract events
}

// ETHRegistrarControllerCaller is an auto generated read-only Go binding around an Ethereum contract.
type ETHRegistrarControllerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ETHRegistrarControllerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ETHRegistrarControllerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ETHRegistrarControllerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ETHRegistrarControllerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ETHRegistrarControllerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ETHRegistrarControllerSession struct {
	Contract     *ETHRegistrarController // Generic contract binding to set the session for
	CallOpts     bind.CallOpts           // Call options to use throughout this session
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// ETHRegistrarControllerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ETHRegistrarControllerCallerSession struct {
	Contract *ETHRegistrarControllerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                 // Call options to use throughout this session
}

// ETHRegistrarControllerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ETHRegistrarControllerTransactorSession struct {
	Contract     *ETHRegistrarControllerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// ETHRegistrarControllerRaw is an auto generated low-level Go binding around an Ethereum contract.
type ETHRegistrarControllerRaw struct {
	Contract *ETHRegistrarController // Generic contract binding to access the raw methods on
}

// ETHRegistrarControllerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ETHRegistrarControllerCallerRaw struct {
	Contract *ETHRegistrarControllerCaller // Generic read-only contract binding to access the raw methods on
}

// ETHRegistrarControllerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ETHRegistrarControllerTransactorRaw struct {
	Contract *ETHRegistrarControllerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewETHRegistrarController creates a new instance of ETHRegistrarController, bound to a specific deployed contract.
func NewETHRegistrarController(address common.Address, backend bind.ContractBackend) (*ETHRegistrarController, error) {
	contract, err := bindETHRegistrarController(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ETHRegistrarController{ETHRegistrarControllerCaller: ETHRegistrarControllerCaller{contract: contract}, ETHRegistrarControllerTransactor: ETHRegistrarControllerTransactor{contract: contract}, ETHRegistrarControllerFilterer: ETHRegistrarControllerFilterer{contract: contract}}, nil
}

// NewETHRegistrarControllerCaller creates a new read-only instance of ETHRegistrarController, bound to a specific deployed contract.
func NewETHRegistrarControllerCaller(address common.Address, caller bind.ContractCaller) (*ETHRegistrarControllerCaller, error) {
	contract, err := bindETHRegistrarController(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ETHRegistrarControllerCaller{contract: contract}, nil
}

// NewETHRegistrarControllerTransactor creates a new write-only instance of ETHRegistrarController, bound to a specific deployed contract.
func NewETHRegistrarControllerTransactor(address common.Address, transactor bind.ContractTransactor) (*ETHRegistrarControllerTransactor, error) {
	contract, err := bindETHRegistrarController(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ETHRegistrarControllerTransactor{contract: contract}, nil
}

// NewETHRegistrarControllerFilterer creates a new log filterer instance of ETHRegistrarController, bound to a specific deployed contract.
func NewETHRegistrarControllerFilterer(address common.Address, filterer bind.ContractFilterer) (*ETHRegistrarControllerFilterer, error) {
	contract, err := bindETHRegistrarController(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ETHRegistrarControllerFilterer{contract: contract}, nil
}

// bindETHRegistrarController binds a generic wrapper to an already deployed contract.
func bindETHRegistrarController(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ETHRegistrarControllerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ETHRegistrarController *ETHRegistrarControllerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ETHRegistrarController.Contract.ETHRegistrarControllerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ETHRegistrarController *ETHRegistrarControllerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ETHRegistrarController.Contract.ETHRegistrarControllerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ETHRegistrarController *ETHRegistrarControllerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ETHRegistrarController.Contract.ETHRegistrarControllerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ETHRegistrarController *ETHRegistrarControllerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ETHRegistrarController.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ETHRegistrarController *ETHRegistrarControllerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ETHRegistrarController.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ETHRegistrarController *ETHRegistrarControllerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ETHRegistrarController.Contract.contract.Transact(opts, method, params...)
}

// MINREGISTRATIONDURATION is a free data retrieval call binding the contract method 0x8a95b09f.
//
// Solidity: function MIN_REGISTRATION_DURATION() view returns(uint256)
func (_ETHRegistrarController *ETHRegistrarControllerCaller) MINREGISTRATIONDURATION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ETHRegistrarController.contract.Call(opts, &out, "MIN_REGISTRATION_DURATION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINREGISTRATIONDURATION is a free data retrieval call binding the contract method 0x8a95b09f.
//
// Solidity: function MIN_REGISTRATION_DURATION() view returns(uint256)
func (_ETHRegistrarController *ETHRegistrarControllerSession) MINREGISTRATIONDURATION() (*big.Int, error) {
	return _ETHRegistrarController.Contract.MINREGISTRATIONDURATION(&_ETHRegistrarController.CallOpts)
}

// MINREGISTRATIONDURATION is a free data retrieval call binding the contract method 0x8a95b09f.
//
// Solidity: function MIN_REGISTRATION_DURATION() view returns(uint256)
func (_ETHRegistrarController *ETHRegistrarControllerCallerSession) MINREGISTRATIONDURATION() (*big.Int, error) {
	return _ETHRegistrarController.Contract.MINREGISTRATIONDURATION(&_ETHRegistrarController.CallOpts)
}

// Available is a free data retrieval call binding the contract method 0xaeb8ce9b.
//
// Solidity: function available(string label) view returns(bool)
func (_ETHRegistrarController *ETHRegistrarControllerCaller) Available(opts *bind.CallOpts, label string) (bool, error) {
	var out []interface{}
	err := _ETHRegistrarController.contract.Call(opts, &out, "available", label)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Available is a free data retrieval call binding the contract method 0xaeb8ce9b.
//
// Solidity: function available(string label) view returns(bool)
func (_ETHRegistrarController *ETHRegistrarControllerSession) Available(label string) (bool, error) {
	return _ETHRegistrarController.Contract.Available(&_ETHRegistrarController.CallOpts, label)
}

// Available is a free data retrieval call binding the contract method 0xaeb8ce9b.
//
// Solidity: function available(string label) view returns(bool)
func (_ETHRegistrarController *ETHRegistrarControllerCallerSession) Available(label string) (bool, error) {
	return _ETHRegistrarController.Contract.Available(&_ETHRegistrarController.CallOpts, label)
}

// Commitments is a free data retrieval call binding the contract method 0x839df945.
//
// Solidity: function commitments(bytes32 ) view returns(uint256)
func (_ETHRegistrarController *ETHRegistrarControllerCaller) Commitments(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _ETHRegistrarController.contract.Call(opts, &out, "commitments", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Commitments is a free data retrieval call binding the contract method 0x839df945.
//
// Solidity: function commitments(bytes32 ) view returns(uint256)
func (_ETHRegistrarController *ETHRegistrarControllerSession) Commitments(arg0 [32]byte) (*big.Int, error) {
	return _ETHRegistrarController.Contract.Commitments(&_ETHRegistrarController.CallOpts, arg0)
}

// Commitments is a free data retrieval call binding the contract method 0x839df945.
//
// Solidity: function commitments(bytes32 ) view returns(uint256)
func (_ETHRegistrarController *ETHRegistrarControllerCallerSession) Commitments(arg0 [32]byte) (*big.Int, error) {
	return _ETHRegistrarController.Contract.Commitments(&_ETHRegistrarController.CallOpts, arg0)
}

// DefaultReverseRegistrar is a free data retrieval call binding the contract method 0x469bf441.
//
// Solidity: function defaultReverseRegistrar() view returns(address)
func (_ETHRegistrarController *ETHRegistrarControllerCaller) DefaultReverseRegistrar(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ETHRegistrarController.contract.Call(opts, &out, "defaultReverseRegistrar")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DefaultReverseRegistrar is a free data retrieval call binding the contract method 0x469bf441.
//
// Solidity: function defaultReverseRegistrar() view returns(address)
func (_ETHRegistrarController *ETHRegistrarControllerSession) DefaultReverseRegistrar() (common.Address, error) {
	return _ETHRegistrarController.Contract.DefaultReverseRegistrar(&_ETHRegistrarController.CallOpts)
}

// DefaultReverseRegistrar is a free data retrieval call binding the contract method 0x469bf441.
//
// Solidity: function defaultReverseRegistrar() view returns(address)
func (_ETHRegistrarController *ETHRegistrarControllerCallerSession) DefaultReverseRegistrar() (common.Address, error) {
	return _ETHRegistrarController.Contract.DefaultReverseRegistrar(&_ETHRegistrarController.CallOpts)
}

// Ens is a free data retrieval call binding the contract method 0x3f15457f.
//
// Solidity: function ens() view returns(address)
func (_ETHRegistrarController *ETHRegistrarControllerCaller) Ens(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ETHRegistrarController.contract.Call(opts, &out, "ens")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Ens is a free data retrieval call binding the contract method 0x3f15457f.
//
// Solidity: function ens() view returns(address)
func (_ETHRegistrarController *ETHRegistrarControllerSession) Ens() (common.Address, error) {
	return _ETHRegistrarController.Contract.Ens(&_ETHRegistrarController.CallOpts)
}

// Ens is a free data retrieval call binding the contract method 0x3f15457f.
//
// Solidity: function ens() view returns(address)
func (_ETHRegistrarController *ETHRegistrarControllerCallerSession) Ens() (common.Address, error) {
	return _ETHRegistrarController.Contract.Ens(&_ETHRegistrarController.CallOpts)
}

// MakeCommitment is a free data retrieval call binding the contract method 0xcf7d6e01.
//
// Solidity: function makeCommitment((string,address,uint256,bytes32,address,bytes[],uint8,bytes32) registration) pure returns(bytes32 commitment)
func (_ETHRegistrarController *ETHRegistrarControllerCaller) MakeCommitment(opts *bind.CallOpts, registration IETHRegistrarControllerRegistration) ([32]byte, error) {
	var out []interface{}
	err := _ETHRegistrarController.contract.Call(opts, &out, "makeCommitment", registration)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MakeCommitment is a free data retrieval call binding the contract method 0xcf7d6e01.
//
// Solidity: function makeCommitment((string,address,uint256,bytes32,address,bytes[],uint8,bytes32) registration) pure returns(bytes32 commitment)
func (_ETHRegistrarController *ETHRegistrarControllerSession) MakeCommitment(registration IETHRegistrarControllerRegistration) ([32]byte, error) {
	return _ETHRegistrarController.Contract.MakeCommitment(&_ETHRegistrarController.CallOpts, registration)
}

// MakeCommitment is a free data retrieval call binding the contract method 0xcf7d6e01.
//
// Solidity: function makeCommitment((string,address,uint256,bytes32,address,bytes[],uint8,bytes32) registration) pure returns(bytes32 commitment)
func (_ETHRegistrarController *ETHRegistrarControllerCallerSession) MakeCommitment(registration IETHRegistrarControllerRegistration) ([32]byte, error) {
	return _ETHRegistrarController.Contract.MakeCommitment(&_ETHRegistrarController.CallOpts, registration)
}

// MaxCommitmentAge is a free data retrieval call binding the contract method 0xce1e09c0.
//
// Solidity: function maxCommitmentAge() view returns(uint256)
func (_ETHRegistrarController *ETHRegistrarControllerCaller) MaxCommitmentAge(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ETHRegistrarController.contract.Call(opts, &out, "maxCommitmentAge")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxCommitmentAge is a free data retrieval call binding the contract method 0xce1e09c0.
//
// Solidity: function maxCommitmentAge() view returns(uint256)
func (_ETHRegistrarController *ETHRegistrarControllerSession) MaxCommitmentAge() (*big.Int, error) {
	return _ETHRegistrarController.Contract.MaxCommitmentAge(&_ETHRegistrarController.CallOpts)
}

// MaxCommitmentAge is a free data retrieval call binding the contract method 0xce1e09c0.
//
// Solidity: function maxCommitmentAge() view returns(uint256)
func (_ETHRegistrarController *ETHRegistrarControllerCallerSession) MaxCommitmentAge() (*big.Int, error) {
	return _ETHRegistrarController.Contract.MaxCommitmentAge(&_ETHRegistrarController.CallOpts)
}

// MinCommitmentAge is a free data retrieval call binding the contract method 0x8d839ffe.
//
// Solidity: function minCommitmentAge() view returns(uint256)
func (_ETHRegistrarController *ETHRegistrarControllerCaller) MinCommitmentAge(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ETHRegistrarController.contract.Call(opts, &out, "minCommitmentAge")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinCommitmentAge is a free data retrieval call binding the contract method 0x8d839ffe.
//
// Solidity: function minCommitmentAge() view returns(uint256)
func (_ETHRegistrarController *ETHRegistrarControllerSession) MinCommitmentAge() (*big.Int, error) {
	return _ETHRegistrarController.Contract.MinCommitmentAge(&_ETHRegistrarController.CallOpts)
}

// MinCommitmentAge is a free data retrieval call binding the contract method 0x8d839ffe.
//
// Solidity: function minCommitmentAge() view returns(uint256)
func (_ETHRegistrarController *ETHRegistrarControllerCallerSession) MinCommitmentAge() (*big.Int, error) {
	return _ETHRegistrarController.Contract.MinCommitmentAge(&_ETHRegistrarController.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ETHRegistrarController *ETHRegistrarControllerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ETHRegistrarController.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ETHRegistrarController *ETHRegistrarControllerSession) Owner() (common.Address, error) {
	return _ETHRegistrarController.Contract.Owner(&_ETHRegistrarController.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ETHRegistrarController *ETHRegistrarControllerCallerSession) Owner() (common.Address, error) {
	return _ETHRegistrarController.Contract.Owner(&_ETHRegistrarController.CallOpts)
}

// Prices is a free data retrieval call binding the contract method 0xd3419bf3.
//
// Solidity: function prices() view returns(address)
func (_ETHRegistrarController *ETHRegistrarControllerCaller) Prices(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ETHRegistrarController.contract.Call(opts, &out, "prices")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Prices is a free data retrieval call binding the contract method 0xd3419bf3.
//
// Solidity: function prices() view returns(address)
func (_ETHRegistrarController *ETHRegistrarControllerSession) Prices() (common.Address, error) {
	return _ETHRegistrarController.Contract.Prices(&_ETHRegistrarController.CallOpts)
}

// Prices is a free data retrieval call binding the contract method 0xd3419bf3.
//
// Solidity: function prices() view returns(address)
func (_ETHRegistrarController *ETHRegistrarControllerCallerSession) Prices() (common.Address, error) {
	return _ETHRegistrarController.Contract.Prices(&_ETHRegistrarController.CallOpts)
}

// RentPrice is a free data retrieval call binding the contract method 0x83e7f6ff.
//
// Solidity: function rentPrice(string label, uint256 duration) view returns((uint256,uint256) price)
func (_ETHRegistrarController *ETHRegistrarControllerCaller) RentPrice(opts *bind.CallOpts, label string, duration *big.Int) (IPriceOraclePrice, error) {
	var out []interface{}
	err := _ETHRegistrarController.contract.Call(opts, &out, "rentPrice", label, duration)

	if err != nil {
		return *new(IPriceOraclePrice), err
	}

	out0 := *abi.ConvertType(out[0], new(IPriceOraclePrice)).(*IPriceOraclePrice)

	return out0, err

}

// RentPrice is a free data retrieval call binding the contract method 0x83e7f6ff.
//
// Solidity: function rentPrice(string label, uint256 duration) view returns((uint256,uint256) price)
func (_ETHRegistrarController *ETHRegistrarControllerSession) RentPrice(label string, duration *big.Int) (IPriceOraclePrice, error) {
	return _ETHRegistrarController.Contract.RentPrice(&_ETHRegistrarController.CallOpts, label, duration)
}

// RentPrice is a free data retrieval call binding the contract method 0x83e7f6ff.
//
// Solidity: function rentPrice(string label, uint256 duration) view returns((uint256,uint256) price)
func (_ETHRegistrarController *ETHRegistrarControllerCallerSession) RentPrice(label string, duration *big.Int) (IPriceOraclePrice, error) {
	return _ETHRegistrarController.Contract.RentPrice(&_ETHRegistrarController.CallOpts, label, duration)
}

// ReverseRegistrar is a free data retrieval call binding the contract method 0x80869853.
//
// Solidity: function reverseRegistrar() view returns(address)
func (_ETHRegistrarController *ETHRegistrarControllerCaller) ReverseRegistrar(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ETHRegistrarController.contract.Call(opts, &out, "reverseRegistrar")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ReverseRegistrar is a free data retrieval call binding the contract method 0x80869853.
//
// Solidity: function reverseRegistrar() view returns(address)
func (_ETHRegistrarController *ETHRegistrarControllerSession) ReverseRegistrar() (common.Address, error) {
	return _ETHRegistrarController.Contract.ReverseRegistrar(&_ETHRegistrarController.CallOpts)
}

// ReverseRegistrar is a free data retrieval call binding the contract method 0x80869853.
//
// Solidity: function reverseRegistrar() view returns(address)
func (_ETHRegistrarController *ETHRegistrarControllerCallerSession) ReverseRegistrar() (common.Address, error) {
	return _ETHRegistrarController.Contract.ReverseRegistrar(&_ETHRegistrarController.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceID) view returns(bool)
func (_ETHRegistrarController *ETHRegistrarControllerCaller) SupportsInterface(opts *bind.CallOpts, interfaceID [4]byte) (bool, error) {
	var out []interface{}
	err := _ETHRegistrarController.contract.Call(opts, &out, "supportsInterface", interfaceID)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceID) view returns(bool)
func (_ETHRegistrarController *ETHRegistrarControllerSession) SupportsInterface(interfaceID [4]byte) (bool, error) {
	return _ETHRegistrarController.Contract.SupportsInterface(&_ETHRegistrarController.CallOpts, interfaceID)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceID) view returns(bool)
func (_ETHRegistrarController *ETHRegistrarControllerCallerSession) SupportsInterface(interfaceID [4]byte) (bool, error) {
	return _ETHRegistrarController.Contract.SupportsInterface(&_ETHRegistrarController.CallOpts, interfaceID)
}

// Valid is a free data retrieval call binding the contract method 0x9791c097.
//
// Solidity: function valid(string label) pure returns(bool)
func (_ETHRegistrarController *ETHRegistrarControllerCaller) Valid(opts *bind.CallOpts, label string) (bool, error) {
	var out []interface{}
	err := _ETHRegistrarController.contract.Call(opts, &out, "valid", label)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Valid is a free data retrieval call binding the contract method 0x9791c097.
//
// Solidity: function valid(string label) pure returns(bool)
func (_ETHRegistrarController *ETHRegistrarControllerSession) Valid(label string) (bool, error) {
	return _ETHRegistrarController.Contract.Valid(&_ETHRegistrarController.CallOpts, label)
}

// Valid is a free data retrieval call binding the contract method 0x9791c097.
//
// Solidity: function valid(string label) pure returns(bool)
func (_ETHRegistrarController *ETHRegistrarControllerCallerSession) Valid(label string) (bool, error) {
	return _ETHRegistrarController.Contract.Valid(&_ETHRegistrarController.CallOpts, label)
}

// Commit is a paid mutator transaction binding the contract method 0xf14fcbc8.
//
// Solidity: function commit(bytes32 commitment) returns()
func (_ETHRegistrarController *ETHRegistrarControllerTransactor) Commit(opts *bind.TransactOpts, commitment [32]byte) (*types.Transaction, error) {
	return _ETHRegistrarController.contract.Transact(opts, "commit", commitment)
}

// Commit is a paid mutator transaction binding the contract method 0xf14fcbc8.
//
// Solidity: function commit(bytes32 commitment) returns()
func (_ETHRegistrarController *ETHRegistrarControllerSession) Commit(commitment [32]byte) (*types.Transaction, error) {
	return _ETHRegistrarController.Contract.Commit(&_ETHRegistrarController.TransactOpts, commitment)
}

// Commit is a paid mutator transaction binding the contract method 0xf14fcbc8.
//
// Solidity: function commit(bytes32 commitment) returns()
func (_ETHRegistrarController *ETHRegistrarControllerTransactorSession) Commit(commitment [32]byte) (*types.Transaction, error) {
	return _ETHRegistrarController.Contract.Commit(&_ETHRegistrarController.TransactOpts, commitment)
}

// RecoverFunds is a paid mutator transaction binding the contract method 0x5d3590d5.
//
// Solidity: function recoverFunds(address _token, address _to, uint256 _amount) returns()
func (_ETHRegistrarController *ETHRegistrarControllerTransactor) RecoverFunds(opts *bind.TransactOpts, _token common.Address, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ETHRegistrarController.contract.Transact(opts, "recoverFunds", _token, _to, _amount)
}

// RecoverFunds is a paid mutator transaction binding the contract method 0x5d3590d5.
//
// Solidity: function recoverFunds(address _token, address _to, uint256 _amount) returns()
func (_ETHRegistrarController *ETHRegistrarControllerSession) RecoverFunds(_token common.Address, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ETHRegistrarController.Contract.RecoverFunds(&_ETHRegistrarController.TransactOpts, _token, _to, _amount)
}

// RecoverFunds is a paid mutator transaction binding the contract method 0x5d3590d5.
//
// Solidity: function recoverFunds(address _token, address _to, uint256 _amount) returns()
func (_ETHRegistrarController *ETHRegistrarControllerTransactorSession) RecoverFunds(_token common.Address, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ETHRegistrarController.Contract.RecoverFunds(&_ETHRegistrarController.TransactOpts, _token, _to, _amount)
}

// Register is a paid mutator transaction binding the contract method 0xef9c8805.
//
// Solidity: function register((string,address,uint256,bytes32,address,bytes[],uint8,bytes32) registration) payable returns()
func (_ETHRegistrarController *ETHRegistrarControllerTransactor) Register(opts *bind.TransactOpts, registration IETHRegistrarControllerRegistration) (*types.Transaction, error) {
	return _ETHRegistrarController.contract.Transact(opts, "register", registration)
}

// Register is a paid mutator transaction binding the contract method 0xef9c8805.
//
// Solidity: function register((string,address,uint256,bytes32,address,bytes[],uint8,bytes32) registration) payable returns()
func (_ETHRegistrarController *ETHRegistrarControllerSession) Register(registration IETHRegistrarControllerRegistration) (*types.Transaction, error) {
	return _ETHRegistrarController.Contract.Register(&_ETHRegistrarController.TransactOpts, registration)
}

// Register is a paid mutator transaction binding the contract method 0xef9c8805.
//
// Solidity: function register((string,address,uint256,bytes32,address,bytes[],uint8,bytes32) registration) payable returns()
func (_ETHRegistrarController *ETHRegistrarControllerTransactorSession) Register(registration IETHRegistrarControllerRegistration) (*types.Transaction, error) {
	return _ETHRegistrarController.Contract.Register(&_ETHRegistrarController.TransactOpts, registration)
}

// Renew is a paid mutator transaction binding the contract method 0x18026ad1.
//
// Solidity: function renew(string label, uint256 duration, bytes32 referrer) payable returns()
func (_ETHRegistrarController *ETHRegistrarControllerTransactor) Renew(opts *bind.TransactOpts, label string, duration *big.Int, referrer [32]byte) (*types.Transaction, error) {
	return _ETHRegistrarController.contract.Transact(opts, "renew", label, duration, referrer)
}

// Renew is a paid mutator transaction binding the contract method 0x18026ad1.
//
// Solidity: function renew(string label, uint256 duration, bytes32 referrer) payable returns()
func (_ETHRegistrarController *ETHRegistrarControllerSession) Renew(label string, duration *big.Int, referrer [32]byte) (*types.Transaction, error) {
	return _ETHRegistrarController.Contract.Renew(&_ETHRegistrarController.TransactOpts, label, duration, referrer)
}

// Renew is a paid mutator transaction binding the contract method 0x18026ad1.
//
// Solidity: function renew(string label, uint256 duration, bytes32 referrer) payable returns()
func (_ETHRegistrarController *ETHRegistrarControllerTransactorSession) Renew(label string, duration *big.Int, referrer [32]byte) (*types.Transaction, error) {
	return _ETHRegistrarController.Contract.Renew(&_ETHRegistrarController.TransactOpts, label, duration, referrer)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ETHRegistrarController *ETHRegistrarControllerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ETHRegistrarController.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ETHRegistrarController *ETHRegistrarControllerSession) RenounceOwnership() (*types.Transaction, error) {
	return _ETHRegistrarController.Contract.RenounceOwnership(&_ETHRegistrarController.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ETHRegistrarController *ETHRegistrarControllerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ETHRegistrarController.Contract.RenounceOwnership(&_ETHRegistrarController.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ETHRegistrarController *ETHRegistrarControllerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ETHRegistrarController.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ETHRegistrarController *ETHRegistrarControllerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ETHRegistrarController.Contract.TransferOwnership(&_ETHRegistrarController.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ETHRegistrarController *ETHRegistrarControllerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ETHRegistrarController.Contract.TransferOwnership(&_ETHRegistrarController.TransactOpts, newOwner)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_ETHRegistrarController *ETHRegistrarControllerTransactor) Withdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ETHRegistrarController.contract.Transact(opts, "withdraw")
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_ETHRegistrarController *ETHRegistrarControllerSession) Withdraw() (*types.Transaction, error) {
	return _ETHRegistrarController.Contract.Withdraw(&_ETHRegistrarController.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_ETHRegistrarController *ETHRegistrarControllerTransactorSession) Withdraw() (*types.Transaction, error) {
	return _ETHRegistrarController.Contract.Withdraw(&_ETHRegistrarController.TransactOpts)
}

// ETHRegistrarControllerNameRegisteredIterator is returned from FilterNameRegistered and is used to iterate over the raw logs and unpacked data for NameRegistered events raised by the ETHRegistrarController contract.
type ETHRegistrarControllerNameRegisteredIterator struct {
	Event *ETHRegistrarControllerNameRegistered // Event containing the contract specifics and raw log

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
func (it *ETHRegistrarControllerNameRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ETHRegistrarControllerNameRegistered)
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
		it.Event = new(ETHRegistrarControllerNameRegistered)
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
func (it *ETHRegistrarControllerNameRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ETHRegistrarControllerNameRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ETHRegistrarControllerNameRegistered represents a NameRegistered event raised by the ETHRegistrarController contract.
type ETHRegistrarControllerNameRegistered struct {
	Label     string
	Labelhash [32]byte
	Owner     common.Address
	BaseCost  *big.Int
	Premium   *big.Int
	Expires   *big.Int
	Referrer  [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterNameRegistered is a free log retrieval operation binding the contract event 0xc2240194853531f1ae318dcef227de79c6ad0fd9d1b0e4fe08568415be2e08a5.
//
// Solidity: event NameRegistered(string label, bytes32 indexed labelhash, address indexed owner, uint256 baseCost, uint256 premium, uint256 expires, bytes32 referrer)
func (_ETHRegistrarController *ETHRegistrarControllerFilterer) FilterNameRegistered(opts *bind.FilterOpts, labelhash [][32]byte, owner []common.Address) (*ETHRegistrarControllerNameRegisteredIterator, error) {

	var labelhashRule []interface{}
	for _, labelhashItem := range labelhash {
		labelhashRule = append(labelhashRule, labelhashItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _ETHRegistrarController.contract.FilterLogs(opts, "NameRegistered", labelhashRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return &ETHRegistrarControllerNameRegisteredIterator{contract: _ETHRegistrarController.contract, event: "NameRegistered", logs: logs, sub: sub}, nil
}

// WatchNameRegistered is a free log subscription operation binding the contract event 0xc2240194853531f1ae318dcef227de79c6ad0fd9d1b0e4fe08568415be2e08a5.
//
// Solidity: event NameRegistered(string label, bytes32 indexed labelhash, address indexed owner, uint256 baseCost, uint256 premium, uint256 expires, bytes32 referrer)
func (_ETHRegistrarController *ETHRegistrarControllerFilterer) WatchNameRegistered(opts *bind.WatchOpts, sink chan<- *ETHRegistrarControllerNameRegistered, labelhash [][32]byte, owner []common.Address) (event.Subscription, error) {

	var labelhashRule []interface{}
	for _, labelhashItem := range labelhash {
		labelhashRule = append(labelhashRule, labelhashItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _ETHRegistrarController.contract.WatchLogs(opts, "NameRegistered", labelhashRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ETHRegistrarControllerNameRegistered)
				if err := _ETHRegistrarController.contract.UnpackLog(event, "NameRegistered", log); err != nil {
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

// ParseNameRegistered is a log parse operation binding the contract event 0xc2240194853531f1ae318dcef227de79c6ad0fd9d1b0e4fe08568415be2e08a5.
//
// Solidity: event NameRegistered(string label, bytes32 indexed labelhash, address indexed owner, uint256 baseCost, uint256 premium, uint256 expires, bytes32 referrer)
func (_ETHRegistrarController *ETHRegistrarControllerFilterer) ParseNameRegistered(log types.Log) (*ETHRegistrarControllerNameRegistered, error) {
	event := new(ETHRegistrarControllerNameRegistered)
	if err := _ETHRegistrarController.contract.UnpackLog(event, "NameRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ETHRegistrarControllerNameRenewedIterator is returned from FilterNameRenewed and is used to iterate over the raw logs and unpacked data for NameRenewed events raised by the ETHRegistrarController contract.
type ETHRegistrarControllerNameRenewedIterator struct {
	Event *ETHRegistrarControllerNameRenewed // Event containing the contract specifics and raw log

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
func (it *ETHRegistrarControllerNameRenewedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ETHRegistrarControllerNameRenewed)
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
		it.Event = new(ETHRegistrarControllerNameRenewed)
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
func (it *ETHRegistrarControllerNameRenewedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ETHRegistrarControllerNameRenewedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ETHRegistrarControllerNameRenewed represents a NameRenewed event raised by the ETHRegistrarController contract.
type ETHRegistrarControllerNameRenewed struct {
	Label     string
	Labelhash [32]byte
	Cost      *big.Int
	Expires   *big.Int
	Referrer  [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterNameRenewed is a free log retrieval operation binding the contract event 0xfa956c3bce4cb4b01166868ecaf0620566bc7e33fc70b0b9c6aef61e37e50b94.
//
// Solidity: event NameRenewed(string label, bytes32 indexed labelhash, uint256 cost, uint256 expires, bytes32 referrer)
func (_ETHRegistrarController *ETHRegistrarControllerFilterer) FilterNameRenewed(opts *bind.FilterOpts, labelhash [][32]byte) (*ETHRegistrarControllerNameRenewedIterator, error) {

	var labelhashRule []interface{}
	for _, labelhashItem := range labelhash {
		labelhashRule = append(labelhashRule, labelhashItem)
	}

	logs, sub, err := _ETHRegistrarController.contract.FilterLogs(opts, "NameRenewed", labelhashRule)
	if err != nil {
		return nil, err
	}
	return &ETHRegistrarControllerNameRenewedIterator{contract: _ETHRegistrarController.contract, event: "NameRenewed", logs: logs, sub: sub}, nil
}

// WatchNameRenewed is a free log subscription operation binding the contract event 0xfa956c3bce4cb4b01166868ecaf0620566bc7e33fc70b0b9c6aef61e37e50b94.
//
// Solidity: event NameRenewed(string label, bytes32 indexed labelhash, uint256 cost, uint256 expires, bytes32 referrer)
func (_ETHRegistrarController *ETHRegistrarControllerFilterer) WatchNameRenewed(opts *bind.WatchOpts, sink chan<- *ETHRegistrarControllerNameRenewed, labelhash [][32]byte) (event.Subscription, error) {

	var labelhashRule []interface{}
	for _, labelhashItem := range labelhash {
		labelhashRule = append(labelhashRule, labelhashItem)
	}

	logs, sub, err := _ETHRegistrarController.contract.WatchLogs(opts, "NameRenewed", labelhashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ETHRegistrarControllerNameRenewed)
				if err := _ETHRegistrarController.contract.UnpackLog(event, "NameRenewed", log); err != nil {
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

// ParseNameRenewed is a log parse operation binding the contract event 0xfa956c3bce4cb4b01166868ecaf0620566bc7e33fc70b0b9c6aef61e37e50b94.
//
// Solidity: event NameRenewed(string label, bytes32 indexed labelhash, uint256 cost, uint256 expires, bytes32 referrer)
func (_ETHRegistrarController *ETHRegistrarControllerFilterer) ParseNameRenewed(log types.Log) (*ETHRegistrarControllerNameRenewed, error) {
	event := new(ETHRegistrarControllerNameRenewed)
	if err := _ETHRegistrarController.contract.UnpackLog(event, "NameRenewed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ETHRegistrarControllerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ETHRegistrarController contract.
type ETHRegistrarControllerOwnershipTransferredIterator struct {
	Event *ETHRegistrarControllerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ETHRegistrarControllerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ETHRegistrarControllerOwnershipTransferred)
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
		it.Event = new(ETHRegistrarControllerOwnershipTransferred)
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
func (it *ETHRegistrarControllerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ETHRegistrarControllerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ETHRegistrarControllerOwnershipTransferred represents a OwnershipTransferred event raised by the ETHRegistrarController contract.
type ETHRegistrarControllerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ETHRegistrarController *ETHRegistrarControllerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ETHRegistrarControllerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ETHRegistrarController.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ETHRegistrarControllerOwnershipTransferredIterator{contract: _ETHRegistrarController.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ETHRegistrarController *ETHRegistrarControllerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ETHRegistrarControllerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ETHRegistrarController.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ETHRegistrarControllerOwnershipTransferred)
				if err := _ETHRegistrarController.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ETHRegistrarController *ETHRegistrarControllerFilterer) ParseOwnershipTransferred(log types.Log) (*ETHRegistrarControllerOwnershipTransferred, error) {
	event := new(ETHRegistrarControllerOwnershipTransferred)
	if err := _ETHRegistrarController.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
