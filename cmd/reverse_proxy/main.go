package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ReverseProxy struct {
	backend url.URL
	proxy   *httputil.ReverseProxy
}

func printHeaders(label string, headers http.Header) {
	fmt.Printf("\n=== %s ===\n", label)
	for key, values := range headers {
		for _, value := range values {
			fmt.Printf("%s: %s\n", key, value)
		}
	}
	fmt.Println("==========================================")
}

func NewReverseProxy(backendURL string) (*ReverseProxy, error) {
	parsedURL, err := url.Parse(backendURL)
	if err != nil {
		return nil, err
	}

	internalProxy := httputil.NewSingleHostReverseProxy(parsedURL)
	internalProxy.Director = nil

	internalProxy.Rewrite = func(proxyreq *httputil.ProxyRequest) {
		printHeaders("Before", proxyreq.In.Header)
		proxyreq.SetURL(parsedURL)

		proxyreq.Out.Header.Set("X-Forwarded-Host", proxyreq.In.Host)
		proxyreq.Out.Header.Set("X-Forwarded-Proto", "http")

		printHeaders("After", proxyreq.Out.Header)
	}

	reverseproxy := &ReverseProxy{
		backend: *parsedURL,
		proxy:   internalProxy,
	}
	return reverseproxy, nil
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

	log.Printf("proxying :8080 %s", backend)
	if err := http.ListenAndServe(":8080", proxy); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
