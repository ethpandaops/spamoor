package seaporttrades

import (
	"context"
	cryptorand "crypto/rand"
	"fmt"
	"math/big"
	mathrand "math/rand"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"

	"github.com/ethpandaops/spamoor/scenarios/seaport-trades/contract"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
)

// fulfillGasLimit is the static gas limit for every fulfillOrder spam tx. A
// fulfillment moves one ERC721 and one ERC20, marks the order filled (a fresh
// Seaport storage slot under the Amsterdam state-creation surcharge), and
// recovers the signature; this was measured at ~197k gas on an Amsterdam devnet.
// The budget keeps comfortable headroom over that so the hot path needs no
// per-tx estimate. Stays well under the EIP-7825 per-tx gas cap (2^24).
const fulfillGasLimit = 500000

// errNoTrade is returned when neither a buy nor a sell can currently be built
// for the chosen wallet (e.g. it holds no NFTs to sell and cannot afford to
// buy). The scenario treats it like any other "could not send" and moves on.
var errNoTrade = fmt.Errorf("no fulfillable trade available")

// tradeResult finalizes the in-memory inventory/balance reservations made while
// building a trade: committing them to their destinations on success, or rolling
// them back on failure. It is invoked from the fulfillment's completion callback.
type tradeResult func(success bool)

// BuildTrade builds one Seaport fulfillOrder transaction for the given taker
// wallet. It chooses buy vs sell from the configured ratio, nudged by the
// wallet's current inventory (a wallet with no NFTs must buy; one holding many
// is pushed to sell), and falls back to the other direction if the preferred one
// is not currently fulfillable. The returned tradeResult must be called with the
// fulfillment's success once it completes so the inventory caches stay accurate.
func (m *Market) BuildTrade(ctx context.Context, taker *spamoor.Wallet, feeCap, tipCap *big.Int) (*types.Transaction, tradeResult, error) {
	takerAddr := taker.GetAddress()
	price := m.randomPrice()

	isBuy := mathrand.Intn(100) < int(m.options.BuyRatio)
	switch {
	case m.walletNFTCount(takerAddr) == 0:
		isBuy = true // nothing to sell
	case m.walletNFTCount(takerAddr) >= int(m.options.SellThreshold):
		isBuy = false // shed a large inventory
	}

	order := []func() (*types.Transaction, tradeResult, bool, error){
		func() (*types.Transaction, tradeResult, bool, error) {
			return m.buildBuy(ctx, taker, feeCap, tipCap, price)
		},
		func() (*types.Transaction, tradeResult, bool, error) {
			return m.buildSell(ctx, taker, feeCap, tipCap, price)
		},
	}
	if !isBuy {
		order[0], order[1] = order[1], order[0]
	}

	for _, build := range order {
		tx, result, ok, err := build()
		if err != nil {
			return nil, nil, err
		}
		if ok {
			return tx, result, nil
		}
	}
	return nil, nil, errNoTrade
}

// buildBuy builds a fulfillment where the market lists one of its NFTs for the
// stablecoin and the taker buys it: the taker pays `price` coin to the market and
// receives the NFT. ok=false when the market has no inventory or the taker cannot
// afford the price.
func (m *Market) buildBuy(ctx context.Context, taker *spamoor.Wallet, feeCap, tipCap, price *big.Int) (*types.Transaction, tradeResult, bool, error) {
	takerAddr := taker.GetAddress()

	tokenID, ok := m.reserveMarketNFT()
	if !ok {
		return nil, nil, false, nil
	}
	if !m.reserveCoin(takerAddr, price) {
		m.giveMarketNFT(tokenID)
		return nil, nil, false, nil
	}

	params := contract.OrderParameters{
		Offerer: m.marketAddr,
		Zone:    common.Address{},
		Offer: []contract.OfferItem{{
			ItemType:             itemTypeERC721,
			Token:                m.deployment.NFTAddr,
			IdentifierOrCriteria: tokenID,
			StartAmount:          big.NewInt(1),
			EndAmount:            big.NewInt(1),
		}},
		Consideration: []contract.ConsiderationItem{{
			ItemType:             itemTypeERC20,
			Token:                m.deployment.CoinAddr,
			IdentifierOrCriteria: big.NewInt(0),
			StartAmount:          price,
			EndAmount:            price,
			Recipient:            m.marketAddr,
		}},
		OrderType:                       orderTypeFullOpen,
		StartTime:                       big.NewInt(0),
		EndTime:                         maxOrderEndTime,
		ZoneHash:                        zeroBytes32,
		Salt:                            randomSalt(),
		ConduitKey:                      zeroBytes32,
		TotalOriginalConsiderationItems: big.NewInt(1),
	}

	tx, err := m.buildFulfill(ctx, taker, feeCap, tipCap, params)
	if err != nil {
		m.giveMarketNFT(tokenID)
		m.creditCoin(takerAddr, price)
		return nil, nil, false, err
	}

	result := func(success bool) {
		if success {
			m.giveWalletNFT(takerAddr, tokenID)
			m.creditCoin(m.marketAddr, price)
		} else {
			m.giveMarketNFT(tokenID)
			m.creditCoin(takerAddr, price)
		}
	}
	return tx, result, true, nil
}

// buildSell builds a fulfillment where the market bids the stablecoin for one of
// the taker's NFTs and the taker accepts: the taker hands over the NFT and
// receives `price` coin from the market. ok=false when the taker holds no NFT or
// the market's coin float is (cache-)exhausted.
func (m *Market) buildSell(ctx context.Context, taker *spamoor.Wallet, feeCap, tipCap, price *big.Int) (*types.Transaction, tradeResult, bool, error) {
	takerAddr := taker.GetAddress()

	tokenID, ok := m.reserveWalletNFT(takerAddr)
	if !ok {
		return nil, nil, false, nil
	}
	if !m.reserveCoin(m.marketAddr, price) {
		m.giveWalletNFT(takerAddr, tokenID)
		return nil, nil, false, nil
	}

	params := contract.OrderParameters{
		Offerer: m.marketAddr,
		Zone:    common.Address{},
		Offer: []contract.OfferItem{{
			ItemType:             itemTypeERC20,
			Token:                m.deployment.CoinAddr,
			IdentifierOrCriteria: big.NewInt(0),
			StartAmount:          price,
			EndAmount:            price,
		}},
		Consideration: []contract.ConsiderationItem{{
			ItemType:             itemTypeERC721,
			Token:                m.deployment.NFTAddr,
			IdentifierOrCriteria: tokenID,
			StartAmount:          big.NewInt(1),
			EndAmount:            big.NewInt(1),
			Recipient:            m.marketAddr,
		}},
		OrderType:                       orderTypeFullOpen,
		StartTime:                       big.NewInt(0),
		EndTime:                         maxOrderEndTime,
		ZoneHash:                        zeroBytes32,
		Salt:                            randomSalt(),
		ConduitKey:                      zeroBytes32,
		TotalOriginalConsiderationItems: big.NewInt(1),
	}

	tx, err := m.buildFulfill(ctx, taker, feeCap, tipCap, params)
	if err != nil {
		m.giveWalletNFT(takerAddr, tokenID)
		m.creditCoin(m.marketAddr, price)
		return nil, nil, false, err
	}

	result := func(success bool) {
		if success {
			m.giveMarketNFT(tokenID)
			m.creditCoin(takerAddr, price)
		} else {
			m.giveWalletNFT(takerAddr, tokenID)
			m.creditCoin(m.marketAddr, price)
		}
	}
	return tx, result, true, nil
}

// buildFulfill signs the market order and wraps it in a taker-submitted
// fulfillOrder transaction with a zero fulfiller conduit key (transfers route
// through Seaport directly against the approvals granted at seed time).
func (m *Market) buildFulfill(ctx context.Context, taker *spamoor.Wallet, feeCap, tipCap *big.Int, params contract.OrderParameters) (*types.Transaction, error) {
	signature, _, err := m.signOrder(params)
	if err != nil {
		return nil, err
	}
	order := contract.Order{Parameters: params, Signature: signature}

	return taker.BuildBoundTx(ctx, &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       fulfillGasLimit,
		Value:     uint256.NewInt(0),
	}, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return m.deployment.Seaport.FulfillOrder(opts, order, zeroBytes32)
	})
}

// randomPrice returns a uniform random stablecoin amount in [MinPrice, MaxPrice].
// crypto/rand keeps it concurrency-safe and handles spans wider than int64.
func (m *Market) randomPrice() *big.Int {
	minP, ok := new(big.Int).SetString(m.options.MinPrice, 10)
	if !ok {
		minP = big.NewInt(1)
	}
	maxP, ok := new(big.Int).SetString(m.options.MaxPrice, 10)
	if !ok || maxP.Cmp(minP) <= 0 {
		return minP
	}
	span := new(big.Int).Add(new(big.Int).Sub(maxP, minP), big.NewInt(1))
	delta, err := cryptorand.Int(cryptorand.Reader, span)
	if err != nil {
		return minP
	}
	return new(big.Int).Add(minP, delta)
}

// randomSalt returns a 256-bit random Seaport order salt, so every order (even
// for the same tokenId and price) hashes uniquely and is independently
// fulfillable.
func randomSalt() *big.Int {
	max := new(big.Int).Lsh(big.NewInt(1), 256)
	salt, err := cryptorand.Int(cryptorand.Reader, max)
	if err != nil {
		return big.NewInt(0)
	}
	return salt
}
