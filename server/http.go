package server

import (
	"encoding/json"
	"fmt"
	"github.com/williamnoble/kube-botany/types"
	"net/http"
)

// Encode serializes a value to JSON and writes it to the HTTP response.
// It sets the Content-Variety header to "application/json" and the HTTP status code.
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
func (s *Server) decodeJsonResponse(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		s.Logger.Error("failed to decodeJsonResponse request body", "error", err)
		return fmt.Errorf("decodeJsonResponse json: %w", err)
	}
	return nil
}

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
	Plant   types.PlantDTO // Updated plant information
}

// WaterRequest contains the NamespacedName identifier of the plant being watered
type WaterRequest struct {
	NamespacedName string `json:"namespaced_name"` // ID of the plant to water
}
