package network

import (
	"errors"
	"net/http"
	"strings"

	"github.com/nice-pink/goutil/pkg/log"
)

type Headers map[string]string

type AuthFn func() (string, error)

type Client struct {
	verbose       bool
	sharedHeaders Headers
	token         string
	basicAuth     string
	authFn        AuthFn
	httpClient    *http.Client
}

// sharedHeaders: will be added to all requests
// basicAuth: username:password
// token: bearer token
// authFn: function to get a token
func NewClient(sharedHeaders Headers, token, basicAuth string, authFn AuthFn, verbose bool) *Client {
	return &Client{
		verbose:       verbose,
		sharedHeaders: sharedHeaders,
		token:         token,
		basicAuth:     basicAuth,
		httpClient:    &http.Client{},
		authFn:        authFn,
	}
}

func (c *Client) ClearToken() {
	c.token = ""
}

func (c *Client) RefreshToken() error {
	if c.basicAuth != "" {
		// use basic auth
		return nil
	}

	// use token
	if c.authFn == nil && c.token == "" {
		// token can't be generated
		return errors.New("no auth function")
	}

	// get token
	token, err := c.authFn()
	if err != nil {
		return err
	}
	c.token = token
	return nil
}

// request

func (c *Client) Request(req *http.Request, headers Headers, authRequired bool) (*http.Response, error) {
	if c.verbose {
		log.Verbose(strings.ToUpper(req.Method), req.URL)
	}

	if authRequired {
		err := c.RefreshToken()
		if err != nil {
			return nil, err
		}
	}

	// set headers
	c.addHeaders(req, headers)

	// request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Err(err, "Could not send request.")
		return nil, err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		if c.verbose {
			log.Info("not authorized -> clear token.", req.URL, resp.StatusCode)
		}
		c.ClearToken()
		// recursive call after clearing token
		return c.Request(req, headers, authRequired)
	}

	return resp, nil
}

func (c *Client) addHeaders(req *http.Request, headers Headers) {
	// bearer token
	if c.token != "" {
		req.Header.Add("Authorization", "Bearer "+c.token)
	}
	// shared headers
	for k, v := range c.sharedHeaders {
		req.Header.Add(k, v)
	}
	// additional headers
	for k, v := range headers {
		req.Header.Add(k, v)
	}
}
