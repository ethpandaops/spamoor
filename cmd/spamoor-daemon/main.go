package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

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
	disableAuditLogs bool
	secondsPerSlot   uint64
	auditUserHeader  string
	startupDelay     uint64
}

func main() {
	cliArgs := CliArgs{}
	flags := pflag.NewFlagSet("main", pflag.PanicOnError)

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
	flags.BoolVar(&cliArgs.disableAuditLogs, "disable-audit-logs", false, "Disable audit logs")
	flags.Uint64Var(&cliArgs.secondsPerSlot, "seconds-per-slot", 12, "Seconds per slot for rate limiting (used for throughput calculation).")
	flags.StringVar(&cliArgs.auditUserHeader, "audit-user-header", "Cf-Access-Authenticated-User-Email", "HTTP header containing the authenticated user email for audit logs")
	flags.Uint64Var(&cliArgs.startupDelay, "startup-delay", 30, "Delay in seconds before starting spammers on daemon startup (to allow cancellation)")
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
		ChainId:    clientPool.GetChainId(),
	})

	// init daemon
	spamoorDaemon = daemon.NewDaemon(ctx, logger.WithField("module", "daemon"), clientPool, txpool, database)
	if cliArgs.fuluActivation > 0 {
		spamoorDaemon.SetGlobalCfg("fulu_activation", cliArgs.fuluActivation)
	}
	if cliArgs.startupDelay > 0 {
		spamoorDaemon.SetStartupDelay(time.Duration(cliArgs.startupDelay) * time.Second)
	}

	// init audit logger
	auditLogger := daemon.NewAuditLogger(spamoorDaemon, cliArgs.auditUserHeader, "user")
	spamoorDaemon.SetAuditLogger(auditLogger)

	// start frontend
	webui.StartHttpServer(&types.FrontendConfig{
		Host:             "0.0.0.0",
		Port:             cliArgs.port,
		SiteName:         "Spamoor",
		Debug:            cliArgs.debug,
		Pprof:            true,
		Minify:           true,
		DisableTxMetrics: cliArgs.disableTxMetrics,
		DisableAuditLogs: cliArgs.disableAuditLogs,
	}, spamoorDaemon)

	// load and apply client configs from database
	err = spamoorDaemon.LoadAndApplyClientConfigs()
	if err != nil {
		logger.Warnf("failed to load client configs: %v", err)
	}

	// prepare clients after DB settings have been applied
	err = clientPool.PrepareClients()
	if err != nil {
		panic(fmt.Errorf("failed to prepare clients: %v", err))
	}

	// init root wallet
	rootWallet, err := spamoor.InitRootWallet(ctx, cliArgs.privkey, clientPool, txpool, logger)
	if err != nil {
		panic(fmt.Errorf("failed to init root wallet: %v", err))
	}
	defer rootWallet.Shutdown()

	if !cliArgs.withoutBatcher {
		rootWallet.InitTxBatcher(ctx, txpool)
	}

	spamoorDaemon.SetRootWallet(rootWallet)

	// start daemon & spammers
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

	// Shutdown components
	spamoorDaemon.Shutdown()
}
