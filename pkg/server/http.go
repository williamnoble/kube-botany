package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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

func (s *Server) decode(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		s.logger.Error("failed to decode water request", "error", err)
		return fmt.Errorf("decode json: %w", err)
	}
	return nil
}
