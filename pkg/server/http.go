package server

import (
	"encoding/json"
	"fmt"
	"github.com/williamnoble/kube-botany/pkg/types"
	"net/http"
)

// Encode serializes a value to JSON and writes it to the HTTP response.
// It sets the Content-Type header to "application/json" and the HTTP status code.
// If wrap is provided, it wraps the value in a JSON object with the wrap string as the key
func (s *Server) encodeJsonResponse(w http.ResponseWriter, r *http.Request, status int, data interface{}, wrap ...string) error {
	if len(wrap) > 0 {
		data = map[string]interface{}{
			wrap[0]: data,
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return fmt.Errorf("encodeJsonResponse json: %w", err)
	}
	return nil
}

// Decode deserializes a JSON request body into a value
// It logs an error if decoding fails
func (s *Server) decodeJsonRequest(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		s.Logger.Error("failed to decode request body", "error", err)
		return fmt.Errorf("decode json request: %w", err)
	}
	return nil
}

// InternalServerErrorResponse is a convenience function for returning an error response with status code 500
func (s *Server) InternalServerErrorResponse(w http.ResponseWriter, err error) {
	s.Logger.Error("internal server error", "error", err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	errorResponse := map[string]string{
		"error": http.StatusText(http.StatusInternalServerError),
	}

	if encodeErr := json.NewEncoder(w).Encode(errorResponse); encodeErr != nil {
		s.Logger.Error("failed to encode error response", "encode_error", encodeErr)
	}
}

// WaterResponse is the response returned by the water endpoint
type WaterResponse struct {
	Message string         `json:"message"` // Message about the watering result
	Plant   types.PlantDTO `json:"plant"`   // Updated plant information
}

// WaterRequest contains the Id identifier of the plant being watered
type WaterRequest struct {
	Id string `json:"id"` // ID of the plant to water
}
