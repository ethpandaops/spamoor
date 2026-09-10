package uniswapswaps

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/scenarios/uniswap-swaps/contract"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/txtypes"
)

// DeploymentInfo holds the deployed Uniswap v2 contract set for the scenario.
// Two factories (each with its own router) are deployed so that every DAI
// token gets a pair with the shared quote token on both factories.
type DeploymentInfo struct {
	// Weth9Addr is only needed because the canonical router requires a WETH
	// address at construction; the scenario never trades ETH/WETH.
	Weth9Addr             common.Address
	QuoteAddr             common.Address
	Quote                 *contract.Dai
	UniswapFactoryAAddr   common.Address
	UniswapFactoryA       *contract.UniswapV2Factory
	UniswapRouterAAddr    common.Address
	UniswapRouterA        *contract.UniswapV2Router02
	UniswapFactoryBAddr   common.Address
	UniswapFactoryB       *contract.UniswapV2Factory
	UniswapRouterBAddr    common.Address
	UniswapRouterB        *contract.UniswapV2Router02
	LiquidityProviderAddr common.Address
	LiquidityProvider     *contract.PairLiquidityProvider
	Pairs                 []PairDeploymentInfo
}

type PairDeploymentInfo struct {
	DaiAddr   common.Address
	Dai       *contract.Dai
	PairAddrA common.Address
	PairA     *contract.UniswapV2Pair
	PairAddrB common.Address
	PairB     *contract.UniswapV2Pair
}

// Token deployment salts. The DAI/quote tokens are deployed "globally" (seed
// without the deployer address) and share identical init code, so the salt is
// all that distinguishes them: the quote token takes salt 0 and the DAI tokens
// take 1..PairCount.
const (
	quoteTokenSalt = 0
	daiTokenSalt   = 1
)

func (u *Uniswap) DeployUniswapPairs(redeploy bool) (*DeploymentInfo, error) {
	client := u.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0),
		spamoor.WithClientGroup(u.options.ClientGroup),
	)
	if client == nil {
		return nil, scenario.ErrNoClients
	}

	deployerWallet := u.walletPool.GetWellKnownWallet("deployer")
	ownerWallet := u.walletPool.GetWellKnownWallet("owner")
	if deployerWallet == nil || ownerWallet == nil {
		return nil, scenario.ErrNoWallet
	}

	deployerSeed := [32]byte{}
	copy(deployerSeed[:], deployerWallet.GetAddress().Bytes())

	if redeploy {
		copy(deployerSeed[20:], []byte(fmt.Sprintf("%x", deployerWallet.GetNonce()+1)))
	}

	baseFeeWei, tipFeeWei := spamoor.ResolveFees(u.options.BaseFee, u.options.TipFee, u.options.BaseFeeWei, u.options.TipFeeWei)
	feeCap, tipCap, err := u.walletPool.GetSuggestedFees(client, baseFeeWei, tipFeeWei)
	if err != nil {
		return nil, fmt.Errorf("could not get tx fee: %w", err)
	}

	deploymentTxs := []*txtypes.Transaction{}
	deploymentInfo := &DeploymentInfo{}
	deployContract := func(metadata *bind.MetaData, global bool, salt uint32, params ...interface{}) (common.Address, error) {
		parsed, err := metadata.GetAbi()
		if err != nil {
			return common.Address{}, err
		}
		if parsed == nil {
			return common.Address{}, fmt.Errorf("GetABI returned nil")
		}

		initCodeBytes := common.FromHex(metadata.Bin)

		packed, err := parsed.Pack("", params...)
		if err != nil {
			return common.Address{}, err
		}

		initCodeBytes = append(initCodeBytes, packed...)

		seed := [32]byte{}
		if !global {
			copy(seed[:], deployerSeed[:])
		}
		if salt != 0 {
			binary.BigEndian.PutUint32(seed[28:], salt)
		}
		addr, tx, err := u.walletPool.GetDeploymentFactory().GetContractDeployment(u.ctx, initCodeBytes, seed, client, deployerWallet, feeCap, tipCap, false)
		if err != nil {
			return common.Address{}, err
		}

		if tx != nil {
			deploymentTxs = append(deploymentTxs, tx)
		}

		return addr, nil
	}

	// deploy WETH9 (router constructor dependency only)
	deploymentInfo.Weth9Addr, err = deployContract(contract.WETH9MetaData, true, 0)
	if err != nil {
		return nil, fmt.Errorf("could not deploy WETH9: %w", err)
	}

	// deploy uniswap factory A
	deploymentInfo.UniswapFactoryAAddr, err = deployContract(contract.UniswapV2FactoryMetaData, false, 0, ownerWallet.GetAddress())
	if err != nil {
		return nil, fmt.Errorf("could not deploy uniswap v2 factory A: %w", err)
	}
	deploymentInfo.UniswapFactoryA, err = contract.NewUniswapV2Factory(deploymentInfo.UniswapFactoryAAddr, client.GetEthClient())
	if err != nil {
		return nil, fmt.Errorf("could not create instance of uniswap v2 factory A: %w", err)
	}

	// deploy uniswap factory B
	deploymentInfo.UniswapFactoryBAddr, err = deployContract(contract.UniswapV2FactoryMetaData, false, 1, ownerWallet.GetAddress())
	if err != nil {
		return nil, fmt.Errorf("could not deploy uniswap v2 factory B: %w", err)
	}
	deploymentInfo.UniswapFactoryB, err = contract.NewUniswapV2Factory(deploymentInfo.UniswapFactoryBAddr, client.GetEthClient())
	if err != nil {
		return nil, fmt.Errorf("could not create instance of uniswap v2 factory B: %w", err)
	}

	// deploy uniswap router A
	deploymentInfo.UniswapRouterAAddr, err = deployContract(contract.UniswapV2Router02MetaData, false, 0, deploymentInfo.UniswapFactoryAAddr, deploymentInfo.Weth9Addr)
	if err != nil {
		return nil, fmt.Errorf("could not deploy uniswap v2 router A: %w", err)
	}
	deploymentInfo.UniswapRouterA, err = contract.NewUniswapV2Router02(deploymentInfo.UniswapRouterAAddr, client.GetEthClient())
	if err != nil {
		return nil, fmt.Errorf("could not create instance of uniswap v2 router A: %w", err)
	}

	// deploy uniswap router B
	deploymentInfo.UniswapRouterBAddr, err = deployContract(contract.UniswapV2Router02MetaData, false, 1, deploymentInfo.UniswapFactoryBAddr, deploymentInfo.Weth9Addr)
	if err != nil {
		return nil, fmt.Errorf("could not deploy uniswap v2 router B: %w", err)
	}
	deploymentInfo.UniswapRouterB, err = contract.NewUniswapV2Router02(deploymentInfo.UniswapRouterBAddr, client.GetEthClient())
	if err != nil {
		return nil, fmt.Errorf("could not create instance of uniswap v2 router B: %w", err)
	}

	// deploy pair liquidity provider (owner-only helper that mints both tokens
	// on demand, so it needs no funding)
	deploymentInfo.LiquidityProviderAddr, err = deployContract(
		contract.PairLiquidityProviderMetaData, false, 0,
		ownerWallet.GetAddress(),
		deploymentInfo.UniswapRouterAAddr,
		deploymentInfo.UniswapRouterBAddr,
	)
	if err != nil {
		return nil, fmt.Errorf("could not deploy pair liquidity provider: %w", err)
	}
	deploymentInfo.LiquidityProvider, err = contract.NewPairLiquidityProvider(deploymentInfo.LiquidityProviderAddr, client.GetEthClient())
	if err != nil {
		return nil, fmt.Errorf("could not create instance of pair liquidity provider: %w", err)
	}

	// deploy the shared quote token
	deploymentInfo.QuoteAddr, err = deployContract(contract.DaiMetaData, true, quoteTokenSalt, deployerWallet.GetChainId())
	if err != nil {
		return nil, fmt.Errorf("could not deploy quote token: %w", err)
	}
	deploymentInfo.Quote, err = contract.NewDai(deploymentInfo.QuoteAddr, client.GetEthClient())
	if err != nil {
		return nil, fmt.Errorf("could not create instance of quote token: %w", err)
	}

	// deploy DAI tokens and derive their pair addresses on both factories
	pairInitCode := common.FromHex(contract.UniswapV2PairBin)
	pairInitHash := crypto.Keccak256(pairInitCode)

	for i := uint64(0); i < u.options.PairCount; i++ {
		pairInfo := PairDeploymentInfo{}

		// deploy Dai
		pairInfo.DaiAddr, err = deployContract(contract.DaiMetaData, true, daiTokenSalt+uint32(i), deployerWallet.GetChainId())
		if err != nil {
			return nil, fmt.Errorf("could not deploy Dai: %w", err)
		}
		pairInfo.Dai, err = contract.NewDai(pairInfo.DaiAddr, client.GetEthClient())
		if err != nil {
			return nil, fmt.Errorf("could not create instance of Dai: %w", err)
		}

		// pair on factory A
		pairInfo.PairAddrA = v2PairAddress(deploymentInfo.UniswapFactoryAAddr, pairInfo.DaiAddr, deploymentInfo.QuoteAddr, pairInitHash)
		pairInfo.PairA, err = contract.NewUniswapV2Pair(pairInfo.PairAddrA, client.GetEthClient())
		if err != nil {
			return nil, fmt.Errorf("could not create instance of uniswap v2 pair A: %w", err)
		}

		// pair on factory B
		pairInfo.PairAddrB = v2PairAddress(deploymentInfo.UniswapFactoryBAddr, pairInfo.DaiAddr, deploymentInfo.QuoteAddr, pairInitHash)
		pairInfo.PairB, err = contract.NewUniswapV2Pair(pairInfo.PairAddrB, client.GetEthClient())
		if err != nil {
			return nil, fmt.Errorf("could not create instance of uniswap v2 pair B: %w", err)
		}

		deploymentInfo.Pairs = append(deploymentInfo.Pairs, pairInfo)
	}

	// submit & await all deployment transactions
	if err := u.sendBatch(deployerWallet, client, deploymentTxs, "deploying contracts v2"); err != nil {
		return nil, err
	}

	// Phase 2: seed liquidity into pairs that don't have any yet. Both tokens
	// are minted by the liquidity provider, so the owner wallet only pays gas.
	// Built only after the deployment batch has been mined so eth_estimateGas
	// dispatches into the real contract code instead of treating the target as
	// an EOA.
	callOpts := &bind.CallOpts{Context: u.ctx}
	liquidityTxs := []*txtypes.Transaction{}

	// the liquidity provider splits the amounts evenly between both factories
	quoteLiquidity := new(big.Int).Mul(u.options.QuoteLiquidityPerPool, big.NewInt(2))
	daiLiquidity := new(big.Int).Mul(quoteLiquidity, new(big.Int).SetUint64(u.options.TokensPerQuote))

	for _, pairInfo := range deploymentInfo.Pairs {
		seeded, err := v2PairHasLiquidity(callOpts, deploymentInfo.UniswapFactoryA, client, pairInfo.DaiAddr, deploymentInfo.QuoteAddr)
		if err != nil {
			return nil, err
		}
		if seeded {
			continue
		}

		daiAddr := pairInfo.DaiAddr
		tx, err := ownerWallet.BuildBoundTxWithEstimate(u.ctx, client, u.walletPool.GetTxPool(), &txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Value:     uint256.NewInt(0),
		}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
			return deploymentInfo.LiquidityProvider.ProvidePairLiquidity(transactOpts, deploymentInfo.QuoteAddr, daiAddr, quoteLiquidity, daiLiquidity)
		})
		if err != nil {
			return nil, fmt.Errorf("could not provide liquidity for dai %v: %w", daiAddr.String(), err)
		}
		liquidityTxs = append(liquidityTxs, tx)
	}

	if err := u.sendBatch(ownerWallet, client, liquidityTxs, "providing liquidity"); err != nil {
		return nil, err
	}

	return deploymentInfo, nil
}

// v2PairAddress computes the CREATE2 address of the tokenA/tokenB pair on the
// given factory, mirroring UniswapV2Library.pairFor.
func v2PairAddress(factory, tokenA, tokenB common.Address, pairInitHash []byte) common.Address {
	token0, token1 := tokenA, tokenB
	if token1.Big().Cmp(token0.Big()) < 0 {
		token0, token1 = token1, token0
	}
	var salt [32]byte
	copy(salt[:], crypto.Keccak256(token0.Bytes(), token1.Bytes()))
	return crypto.CreateAddress2(factory, salt, pairInitHash)
}

// v2PairHasLiquidity reports whether the tokenA/tokenB pair exists on the
// factory and already holds reserves, in which case seeding is skipped so
// re-runs against an existing deployment don't double-seed.
func v2PairHasLiquidity(callOpts *bind.CallOpts, factory *contract.UniswapV2Factory, client *spamoor.Client, tokenA, tokenB common.Address) (bool, error) {
	pairAddr, err := factory.GetPair(callOpts, tokenA, tokenB)
	if err != nil {
		return false, fmt.Errorf("could not check pair existence: %w", err)
	}
	if pairAddr == (common.Address{}) {
		return false, nil
	}

	pair, err := contract.NewUniswapV2Pair(pairAddr, client.GetEthClient())
	if err != nil {
		return false, fmt.Errorf("could not bind pair %v: %w", pairAddr.Hex(), err)
	}
	reserves, err := pair.GetReserves(callOpts)
	if err != nil {
		return false, fmt.Errorf("could not read pair reserves: %w", err)
	}
	return reserves.Reserve0.Sign() > 0 && reserves.Reserve1.Sign() > 0, nil
}
