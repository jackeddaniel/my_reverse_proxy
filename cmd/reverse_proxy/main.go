package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type ReverseProxy struct {
	backend url.URL
	proxy   *httputil.ReverseProxy
}

type LogEntry struct {
	method    string
	path      string
	status    string
	latency   time.Time
	backend   string
	client_ip string
	trace_id  string
}

func generateID() string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "unknown-id"
	}

	return hex.EncodeToString(bytes)
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

		currIP, _, err := net.SplitHostPort(proxyreq.In.RemoteAddr)

		fmt.Println("this is what the remoteaddr stores", proxyreq.In.RemoteAddr)

		if err != nil {
			currIP = proxyreq.In.RemoteAddr
		}

		currXFF := proxyreq.In.Header.Get("X-Forwarded-Host")

		if currXFF != "" {
			currXFF = currXFF + ", " + currIP
		} else {
			currXFF = currIP
		}

		reqId := proxyreq.In.Header.Get("Request-ID")
		if reqId == "" {
			reqId = generateID()
		}

		proxyreq.Out.Header.Set("Request-ID", reqId)
		proxyreq.Out.Header.Set("X-Forwarded-For", currXFF)
		proxyreq.Out.Header.Set("X-Forwarded-Host", proxyreq.In.Host)
		proxyreq.Out.Header.Set("X-Forwarded-Proto", "http")

		printHeaders("After", proxyreq.Out.Header)
	}

	internalProxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Set("Server", "Siddu's proxy")
		resp.Header.Set("Trace-id", generateID())

		return nil
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
