package api

import (
	"github.com/ethpandaops/spamoor/daemon"
)

type APIHandler struct {
	daemon *daemon.Daemon
}

func NewAPIHandler(d *daemon.Daemon) *APIHandler {
	return &APIHandler{
		daemon: d,
	}
}
