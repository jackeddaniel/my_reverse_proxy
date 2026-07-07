package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type ProxyServer struct {
	proxy *httputil.ReverseProxy
}

func NewProxyServer(targetURL string) (*ProxyServer, error) {
	destination, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(destination)

	originalDirector := reverseProxy.Director
	reverseProxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = destination.Host
	}

	return &ProxyServer{proxy: reverseProxy}, nil
}

func (p *ProxyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("[PROXY] %s %s -> %s", r.Method, r.URL.Path, r.RemoteAddr)
	p.proxy.ServeHTTP(w, r)
}

func main() {
	const proxyPort = ":8080"
	const backendTarget = "http://localhost:3490"

	server, err := NewProxyServer(backendTarget)
	if err != nil {
		log.Fatalf("Failed to initialize proxy layout: %v", err)
	}

	httpServer := &http.Server{
		Addr:         proxyPort,
		Handler:      server,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Proxying traffic from http://localhost%s to %s\n", proxyPort, backendTarget)
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
