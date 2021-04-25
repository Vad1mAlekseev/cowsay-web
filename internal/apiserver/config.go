package apiserver

import "time"

// The Config of APIServer.
type Config struct {
	BindAddr          string
	StaticURLPrefix   string
	LogLevel          string
	ConnectionTimeout time.Duration `json:"connectionTimeoutNs,string,omitempty"`
	certPath          string
	keyPath           string
}

// NewConfig get default config.
func NewConfig() *Config {
	return &Config{
		":8080",
		"/static/",
		"debug",
		15 * time.Second,
		"/www/cert.cert",
		"/www/key",
	}
}
