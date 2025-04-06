package plant

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// CharacteristicsByPlantType creates a mapping of plant with
// its requirements.
var CharacteristicsByPlantType = map[Type]Characteristic{
	Bonsai: {
		GrowthRate:       22, // Matures in ~45 days
		WaterConsumption: 5,  // 5% per day
	},
	Cactus: {
		GrowthRate:       16, // Matures in ~60 days
		WaterConsumption: 2,  // 2% per day
	},
	Fern: {
		GrowthRate:       33, // Matures in ~30 days
		WaterConsumption: 15, // 15% per day
	},
}

type Type string

var (
	Bonsai    Type = "bonsai"
	Fern      Type = "fern"
	Cactus    Type = "cactus"
	Sunflower Type = "sunflower"
)

type GrowthStage string

var (
	Seeding   GrowthStage = "seeding"
	Sprouting GrowthStage = "sprouting"
	Growing   GrowthStage = "growing"
	Maturing  GrowthStage = "maturing"
	Dead      GrowthStage = "dead"
)

func (g GrowthStage) String() string {
	return string(g)
}

type Characteristic struct {
	// Growth rate per day
	GrowthRate int64

	// Water consumption per day (percentage)
	WaterConsumption int64
}

// growthStageThreshold defines the stage a plant is in.
var growthStageThreshold = map[GrowthStage]int64{
	Dead:      -1,
	Seeding:   0,
	Sprouting: 100,
	Growing:   500,
	Maturing:  1000,
}

var healthThreshold = map[int64]string{
	20: "Wilting/drooping plant",
	60: "Normal healthy appearance",
	80: "Vibrant, thriving appearance",
}

type Plant struct {
	Id              string //NamespacedName
	Name            string // Fallback to ID when not provided
	Type            Type
	CreationTime    time.Time
	Characteristics Characteristic

	// Whether the plant can die
	CanDie bool

	// Water related properties
	WaterLevel        int
	MinimumWaterLevel int
	LastWatered       time.Time

	// Growth related properties
	GrowthStage GrowthStage
	Growth      int64
	LastUpdated time.Time
}

func NewPlant(
	namespacedName string,
	friendlyName string,
	plantType Type,
	canDie bool) *Plant {
	now := time.Now()
	characteristics := CharacteristicsByPlantType[plantType]

	return &Plant{
		Id:              namespacedName,
		Name:            friendlyName,
		Type:            plantType,
		CreationTime:    now,
		Characteristics: characteristics,
		CanDie:          canDie,

		GrowthStage: Seeding,
		Growth:      0,
		LastUpdated: now,

		WaterLevel:        50,
		MinimumWaterLevel: 20,
		LastWatered:       now,
	}
}

// Update progresses the plant state based on elapsed time
func (p *Plant) Update(currentTime time.Time) {
	// First update water consumption
	p.updateWaterConsumption(currentTime)

	// Update growth only if there's still water available
	if p.WaterLevel > p.MinimumWaterLevel {
		p.updateGrowth(currentTime)
	} else {
		// Just update the timestamp without growing
		p.LastUpdated = currentTime
	}
}

// Water adds 1-5% water to the plant (up to 100%)
func (p *Plant) Water(t time.Time) int {
	// Fixed water increment of 5%
	var waterIncrement = rand.Intn(5) + 1

	// Add water (capped at 100%)
	actualToAdd := waterIncrement
	if p.WaterLevel+waterIncrement > 100 {
		actualToAdd = 100 - p.WaterLevel
	}

	// Add water and update last watered time
	p.WaterLevel += actualToAdd
	p.LastWatered = t

	return actualToAdd
}

// WaterLevelPercent returns a string representation of the water level
func (p *Plant) WaterLevelPercent() string {
	return fmt.Sprintf("%d%%", p.WaterLevel)
}

func (p *Plant) MaxWatered() bool {
	return p.WaterLevel == 100
}

// updateWaterConsumption reduces water level based on elapsed time
func (p *Plant) updateWaterConsumption(currentTime time.Time) {
	elapsed := currentTime.Sub(p.LastWatered)
	days := elapsed.Hours() / 24

	// Calculate water consumed (integer percentage)
	waterConsumed := int(float64(p.Characteristics.WaterConsumption) * days)

	// Reduce water level
	p.WaterLevel -= waterConsumed
	if p.WaterLevel < 0 {
		p.WaterLevel = 0
	}

	// Update LastWatered to current time
	p.LastWatered = currentTime
}

// updateGrowth increases growth if there's water available
func (p *Plant) updateGrowth(currentTime time.Time) {
	// This function is only called when water level > 0
	elapsed := currentTime.Sub(p.LastUpdated)
	days := elapsed.Hours() / 24

	// Calculate growth based on water level
	// Full growth rate when water level is 50% or higher
	// Reduced growth rate when water level is below 50%
	growthMultiplier := 1.0
	if p.WaterLevel < 50 {
		growthMultiplier = float64(p.WaterLevel) / 50
	}

	// Calculate growth with float64 for accuracy, then convert to int64
	growth := float64(p.Characteristics.GrowthRate) * growthMultiplier * days
	p.Growth += int64(math.Floor(growth))

	// Update growth stage
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
	return fmt.Sprintf("%s-%s-.png", formattedDate, p.Id)
}
