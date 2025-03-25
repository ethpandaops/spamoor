package types

type FrontendConfig struct {
	Host     string
	Port     int
	SiteName string
	Debug    bool
	Pprof    bool
	Minify   bool
}
