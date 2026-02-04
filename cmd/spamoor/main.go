package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/scenarios"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/utils"
)

type CliArgs struct {
	verbose          bool
	trace            bool
	rpchosts         []string
	rpchostsFile     string
	privkey          string
	seed             string
	refillAmount     uint64
	refillBalance    uint64
	refillAmountWei  string
	refillBalanceWei string
	refillInterval   uint64
	slotDuration     time.Duration
	fundingGasLimit  uint64
	scenarioDir      string
	scenarioFile     string
}

func main() {
	// Check for subcommands
	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "run":
			RunCommand(os.Args[2:])
			return
		case "validate-scenario":
			ValidateCommand(os.Args[2:])
			return
		}
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
	flags.StringVar(&cliArgs.refillAmountWei, "refill-amount-wei", "", "Amount in Wei to fund each child wallet (overrides --refill-amount).")
	flags.StringVar(&cliArgs.refillBalanceWei, "refill-balance-wei", "", "Min balance in Wei before refilling (overrides --refill-balance).")
	flags.Uint64Var(&cliArgs.refillInterval, "refill-interval", 300, "Interval for child wallet rbalance check and refilling if needed (in sec).")
	flags.DurationVar(&cliArgs.slotDuration, "slot-duration", 12*time.Second, "Duration of a slot/block for rate limiting (e.g., '12s', '250ms'). Use sub-second values for L2 chains.")
	flags.Uint64Var(&cliArgs.fundingGasLimit, "funding-gas-limit", 21000, "Gas limit for wallet funding transactions (use 100000+ for L2s).")
	flags.StringVar(&cliArgs.scenarioDir, "scenario-dir", "", "Directory to load dynamic scenarios from (.go files).")
	flags.StringVar(&cliArgs.scenarioFile, "scenario-file", "", "Path to a single dynamic scenario file (.go file).")

	flags.Parse(os.Args)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// init logger
	logger := logrus.StandardLogger()

	logger.WithFields(logrus.Fields{
		"version":   utils.GetBuildVersion(),
		"buildtime": utils.BuildTime,
	}).Infof("starting spamoor")

	// Load dynamic scenarios if specified
	if cliArgs.scenarioDir != "" {
		scenarios.LoadDynamicScenarios(cliArgs.scenarioDir, logger)
	}
	if cliArgs.scenarioFile != "" {
		if err := scenarios.LoadDynamicScenarioFromFile(cliArgs.scenarioFile, logger); err != nil {
			logger.WithError(err).Fatalf("failed to load explicitly requested scenario file")
		}
	}

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

	// Set global slot duration for rate limiting
	scenario.GlobalSlotDuration = cliArgs.slotDuration

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

	clientOptions := []*spamoor.ClientOptions{}
	for _, rpcHost := range rpcHosts {
		clientOptions = append(clientOptions, &spamoor.ClientOptions{
			RpcHost: rpcHost,
		})
	}

	err := clientPool.InitClients(clientOptions)
	if err != nil {
		panic(fmt.Errorf("failed to init clients: %v", err))
	}

	err = clientPool.PrepareClients()
	if err != nil {
		panic(fmt.Errorf("failed to prepare clients: %v", err))
	}

	// prepare txpool
	txpool := spamoor.NewTxPool(&spamoor.TxPoolOptions{
		Context:    ctx,
		Logger:     logger.WithField("module", "txpool"),
		ClientPool: clientPool,
		ChainId:    clientPool.GetChainId(),
	})

	// init root wallet
	rootWallet, err := spamoor.InitRootWallet(ctx, cliArgs.privkey, clientPool, txpool, logger)
	if err != nil {
		panic(fmt.Errorf("failed to init root wallet: %v", err))
	}
	defer rootWallet.Shutdown()

	// init wallet pool
	walletPool := spamoor.NewWalletPool(ctx, logger.WithField("module", "walletpool"), rootWallet, clientPool, txpool)
	walletPool.SetWalletCount(100)

	// Set refill amount (Wei flag overrides ETH flag)
	if cliArgs.refillAmountWei != "" {
		amount, ok := new(big.Int).SetString(cliArgs.refillAmountWei, 10)
		if !ok {
			panic(fmt.Errorf("invalid refill-amount-wei value: %s", cliArgs.refillAmountWei))
		}
		walletPool.SetRefillAmount(uint256.MustFromBig(amount))
	} else {
		walletPool.SetRefillAmount(utils.EtherToWei(uint256.NewInt(cliArgs.refillAmount)))
	}

	// Set refill balance (Wei flag overrides ETH flag)
	if cliArgs.refillBalanceWei != "" {
		balance, ok := new(big.Int).SetString(cliArgs.refillBalanceWei, 10)
		if !ok {
			panic(fmt.Errorf("invalid refill-balance-wei value: %s", cliArgs.refillBalanceWei))
		}
		walletPool.SetRefillBalance(uint256.MustFromBig(balance))
	} else {
		walletPool.SetRefillBalance(utils.EtherToWei(uint256.NewInt(cliArgs.refillBalance)))
	}

	walletPool.SetRefillInterval(cliArgs.refillInterval)
	walletPool.SetWalletSeed(cliArgs.seed)
	walletPool.SetFundingGasLimit(cliArgs.fundingGasLimit)

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
