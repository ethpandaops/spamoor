package api

import (
	"net/http"

	"github.com/ethpandaops/spamoor/daemon"
	"github.com/ethpandaops/spamoor/webui/handlers/auth"
)

type APIHandler struct {
	daemon      *daemon.Daemon
	authHandler *auth.Handler
}

func NewAPIHandler(d *daemon.Daemon, authHandler *auth.Handler) *APIHandler {
	return &APIHandler{
		daemon:      d,
		authHandler: authHandler,
	}
}

// checkAuth verifies the Authorization header and returns true if authenticated.
// If not authenticated, it writes an error response and returns false.
// Also checks for auth query parameter for SSE connections.
func (ah *APIHandler) checkAuth(w http.ResponseWriter, r *http.Request) bool {
	if ah.authHandler == nil {
		return true // No auth handler configured, allow all
	}

	// First try Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		// Fall back to query parameter for SSE connections
		authHeader = r.URL.Query().Get("auth")
		if authHeader != "" {
			authHeader = "Bearer " + authHeader
		}
	}

	if authHeader == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return false
	}

	token := ah.authHandler.CheckAuthToken(authHeader)
	if token == nil || !token.Valid {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return false
	}

	return true
}

// isAuthenticated checks if the request has valid authentication without returning an error.
// Use this for endpoints that are public but need to conditionally show data.
func (ah *APIHandler) isAuthenticated(r *http.Request) bool {
	if ah.authHandler == nil {
		return true // No auth handler configured, consider authenticated
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return false
	}

	token := ah.authHandler.CheckAuthToken(authHeader)
	return token != nil && token.Valid
}
