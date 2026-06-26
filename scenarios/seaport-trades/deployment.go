package seaporttrades

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/scenarios/seaport-trades/contract"
	"github.com/ethpandaops/spamoor/spamoor"
)

// DeploymentInfo holds the addresses and bound instances of the Seaport stack
// deployed (or resolved) for a scenario run: the marketplace itself, its conduit
// controller, and the mock NFT collection and stablecoin that are traded.
type DeploymentInfo struct {
	ConduitControllerAddr common.Address
	SeaportAddr           common.Address
	Seaport               *contract.Seaport
	NFTAddr               common.Address
	NFT                   *contract.MintableNFT
	CoinAddr              common.Address
	Coin                  *contract.MintableToken

	// domainSeparator is read from Seaport.information() once and reused for every
	// order's EIP-712 digest (it is independent of the offerer).
	domainSeparator [32]byte
}

// DeployContracts deploys the Seaport stack from the "deployer" wallet via the
// CREATE2 deployment factory (so a restarted scenario re-uses the already
// deployed addresses): ConduitController first (Seaport's constructor reads the
// conduit code hashes from it), then Seaport bound to that controller, then the
// mock NFT collection and stablecoin. All are deployed with a zero seed so their
// addresses are stable across runs.
func (s *Scenario) DeployContracts(ctx context.Context) (*DeploymentInfo, error) {
	client := s.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0),
		spamoor.WithClientGroup(s.deployClientGroup()),
	)
	if client == nil {
		return nil, scenario.ErrNoClients
	}

	deployerWallet := s.walletPool.GetWellKnownWallet("deployer")
	if deployerWallet == nil {
		return nil, scenario.ErrNoWallet
	}

	baseFeeWei, tipFeeWei := spamoor.ResolveFees(s.options.BaseFee, s.options.TipFee, s.options.BaseFeeWei, s.options.TipFeeWei)
	feeCap, tipCap, err := s.walletPool.GetSuggestedFees(client, baseFeeWei, tipFeeWei)
	if err != nil {
		return nil, fmt.Errorf("could not get tx fee: %w", err)
	}

	deploymentTxs := []*types.Transaction{}
	info := &DeploymentInfo{}

	// deploy resolves (or builds the creation tx for) one contract via the CREATE2
	// factory with a zero seed, ABI-packing any constructor args onto the creation
	// bytecode (like the erc4337 scenario does for its constructor-taking
	// contracts).
	deploy := func(metadata *bind.MetaData, params ...interface{}) (common.Address, error) {
		parsed, err := metadata.GetAbi()
		if err != nil {
			return common.Address{}, fmt.Errorf("could not parse abi: %w", err)
		}
		packed, err := parsed.Pack("", params...)
		if err != nil {
			return common.Address{}, fmt.Errorf("could not pack constructor args: %w", err)
		}
		initCode := append(common.FromHex(metadata.Bin), packed...)

		seed := [32]byte{}
		addr, tx, err := s.walletPool.GetDeploymentFactory().GetContractDeployment(ctx, initCode, seed, client, deployerWallet, feeCap, tipCap, false)
		if err != nil {
			return common.Address{}, err
		}
		if tx != nil {
			deploymentTxs = append(deploymentTxs, tx)
		}
		return addr, nil
	}

	info.ConduitControllerAddr, err = deploy(contract.ConduitControllerMetaData)
	if err != nil {
		return nil, fmt.Errorf("could not deploy ConduitController: %w", err)
	}

	// Deploy (and await) the ConduitController before building the Seaport
	// deploy: Seaport's constructor reads the conduit code hashes from the
	// controller, so its deploy gas can only be estimated once the controller
	// exists on-chain. Without this the estimate reverts and the deployment
	// factory falls back to a large upper-bound gas (~44M for Seaport's 25KB
	// creation code under Amsterdam) instead of the real ~36.5M, which can
	// exceed an immature devnet's still-ramping block gas limit.
	if len(deploymentTxs) > 0 {
		if _, err := s.walletPool.GetTxPool().SendTransactionBatch(ctx, deployerWallet, deploymentTxs, &spamoor.BatchOptions{
			SendTransactionOptions: spamoor.SendTransactionOptions{
				Client:      client,
				ClientGroup: s.options.ClientGroup,
			},
			MaxRetries:   3,
			PendingLimit: 10,
		}); err != nil {
			return nil, fmt.Errorf("could not deploy ConduitController: %w", err)
		}
		s.logger.Infof("deployed ConduitController")
		deploymentTxs = deploymentTxs[:0]
	}

	info.SeaportAddr, err = deploy(contract.SeaportMetaData, info.ConduitControllerAddr)
	if err != nil {
		return nil, fmt.Errorf("could not deploy Seaport: %w", err)
	}
	info.Seaport, err = contract.NewSeaport(info.SeaportAddr, client.GetEthClient())
	if err != nil {
		return nil, fmt.Errorf("could not create Seaport instance: %w", err)
	}

	info.NFTAddr, err = deploy(contract.MintableNFTMetaData, "Spamoor Seaport NFT", "SPNFT")
	if err != nil {
		return nil, fmt.Errorf("could not deploy MintableNFT: %w", err)
	}
	info.NFT, err = contract.NewMintableNFT(info.NFTAddr, client.GetEthClient())
	if err != nil {
		return nil, fmt.Errorf("could not create MintableNFT instance: %w", err)
	}

	info.CoinAddr, err = deploy(contract.MintableTokenMetaData, "Spamoor Seaport USD", "SPUSD")
	if err != nil {
		return nil, fmt.Errorf("could not deploy MintableToken: %w", err)
	}
	info.Coin, err = contract.NewMintableToken(info.CoinAddr, client.GetEthClient())
	if err != nil {
		return nil, fmt.Errorf("could not create MintableToken instance: %w", err)
	}

	if len(deploymentTxs) > 0 {
		_, err := s.walletPool.GetTxPool().SendTransactionBatch(ctx, deployerWallet, deploymentTxs, &spamoor.BatchOptions{
			SendTransactionOptions: spamoor.SendTransactionOptions{
				Client:      client,
				ClientGroup: s.options.ClientGroup,
			},
			MaxRetries:   3,
			PendingLimit: 10,
			LogFn: func(confirmedCount int, totalCount int) {
				s.logger.Infof("deploying contracts... (%v/%v)", confirmedCount, totalCount)
			},
			LogInterval: 10,
		})
		if err != nil {
			return nil, fmt.Errorf("could not send deployment txs: %w", err)
		}
		s.logger.Infof("deployed seaport stack (%d txs)", len(deploymentTxs))
	} else {
		s.logger.Infof("seaport stack already deployed, reusing")
	}

	// Cache the EIP-712 domain separator (offerer-independent) for order signing.
	infoResult, err := info.Seaport.Information(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, fmt.Errorf("could not read seaport information: %w", err)
	}
	info.domainSeparator = infoResult.DomainSeparator
	s.logger.Infof("seaport %s (version %s), nft %s, coin %s", info.SeaportAddr.Hex(), infoResult.Version, info.NFTAddr.Hex(), info.CoinAddr.Hex())

	return info, nil
}
