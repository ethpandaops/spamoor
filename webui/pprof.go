package webui

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"

	"github.com/ethpandaops/spamoor/daemon"
	"github.com/ethpandaops/spamoor/webui/handlers/auth"
)

// pprofPathPrefix is where net/http/pprof registers itself on
// http.DefaultServeMux (via the blank import in webui.go).
const pprofPathPrefix = "/debug/pprof"

// pprofTokenCookie carries the login token for the /debug/pprof pages in a
// browser. The gate page sets it from the session (path-scoped to
// pprofPathPrefix, SameSite=Strict, expiring with the token); see
// templates/debug/pprof_login.html.
const pprofTokenCookie = "spamoor_debug_token"

// pprofTokenChecker is the part of auth.Handler the pprof gate needs; an
// interface so tests can stand in a fake verifier.
type pprofTokenChecker interface {
	CheckAuthToken(tokenStr, host string) *jwt.Token
}

// pprofAuditFunc records a profile download. profile is the path below
// pprofPathPrefix ("heap", "goroutine", "profile", ...), status the HTTP
// status the pprof handler answered with.
type pprofAuditFunc func(r *http.Request, token *jwt.Token, profile string, status int)

// newPprofHandler serves the pprof pages to requests that carry a valid
// login token, as a bearer header (tools: curl, go tool pprof) or the
// pprofTokenCookie (a browser); always valid when no auth provider is
// configured. A profile both exposes the process internals (RPC
// credentials, private keys, ...) and costs CPU to collect, so every
// profile download lands in the audit log via audit; only the index page
// with its counts does not. A browser without a valid token gets the gate
// page, which sets the cookie from the session and reloads, or asks for a
// login; everything else is a 401.
func newPprofHandler(checker pprofTokenChecker, audit pprofAuditFunc, gate, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			if cookie, err := r.Cookie(pprofTokenCookie); err == nil {
				header = cookie.Value
			}
		}

		token := checker.CheckAuthToken(header, auth.StripPort(r.Host))
		if token != nil && token.Valid {
			profile := strings.Trim(strings.TrimPrefix(r.URL.Path, pprofPathPrefix), "/")
			if profile == "" {
				next.ServeHTTP(w, r)

				return
			}

			recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(recorder, r)

			if audit != nil {
				audit(r, token, profile, recorder.status)
			}

			return
		}

		if r.Method == http.MethodGet && strings.Contains(r.Header.Get("Accept"), "text/html") {
			w.Header().Set("Cache-Control", "no-store")
			gate.ServeHTTP(w, r)

			return
		}

		w.Header().Set("WWW-Authenticate", "Bearer")
		http.Error(w, "login required", http.StatusUnauthorized)
	})
}

// pprofAuditor returns the pprofAuditFunc that writes profile downloads to
// the daemon's audit log. Like the API handlers, it attributes requests
// without an identity (open mode) to "api".
func pprofAuditor(d *daemon.Daemon) pprofAuditFunc {
	return func(r *http.Request, token *jwt.Token, profile string, status int) {
		auditLogger := d.GetAuditLogger()
		if auditLogger == nil {
			return
		}

		userEmail := auth.TokenSubject(token)
		if userEmail == "" {
			userEmail = "api"
		}

		result := "ok"
		if status >= http.StatusBadRequest {
			result = fmt.Sprintf("error: HTTP %d", status)
		}

		err := auditLogger.LogPprofDownload(userEmail, profile, r.URL.RawQuery, result)
		if err != nil {
			logrus.WithError(err).WithField("profile", profile).Warn("failed to audit pprof download")
		}
	}
}

// statusRecorder captures the status code a handler wrote.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

// Flush keeps streaming responses (CPU profile, trace) flushable.
func (s *statusRecorder) Flush() {
	if flusher, ok := s.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

var _ http.Flusher = (*statusRecorder)(nil)
