package httpserv

import (
	"fmt"
	"net/http"

	"freundallein/loadbalancer/bucket"
)

// New - http server constructor
func New(port int, bckt bucket.ServerBucket) *http.Server {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(LoadBalance(bckt)),
	}
	return server
}
