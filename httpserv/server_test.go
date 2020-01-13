package httpserv

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	srv := New(9000, nil)
	observedType := reflect.TypeOf(srv)
	expectedType := reflect.TypeOf(&http.Server{})
	if observedType != expectedType {
		t.Error("Expected", expectedType, "got", observedType)
	}
}
