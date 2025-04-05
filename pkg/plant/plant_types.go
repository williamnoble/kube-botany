package plant

const (
	MinimumGrowthThreshold  = 20
	BiomeGrowthRateModifier = 1.0
)

// Type defines the type of Plant e.g. a bonsai or cactus
type Type string

var (
	Sunflower Type = "sunflower"
	Bonsai    Type = "bonsai"
	Fern      Type = "fern"
	Cactus    Type = "cactus"
	Generic   Type = "generic"
)

// GrowthStage defines the current plant upgradeGrowth stage, it
// proceeds from seeding through sprouting and growing to
// / a fully matured plant.
type GrowthStage string

const (
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
	// upgradeGrowth rate per day
	GrowthRate float64

	// watering
	WaterConsumption int
	OptimalWaterMin  int
	OptimalWaterMax  int
}

// CharacteristicsByPlantType creates a mapping of plant with
// its requirements.
var CharacteristicsByPlantType = map[Type]Characteristic{
	Bonsai: {
		GrowthRate:       22.2, // Matures in ~45 days
		OptimalWaterMin:  30,
		OptimalWaterMax:  70,
		WaterConsumption: 5,
	},
	Cactus: {
		GrowthRate:       16.7, // Matures in ~60 days
		OptimalWaterMin:  10,
		OptimalWaterMax:  30,
		WaterConsumption: 2,
	},
	Fern: {
		GrowthRate:       33.3, // Matures in ~30 days
		OptimalWaterMin:  60,
		OptimalWaterMax:  90,
		WaterConsumption: 15,
	},
	Generic: {
		GrowthRate:       33.3, // Matures in ~30 days
		OptimalWaterMin:  60,
		OptimalWaterMax:  90,
		WaterConsumption: 15,
	},
}

// growthStageThreshold defines the stage a plant is in.
var growthStageThreshold = map[GrowthStage]float64{
	Dead:      -1,
	Seeding:   0,
	Sprouting: 100,
	Growing:   500,
	Maturing:  1000,
}
