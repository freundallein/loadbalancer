package httpserv

import (
	"net/http"

	"github.com/freundallein/loadbalancer/bucket"
)

// LoadBalance - main handler, use servers pool to serve requests
func LoadBalance(buck bucket.ServerBucket) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := buck.Serve(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
		}
	}
}

// Healthz - service healthcheck handler
func Healthz(buck bucket.ServerBucket) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if buck.Size() > 0 {
			w.WriteHeader(200)
			w.Write([]byte("ok"))
		} else {
			w.WriteHeader(500)
        	w.Write([]byte("error: no servers available"))
		}
	}
}
