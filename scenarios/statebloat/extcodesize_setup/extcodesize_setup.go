package extcodesizesetup

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	geas "github.com/fjl/geas/asm"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/utils"
)

// Contract sizes in KB - default set for backwards compatibility
var DefaultContractSizesKB = []float64{0.5, 1, 2, 5, 10, 24}

// AllContractSizesKB includes all supported sizes (including post-EIP-7907 sizes)
var AllContractSizesKB = []float64{0.5, 1, 2, 5, 10, 24, 32, 64}

// MaxContractSize is the maximum contract size in bytes
// Pre-Prague: 24KB (EIP-170)
// Post-Prague: 48KB (EIP-7907), but testnets may allow larger
const MaxContractSize = 65536 // 64KB for testnet support

// InitcodeInfo stores information about a deployed initcode contract
type InitcodeInfo struct {
	SizeKB       float64        `yaml:"size_kb" json:"size_kb"`
	Address      common.Address `yaml:"address" json:"address"`
	Hash         common.Hash    `yaml:"hash" json:"hash"`
	BytecodeSize int            `yaml:"bytecode_size" json:"bytecode_size"`
}

// FactoryInfo stores information about a deployed factory contract
type FactoryInfo struct {
	SizeKB          float64        `yaml:"size_kb" json:"size_kb"`
	Address         common.Address `yaml:"address" json:"address"`
	InitcodeAddress common.Address `yaml:"initcode_address" json:"initcode_address"`
	InitcodeHash    common.Hash    `yaml:"initcode_hash" json:"initcode_hash"`
	InitcodeSize    int            `yaml:"initcode_size" json:"initcode_size"`
}

// PredeployedAddresses allows users to provide pre-deployed contract addresses
type PredeployedAddresses struct {
	Initcode  map[string]common.Address `yaml:"initcode" json:"initcode"`   // e.g., "0_5kb": "0x..."
	Factories map[string]common.Address `yaml:"factories" json:"factories"` // e.g., "0_5kb": "0x..."
}

type ScenarioOptions struct {
	ContractsPerSize     uint64               `yaml:"contracts_per_size" json:"contracts_per_size"`
	ContractSizes        []float64            `yaml:"contract_sizes" json:"contract_sizes"` // Sizes in KB to deploy (e.g., [0.5, 1, 2, 5, 10, 24, 32, 64])
	MaxWallets           uint64               `yaml:"max_wallets" json:"max_wallets"`
	MaxPending           uint64               `yaml:"max_pending" json:"max_pending"` // Maximum pending transactions (0 = auto-calculate based on wallet count)
	BaseFee              float64              `yaml:"base_fee" json:"base_fee"`
	TipFee               float64              `yaml:"tip_fee" json:"tip_fee"`
	Throughput           uint64               `yaml:"throughput" json:"throughput"`
	PredeployedAddresses PredeployedAddresses `yaml:"predeployed" json:"predeployed"`
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	// Deployed contract addresses
	initcodeContracts map[string]*InitcodeInfo
	factoryContracts  map[string]*FactoryInfo
}

var ScenarioName = "extcodesize_setup"
var ScenarioDefaultOptions = ScenarioOptions{
	ContractsPerSize: 1000,
	ContractSizes:    AllContractSizesKB, // [0.5, 1, 2, 5, 10, 24, 32, 64]
	MaxWallets:       50,
	MaxPending:       0, // 0 = auto-calculate based on wallet count
	BaseFee:          20,
	TipFee:           2,
	Throughput:       0, // 0 = auto-calculate based on block gas limit
	PredeployedAddresses: PredeployedAddresses{
		Initcode:  make(map[string]common.Address),
		Factories: make(map[string]common.Address),
	},
}
var ScenarioDescriptor = scenario.Descriptor{
	Name:           ScenarioName,
	Description:    "Deploy extcodesize benchmark setup (initcode contracts, factories, and test contracts)",
	DefaultOptions: ScenarioDefaultOptions,
	NewScenario:    newScenario,
}

func newScenario(logger logrus.FieldLogger) scenario.Scenario {
	return &Scenario{
		options:           ScenarioDefaultOptions,
		logger:            logger.WithField("scenario", ScenarioName),
		initcodeContracts: make(map[string]*InitcodeInfo),
		factoryContracts:  make(map[string]*FactoryInfo),
	}
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
	flags.Uint64VarP(&s.options.ContractsPerSize, "contracts-per-size", "c", ScenarioDefaultOptions.ContractsPerSize, "Number of contracts to deploy per size")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of wallets to use for parallel execution")
	flags.Float64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Base fee per gas in gwei")
	flags.Float64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Tip fee per gas in gwei")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t", ScenarioDefaultOptions.Throughput, "Max deployments per slot (0=auto based on block gas limit)")
	return nil
}

func (s *Scenario) Init(options *scenario.Options) error {
	s.walletPool = options.WalletPool

	if options.Config != "" {
		err := yaml.Unmarshal([]byte(options.Config), &s.options)
		if err != nil {
			return fmt.Errorf("failed to unmarshal config: %w", err)
		}
	}

	// Initialize maps if nil (in case config didn't have them)
	if s.options.PredeployedAddresses.Initcode == nil {
		s.options.PredeployedAddresses.Initcode = make(map[string]common.Address)
	}
	if s.options.PredeployedAddresses.Factories == nil {
		s.options.PredeployedAddresses.Factories = make(map[string]common.Address)
	}

	// Use default contract sizes if not specified
	if len(s.options.ContractSizes) == 0 {
		s.options.ContractSizes = AllContractSizesKB
	}

	// Set up wallets for parallel execution
	if s.options.MaxWallets > 0 {
		s.walletPool.SetWalletCount(s.options.MaxWallets)
	} else {
		s.walletPool.SetWalletCount(50) // Default fallback
	}

	return nil
}

func (s *Scenario) Config() string {
	type ConfigWithState struct {
		ScenarioOptions
		InitcodeContracts map[string]*InitcodeInfo `yaml:"initcode_contracts,omitempty" json:"initcode_contracts,omitempty"`
		FactoryContracts  map[string]*FactoryInfo  `yaml:"factory_contracts,omitempty" json:"factory_contracts,omitempty"`
	}

	cfg := ConfigWithState{
		ScenarioOptions:   s.options,
		InitcodeContracts: s.initcodeContracts,
		FactoryContracts:  s.factoryContracts,
	}

	yamlBytes, _ := yaml.Marshal(&cfg)
	return string(yamlBytes)
}

func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Infof("starting scenario: %s", ScenarioName)
	defer s.logger.Infof("scenario %s finished", ScenarioName)

	// Phase 1: Deploy initcode contracts for all sizes
	s.logger.Info("=== Phase 1: Deploying initcode contracts ===")
	for _, sizeKB := range s.options.ContractSizes {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		sizeKey := sizeKeyFromKB(sizeKB)

		// Check if pre-deployed address is provided
		if predeployed, ok := s.options.PredeployedAddresses.Initcode[sizeKey]; ok && predeployed != (common.Address{}) {
			s.logger.Infof("using pre-deployed initcode for %sKB at %s", sizeKey, predeployed.Hex())
			// We need to fetch the code to compute the hash
			info, err := s.loadInitcodeInfo(ctx, sizeKB, predeployed)
			if err != nil {
				return fmt.Errorf("failed to load pre-deployed initcode for %sKB: %w", sizeKey, err)
			}
			s.initcodeContracts[sizeKey] = info
			continue
		}

		info, err := s.deployInitcode(ctx, sizeKB)
		if err != nil {
			return fmt.Errorf("failed to deploy initcode for %sKB: %w", sizeKey, err)
		}
		s.initcodeContracts[sizeKey] = info
		s.logger.Infof("deployed initcode for %.1fKB at %s (hash: %s)", sizeKB, info.Address.Hex(), info.Hash.Hex()[:18]+"...")
	}

	// Phase 2: Deploy factory contracts for all sizes
	s.logger.Info("=== Phase 2: Deploying factory contracts ===")
	for _, sizeKB := range s.options.ContractSizes {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		sizeKey := sizeKeyFromKB(sizeKB)
		initcodeInfo := s.initcodeContracts[sizeKey]

		// Check if pre-deployed factory address is provided
		if predeployed, ok := s.options.PredeployedAddresses.Factories[sizeKey]; ok && predeployed != (common.Address{}) {
			s.logger.Infof("using pre-deployed factory for %sKB at %s", sizeKey, predeployed.Hex())
			s.factoryContracts[sizeKey] = &FactoryInfo{
				SizeKB:          sizeKB,
				Address:         predeployed,
				InitcodeAddress: initcodeInfo.Address,
				InitcodeHash:    initcodeInfo.Hash,
				InitcodeSize:    initcodeInfo.BytecodeSize,
			}
			continue
		}

		factoryInfo, err := s.deployFactory(ctx, sizeKB, initcodeInfo)
		if err != nil {
			return fmt.Errorf("failed to deploy factory for %sKB: %w", sizeKey, err)
		}
		s.factoryContracts[sizeKey] = factoryInfo
		s.logger.Infof("deployed factory for %.1fKB at %s", sizeKB, factoryInfo.Address.Hex())
	}

	// Phase 3: Deploy contracts via factories
	s.logger.Info("=== Phase 3: Deploying contracts via factories ===")
	for _, sizeKB := range s.options.ContractSizes {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		sizeKey := sizeKeyFromKB(sizeKB)
		factoryInfo := s.factoryContracts[sizeKey]

		err := s.deployContractsViaFactory(ctx, sizeKB, factoryInfo)
		if err != nil {
			return fmt.Errorf("failed to deploy contracts for %sKB: %w", sizeKey, err)
		}
	}

	// Log summary in stubs.json format for execution-spec-tests compatibility
	s.logger.Info("=== Deployment Summary ===")
	s.logger.Info("Factory addresses (stubs.json format):")
	s.logger.Info("{")
	for i, sizeKB := range s.options.ContractSizes {
		sizeKey := sizeKeyFromKB(sizeKB)
		factoryInfo := s.factoryContracts[sizeKey]
		stubKey := fmt.Sprintf("bloatnet_factory_%s", sizeKey)
		comma := ","
		if i == len(s.options.ContractSizes)-1 {
			comma = ""
		}
		s.logger.Infof("  \"%s\": \"%s\"%s", stubKey, factoryInfo.Address.Hex(), comma)
	}
	s.logger.Info("}")

	// Also log in a more readable format
	s.logger.Info("")
	s.logger.Info("=== Factory Addresses ===")
	for _, sizeKB := range s.options.ContractSizes {
		sizeKey := sizeKeyFromKB(sizeKB)
		factoryInfo := s.factoryContracts[sizeKey]
		stubKey := fmt.Sprintf("bloatnet_factory_%s", sizeKey)
		s.logger.Infof("%s: %s", stubKey, factoryInfo.Address.Hex())
	}

	return nil
}

// sizeKeyFromKB converts a size in KB to a string key (e.g., 0.5 -> "0_5kb", 1.0 -> "1kb")
func sizeKeyFromKB(sizeKB float64) string {
	if sizeKB == float64(int(sizeKB)) {
		return fmt.Sprintf("%dkb", int(sizeKB))
	}
	return strings.ReplaceAll(fmt.Sprintf("%.1fkb", sizeKB), ".", "_")
}

// loadInitcodeInfo loads info for a pre-deployed initcode contract
func (s *Scenario) loadInitcodeInfo(ctx context.Context, sizeKB float64, address common.Address) (*InitcodeInfo, error) {
	client := s.walletPool.GetClient(spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0))
	if client == nil {
		return nil, fmt.Errorf("no client available")
	}

	code, err := client.GetEthClient().CodeAt(ctx, address, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get code: %w", err)
	}

	if len(code) == 0 {
		return nil, fmt.Errorf("no code at address %s", address.Hex())
	}

	hash := crypto.Keccak256Hash(code)

	return &InitcodeInfo{
		SizeKB:       sizeKB,
		Address:      address,
		Hash:         hash,
		BytecodeSize: len(code),
	}, nil
}

// deployInitcode deploys an initcode contract for a specific size
func (s *Scenario) deployInitcode(ctx context.Context, sizeKB float64) (*InitcodeInfo, error) {
	client := s.walletPool.GetClient(spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0))
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, 0)

	if client == nil {
		return nil, fmt.Errorf("no client available")
	}
	if wallet == nil {
		return nil, fmt.Errorf("no wallet available")
	}

	// Build initcode bytecode
	initcodeBytecode, err := buildInitcode(sizeKB)
	if err != nil {
		return nil, fmt.Errorf("failed to build initcode: %w", err)
	}

	initcodeHash := crypto.Keccak256Hash(initcodeBytecode)
	s.logger.Debugf("built initcode for %.1fKB: %d bytes, hash: %s", sizeKB, len(initcodeBytecode), initcodeHash.Hex()[:18]+"...")

	// Build deployment bytecode that returns the initcode as contract code
	deploymentBytecode := buildInitcodeDeployer(initcodeBytecode)

	feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(client, s.options.BaseFee, s.options.TipFee)
	if err != nil {
		return nil, fmt.Errorf("failed to get suggested fees: %w", err)
	}

	txData, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       16_000_000, // High gas for large initcode
		To:        nil,
		Value:     uint256.NewInt(0),
		Data:      deploymentBytecode,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to build tx data: %w", err)
	}

	tx, err := wallet.BuildDynamicFeeTx(txData)
	if err != nil {
		return nil, fmt.Errorf("failed to build tx: %w", err)
	}

	receipt, err := s.walletPool.GetTxPool().SendAndAwaitTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		Rebroadcast: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to send tx: %w", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, fmt.Errorf("deployment failed with status %d", receipt.Status)
	}

	return &InitcodeInfo{
		SizeKB:       sizeKB,
		Address:      receipt.ContractAddress,
		Hash:         initcodeHash,
		BytecodeSize: len(initcodeBytecode),
	}, nil
}

// deployFactory deploys a CREATE2 factory contract for a specific size
func (s *Scenario) deployFactory(ctx context.Context, sizeKB float64, initcodeInfo *InitcodeInfo) (*FactoryInfo, error) {
	client := s.walletPool.GetClient(spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0))
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, 0)

	if client == nil {
		return nil, fmt.Errorf("no client available")
	}
	if wallet == nil {
		return nil, fmt.Errorf("no wallet available")
	}

	// Build factory bytecode
	factoryBytecode, err := buildFactory(initcodeInfo.Address, initcodeInfo.Hash, initcodeInfo.BytecodeSize)
	if err != nil {
		return nil, fmt.Errorf("failed to build factory: %w", err)
	}

	s.logger.Debugf("built factory bytecode for %.1fKB: %d bytes", sizeKB, len(factoryBytecode))

	feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(client, s.options.BaseFee, s.options.TipFee)
	if err != nil {
		return nil, fmt.Errorf("failed to get suggested fees: %w", err)
	}

	txData, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       10_000_000,
		To:        nil,
		Value:     uint256.NewInt(0),
		Data:      factoryBytecode,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to build tx data: %w", err)
	}

	tx, err := wallet.BuildDynamicFeeTx(txData)
	if err != nil {
		return nil, fmt.Errorf("failed to build tx: %w", err)
	}

	receipt, err := s.walletPool.GetTxPool().SendAndAwaitTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		Rebroadcast: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to send tx: %w", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, fmt.Errorf("deployment failed with status %d", receipt.Status)
	}

	// Verify factory storage
	counter, err := client.GetEthClient().StorageAt(ctx, receipt.ContractAddress, common.Hash{}, nil)
	if err != nil {
		s.logger.Warnf("failed to verify factory storage: %v", err)
	} else {
		s.logger.Debugf("factory storage slot 0 (counter): %s", hex.EncodeToString(counter))
	}

	return &FactoryInfo{
		SizeKB:          sizeKB,
		Address:         receipt.ContractAddress,
		InitcodeAddress: initcodeInfo.Address,
		InitcodeHash:    initcodeInfo.Hash,
		InitcodeSize:    initcodeInfo.BytecodeSize,
	}, nil
}

// deployContractsViaFactory deploys multiple contracts via a CREATE2 factory using parallel transactions
func (s *Scenario) deployContractsViaFactory(ctx context.Context, sizeKB float64, factoryInfo *FactoryInfo) error {
	sizeKey := sizeKeyFromKB(sizeKB)

	// Get current counter from factory
	client := s.walletPool.GetClient(spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0))
	if client == nil {
		return fmt.Errorf("no client available")
	}

	counterBytes, err := client.GetEthClient().StorageAt(ctx, factoryInfo.Address, common.Hash{}, nil)
	if err != nil {
		return fmt.Errorf("failed to get factory counter: %w", err)
	}
	currentCounter := new(big.Int).SetBytes(counterBytes).Uint64()

	if currentCounter >= s.options.ContractsPerSize {
		s.logger.Infof("factory for %s already has %d contracts (target: %d), skipping", sizeKey, currentCounter, s.options.ContractsPerSize)
		return nil
	}

	remaining := s.options.ContractsPerSize - currentCounter

	// Get block gas limit to calculate optimal throughput
	blockGasLimit, err := s.walletPool.GetTxPool().GetCurrentGasLimitWithInit()
	if err != nil {
		return fmt.Errorf("failed to get block gas limit: %w", err)
	}

	// Calculate gas per deployment based on contract size
	gasPerDeployment := calculateGasForSize(sizeKB, blockGasLimit)

	// Calculate how many transactions can fit per block (use 90% of block gas limit for safety)
	usableGas := uint64(float64(blockGasLimit) * 0.9)
	txsPerBlock := usableGas / gasPerDeployment
	if txsPerBlock < 1 {
		txsPerBlock = 1
	}

	s.logger.Infof("deploying %d contracts for %s: block_gas=%dM, tx_gas=%d, txs_per_block=%d",
		remaining, sizeKey, blockGasLimit/1_000_000, gasPerDeployment, txsPerBlock)

	// Calculate maxPending - use config value if set, otherwise auto-calculate
	maxPending := s.options.MaxPending
	if maxPending == 0 {
		maxPending = s.walletPool.GetConfiguredWalletCount() * 10
		if maxPending < 100 {
			maxPending = 100
		}
	}

	err = scenario.RunTransactionScenario(ctx, scenario.TransactionScenarioOptions{
		TotalCount: remaining,
		Throughput: s.options.Throughput, // 0 = unlimited, send as fast as possible
		MaxPending: maxPending,
		WalletPool: s.walletPool,
		Logger:     s.logger.WithField("size", sizeKey),
		ProcessNextTxFn: func(ctx context.Context, params *scenario.ProcessNextTxParams) error {
			return s.sendFactoryDeployTx(ctx, params, factoryInfo, gasPerDeployment, sizeKey)
		},
	})
	if err != nil {
		return fmt.Errorf("failed to deploy contracts: %w", err)
	}

	// Verify final count
	counterBytes, err = client.GetEthClient().StorageAt(ctx, factoryInfo.Address, common.Hash{}, nil)
	if err != nil {
		s.logger.Warnf("failed to verify final counter: %v", err)
	} else {
		finalCounter := new(big.Int).SetBytes(counterBytes).Uint64()
		s.logger.Infof("%s: final contract count: %d", sizeKey, finalCounter)
	}

	return nil
}

// sendFactoryDeployTx sends a transaction to deploy a contract via the factory
func (s *Scenario) sendFactoryDeployTx(ctx context.Context, params *scenario.ProcessNextTxParams, factoryInfo *FactoryInfo, gasLimit uint64, sizeKey string) error {
	// Get client and wallet using txIdx for distribution across wallets
	client := s.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, int(params.TxIdx)),
	)
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(params.TxIdx))

	if client == nil {
		return scenario.ErrNoClients
	}
	if wallet == nil {
		return scenario.ErrNoWallet
	}

	if err := wallet.ResetNoncesIfNeeded(ctx, client); err != nil {
		return err
	}

	feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(client, s.options.BaseFee, s.options.TipFee)
	if err != nil {
		return err
	}

	// Build transaction to call factory (non-empty calldata triggers CREATE2)
	txData, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       gasLimit,
		To:        &factoryInfo.Address,
		Value:     uint256.NewInt(0),
		Data:      []byte{0x01}, // Non-empty data triggers CREATE2 path
	})
	if err != nil {
		return err
	}

	tx, err := wallet.BuildDynamicFeeTx(txData)
	if err != nil {
		return err
	}

	// Send transaction and wait for receipt
	receipt, err := s.walletPool.GetTxPool().SendAndAwaitTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		Rebroadcast: true,
	})

	params.NotifySubmitted()

	logger := s.logger.WithFields(logrus.Fields{
		"size":   sizeKey,
		"rpc":    client.GetName(),
		"nonce":  tx.Nonce(),
		"wallet": s.walletPool.GetWalletName(wallet.GetAddress()),
	})

	params.OrderedLogCb(func() {
		if err != nil {
			logger.Warnf("tx %d failed to send: %v", params.TxIdx+1, err)
		} else {
			logger.Debugf("sent tx #%d: %v", params.TxIdx+1, tx.Hash().Hex())
		}
	})

	if receipt != nil && receipt.Status != 1 {
		logger.Warnf("tx %d failed with status %d", params.TxIdx+1, receipt.Status)
	}

	// Mark wallet for nonce resync if transaction failed to send
	if err != nil {
		wallet.MarkNeedResync()
	}

	return err
}

// calculateGasForSize calculates gas needed to deploy a contract of given size
// Gas costs breakdown:
// - Transaction intrinsic: 21,000
// - Factory execution overhead: ~5,000
// - CREATE2 overhead: ~32,000
// - Contract bytecode storage: 200 gas per byte
// - Init code execution: scales with size (keccak256 expansion loops)
func calculateGasForSize(sizeKB float64, blockGasLimit uint64) uint64 {
	sizeBytes := uint64(sizeKB * 1024)

	// Base costs
	intrinsicGas := uint64(21_000)
	factoryOverhead := uint64(5_000)
	create2Overhead := uint64(32_000)
	bytecodeStorageCost := sizeBytes * 200

	// Init code execution cost scales with size due to keccak256 expansion loops
	// Larger contracts require more iterations to generate unique bytecode
	var initExecutionCost uint64
	if sizeKB < 1 {
		// Small contracts: simple padding, minimal execution
		initExecutionCost = uint64(sizeKB * 30_000)
	} else {
		// Larger contracts: keccak256 expansion loop
		// Each iteration does SHA3 + XOR operations
		// Roughly: (size / 256) iterations * ~3000 gas per iteration
		iterations := sizeBytes / 256
		if iterations < 1 {
			iterations = 1
		}
		initExecutionCost = iterations * 3_000
	}

	// Total gas with 15% buffer for safety
	totalGas := uint64(float64(intrinsicGas+factoryOverhead+create2Overhead+bytecodeStorageCost+initExecutionCost) * 1.15)

	// Cap at block gas limit (single tx shouldn't exceed block limit)
	if totalGas > blockGasLimit {
		totalGas = blockGasLimit
	}

	// Also cap at EIP-7825 limit if applicable
	if totalGas > utils.MaxGasLimitPerTx {
		totalGas = utils.MaxGasLimitPerTx
	}

	return totalGas
}

// buildInitcode builds initcode bytecode that generates contracts of specific size
func buildInitcode(targetSizeKB float64) ([]byte, error) {
	targetSize := int(targetSizeKB * 1024)
	if targetSize > MaxContractSize {
		targetSize = MaxContractSize
	}

	compiler := geas.NewCompiler(nil)

	if targetSize < 1024 {
		// For small contracts (< 1KB), use simple padding with JUMPDEST opcodes
		// The contract will store ADDRESS at memory[0] for uniqueness, then pad with JUMPDESTs
		paddingSize := targetSize - 33 - 10 // Account for MSTORE(ADDRESS) and RETURN opcodes
		if paddingSize < 0 {
			paddingSize = 0
		}

		// Build padding string
		padding := ""
		for i := 0; i < paddingSize; i++ {
			padding += "jumpdest\n"
		}

		initcodeGeas := fmt.Sprintf(`
			;; Store deployer address for uniqueness
			address
			push 0
			mstore

			;; Padding with JUMPDESTs
			%s

			;; Set first byte to STOP for efficient CALL handling
			push 0x00
			push 0
			mstore8

			;; Return the contract bytecode
			push %d
			push 0
			return
		`, padding, targetSize)

		code := compiler.CompileString(initcodeGeas)
		if code == nil {
			return nil, fmt.Errorf("failed to compile initcode for small contract")
		}
		return code, nil
	}

	// For larger contracts, use keccak256 expansion pattern
	// We generate an XOR table and use SHA3 to expand the bytecode
	xorTableSize := min(256, targetSize/256)
	if xorTableSize < 1 {
		xorTableSize = 1
	}

	// Generate XOR table values
	xorValues := make([]string, xorTableSize)
	for i := 0; i < xorTableSize; i++ {
		// Use keccak256(i) as XOR value
		hash := crypto.Keccak256(big.NewInt(int64(i)).Bytes())
		xorValues[i] = "0x" + hex.EncodeToString(hash)
	}

	// Build XOR expansion loop body
	xorExpansion := ""
	for _, xorVal := range xorValues {
		xorExpansion += fmt.Sprintf(`
			push %s
			xor
			dup1
			msize
			mstore
		`, xorVal)
	}

	initcodeGeas := fmt.Sprintf(`
		;; Store ADDRESS as initial seed - creates uniqueness per deployment
		address
		push 0
		mstore

	loop:
		;; Check if we've reached target size
		msize
		push %d
		lt
		iszero
		jumpi @done

		;; keccak256 of last 32 bytes to expand
		push 32
		msize
		push 32
		sub
		keccak256

		;; XOR expansion
		%s
		pop

		jump @loop

	done:
		;; Set first byte to STOP for efficient CALL handling
		push 0x00
		push 0
		mstore8

		;; Return the full contract
		push %d
		push 0
		return
	`, targetSize, xorExpansion, targetSize)

	code := compiler.CompileString(initcodeGeas)
	if code == nil {
		return nil, fmt.Errorf("failed to compile initcode for large contract")
	}
	return code, nil
}

// buildInitcodeDeployer builds deployment bytecode that returns the initcode as contract code
func buildInitcodeDeployer(initcode []byte) []byte {
	compiler := geas.NewCompiler(nil)

	// Prefix size is 12 bytes (hardcoded based on geas output)
	prefixSize := 12

	deployerGeas := fmt.Sprintf(`
		;; CODECOPY(destOffset=0, offset=prefix_size, size=initcode_size)
		push %d          ;; initcode size
		dup1             ;; duplicate for RETURN
		push %d          ;; offset (after this prefix)
		push 0           ;; dest offset in memory
		codecopy

		;; RETURN(offset=0, size=initcode_size)
		push 0
		return
	`, len(initcode), prefixSize)

	deployerCode := compiler.CompileString(deployerGeas)
	if deployerCode == nil {
		// Fallback to manual bytecode construction
		return buildInitcodeDeployerManual(initcode)
	}

	// Adjust offset if actual deployer size differs
	actualPrefixSize := len(deployerCode)
	if actualPrefixSize != prefixSize {
		deployerGeas = fmt.Sprintf(`
			push %d
			dup1
			push %d
			push 0
			codecopy
			push 0
			return
		`, len(initcode), actualPrefixSize)
		deployerCode = compiler.CompileString(deployerGeas)
	}

	return append(deployerCode, initcode...)
}

// buildInitcodeDeployerManual is a fallback that builds the deployer bytecode manually
func buildInitcodeDeployerManual(initcode []byte) []byte {
	size := len(initcode)

	// PUSH2 size, DUP1, PUSH1 offset, PUSH1 0, CODECOPY, PUSH1 0, RETURN
	// Total: 12 bytes
	deployer := []byte{
		0x61, byte(size >> 8), byte(size), // PUSH2 size
		0x80,       // DUP1
		0x60, 0x0c, // PUSH1 12 (offset after this prefix)
		0x60, 0x00, // PUSH1 0 (dest)
		0x39,       // CODECOPY
		0x60, 0x00, // PUSH1 0
		0xf3, // RETURN
	}

	return append(deployer, initcode...)
}

// buildFactory builds a CREATE2 factory contract bytecode
func buildFactory(initcodeAddr common.Address, initcodeHash common.Hash, initcodeSize int) ([]byte, error) {
	compiler := geas.NewCompiler(nil)

	// Factory layout:
	// - Storage slot 0: Counter (number of deployed contracts)
	// - Storage slot 1: Init code hash for CREATE2 address calculation
	// - Storage slot 2: Init code address
	//
	// Interface:
	// - When called with CALLDATASIZE == 0: Returns (num_deployed_contracts, init_code_hash)
	// - When called otherwise: Deploys a new contract via CREATE2

	// Constructor: store init code hash and address
	constructorGeas := fmt.Sprintf(`
		;; Store init code hash in slot 1
		push %s
		push 1
		sstore

		;; Store initcode address in slot 2
		push %s
		push 2
		sstore
	`, initcodeHash.Hex(), initcodeAddr.Hex())

	// Runtime code
	runtimeGeas := fmt.Sprintf(`
	runtime_start:
		;; Check if this is a getConfig() call (CALLDATASIZE == 0)
		calldatasize
		iszero
		jumpi @getconfig

		;; === CREATE2 DEPLOYMENT PATH ===
		;; Load initcode address from storage slot 2
		push 2
		sload

		;; EXTCODECOPY: copy initcode to memory
		push %d          ;; size
		push 0           ;; source offset
		push 0           ;; dest offset
		dup4             ;; address
		extcodecopy

		;; Prepare for CREATE2
		push %d          ;; size
		swap1
		pop              ;; remove address from stack

		;; CREATE2 with current counter as salt
		push 0           ;; slot 0
		sload            ;; load counter (use as salt)
		swap1            ;; put size on top
		push 0           ;; offset in memory
		push 0           ;; value
		create2          ;; create contract

		;; Store the created address for return
		dup1
		push 0
		mstore

		;; Increment counter
		push 0           ;; slot 0
		dup1             ;; duplicate
		sload            ;; load counter
		push 1           ;; increment
		add              ;; add
		swap1            ;; swap
		sstore           ;; store new counter

		;; Return the created address
		push 32          ;; return 32 bytes
		push 0           ;; from memory position 0
		return

	getconfig:
		;; === GETCONFIG PATH ===
		push 0           ;; slot 0
		sload            ;; load number of deployed contracts
		push 0           ;; memory position 0
		mstore           ;; store in memory

		push 1           ;; slot 1
		sload            ;; load init code hash
		push 32          ;; memory position 32
		mstore           ;; store in memory

		push 64          ;; return 64 bytes (2 * 32)
		push 0           ;; from memory position 0
		return
	`, initcodeSize, initcodeSize)

	// Compile constructor
	constructorCode := compiler.CompileString(constructorGeas)
	if constructorCode == nil {
		return nil, fmt.Errorf("failed to compile factory constructor")
	}

	// Compile runtime
	runtimeCode := compiler.CompileString(runtimeGeas)
	if runtimeCode == nil {
		return nil, fmt.Errorf("failed to compile factory runtime")
	}

	// Build deployer that executes constructor and returns runtime
	runtimeSize := len(runtimeCode)
	constructorSize := len(constructorCode)

	// Deployer code: execute constructor, then copy and return runtime
	deployerGeas := fmt.Sprintf(`
		;; Copy runtime code to memory
		push %d          ;; runtime size
		push %d          ;; offset to runtime
		push 0           ;; dest in memory
		codecopy

		;; Return runtime code
		push %d          ;; size to return
		push 0           ;; offset in memory
		return
	`, runtimeSize, constructorSize+14, runtimeSize) // 14 is approximate deployer size

	deployerCode := compiler.CompileString(deployerGeas)
	if deployerCode == nil {
		return nil, fmt.Errorf("failed to compile factory deployer")
	}

	// Adjust offset based on actual deployer size
	actualDeployerSize := len(deployerCode)
	runtimeOffset := constructorSize + actualDeployerSize

	deployerGeas = fmt.Sprintf(`
		push %d
		push %d
		push 0
		codecopy
		push %d
		push 0
		return
	`, runtimeSize, runtimeOffset, runtimeSize)

	deployerCode = compiler.CompileString(deployerGeas)
	if deployerCode == nil {
		return nil, fmt.Errorf("failed to compile adjusted factory deployer")
	}

	// Final bytecode: constructor + deployer + runtime
	return append(append(constructorCode, deployerCode...), runtimeCode...), nil
}
