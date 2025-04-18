package uniswapswaps

import (
	"context"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"
	"gopkg.in/yaml.v3"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenariotypes"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/utils"
)

type ScenarioOptions struct {
	TotalCount    uint64 `yaml:"total_count"`
	Throughput    uint64 `yaml:"throughput"`
	MaxPending    uint64 `yaml:"max_pending"`
	MaxWallets    uint64 `yaml:"max_wallets"`
	Rebroadcast   uint64 `yaml:"rebroadcast"`
	BaseFee       uint64 `yaml:"base_fee"`
	TipFee        uint64 `yaml:"tip_fee"`
	PairCount     uint64 `yaml:"pair_count"`
	MinSwapAmount string `yaml:"min_swap_amount"`
	MaxSwapAmount string `yaml:"max_swap_amount"`
	BuyRatio      uint64 `yaml:"buy_ratio"`
	Slippage      uint64 `yaml:"slippage"`
	SellThreshold string `yaml:"sell_threshold"`
	ClientGroup   string `yaml:"client_group"`
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	uniswap        *Uniswap
	deploymentInfo *DeploymentInfo

	pendingChan   chan bool
	pendingWGroup sync.WaitGroup
}

var ScenarioName = "uniswap-swaps"
var ScenarioDefaultOptions = ScenarioOptions{
	TotalCount:    0,
	Throughput:    0,
	MaxPending:    0,
	MaxWallets:    0,
	Rebroadcast:   120,
	BaseFee:       20,
	TipFee:        2,
	PairCount:     1,
	MinSwapAmount: "100000000000000000",     // 0.1 DAI
	MaxSwapAmount: "1000000000000000000000", // 1000 DAI
	BuyRatio:      50,
	Slippage:      50,
	SellThreshold: "100000000000000000000000", // 10000 DAI
	ClientGroup:   "",
}
var ScenarioDescriptor = scenariotypes.ScenarioDescriptor{
	Name:           ScenarioName,
	Description:    "Send uniswap v2 swaps with different configurations",
	DefaultOptions: ScenarioDefaultOptions,
	NewScenario:    newScenario,
}

func newScenario(logger logrus.FieldLogger) scenariotypes.Scenario {
	return &Scenario{
		logger: logger.WithField("scenario", ScenarioName),
	}
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
	flags.Uint64VarP(&s.options.TotalCount, "count", "c", ScenarioDefaultOptions.TotalCount, "Total number of transfer transactions to send")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t", ScenarioDefaultOptions.Throughput, "Number of transfer transactions to send per slot")
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", ScenarioDefaultOptions.Rebroadcast, "Number of seconds to wait before re-broadcasting a transaction")
	flags.Uint64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use in transfer transactions (in gwei)")
	flags.Uint64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use in transfer transactions (in gwei)")
	flags.Uint64Var(&s.options.PairCount, "pair-count", ScenarioDefaultOptions.PairCount, "Number of uniswap pairs to deploy")
	flags.StringVar(&s.options.MinSwapAmount, "min-swap", ScenarioDefaultOptions.MinSwapAmount, "Minimum swap amount in wei")
	flags.StringVar(&s.options.MaxSwapAmount, "max-swap", ScenarioDefaultOptions.MaxSwapAmount, "Maximum swap amount in wei")
	flags.Uint64Var(&s.options.BuyRatio, "buy-ratio", ScenarioDefaultOptions.BuyRatio, "Ratio of buy vs sell swaps (0-100)")
	flags.Uint64Var(&s.options.Slippage, "slippage", ScenarioDefaultOptions.Slippage, "Slippage tolerance in basis points")
	flags.StringVar(&s.options.SellThreshold, "sell-threshold", ScenarioDefaultOptions.SellThreshold, "DAI balance threshold to force sell (in wei)")
	flags.StringVar(&s.options.ClientGroup, "client-group", ScenarioDefaultOptions.ClientGroup, "Client group to use for sending transactions")
	return nil
}

func (s *Scenario) Init(walletPool *spamoor.WalletPool, config string) error {
	s.walletPool = walletPool

	if config != "" {
		err := yaml.Unmarshal([]byte(config), &s.options)
		if err != nil {
			return fmt.Errorf("failed to unmarshal config: %w", err)
		}
	}

	if s.options.MaxWallets > 0 {
		s.walletPool.SetWalletCount(s.options.MaxWallets)
	} else if s.options.TotalCount > 0 {
		if s.options.TotalCount < 1000 {
			s.walletPool.SetWalletCount(s.options.TotalCount)
		} else {
			s.walletPool.SetWalletCount(1000)
		}
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

	if s.options.MaxPending > 0 {
		s.pendingChan = make(chan bool, s.options.MaxPending)
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

func (s *Scenario) Config() string {
	yamlBytes, _ := yaml.Marshal(&s.options)
	return string(yamlBytes)
}

func (s *Scenario) Run(ctx context.Context) error {
	txIdxCounter := uint64(0)
	pendingCount := atomic.Int64{}
	txCount := atomic.Uint64{}
	var lastChan chan bool

	s.logger.Infof("starting scenario: %s", ScenarioName)
	defer s.logger.Infof("scenario %s finished.", ScenarioName)

	s.uniswap = NewUniswap(ctx, s.walletPool, s.logger, UniswapOptions{
		BaseFee:             s.options.BaseFee,
		TipFee:              s.options.TipFee,
		DaiPairs:            s.options.PairCount,
		EthLiquidityPerPair: uint256.NewInt(0).Mul(uint256.NewInt(2000), uint256.NewInt(1000000000000000000)),
		DaiLiquidityFactor:  10000,
		ClientGroup:         s.options.ClientGroup,
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

	initialRate := rate.Limit(float64(s.options.Throughput) / float64(utils.SecondsPerSlot))
	if initialRate == 0 {
		initialRate = rate.Inf
	}
	limiter := rate.NewLimiter(initialRate, 1)

	for {
		if err := limiter.Wait(ctx); err != nil {
			if ctx.Err() != nil {
				break
			}

			s.logger.Debugf("rate limited: %s", err.Error())
			time.Sleep(100 * time.Millisecond)
			continue
		}

		txIdx := txIdxCounter
		txIdxCounter++

		if s.pendingChan != nil {
			// await pending transactions
			s.pendingChan <- true
		}
		pendingCount.Add(1)
		currentChan := make(chan bool, 1)

		go func(txIdx uint64, lastChan, currentChan chan bool) {
			defer func() {
				pendingCount.Add(-1)
				currentChan <- true
			}()

			logger := s.logger
			tx, client, wallet, err := s.sendTx(ctx, txIdx, func() {
				if s.pendingChan != nil {
					time.Sleep(100 * time.Millisecond)
					<-s.pendingChan
				}
			})
			if client != nil {
				logger = logger.WithField("rpc", client.GetName())
			}
			if tx != nil {
				logger = logger.WithField("nonce", tx.Nonce())
			}
			if wallet != nil {
				logger = logger.WithField("wallet", s.walletPool.GetWalletName(wallet.GetAddress()))
			}
			if lastChan != nil {
				<-lastChan
				close(lastChan)
			}
			if err != nil {
				logger.Warnf("could not send transaction: %v", err)
				return
			}

			txCount.Add(1)
			logger.Infof("sent tx #%6d: %v", txIdx+1, tx.Hash().String())
		}(txIdx, lastChan, currentChan)

		lastChan = currentChan

		count := txCount.Load() + uint64(pendingCount.Load())
		if s.options.TotalCount > 0 && count >= s.options.TotalCount {
			break
		}
	}

	<-lastChan
	close(lastChan)

	s.logger.Infof("finished sending transactions, awaiting block inclusion...")
	s.pendingWGroup.Wait()

	return nil
}

func (s *Scenario) getTxFee(ctx context.Context, client *txbuilder.Client) (*big.Int, *big.Int, error) {
	var feeCap *big.Int
	var tipCap *big.Int

	if s.options.BaseFee > 0 {
		feeCap = new(big.Int).Mul(big.NewInt(int64(s.options.BaseFee)), big.NewInt(1000000000))
	}
	if s.options.TipFee > 0 {
		tipCap = new(big.Int).Mul(big.NewInt(int64(s.options.TipFee)), big.NewInt(1000000000))
	}

	if feeCap == nil || tipCap == nil {
		var err error
		feeCap, tipCap, err = client.GetSuggestedFee(s.walletPool.GetContext())
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

func (s *Scenario) sendTx(ctx context.Context, txIdx uint64, onComplete func()) (*types.Transaction, *txbuilder.Client, *txbuilder.Wallet, error) {
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, int(txIdx), s.options.ClientGroup)
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(txIdx))
	transactionSubmitted := false

	defer func() {
		if !transactionSubmitted {
			onComplete()
		}
	}()

	if client == nil {
		return nil, client, wallet, fmt.Errorf("no client available")
	}

	feeCap, tipCap, err := s.getTxFee(ctx, client)
	if err != nil {
		return nil, client, wallet, err
	}

	// Select random pair
	pairIdx := mathrand.Intn(len(s.deploymentInfo.Pairs))
	pair := s.deploymentInfo.Pairs[pairIdx]

	// Parse min and max swap amounts
	minAmount, ok := new(big.Int).SetString(s.options.MinSwapAmount, 10)
	if !ok {
		return nil, client, wallet, fmt.Errorf("invalid min swap amount: %s", s.options.MinSwapAmount)
	}

	maxAmount, ok := new(big.Int).SetString(s.options.MaxSwapAmount, 10)
	if !ok {
		return nil, client, wallet, fmt.Errorf("invalid max swap amount: %s", s.options.MaxSwapAmount)
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
		return nil, client, wallet, fmt.Errorf("invalid sell threshold: %s", s.options.SellThreshold)
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
		useWeth := wethBalance.Cmp(randomAmount) >= 0 && mathrand.Intn(100) < 50 // 50% chance to use WETH if available

		if useWeth {
			// Buying DAI with WETH
			// Calculate how much WETH we need to spend to get the desired amount of DAI
			amounts, err := router.GetAmountsIn(&bind.CallOpts{}, randomAmount, []common.Address{s.deploymentInfo.Weth9Addr, pair.DaiAddr})
			if err != nil {
				return nil, nil, wallet, err
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
					return nil, nil, wallet, err
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
				return nil, nil, wallet, err
			}

			ethAmount := amounts[0]

			// Check if we have enough ETH
			if ethBalance.Cmp(ethAmount) < 0 {
				return nil, client, wallet, fmt.Errorf("insufficient ETH balance for swap")
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
				return nil, nil, wallet, err
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
				return nil, client, wallet, fmt.Errorf("insufficient DAI balance for swap")
			}

			// Calculate minimum WETH amount to receive (with slippage)
			amounts, err := router.GetAmountsOut(&bind.CallOpts{}, randomAmount, []common.Address{pair.DaiAddr, s.deploymentInfo.Weth9Addr})
			if err != nil {
				return nil, nil, wallet, err
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
				return nil, nil, wallet, err
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
				return nil, client, wallet, fmt.Errorf("insufficient DAI balance for swap")
			}

			// Calculate minimum ETH amount to receive (with slippage)
			amounts, err := router.GetAmountsOut(&bind.CallOpts{}, randomAmount, []common.Address{pair.DaiAddr, s.deploymentInfo.Weth9Addr})
			if err != nil {
				return nil, nil, wallet, err
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
				return nil, nil, wallet, err
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
		return nil, nil, wallet, err
	}

	rebroadcast := 0
	if s.options.Rebroadcast > 0 {
		rebroadcast = 10
	}

	s.pendingWGroup.Add(1)
	transactionSubmitted = true
	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &txbuilder.SendTransactionOptions{
		Client:              client,
		MaxRebroadcasts:     rebroadcast,
		RebroadcastInterval: time.Duration(s.options.Rebroadcast) * time.Second,
		OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			defer func() {
				onComplete()
				s.pendingWGroup.Done()
			}()

			if err != nil {
				s.logger.WithField("rpc", client.GetName()).Warnf("tx %6d: await receipt failed: %v", txIdx+1, err)
				return
			}
			if receipt == nil {
				return
			}

			effectiveGasPrice := receipt.EffectiveGasPrice
			if effectiveGasPrice == nil {
				effectiveGasPrice = big.NewInt(0)
			}
			feeAmount := new(big.Int).Mul(effectiveGasPrice, big.NewInt(int64(receipt.GasUsed)))
			totalAmount := new(big.Int).Add(tx.Value(), feeAmount)
			wallet.SubBalance(totalAmount)

			gweiTotalFee := new(big.Int).Div(feeAmount, big.NewInt(1000000000))
			gweiBaseFee := new(big.Int).Div(effectiveGasPrice, big.NewInt(1000000000))

			s.logger.WithField("rpc", client.GetName()).Debugf(" transaction %d confirmed in block #%v. total fee: %v gwei (base: %v) logs: %v", txIdx+1, receipt.BlockNumber.String(), gweiTotalFee, gweiBaseFee, len(receipt.Logs))
		},
		LogFn: func(client *txbuilder.Client, retry int, rebroadcast int, err error) {
			logger := s.logger.WithField("rpc", client.GetName())
			if retry > 0 {
				logger = logger.WithField("retry", retry)
			}
			if rebroadcast > 0 {
				logger = logger.WithField("rebroadcast", rebroadcast)
			}
			if err != nil {
				logger.Debugf("failed sending tx %6d: %v", txIdx+1, err)
			} else if retry > 0 || rebroadcast > 0 {
				logger.Debugf("successfully sent tx %6d", txIdx+1)
			}
		},
	})
	if err != nil {
		// reset nonce if tx was not sent
		wallet.ResetPendingNonce(ctx, client)

		return nil, client, wallet, err
	}

	return tx, client, wallet, nil
}
