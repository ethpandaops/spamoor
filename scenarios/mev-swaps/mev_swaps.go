package mevswaps

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

	"github.com/ethpandaops/spamoor/scenarios/mev-swaps/contract"
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
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	deploymentInfo *DeploymentInfo

	pendingChan   chan bool
	pendingWGroup sync.WaitGroup

	// Track DAI balances locally to avoid chain calls
	// Map[walletAddress]Map[tokenAddress]balance
	daiBalances map[common.Address]map[common.Address]*big.Int
	daiMutex    sync.RWMutex

	// Reuse contract instances
	routerA *contract.UniswapV2Router02
	routerB *contract.UniswapV2Router02
	weth    *contract.WETH9
	tokens  map[common.Address]*contract.Dai
}

var ScenarioName = "mev-swaps"
var ScenarioDefaultOptions = ScenarioOptions{
	TotalCount:    0,
	Throughput:    0,
	MaxPending:    0,
	MaxWallets:    0,
	Rebroadcast:   120,
	BaseFee:       20,
	TipFee:        2,
	PairCount:     1,
	MinSwapAmount: "1000000000000000",     // 0.001 DAI
	MaxSwapAmount: "10000000000000000000", // 10 DAI
	BuyRatio:      50,
	Slippage:      50,
	SellThreshold: "100000000000000000000", // 100 DAI
}
var ScenarioDescriptor = scenariotypes.ScenarioDescriptor{
	Name:           ScenarioName,
	Description:    "Send MEV swaps with different configurations",
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

	deploymentInfo, _, err := s.deployUniswapPairs(
		ctx,
		false,
		uint256.NewInt(0).Mul(uint256.NewInt(2000), uint256.NewInt(1000000000000000000)),
		10000,
	)
	if err != nil {
		s.logger.Errorf("could not deploy uniswap pairs: %v", err)
		return err
	}
	if deploymentInfo == nil {
		return fmt.Errorf("could not deploy uniswap pairs: %w", err)
	}
	s.deploymentInfo = deploymentInfo
	s.logger.Infof("deployed uniswap pairs: %v", len(s.deploymentInfo.Pairs))

	// Initialize contract instances and token balances
	s.initializeContracts()
	s.initializeTokenBalances()

	// Set unlimited allowances for all wallets to both routers
	err = s.setUnlimitedAllowances(ctx)
	if err != nil {
		s.logger.Errorf("could not set unlimited allowances: %v", err)
		return err
	}

	// provide token liquidity to the child wallets

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

// Initialize contract instances to reuse
func (s *Scenario) initializeContracts() {
	// Initialize router A
	routerA, err := contract.NewUniswapV2Router02(s.deploymentInfo.UniswapRouterAAddr, s.walletPool.GetClient(spamoor.SelectClientByIndex, 0).GetEthClient())
	if err != nil {
		s.logger.Errorf("could not initialize router A: %v", err)
		return
	}
	s.routerA = routerA

	// Initialize router B
	routerB, err := contract.NewUniswapV2Router02(s.deploymentInfo.UniswapRouterBAddr, s.walletPool.GetClient(spamoor.SelectClientByIndex, 0).GetEthClient())
	if err != nil {
		s.logger.Errorf("could not initialize router B: %v", err)
		return
	}
	s.routerB = routerB

	// Initialize WETH9
	weth, err := contract.NewWETH9(s.deploymentInfo.Weth9Addr, s.walletPool.GetClient(spamoor.SelectClientByIndex, 0).GetEthClient())
	if err != nil {
		s.logger.Errorf("could not initialize WETH9: %v", err)
		return
	}
	s.weth = weth

	// Initialize token contracts
	s.tokens = make(map[common.Address]*contract.Dai)
	for _, pair := range s.deploymentInfo.Pairs {
		token, err := contract.NewDai(pair.DaiAddr, s.walletPool.GetClient(spamoor.SelectClientByIndex, 0).GetEthClient())
		if err != nil {
			s.logger.Errorf("could not initialize token %v: %v", pair.DaiAddr, err)
			continue
		}
		s.tokens[pair.DaiAddr] = token
	}
}

// Initialize token balances for all wallets
func (s *Scenario) initializeTokenBalances() {
	// Initialize the 2D map
	s.daiBalances = make(map[common.Address]map[common.Address]*big.Int)
	s.daiMutex = sync.RWMutex{}

	// Get all wallets
	wallets := s.walletPool.GetAllWallets()

	// Initialize balances for each wallet and token
	for _, wallet := range wallets {
		walletAddr := wallet.GetAddress()

		// Initialize the inner map for this wallet
		s.daiMutex.Lock()
		s.daiBalances[walletAddr] = make(map[common.Address]*big.Int)
		s.daiMutex.Unlock()

		for _, pair := range s.deploymentInfo.Pairs {
			token := s.tokens[pair.DaiAddr]
			if token == nil {
				continue
			}

			balance, err := token.BalanceOf(&bind.CallOpts{}, walletAddr)
			if err != nil {
				s.logger.Errorf("could not get token balance for %v: %v", walletAddr, err)
				continue
			}

			s.daiMutex.Lock()
			s.daiBalances[walletAddr][pair.DaiAddr] = balance
			s.daiMutex.Unlock()
		}

		// Get WETH balance
		wethBalance, err := s.weth.BalanceOf(&bind.CallOpts{}, walletAddr)
		if err != nil {
			s.logger.Errorf("could not get WETH balance for %v: %v", walletAddr, err)
			continue
		}

		// Store WETH balance in the same map
		s.daiMutex.Lock()
		s.daiBalances[walletAddr][s.deploymentInfo.Weth9Addr] = wethBalance
		s.daiMutex.Unlock()
	}
}

// Get DAI balance from local cache
func (s *Scenario) getDaiBalance(walletAddr common.Address, tokenAddr common.Address) *big.Int {
	s.daiMutex.RLock()
	defer s.daiMutex.RUnlock()

	walletBalances, exists := s.daiBalances[walletAddr]
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
func (s *Scenario) updateDaiBalance(walletAddr common.Address, tokenAddr common.Address, newBalance *big.Int) {
	s.daiMutex.Lock()
	defer s.daiMutex.Unlock()

	// Ensure the wallet map exists
	if _, exists := s.daiBalances[walletAddr]; !exists {
		s.daiBalances[walletAddr] = make(map[common.Address]*big.Int)
	}

	s.daiBalances[walletAddr][tokenAddr] = newBalance
}

// Set unlimited allowances for all wallets to both routers
func (s *Scenario) setUnlimitedAllowances(ctx context.Context) error {
	s.logger.Infof("Setting unlimited allowances for all wallets...")

	// Get all wallets
	wallets := s.walletPool.GetAllWallets()

	// Maximum uint256 value for unlimited allowance
	maxAllowance := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))

	// Get a client for fee calculation
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0)
	feeCap, tipCap, err := s.getTxFee(ctx, client)
	if err != nil {
		s.logger.Errorf("could not get tx fee: %v", err)
		return err
	}

	// Track all approval transactions
	var approvalTxs []*types.Transaction

	// For each wallet and token pair
	for _, wallet := range wallets {
		// Set allowances for DAI tokens
		for _, pair := range s.deploymentInfo.Pairs {
			token := s.tokens[pair.DaiAddr]
			if token == nil {
				continue
			}

			// Check if allowance is already set for router A
			allowanceA, err := token.Allowance(&bind.CallOpts{}, wallet.GetAddress(), s.deploymentInfo.UniswapRouterAAddr)
			if err != nil {
				s.logger.Errorf("could not check allowance for %v: %v", wallet.GetAddress(), err)
				continue
			}

			// Check if allowance is already set for router B
			allowanceB, err := token.Allowance(&bind.CallOpts{}, wallet.GetAddress(), s.deploymentInfo.UniswapRouterBAddr)
			if err != nil {
				s.logger.Errorf("could not check allowance for %v: %v", wallet.GetAddress(), err)
				continue
			}

			// Skip if allowance is already set for both routers
			if allowanceA.Cmp(maxAllowance) >= 0 && allowanceB.Cmp(maxAllowance) >= 0 {
				continue
			}

			// Build approval transaction for router A if needed
			if allowanceA.Cmp(maxAllowance) < 0 {
				approveTx, err := wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
					GasFeeCap: uint256.MustFromBig(feeCap),
					GasTipCap: uint256.MustFromBig(tipCap),
					Gas:       100000,
					Value:     uint256.NewInt(0),
				}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
					return token.Approve(transactOpts, s.deploymentInfo.UniswapRouterAAddr, maxAllowance)
				})
				if err != nil {
					s.logger.Errorf("could not build approval tx for %v: %v", wallet.GetAddress(), err)
					continue
				}

				approvalTxs = append(approvalTxs, approveTx)
			}

			// Build approval transaction for router B if needed
			if allowanceB.Cmp(maxAllowance) < 0 {
				approveTx, err := wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
					GasFeeCap: uint256.MustFromBig(feeCap),
					GasTipCap: uint256.MustFromBig(tipCap),
					Gas:       100000,
					Value:     uint256.NewInt(0),
				}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
					return token.Approve(transactOpts, s.deploymentInfo.UniswapRouterBAddr, maxAllowance)
				})
				if err != nil {
					s.logger.Errorf("could not build approval tx for %v: %v", wallet.GetAddress(), err)
					continue
				}

				approvalTxs = append(approvalTxs, approveTx)
			}
		}

		// Set allowances for WETH
		wethAllowanceA, err := s.weth.Allowance(&bind.CallOpts{}, wallet.GetAddress(), s.deploymentInfo.UniswapRouterAAddr)
		if err != nil {
			s.logger.Errorf("could not check WETH allowance for %v: %v", wallet.GetAddress(), err)
			continue
		}

		wethAllowanceB, err := s.weth.Allowance(&bind.CallOpts{}, wallet.GetAddress(), s.deploymentInfo.UniswapRouterBAddr)
		if err != nil {
			s.logger.Errorf("could not check WETH allowance for %v: %v", wallet.GetAddress(), err)
			continue
		}

		// Skip if allowance is already set for both routers
		if wethAllowanceA.Cmp(maxAllowance) >= 0 && wethAllowanceB.Cmp(maxAllowance) >= 0 {
			continue
		}

		// Build approval transaction for router A if needed
		if wethAllowanceA.Cmp(maxAllowance) < 0 {
			approveTx, err := wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
				GasFeeCap: uint256.MustFromBig(feeCap),
				GasTipCap: uint256.MustFromBig(tipCap),
				Gas:       100000,
				Value:     uint256.NewInt(0),
			}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
				return s.weth.Approve(transactOpts, s.deploymentInfo.UniswapRouterAAddr, maxAllowance)
			})
			if err != nil {
				s.logger.Errorf("could not build WETH approval tx for %v: %v", wallet.GetAddress(), err)
				continue
			}

			approvalTxs = append(approvalTxs, approveTx)
		}

		// Build approval transaction for router B if needed
		if wethAllowanceB.Cmp(maxAllowance) < 0 {
			approveTx, err := wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
				GasFeeCap: uint256.MustFromBig(feeCap),
				GasTipCap: uint256.MustFromBig(tipCap),
				Gas:       100000,
				Value:     uint256.NewInt(0),
			}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
				return s.weth.Approve(transactOpts, s.deploymentInfo.UniswapRouterBAddr, maxAllowance)
			})
			if err != nil {
				s.logger.Errorf("could not build WETH approval tx for %v: %v", wallet.GetAddress(), err)
				continue
			}

			approvalTxs = append(approvalTxs, approveTx)
		}
	}

	// Send all approval transactions in parallel
	if len(approvalTxs) > 0 {
		s.logger.Infof("Sending %d approval transactions...", len(approvalTxs))

		// Create a wait group to track all transactions
		var wg sync.WaitGroup

		// Send each transaction to a different client
		for i, tx := range approvalTxs {
			// Get a different client for each transaction
			txClient := s.walletPool.GetClient(spamoor.SelectClientByIndex, i)
			wg.Add(1)

			go func(tx *types.Transaction, client *txbuilder.Client) {
				err := s.walletPool.GetTxPool().SendTransaction(ctx, s.walletPool.GetRootWallet(), tx, &txbuilder.SendTransactionOptions{
					Client: client,
					OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
						if err != nil {
							s.logger.Errorf("approval tx failed: %v", err)
						}
						wg.Done()
					},
					MaxRebroadcasts:     10,
					RebroadcastInterval: 30 * time.Second,
				})
				if err != nil {
					s.logger.Errorf("failed to send approval tx: %v", err)
				}
			}(tx, txClient)
		}

		// Wait for all transactions to be sent
		wg.Wait()
		s.logger.Infof("All approval transactions sent")
	} else {
		s.logger.Infof("No approval transactions needed (allowances already set)")
	}

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
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, int(txIdx))
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(txIdx))
	transactionSubmitted := false

	defer func() {
		if !transactionSubmitted {
			onComplete()
		}
	}()

	feeCap, tipCap, err := s.getTxFee(ctx, client)
	if err != nil {
		return nil, client, wallet, err
	}

	// Select random pair
	pairIdx := mathrand.Intn(len(s.deploymentInfo.Pairs))
	pair := s.deploymentInfo.Pairs[pairIdx]

	// Get token contract from cache
	token := s.tokens[pair.DaiAddr]
	if token == nil {
		return nil, client, wallet, fmt.Errorf("token contract not initialized for %v", pair.DaiAddr)
	}

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
	tokenBalance := s.getDaiBalance(wallet.GetAddress(), pair.DaiAddr)

	// Get current ETH balance from wallet
	ethBalance := wallet.GetBalance()

	// Get current WETH balance
	wethBalance := s.getDaiBalance(wallet.GetAddress(), s.deploymentInfo.Weth9Addr)

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
	router := s.routerA
	if mathrand.Intn(100) < 50 {
		router = s.routerB
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
					s.updateDaiBalance(wallet.GetAddress(), s.deploymentInfo.Weth9Addr, newWethBalance)

					// Add DAI amount
					newDaiBalance := new(big.Int).Add(tokenBalance, randomAmount)
					s.updateDaiBalance(wallet.GetAddress(), pair.DaiAddr, newDaiBalance)
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
				s.updateDaiBalance(wallet.GetAddress(), pair.DaiAddr, newDaiBalance)
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
				s.updateDaiBalance(wallet.GetAddress(), pair.DaiAddr, newDaiBalance)

				// Add WETH amount
				newWethBalance := new(big.Int).Add(wethBalance, amounts[1])
				s.updateDaiBalance(wallet.GetAddress(), s.deploymentInfo.Weth9Addr, newWethBalance)
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
				s.updateDaiBalance(wallet.GetAddress(), pair.DaiAddr, newDaiBalance)

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
