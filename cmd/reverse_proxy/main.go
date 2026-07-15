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
	backends []url.URL
	proxy    *httputil.ReverseProxy
	index    int
}

type logEntry struct {
	Method   string `json:"method"`
	Path     string `json:"path"`
	Status   int    `json:"status"`
	Latency  string `json:"latency"`
	Backend  string `json:"backend"`
	ClientIP string `json:"client_ip"`
	TraceID  string `json:"trace_id"`
}

type statusRecorder struct {
	http.ResponseWriter
	status int
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

func (p *ReverseProxy) nextBackend() *url.URL {
	backend := &p.backends[p.index]
	p.index = (p.index + 1) % len(p.backends)
	return backend
}

func (r *statusRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func NewReverseProxy(backends []string) (*ReverseProxy, error) {
	reverseproxy := &ReverseProxy{
		backends: make([]url.URL, 0, len(backends)),
	}

	for i, backend := range backends {
		parsedURL, err := url.Parse(backend)
		if err != nil {
			fmt.Println("Failed to append the following backend: ", i, backend)
		}

		reverseproxy.backends = append(reverseproxy.backends, *parsedURL)
	}

	reverseproxy.index = 0

	//internalProxy := httputil.NewSingleHostReverseProxy(parsedURL)
	var internalProxy httputil.ReverseProxy

	internalProxy.Rewrite = func(proxyreq *httputil.ProxyRequest) {
		//printHeaders("Before", proxyreq.In.Header)
		proxyreq.SetURL(reverseproxy.nextBackend())

		currIP, _, err := net.SplitHostPort(proxyreq.In.RemoteAddr)

		//fmt.Println("this is what the remoteaddr stores", proxyreq.In.RemoteAddr)

		if err != nil {
			currIP = proxyreq.In.RemoteAddr
		}

		currXFF := proxyreq.In.Header.Get("X-Forwarded-For")

		if currXFF != "" {
			currXFF = currXFF + ", " + currIP
		} else {
			currXFF = currIP
		}

		reqId := proxyreq.In.Header.Get("Request-ID")
		if reqId == "" {
			reqId = generateID()

			proxyreq.In.Header.Set("Request-ID", reqId)
		}

		proxyreq.Out.Header.Set("Request-ID", reqId)
		proxyreq.Out.Header.Set("X-Forwarded-For", currXFF)
		proxyreq.Out.Header.Set("X-Forwarded-Host", proxyreq.In.Host)
		proxyreq.Out.Header.Set("X-Forwarded-Proto", "http")

		// now we set next index
		reverseproxy.index = (reverseproxy.index + 1) % len(reverseproxy.backends)

	}

	internalProxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Set("Server", "Siddu's proxy")

		if resp.Request != nil {
			reqID := resp.Request.Header.Get("Request-ID")

			if reqID != "" {
				resp.Header.Set("Trace-ID", reqID)
			}
		}

		return nil
	}

	reverseproxy.proxy = &internalProxy

	return reverseproxy, nil
}

func (p *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

	p.proxy.ServeHTTP(recorder, r)

	duration := time.Since(startTime)
	latencyStr := fmt.Sprintf("%.2fms", float64(duration.Nanoseconds())/1e6)

	clientIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		clientIP = r.RemoteAddr
	}

	entry := logEntry{
		Method:   r.Method,
		Path:     r.URL.Path,
		Status:   recorder.status,
		Latency:  latencyStr,
		Backend:  p.backends[p.index].String(),
		ClientIP: clientIP,
		TraceID:  r.Header.Get("Request-ID"),
	}

	jsonData, err := json.Marshal(entry)
	if err == nil {
		fmt.Println(string(jsonData))
	} else {
		log.Printf("Failed to marshal log entry: %v", err)
	}
}

func main() {
	config, err := LoadConfig("../../config.yaml")
	if err != nil {
		fmt.Println("Config read failed")
		return
	}

	fmt.Println("This is the backend: ", config.Backends[0])
	fmt.Println("This is the config: ", *config)

	proxy, err := NewReverseProxy(config.Backends)
	if err != nil {
		log.Fatalf("failed to create reverse proxy: %v", err)
	}

	readTimeout, err := time.ParseDuration(config.Timeouts.Read)
	if err != nil {
		log.Fatalf("invalid read timeout %q: %v", config.Timeouts.Read, err)
	}

	server := &http.Server{
		Addr:    config.Port,
		Handler: proxy, ReadTimeout: readTimeout,
	}

	log.Printf("proxying :8080 %s", config.Backends)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
