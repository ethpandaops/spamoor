package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenarios"
	"github.com/ethpandaops/spamoor/scenariotypes"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/utils"
)

type CliArgs struct {
	verbose        bool
	trace          bool
	rpchosts       []string
	rpchostsFile   string
	privkey        string
	seed           string
	refillAmount   uint64
	refillBalance  uint64
	refillInterval uint64
}

func main() {
	cliArgs := CliArgs{}
	flags := pflag.NewFlagSet("main", pflag.ContinueOnError)

	flags.BoolVarP(&cliArgs.verbose, "verbose", "v", false, "Run the script with verbose output")
	flags.BoolVar(&cliArgs.trace, "trace", false, "Run the script with tracing output")
	flags.StringArrayVarP(&cliArgs.rpchosts, "rpchost", "h", []string{}, "The RPC host to send transactions to.")
	flags.StringVar(&cliArgs.rpchostsFile, "rpchost-file", "", "File with a list of RPC hosts to send transactions to.")
	flags.StringVarP(&cliArgs.privkey, "privkey", "p", "", "The private key of the wallet to send funds from.")
	flags.StringVarP(&cliArgs.seed, "seed", "s", "", "The child wallet seed.")
	flags.Uint64Var(&cliArgs.refillAmount, "refill-amount", 5, "Amount of ETH to fund/refill each child wallet with.")
	flags.Uint64Var(&cliArgs.refillBalance, "refill-balance", 2, "Min amount of ETH each child wallet should hold before refilling.")
	flags.Uint64Var(&cliArgs.refillInterval, "refill-interval", 300, "Interval for child wallet rbalance check and refilling if needed (in sec).")

	flags.Parse(os.Args)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// init logger
	if cliArgs.trace {
		logrus.SetLevel(logrus.TraceLevel)
	} else if cliArgs.verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logger := logrus.StandardLogger()

	logger.WithFields(logrus.Fields{
		"version":   utils.GetBuildVersion(),
		"buildtime": utils.BuildTime,
	}).Infof("starting spamoor")

	// load scenario
	invalidScenario := false
	var scenarioName string
	var scenarioDescriptior *scenariotypes.ScenarioDescriptor
	if flags.NArg() < 2 {
		invalidScenario = true
	} else {
		scenarioName = flags.Args()[1]
		scenarioDescriptior = scenarios.GetScenario(scenarioName)
		if scenarioDescriptior == nil {
			invalidScenario = true
		}
	}
	if invalidScenario {
		fmt.Printf("invalid or missing scenario\n\n")
		fmt.Printf("implemented scenarios:\n")
		scenarioNames := scenarios.GetScenarioNames()
		sort.Slice(scenarioNames, func(a int, b int) bool {
			return strings.Compare(scenarioNames[a], scenarioNames[b]) > 0
		})
		for _, name := range scenarioNames {
			fmt.Printf("  %v\n", name)
		}
		return
	}

	scenario := scenarioDescriptior.NewScenario(logger)
	if scenario == nil {
		panic("could not create scenario instance")
	}

	flags.Init(fmt.Sprintf("%v %v", flags.Args()[0], scenarioName), pflag.ExitOnError)
	scenario.Flags(flags)
	cliArgs.rpchosts = nil
	flags.Parse(os.Args)

	// start client pool
	rpcHosts := []string{}
	for _, rpcHost := range strings.Split(strings.Join(cliArgs.rpchosts, ","), ",") {
		if rpcHost != "" {
			rpcHosts = append(rpcHosts, rpcHost)
		}
	}

	if cliArgs.rpchostsFile != "" {
		fileLines, err := utils.ReadFileLinesTrimmed(cliArgs.rpchostsFile)
		if err != nil {
			panic(err)
		}
		rpcHosts = append(rpcHosts, fileLines...)
	}

	clientPool := spamoor.NewClientPool(ctx, rpcHosts, logger.WithField("module", "clientpool"))
	err := clientPool.PrepareClients()
	if err != nil {
		panic(fmt.Errorf("failed to prepare clients: %v", err))
	}

	// init root wallet
	rootWallet, err := spamoor.InitRootWallet(ctx, cliArgs.privkey, clientPool.GetClient(spamoor.SelectClientRandom, 0), logger)
	if err != nil {
		panic(fmt.Errorf("failed to init root wallet: %v", err))
	}

	// prepare txpool
	txpool := txbuilder.NewTxPool(&txbuilder.TxPoolOptions{
		GetClientFn: func(index int, random bool) *txbuilder.Client {
			mode := spamoor.SelectClientByIndex
			if random {
				mode = spamoor.SelectClientRandom
			}

			return clientPool.GetClient(mode, index)
		},
		GetClientCountFn: func() int {
			return len(clientPool.GetAllClients())
		},
	})

	// init wallet pool
	walletPool := spamoor.NewWalletPool(ctx, logger.WithField("module", "walletpool"), rootWallet, clientPool, txpool)
	walletPool.SetWalletCount(100)
	walletPool.SetRefillAmount(utils.EtherToWei(uint256.NewInt(cliArgs.refillAmount)))
	walletPool.SetRefillBalance(utils.EtherToWei(uint256.NewInt(cliArgs.refillBalance)))
	walletPool.SetRefillInterval(cliArgs.refillInterval)
	walletPool.SetWalletSeed(cliArgs.seed)

	// init scenario
	err = scenario.Init(walletPool, "")
	if err != nil {
		panic(err)
	}

	// prepare wallet pool
	err = walletPool.PrepareWallets(true)
	if err != nil {
		panic(fmt.Errorf("failed to prepare wallets: %v", err))
	}

	// start scenario
	err = scenario.Run(ctx)
	if err != nil {
		panic(err)
	}
}
