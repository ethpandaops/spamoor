package erc4337

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/scenarios/erc4337/contract"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
)

// DeploymentInfo holds the addresses and bound instances of the ERC-4337 stack
// deployed (or resolved) for a scenario run.
type DeploymentInfo struct {
	EntryPointAddr common.Address
	EntryPoint     *contract.EntryPoint
	FactoryAddr    common.Address
	Factory        *contract.SimpleAccountFactory
	PaymasterAddr  common.Address
	Paymaster      *contract.AcceptAllPaymaster
	CounterAddr    common.Address
	Counter        *contract.Counter
}

// DeployContracts deploys the ERC-4337 v0.7 stack from the "deployer" wallet:
// EntryPoint, SimpleAccountFactory (bound to the EntryPoint), the accept-all
// Paymaster, and the Counter target. It mirrors the resume-safe nonce-guard
// pattern used by the other multi-contract scenarios so a restarted scenario
// re-uses the already-deployed addresses instead of redeploying.
func (s *Scenario) DeployContracts(ctx context.Context, redeploy bool) (*DeploymentInfo, error) {
	client := s.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0),
		spamoor.WithClientGroup(s.options.ClientGroup),
	)
	if client == nil {
		return nil, scenario.ErrNoClients
	}

	deployerWallet := s.walletPool.GetWellKnownWallet("deployer")
	if deployerWallet == nil {
		return nil, scenario.ErrNoWallet
	}

	deployerSeed := [32]byte{}
	copy(deployerSeed[:], deployerWallet.GetAddress().Bytes())

	if redeploy {
		copy(deployerSeed[20:], []byte(fmt.Sprintf("%x", deployerWallet.GetNonce()+1)))
	}

	baseFeeWei, tipFeeWei := spamoor.ResolveFees(s.options.BaseFee, s.options.TipFee, s.options.BaseFeeWei, s.options.TipFeeWei)
	feeCap, tipCap, err := s.walletPool.GetSuggestedFees(client, baseFeeWei, tipFeeWei)
	if err != nil {
		return nil, fmt.Errorf("could not get tx fee: %w", err)
	}

	deploymentTxs := []*types.Transaction{}
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
			binary.BigEndian.PutUint32(deployerSeed[28:], salt)
		}
		addr, tx, err := s.walletPool.GetDeploymentFactory().GetContractDeployment(ctx, initCodeBytes, deployerSeed, client, deployerWallet, feeCap, tipCap, false)
		if err != nil {
			return common.Address{}, err
		}

		if tx != nil {
			deploymentTxs = append(deploymentTxs, tx)
		}

		return addr, nil
	}

	// deploy EntryPoint
	deploymentInfo.EntryPointAddr, err = deployContract(contract.EntryPointMetaData, true, 0)
	if err != nil {
		return nil, fmt.Errorf("could not deploy EntryPoint: %w", err)
	}
	deploymentInfo.EntryPoint, err = contract.NewEntryPoint(deploymentInfo.EntryPointAddr, client.GetEthClient())
	if err != nil {
		return nil, fmt.Errorf("could not create instance of EntryPoint: %w", err)
	}

	// deploy SimpleAccountFactory
	deploymentInfo.FactoryAddr, err = deployContract(contract.SimpleAccountFactoryMetaData, true, 0, deploymentInfo.EntryPointAddr)
	if err != nil {
		return nil, fmt.Errorf("could not deploy SimpleAccountFactory: %w", err)
	}
	deploymentInfo.Factory, err = contract.NewSimpleAccountFactory(deploymentInfo.FactoryAddr, client.GetEthClient())
	if err != nil {
		return nil, fmt.Errorf("could not create instance of SimpleAccountFactory: %w", err)
	}

	// deploy Paymaster
	deploymentInfo.PaymasterAddr, err = deployContract(contract.AcceptAllPaymasterMetaData, true, 0, deploymentInfo.EntryPointAddr)
	if err != nil {
		return nil, fmt.Errorf("could not deploy Paymaster: %w", err)
	}
	deploymentInfo.Paymaster, err = contract.NewAcceptAllPaymaster(deploymentInfo.PaymasterAddr, client.GetEthClient())
	if err != nil {
		return nil, fmt.Errorf("could not create instance of Paymaster: %w", err)
	}

	// deploy Counter
	deploymentInfo.CounterAddr, err = deployContract(contract.CounterMetaData, true, 0)
	if err != nil {
		return nil, fmt.Errorf("could not deploy Counter: %w", err)
	}
	deploymentInfo.Counter, err = contract.NewCounter(deploymentInfo.CounterAddr, client.GetEthClient())
	if err != nil {
		return nil, fmt.Errorf("could not create instance of Counter: %w", err)
	}

	// submit & await all deployment transactions
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
		s.logger.Infof("contract deployment complete. (%v/%v)", len(deploymentTxs), len(deploymentTxs))
	}

	return deploymentInfo, nil
}

// ensurePaymasterDeposit tops up the paymaster's EntryPoint deposit to at least
// topUpAmount when it falls below minBalance. The deposit funds gas for every
// sponsored UserOperation, so a long run drains it over time; the scenario calls
// this once before sending and periodically afterwards. Funding comes from the
// locked root wallet (deposits can be large) and flows back to the child wallets
// that act as handleOps beneficiaries.
func (s *Scenario) ensurePaymasterDeposit(ctx context.Context, minBalance, topUpAmount *big.Int) error {
	if s.deploymentInfo == nil {
		return fmt.Errorf("contracts not deployed yet")
	}

	client := s.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0),
		spamoor.WithClientGroup(s.options.ClientGroup),
	)
	if client == nil {
		return scenario.ErrNoClients
	}

	deposit, err := s.deploymentInfo.Paymaster.GetDeposit(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("could not read paymaster deposit: %w", err)
	}
	if deposit.Cmp(minBalance) >= 0 {
		return nil
	}

	baseFeeWei, tipFeeWei := spamoor.ResolveFees(s.options.BaseFee, s.options.TipFee, s.options.BaseFeeWei, s.options.TipFeeWei)
	feeCap, tipCap, err := s.walletPool.GetSuggestedFees(client, baseFeeWei, tipFeeWei)
	if err != nil {
		return fmt.Errorf("could not get tx fee: %w", err)
	}

	rootWallet := s.walletPool.GetRootWallet()
	return rootWallet.WithWalletLock(ctx, 1, uint256.MustFromBig(topUpAmount), s.walletPool.GetClientPool(), func(reason string) {
		s.logger.Infof("root wallet is locked, %s", reason)
	}, func() error {
		tx, err := rootWallet.GetWallet().BuildBoundTxWithEstimate(ctx, client, s.walletPool.GetTxPool(), &txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Value:     uint256.MustFromBig(topUpAmount),
		}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
			return s.deploymentInfo.Paymaster.Deposit(transactOpts)
		})
		if err != nil {
			return fmt.Errorf("could not build paymaster deposit tx: %w", err)
		}

		_, err = s.walletPool.GetTxPool().SendAndAwaitTransaction(ctx, rootWallet.GetWallet(), tx, &spamoor.SendTransactionOptions{
			Client:      client,
			ClientGroup: s.options.ClientGroup,
			Rebroadcast: true,
		})
		if err != nil {
			return fmt.Errorf("could not send paymaster deposit tx: %w", err)
		}

		s.logger.Infof("funded paymaster deposit with %v wei (was %v wei)", topUpAmount.String(), deposit.String())
		return nil
	})
}
