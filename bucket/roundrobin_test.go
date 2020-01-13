package bucket

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
)

// type RoundRobinServerBucket struct {
// 	servers []Server     // servers storage
// 	last    uint64       // last used server index
// 	lock    sync.RWMutex // lock for servers slice
// }
type MockServer struct {
	address     *url.URL
	isAvailable bool
	ping        bool
}

func (ms *MockServer) IsAvailable() bool {
	return ms.isAvailable
}
func (ms *MockServer) SetAvailable(status bool) {
	ms.isAvailable = status
}

func (ms *MockServer) Address() *url.URL {
	return ms.address
}
func (ms *MockServer) ReverseProxy() *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(ms.address)
	return proxy
}
func (ms *MockServer) LastSeen() int64 {
	return 1
}

func (ms *MockServer) PingServer() bool {
	return ms.ping
}

func TestAddServer(t *testing.T) {
	bckt := &RoundRobinServerBucket{
		servers: []Server{},
	}

	srv, _ := NewServer("http://testhost:8000")
	bckt.AddServer(srv)
	if len(bckt.servers) != 1 {
		t.Error("Expected", 1, "got", len(bckt.servers))
	}
	if srv.IsAvailable() {
		t.Error("Expected", false, "got", srv.IsAvailable())
	}
	// Check if proxy ErrorHandler was installed
	errHandler := &srv.ReverseProxy().ErrorHandler
	if errHandler == nil {
		t.Error("Expected", "error handler", "got", errHandler)
	}
}

func TestServe(t *testing.T) {
	// TODO:
}
func TestGetNextServer(t *testing.T) {
	bckt := &RoundRobinServerBucket{
		servers: []Server{},
	}
	addrs := []string{"http://testhost1:8000", "http://testhost2:8000", "http://testhost3:8000"}
	for i := 0; i < 3; i++ {
		addr, _ := url.Parse(addrs[i])
		srv := &MockServer{
			address:     addr,
			isAvailable: true,
			ping:        true,
		}
		bckt.AddServer(srv)
	}
	for i := 0; i < 6; i++ {
		srv, _ := bckt.getNextServer()
		host := strings.Split(addrs[i%3], "/")[2]
		if srv.Address().Host != host {
			t.Error("Expected", host, "got", srv.Address().Host)
		}
	}

}
func TestGetNextServerEmpty(t *testing.T) {
	bckt := &RoundRobinServerBucket{
		servers: []Server{},
	}
	srv, err := bckt.getNextServer()
	if err == nil {
		t.Error("Expected", ErrNoServersAvailable, "got", nil)
	}
	if srv != nil {
		t.Error("Expected", nil, "got", srv)
	}
}
func TestGetNextServerUnreachable(t *testing.T) {
	bckt := &RoundRobinServerBucket{
		servers: []Server{},
	}
	addr, _ := url.Parse("http://testhost1:8000")

	bckt.AddServer(&MockServer{address: addr, isAvailable: false})
	srv, err := bckt.getNextServer()
	if err == nil {
		t.Error("Expected", ErrAllServersUnreachable, "got", nil)
	}
	if srv != nil {
		t.Error("Expected", nil, "got", srv)
	}
}
func TestGetErrHandler(t *testing.T) {
	bckt := &RoundRobinServerBucket{
		servers: []Server{},
	}
	addr, _ := url.Parse("http://testhost1:8000")
	srv := &MockServer{address: addr, isAvailable: true}
	bckt.AddServer(srv)
	errHandler := bckt.getErrHandler(srv)
	observedType := reflect.TypeOf(errHandler)
	expectedType := reflect.TypeOf(func(http.ResponseWriter, *http.Request, error) {})
	if observedType != expectedType {
		t.Error("Expected", expectedType, "got", observedType)
	}

}

func TestHealthcheck(t *testing.T) {
	bckt := &RoundRobinServerBucket{
		servers: []Server{},
	}
	addrs := []string{"http://testhost7:8000", "http://testhost8:8000", "http://testhost9:8000"}
	flag := true
	for i := 0; i < 3; i++ {
		addr, _ := url.Parse(addrs[i])
		srv := &MockServer{
			address:     addr,
			isAvailable: true,
			ping:        flag,
		}
		bckt.AddServer(srv)
		flag = !flag
	}
	bckt.Healthcheck()
	flag = true
	for _, srv := range bckt.servers {
		if srv.IsAvailable() != flag {
			t.Error("Expected", srv, "got", srv.IsAvailable())
		}
		flag = !flag
	}

}

func TestRemoveStale(t *testing.T) {
	bckt := &RoundRobinServerBucket{
		servers: []Server{},
	}
	addrs := []string{"http://testhost1:8000", "http://testhost2:8000", "http://testhost3:8000"}
	for i := 0; i < 3; i++ {
		addr, _ := url.Parse(addrs[i])
		srv := &MockServer{
			address:     addr,
			isAvailable: false,
		}
		bckt.AddServer(srv)
	}
	bckt.RemoveStale(time.Second * 0)
	if len(bckt.servers) != 0 {
		t.Error("Expected", 0, "got", len(bckt.servers))
	}
}
func TestRemoveStaleDifferent(t *testing.T) {
	bckt := &RoundRobinServerBucket{
		servers: []Server{},
	}
	addrs := []string{"http://testhost3:8000", "http://testhost4:8000", "http://testhost5:8000"}
	flag := true
	for i := 0; i < 3; i++ {
		addr, _ := url.Parse(addrs[i])
		srv := &MockServer{
			address:     addr,
			isAvailable: flag,
			ping:        flag,
		}
		bckt.AddServer(srv)
		flag = !flag
	}
	bckt.RemoveStale(time.Second * 0)
	if len(bckt.servers) != 2 {
		t.Error("Expected", 2, "got", len(bckt.servers))
	}
}
