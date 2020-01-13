package bucket

import "net/http"

const (
	AttemptsKey = "attempts"
	RetriesKey  = "retries"
)

// GetAttemptsFromContext - extract attempts for request
func GetAttemptsFromContext(r *http.Request) int {
	if attempts, ok := r.Context().Value(AttemptsKey).(int); ok {
		return attempts
	}
	return 1
}

// GetAttemptsFromContext - extract the attempts for request
func GetRetriesFromContext(r *http.Request) int {
	if retry, ok := r.Context().Value(RetriesKey).(int); ok {
		return retry
	}
	return 0
}
