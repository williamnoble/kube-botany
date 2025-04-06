package server

import (
	"fmt"
	"github.com/williamnoble/kube-botany/pkg/plant"
)

func (s *Server) plantByID(id string) (*plant.Plant, error) {
	for _, p := range s.plants {
		if p.Id == id {
			return p, nil
		}
	}

	return nil, fmt.Errorf("plant not found")
}

type UIGrowthStage struct {
	Stage       plant.GrowthStage
	ColorClass  string
	TooltipText string
}

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
