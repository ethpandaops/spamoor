package seaporttrades

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/scenarios/seaport-trades/contract"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/utils"
)

// ScenarioOptions configures the seaport-trades scenario: the standard
// throughput/fee knobs shared by all scenarios plus the marketplace-specific
// price range, buy/sell mix, and self-seeding inventory sizes.
type ScenarioOptions struct {
	TotalCount         uint64  `yaml:"total_count"`
	Throughput         uint64  `yaml:"throughput"`
	MaxPending         uint64  `yaml:"max_pending"`
	MaxWallets         uint64  `yaml:"max_wallets"`
	Rebroadcast        uint64  `yaml:"rebroadcast"`
	BaseFee            float64 `yaml:"base_fee"`
	TipFee             float64 `yaml:"tip_fee"`
	BaseFeeWei         string  `yaml:"base_fee_wei"`
	TipFeeWei          string  `yaml:"tip_fee_wei"`
	BuyRatio           uint64  `yaml:"buy_ratio"`
	MinPrice           string  `yaml:"min_price"`
	MaxPrice           string  `yaml:"max_price"`
	SellThreshold      uint64  `yaml:"sell_threshold"`
	MarketInventory    uint64  `yaml:"market_inventory"`
	WalletInventory    uint64  `yaml:"wallet_inventory"`
	ReplenishThreshold uint64  `yaml:"replenish_threshold"`
	ReplenishBatch     uint64  `yaml:"replenish_batch"`
	Timeout            string  `yaml:"timeout"`
	ClientGroup        string  `yaml:"client_group"`
	DeployClientGroup  string  `yaml:"deploy_client_group"`
	LogTxs             bool    `yaml:"log_txs"`
}

// Scenario drives Seaport order fulfillments against a self-deployed,
// self-seeded marketplace: a mock NFT collection and stablecoin, a single market
// counterparty that signs listings and bids, and child wallets that fulfill them.
type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	deployment *DeploymentInfo
	market     *Market
}

var ScenarioName = "seaport-trades"
var ScenarioDefaultOptions = ScenarioOptions{
	TotalCount:         0,
	Throughput:         10,
	MaxPending:         0,
	MaxWallets:         0,
	Rebroadcast:        1,
	BaseFee:            20,
	TipFee:             2,
	BuyRatio:           50,
	MinPrice:           "10000000000000000",     // 0.01 coin
	MaxPrice:           "100000000000000000000", // 100 coin
	SellThreshold:      20,
	MarketInventory:    50,
	WalletInventory:    5,
	ReplenishThreshold: 10,
	ReplenishBatch:     50,
	Timeout:            "",
	ClientGroup:        "",
	DeployClientGroup:  "",
	LogTxs:             false,
}

var ScenarioDescriptor = scenario.Descriptor{
	Name:           ScenarioName,
	Description:    "Send Seaport (OpenSea) NFT marketplace buy & sell order fulfillments",
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
	flags.Uint64VarP(&s.options.TotalCount, "count", "c", ScenarioDefaultOptions.TotalCount, "Total number of fulfillment transactions to send")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t", ScenarioDefaultOptions.Throughput, "Number of fulfillment transactions to send per slot")
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", ScenarioDefaultOptions.Rebroadcast, "Enable reliable rebroadcast system")
	flags.Float64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use in transactions (in gwei)")
	flags.Float64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use in transactions (in gwei)")
	flags.StringVar(&s.options.BaseFeeWei, "basefee-wei", "", "Max fee per gas in wei (overrides --basefee for L2 sub-gwei fees)")
	flags.StringVar(&s.options.TipFeeWei, "tipfee-wei", "", "Max tip per gas in wei (overrides --tipfee for L2 sub-gwei fees)")
	flags.Uint64Var(&s.options.BuyRatio, "buy-ratio", ScenarioDefaultOptions.BuyRatio, "Ratio of buy vs sell fulfillments (0-100)")
	flags.StringVar(&s.options.MinPrice, "min-price", ScenarioDefaultOptions.MinPrice, "Minimum trade price in stablecoin wei")
	flags.StringVar(&s.options.MaxPrice, "max-price", ScenarioDefaultOptions.MaxPrice, "Maximum trade price in stablecoin wei")
	flags.Uint64Var(&s.options.SellThreshold, "sell-threshold", ScenarioDefaultOptions.SellThreshold, "Force a sell when a wallet holds more than this many NFTs")
	flags.Uint64Var(&s.options.MarketInventory, "market-inventory", ScenarioDefaultOptions.MarketInventory, "Number of NFTs minted to the market counterparty at start")
	flags.Uint64Var(&s.options.WalletInventory, "wallet-inventory", ScenarioDefaultOptions.WalletInventory, "Number of NFTs minted to each trader wallet at start")
	flags.Uint64Var(&s.options.ReplenishThreshold, "replenish-threshold", ScenarioDefaultOptions.ReplenishThreshold, "Self-mint more market NFTs when its inventory drops below this")
	flags.Uint64Var(&s.options.ReplenishBatch, "replenish-batch", ScenarioDefaultOptions.ReplenishBatch, "Number of NFTs to self-mint per replenish")
	flags.StringVar(&s.options.Timeout, "timeout", ScenarioDefaultOptions.Timeout, "Timeout for the scenario (e.g. '1h', '30m', '5s') - empty means no timeout")
	flags.StringVar(&s.options.ClientGroup, "client-group", ScenarioDefaultOptions.ClientGroup, "Client group to use for sending transactions")
	flags.StringVar(&s.options.DeployClientGroup, "deploy-client-group", ScenarioDefaultOptions.DeployClientGroup, "Client group to use for deployments")
	flags.BoolVar(&s.options.LogTxs, "log-txs", ScenarioDefaultOptions.LogTxs, "Log all submitted transactions")
	return nil
}

func (s *Scenario) Init(options *scenario.Options) error {
	s.walletPool = options.WalletPool

	if options.Config != "" {
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

	// The deployer funds the contract deployments; the market is the standing
	// counterparty that signs every order and holds the NFT/coin float.
	s.walletPool.AddWellKnownWallet(&spamoor.WellKnownWalletConfig{
		Name:          "deployer",
		RefillAmount:  uint256.NewInt(2000000000000000000), // 2 ETH
		RefillBalance: uint256.NewInt(1000000000000000000), // 1 ETH
	})
	s.walletPool.AddWellKnownWallet(&spamoor.WellKnownWalletConfig{
		Name:          "market",
		RefillAmount:  uint256.NewInt(2000000000000000000), // 2 ETH
		RefillBalance: uint256.NewInt(1000000000000000000), // 1 ETH
	})

	return nil
}

func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Infof("starting scenario: %s", ScenarioName)
	defer s.logger.Infof("scenario %s finished.", ScenarioName)

	deployment, err := s.DeployContracts(ctx)
	if err != nil {
		s.logger.Errorf("could not deploy seaport contracts: %v", err)
		return err
	}
	s.deployment = deployment

	market, err := NewMarket(s.walletPool, s.logger, &s.options, deployment)
	if err != nil {
		return err
	}
	s.market = market

	if err := market.Seed(ctx); err != nil {
		s.logger.Errorf("could not seed seaport market: %v", err)
		return err
	}

	// Validate the off-chain order hashing against the on-chain Seaport once, so a
	// signature mismatch fails fast instead of every fulfillment reverting.
	if err := s.verifyOrderHashing(ctx); err != nil {
		return err
	}

	// Background self-mint: keep the market's NFT inventory and coin float topped
	// up so the run never stalls for lack of something to trade.
	go s.runReplenisher(ctx)

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

	var timeout time.Duration
	if s.options.Timeout != "" {
		timeout, err = time.ParseDuration(s.options.Timeout)
		if err != nil {
			return fmt.Errorf("invalid timeout value: %v", err)
		}
		s.logger.Infof("Timeout set to %v", timeout)
	}

	return scenario.RunTransactionScenario(ctx, scenario.TransactionScenarioOptions{
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

			if receiptChan != nil {
				if _, err := receiptChan.Wait(ctx); err != nil {
					return err
				}
			}
			return err
		},
	})
}

// sendTx builds and submits one fulfillment for the given tx index, wiring the
// trade's inventory reservations to be committed or rolled back when the tx
// completes.
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

	baseFeeWei, tipFeeWei := spamoor.ResolveFees(s.options.BaseFee, s.options.TipFee, s.options.BaseFeeWei, s.options.TipFeeWei)
	feeCap, tipCap, err := s.walletPool.GetSuggestedFees(client, baseFeeWei, tipFeeWei)
	if err != nil {
		return nil, nil, client, wallet, err
	}

	tx, onResult, err := s.market.BuildTrade(ctx, wallet, feeCap, tipCap)
	if err != nil {
		return nil, nil, client, wallet, err
	}

	receiptChan := make(scenario.ReceiptChan, 1)
	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		ClientGroup: s.options.ClientGroup,
		Rebroadcast: s.options.Rebroadcast > 0,
		OnComplete: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			onResult(err == nil && receipt != nil && receipt.Status == types.ReceiptStatusSuccessful)
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
		// The tx never went out: roll back the reservations and free the nonce.
		onResult(false)
		wallet.MarkSkippedNonce(tx.Nonce())
		return nil, nil, client, wallet, err
	}

	return receiptChan, tx, client, wallet, nil
}

// runReplenisher periodically tops up the market's NFT inventory and coin float
// until the scenario context is cancelled.
func (s *Scenario) runReplenisher(ctx context.Context) {
	ticker := time.NewTicker(scenario.GlobalSlotDuration)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			client := s.walletPool.GetClient(
				spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0),
				spamoor.WithClientGroup(s.deployClientGroup()),
			)
			if client == nil {
				continue
			}
			baseFeeWei, tipFeeWei := spamoor.ResolveFees(s.options.BaseFee, s.options.TipFee, s.options.BaseFeeWei, s.options.TipFeeWei)
			feeCap, tipCap, err := s.walletPool.GetSuggestedFees(client, baseFeeWei, tipFeeWei)
			if err != nil {
				continue
			}
			s.market.maybeReplenish(ctx, client, feeCap, tipCap)
		}
	}
}

// deployClientGroup returns the client group used for deployments, falling back
// to the general client group when unset.
func (s *Scenario) deployClientGroup() string {
	if s.options.DeployClientGroup != "" {
		return s.options.DeployClientGroup
	}
	return s.options.ClientGroup
}

// verifyOrderHashing validates the off-chain Seaport EIP-712 order hashing
// against the on-chain contract once at startup, mirroring the safe-multisig
// self-check. A sample market listing is hashed both ways; a mismatch means the
// signing path would produce invalid signatures, so it fails fast.
func (s *Scenario) verifyOrderHashing(ctx context.Context) error {
	sample := contract.OrderParameters{
		Offerer: s.market.marketAddr,
		Zone:    common.Address{},
		Offer: []contract.OfferItem{{
			ItemType:             itemTypeERC721,
			Token:                s.deployment.NFTAddr,
			IdentifierOrCriteria: big.NewInt(0),
			StartAmount:          big.NewInt(1),
			EndAmount:            big.NewInt(1),
		}},
		Consideration: []contract.ConsiderationItem{{
			ItemType:             itemTypeERC20,
			Token:                s.deployment.CoinAddr,
			IdentifierOrCriteria: big.NewInt(0),
			StartAmount:          big.NewInt(1),
			EndAmount:            big.NewInt(1),
			Recipient:            s.market.marketAddr,
		}},
		OrderType:                       orderTypeFullOpen,
		StartTime:                       big.NewInt(0),
		EndTime:                         maxOrderEndTime,
		ZoneHash:                        zeroBytes32,
		Salt:                            big.NewInt(1),
		ConduitKey:                      zeroBytes32,
		TotalOriginalConsiderationItems: big.NewInt(1),
	}

	goHash, err := computeOrderHash(sample, s.market.counter)
	if err != nil {
		return err
	}
	onchain, err := s.deployment.Seaport.GetOrderHash(&bind.CallOpts{Context: ctx}, toOrderComponents(sample, s.market.counter))
	if err != nil {
		return fmt.Errorf("could not read on-chain order hash: %w", err)
	}
	if onchain != goHash {
		return fmt.Errorf("seaport order hash mismatch: on-chain %x vs computed %x", onchain, goHash)
	}
	s.logger.Debugf("validated off-chain seaport order hashing against on-chain contract")
	return nil
}
