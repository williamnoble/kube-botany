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

type Generator struct {
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
	Variety        string
	CreationTime   time.Time
	LastUpdated    time.Time // TODO: I think we can deprecate LastWatered and retain LastUpdated, or remove both...

	Health Health
}

// Update progresses the plant state based on elapsed time
// Water consumption is calculated, assuming the plant is appropriated watered, it grows.
func (p *Plant) Update(currentTime time.Time, t Variety) {
	p.updateWaterConsumption(currentTime, t)
	if p.Health.CurrentWaterLevel > MinimumWaterLevel {
		p.updateGrowth(currentTime, t)
	} else {
		p.LastUpdated = currentTime
	}
}

// updateWaterConsumption calculates and applies water consumption since the last update.
func (p *Plant) updateWaterConsumption(currentTime time.Time, t Variety) {
	// calculate elapsed time in days, since the last update
	elapsed := currentTime.Sub(p.LastUpdated)
	days := elapsed.Hours() / 24

	//  determining water consumed based on the consumption rate of a particular variety of plant.
	waterConsumed := int(float64(t.WaterConsumptionUnits) * days)

	// reducing the current water level, (bounded at zero).
	p.Health.CurrentWaterLevel -= waterConsumed
	if p.Health.CurrentWaterLevel < 0 {
		p.Health.CurrentWaterLevel = 0
	}

}

// updateGrowth calculates and applies growth progress since the last update.
func (p *Plant) updateGrowth(currentTime time.Time, t Variety) {
	elapsed := currentTime.Sub(p.LastUpdated)
	days := elapsed.Hours() / 24

	// Growth is determined solely by the elapsed time and the plant's growth rate.
	// The growth accumulates in CurrentGrowth, which is used to determine the
	// plant's growth stage.
	growth := float64(t.GrowthRate) * days
	p.Health.CurrentGrowth += int64(math.Round(growth))
	p.LastUpdated = currentTime
}

// GrowthStage returns the growth stage based on the current growth value.
func (p *Plant) GrowthStage() string {
	stages := []GrowthStage{Maturing, Growing, Sprouting, Seeding}
	for _, stage := range stages {
		if p.Health.CurrentGrowth >= growthStageThreshold[stage] {
			return stage.String()
		}
	}
	return "dead"
}

// GrowthPercentage returns the plant's current growth as a percentage of full maturity (capped at 100%).
func (p *Plant) GrowthPercentage() int {
	maturingThreshold := growthStageThreshold[Maturing]
	if maturingThreshold <= 0 {
		return 0
	}

	percentage := (float64(p.Health.CurrentGrowth) / float64(maturingThreshold)) * 100
	if percentage > 100 {
		return 100
	}
	return int(percentage)
}

// DaysToMaturity estimates the number of days until the plant reaches maturity
// based on its current growth and growth rate (0 if already mature)
func (p *Plant) DaysToMaturity(t Variety) int {
	if p.Health.CurrentGrowth >= growthStageThreshold[Maturing] {
		return 0
	}

	remainingGrowth := growthStageThreshold[Maturing] - p.Health.CurrentGrowth
	daysRemaining := float64(remainingGrowth) / float64(t.GrowthRate)
	return int(math.Ceil(daysRemaining))
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

func (p *Plant) WaterLevel() int {
	return p.Health.CurrentWaterLevel
}

func (p *Plant) Validate() error {
	// TODO, implement Zod or Zog? I forgot
	return nil
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
