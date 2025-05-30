package plant

import (
	"fmt"
	"math"
	"time"
)

type Generator struct {
	Backdrop string
	Mascot   string
}

type Health struct {
	// Health State (config-map?)
	CurrentGrowth     int64
	CurrentWaterLevel int
}

type Plant struct {
	NamespacedName string
	FriendlyName   string
	Variety        *Variety
	CreationTime   time.Time
	LastUpdated    time.Time

	Health Health
}

// Update progresses the plant state based on elapsed time
// Water consumption is calculated, assuming the plant is appropriated watered, it grows.
func (p *Plant) Update(currentTime time.Time) {
	p.updateWaterConsumption(currentTime)
	p.updateGrowth(currentTime)
	p.LastUpdated = currentTime
}

// updateWaterConsumption calculates and applies water consumption since the last update.
func (p *Plant) updateWaterConsumption(currentTime time.Time) {
	// calculate elapsed time (days), since the last update
	elapsedDays := elapsedDays(currentTime, p.LastUpdated)
	//elapsed := currentTime.Sub(p.LastUpdated)
	//days := elapsed.Hours() / 24

	//  determining water consumed based on the consumption rate of a particular variety of plant.
	waterConsumed := int(float64(p.Variety.WaterConsumptionUnitsPerDay) * elapsedDays)

	// reducing the current water level, (bounded at zero).
	p.Health.CurrentWaterLevel -= waterConsumed
	if p.Health.CurrentWaterLevel < 0 {
		p.Health.CurrentWaterLevel = 0
	}
}

func elapsedDays(currentTime time.Time, lastUpdatedTime time.Time) float64 {
	elapsed := currentTime.Sub(lastUpdatedTime)
	days := elapsed.Hours() / 24
	return days
}

// updateGrowth calculates and applies growth progress since the last update.
func (p *Plant) updateGrowth(currentTime time.Time) {
	// calculate elapsed time (days), since the last update
	//elapsed := currentTime.Sub(p.LastUpdated)
	//days := elapsed.Hours() / 24
	elapsedDays := elapsedDays(currentTime, p.LastUpdated)

	// growth is determined solely by the elapsed time and the plant's growth rate.
	// the growth accumulates in CurrentGrowth, which is used to determine the
	// plant's growth stage.
	growth := float64(p.Variety.GrowthRatePerDay) * elapsedDays
	p.Health.CurrentGrowth += int64(math.Round(growth))
}

// GrowthStage returns the growth stage based on the current growth value.
func (p *Plant) GrowthStage() string {
	// maps.Keys() is non-deterministic so we'll hardcode
	stages := []GrowthStage{Maturing, Growing, Sprouting, Seeding}
	for _, stage := range stages {
		if p.Health.CurrentGrowth >= growthStageThreshold[stage] {
			return stage.String()
		}
	}
	return Seeding.String()
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
func (p *Plant) DaysToMaturity() int {
	if p.Health.CurrentGrowth >= growthStageThreshold[Maturing] {
		return 0
	}

	remainingGrowth := growthStageThreshold[Maturing] - p.Health.CurrentGrowth
	daysRemaining := float64(remainingGrowth) / float64(p.Variety.GrowthRatePerDay)
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

// AddWater fully waters the plant
func (p *Plant) AddWater() int {
	amountAdded := 100 - p.Health.CurrentWaterLevel
	p.Health.CurrentWaterLevel = 100
	return amountAdded
}

func (p *Plant) CurrentWaterLevel() int {
	return p.Health.CurrentWaterLevel
}

func (p *Plant) CurrentGrowth() int64 { return p.Health.CurrentGrowth }

func (p *Plant) Healthy() bool {
	return p.CurrentWaterLevel() >= p.Variety.MinimumWaterLevel
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
			TooltipText: "Your plant is in active growth. Ensuring proper watering is essential.",
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
