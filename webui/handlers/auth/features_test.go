package auth

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"

	authpkg "github.com/ethpandaops/service-authenticatoor/pkg/auth"
)

func TestFeatureAllowed(t *testing.T) {
	enforcing := &Handler{verifier: &authpkg.JWKSVerifier{}}
	token := func(services string) *jwt.Token {
		return &jwt.Token{Valid: true, Claims: &authpkg.Claims{Services: services}}
	}

	tests := []struct {
		name    string
		handler *Handler
		token   *jwt.Token
		feature string
		want    bool
	}{
		{"open mode allows everything", &Handler{}, openToken, FeaturePlugins, true},
		{"no token", enforcing, nil, FeaturePlugins, false},
		{"invalid token", enforcing, &jwt.Token{Claims: &authpkg.Claims{}}, FeaturePlugins, false},
		{"token without claims type", enforcing, &jwt.Token{Valid: true}, FeaturePlugins, false},
		{"no services claim: authenticatoor without grants", enforcing, token(""), FeaturePlugins, true},
		{"staff all:all", enforcing, token("+all,+all:all"), FeatureWallet, true},
		{"basic spamoor only", enforcing, token("-all,+spamoor"), FeaturePlugins, false},
		{"basic spamoor only, pprof", enforcing, token("-all,+spamoor"), FeaturePprof, false},
		{"granted feature", enforcing, token("-all,+spamoor,+spamoor:clients"), FeatureClients, true},
		{"other feature not granted", enforcing, token("-all,+spamoor,+spamoor:clients"), FeatureWallet, false},
		{"spamoor:all", enforcing, token("-all,+spamoor,+spamoor:all"), FeaturePlugins, true},
		{"spamoor:all minus plugins", enforcing, token("-all,+spamoor,+spamoor:all:-plugins"), FeaturePlugins, false},
		{"spamoor:all minus plugins, wallet", enforcing, token("-all,+spamoor,+spamoor:all:-plugins"), FeatureWallet, true},
		{"faucet only", enforcing, token("-all,+faucet,+faucet:drop_100eth"), FeatureWallet, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.handler.FeatureAllowed(tt.token, tt.feature); got != tt.want {
				t.Errorf("FeatureAllowed(%q) = %v, want %v", tt.feature, got, tt.want)
			}
		})
	}
}
