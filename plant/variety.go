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
	GrowthRatePerDay            int64  `json:"growth_rate"`       // between 4-6 weeks at max growth
	WaterConsumptionUnitsPerDay int64  `json:"water_consumption"` // 0-1 scale per day
	MinimumWaterLevel           int    `json:"minimum_water_level"`
	Type                        string `json:"type,omitempty"` // duplicates key e.g. "bonsai"
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

	// we store a pointer to Variety in the Plant type, for ease
	// create a copy of the key and store in type
	for variety, props := range varieties {
		props.Type = variety
		varieties[variety] = props
	}
	return varieties, nil
}
