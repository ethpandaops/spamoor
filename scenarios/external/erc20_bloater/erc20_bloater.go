package main

import (
	"context"
	"fmt"
	"math/big"
	"time"

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

	// EIP-7825 transaction gas limits
	MaxGasLimitPerTx   = 16777216 // EIP-7825 maximum: exactly 2^24
	FixedGasLimitPerTx = 16700000 // Set slightly below max to ensure transaction success

	// MaxBloatedAddressesPerTx is the maximum number of addresses we can bloat in a single transaction
	// while staying under the EIP-7825 gas limit.
	//
	// Gas cost breakdown per address iteration in bloatStorage():
	//   - SSTORE to balanceOf[targetAddr]:
	//     * Cold address (first time): 22,100 gas (2,900 cold account + 20,000 SSTORE new slot)
	//     * Warm address (subsequent): 2,900 gas (100 warm account + 2,900 SSTORE existing slot)
	//   - SSTORE to allowance[sender][targetAddr]:
	//     * Cold mapping: 22,100 gas (similar to above)
	//     * Warm mapping: 2,900 gas
	//   - SSTORE to balanceOf[sender] (once per tx): 2,900 gas (warm storage)
	//   - Loop overhead (arithmetic, memory): ~200 gas per iteration
	//
	// Total per address (cold): ~44,400 gas
	// Total per address (warm): ~6,000 gas
	//
	// For maximum efficiency with cold addresses:
	// 16,700,000 / 44,400 ≈ 376 addresses
	// We use 370 to leave a safety margin.
	MaxBloatedAddressesPerTx = 370
)

type ScenarioOptions struct {
	TargetStorageGB  float64 `yaml:"target_storage_gb" json:"target_storage_gb"`
	TargetGasRatio   float64 `yaml:"target_gas_ratio" json:"target_gas_ratio"`
	BaseFee          float64 `yaml:"base_fee" json:"base_fee"`
	TipFee           float64 `yaml:"tip_fee" json:"tip_fee"`
	ExistingContract string  `yaml:"existing_contract" json:"existing_contract"` // Optional override for edge cases
	WalletCount      int     `yaml:"wallet_count" json:"wallet_count"`           // Number of wallets to initialize
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
	WalletCount:      50,
}
var ScenarioDescriptor = scenario.Descriptor{
	Name:           ScenarioName,
	Description:    "Bloat ERC20 contract storage to target GB size using sequential addresses",
	DefaultOptions: ScenarioDefaultOptions,
	NewScenario:    newScenario,
}

func newScenario(logger logrus.FieldLogger) scenario.Scenario {
	return &Scenario{
		options: ScenarioDefaultOptions,
		logger:  logger.WithField("scenario", ScenarioName),
	}
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
	flags.Float64Var(&s.options.TargetStorageGB, "target-gb", ScenarioDefaultOptions.TargetStorageGB, "Target storage size in GB")
	flags.Float64Var(&s.options.TargetGasRatio, "target-gas-ratio", ScenarioDefaultOptions.TargetGasRatio, "Target gas usage as ratio of block gas limit (default 0.50 = 50%)")
	flags.Float64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Base fee per gas in gwei")
	flags.Float64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Tip fee per gas in gwei")
	flags.StringVar(&s.options.ExistingContract, "existing-contract", ScenarioDefaultOptions.ExistingContract, "(Optional) Override contract address for edge cases")
	flags.IntVar(&s.options.WalletCount, "wallet-count", ScenarioDefaultOptions.WalletCount, "Number of wallets to initialize for parallel execution")
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

	// Initialize multiple wallets for parallel execution
	if s.options.WalletCount < 1 {
		s.options.WalletCount = 50 // Ensure at least 1 wallet, default to 50
	}
	s.walletPool.SetWalletCount(uint64(s.options.WalletCount))
	s.logger.Infof("initialized %d wallets for parallel execution", s.options.WalletCount)

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

	var nextAddressIndex uint64 = 1 // Default: start from address 0x01 (matches contract's nextStorageSlot)

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
			nextAddressIndex = 1
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

	// Distribute tokens to wallets for parallel execution
	// Calculate how many wallets we might need based on gas limits
	blockGasLimit, err := s.walletPool.GetTxPool().GetCurrentGasLimitWithInit()
	if err != nil {
		return fmt.Errorf("failed to get current gas limit: %w", err)
	}
	totalTargetGas := uint64(float64(blockGasLimit) * s.options.TargetGasRatio)
	maxSplits := (totalTargetGas + utils.MaxGasLimitPerTx - 1) / utils.MaxGasLimitPerTx // ceiling division
	walletsNeeded := int(maxSplits) + 1                                                 // +1 for deployer wallet

	// Use the minimum of walletsNeeded and configured wallet count
	walletsToFund := walletsNeeded
	if walletsToFund > s.options.WalletCount {
		walletsToFund = s.options.WalletCount
	}

	// Distribute tokens if this is a fresh deployment or manual contract
	if wallet.GetNonce() == 1 || s.options.ExistingContract != "" {
		s.logger.Infof("distributing tokens to %d wallets for parallel execution", walletsToFund-1)
		if err := s.distributeTokensToWallets(ctx, walletsToFund); err != nil {
			return fmt.Errorf("failed to distribute tokens: %w", err)
		}
	}

	// Query on-chain progress from contract (if resuming)
	if wallet.GetNonce() > 0 || s.options.ExistingContract != "" {
		nextSlot, err := contractInstance.NextStorageSlot(nil)
		if err != nil {
			return fmt.Errorf("failed to query nextStorageSlot from contract: %w", err)
		}

		nextAddressIndex = nextSlot.Uint64()
		if nextAddressIndex == 0 {
			nextAddressIndex = 1 // Contract not yet bloated, start from slot 1
		}

		// Calculate and log current progress
		// Note: nextAddressIndex is in address units; each address = 2 storage slots = 64 bytes
		targetBytes := uint64(s.options.TargetStorageGB * 1024 * 1024 * 1024)
		targetAddresses := targetBytes / (BytesPerSlot * SlotsPerBloatCycle)
		currentGB := float64(nextAddressIndex*SlotsPerBloatCycle*BytesPerSlot) / (1024 * 1024 * 1024)
		progress := float64(nextAddressIndex) / float64(targetAddresses) * 100

		s.logger.Infof("resuming from on-chain state: contract %s | address %d (%.2f%% complete, %.3f GB / %.3f GB)",
			s.contractAddr.Hex(), nextAddressIndex, progress, currentGB, s.options.TargetStorageGB)
	}

	// Query network gas limit (reuse existing if already fetched)
	if blockGasLimit == 0 {
		blockGasLimit, err = s.walletPool.GetTxPool().GetCurrentGasLimitWithInit()
		if err != nil {
			return fmt.Errorf("failed to get current gas limit: %w", err)
		}
	}

	// Calculate target addresses needed (each address = 2 storage slots = 64 bytes)
	targetBytes := uint64(s.options.TargetStorageGB * 1024 * 1024 * 1024)
	targetAddresses := targetBytes / (BytesPerSlot * SlotsPerBloatCycle)
	s.logger.Infof("target: %.2f GB = %d addresses (%.2f million)",
		s.options.TargetStorageGB, targetAddresses, float64(targetAddresses)/1000000)

	// Start bloating with EIP-7825 compliant transaction splitting
	totalTxCount := uint64(0)
	errorCount := 0

	for nextAddressIndex < targetAddresses {
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
			s.logger.Infof("splitting target gas (%.1fM) across %d transactions (fixed %.1fM gas each)",
				float64(totalTargetGas)/1000000, len(txSplits), float64(FixedGasLimitPerTx)/1000000)
		} else {
			s.logger.Infof("block gas limit: %d, target gas: %d (%.0f%%) - single tx",
				blockGasLimit, totalTargetGas, s.options.TargetGasRatio*100)
		}

		// Process batch of transactions for this round
		roundStartAddressIndex := nextAddressIndex
		roundSuccess := true
		batchTxCount := 0

		// Structure to hold transaction data for parallel processing
		type txBatch struct {
			tx              *types.Transaction
			wallet          *spamoor.Wallet
			numAddresses    uint64
			gasLimit        uint64
			endAddressIndex uint64
		}
		var txBatches []txBatch

		// Build all transactions first
		for i := range txSplits {
			// Use a different wallet for each split transaction to enable parallel submission
			// Wallet 0 is the deployer, so we use wallets 1, 2, 3, ... for bloating
			walletIndex := i + 1

			// Ensure we don't exceed available wallets
			if walletIndex >= s.options.WalletCount {
				s.logger.Errorf("not enough wallets: need %d but only have %d", walletIndex+1, s.options.WalletCount)
				roundSuccess = false
				break
			}

			wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, walletIndex)

			// Use the maximum number of addresses that fit within EIP-7825 limit
			numAddresses := uint64(MaxBloatedAddressesPerTx)

			// Check if we would exceed our target addresses
			endAddressIndex := nextAddressIndex + numAddresses
			if endAddressIndex > targetAddresses {
				endAddressIndex = targetAddresses
				numAddresses = endAddressIndex - nextAddressIndex
			}

			if numAddresses == 0 {
				break // No more addresses to process
			}

			s.logger.Debugf("batch %d/%d: processing %d addresses (max per tx) with %dM gas limit",
				i+1, len(txSplits), numAddresses, FixedGasLimitPerTx/1000000)

			// Build bloating transaction with calculated number of addresses
			// NOTE: nextAddressIndex is already in address units (matching contract's nextStorageSlot)
			tx, err := s.buildBloatTx(ctx, wallet, nextAddressIndex, numAddresses)
			if err != nil {
				s.logger.Errorf("failed to build batch tx %d/%d: %v", i+1, len(txSplits), err)
				roundSuccess = false
				break
			}

			txBatches = append(txBatches, txBatch{
				tx:              tx,
				wallet:          wallet,
				numAddresses:    numAddresses,
				gasLimit:        FixedGasLimitPerTx,
				endAddressIndex: endAddressIndex,
			})

			s.logger.WithFields(logrus.Fields{
				"batch":     fmt.Sprintf("%d/%d", i+1, len(txSplits)),
				"wallet":    s.walletPool.GetWalletName(wallet.GetAddress()),
				"addresses": numAddresses,
				"gas_limit": FixedGasLimitPerTx,
			}).Debugf("built bloat tx")

			nextAddressIndex = endAddressIndex

			// Break if we've reached target
			if nextAddressIndex >= targetAddresses {
				break
			}
		}

		if !roundSuccess {
			// Revert to beginning of round on failure
			nextAddressIndex = roundStartAddressIndex
			errorCount++
			time.Sleep(time.Second * time.Duration(errorCount))
			continue
		}

		// Prepare wallet-to-transactions map for SendMultiTransactionBatch
		walletTxMap := make(map[*spamoor.Wallet][]*types.Transaction)
		for _, batch := range txBatches {
			walletTxMap[batch.wallet] = append(walletTxMap[batch.wallet], batch.tx)

			s.logger.WithFields(logrus.Fields{
				"wallet":    s.walletPool.GetWalletName(batch.wallet.GetAddress()),
				"tx":        batch.tx.Hash().Hex(),
				"nonce":     batch.tx.Nonce(),
				"addresses": batch.numAddresses,
				"gas_limit": batch.gasLimit,
			}).Debugf("prepared bloat tx for batch sending")
		}

		// Send all transactions in parallel using SendMultiTransactionBatch
		s.logger.Infof("sending %d transactions in parallel from %d wallets", len(txBatches), len(walletTxMap))

		client := s.walletPool.GetClient(spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0))

		receipts, err := s.walletPool.GetTxPool().SendMultiTransactionBatch(ctx, walletTxMap, &spamoor.BatchOptions{
			SendTransactionOptions: spamoor.SendTransactionOptions{
				Client:      client,
				Rebroadcast: true,
			},
		})
		if err != nil {
			s.logger.Errorf("failed to send transaction batch: %v", err)
			roundSuccess = false
		} else {
			// Process receipts
			for i, batch := range txBatches {
				walletReceipts := receipts[batch.wallet]
				if len(walletReceipts) == 0 {
					s.logger.Errorf("no receipt for batch tx %d/%d", i+1, len(txBatches))
					roundSuccess = false
					break
				}

				receipt := walletReceipts[0] // Each wallet sends only one tx in our case
				if receipt.Status != types.ReceiptStatusSuccessful {
					s.logger.Errorf("tx failed: %s (gas used: %d, gas limit: %d)",
						batch.tx.Hash().Hex(), receipt.GasUsed, batch.tx.Gas())
					roundSuccess = false
					break
				}

				s.logger.WithFields(logrus.Fields{
					"batch":    fmt.Sprintf("%d/%d", i+1, len(txBatches)),
					"tx":       batch.tx.Hash().Hex(),
					"gas_used": receipt.GasUsed,
					"block":    receipt.BlockNumber.Uint64(),
				}).Infof("bloat tx confirmed")

				batchTxCount++
			}

			if roundSuccess {
				nextAddressIndex = txBatches[len(txBatches)-1].endAddressIndex
				totalTxCount += uint64(batchTxCount)
			}
		}

		if !roundSuccess {
			// Revert to beginning of round on failure
			nextAddressIndex = roundStartAddressIndex
			errorCount++
			time.Sleep(time.Second * time.Duration(errorCount))
			continue
		}

		// Reset error count on successful round
		errorCount = 0

		// Log progress after successful round
		// Note: each address = 2 storage slots = 64 bytes
		currentGB := float64(nextAddressIndex*SlotsPerBloatCycle*BytesPerSlot) / (1024 * 1024 * 1024)
		progress := float64(nextAddressIndex) / float64(targetAddresses) * 100
		s.logger.Infof("progress: %.2f%% | contract: %s | addresses: %d / %d | storage: %.3f GB / %.3f GB | round txs: %d",
			progress, s.contractAddr.Hex(), nextAddressIndex, targetAddresses, currentGB, s.options.TargetStorageGB, batchTxCount)
	}

	// Log completion
	finalGB := float64(nextAddressIndex*SlotsPerBloatCycle*BytesPerSlot) / (1024 * 1024 * 1024)
	s.logger.Infof("bloating complete! contract: %s | total addresses: %d | estimated storage: %.3f GB | total txs: %d",
		s.contractAddr.Hex(), nextAddressIndex, finalGB, totalTxCount)

	return nil
}

// calculateTransactionSplits determines how many transactions are needed for the target gas
// Each transaction uses the fixed gas limit for simplicity and predictability
func (s *Scenario) calculateTransactionSplits(totalTargetGas uint64) []uint64 {
	// Simple calculation: divide total by fixed limit
	numTxs := (totalTargetGas + FixedGasLimitPerTx - 1) / FixedGasLimitPerTx // ceiling division

	// All transactions use the same fixed gas limit
	splits := make([]uint64, numTxs)
	for i := range splits {
		splits[i] = FixedGasLimitPerTx
	}

	return splits
}

// distributeTokensToWallets distributes tokens from wallet 0 to other wallets for parallel execution
func (s *Scenario) distributeTokensToWallets(ctx context.Context, numWallets int) error {
	if numWallets <= 1 {
		return nil // No distribution needed if only using deployer wallet
	}

	client := s.walletPool.GetClient(spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0))
	deployerWallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, 0)

	// 10 million tokens per wallet (with 18 decimals)
	tokensPerWallet := new(big.Int)
	tokensPerWallet.SetString("10000000000000000000000000", 10) // 10M * 10^18

	s.logger.Infof("distributing 10M tokens to each of %d wallets", numWallets-1)

	for i := 1; i < numWallets; i++ {
		recipientWallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, i)
		recipientAddr := recipientWallet.GetAddress()

		// Check if wallet already has tokens (in case of resume)
		balance, err := s.contractInstance.BalanceOf(nil, recipientAddr)
		if err != nil {
			return fmt.Errorf("failed to check balance for wallet %d: %w", i, err)
		}

		if balance.Cmp(tokensPerWallet) >= 0 {
			s.logger.Debugf("wallet %d already has sufficient tokens, skipping", i)
			continue
		}

		// Build transfer transaction
		feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(client, s.options.BaseFee, s.options.TipFee)
		if err != nil {
			return fmt.Errorf("failed to get suggested fees: %w", err)
		}

		tx, err := deployerWallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
			To:        &s.contractAddr,
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Gas:       100000, // Simple transfer shouldn't need much gas
			Value:     uint256.NewInt(0),
		}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
			return s.contractInstance.Transfer(transactOpts, recipientAddr, tokensPerWallet)
		})
		if err != nil {
			return fmt.Errorf("failed to build transfer tx for wallet %d: %w", i, err)
		}

		// Send and wait for confirmation
		receipt, err := s.walletPool.GetTxPool().SendAndAwaitTransaction(ctx, deployerWallet, tx, &spamoor.SendTransactionOptions{
			Client:      client,
			Rebroadcast: true,
		})
		if err != nil {
			return fmt.Errorf("failed to send transfer to wallet %d: %w", i, err)
		}

		if receipt.Status != types.ReceiptStatusSuccessful {
			return fmt.Errorf("token transfer to wallet %d failed", i)
		}

		s.logger.Debugf("transferred 10M tokens to wallet %d (tx: %s)", i, tx.Hash().Hex())
	}

	s.logger.Infof("token distribution complete")
	return nil
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

	tx, err := wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       2000000,
		Value:     uint256.NewInt(0),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		_, deployTx, _, err := contract.DeployERC20Bloater(transactOpts, client.GetEthClient(), initialSupply)
		if err != nil {
			return nil, err
		}
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

	return receipt, tx, nil
}

// buildBloatTx builds a bloating transaction without sending it.
// nextAddressIndex is the starting address index (matching contract's nextStorageSlot semantics).
func (s *Scenario) buildBloatTx(ctx context.Context, wallet *spamoor.Wallet, startAddressIndex uint64, numAddresses uint64) (*types.Transaction, error) {
	client := s.walletPool.GetClient(spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0))
	if client == nil {
		return nil, fmt.Errorf("no client available")
	}

	feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(client, s.options.BaseFee, s.options.TipFee)
	if err != nil {
		return nil, fmt.Errorf("failed to get suggested fees: %w", err)
	}

	// Use fixed gas limit for simplicity and predictability
	gasLimit := uint64(FixedGasLimitPerTx)

	// Build transaction using BuildBoundTx
	tx, err := wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
		To:        &s.contractAddr,
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       gasLimit,
		Value:     uint256.NewInt(0),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		return s.contractInstance.BloatStorage(transactOpts, new(big.Int).SetUint64(startAddressIndex), new(big.Int).SetUint64(numAddresses))
	})
	if err != nil {
		return nil, fmt.Errorf("failed to build bloat tx: %w", err)
	}

	return tx, nil
}

// main is a placeholder required for package main compilation.
// This file is loaded dynamically by Yaegi at runtime, not executed directly.
func main() {}
