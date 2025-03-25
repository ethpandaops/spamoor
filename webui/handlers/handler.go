package handlers

import (
	"github.com/ethpandaops/spamoor/daemon"
)

type FrontendHandler struct {
	daemon *daemon.Daemon
}

func NewFrontendHandler(d *daemon.Daemon) *FrontendHandler {
	return &FrontendHandler{
		daemon: d,
	}
}
