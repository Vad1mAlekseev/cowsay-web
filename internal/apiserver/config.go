package apiserver

import "time"

type Config struct {
	BindAddr          string
	StaticURLPrefix   string
	LogLevel          string
	ConnectionTimeout time.Duration `json:"connectionTimeoutNs,string,omitempty"`
}

func NewConfig() *Config {
	return &Config{
		":8080",
		"/static/",
		"debug",
		15 * time.Second,
	}
}
