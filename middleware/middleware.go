package middleware

import (
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

// Metrics holds request statistics
type Metrics struct {
	TotalRequests   uint64
	SuccessRequests uint64
	FailedRequests  uint64
}

// GlobalMetrics is the global metrics instance
var GlobalMetrics = &Metrics{}

// LoggingMiddleware logs each request
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Create a response writer wrapper to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(wrapped, r)
		
		duration := time.Since(start)
		log.Printf("%s %s %d %v\n", r.Method, r.URL.Path, wrapped.statusCode, duration)
	})
}

// MetricsMiddleware tracks request metrics
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddUint64(&GlobalMetrics.TotalRequests, 1)
		
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wrapped, r)
		
		if wrapped.statusCode >= 200 && wrapped.statusCode < 400 {
			atomic.AddUint64(&GlobalMetrics.SuccessRequests, 1)
		} else {
			atomic.AddUint64(&GlobalMetrics.FailedRequests, 1)
		}
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// GetMetrics returns current metrics
func GetMetrics() map[string]uint64 {
	return map[string]uint64{
		"total_requests":   atomic.LoadUint64(&GlobalMetrics.TotalRequests),
		"success_requests": atomic.LoadUint64(&GlobalMetrics.SuccessRequests),
		"failed_requests":  atomic.LoadUint64(&GlobalMetrics.FailedRequests),
	}
}
