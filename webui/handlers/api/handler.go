package api

import (
	"encoding/json"
	"net/http"

	"github.com/ethpandaops/spamoor/daemon"
	"github.com/ethpandaops/spamoor/webui/handlers/auth"
)

type APIHandler struct {
	daemon          *daemon.Daemon
	authHandler     *auth.Handler
	authProviderURL string
}

func NewAPIHandler(d *daemon.Daemon, authHandler *auth.Handler, authProviderURL string) *APIHandler {
	return &APIHandler{
		daemon:          d,
		authHandler:     authHandler,
		authProviderURL: authProviderURL,
	}
}

// runtimeConfigResponse is the body of GET /api/runtime-config. It tells
// the frontend which auth provider (if any) to load.
type runtimeConfigResponse struct {
	AuthProviderURL string `json:"authProviderURL"`
}

// GetRuntimeConfig serves runtime configuration that the frontend reads
// at boot to wire up the auth provider client. Public — no auth required.
func (ah *APIHandler) GetRuntimeConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	_ = json.NewEncoder(w).Encode(runtimeConfigResponse{
		AuthProviderURL: ah.authProviderURL,
	})
}

// checkAuth verifies the Authorization header and returns true if
// authenticated. In open mode (no auth provider configured) it always
// returns true. On failure it writes a 401 and returns false. The "auth"
// query parameter is honored as a fallback for SSE connections.
func (ah *APIHandler) checkAuth(w http.ResponseWriter, r *http.Request) bool {
	if ah.authHandler.IsOpen() {
		return true
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		if q := r.URL.Query().Get("auth"); q != "" {
			authHeader = "Bearer " + q
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

// getUserEmail extracts the user identity for audit logging directly from
// the verified bearer token. In open mode there's no token to read from,
// so it returns the default "api".
func (ah *APIHandler) getUserEmail(r *http.Request) string {
	if ah.authHandler.IsOpen() {
		return "api"
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		if q := r.URL.Query().Get("auth"); q != "" {
			authHeader = "Bearer " + q
		}
	}
	if subject := ah.authHandler.GetTokenSubject(authHeader); subject != "" {
		return subject
	}
	return "api"
}

// isAuthenticated reports whether the request carries a valid token.
// Used by endpoints that are public but conditionally redact fields.
// In open mode every request is considered authenticated.
func (ah *APIHandler) isAuthenticated(r *http.Request) bool {
	if ah.authHandler.IsOpen() {
		return true
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return false
	}

	token := ah.authHandler.CheckAuthToken(authHeader)
	return token != nil && token.Valid
}
