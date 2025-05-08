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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethpandaops/spamoor/scenarios/statebloat/contract_deploy/contract"
	"github.com/ethpandaops/spamoor/scenariotypes"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/utils"
)

const (
	// GAS_PER_CONTRACT is the estimated gas cost for deploying one StateBloatToken contract
	GAS_PER_CONTRACT = 4970000 // ~4.97M gas per contract
)

type ScenarioOptions struct {
	MaxPending     uint64 `yaml:"max_pending"`
	MaxWallets     uint64 `yaml:"max_wallets"`
	Rebroadcast    uint64 `yaml:"rebroadcast"`
	BaseFee        uint64 `yaml:"base_fee"`
	TipFee         uint64 `yaml:"tip_fee"`
	GasPerBlock    uint64 `yaml:"gas_per_block"`
	ClientGroup    string `yaml:"client_group"`
	ContractsPerTx uint64 `yaml:"contracts_per_tx"`
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
	MaxPending:     0,
	MaxWallets:     0,
	Rebroadcast:    30,
	BaseFee:        20,
	TipFee:         2,
	GasPerBlock:    0,
	ClientGroup:    "default",
	ContractsPerTx: 1,
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
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", ScenarioDefaultOptions.Rebroadcast, "Number of seconds to wait before re-broadcasting a transaction")
	flags.Uint64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use in transactions (in gwei)")
	flags.Uint64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use in transactions (in gwei)")
	flags.Uint64Var(&s.options.GasPerBlock, "gas-per-block", ScenarioDefaultOptions.GasPerBlock, "Target gas to use per block (will calculate number of contracts to deploy)")
	flags.StringVar(&s.options.ClientGroup, "client-group", ScenarioDefaultOptions.ClientGroup, "Client group to use for sending transactions")
	flags.Uint64Var(&s.options.ContractsPerTx, "contracts-per-tx", ScenarioDefaultOptions.ContractsPerTx, "Number of contracts to deploy in a single transaction")
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

	if s.options.GasPerBlock == 0 && s.options.ContractsPerTx == 0 {
		return fmt.Errorf("neither gas per block limit nor contracts per tx set, must define at least one of them (see --help for list of all flags)")
	}

	if s.options.MaxWallets > 0 {
		walletPool.SetWalletCount(s.options.MaxWallets)
	} else {
		if s.options.ContractsPerTx*10 < 1000 {
			walletPool.SetWalletCount(s.options.ContractsPerTx * 10)
		} else {
			walletPool.SetWalletCount(1000)
		}
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

	// Get block gas limit and validate contract deployment gas
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, s.options.ClientGroup)
	if client == nil {
		return fmt.Errorf("no client available")
	}
	block, err := client.GetEthClient().BlockByNumber(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to get latest block: %w", err)
	}
	blockGasLimit := block.GasLimit()

	// Validate gas usage against block limit
	if s.options.GasPerBlock > blockGasLimit {
		return fmt.Errorf("gas per block (%d) exceeds block gas limit (%d)", s.options.GasPerBlock, blockGasLimit)
	}
	if s.options.ContractsPerTx*GAS_PER_CONTRACT > blockGasLimit {
		return fmt.Errorf("contracts per tx (%d) requires %d gas, exceeding block gas limit (%d)",
			s.options.ContractsPerTx, s.options.ContractsPerTx*GAS_PER_CONTRACT, blockGasLimit)
	}

	// Calculate throughput based on gas per block or contracts per tx
	throughput := s.options.ContractsPerTx
	if s.options.GasPerBlock > 0 {
		// Each deployment costs ~4.97M gas
		throughput = s.options.GasPerBlock / GAS_PER_CONTRACT
		if throughput == 0 {
			throughput = 1
		}
		s.logger.Infof("calculated throughput: %d contracts per block (target gas: %d)", throughput, s.options.GasPerBlock)
	} else {
		// Calculate gas needed for the specified number of contracts
		totalGas := GAS_PER_CONTRACT * s.options.ContractsPerTx
		s.logger.Infof("calculated gas: %d per tx for %d contracts", totalGas, throughput)
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
			}
			if err != nil {
				logger.Warnf("could not send transaction: %v", err)
				return
			}

			txCount.Add(1)
			logger.Infof("sent tx #%6d: %v", txIdx+1, tx.Hash().String())
		}(txIdx, lastChan, currentChan)

		lastChan = currentChan
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
	bytecode := common.FromHex(contract.ContractBin)

	// Get ABI
	parsed, err := contract.ContractMetaData.GetAbi()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get ABI: %w", err)
	}

	// Pack constructor arguments
	packedArgs, err := parsed.Pack("", saltInt)
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
