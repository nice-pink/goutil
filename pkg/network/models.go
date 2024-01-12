package network

type Auth struct {
	// either bearer token (preferred) or basic auth!
	BasicUser     string
	BasicPassword string
	BearerToken   string
}

type RequestConfig struct {
	Auth   Auth
	Accept string // e.g. "application/json"
}
