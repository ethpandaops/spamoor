package setcodetx

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"
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
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	pendingChan   chan bool
	pendingWGroup sync.WaitGroup
	delegatorSeed []byte
	delegators    []*txbuilder.Wallet
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
		s.walletPool.SetWalletCount(s.options.MaxWallets)
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

	if s.options.TotalCount == 0 && s.options.Throughput == 0 {
		return fmt.Errorf("neither total count nor throughput limit set, must define at least one of them (see --help for list of all flags)")
	}

	if s.options.MaxPending > 0 {
		s.pendingChan = make(chan bool, s.options.MaxPending)
	}

	s.delegatorSeed = make([]byte, 32)
	rand.Read(s.delegatorSeed)

	if s.options.MaxDelegators > 0 {
		s.delegators = make([]*txbuilder.Wallet, 0, s.options.MaxDelegators)
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
			tx, client, wallet, err := s.sendTx(ctx, txIdx)
			if client != nil {
				logger = logger.WithField("rpc", client.GetName())
			}
			if tx != nil {
				logger = logger.WithField("nonce", tx.Nonce())
			}
			if wallet != nil {
				logger = logger.WithField("wallet", s.walletPool.GetWalletIndex(wallet.GetAddress()))
			}
			if lastChan != nil {
				<-lastChan
				close(lastChan)
			}
			if err != nil {
				logger.Warnf("could not send transaction: %v", err)
				<-s.pendingChan
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

func (s *Scenario) sendTx(ctx context.Context, txIdx uint64) (*types.Transaction, *txbuilder.Client, *txbuilder.Wallet, error) {
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, int(txIdx))
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(txIdx))

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
	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &txbuilder.SendTransactionOptions{
		Client:              client,
		MaxRebroadcasts:     rebroadcast,
		RebroadcastInterval: time.Duration(s.options.Rebroadcast) * time.Second,
		OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			defer func() {
				if s.pendingChan != nil {
					time.Sleep(100 * time.Millisecond)
					<-s.pendingChan
				}
				s.pendingWGroup.Done()
			}()

			if err != nil {
				s.logger.WithField("client", client.GetName()).Warnf("tx %6d: await receipt failed: %v", txIdx+1, err)
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

			s.logger.WithField("client", client.GetName()).Debugf(" transaction %d confirmed in block #%v. total fee: %v gwei (base: %v) logs: %v", txIdx+1, receipt.BlockNumber.String(), gweiTotalFee, gweiBaseFee, len(receipt.Logs))
		},
		LogFn: func(client *txbuilder.Client, retry int, rebroadcast int, err error) {
			logger := s.logger.WithField("client", client.GetName())
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
			codeAddr = s.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(txIdx)).GetAddress()
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
