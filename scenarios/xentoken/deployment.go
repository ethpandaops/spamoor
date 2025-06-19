package xentoken

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethpandaops/spamoor/scenarios/xentoken/contract"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/holiman/uint256"
)

// we need to inject a dynamic address for XENCrypto, so we need to use a global mutex to prevent race conditions
var xenCryptoDeploymentMutex sync.Mutex

type DeploymentInfo struct {
	XENCryptoAddr     common.Address
	SybilAttackerAddr common.Address
	SybilAttacker     *contract.XENSybilAttacker
}

func (s *Scenario) DeployContracts(ctx context.Context, xenTokenAddress *common.Address, redeploy bool) (*DeploymentInfo, error) {
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, s.options.ClientGroup)
	if client == nil {
		return nil, fmt.Errorf("no client available")
	}

	feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(client, s.options.BaseFee, s.options.TipFee)
	if err != nil {
		return nil, fmt.Errorf("could not get tx fee: %w", err)
	}

	deploymentTxs := map[*spamoor.Wallet][]*types.Transaction{}

	if xenTokenAddress == nil {
		deployerWallet := s.walletPool.GetWellKnownWallet("xen-deployer")
		deployerNonce := deployerWallet.GetNonce()
		contractNonce := uint64(0)
		usedNonce := uint64(0)

		// deploy XENMath
		if redeploy || deployerNonce <= contractNonce {
			tx, err := deployerWallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
				GasFeeCap: uint256.MustFromBig(feeCap),
				GasTipCap: uint256.MustFromBig(tipCap),
				Gas:       300000,
				Value:     uint256.NewInt(0),
			}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
				_, deployTx, _, err := contract.DeployXENMath(transactOpts, client.GetEthClient())
				return deployTx, err
			})
			if err != nil {
				return nil, fmt.Errorf("could not deploy XENMath: %w", err)
			}
			deploymentTxs[deployerWallet] = append(deploymentTxs[deployerWallet], tx)
			usedNonce = tx.Nonce()
		} else {
			usedNonce = contractNonce
		}
		contractNonce++

		xenMathAddr := crypto.CreateAddress(deployerWallet.GetAddress(), usedNonce)
		xenMathAddrHex := fmt.Sprintf("%040s", strings.TrimPrefix(xenMathAddr.Hex(), "0x"))

		// deploy XENCrypto
		if redeploy || deployerNonce <= contractNonce {
			buildXenCryptoTx := func() (*types.Transaction, error) {
				re := regexp.MustCompile(`__\$[a-f0-9]{34}\$__`)

				xenCryptoDeploymentMutex.Lock()
				defer xenCryptoDeploymentMutex.Unlock()

				origStr := contract.XENCryptoBin
				contract.XENCryptoBin = re.ReplaceAllString(origStr, xenMathAddrHex)
				defer func() {
					contract.XENCryptoBin = origStr
				}()

				return deployerWallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
					GasFeeCap: uint256.MustFromBig(feeCap),
					GasTipCap: uint256.MustFromBig(tipCap),
					Gas:       3000000,
					Value:     uint256.NewInt(0),
				}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
					_, deployTx, _, err := contract.DeployXENCrypto(transactOpts, client.GetEthClient())
					return deployTx, err
				})
			}

			tx, err := buildXenCryptoTx()
			if err != nil {
				return nil, fmt.Errorf("could not deploy XENCrypto: %w", err)
			}
			deploymentTxs[deployerWallet] = append(deploymentTxs[deployerWallet], tx)
			usedNonce = tx.Nonce()
		} else {
			usedNonce = contractNonce
		}
		contractNonce++

		xenCryptoAddr := crypto.CreateAddress(deployerWallet.GetAddress(), usedNonce)
		xenTokenAddress = &xenCryptoAddr
	}

	deployerWallet := s.walletPool.GetWellKnownWallet("misc-deployer")
	deployerNonce := deployerWallet.GetNonce()
	contractNonce := uint64(0)
	usedNonce := uint64(0)

	// deploy SybilAttacker
	if redeploy || deployerNonce <= contractNonce {
		tx, err := deployerWallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Gas:       2000000,
			Value:     uint256.NewInt(0),
		}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
			_, deployTx, _, err := contract.DeployXENSybilAttacker(transactOpts, client.GetEthClient())
			return deployTx, err
		})
		if err != nil {
			return nil, fmt.Errorf("could not deploy XENMath: %w", err)
		}
		deploymentTxs[deployerWallet] = append(deploymentTxs[deployerWallet], tx)
		usedNonce = tx.Nonce()
	} else {
		usedNonce = contractNonce
	}
	contractNonce++

	sybilAttackerAddr := crypto.CreateAddress(deployerWallet.GetAddress(), usedNonce)
	sybilAttacker, err := contract.NewXENSybilAttacker(sybilAttackerAddr, client.GetEthClient())
	if err != nil {
		return nil, err
	}

	// submit & await all deployment transactions
	if len(deploymentTxs) > 0 {
		s.logger.Infof("deploying contracts... (0/%v)", len(deploymentTxs))
		wg := sync.WaitGroup{}
		hasErrors := atomic.Bool{}
		for deployerWallet, txs := range deploymentTxs {
			wg.Add(1)
			go func(deployerWallet *spamoor.Wallet, txs []*types.Transaction) {
				defer wg.Done()

				for txIdx := 0; txIdx < len(txs); txIdx += 10 {
					endIdx := txIdx + 10
					if txIdx > 0 {
						s.logger.Infof("deploying contracts for %s... (%v/%v)", s.walletPool.GetWalletName(deployerWallet.GetAddress()), txIdx, len(txs))
					}
					if endIdx > len(txs) {
						endIdx = len(txs)
					}
					err := s.walletPool.GetTxPool().SendAndAwaitTxRange(ctx, deployerWallet, txs[txIdx:endIdx], &spamoor.SendTransactionOptions{
						Client: client,
						OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
							if err != nil {
								s.logger.Warnf("could not send deployment tx %v: %v", tx.Hash().String(), err)
							}
						},
					})
					if err != nil {
						s.logger.Warnf("could not send deployment txs: %v", err)
						hasErrors.Store(true)
					}
				}
			}(deployerWallet, txs)
		}
		wg.Wait()
		if hasErrors.Load() {
			return nil, fmt.Errorf("some deployment transactions failed")
		}

		s.logger.Infof("contract deployment complete. (%v/%v)", len(deploymentTxs), len(deploymentTxs))
	}

	return &DeploymentInfo{
		XENCryptoAddr:     *xenTokenAddress,
		SybilAttackerAddr: sybilAttackerAddr,
		SybilAttacker:     sybilAttacker,
	}, nil
}
