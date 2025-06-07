package plant

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func createTempVarietyFile(t *testing.T) (string, func()) {

	tempFile, err := os.CreateTemp("", "tmp-varieties-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// Write test JSON to the file
	testJSON := `{
        "bonsai": {
            "growth_rate": 1,
            "water_consumption": 2,
            "minimum_water_level": 20
        }
    }`

	if _, err := tempFile.Write([]byte(testJSON)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	return tempFile.Name(), func() {
		os.Remove(tempFile.Name())
	}

}

func TestVarietiesFromJson(t *testing.T) {
	t.Parallel()
	file, remove := createTempVarietyFile(t)
	defer remove()
	varieties, err := VarietiesFromJson(file)
	assert.NoError(t, err)
	assert.Equal(t, len(varieties), 1)
	assert.Equal(t, varieties["bonsai"].GrowthRatePerDay, int64(1))
	assert.Equal(t, varieties["bonsai"].WaterConsumptionUnitsPerDay, int64(2))
	assert.Equal(t, varieties["bonsai"].MinimumWaterLevel, int(20))
	assert.Equal(t, varieties["bonsai"].Type, "bonsai") // added field
}
