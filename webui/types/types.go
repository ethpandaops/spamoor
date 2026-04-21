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

	AuthUserHeader    string
	AuthTokenKey      string
	DisableLocalToken bool
	DisableAuth       bool
}
