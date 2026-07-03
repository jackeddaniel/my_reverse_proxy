package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Let's use the client")
	tr := &http.Transport{
		IdleConnTimeout: 30 * time.Second,
	}
	client := &http.Client{Timeout: time.Duration(1) * time.Second, Transport: tr}
	resp, err := client.Get("http://localhost:3490/test.html")

	if err != nil {
		panic(err)
	}

	log.Println(resp)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Printf("Body: %s", body)

}
