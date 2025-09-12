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

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/scenarios"
	"github.com/ethpandaops/spamoor/spamoor"
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
	secondsPerSlot uint64
}

func main() {
	// Check if "run" subcommand is used
	if len(os.Args) >= 2 && os.Args[1] == "run" {
		RunCommand(os.Args[2:])
		return
	}

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
	flags.Uint64Var(&cliArgs.secondsPerSlot, "seconds-per-slot", 12, "Seconds per slot for rate limiting (used for throughput calculation).")

	flags.Parse(os.Args)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// init logger
	logger := logrus.StandardLogger()

	logger.WithFields(logrus.Fields{
		"version":   utils.GetBuildVersion(),
		"buildtime": utils.BuildTime,
	}).Infof("starting spamoor")

	// load scenario
	invalidScenario := false
	var scenarioName string
	var scenarioDescriptior *scenario.Descriptor
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
		fmt.Printf("Usage:\n")
		fmt.Printf("  spamoor <scenario> [options]     Run a specific scenario\n")
		fmt.Printf("  spamoor run <yaml-file> [options] Run multiple scenarios from YAML config\n\n")
		fmt.Printf("Implemented scenarios:\n")
		scenarioNames := scenarios.GetScenarioNames()
		sort.Slice(scenarioNames, func(a int, b int) bool {
			return strings.Compare(scenarioNames[a], scenarioNames[b]) > 0
		})
		for _, name := range scenarioNames {
			fmt.Printf("  %v\n", name)
		}
		return
	}

	newScenario := scenarioDescriptior.NewScenario(logger)
	if newScenario == nil {
		panic("could not create scenario instance")
	}

	flags.Init(fmt.Sprintf("%v %v", flags.Args()[0], scenarioName), pflag.ExitOnError)
	newScenario.Flags(flags)
	cliArgs.rpchosts = nil
	flags.Parse(os.Args)

	if cliArgs.trace {
		logrus.SetLevel(logrus.TraceLevel)
	} else if cliArgs.verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}

	// Set global seconds per slot
	scenario.GlobalSecondsPerSlot = cliArgs.secondsPerSlot

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

	clientPool := spamoor.NewClientPool(ctx, logger.WithField("module", "clientpool"))

	err := clientPool.InitClients(rpcHosts)
	if err != nil {
		panic(fmt.Errorf("failed to init clients: %v", err))
	}

	err = clientPool.PrepareClients()
	if err != nil {
		panic(fmt.Errorf("failed to prepare clients: %v", err))
	}

	// init root wallet
	rootWallet, err := spamoor.InitRootWallet(ctx, cliArgs.privkey, clientPool, logger)
	if err != nil {
		panic(fmt.Errorf("failed to init root wallet: %v", err))
	}
	defer rootWallet.Shutdown()

	// prepare txpool
	var walletPool *spamoor.WalletPool

	txpool := spamoor.NewTxPool(&spamoor.TxPoolOptions{
		Context:    ctx,
		Logger:     logger.WithField("module", "txpool"),
		ClientPool: clientPool,
		GetActiveWalletPools: func() []*spamoor.WalletPool {
			return []*spamoor.WalletPool{walletPool}
		},
	})

	// init wallet pool
	walletPool = spamoor.NewWalletPool(ctx, logger.WithField("module", "walletpool"), rootWallet, clientPool, txpool)
	walletPool.SetWalletCount(100)
	walletPool.SetRefillAmount(utils.EtherToWei(uint256.NewInt(cliArgs.refillAmount)))
	walletPool.SetRefillBalance(utils.EtherToWei(uint256.NewInt(cliArgs.refillBalance)))
	walletPool.SetRefillInterval(cliArgs.refillInterval)
	walletPool.SetWalletSeed(cliArgs.seed)

	// init scenario
	err = newScenario.Init(&scenario.Options{
		WalletPool: walletPool,
	})
	if err != nil {
		panic(err)
	}

	// prepare wallet pool
	err = walletPool.PrepareWallets()
	if err != nil {
		panic(fmt.Errorf("failed to prepare wallets: %v", err))
	}

	// start scenario
	err = newScenario.Run(ctx)
	if err != nil {
		panic(err)
	}
}
