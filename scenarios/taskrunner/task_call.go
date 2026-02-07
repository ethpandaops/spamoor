package taskrunner

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/utils"
	"github.com/holiman/uint256"
)

// CallTask represents a contract function call task
type CallTask struct {
	BaseTask
	Target      string        `yaml:"target" json:"target"`               // Contract address or {contract:name}
	CallData    string        `yaml:"call_data" json:"call_data"`         // Raw hex calldata
	CallABI     string        `yaml:"call_abi" json:"call_abi"`           // JSON ABI string
	CallABIFile string        `yaml:"call_abi_file" json:"call_abi_file"` // Path to ABI file or URL (http/https)
	CallFnName  string        `yaml:"call_fn_name" json:"call_fn_name"`   // Function name
	CallArgs    []interface{} `yaml:"call_args" json:"call_args"`         // Function arguments
	GasLimit    uint64        `yaml:"gas_limit" json:"gas_limit"`         // Gas limit
	Amount      uint64        `yaml:"amount" json:"amount"`               // ETH amount to send (in gwei)
}

// NewCallTask creates a new call task from configuration
func NewCallTask(name string, data map[string]interface{}) (Task, error) {
	task := &CallTask{
		BaseTask: BaseTask{
			Type: "call",
			Name: name,
		},
		GasLimit: 100000, // Default gas limit
	}

	// Parse configuration fields
	var err error
	if task.Target, err = getRequiredString("target", data); err != nil {
		return nil, err
	}

	if err := parseValue("call_data", data, &task.CallData); err != nil {
		return nil, err
	}
	if err := parseValue("call_abi", data, &task.CallABI); err != nil {
		return nil, err
	}
	if err := parseValue("call_abi_file", data, &task.CallABIFile); err != nil {
		return nil, err
	}
	if err := parseValue("call_fn_name", data, &task.CallFnName); err != nil {
		return nil, err
	}
	if err := parseValue("call_args", data, &task.CallArgs); err != nil {
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

// Validate checks if the call task configuration is valid
func (t *CallTask) Validate() error {
	// Must have target
	if t.Target == "" {
		return fmt.Errorf("call task must specify target address or contract reference")
	}

	// Must have either calldata OR (abi + function)
	hasCallData := t.CallData != ""
	hasABIFunction := (t.CallABI != "" || t.CallABIFile != "") && t.CallFnName != ""

	if hasCallData && hasABIFunction {
		return fmt.Errorf("call task cannot specify both call_data and call_abi+call_fn_name")
	}

	// Validate gas limit
	if t.GasLimit == 0 {
		return fmt.Errorf("call task gas_limit must be greater than 0")
	}

	// If calldata is specified, it should be valid hex
	if t.CallData != "" && !common.IsHexAddress("0x"+t.CallData) && !strings.HasPrefix(t.CallData, "0x") {
		return fmt.Errorf("call task call_data must be valid hex string")
	}

	// Validate ABI sources (should have exactly one if using ABI)
	if hasABIFunction {
		sources := 0
		if t.CallABI != "" {
			sources++
		}
		if t.CallABIFile != "" {
			sources++
		}
		if sources != 1 {
			return fmt.Errorf("call task must specify exactly one ABI source (call_abi or call_abi_file)")
		}
	}

	return nil
}

// BuildTransaction creates a contract call transaction
func (t *CallTask) BuildTransaction(ctx context.Context, wallet *spamoor.Wallet, registry *ContractRegistry, execCtx *TaskExecutionContext) (*types.Transaction, error) {
	// Resolve target address (placeholders should already be processed by caller)
	targetAddr, err := registry.ResolveReference(t.Target)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve target address: %w", err)
	}

	// Build call data
	callData, err := t.buildCallData(wallet)
	if err != nil {
		return nil, fmt.Errorf("failed to build call data: %w", err)
	}

	// Get suggested fees from execution context
	baseFeeWei, tipFeeWei := spamoor.ResolveFees(execCtx.BaseFee, execCtx.TipFee, execCtx.BaseFeeWei, execCtx.TipFeeWei)
	feeCap, tipCap, err := execCtx.TxPool.GetSuggestedFees(nil, baseFeeWei, tipFeeWei)
	if err != nil {
		return nil, fmt.Errorf("failed to get suggested fees: %w", err)
	}

	// Create transaction data
	txData := &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       t.GasLimit,
		To:        &targetAddr,
		Data:      callData,
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
		return nil, fmt.Errorf("failed to build call transaction: %w", err)
	}

	return tx, nil
}

// buildCallData constructs the transaction call data
func (t *CallTask) buildCallData(wallet *spamoor.Wallet) ([]byte, error) {
	// If raw calldata is provided, use it directly
	if t.CallData != "" {
		return common.FromHex(t.CallData), nil
	}

	// Otherwise, use ABI + function to build calldata
	if t.CallABI != "" || t.CallABIFile != "" {
		return t.buildABICallData(wallet)
	}

	return nil, nil
}

// buildABICallData builds call data using ABI and function
func (t *CallTask) buildABICallData(wallet *spamoor.Wallet) ([]byte, error) {
	// Load ABI from appropriate source
	abiJSON, err := t.loadABI()
	if err != nil {
		return nil, fmt.Errorf("failed to load ABI: %w", err)
	}

	// Convert arguments to JSON string format expected by ABICallDataBuilder
	argsJSON := "[]"
	if len(t.CallArgs) > 0 {
		// For ABI building, we'll process placeholders separately during actual execution
		// For now, use the args as-is since this is just for ABI validation
		processedArgs := t.CallArgs

		// Convert processed arguments to JSON string
		argsBytes, err := json.Marshal(processedArgs)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal arguments: %w", err)
		}
		argsJSON = string(argsBytes)
	}

	// Create ABI call data builder using the constructor
	abiBuilder, err := utils.NewABICallDataBuilder(string(abiJSON), t.CallFnName, "", argsJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to create ABI call builder: %w", err)
	}

	// Build the call data (BuildCallData expects a txIdx parameter)
	callData, err := abiBuilder.BuildCallData(0) // Use 0 as txIdx since we handle placeholders ourselves
	if err != nil {
		return nil, fmt.Errorf("failed to build call data: %w", err)
	}

	return callData, nil
}

// loadABI loads ABI from the configured source
func (t *CallTask) loadABI() ([]byte, error) {
	// Direct ABI string
	if t.CallABI != "" {
		return []byte(t.CallABI), nil
	}

	// File or URL source
	if t.CallABIFile != "" {
		if t.CallABI != "" {
			return nil, fmt.Errorf("only one of call_abi or call_abi_file can be set")
		}

		var abiBytes []byte
		var err error

		// Check if it's a URL (same logic as calltx)
		if strings.HasPrefix(t.CallABIFile, "https://") || strings.HasPrefix(t.CallABIFile, "http://") {
			resp, err := http.Get(t.CallABIFile)
			if err != nil {
				return nil, fmt.Errorf("could not load ABI file: %w", err)
			}
			defer resp.Body.Close()
			abiBytes, err = io.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("could not read ABI file: %w", err)
			}
		} else {
			// Local file
			abiBytes, err = os.ReadFile(t.CallABIFile)
			if err != nil {
				return nil, fmt.Errorf("could not read ABI file: %w", err)
			}
		}
		return abiBytes, nil
	}

	return nil, fmt.Errorf("no ABI source configured")
}

// processArguments processes function arguments with placeholder substitution
func (t *CallTask) processArguments(args []interface{}, wallet *spamoor.Wallet, registry *ContractRegistry, txIdx uint64, stepIdx int) ([]interface{}, error) {
	processed := make([]interface{}, len(args))

	for i, arg := range args {
		switch v := arg.(type) {
		case string:
			// Apply placeholder substitution
			processedStr, err := t.processPlaceholders(v, registry, txIdx, stepIdx)
			if err != nil {
				return nil, fmt.Errorf("failed to process placeholders in argument %d: %w", i, err)
			}
			processed[i] = processedStr
		default:
			// Use as-is for non-string arguments
			processed[i] = v
		}
	}

	return processed, nil
}

// processPlaceholders processes placeholders in strings
func (t *CallTask) processPlaceholders(str string, registry *ContractRegistry, txIdx uint64, stepIdx int) (string, error) {
	// Process contract placeholders (with 0x prefix for call tasks)
	processed, err := ProcessContractPlaceholders(str, registry, false)
	if err != nil {
		return "", err
	}

	// Process basic placeholders (with 0x prefix for call tasks)
	processed, err = ProcessBasicPlaceholders(processed, txIdx, stepIdx, false)
	if err != nil {
		return "", err
	}

	return processed, nil
}
