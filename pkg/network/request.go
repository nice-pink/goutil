package network

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/nice-pink/goutil/pkg/log"
)

type StreamInfo struct {
	Url       string
	BytesRead uint64
}

type Requester struct {
	config     RequestConfig
	streamInfo StreamInfo
}

func NewRequester(config RequestConfig) *Requester {
	return &Requester{config: config}
}

// request

func (r *Requester) Get(url string, printBody bool) ([]byte, error) {
	if r.config.LogLevel > 0 {
		log.Info("Get:", url)
	}

	// request
	resp, err := r.Request(http.MethodGet, url, false, nil)
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
	if r.config.LogLevel > 0 {
		log.Info("Delete:", url)
	}

	// request
	resp, err := r.Request(http.MethodDelete, url, false, nil)
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

// stream

func (r *Requester) ReadStream(url string, dumpToFile string) error {
	if r.config.LogLevel > 0 {
		log.Info("Read stream:", url)
	}

	r.streamInfo.Url = url

	// request
	resp, err := r.Request(http.MethodGet, url, true, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// open file
	writeToFile := false
	var file *os.File = nil
	if dumpToFile != "" {
		file, err = os.Create(dumpToFile)
		writeToFile = true
		defer func() {
			if err := file.Close(); err != nil {
				log.Err(err, "Could not close file.")
			}
		}()
	}

	// read data
	reader := bufio.NewReader(resp.Body)
	for {
		line, _ := reader.ReadBytes('\n')
		if writeToFile {
			file.Write(line)
		}
		r.streamInfo.BytesRead += uint64(len(line))

		if r.config.MaxBytes > 0 && r.streamInfo.BytesRead > uint64(r.config.MaxBytes) {
			log.Info("Stop: Max bytes read", r.streamInfo.BytesRead)
			break
		}
	}

	return err
}

func (r *Requester) PrintStreamInfo() {
	log.Info()
	log.Info("Url:", r.streamInfo.Url)
	log.Info("Bytes read:", r.streamInfo.BytesRead)
}

// common

func (r *Requester) Request(method string, url string, isStream bool, body io.Reader) (*http.Response, error) {
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

		if r.config.LogLevel > 2 {
			log.Info("Bearer auth:", bearer)
		}
	} else if r.config.Auth.BasicUser != "" && r.config.Auth.BasicPassword != "" {
		// add basic auth
		req.SetBasicAuth(r.config.Auth.BasicUser, r.config.Auth.BasicPassword)

		if r.config.LogLevel > 2 {
			log.Info("Basic auth:", r.config.Auth.BasicUser, r.config.Auth.BasicPassword)
		}
	}

	// header
	if r.config.Accept != "" {
		req.Header.Add("Accept", r.config.Accept)

		if r.config.LogLevel > 1 {
			log.Info("Accept:", r.config.Accept)
		}
	}

	// request
	client := &http.Client{Timeout: r.config.Timeout * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Err(err, "Client error.")
		return nil, err
	}

	return resp, err
}
