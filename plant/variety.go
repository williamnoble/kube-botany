package plant

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Variety is the plant type. Ideally, we'd use `Variety` but similarity
// to in-built is playing with fire:)
type Variety struct {
	GrowthRate            int64 `json:"growthRate"`       // between 4-6 weeks at max growth
	WaterConsumptionUnits int64 `json:"waterRequirement"` // 0-1 scale per day
	MinimumWaterLevel     int   `json:"minimumWaterLevel"`
}

type Varieties = map[string]Variety

func VarietiesFromJson() (Varieties, error) {
	const filePath = "plants.json"
	var varieties Varieties

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return varieties, err
	}

	err = json.Unmarshal(fileData, &varieties)
	if err != nil {
		log.Fatal(err)
	}

	return varieties, nil
}
