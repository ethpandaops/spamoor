package seaporttrades

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethpandaops/spamoor/scenarios/seaport-trades/contract"
)

// EIP-712 type hashes used by Seaport 1.x. The type strings are the exact ones
// Seaport hashes internally, so an order hash computed here matches what
// Seaport.getOrderHash returns on-chain (validated once at startup). The
// referenced struct types follow EIP-712 ordering (ConsiderationItem before
// OfferItem, alphabetical) appended after the OrderComponents definition.
var (
	offerItemTypehash = crypto.Keccak256Hash([]byte(
		"OfferItem(uint8 itemType,address token,uint256 identifierOrCriteria,uint256 startAmount,uint256 endAmount)"))
	considerationItemTypehash = crypto.Keccak256Hash([]byte(
		"ConsiderationItem(uint8 itemType,address token,uint256 identifierOrCriteria,uint256 startAmount,uint256 endAmount,address recipient)"))
	orderTypehash = crypto.Keccak256Hash([]byte(
		"OrderComponents(address offerer,address zone,OfferItem[] offer,ConsiderationItem[] consideration,uint8 orderType,uint256 startTime,uint256 endTime,bytes32 zoneHash,uint256 salt,bytes32 conduitKey,uint256 counter)" +
			"ConsiderationItem(uint8 itemType,address token,uint256 identifierOrCriteria,uint256 startAmount,uint256 endAmount,address recipient)" +
			"OfferItem(uint8 itemType,address token,uint256 identifierOrCriteria,uint256 startAmount,uint256 endAmount)"))
)

// ABI primitive types reused when ABI-encoding the EIP-712 preimages. Parsed
// once - the inputs are constant so the errors can never trigger.
var (
	abiAddressTy = mustNewABIType("address")
	abiUint256Ty = mustNewABIType("uint256")
	abiUint8Ty   = mustNewABIType("uint8")
	abiBytes32Ty = mustNewABIType("bytes32")
)

func mustNewABIType(t string) abi.Type {
	ty, err := abi.NewType(t, "", nil)
	if err != nil {
		panic(fmt.Sprintf("seaporttrades: invalid abi type %q: %v", t, err))
	}
	return ty
}

// hashOfferItem reproduces Seaport's per-offer-item EIP-712 hashing:
// keccak256(abi.encode(OFFER_ITEM_TYPEHASH, itemType, token, identifier,
// startAmount, endAmount)).
func hashOfferItem(item contract.OfferItem) ([32]byte, error) {
	args := abi.Arguments{
		{Type: abiBytes32Ty}, // OFFER_ITEM_TYPEHASH
		{Type: abiUint8Ty},   // itemType
		{Type: abiAddressTy}, // token
		{Type: abiUint256Ty}, // identifierOrCriteria
		{Type: abiUint256Ty}, // startAmount
		{Type: abiUint256Ty}, // endAmount
	}
	packed, err := args.Pack(
		[32]byte(offerItemTypehash),
		item.ItemType,
		item.Token,
		item.IdentifierOrCriteria,
		item.StartAmount,
		item.EndAmount,
	)
	if err != nil {
		return [32]byte{}, fmt.Errorf("could not pack offer item: %w", err)
	}
	return crypto.Keccak256Hash(packed), nil
}

// hashConsiderationItem reproduces Seaport's per-consideration-item EIP-712
// hashing (the offer-item layout plus the recipient field).
func hashConsiderationItem(item contract.ConsiderationItem) ([32]byte, error) {
	args := abi.Arguments{
		{Type: abiBytes32Ty}, // CONSIDERATION_ITEM_TYPEHASH
		{Type: abiUint8Ty},   // itemType
		{Type: abiAddressTy}, // token
		{Type: abiUint256Ty}, // identifierOrCriteria
		{Type: abiUint256Ty}, // startAmount
		{Type: abiUint256Ty}, // endAmount
		{Type: abiAddressTy}, // recipient
	}
	packed, err := args.Pack(
		[32]byte(considerationItemTypehash),
		item.ItemType,
		item.Token,
		item.IdentifierOrCriteria,
		item.StartAmount,
		item.EndAmount,
		item.Recipient,
	)
	if err != nil {
		return [32]byte{}, fmt.Errorf("could not pack consideration item: %w", err)
	}
	return crypto.Keccak256Hash(packed), nil
}

// computeOrderHash reproduces Seaport.getOrderHash for the given order
// parameters and offerer counter: the offer and consideration arrays are hashed
// into rolling keccak digests, then combined with the order fields under the
// ORDER_TYPEHASH. The result matches the on-chain order hash (verified at
// startup) so signatures recover to the offerer.
func computeOrderHash(params contract.OrderParameters, counter *big.Int) ([32]byte, error) {
	offerHashes := make([]byte, 0, len(params.Offer)*32)
	for _, item := range params.Offer {
		h, err := hashOfferItem(item)
		if err != nil {
			return [32]byte{}, err
		}
		offerHashes = append(offerHashes, h[:]...)
	}
	offerHash := crypto.Keccak256Hash(offerHashes)

	considerationHashes := make([]byte, 0, len(params.Consideration)*32)
	for _, item := range params.Consideration {
		h, err := hashConsiderationItem(item)
		if err != nil {
			return [32]byte{}, err
		}
		considerationHashes = append(considerationHashes, h[:]...)
	}
	considerationHash := crypto.Keccak256Hash(considerationHashes)

	args := abi.Arguments{
		{Type: abiBytes32Ty}, // ORDER_TYPEHASH
		{Type: abiAddressTy}, // offerer
		{Type: abiAddressTy}, // zone
		{Type: abiBytes32Ty}, // offerHash
		{Type: abiBytes32Ty}, // considerationHash
		{Type: abiUint8Ty},   // orderType
		{Type: abiUint256Ty}, // startTime
		{Type: abiUint256Ty}, // endTime
		{Type: abiBytes32Ty}, // zoneHash
		{Type: abiUint256Ty}, // salt
		{Type: abiBytes32Ty}, // conduitKey
		{Type: abiUint256Ty}, // counter
	}
	packed, err := args.Pack(
		[32]byte(orderTypehash),
		params.Offerer,
		params.Zone,
		offerHash,
		considerationHash,
		params.OrderType,
		params.StartTime,
		params.EndTime,
		params.ZoneHash,
		params.Salt,
		params.ConduitKey,
		counter,
	)
	if err != nil {
		return [32]byte{}, fmt.Errorf("could not pack order: %w", err)
	}
	return crypto.Keccak256Hash(packed), nil
}

// toOrderComponents converts the OrderParameters used on the fulfillment path
// into the OrderComponents shape Seaport.getOrderHash expects: the
// totalOriginalConsiderationItems count is dropped and replaced by the offerer's
// counter. Used only by the startup hashing self-check.
func toOrderComponents(params contract.OrderParameters, counter *big.Int) contract.OrderComponents {
	return contract.OrderComponents{
		Offerer:       params.Offerer,
		Zone:          params.Zone,
		Offer:         params.Offer,
		Consideration: params.Consideration,
		OrderType:     params.OrderType,
		StartTime:     params.StartTime,
		EndTime:       params.EndTime,
		ZoneHash:      params.ZoneHash,
		Salt:          params.Salt,
		ConduitKey:    params.ConduitKey,
		Counter:       counter,
	}
}

// deriveEIP712Digest builds the final signing digest from the cached domain
// separator and an order hash: keccak256(0x19 || 0x01 || domainSeparator ||
// orderHash).
func deriveEIP712Digest(domainSeparator [32]byte, orderHash [32]byte) [32]byte {
	preimage := make([]byte, 0, 2+32+32)
	preimage = append(preimage, 0x19, 0x01)
	preimage = append(preimage, domainSeparator[:]...)
	preimage = append(preimage, orderHash[:]...)
	return crypto.Keccak256Hash(preimage)
}

// signOrder computes the order hash for the market-offered order, derives the
// EIP-712 digest from the cached Seaport domain separator, and produces a
// 65-byte ECDSA signature (r || s || v) with v in the 27/28 convention Seaport's
// ecrecover path expects. The order is always offered by the market wallet, so
// its counter is used.
func (m *Market) signOrder(params contract.OrderParameters) ([]byte, [32]byte, error) {
	orderHash, err := computeOrderHash(params, m.counter)
	if err != nil {
		return nil, [32]byte{}, err
	}
	digest := deriveEIP712Digest(m.domainSeparator, orderHash)

	sig, err := crypto.Sign(digest[:], m.marketKey)
	if err != nil {
		return nil, [32]byte{}, fmt.Errorf("could not sign order: %w", err)
	}
	sig[64] += 27
	return sig, orderHash, nil
}

// zeroBytes32 is the conduit key / zone hash used throughout: a zero conduit key
// routes token transfers through Seaport directly (approvals are granted to the
// Seaport contract), so the scenario never needs to create a conduit.
var zeroBytes32 = [32]byte{}

// Seaport ItemType enum values (see SeaportInterface): only the two used by this
// scenario are named.
const (
	itemTypeERC20  uint8 = 1
	itemTypeERC721 uint8 = 2
)

// orderTypeFullOpen is Seaport's FULL_OPEN order type: anyone may fulfill, no
// partial fills, no zone validation.
const orderTypeFullOpen uint8 = 0

// maxOrderEndTime is used as every order's endTime. It is far past any realistic
// run, so orders never expire mid-scenario; startTime is always 0.
var maxOrderEndTime = new(big.Int).SetUint64(1 << 62)
