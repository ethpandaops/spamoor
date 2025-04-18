package uniswapswaps

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethpandaops/spamoor/scenarios/uniswap-swaps/contract"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
)

type UniswapOptions struct {
	BaseFee             uint64
	TipFee              uint64
	DaiPairs            uint64
	EthLiquidityPerPair *uint256.Int
	DaiLiquidityFactor  uint64
	ClientGroup         string
}

type Uniswap struct {
	ctx            context.Context
	walletPool     *spamoor.WalletPool
	deploymentInfo *DeploymentInfo
	logger         *logrus.Entry
	options        UniswapOptions

	// local cache of token balances
	tokenBalances      map[common.Address]map[common.Address]*big.Int
	tokenBalancesMutex sync.RWMutex

	// contract instances
	RouterA *contract.UniswapV2Router02
	RouterB *contract.UniswapV2Router02
	Weth    *contract.WETH9
	Tokens  map[common.Address]*contract.Dai
}

func NewUniswap(ctx context.Context, walletPool *spamoor.WalletPool, logger *logrus.Entry, options UniswapOptions) *Uniswap {
	return &Uniswap{
		ctx:           ctx,
		walletPool:    walletPool,
		logger:        logger,
		options:       options,
		tokenBalances: make(map[common.Address]map[common.Address]*big.Int),
	}
}

// Initialize contract instances to reuse
func (u *Uniswap) InitializeContracts(deploymentInfo *DeploymentInfo) error {
	u.deploymentInfo = deploymentInfo

	client := u.walletPool.GetClient(spamoor.SelectClientByIndex, 0, u.options.ClientGroup)
	if client == nil {
		return fmt.Errorf("no client available")
	}

	// Initialize router A
	routerA, err := contract.NewUniswapV2Router02(u.deploymentInfo.UniswapRouterAAddr, client.GetEthClient())
	if err != nil {
		return fmt.Errorf("could not initialize router A: %w", err)
	}
	u.RouterA = routerA

	// Initialize router B
	routerB, err := contract.NewUniswapV2Router02(u.deploymentInfo.UniswapRouterBAddr, client.GetEthClient())
	if err != nil {
		return fmt.Errorf("could not initialize router B: %w", err)
	}
	u.RouterB = routerB

	// Initialize WETH9
	weth, err := contract.NewWETH9(u.deploymentInfo.Weth9Addr, client.GetEthClient())
	if err != nil {
		return fmt.Errorf("could not initialize WETH9: %w", err)
	}
	u.Weth = weth

	// Initialize token contracts
	u.Tokens = make(map[common.Address]*contract.Dai)
	for _, pair := range u.deploymentInfo.Pairs {
		token, err := contract.NewDai(pair.DaiAddr, client.GetEthClient())
		if err != nil {
			return fmt.Errorf("could not initialize token %v: %w", pair.DaiAddr, err)
		}
		u.Tokens[pair.DaiAddr] = token
	}

	return nil
}

// Initialize token balances for all wallets
func (u *Uniswap) InitializeTokenBalances() {
	// Initialize the 2D map
	u.tokenBalances = make(map[common.Address]map[common.Address]*big.Int)
	u.tokenBalancesMutex = sync.RWMutex{}

	// Get all wallets
	wallets := u.walletPool.GetAllWallets()

	// Initialize balances for each wallet and token
	for _, wallet := range wallets {
		walletAddr := wallet.GetAddress()

		// Initialize the inner map for this wallet
		u.tokenBalancesMutex.Lock()
		u.tokenBalances[walletAddr] = make(map[common.Address]*big.Int)
		u.tokenBalancesMutex.Unlock()

		for _, pair := range u.deploymentInfo.Pairs {
			token := u.Tokens[pair.DaiAddr]
			if token == nil {
				continue
			}

			balance, err := token.BalanceOf(&bind.CallOpts{}, walletAddr)
			if err != nil {
				u.logger.Errorf("could not get token balance for %v: %v", walletAddr, err)
				continue
			}

			u.tokenBalancesMutex.Lock()
			u.tokenBalances[walletAddr][pair.DaiAddr] = balance
			u.tokenBalancesMutex.Unlock()
		}

		// Get WETH balance
		wethBalance, err := u.Weth.BalanceOf(&bind.CallOpts{}, walletAddr)
		if err != nil {
			u.logger.Errorf("could not get WETH balance for %v: %v", walletAddr, err)
			continue
		}

		// Store WETH balance in the same map
		u.tokenBalancesMutex.Lock()
		u.tokenBalances[walletAddr][u.deploymentInfo.Weth9Addr] = wethBalance
		u.tokenBalancesMutex.Unlock()
	}
}

// Get DAI balance from local cache
func (u *Uniswap) GetTokenBalance(walletAddr common.Address, tokenAddr common.Address) *big.Int {
	u.tokenBalancesMutex.RLock()
	defer u.tokenBalancesMutex.RUnlock()

	walletBalances, exists := u.tokenBalances[walletAddr]
	if !exists {
		return big.NewInt(0)
	}

	balance, exists := walletBalances[tokenAddr]
	if !exists {
		return big.NewInt(0)
	}
	return balance
}

// Update DAI balance in local cache
func (u *Uniswap) UpdateTokenBalance(walletAddr common.Address, tokenAddr common.Address, newBalance *big.Int) {
	u.tokenBalancesMutex.Lock()
	defer u.tokenBalancesMutex.Unlock()

	// Ensure the wallet map exists
	if _, exists := u.tokenBalances[walletAddr]; !exists {
		u.tokenBalances[walletAddr] = make(map[common.Address]*big.Int)
	}

	u.tokenBalances[walletAddr][tokenAddr] = newBalance
}

func (u *Uniswap) getTxFee(ctx context.Context, client *txbuilder.Client) (*big.Int, *big.Int, error) {
	var feeCap *big.Int
	var tipCap *big.Int

	if u.options.BaseFee > 0 {
		feeCap = new(big.Int).Mul(big.NewInt(int64(u.options.BaseFee)), big.NewInt(1000000000))
	}
	if u.options.TipFee > 0 {
		tipCap = new(big.Int).Mul(big.NewInt(int64(u.options.TipFee)), big.NewInt(1000000000))
	}

	if feeCap == nil || tipCap == nil {
		var err error
		feeCap, tipCap, err = client.GetSuggestedFee(u.ctx)
		if err != nil {
			return nil, nil, err
		}
	}

	if feeCap.Cmp(big.NewInt(1000000000)) < 0 {
		feeCap = big.NewInt(1000000000)
	}
	if tipCap.Cmp(big.NewInt(1000000000)) < 0 {
		tipCap = big.NewInt(1000000000)
	}

	return feeCap, tipCap, nil
}

// Set unlimited allowances for all wallets to both routers
func (u *Uniswap) SetUnlimitedAllowances() error {
	u.logger.Infof("Setting unlimited allowances for all wallets...")

	// Get all wallets
	wallets := u.walletPool.GetAllWallets()

	// Maximum uint256 value for unlimited allowance
	maxAllowance := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))

	// Get a client for fee calculation
	client := u.walletPool.GetClient(spamoor.SelectClientByIndex, 0, u.options.ClientGroup)
	if client == nil {
		return fmt.Errorf("no client available")
	}

	feeCap, tipCap, err := u.getTxFee(u.ctx, client)
	if err != nil {
		return fmt.Errorf("could not get tx fee: %v", err)
	}

	// Track all approval transactions
	var approvalTxs []*types.Transaction
	var approvalWallets []*txbuilder.Wallet

	// For each wallet and token pair
	for _, wallet := range wallets {
		// Set allowances for DAI tokens
		for _, pair := range u.deploymentInfo.Pairs {
			token := u.Tokens[pair.DaiAddr]
			if token == nil {
				continue
			}

			// Check if allowance is already set for router A
			allowanceA, err := token.Allowance(&bind.CallOpts{}, wallet.GetAddress(), u.deploymentInfo.UniswapRouterAAddr)
			if err != nil {
				u.logger.Errorf("could not check allowance for %v: %v", wallet.GetAddress(), err)
				continue
			}

			// Check if allowance is already set for router B
			allowanceB, err := token.Allowance(&bind.CallOpts{}, wallet.GetAddress(), u.deploymentInfo.UniswapRouterBAddr)
			if err != nil {
				u.logger.Errorf("could not check allowance for %v: %v", wallet.GetAddress(), err)
				continue
			}

			// Skip if allowance is already set for both routers
			if allowanceA.Cmp(maxAllowance) >= 0 && allowanceB.Cmp(maxAllowance) >= 0 {
				continue
			}

			// Build approval transaction for router A if needed
			if allowanceA.Cmp(maxAllowance) < 0 {
				approveTx, err := wallet.BuildBoundTx(u.ctx, &txbuilder.TxMetadata{
					GasFeeCap: uint256.MustFromBig(feeCap),
					GasTipCap: uint256.MustFromBig(tipCap),
					Gas:       100000,
					Value:     uint256.NewInt(0),
				}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
					return token.Approve(transactOpts, u.deploymentInfo.UniswapRouterAAddr, maxAllowance)
				})
				if err != nil {
					u.logger.Errorf("could not build approval tx for %v: %v", wallet.GetAddress(), err)
					continue
				}

				approvalTxs = append(approvalTxs, approveTx)
				approvalWallets = append(approvalWallets, wallet)
			}

			// Build approval transaction for router B if needed
			if allowanceB.Cmp(maxAllowance) < 0 {
				approveTx, err := wallet.BuildBoundTx(u.ctx, &txbuilder.TxMetadata{
					GasFeeCap: uint256.MustFromBig(feeCap),
					GasTipCap: uint256.MustFromBig(tipCap),
					Gas:       100000,
					Value:     uint256.NewInt(0),
				}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
					return token.Approve(transactOpts, u.deploymentInfo.UniswapRouterBAddr, maxAllowance)
				})
				if err != nil {
					u.logger.Errorf("could not build approval tx for %v: %v", wallet.GetAddress(), err)
					continue
				}

				approvalTxs = append(approvalTxs, approveTx)
				approvalWallets = append(approvalWallets, wallet)
			}
		}

		// Set allowances for WETH
		wethAllowanceA, err := u.Weth.Allowance(&bind.CallOpts{}, wallet.GetAddress(), u.deploymentInfo.UniswapRouterAAddr)
		if err != nil {
			u.logger.Errorf("could not check WETH allowance for %v: %v", wallet.GetAddress(), err)
			continue
		}

		wethAllowanceB, err := u.Weth.Allowance(&bind.CallOpts{}, wallet.GetAddress(), u.deploymentInfo.UniswapRouterBAddr)
		if err != nil {
			u.logger.Errorf("could not check WETH allowance for %v: %v", wallet.GetAddress(), err)
			continue
		}

		// Skip if allowance is already set for both routers
		if wethAllowanceA.Cmp(maxAllowance) >= 0 && wethAllowanceB.Cmp(maxAllowance) >= 0 {
			continue
		}

		// Build approval transaction for router A if needed
		if wethAllowanceA.Cmp(maxAllowance) < 0 {
			approveTx, err := wallet.BuildBoundTx(u.ctx, &txbuilder.TxMetadata{
				GasFeeCap: uint256.MustFromBig(feeCap),
				GasTipCap: uint256.MustFromBig(tipCap),
				Gas:       100000,
				Value:     uint256.NewInt(0),
			}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
				return u.Weth.Approve(transactOpts, u.deploymentInfo.UniswapRouterAAddr, maxAllowance)
			})
			if err != nil {
				u.logger.Errorf("could not build WETH approval tx for %v: %v", wallet.GetAddress(), err)
				continue
			}

			approvalTxs = append(approvalTxs, approveTx)
			approvalWallets = append(approvalWallets, wallet)
		}

		// Build approval transaction for router B if needed
		if wethAllowanceB.Cmp(maxAllowance) < 0 {
			approveTx, err := wallet.BuildBoundTx(u.ctx, &txbuilder.TxMetadata{
				GasFeeCap: uint256.MustFromBig(feeCap),
				GasTipCap: uint256.MustFromBig(tipCap),
				Gas:       100000,
				Value:     uint256.NewInt(0),
			}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
				return u.Weth.Approve(transactOpts, u.deploymentInfo.UniswapRouterBAddr, maxAllowance)
			})
			if err != nil {
				u.logger.Errorf("could not build WETH approval tx for %v: %v", wallet.GetAddress(), err)
				continue
			}

			approvalTxs = append(approvalTxs, approveTx)
			approvalWallets = append(approvalWallets, wallet)
		}
	}

	// Send all approval transactions in parallel
	if len(approvalTxs) > 0 {
		u.logger.Infof("Sending %d approval transactions...", len(approvalTxs))

		// Create a wait group to track all transactions
		var wg sync.WaitGroup

		// Send each transaction to a different client
		for i, tx := range approvalTxs {
			// Get a different client for each transaction
			txClient := u.walletPool.GetClient(spamoor.SelectClientByIndex, i, u.options.ClientGroup)
			if txClient == nil {
				txClient = client
			}

			wg.Add(1)

			go func(tx *types.Transaction, client *txbuilder.Client, wallet *txbuilder.Wallet) {
				err := u.walletPool.GetTxPool().SendTransaction(u.ctx, wallet, tx, &txbuilder.SendTransactionOptions{
					Client: client,
					OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
						if err != nil {
							u.logger.Errorf("approval tx failed: %v", err)
						}
						wg.Done()
					},
					MaxRebroadcasts:     10,
					RebroadcastInterval: 30 * time.Second,
				})
				if err != nil {
					u.logger.Errorf("failed to send approval tx: %v", err)
				}
			}(tx, txClient, approvalWallets[i])
		}

		// Wait for all transactions to be sent
		wg.Wait()
		u.logger.Infof("All approval transactions sent")
	} else {
		u.logger.Infof("No approval transactions needed (allowances already set)")
	}

	return nil
}
