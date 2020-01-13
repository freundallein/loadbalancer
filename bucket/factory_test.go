package bucket

import (
	"reflect"
	"testing"
)

func TestNewServer(t *testing.T) {
	observed, err := NewServer("http://testhost:8000")
	if err != nil {
		t.Error(err.Error())
	}
	observedType := reflect.TypeOf(observed)
	expectedType := reflect.TypeOf((*Server)(nil)).Elem()
	if reflect.PtrTo(observedType).Implements(expectedType) {
		t.Error("Expected", expectedType, "got", observedType)
	}
}
func TestNewServerInvalidUrl(t *testing.T) {
	_, err := NewServer("\x80testhost:8000")
	if err == nil {
		t.Error("Url validation is broken")
		return
	}
}

func TestNew(t *testing.T) {
	observed, _ := New(RoundRobin)
	observedType := reflect.TypeOf(observed)
	expectedType := reflect.TypeOf((*ServerBucket)(nil)).Elem()
	if reflect.PtrTo(observedType).Implements(expectedType) {
		t.Error("Expected", expectedType, "got", observedType)
	}
}

func TestNewInvalidAlgorithm(t *testing.T) {
	observed, err := New("invalid")
	if err == nil {
		t.Error("Algorithm check is broken")
	}
	if observed != nil {
		t.Error("Expected nil")
	}
}
