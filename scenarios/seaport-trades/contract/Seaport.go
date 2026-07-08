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

// AdditionalRecipient is an auto generated low-level Go binding around an user-defined struct.
type AdditionalRecipient struct {
	Amount    *big.Int
	Recipient common.Address
}

// AdvancedOrder is an auto generated low-level Go binding around an user-defined struct.
type AdvancedOrder struct {
	Parameters  OrderParameters
	Numerator   *big.Int
	Denominator *big.Int
	Signature   []byte
	ExtraData   []byte
}

// BasicOrderParameters is an auto generated low-level Go binding around an user-defined struct.
type BasicOrderParameters struct {
	ConsiderationToken                common.Address
	ConsiderationIdentifier           *big.Int
	ConsiderationAmount               *big.Int
	Offerer                           common.Address
	Zone                              common.Address
	OfferToken                        common.Address
	OfferIdentifier                   *big.Int
	OfferAmount                       *big.Int
	BasicOrderType                    uint8
	StartTime                         *big.Int
	EndTime                           *big.Int
	ZoneHash                          [32]byte
	Salt                              *big.Int
	OffererConduitKey                 [32]byte
	FulfillerConduitKey               [32]byte
	TotalOriginalAdditionalRecipients *big.Int
	AdditionalRecipients              []AdditionalRecipient
	Signature                         []byte
}

// ConsiderationItem is an auto generated low-level Go binding around an user-defined struct.
type ConsiderationItem struct {
	ItemType             uint8
	Token                common.Address
	IdentifierOrCriteria *big.Int
	StartAmount          *big.Int
	EndAmount            *big.Int
	Recipient            common.Address
}

// CriteriaResolver is an auto generated low-level Go binding around an user-defined struct.
type CriteriaResolver struct {
	OrderIndex    *big.Int
	Side          uint8
	Index         *big.Int
	Identifier    *big.Int
	CriteriaProof [][32]byte
}

// Execution is an auto generated low-level Go binding around an user-defined struct.
type Execution struct {
	Item       ReceivedItem
	Offerer    common.Address
	ConduitKey [32]byte
}

// Fulfillment is an auto generated low-level Go binding around an user-defined struct.
type Fulfillment struct {
	OfferComponents         []FulfillmentComponent
	ConsiderationComponents []FulfillmentComponent
}

// FulfillmentComponent is an auto generated low-level Go binding around an user-defined struct.
type FulfillmentComponent struct {
	OrderIndex *big.Int
	ItemIndex  *big.Int
}

// OfferItem is an auto generated low-level Go binding around an user-defined struct.
type OfferItem struct {
	ItemType             uint8
	Token                common.Address
	IdentifierOrCriteria *big.Int
	StartAmount          *big.Int
	EndAmount            *big.Int
}

// Order is an auto generated low-level Go binding around an user-defined struct.
type Order struct {
	Parameters OrderParameters
	Signature  []byte
}

// OrderComponents is an auto generated low-level Go binding around an user-defined struct.
type OrderComponents struct {
	Offerer       common.Address
	Zone          common.Address
	Offer         []OfferItem
	Consideration []ConsiderationItem
	OrderType     uint8
	StartTime     *big.Int
	EndTime       *big.Int
	ZoneHash      [32]byte
	Salt          *big.Int
	ConduitKey    [32]byte
	Counter       *big.Int
}

// OrderParameters is an auto generated low-level Go binding around an user-defined struct.
type OrderParameters struct {
	Offerer                         common.Address
	Zone                            common.Address
	Offer                           []OfferItem
	Consideration                   []ConsiderationItem
	OrderType                       uint8
	StartTime                       *big.Int
	EndTime                         *big.Int
	ZoneHash                        [32]byte
	Salt                            *big.Int
	ConduitKey                      [32]byte
	TotalOriginalConsiderationItems *big.Int
}

// ReceivedItem is an auto generated low-level Go binding around an user-defined struct.
type ReceivedItem struct {
	ItemType   uint8
	Token      common.Address
	Identifier *big.Int
	Amount     *big.Int
	Recipient  common.Address
}

// SpentItem is an auto generated low-level Go binding around an user-defined struct.
type SpentItem struct {
	ItemType   uint8
	Token      common.Address
	Identifier *big.Int
	Amount     *big.Int
}

// SeaportMetaData contains all meta data concerning the Seaport contract.
var SeaportMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"conduitController\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"BadContractSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BadFraction\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BadReturnValueFromERC20OnTransfer\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"}],\"name\":\"BadSignatureV\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotCancelOrder\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConsiderationCriteriaResolverOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConsiderationLengthNotEqualToTotalOriginal\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"orderIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"considerationIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"shortfallAmount\",\"type\":\"uint256\"}],\"name\":\"ConsiderationNotMet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CriteriaNotEnabledForItem\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"identifiers\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"name\":\"ERC1155BatchTransferGenericFailure\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InexactFraction\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientNativeTokensSupplied\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Invalid1155BatchTransferEncoding\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidBasicOrderParameterEncoding\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"conduit\",\"type\":\"address\"}],\"name\":\"InvalidCallToConduit\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"conduitKey\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"conduit\",\"type\":\"address\"}],\"name\":\"InvalidConduit\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"}],\"name\":\"InvalidContractOrder\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"InvalidERC721TransferAmount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidFulfillmentComponentData\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"InvalidMsgValue\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidNativeOfferItem\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidProof\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"}],\"name\":\"InvalidRestrictedOrder\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSigner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endTime\",\"type\":\"uint256\"}],\"name\":\"InvalidTime\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"fulfillmentIndex\",\"type\":\"uint256\"}],\"name\":\"MismatchedFulfillmentOfferAndConsiderationComponents\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"enumSide\",\"name\":\"side\",\"type\":\"uint8\"}],\"name\":\"MissingFulfillmentComponentOnAggregation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MissingItemAmount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MissingOriginalConsiderationItems\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"NativeTokenTransferGenericFailure\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"NoContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoReentrantCalls\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoSpecifiedOrdersAvailable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OfferAndConsiderationRequiredOnFulfillment\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OfferCriteriaResolverOutOfRange\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"}],\"name\":\"OrderAlreadyFilled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"enumSide\",\"name\":\"side\",\"type\":\"uint8\"}],\"name\":\"OrderCriteriaResolverOutOfRange\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"}],\"name\":\"OrderIsCancelled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"}],\"name\":\"OrderPartiallyFilled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PartialFillsNotEnabledForOrder\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TStoreAlreadyActivated\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TStoreNotSupported\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TloadTestContractDeploymentFailed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"identifier\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TokenTransferGenericFailure\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"orderIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"considerationIndex\",\"type\":\"uint256\"}],\"name\":\"UnresolvedConsiderationCriteria\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"orderIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"offerIndex\",\"type\":\"uint256\"}],\"name\":\"UnresolvedOfferCriteria\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnusedItemParameters\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newCounter\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"offerer\",\"type\":\"address\"}],\"name\":\"CounterIncremented\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"offerer\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"zone\",\"type\":\"address\"}],\"name\":\"OrderCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"offerer\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"zone\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumItemType\",\"name\":\"itemType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"identifier\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structSpentItem[]\",\"name\":\"offer\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"enumItemType\",\"name\":\"itemType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"identifier\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structReceivedItem[]\",\"name\":\"consideration\",\"type\":\"tuple[]\"}],\"name\":\"OrderFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"offerer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"zone\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumItemType\",\"name\":\"itemType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"identifierOrCriteria\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endAmount\",\"type\":\"uint256\"}],\"internalType\":\"structOfferItem[]\",\"name\":\"offer\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"enumItemType\",\"name\":\"itemType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"identifierOrCriteria\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endAmount\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"}],\"internalType\":\"structConsiderationItem[]\",\"name\":\"consideration\",\"type\":\"tuple[]\"},{\"internalType\":\"enumOrderType\",\"name\":\"orderType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endTime\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"zoneHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"conduitKey\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"totalOriginalConsiderationItems\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structOrderParameters\",\"name\":\"orderParameters\",\"type\":\"tuple\"}],\"name\":\"OrderValidated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32[]\",\"name\":\"orderHashes\",\"type\":\"bytes32[]\"}],\"name\":\"OrdersMatched\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"__activateTstore\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"offerer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"zone\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumItemType\",\"name\":\"itemType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"identifierOrCriteria\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endAmount\",\"type\":\"uint256\"}],\"internalType\":\"structOfferItem[]\",\"name\":\"offer\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"enumItemType\",\"name\":\"itemType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"identifierOrCriteria\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endAmount\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"}],\"internalType\":\"structConsiderationItem[]\",\"name\":\"consideration\",\"type\":\"tuple[]\"},{\"internalType\":\"enumOrderType\",\"name\":\"orderType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endTime\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"zoneHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"conduitKey\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"counter\",\"type\":\"uint256\"}],\"internalType\":\"structOrderComponents[]\",\"name\":\"orders\",\"type\":\"tuple[]\"}],\"name\":\"cancel\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"cancelled\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"offerer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"zone\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumItemType\",\"name\":\"itemType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"identifierOrCriteria\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endAmount\",\"type\":\"uint256\"}],\"internalType\":\"structOfferItem[]\",\"name\":\"offer\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"enumItemType\",\"name\":\"itemType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"identifierOrCriteria\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endAmount\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"}],\"internalType\":\"structConsiderationItem[]\",\"name\":\"consideration\",\"type\":\"tuple[]\"},{\"internalType\":\"enumOrderType\",\"name\":\"orderType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endTime\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"zoneHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"conduitKey\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"totalOriginalConsiderationItems\",\"type\":\"uint256\"}],\"internalType\":\"structOrderParameters\",\"name\":\"parameters\",\"type\":\"tuple\"},{\"internalType\":\"uint120\",\"name\":\"numerator\",\"type\":\"uint120\"},{\"internalType\":\"uint120\",\"name\":\"denominator\",\"type\":\"uint120\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"internalType\":\"structAdvancedOrder\",\"name\":\"\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"orderIndex\",\"type\":\"uint256\"},{\"internalType\":\"enumSide\",\"name\":\"side\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"identifier\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"criteriaProof\",\"type\":\"bytes32[]\"}],\"internalType\":\"structCriteriaResolver[]\",\"name\":\"\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes32\",\"name\":\"fulfillerConduitKey\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"fulfillAdvancedOrder\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"fulfilled\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"offerer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"zone\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumItemType\",\"name\":\"itemType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"identifierOrCriteria\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endAmount\",\"type\":\"uint256\"}],\"internalType\":\"structOfferItem[]\",\"name\":\"offer\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"enumItemType\",\"name\":\"itemType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"identifierOrCriteria\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endAmount\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"}],\"internalType\":\"structConsiderationItem[]\",\"name\":\"consideration\",\"type\":\"tuple[]\"},{\"internalType\":\"enumOrderType\",\"name\":\"orderType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endTime\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"zoneHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"conduitKey\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"totalOriginalConsiderationItems\",\"type\":\"uint256\"}],\"internalType\":\"structOrderParameters\",\"name\":\"parameters\",\"type\":\"tuple\"},{\"internalType\":\"uint120\",\"name\":\"numerator\",\"type\":\"uint120\"},{\"internalType\":\"uint120\",\"name\":\"denominator\",\"type\":\"uint120\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"internalType\":\"structAdvancedOrder[]\",\"name\":\"\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"orderIndex\",\"type\":\"uint256\"},{\"internalType\":\"enumSide\",\"name\":\"side\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"identifier\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"criteriaProof\",\"type\":\"bytes32[]\"}],\"internalType\":\"structCriteriaResolver[]\",\"name\":\"\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"orderIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"itemIndex\",\"type\":\"uint256\"}],\"internalType\":\"structFulfillmentComponent[][]\",\"name\":\"\",\"type\":\"tuple[][]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"orderIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"itemIndex\",\"type\":\"uint256\"}],\"internalType\":\"structFulfillmentComponent[][]\",\"name\":\"\",\"type\":\"tuple[][]\"},{\"internalType\":\"bytes32\",\"name\":\"fulfillerConduitKey\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"maximumFulfilled\",\"type\":\"uint256\"}],\"name\":\"fulfillAvailableAdvancedOrders\",\"outputs\":[{\"internalType\":\"bool[]\",\"name\":\"\",\"type\":\"bool[]\"},{\"components\":[{\"components\":[{\"internalType\":\"enumItemType\",\"name\":\"itemType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"identifier\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"}],\"internalType\":\"structReceivedItem\",\"name\":\"item\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"offerer\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"conduitKey\",\"type\":\"bytes32\"}],\"internalType\":\"structExecution[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"offerer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"zone\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumItemType\",\"name\":\"itemType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"identifierOrCriteria\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endAmount\",\"type\":\"uint256\"}],\"internalType\":\"structOfferItem[]\",\"name\":\"offer\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"enumItemType\",\"name\":\"itemType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"identifierOrCriteria\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endAmount\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"}],\"internalType\":\"structConsiderationItem[]\",\"name\":\"consideration\",\"type\":\"tuple[]\"},{\"internalType\":\"enumOrderType\",\"name\":\"orderType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endTime\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"zoneHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"conduitKey\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"totalOriginalConsiderationItems\",\"type\":\"uint256\"}],\"internalType\":\"structOrderParameters\",\"name\":\"parameters\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structOrder[]\",\"name\":\"\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"orderIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"itemIndex\",\"type\":\"uint256\"}],\"internalType\":\"structFulfillmentComponent[][]\",\"name\":\"\",\"type\":\"tuple[][]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"orderIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"itemIndex\",\"type\":\"uint256\"}],\"internalType\":\"structFulfillmentComponent[][]\",\"name\":\"\",\"type\":\"tuple[][]\"},{\"internalType\":\"bytes32\",\"name\":\"fulfillerConduitKey\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"maximumFulfilled\",\"type\":\"uint256\"}],\"name\":\"fulfillAvailableOrders\",\"outputs\":[{\"internalType\":\"bool[]\",\"name\":\"\",\"type\":\"bool[]\"},{\"components\":[{\"components\":[{\"internalType\":\"enumItemType\",\"name\":\"itemType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"identifier\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"}],\"internalType\":\"structReceivedItem\",\"name\":\"item\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"offerer\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"conduitKey\",\"type\":\"bytes32\"}],\"internalType\":\"structExecution[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"considerationToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"considerationIdentifier\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"considerationAmount\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"offerer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"zone\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"offerToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"offerIdentifier\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"offerAmount\",\"type\":\"uint256\"},{\"internalType\":\"enumBasicOrderType\",\"name\":\"basicOrderType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endTime\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"zoneHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"offererConduitKey\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"fulfillerConduitKey\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"totalOriginalAdditionalRecipients\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"}],\"internalType\":\"structAdditionalRecipient[]\",\"name\":\"additionalRecipients\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structBasicOrderParameters\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"fulfillBasicOrder\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"fulfilled\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"considerationToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"considerationIdentifier\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"considerationAmount\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"offerer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"zone\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"offerToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"offerIdentifier\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"offerAmount\",\"type\":\"uint256\"},{\"internalType\":\"enumBasicOrderType\",\"name\":\"basicOrderType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endTime\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"zoneHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"offererConduitKey\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"fulfillerConduitKey\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"totalOriginalAdditionalRecipients\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"}],\"internalType\":\"structAdditionalRecipient[]\",\"name\":\"additionalRecipients\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structBasicOrderParameters\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"fulfillBasicOrder_efficient_6GL6yc\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"fulfilled\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"offerer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"zone\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumItemType\",\"name\":\"itemType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"identifierOrCriteria\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endAmount\",\"type\":\"uint256\"}],\"internalType\":\"structOfferItem[]\",\"name\":\"offer\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"enumItemType\",\"name\":\"itemType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"identifierOrCriteria\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endAmount\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"}],\"internalType\":\"structConsiderationItem[]\",\"name\":\"consideration\",\"type\":\"tuple[]\"},{\"internalType\":\"enumOrderType\",\"name\":\"orderType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endTime\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"zoneHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"conduitKey\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"totalOriginalConsiderationItems\",\"type\":\"uint256\"}],\"internalType\":\"structOrderParameters\",\"name\":\"parameters\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structOrder\",\"name\":\"\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"fulfillerConduitKey\",\"type\":\"bytes32\"}],\"name\":\"fulfillOrder\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"fulfilled\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"contractOfferer\",\"type\":\"address\"}],\"name\":\"getContractOffererNonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"offerer\",\"type\":\"address\"}],\"name\":\"getCounter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"counter\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"offerer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"zone\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumItemType\",\"name\":\"itemType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"identifierOrCriteria\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endAmount\",\"type\":\"uint256\"}],\"internalType\":\"structOfferItem[]\",\"name\":\"offer\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"enumItemType\",\"name\":\"itemType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"identifierOrCriteria\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endAmount\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"}],\"internalType\":\"structConsiderationItem[]\",\"name\":\"consideration\",\"type\":\"tuple[]\"},{\"internalType\":\"enumOrderType\",\"name\":\"orderType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endTime\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"zoneHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"conduitKey\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"counter\",\"type\":\"uint256\"}],\"internalType\":\"structOrderComponents\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"getOrderHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"}],\"name\":\"getOrderStatus\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isValidated\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"isCancelled\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"totalFilled\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalSize\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"incrementCounter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"newCounter\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"information\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"bytes32\",\"name\":\"domainSeparator\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"conduitController\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"offerer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"zone\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumItemType\",\"name\":\"itemType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"identifierOrCriteria\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endAmount\",\"type\":\"uint256\"}],\"internalType\":\"structOfferItem[]\",\"name\":\"offer\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"enumItemType\",\"name\":\"itemType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"identifierOrCriteria\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endAmount\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"}],\"internalType\":\"structConsiderationItem[]\",\"name\":\"consideration\",\"type\":\"tuple[]\"},{\"internalType\":\"enumOrderType\",\"name\":\"orderType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endTime\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"zoneHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"conduitKey\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"totalOriginalConsiderationItems\",\"type\":\"uint256\"}],\"internalType\":\"structOrderParameters\",\"name\":\"parameters\",\"type\":\"tuple\"},{\"internalType\":\"uint120\",\"name\":\"numerator\",\"type\":\"uint120\"},{\"internalType\":\"uint120\",\"name\":\"denominator\",\"type\":\"uint120\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"internalType\":\"structAdvancedOrder[]\",\"name\":\"\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"orderIndex\",\"type\":\"uint256\"},{\"internalType\":\"enumSide\",\"name\":\"side\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"identifier\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"criteriaProof\",\"type\":\"bytes32[]\"}],\"internalType\":\"structCriteriaResolver[]\",\"name\":\"\",\"type\":\"tuple[]\"},{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"orderIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"itemIndex\",\"type\":\"uint256\"}],\"internalType\":\"structFulfillmentComponent[]\",\"name\":\"offerComponents\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"orderIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"itemIndex\",\"type\":\"uint256\"}],\"internalType\":\"structFulfillmentComponent[]\",\"name\":\"considerationComponents\",\"type\":\"tuple[]\"}],\"internalType\":\"structFulfillment[]\",\"name\":\"\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"matchAdvancedOrders\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"enumItemType\",\"name\":\"itemType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"identifier\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"}],\"internalType\":\"structReceivedItem\",\"name\":\"item\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"offerer\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"conduitKey\",\"type\":\"bytes32\"}],\"internalType\":\"structExecution[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"offerer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"zone\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumItemType\",\"name\":\"itemType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"identifierOrCriteria\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endAmount\",\"type\":\"uint256\"}],\"internalType\":\"structOfferItem[]\",\"name\":\"offer\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"enumItemType\",\"name\":\"itemType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"identifierOrCriteria\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endAmount\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"}],\"internalType\":\"structConsiderationItem[]\",\"name\":\"consideration\",\"type\":\"tuple[]\"},{\"internalType\":\"enumOrderType\",\"name\":\"orderType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endTime\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"zoneHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"conduitKey\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"totalOriginalConsiderationItems\",\"type\":\"uint256\"}],\"internalType\":\"structOrderParameters\",\"name\":\"parameters\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structOrder[]\",\"name\":\"\",\"type\":\"tuple[]\"},{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"orderIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"itemIndex\",\"type\":\"uint256\"}],\"internalType\":\"structFulfillmentComponent[]\",\"name\":\"offerComponents\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"orderIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"itemIndex\",\"type\":\"uint256\"}],\"internalType\":\"structFulfillmentComponent[]\",\"name\":\"considerationComponents\",\"type\":\"tuple[]\"}],\"internalType\":\"structFulfillment[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"name\":\"matchOrders\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"enumItemType\",\"name\":\"itemType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"identifier\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"}],\"internalType\":\"structReceivedItem\",\"name\":\"item\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"offerer\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"conduitKey\",\"type\":\"bytes32\"}],\"internalType\":\"structExecution[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"offerer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"zone\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumItemType\",\"name\":\"itemType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"identifierOrCriteria\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endAmount\",\"type\":\"uint256\"}],\"internalType\":\"structOfferItem[]\",\"name\":\"offer\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"enumItemType\",\"name\":\"itemType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"identifierOrCriteria\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endAmount\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"}],\"internalType\":\"structConsiderationItem[]\",\"name\":\"consideration\",\"type\":\"tuple[]\"},{\"internalType\":\"enumOrderType\",\"name\":\"orderType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endTime\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"zoneHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"conduitKey\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"totalOriginalConsiderationItems\",\"type\":\"uint256\"}],\"internalType\":\"structOrderParameters\",\"name\":\"parameters\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structOrder[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"name\":\"validate\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x61020060405234801562000011575f80fd5b50604051620061dc380380620061dc83398101604081905262000034916200031b565b808080808080808080806200004862000179565b610120526101005260e05260c081905260a082815246610140819052604080515f9485526020879052948152606091825230608090815292842085825293909152939052610160526001600160a01b038316610180819052630a96ad3960e01b825282519092630a96ad3992600480820193918290030181865afa158015620000d3573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190620000f991906200034a565b506101a052505f90506200010c620002a0565b90506001600160a01b0381166200013657604051632aea588760e01b815260040160405180910390fd5b5f6200014282620002b9565b8015156101c0526001600160a01b0383166101e0529050806200016857600163929eee14555b5050505050505050505050620003e8565b5f8080808080620001a460408051808201909152600781526614d9585c1bdc9d60ca1b602082015290565b8051906020012095506040518060400160405280600381526020016218971b60e91b8152508051906020012094505f6040518060a00160405280606a815260200162006172606a913990505f6040518060c001604052806084815260200162005fc86084913990505f60405180610100016040528060d481526020016200609e60d4913990506040518060800160405280605281526020016200604c6052913980519060200120965082805190602001209550818051906020012094505f81838560405160200162000279939291906200039c565b60405160208183030381529060405290508080519060200120945050505050909192939495565b5f696002601e613d5c3d52f35f52600a60165ff0905090565b5f816001600160a01b0316600a5a620002d39190620003c8565b6040515f8181818686fa925050503d805f81146200030d576040519150601f19603f3d011682016040523d82523d5f602084013e62000312565b606091505b50909392505050565b5f602082840312156200032c575f80fd5b81516001600160a01b038116811462000343575f80fd5b9392505050565b5f80604083850312156200035c575f80fd5b505080516020909101519092909150565b5f81515f5b818110156200038e576020818501810151868301520162000372565b505f93019283525090919050565b5f620003bf620003b8620003b184886200036d565b866200036d565b846200036d565b95945050505050565b5f82620003e357634e487b7160e01b5f52601260045260245ffd5b500490565b60805160a05160c05160e05161010051610120516101405161016051610180516101a0516101c0516101e051615b24620004a45f395f61047f01525f818161036d0152818161043401528181611a7201528181611ac501526124e501525f61309501525f81816112b7015261306501525f612f2d01525f612e7501525f8181610c27015261165201525f8181610bb6015261149b01525f8181610b5001526115e401525f612ea501525f612eee01525f612eca0152615b245ff3fe608060405260043610610103575f3560e01c8063a900866b11610092578063f07ec37311610062578063f07ec373146102f7578063f2d12b1214610316578063f47b774014610329578063fb0f3ee114610116578063fd9f1e101461034c575f80fd5b8063a900866b1461028a578063b3a34c4c146102be578063e7acab24146102d1578063ed98a574146102e4575f80fd5b80637423eb3c116100d85780637423eb3c146101f757806379df72bd1461020b57806387201b411461022a578063881477321461024b578063a81744041461026a575f80fd5b801561011657806306fdde031461013e57806346423aa71461015f5780635b34b966146101d5575f80fd5b366101125761011061036b565b005b5f80fd5b61012961012436600461513d565b610402565b60405190151581526020015b60405180910390f35b348015610149575f80fd5b50610152610411565b60405161013591906151b7565b34801561016a575f80fd5b506101b36101793660046151c9565b5f9081526001602052604090205460ff808216926101008304909116916001600160781b03620100008204811692600160881b9092041690565b6040805194151585529215156020850152918301526060820152608001610135565b3480156101e0575f80fd5b506101e9610420565b604051908152602001610135565b348015610202575f80fd5b50610110610429565b348015610216575f80fd5b506101e96102253660046151e0565b6104ca565b61023d610238366004615274565b610501565b604051610135929190615418565b348015610256575f80fd5b50610129610265366004615467565b61057a565b61027d6102783660046154a5565b610598565b604051610135919061550b565b348015610295575f80fd5b506101e96102a436600461551d565b6001600160a01b03165f9081526002602052604090205490565b6101296102cc366004615536565b610633565b6101296102df36600461557b565b6106ac565b61023d6102f23660046155ff565b6106ea565b348015610302575f80fd5b506101e961031136600461551d565b61078d565b61027d61032436600461569f565b6107aa565b348015610334575f80fd5b5061033d6107f3565b60405161013593929190615741565b348015610357575f80fd5b50610129610366366004615467565b61080a565b7f000000000000000000000000000000000000000000000000000000000000000080156103b057600263929eee145c146103b05763a61be9f05f52346020526024601cfd5b806103ff5763929eee1454806103de57600263929eee145c146103de5763a61be9f05f52346020526024601cfd5b600381141581151516156103fd5763a61be9f05f52346020526024601cfd5b505b50565b5f61040b610815565b92915050565b606061041b6109c9565b905090565b5f61041b6109e1565b63929eee14546001147f00000000000000000000000000000000000000000000000000000000000000008061045c575080155b1561047a57604051630f45b98b60e41b815260040160405180910390fd5b6104a37f0000000000000000000000000000000000000000000000000000000000000000610a4e565b6104c0576040516370a4078f60e01b815260040160405180910390fd5b5f63929eee145550565b5f806104d66004610aac565b90506104fa6104eb82610abc5b63ffffffff16565b610140830135610b39565b3590565b9392505050565b60608061056661051c6105146004610aac565b610c7d6104e3565b61053361052b60046020610ce5565b610d036104e3565b61054a61054260046040610ce5565b610d5b6104e3565b61055961054260046060610ce5565b89338a15028a0189610db3565b915091509b509b9950505050505050505050565b5f6104fa61059361058b6004610aac565b610ded6104e3565b610e45565b60606106286105b26105aa6004610aac565b610f666104e3565b604080515f808252602082019092529061060a565b6105f76040805160a081019091525f808252602082019081526020015f81526020015f8152602001606081525090565b8152602001906001900390816105c75790505b5061062261061a60046020610ce5565b610fbe6104e3565b33611016565b90505b949350505050565b5f6104fa61064c6106446004610aac565b6110546104e3565b604080515f80825260208201909252906106a4565b6106916040805160a081019091525f808252602082019081526020015f81526020015f8152602001606081525090565b8152602001906001900390816106615790505b5084336110ca565b5f6106e06106c56106bd6004610aac565b6112496104e3565b6106d461052b60046020610ce5565b853386150286016110ca565b9695505050505050565b60608061077c6106fd6105aa6004610aac565b604080515f8082526020820190925290610755565b6107426040805160a081019091525f808252602082019081526020015f81526020015f8152602001606081525090565b8152602001906001900390816107125790505b5061076561054260046020610ce5565b61077461054260046040610ce5565b883389610db3565b915091509850989650505050505050565b6001600160a01b0381165f9081526020819052604081205461040b565b60606107e56107bc6105146004610aac565b6107cb61052b60046020610ce5565b6107da61061a60046040610ce5565b338615028601611016565b90505b979650505050505050565b60605f806107ff6112a5565b925092509250909192565b5f6104fa83836112f1565b5f61012435600281901c906003166001821183341582148061083a5761083a3461142c565b506003841160a0810260240135906502030203010160d01b861a905f630101020360d01b881a61086d888289888861143d565b9096509150506101c4600583901b01355f87600581111561089057610890615351565b036108c55760443560243517156108ae57636ab37ce75f526004601cfd5b6108b8848261175a565b6108c0611824565b6109a6565b6040805160208082528183019092525f9160208201818036833701905050905060028a60058111156108f9576108f9615351565b0361091c5761091760c4356084353360e435610104355b87876118a4565b610991565b60038a600581111561093057610930615351565b0361094e5761091760c4356084353360e435610104355b87876118ef565b60048a600581111561096257610962615351565b0361097c5761091760243533608435604435606435610910565b61099160243533608435604435606435610947565b61099b8482611925565b6109a4816119df565b505b6109b1868984611a03565b6109b9611a70565b6001995050505050505050505090565b6060602080526707536561706f727460475260606020f35b5f6109ea611ac3565b600143034060801c335f525f60205260405f208054820192508281555050336001600160a01b03167f721c20121297512b72821b97f5326877ea8ecf4bb9948fea5bfcb6453074d37f82604051610a4391815260200190565b60405180910390a290565b5f816001600160a01b0316600a5a610a66919061579b565b6040515f8181818686fa925050503d805f8114610a9e576040519150601f19603f3d011682016040523d82523d5f602084013e610aa3565b606091505b50909392505050565b5f813563ffffffff16820161040b565b5f610acf61016060408051918201905290565b9050610ade8282610140611b40565b610afb610af4610aef846040610ce5565b611b49565b6040830152565b5f610b0f610b0a846060610ce5565b611b9c565b9050610b1c816060840152565b610b2f610b27825190565b610140840152565b50919050565b0190565b610140820151604080519084015180515f939284927f000000000000000000000000000000000000000000000000000000000000000092602090910190845b81811015610ba5578251601f1901805186825260c082208652905260209384019390920191600101610b78565b508060051b6040512094505050505f7f0000000000000000000000000000000000000000000000000000000000000000915060405160206060890151015f5b86811015610c11578151601f1901805186825260e082208552905260209283019290910190600101610be4565b505060408051600587901b9020601f198a0180517f00000000000000000000000000000000000000000000000000000000000000008252928b01805197815260608c018051938152610140909c019a8b5261018082209390915295909552939097525050925250919050565b5f8063ffffffff8335169050600581901b610ca16020820160408051918201905290565b828152925060208381019085015f5b83811015610cdb57610cd3610ccd610cc88484610ce5565b611249565b82850152565b602001610cb0565b5050505050919050565b5f6104fa63ffffffff610cfd6104f68686610b358516565b16840190565b5f8063ffffffff8335169050600581901b610d276020820160408051918201905290565b828152925060208381019085015f5b83811015610cdb57610d53610ccd610d4e8484610ce5565b611bdf565b602001610d36565b5f8063ffffffff8335169050600581901b610d7f6020820160408051918201905290565b828152925060208381019085015f5b83811015610cdb57610dab610ccd610da68484610ce5565b611c21565b602001610d8e565b60608036155f80610dc78c8c85898b611c64565b91509150610dda8c8b8b8b8b8787612041565b9450945050505097509795505050505050565b5f8063ffffffff8335169050600581901b610e116020820160408051918201905290565b828152925060208381019085015f5b83811015610cdb57610e3d610ccd610e388484610ce5565b612173565b602001610e20565b5f610e4e611ac3565b5f805f80855190505f5b81811015610f59575f878281518110610e7357610e736157ba565b60209081029190910101518051909150600481608001516004811115610e9b57610e9b615351565b03610ea7575050610f51565b80519450610eb4816121b6565b5f8181526001602052604081209850909650610ed690879089903615156121ef565b50865460ff16610f4e5780610140015181606001515114610ef957610ef961227d565b610f088587846020015161228a565b865460ff191660011787556040517ff280791efe782edcf06ce15c8f4dff17601db3b88eb3805a0db7d77faf757f0490610f4590889084906158bd565b60405180910390a15b50505b600101610e58565b5060019695505050505050565b5f8063ffffffff8335169050600581901b610f8a6020820160408051918201905290565b828152925060208381019085015f5b83811015610cdb57610fb6610ccd610fb18484610ce5565b611054565b602001610f99565b5f8063ffffffff8335169050600581901b610fe26020820160408051918201905290565b828152925060208381019085015f5b83811015610cdb5761100e610ccd6110098484610ce5565b612326565b602001610ff1565b60605f36151590505f8061102e8888858b5189611c64565b9150915061103b82612359565b6110488887848885612398565b98975050505050505050565b5f61106761020060408051918201905290565b60a0810180825290915061108361107d84610aac565b8261246b565b61108f60016020840152565b61109b60016040840152565b6110b86110b16110ac856020610ce5565b6124a1565b6060840152565b610b2f6110c36124c9565b6080840152565b835160808101515f91906110f160048260048111156110eb576110eb615351565b146124e3565b5f80806111008a361515612586565b60408051600180825281830190925293965091945092505f9190816020015b61112761504f565b81526020019060019003908161111f5790505090508a815f8151811061114f5761114f6157ba565b6020026020010181905250611164818b61277e565b6111708684848b6128f7565b6040805160018082528183019092525f9160208083019080368337509192505050361515600487818111156111a7576111a7615351565b146111ca576111b88d83885f6129f7565b6111c486868684612a49565b506111dc565b6111d9888e6080015183612b4d565b95505b6111e7888c8c612c1d565b85825f815181106111fa576111fa6157ba565b6020026020010181815250506112118d8388612d40565b61122e86895f01518a602001518d8c604001518d60600151612e0d565b611236611a70565b5060019c9b505050505050505050505050565b5f61125c61020060408051918201905290565b905061127060208381019083016040611b40565b60a0810180825261128361107d84610aac565b6112946110b16110ac856060610ce5565b610b2f6110c36110ac856080610ce5565b60605f805f6112b2612e72565b90505f7f0000000000000000000000000000000000000000000000000000000000000000905060605f5281602052806040526303312e3660635260a05ff35b5f6112fa611ac3565b5f8083815b81811015611411573687878381811061131a5761131a6157ba565b905060200281019061132c919061599f565b90505f61133c602083018361551d565b90505f61134f604084016020850161551d565b90505f61136260a08501608086016159be565b905081331483331417156004821417871796505f61139661138b6113838790565b610abc6104e3565b866101400135610b39565b5f8181526001602052604090819020805461ffff19166101001781559051909a509091506001600160a01b0380851691908616907f6bacc01dbe442496068f7d234edd811f1a5f833243e0aec824f86ab861f3c90d906113f99085815260200190565b60405180910390a385600101955050505050506112ff565b5050801561142157611421612f4f565b506001949350505050565b63a61be9f05f52806020526024601cfd5b5f806114485f6124e3565b611450612f5c565b42610164351115426101443511171561147e576321ccfeb75f5261014435602052610164356040526044601cfd5b610204356102643510156114995763466aa6165f526004601cfd5b7f0000000000000000000000000000000000000000000000000000000000000000608081905260a08790526060602460c037604060646101203760e060802061016052610264356102043560051b6102a0016001820181526020810190508881526080602460208301376101608860a0528760c0525f60e0525f6102043593505f5b8481101561156f578060400261028401602081610100376040816101203760208101358317925060208401935060e0608020845260a0850194508b85528a602086015260408160608701375060010161151b565b6001850160051b610160206060526102643594505b848110156115bf578060400261028401925060a0840193508a8452896020850152604083606086013760208301359190911790600101611584565b506001600160a01b038111156115dc576339f3e3fd5f526004601cfd5b50505050505f7f00000000000000000000000000000000000000000000000000000000000000009050806080528360a052606060c460c0376020610104610120375060c06080205f9081526020812060e05260843590611650826001600160a01b03165f9081526020819052604090205490565b7f000000000000000000000000000000000000000000000000000000000000000060808190529091506040608460a03760605161010052896101205260a061014461014037816101e05261018060802094505050506102043560051b61018001828152336020820152608060408201526101206060820152600160808201528360a0820152606060c460c083013760a061026435026101e00160a4356084357f9d9af8e38d66c62e2c12f0225249fd9d721c54b83f48d9352c97c6cacdcb6f318385a35f60605260608101820160405250505f61172c83612fa7565b90506117388389612ff8565b7101000000000000000000000000000001000182559150509550959350505050565b60c43560843560e4356101043584156117cb576117768161304d565b5f6040519050632671a55160e11b815260206004820152600160248201528660448201528460648201528360848201523360a48201528260c48201528160e48201526117c5868261010461305f565b5061181c565b60028660058111156117df576117df615351565b0361180657806001146117f5576117f581613116565b61180184843385613127565b61181c565b61180f8161304d565b61181c84843385856131dc565b505050505050565b346064356084356102643560061b5f80805b838110156118755761028481013592506102a481013591508683111561185e5761185e6132b3565b828703965061186d82846132c0565b604001611836565b5085851115611886576118866132b3565b61189084866132c0565b8486111561181c5761181c338688036132c0565b6118ae81836132f6565b816118d557826001146118c4576118c483613116565b6118d087878787613127565b6118e6565b6118e6828260028a8a8a8a8a613314565b50505050505050565b6118f88361304d565b61190281836132f6565b81611914576118d087878787876131dc565b6118e6828260038a8a8a8a8a613314565b5f805f805f861561194a57505060843592503391505060c4356101043560e43561195f565b50339350506084359150506024356064356044355b801561196d5761196d613393565b50600586901b6101e403356102643560061b5f80805b838110156119c45761028481013592506102a481013591508a156119ae576119ab83876159dc565b95505b6119bc878a8486898f6133a0565b604001611983565b506119d386898988888e6133a0565b50505050505050505050565b60408151146119eb5750565b5f6119f7826020015190565b90506103fd81836133d5565b611a1f8260a4355b331415600182116004909210919091161690565b15611a6b57805f611a2e825190565b9050608081901c63ffffffff8216611a4684826133f9565b601c840163fb5014fc6060529350611a6260a435888685613409565b5f6060526118e6565b505050565b7f00000000000000000000000000000000000000000000000000000000000000008015611aa1575f63929eee145d50565b63929eee145480611ab7575f63929eee145d5050565b50600163929eee145550565b7f00000000000000000000000000000000000000000000000000000000000000008015611b025763929eee145c15611b0257637fa8a9875f526004601cfd5b806103ff5763929eee145480611b2a5763929eee145c15611b2a57637fa8a9875f526004601cfd5b60018111156103fd57637fa8a9875f526004601cfd5b80838337505050565b5f63ffffffff8235166040519150808252602082018160051b81018060a084026020870183378293505b81841015611b8c5780845260209093019260a001611b73565b60405250919392505050565b9052565b5f63ffffffff8235166040519150808252602082018160051b81018060c084026020870183378293505b81841015611b8c5780845260209093019260c001611bc6565b5f611bf160a060408051918201905290565b9050611bff82826080611b40565b611c1c611c15611c10846080610ce5565b613451565b6080830152565b919050565b5f63ffffffff8235166040519150808252602082018160051b8101808360061b6020870183378293505b81841015611b8c57808452602090930192604001611c4b565b60605f611c7160016124e3565b86515f90600160e61b82351690806001600160401b03811115611c9657611c96615773565b604051908082528060200260200182016040528015611cbf578160200160208202803683370190505b50945060010160051b91505f60205b83811015611ec0575f611ce48c83613cf86104e3565b90505f805f611cf3848e612586565b6001600160781b0382166020880152919450925090505f829003611d1a5750505050611eb8565b6001600160781b0381166040808601919091528a8601849052845160a081015160c0820151608083015192909301518051600184119d909d179c600490931099509092915f5b81811015611e02575f838281518110611d7b57611d7b6157ba565b602002602001015190508b8151108d179c505f611d9d89898460800151613488565b90508160800151826060015103611dba5760608201819052611dcf565b611dc989898460600151613488565b60608301525b5f611de88360600151838a8a611de3361590565b6134c4565b606084018190526080909301929092525050600101611d60565b5087516060015180515f5b81811015611eac575f838281518110611e2857611e286157ba565b602002602001015190505f611e428b8b8460800151613488565b90508160800151826060015103611e5f5760608201819052611e74565b611e6e8b8b8460600151613488565b60608301525b5f611e898360600151838c8c611de336151590565b6060840181905260a0840180516080909501949094529092525050600101611e0d565b50505050505050505050505b602001611cce565b50506001600160e61b018103611ed857611ed8613517565b50611ee3888861277e565b5f8060205b8381101561202657858101519250821561201e575f611f0a8c83613cf86104e3565b9050885f03611f25575f87830181905260209091015261201e565b60048151608001516004811115611f3e57611f3e615351565b14611fc157611f578188866001600587901c038e613524565b611f6d575f87830181905260209091015261201e565b602080820151604083015183516080810151930151611fa69388936001600160781b039081169316913314156001909111168e17612a49565b611fbc575f87830181905260209091015261201e565b611fec565b611fd3815f015182608001518c612b4d565b878301819052935083611fec575f60209091015261201e565b886001900398505f815f0151905061201785825f015183602001518c85604001518660600151612e0d565b6001935050505b602001611ee8565b50806120345761203461357f565b5050509550959350505050565b85518551606091829161205481836159ef565b6001600160401b0381111561206b5761206b615773565b6040519080825280602002602001820160405280156120a457816020015b612091615082565b8152602001906001900390816120895790505b5092505f5b828110156120fc576120d78c5f8d84815181106120c8576120c86157ba565b60200260200101518c8c61358c565b8482815181106120e9576120e96157ba565b60209081029190910101526001016120a9565b505f5b818110156121555761212e8c60018c848151811061211f5761211f6157ba565b60200260200101518c5f61358c565b8484830181518110612142576121426157ba565b60209081029190910101526001016120ff565b506121638b84888a896135db565b9350505097509795505050505050565b5f6121846040808051918201905290565b905061219e61219a61219584610aac565b6138c4565b8252565b611c1c6121af6110ac846020610ce5565b6020830152565b5f6121cb8260600151518361014001516138e3565b81516001600160a01b03165f9081526020819052604090205461040b908390610b39565b82545f90610100900460ff161561221657811561220f5761220f856138f3565b505f61062b565b83546201000090046001600160781b031680156122715783156122415761223c86613904565b612271565b8454600160881b90046001600160781b031681106122715782156122685761226886613915565b5f91505061062b565b50600195945050505050565b632165628a5f526004601cfd5b33831480156122995750505050565b5f6122a2612e72565b61190160f01b5f9081526002828152602287815260428320908390528651939450929190601f601d840116106102e2606219840110161561230c576122e78688613926565b61190160f01b5f9081526002869052602282815260428220919052909750905061230f565b50815b61231c888285858a6139bc565b5050505050505050565b5f6123376040808051918201905290565b905061234861219a610da684610aac565b611c1c6121af610da6846020610ce5565b80518060051b6040019050602082038051602082527f4b9f2d36e1b4c93de62cc077b00b1a91d84b6c31b4a14e012718dcca230689e78383a190525050565b8351606090806001600160401b038111156123b5576123b5615773565b6040519080825280602002602001820160405280156123ee57816020015b6123db615082565b8152602001906001900390816123d35790505b5091505f5b81811015612451575f87828151811061240e5761240e6157ba565b6020026020010151905061242b89825f0151836020015185613b05565b84838151811061243d5761243d6157ba565b6020908102919091010152506001016123f3565b5061245f87838787876135db565b50505b95945050505050565b6124788282610160611b40565b612489610af4610aef846040610ce5565b6103fd61249a610b0a846060610ce5565b6060830152565b6040518135601f0163ffffffe01660200180838337913563ffffffff16815290810160405290565b5f6124db602060408051918201905290565b5f8152905090565b7f000000000000000000000000000000000000000000000000000000000000000080156125315763929eee145c1561252257637fa8a9875f526004601cfd5b8160010163929eee145d6103fd565b63929eee1454806125645763929eee145c1561255457637fa8a9875f526004601cfd5b8260010163929eee145d506103fd565b6001811461257957637fa8a9875f526004601cfd5b505060020163929eee1455565b5f805f80855f015190506125a38160a001518260c0015187613caa565b6125b657505f9250829150819050612777565b602086015160408701516001600160781b0391821694501691505f6004826080015160048111156125e9576125e9615351565b03612616576001838502189050801561260457612604613ccd565b50600193508392508291506127779050565b50818311831517801561262b5761262b613ccd565b608082015160011615848411161561264557612645613cda565b61264e826121b6565b5f81815260016020526040812091965061266c90879083908a6121ef565b61267f57505f9350839250612777915050565b805460ff1661269a5761269a835f0151878a6060015161228a565b8054608881901c806126ae57869150612771565b6001600160781b038260101c169150600186036126d2578181039650809550612771565b8086036126ed57908601858103868211029096039590612771565b80860296810291909502810186810387821102918290039695919003906001600160781b0386111561277157612731565b5f5b8215610b2f57908290069190612720565b61274461273e878461271e565b8861271e565b8015019687900496909504946001600160781b0386111561277157634e487b715f5260116020526024601cfd5b50505050505b9250925092565b805182515f5b8281101561286d575f84828151811061279f5761279f6157ba565b602002602001015190505f815f015190508381106127c4576127c48260200151613ce7565b5f8782815181106127d7576127d76157ba565b6020026020010151905080602001516001600160781b03165f036127fd57505050612865565b80516040808201519085015163bfb3f8ce5f8760200151600181111561282557612825615351565b14612841575f612836856060613cf8565b9350636088d7de9150505b8251821061285257805f526004601cfd5b61285d838389613d03565b505050505050505b600101612784565b505f5b818110156128f0575f85828151811061288b5761288b6157ba565b6020026020010151905080602001516001600160781b03165f036128af57506128e8565b8051608081015160608201516128cc9085908363a8930e9a613db4565b6128e48483604001518363d69293326104e3613db490565b5050505b600101612870565b5050505050565b60a084015160c08501516040860151515f805b82811015612974575f89604001518281518110612929576129296157ba565b602002602001015190505f815f01519050801584179350505f612960826060015183608001518c8c8b8b61295b361590565b613e1b565b60608301525060800186905260010161290a565b506080880151600481108216801561298e5761298e613517565b505050506060860151515f5b8181101561231c575f886060015182815181106129b9576129b96157ba565b602002602001015190505f6129de826060015183608001518b8b8a8a61295b36151590565b60608301525060a081015160809091015260010161299a565b8351608081015160208201513314156001821160049092109190911616156128f0575f80612a2c858489608001518988613e56565b63fb5014fc6060529092509050611a628360200151868484613409565b5f848152600160205260408120805482908290608881901c80612a6e57889150612af6565b6001600160781b038260101c169150808803612a9257908801878111935090612af6565b97880297808802979190910288018781119350906001600160781b038083119089111715612af657612ac4888361271e565b8015019788900497909104906001600160781b038083119089111715612af657634e487b715f5260116020526024601cfd5b508215612b2f578515612b23576040516310fda3e160e01b8152600481018a905260240160405180910390fd5b5f94505050505061062b565b8660881b8160101b1760011782556001945050505050949350505050565b5f83610140015184606001515114612b6757612b6761227d565b83515f8080612b768888613fd4565b915091505f8082845f885af16001600160a01b0385165f908152600260205260409020805460018101909155606086901b189550925082612bd5578515612bc857612bbf61407e565b612bc8856140c5565b505f93506104fa92505050565b505050505f805f612bf2876040015188606001516104e36140d690565b925092509250825f14612c0857612c08846140c5565b60408701919091526060860152509392505050565b6040805160208082528183019092525f916020820181803683375050506040850151519091505f5b81811015612c95575f86604001518281518110612c6457612c646157ba565b60200260200101519050846080820152612c8c81885f0151896101200151876104e361434f90565b50600101612c45565b50506060840151515f90815b81811015612d23575f87606001518281518110612cc057612cc06157ba565b602002602001015190505f6005811115612cdc57612cdc615351565b81516005811115612cef57612cef615351565b03612d0b574793508381606001511115612d0b57612d0b6132b3565b612d1a8133898861434f6104e3565b50600101612ca1565b5050612d2e826119df565b504780156128f0576128f033826132c0565b8251608081015160208201515f92839283928392916004811060019091111633909114151615612da257612d83612d7d61010083015190565b5190565b88614444565b9093509150612d9460208201612d79565b945063fb5014fc9350612dee565b600481608001516004811115612dba57612dba615351565b0361231c57805194505f8560601b9050612ddb87838b608001518b85614485565b639397928596509094509250612dee9050565b612df86060859052565b612e0485878585613409565b5f60605261231c565b60608290506060829050856001600160a01b0316876001600160a01b03167f9d9af8e38d66c62e2c12f0225249fd9d721c54b83f48d9352c97c6cacdcb6f318a888686604051612e609493929190615a3b565b60405180910390a35050505050505050565b5f7f00000000000000000000000000000000000000000000000000000000000000004614612f2a575060408051608080517f00000000000000000000000000000000000000000000000000000000000000005f9081527f00000000000000000000000000000000000000000000000000000000000000006020527f0000000000000000000000000000000000000000000000000000000000000000855246606090815230845260a08220949095529093529190915290565b507f000000000000000000000000000000000000000000000000000000000000000090565b63fed398fc5f526004601cfd5b600435602014610224356102401416610244356102606102643560061b01141660186101243510600160a01b60843560a4351760c4356024351717101616806103ff576103ff614550565b5f8181526001602081905260409091209060843590612fcc90849084903615156121ef565b50815460ff16610b2f57610b2f8184612ff3602463ffffffff6102443516016124a16104e3565b61228a565b5f6130058260a435611a0b565b1561040b575f805f6130168661455d565b63fb5014fc6060529194509250905061303660a43587601c860185613409565b5f60605260209190910160801b1781529392505050565b806103ff576391b3e5145f526004601cfd5b604080517f000000000000000000000000000000000000000000000000000000000000000060ff60a01b175f90815260208690527f000000000000000000000000000000000000000000000000000000000000000083526055600b20919092526001600160a01b031690505f805f805260205f85875f875af191505f519050816130f4576130eb61407e565b6130f483614615565b6001600160e01b03198116632671a55160e11b1461181c5761181c8684614626565b6369f958275f52806020526024601cfd5b833b61313e57635f15d6725f52836020526024601cfd5b6040516323b872dd60e01b5f528360045282602452816044525f8060645f80895af1806131ce573d156131ac57601f3d0160051c8260051c8160030281831115613195578183036003028280028480020360091c01015b5a6020820110156131a8573d5f803e3d5ffd5b5050505b63f486bc875f5285602052846040528360605282608052600160a05260a4601cfd5b5060405250505f6060525050565b843b6131f357635f15d6725f52846020526024601cfd5b60405160805160a05160c051637921219560e11b5f528760045286602452856044528460645260a06084525f60a4525f8060c45f808d5af180613298573d1561327757601f3d0160051c8560051c8160030281831115613260578183036003028280028480020360091c01015b5a602082011015613273573d5f803e3d5ffd5b5050505b63f486bc875f52896020528860405287606052866080528560a05260a4601cfd5b5060809290925260a05260c05260405250505f606052505050565b638ffff9805f526004601cfd5b6132c98161304d565b5f805f805f85875af1905080611a6b576132e161407e565b63bc806b965f5282602052816040526044601cfd5b5f613302836020015190565b9050818114611a6b57611a6b836119df565b5f602088510361334e5750604080885260208089018a9052632671a55160e11b91890191909152604488015260016064880181905261335d565b50606487018051600101908190525b603c60c082028901038781528660208201528560408201528460608201528360808201528260a082015250505050505050505050565b636ab37ce75f526004601cfd5b6133a98361304d565b6133b381836132f6565b816133c4576118018686868661463b565b61181c828260018989895f8a613314565b6064810151604082019060c0026044016133f084838361305f565b50506020905250565b6317b1f9428252600181526103fd565b5f806001600160e01b03198451165f805260205f85875f8b5af15f51909350149050816134425761343861407e565b846080526024607cfd5b8061181c57846080526024607cfd5b5f8063ffffffff83351690506001810160051b6134748160408051918201905290565b9250613481848483611b40565b5050919050565b5f8284036134975750806104fa565b82848309156134ad5763c63cf0895f526004601cfd5b5f6134b88584615ad7565b93909304949350505050565b5f84861461350d57838303428590038082035f6134e1838a615ad7565b6134eb838c615ad7565b6134f591906159ef565b90508584878303040181151502945050505050612462565b5092949350505050565b6312d3f5a35f526004601cfd5b8451608081015160208201515f92916004811060019091111633909114151615610f59575f8061355b87848b608001518b8a613e56565b91509150613575836020015188848463fb5014fc8a61472f565b9350505050612462565b63d5da9a1b5f526004601cfd5b613594615082565b83515f036135a5576135a58561479c565b5f8560018111156135b8576135b8615351565b036135ce576135c9868583856147ad565b612462565b612462868583338761491d565b84516060905f816001600160401b038111156135f9576135f9615773565b604051908082528060200260200182016040528015613622578160200160208202803683370190505b506040805160208082528183019092529192505f9190602082018180368337505089519192505060010160051b60205b818110156136b1575f6136688b83613cf86104e3565b805160608101519192509080156136a657478111825115161561369257638ffff9805f526004601cfd5b6136a682846020015185604001518961434f565b505050602001613652565b50505f5b8381101561381d575f8a82815181106136d0576136d06157ba565b6020026020010151905080602001516001600160781b03165f03613717575f848381518110613701576137016157ba565b9115156020928302919091019091015250613815565b600184838151811061372b5761372b6157ba565b911515602092830291909101909101528051604081015180515f5b818110156137b0575f838281518110613761576137616157ba565b6020026020010151905080606001515f1461379d57608081018051908e905285516101208701516137979184918c61434f6104e3565b60808201525b6080810151606090910152600101613746565b505050606081015180515f5b8181101561380f575f8382815181106137d7576137d76157ba565b602002602001015190505f81606001519050805f146137fb576137fb888483614a54565b5060a08101516060909101526001016137bc565b50505050505b6001016136b5565b50613827816119df565b4780156138385761383833826132c0565b85156138ae575f5b848110156138ac5783818151811061385a5761385a6157ba565b6020026020010151156138a4576138a48b828151811061387c5761387c6157ba565b60200260200101518a8b8481518110613897576138976157ba565b6020026020010151612d40565b600101613840565b505b6138b6611a70565b509098975050505050505050565b5f6138d761016060408051918201905290565b9050611c1c828261246b565b808210156103fd576103fd614a6d565b631a5155745f52806020526024601cfd5b63ee9e0e635f52806020526024601cfd5b6310fda3e15f52806020526024601cfd5b5f805f84516001811660410380820360051c9250808752806020018701915050805160e81c6003820191506001811660051b868152825160208218525060015b838110156139925760405f2082821c60051b602090811691825293840180519190941852600101613966565b50505060405f2091505f6139a582614a7a565b5f9081526020939093525050604090209392505050565b5f805f528151602083038051826041035f60018211613a1f57604087015160608801515f1a8315613a0057601b8260ff1c0190506001600160ff1b03821660408a01525b88528a855260205f60808760015afa508385528588526040880152505f515b8a148a1515169450849050613ae857858552604082526044850380516040870351630b135d3f60e11b835289604089035260205f60648b01858f5afa96508615613adc57630b135d3f60e11b5f5114613adc578b3b15613a8657634f7fb80d5f526004601cfd5b6001866041031115613a9f57638baa579f5f526004601cfd5b64010100000060608901515f1a1a15604187141615613acf57631f003d0a5f5260608801515f1a6020526024601cfd5b63815e1d645f526004601cfd5b8385529152603f198601525b5050508061181c57613af861407e565b634f7fb80d5f526004601cfd5b613b0d615082565b8251158451151715613b26576398e9db6e5f526004601cfd5b613b2e615082565b613b3b8685835f8061491d565b805160608101515f03613b505750905061062b565b613b6087878584608001516147ad565b82516040828101519082015160208085015190840151855185511891181791181715613b975763bced929d5f52846020526024601cfd5b806060015182606001511115613c1e575f865f81518110613bba57613bba6157ba565b60200260200101519050816060015183606001510389825f015181518110613be457613be46157ba565b60200260200101515f015160600151826020015181518110613c0857613c086157ba565b6020026020010151606001818152505050613c9f565b5f875f81518110613c3157613c316157ba565b60200260200101519050826060015182606001510389825f015181518110613c5b57613c5b6157ba565b60200260200101515f015160400151826020015181518110613c7f57613c7f6157ba565b602002602001015160600181815250508260600151826060018181525050505b505050949350505050565b428084111590831116818015613cbe575080155b156104fa576104fa8484614e89565b635a052b325f526004601cfd5b63a11b63ff5f526004601cfd5b63133c37c65f52806020526024601cfd5b5f6104fa8284015190565b5f838381518110613d1657613d166157ba565b602002602001015190505f815f01519050613d318160031090565b613d3d57613d3d614e9e565b60408201518015613d6057613d5b8460600151828660800151614eab565b613d73565b60808401515115613d7357613d73614ef5565b600119820183816005811115613d8b57613d8b615351565b90816005811115613d9e57613d9e615351565b9052505050606090920151604090910152505050565b82515f5b8181101561181c575f858281518110613dd357613dd36157ba565b60209081029190910101518051604082015191925090600382116004881415821515171615613e0d57855f5288602052836040526044601cfd5b505050806001019050613db8565b5f868803613e3557613e2e868689613488565b90506107e8565b6107e5613e4387878b613488565b613e4e88888b613488565b8686866134c4565b5f805f613e61614f02565b6301e4d72a815260208082015260408101898152336060830152601c820194509091508751604082015287613ea1613e9a60a083015190565b60e0840152565b613eb7613eaf60c083015190565b610100840152565b613ecd613ec560e083015190565b610120840152565b610140613edb816060850152565b5f613ee7604084015190565b90505f613ef682848701614f0c565b928301929050613f07836080870152565b5f613f13606086015190565b90505f613f2282868901614f72565b948501949050613f338560a0890152565b5f613f408e878a01614fd7565b959095019450613f518560c0890152565b8685015f613f5f8e83614ff6565b602497019687019a50613f7b9050613f768c8c0190565b615026565b8060408b901b60808b901b17178f610100018181525050613fa58c82611b9890919063ffffffff16565b60058c8e51613fb491906159dc565b613fbf911b8b6159dc565b99505050505050505050509550959350505050565b5f8083613fdf614f02565b639891976581523360208201908152608060408301819052601c9092019450905f61400b604085015190565b90505f61401a82848601614f0c565b92830192905061402b836040860152565b5f614037606087015190565b90505f61404682868801614f0c565b948501949050614057856060880152565b895f61406582898901614fd7565b9a9d96909a016004019b50949950505050505050505050565b3d156140c357601f3d0160051c60405160051c81600302818311156140b0578183036003028280028480020360091c01015b5a602082011015611a6b573d5f803e3d5ffd5b565b63939792855f52806020526024601cfd5b60603d105f8080808080866141485760405f803e5f51935060205192503d60208501113d60208501118082179850505086614148576020845f3e5f51915060208360203e60205190508160071b60208501018160a0026020850101803d10823d101761ffff8486171117985050505f80525b8661417a575f8061415d84602088018d614183565b9250975061416f83602087018c614263565b929092179850909550505b50505050612777565b5f806040519150825160c08602602001830160405285835260208660010160051b8085018360010160051b87016141c1858b81811090829003020190565b60010160051b8a861196505b8085101561422157828589015260808a843e6060830151955085608084015260608201518681116141fe858561430a565b17881797505060808a01995060a08301925060a0820191506020850194506141cd565b50505b81831015614257578083870152608088823e6060810151608082015260808801975060a081019050602083019250614224565b50505050935093915050565b604051815180851190808603818710028101602060e08202850181016040528185526001928301600590811b87019390920190911b908185015b8282101561425757808287015260a088823e60206060890160a083013e606081015160608501516142da6080840151608088015180159114171590565b818311176142e8848861430a565b60a09b909b019a179690961795505060c093840193602092909201910161429d565b5f81516040830151801560038311161561432c57506040840151600119909101905b604085015181148551831460208701516020870151141616159250505092915050565b5f8451600581111561436357614363615351565b036143a057604084015160208501516001600160a01b0316171561438957614389613393565b61439b846080015185606001516132c0565b61443e565b6001845160058111156143b5576143b5615351565b036143e6576040840151156143cc576143cc613393565b61439b8460200151848660800151876060015186866133a0565b6002845160058111156143fb576143fb615351565b0361441f5761439b84602001518486608001518760400151886060015187876118a4565b61443e84602001518486608001518760400151886060015187876118ef565b50505050565b608082901c63ffffffff604084901c8116908085169061446c9084906317b1f94290611b9816565b601c8301925061447c8482614ff6565b50509250929050565b5f8061448f614f02565b63f4dd92ce815287841860a0820152601c8101925060200160a0808252875f6144b9604083015190565b90505f6144c882858701614f0c565b9384019390506144d9846020870152565b5f6144e5606085015190565b90505f6144f482878901614f72565b958601959050614505866040890152565b5f6145128d888a01614fd7565b9687019690506145238760608a0152565b5f6145308d898b01614ff6565b905080880197508760040199505050505050505050509550959350505050565b6339f3e3fd5f526004601cfd5b6301e4d72a6102043560051b6080019081525f808260208082015260408101858152336060830152601c9190910190614597608435610af4565b6145a861014460e083016060611b40565b6101406145b6816060840152565b6145c460a082016080840152565b61016060a06102643581029290920101906145e0908290840152565b6145ea5f82840152565b6020016145f88160c0840152565b5f9181019182526020820196909652939560449095019492505050565b63d13d53d45f52806020526024601cfd5b631cf99b265f5281602052806040526044601cfd5b6040516323b872dd60e01b5f5283600452826024528160445260205f60645f80895af1803d15601f3d1160015f51141617163d151581166147205780873b151516614720578061470f57816146f2573d156146d157601f3d0160051c8360051c81600302818311156146ba578183036003028280028480020360091c01015b5a6020820110156146cd573d5f803e3d5ffd5b5050505b63f486bc875f528660205285604052846060525f6080528360a05260a4601cfd5b63988919235f52866020528560405284606052836080526084601cfd5b635f15d6725f52866020526024601cfd5b505060405250505f6060525050565b5f805f6001600160e01b03198751165f805260205f888a5f8e5af15f519093501490508161477b5783614766575f925050506106e0565b61476e61407e565b845f52876020526024601cfd5b8061478d57845f52876020526024601cfd5b50600198975050505050505050565b63375c24c15f52806020526024601cfd5b5f805f85865160051b87015b808210156148ca576020820191508851825151106147d9576147d96148fe565b81515160051b60208a01015180516020845101515f60408301516020850151158151841015171561480e5750505050506147b9565b8260051b60208201015191505060608101935083518901915083511589831060011b17881797508198505f84528a5193508615600181146148775760608220881860408d01516101208601511860208e0151865118171715614872576148726148fe565b6148c0565b8151855260208201516020860152604082015160408601528a6080860152835160208d015261012084015160408d015260608520975060208d019250868318156148c057865183525b50505050506147b9565b5050508160608551015280156148f757600181036148ef576391b3e5145f526004601cfd5b6148f761490b565b505061443e565b637fda72795f526004601cfd5b634e487b715f5260116020526024601cfd5b5f805f86875160051b88015b80821015614a1f576020820191508151518a51811061494a5761494a6148fe565b8060051b60208c01015190506020835101515f6060835101516020840151158151841015171561497d5750505050614929565b8260051b60208201015191505060608101925082518801915082511588831060011b17871796508197505f83528a5192508515600181146149cf5760a0822087146149ca576149ca6148fe565b614a16565b815184526020820151602085015260408201516040850152608082015160808501528a60208d01528960408d015260a08220965060208d01925085831815614a1657855183525b50505050614929565b50508551606001839052508015614a4d5760018103614a45576391b3e5145f526004601cfd5b614a4d61490b565b50506128f0565b63a5f542085f528260205281604052806060526064601cfd5b63466aa6165f526004601cfd5b5f614e80565b5f6009821015614bd9576005821015614b36576003821015614ae9577f832c58a5b611aadcfa6a082ac9d04bace53d8278387f10040347b7e98eb5b30260018314027fbf8e29b89f29ed9b529c154a63038ffca562f8d7cd1e2545dda53a1b582dde301861040b565b7ff3e8417a785f980bdaf134fa0274a6bf891eeb8195cd94b09d2aa651046e28bc60038314027fa02eb7ff164c884e5e2c336dc85f81c6a93329d8e9adf214b32729b894de2af11861040b565b6007821015614b8c577f25d02425402d882d211a7ab774c0ed6eca048c4d03d9af40132475744753b2a360058314027f1c19f71958cdd8f081b4c31f7caf5c010b29d12950be2fa1c95070dc47e30b551861040b565b7fb58d772fb09b426b9dece637f61ca9065f2b994f1464b51e9207f55f7c8f594860078314027f7ff98d9d4e55d876c5cfac10b43c04039522f3ddfb0ea9bfe70c68cfb5c7cc141861040b565b6011821015614d3157600d821015614c8e57600b821015614c41577f6f0ec38c21f6f583ab7f3c5413c773ffd5344c34fde1d390958e438bf667448f60098314027fd1d97d1ef5eaa37a4ee5fbf234e6f6d64eb511eb562221cd7edfbdde0848da051861040b565b7f32f4e7485d6485f9f6c255929b9905c62ba919758bbe231f231eaeecf33d810c600b8314027fbb98d87cc12922b83759626c5f07d72266da9702d19ffad6a514c73a89002f5f1861040b565b600f821015614ce4577f8df51df98847160517f5b1186b4bc3f418d98b8a7f17f1292f392d79d600d79e600d8314027f6b5b04cbae4fcb1a9d78e7b2dfc51a36933d023cf6e347e03d517b472a8525901861040b565b7fcc4886e37eedd9aacd6c1c2c9247197a621a71282e87a7cbc673f3736d9aa141600f8314027f1da3eed3ecef6ebaa6e5023c057ec2c75150693fd0dac5c90f4a142f9879fde81861040b565b6015821015614ddd576013821015614d90577f2d7a3ed6dab270fdb8e054b2ad525f0ce2a8b89cc76c17f0965434740f673a5560118314027fc3939feff011e53ab8c35ca3370aad54c5df1fc2938cd62543174fa6e7d858771861040b565b7f54b3212a178782f104e0d514b41a9a5c4ca9c980bf6597c3cecbf280917e202a60138314027f5a4f867d3d458dabecad65f6201ceeaba0096df2d0c491cc32e6ea4e643500171861040b565b6017821015614e33577fbb40bf8cea3a5a716e2b6eb08bbdac8ec159f82f380783db3c56904f15a43d0460158314027f3bd8cff538aba49a9c374c806d277181e9651624b3e31111bc0624574f8bca1d1861040b565b7f403be09941a31d05cfc2f896505811353d45d38743288b016630cce39435476a60178314027f1d51df90cba8de7637ca3e8fe1e3511d1dc2f23487d05dbdecb781860c21ac1c1861040b565b61040b82614a80565b6321ccfeb75f5281602052806040526044601cfd5b6394eb6af65f526004601cfd5b5f835f5260205f2060208301835160051b81015b80821015614ee657815180841160051b93845260209384185260405f209290910190614ebf565b5050831490508061443e5761443e5b6309bde3395f526004601cfd5b5f61041b60405190565b5f825180835260208401602084018260051b82015b80831015614f5f5782518051835260208101516020840152604081015160408401526060810151606084015250602083019250608082019150614f21565b5050508060071b60200191505092915050565b5f80614f7c845190565b8084529050602084810190600583901b860181019085015b82821115614fc4575f614fa684615030565b9050614fb4818360a0615039565b506020929092019160a001614f94565b60a0840260200194505050505092915050565b5f63ffffffe0603f614fe7855190565b0116905061040b838383615039565b5f80615000845190565b8084529050600581901b61501b602086810190860183615039565b602001949350505050565b6103ff6040829052565b5f61040b825190565b8082828560045afa80153d15171561443e575f80fd5b6040518060a001604052806150626150c4565b81525f602082018190526040820152606080820181905260809091015290565b60408051610100810182525f606082018181526080830182905260a0830182905260c0830182905260e083018290528252602082018190529181019190915290565b6040518061016001604052805f6001600160a01b031681526020015f6001600160a01b0316815260200160608152602001606081526020015f600481111561510e5761510e615351565b81525f6020820181905260408201819052606082018190526080820181905260a0820181905260c09091015290565b5f6020828403121561514d575f80fd5b81356001600160401b03811115615162575f80fd5b820161024081850312156104fa575f80fd5b5f81518084525f5b818110156151985760208185018101518683018201520161517c565b505f602082860101526020601f19601f83011685010191505092915050565b602081525f6104fa6020830184615174565b5f602082840312156151d9575f80fd5b5035919050565b5f602082840312156151f0575f80fd5b81356001600160401b03811115615205575f80fd5b820161016081850312156104fa575f80fd5b5f8083601f840112615227575f80fd5b5081356001600160401b0381111561523d575f80fd5b6020830191508360208260051b8501011115615257575f80fd5b9250929050565b80356001600160a01b0381168114611c1c575f80fd5b5f805f805f805f805f805f60e08c8e03121561528e575f80fd5b6001600160401b03808d3511156152a3575f80fd5b6152b08e8e358f01615217565b909c509a5060208d01358110156152c5575f80fd5b6152d58e60208f01358f01615217565b909a50985060408d01358110156152ea575f80fd5b6152fa8e60408f01358f01615217565b909850965060608d013581101561530f575f80fd5b506153208d60608e01358e01615217565b909550935060808c0135925061533860a08d0161525e565b915060c08c013590509295989b509295989b9093969950565b634e487b7160e01b5f52602160045260245ffd5b60068110611b9857611b98615351565b615380828251615365565b6020818101516001600160a01b0390811691840191909152604080830151908401526060808301519084015260809182015116910152565b5f815180845260208085019450602084015f5b8381101561540d5781516153e0888251615375565b808401516001600160a01b031660a08901526040015160c088015260e090960195908201906001016153cb565b509495945050505050565b604080825283519082018190525f906020906060840190828701845b82811015615452578151151584529284019290840190600101615434565b50505083810360208501526106e081866153b8565b5f8060208385031215615478575f80fd5b82356001600160401b0381111561548d575f80fd5b61549985828601615217565b90969095509350505050565b5f805f80604085870312156154b8575f80fd5b84356001600160401b03808211156154ce575f80fd5b6154da88838901615217565b909650945060208701359150808211156154f2575f80fd5b506154ff87828801615217565b95989497509550505050565b602081525f6104fa60208301846153b8565b5f6020828403121561552d575f80fd5b6104fa8261525e565b5f8060408385031215615547575f80fd5b82356001600160401b0381111561555c575f80fd5b83016040818603121561556d575f80fd5b946020939093013593505050565b5f805f805f6080868803121561558f575f80fd5b85356001600160401b03808211156155a5575f80fd5b9087019060a0828a0312156155b8575f80fd5b909550602087013590808211156155cd575f80fd5b506155da88828901615217565b909550935050604086013591506155f36060870161525e565b90509295509295909350565b5f805f805f805f8060a0898b031215615616575f80fd5b88356001600160401b038082111561562c575f80fd5b6156388c838d01615217565b909a50985060208b0135915080821115615650575f80fd5b61565c8c838d01615217565b909850965060408b0135915080821115615674575f80fd5b506156818b828c01615217565b999c989b509699959896976060870135966080013595509350505050565b5f805f805f805f6080888a0312156156b5575f80fd5b87356001600160401b03808211156156cb575f80fd5b6156d78b838c01615217565b909950975060208a01359150808211156156ef575f80fd5b6156fb8b838c01615217565b909750955060408a0135915080821115615713575f80fd5b506157208a828b01615217565b909450925061573390506060890161525e565b905092959891949750929550565b606081525f6157536060830186615174565b6020830194909452506001600160a01b0391909116604090910152919050565b634e487b7160e01b5f52604160045260245ffd5b634e487b7160e01b5f52601160045260245ffd5b5f826157b557634e487b7160e01b5f52601260045260245ffd5b500490565b634e487b7160e01b5f52603260045260245ffd5b5f815180845260208085019450602084015f5b8381101561540d5781516157f6888251615365565b838101516001600160a01b03168885015260408082015190890152606080820151908901526080908101519088015260a090960195908201906001016157e1565b5f815180845260208085019450602084015f5b8381101561540d57815161585f888251615365565b808401516001600160a01b0390811689860152604080830151908a0152606080830151908a0152608080830151908a015260a091820151169088015260c0909601959082019060010161584a565b60058110611b9857611b98615351565b828152604060208201526158dd6040820183516001600160a01b03169052565b5f60208301516158f860608401826001600160a01b03169052565b5060408301516101608060808501526159156101a08501836157ce565b91506060850151603f198584030160a08601526159328382615837565b925050608085015161594760c08601826158ad565b5060a085015160e085015260c0850151610100818187015260e0870151915061012082818801528188015192506101409150828288015280880151848801525080870151610180870152505050809150509392505050565b5f823561015e198336030181126159b4575f80fd5b9190910192915050565b5f602082840312156159ce575f80fd5b8135600581106104fa575f80fd5b8181038181111561040b5761040b615787565b8082018082111561040b5761040b615787565b5f815180845260208085019450602084015f5b8381101561540d57615a28878351615375565b60a0969096019590820190600101615a15565b5f6080808301878452602060018060a01b03808916602087015260406080604088015283895180865260a08901915060208b0195505f5b81811015615ab3578651615a87848251615365565b808701518616848801528481015185850152606090810151908401529585019591870191600101615a72565b50508781036060890152615ac7818a615a02565b9c9b505050505050505050505050565b808202811582820484141761040b5761040b61578756fea2646970667358221220c8a534726d5a7361f31fb24d3afb58707e2e475615cdc3e4c9e9d257965805ca64736f6c63430008180033436f6e73696465726174696f6e4974656d2875696e7438206974656d547970652c6164647265737320746f6b656e2c75696e74323536206964656e7469666965724f7243726974657269612c75696e74323536207374617274416d6f756e742c75696e7432353620656e64416d6f756e742c6164647265737320726563697069656e7429454950373132446f6d61696e28737472696e67206e616d652c737472696e672076657273696f6e2c75696e7432353620636861696e49642c6164647265737320766572696679696e67436f6e7472616374294f72646572436f6d706f6e656e74732861646472657373206f6666657265722c61646472657373207a6f6e652c4f666665724974656d5b5d206f666665722c436f6e73696465726174696f6e4974656d5b5d20636f6e73696465726174696f6e2c75696e7438206f72646572547970652c75696e7432353620737461727454696d652c75696e7432353620656e6454696d652c62797465733332207a6f6e65486173682c75696e743235362073616c742c6279746573333220636f6e647569744b65792c75696e7432353620636f756e746572294f666665724974656d2875696e7438206974656d547970652c6164647265737320746f6b656e2c75696e74323536206964656e7469666965724f7243726974657269612c75696e74323536207374617274416d6f756e742c75696e7432353620656e64416d6f756e7429",
}

// SeaportABI is the input ABI used to generate the binding from.
// Deprecated: Use SeaportMetaData.ABI instead.
var SeaportABI = SeaportMetaData.ABI

// SeaportBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SeaportMetaData.Bin instead.
var SeaportBin = SeaportMetaData.Bin

// DeploySeaport deploys a new Ethereum contract, binding an instance of Seaport to it.
func DeploySeaport(auth *bind.TransactOpts, backend bind.ContractBackend, conduitController common.Address) (common.Address, *types.Transaction, *Seaport, error) {
	parsed, err := SeaportMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SeaportBin), backend, conduitController)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Seaport{SeaportCaller: SeaportCaller{contract: contract}, SeaportTransactor: SeaportTransactor{contract: contract}, SeaportFilterer: SeaportFilterer{contract: contract}}, nil
}

// Seaport is an auto generated Go binding around an Ethereum contract.
type Seaport struct {
	SeaportCaller     // Read-only binding to the contract
	SeaportTransactor // Write-only binding to the contract
	SeaportFilterer   // Log filterer for contract events
}

// SeaportCaller is an auto generated read-only Go binding around an Ethereum contract.
type SeaportCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SeaportTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SeaportTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SeaportFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SeaportFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SeaportSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SeaportSession struct {
	Contract     *Seaport          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SeaportCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SeaportCallerSession struct {
	Contract *SeaportCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// SeaportTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SeaportTransactorSession struct {
	Contract     *SeaportTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// SeaportRaw is an auto generated low-level Go binding around an Ethereum contract.
type SeaportRaw struct {
	Contract *Seaport // Generic contract binding to access the raw methods on
}

// SeaportCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SeaportCallerRaw struct {
	Contract *SeaportCaller // Generic read-only contract binding to access the raw methods on
}

// SeaportTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SeaportTransactorRaw struct {
	Contract *SeaportTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSeaport creates a new instance of Seaport, bound to a specific deployed contract.
func NewSeaport(address common.Address, backend bind.ContractBackend) (*Seaport, error) {
	contract, err := bindSeaport(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Seaport{SeaportCaller: SeaportCaller{contract: contract}, SeaportTransactor: SeaportTransactor{contract: contract}, SeaportFilterer: SeaportFilterer{contract: contract}}, nil
}

// NewSeaportCaller creates a new read-only instance of Seaport, bound to a specific deployed contract.
func NewSeaportCaller(address common.Address, caller bind.ContractCaller) (*SeaportCaller, error) {
	contract, err := bindSeaport(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SeaportCaller{contract: contract}, nil
}

// NewSeaportTransactor creates a new write-only instance of Seaport, bound to a specific deployed contract.
func NewSeaportTransactor(address common.Address, transactor bind.ContractTransactor) (*SeaportTransactor, error) {
	contract, err := bindSeaport(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SeaportTransactor{contract: contract}, nil
}

// NewSeaportFilterer creates a new log filterer instance of Seaport, bound to a specific deployed contract.
func NewSeaportFilterer(address common.Address, filterer bind.ContractFilterer) (*SeaportFilterer, error) {
	contract, err := bindSeaport(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SeaportFilterer{contract: contract}, nil
}

// bindSeaport binds a generic wrapper to an already deployed contract.
func bindSeaport(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SeaportMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Seaport *SeaportRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Seaport.Contract.SeaportCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Seaport *SeaportRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Seaport.Contract.SeaportTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Seaport *SeaportRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Seaport.Contract.SeaportTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Seaport *SeaportCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Seaport.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Seaport *SeaportTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Seaport.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Seaport *SeaportTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Seaport.Contract.contract.Transact(opts, method, params...)
}

// GetContractOffererNonce is a free data retrieval call binding the contract method 0xa900866b.
//
// Solidity: function getContractOffererNonce(address contractOfferer) view returns(uint256 nonce)
func (_Seaport *SeaportCaller) GetContractOffererNonce(opts *bind.CallOpts, contractOfferer common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Seaport.contract.Call(opts, &out, "getContractOffererNonce", contractOfferer)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetContractOffererNonce is a free data retrieval call binding the contract method 0xa900866b.
//
// Solidity: function getContractOffererNonce(address contractOfferer) view returns(uint256 nonce)
func (_Seaport *SeaportSession) GetContractOffererNonce(contractOfferer common.Address) (*big.Int, error) {
	return _Seaport.Contract.GetContractOffererNonce(&_Seaport.CallOpts, contractOfferer)
}

// GetContractOffererNonce is a free data retrieval call binding the contract method 0xa900866b.
//
// Solidity: function getContractOffererNonce(address contractOfferer) view returns(uint256 nonce)
func (_Seaport *SeaportCallerSession) GetContractOffererNonce(contractOfferer common.Address) (*big.Int, error) {
	return _Seaport.Contract.GetContractOffererNonce(&_Seaport.CallOpts, contractOfferer)
}

// GetCounter is a free data retrieval call binding the contract method 0xf07ec373.
//
// Solidity: function getCounter(address offerer) view returns(uint256 counter)
func (_Seaport *SeaportCaller) GetCounter(opts *bind.CallOpts, offerer common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Seaport.contract.Call(opts, &out, "getCounter", offerer)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCounter is a free data retrieval call binding the contract method 0xf07ec373.
//
// Solidity: function getCounter(address offerer) view returns(uint256 counter)
func (_Seaport *SeaportSession) GetCounter(offerer common.Address) (*big.Int, error) {
	return _Seaport.Contract.GetCounter(&_Seaport.CallOpts, offerer)
}

// GetCounter is a free data retrieval call binding the contract method 0xf07ec373.
//
// Solidity: function getCounter(address offerer) view returns(uint256 counter)
func (_Seaport *SeaportCallerSession) GetCounter(offerer common.Address) (*big.Int, error) {
	return _Seaport.Contract.GetCounter(&_Seaport.CallOpts, offerer)
}

// GetOrderHash is a free data retrieval call binding the contract method 0x79df72bd.
//
// Solidity: function getOrderHash((address,address,(uint8,address,uint256,uint256,uint256)[],(uint8,address,uint256,uint256,uint256,address)[],uint8,uint256,uint256,bytes32,uint256,bytes32,uint256) ) view returns(bytes32 orderHash)
func (_Seaport *SeaportCaller) GetOrderHash(opts *bind.CallOpts, arg0 OrderComponents) ([32]byte, error) {
	var out []interface{}
	err := _Seaport.contract.Call(opts, &out, "getOrderHash", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetOrderHash is a free data retrieval call binding the contract method 0x79df72bd.
//
// Solidity: function getOrderHash((address,address,(uint8,address,uint256,uint256,uint256)[],(uint8,address,uint256,uint256,uint256,address)[],uint8,uint256,uint256,bytes32,uint256,bytes32,uint256) ) view returns(bytes32 orderHash)
func (_Seaport *SeaportSession) GetOrderHash(arg0 OrderComponents) ([32]byte, error) {
	return _Seaport.Contract.GetOrderHash(&_Seaport.CallOpts, arg0)
}

// GetOrderHash is a free data retrieval call binding the contract method 0x79df72bd.
//
// Solidity: function getOrderHash((address,address,(uint8,address,uint256,uint256,uint256)[],(uint8,address,uint256,uint256,uint256,address)[],uint8,uint256,uint256,bytes32,uint256,bytes32,uint256) ) view returns(bytes32 orderHash)
func (_Seaport *SeaportCallerSession) GetOrderHash(arg0 OrderComponents) ([32]byte, error) {
	return _Seaport.Contract.GetOrderHash(&_Seaport.CallOpts, arg0)
}

// GetOrderStatus is a free data retrieval call binding the contract method 0x46423aa7.
//
// Solidity: function getOrderStatus(bytes32 orderHash) view returns(bool isValidated, bool isCancelled, uint256 totalFilled, uint256 totalSize)
func (_Seaport *SeaportCaller) GetOrderStatus(opts *bind.CallOpts, orderHash [32]byte) (struct {
	IsValidated bool
	IsCancelled bool
	TotalFilled *big.Int
	TotalSize   *big.Int
}, error) {
	var out []interface{}
	err := _Seaport.contract.Call(opts, &out, "getOrderStatus", orderHash)

	outstruct := new(struct {
		IsValidated bool
		IsCancelled bool
		TotalFilled *big.Int
		TotalSize   *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.IsValidated = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.IsCancelled = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.TotalFilled = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.TotalSize = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetOrderStatus is a free data retrieval call binding the contract method 0x46423aa7.
//
// Solidity: function getOrderStatus(bytes32 orderHash) view returns(bool isValidated, bool isCancelled, uint256 totalFilled, uint256 totalSize)
func (_Seaport *SeaportSession) GetOrderStatus(orderHash [32]byte) (struct {
	IsValidated bool
	IsCancelled bool
	TotalFilled *big.Int
	TotalSize   *big.Int
}, error) {
	return _Seaport.Contract.GetOrderStatus(&_Seaport.CallOpts, orderHash)
}

// GetOrderStatus is a free data retrieval call binding the contract method 0x46423aa7.
//
// Solidity: function getOrderStatus(bytes32 orderHash) view returns(bool isValidated, bool isCancelled, uint256 totalFilled, uint256 totalSize)
func (_Seaport *SeaportCallerSession) GetOrderStatus(orderHash [32]byte) (struct {
	IsValidated bool
	IsCancelled bool
	TotalFilled *big.Int
	TotalSize   *big.Int
}, error) {
	return _Seaport.Contract.GetOrderStatus(&_Seaport.CallOpts, orderHash)
}

// Information is a free data retrieval call binding the contract method 0xf47b7740.
//
// Solidity: function information() view returns(string version, bytes32 domainSeparator, address conduitController)
func (_Seaport *SeaportCaller) Information(opts *bind.CallOpts) (struct {
	Version           string
	DomainSeparator   [32]byte
	ConduitController common.Address
}, error) {
	var out []interface{}
	err := _Seaport.contract.Call(opts, &out, "information")

	outstruct := new(struct {
		Version           string
		DomainSeparator   [32]byte
		ConduitController common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Version = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.DomainSeparator = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.ConduitController = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// Information is a free data retrieval call binding the contract method 0xf47b7740.
//
// Solidity: function information() view returns(string version, bytes32 domainSeparator, address conduitController)
func (_Seaport *SeaportSession) Information() (struct {
	Version           string
	DomainSeparator   [32]byte
	ConduitController common.Address
}, error) {
	return _Seaport.Contract.Information(&_Seaport.CallOpts)
}

// Information is a free data retrieval call binding the contract method 0xf47b7740.
//
// Solidity: function information() view returns(string version, bytes32 domainSeparator, address conduitController)
func (_Seaport *SeaportCallerSession) Information() (struct {
	Version           string
	DomainSeparator   [32]byte
	ConduitController common.Address
}, error) {
	return _Seaport.Contract.Information(&_Seaport.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() pure returns(string)
func (_Seaport *SeaportCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Seaport.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() pure returns(string)
func (_Seaport *SeaportSession) Name() (string, error) {
	return _Seaport.Contract.Name(&_Seaport.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() pure returns(string)
func (_Seaport *SeaportCallerSession) Name() (string, error) {
	return _Seaport.Contract.Name(&_Seaport.CallOpts)
}

// ActivateTstore is a paid mutator transaction binding the contract method 0x7423eb3c.
//
// Solidity: function __activateTstore() returns()
func (_Seaport *SeaportTransactor) ActivateTstore(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Seaport.contract.Transact(opts, "__activateTstore")
}

// ActivateTstore is a paid mutator transaction binding the contract method 0x7423eb3c.
//
// Solidity: function __activateTstore() returns()
func (_Seaport *SeaportSession) ActivateTstore() (*types.Transaction, error) {
	return _Seaport.Contract.ActivateTstore(&_Seaport.TransactOpts)
}

// ActivateTstore is a paid mutator transaction binding the contract method 0x7423eb3c.
//
// Solidity: function __activateTstore() returns()
func (_Seaport *SeaportTransactorSession) ActivateTstore() (*types.Transaction, error) {
	return _Seaport.Contract.ActivateTstore(&_Seaport.TransactOpts)
}

// Cancel is a paid mutator transaction binding the contract method 0xfd9f1e10.
//
// Solidity: function cancel((address,address,(uint8,address,uint256,uint256,uint256)[],(uint8,address,uint256,uint256,uint256,address)[],uint8,uint256,uint256,bytes32,uint256,bytes32,uint256)[] orders) returns(bool cancelled)
func (_Seaport *SeaportTransactor) Cancel(opts *bind.TransactOpts, orders []OrderComponents) (*types.Transaction, error) {
	return _Seaport.contract.Transact(opts, "cancel", orders)
}

// Cancel is a paid mutator transaction binding the contract method 0xfd9f1e10.
//
// Solidity: function cancel((address,address,(uint8,address,uint256,uint256,uint256)[],(uint8,address,uint256,uint256,uint256,address)[],uint8,uint256,uint256,bytes32,uint256,bytes32,uint256)[] orders) returns(bool cancelled)
func (_Seaport *SeaportSession) Cancel(orders []OrderComponents) (*types.Transaction, error) {
	return _Seaport.Contract.Cancel(&_Seaport.TransactOpts, orders)
}

// Cancel is a paid mutator transaction binding the contract method 0xfd9f1e10.
//
// Solidity: function cancel((address,address,(uint8,address,uint256,uint256,uint256)[],(uint8,address,uint256,uint256,uint256,address)[],uint8,uint256,uint256,bytes32,uint256,bytes32,uint256)[] orders) returns(bool cancelled)
func (_Seaport *SeaportTransactorSession) Cancel(orders []OrderComponents) (*types.Transaction, error) {
	return _Seaport.Contract.Cancel(&_Seaport.TransactOpts, orders)
}

// FulfillAdvancedOrder is a paid mutator transaction binding the contract method 0xe7acab24.
//
// Solidity: function fulfillAdvancedOrder(((address,address,(uint8,address,uint256,uint256,uint256)[],(uint8,address,uint256,uint256,uint256,address)[],uint8,uint256,uint256,bytes32,uint256,bytes32,uint256),uint120,uint120,bytes,bytes) , (uint256,uint8,uint256,uint256,bytes32[])[] , bytes32 fulfillerConduitKey, address recipient) payable returns(bool fulfilled)
func (_Seaport *SeaportTransactor) FulfillAdvancedOrder(opts *bind.TransactOpts, arg0 AdvancedOrder, arg1 []CriteriaResolver, fulfillerConduitKey [32]byte, recipient common.Address) (*types.Transaction, error) {
	return _Seaport.contract.Transact(opts, "fulfillAdvancedOrder", arg0, arg1, fulfillerConduitKey, recipient)
}

// FulfillAdvancedOrder is a paid mutator transaction binding the contract method 0xe7acab24.
//
// Solidity: function fulfillAdvancedOrder(((address,address,(uint8,address,uint256,uint256,uint256)[],(uint8,address,uint256,uint256,uint256,address)[],uint8,uint256,uint256,bytes32,uint256,bytes32,uint256),uint120,uint120,bytes,bytes) , (uint256,uint8,uint256,uint256,bytes32[])[] , bytes32 fulfillerConduitKey, address recipient) payable returns(bool fulfilled)
func (_Seaport *SeaportSession) FulfillAdvancedOrder(arg0 AdvancedOrder, arg1 []CriteriaResolver, fulfillerConduitKey [32]byte, recipient common.Address) (*types.Transaction, error) {
	return _Seaport.Contract.FulfillAdvancedOrder(&_Seaport.TransactOpts, arg0, arg1, fulfillerConduitKey, recipient)
}

// FulfillAdvancedOrder is a paid mutator transaction binding the contract method 0xe7acab24.
//
// Solidity: function fulfillAdvancedOrder(((address,address,(uint8,address,uint256,uint256,uint256)[],(uint8,address,uint256,uint256,uint256,address)[],uint8,uint256,uint256,bytes32,uint256,bytes32,uint256),uint120,uint120,bytes,bytes) , (uint256,uint8,uint256,uint256,bytes32[])[] , bytes32 fulfillerConduitKey, address recipient) payable returns(bool fulfilled)
func (_Seaport *SeaportTransactorSession) FulfillAdvancedOrder(arg0 AdvancedOrder, arg1 []CriteriaResolver, fulfillerConduitKey [32]byte, recipient common.Address) (*types.Transaction, error) {
	return _Seaport.Contract.FulfillAdvancedOrder(&_Seaport.TransactOpts, arg0, arg1, fulfillerConduitKey, recipient)
}

// FulfillAvailableAdvancedOrders is a paid mutator transaction binding the contract method 0x87201b41.
//
// Solidity: function fulfillAvailableAdvancedOrders(((address,address,(uint8,address,uint256,uint256,uint256)[],(uint8,address,uint256,uint256,uint256,address)[],uint8,uint256,uint256,bytes32,uint256,bytes32,uint256),uint120,uint120,bytes,bytes)[] , (uint256,uint8,uint256,uint256,bytes32[])[] , (uint256,uint256)[][] , (uint256,uint256)[][] , bytes32 fulfillerConduitKey, address recipient, uint256 maximumFulfilled) payable returns(bool[], ((uint8,address,uint256,uint256,address),address,bytes32)[])
func (_Seaport *SeaportTransactor) FulfillAvailableAdvancedOrders(opts *bind.TransactOpts, arg0 []AdvancedOrder, arg1 []CriteriaResolver, arg2 [][]FulfillmentComponent, arg3 [][]FulfillmentComponent, fulfillerConduitKey [32]byte, recipient common.Address, maximumFulfilled *big.Int) (*types.Transaction, error) {
	return _Seaport.contract.Transact(opts, "fulfillAvailableAdvancedOrders", arg0, arg1, arg2, arg3, fulfillerConduitKey, recipient, maximumFulfilled)
}

// FulfillAvailableAdvancedOrders is a paid mutator transaction binding the contract method 0x87201b41.
//
// Solidity: function fulfillAvailableAdvancedOrders(((address,address,(uint8,address,uint256,uint256,uint256)[],(uint8,address,uint256,uint256,uint256,address)[],uint8,uint256,uint256,bytes32,uint256,bytes32,uint256),uint120,uint120,bytes,bytes)[] , (uint256,uint8,uint256,uint256,bytes32[])[] , (uint256,uint256)[][] , (uint256,uint256)[][] , bytes32 fulfillerConduitKey, address recipient, uint256 maximumFulfilled) payable returns(bool[], ((uint8,address,uint256,uint256,address),address,bytes32)[])
func (_Seaport *SeaportSession) FulfillAvailableAdvancedOrders(arg0 []AdvancedOrder, arg1 []CriteriaResolver, arg2 [][]FulfillmentComponent, arg3 [][]FulfillmentComponent, fulfillerConduitKey [32]byte, recipient common.Address, maximumFulfilled *big.Int) (*types.Transaction, error) {
	return _Seaport.Contract.FulfillAvailableAdvancedOrders(&_Seaport.TransactOpts, arg0, arg1, arg2, arg3, fulfillerConduitKey, recipient, maximumFulfilled)
}

// FulfillAvailableAdvancedOrders is a paid mutator transaction binding the contract method 0x87201b41.
//
// Solidity: function fulfillAvailableAdvancedOrders(((address,address,(uint8,address,uint256,uint256,uint256)[],(uint8,address,uint256,uint256,uint256,address)[],uint8,uint256,uint256,bytes32,uint256,bytes32,uint256),uint120,uint120,bytes,bytes)[] , (uint256,uint8,uint256,uint256,bytes32[])[] , (uint256,uint256)[][] , (uint256,uint256)[][] , bytes32 fulfillerConduitKey, address recipient, uint256 maximumFulfilled) payable returns(bool[], ((uint8,address,uint256,uint256,address),address,bytes32)[])
func (_Seaport *SeaportTransactorSession) FulfillAvailableAdvancedOrders(arg0 []AdvancedOrder, arg1 []CriteriaResolver, arg2 [][]FulfillmentComponent, arg3 [][]FulfillmentComponent, fulfillerConduitKey [32]byte, recipient common.Address, maximumFulfilled *big.Int) (*types.Transaction, error) {
	return _Seaport.Contract.FulfillAvailableAdvancedOrders(&_Seaport.TransactOpts, arg0, arg1, arg2, arg3, fulfillerConduitKey, recipient, maximumFulfilled)
}

// FulfillAvailableOrders is a paid mutator transaction binding the contract method 0xed98a574.
//
// Solidity: function fulfillAvailableOrders(((address,address,(uint8,address,uint256,uint256,uint256)[],(uint8,address,uint256,uint256,uint256,address)[],uint8,uint256,uint256,bytes32,uint256,bytes32,uint256),bytes)[] , (uint256,uint256)[][] , (uint256,uint256)[][] , bytes32 fulfillerConduitKey, uint256 maximumFulfilled) payable returns(bool[], ((uint8,address,uint256,uint256,address),address,bytes32)[])
func (_Seaport *SeaportTransactor) FulfillAvailableOrders(opts *bind.TransactOpts, arg0 []Order, arg1 [][]FulfillmentComponent, arg2 [][]FulfillmentComponent, fulfillerConduitKey [32]byte, maximumFulfilled *big.Int) (*types.Transaction, error) {
	return _Seaport.contract.Transact(opts, "fulfillAvailableOrders", arg0, arg1, arg2, fulfillerConduitKey, maximumFulfilled)
}

// FulfillAvailableOrders is a paid mutator transaction binding the contract method 0xed98a574.
//
// Solidity: function fulfillAvailableOrders(((address,address,(uint8,address,uint256,uint256,uint256)[],(uint8,address,uint256,uint256,uint256,address)[],uint8,uint256,uint256,bytes32,uint256,bytes32,uint256),bytes)[] , (uint256,uint256)[][] , (uint256,uint256)[][] , bytes32 fulfillerConduitKey, uint256 maximumFulfilled) payable returns(bool[], ((uint8,address,uint256,uint256,address),address,bytes32)[])
func (_Seaport *SeaportSession) FulfillAvailableOrders(arg0 []Order, arg1 [][]FulfillmentComponent, arg2 [][]FulfillmentComponent, fulfillerConduitKey [32]byte, maximumFulfilled *big.Int) (*types.Transaction, error) {
	return _Seaport.Contract.FulfillAvailableOrders(&_Seaport.TransactOpts, arg0, arg1, arg2, fulfillerConduitKey, maximumFulfilled)
}

// FulfillAvailableOrders is a paid mutator transaction binding the contract method 0xed98a574.
//
// Solidity: function fulfillAvailableOrders(((address,address,(uint8,address,uint256,uint256,uint256)[],(uint8,address,uint256,uint256,uint256,address)[],uint8,uint256,uint256,bytes32,uint256,bytes32,uint256),bytes)[] , (uint256,uint256)[][] , (uint256,uint256)[][] , bytes32 fulfillerConduitKey, uint256 maximumFulfilled) payable returns(bool[], ((uint8,address,uint256,uint256,address),address,bytes32)[])
func (_Seaport *SeaportTransactorSession) FulfillAvailableOrders(arg0 []Order, arg1 [][]FulfillmentComponent, arg2 [][]FulfillmentComponent, fulfillerConduitKey [32]byte, maximumFulfilled *big.Int) (*types.Transaction, error) {
	return _Seaport.Contract.FulfillAvailableOrders(&_Seaport.TransactOpts, arg0, arg1, arg2, fulfillerConduitKey, maximumFulfilled)
}

// FulfillBasicOrder is a paid mutator transaction binding the contract method 0xfb0f3ee1.
//
// Solidity: function fulfillBasicOrder((address,uint256,uint256,address,address,address,uint256,uint256,uint8,uint256,uint256,bytes32,uint256,bytes32,bytes32,uint256,(uint256,address)[],bytes) ) payable returns(bool fulfilled)
func (_Seaport *SeaportTransactor) FulfillBasicOrder(opts *bind.TransactOpts, arg0 BasicOrderParameters) (*types.Transaction, error) {
	return _Seaport.contract.Transact(opts, "fulfillBasicOrder", arg0)
}

// FulfillBasicOrder is a paid mutator transaction binding the contract method 0xfb0f3ee1.
//
// Solidity: function fulfillBasicOrder((address,uint256,uint256,address,address,address,uint256,uint256,uint8,uint256,uint256,bytes32,uint256,bytes32,bytes32,uint256,(uint256,address)[],bytes) ) payable returns(bool fulfilled)
func (_Seaport *SeaportSession) FulfillBasicOrder(arg0 BasicOrderParameters) (*types.Transaction, error) {
	return _Seaport.Contract.FulfillBasicOrder(&_Seaport.TransactOpts, arg0)
}

// FulfillBasicOrder is a paid mutator transaction binding the contract method 0xfb0f3ee1.
//
// Solidity: function fulfillBasicOrder((address,uint256,uint256,address,address,address,uint256,uint256,uint8,uint256,uint256,bytes32,uint256,bytes32,bytes32,uint256,(uint256,address)[],bytes) ) payable returns(bool fulfilled)
func (_Seaport *SeaportTransactorSession) FulfillBasicOrder(arg0 BasicOrderParameters) (*types.Transaction, error) {
	return _Seaport.Contract.FulfillBasicOrder(&_Seaport.TransactOpts, arg0)
}

// FulfillBasicOrderEfficient6GL6yc is a paid mutator transaction binding the contract method 0x00000000.
//
// Solidity: function fulfillBasicOrder_efficient_6GL6yc((address,uint256,uint256,address,address,address,uint256,uint256,uint8,uint256,uint256,bytes32,uint256,bytes32,bytes32,uint256,(uint256,address)[],bytes) ) payable returns(bool fulfilled)
func (_Seaport *SeaportTransactor) FulfillBasicOrderEfficient6GL6yc(opts *bind.TransactOpts, arg0 BasicOrderParameters) (*types.Transaction, error) {
	return _Seaport.contract.Transact(opts, "fulfillBasicOrder_efficient_6GL6yc", arg0)
}

// FulfillBasicOrderEfficient6GL6yc is a paid mutator transaction binding the contract method 0x00000000.
//
// Solidity: function fulfillBasicOrder_efficient_6GL6yc((address,uint256,uint256,address,address,address,uint256,uint256,uint8,uint256,uint256,bytes32,uint256,bytes32,bytes32,uint256,(uint256,address)[],bytes) ) payable returns(bool fulfilled)
func (_Seaport *SeaportSession) FulfillBasicOrderEfficient6GL6yc(arg0 BasicOrderParameters) (*types.Transaction, error) {
	return _Seaport.Contract.FulfillBasicOrderEfficient6GL6yc(&_Seaport.TransactOpts, arg0)
}

// FulfillBasicOrderEfficient6GL6yc is a paid mutator transaction binding the contract method 0x00000000.
//
// Solidity: function fulfillBasicOrder_efficient_6GL6yc((address,uint256,uint256,address,address,address,uint256,uint256,uint8,uint256,uint256,bytes32,uint256,bytes32,bytes32,uint256,(uint256,address)[],bytes) ) payable returns(bool fulfilled)
func (_Seaport *SeaportTransactorSession) FulfillBasicOrderEfficient6GL6yc(arg0 BasicOrderParameters) (*types.Transaction, error) {
	return _Seaport.Contract.FulfillBasicOrderEfficient6GL6yc(&_Seaport.TransactOpts, arg0)
}

// FulfillOrder is a paid mutator transaction binding the contract method 0xb3a34c4c.
//
// Solidity: function fulfillOrder(((address,address,(uint8,address,uint256,uint256,uint256)[],(uint8,address,uint256,uint256,uint256,address)[],uint8,uint256,uint256,bytes32,uint256,bytes32,uint256),bytes) , bytes32 fulfillerConduitKey) payable returns(bool fulfilled)
func (_Seaport *SeaportTransactor) FulfillOrder(opts *bind.TransactOpts, arg0 Order, fulfillerConduitKey [32]byte) (*types.Transaction, error) {
	return _Seaport.contract.Transact(opts, "fulfillOrder", arg0, fulfillerConduitKey)
}

// FulfillOrder is a paid mutator transaction binding the contract method 0xb3a34c4c.
//
// Solidity: function fulfillOrder(((address,address,(uint8,address,uint256,uint256,uint256)[],(uint8,address,uint256,uint256,uint256,address)[],uint8,uint256,uint256,bytes32,uint256,bytes32,uint256),bytes) , bytes32 fulfillerConduitKey) payable returns(bool fulfilled)
func (_Seaport *SeaportSession) FulfillOrder(arg0 Order, fulfillerConduitKey [32]byte) (*types.Transaction, error) {
	return _Seaport.Contract.FulfillOrder(&_Seaport.TransactOpts, arg0, fulfillerConduitKey)
}

// FulfillOrder is a paid mutator transaction binding the contract method 0xb3a34c4c.
//
// Solidity: function fulfillOrder(((address,address,(uint8,address,uint256,uint256,uint256)[],(uint8,address,uint256,uint256,uint256,address)[],uint8,uint256,uint256,bytes32,uint256,bytes32,uint256),bytes) , bytes32 fulfillerConduitKey) payable returns(bool fulfilled)
func (_Seaport *SeaportTransactorSession) FulfillOrder(arg0 Order, fulfillerConduitKey [32]byte) (*types.Transaction, error) {
	return _Seaport.Contract.FulfillOrder(&_Seaport.TransactOpts, arg0, fulfillerConduitKey)
}

// IncrementCounter is a paid mutator transaction binding the contract method 0x5b34b966.
//
// Solidity: function incrementCounter() returns(uint256 newCounter)
func (_Seaport *SeaportTransactor) IncrementCounter(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Seaport.contract.Transact(opts, "incrementCounter")
}

// IncrementCounter is a paid mutator transaction binding the contract method 0x5b34b966.
//
// Solidity: function incrementCounter() returns(uint256 newCounter)
func (_Seaport *SeaportSession) IncrementCounter() (*types.Transaction, error) {
	return _Seaport.Contract.IncrementCounter(&_Seaport.TransactOpts)
}

// IncrementCounter is a paid mutator transaction binding the contract method 0x5b34b966.
//
// Solidity: function incrementCounter() returns(uint256 newCounter)
func (_Seaport *SeaportTransactorSession) IncrementCounter() (*types.Transaction, error) {
	return _Seaport.Contract.IncrementCounter(&_Seaport.TransactOpts)
}

// MatchAdvancedOrders is a paid mutator transaction binding the contract method 0xf2d12b12.
//
// Solidity: function matchAdvancedOrders(((address,address,(uint8,address,uint256,uint256,uint256)[],(uint8,address,uint256,uint256,uint256,address)[],uint8,uint256,uint256,bytes32,uint256,bytes32,uint256),uint120,uint120,bytes,bytes)[] , (uint256,uint8,uint256,uint256,bytes32[])[] , ((uint256,uint256)[],(uint256,uint256)[])[] , address recipient) payable returns(((uint8,address,uint256,uint256,address),address,bytes32)[])
func (_Seaport *SeaportTransactor) MatchAdvancedOrders(opts *bind.TransactOpts, arg0 []AdvancedOrder, arg1 []CriteriaResolver, arg2 []Fulfillment, recipient common.Address) (*types.Transaction, error) {
	return _Seaport.contract.Transact(opts, "matchAdvancedOrders", arg0, arg1, arg2, recipient)
}

// MatchAdvancedOrders is a paid mutator transaction binding the contract method 0xf2d12b12.
//
// Solidity: function matchAdvancedOrders(((address,address,(uint8,address,uint256,uint256,uint256)[],(uint8,address,uint256,uint256,uint256,address)[],uint8,uint256,uint256,bytes32,uint256,bytes32,uint256),uint120,uint120,bytes,bytes)[] , (uint256,uint8,uint256,uint256,bytes32[])[] , ((uint256,uint256)[],(uint256,uint256)[])[] , address recipient) payable returns(((uint8,address,uint256,uint256,address),address,bytes32)[])
func (_Seaport *SeaportSession) MatchAdvancedOrders(arg0 []AdvancedOrder, arg1 []CriteriaResolver, arg2 []Fulfillment, recipient common.Address) (*types.Transaction, error) {
	return _Seaport.Contract.MatchAdvancedOrders(&_Seaport.TransactOpts, arg0, arg1, arg2, recipient)
}

// MatchAdvancedOrders is a paid mutator transaction binding the contract method 0xf2d12b12.
//
// Solidity: function matchAdvancedOrders(((address,address,(uint8,address,uint256,uint256,uint256)[],(uint8,address,uint256,uint256,uint256,address)[],uint8,uint256,uint256,bytes32,uint256,bytes32,uint256),uint120,uint120,bytes,bytes)[] , (uint256,uint8,uint256,uint256,bytes32[])[] , ((uint256,uint256)[],(uint256,uint256)[])[] , address recipient) payable returns(((uint8,address,uint256,uint256,address),address,bytes32)[])
func (_Seaport *SeaportTransactorSession) MatchAdvancedOrders(arg0 []AdvancedOrder, arg1 []CriteriaResolver, arg2 []Fulfillment, recipient common.Address) (*types.Transaction, error) {
	return _Seaport.Contract.MatchAdvancedOrders(&_Seaport.TransactOpts, arg0, arg1, arg2, recipient)
}

// MatchOrders is a paid mutator transaction binding the contract method 0xa8174404.
//
// Solidity: function matchOrders(((address,address,(uint8,address,uint256,uint256,uint256)[],(uint8,address,uint256,uint256,uint256,address)[],uint8,uint256,uint256,bytes32,uint256,bytes32,uint256),bytes)[] , ((uint256,uint256)[],(uint256,uint256)[])[] ) payable returns(((uint8,address,uint256,uint256,address),address,bytes32)[])
func (_Seaport *SeaportTransactor) MatchOrders(opts *bind.TransactOpts, arg0 []Order, arg1 []Fulfillment) (*types.Transaction, error) {
	return _Seaport.contract.Transact(opts, "matchOrders", arg0, arg1)
}

// MatchOrders is a paid mutator transaction binding the contract method 0xa8174404.
//
// Solidity: function matchOrders(((address,address,(uint8,address,uint256,uint256,uint256)[],(uint8,address,uint256,uint256,uint256,address)[],uint8,uint256,uint256,bytes32,uint256,bytes32,uint256),bytes)[] , ((uint256,uint256)[],(uint256,uint256)[])[] ) payable returns(((uint8,address,uint256,uint256,address),address,bytes32)[])
func (_Seaport *SeaportSession) MatchOrders(arg0 []Order, arg1 []Fulfillment) (*types.Transaction, error) {
	return _Seaport.Contract.MatchOrders(&_Seaport.TransactOpts, arg0, arg1)
}

// MatchOrders is a paid mutator transaction binding the contract method 0xa8174404.
//
// Solidity: function matchOrders(((address,address,(uint8,address,uint256,uint256,uint256)[],(uint8,address,uint256,uint256,uint256,address)[],uint8,uint256,uint256,bytes32,uint256,bytes32,uint256),bytes)[] , ((uint256,uint256)[],(uint256,uint256)[])[] ) payable returns(((uint8,address,uint256,uint256,address),address,bytes32)[])
func (_Seaport *SeaportTransactorSession) MatchOrders(arg0 []Order, arg1 []Fulfillment) (*types.Transaction, error) {
	return _Seaport.Contract.MatchOrders(&_Seaport.TransactOpts, arg0, arg1)
}

// Validate is a paid mutator transaction binding the contract method 0x88147732.
//
// Solidity: function validate(((address,address,(uint8,address,uint256,uint256,uint256)[],(uint8,address,uint256,uint256,uint256,address)[],uint8,uint256,uint256,bytes32,uint256,bytes32,uint256),bytes)[] ) returns(bool)
func (_Seaport *SeaportTransactor) Validate(opts *bind.TransactOpts, arg0 []Order) (*types.Transaction, error) {
	return _Seaport.contract.Transact(opts, "validate", arg0)
}

// Validate is a paid mutator transaction binding the contract method 0x88147732.
//
// Solidity: function validate(((address,address,(uint8,address,uint256,uint256,uint256)[],(uint8,address,uint256,uint256,uint256,address)[],uint8,uint256,uint256,bytes32,uint256,bytes32,uint256),bytes)[] ) returns(bool)
func (_Seaport *SeaportSession) Validate(arg0 []Order) (*types.Transaction, error) {
	return _Seaport.Contract.Validate(&_Seaport.TransactOpts, arg0)
}

// Validate is a paid mutator transaction binding the contract method 0x88147732.
//
// Solidity: function validate(((address,address,(uint8,address,uint256,uint256,uint256)[],(uint8,address,uint256,uint256,uint256,address)[],uint8,uint256,uint256,bytes32,uint256,bytes32,uint256),bytes)[] ) returns(bool)
func (_Seaport *SeaportTransactorSession) Validate(arg0 []Order) (*types.Transaction, error) {
	return _Seaport.Contract.Validate(&_Seaport.TransactOpts, arg0)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Seaport *SeaportTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Seaport.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Seaport *SeaportSession) Receive() (*types.Transaction, error) {
	return _Seaport.Contract.Receive(&_Seaport.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Seaport *SeaportTransactorSession) Receive() (*types.Transaction, error) {
	return _Seaport.Contract.Receive(&_Seaport.TransactOpts)
}

// SeaportCounterIncrementedIterator is returned from FilterCounterIncremented and is used to iterate over the raw logs and unpacked data for CounterIncremented events raised by the Seaport contract.
type SeaportCounterIncrementedIterator struct {
	Event *SeaportCounterIncremented // Event containing the contract specifics and raw log

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
func (it *SeaportCounterIncrementedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SeaportCounterIncremented)
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
		it.Event = new(SeaportCounterIncremented)
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
func (it *SeaportCounterIncrementedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SeaportCounterIncrementedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SeaportCounterIncremented represents a CounterIncremented event raised by the Seaport contract.
type SeaportCounterIncremented struct {
	NewCounter *big.Int
	Offerer    common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterCounterIncremented is a free log retrieval operation binding the contract event 0x721c20121297512b72821b97f5326877ea8ecf4bb9948fea5bfcb6453074d37f.
//
// Solidity: event CounterIncremented(uint256 newCounter, address indexed offerer)
func (_Seaport *SeaportFilterer) FilterCounterIncremented(opts *bind.FilterOpts, offerer []common.Address) (*SeaportCounterIncrementedIterator, error) {

	var offererRule []interface{}
	for _, offererItem := range offerer {
		offererRule = append(offererRule, offererItem)
	}

	logs, sub, err := _Seaport.contract.FilterLogs(opts, "CounterIncremented", offererRule)
	if err != nil {
		return nil, err
	}
	return &SeaportCounterIncrementedIterator{contract: _Seaport.contract, event: "CounterIncremented", logs: logs, sub: sub}, nil
}

// WatchCounterIncremented is a free log subscription operation binding the contract event 0x721c20121297512b72821b97f5326877ea8ecf4bb9948fea5bfcb6453074d37f.
//
// Solidity: event CounterIncremented(uint256 newCounter, address indexed offerer)
func (_Seaport *SeaportFilterer) WatchCounterIncremented(opts *bind.WatchOpts, sink chan<- *SeaportCounterIncremented, offerer []common.Address) (event.Subscription, error) {

	var offererRule []interface{}
	for _, offererItem := range offerer {
		offererRule = append(offererRule, offererItem)
	}

	logs, sub, err := _Seaport.contract.WatchLogs(opts, "CounterIncremented", offererRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SeaportCounterIncremented)
				if err := _Seaport.contract.UnpackLog(event, "CounterIncremented", log); err != nil {
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

// ParseCounterIncremented is a log parse operation binding the contract event 0x721c20121297512b72821b97f5326877ea8ecf4bb9948fea5bfcb6453074d37f.
//
// Solidity: event CounterIncremented(uint256 newCounter, address indexed offerer)
func (_Seaport *SeaportFilterer) ParseCounterIncremented(log types.Log) (*SeaportCounterIncremented, error) {
	event := new(SeaportCounterIncremented)
	if err := _Seaport.contract.UnpackLog(event, "CounterIncremented", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SeaportOrderCancelledIterator is returned from FilterOrderCancelled and is used to iterate over the raw logs and unpacked data for OrderCancelled events raised by the Seaport contract.
type SeaportOrderCancelledIterator struct {
	Event *SeaportOrderCancelled // Event containing the contract specifics and raw log

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
func (it *SeaportOrderCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SeaportOrderCancelled)
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
		it.Event = new(SeaportOrderCancelled)
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
func (it *SeaportOrderCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SeaportOrderCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SeaportOrderCancelled represents a OrderCancelled event raised by the Seaport contract.
type SeaportOrderCancelled struct {
	OrderHash [32]byte
	Offerer   common.Address
	Zone      common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterOrderCancelled is a free log retrieval operation binding the contract event 0x6bacc01dbe442496068f7d234edd811f1a5f833243e0aec824f86ab861f3c90d.
//
// Solidity: event OrderCancelled(bytes32 orderHash, address indexed offerer, address indexed zone)
func (_Seaport *SeaportFilterer) FilterOrderCancelled(opts *bind.FilterOpts, offerer []common.Address, zone []common.Address) (*SeaportOrderCancelledIterator, error) {

	var offererRule []interface{}
	for _, offererItem := range offerer {
		offererRule = append(offererRule, offererItem)
	}
	var zoneRule []interface{}
	for _, zoneItem := range zone {
		zoneRule = append(zoneRule, zoneItem)
	}

	logs, sub, err := _Seaport.contract.FilterLogs(opts, "OrderCancelled", offererRule, zoneRule)
	if err != nil {
		return nil, err
	}
	return &SeaportOrderCancelledIterator{contract: _Seaport.contract, event: "OrderCancelled", logs: logs, sub: sub}, nil
}

// WatchOrderCancelled is a free log subscription operation binding the contract event 0x6bacc01dbe442496068f7d234edd811f1a5f833243e0aec824f86ab861f3c90d.
//
// Solidity: event OrderCancelled(bytes32 orderHash, address indexed offerer, address indexed zone)
func (_Seaport *SeaportFilterer) WatchOrderCancelled(opts *bind.WatchOpts, sink chan<- *SeaportOrderCancelled, offerer []common.Address, zone []common.Address) (event.Subscription, error) {

	var offererRule []interface{}
	for _, offererItem := range offerer {
		offererRule = append(offererRule, offererItem)
	}
	var zoneRule []interface{}
	for _, zoneItem := range zone {
		zoneRule = append(zoneRule, zoneItem)
	}

	logs, sub, err := _Seaport.contract.WatchLogs(opts, "OrderCancelled", offererRule, zoneRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SeaportOrderCancelled)
				if err := _Seaport.contract.UnpackLog(event, "OrderCancelled", log); err != nil {
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

// ParseOrderCancelled is a log parse operation binding the contract event 0x6bacc01dbe442496068f7d234edd811f1a5f833243e0aec824f86ab861f3c90d.
//
// Solidity: event OrderCancelled(bytes32 orderHash, address indexed offerer, address indexed zone)
func (_Seaport *SeaportFilterer) ParseOrderCancelled(log types.Log) (*SeaportOrderCancelled, error) {
	event := new(SeaportOrderCancelled)
	if err := _Seaport.contract.UnpackLog(event, "OrderCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SeaportOrderFulfilledIterator is returned from FilterOrderFulfilled and is used to iterate over the raw logs and unpacked data for OrderFulfilled events raised by the Seaport contract.
type SeaportOrderFulfilledIterator struct {
	Event *SeaportOrderFulfilled // Event containing the contract specifics and raw log

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
func (it *SeaportOrderFulfilledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SeaportOrderFulfilled)
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
		it.Event = new(SeaportOrderFulfilled)
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
func (it *SeaportOrderFulfilledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SeaportOrderFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SeaportOrderFulfilled represents a OrderFulfilled event raised by the Seaport contract.
type SeaportOrderFulfilled struct {
	OrderHash     [32]byte
	Offerer       common.Address
	Zone          common.Address
	Recipient     common.Address
	Offer         []SpentItem
	Consideration []ReceivedItem
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOrderFulfilled is a free log retrieval operation binding the contract event 0x9d9af8e38d66c62e2c12f0225249fd9d721c54b83f48d9352c97c6cacdcb6f31.
//
// Solidity: event OrderFulfilled(bytes32 orderHash, address indexed offerer, address indexed zone, address recipient, (uint8,address,uint256,uint256)[] offer, (uint8,address,uint256,uint256,address)[] consideration)
func (_Seaport *SeaportFilterer) FilterOrderFulfilled(opts *bind.FilterOpts, offerer []common.Address, zone []common.Address) (*SeaportOrderFulfilledIterator, error) {

	var offererRule []interface{}
	for _, offererItem := range offerer {
		offererRule = append(offererRule, offererItem)
	}
	var zoneRule []interface{}
	for _, zoneItem := range zone {
		zoneRule = append(zoneRule, zoneItem)
	}

	logs, sub, err := _Seaport.contract.FilterLogs(opts, "OrderFulfilled", offererRule, zoneRule)
	if err != nil {
		return nil, err
	}
	return &SeaportOrderFulfilledIterator{contract: _Seaport.contract, event: "OrderFulfilled", logs: logs, sub: sub}, nil
}

// WatchOrderFulfilled is a free log subscription operation binding the contract event 0x9d9af8e38d66c62e2c12f0225249fd9d721c54b83f48d9352c97c6cacdcb6f31.
//
// Solidity: event OrderFulfilled(bytes32 orderHash, address indexed offerer, address indexed zone, address recipient, (uint8,address,uint256,uint256)[] offer, (uint8,address,uint256,uint256,address)[] consideration)
func (_Seaport *SeaportFilterer) WatchOrderFulfilled(opts *bind.WatchOpts, sink chan<- *SeaportOrderFulfilled, offerer []common.Address, zone []common.Address) (event.Subscription, error) {

	var offererRule []interface{}
	for _, offererItem := range offerer {
		offererRule = append(offererRule, offererItem)
	}
	var zoneRule []interface{}
	for _, zoneItem := range zone {
		zoneRule = append(zoneRule, zoneItem)
	}

	logs, sub, err := _Seaport.contract.WatchLogs(opts, "OrderFulfilled", offererRule, zoneRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SeaportOrderFulfilled)
				if err := _Seaport.contract.UnpackLog(event, "OrderFulfilled", log); err != nil {
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

// ParseOrderFulfilled is a log parse operation binding the contract event 0x9d9af8e38d66c62e2c12f0225249fd9d721c54b83f48d9352c97c6cacdcb6f31.
//
// Solidity: event OrderFulfilled(bytes32 orderHash, address indexed offerer, address indexed zone, address recipient, (uint8,address,uint256,uint256)[] offer, (uint8,address,uint256,uint256,address)[] consideration)
func (_Seaport *SeaportFilterer) ParseOrderFulfilled(log types.Log) (*SeaportOrderFulfilled, error) {
	event := new(SeaportOrderFulfilled)
	if err := _Seaport.contract.UnpackLog(event, "OrderFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SeaportOrderValidatedIterator is returned from FilterOrderValidated and is used to iterate over the raw logs and unpacked data for OrderValidated events raised by the Seaport contract.
type SeaportOrderValidatedIterator struct {
	Event *SeaportOrderValidated // Event containing the contract specifics and raw log

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
func (it *SeaportOrderValidatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SeaportOrderValidated)
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
		it.Event = new(SeaportOrderValidated)
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
func (it *SeaportOrderValidatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SeaportOrderValidatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SeaportOrderValidated represents a OrderValidated event raised by the Seaport contract.
type SeaportOrderValidated struct {
	OrderHash       [32]byte
	OrderParameters OrderParameters
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterOrderValidated is a free log retrieval operation binding the contract event 0xf280791efe782edcf06ce15c8f4dff17601db3b88eb3805a0db7d77faf757f04.
//
// Solidity: event OrderValidated(bytes32 orderHash, (address,address,(uint8,address,uint256,uint256,uint256)[],(uint8,address,uint256,uint256,uint256,address)[],uint8,uint256,uint256,bytes32,uint256,bytes32,uint256) orderParameters)
func (_Seaport *SeaportFilterer) FilterOrderValidated(opts *bind.FilterOpts) (*SeaportOrderValidatedIterator, error) {

	logs, sub, err := _Seaport.contract.FilterLogs(opts, "OrderValidated")
	if err != nil {
		return nil, err
	}
	return &SeaportOrderValidatedIterator{contract: _Seaport.contract, event: "OrderValidated", logs: logs, sub: sub}, nil
}

// WatchOrderValidated is a free log subscription operation binding the contract event 0xf280791efe782edcf06ce15c8f4dff17601db3b88eb3805a0db7d77faf757f04.
//
// Solidity: event OrderValidated(bytes32 orderHash, (address,address,(uint8,address,uint256,uint256,uint256)[],(uint8,address,uint256,uint256,uint256,address)[],uint8,uint256,uint256,bytes32,uint256,bytes32,uint256) orderParameters)
func (_Seaport *SeaportFilterer) WatchOrderValidated(opts *bind.WatchOpts, sink chan<- *SeaportOrderValidated) (event.Subscription, error) {

	logs, sub, err := _Seaport.contract.WatchLogs(opts, "OrderValidated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SeaportOrderValidated)
				if err := _Seaport.contract.UnpackLog(event, "OrderValidated", log); err != nil {
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

// ParseOrderValidated is a log parse operation binding the contract event 0xf280791efe782edcf06ce15c8f4dff17601db3b88eb3805a0db7d77faf757f04.
//
// Solidity: event OrderValidated(bytes32 orderHash, (address,address,(uint8,address,uint256,uint256,uint256)[],(uint8,address,uint256,uint256,uint256,address)[],uint8,uint256,uint256,bytes32,uint256,bytes32,uint256) orderParameters)
func (_Seaport *SeaportFilterer) ParseOrderValidated(log types.Log) (*SeaportOrderValidated, error) {
	event := new(SeaportOrderValidated)
	if err := _Seaport.contract.UnpackLog(event, "OrderValidated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SeaportOrdersMatchedIterator is returned from FilterOrdersMatched and is used to iterate over the raw logs and unpacked data for OrdersMatched events raised by the Seaport contract.
type SeaportOrdersMatchedIterator struct {
	Event *SeaportOrdersMatched // Event containing the contract specifics and raw log

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
func (it *SeaportOrdersMatchedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SeaportOrdersMatched)
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
		it.Event = new(SeaportOrdersMatched)
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
func (it *SeaportOrdersMatchedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SeaportOrdersMatchedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SeaportOrdersMatched represents a OrdersMatched event raised by the Seaport contract.
type SeaportOrdersMatched struct {
	OrderHashes [][32]byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterOrdersMatched is a free log retrieval operation binding the contract event 0x4b9f2d36e1b4c93de62cc077b00b1a91d84b6c31b4a14e012718dcca230689e7.
//
// Solidity: event OrdersMatched(bytes32[] orderHashes)
func (_Seaport *SeaportFilterer) FilterOrdersMatched(opts *bind.FilterOpts) (*SeaportOrdersMatchedIterator, error) {

	logs, sub, err := _Seaport.contract.FilterLogs(opts, "OrdersMatched")
	if err != nil {
		return nil, err
	}
	return &SeaportOrdersMatchedIterator{contract: _Seaport.contract, event: "OrdersMatched", logs: logs, sub: sub}, nil
}

// WatchOrdersMatched is a free log subscription operation binding the contract event 0x4b9f2d36e1b4c93de62cc077b00b1a91d84b6c31b4a14e012718dcca230689e7.
//
// Solidity: event OrdersMatched(bytes32[] orderHashes)
func (_Seaport *SeaportFilterer) WatchOrdersMatched(opts *bind.WatchOpts, sink chan<- *SeaportOrdersMatched) (event.Subscription, error) {

	logs, sub, err := _Seaport.contract.WatchLogs(opts, "OrdersMatched")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SeaportOrdersMatched)
				if err := _Seaport.contract.UnpackLog(event, "OrdersMatched", log); err != nil {
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

// ParseOrdersMatched is a log parse operation binding the contract event 0x4b9f2d36e1b4c93de62cc077b00b1a91d84b6c31b4a14e012718dcca230689e7.
//
// Solidity: event OrdersMatched(bytes32[] orderHashes)
func (_Seaport *SeaportFilterer) ParseOrdersMatched(log types.Log) (*SeaportOrdersMatched, error) {
	event := new(SeaportOrdersMatched)
	if err := _Seaport.contract.UnpackLog(event, "OrdersMatched", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
