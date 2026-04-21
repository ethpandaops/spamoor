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

// getUserEmail extracts the user identity for audit logging.
// It first checks the JWT token subject, then falls back to the audit logger's header-based extraction.
func (ah *APIHandler) getUserEmail(r *http.Request) string {
	// Try to get the user from the JWT token subject
	if ah.authHandler != nil {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			authHeader = r.URL.Query().Get("auth")
			if authHeader != "" {
				authHeader = "Bearer " + authHeader
			}
		}

		if subject := ah.authHandler.GetTokenSubject(authHeader); subject != "" {
			return subject
		}
	}

	// Fall back to audit logger header-based extraction
	if auditLogger := ah.daemon.GetAuditLogger(); auditLogger != nil {
		return auditLogger.GetUserFromRequest(r.Header)
	}

	return "api"
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
