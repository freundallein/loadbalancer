package httpserv

import (
	"fmt"
	"net/http"
	"time"

	"github.com/freundallein/loadbalancer/bucket"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	bucketSize = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "lb_bucket_size",
		Help: "The total number of servers in bucket",
	})
)

func collectCustomMetrics(buck bucket.ServerBucket) {
	for {
		select {
		case <-time.After(5 * time.Second):
			bucketSize.Set(float64(buck.Size()))
		}
	}
}

// New - http server constructor
func New(port int, bckt bucket.ServerBucket) *http.Server {
	go collectCustomMetrics(bckt)
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", LoadBalance(bckt))
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}
	return server
}
