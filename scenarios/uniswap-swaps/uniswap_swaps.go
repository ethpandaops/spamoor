package uniswapswaps

import (
	"context"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/utils"
)

type ScenarioOptions struct {
	TotalCount        uint64  `yaml:"total_count"`
	Throughput        uint64  `yaml:"throughput"`
	MaxPending        uint64  `yaml:"max_pending"`
	MaxWallets        uint64  `yaml:"max_wallets"`
	Rebroadcast       uint64  `yaml:"rebroadcast"`
	BaseFee           float64 `yaml:"base_fee"`
	TipFee            float64 `yaml:"tip_fee"`
	PairCount         uint64  `yaml:"pair_count"`
	MinSwapAmount     string  `yaml:"min_swap_amount"`
	MaxSwapAmount     string  `yaml:"max_swap_amount"`
	BuyRatio          uint64  `yaml:"buy_ratio"`
	Slippage          uint64  `yaml:"slippage"`
	SellThreshold     string  `yaml:"sell_threshold"`
	Timeout           string  `yaml:"timeout"`
	ClientGroup       string  `yaml:"client_group"`
	DeployClientGroup string  `yaml:"deploy_client_group"`
	LogTxs            bool    `yaml:"log_txs"`
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	uniswap        *Uniswap
	deploymentInfo *DeploymentInfo
}

var ScenarioName = "uniswap-swaps"
var ScenarioDefaultOptions = ScenarioOptions{
	TotalCount:        0,
	Throughput:        10,
	MaxPending:        0,
	MaxWallets:        0,
	Rebroadcast:       1,
	BaseFee:           20,
	TipFee:            2,
	PairCount:         1,
	MinSwapAmount:     "100000000000000000",     // 0.1 DAI
	MaxSwapAmount:     "1000000000000000000000", // 1000 DAI
	BuyRatio:          40,
	Slippage:          50,
	SellThreshold:     "50000000000000000000000", // 50000 DAI
	Timeout:           "",
	ClientGroup:       "",
	DeployClientGroup: "",
	LogTxs:            false,
}
var ScenarioDescriptor = scenario.Descriptor{
	Name:           ScenarioName,
	Description:    "Send uniswap v2 swaps with different configurations",
	DefaultOptions: ScenarioDefaultOptions,
	NewScenario:    newScenario,
}

func newScenario(logger logrus.FieldLogger) scenario.Scenario {
	return &Scenario{
		options: ScenarioDefaultOptions,
		logger:  logger.WithField("scenario", ScenarioName),
	}
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
	flags.Uint64VarP(&s.options.TotalCount, "count", "c", ScenarioDefaultOptions.TotalCount, "Total number of transfer transactions to send")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t", ScenarioDefaultOptions.Throughput, "Number of transfer transactions to send per slot")
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", ScenarioDefaultOptions.Rebroadcast, "Enable reliable rebroadcast system")
	flags.Float64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use in transfer transactions (in gwei)")
	flags.Float64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use in transfer transactions (in gwei)")
	flags.Uint64Var(&s.options.PairCount, "pair-count", ScenarioDefaultOptions.PairCount, "Number of uniswap pairs to deploy")
	flags.StringVar(&s.options.MinSwapAmount, "min-swap", ScenarioDefaultOptions.MinSwapAmount, "Minimum swap amount in wei")
	flags.StringVar(&s.options.MaxSwapAmount, "max-swap", ScenarioDefaultOptions.MaxSwapAmount, "Maximum swap amount in wei")
	flags.Uint64Var(&s.options.BuyRatio, "buy-ratio", ScenarioDefaultOptions.BuyRatio, "Ratio of buy vs sell swaps (0-100)")
	flags.Uint64Var(&s.options.Slippage, "slippage", ScenarioDefaultOptions.Slippage, "Slippage tolerance in basis points")
	flags.StringVar(&s.options.SellThreshold, "sell-threshold", ScenarioDefaultOptions.SellThreshold, "DAI balance threshold to force sell (in wei)")
	flags.StringVar(&s.options.Timeout, "timeout", ScenarioDefaultOptions.Timeout, "Timeout for the scenario (e.g. '1h', '30m', '5s') - empty means no timeout")
	flags.StringVar(&s.options.ClientGroup, "client-group", ScenarioDefaultOptions.ClientGroup, "Client group to use for sending transactions")
	flags.StringVar(&s.options.DeployClientGroup, "deploy-client-group", ScenarioDefaultOptions.DeployClientGroup, "Client group to use for deployments")
	flags.BoolVar(&s.options.LogTxs, "log-txs", ScenarioDefaultOptions.LogTxs, "Log all submitted transactions")
	return nil
}

func (s *Scenario) Init(options *scenario.Options) error {
	s.walletPool = options.WalletPool

	if options.Config != "" {
		// Use the generalized config validation and parsing helper
		err := scenario.ParseAndValidateConfig(&ScenarioDescriptor, options.Config, &s.options, s.logger)
		if err != nil {
			return err
		}
	}

	if s.options.MaxWallets > 0 {
		s.walletPool.SetWalletCount(s.options.MaxWallets)
	} else if s.options.TotalCount > 0 {
		maxWallets := s.options.TotalCount / 50
		if maxWallets < 10 {
			maxWallets = 10
		} else if maxWallets > 1000 {
			maxWallets = 1000
		}

		s.walletPool.SetWalletCount(maxWallets)
	} else {
		if s.options.Throughput*10 < 1000 {
			s.walletPool.SetWalletCount(s.options.Throughput * 10)
		} else {
			s.walletPool.SetWalletCount(1000)
		}
	}

	if s.options.TotalCount == 0 && s.options.Throughput == 0 {
		return fmt.Errorf("neither total count nor throughput limit set, must define at least one of them (see --help for list of all flags)")
	}

	// register well known wallets
	s.walletPool.AddWellKnownWallet(&spamoor.WellKnownWalletConfig{
		Name:          "deployer",
		RefillAmount:  uint256.NewInt(2000000000000000000), // 2 ETH
		RefillBalance: uint256.NewInt(1000000000000000000), // 1 ETH
	})
	s.walletPool.AddWellKnownWallet(&spamoor.WellKnownWalletConfig{
		Name:          "owner",
		RefillAmount:  uint256.NewInt(1000000000000000000), // 1 ETH
		RefillBalance: uint256.NewInt(500000000000000000),  // 0.5 ETH
	})

	return nil
}

func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Infof("starting scenario: %s", ScenarioName)
	defer s.logger.Infof("scenario %s finished.", ScenarioName)

	deployClientGroup := s.options.DeployClientGroup
	if deployClientGroup == "" {
		deployClientGroup = s.options.ClientGroup
	}

	// deploy uniswap pairs
	s.uniswap = NewUniswap(ctx, s.walletPool, s.logger, UniswapOptions{
		BaseFee:             s.options.BaseFee,
		TipFee:              s.options.TipFee,
		DaiPairs:            s.options.PairCount,
		EthLiquidityPerPair: uint256.NewInt(0).Mul(uint256.NewInt(2000), uint256.NewInt(1000000000000000000)),
		DaiLiquidityFactor:  10000,
		ClientGroup:         deployClientGroup,
	})

	deploymentInfo, err := s.uniswap.DeployUniswapPairs(false)
	if err != nil {
		s.logger.Errorf("could not deploy uniswap pairs: %v", err)
		return err
	}
	if deploymentInfo == nil {
		return fmt.Errorf("could not deploy uniswap pairs: %w", err)
	}
	s.deploymentInfo = deploymentInfo

	err = s.uniswap.InitializeContracts(deploymentInfo)
	if err != nil {
		s.logger.Errorf("could not initialize uniswap contracts: %v", err)
		return err
	}

	s.uniswap.InitializeTokenBalances()

	// Set unlimited allowances for all wallets to both routers
	err = s.uniswap.SetUnlimitedAllowances()
	if err != nil {
		s.logger.Errorf("could not set unlimited allowances: %v", err)
		return err
	}

	// send transactions
	maxPending := s.options.MaxPending
	if maxPending == 0 {
		maxPending = s.options.Throughput * 10
		if maxPending == 0 {
			maxPending = 4000
		}

		if maxPending > s.walletPool.GetConfiguredWalletCount()*10 {
			maxPending = s.walletPool.GetConfiguredWalletCount() * 10
		}
	}

	// Parse timeout
	var timeout time.Duration
	if s.options.Timeout != "" {
		var err error
		timeout, err = time.ParseDuration(s.options.Timeout)
		if err != nil {
			return fmt.Errorf("invalid timeout value: %v", err)
		}
		s.logger.Infof("Timeout set to %v", timeout)
	}

	err = scenario.RunTransactionScenario(ctx, scenario.TransactionScenarioOptions{
		TotalCount:                  s.options.TotalCount,
		Throughput:                  s.options.Throughput,
		MaxPending:                  maxPending,
		ThroughputIncrementInterval: 0,
		Timeout:                     timeout,
		WalletPool:                  s.walletPool,

		Logger: s.logger,
		ProcessNextTxFn: func(ctx context.Context, params *scenario.ProcessNextTxParams) error {
			logger := s.logger
			receiptChan, tx, client, wallet, err := s.sendTx(ctx, params.TxIdx)
			if client != nil {
				logger = logger.WithField("rpc", client.GetName())
			}
			if tx != nil {
				logger = logger.WithField("nonce", tx.Nonce())
			}
			if wallet != nil {
				logger = logger.WithField("wallet", s.walletPool.GetWalletName(wallet.GetAddress()))
			}

			params.NotifySubmitted()
			params.OrderedLogCb(func() {
				if err != nil {
					logger.Warnf("could not send transaction: %v", err)
				} else if s.options.LogTxs {
					logger.Infof("sent tx #%6d: %v", params.TxIdx+1, tx.Hash().String())
				} else {
					logger.Debugf("sent tx #%6d: %v", params.TxIdx+1, tx.Hash().String())
				}
			})

			// wait for receipt
			if _, err := receiptChan.Wait(ctx); err != nil {
				return err
			}

			return err
		},
	})

	return err
}

func (s *Scenario) sendTx(ctx context.Context, txIdx uint64) (scenario.ReceiptChan, *types.Transaction, *spamoor.Client, *spamoor.Wallet, error) {
	client := s.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, int(txIdx)),
		spamoor.WithClientGroup(s.options.ClientGroup),
	)
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(txIdx))

	if client == nil {
		return nil, nil, client, wallet, scenario.ErrNoClients
	}

	if wallet == nil {
		return nil, nil, client, wallet, scenario.ErrNoWallet
	}

	if err := wallet.ResetNoncesIfNeeded(ctx, client); err != nil {
		return nil, nil, client, wallet, err
	}

	feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(client, s.options.BaseFee, s.options.TipFee)
	if err != nil {
		return nil, nil, client, wallet, err
	}

	// Select random pair
	pairIdx := mathrand.Intn(len(s.deploymentInfo.Pairs))
	pair := s.deploymentInfo.Pairs[pairIdx]

	// Parse min and max swap amounts
	minAmount, ok := new(big.Int).SetString(s.options.MinSwapAmount, 10)
	if !ok {
		return nil, nil, client, wallet, fmt.Errorf("invalid min swap amount: %s", s.options.MinSwapAmount)
	}

	maxAmount, ok := new(big.Int).SetString(s.options.MaxSwapAmount, 10)
	if !ok {
		return nil, nil, client, wallet, fmt.Errorf("invalid max swap amount: %s", s.options.MaxSwapAmount)
	}

	// Calculate random swap amount
	diff := new(big.Int).Sub(maxAmount, minAmount)
	randomAmount := new(big.Int).Add(minAmount, new(big.Int).Rand(mathrand.New(mathrand.NewSource(time.Now().UnixNano())), diff))

	// Get current token balance from cache
	tokenBalance := s.uniswap.GetTokenBalance(wallet.GetAddress(), pair.DaiAddr)

	// Get current ETH balance from wallet
	ethBalance := wallet.GetBalance()

	// Get current WETH balance
	wethBalance := s.uniswap.GetTokenBalance(wallet.GetAddress(), s.deploymentInfo.Weth9Addr)

	// Decide if we're buying or selling based on buy ratio and balances
	isBuy := mathrand.Intn(100) < int(s.options.BuyRatio)

	// Parse sell threshold
	sellThreshold, ok := new(big.Int).SetString(s.options.SellThreshold, 10)
	if !ok {
		return nil, nil, client, wallet, fmt.Errorf("invalid sell threshold: %s", s.options.SellThreshold)
	}

	// If we have a lot of DAI, force a sell to avoid depleting the pool
	if tokenBalance.Cmp(sellThreshold) > 0 {
		isBuy = false
	}

	// If we don't have enough DAI to sell, switch to buy
	if !isBuy && tokenBalance.Cmp(randomAmount) < 0 {
		isBuy = true
	}

	// Alternate between routers based on transaction index
	router := s.uniswap.RouterA
	if mathrand.Intn(100) < 50 {
		router = s.uniswap.RouterB
	}

	var tx *types.Transaction

	if isBuy {
		// Decide whether to use ETH or WETH for buying
		useWeth := mathrand.Intn(100) < 60 // 60% chance to use WETH if available

		if useWeth {
			// Buying DAI with WETH
			// Calculate how much WETH we need to spend to get the desired amount of DAI
			amounts, err := router.GetAmountsIn(&bind.CallOpts{}, randomAmount, []common.Address{s.deploymentInfo.Weth9Addr, pair.DaiAddr})
			if err != nil {
				return nil, nil, client, wallet, err
			}

			wethAmount := amounts[0]

			// Check if we have enough WETH
			if wethBalance.Cmp(wethAmount) < 0 {
				// Fall back to ETH if not enough WETH
				useWeth = false
			} else {
				// Calculate minimum DAI amount to receive (with slippage)
				minDaiAmount := new(big.Int).Mul(randomAmount, big.NewInt(10000-int64(s.options.Slippage)))
				minDaiAmount = minDaiAmount.Div(minDaiAmount, big.NewInt(10000))

				// Build buy transaction with WETH
				tx, err = wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
					GasFeeCap: uint256.MustFromBig(feeCap),
					GasTipCap: uint256.MustFromBig(tipCap),
					Gas:       200000,
					Value:     uint256.NewInt(0),
				}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
					return router.SwapExactTokensForTokens(transactOpts, wethAmount, minDaiAmount, []common.Address{s.deploymentInfo.Weth9Addr, pair.DaiAddr}, wallet.GetAddress(), big.NewInt(time.Now().Unix()+300))
				})
				if err != nil {
					return nil, nil, client, wallet, err
				}

				// Update balances in local cache
				if tx != nil {
					// Subtract WETH amount
					newWethBalance := new(big.Int).Sub(wethBalance, wethAmount)
					s.uniswap.UpdateTokenBalance(wallet.GetAddress(), s.deploymentInfo.Weth9Addr, newWethBalance)

					// Add DAI amount
					newDaiBalance := new(big.Int).Add(tokenBalance, randomAmount)
					s.uniswap.UpdateTokenBalance(wallet.GetAddress(), pair.DaiAddr, newDaiBalance)
				}
			}
		}

		// If not using WETH or not enough WETH, use ETH
		if !useWeth {
			// Buying DAI with ETH
			// Calculate how much ETH we need to spend to get the desired amount of DAI
			amounts, err := router.GetAmountsIn(&bind.CallOpts{}, randomAmount, []common.Address{s.deploymentInfo.Weth9Addr, pair.DaiAddr})
			if err != nil {
				return nil, nil, client, wallet, err
			}

			ethAmount := amounts[0]

			// Check if we have enough ETH
			if ethBalance.Cmp(ethAmount) < 0 {
				return nil, nil, client, wallet, fmt.Errorf("insufficient ETH balance for swap")
			}

			// Calculate minimum DAI amount to receive (with slippage)
			minDaiAmount := new(big.Int).Mul(randomAmount, big.NewInt(10000-int64(s.options.Slippage)))
			minDaiAmount = minDaiAmount.Div(minDaiAmount, big.NewInt(10000))

			// Build buy transaction
			tx, err = wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
				GasFeeCap: uint256.MustFromBig(feeCap),
				GasTipCap: uint256.MustFromBig(tipCap),
				Gas:       200000,
				Value:     uint256.MustFromBig(ethAmount),
			}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
				return router.SwapExactETHForTokens(transactOpts, minDaiAmount, []common.Address{s.deploymentInfo.Weth9Addr, pair.DaiAddr}, wallet.GetAddress(), big.NewInt(time.Now().Unix()+300))
			})
			if err != nil {
				return nil, nil, client, wallet, err
			}

			// Update balances in local cache
			if tx != nil {
				// Subtract ETH amount
				wallet.SubBalance(ethAmount)

				// Add DAI amount
				newDaiBalance := new(big.Int).Add(tokenBalance, randomAmount)
				s.uniswap.UpdateTokenBalance(wallet.GetAddress(), pair.DaiAddr, newDaiBalance)
			}
		}
	} else {
		// Decide whether to keep WETH or convert to ETH
		keepWeth := mathrand.Intn(100) < 30 // 30% chance to keep WETH

		if keepWeth {
			// Selling DAI for WETH
			if tokenBalance.Cmp(randomAmount) < 0 {
				return nil, nil, client, wallet, fmt.Errorf("insufficient DAI balance for swap")
			}

			// Calculate minimum WETH amount to receive (with slippage)
			amounts, err := router.GetAmountsOut(&bind.CallOpts{}, randomAmount, []common.Address{pair.DaiAddr, s.deploymentInfo.Weth9Addr})
			if err != nil {
				return nil, nil, client, wallet, err
			}
			minWethAmount := new(big.Int).Mul(amounts[1], big.NewInt(10000-int64(s.options.Slippage)))
			minWethAmount = minWethAmount.Div(minWethAmount, big.NewInt(10000))

			// Build sell transaction for WETH
			tx, err = wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
				GasFeeCap: uint256.MustFromBig(feeCap),
				GasTipCap: uint256.MustFromBig(tipCap),
				Gas:       200000,
				Value:     uint256.NewInt(0),
			}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
				return router.SwapExactTokensForTokens(transactOpts, randomAmount, minWethAmount, []common.Address{pair.DaiAddr, s.deploymentInfo.Weth9Addr}, wallet.GetAddress(), big.NewInt(time.Now().Unix()+300))
			})
			if err != nil {
				return nil, nil, client, wallet, err
			}

			// Update balances in local cache
			if tx != nil {
				// Subtract DAI amount
				newDaiBalance := new(big.Int).Sub(tokenBalance, randomAmount)
				s.uniswap.UpdateTokenBalance(wallet.GetAddress(), pair.DaiAddr, newDaiBalance)

				// Add WETH amount
				newWethBalance := new(big.Int).Add(wethBalance, amounts[1])
				s.uniswap.UpdateTokenBalance(wallet.GetAddress(), s.deploymentInfo.Weth9Addr, newWethBalance)
			}
		} else {
			// Selling DAI for ETH
			if tokenBalance.Cmp(randomAmount) < 0 {
				return nil, nil, client, wallet, fmt.Errorf("insufficient DAI balance for swap")
			}

			// Calculate minimum ETH amount to receive (with slippage)
			amounts, err := router.GetAmountsOut(&bind.CallOpts{}, randomAmount, []common.Address{pair.DaiAddr, s.deploymentInfo.Weth9Addr})
			if err != nil {
				return nil, nil, client, wallet, err
			}
			minEthAmount := new(big.Int).Mul(amounts[1], big.NewInt(10000-int64(s.options.Slippage)))
			minEthAmount = minEthAmount.Div(minEthAmount, big.NewInt(10000))

			// Build sell transaction
			tx, err = wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
				GasFeeCap: uint256.MustFromBig(feeCap),
				GasTipCap: uint256.MustFromBig(tipCap),
				Gas:       200000,
				Value:     uint256.NewInt(0),
			}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
				return router.SwapExactTokensForETH(transactOpts, randomAmount, minEthAmount, []common.Address{pair.DaiAddr, s.deploymentInfo.Weth9Addr}, wallet.GetAddress(), big.NewInt(time.Now().Unix()+300))
			})
			if err != nil {
				return nil, nil, client, wallet, err
			}

			// Update balances in local cache
			if tx != nil {
				// Subtract DAI amount
				newDaiBalance := new(big.Int).Sub(tokenBalance, randomAmount)
				s.uniswap.UpdateTokenBalance(wallet.GetAddress(), pair.DaiAddr, newDaiBalance)

				// Add ETH amount
				wallet.AddBalance(amounts[1])
			}
		}
	}

	if err != nil {
		return nil, nil, client, wallet, err
	}

	receiptChan := make(scenario.ReceiptChan, 1)
	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		ClientGroup: s.options.ClientGroup,
		Rebroadcast: s.options.Rebroadcast > 0,
		OnComplete: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			receiptChan <- receipt
		},
		OnConfirm: func(tx *types.Transaction, receipt *types.Receipt) {
			txFees := utils.GetTransactionFees(tx, receipt)
			s.logger.WithField("rpc", client.GetName()).Debugf(
				" transaction %d confirmed in block #%v. total fee: %v gwei (base: %v) logs: %v",
				txIdx+1,
				receipt.BlockNumber.String(),
				txFees.TotalFeeGweiString(),
				txFees.TxBaseFeeGweiString(),
				len(receipt.Logs),
			)
		},
		LogFn: spamoor.GetDefaultLogFn(s.logger, "", fmt.Sprintf("%6d", txIdx+1), tx),
	})
	if err != nil {
		// mark nonce as skipped if tx was not sent
		wallet.MarkSkippedNonce(tx.Nonce())

		return nil, nil, client, wallet, err
	}

	return receiptChan, tx, client, wallet, nil
}
