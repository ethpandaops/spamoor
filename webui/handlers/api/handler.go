package api

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"

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

// checkAuth verifies the Authorization header and returns true if
// authenticated. In open mode (no auth provider configured) it always
// returns true. On failure it writes a 401 and returns false. The "auth"
// query parameter is honored as a fallback for SSE connections.
func (ah *APIHandler) checkAuth(w http.ResponseWriter, r *http.Request) bool {
	_, ok := ah.authToken(w, r)
	return ok
}

// checkFeature is checkAuth for a privileged feature (auth.Feature*): a
// valid token without a grant for it gets a 403.
func (ah *APIHandler) checkFeature(w http.ResponseWriter, r *http.Request, feature string) bool {
	token, ok := ah.authToken(w, r)
	if !ok {
		return false
	}

	if !ah.authHandler.FeatureAllowed(token, feature) {
		sendError(w, "this login has no grant for spamoor:"+feature, http.StatusForbidden)
		return false
	}

	return true
}

// authToken returns the verified token of the request, or writes a 401 and
// returns false.
func (ah *APIHandler) authToken(w http.ResponseWriter, r *http.Request) (*jwt.Token, bool) {
	if ah.authHandler.IsOpen() {
		return ah.authHandler.CheckAuthToken("", ""), true
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		if q := r.URL.Query().Get("auth"); q != "" {
			authHeader = "Bearer " + q
		}
	}

	if authHeader == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return nil, false
	}

	token := ah.authHandler.CheckAuthToken(authHeader, auth.StripPort(r.Host))
	if token == nil || !token.Valid {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return nil, false
	}

	return token, true
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
	if subject := ah.authHandler.GetTokenSubject(authHeader, auth.StripPort(r.Host)); subject != "" {
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

	token := ah.authHandler.CheckAuthToken(authHeader, auth.StripPort(r.Host))
	return token != nil && token.Valid
}
