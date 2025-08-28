package warmextcodesize

import (
	"context"
	"crypto/ecdsa"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/scenarios/statebloat/warm-extcodesize/contract"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

const (
	deploymentFile             = "deployments.json"
	estimateGasPerTx           = 30000
	estimateExecGasPerAddr     = 3650
	estimateCallDataGasPerAddr = 1003
)

// TODO(weiihann)
type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	contractAddr common.Address
	contract     *contract.Contract

	curAddrIdx int
}

type ScenarioOptions struct {
	MaxWallets      uint64  `yaml:"max_wallets"`
	BaseFee         float64 `yaml:"base_fee"`
	TipFee          float64 `yaml:"tip_fee"`
	ClientGroup     string  `yaml:"client_group"`
	MaxTransactions uint64  `yaml:"max_transactions"`
	Timeout         string  `yaml:"timeout"`
	Throughput      uint64  `yaml:"throughput"`
	GasLimit        uint64  `yaml:"gas_limit"`
	DeployGasLimit  uint64  `yaml:"deploy_gas_limit"`
	ContractAddress string  `yaml:"contract_address"`
	LogTxs          bool    `yaml:"log_txs"`
}

// BlockDeploymentStats tracks deployment statistics per block
type BlockDeploymentStats struct {
	BlockNumber       uint64
	ContractCount     int
	TotalGasUsed      uint64
	TotalBytecodeSize int
}

type PendingTransaction struct {
	TxHash     common.Hash
	PrivateKey *ecdsa.PrivateKey
	Timestamp  time.Time
}

//go:embed contract/BalanceThenExtCodeSize.bin
var contractBin []byte

var (
	ScenarioName           = "warm-extcodesize"
	ScenarioDefaultOptions = ScenarioOptions{
		MaxWallets:      0,
		MaxTransactions: 0,
		ClientGroup:     "",
		GasLimit:        0,
		DeployGasLimit:  20000000,
		ContractAddress: "",
		LogTxs:          false,
	}
	ScenarioDescriptor = scenario.Descriptor{
		Name:           ScenarioName,
		Description:    "Warm up the extcodesize",
		DefaultOptions: ScenarioDefaultOptions,
		NewScenario:    newScenario,
	}
)

func newScenario(logger logrus.FieldLogger) scenario.Scenario {
	return &Scenario{
		logger: logger.WithField("scenario", ScenarioName),
	}
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.MaxTransactions, "max-transactions", ScenarioDefaultOptions.MaxTransactions, "Maximum number of transactions to send (0 = use rate limiting based on block gas limit)")
	flags.StringVar(&s.options.ClientGroup, "client-group", ScenarioDefaultOptions.ClientGroup, "Client group to use for sending transactions")
	flags.Uint64Var(&s.options.GasLimit, "gas-limit", ScenarioDefaultOptions.GasLimit, "Gas limit for each transaction")
	flags.Uint64Var(&s.options.DeployGasLimit, "deploy-gas-limit", ScenarioDefaultOptions.DeployGasLimit, "Gas limit for deployment transaction")
	flags.StringVar(&s.options.ContractAddress, "contract-address", ScenarioDefaultOptions.ContractAddress, "Contract address to use")
	flags.BoolVar(&s.options.LogTxs, "log-txs", ScenarioDefaultOptions.LogTxs, "Log transactions")
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

	if s.options.MaxWallets > 0 {
		s.walletPool.SetWalletCount(s.options.MaxWallets)
	} else {
		// Use only root wallet by default for better efficiency
		// This avoids child wallet funding overhead
		s.walletPool.SetWalletCount(0)
	}

	return nil
}

func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Infof("starting scenario: %s", ScenarioName)
	defer s.logger.Infof("scenario %s finished.", ScenarioName)

	ethClient := s.walletPool.GetClient(spamoor.SelectClientRoundRobin, 0, "").GetEthClient()
	if ethClient == nil {
		return fmt.Errorf("failed to get eth client")
	}

	deployments, err := loadDeployments()
	if err != nil {
		return fmt.Errorf("failed to load deployments: %w", err)
	}

	// Deploy contract if not already deployed
	if s.options.ContractAddress != "" {
		// Existing contract
		if !common.IsHexAddress(s.options.ContractAddress) {
			return fmt.Errorf("invalid contract address format: %s", s.options.ContractAddress)
		}
		s.contractAddr = common.HexToAddress(s.options.ContractAddress)
		s.logger.Infof("using existing contract: %v", s.contractAddr.String())
	} else {
		addr, err := s.sendDeploymentTx()
		if err != nil {
			return fmt.Errorf("failed to deploy contract: %w", err)
		}
		s.contractAddr = addr
		s.logger.Infof("deployed new contract: %v", addr.String())

	}
	// Parse timeout duration
	var timeout time.Duration
	if s.options.Timeout != "" {
		var err error
		timeout, err = time.ParseDuration(s.options.Timeout)
		if err != nil {
			return fmt.Errorf("invalid timeout format '%s': %w", s.options.Timeout, err)
		}
	}

	ct, err := contract.NewContract(s.contractAddr, ethClient)
	if err != nil {
		return fmt.Errorf("failed to create contract: %w", err)
	}
	s.contract = ct

	err = scenario.RunTransactionScenario(ctx, scenario.TransactionScenarioOptions{
		TotalCount:                  s.options.MaxTransactions,
		Throughput:                  s.options.Throughput,
		MaxPending:                  1, // Sequential tx execution
		ThroughputIncrementInterval: 0,
		Timeout:                     timeout,
		WalletPool:                  s.walletPool,
		Logger:                      s.logger,
		ProcessNextTxFn: func(ctx context.Context, txIdx uint64, onComplete func()) (func(), error) {
			logger := s.logger

			addrs, err := s.getPokeAddrs(deployments)
			if err != nil {
				return nil, err
			}

			tx, client, wallet, err := s.sendTx(ctx, txIdx, onComplete, addrs)
			if client != nil {
				logger = logger.WithField("rpc", client.GetName())
			}
			if tx != nil {
				logger = logger.WithField("nonce", tx.Nonce())
			}
			if wallet != nil {
				logger = logger.WithField("wallet", s.walletPool.GetWalletName(wallet.GetAddress()))
			}

			return func() {
				if err != nil {
					logger.Warnf("could not send transaction: %v", err)
				} else if s.options.LogTxs {
					logger.Infof("sent tx #%6d: %v", txIdx+1, tx.Hash().String())
				} else {
					logger.Debugf("sent tx #%6d: %v", txIdx+1, tx.Hash().String())
				}
			}, err
		},
	})

	return err
}

func (s *Scenario) sendDeploymentTx() (common.Address, error) {
	client := s.walletPool.GetClient(spamoor.SelectClientRandom, 0, s.options.ClientGroup)
	wallet := s.walletPool.GetRootWallet().GetWallet()

	var feeCap *big.Int
	var tipCap *big.Int

	if client == nil {
		return common.Address{}, errors.New("no client available")
	}

	if wallet == nil {
		return common.Address{}, errors.New("no wallet available")
	}

	feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(client, float64(s.options.DeployGasLimit), s.options.TipFee)
	if err != nil {
		return common.Address{}, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(wallet.GetPrivateKey(), wallet.GetChainId())
	if err != nil {
		return common.Address{}, err
	}
	auth.GasFeeCap = feeCap
	auth.GasTipCap = tipCap
	auth.GasLimit = s.options.DeployGasLimit

	a, _, _, err := contract.DeployContract(auth, client.GetEthClient())
	if err != nil {
		return common.Address{}, err
	}

	return a, nil
}

func (s *Scenario) sendTx(ctx context.Context, _ uint64, onComplete func(), addrs []common.Address) (*types.Transaction, *spamoor.Client, *spamoor.Wallet, error) {
	client := s.walletPool.GetClient(spamoor.SelectClientRandom, 0, s.options.ClientGroup)
	wallet := s.walletPool.GetRootWallet().GetWallet()
	transactionSubmitted := false

	defer func() {
		if !transactionSubmitted {
			onComplete()
		}
	}()

	if client == nil {
		return nil, client, wallet, fmt.Errorf("no client available")
	}

	feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(client, s.options.BaseFee, s.options.TipFee)
	if err != nil {
		return nil, client, wallet, err
	}

	tx, err := wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       s.options.GasLimit,
		Value:     uint256.NewInt(0),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		return s.contract.PokeBatch(transactOpts, addrs)
	})
	if err != nil {
		return nil, client, wallet, err
	}

	transactionSubmitted = true
	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		Rebroadcast: true,
		OnComplete: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			onComplete() // CRITICAL: Signal completion for scenario counting
		},
		OnConfirm: func(tx *types.Transaction, receipt *types.Receipt) {
			if receipt.Status == types.ReceiptStatusSuccessful {
				s.logger.WithFields(logrus.Fields{
					"txHash":   tx.Hash().Hex(),
					"block":    receipt.BlockNumber.Uint64(),
					"gasUsed":  receipt.GasUsed,
					"contract": s.contractAddr.Hex(),
				}).Infof("contract interaction confirmed")
			}
		},
	})
	if err != nil {
		// Reset nonce if transaction was not sent
		wallet.ResetPendingNonce(ctx, client)
		return nil, client, wallet, fmt.Errorf("failed to send transaction: %w", err)
	}

	return tx, client, wallet, nil
}

type ContractDeployment struct {
	ContractAddress common.Address `json:"contract_address"`
}

func loadDeployments() ([]common.Address, error) {
	// Check if file exists
	if _, err := os.Stat(deploymentFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("deployments.json file does not exist")
	}

	deploymentsFile, err := os.Open(deploymentFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open deployments.json file: %w", err)
	}
	defer deploymentsFile.Close()

	decoder := json.NewDecoder(deploymentsFile)

	// The actual format is: map[walletAddress][]contractAddress
	var deploymentsByWallet map[string][]string
	err = decoder.Decode(&deploymentsByWallet)
	if err != nil {
		return nil, fmt.Errorf("failed to decode deployments.json: %w", err)
	}

	// Flatten all contract addresses from all wallets
	var allAddresses []common.Address
	for walletAddr, contractAddrs := range deploymentsByWallet {
		for _, addrStr := range contractAddrs {
			if !common.IsHexAddress(addrStr) {
				return nil, fmt.Errorf("invalid contract address format: %s (from wallet %s)", addrStr, walletAddr)
			}
			allAddresses = append(allAddresses, common.HexToAddress(addrStr))
		}
	}

	if len(allAddresses) == 0 {
		return nil, fmt.Errorf("no contract addresses found in deployments.json")
	}

	return allAddresses, nil
}

func (s *Scenario) getGasLimit() (uint64, error) {
	gasLimit := s.options.GasLimit
	if gasLimit == 0 {
		var err error
		gasLimit, err = s.walletPool.GetTxPool().GetCurrentGasLimitWithInit()
		if err != nil {
			s.logger.Warnf("failed to get current gas limit: %v, using fallback", err)
			gasLimit = 45000000
		} else if gasLimit == 0 {
			// Final fallback to a reasonable default if no block gas limit is available
			gasLimit = 45000000
			s.logger.Warnf("no gas limit available, using fallback %v", gasLimit)
		} else {
			s.logger.Debugf("using block gas limit %v", gasLimit)
		}
	}
	return gasLimit, nil
}

func (s *Scenario) getPokeAddrs(deployments []common.Address) ([]common.Address, error) {
	// Handle empty deployments
	if len(deployments) == 0 {
		s.logger.Warnf("no deployments available, cannot poke any addresses")
		return nil, fmt.Errorf("no deployments available")
	}

	gasLimit, err := s.getGasLimit()
	if err != nil {
		return nil, err
	}

	numAddrs := (gasLimit - estimateGasPerTx) / (estimateExecGasPerAddr + estimateCallDataGasPerAddr)
	if numAddrs == 0 {
		s.logger.Errorf("gas limit too low to estimate number of addresses")
		return nil, fmt.Errorf("gas limit %d too low for any addresses", gasLimit)
	}

	// Ensure we don't request more addresses than available
	if int(numAddrs) > len(deployments) {
		numAddrs = uint64(len(deployments))
		s.logger.Debugf("limiting addresses to available deployments: %d", numAddrs)
	}

	// Ensure we have at least one address if deployments are available
	if numAddrs == 0 && len(deployments) > 0 {
		numAddrs = 1
		s.logger.Debugf("using minimum of 1 address despite gas limit constraints")
	}

	// Handle wraparound when we reach the end of deployments
	if s.curAddrIdx >= len(deployments) {
		s.curAddrIdx = 0
		s.logger.Debugf("wrapped around to beginning of deployments")
	}

	// Calculate the end index and handle wraparound
	endIdx := s.curAddrIdx + int(numAddrs)
	var result []common.Address

	if endIdx <= len(deployments) {
		// Simple case: no wraparound needed
		result = make([]common.Address, numAddrs)
		copy(result, deployments[s.curAddrIdx:endIdx])
		s.curAddrIdx = endIdx % len(deployments)
	} else {
		// Wraparound case: take from current index to end, then from beginning
		result = make([]common.Address, numAddrs)

		// Copy from current index to end of array
		firstPart := len(deployments) - s.curAddrIdx
		copy(result[:firstPart], deployments[s.curAddrIdx:])

		// Copy remaining from beginning of array
		remaining := int(numAddrs) - firstPart
		copy(result[firstPart:], deployments[:remaining])

		s.curAddrIdx = remaining
	}

	s.logger.Debugf("selected %d addresses starting from index %d", len(result), s.curAddrIdx-len(result))
	return result, nil
}
