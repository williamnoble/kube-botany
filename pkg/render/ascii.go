package render

import (
	"fmt"
	"github.com/williamnoble/kube-botany/pkg/plant"
)

// ASCIIRenderer render plants as ASCII art
type ASCIIRenderer struct{}

func NewASCIIRenderer() *ASCIIRenderer {
	return &ASCIIRenderer{}
}

func (r *ASCIIRenderer) RenderFern(p *plant.Plant) string {
	fmt.Println("render")
	var art string

	// Render fern based on growth stage
	switch p.GrowthStage {
	case plant.Dead:
		art = `
     x    x
    x \  / x
  x   \|/   x
       |
       |
  _____|_____
 /           \
`
	case plant.Seeding:
		art = `
     .
     |
  ___|___
 /       \
`
	case plant.Sprouting:
		art = `
     ^
     |
  ___|___
 /       \
`
	case plant.Growing:
		art = `
    \|/
     |
     |
  ___|___
 /       \
`
	case plant.Maturing:
		art = `
     /\
  //|  |\\
 // |  | \\
 /  |  |  \
    |  |
    |  |
  __|__|__
 /        \
`
	default:
		art = `
     ?
     |
  ___|___
 /       \
`
	}

	// Add status indicators
	art = art + "\n"

	// Add water level bar
	waterBar := renderBar("AddWater", float64(p.WaterLevel), 100)
	art = art + waterBar + "\n"

	// Add growth information
	growthInfo := fmt.Sprintf("Growth: %d (Stage: %s)",
		p.Growth, p.GrowthStage)
	art = art + growthInfo + "\n"

	// Add water consumption information
	waterInfo := fmt.Sprintf("AddWater consumption: %d%% per day",
		p.WaterConsumptionRatePerDay)
	art = art + waterInfo + "\n"

	return art
}

func renderBar(label string, value, max float64) string {
	const barWidth = 20
	filled := int((value / max) * float64(barWidth))
	if filled < 0 {
		filled = 0
	} else if filled > barWidth {
		filled = barWidth
	}

	bar := label + ": ["
	for i := 0; i < barWidth; i++ {
		if i < filled {
			bar += "#"
		} else {
			bar += " "
		}
	}
	bar += "] " + fmt.Sprintf("%.1f%%", value)
	return bar
}
