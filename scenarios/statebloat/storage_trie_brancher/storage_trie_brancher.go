package storagetriebrancher

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
)

const (
	// Nick's factory address for CREATE2 deployments
	NickFactoryAddress = "0x4e59b44847b379578588920ca78fbf26c0b4956c"

	// Gas limits
	DeployGasLimit = 3000000
	FundGasLimit   = 21000

	// Funding amount for auxiliary accounts (1 wei to create the account)
	FundingAmount = 1
)

type ScenarioOptions struct {
	TotalContracts uint64  `yaml:"total_contracts"`
	StorageDepth   uint64  `yaml:"storage_depth"`
	AccountDepth   uint64  `yaml:"account_depth"`
	MaxWallets     uint64  `yaml:"max_wallets"`
	SkipContracts  bool    `yaml:"skip_contracts"`
	SkipFunding    bool    `yaml:"skip_funding"`
	DataFile       string  `yaml:"data_file"`
	ContractFile   string  `yaml:"contract_file"`
	BaseFee        float64 `yaml:"base_fee"`
	TipFee         float64 `yaml:"tip_fee"`
	ClientGroup    string  `yaml:"client_group"`
	LogTxs         bool    `yaml:"log_txs"`
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	// Deployment data
	deployData     *DeploymentData
	initCode       []byte
	deployments    []DeploymentInfo
	factoryAddress common.Address // Actual factory address (might differ from Nick's canonical address)

	// Progress tracking
	fundedAccounts    map[string]bool
	deployedContracts map[string]bool
	currentIndex      uint64
	phase             string // "funding" or "deploying"
	mu                sync.RWMutex
}

type DeploymentData struct {
	Deployer     string         `json:"deployer"`
	InitCodeHash string         `json:"init_code_hash"`
	TargetDepth  int            `json:"target_depth"`
	NumContracts int            `json:"num_contracts"`
	Contracts    []ContractData `json:"contracts"`
}

type ContractData struct {
	Salt              interface{} `json:"salt"` // Can be number or string
	ContractAddress   string      `json:"contract_address"`
	AuxiliaryAccounts []string    `json:"auxiliary_accounts"`
}

type DeploymentInfo struct {
	Address           string   `json:"address"`
	Salt              string   `json:"salt"`
	AuxiliaryAccounts []string `json:"auxiliary_accounts"`
}

var ScenarioName = "storage-trie-brancher"
var ScenarioDefaultOptions = ScenarioOptions{
	TotalContracts: 1000,
	StorageDepth:   9,
	AccountDepth:   3,
	MaxWallets:     50,
	SkipContracts:  false,
	SkipFunding:    false,
	DataFile:       "",
	ContractFile:   "",
	BaseFee:        20,
	TipFee:         2,
	ClientGroup:    "",
	LogTxs:         false,
}

var ScenarioDescriptor = scenario.Descriptor{
	Name:           ScenarioName,
	Description:    "Deploy worst-case depth attack contracts (deep storage tries with auxiliary accounts)",
	DefaultOptions: ScenarioDefaultOptions,
	NewScenario:    newScenario,
}

func newScenario(logger logrus.FieldLogger) scenario.Scenario {
	return &Scenario{
		options:           ScenarioDefaultOptions,
		logger:            logger.WithField("scenario", ScenarioName),
		fundedAccounts:    make(map[string]bool),
		deployedContracts: make(map[string]bool),
	}
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
	flags.Uint64VarP(&s.options.TotalContracts, "count", "c", ScenarioDefaultOptions.TotalContracts, "Total number of contracts to deploy")
	flags.Uint64Var(&s.options.StorageDepth, "storage-depth", ScenarioDefaultOptions.StorageDepth, "Storage trie depth (9 or 10)")
	flags.Uint64Var(&s.options.AccountDepth, "account-depth", ScenarioDefaultOptions.AccountDepth, "Account trie depth (3, 4, or 5)")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of wallets to use for parallel execution")
	flags.BoolVar(&s.options.SkipContracts, "skip-contracts", ScenarioDefaultOptions.SkipContracts, "Skip contract deployment (only fund EOAs)")
	flags.BoolVar(&s.options.SkipFunding, "skip-funding", ScenarioDefaultOptions.SkipFunding, "Skip EOA funding (only deploy contracts)")
	flags.StringVar(&s.options.DataFile, "data-file", ScenarioDefaultOptions.DataFile, "Path or URL to CREATE2 data JSON file (auto-detected if empty)")
	flags.StringVar(&s.options.ContractFile, "contract-file", ScenarioDefaultOptions.ContractFile, "Path or URL to Solidity contract file (auto-detected if empty)")
	flags.Float64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas (in gwei)")
	flags.Float64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas (in gwei)")
	flags.StringVar(&s.options.ClientGroup, "client-group", ScenarioDefaultOptions.ClientGroup, "Client group to use for sending transactions")
	flags.BoolVar(&s.options.LogTxs, "log-txs", ScenarioDefaultOptions.LogTxs, "Log all submitted transactions")
	return nil
}

func (s *Scenario) Init(options *scenario.Options) error {
	s.walletPool = options.WalletPool

	if options.Config != "" {
		// Use the generalized config validation and parsing helper
		err := scenario.ParseAndValidateConfig(&ScenarioDescriptor, options.Config, &s.options, s.logger)
		if err != nil {
			return err
		}
	}

	// Validate options
	if s.options.SkipContracts && s.options.SkipFunding {
		return fmt.Errorf("cannot skip both contracts and funding - nothing to do")
	}

	if s.options.TotalContracts == 0 {
		return fmt.Errorf("total contracts must be greater than 0")
	}

	// Set up wallets for parallel execution
	// Use the configured MaxWallets, or default to 50 if not specified
	if s.options.MaxWallets > 0 {
		s.walletPool.SetWalletCount(s.options.MaxWallets)
	} else {
		s.walletPool.SetWalletCount(50) // Default fallback
	}

	// Auto-detect data file if not specified
	if s.options.DataFile == "" {
		s.options.DataFile = filepath.Join("scenarios", "statebloat", "storage_trie_brancher",
			fmt.Sprintf("s%d_acc%d.json", s.options.StorageDepth, s.options.AccountDepth))
	}

	// Auto-detect contract file if not specified
	if s.options.ContractFile == "" {
		s.options.ContractFile = filepath.Join("scenarios", "statebloat", "storage_trie_brancher",
			fmt.Sprintf("depth_%d.sol", s.options.StorageDepth))
	}

	// Load deployment data
	if err := s.loadDeploymentData(); err != nil {
		return fmt.Errorf("failed to load deployment data: %w", err)
	}

	// Compile contract to get init code
	if err := s.compileContract(); err != nil {
		return fmt.Errorf("failed to compile contract: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"contracts":     s.options.TotalContracts,
		"storage_depth": s.options.StorageDepth,
		"account_depth": s.options.AccountDepth,
		"data_file":     s.options.DataFile,
		"contract_file": s.options.ContractFile,
	}).Info("initialized storage trie brancher scenario")

	return nil
}

// loadDataFromPathOrURL loads data from a file path or URL
func (s *Scenario) loadDataFromPathOrURL(pathOrURL string) ([]byte, error) {
	if strings.HasPrefix(pathOrURL, "http://") || strings.HasPrefix(pathOrURL, "https://") {
		// Load from URL
		s.logger.WithField("url", pathOrURL).Debug("Loading data from URL")
		resp, err := http.Get(pathOrURL)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch URL: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("HTTP error: %s", resp.Status)
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		}
		return data, nil
	}

	// Load from file
	s.logger.WithField("file", pathOrURL).Debug("Loading data from file")
	return os.ReadFile(pathOrURL)
}

func (s *Scenario) loadDeploymentData() error {
	data, err := s.loadDataFromPathOrURL(s.options.DataFile)
	if err != nil {
		return fmt.Errorf("failed to load data file: %w", err)
	}

	s.deployData = &DeploymentData{}
	if err := json.Unmarshal(data, s.deployData); err != nil {
		return fmt.Errorf("failed to parse JSON data: %w", err)
	}

	// Limit contracts to requested count
	if uint64(len(s.deployData.Contracts)) > s.options.TotalContracts {
		s.deployData.Contracts = s.deployData.Contracts[:s.options.TotalContracts]
	} else if uint64(len(s.deployData.Contracts)) < s.options.TotalContracts {
		s.logger.Warnf("Only %d contracts available in data file, requested %d",
			len(s.deployData.Contracts), s.options.TotalContracts)
		s.options.TotalContracts = uint64(len(s.deployData.Contracts))
	}

	return nil
}

func (s *Scenario) compileContract() error {
	// Check if solc is available
	if _, err := exec.LookPath("solc"); err != nil {
		return fmt.Errorf("solc not found in PATH: %w", err)
	}

	contractPath := s.options.ContractFile
	var tempFile string

	// If it's a URL, download to a temp file
	if strings.HasPrefix(contractPath, "http://") || strings.HasPrefix(contractPath, "https://") {
		s.logger.WithField("url", contractPath).Debug("Downloading contract from URL")
		data, err := s.loadDataFromPathOrURL(contractPath)
		if err != nil {
			return fmt.Errorf("failed to download contract file: %w", err)
		}

		// Create temp file
		tmpFile, err := os.CreateTemp("", "contract-*.sol")
		if err != nil {
			return fmt.Errorf("failed to create temp file: %w", err)
		}
		tempFile = tmpFile.Name()
		defer os.Remove(tempFile) // Clean up temp file

		if _, err := tmpFile.Write(data); err != nil {
			tmpFile.Close()
			return fmt.Errorf("failed to write contract to temp file: %w", err)
		}
		tmpFile.Close()

		contractPath = tempFile
	}

	// Compile with the same flags as the Python script to get matching bytecode
	cmd := exec.Command("solc",
		"--bin",
		"--optimize",
		"--optimize-runs", "200",
		"--metadata-hash", "none", // Critical for reproducible bytecode
		contractPath,
	)

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to compile contract: %w", err)
	}

	// Parse the output to extract bytecode
	lines := strings.Split(string(output), "\n")
	inBinarySection := false
	for _, line := range lines {
		if strings.Contains(line, "Binary:") {
			inBinarySection = true
			continue
		}
		if inBinarySection && strings.TrimSpace(line) != "" && !strings.HasPrefix(line, "=") {
			bytecodeHex := strings.TrimSpace(line)
			bytecodeHex = strings.TrimPrefix(bytecodeHex, "0x")
			s.initCode, err = hex.DecodeString(bytecodeHex)
			if err != nil {
				return fmt.Errorf("failed to decode bytecode: %w", err)
			}
			break
		}
	}

	if len(s.initCode) == 0 {
		return fmt.Errorf("failed to extract bytecode from solc output")
	}

	s.logger.Infof("Compiled contract, init code size: %d bytes", len(s.initCode))
	return nil
}

func (s *Scenario) ensureNicksFactory(ctx context.Context) error {
	s.logger.Info("Phase 2: Checking if Nick's factory exists")

	// Get a client to check factory existence
	client := s.walletPool.GetClient(spamoor.WithClientGroup(s.options.ClientGroup))
	if client == nil {
		return scenario.ErrNoClients
	}

	// First check canonical Nick's factory address
	canonicalAddr := common.HexToAddress(NickFactoryAddress)
	code, err := client.GetEthClient().CodeAt(ctx, canonicalAddr, nil)
	if err != nil {
		return fmt.Errorf("failed to check factory existence: %w", err)
	}

	if len(code) > 0 {
		s.logger.Infof("Nick's factory already exists at canonical address %s", NickFactoryAddress)
		s.factoryAddress = canonicalAddr
		return nil
	}

	s.logger.Infof("Nick's factory not found at canonical address, deploying it now")

	// Deploy Nick's factory at the canonical address
	if err := s.deployNicksFactory(ctx); err != nil {
		return fmt.Errorf("failed to deploy Nick's factory: %w", err)
	}

	// Verify deployment
	code, err = client.GetEthClient().CodeAt(ctx, canonicalAddr, nil)
	if err != nil {
		return fmt.Errorf("failed to verify factory deployment: %w", err)
	}

	if len(code) == 0 {
		return fmt.Errorf("factory deployment failed - no code at canonical address")
	}

	s.factoryAddress = canonicalAddr
	s.logger.Infof("Successfully deployed Nick's factory at canonical address %s", NickFactoryAddress)
	return nil
}

func (s *Scenario) deployNicksFactory(ctx context.Context) error {
	// Nick's factory deployment method:
	// 1. Fund the deployment address 0x3fab184622dc19b6109349b94811493bf2a45362
	// 2. Send the pre-signed deployment transaction

	client := s.walletPool.GetClient(spamoor.WithClientGroup(s.options.ClientGroup))
	if client == nil {
		return scenario.ErrNoClients
	}

	// The deployment address that needs to be funded
	deployerAddr := common.HexToAddress("0x3fab184622dc19b6109349b94811493bf2a45362")

	// Check if deployer already has funds
	balance, err := client.GetEthClient().BalanceAt(ctx, deployerAddr, nil)
	if err != nil {
		return fmt.Errorf("failed to check deployer balance: %w", err)
	}

	// Need at least 0.04 ETH for deployment
	requiredBalance := new(big.Int).Mul(big.NewInt(4), big.NewInt(1e16)) // 0.04 ETH

	if balance.Cmp(requiredBalance) < 0 {
		s.logger.Infof("Funding Nick's factory deployer address with 0.05 ETH")

		// Fund the deployer address
		wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, 0)
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

		// Send 0.05 ETH to the deployer address
		fundAmount := new(big.Int).Mul(big.NewInt(5), big.NewInt(1e16)) // 0.05 ETH
		txData, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
			Gas:       21000,
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			To:        &deployerAddr,
			Value:     uint256.MustFromBig(fundAmount),
		})

		if err != nil {
			return fmt.Errorf("failed to build funding tx: %w", err)
		}

		tx, err := wallet.BuildDynamicFeeTx(txData)
		if err != nil {
			return fmt.Errorf("failed to sign funding tx: %w", err)
		}

		receipt, err := s.walletPool.GetTxPool().SendAndAwaitTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
			Client:      client,
			ClientGroup: s.options.ClientGroup,
			Rebroadcast: true,
		})

		if err != nil {
			return fmt.Errorf("failed to fund deployer address: %w", err)
		}

		if receipt == nil || receipt.Status != 1 {
			return fmt.Errorf("funding transaction failed")
		}

		s.logger.Infof("Funded deployer address in block %s", receipt.BlockNumber.String())
	}

	// Now send the pre-signed deployment transaction
	// This is the actual signed transaction from Nick that deploys the factory
	rawTx := "0xf8a58085174876e800830186a08080b853604580600e600039806000f350fe7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe03601600081602082378035828234f58015156039578182fd5b8082525050506014600cf31ba02222222222222222222222222222222222222222222222222222222222222222a02222222222222222222222222222222222222222222222222222222222222222"

	var tx types.Transaction
	if err := tx.UnmarshalBinary(common.FromHex(rawTx)); err != nil {
		return fmt.Errorf("failed to unmarshal pre-signed transaction: %w", err)
	}

	// Send the pre-signed transaction
	if err := client.GetEthClient().SendTransaction(ctx, &tx); err != nil {
		// Check if the error is because factory already exists
		if strings.Contains(err.Error(), "already known") || strings.Contains(err.Error(), "nonce too low") {
			s.logger.Info("Factory deployment transaction already processed")
			return nil
		}
		return fmt.Errorf("failed to send deployment transaction: %w", err)
	}

	// Wait for the transaction to be mined
	receipt, err := bind.WaitMined(ctx, client.GetEthClient(), &tx)
	if err != nil {
		return fmt.Errorf("failed to wait for deployment: %w", err)
	}

	if receipt.Status != 1 {
		return fmt.Errorf("factory deployment transaction failed")
	}

	s.logger.Infof("Nick's factory deployed successfully in block %s", receipt.BlockNumber.String())
	return nil
}

func (s *Scenario) Run(ctx context.Context) error {
	// Phase 1: Fund auxiliary accounts
	if !s.options.SkipFunding {
		s.phase = "funding"
		s.logger.Info("Phase 1: Funding auxiliary accounts")

		// Count total auxiliary accounts
		auxAccounts := make(map[string]bool)
		for _, contract := range s.deployData.Contracts {
			for _, acc := range contract.AuxiliaryAccounts {
				auxAccounts[acc] = true
			}
		}
		totalAuxAccounts := uint64(len(auxAccounts))
		s.logger.Infof("Total unique auxiliary accounts to fund: %d", totalAuxAccounts)

		// Run funding scenario - deploy all as fast as possible
		err := scenario.RunTransactionScenario(ctx, scenario.TransactionScenarioOptions{
			TotalCount:                  totalAuxAccounts,
			Throughput:                  0,   // No throughput limit - send as fast as possible
			MaxPending:                  100, // Reasonable pending limit
			ThroughputIncrementInterval: 0,
			Timeout:                     0,
			WalletPool:                  s.walletPool,
			Logger:                      s.logger,
			ProcessNextTxFn:             s.sendFundingTx,
		})

		if err != nil {
			s.logger.Warnf("Error during funding phase: %v", err)
		}

		s.logger.Infof("Funded %d auxiliary accounts", len(s.fundedAccounts))
	}

	// Phase 2: Ensure Nick's factory exists
	if !s.options.SkipContracts {
		if err := s.ensureNicksFactory(ctx); err != nil {
			return fmt.Errorf("failed to ensure Nick's factory exists: %w", err)
		}
	}

	// Phase 3: Deploy contracts via Nick's factory
	if !s.options.SkipContracts {
		s.phase = "deploying"
		s.currentIndex = 0
		s.logger.Info("Phase 3: Deploying contracts via Nick's factory")

		err := scenario.RunTransactionScenario(ctx, scenario.TransactionScenarioOptions{
			TotalCount:                  s.options.TotalContracts,
			Throughput:                  0,  // No throughput limit - deploy as fast as possible
			MaxPending:                  50, // Lower limit for deployments due to higher gas usage
			ThroughputIncrementInterval: 0,
			Timeout:                     0,
			WalletPool:                  s.walletPool,
			Logger:                      s.logger,
			ProcessNextTxFn:             s.sendDeploymentTx,
		})

		if err != nil {
			s.logger.Warnf("Error during deployment phase: %v", err)
		}

		s.logger.Infof("Deployed %d contracts", len(s.deployedContracts))
	}

	// Save deployment info
	if err := s.saveDeploymentInfo(); err != nil {
		s.logger.Warnf("Failed to save deployment info: %v", err)
	}

	s.logger.WithFields(logrus.Fields{
		"contracts_deployed": len(s.deployedContracts),
		"accounts_funded":    len(s.fundedAccounts),
	}).Info("storage trie brancher scenario completed")

	return nil
}

func (s *Scenario) sendFundingTx(ctx context.Context, params *scenario.ProcessNextTxParams) error {
	// Collect all unique auxiliary accounts with proper locking
	s.mu.RLock()
	auxAccounts := make([]string, 0)
	for _, contract := range s.deployData.Contracts {
		for _, acc := range contract.AuxiliaryAccounts {
			if !s.fundedAccounts[acc] {
				auxAccounts = append(auxAccounts, acc)
			}
		}
	}
	s.mu.RUnlock()

	if params.TxIdx >= uint64(len(auxAccounts)) {
		return fmt.Errorf("funding index out of range")
	}

	account := auxAccounts[params.TxIdx]

	// Get client and wallet
	client := s.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, int(params.TxIdx)),
		spamoor.WithClientGroup(s.options.ClientGroup),
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

	// Build funding transaction
	toAddr := common.HexToAddress(account)
	txData, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
		Gas:       FundGasLimit,
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		To:        &toAddr,
		Value:     uint256.NewInt(FundingAmount),
	})

	if err != nil {
		return fmt.Errorf("failed to build funding tx for %s: %w", account, err)
	}

	tx, err := wallet.BuildDynamicFeeTx(txData)
	if err != nil {
		return fmt.Errorf("failed to sign funding tx for %s: %w", account, err)
	}

	// Send transaction
	receipt, err := s.walletPool.GetTxPool().SendAndAwaitTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		ClientGroup: s.options.ClientGroup,
		Rebroadcast: true, // Always use rebroadcast for reliability
	})

	params.NotifySubmitted()

	logger := s.logger.WithFields(logrus.Fields{
		"rpc":    client.GetName(),
		"nonce":  tx.Nonce(),
		"wallet": s.walletPool.GetWalletName(wallet.GetAddress()),
	})

	params.OrderedLogCb(func() {
		if err != nil {
			logger.Warnf("could not send funding transaction: %v", err)
		} else if s.options.LogTxs {
			logger.Infof("sent funding tx #%6d: %v to %s", params.TxIdx+1, tx.Hash().String(), account)
		} else {
			logger.Debugf("sent funding tx #%6d: %v to %s", params.TxIdx+1, tx.Hash().String(), account)
		}
	})

	if receipt != nil {
		s.mu.Lock()
		s.fundedAccounts[account] = true
		s.mu.Unlock()
	}

	return err
}

func (s *Scenario) sendDeploymentTx(ctx context.Context, params *scenario.ProcessNextTxParams) error {
	if params.TxIdx >= uint64(len(s.deployData.Contracts)) {
		return fmt.Errorf("deployment index out of range")
	}

	contract := s.deployData.Contracts[params.TxIdx]

	// Skip if already deployed
	s.mu.RLock()
	if s.deployedContracts[contract.ContractAddress] {
		s.mu.RUnlock()
		return nil
	}
	s.mu.RUnlock()

	// Get client and wallet
	client := s.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, int(params.TxIdx)),
		spamoor.WithClientGroup(s.options.ClientGroup),
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

	// Convert salt to bytes
	saltBytes, err := s.saltToBytes(contract.Salt)
	if err != nil {
		return fmt.Errorf("failed to convert salt: %w", err)
	}

	// Build factory call data (salt + init code)
	factoryData := append(saltBytes, s.initCode...)

	// Build deployment transaction
	txData, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
		Gas:       DeployGasLimit,
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		To:        &s.factoryAddress,
		Data:      factoryData,
	})

	if err != nil {
		return fmt.Errorf("failed to build deployment tx: %w", err)
	}

	tx, err := wallet.BuildDynamicFeeTx(txData)
	if err != nil {
		return fmt.Errorf("failed to sign deployment tx: %w", err)
	}

	// Send transaction
	receipt, err := s.walletPool.GetTxPool().SendAndAwaitTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		ClientGroup: s.options.ClientGroup,
		Rebroadcast: true, // Always use rebroadcast for reliability
	})

	params.NotifySubmitted()

	logger := s.logger.WithFields(logrus.Fields{
		"rpc":    client.GetName(),
		"nonce":  tx.Nonce(),
		"wallet": s.walletPool.GetWalletName(wallet.GetAddress()),
	})

	params.OrderedLogCb(func() {
		if err != nil {
			logger.Warnf("could not send deployment transaction: %v", err)
		} else if s.options.LogTxs {
			deployAddr := s.calculateCreate2Address(saltBytes)
			logger.Infof("sent deployment tx #%6d: %v (contract: %s)", params.TxIdx+1, tx.Hash().String(), deployAddr.Hex())
		} else {
			logger.Debugf("sent deployment tx #%6d: %v", params.TxIdx+1, tx.Hash().String())
		}
	})

	if receipt != nil && receipt.Status == 1 {
		// Calculate deployment address
		deployAddr := s.calculateCreate2Address(saltBytes)

		s.mu.Lock()
		s.deployedContracts[contract.ContractAddress] = true
		s.deployments = append(s.deployments, DeploymentInfo{
			Address:           deployAddr.Hex(),
			Salt:              "0x" + hex.EncodeToString(saltBytes),
			AuxiliaryAccounts: contract.AuxiliaryAccounts,
		})
		s.mu.Unlock()
	}

	return err
}

func (s *Scenario) saltToBytes(salt interface{}) ([]byte, error) {
	switch v := salt.(type) {
	case float64:
		// JSON numbers come as float64
		saltBig := big.NewInt(int64(v))
		saltBytes := make([]byte, 32)
		saltBig.FillBytes(saltBytes)
		return saltBytes, nil
	case int:
		saltBig := big.NewInt(int64(v))
		saltBytes := make([]byte, 32)
		saltBig.FillBytes(saltBytes)
		return saltBytes, nil
	case string:
		saltStr := v
		saltStr = strings.TrimPrefix(saltStr, "0x")
		saltBytes, err := hex.DecodeString(saltStr)
		if err != nil {
			return nil, err
		}
		if len(saltBytes) != 32 {
			// Pad with zeros if needed
			padded := make([]byte, 32)
			copy(padded[32-len(saltBytes):], saltBytes)
			return padded, nil
		}
		return saltBytes, nil
	default:
		return nil, fmt.Errorf("unsupported salt type: %T", v)
	}
}

func (s *Scenario) calculateCreate2Address(salt []byte) common.Address {
	// CREATE2 address = keccak256(0xff ++ deployer ++ salt ++ keccak256(init_code))[12:]
	initCodeHash := crypto.Keccak256(s.initCode)

	data := []byte{0xff}
	data = append(data, common.HexToAddress(NickFactoryAddress).Bytes()...)
	data = append(data, salt...)
	data = append(data, initCodeHash...)

	hash := crypto.Keccak256(data)
	return common.BytesToAddress(hash[12:])
}

func (s *Scenario) saveDeploymentInfo() error {
	// Save deployment info to JSON file
	deploymentInfo := map[string]interface{}{
		"storage_depth":      s.options.StorageDepth,
		"account_depth":      s.options.AccountDepth,
		"deployer":           NickFactoryAddress,
		"contracts_deployed": len(s.deployedContracts),
		"accounts_funded":    len(s.fundedAccounts),
		"contracts":          s.deployments,
	}

	outputPath := filepath.Join("scenarios", "statebloat", "storage_trie_brancher", "deployed_contracts.json")
	data, err := json.MarshalIndent(deploymentInfo, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return err
	}

	// Also save just the contract addresses to stubs.json
	addresses := make([]string, 0, len(s.deployments))
	for _, deploy := range s.deployments {
		addresses = append(addresses, deploy.Address)
	}

	stubsPath := filepath.Join("scenarios", "statebloat", "storage_trie_brancher", "stubs.json")
	stubsData, err := json.MarshalIndent(addresses, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(stubsPath, stubsData, 0644); err != nil {
		return err
	}

	s.logger.Infof("Deployment info saved to %s", outputPath)
	s.logger.Infof("Contract addresses saved to %s", stubsPath)

	return nil
}
