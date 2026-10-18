package webui

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v5"

	"github.com/ethpandaops/spamoor/webui/handlers"
	"github.com/ethpandaops/spamoor/webui/server"
	"github.com/ethpandaops/spamoor/webui/types"
)

// fakeChecker accepts exactly one token string (with or without the
// "Bearer " prefix). With accepted == "" it behaves like open mode and
// accepts everything.
type fakeChecker struct {
	accepted string
}

func (f *fakeChecker) CheckAuthToken(tokenStr, _ string) *jwt.Token {
	if f.accepted == "" {
		return &jwt.Token{Valid: true}
	}
	if tokenStr == f.accepted || tokenStr == "Bearer "+f.accepted {
		return &jwt.Token{Valid: true}
	}
	return nil
}

type auditCall struct {
	profile string
	status  int
}

func newTestPprofHandler(checker *fakeChecker, calls *[]auditCall) http.Handler {
	gate := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("gate"))
	})
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == pprofPathPrefix+"/missing" {
			http.Error(w, "unknown profile", http.StatusNotFound)
			return
		}
		_, _ = w.Write([]byte("profile:" + r.URL.Path))
	})
	audit := func(_ *http.Request, _ *jwt.Token, profile string, status int) {
		*calls = append(*calls, auditCall{profile: profile, status: status})
	}

	return newPprofHandler(checker, audit, gate, next)
}

func TestPprofHandlerAuth(t *testing.T) {
	tests := []struct {
		name       string
		checker    *fakeChecker
		method     string
		path       string
		accept     string
		header     string
		cookie     string
		wantStatus int
		wantBody   string
		wantWWW    bool
		wantAudit  []auditCall
	}{
		{
			name:       "open mode serves everything",
			checker:    &fakeChecker{},
			path:       "/debug/pprof/heap",
			wantStatus: http.StatusOK,
			wantBody:   "profile:/debug/pprof/heap",
			wantAudit:  []auditCall{{profile: "heap", status: http.StatusOK}},
		},
		{
			name:       "tool without token gets 401",
			checker:    &fakeChecker{accepted: "tok"},
			path:       "/debug/pprof/heap",
			accept:     "*/*",
			wantStatus: http.StatusUnauthorized,
			wantWWW:    true,
		},
		{
			name:       "browser without token gets the gate page",
			checker:    &fakeChecker{accepted: "tok"},
			path:       "/debug/pprof/",
			accept:     "text/html,application/xhtml+xml",
			wantStatus: http.StatusOK,
			wantBody:   "gate",
		},
		{
			name:       "browser POST without token gets 401, not the gate",
			checker:    &fakeChecker{accepted: "tok"},
			method:     http.MethodPost,
			path:       "/debug/pprof/heap",
			accept:     "text/html",
			wantStatus: http.StatusUnauthorized,
			wantWWW:    true,
		},
		{
			name:       "wrong bearer token gets 401",
			checker:    &fakeChecker{accepted: "tok"},
			path:       "/debug/pprof/heap",
			header:     "Bearer nope",
			wantStatus: http.StatusUnauthorized,
			wantWWW:    true,
		},
		{
			name:       "bearer header serves the profile and audits it",
			checker:    &fakeChecker{accepted: "tok"},
			path:       "/debug/pprof/goroutine",
			header:     "Bearer tok",
			wantStatus: http.StatusOK,
			wantBody:   "profile:/debug/pprof/goroutine",
			wantAudit:  []auditCall{{profile: "goroutine", status: http.StatusOK}},
		},
		{
			name:       "cookie serves the profile and audits it",
			checker:    &fakeChecker{accepted: "tok"},
			path:       "/debug/pprof/heap",
			cookie:     "tok",
			wantStatus: http.StatusOK,
			wantBody:   "profile:/debug/pprof/heap",
			wantAudit:  []auditCall{{profile: "heap", status: http.StatusOK}},
		},
		{
			name:       "wrong cookie gets the gate page in a browser",
			checker:    &fakeChecker{accepted: "tok"},
			path:       "/debug/pprof/heap",
			accept:     "text/html",
			cookie:     "stale",
			wantStatus: http.StatusOK,
			wantBody:   "gate",
		},
		{
			name:       "index page is served without an audit entry",
			checker:    &fakeChecker{accepted: "tok"},
			path:       "/debug/pprof/",
			header:     "Bearer tok",
			wantStatus: http.StatusOK,
			wantBody:   "profile:/debug/pprof/",
		},
		{
			name:       "index page without trailing slash is not audited either",
			checker:    &fakeChecker{accepted: "tok"},
			path:       "/debug/pprof",
			header:     "Bearer tok",
			wantStatus: http.StatusOK,
			wantBody:   "profile:/debug/pprof",
		},
		{
			name:       "failed download is audited with its status",
			checker:    &fakeChecker{accepted: "tok"},
			path:       "/debug/pprof/missing",
			header:     "Bearer tok",
			wantStatus: http.StatusNotFound,
			wantAudit:  []auditCall{{profile: "missing", status: http.StatusNotFound}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := make([]auditCall, 0, 1)
			handler := newTestPprofHandler(tt.checker, &calls)

			method := tt.method
			if method == "" {
				method = http.MethodGet
			}
			req := httptest.NewRequest(method, tt.path, nil)
			if tt.accept != "" {
				req.Header.Set("Accept", tt.accept)
			}
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}
			if tt.cookie != "" {
				req.AddCookie(&http.Cookie{Name: pprofTokenCookie, Value: tt.cookie})
			}

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d (body %q)", rec.Code, tt.wantStatus, rec.Body.String())
			}
			if tt.wantBody != "" && rec.Body.String() != tt.wantBody {
				t.Errorf("body = %q, want %q", rec.Body.String(), tt.wantBody)
			}
			if tt.wantWWW && rec.Header().Get("WWW-Authenticate") != "Bearer" {
				t.Errorf("WWW-Authenticate = %q, want Bearer", rec.Header().Get("WWW-Authenticate"))
			}
			if tt.wantBody == "gate" && rec.Header().Get("Cache-Control") != "no-store" {
				t.Errorf("Cache-Control = %q, want no-store", rec.Header().Get("Cache-Control"))
			}
			if tt.wantAudit == nil {
				if len(calls) != 0 {
					t.Errorf("audit calls = %v, want none", calls)
				}
			} else if !reflect.DeepEqual(calls, tt.wantAudit) {
				t.Errorf("audit calls = %v, want %v", calls, tt.wantAudit)
			}
		})
	}
}

func TestPprofHandlerPrefersHeaderOverCookie(t *testing.T) {
	calls := make([]auditCall, 0, 1)
	handler := newTestPprofHandler(&fakeChecker{accepted: "tok"}, &calls)

	// A stale header must not be rescued by a valid cookie: the header is
	// what the caller asserted, so it is what gets checked.
	req := httptest.NewRequest(http.MethodGet, "/debug/pprof/heap", nil)
	req.Header.Set("Authorization", "Bearer stale")
	req.AddCookie(&http.Cookie{Name: pprofTokenCookie, Value: "tok"})

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
	if len(calls) != 0 {
		t.Errorf("audit calls = %v, want none", calls)
	}
}

func TestPprofLoginGateRenders(t *testing.T) {
	_, err := server.NewFrontend(&types.FrontendConfig{SiteName: "Spamoor", AuthProviderURL: "https://auth.example.test"}, staticEmbedFS, templateEmbedFS)
	if err != nil {
		t.Fatalf("failed to init frontend: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/debug/pprof/", nil)
	rec := httptest.NewRecorder()
	handlers.NewFrontendHandler(nil).PprofLogin(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()
	for _, want := range []string{"<title>Profiling - Spamoor", "pprof-login-required", "Path=' + COOKIE_PATH", "client-v2.js"} {
		if !strings.Contains(body, want) {
			t.Errorf("gate page is missing %q", want)
		}
	}
}
