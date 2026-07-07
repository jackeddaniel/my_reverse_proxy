package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ReverseProxy struct {
	backend url.URL
	proxy   *httputil.ReverseProxy
}

func NewReverseProxy(backendURL string) (*ReverseProxy, error) {
	parsedURL, err := url.Parse(backendURL)
	if err != nil {
		return nil, err
	}

	return &ReverseProxy{
		backend: *parsedURL,
		proxy:   httputil.NewSingleHostReverseProxy(parsedURL),
	}, nil
}

func (p *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s -> %s", r.Method, r.URL.Path, p.backend.Host)
	p.proxy.ServeHTTP(w, r)
}

func main() {
	backend := "http://localhost:3490/"

	proxy, err := NewReverseProxy(backend)
	if err != nil {
		log.Fatalf("failed to create reverse proxy: %v", err)
	}

	log.Printf("proxying :8080", backend)
	if err := http.ListenAndServe(":8080", proxy); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
