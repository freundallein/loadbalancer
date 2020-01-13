package bucket

import (
	"errors"
	"net/http/httputil"
	"net/url"
	"time"
)

var (
	ErrInvalidAlgorithm = errors.New("invalid balancing algorithm chosen.")
)

// NewServer - backend server factory
func NewServer(URL string) (Server, error) {
	addr, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}
	reverseProxy := httputil.NewSingleHostReverseProxy(addr)
	return &DefaultServer{
		address:      addr,
		isAvailable:  true,
		reverseProxy: reverseProxy,
		lastSeen:     time.Now().Unix(),
	}, nil
}

// New - backends pool factory, can use different balancing algorithms in future
func New(algo string) (ServerBucket, error) {
	var bckt ServerBucket
	switch algo {
	case RoundRobin:
		bckt = &RoundRobinServerBucket{
			servers: []Server{},
		}
	}
	if bckt == nil {
		return nil, ErrInvalidAlgorithm
	}
	return bckt, nil
}
