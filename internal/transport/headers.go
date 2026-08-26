package transport

import (
	"fmt"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

var requestSequence atomic.Uint64

func RequestID(r *http.Request) string {
	if value := r.Header.Get("X-Request-ID"); value != "" {
		return value
	}
	sequence := requestSequence.Add(1)
	return fmt.Sprintf("kp-%d-%d", time.Now().UnixNano(), sequence)
}

func WithRequestHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := RequestID(r)
		w.Header().Set("X-Request-ID", requestID)
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		if r.Method == http.MethodOptions {
			w.Header().Set("Allow", "GET, POST, OPTIONS")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func ParsePositive(value string) (int64, error) {
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed <= 0 {
		return 0, fmt.Errorf("value must be a positive integer")
	}
	return parsed, nil
}

func ParseRange(value string, minimum, maximum int) (int, error) {
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("value must be an integer")
	}
	if parsed < minimum || parsed > maximum {
		return 0, fmt.Errorf("value must be between %d and %d", minimum, maximum)
	}
	return parsed, nil
}

func ContentTypeIsJSON(r *http.Request) bool {
	contentType := r.Header.Get("Content-Type")
	return contentType == "" || contentType == "application/json" || len(contentType) >= 16 && contentType[:16] == "application/json"
}
