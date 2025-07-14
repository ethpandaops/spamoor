package taskrunner

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/holiman/uint256"
)

// DeployTask represents a contract deployment task
type DeployTask struct {
	BaseTask
	ContractCode string `yaml:"contract_code" json:"contract_code"` // Hex-encoded bytecode
	ContractFile string `yaml:"contract_file" json:"contract_file"` // Path to bytecode file or URL (http/https)
	ContractArgs string `yaml:"contract_args" json:"contract_args"` // Constructor arguments (hex)
	GasLimit     uint64 `yaml:"gas_limit" json:"gas_limit"`         // Gas limit for deployment
	Amount       uint64 `yaml:"amount" json:"amount"`               // ETH amount to send (in gwei)
}

// NewDeployTask creates a new deploy task from configuration
func NewDeployTask(name string, data map[string]interface{}) (Task, error) {
	task := &DeployTask{
		BaseTask: BaseTask{
			Type: "deploy",
			Name: name,
		},
		GasLimit: 2000000, // Default gas limit
	}

	// Parse configuration fields
	if err := parseValue("contract_code", data, &task.ContractCode); err != nil {
		return nil, err
	}
	if err := parseValue("contract_file", data, &task.ContractFile); err != nil {
		return nil, err
	}
	if err := parseValue("contract_args", data, &task.ContractArgs); err != nil {
		return nil, err
	}
	if err := parseValue("gas_limit", data, &task.GasLimit); err != nil {
		return nil, err
	}
	if err := parseValue("amount", data, &task.Amount); err != nil {
		return nil, err
	}

	return task, nil
}

// Validate checks if the deploy task configuration is valid
func (t *DeployTask) Validate() error {
	// Must have exactly one bytecode source
	sources := 0
	if t.ContractCode != "" {
		sources++
	}
	if t.ContractFile != "" {
		sources++
	}

	if sources == 0 {
		return fmt.Errorf("deploy task must specify bytecode source (contract_code or contract_file)")
	}
	if sources > 1 {
		return fmt.Errorf("deploy task must specify exactly one bytecode source")
	}

	// Validate gas limit
	if t.GasLimit == 0 {
		return fmt.Errorf("deploy task gas_limit must be greater than 0")
	}

	// If args is specified, it should be valid hex
	if t.ContractArgs != "" && !strings.HasPrefix(t.ContractArgs, "0x") && len(t.ContractArgs)%2 != 0 {
		return fmt.Errorf("deploy task contract_args must be valid hex string")
	}

	return nil
}

// BuildTransaction creates a contract deployment transaction
func (t *DeployTask) BuildTransaction(ctx context.Context, wallet *spamoor.Wallet, registry *ContractRegistry, execCtx *TaskExecutionContext) (*types.Transaction, error) {
	// Load bytecode from appropriate source (placeholders already processed)
	bytecode, err := t.loadBytecode()
	if err != nil {
		return nil, fmt.Errorf("failed to load bytecode: %w", err)
	}

	// Prepare deployment data (bytecode + constructor args, placeholders already processed)
	deployData := bytecode
	if t.ContractArgs != "" {
		// Remove 0x prefix if present and convert to bytes
		argBytes := common.FromHex(t.ContractArgs)
		deployData = append(deployData, argBytes...)
	}

	// Get suggested fees from execution context
	feeCap, tipCap, err := execCtx.TxPool.GetSuggestedFees(nil, execCtx.BaseFee, execCtx.TipFee)
	if err != nil {
		return nil, fmt.Errorf("failed to get suggested fees: %w", err)
	}

	// Create transaction data
	txData := &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       t.GasLimit,
		To:        nil, // nil for contract deployment
		Data:      deployData,
	}
	// Convert amount from gwei to wei
	amount := uint256.NewInt(t.Amount)
	amount = amount.Mul(amount, uint256.NewInt(1000000000))
	txData.Value = amount

	// Build dynamic fee transaction
	dynFeeTx, err := txbuilder.DynFeeTx(txData)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic fee transaction: %w", err)
	}

	// Build the transaction
	tx, err := wallet.BuildDynamicFeeTx(dynFeeTx)
	if err != nil {
		return nil, fmt.Errorf("failed to build deployment transaction: %w", err)
	}

	// Calculate contract address for registry
	contractAddr := crypto.CreateAddress(wallet.GetAddress(), tx.Nonce())

	// Register the contract address immediately (pre-calculated)
	if t.Name != "" {
		registry.Set(t.Name, contractAddr)
	}

	return tx, nil
}

// loadBytecode loads bytecode from the configured source (placeholders already processed)
func (t *DeployTask) loadBytecode() ([]byte, error) {
	// Direct bytecode string (placeholders already processed)
	if t.ContractCode != "" {
		return common.FromHex(t.ContractCode), nil
	}

	// File or URL source
	if t.ContractFile != "" {
		var bytecodeBytes []byte

		// Check if it's a URL (same logic as calltx)
		if strings.HasPrefix(t.ContractFile, "https://") || strings.HasPrefix(t.ContractFile, "http://") {
			resp, err := http.Get(t.ContractFile)
			if err != nil {
				return nil, fmt.Errorf("could not load contract file: %w", err)
			}
			defer resp.Body.Close()
			contractCodeHex, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("could not read contract file: %w", err)
			}
			bytecodeBytes = common.FromHex(strings.Trim(string(contractCodeHex), "\r\n\t "))
		} else {
			// Local file
			code, err := os.ReadFile(t.ContractFile)
			if err != nil {
				return nil, fmt.Errorf("could not read contract file: %w", err)
			}
			bytecodeBytes = common.FromHex(strings.Trim(string(code), "\r\n\t "))
		}
		return bytecodeBytes, nil
	}

	return nil, fmt.Errorf("no bytecode source configured")
}

// processPlaceholders processes placeholders in strings, ensuring no 0x prefixes in calldata/bytecode
func (t *DeployTask) processPlaceholders(str string, registry *ContractRegistry, txIdx uint64, stepIdx int) (string, error) {
	// Process contract placeholders (without 0x prefix for deploy bytecode/calldata)
	processed, err := ProcessContractPlaceholders(str, registry, true)
	if err != nil {
		return "", err
	}

	// Process basic placeholders (without 0x prefix for deploy bytecode/calldata)
	processed, err = ProcessBasicPlaceholders(processed, txIdx, stepIdx, true)
	if err != nil {
		return "", err
	}

	return processed, nil
}
