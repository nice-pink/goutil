package network

import "time"

type Auth struct {
	// either bearer token (preferred) or basic auth!
	BasicUser     string
	BasicPassword string
	BearerToken   string
}

type RequestConfig struct {
	LogLevel int
	Auth     Auth
	Accept   string // e.g. "application/json"
	Timeout  time.Duration
	MaxBytes int64 // stop after reading bytes
}

func DefaultRequestConfig() RequestConfig {
	return RequestConfig{Timeout: 15.0, MaxBytes: -1}
}
