package plant

import (
	"fmt"
	"github.com/google/uuid"
	"math"
	"time"
)

const (
	maxWaterLevel = 10000 // 100.00%
)

type Plant struct {
	ID              string
	Name            string
	Type            Type
	CreationTime    time.Time
	Characteristics Characteristic

	// whether the plant can die
	CanDie bool

	// Overall Health of the plant
	Health int

	// Water related properties
	WaterLevel  int
	LastWatered time.Time

	// Growth related properties
	GrowthStage GrowthStage
	Growth      float64
	LastUpdated time.Time
}

func NewPlant(name string, plantType Type) *Plant {
	now := time.Now()
	characteristics := CharacteristicsByPlantType[plantType]

	return &Plant{
		ID:              uuid.New().String(),
		Name:            name,
		Type:            plantType,
		CreationTime:    now,
		Characteristics: characteristics,
		CanDie:          true,

		GrowthStage: Seeding,
		Growth:      0,
		LastUpdated: now,

		Health:      100,
		WaterLevel:  50,
		LastWatered: now,
	}
}

func (p *Plant) Update(currentTime time.Time) {
	p.updateWaterConsumption(currentTime)
	p.updateHealth(currentTime)
	p.upgradeGrowth(currentTime)
}

// upgradeGrowth calculates upgradeGrowth depending on elapsed time and biome
func (p *Plant) upgradeGrowth(currentTime time.Time) {
	if p.IsDead() {
		return
	}

	if p.Health <= MinimumGrowthThreshold {
		p.LastUpdated = currentTime
		return
	}

	growth := calculateGrowthRate(p.Characteristics.GrowthRate, BiomeGrowthRateModifier, p.Health, currentTime.Sub(p.LastUpdated))

	p.Growth += growth
	p.updateGrowthStage()
	p.LastUpdated = currentTime
}

func (p *Plant) Water(amount int, t time.Time) int {
	const maxWaterLevel = 10000 // 100.00%

	// Check if already at maximum capacity
	if p.WaterLevel >= maxWaterLevel {
		p.LastWatered = t
		return 0
	}

	// Calculate how much we can actually add
	spaceAvailable := maxWaterLevel - p.WaterLevel

	// Determine actual amount to add
	actualToAdd := amount
	if actualToAdd > spaceAvailable {
		actualToAdd = spaceAvailable
	}

	// Add water
	p.WaterLevel += actualToAdd

	// Safety check (shouldn't be needed with integer math, but just in case)
	if p.WaterLevel > maxWaterLevel {
		p.WaterLevel = maxWaterLevel
	}

	// Update last watered time
	p.LastWatered = t

	return actualToAdd
}

func (p *Plant) WaterLevelFormatted() float64 {
	return float64(p.WaterLevel) / 100.0
}

func (p *Plant) HealthPercent() string {
	whole := p.Health / 100
	decimal := p.Health % 100

	// Format with one decimal place as per your original format
	if decimal == 0 {
		return fmt.Sprintf("%d%%", whole)
	} else {
		// For one decimal place, divide decimal by 10 and only show the first digit
		return fmt.Sprintf("%d.%d%%", whole, decimal/10)
	}
}

func (p *Plant) WaterLevelPercent() string {
	whole := p.WaterLevel / 100
	decimal := p.WaterLevel % 100

	if decimal == 0 {
		return fmt.Sprintf("%d%%", whole)
	}
	return fmt.Sprintf("%d.%02d%%", whole, decimal)
}

// updateGrowthStage changes the upgradeGrowth stage based on current upgradeGrowth value
func (p *Plant) updateGrowthStage() {
	if p.IsDead() {
		return
	}

	stages := []GrowthStage{Maturing, Growing, Sprouting, Seeding}

	for _, stage := range stages {
		if p.Growth >= growthStageThreshold[stage] {
			p.GrowthStage = stage
			return
		}
	}
}

func calculateGrowthRate(baseGrowthRate float64, biomeGrowthRateModifier float64, health float64, elapsedDuration time.Duration) float64 {
	days := elapsedDuration.Hours() / 24
	growth := baseGrowthRate * biomeGrowthRateModifier * (health / 100) * days
	return growth
}

func (p *Plant) updateWaterConsumption(currentTime time.Time) {
	if p.IsDead() {
		return
	}

	elapsed := currentTime.Sub(p.LastWatered)
	days := elapsed.Hours() / 24
	waterConsumed := p.Characteristics.WaterConsumption * days
	p.WaterLevel -= waterConsumed
	if p.WaterLevel < 0 {
		p.WaterLevel = 0
	}
}

func (p *Plant) updateHealth(currentTime time.Time) {
	if p.IsDead() {
		return
	}

	elapsed := currentTime.Sub(p.LastWatered)
	days := elapsed.Hours() / 24

	if p.OptimalConditions() {
		if p.Health < 100 {
			p.Health += 5 * days
			if p.Health > 100 {
				p.Health = 100
			}
		}
		return
	}

	var distance, maxPossibleDistance, damageRate float64

	if p.WaterLevel < p.Characteristics.OptimalWaterMin {
		distance = p.Characteristics.OptimalWaterMin - p.WaterLevel
		maxPossibleDistance = p.Characteristics.OptimalWaterMin
		damageRate = 10.0
	}

	if p.WaterLevel > p.Characteristics.OptimalWaterMax {
		distance = p.WaterLevel - p.Characteristics.OptimalWaterMax
		maxPossibleDistance = 100 - p.Characteristics.OptimalWaterMax
		damageRate = 8.0
	}

	healthDecrease := (distance / maxPossibleDistance) * damageRate * days
	p.Health -= healthDecrease

	if p.Health < 0 {
		p.Health = 0
		if p.CanDie {
			p.GrowthStage = Dead
		}
	}
}

func (p *Plant) OptimalConditions() bool {
	return p.WaterLevel >= p.Characteristics.OptimalWaterMin &&
		p.WaterLevel <= p.Characteristics.OptimalWaterMax
}

func (p *Plant) IsDead() bool {
	if p.GrowthStage == Dead {
		return true
	}
	return false
}

func RoundToTwoDecimal(val float64) float64 {
	return math.Round(val*100) / 100
}
