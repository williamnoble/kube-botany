package server

import (
	"net/http"
	"time"
)

// statusRecorder is a custom ResponseWriter that keeps track of the response status code
type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code before calling the underlying ResponseWriter's WriteHeader
func (sr *statusRecorder) WriteHeader(code int) {
	sr.statusCode = code
	sr.ResponseWriter.WriteHeader(code)
}

// requestLogger is a middleware that logs information about HTTP requests
// It logs when a request starts with method, path, remote address, and user agent
// It also logs when a request completes with method, path, duration, and status code
func (s *Server) requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		s.Logger.Info("request started",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)

		recorder := &statusRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(recorder, r)

		duration := time.Since(start)

		s.Logger.Info("request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", recorder.statusCode,
			"duration_ms", duration.Milliseconds(),
		)
	})
}
