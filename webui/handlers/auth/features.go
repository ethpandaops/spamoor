package auth

import (
	"github.com/golang-jwt/jwt/v5"

	authpkg "github.com/ethpandaops/service-authenticatoor/pkg/auth"
)

// Features of spamoor. A token needs an explicit grant for each
// ("spamoor:write", "spamoor:all", "all:all"); basic access to spamoor is
// read-only.
const (
	// FeatureWrite: create, change, start, pause, delete and import spammers
	// and spammer groups, and reclaim their funds.
	FeatureWrite = "write"
	// FeaturePlugins: register, reload and delete plugins (code that runs
	// inside spamoor).
	FeaturePlugins = "plugins"
	// FeatureWallet: send transactions from the root wallet.
	FeatureWallet = "wallet"
	// FeatureClients: change the client pool (group, enabled, name, type),
	// which affects every spammer.
	FeatureClients = "clients"
	// FeaturePprof: the pprof pages, whose profiles expose process memory.
	FeaturePprof = "pprof"
)

// FeatureAllowed reports whether a token returned by CheckAuthToken may use a
// privileged feature. Open mode allows everything. A token without a
// "services" claim comes from an authenticatoor without grants and keeps full
// access, as before features existed; one with grants never mints an empty claim.
func (h *Handler) FeatureAllowed(token *jwt.Token, feature string) bool {
	if h.verifier == nil {
		return true
	}
	if token == nil || !token.Valid {
		return false
	}

	claims, ok := token.Claims.(*authpkg.Claims)
	if !ok {
		return false
	}
	if claims.Services == "" {
		return true
	}

	return authpkg.FeatureAllowed(claims.Services, serviceName, feature, authpkg.FeatureDefaultDeny)
}
