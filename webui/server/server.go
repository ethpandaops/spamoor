package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ethpandaops/spamoor/daemon"
	"github.com/ethpandaops/spamoor/webui"
	"github.com/ethpandaops/spamoor/webui/handlers"
	"github.com/ethpandaops/spamoor/webui/handlers/api"
	"github.com/ethpandaops/spamoor/webui/handlers/docs"
	"github.com/ethpandaops/spamoor/webui/types"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

func StartHttpServer(config *types.FrontendConfig, daemon *daemon.Daemon) {
	// init router
	router := mux.NewRouter()

	frontend, err := webui.NewFrontend(config)
	if err != nil {
		logrus.Fatalf("error initializing frontend: %v", err)
	}

	// register frontend routes
	frontendHandler := handlers.NewFrontendHandler(daemon)
	router.HandleFunc("/", frontendHandler.Index).Methods("GET")
	router.HandleFunc("/health", frontendHandler.Health).Methods("GET")
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
	apiRouter.HandleFunc("/spammer/{id}", apiHandler.DeleteSpammer).Methods("DELETE")
	apiRouter.HandleFunc("/spammer/{id}/logs", apiHandler.GetSpammerLogs).Methods("GET")
	apiRouter.HandleFunc("/spammer/{id}", apiHandler.GetSpammerDetails).Methods("GET")
	apiRouter.HandleFunc("/spammer/{id}", apiHandler.UpdateSpammer).Methods("PUT")
	apiRouter.HandleFunc("/spammer/{id}/logs/stream", apiHandler.StreamSpammerLogs).Methods("GET")

	// swagger
	router.PathPrefix("/docs/").Handler(docs.GetSwaggerHandler(logrus.StandardLogger()))

	router.PathPrefix("/").Handler(frontend)

	if config.Pprof {
		// add pprof handler
		router.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)
	}

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
