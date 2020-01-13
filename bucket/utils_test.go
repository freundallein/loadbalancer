package bucket

import (
	"context"
	"net/http"
	"testing"
)

func TestGetAttemptsFromContextEmpty(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/test", nil)
	expected := 1
	observed := GetAttemptsFromContext(request)
	if observed != expected {
		t.Error("Expected", expected, "got ", observed)
	}
}

func TestGetAttemptsFromContextFilled(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/test", nil)
	expected := 4
	ctx := context.WithValue(request.Context(), AttemptsKey, 4)
	observed := GetAttemptsFromContext(request.WithContext(ctx))
	if observed != expected {
		t.Error("Expected", expected, "got ", observed)
	}
}

func TestGetRetriesFromContextEmpty(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/test", nil)
	expected := 0
	observed := GetRetriesFromContext(request)
	if observed != expected {
		t.Error("Expected", expected, "got ", observed)
	}
}

func TestGetRetriesFromContextFilled(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/test", nil)
	expected := 4
	ctx := context.WithValue(request.Context(), RetriesKey, 4)
	observed := GetRetriesFromContext(request.WithContext(ctx))
	if observed != expected {
		t.Error("Expected", expected, "got ", observed)
	}
}
