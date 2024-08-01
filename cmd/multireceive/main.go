package main

import (
	"flag"
	"io"
	"net/http"
	"os"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/nice-pink/goutil/pkg/log"
)

var (
	TIMEOUT         time.Duration
	ENDPOINTS       []string
	ENDPOINTS_ACKED []string
)

func main() {
	endpoints := flag.String("endpoints", "", "Commaseparated endpoints to watch on.")
	method := flag.String("method", "get", "Accepted method.")
	timeout := flag.Int("timeout", 600, "Timeout for accept")
	flag.Parse()

	TIMEOUT = time.Duration(*timeout) * time.Second
	ENDPOINTS = strings.Split(*endpoints, ",")
	sort.Strings(ENDPOINTS)

	for _, e := range ENDPOINTS {
		if strings.ToLower(*method) == "get" {
			log.Info("endpoint listening. get:", e)
			http.HandleFunc("/"+e, handleGet)
		} else if strings.ToLower(*method) == "post" {
			log.Info("endpoint listening. post:", e)
			http.HandleFunc("/"+e, handlePost)
		}
	}

	s := &http.Server{
		Addr:           ":8080",
		ReadTimeout:    TIMEOUT,
		WriteTimeout:   TIMEOUT,
		IdleTimeout:    TIMEOUT,
		MaxHeaderBytes: 1 << 20,
	}

	s.SetKeepAlivesEnabled(false)

	// start server, panic on error
	err := s.ListenAndServe()
	if err != nil {
		log.Err(err, "listen and serve error")
		os.Exit(2)
	}

	log.Info("Hooray")
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	ackEndpoint(r.URL.Path)

	io.WriteString(w, "OK")
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	io.WriteString(w, "OK")
}

func ackEndpoint(e string) bool {
	ack := strings.TrimPrefix(e, "/")
	for _, endpoint := range ENDPOINTS_ACKED {
		if ack == endpoint {
			return true
		}
	}

	log.Info("received request from:", ack)
	ENDPOINTS_ACKED = append(ENDPOINTS_ACKED, ack)
	sort.Strings(ENDPOINTS_ACKED)

	if allAcked() {
		log.Info("All endpoints acked.")
		os.Exit(0)
	}
	return true
}

func allAcked() bool {
	return slices.Equal[[]string](ENDPOINTS, ENDPOINTS_ACKED)
}
