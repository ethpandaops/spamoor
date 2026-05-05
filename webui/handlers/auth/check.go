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
// prefix). In open mode it always returns a valid token; in remote mode
// it delegates to the JWKS verifier and returns nil on any failure.
func (h *Handler) CheckAuthToken(tokenStr string) *jwt.Token {
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

	claims, err := h.verifier.Verify(tokenStr)
	if err != nil {
		return nil
	}
	return &jwt.Token{Valid: true, Claims: claims}
}

// GetTokenSubject extracts the user identity (email) from a verified
// token. Returns "" when the token is missing/invalid or the handler is
// in open mode (no upstream identity to extract).
func (h *Handler) GetTokenSubject(authHeader string) string {
	if h.verifier == nil || authHeader == "" {
		return ""
	}

	token := h.CheckAuthToken(authHeader)
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
