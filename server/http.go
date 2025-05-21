package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Encode serializes a value to JSON and writes it to the HTTP response.
// It sets the Content-Variety header to "application/json" and the HTTP status code.
// If wrap is provided, it wraps the value in a JSON object with the wrap string as the key
func (s *Server) encodeJsonResponse(w http.ResponseWriter, r *http.Request, status int, v interface{}, wrap ...string) error {
	if len(wrap) > 0 {
		v = map[string]interface{}{
			wrap[0]: v,
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encodeJsonResponse json: %w", err)
	}
	return nil
}

// Decode deserializes a JSON request body into a value
// It logs an error if decoding fails
func (s *Server) decodeJsonResponse(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		s.Logger.Error("failed to decodeJsonResponse request body", "error", err)
		return fmt.Errorf("decodeJsonResponse json: %w", err)
	}
	return nil
}
