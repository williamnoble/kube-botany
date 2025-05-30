package render

import (
	"fmt"
	"github.com/williamnoble/kube-botany/plant"
)

// ASCIIRenderer renders the plant as ASCII art
type ASCIIRenderer struct{}

func NewASCIIRenderer() *ASCIIRenderer {
	return &ASCIIRenderer{}
}

func (r *ASCIIRenderer) RenderText(p *plant.Plant) string {
	fmt.Println("rendering")

	asciiArt := map[string]string{
		"seeding":   seeding,
		"sprouting": sprouting,
		"growing":   growing,
		"maturing":  maturing,
		"dead":      dead,
	}

	var text string
	text += fmt.Sprintf("\nName: %s\n", p.NamespacedName)
	text += fmt.Sprintf("FriendlyName: %s\n", p.FriendlyName)
	waterBar := renderBar("Water Level:", p.CurrentWaterLevel(), 100)
	text += waterBar + "\n"
	text += fmt.Sprintf("Growth Stage: %s\n", p.GrowthStage())
	text += fmt.Sprintf("Created: %s\n", p.HumanCreationTime())
	text += fmt.Sprintf("Day: %d (%d days to maturity)\n", p.DaysAlive(), p.DaysToMaturity())
	text += "Image" + "\n"
	text += asciiArt[p.GrowthStage()]
	return text
}

func renderBar(label string, value, max int) string {
	const barWidth = 20
	filled := int(float64(value) / float64(max) * float64(barWidth))
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
	// Fix format specifier to match the int type
	bar += "] " + fmt.Sprintf("%d%%", value)
	return bar
}
