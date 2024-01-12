package network

type Auth struct {
	// either bearer token (preferred) or basic auth!
	BasicUser     string
	BasicPassword string
	BearerToken   string
}

type RequestConfig struct {
	auth   Auth
	accept string // e.g. "application/json"
}
