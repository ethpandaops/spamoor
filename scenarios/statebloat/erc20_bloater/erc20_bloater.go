package erc20bloater

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/scenarios/statebloat/erc20_bloater/contract"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
)

const (
	CheckpointFileName  = ".erc20_bloater_checkpoint.json"
	ConfigFileName      = "config.yaml"
	BytesPerSlot        = 32
	SlotsPerBloatCycle  = 2 // Each iteration: 1 balance + 1 allowance
	DefaultInitialSupply = "115792089237316195423570985008687907853269984665640564039457584007913129639935" // max uint256
)

type ScenarioOptions struct {
	TargetStorageGB  float64 `yaml:"target_storage_gb" json:"target_storage_gb"`
	TargetGasRatio   float64 `yaml:"target_gas_ratio" json:"target_gas_ratio"`
	BaseFee          uint64  `yaml:"base_fee" json:"base_fee"`
	TipFee           uint64  `yaml:"tip_fee" json:"tip_fee"`
	ExistingContract string  `yaml:"existing_contract" json:"existing_contract"`
}

type Checkpoint struct {
	ContractAddress      common.Address `json:"contract_address"`
	LastSuccessfulSlot   uint64         `json:"last_successful_slot"`
	NextSlotToWrite      uint64         `json:"next_slot_to_write"`
	TotalSlotsCreated    uint64         `json:"total_slots_created"`
	EstimatedStorageGB   float64        `json:"estimated_storage_gb"`
	LastBlockNumber      uint64         `json:"last_block_number"`
	LastTxHash           string         `json:"last_transaction_hash"`
	LastUpdateTimestamp  time.Time      `json:"last_update_timestamp"`
	ErrorCount           int            `json:"error_count"`
	LastError            *string        `json:"last_error"`
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	contractAddr     common.Address
	contractInstance *contract.ERC20Bloater

	checkpoint      *Checkpoint
	checkpointMutex sync.Mutex

	chainID      *big.Int
	chainIDOnce  sync.Once
	chainIDError error
}

var ScenarioName = "erc20_bloater"
var ScenarioDefaultOptions = ScenarioOptions{
	TargetStorageGB:  1.0,
	TargetGasRatio:   0.90,
	BaseFee:          20,
	TipFee:           2,
	ExistingContract: "",
}
var ScenarioDescriptor = scenario.Descriptor{
	Name:           ScenarioName,
	Description:    "Bloat ERC20 contract storage to target GB size using sequential addresses",
	DefaultOptions: ScenarioDefaultOptions,
	NewScenario:    newScenario,
}

func newScenario(logger logrus.FieldLogger) scenario.Scenario {
	return &Scenario{
		logger: logger.WithField("scenario", ScenarioName),
	}
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
	flags.Float64Var(&s.options.TargetStorageGB, "target-gb", ScenarioDefaultOptions.TargetStorageGB, "Target storage size in GB")
	flags.Float64Var(&s.options.TargetGasRatio, "target-gas-ratio", ScenarioDefaultOptions.TargetGasRatio, "Target gas usage as ratio of block gas limit (default 0.90 = 90%)")
	flags.Uint64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Base fee per gas in gwei")
	flags.Uint64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Tip fee per gas in gwei")
	flags.StringVar(&s.options.ExistingContract, "existing-contract", ScenarioDefaultOptions.ExistingContract, "Use existing contract address instead of deploying new one")
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

	// Only use 1 wallet for this scenario
	s.walletPool.SetWalletCount(1)

	return nil
}

func (s *Scenario) Config() string {
	yamlBytes, _ := yaml.Marshal(&s.options)
	return string(yamlBytes)
}

func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Infof("starting scenario: %s", ScenarioName)
	defer s.logger.Infof("scenario %s finished", ScenarioName)

	// Try to load checkpoint (always enabled)
	var startSlot uint64
	checkpoint, err := s.loadCheckpoint()
	if err != nil {
		s.logger.Warnf("failed to load checkpoint: %v, starting fresh", err)
	} else if checkpoint != nil {
		s.checkpoint = checkpoint
		s.contractAddr = checkpoint.ContractAddress
		startSlot = checkpoint.NextSlotToWrite
		s.logger.Infof("resuming from checkpoint: contract=%s, next_slot=%d, total_slots=%d",
			checkpoint.ContractAddress.Hex(), checkpoint.NextSlotToWrite, checkpoint.TotalSlotsCreated)
	}

	// Deploy or use existing contract
	if s.contractAddr == (common.Address{}) {
		if s.options.ExistingContract != "" {
			s.contractAddr = common.HexToAddress(s.options.ExistingContract)
			s.logger.Infof("using existing contract: %s", s.contractAddr.Hex())
		} else {
			receipt, _, err := s.deployContract(ctx)
			if err != nil {
				return fmt.Errorf("failed to deploy contract: %w", err)
			}
			s.contractAddr = receipt.ContractAddress
			s.logger.Infof("deployed contract: %s (block #%d)", s.contractAddr.Hex(), receipt.BlockNumber.Uint64())

			// Save config with deployed contract address for future runs
			s.options.ExistingContract = s.contractAddr.Hex()
			if err := s.saveConfig(); err != nil {
				s.logger.Warnf("failed to save config file: %v", err)
			} else {
				s.logger.Infof("saved config to %s for future runs", ConfigFileName)
			}
		}

		// Initialize checkpoint if not resuming
		if s.checkpoint == nil {
			s.checkpoint = &Checkpoint{
				ContractAddress:     s.contractAddr,
				LastSuccessfulSlot:  0,
				NextSlotToWrite:     1, // Start from address 0x0000...0001
				TotalSlotsCreated:   0,
				EstimatedStorageGB:  0,
				LastUpdateTimestamp: time.Now(),
			}
		}
	}

	// Bind to contract
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, "")
	contractInstance, err := contract.NewERC20Bloater(s.contractAddr, client.GetEthClient())
	if err != nil {
		return fmt.Errorf("failed to bind contract: %w", err)
	}
	s.contractInstance = contractInstance

	// Query network gas limit
	blockGasLimit, err := s.getBlockGasLimit(ctx)
	if err != nil {
		return fmt.Errorf("failed to get block gas limit: %w", err)
	}

	// Calculate target gas per transaction (3/4 of block gas limit by default)
	targetGasPerTx := uint64(float64(blockGasLimit) * s.options.TargetGasRatio)
	s.logger.Infof("block gas limit: %d, target gas per tx: %d (%.0f%%)",
		blockGasLimit, targetGasPerTx, s.options.TargetGasRatio*100)

	// Estimate addresses per transaction
	// Each address needs: transfer (2 SSTOREs) + approve (1 SSTORE) ≈ 45,000-50,000 gas
	estimatedGasPerAddress := uint64(50000)
	addressesPerTx := targetGasPerTx / estimatedGasPerAddress
	s.logger.Infof("estimated addresses per tx: %d", addressesPerTx)

	// Calculate target slots needed
	targetBytes := uint64(s.options.TargetStorageGB * 1024 * 1024 * 1024)
	targetSlots := targetBytes / BytesPerSlot
	s.logger.Infof("target: %.2f GB = %d slots (%.2f million addresses)",
		s.options.TargetStorageGB, targetSlots, float64(targetSlots)/float64(SlotsPerBloatCycle)/1000000)

	// Start bloating
	txCount := uint64(0)
	for startSlot < targetSlots {
		select {
		case <-ctx.Done():
			s.logger.Info("context cancelled, saving final checkpoint")
			s.saveCheckpoint()
			return ctx.Err()
		default:
		}

		endSlot := startSlot + addressesPerTx*SlotsPerBloatCycle
		if endSlot > targetSlots {
			endSlot = targetSlots
		}
		numAddresses := (endSlot - startSlot) / SlotsPerBloatCycle

		// Submit bloating transaction
		tx, wallet, receipt, err := s.sendBloatTx(ctx, startSlot/SlotsPerBloatCycle, numAddresses)
		if err != nil {
			s.handleError(err)
			time.Sleep(time.Second * time.Duration(s.checkpoint.ErrorCount))
			continue
		}

		s.logger.WithFields(logrus.Fields{
			"tx":        tx.Hash().Hex(),
			"nonce":     tx.Nonce(),
			"wallet":    s.walletPool.GetWalletName(wallet.GetAddress()),
			"addresses": numAddresses,
			"from_slot": startSlot / SlotsPerBloatCycle,
		}).Infof("sent bloat tx #%d", txCount+1)

		if receipt.Status != types.ReceiptStatusSuccessful {
			s.handleError(fmt.Errorf("tx failed: %s (gas used: %d, gas limit: %d)",
				tx.Hash().Hex(), receipt.GasUsed, tx.Gas()))
			time.Sleep(time.Second * time.Duration(s.checkpoint.ErrorCount))
			continue
		}

		// Update checkpoint on success
		s.checkpointMutex.Lock()
		s.checkpoint.LastSuccessfulSlot = endSlot - 1
		s.checkpoint.NextSlotToWrite = endSlot
		s.checkpoint.TotalSlotsCreated = endSlot
		s.checkpoint.EstimatedStorageGB = float64(endSlot*BytesPerSlot) / (1024 * 1024 * 1024)
		s.checkpoint.LastBlockNumber = receipt.BlockNumber.Uint64()
		s.checkpoint.LastTxHash = tx.Hash().Hex()
		s.checkpoint.ErrorCount = 0
		s.checkpoint.LastError = nil
		s.checkpointMutex.Unlock()

		// Save checkpoint after each confirmed transaction
		if err := s.saveCheckpoint(); err != nil {
			s.logger.Warnf("failed to save checkpoint: %v", err)
		}

		txCount++

		// Log progress
		progress := float64(endSlot) / float64(targetSlots) * 100
		s.logger.Infof("progress: %.2f%% | slots: %d / %d | storage: %.3f GB / %.3f GB",
			progress, endSlot, targetSlots, s.checkpoint.EstimatedStorageGB, s.options.TargetStorageGB)

		startSlot = endSlot
	}

	// Save final checkpoint
	s.saveCheckpoint()
	s.logger.Infof("bloating complete! total slots: %d, estimated storage: %.3f GB",
		s.checkpoint.TotalSlotsCreated, s.checkpoint.EstimatedStorageGB)

	return nil
}

func (s *Scenario) deployContract(ctx context.Context) (*types.Receipt, *types.Transaction, error) {
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, "")
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, 0)

	initialSupply, ok := new(big.Int).SetString(DefaultInitialSupply, 10)
	if !ok {
		return nil, nil, fmt.Errorf("failed to parse initial supply")
	}

	feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(client, s.options.BaseFee, s.options.TipFee)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get suggested fees: %w", err)
	}

	var deployedAddr common.Address
	tx, err := wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       2000000,
		Value:     uint256.NewInt(0),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		addr, deployTx, _, err := contract.DeployERC20Bloater(transactOpts, client.GetEthClient(), initialSupply)
		if err != nil {
			return nil, err
		}
		deployedAddr = addr
		return deployTx, nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build deployment tx: %w", err)
	}

	s.logger.Infof("deployment tx sent: %s, waiting for confirmation...", tx.Hash().Hex())

	receipt, err := s.walletPool.GetTxPool().SendAndAwaitTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		Rebroadcast: true,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to send/confirm deployment: %w", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, nil, fmt.Errorf("deployment tx failed")
	}

	s.contractAddr = deployedAddr

	return receipt, tx, nil
}

func (s *Scenario) sendBloatTx(ctx context.Context, startSlot uint64, numAddresses uint64) (*types.Transaction, *spamoor.Wallet, *types.Receipt, error) {
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, "")
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, 0)

	if client == nil {
		return nil, nil, nil, fmt.Errorf("no client available")
	}

	feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(client, s.options.BaseFee, s.options.TipFee)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get suggested fees: %w", err)
	}

	// Estimate gas using RPC call for accuracy
	var gasLimit uint64

	// Pack the contract call data for gas estimation
	abi, err := contract.ERC20BloaterMetaData.GetAbi()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get contract ABI: %w", err)
	}

	callData, err := abi.Pack("bloatStorage", new(big.Int).SetUint64(startSlot), new(big.Int).SetUint64(numAddresses))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to pack call data: %w", err)
	}

	callMsg := ethereum.CallMsg{
		From: wallet.GetAddress(),
		To:   &s.contractAddr,
		Data: callData,
	}

	gasEstimate, err := client.GetEthClient().EstimateGas(ctx, callMsg)
	if err == nil {
		// Add 5% buffer to estimated gas
		gasLimit = uint64(float64(gasEstimate) * 1.05)
		s.logger.Debugf("estimated gas: %d, using with buffer: %d", gasEstimate, gasLimit)
	} else {
		// Fallback to formula-based calculation if estimation fails
		s.logger.Debugf("gas estimation failed: %v, using fallback calculation", err)
		baseGas := uint64(21000)
		gasPerAddress := uint64(55000)
		calculatedGas := baseGas + (numAddresses * gasPerAddress)
		gasLimit = calculatedGas + (calculatedGas / 10) // 10% buffer for fallback
	}

	// Build transaction using BuildBoundTx
	tx, err := wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
		To:        &s.contractAddr,
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       gasLimit,
		Value:     uint256.NewInt(0),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		return s.contractInstance.BloatStorage(transactOpts, new(big.Int).SetUint64(startSlot), new(big.Int).SetUint64(numAddresses))
	})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to build bloat tx: %w", err)
	}

	// Send transaction using TxPool but with manual receipt polling
	// Note: We bypass SendAndAwaitTransaction because it can miss blocks due to TxPool race conditions
	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		Rebroadcast: false, // Disable rebroadcast since we're handling confirmation manually
	})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to send tx: %w", err)
	}

	s.logger.Debugf("sent tx %s, manually polling for receipt...", tx.Hash().Hex())

	// Manually poll for receipt (more reliable than relying on TxPool block processing)
	ethClient := client.GetEthClient()
	var receipt *types.Receipt
	maxAttempts := 60 // 60 attempts * 2s = 2 minute timeout
	for i := 0; i < maxAttempts; i++ {
		receipt, err = ethClient.TransactionReceipt(ctx, tx.Hash())
		if err == nil && receipt != nil {
			s.logger.Debugf("retrieved receipt for tx %s in block #%d (status: %d)",
				tx.Hash().Hex(), receipt.BlockNumber.Uint64(), receipt.Status)
			return tx, wallet, receipt, nil
		}

		// Check if context was cancelled
		select {
		case <-ctx.Done():
			return nil, nil, nil, ctx.Err()
		default:
		}

		time.Sleep(2 * time.Second)
	}

	return nil, nil, nil, fmt.Errorf("timeout waiting for tx confirmation after %d seconds: %s", maxAttempts*2, tx.Hash().Hex())
}

func (s *Scenario) loadCheckpoint() (*Checkpoint, error) {
	data, err := os.ReadFile(CheckpointFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var checkpoint Checkpoint
	if err := json.Unmarshal(data, &checkpoint); err != nil {
		return nil, err
	}

	return &checkpoint, nil
}

func (s *Scenario) saveCheckpoint() error {
	s.checkpointMutex.Lock()
	defer s.checkpointMutex.Unlock()

	if s.checkpoint == nil {
		return nil
	}

	s.checkpoint.LastUpdateTimestamp = time.Now()

	data, err := json.MarshalIndent(s.checkpoint, "", "  ")
	if err != nil {
		return err
	}

	tempPath := CheckpointFileName + ".tmp"
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return err
	}

	return os.Rename(tempPath, CheckpointFileName)
}

func (s *Scenario) saveConfig() error {
	data, err := yaml.Marshal(&s.options)
	if err != nil {
		return err
	}

	return os.WriteFile(ConfigFileName, data, 0644)
}

func (s *Scenario) handleError(err error) {
	s.checkpointMutex.Lock()
	defer s.checkpointMutex.Unlock()

	if s.checkpoint != nil {
		s.checkpoint.ErrorCount++
		errStr := err.Error()
		s.checkpoint.LastError = &errStr

		if saveErr := s.saveCheckpoint(); saveErr != nil {
			s.logger.Errorf("failed to save checkpoint after error: %v", saveErr)
		}
	}

	s.logger.Errorf("bloating error: %v", err)
}

func (s *Scenario) getChainID(ctx context.Context) (*big.Int, error) {
	s.chainIDOnce.Do(func() {
		client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, "")
		s.chainID, s.chainIDError = client.GetChainId(ctx)
	})
	return s.chainID, s.chainIDError
}

func (s *Scenario) getBlockGasLimit(ctx context.Context) (uint64, error) {
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, "")
	block, err := client.GetEthClient().BlockByNumber(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get latest block: %w", err)
	}
	return block.GasLimit(), nil
}
