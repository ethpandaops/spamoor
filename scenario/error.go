package scenario

import "errors"

var (
	ErrNoClients = errors.New("no clients available")
	ErrNoWallet  = errors.New("no wallet available")
)
