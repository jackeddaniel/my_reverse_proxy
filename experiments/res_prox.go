package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ReverseProxy struct {
	backend url.URL
	current uint64
}

func NewReverseProxy(backendURL string) *ReverseProxy {
	parsedURL, err := url.Parse(backendURL)
	if err != nil {
		fmt.Println("Wrong URL/ broken URL")
		return nil
	}

	return &ReverseProxy{backend: *parsedURL}
}

func (p *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("the req is ", r)
	proxy := httputil.NewSingleHostReverseProxy(&p.backend)
	proxy.ServeHTTP(w, r)
}
func main() {
	backend := "http://localhost:3490/"
	proxy := NewReverseProxy(backend)
	fmt.Println("The reverse proxy struct is is ", proxy)
	http.ListenAndServe(":8080", proxy)
}
