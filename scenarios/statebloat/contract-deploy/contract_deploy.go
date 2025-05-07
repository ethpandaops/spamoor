package contractdeploy

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"
	"gopkg.in/yaml.v3"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenarios/statebloat/contract-deploy/contract"
	"github.com/ethpandaops/spamoor/scenariotypes"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/utils"
)

type ScenarioOptions struct {
	TotalCount  uint64 `yaml:"total_count"`
	Throughput  uint64 `yaml:"throughput"`
	MaxPending  uint64 `yaml:"max_pending"`
	MaxWallets  uint64 `yaml:"max_wallets"`
	Rebroadcast uint64 `yaml:"rebroadcast"`
	BaseFee     uint64 `yaml:"base_fee"`
	TipFee      uint64 `yaml:"tip_fee"`
	GasPerBlock uint64 `yaml:"gas_per_block"`
	ClientGroup string `yaml:"client_group"`
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	pendingChan   chan bool
	pendingWGroup sync.WaitGroup
}

var ScenarioName = "contract-deploy"
var ScenarioDefaultOptions = ScenarioOptions{
	TotalCount:  0,
	Throughput:  0,
	MaxPending:  0,
	MaxWallets:  0,
	Rebroadcast: 30,
	BaseFee:     20,
	TipFee:      2,
	GasPerBlock: 0,
	ClientGroup: "",
}
var ScenarioDescriptor = scenariotypes.ScenarioDescriptor{
	Name:           ScenarioName,
	Description:    "Deploy contracts to create state bloat",
	DefaultOptions: ScenarioDefaultOptions,
	NewScenario:    newScenario,
}

func newScenario(logger logrus.FieldLogger) scenariotypes.Scenario {
	return &Scenario{
		logger: logger.WithField("scenario", ScenarioName),
	}
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
	flags.Uint64VarP(&s.options.TotalCount, "count", "c", ScenarioDefaultOptions.TotalCount, "Total number of contracts to deploy")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t", ScenarioDefaultOptions.Throughput, "Number of contracts to deploy per slot")
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", ScenarioDefaultOptions.Rebroadcast, "Number of seconds to wait before re-broadcasting a transaction")
	flags.Uint64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use in transactions (in gwei)")
	flags.Uint64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use in transactions (in gwei)")
	flags.Uint64Var(&s.options.GasPerBlock, "gas-per-block", ScenarioDefaultOptions.GasPerBlock, "Target gas to use per block (will calculate number of contracts to deploy)")
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
		walletPool.SetWalletCount(s.options.MaxWallets)
	} else if s.options.TotalCount > 0 {
		if s.options.TotalCount < 1000 {
			walletPool.SetWalletCount(s.options.TotalCount)
		} else {
			walletPool.SetWalletCount(1000)
		}
	} else {
		if s.options.Throughput*10 < 1000 {
			walletPool.SetWalletCount(s.options.Throughput * 10)
		} else {
			walletPool.SetWalletCount(1000)
		}
	}

	if s.options.TotalCount == 0 && s.options.Throughput == 0 && s.options.GasPerBlock == 0 {
		return fmt.Errorf("neither total count, throughput, nor gas per block limit set, must define at least one of them (see --help for list of all flags)")
	}

	if s.options.MaxPending > 0 {
		s.pendingChan = make(chan bool, s.options.MaxPending)
	}

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

	// Calculate throughput based on gas per block if specified
	throughput := s.options.Throughput
	if s.options.GasPerBlock > 0 {
		// Each deployment costs ~4.967M gas
		throughput = s.options.GasPerBlock / 4967200
		if throughput == 0 {
			throughput = 1
		}
		s.logger.Infof("calculated throughput: %d contracts per block (target gas: %d)", throughput, s.options.GasPerBlock)
	}

	initialRate := rate.Limit(float64(throughput) / float64(utils.SecondsPerSlot))
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
			if lastChan != nil {
				<-lastChan
				close(lastChan)
			}
			logger.Infof("sent tx #%6d: %v", txIdx+1, tx.Hash().String())
		}(txIdx, lastChan, currentChan)

		lastChan = currentChan

		count := txCount.Load() + uint64(pendingCount.Load())
		if s.options.TotalCount > 0 && count >= s.options.TotalCount {
			break
		}
	}
	s.pendingWGroup.Wait()
	s.logger.Infof("finished sending transactions, awaiting block inclusion...")

	return nil
}

func (s *Scenario) sendTx(ctx context.Context, txIdx uint64, onComplete func()) (*types.Transaction, *txbuilder.Client, *txbuilder.Wallet, error) {
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, int(txIdx), s.options.ClientGroup)
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(txIdx))

	if client == nil {
		return nil, nil, nil, fmt.Errorf("no client available")
	}

	var feeCap *big.Int
	var tipCap *big.Int

	if s.options.BaseFee > 0 {
		feeCap = new(big.Int).Mul(big.NewInt(int64(s.options.BaseFee)), big.NewInt(1000000000))
	}
	if s.options.TipFee > 0 {
		tipCap = new(big.Int).Mul(big.NewInt(int64(s.options.TipFee)), big.NewInt(1000000000))
	}

	// Generate random salt for unique contract
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	saltInt := new(big.Int).SetBytes(salt)

	// Get chain ID
	chainId, err := client.GetChainId(ctx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Deploy contract
	auth, err := bind.NewKeyedTransactorWithChainID(wallet.GetPrivateKey(), chainId)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create auth: %w", err)
	}

	auth.Nonce = big.NewInt(int64(wallet.GetNonce()))
	auth.Value = big.NewInt(0)
	auth.GasLimit = 5000000 // Enough for 24kB contract
	auth.GasPrice = feeCap
	auth.GasTipCap = tipCap

	// Get contract bytecode
	bytecode, err := contract.StateBloatTokenMetaData.GetDeployedBytecode("StateBloatToken")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get bytecode: %w", err)
	}

	// Pack constructor arguments
	packedArgs, err := contract.StateBloatTokenMetaData.Abi.Pack("", saltInt)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to pack constructor args: %w", err)
	}

	// Combine bytecode and constructor args
	fullBytecode := append(bytecode, packedArgs...)

	// Create transaction
	tx := types.NewContractCreation(auth.Nonce.Uint64(), auth.Value, auth.GasLimit, auth.GasPrice, fullBytecode)

	// Sign transaction
	signedTx, err := auth.Signer(auth.From, tx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Send transaction
	err = client.SendTransactionCtx(ctx, signedTx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	// Update nonce
	wallet.SetNonce(wallet.GetNonce() + 1)

	return signedTx, client, wallet, nil
}
