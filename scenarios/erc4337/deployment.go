package erc4337

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
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

	baseFeeWei, tipFeeWei := spamoor.ResolveFees(s.options.BaseFee, s.options.TipFee, s.options.BaseFeeWei, s.options.TipFeeWei)
	feeCap, tipCap, err := s.walletPool.GetSuggestedFees(client, baseFeeWei, tipFeeWei)
	if err != nil {
		return nil, fmt.Errorf("could not get tx fee: %w", err)
	}

	deploymentTxs := make([]*types.Transaction, 0, 4)
	info := &DeploymentInfo{}
	deployerNonce := deployerWallet.GetNonce()
	contractNonce := uint64(0)
	usedNonce := uint64(0)

	// deploy EntryPoint
	if redeploy || deployerNonce <= contractNonce {
		tx, err := deployerWallet.BuildBoundTxWithEstimate(ctx, client, s.walletPool.GetTxPool(), &txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Value:     uint256.NewInt(0),
		}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
			_, deployTx, _, err := contract.DeployEntryPoint(transactOpts, client.GetEthClient())
			return deployTx, err
		})
		if err != nil {
			return nil, fmt.Errorf("could not deploy EntryPoint: %w", err)
		}
		deploymentTxs = append(deploymentTxs, tx)
		usedNonce = tx.Nonce()
	} else {
		usedNonce = contractNonce
	}
	contractNonce++

	info.EntryPointAddr = crypto.CreateAddress(deployerWallet.GetAddress(), usedNonce)
	info.EntryPoint, err = contract.NewEntryPoint(info.EntryPointAddr, client.GetEthClient())
	if err != nil {
		return nil, fmt.Errorf("could not create instance of EntryPoint: %w", err)
	}

	// deploy SimpleAccountFactory
	if redeploy || deployerNonce <= contractNonce {
		tx, err := deployerWallet.BuildBoundTxWithEstimate(ctx, client, s.walletPool.GetTxPool(), &txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Value:     uint256.NewInt(0),
		}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
			_, deployTx, _, err := contract.DeploySimpleAccountFactory(transactOpts, client.GetEthClient(), info.EntryPointAddr)
			return deployTx, err
		})
		if err != nil {
			return nil, fmt.Errorf("could not deploy SimpleAccountFactory: %w", err)
		}
		deploymentTxs = append(deploymentTxs, tx)
		usedNonce = tx.Nonce()
	} else {
		usedNonce = contractNonce
	}
	contractNonce++

	info.FactoryAddr = crypto.CreateAddress(deployerWallet.GetAddress(), usedNonce)
	info.Factory, err = contract.NewSimpleAccountFactory(info.FactoryAddr, client.GetEthClient())
	if err != nil {
		return nil, fmt.Errorf("could not create instance of SimpleAccountFactory: %w", err)
	}

	// deploy Paymaster
	if redeploy || deployerNonce <= contractNonce {
		tx, err := deployerWallet.BuildBoundTxWithEstimate(ctx, client, s.walletPool.GetTxPool(), &txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Value:     uint256.NewInt(0),
		}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
			_, deployTx, _, err := contract.DeployAcceptAllPaymaster(transactOpts, client.GetEthClient(), info.EntryPointAddr)
			return deployTx, err
		})
		if err != nil {
			return nil, fmt.Errorf("could not deploy Paymaster: %w", err)
		}
		deploymentTxs = append(deploymentTxs, tx)
		usedNonce = tx.Nonce()
	} else {
		usedNonce = contractNonce
	}
	contractNonce++

	info.PaymasterAddr = crypto.CreateAddress(deployerWallet.GetAddress(), usedNonce)
	info.Paymaster, err = contract.NewAcceptAllPaymaster(info.PaymasterAddr, client.GetEthClient())
	if err != nil {
		return nil, fmt.Errorf("could not create instance of Paymaster: %w", err)
	}

	// deploy Counter
	if redeploy || deployerNonce <= contractNonce {
		tx, err := deployerWallet.BuildBoundTxWithEstimate(ctx, client, s.walletPool.GetTxPool(), &txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Value:     uint256.NewInt(0),
		}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
			_, deployTx, _, err := contract.DeployCounter(transactOpts, client.GetEthClient())
			return deployTx, err
		})
		if err != nil {
			return nil, fmt.Errorf("could not deploy Counter: %w", err)
		}
		deploymentTxs = append(deploymentTxs, tx)
		usedNonce = tx.Nonce()
	} else {
		usedNonce = contractNonce
	}
	contractNonce++

	info.CounterAddr = crypto.CreateAddress(deployerWallet.GetAddress(), usedNonce)
	info.Counter, err = contract.NewCounter(info.CounterAddr, client.GetEthClient())
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

	return info, nil
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
