package erc20bloater

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/scenarios/statebloat/erc20_bloater/contract"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/utils"
)

const (
	BytesPerSlot         = 32
	SlotsPerBloatCycle   = 2                                                                                // Each iteration: 1 balance + 1 allowance
	DefaultInitialSupply = "115792089237316195423570985008687907853269984665640564039457584007913129639935" // max uint256
)

type ScenarioOptions struct {
	TargetStorageGB  float64 `yaml:"target_storage_gb" json:"target_storage_gb"`
	TargetGasRatio   float64 `yaml:"target_gas_ratio" json:"target_gas_ratio"`
	BaseFee          float64 `yaml:"base_fee" json:"base_fee"`
	TipFee           float64 `yaml:"tip_fee" json:"tip_fee"`
	ExistingContract string  `yaml:"existing_contract" json:"existing_contract"` // Optional override for edge cases
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	contractAddr     common.Address
	contractInstance *contract.ERC20Bloater
}

var ScenarioName = "erc20_bloater"
var ScenarioDefaultOptions = ScenarioOptions{
	TargetStorageGB:  1.0,
	TargetGasRatio:   0.50,
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
	flags.Float64Var(&s.options.TargetGasRatio, "target-gas-ratio", ScenarioDefaultOptions.TargetGasRatio, "Target gas usage as ratio of block gas limit (default 0.50 = 50%)")
	flags.Float64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Base fee per gas in gwei")
	flags.Float64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Tip fee per gas in gwei")
	flags.StringVar(&s.options.ExistingContract, "existing-contract", ScenarioDefaultOptions.ExistingContract, "(Optional) Override contract address for edge cases")
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
	// Include runtime state in config output for web UI visibility
	type ConfigWithState struct {
		ScenarioOptions
		ContractAddress string `yaml:"contract_address,omitempty" json:"contract_address,omitempty"`
	}

	cfg := ConfigWithState{
		ScenarioOptions: s.options,
	}

	// Add contract address if known
	if s.contractAddr != (common.Address{}) {
		cfg.ContractAddress = s.contractAddr.Hex()
	}

	yamlBytes, _ := yaml.Marshal(&cfg)
	return string(yamlBytes)
}

func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Infof("starting scenario: %s", ScenarioName)
	defer s.logger.Infof("scenario %s finished", ScenarioName)

	client := s.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0),
	)
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, 0)

	var startSlot uint64 = 1 // Default: start from address 0x01

	// Determine contract address using nonce-based approach
	if s.options.ExistingContract != "" {
		// Manual override for edge cases
		s.contractAddr = common.HexToAddress(s.options.ExistingContract)
		s.logger.Infof("using manually specified contract: %s", s.contractAddr.Hex())
	} else {
		// Nonce-based automatic detection
		nonce := wallet.GetNonce()

		if nonce == 0 {
			// Fresh wallet - deploy new contract
			s.logger.Infof("wallet nonce is 0, deploying new contract...")
			receipt, _, err := s.deployContract(ctx)
			if err != nil {
				return fmt.Errorf("failed to deploy contract: %w", err)
			}
			s.contractAddr = receipt.ContractAddress
			s.logger.Infof("deployed contract: %s (block #%d)", s.contractAddr.Hex(), receipt.BlockNumber.Uint64())
			s.logger.Infof("to resume later, use same wallet (--seed) - contract will be auto-detected")
			startSlot = 1
		} else {
			// Wallet has history - contract should exist at nonce 0 address
			s.contractAddr = crypto.CreateAddress(wallet.GetAddress(), 0)
			s.logger.Infof("wallet nonce is %d, calculated contract address: %s", nonce, s.contractAddr.Hex())
		}
	}

	// Bind to contract
	contractInstance, err := contract.NewERC20Bloater(s.contractAddr, client.GetEthClient())
	if err != nil {
		return fmt.Errorf("failed to bind to contract: %w", err)
	}
	s.contractInstance = contractInstance

	// Query on-chain progress from contract (if resuming)
	if wallet.GetNonce() > 0 || s.options.ExistingContract != "" {
		nextSlot, err := contractInstance.NextStorageSlot(nil)
		if err != nil {
			return fmt.Errorf("failed to query nextStorageSlot from contract: %w", err)
		}

		startSlot = nextSlot.Uint64()
		if startSlot == 0 {
			startSlot = 1 // Contract not yet bloated, start from slot 1
		}

		// Calculate and log current progress
		targetBytes := uint64(s.options.TargetStorageGB * 1024 * 1024 * 1024)
		targetSlots := targetBytes / BytesPerSlot
		currentGB := float64(startSlot*BytesPerSlot) / (1024 * 1024 * 1024)
		progress := float64(startSlot) / float64(targetSlots) * 100

		s.logger.Infof("resuming from on-chain state: slot %d (%.2f%% complete, %.3f GB / %.3f GB)",
			startSlot, progress, currentGB, s.options.TargetStorageGB)
	}

	// Query network gas limit
	blockGasLimit, err := s.walletPool.GetTxPool().GetCurrentGasLimitWithInit()
	if err != nil {
		return fmt.Errorf("failed to get current gas limit: %w", err)
	}

	// Calculate target slots needed
	targetBytes := uint64(s.options.TargetStorageGB * 1024 * 1024 * 1024)
	targetSlots := targetBytes / BytesPerSlot
	s.logger.Infof("target: %.2f GB = %d slots (%.2f million addresses)",
		s.options.TargetStorageGB, targetSlots, float64(targetSlots)/float64(SlotsPerBloatCycle)/1000000)

	// Start bloating with EIP-7825 compliant transaction splitting
	totalTxCount := uint64(0)
	errorCount := 0
	estimatedGasPerAddress := uint64(50000)

	for startSlot < targetSlots {
		select {
		case <-ctx.Done():
			s.logger.Info("context cancelled, exiting")
			return ctx.Err()
		default:
		}

		// Calculate how to split transactions for EIP-7825 compliance
		totalTargetGas := uint64(float64(blockGasLimit) * s.options.TargetGasRatio)
		txSplits := s.calculateTransactionSplits(totalTargetGas)

		// Log splitting strategy if needed
		if len(txSplits) > 1 {
			s.logger.Infof("EIP-7825: splitting target gas (%.1fM) across %d transactions (%.1fM gas each)",
				float64(totalTargetGas)/1000000, len(txSplits), float64(txSplits[0])/1000000)
		} else {
			s.logger.Infof("block gas limit: %d, target gas: %d (%.0f%%) - single tx",
				blockGasLimit, totalTargetGas, s.options.TargetGasRatio*100)
		}

		// Process batch of transactions for this round
		roundStartSlot := startSlot
		roundSuccess := true
		batchTxCount := 0

		for i, gasLimit := range txSplits {
			// Calculate addresses for this transaction
			addressesForTx := gasLimit / estimatedGasPerAddress
			endSlot := startSlot + addressesForTx*SlotsPerBloatCycle
			if endSlot > targetSlots {
				endSlot = targetSlots
			}

			numAddresses := (endSlot - startSlot) / SlotsPerBloatCycle
			if numAddresses == 0 {
				break // No more addresses to process
			}

			// Submit bloating transaction with explicit gas limit
			tx, w, receipt, err := s.sendBloatTxWithGasLimit(ctx, startSlot/SlotsPerBloatCycle, numAddresses, gasLimit)
			if err != nil {
				s.logger.Errorf("batch tx %d/%d failed: %v", i+1, len(txSplits), err)
				roundSuccess = false
				break
			}

			totalTxCount++
			batchTxCount++

			s.logger.WithFields(logrus.Fields{
				"batch":     fmt.Sprintf("%d/%d", i+1, len(txSplits)),
				"tx":        tx.Hash().Hex(),
				"nonce":     tx.Nonce(),
				"wallet":    s.walletPool.GetWalletName(w.GetAddress()),
				"addresses": numAddresses,
				"gas_used":  receipt.GasUsed,
				"gas_limit": gasLimit,
			}).Infof("sent bloat tx #%d", totalTxCount)

			if receipt.Status != types.ReceiptStatusSuccessful {
				s.logger.Errorf("tx failed: %s (gas used: %d, gas limit: %d)",
					tx.Hash().Hex(), receipt.GasUsed, tx.Gas())
				roundSuccess = false
				break
			}

			startSlot = endSlot

			// Break if we've reached target
			if startSlot >= targetSlots {
				break
			}
		}

		if !roundSuccess {
			// Revert to beginning of round on failure
			startSlot = roundStartSlot
			errorCount++
			time.Sleep(time.Second * time.Duration(errorCount))
			continue
		}

		// Reset error count on successful round
		errorCount = 0

		// Log progress after successful round
		currentGB := float64(startSlot*BytesPerSlot) / (1024 * 1024 * 1024)
		progress := float64(startSlot) / float64(targetSlots) * 100
		s.logger.Infof("progress: %.2f%% | slots: %d / %d | storage: %.3f GB / %.3f GB | round txs: %d",
			progress, startSlot, targetSlots, currentGB, s.options.TargetStorageGB, batchTxCount)
	}

	// Log completion
	finalGB := float64(startSlot*BytesPerSlot) / (1024 * 1024 * 1024)
	s.logger.Infof("bloating complete! total slots: %d, estimated storage: %.3f GB, total txs: %d",
		startSlot, finalGB, totalTxCount)

	return nil
}

// calculateTransactionSplits determines how to split the total target gas across multiple transactions
// to comply with EIP-7825's maximum gas limit per transaction
func (s *Scenario) calculateTransactionSplits(totalTargetGas uint64) []uint64 {
	// Check if we need to split at all
	if totalTargetGas <= utils.MaxGasLimitPerTx {
		return []uint64{totalTargetGas}
	}

	// Calculate number of transactions needed
	numTxs := (totalTargetGas + utils.MaxGasLimitPerTx - 1) / utils.MaxGasLimitPerTx // ceiling division
	gasPerTx := totalTargetGas / numTxs

	// Ensure each transaction doesn't exceed the limit
	if gasPerTx > utils.MaxGasLimitPerTx {
		gasPerTx = utils.MaxGasLimitPerTx
	}

	// Build the splits array
	splits := make([]uint64, numTxs)
	remainingGas := totalTargetGas

	for i := 0; i < int(numTxs); i++ {
		if remainingGas >= gasPerTx {
			splits[i] = gasPerTx
			remainingGas -= gasPerTx
		} else if remainingGas > 0 {
			splits[i] = remainingGas
			remainingGas = 0
		}
	}

	// Adjust last transaction to account for any rounding
	if remainingGas > 0 && len(splits) > 0 {
		lastIdx := len(splits) - 1
		if splits[lastIdx]+remainingGas <= utils.MaxGasLimitPerTx {
			splits[lastIdx] += remainingGas
		}
	}

	return splits
}

func (s *Scenario) deployContract(ctx context.Context) (*types.Receipt, *types.Transaction, error) {
	client := s.walletPool.GetClient(spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0))
	if client == nil {
		return nil, nil, fmt.Errorf("no client available")
	}

	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, 0)
	if wallet == nil {
		return nil, nil, fmt.Errorf("no wallet available")
	}

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

// sendBloatTxWithGasLimit sends a bloating transaction with an explicit gas limit for EIP-7825 compliance
func (s *Scenario) sendBloatTxWithGasLimit(ctx context.Context, startSlot uint64, numAddresses uint64, targetGasLimit uint64) (*types.Transaction, *spamoor.Wallet, *types.Receipt, error) {
	client := s.walletPool.GetClient(spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0))
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, 0)

	if client == nil {
		return nil, nil, nil, fmt.Errorf("no client available")
	}

	feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(client, s.options.BaseFee, s.options.TipFee)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get suggested fees: %w", err)
	}

	var gasLimit uint64

	// If target gas limit is provided and reasonable, use it
	if targetGasLimit > 0 {
		gasLimit = targetGasLimit
		// Ensure we don't exceed EIP-7825 limit
		if gasLimit > utils.MaxGasLimitPerTx {
			s.logger.Warnf("requested gas limit %d exceeds EIP-7825 max of %d, capping", gasLimit, utils.MaxGasLimitPerTx)
			gasLimit = utils.MaxGasLimitPerTx
		}
	} else {
		// Fallback to estimation
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

		// Apply EIP-7825 cap to estimated gas
		if gasLimit > utils.MaxGasLimitPerTx {
			s.logger.Warnf("estimated gas %d exceeds EIP-7825 limit, capping to %d", gasLimit, utils.MaxGasLimitPerTx)
			gasLimit = utils.MaxGasLimitPerTx
		}
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

	// Send transaction and wait for receipt
	receipt, err := s.walletPool.GetTxPool().SendAndAwaitTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		Rebroadcast: true,
	})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to send/confirm tx: %w", err)
	}

	s.logger.Debugf("tx %s confirmed in block #%d (status: %d, gas used: %d/%d)",
		tx.Hash().Hex(), receipt.BlockNumber.Uint64(), receipt.Status, receipt.GasUsed, gasLimit)

	return tx, wallet, receipt, nil
}
