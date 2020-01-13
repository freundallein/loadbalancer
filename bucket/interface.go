package bucket

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

// Available loadbalancing algorithms
const (
	RoundRobin = "round-robin"
)

// Server - common backend server interface
type Server interface {
	Address() *url.URL
	ReverseProxy() *httputil.ReverseProxy

	IsAvailable() bool
	SetAvailable(bool)

	LastSeen() int64

	PingServer() bool
}

// ServerBucket - common servers pool interface
type ServerBucket interface {
	AddServer(Server) error
	Size() int
	Serve(http.ResponseWriter, *http.Request) error
	Healthcheck()
	RemoveStale(time.Duration)
	RunServices(int)
}
