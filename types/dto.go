package types

import (
	"fmt"
	"github.com/williamnoble/kube-botany/plant"
)

// PlantDTO represents a plant in API responses and UI rendering
type PlantDTO struct {
	NamespacedName    string `json:"namespaced_name"`      // Unique identifier for the plant
	FriendlyName      string `json:"friendly_name"`        // Display name for the plant
	Variety           string `json:"variety"`              // Variety of plant (e.g., bonsai, sunflower)
	DaysAlive         int    `json:"days_alive,omitempty"` // Number of days the plant has been alive
	DaysToMaturity    int    `json:"days_to_maturity,omitempty"`
	CurrentWaterLevel int    `json:"current_water_level"` // The current water level
	GrowthStage       string `json:"growth_stage"`        // Derives growth stage from current growth

	Image string `json:"image,omitempty"` // Path to the plant's image
}

// IntoPlantDTO converts a plant.Plant to a PlantDTO for API responses and UI rendering
func IntoPlantDTO(p *plant.Plant) PlantDTO {

	r := PlantDTO{
		NamespacedName:    p.NamespacedName, // Unique ID
		FriendlyName:      p.FriendlyName,
		Variety:           p.Variety.Type,
		DaysAlive:         p.DaysAlive(),
		DaysToMaturity:    p.DaysToMaturity(),
		CurrentWaterLevel: p.WaterLevel(),
		GrowthStage:       p.GrowthStage(),
		Image:             fmt.Sprintf("/static/images/%s", p.Image()),
	}

	return r
}

// FromPlantDTO converts from PlantDTO to *plant.Plant for API responses and UI rendering
func FromPlantDTO(p *plant.Plant) PlantDTO {
	r := PlantDTO{
		NamespacedName:    p.NamespacedName,
		FriendlyName:      p.FriendlyName,
		Variety:           p.Variety.Type,
		DaysAlive:         p.DaysAlive(),
		CurrentWaterLevel: p.WaterLevel(),
		GrowthStage:       p.GrowthStage(),
	}
	return r
}
