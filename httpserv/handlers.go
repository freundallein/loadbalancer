package httpserv

import (
	"net/http"

	"../bucket"
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
