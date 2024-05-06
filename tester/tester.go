package tester

import (
	"fmt"
	"math/big"
	"math/rand"
	"sync"
	"time"

	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/utils"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
)

type Tester struct {
	config         *TesterConfig
	logger         *logrus.Entry
	running        bool
	scenario       string
	chainId        *big.Int
	selectionMutex sync.Mutex
	allClients     []*txbuilder.Client
	goodClients    []*txbuilder.Client
	rrClientIdx    int
	rootWallet     *txbuilder.Wallet
	childWallets   []*txbuilder.Wallet
	rrWalletIdx    int
}

type TesterConfig struct {
	RpcHosts      []string     // rpc host urls to use for blob tests
	WalletPrivkey string       // pre-funded wallet privkey to use for blob tests
	WalletCount   uint64       // number of child wallets to generate & use (based on walletPrivkey)
	WalletPrefund *uint256.Int // amount of funds to send to each child wallet
	WalletMinfund *uint256.Int // min amount of funds child wallets should hold - refill with walletPrefund if lower
}

func NewTester(config *TesterConfig) *Tester {
	return &Tester{
		config: config,
		logger: logrus.NewEntry(logrus.StandardLogger()),
	}
}

func (tester *Tester) SetScenario(name string) {
	tester.scenario = name
	tester.logger = logrus.WithField("tester", name)
}

func (tester *Tester) Start(seed string) error {
	var err error
	if tester.running {
		return fmt.Errorf("already started")
	}
	tester.running = true

	tester.logger.WithFields(logrus.Fields{
		"version": utils.GetBuildVersion(),
	}).Infof("starting blob testing tool")

	// prepare clients
	err = tester.PrepareClients()
	if err != nil {
		return err
	}
	err = tester.watchClientStatus()
	if err != nil {
		return err
	}
	// watch client status
	go tester.watchClientStatusLoop()

	// prepare wallets
	err = tester.PrepareWallets(seed)
	if err != nil {
		return err
	}

	// watch wallet balances
	go tester.watchWalletBalancesLoop()

	return nil
}

func (tester *Tester) Stop() {
	if tester.running {
		tester.running = false
	}
}

func (tester *Tester) watchClientStatusLoop() {
	sleepTime := 2 * time.Minute
	for tester.running {
		time.Sleep(sleepTime)

		err := tester.watchClientStatus()
		if err != nil {
			tester.logger.Warnf("could not check client status: %v", err)
			sleepTime = 10 * time.Second
		} else {
			sleepTime = 2 * time.Minute
		}
	}
}

func (tester *Tester) watchWalletBalancesLoop() {
	sleepTime := 10 * time.Minute
	for tester.running {
		time.Sleep(sleepTime)

		err := tester.resupplyChildWallets()
		if err != nil {
			tester.logger.Warnf("could not check & resupply chile wallets: %v", err)
			sleepTime = 1 * time.Minute
		} else {
			sleepTime = 10 * time.Minute
		}
	}
}

type SelectionMode uint8

var (
	SelectByIndex    SelectionMode = 0
	SelectRandom     SelectionMode = 1
	SelectRoundRobin SelectionMode = 2
)

func (tester *Tester) GetClient(mode SelectionMode, input int) *txbuilder.Client {
	tester.selectionMutex.Lock()
	defer tester.selectionMutex.Unlock()
	switch mode {
	case SelectByIndex:
		input = input % len(tester.goodClients)
	case SelectRandom:
		input = rand.Intn(len(tester.goodClients))
	case SelectRoundRobin:
		input = tester.rrClientIdx
		tester.rrClientIdx++
		if tester.rrClientIdx >= len(tester.goodClients) {
			tester.rrClientIdx = 0
		}
	}
	return tester.goodClients[input]
}

func (tester *Tester) GetWallet(mode SelectionMode, input int) *txbuilder.Wallet {
	tester.selectionMutex.Lock()
	defer tester.selectionMutex.Unlock()
	switch mode {
	case SelectByIndex:
		input = input % len(tester.childWallets)
	case SelectRandom:
		input = rand.Intn(len(tester.childWallets))
	case SelectRoundRobin:
		input = tester.rrWalletIdx
		tester.rrWalletIdx++
		if tester.rrWalletIdx >= len(tester.childWallets) {
			tester.rrWalletIdx = 0
		}
	}
	return tester.childWallets[input]
}

func (tester *Tester) GetRootWallet() *txbuilder.Wallet {
	return tester.rootWallet
}
