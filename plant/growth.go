package plant

// GrowthStage describes the current stage in the growth of a plant.
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

// growthStageThreshold defines the stage of growth.
var growthStageThreshold = map[GrowthStage]int64{
	Seeding:   0,   // growthRate(5) => 10 days
	Sprouting: 50,  // growthRate(5) => 10-30 days
	Growing:   150, // growthRate(5) => 30-50 days
	Maturing:  250, // growthRate(5) => 50 days+
	Dead:      -1,
}

var healthThreshold = map[int64]string{
	20: "Wilting/drooping plant",
	60: "Normal healthy appearance",
	80: "Vibrant, thriving appearance",
}
