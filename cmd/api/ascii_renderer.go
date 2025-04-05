package main

import (
	"fmt"
	"github.com/williamnoble/kube-botany/pkg/plant"
)

// ASCIIRenderer renders plants as ASCII art
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

	// Add health bar
	healthBar := renderBar("Health", p.Health, 100)
	art = art + healthBar + "\n"

	// Add water level bar
	waterBar := renderBar("Water", p.WaterLevel, 100)
	art = art + waterBar + "\n"

	// Add optimal water range
	optimalRange := fmt.Sprintf("Optimal water range: %.1f%% - %.1f%%",
		p.Characteristics.OptimalWaterMin,
		p.Characteristics.OptimalWaterMax)
	art = art + optimalRange + "\n"

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
