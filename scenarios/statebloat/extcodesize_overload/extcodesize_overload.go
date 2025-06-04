package extcodesizeoverload

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"sync"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenariotypes"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
)

type ScenarioOptions struct {
	BaseFee uint64 `yaml:"base_fee"`
	TipFee  uint64 `yaml:"tip_fee"`
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool
}

var ScenarioName = "extcodesize-overload"
var ScenarioDefaultOptions = ScenarioOptions{
	BaseFee: 20,
	TipFee:  2,
}
var ScenarioDescriptor = scenariotypes.ScenarioDescriptor{
	Name:           ScenarioName,
	Description:    "Perform EXTCODESIZE calls to maximize gas usage",
	DefaultOptions: ScenarioDefaultOptions,
	NewScenario:    newScenario,
}

func newScenario(logger logrus.FieldLogger) scenariotypes.Scenario {
	return &Scenario{
		logger: logger.WithField("scenario", ScenarioName),
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

	walletPool.SetWalletCount(1)
	return nil
}

func (s *Scenario) Config() string {
	yamlBytes, _ := yaml.Marshal(&s.options)
	return string(yamlBytes)
}

// loadContractAddresses loads contract addresses from deployments.json
func (s *Scenario) loadContractAddresses() ([]common.Address, error) {
	file, err := os.Open("deployments.json")
	if err != nil {
		return nil, fmt.Errorf("failed to open deployments.json: %w", err)
	}
	defer file.Close()

	var deployments map[string][]string
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&deployments); err != nil {
		return nil, fmt.Errorf("failed to decode deployments.json: %w", err)
	}

	var allAddresses []common.Address

	// Use the first private key's contracts
	for _, addresses := range deployments {
		for _, addr := range addresses {
			allAddresses = append(allAddresses, common.HexToAddress(addr))
		}
		break // Use only the first private key's contracts
	}

	if len(allAddresses) == 0 {
		return nil, fmt.Errorf("no contract addresses found in deployments.json")
	}

	return allAddresses, nil
}

// generateExtcodesizeBytecode generates bytecode that performs EXTCODESIZE calls
func (s *Scenario) generateExtcodesizeBytecode(addresses []common.Address) []byte {
	var bytecode []byte

	for _, addr := range addresses {
		// PUSH20 address (0x73 followed by 20 bytes of address)
		bytecode = append(bytecode, 0x73)
		bytecode = append(bytecode, addr.Bytes()...)

		// EXTCODESIZE (0x3B)
		bytecode = append(bytecode, 0x3B)

		// POP to remove the result from stack (0x50)
		bytecode = append(bytecode, 0x50)
	}

	// Return empty to end execution gracefully
	// PUSH1 0x00, PUSH1 0x00, RETURN
	bytecode = append(bytecode, 0x60, 0x00, 0x60, 0x00, 0xF3)

	return bytecode
}

func (s *Scenario) Run(ctx context.Context) error {
	// Get client for gas limit and transaction execution
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, "")
	if client == nil {
		return fmt.Errorf("no client available")
	}

	ethClient := client.GetEthClient()

	// Get chain gas limit
	block, err := ethClient.BlockByNumber(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to get latest block: %w", err)
	}

	gasLimit := block.GasLimit()
	s.logger.WithField("block_gas_limit", gasLimit).Info("Retrieved block gas limit")

	// Calculate maximum EXTCODESIZE calls
	// Gas breakdown: base tx cost (21000) + cold storage access per EXTCODESIZE (2600)
	baseGasCost := uint64(21000)
	extcodesizeGasCost := uint64(2600)

	maxExtcodesizeCalls := (gasLimit - baseGasCost) / extcodesizeGasCost
	s.logger.WithField("max_extcodesize_calls", maxExtcodesizeCalls).Info("Calculated maximum EXTCODESIZE calls")

	// Load contract addresses from deployments.json
	contractAddresses, err := s.loadContractAddresses()
	if err != nil {
		return fmt.Errorf("failed to load contract addresses: %w", err)
	}

	s.logger.WithField("total_contracts", len(contractAddresses)).Info("Loaded contract addresses from deployments.json")

	// Determine how many contracts to call EXTCODESIZE on
	numCalls := uint64(len(contractAddresses))
	if numCalls > maxExtcodesizeCalls {
		numCalls = maxExtcodesizeCalls
		contractAddresses = contractAddresses[:numCalls]
	}

	s.logger.WithField("actual_calls", numCalls).Info("Number of EXTCODESIZE calls to perform")

	// Generate bytecode for EXTCODESIZE calls
	bytecode := s.generateExtcodesizeBytecode(contractAddresses)

	s.logger.WithField("bytecode_size", len(bytecode)).Info("Generated bytecode")
	s.logger.WithField("bytecode_hex", hex.EncodeToString(bytecode)).Debug("Bytecode hex")

	// Get wallet for sending transaction
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, 0)
	if wallet == nil {
		return fmt.Errorf("no wallet available")
	}

	// Create transaction with the bytecode as data
	baseFeeWei := new(big.Int).Mul(big.NewInt(int64(s.options.BaseFee)), big.NewInt(1e9))
	tipFeeWei := new(big.Int).Mul(big.NewInt(int64(s.options.TipFee)), big.NewInt(1e9))
	maxFeePerGas := new(big.Int).Add(baseFeeWei, tipFeeWei)

	estimatedGas := baseGasCost + (numCalls * extcodesizeGasCost)

	s.logger.WithField("estimated_gas", estimatedGas).Info("Estimated gas for transaction")

	// Build transaction using txbuilder pattern
	txData, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(maxFeePerGas),
		GasTipCap: uint256.MustFromBig(tipFeeWei),
		Gas:       estimatedGas,
		To:        nil, // Contract creation to execute bytecode
		Value:     uint256.NewInt(0),
		Data:      bytecode,
	})
	if err != nil {
		return fmt.Errorf("failed to create transaction data: %w", err)
	}

	tx, err := wallet.BuildDynamicFeeTx(txData)
	if err != nil {
		return fmt.Errorf("failed to build transaction: %w", err)
	}

	// Use a WaitGroup to wait for transaction confirmation
	var wg sync.WaitGroup
	var receipt *types.Receipt
	var txErr error
	wg.Add(1)

	// Send transaction using the framework's transaction pool
	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &txbuilder.SendTransactionOptions{
		Client:              client,
		MaxRebroadcasts:     10,
		RebroadcastInterval: 30 * time.Second,
		OnConfirm: func(tx *types.Transaction, r *types.Receipt, err error) {
			defer wg.Done()
			receipt = r
			txErr = err
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	s.logger.WithField("tx_hash", tx.Hash().Hex()).Info("Transaction sent")

	// Wait for transaction confirmation
	wg.Wait()
	if txErr != nil {
		return fmt.Errorf("transaction failed: %w", txErr)
	}

	if receipt == nil {
		return fmt.Errorf("no receipt received")
	}

	// Calculate results
	totalBytesRead := numCalls * 24576 // Each contract is 24KiB
	gasUsed := receipt.GasUsed
	utilizationPercent := float64(gasUsed) / float64(gasLimit) * 100

	// Log results
	s.logger.WithFields(logrus.Fields{
		"tx_hash":             receipt.TxHash.Hex(),
		"contracts_called":    numCalls,
		"total_bytes_read":    totalBytesRead,
		"gas_used":            gasUsed,
		"gas_limit":           gasLimit,
		"utilization_percent": fmt.Sprintf("%.2f%%", utilizationPercent),
		"status":              receipt.Status,
	}).Info("EXTCODESIZE overload attack completed")

	if receipt.Status != 1 {
		return fmt.Errorf("transaction failed with status %d", receipt.Status)
	}

	return nil
}
