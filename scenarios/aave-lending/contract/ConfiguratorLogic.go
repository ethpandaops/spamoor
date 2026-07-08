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

// ConfiguratorLogicMetaData contains all meta data concerning the ConfiguratorLogic contract.
var ConfiguratorLogicMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proxy\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"ATokenUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"aToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"stableDebtToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"variableDebtToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"interestRateStrategyAddress\",\"type\":\"address\"}],\"name\":\"ReserveInitialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proxy\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"StableDebtTokenUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proxy\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"VariableDebtTokenUpgraded\",\"type\":\"event\"}]",
	Bin: "0x61221c61003a600b82828239805160001a60731461002d57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600436106100565760003560e01c8063b0f093551461005b578063b13c96a81461007d578063df59b8b21461009d578063f5b50e70146100bd575b600080fd5b81801561006757600080fd5b5061007b61007636600461117d565b6100dd565b005b81801561008957600080fd5b5061007b6100983660046111d4565b610439565b8180156100a957600080fd5b5061007b6100b8366004611220565b6106c6565b8180156100c957600080fd5b5061007b6100d836600461117d565b610bd3565b600073ffffffffffffffffffffffffffffffffffffffff83166335ea6a75610108602085018561126d565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b16815273ffffffffffffffffffffffffffffffffffffffff90911660048201526024016101e060405180830381865afa158015610172573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061019691906113a2565b9050600061028573ffffffffffffffffffffffffffffffffffffffff851663c44b11f76101c6602087018761126d565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b16815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401602060405180830381865afa15801561022f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061025391906114c5565b5161ffff80821692601083901c821692602081901c83169260ff603083901c811693604084901c9092169260a81c1690565b50909450600093507fc222ec8a0000000000000000000000000000000000000000000000000000000092508791506102c29050602087018761126d565b6102d2604088016020890161126d565b856102e060408a018a6114e1565b6102ed60608c018c6114e1565b6102fa60a08e018e6114e1565b6040516024016103139a99989796959493929190611596565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909316929092179091526101408401519091506103b3906103ad60a087016080880161126d565b83610e6a565b6103c360a085016080860161126d565b61014084015173ffffffffffffffffffffffffffffffffffffffff91821691166103f0602087018761126d565b73ffffffffffffffffffffffffffffffffffffffff167f9439658a562a5c46b1173589df89cf001483d685bad28aedaff4a88656292d8160405160405180910390a45050505050565b600073ffffffffffffffffffffffffffffffffffffffff83166335ea6a75610464602085018561126d565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b16815273ffffffffffffffffffffffffffffffffffffffff90911660048201526024016101e060405180830381865afa1580156104ce573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104f291906113a2565b9050600061052273ffffffffffffffffffffffffffffffffffffffff851663c44b11f76101c6602087018761126d565b50509350505050600063183fb41360e01b85856020016020810190610547919061126d565b610554602088018861126d565b6105646060890160408a0161126d565b8661057260608b018b6114e1565b61057f60808d018d6114e1565b61058c60c08f018f6114e1565b6040516024016105a69b9a99989796959493929190611617565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090931692909217909152610100840151909150610640906103ad60c0870160a0880161126d565b61065060c0850160a0860161126d565b61010084015173ffffffffffffffffffffffffffffffffffffffff918216911661067d602087018761126d565b73ffffffffffffffffffffffffffffffffffffffff167fa76f65411ec66a7fb6bc467432eb14767900449ae4469fa295e4441fe5e1cb7360405160405180910390a45050505050565b60006108016106d8602084018461126d565b7f183fb413000000000000000000000000000000000000000000000000000000008561070a60e0870160c0880161126d565b61071a60c0880160a0890161126d565b61072b610100890160e08a0161126d565b61073b60808a0160608b016116a4565b6107496101008b018b6114e1565b6107576101208d018d6114e1565b6107656101c08f018f6114e1565b60405160240161077f9b9a999897969594939291906116c7565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090931692909217909152610ef8565b905060006108ae610818604085016020860161126d565b7fc222ec8a000000000000000000000000000000000000000000000000000000008661084a60c0880160a0890161126d565b61085b610100890160e08a0161126d565b61086b60808a0160608b016116a4565b6108796101808b018b6114e1565b6108876101a08d018d6114e1565b6108956101c08f018f6114e1565b60405160240161077f9a9998979695949392919061171b565b905060006109456108c5606086016040870161126d565b7fc222ec8a00000000000000000000000000000000000000000000000000000000876108f760c0890160a08a0161126d565b6109086101008a0160e08b0161126d565b61091860808b0160608c016116a4565b6109266101408c018c6114e1565b6109346101608e018e6114e1565b8e806101c0019061089591906114e1565b905073ffffffffffffffffffffffffffffffffffffffff8516637a708e9261097360c0870160a0880161126d565b85858561098660a08b0160808c0161126d565b60405160e087901b7fffffffff0000000000000000000000000000000000000000000000000000000016815273ffffffffffffffffffffffffffffffffffffffff95861660048201529385166024850152918416604484015283166064830152909116608482015260a401600060405180830381600087803b158015610a0b57600080fd5b505af1158015610a1f573d6000803e3d6000fd5b50506040805160208101909152600081529150610a519050610a4760808701606088016116a4565b829060ff16610fd3565b610a5c81600161107c565b610a678160006110c1565b610a72816000611106565b73ffffffffffffffffffffffffffffffffffffffff861663f51e435b610a9e60c0880160a0890161126d565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b16815273ffffffffffffffffffffffffffffffffffffffff909116600482015283516024820152604401600060405180830381600087803b158015610b0b57600080fd5b505af1158015610b1f573d6000803e3d6000fd5b50505073ffffffffffffffffffffffffffffffffffffffff85169050610b4b60c0870160a0880161126d565b73ffffffffffffffffffffffffffffffffffffffff167f3a0ca721fc364424566385a1aa271ed508cc2c0949c2272575fb3013a163a45f8585610b9460a08b0160808c0161126d565b6040805173ffffffffffffffffffffffffffffffffffffffff9485168152928416602084015292168183015290519081900360600190a3505050505050565b600073ffffffffffffffffffffffffffffffffffffffff83166335ea6a75610bfe602085018561126d565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b16815273ffffffffffffffffffffffffffffffffffffffff90911660048201526024016101e060405180830381865afa158015610c68573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c8c91906113a2565b90506000610cbc73ffffffffffffffffffffffffffffffffffffffff851663c44b11f76101c6602087018761126d565b50909450600093507fc222ec8a000000000000000000000000000000000000000000000000000000009250879150610cf99050602087018761126d565b610d09604088016020890161126d565b85610d1760408a018a6114e1565b610d2460608c018c6114e1565b610d3160a08e018e6114e1565b604051602401610d4a9a99989796959493929190611596565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090931692909217909152610120840151909150610de4906103ad60a087016080880161126d565b610df460a085016080860161126d565b61012084015173ffffffffffffffffffffffffffffffffffffffff9182169116610e21602087018761126d565b73ffffffffffffffffffffffffffffffffffffffff167f7a943a5b6c214bf7726c069a878b1e2a8e7371981d516048b84e03743e67bc2860405160405180910390a45050505050565b6040517f4f1ef286000000000000000000000000000000000000000000000000000000008152839073ffffffffffffffffffffffffffffffffffffffff821690634f1ef28690610ec090869086906004016117d1565b600060405180830381600087803b158015610eda57600080fd5b505af1158015610eee573d6000803e3d6000fd5b5050505050505050565b60008030604051610f089061114b565b73ffffffffffffffffffffffffffffffffffffffff9091168152602001604051809103906000f080158015610f41573d6000803e3d6000fd5b506040517fd1f5789400000000000000000000000000000000000000000000000000000000815290915073ffffffffffffffffffffffffffffffffffffffff82169063d1f5789490610f9990879087906004016117d1565b600060405180830381600087803b158015610fb357600080fd5b505af1158015610fc7573d6000803e3d6000fd5b50929695505050505050565b60408051808201909152600281527f3636000000000000000000000000000000000000000000000000000000000000602082015260ff82111561104c576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016110439190611808565b60405180910390fd5b5081517fffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffff1660309190911b179052565b60388161108a57600061108d565b60015b83517ffffffffffffffffffffffffffffffffffffffffffffffffffeffffffffffffff1660ff9190911690911b1790915250565b603c816110cf5760006110d2565b60015b83517fffffffffffffffffffffffffffffffffffffffffffffffffefffffffffffffff1660ff9190911690911b1790915250565b603981611114576000611117565b60015b83517ffffffffffffffffffffffffffffffffffffffffffffffffffdffffffffffffff1660ff9190911690911b1790915250565b6109cb8061181c83390190565b73ffffffffffffffffffffffffffffffffffffffff8116811461117a57600080fd5b50565b6000806040838503121561119057600080fd5b823561119b81611158565b9150602083013567ffffffffffffffff8111156111b757600080fd5b830160c081860312156111c957600080fd5b809150509250929050565b600080604083850312156111e757600080fd5b82356111f281611158565b9150602083013567ffffffffffffffff81111561120e57600080fd5b830160e081860312156111c957600080fd5b6000806040838503121561123357600080fd5b823561123e81611158565b9150602083013567ffffffffffffffff81111561125a57600080fd5b83016101e081860312156111c957600080fd5b60006020828403121561127f57600080fd5b813561128a81611158565b9392505050565b6040516101e0810167ffffffffffffffff811182821017156112dc577f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405290565b6000602082840312156112f457600080fd5b6040516020810181811067ffffffffffffffff8211171561133e577f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040529151825250919050565b80516fffffffffffffffffffffffffffffffff8116811461136b57600080fd5b919050565b805164ffffffffff8116811461136b57600080fd5b805161ffff8116811461136b57600080fd5b805161136b81611158565b60006101e082840312156113b557600080fd5b6113bd611291565b6113c784846112e2565b81526113d56020840161134b565b60208201526113e66040840161134b565b60408201526113f76060840161134b565b60608201526114086080840161134b565b608082015261141960a0840161134b565b60a082015261142a60c08401611370565b60c082015261143b60e08401611385565b60e082015261010061144e818501611397565b90820152610120611460848201611397565b90820152610140611472848201611397565b90820152610160611484848201611397565b9082015261018061149684820161134b565b908201526101a06114a884820161134b565b908201526101c06114ba84820161134b565b908201529392505050565b6000602082840312156114d757600080fd5b61128a83836112e2565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261151657600080fd5b83018035915067ffffffffffffffff82111561153157600080fd5b60200191503681900382131561154657600080fd5b9250929050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b600073ffffffffffffffffffffffffffffffffffffffff808d168352808c166020840152808b1660408401525088606083015260e060808301526115de60e08301888a61154d565b82810360a08401526115f181878961154d565b905082810360c084015261160681858761154d565b9d9c50505050505050505050505050565b600061010073ffffffffffffffffffffffffffffffffffffffff808f168452808e166020850152808d166040850152808c166060850152508960808401528060a0840152611668818401898b61154d565b905082810360c084015261167d81878961154d565b905082810360e084015261169281858761154d565b9e9d5050505050505050505050505050565b6000602082840312156116b657600080fd5b813560ff8116811461128a57600080fd5b600061010073ffffffffffffffffffffffffffffffffffffffff808f168452808e166020850152808d166040850152808c1660608501525060ff8a1660808401528060a0840152611668818401898b61154d565b600073ffffffffffffffffffffffffffffffffffffffff808d168352808c166020840152808b1660408401525060ff8916606083015260e060808301526115de60e08301888a61154d565b6000815180845260005b8181101561178c57602081850181015186830182015201611770565b8181111561179e576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b73ffffffffffffffffffffffffffffffffffffffff831681526040602082015260006118006040830184611766565b949350505050565b60208152600061128a602083018461176656fe60a060405234801561001057600080fd5b506040516109cb3803806109cb83398101604081905261002f91610040565b6001600160a01b0316608052610070565b60006020828403121561005257600080fd5b81516001600160a01b038116811461006957600080fd5b9392505050565b60805161091d6100ae6000396000818161014f015281816101a101528181610274015281816104110152818161043a01526105a4015261091d6000f3fe60806040526004361061005a5760003560e01c80635c60da1b116100435780635c60da1b14610097578063d1f57894146100d5578063f851a440146100e85761005a565b80633659cfe6146100645780634f1ef28614610084575b6100626100fd565b005b34801561007057600080fd5b5061006261007f36600461067b565b610137565b61006261009236600461069d565b610189565b3480156100a357600080fd5b506100ac61025a565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390f35b6100626100e336600461074f565b6102cb565b3480156100f457600080fd5b506100ac6103f7565b61010561045c565b6101356101307f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5490565b610464565b565b3373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614156101815761017e81610488565b50565b61017e6100fd565b3373ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016141561024d576101d083610488565b60008373ffffffffffffffffffffffffffffffffffffffff1683836040516101f992919061082f565b600060405180830381855af49150503d8060008114610234576040519150601f19603f3d011682016040523d82523d6000602084013e610239565b606091505b505090508061024757600080fd5b50505050565b6102556100fd565b505050565b60003373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614156102c057507f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5490565b6102c86100fd565b90565b60006102f57f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5490565b73ffffffffffffffffffffffffffffffffffffffff161461031557600080fd5b61034060017f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbd61083f565b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc1461036e5761036e61087d565b610377826104d5565b8051156103f35760008273ffffffffffffffffffffffffffffffffffffffff16826040516103a591906108ac565b600060405180830381855af49150503d80600081146103e0576040519150601f19603f3d011682016040523d82523d6000602084013e6103e5565b606091505b505090508061025557600080fd5b5050565b60003373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614156102c057507f000000000000000000000000000000000000000000000000000000000000000090565b61013561058c565b3660008037600080366000845af43d6000803e808015610483573d6000f35b3d6000fd5b610491816104d5565b60405173ffffffffffffffffffffffffffffffffffffffff8216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b803b610568576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603b60248201527f43616e6e6f742073657420612070726f787920696d706c656d656e746174696f60448201527f6e20746f2061206e6f6e2d636f6e74726163742061646472657373000000000060648201526084015b60405180910390fd5b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc55565b3373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161415610135576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603260248201527f43616e6e6f742063616c6c2066616c6c6261636b2066756e6374696f6e20667260448201527f6f6d207468652070726f78792061646d696e0000000000000000000000000000606482015260840161055f565b803573ffffffffffffffffffffffffffffffffffffffff8116811461067657600080fd5b919050565b60006020828403121561068d57600080fd5b61069682610652565b9392505050565b6000806000604084860312156106b257600080fd5b6106bb84610652565b9250602084013567ffffffffffffffff808211156106d857600080fd5b818601915086601f8301126106ec57600080fd5b8135818111156106fb57600080fd5b87602082850101111561070d57600080fd5b6020830194508093505050509250925092565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6000806040838503121561076257600080fd5b61076b83610652565b9150602083013567ffffffffffffffff8082111561078857600080fd5b818501915085601f83011261079c57600080fd5b8135818111156107ae576107ae610720565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f011681019083821181831017156107f4576107f4610720565b8160405282815288602084870101111561080d57600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b8183823760009101908152919050565b600082821015610878577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b500390565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052600160045260246000fd5b6000825160005b818110156108cd57602081860181015185830152016108b3565b818111156108dc576000828501525b50919091019291505056fea2646970667358221220899ba9574e8c52c72539176723d8c74a8618334587150196e2029371e7486a8464736f6c634300080a0033a2646970667358221220722c61fe2602f4b3008e7a9a899fcb2eaf1517e82305577945d327c6cdb6b63864736f6c634300080a0033",
}

// ConfiguratorLogicABI is the input ABI used to generate the binding from.
// Deprecated: Use ConfiguratorLogicMetaData.ABI instead.
var ConfiguratorLogicABI = ConfiguratorLogicMetaData.ABI

// ConfiguratorLogicBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ConfiguratorLogicMetaData.Bin instead.
var ConfiguratorLogicBin = ConfiguratorLogicMetaData.Bin

// DeployConfiguratorLogic deploys a new Ethereum contract, binding an instance of ConfiguratorLogic to it.
func DeployConfiguratorLogic(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ConfiguratorLogic, error) {
	parsed, err := ConfiguratorLogicMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ConfiguratorLogicBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ConfiguratorLogic{ConfiguratorLogicCaller: ConfiguratorLogicCaller{contract: contract}, ConfiguratorLogicTransactor: ConfiguratorLogicTransactor{contract: contract}, ConfiguratorLogicFilterer: ConfiguratorLogicFilterer{contract: contract}}, nil
}

// ConfiguratorLogic is an auto generated Go binding around an Ethereum contract.
type ConfiguratorLogic struct {
	ConfiguratorLogicCaller     // Read-only binding to the contract
	ConfiguratorLogicTransactor // Write-only binding to the contract
	ConfiguratorLogicFilterer   // Log filterer for contract events
}

// ConfiguratorLogicCaller is an auto generated read-only Go binding around an Ethereum contract.
type ConfiguratorLogicCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConfiguratorLogicTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ConfiguratorLogicTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConfiguratorLogicFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ConfiguratorLogicFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConfiguratorLogicSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ConfiguratorLogicSession struct {
	Contract     *ConfiguratorLogic // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ConfiguratorLogicCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ConfiguratorLogicCallerSession struct {
	Contract *ConfiguratorLogicCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// ConfiguratorLogicTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ConfiguratorLogicTransactorSession struct {
	Contract     *ConfiguratorLogicTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// ConfiguratorLogicRaw is an auto generated low-level Go binding around an Ethereum contract.
type ConfiguratorLogicRaw struct {
	Contract *ConfiguratorLogic // Generic contract binding to access the raw methods on
}

// ConfiguratorLogicCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ConfiguratorLogicCallerRaw struct {
	Contract *ConfiguratorLogicCaller // Generic read-only contract binding to access the raw methods on
}

// ConfiguratorLogicTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ConfiguratorLogicTransactorRaw struct {
	Contract *ConfiguratorLogicTransactor // Generic write-only contract binding to access the raw methods on
}

// NewConfiguratorLogic creates a new instance of ConfiguratorLogic, bound to a specific deployed contract.
func NewConfiguratorLogic(address common.Address, backend bind.ContractBackend) (*ConfiguratorLogic, error) {
	contract, err := bindConfiguratorLogic(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ConfiguratorLogic{ConfiguratorLogicCaller: ConfiguratorLogicCaller{contract: contract}, ConfiguratorLogicTransactor: ConfiguratorLogicTransactor{contract: contract}, ConfiguratorLogicFilterer: ConfiguratorLogicFilterer{contract: contract}}, nil
}

// NewConfiguratorLogicCaller creates a new read-only instance of ConfiguratorLogic, bound to a specific deployed contract.
func NewConfiguratorLogicCaller(address common.Address, caller bind.ContractCaller) (*ConfiguratorLogicCaller, error) {
	contract, err := bindConfiguratorLogic(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConfiguratorLogicCaller{contract: contract}, nil
}

// NewConfiguratorLogicTransactor creates a new write-only instance of ConfiguratorLogic, bound to a specific deployed contract.
func NewConfiguratorLogicTransactor(address common.Address, transactor bind.ContractTransactor) (*ConfiguratorLogicTransactor, error) {
	contract, err := bindConfiguratorLogic(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConfiguratorLogicTransactor{contract: contract}, nil
}

// NewConfiguratorLogicFilterer creates a new log filterer instance of ConfiguratorLogic, bound to a specific deployed contract.
func NewConfiguratorLogicFilterer(address common.Address, filterer bind.ContractFilterer) (*ConfiguratorLogicFilterer, error) {
	contract, err := bindConfiguratorLogic(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConfiguratorLogicFilterer{contract: contract}, nil
}

// bindConfiguratorLogic binds a generic wrapper to an already deployed contract.
func bindConfiguratorLogic(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ConfiguratorLogicMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConfiguratorLogic *ConfiguratorLogicRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConfiguratorLogic.Contract.ConfiguratorLogicCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConfiguratorLogic *ConfiguratorLogicRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConfiguratorLogic.Contract.ConfiguratorLogicTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConfiguratorLogic *ConfiguratorLogicRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConfiguratorLogic.Contract.ConfiguratorLogicTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConfiguratorLogic *ConfiguratorLogicCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConfiguratorLogic.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConfiguratorLogic *ConfiguratorLogicTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConfiguratorLogic.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConfiguratorLogic *ConfiguratorLogicTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConfiguratorLogic.Contract.contract.Transact(opts, method, params...)
}

// ConfiguratorLogicATokenUpgradedIterator is returned from FilterATokenUpgraded and is used to iterate over the raw logs and unpacked data for ATokenUpgraded events raised by the ConfiguratorLogic contract.
type ConfiguratorLogicATokenUpgradedIterator struct {
	Event *ConfiguratorLogicATokenUpgraded // Event containing the contract specifics and raw log

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
func (it *ConfiguratorLogicATokenUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfiguratorLogicATokenUpgraded)
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
		it.Event = new(ConfiguratorLogicATokenUpgraded)
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
func (it *ConfiguratorLogicATokenUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConfiguratorLogicATokenUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConfiguratorLogicATokenUpgraded represents a ATokenUpgraded event raised by the ConfiguratorLogic contract.
type ConfiguratorLogicATokenUpgraded struct {
	Asset          common.Address
	Proxy          common.Address
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterATokenUpgraded is a free log retrieval operation binding the contract event 0xa76f65411ec66a7fb6bc467432eb14767900449ae4469fa295e4441fe5e1cb73.
//
// Solidity: event ATokenUpgraded(address indexed asset, address indexed proxy, address indexed implementation)
func (_ConfiguratorLogic *ConfiguratorLogicFilterer) FilterATokenUpgraded(opts *bind.FilterOpts, asset []common.Address, proxy []common.Address, implementation []common.Address) (*ConfiguratorLogicATokenUpgradedIterator, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}
	var proxyRule []interface{}
	for _, proxyItem := range proxy {
		proxyRule = append(proxyRule, proxyItem)
	}
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ConfiguratorLogic.contract.FilterLogs(opts, "ATokenUpgraded", assetRule, proxyRule, implementationRule)
	if err != nil {
		return nil, err
	}
	return &ConfiguratorLogicATokenUpgradedIterator{contract: _ConfiguratorLogic.contract, event: "ATokenUpgraded", logs: logs, sub: sub}, nil
}

// WatchATokenUpgraded is a free log subscription operation binding the contract event 0xa76f65411ec66a7fb6bc467432eb14767900449ae4469fa295e4441fe5e1cb73.
//
// Solidity: event ATokenUpgraded(address indexed asset, address indexed proxy, address indexed implementation)
func (_ConfiguratorLogic *ConfiguratorLogicFilterer) WatchATokenUpgraded(opts *bind.WatchOpts, sink chan<- *ConfiguratorLogicATokenUpgraded, asset []common.Address, proxy []common.Address, implementation []common.Address) (event.Subscription, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}
	var proxyRule []interface{}
	for _, proxyItem := range proxy {
		proxyRule = append(proxyRule, proxyItem)
	}
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ConfiguratorLogic.contract.WatchLogs(opts, "ATokenUpgraded", assetRule, proxyRule, implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConfiguratorLogicATokenUpgraded)
				if err := _ConfiguratorLogic.contract.UnpackLog(event, "ATokenUpgraded", log); err != nil {
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

// ParseATokenUpgraded is a log parse operation binding the contract event 0xa76f65411ec66a7fb6bc467432eb14767900449ae4469fa295e4441fe5e1cb73.
//
// Solidity: event ATokenUpgraded(address indexed asset, address indexed proxy, address indexed implementation)
func (_ConfiguratorLogic *ConfiguratorLogicFilterer) ParseATokenUpgraded(log types.Log) (*ConfiguratorLogicATokenUpgraded, error) {
	event := new(ConfiguratorLogicATokenUpgraded)
	if err := _ConfiguratorLogic.contract.UnpackLog(event, "ATokenUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConfiguratorLogicReserveInitializedIterator is returned from FilterReserveInitialized and is used to iterate over the raw logs and unpacked data for ReserveInitialized events raised by the ConfiguratorLogic contract.
type ConfiguratorLogicReserveInitializedIterator struct {
	Event *ConfiguratorLogicReserveInitialized // Event containing the contract specifics and raw log

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
func (it *ConfiguratorLogicReserveInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfiguratorLogicReserveInitialized)
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
		it.Event = new(ConfiguratorLogicReserveInitialized)
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
func (it *ConfiguratorLogicReserveInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConfiguratorLogicReserveInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConfiguratorLogicReserveInitialized represents a ReserveInitialized event raised by the ConfiguratorLogic contract.
type ConfiguratorLogicReserveInitialized struct {
	Asset                       common.Address
	AToken                      common.Address
	StableDebtToken             common.Address
	VariableDebtToken           common.Address
	InterestRateStrategyAddress common.Address
	Raw                         types.Log // Blockchain specific contextual infos
}

// FilterReserveInitialized is a free log retrieval operation binding the contract event 0x3a0ca721fc364424566385a1aa271ed508cc2c0949c2272575fb3013a163a45f.
//
// Solidity: event ReserveInitialized(address indexed asset, address indexed aToken, address stableDebtToken, address variableDebtToken, address interestRateStrategyAddress)
func (_ConfiguratorLogic *ConfiguratorLogicFilterer) FilterReserveInitialized(opts *bind.FilterOpts, asset []common.Address, aToken []common.Address) (*ConfiguratorLogicReserveInitializedIterator, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}
	var aTokenRule []interface{}
	for _, aTokenItem := range aToken {
		aTokenRule = append(aTokenRule, aTokenItem)
	}

	logs, sub, err := _ConfiguratorLogic.contract.FilterLogs(opts, "ReserveInitialized", assetRule, aTokenRule)
	if err != nil {
		return nil, err
	}
	return &ConfiguratorLogicReserveInitializedIterator{contract: _ConfiguratorLogic.contract, event: "ReserveInitialized", logs: logs, sub: sub}, nil
}

// WatchReserveInitialized is a free log subscription operation binding the contract event 0x3a0ca721fc364424566385a1aa271ed508cc2c0949c2272575fb3013a163a45f.
//
// Solidity: event ReserveInitialized(address indexed asset, address indexed aToken, address stableDebtToken, address variableDebtToken, address interestRateStrategyAddress)
func (_ConfiguratorLogic *ConfiguratorLogicFilterer) WatchReserveInitialized(opts *bind.WatchOpts, sink chan<- *ConfiguratorLogicReserveInitialized, asset []common.Address, aToken []common.Address) (event.Subscription, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}
	var aTokenRule []interface{}
	for _, aTokenItem := range aToken {
		aTokenRule = append(aTokenRule, aTokenItem)
	}

	logs, sub, err := _ConfiguratorLogic.contract.WatchLogs(opts, "ReserveInitialized", assetRule, aTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConfiguratorLogicReserveInitialized)
				if err := _ConfiguratorLogic.contract.UnpackLog(event, "ReserveInitialized", log); err != nil {
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

// ParseReserveInitialized is a log parse operation binding the contract event 0x3a0ca721fc364424566385a1aa271ed508cc2c0949c2272575fb3013a163a45f.
//
// Solidity: event ReserveInitialized(address indexed asset, address indexed aToken, address stableDebtToken, address variableDebtToken, address interestRateStrategyAddress)
func (_ConfiguratorLogic *ConfiguratorLogicFilterer) ParseReserveInitialized(log types.Log) (*ConfiguratorLogicReserveInitialized, error) {
	event := new(ConfiguratorLogicReserveInitialized)
	if err := _ConfiguratorLogic.contract.UnpackLog(event, "ReserveInitialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConfiguratorLogicStableDebtTokenUpgradedIterator is returned from FilterStableDebtTokenUpgraded and is used to iterate over the raw logs and unpacked data for StableDebtTokenUpgraded events raised by the ConfiguratorLogic contract.
type ConfiguratorLogicStableDebtTokenUpgradedIterator struct {
	Event *ConfiguratorLogicStableDebtTokenUpgraded // Event containing the contract specifics and raw log

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
func (it *ConfiguratorLogicStableDebtTokenUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfiguratorLogicStableDebtTokenUpgraded)
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
		it.Event = new(ConfiguratorLogicStableDebtTokenUpgraded)
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
func (it *ConfiguratorLogicStableDebtTokenUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConfiguratorLogicStableDebtTokenUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConfiguratorLogicStableDebtTokenUpgraded represents a StableDebtTokenUpgraded event raised by the ConfiguratorLogic contract.
type ConfiguratorLogicStableDebtTokenUpgraded struct {
	Asset          common.Address
	Proxy          common.Address
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterStableDebtTokenUpgraded is a free log retrieval operation binding the contract event 0x7a943a5b6c214bf7726c069a878b1e2a8e7371981d516048b84e03743e67bc28.
//
// Solidity: event StableDebtTokenUpgraded(address indexed asset, address indexed proxy, address indexed implementation)
func (_ConfiguratorLogic *ConfiguratorLogicFilterer) FilterStableDebtTokenUpgraded(opts *bind.FilterOpts, asset []common.Address, proxy []common.Address, implementation []common.Address) (*ConfiguratorLogicStableDebtTokenUpgradedIterator, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}
	var proxyRule []interface{}
	for _, proxyItem := range proxy {
		proxyRule = append(proxyRule, proxyItem)
	}
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ConfiguratorLogic.contract.FilterLogs(opts, "StableDebtTokenUpgraded", assetRule, proxyRule, implementationRule)
	if err != nil {
		return nil, err
	}
	return &ConfiguratorLogicStableDebtTokenUpgradedIterator{contract: _ConfiguratorLogic.contract, event: "StableDebtTokenUpgraded", logs: logs, sub: sub}, nil
}

// WatchStableDebtTokenUpgraded is a free log subscription operation binding the contract event 0x7a943a5b6c214bf7726c069a878b1e2a8e7371981d516048b84e03743e67bc28.
//
// Solidity: event StableDebtTokenUpgraded(address indexed asset, address indexed proxy, address indexed implementation)
func (_ConfiguratorLogic *ConfiguratorLogicFilterer) WatchStableDebtTokenUpgraded(opts *bind.WatchOpts, sink chan<- *ConfiguratorLogicStableDebtTokenUpgraded, asset []common.Address, proxy []common.Address, implementation []common.Address) (event.Subscription, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}
	var proxyRule []interface{}
	for _, proxyItem := range proxy {
		proxyRule = append(proxyRule, proxyItem)
	}
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ConfiguratorLogic.contract.WatchLogs(opts, "StableDebtTokenUpgraded", assetRule, proxyRule, implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConfiguratorLogicStableDebtTokenUpgraded)
				if err := _ConfiguratorLogic.contract.UnpackLog(event, "StableDebtTokenUpgraded", log); err != nil {
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

// ParseStableDebtTokenUpgraded is a log parse operation binding the contract event 0x7a943a5b6c214bf7726c069a878b1e2a8e7371981d516048b84e03743e67bc28.
//
// Solidity: event StableDebtTokenUpgraded(address indexed asset, address indexed proxy, address indexed implementation)
func (_ConfiguratorLogic *ConfiguratorLogicFilterer) ParseStableDebtTokenUpgraded(log types.Log) (*ConfiguratorLogicStableDebtTokenUpgraded, error) {
	event := new(ConfiguratorLogicStableDebtTokenUpgraded)
	if err := _ConfiguratorLogic.contract.UnpackLog(event, "StableDebtTokenUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConfiguratorLogicVariableDebtTokenUpgradedIterator is returned from FilterVariableDebtTokenUpgraded and is used to iterate over the raw logs and unpacked data for VariableDebtTokenUpgraded events raised by the ConfiguratorLogic contract.
type ConfiguratorLogicVariableDebtTokenUpgradedIterator struct {
	Event *ConfiguratorLogicVariableDebtTokenUpgraded // Event containing the contract specifics and raw log

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
func (it *ConfiguratorLogicVariableDebtTokenUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfiguratorLogicVariableDebtTokenUpgraded)
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
		it.Event = new(ConfiguratorLogicVariableDebtTokenUpgraded)
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
func (it *ConfiguratorLogicVariableDebtTokenUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConfiguratorLogicVariableDebtTokenUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConfiguratorLogicVariableDebtTokenUpgraded represents a VariableDebtTokenUpgraded event raised by the ConfiguratorLogic contract.
type ConfiguratorLogicVariableDebtTokenUpgraded struct {
	Asset          common.Address
	Proxy          common.Address
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterVariableDebtTokenUpgraded is a free log retrieval operation binding the contract event 0x9439658a562a5c46b1173589df89cf001483d685bad28aedaff4a88656292d81.
//
// Solidity: event VariableDebtTokenUpgraded(address indexed asset, address indexed proxy, address indexed implementation)
func (_ConfiguratorLogic *ConfiguratorLogicFilterer) FilterVariableDebtTokenUpgraded(opts *bind.FilterOpts, asset []common.Address, proxy []common.Address, implementation []common.Address) (*ConfiguratorLogicVariableDebtTokenUpgradedIterator, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}
	var proxyRule []interface{}
	for _, proxyItem := range proxy {
		proxyRule = append(proxyRule, proxyItem)
	}
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ConfiguratorLogic.contract.FilterLogs(opts, "VariableDebtTokenUpgraded", assetRule, proxyRule, implementationRule)
	if err != nil {
		return nil, err
	}
	return &ConfiguratorLogicVariableDebtTokenUpgradedIterator{contract: _ConfiguratorLogic.contract, event: "VariableDebtTokenUpgraded", logs: logs, sub: sub}, nil
}

// WatchVariableDebtTokenUpgraded is a free log subscription operation binding the contract event 0x9439658a562a5c46b1173589df89cf001483d685bad28aedaff4a88656292d81.
//
// Solidity: event VariableDebtTokenUpgraded(address indexed asset, address indexed proxy, address indexed implementation)
func (_ConfiguratorLogic *ConfiguratorLogicFilterer) WatchVariableDebtTokenUpgraded(opts *bind.WatchOpts, sink chan<- *ConfiguratorLogicVariableDebtTokenUpgraded, asset []common.Address, proxy []common.Address, implementation []common.Address) (event.Subscription, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}
	var proxyRule []interface{}
	for _, proxyItem := range proxy {
		proxyRule = append(proxyRule, proxyItem)
	}
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ConfiguratorLogic.contract.WatchLogs(opts, "VariableDebtTokenUpgraded", assetRule, proxyRule, implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConfiguratorLogicVariableDebtTokenUpgraded)
				if err := _ConfiguratorLogic.contract.UnpackLog(event, "VariableDebtTokenUpgraded", log); err != nil {
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

// ParseVariableDebtTokenUpgraded is a log parse operation binding the contract event 0x9439658a562a5c46b1173589df89cf001483d685bad28aedaff4a88656292d81.
//
// Solidity: event VariableDebtTokenUpgraded(address indexed asset, address indexed proxy, address indexed implementation)
func (_ConfiguratorLogic *ConfiguratorLogicFilterer) ParseVariableDebtTokenUpgraded(log types.Log) (*ConfiguratorLogicVariableDebtTokenUpgraded, error) {
	event := new(ConfiguratorLogicVariableDebtTokenUpgraded)
	if err := _ConfiguratorLogic.contract.UnpackLog(event, "VariableDebtTokenUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
