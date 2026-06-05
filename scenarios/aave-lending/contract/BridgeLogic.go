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

// BridgeLogicMetaData contains all meta data concerning the BridgeLogic contract.
var BridgeLogicMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"reserve\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"backer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"BackUnbacked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"reserve\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"onBehalfOf\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint16\",\"name\":\"referralCode\",\"type\":\"uint16\"}],\"name\":\"MintUnbacked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"reserve\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"ReserveUsedAsCollateralEnabled\",\"type\":\"event\"}]",
	Bin: "0x6122e261003a600b82828239805160001a60731461002d57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600436106100405760003560e01c80630413c86f146100455780638e74324814610067575b600080fd5b81801561005157600080fd5b50610065610060366004611e11565b610099565b005b81801561007357600080fd5b50610087610082366004611e8a565b6103f7565b60405190815260200160405180910390f35b73ffffffffffffffffffffffffffffffffffffffff84166000908152602088905260408120906100c8826106ee565b90506100d48282610907565b6100df818387610992565b6101c08101515160b081901c640fffffffff169060301c60ff16600061010488610d2a565b60088601805460109061013e90849070010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff16611f01565b92506101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff16021790556fffffffffffffffffffffffffffffffff16905081600a6101949190612055565b61019e9084612061565b8111156040518060400160405280600281526020017f353200000000000000000000000000000000000000000000000000000000000081525090610218576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161020f919061209e565b60405180910390fd5b5061022785858b600080610dd0565b6101e08401516101008501516040517fb3f1c93d00000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff8a81166024830152604482018c90526064820192909252600092919091169063b3f1c93d906084016020604051808303816000875af11580156102bb573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102df9190612111565b9050801561038c576102fe8d8d8d886101c00151896101e00151611111565b1561038c576003860154610332908c907501000000000000000000000000000000000000000000900461ffff166001611351565b8773ffffffffffffffffffffffffffffffffffffffff168a73ffffffffffffffffffffffffffffffffffffffff167e058a56ea94653cdf4f152d227ace22d4c00ad99e2a43f58cb7d9e3feb295f260405160405180910390a35b60408051338152602081018b905261ffff89169173ffffffffffffffffffffffffffffffffffffffff808c1692908e16917ff25af37b3d3ec226063dc9bdc103ece7eb110a50f340fe854bb7bc1b0676d7d0910160405180910390a450505050505050505050505050565b600080610403876106ee565b905061040f8782610907565b600887015460009070010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff16861061047357600888015470010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff16610475565b855b9050600061048386866113e8565b905060006104918288612133565b9050600061049f888561214a565b61010086015160088d0154919250610555916104cf916fffffffffffffffffffffffffffffffff9091169061142b565b866101e0015173ffffffffffffffffffffffffffffffffffffffff166318160ddd6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561051f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105439190612162565b61054d919061214a565b8c9084611482565b61010086018190526105719061056c908590611522565b610d2a565b60088c0180546000906105979084906fffffffffffffffffffffffffffffffff16611f01565b92506101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff1602179055506105d684610d2a565b60088c01805460109061061090849070010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff1661217b565b92506101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff160217905550610660858b8360008f610dd090949392919063ffffffff16565b6101e085015161068a9073ffffffffffffffffffffffffffffffffffffffff8c1690339084611561565b60408051858152602081018a9052339173ffffffffffffffffffffffffffffffffffffffff8d16917f281596e92b2d974beb7d4f124df30a0b39067b096893e95011ce4bdad798b759910160405180910390a3509193505050505b95945050505050565b6106f6611d3f565b6106fe611d3f565b60408051602081018252845481526101c0830181905251901c61ffff166101a082015260018301546fffffffffffffffffffffffffffffffff808216610100840181905260e0840152600285015480821661014085018190526101208501527001000000000000000000000000000000009283900482166101608501528290041661018083015260048085015473ffffffffffffffffffffffffffffffffffffffff9081166101e085015260058601548116610200850152600686015416610220840181905260038601549290920464ffffffffff16610240840152604080517fb1bf962d000000000000000000000000000000000000000000000000000000008152905163b1bf962d928281019260209291908290030181865afa15801561082b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061084f9190612162565b816020018181525081600001818152505080610200015173ffffffffffffffffffffffffffffffffffffffff1663797743386040518163ffffffff1660e01b8152600401608060405180830381865afa1580156108b0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108d491906121ac565b64ffffffffff166102608501526060840181905260808401829052604084019290925260c083015260a082015292915050565b60038201544264ffffffffff908116700100000000000000000000000000000000909204161415610936575050565b6109408282611643565b61094a8282611764565b5060030180547fffffffffffffffffffffff0000000000ffffffffffffffffffffffffffffffff167001000000000000000000000000000000004264ffffffffff1602179055565b60408051808201909152600281527f32360000000000000000000000000000000000000000000000000000000000006020820152816109fe576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161020f919061209e565b506000806000610a55866101c0015151670100000000000000811615159167020000000000000082161515916704000000000000008116151591670800000000000000821615159167100000000000000016151590565b9450505092509250826040518060400160405280600281526020017f323700000000000000000000000000000000000000000000000000000000000081525090610acc576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161020f919061209e565b5060408051808201909152600281527f323900000000000000000000000000000000000000000000000000000000000060208201528115610b3a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161020f919061209e565b5060408051808201909152600281527f323800000000000000000000000000000000000000000000000000000000000060208201528215610ba8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161020f919061209e565b506101c08601515160741c640fffffffff16801580610cb257506101c08701515160301c60ff16610bda90600a612055565b610be49082612061565b85610ca58961010001518960080160009054906101000a90046fffffffffffffffffffffffffffffffff166fffffffffffffffffffffffffffffffff168b6101e0015173ffffffffffffffffffffffffffffffffffffffff1663b1bf962d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610c71573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c959190612162565b610c9f919061214a565b9061142b565b610caf919061214a565b11155b6040518060400160405280600281526020017f353100000000000000000000000000000000000000000000000000000000000081525090610d20576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161020f919061209e565b5050505050505050565b60006fffffffffffffffffffffffffffffffff821115610dcc576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602760248201527f53616665436173743a2076616c756520646f65736e27742066697420696e203160448201527f3238206269747300000000000000000000000000000000000000000000000000606482015260840161020f565b5090565b610dfb6040518060800160405280600081526020016000815260200160008152602001600081525090565b6101408501516020860151610e0f9161142b565b60608083019182526007880154604080516101208101825260088b01546fffffffffffffffffffffffffffffffff7001000000000000000000000000000000009091041681526020810188905280820187905260c0808b0151948201949094529351608085015260a0808a0151908501526101a08901519284019290925273ffffffffffffffffffffffffffffffffffffffff87811660e08501526101e0890151811661010085015291517fa589870900000000000000000000000000000000000000000000000000000000815291169163a589870991610f709190600401600061012082019050825182526020830151602083015260408301516040830152606083015160608301526080830151608083015260a083015160a083015260c083015160c083015260e083015173ffffffffffffffffffffffffffffffffffffffff80821660e0850152610100915080828601511682850152505092915050565b606060405180830381865afa158015610f8d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610fb191906121f7565b60408401526020830152808252610fc790610d2a565b6001870180546fffffffffffffffffffffffffffffffff928316700100000000000000000000000000000000029216919091179055602081015161100a90610d2a565b6003870180547fffffffffffffffffffffffffffffffff00000000000000000000000000000000166fffffffffffffffffffffffffffffffff92909216919091179055604081015161105b90610d2a565b6002870180546fffffffffffffffffffffffffffffffff92831670010000000000000000000000000000000002921691909117905580516020808301516040808501516101008a01516101408b0151835196875294860193909352908401526060830152608082015273ffffffffffffffffffffffffffffffffffffffff8516907f804c9b842b2748a22bb64b345453a3de7ca54a6ca45ce00d415894979e22897a9060a00160405180910390a2505050505050565b815160009060d41c64ffffffffff161561133b5760008273ffffffffffffffffffffffffffffffffffffffff16637535d2466040518163ffffffff1660e01b8152600401602060405180830381865afa158015611172573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906111969190612225565b73ffffffffffffffffffffffffffffffffffffffff16630542975c6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156111e0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112049190612225565b90508073ffffffffffffffffffffffffffffffffffffffff1663707cd7166040518163ffffffff1660e01b8152600401602060405180830381865afa158015611251573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112759190612225565b6040517f91d148540000000000000000000000000000000000000000000000000000000081527fd1d2cf869016112a9af1107bcf43c3759daf22cf734aad47d0c9c726e33bc782600482015233602482015273ffffffffffffffffffffffffffffffffffffffff91909116906391d1485490604401602060405180830381865afa158015611307573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061132b9190612111565b6113395760009150506106e5565b505b611347868686866118e4565b9695505050505050565b60408051808201909152600281527f37340000000000000000000000000000000000000000000000000000000000006020820152608083106113c0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161020f919061209e565b50600182811b81011b81156113da578354811784556113e2565b835481191684555b50505050565b600081157fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffec778390048411151761141d57600080fd5b506127109102611388010490565b600081157ffffffffffffffffffffffffffffffffffffffffffe6268e1b017bfe18bffffff8390048411151761146057600080fd5b506b033b2e3c9fd0803ce800000091026b019d971e4fe8401e74000000010490565b600183015460009081906114ca906fffffffffffffffffffffffffffffffff166b033b2e3c9fd0803ce8000000610c956114bb88611981565b6114c488611981565b90611522565b90506114d581610d2a565b6001860180547fffffffffffffffffffffffffffffffff00000000000000000000000000000000166fffffffffffffffffffffffffffffffff9290921691909117905590505b9392505050565b600081156b033b2e3c9fd0803ce80000006002840419048411171561154657600080fd5b506b033b2e3c9fd0803ce80000009190910260028204010490565b6040517f23b872dd0000000000000000000000000000000000000000000000000000000080825273ffffffffffffffffffffffffffffffffffffffff8581166004840152841660248301526044820183905290600080606483828a5af16115cc573d6000803e3d6000fd5b506115d68561199c565b61163c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f475076323a206661696c6564207472616e7366657246726f6d00000000000000604482015260640161020f565b5050505050565b610160810151156116d3576000611664826101600151836102400151611a68565b905061167d8260e001518261142b90919063ffffffff16565b610100830181905261168e90610d2a565b6001840180547fffffffffffffffffffffffffffffffff00000000000000000000000000000000166fffffffffffffffffffffffffffffffff92909216919091179055505b8051156117605760006116f0826101800151836102400151611aaf565b905061170a8261012001518261142b90919063ffffffff16565b610140830181905261171b90610d2a565b6002840180547fffffffffffffffffffffffffffffffff00000000000000000000000000000000166fffffffffffffffffffffffffffffffff92909216919091179055505b5050565b61179d6040518060c001604052806000815260200160008152602001600081526020016000815260200160008152602001600081525090565b6101a08201516117ac57505050565b61012082015182516117bd9161142b565b602082015261014082015182516117d39161142b565b604082015260608201516102608301516102408401516117fb92919064ffffffffff16611ab8565b6060820181905260408301516118109161142b565b80825260208201516080840151604084015161182c919061214a565b6118369190612133565b6118409190612133565b608082018190526101a083015161185791906113e8565b60a08201819052156118df5761188261056c8361010001518360a0015161152290919063ffffffff16565b6008840180546000906118a89084906fffffffffffffffffffffffffffffffff16611f01565b92506101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff1602179055505b505050565b60006118f2825161ffff1690565b6118fe57506000611979565b60408051602081019091528354908190527faaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa1661193d57506001611979565b60408051602081019091528354815260009061195a908787611bff565b50509050801580156119755750825160d41c64ffffffffff16155b9150505b949350505050565b633b9aca00818102908104821461199757600080fd5b919050565b60006119dc565b7f08c379a00000000000000000000000000000000000000000000000000000000060005260206004528060245250806044525060646000fd5b3d8015611a1b5760208114611a5557611a167f475076323a206d616c666f726d6564207472616e7366657220726573756c7400601f6119a3565b611a62565b823b611a4c57611a4c7f475076323a206e6f74206120636f6e747261637400000000000000000000000060146119a3565b60019150611a62565b3d6000803e600051151591505b50919050565b600080611a7c64ffffffffff841642612133565b611a869085612061565b6301e1338090049050611aa5816b033b2e3c9fd0803ce800000061214a565b9150505b92915050565b600061151b8383425b600080611acc64ffffffffff851684612133565b905080611ae8576b033b2e3c9fd0803ce800000091505061151b565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81016000808060028511611b1e576000611b23565b600285035b925066038882915c4000611b378a8061142b565b81611b4457611b44612242565b0491506301e13380611b56838b61142b565b81611b6357611b63612242565b049050600082611b738688612061565b611b7d9190612061565b60029004905060008285611b91888a612061565b611b9b9190612061565b611ba59190612061565b60069004905080826301e13380611bbc8a8f612061565b611bc69190612271565b611bdc906b033b2e3c9fd0803ce800000061214a565b611be6919061214a565b611bf0919061214a565b9b9a5050505050505050505050565b6000806000611c0d86611cb7565b15611ca4576000611c3e877faaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa611cfb565b6000818152602087815260408083205473ffffffffffffffffffffffffffffffffffffffff168084528a8352818420825193840190925290549182905292935060d41c64ffffffffff1690508015611ca057600195509093509150611cae9050565b5050505b5060009150819050805b93509350939050565b80516000907faaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa16801580159061151b5750611cf3600182612133565b161592915050565b815160009082167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8101198116825b60029190911c9081156106e557600101611d2a565b6040518061028001604052806000815260200160008152602001600081526020016000815260200160008152602001600081526020016000815260200160008152602001600081526020016000815260200160008152602001600081526020016000815260200160008152602001611dc36040518060200160405280600081525090565b815260006020820181905260408201819052606082018190526080820181905260a09091015290565b73ffffffffffffffffffffffffffffffffffffffff81168114611e0e57600080fd5b50565b600080600080600080600060e0888a031215611e2c57600080fd5b8735965060208801359550604088013594506060880135611e4c81611dec565b93506080880135925060a0880135611e6381611dec565b915060c088013561ffff81168114611e7a57600080fd5b8091505092959891949750929550565b600080600080600060a08688031215611ea257600080fd5b853594506020860135611eb481611dec565b94979496505050506040830135926060810135926080909101359150565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006fffffffffffffffffffffffffffffffff808316818516808303821115611f2c57611f2c611ed2565b01949350505050565b600181815b80851115611f8e57817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04821115611f7457611f74611ed2565b80851615611f8157918102915b93841c9390800290611f3a565b509250929050565b600082611fa557506001611aa9565b81611fb257506000611aa9565b8160018114611fc85760028114611fd257611fee565b6001915050611aa9565b60ff841115611fe357611fe3611ed2565b50506001821b611aa9565b5060208310610133831016604e8410600b8410161715612011575081810a611aa9565b61201b8383611f35565b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0482111561204d5761204d611ed2565b029392505050565b600061151b8383611f96565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048311821515161561209957612099611ed2565b500290565b600060208083528351808285015260005b818110156120cb578581018301518582016040015282016120af565b818111156120dd576000604083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016929092016040019392505050565b60006020828403121561212357600080fd5b8151801515811461151b57600080fd5b60008282101561214557612145611ed2565b500390565b6000821982111561215d5761215d611ed2565b500190565b60006020828403121561217457600080fd5b5051919050565b60006fffffffffffffffffffffffffffffffff838116908316818110156121a4576121a4611ed2565b039392505050565b600080600080608085870312156121c257600080fd5b845193506020850151925060408501519150606085015164ffffffffff811681146121ec57600080fd5b939692955090935050565b60008060006060848603121561220c57600080fd5b8351925060208401519150604084015190509250925092565b60006020828403121561223757600080fd5b815161151b81611dec565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000826122a7577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b50049056fea2646970667358221220cee3e27afab70c3d80fa0ed71ff4a16129e046d54609375c52a4db847dacbba164736f6c634300080a0033",
}

// BridgeLogicABI is the input ABI used to generate the binding from.
// Deprecated: Use BridgeLogicMetaData.ABI instead.
var BridgeLogicABI = BridgeLogicMetaData.ABI

// BridgeLogicBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BridgeLogicMetaData.Bin instead.
var BridgeLogicBin = BridgeLogicMetaData.Bin

// DeployBridgeLogic deploys a new Ethereum contract, binding an instance of BridgeLogic to it.
func DeployBridgeLogic(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *BridgeLogic, error) {
	parsed, err := BridgeLogicMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BridgeLogicBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BridgeLogic{BridgeLogicCaller: BridgeLogicCaller{contract: contract}, BridgeLogicTransactor: BridgeLogicTransactor{contract: contract}, BridgeLogicFilterer: BridgeLogicFilterer{contract: contract}}, nil
}

// BridgeLogic is an auto generated Go binding around an Ethereum contract.
type BridgeLogic struct {
	BridgeLogicCaller     // Read-only binding to the contract
	BridgeLogicTransactor // Write-only binding to the contract
	BridgeLogicFilterer   // Log filterer for contract events
}

// BridgeLogicCaller is an auto generated read-only Go binding around an Ethereum contract.
type BridgeLogicCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeLogicTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BridgeLogicTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeLogicFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BridgeLogicFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeLogicSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BridgeLogicSession struct {
	Contract     *BridgeLogic      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BridgeLogicCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BridgeLogicCallerSession struct {
	Contract *BridgeLogicCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// BridgeLogicTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BridgeLogicTransactorSession struct {
	Contract     *BridgeLogicTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// BridgeLogicRaw is an auto generated low-level Go binding around an Ethereum contract.
type BridgeLogicRaw struct {
	Contract *BridgeLogic // Generic contract binding to access the raw methods on
}

// BridgeLogicCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BridgeLogicCallerRaw struct {
	Contract *BridgeLogicCaller // Generic read-only contract binding to access the raw methods on
}

// BridgeLogicTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BridgeLogicTransactorRaw struct {
	Contract *BridgeLogicTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBridgeLogic creates a new instance of BridgeLogic, bound to a specific deployed contract.
func NewBridgeLogic(address common.Address, backend bind.ContractBackend) (*BridgeLogic, error) {
	contract, err := bindBridgeLogic(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BridgeLogic{BridgeLogicCaller: BridgeLogicCaller{contract: contract}, BridgeLogicTransactor: BridgeLogicTransactor{contract: contract}, BridgeLogicFilterer: BridgeLogicFilterer{contract: contract}}, nil
}

// NewBridgeLogicCaller creates a new read-only instance of BridgeLogic, bound to a specific deployed contract.
func NewBridgeLogicCaller(address common.Address, caller bind.ContractCaller) (*BridgeLogicCaller, error) {
	contract, err := bindBridgeLogic(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeLogicCaller{contract: contract}, nil
}

// NewBridgeLogicTransactor creates a new write-only instance of BridgeLogic, bound to a specific deployed contract.
func NewBridgeLogicTransactor(address common.Address, transactor bind.ContractTransactor) (*BridgeLogicTransactor, error) {
	contract, err := bindBridgeLogic(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeLogicTransactor{contract: contract}, nil
}

// NewBridgeLogicFilterer creates a new log filterer instance of BridgeLogic, bound to a specific deployed contract.
func NewBridgeLogicFilterer(address common.Address, filterer bind.ContractFilterer) (*BridgeLogicFilterer, error) {
	contract, err := bindBridgeLogic(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BridgeLogicFilterer{contract: contract}, nil
}

// bindBridgeLogic binds a generic wrapper to an already deployed contract.
func bindBridgeLogic(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BridgeLogicMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BridgeLogic *BridgeLogicRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BridgeLogic.Contract.BridgeLogicCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BridgeLogic *BridgeLogicRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgeLogic.Contract.BridgeLogicTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BridgeLogic *BridgeLogicRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BridgeLogic.Contract.BridgeLogicTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BridgeLogic *BridgeLogicCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BridgeLogic.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BridgeLogic *BridgeLogicTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgeLogic.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BridgeLogic *BridgeLogicTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BridgeLogic.Contract.contract.Transact(opts, method, params...)
}

// BridgeLogicBackUnbackedIterator is returned from FilterBackUnbacked and is used to iterate over the raw logs and unpacked data for BackUnbacked events raised by the BridgeLogic contract.
type BridgeLogicBackUnbackedIterator struct {
	Event *BridgeLogicBackUnbacked // Event containing the contract specifics and raw log

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
func (it *BridgeLogicBackUnbackedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeLogicBackUnbacked)
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
		it.Event = new(BridgeLogicBackUnbacked)
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
func (it *BridgeLogicBackUnbackedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeLogicBackUnbackedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeLogicBackUnbacked represents a BackUnbacked event raised by the BridgeLogic contract.
type BridgeLogicBackUnbacked struct {
	Reserve common.Address
	Backer  common.Address
	Amount  *big.Int
	Fee     *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterBackUnbacked is a free log retrieval operation binding the contract event 0x281596e92b2d974beb7d4f124df30a0b39067b096893e95011ce4bdad798b759.
//
// Solidity: event BackUnbacked(address indexed reserve, address indexed backer, uint256 amount, uint256 fee)
func (_BridgeLogic *BridgeLogicFilterer) FilterBackUnbacked(opts *bind.FilterOpts, reserve []common.Address, backer []common.Address) (*BridgeLogicBackUnbackedIterator, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}
	var backerRule []interface{}
	for _, backerItem := range backer {
		backerRule = append(backerRule, backerItem)
	}

	logs, sub, err := _BridgeLogic.contract.FilterLogs(opts, "BackUnbacked", reserveRule, backerRule)
	if err != nil {
		return nil, err
	}
	return &BridgeLogicBackUnbackedIterator{contract: _BridgeLogic.contract, event: "BackUnbacked", logs: logs, sub: sub}, nil
}

// WatchBackUnbacked is a free log subscription operation binding the contract event 0x281596e92b2d974beb7d4f124df30a0b39067b096893e95011ce4bdad798b759.
//
// Solidity: event BackUnbacked(address indexed reserve, address indexed backer, uint256 amount, uint256 fee)
func (_BridgeLogic *BridgeLogicFilterer) WatchBackUnbacked(opts *bind.WatchOpts, sink chan<- *BridgeLogicBackUnbacked, reserve []common.Address, backer []common.Address) (event.Subscription, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}
	var backerRule []interface{}
	for _, backerItem := range backer {
		backerRule = append(backerRule, backerItem)
	}

	logs, sub, err := _BridgeLogic.contract.WatchLogs(opts, "BackUnbacked", reserveRule, backerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeLogicBackUnbacked)
				if err := _BridgeLogic.contract.UnpackLog(event, "BackUnbacked", log); err != nil {
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

// ParseBackUnbacked is a log parse operation binding the contract event 0x281596e92b2d974beb7d4f124df30a0b39067b096893e95011ce4bdad798b759.
//
// Solidity: event BackUnbacked(address indexed reserve, address indexed backer, uint256 amount, uint256 fee)
func (_BridgeLogic *BridgeLogicFilterer) ParseBackUnbacked(log types.Log) (*BridgeLogicBackUnbacked, error) {
	event := new(BridgeLogicBackUnbacked)
	if err := _BridgeLogic.contract.UnpackLog(event, "BackUnbacked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeLogicMintUnbackedIterator is returned from FilterMintUnbacked and is used to iterate over the raw logs and unpacked data for MintUnbacked events raised by the BridgeLogic contract.
type BridgeLogicMintUnbackedIterator struct {
	Event *BridgeLogicMintUnbacked // Event containing the contract specifics and raw log

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
func (it *BridgeLogicMintUnbackedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeLogicMintUnbacked)
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
		it.Event = new(BridgeLogicMintUnbacked)
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
func (it *BridgeLogicMintUnbackedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeLogicMintUnbackedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeLogicMintUnbacked represents a MintUnbacked event raised by the BridgeLogic contract.
type BridgeLogicMintUnbacked struct {
	Reserve      common.Address
	User         common.Address
	OnBehalfOf   common.Address
	Amount       *big.Int
	ReferralCode uint16
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterMintUnbacked is a free log retrieval operation binding the contract event 0xf25af37b3d3ec226063dc9bdc103ece7eb110a50f340fe854bb7bc1b0676d7d0.
//
// Solidity: event MintUnbacked(address indexed reserve, address user, address indexed onBehalfOf, uint256 amount, uint16 indexed referralCode)
func (_BridgeLogic *BridgeLogicFilterer) FilterMintUnbacked(opts *bind.FilterOpts, reserve []common.Address, onBehalfOf []common.Address, referralCode []uint16) (*BridgeLogicMintUnbackedIterator, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}

	var onBehalfOfRule []interface{}
	for _, onBehalfOfItem := range onBehalfOf {
		onBehalfOfRule = append(onBehalfOfRule, onBehalfOfItem)
	}

	var referralCodeRule []interface{}
	for _, referralCodeItem := range referralCode {
		referralCodeRule = append(referralCodeRule, referralCodeItem)
	}

	logs, sub, err := _BridgeLogic.contract.FilterLogs(opts, "MintUnbacked", reserveRule, onBehalfOfRule, referralCodeRule)
	if err != nil {
		return nil, err
	}
	return &BridgeLogicMintUnbackedIterator{contract: _BridgeLogic.contract, event: "MintUnbacked", logs: logs, sub: sub}, nil
}

// WatchMintUnbacked is a free log subscription operation binding the contract event 0xf25af37b3d3ec226063dc9bdc103ece7eb110a50f340fe854bb7bc1b0676d7d0.
//
// Solidity: event MintUnbacked(address indexed reserve, address user, address indexed onBehalfOf, uint256 amount, uint16 indexed referralCode)
func (_BridgeLogic *BridgeLogicFilterer) WatchMintUnbacked(opts *bind.WatchOpts, sink chan<- *BridgeLogicMintUnbacked, reserve []common.Address, onBehalfOf []common.Address, referralCode []uint16) (event.Subscription, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}

	var onBehalfOfRule []interface{}
	for _, onBehalfOfItem := range onBehalfOf {
		onBehalfOfRule = append(onBehalfOfRule, onBehalfOfItem)
	}

	var referralCodeRule []interface{}
	for _, referralCodeItem := range referralCode {
		referralCodeRule = append(referralCodeRule, referralCodeItem)
	}

	logs, sub, err := _BridgeLogic.contract.WatchLogs(opts, "MintUnbacked", reserveRule, onBehalfOfRule, referralCodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeLogicMintUnbacked)
				if err := _BridgeLogic.contract.UnpackLog(event, "MintUnbacked", log); err != nil {
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

// ParseMintUnbacked is a log parse operation binding the contract event 0xf25af37b3d3ec226063dc9bdc103ece7eb110a50f340fe854bb7bc1b0676d7d0.
//
// Solidity: event MintUnbacked(address indexed reserve, address user, address indexed onBehalfOf, uint256 amount, uint16 indexed referralCode)
func (_BridgeLogic *BridgeLogicFilterer) ParseMintUnbacked(log types.Log) (*BridgeLogicMintUnbacked, error) {
	event := new(BridgeLogicMintUnbacked)
	if err := _BridgeLogic.contract.UnpackLog(event, "MintUnbacked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeLogicReserveUsedAsCollateralEnabledIterator is returned from FilterReserveUsedAsCollateralEnabled and is used to iterate over the raw logs and unpacked data for ReserveUsedAsCollateralEnabled events raised by the BridgeLogic contract.
type BridgeLogicReserveUsedAsCollateralEnabledIterator struct {
	Event *BridgeLogicReserveUsedAsCollateralEnabled // Event containing the contract specifics and raw log

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
func (it *BridgeLogicReserveUsedAsCollateralEnabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeLogicReserveUsedAsCollateralEnabled)
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
		it.Event = new(BridgeLogicReserveUsedAsCollateralEnabled)
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
func (it *BridgeLogicReserveUsedAsCollateralEnabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeLogicReserveUsedAsCollateralEnabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeLogicReserveUsedAsCollateralEnabled represents a ReserveUsedAsCollateralEnabled event raised by the BridgeLogic contract.
type BridgeLogicReserveUsedAsCollateralEnabled struct {
	Reserve common.Address
	User    common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterReserveUsedAsCollateralEnabled is a free log retrieval operation binding the contract event 0x00058a56ea94653cdf4f152d227ace22d4c00ad99e2a43f58cb7d9e3feb295f2.
//
// Solidity: event ReserveUsedAsCollateralEnabled(address indexed reserve, address indexed user)
func (_BridgeLogic *BridgeLogicFilterer) FilterReserveUsedAsCollateralEnabled(opts *bind.FilterOpts, reserve []common.Address, user []common.Address) (*BridgeLogicReserveUsedAsCollateralEnabledIterator, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _BridgeLogic.contract.FilterLogs(opts, "ReserveUsedAsCollateralEnabled", reserveRule, userRule)
	if err != nil {
		return nil, err
	}
	return &BridgeLogicReserveUsedAsCollateralEnabledIterator{contract: _BridgeLogic.contract, event: "ReserveUsedAsCollateralEnabled", logs: logs, sub: sub}, nil
}

// WatchReserveUsedAsCollateralEnabled is a free log subscription operation binding the contract event 0x00058a56ea94653cdf4f152d227ace22d4c00ad99e2a43f58cb7d9e3feb295f2.
//
// Solidity: event ReserveUsedAsCollateralEnabled(address indexed reserve, address indexed user)
func (_BridgeLogic *BridgeLogicFilterer) WatchReserveUsedAsCollateralEnabled(opts *bind.WatchOpts, sink chan<- *BridgeLogicReserveUsedAsCollateralEnabled, reserve []common.Address, user []common.Address) (event.Subscription, error) {

	var reserveRule []interface{}
	for _, reserveItem := range reserve {
		reserveRule = append(reserveRule, reserveItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _BridgeLogic.contract.WatchLogs(opts, "ReserveUsedAsCollateralEnabled", reserveRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeLogicReserveUsedAsCollateralEnabled)
				if err := _BridgeLogic.contract.UnpackLog(event, "ReserveUsedAsCollateralEnabled", log); err != nil {
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

// ParseReserveUsedAsCollateralEnabled is a log parse operation binding the contract event 0x00058a56ea94653cdf4f152d227ace22d4c00ad99e2a43f58cb7d9e3feb295f2.
//
// Solidity: event ReserveUsedAsCollateralEnabled(address indexed reserve, address indexed user)
func (_BridgeLogic *BridgeLogicFilterer) ParseReserveUsedAsCollateralEnabled(log types.Log) (*BridgeLogicReserveUsedAsCollateralEnabled, error) {
	event := new(BridgeLogicReserveUsedAsCollateralEnabled)
	if err := _BridgeLogic.contract.UnpackLog(event, "ReserveUsedAsCollateralEnabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
