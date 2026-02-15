package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Handler handles authentication requests for the assertoor web UI.
type Handler struct {
	userHeader        string
	tokenKey          string
	disableLocalToken bool
}

// NewAuthHandler creates a new authentication handler.
func NewAuthHandler(tokenKey, userHeader string, disableLocalToken bool) *Handler {
	return &Handler{
		userHeader:        userHeader,
		tokenKey:          tokenKey,
		disableLocalToken: disableLocalToken,
	}
}

// GetToken handles the authentication request.
func (h *Handler) GetToken(w http.ResponseWriter, r *http.Request) {
	if h.disableLocalToken {
		http.Error(w, "local token is disabled", http.StatusForbidden)
		return
	}

	headers := r.Header
	authUser := "unauthenticated"

	// Try exact header match first
	if values, ok := headers[h.userHeader]; ok && len(values) > 0 {
		authUser = values[0]
	} else {
		// Try case-insensitive match
		for key, values := range headers {
			if strings.EqualFold(key, h.userHeader) && len(values) > 0 {
				authUser = values[0]
				break
			}
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "spamoor",
		Subject:   authUser,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
	})

	tokenString, err := token.SignedString([]byte(h.tokenKey))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	claims, ok := token.Claims.(jwt.RegisteredClaims)
	if !ok {
		http.Error(w, "invalid token claims", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
		"user":  authUser,
		"expr":  fmt.Sprintf("%d", claims.ExpiresAt.Unix()),
		"now":   fmt.Sprintf("%d", time.Now().Unix()),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

// GetLogin redirects to the index page.
func (h *Handler) GetLogin(w http.ResponseWriter, r *http.Request) {
	// redirect back to the index page
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
