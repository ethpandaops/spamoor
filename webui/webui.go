package webui

import (
	"embed"
	"fmt"
	"net/http"
	"time"

	"github.com/ethpandaops/spamoor/daemon"
	"github.com/ethpandaops/spamoor/webui/handlers"
	"github.com/ethpandaops/spamoor/webui/handlers/api"
	"github.com/ethpandaops/spamoor/webui/handlers/docs"
	"github.com/ethpandaops/spamoor/webui/server"
	"github.com/ethpandaops/spamoor/webui/types"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"

	_ "net/http/pprof"
)

var (
	//go:embed static/*
	staticEmbedFS embed.FS

	//go:embed templates/*
	templateEmbedFS embed.FS
)

func StartHttpServer(config *types.FrontendConfig, daemon *daemon.Daemon) {
	// Initialize metrics collector
	err := daemon.InitializeMetrics()
	if err != nil {
		logrus.Errorf("failed to initialize metrics: %v", err)
	} else {
		logrus.Info("metrics endpoint available at /metrics")
	}

	// init router
	router := mux.NewRouter()

	frontend, err := server.NewFrontend(config, staticEmbedFS, templateEmbedFS)
	if err != nil {
		logrus.Fatalf("error initializing frontend: %v", err)
	}

	// register frontend routes
	frontendHandler := handlers.NewFrontendHandler(daemon)
	router.HandleFunc("/", frontendHandler.Index).Methods("GET")
	router.HandleFunc("/clients", frontendHandler.Clients).Methods("GET")
	router.HandleFunc("/wallets", frontendHandler.Wallets).Methods("GET")

	// API routes
	apiHandler := api.NewAPIHandler(daemon)
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/spammers", apiHandler.GetSpammerList).Methods("GET")
	apiRouter.HandleFunc("/scenarios", apiHandler.GetScenarios).Methods("GET")
	apiRouter.HandleFunc("/scenarios/{name}/config", apiHandler.GetScenarioConfig).Methods("GET")
	apiRouter.HandleFunc("/spammer", apiHandler.CreateSpammer).Methods("POST")
	apiRouter.HandleFunc("/spammer/{id}/start", apiHandler.StartSpammer).Methods("POST")
	apiRouter.HandleFunc("/spammer/{id}/pause", apiHandler.PauseSpammer).Methods("POST")
	apiRouter.HandleFunc("/spammer/{id}/reclaim", apiHandler.ReclaimFunds).Methods("POST")
	apiRouter.HandleFunc("/spammer/{id}", apiHandler.DeleteSpammer).Methods("DELETE")
	apiRouter.HandleFunc("/spammer/{id}/logs", apiHandler.GetSpammerLogs).Methods("GET")
	apiRouter.HandleFunc("/spammer/{id}", apiHandler.GetSpammerDetails).Methods("GET")
	apiRouter.HandleFunc("/spammer/{id}", apiHandler.UpdateSpammer).Methods("PUT")
	apiRouter.HandleFunc("/spammer/{id}/logs/stream", apiHandler.StreamSpammerLogs).Methods("GET")
	apiRouter.HandleFunc("/clients", apiHandler.GetClients).Methods("GET")
	apiRouter.HandleFunc("/client/{index}/group", apiHandler.UpdateClientGroup).Methods("PUT")
	apiRouter.HandleFunc("/client/{index}/enabled", apiHandler.UpdateClientEnabled).Methods("PUT")
	
	// Export/Import routes
	apiRouter.HandleFunc("/spammers/export", apiHandler.ExportSpammers).Methods("POST")
	apiRouter.HandleFunc("/spammers/import", apiHandler.ImportSpammers).Methods("POST")

	// metrics endpoint
	router.Handle("/metrics", promhttp.Handler()).Methods("GET")

	// swagger
	router.PathPrefix("/docs/").Handler(docs.GetSwaggerHandler(logrus.StandardLogger()))

	if config.Pprof {
		// add pprof handler
		router.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)
	}

	router.PathPrefix("/").Handler(frontend)

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	//n.Use(gzip.Gzip(gzip.DefaultCompression))
	n.UseHandler(router)

	if config.Host == "" {
		config.Host = "0.0.0.0"
	}
	if config.Port == 0 {
		config.Port = 8080
	}
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		WriteTimeout: 0,
		ReadTimeout:  0,
		IdleTimeout:  120 * time.Second,
		Handler:      n,
	}

	logrus.Printf("http server listening on %v", srv.Addr)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logrus.WithError(err).Fatal("Error serving frontend")
		}
	}()
}
