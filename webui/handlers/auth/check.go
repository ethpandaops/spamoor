package auth

import (
	"strings"

	"github.com/golang-jwt/jwt/v5"

	authpkg "github.com/ethpandaops/service-authenticatoor/pkg/auth"
)

// openToken is the sentinel returned in open mode. Callers only inspect
// `token != nil && token.Valid`, so an empty *jwt.Token with Valid set is
// enough for them to treat the request as authenticated.
var openToken = &jwt.Token{Valid: true}

// CheckAuthToken validates a bearer token (with or without the "Bearer "
// prefix) for a request bound to host. In open mode it always returns a
// valid token; in remote mode it delegates to the JWKS verifier and
// returns nil on any failure (signature, exp, iss, aud, scope/host,
// services). host should be the request's Host header stripped of any
// port — the verifier matches it against the token's "scope" claim.
func (h *Handler) CheckAuthToken(tokenStr, host string) *jwt.Token {
	if h.verifier == nil {
		return openToken
	}

	parts := strings.SplitN(tokenStr, " ", 2)
	if len(parts) == 2 && strings.EqualFold(parts[0], "bearer") {
		tokenStr = parts[1]
	}
	if tokenStr == "" {
		return nil
	}

	claims, err := h.verifier.Verify(tokenStr, authpkg.WithRequestHost(host))
	if err != nil {
		return nil
	}
	return &jwt.Token{Valid: true, Claims: claims}
}

// GetTokenSubject extracts the user identity (email) from a verified
// token. Returns "" when the token is missing/invalid or the handler is
// in open mode (no upstream identity to extract).
func (h *Handler) GetTokenSubject(authHeader, host string) string {
	if h.verifier == nil || authHeader == "" {
		return ""
	}

	token := h.CheckAuthToken(authHeader, host)
	if token == nil || !token.Valid {
		return ""
	}

	if c, ok := token.Claims.(*authpkg.Claims); ok {
		if c.Email != "" {
			return c.Email
		}
		return c.Subject
	}
	return ""
}

// StripPort removes a trailing ":port" from a request Host header value
// so it can be passed to CheckAuthToken. Handles bracketed IPv6 hosts.
func StripPort(host string) string {
	if host == "" {
		return host
	}
	if host[0] == '[' {
		if i := strings.IndexByte(host, ']'); i >= 0 {
			return host[1:i]
		}
		return host
	}
	if i := strings.LastIndexByte(host, ':'); i >= 0 {
		return host[:i]
	}
	return host
}
