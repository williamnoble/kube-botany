package server

import (
	"net/http"
	"time"
)

// requestLogger is a middleware that logs information about HTTP requests
// It logs when a request starts with method, path, remote address, and user agent
// It also logs when a request completes with method, path, and duration
func (s *Server) requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		s.logger.Info("request started",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)

		next.ServeHTTP(w, r)

		duration := time.Since(start)

		s.logger.Info("request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"duration_ms", duration.Milliseconds(),
		)
	})
}
