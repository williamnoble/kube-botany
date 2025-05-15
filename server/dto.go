package server

import (
	"github.com/williamnoble/kube-botany/plant"
	"time"
)

// PlantDTO represents a plant in API responses and UI rendering
type PlantDTO struct {
	Id           string    `json:"id"`            // Unique identifier for the plant
	Name         string    `json:"name"`          // Display name for the plant
	Type         string    `json:"type"`          // Type of plant (e.g., bonsai, sunflower)
	Age          string    `json:"age"`           // Age of the plant as a formatted string
	Growth       int64     `json:"growth"`        // Growth value of the plant
	GrowthStage  string    `json:"stage"`         // Current growth stage (e.g., seeding, sprouting)
	WaterLevel   int       `json:"water"`         // Current water level (0-100%)
	WateredLast  time.Time `json:"watered_last"`  // Time when the plant was last watered
	CreationTime time.Time `json:"creation_time"` // Time when the plant was created (matches CR)

	// Fields used for UI rendering
	Image     string `json:"image,omitempty"`      // Path to the plant's image
	DaysAlive int    `json:"days_alive,omitempty"` // Number of days the plant has been alive
}

// plantDTO converts a plant.Plant to a PlantDTO for API responses and UI rendering
// It extracts relevant fields from the plant and formats them for display
func (s *Server) plantDTO(p *plant.Plant) PlantDTO {
	r := PlantDTO{
		Id:          p.Id,
		Name:        p.FriendlyName,
		Type:        string(p.Type),
		GrowthStage: p.GrowthStage.String(),
		Age:         time.Since(p.CreationTime).Round(time.Second).String(),
		WaterLevel:  p.CurrentWaterLevel,
		WateredLast: p.LastWatered,
		Image:       "",
		DaysAlive:   p.DaysAlive(),
	}

	return r
}

// WaterResponse is the response returned by the water endpoint
type WaterResponse struct {
	Message string   `json:"message"` // Message about the watering result
	Plant   PlantDTO // Updated plant information
}

// WaterRequest is the request body for the water endpoint
type WaterRequest struct {
	Id string `json:"id"` // ID of the plant to water
}
