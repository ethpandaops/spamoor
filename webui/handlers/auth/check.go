package auth

import (
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (h *Handler) CheckAuthToken(tokenStr string) *jwt.Token {
	// Extract token from "Bearer <token>"
	parts := strings.SplitN(tokenStr, " ", 2)
	if len(parts) == 2 && strings.EqualFold(parts[0], "bearer") {
		tokenStr = parts[1]
	}

	token, _ := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(h.tokenKey), nil
	})

	return token
}

// GetTokenSubject extracts the subject claim from the authorization header's JWT token.
// Returns an empty string if the token is missing, invalid, or has no subject.
func (h *Handler) GetTokenSubject(authHeader string) string {
	if authHeader == "" {
		return ""
	}

	token := h.CheckAuthToken(authHeader)
	if token == nil || !token.Valid {
		return ""
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || claims.Subject == "" {
		return ""
	}

	return claims.Subject
}
