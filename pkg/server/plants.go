package server

import (
	"fmt"
	"github.com/williamnoble/kube-botany/pkg/plant"
)

// plantByID finds a plant by its ID
// It returns the plant if found, or an error if not found
func (s *Server) plantByID(id string) (*plant.Plant, error) {
	for _, p := range s.plants {
		if p.Id == id {
			return p, nil
		}
	}

	return nil, fmt.Errorf("plant not found")
}

// UIGrowthStage represents a growth stage with UI-specific information
type UIGrowthStage struct {
	Stage       plant.GrowthStage // The growth stage
	ColorClass  string            // CSS class for styling
	TooltipText string            // Text to display in a tooltip
}

// RenderGrowthStage maps a plant growth stage to a UIGrowthStage with UI-specific information
func (s *Server) RenderGrowthStage(stage plant.GrowthStage) UIGrowthStage {
	growthStages := map[plant.GrowthStage]UIGrowthStage{
		plant.Seeding: {
			Stage:       plant.Seeding,
			ColorClass:  "yellow",
			TooltipText: "Your plant is in the seeding stage. Keep soil moist and warm.",
		},
		plant.Sprouting: {
			Stage:       plant.Sprouting,
			ColorClass:  "lime",
			TooltipText: "Your plant is sprouting! First signs of growth are visible.",
		},
		plant.Growing: {
			Stage:       plant.Growing,
			ColorClass:  "green",
			TooltipText: "Your plant is in active growth. Ensure proper watering and light is essential.",
		},
		plant.Maturing: {
			Stage:       plant.Maturing,
			ColorClass:  "emerald",
			TooltipText: "Your plant is maturing. It's reaching its full potential!",
		},
		plant.Dead: {
			Stage:       plant.Dead,
			ColorClass:  "red",
			TooltipText: "Sadly, your plant has died. Consider starting over with a new one.",
		},
	}

	return growthStages[stage]
}
