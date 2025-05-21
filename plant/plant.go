package plant

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

var (
	MinimumWaterLevel = 20
)

type TypeCharacteristics struct {
	Type                  string `json:"type"`
	GrowthRate            int64  `json:"growthRate"`       // between 4-6 weeks at max growth
	WaterConsumptionUnits int64  `json:"waterRequirement"` // 0-1 scale per day
}

var TypeCharacteristicsMap = map[string]TypeCharacteristics{}

type GenerationProperties struct {
	// Generation related properties
	Backdrop string
	Mascot   string
}

type Health struct {
	// Health State (config-map?)
	CurrentGrowth     int64
	CurrentWaterLevel int
	LastWatered       time.Time   // TODO: can we remove this? perhaps compute on a 24 hour basis as a bg task.
	GrowthStage       GrowthStage // TODO: this can be computed to avoid storing in the struct
}

type Plant struct {
	NamespacedName string
	FriendlyName   string
	Type           string
	CreationTime   time.Time
	LastUpdated    time.Time // TODO: I think we can deprecate LastWatered and retain LastUpdated, or remove both...

	Health Health
}

// Update progresses the plant state based on elapsed time
func (p *Plant) Update(currentTime time.Time) {
	p.updateWaterConsumption(currentTime)
	if p.Health.CurrentWaterLevel > MinimumWaterLevel {
		p.updateGrowth(currentTime)
	} else {
		p.LastUpdated = currentTime
	}
	// TODO: Add death logic
}

func (p *Plant) updateWaterConsumption(currentTime time.Time) {
	elapsed := currentTime.Sub(p.Health.LastWatered)
	days := elapsed.Hours() / 24
	waterConsumed := int(float64(TypeCharacteristicsMap[p.Type].WaterConsumptionUnits) * days)
	p.Health.CurrentWaterLevel -= waterConsumed
	if p.Health.CurrentWaterLevel < 0 {
		p.Health.CurrentWaterLevel = 0
	}
	p.Health.LastWatered = currentTime
}

func (p *Plant) updateGrowth(currentTime time.Time) {
	elapsed := currentTime.Sub(p.LastUpdated)
	days := elapsed.Hours() / 24

	growthMultiplier := 1.0
	if p.Health.CurrentWaterLevel < 50 {
		growthMultiplier = float64(p.Health.CurrentWaterLevel) / 50
	}
	if p.Health.CurrentWaterLevel < 20 {
		growthMultiplier = 0.0
	}

	growth := float64(TypeCharacteristicsMap[p.Type].GrowthRate) * growthMultiplier * days
	p.Health.CurrentGrowth += int64(math.Round(growth))

	p.updateGrowthStage()
	p.LastUpdated = currentTime
}

// updateGrowthStage updates the growth stage based on current growth value
func (p *Plant) updateGrowthStage() {
	stages := []GrowthStage{Maturing, Growing, Sprouting, Seeding}
	for _, stage := range stages {
		if p.Health.CurrentGrowth >= growthStageThreshold[stage] {
			p.Health.GrowthStage = stage
			return
		}
	}
}

func (p *Plant) Image() string {
	formattedDate := time.Now().Format("2006-01-02")
	return fmt.Sprintf("%s-%s.png", formattedDate, p.NamespacedName)
}

func (p *Plant) DaysAlive() int {
	currentTime := time.Now()
	elapsed := currentTime.Sub(p.CreationTime)

	// Day 1 is the creation day, no matter what time
	if elapsed < 24*time.Hour {
		return 1
	}

	// After 24 hours, calculate normally but add 1 to start from day 1
	days := (elapsed.Hours() / 24) + 1
	return int(days)
}

// AddWater adds 1-5% water to the plant (up to 100%)
func (p *Plant) AddWater(t time.Time) int {
	// Random water increment between 1% and 5%
	var waterIncrement = rand.Intn(5) + 1

	// Add water (capped at 100%)
	actualToAdd := waterIncrement
	if p.Health.CurrentWaterLevel+waterIncrement > 100 {
		actualToAdd = 100 - p.Health.CurrentWaterLevel
	}
	p.Health.CurrentWaterLevel += actualToAdd

	// We need to know when la
	p.Health.LastWatered = t

	return actualToAdd
}

// RenderGrowthStage maps a plant growth stage to a UIGrowthStage with UI-specific information
func (p *Plant) RenderGrowthStage(stage GrowthStage) UIGrowthStage {
	growthStages := map[GrowthStage]UIGrowthStage{
		Seeding: {
			Stage:       Seeding,
			ColorClass:  "yellow",
			TooltipText: "Your plant is in the seeding stage. Keep soil moist and warm.",
		},
		Sprouting: {
			Stage:       Sprouting,
			ColorClass:  "lime",
			TooltipText: "Your plant is sprouting! First signs of growth are visible.",
		},
		Growing: {
			Stage:       Growing,
			ColorClass:  "green",
			TooltipText: "Your plant is in active growth. Ensure proper watering and light is essential.",
		},
		Maturing: {
			Stage:       Maturing,
			ColorClass:  "emerald",
			TooltipText: "Your plant is maturing. It's reaching its full potential!",
		},
		Dead: {
			Stage:       Dead,
			ColorClass:  "red",
			TooltipText: "Sadly, your plant has died. Consider starting over with a new one.",
		},
	}

	return growthStages[stage]
}

// UIGrowthStage represents a growth stage with UI-specific information
type UIGrowthStage struct {
	Stage       GrowthStage // The growth stage
	ColorClass  string      // CSS class for styling
	TooltipText string      // Text to display in a tooltip
}
