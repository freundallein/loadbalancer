package bucket

import (
	"net"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

// DefaultServer - default backend server implementation
type DefaultServer struct {
	address      *url.URL               // server address
	isAvailable  bool                   // current status
	lock         sync.RWMutex           // lock for isAvailable attribute
	reverseProxy *httputil.ReverseProxy // reverse proxy for request forwarding
	lastSeen     int64                  // unixtime for last time, when server was available
}

// IsAvailable - getter for server's availability
func (ds *DefaultServer) IsAvailable() bool {
	ds.lock.RLock()
	status := ds.isAvailable
	ds.lock.RUnlock()
	return status
}

// SetAvailable - setter for server's availability
func (ds *DefaultServer) SetAvailable(status bool) {
	ds.lock.Lock()
	ds.isAvailable = status
	if status {
		ds.lastSeen = time.Now().Unix()
	}
	ds.lock.Unlock()
}

// Address - getter for server address
func (ds *DefaultServer) Address() *url.URL {
	return ds.address
}

// ReverseProxy - getter for reversing proxy
func (ds *DefaultServer) ReverseProxy() *httputil.ReverseProxy {
	return ds.reverseProxy
}

// LastSeen - getter for lastSeen time field
func (ds *DefaultServer) LastSeen() int64 {
	return ds.lastSeen
}

func (ds *DefaultServer) PingServer() bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", ds.address.Host, timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
