package types

type FrontendConfig struct {
	Host             string
	Port             int
	SiteName         string
	Debug            bool
	Pprof            bool
	Minify           bool
	DisableTxMetrics bool
	DisableAuditLogs bool
	DisablePluginAPI bool

	// AuthProviderURL is the canonical URL of a remote authenticatoor
	// service. When empty, the API runs unauthenticated; when set, all
	// protected endpoints require a JWT verified against the service's
	// JWKS, and the layout template injects window.spamoorConfig +
	// loads <url>/client.js so the frontend can drive login.
	AuthProviderURL string
}
