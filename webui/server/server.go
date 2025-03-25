package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ethpandaops/spamoor/daemon"
	"github.com/ethpandaops/spamoor/webui"
	"github.com/ethpandaops/spamoor/webui/handlers"
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
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/scenarios", frontendHandler.GetScenarios).Methods("GET")
	api.HandleFunc("/scenarios/{name}/config", frontendHandler.GetScenarioConfig).Methods("GET")
	api.HandleFunc("/spammer", frontendHandler.CreateSpammer).Methods("POST")
	api.HandleFunc("/spammer/{id}/start", frontendHandler.StartSpammer).Methods("POST")
	api.HandleFunc("/spammer/{id}/pause", frontendHandler.PauseSpammer).Methods("POST")
	api.HandleFunc("/spammer/{id}", frontendHandler.DeleteSpammer).Methods("DELETE")
	api.HandleFunc("/spammer/{id}/logs", frontendHandler.GetSpammerLogs).Methods("GET")
	api.HandleFunc("/spammer/{id}", frontendHandler.GetSpammerDetails).Methods("GET")
	api.HandleFunc("/spammer/{id}", frontendHandler.UpdateSpammer).Methods("PUT")
	api.HandleFunc("/spammer/{id}/logs/stream", frontendHandler.StreamSpammerLogs).Methods("GET")

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
