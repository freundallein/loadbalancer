package main

import (
	"fmt"
	"net/http"
	"os"
)

// Healthchecker for docker container
// We build image from scratch and dont have curl
func main() {
	response, err := http.Get(fmt.Sprintf("http://127.0.0.1:%s/healthz", os.Getenv("PORT")))
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil || (response != nil && response.StatusCode >= 500) {
		os.Exit(1)
	}
}
