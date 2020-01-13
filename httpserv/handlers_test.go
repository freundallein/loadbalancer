package httpserv

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"../bucket"
)

const (
	GoodResponse = "OK"
	BadResponse  = "NOT OK"
	ErrReponse   = "Error"
)

type MockBucket struct {
	response string
}

func (mb *MockBucket) AddServer(bucket.Server) error { return nil }

func (mb *MockBucket) Serve(w http.ResponseWriter, r *http.Request) error {
	switch mb.response {
	case GoodResponse:
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, GoodResponse)
	case BadResponse:
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, BadResponse)
	case ErrReponse:
		return errors.New(ErrReponse)
	}
	return nil

}

func (mb *MockBucket) Healthcheck() {}

func (mb *MockBucket) RemoveStale(time.Duration) {}

func (mb *MockBucket) RunServices(int) {}

func TestBalanceGoodResponse(t *testing.T) {
	handlerFunc := LoadBalance(&MockBucket{GoodResponse})
	observedType := reflect.TypeOf(handlerFunc)
	expectedType := reflect.TypeOf(func(w http.ResponseWriter, r *http.Request) {})
	if observedType != expectedType {
		t.Error("Expected", expectedType, "got", observedType)
	}
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(handlerFunc)
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
func TestBalanceBadResponse(t *testing.T) {
	handlerFunc := LoadBalance(&MockBucket{BadResponse})
	observedType := reflect.TypeOf(handlerFunc)
	expectedType := reflect.TypeOf(func(w http.ResponseWriter, r *http.Request) {})
	if observedType != expectedType {
		t.Error("Expected", expectedType, "got", observedType)
	}
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(handlerFunc)
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
func TestBalanceErrResponse(t *testing.T) {
	handlerFunc := LoadBalance(&MockBucket{ErrReponse})
	observedType := reflect.TypeOf(handlerFunc)
	expectedType := reflect.TypeOf(func(w http.ResponseWriter, r *http.Request) {})
	if observedType != expectedType {
		t.Error("Expected", expectedType, "got", observedType)
	}
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(handlerFunc)
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusServiceUnavailable {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
