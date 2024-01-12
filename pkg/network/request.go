package network

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/nice-pink/goutil/pkg/log"
)

type Requester struct {
	config RequestConfig
}

func NewRequester(config RequestConfig) *Requester {
	return &Requester{config: config}
}

// request

func (r *Requester) Get(url string, printBody bool) ([]byte, error) {
	// request
	resp, err := r.Request(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// read and return
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Err(err, "Read all error.")
		return nil, err
	}

	if printBody {
		log.Info(string(body))
	}
	return body, err
}

func (r *Requester) Delete(url string) (bool, error) {
	// request
	resp, err := r.Request(http.MethodDelete, url, nil)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// read and return
	if resp.StatusCode != 200 && resp.StatusCode != 202 {
		log.Error("Could not delete. Status code:", strconv.Itoa(resp.StatusCode), "Url:", url)
		return false, nil
	}
	log.Info("Success! Deleted:", url)
	return true, nil
}

func (r *Requester) Request(method string, url string, body io.Reader) (*http.Response, error) {
	// build request
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Err(err, "Request error.")
		return nil, err
	}

	// auth
	if r.config.Auth.BearerToken != "" {
		var bearer = "Bearer " + r.config.Auth.BearerToken
		req.Header.Add("Authorization", bearer)
	} else if r.config.Auth.BasicUser != "" && r.config.Auth.BasicPassword != "" {
		// add basic auth
		req.SetBasicAuth(r.config.Auth.BasicUser, r.config.Auth.BasicPassword)
	}

	// header
	if r.config.Accept != "" {
		req.Header.Add("Accept", r.config.Accept)
	}

	// request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Err(err, "Client error.")
		return nil, err
	}

	return resp, err
}
