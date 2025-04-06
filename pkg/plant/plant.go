package plant

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Plant struct {
	Id           string //NamespacedName
	Name         string // Fallback to ID when not provided
	Type         Type
	CreationTime time.Time

	// Whether the plant can die
	CanDie bool

	// AddWater related properties
	WaterConsumptionRatePerDay int // Random, between 2-6% per day
	MinimumWaterLevel          int
	WaterLevel                 int
	LastWatered                time.Time

	// Growth related properties
	GrowthRate  int64 // Random, between 4-6 weeks at max growth
	Growth      int64
	GrowthStage GrowthStage
	LastUpdated time.Time
}

func NewPlant(
	namespacedName string,
	friendlyName string,
	plantType Type,
	creationTime time.Time,
	canDie bool) *Plant {

	// Provide a Random Growth Rate, our plant requires 1000 Growth to reach maturity. Hence,
	// We divide 1000 by the number of days required to reach maturity.
	minDaysToMaturity := 28 // 4 weeks
	maxDaysToMaturity := 42 // 6 weeks
	// Random days to maturity between min and max
	daysToMaturity := rand.Intn(maxDaysToMaturity-minDaysToMaturity+1) + minDaysToMaturity
	// Calculate daily growth rate
	growthRate := int64(1000 / daysToMaturity)

	return &Plant{
		Id:           namespacedName,
		Name:         friendlyName,
		Type:         plantType,
		CreationTime: creationTime,
		CanDie:       canDie,

		GrowthStage: Seeding,
		Growth:      0,
		LastUpdated: creationTime,

		WaterLevel:        50,
		MinimumWaterLevel: 20,
		LastWatered:       creationTime,

		GrowthRate:                 growthRate,
		WaterConsumptionRatePerDay: rand.Intn(6) + 2,
	}
}

// Update progresses the plant state based on elapsed time
func (p *Plant) Update(currentTime time.Time) {
	p.updateWaterConsumption(currentTime)
	if p.WaterLevel > p.MinimumWaterLevel {
		p.updateGrowth(currentTime)
	} else {
		p.LastUpdated = currentTime
	}
	// TODO: Add death logic
}

func (p *Plant) updateWaterConsumption(currentTime time.Time) {
	elapsed := currentTime.Sub(p.LastWatered)
	days := elapsed.Hours() / 24
	waterConsumed := int(float64(p.WaterConsumptionRatePerDay) * days)
	p.WaterLevel -= waterConsumed
	if p.WaterLevel < 0 {
		p.WaterLevel = 0
	}
	p.LastWatered = currentTime
}

func (p *Plant) updateGrowth(currentTime time.Time) {
	elapsed := currentTime.Sub(p.LastUpdated)
	days := elapsed.Hours() / 24

	growthMultiplier := 1.0
	if p.WaterLevel < 50 {
		growthMultiplier = float64(p.WaterLevel) / 50
	}
	if p.WaterLevel < 20 {
		growthMultiplier = 0.0
	}

	growth := float64(p.GrowthRate) * growthMultiplier * days
	p.Growth += int64(math.Round(growth))

	p.updateGrowthStage()
	p.LastUpdated = currentTime
}

// updateGrowthStage updates the growth stage based on current growth value
func (p *Plant) updateGrowthStage() {
	stages := []GrowthStage{Maturing, Growing, Sprouting, Seeding}
	for _, stage := range stages {
		if p.Growth >= growthStageThreshold[stage] {
			p.GrowthStage = stage
			return
		}
	}
}

func (p *Plant) Image() string {
	formattedDate := time.Now().Format("2006-01-02")
	return fmt.Sprintf("%s-%s.png", formattedDate, p.Id)
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
	if p.WaterLevel+waterIncrement > 100 {
		actualToAdd = 100 - p.WaterLevel
	}
	p.WaterLevel += actualToAdd

	// Update last watered
	p.LastWatered = t

	return actualToAdd
}
