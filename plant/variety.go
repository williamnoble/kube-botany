package plant

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Variety is the plant type. Ideally, we'd use `Variety` but similarity
// to in-built is playing with fire:)
type Variety struct {
	GrowthRate            int64 `json:"growth_rate"`                     // between 4-6 weeks at max growth
	WaterConsumptionUnits int64 `json:"water_requirement_units_per_day"` // 0-1 scale per day
	MinimumWaterLevel     int   `json:"minimum_water_level"`
}

type Varieties = map[string]Variety

func VarietiesFromJson() (Varieties, error) {
	filePath := filepath.Join("./plant/", "varieties.json")
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
