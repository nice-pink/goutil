package network

import (
	"encoding/json"
	"errors"
	"io"
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

	if c.token != "" {
		// has token
		return nil
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

func (c *Client) Request(method, url string, body io.Reader, headers Headers, authRequired bool) (*http.Response, error) {
	if c.verbose {
		log.Verbose(strings.ToUpper(method), url)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Err(err, "Could not create request.", method, url)
		return nil, err
	}

	if authRequired {
		err := c.RefreshToken()
		if err != nil {
			return nil, err
		}
	}

	// set headers
	c.addHeaders(req, headers, authRequired)

	// request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		if c.verbose {
			log.Err(err, "Could not send request.")
		}
		return nil, err
	}

	if authRequired && resp.StatusCode == http.StatusUnauthorized {
		if c.verbose {
			log.Info("not authorized -> clear token.", req.URL, resp.StatusCode)
		}
		c.ClearToken()
		// recursive call after clearing token
		return c.Request(method, url, body, headers, authRequired)
	}

	return resp, nil
}

// convenience

func (c *Client) RequestData(method, url string, body io.Reader, headers Headers, authRequired bool) ([]byte, error) {
	resp, err := c.Request(method, url, body, headers, authRequired)
	if err != nil {
		if c.verbose {
			log.Err(err, "response error", url)
		}
		return nil, err
	}
	defer resp.Body.Close()

	// read data
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		if c.verbose {
			log.Err(err, "read body", url)
		}
		return nil, err
	}

	return data, err
}

func (c *Client) RequestMap(method, url string, body io.Reader, headers Headers, authRequired bool) (map[string]any, error) {
	data, err := c.RequestData(method, url, body, headers, authRequired)
	if err != nil {
		if c.verbose {
			log.Err(err, "response error", url)
		}
		return nil, err
	}

	// unmarshal to map
	var m map[string]any
	err = json.Unmarshal(data, &m)
	if err != nil {
		if c.verbose {
			log.Err(err, "unmarshal error", url, string(data))
		}
		return nil, err
	}
	return m, err
}

func (c *Client) RequestType(method, url string, body io.Reader, headers Headers, authRequired bool, output any) error {
	data, err := c.RequestData(method, url, body, headers, authRequired)
	if err != nil {
		if c.verbose {
			log.Err(err, "response error", url)
		}
		return err
	}

	// unmarshal to map
	err = json.Unmarshal(data, &output)
	if err != nil {
		if c.verbose {
			log.Err(err, "unmarshal error", url, string(data))
		}
		return err
	}
	return err
}

// intern

func (c *Client) addHeaders(req *http.Request, headers Headers, authRequired bool) {
	// bearer token
	if authRequired {
		if c.token != "" {
			req.Header.Add("Authorization", "Bearer "+c.token)
		}
		if c.basicAuth != "" {
			req.Header.Add("Authorization", "Basic "+c.basicAuth)
		}
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
