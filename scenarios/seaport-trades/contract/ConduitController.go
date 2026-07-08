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

// ConduitControllerMetaData contains all meta data concerning the ConduitController contract.
var ConduitControllerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"conduit\",\"type\":\"address\"}],\"name\":\"CallerIsNotNewPotentialOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"conduit\",\"type\":\"address\"}],\"name\":\"CallerIsNotOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"conduit\",\"type\":\"address\"}],\"name\":\"ChannelOutOfRange\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"conduit\",\"type\":\"address\"}],\"name\":\"ConduitAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCreator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"conduit\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"newPotentialOwner\",\"type\":\"address\"}],\"name\":\"NewPotentialOwnerAlreadySet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"conduit\",\"type\":\"address\"}],\"name\":\"NewPotentialOwnerIsZeroAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoConduit\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"conduit\",\"type\":\"address\"}],\"name\":\"NoPotentialOwnerCurrentlySet\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"conduit\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"conduitKey\",\"type\":\"bytes32\"}],\"name\":\"NewConduit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"conduit\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newPotentialOwner\",\"type\":\"address\"}],\"name\":\"PotentialOwnerUpdated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"conduit\",\"type\":\"address\"}],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"conduit\",\"type\":\"address\"}],\"name\":\"cancelOwnershipTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"conduitKey\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"initialOwner\",\"type\":\"address\"}],\"name\":\"createConduit\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"conduit\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"conduit\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"channelIndex\",\"type\":\"uint256\"}],\"name\":\"getChannel\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"channel\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"conduit\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"channel\",\"type\":\"address\"}],\"name\":\"getChannelStatus\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isOpen\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"conduit\",\"type\":\"address\"}],\"name\":\"getChannels\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"channels\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"conduitKey\",\"type\":\"bytes32\"}],\"name\":\"getConduit\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"conduit\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"exists\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConduitCodeHashes\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"creationCodeHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"runtimeCodeHash\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"conduit\",\"type\":\"address\"}],\"name\":\"getKey\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"conduitKey\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"conduit\",\"type\":\"address\"}],\"name\":\"getPotentialOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"potentialOwner\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"conduit\",\"type\":\"address\"}],\"name\":\"getTotalChannels\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"totalChannels\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"conduit\",\"type\":\"address\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"conduit\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"newPotentialOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"conduit\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"channel\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isOpen\",\"type\":\"bool\"}],\"name\":\"updateChannel\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60c060405234801561000f575f80fd5b5060405161001f60208201610082565b6020820181038252601f19601f8201166040525080519060200120608081815250505f805f1b60405161005190610082565b8190604051809103905ff590508015801561006e573d5f803e3d5ffd5b506001600160a01b03163f60a0525061008f565b610a6a8061197083390190565b60805160a0516118a46100cc5f395f81816101420152818161081101526108df01525f818161011f015281816107b9015261089b01526118a45ff3fe608060405234801561000f575f80fd5b50600436106100e5575f3560e01c80636d435421116100885780637b37e561116100635780637b37e561146102435780638b9e028b14610256578063906c87cc1461027657806393790f4414610289575f80fd5b80636d435421146101eb5780636e9bfd9f146101fe578063794593bc14610230575f80fd5b806314afd79e116100c357806314afd79e1461018157806333bc8572146101945780634e3f9580146101b757806351710e45146101d8575f80fd5b8063027cc764146100e95780630a96ad391461011957806313ad9cab1461016c575b5f80fd5b6100fc6100f7366004610c39565b61029c565b6040516001600160a01b0390911681526020015b60405180910390f35b604080517f000000000000000000000000000000000000000000000000000000000000000081527f0000000000000000000000000000000000000000000000000000000000000000602082015201610110565b61017f61017a366004610c61565b610339565b005b6100fc61018f366004610ca9565b610531565b6101a76101a2366004610cc9565b61055c565b6040519015158152602001610110565b6101ca6101c5366004610ca9565b610595565b604051908152602001610110565b61017f6101e6366004610ca9565b6105bd565b61017f6101f9366004610cc9565b6106be565b61021161020c366004610cfa565b6107ae565b604080516001600160a01b039093168352901515602083015201610110565b6100fc61023e366004610d11565b610838565b61017f610251366004610ca9565b610a0a565b610269610264366004610ca9565b610aab565b6040516101109190610d32565b6100fc610284366004610ca9565b610b28565b6101ca610297366004610ca9565b610b53565b5f6102a683610b8e565b6001600160a01b0383165f908152602081905260409020600301548083106102f157604051636ceb340b60e01b81526001600160a01b03851660048201526024015b60405180910390fd5b6001600160a01b0384165f90815260208190526040902060030180548490811061031d5761031d610d7e565b5f918252602090912001546001600160a01b0316949350505050565b61034283610bc6565b60405163c4e8fcb560e01b81526001600160a01b038381166004830152821515602483015284169063c4e8fcb5906044015f604051808303815f87803b15801561038a575f80fd5b505af115801561039c573d5f803e3d5ffd5b505050506001600160a01b038381165f908152602081815260408083209386168352600484019091529020548015158380156103d6575080155b15610425576003830180546001810182555f828152602080822090920180546001600160a01b0319166001600160a01b038a169081179091559254928152600486019091526040902055610529565b831580156104305750805b156105295760038301545f198301905f9061044d90600190610d92565b90508181146104d8575f85600301828154811061046c5761046c610d7e565b5f918252602090912001546003870180546001600160a01b03909216925082918590811061049c5761049c610d7e565b5f91825260208083209190910180546001600160a01b0319166001600160a01b0394851617905592909116815260048701909152604090208490555b846003018054806104eb576104eb610db7565b5f828152602080822083015f1990810180546001600160a01b03191690559092019092556001600160a01b0389168252600487019052604081205550505b505050505050565b5f61053b82610b8e565b506001600160a01b039081165f908152602081905260409020600101541690565b5f61056683610b8e565b506001600160a01b039182165f9081526020818152604080832093909416825260049092019091522054151590565b5f61059f82610b8e565b506001600160a01b03165f9081526020819052604090206003015490565b6105c681610b8e565b6001600160a01b038181165f9081526020819052604090206002015416331461060d576040516388c3a11560e01b81526001600160a01b03821660048201526024016102e8565b6040515f907f11a3cf439fb225bfe74225716b6774765670ec1060e3796802e62139d69974da908290a26001600160a01b038082165f818152602081905260408082206002810180546001600160a01b031916905560010154905133949190911692917fc8894f26f396ce8c004245c8b7cd1b92103a6e4302fcbab883987149ac01b7ec91a46001600160a01b03165f90815260208190526040902060010180546001600160a01b03191633179055565b6106c782610bc6565b6001600160a01b0381166106f95760405163a388d26360e01b81526001600160a01b03831660048201526024016102e8565b6001600160a01b038083165f9081526020819052604090206002015481169082160361074b576040516365e0406560e11b81526001600160a01b038084166004830152821660248201526044016102e8565b6040516001600160a01b038216907f11a3cf439fb225bfe74225716b6774765670ec1060e3796802e62139d69974da905f90a26001600160a01b039182165f90815260208190526040902060020180546001600160a01b03191691909216179055565b5f8060ff60f81b30847f00000000000000000000000000000000000000000000000000000000000000006040516020016107eb9493929190610dcb565b60408051601f198184030181529190528051602090910120936001600160a01b0385163f7f0000000000000000000000000000000000000000000000000000000000000000149350915050565b5f6001600160a01b0382166108605760405163267eaa8160e21b815260040160405180910390fd5b606083901c3314610884576040516332db94d160e21b815260040160405180910390fd5b6040516108c3906001600160f81b031990309086907f000000000000000000000000000000000000000000000000000000000000000090602001610dcb565b604051602081830303815290604052805190602001205f1c90507f0000000000000000000000000000000000000000000000000000000000000000816001600160a01b03163f0361093257604051633194665960e11b81526001600160a01b03821660048201526024016102e8565b8260405161093f90610c16565b8190604051809103905ff590508015801561095c573d5f803e3d5ffd5b50506001600160a01b038181165f81815260208181526040918290206001810180546001600160a01b03191695881695909517909455868455815192835282018690527f4397af6128d529b8ae0442f99db1296d5136062597a15bbc61c1b2a6431a7d15910160405180910390a16040516001600160a01b03808516915f918516907fc8894f26f396ce8c004245c8b7cd1b92103a6e4302fcbab883987149ac01b7ec908390a45092915050565b610a1381610bc6565b6001600160a01b038181165f9081526020819052604090206002015416610a58576040516335809b0b60e11b81526001600160a01b03821660048201526024016102e8565b6040515f907f11a3cf439fb225bfe74225716b6774765670ec1060e3796802e62139d69974da908290a26001600160a01b03165f90815260208190526040902060020180546001600160a01b0319169055565b6060610ab682610b8e565b6001600160a01b0382165f908152602081815260409182902060030180548351818402810184019094528084529091830182828015610b1c57602002820191905f5260205f20905b81546001600160a01b03168152600190910190602001808311610afe575b50505050509050919050565b5f610b3282610b8e565b506001600160a01b039081165f908152602081905260409020600201541690565b6001600160a01b0381165f9081526020819052604090205480610b89576040516304ca820960e41b815260040160405180910390fd5b919050565b6001600160a01b0381165f90815260208190526040902054610bc3576040516304ca820960e41b815260040160405180910390fd5b50565b610bcf81610b8e565b6001600160a01b038181165f90815260208190526040902060010154163314610bc35760405163d4ed9a1760e01b81526001600160a01b03821660048201526024016102e8565b610a6a80610e0583390190565b80356001600160a01b0381168114610b89575f80fd5b5f8060408385031215610c4a575f80fd5b610c5383610c23565b946020939093013593505050565b5f805f60608486031215610c73575f80fd5b610c7c84610c23565b9250610c8a60208501610c23565b915060408401358015158114610c9e575f80fd5b809150509250925092565b5f60208284031215610cb9575f80fd5b610cc282610c23565b9392505050565b5f8060408385031215610cda575f80fd5b610ce383610c23565b9150610cf160208401610c23565b90509250929050565b5f60208284031215610d0a575f80fd5b5035919050565b5f8060408385031215610d22575f80fd5b82359150610cf160208401610c23565b602080825282518282018190525f9190848201906040850190845b81811015610d725783516001600160a01b031683529284019291840191600101610d4d565b50909695505050505050565b634e487b7160e01b5f52603260045260245ffd5b81810381811115610db157634e487b7160e01b5f52601160045260245ffd5b92915050565b634e487b7160e01b5f52603160045260245ffd5b6001600160f81b031994909416845260609290921b6bffffffffffffffffffffffff19166001840152601583015260358201526055019056fe60a060405234801561000f575f80fd5b5033608052608051610a3e61002c5f395f6101d20152610a3e5ff3fe608060405234801561000f575f80fd5b506004361061004a575f3560e01c80634ce34aa21461004e578063899e104c1461007e5780638df25d9214610091578063c4e8fcb5146100a4575b5f80fd5b61006161005c366004610834565b6100b9565b6040516001600160e01b0319909116815260200160405180910390f35b61006161008c3660046108b4565b610121565b61006161009f36600461091b565b61018a565b6100b76100b2366004610969565b6101c7565b005b5f335f525f60205260405f20546100dd576349ed56f960e11b5f523360045260245ffd5b815f5b81811015610110576101088585838181106100fd576100fd6109a2565b905060c002016102c4565b6001016100e0565b50632671a55160e11b949350505050565b5f335f525f60205260405f2054610145576349ed56f960e11b5f523360045260245ffd5b835f5b8181101561016d576101658787838181106100fd576100fd6109a2565b600101610148565b506101788484610436565b50632267841360e21b95945050505050565b5f335f525f60205260405f20546101ae576349ed56f960e11b5f523360045260245ffd5b6101b88383610436565b506346f92ec960e11b92915050565b336001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614610210576040516336abb4df60e11b815260040160405180910390fd5b6001600160a01b0382165f9081526020819052604090205481151560ff909116151503610268576040516349271a0f60e11b81526001600160a01b038316600482015281151560248201526044015b60405180910390fd5b6001600160a01b0382165f8181526020818152604091829020805460ff191685151590811790915591519182527fae63067d43ac07563b7eb8db6595635fc77f1578a2a5ea06ba91b63e2afa37e2910160405180910390a25050565b60016102d360208301836109ca565b60038111156102e4576102e46109b6565b03610329576103266102fc60408301602084016109ef565b61030c60608401604085016109ef565b61031c60808501606086016109ef565b8460a0013561056c565b50565b600261033860208301836109ca565b6003811115610349576103496109b6565b036103b6578060a00135600114610379576040516369f9582760e01b815260a0820135600482015260240161025f565b61032661038c60408301602084016109ef565b61039c60608401604085016109ef565b6103ac60808501606086016109ef565b8460800135610660565b60036103c560208301836109ca565b60038111156103d6576103d66109b6565b0361041d576103266103ee60408301602084016109ef565b6103fe60608401604085016109ef565b61040e60808501606086016109ef565b84608001358560a00135610715565b604051631e4cbc7f60e21b815260040160405180910390fd5b808280631759616b60e11b6020525f5b8381101561055f57823582018035803b61046b57635f15d6725f52806020526024601cfd5b60a08201358060051b60c0018060808501351460a0606086013514168185013583141615905080156104a657633ae8821360e21b5f5260045ffd5b506020860195506080602084016024378060061b60400190508060a00160a4525f8160c401528060c4018160a0850160c4375f808260205f875af1935083610550573d1561053057601f3d0160051c91508060051c826003028184111561051a578184036003028280028580020360091c01015b5a60208201101561052d573d5f803e3d5ffd5b50505b6357e222f160e11b5f528260045260c0606452608451602001608452805ffd5b50505050600181019050610446565b5050505060806040525050565b6040516323b872dd60e01b5f5283600452826024528160445260205f60645f80895af1803d15601f3d1160015f51141617163d151581166106515780873b15151661065157806106405781610623573d1561060257601f3d0160051c8360051c81600302818311156105eb578183036003028280028480020360091c01015b5a6020820110156105fe573d5f803e3d5ffd5b5050505b63f486bc875f528660205285604052846060525f6080528360a05260a4601cfd5b63988919235f52866020528560405284606052836080526084601cfd5b635f15d6725f52866020526024601cfd5b505060405250505f6060525050565b833b61067757635f15d6725f52836020526024601cfd5b6040516323b872dd60e01b5f528360045282602452816044525f8060645f80895af180610707573d156106e557601f3d0160051c8260051c81600302818311156106ce578183036003028280028480020360091c01015b5a6020820110156106e1573d5f803e3d5ffd5b5050505b63f486bc875f5285602052846040528360605282608052600160a05260a4601cfd5b5060405250505f6060525050565b843b61072c57635f15d6725f52846020526024601cfd5b60405160805160a05160c051637921219560e11b5f528760045286602452856044528460645260a06084525f60a4525f8060c45f808d5af1806107d1573d156107b057601f3d0160051c8560051c8160030281831115610799578183036003028280028480020360091c01015b5a6020820110156107ac573d5f803e3d5ffd5b5050505b63f486bc875f52896020528860405287606052866080528560a05260a4601cfd5b5060809290925260a05260c05260405250505f606052505050565b5f8083601f8401126107fc575f80fd5b50813567ffffffffffffffff811115610813575f80fd5b60208301915083602060c08302850101111561082d575f80fd5b9250929050565b5f8060208385031215610845575f80fd5b823567ffffffffffffffff81111561085b575f80fd5b610867858286016107ec565b90969095509350505050565b5f8083601f840112610883575f80fd5b50813567ffffffffffffffff81111561089a575f80fd5b6020830191508360208260051b850101111561082d575f80fd5b5f805f80604085870312156108c7575f80fd5b843567ffffffffffffffff808211156108de575f80fd5b6108ea888389016107ec565b90965094506020870135915080821115610902575f80fd5b5061090f87828801610873565b95989497509550505050565b5f806020838503121561092c575f80fd5b823567ffffffffffffffff811115610942575f80fd5b61086785828601610873565b80356001600160a01b0381168114610964575f80fd5b919050565b5f806040838503121561097a575f80fd5b6109838361094e565b915060208301358015158114610997575f80fd5b809150509250929050565b634e487b7160e01b5f52603260045260245ffd5b634e487b7160e01b5f52602160045260245ffd5b5f602082840312156109da575f80fd5b8135600481106109e8575f80fd5b9392505050565b5f602082840312156109ff575f80fd5b6109e88261094e56fea2646970667358221220021644d03af3d878d2df970c834f2e698b8641d72b762e16f09f7873e278882364736f6c63430008180033a26469706673582212202d21d57ec345ef2e545a19cc8167534bf51cc9b6e2ce4e467994984c62edf10964736f6c6343000818003360a060405234801561000f575f80fd5b5033608052608051610a3e61002c5f395f6101d20152610a3e5ff3fe608060405234801561000f575f80fd5b506004361061004a575f3560e01c80634ce34aa21461004e578063899e104c1461007e5780638df25d9214610091578063c4e8fcb5146100a4575b5f80fd5b61006161005c366004610834565b6100b9565b6040516001600160e01b0319909116815260200160405180910390f35b61006161008c3660046108b4565b610121565b61006161009f36600461091b565b61018a565b6100b76100b2366004610969565b6101c7565b005b5f335f525f60205260405f20546100dd576349ed56f960e11b5f523360045260245ffd5b815f5b81811015610110576101088585838181106100fd576100fd6109a2565b905060c002016102c4565b6001016100e0565b50632671a55160e11b949350505050565b5f335f525f60205260405f2054610145576349ed56f960e11b5f523360045260245ffd5b835f5b8181101561016d576101658787838181106100fd576100fd6109a2565b600101610148565b506101788484610436565b50632267841360e21b95945050505050565b5f335f525f60205260405f20546101ae576349ed56f960e11b5f523360045260245ffd5b6101b88383610436565b506346f92ec960e11b92915050565b336001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614610210576040516336abb4df60e11b815260040160405180910390fd5b6001600160a01b0382165f9081526020819052604090205481151560ff909116151503610268576040516349271a0f60e11b81526001600160a01b038316600482015281151560248201526044015b60405180910390fd5b6001600160a01b0382165f8181526020818152604091829020805460ff191685151590811790915591519182527fae63067d43ac07563b7eb8db6595635fc77f1578a2a5ea06ba91b63e2afa37e2910160405180910390a25050565b60016102d360208301836109ca565b60038111156102e4576102e46109b6565b03610329576103266102fc60408301602084016109ef565b61030c60608401604085016109ef565b61031c60808501606086016109ef565b8460a0013561056c565b50565b600261033860208301836109ca565b6003811115610349576103496109b6565b036103b6578060a00135600114610379576040516369f9582760e01b815260a0820135600482015260240161025f565b61032661038c60408301602084016109ef565b61039c60608401604085016109ef565b6103ac60808501606086016109ef565b8460800135610660565b60036103c560208301836109ca565b60038111156103d6576103d66109b6565b0361041d576103266103ee60408301602084016109ef565b6103fe60608401604085016109ef565b61040e60808501606086016109ef565b84608001358560a00135610715565b604051631e4cbc7f60e21b815260040160405180910390fd5b808280631759616b60e11b6020525f5b8381101561055f57823582018035803b61046b57635f15d6725f52806020526024601cfd5b60a08201358060051b60c0018060808501351460a0606086013514168185013583141615905080156104a657633ae8821360e21b5f5260045ffd5b506020860195506080602084016024378060061b60400190508060a00160a4525f8160c401528060c4018160a0850160c4375f808260205f875af1935083610550573d1561053057601f3d0160051c91508060051c826003028184111561051a578184036003028280028580020360091c01015b5a60208201101561052d573d5f803e3d5ffd5b50505b6357e222f160e11b5f528260045260c0606452608451602001608452805ffd5b50505050600181019050610446565b5050505060806040525050565b6040516323b872dd60e01b5f5283600452826024528160445260205f60645f80895af1803d15601f3d1160015f51141617163d151581166106515780873b15151661065157806106405781610623573d1561060257601f3d0160051c8360051c81600302818311156105eb578183036003028280028480020360091c01015b5a6020820110156105fe573d5f803e3d5ffd5b5050505b63f486bc875f528660205285604052846060525f6080528360a05260a4601cfd5b63988919235f52866020528560405284606052836080526084601cfd5b635f15d6725f52866020526024601cfd5b505060405250505f6060525050565b833b61067757635f15d6725f52836020526024601cfd5b6040516323b872dd60e01b5f528360045282602452816044525f8060645f80895af180610707573d156106e557601f3d0160051c8260051c81600302818311156106ce578183036003028280028480020360091c01015b5a6020820110156106e1573d5f803e3d5ffd5b5050505b63f486bc875f5285602052846040528360605282608052600160a05260a4601cfd5b5060405250505f6060525050565b843b61072c57635f15d6725f52846020526024601cfd5b60405160805160a05160c051637921219560e11b5f528760045286602452856044528460645260a06084525f60a4525f8060c45f808d5af1806107d1573d156107b057601f3d0160051c8560051c8160030281831115610799578183036003028280028480020360091c01015b5a6020820110156107ac573d5f803e3d5ffd5b5050505b63f486bc875f52896020528860405287606052866080528560a05260a4601cfd5b5060809290925260a05260c05260405250505f606052505050565b5f8083601f8401126107fc575f80fd5b50813567ffffffffffffffff811115610813575f80fd5b60208301915083602060c08302850101111561082d575f80fd5b9250929050565b5f8060208385031215610845575f80fd5b823567ffffffffffffffff81111561085b575f80fd5b610867858286016107ec565b90969095509350505050565b5f8083601f840112610883575f80fd5b50813567ffffffffffffffff81111561089a575f80fd5b6020830191508360208260051b850101111561082d575f80fd5b5f805f80604085870312156108c7575f80fd5b843567ffffffffffffffff808211156108de575f80fd5b6108ea888389016107ec565b90965094506020870135915080821115610902575f80fd5b5061090f87828801610873565b95989497509550505050565b5f806020838503121561092c575f80fd5b823567ffffffffffffffff811115610942575f80fd5b61086785828601610873565b80356001600160a01b0381168114610964575f80fd5b919050565b5f806040838503121561097a575f80fd5b6109838361094e565b915060208301358015158114610997575f80fd5b809150509250929050565b634e487b7160e01b5f52603260045260245ffd5b634e487b7160e01b5f52602160045260245ffd5b5f602082840312156109da575f80fd5b8135600481106109e8575f80fd5b9392505050565b5f602082840312156109ff575f80fd5b6109e88261094e56fea2646970667358221220021644d03af3d878d2df970c834f2e698b8641d72b762e16f09f7873e278882364736f6c63430008180033",
}

// ConduitControllerABI is the input ABI used to generate the binding from.
// Deprecated: Use ConduitControllerMetaData.ABI instead.
var ConduitControllerABI = ConduitControllerMetaData.ABI

// ConduitControllerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ConduitControllerMetaData.Bin instead.
var ConduitControllerBin = ConduitControllerMetaData.Bin

// DeployConduitController deploys a new Ethereum contract, binding an instance of ConduitController to it.
func DeployConduitController(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ConduitController, error) {
	parsed, err := ConduitControllerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ConduitControllerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ConduitController{ConduitControllerCaller: ConduitControllerCaller{contract: contract}, ConduitControllerTransactor: ConduitControllerTransactor{contract: contract}, ConduitControllerFilterer: ConduitControllerFilterer{contract: contract}}, nil
}

// ConduitController is an auto generated Go binding around an Ethereum contract.
type ConduitController struct {
	ConduitControllerCaller     // Read-only binding to the contract
	ConduitControllerTransactor // Write-only binding to the contract
	ConduitControllerFilterer   // Log filterer for contract events
}

// ConduitControllerCaller is an auto generated read-only Go binding around an Ethereum contract.
type ConduitControllerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConduitControllerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ConduitControllerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConduitControllerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ConduitControllerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConduitControllerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ConduitControllerSession struct {
	Contract     *ConduitController // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ConduitControllerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ConduitControllerCallerSession struct {
	Contract *ConduitControllerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// ConduitControllerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ConduitControllerTransactorSession struct {
	Contract     *ConduitControllerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// ConduitControllerRaw is an auto generated low-level Go binding around an Ethereum contract.
type ConduitControllerRaw struct {
	Contract *ConduitController // Generic contract binding to access the raw methods on
}

// ConduitControllerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ConduitControllerCallerRaw struct {
	Contract *ConduitControllerCaller // Generic read-only contract binding to access the raw methods on
}

// ConduitControllerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ConduitControllerTransactorRaw struct {
	Contract *ConduitControllerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewConduitController creates a new instance of ConduitController, bound to a specific deployed contract.
func NewConduitController(address common.Address, backend bind.ContractBackend) (*ConduitController, error) {
	contract, err := bindConduitController(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ConduitController{ConduitControllerCaller: ConduitControllerCaller{contract: contract}, ConduitControllerTransactor: ConduitControllerTransactor{contract: contract}, ConduitControllerFilterer: ConduitControllerFilterer{contract: contract}}, nil
}

// NewConduitControllerCaller creates a new read-only instance of ConduitController, bound to a specific deployed contract.
func NewConduitControllerCaller(address common.Address, caller bind.ContractCaller) (*ConduitControllerCaller, error) {
	contract, err := bindConduitController(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConduitControllerCaller{contract: contract}, nil
}

// NewConduitControllerTransactor creates a new write-only instance of ConduitController, bound to a specific deployed contract.
func NewConduitControllerTransactor(address common.Address, transactor bind.ContractTransactor) (*ConduitControllerTransactor, error) {
	contract, err := bindConduitController(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConduitControllerTransactor{contract: contract}, nil
}

// NewConduitControllerFilterer creates a new log filterer instance of ConduitController, bound to a specific deployed contract.
func NewConduitControllerFilterer(address common.Address, filterer bind.ContractFilterer) (*ConduitControllerFilterer, error) {
	contract, err := bindConduitController(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConduitControllerFilterer{contract: contract}, nil
}

// bindConduitController binds a generic wrapper to an already deployed contract.
func bindConduitController(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ConduitControllerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConduitController *ConduitControllerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConduitController.Contract.ConduitControllerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConduitController *ConduitControllerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConduitController.Contract.ConduitControllerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConduitController *ConduitControllerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConduitController.Contract.ConduitControllerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConduitController *ConduitControllerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConduitController.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConduitController *ConduitControllerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConduitController.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConduitController *ConduitControllerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConduitController.Contract.contract.Transact(opts, method, params...)
}

// GetChannel is a free data retrieval call binding the contract method 0x027cc764.
//
// Solidity: function getChannel(address conduit, uint256 channelIndex) view returns(address channel)
func (_ConduitController *ConduitControllerCaller) GetChannel(opts *bind.CallOpts, conduit common.Address, channelIndex *big.Int) (common.Address, error) {
	var out []interface{}
	err := _ConduitController.contract.Call(opts, &out, "getChannel", conduit, channelIndex)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetChannel is a free data retrieval call binding the contract method 0x027cc764.
//
// Solidity: function getChannel(address conduit, uint256 channelIndex) view returns(address channel)
func (_ConduitController *ConduitControllerSession) GetChannel(conduit common.Address, channelIndex *big.Int) (common.Address, error) {
	return _ConduitController.Contract.GetChannel(&_ConduitController.CallOpts, conduit, channelIndex)
}

// GetChannel is a free data retrieval call binding the contract method 0x027cc764.
//
// Solidity: function getChannel(address conduit, uint256 channelIndex) view returns(address channel)
func (_ConduitController *ConduitControllerCallerSession) GetChannel(conduit common.Address, channelIndex *big.Int) (common.Address, error) {
	return _ConduitController.Contract.GetChannel(&_ConduitController.CallOpts, conduit, channelIndex)
}

// GetChannelStatus is a free data retrieval call binding the contract method 0x33bc8572.
//
// Solidity: function getChannelStatus(address conduit, address channel) view returns(bool isOpen)
func (_ConduitController *ConduitControllerCaller) GetChannelStatus(opts *bind.CallOpts, conduit common.Address, channel common.Address) (bool, error) {
	var out []interface{}
	err := _ConduitController.contract.Call(opts, &out, "getChannelStatus", conduit, channel)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetChannelStatus is a free data retrieval call binding the contract method 0x33bc8572.
//
// Solidity: function getChannelStatus(address conduit, address channel) view returns(bool isOpen)
func (_ConduitController *ConduitControllerSession) GetChannelStatus(conduit common.Address, channel common.Address) (bool, error) {
	return _ConduitController.Contract.GetChannelStatus(&_ConduitController.CallOpts, conduit, channel)
}

// GetChannelStatus is a free data retrieval call binding the contract method 0x33bc8572.
//
// Solidity: function getChannelStatus(address conduit, address channel) view returns(bool isOpen)
func (_ConduitController *ConduitControllerCallerSession) GetChannelStatus(conduit common.Address, channel common.Address) (bool, error) {
	return _ConduitController.Contract.GetChannelStatus(&_ConduitController.CallOpts, conduit, channel)
}

// GetChannels is a free data retrieval call binding the contract method 0x8b9e028b.
//
// Solidity: function getChannels(address conduit) view returns(address[] channels)
func (_ConduitController *ConduitControllerCaller) GetChannels(opts *bind.CallOpts, conduit common.Address) ([]common.Address, error) {
	var out []interface{}
	err := _ConduitController.contract.Call(opts, &out, "getChannels", conduit)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetChannels is a free data retrieval call binding the contract method 0x8b9e028b.
//
// Solidity: function getChannels(address conduit) view returns(address[] channels)
func (_ConduitController *ConduitControllerSession) GetChannels(conduit common.Address) ([]common.Address, error) {
	return _ConduitController.Contract.GetChannels(&_ConduitController.CallOpts, conduit)
}

// GetChannels is a free data retrieval call binding the contract method 0x8b9e028b.
//
// Solidity: function getChannels(address conduit) view returns(address[] channels)
func (_ConduitController *ConduitControllerCallerSession) GetChannels(conduit common.Address) ([]common.Address, error) {
	return _ConduitController.Contract.GetChannels(&_ConduitController.CallOpts, conduit)
}

// GetConduit is a free data retrieval call binding the contract method 0x6e9bfd9f.
//
// Solidity: function getConduit(bytes32 conduitKey) view returns(address conduit, bool exists)
func (_ConduitController *ConduitControllerCaller) GetConduit(opts *bind.CallOpts, conduitKey [32]byte) (struct {
	Conduit common.Address
	Exists  bool
}, error) {
	var out []interface{}
	err := _ConduitController.contract.Call(opts, &out, "getConduit", conduitKey)

	outstruct := new(struct {
		Conduit common.Address
		Exists  bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Conduit = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Exists = *abi.ConvertType(out[1], new(bool)).(*bool)

	return *outstruct, err

}

// GetConduit is a free data retrieval call binding the contract method 0x6e9bfd9f.
//
// Solidity: function getConduit(bytes32 conduitKey) view returns(address conduit, bool exists)
func (_ConduitController *ConduitControllerSession) GetConduit(conduitKey [32]byte) (struct {
	Conduit common.Address
	Exists  bool
}, error) {
	return _ConduitController.Contract.GetConduit(&_ConduitController.CallOpts, conduitKey)
}

// GetConduit is a free data retrieval call binding the contract method 0x6e9bfd9f.
//
// Solidity: function getConduit(bytes32 conduitKey) view returns(address conduit, bool exists)
func (_ConduitController *ConduitControllerCallerSession) GetConduit(conduitKey [32]byte) (struct {
	Conduit common.Address
	Exists  bool
}, error) {
	return _ConduitController.Contract.GetConduit(&_ConduitController.CallOpts, conduitKey)
}

// GetConduitCodeHashes is a free data retrieval call binding the contract method 0x0a96ad39.
//
// Solidity: function getConduitCodeHashes() view returns(bytes32 creationCodeHash, bytes32 runtimeCodeHash)
func (_ConduitController *ConduitControllerCaller) GetConduitCodeHashes(opts *bind.CallOpts) (struct {
	CreationCodeHash [32]byte
	RuntimeCodeHash  [32]byte
}, error) {
	var out []interface{}
	err := _ConduitController.contract.Call(opts, &out, "getConduitCodeHashes")

	outstruct := new(struct {
		CreationCodeHash [32]byte
		RuntimeCodeHash  [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.CreationCodeHash = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.RuntimeCodeHash = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

// GetConduitCodeHashes is a free data retrieval call binding the contract method 0x0a96ad39.
//
// Solidity: function getConduitCodeHashes() view returns(bytes32 creationCodeHash, bytes32 runtimeCodeHash)
func (_ConduitController *ConduitControllerSession) GetConduitCodeHashes() (struct {
	CreationCodeHash [32]byte
	RuntimeCodeHash  [32]byte
}, error) {
	return _ConduitController.Contract.GetConduitCodeHashes(&_ConduitController.CallOpts)
}

// GetConduitCodeHashes is a free data retrieval call binding the contract method 0x0a96ad39.
//
// Solidity: function getConduitCodeHashes() view returns(bytes32 creationCodeHash, bytes32 runtimeCodeHash)
func (_ConduitController *ConduitControllerCallerSession) GetConduitCodeHashes() (struct {
	CreationCodeHash [32]byte
	RuntimeCodeHash  [32]byte
}, error) {
	return _ConduitController.Contract.GetConduitCodeHashes(&_ConduitController.CallOpts)
}

// GetKey is a free data retrieval call binding the contract method 0x93790f44.
//
// Solidity: function getKey(address conduit) view returns(bytes32 conduitKey)
func (_ConduitController *ConduitControllerCaller) GetKey(opts *bind.CallOpts, conduit common.Address) ([32]byte, error) {
	var out []interface{}
	err := _ConduitController.contract.Call(opts, &out, "getKey", conduit)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetKey is a free data retrieval call binding the contract method 0x93790f44.
//
// Solidity: function getKey(address conduit) view returns(bytes32 conduitKey)
func (_ConduitController *ConduitControllerSession) GetKey(conduit common.Address) ([32]byte, error) {
	return _ConduitController.Contract.GetKey(&_ConduitController.CallOpts, conduit)
}

// GetKey is a free data retrieval call binding the contract method 0x93790f44.
//
// Solidity: function getKey(address conduit) view returns(bytes32 conduitKey)
func (_ConduitController *ConduitControllerCallerSession) GetKey(conduit common.Address) ([32]byte, error) {
	return _ConduitController.Contract.GetKey(&_ConduitController.CallOpts, conduit)
}

// GetPotentialOwner is a free data retrieval call binding the contract method 0x906c87cc.
//
// Solidity: function getPotentialOwner(address conduit) view returns(address potentialOwner)
func (_ConduitController *ConduitControllerCaller) GetPotentialOwner(opts *bind.CallOpts, conduit common.Address) (common.Address, error) {
	var out []interface{}
	err := _ConduitController.contract.Call(opts, &out, "getPotentialOwner", conduit)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetPotentialOwner is a free data retrieval call binding the contract method 0x906c87cc.
//
// Solidity: function getPotentialOwner(address conduit) view returns(address potentialOwner)
func (_ConduitController *ConduitControllerSession) GetPotentialOwner(conduit common.Address) (common.Address, error) {
	return _ConduitController.Contract.GetPotentialOwner(&_ConduitController.CallOpts, conduit)
}

// GetPotentialOwner is a free data retrieval call binding the contract method 0x906c87cc.
//
// Solidity: function getPotentialOwner(address conduit) view returns(address potentialOwner)
func (_ConduitController *ConduitControllerCallerSession) GetPotentialOwner(conduit common.Address) (common.Address, error) {
	return _ConduitController.Contract.GetPotentialOwner(&_ConduitController.CallOpts, conduit)
}

// GetTotalChannels is a free data retrieval call binding the contract method 0x4e3f9580.
//
// Solidity: function getTotalChannels(address conduit) view returns(uint256 totalChannels)
func (_ConduitController *ConduitControllerCaller) GetTotalChannels(opts *bind.CallOpts, conduit common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ConduitController.contract.Call(opts, &out, "getTotalChannels", conduit)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTotalChannels is a free data retrieval call binding the contract method 0x4e3f9580.
//
// Solidity: function getTotalChannels(address conduit) view returns(uint256 totalChannels)
func (_ConduitController *ConduitControllerSession) GetTotalChannels(conduit common.Address) (*big.Int, error) {
	return _ConduitController.Contract.GetTotalChannels(&_ConduitController.CallOpts, conduit)
}

// GetTotalChannels is a free data retrieval call binding the contract method 0x4e3f9580.
//
// Solidity: function getTotalChannels(address conduit) view returns(uint256 totalChannels)
func (_ConduitController *ConduitControllerCallerSession) GetTotalChannels(conduit common.Address) (*big.Int, error) {
	return _ConduitController.Contract.GetTotalChannels(&_ConduitController.CallOpts, conduit)
}

// OwnerOf is a free data retrieval call binding the contract method 0x14afd79e.
//
// Solidity: function ownerOf(address conduit) view returns(address owner)
func (_ConduitController *ConduitControllerCaller) OwnerOf(opts *bind.CallOpts, conduit common.Address) (common.Address, error) {
	var out []interface{}
	err := _ConduitController.contract.Call(opts, &out, "ownerOf", conduit)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x14afd79e.
//
// Solidity: function ownerOf(address conduit) view returns(address owner)
func (_ConduitController *ConduitControllerSession) OwnerOf(conduit common.Address) (common.Address, error) {
	return _ConduitController.Contract.OwnerOf(&_ConduitController.CallOpts, conduit)
}

// OwnerOf is a free data retrieval call binding the contract method 0x14afd79e.
//
// Solidity: function ownerOf(address conduit) view returns(address owner)
func (_ConduitController *ConduitControllerCallerSession) OwnerOf(conduit common.Address) (common.Address, error) {
	return _ConduitController.Contract.OwnerOf(&_ConduitController.CallOpts, conduit)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x51710e45.
//
// Solidity: function acceptOwnership(address conduit) returns()
func (_ConduitController *ConduitControllerTransactor) AcceptOwnership(opts *bind.TransactOpts, conduit common.Address) (*types.Transaction, error) {
	return _ConduitController.contract.Transact(opts, "acceptOwnership", conduit)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x51710e45.
//
// Solidity: function acceptOwnership(address conduit) returns()
func (_ConduitController *ConduitControllerSession) AcceptOwnership(conduit common.Address) (*types.Transaction, error) {
	return _ConduitController.Contract.AcceptOwnership(&_ConduitController.TransactOpts, conduit)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x51710e45.
//
// Solidity: function acceptOwnership(address conduit) returns()
func (_ConduitController *ConduitControllerTransactorSession) AcceptOwnership(conduit common.Address) (*types.Transaction, error) {
	return _ConduitController.Contract.AcceptOwnership(&_ConduitController.TransactOpts, conduit)
}

// CancelOwnershipTransfer is a paid mutator transaction binding the contract method 0x7b37e561.
//
// Solidity: function cancelOwnershipTransfer(address conduit) returns()
func (_ConduitController *ConduitControllerTransactor) CancelOwnershipTransfer(opts *bind.TransactOpts, conduit common.Address) (*types.Transaction, error) {
	return _ConduitController.contract.Transact(opts, "cancelOwnershipTransfer", conduit)
}

// CancelOwnershipTransfer is a paid mutator transaction binding the contract method 0x7b37e561.
//
// Solidity: function cancelOwnershipTransfer(address conduit) returns()
func (_ConduitController *ConduitControllerSession) CancelOwnershipTransfer(conduit common.Address) (*types.Transaction, error) {
	return _ConduitController.Contract.CancelOwnershipTransfer(&_ConduitController.TransactOpts, conduit)
}

// CancelOwnershipTransfer is a paid mutator transaction binding the contract method 0x7b37e561.
//
// Solidity: function cancelOwnershipTransfer(address conduit) returns()
func (_ConduitController *ConduitControllerTransactorSession) CancelOwnershipTransfer(conduit common.Address) (*types.Transaction, error) {
	return _ConduitController.Contract.CancelOwnershipTransfer(&_ConduitController.TransactOpts, conduit)
}

// CreateConduit is a paid mutator transaction binding the contract method 0x794593bc.
//
// Solidity: function createConduit(bytes32 conduitKey, address initialOwner) returns(address conduit)
func (_ConduitController *ConduitControllerTransactor) CreateConduit(opts *bind.TransactOpts, conduitKey [32]byte, initialOwner common.Address) (*types.Transaction, error) {
	return _ConduitController.contract.Transact(opts, "createConduit", conduitKey, initialOwner)
}

// CreateConduit is a paid mutator transaction binding the contract method 0x794593bc.
//
// Solidity: function createConduit(bytes32 conduitKey, address initialOwner) returns(address conduit)
func (_ConduitController *ConduitControllerSession) CreateConduit(conduitKey [32]byte, initialOwner common.Address) (*types.Transaction, error) {
	return _ConduitController.Contract.CreateConduit(&_ConduitController.TransactOpts, conduitKey, initialOwner)
}

// CreateConduit is a paid mutator transaction binding the contract method 0x794593bc.
//
// Solidity: function createConduit(bytes32 conduitKey, address initialOwner) returns(address conduit)
func (_ConduitController *ConduitControllerTransactorSession) CreateConduit(conduitKey [32]byte, initialOwner common.Address) (*types.Transaction, error) {
	return _ConduitController.Contract.CreateConduit(&_ConduitController.TransactOpts, conduitKey, initialOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0x6d435421.
//
// Solidity: function transferOwnership(address conduit, address newPotentialOwner) returns()
func (_ConduitController *ConduitControllerTransactor) TransferOwnership(opts *bind.TransactOpts, conduit common.Address, newPotentialOwner common.Address) (*types.Transaction, error) {
	return _ConduitController.contract.Transact(opts, "transferOwnership", conduit, newPotentialOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0x6d435421.
//
// Solidity: function transferOwnership(address conduit, address newPotentialOwner) returns()
func (_ConduitController *ConduitControllerSession) TransferOwnership(conduit common.Address, newPotentialOwner common.Address) (*types.Transaction, error) {
	return _ConduitController.Contract.TransferOwnership(&_ConduitController.TransactOpts, conduit, newPotentialOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0x6d435421.
//
// Solidity: function transferOwnership(address conduit, address newPotentialOwner) returns()
func (_ConduitController *ConduitControllerTransactorSession) TransferOwnership(conduit common.Address, newPotentialOwner common.Address) (*types.Transaction, error) {
	return _ConduitController.Contract.TransferOwnership(&_ConduitController.TransactOpts, conduit, newPotentialOwner)
}

// UpdateChannel is a paid mutator transaction binding the contract method 0x13ad9cab.
//
// Solidity: function updateChannel(address conduit, address channel, bool isOpen) returns()
func (_ConduitController *ConduitControllerTransactor) UpdateChannel(opts *bind.TransactOpts, conduit common.Address, channel common.Address, isOpen bool) (*types.Transaction, error) {
	return _ConduitController.contract.Transact(opts, "updateChannel", conduit, channel, isOpen)
}

// UpdateChannel is a paid mutator transaction binding the contract method 0x13ad9cab.
//
// Solidity: function updateChannel(address conduit, address channel, bool isOpen) returns()
func (_ConduitController *ConduitControllerSession) UpdateChannel(conduit common.Address, channel common.Address, isOpen bool) (*types.Transaction, error) {
	return _ConduitController.Contract.UpdateChannel(&_ConduitController.TransactOpts, conduit, channel, isOpen)
}

// UpdateChannel is a paid mutator transaction binding the contract method 0x13ad9cab.
//
// Solidity: function updateChannel(address conduit, address channel, bool isOpen) returns()
func (_ConduitController *ConduitControllerTransactorSession) UpdateChannel(conduit common.Address, channel common.Address, isOpen bool) (*types.Transaction, error) {
	return _ConduitController.Contract.UpdateChannel(&_ConduitController.TransactOpts, conduit, channel, isOpen)
}

// ConduitControllerNewConduitIterator is returned from FilterNewConduit and is used to iterate over the raw logs and unpacked data for NewConduit events raised by the ConduitController contract.
type ConduitControllerNewConduitIterator struct {
	Event *ConduitControllerNewConduit // Event containing the contract specifics and raw log

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
func (it *ConduitControllerNewConduitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConduitControllerNewConduit)
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
		it.Event = new(ConduitControllerNewConduit)
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
func (it *ConduitControllerNewConduitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConduitControllerNewConduitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConduitControllerNewConduit represents a NewConduit event raised by the ConduitController contract.
type ConduitControllerNewConduit struct {
	Conduit    common.Address
	ConduitKey [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterNewConduit is a free log retrieval operation binding the contract event 0x4397af6128d529b8ae0442f99db1296d5136062597a15bbc61c1b2a6431a7d15.
//
// Solidity: event NewConduit(address conduit, bytes32 conduitKey)
func (_ConduitController *ConduitControllerFilterer) FilterNewConduit(opts *bind.FilterOpts) (*ConduitControllerNewConduitIterator, error) {

	logs, sub, err := _ConduitController.contract.FilterLogs(opts, "NewConduit")
	if err != nil {
		return nil, err
	}
	return &ConduitControllerNewConduitIterator{contract: _ConduitController.contract, event: "NewConduit", logs: logs, sub: sub}, nil
}

// WatchNewConduit is a free log subscription operation binding the contract event 0x4397af6128d529b8ae0442f99db1296d5136062597a15bbc61c1b2a6431a7d15.
//
// Solidity: event NewConduit(address conduit, bytes32 conduitKey)
func (_ConduitController *ConduitControllerFilterer) WatchNewConduit(opts *bind.WatchOpts, sink chan<- *ConduitControllerNewConduit) (event.Subscription, error) {

	logs, sub, err := _ConduitController.contract.WatchLogs(opts, "NewConduit")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConduitControllerNewConduit)
				if err := _ConduitController.contract.UnpackLog(event, "NewConduit", log); err != nil {
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

// ParseNewConduit is a log parse operation binding the contract event 0x4397af6128d529b8ae0442f99db1296d5136062597a15bbc61c1b2a6431a7d15.
//
// Solidity: event NewConduit(address conduit, bytes32 conduitKey)
func (_ConduitController *ConduitControllerFilterer) ParseNewConduit(log types.Log) (*ConduitControllerNewConduit, error) {
	event := new(ConduitControllerNewConduit)
	if err := _ConduitController.contract.UnpackLog(event, "NewConduit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConduitControllerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ConduitController contract.
type ConduitControllerOwnershipTransferredIterator struct {
	Event *ConduitControllerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ConduitControllerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConduitControllerOwnershipTransferred)
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
		it.Event = new(ConduitControllerOwnershipTransferred)
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
func (it *ConduitControllerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConduitControllerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConduitControllerOwnershipTransferred represents a OwnershipTransferred event raised by the ConduitController contract.
type ConduitControllerOwnershipTransferred struct {
	Conduit       common.Address
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0xc8894f26f396ce8c004245c8b7cd1b92103a6e4302fcbab883987149ac01b7ec.
//
// Solidity: event OwnershipTransferred(address indexed conduit, address indexed previousOwner, address indexed newOwner)
func (_ConduitController *ConduitControllerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, conduit []common.Address, previousOwner []common.Address, newOwner []common.Address) (*ConduitControllerOwnershipTransferredIterator, error) {

	var conduitRule []interface{}
	for _, conduitItem := range conduit {
		conduitRule = append(conduitRule, conduitItem)
	}
	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ConduitController.contract.FilterLogs(opts, "OwnershipTransferred", conduitRule, previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ConduitControllerOwnershipTransferredIterator{contract: _ConduitController.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0xc8894f26f396ce8c004245c8b7cd1b92103a6e4302fcbab883987149ac01b7ec.
//
// Solidity: event OwnershipTransferred(address indexed conduit, address indexed previousOwner, address indexed newOwner)
func (_ConduitController *ConduitControllerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ConduitControllerOwnershipTransferred, conduit []common.Address, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var conduitRule []interface{}
	for _, conduitItem := range conduit {
		conduitRule = append(conduitRule, conduitItem)
	}
	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ConduitController.contract.WatchLogs(opts, "OwnershipTransferred", conduitRule, previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConduitControllerOwnershipTransferred)
				if err := _ConduitController.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0xc8894f26f396ce8c004245c8b7cd1b92103a6e4302fcbab883987149ac01b7ec.
//
// Solidity: event OwnershipTransferred(address indexed conduit, address indexed previousOwner, address indexed newOwner)
func (_ConduitController *ConduitControllerFilterer) ParseOwnershipTransferred(log types.Log) (*ConduitControllerOwnershipTransferred, error) {
	event := new(ConduitControllerOwnershipTransferred)
	if err := _ConduitController.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConduitControllerPotentialOwnerUpdatedIterator is returned from FilterPotentialOwnerUpdated and is used to iterate over the raw logs and unpacked data for PotentialOwnerUpdated events raised by the ConduitController contract.
type ConduitControllerPotentialOwnerUpdatedIterator struct {
	Event *ConduitControllerPotentialOwnerUpdated // Event containing the contract specifics and raw log

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
func (it *ConduitControllerPotentialOwnerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConduitControllerPotentialOwnerUpdated)
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
		it.Event = new(ConduitControllerPotentialOwnerUpdated)
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
func (it *ConduitControllerPotentialOwnerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConduitControllerPotentialOwnerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConduitControllerPotentialOwnerUpdated represents a PotentialOwnerUpdated event raised by the ConduitController contract.
type ConduitControllerPotentialOwnerUpdated struct {
	NewPotentialOwner common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterPotentialOwnerUpdated is a free log retrieval operation binding the contract event 0x11a3cf439fb225bfe74225716b6774765670ec1060e3796802e62139d69974da.
//
// Solidity: event PotentialOwnerUpdated(address indexed newPotentialOwner)
func (_ConduitController *ConduitControllerFilterer) FilterPotentialOwnerUpdated(opts *bind.FilterOpts, newPotentialOwner []common.Address) (*ConduitControllerPotentialOwnerUpdatedIterator, error) {

	var newPotentialOwnerRule []interface{}
	for _, newPotentialOwnerItem := range newPotentialOwner {
		newPotentialOwnerRule = append(newPotentialOwnerRule, newPotentialOwnerItem)
	}

	logs, sub, err := _ConduitController.contract.FilterLogs(opts, "PotentialOwnerUpdated", newPotentialOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ConduitControllerPotentialOwnerUpdatedIterator{contract: _ConduitController.contract, event: "PotentialOwnerUpdated", logs: logs, sub: sub}, nil
}

// WatchPotentialOwnerUpdated is a free log subscription operation binding the contract event 0x11a3cf439fb225bfe74225716b6774765670ec1060e3796802e62139d69974da.
//
// Solidity: event PotentialOwnerUpdated(address indexed newPotentialOwner)
func (_ConduitController *ConduitControllerFilterer) WatchPotentialOwnerUpdated(opts *bind.WatchOpts, sink chan<- *ConduitControllerPotentialOwnerUpdated, newPotentialOwner []common.Address) (event.Subscription, error) {

	var newPotentialOwnerRule []interface{}
	for _, newPotentialOwnerItem := range newPotentialOwner {
		newPotentialOwnerRule = append(newPotentialOwnerRule, newPotentialOwnerItem)
	}

	logs, sub, err := _ConduitController.contract.WatchLogs(opts, "PotentialOwnerUpdated", newPotentialOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConduitControllerPotentialOwnerUpdated)
				if err := _ConduitController.contract.UnpackLog(event, "PotentialOwnerUpdated", log); err != nil {
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

// ParsePotentialOwnerUpdated is a log parse operation binding the contract event 0x11a3cf439fb225bfe74225716b6774765670ec1060e3796802e62139d69974da.
//
// Solidity: event PotentialOwnerUpdated(address indexed newPotentialOwner)
func (_ConduitController *ConduitControllerFilterer) ParsePotentialOwnerUpdated(log types.Log) (*ConduitControllerPotentialOwnerUpdated, error) {
	event := new(ConduitControllerPotentialOwnerUpdated)
	if err := _ConduitController.contract.UnpackLog(event, "PotentialOwnerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
