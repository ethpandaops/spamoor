package sbrandsstore

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
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/scenarios/statebloat/rand_sstore/contract"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/utils"
)

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
	BaseFee         uint64 `yaml:"base_fee"`
	TipFee          uint64 `yaml:"tip_fee"`
	DeploymentsFile string `yaml:"deployments_file"`
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	// Contract state
	contractAddress  common.Address
	contractABI      abi.ABI
	contractInstance *contract.SSTOREStorageBloater // Generated contract binding
	isDeployed       bool
	deployMutex      sync.Mutex

	// Scenario state
	totalSlots  uint64 // Total number of slots created
	roundNumber uint64 // Current round number for SSTORE bloating

	// Adaptive gas tracking
	actualGasPerNewSlotIteration uint64          // Dynamically adjusted based on actual usage
	successfulSlotCounts         map[uint64]bool // Track successful slot counts to avoid retries
}

var ScenarioName = "statebloat-rand-sstore"
var ScenarioDefaultOptions = ScenarioOptions{
	BaseFee:         10, // 10 gwei default
	TipFee:          2,  // 2 gwei default
	DeploymentsFile: "",
}
var ScenarioDescriptor = scenario.Descriptor{
	Name:           ScenarioName,
	Description:    "Maximum state bloat via SSTORE operations using curve25519 prime dispersion",
	DefaultOptions: ScenarioDefaultOptions,
	NewScenario:    newScenario,
}

func newScenario(logger logrus.FieldLogger) scenario.Scenario {
	return &Scenario{
		logger:                       logger.WithField("scenario", ScenarioName),
		actualGasPerNewSlotIteration: GasPerNewSlotIteration, // Start with estimated values
		successfulSlotCounts:         make(map[uint64]bool),
	}
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
	flags.Uint64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Base fee per gas in gwei")
	flags.Uint64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Tip fee per gas in gwei")
	flags.StringVar(&s.options.DeploymentsFile, "deployments-file", ScenarioDefaultOptions.DeploymentsFile, "Deployments file")
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

	s.walletPool.SetWalletCount(1)
	s.walletPool.SetRefillAmount(uint256.NewInt(0).Mul(uint256.NewInt(20), uint256.NewInt(1000000000000000000)))  // 20 ETH
	s.walletPool.SetRefillBalance(uint256.NewInt(0).Mul(uint256.NewInt(10), uint256.NewInt(1000000000000000000))) // 10 ETH

	// register well known wallets
	s.walletPool.AddWellKnownWallet(&spamoor.WellKnownWalletConfig{
		Name:          "deployer",
		RefillAmount:  uint256.NewInt(2000000000000000000), // 2 ETH
		RefillBalance: uint256.NewInt(1000000000000000000), // 1 ETH
	})

	// Parse contract ABI
	parsedABI, err := abi.JSON(strings.NewReader(string(contract.SSTOREStorageBloaterMetaData.ABI)))
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
func (s *Scenario) loadDeploymentFile() (DeploymentFile, error) {
	if s.options.DeploymentsFile == "" {
		return make(DeploymentFile), nil
	}

	data, err := os.ReadFile(s.options.DeploymentsFile)
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
func (s *Scenario) saveDeploymentFile(deployments DeploymentFile) error {
	if s.options.DeploymentsFile == "" {
		return nil
	}

	data, err := json.MarshalIndent(deployments, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal deployment file: %w", err)
	}

	if err := os.WriteFile(s.options.DeploymentsFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write deployment file: %w", err)
	}

	return nil
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

	wallet := s.walletPool.GetWellKnownWallet("deployer")
	if wallet == nil {
		return fmt.Errorf("no wallet available")
	}

	feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(client, s.options.BaseFee, s.options.TipFee)
	if err != nil {
		return err
	}

	tx, err := wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       2000000,
		Value:     uint256.NewInt(0),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		_, deployTx, _, err := contract.DeploySSTOREStorageBloater(transactOpts, client.GetEthClient())
		return deployTx, err
	})

	if err != nil {
		return err
	}

	var txReceipt *types.Receipt
	var txErr error
	txWg := sync.WaitGroup{}
	txWg.Add(1)

	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		Rebroadcast: true,
		OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			defer func() {
				txWg.Done()
			}()

			txErr = err
			txReceipt = receipt
		},
	})
	if err != nil {
		return err
	}

	txWg.Wait()
	if txErr != nil {
		return err
	}

	s.contractAddress = txReceipt.ContractAddress
	s.contractInstance, err = contract.NewSSTOREStorageBloater(s.contractAddress, client.GetEthClient())
	if err != nil {
		return err
	}
	s.isDeployed = true

	// No need to reset nonce - the wallet manager handles it automatically

	// Track deployment in JSON file
	deployments, err := s.loadDeploymentFile()
	if err != nil {
		s.logger.Warnf("failed to load deployment file: %v", err)
		deployments = make(DeploymentFile)
	}

	// Initialize deployment data for this contract
	deployments[s.contractAddress.Hex()] = &DeploymentData{
		StorageRounds: []BlockInfo{},
	}

	if err := s.saveDeploymentFile(deployments); err != nil {
		s.logger.Warnf("failed to save deployment file: %v", err)
	}

	s.logger.WithField("address", s.contractAddress.Hex()).Info("SSTOREStorageBloater contract deployed successfully")

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

		blockGasLimit, err := s.walletPool.GetTxPool().GetCurrentGasLimitWithInit()
		if err != nil {
			s.logger.Warnf("failed to get current gas limit: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

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

	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, 0)
	if wallet == nil {
		return fmt.Errorf("no wallet available")
	}

	// Create transaction options
	feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(client, s.options.BaseFee, s.options.TipFee)
	if err != nil {
		return fmt.Errorf("failed to get suggested fees: %w", err)
	}

	tx, err := wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       targetGas,
		Value:     uint256.NewInt(0),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		return s.contractInstance.CreateSlots(transactOpts, big.NewInt(int64(slotsToCreate)))
	})
	if err != nil {
		return err
	}

	var txReceipt *types.Receipt
	var txErr error
	txWg := sync.WaitGroup{}
	txWg.Add(1)

	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		Rebroadcast: true,
		OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			if receipt != nil {
				txFees := utils.GetTransactionFees(tx, receipt)
				s.logger.WithField("rpc", client.GetName()).Debugf(" transaction confirmed in block #%v. total fee: %v gwei (base: %v) logs: %v", receipt.BlockNumber.String(), txFees.TotalFeeGwei(), txFees.TxBaseFeeGwei(), len(receipt.Logs))
			}
			txErr = err
			txReceipt = receipt
			txWg.Done()
		},
		LogFn: func(client *spamoor.Client, retry int, rebroadcast int, err error) {
			logger := s.logger.WithField("rpc", client.GetName()).WithField("nonce", tx.Nonce())
			if retry == 0 && rebroadcast > 0 {
				logger.Infof("rebroadcasting tx")
			}
			if retry > 0 {
				logger = logger.WithField("retry", retry)
			}
			if rebroadcast > 0 {
				logger = logger.WithField("rebroadcast", rebroadcast)
			}
			if err != nil {
				logger.Debugf("failed sending tx: %v", err)
			} else if retry > 0 || rebroadcast > 0 {
				logger.Debugf("successfully sent tx")
			}
		},
	})
	if err != nil {
		// reset nonce if tx was not sent
		wallet.ResetPendingNonce(ctx, client)

		return err
	}

	txWg.Wait()
	if txErr != nil {
		return fmt.Errorf("transaction failed: %w", txErr)
	}

	if txReceipt == nil || txReceipt.Status != 1 {
		// Increase our gas estimate by 10%
		s.actualGasPerNewSlotIteration = uint64(float64(s.actualGasPerNewSlotIteration) * 1.1)

		return fmt.Errorf("transaction rejected")
	}

	// Mark this slot count as successful
	s.successfulSlotCounts[slotsToCreate] = true

	// Update metrics and adaptive gas tracking
	s.totalSlots += slotsToCreate
	totalOverhead := BaseTxCost + FunctionCallOverhead
	actualGasPerSlotIteration := (txReceipt.GasUsed - totalOverhead) / slotsToCreate

	// Update our gas estimate using exponential moving average
	// New estimate = 0.7 * old estimate + 0.3 * actual
	s.actualGasPerNewSlotIteration = uint64(float64(s.actualGasPerNewSlotIteration)*0.7 + float64(actualGasPerSlotIteration)*0.3)

	// Get previous block info for tracking
	prevBlockNumber := txReceipt.BlockNumber.Uint64() - 1

	prevBlock, err := client.GetEthClient().BlockByNumber(ctx, big.NewInt(int64(prevBlockNumber)))
	if err != nil {
		s.logger.Warnf("failed to get previous block info: %v", err)
	} else {
		// Track this storage round in deployment file
		deployments, err := s.loadDeploymentFile()
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

				if err := s.saveDeploymentFile(deployments); err != nil {
					s.logger.Warnf("failed to save deployment file: %v", err)
				}
			}
		}
	}

	// Calculate MB written in this transaction (64 bytes per slot: 32 byte key + 32 byte value)
	mbWrittenThisTx := float64(slotsToCreate*64) / (1024 * 1024)

	// Calculate block utilization percentage
	blockUtilization := float64(txReceipt.GasUsed) / float64(blockGasLimit) * 100

	s.logger.WithFields(logrus.Fields{
		"block_number":      txReceipt.BlockNumber,
		"gas_used":          txReceipt.GasUsed,
		"slots_created":     slotsToCreate,
		"gas_per_slot":      actualGasPerSlotIteration,
		"total_slots":       s.totalSlots,
		"mb_written":        mbWrittenThisTx,
		"block_utilization": fmt.Sprintf("%.2f%%", blockUtilization),
	}).Info("SSTORE bloating round summary")

	return nil
}
