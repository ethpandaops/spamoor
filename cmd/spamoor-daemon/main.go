package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/daemon"
	"github.com/ethpandaops/spamoor/daemon/db"
	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/utils"
	"github.com/ethpandaops/spamoor/webui"
	"github.com/ethpandaops/spamoor/webui/types"
)

type CliArgs struct {
	verbose          bool
	trace            bool
	debug            bool
	rpchosts         []string
	rpchostsFile     string
	privkey          string
	port             int
	dbFile           string
	startupSpammer   string
	fuluActivation   uint64
	withoutBatcher   bool
	disableTxMetrics bool
	secondsPerSlot   uint64
}

func main() {
	cliArgs := CliArgs{}
	flags := pflag.NewFlagSet("main", pflag.ContinueOnError)

	flags.BoolVarP(&cliArgs.verbose, "verbose", "v", false, "Run the tool with verbose output")
	flags.BoolVar(&cliArgs.trace, "trace", false, "Run the tool with tracing output")
	flags.BoolVar(&cliArgs.debug, "debug", false, "Run the tool in debug mode")
	flags.StringArrayVarP(&cliArgs.rpchosts, "rpchost", "h", []string{}, "The RPC host to send transactions to.")
	flags.StringVar(&cliArgs.rpchostsFile, "rpchost-file", "", "File with a list of RPC hosts to send transactions to.")
	flags.StringVarP(&cliArgs.privkey, "privkey", "p", "", "The private key of the wallet to send funds from.")
	flags.IntVarP(&cliArgs.port, "port", "P", 8080, "The port to run the webui on.")
	flags.StringVarP(&cliArgs.dbFile, "db", "d", "spamoor.db", "The file to store the database in.")
	flags.StringVar(&cliArgs.startupSpammer, "startup-spammer", "", "YAML file or URL with startup spammers configuration")
	flags.Uint64Var(&cliArgs.fuluActivation, "fulu-activation", 0, "The unix timestamp of the Fulu activation (if activated)")
	flags.BoolVar(&cliArgs.withoutBatcher, "without-batcher", false, "Run the tool without batching funding transactions")
	flags.BoolVar(&cliArgs.disableTxMetrics, "disable-tx-metrics", false, "Disable transaction metrics collection and graphs page (keeps Prometheus metrics)")
	flags.Uint64Var(&cliArgs.secondsPerSlot, "seconds-per-slot", 12, "Seconds per slot for rate limiting (used for throughput calculation).")
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
	}).Infof("starting spamoor daemon")

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

	clientPool := spamoor.NewClientPool(ctx, rpcHosts, logger.WithField("module", "clientpool"))
	err := clientPool.PrepareClients()
	if err != nil {
		panic(fmt.Errorf("failed to prepare clients: %v", err))
	}

	// init root wallet
	client := clientPool.GetClient(spamoor.SelectClientRandom, 0, "")
	if client == nil {
		panic(fmt.Errorf("no client available"))
	}

	rootWallet, err := spamoor.InitRootWallet(ctx, cliArgs.privkey, client, logger)
	if err != nil {
		panic(fmt.Errorf("failed to init root wallet: %v", err))
	}

	// init db
	database := db.NewDatabase(&db.SqliteDatabaseConfig{
		File: cliArgs.dbFile,
	}, logger.WithField("module", "db"))
	err = database.Init()
	if err != nil {
		panic(fmt.Errorf("failed to init db: %v", err))
	}

	err = database.ApplyEmbeddedDbSchema(-2)
	if err != nil {
		panic(fmt.Errorf("failed to apply db schema: %v", err))
	}

	// prepare txpool
	var spamoorDaemon *daemon.Daemon

	txpool := spamoor.NewTxPool(&spamoor.TxPoolOptions{
		Context:    ctx,
		Logger:     logger.WithField("module", "txpool"),
		ClientPool: clientPool,
		GetActiveWalletPools: func() []*spamoor.WalletPool {
			walletPools := make([]*spamoor.WalletPool, 0)
			for _, sp := range spamoorDaemon.GetAllSpammers() {
				walletPool := sp.GetWalletPool()
				if walletPool != nil {
					walletPools = append(walletPools, walletPool)
				}
			}
			return walletPools
		},
	})

	if !cliArgs.withoutBatcher {
		rootWallet.InitTxBatcher(ctx, txpool)
	}

	// init daemon
	spamoorDaemon = daemon.NewDaemon(ctx, logger.WithField("module", "daemon"), clientPool, rootWallet, txpool, database)
	if cliArgs.fuluActivation > 0 {
		spamoorDaemon.SetGlobalCfg("fulu_activation", cliArgs.fuluActivation)
	}

	// start frontend
	webui.StartHttpServer(&types.FrontendConfig{
		Host:             "0.0.0.0",
		Port:             cliArgs.port,
		SiteName:         "Spamoor",
		Debug:            cliArgs.debug,
		Pprof:            true,
		Minify:           true,
		DisableTxMetrics: cliArgs.disableTxMetrics,
	}, spamoorDaemon)

	// start daemon
	firstLaunch, err := spamoorDaemon.Run()
	if err != nil {
		panic(err)
	}

	// load startup spammers if configured
	if firstLaunch && cliArgs.startupSpammer != "" {
		err := spamoorDaemon.ImportSpammersOnStartup(cliArgs.startupSpammer, logger.WithField("module", "startup"))
		if err != nil {
			logger.Errorf("failed to import startup spammers: %v", err)
		}
	}

	// wait for ctrl+c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	spamoorDaemon.Shutdown()
}
