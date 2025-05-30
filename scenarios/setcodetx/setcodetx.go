package setcodetx

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"
	"gopkg.in/yaml.v3"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenariotypes"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/utils"
)

type ScenarioOptions struct {
	TotalCount        uint64 `yaml:"total_count"`
	Throughput        uint64 `yaml:"throughput"`
	MaxPending        uint64 `yaml:"max_pending"`
	MaxWallets        uint64 `yaml:"max_wallets"`
	MinAuthorizations uint64 `yaml:"min_authorizations"`
	MaxAuthorizations uint64 `yaml:"max_authorizations"`
	MaxDelegators     uint64 `yaml:"max_delegators"`
	Rebroadcast       uint64 `yaml:"rebroadcast"`
	BaseFee           uint64 `yaml:"base_fee"`
	TipFee            uint64 `yaml:"tip_fee"`
	GasLimit          uint64 `yaml:"gas_limit"`
	Amount            uint64 `yaml:"amount"`
	Data              string `yaml:"data"`
	CodeAddr          string `yaml:"code_addr"`
	RandomAmount      bool   `yaml:"random_amount"`
	RandomTarget      bool   `yaml:"random_target"`
	RandomCodeAddr    bool   `yaml:"random_code_addr"`
	MaxBloating       bool   `yaml:"max_bloating"`
	ClientGroup       string `yaml:"client_group"`
}

// EOAEntry represents a funded EOA account
type EOAEntry struct {
	Address    string `json:"address"`
	PrivateKey string `json:"private_key"`
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	pendingChan   chan bool
	pendingWGroup sync.WaitGroup
	delegatorSeed []byte
	delegators    []*txbuilder.Wallet

	// FIFO queue for funded accounts
	eoaQueue      []EOAEntry
	eoaQueueMutex sync.Mutex

	// Semaphore for worker control
	workerSemaphore chan struct{}
	workerDone      chan struct{}
	workerWg        sync.WaitGroup
}

var ScenarioName = "setcodetx"
var ScenarioDefaultOptions = ScenarioOptions{
	TotalCount:        0,
	Throughput:        0,
	MaxPending:        0,
	MaxWallets:        0,
	MinAuthorizations: 1,
	MaxAuthorizations: 10,
	MaxDelegators:     0,
	Rebroadcast:       120,
	BaseFee:           20,
	TipFee:            2,
	GasLimit:          200000,
	Amount:            20,
	Data:              "",
	CodeAddr:          "",
	RandomAmount:      false,
	RandomTarget:      false,
	RandomCodeAddr:    false,
	MaxBloating:       false,
	ClientGroup:       "",
}
var ScenarioDescriptor = scenariotypes.ScenarioDescriptor{
	Name:           ScenarioName,
	Description:    "Send setcode transactions with different configurations",
	DefaultOptions: ScenarioDefaultOptions,
	NewScenario:    newScenario,
}

func newScenario(logger logrus.FieldLogger) scenariotypes.Scenario {
	return &Scenario{
		logger: logger.WithField("scenario", ScenarioName),
	}
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
	flags.Uint64VarP(&s.options.TotalCount, "count", "c", ScenarioDefaultOptions.TotalCount, "Total number of transfer transactions to send")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t", ScenarioDefaultOptions.Throughput, "Number of transfer transactions to send per slot")
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.MinAuthorizations, "min-authorizations", ScenarioDefaultOptions.MinAuthorizations, "Minimum number of authorizations to send per transaction")
	flags.Uint64Var(&s.options.MaxAuthorizations, "max-authorizations", ScenarioDefaultOptions.MaxAuthorizations, "Maximum number of authorizations to send per transaction")
	flags.Uint64Var(&s.options.MaxDelegators, "max-delegators", ScenarioDefaultOptions.MaxDelegators, "Maximum number of random delegators to use (0 = no delegator gets reused)")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", ScenarioDefaultOptions.Rebroadcast, "Number of seconds to wait before re-broadcasting a transaction")
	flags.Uint64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use in transfer transactions (in gwei)")
	flags.Uint64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use in transfer transactions (in gwei)")
	flags.Uint64Var(&s.options.GasLimit, "gaslimit", ScenarioDefaultOptions.GasLimit, "Gas limit to use in transactions")
	flags.Uint64Var(&s.options.Amount, "amount", ScenarioDefaultOptions.Amount, "Transfer amount per transaction (in gwei)")
	flags.StringVar(&s.options.Data, "data", ScenarioDefaultOptions.Data, "Transaction call data to send")
	flags.StringVar(&s.options.CodeAddr, "code-addr", ScenarioDefaultOptions.CodeAddr, "Code delegation target address to use for transactions")
	flags.BoolVar(&s.options.RandomAmount, "random-amount", ScenarioDefaultOptions.RandomAmount, "Use random amounts for transactions (with --amount as limit)")
	flags.BoolVar(&s.options.RandomTarget, "random-target", ScenarioDefaultOptions.RandomTarget, "Use random to addresses for transactions")
	flags.BoolVar(&s.options.RandomCodeAddr, "random-code-addr", ScenarioDefaultOptions.RandomCodeAddr, "Use random delegation target for transactions")
	flags.BoolVar(&s.options.MaxBloating, "max-bloating", ScenarioDefaultOptions.MaxBloating, "Enable maximum state bloating mode: creates ~960 EOA delegations in a single block-filling transaction")
	flags.StringVar(&s.options.ClientGroup, "client-group", ScenarioDefaultOptions.ClientGroup, "Client group to use for sending transactions")
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

	if s.options.MaxBloating && !s.options.RandomCodeAddr && s.options.CodeAddr == "" {
		s.logger.Infof("no --code-addr specified, using ecrecover precompile as delegate: %s", common.HexToAddress("0x0000000000000000000000000000000000000001"))
	}

	if s.options.MaxWallets > 0 {
		s.walletPool.SetWalletCount(s.options.MaxWallets)
	} else if s.options.MaxBloating {
		// In max-bloating mode, if maxWallets is not set, use 1 wallet (the root wallet)
		s.walletPool.SetWalletCount(1)
	} else if s.options.TotalCount > 0 {
		if s.options.TotalCount < 1000 {
			s.walletPool.SetWalletCount(s.options.TotalCount)
		} else {
			s.walletPool.SetWalletCount(1000)
		}
	} else {
		if s.options.Throughput*10 < 1000 {
			s.walletPool.SetWalletCount(s.options.Throughput * 10)
		} else {
			s.walletPool.SetWalletCount(1000)
		}
	}

	if s.options.MaxPending > 0 {
		s.pendingChan = make(chan bool, s.options.MaxPending)
	}

	s.delegatorSeed = make([]byte, 32)
	rand.Read(s.delegatorSeed)

	if s.options.MaxDelegators > 0 {
		s.delegators = make([]*txbuilder.Wallet, 0, s.options.MaxDelegators)
	}

	// In max-bloating mode, we handle throughput automatically, so skip this validation
	if !s.options.MaxBloating && s.options.TotalCount == 0 && s.options.Throughput == 0 {
		return fmt.Errorf("neither total count nor throughput limit set, must define at least one of them (see --help for list of all flags)")
	}

	// Initialize FIFO queue and worker for EOA management
	s.eoaQueue = make([]EOAEntry, 0)
	s.workerSemaphore = make(chan struct{}, 1) // Buffered channel for semaphore
	s.workerDone = make(chan struct{})

	// Start the worker goroutine for writing EOAs to file
	if s.options.MaxBloating {
		s.workerWg.Add(1)
		go s.eoaWorker()
	}

	return nil
}

func (s *Scenario) Config() string {
	yamlBytes, _ := yaml.Marshal(&s.options)
	return string(yamlBytes)
}

func (s *Scenario) Run(ctx context.Context) error {
	// Check if max bloating mode is enabled
	if s.options.MaxBloating {
		return s.runMaxBloatingMode(ctx)
	}

	// Original implementation for backward compatibility
	txIdxCounter := uint64(0)
	pendingCount := atomic.Int64{}
	txCount := atomic.Uint64{}

	s.logger.Infof("starting scenario: %s", ScenarioName)
	defer s.logger.Infof("scenario %s finished.", ScenarioName)

	var lastChan chan bool

	initialRate := rate.Limit(float64(s.options.Throughput) / float64(utils.SecondsPerSlot))
	if initialRate == 0 {
		initialRate = rate.Inf
	}
	limiter := rate.NewLimiter(initialRate, 1)

	for {
		if err := limiter.Wait(ctx); err != nil {
			if ctx.Err() != nil {
				break
			}

			s.logger.Debugf("rate limited: %s", err.Error())
			time.Sleep(100 * time.Millisecond)
			continue
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
				close(lastChan)
			}
			if err != nil {
				logger.Warnf("could not send transaction: %v", err)
				return
			}

			txCount.Add(1)
			logger.Infof("sent tx #%6d: %v", txIdx+1, tx.Hash().String())
		}(txIdx, lastChan, currentChan)

		lastChan = currentChan

		count := txCount.Load() + uint64(pendingCount.Load())
		if s.options.TotalCount > 0 && count >= s.options.TotalCount {
			break
		}
	}

	<-lastChan
	close(lastChan)

	s.logger.Infof("finished sending transactions, awaiting block inclusion...")
	s.pendingWGroup.Wait()

	return nil
}

func (s *Scenario) sendTx(ctx context.Context, txIdx uint64, onComplete func()) (*types.Transaction, *txbuilder.Client, *txbuilder.Wallet, error) {
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

	var feeCap *big.Int
	var tipCap *big.Int

	if s.options.BaseFee > 0 {
		feeCap = new(big.Int).Mul(big.NewInt(int64(s.options.BaseFee)), big.NewInt(1000000000))
	}
	if s.options.TipFee > 0 {
		tipCap = new(big.Int).Mul(big.NewInt(int64(s.options.TipFee)), big.NewInt(1000000000))
	}

	if feeCap == nil || tipCap == nil {
		var err error
		feeCap, tipCap, err = client.GetSuggestedFee(s.walletPool.GetContext())
		if err != nil {
			return nil, client, wallet, err
		}
	}

	if feeCap.Cmp(big.NewInt(1000000000)) < 0 {
		feeCap = big.NewInt(1000000000)
	}
	if tipCap.Cmp(big.NewInt(1000000000)) < 0 {
		tipCap = big.NewInt(1000000000)
	}

	amount := uint256.NewInt(s.options.Amount)
	amount = amount.Mul(amount, uint256.NewInt(1000000000))
	if s.options.RandomAmount {
		n, err := rand.Int(rand.Reader, amount.ToBig())
		if err == nil {
			amount = uint256.MustFromBig(n)
		}
	}

	toAddr := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(txIdx)+1).GetAddress()
	if s.options.RandomTarget {
		addrBytes := make([]byte, 20)
		rand.Read(addrBytes)
		toAddr = common.Address(addrBytes)
	}

	txCallData := []byte{}

	if s.options.Data != "" {
		dataBytes, err := txbuilder.ParseBlobRefsBytes(strings.Split(s.options.Data, ","), nil)
		if err != nil {
			return nil, nil, wallet, err
		}

		txCallData = dataBytes
	}

	txData, err := txbuilder.SetCodeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       s.options.GasLimit,
		To:        &toAddr,
		Value:     amount,
		Data:      txCallData,
		AuthList:  s.buildSetCodeAuthorizations(txIdx),
	})
	if err != nil {
		return nil, nil, wallet, err
	}

	tx, err := wallet.BuildSetCodeTx(txData)
	if err != nil {
		return nil, nil, wallet, err
	}

	rebroadcast := 0
	if s.options.Rebroadcast > 0 {
		rebroadcast = 10
	}

	s.pendingWGroup.Add(1)
	transactionSubmitted = true
	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &txbuilder.SendTransactionOptions{
		Client:              client,
		MaxRebroadcasts:     rebroadcast,
		RebroadcastInterval: time.Duration(s.options.Rebroadcast) * time.Second,
		OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			defer func() {
				onComplete()
				s.pendingWGroup.Done()
			}()

			if err != nil {
				s.logger.WithField("rpc", client.GetName()).Warnf("tx %6d: await receipt failed: %v", txIdx+1, err)
				return
			}
			if receipt == nil {
				return
			}

			effectiveGasPrice := receipt.EffectiveGasPrice
			if effectiveGasPrice == nil {
				effectiveGasPrice = big.NewInt(0)
			}
			feeAmount := new(big.Int).Mul(effectiveGasPrice, big.NewInt(int64(receipt.GasUsed)))
			totalAmount := new(big.Int).Add(tx.Value(), feeAmount)
			wallet.SubBalance(totalAmount)

			gweiTotalFee := new(big.Int).Div(feeAmount, big.NewInt(1000000000))
			gweiBaseFee := new(big.Int).Div(effectiveGasPrice, big.NewInt(1000000000))

			s.logger.WithField("rpc", client.GetName()).Debugf(" transaction %d confirmed in block #%v. total fee: %v gwei (base: %v) logs: %v", txIdx+1, receipt.BlockNumber.String(), gweiTotalFee, gweiBaseFee, len(receipt.Logs))
		},
		LogFn: func(client *txbuilder.Client, retry int, rebroadcast int, err error) {
			logger := s.logger.WithField("rpc", client.GetName())
			if retry > 0 {
				logger = logger.WithField("retry", retry)
			}
			if rebroadcast > 0 {
				logger = logger.WithField("rebroadcast", rebroadcast)
			}
			if err != nil {
				logger.Warnf("failed sending tx %6d: %v", txIdx+1, err)
			} else if retry > 0 || rebroadcast > 0 {
				logger.Debugf("successfully sent tx %6d", txIdx+1)
			}
		},
	})
	if err != nil {
		// reset nonce if tx was not sent
		wallet.ResetPendingNonce(ctx, client)

		return nil, client, wallet, err
	}

	return tx, client, wallet, nil
}

func (s *Scenario) buildSetCodeAuthorizations(txIdx uint64) []types.SetCodeAuthorization {
	authorizations := []types.SetCodeAuthorization{}

	if s.options.MaxAuthorizations == 0 {
		return authorizations
	}

	n, _ := rand.Int(rand.Reader, big.NewInt(int64(s.options.MaxAuthorizations-s.options.MinAuthorizations+1)))
	authorizationCount := int(n.Int64()) + int(s.options.MinAuthorizations)

	for i := 0; i < authorizationCount; i++ {
		delegatorIndex := (txIdx * s.options.MaxAuthorizations) + uint64(i)
		if s.options.MaxDelegators > 0 {
			delegatorIndex = delegatorIndex % s.options.MaxDelegators
		}

		var delegator *txbuilder.Wallet
		if s.options.MaxDelegators > 0 && len(s.delegators) > int(delegatorIndex) {
			delegator = s.delegators[delegatorIndex]
		} else {
			d, err := s.prepareDelegator(delegatorIndex)
			if err != nil {
				s.logger.Errorf("could not prepare delegator %v: %v", delegatorIndex, err)
				continue
			}

			delegator = d

			if s.options.MaxDelegators > 0 {
				s.delegators = append(s.delegators, delegator)
			}
		}

		var codeAddr common.Address
		if s.options.RandomCodeAddr {
			codeAddr = common.Address(make([]byte, 20))
			rand.Read(codeAddr[:])
		} else if s.options.CodeAddr != "" {
			codeAddr = common.HexToAddress(s.options.CodeAddr)
		} else {
			// Use a fixed delegate contract address for maximum efficiency
			// In max bloating mode, we want all EOAs to delegate to the same existing contract
			// to benefit from reduced gas costs (PER_AUTH_BASE_COST vs PER_EMPTY_ACCOUNT_COST)
			// Precompiles are ideal as they're guaranteed to exist with code on all networks
			codeAddr = common.HexToAddress("0x0000000000000000000000000000000000000001")
		}

		authorization := types.SetCodeAuthorization{
			ChainID: *uint256.NewInt(s.walletPool.GetRootWallet().GetChainId().Uint64()),
			Address: codeAddr,
			Nonce:   delegator.GetNextNonce(),
		}

		authorization, err := types.SignSetCode(delegator.GetPrivateKey(), authorization)
		if err != nil {
			s.logger.Errorf("could not sign set code authorization: %v", err)
			continue
		}

		authorizations = append(authorizations, authorization)
	}

	return authorizations
}

func (s *Scenario) prepareDelegator(delegatorIndex uint64) (*txbuilder.Wallet, error) {
	idxBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(idxBytes, delegatorIndex)
	if s.options.MaxDelegators > 0 {
		seedBytes := []byte(s.delegatorSeed)
		idxBytes = append(idxBytes, seedBytes...)
	}
	childKey := sha256.Sum256(append(common.FromHex(s.walletPool.GetRootWallet().GetAddress().Hex()), idxBytes...))
	return txbuilder.NewWallet(fmt.Sprintf("%x", childKey))
}

func (s *Scenario) buildMaxBloatingAuthorizations(targetCount int, iteration int) []types.SetCodeAuthorization {
	authorizations := make([]types.SetCodeAuthorization, 0, targetCount)

	// Use a fixed delegate contract address for maximum efficiency
	// In max bloating mode, we want all EOAs to delegate to the same existing contract
	// to benefit from reduced gas costs (PER_AUTH_BASE_COST vs PER_EMPTY_ACCOUNT_COST)
	// Precompiles are ideal as they're guaranteed to exist with code on all networks
	var codeAddr common.Address
	if s.options.CodeAddr != "" {
		codeAddr = common.HexToAddress(s.options.CodeAddr)
	} else {
		// Default to using the ecrecover precompile (0x1) as delegate target
		codeAddr = common.HexToAddress("0x0000000000000000000000000000000000000001")
	}

	chainId := s.walletPool.GetRootWallet().GetChainId().Uint64()

	for i := 0; i < targetCount; i++ {
		// Create a unique delegator for each authorization
		// Include iteration counter to ensure different addresses for each iteration
		delegatorIndex := uint64(iteration*targetCount + i)

		delegator, err := s.prepareDelegator(delegatorIndex)
		if err != nil {
			s.logger.Errorf("could not prepare delegator %v: %v", delegatorIndex, err)
			continue
		}

		// Each EOA uses auth_nonce = 0 (assuming first EIP-7702 operation)
		// This creates maximum new state as each EOA gets its first delegation
		authorization := types.SetCodeAuthorization{
			ChainID: *uint256.NewInt(chainId),
			Address: codeAddr,
			Nonce:   0, // First delegation for each EOA
		}

		// Sign the authorization with the delegator's private key
		signedAuth, err := types.SignSetCode(delegator.GetPrivateKey(), authorization)
		if err != nil {
			s.logger.Errorf("could not sign set code authorization for delegator %v: %v", delegatorIndex, err)
			continue
		}

		authorizations = append(authorizations, signedAuth)
	}

	return authorizations
}

func (s *Scenario) runMaxBloatingMode(ctx context.Context) error {
	s.logger.Infof("starting max bloating mode: self-adjusting to target 29.9M gas per block, continuous operation")

	// Dynamic authorization count - starts conservatively and adjusts based on actual performance
	// TODO: This should be set as a constant or similar computed as block_gas_limit/Gas_per_Auth (26000)
	currentAuthorizations := 1000 // Start with known working value
	// TODO: This should be obtained from the network.
	targetGas := uint64(29900000) // Target 29.9M gas

	var blockCounter int

	for {
		select {
		case <-ctx.Done():
			s.logger.Errorf("max bloating mode stopping due to context cancellation")
			return ctx.Err()
		default:
		}

		blockCounter++

		// Fund delegator accounts first
		s.logger.Infof("════════════════ FUNDING PHASE #%d ════════════════", blockCounter)
		err := s.fundMaxBloatingDelegators(ctx, currentAuthorizations, blockCounter)
		if err != nil {
			s.logger.Errorf("failed to fund delegators for iteration %d: %v", blockCounter, err)
			time.Sleep(5 * time.Second) // Wait before retry
			continue
		}

		// Send the max bloating transaction and wait for confirmation
		s.logger.Infof("════════════════ BLOATING PHASE #%d ════════════════", blockCounter)
		actualGasUsed, blockNumber, authCount, gasPerAuth, gasPerByte, gweiTotalFee, err := s.sendMaxBloatingTransaction(ctx, currentAuthorizations, targetGas, blockCounter)
		if err != nil {
			s.logger.Errorf("failed to send max bloating transaction for iteration %d: %v", blockCounter, err)
			time.Sleep(5 * time.Second) // Wait before retry
			continue
		}
		// Open semaphore (green light) during analysis phase to allow worker to process EOA queue
		s.openWorkerSemaphore()

		s.logger.Infof("%%%%%%%%%%%%%%%%%%%% ANALYSIS PHASE #%d %%%%%%%%%%%%%%%%%%%%", blockCounter)
		s.logger.WithField("scenario", "setcodetx").Infof("MAX BLOATING TX MINED - Block #%s, Gas Used: %d, Authorizations: %d, Gas/Auth: %.1f, Gas/Byte: %.1f, Total Fee: %s gwei",
			blockNumber, actualGasUsed, authCount, gasPerAuth, gasPerByte, gweiTotalFee)

		// Self-adjust authorization count based on actual performance
		if actualGasUsed > 0 && authCount > 0 {
			gasPerAuth := float64(actualGasUsed) / float64(authCount)
			targetAuths := int(float64(targetGas) / gasPerAuth)

			// Calculate the adjustment needed
			authDifference := targetAuths - authCount

			if actualGasUsed < targetGas {
				// We're under target, increase authorization count with a slight safety margin
				newAuthorizations := currentAuthorizations + authDifference - 1

				if newAuthorizations > currentAuthorizations {
					s.logger.Infof("Adjusting authorizations: %d → %d (need %d more for target)",
						currentAuthorizations, newAuthorizations, authDifference)
					currentAuthorizations = newAuthorizations
				}
			} else if actualGasUsed > targetGas {
				// We're over target, reduce to reach max block utilization
				excess := actualGasUsed - targetGas
				newAuthorizations := currentAuthorizations - int(excess) + 1

				s.logger.Infof("Reducing authorizations: %d → %d (excess: %d gas)",
					currentAuthorizations, newAuthorizations, excess)
				currentAuthorizations = newAuthorizations

			} else {
				s.logger.Infof("Target achieved! Gas Used: %d / Target: %d", actualGasUsed, targetGas)
			}
		}

		// Transaction confirmed - small delay before next iteration
		time.Sleep(2 * time.Second)
	}
}

func (s *Scenario) fundMaxBloatingDelegators(ctx context.Context, targetCount int, iteration int) error {
	// Close semaphore (red light) during funding phase
	s.closeWorkerSemaphore()

	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, s.options.ClientGroup)
	if client == nil {
		return fmt.Errorf("no client available for funding delegators")
	}

	// Use root wallet since we set child wallet count to 0 in max-bloating mode
	wallet := s.walletPool.GetRootWallet()

	// Get suggested fees for funding transactions
	var feeCap *big.Int
	var tipCap *big.Int

	if s.options.BaseFee > 0 {
		feeCap = new(big.Int).Mul(big.NewInt(int64(s.options.BaseFee)), big.NewInt(1000000000))
	}
	if s.options.TipFee > 0 {
		tipCap = new(big.Int).Mul(big.NewInt(int64(s.options.TipFee)), big.NewInt(1000000000))
	}

	if feeCap == nil || tipCap == nil {
		var err error
		feeCap, tipCap, err = client.GetSuggestedFee(s.walletPool.GetContext())
		if err != nil {
			return fmt.Errorf("failed to get suggested fees for funding: %w", err)
		}
	}

	// Minimum gas prices
	if feeCap.Cmp(big.NewInt(1000000000)) < 0 {
		feeCap = big.NewInt(1000000000)
	}
	if tipCap.Cmp(big.NewInt(1000000000)) < 0 {
		tipCap = big.NewInt(1000000000)
	}

	// Fund with 1 wei as requested by user
	fundingAmount := uint256.NewInt(1)

	var confirmedCount int64
	sentCount := 0
	delegatorIndex := uint64(iteration * 1000000) // Large offset per iteration to avoid conflicts

	// Calculate approximate transactions per block based on gas limit
	// Standard transfer = 21000 gas, typical block = ~30M gas = ~1400 txs per block
	const maxTxsPerBlock = 1400
	const blockTimeSeconds = 10 // Adjust based on your network

	for {
		// Check if we have enough confirmed transactions
		confirmed := atomic.LoadInt64(&confirmedCount)
		if confirmed >= int64(targetCount) {
			// We have minimum required, but let's check if we should fill the current block
			// If we've sent transactions recently, wait a bit to see if block gets filled
			if sentCount > 0 && (sentCount%100) > 50 {
				// Continue to fill the block
			} else {
				break
			}
		}

		// Generate unique delegator address
		delegator, err := s.prepareDelegator(delegatorIndex)
		if err != nil {
			s.logger.Errorf("could not prepare delegator %v for funding: %v", delegatorIndex, err)
			delegatorIndex++
			continue
		}

		// Build funding transaction
		delegatorAddr := delegator.GetAddress()
		txData, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Gas:       21000, // Standard ETH transfer gas
			To:        &delegatorAddr,
			Value:     fundingAmount,
			Data:      []byte{},
		})
		if err != nil {
			s.logger.Errorf("failed to build funding tx for delegator %d: %v", delegatorIndex, err)
			delegatorIndex++
			continue
		}

		tx, err := wallet.BuildDynamicFeeTx(txData)
		if err != nil {
			s.logger.Errorf("failed to build funding transaction for delegator %d: %v", delegatorIndex, err)
			delegatorIndex++
			continue
		}

		// Send funding transaction with no retries to avoid duplicates
		err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &txbuilder.SendTransactionOptions{
			Client:          client,
			MaxRebroadcasts: 0, // No retries to avoid duplicates
			OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
				if err != nil {
					return // Don't log individual failures
				}
				if receipt != nil && receipt.Status == 1 {
					atomic.AddInt64(&confirmedCount, 1)

					// Add successfully funded delegator to EOA queue
					s.addEOAToQueue(delegator.GetAddress().Hex(), fmt.Sprintf("%x", delegator.GetPrivateKey().D))

					// No progress logging - only log when target is reached
				}
			},
			LogFn: func(client *txbuilder.Client, retry int, rebroadcast int, err error) {
				// Only log actual send failures, not confirmation failures
				if err != nil {
					s.logger.Debugf("funding tx send failed: %v", err)
				}
			},
		})

		if err != nil {
			s.logger.Debugf("failed to send funding transaction for delegator %d: %v", delegatorIndex, err)
			delegatorIndex++
			continue
		}

		sentCount++
		delegatorIndex++

		// Check if we should continue filling the block
		confirmed = atomic.LoadInt64(&confirmedCount)
		if confirmed >= int64(targetCount) && sentCount%maxTxsPerBlock < 100 {
			// We have enough confirmed and we're at the end of a block cycle, stop for now
			break
		}

		// Small delay between transactions to ensure proper nonce ordering
		// Reduce delay as we get more efficient
		if sentCount < 100 {
			time.Sleep(10 * time.Millisecond)
		} else {
			time.Sleep(5 * time.Millisecond)
		}

		// Add context cancellation check
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}

	// Wait for any remaining transactions to be included
	s.logger.Debugf("waiting for remaining funding transactions to be confirmed...")
	time.Sleep(3 * time.Second)

	return nil
}

func (s *Scenario) sendMaxBloatingTransaction(ctx context.Context, targetAuthorizations int, targetGasLimit uint64, blockCounter int) (uint64, string, int, float64, float64, string, error) {
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, s.options.ClientGroup)
	if client == nil {
		return 0, "", 0, 0, 0, "", fmt.Errorf("no client available for sending max bloating transaction")
	}

	// Use root wallet since we set child wallet count to 0 in max-bloating mode
	wallet := s.walletPool.GetRootWallet()

	// Get suggested fees or use configured values
	var feeCap *big.Int
	var tipCap *big.Int

	if s.options.BaseFee > 0 {
		feeCap = new(big.Int).Mul(big.NewInt(int64(s.options.BaseFee)), big.NewInt(1000000000))
	}
	if s.options.TipFee > 0 {
		tipCap = new(big.Int).Mul(big.NewInt(int64(s.options.TipFee)), big.NewInt(1000000000))
	}

	if feeCap == nil || tipCap == nil {
		var err error
		feeCap, tipCap, err = client.GetSuggestedFee(s.walletPool.GetContext())
		if err != nil {
			return 0, "", 0, 0, 0, "", fmt.Errorf("failed to get suggested fees: %w", err)
		}
	}

	// Ensure minimum gas prices for inclusion
	if feeCap.Cmp(big.NewInt(1000000000)) < 0 {
		feeCap = big.NewInt(1000000000)
	}
	if tipCap.Cmp(big.NewInt(1000000000)) < 0 {
		tipCap = big.NewInt(1000000000)
	}

	// Use minimal amount for max bloating (focus on authorizations, not value transfer)
	amount := uint256.NewInt(0) // No value transfer needed

	// Target address - use our own wallet for simplicity
	toAddr := wallet.GetAddress()

	// Minimal call data
	txCallData := []byte{}
	if s.options.Data != "" {
		dataBytes, err := txbuilder.ParseBlobRefsBytes(strings.Split(s.options.Data, ","), nil)
		if err != nil {
			return 0, "", 0, 0, 0, "", fmt.Errorf("failed to parse call data: %w", err)
		}
		txCallData = dataBytes
	}

	// Build the authorizations for maximum state bloat
	authorizations := s.buildMaxBloatingAuthorizations(targetAuthorizations, blockCounter)

	txData, err := txbuilder.SetCodeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       targetGasLimit,
		To:        &toAddr,
		Value:     amount,
		Data:      txCallData,
		AuthList:  authorizations,
	})
	if err != nil {
		return 0, "", 0, 0, 0, "", fmt.Errorf("failed to build transaction metadata: %w", err)
	}

	tx, err := wallet.BuildSetCodeTx(txData)
	if err != nil {
		return 0, "", 0, 0, 0, "", fmt.Errorf("failed to build transaction: %w", err)
	}

	// Use channels to capture transaction results
	resultChan := make(chan struct {
		gasUsed      uint64
		blockNumber  string
		authCount    int
		gasPerAuth   float64
		gasPerByte   float64
		gweiTotalFee string
		err          error
	}, 1)

	// Send the transaction
	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &txbuilder.SendTransactionOptions{
		Client:              client,
		MaxRebroadcasts:     10,
		RebroadcastInterval: time.Duration(s.options.Rebroadcast) * time.Second,
		OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			if err != nil {
				s.logger.WithField("rpc", client.GetName()).Errorf("max bloating tx failed: %v", err)
				resultChan <- struct {
					gasUsed      uint64
					blockNumber  string
					authCount    int
					gasPerAuth   float64
					gasPerByte   float64
					gweiTotalFee string
					err          error
				}{0, "", 0, 0, 0, "", err}
				return
			}
			if receipt == nil {
				resultChan <- struct {
					gasUsed      uint64
					blockNumber  string
					authCount    int
					gasPerAuth   float64
					gasPerByte   float64
					gweiTotalFee string
					err          error
				}{0, "", 0, 0, 0, "", fmt.Errorf("no receipt received")}
				return
			}

			effectiveGasPrice := receipt.EffectiveGasPrice
			if effectiveGasPrice == nil {
				effectiveGasPrice = big.NewInt(0)
			}
			feeAmount := new(big.Int).Mul(effectiveGasPrice, big.NewInt(int64(receipt.GasUsed)))
			totalAmount := new(big.Int).Add(tx.Value(), feeAmount)
			wallet.SubBalance(totalAmount)

			gweiTotalFee := new(big.Int).Div(feeAmount, big.NewInt(1000000000))

			// Calculate efficiency metrics
			authCount := len(authorizations)
			gasPerAuth := float64(receipt.GasUsed) / float64(authCount)
			// TODO: This should be a constant.
			estimatedBytesPerAuth := 135.0 // ~135 bytes state change per EOA delegation
			gasPerByte := gasPerAuth / estimatedBytesPerAuth

			resultChan <- struct {
				gasUsed      uint64
				blockNumber  string
				authCount    int
				gasPerAuth   float64
				gasPerByte   float64
				gweiTotalFee string
				err          error
			}{receipt.GasUsed, receipt.BlockNumber.String(),
				authCount,
				gasPerAuth,
				gasPerByte,
				gweiTotalFee.String(), nil}
		},
		LogFn: func(client *txbuilder.Client, retry int, rebroadcast int, err error) {
			logger := s.logger.WithField("rpc", client.GetName())
			if retry > 0 {
				logger = logger.WithField("retry", retry)
			}
			if rebroadcast > 0 {
				logger = logger.WithField("rebroadcast", rebroadcast)
			}
			if err != nil {
				logger.Errorf("failed sending max bloating tx: %v", err)
			} else if retry > 0 || rebroadcast > 0 {
				logger.Infof("successfully sent max bloating tx")
			}
		},
	})

	if err != nil {
		wallet.ResetPendingNonce(ctx, client)
		return 0, "", 0, 0, 0, "", fmt.Errorf("failed to send max bloating transaction: %w", err)
	}

	// Wait for transaction confirmation
	result := <-resultChan
	return result.gasUsed, result.blockNumber, result.authCount, result.gasPerAuth, result.gasPerByte, result.gweiTotalFee, result.err
}

// eoaWorker runs in a separate goroutine and writes funded EOAs to EOAs.json
// when the semaphore is open (green). It sleeps when the semaphore is closed (red).
func (s *Scenario) eoaWorker() {
	defer s.workerWg.Done()

	for {
		select {
		case <-s.workerDone:
			// Shutdown signal received
			return
		case <-s.workerSemaphore:
			// Semaphore is green, process the queue
			s.processEOAQueue()
		}
	}
}

// processEOAQueue drains the EOA queue and writes entries to EOAs.json
func (s *Scenario) processEOAQueue() {
	for {
		// Check if there are items in the queue
		s.eoaQueueMutex.Lock()
		if len(s.eoaQueue) == 0 {
			s.eoaQueueMutex.Unlock()
			return // Queue is empty, exit processing
		}

		// Dequeue all items (FIFO)
		eoasToWrite := make([]EOAEntry, len(s.eoaQueue))
		copy(eoasToWrite, s.eoaQueue)
		s.eoaQueue = s.eoaQueue[:0] // Clear the queue
		s.eoaQueueMutex.Unlock()

		// Write to file
		err := s.writeEOAsToFile(eoasToWrite)
		if err != nil {
			s.logger.Errorf("failed to write EOAs to file: %v", err)
			// Re-queue the items if write failed
			s.eoaQueueMutex.Lock()
			s.eoaQueue = append(eoasToWrite, s.eoaQueue...)
			s.eoaQueueMutex.Unlock()
			return
		}

		s.logger.Debugf("wrote %d EOA entries to EOAs.json", len(eoasToWrite))
		s.logger.Debugf("left %d EOA entries in queue", len(s.eoaQueue))
	}
}

// writeEOAsToFile appends EOA entries to EOAs.json file
func (s *Scenario) writeEOAsToFile(eoas []EOAEntry) error {
	if len(eoas) == 0 {
		return nil
	}

	fileName := "EOAs.json"

	// Read existing entries if file exists
	var existingEntries []EOAEntry
	if data, err := os.ReadFile(fileName); err == nil {
		json.Unmarshal(data, &existingEntries)
	}

	// Append new entries
	allEntries := append(existingEntries, eoas...)

	// Write back to file
	data, err := json.MarshalIndent(allEntries, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal EOA entries: %w", err)
	}

	err = os.WriteFile(fileName, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write EOAs.json: %w", err)
	}

	return nil
}

// addEOAToQueue adds a funded EOA to the queue
func (s *Scenario) addEOAToQueue(address, privateKey string) {
	s.eoaQueueMutex.Lock()
	defer s.eoaQueueMutex.Unlock()

	entry := EOAEntry{
		Address:    address,
		PrivateKey: privateKey,
	}

	s.eoaQueue = append(s.eoaQueue, entry)
}

// openWorkerSemaphore opens the semaphore (green light) allowing the worker to process
func (s *Scenario) openWorkerSemaphore() {
	select {
	case s.workerSemaphore <- struct{}{}:
		// Semaphore opened successfully
	default:
		// Semaphore already open, do nothing
	}
}

// closeWorkerSemaphore closes the semaphore (red light) putting the worker to sleep
func (s *Scenario) closeWorkerSemaphore() {
	select {
	case <-s.workerSemaphore:
		// Semaphore closed successfully
	default:
		// Semaphore already closed, do nothing
	}
}

// shutdownWorker signals the worker to stop and waits for it to finish
func (s *Scenario) shutdownWorker() {
	if s.options.MaxBloating {
		close(s.workerDone)
		s.workerWg.Wait()
	}
}
