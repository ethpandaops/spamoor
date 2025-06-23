package randsstorebloater

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"

	"github.com/ethpandaops/spamoor/scenarios/statebloat/rand_sstore_bloater/contract"
	"github.com/ethpandaops/spamoor/scenariotypes"
	"github.com/ethpandaops/spamoor/spamoor"
)

//go:embed contract/SSTOREStorageBloater.abi
var contractABIBytes []byte

//go:embed contract/SSTOREStorageBloater.bin
var contractBytecodeHex []byte

// Constants for SSTORE operations
const (
	// Base Ethereum transaction cost
	BaseTxCost = uint64(21000)

	// Function call overhead (measured from actual transactions)
	// Includes: function selector, ABI decoding, contract loading, etc.
	FunctionCallOverhead = uint64(1556)

	// Gas cost per iteration (measured from actual transactions)
	// Includes: SSTORE (0→non-zero), MULMOD, loop overhead, stack operations
	// Measured: 22,165 gas per iteration
	GasPerNewSlotIteration = uint64(22165)

	// Contract deployment and call overhead
	EstimatedDeployGas = uint64(500000) // Deployment gas for our contract

	// Safety margins and multipliers
	GasLimitSafetyMargin = 0.99 // Use 99% of block gas limit (1% margin for gas price variations)
	
	// Deployment tracking file
	DeploymentFileName = "deployments_sstore_bloating.json"
)

// BlockInfo stores block information for each storage round
type BlockInfo struct {
	BlockNumber uint64 `json:"block_number"`
	Timestamp   uint64 `json:"timestamp"`
}

// DeploymentData tracks a single contract deployment and its storage rounds
type DeploymentData struct {
	StorageRounds []BlockInfo `json:"storage_rounds"`
}

// DeploymentFile represents the entire deployment tracking file
type DeploymentFile map[string]*DeploymentData // key is contract address

type ScenarioOptions struct {
	BaseFee uint64 `yaml:"base_fee"`
	TipFee  uint64 `yaml:"tip_fee"`
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	// Contract state
	contractAddress  common.Address
	contractABI      abi.ABI
	contractInstance *contract.Contract // Generated contract binding
	isDeployed       bool
	deployMutex      sync.Mutex

	// Scenario state
	totalSlots  uint64 // Total number of slots created
	cycleCount  uint64 // Number of complete create/update cycles
	roundNumber uint64 // Current round number for SSTORE bloating

	// Adaptive gas tracking
	actualGasPerNewSlotIteration uint64          // Dynamically adjusted based on actual usage
	successfulSlotCounts         map[uint64]bool // Track successful slot counts to avoid retries

	// Cached values
	chainID      *big.Int
	chainIDOnce  sync.Once
	chainIDError error
}

var ScenarioName = "rand_sstore_bloater"
var ScenarioDefaultOptions = ScenarioOptions{
	BaseFee: 10, // 10 gwei default
	TipFee:  2,  // 2 gwei default
}
var ScenarioDescriptor = scenariotypes.ScenarioDescriptor{
	Name:           ScenarioName,
	Description:    "Maximum state bloat via SSTORE operations using curve25519 prime dispersion",
	DefaultOptions: ScenarioDefaultOptions,
	NewScenario:    newScenario,
}

func newScenario(logger logrus.FieldLogger) scenariotypes.Scenario {
	return &Scenario{
		logger:                       logger.WithField("scenario", ScenarioName),
		actualGasPerNewSlotIteration: GasPerNewSlotIteration, // Start with estimated values
		successfulSlotCounts:         make(map[uint64]bool),
	}
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
	flags.Uint64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Base fee per gas in gwei")
	flags.Uint64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Tip fee per gas in gwei")
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

	// Use only root wallet for simplicity
	s.walletPool.SetWalletCount(1)

	// Parse contract ABI
	parsedABI, err := abi.JSON(strings.NewReader(string(contractABIBytes)))
	if err != nil {
		return fmt.Errorf("failed to parse contract ABI: %w", err)
	}
	s.contractABI = parsedABI

	return nil
}

func (s *Scenario) Config() string {
	yamlBytes, _ := yaml.Marshal(&s.options)
	return string(yamlBytes)
}

// loadDeploymentFile loads the deployment tracking file or creates an empty one
func loadDeploymentFile() (DeploymentFile, error) {
	data, err := os.ReadFile(DeploymentFileName)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, return empty map
			return make(DeploymentFile), nil
		}
		return nil, fmt.Errorf("failed to read deployment file: %w", err)
	}

	var deployments DeploymentFile
	if err := json.Unmarshal(data, &deployments); err != nil {
		return nil, fmt.Errorf("failed to unmarshal deployment file: %w", err)
	}

	return deployments, nil
}

// saveDeploymentFile saves the deployment tracking file
func saveDeploymentFile(deployments DeploymentFile) error {
	data, err := json.MarshalIndent(deployments, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal deployment file: %w", err)
	}

	if err := os.WriteFile(DeploymentFileName, data, 0644); err != nil {
		return fmt.Errorf("failed to write deployment file: %w", err)
	}

	return nil
}

func (s *Scenario) getChainID(ctx context.Context) (*big.Int, error) {
	s.chainIDOnce.Do(func() {
		client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, "")
		if client == nil {
			s.chainIDError = fmt.Errorf("no client available for chain ID")
			return
		}
		s.chainID, s.chainIDError = client.GetChainId(ctx)
	})
	return s.chainID, s.chainIDError
}

func (s *Scenario) deployContract(ctx context.Context) error {
	s.deployMutex.Lock()
	defer s.deployMutex.Unlock()

	if s.isDeployed {
		return nil
	}

	s.logger.Info("Deploying SSTOREStorageBloater contract...")

	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, "")
	if client == nil {
		return fmt.Errorf("no client available")
	}

	wallet := s.walletPool.GetRootWallet()
	if wallet == nil {
		return fmt.Errorf("no wallet available")
	}

	chainID, err := s.getChainID(ctx)
	if err != nil {
		return err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(wallet.GetPrivateKey(), chainID)
	if err != nil {
		return fmt.Errorf("failed to create transactor: %w", err)
	}

	// Set gas parameters
	auth.GasLimit = EstimatedDeployGas
	auth.GasFeeCap = new(big.Int).Mul(big.NewInt(int64(s.options.BaseFee)), big.NewInt(1000000000))
	auth.GasTipCap = new(big.Int).Mul(big.NewInt(int64(s.options.TipFee)), big.NewInt(1000000000))

	// Deploy contract using generated bindings
	address, tx, contractInstance, err := contract.DeployContract(auth, client.GetEthClient())
	if err != nil {
		return fmt.Errorf("failed to deploy contract: %w", err)
	}

	s.logger.WithField("tx", tx.Hash().Hex()).Info("Contract deployment transaction sent")

	// Wait for deployment
	receipt, err := bind.WaitMined(ctx, client.GetEthClient(), tx)
	if err != nil {
		return fmt.Errorf("failed to wait for deployment: %w", err)
	}

	if receipt.Status != 1 {
		return fmt.Errorf("contract deployment failed")
	}

	s.contractAddress = address
	s.contractInstance = contractInstance
	s.isDeployed = true

	// No need to reset nonce - the wallet manager handles it automatically

	// Track deployment in JSON file
	deployments, err := loadDeploymentFile()
	if err != nil {
		s.logger.Warnf("failed to load deployment file: %v", err)
		deployments = make(DeploymentFile)
	}

	// Initialize deployment data for this contract
	deployments[address.Hex()] = &DeploymentData{
		StorageRounds: []BlockInfo{},
	}

	if err := saveDeploymentFile(deployments); err != nil {
		s.logger.Warnf("failed to save deployment file: %v", err)
	}

	s.logger.WithField("address", address.Hex()).Info("SSTOREStorageBloater contract deployed successfully")

	return nil
}

func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Infof("starting scenario: %s", ScenarioName)
	defer s.logger.Infof("scenario %s finished.", ScenarioName)

	// Deploy the contract if not already deployed
	if !s.isDeployed {
		if err := s.deployContract(ctx); err != nil {
			return fmt.Errorf("failed to deploy contract: %w", err)
		}
	}

	// Get network parameters
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, "")
	if client == nil {
		return fmt.Errorf("no client available")
	}

	// Main loop - alternate between creating and updating slots
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Get current block gas limit
		latestBlock, err := client.GetEthClient().BlockByNumber(ctx, nil)
		if err != nil {
			s.logger.Warnf("failed to get latest block: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		blockGasLimit := latestBlock.GasLimit()
		targetGas := uint64(float64(blockGasLimit) * GasLimitSafetyMargin)

		// Never stop spamming SSTORE operations.
		s.roundNumber++
		if err := s.executeCreateSlots(ctx, targetGas, blockGasLimit); err != nil {
			s.logger.Errorf("failed to create slots: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
	}
}

func (s *Scenario) executeCreateSlots(ctx context.Context, targetGas uint64, blockGasLimit uint64) error {
	// Calculate how many slots we can create with precise gas costs
	// Account for base tx cost and function overhead
	availableGas := targetGas - BaseTxCost - FunctionCallOverhead
	slotsToCreate := availableGas / s.actualGasPerNewSlotIteration // Integer division rounds down

	if slotsToCreate == 0 {
		return fmt.Errorf("not enough gas to create any slots")
	}

	// Get client and wallet
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, "")
	if client == nil {
		return fmt.Errorf("no client available")
	}

	wallet := s.walletPool.GetRootWallet()
	if wallet == nil {
		return fmt.Errorf("no wallet available")
	}

	// Create transaction options
	chainID, err := s.getChainID(ctx)
	if err != nil {
		return err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(wallet.GetPrivateKey(), chainID)
	if err != nil {
		return fmt.Errorf("failed to create transactor: %w", err)
	}

	// Set gas parameters
	auth.GasLimit = targetGas
	auth.GasFeeCap = new(big.Int).Mul(big.NewInt(int64(s.options.BaseFee)), big.NewInt(1000000000))
	auth.GasTipCap = new(big.Int).Mul(big.NewInt(int64(s.options.TipFee)), big.NewInt(1000000000))

	// Execute transaction using contract bindings
	tx, err := s.contractInstance.CreateSlots(auth, big.NewInt(int64(slotsToCreate)))
	if err != nil {
		// Check if it's an out-of-gas error
		if strings.Contains(err.Error(), "out of gas") || strings.Contains(err.Error(), "OutOfGas") {
			// Increase our gas estimate by 10%
			s.actualGasPerNewSlotIteration = uint64(float64(s.actualGasPerNewSlotIteration) * 1.1)
			s.logger.Warnf("Out of gas error detected. Adjusting gas per slot estimate to %d", s.actualGasPerNewSlotIteration)
		}
		return err
	}

	// Wait for transaction confirmation
	receipt, err := bind.WaitMined(ctx, client.GetEthClient(), tx)
	if err != nil {
		return fmt.Errorf("failed to wait for transaction: %w", err)
	}

	if receipt.Status != 1 {
		return fmt.Errorf("transaction failed")
	}

	// Mark this slot count as successful
	s.successfulSlotCounts[slotsToCreate] = true

	// Update metrics and adaptive gas tracking
	s.totalSlots += slotsToCreate
	totalOverhead := BaseTxCost + FunctionCallOverhead
	actualGasPerSlotIteration := (receipt.GasUsed - totalOverhead) / slotsToCreate

	// Update our gas estimate using exponential moving average
	// New estimate = 0.7 * old estimate + 0.3 * actual
	s.actualGasPerNewSlotIteration = uint64(float64(s.actualGasPerNewSlotIteration)*0.7 + float64(actualGasPerSlotIteration)*0.3)

	// Get previous block info for tracking
	prevBlockNumber := receipt.BlockNumber.Uint64() - 1
	prevBlock, err := client.GetEthClient().BlockByNumber(ctx, big.NewInt(int64(prevBlockNumber)))
	if err != nil {
		s.logger.Warnf("failed to get previous block info: %v", err)
	} else {
		// Track this storage round in deployment file
		deployments, err := loadDeploymentFile()
		if err != nil {
			s.logger.Warnf("failed to load deployment file: %v", err)
		} else if deployments != nil {
			contractAddr := s.contractAddress.Hex()
			if deploymentData, exists := deployments[contractAddr]; exists {
				// Append new block info
				deploymentData.StorageRounds = append(deploymentData.StorageRounds, BlockInfo{
					BlockNumber: prevBlockNumber,
					Timestamp:   prevBlock.Time(),
				})
				
				if err := saveDeploymentFile(deployments); err != nil {
					s.logger.Warnf("failed to save deployment file: %v", err)
				}
			}
		}
	}

	// Calculate MB written in this transaction (64 bytes per slot: 32 byte key + 32 byte value)
	mbWrittenThisTx := float64(slotsToCreate*64) / (1024 * 1024)
	
	// Calculate block utilization percentage
	blockUtilization := float64(receipt.GasUsed) / float64(blockGasLimit) * 100

	s.logger.WithFields(logrus.Fields{
		"block_number":      receipt.BlockNumber,
		"gas_used":          receipt.GasUsed,
		"slots_created":     slotsToCreate,
		"gas_per_slot":      actualGasPerSlotIteration,
		"total_slots":       s.totalSlots,
		"mb_written":        mbWrittenThisTx,
		"block_utilization": fmt.Sprintf("%.2f%%", blockUtilization),
	}).Info("SSTORE bloating round summary")

	return nil
}
