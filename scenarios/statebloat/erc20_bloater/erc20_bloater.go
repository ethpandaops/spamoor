package erc20bloater

import (
	"context"
	"fmt"
	"math/big"
	"strings"
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
	BaseFeeWei       string  `yaml:"base_fee_wei"`
	TipFeeWei        string  `yaml:"tip_fee_wei"`
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
		logger: logger.WithField("scenario", ScenarioName),
	}
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
	flags.Float64Var(&s.options.TargetStorageGB, "target-gb", ScenarioDefaultOptions.TargetStorageGB, "Target storage size in GB")
	flags.Float64Var(&s.options.TargetGasRatio, "target-gas-ratio", ScenarioDefaultOptions.TargetGasRatio, "Target gas usage as ratio of block gas limit (default 0.50 = 50%)")
	flags.Float64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Base fee per gas in gwei")
	flags.Float64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Tip fee per gas in gwei")
	flags.StringVar(&s.options.BaseFeeWei, "basefee-wei", "", "Max fee per gas in wei (overrides --basefee for L2 sub-gwei fees)")
	flags.StringVar(&s.options.TipFeeWei, "tipfee-wei", "", "Max tip per gas in wei (overrides --tipfee for L2 sub-gwei fees)")
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
	blockGasLimit := s.walletPool.GetTxPool().GetCurrentGasLimit()
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
		blockGasLimit = s.walletPool.GetTxPool().GetCurrentGasLimit()
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

		// Recomputed every round (not once for the whole run): the chain's
		// live gas limit can change over a long run, and this also lets the
		// gas-limit-too-high fallback below only affect the round it fires
		// on, not the rest of the run.
		gasLimitPerTx, addressesPerTx := s.computeBatchParams()

		newAddressIndex, batchTxCount, err := s.attemptBloatRound(ctx, blockGasLimit, gasLimitPerTx, addressesPerTx, nextAddressIndex, targetAddresses)
		if err != nil && isGasLimitTooHighError(err) && (gasLimitPerTx != FixedGasLimitPerTx || addressesPerTx != MaxBloatedAddressesPerTx) {
			// The chain rejected the batch for exceeding the gas limit, which
			// means it has not actually activated Amsterdam yet even though
			// spamoor's fee model defaults to assuming it has (see
			// TxPool.IsAmsterdam's doc comment - it is an operator-set
			// preference, not a live check against the connected chain).
			// Retry this one round with the pre-Amsterdam sizing; the next
			// round tries the Amsterdam sizing again, so throughput recovers
			// automatically once the chain actually activates it.
			s.logger.Warnf("batch rejected for exceeding the gas limit (%v); the connected chain does not appear to have activated Amsterdam yet. Retrying this round with pre-Amsterdam batch sizing. Pass --pre-amsterdam-fee-model if this chain will not activate Amsterdam during this run.", err)
			newAddressIndex, batchTxCount, err = s.attemptBloatRound(ctx, blockGasLimit, FixedGasLimitPerTx, MaxBloatedAddressesPerTx, nextAddressIndex, targetAddresses)
		}

		if err != nil {
			s.logger.Errorf("round failed: %v", err)
			errorCount++
			time.Sleep(time.Second * time.Duration(errorCount))
			continue
		}

		nextAddressIndex = newAddressIndex
		totalTxCount += uint64(batchTxCount)
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

// computeBatchParams determines the per-transaction gas limit and how many
// addresses fit within it, from the connected chain's live fork/gas state.
func (s *Scenario) computeBatchParams() (gasLimitPerTx uint64, addressesPerTx uint64) {
	txpool := s.walletPool.GetTxPool()
	return calculateBatchParams(txpool.IsAmsterdam(), txpool.MaxTxGas(), txpool.GetCostPerStateByte(), s.logger)
}

// calculateBatchParams is the pure sizing calculation behind computeBatchParams,
// split out so it can be tested without a live chain connection.
//
// Pre-Amsterdam this just returns the existing fixed constants, unchanged.
//
// Under Amsterdam (EIP-8037), creating a new storage slot also charges
// state-creation gas from a second, separate budget, and a transaction only
// receives a nonzero share of that budget once its total requested Gas
// exceeds the EIP-7825 ceiling (utils.MaxGasLimitPerTx) - the excess becomes
// the state-gas reservoir. A fixed, chain-agnostic gas limit at or below that
// ceiling gets a state-gas reservoir of exactly zero, so every new slot's
// state-creation cost spills entirely into regular gas instead, and a batch
// sized for the old cost model runs out of gas partway through. Both values
// are recomputed from the live per-tx gas ceiling, which is itself chain-
// aware (TxPool.MaxTxGas reports the block gas limit on Amsterdam, since only
// the state dimension is capped at the legacy ceiling there).
func calculateBatchParams(isAmsterdam bool, maxTxGas uint64, costPerStateByte uint64, logger *logrus.Entry) (gasLimitPerTx uint64, addressesPerTx uint64) {
	if !isAmsterdam {
		return FixedGasLimitPerTx, MaxBloatedAddressesPerTx
	}

	const (
		// Regular-gas portion of creating one new (cold) storage slot under
		// Amsterdam's restructured SSTORE costs: the SSTORE_RESET_GAS tier,
		// 5,000 gas, replaces the old SSTORE_SET + cold-access combination.
		sstoreResetGas = uint64(5000)
		// Key + value bytes charged as state-creation gas for each new slot.
		perSlotStateBytes = uint64(64)
		// Warm balanceOf[msg.sender] update, loop arithmetic, margin.
		loopOverhead = uint64(400)
		// Function dispatch, state variable reads, margin.
		callOverhead = uint64(300000)
	)

	regularPerAddr := SlotsPerBloatCycle*sstoreResetGas + loopOverhead
	statePerAddr := SlotsPerBloatCycle * perSlotStateBytes * costPerStateByte

	if statePerAddr == 0 || maxTxGas <= utils.MaxGasLimitPerTx {
		// Either state gas isn't actually priced, or the chain's block gas
		// limit leaves no room to exceed the EIP-7825 ceiling at all - there is
		// nowhere for a state-gas reservoir to come from. Fall back to the
		// pre-Amsterdam sizing so the scenario still makes progress, just
		// without a state-gas reservoir (everything spills into regular gas,
		// which the pre-Amsterdam batch size stays safely under).
		if logger != nil {
			logger.Warnf("cannot size an Amsterdam state-gas reservoir (max tx gas %d); falling back to pre-Amsterdam batch sizing", maxTxGas)
		}
		return FixedGasLimitPerTx, MaxBloatedAddressesPerTx
	}

	addressesPerTx = (maxTxGas - utils.MaxGasLimitPerTx) / statePerAddr
	if addressesPerTx == 0 {
		addressesPerTx = 1
	}
	// The regular-gas need for even a large batch is tiny next to the
	// mandatory EIP-7825 regular allotment (state gas dominates by roughly
	// 20x per address), but guard it explicitly rather than assume it can
	// never bind.
	for addressesPerTx > 1 && regularPerAddr*addressesPerTx+callOverhead > utils.MaxGasLimitPerTx {
		addressesPerTx--
	}

	gasLimitPerTx = utils.MaxGasLimitPerTx + addressesPerTx*statePerAddr

	return gasLimitPerTx, addressesPerTx
}

// calculateTransactionSplits determines how many transactions are needed for the target gas.
// Each transaction uses gasLimitPerTx for simplicity and predictability.
func (s *Scenario) calculateTransactionSplits(totalTargetGas uint64, gasLimitPerTx uint64) []uint64 {
	// Simple calculation: divide total by the per-tx limit
	numTxs := (totalTargetGas + gasLimitPerTx - 1) / gasLimitPerTx // ceiling division

	// All transactions use the same gas limit
	splits := make([]uint64, numTxs)
	for i := range splits {
		splits[i] = gasLimitPerTx
	}

	return splits
}

// isGasLimitTooHighError reports whether err is the rejection a real node
// returns when a transaction's declared Gas exceeds the EIP-7825 ceiling on
// a chain that has not activated Amsterdam - go-ethereum's mempool
// validation returns the literal error "transaction gas limit too high" for
// this case (core.ErrGasLimitTooHigh), and that text survives being wrapped
// on its way back through the RPC client and through
// SendMultiTransactionBatch's own error wrapping.
func isGasLimitTooHighError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "gas limit too high")
}

// attemptBloatRound builds and sends one round of bloat transactions using
// gasLimitPerTx and addressesPerTx, starting from startAddressIndex. It
// returns the address index to resume from and the number of confirmed
// transactions on success.
//
// On any failure before transactions are confirmed on-chain (a build
// failure, or a hard rejection at submission), it releases the nonces this
// attempt allocated via wallet.MarkSkippedNonce so a subsequent retry with
// different parameters can reuse them instead of leaving a permanent gap in
// the wallet's nonce sequence. A transaction that reverted on-chain already
// consumed its nonce legitimately, so those are left alone.
func (s *Scenario) attemptBloatRound(ctx context.Context, blockGasLimit, gasLimitPerTx, addressesPerTx, startAddressIndex, targetAddresses uint64) (newAddressIndex uint64, txCount int, err error) {
	totalTargetGas := uint64(float64(blockGasLimit) * s.options.TargetGasRatio)
	txSplits := s.calculateTransactionSplits(totalTargetGas, gasLimitPerTx)

	if len(txSplits) > 1 {
		s.logger.Infof("splitting target gas (%.1fM) across %d transactions (%.1fM gas each)",
			float64(totalTargetGas)/1000000, len(txSplits), float64(gasLimitPerTx)/1000000)
	} else {
		s.logger.Infof("block gas limit: %d, target gas: %d (%.0f%%) - single tx",
			blockGasLimit, totalTargetGas, s.options.TargetGasRatio*100)
	}

	type txBatch struct {
		tx              *types.Transaction
		wallet          *spamoor.Wallet
		numAddresses    uint64
		gasLimit        uint64
		endAddressIndex uint64
	}
	var txBatches []txBatch
	nextAddressIndex := startAddressIndex

	releaseNonces := func() {
		for _, batch := range txBatches {
			batch.wallet.MarkSkippedNonce(batch.tx.Nonce())
		}
	}

	// Build all transactions first
	for i := range txSplits {
		// Use a different wallet for each split transaction to enable parallel submission
		// Wallet 0 is the deployer, so we use wallets 1, 2, 3, ... for bloating
		walletIndex := i + 1

		// Ensure we don't exceed available wallets
		if walletIndex >= s.options.WalletCount {
			releaseNonces()
			return startAddressIndex, 0, fmt.Errorf("not enough wallets: need %d but only have %d", walletIndex+1, s.options.WalletCount)
		}

		wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, walletIndex)

		// Use the maximum number of addresses that fit within gasLimitPerTx
		numAddresses := addressesPerTx

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
			i+1, len(txSplits), numAddresses, gasLimitPerTx/1000000)

		// Build bloating transaction with calculated number of addresses
		// NOTE: nextAddressIndex is already in address units (matching contract's nextStorageSlot)
		tx, err := s.buildBloatTx(ctx, wallet, nextAddressIndex, numAddresses, gasLimitPerTx)
		if err != nil {
			releaseNonces()
			return startAddressIndex, 0, fmt.Errorf("failed to build batch tx %d/%d: %w", i+1, len(txSplits), err)
		}

		txBatches = append(txBatches, txBatch{
			tx:              tx,
			wallet:          wallet,
			numAddresses:    numAddresses,
			gasLimit:        gasLimitPerTx,
			endAddressIndex: endAddressIndex,
		})

		s.logger.WithFields(logrus.Fields{
			"batch":     fmt.Sprintf("%d/%d", i+1, len(txSplits)),
			"wallet":    s.walletPool.GetWalletName(wallet.GetAddress()),
			"addresses": numAddresses,
			"gas_limit": gasLimitPerTx,
		}).Debugf("built bloat tx")

		nextAddressIndex = endAddressIndex

		// Break if we've reached target
		if nextAddressIndex >= targetAddresses {
			break
		}
	}

	if len(txBatches) == 0 {
		return startAddressIndex, 0, nil
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

	receipts, sendErr := s.walletPool.GetTxPool().SendMultiTransactionBatch(ctx, walletTxMap, &spamoor.BatchOptions{
		SendTransactionOptions: spamoor.SendTransactionOptions{
			Client:      client,
			Rebroadcast: true,
		},
	})
	if sendErr != nil {
		releaseNonces()
		return startAddressIndex, 0, fmt.Errorf("failed to send transaction batch: %w", sendErr)
	}

	// Process receipts
	batchTxCount := 0
	for i, batch := range txBatches {
		walletReceipts := receipts[batch.wallet]
		if len(walletReceipts) == 0 {
			releaseNonces()
			return startAddressIndex, 0, fmt.Errorf("no receipt for batch tx %d/%d", i+1, len(txBatches))
		}

		receipt := walletReceipts[0] // Each wallet sends only one tx in our case
		if receipt.Status != types.ReceiptStatusSuccessful {
			// Mined but reverted: the nonce is legitimately consumed
			// on-chain, so it is not released here.
			return startAddressIndex, 0, fmt.Errorf("tx failed: %s (gas used: %d, gas limit: %d)",
				batch.tx.Hash().Hex(), receipt.GasUsed, batch.tx.Gas())
		}

		s.logger.WithFields(logrus.Fields{
			"batch":    fmt.Sprintf("%d/%d", i+1, len(txBatches)),
			"tx":       batch.tx.Hash().Hex(),
			"gas_used": receipt.GasUsed,
			"block":    receipt.BlockNumber.Uint64(),
		}).Infof("bloat tx confirmed")

		batchTxCount++
	}

	return txBatches[len(txBatches)-1].endAddressIndex, batchTxCount, nil
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

	// Build all transfer transactions upfront
	var txList []*types.Transaction
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
		baseFeeWei, tipFeeWei := spamoor.ResolveFees(s.options.BaseFee, s.options.TipFee, s.options.BaseFeeWei, s.options.TipFeeWei)
		feeCap, tipCap, err := s.walletPool.GetSuggestedFees(client, baseFeeWei, tipFeeWei)
		if err != nil {
			return fmt.Errorf("failed to get suggested fees: %w", err)
		}

		tx, err := deployerWallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
			To:        &s.contractAddr,
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			// Every transfer targets a fresh recipient (state bloat by design), so
			// it creates a new balance slot. Under the Amsterdam fee schedule that
			// costs ~133k gas (vs ~50k pre-Amsterdam); keep static headroom.
			Gas:   200000,
			Value: uint256.NewInt(0),
		}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
			return s.contractInstance.Transfer(transactOpts, recipientAddr, tokensPerWallet)
		})
		if err != nil {
			return fmt.Errorf("failed to build transfer tx for wallet %d: %w", i, err)
		}

		txList = append(txList, tx)
	}

	if len(txList) == 0 {
		s.logger.Infof("all wallets already have sufficient tokens")
		return nil
	}

	// Send all transfers in batch with sliding window
	s.logger.Infof("sending %d token transfer transactions in batch", len(txList))
	receipts, err := s.walletPool.GetTxPool().SendTransactionBatch(ctx, deployerWallet, txList, &spamoor.BatchOptions{
		SendTransactionOptions: spamoor.SendTransactionOptions{
			Client:      client,
			Rebroadcast: true,
		},
		MaxRetries:   3,
		PendingLimit: 200,
	})
	if err != nil {
		return fmt.Errorf("failed to send token transfer batch: %w", err)
	}

	// Verify all transfers succeeded
	for i, receipt := range receipts {
		if receipt == nil || receipt.Status != types.ReceiptStatusSuccessful {
			return fmt.Errorf("token transfer %d failed", i)
		}
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

	baseFeeWei, tipFeeWei := spamoor.ResolveFees(s.options.BaseFee, s.options.TipFee, s.options.BaseFeeWei, s.options.TipFeeWei)
	feeCap, tipCap, err := s.walletPool.GetSuggestedFees(client, baseFeeWei, tipFeeWei)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get suggested fees: %w", err)
	}

	tx, err := wallet.BuildBoundTxWithEstimate(ctx, client, s.walletPool.GetTxPool(), &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
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
func (s *Scenario) buildBloatTx(ctx context.Context, wallet *spamoor.Wallet, startAddressIndex uint64, numAddresses uint64, gasLimit uint64) (*types.Transaction, error) {
	client := s.walletPool.GetClient(spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0))
	if client == nil {
		return nil, fmt.Errorf("no client available")
	}

	baseFeeWei, tipFeeWei := spamoor.ResolveFees(s.options.BaseFee, s.options.TipFee, s.options.BaseFeeWei, s.options.TipFeeWei)
	feeCap, tipCap, err := s.walletPool.GetSuggestedFees(client, baseFeeWei, tipFeeWei)
	if err != nil {
		return nil, fmt.Errorf("failed to get suggested fees: %w", err)
	}

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
