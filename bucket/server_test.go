package bucket

import (
	"net/http/httputil"
	"reflect"
	"testing"
	"time"
)

func TestIsAvailable(t *testing.T) {
	srv, _ := NewServer("http://testhost:8000")
	if !srv.IsAvailable() {
		t.Error("Expected true, got false")
	}
}

func TestSetAvailable(t *testing.T) {
	srv, _ := NewServer("http://testhost:8000")
	srv.SetAvailable(false)
	if srv.IsAvailable() {
		t.Error("Expected false, got true")
	}
	srv.SetAvailable(true)
	if !srv.IsAvailable() {
		t.Error("Expected true, got false")
	}
}

func TestAddress(t *testing.T) {
	srv, _ := NewServer("http://testhost:8000")
	url := srv.Address()
	if url.Host != "testhost:8000" {
		t.Error("Expected", "testhost:8000", "got", url.Host)
	}
}

func TestReverseProxy(t *testing.T) {
	srv, _ := NewServer("http://testhost:8000")
	observedType := reflect.TypeOf(srv.ReverseProxy())
	expectedType := reflect.TypeOf(&httputil.ReverseProxy{})
	if observedType != expectedType {
		t.Error("Expected", expectedType, "got", observedType)
	}
}

func TestLastSeen(t *testing.T) {
	srv, _ := NewServer("http://testhost:8000")
	if srv.LastSeen() != time.Now().Unix() {
		t.Error("Expected", time.Now().Unix(), "got", srv.LastSeen())
	}
}
