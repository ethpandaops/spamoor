package contractdeploy

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenarios/statebloat/contract_deploy/contract"
	"github.com/ethpandaops/spamoor/scenariotypes"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
)

const (
	// GAS_PER_CONTRACT is the estimated gas cost for deploying one StateBloatToken contract
	GAS_PER_CONTRACT = 4970000 // ~4.97M gas per contract
)

type ScenarioOptions struct {
	MaxPending      uint64 `yaml:"max_pending"`
	MaxWallets      uint64 `yaml:"max_wallets"`
	Rebroadcast     uint64 `yaml:"rebroadcast"`
	BaseFee         uint64 `yaml:"base_fee"`
	TipFee          uint64 `yaml:"tip_fee"`
	ClientGroup     string `yaml:"client_group"`
	MaxTransactions uint64 `yaml:"max_transactions"` // Maximum number of transactions to send (0 for unlimited)
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	pendingChan   chan bool
	pendingWGroup sync.WaitGroup

	deployedAddresses []string
	addressesMutex    sync.Mutex
}

var ScenarioName = "contract-deploy"
var ScenarioDefaultOptions = ScenarioOptions{
	MaxPending:      0,
	MaxWallets:      0,
	Rebroadcast:     1,
	BaseFee:         20,
	TipFee:          2,
	ClientGroup:     "default",
	MaxTransactions: 0, // Default to unlimited
}
var ScenarioDescriptor = scenariotypes.ScenarioDescriptor{
	Name:           ScenarioName,
	Description:    "Deploy contracts to create state bloat",
	DefaultOptions: ScenarioDefaultOptions,
	NewScenario:    newScenario,
}

func newScenario(logger logrus.FieldLogger) scenariotypes.Scenario {
	return &Scenario{
		logger: logger.WithField("scenario", ScenarioName),
	}
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", ScenarioDefaultOptions.Rebroadcast, "Number of seconds to wait before re-broadcasting a transaction")
	flags.Uint64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use in transactions (in gwei)")
	flags.Uint64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use in transactions (in gwei)")
	flags.StringVar(&s.options.ClientGroup, "client-group", ScenarioDefaultOptions.ClientGroup, "Client group to use for sending transactions")
	flags.Uint64Var(&s.options.MaxTransactions, "max-transactions", ScenarioDefaultOptions.MaxTransactions, "Maximum number of transactions to send (0 for unlimited)")
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

	if s.options.MaxWallets > 0 {
		walletPool.SetWalletCount(s.options.MaxWallets)
	} else {
		walletPool.SetWalletCount(1000)
	}

	if s.options.MaxPending > 0 {
		s.pendingChan = make(chan bool, s.options.MaxPending)
	}

	return nil
}

func (s *Scenario) Config() string {
	yamlBytes, _ := yaml.Marshal(&s.options)
	return string(yamlBytes)
}

func (s *Scenario) Run(ctx context.Context) error {
	txIdxCounter := uint64(0)
	pendingCount := atomic.Int64{}
	txCount := atomic.Uint64{}
	var lastChan chan bool

	s.logger.Infof("starting scenario: %s", ScenarioName)
	defer s.logger.Infof("scenario %s finished.", ScenarioName)

	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, s.options.ClientGroup)
	if client == nil {
		return fmt.Errorf("no client available")
	}

	var deployTxHashesMu sync.Mutex
	var deployTxHashes []common.Hash

	for {
		if s.options.MaxTransactions > 0 && txIdxCounter >= s.options.MaxTransactions {
			s.logger.Infof("reached maximum number of transactions (%d)", s.options.MaxTransactions)
			break
		}

		txIdx := txIdxCounter
		txIdxCounter++

		if s.pendingChan != nil {
			// await pending transactions
			s.pendingChan <- true
		}
		pendingCount.Add(1)
		currentChan := make(chan bool, 1)

		go func(txIdx uint64, lastChan, currentChan chan bool) {
			defer func() {
				pendingCount.Add(-1)
				currentChan <- true
			}()

			logger := s.logger
			tx, client, wallet, err := s.sendTx(ctx, txIdx, func() {
				if s.pendingChan != nil {
					time.Sleep(100 * time.Millisecond)
					<-s.pendingChan
				}
			})
			if client != nil {
				logger = logger.WithField("rpc", client.GetName())
			}
			if tx != nil {
				logger = logger.WithField("nonce", tx.Nonce())
			}
			if wallet != nil {
				logger = logger.WithField("wallet", s.walletPool.GetWalletName(wallet.GetAddress()))
			}
			if lastChan != nil {
				<-lastChan
			}
			if err != nil {
				logger.Warnf("could not send transaction: %v", err)
				return
			}

			// Collect tx hash for later receipt lookup
			deployTxHashesMu.Lock()
			deployTxHashes = append(deployTxHashes, tx.Hash())
			deployTxHashesMu.Unlock()

			txCount.Add(1)
		}(txIdx, lastChan, currentChan)

		lastChan = currentChan
	}
	s.pendingWGroup.Wait()

	s.logger.WithFields(logrus.Fields{
		"test":        "contract-deploy",
		"total_txs":   txCount.Load(),
		"pending_txs": pendingCount.Load(),
	}).Info("finished sending transactions, awaiting block inclusion...")

	// Wait for a new block to be mined
	ethClient := client.GetEthClient()
	startBlock, err := ethClient.BlockByNumber(ctx, nil)
	if err != nil {
		s.logger.Warnf("Failed to get current block: %v", err)
		return err
	}
	startBlockNum := startBlock.NumberU64()
	var latestBlockNum uint64
	for {
		block, err := ethClient.BlockByNumber(ctx, nil)
		if err != nil {
			s.logger.Warnf("Failed to get block: %v", err)
			return err
		}
		latestBlockNum = block.NumberU64()
		if latestBlockNum > startBlockNum {
			break
		}
		time.Sleep(2 * time.Second)
	}

	// For each deployment tx, get the receipt and log contract address and gas used
	var deployedAddresses []string
	for _, txHash := range deployTxHashes {
		receipt, err := ethClient.TransactionReceipt(ctx, txHash)
		if err != nil {
			s.logger.Warnf("Failed to get receipt for tx %s: %v", txHash.Hex(), err)
			continue
		}
		if receipt.ContractAddress != (common.Address{}) {
			addr := receipt.ContractAddress.Hex()
			deployedAddresses = append(deployedAddresses, addr)
			s.logger.Infof("Deployed contract at address: %s (gas used: %d)", addr, receipt.GasUsed)
		}
	}

	// Write deployed addresses to JSON file and log them
	if len(deployedAddresses) > 0 {
		file, err := os.Create("deployed_contracts.json")
		if err == nil {
			_ = json.NewEncoder(file).Encode(deployedAddresses)
			file.Close()
			s.logger.Infof("Wrote %d deployed contract addresses to deployed_contracts.json", len(deployedAddresses))
		} else {
			s.logger.Warnf("Failed to write deployed_contracts.json: %v", err)
		}
		s.logger.Infof("Deployed contract addresses: %v", deployedAddresses)
	}

	return nil
}

func (s *Scenario) sendTx(ctx context.Context, txIdx uint64, onComplete func()) (*types.Transaction, *txbuilder.Client, *txbuilder.Wallet, error) {
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, int(txIdx), s.options.ClientGroup)
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(txIdx))

	if client == nil {
		return nil, nil, nil, fmt.Errorf("no client available")
	}
	if wallet == nil {
		return nil, nil, nil, fmt.Errorf("no wallet available")
	}

	// Get chain ID
	chainId, err := client.GetChainId(ctx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Deploy contract
	auth, err := bind.NewKeyedTransactorWithChainID(wallet.GetPrivateKey(), chainId)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create auth: %w", err)
	}

	auth.Nonce = big.NewInt(int64(wallet.GetNonce()))
	auth.Value = big.NewInt(0)
	auth.GasLimit = GAS_PER_CONTRACT + GAS_PER_CONTRACT*5/100 // Gas for single contract deployment

	// Set EIP-1559 fee parameters
	if s.options.BaseFee > 0 {
		auth.GasFeeCap = new(big.Int).Mul(big.NewInt(int64(s.options.BaseFee)), big.NewInt(1000000000))
	}
	if s.options.TipFee > 0 {
		auth.GasTipCap = new(big.Int).Mul(big.NewInt(int64(s.options.TipFee)), big.NewInt(1000000000))
	}

	// Generate random salt for unique contract
	salt := make([]byte, 32)
	_, err = rand.Read(salt)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	saltInt := new(big.Int).SetBytes(salt)

	// Ensure we have a valid client
	ethClient := client.GetEthClient()
	if ethClient == nil {
		return nil, nil, nil, fmt.Errorf("failed to get eth client")
	}

	// Deploy the StateBloatToken contract
	address, tx, _, err := contract.DeployContract(auth, ethClient, saltInt)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to deploy contract: %w", err)
	}

	// Update nonce
	wallet.SetNonce(wallet.GetNonce() + 1)

	// Calculate bytes written and gas/byte ratio
	txBytes := len(tx.Data())
	gasPerByte := float64(GAS_PER_CONTRACT) / float64(txBytes)

	s.logger.WithFields(logrus.Fields{
		"test":             "contract-deploy",
		"tx_hash":          tx.Hash().Hex(),
		"contract_address": address.Hex(),
		"bytes_written":    txBytes,
		"gas_per_byte":     fmt.Sprintf("%.2f", gasPerByte),
		"contracts":        1,
	}).Info("deployed contract")

	return tx, client, wallet, nil
}
