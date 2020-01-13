package bucket

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

const (
	healthCheckPeriod = 5 * time.Second
	removeStalePeriod = 60 * time.Second
	maxRetries        = 3
	maxAttempts       = 3
)

var (
	ErrInvalidServer         = errors.New("expected Server, got nil")
	ErrNoServersAvailable    = errors.New("no servers are available")
	ErrAllServersUnreachable = errors.New("all servers are unreachable")
	ErrServiceUnavailable    = errors.New("service not available")
)

// RoundRobinServerBucket - round-robin representatino of servers pool
type RoundRobinServerBucket struct {
	servers []Server     // servers storage
	last    uint64       // last used server index
	lock    sync.RWMutex // lock for servers slice
}

// AddServer - collect Server instance
func (sb *RoundRobinServerBucket) AddServer(srv Server) error {
	if srv == nil {
		return ErrInvalidServer
	}
	srv.ReverseProxy().ErrorHandler = sb.getErrHandler(srv)
	status := srv.PingServer()
	srv.SetAvailable(status)
	sb.lock.Lock()
	sb.servers = append(sb.servers, srv)
	sb.lock.Unlock()
	return nil
}

// Serve - serve incoming request with server's proxy
func (sb *RoundRobinServerBucket) Serve(w http.ResponseWriter, r *http.Request) error {
	srv, err := sb.getNextServer()
	if err != nil {
		return err
	}
	proxy := srv.ReverseProxy()
	log.Println("[proxy] to", srv.Address())
	proxy.ServeHTTP(w, r)
	return nil
}

// getNextServer - round-robin algorithm for chosing next server
// Check if current server is available to serve request,
// in other case we just get next, while not find good one
func (sb *RoundRobinServerBucket) getNextServer() (Server, error) {
	srvAmount := uint64(len(sb.servers))
	if srvAmount == 0 {
		return nil, ErrNoServersAvailable
	}
	next := sb.last % srvAmount
	defer atomic.AddUint64(&sb.last, 1)
	full := next + srvAmount
	sb.lock.Lock()
	for pos := next; pos < full; pos++ {
		index := pos % srvAmount
		if sb.servers[index].IsAvailable() {
			if index != next {
				atomic.StoreUint64(&sb.last, pos)
			}
			srv := sb.servers[index]
			sb.lock.Unlock()
			return srv, nil
		}
	}
	sb.lock.Unlock()
	return nil, ErrAllServersUnreachable
}

// getErrHandler - error handler func for reverse proxy instance
// First, we try MAX_RETRIES time to serve request with current server
// Second, we recurrently call Serve func, to switch server
// Count retries for each server separately
// Count attempts for each request
func (sb *RoundRobinServerBucket) getErrHandler(srv Server) func(w http.ResponseWriter, r *http.Request, e error) {
	return func(w http.ResponseWriter, r *http.Request, e error) {
		attempts := GetAttemptsFromContext(r)
		if attempts > maxAttempts {
			log.Printf("[attempt] %s (%s) Too much attempts, refusing\n", r.RemoteAddr, r.URL.Path)
			http.Error(w, ErrServiceUnavailable.Error(), http.StatusServiceUnavailable)
			return
		}
		log.Printf("[%s] %s\n", srv.Address(), e.Error())
		retries := GetRetriesFromContext(r)
		proxy := srv.ReverseProxy()
		if retries < maxRetries {
			select {
			case <-time.After(10 * time.Millisecond):
				ctx := context.WithValue(r.Context(), RetriesKey, retries+1)

				log.Printf("[retry] %s (%s) Retrying server %d\n", r.RemoteAddr, r.URL.Path, attempts)
				proxy.ServeHTTP(w, r.WithContext(ctx))
			}
			return
		}
		srv.SetAvailable(false)
		log.Printf("[attempt] %s (%s) Attempting server %d\n", r.RemoteAddr, r.URL.Path, attempts)
		ctx := context.WithValue(r.Context(), AttemptsKey, attempts+1)
		sb.Serve(w, r.WithContext(ctx))
	}
}

// Healthcheck - passive server's availability checks
func (sb *RoundRobinServerBucket) Healthcheck() {
	for _, srv := range sb.servers {
		msg := "available"
		status := srv.PingServer()
		srv.SetAvailable(status)
		if !status {
			msg = "unreachable"
		}
		log.Printf("[healthcheck] %s (%s)\n", srv.Address(), msg)
	}
}

// RemoveStale - remove stale servers from storage
func (sb *RoundRobinServerBucket) RemoveStale(timeout time.Duration) {
	sb.lock.Lock()
	newServers := []Server{}
	for _, srv := range sb.servers {
		addr := srv.Address()
		timeDiff := time.Since(time.Unix(srv.LastSeen(), 0))
		if !srv.IsAvailable() && timeDiff > timeout {
			log.Printf("[remove] %s is stale and will be removed\n", addr)
			continue
		}
		newServers = append(newServers, srv)
	}
	if len(newServers) != len(sb.servers) {
		sb.servers = newServers
	}
	sb.lock.Unlock()
}

// RunServices - execute servers pool services
func (sb *RoundRobinServerBucket) RunServices(staleTimeout int) {
	go func() {
		for {
			select {
			case <-time.After(healthCheckPeriod):
				sb.Healthcheck()
			}
		}
	}()
	go func() {
		for {
			select {
			case <-time.After(removeStalePeriod):
				sb.RemoveStale(time.Minute * time.Duration(staleTimeout))
			}
		}
	}()
}
