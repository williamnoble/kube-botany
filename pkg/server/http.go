package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// encode serializes a value to JSON and writes it to the HTTP response
// It sets the Content-Type header to "application/json" and the HTTP status code
// If wrap is provided, it wraps the value in a JSON object with the wrap string as the key
func (s *Server) encode(w http.ResponseWriter, r *http.Request, status int, v interface{}, wrap ...string) error {
	if len(wrap) > 0 {
		v = map[string]interface{}{
			wrap[0]: v,
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

// decode deserializes a JSON request body into a value
// It logs an error if decoding fails
func (s *Server) decode(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		s.logger.Error("failed to decode request body", "error", err)
		return fmt.Errorf("decode json: %w", err)
	}
	return nil
}
