package plant

import (
	"encoding/json"
	"fmt"
	"os"
)

// Variety contains some characterists for a given plant type e.g., a sunflower will grow quickly and high
// water consumption level, whereas a bonsai will grow very slowly and have relatively low water consumption.
type Variety struct {
	GrowthRatePerDay            int64  `json:"growth_rate"`       // between 4-6 weeks at max growth
	WaterConsumptionUnitsPerDay int64  `json:"water_consumption"` // 0-1 scale per day
	MinimumWaterLevel           int    `json:"minimum_water_level"`
	Type                        string `json:"type,omitempty"` // duplicates key e.g. "bonsai"
}

type Varieties = map[string]Variety

// VarietiesFromJson reads a JSON file containing plant varieties and characteristics like water requirements
// and returns a map of Variety objects.
func VarietiesFromJson(filePath string) (Varieties, error) {
	var varieties Varieties

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return varieties, err
	}

	err = json.Unmarshal(fileData, &varieties)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal varieties JSON: %w", err)
	}

	// we store a pointer to Variety in the Plant type, for ease
	// create a copy of the key and store in type
	for variety, props := range varieties {
		props.Type = variety
		varieties[variety] = props
	}
	return varieties, nil
}
