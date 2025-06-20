package calltx

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/utils"
)

type ScenarioOptions struct {
	TotalCount       uint64 `yaml:"total_count"`
	Throughput       uint64 `yaml:"throughput"`
	MaxPending       uint64 `yaml:"max_pending"`
	MaxWallets       uint64 `yaml:"max_wallets"`
	Rebroadcast      uint64 `yaml:"rebroadcast"`
	BaseFee          uint64 `yaml:"base_fee"`
	TipFee           uint64 `yaml:"tip_fee"`
	DeployGasLimit   uint64 `yaml:"deploy_gas_limit"`
	GasLimit         uint64 `yaml:"gas_limit"`
	Amount           uint64 `yaml:"amount"`
	RandomAmount     bool   `yaml:"random_amount"`
	RandomTarget     bool   `yaml:"random_target"`
	ContractCode     string `yaml:"contract_code"`
	ContractFile     string `yaml:"contract_file"`
	ContractAddress  string `yaml:"contract_address"`
	ContractArgs     string `yaml:"contract_args"`
	ContractAddrPath string `yaml:"contract_addr_path"`
	CallData         string `yaml:"call_data"`
	CallABI          string `yaml:"call_abi"`
	CallABIFile      string `yaml:"call_abi_file"`
	CallFnName       string `yaml:"call_fn_name"`
	CallFnSig        string `yaml:"call_fn_sig"`
	CallArgs         string `yaml:"call_args"`
	Timeout          string `yaml:"timeout"`
	ClientGroup      string `yaml:"client_group"`
	LogTxs           bool   `yaml:"log_txs"`
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	contractAddr   common.Address
	abiCallBuilder *utils.ABICallDataBuilder
}

var ScenarioName = "calltx"
var ScenarioDefaultOptions = ScenarioOptions{
	TotalCount:       0,
	Throughput:       100,
	MaxPending:       0,
	MaxWallets:       0,
	Rebroadcast:      1,
	BaseFee:          20,
	TipFee:           2,
	DeployGasLimit:   2000000,
	GasLimit:         1000000,
	Amount:           0,
	RandomAmount:     false,
	RandomTarget:     false,
	ContractCode:     "",
	ContractFile:     "",
	ContractAddress:  "",
	ContractArgs:     "",
	ContractAddrPath: "",
	CallData:         "",
	CallABI:          "",
	CallABIFile:      "",
	CallFnName:       "",
	CallFnSig:        "",
	CallArgs:         "",
	Timeout:          "",
	ClientGroup:      "",
	LogTxs:           false,
}
var ScenarioDescriptor = scenario.Descriptor{
	Name:           ScenarioName,
	Description:    "Deploy a contract and repeatedly call a function on it",
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
	flags.Uint64VarP(&s.options.TotalCount, "count", "c", ScenarioDefaultOptions.TotalCount, "Total number of call transactions to send")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t", ScenarioDefaultOptions.Throughput, "Number of call transactions to send per slot")
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", ScenarioDefaultOptions.Rebroadcast, "Enable reliable rebroadcast with unlimited retries and exponential backoff")
	flags.Uint64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use in call and deployment transactions (in gwei)")
	flags.Uint64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use in call and deployment transactions (in gwei)")
	flags.Uint64Var(&s.options.DeployGasLimit, "deploy-gas-limit", ScenarioDefaultOptions.DeployGasLimit, "Gas limit to use for deployment transaction")
	flags.Uint64Var(&s.options.GasLimit, "gas-limit", ScenarioDefaultOptions.GasLimit, "Gas limit to use for call transactions")
	flags.Uint64Var(&s.options.Amount, "amount", ScenarioDefaultOptions.Amount, "Transfer amount per transaction (in gwei)")
	flags.BoolVar(&s.options.RandomAmount, "random-amount", ScenarioDefaultOptions.RandomAmount, "Use random amounts for transactions (with --amount as limit)")
	flags.BoolVar(&s.options.RandomTarget, "random-target", ScenarioDefaultOptions.RandomTarget, "Use random to addresses for transactions")
	flags.StringVar(&s.options.ContractCode, "contract-code", ScenarioDefaultOptions.ContractCode, "Contract code to deploy")
	flags.StringVar(&s.options.ContractFile, "contract-file", ScenarioDefaultOptions.ContractFile, "Contract file to deploy")
	flags.StringVar(&s.options.ContractAddress, "contract-address", ScenarioDefaultOptions.ContractAddress, "Address of already deployed contract (skips deployment)")
	flags.StringVar(&s.options.ContractArgs, "contract-args", ScenarioDefaultOptions.ContractArgs, "Contract arguments to pass to the constructor")
	flags.StringVar(&s.options.ContractAddrPath, "contract-addr-path", ScenarioDefaultOptions.ContractAddrPath, "Path to child contract created during deployment (e.g. '.0.1' for nonce 1 of nonce 0)")
	flags.StringVar(&s.options.CallData, "call-data", ScenarioDefaultOptions.CallData, "Data to pass to the function to call")
	flags.StringVar(&s.options.CallABI, "call-abi", ScenarioDefaultOptions.CallABI, "JSON ABI of the contract for function calls")
	flags.StringVar(&s.options.CallABIFile, "call-abi-file", ScenarioDefaultOptions.CallABIFile, "JSON ABI file of the contract for function calls")
	flags.StringVar(&s.options.CallFnName, "call-fn-name", ScenarioDefaultOptions.CallFnName, "Function name to call (requires --call-abi)")
	flags.StringVar(&s.options.CallFnSig, "call-fn-sig", ScenarioDefaultOptions.CallFnSig, "Function signature to call (alternative to --call-abi)")
	flags.StringVar(&s.options.CallArgs, "call-args", ScenarioDefaultOptions.CallArgs, "JSON array of arguments to pass to the function")
	flags.StringVar(&s.options.ClientGroup, "client-group", ScenarioDefaultOptions.ClientGroup, "Client group to use for sending transactions")
	flags.StringVar(&s.options.Timeout, "timeout", ScenarioDefaultOptions.Timeout, "Timeout for the scenario (e.g. '1h', '30m', '5s') - empty means no timeout")
	flags.BoolVar(&s.options.LogTxs, "log-txs", ScenarioDefaultOptions.LogTxs, "Log all submitted transactions")
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

	if s.options.MaxWallets > 0 {
		s.walletPool.SetWalletCount(s.options.MaxWallets)
	} else if s.options.TotalCount > 0 {
		maxWallets := s.options.TotalCount / 50
		if maxWallets < 10 {
			maxWallets = 10
		} else if maxWallets > 1000 {
			maxWallets = 1000
		}

		s.walletPool.SetWalletCount(maxWallets)
	} else {
		if s.options.Throughput*10 < 1000 {
			s.walletPool.SetWalletCount(s.options.Throughput * 10)
		} else {
			s.walletPool.SetWalletCount(1000)
		}
	}

	if s.options.TotalCount == 0 && s.options.Throughput == 0 {
		return fmt.Errorf("neither total count nor throughput limit set, must define at least one of them (see --help for list of all flags)")
	}

	// Validate contract source options (mutually exclusive)
	contractSources := 0
	if s.options.ContractCode != "" {
		contractSources++
	}
	if s.options.ContractFile != "" {
		contractSources++
	}
	if s.options.ContractAddress != "" {
		contractSources++
	}

	if contractSources == 0 {
		return fmt.Errorf("must specify one of: --contract-code, --contract-file, or --contract-address")
	}
	if contractSources > 1 {
		return fmt.Errorf("only one of --contract-code, --contract-file, or --contract-address can be set")
	}

	// Initialize ABI call builder if ABI options are provided
	if s.options.CallABI != "" || s.options.CallABIFile != "" || s.options.CallFnSig != "" {
		var abiContent string
		var err error

		// Load ABI content from file if specified
		if s.options.CallABIFile != "" {
			if s.options.CallABI != "" {
				return fmt.Errorf("only one of --call-abi or --call-abi-file can be set")
			}

			var abiBytes []byte
			if strings.HasPrefix(s.options.CallABIFile, "https://") || strings.HasPrefix(s.options.CallABIFile, "http://") {
				resp, err := http.Get(s.options.CallABIFile)
				if err != nil {
					return fmt.Errorf("could not load ABI file: %w", err)
				}
				defer resp.Body.Close()
				abiBytes, err = io.ReadAll(resp.Body)
				if err != nil {
					return fmt.Errorf("could not read ABI file: %w", err)
				}
			} else {
				abiBytes, err = os.ReadFile(s.options.CallABIFile)
				if err != nil {
					return fmt.Errorf("could not read ABI file: %w", err)
				}
			}
			abiContent = string(abiBytes)
		} else {
			abiContent = s.options.CallABI
		}

		callArgsWithPlaceholders := s.replaceCallDataPlaceholders(s.options.CallArgs)
		s.abiCallBuilder, err = utils.NewABICallDataBuilder(abiContent, s.options.CallFnName, s.options.CallFnSig, callArgsWithPlaceholders)
		if err != nil {
			return fmt.Errorf("failed to initialize ABI call builder: %w", err)
		}
	}

	return nil
}

func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Infof("starting scenario: %s", ScenarioName)
	defer s.logger.Infof("scenario %s finished.", ScenarioName)

	if s.options.ContractAddress != "" {
		// Existing contract
		if !common.IsHexAddress(s.options.ContractAddress) {
			return fmt.Errorf("invalid contract address format: %s", s.options.ContractAddress)
		}
		s.contractAddr = common.HexToAddress(s.options.ContractAddress)
		s.logger.Infof("using existing contract: %v", s.contractAddr.String())
	} else {
		// New contract
		var contractCode []byte
		if s.options.ContractCode != "" {
			contractCode = common.FromHex(s.options.ContractCode)
		} else if strings.HasPrefix(s.options.ContractFile, "https://") || strings.HasPrefix(s.options.ContractFile, "http://") {
			resp, err := http.Get(s.options.ContractFile)
			if err != nil {
				return fmt.Errorf("could not load contract file: %w", err)
			}
			defer resp.Body.Close()
			contractCodeHex, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("could not read contract file: %w", err)
			}
			contractCode = common.FromHex(strings.Trim(string(contractCodeHex), "\r\n\t "))
		} else {
			code, err := os.ReadFile(s.options.ContractFile)
			if err != nil {
				return fmt.Errorf("could not read contract file: %w", err)
			}
			contractCode = common.FromHex(strings.Trim(string(code), "\r\n\t "))
		}

		// deploy contract
		contractReceipt, _, err := s.sendDeploymentTx(ctx, contractCode)
		if err != nil {
			s.logger.Errorf("could not deploy contract: %v", err)
			return err
		}
		if contractReceipt == nil {
			return fmt.Errorf("could not deploy contract: deployment failed")
		}
		s.contractAddr = contractReceipt.ContractAddress
		s.logger.Infof("deployed contract: %v (confirmed in block #%v)", s.contractAddr.String(), contractReceipt.BlockNumber.String())
	}

	// Calculate child contract address if path is specified
	if s.options.ContractAddrPath != "" {
		childAddr, err := s.calculateChildContractAddress(s.contractAddr, s.options.ContractAddrPath)
		if err != nil {
			return fmt.Errorf("failed to calculate child contract address: %w", err)
		}
		s.logger.Infof("targeting child contract at path %s: %v", s.options.ContractAddrPath, childAddr.String())
		s.contractAddr = childAddr
	}

	// send transactions
	maxPending := s.options.MaxPending
	if maxPending == 0 {
		maxPending = s.options.Throughput * 10
		if maxPending == 0 {
			maxPending = 4000
		}

		if maxPending > s.walletPool.GetConfiguredWalletCount()*10 {
			maxPending = s.walletPool.GetConfiguredWalletCount() * 10
		}
	}

	// Parse timeout
	var timeout time.Duration
	if s.options.Timeout != "" {
		var err error
		timeout, err = time.ParseDuration(s.options.Timeout)
		if err != nil {
			return fmt.Errorf("invalid timeout value: %v", err)
		}
		s.logger.Infof("Timeout set to %v", timeout)
	}

	err := scenario.RunTransactionScenario(ctx, scenario.TransactionScenarioOptions{
		TotalCount:                  s.options.TotalCount,
		Throughput:                  s.options.Throughput,
		MaxPending:                  maxPending,
		ThroughputIncrementInterval: 0,
		Timeout:                     timeout,
		WalletPool:                  s.walletPool,

		Logger: s.logger,
		ProcessNextTxFn: func(ctx context.Context, txIdx uint64, onComplete func()) (func(), error) {
			logger := s.logger
			tx, client, wallet, err := s.sendTx(ctx, txIdx, onComplete)
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

func (s *Scenario) sendDeploymentTx(ctx context.Context, contractCode []byte) (*types.Receipt, *spamoor.Client, error) {
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, s.options.ClientGroup)
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, 0)

	var feeCap *big.Int
	var tipCap *big.Int

	if client == nil {
		return nil, client, fmt.Errorf("no client available")
	}

	feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(client, s.options.BaseFee, s.options.TipFee)
	if err != nil {
		return nil, client, err
	}

	deployData := contractCode
	if s.options.ContractArgs != "" {
		deployData = append(deployData, common.FromHex(s.options.ContractArgs)...)
	}

	txData, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       s.options.DeployGasLimit,
		To:        nil,
		Value:     uint256.NewInt(0),
		Data:      deployData,
	})
	if err != nil {
		return nil, nil, err
	}

	tx, err := wallet.BuildDynamicFeeTx(txData)
	if err != nil {
		return nil, nil, err
	}

	receipt, err := s.walletPool.GetTxPool().SendAndAwaitTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		Rebroadcast: true,
	})
	if err != nil {
		return nil, client, err
	}

	return receipt, client, nil
}

func (s *Scenario) sendTx(ctx context.Context, txIdx uint64, onComplete func()) (*types.Transaction, *spamoor.Client, *spamoor.Wallet, error) {
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, int(txIdx), s.options.ClientGroup)
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(txIdx))
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

	// Determine gas limit: use block gas limit if GasLimit is 0
	gasLimit := s.options.GasLimit
	if gasLimit == 0 {
		var err error
		gasLimit, err = s.walletPool.GetTxPool().GetCurrentGasLimitWithInit()
		if err != nil {
			s.logger.Warnf("tx %6d: failed to fetch current gas limit: %v, using fallback", txIdx+1, err)
			gasLimit = 30000000
		} else if gasLimit == 0 {
			// Final fallback to a reasonable default if no block gas limit is available
			gasLimit = 30000000
			s.logger.Warnf("tx %6d: no gas limit available, using fallback %v", txIdx+1, gasLimit)
		} else {
			s.logger.Debugf("tx %6d: using block gas limit %v", txIdx+1, gasLimit)
		}
	}

	amount := uint256.NewInt(s.options.Amount)
	amount = amount.Mul(amount, uint256.NewInt(1000000000))
	if s.options.RandomAmount {
		n, err := rand.Int(rand.Reader, amount.ToBig())
		if err == nil {
			amount = uint256.MustFromBig(n)
		}
	}

	txCallData := []byte{}

	if s.abiCallBuilder != nil {
		// Use ABI call builder
		var err error
		txCallData, err = s.abiCallBuilder.BuildCallData(txIdx)
		if err != nil {
			return nil, nil, wallet, fmt.Errorf("failed to build ABI call data: %w", err)
		}
	} else if s.options.CallData != "" {
		// Use raw call data with placeholder replacement
		callDataWithPlaceholders := s.replaceCallDataPlaceholders(s.options.CallData)
		dataBytes, err := txbuilder.ParseBlobRefsBytes(strings.Split(callDataWithPlaceholders, ","), nil)
		if err != nil {
			return nil, nil, wallet, err
		}
		txCallData = dataBytes
	}

	txData, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       gasLimit,
		To:        &s.contractAddr,
		Value:     amount,
		Data:      txCallData,
	})
	if err != nil {
		return nil, nil, wallet, err
	}

	tx, err := wallet.BuildDynamicFeeTx(txData)
	if err != nil {
		return nil, nil, wallet, err
	}

	transactionSubmitted = true
	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		Rebroadcast: s.options.Rebroadcast > 0,
		OnComplete: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			onComplete()
		},
		OnConfirm: func(tx *types.Transaction, receipt *types.Receipt) {
			txFees := utils.GetTransactionFees(tx, receipt)
			s.logger.WithField("rpc", client.GetName()).Debugf(
				" transaction %d confirmed in block #%v. total fee: %v gwei (base: %v) logs: %v",
				txIdx+1,
				receipt.BlockNumber.String(),
				txFees.TotalFeeGweiString(),
				txFees.TxBaseFeeGweiString(),
				len(receipt.Logs),
			)
		},
		LogFn: spamoor.GetDefaultLogFn(s.logger, "", fmt.Sprintf("%6d", txIdx+1), tx),
	})
	if err != nil {
		// reset nonce if tx was not sent
		wallet.ResetPendingNonce(ctx, client)

		return nil, client, wallet, err
	}

	return tx, client, wallet, nil
}

// replaceCallDataPlaceholders replaces placeholders in call data with actual values
func (s *Scenario) replaceCallDataPlaceholders(callData string) string {
	// Replace factory address placeholder with well-known CREATE2 factory address
	if strings.Contains(callData, "{factory_address}") {
		// Use the same well-known CREATE2 factory address as factorydeploytx scenario
		// Get the very well known factory deployer wallet address and calculate the factory address
		factoryWalletAddr := s.walletPool.GetVeryWellKnownWalletAddress("create2-factory-deployer")
		factoryAddr := crypto.CreateAddress(factoryWalletAddr, 0)
		result := strings.ReplaceAll(callData, "{factory_address}", factoryAddr.Hex())
		return result
	}

	return callData
}

// calculateChildContractAddress calculates the address of a child contract based on the nonce path
func (s *Scenario) calculateChildContractAddress(parentAddr common.Address, noncePath string) (common.Address, error) {
	if noncePath == "" || !strings.HasPrefix(noncePath, ".") {
		return common.Address{}, fmt.Errorf("invalid child contract path format, must start with '.' (e.g. '.0.1')")
	}

	// Remove the leading dot and split by dots
	pathStr := strings.TrimPrefix(noncePath, ".")
	if pathStr == "" {
		return common.Address{}, fmt.Errorf("empty child contract path")
	}

	nonceParts := strings.Split(pathStr, ".")
	currentAddr := parentAddr

	for _, nonceStr := range nonceParts {
		nonce, err := strconv.ParseUint(nonceStr, 10, 64)
		if err != nil {
			return common.Address{}, fmt.Errorf("invalid nonce value '%s' in path: %w", nonceStr, err)
		}

		// Calculate the child contract address using CREATE opcode formula
		currentAddr = crypto.CreateAddress(currentAddr, nonce)
	}

	return currentAddr, nil
}
