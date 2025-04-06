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

type Type string

var (
	Bonsai    Type = "bonsai"
	Sunflower Type = "sunflower"
)
