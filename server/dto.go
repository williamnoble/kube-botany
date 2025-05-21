package server

import (
	"github.com/williamnoble/kube-botany/plant"
	"time"
)

// PlantDTO represents a plant in API responses and UI rendering
type PlantDTO struct {
	NamespacedName string `json:"namespaced_name"` // Unique identifier for the plant
	FriendlyName   string `json:"friendly_name"`   // Display name for the plant
	Variety        string `json:"variety"`         // Variety of plant (e.g., bonsai, sunflower)

	Age               string `json:"age"`                  // Age of the plant as a formatted string
	DaysAlive         int    `json:"days_alive,omitempty"` // Number of days the plant has been alive
	CurrentWaterLevel int    `json:"current_water_level"`  // The current water level
	GrowthStage       string `json:"growth_stage"`         // Derives growth stage from current growth

	Image string `json:"image,omitempty"` // Path to the plant's image
}

// plantDTO converts a plant.Plant to a PlantDTO for API responses and UI rendering
func (s *Server) plantDTO(p *plant.Plant) PlantDTO {
	r := PlantDTO{
		NamespacedName:    p.NamespacedName,
		FriendlyName:      p.FriendlyName,
		Variety:           p.Variety,
		Age:               time.Since(p.CreationTime).Round(time.Second).String(),
		DaysAlive:         p.DaysAlive(),
		CurrentWaterLevel: p.WaterLevel(),
		GrowthStage:       p.GrowthStage(),
		Image:             "",
	}

	return r
}

// WaterResponse is the response returned by the water endpoint
type WaterResponse struct {
	Message string   `json:"message"` // Message about the watering result
	Plant   PlantDTO // Updated plant information
}

// WaterRequest contains the NamespacedName identifier of the plant being watered
type WaterRequest struct {
	NamespacedName string `json:"namespaced_name"` // ID of the plant to water
}
